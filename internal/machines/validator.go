// Package machines provides machine validation functionality for the spooky codebase.
package machines

import (
	"context"
	"fmt"
	"net"
	"regexp"
	"strings"
	"time"

	spookyinterfaces "spooky/internal/interfaces"
	spookyschemas "spooky/internal/schemas"
	spookytypes "spooky/internal/types"
	spookytypeslogging "spooky/internal/types/logging"
	spookytypesmachines "spooky/internal/types/machines"
	spookytypesschemas "spooky/internal/types/schemas"
)

// Validator implements the MachineValidator interface
type Validator struct {
	logger                spookytypeslogging.Logger
	schemaDrivenValidator *spookyschemas.SchemaDrivenValidator
	enhancedValidator     *spookyschemas.EnhancedValidator
}

// NewValidator creates a new MachineValidator instance
func NewValidator(logger spookytypeslogging.Logger) spookyinterfaces.MachineValidator {
	// Create schema-driven validator for machine configuration validation
	schemaDrivenConfig := &spookyschemas.SchemaDrivenValidationConfig{
		UseEmbeddedSchemas: true,
		StrictValidation:   true,
		AllowUnknownFields: false,
		DetailedErrors:     true,
	}
	schemaDrivenValidator := spookyschemas.NewSchemaDrivenValidator(logger, schemaDrivenConfig)

	// Create enhanced validator for individual machine validation
	enhancedConfig := &spookyschemas.ValidationConfig{
		Mode: spookyschemas.ValidationModeStrict,
		ErrorHandling: &spookyschemas.ErrorHandlingConfig{
			StopOnFirstError:   false,
			MaxErrors:          100,
			IncludeWarnings:    true,
			IncludeContext:     true,
			IncludeSuggestions: true,
		},
		Evolution: &spookyschemas.EvolutionConfig{
			EnableTracking:  true,
			AllowDeprecated: true,
			WarnDeprecated:  true,
			AllowBreaking:   false,
		},
	}
	enhancedValidator := spookyschemas.NewEnhancedValidator(enhancedConfig)

	return &Validator{
		logger:                logger,
		schemaDrivenValidator: schemaDrivenValidator,
		enhancedValidator:     enhancedValidator,
	}
}

// ValidateMachines validates a collection of machines
func (v *Validator) ValidateMachines(ctx context.Context, machines []spookytypes.Machine) (*spookytypes.ValidationResult, error) {
	v.logger.Debug("Validating machines", map[string]interface{}{
		"count": len(machines),
	})

	var errors []spookytypesschemas.SchemaError
	var warnings []spookytypesschemas.SchemaError

	for i := range machines {
		result, err := v.ValidateMachine(ctx, &machines[i])
		if err != nil {
			return nil, fmt.Errorf("failed to validate machine %d: %w", i, err)
		}

		if !result.Valid {
			errors = append(errors, result.Errors...)
		}

		warnings = append(warnings, result.Warnings...)
	}

	valid := len(errors) == 0

	return &spookytypes.ValidationResult{
		Valid:    valid,
		Errors:   errors,
		Warnings: warnings,
	}, nil
}

// ValidateMachine validates a single machine
func (v *Validator) ValidateMachine(ctx context.Context, machine *spookytypes.Machine) (*spookytypes.ValidationResult, error) {
	v.logger.Debug("Validating machine", map[string]interface{}{
		"hostname": machine.Hostname,
		"host":     machine.Host,
	})

	// Get machine schema for enhanced validation
	machineSchema, err := v.getMachineSchema()
	if err != nil {
		return nil, fmt.Errorf("failed to get machine schema: %w", err)
	}

	// Use enhanced validator for comprehensive machine validation
	result, err := v.enhancedValidator.ValidateWithEnhancedFeatures(ctx, machineSchema, machine)
	if err != nil {
		return nil, fmt.Errorf("failed to validate machine with enhanced validator: %w", err)
	}

	// Add additional custom validation for machine-specific rules
	v.addCustomMachineValidation(machine, result)

	return &spookytypes.ValidationResult{
		Valid:    result.Valid,
		Errors:   result.Errors,
		Warnings: result.Warnings,
	}, nil
}

