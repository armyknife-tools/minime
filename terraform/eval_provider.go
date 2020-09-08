package terraform

import (
	"fmt"
	"log"

	"github.com/hashicorp/hcl/v2"

	"github.com/hashicorp/terraform/addrs"
	"github.com/hashicorp/terraform/configs"
	"github.com/hashicorp/terraform/providers"
	"github.com/hashicorp/terraform/tfdiags"
)

func buildProviderConfig(ctx EvalContext, addr addrs.AbsProviderConfig, config *configs.Provider) hcl.Body {
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

// GetProvider returns the providers interface and schema for a given provider.
func GetProvider(ctx EvalContext, addr addrs.AbsProviderConfig) (providers.Interface, *ProviderSchema, error) {
	if addr.Provider.Type == "" {
		// Should never happen
		panic("EvalGetProvider used with uninitialized provider configuration address")
	}
	provider := ctx.Provider(addr)
	if provider == nil {
		return nil, &ProviderSchema{}, fmt.Errorf("provider %s not initialized", addr)
	}
	// Not all callers require a schema, so we will leave checking for a nil
	// schema to the callers.
	schema := ctx.ProviderSchema(addr)
	return provider, schema, nil
}

// EvalConfigProvider is an EvalNode implementation that configures
// a provider that is already initialized and retrieved.
type EvalConfigProvider struct {
	Addr                addrs.AbsProviderConfig
	Provider            *providers.Interface
	Config              *configs.Provider
	VerifyConfigIsKnown bool
}

func (n *EvalConfigProvider) Eval(ctx EvalContext) (interface{}, error) {
	if n.Provider == nil {
		return nil, fmt.Errorf("EvalConfigProvider Provider is nil")
	}

	var diags tfdiags.Diagnostics
	provider := *n.Provider
	config := n.Config

	configBody := buildProviderConfig(ctx, n.Addr, config)

	resp := provider.GetSchema()
	diags = diags.Append(resp.Diagnostics)
	if diags.HasErrors() {
		return nil, diags.NonFatalErr()
	}

	configSchema := resp.Provider.Block
	configVal, configBody, evalDiags := ctx.EvaluateBlock(configBody, configSchema, nil, EvalDataForNoInstanceKey)
	diags = diags.Append(evalDiags)
	if evalDiags.HasErrors() {
		return nil, diags.NonFatalErr()
	}

	if n.VerifyConfigIsKnown && !configVal.IsWhollyKnown() {
		diags = diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Invalid provider configuration",
			Detail:   fmt.Sprintf("The configuration for %s depends on values that cannot be determined until apply.", n.Addr),
			Subject:  &config.DeclRange,
		})
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
	Addr addrs.AbsProviderConfig
}

func (n *EvalInitProvider) Eval(ctx EvalContext) (interface{}, error) {
	return ctx.InitProvider(n.Addr)
}

// EvalCloseProvider is an EvalNode implementation that closes provider
// connections that aren't needed anymore.
type EvalCloseProvider struct {
	Addr addrs.AbsProviderConfig
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
	Output *providers.Interface

	// If non-nil, Schema will be updated after eval to refer to the
	// schema of the provider.
	Schema **ProviderSchema
}

func (n *EvalGetProvider) Eval(ctx EvalContext) (interface{}, error) {
	if n.Addr.Provider.Type == "" {
		// Should never happen
		panic("EvalGetProvider used with uninitialized provider configuration address")
	}

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
