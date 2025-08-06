package engine

import (
	"text/template"
	"time"
)

// EngineManager defines the interface for template engine operations
type EngineManager interface {
	// Core engine operations
	ParseTemplate(content []byte, name string) (*template.Template, error)
	RenderTemplate(tmpl *template.Template, data interface{}) (string, error)
	ValidateTemplate(tmpl *template.Template) error

	// Configuration
	SetDelimiters(left, right string) error
	SetMaxExecutionTime(timeout time.Duration) error
	EnableStrictMode(strict bool) error

	// Utility operations
	GetTemplateFunctions() template.FuncMap
	Close() error
}

// TemplateParser defines the interface for template parsing
type TemplateParser interface {
	Parse(content []byte, name string) (*template.Template, error)
	ParseFile(path string) (*template.Template, error)
	ValidateSyntax(content []byte) error
}

// TemplateRenderer defines the interface for template rendering
type TemplateRenderer interface {
	Render(tmpl *template.Template, data interface{}) (string, error)
	RenderWithTimeout(tmpl *template.Template, data interface{}, timeout time.Duration) (string, error)
	ValidateData(data interface{}) error
}
