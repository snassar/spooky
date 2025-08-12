# Variables System API Reference

## Overview

This document provides a comprehensive API reference for the spooky variables system. It covers all interfaces, types, methods, and implementation details for developers working with the variables system.

## Table of Contents

1. [Core Interfaces](#core-interfaces)
2. [Type Definitions](#type-definitions)
3. [Implementation Details](#implementation-details)
4. [Error Handling](#error-handling)
5. [Validation Rules](#validation-rules)
6. [CLI Integration](#cli-integration)
7. [Examples](#examples)

## Core Interfaces

### VariablesIntegration

The primary interface for variable management operations.

```go
type VariablesIntegration interface {
    // LoadVariables loads variables from the given source
    LoadVariables(ctx context.Context, source string) (map[string]*spookytypes.Variable, error)
    
    // ValidateVariables validates variables
    ValidateVariables(ctx context.Context, variables map[string]*spookytypes.Variable) (*spookytypes.ValidationResult, error)
    
    // ResolveVariables resolves variables with the given context
    ResolveVariables(ctx context.Context, variables map[string]*spookytypes.Variable, context *spookytypes.VariableContext) (*spookytypes.VariableResolutionResult, error)
}
```

**Methods:**

#### LoadVariables
```go
LoadVariables(ctx context.Context, source string) (map[string]*spookytypes.Variable, error)
```
Loads variables from a project directory, supporting both `variables.hcl` files and `variables/` directories.

**Parameters:**
- `ctx`: Context for cancellation and timeouts
- `source`: Project directory path

**Returns:**
- `map[string]*spookytypes.Variable`: Map of loaded variables by name
- `error`: Error if loading fails

**Behavior:**
- Checks for `variables.hcl` file in the project root
- Checks for `variables/` directory and loads all `.hcl` files
- Combines variables from all sources
- Performs duplicate detection and cross-file validation
- Returns error if no variables found or validation fails

#### ValidateVariables
```go
ValidateVariables(ctx context.Context, variables map[string]*spookytypes.Variable) (*spookytypes.ValidationResult, error)
```
Validates variable configurations and returns detailed validation results.

**Parameters:**
- `ctx`: Context for cancellation and timeouts
- `variables`: Map of variables to validate

**Returns:**
- `*spookytypes.ValidationResult`: Validation results with errors and warnings
- `error`: Error if validation process fails

**Validation Rules:**
- Required fields validation (name, type, scope)
- Variable name format validation
- Type consistency validation
- Dependency validation and circular dependency detection
- Constraint validation
- Cross-file duplicate detection

#### ResolveVariables
```go
ResolveVariables(ctx context.Context, variables map[string]*spookytypes.Variable, context *spookytypes.VariableContext) (*spookytypes.VariableResolutionResult, error)
```
Resolves variables with dependencies and context, applying environment overrides.

**Parameters:**
- `ctx`: Context for cancellation and timeouts
- `variables`: Map of variables to resolve
- `context`: Variable resolution context

**Returns:**
- `*spookytypes.VariableResolutionResult`: Resolution results with resolved values
- `error`: Error if resolution fails

**Resolution Process:**
1. Dependency analysis and topological sorting
2. Environment variable override application
3. Default value resolution
4. Type conversion and validation
5. Circular dependency detection

### VariableValidator

Interface for variable validation operations.

```go
type VariableValidator interface {
    // ValidateVariables validates a collection of variables
    ValidateVariables(ctx context.Context, variables map[string]*spookytypes.Variable) (*spookytypes.ValidationResult, error)
    
    // ValidateVariable validates a single variable
    ValidateVariable(ctx context.Context, variable *spookytypes.Variable) (*spookytypes.ValidationResult, error)
}
```

**Methods:**

#### ValidateVariables
```go
ValidateVariables(ctx context.Context, variables map[string]*spookytypes.Variable) (*spookytypes.ValidationResult, error)
```
Validates a collection of variables with comprehensive validation rules.

**Parameters:**
- `ctx`: Context for cancellation and timeouts
- `variables`: Map of variables to validate

**Returns:**
- `*spookytypes.ValidationResult`: Validation results with errors and warnings
- `error`: Error if validation process fails

**Validation Steps:**
1. Individual variable validation
2. Cross-variable validation (dependencies, duplicates)
3. Circular dependency detection
4. Constraint validation
5. Type consistency validation

#### ValidateVariable
```go
ValidateVariable(ctx context.Context, variable *spookytypes.Variable) (*spookytypes.ValidationResult, error)
```
Validates a single variable against all validation rules.

**Parameters:**
- `ctx`: Context for cancellation and timeouts
- `variable`: Variable to validate

**Returns:**
- `*spookytypes.ValidationResult`: Validation results with errors and warnings
- `error`: Error if validation process fails

**Validation Rules:**
- Name format validation
- Type validation
- Scope validation
- Constraint validation
- Validation rule compliance

### VariableLoader

Interface for variable loading operations.

```go
type VariableLoader interface {
    // LoadVariablesFromFile loads variables from a file
    LoadVariablesFromFile(ctx context.Context, filePath string) (map[string]*spookytypes.Variable, error)
    
    // LoadVariablesFromDirectory loads variables from a directory
    LoadVariablesFromDirectory(ctx context.Context, dirPath string) (map[string]*spookytypes.Variable, error)
}
```

**Methods:**

#### LoadVariablesFromFile
```go
LoadVariablesFromFile(ctx context.Context, filePath string) (map[string]*spookytypes.Variable, error)
```
Loads variables from a single HCL file.

**Parameters:**
- `ctx`: Context for cancellation and timeouts
- `filePath`: Path to the HCL file

**Returns:**
- `map[string]*spookytypes.Variable`: Map of loaded variables by name
- `error`: Error if loading fails

**Loading Process:**
1. HCL file parsing
2. Variable block extraction
3. Attribute parsing and type conversion
4. Validation block processing
5. Constraint block processing
6. Source file tracking

#### LoadVariablesFromDirectory
```go
LoadVariablesFromDirectory(ctx context.Context, dirPath string) (map[string]*spookytypes.Variable, error)
```
Loads variables from all HCL files in a directory.

**Parameters:**
- `ctx`: Context for cancellation and timeouts
- `dirPath`: Path to the directory

**Returns:**
- `map[string]*spookytypes.Variable`: Map of loaded variables by name
- `error`: Error if loading fails

**Loading Process:**
1. Directory traversal
2. HCL file discovery
3. Individual file loading
4. Variable merging
5. Duplicate detection
6. Cross-file validation

## Type Definitions

### Core Variable Types

#### Variable
```go
type Variable struct {
    spookytypescommon.CompleteEntity

    // Core variable properties
    Name        string      `hcl:"name" json:"name"`
    Type        VariableType `hcl:"type" json:"type"`
    Description string      `hcl:"description,optional" json:"description,omitempty"`
    Default     interface{} `hcl:"default,optional" json:"default,omitempty"`
    Required    bool        `hcl:"required,optional" json:"required,omitempty"`
    Sensitive   bool        `hcl:"sensitive,optional" json:"sensitive,omitempty"`
    Encrypted   bool        `hcl:"encrypted,optional" json:"encrypted,omitempty"`
    Scope       VariableScope `hcl:"scope,optional" json:"scope,omitempty"`

    // Dependencies and validation
    Dependencies []string           `hcl:"dependencies,optional" json:"dependencies,omitempty"`
    Validation   *VariableValidation `hcl:"validation,block" json:"validation,omitempty"`
    Constraints  *VariableConstraints `hcl:"constraints,block" json:"constraints,omitempty"`

    // Metadata and source information
    SourceFile string                 `hcl:"-" json:"source_file,omitempty"`
    SourceLine int                    `hcl:"-" json:"source_line,omitempty"`
    Metadata   map[string]interface{} `hcl:"metadata,optional" json:"metadata,omitempty"`

    // Runtime state
    ResolvedValue interface{} `hcl:"-" json:"resolved_value,omitempty"`
    IsResolved    bool        `hcl:"-" json:"is_resolved,omitempty"`
    ResolutionError string    `hcl:"-" json:"resolution_error,omitempty"`
}
```

#### VariableType
```go
type VariableType string

const (
    VariableTypeString VariableType = "string"
    VariableTypeNumber VariableType = "number"
    VariableTypeBool   VariableType = "bool"
    VariableTypeList   VariableType = "list"
    VariableTypeMap    VariableType = "map"
)
```

#### VariableScope
```go
type VariableScope string

const (
    VariableScopeProject VariableScope = "project"
    VariableScopeGlobal  VariableScope = "global"
    VariableScopeLocal   VariableScope = "local"
)
```

### Validation Types

#### VariableValidation
```go
type VariableValidation struct {
    AllowedValues []interface{} `hcl:"allowed_values,optional" json:"allowed_values,omitempty"`
    Pattern       *string       `hcl:"pattern,optional" json:"pattern,omitempty"`
    MinValue      *float64      `hcl:"min_value,optional" json:"min_value,omitempty"`
    MaxValue      *float64      `hcl:"max_value,optional" json:"max_value,omitempty"`
    MinLength     *int          `hcl:"min_length,optional" json:"min_length,omitempty"`
    MaxLength     *int          `hcl:"max_length,optional" json:"max_length,omitempty"`
    MinItems      *int          `hcl:"min_items,optional" json:"min_items,omitempty"`
    MaxItems      *int          `hcl:"max_items,optional" json:"max_items,omitempty"`
}
```

#### VariableConstraints
```go
type VariableConstraints struct {
    MinValue  *float64 `hcl:"min_value,optional" json:"min_value,omitempty"`
    MaxValue  *float64 `hcl:"max_value,optional" json:"max_value,omitempty"`
    MinLength *int     `hcl:"min_length,optional" json:"min_length,omitempty"`
    MaxLength *int     `hcl:"max_length,optional" json:"max_length,omitempty"`
    Pattern   *string  `hcl:"pattern,optional" json:"pattern,omitempty"`
    MinItems  *int     `hcl:"min_items,optional" json:"min_items,omitempty"`
    MaxItems  *int     `hcl:"max_items,optional" json:"max_items,omitempty"`
}
```

### Context and Resolution Types

#### VariableContext
```go
type VariableContext struct {
    spookytypescommon.CompleteEntity

    // Context information
    ProjectPath string                 `json:"project_path"`
    Environment string                 `json:"environment,omitempty"`
    Overrides   map[string]interface{} `json:"overrides,omitempty"`
    Metadata    map[string]interface{} `json:"metadata,omitempty"`

    // Resolution settings
    ApplyEnvironmentOverrides bool `json:"apply_environment_overrides"`
    StrictMode                bool `json:"strict_mode"`
    ValidateConstraints       bool `json:"validate_constraints"`
}
```

#### VariableResolutionResult
```go
type VariableResolutionResult struct {
    spookytypescommon.CompleteEntity

    // Resolution results
    ResolvedVariables map[string]*Variable `json:"resolved_variables"`
    UnresolvedVariables []string           `json:"unresolved_variables,omitempty"`
    ResolutionOrder    []string            `json:"resolution_order,omitempty"`

    // Statistics
    TotalVariables     int `json:"total_variables"`
    ResolvedCount      int `json:"resolved_count"`
    UnresolvedCount    int `json:"unresolved_count"`
    OverrideCount      int `json:"override_count"`

    // Errors and warnings
    Errors   []VariableError   `json:"errors,omitempty"`
    Warnings []VariableWarning `json:"warnings,omitempty"`
}
```

### Error and Warning Types

#### VariableError
```go
type VariableError struct {
    spookytypescommon.ErrorDetails

    Type        VariableErrorType `json:"type"`
    VariableName string           `json:"variable_name,omitempty"`
    Field       string            `json:"field,omitempty"`
    Value       interface{}       `json:"value,omitempty"`
    Expected    interface{}       `json:"expected,omitempty"`
}
```

#### VariableErrorType
```go
type VariableErrorType string

const (
    VariableErrorTypeValidation      VariableErrorType = "validation"
    VariableErrorTypeResolution      VariableErrorType = "resolution"
    VariableErrorTypeDependency      VariableErrorType = "dependency"
    VariableErrorTypeType            VariableErrorType = "type"
    VariableErrorTypeConstraint      VariableErrorType = "constraint"
    VariableErrorTypeCircular        VariableErrorType = "circular"
    VariableErrorTypeDuplicate       VariableErrorType = "duplicate"
    VariableErrorTypeRequired        VariableErrorType = "required"
    VariableErrorTypeEnvironment     VariableErrorType = "environment"
)
```

#### VariableWarning
```go
type VariableWarning struct {
    spookytypescommon.ErrorDetails

    Type        VariableWarningType `json:"type"`
    VariableName string             `json:"variable_name,omitempty"`
    Field       string              `json:"field,omitempty"`
    Value       interface{}         `json:"value,omitempty"`
    Suggestion  string              `json:"suggestion,omitempty"`
}
```

#### VariableWarningType
```go
type VariableWarningType string

const (
    VariableWarningTypeDeprecated  VariableWarningType = "deprecated"
    VariableWarningTypeUnused      VariableWarningType = "unused"
    VariableWarningTypeDefault     VariableWarningType = "default"
    VariableWarningTypeOverride    VariableWarningType = "override"
    VariableWarningTypeConstraint  VariableWarningType = "constraint"
)
```

## Implementation Details

### LoadVariables Implementation

The `LoadVariables` method in the `VariablesIntegration` interface orchestrates the loading process:

```go
func (m *Manager) LoadVariables(ctx context.Context, projectPath string) (map[string]*spookytypesvariables.Variable, error) {
    m.logger.Info("Loading variables", "project_path", projectPath)
    
    variables := make(map[string]*spookytypesvariables.Variable)
    
    // Load from variables.hcl if it exists
    variablesHCLPath := filepath.Join(projectPath, "variables.hcl")
    if _, err := os.Stat(variablesHCLPath); err == nil {
        fileVariables, err := m.loader.LoadVariablesFromFile(ctx, variablesHCLPath)
        if err != nil {
            return nil, fmt.Errorf("failed to load variables.hcl: %w", err)
        }
        
        // Merge variables
        for name, variable := range fileVariables {
            if existing, exists := variables[name]; exists {
                return nil, fmt.Errorf("duplicate variable name '%s' found in variables.hcl and %s", name, existing.SourceFile)
            }
            variables[name] = variable
        }
        
        m.logger.Debug("Loaded variables from variables.hcl", "count", len(fileVariables))
    }
    
    // Load from variables/ directory if it exists
    variablesDirPath := filepath.Join(projectPath, "variables")
    if _, err := os.Stat(variablesDirPath); err == nil {
        dirVariables, err := m.loader.LoadVariablesFromDirectory(ctx, variablesDirPath)
        if err != nil {
            return nil, fmt.Errorf("failed to load variables directory: %w", err)
        }
        
        // Merge variables with duplicate detection
        for name, variable := range dirVariables {
            if existing, exists := variables[name]; exists {
                return nil, fmt.Errorf("duplicate variable name '%s' found in %s and %s", name, variable.SourceFile, existing.SourceFile)
            }
            variables[name] = variable
        }
        
        m.logger.Debug("Loaded variables from variables/ directory", "count", len(dirVariables))
    }
    
    if len(variables) == 0 {
        return nil, fmt.Errorf("no variables found in project")
    }
    
    m.logger.Info("Variables loaded successfully", "total_count", len(variables))
    return variables, nil
}
```

### ValidateVariables Implementation

The `ValidateVariables` method performs comprehensive validation:

```go
func (m *Manager) ValidateVariables(ctx context.Context, variables map[string]*spookytypesvariables.Variable) (*spookytypes.ValidationResult, error) {
    m.logger.Info("Validating variables", "count", len(variables))
    
    // Use the validator to perform validation
    result, err := m.validator.ValidateVariables(ctx, variables)
    if err != nil {
        return nil, fmt.Errorf("validation failed: %w", err)
    }
    
    m.logger.Info("Variable validation completed", 
        "total", len(variables),
        "errors", len(result.Errors),
        "warnings", len(result.Warnings))
    
    return result, nil
}
```

### ResolveVariables Implementation

The `ResolveVariables` method handles dependency resolution and environment overrides:

```go
func (m *Manager) ResolveVariables(ctx context.Context, variables map[string]*spookytypesvariables.Variable, context *spookytypesvariables.VariableContext) (*spookytypesvariables.VariableResolutionResult, error) {
    m.logger.Info("Resolving variables", "count", len(variables))
    
    // Create a copy of variables for resolution
    resolvedVariables := make(map[string]*spookytypesvariables.Variable)
    for name, variable := range variables {
        // Create a copy to avoid modifying the original
        resolvedVar := *variable
        resolvedVariables[name] = &resolvedVar
    }
    
    // Apply environment overrides if enabled
    if context.ApplyEnvironmentOverrides {
        m.applyEnvironmentOverrides(resolvedVariables)
    }
    
    // Resolve dependencies using topological sort
    resolutionOrder, err := m.resolveDependencies(resolvedVariables)
    if err != nil {
        return nil, fmt.Errorf("dependency resolution failed: %w", err)
    }
    
    // Resolve variables in dependency order
    for _, varName := range resolutionOrder {
        variable := resolvedVariables[varName]
        if err := m.resolveVariable(variable, resolvedVariables, context); err != nil {
            return nil, fmt.Errorf("failed to resolve variable '%s': %w", varName, err)
        }
    }
    
    // Build result
    result := &spookytypesvariables.VariableResolutionResult{
        ResolvedVariables: resolvedVariables,
        ResolutionOrder:   resolutionOrder,
        TotalVariables:    len(variables),
        ResolvedCount:     len(resolutionOrder),
    }
    
    m.logger.Info("Variable resolution completed", 
        "total", result.TotalVariables,
        "resolved", result.ResolvedCount)
    
    return result, nil
}
```

## Error Handling

### Structured Error Types

The variables system uses structured error types for comprehensive error reporting:

```go
// VariableError provides detailed error information
type VariableError struct {
    spookytypescommon.ErrorDetails

    Type        VariableErrorType `json:"type"`
    VariableName string           `json:"variable_name,omitempty"`
    Field       string            `json:"field,omitempty"`
    Value       interface{}       `json:"value,omitempty"`
    Expected    interface{}       `json:"expected,omitempty"`
}

// NewVariableError creates a new variable error
func NewVariableError(errorType VariableErrorType, message string, variableName string) *VariableError {
    return &VariableError{
        ErrorDetails: spookytypescommon.ErrorDetails{
            Message:   message,
            Timestamp: time.Now(),
        },
        Type:         errorType,
        VariableName: variableName,
    }
}
```

### Error Handling Patterns

#### Validation Error Handling
```go
func (v *Validator) validateVariableName(name string) error {
    if name == "" {
        return fmt.Errorf("variable name cannot be empty")
    }
    
    // Check name format
    nameRegex := regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)
    if !nameRegex.MatchString(name) {
        return fmt.Errorf("invalid variable name '%s': must match pattern ^[a-zA-Z_][a-zA-Z0-9_]*$", name)
    }
    
    return nil
}
```

#### Resolution Error Handling
```go
func (m *Manager) resolveVariable(variable *spookytypesvariables.Variable, allVariables map[string]*spookytypesvariables.Variable, context *spookytypesvariables.VariableContext) error {
    // Check if already resolved
    if variable.IsResolved {
        return nil
    }
    
    // Resolve dependencies first
    for _, depName := range variable.Dependencies {
        depVar, exists := allVariables[depName]
        if !exists {
            return fmt.Errorf("variable '%s' depends on undefined variable '%s'", variable.Name, depName)
        }
        
        if !depVar.IsResolved {
            return fmt.Errorf("circular dependency detected: variable '%s' depends on unresolved variable '%s'", variable.Name, depName)
        }
    }
    
    // Apply default value if not set
    if variable.ResolvedValue == nil && variable.Default != nil {
        variable.ResolvedValue = variable.Default
    }
    
    // Validate constraints if enabled
    if context.ValidateConstraints && variable.Constraints != nil {
        if err := m.validateVariableConstraints(variable); err != nil {
            return fmt.Errorf("constraint validation failed for variable '%s': %w", variable.Name, err)
        }
    }
    
    variable.IsResolved = true
    return nil
}
```

## Validation Rules

### Variable Name Validation

Variable names must follow specific patterns:

```go
// Variable name validation rules
const (
    // Name must start with letter or underscore
    // Can contain letters, numbers, and underscores
    // Must be at least 1 character long
    VariableNamePattern = `^[a-zA-Z_][a-zA-Z0-9_]*$`
    
    // Reserved names that cannot be used
    ReservedVariableNames = []string{
        "type", "name", "description", "default", "required",
        "sensitive", "encrypted", "scope", "dependencies",
        "validation", "constraints", "metadata",
    }
)
```

### Type Validation

Each variable type has specific validation rules:

```go
func (v *Validator) validateVariableType(variable *spookytypesvariables.Variable) error {
    switch variable.Type {
    case spookytypesvariables.VariableTypeString:
        return v.validateStringVariable(variable)
    case spookytypesvariables.VariableTypeNumber:
        return v.validateNumberVariable(variable)
    case spookytypesvariables.VariableTypeBool:
        return v.validateBoolVariable(variable)
    case spookytypesvariables.VariableTypeList:
        return v.validateListVariable(variable)
    case spookytypesvariables.VariableTypeMap:
        return v.validateMapVariable(variable)
    default:
        return fmt.Errorf("unknown variable type: %s", variable.Type)
    }
}
```

### Dependency Validation

Dependencies are validated for circular references and existence:

```go
func (v *Validator) validateDependencies(variables map[string]*spookytypesvariables.Variable) error {
    // Build dependency graph
    graph := make(map[string][]string)
    for name, variable := range variables {
        graph[name] = variable.Dependencies
    }
    
    // Detect circular dependencies using DFS
    visited := make(map[string]bool)
    recStack := make(map[string]bool)
    
    for name := range graph {
        if !visited[name] {
            if v.hasCycle(name, graph, visited, recStack) {
                return fmt.Errorf("circular dependency detected")
            }
        }
    }
    
    return nil
}

func (v *Validator) hasCycle(node string, graph map[string][]string, visited, recStack map[string]bool) bool {
    visited[node] = true
    recStack[node] = true
    
    for _, neighbor := range graph[node] {
        if !visited[neighbor] {
            if v.hasCycle(neighbor, graph, visited, recStack) {
                return true
            }
        } else if recStack[neighbor] {
            return true
        }
    }
    
    recStack[node] = false
    return false
}
```

## CLI Integration

### Command Structure

The variables system integrates with the spooky CLI through the following commands:

```go
// variablesCmd represents the variables command
var variablesCmd = &cobra.Command{
    Use:   "variables",
    Short: "Manage project variables",
    Long: `Manage project variables including listing, validation, and resolution.
    
Variables can be defined in variables.hcl files or in a variables/ directory.
The system supports dependency resolution, environment overrides, and comprehensive validation.`,
}

// variablesListCmd represents the variables list command
var variablesListCmd = &cobra.Command{
    Use:   "list [project-path]",
    Short: "List variables in a project",
    Long: `List all variables in a project, showing their types, values, and sources.
    
Variables are grouped by source file and show their current values,
types, scopes, and whether they are sensitive or encrypted.`,
    Args: cobra.MaximumNArgs(1),
    RunE: func(cmd *cobra.Command, args []string) error {
        projectPath := "."
        if len(args) > 0 {
            projectPath = args[0]
        }
        return handleVariablesList(projectPath)
    },
}

// variablesValidateCmd represents the variables validate command
var variablesValidateCmd = &cobra.Command{
    Use:   "validate [project-path]",
    Short: "Validate variable definitions",
    Long: `Validate variable definitions for syntax, dependencies, and constraints.
    
Performs comprehensive validation including:
- Syntax validation
- Type checking
- Dependency validation
- Circular dependency detection
- Constraint validation
- Cross-file duplicate detection`,
    Args: cobra.MaximumNArgs(1),
    RunE: func(cmd *cobra.Command, args []string) error {
        projectPath := "."
        if len(args) > 0 {
            projectPath = args[0]
        }
        return handleVariablesValidate(projectPath)
    },
}

// variablesResolveCmd represents the variables resolve command
var variablesResolveCmd = &cobra.Command{
    Use:   "resolve [project-path]",
    Short: "Resolve variables with context",
    Long: `Resolve variables with dependencies and context.
    
Resolves variable dependencies, applies environment overrides,
and provides the final resolved values for all variables.`,
    Args: cobra.MaximumNArgs(1),
    RunE: func(cmd *cobra.Command, args []string) error {
        projectPath := "."
        if len(args) > 0 {
            projectPath = args[0]
        }
        return handleVariablesResolve(projectPath)
    },
}
```

### Command Handlers

#### List Handler
```go
func handleVariablesList(projectPath string) error {
    // Initialize dependencies if not already done
    if variablesManager == nil {
        if err := InitializeVariablesDependencies(); err != nil {
            return fmt.Errorf("failed to initialize variables dependencies: %w", err)
        }
    }
    
    // Load variables
    variables, err := variablesManager.LoadVariables(context.Background(), projectPath)
    if err != nil {
        return fmt.Errorf("failed to load variables: %w", err)
    }
    
    // Group variables by source file
    sourceGroups := make(map[string][]*spookytypesvariables.Variable)
    for _, variable := range variables {
        sourceFile := variable.SourceFile
        if sourceFile == "" {
            sourceFile = "unknown"
        }
        sourceGroups[sourceFile] = append(sourceGroups[sourceFile], variable)
    }
    
    // Display variables grouped by source
    fmt.Printf("Variables in %s:\n\n", projectPath)
    
    var totalCount int
    for sourceFile, vars := range sourceGroups {
        fmt.Printf("📁 %s:\n", filepath.Base(sourceFile))
        for _, variable := range vars {
            value := formatVariableValue(variable)
            fmt.Printf("  • %s (%s) = %s [%s]\n", 
                variable.Name, variable.Type, value, variable.Scope)
        }
        fmt.Println()
        totalCount += len(vars)
    }
    
    fmt.Printf("Total: %d variables across %d files\n", totalCount, len(sourceGroups))
    return nil
}
```

#### Validate Handler
```go
func handleVariablesValidate(projectPath string) error {
    // Initialize dependencies if not already done
    if variablesManager == nil {
        if err := InitializeVariablesDependencies(); err != nil {
            return fmt.Errorf("failed to initialize variables dependencies: %w", err)
        }
    }
    
    // Load variables
    variables, err := variablesManager.LoadVariables(context.Background(), projectPath)
    if err != nil {
        return fmt.Errorf("failed to load variables: %w", err)
    }
    
    // Validate variables
    result, err := variablesManager.ValidateVariables(context.Background(), variables)
    if err != nil {
        return fmt.Errorf("validation failed: %w", err)
    }
    
    // Display results
    if len(result.Errors) == 0 {
        fmt.Println("✅ Validation successful")
        fmt.Println()
        fmt.Printf("Variables validated: %d\n", len(variables))
        fmt.Printf("Files processed: %d\n", countSourceFiles(variables))
        fmt.Printf("Warnings: %d\n", len(result.Warnings))
        fmt.Printf("Errors: %d\n", len(result.Errors))
        
        if len(result.Warnings) > 0 {
            fmt.Println("\nWarnings:")
            for _, warning := range result.Warnings {
                fmt.Printf("  • %s\n", warning.Message)
            }
        }
    } else {
        fmt.Println("❌ Validation failed")
        fmt.Println()
        fmt.Printf("Variables validated: %d\n", len(variables))
        fmt.Printf("Errors: %d\n", len(result.Errors))
        fmt.Printf("Warnings: %d\n", len(result.Warnings))
        fmt.Println()
        
        fmt.Println("Errors:")
        for _, error := range result.Errors {
            fmt.Printf("  • %s\n", error.Message)
        }
        
        return fmt.Errorf("validation failed with %d errors", len(result.Errors))
    }
    
    return nil
}
```

#### Resolve Handler
```go
func handleVariablesResolve(projectPath string) error {
    // Initialize dependencies if not already done
    if variablesManager == nil {
        if err := InitializeVariablesDependencies(); err != nil {
            return fmt.Errorf("failed to initialize variables dependencies: %w", err)
        }
    }
    
    // Load variables
    variables, err := variablesManager.LoadVariables(context.Background(), projectPath)
    if err != nil {
        return fmt.Errorf("failed to load variables: %w", err)
    }
    
    // Create resolution context
    context := &spookytypesvariables.VariableContext{
        ProjectPath:              projectPath,
        ApplyEnvironmentOverrides: true,
        StrictMode:               false,
        ValidateConstraints:      true,
    }
    
    // Resolve variables
    result, err := variablesManager.ResolveVariables(context.Background(), variables, context)
    if err != nil {
        return fmt.Errorf("resolution failed: %w", err)
    }
    
    // Display resolved variables
    fmt.Println("Resolved Variables:")
    fmt.Println()
    
    for _, varName := range result.ResolutionOrder {
        variable := result.ResolvedVariables[varName]
        value := formatVariableValue(variable)
        fmt.Printf("%s = %s\n", variable.Name, value)
    }
    
    fmt.Println()
    fmt.Println("Resolution completed successfully")
    fmt.Printf("Dependencies resolved: %d circular dependencies detected\n", len(result.Errors))
    fmt.Printf("Environment overrides applied: %d variables\n", result.OverrideCount)
    
    return nil
}
```

## Examples

### Basic Variable Loading

```go
// Load variables from a project
variables, err := variablesManager.LoadVariables(ctx, "./my-project")
if err != nil {
    log.Fatalf("Failed to load variables: %v", err)
}

// Display loaded variables
for name, variable := range variables {
    fmt.Printf("Variable: %s (Type: %s, Scope: %s)\n", 
        name, variable.Type, variable.Scope)
}
```

### Variable Validation

```go
// Validate variables
result, err := variablesManager.ValidateVariables(ctx, variables)
if err != nil {
    log.Fatalf("Validation failed: %v", err)
}

// Check validation results
if len(result.Errors) > 0 {
    fmt.Println("Validation errors:")
    for _, error := range result.Errors {
        fmt.Printf("  • %s\n", error.Message)
    }
} else {
    fmt.Println("✅ All variables are valid")
}
```

### Variable Resolution

```go
// Create resolution context
context := &spookytypesvariables.VariableContext{
    ProjectPath:              "./my-project",
    Environment:              "production",
    ApplyEnvironmentOverrides: true,
    StrictMode:               false,
}

// Resolve variables
result, err := variablesManager.ResolveVariables(ctx, variables, context)
if err != nil {
    log.Fatalf("Resolution failed: %v", err)
}

// Use resolved variables
for name, variable := range result.ResolvedVariables {
    if variable.IsResolved {
        fmt.Printf("%s = %v\n", name, variable.ResolvedValue)
    }
}
```

### Custom Validation

```go
// Create custom validator
validator := variables.NewValidator(logger)

// Validate individual variable
variable := &spookytypesvariables.Variable{
    Name: "app_port",
    Type: spookytypesvariables.VariableTypeNumber,
    Default: 8080,
    Scope: spookytypesvariables.VariableScopeProject,
    Constraints: &spookytypesvariables.VariableConstraints{
        MinValue: float64Ptr(1024),
        MaxValue: float64Ptr(65535),
    },
}

result, err := validator.ValidateVariable(ctx, variable)
if err != nil {
    log.Fatalf("Validation failed: %v", err)
}

if len(result.Errors) == 0 {
    fmt.Println("Variable is valid")
}
```

This API reference provides comprehensive documentation for the spooky variables system. For usage examples and best practices, refer to the [User Guide](VARIABLES_USER_GUIDE.md), and for troubleshooting help, see the [Troubleshooting Guide](VARIABLES_TROUBLESHOOTING.md).
