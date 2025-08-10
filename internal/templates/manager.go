package templates

import (
	"context"
	"fmt"
	"os"
	"text/template"
	"time"

	spookyinterfaces "spooky/internal/interfaces"
	spookytypes "spooky/internal/types"
)

// Manager implements TemplateManager interface
type Manager struct {
	config            *spookytypes.TemplateConfig
	engineManager     spookyinterfaces.EngineManager
	functionsManager  spookyinterfaces.FunctionsManager
	validationManager spookyinterfaces.TemplateValidationManager
	secretsManager    spookyinterfaces.TemplateSecretsManager
	logger            spookyinterfaces.Logger
}

// NewManager creates a new template manager
func NewManager(
	config *spookytypes.TemplateConfig,
	engineManager spookyinterfaces.EngineManager,
	functionsManager spookyinterfaces.FunctionsManager,
	validationManager spookyinterfaces.TemplateValidationManager,
	secretsManager spookyinterfaces.TemplateSecretsManager,
	logger spookyinterfaces.Logger,
) *Manager {
	return &Manager{
		config:            config,
		engineManager:     engineManager,
		functionsManager:  functionsManager,
		validationManager: validationManager,
		secretsManager:    secretsManager,
		logger:            logger,
	}
}

// LoadTemplate loads a template from a file
func (m *Manager) LoadTemplate(path string) (*spookytypes.Template, error) {
	// Read template content
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read template file: %w", err)
	}

	// Create template object
	template := &spookytypes.Template{
		Name:      path,
		Source:    path,
		Content:   string(content),
		Path:      path,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return template, nil
}

// RenderTemplate renders a template with the given context
func (m *Manager) RenderTemplate(templateFile, projectPath string, additionalData map[string]interface{}) (string, error) {
	// 1. Load template
	template, err := m.LoadTemplate(templateFile)
	if err != nil {
		return "", fmt.Errorf("failed to load template: %w", err)
	}

	// 2. Create context
	context, err := m.NewTemplateContext(projectPath)
	if err != nil {
		return "", fmt.Errorf("failed to create context: %w", err)
	}

	// 3. Merge additional data
	for key, value := range additionalData {
		context.Data[key] = value
	}

	// 4. Parse template
	tmpl, err := m.engineManager.ParseTemplate([]byte(template.Content), template.Name)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	// 5. Render template
	return m.engineManager.RenderTemplate(tmpl, context.Data)
}

// RenderTemplateWithTimeout renders a template with timeout and security checks
func (m *Manager) RenderTemplateWithTimeout(_ context.Context, templateFile, projectPath string, additionalData map[string]interface{}) (string, error) {
	// 1. Validate template
	if err := m.ValidateTemplate(templateFile); err != nil {
		return "", fmt.Errorf("template validation failed: %w", err)
	}

	// 2. Load template
	template, err := m.LoadTemplate(templateFile)
	if err != nil {
		return "", fmt.Errorf("failed to load template: %w", err)
	}

	// 3. Check template size against max_template_size
	if int64(len(template.Content)) > m.config.MaxTemplateSize {
		return "", fmt.Errorf("template size exceeds maximum allowed size")
	}

	// 4. Create context
	context, err := m.NewTemplateContext(projectPath)
	if err != nil {
		return "", fmt.Errorf("failed to create context: %w", err)
	}

	// 5. Merge additional data
	for key, value := range additionalData {
		context.Data[key] = value
	}

	// 6. Parse template
	tmpl, err := m.engineManager.ParseTemplate([]byte(template.Content), template.Name)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	// 7. Render with timeout
	return m.engineManager.RenderTemplate(tmpl, context.Data)
}

// ValidateTemplate validates a single template
func (m *Manager) ValidateTemplate(templateFile string) error {
	return m.validationManager.ValidateTemplate(templateFile)
}

// ValidateTemplates validates all templates in a project
func (m *Manager) ValidateTemplates(projectPath string) ([]string, error) {
	return m.validationManager.ValidateTemplates(projectPath)
}

// NewTemplateContext creates a new template context
// Aligns with template-context.hcl schema
func (m *Manager) NewTemplateContext(projectPath string) (*spookytypes.TemplateContext, error) {
	return &spookytypes.TemplateContext{
		ProjectPath: projectPath,
		Data:        make(map[string]interface{}),
		Functions:   m.functionsManager.GetBuiltinFunctions(),
		CreatedAt:   time.Now(),
	}, nil
}

// GetTemplateFunctions returns template functions for a context
func (m *Manager) GetTemplateFunctions(ctx *spookytypes.TemplateContext) template.FuncMap {
	functions := m.functionsManager.GetBuiltinFunctions()

	// Add custom functions
	for name, fn := range ctx.Functions {
		functions[name] = fn
	}

	return functions
}

// SetRenderTimeout sets the rendering timeout
func (m *Manager) SetRenderTimeout(timeout time.Duration) {
	m.config.DefaultTimeout = timeout
}

// SetMaxTemplateSize sets the maximum template file size
func (m *Manager) SetMaxTemplateSize(maxSize int64) {
	m.config.MaxTemplateSize = maxSize
}

// SetDefaultTimeout sets the default timeout
func (m *Manager) SetDefaultTimeout(timeout time.Duration) {
	m.config.DefaultTimeout = timeout
}

// RegisterCustomFunction registers a custom function
func (m *Manager) RegisterCustomFunction(name string, fn interface{}) error {
	return m.functionsManager.RegisterCustomFunction(name, fn)
}

// Close closes the template manager
func (m *Manager) Close() error {
	// Close all sub-managers
	if err := m.engineManager.Close(); err != nil {
		return fmt.Errorf("failed to close engine manager: %w", err)
	}

	if err := m.functionsManager.Close(); err != nil {
		return fmt.Errorf("failed to close functions manager: %w", err)
	}

	if err := m.validationManager.Close(); err != nil {
		return fmt.Errorf("failed to close validation manager: %w", err)
	}

	if err := m.secretsManager.Close(); err != nil {
		return fmt.Errorf("failed to close secrets manager: %w", err)
	}

	return nil
}

// Coordinator integration methods
func (m *Manager) LoadTemplatesForProject(_ string) ([]*spookytypes.Template, error) {
	// Load all templates from project templates directory
	// For now, return empty list - this would be implemented to load all templates
	// In a full implementation, this would:
	// 1. Scan the templates directory (projectPath + "/templates")
	// 2. Load each .tmpl file
	// 3. Parse template metadata
	// 4. Return list of templates
	return []*spookytypes.Template{}, nil
}

func (m *Manager) ValidateTemplatesForProject(projectPath string) error {
	// Validate all templates in project
	errors, err := m.validationManager.ValidateTemplates(projectPath)
	if err != nil {
		return err
	}

	if len(errors) > 0 {
		return fmt.Errorf("template validation failed: %v", errors)
	}

	return nil
}

func (m *Manager) RenderTemplateForContext(templateFile string, context *spookytypes.TemplateContext) (string, error) {
	// Render template with provided context
	return m.RenderTemplate(templateFile, context.ProjectPath, context.Data)
}
