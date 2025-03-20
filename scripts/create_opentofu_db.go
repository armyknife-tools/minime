// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2023 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/lib/pq"
)

// Database schema definitions
const (
	modulesTableSQL = `
CREATE TABLE IF NOT EXISTS modules (
    id SERIAL PRIMARY KEY,
    namespace VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    provider VARCHAR(255) NOT NULL,
    version VARCHAR(50) NOT NULL,
    description TEXT,
    source_url TEXT,
    published_at TIMESTAMP,
    downloads INTEGER DEFAULT 0,
    UNIQUE(namespace, name, provider, version)
);
CREATE INDEX IF NOT EXISTS idx_modules_namespace ON modules(namespace);
CREATE INDEX IF NOT EXISTS idx_modules_name ON modules(name);
CREATE INDEX IF NOT EXISTS idx_modules_provider ON modules(provider);
CREATE INDEX IF NOT EXISTS idx_modules_version ON modules(version);
`

	providersTableSQL = `
CREATE TABLE IF NOT EXISTS providers (
    id SERIAL PRIMARY KEY,
    namespace VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    version VARCHAR(50) NOT NULL,
    description TEXT,
    source_url TEXT,
    published_at TIMESTAMP,
    downloads INTEGER DEFAULT 0,
    UNIQUE(namespace, name, version)
);
CREATE INDEX IF NOT EXISTS idx_providers_namespace ON providers(namespace);
CREATE INDEX IF NOT EXISTS idx_providers_name ON providers(name);
CREATE INDEX IF NOT EXISTS idx_providers_version ON providers(version);
`

	cacheTableSQL = `
CREATE TABLE IF NOT EXISTS cache (
    id SERIAL PRIMARY KEY,
    cache_key VARCHAR(255) NOT NULL UNIQUE,
    cache_value TEXT NOT NULL,
    expires_at TIMESTAMP NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_cache_expires_at ON cache(expires_at);
`
)

func main() {
	// Get database connection parameters from environment variables or use defaults
	host := getEnv("PGHOST", "localhost")
	port := getEnv("PGPORT", "5432")
	user := getEnv("PGUSER", "postgres")
	password := getEnv("PGPASSWORD", "")
	dbname := getEnv("PGDATABASE", "opentofu")

	// Connect to PostgreSQL server
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s sslmode=disable", 
		host, port, user, password)
	
	fmt.Println("Connecting to PostgreSQL server...")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL server: %v", err)
	}
	defer db.Close()

	// Check connection
	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to ping PostgreSQL server: %v", err)
	}
	fmt.Println("Successfully connected to PostgreSQL server")

	// Check if database exists
	var exists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname=$1)", dbname).Scan(&exists)
	if err != nil {
		log.Fatalf("Failed to check if database exists: %v", err)
	}

	// Create database if it doesn't exist
	if !exists {
		fmt.Printf("Creating database '%s'...\n", dbname)
		_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", dbname))
		if err != nil {
			log.Fatalf("Failed to create database: %v", err)
		}
		fmt.Printf("Database '%s' created successfully\n", dbname)
	} else {
		fmt.Printf("Database '%s' already exists\n", dbname)
	}

	// Connect to the newly created database
	connStr = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", 
		host, port, user, password, dbname)
	
	fmt.Printf("Connecting to database '%s'...\n", dbname)
	dbConn, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbConn.Close()

	// Create tables
	fmt.Println("Creating tables...")
	
	// Create modules table
	_, err = dbConn.Exec(modulesTableSQL)
	if err != nil {
		log.Fatalf("Failed to create modules table: %v", err)
	}
	fmt.Println("Modules table created successfully")

	// Create providers table
	_, err = dbConn.Exec(providersTableSQL)
	if err != nil {
		log.Fatalf("Failed to create providers table: %v", err)
	}
	fmt.Println("Providers table created successfully")

	// Create cache table
	_, err = dbConn.Exec(cacheTableSQL)
	if err != nil {
		log.Fatalf("Failed to create cache table: %v", err)
	}
	fmt.Println("Cache table created successfully")

	fmt.Println("\nDatabase setup completed successfully!")
	fmt.Println("You can now create a dedicated database user with the create_opentofu_user.go script.")
}

// getEnv retrieves an environment variable or returns a default value if not set
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return strings.TrimSpace(value)
}
