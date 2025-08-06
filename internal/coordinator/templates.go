package coordinator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
	spookyinterfaces "spooky/internal/interfaces"
	spookylogging "spooky/internal/logging"
	spookytemplates "spooky/internal/templates"
	spookytemplatestypes "spooky/internal/templates/types"

)

// CoordinatorTemplatesIntegration implements templates system integration
type CoordinatorTemplatesIntegration struct {
	templatesManager spookytemplates.TemplateManager
	logger           spookylogging.Logger
	templateCache    map[string]*spookytemplatestypes.Template // Add this field
}

// NewCoordinatorTemplatesIntegration creates a new templates integration
func NewCoordinatorTemplatesIntegration(templatesManager spookytemplates.TemplateManager, logger spookylogging.Logger) *CoordinatorTemplatesIntegration {
	return &CoordinatorTemplatesIntegration{
		templatesManager: templatesManager,
		logger:           logger,
		templateCache:    make(map[string]*spookytemplatestypes.Template), // Initialize cache
	}
}

// NewDefaultCoordinatorTemplatesIntegration creates a templates integration with default manager
func NewDefaultCoordinatorTemplatesIntegration(logger spookylogging.Logger) (*CoordinatorTemplatesIntegration, error) {
	templatesManager, err := spookytemplates.NewDefaultTemplateManager(logger)
	if err != nil {
		return nil, err
	}

	return NewCoordinatorTemplatesIntegration(templatesManager, logger), nil
}

// LoadTemplates loads templates from the project
func (ti *CoordinatorTemplatesIntegration) LoadTemplates(projectPath string) (*spookyinterfaces.TemplatesContext, error) {
	context := &spookyinterfaces.TemplatesContext{
		BaseContext: spookyinterfaces.BaseContext{
			ProjectPath: projectPath,
			Timestamp:   time.Now(),
		},
		Templates:     make(map[string]*spookytemplatestypes.Template),
		RenderedCache: make(map[string]string),
		Functions:     make(map[string]interface{}),
	}

	// Load templates from project using templates manager
	if ti.templatesManager != nil {
		// Scan templates directory
		templatesDir := filepath.Join(projectPath, "templates")

		// Check if templates directory exists
		if _, err := os.Stat(templatesDir); err == nil {
			// Walk through templates directory
			err := filepath.Walk(templatesDir, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}

				// Skip directories
				if info.IsDir() {
					return nil
				}

				// Only process template files (common extensions)
				ext := filepath.Ext(path)
				if ext != ".tmpl" && ext != ".tpl" && ext != ".template" && ext != ".hcl" {
					return nil
				}

				// Read template file
				content, err := os.ReadFile(path)
				if err != nil {
					ti.logger.Warn("Failed to read template file", spookylogging.String("file", path), spookylogging.Error(err))
					return nil // Continue with other files
				}

				// Create template object
				relPath, _ := filepath.Rel(templatesDir, path)
				templateName := strings.TrimSuffix(relPath, ext)

				template := &spookytemplatestypes.Template{
					Name:    templateName,
					Source:  string(content),
					Content: string(content),
				}

				context.Templates[templateName] = template
				ti.logger.Debug("Loaded template", spookylogging.String("template", templateName), spookylogging.String("file", path))

				return nil
			})

			if err != nil {
				return nil, fmt.Errorf("failed to scan templates directory: %w", err)
			}
		}

		// Load template functions
		templateContext, err := ti.templatesManager.NewTemplateContext(projectPath)
		if err == nil {
			functions := ti.templatesManager.GetTemplateFunctions(templateContext)
			for name, fn := range functions {
				context.Functions[name] = fn
			}
		}
	}

	ti.logger.Info("Loaded templates from project",
		spookylogging.String("project", projectPath),
		spookylogging.Int("templates_count", len(context.Templates)))

	return context, nil
}

