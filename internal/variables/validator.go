// Package variables provides variable validation functionality for the spooky codebase.
package variables

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	spookyinterfaces "spooky/internal/interfaces"
	spookyschemas "spooky/internal/schemas"
	spookytypes "spooky/internal/types"
	spookytypeslogging "spooky/internal/types/logging"
	spookytypesschemas "spooky/internal/types/schemas"
	spookytypesvariables "spooky/internal/types/variables"
)

// Validator implements the VariableValidator interface
type Validator struct {
	logger                spookytypeslogging.Logger
	schemaDrivenValidator *spookyschemas.SchemaDrivenValidator
	enhancedValidator     *spookyschemas.EnhancedValidator
}

// NewValidator creates a new VariableValidator instance
func NewValidator(logger spookytypeslogging.Logger) spookyinterfaces.VariableValidator {
	// Create schema-driven validator for variable configuration validation
	schemaDrivenConfig := &spookyschemas.SchemaDrivenValidationConfig{
		UseEmbeddedSchemas: true,
		StrictValidation:   true,
		AllowUnknownFields: false,
		DetailedErrors:     true,
	}
	schemaDrivenValidator := spookyschemas.NewSchemaDrivenValidator(logger, schemaDrivenConfig)

	// Create enhanced validator for individual variable validation
	enhancedConfig := &spookyschemas.ValidationConfig{
		Mode: spookyschemas.ValidationModeStrict,
		ErrorHandling: &spookyschemas.ErrorHandlingConfig{
			StopOnFirstError:   false,
			MaxErrors:          100,
			IncludeWarnings:    true,
			IncludeContext:     true,
			IncludeSuggestions: true,
		},
		Evolution: &spookyschemas.EvolutionConfig{
			EnableTracking:  true,
			AllowDeprecated: true,
			WarnDeprecated:  true,
			AllowBreaking:   false,
		},
	}
	enhancedValidator := spookyschemas.NewEnhancedValidator(enhancedConfig)

	return &Validator{
		logger:                logger,
		schemaDrivenValidator: schemaDrivenValidator,
		enhancedValidator:     enhancedValidator,
	}
}

// ValidateVariables validates a collection of variables
func (v *Validator) ValidateVariables(ctx context.Context, variables map[string]*spookytypesvariables.Variable) (*spookytypes.ValidationResult, error) {
	v.logger.Debug("Validating variables", map[string]interface{}{
		"count": len(variables),
	})

	var errors []spookytypesschemas.SchemaError
	var warnings []spookytypesschemas.SchemaError

	for name, variable := range variables {
		result, err := v.ValidateVariable(ctx, variable)
		if err != nil {
			return nil, fmt.Errorf("failed to validate variable %s: %w", name, err)
		}

		if !result.Valid {
			errors = append(errors, result.Errors...)
		}

		warnings = append(warnings, result.Warnings...)
	}

	// Validate variable collection (dependencies, duplicates, etc.)
	collectionErrors, collectionWarnings := v.validateVariableCollection(variables)
	errors = append(errors, collectionErrors...)
	warnings = append(warnings, collectionWarnings...)

	valid := len(errors) == 0

	return &spookytypes.ValidationResult{
		Valid:    valid,
		Errors:   errors,
		Warnings: warnings,
	}, nil
}

// ValidateVariable validates a single variable
func (v *Validator) ValidateVariable(ctx context.Context, variable *spookytypesvariables.Variable) (*spookytypes.ValidationResult, error) {
	v.logger.Debug("Validating variable", map[string]interface{}{
		"name": variable.Name,
		"type": variable.Type,
	})

	// Get variable schema for enhanced validation
	variableSchema, err := v.getVariableSchema()
	if err != nil {
		return nil, fmt.Errorf("failed to get variable schema: %w", err)
	}

	// Use enhanced validator for comprehensive variable validation
	result, err := v.enhancedValidator.ValidateWithEnhancedFeatures(ctx, variableSchema, variable)
	if err != nil {
		return nil, fmt.Errorf("failed to validate variable with enhanced validator: %w", err)
	}

	// Add additional custom validation for variable-specific rules
	v.addCustomVariableValidation(variable, result)

	return &spookytypes.ValidationResult{
		Valid:    result.Valid,
		Errors:   result.Errors,
		Warnings: result.Warnings,
	}, nil
}

