# Template System: Comprehensive Implementation Plan

## Overview

This document is the authoritative source for all template system implementation details in spooky. It covers template rendering, context management, validation, preview capabilities, and integration with all other spooky systems.

**Schema Integration**: This template system implements the schema validation patterns and template configuration definitions defined in [Schema System](../schema-system.md) for comprehensive template validation, context schema enforcement, and schema-based template rendering. 

- **Template Metadata Schema**: See [`internal/schemas/schemas/template-metadata.hcl`](../../internal/schemas/schemas/template-metadata.hcl) for the authoritative schema describing template metadata, required variables, and output info.
- **Template Context Schema**: See [`internal/schemas/schemas/template-context.hcl`](../../internal/schemas/schemas/template-context.hcl) for the schema describing the data context available to templates at render time (facts, variables, machines, environment, custom data).
- **Template Functions & Restrictions Schema**: See [`internal/schemas/schemas/template-functions.hcl`](../../internal/schemas/schemas/template-functions.hcl) for the schema describing allowed template functions and restricted patterns for security and best practices.

**Architecture Integration**: Templates integrate with the overall spooky architecture as described in [Spooky Design](../spooky-design.md), providing dynamic content generation and configuration management for all system components.

## System Integration

This template system integrates with other core Spooky systems to provide comprehensive template rendering and dynamic content generation:

### **Actions System Integration**
- **Template Actions**: Actions can run template rendering and deployment (see [Actions System](../actions-system.md))
- **Template Context**: Actions provide template context with facts, variables, and machine data
- **Template Deployment**: Actions can deploy rendered templates to target machines
- **Template Validation**: Actions validate templates before running

### **Facts System Integration**
- **Facts in Templates**: Templates can access machine facts for dynamic content (see [Facts System](../facts-system.md))
- **System Facts**: Templates have access to system facts (OS, hardware, network)
- **Custom Facts**: Templates can use custom facts for project-specific data
- **Facts Functions**: Template functions for accessing facts data (`system()`, `custom()`)

### **Variables System Integration**
- **Variables in Templates**: Templates can access project variables for dynamic configuration (see [Variables System](../variables-system.md))
- **Variable Interpolation**: Templates support variable interpolation and resolution
- **Variable Functions**: Template functions for accessing variables (`var()`, `varOrDefault()`)
- **Variable Precedence**: Templates follow variable resolution precedence

### **Project System Integration**
- **Project Templates**: Templates are stored in project-specific `templates/` directories (see [Project System](../project-system.md))
- **Project Context**: Templates run within project context and configuration
- **Project Data**: Templates have access to project-specific data and configuration
- **Project Validation**: Template validation integrated with project validation

### **CLI System Integration**
- **Template Commands**: Template management through `spooky templates` commands (see [CLI System](../cli-system.md))
- **Command Patterns**: Template commands follow the established `spooky noun verb` CLI pattern
- **Validation Commands**: `spooky templates validate` for template validation
- **Rendering Commands**: `spooky templates render` for template rendering
- **Preview Commands**: `spooky templates render --preview` for template analysis

### **Configuration System Integration**
- **Global Configuration**: Template settings inherit from global configuration (see [Configuration System](../configuration-system.md))
- **Template Settings**: Template rendering settings, timeouts, and limits
- **Security Settings**: Template run security policies and restrictions
- **Performance Settings**: Template caching and optimization settings

### **Machines System Integration**
- **Machine Context**: Templates can access machine-specific data and facts (see [Machines System](../machines-system.md))
- **Machine Functions**: Template functions for accessing machine data
- **Remote Rendering**: Templates can be rendered on remote machines via SSH
- **Machine-Specific Templates**: Templates can be customized per machine

## Current State Analysis

### **What We Have**
- ✅ **Template rendering engine** with Go templates
- ✅ **Template context management** with facts, variables, and environment data
- ✅ **Template CLI commands** for render, validate, and preview
- ✅ **Template functions** for accessing data (`custom()`, `system()`, `env()`, `data()`)
- ✅ **Template validation** and syntax checking
- ✅ **Template preview** capabilities with mock data
- ✅ **Remote template running** via SSH

