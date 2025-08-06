package project

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// ValidationError represents a validation error
type ValidationError struct {
	Field    string `json:"field"`
	Message  string `json:"message"`
	Value    string `json:"value,omitempty"`
	Severity string `json:"severity"`
}

// ValidationResult represents the result of validation
type ValidationResult struct {
	Valid    bool              `json:"valid"`
	Errors   []ValidationError `json:"errors,omitempty"`
	Warnings []ValidationError `json:"warnings,omitempty"`
}

// ProjectValidator provides project validation functionality
type ProjectValidator struct{}

// NewProjectValidator creates a new project validator
func NewProjectValidator() *ProjectValidator {
	return &ProjectValidator{}
}

// ValidateProject validates a project configuration
func (pv *ProjectValidator) ValidateProject(project *Project) *ValidationResult {
	result := &ValidationResult{
		Valid: true,
	}

	// Validate required fields
	if err := pv.validateRequiredFields(project); err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, *err)
	}

	// Validate project name
	if err := pv.validateProjectName(project.Name); err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, *err)
	}

	// Validate project version if provided
	if project.Version != "" {
		if err := pv.validateProjectVersion(project.Version); err != nil {
			result.Valid = false
			result.Errors = append(result.Errors, *err)
		}
	}

	// Validate project description if provided
	if project.Description != "" {
		if err := pv.validateProjectDescription(project.Description); err != nil {
			result.Valid = false
			result.Errors = append(result.Errors, *err)
		}
	}

	// Validate project environment
	if err := pv.validateProjectEnvironment(project.Environment); err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, *err)
	}

	// Validate project tags
	if err := pv.validateProjectTags(project.Tags); err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, *err)
	}

	// Validate project structure
	if err := pv.validateProjectStructure(project.Structure); err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, *err)
	}

	// Validate project isolation settings
	if err := pv.validateProjectIsolation(project.Isolation); err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, *err)
	}

	// Validate project dependencies
	if err := pv.validateProjectDependencies(project.Dependencies); err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, *err)
	}

	// Validate project execution settings
	if err := pv.validateProjectExecution(project.Execution); err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, *err)
	}

	// Validate project contact information
	if project.Contact != nil {
		if err := pv.validateProjectContact(project.Contact); err != nil {
			result.Valid = false
			result.Errors = append(result.Errors, *err)
		}
	}

	return result
}

// ValidateProjectPath validates project directory structure
func (pv *ProjectValidator) ValidateProjectPath(projectPath string) *ValidationResult {
	result := &ValidationResult{
		Valid: true,
	}

	// Check if project path exists
	if _, err := os.Stat(projectPath); os.IsNotExist(err) {
		result.Valid = false
		result.Errors = append(result.Errors, ValidationError{
			Field:    "path",
			Message:  "Project path does not exist",
			Value:    projectPath,
			Severity: "error",
		})
		return result
	}

	// Check if project.hcl exists
	projectHCLPath := filepath.Join(projectPath, "project.hcl")
	if _, err := os.Stat(projectHCLPath); os.IsNotExist(err) {
		result.Valid = false
		result.Errors = append(result.Errors, ValidationError{
			Field:    "project.hcl",
			Message:  "Project configuration file not found",
			Value:    projectHCLPath,
			Severity: "error",
		})
	}

	// Check if facts.db directory exists
	factsDBPath := filepath.Join(projectPath, "facts.db")
	if _, err := os.Stat(factsDBPath); os.IsNotExist(err) {
		result.Warnings = append(result.Warnings, ValidationError{
			Field:    "facts.db",
			Message:  "Facts database directory not found (will be created during initialization)",
			Value:    factsDBPath,
			Severity: "warning",
		})
	}

	// Check for optional directories
	optionalDirs := []string{"templates", "data", "scripts", "logs", "backups", "machines", "actions", "variables"}
	for _, dir := range optionalDirs {
		dirPath := filepath.Join(projectPath, dir)
		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			result.Warnings = append(result.Warnings, ValidationError{
				Field:    dir,
				Message:  fmt.Sprintf("Optional directory '%s' not found", dir),
				Value:    dirPath,
				Severity: "warning",
			})
		}
	}

	return result
}

// validateRequiredFields validates required project fields
func (pv *ProjectValidator) validateRequiredFields(project *Project) *ValidationError {
	if project.Name == "" {
		return &ValidationError{
			Field:    "name",
			Message:  "Project name is required",
			Severity: "error",
		}
	}

	if project.Path == "" {
		return &ValidationError{
			Field:    "path",
			Message:  "Project path is required",
			Severity: "error",
		}
	}

	return nil
}

