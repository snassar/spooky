# Variables System User Guide

## Overview

The spooky variables system provides comprehensive management of configuration variables for automation and orchestration. This guide covers everything from basic variable configuration to advanced features like multi-file variables, dependency resolution, and validation.

## Table of Contents

1. [Getting Started](#getting-started)
2. [Variable Configuration](#variable-configuration)
3. [Variable Management](#variable-management)
4. [Variable Resolution](#variable-resolution)
5. [Validation and Troubleshooting](#validation-and-troubleshooting)
6. [Advanced Features](#advanced-features)
7. [Best Practices](#best-practices)
8. [Examples](#examples)

## Getting Started

### Basic Variable Configuration

A variable configuration in spooky consists of one or more HCL files that define configuration variables with their types, default values, dependencies, and validation rules.

**Single File Configuration (`variables.hcl`):**
```hcl
variables {
  variable "app_name" {
    type = "string"
    description = "Application name"
    default = "my-app"
    scope = "project"
  }

  variable "app_version" {
    type = "string"
    description = "Application version"
    default = "1.0.0"
    scope = "project"
  }

  variable "environment" {
    type = "string"
    description = "Deployment environment"
    default = "development"
    scope = "project"
    
    validation {
      allowed_values = ["development", "staging", "production"]
    }
  }

  variable "port" {
    type = "number"
    description = "Application port"
    default = 8080
    scope = "project"
    
    constraints {
      min_value = 1024
      max_value = 65535
    }
  }

  variable "debug" {
    type = "bool"
    description = "Enable debug mode"
    default = false
    scope = "project"
  }
}
```

### Multi-File Variable Configuration

For larger projects, you can organize variables into multiple files within a `variables/` directory:

**Project Structure:**
```
my-project/
├── project.hcl
├── variables.hcl                    # Global variables
└── variables/
    ├── app.hcl                      # Application variables
    ├── database.hcl                 # Database variables
    └── security.hcl                 # Security variables
```

**Application Variables (`variables/app.hcl`):**
```hcl
variables {
  variable "app_name" {
    type = "string"
    description = "Application name"
    default = "my-app"
    scope = "project"
  }

  variable "app_version" {
    type = "string"
    description = "Application version"
    default = "1.0.0"
    scope = "project"
  }

  variable "app_port" {
    type = "number"
    description = "Application port"
    default = 8080
    scope = "project"
    
    constraints {
      min_value = 1024
      max_value = 65535
    }
  }
}
```

**Database Variables (`variables/database.hcl`):**
```hcl
variables {
  variable "db_host" {
    type = "string"
    description = "Database host"
    default = "localhost"
    scope = "project"
  }

  variable "db_port" {
    type = "number"
    description = "Database port"
    default = 5432
    scope = "project"
  }

  variable "db_name" {
    type = "string"
    description = "Database name"
    default = "myapp"
    scope = "project"
  }

  variable "db_user" {
    type = "string"
    description = "Database user"
    scope = "project"
    sensitive = true
  }

  variable "db_password" {
    type = "string"
    description = "Database password"
    scope = "project"
    sensitive = true
    encrypted = true
  }
}
```

## Variable Configuration

### Variable Types

The variables system supports several data types:

#### String Variables
```hcl
variable "app_name" {
  type = "string"
  description = "Application name"
  default = "my-app"
  scope = "project"
}
```

#### Number Variables
```hcl
variable "port" {
  type = "number"
  description = "Application port"
  default = 8080
  scope = "project"
  
  constraints {
    min_value = 1024
    max_value = 65535
  }
}
```

#### Boolean Variables
```hcl
variable "debug" {
  type = "bool"
  description = "Enable debug mode"
  default = false
  scope = "project"
}
```

#### List Variables
```hcl
variable "allowed_hosts" {
  type = "list"
  description = "List of allowed hosts"
  default = ["localhost", "127.0.0.1"]
  scope = "project"
}
```

#### Map Variables
```hcl
variable "config" {
  type = "map"
  description = "Configuration map"
  default = {
    timeout = 30
    retries = 3
    cache_size = 1000
  }
  scope = "project"
}
```

### Variable Attributes

#### Required Attributes

- **`type`** - The data type of the variable (string, number, bool, list, map)
- **`scope`** - The scope of the variable (project, global, local)

#### Optional Attributes

- **`description`** - Human-readable description of the variable
- **`default`** - Default value for the variable
- **`required`** - Whether the variable is required (default: false)
- **`sensitive`** - Whether the variable contains sensitive data (default: false)
- **`encrypted`** - Whether the variable should be encrypted (default: false)
- **`dependencies`** - List of other variables this variable depends on
- **`validation`** - Validation rules for the variable
- **`constraints`** - Type-specific constraints
- **`metadata`** - Additional metadata for the variable

### Validation Rules

#### Allowed Values
```hcl
variable "environment" {
  type = "string"
  description = "Deployment environment"
  default = "development"
  scope = "project"
  
  validation {
    allowed_values = ["development", "staging", "production"]
  }
}
```

#### Pattern Matching
```hcl
variable "email" {
  type = "string"
  description = "Contact email"
  scope = "project"
  
  validation {
    pattern = "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"
  }
}
```

#### Custom Validation
```hcl
variable "port" {
  type = "number"
  description = "Application port"
  scope = "project"
  
  validation {
    min_value = 1024
    max_value = 65535
  }
}
```

### Constraints

#### String Constraints
```hcl
variable "username" {
  type = "string"
  description = "Username"
  scope = "project"
  
  constraints {
    min_length = 3
    max_length = 20
    pattern = "^[a-zA-Z0-9_-]+$"
  }
}
```

#### Number Constraints
```hcl
variable "timeout" {
  type = "number"
  description = "Timeout in seconds"
  default = 30
  scope = "project"
  
  constraints {
    min_value = 1
    max_value = 300
  }
}
```

#### List Constraints
```hcl
variable "tags" {
  type = "list"
  description = "Resource tags"
  scope = "project"
  
  constraints {
    min_items = 1
    max_items = 10
  }
}
```

## Variable Management

### Listing Variables

Use the `spooky variables list` command to view all variables in a project:

```bash
# List variables in the current directory
spooky variables list

# List variables in a specific project
spooky variables list ./my-project

# List variables with JSON output
spooky variables list --json
```

**Example Output:**
```
Variables in ./my-project:

📁 variables.hcl:
  • app_name (string) = "my-app" [project]
  • app_version (string) = "1.0.0" [project]
  • environment (string) = "development" [project]

📁 variables/app.hcl:
  • app_port (number) = 8080 [project]
  • debug (bool) = false [project]

📁 variables/database.hcl:
  • db_host (string) = "localhost" [project]
  • db_port (number) = 5432 [project]
  • db_name (string) = "myapp" [project]
  • db_user (string) = <sensitive> [project]
  • db_password (string) = <encrypted> [project]

Total: 10 variables across 3 files
```

### Validating Variables

Use the `spooky variables validate` command to validate your variable configuration:

```bash
# Validate variables in the current directory
spooky variables validate

# Validate variables in a specific project
spooky variables validate ./my-project

# Validate with verbose output
spooky variables validate --verbose
```

**Example Output:**
```
✅ Validation successful

Variables validated: 10
Files processed: 3
Warnings: 0
Errors: 0

Validation details:
• All variable names are valid
• No circular dependencies detected
• All constraints are satisfied
• All validation rules passed
```

### Resolving Variables

Use the `spooky variables resolve` command to resolve variables with context:

```bash
# Resolve variables in the current directory
spooky variables resolve

# Resolve variables with specific context
spooky variables resolve --context production

# Resolve with environment overrides
spooky variables resolve --env-override
```

**Example Output:**
```
Resolved Variables:

app_name = "my-app"
app_version = "1.0.0"
environment = "production"  # Overridden by environment
port = 8080
debug = false
db_host = "prod-db.example.com"  # Overridden by environment
db_port = 5432
db_name = "myapp_prod"  # Overridden by environment
db_user = "prod_user"  # Overridden by environment
db_password = "***"  # Encrypted

Resolution completed successfully
Dependencies resolved: 0 circular dependencies detected
Environment overrides applied: 5 variables
```

## Variable Resolution

### Dependency Resolution

Variables can depend on other variables, and spooky automatically resolves dependencies:

```hcl
variables {
  variable "app_name" {
    type = "string"
    description = "Application name"
    default = "my-app"
    scope = "project"
  }

  variable "app_title" {
    type = "string"
    description = "Application title"
    default = "My Application"
    scope = "project"
    dependencies = ["app_name"]
  }

  variable "app_url" {
    type = "string"
    description = "Application URL"
    scope = "project"
    dependencies = ["app_name", "environment"]
  }
}
```

### Environment Variable Overrides

Variables can be overridden by environment variables using the `SPOOKY_VAR_` prefix:

```bash
# Override variables with environment variables
export SPOOKY_VAR_APP_NAME="production-app"
export SPOOKY_VAR_ENVIRONMENT="production"
export SPOOKY_VAR_DB_HOST="prod-db.example.com"

# Resolve variables with overrides
spooky variables resolve --env-override
```

### Context-Based Resolution

Variables can be resolved in different contexts:

```hcl
variables {
  variable "api_url" {
    type = "string"
    description = "API URL"
    scope = "project"
  }

  variable "api_timeout" {
    type = "number"
    description = "API timeout"
    default = 30
    scope = "project"
  }
}
```

```bash
# Resolve with development context
spooky variables resolve --context development

# Resolve with production context
spooky variables resolve --context production
```

## Validation and Troubleshooting

### Common Validation Errors

#### Duplicate Variable Names
```
❌ Validation failed: duplicate variable name 'app_name' found in multiple files
  • variables.hcl:5:1
  • variables/app.hcl:3:1
```

**Solution:** Rename one of the variables or remove the duplicate definition.

#### Circular Dependencies
```
❌ Validation failed: circular dependency detected: app_name -> app_title -> app_name
```

**Solution:** Restructure your variables to remove circular dependencies.

#### Invalid Variable Names
```
❌ Validation failed: invalid variable name 'app-name': must match pattern ^[a-zA-Z_][a-zA-Z0-9_]*$
```

**Solution:** Use valid variable names (letters, numbers, underscores, starting with letter or underscore).

#### Missing Required Variables
```
❌ Validation failed: required variable 'db_password' has no default value
```

**Solution:** Provide a default value or set the variable via environment variable.

### Troubleshooting Commands

#### Verbose Validation
```bash
# Get detailed validation information
spooky variables validate --verbose
```

#### Dependency Analysis
```bash
# Analyze variable dependencies
spooky variables resolve --show-dependencies
```

#### Environment Variable Debugging
```bash
# Show which environment variables are being used
spooky variables resolve --debug-env
```

## Advanced Features

### Variable Scopes

Variables can have different scopes:

#### Project Scope
```hcl
variable "app_name" {
  type = "string"
  scope = "project"  # Available within the project
}
```

#### Global Scope
```hcl
variable "global_config" {
  type = "string"
  scope = "global"  # Available across all projects
}
```

#### Local Scope
```hcl
variable "local_setting" {
  type = "string"
  scope = "local"  # Available only in current context
}
```

### Sensitive Variables

Mark variables as sensitive to protect sensitive data:

```hcl
variable "api_key" {
  type = "string"
  description = "API key"
  scope = "project"
  sensitive = true  # Will be masked in output
}

variable "password" {
  type = "string"
  description = "Database password"
  scope = "project"
  sensitive = true
  encrypted = true  # Will be encrypted in storage
}
```

### Complex Dependencies

Variables can have complex dependency relationships:

```hcl
variables {
  variable "base_url" {
    type = "string"
    description = "Base URL"
    default = "https://api.example.com"
    scope = "project"
  }

  variable "api_version" {
    type = "string"
    description = "API version"
    default = "v1"
    scope = "project"
  }

  variable "api_url" {
    type = "string"
    description = "Full API URL"
    scope = "project"
    dependencies = ["base_url", "api_version"]
  }

  variable "timeout" {
    type = "number"
    description = "Request timeout"
    default = 30
    scope = "project"
  }

  variable "retry_count" {
    type = "number"
    description = "Retry count"
    default = 3
    scope = "project"
  }

  variable "api_config" {
    type = "map"
    description = "API configuration"
    scope = "project"
    dependencies = ["api_url", "timeout", "retry_count"]
  }
}
```

## Best Practices

### Variable Organization

1. **Group by Functionality** - Organize variables by their purpose (app, database, security, etc.)
2. **Use Descriptive Names** - Choose clear, descriptive variable names
3. **Consistent Naming** - Follow consistent naming conventions across your project
4. **Documentation** - Always provide descriptions for your variables

### Security Best Practices

1. **Sensitive Data** - Mark sensitive variables with `sensitive = true`
2. **Encryption** - Use `encrypted = true` for highly sensitive data
3. **Environment Variables** - Use environment variables for secrets
4. **Access Control** - Limit access to variable files containing sensitive data

### Dependency Management

1. **Minimize Dependencies** - Keep dependency graphs simple
2. **Avoid Circular Dependencies** - Design your variables to avoid circular references
3. **Document Dependencies** - Clearly document why variables depend on each other
4. **Test Dependencies** - Regularly test dependency resolution

### Validation and Constraints

1. **Use Validation Rules** - Always validate important variables
2. **Appropriate Constraints** - Use constraints to prevent invalid values
3. **Test Validation** - Test your validation rules with various inputs
4. **Clear Error Messages** - Provide clear error messages for validation failures

## Examples

### Basic Web Application

**`variables.hcl`:**
```hcl
variables {
  variable "app_name" {
    type = "string"
    description = "Application name"
    default = "web-app"
    scope = "project"
  }

  variable "app_port" {
    type = "number"
    description = "Application port"
    default = 8080
    scope = "project"
    
    constraints {
      min_value = 1024
      max_value = 65535
    }
  }

  variable "environment" {
    type = "string"
    description = "Deployment environment"
    default = "development"
    scope = "project"
    
    validation {
      allowed_values = ["development", "staging", "production"]
    }
  }

  variable "debug" {
    type = "bool"
    description = "Enable debug mode"
    default = false
    scope = "project"
  }
}
```

### Database Configuration

**`variables/database.hcl`:**
```hcl
variables {
  variable "db_host" {
    type = "string"
    description = "Database host"
    default = "localhost"
    scope = "project"
  }

  variable "db_port" {
    type = "number"
    description = "Database port"
    default = 5432
    scope = "project"
  }

  variable "db_name" {
    type = "string"
    description = "Database name"
    default = "myapp"
    scope = "project"
  }

  variable "db_user" {
    type = "string"
    description = "Database user"
    scope = "project"
    sensitive = true
  }

  variable "db_password" {
    type = "string"
    description = "Database password"
    scope = "project"
    sensitive = true
    encrypted = true
  }

  variable "db_ssl" {
    type = "bool"
    description = "Enable SSL for database connection"
    default = false
    scope = "project"
  }
}
```

### API Configuration

**`variables/api.hcl`:**
```hcl
variables {
  variable "api_base_url" {
    type = "string"
    description = "API base URL"
    default = "https://api.example.com"
    scope = "project"
  }

  variable "api_version" {
    type = "string"
    description = "API version"
    default = "v1"
    scope = "project"
  }

  variable "api_timeout" {
    type = "number"
    description = "API request timeout in seconds"
    default = 30
    scope = "project"
    
    constraints {
      min_value = 1
      max_value = 300
    }
  }

  variable "api_retries" {
    type = "number"
    description = "Number of API retries"
    default = 3
    scope = "project"
    
    constraints {
      min_value = 0
      max_value = 10
    }
  }

  variable "api_key" {
    type = "string"
    description = "API authentication key"
    scope = "project"
    sensitive = true
  }
}
```

This user guide provides a comprehensive overview of the spooky variables system. For more detailed technical information, refer to the [API Reference](VARIABLES_API_REFERENCE.md), and for troubleshooting help, see the [Troubleshooting Guide](VARIABLES_TROUBLESHOOTING.md).
