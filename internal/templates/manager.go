// Package templates provides template management functionality for the spooky codebase.
package templates

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

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
	tmplData := &spookytypes.Template{
		SourcePath: templatePath,
		Content:    string(data),
		ID:         filepath.Base(templatePath),
	}

	// Validate template
	if err := m.validateTemplate(ctx, tmplData); err != nil {
		return nil, fmt.Errorf("template validation failed: %w", err)
	}

	m.logger.Info("Template loaded successfully", map[string]interface{}{
		"path": templatePath,
		"size": len(data),
	})

	return tmplData, nil
}

// RenderTemplate renders a template with the given data using Go's text/template engine
func (m *Manager) RenderTemplate(_ context.Context, tmplData *spookytypes.Template, data map[string]interface{}) (string, error) {
	m.logger.Info("Rendering template", map[string]interface{}{
		"template":  tmplData.ID,
		"data_keys": len(data),
	})

	// Parse the template using Go's text/template engine
	tmpl, err := template.New(tmplData.ID).Parse(tmplData.Content)
	if err != nil {
		m.logger.Error("Failed to parse template", err, map[string]interface{}{
			"template": tmplData.ID,
			"error":    err.Error(),
		})
		return "", fmt.Errorf("failed to parse template %s: %w", tmplData.ID, err)
	}

	// Execute the template with the provided data
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		m.logger.Error("Failed to execute template", err, map[string]interface{}{
			"template": tmplData.ID,
			"error":    err.Error(),
		})
		return "", fmt.Errorf("failed to execute template %s: %w", tmplData.ID, err)
	}

	result := buf.String()

	m.logger.Info("Template rendered successfully", map[string]interface{}{
		"template":      tmplData.ID,
		"result_length": len(result),
	})

	return result, nil
}

// ValidateTemplate validates a template
func (m *Manager) ValidateTemplate(_ context.Context, tmplData *spookytypes.Template) (*spookytypesschemas.ValidationResult, error) {
	m.logger.Info("Validating template", map[string]interface{}{
		"template": tmplData.ID,
	})

	// Basic validation
	var errors []spookytypesschemas.SchemaError
	var warnings []spookytypesschemas.SchemaError

	// Check if template has content
	if tmplData.Content == "" {
		errors = append(errors, spookytypesschemas.SchemaError{
			Message: "template content cannot be empty",
		})
	}

	// Check for basic template syntax
	if strings.Contains(tmplData.Content, "{{") && !strings.Contains(tmplData.Content, "}}") {
		errors = append(errors, spookytypesschemas.SchemaError{
			Message: "unclosed template variable",
		})
	}

	valid := len(errors) == 0

	m.logger.Info("Template validation completed", map[string]interface{}{
		"template": tmplData.ID,
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
func (m *Manager) validateTemplate(ctx context.Context, tmplData *spookytypes.Template) error {
	result, err := m.ValidateTemplate(ctx, tmplData)
	if err != nil {
		return err
	}

	if !result.Valid {
		return fmt.Errorf("template validation failed: %v", result.Errors)
	}

	return nil
}
