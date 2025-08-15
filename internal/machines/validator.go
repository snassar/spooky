// Package machines provides machine validation functionality for the spooky codebase.
package machines

import (
	"context"
	"fmt"
	"net"
	"regexp"
	"strings"

	spookyinterfaces "spooky/internal/interfaces"
	spookytypes "spooky/internal/types"
	spookytypeslogging "spooky/internal/types/logging"
	spookytypesmachines "spooky/internal/types/machines"
	spookytypesschemas "spooky/internal/types/schemas"
)

// Validator implements the MachineValidator interface
type Validator struct {
	logger spookytypeslogging.Logger
}

// NewValidator creates a new MachineValidator instance
func NewValidator(logger spookytypeslogging.Logger) spookyinterfaces.MachineValidator {
	return &Validator{
		logger: logger,
	}
}

// ValidateMachines validates a collection of machines
func (v *Validator) ValidateMachines(ctx context.Context, machines []spookytypes.Machine) (*spookytypes.ValidationResult, error) {
	v.logger.Debug("Validating machines", map[string]interface{}{
		"count": len(machines),
	})

	var errors []spookytypesschemas.SchemaError
	var warnings []spookytypesschemas.SchemaError

	for i, machine := range machines {
		result, err := v.ValidateMachine(ctx, machine)
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
func (v *Validator) ValidateMachine(_ context.Context, machine spookytypes.Machine) (*spookytypes.ValidationResult, error) {
	v.logger.Debug("Validating machine", map[string]interface{}{
		"hostname": machine.Hostname,
		"host":     machine.Host,
	})

	var errors []spookytypesschemas.SchemaError
	var warnings []spookytypesschemas.SchemaError

	// Validate hostname
	if err := v.validateHostname(machine.Hostname); err != nil {
		errors = append(errors, v.convertErrorToSchemaError(err, "hostname"))
	}

	// Validate host
	if err := v.validateHost(machine.Host); err != nil {
		errors = append(errors, v.convertErrorToSchemaError(err, "host"))
	}

	// Validate port
	if err := v.validatePort(machine.Port); err != nil {
		errors = append(errors, v.convertErrorToSchemaError(err, "port"))
	}

	// Validate user
	if err := v.validateUser(machine.User); err != nil {
		errors = append(errors, v.convertErrorToSchemaError(err, "user"))
	}

	// Validate authentication
	if err := v.validateAuthentication(&machine); err != nil {
		errors = append(errors, v.convertErrorToSchemaError(err, "authentication"))
	}

	// Validate SSH configuration
	if err := v.validateSSHConfig(&machine); err != nil {
		errors = append(errors, v.convertErrorToSchemaError(err, "ssh_config"))
	}

	// Validate resource specifications
	if err := v.validateResources(machine.Resources); err != nil {
		errors = append(errors, v.convertErrorToSchemaError(err, "resources"))
	}

	valid := len(errors) == 0

	return &spookytypes.ValidationResult{
		Valid:    valid,
		Errors:   errors,
		Warnings: warnings,
	}, nil
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
	// Validate connection timeout
	if machine.ConnectionTimeout < 1 || machine.ConnectionTimeout > 300 {
		return spookytypesmachines.NewMachineValidationError(machine.Hostname, "connection_timeout", machine.ConnectionTimeout, "range", "connection_timeout must be between 1 and 300 seconds")
	}

	// Validate command timeout
	if machine.CommandTimeout < 1 || machine.CommandTimeout > 3600 {
		return spookytypesmachines.NewMachineValidationError(machine.Hostname, "command_timeout", machine.CommandTimeout, "range", "command_timeout must be between 1 and 3600 seconds")
	}

	// Validate max connections
	if machine.MaxConnections < 1 || machine.MaxConnections > 100 {
		return spookytypesmachines.NewMachineValidationError(machine.Hostname, "max_connections", machine.MaxConnections, "range", "max_connections must be between 1 and 100")
	}

	// Validate retry attempts
	if machine.RetryAttempts < 0 || machine.RetryAttempts > 10 {
		return spookytypesmachines.NewMachineValidationError(machine.Hostname, "retry_attempts", machine.RetryAttempts, "range", "retry_attempts must be between 0 and 10")
	}

	// Validate retry delay
	if machine.RetryDelay < 1 || machine.RetryDelay > 60 {
		return spookytypesmachines.NewMachineValidationError(machine.Hostname, "retry_delay", machine.RetryDelay, "range", "retry_delay must be between 1 and 60 seconds")
	}

	return nil
}

// validateResources validates the resource specifications
func (v *Validator) validateResources(resources *spookytypesmachines.MachineResources) error {
	if resources == nil {
		return nil // Resources are optional
	}

	// Validate CPU cores
	if resources.CPUCores < 1 || resources.CPUCores > 1024 {
		return spookytypesmachines.NewMachineValidationError("", "cpu_cores", resources.CPUCores, "range", "cpu_cores must be between 1 and 1024")
	}

	// Validate memory
	if resources.MemoryGB < 1 || resources.MemoryGB > 32768 {
		return spookytypesmachines.NewMachineValidationError("", "memory_gb", resources.MemoryGB, "range", "memory_gb must be between 1 and 32768")
	}

	// Validate disk
	if resources.DiskGB < 1 || resources.DiskGB > 1048576 {
		return spookytypesmachines.NewMachineValidationError("", "disk_gb", resources.DiskGB, "range", "disk_gb must be between 1 and 1048576")
	}

	// Validate network speed
	if resources.NetworkSpeed != "" {
		networkSpeedRegex := regexp.MustCompile(`^[0-9]+(Gbps|Mbps)$`)
		if !networkSpeedRegex.MatchString(resources.NetworkSpeed) {
			return spookytypesmachines.NewMachineValidationError("", "network_speed", resources.NetworkSpeed, "format", "network_speed must be in format '10Gbps' or '1Mbps'")
		}
	}

	return nil
}

// convertErrorToSchemaError converts a machine error to a schema error
func (v *Validator) convertErrorToSchemaError(err error, field string) spookytypesschemas.SchemaError {
	if machineErr, ok := err.(*spookytypesmachines.MachineValidationError); ok {
		return spookytypesschemas.SchemaError{
			Code:        "machine_validation_error",
			Message:     machineErr.Message,
			Recoverable: machineErr.Recoverable,
			SchemaName:  "machine",
			SchemaType:  "validation",
			FieldPath:   field,
			Value:       machineErr.Value,
			Severity:    "error",
		}
	}

	return spookytypesschemas.SchemaError{
		Code:        "machine_validation_error",
		Message:     err.Error(),
		Recoverable: true,
		SchemaName:  "machine",
		SchemaType:  "validation",
		FieldPath:   field,
		Severity:    "error",
	}
}
