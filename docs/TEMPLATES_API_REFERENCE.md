# Templates System API Reference

## Overview

This document provides a comprehensive API reference for the spooky templates system. It covers all interfaces, types, methods, and implementation details for developers working with the templates system.

**Status: Partially Implemented** - The templates system has basic functionality but SSH-based template rendering has known issues that need to be addressed.

## Core Interfaces

### TemplatesIntegration Interface

The `TemplatesIntegration` interface provides the primary entry point for templates operations:

```go
type TemplatesIntegration interface {
    // LoadTemplates loads templates from the specified project path
    LoadTemplates(ctx context.Context, projectPath string) ([]interface{}, error)
    
    // ValidateTemplates validates template definitions
    ValidateTemplates(ctx context.Context, templates []interface{}) (*ValidationResult, error)
    
    // RenderTemplate renders a template with data
    RenderTemplate(ctx context.Context, templatePath string, data map[string]interface{}) (string, error)
    
    // RenderTemplates renders multiple templates
    RenderTemplates(ctx context.Context, templates []interface{}, data map[string]interface{}) ([]RenderResult, error)
    
    // ListTemplates lists available templates
    ListTemplates(ctx context.Context, projectPath string) ([]TemplateInfo, error)
}
```

**Implementation Status**: ⚠️ **Partially Implemented** - Basic functionality exists but SSH-based rendering has issues

### TemplatesManager Interface

The `TemplatesManager` interface provides templates management and rendering:

```go
type TemplatesManager interface {
    // LoadTemplates loads templates from project configuration
    LoadTemplates(ctx context.Context, projectPath string) ([]*spookytypestemplates.Template, error)
    
    // ValidateTemplates validates template definitions
    ValidateTemplates(ctx context.Context, templates []*spookytypestemplates.Template) (*ValidationResult, error)
    
    // RenderTemplate renders a template with data
    RenderTemplate(ctx context.Context, templatePath string, data map[string]interface{}) (string, error)
    
    // RenderTemplates renders multiple templates
    RenderTemplates(ctx context.Context, templates []*spookytypestemplates.Template, data map[string]interface{}) ([]*spookytypestemplates.RenderResult, error)
    
    // ListTemplates lists available templates
    ListTemplates(ctx context.Context, projectPath string) ([]*spookytypestemplates.TemplateInfo, error)
}
```

**Implementation Status**: ⚠️ **Partially Implemented** - Basic loading and validation exist but rendering has issues

## Current Implementation Status

### ✅ Working Components

1. **Template Loading**: Loading templates from HCL configuration files
2. **Template Validation**: Basic validation of template definitions
3. **Template Structure**: Proper template type definitions and structures
4. **CLI Integration**: `spooky templates` commands with basic functionality
5. **Project Integration**: Templates loading from project configuration
6. **Basic Validation**: Template definition validation and error handling
7. **List Support**: Template listing and information display
8. **Local Rendering**: Basic local template rendering capabilities

### ⚠️ Known Issues

1. **SSH-Based Rendering**: SSH-based template rendering has implementation issues
2. **Remote Template Reading**: Cannot properly read templates from remote machines
3. **Parallel Processing**: No parallel template rendering support
4. **Advanced Functions**: Limited template function support
5. **Context Integration**: No template context integration
6. **Variable Integration**: No variable integration in templates

### 🔄 In Progress

1. **SSH Rendering Fixes**: Addressing SSH-based template rendering issues
2. **Function Enhancements**: Improving template function support
3. **Context Integration**: Adding template context integration

## Implementation Details

### Template Loading System

The templates system loads templates from HCL configuration files:

```go
type TemplateLoader struct {
    logger spookylogging.Logger
}

func (l *TemplateLoader) LoadTemplates(ctx context.Context, projectPath string) ([]*spookytypestemplates.Template, error) {
    var templates []*spookytypestemplates.Template
    
    // Load templates.hcl file
    templatesPath := filepath.Join(projectPath, "templates.hcl")
    if data, err := os.ReadFile(templatesPath); err == nil {
        if err := l.parseTemplatesFile(data, &templates); err != nil {
            return nil, fmt.Errorf("failed to parse templates.hcl: %w", err)
        }
    }
    
    // Load templates from templates/ directory
    templatesDir := filepath.Join(projectPath, "templates")
    if entries, err := os.ReadDir(templatesDir); err == nil {
        for _, entry := range entries {
            if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".hcl") {
                filePath := filepath.Join(templatesDir, entry.Name())
                if data, err := os.ReadFile(filePath); err == nil {
                    if err := l.parseTemplatesFile(data, &templates); err != nil {
                        return nil, fmt.Errorf("failed to parse %s: %w", entry.Name(), err)
                    }
                }
            }
        }
    }
    
    return templates, nil
}

func (l *TemplateLoader) parseTemplatesFile(data []byte, templates *[]*spookytypestemplates.Template) error {
    var config struct {
        Templates []*spookytypestemplates.Template `hcl:"template,block"`
    }
    
    if err := hcl.Unmarshal(data, &config); err != nil {
        return fmt.Errorf("failed to parse HCL: %w", err)
    }
    
    for _, template := range config.Templates {
        if template.Name == "" {
            return fmt.Errorf("template name is required")
        }
        
        // Check for duplicate names
        for _, existing := range *templates {
            if existing.Name == template.Name {
                return fmt.Errorf("duplicate template name: %s", template.Name)
            }
        }
        
        *templates = append(*templates, template)
    }
    
    return nil
}
```

**Supported Template Sources:**
- **Local Templates**: Templates defined in `templates.hcl` and `templates/*.hcl` files
- **File Templates**: Templates stored as files in the templates directory
- **SSH Templates**: Templates collected from remote machines (has issues)
- **Inline Templates**: Templates defined inline in configuration

### Template Validation System

Templates are validated against schemas and business rules:

```go
type TemplateValidator struct {
    logger spookylogging.Logger
}

func (v *TemplateValidator) ValidateTemplates(ctx context.Context, templates []*spookytypestemplates.Template) (*spookytypes.ValidationResult, error) {
    var errors []spookyschemas.SchemaError
    var warnings []spookyschemas.SchemaError
    
    for _, template := range templates {
        // Validate template name
        if template.Name == "" {
            errors = append(errors, spookyschemas.SchemaError{
                Message: "template name cannot be empty",
            })
            continue
        }
        
        // Validate template structure
        if err := v.validateTemplateStructure(template); err != nil {
            errors = append(errors, spookyschemas.SchemaError{
                Message: fmt.Sprintf("template %s: %s", template.Name, err.Error()),
            })
        }
        
        // Validate template source
        if err := v.validateTemplateSource(template); err != nil {
            errors = append(errors, spookyschemas.SchemaError{
                Message: fmt.Sprintf("template %s: %s", template.Name, err.Error()),
            })
        }
        
        // Validate template functions
        if err := v.validateTemplateFunctions(template); err != nil {
            warnings = append(warnings, spookyschemas.SchemaError{
                Message: fmt.Sprintf("template %s: %s", template.Name, err.Error()),
            })
        }
    }
    
    return &spookytypes.ValidationResult{
        Valid:    len(errors) == 0,
        Errors:   errors,
        Warnings: warnings,
    }, nil
}
```

### Template Rendering System

Templates are rendered using Go's template engine (currently has issues with SSH):

