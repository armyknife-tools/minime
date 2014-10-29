package openstack

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack/compute/v2/servers"
)

func resourceComputeInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceComputeInstanceCreate,
		Read:   resourceComputeInstanceRead,
		Update: resourceComputeInstanceUpdate,
		Delete: resourceComputeInstanceDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},

			"image_ref": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},

			"flavor_ref": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},

			"security_groups": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: false,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set: func(v interface{}) int {
					return hashcode.String(v.(string))
				},
			},

			"availability_zone": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"networks": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"uuid": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},

						"port": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},

						"fixed_ip": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},

			"metadata": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
			},

			"config_drive": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},

			"access_ip_v4": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
				ForceNew: false,
			},

			"access_ip_v6": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
				ForceNew: false,
			},
		},
	}
}

func resourceComputeInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	osClient := config.computeV2Client

	createOpts := &servers.CreateOpts{
		Name:      d.Get("name").(string),
		ImageRef:  d.Get("image_ref").(string),
		FlavorRef: d.Get("flavor_ref").(string),
		//SecurityGroups []string
		AvailabilityZone: d.Get("availability_zone").(string),
		Networks:         resourceInstanceNetworks(d),
		Metadata:         resourceInstanceMetadata(d),
		ConfigDrive:      d.Get("config_drive").(bool),
	}

	log.Printf("[INFO] Requesting instance creation")
	server, err := servers.Create(osClient, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error creating OpenStack server: %s", err)
	}
	log.Printf("[INFO] Instance ID: %s", server.ID)

	// Store the ID now
	d.SetId(server.ID)

	// Wait for the instance to become running so we can get some attributes
	// that aren't available until later.
	log.Printf(
		"[DEBUG] Waiting for instance (%s) to become running",
		server.ID)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"BUILD"},
		Target:     "ACTIVE",
		Refresh:    ServerStateRefreshFunc(osClient, server.ID),
		Timeout:    10 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	serverRaw, err := stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(
			"Error waiting for instance (%s) to become ready: %s",
			server.ID, err)
	}

	server = serverRaw.(*servers.Server)

	return resourceComputeInstanceRead(d, meta)
}

func resourceComputeInstanceRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	osClient := config.computeV2Client

	server, err := servers.Get(osClient, d.Id()).Extract()
	if err != nil {
		return fmt.Errorf("Error retrieving OpenStack server: %s", err)
	}

	log.Printf("[DEBUG] Retreived Server %s: %+v", d.Id(), server)

	d.Set("name", server.Name)
	d.Set("access_ip_v4", server.AccessIPv4)
	d.Set("access_ip_v6", server.AccessIPv6)

	host := server.AccessIPv4
	if host == "" {
		if publicAddressesRaw, ok := server.Addresses["public"]; ok {
			publicAddresses := publicAddressesRaw.([]interface{})
			for _, paRaw := range publicAddresses {
				pa := paRaw.(map[string]interface{})
				if pa["version"].(float64) == 4 {
					host = pa["addr"].(string)
				}
			}
		}
	}

	log.Printf("host: %s", host)

	// Initialize the connection info
	d.SetConnInfo(map[string]string{
		"type": "ssh",
		"host": host,
	})

	d.Set("metadata", server.Metadata)

	return nil
}

func resourceComputeInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	osClient := config.computeV2Client

	var updateOpts servers.UpdateOpts
	// If the Metadata has changed, then update that.
	if d.HasChange("name") {
		updateOpts.Name = d.Get("name").(string)
	}
	if d.HasChange("access_ip_v4") {
		updateOpts.AccessIPv4 = d.Get("access_ip_v4").(string)
	}
	if d.HasChange("access_ip_v6") {
		updateOpts.AccessIPv4 = d.Get("access_ip_v6").(string)
	}

	// If there's nothing to update, don't waste an HTTP call.
	if updateOpts != (servers.UpdateOpts{}) {
		log.Printf("[DEBUG] Updating Server %s with options: %+v", d.Id(), updateOpts)

		_, err := servers.Update(osClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return fmt.Errorf("Error updating Openstack server: %s", err)
		}
	}

	return resourceComputeInstanceRead(d, meta)
}

func resourceComputeInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	osClient := config.computeV2Client

	err := servers.Delete(osClient, d.Id()).ExtractErr()
	if err != nil {
		return fmt.Errorf("Error deleting OpenStack server: %s", err)
	}

	// Wait for the instance to delete before moving on.
	log.Printf(
		"[DEBUG] Waiting for instance (%s) to delete",
		d.Id())

	stateConf := &resource.StateChangeConf{
		Target:     "",
		Refresh:    ServerStateRefreshFunc(osClient, d.Id()),
		Timeout:    10 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(
			"Error waiting for instance (%s) to delete: %s",
			d.Id(), err)
	}

	d.SetId("")
	return nil
}

// ServerStateRefreshFunc returns a resource.StateRefreshFunc that is used to watch
// an OpenStack instance.
func ServerStateRefreshFunc(client *gophercloud.ServiceClient, instanceID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		s, err := servers.Get(client, instanceID).Extract()
		if err != nil {
			return nil, "", err
		}

		return s, s.Status, nil
	}
}

func resourceInstanceNetworks(d *schema.ResourceData) []servers.Network {
	rawNetworks := d.Get("networks").([]interface{})
	networks := make([]servers.Network, len(rawNetworks))
	for i, raw := range rawNetworks {
		rawMap := raw.(map[string]interface{})
		networks[i] = servers.Network{
			UUID:    rawMap["uuid"].(string),
			Port:    rawMap["port"].(string),
			FixedIP: rawMap["fixed_ip"].(string),
		}
	}
	return networks
}

func resourceInstanceMetadata(d *schema.ResourceData) map[string]string {
	m := make(map[string]string)
	for key, val := range d.Get("metadata").(map[string]interface{}) {
		m[key] = val.(string)
	}
	return m
}
