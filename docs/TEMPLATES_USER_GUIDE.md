# Templates System User Guide

## Overview

This user guide provides a comprehensive introduction to using the spooky templates system. It covers getting started, basic and advanced usage patterns, and real-world examples for template management, rendering, and validation.

**Status: Partially Implemented** - The templates system has basic functionality but CLI commands and SSH-based template rendering have known issues that need to be addressed.

## Related Documentation

- [Variables User Guide](VARIABLES_USER_GUIDE.md) - Variable management and resolution for templates
- [Actions User Guide](ACTIONS_USER_GUIDE.md) - Template deployment through actions
- [SSH User Guide](SSH_USER_GUIDE.md) - SSH-based template rendering
- [Machines User Guide](MACHINES_USER_GUIDE.md) - Machine targeting for template deployment
- [Facts User Guide](FACTS_USER_GUIDE.md) - Using machine facts in templates

> **See also**: [User Guides Index](USER_GUIDES_INDEX.md) - Complete overview of all user guides

## Getting Started

### Prerequisites

Before using the templates system, ensure you have:

1. **Spooky Installation**: A working spooky installation
2. **Project Structure**: A valid spooky project with proper structure
3. **Template Files**: Template files in the `templates/` directory
4. **Data Files**: Data files for template variables (optional)

### Basic Project Structure

```
myproject/
├── project.hcl              # Project configuration
├── machines.hcl             # Machine inventory
├── variables.hcl            # Project variables
├── variables/               # Variables directory (optional)
│   ├── production.hcl      # Environment-specific variables
│   └── development.hcl     # Environment-specific variables
├── templates/               # Template directory
│   ├── nginx.conf.tmpl      # Nginx configuration template
│   ├── deploy.sh.tmpl       # Deployment script template
│   └── config.yaml.tmpl     # Configuration template
```

### Creating Your First Template

1. **Create Template Directory**:
   ```bash
   mkdir -p myproject/templates
   ```

2. **Create a Simple Template**:
   ```bash
   # myproject/templates/hello.tmpl
   Hello {{.name}}!
   
   Welcome to {{.environment}} environment.
   
   Server: {{.server_name}}
   Port: {{.port}}
   ```

3. **Create Variables File**:
   ```hcl
   # myproject/variables.hcl
   name = "World"
   environment = "development"
   server_name = "localhost"
   port = 8080
   ```

4. **Render Template** (CLI implementation in progress):
   ```bash
   spooky templates render ./myproject templates/hello.tmpl \
     --output output/hello.txt
   ```
   
   > **Note**: Template rendering CLI commands are partially implemented. For production use, consider using templates through the [Actions User Guide](ACTIONS_USER_GUIDE.md) system for deployment.