// getVariableSchema gets the variable schema for validation
func (v *Validator) getVariableSchema() (*spookytypesschemas.Schema, error) {
	// Try to get schema from embedded schemas first
	if schema, err := v.schemaDrivenValidator.GetEmbeddedSchema("variables"); err == nil {
		return schema, nil
	}

	// Fallback: create a basic variable schema
	return &spookytypesschemas.Schema{
		Name:        "variables",
		Type:        "hcl",
		Version:     "1.0",
		Description: "Variable configuration schema",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Content:     "", // Will be loaded from file if needed
		Metadata:    make(map[string]interface{}),
	}, nil
}

// addCustomVariableValidation adds custom validation rules specific to variables
func (v *Validator) addCustomVariableValidation(variable *spookytypesvariables.Variable, result *spookytypesschemas.ValidationResult) {
	// Validate variable name
	if err := v.validateVariableName(variable.Name); err != nil {
		v.addSchemaError(result, "invalid_name", err.Error(), "error")
	}

	// Validate variable type
	if err := v.validateVariableType(variable.Type); err != nil {
		v.addSchemaError(result, "invalid_type", err.Error(), "error")
	}

	// Validate variable scope
	if err := v.validateVariableScope(variable.Scope); err != nil {
		v.addSchemaError(result, "invalid_scope", err.Error(), "error")
	}

	// Validate variable dependencies
	if err := v.validateVariableDependencies(variable); err != nil {
		v.addSchemaError(result, "invalid_dependencies", err.Error(), "error")
	}

	// Validate variable constraints
	if variable.Constraints != nil {
		if err := v.validateVariableConstraints(variable); err != nil {
			v.addSchemaError(result, "invalid_constraints", err.Error(), "error")
		}
	}
}

// addSchemaError adds a schema error to the validation result
func (v *Validator) addSchemaError(result *spookytypesschemas.ValidationResult, code, message, severity string) {
	schemaError := spookytypesschemas.SchemaError{
		Code:     code,
		Message:  message,
		Severity: severity,
	}
	result.Errors = append(result.Errors, schemaError)
	result.Valid = false
}

// validateVariableName validates the variable name
func (v *Validator) validateVariableName(name string) error {
	if name == "" {
		return fmt.Errorf("variable name cannot be empty")
	}

	// Check naming pattern: lowercase with underscores, starting with a letter
	pattern := regexp.MustCompile(`^[a-z][a-z0-9_]*$`)
	if !pattern.MatchString(name) {
		return fmt.Errorf("variable name must be lowercase with underscores, starting with a letter")
	}

	if len(name) > 100 {
		return fmt.Errorf("variable name cannot exceed 100 characters")
	}

	return nil
}

// validateVariableType validates the variable type
func (v *Validator) validateVariableType(varType spookytypesvariables.VariableType) error {
	validTypes := []spookytypesvariables.VariableType{
		spookytypesvariables.VariableTypeString,
		spookytypesvariables.VariableTypeNumber,
		spookytypesvariables.VariableTypeFloat,
		spookytypesvariables.VariableTypeBool,
		spookytypesvariables.VariableTypeList,
		spookytypesvariables.VariableTypeMap,
		spookytypesvariables.VariableTypeObject,
		spookytypesvariables.VariableTypeDuration,
		spookytypesvariables.VariableTypeIP,
		spookytypesvariables.VariableTypeCIDR,
		spookytypesvariables.VariableTypePath,
		spookytypesvariables.VariableTypeFile,
		spookytypesvariables.VariableTypeSecret,
	}

	for _, validType := range validTypes {
		if varType == validType {
			return nil
		}
	}

	return fmt.Errorf("invalid variable type: %s", varType)
}