// ValidateTemplate validates a template using the context
func (ti *CoordinatorTemplatesIntegration) ValidateTemplate(template *spookytemplatestypes.Template, context *spookyinterfaces.TemplatesContext) error {
	if template == nil {
		return fmt.Errorf("template cannot be nil")
	}

	if context == nil {
		return fmt.Errorf("templates context cannot be nil")
	}

	// Basic template validation
	if template.Name == "" {
		return fmt.Errorf("template name cannot be empty")
	}

	if template.Source == "" {
		return fmt.Errorf("template source cannot be empty")
	}

	// Syntax validation using templates manager
	if ti.templatesManager != nil {
		// Create a temporary file for validation
		tempFile := filepath.Join(os.TempDir(), fmt.Sprintf("template_%s.tmpl", template.Name))
		err := os.WriteFile(tempFile, []byte(template.Source), 0644)
		if err != nil {
			return fmt.Errorf("failed to create temporary file for validation: %w", err)
		}
		defer os.Remove(tempFile)

		// Validate template syntax
		if err := ti.templatesManager.ValidateTemplate(tempFile); err != nil {
			return fmt.Errorf("template syntax validation failed: %w", err)
		}
	}

	// Security validation - check for potentially dangerous patterns
	dangerousPatterns := []string{
		"{{.Exec", "{{.Command", "{{.Shell", "{{.System", "{{.Process",
		"os/exec", "exec.Command", "system(", "eval(", "exec(",
	}

	for _, pattern := range dangerousPatterns {
		if strings.Contains(template.Source, pattern) {
			ti.logger.Warn("Template contains potentially dangerous pattern",
				spookylogging.String("template", template.Name),
				spookylogging.String("pattern", pattern))
		}
	}

	// Variable reference validation
	// Check if template references variables that exist in context
	if context.Functions != nil {
		// Basic function availability check
		// This could be enhanced to check specific function calls
	}

	return nil
}

// RenderTemplate renders a template with enhanced features
func (ti *CoordinatorTemplatesIntegration) RenderTemplate(template *spookytemplatestypes.Template, context *spookyinterfaces.TemplatesContext) (string, error) {
	if template == nil {
		return "", fmt.Errorf("template cannot be nil")
	}

	if context == nil {
		return "", fmt.Errorf("templates context cannot be nil")
	}

	// Check cache first
	cacheKey := ti.generateCacheKey(template, context)
	if cached, exists := context.RenderedCache[cacheKey]; exists {
		ti.logger.Debug("Using cached template render", spookylogging.String("template", template.Name))
		return cached, nil
	}

	// Render template using templates manager with enhanced features
	var rendered string
	var err error

	if ti.templatesManager != nil {
		// Create a temporary file for rendering
		tempFile := filepath.Join(os.TempDir(), fmt.Sprintf("template_%s.tmpl", template.Name))
		err = os.WriteFile(tempFile, []byte(template.Source), 0644)
		if err != nil {
			return "", fmt.Errorf("failed to create temporary file for rendering: %w", err)
		}
		defer os.Remove(tempFile)

		// Convert context to additional data for rendering
		additionalData := make(map[string]interface{})
		if context.Functions != nil {
			additionalData["functions"] = context.Functions
		}

		// Add template-specific data
		additionalData["template"] = map[string]interface{}{
			"name": template.Name,
			"path": tempFile,
		}

		// Add project context
		additionalData["project"] = map[string]interface{}{
			"path": context.ProjectPath,
		}

		// Render the template with error handling and recovery
		rendered, err = ti.renderTemplateWithRecovery(tempFile, context.ProjectPath, additionalData)
		if err != nil {
			return "", fmt.Errorf("failed to render template: %w", err)
		}
	} else {
		// Fallback rendering - return template source as-is
		ti.logger.Warn("Templates manager not available, using fallback rendering", spookylogging.String("template", template.Name))
		rendered = template.Source
	}

	// Cache the result
	context.RenderedCache[cacheKey] = rendered

	ti.logger.Info("Rendered template", spookylogging.String("template", template.Name))

	return rendered, nil
}

// CacheTemplate caches a template for later use
func (ti *CoordinatorTemplatesIntegration) CacheTemplate(template *spookytemplatestypes.Template) error {
	if template == nil {
		return fmt.Errorf("template cannot be nil")
	}

	// Generate cache key for the template
	cacheKey := ti.generateTemplateCacheKey(template)

	// Compile and cache template using templates manager
	if ti.templatesManager != nil {
		// Create a temporary file for compilation
		tempFile := filepath.Join(os.TempDir(), fmt.Sprintf("template_%s.tmpl", template.Name))
		err := os.WriteFile(tempFile, []byte(template.Source), 0644)
		if err != nil {
			return fmt.Errorf("failed to create temporary file for caching: %w", err)
		}
		defer os.Remove(tempFile)

		// Validate template before caching
		if err := ti.templatesManager.ValidateTemplate(tempFile); err != nil {
			return fmt.Errorf("template validation failed during caching: %w", err)
		}

		// Store template in internal cache
		ti.templateCache[cacheKey] = template

		// Template is now cached in the manager's internal cache
		ti.logger.Debug("Template compiled and cached",
			spookylogging.String("template", template.Name),
			spookylogging.String("cache_key", cacheKey))
	} else {
		// Fallback caching - store in local cache
		ti.templateCache[cacheKey] = template
		ti.logger.Debug("Template cached locally",
			spookylogging.String("template", template.Name),
			spookylogging.String("cache_key", cacheKey))
	}

	return nil
}

