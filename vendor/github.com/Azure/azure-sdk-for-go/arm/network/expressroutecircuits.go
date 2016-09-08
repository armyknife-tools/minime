package network

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
// Code generated by Microsoft (R) AutoRest Code Generator 0.17.0.0
// Changes may cause incorrect behavior and will be lost if the code is
// regenerated.

import (
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"net/http"
)

// ExpressRouteCircuitsClient is the the Microsoft Azure Network management
// API provides a RESTful set of web services that interact with Microsoft
// Azure Networks service to manage your network resrources. The API has
// entities that capture the relationship between an end user and the
// Microsoft Azure Networks service.
type ExpressRouteCircuitsClient struct {
	ManagementClient
}

// NewExpressRouteCircuitsClient creates an instance of the
// ExpressRouteCircuitsClient client.
func NewExpressRouteCircuitsClient(subscriptionID string) ExpressRouteCircuitsClient {
	return NewExpressRouteCircuitsClientWithBaseURI(DefaultBaseURI, subscriptionID)
}

// NewExpressRouteCircuitsClientWithBaseURI creates an instance of the
// ExpressRouteCircuitsClient client.
func NewExpressRouteCircuitsClientWithBaseURI(baseURI string, subscriptionID string) ExpressRouteCircuitsClient {
	return ExpressRouteCircuitsClient{NewWithBaseURI(baseURI, subscriptionID)}
}

// CreateOrUpdate the Put ExpressRouteCircuit operation creates/updates a
// ExpressRouteCircuit This method may poll for completion. Polling can be
// canceled by passing the cancel channel argument. The channel will be used
// to cancel polling and any outstanding HTTP requests.
//
// resourceGroupName is the name of the resource group. circuitName is the
// name of the circuit. parameters is parameters supplied to the
// create/delete ExpressRouteCircuit operation
func (client ExpressRouteCircuitsClient) CreateOrUpdate(resourceGroupName string, circuitName string, parameters ExpressRouteCircuit, cancel <-chan struct{}) (result autorest.Response, err error) {
	req, err := client.CreateOrUpdatePreparer(resourceGroupName, circuitName, parameters, cancel)
	if err != nil {
		return result, autorest.NewErrorWithError(err, "network.ExpressRouteCircuitsClient", "CreateOrUpdate", nil, "Failure preparing request")
	}

	resp, err := client.CreateOrUpdateSender(req)
	if err != nil {
		result.Response = resp
		return result, autorest.NewErrorWithError(err, "network.ExpressRouteCircuitsClient", "CreateOrUpdate", resp, "Failure sending request")
	}

	result, err = client.CreateOrUpdateResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "network.ExpressRouteCircuitsClient", "CreateOrUpdate", resp, "Failure responding to request")
	}

	return
}

// CreateOrUpdatePreparer prepares the CreateOrUpdate request.
func (client ExpressRouteCircuitsClient) CreateOrUpdatePreparer(resourceGroupName string, circuitName string, parameters ExpressRouteCircuit, cancel <-chan struct{}) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"circuitName":       autorest.Encode("path", circuitName),
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
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/expressRouteCircuits/{circuitName}", pathParameters),
		autorest.WithJSON(parameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare(&http.Request{Cancel: cancel})
}

// CreateOrUpdateSender sends the CreateOrUpdate request. The method will close the
// http.Response Body if it receives an error.
func (client ExpressRouteCircuitsClient) CreateOrUpdateSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client,
		req,
		azure.DoPollForAsynchronous(client.PollingDelay))
}

// CreateOrUpdateResponder handles the response to the CreateOrUpdate request. The method always
// closes the http.Response Body.
func (client ExpressRouteCircuitsClient) CreateOrUpdateResponder(resp *http.Response) (result autorest.Response, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusCreated, http.StatusOK),
		autorest.ByClosing())
	result.Response = resp
	return
}

