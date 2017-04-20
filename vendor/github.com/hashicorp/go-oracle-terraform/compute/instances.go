package compute

import (
	"errors"
	"fmt"
	"strings"
)

const WaitForInstanceReadyTimeout = 600
const WaitForInstanceDeleteTimeout = 600

// InstancesClient is a client for the Instance functions of the Compute API.
type InstancesClient struct {
	ResourceClient
}

// Instances obtains an InstancesClient which can be used to access to the
// Instance functions of the Compute API
func (c *Client) Instances() *InstancesClient {
	return &InstancesClient{
		ResourceClient: ResourceClient{
			Client:              c,
			ResourceDescription: "instance",
			ContainerPath:       "/launchplan/",
			ResourceRootPath:    "/instance",
		}}
}

type InstanceState string

const (
	InstanceRunning      InstanceState = "running"
	InstanceInitializing InstanceState = "initializing"
	InstancePreparing    InstanceState = "preparing"
	InstanceStopping     InstanceState = "stopping"
	InstanceQueued       InstanceState = "queued"
	InstanceError        InstanceState = "error"
)

// InstanceInfo represents the Compute API's view of the state of an instance.
type InstanceInfo struct {
	// The ID for the instance. Set by the SDK based on the request - not the API.
	ID string

	// A dictionary of attributes to be made available to the instance.
	// A value with the key "userdata" will be made available in an EC2-compatible manner.
	Attributes map[string]interface{} `json:"attributes"`

	// The availability domain for the instance
	AvailabilityDomain string `json:"availability_domain"`

	// Boot order list.
	BootOrder []int `json:"boot_order"`

	// The default domain to use for the hostname and DNS lookups
	Domain string `json:"domain"`

	// Optional ImageListEntry number. Default will be used if not specified
	Entry int `json:"entry"`

	// The reason for the instance going to error state, if available.
	ErrorReason string `json:"error_reason"`

	// SSH Server Fingerprint presented by the instance
	Fingerprint string `json:"fingerprint"`

	// The hostname for the instance
	Hostname string `json:"hostname"`

	// The format of the image
	ImageFormat string `json:"image_format"`

	// Name of imagelist to be launched.
	ImageList string `json:"imagelist"`

	// IP address of the instance.
	IPAddress string `json:"ip"`

	// A label assigned by the user, specifically for defining inter-instance relationships.
	Label string `json:"label"`

	// Name of this instance, generated by the server.
	Name string `json:"name"`

	// Mapping of to network specifiers for virtual NICs to be attached to this instance.
	Networking map[string]NetworkingInfo `json:"networking"`

	// A list of strings specifying arbitrary tags on nodes to be matched on placement.
	PlacementRequirements []string `json:"placement_requirements"`

	// The OS platform for the instance.
	Platform string `json:"platform"`

	// The priority at which this instance will be run
	Priority string `json:"priority"`

	// Reference to the QuotaReservation, to be destroyed with the instance
	QuotaReservation string `json:"quota_reservation"`

	// Array of relationship specifications to be satisfied on this instance's placement
	Relationships []string `json:"relationships"`

	// Resolvers to use instead of the default resolvers
	Resolvers []string `json:"resolvers"`

	// Add PTR records for the hostname
	ReverseDNS bool `json:"reverse_dns"`

	// Type of instance, as defined on site configuration.
	Shape string `json:"shape"`

	// Site to run on
	Site string `json:"site"`

	// ID's of SSH keys that will be exposed to the instance.
	SSHKeys []string `json:"sshkeys"`

	// The start time of the instance
	StartTime string `json:"start_time"`

	// State of the instance.
	State InstanceState `json:"state"`

	// The Storage Attachment information.
	Storage []StorageAttachment `json:"storage_attachments"`

	// Array of tags associated with the instance.
	Tags []string `json:"tags"`

	// vCable for this instance.
	VCableID string `json:"vcable_id"`

	// Specify if the devices created for the instance are virtio devices. If not specified, the default
	// will come from the cluster configuration file
	Virtio bool `json:"virtio,omitempty"`

	// IP Address and port of the VNC console for the instance
	VNC string `json:"vnc"`
}

