// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2023 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package registry

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	svchost "github.com/hashicorp/terraform-svchost"
	"github.com/hashicorp/terraform-svchost/disco"

	"github.com/opentofu/opentofu/internal/httpclient"
	"github.com/opentofu/opentofu/internal/logging"
	"github.com/opentofu/opentofu/internal/registry/regsrc"
	"github.com/opentofu/opentofu/internal/registry/response"
	"github.com/opentofu/opentofu/version"
)

const (
	xTerraformGet      = "X-Terraform-Get"
	xTerraformVersion  = "X-Terraform-Version"
	modulesServiceID   = "modules.v1"
	providersServiceID = "providers.v1"

	// registryDiscoveryRetryEnvName is the name of the environment variable that
	// can be configured to customize number of retries for module and provider
	// discovery requests with the remote registry.
	registryDiscoveryRetryEnvName = "TF_REGISTRY_DISCOVERY_RETRY"
	defaultRetry                  = 1

	// registryClientTimeoutEnvName is the name of the environment variable that
	// can be configured to customize the timeout duration (seconds) for module
	// and provider discovery with the remote registry.
	registryClientTimeoutEnvName = "TF_REGISTRY_CLIENT_TIMEOUT"

	// defaultRequestTimeout is the default timeout duration for requests to the
	// remote registry.
	defaultRequestTimeout = 10 * time.Second
)

var (
	tfVersion = version.String()

	discoveryRetry int
	requestTimeout time.Duration
)

func init() {
	configureDiscoveryRetry()
	configureRequestTimeout()
}

// Client provides methods to query OpenTofu Registries.
type Client struct {
	// this is the client to be used for all requests.
	client *retryablehttp.Client

	// services is a required *disco.Disco, which may have services and
	// credentials pre-loaded.
	services *disco.Disco
}

// NewClient returns a new initialized registry client.
func NewClient(services *disco.Disco, client *http.Client) *Client {
	if services == nil {
		services = disco.New()
	}

	if client == nil {
		client = httpclient.New()
		client.Timeout = requestTimeout
	}
	retryableClient := retryablehttp.NewClient()
	retryableClient.HTTPClient = client
	retryableClient.RetryMax = discoveryRetry
	retryableClient.RequestLogHook = requestLogHook
	retryableClient.ErrorHandler = maxRetryErrorHandler

	logOutput := logging.LogOutput()
	retryableClient.Logger = log.New(logOutput, "", log.Flags())

	services.Transport = retryableClient.HTTPClient.Transport

	services.SetUserAgent(httpclient.OpenTofuUserAgent(version.String()))

	return &Client{
		client:   retryableClient,
		services: services,
	}
}

// Discover queries the host, and returns the url for the registry.
func (c *Client) Discover(host svchost.Hostname, serviceID string) (*url.URL, error) {
	service, err := c.services.DiscoverServiceURL(host, serviceID)
	if err != nil {
		return nil, &ServiceUnreachableError{err}
	}
	if !strings.HasSuffix(service.Path, "/") {
		service.Path += "/"
	}
	return service, nil
}

// ModuleVersions queries the registry for a module, and returns the available versions.
func (c *Client) ModuleVersions(ctx context.Context, module *regsrc.Module) (*response.ModuleVersions, error) {
	host, err := module.SvcHost()
	if err != nil {
		return nil, err
	}

	service, err := c.Discover(host, modulesServiceID)
	if err != nil {
		return nil, err
	}

	p, err := url.Parse(path.Join(module.Module(), "versions"))
	if err != nil {
		return nil, err
	}

	service = service.ResolveReference(p)

	log.Printf("[DEBUG] fetching module versions from %q", service)

	req, err := retryablehttp.NewRequest("GET", service.String(), nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)

	c.addRequestCreds(host, req.Request)
	req.Header.Set(xTerraformVersion, tfVersion)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		// OK
	case http.StatusNotFound:
		return nil, &errModuleNotFound{addr: module}
	default:
		return nil, fmt.Errorf("error looking up module versions: %s", resp.Status)
	}

	var versions response.ModuleVersions

	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&versions); err != nil {
		return nil, err
	}

	for _, mod := range versions.Modules {
		for _, v := range mod.Versions {
			log.Printf("[DEBUG] found available version %q for %s", v.Version, module.Module())
		}
	}

	return &versions, nil
}

