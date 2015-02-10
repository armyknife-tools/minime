package terraform

import (
	"fmt"

	"github.com/hashicorp/terraform/config"
	"github.com/hashicorp/terraform/dag"
)

// ResourceCountTransformer is a GraphTransformer that expands the count
// out for a specific resource.
type ResourceCountTransformer struct {
	Resource *config.Resource
}

func (t *ResourceCountTransformer) Transform(g *Graph) error {
	// Expand the resource count
	count, err := t.Resource.Count()
	if err != nil {
		return err
	}

	// Don't allow the count to be negative
	if count < 0 {
		return fmt.Errorf("negative count: %d", count)
	}

	// For each count, build and add the node
	nodes := make([]dag.Vertex, count)
	for i := 0; i < count; i++ {
		// Save the node for later so we can do connections
		nodes[i] = &graphNodeExpandedResource{
			Index:    i,
			Resource: t.Resource,
		}

		// Add the node now
		g.Add(nodes[i])
	}

	// Make the dependency connections
	for _, n := range nodes {
		// Connect the dependents. We ignore the return value for missing
		// dependents since that should've been caught at a higher level.
		g.ConnectDependent(n)
	}

	return nil
}

type graphNodeExpandedResource struct {
	Index    int
	Resource *config.Resource
}

func (n *graphNodeExpandedResource) Name() string {
	return fmt.Sprintf("%s #%d", n.Resource.Id(), n.Index)
}

// GraphNodeDependable impl.
func (n *graphNodeExpandedResource) DependableName() []string {
	return []string{
		n.Resource.Id(),
		n.stateId(),
	}
}

// GraphNodeDependent impl.
func (n *graphNodeExpandedResource) DependentOn() []string {
	config := &GraphNodeConfigResource{Resource: n.Resource}
	return config.DependentOn()
}

// GraphNodeProviderConsumer
func (n *graphNodeExpandedResource) ProvidedBy() string {
	return resourceProvider(n.Resource.Type)
}

// GraphNodeEvalable impl.
func (n *graphNodeExpandedResource) EvalTree() EvalNode {
	seq := &EvalSequence{Nodes: make([]EvalNode, 0, 5)}

	// Validate the resource
	seq.Nodes = append(seq.Nodes, &EvalValidateResource{
		Provider:     &EvalGetProvider{Name: n.ProvidedBy()},
		Config:       &EvalInterpolate{Config: n.Resource.RawConfig},
		ResourceName: n.Resource.Name,
		ResourceType: n.Resource.Type,
	})

	// Validate all the provisioners
	for _, p := range n.Resource.Provisioners {
		seq.Nodes = append(seq.Nodes, &EvalValidateProvisioner{
			Provisioner: &EvalGetProvisioner{Name: p.Type},
			Config:      &EvalInterpolate{Config: p.RawConfig},
		})
	}

	return seq
}

// stateId is the name used for the state key
func (n *graphNodeExpandedResource) stateId() string {
	return fmt.Sprintf("%s.%d", n.Resource.Id(), n.Index)
}
