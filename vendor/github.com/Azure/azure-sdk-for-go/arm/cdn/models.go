package cdn

// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Code generated by Microsoft (R) AutoRest Code Generator 1.0.1.0
// Changes may cause incorrect behavior and will be lost if the code is
// regenerated.

import (
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/to"
	"net/http"
)

// CustomDomainResourceState enumerates the values for custom domain resource
// state.
type CustomDomainResourceState string

const (
	// Active specifies the active state for custom domain resource state.
	Active CustomDomainResourceState = "Active"
	// Creating specifies the creating state for custom domain resource state.
	Creating CustomDomainResourceState = "Creating"
	// Deleting specifies the deleting state for custom domain resource state.
	Deleting CustomDomainResourceState = "Deleting"
)

// CustomHTTPSProvisioningState enumerates the values for custom https
// provisioning state.
type CustomHTTPSProvisioningState string

const (
	// Disabled specifies the disabled state for custom https provisioning
	// state.
	Disabled CustomHTTPSProvisioningState = "Disabled"
	// Disabling specifies the disabling state for custom https provisioning
	// state.
	Disabling CustomHTTPSProvisioningState = "Disabling"
	// Enabled specifies the enabled state for custom https provisioning state.
	Enabled CustomHTTPSProvisioningState = "Enabled"
	// Enabling specifies the enabling state for custom https provisioning
	// state.
	Enabling CustomHTTPSProvisioningState = "Enabling"
	// Failed specifies the failed state for custom https provisioning state.
	Failed CustomHTTPSProvisioningState = "Failed"
)

// EndpointResourceState enumerates the values for endpoint resource state.
type EndpointResourceState string

const (
	// EndpointResourceStateCreating specifies the endpoint resource state
	// creating state for endpoint resource state.
	EndpointResourceStateCreating EndpointResourceState = "Creating"
	// EndpointResourceStateDeleting specifies the endpoint resource state
	// deleting state for endpoint resource state.
	EndpointResourceStateDeleting EndpointResourceState = "Deleting"
	// EndpointResourceStateRunning specifies the endpoint resource state
	// running state for endpoint resource state.
	EndpointResourceStateRunning EndpointResourceState = "Running"
	// EndpointResourceStateStarting specifies the endpoint resource state
	// starting state for endpoint resource state.
	EndpointResourceStateStarting EndpointResourceState = "Starting"
	// EndpointResourceStateStopped specifies the endpoint resource state
	// stopped state for endpoint resource state.
	EndpointResourceStateStopped EndpointResourceState = "Stopped"
	// EndpointResourceStateStopping specifies the endpoint resource state
	// stopping state for endpoint resource state.
	EndpointResourceStateStopping EndpointResourceState = "Stopping"
)

// GeoFilterActions enumerates the values for geo filter actions.
type GeoFilterActions string

const (
	// Allow specifies the allow state for geo filter actions.
	Allow GeoFilterActions = "Allow"
	// Block specifies the block state for geo filter actions.
	Block GeoFilterActions = "Block"
)

// OriginResourceState enumerates the values for origin resource state.
type OriginResourceState string

const (
	// OriginResourceStateActive specifies the origin resource state active
	// state for origin resource state.
	OriginResourceStateActive OriginResourceState = "Active"
	// OriginResourceStateCreating specifies the origin resource state creating
	// state for origin resource state.
	OriginResourceStateCreating OriginResourceState = "Creating"
	// OriginResourceStateDeleting specifies the origin resource state deleting
	// state for origin resource state.
	OriginResourceStateDeleting OriginResourceState = "Deleting"
)

// ProfileResourceState enumerates the values for profile resource state.
type ProfileResourceState string

const (
	// ProfileResourceStateActive specifies the profile resource state active
	// state for profile resource state.
	ProfileResourceStateActive ProfileResourceState = "Active"
	// ProfileResourceStateCreating specifies the profile resource state
	// creating state for profile resource state.
	ProfileResourceStateCreating ProfileResourceState = "Creating"
	// ProfileResourceStateDeleting specifies the profile resource state
	// deleting state for profile resource state.
	ProfileResourceStateDeleting ProfileResourceState = "Deleting"
	// ProfileResourceStateDisabled specifies the profile resource state
	// disabled state for profile resource state.
	ProfileResourceStateDisabled ProfileResourceState = "Disabled"
)

