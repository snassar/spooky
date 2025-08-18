# Projects Troubleshooting Guide

## Overview

This guide provides solutions for common issues encountered when working with the Projects System in spooky. It covers project initialization, configuration, validation, and operational problems.

## Quick Diagnostic Commands

### Basic Health Check
```bash
# Check project status
spooky project status

# Validate project configuration
spooky project validate

# Show project information
spooky project info

# Check project structure
spooky project structure
```

### Detailed Diagnostics
```bash
# Validate with verbose output
spooky project validate --verbose

# Check specific components
spooky project status --components

# Export configuration for review
spooky project config --export

# Check project logs
spooky project logs
```

## Common Issues and Solutions

### Project Initialization Issues

#### Issue: Project Already Exists
**Error Message:**
```
Error: project "my-project" already exists
```

**Symptoms:**
- Cannot initialize new project
- Directory already contains project files

**Solutions:**

1. **Check if project exists:**
   ```bash
   ls -la my-project/
   ```

2. **Remove existing project (if safe):**
   ```bash
   rm -rf my-project/
   spooky project init my-project
   ```

3. **Initialize in different location:**
   ```bash
   spooky project init my-project --path /different/location
   ```

4. **Use different project name:**
   ```bash
   spooky project init my-project-v2
   ```

#### Issue: Invalid Project Name
**Error Message:**
```
Error: invalid project name "my project"
```

**Symptoms:**
- Project name contains invalid characters
- Project name is too long or too short

**Solutions:**

1. **Use valid project name:**
   ```bash
   spooky project init my-project
   ```

2. **Follow naming conventions:**
   - Use lowercase letters, numbers, and hyphens only
   - Keep names between 3-50 characters
   - Avoid special characters and spaces

3. **Check project name requirements:**
   ```bash
   spooky project init --help
   ```

#### Issue: Permission Denied
**Error Message:**
```
Error: permission denied creating project directory
```

**Symptoms:**
- Cannot create project files
- Insufficient permissions

**Solutions:**

1. **Check directory permissions:**
   ```bash
   ls -ld .
   ```

2. **Fix directory permissions:**
   ```bash
   chmod 755 .
   ```

3. **Use different directory:**
   ```bash
   spooky project init my-project --path ~/projects/
   ```

4. **Run with appropriate permissions:**
   ```bash
   sudo spooky project init my-project
   ```

### Configuration Issues

#### Issue: Invalid HCL Syntax
**Error Message:**
```
Error: invalid HCL syntax in project.hcl
```

**Symptoms:**
- Configuration file cannot be parsed
- Syntax errors in HCL file

**Solutions:**

1. **Validate HCL syntax:**
   ```bash
   spooky project validate --syntax-only
   ```

2. **Check for common syntax errors:**
   - Missing closing braces `}`
   - Incorrect attribute syntax
   - Invalid block structure

3. **Use HCL formatter:**
   ```bash
   hclfmt project.hcl
   ```

4. **Example of correct syntax:**
   ```hcl
   project {
     name = "my-project"
     description = "A sample project"
     
     metadata {
       version = "1.0.0"
       author = "spooky-user"
     }
   }
   ```

#### Issue: Missing Required Fields
**Error Message:**
```
Error: missing required field "name" in project configuration
```

**Symptoms:**
- Configuration validation fails
- Required fields are missing

**Solutions:**

1. **Check required fields:**
   ```bash
   spooky project validate --check required
   ```

2. **Add missing fields:**
   ```hcl
   project {
     name = "my-project"           # Required
     description = "My project"    # Required
     
     metadata {
       version = "1.0.0"          # Required
       author = "spooky-user"     # Required
     }
   }
   ```

3. **Use configuration template:**
   ```bash
   spooky project init my-project --template basic
   ```

#### Issue: Invalid Field Values
**Error Message:**
```
Error: invalid value for "parallel_workers": must be between 1 and 100
```

**Symptoms:**
- Field values are outside valid ranges
- Invalid data types

**Solutions:**

1. **Check field constraints:**
   ```bash
   spooky project validate --check constraints
   ```

2. **Fix invalid values:**
   ```hcl
   run {
     max_parallel = 4        # Must be 1-100
     default_timeout = 300   # Must be positive
     dry_run_default = false
     validate_before_run = true
     backup_before_changes = false
   }
   ```

3. **Use default values:**
   ```hcl
   run {
     # Use defaults for unspecified values
   }
   ```

### Component Issues

#### Issue: Component Not Enabled
**Error Message:**
```
Error: component "facts" is not enabled
```

**Symptoms:**
- Cannot use component functionality
- Component is disabled in configuration

**Solutions:**

1. **Check component status:**
   ```bash
   spooky project status --components
   ```

2. **Enable component in configuration:**
   ```hcl
   components {
     facts = true
     actions = true
     machines = true
   }
   ```

