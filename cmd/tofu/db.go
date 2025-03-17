// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mitchellh/cli"
	"github.com/opentofu/opentofu/internal/command"
	"github.com/opentofu/opentofu/internal/dotenv"
	"github.com/opentofu/opentofu/internal/templates"
	
	// Database drivers
	_ "github.com/lib/pq"           // PostgreSQL driver
	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

// DBCommand is a container for db subcommands
type DBCommand struct {
	Meta command.Meta
}

// DBSetupCommand is the command that handles database setup
type DBSetupCommand struct {
	Meta command.Meta
}

// DBConfigureCommand is the command that handles database configuration
type DBConfigureCommand struct {
	Meta command.Meta
}

// DBTestCommand is the command that handles database testing
type DBTestCommand struct {
	Meta command.Meta
}

// DBMigrateCommand is the command that handles database migrations
type DBMigrateCommand struct {
	Meta command.Meta
}

// Help returns help text for the DB command
func (c *DBCommand) Help() string {
	helpText := `
Usage: tofu db [subcommand]

  This command manages the OpenTofu database operations.

Subcommands:
  setup      Setup the database
  configure  Configure database connection parameters
  test       Test the database connection
  migrate    Migrate the database schema
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short description of the DB command
func (c *DBCommand) Synopsis() string {
	return "Manage OpenTofu database operations"
}

// Run runs the DB command
func (c *DBCommand) Run(args []string) int {
	return cli.RunResultHelp
}

// Help returns help text for the DB setup command
func (c *DBSetupCommand) Help() string {
	helpText := `
Usage: tofu db setup [options]

  This command sets up the OpenTofu database.

Options:
  -type=TYPE    Database type (sqlite or postgres). Default: sqlite
  -path=PATH    Path to SQLite database file. Default: ~/.opentofu/tofu.db
  -host=HOST    PostgreSQL host. Default: localhost
  -port=PORT    PostgreSQL port. Default: 5432
  -user=USER    PostgreSQL user. Default: postgres
  -password=PASSWORD  PostgreSQL password
  -dbname=NAME  PostgreSQL database name. Default: opentofu
  -sslmode=MODE PostgreSQL SSL mode. Default: disable
  -force        Force setup even if database already exists
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short description of the DB setup command
func (c *DBSetupCommand) Synopsis() string {
	return "Setup the OpenTofu database"
}

// Run runs the DB setup command
func (c *DBSetupCommand) Run(args []string) int {
	// Load environment variables from .env file
	vars, err := dotenv.Load("")
	if err != nil {
		c.Meta.Ui.Error(fmt.Sprintf("Error loading .env file: %s", err))
	} else if len(vars) > 0 {
		c.Meta.Ui.Info("Loaded environment variables from .env file")
	}

	var dbType string
	var dbPath string
	var host string
	var port string
	var user string
	var password string
	var dbname string
	var sslmode string
	var force bool

	// Parse command line arguments
	for i := 0; i < len(args); i++ {
		if args[i] == "-force" {
			force = true
		} else if strings.HasPrefix(args[i], "-type=") {
			dbType = args[i][6:]
		} else if strings.HasPrefix(args[i], "-path=") {
			dbPath = args[i][6:]
		} else if strings.HasPrefix(args[i], "-host=") {
			host = args[i][6:]
		} else if strings.HasPrefix(args[i], "-port=") {
			port = args[i][6:]
		} else if strings.HasPrefix(args[i], "-user=") {
			user = args[i][6:]
		} else if strings.HasPrefix(args[i], "-password=") {
			password = args[i][9:]
		} else if strings.HasPrefix(args[i], "-dbname=") {
			dbname = args[i][8:]
		} else if strings.HasPrefix(args[i], "-sslmode=") {
			sslmode = args[i][9:]
		}
	}

	// Set default values
	if dbType == "" {
		dbType = dotenv.GetWithDefault("TOFU_DB_TYPE", "sqlite")
	}
	if dbPath == "" {
		if dbType == "sqlite" {
			dbPath = dotenv.GetWithDefault("TOFU_DB_PATH", filepath.Join(os.Getenv("HOME"), ".opentofu", "tofu.db"))
		}
	}
	if host == "" {
		host = dotenv.GetWithDefault("TOFU_REGISTRY_DB_HOST", "localhost")
	}
	if port == "" {
		port = dotenv.GetWithDefault("TOFU_REGISTRY_DB_PORT", "5432")
	}
	if user == "" {
		user = dotenv.GetWithDefault("TOFU_REGISTRY_DB_USER", "postgres")
	}
	if password == "" {
		password = dotenv.GetWithDefault("TOFU_REGISTRY_DB_PASSWORD", "postgres")
	}
	if dbname == "" {
		dbname = dotenv.GetWithDefault("TOFU_REGISTRY_DB_NAME", "opentofu")
	}
	if sslmode == "" {
		sslmode = dotenv.GetWithDefault("TOFU_REGISTRY_DB_SSLMODE", "disable")
	}

	// Validate database type
	if dbType != "sqlite" && dbType != "postgres" {
		c.Meta.Ui.Error(fmt.Sprintf("Invalid database type: %s. Must be 'sqlite' or 'postgres'", dbType))
		return 1
	}

	// Display setup information
	c.Meta.Ui.Output("\n=== Database Setup ===")
	c.Meta.Ui.Info(fmt.Sprintf("Setting up %s database...", dbType))

	// Setup database based on type
	if dbType == "sqlite" {
		// Create directory if it doesn't exist
		dbDir := filepath.Dir(dbPath)
		if err := os.MkdirAll(dbDir, 0755); err != nil {
			c.Meta.Ui.Error(fmt.Sprintf("Error creating directory for SQLite database: %s", err))
			return 1
		}

		// Check if database file already exists
		if templates.FileExists(dbPath) && !force {
			c.Meta.Ui.Warn(fmt.Sprintf("SQLite database file already exists at %s", dbPath))
			c.Meta.Ui.Warn("Use -force to overwrite existing database")
			return 1
		}

		// Create SQLite database
		db, err := sql.Open("sqlite3", dbPath)
		if err != nil {
			c.Meta.Ui.Error(fmt.Sprintf("Error creating SQLite database: %s", err))
			return 1
		}
		defer db.Close()

		// Create tables
		c.Meta.Ui.Info("Creating database schema...")
		if err := createSQLiteSchema(db); err != nil {
			c.Meta.Ui.Error(fmt.Sprintf("Error creating database schema: %s", err))
			return 1
		}

		c.Meta.Ui.Info(fmt.Sprintf("SQLite database created at %s", dbPath))
	} else {
		// Setup PostgreSQL database
		c.Meta.Ui.Info("Connecting to PostgreSQL server...")
		
		// Connect to PostgreSQL server
		connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=postgres sslmode=%s", 
			host, port, user, password, sslmode)
		db, err := sql.Open("postgres", connStr)
		if err != nil {
			c.Meta.Ui.Error(fmt.Sprintf("Error connecting to PostgreSQL: %s", err))
			return 1
		}
		
		// Check if database exists
		var exists bool
		err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)", dbname).Scan(&exists)
		if err != nil {
			c.Meta.Ui.Error(fmt.Sprintf("Error checking if database exists: %s", err))
			db.Close()
			return 1
		}
		
		// Create database if it doesn't exist or force is specified
		if exists && force {
			c.Meta.Ui.Info(fmt.Sprintf("Dropping existing database '%s'...", dbname))
			_, err = db.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", dbname))
			if err != nil {
				c.Meta.Ui.Error(fmt.Sprintf("Error dropping database: %s", err))
				db.Close()
				return 1
			}
			exists = false
		}
		
		if !exists {
			c.Meta.Ui.Info(fmt.Sprintf("Creating database '%s'...", dbname))
			_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", dbname))
			if err != nil {
				c.Meta.Ui.Error(fmt.Sprintf("Error creating database: %s", err))
				db.Close()
				return 1
			}
		} else {
			c.Meta.Ui.Warn(fmt.Sprintf("Database '%s' already exists", dbname))
			c.Meta.Ui.Warn("Use -force to recreate the database")
			db.Close()
			return 1
		}
		
		// Close connection to postgres database
		db.Close()
		
		// Connect to the newly created database
		connStr = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", 
			host, port, user, password, dbname, sslmode)
		db, err = sql.Open("postgres", connStr)
		if err != nil {
			c.Meta.Ui.Error(fmt.Sprintf("Error connecting to new database: %s", err))
			return 1
		}
		defer db.Close()
		
		// Create tables
		c.Meta.Ui.Info("Creating database schema...")
		if err := createPostgresSchema(db); err != nil {
			c.Meta.Ui.Error(fmt.Sprintf("Error creating database schema: %s", err))
			return 1
		}
		
		c.Meta.Ui.Info(fmt.Sprintf("PostgreSQL database '%s' created successfully", dbname))
	}

	c.Meta.Ui.Output("\nDatabase setup complete!")
	return 0
}

