# Projects System

## Overview

The spooky projects system provides comprehensive project management capabilities including initialization, validation, and encryption. Projects are the foundational unit of organization in spooky, containing all configuration files, templates, and resources needed for automation tasks.

## Related Systems

This system organizes and contains all other spooky systems:

- **[Actions System](ACTIONS_SYSTEM.md)** - Projects contain action definitions in `actions.hcl` files
- **[Facts System](FACTS_SYSTEM.md)** - Projects have fact collections stored in `facts.db`
- **[Machines System](MACHINES_SYSTEM.md)** - Projects contain machine inventory in `machines.hcl` files
- **[Templates System](TEMPLATES_SYSTEM.md)** - Projects contain templates in the `templates/` directory
- **[Variables System](VARIABLES_SYSTEM.md)** - Projects contain variables in `variables.hcl` and `variables/` files
- **[Logging System](LOGGING_SYSTEM.md)** - Projects have logging configuration and generate logs
- **[Integrations System](INTEGRATIONS_SYSTEM.md)** - Projects coordinate all system integrations through the IntegrationManager
- **[SSH System](SSH_SYSTEM.md)** - Projects contain SSH configurations for machine connectivity

## Core Concepts

### Project Structure
A spooky project follows a standardized directory structure defined by the `project-directory.schema.hcl` schema:

```
project-name/
├── project.hcl          # Project metadata and configuration
├── machines.hcl         # Machine inventory (optional)
├── actions.hcl          # Action definitions (optional)
├── variables.hcl        # Variable definitions (optional)
├── variables/           # Modular variable files (optional)
├── templates/           # Template files (optional)
├── data/               # Data files for templates (optional)
└── schemas/            # Custom schemas (optional)
```

### Project Configuration
Projects are configured through `project.hcl` files that define:

- **Project metadata**: Name, description, version, author
- **Settings**: Parallel workers, timeouts, logging configuration
- **Dependencies**: Required schemas and extensions
- **Isolation**: Project-specific settings and access controls

## CLI Commands

### Project Management Commands

#### `spooky project init [project-name]`
Initialize a new spooky project with the specified name.

**Flags:**
- `--name` - Project name (defaults to directory name)
- `--description` - Project description
- `--version` - Project version
- `--author` - Project author
- `--email` - Project email
- `--url` - Project URL

**Examples:**
```bash
# Basic project initialization
spooky project init my-project

# Initialize with metadata
spooky project init my-project --name "My Project" --description "A test project"
```

#### `spooky project validate [project-path]`
Validate a spooky project structure and configuration.

**Examples:**
```bash
# Validate project structure and configuration
spooky project validate ./my-project
```

#### `spooky project encrypt [project-path]`
Encrypt all variables and machines in a project that have `encrypted=true`.

**Flags:**
- `--dry-run` - Show what would be encrypted without making changes

**Examples:**
```bash
# Encrypt project with dry-run
spooky project encrypt ./my-project --dry-run

# Encrypt project
spooky project encrypt ./my-project
```

## Project Configuration

### Project Metadata
Project metadata is defined in the `project` block:

```hcl
project {
  name = "my-spooky-project"
  description = "A comprehensive spooky project for infrastructure management"
  
  metadata {
    version = "1.0.0"
    author = "spooky-user"
    email = "user@example.com"
    url = "https://github.com/user/my-project"
    tags = ["infrastructure", "automation", "testing"]
  }
}
```

### Project Settings
Project settings control behavior and performance:

```hcl
project {
  name = "my-project"
  
  run {
    default_timeout = 300
    max_parallel = 4
    dry_run_default = false
    validate_before_run = true
    backup_before_changes = false
  }
}
```

## Project Validation

### Structure Validation
The project validation system checks:

1. **Directory Structure**: Validates against `project-directory.schema.hcl`
2. **Configuration Files**: Ensures required files exist and are valid HCL
3. **Schema Compliance**: Validates all files against their respective schemas
4. **Dependencies**: Checks for required schemas and extensions
5. **Permissions**: Validates file permissions and access controls

### Validation Process
```bash
# Validate project structure and configuration
spooky project validate ./my-project
```

**Output:**
```
🔍 Validating project: ./my-project
✅ Project validation passed - all required components present
📋 Schema compliance: project-directory.schema.hcl ✅
📋 Schema compliance: project.schema.hcl ✅
```

## Project Encryption

### Encryption Process
The project encryption system:

1. **Scans Variables**: Processes `variables.hcl` and `variables/*.hcl` files
2. **Scans Machines**: Processes `machines.hcl` files
3. **Identifies Encrypted Fields**: Finds fields with `encrypted = true`
4. **Loads Recipients**: Reads from `identities/recipients.txt`
5. **Encrypts Values**: Uses age encryption for sensitive data
6. **Updates Files**: Replaces plaintext values with encrypted armor

### Encryption Configuration
```hcl
# In variables.hcl
variables {
  variable "database_password" {
    value = "secret123"
    encrypted = true
    description = "Database password"
  }
}

# In machines.hcl
machines_inventory {
  machines {
    machine "web-server" {
      hostname = "web.example.com"
      user = "admin"
      
      authentication {
        method = "password"
        password {
          value = "secret456"
          encrypted = true
        }
      }
    }
  }
}
```

