package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"sort"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/fatih/color"
	"github.com/hashicorp/go-hclog"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/opentofu/opentofu/internal/registry"
	"github.com/opentofu/opentofu/internal/registry/regsrc"
	"github.com/opentofu/opentofu/internal/registry/response"
)

// SearchResult represents a search result for modules or providers
type SearchResult struct {
	Type        string    // "module" or "provider"
	Namespace   string    // Owner/organization
	Name        string    // Name of the module/provider
	Provider    string    // Provider (for modules only)
	Description string    // Description
	Version     string    // Latest version
	Downloads   int       // Download count
	Stars       int       // Star count (if available)
	Published   time.Time // Published date
	Verified    bool      // Whether the module/provider is verified
}

func main() {
	// Create a logger
	logger := hclog.New(&hclog.LoggerOptions{
		Name:  "opentofu-registry",
		Level: hclog.Info,
		Color: hclog.AutoColor,
	})

	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		logger.Error("Error loading .env file", "error", err)
		return
	}

	// Test database connection
	logger.Info("Testing database connection...")
	dbClient, err := testDatabaseConnection()
	if err != nil {
		logger.Error("Database connection failed", "error", err)
		return
	}
	logger.Info("Database connection successful", "database", os.Getenv("TOFU_REGISTRY_DB_NAME"))

	// Initialize registry client
	logger.Info("Initializing registry client...")
	client, err := testRegistryClient()
	if err != nil {
		logger.Error("Registry client initialization failed", "error", err)
		return
	}
	logger.Info("Registry client initialized successfully")

	// Perform search
	query := "aws"
	if len(os.Args) > 1 {
		query = os.Args[1]
	}

	logger.Info("Searching for modules and providers", "query", query)
	results, err := performSearch(client, query)
	if err != nil {
		logger.Error("Search failed", "error", err)
		return
	}

	// Display results
	displaySearchResults(results, query)

	logger.Info("Search completed successfully!")
}

// testDatabaseConnection tests the connection to the PostgreSQL database
func testDatabaseConnection() (*registry.DBClient, error) {
	// Get database connection parameters from environment variables
	dbType := os.Getenv("TOFU_REGISTRY_DB_TYPE")
	dbHost := os.Getenv("TOFU_REGISTRY_DB_HOST")
	dbPort := os.Getenv("TOFU_REGISTRY_DB_PORT")
	dbName := os.Getenv("TOFU_REGISTRY_DB_NAME")
	dbUser := os.Getenv("TOFU_REGISTRY_DB_USER")
	dbPassword := os.Getenv("TOFU_REGISTRY_DB_PASSWORD")
	dbSSLMode := os.Getenv("TOFU_REGISTRY_DB_SSLMODE")

	// Connect to PostgreSQL database
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbHost, dbPort, dbUser, dbPassword, dbName, dbSSLMode)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error opening database connection: %w", err)
	}

	// Test connection
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	// Create a new DBClient
	dbClient := &registry.DBClient{
		DB:     db,
		DBType: dbType,
	}

	return dbClient, nil
}

// testRegistryClient initializes and tests the registry client
func testRegistryClient() (*registry.Client, error) {
	// Create an HTTP client
	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Create a new registry client
	client := registry.NewClient(nil, httpClient)
	if client == nil {
		return nil, fmt.Errorf("failed to create registry client")
	}

	return client, nil
}

// performSearch searches for modules and providers
func performSearch(client *registry.Client, query string) ([]SearchResult, error) {
	ctx := context.Background()
	results := []SearchResult{}

	// Search for modules
	fmt.Println("Searching for modules...")
	moduleResults, err := searchModules(ctx, client, query)
	if err != nil {
		return nil, fmt.Errorf("error searching for modules: %w", err)
	}
	results = append(results, moduleResults...)

	// Search for providers
	fmt.Println("Searching for providers...")
	providerResults, err := searchProviders(ctx, client, query)
	if err != nil {
		return nil, fmt.Errorf("error searching for providers: %w", err)
	}
	results = append(results, providerResults...)

	// Sort results by downloads (descending)
	sort.Slice(results, func(i, j int) bool {
		return results[i].Downloads > results[j].Downloads
	})

	return results, nil
}

