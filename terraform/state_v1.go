package terraform

import (
	"fmt"

	"github.com/mitchellh/copystructure"
)

// stateV1 keeps track of a snapshot state-of-the-world that Terraform
// can use to keep track of what real world resources it is actually
// managing.
//
// stateV1 is _only used for the purposes of backwards compatibility
// and is no longer used in Terraform.
type stateV1 struct {
	// Version is the protocol version. "1" for a StateV1.
	Version int `json:"version"`

	// Serial is incremented on any operation that modifies
	// the State file. It is used to detect potentially conflicting
	// updates.
	Serial int64 `json:"serial"`

	// Remote is used to track the metadata required to
	// pull and push state files from a remote storage endpoint.
	Remote *remoteStateV1 `json:"remote,omitempty"`

	// Modules contains all the modules in a breadth-first order
	Modules []*moduleStateV1 `json:"modules"`
}

// upgrade is used to upgrade a V1 state representation
// into a State (current) representation.
func (old *stateV1) upgrade() (*State, error) {
	if old == nil {
		return nil, nil
	}

	remote, err := old.Remote.upgrade()
	if err != nil {
		return nil, fmt.Errorf("Error upgrading State V1: %v", err)
	}

	modules := make([]*ModuleState, len(old.Modules))
	for i, module := range old.Modules {
		upgraded, err := module.upgrade()
		if err != nil {
			return nil, fmt.Errorf("Error upgrading State V1: %v", err)
		}
		modules[i] = upgraded
	}
	if len(modules) == 0 {
		modules = nil
	}

	newState := &State{
		Version: old.Version,
		Serial:  old.Serial,
		Remote:  remote,
		Modules: modules,
	}

	newState.sort()

	return newState, nil
}

type remoteStateV1 struct {
	// Type controls the client we use for the remote state
	Type string `json:"type"`

	// Config is used to store arbitrary configuration that
	// is type specific
	Config map[string]string `json:"config"`
}

func (old *remoteStateV1) upgrade() (*RemoteState, error) {
	if old == nil {
		return nil, nil
	}

	config, err := copystructure.Copy(old.Config)
	if err != nil {
		return nil, fmt.Errorf("Error upgrading RemoteState V1: %v", err)
	}

	return &RemoteState{
		Type:   old.Type,
		Config: config.(map[string]string),
	}, nil
}

type moduleStateV1 struct {
	// Path is the import path from the root module. Modules imports are
	// always disjoint, so the path represents amodule tree
	Path []string `json:"path"`

	// Outputs declared by the module and maintained for each module
	// even though only the root module technically needs to be kept.
	// This allows operators to inspect values at the boundaries.
	Outputs map[string]string `json:"outputs"`

	// Resources is a mapping of the logically named resource to
	// the state of the resource. Each resource may actually have
	// N instances underneath, although a user only needs to think
	// about the 1:1 case.
	Resources map[string]*resourceStateV1 `json:"resources"`

	// Dependencies are a list of things that this module relies on
	// existing to remain intact. For example: an module may depend
	// on a VPC ID given by an aws_vpc resource.
	//
	// Terraform uses this information to build valid destruction
	// orders and to warn the user if they're destroying a module that
	// another resource depends on.
	//
	// Things can be put into this list that may not be managed by
	// Terraform. If Terraform doesn't find a matching ID in the
	// overall state, then it assumes it isn't managed and doesn't
	// worry about it.
	Dependencies []string `json:"depends_on,omitempty"`
}

func (old *moduleStateV1) upgrade() (*ModuleState, error) {
	if old == nil {
		return nil, nil
	}

	path, err := copystructure.Copy(old.Path)
	if err != nil {
		return nil, fmt.Errorf("Error upgrading ModuleState V1: %v", err)
	}

	// Outputs needs upgrading to use the new structure
	outputs := make(map[string]*OutputState)
	for key, output := range old.Outputs {
		outputs[key] = &OutputState{
			Type:      "string",
			Value:     output,
			Sensitive: false,
		}
	}
	if len(outputs) == 0 {
		outputs = nil
	}

	resources := make(map[string]*ResourceState)
	for key, oldResource := range old.Resources {
		upgraded, err := oldResource.upgrade()
		if err != nil {
			return nil, fmt.Errorf("Error upgrading ModuleState V1: %v", err)
		}
		resources[key] = upgraded
	}
	if len(resources) == 0 {
		resources = nil
	}

	dependencies, err := copystructure.Copy(old.Dependencies)
	if err != nil {
		return nil, fmt.Errorf("Error upgrading ModuleState V1: %v", err)
	}

	return &ModuleState{
		Path:         path.([]string),
		Outputs:      outputs,
		Resources:    resources,
		Dependencies: dependencies.([]string),
	}, nil
}

