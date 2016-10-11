package terraform

import (
	"sync"

	"github.com/hashicorp/terraform/config"
)

// ShadowEvalContext is an EvalContext that is used to "shadow" a real
// eval context for comparing whether two separate graph executions result
// in the same output.
//
// This eval context will never communicate with a real provider and will
// never modify real state.
type ShadowEvalContext interface {
	EvalContext

	// Close should be called when the _real_ EvalContext operations
	// are complete. This will immediately end any blocks calls and record
	// any errors.
	//
	// The returned error is the result of the shadow run. If it is nil,
	// then the shadow run seemingly completed successfully. You should
	// still compare the resulting states, diffs from both the real and shadow
	// contexts to verify equivalent end state.
	//
	// If the error is non-nil, then an error occurred during the execution
	// itself. In this scenario, you should not compare diffs/states since
	// they can't be considered accurate since operations during execution
	// failed.
	Close() error
}

// NewShadowEvalContext creates a new shadowed EvalContext. This returns
// the real EvalContext that should be used with the real evaluation and
// will communicate with real providers and write real state as well as
// the ShadowEvalContext that should be used with the test graph.
//
// This should be called before the ctx is ever used in order to ensure
// a consistent shadow state.
func NewShadowEvalContext(ctx EvalContext) (EvalContext, ShadowEvalContext) {
	real := &shadowEvalContextReal{EvalContext: ctx}

	// Copy the diff. We do this using some weird scoping so that the
	// "diff" (real) value never leaks out and can be used.
	var diffCopy *Diff
	{
		diff, lock := ctx.Diff()
		lock.RLock()
		diffCopy = diff
		// TODO: diffCopy = diff.DeepCopy()
		lock.RUnlock()
	}

	// Copy the state. We do this using some weird scoping so that the
	// "state" (real) value never leaks out and can be used.
	var stateCopy *State
	{
		state, lock := ctx.State()
		lock.RLock()
		stateCopy = state.DeepCopy()
		lock.RUnlock()
	}

	// Build the shadow copy. For safety, we don't even give the shadow
	// copy a reference to the real context. This means that it would be
	// very difficult (impossible without some real obvious mistakes) for
	// the shadow context to do "real" work.
	shadow := &shadowEvalContextShadow{
		PathValue:  ctx.Path(),
		StateValue: stateCopy,
		StateLock:  new(sync.RWMutex),
		DiffValue:  diffCopy,
		DiffLock:   new(sync.RWMutex),
	}

	return real, shadow
}

// shadowEvalContextReal is the EvalContext that does real work.
type shadowEvalContextReal struct {
	EvalContext
}

// shadowEvalContextShadow is the EvalContext that shadows the real one
// and leans on that for data.
type shadowEvalContextShadow struct {
	PathValue  []string
	Providers  map[string]ResourceProvider
	DiffValue  *Diff
	DiffLock   *sync.RWMutex
	StateValue *State
	StateLock  *sync.RWMutex

	// Fields relating to closing the context. Closing signals that
	// the execution of the real context completed.
	closeLock sync.Mutex
	closed    bool
	closeCh   chan struct{}
}

// Shared is the shared state between the shadow and real contexts when
// a shadow context is active. This is used by the real context to setup
// some state, trigger condition variables, etc.
type shadowEvalContextShared struct {
	// This lock must be held when modifying just about anything in this
	// structure. It is a "big" lock but the work done here is usually very
	// fast so we do this.
	sync.Mutex

	// Providers is the map of active (initialized) providers.
	//
	// ProviderWaiters is the condition variable associated with waiting
	// for a provider to be initialized. If this is non-nil for a provider,
	// then that will be triggered when the provider is initialized.
	Providers       map[string]struct{}
	ProviderWaiters map[string]*sync.Cond
}

func (c *shadowEvalContextShadow) Close() error {
	// TODO
	return nil
}

func (c *shadowEvalContextShadow) Path() []string {
	return c.PathValue
}

func (c *shadowEvalContextShadow) Hook(f func(Hook) (HookAction, error)) error {
	// Don't do anything on hooks. Mission critical behavior should not
	// depend on hooks and at the time of writing it does not depend on
	// hooks. In the future we could also test hooks but not now.
	return nil
}

func (c *shadowEvalContextShadow) Input() UIInput {
	// TODO
	return nil
}

func (c *shadowEvalContextShadow) InitProvider(n string) (ResourceProvider, error) {
	// Initialize our shadow provider here. We also wait for the
	// real context to initialize the same provider. If it doesn't
	// before close, then an error is reported.

	// TODO: shadow provider

	return nil, nil
}

func (c *shadowEvalContextShadow) Diff() (*Diff, *sync.RWMutex) {
	return c.DiffValue, c.DiffLock
}

func (c *shadowEvalContextShadow) State() (*State, *sync.RWMutex) {
	return c.StateValue, c.StateLock
}

func (c *shadowEvalContextShadow) Provider(n string) ResourceProvider              { return nil }
func (c *shadowEvalContextShadow) CloseProvider(n string) error                    { return nil }
func (c *shadowEvalContextShadow) ConfigureProvider(string, *ResourceConfig) error { return nil }
func (c *shadowEvalContextShadow) SetProviderConfig(string, *ResourceConfig) error { return nil }
func (c *shadowEvalContextShadow) ParentProviderConfig(string) *ResourceConfig     { return nil }
func (c *shadowEvalContextShadow) ProviderInput(string) map[string]interface{}     { return nil }
func (c *shadowEvalContextShadow) SetProviderInput(string, map[string]interface{}) {}
func (c *shadowEvalContextShadow) InitProvisioner(string) (ResourceProvisioner, error) {
	return nil, nil
}
func (c *shadowEvalContextShadow) Provisioner(string) ResourceProvisioner { return nil }
func (c *shadowEvalContextShadow) CloseProvisioner(string) error          { return nil }
func (c *shadowEvalContextShadow) Interpolate(*config.RawConfig, *Resource) (*ResourceConfig, error) {
	return nil, nil
}
func (c *shadowEvalContextShadow) SetVariables(string, map[string]interface{}) {}
