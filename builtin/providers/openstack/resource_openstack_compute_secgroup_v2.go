package openstack

import (
	"bytes"
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/racker/perigee"
	"github.com/rackspace/gophercloud/openstack/compute/v2/extensions/secgroups"
)

func resourceComputeSecGroupV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceComputeSecGroupV2Create,
		Read:   resourceComputeSecGroupV2Read,
		Update: resourceComputeSecGroupV2Update,
		Delete: resourceComputeSecGroupV2Delete,

		Schema: map[string]*schema.Schema{
			"region": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				DefaultFunc: envDefaultFunc("OS_REGION_NAME"),
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"rule": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"from_port": &schema.Schema{
							Type:     schema.TypeInt,
							Required: true,
							ForceNew: false,
						},
						"to_port": &schema.Schema{
							Type:     schema.TypeInt,
							Required: true,
							ForceNew: false,
						},
						"ip_protocol": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
							ForceNew: false,
						},
						"cidr": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: false,
						},
						"from_group_id": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: false,
						},
					},
				},
			},
		},
	}
}

func resourceComputeSecGroupV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	computeClient, err := config.computeV2Client(d.Get("region").(string))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack compute client: %s", err)
	}

	createOpts := secgroups.CreateOpts{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}

	sg, err := secgroups.Create(computeClient, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error creating OpenStack security group: %s", err)
	}

	d.SetId(sg.ID)

	createRuleOptsList := resourceSecGroupRulesV2(d)
	for _, createRuleOpts := range createRuleOptsList {
		_, err := secgroups.CreateRule(computeClient, createRuleOpts).Extract()
		if err != nil {
			return fmt.Errorf("Error creating OpenStack security group rule: %s", err)
		}
	}

	return resourceComputeSecGroupV2Read(d, meta)
}

func resourceComputeSecGroupV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	computeClient, err := config.computeV2Client(d.Get("region").(string))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack compute client: %s", err)
	}

	sg, err := secgroups.Get(computeClient, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "security group")
	}

	d.Set("region", d.Get("region").(string))
	d.Set("name", sg.Name)
	d.Set("description", sg.Description)
	log.Printf("[DEBUG] rulesToMap(sg.Rules): %+v", rulesToMap(sg.Rules))
	d.Set("rule", rulesToMap(sg.Rules))

	return nil
}

func resourceComputeSecGroupV2Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	computeClient, err := config.computeV2Client(d.Get("region").(string))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack compute client: %s", err)
	}

	updateOpts := secgroups.UpdateOpts{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}

	log.Printf("[DEBUG] Updating Security Group (%s) with options: %+v", d.Id(), updateOpts)

	_, err = secgroups.Update(computeClient, d.Id(), updateOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error updating OpenStack security group (%s): %s", d.Id(), err)
	}

	if d.HasChange("rule") {
		oldSGRaw, newSGRaw := d.GetChange("rule")
		oldSGRSlice, newSGRSlice := oldSGRaw.([]interface{}), newSGRaw.([]interface{})
		oldSGRSet := schema.NewSet(secgroupRuleV2Hash, oldSGRSlice)
		newSGRSet := schema.NewSet(secgroupRuleV2Hash, newSGRSlice)
		secgrouprulesToAdd := newSGRSet.Difference(oldSGRSet)
		secgrouprulesToRemove := oldSGRSet.Difference(newSGRSet)

		log.Printf("[DEBUG] Security group rules to add: %v", secgrouprulesToAdd)

		log.Printf("[DEBUG] Security groups rules to remove: %v", secgrouprulesToRemove)

		for _, rawRule := range secgrouprulesToAdd.List() {
			createRuleOpts := resourceSecGroupRuleCreateOptsV2(d, rawRule)
			rule, err := secgroups.CreateRule(computeClient, createRuleOpts).Extract()
			if err != nil {
				return fmt.Errorf("Error adding rule to OpenStack security group (%s): %s", d.Id(), err)
			}
			log.Printf("[DEBUG] Added rule (%s) to OpenStack security group (%s) ", rule.ID, d.Id())
		}

		for _, r := range secgrouprulesToRemove.List() {
			rule := resourceSecGroupRuleV2(d, r)
			err := secgroups.DeleteRule(computeClient, rule.ID).ExtractErr()
			if err != nil {
				errCode, ok := err.(*perigee.UnexpectedResponseCodeError)
				if !ok {
					return fmt.Errorf("Error removing rule (%s) from OpenStack security group (%s): %s", rule.ID, d.Id(), err)
				}
				if errCode.Actual == 404 {
					continue
				} else {
					return fmt.Errorf("Error removing rule (%s) from OpenStack security group (%s)", rule.ID, d.Id())
				}
			} else {
				log.Printf("[DEBUG] Removed rule (%s) from OpenStack security group (%s): %s", rule.ID, d.Id(), err)
			}
		}
	}

	return resourceComputeSecGroupV2Read(d, meta)
}

