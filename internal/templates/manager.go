// Package templates provides template management functionality for the spooky codebase.
package templates

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	spookytypes "spooky/internal/types"
	spookytypeslogging "spooky/internal/types/logging"
	spookytypesschemas "spooky/internal/types/schemas"
)

// Manager provides template management functionality
type Manager struct {
	logger spookytypeslogging.Logger
}

// NewManager creates a new template manager
func NewManager(
	logger spookytypeslogging.Logger,
) *Manager {
	return &Manager{
		logger: logger,
	}
}

// LoadTemplate loads a template from the given path
func (m *Manager) LoadTemplate(ctx context.Context, templatePath string) (*spookytypes.Template, error) {
	m.logger.Info("Loading template", map[string]interface{}{
		"path": templatePath,
	})

	// Read template file
	data, err := os.ReadFile(templatePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read template file: %w", err)
	}

	// Create template structure
	template := &spookytypes.Template{
		SourcePath: templatePath,
		Content:    string(data),
		ID:         filepath.Base(templatePath),
	}

	// Validate template
	if err := m.validateTemplate(ctx, template); err != nil {
		return nil, fmt.Errorf("template validation failed: %w", err)
	}

	m.logger.Info("Template loaded successfully", map[string]interface{}{
		"path": templatePath,
		"size": len(data),
	})

	return template, nil
}

// RenderTemplate renders a template with the given data
func (m *Manager) RenderTemplate(ctx context.Context, template *spookytypes.Template, data map[string]interface{}) (string, error) {
	m.logger.Info("Rendering template", map[string]interface{}{
		"template":  template.ID,
		"data_keys": len(data),
	})

	// Basic template rendering implementation
	// In a real implementation, this would use a proper template engine
	result := template.Content

	// Simple variable substitution for demonstration
	// Replace {{.variable}} patterns with actual values
	for key, value := range data {
		placeholder := fmt.Sprintf("{{.%s}}", key)
		result = strings.ReplaceAll(result, placeholder, fmt.Sprintf("%v", value))
	}

	m.logger.Info("Template rendered successfully", map[string]interface{}{
		"template":      template.ID,
		"result_length": len(result),
	})

	return result, nil
}

// ValidateTemplate validates a template
func (m *Manager) ValidateTemplate(ctx context.Context, template *spookytypes.Template) (*spookytypesschemas.ValidationResult, error) {
	m.logger.Info("Validating template", map[string]interface{}{
		"template": template.ID,
	})

	// Basic validation
	var errors []spookytypesschemas.SchemaError
	var warnings []spookytypesschemas.SchemaError

	// Check if template has content
	if template.Content == "" {
		errors = append(errors, spookytypesschemas.SchemaError{
			Message: "template content cannot be empty",
		})
	}

	// Check for basic template syntax
	if strings.Contains(template.Content, "{{") && !strings.Contains(template.Content, "}}") {
		errors = append(errors, spookytypesschemas.SchemaError{
			Message: "unclosed template variable",
		})
	}

	valid := len(errors) == 0

	m.logger.Info("Template validation completed", map[string]interface{}{
		"template": template.ID,
		"valid":    valid,
		"errors":   len(errors),
		"warnings": len(warnings),
	})

	return &spookytypesschemas.ValidationResult{
		Valid:    valid,
		Errors:   errors,
		Warnings: warnings,
	}, nil
}

// validateTemplate is a helper method for internal validation
func (m *Manager) validateTemplate(ctx context.Context, template *spookytypes.Template) error {
	result, err := m.ValidateTemplate(ctx, template)
	if err != nil {
		return err
	}

	if !result.Valid {
		return fmt.Errorf("template validation failed: %v", result.Errors)
	}

	return nil
}
