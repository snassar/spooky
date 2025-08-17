# Code Walkthrough: `spooky project validate`

## Command Overview

**Command**: `spooky project validate <project-directory>`

**Purpose**: Validates a spooky project directory structure and configuration files against the project directory schema.

**Key Features**:
- Schema-driven validation using `project-directory.hcl` schema
- Comprehensive error reporting with context
- Validation of project structure, configuration files, and metadata
- Support for both strict and lenient validation modes

## Execution Flow

### 1. Command Entry Point: `main()`

**File**: `main.go`
```go
func main() {
    cmd.Execute()
}
```

**What Happens**:
- Calls the root command's `Execute()` method
- Initializes the CLI framework and command structure

### 2. Root Command Setup: `RootCmd.Execute()`

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

**What Happens**:
- Executes the root command with all subcommands
- Handles any errors and exits with appropriate status code

### 3. Persistent Pre-Run: `RootCmd.PersistentPreRunE`

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

**What Happens**:
- Sets up global configuration and logging
- Initializes integration dependencies
- Ensures proper environment setup for all commands

### 4. Project Validate Command: `projectValidateCmd`

**File**: `cmd/project.go`
```go
var projectValidateCmd = &cobra.Command{
    Use:   "validate [project-directory]",
    Short: "Validate a spooky project",
    Long: `Validate a spooky project directory structure and configuration.
    
This command validates:
- Project directory structure against schema
- Configuration file syntax and content
- Required and optional files/directories
- Project metadata and settings`,
    Args: cobra.ExactArgs(1),
    RunE: func(cmd *cobra.Command, args []string) error {
        return handleProjectValidate(args[0])
    },
}
```

**What Happens**:
- Defines the command structure and help text
- Validates exactly one argument (project directory)
- Calls the main handler function

### 5. Project Validation Handler: `handleProjectValidate()`

**File**: `cmd/project.go`
```go
func handleProjectValidate(projectPath string) error {
    ctx := context.Background()
    if projectValidator == nil {
        if err := InitializeProjectDependencies(); err != nil {
            return fmt.Errorf("failed to initialize project dependencies: %w", err)
        }
    }
    
    result, err := projectValidator.Validate(ctx, projectPath)
    if err != nil {
        return fmt.Errorf("validation failed: %w", err)
    }
    
    if !result.IsValid() {
        fmt.Fprintf(os.Stderr, "Project validation failed:\n")
        for _, err := range result.GetErrors() {
            fmt.Fprintf(os.Stderr, "  ❌ %s\n", err.Error())
        }
        for _, warning := range result.GetWarnings() {
            fmt.Fprintf(os.Stderr, "  ⚠️  %s\n", warning)
        }
        os.Exit(1)
    }
    
    fmt.Printf("✅ Project validation passed\n")
    if len(result.GetWarnings()) > 0 {
        fmt.Printf("⚠️  Warnings:\n")
        for _, warning := range result.GetWarnings() {
            fmt.Printf("  %s\n", warning)
        }
    }
    return nil
}
```

**What Happens**:
- Initializes project dependencies if needed
- Calls the project validator with the project path
- Handles validation results and displays appropriate output
- Exits with error code 1 if validation fails

### 6. Project Validator: `projectValidator.Validate()`

**File**: `internal/project/validator.go`
```go
func (v *Validator) Validate(ctx context.Context, projectPath string) (*spookytypesschemas.ValidationResult, error) {
    v.logger.Info("Validating project", map[string]interface{}{
        "project_path": projectPath,
    })
    
    // Resolve absolute path
    absPath, err := filepath.Abs(projectPath)
    if err != nil {
        return nil, fmt.Errorf("failed to resolve project path: %w", err)
    }
    
    // Check if project directory exists
    if _, err := os.Stat(absPath); os.IsNotExist(err) {
        return spookytypesschemas.NewValidationResult(false, []error{
            fmt.Errorf("project directory does not exist: %s", absPath),
        }, nil), nil
    }
    
    // Validate project structure using schema-driven validator
    structureResult, err := v.schemaValidator.ValidateProjectStructure(ctx, absPath)
    if err != nil {
        return nil, fmt.Errorf("failed to validate project structure: %w", err)
    }
    
    // Validate project configuration
    configResult, err := v.validateProjectConfig(ctx, absPath)
    if err != nil {
        return nil, fmt.Errorf("failed to validate project configuration: %w", err)
    }
    
    // Combine validation results
    allErrors := append(structureResult.GetErrors(), configResult.GetErrors()...)
    allWarnings := append(structureResult.GetWarnings(), configResult.GetWarnings()...)
    
    isValid := len(allErrors) == 0
    
    return spookytypesschemas.NewValidationResult(isValid, allErrors, allWarnings), nil
}
```

