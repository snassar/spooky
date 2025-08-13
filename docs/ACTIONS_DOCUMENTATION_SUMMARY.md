# Actions System Documentation Summary

## Overview

This document provides a comprehensive overview of the spooky actions system documentation. It serves as a guide to help you find the right documentation for your needs and understand how all the pieces fit together.

## Documentation Structure

### 📚 Core Documentation

#### 1. [User Guide](ACTIONS_USER_GUIDE.md)
**Audience:** End users, system administrators, DevOps engineers
**Purpose:** Complete guide to using the actions system

**What it covers:**
- Getting started with action configuration
- Basic and advanced usage patterns
- Action types and configuration
- Dependency management and running
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
- Dependency and running problems
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

## Documentation Navigation

### By Use Case

#### Getting Started
- **New to spooky actions?** → [User Guide](ACTIONS_USER_GUIDE.md) - Getting Started section
- **Setting up your first project?** → [User Guide](ACTIONS_USER_GUIDE.md) - Basic Usage section
- **Need examples?** → [Examples Directory](examples/) - Basic examples

#### Configuration
- **Action configuration?** → [User Guide](ACTIONS_USER_GUIDE.md) - Action Configuration section
- **Machine targeting?** → [User Guide](ACTIONS_USER_GUIDE.md) - Action Targeting section
- **Dependency management?** → [User Guide](ACTIONS_USER_GUIDE.md) - Action Dependencies section

#### Integration
- **Using actions with machines?** → [User Guide](ACTIONS_USER_GUIDE.md) - Machine Integration section
- **Using actions with variables?** → [User Guide](ACTIONS_USER_GUIDE.md) - Variables Integration section
- **Using actions with facts?** → [User Guide](ACTIONS_USER_GUIDE.md) - Facts Integration section

#### Troubleshooting
- **Action loading failures?** → [Troubleshooting Guide](ACTIONS_TROUBLESHOOTING.md) - Action Loading Errors section
- **Run failures?** → [Troubleshooting Guide](ACTIONS_TROUBLESHOOTING.md) - Action Run Errors section
- **Dependency issues?** → [Troubleshooting Guide](ACTIONS_TROUBLESHOOTING.md) - Dependency Errors section

#### Development
- **Extending the system?** → [API Reference](ACTIONS_API_REFERENCE.md) - Core Interfaces section
- **Adding new action types?** → [API Reference](ACTIONS_API_REFERENCE.md) - Action Types section
- **Custom validators?** → [API Reference](ACTIONS_API_REFERENCE.md) - ActionValidator Interface section

### By Component

#### Action Configuration
- **Overview** → [User Guide](ACTIONS_USER_GUIDE.md) - Actions System Concepts
- **Configuration** → [User Guide](ACTIONS_USER_GUIDE.md) - Creating Actions
- **Troubleshooting** → [Troubleshooting Guide](ACTIONS_TROUBLESHOOTING.md) - Action Configuration Issues
- **API** → [API Reference](ACTIONS_API_REFERENCE.md) - Action Type

#### Action Running
- **Overview** → [User Guide](ACTIONS_USER_GUIDE.md) - Action Lifecycle
- **Configuration** → [User Guide](ACTIONS_USER_GUIDE.md) - Running Actions
- **Troubleshooting** → [Troubleshooting Guide](ACTIONS_TROUBLESHOOTING.md) - Action Run Errors
- **API** → [API Reference](ACTIONS_API_REFERENCE.md) - ActionRunContext Type

#### Dependency Management
- **Overview** → [User Guide](ACTIONS_USER_GUIDE.md) - Action Dependencies
- **Configuration** → [User Guide](ACTIONS_USER_GUIDE.md) - Dependency Configuration
- **Troubleshooting** → [Troubleshooting Guide](ACTIONS_TROUBLESHOOTING.md) - Dependency Errors
- **API** → [API Reference](ACTIONS_API_REFERENCE.md) - ActionDependency Type

#### CLI Commands
- **Overview** → [User Guide](ACTIONS_USER_GUIDE.md) - CLI Commands
- **Commands** → [User Guide](ACTIONS_USER_GUIDE.md) - Command Options
- **Troubleshooting** → [Troubleshooting Guide](ACTIONS_TROUBLESHOOTING.md) - CLI Issues
- **API** → [API Reference](ACTIONS_API_REFERENCE.md) - CLI Integration

## Key Concepts

### Actions System Architecture

The actions system consists of several key components:

1. **Action** - Individual operations to be run on machines
2. **ActionCollection** - Groups of related actions
3. **ActionPlan** - Run plan with dependency resolution
4. **ActionsIntegration** - Primary interface for action management
5. **ActionValidator** - Validation of action configurations
6. **SSHManager** - SSH connection and running capabilities
7. **CLI Commands** - User interface for action management

### Data Flow

