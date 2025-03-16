# OpenTofu Template Generation System Diagram

```
┌─────────────────────────────────────────────────────────────────────────┐
│                                                                         │
│                         OpenTofu CLI Interface                          │
│                                                                         │
└───────────────────────────────────┬─────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                                                                         │
│                          Template Command                               │
│                        (cmd/tofu/template.go)                           │
│                                                                         │
│  ┌───────────────┐    ┌───────────────┐    ┌───────────────────────┐    │
│  │  Parse Args   │───▶│ Load Env Vars │───▶│ Process Command Flags │    │
│  └───────────────┘    └───────────────┘    └───────────────────────┘    │
│                                                       │                  │
└───────────────────────────────────────────────────────┼──────────────────┘
                                                        │
                                                        ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                                                                         │
│                        Template Database Layer                          │
│                       (cmd/tofu/template_db.go)                         │
│                                                                         │
│  ┌───────────────┐    ┌───────────────┐    ┌───────────────────────┐    │
│  │ GetTemplateDB │───▶│ Connect to DB │───▶│ Execute DB Operations │    │
│  └───────────────┘    └───────────────┘    └───────────────────────┘    │
│                                                       │                  │
└───────────────────────────────────────────────────────┼──────────────────┘
                                                        │
                                                        ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                                                                         │
│                       Database Connection Layer                         │
│                    (internal/templates/templates.go)                    │
│                                                                         │
│  ┌────────────────────┐                     ┌────────────────────────┐  │
│  │                    │                     │                        │  │
│  │    PostgreSQL      │◀────Environment────▶│       SQLite          │  │
│  │    Connection      │      Variables      │     Connection         │  │
│  │                    │                     │                        │  │
│  └────────┬───────────┘                     └────────────┬───────────┘  │
│           │                                              │              │
│           └──────────────────┬───────────────────────────┘              │
│                              │                                          │
└──────────────────────────────┼──────────────────────────────────────────┘
                               │
                               ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                                                                         │
│                         Template Storage                                │
│                                                                         │
│  ┌────────────────────┐                     ┌────────────────────────┐  │
│  │                    │                     │                        │  │
│  │    PostgreSQL      │◀───Fallback if ────▶│       SQLite          │  │
│  │    Database        │   Postgres fails    │      Database         │  │
│  │                    │                     │                        │  │
│  └────────┬───────────┘                     └────────────┬───────────┘  │
│           │                                              │              │
│           └──────────────────┬───────────────────────────┘              │
│                              │                                          │
└──────────────────────────────┼──────────────────────────────────────────┘
                               │
                               ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                                                                         │
│                         Template Generation                             │
│                                                                         │
│  ┌───────────────┐    ┌───────────────┐    ┌───────────────────────┐    │
│  │ Retrieve      │───▶│ Process       │───▶│ Write Template        │    │
│  │ Template      │    │ Template      │    │ to File               │    │
│  └───────────────┘    └───────────────┘    └───────────────────────┘    │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

## Template Generation Flow

1. **User Interaction**:
   - User runs `tofu template [provider]/[resource]` command
   - CLI parses arguments and flags

2. **Environment Setup**:
   - Load environment variables from `.env` file
   - Set up database connection parameters

3. **Database Connection**:
   - Attempt to connect to PostgreSQL using environment variables
   - Fall back to SQLite if PostgreSQL connection fails

4. **Template Retrieval**:
   - Query database for template matching provider/resource
   - Process template content (HCL format)

5. **Template Generation**:
   - Apply any customizations or variables
   - Write template to output file (default: resource_name.tf)

6. **User Feedback**:
   - Display success message with path to generated file

## Database Fallback Mechanism

```
┌────────────────┐     ┌───────────────┐     ┌───────────────┐
│                │     │               │     │               │
│  Start DB      │────▶│  Try          │────▶│  Connected    │
│  Connection    │     │  PostgreSQL   │  ✓  │  Successfully │
│                │     │  Connection   │     │               │
└────────────────┘     └───────┬───────┘     └───────────────┘
                               │
                               │ ✗ Connection Failed
                               ▼
                       ┌───────────────┐     ┌───────────────┐
                       │               │     │               │
                       │  Fall back   │────▶│  Connected    │
                       │  to SQLite   │  ✓  │  Successfully │
                       │              │     │               │
                       └───────┬───────┘     └───────────────┘
                               │
                               │ ✗ Connection Failed
                               ▼
                       ┌───────────────┐
                       │               │
                       │  Return       │
                       │  Error        │
                       │               │
                       └───────────────┘
```

## Template Loading Process

```
┌────────────────┐     ┌───────────────┐     ┌───────────────┐
│                │     │               │     │               │
│  Start         │────▶│  Connect to   │────▶│  Create       │
│  Template Load │     │  Database     │     │  Tables       │
│                │     │               │     │               │
└────────────────┘     └───────────────┘     └───────┬───────┘
                                                     │
                                                     ▼
┌────────────────┐     ┌───────────────┐     ┌───────────────┐
│                │     │               │     │               │
│  Success       │◀────│  Insert       │◀────│  Load Built-in│
│  Message       │     │  Templates    │     │  Templates    │
│                │     │               │     │               │
└────────────────┘     └───────────────┘     └───────────────┘
```
