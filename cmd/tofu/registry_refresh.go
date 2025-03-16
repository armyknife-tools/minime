// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	"github.com/hashicorp/go-hclog"
	svchost "github.com/hashicorp/terraform-svchost"
	"github.com/hashicorp/terraform-svchost/disco"
	"github.com/opentofu/opentofu/internal/command"
	"github.com/opentofu/opentofu/internal/command/cliconfig"
	"github.com/opentofu/opentofu/internal/httpclient"
	"github.com/opentofu/opentofu/internal/registry"
	"github.com/opentofu/opentofu/version"
)

// RegistryRefreshCommand is a CLI command for refreshing the registry cache
type RegistryRefreshCommand struct {
	Meta command.Meta
}

func (c *RegistryRefreshCommand) Help() string {
	helpText := `
Usage: tofu registry refresh [options]

  Refreshes the local cache of registry modules and providers.

Options:

  -hosts=hostname,...     Comma-separated list of registry hostnames to refresh 
                          (default: registry.terraform.io,registry.opentofu.org)
  -cache-dir=path         Directory to store registry cache files
                          (default: $HOME/.terraform.d/registry-cache)
  -refresh-interval=6h    Interval between registry refreshes when running as daemon
  -cleanup-interval=24h   Interval between cache cleanups when running as daemon
  -max-age=168h           Maximum age of cache files before deletion (default: 7 days)
  -daemon                 Run continuously as a background process
`
	return strings.TrimSpace(helpText)
}

func (c *RegistryRefreshCommand) Synopsis() string {
	return "Refresh local registry module and provider cache"
}

func (c *RegistryRefreshCommand) Run(args []string) int {
	var hostsFlag string
	var cacheDirFlag string
	var refreshIntervalFlag time.Duration
	var cleanupIntervalFlag time.Duration
	var maxAgeFlag time.Duration
	var daemonFlag bool

	flags := flag.NewFlagSet("registry refresh", flag.ContinueOnError)
	flags.StringVar(&hostsFlag, "hosts", "registry.terraform.io,registry.opentofu.org", "Comma-separated list of registry hostnames")
	flags.StringVar(&cacheDirFlag, "cache-dir", "", "Directory to store registry cache files")
	flags.DurationVar(&refreshIntervalFlag, "refresh-interval", 6*time.Hour, "Interval between registry refreshes when running as daemon")
	flags.DurationVar(&cleanupIntervalFlag, "cleanup-interval", 24*time.Hour, "Interval between cache cleanups when running as daemon")
	flags.DurationVar(&maxAgeFlag, "max-age", 168*time.Hour, "Maximum age of cache files before deletion")
	flags.BoolVar(&daemonFlag, "daemon", false, "Run continuously as a background process")

	flags.Usage = func() { c.Meta.Ui.Error(c.Help()) }

	if err := flags.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return 0
		}
		c.Meta.Ui.Error(fmt.Sprintf("Error parsing command line arguments: %s", err))
		return 1
	}

	// Use default cache directory if not specified
	if cacheDirFlag == "" {
		var err error
		cacheDirFlag, err = defaultCacheDir()
		if err != nil {
			c.Meta.Ui.Error(fmt.Sprintf("Error determining default cache directory: %s", err))
			return 1
		}
	}

	// Parse the hosts flag
	var hosts []svchost.Hostname
	if hostsFlag != "" {
		hostStrings := strings.Split(hostsFlag, ",")
		hosts = make([]svchost.Hostname, 0, len(hostStrings))

		for _, hostStr := range hostStrings {
			hostStr = strings.TrimSpace(hostStr)
			if hostStr == "" {
				continue
			}

			host, err := svchost.ForComparison(hostStr)
			if err != nil {
				c.Meta.Ui.Error(fmt.Sprintf("Invalid hostname %q: %s", hostStr, err))
				return 1
			}
			hosts = append(hosts, host)
		}
	}

	// Create a registry client
	services := disco.New()
	services.SetUserAgent(httpclient.OpenTofuUserAgent(version.String()))

	// Create a logger
	var logOutput io.Writer = io.Discard
	if c.Meta.Streams != nil && c.Meta.Streams.Stderr != nil {
		logOutput = c.Meta.Streams.Stderr.File
	}
	
	logger := hclog.New(&hclog.LoggerOptions{
		Name:   "registry-refresh",
		Level:  hclog.Info,
		Output: logOutput,
	})

	// Create a caching client
	cachingClient, err := registry.NewCachingClient(registry.NewClient(services, nil), cacheDirFlag, logger)
	if err != nil {
		c.Meta.Ui.Error(fmt.Sprintf("Error creating caching client: %s", err))
		return 1
	}

	// Run in daemon mode if requested
	if daemonFlag {
		ctx := context.Background()
		cachingClient.StartBackgroundRefresh(ctx, hosts, refreshIntervalFlag, cleanupIntervalFlag, maxAgeFlag)

		c.Meta.Ui.Output(fmt.Sprintf("Starting registry refresh daemon with refresh interval %s", refreshIntervalFlag))
		c.Meta.Ui.Output(fmt.Sprintf("Cache directory: %s", cacheDirFlag))
		c.Meta.Ui.Output(fmt.Sprintf("Hosts: %s", hostsFlag))
		c.Meta.Ui.Output("Press Ctrl+C to stop")

		// Block forever (until interrupted)
		select {}
	}

	// Otherwise, do a single refresh
	ctx := context.Background()

	c.Meta.Ui.Output(fmt.Sprintf("Refreshing registry metadata for %d hosts...", len(hosts)))

	for _, host := range hosts {
		c.Meta.Ui.Output(fmt.Sprintf("Refreshing modules for %s...", host))
		if err := cachingClient.RefreshModuleCache(ctx, host); err != nil {
			c.Meta.Ui.Error(fmt.Sprintf("Error refreshing modules for %s: %s", host, err))
		}

		c.Meta.Ui.Output(fmt.Sprintf("Refreshing providers for %s...", host))
		if err := cachingClient.RefreshProviderCache(ctx, host); err != nil {
			c.Meta.Ui.Error(fmt.Sprintf("Error refreshing providers for %s: %s", host, err))
		}
	}

	c.Meta.Ui.Output(fmt.Sprintf("Cleaning up cache files older than %s...", maxAgeFlag))
	if err := cachingClient.CleanupOldCacheFiles(maxAgeFlag); err != nil {
		c.Meta.Ui.Error(fmt.Sprintf("Error cleaning up cache files: %s", err))
	}

	c.Meta.Ui.Output("Registry refresh complete")
	return 0
}

// defaultCacheDir returns the default directory for registry cache files
func defaultCacheDir() (string, error) {
	configDir, err := cliconfig.ConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "registry-cache"), nil
}
