package validation

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	spookyconfigtypes "spooky/internal/types/config"
	spookylogging "spooky/internal/logging"
	spookymachinestypes "spooky/internal/machines/types"
	spookyschemas "spooky/internal/schemas"
)

// Manager implements the ValidationManager interface
type Manager struct {
	config    *spookymachinestypes.ValidationConfig
	validator ValidationValidator
	backend   ValidationBackend
	logger    spookylogging.Logger
}

// NewManager creates a new validation manager with dependency injection
func NewManager(
	config *spookymachinestypes.ValidationConfig,
	validator ValidationValidator,
	backend ValidationBackend,
	logger spookylogging.Logger,
) *Manager {
	return &Manager{
		config:    config,
		validator: validator,
		backend:   backend,
		logger:    logger,
	}
}

// ValidateMachine validates a single machine
func (m *Manager) ValidateMachine(ctx context.Context, machine *spookyconfigtypes.Machine) error {
	if !m.config.Enabled {
		return fmt.Errorf("validation is disabled")
	}

	startTime := time.Now()
	m.logger.Debug("validating machine", spookylogging.String("machine", machine.Name))

	// Validate machine data
	if err := m.validator.ValidateMachineData(ctx, machine); err != nil {
		m.logger.Error("machine validation failed", err,
			spookylogging.String("machine", machine.Name))
		return fmt.Errorf("machine validation failed: %w", err)
	}

	// Validate machine structure
	if err := m.validator.ValidateMachineStructure(ctx, machine); err != nil {
		m.logger.Error("machine structure validation failed", err,
			spookylogging.String("machine", machine.Name))
		return fmt.Errorf("machine structure validation failed: %w", err)
	}

	// Validate machine fields
	if err := m.validator.ValidateMachineFields(ctx, machine); err != nil {
		m.logger.Error("machine fields validation failed", err,
			spookylogging.String("machine", machine.Name))
		return fmt.Errorf("machine fields validation failed: %w", err)
	}

	duration := time.Since(startTime)
	m.logger.Debug("machine validation completed successfully",
		spookylogging.String("machine", machine.Name),
		spookylogging.Int64("duration_ms", int64(duration.Milliseconds())))

	return nil
}

// ValidateMachines validates multiple machines
func (m *Manager) ValidateMachines(ctx context.Context, machines []*spookyconfigtypes.Machine) (*spookymachinestypes.ValidationResult, error) {
	if !m.config.Enabled {
		return nil, fmt.Errorf("validation is disabled")
	}

	startTime := time.Now()
	m.logger.Info("validating machines", spookylogging.Int("count", len(machines)))

	result := &spookymachinestypes.ValidationResult{
		MachineName:    "batch",
		Valid:          true,
		Timestamp:      startTime,
		ValidationType: "batch",
		Errors:         make([]spookymachinestypes.ValidationError, 0),
		Warnings:       make([]spookymachinestypes.ValidationError, 0),
		Info:           make([]spookymachinestypes.ValidationError, 0),
	}

	validCount := 0
	invalidCount := 0

	for _, machine := range machines {
		if err := m.ValidateMachine(ctx, machine); err != nil {
			invalidCount++
			result.Valid = false
			result.Errors = append(result.Errors, spookymachinestypes.ValidationError{
				MachineName: machine.Name,
				Field:       "machine",
				Message:     err.Error(),
				Severity:    "error",
			})

			if m.config.FailFast {
				break
			}
		} else {
			validCount++
		}
	}

	result.Duration = time.Since(startTime)
	result.Metadata = map[string]string{
		"total_machines":   fmt.Sprintf("%d", len(machines)),
		"valid_machines":   fmt.Sprintf("%d", validCount),
		"invalid_machines": fmt.Sprintf("%d", invalidCount),
	}

	m.logger.Info("batch validation completed",
		spookylogging.Int("total", len(machines)),
		spookylogging.Int("valid", validCount),
		spookylogging.Int("invalid", invalidCount),
		spookylogging.Int64("duration_ms", int64(result.Duration.Milliseconds())))

	return result, nil
}

