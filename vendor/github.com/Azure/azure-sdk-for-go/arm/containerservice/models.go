package containerservice

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
// Code generated by Microsoft (R) AutoRest Code Generator 1.0.0.0
// Changes may cause incorrect behavior and will be lost if the code is
// regenerated.

import (
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/to"
	"net/http"
)

// OchestratorTypes enumerates the values for ochestrator types.
type OchestratorTypes string

const (
	// Custom specifies the custom state for ochestrator types.
	Custom OchestratorTypes = "Custom"
	// DCOS specifies the dcos state for ochestrator types.
	DCOS OchestratorTypes = "DCOS"
	// Kubernetes specifies the kubernetes state for ochestrator types.
	Kubernetes OchestratorTypes = "Kubernetes"
	// Swarm specifies the swarm state for ochestrator types.
	Swarm OchestratorTypes = "Swarm"
)

// VMSizeTypes enumerates the values for vm size types.
type VMSizeTypes string

const (
	// StandardA0 specifies the standard a0 state for vm size types.
	StandardA0 VMSizeTypes = "Standard_A0"
	// StandardA1 specifies the standard a1 state for vm size types.
	StandardA1 VMSizeTypes = "Standard_A1"
	// StandardA10 specifies the standard a10 state for vm size types.
	StandardA10 VMSizeTypes = "Standard_A10"
	// StandardA11 specifies the standard a11 state for vm size types.
	StandardA11 VMSizeTypes = "Standard_A11"
	// StandardA2 specifies the standard a2 state for vm size types.
	StandardA2 VMSizeTypes = "Standard_A2"
	// StandardA3 specifies the standard a3 state for vm size types.
	StandardA3 VMSizeTypes = "Standard_A3"
	// StandardA4 specifies the standard a4 state for vm size types.
	StandardA4 VMSizeTypes = "Standard_A4"
	// StandardA5 specifies the standard a5 state for vm size types.
	StandardA5 VMSizeTypes = "Standard_A5"
	// StandardA6 specifies the standard a6 state for vm size types.
	StandardA6 VMSizeTypes = "Standard_A6"
	// StandardA7 specifies the standard a7 state for vm size types.
	StandardA7 VMSizeTypes = "Standard_A7"
	// StandardA8 specifies the standard a8 state for vm size types.
	StandardA8 VMSizeTypes = "Standard_A8"
	// StandardA9 specifies the standard a9 state for vm size types.
	StandardA9 VMSizeTypes = "Standard_A9"
	// StandardD1 specifies the standard d1 state for vm size types.
	StandardD1 VMSizeTypes = "Standard_D1"
	// StandardD11 specifies the standard d11 state for vm size types.
	StandardD11 VMSizeTypes = "Standard_D11"
	// StandardD11V2 specifies the standard d11v2 state for vm size types.
	StandardD11V2 VMSizeTypes = "Standard_D11_v2"
	// StandardD12 specifies the standard d12 state for vm size types.
	StandardD12 VMSizeTypes = "Standard_D12"
	// StandardD12V2 specifies the standard d12v2 state for vm size types.
	StandardD12V2 VMSizeTypes = "Standard_D12_v2"
	// StandardD13 specifies the standard d13 state for vm size types.
	StandardD13 VMSizeTypes = "Standard_D13"
	// StandardD13V2 specifies the standard d13v2 state for vm size types.
	StandardD13V2 VMSizeTypes = "Standard_D13_v2"
	// StandardD14 specifies the standard d14 state for vm size types.
	StandardD14 VMSizeTypes = "Standard_D14"
	// StandardD14V2 specifies the standard d14v2 state for vm size types.
	StandardD14V2 VMSizeTypes = "Standard_D14_v2"
	// StandardD1V2 specifies the standard d1v2 state for vm size types.
	StandardD1V2 VMSizeTypes = "Standard_D1_v2"
	// StandardD2 specifies the standard d2 state for vm size types.
	StandardD2 VMSizeTypes = "Standard_D2"
	// StandardD2V2 specifies the standard d2v2 state for vm size types.
	StandardD2V2 VMSizeTypes = "Standard_D2_v2"
	// StandardD3 specifies the standard d3 state for vm size types.
	StandardD3 VMSizeTypes = "Standard_D3"
	// StandardD3V2 specifies the standard d3v2 state for vm size types.
	StandardD3V2 VMSizeTypes = "Standard_D3_v2"
	// StandardD4 specifies the standard d4 state for vm size types.
	StandardD4 VMSizeTypes = "Standard_D4"
	// StandardD4V2 specifies the standard d4v2 state for vm size types.
	StandardD4V2 VMSizeTypes = "Standard_D4_v2"
	// StandardD5V2 specifies the standard d5v2 state for vm size types.
	StandardD5V2 VMSizeTypes = "Standard_D5_v2"
	// StandardDS1 specifies the standard ds1 state for vm size types.
	StandardDS1 VMSizeTypes = "Standard_DS1"
	// StandardDS11 specifies the standard ds11 state for vm size types.
	StandardDS11 VMSizeTypes = "Standard_DS11"
	// StandardDS12 specifies the standard ds12 state for vm size types.
	StandardDS12 VMSizeTypes = "Standard_DS12"
	// StandardDS13 specifies the standard ds13 state for vm size types.
	StandardDS13 VMSizeTypes = "Standard_DS13"
	// StandardDS14 specifies the standard ds14 state for vm size types.
	StandardDS14 VMSizeTypes = "Standard_DS14"
	// StandardDS2 specifies the standard ds2 state for vm size types.
	StandardDS2 VMSizeTypes = "Standard_DS2"
	// StandardDS3 specifies the standard ds3 state for vm size types.
	StandardDS3 VMSizeTypes = "Standard_DS3"
	// StandardDS4 specifies the standard ds4 state for vm size types.
	StandardDS4 VMSizeTypes = "Standard_DS4"
	// StandardG1 specifies the standard g1 state for vm size types.
	StandardG1 VMSizeTypes = "Standard_G1"
	// StandardG2 specifies the standard g2 state for vm size types.
	StandardG2 VMSizeTypes = "Standard_G2"
	// StandardG3 specifies the standard g3 state for vm size types.
	StandardG3 VMSizeTypes = "Standard_G3"
	// StandardG4 specifies the standard g4 state for vm size types.
	StandardG4 VMSizeTypes = "Standard_G4"
	// StandardG5 specifies the standard g5 state for vm size types.
	StandardG5 VMSizeTypes = "Standard_G5"
	// StandardGS1 specifies the standard gs1 state for vm size types.
	StandardGS1 VMSizeTypes = "Standard_GS1"
	// StandardGS2 specifies the standard gs2 state for vm size types.
	StandardGS2 VMSizeTypes = "Standard_GS2"
	// StandardGS3 specifies the standard gs3 state for vm size types.
	StandardGS3 VMSizeTypes = "Standard_GS3"
	// StandardGS4 specifies the standard gs4 state for vm size types.
	StandardGS4 VMSizeTypes = "Standard_GS4"
	// StandardGS5 specifies the standard gs5 state for vm size types.
	StandardGS5 VMSizeTypes = "Standard_GS5"
)

