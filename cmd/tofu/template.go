// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/opentofu/opentofu/internal/command"
	"github.com/opentofu/opentofu/internal/templates"
)

// TemplateCommand is a command for generating cloud resource templates
type TemplateCommand struct {
	Meta command.Meta
}

func (c *TemplateCommand) Help() string {
	helpText := `
Usage: tofu template [options] PROVIDER/RESOURCE

  This command generates a template for a cloud resource.

  If PROVIDER is specified without RESOURCE, it will list all available 
  resources for that provider.

  If neither PROVIDER nor RESOURCE is specified, it will list all available
  providers.

Options:

  -db=TYPE     Database type to use (sqlite or postgres). Default: sqlite
  -load        Load templates into the database
`
	return strings.TrimSpace(helpText)
}

func (c *TemplateCommand) Synopsis() string {
	return "Generate templates for cloud resources"
}

func (c *TemplateCommand) Run(args []string) int {
	// Parse command-line flags
	cmdFlags := flag.NewFlagSet("template", flag.ContinueOnError)
	cmdFlags.Usage = func() { c.Meta.Ui.Error(c.Help()) }
	dbTypeFlag := cmdFlags.String("db", "sqlite", "Database type: sqlite or postgres")
	loadFlag := cmdFlags.Bool("load", false, "Load templates into the database")
	
	if err := cmdFlags.Parse(args); err != nil {
		c.Meta.Ui.Error(fmt.Sprintf("Error parsing command-line flags: %s", err))
		return 1
	}

	// If the load flag is set, load templates into the database
	if *loadFlag {
		c.Meta.Ui.Output("Loading templates into the database...")
		
		var dbPath string
		if *dbTypeFlag == "sqlite" {
			configDir := filepath.Join(os.Getenv("HOME"), ".opentofu")
			if err := os.MkdirAll(configDir, 0755); err != nil {
				c.Meta.Ui.Error(fmt.Sprintf("Error creating config directory: %v", err))
				return 1
			}
			dbPath = filepath.Join(configDir, "templates.db")
		}
		
		err := templates.LoadTemplates(*dbTypeFlag, dbPath)
		if err != nil {
			c.Meta.Ui.Error(fmt.Sprintf("Error loading templates: %v", err))
			return 1
		}
		
		c.Meta.Ui.Output("Templates loaded successfully!")
		return 0
	}

	// Connect to the template database
	db, err := NewTemplateDB()
	if err != nil {
		c.Meta.Ui.Error(fmt.Sprintf("Error connecting to template database: %s", err))
		return 1
	}
	defer db.Close()

	// Parse remaining arguments
	args = cmdFlags.Args()
	if len(args) == 0 {
		// List all providers
		providers, err := db.GetProviders()
		if err != nil {
			c.Meta.Ui.Error(fmt.Sprintf("Error retrieving providers: %s", err))
			return 1
		}

		c.Meta.Ui.Output("Available providers:")
		for _, provider := range providers {
			c.Meta.Ui.Output(fmt.Sprintf("  %s", provider))
		}
		return 0
	}

	// Parse provider/resource
	parts := strings.Split(args[0], "/")
	provider := parts[0]

	if len(parts) == 1 {
		// List resources for the provider
		resources, err := db.GetResources(provider)
		if err != nil {
			c.Meta.Ui.Error(fmt.Sprintf("Error retrieving resources for provider %s: %s", provider, err))
			return 1
		}

		if len(resources) == 0 {
			c.Meta.Ui.Error(fmt.Sprintf("No resources found for provider %s", provider))
			return 1
		}

		c.Meta.Ui.Output(fmt.Sprintf("Available resources for %s:", provider))
		for _, resource := range resources {
			c.Meta.Ui.Output(fmt.Sprintf("  %s", resource))
		}
		return 0
	}

	resource := parts[1]

	// Get the template content
	content, err := db.GetTemplateContent(provider, resource)
	if err != nil {
		c.Meta.Ui.Error(fmt.Sprintf("Error retrieving template for %s/%s: %s", provider, resource, err))
		return 1
	}

	if content == "" {
		c.Meta.Ui.Error(fmt.Sprintf("No template found for %s/%s", provider, resource))
		return 1
	}

	// Write the template to a file
	filename := fmt.Sprintf("%s.tf", resource)
	err = os.WriteFile(filename, []byte(content), 0644)
	if err != nil {
		c.Meta.Ui.Error(fmt.Sprintf("Error writing template to file: %s", err))
		return 1
	}

	c.Meta.Ui.Output(fmt.Sprintf("Template for %s/%s written to %s", provider, resource, filename))
	return 0
}
