# Templates System

## Overview

The Templates System provides comprehensive template rendering, validation, and management capabilities for the spooky codebase. It enables dynamic content generation through secure template processing with integration to facts, variables, machines, and other system components.

**Status**: **Implemented** - Complete template system with rendering, validation, security, and CLI integration.

## Architecture

### Core Components

#### Template Manager
- **File**: `internal/templates/manager.go`
- **Purpose**: Central template management with caching, security, and performance optimization
- **Features**:
  - Template loading and caching
  - Template rendering with context resolution
  - Template validation with schema support
  - Template metadata management
  - Security sandboxing and pattern filtering
  - Performance optimization and metrics

#### Template Functions Registry
- **File**: `internal/templates/functions.go`
- **Purpose**: Secure template function management with built-in functions and security restrictions
- **Features**:
  - 50+ built-in template functions
  - Function security manager with restrictions
  - Function result caching
  - Category-based function organization
  - Performance monitoring and limits

#### Template Integration
- **File**: `internal/templates/integration.go`
- **Purpose**: Interface implementation for system integration
- **Features**:
  - LoadTemplate - Load templates from paths
  - RenderTemplate - Render templates with data
  - ValidateTemplate - Validate templates with comprehensive results

### Integration Points

#### Facts Integration
- Provides machine facts data to templates
- Supports custom and system fact collections
- Enables dynamic fact-based template rendering

#### Variables Integration
- Provides project variables to templates
- Supports variable resolution and validation
- Enables dynamic variable-based template rendering

#### Machines Integration
- Provides machine inventory data to templates
- Supports machine-specific template rendering
- Enables dynamic machine-based template rendering

#### Secrets Integration
- Provides secure secret management for templates
- Supports encrypted template variables
- Enables secure template rendering

## Template Types

### Template Structure
```go
type Template struct {
    ID              string                 // Template identifier
    SourcePath      string                 // Template file path
    DestinationPath string                 // Output destination
    Type            string                 // Template type
    Scope           string                 // Usage scope
    SecurityLevel   string                 // Security level
    Engine          string                 // Rendering engine
    Variables       map[string]interface{} // Template variables
    ContextData     map[string]interface{} // Context data
    Functions       map[string]interface{} // Available functions
    Metadata        *TemplateMetadata      // Template metadata
    Content         string                 // Template content
    CreatedAt       time.Time              // Creation timestamp
    UpdatedAt       time.Time              // Update timestamp
}
```

### Template Context
```go
type TemplateContext struct {
    Project     map[string]interface{}   // Project information
    Facts       map[string]interface{}   // Machine facts
    Machines    []map[string]interface{} // Inventory data
    Environment map[string]string        // Environment variables
    CustomData  map[string]interface{}   // Custom data
    Variables   map[string]interface{}   // Project variables
}
```

### Template Metadata
```go
type TemplateMetadata struct {
    Name        string   // Template name
    Description string   // Template description
    Author      string   // Template author
    Version     string   // Template version
    Tags        []string // Template tags
    License     string   // Template license
}
```

## Template Functions

### String Functions
- `upper(s)` - Convert string to uppercase
- `lower(s)` - Convert string to lowercase
- `title(s)` - Convert string to title case
- `trim(s)` - Trim whitespace
- `trimLeft(s)` - Trim left whitespace
- `trimRight(s)` - Trim right whitespace
- `replace(s, old, new, n)` - Replace substrings
- `replaceAll(s, old, new)` - Replace all occurrences
- `split(s, sep)` - Split string by separator
- `join(slice, sep)` - Join slice with separator
- `contains(s, substr)` - Check if string contains substring
- `hasPrefix(s, prefix)` - Check if string has prefix
- `hasSuffix(s, suffix)` - Check if string has suffix
- `repeat(s, n)` - Repeat string n times
- `substr(s, start, length)` - Extract substring
- `len(v)` - Get length of string, array, or map

### Mathematical Functions
- `add(a, b)` - Add two numbers
- `sub(a, b)` - Subtract two numbers
- `mul(a, b)` - Multiply two numbers
- `div(a, b)` - Divide two numbers
- `mod(a, b)` - Modulo operation
- `abs(x)` - Absolute value
- `ceil(x)` - Ceiling function
- `floor(x)` - Floor function
- `round(x)` - Round to nearest integer
- `min(...)` - Minimum value
- `max(...)` - Maximum value
- `pow(x, y)` - Power function
- `sqrt(x)` - Square root

