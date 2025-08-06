package project

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// ProjectIdentityManager manages project identity and naming
type ProjectIdentityManager struct {
	validator *ProjectValidator
}

// NewProjectIdentityManager creates a new project identity manager
func NewProjectIdentityManager() *ProjectIdentityManager {
	return &ProjectIdentityManager{
		validator: NewProjectValidator(),
	}
}

// ValidateProjectName validates project name format and uniqueness
func (pim *ProjectIdentityManager) ValidateProjectName(name string, projectPath string) *ValidationResult {
	result := &ValidationResult{
		Valid: true,
	}

	// Validate name format
	if err := pim.validateNameFormat(name); err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, *err)
	}

	// Check for uniqueness in the same directory
	if err := pim.validateNameUniqueness(name, projectPath); err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, *err)
	}

	return result
}

// GenerateProjectName generates a project name based on directory name
func (pim *ProjectIdentityManager) GenerateProjectName(projectPath string) (string, error) {
	// Get the directory name
	dirName := filepath.Base(projectPath)

	// Clean the directory name to make it a valid project name
	cleanName := pim.cleanProjectName(dirName)

	// Validate the generated name
	if err := pim.validateNameFormat(cleanName); err != nil {
		return "", fmt.Errorf("generated project name is invalid: %s", err.Message)
	}

	return cleanName, nil
}

// ValidateEnvironmentAwareName validates environment-aware project naming
func (pim *ProjectIdentityManager) ValidateEnvironmentAwareName(name string, environment string) *ValidationResult {
	result := &ValidationResult{
		Valid: true,
	}

	// Check if name contains environment-specific patterns
	if pim.hasEnvironmentPattern(name) {
		// Validate that the environment in the name matches the specified environment
		if !pim.matchesEnvironment(name, environment) {
			result.Valid = false
			result.Errors = append(result.Errors, ValidationError{
				Field:    "name",
				Message:  fmt.Sprintf("Project name contains environment pattern that doesn't match specified environment '%s'", environment),
				Value:    name,
				Severity: "error",
			})
		}
	}

	return result
}

// GenerateEnvironmentAwareName generates an environment-aware project name
func (pim *ProjectIdentityManager) GenerateEnvironmentAwareName(baseName string, environment string) string {
	// Remove any existing environment suffixes
	cleanName := pim.removeEnvironmentSuffix(baseName)

	// Add environment suffix if not production
	if environment != "production" {
		return fmt.Sprintf("%s-%s", cleanName, environment)
	}

	return cleanName
}

// GetProjectAlias generates a project alias for different environments
func (pim *ProjectIdentityManager) GetProjectAlias(projectName string, targetEnvironment string) string {
	// Remove current environment suffix
	baseName := pim.removeEnvironmentSuffix(projectName)

	// Add target environment suffix if not production
	if targetEnvironment != "production" {
		return fmt.Sprintf("%s-%s", baseName, targetEnvironment)
	}

	return baseName
}

// ValidateProjectPortability validates project portability across environments
func (pim *ProjectIdentityManager) ValidateProjectPortability(project *Project) *ValidationResult {
	result := &ValidationResult{
		Valid: true,
	}

	// Check if project name is portable
	if !pim.isPortableName(project.Name) {
		result.Warnings = append(result.Warnings, ValidationError{
			Field:    "name",
			Message:  "Project name may not be portable across environments",
			Value:    project.Name,
			Severity: "warning",
		})
	}

	// Check if project has environment-specific configuration that might not be portable
	if project.Environment != "" && project.Environment != "development" {
		result.Warnings = append(result.Warnings, ValidationError{
			Field:    "environment",
			Message:  fmt.Sprintf("Project is configured for '%s' environment, may not be portable", project.Environment),
			Value:    project.Environment,
			Severity: "warning",
		})
	}

	return result
}

// validateNameFormat validates project name format
func (pim *ProjectIdentityManager) validateNameFormat(name string) *ValidationError {
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

	// Check for reserved names
	reservedNames := []string{"spooky", "config", "global", "system", "admin", "root"}
	for _, reserved := range reservedNames {
		if strings.EqualFold(name, reserved) {
			return &ValidationError{
				Field:    "name",
				Message:  fmt.Sprintf("Project name '%s' is reserved", reserved),
				Value:    name,
				Severity: "error",
			}
		}
	}

	return nil
}

