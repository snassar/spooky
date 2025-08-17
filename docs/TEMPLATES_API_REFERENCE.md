# Templates System API Reference

## Overview

This document provides a comprehensive API reference for the spooky templates system. It covers all interfaces, types, methods, and implementation details for developers working with the templates system.

**Status**: **Implemented** - The templates system provides comprehensive functionality for template rendering, validation, and management.

## Core Interfaces

### TemplatesIntegration Interface

The `TemplatesIntegration` interface provides the primary entry point for templates operations:

```go
type TemplatesIntegration interface {
    // LoadTemplate loads a template from the given path
    LoadTemplate(ctx context.Context, templatePath string) (*spookytypes.Template, error)
    
    // RenderTemplate renders a template with the given data
    RenderTemplate(ctx context.Context, template *spookytypes.Template, data map[string]interface{}) (string, error)
    
    // ValidateTemplate validates a template
    ValidateTemplate(ctx context.Context, template *spookytypes.Template) (*spookytypes.ValidationResult, error)
}
```

**Implementation Status**: ✅ **Implemented** - Complete functionality for template management and rendering

## Core Types

### Template

```go
type Template struct {
    // Template identification
    ID string `json:"id" hcl:"id"`
    
    // Template source location
    SourcePath string `json:"source_path" hcl:"source_path"`
    
    // Template destination
    DestinationPath string `json:"destination_path,omitempty" hcl:"destination_path,optional"`
    
    // Template type classification
    Type string `json:"type" hcl:"type"`
    
    // Template scope
    Scope string `json:"scope" hcl:"scope"`
    
    // Template security level
    SecurityLevel string `json:"security_level" hcl:"security_level"`
    
    // Template rendering engine
    Engine string `json:"engine" hcl:"engine"`
    
    // Template variables
    Variables map[string]interface{} `json:"variables,omitempty" hcl:"variables,optional"`
    
    // Template context data
    ContextData map[string]interface{} `json:"context_data,omitempty" hcl:"context_data,optional"`
    
    // Template functions and restrictions
    Functions map[string]interface{} `json:"functions,omitempty" hcl:"functions,optional"`
    
    // Template metadata
    Metadata *TemplateMetadata `json:"metadata,omitempty" hcl:"metadata,optional"`
    
    // Template content
    Content string `json:"content,omitempty" hcl:"content,optional"`
    
    // Creation timestamp
    CreatedAt time.Time `json:"created_at" hcl:"created_at"`
    
    // Last update timestamp
    UpdatedAt time.Time `json:"updated_at" hcl:"updated_at"`
}
```

### TemplateMetadata

```go
type TemplateMetadata struct {
    Name        string   `json:"name,omitempty" hcl:"name,optional"`
    Description string   `json:"description,omitempty" hcl:"description,optional"`
    Author      string   `json:"author,omitempty" hcl:"author,optional"`
    Version     string   `json:"version,omitempty" hcl:"version,optional"`
    Tags        []string `json:"tags,omitempty" hcl:"tags,optional"`
    License     string   `json:"license,omitempty" hcl:"license,optional"`
}
```

### TemplateContext

```go
type TemplateContext struct {
    // Project information
    Project map[string]interface{} `json:"project,omitempty"`
    
    // Machine facts
    Facts map[string]interface{} `json:"facts,omitempty"`
    
    // Inventory information
    Machines []map[string]interface{} `json:"machines,omitempty"`
    
    // Environment variables
    Environment map[string]string `json:"environment,omitempty"`
    
    // Custom data
    CustomData map[string]interface{} `json:"custom_data,omitempty"`
    
    // Variables
    Variables map[string]interface{} `json:"variables,omitempty"`
}
```

## Manager Implementation

### Template Manager

The `Manager` struct provides comprehensive template management functionality:

```go
type Manager struct {
    logger             spookytypeslogging.Logger
    cache              TemplateCache
    functions          TemplateFunctionRegistry
    contextResolver    TemplateContextResolver
    metadataManager    TemplateMetadataManager
    validator          TemplateValidator
    securityManager    TemplateSecurityManager
    performanceManager TemplatePerformanceManager
    mu                 sync.RWMutex
}
```

### Key Methods

#### LoadTemplate
```go
func (m *Manager) LoadTemplate(ctx context.Context, templatePath string) (*spookytypes.Template, error)
```
Loads a template from the given path with caching and validation.

**Parameters**:
- `ctx`: Context for the operation
- `templatePath`: Path to the template file

**Returns**:
- `*spookytypes.Template`: Loaded template
- `error`: Error if loading fails

**Example**:
```go
template, err := manager.LoadTemplate(ctx, "templates/nginx.conf.tmpl")
if err != nil {
    return fmt.Errorf("failed to load template: %w", err)
}
```

