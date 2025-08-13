# Actions System User Guide

## Overview

The spooky actions system provides comprehensive action orchestration capabilities for running commands, scripts, templates, and service controls on remote machines. This guide covers everything from basic action configuration to advanced features like dependency management, parallel execution, and service control.

## Getting Started

### Prerequisites

- spooky CLI installed and configured
- SSH access to target machines
- Basic understanding of HCL configuration syntax
- Access to create and modify project files

### Quick Start

1. **Check Available Actions Commands**
   ```bash
   spooky actions --help
   ```

2. **List Actions in a Project**
   ```bash
   spooky actions list ./my-project
   ```

3. **Validate Actions Configuration**
   ```bash
   spooky actions validate ./my-project
   ```

4. **Run Actions**
   ```bash
   spooky actions run ./my-project
   ```

## Actions System Concepts

### What are Actions?

Actions are operations that can be performed on machines, such as:

- **Running commands** - Execute shell commands on remote machines
- **Running scripts** - Execute script files with template processing
- **Deploying templates** - Render and deploy template files
- **Copying files** - Transfer files between machines
- **Controlling services** - Start, stop, restart system services

### Action Types

The actions system supports several types of actions:

1. **command** - Execute shell commands
2. **script** - Run script files (with template support)
3. **template_deploy** - Deploy template files
4. **file_copy** - Copy files between machines
5. **service_control** - Control system services

### Action Lifecycle

The action execution process follows these steps:

1. **Action Loading** - Load actions from `actions.hcl` and `actions/` directory
2. **Dependency Resolution** - Resolve action dependencies and create execution order
3. **Machine Targeting** - Determine target machines based on configuration
4. **Execution Planning** - Create execution plan with `--plan` flag
5. **Action Execution** - Execute actions on target machines via SSH
6. **Result Collection** - Collect and report execution results

## Basic Usage

### Creating Actions

Actions are defined in HCL configuration files. You can use either a single `actions.hcl` file or multiple files in an `actions/` directory.

#### Single File Configuration

```hcl
# actions.hcl
action "restart-nginx" {
  type = "service_control"
  description = "Restart nginx service"
  
  service_control {
    service = "nginx"
    action = "restart"
    wait_for_status = "active"
  }
  
  machines = ["web-server-1", "web-server-2"]
  tags = ["web", "production"]
}

action "deploy-config" {
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
  dependencies = ["restart-nginx"]
}
```

#### Directory Configuration

```bash
my-project/
├── actions.hcl              # Main actions file
└── actions/                 # Additional action files
    ├── web-actions.hcl      # Web server actions
    ├── db-actions.hcl       # Database actions
    └── deploy-actions.hcl   # Deployment actions
```

### Action Configuration

#### Basic Action Structure

```hcl
action "action-name" {
  type = "command"                    # Required: action type
  description = "Action description"  # Optional: human-readable description
  
  # Action-specific configuration
  command = "echo 'Hello World'"
  
  # Targeting configuration
  machines = ["machine1", "machine2"]  # Target specific machines
  tags = ["web", "production"]         # Target machines by tags
  
  # Execution configuration
  parallel = false                     # Run sequentially (default)
  timeout = 300                        # Timeout in seconds
  retries = 3                          # Number of retries
  retry_delay = 10                     # Delay between retries
  
  # Dependencies
  dependencies = ["other-action"]      # Actions that must complete first
  
  # Environment and context
  environment = {
    NODE_ENV = "production"
    DEBUG = "false"
  }
  
  working_directory = "/opt/app"
  user = "app"
  sudo = false
}
```

#### Command Actions

Execute shell commands on remote machines:

```hcl
action "check-disk-space" {
  type = "command"
  description = "Check available disk space"
  
  command = "df -h /"
  
  machines = ["web-server", "db-server"]
  parallel = true
}

action "update-packages" {
  type = "command"
  description = "Update system packages"
  
  command = "apt update && apt upgrade -y"
  
  machines = ["web-server"]
  sudo = true
  timeout = 600
}
```

#### Script Actions

Execute script files with optional template processing:

```hcl
action "backup-database" {
  type = "script"
  description = "Create database backup"
  
  script = "files/backup-db.sh"
  
  variables = {
    backup_dir = "/backups"
    db_name = "myapp"
    retention_days = "7"
  }
  
  machines = ["db-server"]
  working_directory = "/opt/scripts"
}
```

#### Template Deploy Actions

Deploy template files with variable substitution:

```hcl
action "deploy-nginx-config" {
  type = "template_deploy"
  description = "Deploy nginx configuration"
  
  template {
    source = "templates/nginx.conf.tmpl"
    destination = "/etc/nginx/nginx.conf"
    permissions = "0644"
    owner = "root"
    group = "root"
    validate = true
    backup = true
  }
  
  machines = ["web-server"]
  dependencies = ["restart-nginx"]
}
```

#### File Copy Actions

Copy files between machines:

```hcl
action "deploy-certificates" {
  type = "file_copy"
  description = "Deploy SSL certificates"
  
  file_copy {
    source = "files/certs/ssl.crt"
    destination = "/etc/ssl/certs/app.crt"
    permissions = "0644"
    owner = "root"
    group = "root"
    backup = true
  }
  
  machines = ["web-server"]
}
```

#### Service Control Actions

Control system services:

```hcl
action "restart-application" {
  type = "service_control"
  description = "Restart application service"
  
  service_control {
    service = "myapp"
    action = "restart"
    systemd = true
    timeout = 30
    wait_for_status = "active"
    wait_timeout = 60
  }
  
  machines = ["app-server"]
}
```

### Action Targeting

#### Machine Targeting

Target specific machines by name:

```hcl
action "web-deploy" {
  type = "command"
  command = "echo 'Deploying to web servers'"
  
  machines = ["web-1", "web-2", "web-3"]
}
```

#### Tag-Based Targeting

Target machines by tags:

```hcl
action "prod-deploy" {
  type = "command"
  command = "echo 'Deploying to production'"
  
  tags = ["environment=production", "role=web"]
}
```

#### Combined Targeting

Combine machine names and tags:

```hcl
action "full-deploy" {
  type = "command"
  command = "echo 'Full deployment'"
  
  machines = ["web-1", "web-2"]
  tags = ["environment=production"]
}
```

### Action Dependencies

Define dependencies between actions:

```hcl
action "prepare-database" {
  type = "command"
  command = "echo 'Preparing database'"
  machines = ["db-server"]
}

action "deploy-application" {
  type = "command"
  command = "echo 'Deploying application'"
  machines = ["app-server"]
  dependencies = ["prepare-database"]  # Must complete first
}

action "update-load-balancer" {
  type = "command"
  command = "echo 'Updating load balancer'"
  machines = ["lb-server"]
  dependencies = ["deploy-application"]  # Must complete first
}
```

## Advanced Usage

### Parallel Execution

Run actions in parallel across machines:

```hcl
action "parallel-update" {
  type = "command"
  command = "apt update"
  
  machines = ["web-1", "web-2", "web-3", "web-4"]
  parallel = true  # Run on all machines simultaneously
  max_concurrent = 4  # Limit concurrent executions
}
```

### Retry Logic

Configure retry behavior for unreliable operations:

```hcl
action "unreliable-operation" {
  type = "command"
  command = "curl -f https://api.example.com/health"
  
  machines = ["app-server"]
  retries = 5
  retry_delay = 30  # Wait 30 seconds between retries
  timeout = 60
}
```

### Environment Variables

Set environment variables for action execution:

```hcl
action "deploy-with-env" {
  type = "script"
  script = "files/deploy.sh"
  
  environment = {
    NODE_ENV = "production"
    DEBUG = "false"
    API_URL = "https://api.production.com"
    DB_HOST = "db.production.com"
  }
  
  machines = ["app-server"]
}
```

### Resource Limits

Set resource limits for action execution:

```hcl
action "resource-intensive" {
  type = "command"
  command = "build-application"
  
  resource_limits {
    memory_mb = 2048
    cpu_percent = 50
    disk_mb = 1024
  }
  
  machines = ["build-server"]
}
```

### Dry Run Mode

Test actions without actually executing them:

```bash
# Show what would be executed
spooky actions run ./my-project --dry-run

# Show execution plan
spooky actions run ./my-project --plan
```

## CLI Commands

### List Actions

```bash
# List all actions in a project
spooky actions list ./my-project

# List actions with details
spooky actions list ./my-project --verbose

# List actions for specific machines
spooky actions list ./my-project --machine web-server

# List actions by tags
spooky actions list ./my-project --tags "environment=production"
```

### Validate Actions

```bash
# Validate action configuration
spooky actions validate ./my-project

# Validate with detailed output
spooky actions validate ./my-project --verbose

# Validate specific actions
spooky actions validate ./my-project --action restart-nginx
```

### Run Actions

```bash
# Run all actions
spooky actions run ./my-project

# Run specific actions
spooky actions run ./my-project --action restart-nginx

# Run actions with plan mode
spooky actions run ./my-project --plan

# Run actions in dry-run mode
spooky actions run ./my-project --dry-run

# Run actions in parallel
spooky actions run ./my-project --parallel 4

# Run actions with timeout
spooky actions run ./my-project --timeout 600
```