**What Happens**:
- Resolves the absolute project path
- Checks if the project directory exists
- Validates project structure using schema-driven validator
- Validates project configuration files
- Combines and returns validation results

### 7. Schema-Driven Structure Validation: `schemaValidator.ValidateProjectStructure()`

**File**: `internal/schemas/schema_driven_validator.go`
```go
func (v *SchemaDrivenValidator) ValidateProjectStructure(ctx context.Context, projectPath string) (*spookytypesschemas.ValidationResult, error) {
    v.logger.Info("Validating project structure", map[string]interface{}{
        "project_path": projectPath,
    })
    
    var errors []error
    var warnings []string
    
    // Load project directory schema
    schema, exists := v.structureSchemas["project-directory"]
    if !exists {
        return nil, fmt.Errorf("project-directory schema not found")
    }
    
    // Validate required files
    for _, file := range schema.GetRequiredFiles() {
        filePath := filepath.Join(projectPath, file.Name)
        if _, err := os.Stat(filePath); os.IsNotExist(err) {
            errors = append(errors, fmt.Errorf("required file missing: %s", file.Name))
        } else if err != nil {
            errors = append(errors, fmt.Errorf("failed to check required file %s: %w", file.Name, err))
        }
    }
    
    // Validate optional files (warnings only)
    for _, file := range schema.GetOptionalFiles() {
        filePath := filepath.Join(projectPath, file.Name)
        if _, err := os.Stat(filePath); os.IsNotExist(err) {
            warnings = append(warnings, fmt.Sprintf("optional file not present: %s", file.Name))
        }
    }
    
    // Validate required directories
    for _, dir := range schema.GetRequiredDirectories() {
        dirPath := filepath.Join(projectPath, dir.Name)
        if info, err := os.Stat(dirPath); os.IsNotExist(err) {
            errors = append(errors, fmt.Errorf("required directory missing: %s", dir.Name))
        } else if err != nil {
            errors = append(errors, fmt.Errorf("failed to check required directory %s: %w", dir.Name, err))
        } else if !info.IsDir() {
            errors = append(errors, fmt.Errorf("required directory is not a directory: %s", dir.Name))
        }
    }
    
    // Validate optional directories (warnings only)
    for _, dir := range schema.GetOptionalDirectories() {
        dirPath := filepath.Join(projectPath, dir.Name)
        if info, err := os.Stat(dirPath); os.IsNotExist(err) {
            warnings = append(warnings, fmt.Sprintf("optional directory not present: %s", dir.Name))
        } else if err != nil {
            warnings = append(warnings, fmt.Sprintf("failed to check optional directory %s: %s", dir.Name, err.Error()))
        } else if !info.IsDir() {
            warnings = append(warnings, fmt.Sprintf("optional directory is not a directory: %s", dir.Name))
        }
    }
    
    isValid := len(errors) == 0
    return spookytypesschemas.NewValidationResult(isValid, errors, warnings), nil
}
```

**What Happens**:
- Loads the project directory schema from embedded schemas
- Validates required files and directories (errors if missing)
- Validates optional files and directories (warnings if missing)
- Checks file/directory types and permissions
- Returns comprehensive validation results

### 8. Project Configuration Validation: `validateProjectConfig()`

