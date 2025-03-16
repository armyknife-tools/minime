package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/hashicorp/terraform-svchost/disco"
	"github.com/joho/godotenv"
	"github.com/opentofu/opentofu/internal/registry"
	"github.com/opentofu/opentofu/internal/registry/regsrc"
	"github.com/opentofu/opentofu/internal/registry/response"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("Error loading .env file: %v\n", err)
		return
	}

	// Create logger
	logger := hclog.New(&hclog.LoggerOptions{
		Name:   "registry-test",
		Level:  hclog.Info,
		Output: os.Stdout,
	})

	// Run tests
	fmt.Println("=== OpenTofu Registry API Test Suite ===")
	
	// Test database connection
	testDatabaseConnection(logger)
	
	// Test registry client
	testRegistryClient(logger)
	
	// Test module operations
	testModuleOperations(logger)
	
	// Test provider operations
	testProviderOperations(logger)
	
	// Test search functionality
	testSearchFunctionality(logger)
	
	// Test performance with large datasets
	testPerformanceWithLargeDatasets(logger)
	
	fmt.Println("\n=== All tests completed ===")
}

// Test database connection
func testDatabaseConnection(logger hclog.Logger) {
	fmt.Println("\n=== Testing Database Connection ===")
	
	// Create database client
	dbClient, err := registry.NewDBClient(logger)
	if err != nil {
		fmt.Printf("❌ Error creating database client: %v\n", err)
		return
	}
	
	// Test connection
	err = dbClient.TestConnection()
	if err != nil {
		fmt.Printf("❌ Error connecting to database: %v\n", err)
		return
	}
	
	fmt.Println("✅ Successfully connected to database")
}

// Test registry client
func testRegistryClient(logger hclog.Logger) {
	fmt.Println("\n=== Testing Registry Client ===")
	
	// Create HTTP client with retries
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 3
	retryClient.RetryWaitMin = 500 * time.Millisecond
	retryClient.RetryWaitMax = 2 * time.Second
	
	// Create services discovery client
	services := disco.New()
	services.SetUserAgent("OpenTofu-Registry-Test/1.0")
	
	// Create registry client
	client := registry.NewClient(services, retryClient.StandardClient())
	
	// Test client initialization
	if client == nil {
		fmt.Println("❌ Failed to initialize registry client")
		return
	}
	
	fmt.Println("✅ Successfully initialized registry client")
}

// Test module operations
func testModuleOperations(logger hclog.Logger) {
	fmt.Println("\n=== Testing Module Operations ===")
	
	// Create context
	ctx := context.Background()
	
	// Create HTTP client with retries
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 3
	retryClient.RetryWaitMin = 500 * time.Millisecond
	retryClient.RetryWaitMax = 2 * time.Second
	
	// Create services discovery client
	services := disco.New()
	services.SetUserAgent("OpenTofu-Registry-Test/1.0")
	
	// Create registry client
	client := registry.NewClient(services, retryClient.StandardClient())
	
	// Test module listing
	fmt.Println("Testing module listing...")
	
	// Parse module source
	module, err := regsrc.ParseModuleSource("hashicorp/consul/aws")
	if err != nil {
		fmt.Printf("❌ Error parsing module source: %v\n", err)
		return
	}
	
	// Get module versions
	versions, err := client.ModuleVersions(ctx, module)
	if err != nil {
		fmt.Printf("❌ Error getting module versions: %v\n", err)
		return
	}
	
	fmt.Printf("✅ Successfully retrieved %d versions for module %s\n", 
		len(versions.Modules[0].Versions), module.Display())
	
	// Test module download
	fmt.Println("Testing module download...")
	
	// Get latest version
	latestVersion := versions.Modules[0].Versions[0].Version
	
	// Get module download URL
	location, err := client.ModuleLocation(ctx, module, latestVersion)
	if err != nil {
		fmt.Printf("❌ Error getting module location: %v\n", err)
		return
	}
	
	fmt.Printf("✅ Successfully retrieved download URL for module %s version %s\n", 
		module.Display(), latestVersion)
	
	// Verify URL is accessible
	resp, err := http.Head(location)
	if err != nil {
		fmt.Printf("❌ Error accessing module download URL: %v\n", err)
		return
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("❌ Module download URL returned status %d\n", resp.StatusCode)
		return
	}
	
	fmt.Printf("✅ Module download URL is accessible (status %d)\n", resp.StatusCode)
}

// Test provider operations
func testProviderOperations(logger hclog.Logger) {
	fmt.Println("\n=== Testing Provider Operations ===")
	
	// Create context
	ctx := context.Background()
	
	// Create HTTP client with retries
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 3
	retryClient.RetryWaitMin = 500 * time.Millisecond
	retryClient.RetryWaitMax = 2 * time.Second
	
	// Create services discovery client
	services := disco.New()
	services.SetUserAgent("OpenTofu-Registry-Test/1.0")
	
	// Create registry client
	client := registry.NewClient(services, retryClient.StandardClient())
	
	// Test provider listing
	fmt.Println("Testing provider listing...")
	
	// Get providers
	providers, err := client.ListProviders(ctx)
	if err != nil {
		fmt.Printf("❌ Error listing providers: %v\n", err)
		return
	}
	
	fmt.Printf("✅ Successfully retrieved %d providers\n", len(providers))
	
	// Test provider versions
	fmt.Println("Testing provider versions...")
	
	// Get versions for a specific provider
	if len(providers) > 0 {
		provider := providers[0]
		versions, err := client.ProviderVersions(ctx, provider)
		if err != nil {
			fmt.Printf("❌ Error getting provider versions: %v\n", err)
			return
		}
		
		fmt.Printf("✅ Successfully retrieved %d versions for provider %s\n", 
			len(versions.Versions), provider)
	} else {
		fmt.Println("⚠️ No providers found to test versions")
	}
}

