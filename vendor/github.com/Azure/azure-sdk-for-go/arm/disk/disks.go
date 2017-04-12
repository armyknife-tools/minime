package disk

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
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/validation"
	"net/http"
)

// DisksClient is the the Disk Resource Provider Client.
type DisksClient struct {
	ManagementClient
}

// NewDisksClient creates an instance of the DisksClient client.
func NewDisksClient(subscriptionID string) DisksClient {
	return NewDisksClientWithBaseURI(DefaultBaseURI, subscriptionID)
}

// NewDisksClientWithBaseURI creates an instance of the DisksClient client.
func NewDisksClientWithBaseURI(baseURI string, subscriptionID string) DisksClient {
	return DisksClient{NewWithBaseURI(baseURI, subscriptionID)}
}

// CreateOrUpdate creates or updates a disk. This method may poll for
// completion. Polling can be canceled by passing the cancel channel argument.
// The channel will be used to cancel polling and any outstanding HTTP
// requests.
//
// resourceGroupName is the name of the resource group. diskName is the name of
// the disk within the given subscription and resource group. diskParameter is
// disk object supplied in the body of the Put disk operation.
func (client DisksClient) CreateOrUpdate(resourceGroupName string, diskName string, diskParameter Model, cancel <-chan struct{}) (result autorest.Response, err error) {
	if err := validation.Validate([]validation.Validation{
		{TargetValue: diskParameter,
			Constraints: []validation.Constraint{{Target: "diskParameter.Properties", Name: validation.Null, Rule: false,
				Chain: []validation.Constraint{{Target: "diskParameter.Properties.CreationData", Name: validation.Null, Rule: true,
					Chain: []validation.Constraint{{Target: "diskParameter.Properties.CreationData.ImageReference", Name: validation.Null, Rule: false,
						Chain: []validation.Constraint{{Target: "diskParameter.Properties.CreationData.ImageReference.ID", Name: validation.Null, Rule: true, Chain: nil}}},
					}},
					{Target: "diskParameter.Properties.EncryptionSettings", Name: validation.Null, Rule: false,
						Chain: []validation.Constraint{{Target: "diskParameter.Properties.EncryptionSettings.DiskEncryptionKey", Name: validation.Null, Rule: false,
							Chain: []validation.Constraint{{Target: "diskParameter.Properties.EncryptionSettings.DiskEncryptionKey.SourceVault", Name: validation.Null, Rule: true, Chain: nil},
								{Target: "diskParameter.Properties.EncryptionSettings.DiskEncryptionKey.SecretURL", Name: validation.Null, Rule: true, Chain: nil},
							}},
							{Target: "diskParameter.Properties.EncryptionSettings.KeyEncryptionKey", Name: validation.Null, Rule: false,
								Chain: []validation.Constraint{{Target: "diskParameter.Properties.EncryptionSettings.KeyEncryptionKey.SourceVault", Name: validation.Null, Rule: true, Chain: nil},
									{Target: "diskParameter.Properties.EncryptionSettings.KeyEncryptionKey.KeyURL", Name: validation.Null, Rule: true, Chain: nil},
								}},
						}},
				}}}}}); err != nil {
		return result, validation.NewErrorWithValidationError(err, "disk.DisksClient", "CreateOrUpdate")
	}

	req, err := client.CreateOrUpdatePreparer(resourceGroupName, diskName, diskParameter, cancel)
	if err != nil {
		return result, autorest.NewErrorWithError(err, "disk.DisksClient", "CreateOrUpdate", nil, "Failure preparing request")
	}

	resp, err := client.CreateOrUpdateSender(req)
	if err != nil {
		result.Response = resp
		return result, autorest.NewErrorWithError(err, "disk.DisksClient", "CreateOrUpdate", resp, "Failure sending request")
	}

	result, err = client.CreateOrUpdateResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "disk.DisksClient", "CreateOrUpdate", resp, "Failure responding to request")
	}

	return
}

// CreateOrUpdatePreparer prepares the CreateOrUpdate request.
func (client DisksClient) CreateOrUpdatePreparer(resourceGroupName string, diskName string, diskParameter Model, cancel <-chan struct{}) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"diskName":          autorest.Encode("path", diskName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	queryParameters := map[string]interface{}{
		"api-version": client.APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsJSON(),
		autorest.AsPut(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Compute/disks/{diskName}", pathParameters),
		autorest.WithJSON(diskParameter),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare(&http.Request{Cancel: cancel})
}

// CreateOrUpdateSender sends the CreateOrUpdate request. The method will close the
// http.Response Body if it receives an error.
func (client DisksClient) CreateOrUpdateSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client,
		req,
		azure.DoPollForAsynchronous(client.PollingDelay))
}

