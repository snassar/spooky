package engine

import (
	"fmt"
	"strings"
	"text/template"
	"time"

	spookylogging "spooky/internal/logging"
	spookytemplatestypes "spooky/internal/templates/types"
)

// Manager implements EngineManager interface
type Manager struct {
	config   *spookytemplatestypes.EngineConfig
	parser   TemplateParser
	renderer TemplateRenderer
	logger   spookylogging.Logger
}

// NewManager creates a new engine manager
func NewManager(config *spookytemplatestypes.EngineConfig, logger spookylogging.Logger) *Manager {
	return &Manager{
		config:   config,
		parser:   NewTemplateParser(),
		renderer: NewTemplateRenderer(),
		logger:   logger,
	}
}

// ParseTemplate parses template content
func (m *Manager) ParseTemplate(content []byte, name string) (*template.Template, error) {
	// 1. Validate syntax
	if err := m.parser.ValidateSyntax(content); err != nil {
		return nil, fmt.Errorf("syntax validation failed: %w", err)
	}

	// 2. Parse template
	tmpl, err := m.parser.Parse(content, name)
	if err != nil {
		return nil, fmt.Errorf("template parsing failed: %w", err)
	}

	// 3. Validate template
	if err := m.ValidateTemplate(tmpl); err != nil {
		return nil, fmt.Errorf("template validation failed: %w", err)
	}

	return tmpl, nil
}

// RenderTemplate renders a template with data
func (m *Manager) RenderTemplate(tmpl *template.Template, data interface{}) (string, error) {
	// 1. Validate data
	if err := m.renderer.ValidateData(data); err != nil {
		return "", fmt.Errorf("data validation failed: %w", err)
	}

	// 2. Render template
	return m.renderer.Render(tmpl, data)
}

// RenderTemplateWithTimeout renders with timeout
func (m *Manager) RenderTemplateWithTimeout(tmpl *template.Template, data interface{}, timeout time.Duration) (string, error) {
	return m.renderer.RenderWithTimeout(tmpl, data, timeout)
}

// ValidateTemplate validates a template
func (m *Manager) ValidateTemplate(tmpl *template.Template) error {
	// Basic validation - template should not be nil
	if tmpl == nil {
		return fmt.Errorf("template is nil")
	}

	// Additional validation can be added here
	return nil
}

// SetDelimiters sets template delimiters
func (m *Manager) SetDelimiters(left, right string) error {
	if m.config == nil {
		m.config = &spookytemplatestypes.EngineConfig{}
	}
	m.config.Delimiters = []string{left, right}
	return nil
}

// SetMaxExecutionTime sets maximum execution time
func (m *Manager) SetMaxExecutionTime(timeout time.Duration) error {
	if m.config == nil {
		m.config = &spookytemplatestypes.EngineConfig{}
	}
	m.config.MaxExecutionTime = timeout
	return nil
}

// EnableStrictMode enables strict mode
func (m *Manager) EnableStrictMode(strict bool) error {
	if m.config == nil {
		m.config = &spookytemplatestypes.EngineConfig{}
	}
	m.config.StrictMode = strict
	return nil
}

// GetTemplateFunctions returns template functions
func (m *Manager) GetTemplateFunctions() template.FuncMap {
	// Return basic template functions
	return template.FuncMap{
		"upper": strings.ToUpper,
		"lower": strings.ToLower,
		"trim":  strings.TrimSpace,
	}
}

// Close closes the engine manager
func (m *Manager) Close() error {
	// Cleanup resources if needed
	return nil
}