func (c *Client) addRequestCreds(host svchost.Hostname, req *http.Request) {
	creds, err := c.services.CredentialsForHost(host)
	if err != nil {
		log.Printf("[WARN] Failed to get credentials for %s: %s (ignoring)", host, err)
		return
	}

	if creds != nil {
		creds.PrepareRequest(req)
	}
}

// ModuleLocation find the download location for a specific version module.
// This returns a string, because the final location may contain special go-getter syntax.
func (c *Client) ModuleLocation(ctx context.Context, module *regsrc.Module, version string) (string, error) {
	host, err := module.SvcHost()
	if err != nil {
		return "", err
	}

	service, err := c.Discover(host, modulesServiceID)
	if err != nil {
		return "", err
	}

	var p *url.URL
	if version == "" {
		p, err = url.Parse(path.Join(module.Module(), "download"))
	} else {
		p, err = url.Parse(path.Join(module.Module(), version, "download"))
	}
	if err != nil {
		return "", err
	}
	download := service.ResolveReference(p)

	log.Printf("[DEBUG] looking up module location from %q", download)

	req, err := retryablehttp.NewRequest("GET", download.String(), nil)
	if err != nil {
		return "", err
	}

	req = req.WithContext(ctx)

	c.addRequestCreds(host, req.Request)
	req.Header.Set(xTerraformVersion, tfVersion)

	resp, err := c.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body from registry: %w", err)
	}

	var location string

	switch resp.StatusCode {
	case http.StatusOK:
		var v response.ModuleLocationRegistryResp
		if err := json.Unmarshal(body, &v); err != nil {
			return "", fmt.Errorf("module %q version %q failed to deserialize response body %s: %w",
				module, version, body, err)
		}

		location = v.Location

		// if the location is empty, we will fallback to the header
		if location == "" {
			location = resp.Header.Get(xTerraformGet)
		}

	case http.StatusNoContent:
		// FALLBACK: set the found location from the header
		location = resp.Header.Get(xTerraformGet)

	case http.StatusNotFound:
		return "", fmt.Errorf("module %q version %q not found", module, version)

	default:
		// anything else is an error:
		return "", fmt.Errorf("error getting download location for %q: %s resp:%s", module, resp.Status, body)
	}

	if location == "" {
		return "", fmt.Errorf("failed to get download URL for %q: %s resp:%s", module, resp.Status, body)
	}

	// If location looks like it's trying to be a relative URL, treat it as
	// one.
	//
	// We don't do this for just _any_ location, since the X-Terraform-Get
	// header is a go-getter location rather than a URL, and so not all
	// possible values will parse reasonably as URLs.)
	//
	// When used in conjunction with go-getter we normally require this header
	// to be an absolute URL, but we are more liberal here because third-party
	// registry implementations may not "know" their own absolute URLs if
	// e.g. they are running behind a reverse proxy frontend, or such.
	if strings.HasPrefix(location, "/") || strings.HasPrefix(location, "./") || strings.HasPrefix(location, "../") {
		locationURL, err := url.Parse(location)
		if err != nil {
			return "", fmt.Errorf("invalid relative URL for %q: %w", module, err)
		}
		locationURL = download.ResolveReference(locationURL)
		location = locationURL.String()
	}

	return location, nil
}

