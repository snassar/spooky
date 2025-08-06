package validation

import (
	"context"

	spookyconfigtypes "spooky/internal/config/types"
	spookymachinestypes "spooky/internal/machines/types"
)

// ValidationManager defines the interface for machine validation operations
type ValidationManager interface {
	// Core validation operations
	ValidateMachine(ctx context.Context, machine *spookyconfigtypes.Machine) error
	ValidateMachines(ctx context.Context, machines []*spookyconfigtypes.Machine) (*spookymachinestypes.ValidationResult, error)
	ValidateMachineInventory(ctx context.Context, inventory *spookyconfigtypes.InventoryConfig) (*spookymachinestypes.ValidationResult, error)

	// Specific validation types
	ValidateMachineConfiguration(ctx context.Context, machine *spookyconfigtypes.Machine) error
	ValidateMachineConnectivity(ctx context.Context, machine *spookyconfigtypes.Machine) error
	ValidateMachineSecurity(ctx context.Context, machine *spookyconfigtypes.Machine) error

	// Configuration
	ConfigureValidation(ctx context.Context, config *spookymachinestypes.ValidationConfig) error
	GetValidationConfig(ctx context.Context) (*spookymachinestypes.ValidationConfig, error)

	// Status and health
	ValidateValidationHealth(ctx context.Context) error
	GetValidationStatus(ctx context.Context) (*spookymachinestypes.ValidationStatus, error)
}

// ValidationValidator defines the interface for validation logic
type ValidationValidator interface {
	// Core validation operations
	ValidateMachineData(ctx context.Context, machine *spookyconfigtypes.Machine) error
	ValidateMachineStructure(ctx context.Context, machine *spookyconfigtypes.Machine) error
	ValidateMachineFields(ctx context.Context, machine *spookyconfigtypes.Machine) error

	// Specific validation types
	ValidateMachineName(ctx context.Context, name string) error
	ValidateMachineHost(ctx context.Context, host string) error
	ValidateMachinePort(ctx context.Context, port int) error
	ValidateMachineUser(ctx context.Context, user string) error
	ValidateMachineTags(ctx context.Context, tags map[string]string) error
	ValidateMachineGroups(ctx context.Context, groups []string) error
	ValidateMachineMetadata(ctx context.Context, metadata map[string]string) error

	// Configuration validation
	ValidateValidationConfiguration(ctx context.Context, config *spookymachinestypes.ValidationConfig) error
}

// ValidationBackend defines the interface for validation storage and operations
type ValidationBackend interface {
	// Storage operations
	StoreValidationResult(ctx context.Context, result *spookymachinestypes.ValidationResult) error
	LoadValidationResult(ctx context.Context, machineName string) (*spookymachinestypes.ValidationResult, error)
	DeleteValidationResult(ctx context.Context, machineName string) error

	// Configuration
	StoreValidationConfig(ctx context.Context, config *spookymachinestypes.ValidationConfig) error
	LoadValidationConfig(ctx context.Context) (*spookymachinestypes.ValidationConfig, error)

	// History and tracking
	StoreValidationHistory(ctx context.Context, machineName string, history *spookymachinestypes.ValidationHistory) error
	LoadValidationHistory(ctx context.Context, machineName string) ([]*spookymachinestypes.ValidationHistory, error)
}