### Encryption Commands
```bash
# Dry run to see what would be encrypted
spooky project encrypt ./my-project --dry-run

# Actually encrypt the project
spooky project encrypt ./my-project
```

## Project Isolation

### Isolation Features
Project isolation provides:

- **Resource Isolation**: Prevents cross-project resource access
- **Configuration Isolation**: Project-specific settings and overrides
- **Security Boundaries**: Controlled access to external resources
- **Dependency Isolation**: Project-specific schema and extension loading

### Isolation Configuration
```hcl
project {
  name = "isolated-project"
  
  isolation {
    enabled = true
    strict_mode = true
    
    allowed_external_resources = [
      "https://api.example.com",
      "https://registry.example.com"
    ]
    
    blocked_resources = [
      "https://malicious-site.com"
    ]
    
    allowed_schemas = [
      "custom-schema.hcl"
    ]
  }
}
```

## Project Dependencies

### Schema Dependencies
Projects can depend on custom schemas:

```hcl
project {
  name = "schema-dependent-project"
  
  dependencies {
    schemas = [
      "custom-validation.schema.hcl",
      "extended-types.schema.hcl"
    ]
    
    extensions = [
      "custom-functions.hcl"
    ]
  }
}
```

### Dependency Resolution
The system resolves dependencies by:

1. **Loading Schemas**: Loading all required schema files
2. **Validating Dependencies**: Ensuring schemas are valid and accessible
3. **Registering Extensions**: Loading custom functions and validators
4. **Validating Against Schemas**: Using loaded schemas for validation

## Project Lifecycle

### Project States
Projects can exist in different states:

- **Initialized**: Basic structure created, ready for configuration
- **Configured**: All required files present and valid
- **Validated**: Passes all validation checks
- **Encrypted**: Sensitive data encrypted
- **Ready**: Ready for automation tasks

### State Transitions
```bash
# Initialize project
spooky project init my-project

# Configure project (manual file editing)
# ... edit project.hcl, machines.hcl, etc.

# Validate project
spooky project validate ./my-project

# Encrypt sensitive data
spooky project encrypt ./my-project

# Project is now ready for use
```

## Error Handling

### Common Validation Errors
- **Missing Required Files**: `project.hcl` not found
- **Invalid HCL Syntax**: Malformed configuration files
- **Schema Violations**: Configuration doesn't match schemas
- **Permission Issues**: File access or permission problems
- **Dependency Issues**: Missing or invalid dependencies

### Error Recovery
```bash
# Fix validation errors
spooky project validate ./my-project

# Address specific issues
# ... fix configuration files

# Re-validate
spooky project validate ./my-project
```

## Best Practices

### Project Organization
1. **Use Descriptive Names**: Choose clear, descriptive project names
2. **Organize by Purpose**: Group related projects in directories
3. **Version Control**: Use git for project version control
4. **Documentation**: Include README files for complex projects
5. **Templates**: Create reusable project templates

### Security Practices
1. **Encrypt Sensitive Data**: Use `encrypted = true` for secrets
2. **Secure File Permissions**: Set appropriate file permissions
3. **Use SSH Keys**: Prefer SSH keys over passwords
4. **Validate Inputs**: Always validate external inputs
5. **Audit Access**: Regularly audit project access

### Performance Practices
1. **Optimize Parallel Workers**: Set appropriate parallel worker counts
2. **Use Caching**: Enable caching for frequently accessed data
3. **Monitor Resources**: Monitor memory and CPU usage
4. **Optimize Templates**: Use efficient template rendering
5. **Batch Operations**: Group related operations

## Integration with Other Systems

### Integration Points
The projects system integrates with:

- **Facts System**: Provides project context for fact collection
- **Actions System**: Supplies project configuration for action running
- **Machines System**: Uses project machine inventory
- **Variables System**: Processes project variable definitions
- **Templates System**: Renders project-specific templates
- **Secrets System**: Manages project encryption and key management

### Cross-Project Operations
```bash
# Use project context for facts
spooky facts export ./my-project --output facts.hcl

# Run actions in project context
spooky actions run ./my-project --action deploy

# Validate project components
spooky machines validate ./my-project
spooky variables validate ./my-project
```

## Troubleshooting

### Common Issues

#### Project Validation Fails
```bash
# Check project structure
ls -la ./my-project/

# Validate specific components
spooky project validate ./my-project
```

#### Encryption Issues
```bash
# Check recipients file
cat ./my-project/identities/recipients.txt

# Dry run encryption
spooky project encrypt ./my-project --dry-run
```

#### Permission Issues
```bash
# Check file permissions
ls -la ./my-project/

# Fix permissions
chmod 600 ./my-project/identities/private.key
chmod 644 ./my-project/identities/recipients.txt
```

### Debug Information
```bash
# Enable verbose logging
export SPOOKY_LOG_LEVEL=debug

# Run validation with debug output
spooky project validate ./my-project
```

## Future Enhancements

### Planned Features
- **Project Templates**: Reusable project templates
- **Project Migration**: Tools for upgrading project formats
- **Project Backup**: Automated project backup and restore
- **Project Analytics**: Usage analytics and performance metrics
- **Project Collaboration**: Multi-user project management

### Extension Points
- **Custom Validators**: Project-specific validation rules
- **Custom Schemas**: Extended schema definitions
- **Custom Functions**: Project-specific template functions
- **Custom Integrations**: External system integrations
- **Custom Workflows**: Automated project workflows
