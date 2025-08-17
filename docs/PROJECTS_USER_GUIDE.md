# Projects User Guide

## Overview

The Projects System provides comprehensive project management capabilities for the spooky codebase. This guide will help you understand how to create, configure, and manage projects effectively.

> **See also**: [User Guides Index](USER_GUIDES_INDEX.md) - Complete overview of all user guides

## Quick Start

### Creating Your First Project

1. **Initialize a new project:**
   ```bash
   spooky project init my-first-project
   ```

2. **Navigate to the project directory:**
   ```bash
   cd my-first-project
   ```

3. **View project information:**
   ```bash
   spooky project info
   ```

4. **Validate the project:**
   ```bash
   spooky project validate
   ```

## Project Structure

A typical spooky project has the following structure:

```
my-project/
├── project.hcl              # Main project configuration
├── machines/                # Machine inventory
│   ├── inventory.hcl
│   └── config.hcl
├── actions/                 # Action definitions
│   ├── actions.hcl
│   └── scripts/
├── facts/                   # Fact collection
│   ├── facts.hcl
│   └── collectors/
├── templates/               # Template files
│   ├── templates.hcl
│   └── files/
├── variables/               # Variable definitions
│   ├── variables.hcl
│   └── env/
├── secrets/                 # Secret management
│   ├── secrets.hcl
│   └── keys/
└── .spooky/                 # Project metadata
    ├── state/
    └── logs/
```

## Project Configuration

### Basic Project Configuration

Create a `project.hcl` file in your project root:

```hcl
project {
  name = "my-web-application"
  description = "A web application deployment project"
  
  metadata {
    version = "1.0.0"
    author = "spooky-user"
    license = "MIT"
    tags = ["web", "deployment", "production"]
  }
  
  run {
    default_timeout = 300
    max_parallel = 4
    dry_run_default = false
    validate_before_run = true
    backup_before_changes = false
  }
  
  components {
    facts = true
    actions = true
    machines = true
    templates = true
    variables = true
    secrets = true
  }
}
```

### Advanced Project Configuration

For more complex projects, you can configure individual components:

```hcl
project {
  name = "enterprise-infrastructure"
  description = "Enterprise infrastructure management project"
  
  metadata {
    version = "2.0.0"
    author = "enterprise-team"
    license = "Proprietary"
    tags = ["enterprise", "infrastructure", "security"]
  }
  
  run {
    default_timeout = 1800
    max_parallel = 16
    dry_run_default = false
    validate_before_run = true
    backup_before_changes = true
  }
  
  components {
    facts {
      enabled = true
      collection_interval = "15m"
      storage_backend = "badgerdb"
    }
    
    actions {
      enabled = true
      execution_mode = "parallel"
      max_retries = 10
      rollback_enabled = true
    }
    
    machines {
      enabled = true
      connection_pool_size = 50
      health_check_interval = "2m"
    }
    
    templates {
      enabled = true
      cache_enabled = true
      cache_ttl = "30m"
    }
    
    variables {
      enabled = true
      encryption_enabled = true
      validation_enabled = true
    }
    
    secrets {
      enabled = true
      encryption_method = "age"
      key_management = "external"
    }
  }
  
  environment {
    development {
      enabled = true
      machines = ["dev-*"]
      variables = ["dev/*.hcl"]
      security_level = "medium"
    }
    
    staging {
      enabled = true
      machines = ["staging-*"]
      variables = ["staging/*.hcl"]
      security_level = "high"
    }
    
    production {
      enabled = true
      machines = ["prod-*"]
      variables = ["prod/*.hcl"]
      security_level = "maximum"
      backup_enabled = true
    }
  }
  
  security {
    encrypt_sensitive_data = true
    validate_file_permissions = true
    audit_all_operations = true
    access_control_enabled = true
    
    allowed_users = ["admin", "operator", "monitor"]
    allowed_groups = ["devops", "sre", "security"]
    
    encryption {
      method = "age"
      key_path = "/etc/spooky/keys"
      backup_keys = true
      rotation_interval = "90d"
    }
    
    compliance {
      sox_enabled = true
      pci_enabled = true
      gdpr_enabled = true
      audit_retention = "7y"
    }
  }
}
```

## Project Types

### Basic Project

A simple project for basic automation tasks:

```hcl
project {
  name = "basic-automation"
  description = "Basic automation project"
  
  metadata {
    version = "1.0.0"
    author = "spooky-user"
    license = "MIT"
  }
  
  settings {
    parallel_workers = 2
    timeout_seconds = 60
    log_level = "info"
  }
  
  components {
    facts = true
    actions = true
    machines = true
  }
}
```

### Web Application Project

A project for web application deployment:

