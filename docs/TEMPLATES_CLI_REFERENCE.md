# Templates CLI Reference

## Overview

This document provides a comprehensive reference for the spooky templates command-line interface (CLI). It covers all template-related commands, options, and usage patterns for template management, rendering, validation, and discovery.

**Status**: **Implemented** - The templates CLI system has comprehensive functionality with all major commands implemented and working.

## Command Structure

The spooky templates CLI follows a consistent command structure:

```bash
spooky templates [global-options] <command> [command-options] <arguments>
```

### Global Options

```bash
--config <file>           # Specify configuration file path
--log-level <level>       # Set logging level (debug, info, warn, error)
--log-format <format>     # Set logging format (json, text)
--quiet                   # Suppress output
--verbose                 # Enable verbose output
--version                 # Show version information
--help                    # Show help information
```

## Template Commands

### Template Rendering

#### `spooky templates render <project> <template>`

Render a template with the given data and output to a file or stdout.

```bash
# Basic template rendering
spooky templates render ./myproject templates/nginx.conf.tmpl

# Render with data file
spooky templates render ./myproject templates/nginx.conf.tmpl --data data/variables.hcl

# Render with output file
spooky templates render ./myproject templates/nginx.conf.tmpl --output /etc/nginx/nginx.conf

# Preview mode (show what would be rendered)
spooky templates render ./myproject templates/nginx.conf.tmpl --preview

# Dry run mode (show what would be rendered without writing files)
spooky templates render ./myproject templates/nginx.conf.tmpl --dry-run
```

**Arguments**:
- `project`: Path to the spooky project directory
- `template`: Path to the template file relative to the project

**Options**:
- `--data <file>`: Data file (JSON/HCL) for template variables
- `--output <file>`: Output file path (default: stdout)
- `--dry-run`: Show what would be rendered without writing files
- `--preview`: Preview the rendered template

**Examples**:
```bash
# Render nginx configuration with variables
spooky templates render ./myproject templates/nginx.conf.tmpl \
  --data data/variables.hcl \
  --output /etc/nginx/nginx.conf

# Preview template rendering
spooky templates render ./myproject templates/deployment.yaml.tmpl \
  --data data/config.json \
  --preview

# Dry run with custom data
spooky templates render ./myproject templates/config.tmpl \
  --data data/production.hcl \
  --dry-run
```

### Template Validation

#### `spooky templates validate <project>`

Validate templates in the project for syntax and security.

```bash
# Validate all templates in project
spooky templates validate ./myproject

# Validate specific template
spooky templates validate ./myproject --template templates/nginx.conf.tmpl

# Validate with verbose output
spooky templates validate ./myproject --verbose
```

**Arguments**:
- `project`: Path to the spooky project directory

**Options**:
- `--template <path>`: Specific template to validate
- `--verbose`: Show detailed validation information

**Examples**:
```bash
# Validate all templates
spooky templates validate ./myproject

# Validate specific template
spooky templates validate ./myproject --template templates/nginx.conf.tmpl

# Validate with detailed output
spooky templates validate ./myproject --template templates/deployment.yaml.tmpl --verbose
```

### Template Listing

#### `spooky templates list <project>`

List all templates in the project with metadata and information.

```bash
# List all templates
spooky templates list ./myproject

# List templates in JSON format
spooky templates list ./myproject --format json

# List templates in HCL format
spooky templates list ./myproject --format hcl

# List templates in table format
spooky templates list ./myproject --format table
```

**Arguments**:
- `project`: Path to the spooky project directory

**Options**:
- `--format <format>`: Output format (table, json, hcl) [default: table]

**Examples**:
```bash
# List all templates
spooky templates list ./myproject

# List templates in JSON format
spooky templates list ./myproject --format json

# List templates in HCL format
spooky templates list ./myproject --format hcl
```

### Template Search

#### `spooky templates search <project> <query>`

Search templates in the project by name, description, or tags.

```bash
# Search templates by query
spooky templates search ./myproject "nginx"

# Search with tags filter
spooky templates search ./myproject "config" --tags web,nginx

# Search by category
spooky templates search ./myproject "deployment" --category kubernetes
```

**Arguments**:
- `project`: Path to the spooky project directory
- `query`: Search query string

**Options**:
- `--tags <tags>`: Filter by tags (comma-separated)
- `--category <category>`: Filter by category

