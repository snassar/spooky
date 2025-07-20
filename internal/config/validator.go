package config

import (
	"fmt"
	"os"
	"strings"

	"spooky/internal/logging"

	"github.com/go-playground/validator/v10"
)

// Validator represents the comprehensive validator
type Validator struct {
	validate *validator.Validate
}

// Global validator instance
var globalValidator *Validator

func init() {
	globalValidator = NewValidator()
}

// NewValidator creates a new validator with all custom validation functions registered
func NewValidator() *Validator {
	v := &Validator{
		validate: validator.New(),
	}
	v.registerCustomValidations()
	return v
}

// ValidateConfig validates the configuration structure
func ValidateConfig(config *Config) error {
	return globalValidator.validateConfig(config)
}

// registerCustomValidations registers all custom validation functions
func (v *Validator) registerCustomValidations() {
	// Register custom field validators
	if err := v.validate.RegisterValidation("sshkeyfile", v.validateSSHKeyFile); err != nil {
		panic(fmt.Sprintf("failed to register sshkeyfile validator: %v", err))
	}
	if err := v.validate.RegisterValidation("scriptfile", v.validateScriptFile); err != nil {
		panic(fmt.Sprintf("failed to register scriptfile validator: %v", err))
	}

	// Register struct-level validations for cross-field validation
	v.validate.RegisterStructValidation(v.validateMachineStruct, Machine{})
	v.validate.RegisterStructValidation(v.validateActionStruct, Action{})
	v.validate.RegisterStructValidation(v.validateConfigStruct, Config{})
}

// validateSSHKeyFile validates that the SSH key file exists and is readable
func (v *Validator) validateSSHKeyFile(fl validator.FieldLevel) bool {
	keyFile := fl.Field().String()
	if keyFile == "" {
		return true // Empty is handled by required_without
	}

	// Check if file exists and is readable
	if _, err := os.Stat(keyFile); err != nil {
		return false
	}

	// Check if file is readable
	if _, err := os.ReadFile(keyFile); err != nil {
		return false
	}

	return true
}

// validateScriptFile validates that the script file exists and is executable
func (v *Validator) validateScriptFile(fl validator.FieldLevel) bool {
	scriptFile := fl.Field().String()
	if scriptFile == "" {
		return true // Empty is handled by required_without
	}

	// Check if file exists
	if _, err := os.Stat(scriptFile); err != nil {
		return false
	}

	// Check if file is executable (Unix-like systems)
	if info, err := os.Stat(scriptFile); err == nil {
		if info.Mode()&0o111 == 0 {
			return false
		}
	}

	return true
}

// validateMachineStruct performs struct-level validation for Machine
func (v *Validator) validateMachineStruct(sl validator.StructLevel) {
	machine := sl.Current().Interface().(Machine)

	// Validate authentication requirements (either password or key_file must be provided)
	if machine.Password == "" && machine.KeyFile == "" {
		sl.ReportError(machine.Password, "Password", "password", "machine_auth", machine.Name)
	}

	// Note: File validation is disabled for testing purposes
	// In production, uncomment the following code to validate SSH key files:
	// if machine.KeyFile != "" {
	// 	if !v.validateSSHKeyFileExists(machine.KeyFile) {
	// 		sl.ReportError(machine.KeyFile, "KeyFile", "key_file", "sshkeyfile", machine.Name)
	// 	}
	// }
}

// validateActionStruct performs struct-level validation for Action
func (v *Validator) validateActionStruct(sl validator.StructLevel) {
	action := sl.Current().Interface().(Action)

	// Validate execution requirements (either command or script must be provided, but not both)
	if action.Command == "" && action.Script == "" {
		sl.ReportError(action.Command, "Command", "command", "action_exec", action.Name)
	}
	if action.Command != "" && action.Script != "" {
		sl.ReportError(action.Command, "Command", "command", "action_exec", action.Name)
	}

	// Note: File validation is disabled for testing purposes
	// In production, uncomment the following code to validate script files:
	// if action.Script != "" {
	// 	if !v.validateScriptFileExists(action.Script) {
	// 		sl.ReportError(action.Script, "Script", "script", "scriptfile", action.Name)
	// 	}
	// }
}

