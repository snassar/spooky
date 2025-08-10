package interfaces

import (
	"context"
	"io"
	"time"

	spookytypesconfig "spooky/internal/types/config"
	spookytypesmachines "spooky/internal/types/machines"
)

// MachineManager defines the main interface for machine management
type MachineManager interface {
	// Core machine operations
	LoadMachines(projectPath string) ([]*spookytypesconfig.Machine, error)
	GetMachine(name string) (*spookytypesconfig.Machine, error)
	ListMachines() ([]*spookytypesconfig.Machine, error)
	AddMachine(machine *spookytypesconfig.Machine) error
	RemoveMachine(name string) error

	// Connection operations
	ConnectMachine(machine *spookytypesconfig.Machine) error
	DisconnectMachine(machine *spookytypesconfig.Machine) error
	TestConnection(machine *spookytypesconfig.Machine) error

	// Validation operations
	ValidateMachine(machine *spookytypesconfig.Machine) error
	ValidateMachines(machines []*spookytypesconfig.Machine) error

	// Import/Export operations
	ExportMachines(w io.Writer, format string) error
	ImportMachines(r io.Reader, format string) error

	// Utility operations
	Close() error
}

// ConnectivityManager defines the interface for machine connectivity operations
type ConnectivityManager interface {
	// Connection management
	ConnectMachine(machine *spookytypesconfig.Machine) error
	DisconnectMachine(machine *spookytypesconfig.Machine) error
	TestConnection(machine *spookytypesconfig.Machine) error

	// Connection pool management
	GetConnection(machine *spookytypesconfig.Machine) (Connection, error)
	ReturnConnection(connection Connection) error
	CloseAllConnections() error

	// Configuration
	SetConnectionTimeout(timeout time.Duration) error
	SetMaxConnections(max int) error
	EnableConnectionPooling(enabled bool) error

	// Utility operations
	GetConnectionStats() *spookytypesmachines.ConnectionStats
	Close() error
}

// Connection defines the interface for machine connections
type Connection interface {
	// Connection information
	GetMachine() *spookytypesconfig.Machine
	IsConnected() bool
	GetConnectionTime() time.Time

	// Connection operations
	Connect() error
	Disconnect() error
	Test() error

	// Configuration
	SetTimeout(timeout time.Duration) error
	SetRetryAttempts(attempts int) error
}

// ImportExportManager defines the interface for machine import/export operations
type ImportExportManager interface {
	// Export operations
	ExportToJSON(w io.Writer, machines []*spookytypesconfig.Machine) error
	ExportToHCL(w io.Writer, machines []*spookytypesconfig.Machine) error
	ExportToCSV(w io.Writer, machines []*spookytypesconfig.Machine) error

	// Import operations
	ImportFromJSON(r io.Reader) ([]*spookytypesconfig.Machine, error)
	ImportFromHCL(r io.Reader) ([]*spookytypesconfig.Machine, error)
	ImportFromCSV(r io.Reader) ([]*spookytypesconfig.Machine, error)

	// Configuration
	SetExportFormat(format string) error
	SetImportFormat(format string) error
	SetValidationEnabled(enabled bool) error

	// Utility operations
	ValidateImportData(data []byte, format string) error
	GetSupportedFormats() []string
	Close() error
}

// IndexingManager defines the interface for machine indexing operations
type IndexingManager interface {
	// Indexing operations
	IndexMachines(machines []*spookytypesconfig.Machine) error
	SearchMachines(query string) ([]*spookytypesconfig.Machine, error)
	GetMachineByIndex(index string) (*spookytypesconfig.Machine, error)

	// Index management
	CreateIndex(name string, fields []string) error
	DeleteIndex(name string) error
	ListIndexes() ([]string, error)

	// Configuration
	SetIndexPath(path string) error
	EnableAutoIndexing(enabled bool) error
	SetIndexUpdateInterval(interval time.Duration) error

	// Utility operations
	RebuildIndex() error
	GetIndexStats() *spookytypesmachines.IndexStats
	Close() error
}

// ValidationManager defines the interface for machine validation operations
type ValidationManager interface {
	// Core validation operations
	ValidateMachine(ctx context.Context, machine *spookytypesconfig.Machine) (*spookytypesmachines.ValidationResult, error)
	ValidateMachines(ctx context.Context, machines []*spookytypesconfig.Machine) ([]*spookytypesmachines.ValidationResult, error)
	ValidateMachineData(ctx context.Context, machine *spookytypesconfig.Machine) error

	// Validation configuration
	SetValidationRules(rules *spookytypesmachines.ValidationRules) error
	EnableStrictValidation(strict bool) error
	SetValidationTimeout(timeout time.Duration) error

	// Custom validation
	RegisterCustomValidator(name string, validator ValidationValidator) error
	UnregisterCustomValidator(name string) error
	GetCustomValidators() []string

	// Utility operations
	GetValidationErrors() []spookytypesmachines.ValidationError
	ClearValidationErrors() error
	Close() error
}

// ValidationValidator defines the interface for validation logic
type ValidationValidator interface {
	// Core validation operations
	ValidateMachineData(ctx context.Context, machine *spookytypesconfig.Machine) error
	ValidateMachineStructure(ctx context.Context, machine *spookytypesconfig.Machine) error
	ValidateMachineFields(ctx context.Context, machine *spookytypesconfig.Machine) error

	// Specific validation types
	ValidateMachineName(ctx context.Context, name string) error
	ValidateMachineHost(ctx context.Context, host string) error
	ValidateMachinePort(ctx context.Context, port int) error
	ValidateMachineUser(ctx context.Context, user string) error
	ValidateMachineTags(ctx context.Context, tags map[string]string) error
	ValidateMachineGroups(ctx context.Context, groups []string) error
	ValidateMachineMetadata(ctx context.Context, metadata map[string]string) error

	// Configuration validation
	ValidateValidationConfiguration(ctx context.Context, config *spookytypesmachines.ValidationConfig) error
}

// ValidationBackend defines the interface for validation storage and operations
type ValidationBackend interface {
	// Storage operations
	StoreValidationResult(ctx context.Context, result *spookytypesmachines.ValidationResult) error
	LoadValidationResult(ctx context.Context, machineName string) (*spookytypesmachines.ValidationResult, error)
	DeleteValidationResult(ctx context.Context, machineName string) error

	// Configuration
	StoreValidationConfig(ctx context.Context, config *spookytypesmachines.ValidationConfig) error
	LoadValidationConfig(ctx context.Context) (*spookytypesmachines.ValidationConfig, error)

	// History and tracking
	StoreValidationHistory(ctx context.Context, machineName string, history *spookytypesmachines.ValidationHistory) error
	LoadValidationHistory(ctx context.Context, machineName string) ([]*spookytypesmachines.ValidationHistory, error)
}