#### RenderTemplate
```go
func (m *Manager) RenderTemplate(ctx context.Context, tmplData *spookytypes.Template, data map[string]interface{}) (string, error)
```
Renders a template with enhanced features including context resolution and security validation.

**Parameters**:
- `ctx`: Context for the operation
- `tmplData`: Template to render
- `data`: Data to use for rendering

**Returns**:
- `string`: Rendered template content
- `error`: Error if rendering fails

**Example**:
```go
result, err := manager.RenderTemplate(ctx, template, map[string]interface{}{
    "server_name": "example.com",
    "port": 80,
})
if err != nil {
    return fmt.Errorf("failed to render template: %w", err)
}
```

#### ValidateTemplate
```go
func (m *Manager) ValidateTemplate(ctx context.Context, template *spookytypes.Template) (*spookytypes.ValidationResult, error)
```
Validates a template with comprehensive validation including syntax, security, and schema validation.

**Parameters**:
- `ctx`: Context for the operation
- `template`: Template to validate

**Returns**:
- `*spookytypes.ValidationResult`: Validation results
- `error`: Error if validation fails

**Example**:
```go
result, err := manager.ValidateTemplate(ctx, template)
if err != nil {
    return fmt.Errorf("validation failed: %w", err)
}
if !result.Valid {
    return fmt.Errorf("template validation failed: %v", result.Errors)
}
```

#### ResolveTemplateContext
```go
func (m *Manager) ResolveTemplateContext(ctx context.Context, template *spookytypes.Template, data map[string]interface{}) (*spookytypes.TemplateContext, error)
```
Resolves template context with facts, variables, and machines data.

**Parameters**:
- `ctx`: Context for the operation
- `template`: Template for context resolution
- `data`: Additional data for context

**Returns**:
- `*spookytypes.TemplateContext`: Resolved context
- `error`: Error if resolution fails

**Example**:
```go
context, err := manager.ResolveTemplateContext(ctx, template, data)
if err != nil {
    return fmt.Errorf("failed to resolve context: %w", err)
}
```

#### RegisterTemplateFunctions
```go
func (m *Manager) RegisterTemplateFunctions(functions map[string]interface{}) error
```
Registers custom template functions with security validation.

**Parameters**:
- `functions`: Map of function names to function implementations

**Returns**:
- `error`: Error if registration fails

**Example**:
```go
customFunctions := map[string]interface{}{
    "customFunc": func(s string) string { return "custom: " + s },
}
err := manager.RegisterTemplateFunctions(customFunctions)
if err != nil {
    return fmt.Errorf("failed to register functions: %w", err)
}
```

#### GetTemplateMetadata
```go
func (m *Manager) GetTemplateMetadata(ctx context.Context, templatePath string) (*spookytypestemplates.TemplateMetadata, error)
```
Gets template metadata with caching and validation.

**Parameters**:
- `ctx`: Context for the operation
- `templatePath`: Path to the template

**Returns**:
- `*spookytypestemplates.TemplateMetadata`: Template metadata
- `error`: Error if retrieval fails

**Example**:
```go
metadata, err := manager.GetTemplateMetadata(ctx, "templates/nginx.conf.tmpl")
if err != nil {
    return fmt.Errorf("failed to get metadata: %w", err)
}
```

## Function Registry

### BuiltInFunctions

The `BuiltInFunctions` struct provides a comprehensive set of secure template functions:

```go
type BuiltInFunctions struct {
    logger spookytypeslogging.Logger
}
```

### Function Categories

#### String Functions
```go
// String manipulation functions
"upper":      strings.ToUpper,
"lower":      strings.ToLower,
"title":      strings.Title,
"trim":       strings.TrimSpace,
"trimLeft":   strings.TrimLeft,
"trimRight":  strings.TrimRight,
"replace":    strings.Replace,
"replaceAll": strings.ReplaceAll,
"split":      strings.Split,
"join":       strings.Join,
"contains":   strings.Contains,
"hasPrefix":  strings.HasPrefix,
"hasSuffix":  strings.HasSuffix,
"repeat":     strings.Repeat,
"substr":     substring,
"len":        length,
```

#### Mathematical Functions
```go
// Mathematical functions
"add":   add,
"sub":   sub,
"mul":   mul,
"div":   div,
"mod":   mod,
"abs":   math.Abs,
"ceil":  math.Ceil,
"floor": math.Floor,
"round": math.Round,
"min":   min,
"max":   max,
"pow":   math.Pow,
"sqrt":  math.Sqrt,
```

#### Array Functions
```go
// Array manipulation functions
"first":        first,
"last":         last,
"index":        index,
"slice":        slice,
"append":       append,
"prepend":      prepend,
"reverse":      reverse,
"sort":         sort,
"uniq":         uniq,
"containsItem": containsItem,
```