3. **Enable with configuration:**
   ```hcl
   components {
     facts {
       enabled = true
       collection_interval = "30m"
     }
   }
   ```

#### Issue: Component Configuration Error
**Error Message:**
```
Error: invalid configuration for component "facts"
```

**Symptoms:**
- Component configuration is invalid
- Component-specific settings are incorrect

**Solutions:**

1. **Validate component configuration:**
   ```bash
   spooky project validate --component facts
   ```

2. **Check component documentation:**
   ```bash
   spooky facts --help
   ```

3. **Use default component configuration:**
   ```hcl
   components {
     facts = true  # Use defaults
   }
   ```

4. **Fix component-specific settings:**
   ```hcl
   components {
     facts {
       enabled = true
       collection_interval = "30m"  # Valid duration
       storage_backend = "badgerdb" # Valid backend
     }
   }
   ```

### Environment Issues

#### Issue: Environment Not Found
**Error Message:**
```
Error: environment "production" not found
```

**Symptoms:**
- Cannot access environment
- Environment is not configured

**Solutions:**

1. **Check available environments:**
   ```bash
   spooky project config --environments
   ```

2. **Add environment configuration:**
   ```hcl
   environment {
     production {
       enabled = true
       machines = ["prod-*"]
       variables = ["prod/*.hcl"]
     }
   }
   ```

3. **Use existing environment:**
   ```bash
   spooky project status --environment development
   ```

#### Issue: Environment Configuration Error
**Error Message:**
```
Error: invalid environment configuration for "production"
```

**Symptoms:**
- Environment settings are invalid
- Environment-specific configuration errors

**Solutions:**

1. **Validate environment configuration:**
   ```bash
   spooky project validate --environment production
   ```

2. **Check environment settings:**
   ```hcl
   environment {
     production {
       enabled = true
       machines = ["prod-server-1", "prod-server-2"]
       variables = ["prod-vars.hcl"]
       security_level = "high"
     }
   }
   ```

3. **Use environment template:**
   ```bash
   spooky project init my-project --template production
   ```

### Security Issues

#### Issue: Permission Validation Failed
**Error Message:**
```
Error: file permission validation failed for project.hcl
```

**Symptoms:**
- File permissions are too permissive
- Security validation fails

**Solutions:**

1. **Check file permissions:**
   ```bash
   ls -la project.hcl
   ```

2. **Fix file permissions:**
   ```bash
   chmod 644 project.hcl
   chmod 600 secrets.hcl
   ```

3. **Disable permission validation (not recommended):**
   ```hcl
   security {
     validate_file_permissions = false
   }
   ```

#### Issue: Encryption Key Not Found
**Error Message:**
```
Error: encryption key not found at /path/to/keys
```

**Symptoms:**
- Cannot encrypt/decrypt data
- Key file is missing or inaccessible

**Solutions:**

1. **Check key file existence:**
   ```bash
   ls -la /path/to/keys
   ```

2. **Generate new keys:**
   ```bash
   spooky secrets generate-keys --path /path/to/keys
   ```

3. **Update key path in configuration:**
   ```hcl
   security {
     encryption {
       key_path = "/correct/path/to/keys"
     }
   }
   ```

4. **Disable encryption (not recommended):**
   ```hcl
   security {
     encrypt_sensitive_data = false
   }
   ```

### Performance Issues

#### Issue: Project Operations Timeout
**Error Message:**
```
Error: operation timed out after 300 seconds
```

**Symptoms:**
- Operations take too long
- Timeout errors occur

**Solutions:**

1. **Increase timeout settings:**
   ```hcl
   run {
     default_timeout = 600  # Increase timeout
   }
   ```

2. **Optimize parallel workers:**
   ```hcl
   run {
     max_parallel = 8  # Adjust based on system
   }
   ```

3. **Check system resources:**
   ```bash
   top
   df -h
   free -h
   ```

4. **Reduce workload:**
   - Process fewer machines at once
   - Reduce concurrent operations
   - Optimize component settings

#### Issue: High Memory Usage
**Error Message:**
```
Error: memory allocation failed
```

**Symptoms:**
- High memory consumption
- Out of memory errors

**Solutions:**

1. **Check memory usage:**
   ```bash
   free -h
   ps aux | grep spooky
   ```

2. **Reduce memory usage:**
   ```hcl
   run {
     max_parallel = 2  # Reduce workers
   }
   ```

3. **Optimize component settings:**
   ```hcl
   components {
     facts {
       storage_backend = "memory"  # Use memory backend
     }
   }
   ```

4. **Increase system memory or use swap:**
   ```bash
   sudo fallocate -l 2G /swapfile
   sudo chmod 600 /swapfile
   sudo mkswap /swapfile
   sudo swapon /swapfile
   ```

### Integration Issues

#### Issue: Component Integration Failed
**Error Message:**
```
Error: failed to integrate with facts system
```

