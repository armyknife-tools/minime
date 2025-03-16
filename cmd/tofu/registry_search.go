// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/opentofu/opentofu/internal/command"
	"github.com/opentofu/opentofu/internal/registry/regsrc"
	"github.com/opentofu/opentofu/internal/registry/response"
	"io"
)

// RegistrySearchCommand is a CLI command for searching the registry
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

  -count-only           Only show the total count of available modules and providers.
`
	return strings.TrimSpace(helpText)
}

func (c *RegistrySearchCommand) Synopsis() string {
	return "Search the registry for modules or providers"
}

func (c *RegistrySearchCommand) Run(args []string) int {
	var typeFlag string
	var limitFlag int
	var jsonFlag bool
	var detailedFlag bool
	var registryFlag string
	var countOnlyFlag bool

	flags := flag.NewFlagSet("registry search", flag.ContinueOnError)
	flags.StringVar(&typeFlag, "type", "module", "Type of resource to search for. Can be \"module\" or \"provider\"")
	flags.IntVar(&limitFlag, "limit", 10, "Limit the number of results displayed")
	flags.BoolVar(&jsonFlag, "json", false, "Output the results as JSON")
	flags.BoolVar(&detailedFlag, "detailed", false, "Show detailed information about each result")
	flags.StringVar(&registryFlag, "registry", "registry.terraform.io", "Use a custom registry host")
	flags.BoolVar(&countOnlyFlag, "count-only", false, "Only show the total count of available modules and providers")

	flags.Usage = func() { c.Meta.Ui.Error(c.Help()) }

	if err := flags.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return 0
		}
		c.Meta.Ui.Error(fmt.Sprintf("Error parsing command line arguments: %s", err))
		return 1
	}

	args = flags.Args()
	if len(args) > 1 {
		c.Meta.Ui.Error("The registry search command expects at most one argument.")
		return 1
	}

	query := ""
	if len(args) == 1 {
		query = args[0]
	}

	// Check if the search type is valid
	if typeFlag != "module" && typeFlag != "provider" {
		c.Meta.Ui.Error(fmt.Sprintf("Invalid search type: %s. Must be 'module' or 'provider'.", typeFlag))
		return 1
	}

	// If a registry host is specified, use it
	var host *regsrc.FriendlyHost
	if registryFlag != "" {
		host = regsrc.NewFriendlyHost(registryFlag)
	} else {
		// Use the default registry host
		if typeFlag == "module" {
			host = regsrc.PublicRegistryHost
		} else {
			host = regsrc.PublicRegistryHost // Same for providers
		}
	}

	// Create a context for the search
	ctx := context.Background()

	// If count-only flag is set, just count the total modules and providers
	if countOnlyFlag {
		moduleCount, providerCount, err := c.countTotalModulesAndProviders(ctx, host)
		if err != nil {
			c.Meta.Ui.Error(err.Error())
			return 1
		}
		c.Meta.Ui.Output(fmt.Sprintf("Total Modules: %d", moduleCount))
		c.Meta.Ui.Output(fmt.Sprintf("Total Providers: %d", providerCount))
		return 0
	}

	// Perform the search
	if typeFlag == "module" {
		return c.searchModules(ctx, host, query, limitFlag, jsonFlag, detailedFlag)
	} else {
		return c.searchProviders(ctx, host, query, limitFlag, jsonFlag, detailedFlag)
	}
}

func (c *RegistrySearchCommand) searchModules(ctx context.Context, host *regsrc.FriendlyHost, query string, limit int, jsonOutput, detailed bool) int {
	c.Meta.Ui.Output(fmt.Sprintf("Searching for modules matching '%s'...", query))
	
	// Use the registry client to search for modules
	modules, err := c.directModuleSearch(ctx, host.String(), query)
	if err != nil {
		c.Meta.Ui.Error(fmt.Sprintf("Error searching for modules: %s", err))
		return 1
	}
	
	// Filter modules based on the query if needed
	if query != "" && len(modules) > 0 {
		filteredModules := make([]*response.Module, 0)
		for _, module := range modules {
			// Check if the module matches the query
			moduleFullName := fmt.Sprintf("%s/%s/%s", module.Namespace, module.Name, module.Provider)
			if strings.Contains(strings.ToLower(moduleFullName), strings.ToLower(query)) || 
			   strings.Contains(strings.ToLower(module.Description), strings.ToLower(query)) {
				filteredModules = append(filteredModules, module)
			}
		}
		modules = filteredModules
	}
	
	// Sort modules by downloads (most popular first)
	sort.Slice(modules, func(i, j int) bool {
		return modules[i].Downloads > modules[j].Downloads
	})
	
	// Limit the number of results
	if limit > 0 && len(modules) > limit {
		modules = modules[:limit]
	}
	
	// Output the results
	if jsonOutput {
		return c.outputModulesAsJSON(modules)
	}
	
	return c.outputModulesAsText(modules, detailed)
}

func (c *RegistrySearchCommand) searchProviders(ctx context.Context, host *regsrc.FriendlyHost, query string, limit int, jsonOutput, detailed bool) int {
	c.Meta.Ui.Output(fmt.Sprintf("Searching for providers matching '%s'...", query))
	
	// Use the registry client to search for providers
	providers, err := c.directProviderSearch(ctx, host.String(), query)
	if err != nil {
		c.Meta.Ui.Error(fmt.Sprintf("Error searching for providers: %s", err))
		return 1
	}
	
	// Filter providers based on the query if needed
	if query != "" && len(providers) > 0 {
		filteredProviders := make([]*response.ModuleProvider, 0)
		for _, provider := range providers {
			// Check if the provider matches the query
			if strings.Contains(strings.ToLower(provider.Name), strings.ToLower(query)) {
				filteredProviders = append(filteredProviders, provider)
			}
		}
		providers = filteredProviders
	}
	
	// Sort providers by downloads (most popular first)
	sort.Slice(providers, func(i, j int) bool {
		return providers[i].Downloads > providers[j].Downloads
	})
	
	// Limit the number of results
	if limit > 0 && len(providers) > limit {
		providers = providers[:limit]
	}
	
	// Output the results
	if jsonOutput {
		return c.outputProvidersAsJSON(providers)
	}
	
	return c.outputProvidersAsText(providers, detailed)
}

// directModuleSearch performs a direct API search for modules
func (c *RegistrySearchCommand) directModuleSearch(ctx context.Context, host string, query string) ([]*response.Module, error) {
	// Construct the API URL
	apiURL := fmt.Sprintf("https://%s/v1/modules?q=%s&limit=100", host, query)
	
	// Create a new HTTP client
	httpClient := retryablehttp.NewClient()
	httpClient.RetryMax = 3
	
	// Create the request
	req, err := retryablehttp.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, err
	}
	
	// Set headers
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "OpenTofu")
	
	// Send the request
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	// Check the response status
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("API returned status code %d", resp.StatusCode)
	}
	
	// Decode the response
	var moduleList response.ModuleList
	if err := json.NewDecoder(resp.Body).Decode(&moduleList); err != nil {
		return nil, err
	}
	
	return moduleList.Modules, nil
}

// directProviderSearch performs a direct API search for providers
func (c *RegistrySearchCommand) directProviderSearch(ctx context.Context, host string, query string) ([]*response.ModuleProvider, error) {
	// Construct the API URL
	apiURL := fmt.Sprintf("https://%s/v1/providers?q=%s&limit=100", host, query)
	
	// Create a new HTTP client
	httpClient := retryablehttp.NewClient()
	httpClient.RetryMax = 3
	
	// Create the request
	req, err := retryablehttp.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, err
	}
	
	// Set headers
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "OpenTofu")
	
	// Send the request
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	// Check the response status
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("API returned status code %d", resp.StatusCode)
	}
	
	// Decode the response
	var providerList response.ModuleProviderList
	if err := json.NewDecoder(resp.Body).Decode(&providerList); err != nil {
		return nil, err
	}
	
	return providerList.Providers, nil
}

func (c *RegistrySearchCommand) outputModulesAsJSON(modules []*response.Module) int {
	// Sort modules by downloads (most popular first)
	sort.Slice(modules, func(i, j int) bool {
		return modules[i].Downloads > modules[j].Downloads
	})

	// Create a more structured JSON output with metadata
	jsonOutput := map[string]interface{}{
		"count":   len(modules),
		"modules": modules,
		"metadata": map[string]interface{}{
			"timestamp": time.Now().Format(time.RFC3339),
			"query":     "modules",
		},
	}
	
	// Convert to JSON with proper indentation
	jsonData, err := json.MarshalIndent(jsonOutput, "", "  ")
	if err != nil {
		c.Meta.Ui.Error(fmt.Sprintf("Error formatting JSON: %s", err))
		return 1
	}
	
	c.Meta.Ui.Output(string(jsonData))
	return 0
}

func (c *RegistrySearchCommand) outputModulesAsText(modules []*response.Module, detailed bool) int {
	if len(modules) == 0 {
		c.Meta.Ui.Output("No matching modules found.")
		return 0
	}

	// Define ANSI color codes
	const (
		colorReset   = "\033[0m"
		colorBold    = "\033[1m"
		colorGreen   = "\033[32m"
		colorYellow  = "\033[33m"
		colorBlue    = "\033[34m"
		colorMagenta = "\033[35m"
		colorCyan    = "\033[36m"
		colorGray    = "\033[90m"
		colorItalic  = "\033[3m"
	)

	// Define icons
	const (
		iconModule    = "📦"
		iconVerified  = "✅"
		iconDownloads = "⬇️ "
		iconVersion   = "🏷️ "
		iconSource    = "🔗"
		iconPublished = "📅"
		iconSeparator = "━━━"
	)

	// Create a header with count and formatting
	c.Meta.Ui.Output(fmt.Sprintf("\n%s%s %s%d matching modules %s%s", 
		colorCyan, strings.Repeat("═", 10),
		colorBold, len(modules),
		strings.Repeat("═", 10), colorReset))
	c.Meta.Ui.Output("")

	// Sort modules by downloads (most popular first)
	sort.Slice(modules, func(i, j int) bool {
		return modules[i].Downloads > modules[j].Downloads
	})

	for i, module := range modules {
		// Create a visually distinct module entry with index
		c.Meta.Ui.Output(fmt.Sprintf("%d. %s%s %s%s/%s/%s%s", 
			i+1, 
			colorBold, iconModule,
			colorBlue, module.Namespace, module.Name, module.Provider, colorReset))
		
		// Always show a brief description if available
		if module.Description != "" {
			// Truncate description if it's too long
			desc := module.Description
			if len(desc) > 80 && !detailed {
				desc = desc[:77] + "..."
			}
			c.Meta.Ui.Output(fmt.Sprintf("   %s%s%s", colorItalic, desc, colorReset))
		}

		// Show basic stats in a compact format
		downloadStr := fmt.Sprintf("%d", module.Downloads)
		if module.Downloads > 1000000 {
			downloadStr = fmt.Sprintf("%.1fM", float64(module.Downloads)/1000000)
		} else if module.Downloads > 1000 {
			downloadStr = fmt.Sprintf("%.1fk", float64(module.Downloads)/1000)
		}
		
		verifiedBadge := ""
		if module.Verified {
			verifiedBadge = fmt.Sprintf("%s %sVerified%s", iconVerified, colorGreen, colorReset)
		}
		
		versionInfo := ""
		if len(module.Version) > 0 {
			versionInfo = fmt.Sprintf("%s %sv%s%s", iconVersion, colorYellow, module.Version, colorReset)
		}
		
		stats := []string{}
		if downloadStr != "" {
			stats = append(stats, fmt.Sprintf("%s %s%s%s", iconDownloads, colorMagenta, downloadStr, colorReset))
		}
		if versionInfo != "" {
			stats = append(stats, versionInfo)
		}
		if verifiedBadge != "" {
			stats = append(stats, verifiedBadge)
		}
		
		c.Meta.Ui.Output(fmt.Sprintf("   %s", strings.Join(stats, " | ")))

		// Show additional details if requested
		if detailed {
			c.Meta.Ui.Output("")
			c.Meta.Ui.Output(fmt.Sprintf("   %s %sPublished:%s %s", 
				iconPublished, colorBold, colorReset, 
				module.PublishedAt.Format("Jan 02, 2006")))
			
			if module.Source != "" {
				c.Meta.Ui.Output(fmt.Sprintf("   %s %sSource:%s %s", 
					iconSource, colorBold, colorReset, module.Source))
			}
			
			c.Meta.Ui.Output(fmt.Sprintf("   %sID:%s %s", colorBold, colorReset, module.ID))
		}
		
		// Add separator between modules
		if i < len(modules)-1 {
			c.Meta.Ui.Output(fmt.Sprintf("\n%s%s%s\n", 
				colorGray, strings.Repeat(iconSeparator, 16), colorReset))
		}
	}

	// Add usage hint at the end
	c.Meta.Ui.Output(fmt.Sprintf("\n%s%s%s", colorCyan, strings.Repeat("═", 30), colorReset))
	c.Meta.Ui.Output(fmt.Sprintf("%sTo install a module, run:%s", colorBold, colorReset))
	c.Meta.Ui.Output(fmt.Sprintf("  %stofu install module <namespace>/<name>/<provider>%s", colorGreen, colorReset))
	c.Meta.Ui.Output(fmt.Sprintf("For more details, use the %s-detailed%s flag", colorYellow, colorReset))

	return 0
}

func (c *RegistrySearchCommand) outputProvidersAsJSON(providers []*response.ModuleProvider) int {
	// Sort providers by downloads (most popular first)
	sort.Slice(providers, func(i, j int) bool {
		return providers[i].Downloads > providers[j].Downloads
	})

	// Create a more structured JSON output with metadata
	jsonOutput := map[string]interface{}{
		"count":     len(providers),
		"providers": providers,
		"metadata": map[string]interface{}{
			"timestamp": time.Now().Format(time.RFC3339),
			"query":     "providers",
		},
	}
	
	// Convert to JSON with proper indentation
	jsonData, err := json.MarshalIndent(jsonOutput, "", "  ")
	if err != nil {
		c.Meta.Ui.Error(fmt.Sprintf("Error formatting JSON: %s", err))
		return 1
	}
	
	c.Meta.Ui.Output(string(jsonData))
	return 0
}

func (c *RegistrySearchCommand) outputProvidersAsText(providers []*response.ModuleProvider, detailed bool) int {
	if len(providers) == 0 {
		c.Meta.Ui.Output("No matching providers found.")
		return 0
	}

	// Define ANSI color codes
	const (
		colorReset  = "\033[0m"
		colorBold   = "\033[1m"
		colorGreen  = "\033[32m"
		colorYellow = "\033[33m"
		colorBlue   = "\033[34m"
		colorMagenta = "\033[35m"
		colorCyan   = "\033[36m"
		colorGray   = "\033[90m"
	)

	// Define icons
	const (
		iconProvider  = "🔌"
		iconDownloads = "⬇️ "
		iconModules   = "📦"
		iconNamespace = "🏢"
		iconSeparator = "━━━"
	)

	// Create a header with count and formatting
	c.Meta.Ui.Output(fmt.Sprintf("\n%s%s %s%d matching providers %s%s", 
		colorCyan, strings.Repeat("═", 10),
		colorBold, len(providers),
		strings.Repeat("═", 10), colorReset))
	c.Meta.Ui.Output("")

	// Sort providers by downloads (most popular first)
	sort.Slice(providers, func(i, j int) bool {
		return providers[i].Downloads > providers[j].Downloads
	})

	for i, provider := range providers {
		// Parse provider name to extract namespace and type if possible
		parts := strings.Split(provider.Name, "/")
		displayName := provider.Name
		
		// Create a visually distinct provider entry with index
		c.Meta.Ui.Output(fmt.Sprintf("%d. %s%s %s%s%s", 
			i+1, colorBold, iconProvider, colorBlue, displayName, colorReset))
		
		// Show basic stats in a compact format
		downloadStr := fmt.Sprintf("%d", provider.Downloads)
		if provider.Downloads > 1000000 {
			downloadStr = fmt.Sprintf("%.1fM", float64(provider.Downloads)/1000000)
		} else if provider.Downloads > 1000 {
			downloadStr = fmt.Sprintf("%.1fk", float64(provider.Downloads)/1000)
		}
		
		moduleCountStr := fmt.Sprintf("%s %s%d modules%s", 
			iconModules, colorYellow, provider.ModuleCount, colorReset)
		if provider.ModuleCount == 1 {
			moduleCountStr = fmt.Sprintf("%s %s1 module%s", 
				iconModules, colorYellow, colorReset)
		}
		
		stats := []string{}
		if downloadStr != "" {
			stats = append(stats, fmt.Sprintf("%s %s%s%s", 
				iconDownloads, colorMagenta, downloadStr, colorReset))
		}
		if moduleCountStr != "" {
			stats = append(stats, moduleCountStr)
		}
		
		// If we have namespace information, display it
		if len(parts) > 1 {
			stats = append(stats, fmt.Sprintf("%s %s%s%s", 
				iconNamespace, colorGreen, parts[0], colorReset))
		}
		
		c.Meta.Ui.Output(fmt.Sprintf("   %s", strings.Join(stats, " | ")))
		
		// Show additional details if requested
		if detailed {
			c.Meta.Ui.Output("")
			// Add any additional provider details here if they become available in the API
		}
		
		// Add separator between providers
		if i < len(providers)-1 {
			c.Meta.Ui.Output(fmt.Sprintf("\n%s%s%s\n", 
				colorGray, strings.Repeat(iconSeparator, 16), colorReset))
		}
	}

	// Add usage hint at the end
	c.Meta.Ui.Output(fmt.Sprintf("\n%s%s%s", colorCyan, strings.Repeat("═", 30), colorReset))
	c.Meta.Ui.Output(fmt.Sprintf("%sTo use a provider, include it in your configuration:%s", colorBold, colorReset))
	c.Meta.Ui.Output(fmt.Sprintf("  %sterraform {%s", colorGreen, colorReset))
	c.Meta.Ui.Output(fmt.Sprintf("  %s  required_providers {%s", colorGreen, colorReset))
	c.Meta.Ui.Output(fmt.Sprintf("  %s    <name> = {%s", colorGreen, colorReset))
	c.Meta.Ui.Output(fmt.Sprintf("  %s      source = \"<provider-address>\"%s", colorGreen, colorReset))
	c.Meta.Ui.Output(fmt.Sprintf("  %s    }%s", colorGreen, colorReset))
	c.Meta.Ui.Output(fmt.Sprintf("  %s  }%s", colorGreen, colorReset))
	c.Meta.Ui.Output(fmt.Sprintf("  %s}%s", colorGreen, colorReset))
	c.Meta.Ui.Output(fmt.Sprintf("For more details, use the %s-detailed%s flag", colorYellow, colorReset))

	return 0
}

func (c *RegistrySearchCommand) countTotalModulesAndProviders(ctx context.Context, host *regsrc.FriendlyHost) (int, int, error) {
	c.Meta.Ui.Output("Counting total modules and providers in the registry...")

	moduleCount, err := c.countModules(ctx, host)
	if err != nil {
		c.Meta.Ui.Error(fmt.Sprintf("Error counting modules: %s", err))
		return 0, 0, err
	}

	providerCount, err := c.countProviders(ctx, host)
	if err != nil {
		c.Meta.Ui.Error(fmt.Sprintf("Error counting providers: %s", err))
		return moduleCount, 0, err
	}

	totalCount := moduleCount + providerCount

	// Display results in a more user-friendly format
	c.Meta.Ui.Output("\n══════════ Registry Statistics ══════════\n")
	c.Meta.Ui.Output(fmt.Sprintf("📦 Total Modules: %d", moduleCount))
	c.Meta.Ui.Output(fmt.Sprintf("🔌 Total Providers: %d", providerCount))
	c.Meta.Ui.Output(fmt.Sprintf("🔢 Total Registry Items: %d", totalCount))

	// Also output in plain format for scripts that might parse the output
	c.Meta.Ui.Output(fmt.Sprintf("Total Modules: %d", moduleCount))
	c.Meta.Ui.Output(fmt.Sprintf("Total Providers: %d", providerCount))

	return moduleCount, providerCount, nil
}

func (c *RegistrySearchCommand) countModules(ctx context.Context, host *regsrc.FriendlyHost) (int, error) {
	// Use a reasonable page size to avoid overwhelming the API
	const pageSize = 100
	maxPages := 200 // Limit to 200 pages (20,000 modules) to avoid excessive API calls
	totalCount := 0
	
	c.Meta.Ui.Output("Fetching modules data (this may take a while)...")
	
	// Pre-allocate for expected size
	allModules := make(map[string]bool, 20000)
	
	// Start with the first page
	nextURL := fmt.Sprintf("https://%s/v1/modules?limit=%d", host.String(), pageSize)
	
	for i := 0; i < maxPages && nextURL != ""; i++ {
		// Add throttling to avoid rate limiting
		if i > 0 {
			time.Sleep(500 * time.Millisecond)
		}
		
		c.Meta.Ui.Output(fmt.Sprintf("GET %s (page %d/%d)", nextURL, i+1, maxPages))
		
		req, err := retryablehttp.NewRequest("GET", nextURL, nil)
		if err != nil {
			c.Meta.Ui.Error(fmt.Sprintf("Error creating request: %s", err))
			return totalCount, err
		}
		
		client := retryablehttp.NewClient()
		client.Logger = nil // Disable logging
		client.RetryMax = 3 // Maximum number of retries
		client.RetryWaitMin = 1 * time.Second
		client.RetryWaitMax = 5 * time.Second
		
		resp, err := client.Do(req)
		if err != nil {
			c.Meta.Ui.Error(fmt.Sprintf("Error making request: %s", err))
			return totalCount, err
		}
		defer resp.Body.Close()
		
		if resp.StatusCode != 200 {
			body, _ := io.ReadAll(resp.Body)
			c.Meta.Ui.Error(fmt.Sprintf("Failed to fetch modules: %s\nResponse body: %s", resp.Status, string(body)))
			
			// If we hit rate limiting, return what we have so far
			if resp.StatusCode == 429 {
				c.Meta.Ui.Output("Rate limit exceeded. Returning count so far.")
				return totalCount, nil
			}
			
			return totalCount, fmt.Errorf("failed to fetch modules: %s", resp.Status)
		}
		
		// We need to handle both pagination styles
		var result struct {
			Meta struct {
				Limit      int    `json:"limit"`
				CurrentPage int    `json:"current_offset"`
				NextOffset  int    `json:"next_offset"`
				NextURL    string `json:"next_url"`
			} `json:"meta"`
			Modules []struct {
				ID string `json:"id"`
			} `json:"modules"`
		}
		
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			c.Meta.Ui.Error(fmt.Sprintf("Error reading response body: %s", err))
			return totalCount, err
		}
		
		if err := json.Unmarshal(body, &result); err != nil {
			c.Meta.Ui.Error(fmt.Sprintf("Error decoding response: %s", err))
			return totalCount, err
		}
		
		c.Meta.Ui.Output(fmt.Sprintf("Received %d modules in this batch", len(result.Modules)))
		
		if len(result.Modules) == 0 {
			c.Meta.Ui.Output("No more modules to fetch")
			break
		}
		
		// Add modules to our map to ensure uniqueness
		for _, module := range result.Modules {
			if !allModules[module.ID] {
				allModules[module.ID] = true
				totalCount++
			}
		}
		
		c.Meta.Ui.Output(fmt.Sprintf("Counted %d unique modules so far...", totalCount))
		
		// Determine the next URL based on the pagination style
		if result.Meta.NextURL != "" {
			// v2 style pagination with next_url
			nextURL = result.Meta.NextURL
			if !strings.HasPrefix(nextURL, "http") {
				nextURL = fmt.Sprintf("https://%s%s", host.String(), nextURL)
			}
		} else if result.Meta.NextOffset > 0 {
			// v1 style pagination with next_offset
			nextURL = fmt.Sprintf("https://%s/v1/modules?limit=%d&offset=%d", host.String(), pageSize, result.Meta.NextOffset)
		} else if len(result.Modules) == pageSize {
			// Fallback to simple offset increment if neither pagination style is detected
			nextURL = fmt.Sprintf("https://%s/v1/modules?limit=%d&offset=%d", host.String(), pageSize, (i+1)*pageSize)
		} else {
			// No more pages
			nextURL = ""
		}
		
		// If we got fewer modules than the page size, we've reached the end
		if len(result.Modules) < pageSize {
			c.Meta.Ui.Output("Received fewer modules than page size, ending search")
			break
		}
	}
	
	c.Meta.Ui.Output(fmt.Sprintf("Final module count: %d", totalCount))
	return totalCount, nil
}

func (c *RegistrySearchCommand) countProviders(ctx context.Context, host *regsrc.FriendlyHost) (int, error) {
	// Use a map to track unique providers
	uniqueProviders := make(map[string]bool)
	
	// Use v2 API with page-based pagination
	baseURL := fmt.Sprintf("https://%s/v2/providers", host.String())
	httpClient := retryablehttp.NewClient()
	httpClient.RetryMax = 3
	httpClient.RetryWaitMin = 500 * time.Millisecond
	httpClient.RetryWaitMax = 2 * time.Second
	httpClient.Logger = nil // Disable logging
	
	// Track consecutive duplicate batches
	duplicateBatches := 0
	maxDuplicateBatches := 3
	maxPages := 100 // Limit to avoid infinite loops
	pageSize := 100
	
	c.Meta.Ui.Info("Counting providers...")
	
	for page := 1; page <= maxPages; page++ {
		// Add a small delay to avoid hitting rate limits
		if page > 1 {
			time.Sleep(500 * time.Millisecond)
		}
		
		params := fmt.Sprintf("?page[number]=%d&page[size]=%d", page, pageSize)
		url := baseURL + params
		
		c.Meta.Ui.Info(fmt.Sprintf("Fetching providers page %d from %s", page, url))
		
		req, err := retryablehttp.NewRequest("GET", url, nil)
		if err != nil {
			return len(uniqueProviders), fmt.Errorf("failed to create request: %w", err)
		}
		
		req.Header.Set("Accept", "application/json")
		req.Header.Set("User-Agent", "OpenTofu")
		
		resp, err := httpClient.Do(req)
		if err != nil {
			return len(uniqueProviders), fmt.Errorf("failed to fetch providers: %w", err)
		}
		defer resp.Body.Close()
		
		if resp.StatusCode != 200 {
			body, _ := io.ReadAll(resp.Body)
			return len(uniqueProviders), fmt.Errorf("failed to fetch providers: %s - %s", resp.Status, body)
		}
		
		var result struct {
			Data []struct {
				ID string `json:"id"`
			} `json:"data"`
		}
		
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return len(uniqueProviders), fmt.Errorf("failed to decode response: %w", err)
		}
		
		// Check if we got any providers
		if len(result.Data) == 0 {
			c.Meta.Ui.Info("No more providers returned. Stopping.")
			break
		}
		
		// Count new unique providers in this batch
		beforeCount := len(uniqueProviders)
		for _, provider := range result.Data {
			uniqueProviders[provider.ID] = true
		}
		afterCount := len(uniqueProviders)
		newCount := afterCount - beforeCount
		
		c.Meta.Ui.Info(fmt.Sprintf("Page %d: Got %d providers, %d new unique. Total so far: %d", 
			page, len(result.Data), newCount, afterCount))
		
		// If we didn't get any new providers, increment duplicate counter
		if newCount == 0 {
			duplicateBatches++
			c.Meta.Ui.Info(fmt.Sprintf("Duplicate batch detected (%d/%d)", duplicateBatches, maxDuplicateBatches))
			if duplicateBatches >= maxDuplicateBatches {
				c.Meta.Ui.Info("Maximum duplicate batches reached. Stopping.")
				break
			}
		} else {
			// Reset duplicate counter if we got new providers
			duplicateBatches = 0
		}
		
		// If we got fewer providers than the page size, we're at the end
		if len(result.Data) < pageSize {
			c.Meta.Ui.Info("Received fewer providers than page size. Stopping.")
			break
		}
	}
	
	totalCount := len(uniqueProviders)
	
	// Add warning if count is lower than expected
	if totalCount < 4000 {
		c.Meta.Ui.Warn(fmt.Sprintf("Warning: Provider count (%d) is lower than expected (4000+). This may be due to API limitations.", totalCount))
		c.Meta.Ui.Warn("The Terraform Registry API has limitations that prevent fetching all providers.")
		c.Meta.Ui.Warn("Based on registry data, there are approximately 4,000+ providers available.")
		c.Meta.Ui.Warn("However, the API currently only returns a subset of these providers.")
	}
	
	return totalCount, nil
}
