// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2023 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package registry

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/go-hclog"
	svchost "github.com/hashicorp/terraform-svchost"

	"github.com/opentofu/opentofu/internal/registry/response"
	"github.com/opentofu/opentofu/version"
)

// Constants for service discovery
const (
// Service discovery constants are already defined in client.go
)

// RegistryCachingError wraps errors that occur during registry caching operations
type RegistryCachingError struct {
	err error
}

func (e *RegistryCachingError) Error() string {
	return e.err.Error()
}

func (e *RegistryCachingError) Unwrap() error {
	return e.err
}

// CachingClient wraps a registry Client and adds metadata caching capabilities.
type CachingClient struct {
	*Client
	cacheDir     string
	logger       hclog.Logger
	refreshMutex sync.Mutex
}

// NewCachingClient creates a new CachingClient.
func NewCachingClient(client *Client, cacheDir string, logger hclog.Logger) (*CachingClient, error) {
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create cache directory: %w", err)
	}

	return &CachingClient{
		Client:   client,
		cacheDir: cacheDir,
		logger:   logger,
	}, nil
}

// CachedMetadata represents cached registry metadata with timestamp information.
type CachedMetadata struct {
	Timestamp time.Time   `json:"timestamp"`
	Host      string      `json:"host"`
	Type      string      `json:"type"`
	Modules   interface{} `json:"modules,omitempty"`
	Providers interface{} `json:"providers,omitempty"`
}

// RefreshModuleCache fetches and caches all available modules for a given host.
func (c *CachingClient) RefreshModuleCache(ctx context.Context, host svchost.Hostname) error {
	c.logger.Debug("Refreshing module cache", "host", host.String())

	modules, err := c.BulkFetchModules(ctx, host)
	if err != nil {
		return fmt.Errorf("failed to fetch modules: %w", err)
	}

	metadata := &CachedMetadata{
		Timestamp: time.Now(),
		Host:      host.String(),
		Type:      "modules",
		Modules:   modules,
	}

	return c.saveToCache(metadata)
}

// RefreshProviderCache fetches and caches all available providers for a given host.
func (c *CachingClient) RefreshProviderCache(ctx context.Context, host svchost.Hostname) error {
	c.logger.Debug("Refreshing provider cache", "host", host.String())

	providers, err := c.BulkFetchProviders(ctx, host)
	if err != nil {
		return fmt.Errorf("failed to fetch providers: %w", err)
	}

	metadata := &CachedMetadata{
		Timestamp: time.Now(),
		Host:      host.String(),
		Type:      "providers",
		Providers: providers,
	}

	return c.saveToCache(metadata)
}

// SaveModulesToCache saves a list of modules to the cache for a given host.
func (c *CachingClient) SaveModulesToCache(host svchost.Hostname, modules []*response.Module) error {
	c.logger.Debug("Saving modules to cache", "host", host.String(), "count", len(modules))

	metadata := &CachedMetadata{
		Timestamp: time.Now(),
		Host:      host.String(),
		Type:      "modules",
		Modules:   modules,
	}

	return c.saveToCache(metadata)
}

// SaveProvidersToCache saves a list of providers to the cache for a given host.
func (c *CachingClient) SaveProvidersToCache(host svchost.Hostname, providers []*response.ModuleProvider) error {
	c.logger.Debug("Saving providers to cache", "host", host.String(), "count", len(providers))

	metadata := &CachedMetadata{
		Timestamp: time.Now(),
		Host:      host.String(),
		Type:      "providers",
		Providers: providers,
	}

	return c.saveToCache(metadata)
}

// saveToCache serializes and saves metadata to a cache file.
func (c *CachingClient) saveToCache(metadata *CachedMetadata) error {
	c.refreshMutex.Lock()
	defer c.refreshMutex.Unlock()

	// Create host-specific directory
	hostDir := filepath.Join(c.cacheDir, metadata.Host)
	if err := os.MkdirAll(hostDir, 0755); err != nil {
		return fmt.Errorf("failed to create host directory: %w", err)
	}

	// Create timestamp-based filename
	timestamp := metadata.Timestamp.Format("20060102-150405")
	filename := fmt.Sprintf("%s-%s.json", metadata.Type, timestamp)
	filePath := filepath.Join(hostDir, filename)

	// Serialize metadata to JSON
	data, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to serialize metadata: %w", err)
	}

	// Write to file
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write cache file: %w", err)
	}

	// Update symlink to latest cache file
	symlinkPath := filepath.Join(hostDir, fmt.Sprintf("%s-latest.json", metadata.Type))
	_ = os.Remove(symlinkPath) // Ignore error if symlink doesn't exist
	if err := os.Symlink(filename, symlinkPath); err != nil {
		return fmt.Errorf("failed to create symlink to latest cache file: %w", err)
	}

	c.logger.Debug("Saved metadata to cache", "host", metadata.Host, "type", metadata.Type, "file", filePath)
	return nil
}