type StorageAttachment struct {
	// The index number for the volume.
	Index int `json:"index"`

	// The three-part name (/Compute-identity_domain/user/object) of the storage attachment.
	Name string `json:"name"`

	// The three-part name (/Compute-identity_domain/user/object) of the storage volume attached to the instance.
	StorageVolumeName string `json:"storage_volume_name"`
}

func (i *InstanceInfo) getInstanceName() string {
	return fmt.Sprintf(CMP_QUALIFIED_NAME, i.Name, i.ID)
}

type CreateInstanceInput struct {
	// A dictionary of user-defined attributes to be made available to the instance.
	// Optional
	Attributes map[string]interface{} `json:"attributes"`
	// Boot order list
	// Optional
	BootOrder []int `json:"boot_order"`
	// The host name assigned to the instance. On an Oracle Linux instance,
	// this host name is displayed in response to the hostname command.
	// Only relative DNS is supported. The domain name is suffixed to the host name
	// that you specify. The host name must not end with a period. If you don't specify a
	// host name, then a name is generated automatically.
	// Optional
	Hostname string `json:"hostname"`
	// Name of imagelist to be launched.
	// Optional
	ImageList string `json:"imagelist"`
	// A label assigned by the user, specifically for defining inter-instance relationships.
	// Optional
	Label string `json:"label"`
	// Name of this instance, generated by the server.
	// Optional
	Name string `json:"name"`
	// Networking information.
	// Optional
	Networking map[string]NetworkingInfo `json:"networking"`
	// If set to true (default), then reverse DNS records are created.
	// If set to false, no reverse DNS records are created.
	// Optional
	ReverseDNS bool `json:"reverse_dns,omitempty"`
	// Type of instance, as defined on site configuration.
	// Required
	Shape string `json:"shape"`
	// A list of the Storage Attachments you want to associate with the instance.
	// Optional
	Storage []StorageAttachmentInput `json:"storage_attachments"`
	// A list of the SSH public keys that you want to associate with the instance.
	// Optional
	SSHKeys []string `json:"sshkeys"`
	// A list of tags to be supplied to the instance
	// Optional
	Tags []string `json:"tags"`
}

type StorageAttachmentInput struct {
	// The index number for the volume. The allowed range is 1 to 10.
	// If you want to use a storage volume as the boot disk for an instance, you must specify the index number for that volume as 1.
	// The index determines the device name by which the volume is exposed to the instance.
	Index int `json:"index"`
	// The three-part name (/Compute-identity_domain/user/object) of the storage volume that you want to attach to the instance.
	// Note that volumes attached to an instance at launch time can't be detached.
	Volume string `json:"volume"`
}

const ReservationPrefix = "ipreservation"
const ReservationIPPrefix = "network/v1/ipreservation"

type NICModel string

const (
	NICDefaultModel NICModel = "e1000"
)

