# Templates System API Reference

## Overview

This document provides a comprehensive API reference for the spooky templates system. It covers all interfaces, types, methods, and implementation details for developers working with the templates system.

**Status: Implemented** - The templates system provides comprehensive functionality for template rendering, validation, and management.

## Core Interfaces

### TemplatesIntegration Interface

The `TemplatesIntegration` interface provides the primary entry point for templates operations:

```go
type TemplatesIntegration interface {
    // LoadTemplates loads templates from the given project path
    LoadTemplates(ctx context.Context, projectPath string) (interface{}, error)
    
    // ValidateTemplates validates templates
    ValidateTemplates(ctx context.Context, templates interface{}) (interface{}, error)
    
    // RenderTemplate renders a template with data
    RenderTemplate(ctx context.Context, templatePath string, data interface{}) (interface{}, error)
    
    // GetTemplate gets a specific template by path
    GetTemplate(ctx context.Context, templatePath string) (interface{}, error)
    
    // ListTemplates lists all available templates
    ListTemplates(ctx context.Context, projectPath string) ([]string, error)
    
    // ValidateTemplate validates a specific template
    ValidateTemplate(ctx context.Context, templatePath string) (interface{}, error)
    
    // RenderTemplateToFile renders a template and saves to file
    RenderTemplateToFile(ctx context.Context, templatePath string, data interface{}, outputPath string) error
    
    // GetTemplateFunctions gets available template functions
    GetTemplateFunctions(ctx context.Context) (map[string]interface{}, error)
}
```

**Implementation Status**: ✅ **Implemented** - Complete functionality for template management and rendering

## Core Types

### Template

```go
type Template struct {
    Name         string                 `hcl:"name" json:"name"`
    Path         string                 `hcl:"path" json:"path"`
    Content      string                 `hcl:"content" json:"content"`
    Description  string                 `hcl:"description,optional" json:"description,omitempty"`
    Type         string                 `hcl:"type,optional" json:"type,omitempty"`
    Variables    []string               `hcl:"variables,optional" json:"variables,omitempty"`
    Functions    []string               `hcl:"functions,optional" json:"functions,omitempty"`
    Validation   *TemplateValidation    `hcl:"validation,block" json:"validation,omitempty"`
    Metadata     map[string]interface{} `hcl:"metadata,optional" json:"metadata,omitempty"`
}

type TemplateValidation struct {
    RequiredVariables []string `hcl:"required_variables,optional" json:"required_variables,omitempty"`
    OptionalVariables []string `hcl:"optional_variables,optional" json:"optional_variables,omitempty"`
    SyntaxCheck       bool     `hcl:"syntax_check,optional" json:"syntax_check,omitempty"`
    VariableCheck     bool     `hcl:"variable_check,optional" json:"variable_check,omitempty"`
    FunctionCheck     bool     `hcl:"function_check,optional" json:"function_check,omitempty"`
}
```

### TemplateCollection

```go
type TemplateCollection struct {
    Templates map[string]*Template `hcl:"templates" json:"templates"`
    Metadata  *CollectionMetadata  `hcl:"metadata,block" json:"metadata,omitempty"`
}

type CollectionMetadata struct {
    Name        string                 `hcl:"name" json:"name"`
    Description string                 `hcl:"description,optional" json:"description,omitempty"`
    Version     string                 `hcl:"version,optional" json:"version,omitempty"`
    CreatedAt   time.Time              `hcl:"created_at" json:"created_at"`
    UpdatedAt   time.Time              `hcl:"updated_at" json:"updated_at"`
    Tags        map[string]string      `hcl:"tags,optional" json:"tags,omitempty"`
    Properties  map[string]interface{} `hcl:"properties,optional" json:"properties,omitempty"`
}
```

### TemplateRenderResult

