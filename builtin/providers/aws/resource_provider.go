package aws

import (
	"log"

	"github.com/hashicorp/terraform/helper/config"
	"github.com/hashicorp/terraform/helper/multierror"
	"github.com/hashicorp/terraform/terraform"
	"github.com/mitchellh/goamz/autoscaling"
	"github.com/mitchellh/goamz/ec2"
	"github.com/mitchellh/goamz/elb"
)

type ResourceProvider struct {
	Config Config

	ec2conn         *ec2.EC2
	elbconn         *elb.ELB
	autoscalingconn *autoscaling.AutoScaling
}

func (p *ResourceProvider) Validate(c *terraform.ResourceConfig) ([]string, []error) {
	return nil, nil
}

func (p *ResourceProvider) ValidateResource(
	t string, c *terraform.ResourceConfig) ([]string, []error) {
	return resourceMap.Validate(t, c)
}

func (p *ResourceProvider) Configure(c *terraform.ResourceConfig) error {
	if _, err := config.Decode(&p.Config, c.Config); err != nil {
		return err
	}

	// Get the auth and region. This can fail if keys/regions were not
	// specified and we're attempting to use the environment.
	var errs []error
	log.Println("Building AWS auth structure")
	auth, err := p.Config.AWSAuth()
	if err != nil {
		errs = append(errs, err)
	}

	log.Println("Building AWS region structure")
	region, err := p.Config.AWSRegion()
	if err != nil {
		errs = append(errs, err)
	}

	if len(errs) == 0 {
		log.Println("Initializing EC2 connection")
		p.ec2conn = ec2.New(auth, region)
		log.Println("Initializing ELB connection")
		p.elbconn = elb.New(auth, region)
		log.Println("Initializing AutoScaling connection")
		p.autoscalingconn = autoscaling.New(auth, region)
	}

	if len(errs) > 0 {
		return &multierror.Error{Errors: errs}
	}

	return nil
}

func (p *ResourceProvider) Apply(
	s *terraform.ResourceState,
	d *terraform.ResourceDiff) (*terraform.ResourceState, error) {
	return resourceMap.Apply(s, d, p)
}

func (p *ResourceProvider) Diff(
	s *terraform.ResourceState,
	c *terraform.ResourceConfig) (*terraform.ResourceDiff, error) {
	return resourceMap.Diff(s, c, p)
}

func (p *ResourceProvider) Refresh(
	s *terraform.ResourceState) (*terraform.ResourceState, error) {
	return resourceMap.Refresh(s, p)
}

func (p *ResourceProvider) Resources() []terraform.ResourceType {
	return resourceMap.Resources()
}