### **What We Need**
- 🔄 **Enhanced template schema** with comprehensive validation
- 🔄 **Template caching** and performance optimization
- 🔄 **Template inheritance** and composition
- 🔄 **Template security** and sandboxing
- 🔄 **Template debugging** and error reporting
- 🔄 **Template versioning** and change tracking

## Template System Design

### **1. Template Context and Data Access**

See [`template-context.hcl`](../../internal/schemas/schemas/template-context.hcl) for the full schema definition of the template context available at render time.

**Template Context Structure**:
```go
// TemplateContext holds all data available to templates
type TemplateContext struct {
    // Project configuration
    Project *config.ProjectConfig
    
    // Machine facts (from facts.db or JSON)
    Facts map[string]interface{}
    
    // Inventory information
    Machines []*config.Machine
    
    // Actions configuration
    Actions []*config.Action
    
    // Server-specific data (when --server is specified)
    ServerFacts map[string]interface{}
    
    // Environment variables
    Environment map[string]string
    
    // Custom data files
    CustomData map[string]interface{}
    
    // Project variables
    Variables *VariableContext
}
```

**Template Data Access**:
```go
// Template functions for data access
func (tc *TemplateContext) GetTemplateFunctions() map[string]interface{} {
    return map[string]interface{}{
        // Project functions
        "project":            func() *config.ProjectConfig { return tc.Project },
        "projectName":        func() string { return tc.Project.Name },
        "projectDescription": func() string { return tc.Project.Description },
        
        // Facts functions
        "facts": func() map[string]interface{} { return tc.Facts },
        "fact":  func(key string) interface{} { return tc.Facts[key] },
        
        // Machine functions
        "machines": func() []*config.Machine { return tc.Machines },
        "machine": func(name string) *config.Machine {
            for _, m := range tc.Machines {
                if m.Name == name {
                    return m
                }
            }
            return nil
        },
        
        // Environment functions
        "env": func(key string) string { return tc.Environment[key] },
        "envOrDefault": func(key, defaultValue string) string {
            if value, exists := tc.Environment[key]; exists {
                return value
            }
            return defaultValue
        },
        
        // Custom data functions
        "data":    func() map[string]interface{} { return tc.CustomData },
        "dataKey": func(key string) interface{} { return tc.CustomData[key] },
        
        // Variable functions
        "var": func(name string) interface{} { return tc.Variables.Get(name) },
        "varOrDefault": func(name string, defaultValue interface{}) interface{} {
            if value := tc.Variables.Get(name); value != nil {
                return value
            }
            return defaultValue
        },
    }
}
```

### **2. Template Functions and Data Access**

See [`template-functions.hcl`](../../internal/schemas/schemas/template-functions.hcl) for the schema defining allowed template functions and restricted patterns.

**System Facts Access**:
```go
// Template functions for accessing system facts
"system": func(path string) interface{} {
    return tc.getNestedValue(tc.Facts, path)
},

"custom": func(path string) interface{} {
    return tc.getNestedValue(tc.CustomData, path)
},

"env": func(key string) string {
    return tc.Environment[key]
},

"data": func(path string) interface{} {
    return tc.getNestedValue(tc.CustomData, path)
}
```

**Variable Access**:
```go
// Template functions for accessing variables
"var": func(name string) interface{} {
    if tc.Variables == nil {
        return nil
    }
    return tc.Variables.Get(name)
},

"varOrDefault": func(name string, defaultValue interface{}) interface{} {
    if tc.Variables == nil {
        return defaultValue
    }
    if value := tc.Variables.Get(name); value != nil {
        return value
    }
    return defaultValue
},

"varRequired": func(name string) interface{} {
    if tc.Variables == nil {
        panic(fmt.Sprintf("variable '%s' not found", name))
    }
    if value := tc.Variables.Get(name); value != nil {
        return value
    }
    panic(fmt.Sprintf("variable '%s' not found", name))
}
```

### **3. Template Rendering Engine**

