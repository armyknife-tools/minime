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
// Code generated by Microsoft (R) AutoRest Code Generator 0.12.0.0
// Changes may cause incorrect behavior and will be lost if the code is
// regenerated.

import (
	"github.com/Azure/azure-sdk-for-go/Godeps/_workspace/src/github.com/Azure/go-autorest/autorest"
	"net/http"
	"net/url"
)

// RouteTablesClient is the the Microsoft Azure Network management API
// provides a RESTful set of web services that interact with Microsoft Azure
// Networks service to manage your network resrources. The API has entities
// that capture the relationship between an end user and the Microsoft Azure
// Networks service.
type RouteTablesClient struct {
	ManagementClient
}

// NewRouteTablesClient creates an instance of the RouteTablesClient client.
func NewRouteTablesClient(subscriptionID string) RouteTablesClient {
	return NewRouteTablesClientWithBaseURI(DefaultBaseURI, subscriptionID)
}

// NewRouteTablesClientWithBaseURI creates an instance of the
// RouteTablesClient client.
func NewRouteTablesClientWithBaseURI(baseURI string, subscriptionID string) RouteTablesClient {
	return RouteTablesClient{NewWithBaseURI(baseURI, subscriptionID)}
}

// CreateOrUpdate the Put RouteTable operation creates/updates a route tablein
// the specified resource group.
//
// resourceGroupName is the name of the resource group. routeTableName is the
// name of the route table. parameters is parameters supplied to the
// create/update Route Table operation
func (client RouteTablesClient) CreateOrUpdate(resourceGroupName string, routeTableName string, parameters RouteTable) (result RouteTable, ae error) {
	req, err := client.CreateOrUpdatePreparer(resourceGroupName, routeTableName, parameters)
	if err != nil {
		return result, autorest.NewErrorWithError(err, "network/RouteTablesClient", "CreateOrUpdate", autorest.UndefinedStatusCode, "Failure preparing request")
	}

	resp, err := client.CreateOrUpdateSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "network/RouteTablesClient", "CreateOrUpdate", resp.StatusCode, "Failure sending request")
	}

	result, err = client.CreateOrUpdateResponder(resp)
	if err != nil {
		ae = autorest.NewErrorWithError(err, "network/RouteTablesClient", "CreateOrUpdate", resp.StatusCode, "Failure responding to request")
	}

	return
}

// CreateOrUpdatePreparer prepares the CreateOrUpdate request.
func (client RouteTablesClient) CreateOrUpdatePreparer(resourceGroupName string, routeTableName string, parameters RouteTable) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"resourceGroupName": url.QueryEscape(resourceGroupName),
		"routeTableName":    url.QueryEscape(routeTableName),
		"subscriptionId":    url.QueryEscape(client.SubscriptionID),
	}

	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	return autorest.Prepare(&http.Request{},
		autorest.AsJSON(),
		autorest.AsPut(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPath("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/routeTables/{routeTableName}"),
		autorest.WithJSON(parameters),
		autorest.WithPathParameters(pathParameters),
		autorest.WithQueryParameters(queryParameters))
}

// CreateOrUpdateSender sends the CreateOrUpdate request. The method will close the
// http.Response Body if it receives an error.
func (client RouteTablesClient) CreateOrUpdateSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, http.StatusOK, http.StatusCreated)
}

