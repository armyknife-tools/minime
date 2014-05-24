// The config package is responsible for loading and validating the
// configuration.
package config

import (
	"strings"
)

// Config is the configuration that comes from loading a collection
// of Terraform templates.
type Config struct {
	Variables map[string]Variable
	Resources []Resource
}

// A resource represents a single Terraform resource in the configuration.
// A Terraform resource is something that represents some component that
// can be created and managed, and has some properties associated with it.
type Resource struct {
	Name      string
	Type      string
	Config    map[string]interface{}
	Variables map[string]InterpolatedVariable
}

type Variable struct {
	Default     string
	Description string
}

// An InterpolatedVariable is a variable that is embedded within a string
// in the configuration, such as "hello ${world}" (world in this case is
// an interpolated variable).
//
// These variables can come from a variety of sources, represented by
// implementations of this interface.
type InterpolatedVariable interface {
	FullKey() string
}

// A ResourceVariable is a variable that is referencing the field
// of a resource, such as "${aws_instance.foo.ami}"
type ResourceVariable struct {
	Type  string
	Name  string
	Field string

	key string
}

// A UserVariable is a variable that is referencing a user variable
// that is inputted from outside the configuration. This looks like
// "${var.foo}"
type UserVariable struct {
	Name string

	key  string
}

func NewResourceVariable(key string) (*ResourceVariable, error) {
	parts := strings.SplitN(key, ".", 3)
	return &ResourceVariable{
		Type:  parts[0],
		Name:  parts[1],
		Field: parts[2],
		key:   key,
	}, nil
}

func (v *ResourceVariable) FullKey() string {
	return v.key
}

func NewUserVariable(key string) (*UserVariable, error) {
	name := key[len("var."):]
	return &UserVariable{
		key:  key,
		Name: name,
	}, nil
}

func (v *UserVariable) FullKey() string {
	return v.key
}