> **See also**: [Known Issues](KNOWN_ISSUES.md#templates-system-cli-issues) - Detailed information about template CLI command issues and workarounds

## Template Basics

### Template Syntax

Templates use Go template syntax with the following features:

### Integration with Other Systems

Templates integrate with other spooky systems:

- **Variables**: Resolve [variables](VARIABLES_USER_GUIDE.md) in template rendering
- **Facts**: Use [machine facts](FACTS_USER_GUIDE.md) in templates
- **Actions**: Deploy templates through [action orchestration](ACTIONS_USER_GUIDE.md)
- **SSH**: Render templates on remote machines via [SSH](SSH_USER_GUIDE.md)
- **Secrets**: Use [encrypted variables](SECRETS_USER_GUIDE.md) in templates

#### Variables
```bash
{{.variable_name}}           # Simple variable
{{.nested.field}}           # Nested field access
{{.array.0}}                # Array element access
{{.map.key}}                # Map key access
```

#### Functions
```bash
{{.name | upper}}           # String to uppercase
{{.items | len}}            # Array length
{{.value | default "none"}} # Default value
{{.list | join ","}}        # Join array with separator
```

#### Conditionals
```bash
{{if .condition}}
  Content when condition is true
{{else}}
  Content when condition is false
{{end}}
```

#### Loops
```bash
{{range .items}}
  {{.name}}: {{.value}}
{{end}}
```

### Template Types

#### Configuration Templates
```bash
# templates/nginx.conf.tmpl
server {
    listen {{.port}};
    server_name {{.server_name}};
    root {{.root_path}};
    
    location / {
        try_files $uri $uri/ /index.html;
    }
    
    # SSL configuration
    {{if .ssl_enabled}}
    listen {{.port}} ssl;
    ssl_certificate {{.ssl_cert_path}};
    ssl_certificate_key {{.ssl_key_path}};
    {{end}}
    
    # Logging
    access_log /var/log/nginx/{{.server_name}}.access.log;
    error_log /var/log/nginx/{{.server_name}}.error.log;
}
```

#### Script Templates
```bash
# templates/deploy.sh.tmpl
#!/bin/bash

# Deployment script for {{.app_name}}
echo "Deploying {{.app_name}} version {{.version}}"

# Environment variables
export APP_NAME="{{.app_name}}"
export VERSION="{{.version}}"
export ENVIRONMENT="{{.environment}}"

# {{if .database_enabled}}
echo "Setting up database connection"
export DATABASE_URL="{{.database_url}}"
{{end}}

# Deploy application
echo "Starting deployment..."
{{range .deploy_steps}}
{{.}}
{{end}}

echo "Deployment completed successfully"
```

#### Documentation Templates
```bash
# templates/system-docs.md.tmpl
# System Documentation for {{.system_name}}

## Overview
{{.description}}

## System Information
- **Hostname**: {{.facts.hostname}}
- **OS**: {{.facts.os.name}} {{.facts.os.version}}
- **Architecture**: {{.facts.architecture}}
- **CPU Cores**: {{.facts.cpu.cores}}
- **Memory**: {{.facts.memory.total | div 1024 | div 1024}} GB

## Network Configuration
{{range .facts.network.interfaces}}
### Interface {{.name}}
- **IP Address**: {{.ip}}
- **MAC Address**: {{.mac}}
- **Status**: {{.status}}
{{end}}

## Services
{{range .services}}
### {{.name}}
- **Status**: {{.status}}
- **Port**: {{.port}}
- **Description**: {{.description}}
{{end}}
```

## Template Functions

### String Functions

#### Text Manipulation
```bash
{{.text | upper}}           # Convert to uppercase
{{.text | lower}}           # Convert to lowercase
{{.text | title}}           # Convert to title case
{{.text | trim}}            # Trim whitespace
{{.text | trimLeft}}        # Trim left whitespace
{{.text | trimRight}}       # Trim right whitespace
```

#### Text Processing
```bash
{{.text | replace "old" "new" 1}}  # Replace first occurrence
{{.text | replaceAll "old" "new"}} # Replace all occurrences
{{.text | split ","}}              # Split by separator
{{.items | join ","}}              # Join with separator
{{.text | contains "substring"}}   # Check if contains substring
{{.text | hasPrefix "prefix"}}     # Check if has prefix
{{.text | hasSuffix "suffix"}}     # Check if has suffix
```

#### Text Analysis
```bash
{{.text | len}}             # Get string length
{{.text | repeat 3}}        # Repeat string 3 times
{{.text | substr 0 5}}      # Extract substring
```

### Mathematical Functions

#### Basic Operations
```bash
{{add 1 2}}                 # Addition: 3
{{sub 5 2}}                 # Subtraction: 3
{{mul 3 4}}                 # Multiplication: 12
{{div 10 2}}                # Division: 5
{{mod 7 3}}                 # Modulo: 1
```

#### Advanced Math
```bash
{{abs -5}}                  # Absolute value: 5
{{ceil 3.7}}                # Ceiling: 4
{{floor 3.7}}               # Floor: 3
{{round 3.7}}               # Round: 4
{{min 1 2 3}}               # Minimum: 1
{{max 1 2 3}}               # Maximum: 3
{{pow 2 3}}                 # Power: 8
{{sqrt 16}}                 # Square root: 4
```

### Array Functions

#### Array Access
```bash
{{first .items}}            # First element
{{last .items}}             # Last element
{{index .items 0}}          # Element at index 0
{{slice .items 1 3}}        # Elements from index 1 to 3
```

#### Array Manipulation
```bash
{{append .items "new"}}     # Add element to array
{{prepend .items "new"}}    # Add element to beginning
{{reverse .items}}          # Reverse array
{{sort .items}}             # Sort array
{{uniq .items}}             # Remove duplicates
{{containsItem .items "value"}} # Check if array contains value
```

### Hash and Encoding Functions

#### Cryptographic Functions
```bash
{{md5 "hello"}}             # MD5 hash
{{sha256 "hello"}}          # SHA256 hash
```

#### Encoding Functions
```bash
{{base64 "hello"}}          # Base64 encode
{{base64Decode "aGVsbG8="}} # Base64 decode
{{hex "hello"}}             # Hex encode
{{hexDecode "68656c6c6f"}}  # Hex decode
```

### Type Conversion Functions

```bash
{{toString 123}}            # Convert to string: "123"
{{toInt "123"}}             # Convert to integer: 123
{{toFloat "123.45"}}        # Convert to float: 123.45
{{toBool "true"}}           # Convert to boolean: true
```

### JSON Functions

```bash
{{toJSON .data}}            # Convert to JSON
{{fromJSON '{"key":"value"}'}} # Parse JSON
{{prettyJSON .data}}        # Pretty print JSON
```

### Date and Time Functions

```bash
{{now}}                     # Current time
{{formatTime .timestamp "2006-01-02"}} # Format time
{{parseTime "2024-01-01" "2006-01-02"}} # Parse time
{{addDays .date 7}}         # Add 7 days
{{addHours .time 24}}       # Add 24 hours
```

### Utility Functions

```bash
{{default .value "default"}} # Default value if empty
{{coalesce .a .b .c}}       # First non-empty value
{{ternary .condition "yes" "no"}} # Conditional value
{{regexMatch "^[a-z]+$" "hello"}} # Regex match
{{regexReplace "l" "L" "hello"}} # Regex replace
{{random 1 100}}            # Random number 1-100
{{uuid}}                    # Generate UUID
```

## Template Context

### Available Context Data

Templates have access to rich context data from various sources:

#### Project Information
```bash
{{.project.name}}           # Project name
{{.project.description}}    # Project description
{{.project.version}}        # Project version
{{.project.author}}         # Project author
{{.project.tags}}           # Project tags
```

#### Machine Facts
```bash
{{.facts.hostname}}         # Machine hostname
{{.facts.os.name}}          # Operating system name
{{.facts.os.version}}       # Operating system version
{{.facts.architecture}}     # Machine architecture
{{.facts.cpu.cores}}        # CPU cores
{{.facts.memory.total}}     # Total memory
{{.facts.disk.usage_percent}} # Disk usage percentage
```

#### Machine Inventory
```bash
{{range .machines}}
{{.hostname}} - {{.ip}}
{{end}}
```

#### Environment Variables
```bash
{{.environment.HOME}}       # Home directory
{{.environment.USER}}       # Current user
{{.environment.PATH}}       # PATH variable
```

#### Project Variables
```bash
{{.variables.app_name}}     # Application name
{{.variables.environment}}  # Environment name
{{.variables.version}}      # Application version
```

## Advanced Usage

### Template Composition

#### Base Template
```bash
# templates/base.tmpl
<!DOCTYPE html>
<html>
<head>
    <title>{{.title}}</title>
    <meta charset="utf-8">
    {{if .css}}
    <link rel="stylesheet" href="{{.css}}">
    {{end}}
</head>
<body>
    <header>
        <h1>{{.title}}</h1>
    </header>
    
    <main>
        {{template "content" .}}
    </main>
    
    <footer>
        <p>Generated by Spooky Templates</p>
    </footer>
</body>
</html>
```

#### Content Template
```bash
# templates/content.tmpl
{{define "content"}}
<div class="content">
    <h2>{{.heading}}</h2>
    <p>{{.description}}</p>
    
    {{if .items}}
    <ul>
    {{range .items}}
        <li>{{.name}}: {{.value}}</li>
    {{end}}
    </ul>
    {{end}}
</div>
{{end}}
```

### Conditional Rendering

#### Environment-Specific Configuration
```bash
# templates/config.tmpl
# Application Configuration
app_name = "{{.app_name}}"
version = "{{.version}}"

# Database Configuration
{{if eq .environment "production"}}
database_url = "{{.production_db_url}}"
database_pool_size = 20
{{else if eq .environment "staging"}}
database_url = "{{.staging_db_url}}"
database_pool_size = 10
{{else}}
database_url = "{{.development_db_url}}"
database_pool_size = 5
{{end}}

# Logging Configuration
{{if .debug_enabled}}
log_level = "debug"
{{else}}
log_level = "info"
{{end}}

# SSL Configuration
{{if .ssl_enabled}}
ssl_cert_path = "{{.ssl_cert_path}}"
ssl_key_path = "{{.ssl_key_path}}"
{{end}}
```

### Dynamic Content Generation

#### System Information Script
```bash
# templates/system-info.sh.tmpl
#!/bin/bash

# System Information for {{.facts.hostname}}
echo "=== System Information ==="
echo "Hostname: {{.facts.hostname}}"
echo "OS: {{.facts.os.name}} {{.facts.os.version}}"
echo "Architecture: {{.facts.architecture}}"
echo "CPU Cores: {{.facts.cpu.cores}}"
echo "Memory: {{.facts.memory.total | div 1024 | div 1024}} GB"
echo "Disk Usage: {{.facts.disk.usage_percent}}%"

# Network Information
echo "=== Network Information ==="
{{range .facts.network.interfaces}}
echo "Interface {{.name}}: {{.ip}}"
{{end}}

# Process Information
echo "=== Process Information ==="
echo "Total Processes: {{.facts.processes.total}}"
echo "Running Processes: {{.facts.processes.running}}"

# Custom Information
{{if .custom_info}}
echo "=== Custom Information ==="
{{range $key, $value := .custom_info}}
echo "{{$key}}: {{$value}}"
{{end}}
{{end}}
```

## Security Features

### Security Levels

The templates system supports different security levels:

#### Restricted Level
- Minimal function access
- Strict pattern filtering
- Limited resource usage
- Safe for untrusted templates

#### Standard Level
- Normal function access
- Standard pattern filtering
- Standard resource limits
- Default for most templates

#### Elevated Level
- Extended function access
- Relaxed pattern filtering
- Higher resource limits
- For trusted templates only

#### Trusted Level
- Full function access
- Minimal pattern filtering
- High resource limits
- For fully trusted templates only

### Pattern Filtering

The system includes built-in pattern filtering for dangerous content:

```bash
# Forbidden patterns (examples)
"exec", "system", "eval", "shell"
"password", "secret", "key", "token"
"{{.*os\\.Run.*}}"
"{{.*system.*}}"
"{{.*eval.*}}"
```

### Security Best Practices

1. **Use Appropriate Security Levels**:
   ```bash
   # For untrusted templates
   security_level = "restricted"
   
   # For trusted templates
   security_level = "standard"
   ```

2. **Validate Template Inputs**:
   ```bash
   # Validate required variables
   {{if not .required_variable}}
   {{error "required_variable is required"}}
   {{end}}
   ```

3. **Avoid Dangerous Patterns**:
   ```bash
   # Good: Safe template
   {{.user_name | upper}}
   
   # Bad: Dangerous template
   {{.user_input | exec}}
   ```

4. **Use Function Restrictions**:
   ```bash
   # Restrict function access
   allowed_functions = ["upper", "lower", "trim"]
   restricted_functions = ["exec", "system", "eval"]
   ```

## Performance Optimization

### Template Caching

The system provides automatic template caching:

```bash
# Cache configuration
cache_ttl = 300              # Cache for 5 minutes
max_cache_size = 1000        # Maximum cache entries
```

### Template Optimization

1. **Keep Templates Simple**:
   ```bash
   # Good: Simple template
   {{.name | upper}}
   
   # Bad: Complex template
   {{range .items}}{{range .subitems}}{{.value | complex_function}}{{end}}{{end}}
   ```

2. **Use Efficient Functions**:
   ```bash
   # Good: Efficient function
   {{.items | len}}
   
   # Bad: Inefficient function
   {{range .items}}{{end}} # Count manually
   ```

3. **Minimize Context Resolution**:
   ```bash
   # Good: Resolve once
   {{$fact := .facts.hostname}}
   Hostname: {{$fact}}
   
   # Bad: Resolve multiple times
   Hostname: {{.facts.hostname}}
   ```

### Performance Monitoring

Monitor template performance with metrics:

```bash
# Performance metrics
template_render_time = 150ms
template_render_size = 2048 bytes
cache_hit_rate = 85%
memory_usage = 50MB
```

## Integration Examples

### Facts Integration

Use machine facts in templates:

```bash
# templates/facts-report.tmpl
# System Facts Report for {{.facts.hostname}}

## System Overview
- **Hostname**: {{.facts.hostname}}
- **OS**: {{.facts.os.name}} {{.facts.os.version}}
- **Architecture**: {{.facts.architecture}}
- **Uptime**: {{.facts.uptime | formatTime "2006-01-02 15:04:05"}}

## Hardware Information
- **CPU Cores**: {{.facts.cpu.cores}}
- **CPU Model**: {{.facts.cpu.model}}
- **Memory**: {{.facts.memory.total | div 1024 | div 1024}} GB
- **Disk Space**: {{.facts.disk.total | div 1024 | div 1024 | div 1024}} GB

## Network Information
{{range .facts.network.interfaces}}
### {{.name}}
- **IP Address**: {{.ip}}
- **MAC Address**: {{.mac}}
- **Status**: {{.status}}
- **Speed**: {{.speed}} Mbps
{{end}}
```

### Variables Integration

Use project variables in templates:

```bash
# templates/app-config.tmpl
# Application Configuration

## Basic Configuration
app_name = "{{.variables.app_name}}"
version = "{{.variables.version}}"
environment = "{{.variables.environment}}"

## Database Configuration
{{if .variables.database_enabled}}
database_url = "{{.variables.database_url}}"
database_pool_size = {{.variables.database_pool_size | default 10}}
{{end}}

## API Configuration
api_host = "{{.variables.api_host | default "localhost"}}"
api_port = {{.variables.api_port | default 8080}}
api_timeout = {{.variables.api_timeout | default 30}}

## Feature Flags
{{range $key, $value := .variables.feature_flags}}
{{$key}} = {{$value}}
{{end}}
```

### Machines Integration

Use machine inventory in templates:

```bash
# templates/inventory-report.tmpl
# Machine Inventory Report

## Summary
- **Total Machines**: {{.machines | len}}
- **Online Machines**: {{range .machines}}{{if .online}}{{end}}{{end | len}}
- **Offline Machines**: {{range .machines}}{{if not .online}}{{end}}{{end | len}}

## Machine Details
{{range .machines}}
### {{.hostname}}
- **IP Address**: {{.ip}}
- **Status**: {{if .online}}Online{{else}}Offline{{end}}
- **Tags**: {{.tags | join ", "}}
- **Description**: {{.description}}

{{if .online}}
#### Connection Information
- **User**: {{.user}}
- **Port**: {{.port}}
- **Authentication**: {{.auth_method}}
{{end}}

{{end}}
```

## Troubleshooting

### Common Issues

#### Template Not Found
```bash
Error: template file not found: templates/nginx.conf.tmpl
```
**Solution**: Check if the template file exists in the specified path.

#### Template Syntax Error
```bash
Error: invalid template syntax: unexpected "}" in template
```
**Solution**: Check template syntax and ensure all template tags are properly closed.

#### Missing Required Variables
```bash
Error: template validation failed: required variable "server_name" not provided
```
**Solution**: Provide all required variables in the data file or command line.

#### Security Violation
```bash
Error: security violation: dangerous pattern detected: {{.system}}
```
**Solution**: Remove dangerous patterns from the template or adjust security level.

### Debugging

#### Verbose Output
```bash
# Enable verbose output for debugging
spooky templates render ./myproject templates/nginx.conf.tmpl \
  --verbose
```

#### Preview Mode
```bash
# Preview template rendering without writing files
spooky templates render ./myproject templates/nginx.conf.tmpl \
  --preview
```

#### Dry Run Mode
```bash
# Show what would be rendered without making changes
spooky templates render ./myproject templates/nginx.conf.tmpl \
  --dry-run
```

## Best Practices

### Template Organization

1. **Use Descriptive Names**:
   ```bash
   # Good
   nginx-production.conf.tmpl
   deploy-web-app.sh.tmpl
   
   # Bad
   config.tmpl
   script.tmpl
   ```

2. **Group Related Templates**:
   ```bash
   templates/
   ├── nginx/
   │   ├── nginx.conf.tmpl
   │   └── nginx.service.tmpl
   ├── kubernetes/
   │   ├── deployment.yaml.tmpl
   │   └── service.yaml.tmpl
   └── scripts/
       ├── deploy.sh.tmpl
       └── backup.sh.tmpl
   ```

3. **Use Consistent Naming**:
   ```bash
   # Use consistent naming conventions
   {service}-{environment}.{extension}.tmpl
   nginx-production.conf.tmpl
   nginx-staging.conf.tmpl
   ```

### Template Security

1. **Use Appropriate Security Levels**:
   ```bash
   # For untrusted templates
   security_level = "restricted"
   
   # For trusted templates
   security_level = "standard"
   ```

2. **Validate All Inputs**:
   ```bash
   # Validate required variables
   {{if not .required_variable}}
   {{error "required_variable is required"}}
   {{end}}
   ```

3. **Avoid Dangerous Patterns**:
   ```bash
   # Good: Safe template
   {{.user_name | upper}}
   
   # Bad: Dangerous template
   {{.user_input | exec}}
   ```

### Template Performance

1. **Keep Templates Simple**:
   ```bash
   # Good: Simple template
   {{.name | upper}}
   
   # Bad: Complex template
   {{range .items}}{{range .subitems}}{{.value | complex_function}}{{end}}{{end}}
   ```

2. **Use Efficient Functions**:
   ```bash
   # Good: Efficient function
   {{.items | len}}
   
   # Bad: Inefficient function
   {{range .items}}{{end}} # Count manually
   ```

3. **Minimize Context Resolution**:
   ```bash
   # Good: Resolve once
   {{$fact := .facts.hostname}}
   Hostname: {{$fact}}
   
   # Bad: Resolve multiple times
   Hostname: {{.facts.hostname}}
   ```

### Template Maintenance

1. **Document Template Variables**:
   ```bash
   # Template: nginx.conf.tmpl
   # Required Variables:
   #   - server_name: Domain name for the server
   #   - port: Port number to listen on
   #   - root_path: Document root path
   # Optional Variables:
   #   - ssl_enabled: Enable SSL (default: false)
   #   - ssl_cert_path: SSL certificate path
   #   - ssl_key_path: SSL private key path
   ```

2. **Version Template Metadata**:
   ```bash
   # Template metadata
   name = "nginx-config"
   version = "1.0.0"
   description = "Nginx configuration template"
   author = "spooky-user"
   ```

3. **Test Templates Regularly**:
   ```bash
   # Validate templates
   spooky templates validate ./myproject
   
   # Test template rendering
   spooky templates render ./myproject templates/test.tmpl --preview
   ```

## Related Documentation

- [Templates API Reference](TEMPLATES_API_REFERENCE.md) - Complete API documentation
- [Templates CLI Reference](TEMPLATES_CLI_REFERENCE.md) - CLI command reference
- [Templates System](TEMPLATES_SYSTEM.md) - System overview and architecture
- [Templates Documentation Summary](TEMPLATES_DOCUMENTATION_SUMMARY.md) - Documentation overview

## Support

For help with templates:

1. **Check System Status**: Review the current implementation status
2. **Review Documentation**: Use the appropriate documentation for your needs
3. **Test Templates**: Use spooky commands to validate templates
4. **Check Troubleshooting**: Review troubleshooting guides for common issues
5. **Ask Questions**: Use the project's support channels for specific help
