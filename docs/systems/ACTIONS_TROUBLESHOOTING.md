# Actions System Troubleshooting Guide

## Overview

This troubleshooting guide provides solutions for common issues encountered when working with the spooky actions system. It covers error messages, configuration problems, performance issues, and debugging techniques.

**Status: Partially Implemented** - The actions system has basic functionality but SSH-based action orchestration has known issues that need to be addressed.

## Actions System Status

### ⚠️ Partially Functional Actions Infrastructure

The actions system currently has **basic actions infrastructure** with:

- **Action Loading**: Loading actions from HCL configuration files
- **Action Validation**: Basic action validation and error handling
- **CLI Integration**: `spooky actions` commands with basic functionality
- **Project Integration**: Actions loading from project configuration
- **Local Action Orchestration**: Local action orchestration capabilities
- **Basic Planning**: Action planning with dependency resolution

### Known Limitations

- **SSH-Based Orchestration**: SSH-based action orchestration has implementation issues
- **Remote Action Running**: Cannot properly run actions on remote machines
- **Parallel Execution**: No parallel action execution support
- **Action Planning**: Limited action planning capabilities
- **Result Aggregation**: No action result aggregation

### Expected Behavior

When using actions, you can expect:

1. **Action Loading**: Actions can be loaded from HCL configuration files
2. **Action Validation**: Basic action validation with error reporting
3. **CLI Integration**: Actions commands work with project configuration
4. **Local Orchestration**: Local action orchestration capabilities
5. **Basic Planning**: Action planning with dependency resolution

## Common Issues and Solutions

### Action Loading Errors

#### "Failed to load actions: file not found"

**Problem:** The actions file or directory does not exist.

**Solution:**
```bash
# Check if actions file exists
ls -la ./actions.hcl

# Check if actions directory exists
ls -la ./actions/

# Create basic actions file
cat > actions.hcl << 'EOF'
actions {
  action "test-action" {
    command = "echo 'Hello World'"
    machines = ["test-machine"]
  }
}
EOF
```

**Check your project structure:**
```bash
my-project/
├── project.hcl
├── machines.hcl
├── actions.hcl          # Main actions file
└── actions/             # Optional actions directory
    ├── web-actions.hcl
    └── db-actions.hcl
```

#### "Failed to parse actions: HCL syntax error"

**Problem:** The actions file contains invalid HCL syntax.

**Solution:**
```bash
# Validate HCL syntax
spooky actions validate ./my-project

# Check for common syntax errors
cat -n actions.hcl | grep -E "(missing|unexpected|invalid)"

# Fix common syntax issues
# - Missing closing braces
# - Invalid quotes
# - Missing commas
# - Invalid field names
```

**Example of correct syntax:**
```hcl
actions {
  action "restart-nginx" {
    command = "sudo systemctl restart nginx"
    description = "Restart nginx service"
    machines = ["web-server"]
  }
}
```

#### "Failed to validate actions: schema validation failed"

**Problem:** Action configuration does not match the required schema.

**Solution:**
```bash
# Validate against schema
spooky actions validate ./my-project --verbose

# Check required fields
# - name: Action name is required
# - command: Command to run is required
# - machines or tags: Targeting is required
```

**Check your action configuration:**
```hcl
actions {
  action "valid-action" {
    name = "valid-action"           # Required
    command = "echo 'test'"         # Required
    machines = ["web-server"]       # Required (machines or tags)
  }
}
```

### Action Validation Errors

#### "Validation failed: action type 'invalid_type' is not supported"

**Problem:** The specified action type is not supported.

**Solution:**
```bash
# Check supported action types
spooky actions validate ./my-project --verbose

# Valid action types:
# - command (basic shell commands)
# - script (script files)
# - template_deploy (template deployment)
# - file_copy (file copying)
# - service_control (service management)
```

**Fix the action type:**
```hcl
actions {
  action "fix-action" {
    command = "echo 'test'"  # Use valid action type
    machines = ["web-server"]
  }
}
```

#### "Validation failed: circular dependency detected"

**Problem:** Actions have circular dependencies.

**Solution:**
```bash
# Check dependency graph
spooky actions run ./my-project --plan

# Identify circular dependencies
spooky actions validate ./my-project --verbose
```

