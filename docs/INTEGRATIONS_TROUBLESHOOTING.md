# Integrations Troubleshooting Guide

## Overview

This guide helps you diagnose and resolve issues with spooky integrations. The integrations system provides centralized coordination for all system components, and problems can manifest in various ways.

## Quick Diagnostic Commands

Start with these commands to assess the current state:

```bash
# Check integration status
spooky integrations list

# Validate all integrations
spooky integrations validate

# Check spooky logs for integration errors
spooky --log-level debug integrations validate
```

## Common Error Patterns

### All Integrations Unavailable

**Symptoms:**
```
Available Integrations:
=======================
facts        ❌ unavailable
actions      ❌ unavailable
variables    ❌ unavailable
templates    ❌ unavailable
machines     ❌ unavailable
secrets      ❌ unavailable
config       ❌ unavailable
```

**Possible Causes:**
1. **IntegrationManager not initialized**: The factory failed to create integrations
2. **Logger not configured**: IntegrationManager requires a logger
3. **Factory configuration error**: Factory is not properly configured
4. **System resource issues**: Memory, disk, or network problems

**Diagnostic Steps:**
1. Check if IntegrationManager is being created correctly
2. Verify logger configuration
3. Check system resources
4. Review factory configuration

**Solutions:**
```go
// Ensure proper factory setup
factory := integration.NewFactory(logger)
manager := factory.CreateIntegrationManager()

// Verify manager creation
if manager == nil {
    // Handle factory failure
}
```

### Specific Integration Unavailable

**Symptoms:**
```
Available Integrations:
=======================
facts        ✅ available
actions      ❌ unavailable  // Only this one is unavailable
variables    ✅ available
templates    ✅ available
machines     ✅ available
secrets      ✅ available
config       ✅ available
```

**Diagnostic Steps:**
1. Check integration-specific configuration
2. Verify dependencies for the specific integration
3. Review integration initialization logs
4. Test integration functionality directly

## Integration-Specific Issues

### Facts Integration Issues

**Common Problems:**
- Facts database not accessible
- Storage backend configuration errors
- Facts collection failures
- Database corruption

**Diagnostic Commands:**
```bash
# Check facts database accessibility
spooky facts list /path/to/project

# Validate facts configuration
spooky facts validate /path/to/project

# Check storage backend
spooky facts gather /path/to/project --dry-run
```

**Solutions:**
1. **Database not accessible**: Check file permissions and disk space
2. **Configuration errors**: Validate facts configuration files
3. **Collection failures**: Check network connectivity and SSH settings
4. **Database corruption**: Recreate facts database

### Actions Integration Issues

**Status: ✅ Fully Functional** - The actions integration is complete with full SSH-based execution, dependency resolution, and comprehensive error handling.

**Common Problems:**
- Actions configuration file not found
- Action validation failures
- SSH connectivity issues during execution
- Action dependencies not met
- Machine targeting problems

**Diagnostic Commands:**
```bash
# Check actions configuration
spooky actions list /path/to/project

# Validate actions
spooky actions validate /path/to/project

# Test action execution (dry-run)
spooky actions run /path/to/project --dry-run

# Check action plan
spooky actions run /path/to/project --plan
```

**Solutions:**
1. **Missing configuration**: Ensure `actions.hcl` exists and is valid
2. **Validation failures**: Fix action syntax and dependencies
3. **SSH issues**: Check SSH connectivity and authentication
4. **Dependency issues**: Resolve action dependencies and circular references
5. **Machine targeting**: Verify machine names and tags are correct

### Variables Integration Issues

**Common Problems:**
- Variables file not found
- Variable resolution failures
- Circular dependencies
- Variable validation errors

**Diagnostic Commands:**
```bash
# Check variables configuration
spooky variables list /path/to/project

# Validate variables
spooky variables validate /path/to/project

# Test variable resolution
spooky variables resolve /path/to/project
```

**Solutions:**
1. **Missing files**: Ensure `variables.hcl` or `variables/` directory exists
2. **Resolution failures**: Check variable dependencies and references
3. **Circular dependencies**: Break dependency cycles
4. **Validation errors**: Fix variable syntax and constraints

### Templates Integration Issues

**Common Problems:**
- Template files not found
- Template syntax errors
- Variable injection failures
- Template function errors

**Diagnostic Commands:**
```bash
# Check templates
spooky templates list /path/to/project

# Validate templates
spooky templates validate /path/to/project

# Test template rendering
spooky templates render /path/to/project template.tmpl --dry-run
```

**Solutions:**
1. **Missing templates**: Ensure template files exist in `templates/` directory
2. **Syntax errors**: Fix template syntax and variable references
3. **Variable issues**: Check variable resolution and injection
4. **Function errors**: Verify template functions are available

### Machines Integration Issues

**Common Problems:**
- Machine inventory not found
- SSH connectivity failures
- Authentication errors
- Machine validation failures

**Diagnostic Commands:**
```bash
# Check machine inventory
spooky machines list /path/to/project

# Validate machines
spooky machines validate /path/to/project

# Test connectivity
spooky machines ping /path/to/project
```

**Solutions:**
1. **Missing inventory**: Ensure `machines.hcl` exists and is valid
2. **SSH issues**: Check SSH configuration and connectivity
3. **Authentication**: Verify SSH keys and authentication methods
4. **Validation errors**: Fix machine configuration syntax

### Secrets Integration Issues

**Common Problems:**
- Encryption keys not available
- Key validation failures
- Encryption/decryption errors
- Key permissions issues

**Diagnostic Commands:**
```bash
# Check secrets configuration
spooky secrets validate

# Test key validation
spooky secrets test-key

# Check encryption
spooky secrets encrypt --test
```

