package terraform

import (
	"reflect"
	"strings"
	"testing"
)

func TestGraph(t *testing.T) {
	config := testConfig(t, "graph-basic")

	g, err := Graph(&GraphOpts{Config: config})
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	actual := strings.TrimSpace(g.String())
	expected := strings.TrimSpace(testTerraformGraphStr)
	if actual != expected {
		t.Fatalf("bad:\n\n%s", actual)
	}
}

func TestGraph_configRequired(t *testing.T) {
	if _, err := Graph(new(GraphOpts)); err == nil {
		t.Fatal("should error")
	}
}

func TestGraph_cycle(t *testing.T) {
	config := testConfig(t, "graph-cycle")

	_, err := Graph(&GraphOpts{Config: config})
	if err == nil {
		t.Fatal("should error")
	}
}

func TestGraph_state(t *testing.T) {
	config := testConfig(t, "graph-basic")
	state := &State{
		Resources: map[string]*ResourceState{
			"aws_instance.old": &ResourceState{
				ID:   "foo",
				Type: "aws_instance",
			},
		},
	}

	g, err := Graph(&GraphOpts{Config: config, State: state})
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	actual := strings.TrimSpace(g.String())
	expected := strings.TrimSpace(testTerraformGraphStateStr)
	if actual != expected {
		t.Fatalf("bad:\n\n%s", actual)
	}
}

func TestGraphFull(t *testing.T) {
	rpAws := new(MockResourceProvider)
	rpOS := new(MockResourceProvider)

	rpAws.ResourcesReturn = []ResourceType{
		ResourceType{Name: "aws_instance"},
		ResourceType{Name: "aws_load_balancer"},
		ResourceType{Name: "aws_security_group"},
	}
	rpOS.ResourcesReturn = []ResourceType{
		ResourceType{Name: "openstack_floating_ip"},
	}

	ps := map[string]ResourceProviderFactory{
		"aws":  testProviderFuncFixed(rpAws),
		"open": testProviderFuncFixed(rpOS),
	}

	c := testConfig(t, "graph-basic")
	g, err := Graph(&GraphOpts{Config: c, Providers: ps})
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// A helper to help get us the provider for a resource.
	graphProvider := func(n string) ResourceProvider {
		return g.Noun(n).Meta.(*GraphNodeResource).Resource.Provider
	}

	// Test a couple
	if graphProvider("aws_instance.web") != rpAws {
		t.Fatalf("bad: %#v", graphProvider("aws_instance.web"))
	}
	if graphProvider("openstack_floating_ip.random") != rpOS {
		t.Fatalf("bad: %#v", graphProvider("openstack_floating_ip.random"))
	}

	// Test that all providers have been set
	for _, n := range g.Nouns {
		switch m := n.Meta.(type) {
		case *GraphNodeResource:
			if m.Resource.Provider == nil {
				t.Fatalf("bad: %#v", m)
			}
		case *GraphNodeResourceProvider:
			if len(m.Providers) == 0 {
				t.Fatalf("bad: %#v", m)
			}
		default:
			continue
		}
	}
}

func TestGraphAddDiff(t *testing.T) {
	config := testConfig(t, "graph-diff")
	diff := &Diff{
		Resources: map[string]*ResourceDiff{
			"aws_instance.foo": &ResourceDiff{
				Attributes: map[string]*ResourceAttrDiff{
					"foo": &ResourceAttrDiff{
						New: "bar",
					},
				},
			},
		},
	}

	g, err := Graph(&GraphOpts{Config: config, Diff: diff})
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	actual := strings.TrimSpace(g.String())
	expected := strings.TrimSpace(testTerraformGraphDiffStr)
	if actual != expected {
		t.Fatalf("bad:\n\n%s", actual)
	}

	// Verify that the state has been added
	n := g.Noun("aws_instance.foo")
	rn := n.Meta.(*GraphNodeResource)

	expected2 := diff.Resources["aws_instance.foo"]
	actual2 := rn.Resource.Diff
	if !reflect.DeepEqual(actual2, expected2) {
		t.Fatalf("bad: %#v", actual2)
	}
}

func TestGraphAddDiff_destroy(t *testing.T) {
	config := testConfig(t, "graph-diff-destroy")
	diff := &Diff{
		Resources: map[string]*ResourceDiff{
			"aws_instance.foo": &ResourceDiff{
				Destroy: true,
			},
			"aws_instance.bar": &ResourceDiff{
				Destroy: true,
			},
		},
	}
	state := &State{
		Resources: map[string]*ResourceState{
			"aws_instance.foo": &ResourceState{
				ID:   "foo",
				Type: "aws_instance",
			},

			"aws_instance.bar": &ResourceState{
				ID:   "bar",
				Type: "aws_instance",
				Dependencies: []ResourceDependency{
					ResourceDependency{
						ID: "foo",
					},
				},
			},
		},
	}

	g, err := Graph(&GraphOpts{
		Config: config,
		Diff:   diff,
		State:  state,
	})
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	actual := strings.TrimSpace(g.String())
	expected := strings.TrimSpace(testTerraformGraphDiffDestroyStr)
	if actual != expected {
		t.Fatalf("bad:\n\n%s", actual)
	}

	// Verify that the state has been added
	n := g.Noun("aws_instance.foo")
	rn := n.Meta.(*GraphNodeResource)

	expected2 := diff.Resources["aws_instance.foo"]
	actual2 := rn.Resource.Diff
	if !reflect.DeepEqual(actual2, expected2) {
		t.Fatalf("bad: %#v", actual2)
	}
}

const testTerraformGraphStr = `
root: root
aws_instance.web
  aws_instance.web -> aws_security_group.firewall
  aws_instance.web -> provider.aws
aws_load_balancer.weblb
  aws_load_balancer.weblb -> aws_instance.web
  aws_load_balancer.weblb -> provider.aws
aws_security_group.firewall
  aws_security_group.firewall -> provider.aws
openstack_floating_ip.random
provider.aws
  provider.aws -> openstack_floating_ip.random
root
  root -> aws_instance.web
  root -> aws_load_balancer.weblb
  root -> aws_security_group.firewall
  root -> openstack_floating_ip.random
  root -> provider.aws
`

const testTerraformGraphDiffStr = `
root: root
aws_instance.foo
root
  root -> aws_instance.foo
`

const testTerraformGraphDiffDestroyStr = `
root: root
aws_instance.bar
  aws_instance.bar -> aws_instance.bar (destroy)
aws_instance.bar (destroy)
aws_instance.foo
  aws_instance.foo -> aws_instance.foo (destroy)
aws_instance.foo (destroy)
  aws_instance.foo (destroy) -> aws_instance.bar (destroy)
root
  root -> aws_instance.bar
  root -> aws_instance.foo
`

const testTerraformGraphStateStr = `
root: root
aws_instance.old
  aws_instance.old -> provider.aws
aws_instance.web
  aws_instance.web -> aws_security_group.firewall
  aws_instance.web -> provider.aws
aws_load_balancer.weblb
  aws_load_balancer.weblb -> aws_instance.web
  aws_load_balancer.weblb -> provider.aws
aws_security_group.firewall
  aws_security_group.firewall -> provider.aws
openstack_floating_ip.random
provider.aws
  provider.aws -> openstack_floating_ip.random
root
  root -> aws_instance.old
  root -> aws_instance.web
  root -> aws_load_balancer.weblb
  root -> aws_security_group.firewall
  root -> openstack_floating_ip.random
  root -> provider.aws
`
