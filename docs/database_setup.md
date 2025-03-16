# OpenTofu Registry Database Setup Guide

This document provides instructions for setting up and configuring the PostgreSQL database for the OpenTofu Registry API.

## Prerequisites

- PostgreSQL 12.0 or higher
- Go 1.18 or higher
- Access to create databases and users in PostgreSQL

## Database Configuration

The OpenTofu Registry API uses PostgreSQL for storing module and provider metadata. Follow these steps to set up your database:

### 1. Create a PostgreSQL Database

You can create the database manually or use the provided script:

```bash
go run ./scripts/create_opentofu_db.go
```

This script will:
- Connect to your PostgreSQL server
- Create a database named `opentofu` if it doesn't exist
- Create the necessary tables for modules, providers, and caching

### 2. Create a Dedicated Database User

For security best practices, create a dedicated user for the OpenTofu Registry API:

```bash
go run ./scripts/create_opentofu_user.go
```

This script will:
- Create a user named `opentofu_user` with a secure random password
- Grant the necessary permissions to this user
- Update your `.env` file with the new credentials

### 3. Configure Environment Variables

The OpenTofu Registry API uses the following environment variables for database configuration:

```
TOFU_REGISTRY_DB_TYPE=postgres
TOFU_REGISTRY_DB_HOST=<your_host>
TOFU_REGISTRY_DB_PORT=<your_port>
TOFU_REGISTRY_DB_NAME=opentofu
TOFU_REGISTRY_DB_USER=opentofu_user
TOFU_REGISTRY_DB_PASSWORD=<your_password>
TOFU_REGISTRY_DB_SSLMODE=require
```

Create a `.env` file in the root directory of the project with these variables.

### 4. Test the Database Connection

Verify your database configuration:

```bash
go run ./scripts/test_opentofu_user_connection.go
```

This script will test:
- Connection to the database
- Table creation
- Data insertion and retrieval

## Database Schema

The OpenTofu Registry API uses the following tables:

### Modules Table

```sql
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
```

### Providers Table

```sql
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
```

### Cache Table

```sql
CREATE TABLE IF NOT EXISTS cache (
    id SERIAL PRIMARY KEY,
    cache_key VARCHAR(255) NOT NULL UNIQUE,
    cache_value TEXT NOT NULL,
    expires_at TIMESTAMP NOT NULL
)
```

## Performance Considerations

The OpenTofu Registry API is designed to handle a large volume of data:
- Approximately 4,000 providers
- Approximately 18,000 modules

To optimize performance:
- Pre-allocated data structures are used based on these known registry sizes
- Caching mechanisms are implemented to reduce API calls
- Indexes are created on frequently queried columns

## Security Considerations

The database setup follows these security best practices:
- Uses a dedicated database for the OpenTofu Registry API
- Uses a dedicated user with limited permissions
- Requires SSL for database connections
- Stores credentials in environment variables, not in code

## Troubleshooting

### Connection Issues

If you encounter connection issues:
1. Verify your database host and port
2. Check that your database user has the necessary permissions
3. Ensure SSL is properly configured if required

### Permission Issues

If you encounter permission issues:
1. Verify that your database user has been granted the necessary permissions
2. Check that the user has access to the `public` schema
3. Ensure the user can create, read, update, and delete data in the tables

## Additional Resources

- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [Go Database/SQL Package](https://golang.org/pkg/database/sql/)
- [OpenTofu Registry API Documentation](https://opentofu.org/docs/registry/api/)
