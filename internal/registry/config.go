// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package registry

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/go-hclog"
	"github.com/joho/godotenv"
)

const (
	// Default environment file names
	defaultEnvFile     = ".env"
	registryEnvFile    = ".env.registry"
	
	// Environment variable names
	dbURLEnvVar        = "TOFU_REGISTRY_DB_URL"
	dbTypeEnvVar       = "TOFU_REGISTRY_DB_TYPE"
	dbHostEnvVar       = "TOFU_REGISTRY_DB_HOST"
	dbPortEnvVar       = "TOFU_REGISTRY_DB_PORT"
	dbNameEnvVar       = "TOFU_REGISTRY_DB_NAME"
	dbUserEnvVar       = "TOFU_REGISTRY_DB_USER"
	dbPasswordEnvVar   = "TOFU_REGISTRY_DB_PASSWORD"
	dbSSLModeEnvVar    = "TOFU_REGISTRY_DB_SSLMODE"
	
	// Default database values
	defaultDBType      = "sqlite"
	defaultDBHost      = "localhost"
	defaultDBPort      = "5432"
	defaultDBName      = "opentofu_registry"
	defaultDBUser      = "opentofu"
	defaultDBPassword  = ""
	defaultDBSSLMode   = "disable"
	
	// SQLite default path
	defaultSQLitePath  = "registry.db"
)

// DBConfig holds the database configuration
type DBConfig struct {
	// Database connection URL (takes precedence if set)
	URL string
	
	// Individual connection parameters (used if URL is not set)
	Type     string // "postgres" or "sqlite"
	Host     string
	Port     string
	Name     string
	User     string
	Password string
	SSLMode  string
	
	// SQLite specific
	SQLitePath string
	
	// Logger
	Logger hclog.Logger
}

// NewDBConfig creates a new database configuration
func NewDBConfig(logger hclog.Logger) *DBConfig {
	return &DBConfig{
		Type:       defaultDBType,
		Host:       defaultDBHost,
		Port:       defaultDBPort,
		Name:       defaultDBName,
		User:       defaultDBUser,
		Password:   defaultDBPassword,
		SSLMode:    defaultDBSSLMode,
		SQLitePath: defaultSQLitePath,
		Logger:     logger,
	}
}

// LoadFromEnv loads database configuration from environment variables
func (c *DBConfig) LoadFromEnv() error {
	// Load from .env files if they exist
	c.loadEnvFiles()
	
	// Check for direct connection URL
	if url := os.Getenv(dbURLEnvVar); url != "" {
		c.URL = url
		c.Logger.Debug("Loaded database URL from environment variable", "var", dbURLEnvVar)
		return nil
	}
	
	// Load individual connection parameters
	if dbType := os.Getenv(dbTypeEnvVar); dbType != "" {
		c.Type = dbType
	}
	
	if dbHost := os.Getenv(dbHostEnvVar); dbHost != "" {
		c.Host = dbHost
	}
	
	if dbPort := os.Getenv(dbPortEnvVar); dbPort != "" {
		c.Port = dbPort
	}
	
	if dbName := os.Getenv(dbNameEnvVar); dbName != "" {
		c.Name = dbName
	}
	
	if dbUser := os.Getenv(dbUserEnvVar); dbUser != "" {
		c.User = dbUser
	}
	
	if dbPassword := os.Getenv(dbPasswordEnvVar); dbPassword != "" {
		c.Password = dbPassword
	}
	
	if dbSSLMode := os.Getenv(dbSSLModeEnvVar); dbSSLMode != "" {
		c.SSLMode = dbSSLMode
	}
	
	return nil
}

// loadEnvFiles attempts to load environment variables from .env files
func (c *DBConfig) loadEnvFiles() {
	// Try to load from common .env file locations
	envFiles := []string{
		registryEnvFile,  // .env.registry in current directory
		defaultEnvFile,   // .env in current directory
		filepath.Join(os.Getenv("HOME"), registryEnvFile), // ~/.env.registry
		filepath.Join(os.Getenv("HOME"), defaultEnvFile),  // ~/.env
	}
	
	for _, file := range envFiles {
		if _, err := os.Stat(file); err == nil {
			if err := godotenv.Load(file); err == nil {
				c.Logger.Debug("Loaded environment variables from file", "file", file)
				return
			}
		}
	}
}

// ConnectionString returns the database connection string
func (c *DBConfig) ConnectionString() string {
	// If URL is set, use it directly
	if c.URL != "" {
		return c.URL
	}
	
	// Otherwise, build connection string based on type
	switch strings.ToLower(c.Type) {
	case "postgres":
		return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
			c.User, c.Password, c.Host, c.Port, c.Name, c.SSLMode)
	case "sqlite":
		return fmt.Sprintf("file:%s?cache=shared&mode=rwc", c.SQLitePath)
	default:
		// Default to SQLite
		return fmt.Sprintf("file:%s?cache=shared&mode=rwc", c.SQLitePath)
	}
}

// IsSQLite returns true if the database type is SQLite
func (c *DBConfig) IsSQLite() bool {
	return strings.ToLower(c.Type) == "sqlite"
}

// IsPostgres returns true if the database type is PostgreSQL
func (c *DBConfig) IsPostgres() bool {
	return strings.ToLower(c.Type) == "postgres"
}