// CreateOrUpdateResponder handles the response to the CreateOrUpdate request. The method always
// closes the http.Response Body.
func (client DisksClient) CreateOrUpdateResponder(resp *http.Response) (result autorest.Response, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusAccepted),
		autorest.ByClosing())
	result.Response = resp
	return
}

// Delete deletes a disk. This method may poll for completion. Polling can be
// canceled by passing the cancel channel argument. The channel will be used to
// cancel polling and any outstanding HTTP requests.
//
// resourceGroupName is the name of the resource group. diskName is the name of
// the disk within the given subscription and resource group.
func (client DisksClient) Delete(resourceGroupName string, diskName string, cancel <-chan struct{}) (result autorest.Response, err error) {
	req, err := client.DeletePreparer(resourceGroupName, diskName, cancel)
	if err != nil {
		return result, autorest.NewErrorWithError(err, "disk.DisksClient", "Delete", nil, "Failure preparing request")
	}

	resp, err := client.DeleteSender(req)
	if err != nil {
		result.Response = resp
		return result, autorest.NewErrorWithError(err, "disk.DisksClient", "Delete", resp, "Failure sending request")
	}

	result, err = client.DeleteResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "disk.DisksClient", "Delete", resp, "Failure responding to request")
	}

	return
}

// DeletePreparer prepares the Delete request.
func (client DisksClient) DeletePreparer(resourceGroupName string, diskName string, cancel <-chan struct{}) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"diskName":          autorest.Encode("path", diskName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	queryParameters := map[string]interface{}{
		"api-version": client.APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsDelete(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Compute/disks/{diskName}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare(&http.Request{Cancel: cancel})
}

// DeleteSender sends the Delete request. The method will close the
// http.Response Body if it receives an error.
func (client DisksClient) DeleteSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client,
		req,
		azure.DoPollForAsynchronous(client.PollingDelay))
}

// DeleteResponder handles the response to the Delete request. The method always
// closes the http.Response Body.
func (client DisksClient) DeleteResponder(resp *http.Response) (result autorest.Response, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusAccepted, http.StatusNoContent),
		autorest.ByClosing())
	result.Response = resp
	return
}

// Get gets information about a disk.
//
// resourceGroupName is the name of the resource group. diskName is the name of
// the disk within the given subscription and resource group.
func (client DisksClient) Get(resourceGroupName string, diskName string) (result Model, err error) {
	req, err := client.GetPreparer(resourceGroupName, diskName)
	if err != nil {
		return result, autorest.NewErrorWithError(err, "disk.DisksClient", "Get", nil, "Failure preparing request")
	}

	resp, err := client.GetSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "disk.DisksClient", "Get", resp, "Failure sending request")
	}

	result, err = client.GetResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "disk.DisksClient", "Get", resp, "Failure responding to request")
	}

	return
}

// GetPreparer prepares the Get request.
func (client DisksClient) GetPreparer(resourceGroupName string, diskName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"diskName":          autorest.Encode("path", diskName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	queryParameters := map[string]interface{}{
		"api-version": client.APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Compute/disks/{diskName}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare(&http.Request{})
}

// GetSender sends the Get request. The method will close the
// http.Response Body if it receives an error.
func (client DisksClient) GetSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req)
}

// GetResponder handles the response to the Get request. The method always
// closes the http.Response Body.
func (client DisksClient) GetResponder(resp *http.Response) (result Model, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// GrantAccess grants access to a disk. This method may poll for completion.
// Polling can be canceled by passing the cancel channel argument. The channel
// will be used to cancel polling and any outstanding HTTP requests.
//
// resourceGroupName is the name of the resource group. diskName is the name of
// the disk within the given subscription and resource group. grantAccessData
// is access data object supplied in the body of the get disk access operation.
func (client DisksClient) GrantAccess(resourceGroupName string, diskName string, grantAccessData GrantAccessData, cancel <-chan struct{}) (result autorest.Response, err error) {
	if err := validation.Validate([]validation.Validation{
		{TargetValue: grantAccessData,
			Constraints: []validation.Constraint{{Target: "grantAccessData.DurationInSeconds", Name: validation.Null, Rule: true, Chain: nil}}}}); err != nil {
		return result, validation.NewErrorWithValidationError(err, "disk.DisksClient", "GrantAccess")
	}

	req, err := client.GrantAccessPreparer(resourceGroupName, diskName, grantAccessData, cancel)
	if err != nil {
		return result, autorest.NewErrorWithError(err, "disk.DisksClient", "GrantAccess", nil, "Failure preparing request")
	}

	resp, err := client.GrantAccessSender(req)
	if err != nil {
		result.Response = resp
		return result, autorest.NewErrorWithError(err, "disk.DisksClient", "GrantAccess", resp, "Failure sending request")
	}

	result, err = client.GrantAccessResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "disk.DisksClient", "GrantAccess", resp, "Failure responding to request")
	}

	return
}

