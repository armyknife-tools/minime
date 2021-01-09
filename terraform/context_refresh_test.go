package terraform

import (
	"reflect"
	"regexp"
	"sort"
	"strings"
	"sync"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/zclconf/go-cty/cty"

	"github.com/hashicorp/terraform/addrs"
	"github.com/hashicorp/terraform/configs/configschema"
	"github.com/hashicorp/terraform/configs/hcl2shim"
	"github.com/hashicorp/terraform/providers"
	"github.com/hashicorp/terraform/states"
)

func TestContext2Refresh(t *testing.T) {
	p := testProvider("aws")
	m := testModule(t, "refresh-basic")

	state := states.NewState()
	root := state.EnsureModule(addrs.RootModuleInstance)
	root.SetResourceInstanceCurrent(
		mustResourceInstanceAddr("aws_instance.web").Resource,
		&states.ResourceInstanceObjectSrc{
			Status:    states.ObjectReady,
			AttrsJSON: []byte(`{"id":"foo","foo":"bar"}`),
		},
		mustProviderConfig(`provider["registry.terraform.io/hashicorp/aws"]`),
	)

	ctx := testContext2(t, &ContextOpts{
		Config: m,
		Providers: map[addrs.Provider]providers.Factory{
			addrs.NewDefaultProvider("aws"): testProviderFuncFixed(p),
		},
		State: state,
	})

	schema := p.GetSchemaReturn.ResourceTypes["aws_instance"]
	ty := schema.ImpliedType()
	readState, err := hcl2shim.HCL2ValueFromFlatmap(map[string]string{"id": "foo", "foo": "baz"}, ty)
	if err != nil {
		t.Fatal(err)
	}

	p.ReadResourceFn = nil
	p.ReadResourceResponse = &providers.ReadResourceResponse{
		NewState: readState,
	}
	p.PlanResourceChangeFn = testDiffFn

	s, diags := ctx.Refresh()
	if diags.HasErrors() {
		t.Fatal(diags.Err())
	}

	if !p.ReadResourceCalled {
		t.Fatal("ReadResource should be called")
	}

	mod := s.RootModule()
	fromState, err := mod.Resources["aws_instance.web"].Instances[addrs.NoKey].Current.Decode(ty)
	if err != nil {
		t.Fatal(err)
	}

	newState, err := schema.CoerceValue(fromState.Value)
	if err != nil {
		t.Fatal(err)
	}

	if !cmp.Equal(readState, newState, valueComparer) {
		t.Fatal(cmp.Diff(readState, newState, valueComparer, equateEmpty))
	}
}

func TestContext2Refresh_dynamicAttr(t *testing.T) {
	m := testModule(t, "refresh-dynamic")

	startingState := states.BuildState(func(ss *states.SyncState) {
		ss.SetResourceInstanceCurrent(
			addrs.Resource{
				Mode: addrs.ManagedResourceMode,
				Type: "test_instance",
				Name: "foo",
			}.Instance(addrs.NoKey).Absolute(addrs.RootModuleInstance),
			&states.ResourceInstanceObjectSrc{
				Status:    states.ObjectReady,
				AttrsJSON: []byte(`{"dynamic":{"type":"string","value":"hello"}}`),
			},
			addrs.AbsProviderConfig{
				Provider: addrs.NewDefaultProvider("test"),
				Module:   addrs.RootModule,
			},
		)
	})

	readStateVal := cty.ObjectVal(map[string]cty.Value{
		"dynamic": cty.EmptyTupleVal,
	})

	p := testProvider("test")
	p.GetSchemaReturn = &ProviderSchema{
		ResourceTypes: map[string]*configschema.Block{
			"test_instance": {
				Attributes: map[string]*configschema.Attribute{
					"dynamic": {Type: cty.DynamicPseudoType, Optional: true},
				},
			},
		},
	}
	p.ReadResourceFn = func(req providers.ReadResourceRequest) providers.ReadResourceResponse {
		return providers.ReadResourceResponse{
			NewState: readStateVal,
		}
	}
	p.PlanResourceChangeFn = func(req providers.PlanResourceChangeRequest) (resp providers.PlanResourceChangeResponse) {
		resp.PlannedState = req.ProposedNewState
		return resp
	}

	ctx := testContext2(t, &ContextOpts{
		Config: m,
		Providers: map[addrs.Provider]providers.Factory{
			addrs.NewDefaultProvider("test"): testProviderFuncFixed(p),
		},
		State: startingState,
	})

	schema := p.GetSchemaReturn.ResourceTypes["test_instance"]
	ty := schema.ImpliedType()

	s, diags := ctx.Refresh()
	if diags.HasErrors() {
		t.Fatal(diags.Err())
	}

	if !p.ReadResourceCalled {
		t.Fatal("ReadResource should be called")
	}

	mod := s.RootModule()
	newState, err := mod.Resources["test_instance.foo"].Instances[addrs.NoKey].Current.Decode(ty)
	if err != nil {
		t.Fatal(err)
	}

	if !cmp.Equal(readStateVal, newState.Value, valueComparer) {
		t.Error(cmp.Diff(newState.Value, readStateVal, valueComparer, equateEmpty))
	}
}

func TestContext2Refresh_dataComputedModuleVar(t *testing.T) {
	p := testProvider("aws")
	m := testModule(t, "refresh-data-module-var")
	p.PlanResourceChangeFn = func(req providers.PlanResourceChangeRequest) (resp providers.PlanResourceChangeResponse) {
		obj := req.ProposedNewState.AsValueMap()
		obj["id"] = cty.UnknownVal(cty.String)
		resp.PlannedState = cty.ObjectVal(obj)
		return resp
	}
	p.ReadDataSourceFn = func(req providers.ReadDataSourceRequest) (resp providers.ReadDataSourceResponse) {
		resp.State = req.Config
		return resp
	}

	p.GetSchemaReturn = &ProviderSchema{
		Provider: &configschema.Block{},
		ResourceTypes: map[string]*configschema.Block{
			"aws_instance": {
				Attributes: map[string]*configschema.Attribute{
					"foo": {
						Type:     cty.String,
						Optional: true,
					},
					"id": {
						Type:     cty.String,
						Computed: true,
					},
				},
			},
		},
		DataSources: map[string]*configschema.Block{
			"aws_data_source": {
				Attributes: map[string]*configschema.Attribute{
					"id": {
						Type:     cty.String,
						Optional: true,
					},
					"output": {
						Type:     cty.String,
						Computed: true,
					},
				},
			},
		},
	}

	ctx := testContext2(t, &ContextOpts{
		Config: m,
		Providers: map[addrs.Provider]providers.Factory{
			addrs.NewDefaultProvider("aws"): testProviderFuncFixed(p),
		},
	})

	s, diags := ctx.Refresh()
	if diags.HasErrors() {
		t.Fatalf("refresh errors: %s", diags.Err())
	}

	checkStateString(t, s, `
<no state>
`)
}

