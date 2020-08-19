package graphrbac

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
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/Azure/go-autorest/tracing"
	"net/http"
)

// DeletedApplicationsClient is the the Graph RBAC Management Client
type DeletedApplicationsClient struct {
	BaseClient
}

// NewDeletedApplicationsClient creates an instance of the DeletedApplicationsClient client.
func NewDeletedApplicationsClient(tenantID string) DeletedApplicationsClient {
	return NewDeletedApplicationsClientWithBaseURI(DefaultBaseURI, tenantID)
}

// NewDeletedApplicationsClientWithBaseURI creates an instance of the DeletedApplicationsClient client using a custom
// endpoint.  Use this when interacting with an Azure cloud that uses a non-standard base URI (sovereign clouds, Azure
// stack).
func NewDeletedApplicationsClientWithBaseURI(baseURI string, tenantID string) DeletedApplicationsClient {
	return DeletedApplicationsClient{NewWithBaseURI(baseURI, tenantID)}
}

// HardDelete hard-delete an application.
// Parameters:
// applicationObjectID - application object ID.
func (client DeletedApplicationsClient) HardDelete(ctx context.Context, applicationObjectID string) (result autorest.Response, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/DeletedApplicationsClient.HardDelete")
		defer func() {
			sc := -1
			if result.Response != nil {
				sc = result.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.HardDeletePreparer(ctx, applicationObjectID)
	if err != nil {
		err = autorest.NewErrorWithError(err, "graphrbac.DeletedApplicationsClient", "HardDelete", nil, "Failure preparing request")
		return
	}

	resp, err := client.HardDeleteSender(req)
	if err != nil {
		result.Response = resp
		err = autorest.NewErrorWithError(err, "graphrbac.DeletedApplicationsClient", "HardDelete", resp, "Failure sending request")
		return
	}

	result, err = client.HardDeleteResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "graphrbac.DeletedApplicationsClient", "HardDelete", resp, "Failure responding to request")
	}

	return
}

// HardDeletePreparer prepares the HardDelete request.
func (client DeletedApplicationsClient) HardDeletePreparer(ctx context.Context, applicationObjectID string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"applicationObjectId": autorest.Encode("path", applicationObjectID),
		"tenantID":            autorest.Encode("path", client.TenantID),
	}

	const APIVersion = "1.6"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsDelete(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/{tenantID}/deletedApplications/{applicationObjectId}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// HardDeleteSender sends the HardDelete request. The method will close the
// http.Response Body if it receives an error.
func (client DeletedApplicationsClient) HardDeleteSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
}

// HardDeleteResponder handles the response to the HardDelete request. The method always
// closes the http.Response Body.
func (client DeletedApplicationsClient) HardDeleteResponder(resp *http.Response) (result autorest.Response, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusNoContent),
		autorest.ByClosing())
	result.Response = resp
	return
}

// List gets a list of deleted applications in the directory.
// Parameters:
// filter - the filter to apply to the operation.
func (client DeletedApplicationsClient) List(ctx context.Context, filter string) (result ApplicationListResultPage, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/DeletedApplicationsClient.List")
		defer func() {
			sc := -1
			if result.alr.Response.Response != nil {
				sc = result.alr.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	result.fn = func(ctx context.Context, lastResult ApplicationListResult) (ApplicationListResult, error) {
		if lastResult.OdataNextLink == nil || len(to.String(lastResult.OdataNextLink)) < 1 {
			return ApplicationListResult{}, nil
		}
		return client.ListNext(ctx, *lastResult.OdataNextLink)
	}
	req, err := client.ListPreparer(ctx, filter)
	if err != nil {
		err = autorest.NewErrorWithError(err, "graphrbac.DeletedApplicationsClient", "List", nil, "Failure preparing request")
		return
	}

	resp, err := client.ListSender(req)
	if err != nil {
		result.alr.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "graphrbac.DeletedApplicationsClient", "List", resp, "Failure sending request")
		return
	}

	result.alr, err = client.ListResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "graphrbac.DeletedApplicationsClient", "List", resp, "Failure responding to request")
	}
	if result.alr.hasNextLink() && result.alr.IsEmpty() {
		err = result.NextWithContext(ctx)
	}

	return
}

