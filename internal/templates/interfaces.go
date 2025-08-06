package templates

import (
	"context"
	"text/template"
	"time"

	spookytemplatestypes "spooky/internal/templates/types"
)

// TemplateManager defines the main interface for template management
type TemplateManager interface {
	// Core template operations
	LoadTemplate(path string) (*spookytemplatestypes.Template, error)
	RenderTemplate(templateFile, projectPath string, additionalData map[string]interface{}) (string, error)
	ValidateTemplate(templateFile string) error
	ValidateTemplates(projectPath string) ([]string, error)

	// Context management
	NewTemplateContext(projectPath string) (*spookytemplatestypes.TemplateContext, error)
	GetTemplateFunctions(ctx *spookytemplatestypes.TemplateContext) template.FuncMap

	// Enhanced operations
	RenderTemplateWithTimeout(ctx context.Context, templateFile, projectPath string, additionalData map[string]interface{}) (string, error)
	SetRenderTimeout(timeout time.Duration)
	SetMaxTemplateSize(maxSize int64)

	// Configuration
	SetDefaultTimeout(timeout time.Duration)
	RegisterCustomFunction(name string, fn interface{}) error

	// Utility operations
	Close() error
}

// TemplateEngine defines the interface for template engine operations
type TemplateEngine interface {
	ParseTemplate(content []byte, name string) (*template.Template, error)
	RenderTemplate(tmpl *template.Template, data interface{}) (string, error)
	ValidateTemplate(tmpl *template.Template) error
	GetTemplateFunctions() template.FuncMap
}

// TemplateValidator defines the interface for template validation
type TemplateValidator interface {
	ValidateTemplate(path string) error
	ValidateTemplates(projectPath string) ([]string, error)
	ValidateSyntax(content []byte) error
	ValidateFunctions(tmpl *template.Template) error
}

// TemplateFunctions defines the interface for template functions
type TemplateFunctions interface {
	GetBuiltinFunctions() template.FuncMap
	RegisterCustomFunction(name string, fn interface{}) error
	ValidateFunction(name string, fn interface{}) error
}

// TemplateSecrets defines the interface for template secrets integration
type TemplateSecrets interface {
	EncryptValue(value string) (string, error)
	DecryptValue(encryptedValue string) (string, error)
	ValidateEncryptionKey(key string) error
}
