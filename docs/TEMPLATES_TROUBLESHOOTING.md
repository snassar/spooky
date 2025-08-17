# Templates System Troubleshooting Guide

## Overview

This troubleshooting guide provides comprehensive solutions for common issues encountered when using the spooky templates system. It covers error diagnosis, debugging techniques, and step-by-step resolution procedures.

**Status: Partially Implemented** - The templates system has basic functionality but CLI commands and SSH-based template rendering have known issues that need to be addressed.

> **See also**: [Known Issues](KNOWN_ISSUES.md#templates-system-ssh-issues) - Comprehensive documentation of all known issues and workarounds

## Quick Diagnosis

### Common Error Patterns

| Error Pattern | Likely Cause | Quick Fix |
|---------------|--------------|-----------|
| `template file not found` | Missing template file | Check file path and existence |
| `invalid template syntax` | Template syntax error | Validate template syntax |
| `required variable not provided` | Missing data | Check data file and variables |
| `security violation` | Dangerous pattern detected | Review template content |
| `failed to load template` | File access issue | Check permissions and path |
| `template validation failed` | Template validation error | Review validation rules |

### System Status Check

Before troubleshooting, verify the system status:

```bash
# Check spooky installation
spooky --version

# Check project structure
ls -la myproject/

# Check template directory
ls -la myproject/templates/

# Check variables files
ls -la myproject/variables/
```

## Template Loading Issues

### Template File Not Found

#### Error Message
```bash
Error: template file not found: templates/nginx.conf.tmpl
```

#### Diagnosis Steps

1. **Check File Existence**:
   ```bash
   ls -la myproject/templates/nginx.conf.tmpl
   ```

2. **Check File Path**:
   ```bash
   # Verify working directory
   pwd
   
   # Check relative path
   find . -name "nginx.conf.tmpl"
   ```

3. **Check File Permissions**:
   ```bash
   ls -la myproject/templates/
   ```

#### Solutions

**File Doesn't Exist**:
```bash
# Create the template file
mkdir -p myproject/templates
touch myproject/templates/nginx.conf.tmpl
```

**Wrong Path**:
```bash
# Use correct path
spooky templates render ./myproject templates/nginx.conf.tmpl
```

**Permission Issues**:
```bash
# Fix permissions
chmod 644 myproject/templates/nginx.conf.tmpl
```

### Template Loading Failed

#### Error Message
```bash
Error: failed to load template templates/nginx.conf.tmpl: invalid template syntax
```

#### Diagnosis Steps

1. **Check Template Syntax**:
   ```bash
   # Validate template syntax
   spooky templates validate ./myproject --template templates/nginx.conf.tmpl
   ```

2. **Check for Common Syntax Errors**:
   ```bash
   # Look for unclosed tags
   grep -n "{{" myproject/templates/nginx.conf.tmpl
   grep -n "}}" myproject/templates/nginx.conf.tmpl
   ```

3. **Check for Invalid Characters**:
   ```bash
   # Check for encoding issues
   file myproject/templates/nginx.conf.tmpl
   ```

#### Solutions

**Unclosed Template Tags**:
```bash
# Fix unclosed tags
# Before: {{.variable
# After:  {{.variable}}
```

**Invalid Template Syntax**:
```bash
# Fix syntax errors
# Before: {{.variable | function arg1 arg2}}
# After:  {{.variable | function "arg1" "arg2"}}
```

**Encoding Issues**:
```bash
# Convert to UTF-8
iconv -f ISO-8859-1 -t UTF-8 template.tmpl > template_utf8.tmpl
mv template_utf8.tmpl template.tmpl
```

## Template Rendering Issues

### Missing Required Variables

#### Error Message
```bash
Error: template validation failed: required variable "server_name" not provided
```

#### Diagnosis Steps

1. **Check Variables File**:
   ```bash
   cat myproject/variables.hcl
   ```

2. **Check Variable Names**:
   ```bash
   # Look for variable definitions
   grep -n "server_name" myproject/variables.hcl
   ```

3. **Check Template Variables**:
   ```bash
   # Look for variable usage in template
   grep -n "server_name" myproject/templates/nginx.conf.tmpl
   ```

#### Solutions

**Missing Variable in Variables File**:
```hcl
# Add missing variable to variables.hcl
server_name = "example.com"
port = 80
root_path = "/var/www/html"
```

**Variable Name Mismatch**:
```bash
# Fix variable name in template
# Before: {{.server_name}}
# After:  {{.serverName}}  # or fix data file to match
```

**Variables File Not Loaded**:
```bash
# Variables are loaded automatically from variables.hcl and variables/ directory
spooky templates render ./myproject templates/nginx.conf.tmpl
```

### Template Rendering Failed

#### Error Message
```bash
Error: failed to render template: template execution failed
```

#### Diagnosis Steps

1. **Check Template Content**:
   ```bash
   cat myproject/templates/nginx.conf.tmpl
   ```

2. **Check Variables Structure**:
   ```bash
   # Validate variables file
   spooky config validate variables.hcl
   ```

3. **Check for Function Errors**:
   ```bash
   # Look for function usage
   grep -n "|" myproject/templates/nginx.conf.tmpl
   ```

#### Solutions

**Invalid Function Usage**:
```bash
# Fix function syntax
# Before: {{.items | len | add 1}}
# After:  {{add (len .items) 1}}
```

**Data Type Mismatch**:
```bash
# Fix data types
# Before: {{.port | add "80"}}
# After:  {{add .port 80}}
```

**Missing Function**:
```bash
# Use available functions
# Before: {{.text | customFunction}}
# After:  {{.text | upper}}
```

## Security Issues

### Security Violation Detected

#### Error Message
```bash
Error: security violation: dangerous pattern detected: {{.system}}
```

#### Diagnosis Steps

1. **Check Template Content**:
   ```bash
   # Look for dangerous patterns
   grep -n "system\|exec\|eval" myproject/templates/nginx.conf.tmpl
   ```

2. **Check Security Level**:
   ```bash
   # Check template metadata
   cat myproject/templates/nginx.conf.tmpl.meta
   ```

3. **Check Forbidden Patterns**:
   ```bash
   # Common forbidden patterns
   # - exec, system, eval, shell
   # - password, secret, key, token
   # - {{.*os\\.Run.*}}
   # - {{.*system.*}}
   # - {{.*eval.*}}
   ```

#### Solutions

**Remove Dangerous Patterns**:
```bash
# Remove or replace dangerous patterns
# Before: {{.user_input | exec}}
# After:  {{.user_input | upper}}
```

**Adjust Security Level**:
```hcl
# Set appropriate security level in template metadata
security_level = "standard"  # or "elevated" for trusted templates
```

**Use Safe Alternatives**:
```bash
# Use safe functions instead
# Before: {{.command | exec}}
# After:  {{.command | upper}}
```

### Function Access Denied

#### Error Message
```bash
Error: function access denied: "system" not allowed in restricted mode
```

#### Diagnosis Steps

1. **Check Function Restrictions**:
   ```bash
   # Look for restricted functions
   grep -n "system\|exec\|eval" myproject/templates/nginx.conf.tmpl
   ```

2. **Check Security Configuration**:
   ```bash
   # Check template configuration
   cat myproject/templates/config.hcl
   ```

#### Solutions

**Use Allowed Functions**:
```bash
# Replace restricted functions with allowed ones
# Before: {{.text | system}}
# After:  {{.text | upper}}
```

**Adjust Security Level**:
```hcl
# Increase security level for trusted templates
security_level = "elevated"
```

**Register Custom Functions**:
```go
// Register safe custom functions
customFunctions := map[string]interface{}{
    "safeFunction": func(s string) string { return "safe: " + s },
}
```

## Performance Issues

### Template Rendering Slow

#### Symptoms
- Template rendering takes more than 5 seconds
- High memory usage during rendering
- Slow response times

#### Diagnosis Steps

1. **Check Template Size**:
   ```bash
   wc -l myproject/templates/nginx.conf.tmpl
   du -h myproject/templates/nginx.conf.tmpl
   ```

2. **Check Template Complexity**:
   ```bash
   # Count template tags
   grep -c "{{" myproject/templates/nginx.conf.tmpl
   grep -c "}}" myproject/templates/nginx.conf.tmpl
   ```

3. **Check Variables Size**:
   ```bash
   du -h myproject/variables.hcl
   ```

#### Solutions

**Optimize Template Structure**:
```bash
# Break large templates into smaller ones
# Before: One large template with 1000+ lines
# After:  Multiple smaller templates
```

**Use Efficient Functions**:
```bash
# Use efficient functions
# Before: {{range .items}}{{.value | complex_function}}{{end}}
# After:  {{.items | efficient_function}}
```

**Enable Caching**:
```hcl
# Enable template caching
cache_ttl = 300              # Cache for 5 minutes
max_cache_size = 1000        # Maximum cache entries
```

### Memory Usage Issues

#### Symptoms
- High memory usage during template rendering
- Out of memory errors
- Slow performance

#### Diagnosis Steps

1. **Check Memory Usage**:
   ```bash
   # Monitor memory usage
   ps aux | grep spooky
   ```

2. **Check Template Variables**:
   ```bash
   # Check variables file size
   du -h myproject/variables/
   ```

#### Solutions

**Optimize Data Structure**:
```hcl
# Use efficient data structures
# Before: Large nested objects
# After:  Flattened data structures
```

**Limit Template Size**:
```bash
# Break large templates into smaller ones
# Before: One large template
# After:  Multiple smaller templates
```

**Use Streaming**:
```bash
# Use streaming for large templates
# Process templates in chunks
```

## Integration Issues

### Facts Integration Problems

#### Error Message
```bash
Error: failed to resolve facts context: facts integration not available
```

#### Diagnosis Steps

1. **Check Facts Integration**:
   ```bash
   # Check if facts system is working
   spooky facts gather ./myproject
   ```

2. **Check Facts Data**:
   ```bash
   # Check facts database
   ls -la ~/.local/state/spooky/
   ```

#### Solutions

**Enable Facts Integration**:
```go
// Set facts integration in template manager
manager.SetFactsIntegration(factsIntegration)
```

**Check Facts Data**:
```bash
# Verify facts are available
spooky facts list ./myproject
```

### Variables Integration Problems

#### Error Message
```bash
Error: failed to resolve variables context: variables integration not available
```

#### Diagnosis Steps

1. **Check Variables Integration**:
   ```bash
   # Check variables system
   spooky variables list ./myproject
   ```

2. **Check Variables Files**:
   ```bash
   # Check variables files
   ls -la myproject/variables/
   cat myproject/variables.hcl
   ```

#### Solutions

**Enable Variables Integration**:
```go
// Set variables integration in template manager
manager.SetVariablesIntegration(variablesIntegration)
```

**Check Variables Files**:
```bash
# Validate variables files
spooky variables validate ./myproject
```

### Machines Integration Problems

#### Error Message
```bash
Error: failed to resolve machines context: machines integration not available
```

#### Diagnosis Steps

1. **Check Machines Integration**:
   ```bash
   # Check machines system
   spooky machines list ./myproject
   ```

2. **Check Machine Inventory**:
   ```bash
   # Check machine inventory
   cat myproject/machines.hcl
   ```

#### Solutions

**Enable Machines Integration**:
```go
// Set machines integration in template manager
manager.SetMachinesIntegration(machinesIntegration)
```

**Check Machine Inventory**:
```bash
# Validate machine inventory
spooky machines validate ./myproject
```

## CLI Issues

### Command Not Found

#### Error Message
```bash
Error: unknown command "templates"
```

#### Diagnosis Steps

1. **Check Spooky Installation**:
   ```bash
   spooky --version
   spooky --help
   ```

2. **Check Command Availability**:
   ```bash
   # List available commands
   spooky --help | grep templates
   ```

#### Solutions

**Update Spooky Installation**:
```bash
# Update to latest version
go install github.com/your-org/spooky@latest
```

**Check Command Implementation**:
```bash
# Templates CLI commands may not be implemented yet
# Use API directly or wait for CLI implementation
```

### Command Options Issues

#### Error Message
```bash
Error: unknown flag: --data
```

#### Diagnosis Steps

1. **Check Command Help**:
   ```bash
   spooky templates render --help
   ```

2. **Check Available Options**:
   ```bash
   # List available options
   spooky templates render --help | grep -A 10 "Flags"
   ```

#### Solutions

**Use Correct Options**:
```bash
# Use correct option names
# Before: --data file.hcl
# After:  --variables file.hcl  # or check actual option name
```

**Check Option Implementation**:
```bash
# Some options may not be implemented yet
# Use basic functionality or wait for implementation
```

## Debugging Techniques

### Verbose Output

#### Enable Verbose Mode
```bash
# Enable verbose output for debugging
spooky templates render ./myproject templates/nginx.conf.tmpl \
  --verbose
```

#### Analyze Verbose Output
```bash
# Look for specific error messages
spooky templates render ./myproject templates/nginx.conf.tmpl \
  --verbose 2>&1 | grep -i error
```

### Template Validation

#### Validate Template Syntax
```bash
# Validate template syntax
spooky templates validate ./myproject --template templates/nginx.conf.tmpl
```

#### Validate Template Data
```bash
# Validate variables file
spooky config validate variables.hcl
```

### Preview Mode

#### Preview Template Rendering
```bash
# Preview template without writing files
spooky templates render ./myproject templates/nginx.conf.tmpl \
  --preview
```

#### Dry Run Mode
```bash
# Show what would be rendered without making changes
spooky templates render ./myproject templates/nginx.conf.tmpl \
  --dry-run
```

## Common Solutions

### Template Syntax Fixes

#### Fix Unclosed Tags
```bash
# Before: {{.variable
# After:  {{.variable}}
```

#### Fix Function Syntax
```bash
# Before: {{.items | len | add 1}}
# After:  {{add (len .items) 1}}
```

#### Fix Conditional Syntax
```bash
# Before: {{if .condition}}
# After:  {{if .condition}}
#         content
#         {{end}}
```

### Data Structure Fixes

#### Fix Variable Names
```bash
# Ensure variable names match between template and data
# Template: {{.server_name}}
# Data:     server_name = "example.com"
```

#### Fix Data Types
```bash
# Use correct data types
# Template: {{add .port 80}}
# Data:     port = 80  # integer, not string
```

#### Fix Nested Data
```bash
# Use correct nested data structure
# Template: {{.database.url}}
# Data:     database = { url = "localhost" }
```

### Security Fixes

#### Remove Dangerous Patterns
```bash
# Replace dangerous patterns with safe alternatives
# Before: {{.user_input | exec}}
# After:  {{.user_input | upper}}
```

#### Adjust Security Level
```hcl
# Set appropriate security level
security_level = "standard"  # or "elevated" for trusted templates
```

#### Use Safe Functions
```bash
# Use only allowed functions
# Allowed: upper, lower, trim, len, add, sub, etc.
# Restricted: exec, system, eval, etc.
```

## Getting Help

### Documentation Resources

1. **API Reference**: [Templates API Reference](TEMPLATES_API_REFERENCE.md)
2. **CLI Reference**: [Templates CLI Reference](TEMPLATES_CLI_REFERENCE.md)
3. **System Overview**: [Templates System](TEMPLATES_SYSTEM.md)
4. **User Guide**: [Templates User Guide](TEMPLATES_USER_GUIDE.md)

### Support Channels

1. **Project Issues**: Report issues on the project repository
2. **Documentation**: Check the comprehensive documentation
3. **Community**: Ask questions in the project community
4. **Examples**: Review example templates and configurations

### Reporting Issues

When reporting issues, include:

1. **Error Message**: Complete error message
2. **Template Content**: Relevant template content
3. **Data Content**: Relevant data file content
4. **System Information**: OS, spooky version, etc.
5. **Steps to Reproduce**: Clear steps to reproduce the issue
6. **Expected Behavior**: What you expected to happen
7. **Actual Behavior**: What actually happened

## Related Documentation

- [Templates API Reference](TEMPLATES_API_REFERENCE.md) - Complete API documentation
- [Templates CLI Reference](TEMPLATES_CLI_REFERENCE.md) - CLI command reference
- [Templates System](TEMPLATES_SYSTEM.md) - System overview and architecture
- [Templates User Guide](TEMPLATES_USER_GUIDE.md) - User guide and examples
- [Templates Documentation Summary](TEMPLATES_DOCUMENTATION_SUMMARY.md) - Documentation overview
