package facts

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"spooky/internal/schemas"
)

// FormatDetector detects the format of a file based on content and schema validation
type FormatDetector struct {
	validator *schemas.SchemaValidator
}

// NewFormatDetector creates a new format detector
func NewFormatDetector() *FormatDetector {
	return &FormatDetector{
		validator: schemas.NewSchemaValidator(),
	}
}

// DetectFormat determines if a file is JSON or HCL based on content analysis and schema validation
func (d *FormatDetector) DetectFormat(r io.Reader) (string, error) {
	// Read the content
	content, err := io.ReadAll(r)
	if err != nil {
		return "", fmt.Errorf("failed to read file content: %w", err)
	}

	// Try to parse as JSON first
	if d.isValidJSON(content) {
		// Validate against JSON schema
		if d.validateAgainstJSONSchema(content) {
			return "json", nil
		}
	}

	// Try to parse as HCL
	if d.isValidHCL(content) {
		// Validate against HCL schema
		if d.validateAgainstHCLSchema(content) {
			return "hcl", nil
		}
	}

	return "", fmt.Errorf("unable to determine file format: content does not match JSON or HCL schemas")
}

// isValidJSON checks if content is valid JSON
func (d *FormatDetector) isValidJSON(content []byte) bool {
	var js json.RawMessage
	return json.Unmarshal(content, &js) == nil
}

// isValidHCL checks if content is valid HCL
func (d *FormatDetector) isValidHCL(content []byte) bool {
	// Check for HCL-specific patterns
	contentStr := string(content)

	// Look for HCL block patterns - updated for unified schema
	hclPatterns := []string{
		"facts_export = {",
		"global_facts = [",
		"project_facts = [",
		"machine_id =",
		"collected_at =",
		"ttl =",
		"metadata = {",
	}

	for _, pattern := range hclPatterns {
		if strings.Contains(contentStr, pattern) {
			return true
		}
	}

	return false
}

// validateAgainstJSONSchema validates content against JSON facts schema
func (d *FormatDetector) validateAgainstJSONSchema(content []byte) bool {
	// Load the facts structure schema
	if err := d.validator.LoadSchema(schemas.SchemaTypeFactsStructure); err != nil {
		return false
	}

	// Parse JSON content
	var data interface{}
	if err := json.Unmarshal(content, &data); err != nil {
		return false
	}

	// Validate against unified schema
	if err := d.validator.ValidateData(data, "facts_export"); err != nil {
		return false
	}

	return true
}

// validateAgainstHCLSchema validates content against HCL facts schema
func (d *FormatDetector) validateAgainstHCLSchema(content []byte) bool {
	// Load the facts structure schema
	if err := d.validator.LoadSchema(schemas.SchemaTypeFactsStructure); err != nil {
		return false
	}

	// For HCL validation, we need to parse the HCL content
	// This is a simplified validation - in a full implementation,
	// we would parse the HCL AST and validate against the schema

	contentStr := string(content)

	// Check for required HCL structure patterns for unified schema
	requiredPatterns := []string{
		"facts_export = {",
		"metadata = {",
		"global_facts = [",
		"project_facts = [",
	}

	for _, pattern := range requiredPatterns {
		if !strings.Contains(contentStr, pattern) {
			return false
		}
	}

	return true
}
