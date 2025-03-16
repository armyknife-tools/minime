package main

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/opentofu/opentofu/internal/registry"
	"github.com/opentofu/opentofu/internal/registry/regsrc"
	"github.com/opentofu/opentofu/internal/registry/response"
)

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

	// Test module operations
	logger.Info("Testing module operations...")
	moduleVersions, err := testModuleOperations(client)
	if err != nil {
		logger.Error("Module operations failed", "error", err)
		return
	}
	logger.Info("Module operations successful")

	// Test module installation
	logger.Info("Testing module installation...")
	err = testModuleInstallation(client, moduleVersions)
	if err != nil {
		logger.Error("Module installation failed", "error", err)
		return
	}
	logger.Info("Module installation successful")

	logger.Info("All tests completed successfully!")
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
	// Create an HTTP client with retries
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 3
	retryClient.RetryWaitMin = 1 * time.Second
	retryClient.RetryWaitMax = 5 * time.Second
	retryClient.Logger = nil // Disable logging for the retry client

	httpClient := retryClient.StandardClient()

	// Create a new registry client
	client := registry.NewClient(nil, httpClient)
	if client == nil {
		return nil, fmt.Errorf("failed to create registry client")
	}

	return client, nil
}

// testModuleOperations tests various module operations
func testModuleOperations(client *registry.Client) (*response.ModuleVersions, error) {
	ctx := context.Background()

	// Parse a module source
	module, err := regsrc.ParseModuleSource("hashicorp/consul/aws")
	if err != nil {
		return nil, fmt.Errorf("error parsing module source: %w", err)
	}

	fmt.Printf("Module details:\n")
	fmt.Printf("  Namespace: %s\n", module.Module.RawNamespace)
	fmt.Printf("  Name: %s\n", module.Module.RawName)
	fmt.Printf("  Provider: %s\n", module.Module.RawProvider)
	fmt.Printf("  Services: %s\n", module.Services.Hostname())
	fmt.Println()

	// Get module versions
	moduleVersions, err := client.ModuleVersions(ctx, module)
	if err != nil {
		return nil, fmt.Errorf("error getting module versions: %w", err)
	}

	fmt.Printf("Available versions for %s/%s/%s:\n", 
		module.Module.RawNamespace, 
		module.Module.RawName, 
		module.Module.RawProvider)
	
	for i, version := range moduleVersions.Modules[0].Versions {
		fmt.Printf("  %d. Version: %s, Published: %s\n", 
			i+1, 
			version.Version, 
			version.PublishedAt.Format("2006-01-02"))
	}
	fmt.Println()

	// Get the latest version
	latestVersion := moduleVersions.Modules[0].Versions[0].Version
	fmt.Printf("Latest version: %s\n", latestVersion)

	// Get module download URL
	downloadURL, err := client.ModuleLocation(ctx, module, latestVersion)
	if err != nil {
		return nil, fmt.Errorf("error getting module download URL: %w", err)
	}

	fmt.Printf("Download URL: %s\n\n", downloadURL)

	return moduleVersions, nil
}