**Core Rendering Engine**:
```go
// TemplateEngine handles template rendering with fact integration
type TemplateEngine struct {
    manager *facts.Manager
}

// RenderTemplate renders a template with fact data
func (te *TemplateEngine) RenderTemplate(templateFile string, server string, additionalData map[string]interface{}) (string, error) {
    // Read template file
    content, err := os.ReadFile(templateFile)
    if err != nil {
        return "", fmt.Errorf("error reading template file: %w", err)
    }
    
    // Preprocess template content
    contentStr := te.preprocessTemplateContent(string(content))
    
    // Create template data
    data := &TemplateData{
        System: make(map[string]interface{}),
        Custom: make(map[string]interface{}),
        Env:    make(map[string]string),
        Data:   additionalData,
    }
    
    // Load system facts if manager is available
    if te.manager != nil {
        if collection, err := te.manager.CollectAllFacts(server); err == nil {
            for key, fact := range collection.Facts {
                data.System[key] = fact.Value
            }
        }
        
        if customFacts, err := te.manager.GetCustomFacts(server); err == nil {
            data.Custom = customFacts
        }
    }
    
    // Load environment variables
    for _, env := range os.Environ() {
        pair := strings.SplitN(env, "=", 2)
        if len(pair) == 2 {
            data.Env[pair[0]] = pair[1]
        }
    }
    
    // Create template with custom functions
    tmpl, err := template.New(filepath.Base(templateFile)).Funcs(template.FuncMap{
        "custom": func(path string) interface{} {
            return te.getNestedValue(data.Custom, path)
        },
        "system": func(path string) interface{} {
            return te.getNestedValue(data.System, path)
        },
        "env": func(key string) string {
            return data.Env[key]
        },
        "data": func(path string) interface{} {
            return te.getNestedValue(data.Data, path)
        },
    }).Parse(contentStr)
    
    if err != nil {
        return "", fmt.Errorf("error parsing template: %w", err)
    }
    
    // Render template
    var buf bytes.Buffer
    if err := tmpl.Execute(&buf, data); err != nil {
        return "", fmt.Errorf("error rendering template: %w", err)
    }
    
    return buf.String(), nil
}
```

### **4. Template Examples and Use Cases**

**Basic Template Example**:
```go
// templates/nginx.conf.tmpl
server {
    listen {{var "server_port"}};
    server_name {{var "server_name"}};
    
    root /var/www/{{var "app_name"}};
    index index.html;
    
    location / {
        try_files $uri $uri/ =404;
    }
    
    # Logging
    access_log /var/log/nginx/{{var "app_name"}}_access.log;
    error_log /var/log/nginx/{{var "app_name"}}_error.log;
}
```

**Template with Facts**:
```go
// templates/systemd.service.tmpl
[Unit]
Description={{var "app_name"}} Service
After=network.target

[Service]
Type=simple
User={{var "service_user"}}
WorkingDirectory=/opt/{{var "app_name"}}
ExecStart=/opt/{{var "app_name"}}/{{var "app_name"}}
Restart=always
RestartSec=5

# Environment variables
Environment=NODE_ENV={{var "environment"}}
Environment=PORT={{var "server_port"}}
Environment=HOSTNAME={{system "hostname"}}
Environment=CPU_CORES={{system "cpu.cores"}}

[Install]
WantedBy=multi-user.target
```

**Template with Machine Data**:
```go
// templates/hosts.tmpl
# Generated hosts file for {{projectName}}

# Localhost
127.0.0.1 localhost

# Project machines
{{range machines}}
{{.Host}} {{.Name}} {{.Name}}.local
{{end}}

# Database connections
{{range machines}}
{{if eq (index .Tags "role") "database"}}
{{.Host}} db-{{.Name}} db-{{.Name}}.local
{{end}}
{{end}}
```

