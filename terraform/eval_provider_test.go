package terraform

import (
	"testing"
)

func TestEvalInitProvider_impl(t *testing.T) {
	var _ EvalNode = new(EvalInitProvider)
}

func TestEvalInitProvider(t *testing.T) {
	n := &EvalInitProvider{Name: "foo"}
	provider := &MockResourceProvider{}
	ctx := &MockEvalContext{InitProviderProvider: provider}
	if actual, err := n.Eval(ctx, nil); err != nil {
		t.Fatalf("err: %s", err)
	} else if actual != provider {
		t.Fatalf("bad: %#v", actual)
	}

	if !ctx.InitProviderCalled {
		t.Fatal("should be called")
	}
	if ctx.InitProviderName != "foo" {
		t.Fatalf("bad: %#v", ctx.InitProviderName)
	}
}

func TestEvalGetProvider_impl(t *testing.T) {
	var _ EvalNode = new(EvalGetProvider)
}

func TestEvalGetProvider(t *testing.T) {
	n := &EvalGetProvider{Name: "foo"}
	provider := &MockResourceProvider{}
	ctx := &MockEvalContext{ProviderProvider: provider}
	if actual, err := n.Eval(ctx, nil); err != nil {
		t.Fatalf("err: %s", err)
	} else if actual != provider {
		t.Fatalf("bad: %#v", actual)
	}

	if !ctx.ProviderCalled {
		t.Fatal("should be called")
	}
	if ctx.ProviderName != "foo" {
		t.Fatalf("bad: %#v", ctx.ProviderName)
	}
}
