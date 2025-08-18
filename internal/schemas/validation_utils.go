// Package schemas provides shared validation utilities
package schemas

import (
	"fmt"
	"regexp"
	"strings"

	spookytypesschemas "spooky/internal/types/schemas"
)

// Shared validation constants
const (
	ValidationEmail    = "email"
	ValidationRequired = "required"
	ValidationError    = "error"
	ValidationWarning  = "warning"
)

// ValidationUtils provides shared validation functionality
type ValidationUtils struct{}

// NewValidationUtils creates a new validation utilities instance
func NewValidationUtils() *ValidationUtils {
	return &ValidationUtils{}
}

// ValidateStringConstraints validates string-specific constraints
func (u *ValidationUtils) ValidateStringConstraints(
	constraints *spookytypesschemas.FieldConstraints,
	strValue string,
	fieldPath string,
	result *spookytypesschemas.ValidationResult,
	errorCreator func(code, message, suggestion, severity string),
) error {
	// Length constraints
	if constraints.MinLength != nil && len(strValue) < *constraints.MinLength {
		errorCreator(
			"string_too_short",
			fmt.Sprintf("Field '%s' length %d is less than minimum %d", fieldPath, len(strValue), *constraints.MinLength),
			fmt.Sprintf("Increase the length of field '%s' to at least %d characters", fieldPath, *constraints.MinLength),
			ValidationError,
		)
	}

	if constraints.MaxLength != nil && len(strValue) > *constraints.MaxLength {
		errorCreator(
			"string_too_long",
			fmt.Sprintf("Field '%s' length %d exceeds maximum %d", fieldPath, len(strValue), *constraints.MaxLength),
			fmt.Sprintf("Reduce the length of field '%s' to at most %d characters", fieldPath, *constraints.MaxLength),
			ValidationError,
		)
	}

	// Pattern constraint
	if constraints.Pattern != nil {
		if err := u.ValidatePatternConstraint(constraints.Pattern, strValue, fieldPath, result, errorCreator); err != nil {
			return err
		}
	}

	// Format constraint
	if constraints.Format != nil {
		if err := u.ValidateFormat(*constraints.Format, strValue, fieldPath, result, errorCreator); err != nil {
			return err
		}
	}

	return nil
}

// ValidateNumericConstraints validates numeric-specific constraints
func (u *ValidationUtils) ValidateNumericConstraints(
	constraints *spookytypesschemas.FieldConstraints,
	numValue float64,
	fieldPath string,
	result *spookytypesschemas.ValidationResult,
	errorCreator func(code, message, suggestion, severity string),
) error {
	if constraints.Min != nil && numValue < *constraints.Min {
		errorCreator(
			"number_too_small",
			fmt.Sprintf("Field '%s' value %f is less than minimum %f", fieldPath, numValue, *constraints.Min),
			fmt.Sprintf("Increase the value of field '%s' to at least %f", fieldPath, *constraints.Min),
			ValidationError,
		)
	}

	if constraints.Max != nil && numValue > *constraints.Max {
		errorCreator(
			"number_too_large",
			fmt.Sprintf("Field '%s' value %f exceeds maximum %f", fieldPath, numValue, *constraints.Max),
			fmt.Sprintf("Reduce the value of field '%s' to at most %f", fieldPath, *constraints.Max),
			ValidationError,
		)
	}

	return nil
}

// ValidateEnumConstraints validates enum constraints
func (u *ValidationUtils) ValidateEnumConstraints(
	constraints *spookytypesschemas.FieldConstraints,
	value interface{},
	fieldPath string,
	result *spookytypesschemas.ValidationResult,
	errorCreator func(code, message, suggestion, severity string),
) error {
	if len(constraints.Enum) > 0 {
		found := false
		for _, enumValue := range constraints.Enum {
			if value == enumValue {
				found = true
				break
			}
		}
		if !found {
			errorCreator(
				"invalid_enum_value",
				fmt.Sprintf("Field '%s' value '%v' is not in allowed enum values", fieldPath, value),
				fmt.Sprintf("Use one of the allowed values for field '%s': %v", fieldPath, constraints.Enum),
				ValidationError,
			)
		}
	}

	return nil
}

// ValidatePatternConstraint validates pattern constraint
func (u *ValidationUtils) ValidatePatternConstraint(
	pattern *string,
	strValue string,
	fieldPath string,
	result *spookytypesschemas.ValidationResult,
	errorCreator func(code, message, suggestion, severity string),
) error {
	matched, err := regexp.MatchString(*pattern, strValue)
	if err != nil {
		errorCreator(
			"invalid_regex",
			fmt.Sprintf("Invalid regex pattern for field '%s': %v", fieldPath, err),
			"",
			ValidationError,
		)
		return err
	}

	if !matched {
		errorCreator(
			"pattern_mismatch",
			fmt.Sprintf("Field '%s' value '%s' does not match pattern '%s'", fieldPath, strValue, *pattern),
			fmt.Sprintf("Ensure field '%s' matches the required pattern", fieldPath),
			ValidationError,
		)
	}

	return nil
}

// ValidateFormat validates format constraints
func (u *ValidationUtils) ValidateFormat(
	format, value, fieldPath string,
	result *spookytypesschemas.ValidationResult,
	errorCreator func(code, message, suggestion, severity string),
) error {
	switch format {
	case ValidationEmail:
		emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
		if !emailRegex.MatchString(value) {
			errorCreator(
				"invalid_email",
				fmt.Sprintf("Field '%s' value '%s' is not a valid email address", fieldPath, value),
				"Enter a valid email address (e.g., user@example.com)",
				ValidationError,
			)
		}
	case "uri":
		if !strings.HasPrefix(value, "http://") && !strings.HasPrefix(value, "https://") {
			errorCreator(
				"invalid_uri",
				fmt.Sprintf("Field '%s' value '%s' is not a valid URI", fieldPath, value),
				"Enter a valid URI starting with http:// or https://",
				ValidationError,
			)
		}
	case "ipv4":
		ipv4Regex := regexp.MustCompile(`^(\d{1,3}\.){3}\d{1,3}$`)
		if !ipv4Regex.MatchString(value) {
			errorCreator(
				"invalid_ipv4",
				fmt.Sprintf("Field '%s' value '%s' is not a valid IPv4 address", fieldPath, value),
				"Enter a valid IPv4 address (e.g., 192.168.1.1)",
				ValidationError,
			)
		}
	}
	return nil
}

// ValidateFieldConstraints validates field constraints for any value type
func (u *ValidationUtils) ValidateFieldConstraints(
	constraints *spookytypesschemas.FieldConstraints,
	value interface{},
	fieldPath string,
	result *spookytypesschemas.ValidationResult,
	errorCreator func(code, message, suggestion, severity string),
	toFloat64 func(interface{}) (float64, bool),
) error {
	// String constraints
	if strValue, ok := value.(string); ok {
		if err := u.ValidateStringConstraints(constraints, strValue, fieldPath, result, errorCreator); err != nil {
			return err
		}
	}

	// Numeric constraints
	if numValue, ok := toFloat64(value); ok {
		if err := u.ValidateNumericConstraints(constraints, numValue, fieldPath, result, errorCreator); err != nil {
			return err
		}
	}

	// Enum constraints
	if err := u.ValidateEnumConstraints(constraints, value, fieldPath, result, errorCreator); err != nil {
		return err
	}

	return nil
}