// ValidateMachineInventory validates a machine inventory
func (m *Manager) ValidateMachineInventory(ctx context.Context, inventory *spookyconfigtypes.InventoryConfig) (*spookymachinestypes.ValidationResult, error) {
	if !m.config.Enabled {
		return nil, fmt.Errorf("validation is disabled")
	}

	m.logger.Info("validating machine inventory", spookylogging.Int("count", len(inventory.Machines)))

	// Convert inventory machines to pointers for validation
	machines := make([]*spookyconfigtypes.Machine, len(inventory.Machines))
	for i := range inventory.Machines {
		machines[i] = &inventory.Machines[i]
	}

	result, err := m.ValidateMachines(ctx, machines)
	if err != nil {
		return nil, err
	}

	result.MachineName = "inventory"
	result.ValidationType = "inventory"

	m.logger.Info("inventory validation completed",
		spookylogging.Int("total", len(inventory.Machines)),
		spookylogging.Bool("valid", result.Valid),
		spookylogging.Int64("duration_ms", int64(result.Duration.Milliseconds())))

	return result, nil
}

// ValidateMachineConfiguration validates machine configuration
func (m *Manager) ValidateMachineConfiguration(ctx context.Context, machine *spookyconfigtypes.Machine) error {
	if !m.config.ValidateConfiguration {
		return nil
	}

	m.logger.Debug("validating machine configuration", spookylogging.String("machine", machine.Name))

	// Validate individual configuration fields
	if err := m.validator.ValidateMachineName(ctx, machine.Name); err != nil {
		return fmt.Errorf("invalid machine name: %w", err)
	}

	if err := m.validator.ValidateMachineHost(ctx, machine.Host); err != nil {
		return fmt.Errorf("invalid machine host: %w", err)
	}

	if err := m.validator.ValidateMachinePort(ctx, machine.Port); err != nil {
		return fmt.Errorf("invalid machine port: %w", err)
	}

	if err := m.validator.ValidateMachineUser(ctx, machine.User); err != nil {
		return fmt.Errorf("invalid machine user: %w", err)
	}

	if err := m.validator.ValidateMachineTags(ctx, machine.Tags); err != nil {
		return fmt.Errorf("invalid machine tags: %w", err)
	}

	if err := m.validator.ValidateMachineGroups(ctx, machine.Groups); err != nil {
		return fmt.Errorf("invalid machine groups: %w", err)
	}

	if err := m.validator.ValidateMachineMetadata(ctx, machine.Metadata); err != nil {
		return fmt.Errorf("invalid machine metadata: %w", err)
	}

	m.logger.Debug("machine configuration validation completed", spookylogging.String("machine", machine.Name))
	return nil
}

// ValidateMachineConnectivity validates machine connectivity
func (m *Manager) ValidateMachineConnectivity(ctx context.Context, machine *spookyconfigtypes.Machine) error {
	if !m.config.ValidateConnectivity {
		return nil
	}

	m.logger.Debug("validating machine connectivity", spookylogging.String("machine", machine.Name))

	// Validate host configuration
	if machine.Host == "" {
		return fmt.Errorf("machine host is required")
	}

	// Validate port configuration
	if machine.Port <= 0 || machine.Port > 65535 {
		return fmt.Errorf("machine port must be between 1 and 65535")
	}

	// Validate hostname format
	if !isValidHostname(machine.Host) && !isValidIPAddress(machine.Host) {
		return fmt.Errorf("invalid host format: %s", machine.Host)
	}

	// TODO: Implement actual connectivity test using SSH client
	// This would involve creating a temporary SSH connection to verify connectivity
	// For now, we'll just validate the configuration

	m.logger.Debug("machine connectivity validation completed", spookylogging.String("machine", machine.Name))
	return nil
}

// ValidateMachineSecurity validates machine security settings
func (m *Manager) ValidateMachineSecurity(ctx context.Context, machine *spookyconfigtypes.Machine) error {
	if !m.config.ValidateSecurity {
		return nil
	}

	m.logger.Debug("validating machine security", spookylogging.String("machine", machine.Name))

	// Validate user configuration
	if machine.User == "" {
		return fmt.Errorf("machine user is required for security")
	}

	// Validate authentication method
	if machine.Password == "" && machine.KeyFile == "" {
		return fmt.Errorf("either password or key file must be provided for authentication")
	}

	// Validate key file if provided
	if machine.KeyFile != "" {
		if err := validateKeyFile(machine.KeyFile); err != nil {
			return fmt.Errorf("invalid key file: %w", err)
		}
	}

	// Validate user permissions
	if machine.User == "root" {
		return fmt.Errorf("root user is not allowed for security reasons")
	}

	m.logger.Debug("machine security validation completed", spookylogging.String("machine", machine.Name))
	return nil
}

// Helper functions for validation

