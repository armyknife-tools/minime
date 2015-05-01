package terraform

import (
	"fmt"

	"github.com/hashicorp/terraform/config"
	"github.com/hashicorp/terraform/dag"
)

// GraphNodeVariable is the interface that must be implemented by anything
// that is a variable.
type GraphNodeVariable interface {
	VariableName() string
	SetVariableValue(*config.RawConfig)
}

// GraphNodeConfigVariable represents a Variable in the config.
type GraphNodeConfigVariable struct {
	Variable *config.Variable

	// Value, if non-nil, will be used to set the value of the variable
	// during evaluation. If this is nil, evaluation will do nothing.
	Value *config.RawConfig

	depPrefix string
}

func (n *GraphNodeConfigVariable) Name() string {
	return fmt.Sprintf("var.%s", n.Variable.Name)
}

func (n *GraphNodeConfigVariable) ConfigType() GraphNodeConfigType {
	return GraphNodeConfigTypeVariable
}

func (n *GraphNodeConfigVariable) DependableName() []string {
	return []string{n.Name()}
}

func (n *GraphNodeConfigVariable) DependentOn() []string {
	// If we don't have any value set, we don't depend on anything
	if n.Value == nil {
		return nil
	}

	// Get what we depend on based on our value
	vars := n.Value.Variables
	result := make([]string, 0, len(vars))
	for _, v := range vars {
		if vn := varNameForVar(v); vn != "" {
			result = append(result, vn)
		}
	}

	return result
}

// GraphNodeVariable impl.
func (n *GraphNodeConfigVariable) VariableName() string {
	return n.Variable.Name
}

// GraphNodeVariable impl.
func (n *GraphNodeConfigVariable) SetVariableValue(v *config.RawConfig) {
	n.Value = v
}

// GraphNodeProxy impl.
func (n *GraphNodeConfigVariable) Proxy() bool {
	return true
}

// GraphNodeEvalable impl.
func (n *GraphNodeConfigVariable) EvalTree() EvalNode {
	// If we have no value, do nothing
	if n.Value == nil {
		return &EvalNoop{}
	}

	// Otherwise, interpolate the value of this variable and set it
	// within the variables mapping.
	var config *ResourceConfig
	variables := make(map[string]string)
	return &EvalSequence{
		Nodes: []EvalNode{
			&EvalInterpolate{
				Config: n.Value,
				Output: &config,
			},

			&EvalVariableBlock{
				Config:    &config,
				Variables: variables,
			},

			&EvalSetVariables{
				Variables: variables,
			},
		},
	}
}

// GraphNodeFlattenable impl.
func (n *GraphNodeConfigVariable) Flatten(p []string) (dag.Vertex, error) {
	return &GraphNodeConfigVariableFlat{
		GraphNodeConfigVariable: n,
		PathValue:               p,
	}, nil
}

type GraphNodeConfigVariableFlat struct {
	*GraphNodeConfigVariable

	PathValue []string
}

func (n *GraphNodeConfigVariableFlat) Name() string {
	return fmt.Sprintf(
		"%s.%s", modulePrefixStr(n.PathValue), n.GraphNodeConfigVariable.Name())
}

func (n *GraphNodeConfigVariableFlat) DependableName() []string {
	return []string{n.Name()}
}

func (n *GraphNodeConfigVariableFlat) DependentOn() []string {
	// We only wrap the dependencies and such if we have a path that is
	// longer than 2 elements (root, child, more). This is because when
	// flattened, variables can point outside the graph.
	prefix := ""
	if len(n.PathValue) <= 2 {
		prefix = modulePrefixStr(n.PathValue[:len(n.PathValue)-1])
	}

	return modulePrefixList(
		n.GraphNodeConfigVariable.DependentOn(),
		prefix)
}

func (n *GraphNodeConfigVariableFlat) Path() []string {
	return n.PathValue
}
