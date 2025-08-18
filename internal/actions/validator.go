package spookyactions

import (
	"context"
	"fmt"
	"time"

	spookyinterfaces "spooky/internal/interfaces"
	spookyschemas "spooky/internal/schemas"
	spookytypes "spooky/internal/types"
	spookytypesactions "spooky/internal/types/actions"
	spookytypeslogging "spooky/internal/types/logging"
	spookytypesschemas "spooky/internal/types/schemas"
)

// Validator implements the ActionValidator interface
type Validator struct {
	logger                spookytypeslogging.Logger
	schemaDrivenValidator *spookyschemas.SchemaDrivenValidator
	enhancedValidator     *spookyschemas.EnhancedValidator
}

// NewValidator creates a new action validator
func NewValidator(logger spookytypeslogging.Logger) spookyinterfaces.ActionValidator {
	// Create schema-driven validator for action configuration validation
	schemaDrivenConfig := &spookyschemas.SchemaDrivenValidationConfig{
		UseEmbeddedSchemas: true,
		StrictValidation:   true,
		AllowUnknownFields: false,
		DetailedErrors:     true,
	}
	schemaDrivenValidator := spookyschemas.NewSchemaDrivenValidator(logger, schemaDrivenConfig)

	// Create enhanced validator for individual action validation
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

// ValidateActions validates a collection of actions
func (v *Validator) ValidateActions(ctx context.Context, actions []spookytypes.Action) (*spookytypes.ValidationResult, error) {
	v.logger.Info("Validating actions", map[string]interface{}{
		"count": len(actions),
	})

	result := &spookytypes.ValidationResult{
		Valid:       true,
		ValidatedAt: time.Now(),
		Errors:      []spookytypes.SchemaError{},
		Warnings:    []spookytypes.SchemaError{},
		Info:        []spookytypes.SchemaError{},
	}

	// Validate each action
	for i := range actions {
		actionResult, err := v.ValidateAction(ctx, &actions[i])
		if err != nil {
			return nil, fmt.Errorf("failed to validate action %s: %w", actions[i].Name, err)
		}

		if !actionResult.Valid {
			result.Valid = false
			result.Errors = append(result.Errors, actionResult.Errors...)
		}

		result.Warnings = append(result.Warnings, actionResult.Warnings...)
	}

	// Validate dependencies across all actions
	if err := v.validateDependencies(actions); err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, spookytypes.SchemaError{
			Code:        "dependency_error",
			Message:     err.Error(),
			Severity:    "error",
			Recoverable: false,
		})
	}

	v.logger.Info("Actions validation completed", map[string]interface{}{
		"valid":         result.Valid,
		"error_count":   len(result.Errors),
		"warning_count": len(result.Warnings),
	})

	return result, nil
}

// ValidateAction validates a single action
func (v *Validator) ValidateAction(ctx context.Context, action *spookytypes.Action) (*spookytypes.ValidationResult, error) {
	v.logger.Debug("Validating action", map[string]interface{}{
		"action": action.Name,
		"type":   action.Type,
	})

	// Get action schema for enhanced validation
	actionSchema, err := v.getActionSchema()
	if err != nil {
		return nil, fmt.Errorf("failed to get action schema: %w", err)
	}

	// Use enhanced validator for comprehensive action validation
	result, err := v.enhancedValidator.ValidateWithEnhancedFeatures(ctx, actionSchema, action)
	if err != nil {
		return nil, fmt.Errorf("failed to validate action with enhanced validator: %w", err)
	}

	// Add additional custom validation for action-specific rules
	v.addCustomActionValidation(action, result)

	return &spookytypes.ValidationResult{
		Valid:    result.Valid,
		Errors:   result.Errors,
		Warnings: result.Warnings,
	}, nil
}

// getActionSchema gets the action schema for validation
func (v *Validator) getActionSchema() (*spookytypesschemas.Schema, error) {
	// Try to get schema from embedded schemas first
	if schema, err := v.schemaDrivenValidator.GetEmbeddedSchema("actions"); err == nil {
		return schema, nil
	}

	// Fallback: create a basic action schema
	return &spookytypesschemas.Schema{
		Name:        "actions",
		Type:        "hcl",
		Version:     "1.0",
		Description: "Action configuration schema",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Content:     "", // Will be loaded from file if needed
		Metadata:    make(map[string]interface{}),
	}, nil
}

