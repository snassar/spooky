package templates

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"spooky/internal/encryption"

	"github.com/pkg/errors"
)

// TemplateRenderer handles template rendering with transparent decryption
type TemplateRenderer struct {
	transparentDecryptor *encryption.TransparentDecryptor
}

// NewTemplateRenderer creates a new template renderer
func NewTemplateRenderer(transparentDecryptor *encryption.TransparentDecryptor) *TemplateRenderer {
	return &TemplateRenderer{
		transparentDecryptor: transparentDecryptor,
	}
}

// TemplateContext represents the data available to templates
type TemplateContext struct {
	// Variables (automatically decrypted)
	Variables map[string]interface{} `json:"variables"`

	// Machines (with decrypted authentication and variables)
	Machines map[string]interface{} `json:"machines"`

	// Facts (collected from machines)
	Facts map[string]interface{} `json:"facts"`

	// Environment variables
	Environment map[string]string `json:"environment"`

	// Project metadata
	Project map[string]interface{} `json:"project"`
}

// RenderTemplate renders a template file with the given context
func (tr *TemplateRenderer) RenderTemplate(templatePath string, context *TemplateContext) (string, error) {
	// Read the template file
	templateContent, err := os.ReadFile(templatePath)
	if err != nil {
		return "", errors.Wrapf(err, "failed to read template file: %s", templatePath)
	}

	// Parse the template
	tmpl, err := template.New(filepath.Base(templatePath)).Parse(string(templateContent))
	if err != nil {
		return "", errors.Wrapf(err, "failed to parse template: %s", templatePath)
	}

	// Prepare the context with decrypted values
	preparedContext, err := tr.prepareContext(context)
	if err != nil {
		return "", errors.Wrap(err, "failed to prepare template context")
	}

	// Execute the template
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, preparedContext); err != nil {
		return "", errors.Wrapf(err, "failed to execute template: %s", templatePath)
	}

	return buf.String(), nil
}

// prepareContext prepares the template context with decrypted values
func (tr *TemplateRenderer) prepareContext(context *TemplateContext) (map[string]interface{}, error) {
	prepared := make(map[string]interface{})

	// Prepare variables with transparent decryption
	if context.Variables != nil {
		decryptedVariables, err := tr.transparentDecryptor.DecryptVariablesMap(context.Variables)
		if err != nil {
			return nil, errors.Wrap(err, "failed to decrypt variables")
		}
		prepared["Variables"] = decryptedVariables
	} else {
		prepared["Variables"] = make(map[string]interface{})
	}

	// Prepare machines with transparent decryption
	if context.Machines != nil {
		decryptedMachines, err := tr.transparentDecryptor.DecryptMachinesMap(context.Machines)
		if err != nil {
			return nil, errors.Wrap(err, "failed to decrypt machines")
		}
		prepared["Machines"] = decryptedMachines
	} else {
		prepared["Machines"] = make(map[string]interface{})
	}

	// Facts don't need decryption (they're collected from machines)
	if context.Facts != nil {
		prepared["Facts"] = context.Facts
	} else {
		prepared["Facts"] = make(map[string]interface{})
	}

	// Environment variables
	if context.Environment != nil {
		prepared["Environment"] = context.Environment
	} else {
		prepared["Environment"] = make(map[string]string)
	}

	// Project metadata
	if context.Project != nil {
		prepared["Project"] = context.Project
	} else {
		prepared["Project"] = make(map[string]interface{})
	}

	// Add flattened variables for easy access
	// This allows templates to use {{ .database_password }} instead of {{ .Variables.database_password }}
	if variables, ok := prepared["Variables"].(map[string]interface{}); ok {
		for key, value := range variables {
			prepared[key] = value
		}
	}

	return prepared, nil
}

// RenderTemplateToFile renders a template and writes the result to a file
func (tr *TemplateRenderer) RenderTemplateToFile(templatePath, outputPath string, context *TemplateContext) error {
	// Render the template
	rendered, err := tr.RenderTemplate(templatePath, context)
	if err != nil {
		return errors.Wrap(err, "failed to render template")
	}

	// Create output directory if it doesn't exist
	outputDir := filepath.Dir(outputPath)
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return errors.Wrapf(err, "failed to create output directory: %s", outputDir)
	}

	// Write the rendered content to file
	if err := os.WriteFile(outputPath, []byte(rendered), 0o644); err != nil {
		return errors.Wrapf(err, "failed to write rendered template to: %s", outputPath)
	}

	return nil
}

// RenderTemplateString renders a template string with the given context
func (tr *TemplateRenderer) RenderTemplateString(templateString string, context *TemplateContext) (string, error) {
	// Parse the template string
	tmpl, err := template.New("template").Parse(templateString)
	if err != nil {
		return "", errors.Wrap(err, "failed to parse template string")
	}

	// Prepare the context with decrypted values
	preparedContext, err := tr.prepareContext(context)
	if err != nil {
		return "", errors.Wrap(err, "failed to prepare template context")
	}

	// Execute the template
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, preparedContext); err != nil {
		return "", errors.Wrap(err, "failed to execute template string")
	}

	return buf.String(), nil
}

// ValidateTemplate validates a template file without rendering it
func (tr *TemplateRenderer) ValidateTemplate(templatePath string) error {
	// Read the template file
	templateContent, err := os.ReadFile(templatePath)
	if err != nil {
		return errors.Wrapf(err, "failed to read template file: %s", templatePath)
	}

	// Try to parse the template
	_, err = template.New(filepath.Base(templatePath)).Parse(string(templateContent))
	if err != nil {
		return errors.Wrapf(err, "failed to parse template: %s", templatePath)
	}

	return nil
}

// GetTemplateFunctions returns custom template functions
func (tr *TemplateRenderer) GetTemplateFunctions() template.FuncMap {
	return template.FuncMap{
		"join":    strings.Join,
		"split":   strings.Split,
		"replace": strings.Replace,
		"upper":   strings.ToUpper,
		"lower":   strings.ToLower,
		"trim":    strings.TrimSpace,
		"contains": func(s, substr string) bool {
			return strings.Contains(s, substr)
		},
		"hasPrefix": func(s, prefix string) bool {
			return strings.HasPrefix(s, prefix)
		},
		"hasSuffix": func(s, suffix string) bool {
			return strings.HasSuffix(s, suffix)
		},
		"default": func(defaultValue, value interface{}) interface{} {
			if value == nil || value == "" {
				return defaultValue
			}
			return value
		},
		"env": func(key string) string {
			return os.Getenv(key)
		},
		"envOrDefault": func(key, defaultValue string) string {
			if value := os.Getenv(key); value != "" {
				return value
			}
			return defaultValue
		},
	}
}

// CreateTemplateContext creates a new template context with the given data
func CreateTemplateContext(variables, machines, facts map[string]interface{}, environment map[string]string, project map[string]interface{}) *TemplateContext {
	return &TemplateContext{
		Variables:   variables,
		Machines:    machines,
		Facts:       facts,
		Environment: environment,
		Project:     project,
	}
}
