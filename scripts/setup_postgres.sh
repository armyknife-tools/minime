#!/bin/bash
# Setup Database for OpenTofu Registry API
# This script helps users set up either a PostgreSQL or SQLite database for the OpenTofu Registry API

set -e

# Default values
DB_TYPE="sqlite"  # Default to SQLite as it's simpler to set up
DB_NAME="opentofu_registry"
DB_USER="opentofu"
DB_PASSWORD="opentofu"
DB_HOST="localhost"
DB_PORT="5432"
DB_SCHEMA="registry"
SQLITE_PATH="registry.db"

# Text formatting
BOLD="\033[1m"
GREEN="\033[0;32m"
YELLOW="\033[0;33m"
RED="\033[0;31m"
RESET="\033[0m"

# Print header
echo -e "${BOLD}OpenTofu Registry Database Setup${RESET}"
echo "This script will help you set up a database for the OpenTofu Registry API."
echo

# Ask for database type
echo -e "${BOLD}Select database type:${RESET}"
echo "1) SQLite (simpler, no additional setup required)"
echo "2) PostgreSQL (more robust, requires PostgreSQL server)"
read -p "Enter your choice [1]: " DB_TYPE_CHOICE
echo

case $DB_TYPE_CHOICE in
    2)
        DB_TYPE="postgres"
        ;;
    *)
        DB_TYPE="sqlite"
        ;;
esac

# SQLite setup
if [ "$DB_TYPE" = "sqlite" ]; then
    echo -e "${BOLD}Setting up SQLite database${RESET}"
    
    # Get SQLite database path
    read -p "SQLite database path [$SQLITE_PATH]: " input
    SQLITE_PATH=${input:-$SQLITE_PATH}
    
    # Create directory for SQLite database if needed
    mkdir -p "$(dirname "$SQLITE_PATH")"
    
    # Create .env file
    ENV_FILE="$(pwd)/.env.registry"
    echo "TOFU_REGISTRY_DB_TYPE=sqlite" > $ENV_FILE
    echo "TOFU_REGISTRY_DB_URL=\"file:$SQLITE_PATH?cache=shared&mode=rwc\"" >> $ENV_FILE
    
    echo -e "${GREEN}SQLite database configuration complete!${RESET}"
    echo "Database will be created at: $SQLITE_PATH"
    echo
    echo "To use the database with OpenTofu, run:"
    echo "export TOFU_REGISTRY_DB_TYPE=sqlite"
    echo "export TOFU_REGISTRY_DB_URL=\"file:$SQLITE_PATH?cache=shared&mode=rwc\""
    echo "or source the created .env file:"
    echo "source $ENV_FILE"
    exit 0
fi

# PostgreSQL setup
echo -e "${BOLD}Setting up PostgreSQL database${RESET}"

# Check if PostgreSQL is installed
if ! command -v psql &> /dev/null; then
    echo -e "${RED}PostgreSQL is not installed.${RESET}"
    echo "Please install PostgreSQL before running this script."
    echo "On Ubuntu/Debian: sudo apt-get install postgresql postgresql-contrib"
    echo "On RHEL/CentOS: sudo yum install postgresql-server postgresql-contrib"
    echo "On macOS with Homebrew: brew install postgresql"
    exit 1
fi

echo -e "${GREEN}PostgreSQL is installed.${RESET}"

# Get PostgreSQL connection details
read -p "Database name [$DB_NAME]: " input
DB_NAME=${input:-$DB_NAME}

read -p "Database user [$DB_USER]: " input
DB_USER=${input:-$DB_USER}

read -p "Database password [$DB_PASSWORD]: " input
DB_PASSWORD=${input:-$DB_PASSWORD}

read -p "Database host [$DB_HOST]: " input
DB_HOST=${input:-$DB_HOST}

read -p "Database port [$DB_PORT]: " input
DB_PORT=${input:-$DB_PORT}

read -p "Database schema [$DB_SCHEMA]: " input
DB_SCHEMA=${input:-$DB_SCHEMA}

# Confirm settings
echo
echo -e "${BOLD}Database Settings:${RESET}"
echo "Database type: PostgreSQL"
echo "Database name: $DB_NAME"
echo "Database user: $DB_USER"
echo "Database password: $DB_PASSWORD"
echo "Database host: $DB_HOST"
echo "Database port: $DB_PORT"
echo "Database schema: $DB_SCHEMA"
echo

read -p "Are these settings correct? (y/n) " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo -e "${YELLOW}Setup cancelled.${RESET}"
    exit 1
fi

# Create database and user
echo -e "${BOLD}Creating database and user...${RESET}"

