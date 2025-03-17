// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/opentofu/opentofu/internal/templates"
)

// TemplateDB handles database operations for templates
type TemplateDB struct {
	DB     *sql.DB
	DBType string
}

// Template represents a resource template
type Template struct {
	ID          int
	Provider    string
	Resource    string
	DisplayName string
	Content     string
	Description string
	Category    string
	Tags        string
}

// GetTemplateDB returns a TemplateDB instance for the specified database type
func GetTemplateDB(dbType string, ui cli.Ui) (*TemplateDB, error) {
	var dbPath string
	
	// Load environment variables from .env file if it exists
	envFile := ".env"
	if _, err := os.Stat(envFile); err == nil {
		ui.Output("Loading environment variables from .env file...")
		envContent, err := os.ReadFile(envFile)
		if err == nil {
			lines := strings.Split(string(envContent), "\n")
			for _, line := range lines {
				line = strings.TrimSpace(line)
				if line == "" || strings.HasPrefix(line, "#") {
					continue
				}
				parts := strings.SplitN(line, "=", 2)
				if len(parts) == 2 {
					os.Setenv(parts[0], parts[1])
					ui.Output(fmt.Sprintf("Set environment variable: %s", parts[0]))
				}
			}
		}
	}
	
	if dbType == "sqlite" {
		configDir := filepath.Join(os.Getenv("HOME"), ".opentofu")
		if err := os.MkdirAll(configDir, 0755); err != nil {
			return nil, fmt.Errorf("error creating config directory: %v", err)
		}
		dbPath = filepath.Join(configDir, "templates.db")
	}
	
	// Try to connect to the database
	ui.Output(fmt.Sprintf("Connecting to %s database...", dbType))
	db, err := templates.ConnectToDatabase(dbType, dbPath)
	if err != nil {
		ui.Error(fmt.Sprintf("Error connecting to %s database: %v", dbType, err))
		if dbType == "postgres" {
			ui.Output("Falling back to SQLite database...")
			return GetTemplateDB("sqlite", ui)
		}
		return nil, err
	}
	
	return &TemplateDB{
		DB:     db,
		DBType: dbType,
	}, nil
}

// Close closes the database connection
func (tdb *TemplateDB) Close() error {
	return tdb.DB.Close()
}

// GetTemplate retrieves a template from the database
func (tdb *TemplateDB) GetTemplate(provider, resource string) (*Template, error) {
	query := `
		SELECT id, provider, resource, display_name, content, description, category, tags
		FROM templates
		WHERE provider = $1 AND resource = $2
	`

	var template Template
	err := tdb.DB.QueryRow(query, provider, resource).Scan(
		&template.ID,
		&template.Provider,
		&template.Resource,
		&template.DisplayName,
		&template.Content,
		&template.Description,
		&template.Category,
		&template.Tags,
	)

	if err != nil {
		return nil, err
	}

	return &template, nil
}

// ListTemplates returns all templates
func (tdb *TemplateDB) ListTemplates() ([]Template, error) {
	query := `
		SELECT id, provider, resource, display_name, description, category, tags
		FROM templates
		ORDER BY provider, resource
	`

	rows, err := tdb.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var templates []Template
	for rows.Next() {
		var template Template
		err := rows.Scan(
			&template.ID,
			&template.Provider,
			&template.Resource,
			&template.DisplayName,
			&template.Description,
			&template.Category,
			&template.Tags,
		)
		if err != nil {
			return nil, err
		}
		templates = append(templates, template)
	}

	return templates, nil
}

// ListProviderTemplates returns all templates for a specific provider
func (tdb *TemplateDB) ListProviderTemplates(provider string) ([]Template, error) {
	query := `
		SELECT id, provider, resource, display_name, description, category, tags
		FROM templates
		WHERE provider = $1
		ORDER BY resource
	`

	rows, err := tdb.DB.Query(query, provider)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var templates []Template
	for rows.Next() {
		var template Template
		err := rows.Scan(
			&template.ID,
			&template.Provider,
			&template.Resource,
			&template.DisplayName,
			&template.Description,
			&template.Category,
			&template.Tags,
		)
		if err != nil {
			return nil, err
		}
		templates = append(templates, template)
	}

	return templates, nil
}

// SaveTemplate saves a template to the database
func (tdb *TemplateDB) SaveTemplate(template *Template) error {
	query := `
		INSERT INTO templates (provider, resource, display_name, content, description, category, tags)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (provider, resource) DO UPDATE
		SET display_name = $3, content = $4, description = $5, category = $6, tags = $7, updated_at = CURRENT_TIMESTAMP
	`

	_, err := tdb.DB.Exec(
		query,
		template.Provider,
		template.Resource,
		template.DisplayName,
		template.Content,
		template.Description,
		template.Category,
		template.Tags,
	)

	return err
}

// DeleteTemplate deletes a template from the database
func (tdb *TemplateDB) DeleteTemplate(provider, resource string) error {
	query := `DELETE FROM templates WHERE provider = $1 AND resource = $2`
	_, err := tdb.DB.Exec(query, provider, resource)
	return err
}

// GetProviders returns a list of all providers that have templates
func (tdb *TemplateDB) GetProviders() ([]string, error) {
	query := `
		SELECT DISTINCT provider
		FROM templates
		ORDER BY provider
	`

	rows, err := tdb.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var providers []string
	for rows.Next() {
		var provider string
		if err := rows.Scan(&provider); err != nil {
			return nil, err
		}
		providers = append(providers, provider)
	}

	return providers, nil
}

// GetResources returns a list of all resources for a specific provider
func (tdb *TemplateDB) GetResources(provider string) ([]string, error) {
	query := `
		SELECT resource
		FROM templates
		WHERE provider = $1
		ORDER BY resource
	`

	rows, err := tdb.DB.Query(query, provider)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var resources []string
	for rows.Next() {
		var resource string
		if err := rows.Scan(&resource); err != nil {
			return nil, err
		}
		resources = append(resources, resource)
	}

	return resources, nil
}

// GetTemplateContent returns the content of a specific template
func (tdb *TemplateDB) GetTemplateContent(provider, resource string) (string, error) {
	query := `
		SELECT content
		FROM templates
		WHERE provider = $1 AND resource = $2
	`

	var content string
	err := tdb.DB.QueryRow(query, provider, resource).Scan(&content)
	if err != nil {
		return "", err
	}

	return content, nil
}
