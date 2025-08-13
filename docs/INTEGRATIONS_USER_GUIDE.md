# Integrations User Guide

## Overview

The Integrations system in spooky provides centralized coordination for all system components. It allows you to check the status of various integrations and validate that they're working correctly.

## Available Integrations

spooky includes the following integrations:

- **facts**: Collect and manage system facts
- **actions**: Load and run actions
- **variables**: Manage and resolve variables
- **templates**: Load and render templates
- **machines**: Manage machine inventory and connectivity
- **secrets**: Handle encryption and key management
- **config**: Load and validate configuration

## CLI Commands

### List Integrations

Check which integrations are available and their current status:

```bash
spooky integrations list
```

**Example Output:**
```
Available Integrations:
=======================
facts        ✅ available
actions      ✅ available
variables    ❌ unavailable
templates    ✅ available
machines     ✅ available
secrets      ✅ available
config       ✅ available
```

**Status Indicators:**
- ✅ **available**: Integration is properly initialized and functional
- ❌ **unavailable**: Integration is not configured or not working

### Validate Integrations

Validate that all integrations are working correctly:

```bash
spooky integrations validate
```

**Example Output (Success):**
```
✅ All integrations are working correctly
```

**Example Output (Failure):**
```
❌ Integration validation failed:
  - integration facts is not available
  - integration actions is not available
  - integration variables is not available
  - integration machines is not available
```

## Common Use Cases

### 1. System Health Check

Before running critical operations, validate that all integrations are working:

```bash
# Check integration status
spooky integrations list

# Validate all integrations
spooky integrations validate
```

### 2. Troubleshooting

If you encounter issues with spooky operations, check integration status:

```bash
# List all integrations to see which ones are unavailable
spooky integrations list

# Get detailed validation errors
spooky integrations validate
```

### 3. Development and Testing

During development, verify that all required integrations are properly configured:

```bash
# Quick status check
spooky integrations list

# Full validation
spooky integrations validate
```

## Integration Status

### Available Integrations

When an integration shows as "available", it means:

- The integration is properly initialized
- All required dependencies are configured
- The integration can perform its core functions
- Health checks pass successfully

### Unavailable Integrations

When an integration shows as "unavailable", it means:

- The integration is not properly configured
- Required dependencies are missing
- The integration failed to initialize
- Health checks are failing

## Troubleshooting

### Integration Not Available

If an integration shows as unavailable:

1. **Check configuration**: Ensure all required configuration files are present
2. **Verify dependencies**: Make sure all required services and libraries are available
3. **Check logs**: Look for error messages in spooky logs
4. **Validate manually**: Try using the specific integration directly

### Validation Failures

If `spooky integrations validate` fails:

1. **Review error messages**: Each error provides specific information about what's wrong
2. **Check individual integrations**: Use `spooky integrations list` to see which integrations are unavailable
3. **Verify configuration**: Ensure all required configuration is correct
4. **Check system resources**: Verify that system resources (memory, disk, network) are sufficient

### Common Issues

#### Facts Integration Unavailable
- Check that facts database is accessible
- Verify facts collection configuration
- Ensure storage backend is working

#### Actions Integration Unavailable
- Verify actions configuration files exist
- Check that action templates are valid
- Ensure action dependencies are met

#### Variables Integration Unavailable
- Check variables configuration files
- Verify variable resolution dependencies
- Ensure variable validation passes

#### Templates Integration Unavailable
- Verify template files exist and are accessible
- Check template syntax and validation
- Ensure template functions are available

#### Machines Integration Unavailable
- Check machine inventory configuration
- Verify SSH connectivity settings
- Ensure machine authentication is configured

#### Secrets Integration Unavailable
- Verify encryption keys are available
- Check key permissions and access
- Ensure encryption libraries are working

#### Config Integration Unavailable
- Check configuration file syntax
- Verify configuration validation
- Ensure configuration dependencies are met

## Best Practices

### 1. Regular Health Checks

Perform regular integration health checks:

```bash
# Daily health check
spooky integrations validate

# Quick status check
spooky integrations list
```

### 2. Pre-Operation Validation

Before running critical operations, validate integrations:

```bash
# Validate before running actions
spooky integrations validate && spooky actions run /path/to/project
```

### 3. Monitoring Integration Status

Monitor integration status in production environments:

- Set up automated health checks
- Alert on integration failures
- Log integration status changes

### 4. Configuration Management

- Keep integration configurations in version control
- Use consistent configuration patterns
- Validate configurations before deployment

## Integration Dependencies

Some integrations depend on others:

- **Actions** may depend on **Variables** for variable resolution
- **Templates** may depend on **Variables** for data injection
- **Machines** may depend on **Secrets** for authentication
- **Facts** may depend on **Machines** for data collection

## Error Handling

Integration errors provide detailed information:

- **Error messages** describe the specific problem
- **Field information** indicates which configuration is problematic
- **Validation results** aggregate multiple errors for comprehensive reporting

## Performance Considerations

- Integration health checks are lightweight and fast
- Validation operations may take longer depending on integration complexity
- Status checks use cached information when possible
- Health monitoring has minimal performance impact

## Security Considerations

- Integration status information may reveal system configuration details
- Validation operations may access sensitive configuration data
- Health checks should be performed in secure environments
- Integration failures should be logged appropriately without exposing sensitive information