**Examples**:
```bash
# Search for web templates
spooky templates search ./myproject "web" --tags nginx,apache

# Search for deployment templates
spooky templates search ./myproject "deployment" --category kubernetes

# Search for configuration templates
spooky templates search ./myproject "config" --tags production
```

## Template File Formats

### Template Files

Templates use Go template syntax with `.tmpl` extension:

```bash
# templates/nginx.conf.tmpl
server {
    listen {{.port}};
    server_name {{.server_name}};
    root {{.root_path}};
    
    location / {
        try_files $uri $uri/ /index.html;
    }
    
    # Logging
    access_log /var/log/nginx/{{.server_name}}.access.log;
    error_log /var/log/nginx/{{.server_name}}.error.log;
}
```

### Data Files

Data files can be in HCL or JSON format:

#### HCL Data File
```hcl
# data/variables.hcl
server_name = "example.com"
port = 80
root_path = "/var/www/html"
environment = "production"
```

#### JSON Data File
```json
{
  "server_name": "example.com",
  "port": 80,
  "root_path": "/var/www/html",
  "environment": "production"
}
```

### Template Metadata

Template metadata can be defined in `.meta` files:

```hcl
# templates/nginx.conf.tmpl.meta
name = "nginx-config"
description = "Nginx configuration template"
author = "spooky-user"
version = "1.0.0"
tags = ["web", "nginx", "config"]
license = "MIT"

type = "config"
scope = "project"
security_level = "standard"
engine = "go-template"

required_variables = [
  "server_name",
  "port",
  "root_path"
]

output_format = "nginx.conf"
```

## Template Functions

### Available Functions

The templates system provides a comprehensive set of built-in functions:

#### String Functions
```bash
{{upper "hello"}}           # HELLO
{{lower "WORLD"}}           # world
{{title "hello world"}}     # Hello World
{{trim "  hello  "}}        # hello
{{replace "hello" "l" "L" 1}} # heLlo
{{split "a,b,c" ","}}       # [a b c]
{{join .items ","}}         # item1,item2,item3
{{contains "hello" "ll"}}   # true
{{len "hello"}}             # 5
```

#### Mathematical Functions
```bash
{{add 1 2}}                 # 3
{{sub 5 2}}                 # 3
{{mul 3 4}}                 # 12
{{div 10 2}}                # 5
{{mod 7 3}}                 # 1
{{abs -5}}                  # 5
{{min 1 2 3}}               # 1
{{max 1 2 3}}               # 3
{{pow 2 3}}                 # 8
{{sqrt 16}}                 # 4
```

#### Array Functions
```bash
{{first .items}}            # First item
{{last .items}}             # Last item
{{index .items 0}}          # Item at index 0
{{slice .items 1 3}}        # Items from index 1 to 3
{{append .items "new"}}     # Add item to array
{{sort .items}}             # Sort array
{{uniq .items}}             # Remove duplicates
```

#### Hash and Encoding Functions
```bash
{{md5 "hello"}}             # 5d41402abc4b2a76b9719d911017c592
{{sha256 "hello"}}          # 2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824
{{base64 "hello"}}          # aGVsbG8=
{{base64Decode "aGVsbG8="}} # hello
{{hex "hello"}}             # 68656c6c6f
{{hexDecode "68656c6c6f"}}  # hello
```

#### Type Conversion Functions
```bash
{{toString 123}}            # "123"
{{toInt "123"}}             # 123
{{toFloat "123.45"}}        # 123.45
{{toBool "true"}}           # true
```

#### JSON Functions
```bash
{{toJSON .data}}            # {"key":"value"}
{{fromJSON '{"key":"value"}'}} # map[key:value]
{{prettyJSON .data}}        # Pretty formatted JSON
```

#### Date and Time Functions
```bash
{{now}}                     # Current time
{{formatTime .timestamp "2006-01-02"}} # Formatted date
{{parseTime "2024-01-01" "2006-01-02"}} # Parsed time
{{addDays .date 7}}         # Add 7 days
{{addHours .time 24}}       # Add 24 hours
```

#### Utility Functions
```bash
{{default .value "default"}} # Default value if empty
{{coalesce .a .b .c}}       # First non-empty value
{{ternary .condition "yes" "no"}} # Conditional value
{{regexMatch "^[a-z]+$" "hello"}} # true
{{regexReplace "l" "L" "hello"}} # heLLo
{{random 1 100}}            # Random number 1-100
{{uuid}}                    # Generate UUID
```

