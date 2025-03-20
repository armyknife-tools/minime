// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2023 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

// Test data
const (
	testModuleNamespace = "test"
	testModuleName      = "connection-test"
	testModuleProvider  = "aws"
	testModuleVersion   = "1.0.0"
	testModuleDesc      = "Test module for connection testing"
	testModuleURL       = "https://github.com/opentofu/opentofu"

	testProviderNamespace = "test"
	testProviderName      = "connection-test"
	testProviderVersion   = "1.0.0"
	testProviderDesc      = "Test provider for connection testing"
	testProviderURL       = "https://github.com/opentofu/opentofu"

	testCacheKey   = "test-connection"
	testCacheValue = "Connection test successful"
)

func main() {
	// Load environment variables from .env file
	loadEnvFile()

	// Get database connection parameters from environment variables
	dbType := getEnv("TOFU_REGISTRY_DB_TYPE", "")
	
	var db *sql.DB
	var err error

	// Connect to database based on type
	switch strings.ToLower(dbType) {
	case "postgres":
		db, err = connectToPostgres()
	case "sqlite":
		db, err = connectToSQLite()
	default:
		log.Fatalf("Unsupported database type: %s. Supported types are 'postgres' and 'sqlite'", dbType)
	}

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	fmt.Printf("Successfully connected to %s database\n", dbType)

	// Test table creation
	fmt.Println("\nTesting table creation...")
	err = testTableCreation(db, dbType)
	if err != nil {
		log.Fatalf("Failed to create tables: %v", err)
	}
	fmt.Println("Table creation successful")

	// Test data insertion
	fmt.Println("\nTesting data insertion...")
	err = testDataInsertion(db, dbType)
	if err != nil {
		log.Fatalf("Failed to insert data: %v", err)
	}
	fmt.Println("Data insertion successful")

	// Test data retrieval
	fmt.Println("\nTesting data retrieval...")
	err = testDataRetrieval(db, dbType)
	if err != nil {
		log.Fatalf("Failed to retrieve data: %v", err)
	}
	fmt.Println("Data retrieval successful")

	// Test data deletion
	fmt.Println("\nTesting data deletion...")
	err = testDataDeletion(db, dbType)
	if err != nil {
		log.Fatalf("Failed to delete data: %v", err)
	}
	fmt.Println("Data deletion successful")

	fmt.Println("\nAll tests passed successfully!")
	fmt.Println("Your database connection is working correctly.")
}

// connectToPostgres connects to a PostgreSQL database
func connectToPostgres() (*sql.DB, error) {
	host := getEnv("TOFU_REGISTRY_DB_HOST", "localhost")
	port := getEnv("TOFU_REGISTRY_DB_PORT", "5432")
	user := getEnv("TOFU_REGISTRY_DB_USER", "opentofu_user")
	password := getEnv("TOFU_REGISTRY_DB_PASSWORD", "")
	dbname := getEnv("TOFU_REGISTRY_DB_NAME", "opentofu")
	sslmode := getEnv("TOFU_REGISTRY_DB_SSLMODE", "disable")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", 
		host, port, user, password, dbname, sslmode)
	
	fmt.Printf("Connecting to PostgreSQL database %s on %s:%s as %s...\n", 
		dbname, host, port, user)
	
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %v", err)
	}

	// Check connection
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	return db, nil
}

// connectToSQLite connects to a SQLite database
func connectToSQLite() (*sql.DB, error) {
	dbURL := getEnv("TOFU_REGISTRY_DB_URL", "file:registry.db?cache=shared&mode=rwc")
	
	// Extract file path from URL
	dbPath := strings.TrimPrefix(dbURL, "file:")
	dbPath = strings.Split(dbPath, "?")[0]
	
	fmt.Printf("Connecting to SQLite database at %s...\n", dbPath)
	
	db, err := sql.Open("sqlite3", dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %v", err)
	}

	// Check connection
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	return db, nil
}

// testTableCreation tests creating the necessary tables
func testTableCreation(db *sql.DB, dbType string) error {
	// SQL statements for creating tables
	var modulesTableSQL, providersTableSQL, cacheTableSQL string

	// Adjust SQL syntax based on database type
	if strings.ToLower(dbType) == "postgres" {
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
);`

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
);`

		cacheTableSQL = `
CREATE TABLE IF NOT EXISTS cache (
    id SERIAL PRIMARY KEY,
    cache_key VARCHAR(255) NOT NULL UNIQUE,
    cache_value TEXT NOT NULL,
    expires_at TIMESTAMP NOT NULL
);`
	} else {
		// SQLite syntax
		modulesTableSQL = `
CREATE TABLE IF NOT EXISTS modules (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    namespace TEXT NOT NULL,
    name TEXT NOT NULL,
    provider TEXT NOT NULL,
    version TEXT NOT NULL,
    description TEXT,
    source_url TEXT,
    published_at TIMESTAMP,
    downloads INTEGER DEFAULT 0,
    UNIQUE(namespace, name, provider, version)
);`

		providersTableSQL = `
CREATE TABLE IF NOT EXISTS providers (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    namespace TEXT NOT NULL,
    name TEXT NOT NULL,
    version TEXT NOT NULL,
    description TEXT,
    source_url TEXT,
    published_at TIMESTAMP,
    downloads INTEGER DEFAULT 0,
    UNIQUE(namespace, name, version)
);`

		cacheTableSQL = `
CREATE TABLE IF NOT EXISTS cache (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    cache_key TEXT NOT NULL UNIQUE,
    cache_value TEXT NOT NULL,
    expires_at TIMESTAMP NOT NULL
);`
	}

	// Create modules table
	_, err := db.Exec(modulesTableSQL)
	if err != nil {
		return fmt.Errorf("failed to create modules table: %v", err)
	}

	// Create providers table
	_, err = db.Exec(providersTableSQL)
	if err != nil {
		return fmt.Errorf("failed to create providers table: %v", err)
	}

	// Create cache table
	_, err = db.Exec(cacheTableSQL)
	if err != nil {
		return fmt.Errorf("failed to create cache table: %v", err)
	}

	return nil
}

