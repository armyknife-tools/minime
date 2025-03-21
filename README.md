# OpenTofu

- [Manifesto](https://opentofu.org/manifesto)
- [About the OpenTofu fork](https://opentofu.org/fork)
- [How to install](https://opentofu.org/docs/intro/install)
- [Join our Slack community!](https://opentofu.org/slack)
- [Weekly OpenTofu Status Updates](WEEKLY_UPDATES.md)

![](https://raw.githubusercontent.com/opentofu/brand-artifacts/main/full/transparent/SVG/on-dark.svg#gh-dark-mode-only)
![](https://raw.githubusercontent.com/opentofu/brand-artifacts/main/full/transparent/SVG/on-light.svg#gh-light-mode-only)

OpenTofu is an OSS tool for building, changing, and versioning infrastructure safely and efficiently. OpenTofu can manage existing and popular service providers as well as custom in-house solutions.

The key features of OpenTofu are:

- **Infrastructure as Code**: Infrastructure is described using a high-level configuration syntax. This allows a blueprint of your datacenter to be versioned and treated as you would any other code. Additionally, infrastructure can be shared and re-used.

- **Execution Plans**: OpenTofu has a "planning" step where it generates an execution plan. The execution plan shows what OpenTofu will do when you call apply. This lets you avoid any surprises when OpenTofu manipulates infrastructure.

- **Resource Graph**: OpenTofu builds a graph of all your resources, and parallelizes the creation and modification of any non-dependent resources. Because of this, OpenTofu builds infrastructure as efficiently as possible, and operators get insight into dependencies in their infrastructure.

- **Change Automation**: Complex changesets can be applied to your infrastructure with minimal human interaction. With the previously mentioned execution plan and resource graph, you know exactly what OpenTofu will change and in what order, avoiding many possible human errors.

## Getting help and contributing

- Have a question? Post it in [GitHub Discussions](https://github.com/orgs/opentofu/discussions) or on the [OpenTofu Slack](https://opentofu.org/slack/)!
- Want to contribute? Please read the [Contribution Guide](CONTRIBUTING.md).
- Want to stay up to date? Read the [weekly updates](WEEKLY_UPDATES.md), [TSC summary](TSC_SUMMARY.md), or join the [community meetings](https://meet.google.com/xfm-cgms-has) on Wednesdays at 14:30 CET / 8:30 AM Eastern / 5:30 AM Western / 19:00 India time on this link: https://meet.google.com/xfm-cgms-has ([📅 calendar link](https://calendar.google.com/calendar/event?eid=NDg0aWl2Y3U1aHFva3N0bGhyMHBhNzdpZmsgY18zZjJkZDNjMWZlMGVmNGU5M2VmM2ZjNDU2Y2EyZGQyMTlhMmU4ZmQ4NWY2YjQwNzUwYWYxNmMzZGYzNzBiZjkzQGc))

> [!TIP]
> For more OpenTofu events, subscribe to the [OpenTofu Events Calendar](https://calendar.google.com/calendar/embed?src=c_3f2dd3c1fe0ef4e93ef3fc456ca2dd219a2e8fd85f6b40750af16c3df370bf93%40group.calendar.google.com)!

## Recent Improvements

### Registry API Enhancements

The OpenTofu Registry API has been significantly improved to enhance reliability and performance:

- **Robust Provider Fetching**: Implemented a multi-stage approach to fetching providers from the Terraform Registry
- **Enhanced Module Search**: Improved the search functionality for modules with better caching and error handling
- **Performance Optimizations**: Pre-allocated data structures based on known registry sizes (4,000 providers, 18,000 modules)
- **Better Error Handling**: Added comprehensive error handling and logging for registry operations

For more details, see the [Registry API Improvements documentation](docs/registry_api_improvements.md).

## Reporting security vulnerabilities
If you've found a vulnerability or a potential vulnerability in OpenTofu please follow [Security Policy](https://github.com/opentofu/opentofu/security/policy). We'll send a confirmation email to acknowledge your report, and we'll send an additional email when we've identified the issue positively or negatively.

## Reporting possible copyright issues

If you believe you have found any possible copyright or intellectual property issues, please contact liaison@opentofu.org. We'll send a confirmation email to acknowledge your report.

## Registry Access

In an effort to comply with applicable sanctions, we block access from specific countries of origin.

## Template System

OpenTofu includes a template system that allows you to generate infrastructure-as-code templates for various cloud providers. This feature helps you quickly scaffold common cloud resources with best practices built-in.

### Using Templates

#### Listing Available Templates

To see all available template providers:
```bash
tofu template
```

To list templates for a specific provider:
```bash
tofu template aws
```

#### Generating a Template

To generate a template for a specific resource:
```bash
tofu template aws/s3
```

This will create a file named after the resource (e.g., `s3.tf`) in your current directory.

### Creating Custom Templates

You can create custom templates by adding them to the template database. Templates are defined in the `internal/templates` package.

1. **Create a Template Definition**:
   
   Templates are defined in provider-specific files (e.g., `aws_templates.go`, `azure_templates.go`, `gcp_templates.go`). To add a new template:

   ```go
   // Add to the appropriate generateXTemplates function
   templates = append(templates, Template{
       Provider:    "aws",
       Resource:    "my_resource",
       DisplayName: "My Resource",
       Description: "Description of my resource",
       Category:    "Category",
       Tags:        "tag1,tag2",
       Content:     `# My Resource Template
   
   resource "aws_my_resource" "example" {
     name = "example"
     # ... other properties
   }`,
   })
   ```

2. **Build OpenTofu**:
   ```bash
   go build -o tofu ./cmd/tofu
   ```

3. **Load Templates into the Database**:
   ```bash
   tofu template -load
   ```

   This command supports both SQLite (default) and PostgreSQL databases:
   ```bash
   # For SQLite (default)
   tofu template -load

   # For PostgreSQL (using environment variables)
   export TOFU_REGISTRY_DB_TYPE=postgres
   export TOFU_REGISTRY_DB_HOST=your-postgres-host
   export TOFU_REGISTRY_DB_PORT=5432
   export TOFU_REGISTRY_DB_USER=your-user
   export TOFU_REGISTRY_DB_PASSWORD=your-password
   export TOFU_REGISTRY_DB_NAME=your-dbname
   export TOFU_REGISTRY_DB_SSLMODE=require
   tofu template -load
   ```

4. **Verify Your Template**:
   ```bash
   # List templates for your provider
   tofu template your-provider
   
   # Generate your template
   tofu template your-provider/your-resource
   ```

### Database Configuration

The template system supports both SQLite and PostgreSQL databases:

- **SQLite**: Used by default, with database stored at `~/.opentofu/tofu.db`
- **PostgreSQL**: Configured using environment variables:
  - `TOFU_REGISTRY_DB_TYPE=postgres`
  - `TOFU_REGISTRY_DB_HOST=host`
  - `TOFU_REGISTRY_DB_PORT=port`
  - `TOFU_REGISTRY_DB_USER=user`
  - `TOFU_REGISTRY_DB_PASSWORD=password`
  - `TOFU_REGISTRY_DB_NAME=dbname`
  - `TOFU_REGISTRY_DB_SSLMODE=sslmode`

You can also use the database configuration commands:
```bash
# Set up the database
tofu db setup

# Configure database connection
tofu db configure
```

## License

[Mozilla Public License v2.0](https://github.com/opentofu/opentofu/blob/main/LICENSE)