## Template Context

### Available Context Data

Templates have access to rich context data:

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

## Examples

### Nginx Configuration Template

```bash
# templates/nginx.conf.tmpl
server {
    listen {{.port}};
    server_name {{.server_name}};
    root {{.root_path}};
    
    # SSL configuration
    {{if .ssl_enabled}}
    listen {{.port}} ssl;
    ssl_certificate {{.ssl_cert_path}};
    ssl_certificate_key {{.ssl_key_path}};
    {{end}}
    
    location / {
        try_files $uri $uri/ /index.html;
    }
    
    # API proxy
    {{if .api_enabled}}
    location /api/ {
        proxy_pass {{.api_backend}};
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
    {{end}}
    
    # Logging
    access_log /var/log/nginx/{{.server_name}}.access.log;
    error_log /var/log/nginx/{{.server_name}}.error.log;
}
```

**Usage**:
```bash
# Render nginx configuration
spooky templates render ./myproject templates/nginx.conf.tmpl \
  --data data/nginx.hcl \
  --output /etc/nginx/sites-available/example.com
```

### Kubernetes Deployment Template

```bash
# templates/deployment.yaml.tmpl
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{.app_name | lower}}
  labels:
    app: {{.app_name | lower}}
    version: {{.version}}
    environment: {{.environment}}
spec:
  replicas: {{.replicas | default 3}}
  selector:
    matchLabels:
      app: {{.app_name | lower}}
  template:
    metadata:
      labels:
        app: {{.app_name | lower}}
    spec:
      containers:
      - name: {{.app_name | lower}}
        image: {{.image}}
        ports:
        - containerPort: {{.port}}
        env:
        - name: APP_NAME
          value: {{.app_name | upper}}
        - name: VERSION
          value: {{.version}}
        - name: ENVIRONMENT
          value: {{.environment}}
        {{if .secrets_enabled}}
        - name: DATABASE_URL
          valueFrom:
            secretKeyRef:
              name: {{.app_name}}-secrets
              key: database-url
        {{end}}
        resources:
          requests:
            memory: {{.memory_request | default "128Mi"}}
            cpu: {{.cpu_request | default "100m"}}
          limits:
            memory: {{.memory_limit | default "256Mi"}}
            cpu: {{.cpu_limit | default "200m"}}
```

**Usage**:
```bash
# Render deployment configuration
spooky templates render ./myproject templates/deployment.yaml.tmpl \
  --data data/production.hcl \
  --output k8s/deployment.yaml
```

### System Information Script

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

**Usage**:
```bash
# Render system information script
spooky templates render ./myproject templates/system-info.sh.tmpl \
  --data data/system.hcl \
  --output scripts/system-info.sh

# Make executable
chmod +x scripts/system-info.sh
```

## Error Handling

### Common Errors

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
  --data data/variables.hcl \
  --verbose
```

#### Preview Mode
```bash
# Preview template rendering without writing files
spooky templates render ./myproject templates/nginx.conf.tmpl \
  --data data/variables.hcl \
  --preview
```

#### Dry Run Mode
```bash
# Show what would be rendered without making changes
spooky templates render ./myproject templates/nginx.conf.tmpl \
  --data data/variables.hcl \
  --dry-run
```

## Best Practices

### Template Organization

- Organize templates by purpose and scope
- Use descriptive template names
- Group related templates in subdirectories
- Use consistent naming conventions

### Template Security

- Use appropriate security levels
- Validate all template inputs
- Avoid dangerous patterns
- Use restricted mode for untrusted templates

### Template Performance

- Keep templates simple and focused
- Use caching for frequently rendered templates
- Optimize template complexity
- Monitor template rendering performance

### Template Maintenance

- Document template variables
- Version template metadata
- Use consistent formatting
- Test templates regularly

## Related Documentation

- [Templates System](TEMPLATES_SYSTEM.md) - System overview and architecture
- [Templates API Reference](TEMPLATES_API_REFERENCE.md) - Complete API documentation
- [Template System Design](design/systems/template-system.md) - Design documentation
- [Schema System](../schema-system.md) - Schema validation and configuration
- [Facts System](FACTS_SYSTEM.md) - Facts integration
- [Variables System](VARIABLES_SYSTEM.md) - Variables integration
- [Machines System](MACHINES_SYSTEM.md) - Machines integration
