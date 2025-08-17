# Code Walkthrough: `spooky actions validate`

## Overview

This document provides a detailed function-by-function walkthrough of the `spooky actions validate` command execution flow. This command validates action configurations for syntax and dependencies.

## Command Structure

```bash
spooky actions validate <project-path>
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

### 4. Command Routing: `actionsValidateCmd.RunE`

**File**: `cmd/actions.go`
```go
var actionsValidateCmd = &cobra.Command{
    Use:   "validate",
    Short: "Validate action configurations",
    Long:  `Validate action configurations for syntax and dependencies.`,
    RunE: func(_ *cobra.Command, args []string) error {
        if len(args) == 0 {
            return fmt.Errorf("project directory is required")
        }

        projectPath := args[0]
        if err := validateProjectPath(projectPath); err != nil {
            return err
        }

        // Get actions integration
        actionsIntegration := getIntegrationManager().GetActionsIntegration()
        if actionsIntegration == nil {
            return fmt.Errorf("actions integration not available")
        }

        ctx := context.Background()

        // Load actions from project
        actions, err := actionsIntegration.LoadActions(ctx, projectPath)
        if err != nil {
            return fmt.Errorf("failed to load actions: %w", err)
        }

        if len(actions) == 0 {
            fmt.Println("No actions found to validate")
            return nil
        }

        // Validate actions
        result, err := actionsIntegration.ValidateActions(ctx, actions)
        if err != nil {
            return fmt.Errorf("validation failed: %w", err)
        }

        if result.Valid {
            fmt.Printf("✅ All %d actions are valid!\n", len(actions))
        } else {
            fmt.Printf("❌ Validation failed for %d actions:\n", len(result.Errors))
            for _, err := range result.Errors {
                fmt.Printf("  - %s\n", err.Message)
            }
            return fmt.Errorf("action validation failed")
        }

        return nil
    },
}
```

**Purpose**: 
- Validates project path
- Gets actions integration
- Loads and validates actions
- Reports results to the user

### 5. Project Path Validation: `validateProjectPath()`

**File**: `cmd/actions.go`
```go
func validateProjectPath(projectPath string) error {
    if projectPath == "" {
        return fmt.Errorf("project path cannot be empty")
    }

    // Check if project path exists
    if _, err := os.Stat(projectPath); os.IsNotExist(err) {
        return fmt.Errorf("project path does not exist: %s", projectPath)
    }

    // Check if it's a directory
    if info, err := os.Stat(projectPath); err == nil && !info.IsDir() {
        return fmt.Errorf("project path must be a directory: %s", projectPath)
    }

    // Check for project.hcl file
    projectFile := filepath.Join(projectPath, "project.hcl")
    if _, err := os.Stat(projectFile); os.IsNotExist(err) {
        return fmt.Errorf("project.hcl not found in: %s", projectPath)
    }

    return nil
}
```

**Purpose**: 
- Validates project path exists
- Ensures it's a directory
- Checks for required project.hcl file

### 6. Integration Manager Access: `getIntegrationManager()`

**File**: `cmd/actions.go`
```go
func getIntegrationManager() spookyinterfaces.IntegrationManager {
    // Return the global integration manager from cmd package
    return GetIntegrationManager()
}
```

**Purpose**: Returns the global integration manager for accessing system integrations.

### 7. Actions Integration Access: `GetActionsIntegration()`

**File**: `cmd/integrations.go`
```go
func (im *IntegrationManager) GetActionsIntegration() spookyinterfaces.ActionsIntegration {
    return im.actionsIntegration
}
```

**Purpose**: Returns the actions integration for loading and validating actions.

### 8. Actions Loading: `actionsIntegration.LoadActions()`

**File**: `internal/actions/integration.go`
```go
func (i *ActionsIntegration) LoadActions(ctx context.Context, projectPath string) ([]spookytypes.Action, error) {
    i.logger.Info("Loading actions from project", map[string]interface{}{
        "project_path": projectPath,
    })

    // Use the manager to load actions
    actions, err := i.manager.LoadActions(ctx, projectPath)
    if err != nil {
        return nil, fmt.Errorf("failed to load actions: %w", err)
    }

    i.logger.Info("Actions loaded successfully", map[string]interface{}{
        "project_path": projectPath,
        "action_count": len(actions),
    })

    return actions, nil
}
```

**Purpose**: Delegates action loading to the actions manager.

### 9. Deep Actions Loading: `manager.LoadActions()`

**File**: `internal/actions/manager.go`
```go
func (m *Manager) LoadActions(ctx context.Context, projectPath string) ([]spookytypes.Action, error) {
    m.logger.Info("Loading actions from project", map[string]interface{}{
        "project_path": projectPath,
    })

    // Use the loader to load actions
    actions, err := m.loader.LoadActions(ctx, projectPath)
    if err != nil {
        return nil, fmt.Errorf("failed to load actions: %w", err)
    }

    m.logger.Info("Actions loaded successfully", map[string]interface{}{
        "project_path": projectPath,
        "action_count": len(actions),
    })

    return actions, nil
}
```

**Purpose**: Delegates action loading to the actions loader.

### 10. Actions File Loading: `loader.LoadActions()`

**File**: `internal/actions/loader.go`
```go
func (l *Loader) LoadActions(ctx context.Context, projectPath string) ([]spookytypes.Action, error) {
    l.logger.Info("Loading actions from project", map[string]interface{}{
        "project_path": projectPath,
    })

    var allActions []spookytypes.Action
    var loadErrors []string

    // Check for actions.hcl file
    actionsFile := filepath.Join(projectPath, "actions.hcl")
    if _, err := os.Stat(actionsFile); err == nil {
        l.logger.Debug("Found actions.hcl file", map[string]interface{}{
            "file_path": actionsFile,
        })

        actions, err := l.loadActionsFromFile(ctx, actionsFile)
        if err != nil {
            loadErrors = append(loadErrors, fmt.Sprintf("actions.hcl: %v", err))
        } else {
            allActions = append(allActions, actions...)
            l.logger.Info("Loaded actions from file", map[string]interface{}{
                "file_path": actionsFile,
                "count":     len(actions),
            })
        }
    }

    // Check for actions/ directory
    actionsDir := filepath.Join(projectPath, "actions")
    if _, err := os.Stat(actionsDir); err == nil {
        l.logger.Debug("Found actions directory", map[string]interface{}{
            "dir_path": actionsDir,
        })

        actions, err := l.loadActionsFromDirectory(ctx, actionsDir)
        if err != nil {
            loadErrors = append(loadErrors, fmt.Sprintf("actions/ directory: %v", err))
        } else {
            allActions = append(allActions, actions...)
            l.logger.Info("Loaded actions from directory", map[string]interface{}{
                "dir_path": actionsDir,
                "count":    len(actions),
            })
        }
    }

    // Check if no actions were found
    if len(allActions) == 0 {
        if len(loadErrors) > 0 {
            return nil, fmt.Errorf("failed to load actions: %s", strings.Join(loadErrors, "; "))
        }
        return nil, fmt.Errorf("no actions found in project: %s (neither actions.hcl nor actions/ directory found)", projectPath)
    }

    // Validate for duplicates and consistency
    if err := l.validateActionCollection(allActions); err != nil {
        return nil, fmt.Errorf("action collection validation failed: %w", err)
    }

    l.logger.Info("Actions loaded successfully", map[string]interface{}{
        "project_path": projectPath,
        "total_count":  len(allActions),
    })

    return allActions, nil
}
```

**Purpose**: 
- Checks for actions.hcl file
- Checks for actions/ directory
- Loads actions from both sources
- Validates action collection
- Handles loading errors

### 11. Actions Validation: `actionsIntegration.ValidateActions()`

**File**: `internal/actions/integration.go`
```go
func (i *ActionsIntegration) ValidateActions(ctx context.Context, actions []spookytypes.Action) (*spookytypes.ValidationResult, error) {
    i.logger.Info("Validating actions", map[string]interface{}{
        "action_count": len(actions),
    })

    // Use the manager to validate actions
    result, err := i.manager.ValidateActions(ctx, actions)
    if err != nil {
        return nil, fmt.Errorf("failed to validate actions: %w", err)
    }

    if result.Valid {
        i.logger.Info("Actions validation passed", map[string]interface{}{
            "action_count": len(actions),
        })
    } else {
        i.logger.Warn("Actions validation failed", map[string]interface{}{
            "action_count": len(actions),
            "errors":       len(result.Errors),
            "warnings":     len(result.Warnings),
        })
    }

    return result, nil
}
```

**Purpose**: Delegates action validation to the actions manager.

### 12. Deep Actions Validation: `manager.ValidateActions()`

**File**: `internal/actions/manager.go`
```go
func (m *Manager) ValidateActions(ctx context.Context, actions []spookytypes.Action) (*spookytypes.ValidationResult, error) {
    m.logger.Info("Validating actions", map[string]interface{}{
        "action_count": len(actions),
    })

    // Use the validator to validate actions
    result, err := m.validator.ValidateActions(ctx, actions)
    if err != nil {
        return nil, fmt.Errorf("failed to validate actions: %w", err)
    }

    if result.Valid {
        m.logger.Info("Actions validation passed", map[string]interface{}{
            "action_count": len(actions),
        })
    } else {
        m.logger.Warn("Actions validation failed", map[string]interface{}{
            "action_count": len(actions),
            "errors":       len(result.Errors),
            "warnings":     len(result.Warnings),
        })
    }

    return result, nil
}
```

**Purpose**: Delegates action validation to the actions validator.

### 13. Actions Validator: `validator.ValidateActions()`

**File**: `internal/actions/validator.go`
```go
func (v *Validator) ValidateActions(ctx context.Context, actions []spookytypes.Action) (*spookytypes.ValidationResult, error) {
    v.logger.Info("Validating actions", map[string]interface{}{
        "action_count": len(actions),
    })

    result := &spookytypes.ValidationResult{
        Valid:       true,
        ValidatedAt: time.Now(),
        Errors:      []spookytypesschemas.SchemaError{},
        Warnings:    []spookytypesschemas.SchemaError{},
    }

    // Validate each action individually
    for i := range actions {
        actionResult, err := v.ValidateAction(ctx, &actions[i])
        if err != nil {
            return nil, fmt.Errorf("failed to validate action %d: %w", i, err)
        }

        // Merge validation results
        if !actionResult.Valid {
            result.Valid = false
        }
        result.Errors = append(result.Errors, actionResult.Errors...)
        result.Warnings = append(result.Warnings, actionResult.Warnings...)
    }

    // Validate cross-action dependencies
    if err := v.validateActionDependencies(ctx, actions, result); err != nil {
        return nil, fmt.Errorf("failed to validate action dependencies: %w", err)
    }

    // Validate action orchestration rules
    if err := v.validateOrchestrationRules(ctx, actions, result); err != nil {
        return nil, fmt.Errorf("failed to validate orchestration rules: %w", err)
    }

    return result, nil
}
```

**Purpose**: 
- Validates each action individually
- Validates cross-action dependencies
- Validates orchestration rules
- Aggregates validation results

### 14. Individual Action Validation: `validator.ValidateAction()`

**File**: `internal/actions/validator.go`
```go
func (v *Validator) ValidateAction(ctx context.Context, action *spookytypes.Action) (*spookytypes.ValidationResult, error) {
    v.logger.Debug("Validating action", map[string]interface{}{
        "action_name": action.Name,
    })

    result := &spookytypes.ValidationResult{
        Valid:       true,
        ValidatedAt: time.Now(),
        Errors:      []spookytypesschemas.SchemaError{},
        Warnings:    []spookytypesschemas.SchemaError{},
    }

    // Validate action name
    if err := v.validateActionName(action.Name); err != nil {
        result.Valid = false
        result.Errors = append(result.Errors, spookytypesschemas.SchemaError{
            Code:     "invalid_action_name",
            Message:  fmt.Sprintf("Invalid action name: %s", err.Error()),
            Severity: "error",
        })
    }

    // Validate action description
    if action.Description == "" {
        result.Warnings = append(result.Warnings, spookytypesschemas.SchemaError{
            Code:     "missing_action_description",
            Message:  "Action description is recommended",
            Severity: "warning",
        })
    }

    // Validate action machines
    if err := v.validateActionMachines(action.Machines); err != nil {
        result.Valid = false
        result.Errors = append(result.Errors, spookytypesschemas.SchemaError{
            Code:     "invalid_action_machines",
            Message:  fmt.Sprintf("Invalid action machines: %s", err.Error()),
            Severity: "error",
        })
    }

    // Validate action template
    if err := v.validateActionTemplate(action.Template); err != nil {
        result.Valid = false
        result.Errors = append(result.Errors, spookytypesschemas.SchemaError{
            Code:     "invalid_action_template",
            Message:  fmt.Sprintf("Invalid action template: %s", err.Error()),
            Severity: "error",
        })
    }

    // Validate action command
    if err := v.validateActionCommand(action.Command); err != nil {
        result.Valid = false
        result.Errors = append(result.Errors, spookytypesschemas.SchemaError{
            Code:     "invalid_action_command",
            Message:  fmt.Sprintf("Invalid action command: %s", err.Error()),
            Severity: "error",
        })
    }

    // Validate action variables
    if err := v.validateActionVariables(action.Variables); err != nil {
        result.Valid = false
        result.Errors = append(result.Errors, spookytypesschemas.SchemaError{
            Code:     "invalid_action_variables",
            Message:  fmt.Sprintf("Invalid action variables: %s", err.Error()),
            Severity: "error",
        })
    }

    return result, nil
}
```

**Purpose**: 
- Validates action name format
- Validates action description
- Validates target machines
- Validates template configuration
- Validates command configuration
- Validates variable definitions

## Key Components

### Actions Integration
- **Purpose**: Provides actions functionality to the CLI
- **Responsibilities**: Loading, validation, orchestration
- **Interface**: `spookyinterfaces.ActionsIntegration`

### Actions Manager
- **Purpose**: Coordinates actions operations
- **Responsibilities**: Loading, validation, orchestration
- **Interface**: `spookyinterfaces.ActionsManager`

### Actions Validator
- **Purpose**: Validates action configurations
- **Responsibilities**: Syntax validation, dependency validation, rule validation
- **Features**: Individual and collection validation

### Actions Loader
- **Purpose**: Loads actions from disk
- **Responsibilities**: HCL parsing, file loading, directory scanning
- **Features**: Multiple source support (file and directory)

## Validation Process

1. **Project Validation**
   - Validates project path exists
   - Ensures project.hcl is present
   - Checks project structure

2. **Actions Loading**
   - Loads from actions.hcl file
   - Loads from actions/ directory
   - Merges actions from multiple sources

3. **Individual Action Validation**
   - Validates action name format
   - Validates action description
   - Validates target machines
   - Validates template configuration
   - Validates command configuration
   - Validates variable definitions

4. **Collection Validation**
   - Validates cross-action dependencies
   - Validates orchestration rules
   - Checks for duplicate actions
   - Validates action relationships

## Error Handling

- **Graceful Loading**: Handles missing action files gracefully
- **Detailed Error Reporting**: Provides specific error messages with context
- **Warning System**: Reports non-critical issues as warnings
- **Collection Validation**: Validates entire action collection

## Output Format

### Success Case
```
✅ All 3 actions are valid!
```

### Error Case
```
❌ Validation failed for 2 actions:
  - Invalid action name: deploy-web (contains invalid characters)
  - Missing required field: command
```

## Architecture Patterns

- **Interface-Based Design**: Uses interfaces for loose coupling
- **Dependency Injection**: Injects dependencies through constructors
- **Layered Architecture**: Separates concerns across multiple layers
- **Error Aggregation**: Collects and reports all validation issues
- **Multiple Source Loading**: Supports both file and directory-based action definitions