// createSQLiteSchema creates the schema for SQLite database
func createSQLiteSchema(db *sql.DB) error {
	// Create templates table
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS templates (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			provider TEXT NOT NULL,
			resource TEXT NOT NULL,
			content TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return err
	}

	// Create unique index on provider and resource
	_, err = db.Exec(`
		CREATE UNIQUE INDEX IF NOT EXISTS idx_templates_provider_resource
		ON templates(provider, resource)
	`)
	if err != nil {
		return err
	}

	return nil
}

// createPostgresSchema creates the schema for PostgreSQL database
func createPostgresSchema(db *sql.DB) error {
	// Create templates table
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS templates (
			id SERIAL PRIMARY KEY,
			provider TEXT NOT NULL,
			resource TEXT NOT NULL,
			content TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return err
	}

	// Create unique index on provider and resource
	_, err = db.Exec(`
		CREATE UNIQUE INDEX IF NOT EXISTS idx_templates_provider_resource
		ON templates(provider, resource)
	`)
	if err != nil {
		return err
	}

	return nil
}

// Help returns help text for the DB configure command
func (c *DBConfigureCommand) Help() string {
	helpText := `
Usage: tofu db configure [options]

  This command configures the OpenTofu database connection parameters.

Options:
  -type=TYPE              Database type (sqlite or postgres). Default: sqlite
  -path=PATH              Path to SQLite database file. Default: ~/.opentofu/tofu.db
  -host=HOST              PostgreSQL host. Default: localhost
  -port=PORT              PostgreSQL port. Default: 5432
  -user=USER              PostgreSQL user. Default: postgres
  -password=PASSWORD      PostgreSQL password
  -dbname=NAME            PostgreSQL database name. Default: opentofu
  -sslmode=MODE           PostgreSQL SSL mode. Default: disable
  -save                   Save configuration to .env file
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short description of the DB configure command
func (c *DBConfigureCommand) Synopsis() string {
	return "Configure OpenTofu database connection parameters"
}

// Run runs the DB configure command
func (c *DBConfigureCommand) Run(args []string) int {
	// Load environment variables from .env file
	vars, err := dotenv.Load("")
	if err != nil {
		c.Meta.Ui.Error(fmt.Sprintf("Error loading .env file: %s", err))
	} else if len(vars) > 0 {
		c.Meta.Ui.Info("Loaded environment variables from .env file")
	}

	var dbType string
	var dbPath string
	var host string
	var port string
	var user string
	var password string
	var dbname string
	var sslmode string
	var save bool

	// Parse command line arguments
	for i := 0; i < len(args); i++ {
		if args[i] == "-save" {
			save = true
		} else if strings.HasPrefix(args[i], "-type=") {
			dbType = args[i][6:]
		} else if strings.HasPrefix(args[i], "-path=") {
			dbPath = args[i][6:]
		} else if strings.HasPrefix(args[i], "-host=") {
			host = args[i][6:]
		} else if strings.HasPrefix(args[i], "-port=") {
			port = args[i][6:]
		} else if strings.HasPrefix(args[i], "-user=") {
			user = args[i][6:]
		} else if strings.HasPrefix(args[i], "-password=") {
			password = args[i][9:]
		} else if strings.HasPrefix(args[i], "-dbname=") {
			dbname = args[i][8:]
		} else if strings.HasPrefix(args[i], "-sslmode=") {
			sslmode = args[i][9:]
		}
	}

	// Set default values
	if dbType == "" {
		dbType = dotenv.GetWithDefault("TOFU_DB_TYPE", "sqlite")
	}
	if dbPath == "" {
		if dbType == "sqlite" {
			dbPath = dotenv.GetWithDefault("TOFU_DB_PATH", filepath.Join(os.Getenv("HOME"), ".opentofu", "tofu.db"))
		}
	}
	if host == "" {
		host = dotenv.GetWithDefault("TOFU_REGISTRY_DB_HOST", "localhost")
	}
	if port == "" {
		port = dotenv.GetWithDefault("TOFU_REGISTRY_DB_PORT", "5432")
	}
	if user == "" {
		user = dotenv.GetWithDefault("TOFU_REGISTRY_DB_USER", "postgres")
	}
	if password == "" {
		password = dotenv.GetWithDefault("TOFU_REGISTRY_DB_PASSWORD", "postgres")
	}
	if dbname == "" {
		dbname = dotenv.GetWithDefault("TOFU_REGISTRY_DB_NAME", "opentofu")
	}
	if sslmode == "" {
		sslmode = dotenv.GetWithDefault("TOFU_REGISTRY_DB_SSLMODE", "disable")
	}

	// Validate database type
	if dbType != "sqlite" && dbType != "postgres" {
		c.Meta.Ui.Error(fmt.Sprintf("Invalid database type: %s. Must be 'sqlite' or 'postgres'", dbType))
		return 1
	}

	// Display configuration
	c.Meta.Ui.Info("Database Configuration:")
	c.Meta.Ui.Info(fmt.Sprintf("  Type: %s", dbType))
	if dbType == "sqlite" {
		c.Meta.Ui.Info(fmt.Sprintf("  Path: %s", dbPath))
	} else {
		c.Meta.Ui.Info(fmt.Sprintf("  Host: %s", host))
		c.Meta.Ui.Info(fmt.Sprintf("  Port: %s", port))
		c.Meta.Ui.Info(fmt.Sprintf("  User: %s", user))
		c.Meta.Ui.Info(fmt.Sprintf("  Password: %s", strings.Repeat("*", len(password))))
		c.Meta.Ui.Info(fmt.Sprintf("  Database: %s", dbname))
		c.Meta.Ui.Info(fmt.Sprintf("  SSL Mode: %s", sslmode))
	}

	// Test connection if requested
	if dbType == "sqlite" {
		// Create directory if it doesn't exist
		dbDir := filepath.Dir(dbPath)
		if err := os.MkdirAll(dbDir, 0755); err != nil {
			c.Meta.Ui.Error(fmt.Sprintf("Error creating directory for SQLite database: %s", err))
			return 1
		}
	} else {
		// Test PostgreSQL connection
		connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", 
			host, port, user, password, dbname, sslmode)
		db, err := sql.Open("postgres", connStr)
		if err != nil {
			c.Meta.Ui.Error(fmt.Sprintf("Error connecting to PostgreSQL: %s", err))
			c.Meta.Ui.Info("Configuration saved but connection test failed.")
		} else {
			err = db.Ping()
			if err != nil {
				c.Meta.Ui.Error(fmt.Sprintf("Error pinging PostgreSQL: %s", err))
				c.Meta.Ui.Info("Configuration saved but connection test failed.")
			} else {
				c.Meta.Ui.Info("Successfully connected to PostgreSQL database.")
				db.Close()
			}
		}
	}

	// Save configuration to .env file if requested
	if save {
		var envContent string
		if dbType == "sqlite" {
			envContent = fmt.Sprintf("TOFU_DB_TYPE=%s\nTOFU_DB_PATH=%s\n", dbType, dbPath)
		} else {
			envContent = fmt.Sprintf("TOFU_DB_TYPE=%s\nTOFU_REGISTRY_DB_HOST=%s\nTOFU_REGISTRY_DB_PORT=%s\nTOFU_REGISTRY_DB_USER=%s\nTOFU_REGISTRY_DB_PASSWORD=%s\nTOFU_REGISTRY_DB_NAME=%s\nTOFU_REGISTRY_DB_SSLMODE=%s\n",
				dbType, host, port, user, password, dbname, sslmode)
		}

		// Write to .env file
		envFile := ".env"
		if err := os.WriteFile(envFile, []byte(envContent), 0644); err != nil {
			c.Meta.Ui.Error(fmt.Sprintf("Error writing to .env file: %s", err))
			return 1
		}
		c.Meta.Ui.Info(fmt.Sprintf("Configuration saved to %s", envFile))
	}

	c.Meta.Ui.Info("Database configuration complete")
	return 0
}

// Help returns help text for the DB test command
func (c *DBTestCommand) Help() string {
	helpText := `
Usage: tofu db test [options]

  This command tests the OpenTofu database connection and functionality.

Options:
  -type=TYPE    Database type (sqlite or postgres). Default: sqlite
  -path=PATH    Path to SQLite database file. Default: ~/.opentofu/tofu.db
  -verbose      Show detailed test results
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short description of the DB test command
func (c *DBTestCommand) Synopsis() string {
	return "Test the OpenTofu database connection and functionality"
}

// Run runs the DB test command
func (c *DBTestCommand) Run(args []string) int {
	// Load environment variables from .env file
	vars, err := dotenv.Load("")
	if err != nil {
		c.Meta.Ui.Error(fmt.Sprintf("Error loading .env file: %s", err))
	} else if len(vars) > 0 {
		c.Meta.Ui.Info("Loaded environment variables from .env file")
	}

	var dbType string
	var dbPath string
	var verbose bool

	// Parse command line arguments
	for i := 0; i < len(args); i++ {
		if strings.HasPrefix(args[i], "-type=") {
			dbType = args[i][6:]
		} else if strings.HasPrefix(args[i], "-path=") {
			dbPath = args[i][6:]
		} else if args[i] == "-verbose" {
			verbose = true
		}
	}

	// Set default values
	if dbType == "" {
		dbType = dotenv.GetWithDefault("TOFU_DB_TYPE", "sqlite")
	}

	// Validate database type
	if dbType != "sqlite" && dbType != "postgres" {
		c.Meta.Ui.Error(fmt.Sprintf("Invalid database type: %s. Must be 'sqlite' or 'postgres'", dbType))
		return 1
	}

	// Handle SQLite-specific setup
	if dbType == "sqlite" && dbPath == "" {
		configDir := filepath.Join(os.Getenv("HOME"), ".opentofu")
		dbPath = filepath.Join(configDir, "tofu.db")
	}

	// Test database connection
	c.Meta.Ui.Info(fmt.Sprintf("Testing %s database connection...", dbType))
	startTime := time.Now()
	db, err := templates.ConnectToDatabase(dbType, dbPath)
	if err != nil {
		c.Meta.Ui.Error(fmt.Sprintf("Connection test failed: %s", err))
		return 1
	}
	defer db.Close()
	connectionTime := time.Since(startTime)
	c.Meta.Ui.Info(fmt.Sprintf("Connection successful (%.2f ms)", float64(connectionTime.Milliseconds())))

	// Test database ping
	c.Meta.Ui.Info("Testing database ping...")
	startTime = time.Now()
	err = db.Ping()
	if err != nil {
		c.Meta.Ui.Error(fmt.Sprintf("Ping test failed: %s", err))
		return 1
	}
	pingTime := time.Since(startTime)
	c.Meta.Ui.Info(fmt.Sprintf("Ping successful (%.2f ms)", float64(pingTime.Milliseconds())))

	// Test database schema
	c.Meta.Ui.Info("Testing database schema...")
	if err := testDatabaseSchema(db, verbose); err != nil {
		c.Meta.Ui.Error(fmt.Sprintf("Schema test failed: %s", err))
		return 1
	}
	c.Meta.Ui.Info("Schema test successful")

	// Test CRUD operations
	c.Meta.Ui.Info("Testing CRUD operations...")
	if err := testCRUDOperations(db, verbose); err != nil {
		c.Meta.Ui.Error(fmt.Sprintf("CRUD test failed: %s", err))
		return 1
	}
	c.Meta.Ui.Info("CRUD test successful")

	c.Meta.Ui.Info("All database tests passed successfully")
	return 0
}

// Help returns help text for the DB migrate command
func (c *DBMigrateCommand) Help() string {
	helpText := `
Usage: tofu db migrate [options]

  This command migrates the OpenTofu database schema.

Options:
  -type=TYPE    Database type (sqlite or postgres). Default: sqlite
  -path=PATH    Path to SQLite database file. Default: ~/.opentofu/tofu.db
  -backup       Create a backup before migration
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short description of the DB migrate command
func (c *DBMigrateCommand) Synopsis() string {
	return "Migrate the OpenTofu database schema"
}

// Run runs the DB migrate command
func (c *DBMigrateCommand) Run(args []string) int {
	// Load environment variables from .env file
	vars, err := dotenv.Load("")
	if err != nil {
		c.Meta.Ui.Error(fmt.Sprintf("Error loading .env file: %s", err))
	} else if len(vars) > 0 {
		c.Meta.Ui.Info("Loaded environment variables from .env file")
	}

	var dbType string
	var dbPath string
	var backup bool

	// Parse command line arguments
	for i := 0; i < len(args); i++ {
		if strings.HasPrefix(args[i], "-type=") {
			dbType = args[i][6:]
		} else if strings.HasPrefix(args[i], "-path=") {
			dbPath = args[i][6:]
		} else if args[i] == "-backup" {
			backup = true
		}
	}

	// Set default values
	if dbType == "" {
		dbType = dotenv.GetWithDefault("TOFU_DB_TYPE", "sqlite")
	}

	// Validate database type
	if dbType != "sqlite" && dbType != "postgres" {
		c.Meta.Ui.Error(fmt.Sprintf("Invalid database type: %s. Must be 'sqlite' or 'postgres'", dbType))
		return 1
	}

	// Handle SQLite-specific setup
	if dbType == "sqlite" && dbPath == "" {
		configDir := filepath.Join(os.Getenv("HOME"), ".opentofu")
		dbPath = filepath.Join(configDir, "tofu.db")
	}

	// Check if database exists
	if dbType == "sqlite" {
		if _, err := os.Stat(dbPath); os.IsNotExist(err) {
			c.Meta.Ui.Error(fmt.Sprintf("Database file does not exist at %s", dbPath))
			return 1
		}
	}

	// Create backup if requested
	if backup && dbType == "sqlite" {
		backupPath := dbPath + ".backup." + time.Now().Format("20060102150405")
		c.Meta.Ui.Info(fmt.Sprintf("Creating backup at %s", backupPath))
		
		data, err := os.ReadFile(dbPath)
		if err != nil {
			c.Meta.Ui.Error(fmt.Sprintf("Error reading database file: %s", err))
			return 1
		}
		
		if err := os.WriteFile(backupPath, data, 0644); err != nil {
			c.Meta.Ui.Error(fmt.Sprintf("Error writing backup file: %s", err))
			return 1
		}
	}

	// Connect to database
	c.Meta.Ui.Info(fmt.Sprintf("Connecting to %s database...", dbType))
	db, err := templates.ConnectToDatabase(dbType, dbPath)
	if err != nil {
		c.Meta.Ui.Error(fmt.Sprintf("Error connecting to database: %s", err))
		return 1
	}
	defer db.Close()

	// Perform migration
	c.Meta.Ui.Info("Migrating database schema...")
	if err := migrateDatabaseSchema(db); err != nil {
		c.Meta.Ui.Error(fmt.Sprintf("Error migrating database schema: %s", err))
		return 1
	}

	c.Meta.Ui.Info("Database migration completed successfully")
	return 0
}

// createDatabaseTables creates the necessary tables in the database
func createDatabaseTables(db *sql.DB) error {
	// Create templates table
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
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		UNIQUE(provider, resource)
	);
	`
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("error creating templates table: %v", err)
	}

	// Create schema_migrations table to track migrations
	query = `
	CREATE TABLE IF NOT EXISTS schema_migrations (
		version TEXT PRIMARY KEY,
		applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`
	_, err = db.Exec(query)
	if err != nil {
		return fmt.Errorf("error creating schema_migrations table: %v", err)
	}

	// Insert initial migration version
	query = `
	INSERT INTO schema_migrations (version) 
	VALUES ('20250317000001') 
	ON CONFLICT (version) DO NOTHING;
	`
	_, err = db.Exec(query)
	if err != nil {
		return fmt.Errorf("error inserting initial migration version: %v", err)
	}

	return nil
}

