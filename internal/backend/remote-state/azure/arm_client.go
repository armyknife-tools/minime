// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package azure

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/profiles/2017-03-09/resources/mgmt/resources"
	armStorage "github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2021-01-01/storage"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/authentication"
	"github.com/hashicorp/go-azure-helpers/sender"
	"github.com/manicminer/hamilton/environments"
	"github.com/placeholderplaceholderplaceholder/opentf/internal/httpclient"
	"github.com/placeholderplaceholderplaceholder/opentf/version"
	"github.com/tombuildsstuff/giovanni/storage/2018-11-09/blob/blobs"
	"github.com/tombuildsstuff/giovanni/storage/2018-11-09/blob/containers"
)

type ArmClient struct {
	// These Clients are only initialized if an Access Key isn't provided
	groupsClient          *resources.GroupsClient
	storageAccountsClient *armStorage.AccountsClient
	containersClient      *containers.Client
	blobsClient           *blobs.Client

	// azureAdStorageAuth is only here if we're using AzureAD Authentication but is an Authorizer for Storage
	azureAdStorageAuth *autorest.Authorizer

	accessKey          string
	environment        azure.Environment
	resourceGroupName  string
	storageAccountName string
	sasToken           string
}

func buildArmClient(ctx context.Context, config BackendConfig) (*ArmClient, error) {
	env, err := authentication.AzureEnvironmentByNameFromEndpoint(ctx, config.MetadataHost, config.Environment)
	if err != nil {
		return nil, err
	}

	client := ArmClient{
		environment:        *env,
		resourceGroupName:  config.ResourceGroupName,
		storageAccountName: config.StorageAccountName,
	}

	// if we have an Access Key - we don't need the other clients
	if config.AccessKey != "" {
		client.accessKey = config.AccessKey
		return &client, nil
	}

	// likewise with a SAS token
	if config.SasToken != "" {
		client.sasToken = config.SasToken
		return &client, nil
	}

	builder := authentication.Builder{
		ClientID:                      config.ClientID,
		SubscriptionID:                config.SubscriptionID,
		TenantID:                      config.TenantID,
		CustomResourceManagerEndpoint: config.CustomResourceManagerEndpoint,
		MetadataHost:                  config.MetadataHost,
		Environment:                   config.Environment,
		ClientSecretDocsLink:          "https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/guides/service_principal_client_secret",

		// Service Principal (Client Certificate)
		ClientCertPassword: config.ClientCertificatePassword,
		ClientCertPath:     config.ClientCertificatePath,

		// Service Principal (Client Secret)
		ClientSecret: config.ClientSecret,

		// Managed Service Identity
		MsiEndpoint: config.MsiEndpoint,

		// OIDC
		IDToken:             config.OIDCToken,
		IDTokenFilePath:     config.OIDCTokenFilePath,
		IDTokenRequestURL:   config.OIDCRequestURL,
		IDTokenRequestToken: config.OIDCRequestToken,

		// Feature Toggles
		SupportsAzureCliToken:          true,
		SupportsClientCertAuth:         true,
		SupportsClientSecretAuth:       true,
		SupportsManagedServiceIdentity: config.UseMsi,
		SupportsOIDCAuth:               config.UseOIDC,
		UseMicrosoftGraph:              true,
	}
	armConfig, err := builder.Build()
	if err != nil {
		return nil, fmt.Errorf("Error building ARM Config: %w", err)
	}

	oauthConfig, err := armConfig.BuildOAuthConfig(env.ActiveDirectoryEndpoint)
	if err != nil {
		return nil, err
	}

	hamiltonEnv, err := environments.EnvironmentFromString(config.Environment)
	if err != nil {
		return nil, err
	}

	sender := sender.BuildSender("backend/remote-state/azure")
	log.Printf("[DEBUG] Obtaining an MSAL / Microsoft Graph token for Resource Manager..")
	auth, err := armConfig.GetMSALToken(ctx, hamiltonEnv.ResourceManager, sender, oauthConfig, env.TokenAudience)
	if err != nil {
		return nil, err
	}

	if config.UseAzureADAuthentication {
		log.Printf("[DEBUG] Obtaining an MSAL / Microsoft Graph token for Storage..")
		storageAuth, err := armConfig.GetMSALToken(ctx, hamiltonEnv.Storage, sender, oauthConfig, env.ResourceIdentifiers.Storage)
		if err != nil {
			return nil, err
		}
		client.azureAdStorageAuth = &storageAuth
	}

	accountsClient := armStorage.NewAccountsClientWithBaseURI(env.ResourceManagerEndpoint, armConfig.SubscriptionID)
	client.configureClient(&accountsClient.Client, auth)
	client.storageAccountsClient = &accountsClient

	groupsClient := resources.NewGroupsClientWithBaseURI(env.ResourceManagerEndpoint, armConfig.SubscriptionID)
	client.configureClient(&groupsClient.Client, auth)
	client.groupsClient = &groupsClient

	return &client, nil
}