func resourceComputeSecGroupV2Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	computeClient, err := config.computeV2Client(d.Get("region").(string))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack compute client: %s", err)
	}

	err = secgroups.Delete(computeClient, d.Id()).ExtractErr()
	if err != nil {
		return fmt.Errorf("Error deleting OpenStack security group: %s", err)
	}
	d.SetId("")
	return nil
}

func resourceSecGroupRulesV2(d *schema.ResourceData) []secgroups.CreateRuleOpts {
	rawRules := (d.Get("rule")).([]interface{})
	createRuleOptsList := make([]secgroups.CreateRuleOpts, len(rawRules))
	for i, raw := range rawRules {
		rawMap := raw.(map[string]interface{})
		createRuleOptsList[i] = secgroups.CreateRuleOpts{
			ParentGroupID: d.Id(),
			FromPort:      rawMap["from_port"].(int),
			ToPort:        rawMap["to_port"].(int),
			IPProtocol:    rawMap["ip_protocol"].(string),
			CIDR:          rawMap["cidr"].(string),
			FromGroupID:   rawMap["from_group_id"].(string),
		}
	}
	return createRuleOptsList
}

func resourceSecGroupRuleCreateOptsV2(d *schema.ResourceData, raw interface{}) secgroups.CreateRuleOpts {
	rawMap := raw.(map[string]interface{})
	return secgroups.CreateRuleOpts{
		ParentGroupID: d.Id(),
		FromPort:      rawMap["from_port"].(int),
		ToPort:        rawMap["to_port"].(int),
		IPProtocol:    rawMap["ip_protocol"].(string),
		CIDR:          rawMap["cidr"].(string),
		FromGroupID:   rawMap["from_group_id"].(string),
	}
}

func resourceSecGroupRuleV2(d *schema.ResourceData, raw interface{}) secgroups.Rule {
	rawMap := raw.(map[string]interface{})
	return secgroups.Rule{
		ID:            rawMap["id"].(string),
		ParentGroupID: d.Id(),
		FromPort:      rawMap["from_port"].(int),
		ToPort:        rawMap["to_port"].(int),
		IPProtocol:    rawMap["ip_protocol"].(string),
		IPRange:       secgroups.IPRange{CIDR: rawMap["cidr"].(string)},
	}
}

func rulesToMap(sgrs []secgroups.Rule) []map[string]interface{} {
	sgrMap := make([]map[string]interface{}, len(sgrs))
	for i, sgr := range sgrs {
		sgrMap[i] = map[string]interface{}{
			"id":          sgr.ID,
			"from_port":   sgr.FromPort,
			"to_port":     sgr.ToPort,
			"ip_protocol": sgr.IPProtocol,
			"cidr":        sgr.IPRange.CIDR,
		}
	}
	return sgrMap
}

func secgroupRuleV2Hash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	buf.WriteString(fmt.Sprintf("%d-", m["from_port"].(int)))
	buf.WriteString(fmt.Sprintf("%d-", m["to_port"].(int)))
	buf.WriteString(fmt.Sprintf("%s-", m["ip_protocol"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["cidr"].(string)))

	return hashcode.String(buf.String())
}
