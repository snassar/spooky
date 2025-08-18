// Package project provides project management functionality for spooky.
package project

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	spookyinterfaces "spooky/internal/interfaces"
	spookyschemas "spooky/internal/schemas"
	spookytypes "spooky/internal/types"
	spookytypescommon "spooky/internal/types/common"
	spookytypeslogging "spooky/internal/types/logging"
	spookytypesschemas "spooky/internal/types/schemas"
)

// Validator implements the ProjectValidator interface
type Validator struct {
	logger                spookytypeslogging.Logger
	schemaDrivenValidator *spookyschemas.SchemaDrivenValidator
	enhancedValidator     *spookyschemas.EnhancedValidator
}

// NewValidator creates a new ProjectValidator instance
func NewValidator(logger spookytypeslogging.Logger) spookyinterfaces.ProjectValidator {
	// Create schema-driven validator for project structure validation
	schemaDrivenConfig := &spookyschemas.SchemaDrivenValidationConfig{
		UseEmbeddedSchemas: true,
		StrictValidation:   true,
		AllowUnknownFields: false,
		DetailedErrors:     true,
	}
	schemaDrivenValidator := spookyschemas.NewSchemaDrivenValidator(logger, schemaDrivenConfig)

	// Create enhanced validator for configuration validation
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

// ValidateProject validates a project structure
func (v *Validator) ValidateProject(ctx context.Context, project *spookytypes.Project) (*spookytypesschemas.ValidationResult, error) {
	v.logger.Info("Validating project structure", map[string]interface{}{
		"project_path": project.Path,
		"project_name": project.Config.Name,
	})

	result := &spookytypesschemas.ValidationResult{
		Valid:       true,
		ValidatedAt: project.CreatedAt,
		Errors:      []spookytypesschemas.SchemaError{},
		Warnings:    []spookytypesschemas.SchemaError{},
	}

	// Validate project directory structure using schema-driven validator
	dirResult, err := v.schemaDrivenValidator.ValidateProjectStructure(ctx, project.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to validate project directory: %w", err)
	}

	// Merge validation results
	if !dirResult.Valid {
		result.Valid = false
	}
	result.Errors = append(result.Errors, dirResult.Errors...)
	result.Warnings = append(result.Warnings, dirResult.Warnings...)

	// Validate project configuration using enhanced validator
	if project.Config != nil {
		// Get project schema for enhanced validation
		projectSchema, err := v.getProjectSchema()
		if err != nil {
			return nil, fmt.Errorf("failed to get project schema: %w", err)
		}

		configResult, err := v.enhancedValidator.ValidateWithEnhancedFeatures(ctx, projectSchema, project.Config)
		if err != nil {
			return nil, fmt.Errorf("failed to validate project config: %w", err)
		}

		// Merge validation results
		if !configResult.Valid {
			result.Valid = false
		}
		result.Errors = append(result.Errors, configResult.Errors...)
		result.Warnings = append(result.Warnings, configResult.Warnings...)
	} else {
		// Project config is missing
		result.Valid = false
		result.Errors = append(result.Errors, spookytypesschemas.SchemaError{
			Code:     "missing_project_config",
			Message:  "Project configuration is missing",
			Severity: "error",
		})
	}

	return result, nil
}

// ValidateProjectDirectory validates project directory structure
func (v *Validator) ValidateProjectDirectory(_ context.Context, projectPath string) (*spookytypesschemas.ValidationResult, error) {
	v.logger.Info("Validating project directory structure", map[string]interface{}{
		"project_path": projectPath,
	})

	result := &spookytypesschemas.ValidationResult{
		Valid:       true,
		ValidatedAt: time.Now(),
		Errors:      []spookytypesschemas.SchemaError{},
		Warnings:    []spookytypesschemas.SchemaError{},
	}

	// Resolve absolute path
	absPath, err := filepath.Abs(projectPath)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve project path: %w", err)
	}

	// Check if project directory exists
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		result.Valid = false
		result.Errors = append(result.Errors, spookytypesschemas.SchemaError{
			Code:     "project_directory_not_found",
			Message:  fmt.Sprintf("Project directory does not exist: %s", absPath),
			Severity: "error",
		})
		return result, nil
	}

	// Check required files according to project-directory.schema.hcl
	requiredFiles := []string{
		"project.hcl", // Required by schema
	}
	for _, file := range requiredFiles {
		filePath := filepath.Join(absPath, file)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			result.Valid = false
			result.Errors = append(result.Errors, spookytypesschemas.SchemaError{
				Code:     "missing_required_file",
				Message:  fmt.Sprintf("Missing required file: %s", file),
				Severity: "error",
			})
		}
	}

	// Check required directories according to project-directory.schema.hcl
	requiredDirs := []string{
		// No required directories
	}
	for _, dir := range requiredDirs {
		dirPath := filepath.Join(absPath, dir)
		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			result.Valid = false
			result.Errors = append(result.Errors, spookytypesschemas.SchemaError{
				Code:     "missing_required_directory",
				Message:  fmt.Sprintf("Missing required directory: %s", dir),
				Severity: "error",
			})
		}
	}

	// Check optional files and directories
	optionalFiles := []string{
		"machines.hcl",
		"actions.hcl",
		"variables.hcl",
		"README.md",
	}
	for _, file := range optionalFiles {
		filePath := filepath.Join(absPath, file)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			result.Warnings = append(result.Warnings, spookytypesschemas.SchemaError{
				Code:     "optional_file_not_found",
				Message:  fmt.Sprintf("Optional file not found: %s", file),
				Severity: "warning",
			})
		}
	}

	optionalDirs := []string{
		"machines",
		"actions",
		"variables",
		"templates",
		"files",
		"logs",
	}
	for _, dir := range optionalDirs {
		dirPath := filepath.Join(absPath, dir)
		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			result.Warnings = append(result.Warnings, spookytypesschemas.SchemaError{
				Code:     "optional_directory_not_found",
				Message:  fmt.Sprintf("Optional directory not found: %s", dir),
				Severity: "warning",
			})
		}
	}

	return result, nil
}

