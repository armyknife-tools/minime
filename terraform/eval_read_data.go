package terraform

import (
	"fmt"
)

// EvalReadDataDiff is an EvalNode implementation that executes a data
// resource's ReadDataDiff method to discover what attributes it exports.
type EvalReadDataDiff struct {
	Provider    *ResourceProvider
	Output      **InstanceDiff
	OutputState **InstanceState
	Config      **ResourceConfig
	Info        *InstanceInfo
}

func (n *EvalReadDataDiff) Eval(ctx EvalContext) (interface{}, error) {
	// TODO: test
	provider := *n.Provider
	config := *n.Config

	err := ctx.Hook(func(h Hook) (HookAction, error) {
		return h.PreDiff(n.Info, nil)
	})
	if err != nil {
		return nil, err
	}

	diff, err := provider.ReadDataDiff(n.Info, config)
	if err != nil {
		return nil, err
	}
	if diff == nil {
		diff = new(InstanceDiff)
	}

	// id is always computed, because we're always "creating a new resource"
	diff.init()
	diff.Attributes["id"] = &ResourceAttrDiff{
		Old:         "",
		NewComputed: true,
		RequiresNew: true,
		Type:        DiffAttrOutput,
	}

	err = ctx.Hook(func(h Hook) (HookAction, error) {
		return h.PostDiff(n.Info, diff)
	})
	if err != nil {
		return nil, err
	}

	*n.Output = diff

	if n.OutputState != nil {
		state := &InstanceState{}
		*n.OutputState = state

		// Apply the diff to the returned state, so the state includes
		// any attribute values that are not computed.
		if !diff.Empty() && n.OutputState != nil {
			*n.OutputState = state.MergeDiff(diff)
		}
	}

	return nil, nil
}

// EvalReadDataApply is an EvalNode implementation that executes a data
// resource's ReadDataApply method to read data from the data source.
type EvalReadDataApply struct {
	Provider *ResourceProvider
	Output   **InstanceState
	Diff     **InstanceDiff
	Info     *InstanceInfo
}

func (n *EvalReadDataApply) Eval(ctx EvalContext) (interface{}, error) {
	// TODO: test
	provider := *n.Provider
	diff := *n.Diff

	// For the purpose of external hooks we present a data apply as a
	// "Refresh" rather than an "Apply" because creating a data source
	// is presented to users/callers as a "read" operation.
	err := ctx.Hook(func(h Hook) (HookAction, error) {
		// We don't have a state yet, so we'll just give the hook an
		// empty one to work with.
		return h.PreRefresh(n.Info, &InstanceState{})
	})
	if err != nil {
		return nil, err
	}

	state, err := provider.ReadDataApply(n.Info, diff)
	if err != nil {
		return nil, fmt.Errorf("%s: %s", n.Info.Id, err)
	}

	err = ctx.Hook(func(h Hook) (HookAction, error) {
		return h.PostRefresh(n.Info, state)
	})
	if err != nil {
		return nil, err
	}

	if n.Output != nil {
		*n.Output = state
	}

	return nil, nil
}