// QueryStringCachingBehavior enumerates the values for query string caching
// behavior.
type QueryStringCachingBehavior string

const (
	// BypassCaching specifies the bypass caching state for query string
	// caching behavior.
	BypassCaching QueryStringCachingBehavior = "BypassCaching"
	// IgnoreQueryString specifies the ignore query string state for query
	// string caching behavior.
	IgnoreQueryString QueryStringCachingBehavior = "IgnoreQueryString"
	// NotSet specifies the not set state for query string caching behavior.
	NotSet QueryStringCachingBehavior = "NotSet"
	// UseQueryString specifies the use query string state for query string
	// caching behavior.
	UseQueryString QueryStringCachingBehavior = "UseQueryString"
)

// ResourceType enumerates the values for resource type.
type ResourceType string

const (
	// MicrosoftCdnProfilesEndpoints specifies the microsoft cdn profiles
	// endpoints state for resource type.
	MicrosoftCdnProfilesEndpoints ResourceType = "Microsoft.Cdn/Profiles/Endpoints"
)

// SkuName enumerates the values for sku name.
type SkuName string

const (
	// CustomVerizon specifies the custom verizon state for sku name.
	CustomVerizon SkuName = "Custom_Verizon"
	// PremiumVerizon specifies the premium verizon state for sku name.
	PremiumVerizon SkuName = "Premium_Verizon"
	// StandardAkamai specifies the standard akamai state for sku name.
	StandardAkamai SkuName = "Standard_Akamai"
	// StandardChinaCdn specifies the standard china cdn state for sku name.
	StandardChinaCdn SkuName = "Standard_ChinaCdn"
	// StandardVerizon specifies the standard verizon state for sku name.
	StandardVerizon SkuName = "Standard_Verizon"
)

// CheckNameAvailabilityInput is input of CheckNameAvailability API.
type CheckNameAvailabilityInput struct {
	Name *string `json:"name,omitempty"`
	Type *string `json:"type,omitempty"`
}

// CheckNameAvailabilityOutput is output of check name availability API.
type CheckNameAvailabilityOutput struct {
	autorest.Response `json:"-"`
	NameAvailable     *bool   `json:"nameAvailable,omitempty"`
	Reason            *string `json:"reason,omitempty"`
	Message           *string `json:"message,omitempty"`
}

// CidrIPAddress is cIDR Ip address
type CidrIPAddress struct {
	BaseIPAddress *string `json:"baseIpAddress,omitempty"`
	PrefixLength  *int32  `json:"prefixLength,omitempty"`
}

// CustomDomain is customer provided domain for branding purposes, e.g.
// www.consoto.com.
type CustomDomain struct {
	autorest.Response       `json:"-"`
	ID                      *string             `json:"id,omitempty"`
	Name                    *string             `json:"name,omitempty"`
	Type                    *string             `json:"type,omitempty"`
	Location                *string             `json:"location,omitempty"`
	Tags                    *map[string]*string `json:"tags,omitempty"`
	*CustomDomainProperties `json:"properties,omitempty"`
}

// CustomDomainListResult is result of the request to list custom domains. It
// contains a list of custom domain objects and a URL link to get the next set
// of results.
type CustomDomainListResult struct {
	autorest.Response `json:"-"`
	Value             *[]CustomDomain `json:"value,omitempty"`
	NextLink          *string         `json:"nextLink,omitempty"`
}

// CustomDomainListResultPreparer prepares a request to retrieve the next set of results. It returns
// nil if no more results exist.
func (client CustomDomainListResult) CustomDomainListResultPreparer() (*http.Request, error) {
	if client.NextLink == nil || len(to.String(client.NextLink)) <= 0 {
		return nil, nil
	}
	return autorest.Prepare(&http.Request{},
		autorest.AsJSON(),
		autorest.AsGet(),
		autorest.WithBaseURL(to.String(client.NextLink)))
}