#### Hash and Encoding Functions
```go
// Hash and encoding functions
"md5":          md5,
"sha256":       sha256,
"base64":       base64,
"base64Decode": base64Decode,
"hex":          hex,
"hexDecode":    hexDecode,
```

#### Type Conversion Functions
```go
// Type conversion functions
"toString": toString,
"toInt":    toInt,
"toFloat":  toFloat,
"toBool":   toBool,
```

#### JSON Functions
```go
// JSON functions
"toJSON":     toJSON,
"fromJSON":   fromJSON,
"prettyJSON": prettyJSON,
```

#### Date and Time Functions
```go
// Date and time functions
"now":        now,
"formatTime": formatTime,
"parseTime":  parseTime,
"addDays":    addDays,
"addHours":   addHours,
```

#### Utility Functions
```go
// Utility functions
"default":      defaultValue,
"coalesce":     coalesce,
"ternary":      ternary,
"regexMatch":   regexMatch,
"regexReplace": regexReplace,
"random":       random,
"uuid":         uuid,
```

### Function Security Manager

The `FunctionSecurityManager` provides security restrictions for template functions:

```go
type FunctionSecurityManager struct {
    restrictedMode     bool
    allowedFunctions   map[string]bool
    forbiddenFunctions map[string]bool
    mu                 sync.RWMutex
}
```

#### Key Methods

##### SetRestrictedMode
```go
func (f *FunctionSecurityManager) SetRestrictedMode(restricted bool)
```
Sets the restricted mode for function access.

##### AllowFunction
```go
func (f *FunctionSecurityManager) AllowFunction(name string)
```
Explicitly allows a specific function.

##### ForbidFunction
```go
func (f *FunctionSecurityManager) ForbidFunction(name string)
```
Explicitly forbids a specific function.

##### IsFunctionAllowed
```go
func (f *FunctionSecurityManager) IsFunctionAllowed(name string) bool
```
Checks if a function is allowed in the current security context.

## Security Features

### Template Security Manager

The `TemplateSecurityManager` provides comprehensive security features:

```go
type TemplateSecurityManager struct {
    sandbox         TemplateSandbox
    accessControl   AccessController
    patternFilter   PatternFilter
    resourceMonitor ResourceMonitor
    auditLogger     AuditLogger
    mu              sync.RWMutex
}
```

#### Security Levels

- **Restricted**: Minimal function access, strict limits
- **Standard**: Normal function access, standard limits
- **Elevated**: Extended function access, relaxed limits
- **Trusted**: Full function access, minimal restrictions

#### Pattern Filtering

The system includes built-in pattern filtering for dangerous content:

```go
forbiddenPatterns: []string{
    "exec", "system", "eval", "shell",
    "password", "secret", "key", "token",
}
```

#### Sandboxing

Template sandboxing provides execution isolation:

```go
type TemplateSandbox struct {
    enabled          bool
    maxExecutionTime time.Duration
    maxMemoryUsage   int64
    allowedFunctions map[string]bool
    mu               sync.RWMutex
}
```

## Performance Features

### Template Performance Manager

The `TemplatePerformanceManager` provides performance optimization:

```go
type TemplatePerformanceManager struct {
    cache        TemplateResultCache
    compiler     TemplateCompiler
    parallelizer ParallelProcessor
    monitor      PerformanceMonitor
    mu           sync.RWMutex
}
```

#### Caching

- Template caching with TTL
- Result caching for rendered templates
- Function result caching
- Context resolution caching

#### Optimization

- Template compilation optimization
- Parallel processing support
- Memory usage optimization
- Performance metrics collection

## Integration Examples

### Basic Template Usage

```go
// Create template manager
logger := spookylogging.GetLogger()
manager := spookytemplates.NewManager(logger)

// Load template
template, err := manager.LoadTemplate(ctx, "templates/nginx.conf.tmpl")
if err != nil {
    return fmt.Errorf("failed to load template: %w", err)
}

// Render template
result, err := manager.RenderTemplate(ctx, template, map[string]interface{}{
    "server_name": "example.com",
    "port": 80,
    "root_path": "/var/www/html",
})
if err != nil {
    return fmt.Errorf("failed to render template: %w", err)
}

fmt.Println(result)
```

### Template with Facts Integration

```go
// Set facts integration
manager.SetFactsIntegration(factsIntegration)

// Load template
template, err := manager.LoadTemplate(ctx, "templates/system-info.sh.tmpl")
if err != nil {
    return fmt.Errorf("failed to load template: %w", err)
}

// Render with facts context
result, err := manager.RenderTemplate(ctx, template, map[string]interface{}{
    "machine": "web-server",
})
if err != nil {
    return fmt.Errorf("failed to render template: %w", err)
}
```

### Template with Variables Integration

