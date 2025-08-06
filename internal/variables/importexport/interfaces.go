package importexport

import (
	"context"
	"io"
	"spooky/internal/variables/types"
)

// ImportExportManager defines the interface for import/export operations
type ImportExportManager interface {
	// Export operations
	ExportToHCL(ctx context.Context, variables []*types.Variable, path string) error
	ExportToJSON(ctx context.Context, variables []*types.Variable, path string) error
	ExportToWriter(ctx context.Context, variables []*types.Variable, writer io.Writer, format types.ExportFormat) error

	// Import operations
	ImportFromHCL(ctx context.Context, path string) ([]*types.Variable, error)
	ImportFromJSON(ctx context.Context, path string) ([]*types.Variable, error)
	ImportFromReader(ctx context.Context, reader io.Reader, format types.ImportFormat) ([]*types.Variable, error)

	// Configuration
	SetExportOptions(options *types.ExportOptions) error
	SetImportOptions(options *types.ImportOptions) error
	SetDefaultFormat(format types.ExportFormat) error

	// Utility operations
	GetSupportedExportFormats() []types.ExportFormat
	GetSupportedImportFormats() []types.ImportFormat
	ValidateExportPath(path string) error
	ValidateImportPath(path string) error
	Close() error
}

// VariableExporter defines the interface for specific export formats
type VariableExporter interface {
	Export(ctx context.Context, variables []*types.Variable, writer io.Writer) error
	GetName() string
	GetExtension() string
}

// VariableImporter defines the interface for specific import formats
type VariableImporter interface {
	Import(ctx context.Context, reader io.Reader) ([]*types.Variable, error)
	GetName() string
	GetExtension() string
}