// CustomDomainParameters is the customDomain JSON object required for custom
// domain creation or update.
type CustomDomainParameters struct {
	*CustomDomainPropertiesParameters `json:"properties,omitempty"`
}

// CustomDomainProperties is the JSON object that contains the properties of
// the custom domain to create.
type CustomDomainProperties struct {
	HostName                     *string                      `json:"hostName,omitempty"`
	ResourceState                CustomDomainResourceState    `json:"resourceState,omitempty"`
	CustomHTTPSProvisioningState CustomHTTPSProvisioningState `json:"customHttpsProvisioningState,omitempty"`
	ValidationData               *string                      `json:"validationData,omitempty"`
	ProvisioningState            *string                      `json:"provisioningState,omitempty"`
}

// CustomDomainPropertiesParameters is the JSON object that contains the
// properties of the custom domain to create.
type CustomDomainPropertiesParameters struct {
	HostName *string `json:"hostName,omitempty"`
}

// DeepCreatedOrigin is origin to be added when creating a CDN endpoint.
type DeepCreatedOrigin struct {
	Name                         *string `json:"name,omitempty"`
	*DeepCreatedOriginProperties `json:"properties,omitempty"`
}

// DeepCreatedOriginProperties is properties of origin Properties of the origin
// created on the CDN endpoint.
type DeepCreatedOriginProperties struct {
	HostName  *string `json:"hostName,omitempty"`
	HTTPPort  *int32  `json:"httpPort,omitempty"`
	HTTPSPort *int32  `json:"httpsPort,omitempty"`
}

// EdgeNode is edge node of CDN service.
type EdgeNode struct {
	ID                  *string             `json:"id,omitempty"`
	Name                *string             `json:"name,omitempty"`
	Type                *string             `json:"type,omitempty"`
	Location            *string             `json:"location,omitempty"`
	Tags                *map[string]*string `json:"tags,omitempty"`
	*EdgeNodeProperties `json:"properties,omitempty"`
}

// EdgeNodeProperties is the JSON object that contains the properties required
// to create an edgenode.
type EdgeNodeProperties struct {
	IPAddressGroups *[]IPAddressGroup `json:"ipAddressGroups,omitempty"`
}

// EdgenodeResult is result of the request to list CDN edgenodes. It contains a
// list of ip address group and a URL link to get the next set of results.
type EdgenodeResult struct {
	autorest.Response `json:"-"`
	Value             *[]EdgeNode `json:"value,omitempty"`
	NextLink          *string     `json:"nextLink,omitempty"`
}

// EdgenodeResultPreparer prepares a request to retrieve the next set of results. It returns
// nil if no more results exist.
func (client EdgenodeResult) EdgenodeResultPreparer() (*http.Request, error) {
	if client.NextLink == nil || len(to.String(client.NextLink)) <= 0 {
		return nil, nil
	}
	return autorest.Prepare(&http.Request{},
		autorest.AsJSON(),
		autorest.AsGet(),
		autorest.WithBaseURL(to.String(client.NextLink)))
}

// Endpoint is cDN endpoint is the entity within a CDN profile containing
// configuration information such as origin, protocol, content caching and
// delivery behavior. The CDN endpoint uses the URL format
// <endpointname>.azureedge.net.
type Endpoint struct {
	autorest.Response   `json:"-"`
	ID                  *string             `json:"id,omitempty"`
	Name                *string             `json:"name,omitempty"`
	Type                *string             `json:"type,omitempty"`
	Location            *string             `json:"location,omitempty"`
	Tags                *map[string]*string `json:"tags,omitempty"`
	*EndpointProperties `json:"properties,omitempty"`
}

// EndpointListResult is result of the request to list endpoints. It contains a
// list of endpoint objects and a URL link to get the the next set of results.
type EndpointListResult struct {
	autorest.Response `json:"-"`
	Value             *[]Endpoint `json:"value,omitempty"`
	NextLink          *string     `json:"nextLink,omitempty"`
}

