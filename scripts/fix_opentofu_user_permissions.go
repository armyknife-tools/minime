package main

import (
	"database/sql"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

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
	dbSSLMode := os.Getenv("TOFU_REGISTRY_DB_SSLMODE")

	// Connect to the database with admin credentials
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=opentofu sslmode=%s",
		dbHost, dbPort, "vultradmin", "AVNS_bwffjRd1Lg9SFhCqN83", dbSSLMode)

	fmt.Println("Connecting to PostgreSQL database with admin credentials...")
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

	// Generate a secure random password
	password := generateSecurePassword(16)

	// Fix permissions for opentofu_user
	fmt.Println("Fixing permissions for 'opentofu_user'...")
	
	// Grant schema ownership to the admin user first
	_, err = db.Exec("ALTER SCHEMA public OWNER TO vultradmin")
	if err != nil {
		fmt.Printf("Error changing schema owner: %v\n", err)
		return
	}
	
	// Grant create permission on schema
	_, err = db.Exec("GRANT CREATE ON SCHEMA public TO opentofu_user")
	if err != nil {
		fmt.Printf("Error granting create permission: %v\n", err)
		return
	}
	
	// Grant all privileges on all tables
	_, err = db.Exec("GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO opentofu_user")
	if err != nil {
		fmt.Printf("Error granting table privileges: %v\n", err)
		return
	}
	
	// Grant all privileges on all sequences
	_, err = db.Exec("GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO opentofu_user")
	if err != nil {
		fmt.Printf("Error granting sequence privileges: %v\n", err)
		return
	}
	
	// Update user password
	_, err = db.Exec(fmt.Sprintf("ALTER USER opentofu_user WITH PASSWORD '%s'", password))
	if err != nil {
		fmt.Printf("Error updating user password: %v\n", err)
		return
	}
	
	fmt.Println("Successfully fixed permissions for 'opentofu_user'!")

	// Test connection with updated user
	connStrOpenTofuUser := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=opentofu sslmode=%s",
		dbHost, dbPort, "opentofu_user", password, dbSSLMode)

	dbOpenTofuUser, err := sql.Open("postgres", connStrOpenTofuUser)
	if err != nil {
		fmt.Printf("Error opening connection with opentofu_user: %v\n", err)
		return
	}
	defer dbOpenTofuUser.Close()

	// Test connection with opentofu_user
	err = dbOpenTofuUser.Ping()
	if err != nil {
		fmt.Printf("Error connecting with opentofu_user: %v\n", err)
		return
	}

	fmt.Println("Successfully connected with 'opentofu_user'!")

	// Test permissions by creating a test table
	_, err = dbOpenTofuUser.Exec(`
		CREATE TABLE IF NOT EXISTS user_test (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		fmt.Printf("Error creating test table with opentofu_user: %v\n", err)
		return
	}

	fmt.Println("Successfully created test table with 'opentofu_user'!")

	// Update .env file with new credentials
	fmt.Println("\nUpdating .env file with new credentials...")
	
	envContent, err := os.ReadFile(".env")
	if err != nil {
		fmt.Printf("Error reading .env file: %v\n", err)
		return
	}
	
	lines := strings.Split(string(envContent), "\n")
	newLines := []string{}
	
	for _, line := range lines {
		if strings.HasPrefix(line, "TOFU_REGISTRY_DB_USER=") {
			newLines = append(newLines, "TOFU_REGISTRY_DB_USER=opentofu_user")
		} else if strings.HasPrefix(line, "TOFU_REGISTRY_DB_PASSWORD=") {
			newLines = append(newLines, fmt.Sprintf("TOFU_REGISTRY_DB_PASSWORD=%s", password))
		} else {
			newLines = append(newLines, line)
		}
	}
	
	err = os.WriteFile(".env", []byte(strings.Join(newLines, "\n")), 0644)
	if err != nil {
		fmt.Printf("Error writing to .env file: %v\n", err)
		return
	}
	
	fmt.Println("Successfully updated .env file with new credentials!")
	fmt.Println("\nThe 'opentofu_user' is now set up and ready to use with the OpenTofu Registry API.")
	fmt.Printf("Username: opentofu_user\n")
	fmt.Printf("Password: %s\n", password)
	fmt.Println("\nThese credentials have been saved to your .env file.")
}

// Generate a secure random password
func generateSecurePassword(length int) string {
	rand.Seed(time.Now().UnixNano())
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()-_=+[]{}|;:,.<>?"
	password := make([]byte, length)
	for i := 0; i < length; i++ {
		password[i] = chars[rand.Intn(len(chars))]
	}
	return string(password)
}