// validateProjectName validates project name format
func (pv *ProjectValidator) validateProjectName(name string) *ValidationError {
	if name == "" {
		return &ValidationError{
			Field:    "name",
			Message:  "Project name cannot be empty",
			Severity: "error",
		}
	}

	if len(name) > 100 {
		return &ValidationError{
			Field:    "name",
			Message:  "Project name cannot exceed 100 characters",
			Value:    name,
			Severity: "error",
		}
	}

	// Check pattern: must start with letter, contain only alphanumeric, dots, underscores, hyphens
	validPattern := regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9._-]*$`)
	if !validPattern.MatchString(name) {
		return &ValidationError{
			Field:    "name",
			Message:  "Project name must start with a letter and contain only alphanumeric characters, dots, underscores, and hyphens",
			Value:    name,
			Severity: "error",
		}
	}

	return nil
}

// validateProjectVersion validates project version format
func (pv *ProjectValidator) validateProjectVersion(version string) *ValidationError {
	// Check semantic versioning pattern
	validPattern := regexp.MustCompile(`^[0-9]+\.[0-9]+\.[0-9]+(-[a-zA-Z0-9._-]+)?$`)
	if !validPattern.MatchString(version) {
		return &ValidationError{
			Field:    "version",
			Message:  "Version must follow semantic versioning format (e.g., 1.0.0 or 1.0.0-beta)",
			Value:    version,
			Severity: "error",
		}
	}

	return nil
}

// validateProjectDescription validates project description
func (pv *ProjectValidator) validateProjectDescription(description string) *ValidationError {
	if len(description) > 500 {
		return &ValidationError{
			Field:    "description",
			Message:  "Project description cannot exceed 500 characters",
			Value:    description,
			Severity: "error",
		}
	}

	return nil
}

// validateProjectEnvironment validates project environment
func (pv *ProjectValidator) validateProjectEnvironment(environment string) *ValidationError {
	validEnvironments := []string{"production", "staging", "development", "testing"}

	for _, valid := range validEnvironments {
		if environment == valid {
			return nil
		}
	}

	return &ValidationError{
		Field:    "environment",
		Message:  fmt.Sprintf("Environment must be one of: %s", strings.Join(validEnvironments, ", ")),
		Value:    environment,
		Severity: "error",
	}
}

// validateProjectTags validates project tags
func (pv *ProjectValidator) validateProjectTags(tags []string) *ValidationError {
	if tags == nil {
		return nil // Tags are optional
	}

	validPattern := regexp.MustCompile(`^[a-zA-Z0-9._-]+$`)

	for i, tag := range tags {
		if tag == "" {
			return &ValidationError{
				Field:    fmt.Sprintf("tags[%d]", i),
				Message:  "Tag cannot be empty",
				Severity: "error",
			}
		}

		if !validPattern.MatchString(tag) {
			return &ValidationError{
				Field:    fmt.Sprintf("tags[%d]", i),
				Message:  "Tag contains invalid characters",
				Value:    tag,
				Severity: "error",
			}
		}
	}

	return nil
}

// validateProjectStructure validates project structure configuration
func (pv *ProjectValidator) validateProjectStructure(structure *ProjectStructure) *ValidationError {
	if structure == nil {
		return nil // Structure is optional
	}

	// Validate directory names
	validPattern := regexp.MustCompile(`^[a-zA-Z0-9/._-]+$`)

	if structure.TemplatesDir != "" && !validPattern.MatchString(structure.TemplatesDir) {
		return &ValidationError{
			Field:    "structure.templates_dir",
			Message:  "Templates directory contains invalid characters",
			Value:    structure.TemplatesDir,
			Severity: "error",
		}
	}

	if structure.DataDir != "" && !validPattern.MatchString(structure.DataDir) {
		return &ValidationError{
			Field:    "structure.data_dir",
			Message:  "Data directory contains invalid characters",
			Value:    structure.DataDir,
			Severity: "error",
		}
	}

	if structure.ScriptsDir != "" && !validPattern.MatchString(structure.ScriptsDir) {
		return &ValidationError{
			Field:    "structure.scripts_dir",
			Message:  "Scripts directory contains invalid characters",
			Value:    structure.ScriptsDir,
			Severity: "error",
		}
	}

	if structure.LogsDir != "" && !validPattern.MatchString(structure.LogsDir) {
		return &ValidationError{
			Field:    "structure.logs_dir",
			Message:  "Logs directory contains invalid characters",
			Value:    structure.LogsDir,
			Severity: "error",
		}
	}

	if structure.BackupsDir != "" && !validPattern.MatchString(structure.BackupsDir) {
		return &ValidationError{
			Field:    "structure.backups_dir",
			Message:  "Backups directory contains invalid characters",
			Value:    structure.BackupsDir,
			Severity: "error",
		}
	}

	return nil
}

// validateProjectIsolation validates project isolation settings
func (pv *ProjectValidator) validateProjectIsolation(isolation *ProjectIsolation) *ValidationError {
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
		return &ValidationError{
			Field:    "isolation.facts_scope",
			Message:  fmt.Sprintf("Facts scope must be one of: %s", strings.Join(validFactsScopes, ", ")),
			Value:    isolation.FactsScope,
			Severity: "error",
		}
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
		return &ValidationError{
			Field:    "isolation.variables_scope",
			Message:  fmt.Sprintf("Variables scope must be one of: %s", strings.Join(validVariablesScopes, ", ")),
			Value:    isolation.VariablesScope,
			Severity: "error",
		}
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
		return &ValidationError{
			Field:    "isolation.machine_access",
			Message:  fmt.Sprintf("Machine access must be one of: %s", strings.Join(validMachineAccess, ", ")),
			Value:    isolation.MachineAccess,
			Severity: "error",
		}
	}

	// Validate explicit machine access requirements
	if isolation.MachineAccess == "explicit" {
		if len(isolation.AllowedMachines) == 0 && len(isolation.AllowedTags) == 0 {
			return &ValidationError{
				Field:    "isolation",
				Message:  "Explicit machine access requires allowed_machines or allowed_tags to be specified",
				Severity: "error",
			}
		}
	}

	// Validate allowed machines
	validMachinePattern := regexp.MustCompile(`^[a-zA-Z0-9._-]+$`)
	for i, machine := range isolation.AllowedMachines {
		if !validMachinePattern.MatchString(machine) {
			return &ValidationError{
				Field:    fmt.Sprintf("isolation.allowed_machines[%d]", i),
				Message:  "Allowed machine contains invalid characters",
				Value:    machine,
				Severity: "error",
			}
		}
	}

	// Validate allowed tags
	validTagPattern := regexp.MustCompile(`^[a-zA-Z0-9._-]+$`)
	for i, tag := range isolation.AllowedTags {
		if !validTagPattern.MatchString(tag) {
			return &ValidationError{
				Field:    fmt.Sprintf("isolation.allowed_tags[%d]", i),
				Message:  "Allowed tag contains invalid characters",
				Value:    tag,
				Severity: "error",
			}
		}
	}

	return nil
}

// validateProjectDependencies validates project dependencies
func (pv *ProjectValidator) validateProjectDependencies(dependencies *ProjectDependencies) *ValidationError {
	if dependencies == nil {
		return nil // Dependencies are optional
	}

	// Validate imports
	validImportPattern := regexp.MustCompile(`^[a-zA-Z0-9/._-]+$`)
	for i, importPath := range dependencies.Imports {
		if !validImportPattern.MatchString(importPath) {
			return &ValidationError{
				Field:    fmt.Sprintf("dependencies.imports[%d]", i),
				Message:  "Import path contains invalid characters",
				Value:    importPath,
				Severity: "error",
			}
		}
	}

	// Validate shared variables
	validVarPattern := regexp.MustCompile(`^[a-zA-Z0-9._-]+$`)
	for i, varName := range dependencies.SharedVariables {
		if !validVarPattern.MatchString(varName) {
			return &ValidationError{
				Field:    fmt.Sprintf("dependencies.shared_variables[%d]", i),
				Message:  "Shared variable contains invalid characters",
				Value:    varName,
				Severity: "error",
			}
		}
	}

	// Validate shared facts
	validFactPattern := regexp.MustCompile(`^[a-zA-Z0-9._-]+$`)
	for i, factName := range dependencies.SharedFacts {
		if !validFactPattern.MatchString(factName) {
			return &ValidationError{
				Field:    fmt.Sprintf("dependencies.shared_facts[%d]", i),
				Message:  "Shared fact contains invalid characters",
				Value:    factName,
				Severity: "error",
			}
		}
	}

	return nil
}

// validateProjectExecution validates project execution settings
func (pv *ProjectValidator) validateProjectExecution(execution *ProjectExecution) *ValidationError {
	if execution == nil {
		return nil // Execution settings are optional
	}

	// Validate default timeout
	if execution.DefaultTimeout < 1 || execution.DefaultTimeout > 3600 {
		return &ValidationError{
			Field:    "execution.default_timeout",
			Message:  "Default timeout must be between 1 and 3600 seconds",
			Value:    fmt.Sprintf("%d", execution.DefaultTimeout),
			Severity: "error",
		}
	}

	// Validate max parallel
	if execution.MaxParallel < 1 || execution.MaxParallel > 100 {
		return &ValidationError{
			Field:    "execution.max_parallel",
			Message:  "Maximum parallel executions must be between 1 and 100",
			Value:    fmt.Sprintf("%d", execution.MaxParallel),
			Severity: "error",
		}
	}

	return nil
}

// validateProjectContact validates project contact information
func (pv *ProjectValidator) validateProjectContact(contact *ProjectContact) *ValidationError {
	// Validate email format if provided
	if contact.Email != "" {
		emailPattern := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
		if !emailPattern.MatchString(contact.Email) {
			return &ValidationError{
				Field:    "contact.email",
				Message:  "Contact email must be a valid email address",
				Value:    contact.Email,
				Severity: "error",
			}
		}
	}

	// Validate URL format if provided
	if contact.URL != "" {
		urlPattern := regexp.MustCompile(`^https?://[^\s/$.?#].[^\s]*$`)
		if !urlPattern.MatchString(contact.URL) {
			return &ValidationError{
				Field:    "contact.url",
				Message:  "Contact URL must be a valid URI",
				Value:    contact.URL,
				Severity: "error",
			}
		}
	}

	return nil
}
