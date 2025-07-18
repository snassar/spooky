package config

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Validator represents the comprehensive validator
type Validator struct {
	validate *validator.Validate
}

// NewValidator creates a new validator with all custom validation functions registered
func NewValidator() *Validator {
	v := &Validator{
		validate: validator.New(),
	}
	v.registerCustomValidations()
	return v
}

// Global validator instance
var globalValidator *Validator

func init() {
	globalValidator = NewValidator()
}

// ValidateConfig validates the configuration structure
func ValidateConfig(config *Config) error {
	return globalValidator.validateConfig(config)
}

// registerCustomValidations registers all custom validation functions
func (v *Validator) registerCustomValidations() {
	// Register struct-level validations for cross-field validation
	v.validate.RegisterStructValidation(v.validateServerStruct, Server{})
	v.validate.RegisterStructValidation(v.validateActionStruct, Action{})
	v.validate.RegisterStructValidation(v.validateConfigStruct, Config{})
}

// validateServerStruct performs struct-level validation for Server
func (v *Validator) validateServerStruct(sl validator.StructLevel) {
	server := sl.Current().Interface().(Server)

	// Validate authentication requirements
	if server.Password == "" && server.KeyFile == "" {
		sl.ReportError(server.Password, "Password", "password", "server_auth", server.Name)
	}
}

// validateActionStruct performs struct-level validation for Action
func (v *Validator) validateActionStruct(sl validator.StructLevel) {
	action := sl.Current().Interface().(Action)

	// Validate execution requirements
	if action.Command == "" && action.Script == "" {
		sl.ReportError(action.Command, "Command", "command", "action_exec", action.Name)
	}
	if action.Command != "" && action.Script != "" {
		sl.ReportError(action.Command, "Command", "command", "action_exec", action.Name)
	}
}

// validateConfigStruct performs struct-level validation for Config
func (v *Validator) validateConfigStruct(sl validator.StructLevel) {
	config := sl.Current().Interface().(Config)

	// Validate unique server names
	serverNames := make(map[string]bool)
	for _, server := range config.Servers {
		if serverNames[server.Name] {
			sl.ReportError(server.Name, "Name", "name", "unique_server", server.Name)
		}
		serverNames[server.Name] = true
	}

	// Validate unique action names
	actionNames := make(map[string]bool)
	for i := range config.Actions {
		action := &config.Actions[i]
		if actionNames[action.Name] {
			sl.ReportError(action.Name, "Name", "name", "unique_action", action.Name)
		}
		actionNames[action.Name] = true
	}

	// Validate server references in actions
	for i := range config.Actions {
		action := &config.Actions[i]
		for _, serverRef := range action.Servers {
			if !serverNames[serverRef] {
				sl.ReportError(serverRef, "Servers", "servers", "valid_servers", action.Name)
			}
		}
	}
}

// validateConfig validates the entire configuration
func (v *Validator) validateConfig(config *Config) error {
	// Set defaults before validation
	SetDefaults(config)

	// Perform validation
	if err := v.validate.Struct(config); err != nil {
		return v.formatValidationErrors(err)
	}

	return nil
}

// formatValidationErrors converts validator errors to user-friendly messages
func (v *Validator) formatValidationErrors(err error) error {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		var messages []string
		for _, e := range validationErrors {
			message := v.formatValidationError(e)
			if message != "" {
				messages = append(messages, message)
			}
		}
		if len(messages) > 0 {
			return fmt.Errorf("validation failed: %s", strings.Join(messages, "; "))
		}
	}
	return err
}

// formatValidationError formats a single validation error
func (v *Validator) formatValidationError(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", e.Field())
	case "min":
		if e.Field() == "Servers" {
			return "at least one server must be defined"
		}
		return fmt.Sprintf("%s must be at least %s characters long", e.Field(), e.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s", e.Field(), e.Param())
	case "server_auth":
		return fmt.Sprintf("either password or key_file must be specified for server %s", e.Param())
	case "action_exec":
		return fmt.Sprintf("either command or script must be specified for action %s (but not both)", e.Param())
	case "unique_server":
		return fmt.Sprintf("duplicate server name: %s", e.Param())
	case "unique_action":
		return fmt.Sprintf("duplicate action name: %s", e.Param())
	case "valid_port":
		return fmt.Sprintf("port must be between 1 and 65535 for server %s", e.Param())
	case "valid_timeout":
		return fmt.Sprintf("timeout must be between 1 and 3600 seconds for action %s", e.Param())
	case "valid_servers":
		return fmt.Sprintf("server reference '%s' in action '%s' does not exist", e.Value(), e.Param())
	default:
		return fmt.Sprintf("%s failed validation: %s", e.Field(), e.Tag())
	}
}

// ValidateServer validates a single server
func (v *Validator) ValidateServer(server *Server) error {
	if err := v.validate.Struct(server); err != nil {
		return v.formatValidationErrors(err)
	}
	return nil
}

// ValidateAction validates a single action
func (v *Validator) ValidateAction(action *Action) error {
	if err := v.validate.Struct(action); err != nil {
		return v.formatValidationErrors(err)
	}
	return nil
}