func TestContext2Refresh_targeted(t *testing.T) {
	p := testProvider("aws")
	p.GetSchemaReturn = &ProviderSchema{
		Provider: &configschema.Block{},
		ResourceTypes: map[string]*configschema.Block{
			"aws_elb": {
				Attributes: map[string]*configschema.Attribute{
					"instances": {
						Type:     cty.Set(cty.String),
						Optional: true,
					},
				},
			},
			"aws_instance": {
				Attributes: map[string]*configschema.Attribute{
					"id": {
						Type:     cty.String,
						Computed: true,
					},
					"vpc_id": {
						Type:     cty.String,
						Optional: true,
					},
				},
			},
			"aws_vpc": {
				Attributes: map[string]*configschema.Attribute{
					"id": {
						Type:     cty.String,
						Computed: true,
					},
				},
			},
		},
	}

	state := states.NewState()
	root := state.EnsureModule(addrs.RootModuleInstance)
	testSetResourceInstanceCurrent(root, "aws_vpc.metoo", `{"id":"vpc-abc123"}`, `provider["registry.terraform.io/hashicorp/aws"]`)
	testSetResourceInstanceCurrent(root, "aws_instance.notme", `{"id":"i-bcd345"}`, `provider["registry.terraform.io/hashicorp/aws"]`)
	testSetResourceInstanceCurrent(root, "aws_instance.me", `{"id":"i-abc123"}`, `provider["registry.terraform.io/hashicorp/aws"]`)
	testSetResourceInstanceCurrent(root, "aws_elb.meneither", `{"id":"lb-abc123"}`, `provider["registry.terraform.io/hashicorp/aws"]`)

	m := testModule(t, "refresh-targeted")
	ctx := testContext2(t, &ContextOpts{
		Config: m,
		Providers: map[addrs.Provider]providers.Factory{
			addrs.NewDefaultProvider("aws"): testProviderFuncFixed(p),
		},
		State: state,
		Targets: []addrs.Targetable{
			addrs.RootModuleInstance.Resource(
				addrs.ManagedResourceMode, "aws_instance", "me",
			),
		},
	})

	refreshedResources := make([]string, 0, 2)
	p.ReadResourceFn = func(req providers.ReadResourceRequest) providers.ReadResourceResponse {
		refreshedResources = append(refreshedResources, req.PriorState.GetAttr("id").AsString())
		return providers.ReadResourceResponse{
			NewState: req.PriorState,
		}
	}
	p.PlanResourceChangeFn = testDiffFn

	_, diags := ctx.Refresh()
	if diags.HasErrors() {
		t.Fatalf("refresh errors: %s", diags.Err())
	}

	expected := []string{"vpc-abc123", "i-abc123"}
	if !reflect.DeepEqual(refreshedResources, expected) {
		t.Fatalf("expected: %#v, got: %#v", expected, refreshedResources)
	}
}

func TestContext2Refresh_targetedCount(t *testing.T) {
	p := testProvider("aws")
	p.GetSchemaReturn = &ProviderSchema{
		Provider: &configschema.Block{},
		ResourceTypes: map[string]*configschema.Block{
			"aws_elb": {
				Attributes: map[string]*configschema.Attribute{
					"instances": {
						Type:     cty.Set(cty.String),
						Optional: true,
					},
				},
			},
			"aws_instance": {
				Attributes: map[string]*configschema.Attribute{
					"id": {
						Type:     cty.String,
						Computed: true,
					},
					"vpc_id": {
						Type:     cty.String,
						Optional: true,
					},
				},
			},
			"aws_vpc": {
				Attributes: map[string]*configschema.Attribute{
					"id": {
						Type:     cty.String,
						Computed: true,
					},
				},
			},
		},
	}

	state := states.NewState()
	root := state.EnsureModule(addrs.RootModuleInstance)
	testSetResourceInstanceCurrent(root, "aws_vpc.metoo", `{"id":"vpc-abc123"}`, `provider["registry.terraform.io/hashicorp/aws"]`)
	testSetResourceInstanceCurrent(root, "aws_instance.notme", `{"id":"i-bcd345"}`, `provider["registry.terraform.io/hashicorp/aws"]`)
	testSetResourceInstanceCurrent(root, "aws_instance.me[0]", `{"id":"i-abc123"}`, `provider["registry.terraform.io/hashicorp/aws"]`)
	testSetResourceInstanceCurrent(root, "aws_instance.me[1]", `{"id":"i-cde567"}`, `provider["registry.terraform.io/hashicorp/aws"]`)
	testSetResourceInstanceCurrent(root, "aws_instance.me[2]", `{"id":"i-cde789"}`, `provider["registry.terraform.io/hashicorp/aws"]`)
	testSetResourceInstanceCurrent(root, "aws_elb.meneither", `{"id":"lb-abc123"}`, `provider["registry.terraform.io/hashicorp/aws"]`)

	m := testModule(t, "refresh-targeted-count")
	ctx := testContext2(t, &ContextOpts{
		Config: m,
		Providers: map[addrs.Provider]providers.Factory{
			addrs.NewDefaultProvider("aws"): testProviderFuncFixed(p),
		},
		State: state,
		Targets: []addrs.Targetable{
			addrs.RootModuleInstance.Resource(
				addrs.ManagedResourceMode, "aws_instance", "me",
			),
		},
	})

	refreshedResources := make([]string, 0, 2)
	p.ReadResourceFn = func(req providers.ReadResourceRequest) providers.ReadResourceResponse {
		refreshedResources = append(refreshedResources, req.PriorState.GetAttr("id").AsString())
		return providers.ReadResourceResponse{
			NewState: req.PriorState,
		}
	}
	p.PlanResourceChangeFn = testDiffFn

	_, diags := ctx.Refresh()
	if diags.HasErrors() {
		t.Fatalf("refresh errors: %s", diags.Err())
	}

	// Target didn't specify index, so we should get all our instances
	expected := []string{
		"vpc-abc123",
		"i-abc123",
		"i-cde567",
		"i-cde789",
	}
	sort.Strings(expected)
	sort.Strings(refreshedResources)
	if !reflect.DeepEqual(refreshedResources, expected) {
		t.Fatalf("wrong result\ngot:  %#v\nwant: %#v", refreshedResources, expected)
	}
}

func TestContext2Refresh_targetedCountIndex(t *testing.T) {
	p := testProvider("aws")
	p.GetSchemaReturn = &ProviderSchema{
		Provider: &configschema.Block{},
		ResourceTypes: map[string]*configschema.Block{
			"aws_elb": {
				Attributes: map[string]*configschema.Attribute{
					"instances": {
						Type:     cty.Set(cty.String),
						Optional: true,
					},
				},
			},
			"aws_instance": {
				Attributes: map[string]*configschema.Attribute{
					"id": {
						Type:     cty.String,
						Computed: true,
					},
					"vpc_id": {
						Type:     cty.String,
						Optional: true,
					},
				},
			},
			"aws_vpc": {
				Attributes: map[string]*configschema.Attribute{
					"id": {
						Type:     cty.String,
						Computed: true,
					},
				},
			},
		},
	}

	state := states.NewState()
	root := state.EnsureModule(addrs.RootModuleInstance)
	testSetResourceInstanceCurrent(root, "aws_vpc.metoo", `{"id":"vpc-abc123"}`, `provider["registry.terraform.io/hashicorp/aws"]`)
	testSetResourceInstanceCurrent(root, "aws_instance.notme", `{"id":"i-bcd345"}`, `provider["registry.terraform.io/hashicorp/aws"]`)
	testSetResourceInstanceCurrent(root, "aws_instance.me[0]", `{"id":"i-abc123"}`, `provider["registry.terraform.io/hashicorp/aws"]`)
	testSetResourceInstanceCurrent(root, "aws_instance.me[1]", `{"id":"i-cde567"}`, `provider["registry.terraform.io/hashicorp/aws"]`)
	testSetResourceInstanceCurrent(root, "aws_instance.me[2]", `{"id":"i-cde789"}`, `provider["registry.terraform.io/hashicorp/aws"]`)
	testSetResourceInstanceCurrent(root, "aws_elb.meneither", `{"id":"lb-abc123"}`, `provider["registry.terraform.io/hashicorp/aws"]`)

	m := testModule(t, "refresh-targeted-count")
	ctx := testContext2(t, &ContextOpts{
		Config: m,
		Providers: map[addrs.Provider]providers.Factory{
			addrs.NewDefaultProvider("aws"): testProviderFuncFixed(p),
		},
		State: state,
		Targets: []addrs.Targetable{
			addrs.RootModuleInstance.ResourceInstance(
				addrs.ManagedResourceMode, "aws_instance", "me", addrs.IntKey(0),
			),
		},
	})

	refreshedResources := make([]string, 0, 2)
	p.ReadResourceFn = func(req providers.ReadResourceRequest) providers.ReadResourceResponse {
		refreshedResources = append(refreshedResources, req.PriorState.GetAttr("id").AsString())
		return providers.ReadResourceResponse{
			NewState: req.PriorState,
		}
	}
	p.PlanResourceChangeFn = testDiffFn

	_, diags := ctx.Refresh()
	if diags.HasErrors() {
		t.Fatalf("refresh errors: %s", diags.Err())
	}

	expected := []string{"vpc-abc123", "i-abc123"}
	if !reflect.DeepEqual(refreshedResources, expected) {
		t.Fatalf("wrong result\ngot:  %#v\nwant: %#v", refreshedResources, expected)
	}
}