type resourceStateV1 struct {
	// This is filled in and managed by Terraform, and is the resource
	// type itself such as "mycloud_instance". If a resource provider sets
	// this value, it won't be persisted.
	Type string `json:"type"`

	// Dependencies are a list of things that this resource relies on
	// existing to remain intact. For example: an AWS instance might
	// depend on a subnet (which itself might depend on a VPC, and so
	// on).
	//
	// Terraform uses this information to build valid destruction
	// orders and to warn the user if they're destroying a resource that
	// another resource depends on.
	//
	// Things can be put into this list that may not be managed by
	// Terraform. If Terraform doesn't find a matching ID in the
	// overall state, then it assumes it isn't managed and doesn't
	// worry about it.
	Dependencies []string `json:"depends_on,omitempty"`

	// Primary is the current active instance for this resource.
	// It can be replaced but only after a successful creation.
	// This is the instances on which providers will act.
	Primary *instanceStateV1 `json:"primary"`

	// Tainted is used to track any underlying instances that
	// have been created but are in a bad or unknown state and
	// need to be cleaned up subsequently.  In the
	// standard case, there is only at most a single instance.
	// However, in pathological cases, it is possible for the number
	// of instances to accumulate.
	Tainted []*instanceStateV1 `json:"tainted,omitempty"`

	// Deposed is used in the mechanics of CreateBeforeDestroy: the existing
	// Primary is Deposed to get it out of the way for the replacement Primary to
	// be created by Apply. If the replacement Primary creates successfully, the
	// Deposed instance is cleaned up. If there were problems creating the
	// replacement, the instance remains in the Deposed list so it can be
	// destroyed in a future run. Functionally, Deposed instances are very
	// similar to Tainted instances in that Terraform is only tracking them in
	// order to remember to destroy them.
	Deposed []*instanceStateV1 `json:"deposed,omitempty"`

	// Provider is used when a resource is connected to a provider with an alias.
	// If this string is empty, the resource is connected to the default provider,
	// e.g. "aws_instance" goes with the "aws" provider.
	// If the resource block contained a "provider" key, that value will be set here.
	Provider string `json:"provider,omitempty"`
}

func (old *resourceStateV1) upgrade() (*ResourceState, error) {
	if old == nil {
		return nil, nil
	}

	dependencies, err := copystructure.Copy(old.Dependencies)
	if err != nil {
		return nil, fmt.Errorf("Error upgrading ResourceState V1: %v", err)
	}

	primary, err := old.Primary.upgrade()
	if err != nil {
		return nil, fmt.Errorf("Error upgrading ResourceState V1: %v", err)
	}

	tainted := make([]*InstanceState, len(old.Tainted))
	for i, v := range old.Tainted {
		upgraded, err := v.upgrade()
		if err != nil {
			return nil, fmt.Errorf("Error upgrading ResourceState V1: %v", err)
		}
		tainted[i] = upgraded
	}
	if len(tainted) == 0 {
		tainted = nil
	}

	deposed := make([]*InstanceState, len(old.Deposed))
	for i, v := range old.Deposed {
		upgraded, err := v.upgrade()
		if err != nil {
			return nil, fmt.Errorf("Error upgrading ResourceState V1: %v", err)
		}
		deposed[i] = upgraded
	}
	if len(deposed) == 0 {
		deposed = nil
	}

	return &ResourceState{
		Type:         old.Type,
		Dependencies: dependencies.([]string),
		Primary:      primary,
		Tainted:      tainted,
		Deposed:      deposed,
		Provider:     old.Provider,
	}, nil
}

type instanceStateV1 struct {
	// A unique ID for this resource. This is opaque to Terraform
	// and is only meant as a lookup mechanism for the providers.
	ID string `json:"id"`

	// Attributes are basic information about the resource. Any keys here
	// are accessible in variable format within Terraform configurations:
	// ${resourcetype.name.attribute}.
	Attributes map[string]string `json:"attributes,omitempty"`

	// Ephemeral is used to store any state associated with this instance
	// that is necessary for the Terraform run to complete, but is not
	// persisted to a state file.
	Ephemeral ephemeralStateV1 `json:"-"`

	// Meta is a simple K/V map that is persisted to the State but otherwise
	// ignored by Terraform core. It's meant to be used for accounting by
	// external client code.
	Meta map[string]string `json:"meta,omitempty"`
}

func (old *instanceStateV1) upgrade() (*InstanceState, error) {
	if old == nil {
		return nil, nil
	}

	attributes, err := copystructure.Copy(old.Attributes)
	if err != nil {
		return nil, fmt.Errorf("Error upgrading InstanceState V1: %v", err)
	}
	ephemeral, err := old.Ephemeral.upgrade()
	if err != nil {
		return nil, fmt.Errorf("Error upgrading InstanceState V1: %v", err)
	}
	meta, err := copystructure.Copy(old.Meta)
	if err != nil {
		return nil, fmt.Errorf("Error upgrading InstanceState V1: %v", err)
	}

	return &InstanceState{
		ID:         old.ID,
		Attributes: attributes.(map[string]string),
		Ephemeral:  *ephemeral,
		Meta:       meta.(map[string]string),
	}, nil
}

type ephemeralStateV1 struct {
	// ConnInfo is used for the providers to export information which is
	// used to connect to the resource for provisioning. For example,
	// this could contain SSH or WinRM credentials.
	ConnInfo map[string]string `json:"-"`
}

func (old *ephemeralStateV1) upgrade() (*EphemeralState, error) {
	connInfo, err := copystructure.Copy(old.ConnInfo)
	if err != nil {
		return nil, fmt.Errorf("Error upgrading EphemeralState V1: %v", err)
	}
	return &EphemeralState{
		ConnInfo: connInfo.(map[string]string),
	}, nil
}
