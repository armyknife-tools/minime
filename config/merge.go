package config

import (
	"fmt"
)

// Merge merges two configurations into a single configuration.
//
// Merge allows for the two configurations to have duplicate resources,
// because the resources will be merged. This differs from a single
// Config which must only have unique resources.
func Merge(c1, c2 *Config) (*Config, error) {
	c := new(Config)

	// Merge unknown keys
	unknowns := make(map[string]struct{})
	for _, k := range c1.unknownKeys {
		unknowns[k] = struct{}{}
	}
	for _, k := range c2.unknownKeys {
		unknowns[k] = struct{}{}
	}
	for k, _ := range unknowns {
		c.unknownKeys = append(c.unknownKeys, k)
	}

	// Merge variables: Variable merging is quite simple. Set fields in
	// later set variables override those earlier.
	if len(c1.Variables) > 0 || len(c2.Variables) > 0 {
		c.Variables = make([]*Variable, 0, len(c1.Variables)+len(c2.Variables))
		varMap := make(map[string]*Variable)
		for _, v := range c1.Variables {
			varMap[v.Name] = v
		}
		for _, v2 := range c2.Variables {
			v1, ok := varMap[v2.Name]
			if ok {
				if v2.Default == "" {
					v2.Default = v1.Default
				}
				if v2.Description == "" {
					v2.Description = v1.Description
				}
			}

			varMap[v2.Name] = v2
		}
		for _, v := range varMap {
			c.Variables = append(c.Variables, v)
		}
	}

	// Merge outputs: If they collide, just take the latest one for now. In
	// the future, we might provide smarter merge functionality.
	if len(c1.Outputs) > 0 || len(c2.Outputs) > 0 {
		c.Outputs = make([]*Output, 0, len(c1.Outputs)+len(c2.Outputs))
		m := make(map[string]*Output)
		for _, v := range c1.Outputs {
			m[v.Name] = v
		}
		for _, v := range c2.Outputs {
			m[v.Name] = v
		}
		for _, v := range m {
			c.Outputs = append(c.Outputs, v)
		}
	}

	// Merge provider configs: If they collide, we just take the latest one
	// for now. In the future, we might provide smarter merge functionality.
	if len(c1.ProviderConfigs) > 0 || len(c2.ProviderConfigs) > 0 {
		m := make(map[string]*ProviderConfig)
		for _, v := range c1.ProviderConfigs {
			m[v.Name] = v
		}
		for _, v := range c2.ProviderConfigs {
			m[v.Name] = v
		}

		c.ProviderConfigs = make([]*ProviderConfig, 0, len(m))
		for _, v := range m {
			c.ProviderConfigs = append(c.ProviderConfigs, v)
		}
	}

	// Merge resources: If they collide, we just take the latest one
	// for now. In the future, we might provide smarter merge functionality.
	resources := make(map[string]*Resource)
	for _, r := range c1.Resources {
		id := fmt.Sprintf("%s[%s]", r.Type, r.Name)
		resources[id] = r
	}
	for _, r := range c2.Resources {
		id := fmt.Sprintf("%s[%s]", r.Type, r.Name)
		resources[id] = r
	}

	c.Resources = make([]*Resource, 0, len(resources))
	for _, r := range resources {
		c.Resources = append(c.Resources, r)
	}

	return c, nil
}