func TestContext2Refresh_moduleComputedVar(t *testing.T) {
	p := testProvider("aws")
	p.GetSchemaReturn = &ProviderSchema{
		Provider: &configschema.Block{},
		ResourceTypes: map[string]*configschema.Block{
			"aws_instance": {
				Attributes: map[string]*configschema.Attribute{
					"id": {
						Type:     cty.String,
						Computed: true,
					},
					"value": {
						Type:     cty.String,
						Optional: true,
					},
				},
			},
		},
	}
	p.PlanResourceChangeFn = testDiffFn

	m := testModule(t, "refresh-module-computed-var")
	ctx := testContext2(t, &ContextOpts{
		Config: m,
		Providers: map[addrs.Provider]providers.Factory{
			addrs.NewDefaultProvider("aws"): testProviderFuncFixed(p),
		},
	})

	// This was failing (see GH-2188) at some point, so this test just
	// verifies that the failure goes away.
	if _, diags := ctx.Refresh(); diags.HasErrors() {
		t.Fatalf("refresh errs: %s", diags.Err())
	}
}

func TestContext2Refresh_delete(t *testing.T) {
	p := testProvider("aws")
	m := testModule(t, "refresh-basic")

	state := states.NewState()
	root := state.EnsureModule(addrs.RootModuleInstance)
	testSetResourceInstanceCurrent(root, "aws_instance.web", `{"id":"foo"}`, `provider["registry.terraform.io/hashicorp/aws"]`)

	ctx := testContext2(t, &ContextOpts{
		Config: m,
		Providers: map[addrs.Provider]providers.Factory{
			addrs.NewDefaultProvider("aws"): testProviderFuncFixed(p),
		},
		State: state,
	})

	p.ReadResourceFn = nil
	p.ReadResourceResponse = &providers.ReadResourceResponse{
		NewState: cty.NullVal(p.GetSchemaReturn.ResourceTypes["aws_instance"].ImpliedType()),
	}
	p.PlanResourceChangeFn = testDiffFn

	s, diags := ctx.Refresh()
	if diags.HasErrors() {
		t.Fatalf("refresh errors: %s", diags.Err())
	}

	mod := s.RootModule()
	if len(mod.Resources) > 0 {
		t.Fatal("resources should be empty")
	}
}

func TestContext2Refresh_ignoreUncreated(t *testing.T) {
	p := testProvider("aws")
	m := testModule(t, "refresh-basic")
	ctx := testContext2(t, &ContextOpts{
		Config: m,
		Providers: map[addrs.Provider]providers.Factory{
			addrs.NewDefaultProvider("aws"): testProviderFuncFixed(p),
		},
		State: nil,
	})

	p.ReadResourceFn = nil
	p.ReadResourceResponse = &providers.ReadResourceResponse{
		NewState: cty.ObjectVal(map[string]cty.Value{
			"id": cty.StringVal("foo"),
		}),
	}
	p.PlanResourceChangeFn = testDiffFn

	_, diags := ctx.Refresh()
	if diags.HasErrors() {
		t.Fatalf("refresh errors: %s", diags.Err())
	}
	if p.ReadResourceCalled {
		t.Fatal("refresh should not be called")
	}
}

func TestContext2Refresh_hook(t *testing.T) {
	h := new(MockHook)
	p := testProvider("aws")
	p.PlanResourceChangeFn = testDiffFn
	m := testModule(t, "refresh-basic")

	state := states.NewState()
	root := state.EnsureModule(addrs.RootModuleInstance)
	testSetResourceInstanceCurrent(root, "aws_instance.web", `{"id":"foo"}`, `provider["registry.terraform.io/hashicorp/aws"]`)

	ctx := testContext2(t, &ContextOpts{
		Config: m,
		Hooks:  []Hook{h},
		Providers: map[addrs.Provider]providers.Factory{
			addrs.NewDefaultProvider("aws"): testProviderFuncFixed(p),
		},
		State: state,
	})

	if _, diags := ctx.Refresh(); diags.HasErrors() {
		t.Fatalf("refresh errs: %s", diags.Err())
	}
	if !h.PreRefreshCalled {
		t.Fatal("should be called")
	}
	if !h.PostRefreshCalled {
		t.Fatal("should be called")
	}
}

func TestContext2Refresh_modules(t *testing.T) {
	p := testProvider("aws")
	m := testModule(t, "refresh-modules")

	state := states.NewState()
	root := state.EnsureModule(addrs.RootModuleInstance)
	testSetResourceInstanceTainted(root, "aws_instance.web", `{"id":"bar"}`, `provider["registry.terraform.io/hashicorp/aws"]`)
	child := state.EnsureModule(addrs.RootModuleInstance.Child("child", addrs.NoKey))
	testSetResourceInstanceCurrent(child, "aws_instance.web", `{"id":"baz"}`, `provider["registry.terraform.io/hashicorp/aws"]`)

	ctx := testContext2(t, &ContextOpts{
		Config: m,
		Providers: map[addrs.Provider]providers.Factory{
			addrs.NewDefaultProvider("aws"): testProviderFuncFixed(p),
		},
		State: state,
	})

	p.ReadResourceFn = func(req providers.ReadResourceRequest) providers.ReadResourceResponse {
		if !req.PriorState.GetAttr("id").RawEquals(cty.StringVal("baz")) {
			return providers.ReadResourceResponse{
				NewState: req.PriorState,
			}
		}

		new, _ := cty.Transform(req.PriorState, func(path cty.Path, v cty.Value) (cty.Value, error) {
			if len(path) == 1 && path[0].(cty.GetAttrStep).Name == "id" {
				return cty.StringVal("new"), nil
			}
			return v, nil
		})
		return providers.ReadResourceResponse{
			NewState: new,
		}
	}
	p.PlanResourceChangeFn = testDiffFn

	s, diags := ctx.Refresh()
	if diags.HasErrors() {
		t.Fatalf("refresh errors: %s", diags.Err())
	}

	actual := strings.TrimSpace(s.String())
	expected := strings.TrimSpace(testContextRefreshModuleStr)
	if actual != expected {
		t.Fatalf("wrong result\n\ngot:\n%s\n\nwant:\n%s", actual, expected)
	}
}