**Template with Conditional Logic**:
```go
// templates/config.tmpl
# Application configuration for {{var "app_name"}}

# Environment-specific settings
{{if eq (var "environment") "production"}}
DEBUG=false
LOG_LEVEL=warn
{{else}}
DEBUG=true
LOG_LEVEL=debug
{{end}}

# Database configuration
{{if eq (var "database_type") "postgresql"}}
DATABASE_URL=postgresql://{{var "db_user"}}:{{var "db_password"}}@{{var "db_host"}}:{{var "db_port"}}/{{var "db_name"}}
{{else if eq (var "database_type") "mysql"}}
DATABASE_URL=mysql://{{var "db_user"}}:{{var "db_password"}}@{{var "db_host"}}:{{var "db_port"}}/{{var "db_name"}}
{{end}}

# Machine-specific settings
HOSTNAME={{system "hostname"}}
CPU_CORES={{system "cpu.cores"}}
MEMORY_TOTAL={{system "memory.total"}}
```

### **5. Template Validation and Preview**

**Template Validation**:
```go
// ValidateTemplate validates template syntax and structure
func ValidateTemplate(templateFile string) error {
    // Read template file
    content, err := os.ReadFile(templateFile)
    if err != nil {
        return fmt.Errorf("error reading template file: %w", err)
    }
    
    // Create template with basic functions
    tmpl, err := template.New(filepath.Base(templateFile)).Funcs(template.FuncMap{
        "custom": func(path string) interface{} { return "mock-value" },
        "system": func(path string) interface{} { return "mock-value" },
        "env": func(key string) string { return "mock-value" },
        "data": func(path string) interface{} { return "mock-value" },
        "var": func(name string) interface{} { return "mock-value" },
    }).Parse(string(content))
    
    if err != nil {
        return fmt.Errorf("template syntax error: %w", err)
    }
    
    return nil
}
```

**Template Preview**:
```go
// PreviewTemplate provides template analysis and sample rendering
func PreviewTemplate(templateFile string) error {
    // Read template content
    content, err := os.ReadFile(templateFile)
    if err != nil {
        return fmt.Errorf("error reading template file: %w", err)
    }
    
    fmt.Println("=== TEMPLATE PREVIEW ===")
    fmt.Printf("Template: %s\n", filepath.Base(templateFile))
    fmt.Printf("Size: %d bytes\n", len(content))
    
    // Analyze template variables and functions
    variables := detectTemplateVariables(string(content))
    if len(variables) > 0 {
        fmt.Println("Variables/Functions detected:")
        for _, v := range variables {
            fmt.Printf("  - %s\n", v)
        }
    }
    
    // Show available template functions
    fmt.Println("Available template functions:")
    fmt.Println("  - custom(path)     - Access custom facts")
    fmt.Println("  - system(path)     - Access system facts")
    fmt.Println("  - env(key)         - Access environment variables")
    fmt.Println("  - data(path)       - Access additional data")
    fmt.Println("  - var(name)        - Access project variables")
    
    // Render with mock data
    mockData := createMockTemplateData()
    rendered, err := renderWithMockData(string(content), mockData)
    if err != nil {
        fmt.Printf("Error rendering with mock data: %v\n", err)
    } else {
        fmt.Println("Sample rendering (with mock data):")
        fmt.Println(rendered)
    }
    
    return nil
}
```

## CLI Integration

### **1. Template Management Commands**

**List Templates**:
```bash
# List all templates in project
spooky templates list ./my-project

# List templates with details
spooky templates list ./my-project --verbose

# List templates by pattern
spooky templates list ./my-project --pattern "*.conf.tmpl"

# List templates with size information
spooky templates list ./my-project --show-size
```

**Validate Templates**:
```bash
# Validate all templates in project
spooky templates validate ./my-project

# Validate specific template
spooky templates validate ./my-project --template templates/nginx.conf.tmpl

# Validate multiple templates
spooky templates validate ./my-project --templates templates/*.tmpl

# Validate with detailed output
spooky templates validate ./my-project --verbose
```

