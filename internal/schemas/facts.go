package schemas

import (
	"fmt"
	"regexp"

	"spooky/internal/facts/types"
)

// FactsValidator validates facts using schemas
type FactsValidator struct {
	baseValidator *SchemaValidator
}

// NewFactsValidator creates a new facts validator
func NewFactsValidator() *FactsValidator {
	return &FactsValidator{
		baseValidator: NewSchemaValidator(),
	}
}

// ValidateCustomFacts validates custom facts format
func (fv *FactsValidator) ValidateCustomFacts(facts map[string]*types.FactCollection) *ValidationResult {
	result := &ValidationResult{
		Valid:  true,
		Schema: "custom-facts",
	}

	for serverID, factCollection := range facts {
		// Validate server ID
		if err := fv.ValidateServerID(serverID); err != nil {
			result.Valid = false
			result.Errors = append(result.Errors, ValidationError{
				Field:    "server_id",
				Message:  err.Error(),
				Value:    serverID,
				Severity: "error",
			})
		}

		// Validate custom facts
		if factCollection.CustomFacts != nil {
			if err := fv.ValidateCustomFactsStructure(factCollection.CustomFacts); err != nil {
				result.Valid = false
				result.Errors = append(result.Errors, ValidationError{
					Field:    "custom",
					Message:  err.Error(),
					Value:    fmt.Sprintf("%v", factCollection.CustomFacts),
					Severity: "error",
				})
			}
		}
	}

	return result
}

// ValidateServerID validates server identifier
func (fv *FactsValidator) ValidateServerID(serverID string) error {
	if serverID == "" {
		return fmt.Errorf("server ID cannot be empty")
	}

	// Check for valid characters
	validPattern := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	if !validPattern.MatchString(serverID) {
		return fmt.Errorf("server ID contains invalid characters")
	}

	return nil
}

// ValidateCustomFactsStructure validates custom facts structure
func (fv *FactsValidator) ValidateCustomFactsStructure(custom map[string]map[string]interface{}) error {
	for filename, facts := range custom {
		if filename == "" {
			return fmt.Errorf("filename cannot be empty")
		}

		for key, value := range facts {
			if key == "" {
				return fmt.Errorf("fact key cannot be empty in file %s", filename)
			}

			if value == nil {
				return fmt.Errorf("fact value cannot be nil for %s.%s", filename, key)
			}
		}
	}

	return nil
}

// ValidateOverrideStructure validates override structure
func (fv *FactsValidator) ValidateOverrideStructure(overrides map[string]interface{}) error {
	for category, facts := range overrides {
		if category == "" {
			return fmt.Errorf("override category cannot be empty")
		}

		if factsMap, ok := facts.(map[string]interface{}); ok {
			for key, value := range factsMap {
				if key == "" {
					return fmt.Errorf("override key cannot be empty in category %s", category)
				}

				if value == nil {
					return fmt.Errorf("override value cannot be nil for %s.%s", category, key)
				}
			}
		} else {
			return fmt.Errorf("override category %s must be an object", category)
		}
	}

	return nil
}

// ValidateFactsWithSchema validates facts with schema composition
func (fv *FactsValidator) ValidateFactsWithSchema(factsPath string, format string) *ValidationResult {
	result := &ValidationResult{
		Valid:  true,
		Schema: "facts-with-schema",
	}

	// Use the base validator for facts validation
	factsResult := fv.baseValidator.ValidateFacts(factsPath, format)
	if !factsResult.Valid {
		result.Valid = false
		result.Errors = append(result.Errors, factsResult.Errors...)
		result.Warnings = append(result.Warnings, factsResult.Warnings...)
	}

	// Add schema composition validation if needed
	if err := fv.composeFactsSchema(format); err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, ValidationError{
			Field:    "facts-schema",
			Message:  fmt.Sprintf("Facts schema composition failed: %v", err),
			Severity: "error",
		})
	}

	return result
}

// composeFactsSchema composes facts schema for different storage formats
func (fv *FactsValidator) composeFactsSchema(storageFormat string) error {
	// This is a placeholder for schema composition logic
	// In a real implementation, this would compose the base facts schema
	// with format-specific validation rules
	switch storageFormat {
	case "hcl":
		// Compose with HCL-specific rules
		return nil
	case "json":
		// Compose with JSON-specific rules
		return nil
	case "badgerdb":
		// Compose with BadgerDB-specific rules
		return nil
	default:
		return fmt.Errorf("unsupported storage format: %s", storageFormat)
	}
}

// ConvertFactsValidationResult converts facts validation results to unified format
// This is a placeholder since facts validation is now handled in the facts package
func (fv *FactsValidator) ConvertFactsValidationResult(factsResult interface{}) *ValidationResult {
	return ConvertFactsValidationResult(factsResult)
}