// addCustomActionValidation adds custom validation rules specific to actions
func (v *Validator) addCustomActionValidation(action *spookytypes.Action, result *spookytypesschemas.ValidationResult) {
	// Validate required fields
	if action.Name == "" {
		v.addSchemaError(result, "missing_name", "Action name is required", "error")
	}

	if action.Type == "" {
		v.addSchemaError(result, "missing_type", "Action type is required", "error")
	}

	// Validate action type
	if err := v.validateActionType(action.Type); err != nil {
		v.addSchemaError(result, "invalid_type", err.Error(), "error")
	}

	// Validate action parameters
	if err := v.validateActionParameters(action); err != nil {
		v.addSchemaError(result, "invalid_parameters", err.Error(), "error")
	}

	// Validate action dependencies
	if err := v.validateActionDependencies(action); err != nil {
		v.addSchemaError(result, "invalid_dependencies", err.Error(), "error")
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

// validateDependencies validates dependencies across all actions
func (v *Validator) validateDependencies(actions []spookytypes.Action) error {
	// Build dependency graph
	dependencies := make(map[string][]string)
	for i := range actions {
		dependencies[actions[i].Name] = actions[i].Dependencies
	}

	// Check for circular dependencies
	if err := v.checkCircularDependencies(dependencies); err != nil {
		return fmt.Errorf("circular dependency detected: %w", err)
	}

	return nil
}

// checkCircularDependencies checks for circular dependencies using DFS
func (v *Validator) checkCircularDependencies(dependencies map[string][]string) error {
	visited := make(map[string]bool)
	recStack := make(map[string]bool)

	for action := range dependencies {
		if !visited[action] {
			if v.hasCycle(action, dependencies, visited, recStack) {
				return fmt.Errorf("circular dependency detected")
			}
		}
	}

	return nil
}

// hasCycle performs DFS to detect cycles
func (v *Validator) hasCycle(action string, dependencies map[string][]string, visited, recStack map[string]bool) bool {
	visited[action] = true
	recStack[action] = true

	for _, dep := range dependencies[action] {
		if !visited[dep] {
			if v.hasCycle(dep, dependencies, visited, recStack) {
				return true
			}
		} else if recStack[dep] {
			return true
		}
	}

	recStack[action] = false
	return false
}

// validateActionType validates the action type
func (v *Validator) validateActionType(actionType string) error {
	validTypes := []string{
		string(spookytypesactions.ActionTypeCommand),
		string(spookytypesactions.ActionTypeScript),
		string(spookytypesactions.ActionTypeTemplateDeploy),
		string(spookytypesactions.ActionTypeFileCopy),
		string(spookytypesactions.ActionTypeServiceControl),
	}

	for _, validType := range validTypes {
		if actionType == validType {
			return nil
		}
	}

	return fmt.Errorf("invalid action type: %s", actionType)
}

// validateActionParameters validates action parameters based on type
func (v *Validator) validateActionParameters(action *spookytypes.Action) error {
	validationMap := map[string]func(*spookytypes.Action) error{
		string(spookytypesactions.ActionTypeCommand):        v.validateCommandAction,
		string(spookytypesactions.ActionTypeScript):         v.validateScriptAction,
		string(spookytypesactions.ActionTypeTemplateDeploy): v.validateTemplateAction,
		string(spookytypesactions.ActionTypeFileCopy):       v.validateFileCopyAction,
		string(spookytypesactions.ActionTypeServiceControl): v.validateServiceControlAction,
	}

	if validator, exists := validationMap[action.Type]; exists {
		return validator(action)
	}

	return fmt.Errorf("unknown action type: %s", action.Type)
}

// validateCommandAction validates command action parameters
func (v *Validator) validateCommandAction(action *spookytypes.Action) error {
	if action.Command == nil {
		return fmt.Errorf("command action requires a command")
	}
	return nil
}

// validateScriptAction validates script action parameters
func (v *Validator) validateScriptAction(action *spookytypes.Action) error {
	if action.Script == nil {
		return fmt.Errorf("script action requires a script path")
	}
	return nil
}

// validateTemplateAction validates template action parameters
func (v *Validator) validateTemplateAction(action *spookytypes.Action) error {
	if action.Template == nil {
		return fmt.Errorf("template action requires template configuration")
	}
	return nil
}

// validateFileCopyAction validates file copy action parameters
func (v *Validator) validateFileCopyAction(action *spookytypes.Action) error {
	if action.FileCopy == nil {
		return fmt.Errorf("file copy action requires file copy configuration")
	}
	return nil
}

// validateServiceControlAction validates service control action parameters
func (v *Validator) validateServiceControlAction(action *spookytypes.Action) error {
	if action.ServiceControl == nil {
		return fmt.Errorf("service control action requires service control configuration")
	}
	return nil
}

// validateActionDependencies validates action dependencies
func (v *Validator) validateActionDependencies(action *spookytypes.Action) error {
	// Validate timeout
	if action.Timeout <= 0 {
		return fmt.Errorf("timeout must be positive")
	}

	// Validate retries
	if action.Retries < 0 {
		return fmt.Errorf("retries must be non-negative")
	}

	return nil
}