**Render Templates**:
```bash
# Basic rendering
spooky templates render ./my-project templates/nginx.conf.tmpl

# Render with specific data file
spooky templates render ./my-project templates/config.tmpl --data data/variables.hcl

# Render with data directory
spooky templates render ./my-project templates/config.tmpl --data data/

# Render with machine facts
spooky templates render ./my-project templates/config.tmpl --machine web-01

# Render with output file
spooky templates render ./my-project templates/nginx.conf.tmpl --output /etc/nginx/nginx.conf

# Preview mode
spooky templates render ./my-project templates/config.tmpl --preview

# Dry run mode
spooky templates render ./my-project templates/config.tmpl --dry-run

# Combined example
spooky templates render ./my-project templates/config.tmpl --machine web-01 --data data/facts.hcl --dry-run
```

### **2. Template Analysis Commands**

**Template Analysis**:
```bash
# Analyze template structure
spooky templates analyze ./my-project templates/nginx.conf.tmpl

# Analyze all templates
spooky templates analyze ./my-project

# Analyze with dependency graph
spooky templates analyze ./my-project --show-dependencies

# Analyze with complexity metrics
spooky templates analyze ./my-project --show-complexity
```

**Template Debugging**:
```bash
# Debug template rendering
spooky templates debug ./my-project templates/config.tmpl

# Debug with specific data
spooky templates debug ./my-project templates/config.tmpl --data data/variables.hcl

# Debug with step-by-step running
spooky templates debug ./my-project templates/config.tmpl --step-by-step
```

### **3. Template Development Commands**

**Template Scaffolding**:
```bash
# Create new template
spooky templates create ./my-project templates/new-config.tmpl

# Create template from example
spooky templates create ./my-project templates/nginx.conf.tmpl --from-example nginx

# Create template with variables
spooky templates create ./my-project templates/app.conf.tmpl --with-variables
```

**Template Testing**:
```bash
# Test template with sample data
spooky templates test ./my-project templates/config.tmpl

# Test template with multiple datasets
spooky templates test ./my-project templates/config.tmpl --test-data test-data/

# Test template validation
spooky templates test ./my-project templates/config.tmpl --validate-only
```

## Implementation Details

### **1. Template Context Management**

**Context Creation**:
```go
// NewTemplateContext creates a new template context for a project
func NewTemplateContext(logger logging.Logger, projectPath string) (*TemplateContext, error) {
    ctx := &TemplateContext{
        Facts:       make(map[string]interface{}),
        Environment: make(map[string]string),
        CustomData:  make(map[string]interface{}),
    }
    
    logger.Info("Creating template context",
        logging.String("project_path", projectPath))
    
    // Load project configuration
    if err := ctx.loadProjectConfig(logger, projectPath); err != nil {
        return nil, fmt.Errorf("failed to load project config: %w", err)
    }
    
    // Load facts
    if err := ctx.loadFacts(logger, projectPath); err != nil {
        logger.Warn("Failed to load facts", logging.String("error", err.Error()))
    }
    
    // Load inventory
    if err := ctx.loadInventory(logger, projectPath); err != nil {
        logger.Warn("Failed to load inventory", logging.String("error", err.Error()))
    }
    
    // Load actions
    if err := ctx.loadActions(logger, projectPath); err != nil {
        logger.Warn("Failed to load actions", logging.String("error", err.Error()))
    }
    
    // Load environment variables
    ctx.loadEnvironment()
    
    // Load custom data files
    if err := ctx.loadCustomData(logger, projectPath); err != nil {
        logger.Warn("Failed to load custom data", logging.String("error", err.Error()))
    }
    
    // Load variables
    if err := ctx.loadVariables(logger, projectPath); err != nil {
        logger.Warn("Failed to load variables", logging.String("error", err.Error()))
    }
    
    return ctx, nil
}
```

