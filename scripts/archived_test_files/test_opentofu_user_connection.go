package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("Error loading .env file: %v\n", err)
		return
	}

	// Get database connection parameters from environment variables
	dbType := os.Getenv("TOFU_REGISTRY_DB_TYPE")
	dbHost := os.Getenv("TOFU_REGISTRY_DB_HOST")
	dbPort := os.Getenv("TOFU_REGISTRY_DB_PORT")
	dbName := os.Getenv("TOFU_REGISTRY_DB_NAME")
	dbUser := os.Getenv("TOFU_REGISTRY_DB_USER")
	dbPassword := os.Getenv("TOFU_REGISTRY_DB_PASSWORD")
	dbSSLMode := os.Getenv("TOFU_REGISTRY_DB_SSLMODE")

	// Print connection parameters (without password)
	fmt.Printf("Database Type: %s\n", dbType)
	fmt.Printf("Host: %s\n", dbHost)
	fmt.Printf("Port: %s\n", dbPort)
	fmt.Printf("Database: %s\n", dbName)
	fmt.Printf("User: %s\n", dbUser)
	fmt.Printf("SSL Mode: %s\n", dbSSLMode)
	fmt.Println()

	// Connect to PostgreSQL database
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbHost, dbPort, dbUser, dbPassword, dbName, dbSSLMode)

	fmt.Println("Connecting to PostgreSQL database with opentofu_user...")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Printf("Error opening database connection: %v\n", err)
		return
	}
	defer db.Close()

	// Test connection
	err = db.Ping()
	if err != nil {
		fmt.Printf("Error connecting to database: %v\n", err)
		return
	}

	fmt.Println("Successfully connected to PostgreSQL database!")

	// Get PostgreSQL version
	var version string
	err = db.QueryRow("SELECT version()").Scan(&version)
	if err != nil {
		fmt.Printf("Error getting PostgreSQL version: %v\n", err)
		return
	}
	fmt.Printf("PostgreSQL version: %s\n\n", version)

	// Test table creation
	fmt.Println("Testing table creation...")
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS opentofu_test (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		fmt.Printf("Error creating test table: %v\n", err)
		return
	}
	fmt.Println("Successfully created test table!")

	// Test data insertion
	fmt.Println("\nTesting data insertion...")
	_, err = db.Exec(`
		INSERT INTO opentofu_test (name) 
		VALUES ('OpenTofu Registry Test with opentofu_user')
	`)
	if err != nil {
		fmt.Printf("Error inserting test data: %v\n", err)
		return
	}
	fmt.Println("Successfully inserted test data!")

	// Test data retrieval
	fmt.Println("\nTesting data retrieval...")
	var id int
	var name string
	var createdAt string
	err = db.QueryRow(`
		SELECT id, name, created_at 
		FROM opentofu_test 
		ORDER BY id DESC 
		LIMIT 1
	`).Scan(&id, &name, &createdAt)
	if err != nil {
		fmt.Printf("Error retrieving test data: %v\n", err)
		return
	}
	fmt.Printf("Retrieved row: ID=%d, Name=%s, CreatedAt=%s\n\n", id, name, createdAt)

	// Test creating a table for modules
	fmt.Println("Testing module table creation...")
	_, err = db.Exec(`
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
		)
	`)
	if err != nil {
		fmt.Printf("Error creating modules table: %v\n", err)
		return
	}
	fmt.Println("Successfully created modules table!")

	// Test creating a table for providers
	fmt.Println("\nTesting provider table creation...")
	_, err = db.Exec(`
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
		)
	`)
	if err != nil {
		fmt.Printf("Error creating providers table: %v\n", err)
		return
	}
	fmt.Println("Successfully created providers table!")

	// Test creating a cache table
	fmt.Println("\nTesting cache table creation...")
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS cache (
			id SERIAL PRIMARY KEY,
			cache_key VARCHAR(255) NOT NULL UNIQUE,
			cache_value TEXT NOT NULL,
			expires_at TIMESTAMP NOT NULL
		)
	`)
	if err != nil {
		fmt.Printf("Error creating cache table: %v\n", err)
		return
	}
	fmt.Println("Successfully created cache table!")

	fmt.Println("\nAll PostgreSQL connection tests passed successfully!")
	fmt.Println("Your database configuration with opentofu_user is correct and ready for the OpenTofu Registry API.")
}
