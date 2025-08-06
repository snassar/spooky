package engine

import (
	"bytes"
	"context"
	"fmt"
	"text/template"
	"time"
)

// Renderer implements TemplateRenderer interface
type Renderer struct{}

// NewTemplateRenderer creates a new template renderer
func NewTemplateRenderer() TemplateRenderer {
	return &Renderer{}
}

// Render renders a template with data
func (r *Renderer) Render(tmpl *template.Template, data interface{}) (string, error) {
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to render template: %w", err)
	}
	return buf.String(), nil
}

// RenderWithTimeout renders a template with timeout
func (r *Renderer) RenderWithTimeout(tmpl *template.Template, data interface{}, timeout time.Duration) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	var result string
	var err error

	// Create a channel to receive the result
	done := make(chan struct{})

	go func() {
		result, err = r.Render(tmpl, data)
		close(done)
	}()

	// Wait for either completion or timeout
	select {
	case <-done:
		return result, err
	case <-ctx.Done():
		return "", fmt.Errorf("template rendering timed out after %v", timeout)
	}
}

// ValidateData validates template data
func (r *Renderer) ValidateData(data interface{}) error {
	// Basic validation - data should not be nil
	if data == nil {
		return fmt.Errorf("template data is nil")
	}

	// Additional validation can be added here
	return nil
}
