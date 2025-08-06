package project

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// ProjectPortabilityEngine manages project portability across environments
type ProjectPortabilityEngine struct{}

// NewProjectPortabilityEngine creates a new project portability engine
func NewProjectPortabilityEngine() *ProjectPortabilityEngine {
	return &ProjectPortabilityEngine{}
}

// PortabilityContext represents the portability context for a project
type PortabilityContext struct {
	Project     *Project
	Environment string
	Region      string
	Aliases     map[string]string
	Portable    bool
	Validation  *PortabilityValidation
}

// PortabilityValidation represents portability validation results
type PortabilityValidation struct {
	Valid    bool
	Errors   []ValidationError
	Warnings []ValidationError
	Aliases  map[string]string
}

// EnvironmentMapping represents environment-specific mappings
type EnvironmentMapping struct {
	Environment string
	Aliases     map[string]string
	Rules       []PortabilityRule
}

// PortabilityRule represents a portability rule
type PortabilityRule struct {
	Name        string
	Pattern     string
	Replacement string
	Condition   string
}

// ValidateProjectPortability validates project portability across environments
func (pe *ProjectPortabilityEngine) ValidateProjectPortability(project *Project) *PortabilityValidation {
	result := &PortabilityValidation{
		Valid:   true,
		Aliases: make(map[string]string),
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

	// Validate project name format
	if err := pe.validateProjectNameFormat(project.Name); err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, ValidationError{
			Field:    "name",
			Message:  err.Error(),
			Severity: "error",
		})
	}

	// Validate environment-specific naming
	if project.Environment != "" {
		if err := pe.validateEnvironmentName(project.Environment); err != nil {
			result.Valid = false
			result.Errors = append(result.Errors, ValidationError{
				Field:    "environment",
				Message:  err.Error(),
				Severity: "error",
			})
		}
	}

	// Generate environment aliases
	aliases := pe.generateEnvironmentAliases(project)
	result.Aliases = aliases

	// Check for portability issues
	warnings := pe.checkPortabilityWarnings(project)
	result.Warnings = warnings

	return result
}

// GenerateEnvironmentAliases generates environment-specific aliases for a project
func (pe *ProjectPortabilityEngine) GenerateEnvironmentAliases(project *Project) map[string]string {
	return pe.generateEnvironmentAliases(project)
}

// GetPortableName returns the portable name for a project in a specific environment
func (pe *ProjectPortabilityEngine) GetPortableName(project *Project, environment string) (string, error) {
	if project == nil {
		return "", fmt.Errorf("project cannot be nil")
	}

	if err := pe.validateEnvironmentName(environment); err != nil {
		return "", fmt.Errorf("invalid environment: %w", err)
	}

	// Generate environment-specific name
	portableName := pe.applyEnvironmentNaming(project.Name, environment)

	return portableName, nil
}

// CreatePortabilityContext creates a portability context for a project
func (pe *ProjectPortabilityEngine) CreatePortabilityContext(project *Project, environment string) (*PortabilityContext, error) {
	if project == nil {
		return nil, fmt.Errorf("project cannot be nil")
	}

	// Validate environment
	if err := pe.validateEnvironmentName(environment); err != nil {
		return nil, fmt.Errorf("invalid environment: %w", err)
	}

	// Generate aliases
	aliases := pe.generateEnvironmentAliases(project)

	// Validate portability
	validation := pe.ValidateProjectPortability(project)

	context := &PortabilityContext{
		Project:     project,
		Environment: environment,
		Region:      project.Region,
		Aliases:     aliases,
		Portable:    validation.Valid,
		Validation:  validation,
	}

	return context, nil
}

// MigrateProject migrates a project to a different environment
func (pe *ProjectPortabilityEngine) MigrateProject(project *Project, targetEnvironment string) (*Project, error) {
	if project == nil {
		return nil, fmt.Errorf("project cannot be nil")
	}

	// Validate target environment
	if err := pe.validateEnvironmentName(targetEnvironment); err != nil {
		return nil, fmt.Errorf("invalid target environment: %w", err)
	}

	// Create a copy of the project
	migratedProject := *project

	// Update environment
	migratedProject.Environment = targetEnvironment

	// Generate new portable name
	newName, err := pe.GetPortableName(&migratedProject, targetEnvironment)
	if err != nil {
		return nil, fmt.Errorf("failed to generate portable name: %w", err)
	}
	migratedProject.Name = newName

	// Update project path if needed
	if project.Path != "" {
		projectDir := filepath.Base(project.Path)
		newProjectDir := pe.applyEnvironmentNaming(projectDir, targetEnvironment)
		migratedProject.Path = filepath.Join(filepath.Dir(project.Path), newProjectDir)
	}

	// Validate migrated project
	validation := pe.ValidateProjectPortability(&migratedProject)
	if !validation.Valid {
		return nil, fmt.Errorf("migrated project validation failed: %v", validation.Errors)
	}

	return &migratedProject, nil
}

