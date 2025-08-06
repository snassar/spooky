package engine

import (
	"fmt"
	"os"
	"text/template"
)

// Parser implements TemplateParser interface
type Parser struct{}

// NewTemplateParser creates a new template parser
func NewTemplateParser() TemplateParser {
	return &Parser{}
}

// Parse parses template content
func (p *Parser) Parse(content []byte, name string) (*template.Template, error) {
	tmpl, err := template.New(name).Parse(string(content))
	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %w", err)
	}
	return tmpl, nil
}

// ParseFile parses a template file
func (p *Parser) ParseFile(path string) (*template.Template, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read template file: %w", err)
	}

	return p.Parse(content, path)
}

// ValidateSyntax validates template syntax
func (p *Parser) ValidateSyntax(content []byte) error {
	_, err := template.New("validation").Parse(string(content))
	if err != nil {
		return fmt.Errorf("template syntax error: %w", err)
	}
	return nil
}