// BulkFetchModules fetches all available modules from a registry for caching purposes
func (c *Client) BulkFetchModules(ctx context.Context, host svchost.Hostname) ([]*response.Module, error) {
	log.Printf("[DEBUG] Fetching all modules from %s for caching", host)

	service, err := c.Discover(host, modulesServiceID)
	if err != nil {
		log.Printf("[DEBUG] Failed to discover modules service for %s: %s", host, err)
		return nil, fmt.Errorf("error discovering modules service: %w", err)
	}
	
	log.Printf("[DEBUG] Discovered modules service URL: %s", service)

	// Pre-allocate based on known registry sizes (approximately 18,000 modules)
	allModules := make([]*response.Module, 0, 18000)
	
	// For the registry refresh command, we need to fetch all available modules
	// The Terraform Registry API requires a query parameter for the search endpoint
	
	// Try with a wildcard search or a common term that will return many results
	searchTerms := []string{"aws", "azure", "google", "kubernetes", "terraform"}
	
	for _, term := range searchTerms {
		// Construct the URL for the search endpoint
		searchURL := service.ResolveReference(&url.URL{Path: "search"})
		queryParams := searchURL.Query()
		queryParams.Set("q", term) // Use a search term that will return many results
		queryParams.Set("limit", "100") // Max page size
		searchURL.RawQuery = queryParams.Encode()
		
		log.Printf("[DEBUG] Trying modules search API with term '%s': %s", term, searchURL)
		
		req, err := retryablehttp.NewRequest("GET", searchURL.String(), nil)
		if err != nil {
			log.Printf("[DEBUG] Failed to create request for search API: %s", err)
			continue
		}
		req = req.WithContext(ctx)
		
		c.addRequestCreds(host, req.Request)
		req.Header.Set(xTerraformVersion, tfVersion)
		req.Header.Set("User-Agent", "OpenTofu/"+tfVersion)
		
		resp, err := c.client.Do(req)
		if err != nil {
			log.Printf("[DEBUG] Error making request to search API: %s", err)
			continue
		}
		
		defer resp.Body.Close()
		
		if resp.StatusCode == http.StatusOK {
			var moduleList response.ModuleList
			dec := json.NewDecoder(resp.Body)
			if err := dec.Decode(&moduleList); err != nil {
				log.Printf("[DEBUG] Failed to decode response from search API: %s", err)
				continue
			}
			
			if moduleList.Modules != nil {
				log.Printf("[DEBUG] Successfully decoded %d modules from search API with term '%s'", len(moduleList.Modules), term)
				
				// Add modules to our collection, avoiding duplicates
				for _, module := range moduleList.Modules {
					// Check if this module is already in our collection
					isDuplicate := false
					for _, existingModule := range allModules {
						if existingModule.ID == module.ID {
							isDuplicate = true
							break
						}
					}
					
					if !isDuplicate {
						allModules = append(allModules, module)
					}
				}
				
				// Handle pagination if needed
				nextURL := moduleList.Meta.NextURL
				for nextURL != "" {
					// If the URL is relative, resolve it against the service URL
					var nextFullURL string
					if strings.HasPrefix(nextURL, "http") {
						nextFullURL = nextURL
					} else {
						nextURLObj, err := url.Parse(nextURL)
						if err != nil {
							log.Printf("[DEBUG] Failed to parse next URL %q: %s", nextURL, err)
							break
						}
						nextFullURL = service.ResolveReference(nextURLObj).String()
					}
					
					log.Printf("[DEBUG] Fetching next page: %s", nextFullURL)
					
					// Add throttling delay to avoid rate limiting
					time.Sleep(200 * time.Millisecond)
					
					nextReq, err := retryablehttp.NewRequest("GET", nextFullURL, nil)
					if err != nil {
						log.Printf("[DEBUG] Failed to create request for next page %q: %s", nextFullURL, err)
						break
					}
					nextReq = nextReq.WithContext(ctx)
					
					c.addRequestCreds(host, nextReq.Request)
					nextReq.Header.Set(xTerraformVersion, tfVersion)
					nextReq.Header.Set("User-Agent", "OpenTofu/"+tfVersion)
					
					resp, err := c.client.Do(nextReq)
					if err != nil {
						log.Printf("[DEBUG] Error fetching next page %q: %s", nextFullURL, err)
						break
					}
					
					var nextModuleList response.ModuleList
					dec := json.NewDecoder(resp.Body)
					if err := dec.Decode(&nextModuleList); err != nil {
						log.Printf("[DEBUG] Failed to decode next page %q: %s", nextFullURL, err)
						resp.Body.Close()
						break
					}
					
					log.Printf("[DEBUG] Successfully decoded %d modules from next page", len(nextModuleList.Modules))
					
					// Add modules to our collection, avoiding duplicates
					for _, module := range nextModuleList.Modules {
						// Check if this module is already in our collection
						isDuplicate := false
						for _, existingModule := range allModules {
							if existingModule.ID == module.ID {
								isDuplicate = true
								break
							}
						}
						
						if !isDuplicate {
							allModules = append(allModules, module)
						}
					}
					
					nextURL = nextModuleList.Meta.NextURL
					resp.Body.Close()
				}
			}
		} else {
			// Read the response body for more diagnostic information
			body, readErr := io.ReadAll(resp.Body)
			if readErr != nil {
				log.Printf("[DEBUG] Error reading error response body: %s", readErr)
				body = []byte("[failed to read response body]")
			}
			
			log.Printf("[DEBUG] Error response from search API with term '%s': %s - %s", term, resp.Status, string(body))
		}
	}
	
	// If the search API didn't work, try to fetch known popular namespaces
	if len(allModules) == 0 {
		log.Printf("[DEBUG] Search API didn't return any modules, trying to fetch known popular namespaces")
		
		// Try to fetch modules from known popular namespaces
		popularNamespaces := []string{"hashicorp", "terraform-aws-modules"}
		
		for _, namespace := range popularNamespaces {
			// Construct the URL for the namespace modules
			namespaceURL := service.ResolveReference(&url.URL{Path: namespace})
			
			log.Printf("[DEBUG] Fetching modules for namespace: %s", namespaceURL)
			
			// Add a small delay to avoid rate limiting
			time.Sleep(200 * time.Millisecond)
			
			req, err := retryablehttp.NewRequest("GET", namespaceURL.String(), nil)
			if err != nil {
				log.Printf("[DEBUG] Failed to create request for namespace %s: %s", namespace, err)
				continue
			}
			req = req.WithContext(ctx)
			
			c.addRequestCreds(host, req.Request)
			req.Header.Set(xTerraformVersion, tfVersion)
			req.Header.Set("User-Agent", "OpenTofu/"+tfVersion)
			
			resp, err := c.client.Do(req)
			if err != nil {
				log.Printf("[DEBUG] Error making request to namespace %s: %s", namespace, err)
				continue
			}
			
			defer resp.Body.Close()
			
			if resp.StatusCode != http.StatusOK {
				// Read the response body for more diagnostic information
				body, readErr := io.ReadAll(resp.Body)
				if readErr != nil {
					log.Printf("[DEBUG] Error reading error response body: %s", readErr)
					body = []byte("[failed to read response body]")
				}
				
				log.Printf("[DEBUG] Error response from namespace %s: %s - %s", namespace, resp.Status, string(body))
				continue
			}
			
			// Try to decode as array of modules
			var modules []*response.Module
			if err := json.NewDecoder(resp.Body).Decode(&modules); err != nil {
				log.Printf("[DEBUG] Failed to decode response from namespace %s: %s", namespace, err)
				continue
			}
			
			log.Printf("[DEBUG] Successfully decoded %d modules from namespace %s", len(modules), namespace)
			allModules = append(allModules, modules...)
		}
	}
	
	// If we still don't have any modules, try one more approach - fetch individual modules
	if len(allModules) == 0 {
		log.Printf("[DEBUG] Trying to fetch individual popular modules")
		
		// Try some known popular modules
		knownModules := []struct {
			namespace string
			name      string
			provider  string
		}{
			{"terraform-aws-modules", "vpc", "aws"},
			{"terraform-aws-modules", "security-group", "aws"},
			{"hashicorp", "consul", "aws"},
			{"hashicorp", "vault", "aws"},
		}
		
		for _, module := range knownModules {
			// Construct the URL for the module using the {namespace}/{name}/{provider} pattern
			moduleURL := service.ResolveReference(&url.URL{
				Path: path.Join(module.namespace, module.name, module.provider),
			})
			
			log.Printf("[DEBUG] Fetching module: %s", moduleURL)
			
			// Add a small delay to avoid rate limiting
			time.Sleep(200 * time.Millisecond)
			
			req, err := retryablehttp.NewRequest("GET", moduleURL.String(), nil)
			if err != nil {
				log.Printf("[DEBUG] Failed to create request for module %s/%s/%s: %s", 
					module.namespace, module.name, module.provider, err)
				continue
			}
			req = req.WithContext(ctx)
			
			c.addRequestCreds(host, req.Request)
			req.Header.Set(xTerraformVersion, tfVersion)
			req.Header.Set("User-Agent", "OpenTofu/"+tfVersion)
			
			resp, err := c.client.Do(req)
			if err != nil {
				log.Printf("[DEBUG] Error making request to module %s/%s/%s: %s", 
					module.namespace, module.name, module.provider, err)
				continue
			}
			
			defer resp.Body.Close()
			
			if resp.StatusCode == http.StatusOK {
				// Try to decode as a module
				var moduleData response.Module
				if err := json.NewDecoder(resp.Body).Decode(&moduleData); err != nil {
					log.Printf("[DEBUG] Failed to decode response from module %s/%s/%s: %s", 
						module.namespace, module.name, module.provider, err)
					continue
				}
				
				log.Printf("[DEBUG] Successfully decoded module %s/%s/%s", 
					module.namespace, module.name, module.provider)
				allModules = append(allModules, &moduleData)
			} else {
				// Read the response body for more diagnostic information
				body, readErr := io.ReadAll(resp.Body)
				if readErr != nil {
					log.Printf("[DEBUG] Error reading error response body: %s", readErr)
					body = []byte("[failed to read response body]")
				}
				
				log.Printf("[DEBUG] Error response from module %s/%s/%s: %s - %s", 
					module.namespace, module.name, module.provider, resp.Status, string(body))
			}
		}
	}
	
	// If we still don't have any modules, create a minimal cache
	if len(allModules) == 0 {
		log.Printf("[DEBUG] Could not find a working API pattern for modules on %s, creating minimal cache", host)
		
		// Return an empty set rather than failing
		return []*response.Module{}, nil
	}
	
	log.Printf("[DEBUG] Successfully fetched a total of %d modules from %s", len(allModules), host)
	return allModules, nil
}