// AgentPoolProfile is profile for the container service agent pool.
type AgentPoolProfile struct {
	Name      *string     `json:"name,omitempty"`
	Count     *int32      `json:"count,omitempty"`
	VMSize    VMSizeTypes `json:"vmSize,omitempty"`
	DNSPrefix *string     `json:"dnsPrefix,omitempty"`
	Fqdn      *string     `json:"fqdn,omitempty"`
}

// ContainerService is container service.
type ContainerService struct {
	autorest.Response `json:"-"`
	ID                *string             `json:"id,omitempty"`
	Name              *string             `json:"name,omitempty"`
	Type              *string             `json:"type,omitempty"`
	Location          *string             `json:"location,omitempty"`
	Tags              *map[string]*string `json:"tags,omitempty"`
	*Properties       `json:"properties,omitempty"`
}

// CustomProfile is properties to configure a custom container service cluster.
type CustomProfile struct {
	Orchestrator *string `json:"orchestrator,omitempty"`
}

// DiagnosticsProfile is
type DiagnosticsProfile struct {
	VMDiagnostics *VMDiagnostics `json:"vmDiagnostics,omitempty"`
}

// LinuxProfile is profile for Linux VMs in the container service cluster.
type LinuxProfile struct {
	AdminUsername *string           `json:"adminUsername,omitempty"`
	SSH           *SSHConfiguration `json:"ssh,omitempty"`
}