**Fix circular dependencies:**
```hcl
# ❌ WRONG - Circular dependency
actions {
  action "action1" {
    command = "echo 'action1'"
    dependencies = ["action2"]
    machines = ["web-server"]
  }

  action "action2" {
    command = "echo 'action2'"
    dependencies = ["action1"]  # Creates circular dependency
    machines = ["web-server"]
  }
}

# ✅ CORRECT - No circular dependency
actions {
  action "action1" {
    command = "echo 'action1'"
    machines = ["web-server"]
  }

  action "action2" {
    command = "echo 'action2'"
    dependencies = ["action1"]  # action1 must complete first
    machines = ["web-server"]
  }
}
```

#### "Validation failed: target machine 'invalid-machine' not found"

**Problem:** Action targets a machine that doesn't exist in the inventory.

**Solution:**
```bash
# Check available machines
spooky machines list ./my-project

# Check machine inventory
cat machines.hcl

# Fix machine targeting
spooky actions validate ./my-project --verbose
```

**Check your machine inventory:**
```hcl
machines {
  machine "web-server" {
    hostname = "web.example.com"
    port = 22
    user = "admin"
    
    authentication {
      method = "ssh_key"
      key_path = "~/.ssh/id_rsa"
    }
  }
}
```

### Action Run Errors

#### "Failed to run action: SSH connection failed"

**Problem:** Cannot establish SSH connection to target machine.

**Solution:**
```bash
# Test SSH connectivity
spooky machines ping ./my-project

# Check SSH configuration
ssh -i ~/.ssh/id_rsa user@machine.example.com

# Verify machine configuration
cat machines.hcl
```

**Check your machine configuration:**
```hcl
machines {
  machine "web-server" {
    hostname = "web.example.com"  # Verify hostname is correct
    port = 22                     # Verify SSH port
    user = "admin"                # Verify username
    
    authentication {
      method = "ssh_key"
      key_path = "~/.ssh/id_rsa"  # Verify key path
    }
  }
}
```

#### "Failed to run action: command run failed"

**Problem:** Command run failed on target machine.

**Solution:**
```bash
# Check command syntax
spooky actions run ./my-project --action test-action --verbose

# Test command manually
ssh user@machine.example.com "your-command"

# Check for permission issues
spooky actions run ./my-project --action test-action --sudo
```

**Check your command configuration:**
```hcl
actions {
  action "test-command" {
    command = "echo 'Hello World'"  # Use simple commands first
    
    machines = ["web-server"]
    sudo = false  # Only use sudo if necessary
    timeout = 300
  }
}
```

#### "Failed to run action: service control failed"

**Problem:** Service control operation failed.

**Solution:**
```bash
# Check service status manually
ssh user@machine.example.com "systemctl status nginx"

# Check service permissions
ssh user@machine.example.com "sudo systemctl status nginx"

# Test service control manually
ssh user@machine.example.com "sudo systemctl restart nginx"
```

**Check your service configuration:**
```hcl
actions {
  action "restart-nginx" {
    command = "sudo systemctl restart nginx"
    description = "Restart nginx service"
    
    machines = ["web-server"]
    sudo = true  # Service control usually requires sudo
    timeout = 30
  }
}
```

### Dependency Errors

#### "Failed to resolve dependencies: action 'missing-action' not found"

**Problem:** Action depends on a non-existent action.

**Solution:**
```bash
# List all available actions
spooky actions list ./my-project

# Check action names
spooky actions validate ./my-project --verbose

# Fix dependency reference
```

**Check your dependencies:**
```hcl
# ✅ CORRECT - Valid dependency
actions {
  action "deploy-app" {
    command = "echo 'deploying app'"
    dependencies = ["prepare-database"]  # This action must exist
    machines = ["app-server"]
  }

  action "prepare-database" {
    command = "echo 'preparing database'"
    machines = ["db-server"]
  }
}
```

#### "Failed to resolve dependencies: circular dependency detected"

**Problem:** Actions have circular dependencies.

**Solution:**
```bash
# Check dependency graph
spooky actions run ./my-project --plan

# Identify the circular dependency
spooky actions validate ./my-project --verbose
```

**Fix circular dependencies by restructuring:**
```hcl
# ❌ WRONG - Circular dependency
actions {
  action "action1" {
    command = "echo 'action1'"
    dependencies = ["action2"]
    machines = ["web-server"]
  }

  action "action2" {
    command = "echo 'action2'"
    dependencies = ["action1"]  # Creates circular dependency
    machines = ["web-server"]
  }
}

# ✅ CORRECT - Sequential running
actions {
  action "action1" {
    command = "echo 'action1'"
    machines = ["web-server"]
  }

  action "action2" {
    command = "echo 'action2'"
    dependencies = ["action1"]  # action1 must complete first
    machines = ["web-server"]
  }

  action "action3" {
    command = "echo 'action3'"
    dependencies = ["action2"]  # action2 must complete first
    machines = ["web-server"]
  }
}
```

