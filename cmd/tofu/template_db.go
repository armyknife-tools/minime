// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

// TemplateDB handles database operations for templates
type TemplateDB struct {
	db *sql.DB
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

// NewTemplateDB creates a new TemplateDB instance
func NewTemplateDB() (*TemplateDB, error) {
	// Try PostgreSQL first if environment variables are set
	if os.Getenv("POSTGRES_HOST") != "" || fileExists(".env") {
		db, err := connectToPostgres()
		if err == nil {
			return &TemplateDB{db: db}, nil
		}
		// If PostgreSQL connection fails, fall back to SQLite
		fmt.Println("Warning: Could not connect to PostgreSQL, falling back to SQLite")
	}

	// Use SQLite as fallback
	return connectToSQLite()
}

// connectToPostgres connects to PostgreSQL using environment variables
func connectToPostgres() (*sql.DB, error) {
	// Load .env file if it exists
	if fileExists(".env") {
		godotenv.Load()
	}

	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")

	if host == "" {
		host = "localhost"
	}
	if port == "" {
		port = "5432"
	}
	if user == "" {
		user = "postgres"
	}
	if dbname == "" {
		dbname = "opentofu"
	}

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	// Create templates table if it doesn't exist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS templates (
			id SERIAL PRIMARY KEY,
			provider VARCHAR(50) NOT NULL,
			resource VARCHAR(100) NOT NULL,
			display_name VARCHAR(200) NOT NULL,
			content TEXT NOT NULL,
			description TEXT,
			category VARCHAR(100),
			tags VARCHAR(200),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(provider, resource)
		)
	`)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// connectToSQLite connects to a local SQLite database
func connectToSQLite() (*TemplateDB, error) {
	// Create the directory if it doesn't exist
	configDir := filepath.Join(os.Getenv("HOME"), ".opentofu")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, err
	}

	dbPath := filepath.Join(configDir, "templates.db")
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	// Create templates table if it doesn't exist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS templates (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			provider TEXT NOT NULL,
			resource TEXT NOT NULL,
			display_name TEXT NOT NULL,
			content TEXT NOT NULL,
			description TEXT,
			category TEXT,
			tags TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(provider, resource)
		)
	`)
	if err != nil {
		return nil, err
	}

	return &TemplateDB{db: db}, nil
}

// Close closes the database connection
func (tdb *TemplateDB) Close() error {
	return tdb.db.Close()
}

// GetTemplate retrieves a template from the database
func (tdb *TemplateDB) GetTemplate(provider, resource string) (*Template, error) {
	query := `
		SELECT id, provider, resource, display_name, content, description, category, tags
		FROM templates
		WHERE provider = $1 AND resource = $2
	`

	var template Template
	err := tdb.db.QueryRow(query, provider, resource).Scan(
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

	rows, err := tdb.db.Query(query)
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

	rows, err := tdb.db.Query(query, provider)
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

	_, err := tdb.db.Exec(
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
	_, err := tdb.db.Exec(query, provider, resource)
	return err
}

// GetProviders returns a list of all providers that have templates
func (tdb *TemplateDB) GetProviders() ([]string, error) {
	query := `
		SELECT DISTINCT provider
		FROM templates
		ORDER BY provider
	`

	rows, err := tdb.db.Query(query)
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

	rows, err := tdb.db.Query(query, provider)
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
	err := tdb.db.QueryRow(query, provider, resource).Scan(&content)
	if err != nil {
		return "", err
	}

	return content, nil
}

// Helper function to check if a file exists
func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}
