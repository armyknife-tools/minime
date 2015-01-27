package terraform

import (
	"errors"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform/config"
	"github.com/hashicorp/terraform/config/module"
)

// ConfigTransformer is a GraphTransformer that adds the configuration
// to the graph. It is assumed that the module tree given in Module matches
// the Path attribute of the Graph being transformed. If this is not the case,
// the behavior is unspecified, but unlikely to be what you want.
type ConfigTransformer struct {
	Module *module.Tree
}

func (t *ConfigTransformer) Transform(g *Graph) error {
	// A module is required and also must be completely loaded.
	if t.Module == nil {
		return errors.New("module must not be nil")
	}
	if !t.Module.Loaded() {
		return errors.New("module must be loaded")
	}

	// Get the configuration for this module
	config := t.Module.Config()

	// Create the node list we'll use for the graph
	nodes := make([]graphNodeConfig, 0,
		(len(config.ProviderConfigs)+len(config.Modules)+len(config.Resources))*2)

	// Write all the provider configs out
	for _, pc := range config.ProviderConfigs {
		nodes = append(nodes, &GraphNodeConfigProvider{Provider: pc})
	}

	// Write all the resources out
	for _, r := range config.Resources {
		nodes = append(nodes, &GraphNodeConfigResource{Resource: r})
	}

	// Write all the modules out
	children := t.Module.Children()
	for _, m := range config.Modules {
		nodes = append(nodes, &GraphNodeConfigModule{
			Module: m,
			Tree:   children[m.Name],
		})
	}

	// Err is where the final error value will go if there is one
	var err error

	// Build the graph vertices
	for _, n := range nodes {
		g.Add(n)
	}

	// Build up the dependencies. We have to do this outside of the above
	// loop since the nodes need to be in place for us to build the deps.
	for _, n := range nodes {
		if missing := g.ConnectDependent(n); len(missing) > 0 {
			for _, m := range missing {
				err = multierror.Append(err, fmt.Errorf(
					"%s: missing dependency: %s", n.Name(), m))
			}
		}
	}

	return err
}

// varNameForVar returns the VarName value for an interpolated variable.
// This value is compared to the VarName() value for the nodes within the
// graph to build the graph edges.
func varNameForVar(raw config.InterpolatedVariable) string {
	switch v := raw.(type) {
	case *config.ModuleVariable:
		return fmt.Sprintf("module.%s", v.Name)
	case *config.ResourceVariable:
		return v.ResourceId()
	default:
		return ""
	}
}
