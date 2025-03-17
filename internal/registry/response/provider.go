// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2023 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package response

// Provider represents a provider in the registry.
type Provider struct {
	ID          string `json:"id"`
	Namespace   string `json:"namespace"`
	Name        string `json:"name"`
	Version     string `json:"version,omitempty"`
	Owner       string `json:"owner,omitempty"`
	Description string `json:"description,omitempty"`
	Downloads   int    `json:"downloads,omitempty"`
	Logo        string `json:"logo_url,omitempty"`
	Source      string `json:"source,omitempty"`
	Tier        string `json:"tier,omitempty"`
	Published   string `json:"published_at,omitempty"`
	// Additional fields can be added as needed
}

// ProviderList is the response structure for a pageable list of providers.
type ProviderList struct {
	Meta      PaginationMeta `json:"meta"`
	Providers []Provider     `json:"providers"`
}