```go
type TemplateRenderer struct {
    logger spookylogging.Logger
}

func (r *TemplateRenderer) RenderTemplate(ctx context.Context, templatePath string, data map[string]interface{}) (string, error) {
    // Load template content
    content, err := r.loadTemplateContent(templatePath)
    if err != nil {
        return "", fmt.Errorf("failed to load template content: %w", err)
    }
    
    // Create template
    tmpl, err := template.New(templatePath).Parse(content)
    if err != nil {
        return "", fmt.Errorf("failed to parse template: %w", err)
    }
    
    // Add template functions
    tmpl.Funcs(r.getTemplateFunctions())
    
    // Render template
    var buf bytes.Buffer
    if err := tmpl.Execute(&buf, data); err != nil {
        return "", fmt.Errorf("failed to render template: %w", err)
    }
    
    return buf.String(), nil
}

func (r *TemplateRenderer) loadTemplateContent(templatePath string) (string, error) {
    // Check if it's a local file
    if strings.HasPrefix(templatePath, "/") || strings.HasPrefix(templatePath, "./") {
        return r.loadLocalTemplate(templatePath)
    }
    
    // Check if it's an SSH path
    if strings.Contains(templatePath, ":") {
        return r.loadSSHTemplate(templatePath)
    }
    
    // Default to local file
    return r.loadLocalTemplate(templatePath)
}

func (r *TemplateRenderer) loadLocalTemplate(templatePath string) (string, error) {
    data, err := os.ReadFile(templatePath)
    if err != nil {
        return "", fmt.Errorf("failed to read template file: %w", err)
    }
    return string(data), nil
}

func (r *TemplateRenderer) loadSSHTemplate(templatePath string) (string, error) {
    // Parse SSH path (host:path)
    parts := strings.SplitN(templatePath, ":", 2)
    if len(parts) != 2 {
        return "", fmt.Errorf("invalid SSH template path: %s", templatePath)
    }
    
    host := parts[0]
    path := parts[1]
    
    // TODO: Implement SSH template loading
    // This is where the SSH implementation would go
    return "", fmt.Errorf("SSH template loading not implemented yet")
}

func (r *TemplateRenderer) getTemplateFunctions() template.FuncMap {
    return template.FuncMap{
        "upper": strings.ToUpper,
        "lower": strings.ToLower,
        "title": strings.Title,
        "trim":  strings.TrimSpace,
        "split": strings.Split,
        "join":  strings.Join,
        "replace": strings.Replace,
        "contains": strings.Contains,
        "hasPrefix": strings.HasPrefix,
        "hasSuffix": strings.HasSuffix,
        "len": func(v interface{}) int {
            switch v := v.(type) {
            case string:
                return len(v)
            case []interface{}:
                return len(v)
            case map[string]interface{}:
                return len(v)
            default:
                return 0
            }
        },
        "default": func(defaultVal interface{}, val interface{}) interface{} {
            if val == nil || val == "" {
                return defaultVal
            }
            return val
        },
    }
}
```

## Type Definitions

### Template Types

```go
// Template represents a template definition
type Template struct {
    // Template name (required)
    Name string `json:"name" hcl:"name"`
    
    // Template description (optional)
    Description string `json:"description,omitempty" hcl:"description,optional"`
    
    // Template source (file path, SSH path, or inline content)
    Source string `json:"source" hcl:"source"`
    
    // Template type (file, ssh, inline)
    Type string `json:"type" hcl:"type"`
    
    // Template content (for inline templates)
    Content string `json:"content,omitempty" hcl:"content,optional"`
    
    // Template functions
    Functions map[string]interface{} `json:"functions,omitempty" hcl:"functions,optional"`
    
    // Template variables
    Variables map[string]interface{} `json:"variables,omitempty" hcl:"variables,optional"`
    
    // Template metadata
    Metadata map[string]interface{} `json:"metadata,omitempty" hcl:"metadata,optional"`
}

// TemplateInfo represents template information
type TemplateInfo struct {
    // Template name
    Name string `json:"name" hcl:"name"`
    
    // Template description
    Description string `json:"description,omitempty" hcl:"description,optional"`
    
    // Template source
    Source string `json:"source" hcl:"source"`
    
    // Template type
    Type string `json:"type" hcl:"type"`
    
    // Template size (for file templates)
    Size int64 `json:"size,omitempty" hcl:"size,optional"`
    
    // Template modification time
    ModTime time.Time `json:"mod_time,omitempty" hcl:"mod_time,optional"`
    
    // Template metadata
    Metadata map[string]interface{} `json:"metadata,omitempty" hcl:"metadata,optional"`
}

// RenderResult represents the result of template rendering
type RenderResult struct {
    // Template name
    Template string `json:"template" hcl:"template"`
    
    // Template source
    Source string `json:"source" hcl:"source"`
    
    // Rendered content
    Content string `json:"content" hcl:"content"`
    
    // Render success
    Success bool `json:"success" hcl:"success"`
    
    // Render error
    Error string `json:"error,omitempty" hcl:"error,optional"`
    
    // Render duration
    Duration time.Duration `json:"duration" hcl:"duration"`
    
    // Render timestamp
    Timestamp time.Time `json:"timestamp" hcl:"timestamp"`
}
```

### Template Context Types