// CreateOrUpdateResponder handles the response to the CreateOrUpdate request. The method always
// closes the http.Response Body.
func (client RouteTablesClient) CreateOrUpdateResponder(resp *http.Response) (result RouteTable, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		autorest.WithErrorUnlessStatusCode(http.StatusOK, http.StatusCreated),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// Delete the Delete RouteTable operation deletes the specifed Route Table
//
// resourceGroupName is the name of the resource group. routeTableName is the
// name of the route table.
func (client RouteTablesClient) Delete(resourceGroupName string, routeTableName string) (result autorest.Response, ae error) {
	req, err := client.DeletePreparer(resourceGroupName, routeTableName)
	if err != nil {
		return result, autorest.NewErrorWithError(err, "network/RouteTablesClient", "Delete", autorest.UndefinedStatusCode, "Failure preparing request")
	}

	resp, err := client.DeleteSender(req)
	if err != nil {
		result.Response = resp
		return result, autorest.NewErrorWithError(err, "network/RouteTablesClient", "Delete", resp.StatusCode, "Failure sending request")
	}

	result, err = client.DeleteResponder(resp)
	if err != nil {
		ae = autorest.NewErrorWithError(err, "network/RouteTablesClient", "Delete", resp.StatusCode, "Failure responding to request")
	}

	return
}

// DeletePreparer prepares the Delete request.
func (client RouteTablesClient) DeletePreparer(resourceGroupName string, routeTableName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"resourceGroupName": url.QueryEscape(resourceGroupName),
		"routeTableName":    url.QueryEscape(routeTableName),
		"subscriptionId":    url.QueryEscape(client.SubscriptionID),
	}

	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	return autorest.Prepare(&http.Request{},
		autorest.AsJSON(),
		autorest.AsDelete(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPath("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/routeTables/{routeTableName}"),
		autorest.WithPathParameters(pathParameters),
		autorest.WithQueryParameters(queryParameters))
}

// DeleteSender sends the Delete request. The method will close the
// http.Response Body if it receives an error.
func (client RouteTablesClient) DeleteSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, http.StatusNoContent, http.StatusOK, http.StatusAccepted)
}

// DeleteResponder handles the response to the Delete request. The method always
// closes the http.Response Body.
func (client RouteTablesClient) DeleteResponder(resp *http.Response) (result autorest.Response, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		autorest.WithErrorUnlessStatusCode(http.StatusNoContent, http.StatusOK, http.StatusAccepted),
		autorest.ByClosing())
	result.Response = resp
	return
}

// Get the Get RouteTables operation retrieves information about the specified
// route table.
//
// resourceGroupName is the name of the resource group. routeTableName is the
// name of the route table. expand is expand references resources.
func (client RouteTablesClient) Get(resourceGroupName string, routeTableName string, expand string) (result RouteTable, ae error) {
	req, err := client.GetPreparer(resourceGroupName, routeTableName, expand)
	if err != nil {
		return result, autorest.NewErrorWithError(err, "network/RouteTablesClient", "Get", autorest.UndefinedStatusCode, "Failure preparing request")
	}

	resp, err := client.GetSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "network/RouteTablesClient", "Get", resp.StatusCode, "Failure sending request")
	}

	result, err = client.GetResponder(resp)
	if err != nil {
		ae = autorest.NewErrorWithError(err, "network/RouteTablesClient", "Get", resp.StatusCode, "Failure responding to request")
	}

	return
}

// GetPreparer prepares the Get request.
func (client RouteTablesClient) GetPreparer(resourceGroupName string, routeTableName string, expand string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"resourceGroupName": url.QueryEscape(resourceGroupName),
		"routeTableName":    url.QueryEscape(routeTableName),
		"subscriptionId":    url.QueryEscape(client.SubscriptionID),
	}

	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}
	if len(expand) > 0 {
		queryParameters["$expand"] = expand
	}

	return autorest.Prepare(&http.Request{},
		autorest.AsJSON(),
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPath("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/routeTables/{routeTableName}"),
		autorest.WithPathParameters(pathParameters),
		autorest.WithQueryParameters(queryParameters))
}

// GetSender sends the Get request. The method will close the
// http.Response Body if it receives an error.
func (client RouteTablesClient) GetSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, http.StatusOK)
}

// GetResponder handles the response to the Get request. The method always
// closes the http.Response Body.
func (client RouteTablesClient) GetResponder(resp *http.Response) (result RouteTable, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		autorest.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// List the list RouteTables returns all route tables in a resource group
//
// resourceGroupName is the name of the resource group.
func (client RouteTablesClient) List(resourceGroupName string) (result RouteTableListResult, ae error) {
	req, err := client.ListPreparer(resourceGroupName)
	if err != nil {
		return result, autorest.NewErrorWithError(err, "network/RouteTablesClient", "List", autorest.UndefinedStatusCode, "Failure preparing request")
	}

	resp, err := client.ListSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "network/RouteTablesClient", "List", resp.StatusCode, "Failure sending request")
	}

	result, err = client.ListResponder(resp)
	if err != nil {
		ae = autorest.NewErrorWithError(err, "network/RouteTablesClient", "List", resp.StatusCode, "Failure responding to request")
	}

	return
}