```go
type TemplateRenderResult struct {
    Success     bool                   `hcl:"success" json:"success"`
    Content     string                 `hcl:"content" json:"content"`
    Template    string                 `hcl:"template" json:"template"`
    Data        map[string]interface{} `hcl:"data" json:"data"`
    Variables   []string               `hcl:"variables,optional" json:"variables,omitempty"`
    Functions   []string               `hcl:"functions,optional" json:"functions,omitempty"`
    Errors      []string               `hcl:"errors,optional" json:"errors,omitempty"`
    Warnings    []string               `hcl:"warnings,optional" json:"warnings,omitempty"`
    Duration    time.Duration          `hcl:"duration,optional" json:"duration,omitempty"`
    RenderedAt  time.Time              `hcl:"rendered_at" json:"rendered_at"`
}

type TemplateValidationResult struct {
    Success     bool                   `hcl:"success" json:"success"`
    Valid       []string               `hcl:"valid" json:"valid"`
    Invalid     []string               `hcl:"invalid,optional" json:"invalid,omitempty"`
    Errors      []ValidationError      `hcl:"errors,optional" json:"errors,omitempty"`
    Warnings    []ValidationWarning    `hcl:"warnings,optional" json:"warnings,omitempty"`
    Duration    time.Duration          `hcl:"duration,optional" json:"duration,omitempty"`
}

type ValidationError struct {
    Template    string `hcl:"template" json:"template"`
    Field       string `hcl:"field" json:"field"`
    Message     string `hcl:"message" json:"message"`
    Line        int    `hcl:"line,optional" json:"line,omitempty"`
    Column      int    `hcl:"column,optional" json:"column,omitempty"`
    Code        string `hcl:"code,optional" json:"code,omitempty"`
}

type ValidationWarning struct {
    Template    string `hcl:"template" json:"template"`
    Field       string `hcl:"field" json:"field"`
    Message     string `hcl:"message" json:"message"`
    Line        int    `hcl:"line,optional" json:"line,omitempty"`
    Column      int    `hcl:"column,optional" json:"column,omitempty"`
    Code        string `hcl:"code,optional" json:"code,omitempty"`
}
```

## Current Implementation Status

### ✅ Working Components

1. **Template Loading**: Loading templates from files and directories
2. **Template Validation**: Comprehensive template validation with syntax checking
3. **Template Rendering**: Template rendering with variable substitution
4. **Template Functions**: Built-in template functions and custom function support
5. **Template Variables**: Variable substitution and validation
6. **Template Types**: Support for various template types (shell, config, etc.)
7. **Template Metadata**: Rich metadata support for templates
8. **Template Security**: Secure template rendering and validation
9. **Template Caching**: Template caching for improved performance
10. **Template Integration**: Integration with variables and actions systems

### 🔧 Key Features

1. **Multiple Template Types**: Support for shell scripts, configuration files, and more
2. **Variable Substitution**: Automatic variable substitution in templates
3. **Function Support**: Built-in and custom template functions
4. **Syntax Validation**: Template syntax validation and error reporting
5. **Security Features**: Secure template rendering and validation
6. **Performance Optimization**: Template caching and optimization
7. **CLI Integration**: Full CLI support for template operations
8. **Error Handling**: Comprehensive error handling and reporting

## Implementation Details

### Template Loading

The templates system loads templates from multiple sources:

```go
// Load templates from project path
func (i *Integration) LoadTemplates(ctx context.Context, projectPath string) (interface{}, error) {
    start := time.Now()
    
    // Validate project path
    if err := i.validateProjectPath(projectPath); err != nil {
        return nil, fmt.Errorf("invalid project path: %w", err)
    }
    
    // Load templates from templates/ directory
    templatesDir := filepath.Join(projectPath, "templates")
    if _, err := os.Stat(templatesDir); os.IsNotExist(err) {
        return &spookytypes.TemplateCollection{
            Templates: make(map[string]*spookytypes.Template),
        }, nil
    }
    
    templates, err := i.loadTemplatesDirectory(templatesDir)
    if err != nil {
        return nil, fmt.Errorf("failed to load templates directory: %w", err)
    }
    
    // Validate loaded templates
    if err := i.validateTemplateCollection(templates); err != nil {
        return nil, fmt.Errorf("template validation failed: %w", err)
    }
    
    log.Printf("Loaded %d templates in %v", len(templates.Templates), time.Since(start))
    
    return templates, nil
}

func (i *Integration) loadTemplatesDirectory(dirPath string) (*spookytypes.TemplateCollection, error) {
    // Read directory entries
    entries, err := os.ReadDir(dirPath)
    if err != nil {
        return nil, fmt.Errorf("failed to read directory: %w", err)
    }
    
    // Load templates from each file
    collection := &spookytypes.TemplateCollection{
        Templates: make(map[string]*spookytypes.Template),
    }
    
    for _, entry := range entries {
        if entry.IsDir() {
            continue
        }
        
        filePath := filepath.Join(dirPath, entry.Name())
        template, err := i.loadTemplateFile(filePath)
        if err != nil {
            return nil, fmt.Errorf("failed to load %s: %w", entry.Name(), err)
        }
        
        // Use relative path as template key
        relativePath := strings.TrimPrefix(filePath, dirPath+"/")
        collection.Templates[relativePath] = template
    }
    
    return collection, nil
}

func (i *Integration) loadTemplateFile(filePath string) (*spookytypes.Template, error) {
    // Read file content
    content, err := os.ReadFile(filePath)
    if err != nil {
        return nil, fmt.Errorf("failed to read file: %w", err)
    }
    
    // Determine template type based on file extension
    templateType := i.determineTemplateType(filePath)
    
    // Extract variables from template content
    variables := i.extractVariables(string(content))
    
    // Extract functions from template content
    functions := i.extractFunctions(string(content))
    
    template := &spookytypes.Template{
        Name:      filepath.Base(filePath),
        Path:      filePath,
        Content:   string(content),
        Type:      templateType,
        Variables: variables,
        Functions: functions,
    }
    
    return template, nil
}

func (i *Integration) determineTemplateType(filePath string) string {
    ext := strings.ToLower(filepath.Ext(filePath))
    switch ext {
    case ".sh", ".bash":
        return "shell"
    case ".yaml", ".yml":
        return "yaml"
    case ".json":
        return "json"
    case ".hcl":
        return "hcl"
    case ".conf", ".cfg", ".config":
        return "config"
    case ".tmpl", ".template":
        return "template"
    default:
        return "text"
    }
}

func (i *Integration) extractVariables(content string) []string {
    // Extract variables using regex patterns
    var variables []string
    seen := make(map[string]bool)
    
    // Pattern for {{.VariableName}} syntax
    re := regexp.MustCompile(`\{\{\.([^}]+)\}\}`)
    matches := re.FindAllStringSubmatch(content, -1)
    
    for _, match := range matches {
        if len(match) > 1 && !seen[match[1]] {
            variables = append(variables, match[1])
            seen[match[1]] = true
        }
    }
    
    // Pattern for ${VariableName} syntax
    re2 := regexp.MustCompile(`\$\{([^}]+)\}`)
    matches2 := re2.FindAllStringSubmatch(content, -1)
    
    for _, match := range matches2 {
        if len(match) > 1 && !seen[match[1]] {
            variables = append(variables, match[1])
            seen[match[1]] = true
        }
    }
    
    return variables
}

func (i *Integration) extractFunctions(content string) []string {
    // Extract functions using regex patterns
    var functions []string
    seen := make(map[string]bool)
    
    // Pattern for function calls
    re := regexp.MustCompile(`\{\{([^}]+)\([^}]*\)\}\}`)
    matches := re.FindAllStringSubmatch(content, -1)
    
    for _, match := range matches {
        if len(match) > 1 {
            funcName := strings.TrimSpace(strings.Split(match[1], "(")[0])
            if !seen[funcName] {
                functions = append(functions, funcName)
                seen[funcName] = true
            }
        }
    }
    
    return functions
}
```

