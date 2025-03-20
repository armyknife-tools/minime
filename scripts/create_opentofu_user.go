// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2023 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"bufio"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/lib/pq"
)

func main() {
	// Get database connection parameters from environment variables or use defaults
	host := getEnv("PGHOST", "localhost")
	port := getEnv("PGPORT", "5432")
	user := getEnv("PGUSER", "postgres")
	password := getEnv("PGPASSWORD", "")
	dbname := getEnv("PGDATABASE", "opentofu")
	newUser := getEnv("TOFU_DB_USER", "opentofu_user")

	// Generate a secure random password
	newPassword, err := generateSecurePassword(20)
	if err != nil {
		log.Fatalf("Failed to generate secure password: %v", err)
	}

	// Connect to PostgreSQL server
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", 
		host, port, user, password, dbname)
	
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

	// Check if user exists
	var exists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM pg_roles WHERE rolname=$1)", newUser).Scan(&exists)
	if err != nil {
		log.Fatalf("Failed to check if user exists: %v", err)
	}

	// Create user if it doesn't exist, otherwise alter user
	if !exists {
		fmt.Printf("Creating user '%s'...\n", newUser)
		_, err = db.Exec(fmt.Sprintf("CREATE USER %s WITH PASSWORD '%s'", newUser, newPassword))
		if err != nil {
			log.Fatalf("Failed to create user: %v", err)
		}
		fmt.Printf("User '%s' created successfully\n", newUser)
	} else {
		fmt.Printf("User '%s' already exists. Updating password...\n", newUser)
		_, err = db.Exec(fmt.Sprintf("ALTER USER %s WITH PASSWORD '%s'", newUser, newPassword))
		if err != nil {
			log.Fatalf("Failed to update user password: %v", err)
		}
		fmt.Printf("Password for user '%s' updated successfully\n", newUser)
	}

	// Grant privileges to the user
	fmt.Printf("Granting privileges to user '%s'...\n", newUser)
	
	// Grant privileges on database
	_, err = db.Exec(fmt.Sprintf("GRANT ALL PRIVILEGES ON DATABASE %s TO %s", dbname, newUser))
	if err != nil {
		log.Fatalf("Failed to grant database privileges: %v", err)
	}

	// Grant privileges on schema public
	_, err = db.Exec(fmt.Sprintf("GRANT ALL PRIVILEGES ON SCHEMA public TO %s", newUser))
	if err != nil {
		log.Fatalf("Failed to grant schema privileges: %v", err)
	}

	// Grant privileges on all tables in schema public
	_, err = db.Exec(fmt.Sprintf("GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO %s", newUser))
	if err != nil {
		log.Fatalf("Failed to grant table privileges: %v", err)
	}

	// Grant privileges on all sequences in schema public
	_, err = db.Exec(fmt.Sprintf("GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO %s", newUser))
	if err != nil {
		log.Fatalf("Failed to grant sequence privileges: %v", err)
	}

	fmt.Printf("Privileges granted to user '%s' successfully\n", newUser)

	// Update .env file with new credentials
	updateEnvFile(host, port, dbname, newUser, newPassword)

	fmt.Println("\nUser setup completed successfully!")
	fmt.Printf("Database user: %s\n", newUser)
	fmt.Printf("Database password: %s\n", newPassword)
	fmt.Println("These credentials have been saved to the .env file in the current directory.")
	fmt.Println("You can now test the connection with the test_opentofu_user_connection.go script.")
}

// generateSecurePassword generates a cryptographically secure random password
func generateSecurePassword(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes)[:length], nil
}

// updateEnvFile creates or updates the .env file with the database credentials
func updateEnvFile(host, port, dbname, user, password string) {
	envFile := ".env"
	
	// Check if .env file exists
	_, err := os.Stat(envFile)
	if os.IsNotExist(err) {
		// Create new .env file
		file, err := os.Create(envFile)
		if err != nil {
			log.Printf("Failed to create .env file: %v", err)
			return
		}
		defer file.Close()

		// Write database credentials to .env file
		writer := bufio.NewWriter(file)
		fmt.Fprintln(writer, "# OpenTofu Registry API Database Configuration")
		fmt.Fprintln(writer, "TOFU_REGISTRY_DB_TYPE=postgres")
		fmt.Fprintf(writer, "TOFU_REGISTRY_DB_HOST=%s\n", host)
		fmt.Fprintf(writer, "TOFU_REGISTRY_DB_PORT=%s\n", port)
		fmt.Fprintf(writer, "TOFU_REGISTRY_DB_NAME=%s\n", dbname)
		fmt.Fprintf(writer, "TOFU_REGISTRY_DB_USER=%s\n", user)
		fmt.Fprintf(writer, "TOFU_REGISTRY_DB_PASSWORD=%s\n", password)
		fmt.Fprintln(writer, "TOFU_REGISTRY_DB_SSLMODE=require")
		
		writer.Flush()
		fmt.Printf("Created .env file with database credentials at %s\n", filepath.Abs(envFile))
	} else {
		// Read existing .env file
		file, err := os.Open(envFile)
		if err != nil {
			log.Printf("Failed to open .env file: %v", err)
			return
		}
		
		lines := []string{}
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			if !strings.HasPrefix(line, "TOFU_REGISTRY_DB_") {
				lines = append(lines, line)
			}
		}
		file.Close()

		// Add database credentials
		lines = append(lines, "# OpenTofu Registry API Database Configuration")
		lines = append(lines, "TOFU_REGISTRY_DB_TYPE=postgres")
		lines = append(lines, fmt.Sprintf("TOFU_REGISTRY_DB_HOST=%s", host))
		lines = append(lines, fmt.Sprintf("TOFU_REGISTRY_DB_PORT=%s", port))
		lines = append(lines, fmt.Sprintf("TOFU_REGISTRY_DB_NAME=%s", dbname))
		lines = append(lines, fmt.Sprintf("TOFU_REGISTRY_DB_USER=%s", user))
		lines = append(lines, fmt.Sprintf("TOFU_REGISTRY_DB_PASSWORD=%s", password))
		lines = append(lines, "TOFU_REGISTRY_DB_SSLMODE=require")

		// Write updated .env file
		file, err = os.Create(envFile)
		if err != nil {
			log.Printf("Failed to update .env file: %v", err)
			return
		}
		defer file.Close()

		writer := bufio.NewWriter(file)
		for _, line := range lines {
			fmt.Fprintln(writer, line)
		}
		writer.Flush()
		fmt.Printf("Updated .env file with database credentials at %s\n", envFile)
	}
}

// getEnv retrieves an environment variable or returns a default value if not set
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return strings.TrimSpace(value)
}
