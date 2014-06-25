package terraform

import (
	"fmt"
	"log"
	"sync"

	"github.com/hashicorp/terraform/config"
	"github.com/hashicorp/terraform/depgraph"
)

// Terraform is the primary structure that is used to interact with
// Terraform from code, and can perform operations such as returning
// all resources, a resource tree, a specific resource, etc.
type Terraform struct {
	providers map[string]ResourceProviderFactory
}

// terraformProvider contains internal state information about a resource
// provider for Terraform.
type terraformProvider struct {
	Provider ResourceProvider
	Config   *config.ProviderConfig

	sync.Once
}

// This is a function type used to implement a walker for the resource
// tree internally on the Terraform structure.
type genericWalkFunc func(*Resource) (map[string]string, error)

// Config is the configuration that must be given to instantiate
// a Terraform structure.
type Config struct {
	Config    *config.Config
	Providers map[string]ResourceProviderFactory
	Variables map[string]string
}

// New creates a new Terraform structure, initializes resource providers
// for the given configuration, etc.
//
// Semantic checks of the entire configuration structure are done at this
// time, as well as richer checks such as verifying that the resource providers
// can be properly initialized, can be configured, etc.
func New(c *Config) (*Terraform, error) {
	var errs []error

	if c.Config != nil {
		// Validate that all required variables have values
		if err := smcVariables(c); err != nil {
			errs = append(errs, err...)
		}

		// Match all the resources with a provider and initialize the providers
		mapping, err := smcProviders(c)
		if err != nil {
			errs = append(errs, err...)
		}

		// Validate all the configurations, once.
		tps := make(map[*terraformProvider]struct{})
		for _, tp := range mapping {
			if _, ok := tps[tp]; !ok {
				tps[tp] = struct{}{}
			}
		}
		for tp, _ := range tps {
			var rc *ResourceConfig
			if tp.Config != nil {
				rc = NewResourceConfig(tp.Config.RawConfig)
			}

			_, tpErrs := tp.Provider.Validate(rc)
			if len(tpErrs) > 0 {
				errs = append(errs, tpErrs...)
			}
		}
	}

	// If we accumulated any errors, then return them all
	if len(errs) > 0 {
		return nil, &MultiError{Errors: errs}
	}

	return &Terraform{
		providers: c.Providers,
	}, nil
}

func (t *Terraform) Apply(p *Plan) (*State, error) {
	graph, err := t.Graph(p.Config, p.State)
	if err != nil {
		return nil, err
	}

	result := new(State)
	err = graph.Walk(t.applyWalkFn(p, result))
	return result, err
}

// Graph returns the dependency graph for the given configuration and
// state file.
//
// The resulting graph may have more resources than the configuration, because
// it can contain resources in the state file that need to be modified.
func (t *Terraform) Graph(c *config.Config, s *State) (*depgraph.Graph, error) {
	// Get the basic graph with the raw metadata
	g := Graph(c, s)
	if err := g.Validate(); err != nil {
		return nil, err
	}

	// Fill the graph with the providers
	if err := GraphFull(g, t.providers); err != nil {
		return nil, err
	}

	// Validate the graph so that it can setup a root and such
	if err := g.Validate(); err != nil {
		return nil, err
	}

	return g, nil
}