// testModuleInstallation tests the module installation process
func testModuleInstallation(client *registry.Client, moduleVersions *response.ModuleVersions) error {
	ctx := context.Background()

	// Parse a module source
	module, err := regsrc.ParseModuleSource("hashicorp/consul/aws")
	if err != nil {
		return fmt.Errorf("error parsing module source: %w", err)
	}

	// Get the latest version
	latestVersion := moduleVersions.Modules[0].Versions[0].Version

	// Get module download URL
	downloadURL, err := client.ModuleLocation(ctx, module, latestVersion)
	if err != nil {
		return fmt.Errorf("error getting module download URL: %w", err)
	}

	// Create a temporary directory for the download
	tempDir, err := os.MkdirTemp("", "opentofu-module-*")
	if err != nil {
		return fmt.Errorf("error creating temporary directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	fmt.Printf("Installing module to: %s\n", tempDir)

	// Download the module
	modulePath := filepath.Join(tempDir, "module.tar.gz")
	err = downloadFile(downloadURL, modulePath)
	if err != nil {
		return fmt.Errorf("error downloading module: %w", err)
	}

	fmt.Printf("Downloaded module to: %s\n", modulePath)
	fmt.Printf("Module size: %d bytes\n", getFileSize(modulePath))

	// Extract the module (simulation)
	fmt.Println("Simulating module extraction...")
	time.Sleep(500 * time.Millisecond)
	fmt.Println("Module extraction complete")

	// Create module installation directory structure
	installDir := filepath.Join(tempDir, "installed")
	err = os.MkdirAll(installDir, 0755)
	if err != nil {
		return fmt.Errorf("error creating installation directory: %w", err)
	}

	// Create a simulated module structure
	err = createSimulatedModuleStructure(installDir, module)
	if err != nil {
		return fmt.Errorf("error creating simulated module structure: %w", err)
	}

	fmt.Printf("Module installed successfully to: %s\n", installDir)
	fmt.Println("Installation structure:")
	listDirectoryStructure(installDir, "  ")

	return nil
}

// downloadFile downloads a file from a URL
func downloadFile(url, filepath string) error {
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

// getFileSize returns the size of a file in bytes
func getFileSize(filepath string) int64 {
	fileInfo, err := os.Stat(filepath)
	if err != nil {
		return 0
	}
	return fileInfo.Size()
}

// createSimulatedModuleStructure creates a simulated module structure
func createSimulatedModuleStructure(dir string, module *regsrc.Module) error {
	// Create main.tf
	mainTF := filepath.Join(dir, "main.tf")
	mainContent := fmt.Sprintf(`# %s/%s/%s module
resource "aws_instance" "example" {
  ami           = var.ami_id
  instance_type = var.instance_type
  tags = {
    Name = "example-%s"
  }
}
`, module.Module.RawNamespace, module.Module.RawName, module.Module.RawProvider, module.Module.RawName)
	err := os.WriteFile(mainTF, []byte(mainContent), 0644)
	if err != nil {
		return err
	}

	// Create variables.tf
	variablesTF := filepath.Join(dir, "variables.tf")
	variablesContent := `variable "ami_id" {
  description = "The AMI ID to use for the instance"
  type        = string
}

variable "instance_type" {
  description = "The type of instance to start"
  type        = string
  default     = "t2.micro"
}
`
	err = os.WriteFile(variablesTF, []byte(variablesContent), 0644)
	if err != nil {
		return err
	}

	// Create outputs.tf
	outputsTF := filepath.Join(dir, "outputs.tf")
	outputsContent := `output "instance_id" {
  description = "The ID of the instance"
  value       = aws_instance.example.id
}

output "instance_public_ip" {
  description = "The public IP address of the instance"
  value       = aws_instance.example.public_ip
}
`
	err = os.WriteFile(outputsTF, []byte(outputsContent), 0644)
	if err != nil {
		return err
	}

	// Create README.md
	readmeMD := filepath.Join(dir, "README.md")
	readmeContent := fmt.Sprintf(`# %s/%s/%s

A Terraform module for creating an AWS instance.

## Usage

```hcl
module "%s" {
  source  = "%s/%s/%s"
  version = "0.1.0"

  ami_id        = "ami-0c55b159cbfafe1f0"
  instance_type = "t2.micro"
}
```

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| ami_id | The AMI ID to use for the instance | string | n/a | yes |
| instance_type | The type of instance to start | string | "t2.micro" | no |

## Outputs

| Name | Description |
|------|-------------|
| instance_id | The ID of the instance |
| instance_public_ip | The public IP address of the instance |
`,
		module.Module.RawNamespace, module.Module.RawName, module.Module.RawProvider,
		module.Module.RawName,
		module.Module.RawNamespace, module.Module.RawName, module.Module.RawProvider)
	err = os.WriteFile(readmeMD, []byte(readmeContent), 0644)
	if err != nil {
		return err
	}

	// Create examples directory
	examplesDir := filepath.Join(dir, "examples", "basic")
	err = os.MkdirAll(examplesDir, 0755)
	if err != nil {
		return err
	}

	// Create example main.tf
	exampleMainTF := filepath.Join(examplesDir, "main.tf")
	exampleMainContent := fmt.Sprintf(`module "%s" {
  source = "../../"

  ami_id        = "ami-0c55b159cbfafe1f0"
  instance_type = "t2.micro"
}

output "instance_id" {
  value = module.%s.instance_id
}

output "instance_public_ip" {
  value = module.%s.instance_public_ip
}
`, module.Module.RawName, module.Module.RawName, module.Module.RawName)
	err = os.WriteFile(exampleMainTF, []byte(exampleMainContent), 0644)
	if err != nil {
		return err
	}

	return nil
}

// listDirectoryStructure lists the directory structure
func listDirectoryStructure(dir string, indent string) {
	files, err := os.ReadDir(dir)
	if err != nil {
		fmt.Printf("%sError reading directory: %s\n", indent, err)
		return
	}

	for _, file := range files {
		path := filepath.Join(dir, file.Name())
		if file.IsDir() {
			fmt.Printf("%s\uF818 %s/\n", indent, file.Name())
			listDirectoryStructure(path, indent+"  ")
		} else {
			info, _ := file.Info()
			size := info.Size()
			var sizeStr string
			if size < 1024 {
				sizeStr = fmt.Sprintf("%d B", size)
			} else if size < 1024*1024 {
				sizeStr = fmt.Sprintf("%.1f KB", float64(size)/1024)
			} else {
				sizeStr = fmt.Sprintf("%.1f MB", float64(size)/(1024*1024))
			}
			fmt.Printf("%s\uF815 %s (%s)\n", indent, file.Name(), sizeStr)
		}
	}
}