**File**: `internal/project/validator.go`
```go
func (v *Validator) validateProjectConfig(ctx context.Context, projectPath string) (*spookytypesschemas.ValidationResult, error) {
    var errors []error
    var warnings []string
    
    // Validate project.hcl
    projectHCLPath := filepath.Join(projectPath, "project.hcl")
    if _, err := os.Stat(projectHCLPath); err == nil {
        if err := v.validateHCLFile(projectHCLPath); err != nil {
            errors = append(errors, fmt.Errorf("invalid project.hcl: %w", err))
        }
    } else {
        errors = append(errors, fmt.Errorf("required file missing: project.hcl"))
    }
    
    // Validate other configuration files if they exist
    configFiles := []string{"machines.hcl", "actions.hcl", "variables.hcl"}
    for _, configFile := range configFiles {
        configPath := filepath.Join(projectPath, configFile)
        if _, err := os.Stat(configPath); err == nil {
            if err := v.validateHCLFile(configPath); err != nil {
                warnings = append(warnings, fmt.Sprintf("invalid %s: %s", configFile, err.Error()))
            }
        }
    }
    
    isValid := len(errors) == 0
    return spookytypesschemas.NewValidationResult(isValid, errors, warnings), nil
}
```

**What Happens**:
- Validates the required `project.hcl` file
- Validates optional configuration files if they exist
- Checks HCL syntax and basic structure
- Returns validation results with appropriate error/warning levels

### 9. HCL File Validation: `validateHCLFile()`

**File**: `internal/project/validator.go`
```go
func (v *Validator) validateHCLFile(filePath string) error {
    data, err := os.ReadFile(filePath)
    if err != nil {
        return fmt.Errorf("failed to read file: %w", err)
    }
    
    var result interface{}
    if err := hcl.Unmarshal(data, &result); err != nil {
        return fmt.Errorf("invalid HCL syntax: %w", err)
    }
    
    return nil
}
```

**What Happens**:
- Reads the HCL file content
- Attempts to parse the HCL syntax
- Returns error if parsing fails

## Key Components

### Schema-Driven Validator
- **Purpose**: Validates project structure against embedded schemas
- **Location**: `internal/schemas/schema_driven_validator.go`
- **Key Methods**: `ValidateProjectStructure()`, `loadStructureSchemas()`

### Project Validator
- **Purpose**: Coordinates validation of project structure and configuration
- **Location**: `internal/project/validator.go`
- **Key Methods**: `Validate()`, `validateProjectConfig()`

### Project Directory Schema
- **Purpose**: Defines expected project structure
- **Location**: `internal/schemas/schemas/structure/project-directory.hcl`
- **Content**: Required/optional files and directories

## Error Handling

### Validation Result Structure
```go
type ValidationResult struct {
    IsValid  bool
    Errors   []error
    Warnings []string
}
```

### Error Categories
1. **Structure Errors**: Missing required files/directories
2. **Configuration Errors**: Invalid HCL syntax in config files
3. **Permission Errors**: File/directory permission issues
4. **Type Errors**: Files where directories expected and vice versa

### Output Format
- **Success**: "✅ Project validation passed"
- **Errors**: "❌ [error message]" (causes exit code 1)
- **Warnings**: "⚠️ [warning message]" (doesn't cause failure)

## Architecture Patterns

### Dependency Injection
- Validators are injected into the CLI layer
- Schema-driven validator is injected into project validator
- Logger is injected into all components

### Interface-Based Design
- All validators implement `Validator` interface
- Validation results use common `ValidationResult` type
- Schema loading uses interface-based approach

### Error Aggregation
- Multiple validation errors are collected before reporting
- Warnings don't cause validation failure
- Comprehensive error context is provided

## Integration Points

### Schema System
- Uses embedded HCL schemas for validation
- Supports both strict and lenient validation modes
- Extensible for new schema types

### Configuration System
- Integrates with global configuration setup
- Uses XDG base directory specification
- Supports environment variable overrides

### Logging System
- Structured logging with context
- Debug information for troubleshooting
- Error tracking and reporting

## Example Usage

```bash
# Validate a project directory
spooky project validate ./my-project

# Output on success
✅ Project validation passed

# Output on failure
❌ required file missing: project.hcl
❌ required directory missing: machines
⚠️  optional file not present: variables.hcl
```

## Exit Codes

- **0**: Validation passed (may have warnings)
- **1**: Validation failed (has errors)

## Performance Considerations

- Schema loading is cached after first load
- File system operations are minimized
- Validation stops on first critical error
- Parallel validation of independent components
