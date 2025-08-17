# Templates System

## Overview

The spooky templates system provides comprehensive template management including loading, rendering, validation, and discovery capabilities. Templates are used to generate dynamic content for configuration files, scripts, and other resources based on project variables, facts, and machine data.

## Related Systems

This system integrates with and depends on several other spooky systems:

- **[Variables System](VARIABLES_SYSTEM.md)** - Templates use variables for dynamic content rendering
- **[Facts System](FACTS_SYSTEM.md)** - Templates can use facts as data for rendering
- **[Actions System](ACTIONS_SYSTEM.md)** - Templates are used by actions for dynamic content generation
- **[Logging System](LOGGING_SYSTEM.md)** - Template rendering generates comprehensive logs for monitoring and debugging
- **[Projects System](PROJECTS_SYSTEM.md)** - Templates are organized within projects
- **[Machines System](MACHINES_SYSTEM.md)** - Templates can use machine information for rendering
- **[Integrations System](INTEGRATIONS_SYSTEM.md)** - Templates integrate with other systems through the IntegrationManager
- **[SSH System](SSH_SYSTEM.md)** - Templates can be rendered and deployed via SSH

## Core Concepts

### Template Definition
Templates are defined using Go's text/template package and can include:

- **Variable Substitution**: Replace variables with actual values
- **Conditional Logic**: Include/exclude content based on conditions
- **Loops and Iteration**: Generate content for multiple items
- **Function Calls**: Use built-in and custom functions
- **Nested Templates**: Include other templates

### Template Rendering
The system supports multiple rendering modes:

- **Variable-based Rendering**: Render with project variables
- **Fact-based Rendering**: Render with machine facts
- **Context-based Rendering**: Render with combined context
- **Validation Rendering**: Render for validation purposes

### Template Validation
Templates are validated before use:

- **Syntax Validation**: Ensures valid template syntax
- **Variable Validation**: Validates required variables
- **Function Validation**: Validates template functions
- **Context Validation**: Validates rendering context

## CLI Commands

### Template Management Commands

#### `spooky templates render [project] [template]`
Render a template with the given data and output to a file.

**Flags:**
- `--data` - Data file (JSON)
- `--output` - Output file path
- `--dry-run` - Show rendering result without writing to file
- `--preview` - Preview rendering result

**Examples:**
```bash
# Render template with data file
spooky templates render ./my-project templates/nginx.conf.tmpl --data data.json --output nginx.conf

# Render with preview
spooky templates render ./my-project templates/nginx.conf.tmpl --data data.json --preview

# Render with dry-run
spooky templates render ./my-project templates/nginx.conf.tmpl --data data.json --dry-run
```

#### `spooky templates validate [project-path]`
Validate template syntax and variables.

**Flags:**
- `--template` - Validate specific template file
- `--data` - Data file for validation context

**Examples:**
```bash
# Validate all templates
spooky templates validate ./my-project

# Validate specific template
spooky templates validate ./my-project --template templates/nginx.conf.tmpl

# Validate with data context
spooky templates validate ./my-project --template templates/nginx.conf.tmpl --data data.json
```

#### `spooky templates list [project-path]`
List all available templates in the project.

**Examples:**
```bash
# List all templates
spooky templates list ./my-project
```

#### `spooky templates search [project-path]`
Search for templates by content or metadata.

**Flags:**
- `--query` - Search query
- `--type` - Template type filter

**Examples:**
```bash
# Search for templates
spooky templates search ./my-project --query "nginx"

# Search by type
spooky templates search ./my-project --type "config"
```

## Template Configuration

### Basic Template Structure
```hcl
# templates/nginx.conf.tmpl
server {
    listen {{ .port | default 80 }};
    server_name {{ .server_name }};
    
    location / {
        root {{ .document_root | default "/var/www/html" }};
        index index.html index.htm;
    }
    
    {{ if .ssl_enabled }}
    listen 443 ssl;
    ssl_certificate {{ .ssl_certificate }};
    ssl_certificate_key {{ .ssl_certificate_key }};
    {{ end }}
}
```

### Template with Functions
```hcl
# templates/deploy.sh.tmpl
#!/bin/bash

# Deploy script for {{ .app_name }}
set -e

echo "Deploying {{ .app_name }} version {{ .version }}"

{{ range .servers }}
echo "Deploying to {{ .hostname }}"
ssh {{ .user }}@{{ .hostname }} "cd {{ .deploy_path }} && git pull origin {{ $.branch }}"
{{ end }}

echo "Deployment completed successfully"
```

### Template with Conditions
```hcl
# templates/config.yaml.tmpl
apiVersion: v1
kind: ConfigMap
metadata:
  name: {{ .app_name }}-config
data:
  {{ if eq .environment "production" }}
  log_level: "info"
  debug: "false"
  {{ else }}
  log_level: "debug"
  debug: "true"
  {{ end }}
  
  database_url: {{ .database_url }}
  redis_url: {{ .redis_url }}
```