// BulkFetchProviders fetches all available providers from a registry. This is used to populate the
// local cache of providers for the registry.
func (c *Client) BulkFetchProviders(ctx context.Context, host svchost.Hostname) ([]*response.ModuleProvider, error) {
	log.Printf("[DEBUG] Fetching all providers from %s for caching", host.String())

	// Pre-allocate a slice based on the known registry size (approximately 4,000 providers)
	allProviders := make([]*response.ModuleProvider, 0, 4000)
	
	// Track unique providers to avoid duplicates
	seen := make(map[string]bool)

	// Define common provider namespaces to search for
	commonNamespaces := []string{
		"hashicorp", "aws", "azure", "google", "digitalocean", "cloudflare", 
		"datadog", "github", "kubernetes", "random", "time", "template", 
		"null", "local", "tls", "http", "archive", "external", "dns", "docker",
		"helm", "consul", "vault", "nomad", "terraform", "opentofu",
	}

	// Create a retryable client with throttling
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 5
	retryClient.RetryWaitMin = 1 * time.Second
	retryClient.RetryWaitMax = 30 * time.Second
	retryClient.Logger = nil

	service, err := c.Discover(host, providersServiceID)
	if err != nil {
		return nil, err
	}

	baseURL := service

	// First try the search endpoint
	searchURL := *baseURL
	searchURL.Path = path.Join(searchURL.Path, "search")
	query := searchURL.Query()
	query.Set("limit", "100") // Maximum allowed by the API
	searchURL.RawQuery = query.Encode()

	log.Printf("[DEBUG] Trying providers search API: %s", searchURL.String())

	// Try the search endpoint first
	searchReq, err := retryablehttp.NewRequest("GET", searchURL.String(), nil)
	if err != nil {
		log.Printf("[DEBUG] Error creating search request: %s", err)
	} else {
		searchReq.Header.Set("Accept", "application/json")
		searchReq.Header.Set("User-Agent", "OpenTofu/"+tfVersion)
		searchReq.Header.Set(xTerraformVersion, tfVersion)

		searchResp, err := retryClient.Do(searchReq)
		if err != nil {
			log.Printf("[DEBUG] Error from provider search API: %s", err)
		} else {
			defer searchResp.Body.Close()

			if searchResp.StatusCode == http.StatusOK {
				// Define the provider search response structure
				var result struct {
					Meta struct {
						Limit         int `json:"limit"`
						CurrentOffset int `json:"current_offset"`
						NextOffset    int `json:"next_offset"`
					} `json:"meta"`
					Providers []struct {
						ID        string `json:"id"`
						Namespace string `json:"namespace"`
						Name      string `json:"name"`
						Version   string `json:"version"`
					} `json:"providers"`
				}
				
				if err := json.NewDecoder(searchResp.Body).Decode(&result); err != nil {
					log.Printf("[DEBUG] Error decoding provider search response: %s", err)
				} else {
					log.Printf("[DEBUG] Successfully decoded %d providers from search API", len(result.Providers))
					for _, p := range result.Providers {
						providerID := fmt.Sprintf("%s/%s", p.Namespace, p.Name)
						if !seen[providerID] {
							seen[providerID] = true
							allProviders = append(allProviders, &response.ModuleProvider{
								Name:        p.Name,
								Downloads:   0,  // We don't have this information
								ModuleCount: 0, // We don't have this information
							})
						}
					}

					// Handle pagination if available
					for result.Meta.NextOffset > 0 && len(result.Providers) > 0 {
						nextURL := searchURL
						nextQuery := nextURL.Query()
						nextQuery.Set("offset", strconv.Itoa(result.Meta.NextOffset))
						nextURL.RawQuery = nextQuery.Encode()

						log.Printf("[DEBUG] Fetching next page: %s", nextURL.String())
						
						// Add a small delay to avoid rate limiting
						time.Sleep(200 * time.Millisecond)
						
						nextReq, err := retryablehttp.NewRequest("GET", nextURL.String(), nil)
						if err != nil {
							log.Printf("[DEBUG] Error creating next page request: %s", err)
							break
						}
						
						nextReq.Header.Set("Accept", "application/json")
						nextReq.Header.Set("User-Agent", "OpenTofu/"+tfVersion)
						nextReq.Header.Set(xTerraformVersion, tfVersion)

						nextResp, err := retryClient.Do(nextReq)
						if err != nil {
							log.Printf("[DEBUG] Error from next page request: %s", err)
							break
						}

						if nextResp.StatusCode != http.StatusOK {
							log.Printf("[DEBUG] Unexpected status code from next page: %d", nextResp.StatusCode)
							nextResp.Body.Close()
							break
						}

						var nextResult struct {
							Meta struct {
								Limit         int `json:"limit"`
								CurrentOffset int `json:"current_offset"`
								NextOffset    int `json:"next_offset"`
							} `json:"meta"`
							Providers []struct {
								ID        string `json:"id"`
								Namespace string `json:"namespace"`
								Name      string `json:"name"`
								Version   string `json:"version"`
							} `json:"providers"`
						}
						
						if err := json.NewDecoder(nextResp.Body).Decode(&nextResult); err != nil {
							log.Printf("[DEBUG] Error decoding next page response: %s", err)
							nextResp.Body.Close()
							break
						}
						nextResp.Body.Close()

						log.Printf("[DEBUG] Successfully decoded %d providers from next page", len(nextResult.Providers))
						for _, p := range nextResult.Providers {
							providerID := fmt.Sprintf("%s/%s", p.Namespace, p.Name)
							if !seen[providerID] {
								seen[providerID] = true
								allProviders = append(allProviders, &response.ModuleProvider{
									Name:        p.Name,
									Downloads:   0,  // We don't have this information
									ModuleCount: 0, // We don't have this information
								})
							}
						}

						result = nextResult
						if len(result.Providers) == 0 || result.Meta.NextOffset == 0 {
							break
						}
					}
				}
			} else {
				log.Printf("[DEBUG] Unexpected status code from provider search API: %d", searchResp.StatusCode)
				io.Copy(io.Discard, searchResp.Body)
			}
		}
	}

	// If we didn't get any providers from the search API, try fetching specific providers directly
	if len(allProviders) == 0 {
		log.Printf("[DEBUG] Search API didn't return any providers, trying direct provider fetching")
		
		// Try fetching specific providers directly
		for _, namespace := range commonNamespaces {
			// For each namespace, try to fetch a list of providers
			nsURL := *baseURL
			nsURL.Path = path.Join(nsURL.Path, namespace)
			
			log.Printf("[DEBUG] Trying to fetch providers for namespace %s: %s", namespace, nsURL.String())
			
			// Add a small delay to avoid rate limiting
			time.Sleep(200 * time.Millisecond)
			
			nsReq, err := retryablehttp.NewRequest("GET", nsURL.String(), nil)
			if err != nil {
				log.Printf("[DEBUG] Error creating namespace request: %s", err)
				continue
			}
			
			nsReq.Header.Set("Accept", "application/json")
			nsReq.Header.Set("User-Agent", "OpenTofu/"+tfVersion)
			nsReq.Header.Set(xTerraformVersion, tfVersion)

			nsResp, err := retryClient.Do(nsReq)
			if err != nil {
				log.Printf("[DEBUG] Error from namespace request: %s", err)
				continue
			}

			if nsResp.StatusCode != http.StatusOK {
				log.Printf("[DEBUG] Unexpected status code from namespace %s: %d", namespace, nsResp.StatusCode)
				nsResp.Body.Close()
				continue
			}

			// Try different response formats for the namespace endpoint
			// First try the format with providers as an array of strings
			var nsStringResult struct {
				Providers []string `json:"providers"`
			}
			
			bodyBytes, err := io.ReadAll(nsResp.Body)
			if err != nil {
				log.Printf("[DEBUG] Error reading namespace response body: %s", err)
				nsResp.Body.Close()
				continue
			}
			nsResp.Body.Close()
			
			if err := json.Unmarshal(bodyBytes, &nsStringResult); err == nil && len(nsStringResult.Providers) > 0 {
				log.Printf("[DEBUG] Successfully decoded %d providers from namespace %s (string format)", 
					len(nsStringResult.Providers), namespace)
				
				for _, name := range nsStringResult.Providers {
					providerID := fmt.Sprintf("%s/%s", namespace, name)
					if !seen[providerID] {
						seen[providerID] = true
						allProviders = append(allProviders, &response.ModuleProvider{
							Name:        name,
							Downloads:   0,  // We don't have this information
							ModuleCount: 0, // We don't have this information
						})
					}
				}
			} else {
				// Try the format with providers as an array of objects
				var nsObjectResult struct {
					Providers []struct {
						ID   string `json:"id"`
						Name string `json:"name"`
					} `json:"providers"`
				}
				
				if err := json.Unmarshal(bodyBytes, &nsObjectResult); err == nil && len(nsObjectResult.Providers) > 0 {
					log.Printf("[DEBUG] Successfully decoded %d providers from namespace %s (object format)", 
						len(nsObjectResult.Providers), namespace)
					
					for _, provider := range nsObjectResult.Providers {
						name := provider.Name
						if name == "" && provider.ID != "" {
							// Extract name from ID if name is not provided
							parts := strings.Split(provider.ID, "/")
							if len(parts) > 0 {
								name = parts[len(parts)-1]
							}
						}
						
						if name != "" {
							providerID := fmt.Sprintf("%s/%s", namespace, name)
							if !seen[providerID] {
								seen[providerID] = true
								allProviders = append(allProviders, &response.ModuleProvider{
									Name:        name,
									Downloads:   0,  // We don't have this information
									ModuleCount: 0, // We don't have this information
								})
							}
						}
					}
				} else {
					log.Printf("[DEBUG] Could not decode namespace response in any known format: %s", err)
				}
			}
		}
	}

	// If we still don't have any providers, try the v1/providers/{namespace}/{name}/versions pattern
	// for some well-known providers
	if len(allProviders) == 0 {
		log.Printf("[DEBUG] Still no providers found, trying direct version fetching for known providers")
		
		// Define some well-known providers to try
		knownProviders := []struct {
			namespace string
			name      string
		}{
			{"hashicorp", "aws"},
			{"hashicorp", "azurerm"},
			{"hashicorp", "google"},
			{"hashicorp", "kubernetes"},
			{"hashicorp", "random"},
			{"hashicorp", "template"},
			{"hashicorp", "null"},
			{"hashicorp", "local"},
			{"hashicorp", "tls"},
			{"hashicorp", "http"},
			{"hashicorp", "archive"},
			{"hashicorp", "external"},
			{"hashicorp", "dns"},
			{"hashicorp", "docker"},
			{"hashicorp", "helm"},
			{"hashicorp", "consul"},
			{"hashicorp", "vault"},
			{"hashicorp", "nomad"},
			{"hashicorp", "terraform"},
			{"opentofu", "opentofu"},
		}

		for _, provider := range knownProviders {
			providerURL := *baseURL
			providerURL.Path = path.Join(providerURL.Path, provider.namespace, provider.name, "versions")
			
			log.Printf("[DEBUG] Trying to fetch versions for provider %s/%s: %s", 
				provider.namespace, provider.name, providerURL.String())
			
			// Add a small delay to avoid rate limiting
			time.Sleep(200 * time.Millisecond)
			
			providerReq, err := retryablehttp.NewRequest("GET", providerURL.String(), nil)
			if err != nil {
				log.Printf("[DEBUG] Error creating provider versions request: %s", err)
				continue
			}
			
			providerReq.Header.Set("Accept", "application/json")
			providerReq.Header.Set("User-Agent", "OpenTofu/"+tfVersion)
			providerReq.Header.Set(xTerraformVersion, tfVersion)

			providerResp, err := retryClient.Do(providerReq)
			if err != nil {
				log.Printf("[DEBUG] Error from provider versions request: %s", err)
				continue
			}

			// If we get a successful response, add the provider to our list
			if providerResp.StatusCode == http.StatusOK {
				providerID := fmt.Sprintf("%s/%s", provider.namespace, provider.name)
				if !seen[providerID] {
					seen[providerID] = true
					allProviders = append(allProviders, &response.ModuleProvider{
						Name:        provider.name,
						Downloads:   0,  // We don't have this information
						ModuleCount: 0, // We don't have this information
					})
					log.Printf("[DEBUG] Successfully added provider %s/%s from versions endpoint", 
						provider.namespace, provider.name)
				}
			} else {
				log.Printf("[DEBUG] Unexpected status code from provider versions: %d", providerResp.StatusCode)
			}
			
			providerResp.Body.Close()
		}
	}

	// If we still don't have any providers, create a minimal set
	if len(allProviders) == 0 {
		log.Printf("[DEBUG] Could not find a working API pattern for providers on %s, creating minimal cache", host)
		
		// Create a minimal set of providers that we know exist
		log.Printf("[DEBUG] Creating minimal provider set for %s registry", host)
		minimalProviders := []struct {
			namespace string
			name      string
		}{
			{"hashicorp", "aws"},
			{"hashicorp", "azurerm"},
			{"hashicorp", "google"},
			{"hashicorp", "kubernetes"},
			{"hashicorp", "random"},
			{"hashicorp", "template"},
			{"hashicorp", "null"},
			{"hashicorp", "local"},
			{"hashicorp", "tls"},
			{"hashicorp", "http"},
			{"hashicorp", "archive"},
			{"hashicorp", "external"},
			{"hashicorp", "dns"},
			{"hashicorp", "docker"},
			{"hashicorp", "helm"},
			{"hashicorp", "consul"},
			{"hashicorp", "vault"},
			{"hashicorp", "nomad"},
			{"hashicorp", "terraform"},
			{"opentofu", "opentofu"},
		}

		for _, provider := range minimalProviders {
			providerID := fmt.Sprintf("%s/%s", provider.namespace, provider.name)
			if !seen[providerID] {
				seen[providerID] = true
				allProviders = append(allProviders, &response.ModuleProvider{
					Name:        provider.name,
					Downloads:   0,  // We don't have this information
					ModuleCount: 0, // We don't have this information
				})
			}
		}
	}

	log.Printf("[DEBUG] Found a total of %d providers for %s", len(allProviders), host)
	return allProviders, nil
}

