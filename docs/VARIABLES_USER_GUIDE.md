# Variables System User Guide

## Overview

The spooky variables system provides comprehensive variable management, resolution, and validation capabilities. This guide covers everything from basic variable configuration to advanced features like variable dependencies, encryption, and integration with other systems.

**Status: Production Ready** - The variables system is fully implemented with comprehensive variable management, resolution, and validation capabilities.

## Related Documentation

- [Actions User Guide](ACTIONS_USER_GUIDE.md) - Using variables in actions
- [Templates User Guide](TEMPLATES_USER_GUIDE.md) - Variable resolution in templates
- [Facts User Guide](FACTS_USER_GUIDE.md) - Using facts as variables
- [Secrets User Guide](SECRETS_USER_GUIDE.md) - Encrypted variable management
- [Machines User Guide](MACHINES_USER_GUIDE.md) - Machine-specific variables

> **See also**: [User Guides Index](USER_GUIDES_INDEX.md) - Complete overview of all user guides

## Getting Started

### Prerequisites

- spooky CLI installed and configured
- Basic understanding of HCL configuration syntax
- Access to create and modify project files

### Quick Start

1. **Check Available Variables Commands**
   ```bash
   spooky variables --help
   ```

2. **List Variables in a Project**
   ```bash
   spooky variables list ./my-project
   ```

3. **Validate Variable Configuration**
   ```bash
   spooky variables validate ./my-project
   ```

4. **Resolve Variables**
   ```bash
   spooky variables resolve ./my-project
   ```

## Variables System Concepts

### What are Variables?

Variables are configuration values that can be used across your spooky project. They provide:

- **Dynamic Configuration**: Values that change based on environment or context
- **Reusability**: Define once, use everywhere
- **Security**: Encrypt sensitive values with age encryption
- **Integration**: Work with facts, machines, and other systems

### Variable Sources

Variables can be defined from multiple sources:

1. **Project Variables**: Defined in `variables.hcl` and `variables/` directory
2. **Environment Variables**: System environment variables
3. **Facts Variables**: Derived from machine facts
4. **Default Variables**: System-provided defaults

### Variable Resolution

Variables follow a specific resolution order:

1. **Environment Variables** (highest priority)
2. **Project Variables** (from variables.hcl and variables/*.hcl)
3. **Facts Variables** (from machine facts)
4. **Default Variables** (lowest priority)

### Integration with Other Systems

Variables integrate with other spooky systems:

- **Actions**: Use variables in [action definitions](ACTIONS_USER_GUIDE.md)
- **Templates**: Resolve variables in [template rendering](TEMPLATES_USER_GUIDE.md)
- **Facts**: Reference [machine facts](FACTS_USER_GUIDE.md) as variables
- **Secrets**: Use [encrypted variables](SECRETS_USER_GUIDE.md) for sensitive data

## Current Implementation Status

### ✅ Working Features

- **Variable Loading**: Comprehensive variable loading from HCL configuration files
- **Variable Validation**: Complete variable validation and error handling
- **Variable Resolution**: Full variable resolution with dependency management
- **CLI Integration**: Complete CLI integration with all variable commands
- **Project Integration**: Full integration with project configuration
- **Encryption Support**: Age encryption support for sensitive variables

### 🔧 Advanced Features

- **Variable Dependencies**: Support for variable dependencies and references
- **Type Validation**: Comprehensive type validation for all variable types
- **Default Values**: Support for default values and fallbacks
- **Environment Integration**: Full integration with system environment variables
- **Facts Integration**: Integration with machine facts for dynamic values

## Basic Usage

### Listing Variables

List all variables in a project:

```bash
# List all variables
spooky variables list ./my-project
```

**Example Output**:
```
Variables in project (5 found):
1. app_name (string) - Application name
2. app_version (string) - Application version
3. environment (string) - Deployment environment
4. database_host (string) - Database hostname
5. api_port (number) - API server port

Total: 5 variables
```

### Validating Variables

Validate variable configurations:

```bash
# Validate all variables
spooky variables validate ./my-project
```

**Example Output**:
```
✅ Variable validation passed

Validated 5 variables:
- app_name: ✅ Valid
- app_version: ✅ Valid
- environment: ✅ Valid
- database_host: ✅ Valid
- api_port: ✅ Valid
```

### Resolving Variables

Resolve variables with their final values:

```bash
# Resolve all variables
spooky variables resolve ./my-project
```

**Example Output**:
```
Variable resolution for project: ./my-project
Total variables: 5
Resolved variables: 5
Resolution time: 2.3ms

Resolved values:
- app_name: "my-application"
- app_version: "1.0.0"
- environment: "production"
- database_host: "db.example.com"
- api_port: 8080
```

## Configuration

### Basic Variable Configuration

Define variables in your `variables.hcl` file:

```hcl
variables {
  variable "app_name" {
    value = "my-application"
    description = "Application name"
  }
  
  variable "app_version" {
    value = "1.0.0"
    description = "Application version"
  }
  
  variable "environment" {
    value = "production"
    description = "Deployment environment"
  }
  
  variable "database_host" {
    value = "db.example.com"
    description = "Database hostname"
  }
  
  variable "api_port" {
    value = 8080
    description = "API server port"
  }
}
```

### Advanced Variable Configuration

Use advanced variable features:

```hcl
variables {
  # String variable
  variable "app_name" {
    value = "my-application"
    description = "Application name"
  }
  
  # Number variable
  variable "api_port" {
    value = 8080
    description = "API server port"
  }
  
  # Boolean variable
  variable "debug_mode" {
    value = false
    description = "Enable debug mode"
  }
  
  # List variable
  variable "allowed_ips" {
    value = ["192.168.1.0/24", "10.0.0.0/8"]
    description = "Allowed IP ranges"
  }
  
  # Map variable
  variable "database_config" {
    value = {
      host = "db.example.com"
      port = 5432
      name = "myapp"
    }
    description = "Database configuration"
  }
}
```

### Environment Variable Integration

Use environment variables:

```hcl
variables {
  variable "user_home" {
    value = "{{.env.HOME}}"
    description = "User home directory"
  }
  
  variable "node_env" {
    value = "{{.env.NODE_ENV}}"
    description = "Node.js environment"
  }
  
  variable "api_key" {
    value = "{{.env.API_KEY}}"
    required = true
    description = "API key from environment"
  }
}
```

### Machine Integration

Use machine information in variables:

```hcl
variables {
  variable "web_servers" {
    value = "{{.machines.tags.role=web.hostname}}"
    description = "List of web server hostnames"
  }
  
  variable "db_servers" {
    value = "{{.machines.tags.role=database.hostname}}"
    description = "List of database server hostnames"
  }
  
  variable "production_machines" {
    value = "{{.machines.tags.environment=production.hostname}}"
    description = "Production machine hostnames"
  }
}
```

### Variable Encryption

Encrypt sensitive variables using the `spooky variables armor` command:

```bash
# Encrypt a variable value
echo "my-secret-password" | spooky variables armor --recipient age1... > encrypted_value.txt

# Use the encrypted value in variables.hcl
variable "database_password" {
  value = "age1..."  # Content from encrypted_value.txt
  encrypted = true
  description = "Database password"
}
```

### Decrypting Variables

Decrypt variables during action running:

```bash
# Run actions with decryption
spooky actions run ./my-project --decrypt
```

**Example Output**:
```
🔓 Decryption mode enabled - encrypted variables and facts will be decrypted in-memory
🔑 Using identity file: ~/.config/spooky/keys/identity.txt
✅ Decrypted 3 variables in-memory
```

## Troubleshooting

### Common Issues

#### Variable Resolution Failures

**Problem**: Variables fail to resolve

**Solutions**:
1. Check variable definitions:
   ```bash
   spooky variables validate ./my-project
   ```

2. Check for missing dependencies:
   ```bash
   spooky variables resolve ./my-project --verbose
   ```

3. Verify environment variables:
   ```bash
   spooky variables resolve ./my-project --show-env
   ```

#### Validation Errors

**Problem**: Variable validation fails

**Solutions**:
1. Check variable types:
   ```hcl
   variable "port" {
     value = 8080  # Must be number, not string
     description = "API port"
   }
   ```

2. Check required fields:
   ```hcl
   variable "api_key" {
     value = "{{.env.API_KEY}}"
     required = true  # Will fail if env var is missing
     description = "API key"
   }
   ```

3. Check validation rules:
   ```hcl
   variable "port" {
     value = 8080
     validation {
       min = 1024
       max = 65535
     }
     description = "API port"
   }
   ```

#### Encryption Issues

**Problem**: Encrypted variables fail to decrypt

**Solutions**:
1. Check age key configuration:
   ```bash
   spooky secrets validate
   ```

2. Verify recipient configuration:
   ```bash
   spooky variables armor ./my-project --dry-run
   ```

3. Check key permissions:
   ```bash
   ls -la ~/.config/spooky/keys/
   ```

## Integration with Other Systems

### Actions Integration

Variables integrate with the actions system:

```hcl
actions {
  action "deploy-application" {
    description = "Deploy application with variables"
    
    machines = ["web-server"]
    parallel = true
    
    variables {
      app_name = "{{.variables.app_name}}"
      app_version = "{{.variables.app_version}}"
      environment = "{{.variables.environment}}"
    }
    
    command = "deploy.sh"
  }
}
```

### Templates Integration

Variables integrate with the templates system:

```hcl
templates {
  template "config.tmpl" {
    source = "templates/config.tmpl"
    destination = "/etc/app/config.conf"
    
    variables {
      app_name = "{{.variables.app_name}}"
      api_port = "{{.variables.api_port}}"
      database_host = "{{.variables.database_host}}"
    }
  }
}
```

### Facts Integration

Variables can reference machine facts:

```hcl
variables {
  variable "machine_count" {
    value = "{{.facts.machines.count}}"
    description = "Number of machines in inventory"
  }
  
  variable "os_distribution" {
    value = "{{.facts.machines.os.distribution}}"
    description = "Most common OS distribution"
  }
}
```

## Examples

### Basic Project Variables

```hcl
# variables.hcl
variables {
  variable "app_name" {
    value = "my-application"
    description = "Application name"
  }
  
  variable "app_version" {
    value = "1.0.0"
    description = "Application version"
  }
  
  variable "environment" {
    value = "production"
    description = "Deployment environment"
  }
}
```

### Advanced Project Variables

```hcl
# variables.hcl
variables {
  # Application configuration
  variable "app_name" {
    value = "my-application"
    description = "Application name"
    required = true
  }
  
  variable "app_version" {
    value = "1.0.0"
    description = "Application version"
    required = true
  }
  
  # Environment configuration
  variable "environment" {
    value = "{{.env.NODE_ENV}}"
    default = "development"
    description = "Deployment environment"
  }
  
  # Database configuration
  variable "database_config" {
    value = {
      host = "{{.env.DB_HOST}}"
      port = 5432
      name = "{{.variables.app_name}}"
    }
    description = "Database configuration"
  }
  
  # Network configuration
  variable "api_port" {
    value = 8080
    description = "API server port"
    validation {
      min = 1024
      max = 65535
    }
  }
  
  variable "allowed_ips" {
    value = ["192.168.1.0/24", "10.0.0.0/8"]
    description = "Allowed IP ranges"
  }
  
  # Encrypted secrets
  variable "database_password" {
    value = "age1..."
    encrypted = true
    description = "Database password"
  }
}
```

### Modular Variable Organization

```hcl
# variables/app.hcl
variables {
  variable "app_name" {
    value = "my-application"
    description = "Application name"
  }
  
  variable "app_version" {
    value = "1.0.0"
    description = "Application version"
  }
}

# variables/database.hcl
variables {
  variable "database_host" {
    value = "db.example.com"
    description = "Database hostname"
  }
  
  variable "database_port" {
    value = 5432
    description = "Database port"
  }
}

# variables/network.hcl
variables {
  variable "api_port" {
    value = 8080
    description = "API server port"
  }
  
  variable "allowed_ips" {
    value = ["192.168.1.0/24", "10.0.0.0/8"]
    description = "Allowed IP ranges"
  }
}
```

## Best Practices

### Variable Organization

- Use descriptive variable names
- Group related variables together
- Use modular organization with variables/ directory
- Include comprehensive descriptions

### Variable Security

- Encrypt sensitive variables with age encryption
- Use environment variables for secrets when possible
- Never commit unencrypted secrets to version control
- Use appropriate validation rules

### Variable Dependencies

- Minimize variable dependencies
- Use clear dependency chains
- Avoid circular dependencies
- Document complex dependencies

### Variable Validation

- Use appropriate types for variables
- Add validation rules where needed
- Use required fields for critical variables
- Provide sensible defaults

## Next Steps

- Explore the [Variables API Reference](VARIABLES_API_REFERENCE.md) for detailed technical information
- Check the [Variables Troubleshooting Guide](VARIABLES_TROUBLESHOOTING.md) for common issues
- Review the [Variables Documentation Summary](VARIABLES_DOCUMENTATION_SUMMARY.md) for implementation details
- Learn about [Variables Integration Patterns](INTEGRATIONS_USER_GUIDE.md) for advanced usage
