// Package variables provides variable management functionality for the spooky codebase.
package variables

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	spookyinterfaces "spooky/internal/interfaces"
	spookytypes "spooky/internal/types"
	spookytypeslogging "spooky/internal/types/logging"
	spookytypesvariables "spooky/internal/types/variables"
)

// Manager provides variable management functionality
type Manager struct {
	logger    spookytypeslogging.Logger
	loader    spookyinterfaces.VariableLoader
	validator spookyinterfaces.VariableValidator
}

// NewManager creates a new variable manager
func NewManager(
	logger spookytypeslogging.Logger,
	loader spookyinterfaces.VariableLoader,
	validator spookyinterfaces.VariableValidator,
) spookyinterfaces.VariablesIntegration {
	return &Manager{
		logger:    logger,
		loader:    loader,
		validator: validator,
	}
}

// LoadVariables loads variables from the given project path
// Supports both variables.hcl file and variables/ directory
func (m *Manager) LoadVariables(ctx context.Context, projectPath string) (map[string]*spookytypesvariables.Variable, error) {
	m.logger.Debug("Loading variables from project", map[string]interface{}{
		"project_path": projectPath,
	})

	var allVariables = make(map[string]*spookytypesvariables.Variable)
	var loadErrors []string

	// Check for variables.hcl file
	variablesFile := filepath.Join(projectPath, "variables.hcl")
	if _, err := os.Stat(variablesFile); err == nil {
		m.logger.Debug("Found variables.hcl file", map[string]interface{}{
			"file_path": variablesFile,
		})

		variables, err := m.loader.LoadVariablesFromFile(ctx, variablesFile)
		if err != nil {
			loadErrors = append(loadErrors, fmt.Sprintf("variables.hcl: %v", err))
		} else {
			// Merge variables from file
			for name, variable := range variables {
				allVariables[name] = variable
			}
			m.logger.Info("Loaded variables from file", map[string]interface{}{
				"file_path": variablesFile,
				"count":     len(variables),
			})
		}
	}

	// Check for variables/ directory
	variablesDir := filepath.Join(projectPath, "variables")
	if _, err := os.Stat(variablesDir); err == nil {
		m.logger.Debug("Found variables directory", map[string]interface{}{
			"dir_path": variablesDir,
		})

		variables, err := m.loader.LoadVariablesFromDirectory(ctx, variablesDir)
		if err != nil {
			loadErrors = append(loadErrors, fmt.Sprintf("variables/ directory: %v", err))
		} else {
			// Merge variables from directory
			for name, variable := range variables {
				allVariables[name] = variable
			}
			m.logger.Info("Loaded variables from directory", map[string]interface{}{
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
	if err := m.validateVariableCollection(allVariables); err != nil {
		return nil, fmt.Errorf("variable validation failed: %w", err)
	}

	m.logger.Info("Variables loaded successfully", map[string]interface{}{
		"project_path": projectPath,
		"total_count":  len(allVariables),
		"load_errors":  len(loadErrors),
	})

	return allVariables, nil
}

// ValidateVariables validates variables
func (m *Manager) ValidateVariables(ctx context.Context, variables map[string]*spookytypesvariables.Variable) (*spookytypes.ValidationResult, error) {
	m.logger.Debug("Validating variables", map[string]interface{}{
		"count": len(variables),
	})

	result, err := m.validator.ValidateVariables(ctx, variables)
	if err != nil {
		return nil, fmt.Errorf("failed to validate variables: %w", err)
	}

	m.logger.Info("Variable validation completed", map[string]interface{}{
		"count": len(variables),
		"valid": len(result.Errors) == 0,
	})

	return result, nil
}

// ResolveVariables resolves variables with the given context
func (m *Manager) ResolveVariables(_ context.Context, variables map[string]*spookytypesvariables.Variable, context *spookytypesvariables.VariableContext) (*spookytypesvariables.VariableResolutionResult, error) {
	m.logger.Info("Resolving variables", map[string]interface{}{
		"count": len(variables),
	})

	startTime := time.Now()
	resolved := make(map[string]*spookytypesvariables.Variable)
	var errors []spookytypesvariables.VariableError
	var warnings []spookytypesvariables.VariableWarning

	// Create a copy of variables for resolution
	resolutionVars := make(map[string]*spookytypesvariables.Variable)
	for name, variable := range variables {
		// Create a copy to avoid modifying the original
		varCopy := *variable
		resolutionVars[name] = &varCopy
	}

	// Resolve variables in dependency order
	resolvedVars, resolutionErrors, resolutionWarnings := m.resolveVariableDependencies(resolutionVars, context)

	// Merge results
	for name, variable := range resolvedVars {
		resolved[name] = variable
	}
	errors = append(errors, resolutionErrors...)
	warnings = append(warnings, resolutionWarnings...)

	// Apply environment variable overrides
	envErrors, envWarnings := m.applyEnvironmentOverrides(resolved, context)
	errors = append(errors, envErrors...)
	warnings = append(warnings, envWarnings...)

	// Validate resolved variables
	for name, variable := range resolved {
		if err := m.validateResolvedVariable(variable); err != nil {
			errors = append(errors, spookytypesvariables.VariableError{
				VariableName: name,
				ErrorType:    spookytypesvariables.VariableErrorTypeValidation,
				ErrorDetails: spookytypes.ErrorDetails{
					Message: err.Error(),
				},
			})
		}
	}

	duration := time.Since(startTime)

	m.logger.Info("Variable resolution completed", map[string]interface{}{
		"count":    len(variables),
		"resolved": len(resolved),
		"errors":   len(errors),
		"warnings": len(warnings),
		"duration": duration,
	})

	return &spookytypesvariables.VariableResolutionResult{
		Variables: resolved,
		Resolved:  m.extractResolvedValues(resolved),
		Errors:    errors,
		Warnings:  warnings,
		Duration:  duration,
	}, nil
}

// resolveVariableDependencies resolves variables in dependency order
func (m *Manager) resolveVariableDependencies(variables map[string]*spookytypesvariables.Variable, context *spookytypesvariables.VariableContext) (map[string]*spookytypesvariables.Variable, []spookytypesvariables.VariableError, []spookytypesvariables.VariableWarning) {
	var errors []spookytypesvariables.VariableError
	var warnings []spookytypesvariables.VariableWarning

	// Build dependency graph
	graph := m.buildDependencyGraph(variables)

	// Detect circular dependencies
	cycles := m.detectCircularDependencies(graph)
	if len(cycles) > 0 {
		for _, cycle := range cycles {
			errors = append(errors, spookytypesvariables.VariableError{
				ErrorType: spookytypesvariables.VariableErrorTypeCircular,
				ErrorDetails: spookytypes.ErrorDetails{
					Message: fmt.Sprintf("circular dependency detected: %s", strings.Join(cycle, " -> ")),
				},
			})
		}
		return variables, errors, warnings
	}

	// Resolve variables in topological order
	resolved := make(map[string]*spookytypesvariables.Variable)
	visited := make(map[string]bool)

	for name := range variables {
		if !visited[name] {
			varResolved, varErrors, varWarnings := m.resolveVariableRecursive(name, variables, graph, resolved, visited, context)
			if varResolved != nil {
				resolved[name] = varResolved
			}
			errors = append(errors, varErrors...)
			warnings = append(warnings, varWarnings...)
		}
	}

	return resolved, errors, warnings
}

// resolveVariableRecursive recursively resolves a variable and its dependencies
func (m *Manager) resolveVariableRecursive(name string, variables map[string]*spookytypesvariables.Variable, graph map[string][]string, resolved map[string]*spookytypesvariables.Variable, visited map[string]bool, context *spookytypesvariables.VariableContext) (*spookytypesvariables.Variable, []spookytypesvariables.VariableError, []spookytypesvariables.VariableWarning) {
	var errors []spookytypesvariables.VariableError
	var warnings []spookytypesvariables.VariableWarning

	if visited[name] {
		return resolved[name], errors, warnings
	}

	visited[name] = true
	variable, exists := variables[name]
	if !exists {
		errors = append(errors, spookytypesvariables.VariableError{
			VariableName: name,
			ErrorType:    spookytypesvariables.VariableErrorTypeMissing,
			ErrorDetails: spookytypes.ErrorDetails{
				Message: fmt.Sprintf("variable '%s' not found", name),
			},
		})
		return nil, errors, warnings
	}

	// Resolve dependencies first
	for _, dep := range variable.Dependencies {
		if !visited[dep] {
			_, depErrors, depWarnings := m.resolveVariableRecursive(dep, variables, graph, resolved, visited, context)
			errors = append(errors, depErrors...)
			warnings = append(warnings, depWarnings...)
		}
	}

	// Resolve the variable value
	resolvedValue, resolveErrors, resolveWarnings := m.resolveVariableValue(variable, resolved, context)
	errors = append(errors, resolveErrors...)
	warnings = append(warnings, resolveWarnings...)

	// Create resolved variable
	resolvedVar := *variable
	resolvedVar.ResolvedValue = resolvedValue
	resolvedVar.IsResolved = len(resolveErrors) == 0
	if len(resolveErrors) > 0 {
		resolvedVar.ResolutionError = resolveErrors[0].ErrorDetails.Message
	}

	return &resolvedVar, errors, warnings
}

// resolveVariableValue resolves the value of a single variable
func (m *Manager) resolveVariableValue(variable *spookytypesvariables.Variable, resolved map[string]*spookytypesvariables.Variable, context *spookytypesvariables.VariableContext) (interface{}, []spookytypesvariables.VariableError, []spookytypesvariables.VariableWarning) {
	var errors []spookytypesvariables.VariableError
	var warnings []spookytypesvariables.VariableWarning

	// Start with default value
	value := variable.Default

	// Apply dependencies if any
	if len(variable.Dependencies) > 0 {
		for _, depName := range variable.Dependencies {
			depVar, exists := resolved[depName]
			if !exists {
				errors = append(errors, spookytypesvariables.VariableError{
					VariableName: variable.Name,
					ErrorType:    spookytypesvariables.VariableErrorTypeDependency,
					ErrorDetails: spookytypes.ErrorDetails{
						Message: fmt.Sprintf("dependency '%s' not resolved", depName),
					},
				})
				continue
			}

			if !depVar.IsResolved {
				errors = append(errors, spookytypesvariables.VariableError{
					VariableName: variable.Name,
					ErrorType:    spookytypesvariables.VariableErrorTypeDependency,
					ErrorDetails: spookytypes.ErrorDetails{
						Message: fmt.Sprintf("dependency '%s' resolution failed", depName),
					},
				})
				continue
			}

			// Use dependency value (simplified - in real implementation, this would be more complex)
			value = depVar.ResolvedValue
		}
	}

	// Apply context-based resolution
	if context != nil {
		contextValue := m.resolveFromContext(variable, context)
		if contextValue != nil {
			value = contextValue
		}
	}

	return value, errors, warnings
}

// resolveFromContext resolves variable value from context
func (m *Manager) resolveFromContext(variable *spookytypesvariables.Variable, context *spookytypesvariables.VariableContext) interface{} {
	// Check environment variables
	if context.Environment != nil {
		if envValue, exists := context.Environment[variable.Name]; exists {
			return envValue
		}
	}

	// Check facts
	if context.Facts != nil {
		if factValue, exists := context.Facts[variable.Name]; exists {
			return factValue
		}
	}

	// Check machines
	if context.Machines != nil {
		if machineValue, exists := context.Machines[variable.Name]; exists {
			return machineValue
		}
	}

	// Check user data
	if context.UserData != nil {
		if userValue, exists := context.UserData[variable.Name]; exists {
			return userValue
		}
	}

	return nil
}

// applyEnvironmentOverrides applies environment variable overrides
func (m *Manager) applyEnvironmentOverrides(variables map[string]*spookytypesvariables.Variable, context *spookytypesvariables.VariableContext) ([]spookytypesvariables.VariableError, []spookytypesvariables.VariableWarning) {
	var errors []spookytypesvariables.VariableError
	var warnings []spookytypesvariables.VariableWarning

	for name, variable := range variables {
		envName := fmt.Sprintf("SPOOKY_%s", strings.ToUpper(name))
		if envValue := os.Getenv(envName); envValue != "" {
			// Override with environment variable
			variable.ResolvedValue = envValue
			variable.IsResolved = true

			warnings = append(warnings, spookytypesvariables.VariableWarning{
				VariableName: name,
				WarningType:  spookytypesvariables.VariableWarningTypeEnvironment,
				ErrorDetails: spookytypes.ErrorDetails{
					Message: fmt.Sprintf("variable overridden by environment variable %s", envName),
				},
			})
		}
	}

	return errors, warnings
}

// validateResolvedVariable validates a resolved variable
func (m *Manager) validateResolvedVariable(variable *spookytypesvariables.Variable) error {
	if !variable.IsResolved {
		return fmt.Errorf("variable is not resolved")
	}

	// Type validation
	if err := m.validateVariableType(variable.ResolvedValue, variable.Type); err != nil {
		return fmt.Errorf("type validation failed: %w", err)
	}

	// Constraint validation
	if variable.Constraints != nil {
		if err := m.validateVariableConstraints(variable.ResolvedValue, variable.Constraints); err != nil {
			return fmt.Errorf("constraint validation failed: %w", err)
		}
	}

	return nil
}

// validateVariableType validates that a value matches the expected type
func (m *Manager) validateVariableType(value interface{}, varType spookytypesvariables.VariableType) error {
	// Basic type checking - in a real implementation, this would be more comprehensive
	switch varType {
	case spookytypesvariables.VariableTypeString:
		if _, ok := value.(string); !ok {
			return fmt.Errorf("expected string, got %T", value)
		}
	case spookytypesvariables.VariableTypeNumber, spookytypesvariables.VariableTypeFloat:
		switch value.(type) {
		case int, int64, float64:
			// Valid numeric types
		default:
			return fmt.Errorf("expected number, got %T", value)
		}
	case spookytypesvariables.VariableTypeBool:
		if _, ok := value.(bool); !ok {
			return fmt.Errorf("expected bool, got %T", value)
		}
		// Add more type validations as needed
	}

	return nil
}

// validateVariableConstraints validates variable constraints
func (m *Manager) validateVariableConstraints(value interface{}, constraints *spookytypesvariables.VariableConstraints) error {
	// String constraints
	if strValue, ok := value.(string); ok {
		if constraints.MinLength != nil && len(strValue) < *constraints.MinLength {
			return fmt.Errorf("string length %d is less than minimum %d", len(strValue), *constraints.MinLength)
		}
		if constraints.MaxLength != nil && len(strValue) > *constraints.MaxLength {
			return fmt.Errorf("string length %d is greater than maximum %d", len(strValue), *constraints.MaxLength)
		}
		if constraints.Pattern != nil {
			matched, err := regexp.MatchString(*constraints.Pattern, strValue)
			if err != nil {
				return fmt.Errorf("invalid pattern: %w", err)
			}
			if !matched {
				return fmt.Errorf("string does not match pattern %s", *constraints.Pattern)
			}
		}
	}

	// Add more constraint validations as needed

	return nil
}

// buildDependencyGraph builds a dependency graph for variables
func (m *Manager) buildDependencyGraph(variables map[string]*spookytypesvariables.Variable) map[string][]string {
	graph := make(map[string][]string)

	for name, variable := range variables {
		graph[name] = variable.Dependencies
	}

	return graph
}

// detectCircularDependencies detects circular dependencies in the graph
func (m *Manager) detectCircularDependencies(graph map[string][]string) [][]string {
	var cycles [][]string
	visited := make(map[string]bool)
	recStack := make(map[string]bool)

	for name := range graph {
		if !visited[name] {
			cycle := m.dfsForCycles(graph, name, visited, recStack, []string{})
			if len(cycle) > 0 {
				cycles = append(cycles, cycle)
			}
		}
	}

	return cycles
}

// dfsForCycles performs DFS to detect cycles
func (m *Manager) dfsForCycles(graph map[string][]string, name string, visited, recStack map[string]bool, path []string) []string {
	visited[name] = true
	recStack[name] = true
	path = append(path, name)

	for _, dep := range graph[name] {
		if !visited[dep] {
			cycle := m.dfsForCycles(graph, dep, visited, recStack, path)
			if len(cycle) > 0 {
				return cycle
			}
		} else if recStack[dep] {
			// Found a cycle
			cycleStart := -1
			for i, p := range path {
				if p == dep {
					cycleStart = i
					break
				}
			}
			if cycleStart != -1 {
				return path[cycleStart:]
			}
		}
	}

	recStack[name] = false
	return nil
}

// validateVariableCollection validates the entire variable collection
func (m *Manager) validateVariableCollection(variables map[string]*spookytypesvariables.Variable) error {
	// Check for duplicate variable names
	nameCount := make(map[string]int)
	for name := range variables {
		nameCount[name]++
	}

	var duplicates []string
	for name, count := range nameCount {
		if count > 1 {
			duplicates = append(duplicates, name)
		}
	}

	if len(duplicates) > 0 {
		return fmt.Errorf("duplicate variable names found: %s", strings.Join(duplicates, ", "))
	}

	return nil
}

// extractResolvedValues extracts resolved values from variables
func (m *Manager) extractResolvedValues(variables map[string]*spookytypesvariables.Variable) map[string]interface{} {
	resolved := make(map[string]interface{})
	for name, variable := range variables {
		if variable.IsResolved {
			resolved[name] = variable.ResolvedValue
		}
	}
	return resolved
}
