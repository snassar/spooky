# Actions System User Guide

## Overview

The spooky actions system provides comprehensive action orchestration capabilities for running commands, scripts, templates, and service controls on remote machines. This guide covers everything from basic action configuration to advanced features like dependency management, parallel running, and service control.

**Status: Partially Implemented** - The actions system has basic functionality but SSH-based action orchestration has known issues that need to be addressed.

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

- **Running commands** - Run shell commands on remote machines
- **Running scripts** - Run script files with template processing
- **Deploying templates** - Render and deploy template files
- **Copying files** - Transfer files between machines
- **Controlling services** - Start, stop, restart system services

### Action Types

The actions system supports several types of actions:

1. **command** - Run shell commands
2. **script** - Run script files (with template support)
3. **template_deploy** - Deploy template files
4. **file_copy** - Copy files between machines
5. **service_control** - Control system services

### Action Lifecycle

The action run process follows these steps:

1. **Action Loading** - Load actions from `actions.hcl` and `actions/` directory
2. **Dependency Resolution** - Resolve action dependencies and create run order
3. **Machine Targeting** - Determine target machines based on configuration
4. **Run Planning** - Create run plan with `--plan` flag
5. **Action Running** - Run actions on target machines via SSH
6. **Result Collection** - Collect and aggregate results from all machines
7. **Error Handling** - Handle failures and provide detailed error reporting

## Current Implementation Status

### ✅ Working Features

- **Basic Action Loading**: Action loading from HCL configuration files
- **Action Validation**: Basic action validation and error handling
- **CLI Integration**: `spooky actions` commands with basic functionality
- **Project Integration**: Actions loading from project configuration
- **Local Action Orchestration**: Local action orchestration capabilities

### ⚠️ Known Issues

- **SSH-Based Orchestration**: SSH-based action orchestration has implementation issues
- **Remote Action Running**: Cannot properly run actions on remote machines
- **Parallel Execution**: No parallel action execution support
- **Action Planning**: Limited action planning capabilities
- **Result Aggregation**: No action result aggregation

### 🔧 Current Workarounds

- **Local Running**: Use local action running for immediate needs
- **Manual Execution**: Run actions manually on remote machines if needed
- **Filtered Running**: Use machine filtering to limit running scope
- **Monitor Updates**: Watch for improvements to SSH-based orchestration

## Basic Usage

### Listing Actions

List all available actions in a project:

```bash
# List all actions
spooky actions list ./my-project
```

**Example Output**:
```
Available actions (3 found):
1. update-system - Update system packages
2. deploy-app - Deploy application
3. restart-services - Restart system services
```

### Validating Actions

Validate action configurations before running:

```bash
# Validate all actions
spooky actions validate ./my-project
```

**Example Output**:
```
✅ All actions are valid
```

### Running Actions

Run actions on target machines:

```bash
# Run all actions
spooky actions run ./my-project

# Run with dry-run mode
spooky actions run ./my-project --dry-run

# Run with plan mode
spooky actions run ./my-project --plan

# Run with parallel execution
spooky actions run ./my-project --parallel 4
```

### Machine Targeting

Target specific machines for action running:

```bash
# Run actions on specific machines
spooky actions run ./my-project --machine web-server

# Run on multiple machines
spooky actions run ./my-project --machine web-server --machine db-server

# Run on machines with specific tags
spooky actions run ./my-project --tags environment=production

# Run with complex filtering
spooky actions run ./my-project --filter "environment=production AND role=web"
```

## Project Configuration

### Basic Action Configuration

Create actions in your project:

```hcl
# actions.hcl
actions {
  action "update-system" {
    description = "Update system packages"
    type = "command"
    
    command = "apt update && apt upgrade -y"
    
    machines = ["web-server", "db-server"]
    
    tags = ["maintenance"]
  }
  
  action "deploy-app" {
    description = "Deploy application"
    type = "script"
    
    script {
      source = "scripts/deploy.sh"
      arguments = ["--version", "1.2.3"]
    }
    
    machines = ["web-server"]
    
    dependencies = ["update-system"]
  }
  
  action "restart-services" {
    description = "Restart system services"
    type = "service_control"
    
    service_control {
      service = "nginx"
      action = "restart"
    }
    
    machines = ["web-server"]
    
    dependencies = ["deploy-app"]
  }
}
```

### Action Types Configuration

#### Command Actions

Run shell commands on remote machines:

```hcl
action "check-disk" {
  type = "command"
  
  command = "df -h"
  
  machines = ["web-server"]
  
  timeout = 30
  retries = 3
}
```

#### Script Actions

Run script files with template support:

```hcl
action "deploy-config" {
  type = "script"
  
  script {
    source = "scripts/deploy-config.sh"
    arguments = ["--env", "production"]
    environment = {
      APP_VERSION = "1.2.3"
      ENVIRONMENT = "production"
    }
  }
  
  machines = ["web-server"]
}
```

#### Template Deploy Actions

Deploy rendered templates to remote machines:

```hcl
action "deploy-nginx-config" {
  type = "template_deploy"
  
  template {
    source = "templates/nginx.conf.tmpl"
    destination = "/etc/nginx/nginx.conf"
    permissions = "0644"
  }
  
  machines = ["web-server"]
}
```

#### File Copy Actions

Copy files between machines:

```hcl
action "copy-ssl-cert" {
  type = "file_copy"
  
  file_copy {
    source = "certs/ssl.crt"
    destination = "/etc/ssl/certs/ssl.crt"
    permissions = "0644"
  }
  
  machines = ["web-server"]
}
```

#### Service Control Actions

Control system services:

```hcl
action "restart-nginx" {
  type = "service_control"
  
  service_control {
    service = "nginx"
    action = "restart"
  }
  
  machines = ["web-server"]
}
```

### Advanced Configuration

#### Dependencies

Define action dependencies for proper run order:

```hcl
action "deploy-app" {
  type = "script"
  
  script {
    source = "scripts/deploy.sh"
  }
  
  machines = ["web-server"]
  
  dependencies = ["update-system", "deploy-config"]
}
```

#### Parallel Execution

Configure parallel execution for actions:

```hcl
action "update-all-servers" {
  type = "command"
  
  command = "apt update && apt upgrade -y"
  
  machines = ["web-server", "db-server", "cache-server"]
  
  parallel = true
  max_concurrent = 3
}
```

#### Resource Limits

Set resource limits for actions:

```hcl
action "heavy-process" {
  type = "script"
  
  script {
    source = "scripts/heavy-process.sh"
  }
  
  machines = ["web-server"]
  
  resource_limits {
    memory_mb = 1024
    cpu_percent = 50
    timeout_seconds = 300
  }
}
```

#### Error Handling

Configure error handling behavior:

```hcl
action "critical-update" {
  type = "command"
  
  command = "critical-system-update"
  
  machines = ["web-server"]
  
  allow_failure = false
  stop_on_failure = true
  retries = 3
  retry_delay = 30
}
```

## Running Modes

### Normal Running

Run actions normally on target machines:

```bash
spooky actions run ./my-project
```

### Dry Run Mode

Simulate running without making changes:

```bash
spooky actions run ./my-project --dry-run
```

**Example Output**:
```
🔍 Dry run mode - no changes will be made

Would run the following actions:
1. update-system on web-server, db-server
2. deploy-app on web-server
3. restart-services on web-server

Total: 3 actions across 2 machines
```

### Plan Mode

Show running plan without running:

```bash
spooky actions run ./my-project --plan
```

**Example Output**:
```
📋 Running Plan:

Phase 1: update-system
  - web-server: apt update && apt upgrade -y
  - db-server: apt update && apt upgrade -y

Phase 2: deploy-app (depends on: update-system)
  - web-server: scripts/deploy.sh --version 1.2.3

Phase 3: restart-services (depends on: deploy-app)
  - web-server: systemctl restart nginx

Dependencies resolved: ✅
Target machines: web-server, db-server
Estimated time: 5-10 minutes
```

### Parallel Running

Run actions in parallel for better performance:

```bash
spooky actions run ./my-project --parallel 4
```

## Machine Targeting

### Machine Names

Target specific machines by name:

```bash
# Single machine
spooky actions run ./my-project --machine web-server

# Multiple machines
spooky actions run ./my-project --machine web-server --machine db-server
```

### Tag-Based Targeting

Target machines using tags:

```bash
# Single tag
spooky actions run ./my-project --tags environment=production

# Multiple tags (AND logic)
spooky actions run ./my-project --tags environment=production --tags role=web

# Multiple tag values (OR logic)
spooky actions run ./my-project --tags environment=production --tags environment=staging
```

### Complex Filtering

Use complex filter expressions:

```bash
# Environment and role filter
spooky actions run ./my-project --filter "environment=production AND role=web"

# Multiple conditions
spooky actions run ./my-project --filter "environment=production AND (role=web OR role=api)"

# Pattern matching
spooky actions run ./my-project --filter "hostname LIKE 'web-%'"
```

## Troubleshooting

### Common Issues

#### SSH Connection Failures

**Problem**: SSH connections fail during action running

**Solutions**:
1. Verify SSH connectivity manually:
   ```bash
   ssh admin@web.example.com
   ```

2. Check SSH key permissions:
   ```bash
   chmod 600 ~/.ssh/id_ed25519
   ```

3. Verify machine inventory configuration:
   ```bash
   spooky machines validate ./my-project
   ```

#### Action Validation Errors

**Problem**: Actions fail validation

**Solutions**:
1. Check action configuration syntax:
   ```bash
   spooky actions validate ./my-project
   ```

2. Verify required fields are present:
   ```hcl
   action "my-action" {
     name = "my-action"        # Required
     type = "command"          # Required
     command = "echo hello"    # Required for command type
     machines = ["web-server"] # Required
   }
   ```

3. Check for circular dependencies:
   ```bash
   spooky actions validate ./my-project --check-dependencies
   ```

#### Dependency Resolution Issues

**Problem**: Action dependencies cannot be resolved

**Solutions**:
1. Check dependency names match action names exactly
2. Verify no circular dependencies exist
3. Use plan mode to see dependency resolution:
   ```bash
   spooky actions run ./my-project --plan
   ```

### Debugging

Enable verbose output for debugging:

```bash
# Enable debug logging
export SPOOKY_LOG_LEVEL=debug

# Run with verbose output
spooky actions run ./my-project --dry-run
```

## Best Practices

### Action Design

1. **Use Descriptive Names**: Choose clear, descriptive action names
2. **Add Descriptions**: Include descriptions for complex actions
3. **Group Related Actions**: Use tags to group related actions
4. **Set Appropriate Timeouts**: Configure timeouts based on action complexity
5. **Handle Errors Gracefully**: Use retries and error handling appropriately

### Security Considerations

1. **SSH Key Management**: Use dedicated SSH keys for action running
2. **Key Permissions**: Ensure SSH keys have correct permissions (600)
3. **Network Security**: Use VPN or secure networks for action running
4. **Sensitive Commands**: Be careful with commands that handle sensitive data

### Performance Optimization

1. **Use Parallel Running**: Use parallel execution for independent actions
2. **Optimize Dependencies**: Minimize unnecessary dependencies
3. **Target Appropriately**: Use machine filtering to limit scope
4. **Monitor Resources**: Watch for resource usage during running

## Integration with Other Systems

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

## CLI Reference

### `spooky actions list`

List available actions in a project.

**Syntax**:
```bash
spooky actions list <project-path>
```

**Examples**:
```bash
# List all actions
spooky actions list ./my-project
```

### `spooky actions validate`

Validate actions in a project.

**Syntax**:
```bash
spooky actions validate <project-path>
```

**Examples**:
```bash
# Validate all actions
spooky actions validate ./my-project
```

### `spooky actions run`

Run actions on machines.

**Syntax**:
```bash
spooky actions run <project-path> [options]
```

**Options**:
- `--machine <list>` - Target specific machines
- `--tags <list>` - Target machines by tags
- `--filter <query>` - Use complex filter query
- `--parallel <number>` - Number of parallel workers (minimum 2)
- `--dry-run` - Simulate running without making changes
- `--plan` - Show running plan without running
- `--decrypt` - Decrypt encrypted variables and facts in-memory for debugging

**Examples**:
```bash
# Run all actions
spooky actions run ./my-project

# Run with dry-run mode
spooky actions run ./my-project --dry-run

# Run with plan mode
spooky actions run ./my-project --plan

# Run with parallel execution
spooky actions run ./my-project --parallel 4

# Run on specific machines
spooky actions run ./my-project --machine web-server

# Run with decryption
spooky actions run ./my-project --decrypt
```

## Remember

**Good actions system usage enables:**
- Efficient automation of machine operations
- Consistent deployment and configuration
- Scalable infrastructure management
- Integration with other spooky systems
- Reliable and repeatable operations

**Always be aware of current SSH-based orchestration limitations and use appropriate workarounds until these issues are resolved.**