### Command Options

Common options across actions commands:

- `--action`: Target specific action
- `--machine`: Target specific machine
- `--tags`: Filter by tags
- `--parallel`: Number of parallel workers
- `--timeout`: Action timeout in seconds
- `--verbose`: Verbose output
- `--plan`: Show execution plan without running
- `--dry-run`: Simulate execution without running

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
4. **Validate Before Running** - Always validate actions before execution
5. **Use Dry Run Mode** - Test actions with dry-run mode first

### Performance Optimization

1. **Use Parallel Execution** - Enable parallel execution when possible
2. **Limit Concurrent Operations** - Use max_concurrent to prevent overload
3. **Optimize Dependencies** - Minimize unnecessary dependencies
4. **Use Resource Limits** - Set appropriate resource limits
5. **Monitor Execution** - Monitor action execution performance

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

action "update-load-balancer" {
  type = "command"
  description = "Update load balancer configuration"
  
  command = "systemctl reload haproxy"
  
  machines = ["lb-server"]
  dependencies = ["start-application"]
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

action "cleanup-logs" {
  type = "command"
  description = "Clean up old log files"
  
  command = "find /var/log -name '*.log' -mtime +30 -delete"
  
  machines = ["web-server", "app-server", "db-server"]
  parallel = true
  sudo = true
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

action "check-application-health" {
  type = "command"
  description = "Check application health endpoint"
  
  command = "curl -f http://localhost:8080/health"
  
  machines = ["app-server"]
  retries = 3
  retry_delay = 10
}
```

## Troubleshooting

### Common Issues

#### Action Validation Failures

```bash
# Check action configuration
spooky actions validate ./my-project --verbose

# Check for syntax errors
spooky actions validate ./my-project --strict
```

#### Execution Failures

```bash
# Check SSH connectivity
spooky machines ping ./my-project

# Test specific action
spooky actions run ./my-project --action test-action --verbose

# Check action dependencies
spooky actions run ./my-project --plan
```

#### Performance Issues

```bash
# Reduce parallel workers
spooky actions run ./my-project --parallel 2

# Increase timeout
spooky actions run ./my-project --timeout 1200

# Use dry-run to test
spooky actions run ./my-project --dry-run
```

### Debug Information

When troubleshooting action issues, enable verbose output:

```bash
# Enable verbose output
spooky actions run ./my-project --verbose

# Check action plan
spooky actions run ./my-project --plan --verbose

# Validate with details
spooky actions validate ./my-project --verbose
```

## Integration with Other Systems

### Project Integration

Actions integrate with the project system:

- **Machine Inventory**: Uses machines.hcl for target identification
- **Project Structure**: Follows project directory structure
- **Configuration**: Uses project-specific configuration

### Variables Integration

Use variables in action configurations:

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

Use machine facts in actions:

```hcl
action "os-specific-command" {
  type = "command"
  
  command = "if [ \"${facts.system.os.name}\" = \"Ubuntu\" ]; then apt update; else yum update; fi"
  
  machines = ["web-server"]
}
```

## Current Implementation Status

### ✅ Implemented Features

- **Basic Action Types**: Command, script, template deploy, file copy, service control
- **Action Configuration**: HCL-based configuration with validation
- **Dependency Management**: Action dependency resolution and execution order
- **Machine Targeting**: Support for machine names and tags
- **Parallel Execution**: Parallel action execution across machines
- **CLI Integration**: Complete CLI command set
- **Validation**: Action configuration validation
- **Planning Mode**: Execution planning with `--plan` flag
- **Dry Run Mode**: Simulation mode with `--dry-run` flag

### 📋 Future Enhancements

- **Advanced Scripting**: Enhanced script execution with better error handling
- **Template Functions**: Additional template functions and helpers
- **Action History**: Track action execution history
- **Rollback Support**: Automatic rollback capabilities
- **Advanced Scheduling**: Time-based action scheduling
- **Action Metrics**: Performance metrics and monitoring
- **Web Interface**: Web-based action management interface

## Conclusion

The spooky actions system provides comprehensive action orchestration capabilities for managing remote machine operations. The system is designed to be flexible, reliable, and easy to use while supporting complex deployment scenarios.

Start with simple actions and gradually build more complex workflows. Always validate your configurations and use dry-run mode to test before actual execution.

For more advanced usage and troubleshooting, refer to the [API Reference](ACTIONS_API_REFERENCE.md) and [Troubleshooting Guide](ACTIONS_TROUBLESHOOTING.md).
