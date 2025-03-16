// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/hashicorp/terraform-svchost/disco"
	"github.com/opentofu/opentofu/internal/command"
	"github.com/opentofu/opentofu/internal/registry"
	"github.com/opentofu/opentofu/internal/registry/regsrc"
	"github.com/opentofu/opentofu/internal/registry/response"
	"github.com/opentofu/opentofu/internal/tfdiags"
)

// RegistrySearchCommand is a Command implementation that searches the registry
// for modules or providers.
type RegistrySearchCommand struct {
	Meta command.Meta
}

func (c *RegistrySearchCommand) Help() string {
	helpText := `
Usage: tofu registry search [options] [QUERY]

  Search the registry for modules or providers.

Options:

  -type=TYPE            Type of resource to search for. Can be "module" or "provider".
                        Default: "module"

  -limit=N              Limit the number of results displayed. Default: 10

  -json                 Output the results as JSON.

  -detailed             Show detailed information about each result.

  -registry=hostname    Use a custom registry host. By default, public registry
                        hosts are used based on the resource type.
`
	return strings.TrimSpace(helpText)
}

func (c *RegistrySearchCommand) Synopsis() string {
	return "Search the registry for modules or providers"
}

func (c *RegistrySearchCommand) Run(args []string) int {
	var searchType string
	var limit int
	var jsonOutput bool
	var detailed bool
	var registryHost string

	cmdFlags := flag.NewFlagSet("registry search", flag.ContinueOnError)
	cmdFlags.StringVar(&searchType, "type", "module", "Type of resource to search for")
	cmdFlags.IntVar(&limit, "limit", 10, "Limit the number of results")
	cmdFlags.BoolVar(&jsonOutput, "json", false, "Output the results as JSON")
	cmdFlags.BoolVar(&detailed, "detailed", false, "Show detailed information")
	cmdFlags.StringVar(&registryHost, "registry", "", "Registry host")
	cmdFlags.Usage = func() { c.Meta.Ui.Error(c.Help()) }

	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}

	args = cmdFlags.Args()
	if len(args) > 1 {
		c.Meta.Ui.Error("The registry search command expects at most one argument.")
		return 1
	}

	var query string
	if len(args) == 1 {
		query = args[0]
	}

	// Check if the search type is valid
	if searchType != "module" && searchType != "provider" {
		c.Meta.Ui.Error(fmt.Sprintf("Invalid search type: %s. Must be 'module' or 'provider'.", searchType))
		return 1
	}

	// Create a registry client
	client, err := c.createRegistryClient()
	if err != nil {
		c.Meta.Ui.Error(fmt.Sprintf("Error initializing registry client: %s", err))
		return 1
	}

	// If a registry host is specified, use it
	var host *regsrc.FriendlyHost
	if registryHost != "" {
		host = regsrc.NewFriendlyHost(registryHost)
	} else {
		// Use the default registry host
		if searchType == "module" {
			host = regsrc.PublicRegistryHost
		} else {
			host = regsrc.PublicRegistryHost // Same for providers
		}
	}

	// Refresh the registry cache if needed
	ctx := context.Background()
	diags := c.refreshRegistryIfNeeded(ctx, client, host)
	if diags.HasErrors() {
		c.Meta.Ui.Error(fmt.Sprintf("Error refreshing registry: %s", diags.Err()))
		return 1
	}

	// Perform the search
	if searchType == "module" {
		return c.searchModules(ctx, client, host, query, limit, jsonOutput, detailed)
	} else {
		return c.searchProviders(ctx, client, host, query, limit, jsonOutput, detailed)
	}
}

// createRegistryClient creates a new registry client
func (c *RegistrySearchCommand) createRegistryClient() (*registry.Client, error) {
	// Create an HTTP client for the registry
	httpClient := retryablehttp.NewClient()
	httpClient.RetryMax = 3
	httpClient.Logger = hclog.NewNullLogger()

	// Create a services discovery client
	services := disco.New()
	services.SetUserAgent("OpenTofu")

	// Create a registry client
	client := registry.NewClient(services, httpClient.StandardClient())

	// Set up caching
	cacheDir := filepath.Join(os.TempDir(), "opentofu-registry-cache")
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create registry cache directory: %w", err)
	}

	return client, nil
}

func (c *RegistrySearchCommand) refreshRegistryIfNeeded(ctx context.Context, client *registry.Client, host *regsrc.FriendlyHost) tfdiags.Diagnostics {
	var diags tfdiags.Diagnostics

	// Check if we need to refresh the registry
	cacheFile := registry.CacheFilename("modules.json")
	if registry.ShouldRefreshCache(cacheFile, 24*time.Hour) {
		c.Meta.Ui.Output(fmt.Sprintf("Refreshing registry cache for %s...", host))
		err := client.RefreshModules(ctx, host.String())
		if err != nil {
			diags = diags.Append(fmt.Errorf("Failed to refresh modules for %s: %s", host, err))
		}
	}

	cacheFile = registry.CacheFilename("providers.json")
	if registry.ShouldRefreshCache(cacheFile, 24*time.Hour) {
		c.Meta.Ui.Output(fmt.Sprintf("Refreshing provider cache for %s...", host))
		err := client.RefreshProviders(ctx, host.String())
		if err != nil {
			diags = diags.Append(fmt.Errorf("Failed to refresh providers for %s: %s", host, err))
		}
	}

	return diags
}

