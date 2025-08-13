# Actions System Troubleshooting Guide

## Overview

This troubleshooting guide provides solutions for common issues encountered when working with the spooky actions system. It covers error messages, configuration problems, performance issues, and debugging techniques.

**Status: Production Ready** - The actions system is fully implemented with complete acting infrastructure, SSH integration, and comprehensive error handling.

## Acting Infrastructure Status

### ✅ Fully Functional Acting Infrastructure

The actions system now has **complete acting infrastructure** with:

- **SSH-Based Execution**: All action types execute via SSH with proper connection management
- **Machine Targeting**: Full support for machine names and tags with proper filtering
- **Dependency Resolution**: Complete dependency resolution and run order planning
- **Parallel Execution**: Parallel action running with dependency resolution
- **Session Management**: Full session lifecycle management and progress tracking
- **Error Handling**: Comprehensive error handling and result aggregation
- **Resource Management**: Proper resource cleanup and timeout handling

### What This Means for Users

- **No More Stubs**: All functionality is fully implemented - no placeholder code
- **Production Ready**: The system is ready for production use
- **Complete Feature Set**: All documented features are functional
- **Reliable Execution**: Robust error handling and recovery mechanisms
- **Performance Optimized**: Efficient execution with proper resource management

### Expected Behavior

When running actions, you can expect:

1. **Proper Machine Targeting**: Actions will run on the correct machines based on your configuration
2. **Dependency Resolution**: Actions will run in the correct order based on dependencies
3. **SSH Execution**: All actions execute via SSH with proper authentication
4. **Result Reporting**: Comprehensive results with success/failure status
5. **Error Recovery**: Proper error handling with detailed error messages
6. **Resource Cleanup**: Automatic cleanup of SSH connections and resources