```go
// TemplateContext provides context for template operations
type TemplateContext struct {
    // Project path
    ProjectPath string `json:"project_path" hcl:"project_path"`
    
    // Template being processed
    Template *Template `json:"template" hcl:"template"`
    
    // Template data
    Data map[string]interface{} `json:"data" hcl:"data"`
    
    // Operation timestamp
    Timestamp time.Time `json:"timestamp" hcl:"timestamp"`
    
    // Operation metadata
    Metadata map[string]interface{} `json:"metadata,omitempty" hcl:"metadata,optional"`
}

// TemplateResult represents the result of template operations
type TemplateResult struct {
    // Template context
    Context *TemplateContext `json:"context" hcl:"context"`
    
    // Operation success
    Success bool `json:"success" hcl:"success"`
    
    // Rendered content
    Content string `json:"content,omitempty" hcl:"content,optional"`
    
    // Operation error
    Error string `json:"error,omitempty" hcl:"error,optional"`
    
    // Operation duration
    Duration time.Duration `json:"duration" hcl:"duration"`
}
```

## Error Handling

### Template Errors

```go
// TemplateError represents template operation errors
type TemplateError struct {
    TemplateName string `json:"template_name" hcl:"template_name"`
    Error        string `json:"error" hcl:"error"`
    Details      string `json:"details,omitempty" hcl:"details,optional"`
}

// TemplateValidationError represents template validation errors
type TemplateValidationError struct {
    Field   string `json:"field" hcl:"field"`
    Message string `json:"message" hcl:"message"`
    Value   string `json:"value,omitempty" hcl:"value,optional"`
}
```

### Validation Implementation

```go
// ValidateTemplate validates a single template
func (v *TemplateValidator) ValidateTemplate(template *spookytypestemplates.Template) error {
    if template == nil {
        return fmt.Errorf("template cannot be nil")
    }
    
    // Validate required fields
    if template.Name == "" {
        return fmt.Errorf("template name is required")
    }
    
    if template.Source == "" && template.Content == "" {
        return fmt.Errorf("template source or content is required")
    }
    
    // Validate template type
    validTypes := []string{"file", "ssh", "inline"}
    valid := false
    for _, t := range validTypes {
        if template.Type == t {
            valid = true
            break
        }
    }
    if !valid {
        return fmt.Errorf("invalid template type: %s (valid types: %v)", template.Type, validTypes)
    }
    
    // Validate template source based on type
    switch template.Type {
    case "file":
        if template.Source == "" {
            return fmt.Errorf("source is required for file templates")
        }
    case "ssh":
        if template.Source == "" {
            return fmt.Errorf("source is required for SSH templates")
        }
        if !strings.Contains(template.Source, ":") {
            return fmt.Errorf("SSH source must be in format 'host:path'")
        }
    case "inline":
        if template.Content == "" {
            return fmt.Errorf("content is required for inline templates")
        }
    }
    
    return nil
}
```

## CLI Commands

### Templates List Command

```bash
# List all templates in a project
spooky templates list ./my-project

# List templates with specific types
spooky templates list ./my-project --type file,inline

# List templates with specific names
spooky templates list ./my-project --names web-config,db-config

# List templates with verbose output
spooky templates list ./my-project --verbose
```

### Templates Render Command

```bash
# Render a template with data
spooky templates render ./my-project templates/web-config.tmpl --data data.json

# Render a template with inline data
spooky templates render ./my-project templates/web-config.tmpl --data '{"port": 8080, "host": "localhost"}'

# Render a template with output to file
spooky templates render ./my-project templates/web-config.tmpl --data data.json --output web-config.conf

# Render a template with preview (dry-run)
spooky templates render ./my-project templates/web-config.tmpl --data data.json --preview
```

### Templates Validation Command

```bash
# Validate templates in a project
spooky templates validate ./my-project

# Validate templates with verbose output
spooky templates validate ./my-project --verbose

# Validate specific templates
spooky templates validate ./my-project --template templates/web-config.tmpl
```

## Integration Examples

### Basic Template Definition

```hcl
# templates.hcl
templates {
  template "web-config" {
    description = "Web server configuration template"
    type = "file"
    source = "templates/web-config.tmpl"
    
    variables = {
      default_port = 8080
      default_host = "localhost"
    }
    
    metadata = {
      category = "configuration"
      target = "web-server"
    }
  }
  
  template "db-config" {
    description = "Database configuration template"
    type = "file"
    source = "templates/db-config.tmpl"
    
    variables = {
      default_port = 5432
      default_host = "localhost"
    }
    
    metadata = {
      category = "configuration"
      target = "database"
    }
  }
  
  template "inline-script" {
    description = "Inline script template"
    type = "inline"
    content = "#!/bin/bash\necho 'Hello {{.name}}!'\necho 'Port: {{.port}}'"
    
    variables = {
      default_name = "World"
      default_port = 8080
    }
    
    metadata = {
      category = "script"
      target = "general"
    }
  }
}
```