// Struct of Networking info from a populated instance, or to be used as input to create an instance
type NetworkingInfo struct {
	// The DNS name for the Shared network (Required)
	// DNS A Record for an IP Network (Optional)
	DNS []string `json:"dns,omitempty"`
	// IP Network only.
	// If you want to associate a static private IP Address,
	// specify that here within the range of the supplied IPNetwork attribute.
	// Optional
	IPAddress string `json:"ip,omitempty"`
	// IP Network only.
	// The name of the IP Network you want to add the instance to.
	// Required
	IPNetwork string `json:"ipnetwork,omitempty"`
	// IP Network only.
	// The hexadecimal MAC Address of the interface
	// Optional
	MACAddress string `json:"address,omitempty"`
	// Shared Network only.
	// The type of NIC used. Must be set to 'e1000'
	// Required
	Model NICModel `json:"model,omitempty"`
	// IP Network and Shared Network
	// The name servers that are sent through DHCP as option 6.
	// You can specify a maximum of eight name server IP addresses per interface.
	// Optional
	NameServers []string `json:"name_servers,omitempty"`
	// The names of an IP Reservation to associate in an IP Network (Optional)
	// Indicates whether a temporary or permanent public IP Address should be assigned
	// in a Shared Network (Required)
	Nat []string `json:"nat,omitempty"`
	// IP Network and Shared Network
	// The search domains that should be sent through DHCP as option 119.
	// You can enter a maximum of eight search domain zones per interface.
	// Optional
	SearchDomains []string `json:"search_domains,omitempty"`
	// Shared Network only.
	// The security lists that you want to add the instance to
	// Required
	SecLists []string `json:"seclists,omitempty"`
	// IP Network Only
	// The name of the vNIC
	// Optional
	Vnic string `json:"vnic,omitempty"`
	// IP Network only.
	// The names of the vNICSets you want to add the interface to.
	// Optional
	VnicSets []string `json:"vnicsets,omitempty"`
}

// LaunchPlan defines a launch plan, used to launch instances with the supplied InstanceSpec(s)
type LaunchPlanInput struct {
	// Describes an array of instances which should be launched
	Instances []CreateInstanceInput `json:"instances"`
}

type LaunchPlanResponse struct {
	// An array of instances which have been launched
	Instances []InstanceInfo `json:"instances"`
}

// LaunchInstance creates and submits a LaunchPlan to launch a new instance.
func (c *InstancesClient) CreateInstance(input *CreateInstanceInput) (*InstanceInfo, error) {
	qualifiedSSHKeys := []string{}
	for _, key := range input.SSHKeys {
		qualifiedSSHKeys = append(qualifiedSSHKeys, c.getQualifiedName(key))
	}

	input.SSHKeys = qualifiedSSHKeys

	qualifiedStorageAttachments := []StorageAttachmentInput{}
	for _, attachment := range input.Storage {
		qualifiedStorageAttachments = append(qualifiedStorageAttachments, StorageAttachmentInput{
			Index:  attachment.Index,
			Volume: c.getQualifiedName(attachment.Volume),
		})
	}
	input.Storage = qualifiedStorageAttachments

	input.Networking = c.qualifyNetworking(input.Networking)

	input.Name = fmt.Sprintf(CMP_QUALIFIED_NAME, c.getUserName(), input.Name)

	plan := LaunchPlanInput{Instances: []CreateInstanceInput{*input}}

	var responseBody LaunchPlanResponse
	if err := c.createResource(&plan, &responseBody); err != nil {
		return nil, err
	}

	if len(responseBody.Instances) == 0 {
		return nil, fmt.Errorf("No instance information returned: %#v", responseBody)
	}

	// Call wait for instance ready now, as creating the instance is an eventually consistent operation
	getInput := &GetInstanceInput{
		Name: input.Name,
		ID:   responseBody.Instances[0].ID,
	}

	// Wait for instance to be ready and return the result
	// Don't have to unqualify any objects, as the GetInstance method will handle that
	return c.WaitForInstanceRunning(getInput, WaitForInstanceReadyTimeout)
}

// Both of these fields are required. If they're not provided, things go wrong in
// incredibly amazing ways.
type GetInstanceInput struct {
	// The Unqualified Name of this Instance
	Name string
	// The Unqualified ID of this Instance
	ID string
}

func (g *GetInstanceInput) String() string {
	return fmt.Sprintf(CMP_QUALIFIED_NAME, g.Name, g.ID)
}

