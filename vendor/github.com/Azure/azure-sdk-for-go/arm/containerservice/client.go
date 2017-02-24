// Package containerservice implements the Azure ARM Containerservice service
// API version 2016-09-30.
//
// The Container Service Client.
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
)

const (
	// APIVersion is the version of the Containerservice
	APIVersion = "2016-09-30"

	// DefaultBaseURI is the default URI used for the service Containerservice
	DefaultBaseURI = "https://management.azure.com"
)

// ManagementClient is the base client for Containerservice.
type ManagementClient struct {
	autorest.Client
	BaseURI        string
	APIVersion     string
	SubscriptionID string
}

// New creates an instance of the ManagementClient client.
func New(subscriptionID string) ManagementClient {
	return NewWithBaseURI(DefaultBaseURI, subscriptionID)
}

// NewWithBaseURI creates an instance of the ManagementClient client.
func NewWithBaseURI(baseURI string, subscriptionID string) ManagementClient {
	return ManagementClient{
		Client:         autorest.NewClientWithUserAgent(UserAgent()),
		BaseURI:        baseURI,
		APIVersion:     APIVersion,
		SubscriptionID: subscriptionID,
	}
}