// Delete the delete ExpressRouteCircuit operation deletes the specified
// ExpressRouteCircuit. This method may poll for completion. Polling can be
// canceled by passing the cancel channel argument. The channel will be used
// to cancel polling and any outstanding HTTP requests.
//
// resourceGroupName is the name of the resource group. circuitName is the
// name of the express route Circuit.
func (client ExpressRouteCircuitsClient) Delete(resourceGroupName string, circuitName string, cancel <-chan struct{}) (result autorest.Response, err error) {
	req, err := client.DeletePreparer(resourceGroupName, circuitName, cancel)
	if err != nil {
		return result, autorest.NewErrorWithError(err, "network.ExpressRouteCircuitsClient", "Delete", nil, "Failure preparing request")
	}

	resp, err := client.DeleteSender(req)
	if err != nil {
		result.Response = resp
		return result, autorest.NewErrorWithError(err, "network.ExpressRouteCircuitsClient", "Delete", resp, "Failure sending request")
	}

	result, err = client.DeleteResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "network.ExpressRouteCircuitsClient", "Delete", resp, "Failure responding to request")
	}

	return
}

// DeletePreparer prepares the Delete request.
func (client ExpressRouteCircuitsClient) DeletePreparer(resourceGroupName string, circuitName string, cancel <-chan struct{}) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"circuitName":       autorest.Encode("path", circuitName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	queryParameters := map[string]interface{}{
		"api-version": client.APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsDelete(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/expressRouteCircuits/{circuitName}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare(&http.Request{Cancel: cancel})
}

// DeleteSender sends the Delete request. The method will close the
// http.Response Body if it receives an error.
func (client ExpressRouteCircuitsClient) DeleteSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client,
		req,
		azure.DoPollForAsynchronous(client.PollingDelay))
}

// DeleteResponder handles the response to the Delete request. The method always
// closes the http.Response Body.
func (client ExpressRouteCircuitsClient) DeleteResponder(resp *http.Response) (result autorest.Response, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusNoContent, http.StatusAccepted, http.StatusOK),
		autorest.ByClosing())
	result.Response = resp
	return
}

// Get the Get ExpressRouteCircuit operation retreives information about the
// specified ExpressRouteCircuit.
//
// resourceGroupName is the name of the resource group. circuitName is the
// name of the circuit.
func (client ExpressRouteCircuitsClient) Get(resourceGroupName string, circuitName string) (result ExpressRouteCircuit, err error) {
	req, err := client.GetPreparer(resourceGroupName, circuitName)
	if err != nil {
		return result, autorest.NewErrorWithError(err, "network.ExpressRouteCircuitsClient", "Get", nil, "Failure preparing request")
	}

	resp, err := client.GetSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "network.ExpressRouteCircuitsClient", "Get", resp, "Failure sending request")
	}

	result, err = client.GetResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "network.ExpressRouteCircuitsClient", "Get", resp, "Failure responding to request")
	}

	return
}

// GetPreparer prepares the Get request.
func (client ExpressRouteCircuitsClient) GetPreparer(resourceGroupName string, circuitName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"circuitName":       autorest.Encode("path", circuitName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	queryParameters := map[string]interface{}{
		"api-version": client.APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/expressRouteCircuits/{circuitName}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare(&http.Request{})
}

// GetSender sends the Get request. The method will close the
// http.Response Body if it receives an error.
func (client ExpressRouteCircuitsClient) GetSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req)
}