**Solutions:**
1. **Missing keys**: Import encryption keys (spooky does not generate keys)
2. **Key validation**: Fix key format and permissions
3. **Encryption errors**: Check encryption libraries and algorithms
4. **Permission issues**: Fix key file permissions (600)

### Config Integration Issues

**Common Problems:**
- Configuration file not found
- Configuration validation failures
- Configuration syntax errors
- Configuration dependencies not met

**Diagnostic Commands:**
```bash
# Check configuration
spooky config show

# Validate configuration
spooky config validate

# Check configuration path
spooky config show --path
```

**Solutions:**
1. **Missing config**: Ensure configuration files exist
2. **Validation errors**: Fix configuration syntax and validation rules
3. **Syntax errors**: Check HCL/JSON syntax
4. **Dependencies**: Install required dependencies

## Health Monitoring Issues

### Health Status Not Updating

**Symptoms:**
- Integration status shows as unavailable even after fixing issues
- Health status doesn't reflect current state
- Status changes not visible in CLI output

**Diagnostic Steps:**
1. Check if health status is being updated correctly
2. Verify thread safety of health status updates
3. Check for race conditions in status updates
4. Review health monitoring implementation

**Solutions:**
```go
// Ensure proper health status updates
manager.UpdateHealthStatus("facts", true)

// Verify health status
status := manager.GetHealthStatus()
if status["facts"] {
    // Facts integration is healthy
}
```

### Validation Failures

**Symptoms:**
```
❌ Integration validation failed:
  - integration facts is not available
  - integration actions is not available
```

**Diagnostic Steps:**
1. Check individual integration health
2. Review validation error messages
3. Verify integration initialization
4. Check for configuration errors

**Solutions:**
1. **Fix integration initialization**: Ensure integrations are properly created
2. **Update health status**: Mark integrations as healthy after fixing issues
3. **Check configuration**: Verify all required configuration is correct
4. **Review dependencies**: Ensure all dependencies are available

## Performance Issues

### Slow Integration Validation

**Symptoms:**
- `spooky integrations validate` takes too long
- Health checks are slow
- Integration operations are sluggish

**Diagnostic Steps:**
1. Check system resources (CPU, memory, disk)
2. Review integration initialization time
3. Check for blocking operations
4. Verify network connectivity

**Solutions:**
1. **Resource optimization**: Increase system resources
2. **Async operations**: Use asynchronous initialization where possible
3. **Caching**: Implement health status caching
4. **Connection pooling**: Use connection pooling for network operations

### Memory Issues

**Symptoms:**
- High memory usage during integration operations
- Memory leaks in integration manager
- Out of memory errors

**Diagnostic Steps:**
1. Monitor memory usage during operations
2. Check for memory leaks in integrations
3. Review resource cleanup
4. Verify garbage collection

**Solutions:**
1. **Resource cleanup**: Ensure proper cleanup of resources
2. **Memory optimization**: Optimize memory usage in integrations
3. **Connection management**: Properly close connections and pools
4. **Garbage collection**: Force garbage collection if needed

## Logging and Debugging

### Enable Debug Logging

```bash
# Enable debug logging for integrations
spooky --log-level debug integrations validate

# Check integration manager logs
spooky --log-level debug integrations list
```

### Common Log Messages

**Integration Initialization:**
```
INFO IntegrationManager initialized component=test facts_available=false actions_available=false
```

**Health Status Updates:**
```
INFO Integration health status updated integration=facts healthy=true
```

**Validation Results:**
```
INFO System health validation completed component=test valid=false errors=3 warnings=0
```

**Coordinated Operations:**
```
INFO Starting coordinated operation
INFO Coordinated operation completed successfully
```

### Debug Integration Issues

```go
// Enable debug logging in code
logger.SetLevel("debug")

// Add debug logging to integrations
logger.Debug("Integration operation started", map[string]interface{}{
    "integration": "facts",
    "operation":   "CollectFacts",
})
```

## Recovery Procedures

### Integration Recovery

1. **Stop affected operations**
2. **Diagnose the issue** using diagnostic commands
3. **Fix the root cause** (configuration, dependencies, etc.)
4. **Update health status** to reflect the fix
5. **Validate the fix** using integration commands
6. **Resume operations**

### System Recovery

1. **Check all integrations** using `spooky integrations list`
2. **Validate system health** using `spooky integrations validate`
3. **Fix critical issues** first (facts, actions, machines)
4. **Fix secondary issues** (variables, templates, secrets, config)
5. **Verify full recovery** by running validation again

### Configuration Recovery

1. **Backup current configuration**
2. **Restore from known good configuration**
3. **Validate configuration** using appropriate commands
4. **Test functionality** with simple operations
5. **Gradually restore customizations**

## Prevention

### Best Practices

1. **Regular health checks**: Run `spooky integrations validate` regularly
2. **Configuration validation**: Validate configurations before deployment
3. **Monitoring**: Set up monitoring for integration health
4. **Backup**: Keep backups of critical configurations
5. **Testing**: Test integrations in staging environments

### Monitoring

```bash
# Set up automated health checks
crontab -e
# Add: 0 */6 * * * spooky integrations validate

# Monitor integration status
watch -n 30 'spooky integrations list'

# Log integration status
spooky integrations validate >> /var/log/spooky/integrations.log 2>&1
```

### Alerting

Set up alerts for:
- Integration failures
- Health check failures
- Configuration validation errors
- Performance degradation

## Getting Help

If you're still experiencing issues:

1. **Check logs**: Review all relevant log files
2. **Gather diagnostics**: Run diagnostic commands and save output
3. **Document steps**: Document the exact steps to reproduce the issue
4. **Check documentation**: Review relevant documentation
5. **Seek support**: Contact support with detailed information