// GrantAccessPreparer prepares the GrantAccess request.
func (client DisksClient) GrantAccessPreparer(resourceGroupName string, diskName string, grantAccessData GrantAccessData, cancel <-chan struct{}) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"diskName":          autorest.Encode("path", diskName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	queryParameters := map[string]interface{}{
		"api-version": client.APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsJSON(),
		autorest.AsPost(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Compute/disks/{diskName}/beginGetAccess", pathParameters),
		autorest.WithJSON(grantAccessData),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare(&http.Request{Cancel: cancel})
}

// GrantAccessSender sends the GrantAccess request. The method will close the
// http.Response Body if it receives an error.
func (client DisksClient) GrantAccessSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client,
		req,
		azure.DoPollForAsynchronous(client.PollingDelay))
}

// GrantAccessResponder handles the response to the GrantAccess request. The method always
// closes the http.Response Body.
func (client DisksClient) GrantAccessResponder(resp *http.Response) (result autorest.Response, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusAccepted),
		autorest.ByClosing())
	result.Response = resp
	return
}

// List lists all the disks under a subscription.
func (client DisksClient) List() (result ListType, err error) {
	req, err := client.ListPreparer()
	if err != nil {
		return result, autorest.NewErrorWithError(err, "disk.DisksClient", "List", nil, "Failure preparing request")
	}

	resp, err := client.ListSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "disk.DisksClient", "List", resp, "Failure sending request")
	}

	result, err = client.ListResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "disk.DisksClient", "List", resp, "Failure responding to request")
	}

	return
}

// ListPreparer prepares the List request.
func (client DisksClient) ListPreparer() (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"subscriptionId": autorest.Encode("path", client.SubscriptionID),
	}

	queryParameters := map[string]interface{}{
		"api-version": client.APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/providers/Microsoft.Compute/disks", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare(&http.Request{})
}

// ListSender sends the List request. The method will close the
// http.Response Body if it receives an error.
func (client DisksClient) ListSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req)
}

// ListResponder handles the response to the List request. The method always
// closes the http.Response Body.
func (client DisksClient) ListResponder(resp *http.Response) (result ListType, err error) {
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
func (client DisksClient) ListNextResults(lastResults ListType) (result ListType, err error) {
	req, err := lastResults.ListTypePreparer()
	if err != nil {
		return result, autorest.NewErrorWithError(err, "disk.DisksClient", "List", nil, "Failure preparing next results request")
	}
	if req == nil {
		return
	}

	resp, err := client.ListSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "disk.DisksClient", "List", resp, "Failure sending next results request")
	}

	result, err = client.ListResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "disk.DisksClient", "List", resp, "Failure responding to next results request")
	}

	return
}

// ListByResourceGroup lists all the disks under a resource group.
//
// resourceGroupName is the name of the resource group.
func (client DisksClient) ListByResourceGroup(resourceGroupName string) (result ListType, err error) {
	req, err := client.ListByResourceGroupPreparer(resourceGroupName)
	if err != nil {
		return result, autorest.NewErrorWithError(err, "disk.DisksClient", "ListByResourceGroup", nil, "Failure preparing request")
	}

	resp, err := client.ListByResourceGroupSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "disk.DisksClient", "ListByResourceGroup", resp, "Failure sending request")
	}

	result, err = client.ListByResourceGroupResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "disk.DisksClient", "ListByResourceGroup", resp, "Failure responding to request")
	}

	return
}

// ListByResourceGroupPreparer prepares the ListByResourceGroup request.
func (client DisksClient) ListByResourceGroupPreparer(resourceGroupName string) (*http.Request, error) {
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
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Compute/disks", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare(&http.Request{})
}

// ListByResourceGroupSender sends the ListByResourceGroup request. The method will close the
// http.Response Body if it receives an error.
func (client DisksClient) ListByResourceGroupSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req)
}

