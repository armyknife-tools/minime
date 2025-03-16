# OpenTofu Registry API User Guide

## Introduction

The OpenTofu Registry API allows users to search for and import modules and providers from the Terraform Registry. This guide explains how to use the registry search and import functionality.

## Prerequisites

- OpenTofu installed on your system
- PostgreSQL database configured (if using import functionality)
- Internet connection to access the Terraform Registry API

## Configuration

The Registry API uses environment variables for configuration. These can be set in a `.env` file in the current directory:

```
TOFU_REGISTRY_DB_TYPE=postgres
TOFU_REGISTRY_DB_HOST=<hostname>
TOFU_REGISTRY_DB_PORT=<port>
TOFU_REGISTRY_DB_NAME=<database>
TOFU_REGISTRY_DB_USER=<username>
TOFU_REGISTRY_DB_PASSWORD=<password>
TOFU_REGISTRY_DB_SSLMODE=<sslmode>
```

## Searching for Modules

To search for modules in the Terraform Registry, use the following command:

```bash
tofu registry search [options] [keyword]
```

### Options:

- `-detailed`: Show detailed information about each module
- `-json`: Output results in JSON format
- `-limit=n`: Limit the number of results (default: 10)
- `-provider=provider`: Filter modules by provider

### Examples:

```bash
# Search for modules containing "vpc"
tofu registry search vpc

# Search for AWS VPC modules with detailed information
tofu registry search -detailed -provider=aws vpc

# Search for modules and output in JSON format
tofu registry search -json kubernetes

# Limit search results to 5
tofu registry search -limit=5 database
```

## Searching for Providers

To search for providers in the Terraform Registry, use the following command:

```bash
tofu registry search -providers [options] [keyword]
```

### Options:

- `-detailed`: Show detailed information about each provider
- `-json`: Output results in JSON format
- `-limit=n`: Limit the number of results (default: 10)

### Examples:

```bash
# Search for providers containing "aws"
tofu registry search -providers aws

# Search for providers with detailed information
tofu registry search -providers -detailed azure

# Search for providers and output in JSON format
tofu registry search -providers -json google

# Limit provider search results to 5
tofu registry search -providers -limit=5 kubernetes
```

## Importing Modules and Providers to PostgreSQL

To import all modules and providers from the Terraform Registry to a PostgreSQL database, use the following command:

```bash
tofu registry search -import-to-postgres
```

This command will:

1. Count the total number of modules and providers in the Terraform Registry
2. Fetch all module IDs and their versions
3. Fetch all provider IDs and their versions
4. Import all modules and providers into the PostgreSQL database

The import process may take several hours to complete, as it needs to fetch approximately 18,000 modules and 4,000 providers, each with multiple versions.

## Verifying Database Counts

To verify the counts of modules and providers in the PostgreSQL database, use the following command:

```bash
tofu registry search -verify-db-counts
```

This command will query the PostgreSQL database and display the counts of modules and providers.

## Troubleshooting

### Rate Limiting

The Terraform Registry API has rate limits. If you encounter rate limiting errors, the OpenTofu Registry API will automatically retry with exponential backoff.

### Database Connection Issues

If you encounter database connection issues, check your database configuration in the `.env` file and ensure that the database server is running and accessible.

### Import Process Hangs

If the import process appears to hang, it may be due to rate limiting or network issues. The process includes safeguards against infinite loops, but in rare cases, you may need to restart the import process.

## Performance Considerations

- The Registry API pre-allocates data structures based on known registry sizes (4,000 providers, 18,000 modules)
- Pagination is used to handle large datasets
- Throttling is implemented to avoid rate limiting
- Results are cached to reduce API calls

## Advanced Usage

### Using a Different Registry

By default, the Registry API uses the Terraform Registry (`registry.terraform.io`). To use a different registry, you can specify the registry host:

```bash
tofu registry search -registry=my-registry.example.com vpc
```

### Debugging

To enable debug logging, set the `TF_LOG` environment variable:

```bash
TF_LOG=debug tofu registry search vpc
```

## Conclusion

The OpenTofu Registry API provides a powerful way to search for and import modules and providers from the Terraform Registry. By following this guide, you can effectively use the registry search and import functionality to enhance your OpenTofu workflows.