// CleanupOldCacheFiles removes cache files older than maxAge.
func (c *CachingClient) CleanupOldCacheFiles(maxAge time.Duration) error {
	c.logger.Debug("Cleaning up old cache files", "maxAge", maxAge)

	// Get the current time
	now := time.Now()

	// Walk through the cache directory
	return filepath.Walk(c.cacheDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories and symlinks
		if info.IsDir() || (info.Mode()&os.ModeSymlink != 0) {
			return nil
		}

		// Skip files that don't match our naming pattern
		if !strings.HasSuffix(info.Name(), ".json") || strings.HasSuffix(info.Name(), "-latest.json") {
			return nil
		}

		// Check if the file is older than maxAge
		if now.Sub(info.ModTime()) > maxAge {
			c.logger.Debug("Removing old cache file", "file", path)
			return os.Remove(path)
		}

		return nil
	})
}

// StartBackgroundRefresh starts background goroutines to periodically refresh the registry cache.
func (c *CachingClient) StartBackgroundRefresh(ctx context.Context, hosts []svchost.Hostname, refreshInterval, cleanupInterval, maxAge time.Duration) {
	c.logger.Debug("Starting background registry refresh",
		"refreshInterval", refreshInterval,
		"cleanupInterval", cleanupInterval,
		"maxAge", maxAge)

	// Start a goroutine to periodically refresh the registry cache
	go func() {
		ticker := time.NewTicker(refreshInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				c.logger.Debug("Background registry refresh stopped")
				return
			case <-ticker.C:
				c.logger.Debug("Running scheduled registry refresh")
				for _, host := range hosts {
					// Refresh modules
					if err := c.RefreshModuleCache(ctx, host); err != nil {
						c.logger.Error("Failed to refresh module cache",
							"host", host.String(),
							"error", err)
					}

					// Refresh providers
					if err := c.RefreshProviderCache(ctx, host); err != nil {
						c.logger.Error("Failed to refresh provider cache",
							"host", host.String(),
							"error", err)
					}
				}
			}
		}
	}()

	// Start a goroutine to periodically clean up old cache files
	go func() {
		ticker := time.NewTicker(cleanupInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				c.logger.Debug("Running scheduled cache cleanup")
				if err := c.CleanupOldCacheFiles(maxAge); err != nil {
					c.logger.Error("Failed to clean up old cache files", "error", err)
				}
			}
		}
	}()
}

// Constants for registry types
const (
	TerraformRegistryType = "terraform"
	OpenTofuRegistryType  = "opentofu"
	ModuleType            = "modules"
	ProviderType          = "providers"
)

// RegistryCache represents cached registry metadata with timestamp information.
type RegistryCache struct {
	Timestamp time.Time                  `json:"timestamp"`
	Host      string                     `json:"host"`
	Type      string                     `json:"type"`
	Modules   []*response.Module         `json:"modules,omitempty"`
	Providers []*response.ModuleProvider `json:"providers,omitempty"`
}

// RegistryCachingClient wraps a registry Client and adds metadata caching capabilities
// with special handling for registry-specific API calls.
type RegistryCachingClient struct {
	*CachingClient
	rateLimitDelay time.Duration // Delay between API calls to avoid rate limiting
}

// NewRegistryCachingClient creates a new RegistryCachingClient.
func NewRegistryCachingClient(client *Client, cacheDir string, logger hclog.Logger) (*RegistryCachingClient, error) {
	cachingClient, err := NewCachingClient(client, cacheDir, logger)
	if err != nil {
		return nil, err
	}

	return &RegistryCachingClient{
		CachingClient:  cachingClient,
		rateLimitDelay: 200 * time.Millisecond, // Default delay between API calls
	}, nil
}