1. **Action Loading** - Load actions from HCL configuration files
2. **Validation** - Validate action configurations against schema
3. **Dependency Resolution** - Resolve action dependencies and create run order
4. **Machine Targeting** - Determine target machines for each action
5. **Run Planning** - Create run plan with proper ordering
6. **Action Running** - Run actions on target machines via SSH
7. **Result Collection** - Collect and aggregate run results

### Action Types

The system supports several types of actions:

- **Command Actions** - Run shell commands on remote machines
- **Script Actions** - Run script files with template processing
- **Template Deploy Actions** - Deploy template files with variable substitution
- **File Copy Actions** - Copy files between machines
- **Service Control Actions** - Control system services

## Common Patterns

### Basic Action Configuration

```bash
# List actions in a project
spooky actions list ./my-project

# Validate action configuration
spooky actions validate ./my-project

# Run actions
spooky actions run ./my-project
```

### Advanced Configuration

```hcl
action "deploy-application" {
  type = "template_deploy"
  description = "Deploy application configuration"
  
  template {
    source = "templates/app.conf.tmpl"
    destination = "/etc/app/app.conf"
    permissions = "0644"
    owner = "app"
    group = "app"
  }
  
  machines = ["app-server"]
  dependencies = ["prepare-database"]
  parallel = false
  timeout = 300
}
```

## Implementation Status

### ✅ Completed Features

- **Core Action Management**
  - Action configuration with HCL schema validation
  - Multiple action types (command, script, template_deploy, file_copy, service_control)
  - Action dependency resolution and run ordering
  - Machine targeting by names and tags
  - Parallel running support
  - Comprehensive validation and error handling

- **Advanced Features**
  - Template deployment with variable substitution
  - Service control with systemd integration
  - File copy operations with permissions
  - Script running with template processing
  - Resource limits and timeout management
  - Retry logic and error recovery

- **CLI Integration**
  - `spooky actions list` - List actions in a project
  - `spooky actions validate` - Validate action configuration
  - `spooky actions run` - Run actions with plan and dry-run modes

### 🚧 In Progress / Planned Features

- **Advanced Scripting** - Enhanced script running with better error handling
- **Template Functions** - Additional template functions and helpers
- **Action History** - Track action run history
- **Rollback Support** - Automatic rollback capabilities
- **Advanced Scheduling** - Time-based action scheduling
- **Action Metrics** - Performance metrics and monitoring
- **Web Interface** - Web-based action management interface

### 📋 Future Enhancements

- **Action Analytics** - Advanced action analysis and reporting
- **Custom Action Types** - User-defined action types
- **Performance Testing** - Benchmarks and optimization
- **Comprehensive Testing** - Integration and performance tests

## Getting Help

### Documentation Resources

1. **User Guide** - For usage questions and best practices
2. **API Reference** - For technical implementation details
3. **Troubleshooting Guide** - For problem resolution
4. **Examples** - For configuration patterns and use cases

### Common Questions

#### "How do I get started?"
Start with the [User Guide](ACTIONS_USER_GUIDE.md) and copy an example from the [examples/](examples/) directory.