## Configuration Problems

### Action Configuration Issues

#### "Invalid action configuration: missing required field"

**Problem:** Action configuration is missing required fields.

**Solution:**
```hcl
# ✅ CORRECT - All required fields present
actions {
  action "valid-action" {
    name = "valid-action"           # Required
    command = "echo 'test'"         # Required
    machines = ["web-server"]       # Required (machines or tags)
  }
}
```

#### "Invalid action configuration: unsupported action type"

**Problem:** Action type is not supported.

**Solution:**
```hcl
# Supported action types:
actions {
  action "command-action" {
    command = "echo 'test'"         # Run shell commands
    machines = ["web-server"]
  }

  action "script-action" {
    command = "bash /path/to/script.sh"  # Run script files
    machines = ["web-server"]
  }

  action "service-action" {
    command = "sudo systemctl restart nginx"  # Control services
    machines = ["web-server"]
  }
}
```

### Machine Targeting Issues

#### "No target machines found for action"

**Problem:** Action has no valid target machines.

**Solution:**
```bash
# Check available machines
spooky machines list ./my-project

# Check machine inventory
cat machines.hcl

# Fix machine targeting
```

**Check your machine targeting:**
```hcl
# ✅ CORRECT - Target by machine names
actions {
  action "web-action" {
    command = "echo 'web action'"
    machines = ["web-server", "web-server-2"]  # Must exist in machines.hcl
  }
}

# ✅ CORRECT - Target by tags
actions {
  action "prod-action" {
    command = "echo 'production action'"
    machines = ["web-server"]  # Machines must have these tags
  }
}

# ✅ CORRECT - Target all machines
actions {
  action "all-action" {
    command = "echo 'all machines'"
    machines = ["web-server"]  # Target specific machines
  }
}
```

### Template Configuration Issues

#### "Template file not found: templates/test.conf.tmpl"

**Problem:** Template file does not exist.

**Solution:**
```bash
# Check if template file exists
ls -la templates/test.conf.tmpl

# Create template file
mkdir -p templates
cat > templates/test.conf.tmpl << 'EOF'
# Test configuration
server_name = "{{.server_name}}"
port = {{.port}}
EOF
```

**Check your template configuration:**
```hcl
actions {
  action "deploy-config" {
    command = "cp /tmp/test.conf /etc/test.conf"
    description = "Deploy configuration file"
    
    machines = ["web-server"]
  }
}
```

### Script Configuration Issues

#### "Script file not found: files/test.sh"

**Problem:** Script file does not exist.

**Solution:**
```bash
# Check if script file exists
ls -la files/test.sh

# Create script file
mkdir -p files
cat > files/test.sh << 'EOF'
#!/bin/bash
echo "Test script run"
EOF
chmod +x files/test.sh
```

**Check your script configuration:**
```hcl
actions {
  action "run-script" {
    command = "bash files/test.sh"  # File must exist in files/ or templates/
    
    machines = ["web-server"]
    working_directory = "/tmp"
  }
}
```

## Performance Issues

### Slow Action Running

**Problem:** Action running is taking too long.

**Solution:**
```bash
# Use parallel running (if supported)
spooky actions run ./my-project --parallel 4

# Increase timeout
spooky actions run ./my-project --timeout 600

# Profile running
spooky actions run ./my-project --verbose
```

**Optimize your action configuration:**
```hcl
actions {
  action "optimized-action" {
    command = "echo 'optimized'"
    
    machines = ["web-server", "web-server-2", "web-server-3"]
    timeout = 300                # Set appropriate timeout
  }
}
```

### High Resource Usage

**Problem:** Action running is consuming too many resources.

**Solution:**
```bash
# Reduce parallel workers
spooky actions run ./my-project --parallel 2

# Monitor resource usage
top -p $(pgrep spooky)

# Use resource limits
```

**Set resource limits:**
```hcl
actions {
  action "resource-limited" {
    command = "resource-intensive-command"
    
    machines = ["web-server"]
    timeout = 300  # Set timeout to limit resource usage
  }
}
```

### Network Connectivity Issues

**Problem:** Network issues affecting action running.

**Solution:**
```bash
# Test network connectivity
ping machine.example.com

# Test SSH connectivity
ssh -o ConnectTimeout=10 user@machine.example.com

# Check DNS resolution
nslookup machine.example.com

# Increase SSH timeout
```