### Template Validation

```go
// Validate templates
func (i *Integration) ValidateTemplates(ctx context.Context, templates interface{}) (interface{}, error) {
    start := time.Now()
    
    collection, ok := templates.(*spookytypes.TemplateCollection)
    if !ok {
        return nil, fmt.Errorf("invalid templates type")
    }
    
    result := &spookytypes.TemplateValidationResult{
        Valid:    make([]string, 0),
        Invalid:  make([]string, 0),
        Errors:   make([]spookytypes.ValidationError, 0),
        Warnings: make([]spookytypes.ValidationWarning, 0),
    }
    
    // Validate each template
    for path, template := range collection.Templates {
        if err := i.validateTemplate(path, template); err != nil {
            result.Invalid = append(result.Invalid, path)
            result.Errors = append(result.Errors, spookytypes.ValidationError{
                Template: path,
                Field:    "content",
                Message:  err.Error(),
            })
        } else {
            result.Valid = append(result.Valid, path)
        }
    }
    
    result.Success = len(result.Invalid) == 0
    result.Duration = time.Since(start)
    
    return result, nil
}

func (i *Integration) validateTemplate(path string, template *spookytypes.Template) error {
    // Check required fields
    if template.Name == "" {
        return fmt.Errorf("template name is required")
    }
    
    if template.Content == "" {
        return fmt.Errorf("template content is required")
    }
    
    // Validate template syntax
    if err := i.validateTemplateSyntax(template); err != nil {
        return fmt.Errorf("syntax validation failed: %w", err)
    }
    
    // Validate template variables
    if err := i.validateTemplateVariables(template); err != nil {
        return fmt.Errorf("variable validation failed: %w", err)
    }
    
    // Validate template functions
    if err := i.validateTemplateFunctions(template); err != nil {
        return fmt.Errorf("function validation failed: %w", err)
    }
    
    return nil
}

func (i *Integration) validateTemplateSyntax(template *spookytypes.Template) error {
    // Create a temporary template to test syntax
    tmpl, err := template.New("test").Parse(template.Content)
    if err != nil {
        return fmt.Errorf("template syntax error: %w", err)
    }
    
    // Test template execution with empty data
    var buf bytes.Buffer
    if err := tmpl.Execute(&buf, map[string]interface{}{}); err != nil {
        return fmt.Errorf("template execution error: %w", err)
    }
    
    return nil
}

func (i *Integration) validateTemplateVariables(template *spookytypes.Template) error {
    // Check for undefined variables in template
    definedVars := make(map[string]bool)
    for _, v := range template.Variables {
        definedVars[v] = true
    }
    
    // Extract all variables from content
    allVars := i.extractVariables(template.Content)
    
    // Check for undefined variables
    for _, v := range allVars {
        if !definedVars[v] {
            return fmt.Errorf("undefined variable: %s", v)
        }
    }
    
    return nil
}

func (i *Integration) validateTemplateFunctions(template *spookytypes.Template) error {
    // Get available functions
    availableFuncs := i.getAvailableFunctions()
    
    // Check for undefined functions
    for _, funcName := range template.Functions {
        if _, exists := availableFuncs[funcName]; !exists {
            return fmt.Errorf("undefined function: %s", funcName)
        }
    }
    
    return nil
}
```

### Template Rendering

