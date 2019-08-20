package resources

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
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

import (
	"context"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"net/http"
)

// ProvidersClient is the client for the Providers methods of the Resources service.
type ProvidersClient struct {
	BaseClient
}

// NewProvidersClient creates an instance of the ProvidersClient client.
func NewProvidersClient(subscriptionID string) ProvidersClient {
	return NewProvidersClientWithBaseURI(DefaultBaseURI, subscriptionID)
}

// NewProvidersClientWithBaseURI creates an instance of the ProvidersClient client.
func NewProvidersClientWithBaseURI(baseURI string, subscriptionID string) ProvidersClient {
	return ProvidersClient{NewWithBaseURI(baseURI, subscriptionID)}
}

// Get gets a resource provider.
// Parameters:
// resourceProviderNamespace - namespace of the resource provider.
// expand - the $expand query parameter. e.g. To include property aliases in response, use
// $expand=resourceTypes/aliases.
func (client ProvidersClient) Get(ctx context.Context, resourceProviderNamespace string, expand string) (result Provider, err error) {
	req, err := client.GetPreparer(ctx, resourceProviderNamespace, expand)
	if err != nil {
		err = autorest.NewErrorWithError(err, "resources.ProvidersClient", "Get", nil, "Failure preparing request")
		return
	}

	resp, err := client.GetSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "resources.ProvidersClient", "Get", resp, "Failure sending request")
		return
	}

	result, err = client.GetResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "resources.ProvidersClient", "Get", resp, "Failure responding to request")
	}

	return
}

// GetPreparer prepares the Get request.
func (client ProvidersClient) GetPreparer(ctx context.Context, resourceProviderNamespace string, expand string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"resourceProviderNamespace": autorest.Encode("path", resourceProviderNamespace),
		"subscriptionId":            autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2016-02-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}
	if len(expand) > 0 {
		queryParameters["$expand"] = autorest.Encode("query", expand)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/providers/{resourceProviderNamespace}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// GetSender sends the Get request. The method will close the
// http.Response Body if it receives an error.
func (client ProvidersClient) GetSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req,
		azure.DoRetryWithRegistration(client.Client))
}

// GetResponder handles the response to the Get request. The method always
// closes the http.Response Body.
func (client ProvidersClient) GetResponder(resp *http.Response) (result Provider, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// List gets a list of resource providers.
// Parameters:
// top - query parameters. If null is passed returns all deployments.
// expand - the $expand query parameter. e.g. To include property aliases in response, use
// $expand=resourceTypes/aliases.
func (client ProvidersClient) List(ctx context.Context, top *int32, expand string) (result ProviderListResultPage, err error) {
	result.fn = client.listNextResults
	req, err := client.ListPreparer(ctx, top, expand)
	if err != nil {
		err = autorest.NewErrorWithError(err, "resources.ProvidersClient", "List", nil, "Failure preparing request")
		return
	}

	resp, err := client.ListSender(req)
	if err != nil {
		result.plr.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "resources.ProvidersClient", "List", resp, "Failure sending request")
		return
	}

	result.plr, err = client.ListResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "resources.ProvidersClient", "List", resp, "Failure responding to request")
	}

	return
}

// ListPreparer prepares the List request.
func (client ProvidersClient) ListPreparer(ctx context.Context, top *int32, expand string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"subscriptionId": autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2016-02-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}
	if top != nil {
		queryParameters["$top"] = autorest.Encode("query", *top)
	}
	if len(expand) > 0 {
		queryParameters["$expand"] = autorest.Encode("query", expand)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/providers", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// ListSender sends the List request. The method will close the
// http.Response Body if it receives an error.
func (client ProvidersClient) ListSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req,
		azure.DoRetryWithRegistration(client.Client))
}