```hcl
project {
  name = "web-app-deployment"
  description = "Web application deployment project"
  
  metadata {
    version = "2.0.0"
    author = "web-team"
    license = "Apache-2.0"
    tags = ["web", "deployment", "automation"]
  }
  
  settings {
    parallel_workers = 8
    timeout_seconds = 600
    log_level = "debug"
  }
  
  components {
    facts {
      enabled = true
      collection_interval = "30m"
    }
    
    actions {
      enabled = true
      execution_mode = "parallel"
      max_retries = 5
    }
    
    machines {
      enabled = true
      connection_pool_size = 20
    }
    
    templates {
      enabled = true
      cache_enabled = true
    }
    
    variables {
      enabled = true
      encryption_enabled = true
    }
  }
  
  environment {
    development {
      enabled = true
      machines = ["dev-web-1", "dev-web-2"]
    }
    
    staging {
      enabled = true
      machines = ["staging-web-1", "staging-web-2"]
    }
    
    production {
      enabled = true
      machines = ["prod-web-1", "prod-web-2", "prod-web-3"]
      security_level = "high"
    }
  }
}
```

### Database Project

A project for database management:

```hcl
project {
  name = "database-management"
  description = "Database management and maintenance project"
  
  metadata {
    version = "1.5.0"
    author = "db-team"
    license = "MIT"
    tags = ["database", "maintenance", "backup"]
  }
  
  settings {
    parallel_workers = 4
    timeout_seconds = 1200
    log_level = "info"
  }
  
  components {
    facts {
      enabled = true
      collection_interval = "1h"
      storage_backend = "badgerdb"
    }
    
    actions {
      enabled = true
      execution_mode = "sequential"
      max_retries = 3
    }
    
    machines {
      enabled = true
      connection_pool_size = 10
    }
    
    secrets {
      enabled = true
      encryption_method = "age"
    }
  }
  
  environment {
    development {
      enabled = true
      machines = ["dev-db-1"]
      security_level = "low"
    }
    
    production {
      enabled = true
      machines = ["prod-db-1", "prod-db-2"]
      security_level = "high"
      backup_enabled = true
    }
  }
}
```

## Environment Configuration

### Development Environment

```hcl
environment {
  development {
    enabled = true
    machines = ["dev-server-1", "dev-server-2"]
    variables = ["dev-vars.hcl"]
    security_level = "low"
    
    metadata {
      purpose = "development"
      team = "dev-team"
      cost_center = "dev-ops"
    }
  }
}
```

### Staging Environment

```hcl
environment {
  staging {
    enabled = true
    machines = ["staging-server-1", "staging-server-2"]
    variables = ["staging-vars.hcl"]
    security_level = "medium"
    
    metadata {
      purpose = "testing"
      team = "qa-team"
      cost_center = "qa-ops"
    }
  }
}
```

### Production Environment

```hcl
environment {
  production {
    enabled = true
    machines = ["prod-server-1", "prod-server-2", "prod-server-3"]
    variables = ["prod-vars.hcl"]
    security_level = "high"
    backup_enabled = true
    
    metadata {
      purpose = "production"
      team = "ops-team"
      cost_center = "prod-ops"
      sla = "99.9%"
    }
  }
}
```

## Security Configuration

### Basic Security

```hcl
security {
  encrypt_sensitive_data = true
  validate_file_permissions = true
  audit_logging = true
  access_control_enabled = true
  
  allowed_users = ["admin", "operator"]
  allowed_groups = ["devops", "sre"]
}
```

### Advanced Security

```hcl
security {
  encrypt_sensitive_data = true
  validate_file_permissions = true
  audit_all_operations = true
  access_control_enabled = true
  
  allowed_users = ["admin", "operator", "monitor"]
  allowed_groups = ["devops", "sre", "security"]
  
  encryption {
    method = "age"
    key_path = "/etc/spooky/keys"
    backup_keys = true
    rotation_interval = "90d"
  }
  
  compliance {
    sox_enabled = true
    pci_enabled = true
    gdpr_enabled = true
    audit_retention = "7y"
  }
}
```

## CLI Commands

### Project Initialization

```bash
# Initialize a new project
spooky project init my-project

# Initialize with template
spooky project init my-project --template web-app

# Initialize with configuration
spooky project init my-project --config config.hcl

# Initialize in specific directory
spooky project init my-project --path /path/to/project
```

### Project Management

```bash
# Show project information
spooky project info

# Show project status
spooky project status

# Show project configuration
spooky project config

# Show project structure
spooky project structure
```

### Project Validation

```bash
# Validate project
spooky project validate

# Validate with specific checks
spooky project validate --check structure,config

# Validate with verbose output
spooky project validate --verbose

# Validate and fix issues
spooky project validate --fix
```

### Project Operations

```bash
# Start project
spooky project start

# Stop project
spooky project stop

# Restart project
spooky project restart

# Clean project
spooky project clean
```

## Best Practices

### Project Organization

1. **Use descriptive project names** that clearly indicate the purpose
2. **Organize projects by function** (e.g., web-apps, databases, infrastructure)
3. **Use consistent naming conventions** across all projects
4. **Document project purpose** in the description field
5. **Use tags** to categorize and search projects

### Configuration Management

1. **Use version control** for all project configurations
2. **Separate environments** with different configuration files
3. **Use variables** for environment-specific values
4. **Validate configurations** before deployment
5. **Document configuration changes** in commit messages

### Security