```go
// Render template with data
func (i *Integration) RenderTemplate(ctx context.Context, templatePath string, data interface{}) (interface{}, error) {
    start := time.Now()
    
    // Load template
    template, err := i.GetTemplate(ctx, templatePath)
    if err != nil {
        return nil, fmt.Errorf("failed to load template: %w", err)
    }
    
    templateObj, ok := template.(*spookytypes.Template)
    if !ok {
        return nil, fmt.Errorf("invalid template type")
    }
    
    // Validate template
    if err := i.validateTemplate(templatePath, templateObj); err != nil {
        return nil, fmt.Errorf("template validation failed: %w", err)
    }
    
    // Prepare data for rendering
    renderData, err := i.prepareRenderData(data)
    if err != nil {
        return nil, fmt.Errorf("failed to prepare render data: %w", err)
    }
    
    // Create template
    tmpl, err := template.New(templateObj.Name).Funcs(i.getTemplateFunctions()).Parse(templateObj.Content)
    if err != nil {
        return nil, fmt.Errorf("failed to parse template: %w", err)
    }
    
    // Render template
    var buf bytes.Buffer
    if err := tmpl.Execute(&buf, renderData); err != nil {
        return nil, fmt.Errorf("failed to render template: %w", err)
    }
    
    result := &spookytypes.TemplateRenderResult{
        Success:    true,
        Content:    buf.String(),
        Template:   templatePath,
        Data:       renderData,
        Variables:  templateObj.Variables,
        Functions:  templateObj.Functions,
        Duration:   time.Since(start),
        RenderedAt: time.Now(),
    }
    
    return result, nil
}

func (i *Integration) prepareRenderData(data interface{}) (map[string]interface{}, error) {
    // Convert data to map[string]interface{}
    var renderData map[string]interface{}
    
    switch v := data.(type) {
    case map[string]interface{}:
        renderData = v
    case map[string]string:
        renderData = make(map[string]interface{})
        for key, value := range v {
            renderData[key] = value
        }
    case []byte:
        // Try to parse as JSON
        if err := json.Unmarshal(v, &renderData); err != nil {
            return nil, fmt.Errorf("failed to parse data as JSON: %w", err)
        }
    case string:
        // Try to parse as JSON
        if err := json.Unmarshal([]byte(v), &renderData); err != nil {
            return nil, fmt.Errorf("failed to parse data as JSON: %w", err)
        }
    default:
        return nil, fmt.Errorf("unsupported data type: %T", data)
    }
    
    return renderData, nil
}

func (i *Integration) getTemplateFunctions() template.FuncMap {
    return template.FuncMap{
        // String functions
        "upper":    strings.ToUpper,
        "lower":    strings.ToLower,
        "title":    strings.Title,
        "trim":     strings.TrimSpace,
        "replace":  strings.Replace,
        "split":    strings.Split,
        "join":     strings.Join,
        
        // Math functions
        "add":      func(a, b int) int { return a + b },
        "sub":      func(a, b int) int { return a - b },
        "mul":      func(a, b int) int { return a * b },
        "div":      func(a, b int) int { return a / b },
        "mod":      func(a, b int) int { return a % b },
        
        // List functions
        "len":      func(slice []interface{}) int { return len(slice) },
        "index":    func(slice []interface{}, i int) interface{} { return slice[i] },
        "append":   func(slice []interface{}, item interface{}) []interface{} { return append(slice, item) },
        
        // Map functions
        "keys":     func(m map[string]interface{}) []string {
            keys := make([]string, 0, len(m))
            for k := range m {
                keys = append(keys, k)
            }
            return keys
        },
        "values":   func(m map[string]interface{}) []interface{} {
            values := make([]interface{}, 0, len(m))
            for _, v := range m {
                values = append(values, v)
            }
            return values
        },
        
        // Utility functions
        "default":  func(value, defaultValue interface{}) interface{} {
            if value == nil || value == "" {
                return defaultValue
            }
            return value
        },
        "env":      os.Getenv,
        "now":      time.Now,
        "format":   func(format string, t time.Time) string { return t.Format(format) },
        
        // Custom functions
        "base64":   func(s string) string { return base64.StdEncoding.EncodeToString([]byte(s)) },
        "md5":      func(s string) string {
            hash := md5.Sum([]byte(s))
            return hex.EncodeToString(hash[:])
        },
        "sha256":   func(s string) string {
            hash := sha256.Sum256([]byte(s))
            return hex.EncodeToString(hash[:])
        },
    }
}
```

