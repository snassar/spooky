package config

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
	// Register custom struct-level validation for Server
	validate.RegisterStructValidation(ServerStructLevelValidation, Server{})
	// Register custom struct-level validation for Action
	validate.RegisterStructValidation(ActionStructLevelValidation, Action{})
}

// Custom validation for Server: required fields and either Password or KeyFile must be set
func ServerStructLevelValidation(sl validator.StructLevel) {
	server := sl.Current().Interface().(Server)
	if server.Name == "" {
		sl.ReportError(server.Name, "Name", "name", "namerequired", "server")
	}
	if server.Host == "" {
		sl.ReportError(server.Host, "Host", "host", "hostrequired", server.Name)
	}
	if server.User == "" {
		sl.ReportError(server.User, "User", "user", "userrequired", server.Name)
	}
	if server.Password == "" && server.KeyFile == "" {
		sl.ReportError(server.Password, "Password", "password", "passwordorkeyfile", server.Name)
	}
}

// Custom validation for Action: required Name, and either Command or Script, but not both or neither
func ActionStructLevelValidation(sl validator.StructLevel) {
	action := sl.Current().Interface().(Action)
	if action.Name == "" {
		sl.ReportError(action.Name, "Name", "name", "namerequired", "action")
	}
	if action.Command == "" && action.Script == "" {
		sl.ReportError(action.Command, "Command", "command", "commandorscript", action.Name)
	}
	if action.Command != "" && action.Script != "" {
		sl.ReportError(action.Command, "Command", "command", "commandorscriptmutual", action.Name)
	}
}

// ValidateConfig validates the configuration structure using schema-based validation
func ValidateConfig(config *Config) error {
	// Set defaults before validation
	SetDefaults(config)

	// Use schema-based validation for basic field validation and custom rules
	if err := validate.Struct(config); err != nil {
		return formatValidationErrors(err)
	}

	// Validate unique names
	if err := validateUniqueNames(config); err != nil {
		return err
	}

	return nil
}

// validateUniqueNames ensures server and action names are unique
func validateUniqueNames(config *Config) error {
	serverNames := make(map[string]bool)
	for _, server := range config.Servers {
		if serverNames[server.Name] {
			return fmt.Errorf("duplicate server name: %s", server.Name)
		}
		serverNames[server.Name] = true
	}

	actionNames := make(map[string]bool)
	for i := range config.Actions {
		action := &config.Actions[i]
		if actionNames[action.Name] {
			return fmt.Errorf("duplicate action name: %s", action.Name)
		}
		actionNames[action.Name] = true
	}

	return nil
}

// formatValidationErrors converts validator errors to user-friendly messages
func formatValidationErrors(err error) error {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		var messages []string
		for _, e := range validationErrors {
			switch e.Tag() {
			case "namerequired":
				switch e.Param() {
				case "server":
					messages = append(messages, "server name cannot be empty")
				case "action":
					messages = append(messages, "action name cannot be empty")
				default:
					messages = append(messages, "name cannot be empty")
				}
			case "hostrequired":
				messages = append(messages, fmt.Sprintf("server host cannot be empty for server %s", e.Param()))
			case "userrequired":
				messages = append(messages, fmt.Sprintf("server user cannot be empty for server %s", e.Param()))
			case "passwordorkeyfile":
				messages = append(messages, fmt.Sprintf("either password or key_file must be specified for server %s", e.Param()))
			case "commandorscript":
				messages = append(messages, fmt.Sprintf("either command or script must be specified for action %s", e.Param()))
			case "commandorscriptmutual":
				messages = append(messages, fmt.Sprintf("cannot specify both command and script for action %s", e.Param()))
			case "min":
				if e.Field() == "Servers" {
					messages = append(messages, "at least one server must be defined")
				} else {
					messages = append(messages, fmt.Sprintf("%s must be at least %s", e.Field(), e.Param()))
				}
			case "max":
				messages = append(messages, fmt.Sprintf("%s must be at most %s", e.Field(), e.Param()))
			default:
				messages = append(messages, fmt.Sprintf("%s failed validation: %s", e.Field(), e.Tag()))
			}
		}
		return fmt.Errorf("validation failed: %s", strings.Join(messages, "; "))
	}
	return err
}