// SetRateLimitDelay sets the delay between API calls to avoid rate limiting.
func (c *RegistryCachingClient) SetRateLimitDelay(delay time.Duration) {
	c.rateLimitDelay = delay
}

// saveToCache serializes and saves registry cache to a file.
func (c *RegistryCachingClient) saveToCache(cache *RegistryCache, registryType string) error {
	c.refreshMutex.Lock()
	defer c.refreshMutex.Unlock()

	// Create host-specific directory
	hostDir := filepath.Join(c.cacheDir, cache.Host, registryType)
	if err := os.MkdirAll(hostDir, 0755); err != nil {
		return fmt.Errorf("failed to create host directory: %w", err)
	}

	// Create timestamp-based filename
	timestamp := cache.Timestamp.Format("20060102-150405")
	filename := fmt.Sprintf("%s-%s.json", cache.Type, timestamp)
	filePath := filepath.Join(hostDir, filename)

	// Serialize metadata to JSON
	data, err := json.MarshalIndent(cache, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to serialize metadata: %w", err)
	}

	// Write to file
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write cache file: %w", err)
	}

	// Update symlink to latest cache file
	symlinkPath := filepath.Join(hostDir, fmt.Sprintf("%s-latest.json", cache.Type))
	_ = os.Remove(symlinkPath) // Ignore error if symlink doesn't exist
	if err := os.Symlink(filename, symlinkPath); err != nil {
		return fmt.Errorf("failed to create symlink to latest cache file: %w", err)
	}

	c.logger.Debug("Saved registry cache", "host", cache.Host, "type", cache.Type, "file", filePath)
	return nil
}

// RefreshModuleCache fetches and caches all available modules for a given host
func (c *RegistryCachingClient) RefreshModuleCache(ctx context.Context, host svchost.Hostname) error {
	c.logger.Debug("Refreshing module cache", "host", host.String())

	// First try the standard client method
	modules, err := c.Client.BulkFetchModules(ctx, host)
	if err == nil && len(modules) > 0 {
		c.logger.Info("Successfully fetched modules using standard client", "host", host.String(), "count", len(modules))

		// Create the cache entry
		registryType := TerraformRegistryType
		if host.String() == "registry.opentofu.org" {
			registryType = OpenTofuRegistryType
		}

		cache := &RegistryCache{
			Timestamp: time.Now(),
			Host:      host.String(),
			Type:      ModuleType,
			Modules:   modules,
		}

		return c.saveToCache(cache, registryType)
	}

	c.logger.Debug("Standard module fetch failed, trying alternative approach", "error", err)

	// Based on the existing client implementation, build a direct URL
	service, err := c.Client.Discover(host, modulesServiceID)
	if err != nil {
		return &RegistryCachingError{fmt.Errorf("failed to discover modules service: %w", err)}
	}

	// Pre-allocate based on known registry sizes (approximately 18,000 modules)
	allModules := make([]*response.Module, 0, 18000)

	// Try multiple API patterns that might work
	apiPatterns := []string{
		"modules",
		"v1/modules",
	}

	var successfulURL string
	var moduleList response.ModuleList

	for _, pattern := range apiPatterns {
		// Create URL with proper query parameters
		p, err := url.Parse(pattern)
		if err != nil {
			c.logger.Debug("Failed to parse API pattern", "pattern", pattern, "error", err)
			continue
		}

		listURL := service.ResolveReference(p)
		queryParams := listURL.Query()
		queryParams.Set("limit", "100") // Max page size
		listURL.RawQuery = queryParams.Encode()

		c.logger.Debug("Trying modules API pattern", "url", listURL.String())

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, listURL.String(), nil)
		if err != nil {
			c.logger.Debug("Failed to create request", "url", listURL.String(), "error", err)
			continue
		}

		// Set Terraform version header
		req.Header.Set(xTerraformVersion, version.Version)

		// Add authentication if needed
		c.Client.addRequestCreds(host, req)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			c.logger.Debug("Failed to fetch modules", "url", listURL.String(), "error", err)
			continue
		}

		if resp.StatusCode != http.StatusOK {
			c.logger.Debug("Received non-OK status", "url", listURL.String(), "status", resp.StatusCode)
			resp.Body.Close()
			continue
		}

		if err := json.NewDecoder(resp.Body).Decode(&moduleList); err != nil {
			c.logger.Debug("Failed to decode response", "url", listURL.String(), "error", err)
			resp.Body.Close()
			continue
		}
		resp.Body.Close()

		// If we get here, we've found a working URL pattern
		successfulURL = listURL.String()
		allModules = append(allModules, moduleList.Modules...)
		c.logger.Info("Found working modules API pattern", "url", listURL.String())
		break
	}

	if successfulURL == "" {
		return &RegistryCachingError{fmt.Errorf("failed to find working modules API pattern for host %s", host.String())}
	}

	// Process pagination using the successful URL pattern as a base
	nextURL := moduleList.Meta.NextURL

	// Handle pagination if needed
	for nextURL != "" {
		// If the URL is relative, resolve it against the service URL
		var nextFullURL string
		if strings.HasPrefix(nextURL, "http") {
			nextFullURL = nextURL
		} else {
			nextURLObj, err := url.Parse(nextURL)
			if err != nil {
				c.logger.Warn("Failed to parse next URL", "url", nextURL, "error", err)
				break
			}

			nextFullURL = service.ResolveReference(nextURLObj).String()
		}

		// Add throttling delay to avoid rate limiting
		time.Sleep(c.rateLimitDelay)

		c.logger.Debug("Fetching next page", "url", nextFullURL)

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, nextFullURL, nil)
		if err != nil {
			return &RegistryCachingError{fmt.Errorf("failed to create request for next page: %w", err)}
		}

		// Set Terraform version header
		req.Header.Set(xTerraformVersion, version.Version)

		// Add authentication if needed
		c.Client.addRequestCreds(host, req)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return &RegistryCachingError{fmt.Errorf("failed to fetch next page: %w", err)}
		}

		var nextList response.ModuleList
		if err := json.NewDecoder(resp.Body).Decode(&nextList); err != nil {
			resp.Body.Close()
			return &RegistryCachingError{fmt.Errorf("failed to decode next page response: %w", err)}
		}
		resp.Body.Close()

		allModules = append(allModules, nextList.Modules...)

		// Get the next URL for pagination
		nextURL = nextList.Meta.NextURL
	}

	c.logger.Info("Fetched modules", "host", host.String(), "count", len(allModules))

	// Create the cache entry
	registryType := TerraformRegistryType
	if host.String() == "registry.opentofu.org" {
		registryType = OpenTofuRegistryType
	}

	cache := &RegistryCache{
		Timestamp: time.Now(),
		Host:      host.String(),
		Type:      ModuleType,
		Modules:   allModules,
	}

	return c.saveToCache(cache, registryType)
}