// getMachineSchema gets the machine schema for validation
func (v *Validator) getMachineSchema() (*spookytypesschemas.Schema, error) {
	// Try to get schema from embedded schemas first
	if schema, err := v.schemaDrivenValidator.GetEmbeddedSchema("machines"); err == nil {
		return schema, nil
	}

	// Fallback: create a basic machine schema
	return &spookytypesschemas.Schema{
		Name:        "machines",
		Type:        "hcl",
		Version:     "1.0",
		Description: "Machine configuration schema",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Content:     "", // Will be loaded from file if needed
		Metadata:    make(map[string]interface{}),
	}, nil
}

// addCustomMachineValidation adds custom validation rules specific to machines
func (v *Validator) addCustomMachineValidation(machine *spookytypes.Machine, result *spookytypesschemas.ValidationResult) {
	// Validate hostname format
	if err := v.validateHostname(machine.Hostname); err != nil {
		v.addSchemaError(result, "invalid_hostname", err.Error(), "error")
	}

	// Validate host format
	if err := v.validateHost(machine.Host); err != nil {
		v.addSchemaError(result, "invalid_host", err.Error(), "error")
	}

	// Validate port range
	if err := v.validatePort(machine.Port); err != nil {
		v.addSchemaError(result, "invalid_port", err.Error(), "error")
	}

	// Validate user
	if err := v.validateUser(machine.User); err != nil {
		v.addSchemaError(result, "invalid_user", err.Error(), "error")
	}

	// Validate authentication
	if err := v.validateAuthentication(machine); err != nil {
		v.addSchemaError(result, "invalid_authentication", err.Error(), "error")
	}

	// Validate SSH configuration
	if err := v.validateSSHConfig(machine); err != nil {
		v.addSchemaError(result, "invalid_ssh_config", err.Error(), "error")
	}
}

// addSchemaError adds a schema error to the validation result
func (v *Validator) addSchemaError(result *spookytypesschemas.ValidationResult, code, message, severity string) {
	schemaError := spookytypesschemas.SchemaError{
		Code:     code,
		Message:  message,
		Severity: severity,
	}
	result.Errors = append(result.Errors, schemaError)
	result.Valid = false
}

// validateHostname validates the machine hostname
func (v *Validator) validateHostname(hostname string) error {
	if hostname == "" {
		return spookytypesmachines.NewMachineValidationError("", "hostname", hostname, "required", "hostname is required")
	}

	// Check hostname format
	hostnameRegex := regexp.MustCompile(`^[a-zA-Z0-9.-]+$`)
	if !hostnameRegex.MatchString(hostname) {
		return spookytypesmachines.NewMachineValidationError(hostname, "hostname", hostname, "format", "hostname must contain only alphanumeric characters, dots, and hyphens")
	}

	return nil
}

// validateHost validates the machine host
func (v *Validator) validateHost(host string) error {
	if host == "" {
		return spookytypesmachines.NewMachineValidationError("", "host", host, "required", "host is required")
	}

	// Try to parse as IP address
	if ip := net.ParseIP(host); ip != nil {
		return nil // Valid IP address
	}

	// Check if it's a valid hostname
	hostnameRegex := regexp.MustCompile(`^[a-zA-Z0-9.-]+$`)
	if hostnameRegex.MatchString(host) {
		return nil // Valid hostname
	}

	return spookytypesmachines.NewMachineValidationError("", "host", host, "format", "host must be a valid IP address or hostname")
}

// validatePort validates the machine port
func (v *Validator) validatePort(port int) error {
	if port < 1 || port > 65535 {
		return spookytypesmachines.NewMachineValidationError("", "port", port, "range", "port must be between 1 and 65535")
	}

	return nil
}

