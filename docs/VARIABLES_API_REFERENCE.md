# Variables System API Reference

## Overview

This document provides a comprehensive API reference for the spooky variables system. It covers all interfaces, types, methods, and implementation details for developers working with the variables system.

**Status: Partially Implemented** - The variables system has basic functionality but SSH-based variable collection has known issues that need to be addressed.

## Core Interfaces

### VariablesIntegration Interface

The `VariablesIntegration` interface provides the primary entry point for variables operations:

```go
type VariablesIntegration interface {
    // LoadVariables loads variables from the specified project path
    LoadVariables(ctx context.Context, projectPath string) (map[string]interface{}, error)
    
    // ValidateVariables validates variable definitions
    ValidateVariables(ctx context.Context, variables map[string]interface{}) (*ValidationResult, error)
    
    // ResolveVariables resolves variable dependencies and values
    ResolveVariables(ctx context.Context, variables map[string]interface{}) (map[string]interface{}, error)
}
```

**Implementation Status**: ⚠️ **Partially Implemented** - Basic functionality exists but SSH-based collection has issues

### VariableManager Interface

The `VariableManager` interface provides variable management and resolution:

```go
type VariableManager interface {
    // LoadVariables loads variables from project configuration
    LoadVariables(ctx context.Context, projectPath string) (map[string]*spookytypesvariables.Variable, error)
    
    // ValidateVariables validates variable definitions
    ValidateVariables(ctx context.Context, variables map[string]*spookytypesvariables.Variable) (*ValidationResult, error)
    
    // ResolveVariables resolves variable dependencies and values
    ResolveVariables(ctx context.Context, variables map[string]*spookytypesvariables.Variable) (map[string]interface{}, error)
}
```

**Implementation Status**: ⚠️ **Partially Implemented** - Basic loading and validation exist but resolution has issues

## Current Implementation Status

### ✅ Working Components

1. **Variable Loading**: Loading variables from HCL configuration files
2. **Variable Validation**: Basic validation of variable definitions
3. **Variable Structure**: Proper variable type definitions and structures
4. **CLI Integration**: `spooky variables list` command with filtering options
5. **Project Integration**: Variables loading from project configuration
6. **Basic Validation**: Variable definition validation and error handling
7. **Filtering Support**: Support for variable name and type filtering
8. **Export Support**: Variable export to JSON and HCL formats

### ⚠️ Known Issues

1. **SSH-Based Collection**: SSH-based variable collection has implementation issues
2. **Variable Resolution**: Variable dependency resolution has problems
3. **Remote Variable Reading**: Cannot properly read variables from remote machines
4. **Parallel Processing**: No parallel variable collection support
5. **Import Functionality**: No variable import capabilities
6. **Template Integration**: No template variable integration

### 🔄 In Progress

1. **SSH Collection Fixes**: Addressing SSH-based variable collection issues
2. **Resolution Improvements**: Implementing proper variable resolution
3. **Collection Enhancements**: Improving variable collection reliability

## Implementation Details

### Variable Loading System

The variables system loads variables from HCL configuration files:

```go
type VariableLoader struct {
    logger spookylogging.Logger
}

func (l *VariableLoader) LoadVariables(ctx context.Context, projectPath string) (map[string]*spookytypesvariables.Variable, error) {
    variables := make(map[string]*spookytypesvariables.Variable)
    
    // Load variables.hcl file
    variablesPath := filepath.Join(projectPath, "variables.hcl")
    if data, err := os.ReadFile(variablesPath); err == nil {
        if err := l.parseVariablesFile(data, variables); err != nil {
            return nil, fmt.Errorf("failed to parse variables.hcl: %w", err)
        }
    }
    
    // Load variables from variables/ directory
    variablesDir := filepath.Join(projectPath, "variables")
    if entries, err := os.ReadDir(variablesDir); err == nil {
        for _, entry := range entries {
            if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".hcl") {
                filePath := filepath.Join(variablesDir, entry.Name())
                if data, err := os.ReadFile(filePath); err == nil {
                    if err := l.parseVariablesFile(data, variables); err != nil {
                        return nil, fmt.Errorf("failed to parse %s: %w", entry.Name(), err)
                    }
                }
            }
        }
    }
    
    return variables, nil
}

func (l *VariableLoader) parseVariablesFile(data []byte, variables map[string]*spookytypesvariables.Variable) error {
    var config struct {
        Variables []*spookytypesvariables.Variable `hcl:"variable,block"`
    }
    
    if err := hcl.Unmarshal(data, &config); err != nil {
        return fmt.Errorf("failed to parse HCL: %w", err)
    }
    
    for _, variable := range config.Variables {
        if variable.Name == "" {
            return fmt.Errorf("variable name is required")
        }
        
        if _, exists := variables[variable.Name]; exists {
            return fmt.Errorf("duplicate variable name: %s", variable.Name)
        }
        
        variables[variable.Name] = variable
    }
    
    return nil
}
```

