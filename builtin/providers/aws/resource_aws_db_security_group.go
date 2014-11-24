package aws

import (
	"bytes"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/multierror"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/mitchellh/goamz/rds"
)

func resourceAwsDbSecurityGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAwsDbSecurityGroupCreate,
		Read:   resourceAwsDbSecurityGroupRead,
		Delete: resourceAwsDbSecurityGroupDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"description": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"ingress": &schema.Schema{
				Type:     schema.TypeSet,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cidr": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},

						"security_group_name": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},

						"security_group_id": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},

						"security_group_owner_id": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
				Set: resourceAwsDbSecurityGroupIngressHash,
			},
		},
	}
}

func resourceAwsDbSecurityGroupCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).rdsconn

	var err error
	var errs []error

	opts := rds.CreateDBSecurityGroup{
		DBSecurityGroupName:        d.Get("name").(string),
		DBSecurityGroupDescription: d.Get("description").(string),
	}

	log.Printf("[DEBUG] DB Security Group create configuration: %#v", opts)
	_, err = conn.CreateDBSecurityGroup(&opts)
	if err != nil {
		return fmt.Errorf("Error creating DB Security Group: %s", err)
	}

	d.SetId(d.Get("name").(string))

	log.Printf("[INFO] DB Security Group ID: %s", d.Id())

	sg, err := resourceAwsDbSecurityGroupRetrieve(d, meta)
	if err != nil {
		return err
	}

	ingresses := d.Get("ingress").(*schema.Set)
	for _, ing := range ingresses.List() {
		err = resourceAwsDbSecurityGroupAuthorizeRule(ing, sg.Name, conn)

		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return &multierror.Error{Errors: errs}
	}

	log.Println(
		"[INFO] Waiting for Ingress Authorizations to be authorized")

	stateConf := &resource.StateChangeConf{
		Pending: []string{"authorizing"},
		Target:  "authorized",
		Refresh: resourceAwsDbSecurityGroupStateRefreshFunc(d, meta),
		Timeout: 10 * time.Minute,
	}

	// Wait, catching any errors
	_, err = stateConf.WaitForState()
	if err != nil {
		return err
	}

	return resourceAwsDbSecurityGroupRead(d, meta)
}

func resourceAwsDbSecurityGroupRead(d *schema.ResourceData, meta interface{}) error {
	sg, err := resourceAwsDbSecurityGroupRetrieve(d, meta)
	if err != nil {
		return err
	}

	d.Set("name", sg.Name)
	d.Set("description", sg.Description)

	// Create an empty schema.Set to hold all ingress rules
	rules := &schema.Set{
		F: resourceAwsDbSecurityGroupIngressHash,
	}

	for _, v := range sg.CidrIps {
		rule := map[string]interface{}{"cidr": v}
		rules.Add(rule)
	}

	for i, _ := range sg.EC2SecurityGroupOwnerIds {
		rule := map[string]interface{}{
			"security_group_name":     sg.EC2SecurityGroupNames[i],
			"security_group_id":       sg.EC2SecurityGroupIds[i],
			"security_group_owner_id": sg.EC2SecurityGroupOwnerIds[i],
		}
		rules.Add(rule)
	}

	d.Set("ingress", rules)

	return nil
}

func resourceAwsDbSecurityGroupDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).rdsconn

	log.Printf("[DEBUG] DB Security Group destroy: %v", d.Id())

	opts := rds.DeleteDBSecurityGroup{DBSecurityGroupName: d.Id()}

	log.Printf("[DEBUG] DB Security Group destroy configuration: %v", opts)
	_, err := conn.DeleteDBSecurityGroup(&opts)

	if err != nil {
		newerr, ok := err.(*rds.Error)
		if ok && newerr.Code == "InvalidDBSecurityGroup.NotFound" {
			return nil
		}
		return err
	}

	return nil
}

func resourceAwsDbSecurityGroupRetrieve(d *schema.ResourceData, meta interface{}) (*rds.DBSecurityGroup, error) {
	conn := meta.(*AWSClient).rdsconn

	opts := rds.DescribeDBSecurityGroups{
		DBSecurityGroupName: d.Id(),
	}

	log.Printf("[DEBUG] DB Security Group describe configuration: %#v", opts)

	resp, err := conn.DescribeDBSecurityGroups(&opts)

	if err != nil {
		return nil, fmt.Errorf("Error retrieving DB Security Groups: %s", err)
	}

	if len(resp.DBSecurityGroups) != 1 ||
		resp.DBSecurityGroups[0].Name != d.Id() {
		if err != nil {
			return nil, fmt.Errorf("Unable to find DB Security Group: %#v", resp.DBSecurityGroups)
		}
	}

	v := resp.DBSecurityGroups[0]

	return &v, nil
}

// Authorizes the ingress rule on the db security group
func resourceAwsDbSecurityGroupAuthorizeRule(ingress interface{}, dbSecurityGroupName string, conn *rds.Rds) error {
	ing := ingress.(map[string]interface{})

	opts := rds.AuthorizeDBSecurityGroupIngress{
		DBSecurityGroupName: dbSecurityGroupName,
	}

	if attr, ok := ing["cidr"]; ok && attr != "" {
		opts.Cidr = attr.(string)
	}

	if attr, ok := ing["security_group_name"]; ok && attr != "" {
		opts.EC2SecurityGroupName = attr.(string)
	}

	if attr, ok := ing["security_group_id"]; ok && attr != "" {
		opts.EC2SecurityGroupId = attr.(string)
	}

	if attr, ok := ing["security_group_owner_id"]; ok && attr != "" {
		opts.EC2SecurityGroupOwnerId = attr.(string)
	}

	log.Printf("[DEBUG] Authorize ingress rule configuration: %#v", opts)

	_, err := conn.AuthorizeDBSecurityGroupIngress(&opts)

	if err != nil {
		return fmt.Errorf("Error authorizing security group ingress: %s", err)
	}

	return nil
}

func resourceAwsDbSecurityGroupIngressHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})

	if v, ok := m["cidr"]; ok {
		buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	}

	if v, ok := m["security_group_name"]; ok {
		buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	}

	if v, ok := m["security_group_id"]; ok {
		buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	}

	if v, ok := m["security_group_owner_id"]; ok {
		buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	}

	return hashcode.String(buf.String())
}

func resourceAwsDbSecurityGroupStateRefreshFunc(
	d *schema.ResourceData, meta interface{}) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		v, err := resourceAwsDbSecurityGroupRetrieve(d, meta)

		if err != nil {
			log.Printf("Error on retrieving DB Security Group when waiting: %s", err)
			return nil, "", err
		}

		statuses := append(v.EC2SecurityGroupStatuses, v.CidrStatuses...)

		for _, stat := range statuses {
			// Not done
			if stat != "authorized" {
				return nil, "authorizing", nil
			}
		}

		return v, "authorized", nil
	}
}