**Context Loading**:
```go
// loadCustomData loads custom data files from data/ directory
func (ctx *TemplateContext) loadCustomData(logger logging.Logger, projectPath string) error {
    dataPath := filepath.Join(projectPath, "data")
    
    // Check if data directory exists
    if _, err := os.Stat(dataPath); os.IsNotExist(err) {
        logger.Info("Data directory does not exist, skipping custom data loading")
        return nil
    }
    
    // Load all HCL files from data directory
    files, err := filepath.Glob(filepath.Join(dataPath, "*.hcl"))
    if err != nil {
        return fmt.Errorf("error finding data files: %w", err)
    }
    
    for _, file := range files {
        logger.Info("Loading custom data file", logging.String("file", file))
        
        data, err := ParseHCLDataFile(file)
        if err != nil {
            logger.Error("Failed to parse data file", err, logging.String("file", file))
            continue
        }
        
        // Use filename (without extension) as namespace
        namespace := strings.TrimSuffix(filepath.Base(file), ".hcl")
        ctx.CustomData[namespace] = data
    }
    
    logger.Info("Custom data loaded successfully",
        logging.Int("files", len(files)),
        logging.Int("namespaces", len(ctx.CustomData)))
    
    return nil
}
```

### **2. Template Rendering Engine**

**Template Preprocessing**:
```go
// preprocessTemplateContent converts single quotes to double quotes in template functions
func (te *TemplateEngine) preprocessTemplateContent(content string) string {
    // Replace {{custom '...'}} with {{custom "..."}}
    content = strings.ReplaceAll(content, "{{custom '", "{{custom \"")
    content = strings.ReplaceAll(content, "'}}", "\"}}")
    
    // Replace {{system '...'}} with {{system "..."}}
    content = strings.ReplaceAll(content, "{{system '", "{{system \"")
    
    // Replace {{env '...'}} with {{env "..."}}
    content = strings.ReplaceAll(content, "{{env '", "{{env \"")
    
    // Replace {{data '...'}} with {{data "..."}}
    content = strings.ReplaceAll(content, "{{data '", "{{data \"")
    
    return content
}
```

**Nested Value Access**:
```go
// getNestedValue gets a nested value from a map using dot notation
func (te *TemplateEngine) getNestedValue(data map[string]interface{}, path string) interface{} {
    parts := strings.Split(path, ".")
    current := data
    
    for i, part := range parts {
        if current == nil {
            return nil
        }
        
        if i == len(parts)-1 {
            // Last part - return the value
            return current[part]
        }
        
        // Navigate deeper
        if next, ok := current[part].(map[string]interface{}); ok {
            current = next
        } else {
            return nil
        }
    }
    
    return nil
}
```

### **3. Template Validation and Error Handling**

**Template Validation**:
```go
// ValidateTemplate validates template syntax and structure
func ValidateTemplate(templateFile string) ([]string, error) {
    var errors []string
    
    // Read template file
    content, err := os.ReadFile(templateFile)
    if err != nil {
        return nil, fmt.Errorf("error reading template file: %w", err)
    }
    
    // Check for basic syntax issues
    if strings.Contains(string(content), "{{") && !strings.Contains(string(content), "}}") {
        errors = append(errors, "unclosed template expression")
    }
    
    if strings.Contains(string(content), "}}") && !strings.Contains(string(content), "{{") {
        errors = append(errors, "unopened template expression")
    }
    
    // Try to parse template
    tmpl, err := template.New(filepath.Base(templateFile)).Funcs(template.FuncMap{
        "custom": func(path string) interface{} { return "mock-value" },
        "system": func(path string) interface{} { return "mock-value" },
        "env": func(key string) string { return "mock-value" },
        "data": func(path string) interface{} { return "mock-value" },
        "var": func(name string) interface{} { return "mock-value" },
    }).Parse(string(content))
    
    if err != nil {
        errors = append(errors, fmt.Sprintf("template syntax error: %v", err))
    }
    
    // Check for undefined functions
    undefinedFuncs := detectUndefinedFunctions(string(content))
    for _, funcName := range undefinedFuncs {
        errors = append(errors, fmt.Sprintf("undefined function: %s", funcName))
    }
    
    return errors, nil
}
```

**Error Reporting**:
```go
// ReportTemplateErrors reports template errors with context
func ReportTemplateErrors(templateFile string, errors []string) {
    fmt.Printf("Template validation failed for %s:\n", templateFile)
    
    for i, err := range errors {
        fmt.Printf("  %d. %s\n", i+1, err)
    }
    
    if len(errors) > 0 {
        fmt.Printf("\nTotal errors: %d\n", len(errors))
    }
}
```

