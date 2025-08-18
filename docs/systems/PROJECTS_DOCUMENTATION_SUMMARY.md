# Projects Documentation Summary

## Overview

This document provides a comprehensive overview of spooky project management, including project structure, configuration, and usage patterns.

**Status: Implemented** - The project system is fully implemented with comprehensive functionality for project initialization, validation, and management.

## Project Structure

A spooky project follows a standardized directory structure defined by the `project-directory.schema.hcl` schema:

```
my-project/
├── project.hcl              # Project metadata and configuration
├── machines.hcl             # Machine inventory (optional)
├── actions.hcl              # Action definitions (optional)
├── variables.hcl            # Project variables (optional)
├── variables/               # Additional variable files (optional)
│   ├── environment.hcl
│   └── secrets.hcl
├── templates/               # Template files (optional)
│   ├── deployment.sh.tmpl
│   └── config.yaml.tmpl
└── schemas/                 # Project-specific schemas (optional)
    └── custom.schema.hcl
```

## Project Configuration

### project.hcl

The main project configuration file contains metadata and settings:

```hcl
project {
  name = "my-automation-project"
  description = "Automation project for web application deployment"
  
  metadata {
    version = "1.0.0"
    author = "John Doe"
    email = "john@example.com"
    url = "https://github.com/example/my-automation"
    tags = ["web", "deployment", "production"]
  }
  
  run {
    default_timeout = 300
    max_parallel = 4
    dry_run_default = false
    validate_before_run = true
    backup_before_changes = false
  }
}
```

### machines.hcl

Machine inventory configuration:

```hcl
machines_inventory {
  machines {
    machine "web-server" {
      hostname = "web.example.com"
      port = 22
      user = "admin"
      
      authentication {
        method = "ssh_key"
        key_path = "~/.ssh/id_rsa"
      }
      
      tags = ["production", "web"]
      
      metadata {
        environment = "production"
        role = "web"
      }
    }
  }
}
```

### actions.hcl

Action definitions for automation:

```hcl
actions {
  action "deploy-web" {
    description = "Deploy web application"
    
    machines = ["web-server"]
    parallel = true
    
    template {
      source = "templates/deploy.sh.tmpl"
      destination = "/tmp/deploy.sh"
      permissions = "0755"
    }
    
    command = "/tmp/deploy.sh"
  }
}
```

### variables.hcl

Project variables and configuration:

```hcl
variables {
  variable "app_version" {
    value = "1.0.0"
    description = "Application version to deploy"
  }
  
  variable "database_url" {
    value = "postgresql://user:pass@localhost:5432/mydb"
    description = "Database connection URL"
    sensitive = true
    encrypted = true
  }
}
```

## Project Commands

### Initialize Project

```bash
# Initialize a new project
spooky project init my-automation

# Initialize with specific metadata
spooky project init my-automation \
  --name "My Automation Project" \
  --description "Web application deployment automation" \
  --version "1.0.0" \
  --author "John Doe"
```

### Validate Project

```bash
# Validate project structure and configuration
spooky project validate my-automation
```

This command validates:
- Project directory structure compliance
- Configuration file syntax and schema validation
- Required file presence and format
- Cross-file consistency and dependencies

### Encrypt Project

```bash
# Encrypt all variables and machines with encrypted=true
spooky project encrypt my-automation

# Show what would be encrypted without making changes
spooky project encrypt my-automation --dry-run
```

## Project Lifecycle

### 1. Initialization

```bash
# Create new project
spooky project init my-automation
cd my-automation

# Edit configuration files
# - project.hcl (metadata and settings)
# - machines.hcl (inventory)
# - actions.hcl (automation)
# - variables.hcl (configuration)
```

### 2. Validation

```bash
# Validate project structure
spooky project validate .

# Validate specific components
spooky machines validate .
spooky actions validate .
spooky variables validate .
```

### 3. Testing

```bash
# Test machine connectivity
spooky machines ping .

# Test with authentication
spooky machines ping . --auth

# Export facts for analysis
spooky facts export . --output facts.hcl
```

### 4. Execution

```bash
# Run actions with dry-run
spooky actions run . --dry-run

# Run actions with plan
spooky actions run . --plan

# Execute actions
spooky actions run .
```

## Project Configuration Options

### Metadata Fields

- `name` - Project name (required)
- `description` - Project description
- `version` - Project version
- `author` - Project author
- `email` - Author email
- `url` - Project URL
- `tags` - Project tags

### Settings Fields

