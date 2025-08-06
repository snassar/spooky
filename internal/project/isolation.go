package project

import (
	"fmt"
	"strings"
)

// ProjectIsolationEngine manages project isolation and security settings
type ProjectIsolationEngine struct{}

// NewProjectIsolationEngine creates a new project isolation engine
func NewProjectIsolationEngine() *ProjectIsolationEngine {
	return &ProjectIsolationEngine{}
}

// IsolationContext represents the isolation context for a project
type IsolationContext struct {
	Project           *Project
	IsolationSettings *ProjectIsolation
	FactsScope        string
	VariablesScope    string
	MachineAccess     string
	AllowedMachines   []string
	AllowedTags       []string
	IsolationEnabled  bool
}

// IsolationResult represents the result of isolation validation
type IsolationResult struct {
	Valid    bool
	Errors   []ValidationError
	Warnings []ValidationError
}

// ApplyIsolation applies isolation settings to a project
func (ie *ProjectIsolationEngine) ApplyIsolation(project *Project) (*IsolationContext, error) {
	if project == nil {
		return nil, fmt.Errorf("project cannot be nil")
	}

	// Create isolation context with defaults
	context := &IsolationContext{
		Project:          project,
		IsolationEnabled: true,
		FactsScope:       DefaultFactsScope,
		VariablesScope:   DefaultVariablesScope,
		MachineAccess:    DefaultMachineAccess,
		AllowedMachines:  []string{}, // always non-nil
		AllowedTags:      []string{}, // always non-nil
	}

	// Apply project-specific isolation settings if available
	if project.Isolation != nil {
		context.IsolationSettings = project.Isolation
		context.IsolationEnabled = project.Isolation.Enabled

		if project.Isolation.FactsScope != "" {
			context.FactsScope = project.Isolation.FactsScope
		}

		if project.Isolation.VariablesScope != "" {
			context.VariablesScope = project.Isolation.VariablesScope
		}

		if project.Isolation.MachineAccess != "" {
			context.MachineAccess = project.Isolation.MachineAccess
		}

		context.AllowedMachines = project.Isolation.AllowedMachines
		context.AllowedTags = project.Isolation.AllowedTags
	}
	// Ensure AllowedMachines and AllowedTags are never nil
	if context.AllowedMachines == nil {
		context.AllowedMachines = []string{}
	}
	if context.AllowedTags == nil {
		context.AllowedTags = []string{}
	}

	// Validate isolation context
	if err := ie.validateIsolationContext(context); err != nil {
		return nil, fmt.Errorf("isolation context validation failed: %w", err)
	}

	return context, nil
}

// ValidateIsolation validates project isolation configuration
func (ie *ProjectIsolationEngine) ValidateIsolation(project *Project) *IsolationResult {
	result := &IsolationResult{
		Valid: true,
	}

	if project == nil {
		result.Valid = false
		result.Errors = append(result.Errors, ValidationError{
			Field:    "project",
			Message:  "project cannot be nil",
			Severity: "error",
		})
		return result
	}

	// Validate isolation settings if present
	if project.Isolation != nil {
		// Validate facts scope
		if project.Isolation.FactsScope != "" {
			validScopes := []string{"global", "project", "hybrid"}
			valid := false
			for _, scope := range validScopes {
				if project.Isolation.FactsScope == scope {
					valid = true
					break
				}
			}
			if !valid {
				result.Valid = false
				result.Errors = append(result.Errors, ValidationError{
					Field:    "isolation.facts_scope",
					Message:  fmt.Sprintf("invalid facts scope: %s (must be one of: %s)", project.Isolation.FactsScope, strings.Join(validScopes, ", ")),
					Severity: "error",
				})
			}
		}

		// Validate variables scope
		if project.Isolation.VariablesScope != "" {
			validScopes := []string{"project", "inherited"}
			valid := false
			for _, scope := range validScopes {
				if project.Isolation.VariablesScope == scope {
					valid = true
					break
				}
			}
			if !valid {
				result.Valid = false
				result.Errors = append(result.Errors, ValidationError{
					Field:    "isolation.variables_scope",
					Message:  fmt.Sprintf("invalid variables scope: %s (must be one of: %s)", project.Isolation.VariablesScope, strings.Join(validScopes, ", ")),
					Severity: "error",
				})
			}
		}

		// Validate machine access
		if project.Isolation.MachineAccess != "" {
			validAccess := []string{"all", "tagged", "explicit"}
			valid := false
			for _, access := range validAccess {
				if project.Isolation.MachineAccess == access {
					valid = true
					break
				}
			}
			if !valid {
				result.Valid = false
				result.Errors = append(result.Errors, ValidationError{
					Field:    "isolation.machine_access",
					Message:  fmt.Sprintf("invalid machine access: %s (must be one of: %s)", project.Isolation.MachineAccess, strings.Join(validAccess, ", ")),
					Severity: "error",
				})
			}
		}

		// Validate explicit machine access configuration
		if project.Isolation.MachineAccess == "explicit" {
			if len(project.Isolation.AllowedMachines) == 0 && len(project.Isolation.AllowedTags) == 0 {
				result.Valid = false
				result.Errors = append(result.Errors, ValidationError{
					Field:    "isolation",
					Message:  "explicit machine access requires allowed_machines or allowed_tags to be specified",
					Severity: "error",
				})
			}
		}

		// Validate allowed machines format
		for i, machine := range project.Isolation.AllowedMachines {
			if !isValidMachineName(machine) {
				result.Valid = false
				result.Errors = append(result.Errors, ValidationError{
					Field:    fmt.Sprintf("isolation.allowed_machines[%d]", i),
					Message:  fmt.Sprintf("invalid machine name format: %s", machine),
					Severity: "error",
				})
			}
		}

		// Validate allowed tags format
		for i, tag := range project.Isolation.AllowedTags {
			if !isValidTagName(tag) {
				result.Valid = false
				result.Errors = append(result.Errors, ValidationError{
					Field:    fmt.Sprintf("isolation.allowed_tags[%d]", i),
					Message:  fmt.Sprintf("invalid tag format: %s", tag),
					Severity: "error",
				})
			}
		}
	}

	return result
}

