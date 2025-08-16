# Variables System Troubleshooting Guide

## Overview

This troubleshooting guide provides solutions for common issues encountered when working with the spooky variables system. It covers error messages, variable resolution problems, validation issues, and performance problems.

**Status: Production Ready** - The variables system is fully implemented with comprehensive variable management, resolution, validation, and encryption capabilities.

## Variables System Status

### ✅ Fully Functional Variables Infrastructure

The variables system now has **complete variables infrastructure** with:

- **Variable Loading**: Comprehensive variable loading from HCL files and directories
- **Variable Resolution**: Complete variable resolution with dependency management
- **Variable Validation**: Comprehensive variable validation and error handling
- **CLI Integration**: Full CLI integration with `spooky variables` commands
- **Project Integration**: Variable loading from project configuration
- **Encryption Support**: Variable encryption and decryption capabilities
- **Export Functionality**: Variable export to JSON and HCL formats
- **Error Handling**: Comprehensive error handling and reporting

### What This Means for Users

- **No More Stubs**: All functionality is fully implemented - no placeholder code
- **Production Ready**: The system is ready for production use
- **Complete Feature Set**: All documented features are functional
- **Reliable Resolution**: Robust variable resolution and dependency management
- **Performance Optimized**: Efficient variable loading and resolution

### Expected Behavior

When using variables, you can expect:

1. **Proper Variable Loading**: Variables load correctly from HCL files and directories
2. **Dependency Resolution**: Variable dependencies are resolved correctly
3. **Validation**: Comprehensive variable validation and error reporting
4. **Encryption**: Variable encryption and decryption functionality
5. **Export Functionality**: Variable export to various formats
6. **Error Handling**: Clear error messages with actionable information

## Common Issues and Solutions

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
```hcl
# ✅ CORRECT - Valid variable name
variables {
  variable "app_name" {  # Use underscores, not hyphens
    type = "string"
    default = "my-app"
  }
}

# ❌ WRONG - Invalid variable name
variables {
  variable "app-name" {  # Hyphens not allowed
    type = "string"
    default = "my-app"
  }
}
```

#### "Validation failed: missing required field 'type'"

**Cause:** Variable configuration is missing required fields.

**Solution:**
```hcl
# ✅ CORRECT - All required fields present
variables {
  variable "app_name" {
    type = "string"        # Required
    description = "App name"  # Optional but recommended
    default = "my-app"     # Optional
    scope = "project"      # Optional
  }
}
```

#### "Validation failed: invalid variable type"

**Cause:** Variable type is not supported.

**Solution:**
```hcl
# ✅ CORRECT - Valid variable types
variables {
  variable "app_name" {
    type = "string"        # Valid
    default = "my-app"
  }
  
  variable "app_port" {
    type = "number"        # Valid
    default = 8080
  }
  
  variable "app_enabled" {
    type = "bool"          # Valid
    default = true
  }
  
  variable "app_tags" {
    type = "list"          # Valid
    default = ["web", "production"]
  }
}
```

### Resolution Issues

#### "Circular dependency detected in variable resolution"

**Cause:** Variables have circular dependencies.

**Solution:**
```hcl
# ❌ WRONG - Circular dependency
variables {
  variable "app_name" {
    type = "string"
    default = "${app_display_name}"  # Depends on app_display_name
  }
  
  variable "app_display_name" {
    type = "string"
    default = "${app_name}"  # Depends on app_name - CIRCULAR!
  }
}

# ✅ CORRECT - No circular dependency
variables {
  variable "app_name" {
    type = "string"
    default = "my-app"
  }
  
  variable "app_display_name" {
    type = "string"
    default = "${app_name}"  # Depends on app_name - OK
  }
}
```

#### "Variable 'undefined_var' not found"

**Cause:** Variable is referenced but not defined.

**Solution:**
```hcl
# ✅ CORRECT - Define all referenced variables
variables {
  variable "app_name" {
    type = "string"
    default = "my-app"
  }
  
  variable "app_url" {
    type = "string"
    default = "https://${app_name}.example.com"  # app_name is defined
  }
}
```

### Encryption Issues

#### "Failed to decrypt variable: invalid encryption key"

**Cause:** Encryption key is incorrect or missing.