func (c ArmClient) getBlobClient(ctx context.Context) (*blobs.Client, error) {
	if c.sasToken != "" {
		log.Printf("[DEBUG] Building the Blob Client from a SAS Token")
		storageAuth, err := autorest.NewSASTokenAuthorizer(c.sasToken)
		if err != nil {
			return nil, fmt.Errorf("Error building SAS Token Authorizer: %w", err)
		}

		blobsClient := blobs.NewWithEnvironment(c.environment)
		c.configureClient(&blobsClient.Client, storageAuth)
		return &blobsClient, nil
	}

	if c.azureAdStorageAuth != nil {
		blobsClient := blobs.NewWithEnvironment(c.environment)
		c.configureClient(&blobsClient.Client, *c.azureAdStorageAuth)
		return &blobsClient, nil
	}

	accessKey := c.accessKey
	if accessKey == "" {
		log.Printf("[DEBUG] Building the Blob Client from an Access Token (using user credentials)")
		keys, err := c.storageAccountsClient.ListKeys(ctx, c.resourceGroupName, c.storageAccountName, "")
		if err != nil {
			return nil, fmt.Errorf("Error retrieving keys for Storage Account %q: %w", c.storageAccountName, err)
		}

		if keys.Keys == nil {
			return nil, fmt.Errorf("Nil key returned for storage account %q", c.storageAccountName)
		}

		accessKeys := *keys.Keys
		accessKey = *accessKeys[0].Value
	}

	storageAuth, err := autorest.NewSharedKeyAuthorizer(c.storageAccountName, accessKey, autorest.SharedKey)
	if err != nil {
		return nil, fmt.Errorf("Error building Shared Key Authorizer: %w", err)
	}

	blobsClient := blobs.NewWithEnvironment(c.environment)
	c.configureClient(&blobsClient.Client, storageAuth)
	return &blobsClient, nil
}

func (c ArmClient) getContainersClient(ctx context.Context) (*containers.Client, error) {
	if c.sasToken != "" {
		log.Printf("[DEBUG] Building the Container Client from a SAS Token")
		storageAuth, err := autorest.NewSASTokenAuthorizer(c.sasToken)
		if err != nil {
			return nil, fmt.Errorf("Error building SAS Token Authorizer: %w", err)
		}

		containersClient := containers.NewWithEnvironment(c.environment)
		c.configureClient(&containersClient.Client, storageAuth)
		return &containersClient, nil
	}

	if c.azureAdStorageAuth != nil {
		containersClient := containers.NewWithEnvironment(c.environment)
		c.configureClient(&containersClient.Client, *c.azureAdStorageAuth)
		return &containersClient, nil
	}

	accessKey := c.accessKey
	if accessKey == "" {
		log.Printf("[DEBUG] Building the Container Client from an Access Token (using user credentials)")
		keys, err := c.storageAccountsClient.ListKeys(ctx, c.resourceGroupName, c.storageAccountName, "")
		if err != nil {
			return nil, fmt.Errorf("Error retrieving keys for Storage Account %q: %w", c.storageAccountName, err)
		}

		if keys.Keys == nil {
			return nil, fmt.Errorf("Nil key returned for storage account %q", c.storageAccountName)
		}

		accessKeys := *keys.Keys
		accessKey = *accessKeys[0].Value
	}

	storageAuth, err := autorest.NewSharedKeyAuthorizer(c.storageAccountName, accessKey, autorest.SharedKey)
	if err != nil {
		return nil, fmt.Errorf("Error building Shared Key Authorizer: %w", err)
	}

	containersClient := containers.NewWithEnvironment(c.environment)
	c.configureClient(&containersClient.Client, storageAuth)
	return &containersClient, nil
}

func (c *ArmClient) configureClient(client *autorest.Client, auth autorest.Authorizer) {
	client.UserAgent = buildUserAgent()
	client.Authorizer = auth
	client.Sender = buildSender()
	client.SkipResourceProviderRegistration = false
	client.PollingDuration = 60 * time.Minute
}

func buildUserAgent() string {
	userAgent := httpclient.OpenTfUserAgent(version.Version)

	// append the CloudShell version to the user agent if it exists
	if azureAgent := os.Getenv("AZURE_HTTP_USER_AGENT"); azureAgent != "" {
		userAgent = fmt.Sprintf("%s %s", userAgent, azureAgent)
	}

	return userAgent
}
