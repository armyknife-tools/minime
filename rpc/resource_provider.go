package rpc

import (
	"net/rpc"

	"github.com/hashicorp/terraform/terraform"
)

// ResourceProvider is an implementation of terraform.ResourceProvider
// that communicates over RPC.
type ResourceProvider struct {
	Client *rpc.Client
	Name   string
}

func (p *ResourceProvider) Validate(c *terraform.ResourceConfig) ([]string, []error) {
	var resp ResourceProviderValidateResponse
	err := p.Client.Call(p.Name+".Validate", c, &resp)
	if err != nil {
		return nil, []error{err}
	}

	var errs []error
	if len(resp.Errors) > 0 {
		errs = make([]error, len(resp.Errors))
		for i, err := range resp.Errors {
			errs[i] = err
		}
	}

	return resp.Warnings, errs
}

func (p *ResourceProvider) Configure(c *terraform.ResourceConfig) error {
	var resp ResourceProviderConfigureResponse
	err := p.Client.Call(p.Name+".Configure", c, &resp)
	if err != nil {
		return err
	}
	if resp.Error != nil {
		err = resp.Error
	}

	return err
}

func (p *ResourceProvider) Diff(
	s *terraform.ResourceState,
	c *terraform.ResourceConfig) (*terraform.ResourceDiff, error) {
	var resp ResourceProviderDiffResponse
	args := &ResourceProviderDiffArgs{
		State:  s,
		Config: c,
	}
	err := p.Client.Call(p.Name+".Diff", args, &resp)
	if err != nil {
		return nil, err
	}
	if resp.Error != nil {
		err = resp.Error
	}

	return resp.Diff, err
}

func (p *ResourceProvider) Resources() []terraform.ResourceType {
	var result []terraform.ResourceType

	err := p.Client.Call(p.Name+".Resources", new(interface{}), &result)
	if err != nil {
		// TODO: panic, log, what?
		return nil
	}

	return result
}

// ResourceProviderServer is a net/rpc compatible structure for serving
// a ResourceProvider. This should not be used directly.
type ResourceProviderServer struct {
	Provider terraform.ResourceProvider
}

type ResourceProviderConfigureResponse struct {
	Error *BasicError
}

type ResourceProviderDiffArgs struct {
	State  *terraform.ResourceState
	Config *terraform.ResourceConfig
}

type ResourceProviderDiffResponse struct {
	Diff  *terraform.ResourceDiff
	Error *BasicError
}

type ResourceProviderValidateResponse struct {
	Warnings []string
	Errors   []*BasicError
}

func (s *ResourceProviderServer) Validate(
	config *terraform.ResourceConfig,
	reply *ResourceProviderValidateResponse) error {
	warns, errs := s.Provider.Validate(config)
	berrs := make([]*BasicError, len(errs))
	for i, err := range errs {
		berrs[i] = NewBasicError(err)
	}
	*reply = ResourceProviderValidateResponse{
		Warnings: warns,
		Errors:   berrs,
	}
	return nil
}

func (s *ResourceProviderServer) Configure(
	config *terraform.ResourceConfig,
	reply *ResourceProviderConfigureResponse) error {
	err := s.Provider.Configure(config)
	*reply = ResourceProviderConfigureResponse{
		Error: NewBasicError(err),
	}
	return nil
}

func (s *ResourceProviderServer) Diff(
	args *ResourceProviderDiffArgs,
	result *ResourceProviderDiffResponse) error {
	diff, err := s.Provider.Diff(args.State, args.Config)
	*result = ResourceProviderDiffResponse{
		Diff:  diff,
		Error: NewBasicError(err),
	}
	return nil
}

func (s *ResourceProviderServer) Resources(
	nothing interface{},
	result *[]terraform.ResourceType) error {
	*result = s.Provider.Resources()
	return nil
}
