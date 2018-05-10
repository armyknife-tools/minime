package terraform

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform/addrs"
)

func TestRefreshGraphBuilder_configOrphans(t *testing.T) {

	m := testModule(t, "refresh-config-orphan")

	state := &State{
		Modules: []*ModuleState{
			&ModuleState{
				Path: rootModulePath,
				Resources: map[string]*ResourceState{
					"test_object.foo.0": &ResourceState{
						Type: "test_object",
						Deposed: []*InstanceState{
							&InstanceState{
								ID: "foo",
							},
						},
					},
					"test_object.foo.1": &ResourceState{
						Type: "test_object",
						Deposed: []*InstanceState{
							&InstanceState{
								ID: "bar",
							},
						},
					},
					"test_object.foo.2": &ResourceState{
						Type: "test_object",
						Deposed: []*InstanceState{
							&InstanceState{
								ID: "baz",
							},
						},
					},
					"data.test_object.foo.0": &ResourceState{
						Type: "test_object",
						Deposed: []*InstanceState{
							&InstanceState{
								ID: "foo",
							},
						},
					},
					"data.test_object.foo.1": &ResourceState{
						Type: "test_object",
						Deposed: []*InstanceState{
							&InstanceState{
								ID: "bar",
							},
						},
					},
					"data.test_object.foo.2": &ResourceState{
						Type: "test_object",
						Deposed: []*InstanceState{
							&InstanceState{
								ID: "baz",
							},
						},
					},
				},
			},
		},
	}

	b := &RefreshGraphBuilder{
		Config:     m,
		State:      state,
		Components: simpleMockComponentFactory(),
	}
	g, err := b.Build(addrs.RootModuleInstance)
	if err != nil {
		t.Fatalf("Error building graph: %s", err)
	}

	actual := strings.TrimSpace(g.StringWithNodeTypes())
	expected := strings.TrimSpace(`
data.test_object.foo[0] - *terraform.NodeRefreshableManagedResourceInstance
  provider.test - *terraform.NodeApplyableProvider
data.test_object.foo[1] - *terraform.NodeRefreshableManagedResourceInstance
  provider.test - *terraform.NodeApplyableProvider
data.test_object.foo[2] - *terraform.NodeRefreshableManagedResourceInstance
  provider.test - *terraform.NodeApplyableProvider
provider.test - *terraform.NodeApplyableProvider
provider.test (close) - *terraform.graphNodeCloseProvider
  data.test_object.foo[0] - *terraform.NodeRefreshableManagedResourceInstance
  data.test_object.foo[1] - *terraform.NodeRefreshableManagedResourceInstance
  data.test_object.foo[2] - *terraform.NodeRefreshableManagedResourceInstance
  test_object.foo - *terraform.NodeRefreshableManagedResource
test_object.foo - *terraform.NodeRefreshableManagedResource
  provider.test - *terraform.NodeApplyableProvider
`)
	if expected != actual {
		t.Fatalf("wrong result\n\ngot:\n%s\n\nwant:\n%s", actual, expected)
	}
}
