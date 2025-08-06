package environment

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Validator implements EnvironmentValidator interface
type Validator struct {
	validationRules map[string]string
}

// NewEnvironmentValidator creates a new environment validator
func NewEnvironmentValidator() *Validator {
	return &Validator{
		validationRules: make(map[string]string),
	}
}

// ValidateVariable validates an environment variable
func (v *Validator) ValidateVariable(name, value string) error {
	if name == "" {
		return fmt.Errorf("environment variable name cannot be empty")
	}

	// Check for invalid characters in name
	if !regexp.MustCompile(`^[A-Z_][A-Z0-9_]*$`).MatchString(name) {
		return fmt.Errorf("environment variable name must be uppercase letters, numbers, and underscores only")
	}

	// Check for reserved names
	reservedNames := []string{"PATH", "HOME", "USER", "SHELL", "TERM"}
	for _, reserved := range reservedNames {
		if strings.EqualFold(name, reserved) {
			return fmt.Errorf("environment variable name '%s' is reserved", name)
		}
	}

	// Basic value validation - ensure it's not empty for required variables
	if value == "" && strings.HasSuffix(name, "_REQUIRED") {
		return fmt.Errorf("environment variable '%s' is required but has empty value", name)
	}

	return nil
}

// ValidateVariableType validates an environment variable type
func (v *Validator) ValidateVariableType(name, value, expectedType string) error {
	switch expectedType {
	case "string":
		// String is always valid
		return nil
	case "int":
		if _, err := strconv.Atoi(value); err != nil {
			return fmt.Errorf("environment variable '%s' must be an integer", name)
		}
	case "bool":
		validBools := []string{"true", "false", "1", "0", "yes", "no", "on", "off"}
		found := false
		for _, valid := range validBools {
			if strings.EqualFold(value, valid) {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("environment variable '%s' must be a boolean value", name)
		}
	case "float":
		if _, err := strconv.ParseFloat(value, 64); err != nil {
			return fmt.Errorf("environment variable '%s' must be a float", name)
		}
	default:
		return fmt.Errorf("unsupported type: %s", expectedType)
	}

	return nil
}

// ValidateVariableFormat validates an environment variable format
func (v *Validator) ValidateVariableFormat(name, value, format string) error {
	switch format {
	case "email":
		emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
		if !emailRegex.MatchString(value) {
			return fmt.Errorf("environment variable '%s' must be a valid email address", name)
		}
	case "url":
		urlRegex := regexp.MustCompile(`^https?://[^\s/$.?#].\S*$`)
		if !urlRegex.MatchString(value) {
			return fmt.Errorf("environment variable '%s' must be a valid URL", name)
		}
	case "path":
		if value == "" {
			return fmt.Errorf("environment variable '%s' path cannot be empty", name)
		}
		// Basic path validation - could be enhanced
		if strings.Contains(value, "..") {
			return fmt.Errorf("environment variable '%s' path contains invalid '..'", name)
		}
	case "ip":
		ipRegex := regexp.MustCompile(`^(\d{1,3}\.){3}\d{1,3}$`)
		if !ipRegex.MatchString(value) {
			return fmt.Errorf("environment variable '%s' must be a valid IP address", name)
		}
		// Additional IP validation could be added here
	case "port":
		if port, err := strconv.Atoi(value); err != nil || port < 1 || port > 65535 {
			return fmt.Errorf("environment variable '%s' must be a valid port number (1-65535)", name)
		}
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}

	return nil
}
