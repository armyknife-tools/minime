// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package registry

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-hclog"
	_ "github.com/lib/pq" // PostgreSQL driver
	_ "github.com/mattn/go-sqlite3" // SQLite driver
	svchost "github.com/hashicorp/terraform-svchost"
	"github.com/opentofu/opentofu/internal/registry/response"
)

const (
	// Default schema name
	defaultSchemaName = "registry"

	// Default connection pool settings
	defaultMaxOpenConns = 10
	defaultMaxIdleConns = 5
	defaultConnMaxLifetime = 1 * time.Hour
)

// DBClient provides database operations for registry data
type DBClient struct {
	db     *sql.DB
	logger hclog.Logger
	config *DBConfig
	isSQLite bool
}

// NewDBClient creates a new database client for registry operations
func NewDBClient(logger hclog.Logger) (*DBClient, error) {
	// Create and load database configuration
	config := NewDBConfig(logger)
	if err := config.LoadFromEnv(); err != nil {
		return nil, fmt.Errorf("failed to load database configuration: %w", err)
	}

	// Get database connection string
	dbURL := config.ConnectionString()
	isSQLite := config.IsSQLite()

	// Determine the driver based on the database type
	driver := "postgres"
	if isSQLite {
		driver = "sqlite3"
	}

	// Connect to the database
	db, err := sql.Open(driver, dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(defaultMaxOpenConns)
	db.SetMaxIdleConns(defaultMaxIdleConns)
	db.SetConnMaxLifetime(defaultConnMaxLifetime)

	// Test the connection
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	client := &DBClient{
		db:       db,
		logger:   logger,
		config:   config,
		isSQLite: isSQLite,
	}

	// Initialize the database schema
	if err := client.initSchema(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to initialize database schema: %w", err)
	}

	return client, nil
}

// Close closes the database connection
func (c *DBClient) Close() error {
	return c.db.Close()
}

// initSchema initializes the database schema
func (c *DBClient) initSchema() error {
	if c.isSQLite {
		return c.initSQLiteSchema()
	}
	return c.initPostgresSchema()
}

// initPostgresSchema initializes the PostgreSQL database schema
func (c *DBClient) initPostgresSchema() error {
	// Create the schema if it doesn't exist
	_, err := c.db.Exec(fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s", defaultSchemaName))
	if err != nil {
		return fmt.Errorf("failed to create schema: %w", err)
	}

	// Create the modules table
	_, err = c.db.Exec(fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s.modules (
			id SERIAL PRIMARY KEY,
			host TEXT NOT NULL,
			namespace TEXT NOT NULL,
			name TEXT NOT NULL,
			provider TEXT NOT NULL,
			version TEXT,
			downloads INTEGER NOT NULL DEFAULT 0,
			verified BOOLEAN NOT NULL DEFAULT FALSE,
			description TEXT,
			source TEXT,
			published_at TIMESTAMP WITH TIME ZONE,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(host, namespace, name, provider)
		)
	`, defaultSchemaName))
	if err != nil {
		return fmt.Errorf("failed to create modules table: %w", err)
	}

	// Create the providers table
	_, err = c.db.Exec(fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s.providers (
			id SERIAL PRIMARY KEY,
			host TEXT NOT NULL,
			namespace TEXT NOT NULL,
			name TEXT NOT NULL,
			downloads INTEGER NOT NULL DEFAULT 0,
			module_count INTEGER NOT NULL DEFAULT 0,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(host, namespace, name)
		)
	`, defaultSchemaName))
	if err != nil {
		return fmt.Errorf("failed to create providers table: %w", err)
	}

	// Create indexes for efficient searching
	_, err = c.db.Exec(fmt.Sprintf(`
		CREATE INDEX IF NOT EXISTS idx_modules_namespace_name_provider ON %s.modules(namespace, name, provider);
		CREATE INDEX IF NOT EXISTS idx_modules_downloads ON %s.modules(downloads DESC);
		CREATE INDEX IF NOT EXISTS idx_providers_namespace_name ON %s.providers(namespace, name);
		CREATE INDEX IF NOT EXISTS idx_providers_downloads ON %s.providers(downloads DESC);
	`, defaultSchemaName, defaultSchemaName, defaultSchemaName, defaultSchemaName))
	if err != nil {
		return fmt.Errorf("failed to create indexes: %w", err)
	}

	return nil
}

// initSQLiteSchema initializes the SQLite database schema
func (c *DBClient) initSQLiteSchema() error {
	// SQLite doesn't have schemas like PostgreSQL, so we'll use table prefixes instead
	// Create the modules table
	_, err := c.db.Exec(`
		CREATE TABLE IF NOT EXISTS registry_modules (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			host TEXT NOT NULL,
			namespace TEXT NOT NULL,
			name TEXT NOT NULL,
			provider TEXT NOT NULL,
			version TEXT,
			downloads INTEGER NOT NULL DEFAULT 0,
			verified BOOLEAN NOT NULL DEFAULT FALSE,
			description TEXT,
			source TEXT,
			published_at TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(host, namespace, name, provider)
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create modules table: %w", err)
	}

	// Create the providers table
	_, err = c.db.Exec(`
		CREATE TABLE IF NOT EXISTS registry_providers (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			host TEXT NOT NULL,
			namespace TEXT NOT NULL,
			name TEXT NOT NULL,
			downloads INTEGER NOT NULL DEFAULT 0,
			module_count INTEGER NOT NULL DEFAULT 0,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(host, namespace, name)
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create providers table: %w", err)
	}

	// Create indexes for efficient searching
	_, err = c.db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_modules_namespace_name_provider ON registry_modules(namespace, name, provider);
		CREATE INDEX IF NOT EXISTS idx_modules_downloads ON registry_modules(downloads DESC);
		CREATE INDEX IF NOT EXISTS idx_providers_namespace_name ON registry_providers(namespace, name);
		CREATE INDEX IF NOT EXISTS idx_providers_downloads ON registry_providers(downloads DESC);
	`)
	if err != nil {
		return fmt.Errorf("failed to create indexes: %w", err)
	}

	return nil
}

// getModulesTable returns the appropriate modules table name based on database type
func (c *DBClient) getModulesTable() string {
	if c.isSQLite {
		return "registry_modules"
	}
	return fmt.Sprintf("%s.modules", defaultSchemaName)
}

// getProvidersTable returns the appropriate providers table name based on database type
func (c *DBClient) getProvidersTable() string {
	if c.isSQLite {
		return "registry_providers"
	}
	return fmt.Sprintf("%s.providers", defaultSchemaName)
}

// SaveModules saves modules to the database
func (c *DBClient) SaveModules(ctx context.Context, host svchost.Hostname, modules []*response.Module) error {
	// Begin a transaction
	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Prepare the statement for inserting modules
	stmt, err := tx.PrepareContext(ctx, fmt.Sprintf(`
		INSERT INTO %s (
			host, namespace, name, provider, version, downloads, verified, description, source, published_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10
		) ON CONFLICT (host, namespace, name, provider) DO UPDATE SET
			version = EXCLUDED.version,
			downloads = EXCLUDED.downloads,
			verified = EXCLUDED.verified,
			description = EXCLUDED.description,
			source = EXCLUDED.source,
			published_at = EXCLUDED.published_at,
			updated_at = CURRENT_TIMESTAMP
	`, c.getModulesTable()))
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	// Insert modules
	hostStr := host.String()
	for _, module := range modules {
		// Extract version from the versions array if available
		var version string
		if len(module.Version) > 0 {
			version = module.Version
		}

		// Get the published_at time
		publishedAt := module.PublishedAt

		_, err = stmt.ExecContext(ctx,
			hostStr,
			module.Namespace,
			module.Name,
			module.Provider,
			version,
			module.Downloads,
			module.Verified,
			module.Description,
			module.Source,
			publishedAt,
		)
		if err != nil {
			return fmt.Errorf("failed to insert module %s/%s/%s: %w",
				module.Namespace, module.Name, module.Provider, err)
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	c.logger.Info("Saved modules to database", "host", hostStr, "count", len(modules))
	return nil
}

// SaveProviders saves providers to the database
func (c *DBClient) SaveProviders(ctx context.Context, host svchost.Hostname, providers []*response.ModuleProvider) error {
	// Begin a transaction
	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Prepare the statement for inserting providers
	stmt, err := tx.PrepareContext(ctx, fmt.Sprintf(`
		INSERT INTO %s (
			host, namespace, name, downloads, module_count
		) VALUES (
			$1, $2, $3, $4, $5
		) ON CONFLICT (host, namespace, name) DO UPDATE SET
			downloads = EXCLUDED.downloads,
			module_count = EXCLUDED.module_count,
			updated_at = CURRENT_TIMESTAMP
	`, c.getProvidersTable()))
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	// Insert providers
	hostStr := host.String()
	for _, provider := range providers {
		// Extract namespace and name from the provider name
		namespace, name, err := parseProviderID(provider.Name)
		if err != nil {
			c.logger.Warn("Failed to parse provider ID", "provider", provider.Name, "error", err)
			continue
		}

		_, err = stmt.ExecContext(ctx,
			hostStr,
			namespace,
			name,
			provider.Downloads,
			provider.ModuleCount,
		)
		if err != nil {
			return fmt.Errorf("failed to insert provider %s/%s: %w", namespace, name, err)
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	c.logger.Info("Saved providers to database", "host", hostStr, "count", len(providers))
	return nil
}

// GetModules retrieves modules from the database
func (c *DBClient) GetModules(ctx context.Context, host svchost.Hostname, query string, limit int) ([]*response.Module, error) {
	// Prepare the SQL query
	sqlQuery := fmt.Sprintf(`
		SELECT 
			namespace, name, provider, version, downloads, verified, description, source, published_at
		FROM 
			%s
		WHERE 
			host = $1
	`, c.getModulesTable())

	// Add search condition if query is provided
	args := []interface{}{host.String()}
	if query != "" {
		sqlQuery += ` AND (
			namespace ILIKE $2 OR
			name ILIKE $2 OR
			provider ILIKE $2 OR
			description ILIKE $2
		)`
		args = append(args, "%"+query+"%")
	}

	// Add order by and limit
	sqlQuery += ` ORDER BY downloads DESC`
	if limit > 0 {
		sqlQuery += fmt.Sprintf(` LIMIT %d`, limit)
	}

	// Execute the query
	rows, err := c.db.QueryContext(ctx, sqlQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query modules: %w", err)
	}
	defer rows.Close()

	// Process the results
	var modules []*response.Module
	for rows.Next() {
		var module response.Module
		var version sql.NullString
		var publishedAt sql.NullTime

		err := rows.Scan(
			&module.Namespace,
			&module.Name,
			&module.Provider,
			&version,
			&module.Downloads,
			&module.Verified,
			&module.Description,
			&module.Source,
			&publishedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan module row: %w", err)
		}

		// Set version if available
		if version.Valid {
			module.Version = version.String
		}

		// Set published_at if available
		if publishedAt.Valid {
			module.PublishedAt = publishedAt.Time
		}

		modules = append(modules, &module)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating module rows: %w", err)
	}

	return modules, nil
}

// GetProviders retrieves providers from the database
func (c *DBClient) GetProviders(ctx context.Context, host svchost.Hostname, query string, limit int) ([]*response.ModuleProvider, error) {
	// Prepare the SQL query
	sqlQuery := fmt.Sprintf(`
		SELECT 
			namespace, name, downloads, module_count
		FROM 
			%s
		WHERE 
			host = $1
	`, c.getProvidersTable())

	// Add search condition if query is provided
	args := []interface{}{host.String()}
	if query != "" {
		sqlQuery += ` AND (
			namespace ILIKE $2 OR
			name ILIKE $2
		)`
		args = append(args, "%"+query+"%")
	}

	// Add order by and limit
	sqlQuery += ` ORDER BY downloads DESC`
	if limit > 0 {
		sqlQuery += fmt.Sprintf(` LIMIT %d`, limit)
	}

	// Execute the query
	rows, err := c.db.QueryContext(ctx, sqlQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query providers: %w", err)
	}
	defer rows.Close()

	// Process the results
	var providers []*response.ModuleProvider
	for rows.Next() {
		var namespace, name string
		var downloads, moduleCount int

		err := rows.Scan(
			&namespace,
			&name,
			&downloads,
			&moduleCount,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan provider row: %w", err)
		}

		// Create the provider object
		provider := &response.ModuleProvider{
			Name:        fmt.Sprintf("%s/%s", namespace, name),
			Downloads:   downloads,
			ModuleCount: moduleCount,
		}

		providers = append(providers, provider)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating provider rows: %w", err)
	}

	return providers, nil
}

// parseProviderID parses a provider ID in the format "namespace/name"
func parseProviderID(id string) (namespace, name string, err error) {
	parts := strings.Split(id, "/")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid provider ID format: %s", id)
	}
	return parts[0], parts[1], nil
}