func TestContext2Refresh_moduleInputComputedOutput(t *testing.T) {
	m := testModule(t, "refresh-module-input-computed-output")
	p := testProvider("aws")
	p.PlanResourceChangeFn = testDiffFn
	p.GetSchemaReturn = &ProviderSchema{
		Provider: &configschema.Block{},
		ResourceTypes: map[string]*configschema.Block{
			"aws_instance": {
				Attributes: map[string]*configschema.Attribute{
					"foo": {
						Type:     cty.String,
						Optional: true,
						Computed: true,
					},
					"compute": {
						Type:     cty.String,
						Optional: true,
					},
				},
			},
		},
	}

	ctx := testContext2(t, &ContextOpts{
		Config: m,
		Providers: map[addrs.Provider]providers.Factory{
			addrs.NewDefaultProvider("aws"): testProviderFuncFixed(p),
		},
	})

	if _, diags := ctx.Refresh(); diags.HasErrors() {
		t.Fatalf("refresh errs: %s", diags.Err())
	}
}

func TestContext2Refresh_moduleVarModule(t *testing.T) {
	m := testModule(t, "refresh-module-var-module")
	p := testProvider("aws")
	p.PlanResourceChangeFn = testDiffFn
	ctx := testContext2(t, &ContextOpts{
		Config: m,
		Providers: map[addrs.Provider]providers.Factory{
			addrs.NewDefaultProvider("aws"): testProviderFuncFixed(p),
		},
	})

	if _, diags := ctx.Refresh(); diags.HasErrors() {
		t.Fatalf("refresh errs: %s", diags.Err())
	}
}

// GH-70
func TestContext2Refresh_noState(t *testing.T) {
	p := testProvider("aws")
	m := testModule(t, "refresh-no-state")
	ctx := testContext2(t, &ContextOpts{
		Config: m,
		Providers: map[addrs.Provider]providers.Factory{
			addrs.NewDefaultProvider("aws"): testProviderFuncFixed(p),
		},
	})

	p.ReadResourceFn = nil
	p.ReadResourceResponse = &providers.ReadResourceResponse{
		NewState: cty.ObjectVal(map[string]cty.Value{
			"id": cty.StringVal("foo"),
		}),
	}
	p.PlanResourceChangeFn = testDiffFn

	if _, diags := ctx.Refresh(); diags.HasErrors() {
		t.Fatalf("refresh errs: %s", diags.Err())
	}
}

func TestContext2Refresh_output(t *testing.T) {
	p := testProvider("aws")
	p.GetSchemaReturn = &ProviderSchema{
		Provider: &configschema.Block{},
		ResourceTypes: map[string]*configschema.Block{
			"aws_instance": {
				Attributes: map[string]*configschema.Attribute{
					"id": {
						Type:     cty.String,
						Computed: true,
					},
					"foo": {
						Type:     cty.String,
						Optional: true,
						Computed: true,
					},
				},
			},
		},
	}
	p.PlanResourceChangeFn = testDiffFn

	m := testModule(t, "refresh-output")

	state := states.NewState()
	root := state.EnsureModule(addrs.RootModuleInstance)
	testSetResourceInstanceCurrent(root, "aws_instance.web", `{"id":"foo","foo":"bar"}`, `provider["registry.terraform.io/hashicorp/aws"]`)
	root.SetOutputValue("foo", cty.StringVal("foo"), false)

	ctx := testContext2(t, &ContextOpts{
		Config: m,
		Providers: map[addrs.Provider]providers.Factory{
			addrs.NewDefaultProvider("aws"): testProviderFuncFixed(p),
		},
		State: state,
	})

	s, diags := ctx.Refresh()
	if diags.HasErrors() {
		t.Fatalf("refresh errors: %s", diags.Err())
	}

	actual := strings.TrimSpace(s.String())
	expected := strings.TrimSpace(testContextRefreshOutputStr)
	if actual != expected {
		t.Fatalf("wrong result\n\ngot:\n%q\n\nwant:\n%q", actual, expected)
	}
}

func TestContext2Refresh_outputPartial(t *testing.T) {
	p := testProvider("aws")
	m := testModule(t, "refresh-output-partial")

	// Refresh creates a partial plan for any instances that don't have
	// remote objects yet, to get stub values for interpolation. Therefore
	// we need to make DiffFn available to let that complete.
	p.PlanResourceChangeFn = testDiffFn

	p.GetSchemaReturn = &ProviderSchema{
		Provider: &configschema.Block{},
		ResourceTypes: map[string]*configschema.Block{
			"aws_instance": {
				Attributes: map[string]*configschema.Attribute{
					"foo": {
						Type:     cty.String,
						Computed: true,
					},
				},
			},
		},
	}

	p.ReadResourceFn = nil
	p.ReadResourceResponse = &providers.ReadResourceResponse{
		NewState: cty.NullVal(p.GetSchemaReturn.ResourceTypes["aws_instance"].ImpliedType()),
	}

	state := states.NewState()
	root := state.EnsureModule(addrs.RootModuleInstance)
	testSetResourceInstanceCurrent(root, "aws_instance.foo", `{}`, `provider["registry.terraform.io/hashicorp/aws"]`)

	ctx := testContext2(t, &ContextOpts{
		Config: m,
		Providers: map[addrs.Provider]providers.Factory{
			addrs.NewDefaultProvider("aws"): testProviderFuncFixed(p),
		},
		State: state,
	})

	s, diags := ctx.Refresh()
	if diags.HasErrors() {
		t.Fatalf("refresh errors: %s", diags.Err())
	}

	actual := strings.TrimSpace(s.String())
	expected := strings.TrimSpace(testContextRefreshOutputPartialStr)
	if actual != expected {
		t.Fatalf("wrong result\n\ngot:\n%s\n\nwant:\n%s", actual, expected)
	}
}

func TestContext2Refresh_stateBasic(t *testing.T) {
	p := testProvider("aws")
	m := testModule(t, "refresh-basic")

	state := states.NewState()
	root := state.EnsureModule(addrs.RootModuleInstance)
	testSetResourceInstanceCurrent(root, "aws_instance.web", `{"id":"bar"}`, `provider["registry.terraform.io/hashicorp/aws"]`)

	ctx := testContext2(t, &ContextOpts{
		Config: m,
		Providers: map[addrs.Provider]providers.Factory{
			addrs.NewDefaultProvider("aws"): testProviderFuncFixed(p),
		},
		State: state,
	})

	schema := p.GetSchemaReturn.ResourceTypes["aws_instance"]
	ty := schema.ImpliedType()

	readStateVal, err := schema.CoerceValue(cty.ObjectVal(map[string]cty.Value{
		"id": cty.StringVal("foo"),
	}))
	if err != nil {
		t.Fatal(err)
	}

	p.ReadResourceFn = nil
	p.PlanResourceChangeFn = testDiffFn
	p.ReadResourceResponse = &providers.ReadResourceResponse{
		NewState: readStateVal,
	}

	s, diags := ctx.Refresh()
	if diags.HasErrors() {
		t.Fatalf("refresh errors: %s", diags.Err())
	}

	if !p.ReadResourceCalled {
		t.Fatal("read resource should be called")
	}

	mod := s.RootModule()
	newState, err := mod.Resources["aws_instance.web"].Instances[addrs.NoKey].Current.Decode(ty)
	if err != nil {
		t.Fatal(err)
	}

	if !cmp.Equal(readStateVal, newState.Value, valueComparer, equateEmpty) {
		t.Fatal(cmp.Diff(readStateVal, newState.Value, valueComparer, equateEmpty))
	}
}

