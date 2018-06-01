package terraform

import (
	"fmt"
	"log"

	"github.com/hashicorp/hcl2/hcl"

	"github.com/hashicorp/terraform/addrs"
	"github.com/hashicorp/terraform/configs"
	"github.com/hashicorp/terraform/tfdiags"
)

func buildProviderConfig(ctx EvalContext, addr addrs.ProviderConfig, config *configs.Provider) hcl.Body {
	var configBody hcl.Body
	if config != nil {
		configBody = config.Config
	}

	var inputBody hcl.Body
	inputConfig := ctx.ProviderInput(addr)
	if len(inputConfig) > 0 {
		inputBody = configs.SynthBody("<input-prompt>", inputConfig)
	}

	switch {
	case configBody != nil && inputBody != nil:
		log.Printf("[TRACE] buildProviderConfig for %s: merging explicit config and input", addr)
		// Note that the inputBody is the _base_ here, because configs.MergeBodies
		// expects the base have all of the required fields, while these are
		// forced to be optional for the override. The input process should
		// guarantee that we have a value for each of the required arguments and
		// that in practice the sets of attributes in each body will be
		// disjoint.
		return configs.MergeBodies(inputBody, configBody)
	case configBody != nil:
		log.Printf("[TRACE] buildProviderConfig for %s: using explicit config only", addr)
		return configBody
	case inputBody != nil:
		log.Printf("[TRACE] buildProviderConfig for %s: using input only", addr)
		return inputBody
	default:
		log.Printf("[TRACE] buildProviderConfig for %s: no configuration at all", addr)
		return hcl.EmptyBody()
	}
}

// EvalConfigProvider is an EvalNode implementation that configures
// a provider that is already initialized and retrieved.
type EvalConfigProvider struct {
	Addr     addrs.ProviderConfig
	Provider *ResourceProvider
	Config   *configs.Provider
}

func (n *EvalConfigProvider) Eval(ctx EvalContext) (interface{}, error) {
	if n.Provider == nil {
		return nil, fmt.Errorf("EvalConfigProvider Provider is nil")
	}

	var diags tfdiags.Diagnostics
	provider := *n.Provider
	config := n.Config

	configBody := buildProviderConfig(ctx, n.Addr, config)

	schema, err := provider.GetSchema(&ProviderSchemaRequest{})
	if err != nil {
		diags = diags.Append(err)
		return nil, diags.NonFatalErr()
	}
	if schema == nil {
		return nil, fmt.Errorf("schema not available for %s", n.Addr)
	}

	configSchema := schema.Provider
	configVal, configBody, evalDiags := ctx.EvaluateBlock(configBody, configSchema, nil, addrs.NoKey)
	diags = diags.Append(evalDiags)
	if evalDiags.HasErrors() {
		return nil, diags.NonFatalErr()
	}

	configDiags := ctx.ConfigureProvider(n.Addr, configVal)
	configDiags = configDiags.InConfigBody(configBody)

	return nil, configDiags.ErrWithWarnings()
}

// EvalInitProvider is an EvalNode implementation that initializes a provider
// and returns nothing. The provider can be retrieved again with the
// EvalGetProvider node.
type EvalInitProvider struct {
	TypeName string
	Addr     addrs.ProviderConfig
}

func (n *EvalInitProvider) Eval(ctx EvalContext) (interface{}, error) {
	return ctx.InitProvider(n.TypeName, n.Addr)
}

// EvalCloseProvider is an EvalNode implementation that closes provider
// connections that aren't needed anymore.
type EvalCloseProvider struct {
	Addr addrs.ProviderConfig
}

func (n *EvalCloseProvider) Eval(ctx EvalContext) (interface{}, error) {
	ctx.CloseProvider(n.Addr)
	return nil, nil
}

// EvalGetProvider is an EvalNode implementation that retrieves an already
// initialized provider instance for the given name.
//
// Unlike most eval nodes, this takes an _absolute_ provider configuration,
// because providers can be passed into and inherited between modules.
// Resource nodes must therefore know the absolute path of the provider they
// will use, which is usually accomplished by implementing
// interface GraphNodeProviderConsumer.
type EvalGetProvider struct {
	Addr   addrs.AbsProviderConfig
	Output *ResourceProvider

	// If non-nil, Schema will be updated after eval to refer to the
	// schema of the provider.
	Schema **ProviderSchema
}

func (n *EvalGetProvider) Eval(ctx EvalContext) (interface{}, error) {
	result := ctx.Provider(n.Addr)
	if result == nil {
		return nil, fmt.Errorf("provider %s not initialized", n.Addr)
	}

	if n.Output != nil {
		*n.Output = result
	}

	if n.Schema != nil {
		*n.Schema = ctx.ProviderSchema(n.Addr)
	}

	return nil, nil
}

// EvalInputProvider is an EvalNode implementation that asks for input
// for the given provider configurations.
type EvalInputProvider struct {
	Addr     addrs.ProviderConfig
	Provider *ResourceProvider
	Config   *configs.Provider
}

func (n *EvalInputProvider) Eval(ctx EvalContext) (interface{}, error) {
	// This is currently disabled. It used to interact with a provider method
	// called Input, allowing the provider to capture input interactively
	// itself, but once re-implemented we'll have this instead use the
	// provider's configuration schema to automatically infer what we need
	// to prompt for.
	var diags tfdiags.Diagnostics
	diag := &hcl.Diagnostic{
		Severity: hcl.DiagWarning,
		Summary:  "Provider input is temporarily disabled",
		Detail:   fmt.Sprintf("Skipped gathering input for %s because the input step is currently disabled pending a change to the provider API.", n.Addr),
	}
	if n.Config != nil {
		diag.Subject = n.Config.DeclRange.Ptr()
	}
	diags = diags.Append(diag)
	return nil, diags.ErrWithWarnings()
}
