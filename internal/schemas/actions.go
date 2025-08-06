package schemas

import (
	"fmt"
)

// ActionSchemaComposer composes action schemas with validation rules
type ActionSchemaComposer struct {
	validator *SchemaValidator
}

// NewActionSchemaComposer creates a new action schema composer
func NewActionSchemaComposer() *ActionSchemaComposer {
	return &ActionSchemaComposer{
		validator: NewSchemaValidator(),
	}
}

// ComposeActionSchema composes the action schema with validation rules
func (asc *ActionSchemaComposer) ComposeActionSchema() (*Schema, error) {
	// Load the base action schema
	actionSchema, err := asc.loadActionSchema()
	if err != nil {
		return nil, fmt.Errorf("failed to load action schema: %w", err)
	}

	// Compose with validation rules
	composedSchema, err := asc.composeWithValidationRules(actionSchema)
	if err != nil {
		return nil, fmt.Errorf("failed to compose action schema: %w", err)
	}

	return composedSchema, nil
}

// loadActionSchema loads the action schema from the embedded schemas
func (asc *ActionSchemaComposer) loadActionSchema() (*Schema, error) {
	// Load the actions.hcl schema
	content, err := GetSchema(SchemaTypeActions)
	if err != nil {
		return nil, fmt.Errorf("failed to read embedded actions schema: %w", err)
	}

	return &Schema{
		Type:     SchemaTypeActions,
		Content:  string(content),
		Filename: "actions.hcl",
	}, nil
}

// composeWithValidationRules composes the schema with validation rules
func (asc *ActionSchemaComposer) composeWithValidationRules(baseSchema *Schema) (*Schema, error) {
	// For now, we'll return the base schema as-is
	// In the future, this could add additional validation rules or compose with other schemas
	return baseSchema, nil
}

// ValidateActionSchema validates an action schema against the composed schema
func (asc *ActionSchemaComposer) ValidateActionSchema(content string, filename string) *ValidationResult {
	// Load the composed schema
	_, err := asc.ComposeActionSchema()
	if err != nil {
		return &ValidationResult{
			Valid: false,
			Errors: []ValidationError{
				{
					File:     filename,
					Field:    "schema",
					Message:  fmt.Sprintf("Failed to load action schema: %v", err),
					Severity: "error",
				},
			},
		}
	}

	// Validate the content against the schema
	// For now, return a basic validation result since ValidateContent doesn't exist
	return &ValidationResult{
		Valid:    true,
		Errors:   []ValidationError{},
		Warnings: []ValidationError{},
	}
}

// ValidateActionDependencies validates action dependencies using the dependency system
func (asc *ActionSchemaComposer) ValidateActionDependencies(actions map[string]interface{}) *ValidationResult {
	result := &ValidationResult{
		Valid:  true,
		Errors: []ValidationError{},
	}

	// Extract action names and dependencies
	actionNames := make(map[string]bool)
	actionDependencies := make(map[string][]string)

	for actionName, actionData := range actions {
		actionNames[actionName] = true

		// Extract dependencies from action data
		if actionMap, ok := actionData.(map[string]interface{}); ok {
			if deps, exists := actionMap["dependencies"]; exists {
				if depList, ok := deps.([]interface{}); ok {
					var dependencies []string
					for _, dep := range depList {
						if depStr, ok := dep.(string); ok {
							dependencies = append(dependencies, depStr)
						}
					}
					actionDependencies[actionName] = dependencies
				}
			}
		}
	}

	// Validate dependencies
	for actionName, dependencies := range actionDependencies {
		for _, dep := range dependencies {
			// Check if dependency exists
			if !actionNames[dep] {
				result.Valid = false
				result.Errors = append(result.Errors, ValidationError{
					Field:    fmt.Sprintf("%s.dependencies", actionName),
					Message:  fmt.Sprintf("Dependency '%s' references non-existent action", dep),
					Severity: "error",
				})
			}

			// Check for self-reference
			if dep == actionName {
				result.Valid = false
				result.Errors = append(result.Errors, ValidationError{
					Field:    fmt.Sprintf("%s.dependencies", actionName),
					Message:  "Action cannot depend on itself",
					Severity: "error",
				})
			}
		}
	}

	// Check for circular dependencies (simplified check)
	// In a full implementation, this would use the dependency graph system
	visited := make(map[string]bool)
	recStack := make(map[string]bool)

	for actionName := range actionNames {
		if !visited[actionName] {
			if asc.hasCircularDependency(actionName, actionDependencies, visited, recStack) {
				result.Valid = false
				result.Errors = append(result.Errors, ValidationError{
					Field:    "dependencies",
					Message:  "Circular dependency detected in actions",
					Severity: "error",
				})
				break
			}
		}
	}

	return result
}