// ListResult is the response from the List Container Services operation.
type ListResult struct {
	autorest.Response `json:"-"`
	Value             *[]ContainerService `json:"value,omitempty"`
	NextLink          *string             `json:"nextLink,omitempty"`
}

// ListResultPreparer prepares a request to retrieve the next set of results. It returns
// nil if no more results exist.
func (client ListResult) ListResultPreparer() (*http.Request, error) {
	if client.NextLink == nil || len(to.String(client.NextLink)) <= 0 {
		return nil, nil
	}
	return autorest.Prepare(&http.Request{},
		autorest.AsJSON(),
		autorest.AsGet(),
		autorest.WithBaseURL(to.String(client.NextLink)))
}

// MasterProfile is profile for the container service master.
type MasterProfile struct {
	Count     *int32  `json:"count,omitempty"`
	DNSPrefix *string `json:"dnsPrefix,omitempty"`
	Fqdn      *string `json:"fqdn,omitempty"`
}

// OrchestratorProfile is profile for the container service orchestrator.
type OrchestratorProfile struct {
	OrchestratorType OchestratorTypes `json:"orchestratorType,omitempty"`
}

// Properties is properties of the container service.
type Properties struct {
	ProvisioningState       *string                  `json:"provisioningState,omitempty"`
	OrchestratorProfile     *OrchestratorProfile     `json:"orchestratorProfile,omitempty"`
	CustomProfile           *CustomProfile           `json:"customProfile,omitempty"`
	ServicePrincipalProfile *ServicePrincipalProfile `json:"servicePrincipalProfile,omitempty"`
	MasterProfile           *MasterProfile           `json:"masterProfile,omitempty"`
	AgentPoolProfiles       *[]AgentPoolProfile      `json:"agentPoolProfiles,omitempty"`
	WindowsProfile          *WindowsProfile          `json:"windowsProfile,omitempty"`
	LinuxProfile            *LinuxProfile            `json:"linuxProfile,omitempty"`
	DiagnosticsProfile      *DiagnosticsProfile      `json:"diagnosticsProfile,omitempty"`
}

// Resource is the Resource model definition.
type Resource struct {
	ID       *string             `json:"id,omitempty"`
	Name     *string             `json:"name,omitempty"`
	Type     *string             `json:"type,omitempty"`
	Location *string             `json:"location,omitempty"`
	Tags     *map[string]*string `json:"tags,omitempty"`
}

// ServicePrincipalProfile is information about a service principal identity
// for the cluster to use for manipulating Azure APIs.
type ServicePrincipalProfile struct {
	ClientID *string `json:"clientId,omitempty"`
	Secret   *string `json:"secret,omitempty"`
}

// SSHConfiguration is sSH configuration for Linux-based VMs running on Azure.
type SSHConfiguration struct {
	PublicKeys *[]SSHPublicKey `json:"publicKeys,omitempty"`
}

// SSHPublicKey is contains information about SSH certificate public key data.
type SSHPublicKey struct {
	KeyData *string `json:"keyData,omitempty"`
}

// VMDiagnostics is profile for diagnostics on the container service VMs.
type VMDiagnostics struct {
	Enabled    *bool   `json:"enabled,omitempty"`
	StorageURI *string `json:"storageUri,omitempty"`
}

// WindowsProfile is profile for Windows VMs in the container service cluster.
type WindowsProfile struct {
	AdminUsername *string `json:"adminUsername,omitempty"`
	AdminPassword *string `json:"adminPassword,omitempty"`
}