// searchModules searches for modules
func searchModules(ctx context.Context, client *registry.Client, query string) ([]SearchResult, error) {
	// Pre-allocate based on known registry size (approximately 18,000 modules)
	results := make([]SearchResult, 0, 100)

	// Simulate module search (replace with actual implementation)
	// In a real implementation, you would use client.SearchModules or similar
	
	// For demonstration purposes, we'll create some sample results
	sampleModules := []struct {
		namespace   string
		name        string
		provider    string
		description string
		version     string
		downloads   int
		stars       int
		published   time.Time
		verified    bool
	}{
		{
			namespace:   "hashicorp",
			name:        "consul",
			provider:    "aws",
			description: "Terraform module for deploying Consul on AWS",
			version:     "0.8.0",
			downloads:   125000,
			stars:       120,
			published:   time.Now().AddDate(0, -2, 0),
			verified:    true,
		},
		{
			namespace:   "terraform-aws-modules",
			name:        "vpc",
			provider:    "aws",
			description: "Terraform module which creates VPC resources on AWS",
			version:     "3.14.0",
			downloads:   980000,
			stars:       450,
			published:   time.Now().AddDate(0, -1, -15),
			verified:    true,
		},
		{
			namespace:   "terraform-aws-modules",
			name:        "security-group",
			provider:    "aws",
			description: "Terraform module which creates security group resources on AWS",
			version:     "4.9.0",
			downloads:   750000,
			stars:       320,
			published:   time.Now().AddDate(0, -3, -5),
			verified:    true,
		},
		{
			namespace:   "cloudposse",
			name:        "vpc",
			provider:    "aws",
			description: "Terraform module to provision a VPC with Internet Gateway",
			version:     "0.28.1",
			downloads:   320000,
			stars:       180,
			published:   time.Now().AddDate(0, -5, -10),
			verified:    false,
		},
		{
			namespace:   "opentofu",
			name:        "registry",
			provider:    "aws",
			description: "OpenTofu module for deploying a registry on AWS",
			version:     "0.1.0",
			downloads:   5000,
			stars:       25,
			published:   time.Now().AddDate(0, 0, -5),
			verified:    false,
		},
	}

	// Filter modules based on the query
	for _, module := range sampleModules {
		if strings.Contains(strings.ToLower(module.namespace), strings.ToLower(query)) ||
			strings.Contains(strings.ToLower(module.name), strings.ToLower(query)) ||
			strings.Contains(strings.ToLower(module.provider), strings.ToLower(query)) ||
			strings.Contains(strings.ToLower(module.description), strings.ToLower(query)) {
			results = append(results, SearchResult{
				Type:        "module",
				Namespace:   module.namespace,
				Name:        module.name,
				Provider:    module.provider,
				Description: module.description,
				Version:     module.version,
				Downloads:   module.downloads,
				Stars:       module.stars,
				Published:   module.published,
				Verified:    module.verified,
			})
		}
	}

	return results, nil
}

