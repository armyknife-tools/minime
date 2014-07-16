package aws

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/flatmap"
	"github.com/hashicorp/terraform/helper/config"
	"github.com/hashicorp/terraform/helper/diff"
	"github.com/hashicorp/terraform/terraform"
	"github.com/mitchellh/goamz/elb"
)

func resource_aws_elb_create(
	s *terraform.ResourceState,
	d *terraform.ResourceDiff,
	meta interface{}) (*terraform.ResourceState, error) {
	p := meta.(*ResourceProvider)
	elbconn := p.elbconn

	// Merge the diff into the state so that we have all the attributes
	// properly.
	rs := s.MergeDiff(d)

	// The name specified for the ELB. This is also our unique ID
	// we save to state if the creation is succesful (amazon verifies
	// it is unique)
	elbName := rs.Attributes["name"]

	// Expand the "listener" array to goamz compat []elb.Listener
	v := flatmap.Expand(rs.Attributes, "listener").([]interface{})
	listeners := expandListeners(v)

	v = flatmap.Expand(rs.Attributes, "availability_zones").([]interface{})
	zones := expandStringList(v)

	// Provision the elb
	elbOpts := &elb.CreateLoadBalancer{
		LoadBalancerName: elbName,
		Listeners:        listeners,
		AvailZone:        zones,
	}

	log.Printf("[DEBUG] ELB create configuration: %#v", elbOpts)

	_, err := elbconn.CreateLoadBalancer(elbOpts)
	if err != nil {
		return nil, fmt.Errorf("Error creating ELB: %s", err)
	}

	// Assign the elb's unique identifer for use later
	rs.ID = elbName
	log.Printf("[INFO] ELB ID: %s", elbName)

	// If we have any instances, we need to register them
	v = flatmap.Expand(rs.Attributes, "instances").([]interface{})
	instances := expandStringList(v)

	if len(instances) > 0 {
		registerInstancesOpts := elb.RegisterInstancesWithLoadBalancer{
			LoadBalancerName: elbName,
			Instances:        instances,
		}

		_, err := elbconn.RegisterInstancesWithLoadBalancer(&registerInstancesOpts)

		if err != nil {
			return rs, fmt.Errorf("Failure registering instances: %s", err)
		}
	}

	loadBalancer, err := resource_aws_elb_retrieve_balancer(rs.ID, elbconn)
	if err != nil {
		return rs, err
	}

	return resource_aws_elb_update_state(rs, loadBalancer)
}

func resource_aws_elb_update(
	s *terraform.ResourceState,
	d *terraform.ResourceDiff,
	meta interface{}) (*terraform.ResourceState, error) {
	p := meta.(*ResourceProvider)
	elbconn := p.elbconn

	rs := s.MergeDiff(d)

	// If we currently have instances, or did have instances,
	// we want to figure out what to add and remove from the load
	// balancer
	if attr, ok := d.Attributes["instances.#"]; ok && attr.Old != "" {
		// The new state of instances merged with the diff
		mergedInstances := expandStringList(flatmap.Expand(
			rs.Attributes, "instances").([]interface{}))

		// The state before the diff merge
		previousInstances := expandStringList(flatmap.Expand(
			s.Attributes, "instances").([]interface{}))

		// keep track of what instances we are removing, and which
		// we are adding
		var toRemove []string
		var toAdd []string

		for _, instanceId := range mergedInstances {
			for _, prevId := range previousInstances {
				// If the merged instance ID existed
				// previously, we don't have to do anything
				if instanceId == prevId {
					continue
					// Otherwise, we need to add it to the load balancer
				} else {
					toAdd = append(toAdd, instanceId)
				}
			}
		}

		for i, instanceId := range toAdd {
			for _, prevId := range previousInstances {
				// If the instance ID we are adding existed
				// previously, we want to not add it, but rather remove
				// it
				if instanceId == prevId {
					toRemove = append(toRemove, instanceId)
					toAdd = append(toAdd[:i], toAdd[i+1:]...)
					// Otherwise, we continue adding it to the ELB
				} else {
					continue
				}
			}
		}

		if len(toAdd) > 0 {
			registerInstancesOpts := elb.RegisterInstancesWithLoadBalancer{
				LoadBalancerName: rs.ID,
				Instances:        toAdd,
			}

			_, err := elbconn.RegisterInstancesWithLoadBalancer(&registerInstancesOpts)

			if err != nil {
				return s, fmt.Errorf("Failure registering instances: %s", err)
			}
		}

		if len(toRemove) > 0 {
			deRegisterInstancesOpts := elb.DeregisterInstancesFromLoadBalancer{
				LoadBalancerName: rs.ID,
				Instances:        toRemove,
			}

			_, err := elbconn.DeregisterInstancesFromLoadBalancer(&deRegisterInstancesOpts)

			if err != nil {
				return s, fmt.Errorf("Failure deregistering instances: %s", err)
			}
		}
	}

	loadBalancer, err := resource_aws_elb_retrieve_balancer(rs.ID, elbconn)

	if err != nil {
		return s, err
	}

	return resource_aws_elb_update_state(rs, loadBalancer)
}