// validateVariableScope validates the variable scope
func (v *Validator) validateVariableScope(scope spookytypesvariables.VariableScope) error {
	validScopes := []spookytypesvariables.VariableScope{
		spookytypesvariables.VariableScopeProject,
		spookytypesvariables.VariableScopeGlobal,
		spookytypesvariables.VariableScopeInherited,
	}

	for _, validScope := range validScopes {
		if scope == validScope {
			return nil
		}
	}

	return fmt.Errorf("invalid variable scope: %s", scope)
}

// validateVariableDependencies validates variable dependencies
func (v *Validator) validateVariableDependencies(variable *spookytypesvariables.Variable) error {
	if len(variable.Dependencies) == 0 {
		return nil
	}

	// Check for self-reference
	for _, dep := range variable.Dependencies {
		if dep == variable.Name {
			return fmt.Errorf("variable cannot depend on itself")
		}
	}

	// Check for duplicate dependencies
	seen := make(map[string]bool)
	for _, dep := range variable.Dependencies {
		if seen[dep] {
			return fmt.Errorf("duplicate dependency: %s", dep)
		}
		seen[dep] = true
	}

	return nil
}

// validateVariableConstraints validates variable constraints
func (v *Validator) validateVariableConstraints(variable *spookytypesvariables.Variable) error {
	if err := v.validateStringConstraints(variable); err != nil {
		return err
	}
	if err := v.validateNumericConstraints(variable); err != nil {
		return err
	}
	if err := v.validateListConstraints(variable); err != nil {
		return err
	}
	if err := v.validateFileConstraints(variable); err != nil {
		return err
	}
	return nil
}

// validateStringConstraints validates string-specific constraints
func (v *Validator) validateStringConstraints(variable *spookytypesvariables.Variable) error {
	if variable.Type != spookytypesvariables.VariableTypeString {
		return nil
	}

	constraints := variable.Constraints

	// Validate length constraints
	if err := v.validateStringLengthConstraints(constraints); err != nil {
		return err
	}

	// Validate pattern constraint
	if err := v.validateStringPatternConstraint(constraints); err != nil {
		return err
	}

	return nil
}

// validateStringLengthConstraints validates string length constraints
func (v *Validator) validateStringLengthConstraints(constraints *spookytypesvariables.VariableConstraints) error {
	if constraints.MinLength != nil && constraints.MaxLength != nil {
		if *constraints.MinLength > *constraints.MaxLength {
			return fmt.Errorf("min_length cannot be greater than max_length")
		}
	}

	if constraints.MinLength != nil && *constraints.MinLength < 0 {
		return fmt.Errorf("min_length cannot be negative")
	}

	if constraints.MaxLength != nil && *constraints.MaxLength < 1 {
		return fmt.Errorf("max_length must be at least 1")
	}

	return nil
}

// validateStringPatternConstraint validates string pattern constraint
func (v *Validator) validateStringPatternConstraint(constraints *spookytypesvariables.VariableConstraints) error {
	if constraints.Pattern != nil {
		if _, err := regexp.Compile(*constraints.Pattern); err != nil {
			return fmt.Errorf("invalid pattern: %w", err)
		}
	}

	return nil
}

// validateNumericConstraints validates numeric-specific constraints
func (v *Validator) validateNumericConstraints(variable *spookytypesvariables.Variable) error {
	if variable.Type != spookytypesvariables.VariableTypeNumber && variable.Type != spookytypesvariables.VariableTypeFloat {
		return nil
	}

	constraints := variable.Constraints

	if constraints.MinValue != nil && constraints.MaxValue != nil {
		if *constraints.MinValue > *constraints.MaxValue {
			return fmt.Errorf("min_value cannot be greater than max_value")
		}
	}

	return nil
}

