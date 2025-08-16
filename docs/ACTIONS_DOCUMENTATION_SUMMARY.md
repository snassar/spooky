# Actions System Documentation Summary

## Overview

This document provides a comprehensive overview of the spooky actions system documentation. It serves as a guide to help you find the right documentation for your needs and understand how all the pieces fit together.

**Status: Implemented** - The actions system is fully implemented with comprehensive functionality for action definition, validation, and execution.

## Documentation Structure

### 📚 Core Documentation

#### 1. [User Guide](ACTIONS_USER_GUIDE.md)
**Audience:** End users, system administrators, DevOps engineers
**Purpose:** Complete guide to using the actions system

**What it covers:**
- Getting started with action configuration
- Basic and advanced usage patterns
- Action types and configuration
- Dependency management and execution
- CLI commands and options
- Real-world examples and use cases

**When to use:** Start here if you're new to spooky actions or need to understand how to use the system effectively.

#### 2. [API Reference](ACTIONS_API_REFERENCE.md)
**Audience:** Developers, system integrators, contributors
**Purpose:** Technical reference for the actions system APIs and implementation

**What it covers:**
- Core interfaces and type definitions
- Implementation details and algorithms
- Error handling patterns
- Configuration rules and schemas
- CLI integration details
- Code examples and patterns

**When to use:** Use this when developing with the actions system, extending functionality, or debugging implementation issues.

#### 3. [Troubleshooting Guide](ACTIONS_TROUBLESHOOTING.md)
**Audience:** System administrators, support engineers, users experiencing issues
**Purpose:** Solutions for common problems and debugging techniques

**What it covers:**
- Common error messages and solutions
- Configuration problems and fixes
- Performance issues and optimization
- Network and connectivity issues
- Dependency and execution problems
- Recovery procedures and prevention strategies

**When to use:** Use this when encountering problems or need to debug issues with the actions system.

### 📁 Examples Directory

#### [Examples Overview](examples/README.md)
**Audience:** All users
**Purpose:** Practical examples and configuration patterns

**What it covers:**
- Basic action configuration
- Project setup examples
- Machine-specific configurations
- Integration examples with other components
- Best practices and patterns

**Example Files:**
- [`basic-actions-project.hcl`](examples/basic-actions-project.hcl) - Basic action configuration
- [`advanced-actions-config.hcl`](examples/advanced-actions-config.hcl) - Advanced configuration
- [`actions-integration.hcl`](examples/actions-integration.hcl) - Integration with other components

**When to use:** Use these as starting points for your own configurations or to learn best practices.

## Quick Start Guide

### For New Users

1. **Read the User Guide** - Start with [ACTIONS_USER_GUIDE.md](ACTIONS_USER_GUIDE.md) to understand the basics
2. **Try the Examples** - Copy and customize examples from the [examples/](examples/) directory
3. **Test Your Setup** - Use `spooky actions validate` and `spooky actions run --dry-run` to test
4. **Check Troubleshooting** - If you encounter issues, refer to [ACTIONS_TROUBLESHOOTING.md](ACTIONS_TROUBLESHOOTING.md)

### For Developers

1. **Review the API Reference** - Understand the interfaces and implementation in [ACTIONS_API_REFERENCE.md](ACTIONS_API_REFERENCE.md)
2. **Study the Examples** - See how the APIs are used in practice
3. **Check the Code** - Review the actual implementation in `internal/actions/`
4. **Test Your Changes** - Use the examples to test your modifications

### For System Administrators

1. **Start with User Guide** - Understand the system capabilities
2. **Review Examples** - See real-world configuration patterns
3. **Plan Your Actions** - Design your action orchestration strategy
4. **Implement Gradually** - Start with basic actions and expand
5. **Monitor and Validate** - Use validation and testing regularly

## Actions System Overview

### Core Concepts

The actions system provides a declarative way to define and execute automation tasks across multiple machines. Actions are defined in HCL configuration files and can include:

- **Command execution** - Run commands on target machines
- **Template rendering** - Generate dynamic content from templates
- **File operations** - Upload, download, and manage files
- **Conditional execution** - Execute actions based on conditions
- **Parallel execution** - Run actions concurrently across machines