func resource_aws_elb_destroy(
	s *terraform.ResourceState,
	meta interface{}) error {
	p := meta.(*ResourceProvider)
	elbconn := p.elbconn

	log.Printf("[INFO] Deleting ELB: %s", s.ID)

	// Destroy the load balancer
	deleteElbOpts := elb.DeleteLoadBalancer{
		LoadBalancerName: s.ID,
	}
	_, err := elbconn.DeleteLoadBalancer(&deleteElbOpts)

	if err != nil {
		return fmt.Errorf("Error deleting ELB: %s", err)
	}

	return nil
}

func resource_aws_elb_refresh(
	s *terraform.ResourceState,
	meta interface{}) (*terraform.ResourceState, error) {
	p := meta.(*ResourceProvider)
	elbconn := p.elbconn

	loadBalancer, err := resource_aws_elb_retrieve_balancer(s.ID, elbconn)
	if err != nil {
		return nil, err
	}

	return resource_aws_elb_update_state(s, loadBalancer)
}

func resource_aws_elb_diff(
	s *terraform.ResourceState,
	c *terraform.ResourceConfig,
	meta interface{}) (*terraform.ResourceDiff, error) {

	b := &diff.ResourceBuilder{
		Attrs: map[string]diff.AttrType{
			"name":              diff.AttrTypeCreate,
			"availability_zone": diff.AttrTypeCreate,
			"listener":          diff.AttrTypeCreate,
			"instances":         diff.AttrTypeUpdate,
		},

		ComputedAttrs: []string{
			"dns_name",
		},
	}

	return b.Diff(s, c)
}

func resource_aws_elb_update_state(
	s *terraform.ResourceState,
	balancer *elb.LoadBalancer) (*terraform.ResourceState, error) {

	s.Attributes["name"] = balancer.LoadBalancerName
	s.Attributes["dns_name"] = balancer.DNSName

	// Flatten our group values
	toFlatten := make(map[string]interface{})

	if len(balancer.Instances) > 0 && balancer.Instances[0].InstanceId != "" {
		toFlatten["instances"] = flattenInstances(balancer.Instances)
		for k, v := range flatmap.Flatten(toFlatten) {
			s.Attributes[k] = v
		}
	}

	return s, nil
}

// retrieves an ELB by it's ID
func resource_aws_elb_retrieve_balancer(id string, elbconn *elb.ELB) (*elb.LoadBalancer, error) {
	describeElbOpts := &elb.DescribeLoadBalancer{
		Names: []string{id},
	}

	// Retrieve the ELB properties for updating the state
	describeResp, err := elbconn.DescribeLoadBalancers(describeElbOpts)
	if err != nil {
		return nil, fmt.Errorf("Error retrieving ELB: %s", err)
	}

	loadBalancer := describeResp.LoadBalancers[0]

	// Verify AWS returned our ELB
	if len(describeResp.LoadBalancers) != 1 ||
		describeResp.LoadBalancers[0].LoadBalancerName != id {
		if err != nil {
			return nil, fmt.Errorf("Unable to find ELB: %#v", describeResp.LoadBalancers)
		}
	}

	return &loadBalancer, nil
}

func resource_aws_elb_validation() *config.Validator {
	return &config.Validator{
		Required: []string{
			"name",
			"availability_zones.*",
			"listener.*",
			"listener.*.instance_port",
			"listener.*.instance_protocol",
			"listener.*.lb_port",
			"listener.*.lb_protocol",
		},
		Optional: []string{
			"instances.*",
		},
	}
}