// ListResponder handles the response to the List request. The method always
// closes the http.Response Body.
func (client ProvidersClient) ListResponder(resp *http.Response) (result ProviderListResult, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// listNextResults retrieves the next set of results, if any.
func (client ProvidersClient) listNextResults(lastResults ProviderListResult) (result ProviderListResult, err error) {
	req, err := lastResults.providerListResultPreparer()
	if err != nil {
		return result, autorest.NewErrorWithError(err, "resources.ProvidersClient", "listNextResults", nil, "Failure preparing next results request")
	}
	if req == nil {
		return
	}
	resp, err := client.ListSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "resources.ProvidersClient", "listNextResults", resp, "Failure sending next results request")
	}
	result, err = client.ListResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "resources.ProvidersClient", "listNextResults", resp, "Failure responding to next results request")
	}
	return
}

// ListComplete enumerates all values, automatically crossing page boundaries as required.
func (client ProvidersClient) ListComplete(ctx context.Context, top *int32, expand string) (result ProviderListResultIterator, err error) {
	result.page, err = client.List(ctx, top, expand)
	return
}

// Register registers provider to be used with a subscription.
// Parameters:
// resourceProviderNamespace - namespace of the resource provider.
func (client ProvidersClient) Register(ctx context.Context, resourceProviderNamespace string) (result Provider, err error) {
	req, err := client.RegisterPreparer(ctx, resourceProviderNamespace)
	if err != nil {
		err = autorest.NewErrorWithError(err, "resources.ProvidersClient", "Register", nil, "Failure preparing request")
		return
	}

	resp, err := client.RegisterSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "resources.ProvidersClient", "Register", resp, "Failure sending request")
		return
	}

	result, err = client.RegisterResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "resources.ProvidersClient", "Register", resp, "Failure responding to request")
	}

	return
}

// RegisterPreparer prepares the Register request.
func (client ProvidersClient) RegisterPreparer(ctx context.Context, resourceProviderNamespace string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"resourceProviderNamespace": autorest.Encode("path", resourceProviderNamespace),
		"subscriptionId":            autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2016-02-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsPost(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/providers/{resourceProviderNamespace}/register", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// RegisterSender sends the Register request. The method will close the
// http.Response Body if it receives an error.
func (client ProvidersClient) RegisterSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req,
		azure.DoRetryWithRegistration(client.Client))
}

// RegisterResponder handles the response to the Register request. The method always
// closes the http.Response Body.
func (client ProvidersClient) RegisterResponder(resp *http.Response) (result Provider, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// Unregister unregisters provider from a subscription.
// Parameters:
// resourceProviderNamespace - namespace of the resource provider.
func (client ProvidersClient) Unregister(ctx context.Context, resourceProviderNamespace string) (result Provider, err error) {
	req, err := client.UnregisterPreparer(ctx, resourceProviderNamespace)
	if err != nil {
		err = autorest.NewErrorWithError(err, "resources.ProvidersClient", "Unregister", nil, "Failure preparing request")
		return
	}

	resp, err := client.UnregisterSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "resources.ProvidersClient", "Unregister", resp, "Failure sending request")
		return
	}

	result, err = client.UnregisterResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "resources.ProvidersClient", "Unregister", resp, "Failure responding to request")
	}

	return
}

// UnregisterPreparer prepares the Unregister request.
func (client ProvidersClient) UnregisterPreparer(ctx context.Context, resourceProviderNamespace string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"resourceProviderNamespace": autorest.Encode("path", resourceProviderNamespace),
		"subscriptionId":            autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2016-02-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsPost(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/providers/{resourceProviderNamespace}/unregister", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// UnregisterSender sends the Unregister request. The method will close the
// http.Response Body if it receives an error.
func (client ProvidersClient) UnregisterSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req,
		azure.DoRetryWithRegistration(client.Client))
}

// UnregisterResponder handles the response to the Unregister request. The method always
// closes the http.Response Body.
func (client ProvidersClient) UnregisterResponder(resp *http.Response) (result Provider, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}
