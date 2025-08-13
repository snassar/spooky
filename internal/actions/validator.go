package spookyactions

import (
	"context"
	"fmt"
	"time"

	spookyinterfaces "spooky/internal/interfaces"
	spookylogging "spooky/internal/logging"
	spookytypes "spooky/internal/types"
	spookytypesactions "spooky/internal/types/actions"
)

// Validator implements the ActionValidator interface
type Validator struct {
	logger spookylogging.Logger
}

// NewValidator creates a new action validator
func NewValidator(logger spookylogging.Logger) spookyinterfaces.ActionValidator {
	return &Validator{
		logger: logger,
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
	for i, action := range actions {
		actionResult, err := v.ValidateAction(ctx, action)
		if err != nil {
			return nil, fmt.Errorf("failed to validate action %d: %w", i, err)
		}

		// Aggregate results
		if !actionResult.Valid {
			result.Valid = false
		}
		result.Errors = append(result.Errors, actionResult.Errors...)
		result.Warnings = append(result.Warnings, actionResult.Warnings...)
		result.Info = append(result.Info, actionResult.Info...)
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
func (v *Validator) ValidateAction(ctx context.Context, action spookytypes.Action) (*spookytypes.ValidationResult, error) {
	v.logger.Debug("Validating action", map[string]interface{}{
		"action": action.Name,
	})

	result := &spookytypes.ValidationResult{
		Valid:       true,
		ValidatedAt: time.Now(),
		Errors:      []spookytypes.SchemaError{},
		Warnings:    []spookytypes.SchemaError{},
		Info:        []spookytypes.SchemaError{},
	}

	// Validate required fields
	if action.Name == "" {
		result.Valid = false
		result.Errors = append(result.Errors, spookytypes.SchemaError{
			Code:        "missing_name",
			Message:     "Action name is required",
			Severity:    "error",
			Recoverable: false,
		})
	}

	if action.Type == "" {
		result.Valid = false
		result.Errors = append(result.Errors, spookytypes.SchemaError{
			Code:        "missing_type",
			Message:     "Action type is required",
			Severity:    "error",
			Recoverable: false,
		})
	}

	// Validate action type-specific requirements
	switch action.Type {
	case string(spookytypesactions.ActionTypeCommand):
		if action.Command == "" {
			result.Valid = false
			result.Errors = append(result.Errors, spookytypes.SchemaError{
				Code:        "missing_command",
				Message:     "Command action requires a command",
				Severity:    "error",
				Recoverable: false,
			})
		}

	case string(spookytypesactions.ActionTypeScript):
		if action.Script == "" {
			result.Valid = false
			result.Errors = append(result.Errors, spookytypes.SchemaError{
				Code:        "missing_script",
				Message:     "Script action requires a script path",
				Severity:    "error",
				Recoverable: false,
			})
		}

	case string(spookytypesactions.ActionTypeTemplateDeploy):
		if action.Template == nil {
			result.Valid = false
			result.Errors = append(result.Errors, spookytypes.SchemaError{
				Code:        "missing_template",
				Message:     "Template action requires template configuration",
				Severity:    "error",
				Recoverable: false,
			})
		}

	case string(spookytypesactions.ActionTypeFileCopy):
		if action.FileCopy == nil {
			result.Valid = false
			result.Errors = append(result.Errors, spookytypes.SchemaError{
				Code:        "missing_file_copy",
				Message:     "File copy action requires file copy configuration",
				Severity:    "error",
				Recoverable: false,
			})
		}

	case string(spookytypesactions.ActionTypeServiceControl):
		if action.ServiceControl == nil {
			result.Valid = false
			result.Errors = append(result.Errors, spookytypes.SchemaError{
				Code:        "missing_service_control",
				Message:     "Service control action requires service control configuration",
				Severity:    "error",
				Recoverable: false,
			})
		}

	default:
		result.Valid = false
		result.Errors = append(result.Errors, spookytypes.SchemaError{
			Code:        "invalid_type",
			Message:     fmt.Sprintf("Invalid action type: %s", action.Type),
			Severity:    "error",
			Recoverable: false,
		})
	}

	// Validate timeout
	if action.Timeout <= 0 {
		result.Warnings = append(result.Warnings, spookytypes.SchemaError{
			Code:        "no_timeout",
			Message:     "No timeout specified, using default",
			Severity:    "warning",
			Recoverable: true,
		})
	}

	// Validate retries
	if action.Retries < 0 {
		result.Valid = false
		result.Errors = append(result.Errors, spookytypes.SchemaError{
			Code:        "invalid_retries",
			Message:     "Retries must be non-negative",
			Severity:    "error",
			Recoverable: false,
		})
	}

	return result, nil
}

// validateDependencies validates dependencies across all actions
func (v *Validator) validateDependencies(actions []spookytypes.Action) error {
	// Build dependency graph
	dependencies := make(map[string][]string)
	actionNames := make(map[string]bool)

	// Collect all action names
	for _, action := range actions {
		actionNames[action.Name] = true
	}

	// Build dependency map
	for _, action := range actions {
		for _, dep := range action.Dependencies {
			if !actionNames[dep] {
				return fmt.Errorf("action '%s' depends on undefined action '%s'", action.Name, dep)
			}
			dependencies[action.Name] = append(dependencies[action.Name], dep)
		}
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