# Check if the user can connect to PostgreSQL as postgres user
if ! psql -h $DB_HOST -p $DB_PORT -U postgres -c '\q' 2>/dev/null; then
    echo -e "${YELLOW}Cannot connect to PostgreSQL as postgres user.${RESET}"
    echo "You may need to run this script with sudo or as a user with PostgreSQL admin privileges."
    
    # Try to connect as the current user
    if ! psql -h $DB_HOST -p $DB_PORT -c '\q' 2>/dev/null; then
        echo -e "${RED}Cannot connect to PostgreSQL. Please check your PostgreSQL installation.${RESET}"
        exit 1
    else
        # Connected as current user, use that
        PSQL_CMD="psql -h $DB_HOST -p $DB_PORT"
    fi
else
    # Connected as postgres user
    PSQL_CMD="psql -h $DB_HOST -p $DB_PORT -U postgres"
fi

# Create user if it doesn't exist
$PSQL_CMD -c "SELECT 1 FROM pg_roles WHERE rolname='$DB_USER'" | grep -q 1 || \
    $PSQL_CMD -c "CREATE USER $DB_USER WITH PASSWORD '$DB_PASSWORD';"

# Create database if it doesn't exist
$PSQL_CMD -c "SELECT 1 FROM pg_database WHERE datname='$DB_NAME'" | grep -q 1 || \
    $PSQL_CMD -c "CREATE DATABASE $DB_NAME OWNER $DB_USER;"

# Grant privileges
$PSQL_CMD -c "GRANT ALL PRIVILEGES ON DATABASE $DB_NAME TO $DB_USER;"

# Connect to the database and create schema
PSQL_DB_CMD="psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME"

# Check if we need a password
if [ "$DB_PASSWORD" != "" ]; then
    export PGPASSWORD="$DB_PASSWORD"
fi

# Create schema if it doesn't exist
$PSQL_DB_CMD -c "CREATE SCHEMA IF NOT EXISTS $DB_SCHEMA;"
$PSQL_DB_CMD -c "GRANT ALL ON SCHEMA $DB_SCHEMA TO $DB_USER;"

echo -e "${GREEN}Database and user created successfully.${RESET}"

# Generate connection string
DB_URL="postgres://$DB_USER:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=disable"

# Create environment variable setup
echo
echo -e "${BOLD}Environment Variable Setup${RESET}"
echo "Add the following lines to your ~/.bashrc or ~/.zshrc file:"
echo
echo "export TOFU_REGISTRY_DB_TYPE=\"postgres\""
echo "export TOFU_REGISTRY_DB_URL=\"$DB_URL\""
echo

# Create a local .env file
echo -e "${BOLD}Creating .env file...${RESET}"
ENV_FILE="$(pwd)/.env.registry"
echo "TOFU_REGISTRY_DB_TYPE=\"postgres\"" > $ENV_FILE
echo "TOFU_REGISTRY_DB_URL=\"$DB_URL\"" >> $ENV_FILE
echo -e "${GREEN}Created $ENV_FILE${RESET}"
echo "You can source this file before running OpenTofu:"
echo "source $ENV_FILE"

# Test connection
echo
echo -e "${BOLD}Testing database connection...${RESET}"
if $PSQL_DB_CMD -c "SELECT 1;" > /dev/null 2>&1; then
    echo -e "${GREEN}Connection successful!${RESET}"
else
    echo -e "${RED}Connection failed. Please check your settings.${RESET}"
    exit 1
fi

# Create database schema
echo
echo -e "${BOLD}Creating database schema...${RESET}"
$PSQL_DB_CMD << EOF
CREATE TABLE IF NOT EXISTS $DB_SCHEMA.modules (
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
);

CREATE TABLE IF NOT EXISTS $DB_SCHEMA.providers (
    id SERIAL PRIMARY KEY,
    host TEXT NOT NULL,
    namespace TEXT NOT NULL,
    name TEXT NOT NULL,
    downloads INTEGER NOT NULL DEFAULT 0,
    module_count INTEGER NOT NULL DEFAULT 0,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(host, namespace, name)
);

CREATE INDEX IF NOT EXISTS idx_modules_namespace_name_provider ON $DB_SCHEMA.modules(namespace, name, provider);
CREATE INDEX IF NOT EXISTS idx_modules_downloads ON $DB_SCHEMA.modules(downloads DESC);
CREATE INDEX IF NOT EXISTS idx_providers_namespace_name ON $DB_SCHEMA.providers(namespace, name);
CREATE INDEX IF NOT EXISTS idx_providers_downloads ON $DB_SCHEMA.providers(downloads DESC);
EOF

echo -e "${GREEN}Database schema created successfully.${RESET}"

# Final instructions
echo
echo -e "${BOLD}Setup Complete!${RESET}"
echo "Your PostgreSQL database is now set up for the OpenTofu Registry API."
echo
echo "To use the database with OpenTofu, run:"
echo "export TOFU_REGISTRY_DB_TYPE=\"postgres\""
echo "export TOFU_REGISTRY_DB_URL=\"$DB_URL\""
echo "or source the created .env file:"
echo "source $ENV_FILE"
echo
echo -e "${YELLOW}Note:${RESET} Make sure to secure your database credentials in production environments."
echo "Consider using a more secure password and restricting database access."