### Array Functions
- `first(slice)` - Get first element
- `last(slice)` - Get last element
- `index(slice, i)` - Get element at index
- `slice(slice, start, end)` - Extract slice
- `append(slice, ...)` - Append elements
- `prepend(slice, ...)` - Prepend elements
- `reverse(slice)` - Reverse array
- `sort(slice)` - Sort array
- `uniq(slice)` - Remove duplicates
- `containsItem(slice, item)` - Check if array contains item

### Hash and Encoding Functions
- `md5(s)` - MD5 hash
- `sha256(s)` - SHA256 hash
- `base64(s)` - Base64 encode
- `base64Decode(s)` - Base64 decode
- `hex(s)` - Hex encode
- `hexDecode(s)` - Hex decode

### Type Conversion Functions
- `toString(v)` - Convert to string
- `toInt(v)` - Convert to integer
- `toFloat(v)` - Convert to float
- `toBool(v)` - Convert to boolean

### JSON Functions
- `toJSON(v)` - Convert to JSON
- `fromJSON(s)` - Parse JSON
- `prettyJSON(v)` - Pretty print JSON

### Date and Time Functions
- `now()` - Current time
- `formatTime(t, layout)` - Format time
- `parseTime(s, layout)` - Parse time
- `addDays(t, days)` - Add days
- `addHours(t, hours)` - Add hours

### Utility Functions
- `default(value, defaultValue)` - Default value
- `coalesce(...)` - First non-empty value
- `ternary(condition, trueValue, falseValue)` - Conditional value
- `regexMatch(pattern, s)` - Regex match
- `regexReplace(pattern, replacement, s)` - Regex replace
- `random(min, max)` - Random number
- `uuid()` - Generate UUID

## Security Features

### Template Sandboxing
- Execution time limits
- Memory usage limits
- Function access restrictions
- Resource monitoring

### Pattern Filtering
- Dangerous pattern detection
- Security violation logging
- Access control enforcement
- Audit trail maintenance

### Security Levels
- **Restricted**: Minimal function access, strict limits
- **Standard**: Normal function access, standard limits
- **Elevated**: Extended function access, relaxed limits
- **Trusted**: Full function access, minimal restrictions

## Performance Features

### Caching
- Template caching with TTL
- Result caching for rendered templates
- Function result caching
- Context resolution caching

### Optimization
- Template compilation optimization
- Parallel processing support
- Memory usage optimization
- Performance metrics collection

### Monitoring
- Template render time tracking
- Memory usage monitoring
- Cache hit rate tracking
- Performance bottleneck detection

## CLI Commands

### Template Rendering
```bash
# Basic template rendering
spooky templates render <project> <template>

# Render with data file
spooky templates render <project> <template> --data <file>

# Render with output file
spooky templates render <project> <template> --output <file>

# Preview mode
spooky templates render <project> <template> --preview

# Dry run mode
spooky templates render <project> <template> --dry-run
```

### Template Validation
```bash
# Validate all templates in project
spooky templates validate <project>

# Validate specific template
spooky templates validate <project> --template <path>
```

### Template Listing
```bash
# List all templates
spooky templates list <project>

# List with specific format
spooky templates list <project> --format json
spooky templates list <project> --format hcl
```

### Template Search
```bash
# Search templates by query
spooky templates search <project> <query>

# Search with tags
spooky templates search <project> <query> --tags <tags>

# Search by category
spooky templates search <project> <query> --category <category>
```

## Configuration

### Template Configuration
```hcl
# templates/config.hcl
template_config {
  # Security settings
  security {
    default_level = "standard"
    max_execution_time = 30000  # milliseconds
    max_memory_usage = 104857600  # 100MB
    restricted_patterns = [
      "{{.*os\\.Run.*}}",
      "{{.*system.*}}",
      "{{.*eval.*}}"
    ]
  }
  
  # Performance settings
  performance {
    cache_ttl = 300  # seconds
    max_cache_size = 1000
    parallel_workers = 4
  }
  
  # Function settings
  functions {
    allowed_categories = ["string", "math", "array", "utility"]
    restricted_functions = ["system", "eval", "exec"]
  }
}
```

### Template Metadata
```hcl
# templates/metadata.hcl
template_metadata {
  name = "nginx-config"
  description = "Nginx configuration template"
  author = "spooky-user"
  version = "1.0.0"
  tags = ["web", "nginx", "config"]
  license = "MIT"
  
  # Template properties
  type = "config"
  scope = "project"
  security_level = "standard"
  engine = "go-template"
  
  # Required variables
  required_variables = [
    "server_name",
    "port",
    "root_path"
  ]
  
  # Output format
  output_format = "nginx.conf"
}
```