// GetInstance retrieves information about an instance.
func (c *InstancesClient) GetInstance(input *GetInstanceInput) (*InstanceInfo, error) {
	if input.ID == "" || input.Name == "" {
		return nil, errors.New("Both instance name and ID need to be specified")
	}

	var responseBody InstanceInfo
	if err := c.getResource(input.String(), &responseBody); err != nil {
		return nil, err
	}

	if responseBody.Name == "" {
		return nil, fmt.Errorf("Empty response body when requesting instance %s", input.Name)
	}

	// The returned 'Name' attribute is the fully qualified instance name + "/" + ID
	// Split these out to accurately populate the fields
	nID := strings.Split(c.getUnqualifiedName(responseBody.Name), "/")
	responseBody.Name = nID[0]
	responseBody.ID = nID[1]

	c.unqualify(&responseBody.VCableID)

	// Unqualify SSH Key names
	sshKeyNames := []string{}
	for _, sshKeyRef := range responseBody.SSHKeys {
		sshKeyNames = append(sshKeyNames, c.getUnqualifiedName(sshKeyRef))
	}
	responseBody.SSHKeys = sshKeyNames

	var networkingErr error
	responseBody.Networking, networkingErr = c.unqualifyNetworking(responseBody.Networking)
	if networkingErr != nil {
		return nil, networkingErr
	}
	responseBody.Storage = c.unqualifyStorage(responseBody.Storage)

	return &responseBody, nil
}

type DeleteInstanceInput struct {
	// The Unqualified Name of this Instance
	Name string
	// The Unqualified ID of this Instance
	ID string
}

func (d *DeleteInstanceInput) String() string {
	return fmt.Sprintf(CMP_QUALIFIED_NAME, d.Name, d.ID)
}

// DeleteInstance deletes an instance.
func (c *InstancesClient) DeleteInstance(input *DeleteInstanceInput) error {
	// Call to delete the instance
	if err := c.deleteResource(input.String()); err != nil {
		return err
	}
	// Wait for instance to be deleted
	return c.WaitForInstanceDeleted(input, WaitForInstanceDeleteTimeout)
}

// WaitForInstanceRunning waits for an instance to be completely initialized and available.
func (c *InstancesClient) WaitForInstanceRunning(input *GetInstanceInput, timeoutSeconds int) (*InstanceInfo, error) {
	var info *InstanceInfo
	var getErr error
	err := c.waitFor("instance to be ready", timeoutSeconds, func() (bool, error) {
		info, getErr = c.GetInstance(input)
		if getErr != nil {
			return false, getErr
		}
		switch s := info.State; s {
		case InstanceError:
			return false, fmt.Errorf("Error initializing instance: %s", info.ErrorReason)
		case InstanceRunning:
			c.debugLogString("Instance Running")
			return true, nil
		case InstanceQueued:
			c.debugLogString("Instance Queuing")
			return false, nil
		case InstanceInitializing:
			c.debugLogString("Instance Initializing")
			return false, nil
		case InstancePreparing:
			c.debugLogString("Instance Preparing")
			return false, nil
		default:
			c.debugLogString(fmt.Sprintf("Unknown instance state: %s, waiting", s))
			return false, nil
		}
	})
	return info, err
}

// WaitForInstanceDeleted waits for an instance to be fully deleted.
func (c *InstancesClient) WaitForInstanceDeleted(input *DeleteInstanceInput, timeoutSeconds int) error {
	return c.waitFor("instance to be deleted", timeoutSeconds, func() (bool, error) {
		var info InstanceInfo
		if err := c.getResource(input.String(), &info); err != nil {
			if WasNotFoundError(err) {
				// Instance could not be found, thus deleted
				return true, nil
			}
			// Some other error occurred trying to get instance, exit
			return false, err
		}
		switch s := info.State; s {
		case InstanceError:
			return false, fmt.Errorf("Error stopping instance: %s", info.ErrorReason)
		case InstanceStopping:
			c.debugLogString("Instance stopping")
			return false, nil
		default:
			c.debugLogString(fmt.Sprintf("Unknown instance state: %s, waiting", s))
			return false, nil
		}
	})
}

