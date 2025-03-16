# OpenTofu Template System

This document describes the technical implementation of the OpenTofu template system, which allows users to generate standardized infrastructure templates for various cloud providers and resources.

## Overview

The template system consists of the following components:

1. **Template Command**: A CLI command that provides a user interface for generating and managing templates.
2. **Template Database**: A database backend that stores template definitions, supporting both PostgreSQL and SQLite.
3. **Template Engine**: A system for rendering templates with customizable variables.

## Database Integration

The template system supports two database backends:

### PostgreSQL

PostgreSQL is the primary database backend for production environments and shared template repositories. The connection details are specified using environment variables, which can be set in a `.env` file:

```
TOFU_REGISTRY_DB_TYPE=postgres
TOFU_REGISTRY_DB_HOST=host.example.com
TOFU_REGISTRY_DB_PORT=5432
TOFU_REGISTRY_DB_NAME=opentofu
TOFU_REGISTRY_DB_USER=username
TOFU_REGISTRY_DB_PASSWORD=password
TOFU_REGISTRY_DB_SSLMODE=require
```

The system will attempt to connect to PostgreSQL first, and if the connection fails, it will fall back to SQLite.

### SQLite

SQLite is used as a fallback database when PostgreSQL is not available or for local development. The SQLite database is stored in `~/.opentofu/templates.db`.

## Code Structure

### Template Command

The template command is implemented in `cmd/tofu/template.go` and provides the CLI interface for the template system. It handles command-line arguments, connects to the database, and manages template generation.

### Template Database

The template database interface is defined in `cmd/tofu/template_db.go` and provides methods for retrieving, storing, and managing templates. It includes functions for:

- Getting a list of available providers
- Getting a list of resources for a provider
- Retrieving a template for a specific resource
- Loading templates into the database

### Template Engine

The template engine is implemented in `internal/templates/templates.go` and provides functions for loading and rendering templates. It includes:

- `ConnectToDatabase`: Connects to either a PostgreSQL or SQLite database
- `LoadTemplates`: Loads built-in templates into the database
- Helper functions for template rendering

## Database Schema

The template database uses the following schema:

```sql
CREATE TABLE IF NOT EXISTS templates (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    provider TEXT NOT NULL,
    resource TEXT NOT NULL,
    display_name TEXT NOT NULL,
    content TEXT NOT NULL,
    description TEXT,
    category TEXT,
    tags TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(provider, resource)
)
```

## Template Format

Templates are stored as HCL (HashiCorp Configuration Language) strings in the database. They include common configuration options for the specified resource, which can be customized by the user.

## Extending the Template System

### Adding New Templates

To add a new template:

1. Create a new HCL file with the template content
2. Add the template to the database using the `template -load` command
3. Verify that the template is available by running `template <provider>/<resource>`

### Custom Template Variables

Templates can include variables that are replaced when the template is generated. Variables are defined using the standard HCL variable syntax:

```hcl
variable "bucket_name" {
  description = "Name of the S3 bucket"
  type        = string
  default     = "my-bucket"
}
```

## Future Enhancements

Planned enhancements for the template system include:

1. **Template Versioning**: Track changes to templates over time
2. **Template Search**: Search for templates by keywords or tags
3. **Template Validation**: Ensure templates follow best practices
4. **User-Contributed Templates**: Allow users to contribute their own templates
