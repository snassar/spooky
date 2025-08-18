# Variables System Documentation Summary

## Overview

This document provides a comprehensive overview of the spooky variables system documentation. It serves as a guide to help you find the right documentation for your needs and understand how all the pieces fit together.

**Status: Implemented** - The variables system is fully implemented with comprehensive functionality for variable management, encryption, and integration with other systems.

## Documentation Structure

### 📚 Core Documentation

#### 1. [User Guide](VARIABLES_USER_GUIDE.md)
**Audience:** End users, system administrators, DevOps engineers
**Purpose:** Complete guide to using the variables system

**What it covers:**
- Getting started with variable configuration
- Variable definition and management
- Encryption and security features
- Integration with actions and templates
- Real-world examples and use cases

**When to use:** Start here if you're new to spooky variables or need to understand how to use the system effectively.

#### 2. [API Reference](VARIABLES_API_REFERENCE.md)
**Audience:** Developers, system integrators, contributors
**Purpose:** Technical reference for the variables system APIs and implementation

**What it covers:**
- Core interfaces and type definitions
- Implementation details and algorithms
- Error handling patterns
- Configuration rules and schemas
- CLI integration details
- Code examples and patterns

**When to use:** Use this when developing with the variables system, extending functionality, or debugging implementation issues.

#### 3. [Troubleshooting Guide](VARIABLES_TROUBLESHOOTING.md)
**Audience:** System administrators, support engineers, users experiencing issues
**Purpose:** Solutions for common problems and debugging techniques

**What it covers:**
- Common error messages and solutions
- Encryption and decryption issues
- Variable resolution problems
- Integration issues with other systems
- Configuration problems and debugging
- Best practices for troubleshooting

**When to use:** Use this when encountering problems or need to debug issues with the variables system.

### 📁 Examples Directory

#### [Examples Overview](examples/README.md)
**Audience:** All users
**Purpose:** Quick reference for available examples and use cases

**What it covers:**
- Available variable configuration examples
- Example configurations and scripts
- Common use case patterns
- Integration examples with other systems

**When to use:** Use this to quickly find relevant examples for your use case.

## Key Concepts

### Core Features

1. **Variable Definition** - Define variables in HCL format with support for multiple data types
2. **Variable Encryption** - Encrypt sensitive variables using age encryption
3. **Variable Resolution** - Resolve variables with support for dependencies and references
4. **Variable Validation** - Validate variable definitions and values
5. **Integration Support** - Seamless integration with actions, templates, and other systems
6. **CLI Management** - Comprehensive CLI commands for variable management
7. **Security Features** - Secure handling of sensitive variable data

### Architecture Principles

1. **Interface-First Design** - All functionality through well-defined interfaces
2. **Dependency Injection** - Loose coupling through interface-based dependencies
3. **Security by Default** - Encryption and secure handling of sensitive data
4. **Extensible Design** - Easy to add new variable types and features
5. **Performance Optimized** - Efficient variable resolution and caching

### Best Practices

1. **Use Descriptive Names** - Use clear, descriptive variable names
2. **Encrypt Sensitive Data** - Encrypt passwords, keys, and other sensitive information
3. **Organize Variables** - Use logical organization and grouping
4. **Validate Variables** - Always validate variable definitions
5. **Use Appropriate Types** - Choose appropriate data types for your variables
6. **Document Variables** - Include descriptions and usage information

## Variables System Overview

### Core Concepts

The variables system provides a comprehensive solution for managing configuration variables in spooky projects. Variables can be:

- **Simple values** - Strings, numbers, booleans
- **Complex structures** - Maps, lists, nested objects
- **Encrypted values** - Sensitive data encrypted with age
- **Referenced values** - Variables that reference other variables
- **Computed values** - Values computed from other variables or external sources

### Variable Definition Structure

Variables are defined in `variables.hcl` files or in the `variables/` directory:

```hcl
variables {
  # Simple string variable
  variable "app_name" {
    value = "my-application"
    description = "Application name"
  }
  
  # Number variable
  variable "port" {
    value = 8080
    description = "Application port"
  }
  
  # Boolean variable
  variable "debug_mode" {
    value = false
    description = "Enable debug mode"
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
  
  # List variable
  variable "allowed_ips" {
    value = ["192.168.1.0/24", "10.0.0.0/8"]
    description = "Allowed IP ranges"
  }
  
  # Encrypted variable
  variable "db_password" {
    value = "age1ql3z7hjy54pw3hyww5ayyfg7zqgvc7w3j2elw8zmrj2kg5sfn9aqmcac8p"
    encrypted = true
    description = "Database password (encrypted)"
  }
  
  # Referenced variable
  variable "db_url" {
    value = "postgresql://user:${db_password}@${database_config.host}:${database_config.port}/${database_config.name}"
    description = "Database connection URL"
  }
}
```

### CLI Commands

The variables system provides comprehensive CLI commands:

```bash
# List all variables in a project
spooky variables list ./my-project

# List variables with filtering
spooky variables list ./my-project --variable app_name

# Validate variable definitions
spooky variables validate ./my-project

# Validate with verbose output
spooky variables validate ./my-project --verbose

# Armor (encrypt) a variable
spooky variables armor ./my-project --variable db_password --value "secret123"

# Decrypt variables during action execution
spooky actions run ./my-project --decrypt

# Decrypt with dry-run
spooky actions run ./my-project --decrypt --dry-run
```