func (c *InstancesClient) qualifyNetworking(info map[string]NetworkingInfo) map[string]NetworkingInfo {
	qualifiedNetworks := map[string]NetworkingInfo{}
	for k, v := range info {
		qfd := v
		sharedNetwork := false
		if v.IPNetwork != "" {
			// Network interface is for an IP Network
			qfd.IPNetwork = c.getQualifiedName(v.IPNetwork)
			sharedNetwork = true
		}
		if v.Vnic != "" {
			qfd.Vnic = c.getQualifiedName(v.Vnic)
		}
		if v.Nat != nil {
			qfd.Nat = c.qualifyNat(v.Nat, sharedNetwork)
		}
		if v.VnicSets != nil {
			qfd.VnicSets = c.getQualifiedList(v.VnicSets)
		}
		if v.SecLists != nil {
			// Network interface is for the shared network
			secLists := []string{}
			for _, v := range v.SecLists {
				secLists = append(secLists, c.getQualifiedName(v))
			}
			qfd.SecLists = secLists
		}

		qualifiedNetworks[k] = qfd
	}
	return qualifiedNetworks
}

func (c *InstancesClient) unqualifyNetworking(info map[string]NetworkingInfo) (map[string]NetworkingInfo, error) {
	// Unqualify ip network
	var err error
	unqualifiedNetworks := map[string]NetworkingInfo{}
	for k, v := range info {
		unq := v
		if v.IPNetwork != "" {
			unq.IPNetwork = c.getUnqualifiedName(v.IPNetwork)
		}
		if v.Vnic != "" {
			unq.Vnic = c.getUnqualifiedName(v.Vnic)
		}
		if v.Nat != nil {
			unq.Nat, err = c.unqualifyNat(v.Nat)
			if err != nil {
				return nil, err
			}
		}
		if v.VnicSets != nil {
			unq.VnicSets = c.getUnqualifiedList(v.VnicSets)
		}
		if v.SecLists != nil {
			secLists := []string{}
			for _, v := range v.SecLists {
				secLists = append(secLists, c.getUnqualifiedName(v))
			}
			v.SecLists = secLists
		}
		unqualifiedNetworks[k] = unq
	}
	return unqualifiedNetworks, nil
}

func (c *InstancesClient) qualifyNat(nat []string, shared bool) []string {
	qualifiedNats := []string{}
	for _, v := range nat {
		if strings.HasPrefix(v, "ippool:/oracle") {
			qualifiedNats = append(qualifiedNats, v)
			continue
		}
		prefix := ReservationPrefix
		if shared {
			prefix = ReservationIPPrefix
		}
		qualifiedNats = append(qualifiedNats, fmt.Sprintf("%s:%s", prefix, c.getQualifiedName(v)))
	}
	return qualifiedNats
}

func (c *InstancesClient) unqualifyNat(nat []string) ([]string, error) {
	unQualifiedNats := []string{}
	for _, v := range nat {
		if strings.HasPrefix(v, "ippool:/oracle") {
			unQualifiedNats = append(unQualifiedNats, v)
			continue
		}
		n := strings.Split(v, ":")
		if len(n) < 1 {
			return nil, fmt.Errorf("Error unqualifying NAT: %s", v)
		}
		u := n[1]
		unQualifiedNats = append(unQualifiedNats, c.getUnqualifiedName(u))
	}
	return unQualifiedNats, nil
}

func (c *InstancesClient) unqualifyStorage(attachments []StorageAttachment) []StorageAttachment {
	unqAttachments := []StorageAttachment{}
	for _, v := range attachments {
		if v.StorageVolumeName != "" {
			v.StorageVolumeName = c.getUnqualifiedName(v.StorageVolumeName)
		}
		unqAttachments = append(unqAttachments, v)
	}

	return unqAttachments
}