// createValidationResult creates a new validation result with default values
func createValidationResult() *spookytypesschemas.ValidationResult {
	return &spookytypesschemas.ValidationResult{
		Valid:       true,
		ValidatedAt: time.Now(),
		Errors:      []spookytypesschemas.SchemaError{},
		Warnings:    []spookytypesschemas.SchemaError{},
	}
}

// validateProjectName validates the project name field
func validateProjectName(config *spookytypes.ProjectConfig, result *spookytypesschemas.ValidationResult) {
	// Validate required project name
	if config.Name == "" {
		result.Valid = false
		result.Errors = append(result.Errors, spookytypesschemas.SchemaError{
			Code:     "missing_project_name",
			Message:  "Project name is required",
			Severity: "error",
		})
		return
	}

	// Validate project name length
	if len(config.Name) > 128 {
		result.Valid = false
		result.Errors = append(result.Errors, spookytypesschemas.SchemaError{
			Code:     "project_name_too_long",
			Message:  "Project name must be 128 characters or less",
			Severity: "error",
		})
	}
}

// validateProjectDescription validates the project description field
func validateProjectDescription(config *spookytypes.ProjectConfig, result *spookytypesschemas.ValidationResult) {
	if config.Description != "" && len(config.Description) > 1024 {
		result.Warnings = append(result.Warnings, spookytypesschemas.SchemaError{
			Code:     "description_too_long",
			Message:  "Project description should be 1024 characters or less",
			Severity: "warning",
		})
	}
}

// validateMetadataVersion validates the metadata version field
func validateMetadataVersion(metadata *spookytypes.ProjectMetadata, result *spookytypesschemas.ValidationResult) {
	if metadata.Version == "" {
		return
	}

	if !spookytypescommon.IsValidScalVerFormat(metadata.Version) {
		result.Valid = false
		result.Errors = append(result.Errors, spookytypesschemas.SchemaError{
			Code:     "invalid_version_format",
			Message:  "Project version must be in ScalVer format (MAJOR.DATE.PATCH)",
			Severity: "error",
		})
	}
}