// RefreshProviderCache fetches and caches all available providers for a given host
func (c *RegistryCachingClient) RefreshProviderCache(ctx context.Context, host svchost.Hostname) error {
	c.logger.Debug("Refreshing provider cache", "host", host.String())

	// First try the standard client method
	providers, err := c.Client.BulkFetchProviders(ctx, host)
	if err == nil && len(providers) > 0 {
		c.logger.Info("Successfully fetched providers using standard client", "host", host.String(), "count", len(providers))

		// Create the cache entry
		registryType := TerraformRegistryType
		if host.String() == "registry.opentofu.org" {
			registryType = OpenTofuRegistryType
		}

		cache := &RegistryCache{
			Timestamp: time.Now(),
			Host:      host.String(),
			Type:      ProviderType,
			Providers: providers,
		}

		return c.saveToCache(cache, registryType)
	}

	c.logger.Debug("Standard provider fetch failed, trying alternative approach", "error", err)

	// Based on the existing client implementation, build a direct URL
	service, err := c.Client.Discover(host, providersServiceID)
	if err != nil {
		return &RegistryCachingError{fmt.Errorf("failed to discover providers service: %w", err)}
	}

	// Pre-allocate based on known registry sizes (approximately 4,000 providers)
	allProviders := make([]*response.ModuleProvider, 0, 4000)

	// Try multiple API patterns that might work
	apiPatterns := []string{
		"providers",
		"v1/providers",
	}

	var successfulURL string
	var providerList response.ModuleProviderList

	for _, pattern := range apiPatterns {
		// Create URL with proper query parameters
		p, err := url.Parse(pattern)
		if err != nil {
			c.logger.Debug("Failed to parse API pattern", "pattern", pattern, "error", err)
			continue
		}

		listURL := service.ResolveReference(p)
		queryParams := listURL.Query()
		queryParams.Set("limit", "100") // Max page size
		listURL.RawQuery = queryParams.Encode()

		c.logger.Debug("Trying providers API pattern", "url", listURL.String())

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, listURL.String(), nil)
		if err != nil {
			c.logger.Debug("Failed to create request", "url", listURL.String(), "error", err)
			continue
		}

		// Set Terraform version header
		req.Header.Set(xTerraformVersion, version.Version)

		// Add authentication if needed
		c.Client.addRequestCreds(host, req)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			c.logger.Debug("Failed to fetch providers", "url", listURL.String(), "error", err)
			continue
		}

		if resp.StatusCode != http.StatusOK {
			c.logger.Debug("Received non-OK status", "url", listURL.String(), "status", resp.StatusCode)
			resp.Body.Close()
			continue
		}

		if err := json.NewDecoder(resp.Body).Decode(&providerList); err != nil {
			c.logger.Debug("Failed to decode response", "url", listURL.String(), "error", err)
			resp.Body.Close()
			continue
		}
		resp.Body.Close()

		// If we get here, we've found a working URL pattern
		successfulURL = listURL.String()
		allProviders = append(allProviders, providerList.Providers...)
		c.logger.Info("Found working providers API pattern", "url", listURL.String())
		break
	}

	if successfulURL == "" {
		return &RegistryCachingError{fmt.Errorf("failed to find working providers API pattern for host %s", host.String())}
	}

	// Process pagination using the successful URL pattern as a base
	nextURL := providerList.Meta.NextURL

	// Handle pagination if needed
	for nextURL != "" {
		// If the URL is relative, resolve it against the service URL
		var nextFullURL string
		if strings.HasPrefix(nextURL, "http") {
			nextFullURL = nextURL
		} else {
			nextURLObj, err := url.Parse(nextURL)
			if err != nil {
				c.logger.Warn("Failed to parse next URL", "url", nextURL, "error", err)
				break
			}

			nextFullURL = service.ResolveReference(nextURLObj).String()
		}

		// Add throttling delay to avoid rate limiting
		time.Sleep(c.rateLimitDelay)

		c.logger.Debug("Fetching next page", "url", nextFullURL)

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, nextFullURL, nil)
		if err != nil {
			return &RegistryCachingError{fmt.Errorf("failed to create request for next page: %w", err)}
		}

		// Set Terraform version header
		req.Header.Set(xTerraformVersion, version.Version)

		// Add authentication if needed
		c.Client.addRequestCreds(host, req)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return &RegistryCachingError{fmt.Errorf("failed to fetch next page: %w", err)}
		}

		var nextList response.ModuleProviderList
		if err := json.NewDecoder(resp.Body).Decode(&nextList); err != nil {
			resp.Body.Close()
			return &RegistryCachingError{fmt.Errorf("failed to decode next page response: %w", err)}
		}
		resp.Body.Close()

		allProviders = append(allProviders, nextList.Providers...)

		// Get the next URL for pagination
		nextURL = nextList.Meta.NextURL
	}

	c.logger.Info("Fetched providers", "host", host.String(), "count", len(allProviders))

	// Create the cache entry
	registryType := TerraformRegistryType
	if host.String() == "registry.opentofu.org" {
		registryType = OpenTofuRegistryType
	}

	cache := &RegistryCache{
		Timestamp: time.Now(),
		Host:      host.String(),
		Type:      ProviderType,
		Providers: allProviders,
	}

	return c.saveToCache(cache, registryType)
}