**Optimize network configuration:**
```hcl
machines {
  machine "web-server" {
    hostname = "web.example.com"
    port = 22
    user = "admin"
    
    authentication {
      method = "ssh_key"
      key_path = "~/.ssh/id_rsa"
    }
  }
}
```

## Debugging Techniques

### Enable Verbose Output

Enable verbose output to get detailed information:

```bash
# Enable verbose output for all commands
spooky actions list ./my-project --verbose
spooky actions validate ./my-project --verbose
spooky actions run ./my-project --verbose

# Enable debug logging
export SPOOKY_LOG_LEVEL=debug
spooky actions run ./my-project
```

### Use Planning Mode

Use planning mode to understand what will be run:

```bash
# Show run plan
spooky actions run ./my-project --plan

# Show plan with details
spooky actions run ./my-project --plan --verbose

# Check dependencies
spooky actions run ./my-project --plan | grep -A 5 -B 5 "dependency"
```

### Use Dry Run Mode

Use dry run mode to test without actual running:

```bash
# Simulate running
spooky actions run ./my-project --dry-run

# Simulate with verbose output
spooky actions run ./my-project --dry-run --verbose

# Test specific action
spooky actions run ./my-project --action test-action --dry-run
```

### Check Action Configuration

Validate and check action configuration:

```bash
# Validate all actions
spooky actions validate ./my-project

# Validate specific action
spooky actions validate ./my-project --action test-action

# Check action details
spooky actions list ./my-project --verbose
```

### Network Diagnostics

Diagnose network connectivity issues:

```bash
# Test SSH connectivity
spooky machines ping ./my-project

# Test specific machine
ssh -v user@machine.example.com

# Check SSH configuration
ssh -o BatchMode=yes -o ConnectTimeout=10 user@machine.example.com echo "SSH working"
```

### File System Diagnostics

Check file system issues:

```bash
# Check if files exist
ls -la actions.hcl
ls -la actions/
ls -la templates/
ls -la files/

# Check file permissions
ls -la ~/.ssh/id_rsa

# Check disk space
df -h
```

## Recovery Procedures

### Action Run Recovery

Recover from failed action running:

```bash
# Check run status
spooky actions run ./my-project --verbose

# Retry failed actions
spooky actions run ./my-project --action failed-action

# Skip failed actions
spooky actions run ./my-project --skip-failed

# Clean up and restart
pkill spooky
spooky actions run ./my-project
```

### Configuration Recovery

Recover from configuration issues:

```bash
# Backup configuration
cp actions.hcl actions.hcl.backup

# Validate configuration
spooky actions validate ./my-project

# Restore from backup if needed
cp actions.hcl.backup actions.hcl

# Test configuration
spooky actions validate ./my-project
```

### SSH Connection Recovery

Recover from SSH connection issues:

```bash
# Test SSH connectivity
spooky machines ping ./my-project

# Check SSH keys
ls -la ~/.ssh/
chmod 600 ~/.ssh/id_rsa

# Note: spooky does not generate SSH keys - use openssh to regenerate keys if needed
ssh-keygen -t rsa -b 4096 -f ~/.ssh/id_rsa

# Copy keys to machines
ssh-copy-id -i ~/.ssh/id_rsa user@machine.example.com
```

## Prevention Strategies

### Regular Validation

Regularly validate action configurations:

```bash
# Schedule regular validation
crontab -e
# Add: 0 2 * * * /usr/local/bin/spooky actions validate /path/to/project

# Validate before deployment
spooky actions validate ./my-project

# Validate in CI/CD pipeline
spooky actions validate ./my-project --strict
```

### Testing Strategy

Implement comprehensive testing:

```bash
# Test with dry run
spooky actions run ./my-project --dry-run

# Test with plan mode
spooky actions run ./my-project --plan

# Test specific actions
spooky actions run ./my-project --action test-action --dry-run

# Test on staging environment first
spooky actions run ./staging-project
```

### Monitoring

Monitor action running:

```bash
# Monitor run logs
tail -f /var/log/spooky/actions.log

# Monitor system resources
top -p $(pgrep spooky)

# Monitor SSH connections
netstat -an | grep :22
```

### Backup Strategy

Maintain configuration backups:

```bash
# Backup action configurations
cp actions.hcl actions.hcl.$(date +%Y%m%d)

# Version control configurations
git add actions.hcl
git commit -m "Update action configuration"

# Backup project structure
tar -czf project-backup-$(date +%Y%m%d).tar.gz ./
```