**Supported Variable Sources:**
- **Local Variables**: Variables defined in `variables.hcl` and `variables/*.hcl` files
- **Environment Variables**: Variables from environment
- **SSH Variables**: Variables collected from remote machines (has issues)
- **Computed Variables**: Variables computed from other variables

### Variable Validation System

Variables are validated against schemas and business rules:

```go
type VariableValidator struct {
    logger spookylogging.Logger
}

func (v *VariableValidator) ValidateVariables(ctx context.Context, variables map[string]*spookytypesvariables.Variable) (*spookytypes.ValidationResult, error) {
    var errors []spookyschemas.SchemaError
    var warnings []spookyschemas.SchemaError
    
    for name, variable := range variables {
        // Validate variable name
        if name == "" {
            errors = append(errors, spookyschemas.SchemaError{
                Message: "variable name cannot be empty",
            })
            continue
        }
        
        // Validate variable structure
        if err := v.validateVariableStructure(variable); err != nil {
            errors = append(errors, spookyschemas.SchemaError{
                Message: fmt.Sprintf("variable %s: %s", name, err.Error()),
            })
        }
        
        // Validate variable type
        if err := v.validateVariableType(variable); err != nil {
            errors = append(errors, spookyschemas.SchemaError{
                Message: fmt.Sprintf("variable %s: %s", name, err.Error()),
            })
        }
        
        // Validate variable constraints
        if err := v.validateVariableConstraints(variable); err != nil {
            errors = append(errors, spookyschemas.SchemaError{
                Message: fmt.Sprintf("variable %s: %s", name, err.Error()),
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

### Variable Resolution System

Variables are resolved through dependency resolution (currently has issues):

```go
type VariableResolver struct {
    logger spookylogging.Logger
}

func (r *VariableResolver) ResolveVariables(ctx context.Context, variables map[string]*spookytypesvariables.Variable) (map[string]interface{}, error) {
    resolved := make(map[string]interface{})
    
    // Build dependency graph
    graph := r.buildDependencyGraph(variables)
    
    // Detect circular dependencies
    if r.hasCircularDependencies(graph) {
        return nil, fmt.Errorf("circular dependencies detected in variables")
    }
    
    // Resolve variables in dependency order
    resolvedOrder := r.topologicalSort(graph)
    
    for _, varName := range resolvedOrder {
        variable := variables[varName]
        
        // Resolve variable value
        value, err := r.resolveVariableValue(ctx, variable, resolved)
        if err != nil {
            return nil, fmt.Errorf("failed to resolve variable %s: %w", varName, err)
        }
        
        resolved[varName] = value
    }
    
    return resolved, nil
}