// EndpointListResultPreparer prepares a request to retrieve the next set of results. It returns
// nil if no more results exist.
func (client EndpointListResult) EndpointListResultPreparer() (*http.Request, error) {
	if client.NextLink == nil || len(to.String(client.NextLink)) <= 0 {
		return nil, nil
	}
	return autorest.Prepare(&http.Request{},
		autorest.AsJSON(),
		autorest.AsGet(),
		autorest.WithBaseURL(to.String(client.NextLink)))
}

// EndpointProperties is the JSON object that contains the properties required
// to create an endpoint.
type EndpointProperties struct {
	OriginHostHeader           *string                    `json:"originHostHeader,omitempty"`
	OriginPath                 *string                    `json:"originPath,omitempty"`
	ContentTypesToCompress     *[]string                  `json:"contentTypesToCompress,omitempty"`
	IsCompressionEnabled       *bool                      `json:"isCompressionEnabled,omitempty"`
	IsHTTPAllowed              *bool                      `json:"isHttpAllowed,omitempty"`
	IsHTTPSAllowed             *bool                      `json:"isHttpsAllowed,omitempty"`
	QueryStringCachingBehavior QueryStringCachingBehavior `json:"queryStringCachingBehavior,omitempty"`
	OptimizationType           *string                    `json:"optimizationType,omitempty"`
	GeoFilters                 *[]GeoFilter               `json:"geoFilters,omitempty"`
	HostName                   *string                    `json:"hostName,omitempty"`
	Origins                    *[]DeepCreatedOrigin       `json:"origins,omitempty"`
	ResourceState              EndpointResourceState      `json:"resourceState,omitempty"`
	ProvisioningState          *string                    `json:"provisioningState,omitempty"`
}

// EndpointPropertiesUpdateParameters is result of the request to list
// endpoints. It contains a list of endpoints and a URL link to get the next
// set of results.
type EndpointPropertiesUpdateParameters struct {
	OriginHostHeader           *string                    `json:"originHostHeader,omitempty"`
	OriginPath                 *string                    `json:"originPath,omitempty"`
	ContentTypesToCompress     *[]string                  `json:"contentTypesToCompress,omitempty"`
	IsCompressionEnabled       *bool                      `json:"isCompressionEnabled,omitempty"`
	IsHTTPAllowed              *bool                      `json:"isHttpAllowed,omitempty"`
	IsHTTPSAllowed             *bool                      `json:"isHttpsAllowed,omitempty"`
	QueryStringCachingBehavior QueryStringCachingBehavior `json:"queryStringCachingBehavior,omitempty"`
	OptimizationType           *string                    `json:"optimizationType,omitempty"`
	GeoFilters                 *[]GeoFilter               `json:"geoFilters,omitempty"`
}

// EndpointUpdateParameters is properties required to create a new endpoint.
type EndpointUpdateParameters struct {
	Tags                                *map[string]*string `json:"tags,omitempty"`
	*EndpointPropertiesUpdateParameters `json:"properties,omitempty"`
}

// ErrorResponse is error reponse indicates CDN service is not able to process
// the incoming request. The reason is provided in the error message.
type ErrorResponse struct {
	Code    *string `json:"code,omitempty"`
	Message *string `json:"message,omitempty"`
}

// GeoFilter is rules defining user geo access within a CDN endpoint.
type GeoFilter struct {
	RelativePath *string          `json:"relativePath,omitempty"`
	Action       GeoFilterActions `json:"action,omitempty"`
	CountryCodes *[]string        `json:"countryCodes,omitempty"`
}

// IPAddressGroup is cDN Ip address group
type IPAddressGroup struct {
	DeliveryRegion *string          `json:"deliveryRegion,omitempty"`
	Ipv4Addresses  *[]CidrIPAddress `json:"ipv4Addresses,omitempty"`
	Ipv6Addresses  *[]CidrIPAddress `json:"ipv6Addresses,omitempty"`
}