## Performance Considerations

### **1. Template Caching**

**Template Cache**:
```go
// TemplateCache provides caching for parsed templates
type TemplateCache struct {
    templates map[string]*template.Template
    mutex     sync.RWMutex
    maxSize   int
}

// GetTemplate gets or creates a cached template
func (tc *TemplateCache) GetTemplate(templateFile string, funcMap template.FuncMap) (*template.Template, error) {
    tc.mutex.RLock()
    if tmpl, exists := tc.templates[templateFile]; exists {
        defer tc.mutex.RUnlock()
        return tmpl, nil
    }
    tc.mutex.RUnlock()
    
    tc.mutex.Lock()
    defer tc.mutex.Unlock()
    
    // Check again after acquiring write lock
    if tmpl, exists := tc.templates[templateFile]; exists {
        return tmpl, nil
    }
    
    // Parse template
    content, err := os.ReadFile(templateFile)
    if err != nil {
        return nil, err
    }
    
    tmpl, err := template.New(filepath.Base(templateFile)).Funcs(funcMap).Parse(string(content))
    if err != nil {
        return nil, err
    }
    
    // Cache template
    if len(tc.templates) >= tc.maxSize {
        // Remove oldest template (simple LRU)
        for key := range tc.templates {
            delete(tc.templates, key)
            break
        }
    }
    
    tc.templates[templateFile] = tmpl
    return tmpl, nil
}
```

### **2. Template Optimization**

**Template Compilation**:
```go
// CompileTemplate compiles template for better performance
func CompileTemplate(templateFile string) (*template.Template, error) {
    content, err := os.ReadFile(templateFile)
    if err != nil {
        return nil, err
    }
    
    // Preprocess template
    contentStr := preprocessTemplateContent(string(content))
    
    // Create template with optimized functions
    tmpl, err := template.New(filepath.Base(templateFile)).Funcs(template.FuncMap{
        "custom": optimizedCustomFunc,
        "system": optimizedSystemFunc,
        "env":    optimizedEnvFunc,
        "data":   optimizedDataFunc,
        "var":    optimizedVarFunc,
    }).Parse(contentStr)
    
    if err != nil {
        return nil, err
    }
    
    return tmpl, nil
}
```

### **3. Parallel Template Rendering**

**Parallel Rendering**:
```go
// RenderTemplatesParallel renders multiple templates in parallel
func RenderTemplatesParallel(templates []string, context *TemplateContext, maxParallel int) (map[string]string, error) {
    semaphore := make(chan struct{}, maxParallel)
    var wg sync.WaitGroup
    results := make(map[string]string)
    errors := make(chan error, len(templates))
    resultMutex := sync.Mutex{}
    
    for _, templateFile := range templates {
        wg.Add(1)
        go func(tmpl string) {
            defer wg.Done()
            semaphore <- struct{}{}
            defer func() { <-semaphore }()
            
            rendered, err := RenderTemplate(tmpl, context)
            if err != nil {
                errors <- fmt.Errorf("template %s: %w", tmpl, err)
                return
            }
            
            resultMutex.Lock()
            results[tmpl] = rendered
            resultMutex.Unlock()
        }(templateFile)
    }
    
    wg.Wait()
    close(errors)
    
    // Collect errors
    var errs []error
    for err := range errors {
        errs = append(errs, err)
    }
    
    if len(errs) > 0 {
        return results, fmt.Errorf("parallel rendering failed: %v", errs)
    }
    
    return results, nil
}
```

## Security Considerations

### **1. Template Sandboxing**

**Function Restrictions**:
```go
// SafeTemplateFunctions provides safe template functions
func SafeTemplateFunctions() template.FuncMap {
    return template.FuncMap{
        // Safe data access functions
        "custom": safeCustomFunc,
        "system": safeSystemFunc,
        "env":    safeEnvFunc,
        "data":   safeDataFunc,
        "var":    safeVarFunc,
        
        // Safe string functions
        "upper": strings.ToUpper,
        "lower": strings.ToLower,
        "trim":  strings.TrimSpace,
        
        // Safe math functions
        "add": func(a, b int) int { return a + b },
        "sub": func(a, b int) int { return a - b },
        "mul": func(a, b int) int { return a * b },
        "div": func(a, b int) int { return a / b },
    }
}
```