func (r *VariableResolver) resolveVariableValue(ctx context.Context, variable *spookytypesvariables.Variable, resolved map[string]interface{}) (interface{}, error) {
    switch variable.Type {
    case "string":
        return r.resolveStringVariable(variable, resolved)
    case "number":
        return r.resolveNumberVariable(variable, resolved)
    case "boolean":
        return r.resolveBooleanVariable(variable, resolved)
    case "list":
        return r.resolveListVariable(variable, resolved)
    case "map":
        return r.resolveMapVariable(variable, resolved)
    case "ssh":
        return r.resolveSSHVariable(ctx, variable)
    default:
        return nil, fmt.Errorf("unsupported variable type: %s", variable.Type)
    }
}
```

## Type Definitions

### Variable Types

```go
// Variable represents a variable definition
type Variable struct {
    // Variable name (required)
    Name string `json:"name" hcl:"name"`
    
    // Variable description (optional)
    Description string `json:"description,omitempty" hcl:"description,optional"`
    
    // Variable type (string, number, boolean, list, map, ssh)
    Type string `json:"type" hcl:"type"`
    
    // Variable value (for local variables)
    Value interface{} `json:"value,omitempty" hcl:"value,optional"`
    
    // Variable source (for remote variables)
    Source *VariableSource `json:"source,omitempty" hcl:"source,optional"`
    
    // Variable constraints
    Constraints *VariableConstraints `json:"constraints,omitempty" hcl:"constraints,optional"`
    
    // Variable dependencies
    Dependencies []string `json:"dependencies,omitempty" hcl:"dependencies,optional"`
    
    // Variable metadata
    Metadata map[string]interface{} `json:"metadata,omitempty" hcl:"metadata,optional"`
}

// VariableSource defines how to obtain variable values
type VariableSource struct {
    // Source type (environment, ssh, computed)
    Type string `json:"type" hcl:"type"`
    
    // Source configuration
    Config map[string]interface{} `json:"config,omitempty" hcl:"config,optional"`
    
    // SSH source configuration
    SSH *SSHVariableSource `json:"ssh,omitempty" hcl:"ssh,optional"`
    
    // Environment source configuration
    Environment *EnvironmentVariableSource `json:"environment,omitempty" hcl:"environment,optional"`
}

// SSHVariableSource defines SSH-based variable collection
type SSHVariableSource struct {
    // Target machine
    Machine string `json:"machine" hcl:"machine"`
    
    // SSH command to execute
    Command string `json:"command" hcl:"command"`
    
    // File path to read (alternative to command)
    File string `json:"file,omitempty" hcl:"file,optional"`
    
    // Parse configuration
    Parse *ParseConfig `json:"parse,omitempty" hcl:"parse,optional"`
}

// EnvironmentVariableSource defines environment variable collection
type EnvironmentVariableSource struct {
    // Environment variable name
    Name string `json:"name" hcl:"name"`
    
    // Default value if not set
    Default interface{} `json:"default,omitempty" hcl:"default,optional"`
}

// VariableConstraints defines variable validation constraints
type VariableConstraints struct {
    // Required constraint
    Required bool `json:"required,omitempty" hcl:"required,optional"`
    
    // String constraints
    MinLength int    `json:"min_length,omitempty" hcl:"min_length,optional"`
    MaxLength int    `json:"max_length,omitempty" hcl:"max_length,optional"`
    Pattern   string `json:"pattern,omitempty" hcl:"pattern,optional"`
    
    // Number constraints
    MinValue float64 `json:"min_value,omitempty" hcl:"min_value,optional"`
    MaxValue float64 `json:"max_value,omitempty" hcl:"max_value,optional"`
    
    // List constraints
    MinItems int `json:"min_items,omitempty" hcl:"min_items,optional"`
    MaxItems int `json:"max_items,omitempty" hcl:"max_items,optional"`
    
    // Allowed values
    AllowedValues []interface{} `json:"allowed_values,omitempty" hcl:"allowed_values,optional"`
}

// ParseConfig defines how to parse variable values
type ParseConfig struct {
    // Parse format (json, yaml, hcl, text)
    Format string `json:"format" hcl:"format"`
    
    // Parse path (for nested values)
    Path string `json:"path,omitempty" hcl:"path,optional"`
    
    // Parse options
    Options map[string]interface{} `json:"options,omitempty" hcl:"options,optional"`
}
```

### Variable Context Types

```go
// VariableContext provides context for variable resolution
type VariableContext struct {
    // Project path
    ProjectPath string `json:"project_path" hcl:"project_path"`
    
    // Variable being resolved
    Variable *Variable `json:"variable" hcl:"variable"`
    
    // Resolution timestamp
    Timestamp time.Time `json:"timestamp" hcl:"timestamp"`
    
    // Resolution metadata
    Metadata map[string]interface{} `json:"metadata,omitempty" hcl:"metadata,optional"`
}