// LoadParameters is parameters required for content load.
type LoadParameters struct {
	ContentPaths *[]string `json:"contentPaths,omitempty"`
}

// Operation is cDN REST API operation
type Operation struct {
	Name    *string           `json:"name,omitempty"`
	Display *OperationDisplay `json:"display,omitempty"`
}

// OperationDisplay is the object that represents the operation.
type OperationDisplay struct {
	Provider  *string `json:"provider,omitempty"`
	Resource  *string `json:"resource,omitempty"`
	Operation *string `json:"operation,omitempty"`
}

// OperationListResult is result of the request to list CDN operations. It
// contains a list of operations and a URL link to get the next set of results.
type OperationListResult struct {
	autorest.Response `json:"-"`
	Value             *[]Operation `json:"value,omitempty"`
	NextLink          *string      `json:"nextLink,omitempty"`
}

// OperationListResultPreparer prepares a request to retrieve the next set of results. It returns
// nil if no more results exist.
func (client OperationListResult) OperationListResultPreparer() (*http.Request, error) {
	if client.NextLink == nil || len(to.String(client.NextLink)) <= 0 {
		return nil, nil
	}
	return autorest.Prepare(&http.Request{},
		autorest.AsJSON(),
		autorest.AsGet(),
		autorest.WithBaseURL(to.String(client.NextLink)))
}

// Origin is cDN origin is the source of the content being delivered via CDN.
// When the edge nodes represented by an endpoint do not have the requested
// content cached, they attempt to fetch it from one or more of the configured
// origins.
type Origin struct {
	autorest.Response `json:"-"`
	ID                *string             `json:"id,omitempty"`
	Name              *string             `json:"name,omitempty"`
	Type              *string             `json:"type,omitempty"`
	Location          *string             `json:"location,omitempty"`
	Tags              *map[string]*string `json:"tags,omitempty"`
	*OriginProperties `json:"properties,omitempty"`
}

// OriginListResult is result of the request to list origins. It contains a
// list of origin objects and a URL link to get the next set of results.
type OriginListResult struct {
	autorest.Response `json:"-"`
	Value             *[]Origin `json:"value,omitempty"`
	NextLink          *string   `json:"nextLink,omitempty"`
}

// OriginListResultPreparer prepares a request to retrieve the next set of results. It returns
// nil if no more results exist.
func (client OriginListResult) OriginListResultPreparer() (*http.Request, error) {
	if client.NextLink == nil || len(to.String(client.NextLink)) <= 0 {
		return nil, nil
	}
	return autorest.Prepare(&http.Request{},
		autorest.AsJSON(),
		autorest.AsGet(),
		autorest.WithBaseURL(to.String(client.NextLink)))
}

// OriginProperties is the JSON object that contains the properties of the
// origin to create.
type OriginProperties struct {
	HostName          *string             `json:"hostName,omitempty"`
	HTTPPort          *int32              `json:"httpPort,omitempty"`
	HTTPSPort         *int32              `json:"httpsPort,omitempty"`
	ResourceState     OriginResourceState `json:"resourceState,omitempty"`
	ProvisioningState *string             `json:"provisioningState,omitempty"`
}

// OriginPropertiesParameters is the JSON object that contains the properties
// of the origin to create.
type OriginPropertiesParameters struct {
	HostName  *string `json:"hostName,omitempty"`
	HTTPPort  *int32  `json:"httpPort,omitempty"`
	HTTPSPort *int32  `json:"httpsPort,omitempty"`
}

// OriginUpdateParameters is origin properties needed for origin creation or
// update.
type OriginUpdateParameters struct {
	*OriginPropertiesParameters `json:"properties,omitempty"`
}

// Profile is cDN profile represents the top level resource and the entry point
// into the CDN API. This allows users to set up a logical grouping of
// endpoints in addition to creating shared configuration settings and
// selecting pricing tiers and providers.
type Profile struct {
	autorest.Response  `json:"-"`
	ID                 *string             `json:"id,omitempty"`
	Name               *string             `json:"name,omitempty"`
	Type               *string             `json:"type,omitempty"`
	Location           *string             `json:"location,omitempty"`
	Tags               *map[string]*string `json:"tags,omitempty"`
	Sku                *Sku                `json:"sku,omitempty"`
	*ProfileProperties `json:"properties,omitempty"`
}