// ListByResourceGroupResponder handles the response to the ListByResourceGroup request. The method always
// closes the http.Response Body.
func (client DisksClient) ListByResourceGroupResponder(resp *http.Response) (result ListType, err error) {
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
func (client DisksClient) ListByResourceGroupNextResults(lastResults ListType) (result ListType, err error) {
	req, err := lastResults.ListTypePreparer()
	if err != nil {
		return result, autorest.NewErrorWithError(err, "disk.DisksClient", "ListByResourceGroup", nil, "Failure preparing next results request")
	}
	if req == nil {
		return
	}

	resp, err := client.ListByResourceGroupSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "disk.DisksClient", "ListByResourceGroup", resp, "Failure sending next results request")
	}

	result, err = client.ListByResourceGroupResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "disk.DisksClient", "ListByResourceGroup", resp, "Failure responding to next results request")
	}

	return
}

// RevokeAccess revokes access to a disk. This method may poll for completion.
// Polling can be canceled by passing the cancel channel argument. The channel
// will be used to cancel polling and any outstanding HTTP requests.
//
// resourceGroupName is the name of the resource group. diskName is the name of
// the disk within the given subscription and resource group.
func (client DisksClient) RevokeAccess(resourceGroupName string, diskName string, cancel <-chan struct{}) (result autorest.Response, err error) {
	req, err := client.RevokeAccessPreparer(resourceGroupName, diskName, cancel)
	if err != nil {
		return result, autorest.NewErrorWithError(err, "disk.DisksClient", "RevokeAccess", nil, "Failure preparing request")
	}

	resp, err := client.RevokeAccessSender(req)
	if err != nil {
		result.Response = resp
		return result, autorest.NewErrorWithError(err, "disk.DisksClient", "RevokeAccess", resp, "Failure sending request")
	}

	result, err = client.RevokeAccessResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "disk.DisksClient", "RevokeAccess", resp, "Failure responding to request")
	}

	return
}

// RevokeAccessPreparer prepares the RevokeAccess request.
func (client DisksClient) RevokeAccessPreparer(resourceGroupName string, diskName string, cancel <-chan struct{}) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"diskName":          autorest.Encode("path", diskName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	queryParameters := map[string]interface{}{
		"api-version": client.APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsPost(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Compute/disks/{diskName}/endGetAccess", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare(&http.Request{Cancel: cancel})
}

// RevokeAccessSender sends the RevokeAccess request. The method will close the
// http.Response Body if it receives an error.
func (client DisksClient) RevokeAccessSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client,
		req,
		azure.DoPollForAsynchronous(client.PollingDelay))
}

// RevokeAccessResponder handles the response to the RevokeAccess request. The method always
// closes the http.Response Body.
func (client DisksClient) RevokeAccessResponder(resp *http.Response) (result autorest.Response, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusAccepted),
		autorest.ByClosing())
	result.Response = resp
	return
}

// Update updates (patches) a disk. This method may poll for completion.
// Polling can be canceled by passing the cancel channel argument. The channel
// will be used to cancel polling and any outstanding HTTP requests.
//
// resourceGroupName is the name of the resource group. diskName is the name of
// the disk within the given subscription and resource group. diskParameter is
// disk object supplied in the body of the Patch disk operation.
func (client DisksClient) Update(resourceGroupName string, diskName string, diskParameter UpdateType, cancel <-chan struct{}) (result autorest.Response, err error) {
	req, err := client.UpdatePreparer(resourceGroupName, diskName, diskParameter, cancel)
	if err != nil {
		return result, autorest.NewErrorWithError(err, "disk.DisksClient", "Update", nil, "Failure preparing request")
	}

	resp, err := client.UpdateSender(req)
	if err != nil {
		result.Response = resp
		return result, autorest.NewErrorWithError(err, "disk.DisksClient", "Update", resp, "Failure sending request")
	}

	result, err = client.UpdateResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "disk.DisksClient", "Update", resp, "Failure responding to request")
	}

	return
}

// UpdatePreparer prepares the Update request.
func (client DisksClient) UpdatePreparer(resourceGroupName string, diskName string, diskParameter UpdateType, cancel <-chan struct{}) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"diskName":          autorest.Encode("path", diskName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	queryParameters := map[string]interface{}{
		"api-version": client.APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsJSON(),
		autorest.AsPatch(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Compute/disks/{diskName}", pathParameters),
		autorest.WithJSON(diskParameter),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare(&http.Request{Cancel: cancel})
}

// UpdateSender sends the Update request. The method will close the
// http.Response Body if it receives an error.
func (client DisksClient) UpdateSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client,
		req,
		azure.DoPollForAsynchronous(client.PollingDelay))
}

// UpdateResponder handles the response to the Update request. The method always
// closes the http.Response Body.
func (client DisksClient) UpdateResponder(resp *http.Response) (result autorest.Response, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusAccepted),
		autorest.ByClosing())
	result.Response = resp
	return
}