// GetIsolationSummary returns a summary of project isolation settings
func (ie *ProjectIsolationEngine) GetIsolationSummary(project *Project) map[string]interface{} {
	summary := map[string]interface{}{
		"enabled":         true,
		"facts_scope":     DefaultFactsScope,
		"variables_scope": DefaultVariablesScope,
		"machine_access":  DefaultMachineAccess,
	}

	if project.Isolation != nil {
		summary["enabled"] = project.Isolation.Enabled
		if project.Isolation.FactsScope != "" {
			summary["facts_scope"] = project.Isolation.FactsScope
		}
		if project.Isolation.VariablesScope != "" {
			summary["variables_scope"] = project.Isolation.VariablesScope
		}
		if project.Isolation.MachineAccess != "" {
			summary["machine_access"] = project.Isolation.MachineAccess
		}
		if len(project.Isolation.AllowedMachines) > 0 {
			summary["allowed_machines"] = project.Isolation.AllowedMachines
		}
		if len(project.Isolation.AllowedTags) > 0 {
			summary["allowed_tags"] = project.Isolation.AllowedTags
		}
	}

	return summary
}

// CheckMachineAccess checks if a machine is accessible for a project
func (ie *ProjectIsolationEngine) CheckMachineAccess(project *Project, machineName string, machineTags []string) bool {
	if project == nil || project.Isolation == nil {
		return true // Default to allowing access if no isolation settings
	}

	// If isolation is disabled, allow all access
	if !project.Isolation.Enabled {
		return true
	}

	// Check machine access level
	switch project.Isolation.MachineAccess {
	case "all":
		return true
	case "tagged":
		// Check if machine has any of the allowed tags
		if len(project.Isolation.AllowedTags) == 0 {
			return false
		}
		for _, allowedTag := range project.Isolation.AllowedTags {
			for _, machineTag := range machineTags {
				if allowedTag == machineTag {
					return true
				}
			}
		}
		return false
	case "explicit":
		// Check if machine is explicitly allowed
		for _, allowedMachine := range project.Isolation.AllowedMachines {
			if allowedMachine == machineName {
				return true
			}
		}
		// Also check if machine has any of the allowed tags
		for _, allowedTag := range project.Isolation.AllowedTags {
			for _, machineTag := range machineTags {
				if allowedTag == machineTag {
					return true
				}
			}
		}
		return false
	default:
		return true // Default to allowing access for unknown access levels
	}
}

// validateIsolationContext validates the isolation context
func (ie *ProjectIsolationEngine) validateIsolationContext(context *IsolationContext) error {
	if context == nil {
		return fmt.Errorf("isolation context cannot be nil")
	}

	if context.Project == nil {
		return fmt.Errorf("project cannot be nil in isolation context")
	}

	// Validate facts scope
	validFactsScopes := []string{"global", "project", "hybrid"}
	valid := false
	for _, scope := range validFactsScopes {
		if context.FactsScope == scope {
			valid = true
			break
		}
	}
	if !valid {
		return fmt.Errorf("invalid facts scope: %s", context.FactsScope)
	}

	// Validate variables scope
	validVariablesScopes := []string{"project", "inherited"}
	valid = false
	for _, scope := range validVariablesScopes {
		if context.VariablesScope == scope {
			valid = true
			break
		}
	}
	if !valid {
		return fmt.Errorf("invalid variables scope: %s", context.VariablesScope)
	}

	// Validate machine access
	validMachineAccess := []string{"all", "tagged", "explicit"}
	valid = false
	for _, access := range validMachineAccess {
		if context.MachineAccess == access {
			valid = true
			break
		}
	}
	if !valid {
		return fmt.Errorf("invalid machine access: %s", context.MachineAccess)
	}

	return nil
}

// isValidMachineName validates machine name format
func isValidMachineName(name string) bool {
	// Machine names must contain only alphanumeric characters, dots, underscores, and hyphens
	for _, char := range name {
		if !((char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9') || char == '.' || char == '_' || char == '-') {
			return false
		}
	}
	return len(name) > 0
}

// isValidTagName validates tag name format
func isValidTagName(tag string) bool {
	// Tags must contain only alphanumeric characters, dots, underscores, and hyphens
	for _, char := range tag {
		if !((char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9') || char == '.' || char == '_' || char == '-') {
			return false
		}
	}
	return len(tag) > 0
}
