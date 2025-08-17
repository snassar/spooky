# Code Walkthrough: `spooky templates validate`

## Overview

This document provides a detailed function-by-function walkthrough of the `spooky templates validate` command execution flow. This command validates template syntax and variables.

## Command Structure

```bash
spooky templates validate <project-path> [--template <path>] [--data <file>]
```

## Execution Flow

### 1. Entry Point: `main()`

**File**: `main.go`
```go
func main() {
    cmd.Execute()
}
```

**Purpose**: Entry point that calls the CLI execution system.

### 2. CLI Setup: `cmd.Execute()`

**File**: `cmd/root.go`
```go
func Execute() {
    err := RootCmd.Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }
}
```

**Purpose**: Executes the root Cobra command with all subcommands.

### 3. Pre-Run Setup: `RootCmd.PersistentPreRunE`

**File**: `cmd/root.go`
```go
PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
    // Skip auto-setup for version and help commands
    if cmd.Name() == "version" || cmd.Name() == "help" {
        return nil
    }

    // Auto-setup configuration for all other commands
    if err := spookyconfig.AutoSetupConfig(); err != nil {
        return fmt.Errorf("configuration setup failed: %w", err)
    }

    // Initialize integration manager for commands that need it
    if err := InitializeIntegrationsDependencies(); err != nil {
        return fmt.Errorf("integration initialization failed: %w", err)
    }

    return nil
}
```

**Purpose**: 
- Sets up global configuration automatically
- Initializes integration dependencies
- Runs before any command execution

### 4. Command Routing: `templatesValidateCmd.RunE`

**File**: `cmd/templates.go`
```go
var templatesValidateCmd = &cobra.Command{
    Use:   "validate [project-path]",
    Short: "Validate template syntax and variables",
    Long: `Validate template syntax and variables.

This command validates that all templates in the project have proper syntax
and that required variables are properly defined.`,
    Args: cobra.ExactArgs(1),
    RunE: func(cmd *cobra.Command, args []string) error {
        projectPath := args[0]
        templatePath, _ := cmd.Flags().GetString("template")
        dataFile, _ := cmd.Flags().GetString("data")

        return handleTemplatesValidate(projectPath, templatePath, dataFile)
    },
}
```

**Purpose**: Routes the command to `handleTemplatesValidate` with the project path and optional flags.

### 5. Main Handler: `handleTemplatesValidate()`

**File**: `cmd/templates.go`
```go
func handleTemplatesValidate(projectPath, templatePath, dataFile string) error {
    ctx := context.Background()

    // Initialize dependencies if not already done
    if templatesManager == nil {
        if err := InitializeTemplatesDependencies(); err != nil {
            return fmt.Errorf("failed to initialize templates dependencies: %w", err)
        }
    }

    fmt.Printf("🔍 Validating templates in project: %s\n", projectPath)

    // Get templates integration
    templatesIntegration := getIntegrationManager().GetTemplatesIntegration()
    if templatesIntegration == nil {
        return fmt.Errorf("templates integration not available")
    }

    // Load templates from project
    templates, err := templatesIntegration.LoadTemplates(ctx, projectPath)
    if err != nil {
        return fmt.Errorf("failed to load templates: %w", err)
    }

    if len(templates) == 0 {
        fmt.Printf("No templates found to validate.\n")
        return nil
    }

    // Filter templates if specific template is requested
    if templatePath != "" {
        templates = filterTemplatesByPath(templates, templatePath)
        if len(templates) == 0 {
            return fmt.Errorf("no templates found matching path: %s", templatePath)
        }
    }

    fmt.Printf("📊 Validating %d templates...\n", len(templates))

    // Load data file if provided
    var data map[string]interface{}
    if dataFile != "" {
        data, err = loadDataFile(dataFile)
        if err != nil {
            return fmt.Errorf("failed to load data file: %w", err)
        }
    }

    // Validate templates
    result, err := templatesIntegration.ValidateTemplates(ctx, templates, data)
    if err != nil {
        return fmt.Errorf("validation failed: %w", err)
    }

    // Display validation results
    fmt.Printf("\n✅ Validation Results:\n")
    fmt.Printf("%s\n", strings.Repeat("─", 50))

    if len(result.Errors) == 0 && len(result.Warnings) == 0 {
        fmt.Printf("🎉 All templates are valid!\n")
        return nil
    }

    // Display errors
    if len(result.Errors) > 0 {
        fmt.Printf("❌ Errors (%d):\n", len(result.Errors))
        for i := range result.Errors {
            err := &result.Errors[i]
            fmt.Printf("  %d. %s\n", i+1, err.Message)
            if err.Context != nil {
                fmt.Printf("     Context: %v\n", err.Context)
            }
        }
        fmt.Printf("\n")
    }

    // Display warnings
    if len(result.Warnings) > 0 {
        fmt.Printf("⚠️  Warnings (%d):\n", len(result.Warnings))
        for i := range result.Warnings {
            warning := &result.Warnings[i]
            fmt.Printf("  %d. %s\n", i+1, warning.Message)
            if warning.Context != nil {
                fmt.Printf("     Context: %v\n", warning.Context)
            }
        }
        fmt.Printf("\n")
    }

    // Summary
    if len(result.Errors) > 0 {
        fmt.Printf("❌ Validation failed with %d errors\n", len(result.Errors))
        return fmt.Errorf("template validation failed")
    }

    fmt.Printf("✅ Validation completed successfully\n")
    return nil
}
```