**Access Control**:
```go
// TemplateAccessControl controls template access
type TemplateAccessControl struct {
    allowedPaths []string
    blockedPaths []string
    maxFileSize  int64
}

// ValidateTemplateAccess validates template access
func (tac *TemplateAccessControl) ValidateTemplateAccess(templateFile string) error {
    // Check file size
    info, err := os.Stat(templateFile)
    if err != nil {
        return err
    }
    
    if info.Size() > tac.maxFileSize {
        return fmt.Errorf("template file too large: %d bytes", info.Size())
    }
    
    // Check allowed paths
    for _, allowedPath := range tac.allowedPaths {
        if strings.HasPrefix(templateFile, allowedPath) {
            return nil
        }
    }
    
    // Check blocked paths
    for _, blockedPath := range tac.blockedPaths {
        if strings.HasPrefix(templateFile, blockedPath) {
            return fmt.Errorf("access to template file blocked: %s", templateFile)
        }
    }
    
    return fmt.Errorf("template file not in allowed paths: %s", templateFile)
}
```

### **2. Template Injection Prevention**

**Input Validation**:
```go
// ValidateTemplateInput validates template input data
func ValidateTemplateInput(data map[string]interface{}) error {
    for key, value := range data {
        // Check for potentially dangerous values
        if str, ok := value.(string); ok {
            if containsDangerousContent(str) {
                return fmt.Errorf("dangerous content detected in key '%s'", key)
            }
        }
        
        // Recursively check nested maps
        if nested, ok := value.(map[string]interface{}); ok {
            if err := ValidateTemplateInput(nested); err != nil {
                return fmt.Errorf("nested key '%s': %w", key, err)
            }
        }
    }
    
    return nil
}
```

**Content Sanitization**:
```go
// SanitizeTemplateContent sanitizes template content
func SanitizeTemplateContent(content string) string {
    // Remove potentially dangerous patterns
    dangerousPatterns := []string{
        `{{.*os\.Exec.*}}`,
        `{{.*run.*}}`,
        `{{.*system.*}}`,
        `{{.*eval.*}}`,
    }
    
    for _, pattern := range dangerousPatterns {
        re := regexp.MustCompile(pattern)
        content = re.ReplaceAllString(content, "")
    }
    
    return content
}
```

## Future Enhancements

### **1. Advanced Template Features**

**Template Inheritance**:
- **Base Templates**: Base templates for common functionality
- **Template Composition**: Template composition and reuse
- **Template Macros**: Template macros for common patterns
- **Template Libraries**: Template libraries and repositories

**Template Versioning**:
- **Template History**: Template change history and versioning
- **Template Rollback**: Template rollback capabilities
- **Template Comparison**: Template comparison and diff tools
- **Template Migration**: Template migration and upgrade tools

### **2. Template Development Tools**

**Template IDE Integration**:
- **Syntax Highlighting**: Template syntax highlighting
- **Auto-Completion**: Template function auto-completion
- **Error Detection**: Real-time template error detection
- **Debugging Tools**: Template debugging and step-through

**Template Testing Framework**:
- **Unit Testing**: Template unit testing framework
- **Integration Testing**: Template integration testing
- **Mock Data**: Mock data generation for testing
- **Test Coverage**: Template test coverage reporting

### **3. Template Performance Enhancements**

**Advanced Caching**:
- **Intelligent Caching**: Intelligent template caching strategies
- **Cache Invalidation**: Automatic cache invalidation
- **Distributed Caching**: Distributed template caching
- **Cache Analytics**: Template cache performance analytics

**Template Optimization**:
- **Template Compilation**: Template compilation for better performance
- **Lazy Loading**: Lazy loading of template data
- **Parallel Processing**: Enhanced parallel template processing
- **Memory Optimization**: Template memory usage optimization