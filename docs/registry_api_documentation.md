# OpenTofu Registry API Documentation

## Overview

The OpenTofu Registry API provides a centralized repository for modules and providers that can be used in OpenTofu configurations. This documentation covers the features, usage, and architecture of the Registry API implementation.

## Features

### Module and Provider Search
- Search for modules and providers using keywords
- Filter results by namespace, name, or provider
- Display detailed information about modules and providers
- Sort results by popularity, downloads, or relevance

### Module and Provider Import
- Import modules and providers from the Terraform Registry API
- Store module and provider data in a PostgreSQL database
- Import all versions of modules and providers
- Verify database counts after import

### Caching
- Cache search results to improve performance
- Automatically refresh cache when needed
- Configurable cache expiration

### Rate Limiting Protection
- Throttling for API requests to avoid rate limiting
- Automatic retry with exponential backoff
- Error handling for rate limiting scenarios

## Command-Line Usage

### Search for Modules

```bash
tofu registry search [options] [keyword]
```

#### Options:
- `-detailed`: Show detailed information about each module
- `-json`: Output results in JSON format
- `-limit=n`: Limit the number of results (default: 10)
- `-provider=provider`: Filter modules by provider

### Search for Providers

```bash
tofu registry search -providers [options] [keyword]
```

#### Options:
- `-detailed`: Show detailed information about each provider
- `-json`: Output results in JSON format
- `-limit=n`: Limit the number of results (default: 10)

### Import Modules and Providers to PostgreSQL

```bash
tofu registry search -import-to-postgres
```

### Verify Database Counts

```bash
tofu registry search -verify-db-counts
```

## Architecture

The Registry API implementation consists of several components:

1. **Search Command**: The main entry point for searching modules and providers
2. **Registry Client**: Handles communication with the Terraform Registry API
3. **Database Client**: Manages connections to the PostgreSQL database
4. **Cache Manager**: Handles caching of search results
5. **Import Manager**: Imports modules and providers from the Terraform Registry API

## Database Schema

### Modules Table

| Column       | Type         | Description                              |
|--------------|--------------|------------------------------------------|
| namespace    | VARCHAR(255) | Module namespace                         |
| name         | VARCHAR(255) | Module name                              |
| provider     | VARCHAR(255) | Provider name                            |
| version      | VARCHAR(50)  | Module version                           |
| download_url | TEXT         | URL to download the module               |
| published_at | TIMESTAMP    | Publication date                         |

Primary Key: (namespace, name, provider, version)

### Providers Table

| Column       | Type         | Description                              |
|--------------|--------------|------------------------------------------|
| namespace    | VARCHAR(255) | Provider namespace                       |
| name         | VARCHAR(255) | Provider name                            |
| version      | VARCHAR(50)  | Provider version                         |
| platforms    | TEXT         | JSON array of supported platforms        |
| download_url | TEXT         | URL to download the provider             |

Primary Key: (namespace, name, version)

## Configuration

The Registry API uses environment variables for configuration:

```
TOFU_REGISTRY_DB_TYPE=postgres
TOFU_REGISTRY_DB_HOST=<hostname>
TOFU_REGISTRY_DB_PORT=<port>
TOFU_REGISTRY_DB_NAME=<database>
TOFU_REGISTRY_DB_USER=<username>
TOFU_REGISTRY_DB_PASSWORD=<password>
TOFU_REGISTRY_DB_SSLMODE=<sslmode>
```

These variables can be set in a `.env` file in the current directory.

## Error Handling

The Registry API implements robust error handling:

1. **API Errors**: Retries with exponential backoff
2. **Database Errors**: Logs errors and continues with other operations
3. **Rate Limiting**: Waits and retries after appropriate delays
4. **Pagination Errors**: Implements safeguards against infinite loops

## Performance Considerations

- Pre-allocates data structures based on known registry sizes (4,000 providers, 18,000 modules)
- Implements pagination with appropriate page sizes
- Uses throttling to avoid rate limiting
- Caches results to reduce API calls