// validateMetadataAuthor validates the metadata author field
func validateMetadataAuthor(metadata *spookytypes.ProjectMetadata, result *spookytypesschemas.ValidationResult) {
	if metadata.Author != "" && len(metadata.Author) > 128 {
		result.Warnings = append(result.Warnings, spookytypesschemas.SchemaError{
			Code:     "author_too_long",
			Message:  "Project author should be 128 characters or less",
			Severity: "warning",
		})
	}
}

// validateProjectMetadata validates the project metadata fields
func validateProjectMetadata(config *spookytypes.ProjectConfig, result *spookytypesschemas.ValidationResult) {
	if config.Metadata == nil {
		return
	}

	validateMetadataVersion(config.Metadata, result)
	validateMetadataAuthor(config.Metadata, result)
}

// validateParallelWorkers validates the parallel workers setting
func validateParallelWorkers(settings *spookytypes.ProjectSettings, result *spookytypesschemas.ValidationResult) {
	if settings.ParallelWorkers < 1 || settings.ParallelWorkers > 100 {
		result.Warnings = append(result.Warnings, spookytypesschemas.SchemaError{
			Code:     "invalid_parallel_workers",
			Message:  "Parallel workers should be between 1 and 100",
			Severity: "warning",
		})
	}
}

// validateTimeoutSettings validates the timeout settings
func validateTimeoutSettings(settings *spookytypes.ProjectSettings, result *spookytypesschemas.ValidationResult) {
	if settings.TimeoutSeconds < 1 || settings.TimeoutSeconds > 3600 {
		result.Warnings = append(result.Warnings, spookytypesschemas.SchemaError{
			Code:     "invalid_timeout",
			Message:  "Timeout should be between 1 and 3600 seconds",
			Severity: "warning",
		})
	}
}

// validateRetrySettings validates the retry settings
func validateRetrySettings(settings *spookytypes.ProjectSettings, result *spookytypesschemas.ValidationResult) {
	if settings.MaxRetries < 0 || settings.MaxRetries > 10 {
		result.Warnings = append(result.Warnings, spookytypesschemas.SchemaError{
			Code:     "invalid_max_retries",
			Message:  "Max retries should be between 0 and 10",
			Severity: "warning",
		})
	}
}

// validateProjectSettings validates the project settings fields
func validateProjectSettings(config *spookytypes.ProjectConfig, result *spookytypesschemas.ValidationResult) {
	if config.Settings == nil {
		return
	}

	validateParallelWorkers(config.Settings, result)
	validateTimeoutSettings(config.Settings, result)
	validateRetrySettings(config.Settings, result)
}

// ValidateProjectConfig validates project configuration
func (v *Validator) ValidateProjectConfig(_ context.Context, config *spookytypes.ProjectConfig) (*spookytypesschemas.ValidationResult, error) {
	v.logger.Info("Validating project configuration", map[string]interface{}{
		"project_name": config.Name,
	})

	result := createValidationResult()

	// Validate project name
	validateProjectName(config, result)

	// Validate project description
	validateProjectDescription(config, result)

	// Validate project metadata
	validateProjectMetadata(config, result)

	// Validate project settings
	validateProjectSettings(config, result)

	return result, nil
}

// getProjectSchema gets the project schema for validation
func (v *Validator) getProjectSchema() (*spookytypesschemas.Schema, error) {
	// Try to get schema from embedded schemas first
	if schema, err := v.schemaDrivenValidator.GetEmbeddedSchema("project"); err == nil {
		return schema, nil
	}

	// Fallback: create a basic schema
	return &spookytypesschemas.Schema{
		Name:        "project",
		Type:        "hcl",
		Version:     "1.0",
		Description: "Project configuration schema",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Content:     "", // Will be loaded from file if needed
		Metadata:    make(map[string]interface{}),
	}, nil
}