1. **Encrypt sensitive data** using the secrets system
2. **Use appropriate security levels** for each environment
3. **Implement access controls** for project operations
4. **Audit project operations** for compliance
5. **Rotate encryption keys** regularly

### Performance

1. **Configure appropriate timeouts** for operations
2. **Use parallel workers** for concurrent operations
3. **Enable caching** for frequently accessed data
4. **Monitor project performance** and adjust settings
5. **Optimize resource usage** based on workload

## Common Use Cases

### Web Application Deployment

```hcl
project {
  name = "web-app-deployment"
  description = "Deploy web applications to multiple environments"
  
  components {
    facts {
      enabled = true
      collection_interval = "30m"
    }
    
    actions {
      enabled = true
      execution_mode = "parallel"
      max_retries = 3
    }
    
    machines {
      enabled = true
      connection_pool_size = 20
    }
    
    templates {
      enabled = true
      cache_enabled = true
    }
    
    variables {
      enabled = true
      encryption_enabled = true
    }
  }
  
  environment {
    development {
      enabled = true
      machines = ["dev-web-1", "dev-web-2"]
      variables = ["dev-vars.hcl"]
    }
    
    staging {
      enabled = true
      machines = ["staging-web-1", "staging-web-2"]
      variables = ["staging-vars.hcl"]
    }
    
    production {
      enabled = true
      machines = ["prod-web-1", "prod-web-2", "prod-web-3"]
      variables = ["prod-vars.hcl"]
      security_level = "high"
    }
  }
}
```

### Infrastructure Management

```hcl
project {
  name = "infrastructure-management"
  description = "Manage infrastructure components and services"
  
  components {
    facts {
      enabled = true
      collection_interval = "15m"
      storage_backend = "badgerdb"
    }
    
    actions {
      enabled = true
      execution_mode = "parallel"
      max_retries = 5
      rollback_enabled = true
    }
    
    machines {
      enabled = true
      connection_pool_size = 50
      health_check_interval = "2m"
    }
    
    secrets {
      enabled = true
      encryption_method = "age"
      key_management = "external"
    }
  }
  
  environment {
    development {
      enabled = true
      machines = ["dev-*"]
      variables = ["dev/*.hcl"]
      security_level = "medium"
    }
    
    production {
      enabled = true
      machines = ["prod-*"]
      variables = ["prod/*.hcl"]
      security_level = "maximum"
      backup_enabled = true
    }
  }
  
  security {
    encrypt_sensitive_data = true
    validate_file_permissions = true
    audit_all_operations = true
    access_control_enabled = true
    
    allowed_users = ["admin", "operator"]
    allowed_groups = ["devops", "sre"]
    
    encryption {
      method = "age"
      key_path = "/etc/spooky/keys"
      backup_keys = true
      rotation_interval = "90d"
    }
  }
}
```

### Database Management

```hcl
project {
  name = "database-management"
  description = "Database maintenance and backup operations"
  
  components {
    facts {
      enabled = true
      collection_interval = "1h"
    }
    
    actions {
      enabled = true
      execution_mode = "sequential"
      max_retries = 3
    }
    
    machines {
      enabled = true
      connection_pool_size = 10
    }
    
    secrets {
      enabled = true
      encryption_method = "age"
    }
  }
  
  environment {
    development {
      enabled = true
      machines = ["dev-db-1"]
      security_level = "low"
    }
    
    production {
      enabled = true
      machines = ["prod-db-1", "prod-db-2"]
      security_level = "high"
      backup_enabled = true
    }
  }
  
  security {
    encrypt_sensitive_data = true
    audit_all_operations = true
    access_control_enabled = true
    
    allowed_users = ["db-admin"]
    allowed_groups = ["database-team"]
    
    compliance {
      sox_enabled = true
      audit_retention = "7y"
    }
  }
}
```

## Troubleshooting

### Common Issues

#### Project Not Found
```bash
# Check if project exists
ls -la my-project/

# Check project path
spooky project info --path /full/path/to/project

# Validate project structure
spooky project validate
```

#### Configuration Errors
```bash
# Check configuration syntax
spooky project validate --syntax-only

# Show configuration details
spooky project config --verbose

# Export configuration for review
spooky project config --export
```

#### Permission Issues
```bash
# Check file permissions
ls -la project.hcl

# Fix permissions
chmod 644 project.hcl

# Check user access
spooky project info --user
```

#### Component Issues
```bash
# Check component status
spooky project status --components

# Validate specific component
spooky project validate --component facts

# Restart component
spooky project restart --component actions
```

## Related Documentation

- [Projects System](PROJECTS_SYSTEM.md) - Complete system overview
- [Projects API Reference](PROJECTS_API_REFERENCE.md) - API documentation
- [Projects Troubleshooting](PROJECTS_TROUBLESHOOTING.md) - Troubleshooting guide
- [Configuration Management](CONFIGURATION_SYSTEM.md) - Configuration integration
- [Schema System](SCHEMA_SYSTEM.md) - Schema integration
- [CLI Reference](CLI_REFERENCE.md) - CLI integration
