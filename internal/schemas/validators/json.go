package validators

import (
	"encoding/json"
	"fmt"
	spookytypesschemas "spooky/internal/types/schemas"
)

// JSONValidator implements JSON validation
type JSONValidator struct{}

// NewJSONValidator creates a new JSON validator
func NewJSONValidator() *JSONValidator {
	return &JSONValidator{}
}

// Validate validates JSON data against a schema
func (v *JSONValidator) Validate(data interface{}, schema *spookytypesschemas.Schema) *spookytypesschemas.ValidationResult {
	result := &spookytypesschemas.ValidationResult{
		Valid:  true,
		Schema: string(schema.Type),
	}

	// For JSON validation, we'll check if the data can be marshaled to JSON
	if data != nil {
		if _, err := json.Marshal(data); err != nil {
			result.Valid = false
			result.Errors = append(result.Errors, spookytypesschemas.ValidationError{
				File:     schema.Filename,
				Message:  fmt.Sprintf("invalid JSON data: %v", err),
				Severity: "error",
			})
		}
	}

	return result
}

// GetName returns the validator name
func (v *JSONValidator) GetName() string {
	return "json"
}
