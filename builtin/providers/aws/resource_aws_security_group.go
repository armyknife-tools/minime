package aws

import (
	"bytes"
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/mitchellh/goamz/ec2"
)

func resourceAwsSecurityGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAwsSecurityGroupCreate,
		Read:   resourceAwsSecurityGroupRead,
		Update: resourceAwsSecurityGroupUpdate,
		Delete: resourceAwsSecurityGroupDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"description": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"vpc_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"ingress": &schema.Schema{
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"from_port": &schema.Schema{
							Type:     schema.TypeInt,
							Required: true,
						},

						"to_port": &schema.Schema{
							Type:     schema.TypeInt,
							Required: true,
						},

						"protocol": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},

						"cidr_blocks": &schema.Schema{
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},

						"security_groups": &schema.Schema{
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
				Set: resourceAwsSecurityGroupIngressHash,
			},

			"owner_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAwsSecurityGroupIngressHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	buf.WriteString(fmt.Sprintf("%d-", m["from_port"].(int)))
	buf.WriteString(fmt.Sprintf("%d-", m["to_port"].(int)))
	buf.WriteString(fmt.Sprintf("%d-", m["protocol"].(string)))

	// We need to make sure to sort the strings below so that we always
	// generate the same hash code no matter what is in the set.
	if v, ok := m["cidr_blocks"]; ok {
		vs := v.([]interface{})
		s := make([]string, len(vs))
		for i, raw := range vs {
			s[i] = raw.(string)
		}
		sort.Strings(s)

		for _, v := range s {
			buf.WriteString(fmt.Sprintf("%s-", v))
		}
	}
	if v, ok := m["security_groups"]; ok {
		vs := v.([]interface{})
		s := make([]string, len(vs))
		for i, raw := range vs {
			s[i] = raw.(string)
		}
		sort.Strings(s)

		for _, v := range s {
			buf.WriteString(fmt.Sprintf("%s-", v))
		}
	}

	return hashcode.String(buf.String())
}

func resourceAwsSecurityGroupCreate(d *schema.ResourceData, meta interface{}) error {
	p := meta.(*ResourceProvider)
	ec2conn := p.ec2conn

	securityGroupOpts := ec2.SecurityGroup{
		Name: d.Get("name").(string),
	}

	if v := d.Get("vpc_id"); v != nil {
		securityGroupOpts.VpcId = v.(string)
	}

	if v := d.Get("description"); v != nil {
		securityGroupOpts.Description = v.(string)
	}

	log.Printf(
		"[DEBUG] Security Group create configuration: %#v", securityGroupOpts)
	createResp, err := ec2conn.CreateSecurityGroup(securityGroupOpts)
	if err != nil {
		return fmt.Errorf("Error creating Security Group: %s", err)
	}

	d.SetId(createResp.Id)
	group := createResp.SecurityGroup

	log.Printf("[INFO] Security Group ID: %s", d.Id())

	// Wait for the security group to truly exist
	log.Printf(
		"[DEBUG] Waiting for Security Group (%s) to exist",
		d.Id())
	stateConf := &resource.StateChangeConf{
		Pending: []string{""},
		Target:  "exists",
		Refresh: SGStateRefreshFunc(ec2conn, d.Id()),
		Timeout: 1 * time.Minute,
	}
	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf(
			"Error waiting for Security Group (%s) to become available: %s",
			d.Id(), err)
	}

	// Expand the "ingress" array to goamz compat []ec2.IPPerm
	ingressRaw := d.Get("ingress")
	if ingressRaw == nil {
		ingressRaw = new(schema.Set)
	}
	ingressList := ingressRaw.(*schema.Set).List()
	if len(ingressList) > 0 {
		ingressRules := expandIPPerms(ingressList)
		_, err = ec2conn.AuthorizeSecurityGroup(group, ingressRules)
		if err != nil {
			return fmt.Errorf("Error authorizing security group ingress rules: %s", err)
		}
	}

	return resourceAwsSecurityGroupRead(d, meta)
}

func resourceAwsSecurityGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	p := meta.(*ResourceProvider)
	ec2conn := p.ec2conn

	sgRaw, _, err := SGStateRefreshFunc(ec2conn, d.Id())()
	if err != nil {
		return err
	}
	if sgRaw == nil {
		d.SetId("")
		return nil
	}
	group := sgRaw.(*ec2.SecurityGroupInfo).SecurityGroup

	if d.HasChange("ingress") {
		o, n := d.GetChange("ingress")
		if o == nil {
			o = new(schema.Set)
		}
		if n == nil {
			n = new(schema.Set)
		}

		os := o.(*schema.Set)
		ns := n.(*schema.Set)

		remove := expandIPPerms(os.Difference(ns).List())
		add := expandIPPerms(ns.Difference(os).List())

		// TODO: We need to handle partial state better in the in-between
		// in this update.

		// TODO: It'd be nicer to authorize before removing, but then we have
		// to deal with complicated unrolling to get individual CIDR blocks
		// to avoid authorizing already authorized sources. Removing before
		// adding is easier here, and Terraform should be fast enough to
		// not have service issues.

		if len(remove) > 0 {
			// Revoke the old rules
			_, err = ec2conn.RevokeSecurityGroup(group, remove)
			if err != nil {
				return fmt.Errorf("Error authorizing security group ingress rules: %s", err)
			}
		}

		if len(add) > 0 {
			// Authorize the new rules
			_, err := ec2conn.AuthorizeSecurityGroup(group, add)
			if err != nil {
				return fmt.Errorf("Error authorizing security group ingress rules: %s", err)
			}
		}
	}

	return nil
}

func resourceAwsSecurityGroupDelete(d *schema.ResourceData, meta interface{}) error {
	p := meta.(*ResourceProvider)
	ec2conn := p.ec2conn

	log.Printf("[DEBUG] Security Group destroy: %v", d.Id())

	_, err := ec2conn.DeleteSecurityGroup(ec2.SecurityGroup{Id: d.Id()})
	if err != nil {
		ec2err, ok := err.(*ec2.Error)
		if ok && ec2err.Code == "InvalidGroup.NotFound" {
			return nil
		}
	}

	return err
}

func resourceAwsSecurityGroupRead(d *schema.ResourceData, meta interface{}) error {
	p := meta.(*ResourceProvider)
	ec2conn := p.ec2conn

	sgRaw, _, err := SGStateRefreshFunc(ec2conn, d.Id())()
	if err != nil {
		return err
	}
	if sgRaw == nil {
		d.SetId("")
		return nil
	}

	sg := sgRaw.(*ec2.SecurityGroupInfo)

	var deps []terraform.ResourceDependency

	// Gather our ingress rules
	ingressRules := make([]map[string]interface{}, len(sg.IPPerms))
	for i, perm := range sg.IPPerms {
		n := make(map[string]interface{})
		n["from_port"] = perm.FromPort
		n["protocol"] = perm.Protocol
		n["to_port"] = perm.ToPort

		if len(perm.SourceIPs) > 0 {
			n["cidr_blocks"] = perm.SourceIPs
		}

		if len(perm.SourceGroups) > 0 {
			// We depend on other security groups
			for _, v := range perm.SourceGroups {
				deps = append(deps,
					terraform.ResourceDependency{ID: v.Id},
				)
			}

			n["security_groups"] = flattenSecurityGroups(perm.SourceGroups)
		}

		ingressRules[i] = n
	}

	if v := d.Get("vpc_id"); v != nil && v.(string) != "" {
		deps = append(deps,
			terraform.ResourceDependency{ID: v.(string)},
		)
	}

	d.Set("description", sg.Description)
	d.Set("name", sg.Name)
	d.Set("vpc_id", sg.VpcId)
	d.Set("owner_id", sg.OwnerId)
	d.Set("ingress", ingressRules)
	d.SetDependencies(deps)

	return nil
}

// SGStateRefreshFunc returns a resource.StateRefreshFunc that is used to watch
// a security group.
func SGStateRefreshFunc(conn *ec2.EC2, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		sgs := []ec2.SecurityGroup{ec2.SecurityGroup{Id: id}}
		resp, err := conn.SecurityGroups(sgs, nil)
		if err != nil {
			if ec2err, ok := err.(*ec2.Error); ok {
				if ec2err.Code == "InvalidSecurityGroupID.NotFound" ||
					ec2err.Code == "InvalidGroup.NotFound" {
					resp = nil
					err = nil
				}
			}

			if err != nil {
				log.Printf("Error on SGStateRefresh: %s", err)
				return nil, "", err
			}
		}

		if resp == nil {
			return nil, "", nil
		}

		group := &resp.Groups[0]
		return group, "exists", nil
	}
}