// GetResponder handles the response to the Get request. The method always
// closes the http.Response Body.
func (client ExpressRouteCircuitsClient) GetResponder(resp *http.Response) (result ExpressRouteCircuit, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// GetPeeringStats the Liststats ExpressRouteCircuit operation retrieves all
// the stats from a ExpressRouteCircuits in a resource group.
//
// resourceGroupName is the name of the resource group. circuitName is the
// name of the circuit. peeringName is the name of the peering.
func (client ExpressRouteCircuitsClient) GetPeeringStats(resourceGroupName string, circuitName string, peeringName string) (result ExpressRouteCircuitStats, err error) {
	req, err := client.GetPeeringStatsPreparer(resourceGroupName, circuitName, peeringName)
	if err != nil {
		return result, autorest.NewErrorWithError(err, "network.ExpressRouteCircuitsClient", "GetPeeringStats", nil, "Failure preparing request")
	}

	resp, err := client.GetPeeringStatsSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "network.ExpressRouteCircuitsClient", "GetPeeringStats", resp, "Failure sending request")
	}

	result, err = client.GetPeeringStatsResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "network.ExpressRouteCircuitsClient", "GetPeeringStats", resp, "Failure responding to request")
	}

	return
}

// GetPeeringStatsPreparer prepares the GetPeeringStats request.
func (client ExpressRouteCircuitsClient) GetPeeringStatsPreparer(resourceGroupName string, circuitName string, peeringName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"circuitName":       autorest.Encode("path", circuitName),
		"peeringName":       autorest.Encode("path", peeringName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	queryParameters := map[string]interface{}{
		"api-version": client.APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/expressRouteCircuits/{circuitName}/peerings/{peeringName}/stats", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare(&http.Request{})
}

// GetPeeringStatsSender sends the GetPeeringStats request. The method will close the
// http.Response Body if it receives an error.
func (client ExpressRouteCircuitsClient) GetPeeringStatsSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req)
}

// GetPeeringStatsResponder handles the response to the GetPeeringStats request. The method always
// closes the http.Response Body.
func (client ExpressRouteCircuitsClient) GetPeeringStatsResponder(resp *http.Response) (result ExpressRouteCircuitStats, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// GetStats the Liststats ExpressRouteCircuit operation retrieves all the
// stats from a ExpressRouteCircuits in a resource group.
//
// resourceGroupName is the name of the resource group. circuitName is the
// name of the circuit.
func (client ExpressRouteCircuitsClient) GetStats(resourceGroupName string, circuitName string) (result ExpressRouteCircuitStats, err error) {
	req, err := client.GetStatsPreparer(resourceGroupName, circuitName)
	if err != nil {
		return result, autorest.NewErrorWithError(err, "network.ExpressRouteCircuitsClient", "GetStats", nil, "Failure preparing request")
	}

	resp, err := client.GetStatsSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "network.ExpressRouteCircuitsClient", "GetStats", resp, "Failure sending request")
	}

	result, err = client.GetStatsResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "network.ExpressRouteCircuitsClient", "GetStats", resp, "Failure responding to request")
	}

	return
}

// GetStatsPreparer prepares the GetStats request.
func (client ExpressRouteCircuitsClient) GetStatsPreparer(resourceGroupName string, circuitName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"circuitName":       autorest.Encode("path", circuitName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	queryParameters := map[string]interface{}{
		"api-version": client.APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/expressRouteCircuits/{circuitName}/stats", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare(&http.Request{})
}

// GetStatsSender sends the GetStats request. The method will close the
// http.Response Body if it receives an error.
func (client ExpressRouteCircuitsClient) GetStatsSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req)
}

// GetStatsResponder handles the response to the GetStats request. The method always
// closes the http.Response Body.
func (client ExpressRouteCircuitsClient) GetStatsResponder(resp *http.Response) (result ExpressRouteCircuitStats, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// List the List ExpressRouteCircuit operation retrieves all the
// ExpressRouteCircuits in a resource group.
//
// resourceGroupName is the name of the resource group.
func (client ExpressRouteCircuitsClient) List(resourceGroupName string) (result ExpressRouteCircuitListResult, err error) {
	req, err := client.ListPreparer(resourceGroupName)
	if err != nil {
		return result, autorest.NewErrorWithError(err, "network.ExpressRouteCircuitsClient", "List", nil, "Failure preparing request")
	}

	resp, err := client.ListSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "network.ExpressRouteCircuitsClient", "List", resp, "Failure sending request")
	}

	result, err = client.ListResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "network.ExpressRouteCircuitsClient", "List", resp, "Failure responding to request")
	}

	return
}