// CacheFilename returns the filename for a specific cache type
func CacheFilename(cacheType string) string {
	return cacheType
}

// ShouldRefreshCache checks if the cache file is older than maxAge
func ShouldRefreshCache(cacheFile string, maxAge time.Duration) bool {
	info, err := os.Stat(cacheFile)
	if err != nil {
		// If the file doesn't exist or can't be accessed, refresh the cache
		return true
	}
	
	// Check if the file is older than maxAge
	return time.Since(info.ModTime()) > maxAge
}

// GetModulesFromCache retrieves modules from the cache for a given host
func (c *RegistryCachingClient) GetModulesFromCache(host svchost.Hostname) ([]*response.Module, error) {
	cacheFile := filepath.Join(c.cacheDir, fmt.Sprintf("modules_%s.json", host.String()))
	
	data, err := os.ReadFile(cacheFile)
	if err != nil {
		return nil, &RegistryCachingError{err: fmt.Errorf("failed to read module cache file: %w", err)}
	}
	
	var cache RegistryCache
	if err := json.Unmarshal(data, &cache); err != nil {
		return nil, &RegistryCachingError{err: fmt.Errorf("failed to unmarshal module cache: %w", err)}
	}
	
	return cache.Modules, nil
}

// GetProvidersFromCache retrieves providers from the cache for a given host
func (c *RegistryCachingClient) GetProvidersFromCache(host svchost.Hostname) ([]*response.ModuleProvider, error) {
	cacheFile := filepath.Join(c.cacheDir, fmt.Sprintf("providers_%s.json", host.String()))
	
	data, err := os.ReadFile(cacheFile)
	if err != nil {
		return nil, &RegistryCachingError{err: fmt.Errorf("failed to read provider cache file: %w", err)}
	}
	
	var cache RegistryCache
	if err := json.Unmarshal(data, &cache); err != nil {
		return nil, &RegistryCachingError{err: fmt.Errorf("failed to unmarshal provider cache: %w", err)}
	}
	
	return cache.Providers, nil
}