## Examples

### Basic Template
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

### Template with Functions
```bash
# templates/deployment.yaml.tmpl
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{.app_name | lower}}
  labels:
    app: {{.app_name | lower}}
    version: {{.version}}
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
```

### Template with Facts Integration
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
```

## Integration Examples

### Actions Integration
```go
// Use templates in actions
action := &spookytypes.Action{
    Name: "deploy-config",
    Template: "templates/nginx.conf.tmpl",
    Variables: map[string]interface{}{
        "server_name": "example.com",
        "port": 80,
        "root_path": "/var/www/html",
    },
}
```

### Facts Integration
```go
// Use facts in templates
facts, err := factsManager.CollectFacts("web-server")
if err != nil {
    return err
}

template, err := templatesManager.LoadTemplate(ctx, "templates/system-info.sh.tmpl")
if err != nil {
    return err
}

result, err := templatesManager.RenderTemplate(ctx, template, map[string]interface{}{
    "facts": facts,
})
```

### Variables Integration
```go
// Use variables in templates
variables, err := variablesManager.LoadVariables(ctx, "")
if err != nil {
    return err
}

template, err := templatesManager.LoadTemplate(ctx, "templates/config.yaml.tmpl")
if err != nil {
    return err
}

result, err := templatesManager.RenderTemplate(ctx, template, map[string]interface{}{
    "variables": variables,
})
```

## Best Practices

### Security
- Use appropriate security levels for templates
- Validate all template inputs
- Avoid dangerous patterns in templates
- Use restricted mode for untrusted templates

### Performance
- Cache frequently used templates
- Use appropriate TTL values
- Monitor template performance
- Optimize template complexity

### Maintainability
- Use descriptive template names
- Document template variables
- Version template metadata
- Use consistent naming conventions

### Integration
- Leverage facts for dynamic content
- Use variables for configuration
- Integrate with machine inventory
- Use secrets for sensitive data

## Troubleshooting

### Common Issues

#### Template Not Found
```bash
# Check template path
ls -la templates/

# Validate template file
spooky templates validate <project> --template <path>
```

#### Template Rendering Errors
```bash
# Check template syntax
spooky templates validate <project> --template <path>

# Check required variables
spooky templates render <project> <template> --preview
```

#### Performance Issues
```bash
# Check template size
wc -l templates/<template>

# Monitor render time
spooky templates render <project> <template> --verbose
```

#### Security Violations
```bash
# Check security level
cat templates/metadata.hcl

# Review restricted patterns
spooky templates validate <project> --strict
```

## API Reference

### TemplatesIntegration Interface
```go
type TemplatesIntegration interface {
    LoadTemplate(ctx context.Context, templatePath string) (*spookytypes.Template, error)
    RenderTemplate(ctx context.Context, template *spookytypes.Template, data map[string]interface{}) (string, error)
    ValidateTemplate(ctx context.Context, template *spookytypes.Template) (*spookytypes.ValidationResult, error)
}
```

### Template Manager Methods
```go
// Load and render templates
LoadTemplate(ctx context.Context, templatePath string) (*spookytypes.Template, error)
RenderTemplate(ctx context.Context, template *spookytypes.Template, data map[string]interface{}) (string, error)
ValidateTemplate(ctx context.Context, template *spookytypes.Template) (*spookytypes.ValidationResult, error)

// Context resolution
ResolveTemplateContext(ctx context.Context, template *spookytypes.Template, data map[string]interface{}) (*spookytypes.TemplateContext, error)

// Function management
RegisterTemplateFunctions(functions map[string]interface{}) error

// Metadata management
GetTemplateMetadata(ctx context.Context, templatePath string) (*spookytypestemplates.TemplateMetadata, error)
```

## Related Documentation

- [Templates API Reference](TEMPLATES_API_REFERENCE.md) - Complete API documentation
- [Template System Design](design/systems/template-system.md) - Design documentation
- [Template Enhanced Composition](design/TEMPLATE_ENHANCED_COMPOSITION.md) - Composition patterns
- [Schema System](../schema-system.md) - Schema validation and configuration
- [Facts System](FACTS_SYSTEM.md) - Facts integration
- [Variables System](VARIABLES_SYSTEM.md) - Variables integration
- [Machines System](MACHINES_SYSTEM.md) - Machines integration
