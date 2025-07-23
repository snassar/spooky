package facts

import (
	"fmt"
	"regexp"
)

// ValidationError represents a validation error
type ValidationError struct {
	Field   string      `json:"field"`
	Message string      `json:"message"`
	Value   interface{} `json:"value,omitempty"`
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error in %s: %s", e.Field, e.Message)
}

// ValidationResult contains validation results
type ValidationResult struct {
	Valid    bool               `json:"valid"`
	Errors   []*ValidationError `json:"errors,omitempty"`
	Warnings []*ValidationError `json:"warnings,omitempty"`
}

// ValidateCustomFacts validates custom facts format
func ValidateCustomFacts(facts map[string]*CustomFacts) *ValidationResult {
	result := &ValidationResult{Valid: true}

	for serverID, customFacts := range facts {
		// Validate server ID
		if err := validateServerID(serverID); err != nil {
			result.Valid = false
			result.Errors = append(result.Errors, &ValidationError{
				Field:   "server_id",
				Message: err.Error(),
				Value:   serverID,
			})
		}

		// Validate custom facts
		if customFacts.Custom != nil {
			if err := validateCustomFactStructure(customFacts.Custom); err != nil {
				result.Valid = false
				result.Errors = append(result.Errors, &ValidationError{
					Field:   "custom",
					Message: err.Error(),
					Value:   customFacts.Custom,
				})
			}
		}

		// Validate overrides
		if customFacts.Overrides != nil {
			if err := validateOverrideStructure(customFacts.Overrides); err != nil {
				result.Valid = false
				result.Errors = append(result.Errors, &ValidationError{
					Field:   "overrides",
					Message: err.Error(),
					Value:   customFacts.Overrides,
				})
			}
		}
	}

	return result
}

// validateServerID validates server identifier
func validateServerID(serverID string) error {
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

// validateCustomFactStructure validates custom fact structure
func validateCustomFactStructure(custom map[string]interface{}) error {
	for category, facts := range custom {
		if category == "" {
			return fmt.Errorf("category name cannot be empty")
		}

		if factsMap, ok := facts.(map[string]interface{}); ok {
			for key, value := range factsMap {
				if key == "" {
					return fmt.Errorf("fact key cannot be empty in category %s", category)
				}

				if value == nil {
					return fmt.Errorf("fact value cannot be nil for %s.%s", category, key)
				}
			}
		} else {
			return fmt.Errorf("category %s must be an object", category)
		}
	}

	return nil
}

// validateOverrideStructure validates override structure
func validateOverrideStructure(overrides map[string]interface{}) error {
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