// searchProviders searches for providers
func searchProviders(ctx context.Context, client *registry.Client, query string) ([]SearchResult, error) {
	// Pre-allocate based on known registry size (approximately 4,000 providers)
	results := make([]SearchResult, 0, 50)

	// Simulate provider search (replace with actual implementation)
	// In a real implementation, you would use client.SearchProviders or similar
	
	// For demonstration purposes, we'll create some sample results
	sampleProviders := []struct {
		namespace   string
		name        string
		description string
		version     string
		downloads   int
		stars       int
		published   time.Time
		verified    bool
	}{
		{
			namespace:   "hashicorp",
			name:        "aws",
			description: "Terraform AWS provider",
			version:     "4.67.0",
			downloads:   5000000,
			stars:       750,
			published:   time.Now().AddDate(0, -1, 0),
			verified:    true,
		},
		{
			namespace:   "hashicorp",
			name:        "azurerm",
			description: "Terraform Azure provider",
			version:     "3.65.0",
			downloads:   3500000,
			stars:       550,
			published:   time.Now().AddDate(0, -1, -10),
			verified:    true,
		},
		{
			namespace:   "hashicorp",
			name:        "google",
			description: "Terraform Google Cloud provider",
			version:     "4.80.0",
			downloads:   2800000,
			stars:       480,
			published:   time.Now().AddDate(0, -2, -5),
			verified:    true,
		},
		{
			namespace:   "digitalocean",
			name:        "digitalocean",
			description: "Terraform DigitalOcean provider",
			version:     "2.30.0",
			downloads:   950000,
			stars:       220,
			published:   time.Now().AddDate(0, -3, -15),
			verified:    true,
		},
		{
			namespace:   "opentofu",
			name:        "registry",
			description: "OpenTofu Registry provider",
			version:     "0.1.0",
			downloads:   8000,
			stars:       35,
			published:   time.Now().AddDate(0, 0, -3),
			verified:    false,
		},
	}

	// Filter providers based on the query
	for _, provider := range sampleProviders {
		if strings.Contains(strings.ToLower(provider.namespace), strings.ToLower(query)) ||
			strings.Contains(strings.ToLower(provider.name), strings.ToLower(query)) ||
			strings.Contains(strings.ToLower(provider.description), strings.ToLower(query)) {
			results = append(results, SearchResult{
				Type:        "provider",
				Namespace:   provider.namespace,
				Name:        provider.name,
				Provider:    "",
				Description: provider.description,
				Version:     provider.version,
				Downloads:   provider.downloads,
				Stars:       provider.stars,
				Published:   provider.published,
				Verified:    provider.verified,
			})
		}
	}

	return results, nil
}

// displaySearchResults displays the search results in a formatted table
func displaySearchResults(results []SearchResult, query string) {
	if len(results) == 0 {
		fmt.Printf("No results found for query: %s\n", query)
		return
	}

	fmt.Printf("Search results for: %s\n\n", query)

	// Create a tabwriter for aligned output
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	// Define colors
	bold := color.New(color.Bold).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	blue := color.New(color.FgBlue).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	cyan := color.New(color.FgCyan).SprintFunc()

	// Print header
	fmt.Fprintln(w, bold("TYPE\tNAMESPACE/NAME\tVERSION\tDOWNLOADS\tPUBLISHED\tVERIFIED"))
	fmt.Fprintln(w, "----\t--------------\t-------\t---------\t---------\t--------")

	// Print results
	for _, result := range results {
		// Format the name based on type
		var name string
		if result.Type == "module" {
			name = fmt.Sprintf("%s/%s/%s", result.Namespace, result.Name, result.Provider)
		} else {
			name = fmt.Sprintf("%s/%s", result.Namespace, result.Name)
		}

		// Format downloads
		var downloads string
		if result.Downloads >= 1000000 {
			downloads = fmt.Sprintf("%.1fM", float64(result.Downloads)/1000000)
		} else if result.Downloads >= 1000 {
			downloads = fmt.Sprintf("%.1fK", float64(result.Downloads)/1000)
		} else {
			downloads = fmt.Sprintf("%d", result.Downloads)
		}

		// Format published date
		published := result.Published.Format("2006-01-02")

		// Format verified status
		var verified string
		if result.Verified {
			verified = green("✓")
		} else {
			verified = ""
		}

		// Format type
		var typeStr string
		if result.Type == "module" {
			typeStr = blue("module")
		} else {
			typeStr = yellow("provider")
		}

		// Print the row
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n",
			typeStr,
			cyan(name),
			result.Version,
			downloads,
			published,
			verified)
	}

	w.Flush()

	// Print summary
	moduleCount := 0
	providerCount := 0
	for _, result := range results {
		if result.Type == "module" {
			moduleCount++
		} else {
			providerCount++
		}
	}

	fmt.Printf("\nFound %d results (%d modules, %d providers)\n", 
		len(results), moduleCount, providerCount)
	
	// Print usage instructions
	fmt.Println("\nTo use a module:")
	fmt.Println("  module \"example\" {")
	fmt.Println("    source  = \"namespace/name/provider\"")
	fmt.Println("    version = \"x.y.z\"")
	fmt.Println("  }")
	
	fmt.Println("\nTo use a provider:")
	fmt.Println("  terraform {")
	fmt.Println("    required_providers {")
	fmt.Println("      name = {")
	fmt.Println("        source  = \"namespace/name\"")
	fmt.Println("        version = \"x.y.z\"")
	fmt.Println("      }")
	fmt.Println("    }")
	fmt.Println("  }")
}
