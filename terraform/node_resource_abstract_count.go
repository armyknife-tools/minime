package terraform

// NodeAbstractCountResource should be embedded instead of NodeAbstractResource
// if the resource has a `count` value that needs to be expanded.
//
// The embedder should implement `DynamicExpand` to process the count.
type NodeAbstractCountResource struct {
	*NodeAbstractResource

	// Validate, if true, will perform the validation for the count.
	// This should only be turned on for the "validate" operation.
	Validate bool
}

// GraphNodeEvalable
func (n *NodeAbstractCountResource) EvalTree() EvalNode {
	return &EvalSequence{
		Nodes: []EvalNode{
			// The EvalTree for a plannable resource primarily involves
			// interpolating the count since it can contain variables
			// we only just received access to.
			//
			// With the interpolated count, we can then DynamicExpand
			// into the proper number of instances.
			&EvalInterpolate{Config: n.Config.RawCount},

			&EvalCountCheckComputed{Resource: n.Config},

			// If validation is enabled, perform the validation
			&EvalIf{
				If: func(ctx EvalContext) (bool, error) {
					return n.Validate, nil
				},

				Then: &EvalValidateCount{Resource: n.Config},
			},

			&EvalCountFixZeroOneBoundary{Resource: n.Config},
		},
	}
}