// ProfileListResult is result of the request to list profiles. It contains a
// list of profile objects and a URL link to get the the next set of results.
type ProfileListResult struct {
	autorest.Response `json:"-"`
	Value             *[]Profile `json:"value,omitempty"`
	NextLink          *string    `json:"nextLink,omitempty"`
}

// ProfileListResultPreparer prepares a request to retrieve the next set of results. It returns
// nil if no more results exist.
func (client ProfileListResult) ProfileListResultPreparer() (*http.Request, error) {
	if client.NextLink == nil || len(to.String(client.NextLink)) <= 0 {
		return nil, nil
	}
	return autorest.Prepare(&http.Request{},
		autorest.AsJSON(),
		autorest.AsGet(),
		autorest.WithBaseURL(to.String(client.NextLink)))
}

// ProfileProperties is the JSON object that contains the properties required
// to create a profile.
type ProfileProperties struct {
	ResourceState     ProfileResourceState `json:"resourceState,omitempty"`
	ProvisioningState *string              `json:"provisioningState,omitempty"`
}

// ProfileUpdateParameters is properties required to update a profile.
type ProfileUpdateParameters struct {
	Tags *map[string]*string `json:"tags,omitempty"`
}

// PurgeParameters is parameters required for content purge.
type PurgeParameters struct {
	ContentPaths *[]string `json:"contentPaths,omitempty"`
}

// Resource is the Resource definition.
type Resource struct {
	ID       *string             `json:"id,omitempty"`
	Name     *string             `json:"name,omitempty"`
	Type     *string             `json:"type,omitempty"`
	Location *string             `json:"location,omitempty"`
	Tags     *map[string]*string `json:"tags,omitempty"`
}

// ResourceUsage is output of check resource usage API.
type ResourceUsage struct {
	ResourceType *string `json:"resourceType,omitempty"`
	Unit         *string `json:"unit,omitempty"`
	CurrentValue *int32  `json:"currentValue,omitempty"`
	Limit        *int32  `json:"limit,omitempty"`
}

// ResourceUsageListResult is output of check resource usage API.
type ResourceUsageListResult struct {
	autorest.Response `json:"-"`
	Value             *[]ResourceUsage `json:"value,omitempty"`
	NextLink          *string          `json:"nextLink,omitempty"`
}

// ResourceUsageListResultPreparer prepares a request to retrieve the next set of results. It returns
// nil if no more results exist.
func (client ResourceUsageListResult) ResourceUsageListResultPreparer() (*http.Request, error) {
	if client.NextLink == nil || len(to.String(client.NextLink)) <= 0 {
		return nil, nil
	}
	return autorest.Prepare(&http.Request{},
		autorest.AsJSON(),
		autorest.AsGet(),
		autorest.WithBaseURL(to.String(client.NextLink)))
}

// Sku is the pricing tier (defines a CDN provider, feature list and rate) of
// the CDN profile.
type Sku struct {
	Name SkuName `json:"name,omitempty"`
}

// SsoURI is sSO URI required to login to the supplemental portal.
type SsoURI struct {
	autorest.Response `json:"-"`
	SsoURIValue       *string `json:"ssoUriValue,omitempty"`
}

// ValidateCustomDomainInput is input of the custom domain to be validated for
// DNS mapping.
type ValidateCustomDomainInput struct {
	HostName *string `json:"hostName,omitempty"`
}

// ValidateCustomDomainOutput is output of custom domain validation.
type ValidateCustomDomainOutput struct {
	autorest.Response     `json:"-"`
	CustomDomainValidated *bool   `json:"customDomainValidated,omitempty"`
	Reason                *string `json:"reason,omitempty"`
	Message               *string `json:"message,omitempty"`
}
