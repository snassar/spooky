package interfaces

import "spooky/internal/templates/types"

// TemplatesIntegration defines the interface for templates system integration
type TemplatesIntegration interface {
	// LoadTemplates loads templates from the project
	LoadTemplates(projectPath string) (*TemplatesContext, error)

	// ValidateTemplate validates a template using the context
	ValidateTemplate(template *types.Template, context *TemplatesContext) error

	// RenderTemplate renders a template with the given context
	RenderTemplate(template *types.Template, context *TemplatesContext) (string, error)

	// CacheTemplate caches a template for later use
	CacheTemplate(template *types.Template) error

	// GetTemplate gets a specific template by name
	GetTemplate(name string, context *TemplatesContext) (*types.Template, error)

	// ListTemplates lists all available templates
	ListTemplates(context *TemplatesContext) (map[string]*types.Template, error)

	// AddTemplate adds a new template to the project
	AddTemplate(name string, template *types.Template, context *TemplatesContext) error

	// RemoveTemplate removes a template from the project
	RemoveTemplate(name string, context *TemplatesContext) error
}
