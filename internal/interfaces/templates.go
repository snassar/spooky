package interfaces

import (
	"context"
	"text/template"
	"time"

	spookytypestemplates "spooky/internal/types/templates"
)

// TemplateManager defines the main interface for template management
type TemplateManager interface {
	// Core template operations
	LoadTemplate(path string) (*spookytypestemplates.Template, error)
	RenderTemplate(templateFile, projectPath string, additionalData map[string]interface{}) (string, error)
	ValidateTemplate(templateFile string) error
	ValidateTemplates(projectPath string) ([]string, error)

	// Context management
	NewTemplateContext(projectPath string) (*spookytypestemplates.TemplateContext, error)
	GetTemplateFunctions(ctx *spookytypestemplates.TemplateContext) template.FuncMap

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

// TemplateValidatorCore defines the interface for template validation
type TemplateValidatorCore interface {
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

// FunctionsManager defines the interface for template functions
type FunctionsManager interface {
	// Core functions operations
	GetBuiltinFunctions() template.FuncMap
	RegisterCustomFunction(name string, fn interface{}) error
	ValidateFunction(name string, fn interface{}) error

	// Function management
	GetFunction(name string) (interface{}, bool)
	ListFunctions() []string
	RemoveFunction(name string) error

	// Configuration
	EnableBuiltinFunctions(enabled bool) error
	SetFunctionTimeout(timeout time.Duration) error

	// Utility operations
	Close() error
}

// FunctionValidator defines the interface for function validation
type FunctionValidator interface {
	ValidateFunction(name string, fn interface{}) error
	ValidateFunctionSignature(fn interface{}) error
	ValidateFunctionReturnType(fn interface{}) error
}

// TemplateSecretsManager defines the interface for template secrets integration
type TemplateSecretsManager interface {
	// Core secrets operations
	EncryptValue(value string) (string, error)
	DecryptValue(encryptedValue string) (string, error)
	ValidateEncryptionKey(key string) error

	// Configuration
	SetEncryptionKey(key string) error
	SetEncryptionAlgorithm(algorithm string) error
	EnableEncryption(enabled bool) error

	// Utility operations
	IsEncrypted(value string) bool
	GetEncryptionAlgorithm() string
	Close() error
}

// EncryptionProvider defines the interface for encryption providers
type EncryptionProvider interface {
	Encrypt(data []byte) ([]byte, error)
	Decrypt(data []byte) ([]byte, error)
	GetAlgorithm() string
}

// TemplateValidationManager defines the interface for template validation
type TemplateValidationManager interface {
	// Core validation operations
	ValidateTemplate(path string) error
	ValidateTemplates(projectPath string) ([]string, error)
	ValidateSyntax(content []byte) error
	ValidateFunctions(tmpl *template.Template) error

	// Schema validation
	ValidateAgainstSchema(template *spookytypestemplates.Template, schemaName string) error

	// Configuration
	SetValidationRules(rules *spookytypestemplates.ValidationRules) error
	EnableStrictValidation(strict bool) error

	// Utility operations
	GetValidationErrors() []spookytypestemplates.ValidationError
	ClearValidationErrors() error
	Close() error
}

// TemplateValidator defines the interface for specific validators
type TemplateValidator interface {
	Validate(template *spookytypestemplates.Template) error
	GetName() string
	GetDescription() string
}