// ListPreparer prepares the List request.
func (client DeletedApplicationsClient) ListPreparer(ctx context.Context, filter string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"tenantID": autorest.Encode("path", client.TenantID),
	}

	const APIVersion = "1.6"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}
	if len(filter) > 0 {
		queryParameters["$filter"] = autorest.Encode("query", filter)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/{tenantID}/deletedApplications", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// ListSender sends the List request. The method will close the
// http.Response Body if it receives an error.
func (client DeletedApplicationsClient) ListSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
}

// ListResponder handles the response to the List request. The method always
// closes the http.Response Body.
func (client DeletedApplicationsClient) ListResponder(resp *http.Response) (result ApplicationListResult, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// ListComplete enumerates all values, automatically crossing page boundaries as required.
func (client DeletedApplicationsClient) ListComplete(ctx context.Context, filter string) (result ApplicationListResultIterator, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/DeletedApplicationsClient.List")
		defer func() {
			sc := -1
			if result.Response().Response.Response != nil {
				sc = result.page.Response().Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	result.page, err = client.List(ctx, filter)
	return
}

// ListNext gets a list of deleted applications in the directory.
// Parameters:
// nextLink - next link for the list operation.
func (client DeletedApplicationsClient) ListNext(ctx context.Context, nextLink string) (result ApplicationListResult, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/DeletedApplicationsClient.ListNext")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.ListNextPreparer(ctx, nextLink)
	if err != nil {
		err = autorest.NewErrorWithError(err, "graphrbac.DeletedApplicationsClient", "ListNext", nil, "Failure preparing request")
		return
	}

	resp, err := client.ListNextSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "graphrbac.DeletedApplicationsClient", "ListNext", resp, "Failure sending request")
		return
	}

	result, err = client.ListNextResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "graphrbac.DeletedApplicationsClient", "ListNext", resp, "Failure responding to request")
	}

	return
}

// ListNextPreparer prepares the ListNext request.
func (client DeletedApplicationsClient) ListNextPreparer(ctx context.Context, nextLink string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"nextLink": nextLink,
		"tenantID": autorest.Encode("path", client.TenantID),
	}

	const APIVersion = "1.6"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/{tenantID}/{nextLink}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// ListNextSender sends the ListNext request. The method will close the
// http.Response Body if it receives an error.
func (client DeletedApplicationsClient) ListNextSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
}

// ListNextResponder handles the response to the ListNext request. The method always
// closes the http.Response Body.
func (client DeletedApplicationsClient) ListNextResponder(resp *http.Response) (result ApplicationListResult, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// Restore restores the deleted application in the directory.
// Parameters:
// objectID - application object ID.
func (client DeletedApplicationsClient) Restore(ctx context.Context, objectID string) (result Application, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/DeletedApplicationsClient.Restore")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.RestorePreparer(ctx, objectID)
	if err != nil {
		err = autorest.NewErrorWithError(err, "graphrbac.DeletedApplicationsClient", "Restore", nil, "Failure preparing request")
		return
	}

	resp, err := client.RestoreSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "graphrbac.DeletedApplicationsClient", "Restore", resp, "Failure sending request")
		return
	}

	result, err = client.RestoreResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "graphrbac.DeletedApplicationsClient", "Restore", resp, "Failure responding to request")
	}

	return
}

// RestorePreparer prepares the Restore request.
func (client DeletedApplicationsClient) RestorePreparer(ctx context.Context, objectID string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"objectId": autorest.Encode("path", objectID),
		"tenantID": autorest.Encode("path", client.TenantID),
	}

	const APIVersion = "1.6"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsPost(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/{tenantID}/deletedApplications/{objectId}/restore", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// RestoreSender sends the Restore request. The method will close the
// http.Response Body if it receives an error.
func (client DeletedApplicationsClient) RestoreSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
}

// RestoreResponder handles the response to the Restore request. The method always
// closes the http.Response Body.
func (client DeletedApplicationsClient) RestoreResponder(resp *http.Response) (result Application, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}