func (c *RegistrySearchCommand) searchModules(ctx context.Context, client *registry.Client, host *regsrc.FriendlyHost, query string, limit int, jsonOutput, detailed bool) int {
	// Get all modules from the registry
	modules, err := client.GetModules(ctx, host.String())
	if err != nil {
		c.Meta.Ui.Error(fmt.Sprintf("Error fetching modules: %s", err))
		return 1
	}

	// Filter modules based on the query
	var filteredModules []*response.Module
	if query != "" {
		for _, module := range modules {
			// Search in namespace, name, provider, and description
			searchText := strings.ToLower(fmt.Sprintf("%s/%s/%s %s", 
				module.Namespace, module.Name, module.Provider, module.Description))
			if strings.Contains(searchText, strings.ToLower(query)) {
				filteredModules = append(filteredModules, module)
			}
		}
	} else {
		filteredModules = modules
	}

	// Sort modules by downloads (most popular first)
	sort.Slice(filteredModules, func(i, j int) bool {
		return filteredModules[i].Downloads > filteredModules[j].Downloads
	})

	// Limit the number of results
	if limit > 0 && limit < len(filteredModules) {
		filteredModules = filteredModules[:limit]
	}

	// Output the results
	if jsonOutput {
		return c.outputModulesAsJSON(filteredModules)
	} else {
		return c.outputModulesAsText(filteredModules, detailed)
	}
}

func (c *RegistrySearchCommand) searchProviders(ctx context.Context, client *registry.Client, host *regsrc.FriendlyHost, query string, limit int, jsonOutput, detailed bool) int {
	// Get all providers from the registry
	providers, err := client.GetProviders(ctx, host.String())
	if err != nil {
		c.Meta.Ui.Error(fmt.Sprintf("Error fetching providers: %s", err))
		return 1
	}

	// Filter providers based on the query
	var filteredProviders []*response.ModuleProvider
	if query != "" {
		for _, provider := range providers {
			// Search in name
			if strings.Contains(strings.ToLower(provider.Name), strings.ToLower(query)) {
				filteredProviders = append(filteredProviders, provider)
			}
		}
	} else {
		filteredProviders = providers
	}

	// Sort providers by downloads (most popular first)
	sort.Slice(filteredProviders, func(i, j int) bool {
		return filteredProviders[i].Downloads > filteredProviders[j].Downloads
	})

	// Limit the number of results
	if limit > 0 && limit < len(filteredProviders) {
		filteredProviders = filteredProviders[:limit]
	}

	// Output the results
	if jsonOutput {
		return c.outputProvidersAsJSON(filteredProviders)
	} else {
		return c.outputProvidersAsText(filteredProviders, detailed)
	}
}

func (c *RegistrySearchCommand) outputModulesAsJSON(modules []*response.Module) int {
	// Output the results as JSON
	jsonOutput := map[string]interface{}{
		"modules": modules,
	}
	c.Meta.Ui.Output(fmt.Sprintf("%v", jsonOutput))
	return 0
}

func (c *RegistrySearchCommand) outputModulesAsText(modules []*response.Module, detailed bool) int {
	if len(modules) == 0 {
		c.Meta.Ui.Output("No matching modules found.")
		return 0
	}

	c.Meta.Ui.Output(fmt.Sprintf("Found %d matching modules:", len(modules)))
	c.Meta.Ui.Output("")

	for _, module := range modules {
		c.Meta.Ui.Output(fmt.Sprintf("* %s/%s/%s", module.Namespace, module.Name, module.Provider))
		
		if detailed {
			if module.Description != "" {
				c.Meta.Ui.Output(fmt.Sprintf("    Description: %s", module.Description))
			}
			c.Meta.Ui.Output(fmt.Sprintf("    Downloads: %d", module.Downloads))
			c.Meta.Ui.Output(fmt.Sprintf("    Verified: %t", module.Verified))
			if len(module.Version) > 0 {
				c.Meta.Ui.Output(fmt.Sprintf("    Latest Version: %s", module.Version))
			}
			c.Meta.Ui.Output("")
		}
	}

	return 0
}

func (c *RegistrySearchCommand) outputProvidersAsJSON(providers []*response.ModuleProvider) int {
	// Output the results as JSON
	jsonOutput := map[string]interface{}{
		"providers": providers,
	}
	c.Meta.Ui.Output(fmt.Sprintf("%v", jsonOutput))
	return 0
}

func (c *RegistrySearchCommand) outputProvidersAsText(providers []*response.ModuleProvider, detailed bool) int {
	if len(providers) == 0 {
		c.Meta.Ui.Output("No matching providers found.")
		return 0
	}

	c.Meta.Ui.Output(fmt.Sprintf("Found %d matching providers:", len(providers)))
	c.Meta.Ui.Output("")

	for _, provider := range providers {
		c.Meta.Ui.Output(fmt.Sprintf("* %s", provider.Name))
		
		if detailed {
			c.Meta.Ui.Output(fmt.Sprintf("    Downloads: %d", provider.Downloads))
			c.Meta.Ui.Output(fmt.Sprintf("    Module Count: %d", provider.ModuleCount))
			c.Meta.Ui.Output("")
		}
	}

	return 0
}