### Action Configuration

Actions are defined in `actions.hcl` files within spooky projects:

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
  
  action "restart-services" {
    description = "Restart application services"
    
    machines = ["web-server", "db-server"]
    parallel = false
    
    command = "sudo systemctl restart myapp"
  }
}
```

### CLI Commands

The actions system provides comprehensive CLI commands:

```bash
# List actions in a project
spooky actions list ./my-project

# Validate action configurations
spooky actions validate ./my-project

# Run actions with dry-run mode
spooky actions run ./my-project --dry-run

# Run actions with plan mode
spooky actions run ./my-project --plan

# Execute actions
spooky actions run ./my-project

# Run with parallel execution
spooky actions run ./my-project --parallel 4

# Run with decryption for debugging
spooky actions run ./my-project --decrypt
```

### Action Execution Modes

1. **Dry-run Mode** - Simulate execution without making changes
2. **Plan Mode** - Show execution plan without running actions
3. **Normal Mode** - Execute actions on target machines
4. **Parallel Mode** - Execute actions concurrently across machines

### Machine Targeting

Actions can be targeted to specific machines using:

- **Machine names** - Target specific machines by name
- **Tags** - Target machines with specific tags
- **Complex filters** - Use filter expressions for advanced targeting

```bash
# Target specific machines
spooky actions run ./my-project --machine web-server

# Target machines by tags
spooky actions run ./my-project --tags environment=production

# Use complex filters
spooky actions run ./my-project --filter "environment=production AND role=web"
```

## Implementation Details

### Core Components

1. **Action Manager** - Manages action lifecycle and execution
2. **Action Validator** - Validates action configurations
3. **Action Loader** - Loads actions from configuration files
4. **SSH Manager** - Handles SSH connections and command execution
5. **Template Engine** - Renders templates with variables

### Integration Points

The actions system integrates with:

- **Machines System** - For machine inventory and connectivity
- **Variables System** - For dynamic configuration values
- **Templates System** - For template rendering and management
- **Secrets System** - For encrypted configuration and credentials
- **Facts System** - For machine-specific data and context

### Error Handling

The actions system provides comprehensive error handling:

- **Validation errors** - Configuration and syntax validation
- **Connection errors** - SSH connectivity issues
- **Execution errors** - Command execution failures
- **Template errors** - Template rendering issues
- **Dependency errors** - Missing dependencies and resources

## Best Practices

### Action Design

1. **Use descriptive names** for actions
2. **Provide clear descriptions** for action purpose
3. **Group related actions** logically
4. **Use templates** for dynamic content
5. **Test actions** with dry-run mode

### Configuration Management

1. **Validate configurations** before execution
2. **Use variables** for configuration values
3. **Encrypt sensitive data** using age encryption
4. **Version control** action configurations
5. **Document action dependencies** and requirements

### Execution Strategy

1. **Start with dry-run** to preview changes
2. **Use parallel execution** for efficiency
3. **Monitor execution** progress and results
4. **Handle errors** gracefully with retry logic
5. **Log execution** details for audit trails

## Troubleshooting

### Common Issues

1. **Action validation fails** - Check syntax and schema compliance
2. **SSH connection errors** - Verify machine connectivity and authentication
3. **Template rendering errors** - Check template syntax and variable availability
4. **Command execution failures** - Verify command paths and permissions
5. **Parallel execution issues** - Check resource limits and connection pools

### Debug Commands

```bash
# Enable verbose logging
export SPOOKY_LOG_LEVEL=debug

# Validate with detailed output
spooky actions validate ./my-project --verbose

# Test with dry-run and verbose output
spooky actions run ./my-project --dry-run --verbose

# Check machine connectivity
spooky machines ping ./my-project --auth
```

## Related Documentation

- [Actions User Guide](ACTIONS_USER_GUIDE.md) - Complete user guide
- [Actions API Reference](ACTIONS_API_REFERENCE.md) - Technical reference
- [Actions Troubleshooting](ACTIONS_TROUBLESHOOTING.md) - Troubleshooting guide
- [System Design](../design/systems/actions-system.md) - System design documentation
- [CLI Reference](CLI_REFERENCE.md) - CLI command reference
