# Variables System Troubleshooting Guide

## Overview

This troubleshooting guide provides solutions for common issues encountered when working with the spooky variables system. It covers error messages, variable resolution problems, validation issues, and performance problems.

## Table of Contents

1. [Common Error Messages](#common-error-messages)
2. [Variable Resolution Issues](#variable-resolution-issues)
3. [Validation Problems](#validation-problems)
4. [Performance Issues](#performance-issues)
5. [Configuration Problems](#configuration-problems)
6. [Debugging Techniques](#debugging-techniques)
7. [Best Practices for Troubleshooting](#best-practices-for-troubleshooting)

## Common Error Messages

### Loading Errors

#### "Failed to load variables: no variables found in project"

**Cause:** No variable configuration files found in the project directory.

**Solution:**
```bash
# Check if variables.hcl exists
ls -la ./my-project/variables.hcl

# Check if variables/ directory exists
ls -la ./my-project/variables/

# Create a basic variables.hcl file
cat > ./my-project/variables.hcl << 'EOF'
variables {
  variable "app_name" {
    type = "string"
    description = "Application name"
    default = "my-app"
    scope = "project"
  }
}
EOF
```

#### "Failed to parse variable block: Unexpected block type"

**Cause:** Invalid HCL syntax in variable configuration.

**Solution:**
```bash
# Check HCL syntax
spooky variables validate ./my-project

# Common syntax errors to fix:
# 1. Missing quotes around strings
# 2. Incorrect block structure
# 3. Invalid attribute names

# Example of correct syntax:
variables {
  variable "app_name" {
    type = "string"  # Use quotes for strings
    description = "Application name"
    default = "my-app"
    scope = "project"
  }
}
```

#### "duplicate variable name 'app_name' found in multiple files"

**Cause:** Same variable name defined in multiple configuration files.

**Solution:**
```bash
# Find all occurrences of the duplicate variable name
grep -r "app_name" ./my-project/variables/

# Rename one of the variables to make it unique
# In variables.hcl:
variable "app_name" {
  type = "string"
  default = "my-app"
}

# In variables/app.hcl:
variable "app_name_alt" {  # Changed from "app_name"
  type = "string"
  default = "my-app-alt"
}
```

### Validation Errors

#### "invalid variable name 'app-name': must match pattern"

**Cause:** Variable name contains invalid characters.

**Solution:**
```bash
# Variable names must:
# - Start with letter or underscore
# - Contain only letters, numbers, and underscores
# - Not contain hyphens or other special characters

# ❌ WRONG
variable "app-name" {
  type = "string"
}

# ✅ CORRECT
variable "app_name" {
  type = "string"
}

# ✅ CORRECT
variable "appName" {
  type = "string"
}
```

#### "unknown variable type: invalid_type"

**Cause:** Invalid variable type specified.

**Solution:**
```bash
# Valid variable types are:
# - string
# - number
# - bool
# - list
# - map

# ❌ WRONG
variable "port" {
  type = "integer"  # Invalid type
  default = 8080
}

# ✅ CORRECT
variable "port" {
  type = "number"  # Valid type
  default = 8080
}
```

#### "invalid variable scope: invalid_scope"

**Cause:** Invalid variable scope specified.

**Solution:**
```bash
# Valid variable scopes are:
# - project
# - global
# - local

# ❌ WRONG
variable "config" {
  type = "string"
  scope = "environment"  # Invalid scope
}

# ✅ CORRECT
variable "config" {
  type = "string"
  scope = "project"  # Valid scope
}
```

## Variable Resolution Issues

### Circular Dependency Errors

#### "circular dependency detected: app_name -> app_title -> app_name"

**Cause:** Variables have circular dependencies that cannot be resolved.

**Solution:**
```bash
# ❌ PROBLEMATIC - Circular dependency
variable "app_name" {
  type = "string"
  default = "my-app"
  dependencies = ["app_title"]
}

variable "app_title" {
  type = "string"
  default = "My Application"
  dependencies = ["app_name"]
}

# ✅ SOLUTION 1 - Remove dependency
variable "app_name" {
  type = "string"
  default = "my-app"
  # No dependencies
}

variable "app_title" {
  type = "string"
  default = "My Application"
  dependencies = ["app_name"]  # Only one-way dependency
}

# ✅ SOLUTION 2 - Use computed value
variable "app_name" {
  type = "string"
  default = "my-app"
}

variable "app_title" {
  type = "string"
  default = "My Application"
  # Compute title based on name without dependency
}
```

### Missing Dependency Errors

#### "variable 'app_url' depends on undefined variable 'base_url'"

**Cause:** Variable depends on another variable that doesn't exist.

**Solution:**
```bash
# ❌ PROBLEMATIC - Missing dependency
variable "app_url" {
  type = "string"
  dependencies = ["base_url"]  # base_url doesn't exist
}

# ✅ SOLUTION - Define the missing variable
variable "base_url" {
  type = "string"
  default = "https://api.example.com"
}

variable "app_url" {
  type = "string"
  dependencies = ["base_url"]
}
```

### Environment Variable Override Issues

#### "Environment variable override failed for variable 'db_password'"

**Cause:** Environment variable override configuration is incorrect.

**Solution:**
```bash
# Environment variables must use the SPOOKY_VAR_ prefix
# Variable names are converted to uppercase with underscores

# ❌ WRONG
export DB_PASSWORD="secret123"

# ✅ CORRECT
export SPOOKY_VAR_DB_PASSWORD="secret123"

# For variables with underscores, use underscores in env var
variable "api_key" {
  type = "string"
  sensitive = true
}

# Environment variable: SPOOKY_VAR_API_KEY
export SPOOKY_VAR_API_KEY="your-api-key"
```

## Validation Problems

### Constraint Validation Errors

#### "constraint validation failed: value 80 is less than minimum 1024"

**Cause:** Variable value violates defined constraints.

**Solution:**
```bash
# ❌ PROBLEMATIC - Value violates constraint
variable "port" {
  type = "number"
  default = 80  # Too low
  constraints {
    min_value = 1024
    max_value = 65535
  }
}

# ✅ SOLUTION - Use valid value
variable "port" {
  type = "number"
  default = 8080  # Valid value
  constraints {
    min_value = 1024
    max_value = 65535
  }
}
```

### Type Validation Errors

#### "type validation failed: expected string, got number"

**Cause:** Variable value doesn't match declared type.

**Solution:**
```bash
# ❌ PROBLEMATIC - Type mismatch
variable "port" {
  type = "string"
  default = 8080  # Number instead of string
}

# ✅ SOLUTION - Match type and value
variable "port" {
  type = "string"
  default = "8080"  # String value
}

# OR use correct type
variable "port" {
  type = "number"
  default = 8080  # Number value
}
```

### Pattern Validation Errors

#### "pattern validation failed: string does not match pattern"

**Cause:** String value doesn't match required pattern.

**Solution:**
```bash
# ❌ PROBLEMATIC - Pattern mismatch
variable "email" {
  type = "string"
  default = "invalid-email"
  validation {
    pattern = "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"
  }
}

# ✅ SOLUTION - Use valid pattern
variable "email" {
  type = "string"
  default = "user@example.com"  # Valid email
  validation {
    pattern = "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"
  }
}
```

## Performance Issues

### Slow Variable Loading

#### "Variable loading is taking too long"

**Cause:** Large number of variable files or complex dependencies.

**Solution:**
```bash
# 1. Optimize file organization
# Instead of many small files, use fewer larger files
# Group related variables together

# 2. Reduce dependency complexity
# Minimize cross-file dependencies
# Use simpler dependency graphs

# 3. Use caching
# Enable variable caching for repeated operations
spooky variables resolve --cache-enabled

# 4. Profile performance
spooky variables resolve --profile
```

### Memory Usage Issues

#### "High memory usage during variable resolution"

**Cause:** Large variable sets or complex dependency graphs.

**Solution:**
```bash
# 1. Limit variable scope
# Use local scope for temporary variables
# Avoid global scope for large variable sets

# 2. Optimize data structures
# Use appropriate variable types
# Avoid large default values

# 3. Stream processing
# Process variables in batches
spooky variables resolve --batch-size 100

# 4. Monitor memory usage
spooky variables resolve --memory-profile
```

## Configuration Problems

### HCL Parsing Issues

#### "HCL parsing failed: unexpected token"

**Cause:** Invalid HCL syntax in configuration files.

**Solution:**
```bash
# Common HCL syntax issues:

# 1. Missing quotes around strings
# ❌ WRONG
variable "name" {
  type = string
}

# ✅ CORRECT
variable "name" {
  type = "string"
}

# 2. Incorrect block structure
# ❌ WRONG
variables {
  variable "name" {
    type = "string"
  }
  variable "port" {
    type = "number"
  }
}

# ✅ CORRECT
variables {
  variable "name" {
    type = "string"
  }
}

variables {
  variable "port" {
    type = "number"
  }
}

# 3. Invalid attribute names
# ❌ WRONG
variable "name" {
  Type = "string"  # Capitalized
}

# ✅ CORRECT
variable "name" {
  type = "string"  # Lowercase
}
```

### File Organization Issues

#### "Conflicting variable definitions across files"

**Cause:** Poor file organization leading to conflicts.

**Solution:**
```bash
# 1. Use clear naming conventions
# variables/app.hcl - Application variables
# variables/database.hcl - Database variables
# variables/security.hcl - Security variables

# 2. Avoid overlapping variable names
# Use prefixes to avoid conflicts
# app_name, db_name, api_name

# 3. Use hierarchical organization
variables/
├── app/
│   ├── web.hcl
│   └── api.hcl
├── database/
│   ├── postgres.hcl
│   └── redis.hcl
└── security/
    ├── auth.hcl
    └── encryption.hcl
```

## Debugging Techniques

### Verbose Output

Use verbose output to get detailed information about variable processing:

```bash
# Verbose validation
spooky variables validate --verbose

# Verbose resolution
spooky variables resolve --verbose

# Verbose loading
spooky variables list --verbose
```

### Dependency Analysis

Analyze variable dependencies to understand resolution order:

```bash
# Show dependency graph
spooky variables resolve --show-dependencies

# Show resolution order
spooky variables resolve --show-order

# Show dependency conflicts
spooky variables validate --show-dependencies
```

### Environment Variable Debugging

Debug environment variable overrides:

```bash
# Show which environment variables are being used
spooky variables resolve --debug-env

# Show environment variable mapping
spooky variables resolve --env-mapping

# Test environment variable overrides
export SPOOKY_VAR_DEBUG=true
spooky variables resolve --env-override
```

### File-Level Debugging

Debug specific file issues:

```bash
# Validate specific file
spooky variables validate --file variables/app.hcl

# Show file parsing details
spooky variables validate --parse-details

# Show file dependencies
spooky variables validate --file-dependencies
```

### Performance Profiling

Profile variable system performance:

```bash
# Profile loading performance
spooky variables list --profile

# Profile validation performance
spooky variables validate --profile

# Profile resolution performance
spooky variables resolve --profile

# Memory profiling
spooky variables resolve --memory-profile
```

## Best Practices for Troubleshooting

### Systematic Approach

1. **Identify the Problem**
   - Read error messages carefully
   - Note the specific variable and file involved
   - Understand the context (loading, validation, resolution)

2. **Isolate the Issue**
   - Test with minimal configuration
   - Remove variables one by one to identify the problematic one
   - Check file syntax independently

3. **Check Dependencies**
   - Verify all required variables are defined
   - Check for circular dependencies
   - Validate dependency order

4. **Validate Configuration**
   - Use `spooky variables validate` to check syntax
   - Verify variable types and constraints
   - Check for duplicate definitions

5. **Test Resolution**
   - Use `spooky variables resolve` to test resolution
   - Check environment variable overrides
   - Verify final values

### Common Debugging Commands

```bash
# Basic validation
spooky variables validate ./my-project

# Detailed validation
spooky variables validate ./my-project --verbose

# List variables to see what's loaded
spooky variables list ./my-project

# Test resolution
spooky variables resolve ./my-project

# Debug environment variables
spooky variables resolve ./my-project --debug-env

# Profile performance
spooky variables resolve ./my-project --profile
```

### Error Message Interpretation

#### Understanding Error Context

```bash
# Error: "variable 'app_url' depends on undefined variable 'base_url'"
# This means:
# 1. Variable 'app_url' has a dependency on 'base_url'
# 2. Variable 'base_url' is not defined anywhere
# 3. Solution: Define 'base_url' or remove the dependency

# Error: "circular dependency detected: A -> B -> C -> A"
# This means:
# 1. Variable A depends on B
# 2. Variable B depends on C
# 3. Variable C depends on A
# 4. This creates a cycle that cannot be resolved
# 5. Solution: Break the cycle by removing one dependency
```

#### Error Severity Levels

- **Fatal Errors**: Must be fixed before variables can be used
  - Syntax errors
  - Missing required variables
  - Circular dependencies

- **Validation Errors**: Prevent proper variable resolution
  - Type mismatches
  - Constraint violations
  - Invalid patterns

- **Warnings**: May cause issues but don't prevent operation
  - Deprecated syntax
  - Unused variables
  - Performance issues

### Getting Help

When troubleshooting becomes difficult:

1. **Check Documentation**
   - Review the [User Guide](VARIABLES_USER_GUIDE.md)
   - Consult the [API Reference](VARIABLES_API_REFERENCE.md)
   - Look at examples in the examples directory

2. **Use Debugging Tools**
   - Enable verbose output
   - Use profiling tools
   - Check system logs

3. **Simplify the Problem**
   - Create minimal reproduction case
   - Test with basic configuration
   - Remove complexity until it works

4. **Check for Known Issues**
   - Review recent changes
   - Check for similar problems
   - Consult community resources

This troubleshooting guide should help you resolve most issues with the spooky variables system. If you continue to experience problems, refer to the [User Guide](VARIABLES_USER_GUIDE.md) for more detailed information or the [API Reference](VARIABLES_API_REFERENCE.md) for technical implementation details.