```go
// Set variables integration
manager.SetVariablesIntegration(variablesIntegration)

// Load template
template, err := manager.LoadTemplate(ctx, "templates/config.yaml.tmpl")
if err != nil {
    return fmt.Errorf("failed to load template: %w", err)
}

// Render with variables context
result, err := manager.RenderTemplate(ctx, template, map[string]interface{}{
    "environment": "production",
})
if err != nil {
    return fmt.Errorf("failed to render template: %w", err)
}
```

### Template Validation

```go
// Load template
template, err := manager.LoadTemplate(ctx, "templates/nginx.conf.tmpl")
if err != nil {
    return fmt.Errorf("failed to load template: %w", err)
}

// Validate template
result, err := manager.ValidateTemplate(ctx, template)
if err != nil {
    return fmt.Errorf("validation failed: %w", err)
}

if !result.Valid {
    fmt.Printf("Template validation failed:\n")
    for _, err := range result.Errors {
        fmt.Printf("  - %s\n", err.Message)
    }
    return fmt.Errorf("template validation failed")
}

fmt.Println("Template validation passed")
```

### Custom Template Functions

```go
// Define custom functions
customFunctions := map[string]interface{}{
    "formatSize": func(bytes int64) string {
        const unit = 1024
        if bytes < unit {
            return fmt.Sprintf("%d B", bytes)
        }
        div, exp := int64(unit), 0
        for n := bytes / unit; n >= unit; n /= unit {
            div *= unit
            exp++
        }
        return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
    },
    "isProduction": func(env string) bool {
        return env == "production"
    },
}

// Register custom functions
err := manager.RegisterTemplateFunctions(customFunctions)
if err != nil {
    return fmt.Errorf("failed to register functions: %w", err)
}

// Use custom functions in template
template := `File size: {{formatSize .file_size}}
Environment: {{.environment}}
Is Production: {{isProduction .environment}}`

result, err := manager.RenderTemplate(ctx, &spookytypes.Template{
    Content: template,
}, map[string]interface{}{
    "file_size": 1048576,
    "environment": "production",
})
```

## CLI Integration

### Template Commands

The templates system provides comprehensive CLI integration:

```bash
# Render template
spooky templates render <project> <template> --data <file> --output <file>

# Validate template
spooky templates validate <project> [--template <path>]

# List templates
spooky templates list <project> [--format json|hcl]

# Search templates
spooky templates search <project> <query> [--tags <tags>] [--category <category>]
```

### Example CLI Usage

```bash
# Render nginx configuration
spooky templates render ./myproject templates/nginx.conf.tmpl \
  --output /etc/nginx/nginx.conf

# Validate all templates
spooky templates validate ./myproject

# List templates in JSON format
spooky templates list ./myproject --format json

# Search for web templates
spooky templates search ./myproject "web" --tags nginx,apache
```

## Error Handling

### Template Errors

The templates system provides comprehensive error handling:

```go
// Template loading errors
if err != nil {
    return fmt.Errorf("failed to load template %s: %w", templatePath, err)
}

// Template rendering errors
if err != nil {
    return fmt.Errorf("failed to render template %s: %w", template.ID, err)
}

// Template validation errors
if !result.Valid {
    return fmt.Errorf("template validation failed: %v", result.Errors)
}

// Security violation errors
if err != nil {
    return fmt.Errorf("security violation: %w", err)
}
```

### Validation Results

```go
type ValidationResult struct {
    Valid    bool             `json:"valid"`
    Errors   []SchemaError    `json:"errors,omitempty"`
    Warnings []SchemaError    `json:"warnings,omitempty"`
    Info     []SchemaError    `json:"info,omitempty"`
    Details  map[string]interface{} `json:"details,omitempty"`
}
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

## Best Practices

### Security

- Use appropriate security levels for templates
- Validate all template inputs
- Avoid dangerous patterns in templates
- Use restricted mode for untrusted templates
- Monitor template execution and resource usage

### Performance

- Cache frequently used templates
- Use appropriate TTL values
- Monitor template performance
- Optimize template complexity
- Use parallel processing for multiple templates

### Maintainability

- Use descriptive template names
- Document template variables
- Version template metadata
- Use consistent naming conventions
- Organize templates by purpose and scope

### Integration

- Leverage facts for dynamic content
- Use variables for configuration
- Integrate with machine inventory
- Use secrets for sensitive data
- Validate template outputs

## Related Documentation

- [Templates System](TEMPLATES_SYSTEM.md) - System overview and architecture
- [Template System Design](design/systems/template-system.md) - Design documentation
- [Template Enhanced Composition](design/TEMPLATE_ENHANCED_COMPOSITION.md) - Composition patterns
- [Schema System](../schema-system.md) - Schema validation and configuration
- [Facts System](FACTS_SYSTEM.md) - Facts integration
- [Variables System](VARIABLES_SYSTEM.md) - Variables integration
- [Machines System](MACHINES_SYSTEM.md) - Machines integration