**Solution:**
```bash
# Check if encryption key is set
echo $SPOOKY_ENCRYPTION_KEY

# Set encryption key
export SPOOKY_ENCRYPTION_KEY="your-encryption-key"

# Or use age encryption
spooky variables armor ./my-project --key age1...
```

#### "Failed to encrypt variable: encryption not configured"

**Cause:** Encryption is not properly configured.

**Solution:**
```bash
# Configure age encryption
spooky variables armor ./my-project --key age1...

# Or set encryption key
export SPOOKY_ENCRYPTION_KEY="your-key"
```

### Performance Issues

#### "Variable resolution is very slow"

**Cause:** Too many variables or complex dependencies.

**Solution:**
```bash
# Use variable filtering
spooky variables list ./my-project --variables app_name,app_port

# Optimize variable dependencies
# Reduce circular dependencies
# Use simpler variable types
```

#### "High memory usage during variable operations"

**Cause:** Too many variables or inefficient operations.

**Solution:**
```bash
# Process fewer variables at once
spooky variables validate ./my-project --variables var1,var2

# Monitor memory usage
top -p $(pgrep spooky)
```

## Configuration Problems

### Project Configuration Issues

#### "Invalid variable configuration: unknown field"

**Cause:** Unknown fields in variable configuration.

**Solution:**
```hcl
# ✅ CORRECT - Valid variable configuration
variables {
  variable "app_name" {
    type = "string"
    description = "Application name"
    default = "my-app"
    scope = "project"
    
    validation {
      condition = length(var.app_name) > 0
      error_message = "App name cannot be empty"
    }
  }
}
```

#### "Invalid validation configuration: missing condition"

**Cause:** Validation configuration is incomplete.

**Solution:**
```hcl
# ✅ CORRECT - Complete validation configuration
variables {
  variable "app_port" {
    type = "number"
    default = 8080
    
    validation {
      condition = var.app_port > 0 && var.app_port <= 65535
      error_message = "Port must be between 1 and 65535"
    }
  }
}
```

### File Organization Issues

#### "Invalid variable file structure"

**Cause:** Variable files are not organized correctly.

**Solution:**
```bash
# ✅ CORRECT - Proper file structure
./my-project/
├── variables.hcl          # Main variables file
└── variables/             # Variables directory
    ├── app.hcl           # App-specific variables
    ├── database.hcl      # Database variables
    └── network.hcl       # Network variables

# ❌ WRONG - Invalid structure
./my-project/
├── variables.hcl
└── vars/                 # Wrong directory name
    └── app.hcl
```

## Debugging Techniques

### Enable Verbose Output

```bash
# Enable verbose output for variable operations
spooky variables list ./my-project --verbose
spooky variables validate ./my-project --verbose

# Enable debug logging
export SPOOKY_LOG_LEVEL=debug
spooky variables list ./my-project
```

### Test Variable Resolution

```bash
# Test variable resolution
spooky variables resolve ./my-project

# Test specific variables
spooky variables resolve ./my-project --variables app_name,app_port

# Export variables for inspection
spooky variables export ./my-project --output variables.json
```

### Validate Configuration

```bash
# Validate variable configuration
spooky variables validate ./my-project

# Check specific variables
spooky variables validate ./my-project --variables app_name

# Check for circular dependencies
spooky variables validate ./my-project --check-dependencies
```

## Recovery Procedures

### Variable Configuration Recovery

```bash
# Backup configuration
cp variables.hcl variables.hcl.backup

# Validate configuration
spooky variables validate ./my-project

# Restore from backup if needed
cp variables.hcl.backup variables.hcl
```

### Encryption Recovery

```bash
# Test encryption
spooky variables armor ./my-project --test

# Re-encrypt variables
spooky variables armor ./my-project --key age1...

# Decrypt variables
spooky variables unarmor ./my-project
```

## Prevention Strategies

### Regular Validation

```bash
# Schedule regular validation
crontab -e
# Add: 0 2 * * * /usr/local/bin/spooky variables validate /path/to/project

# Validate before operations
spooky variables validate ./my-project

# Validate in CI/CD pipeline
spooky variables validate ./my-project --strict
```

### Monitoring

