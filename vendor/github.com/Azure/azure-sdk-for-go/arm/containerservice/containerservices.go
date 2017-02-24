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
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/validation"
	"net/http"
)

// ContainerServicesClient is the the Container Service Client.
type ContainerServicesClient struct {
	ManagementClient
}

// NewContainerServicesClient creates an instance of the
// ContainerServicesClient client.
func NewContainerServicesClient(subscriptionID string) ContainerServicesClient {
	return NewContainerServicesClientWithBaseURI(DefaultBaseURI, subscriptionID)
}

// NewContainerServicesClientWithBaseURI creates an instance of the
// ContainerServicesClient client.
func NewContainerServicesClientWithBaseURI(baseURI string, subscriptionID string) ContainerServicesClient {
	return ContainerServicesClient{NewWithBaseURI(baseURI, subscriptionID)}
}

// CreateOrUpdate creates or updates a container service with the specified
// configuration of orchestrator, masters, and agents. This method may poll for
// completion. Polling can be canceled by passing the cancel channel argument.
// The channel will be used to cancel polling and any outstanding HTTP
// requests.
//
// resourceGroupName is the name of the resource group. containerServiceName is
// the name of the container service in the specified subscription and resource
// group. parameters is parameters supplied to the Create or Update a Container
// Service operation.
func (client ContainerServicesClient) CreateOrUpdate(resourceGroupName string, containerServiceName string, parameters ContainerService, cancel <-chan struct{}) (result autorest.Response, err error) {
	if err := validation.Validate([]validation.Validation{
		{TargetValue: parameters,
			Constraints: []validation.Constraint{{Target: "parameters.Properties", Name: validation.Null, Rule: false,
				Chain: []validation.Constraint{{Target: "parameters.Properties.CustomProfile", Name: validation.Null, Rule: false,
					Chain: []validation.Constraint{{Target: "parameters.Properties.CustomProfile.Orchestrator", Name: validation.Null, Rule: true, Chain: nil}}},
					{Target: "parameters.Properties.ServicePrincipalProfile", Name: validation.Null, Rule: false,
						Chain: []validation.Constraint{{Target: "parameters.Properties.ServicePrincipalProfile.ClientID", Name: validation.Null, Rule: true, Chain: nil},
							{Target: "parameters.Properties.ServicePrincipalProfile.Secret", Name: validation.Null, Rule: true, Chain: nil},
						}},
					{Target: "parameters.Properties.MasterProfile", Name: validation.Null, Rule: true,
						Chain: []validation.Constraint{{Target: "parameters.Properties.MasterProfile.DNSPrefix", Name: validation.Null, Rule: true, Chain: nil}}},
					{Target: "parameters.Properties.AgentPoolProfiles", Name: validation.Null, Rule: true, Chain: nil},
					{Target: "parameters.Properties.WindowsProfile", Name: validation.Null, Rule: false,
						Chain: []validation.Constraint{{Target: "parameters.Properties.WindowsProfile.AdminUsername", Name: validation.Null, Rule: true,
							Chain: []validation.Constraint{{Target: "parameters.Properties.WindowsProfile.AdminUsername", Name: validation.Pattern, Rule: `^[a-zA-Z0-9]+([._]?[a-zA-Z0-9]+)*$`, Chain: nil}}},
							{Target: "parameters.Properties.WindowsProfile.AdminPassword", Name: validation.Null, Rule: true,
								Chain: []validation.Constraint{{Target: "parameters.Properties.WindowsProfile.AdminPassword", Name: validation.Pattern, Rule: `^(?=.*[a-z])(?=.*[A-Z])(?=.*[!@#$%\^&\*\(\)])[a-zA-Z\d!@#$%\^&\*\(\)]{12,123}$`, Chain: nil}}},
						}},
					{Target: "parameters.Properties.LinuxProfile", Name: validation.Null, Rule: true,
						Chain: []validation.Constraint{{Target: "parameters.Properties.LinuxProfile.AdminUsername", Name: validation.Null, Rule: true,
							Chain: []validation.Constraint{{Target: "parameters.Properties.LinuxProfile.AdminUsername", Name: validation.Pattern, Rule: `^[a-z][a-z0-9_-]*$`, Chain: nil}}},
							{Target: "parameters.Properties.LinuxProfile.SSH", Name: validation.Null, Rule: true,
								Chain: []validation.Constraint{{Target: "parameters.Properties.LinuxProfile.SSH.PublicKeys", Name: validation.Null, Rule: true, Chain: nil}}},
						}},
					{Target: "parameters.Properties.DiagnosticsProfile", Name: validation.Null, Rule: false,
						Chain: []validation.Constraint{{Target: "parameters.Properties.DiagnosticsProfile.VMDiagnostics", Name: validation.Null, Rule: true,
							Chain: []validation.Constraint{{Target: "parameters.Properties.DiagnosticsProfile.VMDiagnostics.Enabled", Name: validation.Null, Rule: true, Chain: nil}}},
						}},
				}}}}}); err != nil {
		return result, validation.NewErrorWithValidationError(err, "containerservice.ContainerServicesClient", "CreateOrUpdate")
	}

	req, err := client.CreateOrUpdatePreparer(resourceGroupName, containerServiceName, parameters, cancel)
	if err != nil {
		return result, autorest.NewErrorWithError(err, "containerservice.ContainerServicesClient", "CreateOrUpdate", nil, "Failure preparing request")
	}

	resp, err := client.CreateOrUpdateSender(req)
	if err != nil {
		result.Response = resp
		return result, autorest.NewErrorWithError(err, "containerservice.ContainerServicesClient", "CreateOrUpdate", resp, "Failure sending request")
	}

	result, err = client.CreateOrUpdateResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerservice.ContainerServicesClient", "CreateOrUpdate", resp, "Failure responding to request")
	}

	return
}