// validateListConstraints validates list-specific constraints
func (v *Validator) validateListConstraints(variable *spookytypesvariables.Variable) error {
	if variable.Type != spookytypesvariables.VariableTypeList {
		return nil
	}

	constraints := variable.Constraints

	if constraints.MinItems != nil && constraints.MaxItems != nil {
		if *constraints.MinItems > *constraints.MaxItems {
			return fmt.Errorf("min_items cannot be greater than max_items")
		}
	}

	if constraints.MinItems != nil && *constraints.MinItems < 0 {
		return fmt.Errorf("min_items cannot be negative")
	}

	if constraints.MaxItems != nil && *constraints.MaxItems < 1 {
		return fmt.Errorf("max_items must be at least 1")
	}

	return nil
}

// validateFileConstraints validates file-specific constraints
func (v *Validator) validateFileConstraints(variable *spookytypesvariables.Variable) error {
	if variable.Type != spookytypesvariables.VariableTypeFile {
		return nil
	}

	constraints := variable.Constraints

	if constraints.FileSizeMax != nil {
		if err := v.validateFileSizeConstraint(*constraints.FileSizeMax); err != nil {
			return fmt.Errorf("invalid file_size_max: %w", err)
		}
	}

	return nil
}

// validateVariableCollection validates the entire variable collection
func (v *Validator) validateVariableCollection(variables map[string]*spookytypesvariables.Variable) (errors, warnings []spookytypesschemas.SchemaError) {
	// Check for circular dependencies
	cycles := v.detectCircularDependencies(variables)
	if len(cycles) > 0 {
		for _, cycle := range cycles {
			schemaError := spookytypesschemas.NewSchemaError(
				"Variable",
				"dependencies",
				fmt.Sprintf("circular dependency detected: %s", strings.Join(cycle, " -> ")),
			)
			errors = append(errors, *schemaError)
		}
	}

	// Check for unused variables
	unusedVars := v.detectUnusedVariables(variables)
	for _, unusedVar := range unusedVars {
		warning := spookytypesschemas.NewSchemaError(
			"Variable",
			"name",
			fmt.Sprintf("variable '%s' is defined but not used", unusedVar),
		)
		warnings = append(warnings, *warning)
	}

	return errors, warnings
}

// detectCircularDependencies detects circular dependencies in variables
func (v *Validator) detectCircularDependencies(variables map[string]*spookytypesvariables.Variable) [][]string {
	var cycles [][]string
	visited := make(map[string]bool)
	recStack := make(map[string]bool)

	for name := range variables {
		if !visited[name] {
			cycle := v.dfsForCycles(variables, name, visited, recStack, []string{})
			if len(cycle) > 0 {
				cycles = append(cycles, cycle)
			}
		}
	}

	return cycles
}

// dfsForCycles performs DFS to detect cycles
func (v *Validator) dfsForCycles(variables map[string]*spookytypesvariables.Variable, name string, visited, recStack map[string]bool, path []string) []string {
	visited[name] = true
	recStack[name] = true
	path = append(path, name)

	variable, exists := variables[name]
	if exists {
		for _, dep := range variable.Dependencies {
			if !visited[dep] {
				cycle := v.dfsForCycles(variables, dep, visited, recStack, path)
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
	}

	recStack[name] = false
	return nil
}

// detectUnusedVariables detects variables that are not used by other variables
func (v *Validator) detectUnusedVariables(variables map[string]*spookytypesvariables.Variable) []string {
	used := make(map[string]bool)

	// Mark all variables that are dependencies of other variables
	for _, variable := range variables {
		for _, dep := range variable.Dependencies {
			used[dep] = true
		}
	}

	// Find variables that are not used
	var unused []string
	for name := range variables {
		if !used[name] {
			unused = append(unused, name)
		}
	}

	return unused
}

// validateFileSizeConstraint validates file size constraint format
func (v *Validator) validateFileSizeConstraint(sizeStr string) error {
	// Pattern: ^\d+[KMGT]?B$ or ^\d+\s*[KMGT]?B$
	pattern := regexp.MustCompile(`^\d+\s*[KMGT]?B$`)
	if !pattern.MatchString(sizeStr) {
		return fmt.Errorf("file size must be in format: number[K|M|G|T]B (e.g., '10MB', '1GB', '500KB', '2TB')")
	}

	return nil
}