func TestContext2Refresh_dataCount(t *testing.T) {
	p := testProvider("test")
	m := testModule(t, "refresh-data-count")

	p.PlanResourceChangeFn = func(req providers.PlanResourceChangeRequest) (resp providers.PlanResourceChangeResponse) {
		m := req.ProposedNewState.AsValueMap()
		m["things"] = cty.ListVal([]cty.Value{cty.StringVal("foo")})
		resp.PlannedState = cty.ObjectVal(m)
		return resp
	}
	p.GetSchemaReturn = &ProviderSchema{
		ResourceTypes: map[string]*configschema.Block{
			"test": {
				Attributes: map[string]*configschema.Attribute{
					"id":     {Type: cty.String, Computed: true},
					"things": {Type: cty.List(cty.String), Computed: true},
				},
			},
		},
		DataSources: map[string]*configschema.Block{
			"test": {},
		},
	}

	p.ReadDataSourceFn = func(req providers.ReadDataSourceRequest) providers.ReadDataSourceResponse {
		return providers.ReadDataSourceResponse{
			State: req.Config,
		}
	}

	ctx := testContext2(t, &ContextOpts{
		Providers: map[addrs.Provider]providers.Factory{
			addrs.NewDefaultProvider("test"): testProviderFuncFixed(p),
		},
		Config: m,
	})

	s, diags := ctx.Refresh()

	if diags.HasErrors() {
		t.Fatalf("refresh errors: %s", diags.Err())
	}

	checkStateString(t, s, `<no state>`)
}

func TestContext2Refresh_dataState(t *testing.T) {
	m := testModule(t, "refresh-data-resource-basic")
	state := states.NewState()
	schema := &configschema.Block{
		Attributes: map[string]*configschema.Attribute{
			"inputs": {
				Type:     cty.Map(cty.String),
				Optional: true,
			},
		},
	}

	p := testProvider("null")
	p.GetSchemaReturn = &ProviderSchema{
		Provider: &configschema.Block{},
		DataSources: map[string]*configschema.Block{
			"null_data_source": schema,
		},
	}

	ctx := testContext2(t, &ContextOpts{
		Config: m,
		Providers: map[addrs.Provider]providers.Factory{
			addrs.NewDefaultProvider("null"): testProviderFuncFixed(p),
		},
		State: state,
	})

	var readStateVal cty.Value

	p.ReadDataSourceFn = func(req providers.ReadDataSourceRequest) providers.ReadDataSourceResponse {
		m := req.Config.AsValueMap()
		readStateVal = cty.ObjectVal(m)

		return providers.ReadDataSourceResponse{
			State: readStateVal,
		}
	}
	p.PlanResourceChangeFn = testDiffFn

	s, diags := ctx.Refresh()
	if diags.HasErrors() {
		t.Fatalf("refresh errors: %s", diags.Err())
	}

	if !p.ReadDataSourceCalled {
		t.Fatal("ReadDataSource should have been called")
	}

	mod := s.RootModule()

	newState, err := mod.Resources["data.null_data_source.testing"].Instances[addrs.NoKey].Current.Decode(schema.ImpliedType())
	if err != nil {
		t.Fatal(err)
	}

	if !cmp.Equal(readStateVal, newState.Value, valueComparer, equateEmpty) {
		t.Fatal(cmp.Diff(readStateVal, newState.Value, valueComparer, equateEmpty))
	}
}

func TestContext2Refresh_dataStateRefData(t *testing.T) {
	p := testProvider("null")
	p.GetSchemaReturn = &ProviderSchema{
		Provider: &configschema.Block{},
		DataSources: map[string]*configschema.Block{
			"null_data_source": {
				Attributes: map[string]*configschema.Attribute{
					"id": {
						Type:     cty.String,
						Computed: true,
					},
					"foo": {
						Type:     cty.String,
						Optional: true,
					},
					"bar": {
						Type:     cty.String,
						Optional: true,
					},
				},
			},
		},
	}

	m := testModule(t, "refresh-data-ref-data")
	state := states.NewState()
	ctx := testContext2(t, &ContextOpts{
		Config: m,
		Providers: map[addrs.Provider]providers.Factory{
			addrs.NewDefaultProvider("null"): testProviderFuncFixed(p),
		},
		State: state,
	})

	p.ReadDataSourceFn = func(req providers.ReadDataSourceRequest) providers.ReadDataSourceResponse {
		// add the required id
		m := req.Config.AsValueMap()
		m["id"] = cty.StringVal("foo")

		return providers.ReadDataSourceResponse{
			State: cty.ObjectVal(m),
		}
	}
	p.PlanResourceChangeFn = testDiffFn

	s, diags := ctx.Refresh()
	if diags.HasErrors() {
		t.Fatalf("refresh errors: %s", diags.Err())
	}

	actual := strings.TrimSpace(s.String())
	expected := strings.TrimSpace(testTerraformRefreshDataRefDataStr)
	if actual != expected {
		t.Fatalf("wrong result\n\ngot:\n%s\n\nwant:\n%s", actual, expected)
	}
}

func TestContext2Refresh_tainted(t *testing.T) {
	p := testProvider("aws")
	m := testModule(t, "refresh-basic")

	state := states.NewState()
	root := state.EnsureModule(addrs.RootModuleInstance)
	testSetResourceInstanceTainted(root, "aws_instance.web", `{"id":"bar"}`, `provider["registry.terraform.io/hashicorp/aws"]`)

	ctx := testContext2(t, &ContextOpts{
		Config: m,
		Providers: map[addrs.Provider]providers.Factory{
			addrs.NewDefaultProvider("aws"): testProviderFuncFixed(p),
		},
		State: state,
	})
	p.ReadResourceFn = func(req providers.ReadResourceRequest) providers.ReadResourceResponse {
		// add the required id
		m := req.PriorState.AsValueMap()
		m["id"] = cty.StringVal("foo")

		return providers.ReadResourceResponse{
			NewState: cty.ObjectVal(m),
		}
	}
	p.PlanResourceChangeFn = testDiffFn

	s, diags := ctx.Refresh()
	if diags.HasErrors() {
		t.Fatalf("refresh errors: %s", diags.Err())
	}
	if !p.ReadResourceCalled {
		t.Fatal("ReadResource was not called; should have been")
	}

	actual := strings.TrimSpace(s.String())
	expected := strings.TrimSpace(testContextRefreshTaintedStr)
	if actual != expected {
		t.Fatalf("wrong result\n\ngot:\n%s\n\nwant:\n%s", actual, expected)
	}
}

// Doing a Refresh (or any operation really, but Refresh usually
// happens first) with a config with an unknown provider should result in
// an error. The key bug this found was that this wasn't happening if
// Providers was _empty_.
func TestContext2Refresh_unknownProvider(t *testing.T) {
	m := testModule(t, "refresh-unknown-provider")
	p := testProvider("aws")
	p.ApplyResourceChangeFn = testApplyFn
	p.PlanResourceChangeFn = testDiffFn

	state := states.NewState()
	root := state.EnsureModule(addrs.RootModuleInstance)
	testSetResourceInstanceCurrent(root, "aws_instance.web", `{"id":"foo"}`, `provider["registry.terraform.io/hashicorp/aws"]`)

	_, diags := NewContext(&ContextOpts{
		Config:    m,
		Providers: map[addrs.Provider]providers.Factory{},
		State:     state,
	})

	if !diags.HasErrors() {
		t.Fatal("successfully created context; want error")
	}

	if !regexp.MustCompile(`failed to instantiate provider ".+"`).MatchString(diags.Err().Error()) {
		t.Fatalf("wrong error: %s", diags.Err())
	}
}

