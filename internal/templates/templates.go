// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

// Package templates provides functionality for generating and loading cloud resource templates
package templates

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

// Template represents a cloud resource template
type Template struct {
	Provider    string
	Resource    string
	DisplayName string
	Description string
	Category    string
	Tags        string
	Content     string
}

// LoadTemplates loads templates into the specified database
func LoadTemplates(dbType, dbPath string) error {
	// Connect to the database
	db, err := connectToDatabase(dbType, dbPath)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}
	defer db.Close()

	// Create the templates table if it doesn't exist
	err = createTemplatesTable(db)
	if err != nil {
		return fmt.Errorf("failed to create templates table: %v", err)
	}

	// Generate templates for each provider
	templates := []Template{}
	templates = append(templates, generateAWSTemplates()...)
	templates = append(templates, generateAzureTemplates()...)
	templates = append(templates, generateGCPTemplates()...)

	// Insert templates into the database
	for _, template := range templates {
		err = insertTemplate(db, template)
		if err != nil {
			return fmt.Errorf("failed to insert template: %v", err)
		}
	}

	return nil
}

// connectToDatabase connects to either a SQLite or PostgreSQL database
func connectToDatabase(dbType, dbPath string) (*sql.DB, error) {
	var db *sql.DB
	var err error

	switch strings.ToLower(dbType) {
	case "sqlite":
		db, err = sql.Open("sqlite3", dbPath)
		if err != nil {
			return nil, fmt.Errorf("failed to open SQLite database: %v", err)
		}
	case "postgres":
		// Get PostgreSQL connection details from environment variables
		host := getEnv("POSTGRES_HOST", "localhost")
		port := getEnv("POSTGRES_PORT", "5432")
		user := getEnv("POSTGRES_USER", "postgres")
		password := getEnv("POSTGRES_PASSWORD", "postgres")
		dbname := getEnv("POSTGRES_DB", "opentofu")

		connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			host, port, user, password, dbname)
		db, err = sql.Open("postgres", connStr)
		if err != nil {
			return nil, fmt.Errorf("failed to open PostgreSQL database: %v", err)
		}
	default:
		return nil, fmt.Errorf("unsupported database type: %s", dbType)
	}

	// Test the connection
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	return db, nil
}

// createTemplatesTable creates the templates table if it doesn't exist
func createTemplatesTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS templates (
		id SERIAL PRIMARY KEY,
		provider TEXT NOT NULL,
		resource TEXT NOT NULL,
		display_name TEXT NOT NULL,
		description TEXT NOT NULL,
		category TEXT NOT NULL,
		tags TEXT NOT NULL,
		content TEXT NOT NULL,
		UNIQUE(provider, resource)
	);
	`
	_, err := db.Exec(query)
	return err
}

// insertTemplate inserts a template into the database
func insertTemplate(db *sql.DB, template Template) error {
	query := `
	INSERT INTO templates (provider, resource, display_name, description, category, tags, content)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	ON CONFLICT (provider, resource) DO UPDATE
	SET display_name = $3, description = $4, category = $5, tags = $6, content = $7;
	`
	_, err := db.Exec(query, template.Provider, template.Resource, template.DisplayName,
		template.Description, template.Category, template.Tags, template.Content)
	return err
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
