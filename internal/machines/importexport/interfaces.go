package importexport

import (
	"context"
	"io"

	spookyconfigtypes "spooky/internal/config/types"
	spookymachinestypes "spooky/internal/machines/types"
)

// ImportExportManager defines the interface for machine import/export operations
type ImportExportManager interface {
	// Export operations
	ExportMachines(ctx context.Context, machines []*spookyconfigtypes.Machine, format spookymachinestypes.ExportFormat, output io.Writer) error
	ExportMachinesToFile(ctx context.Context, machines []*spookyconfigtypes.Machine, format spookymachinestypes.ExportFormat, filePath string) error
	ExportMachinesFiltered(ctx context.Context, filter *spookymachinestypes.MachineFilter, format spookymachinestypes.ExportFormat, output io.Writer) (*spookymachinestypes.ExportResult, error)

	// Import operations
	ImportMachines(ctx context.Context, input io.Reader, format spookymachinestypes.ExportFormat) ([]*spookyconfigtypes.Machine, error)
	ImportMachinesFromFile(ctx context.Context, filePath string, format spookymachinestypes.ExportFormat) ([]*spookyconfigtypes.Machine, error)
	ImportMachinesWithValidation(ctx context.Context, input io.Reader, format spookymachinestypes.ExportFormat, validate bool) ([]*spookyconfigtypes.Machine, error)

	// Configuration
	ConfigureExport(ctx context.Context, config *spookymachinestypes.ExportConfig) error
	GetExportConfig(ctx context.Context) (*spookymachinestypes.ExportConfig, error)

	// Status and health
	ValidateExportFormat(ctx context.Context, format spookymachinestypes.ExportFormat) error
	GetSupportedFormats(ctx context.Context) []spookymachinestypes.ExportFormat
}

// ImportExportValidator defines the interface for import/export validation
type ImportExportValidator interface {
	// Validation operations
	ValidateExportData(ctx context.Context, machines []*spookyconfigtypes.Machine) error
	ValidateImportData(ctx context.Context, data []byte, format spookymachinestypes.ExportFormat) error
	ValidateMachineData(ctx context.Context, machine *spookyconfigtypes.Machine) error

	// Format validation
	ValidateExportFormat(ctx context.Context, format spookymachinestypes.ExportFormat) error
	ValidateImportFormat(ctx context.Context, format spookymachinestypes.ExportFormat) error

	// Data integrity checks
	ValidateDataIntegrity(ctx context.Context, data []byte, format spookymachinestypes.ExportFormat) error
	ValidateMachineIntegrity(ctx context.Context, machine *spookyconfigtypes.Machine) error
}

// ImportExportBackend defines the interface for import/export storage and operations
type ImportExportBackend interface {
	// File operations
	WriteExportFile(ctx context.Context, filePath string, data []byte) error
	ReadImportFile(ctx context.Context, filePath string) ([]byte, error)
	ValidateFilePath(ctx context.Context, filePath string) error

	// Data operations
	SerializeMachines(ctx context.Context, machines []*spookyconfigtypes.Machine, format spookymachinestypes.ExportFormat) ([]byte, error)
	DeserializeMachines(ctx context.Context, data []byte, format spookymachinestypes.ExportFormat) ([]*spookyconfigtypes.Machine, error)

	// Configuration
	StoreExportConfig(ctx context.Context, config *spookymachinestypes.ExportConfig) error
	LoadExportConfig(ctx context.Context) (*spookymachinestypes.ExportConfig, error)
}