func TestContext2Refresh_vars(t *testing.T) {
	p := testProvider("aws")

	schema := &configschema.Block{
		Attributes: map[string]*configschema.Attribute{
			"ami": {
				Type:     cty.String,
				Optional: true,
			},
			"id": {
				Type:     cty.String,
				Computed: true,
			},
		},
	}

	p.GetSchemaReturn = &ProviderSchema{
		Provider:      &configschema.Block{},
		ResourceTypes: map[string]*configschema.Block{"aws_instance": schema},
	}

	m := testModule(t, "refresh-vars")
	state := states.NewState()
	root := state.EnsureModule(addrs.RootModuleInstance)
	testSetResourceInstanceCurrent(root, "aws_instance.web", `{"id":"foo"}`, `provider["registry.terraform.io/hashicorp/aws"]`)

	ctx := testContext2(t, &ContextOpts{
		Config: m,
		Providers: map[addrs.Provider]providers.Factory{
			addrs.NewDefaultProvider("aws"): testProviderFuncFixed(p),
		},
		State: state,
	})

	readStateVal, err := schema.CoerceValue(cty.ObjectVal(map[string]cty.Value{
		"id": cty.StringVal("foo"),
	}))
	if err != nil {
		t.Fatal(err)
	}

	p.ReadResourceFn = nil
	p.PlanResourceChangeFn = testDiffFn
	p.ReadResourceResponse = &providers.ReadResourceResponse{
		NewState: readStateVal,
	}

	p.PlanResourceChangeFn = func(req providers.PlanResourceChangeRequest) providers.PlanResourceChangeResponse {
		return providers.PlanResourceChangeResponse{
			PlannedState: req.ProposedNewState,
		}
	}

	s, diags := ctx.Refresh()
	if diags.HasErrors() {
		t.Fatalf("refresh errors: %s", diags.Err())
	}

	if !p.ReadResourceCalled {
		t.Fatal("read resource should be called")
	}

	mod := s.RootModule()

	newState, err := mod.Resources["aws_instance.web"].Instances[addrs.NoKey].Current.Decode(schema.ImpliedType())
	if err != nil {
		t.Fatal(err)
	}

	if !cmp.Equal(readStateVal, newState.Value, valueComparer, equateEmpty) {
		t.Fatal(cmp.Diff(readStateVal, newState.Value, valueComparer, equateEmpty))
	}

	for _, r := range mod.Resources {
		if r.Addr.Resource.Type == "" {
			t.Fatalf("no type: %#v", r)
		}
	}
}

func TestContext2Refresh_orphanModule(t *testing.T) {
	p := testProvider("aws")
	m := testModule(t, "refresh-module-orphan")

	// Create a custom refresh function to track the order they were visited
	var order []string
	var orderLock sync.Mutex
	p.ReadResourceFn = func(req providers.ReadResourceRequest) providers.ReadResourceResponse {
		orderLock.Lock()
		defer orderLock.Unlock()

		order = append(order, req.PriorState.GetAttr("id").AsString())
		return providers.ReadResourceResponse{
			NewState: req.PriorState,
		}
	}
	p.PlanResourceChangeFn = testDiffFn

	state := states.NewState()
	root := state.EnsureModule(addrs.RootModuleInstance)
	root.SetResourceInstanceCurrent(
		mustResourceInstanceAddr("aws_instance.foo").Resource,
		&states.ResourceInstanceObjectSrc{
			Status:    states.ObjectReady,
			AttrsJSON: []byte(`{"id":"i-abc123"}`),
			Dependencies: []addrs.ConfigResource{
				addrs.ConfigResource{Module: addrs.Module{"module.child"}},
				addrs.ConfigResource{Module: addrs.Module{"module.child"}},
			},
		},
		mustProviderConfig(`provider["registry.terraform.io/hashicorp/aws"]`),
	)
	child := state.EnsureModule(addrs.RootModuleInstance.Child("child", addrs.NoKey))
	child.SetResourceInstanceCurrent(
		mustResourceInstanceAddr("aws_instance.bar").Resource,
		&states.ResourceInstanceObjectSrc{
			Status:       states.ObjectReady,
			AttrsJSON:    []byte(`{"id":"i-bcd23"}`),
			Dependencies: []addrs.ConfigResource{addrs.ConfigResource{Module: addrs.Module{"module.grandchild"}}},
		},
		mustProviderConfig(`provider["registry.terraform.io/hashicorp/aws"]`),
	)
	grandchild := state.EnsureModule(addrs.RootModuleInstance.Child("child", addrs.NoKey).Child("grandchild", addrs.NoKey))
	testSetResourceInstanceCurrent(grandchild, "aws_instance.baz", `{"id":"i-cde345"}`, `provider["registry.terraform.io/hashicorp/aws"]`)

	ctx := testContext2(t, &ContextOpts{
		Config: m,
		Providers: map[addrs.Provider]providers.Factory{
			addrs.NewDefaultProvider("aws"): testProviderFuncFixed(p),
		},
		State: state,
	})

	testCheckDeadlock(t, func() {
		_, err := ctx.Refresh()
		if err != nil {
			t.Fatalf("err: %s", err.Err())
		}

		// TODO: handle order properly for orphaned modules / resources
		// expected := []string{"i-abc123", "i-bcd234", "i-cde345"}
		// if !reflect.DeepEqual(order, expected) {
		// 	t.Fatalf("expected: %#v, got: %#v", expected, order)
		// }
	})
}

func TestContext2Validate(t *testing.T) {
	p := testProvider("aws")
	p.GetSchemaReturn = &ProviderSchema{
		Provider: &configschema.Block{},
		ResourceTypes: map[string]*configschema.Block{
			"aws_instance": {
				Attributes: map[string]*configschema.Attribute{
					"foo": {
						Type:     cty.String,
						Optional: true,
					},
					"num": {
						Type:     cty.String,
						Optional: true,
					},
				},
			},
		},
	}
	p.PlanResourceChangeFn = testDiffFn

	m := testModule(t, "validate-good")
	c := testContext2(t, &ContextOpts{
		Config: m,
		Providers: map[addrs.Provider]providers.Factory{
			addrs.NewDefaultProvider("aws"): testProviderFuncFixed(p),
		},
	})

	diags := c.Validate()
	if len(diags) != 0 {
		t.Fatalf("unexpected error: %#v", diags.ErrWithWarnings())
	}
}

func TestContext2Refresh_updateProviderInState(t *testing.T) {
	m := testModule(t, "update-resource-provider")
	p := testProvider("aws")
	p.PlanResourceChangeFn = testDiffFn
	p.ApplyResourceChangeFn = testApplyFn

	state := states.NewState()
	root := state.EnsureModule(addrs.RootModuleInstance)
	testSetResourceInstanceCurrent(root, "aws_instance.bar", `{"id":"foo"}`, `provider["registry.terraform.io/hashicorp/aws"].baz`)

	ctx := testContext2(t, &ContextOpts{
		Config: m,
		Providers: map[addrs.Provider]providers.Factory{
			addrs.NewDefaultProvider("aws"): testProviderFuncFixed(p),
		},
		State: state,
	})

	expected := strings.TrimSpace(`
aws_instance.bar:
  ID = foo
  provider = provider["registry.terraform.io/hashicorp/aws"].foo`)

	s, diags := ctx.Refresh()
	if diags.HasErrors() {
		t.Fatal(diags.Err())
	}

	actual := s.String()
	if actual != expected {
		t.Fatalf("expected:\n%s\n\ngot:\n%s", expected, actual)
	}
}