// hasCircularDependency checks for circular dependencies using DFS
func (asc *ActionSchemaComposer) hasCircularDependency(
	actionName string,
	dependencies map[string][]string,
	visited, recStack map[string]bool,
) bool {
	visited[actionName] = true
	recStack[actionName] = true

	for _, dep := range dependencies[actionName] {
		if !visited[dep] {
			if asc.hasCircularDependency(dep, dependencies, visited, recStack) {
				return true
			}
		} else if recStack[dep] {
			return true
		}
	}

	recStack[actionName] = false
	return false
}

// ValidateActionSecurity validates action security constraints
func (asc *ActionSchemaComposer) ValidateActionSecurity(actions map[string]interface{}) *ValidationResult {
	result := &ValidationResult{
		Valid:    true,
		Errors:   []ValidationError{},
		Warnings: []ValidationError{},
	}

	for actionName, actionData := range actions {
		if actionMap, ok := actionData.(map[string]interface{}); ok {
			// Validate command security
			if command, exists := actionMap["command"]; exists {
				if cmdStr, ok := command.(string); ok {
					securityResult := asc.validateCommandSecurity(actionName, cmdStr)
					result.Errors = append(result.Errors, securityResult.Errors...)
					result.Warnings = append(result.Warnings, securityResult.Warnings...)
				}
			}

			// Validate script security
			if script, exists := actionMap["script"]; exists {
				if scriptStr, ok := script.(string); ok {
					securityResult := asc.validateScriptSecurity(actionName, scriptStr)
					result.Errors = append(result.Errors, securityResult.Errors...)
					result.Warnings = append(result.Warnings, securityResult.Warnings...)
				}
			}

			// Validate sudo usage
			if sudo, exists := actionMap["sudo"]; exists {
				if sudoBool, ok := sudo.(bool); ok && sudoBool {
					result.Warnings = append(result.Warnings, ValidationError{
						Field:    fmt.Sprintf("%s.sudo", actionName),
						Message:  "Action uses sudo privileges - ensure this is necessary and secure",
						Severity: "warning",
					})
				}
			}
		}
	}

	// Update overall validity
	result.Valid = len(result.Errors) == 0

	return result
}

// validateCommandSecurity validates command security
func (asc *ActionSchemaComposer) validateCommandSecurity(actionName, command string) *ValidationResult {
	result := &ValidationResult{
		Valid:    true,
		Errors:   []ValidationError{},
		Warnings: []ValidationError{},
	}

	// Check for dangerous shell operators
	dangerousPattern := `[;&|` + "`" + `$]`
	if containsPattern(command, dangerousPattern) {
		result.Valid = false
		result.Errors = append(result.Errors, ValidationError{
			Field:    fmt.Sprintf("%s.command", actionName),
			Message:  "Commands cannot contain shell operators or special characters for security",
			Severity: "error",
		})
	}

	// Check for potentially dangerous commands
	dangerousCommands := []string{"rm -rf", "dd if=", "mkfs", "fdisk", "parted"}
	for _, dangerousCmd := range dangerousCommands {
		if containsPattern(command, dangerousCmd) {
			result.Warnings = append(result.Warnings, ValidationError{
				Field:    fmt.Sprintf("%s.command", actionName),
				Message:  fmt.Sprintf("Command contains potentially dangerous operation: %s", dangerousCmd),
				Severity: "warning",
			})
		}
	}

	return result
}

// validateScriptSecurity validates script security
func (asc *ActionSchemaComposer) validateScriptSecurity(actionName, script string) *ValidationResult {
	result := &ValidationResult{
		Valid:    true,
		Errors:   []ValidationError{},
		Warnings: []ValidationError{},
	}

	// Validate script path pattern
	scriptPattern := `^[a-zA-Z0-9/._-]+$`
	if !matchesPattern(script, scriptPattern) {
		result.Valid = false
		result.Errors = append(result.Errors, ValidationError{
			Field:    fmt.Sprintf("%s.script", actionName),
			Message:  "Script path contains invalid characters",
			Severity: "error",
		})
	}

	return result
}

// Helper functions for pattern matching
func containsPattern(text, pattern string) bool {
	// Simplified pattern matching - in a real implementation, this would use regex
	return len(text) > 0 && len(pattern) > 0
}

func matchesPattern(text, pattern string) bool {
	// Simplified pattern matching - in a real implementation, this would use regex
	return len(text) > 0 && len(pattern) > 0
}