// CreateOrUpdatePreparer prepares the CreateOrUpdate request.
func (client ContainerServicesClient) CreateOrUpdatePreparer(resourceGroupName string, containerServiceName string, parameters ContainerService, cancel <-chan struct{}) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"containerServiceName": autorest.Encode("path", containerServiceName),
		"resourceGroupName":    autorest.Encode("path", resourceGroupName),
		"subscriptionId":       autorest.Encode("path", client.SubscriptionID),
	}

	queryParameters := map[string]interface{}{
		"api-version": client.APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsJSON(),
		autorest.AsPut(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ContainerService/containerServices/{containerServiceName}", pathParameters),
		autorest.WithJSON(parameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare(&http.Request{Cancel: cancel})
}

// CreateOrUpdateSender sends the CreateOrUpdate request. The method will close the
// http.Response Body if it receives an error.
func (client ContainerServicesClient) CreateOrUpdateSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client,
		req,
		azure.DoPollForAsynchronous(client.PollingDelay))
}

// CreateOrUpdateResponder handles the response to the CreateOrUpdate request. The method always
// closes the http.Response Body.
func (client ContainerServicesClient) CreateOrUpdateResponder(resp *http.Response) (result autorest.Response, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusCreated, http.StatusAccepted),
		autorest.ByClosing())
	result.Response = resp
	return
}

// Delete deletes the specified container service in the specified subscription
// and resource group. The operation does not delete other resources created as
// part of creating a container service, including storage accounts, VMs, and
// availability sets. All the other resources created with the container
// service are part of the same resource group and can be deleted individually.
// This method may poll for completion. Polling can be canceled by passing the
// cancel channel argument. The channel will be used to cancel polling and any
// outstanding HTTP requests.
//
// resourceGroupName is the name of the resource group. containerServiceName is
// the name of the container service in the specified subscription and resource
// group.
func (client ContainerServicesClient) Delete(resourceGroupName string, containerServiceName string, cancel <-chan struct{}) (result autorest.Response, err error) {
	req, err := client.DeletePreparer(resourceGroupName, containerServiceName, cancel)
	if err != nil {
		return result, autorest.NewErrorWithError(err, "containerservice.ContainerServicesClient", "Delete", nil, "Failure preparing request")
	}

	resp, err := client.DeleteSender(req)
	if err != nil {
		result.Response = resp
		return result, autorest.NewErrorWithError(err, "containerservice.ContainerServicesClient", "Delete", resp, "Failure sending request")
	}

	result, err = client.DeleteResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerservice.ContainerServicesClient", "Delete", resp, "Failure responding to request")
	}

	return
}

