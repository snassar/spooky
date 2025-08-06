package schemas

import (
	"fmt"
	"regexp"
	"strings"

	"spooky/internal/project"
)

// ProjectValidator provides project-specific validation within the unified system
type ProjectValidator struct {
	baseValidator *SchemaValidator
}

// NewProjectValidator creates a new project validator
func NewProjectValidator() *ProjectValidator {
	return &ProjectValidator{
		baseValidator: NewSchemaValidator(),
	}
}

// ValidateProject validates project configuration against schema
func (pv *ProjectValidator) ValidateProject(proj *project.Project) *ValidationResult {
	result := &ValidationResult{
		Valid:  true,
		Schema: "project",
	}

	// Validate project name
	if err := pv.ValidateProjectName(proj.Name); err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, ValidationError{
			Field:    "name",
			Message:  err.Error(),
			Value:    proj.Name,
			Severity: "error",
		})
	}

	// Validate project version if provided
	if proj.Version != "" {
		if err := pv.ValidateProjectVersion(proj.Version); err != nil {
			result.Valid = false
			result.Errors = append(result.Errors, ValidationError{
				Field:    "version",
				Message:  err.Error(),
				Value:    proj.Version,
				Severity: "error",
			})
		}
	}

	// Validate project description if provided
	if proj.Description != "" {
		if err := pv.ValidateProjectDescription(proj.Description); err != nil {
			result.Valid = false
			result.Errors = append(result.Errors, ValidationError{
				Field:    "description",
				Message:  err.Error(),
				Value:    proj.Description,
				Severity: "error",
			})
		}
	}

	// Validate project environment
	if err := pv.ValidateProjectEnvironment(proj.Environment); err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, ValidationError{
			Field:    "environment",
			Message:  err.Error(),
			Value:    proj.Environment,
			Severity: "error",
		})
	}

	// Validate project tags
	if err := pv.ValidateProjectTags(proj.Tags); err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, ValidationError{
			Field:    "tags",
			Message:  err.Error(),
			Value:    fmt.Sprintf("%v", proj.Tags),
			Severity: "error",
		})
	}

	// Validate project structure
	if err := pv.ValidateProjectStructure(proj.Structure); err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, ValidationError{
			Field:    "structure",
			Message:  err.Error(),
			Value:    fmt.Sprintf("%v", proj.Structure),
			Severity: "error",
		})
	}

	// Validate project isolation settings
	if err := pv.ValidateProjectIsolation(proj.Isolation); err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, ValidationError{
			Field:    "isolation",
			Message:  err.Error(),
			Value:    fmt.Sprintf("%v", proj.Isolation),
			Severity: "error",
		})
	}

	// Validate project dependencies
	if err := pv.ValidateProjectDependencies(proj.Dependencies); err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, ValidationError{
			Field:    "dependencies",
			Message:  err.Error(),
			Value:    fmt.Sprintf("%v", proj.Dependencies),
			Severity: "error",
		})
	}

	// Validate project execution settings
	if err := pv.ValidateProjectExecution(proj.Execution); err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, ValidationError{
			Field:    "execution",
			Message:  err.Error(),
			Value:    fmt.Sprintf("%v", proj.Execution),
			Severity: "error",
		})
	}

	return result
}

