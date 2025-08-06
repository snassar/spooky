package validation

import (
	"text/template"

	"spooky/internal/templates/types"
)

// ValidationManager defines the interface for template validation
type ValidationManager interface {
	// Core validation operations
	ValidateTemplate(path string) error
	ValidateTemplates(projectPath string) ([]string, error)
	ValidateSyntax(content []byte) error
	ValidateFunctions(tmpl *template.Template) error

	// Schema validation
	ValidateAgainstSchema(template *types.Template, schemaName string) error

	// Configuration
	SetValidationRules(rules *types.ValidationRules) error
	EnableStrictValidation(strict bool) error

	// Utility operations
	GetValidationErrors() []types.ValidationError
	ClearValidationErrors() error
	Close() error
}

// TemplateValidator defines the interface for specific validators
type TemplateValidator interface {
	Validate(template *types.Template) error
	GetName() string
	GetDescription() string
}
