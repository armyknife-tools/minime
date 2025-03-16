#!/bin/bash
# Test script for OpenTofu Registry Database Integration
# This script tests both PostgreSQL and SQLite database functionality

set -e

# Text formatting
BOLD="\033[1m"
GREEN="\033[0;32m"
YELLOW="\033[0;33m"
RED="\033[0;31m"
RESET="\033[0m"

# Print header
echo -e "${BOLD}OpenTofu Registry Database Integration Test${RESET}"
echo "This script will test both PostgreSQL and SQLite database functionality."
echo

# Create test directory
TEST_DIR="/tmp/opentofu_db_test"
mkdir -p $TEST_DIR

# Function to run a test
run_test() {
    local test_name=$1
    local test_cmd=$2
    
    echo -e "${BOLD}Running test:${RESET} $test_name"
    if eval "$test_cmd"; then
        echo -e "${GREEN}✓ Test passed: $test_name${RESET}"
        return 0
    else
        echo -e "${RED}✗ Test failed: $test_name${RESET}"
        return 1
    fi
}

# Test SQLite integration
test_sqlite() {
    echo -e "\n${BOLD}Testing SQLite Integration${RESET}"
    
    # Set up SQLite environment
    export TOFU_REGISTRY_DB_TYPE="sqlite"
    export TOFU_REGISTRY_DB_URL="file:$TEST_DIR/registry_test.db?cache=shared&mode=rwc"
    
    # Test database connection
    run_test "SQLite Connection" "go run ./cmd/tofu/main.go registry refresh --test-connection"
    
    # Test module refresh
    run_test "SQLite Module Refresh" "go run ./cmd/tofu/main.go registry refresh --type=module --limit=10"
    
    # Test provider refresh
    run_test "SQLite Provider Refresh" "go run ./cmd/tofu/main.go registry refresh --type=provider --limit=10"
    
    # Test search functionality
    run_test "SQLite Module Search" "go run ./cmd/tofu/main.go registry search aws --type=module --limit=5"
    run_test "SQLite Provider Search" "go run ./cmd/tofu/main.go registry search aws --type=provider --limit=5"
    
    # Test install functionality
    run_test "SQLite Module Install" "go run ./cmd/tofu/main.go registry install -type=module hashicorp/consul/aws"
    
    echo -e "\n${GREEN}SQLite tests completed successfully${RESET}"
}

# Test PostgreSQL integration
test_postgres() {
    echo -e "\n${BOLD}Testing PostgreSQL Integration${RESET}"
    
    # Check if PostgreSQL is installed
    if ! command -v psql &> /dev/null; then
        echo -e "${YELLOW}PostgreSQL is not installed. Skipping PostgreSQL tests.${RESET}"
        return 0
    fi
    
    # Create test database
    echo "Creating test PostgreSQL database..."
    PGPASSWORD=postgres psql -h localhost -U postgres -c "DROP DATABASE IF EXISTS opentofu_test;" || true
    PGPASSWORD=postgres psql -h localhost -U postgres -c "CREATE DATABASE opentofu_test;"
    PGPASSWORD=postgres psql -h localhost -U postgres -c "CREATE USER opentofu_test WITH PASSWORD 'opentofu_test';"
    PGPASSWORD=postgres psql -h localhost -U postgres -c "GRANT ALL PRIVILEGES ON DATABASE opentofu_test TO opentofu_test;"
    
    # Set up PostgreSQL environment
    export TOFU_REGISTRY_DB_TYPE="postgres"
    export TOFU_REGISTRY_DB_URL="postgres://opentofu_test:opentofu_test@localhost:5432/opentofu_test?sslmode=disable"
    
    # Test database connection
    run_test "PostgreSQL Connection" "go run ./cmd/tofu/main.go registry refresh --test-connection"
    
    # Test module refresh
    run_test "PostgreSQL Module Refresh" "go run ./cmd/tofu/main.go registry refresh --type=module --limit=10"
    
    # Test provider refresh
    run_test "PostgreSQL Provider Refresh" "go run ./cmd/tofu/main.go registry refresh --type=provider --limit=10"
    
    # Test search functionality
    run_test "PostgreSQL Module Search" "go run ./cmd/tofu/main.go registry search aws --type=module --limit=5"
    run_test "PostgreSQL Provider Search" "go run ./cmd/tofu/main.go registry search aws --type=provider --limit=5"
    
    # Test install functionality
    run_test "PostgreSQL Module Install" "go run ./cmd/tofu/main.go registry install -type=module hashicorp/consul/aws"
    
    # Clean up
    echo "Cleaning up PostgreSQL test database..."
    PGPASSWORD=postgres psql -h localhost -U postgres -c "DROP DATABASE IF EXISTS opentofu_test;"
    PGPASSWORD=postgres psql -h localhost -U postgres -c "DROP USER IF EXISTS opentofu_test;"
    
    echo -e "\n${GREEN}PostgreSQL tests completed successfully${RESET}"
}

