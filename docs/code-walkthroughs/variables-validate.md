# Code Walkthrough: `spooky variables validate`

## Overview

This document provides a detailed function-by-function walkthrough of the `spooky variables validate` command execution flow. This command validates variable definitions and dependencies.

## Command Structure

```bash
spooky variables validate <project-path>
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

### 4. Command Routing: `variablesValidateCmd.RunE`

**File**: `cmd/variables.go`
```go
var variablesValidateCmd = &cobra.Command{
    Use:   "validate [project-path]",
    Short: "Validate variable definitions",
    Long: `Validate variable definitions and dependencies.

This command validates that all variables in the project have proper configuration
including required fields, valid types, and dependency relationships.`,
    Args: cobra.ExactArgs(1),
    RunE: func(_ *cobra.Command, args []string) error {
        return handleVariablesValidate(args[0])
    },
}
```

**Purpose**: Routes the command to `handleVariablesValidate` with the project path.

### 5. Main Handler: `handleVariablesValidate()`

**File**: `cmd/variables.go`
```go
func handleVariablesValidate(projectPath string) error {
    ctx := context.Background()

    // Initialize dependencies if not already done
    if variablesManager == nil {
        if err := InitializeVariablesDependencies(); err != nil {
            return fmt.Errorf("failed to initialize variables dependencies: %w", err)
        }
    }

    fmt.Printf("🔍 Validating variables in project: %s\n", projectPath)

    // Get variables integration
    variablesIntegration := getIntegrationManager().GetVariablesIntegration()
    if variablesIntegration == nil {
        return fmt.Errorf("variables integration not available")
    }

    // Load variables from project
    variables, err := variablesIntegration.LoadVariables(ctx, projectPath)
    if err != nil {
        return fmt.Errorf("failed to load variables: %w", err)
    }

    if len(variables) == 0 {
        fmt.Printf("No variables found to validate.\n")
        return nil
    }

    fmt.Printf("📊 Validating %d variables...\n", len(variables))

    // Validate variables
    result, err := variablesIntegration.ValidateVariables(ctx, variables)
    if err != nil {
        return fmt.Errorf("validation failed: %w", err)
    }

    // Display validation results
    fmt.Printf("\n✅ Validation Results:\n")
    fmt.Printf("%s\n", strings.Repeat("─", 50))

    if len(result.Errors) == 0 && len(result.Warnings) == 0 {
        fmt.Printf("🎉 All variables are valid!\n")
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
        return fmt.Errorf("variables validation failed")
    }

    fmt.Printf("✅ Validation completed successfully\n")
    return nil
}
```

**Purpose**: 
- Initializes variables dependencies
- Loads variables from project
- Validates variables
- Reports detailed results to the user

### 6. Dependency Initialization: `InitializeVariablesDependencies()`

**File**: `cmd/variables.go`
```go
func InitializeVariablesDependencies() error {
    // Create log manager for variables component
    logManager := spookylogging.NewLogManager()
    variablesLogger = logManager.GetLogger("variables")

    // Initialize variables components
    loader := spookyvariables.NewLoader(variablesLogger)
    validator := spookyvariables.NewValidator(variablesLogger)
    resolver := spookyvariables.NewResolver(variablesLogger)
    manager := spookyvariables.NewManager(loader, validator, resolver, variablesLogger)

    // Create variables integration
    variablesManager = spookyvariables.NewIntegration(manager)

    return nil
}
```

**Purpose**: Creates logger, loader, validator, resolver, and manager instances for dependency injection.

### 7. Variables Loading: `variablesIntegration.LoadVariables()`

**File**: `internal/variables/integration.go`
```go
func (i *VariablesIntegration) LoadVariables(ctx context.Context, projectPath string) ([]spookytypes.Variable, error) {
    i.logger.Info("Loading variables from project", map[string]interface{}{
        "project_path": projectPath,
    })

    // Use the manager to load variables
    variables, err := i.manager.LoadVariables(ctx, projectPath)
    if err != nil {
        return nil, fmt.Errorf("failed to load variables: %w", err)
    }

    i.logger.Info("Variables loaded successfully", map[string]interface{}{
        "project_path": projectPath,
        "variable_count": len(variables),
    })

    return variables, nil
}
```

**Purpose**: Delegates variable loading to the variables manager.

### 8. Deep Variables Loading: `manager.LoadVariables()`

**File**: `internal/variables/manager.go`
```go
func (m *Manager) LoadVariables(ctx context.Context, projectPath string) ([]spookytypes.Variable, error) {
    m.logger.Info("Loading variables from project", map[string]interface{}{
        "project_path": projectPath,
    })

    // Use the loader to load variables
    variables, err := m.loader.LoadVariables(ctx, projectPath)
    if err != nil {
        return nil, fmt.Errorf("failed to load variables: %w", err)
    }

    m.logger.Info("Variables loaded successfully", map[string]interface{}{
        "project_path": projectPath,
        "variable_count": len(variables),
    })

    return variables, nil
}
```

**Purpose**: Delegates variable loading to the variables loader.

### 9. Variables File Loading: `loader.LoadVariables()`

**File**: `internal/variables/loader.go`
```go
func (l *Loader) LoadVariables(ctx context.Context, projectPath string) ([]spookytypes.Variable, error) {
    l.logger.Info("Loading variables from project", map[string]interface{}{
        "project_path": projectPath,
    })

    var allVariables []spookytypes.Variable
    var loadErrors []string

    // Check for variables.hcl file
    variablesFile := filepath.Join(projectPath, "variables.hcl")
    if _, err := os.Stat(variablesFile); err == nil {
        l.logger.Debug("Found variables.hcl file", map[string]interface{}{
            "file_path": variablesFile,
        })

        variables, err := l.loadVariablesFromFile(ctx, variablesFile)
        if err != nil {
            loadErrors = append(loadErrors, fmt.Sprintf("variables.hcl: %v", err))
        } else {
            allVariables = append(allVariables, variables...)
            l.logger.Info("Loaded variables from file", map[string]interface{}{
                "file_path": variablesFile,
                "count":     len(variables),
            })
        }
    }

    // Check for variables/ directory
    variablesDir := filepath.Join(projectPath, "variables")
    if _, err := os.Stat(variablesDir); err == nil {
        l.logger.Debug("Found variables directory", map[string]interface{}{
            "dir_path": variablesDir,
        })

        variables, err := l.loadVariablesFromDirectory(ctx, variablesDir)
        if err != nil {
            loadErrors = append(loadErrors, fmt.Sprintf("variables/ directory: %v", err))
        } else {
            allVariables = append(allVariables, variables...)
            l.logger.Info("Loaded variables from directory", map[string]interface{}{
                "dir_path": variablesDir,
                "count":    len(variables),
            })
        }
    }

    // Check if no variables were found
    if len(allVariables) == 0 {
        if len(loadErrors) > 0 {
            return nil, fmt.Errorf("failed to load variables: %s", strings.Join(loadErrors, "; "))
        }
        return nil, fmt.Errorf("no variables found in project: %s (neither variables.hcl nor variables/ directory found)", projectPath)
    }

    // Validate for duplicates and consistency
    if err := l.validateVariableCollection(allVariables); err != nil {
        return nil, fmt.Errorf("variable collection validation failed: %w", err)
    }

    l.logger.Info("Variables loaded successfully", map[string]interface{}{
        "project_path": projectPath,
        "total_count":  len(allVariables),
    })

    return allVariables, nil
}
```

**Purpose**: 
- Checks for variables.hcl file
- Checks for variables/ directory
- Loads variables from both sources
- Validates variable collection
- Handles loading errors

### 10. Variables Validation: `variablesIntegration.ValidateVariables()`

**File**: `internal/variables/integration.go`
```go
func (i *VariablesIntegration) ValidateVariables(ctx context.Context, variables []spookytypes.Variable) (*spookytypes.ValidationResult, error) {
    i.logger.Info("Validating variables", map[string]interface{}{
        "variable_count": len(variables),
    })

    // Use the manager to validate variables
    result, err := i.manager.ValidateVariables(ctx, variables)
    if err != nil {
        return nil, fmt.Errorf("failed to validate variables: %w", err)
    }

    if result.Valid {
        i.logger.Info("Variables validation passed", map[string]interface{}{
            "variable_count": len(variables),
        })
    } else {
        i.logger.Warn("Variables validation failed", map[string]interface{}{
            "variable_count": len(variables),
            "errors":         len(result.Errors),
            "warnings":       len(result.Warnings),
        })
    }

    return result, nil
}
```

**Purpose**: Delegates variable validation to the variables manager.

### 11. Deep Variables Validation: `manager.ValidateVariables()`

**File**: `internal/variables/manager.go`
```go
func (m *Manager) ValidateVariables(ctx context.Context, variables []spookytypes.Variable) (*spookytypes.ValidationResult, error) {
    m.logger.Info("Validating variables", map[string]interface{}{
        "variable_count": len(variables),
    })

    // Use the validator to validate variables
    result, err := m.validator.ValidateVariables(ctx, variables)
    if err != nil {
        return nil, fmt.Errorf("failed to validate variables: %w", err)
    }

    if result.Valid {
        m.logger.Info("Variables validation passed", map[string]interface{}{
            "variable_count": len(variables),
        })
    } else {
        m.logger.Warn("Variables validation failed", map[string]interface{}{
            "variable_count": len(variables),
            "errors":         len(result.Errors),
            "warnings":       len(result.Warnings),
        })
    }

    return result, nil
}
```

**Purpose**: Delegates variable validation to the variables validator.

### 12. Variables Validator: `validator.ValidateVariables()`

**File**: `internal/variables/validator.go`
```go
func (v *Validator) ValidateVariables(ctx context.Context, variables []spookytypes.Variable) (*spookytypes.ValidationResult, error) {
    v.logger.Info("Validating variables", map[string]interface{}{
        "variable_count": len(variables),
    })

    result := &spookytypes.ValidationResult{
        Valid:       true,
        ValidatedAt: time.Now(),
        Errors:      []spookytypesschemas.SchemaError{},
        Warnings:    []spookytypesschemas.SchemaError{},
    }

    // Validate each variable individually
    for i := range variables {
        variableResult, err := v.ValidateVariable(ctx, &variables[i])
        if err != nil {
            return nil, fmt.Errorf("failed to validate variable %d: %w", i, err)
        }

        // Merge validation results
        if !variableResult.Valid {
            result.Valid = false
        }
        result.Errors = append(result.Errors, variableResult.Errors...)
        result.Warnings = append(result.Warnings, variableResult.Warnings...)
    }

    // Validate cross-variable dependencies
    if err := v.validateVariableDependencies(ctx, variables, result); err != nil {
        return nil, fmt.Errorf("failed to validate variable dependencies: %w", err)
    }

    // Validate variable resolution
    if err := v.validateVariableResolution(ctx, variables, result); err != nil {
        return nil, fmt.Errorf("failed to validate variable resolution: %w", err)
    }

    return result, nil
}
```

**Purpose**: 
- Validates each variable individually
- Validates cross-variable dependencies
- Validates variable resolution
- Aggregates validation results

### 13. Individual Variable Validation: `validator.ValidateVariable()`

**File**: `internal/variables/validator.go`
```go
func (v *Validator) ValidateVariable(ctx context.Context, variable *spookytypes.Variable) (*spookytypes.ValidationResult, error) {
    v.logger.Debug("Validating variable", map[string]interface{}{
        "variable_name": variable.Name,
    })

    result := &spookytypes.ValidationResult{
        Valid:       true,
        ValidatedAt: time.Now(),
        Errors:      []spookytypesschemas.SchemaError{},
        Warnings:    []spookytypesschemas.SchemaError{},
    }

    // Validate variable name
    if err := v.validateVariableName(variable.Name); err != nil {
        result.Valid = false
        result.Errors = append(result.Errors, spookytypesschemas.SchemaError{
            Code:     "invalid_variable_name",
            Message:  fmt.Sprintf("Invalid variable name: %s", err.Error()),
            Severity: "error",
        })
    }

    // Validate variable type
    if err := v.validateVariableType(variable.Type); err != nil {
        result.Valid = false
        result.Errors = append(result.Errors, spookytypesschemas.SchemaError{
            Code:     "invalid_variable_type",
            Message:  fmt.Sprintf("Invalid variable type: %s", err.Error()),
            Severity: "error",
        })
    }

    // Validate variable value
    if err := v.validateVariableValue(variable.Value, variable.Type); err != nil {
        result.Valid = false
        result.Errors = append(result.Errors, spookytypesschemas.SchemaError{
            Code:     "invalid_variable_value",
            Message:  fmt.Sprintf("Invalid variable value: %s", err.Error()),
            Severity: "error",
        })
    }

    // Validate variable description
    if variable.Description == "" {
        result.Warnings = append(result.Warnings, spookytypesschemas.SchemaError{
            Code:     "missing_variable_description",
            Message:  "Variable description is recommended",
            Severity: "warning",
        })
    }

    // Validate variable constraints
    if err := v.validateVariableConstraints(variable.Constraints, variable.Type); err != nil {
        result.Valid = false
        result.Errors = append(result.Errors, spookytypesschemas.SchemaError{
            Code:     "invalid_variable_constraints",
            Message:  fmt.Sprintf("Invalid variable constraints: %s", err.Error()),
            Severity: "error",
        })
    }

    // Validate variable dependencies
    if err := v.validateVariableDependencies(variable.Dependencies); err != nil {
        result.Valid = false
        result.Errors = append(result.Errors, spookytypesschemas.SchemaError{
            Code:     "invalid_variable_dependencies",
            Message:  fmt.Sprintf("Invalid variable dependencies: %s", err.Error()),
            Severity: "error",
        })
    }

    // Validate encrypted variables
    if variable.Encrypted {
        if err := v.validateEncryptedVariable(variable); err != nil {
            result.Valid = false
            result.Errors = append(result.Errors, spookytypesschemas.SchemaError{
                Code:     "invalid_encrypted_variable",
                Message:  fmt.Sprintf("Invalid encrypted variable: %s", err.Error()),
                Severity: "error",
            })
        }
    }

    return result, nil
}
```

**Purpose**: 
- Validates variable name format
- Validates variable type
- Validates variable value
- Validates variable description
- Validates variable constraints
- Validates variable dependencies
- Validates encrypted variables

## Key Components

### Variables Integration
- **Purpose**: Provides variables functionality to the CLI
- **Responsibilities**: Loading, validation, resolution
- **Interface**: `spookyinterfaces.VariablesIntegration`

### Variables Manager
- **Purpose**: Coordinates variables operations
- **Responsibilities**: Loading, validation, resolution
- **Interface**: `spookyinterfaces.VariablesManager`

### Variables Validator
- **Purpose**: Validates variable configurations
- **Responsibilities**: Syntax validation, dependency validation, constraint validation
- **Features**: Individual and collection validation

### Variables Loader
- **Purpose**: Loads variables from disk
- **Responsibilities**: HCL parsing, file loading, directory scanning
- **Features**: Multiple source support (file and directory)

## Validation Process

1. **Project Validation**
   - Validates project path exists
   - Ensures project.hcl is present
   - Checks project structure

2. **Variables Loading**
   - Loads from variables.hcl file
   - Loads from variables/ directory
   - Merges variables from multiple sources

3. **Individual Variable Validation**
   - Validates variable name format
   - Validates variable type
   - Validates variable value
   - Validates variable description
   - Validates variable constraints
   - Validates variable dependencies
   - Validates encrypted variables

4. **Collection Validation**
   - Validates cross-variable dependencies
   - Validates variable relationships
   - Checks for consistency across variables

5. **Resolution Validation**
   - Tests variable resolution
   - Validates dependency resolution
   - Checks for circular dependencies

## Error Handling

- **Graceful Loading**: Handles missing variable files gracefully
- **Detailed Error Reporting**: Provides specific error messages with context
- **Warning System**: Reports non-critical issues as warnings
- **Collection Validation**: Validates entire variable collection

## Output Format

### Success Case
```
🔍 Validating variables in project: /path/to/project
📊 Validating 10 variables...

✅ Validation Results:
──────────────────────────────────────────────────
🎉 All variables are valid!
✅ Validation completed successfully
```

### Error Case
```
🔍 Validating variables in project: /path/to/project
📊 Validating 10 variables...

✅ Validation Results:
──────────────────────────────────────────────────
❌ Errors (2):
  1. Invalid variable name: db-password (contains invalid characters)
     Context: variable[0]
  2. Invalid variable value: not a valid integer
     Context: variable[1]

⚠️  Warnings (1):
  1. Variable description is recommended
     Context: variable[2]

❌ Validation failed with 2 errors
```

## Architecture Patterns

- **Interface-Based Design**: Uses interfaces for loose coupling
- **Dependency Injection**: Injects dependencies through constructors
- **Layered Architecture**: Separates concerns across multiple layers
- **Error Aggregation**: Collects and reports all validation issues
- **Multiple Source Loading**: Supports both file and directory-based variable definitions
- **Dependency Resolution**: Validates variable dependencies and resolution