// GetEnvironmentMappings returns environment-specific mappings
func (pe *ProjectPortabilityEngine) GetEnvironmentMappings() map[string]*EnvironmentMapping {
	mappings := make(map[string]*EnvironmentMapping)

	// Production environment
	mappings["production"] = &EnvironmentMapping{
		Environment: "production",
		Aliases: map[string]string{
			"prod": "production",
			"live": "production",
		},
		Rules: []PortabilityRule{
			{
				Name:        "remove_dev_suffix",
				Pattern:     "-dev$",
				Replacement: "",
				Condition:   "environment == production",
			},
			{
				Name:        "add_prod_suffix",
				Pattern:     "$",
				Replacement: "-prod",
				Condition:   "environment == production && !contains(name, '-prod')",
			},
		},
	}

	// Staging environment
	mappings["staging"] = &EnvironmentMapping{
		Environment: "staging",
		Aliases: map[string]string{
			"stage":   "staging",
			"preprod": "staging",
		},
		Rules: []PortabilityRule{
			{
				Name:        "remove_dev_suffix",
				Pattern:     "-dev$",
				Replacement: "",
				Condition:   "environment == staging",
			},
			{
				Name:        "add_staging_suffix",
				Pattern:     "$",
				Replacement: "-staging",
				Condition:   "environment == staging && !contains(name, '-staging')",
			},
		},
	}

	// Development environment
	mappings["development"] = &EnvironmentMapping{
		Environment: "development",
		Aliases: map[string]string{
			"dev":   "development",
			"local": "development",
		},
		Rules: []PortabilityRule{
			{
				Name:        "add_dev_suffix",
				Pattern:     "$",
				Replacement: "-dev",
				Condition:   "environment == development && !contains(name, '-dev')",
			},
		},
	}

	// Testing environment
	mappings["testing"] = &EnvironmentMapping{
		Environment: "testing",
		Aliases: map[string]string{
			"test": "testing",
			"qa":   "testing",
		},
		Rules: []PortabilityRule{
			{
				Name:        "add_test_suffix",
				Pattern:     "$",
				Replacement: "-test",
				Condition:   "environment == testing && !contains(name, '-test')",
			},
		},
	}

	return mappings
}