## Template Rendering Process

### Rendering Workflow
1. **Load Template**: Load template file from project
2. **Parse Template**: Parse template syntax
3. **Load Context**: Load variables, facts, and machine data
4. **Validate Context**: Validate required variables
5. **Render Template**: Generate final content
6. **Validate Output**: Validate rendered content
7. **Write Output**: Write to file or display

### Context Loading
The system loads context from multiple sources:

```bash
# Load context from data file
spooky templates render ./my-project template.tmpl --data context.json

# Load context from project variables
spooky templates render ./my-project template.tmpl

# Load context from facts
spooky templates render ./my-project template.tmpl --facts facts.json
```

### Variable Resolution
Variables are resolved in the following order:

1. **Data File Variables**: Variables from `--data` file
2. **Project Variables**: Variables from `variables.hcl`
3. **Environment Variables**: System environment variables
4. **Default Values**: Template default values

## Template Functions

### Built-in Functions
The system provides several built-in functions:

**String Functions:**
```hcl
{{ .name | upper }}           # Convert to uppercase
{{ .name | lower }}           # Convert to lowercase
{{ .name | title }}           # Convert to title case
{{ .name | trim }}            # Trim whitespace
{{ .name | replace "old" "new" }}  # Replace text
```

**Number Functions:**
```hcl
{{ .count | add 1 }}          # Add numbers
{{ .count | sub 1 }}          # Subtract numbers
{{ .count | mul 2 }}          # Multiply numbers
{{ .count | div 2 }}          # Divide numbers
```

**List Functions:**
```hcl
{{ .items | len }}            # Get list length
{{ .items | first }}          # Get first item
{{ .items | last }}           # Get last item
{{ .items | join "," }}       # Join list items
```

**Conditional Functions:**
```hcl
{{ .value | default "default" }}  # Default value
{{ .value | empty | not }}    # Check if not empty
{{ .value | eq "test" }}      # Equality check
```

### Custom Functions
The system supports custom template functions:

```hcl
# Custom function definition
{{ define "format_size" }}
{{ if gt . 1073741824 }}
{{ div . 1073741824 }}GB
{{ else if gt . 1048576 }}
{{ div . 1048576 }}MB
{{ else if gt . 1024 }}
{{ div . 1024 }}KB
{{ else }}
{{ . }}B
{{ end }}
{{ end }}

# Usage
{{ .file_size | format_size }}
```

## Template Validation

### Validation Process
The template validation system performs:

1. **Syntax Validation**: Validates template syntax
2. **Variable Validation**: Checks for undefined variables
3. **Function Validation**: Validates function calls
4. **Context Validation**: Validates rendering context
5. **Output Validation**: Validates rendered output

### Validation Output
```bash
# Validate template
spooky templates validate ./my-project --template templates/nginx.conf.tmpl
```

**Output:**
```
🔍 Validating template: templates/nginx.conf.tmpl
✅ Template validation passed
📋 Syntax validation: Valid ✅
📋 Variable validation: Valid ✅
📋 Function validation: Valid ✅
```

### Context Validation
```bash
# Validate with context
spooky templates validate ./my-project --template templates/nginx.conf.tmpl --data context.json
```

**Output:**
```
🔍 Validating template with context: templates/nginx.conf.tmpl
✅ Template validation passed
📋 Context validation: Valid ✅
📋 Variable resolution: Valid ✅
📋 Output validation: Valid ✅
```

## Template Discovery

### Template Listing
```bash
# List all templates
spooky templates list ./my-project
```

**Output:**
```
📋 Available templates in ./my-project:
templates/
├── nginx.conf.tmpl
├── deploy.sh.tmpl
├── config.yaml.tmpl
└── docker-compose.yml.tmpl
```

### Template Search
```bash
# Search for templates
spooky templates search ./my-project --query "nginx"
```

**Output:**
```
🔍 Search results for "nginx":
templates/nginx.conf.tmpl - Nginx configuration template
templates/nginx-ssl.conf.tmpl - Nginx SSL configuration template
```

## Template Context

### Context Structure
The template context includes:

```json
{
  "variables": {
    "app_name": "my-app",
    "version": "1.0.0",
    "environment": "production"
  },
  "facts": {
    "machine": "web-server",
    "os": "Ubuntu 22.04",
    "memory": "8GB"
  },
  "machines": [
    {
      "hostname": "web-server",
      "ip": "192.168.1.10",
      "user": "admin"
    }
  ],
  "project": {
    "name": "my-project",
    "path": "./my-project"
  }
}
```

### Context Loading
```bash
# Load context from multiple sources
spooky templates render ./my-project template.tmpl \
  --data variables.json \
  --facts facts.json \
  --machines machines.json
```

## Template Examples

