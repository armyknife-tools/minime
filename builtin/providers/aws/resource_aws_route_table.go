package aws

import (
	"bytes"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/aws-sdk-go/aws"
	"github.com/hashicorp/aws-sdk-go/gen/ec2"
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAwsRouteTable() *schema.Resource {
	return &schema.Resource{
		Create: resourceAwsRouteTableCreate,
		Read:   resourceAwsRouteTableRead,
		Update: resourceAwsRouteTableUpdate,
		Delete: resourceAwsRouteTableDelete,

		Schema: map[string]*schema.Schema{
			"vpc_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"tags": tagsSchema(),

			"route": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cidr_block": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},

						"gateway_id": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},

						"instance_id": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},

						"vpc_peering_connection_id": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
				Set: resourceAwsRouteTableHash,
			},
		},
	}
}

func resourceAwsRouteTableCreate(d *schema.ResourceData, meta interface{}) error {
	ec2conn := meta.(*AWSClient).awsEC2conn

	// Create the routing table
	createOpts := &ec2.CreateRouteTableRequest{
		VPCID: aws.String(d.Get("vpc_id").(string)),
	}
	log.Printf("[DEBUG] RouteTable create config: %#v", createOpts)

	resp, err := ec2conn.CreateRouteTable(createOpts)
	if err != nil {
		return fmt.Errorf("Error creating route table: %s", err)
	}

	// Get the ID and store it
	rt := resp.RouteTable
	d.SetId(*rt.RouteTableID)
	log.Printf("[INFO] Route Table ID: %s", d.Id())

	// Wait for the route table to become available
	log.Printf(
		"[DEBUG] Waiting for route table (%s) to become available",
		d.Id())
	stateConf := &resource.StateChangeConf{
		Pending: []string{"pending"},
		Target:  "ready",
		Refresh: resourceAwsRouteTableStateRefreshFuncSDK(ec2conn, d.Id()),
		Timeout: 1 * time.Minute,
	}
	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf(
			"Error waiting for route table (%s) to become available: %s",
			d.Id(), err)
	}

	return resourceAwsRouteTableUpdate(d, meta)
}

func resourceAwsRouteTableRead(d *schema.ResourceData, meta interface{}) error {
	ec2conn := meta.(*AWSClient).awsEC2conn

	rtRaw, _, err := resourceAwsRouteTableStateRefreshFuncSDK(ec2conn, d.Id())()
	if err != nil {
		return err
	}
	if rtRaw == nil {
		return nil
	}

	rt := rtRaw.(*ec2.RouteTable)
	d.Set("vpc_id", rt.VPCID)

	// Create an empty schema.Set to hold all routes
	route := &schema.Set{F: resourceAwsRouteTableHash}

	// Loop through the routes and add them to the set
	for _, r := range rt.Routes {
		if r.GatewayID != nil && *r.GatewayID == "local" {
			continue
		}

		if r.Origin != nil && *r.Origin == "EnableVgwRoutePropagation" {
			continue
		}

		m := make(map[string]interface{})

		if r.DestinationCIDRBlock != nil {
			m["cidr_block"] = *r.DestinationCIDRBlock
		}
		if r.GatewayID != nil {
			m["gateway_id"] = *r.GatewayID
		}
		if r.InstanceID != nil {
			m["instance_id"] = *r.InstanceID
		}
		if r.VPCPeeringConnectionID != nil {
			m["vpc_peering_connection_id"] = *r.VPCPeeringConnectionID
		}

		route.Add(m)
	}
	d.Set("route", route)

	// Tags
	d.Set("tags", tagsToMapSDK(rt.Tags))

	return nil
}