// testDataInsertion tests inserting data into the tables
func testDataInsertion(db *sql.DB, dbType string) error {
	// Insert test module
	_, err := db.Exec(`
INSERT INTO modules (namespace, name, provider, version, description, source_url, published_at)
VALUES ($1, $2, $3, $4, $5, $6, $7)
ON CONFLICT (namespace, name, provider, version) DO UPDATE
SET description = $5, source_url = $6, published_at = $7`,
		testModuleNamespace, testModuleName, testModuleProvider, testModuleVersion,
		testModuleDesc, testModuleURL, time.Now())
	if err != nil {
		return fmt.Errorf("failed to insert test module: %v", err)
	}

	// Insert test provider
	_, err = db.Exec(`
INSERT INTO providers (namespace, name, version, description, source_url, published_at)
VALUES ($1, $2, $3, $4, $5, $6)
ON CONFLICT (namespace, name, version) DO UPDATE
SET description = $4, source_url = $5, published_at = $6`,
		testProviderNamespace, testProviderName, testProviderVersion,
		testProviderDesc, testProviderURL, time.Now())
	if err != nil {
		return fmt.Errorf("failed to insert test provider: %v", err)
	}

	// Insert test cache entry
	_, err = db.Exec(`
INSERT INTO cache (cache_key, cache_value, expires_at)
VALUES ($1, $2, $3)
ON CONFLICT (cache_key) DO UPDATE
SET cache_value = $2, expires_at = $3`,
		testCacheKey, testCacheValue, time.Now().Add(24*time.Hour))
	if err != nil {
		return fmt.Errorf("failed to insert test cache entry: %v", err)
	}

	return nil
}

// testDataRetrieval tests retrieving data from the tables
func testDataRetrieval(db *sql.DB, dbType string) error {
	// Retrieve test module
	var moduleCount int
	err := db.QueryRow(`
SELECT COUNT(*) FROM modules
WHERE namespace = $1 AND name = $2 AND provider = $3 AND version = $4`,
		testModuleNamespace, testModuleName, testModuleProvider, testModuleVersion).Scan(&moduleCount)
	if err != nil {
		return fmt.Errorf("failed to retrieve test module: %v", err)
	}
	if moduleCount != 1 {
		return fmt.Errorf("expected 1 test module, got %d", moduleCount)
	}

	// Retrieve test provider
	var providerCount int
	err = db.QueryRow(`
SELECT COUNT(*) FROM providers
WHERE namespace = $1 AND name = $2 AND version = $3`,
		testProviderNamespace, testProviderName, testProviderVersion).Scan(&providerCount)
	if err != nil {
		return fmt.Errorf("failed to retrieve test provider: %v", err)
	}
	if providerCount != 1 {
		return fmt.Errorf("expected 1 test provider, got %d", providerCount)
	}

	// Retrieve test cache entry
	var cacheValue string
	err = db.QueryRow(`
SELECT cache_value FROM cache
WHERE cache_key = $1`,
		testCacheKey).Scan(&cacheValue)
	if err != nil {
		return fmt.Errorf("failed to retrieve test cache entry: %v", err)
	}
	if cacheValue != testCacheValue {
		return fmt.Errorf("expected cache value '%s', got '%s'", testCacheValue, cacheValue)
	}

	return nil
}

// testDataDeletion tests deleting data from the tables
func testDataDeletion(db *sql.DB, dbType string) error {
	// Delete test module
	_, err := db.Exec(`
DELETE FROM modules
WHERE namespace = $1 AND name = $2 AND provider = $3 AND version = $4`,
		testModuleNamespace, testModuleName, testModuleProvider, testModuleVersion)
	if err != nil {
		return fmt.Errorf("failed to delete test module: %v", err)
	}

	// Delete test provider
	_, err = db.Exec(`
DELETE FROM providers
WHERE namespace = $1 AND name = $2 AND version = $3`,
		testProviderNamespace, testProviderName, testProviderVersion)
	if err != nil {
		return fmt.Errorf("failed to delete test provider: %v", err)
	}

	// Delete test cache entry
	_, err = db.Exec(`
DELETE FROM cache
WHERE cache_key = $1`,
		testCacheKey)
	if err != nil {
		return fmt.Errorf("failed to delete test cache entry: %v", err)
	}

	return nil
}

// loadEnvFile loads environment variables from .env file
func loadEnvFile() {
	envFile := ".env"
	
	// Check if .env file exists
	_, err := os.Stat(envFile)
	if os.IsNotExist(err) {
		fmt.Printf("No .env file found at %s. Using environment variables.\n", envFile)
		return
	}

	// Open .env file
	file, err := os.Open(envFile)
	if err != nil {
		fmt.Printf("Failed to open .env file: %v. Using environment variables.\n", err)
		return
	}
	defer file.Close()

	// Read .env file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		
		// Skip comments and empty lines
		if strings.HasPrefix(line, "#") || strings.TrimSpace(line) == "" {
			continue
		}

		// Parse key-value pairs
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		
		// Remove quotes if present
		value = strings.Trim(value, "\"'")

		// Set environment variable
		os.Setenv(key, value)
	}

	fmt.Printf("Loaded environment variables from %s\n", envFile)
}

// getEnv retrieves an environment variable or returns a default value if not set
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return strings.TrimSpace(value)
}