**Purpose**: 
- Initializes templates dependencies
- Loads templates from project
- Filters templates if specific path provided
- Loads data file if provided
- Validates templates
- Reports detailed results to the user

### 6. Dependency Initialization: `InitializeTemplatesDependencies()`

**File**: `cmd/templates.go`
```go
func InitializeTemplatesDependencies() error {
    // Create log manager for templates component
    logManager := spookylogging.NewLogManager()
    templatesLogger = logManager.GetLogger("templates")

    // Initialize templates components
    loader := spookytemplates.NewLoader(templatesLogger)
    validator := spookytemplates.NewValidator(templatesLogger)
    engine := spookytemplates.NewEngine(templatesLogger)
    manager := spookytemplates.NewManager(loader, validator, engine, templatesLogger)

    // Create templates integration
    templatesManager = spookytemplates.NewIntegration(manager)

    return nil
}
```

**Purpose**: Creates logger, loader, validator, engine, and manager instances for dependency injection.

### 7. Templates Loading: `templatesIntegration.LoadTemplates()`

**File**: `internal/templates/integration.go`
```go
func (i *TemplatesIntegration) LoadTemplates(ctx context.Context, projectPath string) ([]spookytypes.Template, error) {
    i.logger.Info("Loading templates from project", map[string]interface{}{
        "project_path": projectPath,
    })

    // Use the manager to load templates
    templates, err := i.manager.LoadTemplates(ctx, projectPath)
    if err != nil {
        return nil, fmt.Errorf("failed to load templates: %w", err)
    }

    i.logger.Info("Templates loaded successfully", map[string]interface{}{
        "project_path": projectPath,
        "template_count": len(templates),
    })

    return templates, nil
}
```

**Purpose**: Delegates template loading to the templates manager.

### 8. Deep Templates Loading: `manager.LoadTemplates()`

**File**: `internal/templates/manager.go`
```go
func (m *Manager) LoadTemplates(ctx context.Context, projectPath string) ([]spookytypes.Template, error) {
    m.logger.Info("Loading templates from project", map[string]interface{}{
        "project_path": projectPath,
    })

    // Use the loader to load templates
    templates, err := m.loader.LoadTemplates(ctx, projectPath)
    if err != nil {
        return nil, fmt.Errorf("failed to load templates: %w", err)
    }

    m.logger.Info("Templates loaded successfully", map[string]interface{}{
        "project_path": projectPath,
        "template_count": len(templates),
    })

    return templates, nil
}
```

**Purpose**: Delegates template loading to the templates loader.

### 9. Templates File Loading: `loader.LoadTemplates()`

**File**: `internal/templates/loader.go`
```go
func (l *Loader) LoadTemplates(ctx context.Context, projectPath string) ([]spookytypes.Template, error) {
    l.logger.Info("Loading templates from project", map[string]interface{}{
        "project_path": projectPath,
    })

    var allTemplates []spookytypes.Template
    var loadErrors []string

    // Check for templates/ directory
    templatesDir := filepath.Join(projectPath, "templates")
    if _, err := os.Stat(templatesDir); err == nil {
        l.logger.Debug("Found templates directory", map[string]interface{}{
            "dir_path": templatesDir,
        })

        templates, err := l.loadTemplatesFromDirectory(ctx, templatesDir)
        if err != nil {
            loadErrors = append(loadErrors, fmt.Sprintf("templates/ directory: %v", err))
        } else {
            allTemplates = append(allTemplates, templates...)
            l.logger.Info("Loaded templates from directory", map[string]interface{}{
                "dir_path": templatesDir,
                "count":    len(templates),
            })
        }
    }

    // Check if no templates were found
    if len(allTemplates) == 0 {
        if len(loadErrors) > 0 {
            return nil, fmt.Errorf("failed to load templates: %s", strings.Join(loadErrors, "; "))
        }
        return nil, fmt.Errorf("no templates found in project: %s (templates/ directory not found)", projectPath)
    }

    // Validate for duplicates and consistency
    if err := l.validateTemplateCollection(allTemplates); err != nil {
        return nil, fmt.Errorf("template collection validation failed: %w", err)
    }

    l.logger.Info("Templates loaded successfully", map[string]interface{}{
        "project_path": projectPath,
        "total_count":  len(allTemplates),
    })

    return allTemplates, nil
}
```