- `parallel_workers` - Number of parallel workers (default: 1)
- `timeout_seconds` - Operation timeout in seconds (default: 300)
- `log_level` - Logging level (debug, info, warn, error)
- `log_format` - Logging format (text, json)

## Project Validation

The project validation system checks:

### Structure Validation

- Directory structure compliance with schema
- Required files presence
- Optional directories and files
- File naming conventions

### Configuration Validation

- HCL syntax validation
- Schema compliance validation
- Cross-reference validation
- Dependency validation

### Content Validation

- Machine configuration validation
- Action definition validation
- Variable definition validation
- Template validation

## Project Encryption

Spooky supports encryption of sensitive data using age encryption:

### Encrypted Fields

- Variable values with `encrypted = true`
- Machine authentication credentials with `encrypted = true`
- Template content with `encrypted = true`

### Encryption Configuration

```hcl
# Global configuration (spooky.hcl)
age {
  identities = "~/.config/spooky/identities"
  recipients = "~/.config/spooky/recipients.txt"
}
```

### Encryption Commands

```bash
# Encrypt project variables
spooky variables armor .

# Encrypt project machines
spooky machines encrypt .

# Encrypt entire project
spooky project encrypt .
```

## Project Examples

### Basic Web Application

```bash
# Initialize project
spooky project init web-app

# Configure machines
cat > machines.hcl << 'EOF'
machines {
  machine "web-server" {
    hostname = "web.example.com"
    host = "192.168.1.100"
    port = 22
    user = "admin"
    
    authentication {
      method = "ssh_key"
      key_path = "~/.ssh/id_rsa"
    }
    
    tags = {
      environment = "production"
      role = "web"
    }
  }
}
EOF

# Configure actions
cat > actions.hcl << 'EOF'
actions {
  action "deploy" {
    description = "Deploy web application"
    
    machines = ["web-server"]
    
    template {
      source = "templates/deploy.sh.tmpl"
      destination = "/tmp/deploy.sh"
      permissions = "0755"
    }
    
    command = "/tmp/deploy.sh"
  }
}
EOF

# Configure variables
cat > variables.hcl << 'EOF'
variables {
  variable "app_version" {
    type = "string"
    description = "Application version"
    default = "1.0.0"
  }
}
EOF

# Validate and run
spooky project validate .
spooky machines ping .
spooky actions run . --dry-run
spooky actions run .
```

### Multi-Environment Project

```bash
# Initialize project
spooky project init multi-env

# Configure machines for multiple environments
cat > machines.hcl << 'EOF'
machines {
  machine "web-prod" {
    hostname = "web-prod.example.com"
    host = "10.0.1.100"
    user = "admin"
    
    authentication {
      method = "ssh_key"
      key_path = "~/.ssh/prod_key"
    }
    
    tags = {
      environment = "production"
      role = "web"
    }
  }
  
  machine "web-staging" {
    hostname = "web-staging.example.com"
    host = "10.0.2.100"
    user = "admin"
    
    authentication {
      method = "ssh_key"
      key_path = "~/.ssh/staging_key"
    }
    
    tags = {
      environment = "staging"
      role = "web"
    }
  }
}
EOF

# Run actions on specific environments
spooky actions run . --tags environment=staging
spooky actions run . --tags environment=production
```

## Best Practices

### Project Organization

1. **Use descriptive names** for projects and components
2. **Group related machines** using tags and groups
3. **Separate environments** using tags and different configurations
4. **Use variables** for configuration values
5. **Encrypt sensitive data** using age encryption

### Configuration Management

1. **Validate projects** before execution
2. **Test connectivity** before running actions
3. **Use dry-run mode** to preview changes
4. **Version control** project configurations
5. **Document project purpose** in descriptions

### Security

1. **Encrypt sensitive variables** and credentials
2. **Use SSH keys** instead of passwords
3. **Limit machine access** to necessary users
4. **Regularly rotate** authentication credentials
5. **Validate configurations** for security issues

## Troubleshooting

### Common Issues

1. **Project validation fails** - Check file syntax and schema compliance
2. **Machine connectivity issues** - Verify SSH configuration and network access
3. **Action execution fails** - Check template syntax and command paths
4. **Variable resolution errors** - Verify variable definitions and dependencies

### Debug Commands

```bash
# Enable verbose logging
export SPOOKY_LOG_LEVEL=debug

# Validate with detailed output
spooky project validate . --verbose

# Test connectivity with authentication
spooky machines ping . --auth --verbose

# Run actions with detailed output
spooky actions run . --dry-run --verbose
```