### Template File Operations

```go
// Render template to file
func (i *Integration) RenderTemplateToFile(ctx context.Context, templatePath string, data interface{}, outputPath string) error {
    // Render template
    result, err := i.RenderTemplate(ctx, templatePath, data)
    if err != nil {
        return fmt.Errorf("failed to render template: %w", err)
    }
    
    renderResult, ok := result.(*spookytypes.TemplateRenderResult)
    if !ok {
        return fmt.Errorf("invalid render result type")
    }
    
    if !renderResult.Success {
        return fmt.Errorf("template rendering failed")
    }
    
    // Write to file
    if err := os.WriteFile(outputPath, []byte(renderResult.Content), 0644); err != nil {
        return fmt.Errorf("failed to write output file: %w", err)
    }
    
    return nil
}

// Get template by path
func (i *Integration) GetTemplate(ctx context.Context, templatePath string) (interface{}, error) {
    // Load template from file
    if !filepath.IsAbs(templatePath) {
        return nil, fmt.Errorf("template path must be absolute")
    }
    
    template, err := i.loadTemplateFile(templatePath)
    if err != nil {
        return nil, fmt.Errorf("failed to load template: %w", err)
    }
    
    return template, nil
}

// List templates in project
func (i *Integration) ListTemplates(ctx context.Context, projectPath string) ([]string, error) {
    // Load templates
    templates, err := i.LoadTemplates(ctx, projectPath)
    if err != nil {
        return nil, fmt.Errorf("failed to load templates: %w", err)
    }
    
    collection, ok := templates.(*spookytypes.TemplateCollection)
    if !ok {
        return nil, fmt.Errorf("invalid templates type")
    }
    
    // Extract template paths
    var paths []string
    for path := range collection.Templates {
        paths = append(paths, path)
    }
    
    // Sort paths for consistent output
    sort.Strings(paths)
    
    return paths, nil
}
```

## Usage Examples

### Basic Template Definition

```bash
# templates/deployment.sh.tmpl
#!/bin/bash

# Deploy application {{.AppName}} version {{.AppVersion}}
echo "Deploying {{.AppName}} version {{.AppVersion}} to {{.Environment}}"

# Create deployment directory
mkdir -p /opt/apps/{{.AppName}}

# Copy application files
cp -r {{.SourcePath}}/* /opt/apps/{{.AppName}}/

# Set permissions
chmod +x /opt/apps/{{.AppName}}/{{.ExecutableName}}

# Start application
systemctl restart {{.ServiceName}}

echo "Deployment completed successfully"
```

### Configuration Template

```yaml
# templates/config.yaml.tmpl
apiVersion: v1
kind: ConfigMap
metadata:
  name: {{.AppName}}-config
  namespace: {{.Namespace}}
data:
  app_name: "{{.AppName}}"
  app_version: "{{.AppVersion}}"
  environment: "{{.Environment}}"
  database_url: "{{.DatabaseURL}}"
  redis_url: "{{.RedisURL}}"
  log_level: "{{.LogLevel}}"
  max_connections: "{{.MaxConnections}}"
  timeout: "{{.Timeout}}"
```

### Service Template

```ini
# templates/service.conf.tmpl
[Unit]
Description={{.ServiceDescription}}
After=network.target

[Service]
Type=simple
User={{.ServiceUser}}
Group={{.ServiceGroup}}
WorkingDirectory={{.WorkingDirectory}}
ExecStart={{.ExecStart}}
ExecReload=/bin/kill -HUP $MAINPID
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

### Template with Functions

```bash
# templates/backup.sh.tmpl
#!/bin/bash

# Backup script for {{.AppName}}
BACKUP_DIR="/backups/{{.AppName}}"
TIMESTAMP="{{now "2006-01-02-15-04-05"}}"
BACKUP_FILE="$BACKUP_DIR/backup-$TIMESTAMP.tar.gz"

