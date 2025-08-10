package validation

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	spookyinterfaces "spooky/internal/interfaces"
	spookytypes "spooky/internal/types"
)

// Manager implements ValidationManager interface
type Manager struct {
	config     *spookytypes.TemplateValidationConfig
	validators map[string]spookyinterfaces.TemplateValidator
	logger     spookyinterfaces.Logger
	errors     []spookytypes.TemplateValidationError
}

// NewManager creates a new validation manager
func NewManager(config *spookytypes.TemplateValidationConfig, logger spookyinterfaces.Logger) *Manager {
	return &Manager{
		config:     config,
		validators: make(map[string]spookyinterfaces.TemplateValidator),
		logger:     logger,
		errors:     make([]spookytypes.TemplateValidationError, 0),
	}
}

// ValidateTemplate validates a single template
func (m *Manager) ValidateTemplate(path string) error {
	// 1. Check file exists and is readable
	if err := m.validateFileExists(path); err != nil {
		return fmt.Errorf("file validation failed: %w", err)
	}

	// 2. Load template content
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read template: %w", err)
	}

	// 3. Validate syntax
	if err := m.ValidateSyntax(content); err != nil {
		return fmt.Errorf("syntax validation failed: %w", err)
	}

	// 4. Parse template to check functions
	tmpl, err := template.New("validation").Parse(string(content))
	if err != nil {
		return fmt.Errorf("template parsing failed: %w", err)
	}

	// 5. Validate functions
	if err := m.ValidateFunctions(tmpl); err != nil {
		return fmt.Errorf("function validation failed: %w", err)
	}

	return nil
}

// ValidateTemplates validates all templates in a project
func (m *Manager) ValidateTemplates(projectPath string) ([]string, error) {
	templatesDir := filepath.Join(projectPath, "templates")

	// Find all template files (.tmpl extension)
	templateFiles, err := filepath.Glob(filepath.Join(templatesDir, "*.tmpl"))
	if err != nil {
		return nil, fmt.Errorf("failed to find template files: %w", err)
	}

	var errors []string

	// Validate each template
	for _, templateFile := range templateFiles {
		if err := m.ValidateTemplate(templateFile); err != nil {
			errors = append(errors, fmt.Sprintf("%s: %v", templateFile, err))
		}
	}

	return errors, nil
}

// ValidateSyntax validates template syntax
func (m *Manager) ValidateSyntax(content []byte) error {
	// Parse template to check syntax
	tmpl, err := template.New("validation").Parse(string(content))
	if err != nil {
		return fmt.Errorf("template syntax error: %w", err)
	}

	// Validate functions
	if err := m.ValidateFunctions(tmpl); err != nil {
		return fmt.Errorf("function validation failed: %w", err)
	}

	return nil
}

// ValidateFunctions validates template functions
func (m *Manager) ValidateFunctions(tmpl *template.Template) error {
	// Basic function validation - check for restricted patterns
	content := tmpl.Tree.Root.String()

	// Check for potentially dangerous patterns
	restrictedPatterns := []string{
		"os.Exec",
		"os.Command",
		"exec.Command",
		"system(",
		"eval(",
	}

	for _, pattern := range restrictedPatterns {
		if strings.Contains(content, pattern) {
			return fmt.Errorf("template contains restricted pattern: %s", pattern)
		}
	}

	return nil
}

// ValidateAgainstSchema validates a template against a schema
func (m *Manager) ValidateAgainstSchema(template *spookytypes.Template, _ string) error {
	// Basic schema validation
	if template == nil {
		return fmt.Errorf("template is nil")
	}

	if template.Name == "" {
		return fmt.Errorf("template name is required")
	}

	if template.Source == "" {
		return fmt.Errorf("template source is required")
	}

	return nil
}

// SetValidationRules sets validation rules
func (m *Manager) SetValidationRules(rules *spookytypes.TemplateValidationRules) error {
	if m.config == nil {
		m.config = &spookytypes.TemplateValidationConfig{}
	}
	m.config.ValidationRules = rules
	return nil
}

// EnableStrictValidation enables strict validation
func (m *Manager) EnableStrictValidation(strict bool) error {
	if m.config == nil {
		m.config = &spookytypes.TemplateValidationConfig{}
	}
	m.config.StrictValidation = strict
	return nil
}

// GetValidationErrors returns validation errors
func (m *Manager) GetValidationErrors() []spookytypes.TemplateValidationError {
	return m.errors
}

// ClearValidationErrors clears validation errors
func (m *Manager) ClearValidationErrors() error {
	m.errors = make([]spookytypes.TemplateValidationError, 0)
	return nil
}

// Close closes the validation manager
func (m *Manager) Close() error {
	// Cleanup resources if needed
	return nil
}

// validateFileExists checks if a file exists and is readable
func (m *Manager) validateFileExists(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("file does not exist or is not accessible: %w", err)
	}

	if info.IsDir() {
		return fmt.Errorf("path is a directory, not a file")
	}

	return nil
}