// RefreshModules refreshes the module cache for a given host
func (c *Client) RefreshModules(ctx context.Context, host string) error {
	hostname, err := svchost.ForComparison(host)
	if err != nil {
		return err
	}
	
	// Create a caching client
	cachingClient, err := NewRegistryCachingClient(c, os.TempDir(), hclog.New(&hclog.LoggerOptions{
		Name:   "registry-cache",
		Level:  hclog.Info,
		Output: os.Stderr,
	}))
	if err != nil {
		return err
	}
	
	// Refresh the module cache
	return cachingClient.RefreshModuleCache(ctx, hostname)
}

// RefreshProviders refreshes the provider cache for a given host
func (c *Client) RefreshProviders(ctx context.Context, host string) error {
	hostname, err := svchost.ForComparison(host)
	if err != nil {
		return err
	}
	
	// Create a caching client
	cachingClient, err := NewRegistryCachingClient(c, os.TempDir(), hclog.New(&hclog.LoggerOptions{
		Name:   "registry-cache",
		Level:  hclog.Info,
		Output: os.Stderr,
	}))
	if err != nil {
		return err
	}
	
	// Refresh the provider cache
	return cachingClient.RefreshProviderCache(ctx, hostname)
}

// GetModules retrieves modules from the cache for a given host
func (c *Client) GetModules(ctx context.Context, host string) ([]*response.Module, error) {
	hostname, err := svchost.ForComparison(host)
	if err != nil {
		return nil, err
	}
	
	// Create a caching client
	cachingClient, err := NewRegistryCachingClient(c, os.TempDir(), hclog.New(&hclog.LoggerOptions{
		Name:   "registry-cache",
		Level:  hclog.Info,
		Output: os.Stderr,
	}))
	if err != nil {
		return nil, err
	}
	
	// Get modules from the cache
	modules, err := cachingClient.GetModulesFromCache(hostname)
	if err != nil {
		// If there's an error getting modules from the cache, try refreshing
		if err := cachingClient.RefreshModuleCache(ctx, hostname); err != nil {
			return nil, err
		}
		
		// Try again after refreshing
		return cachingClient.GetModulesFromCache(hostname)
	}
	
	return modules, nil
}

// GetProviders retrieves providers from the cache for a given host
func (c *Client) GetProviders(ctx context.Context, host string) ([]*response.ModuleProvider, error) {
	hostname, err := svchost.ForComparison(host)
	if err != nil {
		return nil, err
	}
	
	// Create a caching client
	cachingClient, err := NewRegistryCachingClient(c, os.TempDir(), hclog.New(&hclog.LoggerOptions{
		Name:   "registry-cache",
		Level:  hclog.Info,
		Output: os.Stderr,
	}))
	if err != nil {
		return nil, err
	}
	
	// Get providers from the cache
	providers, err := cachingClient.GetProvidersFromCache(hostname)
	if err != nil {
		// If there's an error getting providers from the cache, try refreshing
		if err := cachingClient.RefreshProviderCache(ctx, hostname); err != nil {
			return nil, err
		}
		
		// Try again after refreshing
		return cachingClient.GetProvidersFromCache(hostname)
	}
	
	return providers, nil
}
