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
	dbHost := os.Getenv("TOFU_REGISTRY_DB_HOST")
	dbPort := os.Getenv("TOFU_REGISTRY_DB_PORT")
	dbUser := os.Getenv("TOFU_REGISTRY_DB_USER")
	dbPassword := os.Getenv("TOFU_REGISTRY_DB_PASSWORD")
	dbSSLMode := os.Getenv("TOFU_REGISTRY_DB_SSLMODE")

	// Connect to the default database first
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=defaultdb sslmode=%s",
		dbHost, dbPort, dbUser, dbPassword, dbSSLMode)

	fmt.Println("Connecting to PostgreSQL database...")
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

	// Check if opentofu database already exists
	var exists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = 'opentofu')").Scan(&exists)
	if err != nil {
		fmt.Printf("Error checking if database exists: %v\n", err)
		return
	}

	if exists {
		fmt.Println("The 'opentofu' database already exists.")
	} else {
		// Create the opentofu database
		fmt.Println("Creating 'opentofu' database...")
		_, err = db.Exec("CREATE DATABASE opentofu")
		if err != nil {
			fmt.Printf("Error creating database: %v\n", err)
			return
		}
		fmt.Println("Successfully created 'opentofu' database!")
	}

	// Connect to the new opentofu database to test it
	connStrOpenTofu := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=opentofu sslmode=%s",
		dbHost, dbPort, dbUser, dbPassword, dbSSLMode)

	dbOpenTofu, err := sql.Open("postgres", connStrOpenTofu)
	if err != nil {
		fmt.Printf("Error opening connection to opentofu database: %v\n", err)
		return
	}
	defer dbOpenTofu.Close()

	// Test connection to opentofu database
	err = dbOpenTofu.Ping()
	if err != nil {
		fmt.Printf("Error connecting to opentofu database: %v\n", err)
		return
	}

	fmt.Println("Successfully connected to 'opentofu' database!")

	// Create necessary tables for the OpenTofu Registry
	fmt.Println("Creating registry tables...")
	
	// Create modules table
	_, err = dbOpenTofu.Exec(`
		CREATE TABLE IF NOT EXISTS modules (
			id SERIAL PRIMARY KEY,
			namespace VARCHAR(255) NOT NULL,
			name VARCHAR(255) NOT NULL,
			provider VARCHAR(255) NOT NULL,
			version VARCHAR(50) NOT NULL,
			download_url TEXT NOT NULL,
			published_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(namespace, name, provider, version)
		)
	`)
	if err != nil {
		fmt.Printf("Error creating modules table: %v\n", err)
		return
	}
	
	// Create providers table
	_, err = dbOpenTofu.Exec(`
		CREATE TABLE IF NOT EXISTS providers (
			id SERIAL PRIMARY KEY,
			namespace VARCHAR(255) NOT NULL,
			name VARCHAR(255) NOT NULL,
			version VARCHAR(50) NOT NULL,
			platforms JSONB NOT NULL,
			download_url TEXT NOT NULL,
			published_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(namespace, name, version)
		)
	`)
	if err != nil {
		fmt.Printf("Error creating providers table: %v\n", err)
		return
	}
	
	// Create cache table
	_, err = dbOpenTofu.Exec(`
		CREATE TABLE IF NOT EXISTS cache (
			id SERIAL PRIMARY KEY,
			key VARCHAR(255) NOT NULL UNIQUE,
			value JSONB NOT NULL,
			expires_at TIMESTAMP NOT NULL
		)
	`)
	if err != nil {
		fmt.Printf("Error creating cache table: %v\n", err)
		return
	}

	fmt.Println("Successfully created all registry tables!")
	fmt.Println("\nThe 'opentofu' database is now ready to use with the OpenTofu Registry API.")
	fmt.Println("Update your .env file to use this database instead of 'defaultdb'.")
}