# Test environment variable loading
test_env_loading() {
    echo -e "\n${BOLD}Testing Environment Variable Loading${RESET}"
    
    # Create test .env file
    TEST_ENV_FILE="$TEST_DIR/.env.test"
    echo "TOFU_REGISTRY_DB_TYPE=sqlite" > $TEST_ENV_FILE
    echo "TOFU_REGISTRY_DB_URL=file:$TEST_DIR/env_test.db?cache=shared&mode=rwc" >> $TEST_ENV_FILE
    
    # Test loading from .env file
    run_test "Env File Loading" "TOFU_ENV_FILE=$TEST_ENV_FILE go run ./cmd/tofu/main.go registry refresh --test-connection"
    
    echo -e "\n${GREEN}Environment variable loading tests completed successfully${RESET}"
}

# Test database fallback
test_fallback() {
    echo -e "\n${BOLD}Testing Database Fallback${RESET}"
    
    # Test fallback to SQLite when PostgreSQL is unavailable
    export TOFU_REGISTRY_DB_TYPE="postgres"
    export TOFU_REGISTRY_DB_URL="postgres://invalid:invalid@localhost:5432/invalid?sslmode=disable"
    export TOFU_REGISTRY_DB_FALLBACK="true"
    export TOFU_REGISTRY_DB_FALLBACK_URL="file:$TEST_DIR/fallback.db?cache=shared&mode=rwc"
    
    run_test "Database Fallback" "go run ./cmd/tofu/main.go registry refresh --test-connection"
    
    echo -e "\n${GREEN}Database fallback tests completed successfully${RESET}"
}

# Test database schema
test_schema() {
    echo -e "\n${BOLD}Testing Database Schema${RESET}"
    
    # Set up SQLite for schema test
    export TOFU_REGISTRY_DB_TYPE="sqlite"
    export TOFU_REGISTRY_DB_URL="file:$TEST_DIR/schema_test.db?cache=shared&mode=rwc"
    
    # Initialize database
    go run ./cmd/tofu/main.go registry refresh --init-only
    
    # Check if tables exist
    run_test "Module Table Exists" "sqlite3 $TEST_DIR/schema_test.db 'SELECT name FROM sqlite_master WHERE type=\"table\" AND name=\"modules\";' | grep -q modules"
    run_test "Provider Table Exists" "sqlite3 $TEST_DIR/schema_test.db 'SELECT name FROM sqlite_master WHERE type=\"table\" AND name=\"providers\";' | grep -q providers"
    
    echo -e "\n${GREEN}Database schema tests completed successfully${RESET}"
}

# Test performance with large datasets
test_performance() {
    echo -e "\n${BOLD}Testing Performance with Large Datasets${RESET}"
    
    # Set up SQLite for performance test
    export TOFU_REGISTRY_DB_TYPE="sqlite"
    export TOFU_REGISTRY_DB_URL="file:$TEST_DIR/perf_test.db?cache=shared&mode=rwc"
    
    # Test performance with pre-allocated slices
    run_test "Pre-allocation Performance" "time go run ./cmd/tofu/main.go registry refresh --type=module --limit=100"
    
    echo -e "\n${GREEN}Performance tests completed successfully${RESET}"
}

# Run all tests
run_all_tests() {
    test_sqlite
    test_postgres
    test_env_loading
    test_fallback
    test_schema
    test_performance
    
    echo -e "\n${BOLD}${GREEN}All tests completed successfully!${RESET}"
}

# Clean up
cleanup() {
    echo -e "\n${BOLD}Cleaning up test environment${RESET}"
    rm -rf $TEST_DIR
    echo -e "${GREEN}Cleanup completed${RESET}"
}

# Run tests and clean up
run_all_tests
cleanup

echo -e "\n${BOLD}Testing completed.${RESET}"
