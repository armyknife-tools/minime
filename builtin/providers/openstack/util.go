package openstack

import (
	"fmt"

	"github.com/gophercloud/gophercloud"
	"github.com/hashicorp/terraform/helper/schema"
)

// CheckDeleted checks the error to see if it's a 404 (Not Found) and, if so,
// sets the resource ID to the empty string instead of throwing an error.
func CheckDeleted(d *schema.ResourceData, err error, msg string) error {
	if _, ok := err.(gophercloud.ErrDefault404); ok {
		d.SetId("")
		return nil
	}

	return fmt.Errorf("%s: %s", msg, err)
}

// BuildRequest takes an opts struct and builds a request body for
// Gophercloud to execute
func BuildRequest(opts interface{}, parent string) (map[string]interface{}, error) {
	b, err := gophercloud.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	if opts.ValueSpecs != nil {
		for k, v := range opts.ValueSpecs {
			b[k] = v
		}
		delete(b, "value_specs")
	}

	return map[string]interface{}{parent: b}, nil
}
