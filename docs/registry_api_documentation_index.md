# OpenTofu Registry API Documentation

## Overview

This documentation provides comprehensive information about the OpenTofu Registry API, which enables searching for and importing modules and providers from the Terraform Registry.

## Table of Contents

### Documentation Files

1. [Registry API Documentation](./registry_api_documentation.md) - Comprehensive overview of the Registry API features and architecture
2. [Registry API User Guide](./registry_api_user_guide.md) - Guide for using the Registry API search and import functionality
3. [Diagrams](#diagrams) - Visual representations of the Registry API architecture and workflows

### Diagrams

The following diagrams are available in XML format for use with draw.io:

1. [Registry API Architecture Diagram](./registry_api_architecture_diagram.xml) - Overview of the Registry API architecture
2. [Registry API Import Workflow](./registry_api_import_workflow.xml) - Workflow for importing modules and providers
3. [Registry API Search Workflow](./registry_api_search_workflow.xml) - Workflow for searching modules and providers
4. [Registry API Database Schema](./registry_api_database_schema.xml) - Schema for the PostgreSQL database

## How to Use the Diagrams

1. Go to [draw.io](https://app.diagrams.net/)
2. Click on "Open Existing Diagram"
3. Select "Open from Device"
4. Navigate to the diagram XML file
5. Edit or view the diagram as needed

## Key Features

- Search for modules and providers in the Terraform Registry
- Import modules and providers to a PostgreSQL database
- Verify database counts
- Handle pagination for large datasets
- Implement throttling to avoid rate limiting
- Cache results to improve performance

## Command-Line Usage

```bash
# Search for modules
tofu registry search [options] [keyword]

# Search for providers
tofu registry search -providers [options] [keyword]

# Import modules and providers to PostgreSQL
tofu registry search -import-to-postgres

# Verify database counts
tofu registry search -verify-db-counts
```

## Database Configuration

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

## Performance Considerations

- Pre-allocates data structures based on known registry sizes (4,000 providers, 18,000 modules)
- Implements pagination with appropriate page sizes
- Uses throttling to avoid rate limiting
- Caches results to reduce API calls

## Troubleshooting

See the [Registry API User Guide](./registry_api_user_guide.md#troubleshooting) for troubleshooting information.
