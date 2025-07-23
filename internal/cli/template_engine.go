package cli

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"spooky/internal/facts"
)

// TemplateEngine handles template rendering with fact integration
type TemplateEngine struct {
	manager *facts.Manager
}

// NewTemplateEngine creates a new template engine
func NewTemplateEngine(manager *facts.Manager) *TemplateEngine {
	return &TemplateEngine{
		manager: manager,
	}
}

// TemplateData represents the data available in templates
type TemplateData struct {
	// System facts
	System map[string]interface{}

	// Custom facts
	Custom map[string]interface{}

	// Environment variables
	Env map[string]string

	// Additional data
	Data map[string]interface{}
}

// RenderTemplate renders a template with fact data
//
//nolint:gocritic // paramTypeCombine - function signature is correct and follows Go conventions
func (te *TemplateEngine) RenderTemplate(templateFile string, server string, additionalData map[string]interface{}) (string, error) {
	// Read template file
	content, err := os.ReadFile(templateFile)
	if err != nil {
		return "", fmt.Errorf("error reading template file: %w", err)
	}

	// Preprocess template content to handle single quotes in template functions
	// Convert {{custom 'path'}} to {{custom "path"}} for Go template parsing
	contentStr := string(content)
	contentStr = te.preprocessTemplateContent(contentStr)

	// Create template data
	data := &TemplateData{
		System: make(map[string]interface{}),
		Custom: make(map[string]interface{}),
		Env:    make(map[string]string),
		Data:   additionalData,
	}

	// Load system facts if manager is available
	if te.manager != nil {
		// Get system facts
		if collection, err := te.manager.CollectAllFacts(server); err == nil {
			for key, fact := range collection.Facts {
				data.System[key] = fact.Value
			}
		}

		// Get custom facts
		if customFacts, err := te.manager.GetCustomFacts(server); err == nil {
			data.Custom = customFacts
		}
	}

	// Load environment variables
	for _, env := range os.Environ() {
		pair := strings.SplitN(env, "=", 2)
		if len(pair) == 2 {
			data.Env[pair[0]] = pair[1]
		}
	}

	// Create template with custom functions
	tmpl, err := template.New(filepath.Base(templateFile)).Funcs(template.FuncMap{
		"custom": func(path string) interface{} {
			return te.getNestedValue(data.Custom, path)
		},
		"system": func(path string) interface{} {
			return te.getNestedValue(data.System, path)
		},
		"env": func(key string) string {
			return data.Env[key]
		},
		"data": func(path string) interface{} {
			return te.getNestedValue(data.Data, path)
		},
	}).Parse(contentStr)

	if err != nil {
		return "", fmt.Errorf("error parsing template: %w", err)
	}

	// Render template
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("error rendering template: %w", err)
	}

	return buf.String(), nil
}

// preprocessTemplateContent converts single quotes to double quotes in template functions
func (te *TemplateEngine) preprocessTemplateContent(content string) string {
	// This is a simple regex-like replacement for template functions
	// Convert {{custom 'path'}} to {{custom "path"}}

	// Find all template function calls with single quotes and convert them
	// This is a simplified approach - in production you might want a more robust parser

	// Replace {{custom '...'}} with {{custom "..."}}
	content = strings.ReplaceAll(content, "{{custom '", "{{custom \"")
	content = strings.ReplaceAll(content, "'}}", "\"}}")

	// Replace {{system '...'}} with {{system "..."}}
	content = strings.ReplaceAll(content, "{{system '", "{{system \"")

	// Replace {{env '...'}} with {{env "..."}}
	content = strings.ReplaceAll(content, "{{env '", "{{env \"")

	// Replace {{data '...'}} with {{data "..."}}
	content = strings.ReplaceAll(content, "{{data '", "{{data \"")

	return content
}

// getNestedValue gets a nested value from a map using dot notation
func (te *TemplateEngine) getNestedValue(data map[string]interface{}, path string) interface{} {
	parts := strings.Split(path, ".")
	current := data

	for i, part := range parts {
		if current == nil {
			return nil
		}

		if i == len(parts)-1 {
			// Last part - return the value
			return current[part]
		}

		// Navigate deeper
		if next, ok := current[part].(map[string]interface{}); ok {
			current = next
		} else {
			return nil
		}
	}

	return nil
}