// validateProjectNameFormat validates project name format
func (pe *ProjectPortabilityEngine) validateProjectNameFormat(name string) error {
	if name == "" {
		return fmt.Errorf("project name cannot be empty")
	}

	// Check length
	if len(name) > 100 {
		return fmt.Errorf("project name too long (max 100 characters)")
	}

	// Check format (alphanumeric, dots, underscores, hyphens)
	validNamePattern := regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9._-]*$`)
	if !validNamePattern.MatchString(name) {
		return fmt.Errorf("project name must start with a letter and contain only alphanumeric characters, dots, underscores, and hyphens")
	}

	// Check for reserved names
	reservedNames := []string{"spooky", "config", "system", "admin", "root"}
	for _, reserved := range reservedNames {
		if strings.EqualFold(name, reserved) {
			return fmt.Errorf("project name '%s' is reserved", name)
		}
	}

	return nil
}

// validateEnvironmentName validates environment name
func (pe *ProjectPortabilityEngine) validateEnvironmentName(environment string) error {
	if environment == "" {
		return fmt.Errorf("environment cannot be empty")
	}

	validEnvironments := []string{"production", "staging", "development", "testing"}
	for _, valid := range validEnvironments {
		if environment == valid {
			return nil
		}
	}

	return fmt.Errorf("invalid environment '%s' (must be one of: %s)", environment, strings.Join(validEnvironments, ", "))
}

// generateEnvironmentAliases generates environment-specific aliases
func (pe *ProjectPortabilityEngine) generateEnvironmentAliases(project *Project) map[string]string {
	aliases := make(map[string]string)

	if project == nil || project.Name == "" {
		return aliases
	}

	mappings := pe.GetEnvironmentMappings()

	// Generate aliases for each environment
	for envName, mapping := range mappings {
		portableName := pe.applyEnvironmentNaming(project.Name, envName)
		aliases[envName] = portableName

		// Add environment aliases
		for alias, target := range mapping.Aliases {
			if target == envName {
				aliases[alias] = portableName
			}
		}
	}

	return aliases
}

// applyEnvironmentNaming applies environment-specific naming rules
func (pe *ProjectPortabilityEngine) applyEnvironmentNaming(name string, environment string) string {
	if name == "" || environment == "" {
		return name
	}

	mappings := pe.GetEnvironmentMappings()
	mapping, exists := mappings[environment]
	if !exists {
		return name
	}

	result := name

	// Apply rules in order
	for _, rule := range mapping.Rules {
		// Simple rule application (in a real implementation, this would be more sophisticated)
		if strings.Contains(rule.Condition, "environment == "+environment) {
			re := regexp.MustCompile(rule.Pattern)
			result = re.ReplaceAllString(result, rule.Replacement)
		}
	}

	return result
}

// checkPortabilityWarnings checks for portability warnings
func (pe *ProjectPortabilityEngine) checkPortabilityWarnings(project *Project) []ValidationError {
	var warnings []ValidationError

	if project == nil {
		return warnings
	}

	// Check for environment-specific issues
	if project.Environment == "production" {
		if strings.Contains(project.Name, "-dev") {
			warnings = append(warnings, ValidationError{
				Field:    "name",
				Message:  "production project contains '-dev' suffix",
				Severity: "warning",
			})
		}
	}

	// Check for region-specific issues
	if project.Region != "" {
		// Could add region-specific validation here
	}

	// Check for version issues
	if project.Version == "" {
		warnings = append(warnings, ValidationError{
			Field:    "version",
			Message:  "project version not specified (recommended for portability)",
			Severity: "warning",
		})
	}

	return warnings
}

// IsPortable checks if a project is portable across environments
func (pe *ProjectPortabilityEngine) IsPortable(project *Project) bool {
	if project == nil {
		return false
	}

	validation := pe.ValidateProjectPortability(project)
	return validation.Valid
}

// GetPortabilityIssues returns a list of portability issues for a project
func (pe *ProjectPortabilityEngine) GetPortabilityIssues(project *Project) []string {
	var issues []string

	if project == nil {
		return append(issues, "project is nil")
	}

	validation := pe.ValidateProjectPortability(project)

	// Add errors
	for _, err := range validation.Errors {
		issues = append(issues, fmt.Sprintf("ERROR: %s - %s", err.Field, err.Message))
	}

	// Add warnings
	for _, warning := range validation.Warnings {
		issues = append(issues, fmt.Sprintf("WARNING: %s - %s", warning.Field, warning.Message))
	}

	return issues
}

// CreatePortableProject creates a new project with portable naming
func (pe *ProjectPortabilityEngine) CreatePortableProject(name string, environment string, region string) (*Project, error) {
	// Validate inputs
	if err := pe.validateProjectNameFormat(name); err != nil {
		return nil, fmt.Errorf("invalid project name: %w", err)
	}

	if err := pe.validateEnvironmentName(environment); err != nil {
		return nil, fmt.Errorf("invalid environment: %w", err)
	}

	// Create portable name
	portableName := pe.applyEnvironmentNaming(name, environment)

	// Create project
	project := &Project{
		Name:        portableName,
		Environment: environment,
		Region:      region,
	}

	// Set defaults
	project.SetDefaults()

	// Validate portability
	validation := pe.ValidateProjectPortability(project)
	if !validation.Valid {
		return nil, fmt.Errorf("created project is not portable: %v", validation.Errors)
	}

	return project, nil
}

// GetEnvironmentFromPath attempts to determine environment from project path
func (pe *ProjectPortabilityEngine) GetEnvironmentFromPath(projectPath string) string {
	if projectPath == "" {
		return ""
	}

	// Check for environment indicators in path
	pathLower := strings.ToLower(projectPath)

	if strings.Contains(pathLower, "/prod/") || strings.Contains(pathLower, "-prod") {
		return "production"
	}
	if strings.Contains(pathLower, "/staging/") || strings.Contains(pathLower, "-staging") {
		return "staging"
	}
	if strings.Contains(pathLower, "/dev/") || strings.Contains(pathLower, "-dev") {
		return "development"
	}
	if strings.Contains(pathLower, "/test/") || strings.Contains(pathLower, "-test") {
		return "testing"
	}

	// Check environment variable
	if env := os.Getenv("SPOOKY_ENVIRONMENT"); env != "" {
		return env
	}

	// Default to development
	return "development"
}