func TestContext2Refresh_schemaUpgradeFlatmap(t *testing.T) {
	m := testModule(t, "refresh-schema-upgrade")
	p := testProvider("test")
	p.GetSchemaReturn = &ProviderSchema{
		ResourceTypes: map[string]*configschema.Block{
			"test_thing": {
				Attributes: map[string]*configschema.Attribute{
					"name": { // imagining we renamed this from "id"
						Type:     cty.String,
						Optional: true,
					},
				},
			},
		},
		ResourceTypeSchemaVersions: map[string]uint64{
			"test_thing": 5,
		},
	}
	p.UpgradeResourceStateResponse = &providers.UpgradeResourceStateResponse{
		UpgradedState: cty.ObjectVal(map[string]cty.Value{
			"name": cty.StringVal("foo"),
		}),
	}
	p.PlanResourceChangeFn = testDiffFn

	s := states.BuildState(func(s *states.SyncState) {
		s.SetResourceInstanceCurrent(
			addrs.Resource{
				Mode: addrs.ManagedResourceMode,
				Type: "test_thing",
				Name: "bar",
			}.Instance(addrs.NoKey).Absolute(addrs.RootModuleInstance),
			&states.ResourceInstanceObjectSrc{
				Status:        states.ObjectReady,
				SchemaVersion: 3,
				AttrsFlat: map[string]string{
					"id": "foo",
				},
			},
			addrs.AbsProviderConfig{
				Provider: addrs.NewDefaultProvider("test"),
				Module:   addrs.RootModule,
			},
		)
	})

	ctx := testContext2(t, &ContextOpts{
		Config: m,
		Providers: map[addrs.Provider]providers.Factory{
			addrs.NewDefaultProvider("test"): testProviderFuncFixed(p),
		},
		State: s,
	})

	state, diags := ctx.Refresh()
	if diags.HasErrors() {
		t.Fatal(diags.Err())
	}

	{
		got := p.UpgradeResourceStateRequest
		want := providers.UpgradeResourceStateRequest{
			TypeName: "test_thing",
			Version:  3,
			RawStateFlatmap: map[string]string{
				"id": "foo",
			},
		}
		if !cmp.Equal(got, want) {
			t.Errorf("wrong upgrade request\n%s", cmp.Diff(want, got))
		}
	}

	{
		got := state.String()
		want := strings.TrimSpace(`
test_thing.bar:
  ID = 
  provider = provider["registry.terraform.io/hashicorp/test"]
  name = foo
`)
		if got != want {
			t.Fatalf("wrong result state\ngot:\n%s\n\nwant:\n%s", got, want)
		}
	}
}

func TestContext2Refresh_schemaUpgradeJSON(t *testing.T) {
	m := testModule(t, "refresh-schema-upgrade")
	p := testProvider("test")
	p.GetSchemaReturn = &ProviderSchema{
		ResourceTypes: map[string]*configschema.Block{
			"test_thing": {
				Attributes: map[string]*configschema.Attribute{
					"name": { // imagining we renamed this from "id"
						Type:     cty.String,
						Optional: true,
					},
				},
			},
		},
		ResourceTypeSchemaVersions: map[string]uint64{
			"test_thing": 5,
		},
	}
	p.UpgradeResourceStateResponse = &providers.UpgradeResourceStateResponse{
		UpgradedState: cty.ObjectVal(map[string]cty.Value{
			"name": cty.StringVal("foo"),
		}),
	}
	p.PlanResourceChangeFn = testDiffFn

	s := states.BuildState(func(s *states.SyncState) {
		s.SetResourceInstanceCurrent(
			addrs.Resource{
				Mode: addrs.ManagedResourceMode,
				Type: "test_thing",
				Name: "bar",
			}.Instance(addrs.NoKey).Absolute(addrs.RootModuleInstance),
			&states.ResourceInstanceObjectSrc{
				Status:        states.ObjectReady,
				SchemaVersion: 3,
				AttrsJSON:     []byte(`{"id":"foo"}`),
			},
			addrs.AbsProviderConfig{
				Provider: addrs.NewDefaultProvider("test"),
				Module:   addrs.RootModule,
			},
		)
	})

	ctx := testContext2(t, &ContextOpts{
		Config: m,
		Providers: map[addrs.Provider]providers.Factory{
			addrs.NewDefaultProvider("test"): testProviderFuncFixed(p),
		},
		State: s,
	})

	state, diags := ctx.Refresh()
	if diags.HasErrors() {
		t.Fatal(diags.Err())
	}

	{
		got := p.UpgradeResourceStateRequest
		want := providers.UpgradeResourceStateRequest{
			TypeName:     "test_thing",
			Version:      3,
			RawStateJSON: []byte(`{"id":"foo"}`),
		}
		if !cmp.Equal(got, want) {
			t.Errorf("wrong upgrade request\n%s", cmp.Diff(want, got))
		}
	}

	{
		got := state.String()
		want := strings.TrimSpace(`
test_thing.bar:
  ID = 
  provider = provider["registry.terraform.io/hashicorp/test"]
  name = foo
`)
		if got != want {
			t.Fatalf("wrong result state\ngot:\n%s\n\nwant:\n%s", got, want)
		}
	}
}

func TestContext2Refresh_dataValidation(t *testing.T) {
	m := testModuleInline(t, map[string]string{
		"main.tf": `
data "aws_data_source" "foo" {
  foo = "bar"
}
`,
	})

	p := testProvider("aws")
	p.PlanResourceChangeFn = func(req providers.PlanResourceChangeRequest) (resp providers.PlanResourceChangeResponse) {
		resp.PlannedState = req.ProposedNewState
		return
	}
	p.ReadDataSourceFn = func(req providers.ReadDataSourceRequest) (resp providers.ReadDataSourceResponse) {
		resp.State = req.Config
		return
	}
	p.PlanResourceChangeFn = testDiffFn

	ctx := testContext2(t, &ContextOpts{
		Config: m,
		Providers: map[addrs.Provider]providers.Factory{
			addrs.NewDefaultProvider("aws"): testProviderFuncFixed(p),
		},
	})

	_, diags := ctx.Refresh()
	if diags.HasErrors() {
		// Should get this error:
		// Unsupported attribute: This object does not have an attribute named "missing"
		t.Fatal(diags.Err())
	}

	if !p.ValidateDataSourceConfigCalled {
		t.Fatal("ValidateDataSourceConfig not called during plan")
	}
}

func TestContext2Refresh_dataResourceDependsOn(t *testing.T) {
	m := testModule(t, "plan-data-depends-on")
	p := testProvider("test")
	p.GetSchemaReturn = &ProviderSchema{
		ResourceTypes: map[string]*configschema.Block{
			"test_resource": {
				Attributes: map[string]*configschema.Attribute{
					"id":  {Type: cty.String, Computed: true},
					"foo": {Type: cty.String, Optional: true},
				},
			},
		},
		DataSources: map[string]*configschema.Block{
			"test_data": {
				Attributes: map[string]*configschema.Attribute{
					"compute": {Type: cty.String, Computed: true},
				},
			},
		},
	}
	p.PlanResourceChangeFn = testDiffFn
	p.ReadDataSourceResponse = &providers.ReadDataSourceResponse{
		State: cty.ObjectVal(map[string]cty.Value{
			"compute": cty.StringVal("value"),
		}),
	}

	state := states.NewState()
	root := state.EnsureModule(addrs.RootModuleInstance)
	testSetResourceInstanceCurrent(root, "test_resource.a", `{"id":"a"}`, `provider["registry.terraform.io/hashicorp/test"]`)

	ctx := testContext2(t, &ContextOpts{
		Config: m,
		Providers: map[addrs.Provider]providers.Factory{
			addrs.NewDefaultProvider("test"): testProviderFuncFixed(p),
		},
		State: state,
	})

	_, diags := ctx.Refresh()
	if diags.HasErrors() {
		t.Fatalf("unexpected errors: %s", diags.Err())
	}
}