// isValidHostname validates if a string is a valid hostname
func isValidHostname(hostname string) bool {
	if len(hostname) == 0 || len(hostname) > 253 {
		return false
	}

	// Basic hostname validation - check for valid characters
	for _, char := range hostname {
		if !((char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9') || char == '-' || char == '.') {
			return false
		}
	}

	return true
}

// isValidIPAddress validates if a string is a valid IP address
func isValidIPAddress(ip string) bool {
	// Basic IP validation - check for IPv4 format
	parts := strings.Split(ip, ".")
	if len(parts) != 4 {
		return false
	}

	for _, part := range parts {
		if len(part) == 0 || len(part) > 3 {
			return false
		}
		for _, char := range part {
			if char < '0' || char > '9' {
				return false
			}
		}
		num := 0
		for _, char := range part {
			num = num*10 + int(char-'0')
		}
		if num < 0 || num > 255 {
			return false
		}
	}

	return true
}

// validateKeyFile validates if a key file exists and is readable
func validateKeyFile(keyFile string) error {
	// Check if file exists
	if _, err := os.Stat(keyFile); os.IsNotExist(err) {
		return fmt.Errorf("key file does not exist: %s", keyFile)
	}

	// Check if file is readable
	if _, err := os.ReadFile(keyFile); err != nil {
		return fmt.Errorf("key file is not readable: %w", err)
	}

	return nil
}

// ConfigureValidation configures the validation system
func (m *Manager) ConfigureValidation(ctx context.Context, config *spookymachinestypes.ValidationConfig) error {
	m.logger.Info("configuring validation system")

	// Validate configuration
	if err := m.validator.ValidateValidationConfiguration(ctx, config); err != nil {
		return fmt.Errorf("validation configuration validation failed: %w", err)
	}

	// Store configuration
	if err := m.backend.StoreValidationConfig(ctx, config); err != nil {
		return fmt.Errorf("failed to store validation configuration: %w", err)
	}

	m.config = config
	m.logger.Info("validation system configured successfully")
	return nil
}

// GetValidationConfig retrieves the current validation configuration
func (m *Manager) GetValidationConfig(ctx context.Context) (*spookymachinestypes.ValidationConfig, error) {
	return m.backend.LoadValidationConfig(ctx)
}

// ValidateValidationHealth validates the overall health of the validation system
func (m *Manager) ValidateValidationHealth(ctx context.Context) error {
	if !m.config.Enabled {
		return fmt.Errorf("validation system is disabled")
	}

	m.logger.Debug("validating validation system health")

	// This is a placeholder implementation
	// In a real implementation, this would check various health indicators
	// For now, we'll just return success

	m.logger.Debug("validation system health validated successfully")
	return nil
}

// GetValidationStatus retrieves the status of the validation system
func (m *Manager) GetValidationStatus(ctx context.Context) (*spookymachinestypes.ValidationStatus, error) {
	if !m.config.Enabled {
		return &spookymachinestypes.ValidationStatus{Enabled: false}, nil
	}

	// This is a placeholder implementation
	// In a real implementation, this would gather actual status information
	status := &spookymachinestypes.ValidationStatus{
		Enabled:               m.config.Enabled,
		TotalMachines:         0,
		ValidMachines:         0,
		InvalidMachines:       0,
		PendingValidations:    0,
		AverageValidationTime: 0,
		ValidationErrors:      0,
		ConfigurationValid:    m.config.ValidateConfiguration,
		ConnectivityValid:     m.config.ValidateConnectivity,
		SecurityValid:         m.config.ValidateSecurity,
		StructureValid:        m.config.ValidateStructure,
		RecentResults:         make([]spookymachinestypes.ValidationResult, 0),
	}

	return status, nil
}

// validateMachinesComposed validates complex machines using schema system
func (m *Manager) validateMachinesComposed(content []byte) error {
	validator := spookyschemas.NewSchemaValidator()
	if err := validator.LoadSchema(spookyschemas.SchemaTypeMachinesComposed); err != nil {
		return fmt.Errorf("failed to load composed machines schema: %w", err)
	}

	// Parse content to interface{} for validation
	var data interface{}
	if err := json.Unmarshal(content, &data); err != nil {
		return fmt.Errorf("failed to parse machines for validation: %w", err)
	}

	if err := validator.ValidateData(data, "machines-composed"); err != nil {
		return fmt.Errorf("composed machines validation failed: %w", err)
	}
	return nil
}
