# Code Walkthrough: `spooky machines validate`

## Overview

This document provides a detailed function-by-function walkthrough of the `spooky machines validate` command execution flow. This command validates machine inventory configurations for syntax and connectivity.

## Command Structure

```bash
spooky machines validate <project-path>
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

### 4. Command Routing: `machinesValidateCmd.RunE`

**File**: `cmd/machines.go`
```go
var machinesValidateCmd = &cobra.Command{
    Use:   "validate [project-path]",
    Short: "Validate machine inventory",
    Long: `Validate machine inventory configuration and connectivity.

This command validates that all machines in the inventory have proper configuration
including required fields, valid SSH settings, and authentication methods.`,
    Args: cobra.ExactArgs(1),
    RunE: func(_ *cobra.Command, args []string) error {
        return handleMachinesValidate(args[0])
    },
}
```

**Purpose**: Routes the command to `handleMachinesValidate` with the project path.

### 5. Main Handler: `handleMachinesValidate()`

**File**: `cmd/machines.go`
```go
func handleMachinesValidate(projectPath string) error {
    ctx := context.Background()

    // Initialize dependencies if not already done
    if machinesManager == nil {
        if err := InitializeMachinesDependencies(); err != nil {
            return fmt.Errorf("failed to initialize machines dependencies: %w", err)
        }
    }

    fmt.Printf("🔍 Validating machines in project: %s\n", projectPath)

    // Load machines using the enhanced manager
    machines, err := machinesManager.LoadMachines(ctx, projectPath)
    if err != nil {
        return fmt.Errorf("failed to load machines: %w", err)
    }

    if len(machines) == 0 {
        fmt.Printf("No machines found in inventory.\n")
        return nil
    }

    fmt.Printf("📊 Validating %d machines...\n", len(machines))

    // Validate machines using the manager
    result, err := machinesManager.ValidateMachines(ctx, machines)
    if err != nil {
        return fmt.Errorf("failed to validate machines: %w", err)
    }

    // Display validation results
    fmt.Printf("\n✅ Validation Results:\n")
    fmt.Printf("%s\n", strings.Repeat("─", 50))

    if len(result.Errors) == 0 && len(result.Warnings) == 0 {
        fmt.Printf("🎉 All machines are valid!\n")
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
        return fmt.Errorf("machine validation failed")
    }

    fmt.Printf("✅ Validation completed successfully\n")
    return nil
}
```

**Purpose**: 
- Initializes machines dependencies
- Loads machines from project
- Validates machines
- Reports detailed results to the user

### 6. Dependency Initialization: `InitializeMachinesDependencies()`

**File**: `cmd/machines.go`
```go
func InitializeMachinesDependencies() error {
    // Create log manager for machines component
    logManager := spookylogging.NewLogManager()
    machinesLogger = logManager.GetLogger("machines")

    // Create SSH manager for connectivity testing
    sshManager := spookyssh.NewManager(machinesLogger)

    // Initialize machines components
    loader := spookymachines.NewLoader(machinesLogger)
    validator := spookymachines.NewValidator(machinesLogger)
    manager := spookymachines.NewManager(loader, validator, machinesLogger)

    // Create machines integration
    machinesManager = spookymachines.NewIntegration(manager)

    return nil
}
```

**Purpose**: Creates logger, SSH manager, loader, validator, and manager instances for dependency injection.

### 7. Machines Loading: `machinesManager.LoadMachines()`

**File**: `internal/machines/integration.go`
```go
func (i *MachinesIntegration) LoadMachines(ctx context.Context, projectPath string) ([]spookytypes.Machine, error) {
    i.logger.Info("Loading machines from project", map[string]interface{}{
        "project_path": projectPath,
    })

    // Use the manager to load machines
    machines, err := i.manager.LoadMachines(ctx, projectPath)
    if err != nil {
        return nil, fmt.Errorf("failed to load machines: %w", err)
    }

    i.logger.Info("Machines loaded successfully", map[string]interface{}{
        "project_path": projectPath,
        "machine_count": len(machines),
    })

    return machines, nil
}
```

**Purpose**: Delegates machine loading to the machines manager.

### 8. Deep Machines Loading: `manager.LoadMachines()`

**File**: `internal/machines/manager.go`
```go
func (m *Manager) LoadMachines(ctx context.Context, projectPath string) ([]spookytypes.Machine, error) {
    m.logger.Debug("Loading machines from project", map[string]interface{}{
        "project_path": projectPath,
    })

    var allMachines []spookytypes.Machine
    var loadErrors []string

    // Check for machines.hcl file
    machinesFile := filepath.Join(projectPath, "machines.hcl")
    if _, err := os.Stat(machinesFile); err == nil {
        m.logger.Debug("Found machines.hcl file", map[string]interface{}{
            "file_path": machinesFile,
        })

        machines, err := m.loader.LoadMachinesFromFile(ctx, machinesFile)
        if err != nil {
            loadErrors = append(loadErrors, fmt.Sprintf("machines.hcl: %v", err))
        } else {
            allMachines = append(allMachines, machines...)
            m.logger.Info("Loaded machines from file", map[string]interface{}{
                "file_path": machinesFile,
                "count":     len(machines),
            })
        }
    }

    // Check for machines/ directory
    machinesDir := filepath.Join(projectPath, "machines")
    if _, err := os.Stat(machinesDir); err == nil {
        m.logger.Debug("Found machines directory", map[string]interface{}{
            "dir_path": machinesDir,
        })

        machines, err := m.loader.LoadMachinesFromDirectory(ctx, machinesDir)
        if err != nil {
            loadErrors = append(loadErrors, fmt.Sprintf("machines/ directory: %v", err))
        } else {
            allMachines = append(allMachines, machines...)
            m.logger.Info("Loaded machines from directory", map[string]interface{}{
                "dir_path": machinesDir,
                "count":    len(machines),
            })
        }
    }

    // Check if no machines were found
    if len(allMachines) == 0 {
        if len(loadErrors) > 0 {
            return nil, fmt.Errorf("failed to load machines: %s", strings.Join(loadErrors, "; "))
        }
        return nil, fmt.Errorf("no machines found in project: %s (neither machines.hcl nor machines/ directory found)", projectPath)
    }

    // Validate for duplicates and consistency
    if err := m.validateMachineCollection(allMachines); err != nil {
        return nil, fmt.Errorf("machine collection validation failed: %w", err)
    }

    m.logger.Info("Machines loaded successfully", map[string]interface{}{
        "project_path": projectPath,
        "total_count":  len(allMachines),
    })

    return allMachines, nil
}
```

**Purpose**: 
- Checks for machines.hcl file
- Checks for machines/ directory
- Loads machines from both sources
- Validates machine collection
- Handles loading errors

### 9. Machines Validation: `machinesManager.ValidateMachines()`

**File**: `internal/machines/integration.go`
```go
func (i *MachinesIntegration) ValidateMachines(ctx context.Context, machines []spookytypes.Machine) (*spookytypes.ValidationResult, error) {
    i.logger.Info("Validating machines", map[string]interface{}{
        "machine_count": len(machines),
    })

    // Use the manager to validate machines
    result, err := i.manager.ValidateMachines(ctx, machines)
    if err != nil {
        return nil, fmt.Errorf("failed to validate machines: %w", err)
    }

    if result.Valid {
        i.logger.Info("Machines validation passed", map[string]interface{}{
            "machine_count": len(machines),
        })
    } else {
        i.logger.Warn("Machines validation failed", map[string]interface{}{
            "machine_count": len(machines),
            "errors":        len(result.Errors),
            "warnings":      len(result.Warnings),
        })
    }

    return result, nil
}
```

**Purpose**: Delegates machine validation to the machines manager.

### 10. Deep Machines Validation: `manager.ValidateMachines()`

**File**: `internal/machines/manager.go`
```go
func (m *Manager) ValidateMachines(ctx context.Context, machines []spookytypes.Machine) (*spookytypes.ValidationResult, error) {
    m.logger.Info("Validating machines", map[string]interface{}{
        "machine_count": len(machines),
    })

    // Use the validator to validate machines
    result, err := m.validator.ValidateMachines(ctx, machines)
    if err != nil {
        return nil, fmt.Errorf("failed to validate machines: %w", err)
    }

    if result.Valid {
        m.logger.Info("Machines validation passed", map[string]interface{}{
            "machine_count": len(machines),
        })
    } else {
        m.logger.Warn("Machines validation failed", map[string]interface{}{
            "machine_count": len(machines),
            "errors":        len(result.Errors),
            "warnings":      len(result.Warnings),
        })
    }

    return result, nil
}
```

**Purpose**: Delegates machine validation to the machines validator.

### 11. Machines Validator: `validator.ValidateMachines()`

**File**: `internal/machines/validator.go`
```go
func (v *Validator) ValidateMachines(ctx context.Context, machines []spookytypes.Machine) (*spookytypes.ValidationResult, error) {
    v.logger.Debug("Validating machines", map[string]interface{}{
        "count": len(machines),
    })

    var errors []spookytypesschemas.SchemaError
    var warnings []spookytypesschemas.SchemaError

    for i := range machines {
        result, err := v.ValidateMachine(ctx, &machines[i])
        if err != nil {
            return nil, fmt.Errorf("failed to validate machine %d: %w", i, err)
        }

        if !result.Valid {
            errors = append(errors, result.Errors...)
        }

        warnings = append(warnings, result.Warnings...)
    }

    valid := len(errors) == 0

    return &spookytypes.ValidationResult{
        Valid:    valid,
        Errors:   errors,
        Warnings: warnings,
    }, nil
}
```

**Purpose**: 
- Validates each machine individually
- Aggregates validation results
- Returns comprehensive validation status

### 12. Individual Machine Validation: `validator.ValidateMachine()`

**File**: `internal/machines/validator.go`
```go
func (v *Validator) ValidateMachine(ctx context.Context, machine *spookytypes.Machine) (*spookytypes.ValidationResult, error) {
    v.logger.Debug("Validating machine", map[string]interface{}{
        "hostname": machine.Hostname,
        "host":     machine.Host,
    })

    // Get machine schema for enhanced validation
    machineSchema, err := v.getMachineSchema()
    if err != nil {
        return nil, fmt.Errorf("failed to get machine schema: %w", err)
    }

    // Use enhanced validator for comprehensive machine validation
    result, err := v.enhancedValidator.ValidateWithEnhancedFeatures(ctx, machineSchema, machine)
    if err != nil {
        return nil, fmt.Errorf("failed to validate machine with enhanced validator: %w", err)
    }

    // Add additional custom validation for machine-specific rules
    v.addCustomMachineValidation(machine, result)

    return &spookytypes.ValidationResult{
        Valid:    result.Valid,
        Errors:   result.Errors,
        Warnings: result.Warnings,
    }, nil
}
```

**Purpose**: 
- Uses enhanced validator for schema validation
- Adds custom machine-specific validation rules
- Returns detailed validation results

### 13. Custom Machine Validation: `addCustomMachineValidation()`

**File**: `internal/machines/validator.go`
```go
func (v *Validator) addCustomMachineValidation(machine *spookytypes.Machine, result *spookytypes.ValidationResult) {
    // Validate hostname format
    if machine.Hostname != "" {
        if err := v.validateHostname(machine.Hostname); err != nil {
            result.Valid = false
            result.Errors = append(result.Errors, spookytypesschemas.SchemaError{
                Code:     "invalid_hostname",
                Message:  fmt.Sprintf("Invalid hostname format: %s", err.Error()),
                Severity: "error",
            })
        }
    }

    // Validate host format
    if machine.Host != "" {
        if err := v.validateHost(machine.Host); err != nil {
            result.Valid = false
            result.Errors = append(result.Errors, spookytypesschemas.SchemaError{
                Code:     "invalid_host",
                Message:  fmt.Sprintf("Invalid host format: %s", err.Error()),
                Severity: "error",
            })
        }
    }

    // Validate port range
    if machine.Port <= 0 || machine.Port > 65535 {
        result.Valid = false
        result.Errors = append(result.Errors, spookytypesschemas.SchemaError{
            Code:     "invalid_port",
            Message:  fmt.Sprintf("Port must be between 1 and 65535, got: %d", machine.Port),
            Severity: "error",
        })
    }

    // Validate user
    if machine.User == "" {
        result.Valid = false
        result.Errors = append(result.Errors, spookytypesschemas.SchemaError{
            Code:     "missing_user",
            Message:  "SSH user is required",
            Severity: "error",
        })
    }

    // Validate authentication method
    if err := v.validateAuthentication(machine.Authentication); err != nil {
        result.Valid = false
        result.Errors = append(result.Errors, spookytypesschemas.SchemaError{
            Code:     "invalid_authentication",
            Message:  fmt.Sprintf("Invalid authentication configuration: %s", err.Error()),
            Severity: "error",
        })
    }

    // Validate tags format
    if err := v.validateTags(machine.Tags); err != nil {
        result.Warnings = append(result.Warnings, spookytypesschemas.SchemaError{
            Code:     "invalid_tags",
            Message:  fmt.Sprintf("Invalid tags format: %s", err.Error()),
            Severity: "warning",
        })
    }
}
```

**Purpose**: 
- Validates hostname format
- Validates host format
- Validates port range
- Validates user configuration
- Validates authentication method
- Validates tags format

## Key Components

### Machines Integration
- **Purpose**: Provides machines functionality to the CLI
- **Responsibilities**: Loading, validation, connectivity testing
- **Interface**: `spookyinterfaces.MachinesIntegration`

### Machines Manager
- **Purpose**: Coordinates machines operations
- **Responsibilities**: Loading, validation, collection management
- **Interface**: `spookyinterfaces.MachinesManager`

### Machines Validator
- **Purpose**: Validates machine configurations
- **Responsibilities**: Schema validation, custom validation rules
- **Features**: Individual and collection validation

### Machines Loader
- **Purpose**: Loads machines from disk
- **Responsibilities**: HCL parsing, file loading, directory scanning
- **Features**: Multiple source support (file and directory)

## Validation Process

1. **Project Validation**
   - Validates project path exists
   - Ensures project.hcl is present
   - Checks project structure

2. **Machines Loading**
   - Loads from machines.hcl file
   - Loads from machines/ directory
   - Merges machines from multiple sources

3. **Individual Machine Validation**
   - Validates hostname format
   - Validates host format
   - Validates port range
   - Validates user configuration
   - Validates authentication method
   - Validates tags format

4. **Collection Validation**
   - Validates for duplicate machines
   - Validates machine relationships
   - Checks for consistency across machines

## Error Handling

- **Graceful Loading**: Handles missing machine files gracefully
- **Detailed Error Reporting**: Provides specific error messages with context
- **Warning System**: Reports non-critical issues as warnings
- **Collection Validation**: Validates entire machine collection

## Output Format

### Success Case
```
🔍 Validating machines in project: /path/to/project
📊 Validating 3 machines...

✅ Validation Results:
──────────────────────────────────────────────────
🎉 All machines are valid!
✅ Validation completed successfully
```

### Error Case
```
🔍 Validating machines in project: /path/to/project
📊 Validating 3 machines...

✅ Validation Results:
──────────────────────────────────────────────────
❌ Errors (2):
  1. Invalid hostname format: web-server-01 (contains invalid characters)
     Context: machine[0]
  2. Port must be between 1 and 65535, got: 70000
     Context: machine[1]

⚠️  Warnings (1):
  1. Invalid tags format: environment=production,role=web
     Context: machine[2]

❌ Validation failed with 2 errors
```

## Architecture Patterns

- **Interface-Based Design**: Uses interfaces for loose coupling
- **Dependency Injection**: Injects dependencies through constructors
- **Layered Architecture**: Separates concerns across multiple layers
- **Error Aggregation**: Collects and reports all validation issues
- **Multiple Source Loading**: Supports both file and directory-based machine definitions
- **Enhanced Validation**: Uses schema-driven validation with custom rules