**Symptoms:**
- Components cannot communicate
- Integration errors occur

**Solutions:**

1. **Check component status:**
   ```bash
   spooky project status --components
   ```

2. **Validate component configuration:**
   ```bash
   spooky project validate --component facts
   spooky project validate --component actions
   ```

3. **Restart components:**
   ```bash
   spooky project restart --component facts
   ```

4. **Check component dependencies:**
   ```hcl
   components {
     facts = true      # Enable facts first
     actions = true    # Then actions
     machines = true   # Then machines
   }
   ```

#### Issue: External Service Unavailable
**Error Message:**
```
Error: external service "database" unavailable
```

**Symptoms:**
- Cannot connect to external services
- Service dependencies are down

**Solutions:**

1. **Check service availability:**
   ```bash
   ping database-server
   telnet database-server 5432
   ```

2. **Configure service retry:**
   ```hcl
   run {
     # Retry settings are handled at the component level
   }
   ```

3. **Use fallback configuration:**
   ```hcl
   components {
     facts {
       storage_backend = "memory"  # Fallback to memory
     }
   }
   ```

4. **Check service configuration:**
   - Verify connection strings
   - Check authentication credentials
   - Validate network connectivity

## Diagnostic Tools

### Project Validation Tool
```bash
# Comprehensive validation
spooky project validate --all

# Validate specific aspects
spooky project validate --check syntax,structure,config

# Generate validation report
spooky project validate --report validation-report.json
```

### Project Status Tool
```bash
# Show overall status
spooky project status

# Show detailed status
spooky project status --verbose

# Show component status
spooky project status --components

# Show environment status
spooky project status --environments
```

### Project Configuration Tool
```bash
# Show current configuration
spooky project config

# Export configuration
spooky project config --export config.json

# Validate configuration
spooky project config --validate

# Show configuration differences
spooky project config --diff
```

### Project Logs Tool
```bash
# Show recent logs
spooky project logs

# Show logs with timestamps
spooky project logs --timestamps

# Filter logs by level
spooky project logs --level error

# Follow logs in real-time
spooky project logs --follow
```

## Recovery Procedures

### Project Recovery
```bash
# Stop project
spooky project stop

# Clean project state
spooky project clean

# Validate project
spooky project validate

# Start project
spooky project start
```

### Configuration Recovery
```bash
# Backup current configuration
cp project.hcl project.hcl.backup

# Restore from backup
cp project.hcl.backup project.hcl

# Validate restored configuration
spooky project validate
```

### Component Recovery
```bash
# Stop specific component
spooky project stop --component facts

# Clean component state
spooky project clean --component facts

# Restart component
spooky project start --component facts
```

## Prevention Best Practices

### Configuration Management
1. **Use version control** for all project configurations
2. **Validate configurations** before deployment
3. **Use configuration templates** for consistency
4. **Document configuration changes** in commit messages
5. **Test configurations** in development environment first

### Monitoring and Alerting
1. **Monitor project status** regularly
2. **Set up alerts** for critical issues
3. **Log all operations** for debugging
4. **Track performance metrics** over time
5. **Review logs** periodically

### Security Practices
1. **Use appropriate file permissions** for configuration files
2. **Encrypt sensitive data** using the secrets system
3. **Rotate encryption keys** regularly
4. **Audit project operations** for compliance
5. **Implement access controls** for project operations

### Performance Optimization
1. **Monitor resource usage** during operations
2. **Optimize parallel worker settings** based on system capacity
3. **Use appropriate timeouts** for operations
4. **Enable caching** for frequently accessed data
5. **Scale components** based on workload requirements

## Getting Help

### Documentation
- [Projects System](PROJECTS_SYSTEM.md) - Complete system overview
- [Projects API Reference](PROJECTS_API_REFERENCE.md) - API documentation
- [Projects User Guide](PROJECTS_USER_GUIDE.md) - User guide and examples

### Debugging Commands
```bash
# Enable debug mode
spooky project validate --debug

# Show detailed error information
spooky project status --verbose

# Export diagnostic information
spooky project diagnose --output diagnostic.json
```

### Log Analysis
```bash
# Search for errors in logs
grep -i error .spooky/logs/*.log

# Search for specific component logs
grep -i facts .spooky/logs/*.log

# Show recent errors
tail -f .spooky/logs/*.log | grep -i error
```

## Related Documentation

- [Projects System](PROJECTS_SYSTEM.md) - Complete system overview
- [Projects API Reference](PROJECTS_API_REFERENCE.md) - API documentation
- [Projects User Guide](PROJECTS_USER_GUIDE.md) - User guide and examples
- [Configuration Management](CONFIGURATION_SYSTEM.md) - Configuration integration
- [Schema System](SCHEMA_SYSTEM.md) - Schema integration
- [CLI Reference](CLI_REFERENCE.md) - CLI integration