// validateUser validates the machine user
func (v *Validator) validateUser(user string) error {
	if user == "" {
		return spookytypesmachines.NewMachineValidationError("", "user", user, "required", "user is required")
	}

	// Check username format
	userRegex := regexp.MustCompile(`^[a-zA-Z0-9._-]+$`)
	if !userRegex.MatchString(user) {
		return spookytypesmachines.NewMachineValidationError("", "user", user, "format", "username must contain only alphanumeric characters, dots, underscores, and hyphens")
	}

	return nil
}

// validateAuthentication validates the machine authentication
func (v *Validator) validateAuthentication(machine *spookytypes.Machine) error {
	// Check that either password or key_file is provided
	if machine.Password == "" && machine.KeyFile == "" {
		return spookytypesmachines.NewMachineValidationError(machine.Hostname, "authentication", nil, "required", "either password or key_file must be provided")
	}

	// Check that both password and key_file are not provided
	if machine.Password != "" && machine.KeyFile != "" {
		return spookytypesmachines.NewMachineValidationError(machine.Hostname, "authentication", nil, "mutual_exclusive", "password and key_file are mutually exclusive")
	}

	// Validate key_file path if provided
	if machine.KeyFile != "" {
		if strings.HasPrefix(machine.KeyFile, "/") {
			return spookytypesmachines.NewMachineValidationError(machine.Hostname, "key_file", machine.KeyFile, "relative_path", "key_file must be a relative path")
		}
	}

	return nil
}

// validateSSHConfig validates the SSH configuration
func (v *Validator) validateSSHConfig(machine *spookytypes.Machine) error {
	validators := []func(*spookytypes.Machine) error{
		v.validateConnectionTimeout,
		v.validateCommandTimeout,
		v.validateMaxConnections,
		v.validateRetryAttempts,
		v.validateRetryDelay,
	}

	for _, validator := range validators {
		if err := validator(machine); err != nil {
			return err
		}
	}

	return nil
}

func (v *Validator) validateConnectionTimeout(machine *spookytypes.Machine) error {
	if machine.ConnectionTimeout < 1 || machine.ConnectionTimeout > 300 {
		return spookytypesmachines.NewMachineValidationError(machine.Hostname, "connection_timeout", machine.ConnectionTimeout, "range", "connection_timeout must be between 1 and 300 seconds")
	}
	return nil
}

func (v *Validator) validateCommandTimeout(machine *spookytypes.Machine) error {
	if machine.CommandTimeout < 1 || machine.CommandTimeout > 3600 {
		return spookytypesmachines.NewMachineValidationError(machine.Hostname, "command_timeout", machine.CommandTimeout, "range", "command_timeout must be between 1 and 3600 seconds")
	}
	return nil
}

func (v *Validator) validateMaxConnections(machine *spookytypes.Machine) error {
	if machine.MaxConnections < 1 || machine.MaxConnections > 100 {
		return spookytypesmachines.NewMachineValidationError(machine.Hostname, "max_connections", machine.MaxConnections, "range", "max_connections must be between 1 and 100")
	}
	return nil
}

func (v *Validator) validateRetryAttempts(machine *spookytypes.Machine) error {
	if machine.RetryAttempts < 0 || machine.RetryAttempts > 10 {
		return spookytypesmachines.NewMachineValidationError(machine.Hostname, "retry_attempts", machine.RetryAttempts, "range", "retry_attempts must be between 0 and 10")
	}
	return nil
}

func (v *Validator) validateRetryDelay(machine *spookytypes.Machine) error {
	if machine.RetryDelay < 1 || machine.RetryDelay > 60 {
		return spookytypesmachines.NewMachineValidationError(machine.Hostname, "retry_delay", machine.RetryDelay, "range", "retry_delay must be between 1 and 60 seconds")
	}
	return nil
}