// VariableResult represents the result of variable resolution
type VariableResult struct {
    // Variable context
    Context *VariableContext `json:"context" hcl:"context"`
    
    // Resolution success
    Success bool `json:"success" hcl:"success"`
    
    // Resolved value
    Value interface{} `json:"value,omitempty" hcl:"value,optional"`
    
    // Resolution error
    Error string `json:"error,omitempty" hcl:"error,optional"`
    
    // Resolution duration
    Duration time.Duration `json:"duration" hcl:"duration"`
}
```

## Error Handling

### Variable Errors

```go
// VariableError represents variable resolution errors
type VariableError struct {
    VariableName string `json:"variable_name" hcl:"variable_name"`
    Error        string `json:"error" hcl:"error"`
    Details      string `json:"details,omitempty" hcl:"details,optional"`
}

// VariableValidationError represents variable validation errors
type VariableValidationError struct {
    Field   string `json:"field" hcl:"field"`
    Message string `json:"message" hcl:"message"`
    Value   string `json:"value,omitempty" hcl:"value,optional"`
}
```

### Validation Implementation

```go
// ValidateVariable validates a single variable
func (v *VariableValidator) ValidateVariable(variable *spookytypesvariables.Variable) error {
    if variable == nil {
        return fmt.Errorf("variable cannot be nil")
    }
    
    // Validate required fields
    if variable.Name == "" {
        return fmt.Errorf("variable name is required")
    }
    
    if variable.Type == "" {
        return fmt.Errorf("variable type is required")
    }
    
    // Validate variable type
    validTypes := []string{"string", "number", "boolean", "list", "map", "ssh"}
    valid := false
    for _, t := range validTypes {
        if variable.Type == t {
            valid = true
            break
        }
    }
    if !valid {
        return fmt.Errorf("invalid variable type: %s (valid types: %v)", variable.Type, validTypes)
    }
    
    // Validate variable source
    if variable.Source != nil {
        if err := v.validateVariableSource(variable.Source); err != nil {
            return fmt.Errorf("invalid variable source: %w", err)
        }
    }
    
    // Validate variable constraints
    if variable.Constraints != nil {
        if err := v.validateVariableConstraints(variable.Constraints); err != nil {
            return fmt.Errorf("invalid variable constraints: %w", err)
        }
    }
    
    return nil
}
```

## CLI Commands

### Variables List Command

```bash
# List all variables in a project
spooky variables list ./my-project

# List variables with specific types
spooky variables list ./my-project --type string,number

# List variables with specific names
spooky variables list ./my-project --names app_version,db_host

# List variables with verbose output
spooky variables list ./my-project --verbose
```

### Variables Export Command

```bash
# Export variables to JSON format
spooky variables export ./my-project --format json --output variables.json

# Export variables to HCL format
spooky variables export ./my-project --format hcl --output variables.hcl

# Export variables with specific types
spooky variables export ./my-project --type string,number --format json --output string-vars.json
```

### Variables Validation Command

```bash
# Validate variables in a project
spooky variables validate ./my-project