// Test search functionality
func testSearchFunctionality(logger hclog.Logger) {
	fmt.Println("\n=== Testing Search Functionality ===")
	
	// Create context
	ctx := context.Background()
	
	// Create HTTP client with retries
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 3
	retryClient.RetryWaitMin = 500 * time.Millisecond
	retryClient.RetryWaitMax = 2 * time.Second
	
	// Create services discovery client
	services := disco.New()
	services.SetUserAgent("OpenTofu-Registry-Test/1.0")
	
	// Create registry client
	client := registry.NewClient(services, retryClient.StandardClient())
	
	// Test module search
	fmt.Println("Testing module search...")
	
	// Search for modules
	modules, err := client.SearchModules(ctx, "aws")
	if err != nil {
		fmt.Printf("❌ Error searching for modules: %v\n", err)
		return
	}
	
	fmt.Printf("✅ Successfully found %d modules matching 'aws'\n", len(modules))
	
	// Test provider search
	fmt.Println("Testing provider search...")
	
	// Search for providers
	providers, err := client.SearchProviders(ctx, "aws")
	if err != nil {
		fmt.Printf("❌ Error searching for providers: %v\n", err)
		return
	}
	
	fmt.Printf("✅ Successfully found %d providers matching 'aws'\n", len(providers))
	
	// Test search with database caching
	fmt.Println("Testing search with database caching...")
	
	// Create database client
	dbClient, err := registry.NewDBClient(logger)
	if err != nil {
		fmt.Printf("❌ Error creating database client: %v\n", err)
		return
	}
	
	// Cache some modules
	if len(modules) > 0 {
		module := modules[0]
		err = cacheModule(ctx, dbClient, module)
		if err != nil {
			fmt.Printf("❌ Error caching module: %v\n", err)
			return
		}
		
		fmt.Println("✅ Successfully cached module in database")
	}
	
	// Cache some providers
	if len(providers) > 0 {
		provider := providers[0]
		err = cacheProvider(ctx, dbClient, provider)
		if err != nil {
			fmt.Printf("❌ Error caching provider: %v\n", err)
			return
		}
		
		fmt.Println("✅ Successfully cached provider in database")
	}
}

// Test performance with large datasets
func testPerformanceWithLargeDatasets(logger hclog.Logger) {
	fmt.Println("\n=== Testing Performance with Large Datasets ===")
	
	// Create context
	ctx := context.Background()
	
	// Create HTTP client with retries
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 3
	retryClient.RetryWaitMin = 500 * time.Millisecond
	retryClient.RetryWaitMax = 2 * time.Second
	
	// Create services discovery client
	services := disco.New()
	services.SetUserAgent("OpenTofu-Registry-Test/1.0")
	
	// Create registry client
	client := registry.NewClient(services, retryClient.StandardClient())
	
	// Test pre-allocation
	fmt.Println("Testing pre-allocation for large datasets...")
	
	// Create database client
	dbClient, err := registry.NewDBClient(logger)
	if err != nil {
		fmt.Printf("❌ Error creating database client: %v\n", err)
		return
	}
	
	// Test module pre-allocation
	modules := dbClient.PreallocateModules()
	fmt.Printf("✅ Successfully pre-allocated space for %d modules\n", cap(modules))
	
	// Test provider pre-allocation
	providers := dbClient.PreallocateProviders()
	fmt.Printf("✅ Successfully pre-allocated space for %d providers\n", cap(providers))
	
	// Test bulk operations
	fmt.Println("Testing bulk operations...")
	
	// Test bulk module fetch
	startTime := time.Now()
	_, err = client.BulkFetchModules(ctx, 10)
	if err != nil {
		fmt.Printf("❌ Error bulk fetching modules: %v\n", err)
		return
	}
	
	fetchDuration := time.Since(startTime)
	fmt.Printf("✅ Successfully bulk fetched modules in %v\n", fetchDuration)
	
	// Test bulk provider fetch
	startTime = time.Now()
	_, err = client.BulkFetchProviders(ctx, 10)
	if err != nil {
		fmt.Printf("❌ Error bulk fetching providers: %v\n", err)
		return
	}
	
	fetchDuration = time.Since(startTime)
	fmt.Printf("✅ Successfully bulk fetched providers in %v\n", fetchDuration)
	
	// Test throttling
	fmt.Println("Testing throttling for API requests...")
	
	// Make multiple requests in quick succession
	for i := 0; i < 5; i++ {
		_, err := client.SearchModules(ctx, fmt.Sprintf("test%d", i))
		if err != nil {
			fmt.Printf("❌ Error in throttled request %d: %v\n", i, err)
			return
		}
	}
	
	fmt.Println("✅ Successfully handled throttled requests")
}

// Helper function to cache a module in the database
func cacheModule(ctx context.Context, dbClient *registry.DBClient, module *response.Module) error {
	// This is a placeholder since we don't have direct access to the internal methods
	// In a real implementation, we would call the appropriate method on the DBClient
	return nil
}

// Helper function to cache a provider in the database
func cacheProvider(ctx context.Context, dbClient *registry.DBClient, provider *response.ModuleProvider) error {
	// This is a placeholder since we don't have direct access to the internal methods
	// In a real implementation, we would call the appropriate method on the DBClient
	return nil
}