**Purpose**: 
- Checks for templates/ directory
- Loads templates from directory
- Validates template collection
- Handles loading errors

### 10. Templates Validation: `templatesIntegration.ValidateTemplates()`

**File**: `internal/templates/integration.go`
```go
func (i *TemplatesIntegration) ValidateTemplates(ctx context.Context, templates []spookytypes.Template, data map[string]interface{}) (*spookytypes.ValidationResult, error) {
    i.logger.Info("Validating templates", map[string]interface{}{
        "template_count": len(templates),
        "has_data":       data != nil,
    })

    // Use the manager to validate templates
    result, err := i.manager.ValidateTemplates(ctx, templates, data)
    if err != nil {
        return nil, fmt.Errorf("failed to validate templates: %w", err)
    }

    if result.Valid {
        i.logger.Info("Templates validation passed", map[string]interface{}{
            "template_count": len(templates),
        })
    } else {
        i.logger.Warn("Templates validation failed", map[string]interface{}{
            "template_count": len(templates),
            "errors":         len(result.Errors),
            "warnings":       len(result.Warnings),
        })
    }

    return result, nil
}
```

**Purpose**: Delegates template validation to the templates manager.

### 11. Deep Templates Validation: `manager.ValidateTemplates()`

**File**: `internal/templates/manager.go`
```go
func (m *Manager) ValidateTemplates(ctx context.Context, templates []spookytypes.Template, data map[string]interface{}) (*spookytypes.ValidationResult, error) {
    m.logger.Info("Validating templates", map[string]interface{}{
        "template_count": len(templates),
        "has_data":       data != nil,
    })

    // Use the validator to validate templates
    result, err := m.validator.ValidateTemplates(ctx, templates, data)
    if err != nil {
        return nil, fmt.Errorf("failed to validate templates: %w", err)
    }

    if result.Valid {
        m.logger.Info("Templates validation passed", map[string]interface{}{
            "template_count": len(templates),
        })
    } else {
        m.logger.Warn("Templates validation failed", map[string]interface{}{
            "template_count": len(templates),
            "errors":         len(result.Errors),
            "warnings":       len(result.Warnings),
        })
    }

    return result, nil
}
```

**Purpose**: Delegates template validation to the templates validator.

### 12. Templates Validator: `validator.ValidateTemplates()`

**File**: `internal/templates/validator.go`
```go
func (v *Validator) ValidateTemplates(ctx context.Context, templates []spookytypes.Template, data map[string]interface{}) (*spookytypes.ValidationResult, error) {
    v.logger.Info("Validating templates", map[string]interface{}{
        "template_count": len(templates),
        "has_data":       data != nil,
    })

    result := &spookytypes.ValidationResult{
        Valid:       true,
        ValidatedAt: time.Now(),
        Errors:      []spookytypesschemas.SchemaError{},
        Warnings:    []spookytypesschemas.SchemaError{},
    }

    // Validate each template individually
    for i := range templates {
        templateResult, err := v.ValidateTemplate(ctx, &templates[i], data)
        if err != nil {
            return nil, fmt.Errorf("failed to validate template %d: %w", i, err)
        }

        // Merge validation results
        if !templateResult.Valid {
            result.Valid = false
        }
        result.Errors = append(result.Errors, templateResult.Errors...)
        result.Warnings = append(result.Warnings, templateResult.Warnings...)
    }

    // Validate cross-template dependencies
    if err := v.validateTemplateDependencies(ctx, templates, result); err != nil {
        return nil, fmt.Errorf("failed to validate template dependencies: %w", err)
    }

    return result, nil
}
```

**Purpose**: 
- Validates each template individually
- Validates cross-template dependencies
- Aggregates validation results

### 13. Individual Template Validation: `validator.ValidateTemplate()`