### Configuration Templates
```hcl
# templates/application.conf.tmpl
[application]
name = {{ .app_name }}
version = {{ .version }}
environment = {{ .environment }}

[database]
host = {{ .database.host }}
port = {{ .database.port | default 5432 }}
name = {{ .database.name }}
user = {{ .database.user }}

[logging]
level = {{ .log_level | default "info" }}
file = {{ .log_file | default "/var/log/app.log" }}
```

### Script Templates
```hcl
# templates/backup.sh.tmpl
#!/bin/bash

# Backup script for {{ .app_name }}
set -e

BACKUP_DIR="{{ .backup_dir | default "/backups" }}"
DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="$BACKUP_DIR/{{ .app_name }}_$DATE.tar.gz"

echo "Creating backup: $BACKUP_FILE"

{{ range .databases }}
echo "Backing up database: {{ .name }}"
pg_dump -h {{ .host }} -p {{ .port | default 5432 }} -U {{ .user }} {{ .name }} > {{ .name }}.sql
{{ end }}

tar -czf "$BACKUP_FILE" *.sql
rm *.sql

echo "Backup completed: $BACKUP_FILE"
```

### Docker Templates
```hcl
# templates/docker-compose.yml.tmpl
version: '3.8'

services:
  {{ .app_name }}:
    image: {{ .image }}:{{ .version }}
    container_name: {{ .app_name }}
    ports:
      - "{{ .port }}:{{ .port }}"
    environment:
      - NODE_ENV={{ .environment }}
      - DATABASE_URL={{ .database_url }}
    volumes:
      - {{ .data_dir }}:/app/data
    {{ if .restart_policy }}
    restart: {{ .restart_policy }}
    {{ end }}
```

## Integration with Other Systems

### Variables Integration
Templates use project variables for dynamic content:

```bash
# Render with project variables
spooky templates render ./my-project template.tmpl
```

### Facts Integration
Templates can use machine facts:

```hcl
# Template using facts
server {
    listen {{ .facts.port | default 80 }};
    server_name {{ .facts.hostname }};
}
```

### Actions Integration
Templates are used in actions for dynamic scripts:

```hcl
# Action using template
action "deploy" {
  description = "Deploy application"
  
  template {
    source = "templates/deploy.sh.tmpl"
    destination = "/tmp/deploy.sh"
    permissions = "0755"
  }
  
  command = "/tmp/deploy.sh"
}
```

## Error Handling

### Common Template Errors
- **Syntax Errors**: Invalid template syntax
- **Variable Errors**: Undefined or invalid variables
- **Function Errors**: Invalid function calls
- **Context Errors**: Missing or invalid context data

### Error Recovery
```bash
# Validate template for errors
spooky templates validate ./my-project --template template.tmpl

# Check template syntax
spooky templates render ./my-project template.tmpl --preview
```

## Troubleshooting

### Common Issues

#### Template Syntax Issues
```bash
# Validate template syntax
spooky templates validate ./my-project --template template.tmpl

# Check for syntax errors
spooky templates render ./my-project template.tmpl --preview
```

#### Variable Issues
```bash
# Check available variables
spooky variables list ./my-project

# Validate variable context
spooky templates validate ./my-project --template template.tmpl --data context.json
```

#### Rendering Issues
```bash
# Preview rendering result
spooky templates render ./my-project template.tmpl --preview

# Check template functions
spooky templates validate ./my-project --template template.tmpl
```

### Debug Information
```bash
# Enable debug logging
export SPOOKY_LOG_LEVEL=debug

# Render with debug output
spooky templates render ./my-project template.tmpl --preview
```

## Best Practices

### Template Design
1. **Use Descriptive Names**: Choose clear, descriptive template names
2. **Include Comments**: Add comments for complex templates
3. **Use Default Values**: Provide sensible default values
4. **Validate Inputs**: Validate template inputs
5. **Keep Templates Simple**: Avoid overly complex templates

### Variable Management
1. **Use Consistent Naming**: Use consistent variable naming conventions
2. **Document Variables**: Document all template variables
3. **Provide Defaults**: Provide default values for optional variables
4. **Validate Variables**: Validate required variables
5. **Use Type Safety**: Use appropriate variable types

### Security Practices
1. **Validate Inputs**: Validate all template inputs
2. **Escape Output**: Escape user-provided content
3. **Limit Functions**: Limit available template functions
4. **Audit Templates**: Regularly audit template content
5. **Secure Context**: Secure template context data

## Future Enhancements

### Planned Features
- **Template Caching**: Cache parsed templates for performance
- **Template Versioning**: Version control for templates
- **Template Inheritance**: Template inheritance and composition
- **Template Testing**: Automated template testing
- **Template Analytics**: Template usage analytics

### Extension Points
- **Custom Functions**: User-defined template functions
- **Template Plugins**: Pluggable template engines
- **External Integrations**: Integration with external template systems
- **Template APIs**: REST API for template management
- **Template Webhooks**: Webhook notifications for template events
