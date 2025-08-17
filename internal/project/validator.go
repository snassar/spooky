// Package project provides project management functionality for spooky.
package project

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	spookyinterfaces "spooky/internal/interfaces"
	spookytypes "spooky/internal/types"
	spookytypescommon "spooky/internal/types/common"
	spookytypeslogging "spooky/internal/types/logging"
	spookytypesschemas "spooky/internal/types/schemas"
)

// Validator implements the ProjectValidator interface
type Validator struct {
	logger spookytypeslogging.Logger
}

// NewValidator creates a new ProjectValidator instance
func NewValidator(logger spookytypeslogging.Logger) spookyinterfaces.ProjectValidator {
	return &Validator{
		logger: logger,
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

	// Validate project directory structure
	dirResult, err := v.ValidateProjectDirectory(ctx, project.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to validate project directory: %w", err)
	}

	// Merge validation results
	if !dirResult.Valid {
		result.Valid = false
	}
	result.Errors = append(result.Errors, dirResult.Errors...)
	result.Warnings = append(result.Warnings, dirResult.Warnings...)

	// Validate project configuration
	if project.Config != nil {
		configResult, err := v.ValidateProjectConfig(ctx, project.Config)
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

// ValidateProjectConfig validates project configuration
func (v *Validator) ValidateProjectConfig(_ context.Context, config *spookytypes.ProjectConfig) (*spookytypesschemas.ValidationResult, error) {
	v.logger.Info("Validating project configuration", map[string]interface{}{
		"project_name": config.Name,
	})

	result := &spookytypesschemas.ValidationResult{
		Valid:       true,
		ValidatedAt: time.Now(),
		Errors:      []spookytypesschemas.SchemaError{},
		Warnings:    []spookytypesschemas.SchemaError{},
	}

	// Validate project name
	if config.Name == "" {
		result.Valid = false
		result.Errors = append(result.Errors, spookytypesschemas.SchemaError{
			Code:     "missing_project_name",
			Message:  "Project name is required",
			Severity: "error",
		})
	}

	// Validate project name pattern (alphanumeric, dots, underscores, hyphens)
	if config.Name != "" {
		// Basic validation - could be enhanced with regex
		if len(config.Name) > 128 {
			result.Valid = false
			result.Errors = append(result.Errors, spookytypesschemas.SchemaError{
				Code:     "project_name_too_long",
				Message:  "Project name must be 128 characters or less",
				Severity: "error",
			})
		}
	}

	// Validate description length
	if config.Description != "" && len(config.Description) > 1024 {
		result.Warnings = append(result.Warnings, spookytypesschemas.SchemaError{
			Code:     "description_too_long",
			Message:  "Project description should be 1024 characters or less",
			Severity: "warning",
		})
	}

	// Validate metadata if present
	if config.Metadata != nil {
		if config.Metadata.Version != "" {
			// Use centralized ScalVer validation
			if !spookytypescommon.IsValidScalVerFormat(config.Metadata.Version) {
				result.Valid = false
				result.Errors = append(result.Errors, spookytypesschemas.SchemaError{
					Code:     "invalid_version_format",
					Message:  "Project version must be in ScalVer format (MAJOR.DATE.PATCH)",
					Severity: "error",
				})
			}
		}

		if config.Metadata.Author != "" && len(config.Metadata.Author) > 128 {
			result.Warnings = append(result.Warnings, spookytypesschemas.SchemaError{
				Code:     "author_too_long",
				Message:  "Project author should be 128 characters or less",
				Severity: "warning",
			})
		}
	}

	// Validate settings if present
	if config.Settings != nil {
		if config.Settings.ParallelWorkers < 1 || config.Settings.ParallelWorkers > 100 {
			result.Warnings = append(result.Warnings, spookytypesschemas.SchemaError{
				Code:     "invalid_parallel_workers",
				Message:  "Parallel workers should be between 1 and 100",
				Severity: "warning",
			})
		}

		if config.Settings.TimeoutSeconds < 1 || config.Settings.TimeoutSeconds > 3600 {
			result.Warnings = append(result.Warnings, spookytypesschemas.SchemaError{
				Code:     "invalid_timeout",
				Message:  "Timeout should be between 1 and 3600 seconds",
				Severity: "warning",
			})
		}

		if config.Settings.MaxRetries < 0 || config.Settings.MaxRetries > 10 {
			result.Warnings = append(result.Warnings, spookytypesschemas.SchemaError{
				Code:     "invalid_max_retries",
				Message:  "Max retries should be between 0 and 10",
				Severity: "warning",
			})
		}
	}

	return result, nil
}