### Variable Types

The variables system supports multiple data types:

#### Basic Types
```hcl
# String
variable "app_name" {
  value = "my-application"
}

# Number
variable "port" {
  value = 8080
}

# Boolean
variable "debug_mode" {
  value = false
}
```

#### Complex Types
```hcl
# Map
variable "config" {
  value = {
    host = "example.com"
    port = 443
    ssl = true
  }
}

# List
variable "servers" {
  value = ["server1", "server2", "server3"]
}

# Nested structures
variable "database" {
  value = {
    primary = {
      host = "db1.example.com"
      port = 5432
    }
    replica = {
      host = "db2.example.com"
      port = 5432
    }
  }
}
```

### Encryption Support

The variables system supports encryption using age:

```hcl
# Encrypted variable
variable "secret_key" {
  value = "age1ql3z7hjy54pw3hyww5ayyfg7zqgvc7w3j2elw8zmrj2kg5sfn9aqmcac8p"
  encrypted = true
  description = "Secret API key"
}
```

#### Encryption Commands
```bash
# Encrypt a variable value
spooky variables armor ./my-project --variable secret_key --value "my-secret-value"

# Decrypt during action execution
spooky actions run ./my-project --decrypt
```

### Variable Resolution

Variables support complex resolution patterns:

```hcl
# Basic reference
variable "app_name" {
  value = "my-application"
}

variable "app_version" {
  value = "1.0.0"
}

# Reference other variables
variable "full_name" {
  value = "${app_name}-${app_version}"
}

# Nested references
variable "database" {
  value = {
    host = "db.example.com"
    port = 5432
  }
}

variable "db_url" {
  value = "postgresql://user:pass@${database.host}:${database.port}/app"
}
```

### Multi-File Support

Variables can be organized across multiple files:

```bash
my-project/
├── variables.hcl          # Main variables file
└── variables/
    ├── database.hcl       # Database-related variables
    ├── network.hcl        # Network configuration
    └── secrets.hcl        # Encrypted secrets
```

## Implementation Details

### Core Components

1. **Variable Loader** - Loads variables from HCL files
2. **Variable Validator** - Validates variable definitions and values
3. **Variable Resolver** - Resolves variable references and dependencies
4. **Encryption Manager** - Handles variable encryption and decryption
5. **Variable Integration** - Provides integration with other system components

### Integration Points

The variables system integrates with:

- **Actions System** - For variable injection during action execution
- **Templates System** - For variable substitution in templates
- **Secrets System** - For encryption and decryption operations
- **CLI System** - For user interface and command execution

### Error Handling

The variables system provides comprehensive error handling:

- **Validation errors** - Invalid variable definitions or values
- **Resolution errors** - Circular dependencies or missing references
- **Encryption errors** - Encryption/decryption failures
- **Type errors** - Type mismatches and conversion issues
- **File errors** - File I/O and parsing issues

## Best Practices

### Variable Organization

1. **Use logical grouping** to organize related variables
2. **Separate sensitive data** into dedicated files
3. **Use descriptive names** for easy identification
4. **Include documentation** with descriptions
5. **Validate configurations** before use

### Security

1. **Encrypt sensitive data** using age encryption
2. **Use secure variable names** that don't reveal sensitive information
3. **Limit access** to variable files
4. **Rotate encrypted values** regularly
5. **Use environment-specific** variable files

### Performance

1. **Minimize variable dependencies** to improve resolution speed
2. **Use appropriate data types** for your use case
3. **Cache resolved values** when possible
4. **Validate early** to catch issues quickly
5. **Use efficient resolution** patterns

## Troubleshooting

### Common Issues

1. **Variable resolution errors** - Check for circular dependencies and missing references
2. **Encryption errors** - Verify age keys and encryption configuration
3. **Type errors** - Ensure variable types match expected values
4. **File parsing errors** - Validate HCL syntax and file permissions
5. **Integration errors** - Check variable injection in actions and templates

### Debug Commands

```bash
# Enable verbose logging
export SPOOKY_LOG_LEVEL=debug

# List variables with details
spooky variables list ./my-project --verbose

# Validate with detailed output
spooky variables validate ./my-project --verbose

# Test variable resolution
spooky variables resolve ./my-project --variable db_url

# Check encryption status
spooky variables list ./my-project --encrypted
```

### Common Patterns

1. **Environment-specific variables** - Use different files for different environments
2. **Sensitive data handling** - Encrypt all sensitive information
3. **Variable composition** - Build complex values from simple components
4. **Default values** - Provide sensible defaults for optional variables
5. **Validation rules** - Use validation to ensure data quality

## Related Documentation

- [Variables User Guide](VARIABLES_USER_GUIDE.md) - Complete user guide
- [Variables API Reference](VARIABLES_API_REFERENCE.md) - Technical reference
- [Variables Troubleshooting](VARIABLES_TROUBLESHOOTING.md) - Troubleshooting guide
- [System Design](../design/systems/variables-system.md) - System design documentation
- [CLI Reference](CLI_REFERENCE.md) - CLI command reference