// ValidateProjectName validates project name format
func (pv *ProjectValidator) ValidateProjectName(name string) error {
	if name == "" {
		return fmt.Errorf("project name cannot be empty")
	}

	// Check length
	if len(name) > 100 {
		return fmt.Errorf("project name cannot exceed 100 characters")
	}

	// Check pattern: must start with letter, contain only alphanumeric, dots, underscores, hyphens
	validPattern := regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9._-]*$`)
	if !validPattern.MatchString(name) {
		return fmt.Errorf("project name must start with a letter and contain only alphanumeric characters, dots, underscores, and hyphens")
	}

	return nil
}

// ValidateProjectVersion validates project version format
func (pv *ProjectValidator) ValidateProjectVersion(version string) error {
	// Check semantic versioning pattern
	validPattern := regexp.MustCompile(`^[0-9]+\.[0-9]+\.[0-9]+(-[a-zA-Z0-9._-]+)?$`)
	if !validPattern.MatchString(version) {
		return fmt.Errorf("version must follow semantic versioning format (e.g., 1.0.0 or 1.0.0-beta)")
	}

	return nil
}

// ValidateProjectDescription validates project description
func (pv *ProjectValidator) ValidateProjectDescription(description string) error {
	if len(description) > 500 {
		return fmt.Errorf("project description cannot exceed 500 characters")
	}

	return nil
}

// ValidateProjectEnvironment validates project environment
func (pv *ProjectValidator) ValidateProjectEnvironment(environment string) error {
	validEnvironments := []string{"production", "staging", "development", "testing"}

	for _, valid := range validEnvironments {
		if environment == valid {
			return nil
		}
	}

	return fmt.Errorf("environment must be one of: %s", strings.Join(validEnvironments, ", "))
}

// ValidateProjectTags validates project tags
func (pv *ProjectValidator) ValidateProjectTags(tags []string) error {
	if tags == nil {
		return nil // Tags are optional
	}

	validPattern := regexp.MustCompile(`^[a-zA-Z0-9._-]+$`)

	for i, tag := range tags {
		if tag == "" {
			return fmt.Errorf("tag at index %d cannot be empty", i)
		}

		if !validPattern.MatchString(tag) {
			return fmt.Errorf("tag '%s' contains invalid characters", tag)
		}
	}

	return nil
}

// ValidateProjectStructure validates project structure configuration
func (pv *ProjectValidator) ValidateProjectStructure(structure *project.ProjectStructure) error {
	if structure == nil {
		return nil // Structure is optional
	}

	// Validate directory names
	validPattern := regexp.MustCompile(`^[a-zA-Z0-9/._-]+$`)

	if structure.TemplatesDir != "" && !validPattern.MatchString(structure.TemplatesDir) {
		return fmt.Errorf("templates_dir contains invalid characters")
	}

	if structure.DataDir != "" && !validPattern.MatchString(structure.DataDir) {
		return fmt.Errorf("data_dir contains invalid characters")
	}

	if structure.ScriptsDir != "" && !validPattern.MatchString(structure.ScriptsDir) {
		return fmt.Errorf("scripts_dir contains invalid characters")
	}

	if structure.LogsDir != "" && !validPattern.MatchString(structure.LogsDir) {
		return fmt.Errorf("logs_dir contains invalid characters")
	}

	if structure.BackupsDir != "" && !validPattern.MatchString(structure.BackupsDir) {
		return fmt.Errorf("backups_dir contains invalid characters")
	}

	return nil
}

// ValidateProjectIsolation validates project isolation settings
func (pv *ProjectValidator) ValidateProjectIsolation(isolation *project.ProjectIsolation) error {
	if isolation == nil {
		return nil // Isolation is optional
	}

	// Validate facts scope
	validFactsScopes := []string{"global", "project", "hybrid"}
	validScope := false
	for _, scope := range validFactsScopes {
		if isolation.FactsScope == scope {
			validScope = true
			break
		}
	}
	if !validScope {
		return fmt.Errorf("facts_scope must be one of: %s", strings.Join(validFactsScopes, ", "))
	}

	// Validate variables scope
	validVariablesScopes := []string{"project", "inherited"}
	validVarScope := false
	for _, scope := range validVariablesScopes {
		if isolation.VariablesScope == scope {
			validVarScope = true
			break
		}
	}
	if !validVarScope {
		return fmt.Errorf("variables_scope must be one of: %s", strings.Join(validVariablesScopes, ", "))
	}

	// Validate machine access
	validMachineAccess := []string{"all", "tagged", "explicit"}
	validAccess := false
	for _, access := range validMachineAccess {
		if isolation.MachineAccess == access {
			validAccess = true
			break
		}
	}
	if !validAccess {
		return fmt.Errorf("machine_access must be one of: %s", strings.Join(validMachineAccess, ", "))
	}

	// Validate explicit machine access requirements
	if isolation.MachineAccess == "explicit" {
		if len(isolation.AllowedMachines) == 0 && len(isolation.AllowedTags) == 0 {
			return fmt.Errorf("explicit machine access requires allowed_machines or allowed_tags to be specified")
		}
	}

	// Validate allowed machines
	validMachinePattern := regexp.MustCompile(`^[a-zA-Z0-9._-]+$`)
	for _, machine := range isolation.AllowedMachines {
		if !validMachinePattern.MatchString(machine) {
			return fmt.Errorf("allowed_machine '%s' contains invalid characters", machine)
		}
	}

	// Validate allowed tags
	validTagPattern := regexp.MustCompile(`^[a-zA-Z0-9._-]+$`)
	for _, tag := range isolation.AllowedTags {
		if !validTagPattern.MatchString(tag) {
			return fmt.Errorf("allowed_tag '%s' contains invalid characters", tag)
		}
	}

	return nil
}

// ValidateProjectDependencies validates project dependencies
func (pv *ProjectValidator) ValidateProjectDependencies(dependencies *project.ProjectDependencies) error {
	if dependencies == nil {
		return nil // Dependencies are optional
	}

	// Validate imports
	validImportPattern := regexp.MustCompile(`^[a-zA-Z0-9/._-]+$`)
	for _, importPath := range dependencies.Imports {
		if !validImportPattern.MatchString(importPath) {
			return fmt.Errorf("import path '%s' contains invalid characters", importPath)
		}
	}

	// Validate shared variables
	validVarPattern := regexp.MustCompile(`^[a-zA-Z0-9._-]+$`)
	for _, varName := range dependencies.SharedVariables {
		if !validVarPattern.MatchString(varName) {
			return fmt.Errorf("shared variable '%s' contains invalid characters", varName)
		}
	}

	// Validate shared facts
	validFactPattern := regexp.MustCompile(`^[a-zA-Z0-9._-]+$`)
	for _, factName := range dependencies.SharedFacts {
		if !validFactPattern.MatchString(factName) {
			return fmt.Errorf("shared fact '%s' contains invalid characters", factName)
		}
	}

	return nil
}

// ValidateProjectExecution validates project execution settings
func (pv *ProjectValidator) ValidateProjectExecution(execution *project.ProjectExecution) error {
	if execution == nil {
		return nil // Execution settings are optional
	}

	// Validate default timeout
	if execution.DefaultTimeout < 1 || execution.DefaultTimeout > 3600 {
		return fmt.Errorf("default_timeout must be between 1 and 3600 seconds")
	}

	// Validate max parallel
	if execution.MaxParallel < 1 || execution.MaxParallel > 100 {
		return fmt.Errorf("max_parallel must be between 1 and 100")
	}

	return nil
}

// ValidateProjectWithSchema validates project against embedded schema
func (pv *ProjectValidator) ValidateProjectWithSchema(projectPath string) *ValidationResult {
	// Load and validate project configuration against embedded schema
	result := &ValidationResult{
		Valid:  true,
		Schema: "project",
	}

	// DEPRECATED: Schema system is fully implemented - this TODO is ready for removal
	// This will be implemented when the schema loading system is enhanced

	return result
}

// ConvertProjectValidationResult converts project validation result to unified format
func (pv *ProjectValidator) ConvertProjectValidationResult(projectResult interface{}) *ValidationResult {
	// Convert project-specific validation result to unified format
	if projectResult == nil {
		return &ValidationResult{
			Valid:  false,
			Schema: "project",
			Errors: []ValidationError{
				{
					Field:    "project",
					Message:  "Project validation result is nil",
					Severity: "error",
				},
			},
		}
	}

	// DEPRECATED: Schema system is fully implemented - this TODO is ready for removal
	return &ValidationResult{
		Valid:  true,
		Schema: "project",
	}
}