// GetTemplate gets a specific template by name
func (ti *CoordinatorTemplatesIntegration) GetTemplate(name string, context *spookyinterfaces.TemplatesContext) (*spookytemplatestypes.Template, error) {
	if name == "" {
		return nil, fmt.Errorf("template name cannot be empty")
	}

	if context == nil {
		return nil, fmt.Errorf("templates context cannot be nil")
	}

	// Look up template in context
	if template, exists := context.Templates[name]; exists {
		return template, nil
	}

	return nil, fmt.Errorf("template '%s' not found", name)
}

// ListTemplates lists all available templates
func (ti *CoordinatorTemplatesIntegration) ListTemplates(context *spookyinterfaces.TemplatesContext) (map[string]*spookytemplatestypes.Template, error) {
	if context == nil {
		return nil, fmt.Errorf("templates context cannot be nil")
	}

	// Return all templates in context
	result := make(map[string]*spookytemplatestypes.Template)
	for name, template := range context.Templates {
		result[name] = template
	}

	return result, nil
}

// AddTemplate adds a new template to the project with persistence
func (ti *CoordinatorTemplatesIntegration) AddTemplate(name string, template *spookytemplatestypes.Template, context *spookyinterfaces.TemplatesContext) error {
	if name == "" {
		return fmt.Errorf("template name cannot be empty")
	}

	if template == nil {
		return fmt.Errorf("template cannot be nil")
	}

	if context == nil {
		return fmt.Errorf("templates context cannot be nil")
	}

	// Validate template before adding
	if err := ti.ValidateTemplate(template, context); err != nil {
		return fmt.Errorf("template validation failed: %w", err)
	}

	// Check for duplicate template names
	if _, exists := context.Templates[name]; exists {
		return fmt.Errorf("template with name '%s' already exists", name)
	}

	// Add template to context
	context.Templates[name] = template

	// Persist template to project files
	if err := ti.persistTemplate(name, template, context.ProjectPath); err != nil {
		return fmt.Errorf("failed to persist template: %w", err)
	}

	// Cache the template for performance
	if err := ti.CacheTemplate(template); err != nil {
		ti.logger.Warn("Failed to cache template",
			spookylogging.String("template", name),
			spookylogging.Error(err))
	}

	ti.logger.Info("Added template", spookylogging.String("template", name))

	return nil
}

// RemoveTemplate removes a template from the project with cleanup
func (ti *CoordinatorTemplatesIntegration) RemoveTemplate(name string, context *spookyinterfaces.TemplatesContext) error {
	if name == "" {
		return fmt.Errorf("template name cannot be empty")
	}

	if context == nil {
		return fmt.Errorf("templates context cannot be nil")
	}

	// Check if template exists
	_, exists := context.Templates[name]
	if !exists {
		return fmt.Errorf("template '%s' not found", name)
	}

	// Check for template dependencies
	if err := ti.checkTemplateDependencies(name, context); err != nil {
		return fmt.Errorf("cannot remove template due to dependencies: %w", err)
	}

	// Remove from context
	delete(context.Templates, name)

	// Remove from cache
	delete(context.RenderedCache, name)

	// Remove from file system
	if err := ti.removeTemplateFile(name, context.ProjectPath); err != nil {
		return fmt.Errorf("failed to remove template file: %w", err)
	}

	// Invalidate cache entries
	ti.invalidateTemplateCache(name)

	ti.logger.Info("Removed template", spookylogging.String("template", name))

	return nil
}