## Common Error Messages

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
action "test-action" {
  type = "command"
  command = "echo 'Hello World'"
  machines = ["test-machine"]
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
action "restart-nginx" {
  type = "service_control"
  description = "Restart nginx service"
  
  service_control {
    service = "nginx"
    action = "restart"
  }
  
  machines = ["web-server"]
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
# - type: Action type is required
# - machines or tags: Targeting is required
```

**Check your action configuration:**
```hcl
action "valid-action" {
  name = "valid-action"           # Required
  type = "command"                # Required
  command = "echo 'test'"         # Required for command type
  
  machines = ["web-server"]       # Required (machines or tags)
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
# - command
# - script
# - template_deploy
# - file_copy
# - service_control
```

**Fix the action type:**
```hcl
action "fix-action" {
  type = "command"  # Use valid action type
  command = "echo 'test'"
  machines = ["web-server"]
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
action "action1" {
  type = "command"
  command = "echo 'action1'"
  dependencies = ["action2"]
  machines = ["web-server"]
}

action "action2" {
  type = "command"
  command = "echo 'action2'"
  dependencies = ["action1"]  # Creates circular dependency
  machines = ["web-server"]
}

# ✅ CORRECT - No circular dependency
action "action1" {
  type = "command"
  command = "echo 'action1'"
  machines = ["web-server"]
}

action "action2" {
  type = "command"
  command = "echo 'action2'"
  dependencies = ["action1"]  # action1 must complete first
  machines = ["web-server"]
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
action "test-command" {
  type = "command"
  command = "echo 'Hello World'"  # Use simple commands first
  
  machines = ["web-server"]
  sudo = false  # Only use sudo if necessary
  timeout = 300
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
action "restart-nginx" {
  type = "service_control"
  
  service_control {
    service = "nginx"           # Verify service name
    action = "restart"          # Valid actions: start, stop, restart, reload, enable, disable, status
    systemd = true              # Use systemd (default)
    timeout = 30                # Service operation timeout
  }
  
  machines = ["web-server"]
  sudo = true  # Service control usually requires sudo
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
action "deploy-app" {
  type = "command"
  command = "echo 'deploying app'"
  dependencies = ["prepare-database"]  # This action must exist
  machines = ["app-server"]
}

action "prepare-database" {
  type = "command"
  command = "echo 'preparing database'"
  machines = ["db-server"]
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
action "action1" {
  type = "command"
  command = "echo 'action1'"
  dependencies = ["action2"]
  machines = ["web-server"]
}

action "action2" {
  type = "command"
  command = "echo 'action2'"
  dependencies = ["action1"]  # Creates circular dependency
  machines = ["web-server"]
}

# ✅ CORRECT - Sequential running
action "action1" {
  type = "command"
  command = "echo 'action1'"
  machines = ["web-server"]
}

action "action2" {
  type = "command"
  command = "echo 'action2'"
  dependencies = ["action1"]  # action1 must complete first
  machines = ["web-server"]
}

action "action3" {
  type = "command"
  command = "echo 'action3'"
  dependencies = ["action2"]  # action2 must complete first
  machines = ["web-server"]
}
```

## Configuration Problems

### Action Configuration Issues

#### "Invalid action configuration: missing required field"

**Problem:** Action configuration is missing required fields.

**Solution:**
```hcl
# ✅ CORRECT - All required fields present
action "valid-action" {
  name = "valid-action"           # Required
  type = "command"                # Required
  command = "echo 'test'"         # Required for command type
  machines = ["web-server"]       # Required (machines or tags)
}
```

#### "Invalid action configuration: unsupported action type"

**Problem:** Action type is not supported.

**Solution:**
```hcl
# Supported action types:
action "command-action" {
  type = "command"                # Run shell commands
  command = "echo 'test'"
  machines = ["web-server"]
}

action "script-action" {
  type = "script"                 # Run script files
  script = "files/test.sh"
  machines = ["web-server"]
}

action "template-action" {
  type = "template_deploy"        # Deploy template files
  template {
    source = "templates/test.conf.tmpl"
    destination = "/etc/test.conf"
  }
  machines = ["web-server"]
}

action "file-action" {
  type = "file_copy"              # Copy files
  file_copy {
    source = "files/test.txt"
    destination = "/tmp/test.txt"
  }
  machines = ["web-server"]
}

action "service-action" {
  type = "service_control"        # Control services
  service_control {
    service = "nginx"
    action = "restart"
  }
  machines = ["web-server"]
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
action "web-action" {
  type = "command"
  command = "echo 'web action'"
  machines = ["web-server", "web-server-2"]  # Must exist in machines.hcl
}

# ✅ CORRECT - Target by tags
action "prod-action" {
  type = "command"
  command = "echo 'production action'"
  tags = ["environment=production", "role=web"]  # Machines must have these tags
}

# ✅ CORRECT - Target all machines
action "all-action" {
  type = "command"
  command = "echo 'all machines'"
  # No machines or tags specified = target all machines
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
action "deploy-config" {
  type = "template_deploy"
  
  template {
    source = "templates/test.conf.tmpl"      # File must exist
    destination = "/etc/test.conf"
    permissions = "0644"
    owner = "root"
    group = "root"
  }
  
  machines = ["web-server"]
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
action "run-script" {
  type = "script"
  script = "files/test.sh"  # File must exist in files/ or templates/
  
  machines = ["web-server"]
  working_directory = "/tmp"
}
```

## Performance Issues

### Slow Action Running

**Problem:** Action running is taking too long.

**Solution:**
```bash
# Use parallel running
spooky actions run ./my-project --parallel 4

# Increase timeout
spooky actions run ./my-project --timeout 600

# Profile running
spooky actions run ./my-project --verbose
```

**Optimize your action configuration:**
```hcl
action "optimized-action" {
  type = "command"
  command = "echo 'optimized'"
  
  machines = ["web-server", "web-server-2", "web-server-3"]
  parallel = true              # Run in parallel
  timeout = 300                # Set appropriate timeout
  max_concurrent = 4           # Limit concurrent runs
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
action "resource-limited" {
  type = "command"
  command = "resource-intensive-command"
  
  machines = ["web-server"]
  
  resource_limits {
    memory_mb = 1024           # Limit memory usage
    cpu_percent = 50           # Limit CPU usage
    disk_mb = 512              # Limit disk usage
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
    
    # Add connection options
    connection {
      timeout = 30
      retries = 3
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

# Regenerate SSH keys if needed
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

## Best Practices for Troubleshooting

### 1. Start Simple

Begin with simple actions and add complexity gradually:

```hcl
# Start with simple command
action "test-simple" {
  type = "command"
  command = "echo 'Hello World'"
  machines = ["web-server"]
}

# Then add complexity
action "test-complex" {
  type = "service_control"
  service_control {
    service = "nginx"
    action = "restart"
  }
  machines = ["web-server"]
  dependencies = ["test-simple"]
}
```

### 2. Use Descriptive Names

Use clear, descriptive action names:

```hcl
# ✅ GOOD - Descriptive names
action "restart-nginx-service" {
  type = "service_control"
  service_control {
    service = "nginx"
    action = "restart"
  }
  machines = ["web-server"]
}

# ❌ BAD - Unclear names
action "action1" {
  type = "command"
  command = "systemctl restart nginx"
  machines = ["web-server"]
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
action "robust-action" {
  type = "command"
  command = "your-command || exit 1"
  
  machines = ["web-server"]
  retries = 3
  retry_delay = 10
  timeout = 300
  allow_failure = false
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
# Use parallel running
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

This comprehensive troubleshooting guide provides solutions for the most common issues encountered with the spooky actions system, along with prevention strategies and recovery procedures.