// testDatabaseSchema tests if the database schema is valid
func testDatabaseSchema(db *sql.DB, verbose bool) error {
	// Check if templates table exists
	var tableName string
	query := `
	SELECT name FROM sqlite_master WHERE type='table' AND name='templates';
	`
	err := db.QueryRow(query).Scan(&tableName)
	if err != nil {
		return fmt.Errorf("templates table not found: %v", err)
	}

	if verbose {
		// Get table schema
		rows, err := db.Query("PRAGMA table_info(templates);")
		if err != nil {
			return fmt.Errorf("error getting table schema: %v", err)
		}
		defer rows.Close()

		fmt.Println("Templates table schema:")
		for rows.Next() {
			var cid int
			var name string
			var dataType string
			var notNull int
			var dfltValue interface{}
			var pk int
			err = rows.Scan(&cid, &name, &dataType, &notNull, &dfltValue, &pk)
			if err != nil {
				return fmt.Errorf("error scanning row: %v", err)
			}
			fmt.Printf("  Column: %s, Type: %s, NotNull: %d, PK: %d\n", name, dataType, notNull, pk)
		}
	}

	return nil
}

// testCRUDOperations tests basic CRUD operations on the database
func testCRUDOperations(db *sql.DB, verbose bool) error {
	// Test data
	testTemplate := Template{
		Provider:    "test",
		Resource:    "test_resource",
		DisplayName: "Test Resource",
		Content:     "resource \"test_resource\" \"example\" {\n  name = \"test\"\n}",
		Description: "A test resource for database testing",
		Category:    "test",
		Tags:        "test,database",
	}

	// Create
	query := `
	INSERT INTO templates (provider, resource, display_name, content, description, category, tags)
	VALUES (?, ?, ?, ?, ?, ?, ?)
	ON CONFLICT (provider, resource) DO UPDATE
	SET display_name = ?, content = ?, description = ?, category = ?, tags = ?;
	`
	_, err := db.Exec(
		query,
		testTemplate.Provider,
		testTemplate.Resource,
		testTemplate.DisplayName,
		testTemplate.Content,
		testTemplate.Description,
		testTemplate.Category,
		testTemplate.Tags,
		testTemplate.DisplayName,
		testTemplate.Content,
		testTemplate.Description,
		testTemplate.Category,
		testTemplate.Tags,
	)
	if err != nil {
		return fmt.Errorf("error inserting test template: %v", err)
	}

	if verbose {
		fmt.Println("Created test template:")
		fmt.Printf("  Provider: %s\n", testTemplate.Provider)
		fmt.Printf("  Resource: %s\n", testTemplate.Resource)
	}

	// Read
	var template Template
	query = `
	SELECT id, provider, resource, display_name, content, description, category, tags
	FROM templates
	WHERE provider = ? AND resource = ?
	`
	err = db.QueryRow(
		query,
		testTemplate.Provider,
		testTemplate.Resource,
	).Scan(
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
		return fmt.Errorf("error reading test template: %v", err)
	}

	if verbose {
		fmt.Println("Read test template:")
		fmt.Printf("  ID: %d\n", template.ID)
		fmt.Printf("  Provider: %s\n", template.Provider)
		fmt.Printf("  Resource: %s\n", template.Resource)
		fmt.Printf("  DisplayName: %s\n", template.DisplayName)
	}

	// Update
	testTemplate.DisplayName = "Updated Test Resource"
	query = `
	UPDATE templates
	SET display_name = ?
	WHERE provider = ? AND resource = ?
	`
	_, err = db.Exec(
		query,
		testTemplate.DisplayName,
		testTemplate.Provider,
		testTemplate.Resource,
	)
	if err != nil {
		return fmt.Errorf("error updating test template: %v", err)
	}

	if verbose {
		fmt.Println("Updated test template:")
		fmt.Printf("  DisplayName: %s\n", testTemplate.DisplayName)
	}

	// Verify update
	err = db.QueryRow(
		`SELECT display_name FROM templates WHERE provider = ? AND resource = ?`,
		testTemplate.Provider,
		testTemplate.Resource,
	).Scan(&template.DisplayName)
	if err != nil {
		return fmt.Errorf("error verifying update: %v", err)
	}
	if template.DisplayName != testTemplate.DisplayName {
		return fmt.Errorf("update verification failed: expected %s, got %s", testTemplate.DisplayName, template.DisplayName)
	}

	// Delete
	query = `
	DELETE FROM templates
	WHERE provider = ? AND resource = ?
	`
	_, err = db.Exec(
		query,
		testTemplate.Provider,
		testTemplate.Resource,
	)
	if err != nil {
		return fmt.Errorf("error deleting test template: %v", err)
	}

	if verbose {
		fmt.Println("Deleted test template")
	}

	// Verify deletion
	var count int
	err = db.QueryRow(
		`SELECT COUNT(*) FROM templates WHERE provider = ? AND resource = ?`,
		testTemplate.Provider,
		testTemplate.Resource,
	).Scan(&count)
	if err != nil {
		return fmt.Errorf("error verifying deletion: %v", err)
	}
	if count != 0 {
		return fmt.Errorf("deletion verification failed: expected 0 records, found %d", count)
	}

	return nil
}

// migrateDatabaseSchema migrates the database schema to the latest version
func migrateDatabaseSchema(db *sql.DB) error {
	// Check current schema version
	var currentVersion string
	err := db.QueryRow("SELECT MAX(version) FROM schema_migrations").Scan(&currentVersion)
	if err != nil {
		return fmt.Errorf("error getting current schema version: %v", err)
	}

	// Define migrations
	migrations := map[string]string{
		"20250317000001": "", // Initial schema, already applied in createDatabaseTables
		"20250317000002": `
			ALTER TABLE templates ADD COLUMN IF NOT EXISTS version TEXT DEFAULT '1.0.0';
		`,
		// Add more migrations here as needed
	}

	// Apply migrations in order
	for version, migration := range migrations {
		if version <= currentVersion {
			continue
		}

		// Apply migration
		if migration != "" {
			_, err := db.Exec(migration)
			if err != nil {
				return fmt.Errorf("error applying migration %s: %v", version, err)
			}
		}

		// Update schema_migrations table
		_, err := db.Exec("INSERT INTO schema_migrations (version) VALUES (?)", version)
		if err != nil {
			return fmt.Errorf("error updating schema_migrations table: %v", err)
		}
	}

	return nil
}