// ListPreparer prepares the List request.
func (client ExpressRouteCircuitsClient) ListPreparer(resourceGroupName string) (*http.Request, error) {
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
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/expressRouteCircuits", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare(&http.Request{})
}

// ListSender sends the List request. The method will close the
// http.Response Body if it receives an error.
func (client ExpressRouteCircuitsClient) ListSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req)
}

// ListResponder handles the response to the List request. The method always
// closes the http.Response Body.
func (client ExpressRouteCircuitsClient) ListResponder(resp *http.Response) (result ExpressRouteCircuitListResult, err error) {
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
func (client ExpressRouteCircuitsClient) ListNextResults(lastResults ExpressRouteCircuitListResult) (result ExpressRouteCircuitListResult, err error) {
	req, err := lastResults.ExpressRouteCircuitListResultPreparer()
	if err != nil {
		return result, autorest.NewErrorWithError(err, "network.ExpressRouteCircuitsClient", "List", nil, "Failure preparing next results request request")
	}
	if req == nil {
		return
	}

	resp, err := client.ListSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "network.ExpressRouteCircuitsClient", "List", resp, "Failure sending next results request request")
	}

	result, err = client.ListResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "network.ExpressRouteCircuitsClient", "List", resp, "Failure responding to next results request request")
	}

	return
}

// ListAll the List ExpressRouteCircuit operation retrieves all the
// ExpressRouteCircuits in a subscription.
func (client ExpressRouteCircuitsClient) ListAll() (result ExpressRouteCircuitListResult, err error) {
	req, err := client.ListAllPreparer()
	if err != nil {
		return result, autorest.NewErrorWithError(err, "network.ExpressRouteCircuitsClient", "ListAll", nil, "Failure preparing request")
	}

	resp, err := client.ListAllSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "network.ExpressRouteCircuitsClient", "ListAll", resp, "Failure sending request")
	}

	result, err = client.ListAllResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "network.ExpressRouteCircuitsClient", "ListAll", resp, "Failure responding to request")
	}

	return
}

// ListAllPreparer prepares the ListAll request.
func (client ExpressRouteCircuitsClient) ListAllPreparer() (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"subscriptionId": autorest.Encode("path", client.SubscriptionID),
	}

	queryParameters := map[string]interface{}{
		"api-version": client.APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/providers/Microsoft.Network/expressRouteCircuits", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare(&http.Request{})
}

// ListAllSender sends the ListAll request. The method will close the
// http.Response Body if it receives an error.
func (client ExpressRouteCircuitsClient) ListAllSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req)
}

// ListAllResponder handles the response to the ListAll request. The method always
// closes the http.Response Body.
func (client ExpressRouteCircuitsClient) ListAllResponder(resp *http.Response) (result ExpressRouteCircuitListResult, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// ListAllNextResults retrieves the next set of results, if any.
func (client ExpressRouteCircuitsClient) ListAllNextResults(lastResults ExpressRouteCircuitListResult) (result ExpressRouteCircuitListResult, err error) {
	req, err := lastResults.ExpressRouteCircuitListResultPreparer()
	if err != nil {
		return result, autorest.NewErrorWithError(err, "network.ExpressRouteCircuitsClient", "ListAll", nil, "Failure preparing next results request request")
	}
	if req == nil {
		return
	}

	resp, err := client.ListAllSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "network.ExpressRouteCircuitsClient", "ListAll", resp, "Failure sending next results request request")
	}

	result, err = client.ListAllResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "network.ExpressRouteCircuitsClient", "ListAll", resp, "Failure responding to next results request request")
	}

	return
}