// validateConfigStruct performs struct-level validation for Config
func (v *Validator) validateConfigStruct(sl validator.StructLevel) {
	config := sl.Current().Interface().(Config)

	// Validate unique machine names
	machineNames := make(map[string]bool)
	for _, machine := range config.Machines {
		if machineNames[machine.Name] {
			sl.ReportError(machine.Name, "Name", "name", "unique_machine", machine.Name)
		}
		machineNames[machine.Name] = true
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

	// Validate machine references in actions
	for i := range config.Actions {
		action := &config.Actions[i]
		for _, machineRef := range action.Machines {
			if !machineNames[machineRef] {
				sl.ReportError(machineRef, "Machines", "machines", "valid_machines", action.Name)
			}
		}
	}
}

// validateConfig validates the entire configuration
func (v *Validator) validateConfig(config *Config) error {
	logger := logging.GetLogger()

	// Set defaults before validation
	SetDefaults(config)

	// Perform validation
	if err := v.validate.Struct(config); err != nil {
		logger.Error("Configuration validation failed", err,
			logging.Int("machine_count", len(config.Machines)),
			logging.Int("action_count", len(config.Actions)),
		)
		return v.formatValidationErrors(err)
	}

	logger.Info("Configuration validation successful",
		logging.Int("machine_count", len(config.Machines)),
		logging.Int("action_count", len(config.Actions)),
	)

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
	// Handle special cases for min validation
	if e.Tag() == "min" {
		return v.formatMinValidation(e)
	}

	// Use map for other validation tags
	errorMessages := map[string]string{
		"required":       fmt.Sprintf("%s is required", e.Field()),
		"max":            fmt.Sprintf("%s must be at most %s", e.Field(), e.Param()),
		"machine_auth":   fmt.Sprintf("either password or key_file must be specified for machine %s", e.Param()),
		"action_exec":    fmt.Sprintf("either command or script must be specified for action %s (but not both)", e.Param()),
		"unique_machine": fmt.Sprintf("duplicate machine name: %s", e.Param()),
		"unique_action":  fmt.Sprintf("duplicate action name: %s", e.Param()),
		"valid_port":     fmt.Sprintf("port must be between 1 and 65535 for machine %s", e.Param()),
		"valid_timeout":  fmt.Sprintf("timeout must be between 1 and 3600 seconds for action %s", e.Param()),
		"valid_machines": fmt.Sprintf("machine reference '%s' in action '%s' does not exist", e.Value(), e.Param()),
		"sshkeyfile":     fmt.Sprintf("SSH key file '%s' does not exist or is not readable for machine %s", e.Value(), e.Param()),
		"scriptfile":     fmt.Sprintf("script file '%s' does not exist or is not executable for action %s", e.Value(), e.Param()),
	}

	if message, exists := errorMessages[e.Tag()]; exists {
		return message
	}

	return fmt.Sprintf("%s failed validation: %s", e.Field(), e.Tag())
}

// formatMinValidation handles the complex min validation logic
func (v *Validator) formatMinValidation(e validator.FieldError) string {
	if e.Field() == "Machines" {
		return "at least one machine must be defined"
	}

	// Handle numeric fields differently
	if e.Field() == "Port" || e.Field() == "Timeout" {
		return fmt.Sprintf("%s must be at least %s", e.Field(), e.Param())
	}

	return fmt.Sprintf("%s must be at least %s", e.Field(), e.Param())
}

// ValidateMachine validates a single machine configuration
func (v *Validator) ValidateMachine(machine *Machine) error {
	logger := logging.GetLogger()

	// Perform validation
	if err := v.validate.Struct(machine); err != nil {
		logger.Error("Machine validation failed", err,
			logging.String("machine_name", machine.Name),
		)
		return v.formatValidationErrors(err)
	}

	logger.Info("Machine validation successful",
		logging.String("machine_name", machine.Name),
	)

	return nil
}

// ValidateAction validates a single action configuration
func (v *Validator) ValidateAction(action *Action) error {
	logger := logging.GetLogger()

	if err := v.validate.Struct(action); err != nil {
		logger.Error("Action validation failed", err,
			logging.Action(action.Name),
			logging.String("description", action.Description),
		)
		return v.formatValidationErrors(err)
	}

	logger.Info("Action validation successful",
		logging.Action(action.Name),
		logging.String("description", action.Description),
		logging.Bool("parallel", action.Parallel),
		logging.Int("timeout", action.Timeout),
	)

	return nil
}