### Template Loading and Validation

```go
// Template loading and validation example
func loadAndValidateTemplates(projectPath string) error {
    ctx := context.Background()
    
    // Create template manager
    manager := spookytemplates.NewManager(loader, validator, logger)
    
    // Load templates
    templates, err := manager.LoadTemplates(ctx, projectPath)
    if err != nil {
        return fmt.Errorf("failed to load templates: %w", err)
    }
    
    // Validate templates
    result, err := manager.ValidateTemplates(ctx, templates)
    if err != nil {
        return fmt.Errorf("failed to validate templates: %w", err)
    }
    
    if !result.Valid {
        fmt.Println("Template validation failed:")
        for _, error := range result.Errors {
            fmt.Printf("  - %s\n", error.Message)
        }
        return fmt.Errorf("template validation failed")
    }
    
    fmt.Printf("Loaded and validated %d templates\n", len(templates))
    return nil
}
```

### Template Rendering

```go
// Template rendering example
func renderTemplate(projectPath string, templateName string, data map[string]interface{}) error {
    ctx := context.Background()
    
    // Create template manager
    manager := spookytemplates.NewManager(loader, validator, logger)
    
    // Load templates
    templates, err := manager.LoadTemplates(ctx, projectPath)
    if err != nil {
        return fmt.Errorf("failed to load templates: %w", err)
    }
    
    // Find template by name
    var targetTemplate *spookytypestemplates.Template
    for _, template := range templates {
        if template.Name == templateName {
            targetTemplate = template
            break
        }
    }
    
    if targetTemplate == nil {
        return fmt.Errorf("template not found: %s", templateName)
    }
    
    // Render template
    content, err := manager.RenderTemplate(ctx, targetTemplate.Source, data)
    if err != nil {
        return fmt.Errorf("failed to render template: %w", err)
    }
    
    fmt.Printf("Rendered template %s:\n%s\n", templateName, content)
    return nil
}
```

## Current Limitations

### Rendering Limitations

1. **SSH Integration Issues**: SSH-based template rendering has known problems
2. **No Parallel Rendering**: Templates are rendered sequentially, not in parallel
3. **No Result Caching**: Rendered templates are not cached between operations
4. **No Incremental Rendering**: Always renders complete templates
5. **Limited Function Support**: Only basic template functions are supported

### Integration Limitations

1. **No Variable Integration**: Templates are not integrated with variable system
2. **No Context Integration**: Templates are not integrated with context system
3. **No Action Integration**: Templates are not used in action system
4. **No Facts Integration**: Templates are not integrated with facts system
5. **No Conditional Rendering**: No conditional template rendering

### Function Limitations

1. **Basic Functions Only**: Only basic string and utility functions
2. **No Custom Functions**: No support for custom template functions
3. **No SSH Functions**: No SSH-specific template functions
4. **No Data Functions**: No data manipulation functions
5. **No Security Functions**: No security-related template functions

## Future Enhancements

### Planned Features

1. **SSH Rendering Fixes**: Resolve SSH-based template rendering issues
2. **Function Enhancements**: Add more template functions
3. **Parallel Rendering**: Implement parallel template rendering
4. **Result Caching**: Add template result caching
5. **Incremental Rendering**: Support incremental template rendering
6. **Advanced Functions**: Support advanced template functions

### Integration Enhancements

1. **Variable Integration**: Integrate templates with variable system
2. **Context Integration**: Add template context integration
3. **Action Integration**: Use templates in action system
4. **Facts Integration**: Integrate templates with facts system
5. **Advanced Rendering**: Support advanced rendering features

## Summary

The templates system provides basic template loading and validation capabilities but has significant limitations with SSH-based rendering and function support that need to be addressed. The system is functional for basic use cases but requires improvements for production use.

**Status**: ⚠️ **Partially Implemented** - Basic functionality exists but SSH-based rendering and function support have issues that need to be resolved.