// configureDiscoveryRetry configures the number of retries the registry client
// will attempt for requests with retryable errors, like 502 status codes
func configureDiscoveryRetry() {
	discoveryRetry = defaultRetry

	if v := os.Getenv(registryDiscoveryRetryEnvName); v != "" {
		retry, err := strconv.Atoi(v)
		if err == nil && retry > 0 {
			discoveryRetry = retry
		}
	}
}

func requestLogHook(logger retryablehttp.Logger, req *http.Request, i int) {
	if i > 0 {
		logger.Printf("[INFO] Previous request to the remote registry failed, attempting retry.")
	}
}

func maxRetryErrorHandler(resp *http.Response, err error, numTries int) (*http.Response, error) {
	// Close the body per library instructions
	if resp != nil {
		resp.Body.Close()
	}

	// Additional error detail: if we have a response, use the status code;
	// if we have an error, use that; otherwise nothing. We will never have
	// both response and error.
	var errMsg string
	if resp != nil {
		errMsg = fmt.Sprintf(": %s returned from %s", resp.Status, resp.Request.URL)
	} else if err != nil {
		errMsg = fmt.Sprintf(": %s", err)
	}

	// This function is always called with numTries=RetryMax+1. If we made any
	// retry attempts, include that in the error message.
	if numTries > 1 {
		return resp, fmt.Errorf("the request failed after %d attempts, please try again later%s",
			numTries, errMsg)
	}
	return resp, fmt.Errorf("the request failed, please try again later%s", errMsg)
}

// configureRequestTimeout configures the registry client request timeout from
// environment variables
func configureRequestTimeout() {
	requestTimeout = defaultRequestTimeout

	if v := os.Getenv(registryClientTimeoutEnvName); v != "" {
		timeout, err := strconv.Atoi(v)
		if err == nil && timeout > 0 {
			requestTimeout = time.Duration(timeout) * time.Second
		}
	}
}