# Validate variables with verbose output
spooky variables validate ./my-project --verbose
```

## Integration Examples

### Basic Variable Definition

```hcl
# variables.hcl
variables {
  variable "app_version" {
    description = "Application version"
    type = "string"
    value = "1.0.0"
    
    constraints {
      required = true
      pattern = "^\\d+\\.\\d+\\.\\d+$"
    }
  }
  
  variable "db_host" {
    description = "Database host"
    type = "string"
    
    source {
      type = "ssh"
      ssh {
        machine = "db-server"
        command = "hostname -f"
      }
    }
    
    constraints {
      required = true
      min_length = 1
    }
  }
  
  variable "max_connections" {
    description = "Maximum database connections"
    type = "number"
    value = 100
    
    constraints {
      min_value = 1
      max_value = 1000
    }
  }
  
  variable "debug_mode" {
    description = "Enable debug mode"
    type = "boolean"
    
    source {
      type = "environment"
      environment {
        name = "DEBUG_MODE"
        default = false
      }
    }
  }
}
```

### Variable Loading and Validation

```go
// Variable loading and validation example
func loadAndValidateVariables(projectPath string) error {
    ctx := context.Background()
    
    // Create variable manager
    manager := spookyvariables.NewManager(loader, validator, logger)
    
    // Load variables
    variables, err := manager.LoadVariables(ctx, projectPath)
    if err != nil {
        return fmt.Errorf("failed to load variables: %w", err)
    }
    
    // Validate variables
    result, err := manager.ValidateVariables(ctx, variables)
    if err != nil {
        return fmt.Errorf("failed to validate variables: %w", err)
    }
    
    if !result.Valid {
        fmt.Println("Variable validation failed:")
        for _, error := range result.Errors {
            fmt.Printf("  - %s\n", error.Message)
        }
        return fmt.Errorf("variable validation failed")
    }
    
    fmt.Printf("Loaded and validated %d variables\n", len(variables))
    return nil
}
```

### Variable Resolution

```go
// Variable resolution example
func resolveVariables(projectPath string) error {
    ctx := context.Background()
    
    // Create variable manager
    manager := spookyvariables.NewManager(loader, validator, logger)
    
    // Load variables
    variables, err := manager.LoadVariables(ctx, projectPath)
    if err != nil {
        return fmt.Errorf("failed to load variables: %w", err)
    }
    
    // Resolve variables
    resolved, err := manager.ResolveVariables(ctx, variables)
    if err != nil {
        return fmt.Errorf("failed to resolve variables: %w", err)
    }
    
    // Print resolved variables
    for name, value := range resolved {
        fmt.Printf("%s = %v\n", name, value)
    }
    
    return nil
}
```

## Current Limitations

### Collection Limitations

1. **SSH Integration Issues**: SSH-based variable collection has known problems
2. **No Parallel Collection**: Variables are collected sequentially, not in parallel
3. **No Result Caching**: Variable values are not cached between operations
4. **No Incremental Collection**: Always collects all variables
5. **Limited Source Types**: Only basic source types are supported

### Resolution Limitations

1. **Dependency Issues**: Variable dependency resolution has problems
2. **No Circular Detection**: Circular dependencies may not be properly detected
3. **No Error Recovery**: No error recovery or retry mechanisms
4. **No Progress Tracking**: No progress tracking or status reporting

### Integration Limitations

1. **No Template Integration**: Variables are not integrated with template system
2. **No Action Integration**: Variables are not used in action system
3. **No Facts Integration**: Variables are not integrated with facts system
4. **No Conditional Resolution**: No conditional variable resolution

## Future Enhancements

### Planned Features

1. **SSH Collection Fixes**: Resolve SSH-based variable collection issues
2. **Parallel Collection**: Implement parallel variable collection
3. **Result Caching**: Add variable value caching
4. **Incremental Collection**: Support incremental variable collection
5. **Advanced Sources**: Support more variable source types
6. **Conditional Resolution**: Support conditional variable resolution

### Integration Enhancements

1. **Template Integration**: Integrate variables with template system
2. **Action Integration**: Use variables in action system
3. **Facts Integration**: Integrate variables with facts system
4. **Advanced Resolution**: Support advanced resolution features

## Summary

The variables system provides basic variable loading and validation capabilities but has significant limitations with SSH-based collection and dependency resolution that need to be addressed. The system is functional for basic use cases but requires improvements for production use.

**Status**: ⚠️ **Partially Implemented** - Basic functionality exists but SSH-based collection and dependency resolution have issues that need to be resolved.