// validateNameUniqueness validates project name uniqueness
func (pim *ProjectIdentityManager) validateNameUniqueness(name string, projectPath string) *ValidationError {
	// Check if there's already a project.hcl file with the same name in the same directory
	projectHCLPath := filepath.Join(projectPath, "project.hcl")
	if _, err := os.Stat(projectHCLPath); err == nil {
		// TODO: Parse existing project.hcl to check name
		// For now, we'll assume it's unique if the file exists
		return nil
	}

	// Check parent directory for other projects with the same name
	parentDir := filepath.Dir(projectPath)
	if parentDir != "." && parentDir != "/" {
		entries, err := os.ReadDir(parentDir)
		if err == nil {
			for _, entry := range entries {
				if entry.IsDir() && entry.Name() != filepath.Base(projectPath) {
					otherProjectHCL := filepath.Join(parentDir, entry.Name(), "project.hcl")
					if _, err := os.Stat(otherProjectHCL); err == nil {
						// TODO: Parse other project.hcl to check name
						// For now, we'll assume it's unique
					}
				}
			}
		}
	}

	return nil
}

// cleanProjectName cleans a directory name to make it a valid project name
func (pim *ProjectIdentityManager) cleanProjectName(dirName string) string {
	// Remove any file extensions
	cleanName := strings.TrimSuffix(dirName, filepath.Ext(dirName))

	// Replace invalid characters with hyphens
	invalidChars := regexp.MustCompile(`[^a-zA-Z0-9._-]`)
	cleanName = invalidChars.ReplaceAllString(cleanName, "-")

	// Remove multiple consecutive hyphens
	multipleHyphens := regexp.MustCompile(`-+`)
	cleanName = multipleHyphens.ReplaceAllString(cleanName, "-")

	// Remove leading/trailing hyphens
	cleanName = strings.Trim(cleanName, "-")

	// Ensure it starts with a letter
	if cleanName != "" && !regexp.MustCompile(`^[a-zA-Z]`).MatchString(cleanName) {
		cleanName = "project-" + cleanName
	}

	// If empty after cleaning, use default name
	if cleanName == "" {
		cleanName = "spooky-project"
	}

	// Convert to lowercase for consistency
	cleanName = strings.ToLower(cleanName)

	return cleanName
}

// hasEnvironmentPattern checks if a name contains environment-specific patterns
func (pim *ProjectIdentityManager) hasEnvironmentPattern(name string) bool {
	environmentPatterns := []string{"-prod", "-staging", "-dev", "-test", "-development", "-testing"}
	for _, pattern := range environmentPatterns {
		if strings.Contains(name, pattern) {
			return true
		}
	}
	return false
}

// matchesEnvironment checks if a name matches the specified environment
func (pim *ProjectIdentityManager) matchesEnvironment(name string, environment string) bool {
	environmentPatterns := map[string][]string{
		"production":  {"-prod", "-production"},
		"staging":     {"-staging"},
		"development": {"-dev", "-development"},
		"testing":     {"-test", "-testing"},
	}

	if patterns, exists := environmentPatterns[environment]; exists {
		for _, pattern := range patterns {
			if strings.Contains(name, pattern) {
				return true
			}
		}
	}

	return false
}

// removeEnvironmentSuffix removes environment suffix from project name
func (pim *ProjectIdentityManager) removeEnvironmentSuffix(name string) string {
	environmentSuffixes := []string{
		"-prod", "-production",
		"-staging",
		"-dev", "-development",
		"-test", "-testing",
	}

	for _, suffix := range environmentSuffixes {
		if strings.HasSuffix(name, suffix) {
			return strings.TrimSuffix(name, suffix)
		}
	}

	return name
}

// isPortableName checks if a project name is portable across environments
func (pim *ProjectIdentityManager) isPortableName(name string) bool {
	// Names without environment patterns are portable
	return !pim.hasEnvironmentPattern(name)
}

// GetProjectNameVariants generates all possible name variants for a project
func (pim *ProjectIdentityManager) GetProjectNameVariants(baseName string) map[string]string {
	variants := make(map[string]string)

	// Clean base name
	cleanName := pim.removeEnvironmentSuffix(baseName)

	// Generate variants for each environment
	environments := []string{"production", "staging", "development", "testing"}
	for _, env := range environments {
		variants[env] = pim.GenerateEnvironmentAwareName(cleanName, env)
	}

	return variants
}