```bash
# Monitor variable resolution
spooky variables resolve ./my-project

# Monitor encryption status
spooky variables armor ./my-project --status

# Monitor system resources
top -p $(pgrep spooky)
```

### Backup Strategy

```bash
# Backup variable configurations
cp variables.hcl variables.hcl.$(date +%Y%m%d)

# Version control configurations
git add variables.hcl
git commit -m "Update variable configuration"

# Backup project structure
tar -czf project-backup-$(date +%Y%m%d).tar.gz ./
```

## Best Practices for Troubleshooting

### 1. Start Simple

Begin with simple variable configurations and add complexity gradually:

```hcl
# Start with basic configuration
variables {
  variable "app_name" {
    type = "string"
    default = "my-app"
  }
}

# Then add complexity
variables {
  variable "app_name" {
    type = "string"
    description = "Application name"
    default = "my-app"
    scope = "project"
    
    validation {
      condition = length(var.app_name) > 0
      error_message = "App name cannot be empty"
    }
  }
}
```

### 2. Use Descriptive Names

Use clear, descriptive variable names:

```hcl
# ✅ GOOD - Descriptive names
variables {
  variable "production_database_url" {
    type = "string"
    default = "postgresql://user:pass@host:5432/db"
  }
}

# ❌ BAD - Unclear names
variables {
  variable "db_url" {
    type = "string"
    default = "postgresql://user:pass@host:5432/db"
  }
}
```

### 3. Validate Early and Often

Validate configurations frequently:

```bash
# Validate after every change
spooky variables validate ./my-project

# Validate before operations
spooky variables validate ./my-project && spooky variables resolve ./my-project

# Validate in scripts
#!/bin/bash
if spooky variables validate ./my-project; then
    spooky variables resolve ./my-project
else
    echo "Validation failed"
    exit 1
fi
```

### 4. Use Proper Error Handling

Implement proper error handling in configurations:

```hcl
# Use proper validation
variables {
  variable "app_port" {
    type = "number"
    default = 8080
    
    validation {
      condition = var.app_port > 0 && var.app_port <= 65535
      error_message = "Port must be between 1 and 65535"
    }
  }
}
```

### 5. Monitor and Log

Monitor variable operations and maintain logs:

```bash
# Enable verbose logging
spooky variables resolve ./my-project --verbose

# Monitor operations
watch -n 1 'ps aux | grep spooky'

# Check logs
tail -f /var/log/spooky/variables.log
```

## Getting Help

### Documentation Resources

1. **User Guide** - For usage questions and best practices
2. **API Reference** - For technical implementation details
3. **Examples** - For configuration patterns and use cases

### Common Questions

#### "Why can't I resolve my variables?"

1. Check variable configuration
2. Verify variable dependencies
3. Check for circular dependencies
4. Validate variable syntax

#### "How do I debug variable resolution issues?"

```bash
# Enable verbose output
spooky variables resolve ./my-project --verbose

# Test specific variables
spooky variables resolve ./my-project --variables var1,var2

# Check dependencies
spooky variables validate ./my-project --check-dependencies
```

#### "How do I fix validation issues?"

```bash
# Validate configuration
spooky variables validate ./my-project --verbose

# Check specific errors
spooky variables validate ./my-project --variables problematic-var

# Fix configuration issues
# Update variable configuration based on error messages
```

#### "How do I optimize variable operations?"

```bash
# Use variable filtering
spooky variables list ./my-project --variables var1,var2

# Use parallel operations
spooky variables resolve ./my-project --parallel 4

# Monitor resource usage
top -p $(pgrep spooky)
```

### When to Seek Additional Help

- Configuration validation passes but resolution still fails
- Performance issues persist after optimization
- Unusual error messages not covered in this guide
- Integration issues with other spooky components

For additional help, refer to the [User Guide](VARIABLES_USER_GUIDE.md) and [API Reference](VARIABLES_API_REFERENCE.md), or check the project documentation for more advanced troubleshooting techniques.

## Conclusion

The variables system provides robust, reliable variable management with comprehensive resolution, validation, and encryption capabilities. Most issues can be resolved by following the troubleshooting steps outlined in this guide. For persistent issues, enable verbose output and collect diagnostic information for further analysis.
