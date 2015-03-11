package schema

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform/terraform"
)

// Resource represents a thing in Terraform that has a set of configurable
// attributes and a lifecycle (create, read, update, delete).
//
// The Resource schema is an abstraction that allows provider writers to
// worry only about CRUD operations while off-loading validation, diff
// generation, etc. to this higher level library.
type Resource struct {
	// Schema is the schema for the configuration of this resource.
	//
	// The keys of this map are the configuration keys, and the values
	// describe the schema of the configuration value.
	//
	// The schema is used to represent both configurable data as well
	// as data that might be computed in the process of creating this
	// resource.
	Schema map[string]*Schema

	// SchemaVersion is the version number for this resource's Schema
	// definition. The current SchemaVersion stored in the state for each
	// resource. Provider authors can increment this version number
	// when Schema semantics change. If the State's SchemaVersion is less than
	// the current SchemaVersion, the InstanceState is yielded to the
	// MigrateState callback, where the provider can make whatever changes it
	// needs to update the state to be compatible to the latest version of the
	// Schema.
	//
	// When unset, SchemaVersion defaults to 0, so provider authors can start
	// their Versioning at any integer >= 1
	SchemaVersion int

	// MigrateState is responsible for updating an InstanceState with an old
	// version to the format expected by the current version of the Schema.
	//
	// It is called during Refresh if the State's stored SchemaVersion is less
	// than the current SchemaVersion of the Resource.
	//
	// The function is yielded the state's stored SchemaVersion and a pointer to
	// the InstanceState that needs updating, as well as the configured
	// provider's configured meta interface{}, in case the migration process
	// needs to make any remote API calls.
	MigrateState StateMigrateFunc

	// The functions below are the CRUD operations for this resource.
	//
	// The only optional operation is Update. If Update is not implemented,
	// then updates will not be supported for this resource.
	//
	// The ResourceData parameter in the functions below are used to
	// query configuration and changes for the resource as well as to set
	// the ID, computed data, etc.
	//
	// The interface{} parameter is the result of the ConfigureFunc in
	// the provider for this resource. If the provider does not define
	// a ConfigureFunc, this will be nil. This parameter should be used
	// to store API clients, configuration structures, etc.
	//
	// If any errors occur during each of the operation, an error should be
	// returned. If a resource was partially updated, be careful to enable
	// partial state mode for ResourceData and use it accordingly.
	//
	// Exists is a function that is called to check if a resource still
	// exists. If this returns false, then this will affect the diff
	// accordingly. If this function isn't set, it will not be called. It
	// is highly recommended to set it. The *ResourceData passed to Exists
	// should _not_ be modified.
	Create CreateFunc
	Read   ReadFunc
	Update UpdateFunc
	Delete DeleteFunc
	Exists ExistsFunc
}

// See Resource documentation.
type CreateFunc func(*ResourceData, interface{}) error

// See Resource documentation.
type ReadFunc func(*ResourceData, interface{}) error

// See Resource documentation.
type UpdateFunc func(*ResourceData, interface{}) error

// See Resource documentation.
type DeleteFunc func(*ResourceData, interface{}) error

// See Resource documentation.
type ExistsFunc func(*ResourceData, interface{}) (bool, error)

// See Resource documentation.
type StateMigrateFunc func(
	int, *terraform.InstanceState, interface{}) (*terraform.InstanceState, error)

// Apply creates, updates, and/or deletes a resource.
func (r *Resource) Apply(
	s *terraform.InstanceState,
	d *terraform.InstanceDiff,
	meta interface{}) (*terraform.InstanceState, error) {
	data, err := schemaMap(r.Schema).Data(s, d)
	if err != nil {
		return s, err
	}

	if s == nil {
		// The Terraform API dictates that this should never happen, but
		// it doesn't hurt to be safe in this case.
		s = new(terraform.InstanceState)
	}

	if d.Destroy || d.RequiresNew() {
		if s.ID != "" {
			// Destroy the resource since it is created
			if err := r.Delete(data, meta); err != nil {
				return data.State(), err
			}

			// Make sure the ID is gone.
			data.SetId("")
		}

		// If we're only destroying, and not creating, then return
		// now since we're done!
		if !d.RequiresNew() {
			return nil, nil
		}

		// Reset the data to be stateless since we just destroyed
		data, err = schemaMap(r.Schema).Data(nil, d)
		if err != nil {
			return nil, err
		}
	}

	err = nil
	if data.Id() == "" {
		// We're creating, it is a new resource.
		err = r.Create(data, meta)
	} else {
		if r.Update == nil {
			return s, fmt.Errorf("doesn't support update")
		}

		err = r.Update(data, meta)
	}

	return data.State(), err
}

// Diff returns a diff of this resource and is API compatible with the
// ResourceProvider interface.
func (r *Resource) Diff(
	s *terraform.InstanceState,
	c *terraform.ResourceConfig) (*terraform.InstanceDiff, error) {
	return schemaMap(r.Schema).Diff(s, c)
}

// Validate validates the resource configuration against the schema.
func (r *Resource) Validate(c *terraform.ResourceConfig) ([]string, []error) {
	return schemaMap(r.Schema).Validate(c)
}

// Refresh refreshes the state of the resource.
func (r *Resource) Refresh(
	s *terraform.InstanceState,
	meta interface{}) (*terraform.InstanceState, error) {
	if r.Exists != nil {
		// Make a copy of data so that if it is modified it doesn't
		// affect our Read later.
		data, err := schemaMap(r.Schema).Data(s, nil)
		if err != nil {
			return s, err
		}

		exists, err := r.Exists(data, meta)
		if err != nil {
			return s, err
		}
		if !exists {
			return nil, nil
		}
	}

	needsMigration, stateSchemaVersion := r.checkSchemaVersion(s)
	if needsMigration && r.MigrateState != nil {
		s, err := r.MigrateState(stateSchemaVersion, s, meta)
		if err != nil {
			return s, err
		}
	}

	data, err := schemaMap(r.Schema).Data(s, nil)
	if err != nil {
		return s, err
	}

	err = r.Read(data, meta)
	state := data.State()
	if state != nil && state.ID == "" {
		state = nil
	}

	if state != nil && r.SchemaVersion > 0 {
		if state.Meta == nil {
			state.Meta = make(map[string]string)
		}
		state.Meta["schema_version"] = strconv.Itoa(r.SchemaVersion)
	}

	return state, err
}

// InternalValidate should be called to validate the structure
// of the resource.
//
// This should be called in a unit test for any resource to verify
// before release that a resource is properly configured for use with
// this library.
//
// Provider.InternalValidate() will automatically call this for all of
// the resources it manages, so you don't need to call this manually if it
// is part of a Provider.
func (r *Resource) InternalValidate() error {
	if r == nil {
		return errors.New("resource is nil")
	}

	return schemaMap(r.Schema).InternalValidate()
}

// Determines if a given InstanceState needs to be migrated by checking the
// stored version number with the current SchemaVersion
func (r *Resource) checkSchemaVersion(is *terraform.InstanceState) (bool, int) {
	stateSchemaVersion, _ := strconv.Atoi(is.Meta["schema_version"])
	return stateSchemaVersion < r.SchemaVersion, stateSchemaVersion
}