// ListArpTable the ListArpTable from ExpressRouteCircuit operation retrieves
// the currently advertised arp table associated with the
// ExpressRouteCircuits in a resource group. This method may poll for
// completion. Polling can be canceled by passing the cancel channel
// argument. The channel will be used to cancel polling and any outstanding
// HTTP requests.
//
// resourceGroupName is the name of the resource group. circuitName is the
// name of the circuit. peeringName is the name of the peering. devicePath is
// the path of the device.
func (client ExpressRouteCircuitsClient) ListArpTable(resourceGroupName string, circuitName string, peeringName string, devicePath string, cancel <-chan struct{}) (result autorest.Response, err error) {
	req, err := client.ListArpTablePreparer(resourceGroupName, circuitName, peeringName, devicePath, cancel)
	if err != nil {
		return result, autorest.NewErrorWithError(err, "network.ExpressRouteCircuitsClient", "ListArpTable", nil, "Failure preparing request")
	}

	resp, err := client.ListArpTableSender(req)
	if err != nil {
		result.Response = resp
		return result, autorest.NewErrorWithError(err, "network.ExpressRouteCircuitsClient", "ListArpTable", resp, "Failure sending request")
	}

	result, err = client.ListArpTableResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "network.ExpressRouteCircuitsClient", "ListArpTable", resp, "Failure responding to request")
	}

	return
}

// ListArpTablePreparer prepares the ListArpTable request.
func (client ExpressRouteCircuitsClient) ListArpTablePreparer(resourceGroupName string, circuitName string, peeringName string, devicePath string, cancel <-chan struct{}) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"circuitName":       autorest.Encode("path", circuitName),
		"devicePath":        autorest.Encode("path", devicePath),
		"peeringName":       autorest.Encode("path", peeringName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	queryParameters := map[string]interface{}{
		"api-version": client.APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsPost(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/expressRouteCircuits/{circuitName}/peerings/{peeringName}/arpTables/{devicePath}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare(&http.Request{Cancel: cancel})
}

// ListArpTableSender sends the ListArpTable request. The method will close the
// http.Response Body if it receives an error.
func (client ExpressRouteCircuitsClient) ListArpTableSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client,
		req,
		azure.DoPollForAsynchronous(client.PollingDelay))
}

// ListArpTableResponder handles the response to the ListArpTable request. The method always
// closes the http.Response Body.
func (client ExpressRouteCircuitsClient) ListArpTableResponder(resp *http.Response) (result autorest.Response, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusAccepted),
		autorest.ByClosing())
	result.Response = resp
	return
}

// ListRoutesTable the ListRoutesTable from ExpressRouteCircuit operation
// retrieves the currently advertised routes table associated with the
// ExpressRouteCircuits in a resource group. This method may poll for
// completion. Polling can be canceled by passing the cancel channel
// argument. The channel will be used to cancel polling and any outstanding
// HTTP requests.
//
// resourceGroupName is the name of the resource group. circuitName is the
// name of the circuit. peeringName is the name of the peering. devicePath is
// the path of the device.
func (client ExpressRouteCircuitsClient) ListRoutesTable(resourceGroupName string, circuitName string, peeringName string, devicePath string, cancel <-chan struct{}) (result autorest.Response, err error) {
	req, err := client.ListRoutesTablePreparer(resourceGroupName, circuitName, peeringName, devicePath, cancel)
	if err != nil {
		return result, autorest.NewErrorWithError(err, "network.ExpressRouteCircuitsClient", "ListRoutesTable", nil, "Failure preparing request")
	}

	resp, err := client.ListRoutesTableSender(req)
	if err != nil {
		result.Response = resp
		return result, autorest.NewErrorWithError(err, "network.ExpressRouteCircuitsClient", "ListRoutesTable", resp, "Failure sending request")
	}

	result, err = client.ListRoutesTableResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "network.ExpressRouteCircuitsClient", "ListRoutesTable", resp, "Failure responding to request")
	}

	return
}