func (t *Terraform) Plan(s *State) (*Plan, error) {
	// TODO: add config param
	graph, err := t.Graph(nil, s)
	if err != nil {
		return nil, err
	}

	result := new(Plan)
	err = graph.Walk(t.planWalkFn(s, result))
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Refresh goes through all the resources in the state and refreshes them
// to their latest status.
func (t *Terraform) Refresh(c *config.Config, s *State) (*State, error) {
	g, err := t.Graph(c, s)
	if err != nil {
		return s, err
	}

	return t.refresh(g)
}

func (t *Terraform) refresh(g *depgraph.Graph) (*State, error) {
	s := new(State)
	err := g.Walk(t.refreshWalkFn(s))
	return s, err
}

func (t *Terraform) refreshWalkFn(result *State) depgraph.WalkFunc {
	var l sync.Mutex

	// Initialize the result so we don't have to nil check everywhere
	result.init()

	cb := func(r *Resource) (map[string]string, error) {
		rs, err := r.Provider.Refresh(r.State)
		if err != nil {
			return nil, err
		}

		l.Lock()
		result.Resources[r.Id] = rs
		l.Unlock()

		return nil, nil
	}

	return t.genericWalkFn(nil, cb)
}

func (t *Terraform) applyWalkFn(
	p *Plan,
	result *State) depgraph.WalkFunc {
	var l sync.Mutex

	// Initialize the result
	result.init()

	cb := func(r *Resource) (map[string]string, error) {
		// Get the latest diff since there are no computed values anymore
		diff, err := r.Provider.Diff(r.State, r.Config)
		if err != nil {
			return nil, err
		}

		// TODO(mitchellh): we need to verify the diff doesn't change
		// anything and that the diff has no computed values (pre-computed)

		// With the completed diff, apply!
		rs, err := r.Provider.Apply(r.State, diff)
		if err != nil {
			return nil, err
		}

		// If no state was returned, then no variables were updated so
		// just return.
		if rs == nil {
			return nil, nil
		}

		var errs []error
		for ak, av := range rs.Attributes {
			// If the value is the unknown variable value, then it is an error.
			// In this case we record the error and remove it from the state
			if av == config.UnknownVariableValue {
				errs = append(errs, fmt.Errorf(
					"Attribute with unknown value: %s", ak))
				delete(rs.Attributes, ak)
			}
		}

		// Update the resulting diff
		l.Lock()
		result.Resources[r.Id] = rs
		l.Unlock()

		// Determine the new state and update variables
		vars := make(map[string]string)
		for ak, av := range rs.Attributes {
			vars[fmt.Sprintf("%s.%s", r.Id, ak)] = av
		}

		err = nil
		if len(errs) > 0 {
			err = &MultiError{Errors: errs}
		}

		return vars, err
	}

	return t.genericWalkFn(p.Vars, cb)
}

func (t *Terraform) planWalkFn(
	state *State, result *Plan) depgraph.WalkFunc {
	var l sync.Mutex

	// Initialize the result diff so we can write to it
	result.init()

	// Write our configuration out
	//result.Config = t.config

	// Copy the variables
	/*
		result.Vars = make(map[string]string)
		for k, v := range t.variables {
			result.Vars[k] = v
		}
	*/

	cb := func(r *Resource) (map[string]string, error) {
		// Refresh the state so we're working with the latest resource info
		newState, err := r.Provider.Refresh(r.State)
		if err != nil {
			return nil, err
		}

		// Make sure the state is set to at the very least the empty state
		if newState == nil {
			newState = new(ResourceState)
		}

		// Set the type, the provider shouldn't modify this
		newState.Type = r.State.Type

		// Get a diff from the newest state
		diff, err := r.Provider.Diff(newState, r.Config)
		if err != nil {
			return nil, err
		}

		l.Lock()
		if !diff.Empty() {
			result.Diff.Resources[r.Id] = diff
		}
		result.State.Resources[r.Id] = newState
		l.Unlock()

		// Determine the new state and update variables
		vars := make(map[string]string)
		rs := newState
		if !diff.Empty() {
			rs = r.State.MergeDiff(diff)
		}
		if rs != nil {
			for ak, av := range rs.Attributes {
				vars[fmt.Sprintf("%s.%s", r.Id, ak)] = av
			}
		}

		return vars, nil
	}

	return t.genericWalkFn(nil, cb)
}

func (t *Terraform) genericWalkFn(
	invars map[string]string,
	cb genericWalkFunc) depgraph.WalkFunc {
	var l sync.Mutex

	// Initialize the variables for application
	vars := make(map[string]string)
	for k, v := range invars {
		vars[fmt.Sprintf("var.%s", k)] = v
	}

	return func(n *depgraph.Noun) error {
		// If it is the root node, ignore
		if n.Name == GraphRootNode {
			return nil
		}

		switch m := n.Meta.(type) {
		case *GraphNodeResource:
		case *GraphNodeResourceProvider:
			var rc *ResourceConfig
			if m.Config != nil {
				if err := m.Config.RawConfig.Interpolate(vars); err != nil {
					panic(err)
				}
				rc = NewResourceConfig(m.Config.RawConfig)
			}

			for k, p := range m.Providers {
				log.Printf("Configuring provider: %s", k)
				err := p.Configure(rc)
				if err != nil {
					return err
				}
			}

			return nil
		}

		rn := n.Meta.(*GraphNodeResource)
		if len(vars) > 0 && rn.Config != nil {
			if err := rn.Config.RawConfig.Interpolate(vars); err != nil {
				panic(fmt.Sprintf("Interpolate error: %s", err))
			}

			// Set the config
			rn.Resource.Config = NewResourceConfig(rn.Config.RawConfig)
		}

		// Call the callack
		newVars, err := cb(rn.Resource)
		if err != nil {
			return err
		}

		if len(newVars) > 0 {
			// Acquire a lock since this function is called in parallel
			l.Lock()
			defer l.Unlock()

			// Update variables
			for k, v := range newVars {
				vars[k] = v
			}
		}

		return nil
	}
}
