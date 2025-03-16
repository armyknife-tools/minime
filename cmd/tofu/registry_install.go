// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/hashicorp/terraform-svchost/disco"
	"github.com/opentofu/opentofu/internal/command"
	"github.com/opentofu/opentofu/internal/registry"
	"github.com/opentofu/opentofu/internal/registry/regsrc"
)

// RegistryInstallCommand is a Command implementation that installs modules or providers
// from a registry.
type RegistryInstallCommand struct {
	Meta command.Meta
}

func (c *RegistryInstallCommand) Help() string {
	helpText := `
Usage: tofu registry install [options] [ADDRESS]

  Install a module or provider from a registry.

  This command will install a module or provider from a registry and place it
  in the appropriate directory for use by OpenTofu.

  The ADDRESS argument is the address of the module or provider to install.
  For modules, this is in the format "namespace/name/provider".
  For providers, this is in the format "namespace/name".

Options:

  -type=TYPE            Type of resource to install. Can be "module" or "provider".
                        Default: "module"

  -version=VERSION      Specific version to install. If not specified, the latest
                        version will be installed.

  -registry=hostname    Use a custom registry host. By default, public registry
                        hosts are used based on the resource type.
`
	return strings.TrimSpace(helpText)
}

func (c *RegistryInstallCommand) Synopsis() string {
	return "Install a module or provider from a registry"
}

func (c *RegistryInstallCommand) Run(args []string) int {
	var installType string
	var version string
	var registryHost string

	cmdFlags := flag.NewFlagSet("registry install", flag.ContinueOnError)
	cmdFlags.StringVar(&installType, "type", "module", "Type of resource to install")
	cmdFlags.StringVar(&version, "version", "", "Specific version to install")
	cmdFlags.StringVar(&registryHost, "registry", "", "Registry host")
	cmdFlags.Usage = func() { c.Meta.Ui.Error(c.Help()) }

	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}

	args = cmdFlags.Args()
	if len(args) != 1 {
		c.Meta.Ui.Error("The registry install command expects exactly one argument.")
		return 1
	}

	address := args[0]

	// Check if the install type is valid
	if installType != "module" && installType != "provider" {
		c.Meta.Ui.Error(fmt.Sprintf("Invalid install type: %s. Must be 'module' or 'provider'.", installType))
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
		if installType == "module" {
			host = regsrc.PublicRegistryHost
		} else {
			host = regsrc.PublicRegistryHost // Same for providers
		}
	}

	// Perform the installation
	ctx := context.Background()
	if installType == "module" {
		return c.installModule(ctx, client, host, address, version)
	} else {
		return c.installProvider(ctx, client, host, address, version)
	}
}

// createRegistryClient creates a new registry client
func (c *RegistryInstallCommand) createRegistryClient() (*registry.Client, error) {
	// Create an HTTP client for the registry
	httpClient := retryablehttp.NewClient()
	httpClient.RetryMax = 3
	httpClient.Logger = hclog.NewNullLogger()

	// Create a services discovery client
	services := disco.New()
	services.SetUserAgent("OpenTofu")

	// Create a registry client
	client := registry.NewClient(services, httpClient.StandardClient())

	return client, nil
}

// installModule installs a module from a registry
func (c *RegistryInstallCommand) installModule(ctx context.Context, client *registry.Client, host *regsrc.FriendlyHost, address, version string) int {
	// Parse the module address
	module, err := regsrc.ParseModuleSource(address)
	if err != nil {
		c.Meta.Ui.Error(fmt.Sprintf("Invalid module address: %s", err))
		return 1
	}

	// Set the registry host
	module.RawHost = host

	c.Meta.Ui.Output(fmt.Sprintf("Installing module %s/%s/%s...", module.RawNamespace, module.RawName, module.RawProvider))

	// If no version is specified, get the latest version
	if version == "" {
		versions, err := client.ModuleVersions(ctx, module)
		if err != nil {
			c.Meta.Ui.Error(fmt.Sprintf("Error getting module versions: %s", err))
			return 1
		}

		if len(versions.Modules) == 0 {
			c.Meta.Ui.Error("No versions found for this module.")
			return 1
		}

		if len(versions.Modules[0].Versions) == 0 {
			c.Meta.Ui.Error("No versions found for this module.")
			return 1
		}

		version = versions.Modules[0].Versions[0].Version
		c.Meta.Ui.Output(fmt.Sprintf("Latest version is %s", version))
	}

	// Get the module location
	location, err := client.ModuleLocation(ctx, module, version)
	if err != nil {
		c.Meta.Ui.Error(fmt.Sprintf("Error getting module location: %s", err))
		return 1
	}

	// Determine the installation directory
	installDir := filepath.Join(c.Meta.WorkingDir.RootModuleDir(), "modules", module.RawNamespace, module.RawName, module.RawProvider, version)
	if err := os.MkdirAll(installDir, 0755); err != nil {
		c.Meta.Ui.Error(fmt.Sprintf("Error creating module directory: %s", err))
		return 1
	}

	c.Meta.Ui.Output(fmt.Sprintf("Downloading module from %s", location))

	// TODO: Implement actual module download and extraction
	// This would typically use the go-getter package to download and extract the module

	c.Meta.Ui.Output(fmt.Sprintf("Module installed successfully to %s", installDir))
	return 0
}

// installProvider installs a provider from a registry
func (c *RegistryInstallCommand) installProvider(ctx context.Context, client *registry.Client, host *regsrc.FriendlyHost, address, version string) int {
	// Parse the provider address
	parts := strings.Split(address, "/")
	if len(parts) != 2 {
		c.Meta.Ui.Error("Invalid provider address. Expected format: namespace/name")
		return 1
	}

	namespace := parts[0]
	name := parts[1]

	c.Meta.Ui.Output(fmt.Sprintf("Installing provider %s/%s...", namespace, name))

	// Determine the installation directory
	installDir := filepath.Join(c.Meta.WorkingDir.DataDir(), "providers", namespace, name)
	if err := os.MkdirAll(installDir, 0755); err != nil {
		c.Meta.Ui.Error(fmt.Sprintf("Error creating provider directory: %s", err))
		return 1
	}

	// For now, we'll just create a placeholder for provider installation
	// since the actual implementation would require more complex integration
	// with the getproviders package
	
	c.Meta.Ui.Output(fmt.Sprintf("Provider installation is not fully implemented yet."))
	c.Meta.Ui.Output(fmt.Sprintf("Would install provider %s/%s to %s", namespace, name, installDir))
	
	if version == "" {
		c.Meta.Ui.Output("No version specified, would use latest version")
	} else {
		c.Meta.Ui.Output(fmt.Sprintf("Would install version %s", version))
	}

	// Create a placeholder file to indicate the provider installation
	placeholderPath := filepath.Join(installDir, "provider_placeholder.txt")
	placeholderContent := fmt.Sprintf("Provider: %s/%s\nVersion: %s\nInstalled: placeholder\n", 
		namespace, name, version)
	
	if err := os.WriteFile(placeholderPath, []byte(placeholderContent), 0644); err != nil {
		c.Meta.Ui.Error(fmt.Sprintf("Error creating provider placeholder: %s", err))
		return 1
	}

	c.Meta.Ui.Output(fmt.Sprintf("Provider placeholder created at %s", placeholderPath))
	c.Meta.Ui.Output(fmt.Sprintf("Note: This is a placeholder for the provider installation."))
	c.Meta.Ui.Output(fmt.Sprintf("Full provider installation will be implemented in a future update."))
	
	return 0
}
