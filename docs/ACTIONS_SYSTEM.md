# Actions System

## Overview

The spooky actions system provides comprehensive action management and orchestration capabilities. Actions are defined in `actions.hcl` files within spooky projects and can be run on target machines with support for parallel orchestration, dependency management, and conditional running.

## Related Systems

This system integrates with and depends on several other spooky systems:

- **[Facts System](FACTS_SYSTEM.md)** - Actions use facts for context and decision-making
- **[Machines System](MACHINES_SYSTEM.md)** - Actions run on machines defined in the inventory
- **[Templates System](TEMPLATES_SYSTEM.md)** - Actions use templates for rendering dynamic content
- **[Variables System](VARIABLES_SYSTEM.md)** - Actions use variables for dynamic configuration
- **[SSH System](SSH_SYSTEM.md)** - Actions use SSH for remote machine connectivity and command running
- **[Logging System](LOGGING_SYSTEM.md)** - Actions generate comprehensive logs for monitoring and debugging
- **[Integrations System](INTEGRATIONS_SYSTEM.md)** - Actions integrate with other systems through the IntegrationManager
- **[Projects System](PROJECTS_SYSTEM.md)** - Actions are organized within projects

## Core Concepts

### Action Definition
Actions are defined in `actions.hcl` files that specify:

- **Action Metadata**: Name, description, and version information
- **Running Commands**: Scripts or commands to run
- **Target Machines**: Which machines to run the action on
- **Dependencies**: Other actions that must complete first
- **Conditions**: When the action should run
- **Configuration**: Timeouts, retry settings, and other parameters

### Action Orchestration
The system supports multiple orchestration modes:

- **Sequential Orchestration**: Run actions one after another
- **Parallel Orchestration**: Run actions simultaneously across machines
- **Dependency-based Orchestration**: Run actions based on dependencies
- **Conditional Orchestration**: Run actions based on conditions

### Action Validation
Actions are validated before running:

- **Syntax Validation**: Ensures valid HCL syntax
- **Schema Validation**: Validates against action schema
- **Dependency Validation**: Validates action dependencies
- **Machine Validation**: Validates target machine availability

## CLI Commands

### Action Management Commands

#### `spooky actions list [project-path]`
List all available actions in the project.

**Examples:**
```bash
# List all actions in project
spooky actions list ./my-project
```

#### `spooky actions run [project-path]`
Run actions on target machines.

**Flags:**
- `--action` - Run specific action by name
- `--machine` - Target specific machines
- `--tags` - Target machines with specific tags
- `--filter` - Complex filter expression
- `--parallel` - Number of parallel workers (minimum 2)
- `--plan` - Show running plan without running
- `--dry-run` - Simulate running without making changes
- `--decrypt` - Decrypt encrypted variables and facts in-memory for debugging

**Examples:**
```bash
# Run all actions
spooky actions run ./my-project

# Run specific action
spooky actions run ./my-project --action deploy

# Run with dry-run
spooky actions run ./my-project --action deploy --dry-run

# Run with plan
spooky actions run ./my-project --action deploy --plan

# Run with parallel orchestration
spooky actions run ./my-project --action deploy --parallel 4

# Run on specific machines
spooky actions run ./my-project --action deploy --machine web-server

# Run on machines with tags
spooky actions run ./my-project --action deploy --tags production
```

#### `spooky actions validate [project-path]`
Validate action definitions and dependencies.

**Examples:**
```bash
# Validate all actions
spooky actions validate ./my-project

# Validate specific action
spooky actions validate ./my-project --action deploy
```

## Action Configuration

### Basic Action Configuration
```hcl
# actions.hcl
actions {
  action "deploy" {
    description = "Deploy application to target machines"
    
    machines = ["web-server", "db-server"]
    parallel = true
    
    command = "deploy.sh"
    timeout_seconds = 300
    
    tags = ["deployment", "production"]
  }
}
```

### Advanced Action Configuration
```hcl
# actions.hcl
actions {
  action "system-update" {
    description = "Update system packages"
    
    machines = ["web-server", "db-server"]
    parallel = true
    
    command = "apt update && apt upgrade -y"
    timeout_seconds = 600
    
    dependencies = ["backup"]
    
    conditions {
      maintenance_window = true
      system_load < 80
    }
    
    retry {
      attempts = 3
      delay_seconds = 30
    }
    
    tags = ["maintenance", "system"]
  }
  
  action "backup" {
    description = "Create system backup"
    
    machines = ["db-server"]
    parallel = false
    
    command = "backup.sh"
    timeout_seconds = 1800
    
    tags = ["backup", "critical"]
  }
}
```

