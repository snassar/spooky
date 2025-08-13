// Package templates provides template management functionality for the spooky codebase.
package templates

import (
	"context"
	"fmt"
	"os"

	spookyinterfaces "spooky/internal/interfaces"
	spookytypes "spooky/internal/types"
	spookytypesschemas "spooky/internal/types/schemas"
)

// Integration implements the TemplatesIntegration interface
type Integration struct {
	manager *Manager
}

// NewIntegration creates a new templates integration
func NewIntegration(manager *Manager) spookyinterfaces.TemplatesIntegration {
	return &Integration{
		manager: manager,
	}
}

// LoadTemplate loads a template from the given path
func (i *Integration) LoadTemplate(ctx context.Context, templatePath string) (*spookytypes.Template, error) {
	if templatePath == "" {
		return nil, fmt.Errorf("template path cannot be empty")
	}

	// Check if template file exists
	if _, err := os.Stat(templatePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("template file not found: %s", templatePath)
	}

	// Load template using the manager
	template, err := i.manager.LoadTemplate(ctx, templatePath)
	if err != nil {
		return nil, fmt.Errorf("failed to load template %s: %w", templatePath, err)
	}

	return template, nil
}

// RenderTemplate renders a template with the given data
func (i *Integration) RenderTemplate(ctx context.Context, template *spookytypes.Template, data map[string]interface{}) (string, error) {
	if template == nil {
		return "", fmt.Errorf("template cannot be nil")
	}

	if data == nil {
		data = make(map[string]interface{})
	}

	// Render template using the manager
	result, err := i.manager.RenderTemplate(ctx, template, data)
	if err != nil {
		return "", fmt.Errorf("failed to render template: %w", err)
	}

	return result, nil
}

// ValidateTemplate validates a template
func (i *Integration) ValidateTemplate(ctx context.Context, template *spookytypes.Template) (*spookytypes.ValidationResult, error) {
	if template == nil {
		return &spookytypes.ValidationResult{
			Valid:    false,
			Errors:   []spookytypesschemas.SchemaError{{Message: "template cannot be nil"}},
			Warnings: []spookytypesschemas.SchemaError{},
		}, nil
	}

	// Validate template using the manager
	result, err := i.manager.ValidateTemplate(ctx, template)
	if err != nil {
		return &spookytypes.ValidationResult{
			Valid:    false,
			Errors:   []spookytypesschemas.SchemaError{{Message: fmt.Sprintf("template validation failed: %v", err)}},
			Warnings: []spookytypesschemas.SchemaError{},
		}, nil
	}

	// Convert schema validation result to interface validation result
	return &spookytypes.ValidationResult{
		Valid:    result.Valid,
		Errors:   result.Errors,
		Warnings: result.Warnings,
	}, nil
}