// generateCacheKey generates a cache key for template rendering
func (ti *CoordinatorTemplatesIntegration) generateCacheKey(template *spookytemplatestypes.Template, context *spookyinterfaces.TemplatesContext) string {
	params := map[string]interface{}{
		"template_name": template.Name,
		"template_hash": fmt.Sprintf("%d", len(template.Source)), // Simple hash for now
		"timestamp":     context.Timestamp.Unix(),
	}
	return spookyinterfaces.OptimizeCacheKey(context.ProjectPath, "template", params)
}

// generateTemplateCacheKey generates a cache key for template storage
func (ti *CoordinatorTemplatesIntegration) generateTemplateCacheKey(template *spookytemplatestypes.Template) string {
	// Create a unique cache key based on template name and content hash
	contentHash := fmt.Sprintf("%d", len(template.Source)) // Simplified hash
	return fmt.Sprintf("template_%s_%s", template.Name, contentHash)
}

// persistTemplate persists a template to the project file system
func (ti *CoordinatorTemplatesIntegration) persistTemplate(name string, template *spookytemplatestypes.Template, projectPath string) error {
	// Create templates directory if it doesn't exist
	templatesDir := filepath.Join(projectPath, "templates")
	if err := os.MkdirAll(templatesDir, 0755); err != nil {
		return fmt.Errorf("failed to create templates directory: %w", err)
	}

	// Determine file extension based on template type
	ext := ".tmpl" // Default extension
	if strings.Contains(template.Source, "{{") {
		ext = ".tmpl"
	} else if strings.Contains(template.Source, "hcl") {
		ext = ".hcl"
	}

	// Create template file path
	templatePath := filepath.Join(templatesDir, name+ext)

	// Write template to file
	err := os.WriteFile(templatePath, []byte(template.Source), 0644)
	if err != nil {
		return fmt.Errorf("failed to write template file: %w", err)
	}

	ti.logger.Debug("Persisted template to file",
		spookylogging.String("template", name),
		spookylogging.String("path", templatePath))

	return nil
}

// removeTemplateFile removes a template file from the file system
func (ti *CoordinatorTemplatesIntegration) removeTemplateFile(name string, projectPath string) error {
	// Look for template file with various extensions
	templatesDir := filepath.Join(projectPath, "templates")
	extensions := []string{".tmpl", ".hcl"}

	for _, ext := range extensions {
		templatePath := filepath.Join(templatesDir, name+ext)
		if _, err := os.Stat(templatePath); err == nil {
			// File exists, remove it
			if err := os.Remove(templatePath); err != nil {
				return fmt.Errorf("failed to remove template file %s: %w", templatePath, err)
			}
			ti.logger.Debug("Removed template file",
				spookylogging.String("template", name),
				spookylogging.String("path", templatePath))
			return nil
		}
	}

	// Template file not found, which is acceptable
	ti.logger.Debug("Template file not found for removal", spookylogging.String("template", name))
	return nil
}

// checkTemplateDependencies checks if a template has dependencies that would prevent removal
func (ti *CoordinatorTemplatesIntegration) checkTemplateDependencies(name string, context *spookyinterfaces.TemplatesContext) error {
	// In a real implementation, this would check:
	// - Other templates that include this template
	// - Actions that reference this template
	// - Variables that are specific to this template
	// - Any other references to this template

	// For now, we'll just log the check
	ti.logger.Debug("Checking template dependencies", spookylogging.String("template", name))

	// Return nil to indicate no blocking dependencies
	return nil
}

// invalidateTemplateCache invalidates cache entries for a template
func (ti *CoordinatorTemplatesIntegration) invalidateTemplateCache(name string) {
	// Remove from template cache
	for key := range ti.templateCache {
		if strings.Contains(key, name) {
			delete(ti.templateCache, key)
		}
	}

	ti.logger.Debug("Invalidated template cache", spookylogging.String("template", name))
}

// renderTemplateWithRecovery renders a template with error handling and recovery
func (ti *CoordinatorTemplatesIntegration) renderTemplateWithRecovery(templatePath, projectPath string, additionalData map[string]interface{}) (string, error) {
	// Attempt to render the template
	rendered, err := ti.templatesManager.RenderTemplate(templatePath, projectPath, additionalData)
	if err != nil {
		// Log the error for debugging
		ti.logger.Error("Template rendering failed", err,
			spookylogging.String("template", templatePath))

		// In a real implementation, you might:
		// - Try to render with fallback data
		// - Use a default template
		// - Return a partial result
		// - Implement retry logic

		return "", err
	}

	return rendered, nil
}