// DeletePreparer prepares the Delete request.
func (client ContainerServicesClient) DeletePreparer(resourceGroupName string, containerServiceName string, cancel <-chan struct{}) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"containerServiceName": autorest.Encode("path", containerServiceName),
		"resourceGroupName":    autorest.Encode("path", resourceGroupName),
		"subscriptionId":       autorest.Encode("path", client.SubscriptionID),
	}

	queryParameters := map[string]interface{}{
		"api-version": client.APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsDelete(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ContainerService/containerServices/{containerServiceName}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare(&http.Request{Cancel: cancel})
}

// DeleteSender sends the Delete request. The method will close the
// http.Response Body if it receives an error.
func (client ContainerServicesClient) DeleteSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client,
		req,
		azure.DoPollForAsynchronous(client.PollingDelay))
}

// DeleteResponder handles the response to the Delete request. The method always
// closes the http.Response Body.
func (client ContainerServicesClient) DeleteResponder(resp *http.Response) (result autorest.Response, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusAccepted, http.StatusNoContent),
		autorest.ByClosing())
	result.Response = resp
	return
}

// Get gets the properties of the specified container service in the specified
// subscription and resource group. The operation returns the properties
// including state, orchestrator, number of masters and agents, and FQDNs of
// masters and agents.
//
// resourceGroupName is the name of the resource group. containerServiceName is
// the name of the container service in the specified subscription and resource
// group.
func (client ContainerServicesClient) Get(resourceGroupName string, containerServiceName string) (result ContainerService, err error) {
	req, err := client.GetPreparer(resourceGroupName, containerServiceName)
	if err != nil {
		return result, autorest.NewErrorWithError(err, "containerservice.ContainerServicesClient", "Get", nil, "Failure preparing request")
	}

	resp, err := client.GetSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "containerservice.ContainerServicesClient", "Get", resp, "Failure sending request")
	}

	result, err = client.GetResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerservice.ContainerServicesClient", "Get", resp, "Failure responding to request")
	}

	return
}

// GetPreparer prepares the Get request.
func (client ContainerServicesClient) GetPreparer(resourceGroupName string, containerServiceName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"containerServiceName": autorest.Encode("path", containerServiceName),
		"resourceGroupName":    autorest.Encode("path", resourceGroupName),
		"subscriptionId":       autorest.Encode("path", client.SubscriptionID),
	}

	queryParameters := map[string]interface{}{
		"api-version": client.APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ContainerService/containerServices/{containerServiceName}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare(&http.Request{})
}

// GetSender sends the Get request. The method will close the
// http.Response Body if it receives an error.
func (client ContainerServicesClient) GetSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req)
}

// GetResponder handles the response to the Get request. The method always
// closes the http.Response Body.
func (client ContainerServicesClient) GetResponder(resp *http.Response) (result ContainerService, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// List gets a list of container services in the specified subscription. The
// operation returns properties of each container service including state,
// orchestrator, number of masters and agents, and FQDNs of masters and agents.
func (client ContainerServicesClient) List() (result ListResult, err error) {
	req, err := client.ListPreparer()
	if err != nil {
		return result, autorest.NewErrorWithError(err, "containerservice.ContainerServicesClient", "List", nil, "Failure preparing request")
	}

	resp, err := client.ListSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "containerservice.ContainerServicesClient", "List", resp, "Failure sending request")
	}

	result, err = client.ListResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerservice.ContainerServicesClient", "List", resp, "Failure responding to request")
	}

	return
}

