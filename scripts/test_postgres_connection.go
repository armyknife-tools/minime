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

	// Print connection details (without password)
	fmt.Printf("Database Type: %s\n", dbType)
	fmt.Printf("Host: %s\n", dbHost)
	fmt.Printf("Port: %s\n", dbPort)
	fmt.Printf("Database: %s\n", dbName)
	fmt.Printf("User: %s\n", dbUser)
	fmt.Printf("SSL Mode: %s\n", dbSSLMode)

	// Construct connection string
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbHost, dbPort, dbUser, dbPassword, dbName, dbSSLMode)

	// Connect to database
	fmt.Println("\nConnecting to PostgreSQL database...")
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

	// Get database version
	var version string
	err = db.QueryRow("SELECT version()").Scan(&version)
	if err != nil {
		fmt.Printf("Error querying database version: %v\n", err)
		return
	}

	fmt.Printf("PostgreSQL version: %s\n", version)

	// Test creating a table
	fmt.Println("\nTesting table creation...")
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS registry_test (
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

	// Test inserting data
	fmt.Println("\nTesting data insertion...")
	_, err = db.Exec(`
		INSERT INTO registry_test (name) VALUES ('OpenTofu Registry Test')
	`)
	if err != nil {
		fmt.Printf("Error inserting test data: %v\n", err)
		return
	}

	fmt.Println("Successfully inserted test data!")

	// Test querying data
	fmt.Println("\nTesting data retrieval...")
	rows, err := db.Query(`
		SELECT id, name, created_at FROM registry_test ORDER BY id DESC LIMIT 1
	`)
	if err != nil {
		fmt.Printf("Error querying test data: %v\n", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var name string
		var createdAt string
		err = rows.Scan(&id, &name, &createdAt)
		if err != nil {
			fmt.Printf("Error scanning row: %v\n", err)
			return
		}
		fmt.Printf("Retrieved row: ID=%d, Name=%s, CreatedAt=%s\n", id, name, createdAt)
	}

	fmt.Println("\nAll PostgreSQL connection tests passed successfully!")
	fmt.Println("Your database configuration is correct and ready for the OpenTofu Registry API.")
}