// ListPreparer prepares the List request.
func (client RouteTablesClient) ListPreparer(resourceGroupName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"resourceGroupName": url.QueryEscape(resourceGroupName),
		"subscriptionId":    url.QueryEscape(client.SubscriptionID),
	}

	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	return autorest.Prepare(&http.Request{},
		autorest.AsJSON(),
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPath("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/routeTables"),
		autorest.WithPathParameters(pathParameters),
		autorest.WithQueryParameters(queryParameters))
}

// ListSender sends the List request. The method will close the
// http.Response Body if it receives an error.
func (client RouteTablesClient) ListSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, http.StatusOK)
}

// ListResponder handles the response to the List request. The method always
// closes the http.Response Body.
func (client RouteTablesClient) ListResponder(resp *http.Response) (result RouteTableListResult, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		autorest.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// ListNextResults retrieves the next set of results, if any.
func (client RouteTablesClient) ListNextResults(lastResults RouteTableListResult) (result RouteTableListResult, ae error) {
	req, err := lastResults.RouteTableListResultPreparer()
	if err != nil {
		return result, autorest.NewErrorWithError(err, "network/RouteTablesClient", "List", autorest.UndefinedStatusCode, "Failure preparing next results request request")
	}
	if req == nil {
		return
	}

	resp, err := client.ListSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "network/RouteTablesClient", "List", resp.StatusCode, "Failure sending next results request request")
	}

	result, err = client.ListResponder(resp)
	if err != nil {
		ae = autorest.NewErrorWithError(err, "network/RouteTablesClient", "List", resp.StatusCode, "Failure responding to next results request request")
	}

	return
}

// ListAll the list RouteTables returns all route tables in a subscription
func (client RouteTablesClient) ListAll() (result RouteTableListResult, ae error) {
	req, err := client.ListAllPreparer()
	if err != nil {
		return result, autorest.NewErrorWithError(err, "network/RouteTablesClient", "ListAll", autorest.UndefinedStatusCode, "Failure preparing request")
	}

	resp, err := client.ListAllSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "network/RouteTablesClient", "ListAll", resp.StatusCode, "Failure sending request")
	}

	result, err = client.ListAllResponder(resp)
	if err != nil {
		ae = autorest.NewErrorWithError(err, "network/RouteTablesClient", "ListAll", resp.StatusCode, "Failure responding to request")
	}

	return
}

// ListAllPreparer prepares the ListAll request.
func (client RouteTablesClient) ListAllPreparer() (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"subscriptionId": url.QueryEscape(client.SubscriptionID),
	}

	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	return autorest.Prepare(&http.Request{},
		autorest.AsJSON(),
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPath("/subscriptions/{subscriptionId}/providers/Microsoft.Network/routeTables"),
		autorest.WithPathParameters(pathParameters),
		autorest.WithQueryParameters(queryParameters))
}

// ListAllSender sends the ListAll request. The method will close the
// http.Response Body if it receives an error.
func (client RouteTablesClient) ListAllSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, http.StatusOK)
}

// ListAllResponder handles the response to the ListAll request. The method always
// closes the http.Response Body.
func (client RouteTablesClient) ListAllResponder(resp *http.Response) (result RouteTableListResult, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		autorest.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// ListAllNextResults retrieves the next set of results, if any.
func (client RouteTablesClient) ListAllNextResults(lastResults RouteTableListResult) (result RouteTableListResult, ae error) {
	req, err := lastResults.RouteTableListResultPreparer()
	if err != nil {
		return result, autorest.NewErrorWithError(err, "network/RouteTablesClient", "ListAll", autorest.UndefinedStatusCode, "Failure preparing next results request request")
	}
	if req == nil {
		return
	}

	resp, err := client.ListAllSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "network/RouteTablesClient", "ListAll", resp.StatusCode, "Failure sending next results request request")
	}

	result, err = client.ListAllResponder(resp)
	if err != nil {
		ae = autorest.NewErrorWithError(err, "network/RouteTablesClient", "ListAll", resp.StatusCode, "Failure responding to next results request request")
	}

	return
}