#### "How do I configure actions?"
See the [Action Configuration](ACTIONS_USER_GUIDE.md#creating-actions) section and the [basic configuration example](examples/basic-actions-project.hcl).

#### "How do I troubleshoot action issues?"
Check the [Action Loading Errors](ACTIONS_TROUBLESHOOTING.md#action-loading-errors) section in the troubleshooting guide.

#### "How do I validate my configuration?"
Use `spooky actions validate` and see the [Action Validation](ACTIONS_USER_GUIDE.md#validate-actions) section.

#### "How do I integrate with other systems?"
Review the [API Reference](ACTIONS_API_REFERENCE.md) for integration patterns and the planned features section above.

### Contributing

When contributing to the actions system:

1. **Follow Interface Patterns** - Use the established interface architecture
2. **Add Comprehensive Tests** - Include unit and integration tests
3. **Update Documentation** - Keep documentation current with changes
4. **Follow Error Handling** - Use structured error types and patterns
5. **Consider Performance** - Optimize for efficient action running

## Conclusion

The spooky actions system provides a comprehensive solution for action orchestration across all spooky components. The documentation is structured to support users at all levels, from beginners getting started to advanced users implementing complex integrations.

Start with the User Guide to understand the basics, use the examples as templates for your configurations, and refer to the troubleshooting guide when you encounter issues. The API reference provides the technical details needed for development and integration work.

The system is designed to be extensible and maintainable, following interface-first architecture principles and comprehensive configuration patterns. As the system evolves, new features will be added while maintaining backward compatibility and following established patterns.

For the most up-to-date information and examples, always refer to the latest version of the documentation and test your configurations with the current spooky release.

## Integration with Other Systems

### Project Integration

Actions integrate with the project system:

- **Machine Inventory**: Uses machines.hcl for target identification
- **Project Structure**: Follows project directory structure
- **Configuration**: Uses project-specific configuration

### Variables Integration

Actions can use variables for dynamic configuration:

```hcl
action "deploy-with-vars" {
  type = "template_deploy"
  
  template {
    source = "templates/app.conf.tmpl"
    destination = "/etc/app/app.conf"
  }
  
  environment = {
    APP_VERSION = "${variables.app_version}"
    ENVIRONMENT = "${variables.environment}"
  }
  
  machines = ["${variables.target_machines}"]
}
```

### Facts Integration

Actions can use machine facts for conditional running:

```hcl
action "os-specific-command" {
  type = "command"
  
  command = "if [ \"${facts.system.os.name}\" = \"Ubuntu\" ]; then apt update; else yum update; fi"
  
  machines = ["web-server"]
}
```

### SSH Integration

Actions use SSH for remote running:

- **SSH Connections**: Secure SSH connections to target machines
- **Authentication**: SSH key and password authentication
- **Connection Pooling**: Efficient SSH connection management
- **Error Handling**: Robust SSH error handling and recovery

### Logging Integration

Actions integrate with the logging system:

- **Component Logging**: Actions use the "actions" logging component
- **Run Logging**: Log action run details and results
- **Error Logging**: Log action errors and failures
- **Performance Logging**: Log action run performance metrics

## Best Practices

### Action Design

1. **Use Descriptive Names** - Choose clear, descriptive action names
2. **Add Descriptions** - Include human-readable descriptions
3. **Group Related Actions** - Use dependencies to group related operations
4. **Use Tags for Organization** - Use tags to organize actions by purpose
5. **Keep Actions Focused** - Each action should have a single, clear purpose

### Configuration Management

1. **Use Template Variables** - Use variables for configuration values
2. **Validate Templates** - Enable template validation for critical files
3. **Use Backup Options** - Enable backups for file operations
4. **Set Appropriate Permissions** - Use correct file permissions and ownership
5. **Use Working Directories** - Set appropriate working directories

### Error Handling

1. **Configure Retries** - Use retries for unreliable operations
2. **Set Timeouts** - Use appropriate timeouts for long-running operations
3. **Handle Dependencies** - Use dependencies to ensure proper order
4. **Validate Before Running** - Always validate actions before running
5. **Use Dry Run Mode** - Test actions with dry-run mode first

### Performance Optimization

1. **Use Parallel Running** - Enable parallel running when possible
2. **Limit Concurrent Operations** - Use max_concurrent to prevent overload
3. **Optimize Dependencies** - Minimize unnecessary dependencies
4. **Use Resource Limits** - Set appropriate resource limits
5. **Monitor Running** - Monitor action run performance

### Security Considerations

1. **Use Sudo Sparingly** - Only use sudo when necessary
2. **Set File Permissions** - Use appropriate file permissions
3. **Validate Input** - Validate all input and configuration
4. **Use Secure Connections** - Ensure SSH connections are secure
5. **Audit Actions** - Regularly audit action configurations

## Examples

### Web Application Deployment

```hcl
# actions.hcl
action "stop-application" {
  type = "service_control"
  description = "Stop application service"
  
  service_control {
    service = "myapp"
    action = "stop"
  }
  
  machines = ["app-server"]
}

action "backup-database" {
  type = "script"
  description = "Create database backup"
  
  script = "files/backup-db.sh"
  variables = {
    backup_dir = "/backups"
    db_name = "myapp"
  }
  
  machines = ["db-server"]
}

action "deploy-application" {
  type = "template_deploy"
  description = "Deploy application files"
  
  template {
    source = "templates/app.conf.tmpl"
    destination = "/etc/myapp/app.conf"
    permissions = "0644"
    owner = "myapp"
    group = "myapp"
    backup = true
  }
  
  machines = ["app-server"]
  dependencies = ["stop-application", "backup-database"]
}

action "start-application" {
  type = "service_control"
  description = "Start application service"
  
  service_control {
    service = "myapp"
    action = "start"
    wait_for_status = "active"
  }
  
  machines = ["app-server"]
  dependencies = ["deploy-application"]
}
```

### Infrastructure Maintenance

```hcl
action "update-packages" {
  type = "command"
  description = "Update system packages"
  
  command = "apt update && apt upgrade -y"
  
  machines = ["web-server", "app-server", "db-server"]
  parallel = true
  sudo = true
  timeout = 600
}

action "restart-services" {
  type = "service_control"
  description = "Restart critical services"
  
  service_control {
    service = "nginx"
    action = "restart"
    wait_for_status = "active"
  }
  
  machines = ["web-server"]
  dependencies = ["update-packages"]
}
```

### Monitoring and Health Checks

```hcl
action "check-disk-space" {
  type = "command"
  description = "Check available disk space"
  
  command = "df -h"
  
  machines = ["web-server", "app-server", "db-server"]
  parallel = true
}

action "check-service-status" {
  type = "service_control"
  description = "Check service status"
  
  service_control {
    service = "nginx"
    action = "status"
  }
  
  machines = ["web-server"]
}
```

This comprehensive documentation summary provides an overview of the spooky actions system documentation and helps users find the right information for their needs.