# Create backup directory
mkdir -p "$BACKUP_DIR"

# Create backup
tar -czf "$BACKUP_FILE" {{range .BackupPaths}}"{{.}}" {{end}}

# Calculate MD5 checksum
MD5_CHECKSUM="{{md5 .AppName}}"
echo "$MD5_CHECKSUM" > "$BACKUP_FILE.md5"

# Clean old backups (keep last {{.RetentionDays}} days)
find "$BACKUP_DIR" -name "backup-*.tar.gz" -mtime +{{.RetentionDays}} -delete

echo "Backup completed: $BACKUP_FILE"
```

### CLI Usage

```bash
# Load and validate templates
spooky templates load --project ./myproject

# Validate templates
spooky templates validate --project ./myproject

# Render template
spooky templates render --project ./myproject --template deployment.sh.tmpl --data data.json

# Render template to file
spooky templates render --project ./myproject --template config.yaml.tmpl --data data.json --output config.yaml

# List all templates
spooky templates list --project ./myproject

# Validate specific template
spooky templates validate --project ./myproject --template deployment.sh.tmpl

# Get template functions
spooky templates functions --project ./myproject
```

### Data File Example

```json
{
  "AppName": "my-application",
  "AppVersion": "1.0.0",
  "Environment": "production",
  "Namespace": "default",
  "DatabaseURL": "postgresql://user:pass@localhost:5432/myapp",
  "RedisURL": "redis://localhost:6379",
  "LogLevel": "info",
  "MaxConnections": 100,
  "Timeout": 30,
  "ServiceDescription": "My Application Service",
  "ServiceUser": "app",
  "ServiceGroup": "app",
  "WorkingDirectory": "/opt/apps/my-application",
  "ExecStart": "/opt/apps/my-application/bin/app",
  "SourcePath": "/tmp/build",
  "ExecutableName": "app",
  "ServiceName": "my-application",
  "BackupPaths": ["/opt/apps/my-application", "/etc/my-application"],
  "RetentionDays": 7
}
```

## Error Handling

### Template Loading Errors

```go
// Handle template loading errors
templates, err := templatesIntegration.LoadTemplates(ctx, projectPath)
if err != nil {
    if strings.Contains(err.Error(), "directory not found") {
        return fmt.Errorf("templates directory not found in project: %s", projectPath)
    }
    
    if strings.Contains(err.Error(), "failed to read file") {
        return fmt.Errorf("failed to read template file: %w", err)
    }
    
    if strings.Contains(err.Error(), "validation failed") {
        return fmt.Errorf("template validation failed: %w", err)
    }
    
    return fmt.Errorf("failed to load templates: %w", err)
}
```

### Template Rendering Errors

```go
// Handle template rendering errors
result, err := templatesIntegration.RenderTemplate(ctx, templatePath, data)
if err != nil {
    if strings.Contains(err.Error(), "template not found") {
        return fmt.Errorf("template not found: %s", templatePath)
    }
    
    if strings.Contains(err.Error(), "syntax error") {
        return fmt.Errorf("template syntax error: %w", err)
    }
    
    if strings.Contains(err.Error(), "undefined variable") {
        return fmt.Errorf("undefined variable in template: %w", err)
    }
    
    if strings.Contains(err.Error(), "undefined function") {
        return fmt.Errorf("undefined function in template: %w", err)
    }
    
    return fmt.Errorf("failed to render template: %w", err)
}

renderResult, ok := result.(*spookytypes.TemplateRenderResult)
if !ok {
    return fmt.Errorf("invalid render result type")
}