// ListRoutesTablePreparer prepares the ListRoutesTable request.
func (client ExpressRouteCircuitsClient) ListRoutesTablePreparer(resourceGroupName string, circuitName string, peeringName string, devicePath string, cancel <-chan struct{}) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"circuitName":       autorest.Encode("path", circuitName),
		"devicePath":        autorest.Encode("path", devicePath),
		"peeringName":       autorest.Encode("path", peeringName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	queryParameters := map[string]interface{}{
		"api-version": client.APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsPost(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/expressRouteCircuits/{circuitName}/peerings/{peeringName}/routeTables/{devicePath}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare(&http.Request{Cancel: cancel})
}

// ListRoutesTableSender sends the ListRoutesTable request. The method will close the
// http.Response Body if it receives an error.
func (client ExpressRouteCircuitsClient) ListRoutesTableSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client,
		req,
		azure.DoPollForAsynchronous(client.PollingDelay))
}

// ListRoutesTableResponder handles the response to the ListRoutesTable request. The method always
// closes the http.Response Body.
func (client ExpressRouteCircuitsClient) ListRoutesTableResponder(resp *http.Response) (result autorest.Response, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusAccepted),
		autorest.ByClosing())
	result.Response = resp
	return
}

// ListRoutesTableSummary the ListRoutesTable from ExpressRouteCircuit
// operation retrieves the currently advertised routes table associated with
// the ExpressRouteCircuits in a resource group. This method may poll for
// completion. Polling can be canceled by passing the cancel channel
// argument. The channel will be used to cancel polling and any outstanding
// HTTP requests.
//
// resourceGroupName is the name of the resource group. circuitName is the
// name of the circuit. peeringName is the name of the peering. devicePath is
// the path of the device.
func (client ExpressRouteCircuitsClient) ListRoutesTableSummary(resourceGroupName string, circuitName string, peeringName string, devicePath string, cancel <-chan struct{}) (result autorest.Response, err error) {
	req, err := client.ListRoutesTableSummaryPreparer(resourceGroupName, circuitName, peeringName, devicePath, cancel)
	if err != nil {
		return result, autorest.NewErrorWithError(err, "network.ExpressRouteCircuitsClient", "ListRoutesTableSummary", nil, "Failure preparing request")
	}

	resp, err := client.ListRoutesTableSummarySender(req)
	if err != nil {
		result.Response = resp
		return result, autorest.NewErrorWithError(err, "network.ExpressRouteCircuitsClient", "ListRoutesTableSummary", resp, "Failure sending request")
	}

	result, err = client.ListRoutesTableSummaryResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "network.ExpressRouteCircuitsClient", "ListRoutesTableSummary", resp, "Failure responding to request")
	}

	return
}

// ListRoutesTableSummaryPreparer prepares the ListRoutesTableSummary request.
func (client ExpressRouteCircuitsClient) ListRoutesTableSummaryPreparer(resourceGroupName string, circuitName string, peeringName string, devicePath string, cancel <-chan struct{}) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"circuitName":       autorest.Encode("path", circuitName),
		"devicePath":        autorest.Encode("path", devicePath),
		"peeringName":       autorest.Encode("path", peeringName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	queryParameters := map[string]interface{}{
		"api-version": client.APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsPost(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/expressRouteCircuits/{circuitName}/peerings/{peeringName}/routeTablesSummary/{devicePath}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare(&http.Request{Cancel: cancel})
}

// ListRoutesTableSummarySender sends the ListRoutesTableSummary request. The method will close the
// http.Response Body if it receives an error.
func (client ExpressRouteCircuitsClient) ListRoutesTableSummarySender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client,
		req,
		azure.DoPollForAsynchronous(client.PollingDelay))
}

// ListRoutesTableSummaryResponder handles the response to the ListRoutesTableSummary request. The method always
// closes the http.Response Body.
func (client ExpressRouteCircuitsClient) ListRoutesTableSummaryResponder(resp *http.Response) (result autorest.Response, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusAccepted),
		autorest.ByClosing())
	result.Response = resp
	return
}