### Action Dependencies
Actions can depend on other actions:

```hcl
actions {
  action "deploy" {
    description = "Deploy application"
    command = "deploy.sh"
    dependencies = ["backup", "test"]
  }
  
  action "backup" {
    description = "Create backup"
    command = "backup.sh"
  }
  
  action "test" {
    description = "Run tests"
    command = "test.sh"
    dependencies = ["backup"]
  }
}
```

### Conditional Actions
Actions can have conditions for execution:

```hcl
actions {
  action "deploy" {
    description = "Deploy application"
    command = "deploy.sh"
    
    conditions {
      environment = "production"
      maintenance_window = true
      system_load < 80
    }
  }
}
```

## Action Running Process

### Running Workflow
1. **Load Actions**: Load action definitions from `actions.hcl`
2. **Validate Actions**: Validate syntax, schema, and dependencies
3. **Resolve Dependencies**: Determine running order
4. **Target Machines**: Identify target machines for each action
5. **Run Actions**: Run actions on target machines
6. **Monitor Progress**: Track running progress and results
7. **Handle Errors**: Manage errors and retries
8. **Report Results**: Provide running results and summaries

### Parallel Running
```bash
# Run actions in parallel
spooky actions run ./my-project --parallel 4
```

### Sequential Running
```bash
# Run actions sequentially
spooky actions run ./my-project --parallel 1
```

## Action Validation

### Validation Process
The action validation system performs:

1. **Syntax Validation**: Validates HCL syntax
2. **Schema Validation**: Validates against action schema
3. **Dependency Validation**: Validates action dependencies
4. **Machine Validation**: Validates target machine availability
5. **Command Validation**: Validates command syntax and availability

### Validation Output
```bash
# Validate actions
spooky actions validate ./my-project
```

**Output:**
```
🔍 Validating actions: ./my-project
✅ Action validation passed - all actions valid
📋 Schema compliance: actions.schema.hcl ✅
📋 Dependency validation: Valid ✅
📋 Machine targeting: Valid ✅
```

## Action Planning

### Plan Mode
The plan mode shows what would be run without actually running:

```bash
# Show running plan
spooky actions run ./my-project --action deploy --plan
```

**Output:**
```
📋 Running Plan for action: deploy
Target Machines: web-server, db-server
Running Mode: Parallel (4 workers)
Dependencies: backup, test
Estimated Duration: 5 minutes

Actions to run:
1. backup (prerequisite)
2. test (prerequisite)
3. deploy (target action)

Would run commands:
- backup.sh on db-server
- test.sh on web-server, db-server
- deploy.sh on web-server, db-server
```

### Dry Run Mode
The dry run mode simulates running without making changes:

```bash
# Simulate running
spooky actions run ./my-project --action deploy --dry-run
```

## Action Targeting

### Machine Targeting
Actions can target specific machines:

```bash
# Target specific machine
spooky actions run ./my-project --action deploy --machine web-server

# Target multiple machines
spooky actions run ./my-project --action deploy --machine web-server --machine db-server
```

### Tag-based Targeting
Actions can target machines by tags:

```bash
# Target by single tag
spooky actions run ./my-project --action deploy --tags production

# Target by multiple tags
spooky actions run ./my-project --action deploy --tags web,production
```

### Complex Filtering
Actions support complex filter expressions:

```bash
# Complex filter
spooky actions run ./my-project --action deploy --filter "environment=production AND role=web"
```

## Action Dependencies

### Dependency Resolution
The system automatically resolves action dependencies:

```hcl
actions {
  action "deploy" {
    description = "Deploy application"
    command = "deploy.sh"
    dependencies = ["backup", "test"]
  }
  
  action "backup" {
    description = "Create backup"
    command = "backup.sh"
  }
  
  action "test" {
    description = "Run tests"
    command = "test.sh"
    dependencies = ["backup"]
  }
}
```

**Running Order:**
1. `backup` (no dependencies)
2. `test` (depends on backup)
3. `deploy` (depends on backup and test)

### Circular Dependency Detection
The system detects and reports circular dependencies:

```bash
# Validate dependencies
spooky actions validate ./my-project
```

**Error Output:**
```
❌ Circular dependency detected:
deploy -> test -> backup -> deploy
```

## Action Conditions

### Conditional Running
Actions can have conditions for running:

```hcl
actions {
  action "deploy" {
    description = "Deploy application"
    command = "deploy.sh"
    
    conditions {
      environment = "production"
      maintenance_window = true
      system_load < 80
      disk_space > 10GB
    }
  }
}
```

### Condition Types
The system supports various condition types:

- **Environment Conditions**: Check environment variables
- **System Conditions**: Check system resources
- **Time Conditions**: Check time-based conditions
- **Custom Conditions**: User-defined conditions

## Action Retry Logic

### Retry Configuration
Actions can have retry logic for failed runs:

```hcl
actions {
  action "deploy" {
    description = "Deploy application"
    command = "deploy.sh"
    
    retry {
      attempts = 3
      delay_seconds = 30
      backoff_multiplier = 2
      max_delay_seconds = 300
    }
  }
}
```

### Retry Behavior
- **Immediate Retry**: Retry immediately on failure
- **Exponential Backoff**: Increase delay between retries
- **Maximum Retries**: Limit total retry attempts
- **Conditional Retry**: Retry only for specific error types

## Action Timeouts

### Timeout Configuration
Actions can have running timeouts:

```hcl
actions {
  action "deploy" {
    description = "Deploy application"
    command = "deploy.sh"
    timeout_seconds = 300  # 5 minutes
  }
  
  action "backup" {
    description = "Create backup"
    command = "backup.sh"
    timeout_seconds = 1800  # 30 minutes
  }
}
```

### Timeout Handling
- **Command Timeout**: Timeout for individual commands
- **Action Timeout**: Timeout for entire action running
- **Graceful Termination**: Clean shutdown on timeout
- **Timeout Reporting**: Report timeout reasons

## Integration with Other Systems

### SSH Integration
Actions use SSH for remote running:

```bash
# Actions use SSH connections
spooky actions run ./my-project --action deploy
```

### Facts Integration
Actions can use facts for conditional running:

```hcl
actions {
  action "deploy" {
    description = "Deploy application"
    command = "deploy.sh"
    
    conditions {
      facts.memory_available > 2GB
      facts.disk_space > 10GB
    }
  }
}
```

### Variables Integration
Actions can use variables for dynamic configuration:

```hcl
actions {
  action "deploy" {
    description = "Deploy application"
    command = "deploy.sh --version {{ variables.app_version }}"
  }
}
```

## Error Handling

### Common Running Errors
- **SSH Connection Errors**: Network connectivity issues
- **Authentication Errors**: SSH key or password issues
- **Command Running Errors**: Commands failing on remote machines
- **Timeout Errors**: Actions exceeding time limits
- **Dependency Errors**: Failed prerequisite actions

### Error Recovery
```bash
# Retry failed actions
spooky actions run ./my-project --action deploy

# Check action status
spooky actions validate ./my-project
```

## Troubleshooting

### Common Issues

#### Running Issues
```bash
# Check action configuration
spooky actions validate ./my-project

# Test SSH connectivity
spooky machines ping ./my-project

# Run with verbose output
spooky actions run ./my-project --action deploy --verbose
```

#### Dependency Issues
```bash
# Validate dependencies
spooky actions validate ./my-project

# Check dependency order
spooky actions run ./my-project --action deploy --plan
```

#### Performance Issues
```bash
# Adjust parallel workers
spooky actions run ./my-project --action deploy --parallel 2

# Check system resources
spooky machines ping ./my-project
```

### Debug Information
```bash
# Enable debug logging
export SPOOKY_LOG_LEVEL=debug

# Run with debug output
spooky actions run ./my-project --action deploy --verbose
```

## Best Practices

### Action Design
1. **Use Descriptive Names**: Choose clear, descriptive action names
2. **Include Descriptions**: Add descriptions for all actions
3. **Set Appropriate Timeouts**: Configure reasonable timeouts
4. **Use Tags**: Apply meaningful tags for organization
5. **Document Dependencies**: Clearly document action dependencies

### Running Practices
1. **Test Actions**: Test actions before production use
2. **Use Dry Run**: Use dry-run mode for testing
3. **Monitor Running**: Monitor action running progress
4. **Handle Errors**: Implement proper error handling
5. **Use Retry Logic**: Configure retry logic for transient failures

### Security Practices
1. **Validate Commands**: Validate all commands before running
2. **Use Least Privilege**: Use minimal required privileges
3. **Audit Actions**: Log and audit action running
4. **Secure Credentials**: Use secure authentication methods
5. **Validate Inputs**: Validate all action inputs

## Future Enhancements

### Planned Features
- **Action Templates**: Reusable action templates
- **Action Scheduling**: Scheduled action running
- **Action Rollback**: Automatic action rollback
- **Action Analytics**: Running analytics and reporting
- **Action Collaboration**: Multi-user action management

### Extension Points
- **Custom Actions**: User-defined action types
- **External Integrations**: Integration with external systems
- **Action Plugins**: Pluggable action running modules
- **Action APIs**: REST API for action management
- **Action Webhooks**: Webhook notifications for action events