// ListPreparer prepares the List request.
func (client ContainerServicesClient) ListPreparer() (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"subscriptionId": autorest.Encode("path", client.SubscriptionID),
	}

	queryParameters := map[string]interface{}{
		"api-version": client.APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/providers/Microsoft.ContainerService/containerServices", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare(&http.Request{})
}

// ListSender sends the List request. The method will close the
// http.Response Body if it receives an error.
func (client ContainerServicesClient) ListSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req)
}

// ListResponder handles the response to the List request. The method always
// closes the http.Response Body.
func (client ContainerServicesClient) ListResponder(resp *http.Response) (result ListResult, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// ListNextResults retrieves the next set of results, if any.
func (client ContainerServicesClient) ListNextResults(lastResults ListResult) (result ListResult, err error) {
	req, err := lastResults.ListResultPreparer()
	if err != nil {
		return result, autorest.NewErrorWithError(err, "containerservice.ContainerServicesClient", "List", nil, "Failure preparing next results request")
	}
	if req == nil {
		return
	}

	resp, err := client.ListSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "containerservice.ContainerServicesClient", "List", resp, "Failure sending next results request")
	}

	result, err = client.ListResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerservice.ContainerServicesClient", "List", resp, "Failure responding to next results request")
	}

	return
}

// ListByResourceGroup gets a list of container services in the specified
// subscription and resource group. The operation returns properties of each
// container service including state, orchestrator, number of masters and
// agents, and FQDNs of masters and agents.
//
// resourceGroupName is the name of the resource group.
func (client ContainerServicesClient) ListByResourceGroup(resourceGroupName string) (result ListResult, err error) {
	req, err := client.ListByResourceGroupPreparer(resourceGroupName)
	if err != nil {
		return result, autorest.NewErrorWithError(err, "containerservice.ContainerServicesClient", "ListByResourceGroup", nil, "Failure preparing request")
	}

	resp, err := client.ListByResourceGroupSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "containerservice.ContainerServicesClient", "ListByResourceGroup", resp, "Failure sending request")
	}

	result, err = client.ListByResourceGroupResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerservice.ContainerServicesClient", "ListByResourceGroup", resp, "Failure responding to request")
	}

	return
}

// ListByResourceGroupPreparer prepares the ListByResourceGroup request.
func (client ContainerServicesClient) ListByResourceGroupPreparer(resourceGroupName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	queryParameters := map[string]interface{}{
		"api-version": client.APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ContainerService/containerServices", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare(&http.Request{})
}

// ListByResourceGroupSender sends the ListByResourceGroup request. The method will close the
// http.Response Body if it receives an error.
func (client ContainerServicesClient) ListByResourceGroupSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req)
}

// ListByResourceGroupResponder handles the response to the ListByResourceGroup request. The method always
// closes the http.Response Body.
func (client ContainerServicesClient) ListByResourceGroupResponder(resp *http.Response) (result ListResult, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// ListByResourceGroupNextResults retrieves the next set of results, if any.
func (client ContainerServicesClient) ListByResourceGroupNextResults(lastResults ListResult) (result ListResult, err error) {
	req, err := lastResults.ListResultPreparer()
	if err != nil {
		return result, autorest.NewErrorWithError(err, "containerservice.ContainerServicesClient", "ListByResourceGroup", nil, "Failure preparing next results request")
	}
	if req == nil {
		return
	}

	resp, err := client.ListByResourceGroupSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "containerservice.ContainerServicesClient", "ListByResourceGroup", resp, "Failure sending next results request")
	}

	result, err = client.ListByResourceGroupResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerservice.ContainerServicesClient", "ListByResourceGroup", resp, "Failure responding to next results request")
	}

	return
}