## Known Issues and Workarounds

> **See also**: [Known Issues](KNOWN_ISSUES.md#actions-system-ssh-issues) - Comprehensive documentation of all known issues and workarounds

### SSH-Based Orchestration Issues

**Issue:** SSH-based action orchestration has implementation problems.

**Workaround:**
```bash
# Use local action running for immediate needs
spooky actions run ./my-project --local-only

# Run actions manually on remote machines
ssh user@machine.example.com "your-command"
```

### Parallel Execution Issues

**Issue:** Parallel action execution may fail with many machines.

**Workaround:**
```bash
# Use sequential execution
spooky actions run ./my-project --parallel 1

# Run actions in smaller batches
spooky actions run ./my-project --machine machine1,machine2
spooky actions run ./my-project --machine machine3,machine4
```

### Action Planning Issues

**Issue:** Action planning may not work correctly with complex dependencies.

**Workaround:**
```bash
# Use simple action dependencies
# Avoid complex dependency chains

# Test planning mode
spooky actions run ./my-project --plan --verbose
```

## Best Practices for Troubleshooting

### 1. Start Simple

Begin with simple actions and add complexity gradually:

```hcl
# Start with simple command
actions {
  action "test-simple" {
    command = "echo 'Hello World'"
    machines = ["web-server"]
  }
}

# Then add complexity
actions {
  action "test-complex" {
    command = "sudo systemctl restart nginx"
    machines = ["web-server"]
    dependencies = ["test-simple"]
  }
}
```

### 2. Use Descriptive Names

Use clear, descriptive action names:

```hcl
# ✅ GOOD - Descriptive names
actions {
  action "restart-nginx-service" {
    command = "sudo systemctl restart nginx"
    machines = ["web-server"]
  }
}

# ❌ BAD - Unclear names
actions {
  action "action1" {
    command = "systemctl restart nginx"
    machines = ["web-server"]
  }
}
```

### 3. Validate Early and Often

Validate configurations frequently:

```bash
# Validate after every change
spooky actions validate ./my-project

# Validate before running
spooky actions validate ./my-project && spooky actions run ./my-project

# Validate in scripts
#!/bin/bash
if spooky actions validate ./my-project; then
    spooky actions run ./my-project
else
    echo "Validation failed"
    exit 1
fi
```

### 4. Use Proper Error Handling

Implement proper error handling in actions:

```hcl
actions {
  action "robust-action" {
    command = "your-command || exit 1"
    
    machines = ["web-server"]
    timeout = 300
    allow_failure = false
  }
}
```

### 5. Monitor and Log

Monitor action running and maintain logs:

```bash
# Enable verbose logging
spooky actions run ./my-project --verbose

# Monitor running
watch -n 1 'ps aux | grep spooky'

# Check logs
tail -f /var/log/spooky/actions.log
```

## Getting Help

### Documentation Resources

1. **User Guide** - For usage questions and best practices
2. **API Reference** - For technical implementation details
3. **Examples** - For configuration patterns and use cases

### Common Questions

#### "Why aren't my actions running?"

1. Check action validation
2. Verify machine connectivity
3. Check SSH configuration
4. Validate action syntax

#### "How do I debug action failures?"

```bash
# Enable verbose output
spooky actions run ./my-project --verbose

# Use dry run mode
spooky actions run ./my-project --dry-run

# Check specific action
spooky actions run ./my-project --action problematic-action --verbose
```

#### "How do I fix dependency issues?"

```bash
# Check dependency graph
spooky actions run ./my-project --plan

# Validate dependencies
spooky actions validate ./my-project --verbose

# Fix circular dependencies
# Remove or restructure circular dependencies
```

#### "How do I optimize action performance?"

```bash
# Use parallel running (if supported)
spooky actions run ./my-project --parallel 4

# Set resource limits
# Use appropriate timeouts
# Monitor resource usage
```

### When to Seek Additional Help

- Configuration validation passes but actions still fail
- Performance issues persist after optimization
- Unusual error messages not covered in this guide
- Integration issues with other spooky components

For additional help, refer to the [User Guide](ACTIONS_USER_GUIDE.md) and [API Reference](ACTIONS_API_REFERENCE.md), or check the project documentation for more advanced troubleshooting techniques.

## Conclusion

The actions system provides basic action management and orchestration capabilities with some limitations. Most issues can be resolved by following the troubleshooting steps outlined in this guide. For persistent issues, enable verbose output and collect diagnostic information for further analysis. SSH-based action orchestration improvements are planned for future releases.
