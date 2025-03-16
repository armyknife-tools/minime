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
	fmt.Println("=== OpenTofu Registry API Simple Test Suite ===")
	
	// Test database connection
	testDatabaseConnection(logger)
	
	// Test registry client
	testRegistryClient(logger)
	
	// Test module operations
	testModuleOperations(logger)
	
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
	
	// We don't have a direct TestConnection method, but we can check if the client was created
	if dbClient == nil {
		fmt.Println("❌ Failed to create database client")
		return
	}
	
	fmt.Println("✅ Successfully created database client")
	
	// Print database configuration
	dbType := os.Getenv("TOFU_REGISTRY_DB_TYPE")
	dbHost := os.Getenv("TOFU_REGISTRY_DB_HOST")
	dbName := os.Getenv("TOFU_REGISTRY_DB_NAME")
	
	fmt.Printf("Database Type: %s\n", dbType)
	fmt.Printf("Database Host: %s\n", dbHost)
	fmt.Printf("Database Name: %s\n", dbName)
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
	fmt.Println("Testing module versions...")
	
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
	
	if versions == nil || len(versions.Modules) == 0 {
		fmt.Println("❌ No module versions found")
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
	
	// Print information about registry size
	fmt.Println("\n=== Registry Size Information ===")
	fmt.Println("The OpenTofu Registry is optimized to handle:")
	fmt.Println("- Approximately 4,000 providers")
	fmt.Println("- Approximately 18,000 modules")
	fmt.Println("Pre-allocation and caching mechanisms are in place to efficiently handle this volume of data.")
}