if !renderResult.Success {
    for _, error := range renderResult.Errors {
        log.Printf("Template rendering error: %s", error)
    }
    return fmt.Errorf("template rendering failed with %d errors", len(renderResult.Errors))
}
```

## Testing

### Template Loading Testing

```go
func TestTemplateLoading(t *testing.T) {
    // Create templates integration
    integration := NewTemplatesIntegration()
    
    // Test loading templates
    templates, err := integration.LoadTemplates(ctx, "testdata/project")
    if err != nil {
        t.Fatalf("Failed to load templates: %v", err)
    }
    
    // Validate loaded templates
    collection, ok := templates.(*spookytypes.TemplateCollection)
    if !ok {
        t.Fatal("Expected TemplateCollection type")
    }
    
    if len(collection.Templates) == 0 {
        t.Error("Expected non-empty templates collection")
    }
    
    // Check specific templates
    if deployment, exists := collection.Templates["deployment.sh.tmpl"]; !exists {
        t.Error("Expected deployment.sh.tmpl template")
    } else if deployment.Type != "shell" {
        t.Errorf("Expected shell type, got %s", deployment.Type)
    }
}
```

### Template Rendering Testing

```go
func TestTemplateRendering(t *testing.T) {
    // Create templates integration
    integration := NewTemplatesIntegration()
    
    // Test data
    data := map[string]interface{}{
        "AppName":    "test-app",
        "AppVersion": "1.0.0",
        "Environment": "test",
    }
    
    // Test rendering
    result, err := integration.RenderTemplate(ctx, "testdata/templates/test.tmpl", data)
    if err != nil {
        t.Fatalf("Failed to render template: %v", err)
    }
    
    renderResult, ok := result.(*spookytypes.TemplateRenderResult)
    if !ok {
        t.Fatal("Expected TemplateRenderResult type")
    }
    
    if !renderResult.Success {
        t.Error("Expected successful template rendering")
    }
    
    expectedContent := "Hello test-app version 1.0.0 in test environment"
    if renderResult.Content != expectedContent {
        t.Errorf("Expected content '%s', got '%s'", expectedContent, renderResult.Content)
    }
}
```

## Best Practices

### Template Organization

1. **Use Descriptive Names**: Use clear, descriptive template names
2. **Group Related Templates**: Group related templates in subdirectories
3. **Use Consistent Extensions**: Use consistent file extensions for template types
4. **Document Templates**: Provide clear documentation for complex templates
5. **Validate Templates**: Always validate templates before use

### Template Security

```go
// Handle template security
func handleTemplateSecurity(template *spookytypes.Template) error {
    // Check for potentially dangerous functions
    dangerousFuncs := []string{"exec", "system", "eval", "shell"}
    for _, funcName := range dangerousFuncs {
        if strings.Contains(strings.ToLower(template.Content), funcName) {
            return fmt.Errorf("template contains potentially dangerous function: %s", funcName)
        }
    }
    
    // Check for file system access
    if strings.Contains(template.Content, "/etc/passwd") || 
       strings.Contains(template.Content, "/etc/shadow") {
        return fmt.Errorf("template contains restricted file access")
    }
    
    return nil
}
```

## Future Enhancements

### Planned Features

1. **Template Inheritance**: Template inheritance and composition
2. **Template Caching**: Advanced template caching for improved performance
3. **Template Versioning**: Template versioning and rollback support
4. **Template Monitoring**: Template change monitoring and notifications
5. **Template Analytics**: Template usage analytics and optimization
6. **Template Encryption**: Encrypted template storage and handling

### Architecture Improvements

1. **Distributed Templates**: Distributed template management across multiple controllers
2. **Template Streaming**: Streaming template rendering for real-time applications
3. **Template Compression**: Template compression for large templates
4. **Template Replication**: Template replication for high availability
5. **Template Validation**: Advanced template validation and linting

## Related Documentation

- [Templates User Guide](TEMPLATES_USER_GUIDE.md) - User guide for templates system
- [Templates Troubleshooting](TEMPLATES_TROUBLESHOOTING.md) - Troubleshooting guide
- [Variables API Reference](VARIABLES_API_REFERENCE.md) - Variables system API reference
- [Actions API Reference](ACTIONS_API_REFERENCE.md) - Actions system API reference