func resourceAwsRouteTableUpdate(d *schema.ResourceData, meta interface{}) error {
	ec2conn := meta.(*AWSClient).awsEC2conn

	// Check if the route set as a whole has changed
	if d.HasChange("route") {
		o, n := d.GetChange("route")
		ors := o.(*schema.Set).Difference(n.(*schema.Set))
		nrs := n.(*schema.Set).Difference(o.(*schema.Set))

		// Now first loop through all the old routes and delete any obsolete ones
		for _, route := range ors.List() {
			m := route.(map[string]interface{})

			// Delete the route as it no longer exists in the config
			log.Printf(
				"[INFO] Deleting route from %s: %s",
				d.Id(), m["cidr_block"].(string))
			err := ec2conn.DeleteRoute(&ec2.DeleteRouteRequest{
				RouteTableID:         aws.String(d.Id()),
				DestinationCIDRBlock: aws.String(m["cidr_block"].(string)),
			})
			if err != nil {
				return err
			}
		}

		// Make sure we save the state of the currently configured rules
		routes := o.(*schema.Set).Intersection(n.(*schema.Set))
		d.Set("route", routes)

		// Then loop through al the newly configured routes and create them
		for _, route := range nrs.List() {
			m := route.(map[string]interface{})

			opts := ec2.CreateRouteRequest{
				RouteTableID:           aws.String(d.Id()),
				DestinationCIDRBlock:   aws.String(m["cidr_block"].(string)),
				GatewayID:              aws.String(m["gateway_id"].(string)),
				InstanceID:             aws.String(m["instance_id"].(string)),
				VPCPeeringConnectionID: aws.String(m["vpc_peering_connection_id"].(string)),
			}

			log.Printf("[INFO] Creating route for %s: %#v", d.Id(), opts)
			if err := ec2conn.CreateRoute(&opts); err != nil {
				return err
			}

			routes.Add(route)
			d.Set("route", routes)
		}
	}

	if err := setTagsSDK(ec2conn, d); err != nil {
		return err
	} else {
		d.SetPartial("tags")
	}

	return resourceAwsRouteTableRead(d, meta)
}

func resourceAwsRouteTableDelete(d *schema.ResourceData, meta interface{}) error {
	ec2conn := meta.(*AWSClient).awsEC2conn

	// First request the routing table since we'll have to disassociate
	// all the subnets first.
	rtRaw, _, err := resourceAwsRouteTableStateRefreshFuncSDK(ec2conn, d.Id())()
	if err != nil {
		return err
	}
	if rtRaw == nil {
		return nil
	}
	rt := rtRaw.(*ec2.RouteTable)

	// Do all the disassociations
	for _, a := range rt.Associations {
		log.Printf("[INFO] Disassociating association: %s", *a.RouteTableAssociationID)
		err := ec2conn.DisassociateRouteTable(&ec2.DisassociateRouteTableRequest{
			AssociationID: a.RouteTableAssociationID,
		})
		if err != nil {
			return err
		}
	}

	// Delete the route table
	log.Printf("[INFO] Deleting Route Table: %s", d.Id())
	err = ec2conn.DeleteRouteTable(&ec2.DeleteRouteTableRequest{
		RouteTableID: aws.String(d.Id()),
	})
	if err != nil {
		ec2err, ok := err.(aws.APIError)
		if ok && ec2err.Code == "InvalidRouteTableID.NotFound" {
			return nil
		}

		return fmt.Errorf("Error deleting route table: %s", err)
	}

	// Wait for the route table to really destroy
	log.Printf(
		"[DEBUG] Waiting for route table (%s) to become destroyed",
		d.Id())

	stateConf := &resource.StateChangeConf{
		Pending: []string{"ready"},
		Target:  "",
		Refresh: resourceAwsRouteTableStateRefreshFuncSDK(ec2conn, d.Id()),
		Timeout: 1 * time.Minute,
	}
	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf(
			"Error waiting for route table (%s) to become destroyed: %s",
			d.Id(), err)
	}

	return nil
}

func resourceAwsRouteTableHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	buf.WriteString(fmt.Sprintf("%s-", m["cidr_block"].(string)))

	if v, ok := m["gateway_id"]; ok {
		buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	}

	if v, ok := m["instance_id"]; ok {
		buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	}

	if v, ok := m["vpc_peering_connection_id"]; ok {
		buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	}

	return hashcode.String(buf.String())
}

// resourceAwsRouteTableStateRefreshFuncSDK returns a resource.StateRefreshFunc that is used to watch
// a RouteTable.
func resourceAwsRouteTableStateRefreshFuncSDK(conn *ec2.EC2, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := conn.DescribeRouteTables(&ec2.DescribeRouteTablesRequest{
			RouteTableIDs: []string{id},
		})
		if err != nil {
			if ec2err, ok := err.(aws.APIError); ok && ec2err.Code == "InvalidRouteTableID.NotFound" {
				resp = nil
			} else {
				log.Printf("Error on RouteTableStateRefresh: %s", err)
				return nil, "", err
			}
		}

		if resp == nil {
			// Sometimes AWS just has consistency issues and doesn't see
			// our instance yet. Return an empty state.
			return nil, "", nil
		}

		rt := &resp.RouteTables[0]
		return rt, "ready", nil
	}
}