// verify that dependencies are updated in the state during refresh
func TestRefresh_updateDependencies(t *testing.T) {
	state := states.NewState()
	root := state.EnsureModule(addrs.RootModuleInstance)
	root.SetResourceInstanceCurrent(
		addrs.Resource{
			Mode: addrs.ManagedResourceMode,
			Type: "aws_instance",
			Name: "foo",
		}.Instance(addrs.NoKey),
		&states.ResourceInstanceObjectSrc{
			Status:    states.ObjectReady,
			AttrsJSON: []byte(`{"id":"foo"}`),
			Dependencies: []addrs.ConfigResource{
				// Existing dependencies should be removed when overridden by the config
				{
					Module: addrs.RootModule,
					Resource: addrs.Resource{
						Mode: addrs.ManagedResourceMode,
						Type: "aws_instance",
						Name: "baz",
					},
				},
			},
		},
		addrs.AbsProviderConfig{
			Provider: addrs.NewDefaultProvider("aws"),
			Module:   addrs.RootModule,
		},
	)
	root.SetResourceInstanceCurrent(
		addrs.Resource{
			Mode: addrs.ManagedResourceMode,
			Type: "aws_instance",
			Name: "baz",
		}.Instance(addrs.NoKey),
		&states.ResourceInstanceObjectSrc{
			Status:    states.ObjectReady,
			AttrsJSON: []byte(`{"id":"baz"}`),
			Dependencies: []addrs.ConfigResource{
				// Existing dependencies should not be removed from orphaned instances
				{
					Module: addrs.RootModule,
					Resource: addrs.Resource{
						Mode: addrs.ManagedResourceMode,
						Type: "aws_instance",
						Name: "bam",
					},
				},
			},
		},
		addrs.AbsProviderConfig{
			Provider: addrs.NewDefaultProvider("aws"),
			Module:   addrs.RootModule,
		},
	)
	root.SetResourceInstanceCurrent(
		addrs.Resource{
			Mode: addrs.ManagedResourceMode,
			Type: "aws_instance",
			Name: "bar",
		}.Instance(addrs.NoKey),
		&states.ResourceInstanceObjectSrc{
			Status:    states.ObjectReady,
			AttrsJSON: []byte(`{"id":"bar","foo":"foo"}`),
		},
		addrs.AbsProviderConfig{
			Provider: addrs.NewDefaultProvider("aws"),
			Module:   addrs.RootModule,
		},
	)

	m := testModuleInline(t, map[string]string{
		"main.tf": `
resource "aws_instance" "bar" {
  foo = aws_instance.foo.id
}

resource "aws_instance" "foo" {
}
`,
	})

	p := testProvider("aws")
	p.ApplyResourceChangeFn = testApplyFn
	p.PlanResourceChangeFn = testDiffFn

	ctx := testContext2(t, &ContextOpts{
		Config: m,
		Providers: map[addrs.Provider]providers.Factory{
			addrs.NewDefaultProvider("aws"): testProviderFuncFixed(p),
		},
		State: state,
	})

	result, diags := ctx.Refresh()
	if diags.HasErrors() {
		t.Fatalf("plan errors: %s", diags.Err())
	}

	expect := strings.TrimSpace(`
aws_instance.bar:
  ID = bar
  provider = provider["registry.terraform.io/hashicorp/aws"]
  foo = foo

  Dependencies:
    aws_instance.foo
aws_instance.baz:
  ID = baz
  provider = provider["registry.terraform.io/hashicorp/aws"]

  Dependencies:
    aws_instance.bam
aws_instance.foo:
  ID = foo
  provider = provider["registry.terraform.io/hashicorp/aws"]
`)

	checkStateString(t, result, expect)
}

// verify that create_before_destroy is updated in the state during refresh
func TestRefresh_updateLifecycle(t *testing.T) {
	state := states.NewState()
	root := state.EnsureModule(addrs.RootModuleInstance)
	root.SetResourceInstanceCurrent(
		addrs.Resource{
			Mode: addrs.ManagedResourceMode,
			Type: "aws_instance",
			Name: "bar",
		}.Instance(addrs.NoKey),
		&states.ResourceInstanceObjectSrc{
			Status:    states.ObjectReady,
			AttrsJSON: []byte(`{"id":"bar"}`),
		},
		addrs.AbsProviderConfig{
			Provider: addrs.NewDefaultProvider("aws"),
			Module:   addrs.RootModule,
		},
	)

	m := testModuleInline(t, map[string]string{
		"main.tf": `
resource "aws_instance" "bar" {
  lifecycle {
    create_before_destroy = true
  }
}
`,
	})

	p := testProvider("aws")
	p.ApplyResourceChangeFn = testApplyFn
	p.PlanResourceChangeFn = testDiffFn

	ctx := testContext2(t, &ContextOpts{
		Config: m,
		Providers: map[addrs.Provider]providers.Factory{
			addrs.NewDefaultProvider("aws"): testProviderFuncFixed(p),
		},
		State: state,
	})

	state, diags := ctx.Refresh()
	if diags.HasErrors() {
		t.Fatalf("plan errors: %s", diags.Err())
	}

	r := state.ResourceInstance(mustResourceInstanceAddr("aws_instance.bar"))
	if !r.Current.CreateBeforeDestroy {
		t.Fatal("create_before_destroy not updated in instance state")
	}
}

func TestContext2Refresh_dataSourceOrphan(t *testing.T) {
	m := testModuleInline(t, map[string]string{
		"main.tf": ``,
	})

	state := states.NewState()
	root := state.EnsureModule(addrs.RootModuleInstance)
	root.SetResourceInstanceCurrent(
		addrs.Resource{
			Mode: addrs.DataResourceMode,
			Type: "test_data_source",
			Name: "foo",
		}.Instance(addrs.NoKey),
		&states.ResourceInstanceObjectSrc{
			Status:       states.ObjectReady,
			AttrsJSON:    []byte(`{"id":"foo"}`),
			Dependencies: []addrs.ConfigResource{},
		},
		addrs.AbsProviderConfig{
			Provider: addrs.NewDefaultProvider("test"),
			Module:   addrs.RootModule,
		},
	)
	p := testProvider("test")
	p.ReadDataSourceFn = func(req providers.ReadDataSourceRequest) (resp providers.ReadDataSourceResponse) {
		resp.State = cty.NullVal(req.Config.Type())
		return
	}

	ctx := testContext2(t, &ContextOpts{
		Config: m,
		Providers: map[addrs.Provider]providers.Factory{
			addrs.NewDefaultProvider("test"): testProviderFuncFixed(p),
		},
		State: state,
	})

	_, diags := ctx.Refresh()
	if diags.HasErrors() {
		t.Fatal(diags.Err())
	}

	if p.ReadResourceCalled {
		t.Fatal("there are no managed resources to read")
	}

	if p.ReadDataSourceCalled {
		t.Fatal("orphaned data source instance should not be read")
	}
}
