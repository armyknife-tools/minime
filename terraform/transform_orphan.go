package terraform

import (
	"fmt"

	"github.com/hashicorp/terraform/config"
	"github.com/hashicorp/terraform/dag"
)

// OrphanTransformer is a GraphTransformer that adds orphans to the
// graph. This transformer adds both resource and module orphans.
type OrphanTransformer struct {
	// State is the global state. We require the global state to
	// properly find module orphans at our path.
	State *State

	// Config is just the configuration of our current module.
	Config *config.Config
}

func (t *OrphanTransformer) Transform(g *Graph) error {
	state := t.State.ModuleByPath(g.Path)
	if state == nil {
		// If there is no state for our module, there can't be any orphans
		return nil
	}

	// Go over each resource orphan and add it to the graph.
	resourceOrphans := state.Orphans(t.Config)
	resourceVertexes := make([]dag.Vertex, len(resourceOrphans))
	for i, k := range resourceOrphans {
		resourceVertexes[i] = g.Add(&graphNodeOrphanResource{ResourceName: k})
	}

	// Go over each module orphan and add it to the graph. We store the
	// vertexes and states outside so that we can connect dependencies later.
	moduleOrphans := t.State.ModuleOrphans(g.Path, t.Config)
	moduleVertexes := make([]dag.Vertex, len(moduleOrphans))
	moduleStates := make([]*ModuleState, len(moduleVertexes))
	for i, path := range moduleOrphans {
		moduleVertexes[i] = g.Add(&graphNodeOrphanModule{Path: path})
		moduleStates[i] = t.State.ModuleByPath(path)
	}

	// Now do the dependencies. We do this _after_ adding all the orphan
	// nodes above because there are cases in which the orphans themselves
	// depend on other orphans.

	// Resource dependencies
	for i, v := range resourceVertexes {
		g.ConnectTo(v, state.Resources[resourceOrphans[i]].Dependencies)
	}

	// Module dependencies
	for i, v := range moduleVertexes {
		g.ConnectTo(v, moduleStates[i].Dependencies)
	}

	return nil
}

// graphNodeOrphanModule is the graph vertex representing an orphan resource..
type graphNodeOrphanModule struct {
	Path []string
}

func (n *graphNodeOrphanModule) Name() string {
	return fmt.Sprintf("module.%s (orphan)", n.Path[len(n.Path)-1])
}

// graphNodeOrphanResource is the graph vertex representing an orphan resource..
type graphNodeOrphanResource struct {
	ResourceName string
}

func (n *graphNodeOrphanResource) Name() string {
	return fmt.Sprintf("%s (orphan)", n.ResourceName)
}