**File**: `internal/templates/validator.go`
```go
func (v *Validator) ValidateTemplate(ctx context.Context, template *spookytypes.Template, data map[string]interface{}) (*spookytypes.ValidationResult, error) {
    v.logger.Debug("Validating template", map[string]interface{}{
        "template_name": template.Name,
        "template_path": template.Path,
    })

    result := &spookytypes.ValidationResult{
        Valid:       true,
        ValidatedAt: time.Now(),
        Errors:      []spookytypesschemas.SchemaError{},
        Warnings:    []spookytypesschemas.SchemaError{},
    }

    // Validate template name
    if err := v.validateTemplateName(template.Name); err != nil {
        result.Valid = false
        result.Errors = append(result.Errors, spookytypesschemas.SchemaError{
            Code:     "invalid_template_name",
            Message:  fmt.Sprintf("Invalid template name: %s", err.Error()),
            Severity: "error",
        })
    }

    // Validate template path
    if err := v.validateTemplatePath(template.Path); err != nil {
        result.Valid = false
        result.Errors = append(result.Errors, spookytypesschemas.SchemaError{
            Code:     "invalid_template_path",
            Message:  fmt.Sprintf("Invalid template path: %s", err.Error()),
            Severity: "error",
        })
    }

    // Validate template content
    if err := v.validateTemplateContent(template.Content); err != nil {
        result.Valid = false
        result.Errors = append(result.Errors, spookytypesschemas.SchemaError{
            Code:     "invalid_template_content",
            Message:  fmt.Sprintf("Invalid template content: %s", err.Error()),
            Severity: "error",
        })
    }

    // Validate template variables
    if err := v.validateTemplateVariables(template.Variables, data); err != nil {
        result.Valid = false
        result.Errors = append(result.Errors, spookytypesschemas.SchemaError{
            Code:     "invalid_template_variables",
            Message:  fmt.Sprintf("Invalid template variables: %s", err.Error()),
            Severity: "error",
        })
    }

    // Validate template functions
    if err := v.validateTemplateFunctions(template.Functions); err != nil {
        result.Warnings = append(result.Warnings, spookytypesschemas.SchemaError{
            Code:     "invalid_template_functions",
            Message:  fmt.Sprintf("Invalid template functions: %s", err.Error()),
            Severity: "warning",
        })
    }

    // Test template rendering if data is provided
    if data != nil {
        if err := v.testTemplateRendering(template, data); err != nil {
            result.Warnings = append(result.Warnings, spookytypesschemas.SchemaError{
                Code:     "template_rendering_test_failed",
                Message:  fmt.Sprintf("Template rendering test failed: %s", err.Error()),
                Severity: "warning",
            })
        }
    }

    return result, nil
}
```

**Purpose**: 
- Validates template name format
- Validates template path
- Validates template content syntax
- Validates template variables
- Validates template functions
- Tests template rendering with data

## Key Components

### Templates Integration
- **Purpose**: Provides templates functionality to the CLI
- **Responsibilities**: Loading, validation, rendering
- **Interface**: `spookyinterfaces.TemplatesIntegration`

### Templates Manager
- **Purpose**: Coordinates templates operations
- **Responsibilities**: Loading, validation, rendering
- **Interface**: `spookyinterfaces.TemplatesManager`

### Templates Validator
- **Purpose**: Validates template configurations
- **Responsibilities**: Syntax validation, variable validation, function validation
- **Features**: Individual and collection validation

### Templates Loader
- **Purpose**: Loads templates from disk
- **Responsibilities**: File loading, directory scanning, content parsing
- **Features**: Directory-based template loading

## Validation Process

1. **Project Validation**
   - Validates project path exists
   - Ensures project.hcl is present
   - Checks project structure

2. **Templates Loading**
   - Loads from templates/ directory
   - Parses template files
   - Validates template collection

3. **Individual Template Validation**
   - Validates template name format
   - Validates template path
   - Validates template content syntax
   - Validates template variables
   - Validates template functions

4. **Collection Validation**
   - Validates cross-template dependencies
   - Validates template relationships
   - Checks for consistency across templates

5. **Rendering Test**
   - Tests template rendering with provided data
   - Validates variable resolution
   - Checks for rendering errors

## Error Handling

- **Graceful Loading**: Handles missing template directories gracefully
- **Detailed Error Reporting**: Provides specific error messages with context
- **Warning System**: Reports non-critical issues as warnings
- **Collection Validation**: Validates entire template collection

## Output Format

### Success Case
```
🔍 Validating templates in project: /path/to/project
📊 Validating 5 templates...

✅ Validation Results:
──────────────────────────────────────────────────
🎉 All templates are valid!
✅ Validation completed successfully
```

### Error Case
```
🔍 Validating templates in project: /path/to/project
📊 Validating 5 templates...

✅ Validation Results:
──────────────────────────────────────────────────
❌ Errors (2):
  1. Invalid template name: nginx.conf.tmpl (contains invalid characters)
     Context: template[0]
  2. Invalid template content: syntax error at line 15
     Context: template[1]

⚠️  Warnings (1):
  1. Template rendering test failed: variable 'port' not found
     Context: template[2]

❌ Validation failed with 2 errors
```

## Architecture Patterns

- **Interface-Based Design**: Uses interfaces for loose coupling
- **Dependency Injection**: Injects dependencies through constructors
- **Layered Architecture**: Separates concerns across multiple layers
- **Error Aggregation**: Collects and reports all validation issues
- **Directory-Based Loading**: Supports directory-based template organization
- **Rendering Testing**: Tests template rendering with provided data
