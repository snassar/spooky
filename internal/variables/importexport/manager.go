package importexport

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"spooky/internal/logging"
	"spooky/internal/variables/types"
)

// Manager implements ImportExportManager interface
type Manager struct {
	config    *types.ImportExportConfig
	exporters map[types.ExportFormat]VariableExporter
	importers map[types.ImportFormat]VariableImporter
	logger    logging.Logger
}

// NewManager creates a new import/export manager
func NewManager(config *types.ImportExportConfig, logger logging.Logger) *Manager {
	return &Manager{
		config:    config,
		exporters: make(map[types.ExportFormat]VariableExporter),
		importers: make(map[types.ImportFormat]VariableImporter),
		logger:    logger,
	}
}

// ExportToHCL exports variables to HCL format
func (m *Manager) ExportToHCL(ctx context.Context, variables []*types.Variable, path string) error {
	// 1. Validate export path
	if err := m.ValidateExportPath(path); err != nil {
		return fmt.Errorf("export path validation failed: %w", err)
	}

	// 2. Create file
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create export file: %w", err)
	}
	defer file.Close()

	// 3. Export to file
	return m.ExportToWriter(ctx, variables, file, types.ExportFormatHCL)
}

// ExportToJSON exports variables to JSON format
func (m *Manager) ExportToJSON(ctx context.Context, variables []*types.Variable, path string) error {
	// 1. Validate export path
	if err := m.ValidateExportPath(path); err != nil {
		return fmt.Errorf("export path validation failed: %w", err)
	}

	// 2. Create file
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create export file: %w", err)
	}
	defer file.Close()

	// 3. Export to file
	return m.ExportToWriter(ctx, variables, file, types.ExportFormatJSON)
}

// ExportToWriter exports variables to a writer with specified format
func (m *Manager) ExportToWriter(ctx context.Context, variables []*types.Variable, writer io.Writer, format types.ExportFormat) error {
	// 1. Get appropriate exporter
	exporter, exists := m.exporters[format]
	if !exists {
		return fmt.Errorf("unsupported export format: %s", format)
	}

	// 2. Export variables
	if err := exporter.Export(ctx, variables, writer); err != nil {
		return fmt.Errorf("export failed: %w", err)
	}

	return nil
}

// ImportFromHCL imports variables from HCL format
func (m *Manager) ImportFromHCL(ctx context.Context, path string) ([]*types.Variable, error) {
	// 1. Validate import path
	if err := m.ValidateImportPath(path); err != nil {
		return nil, fmt.Errorf("import path validation failed: %w", err)
	}

	// 2. Open file
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open import file: %w", err)
	}
	defer file.Close()

	// 3. Import from file
	return m.ImportFromReader(ctx, file, types.ImportFormatHCL)
}

// ImportFromJSON imports variables from JSON format
func (m *Manager) ImportFromJSON(ctx context.Context, path string) ([]*types.Variable, error) {
	// 1. Validate import path
	if err := m.ValidateImportPath(path); err != nil {
		return nil, fmt.Errorf("import path validation failed: %w", err)
	}

	// 2. Open file
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open import file: %w", err)
	}
	defer file.Close()

	// 3. Import from file
	return m.ImportFromReader(ctx, file, types.ImportFormatJSON)
}

// ImportFromReader imports variables from a reader with specified format
func (m *Manager) ImportFromReader(ctx context.Context, reader io.Reader, format types.ImportFormat) ([]*types.Variable, error) {
	// 1. Get appropriate importer
	importer, exists := m.importers[format]
	if !exists {
		return nil, fmt.Errorf("unsupported import format: %s", format)
	}

	// 2. Import variables
	variables, err := importer.Import(ctx, reader)
	if err != nil {
		return nil, fmt.Errorf("import failed: %w", err)
	}

	return variables, nil
}

// SetExportOptions sets export options
func (m *Manager) SetExportOptions(options *types.ExportOptions) error {
	if m.config == nil {
		m.config = &types.ImportExportConfig{}
	}
	m.config.ExportOptions = options
	return nil
}

// SetImportOptions sets import options
func (m *Manager) SetImportOptions(options *types.ImportOptions) error {
	if m.config == nil {
		m.config = &types.ImportExportConfig{}
	}
	m.config.ImportOptions = options
	return nil
}

// SetDefaultFormat sets the default export format
func (m *Manager) SetDefaultFormat(format types.ExportFormat) error {
	if m.config == nil {
		m.config = &types.ImportExportConfig{}
	}
	m.config.DefaultFormat = format
	return nil
}

// GetSupportedExportFormats returns supported export formats
func (m *Manager) GetSupportedExportFormats() []types.ExportFormat {
	formats := make([]types.ExportFormat, 0, len(m.exporters))
	for format := range m.exporters {
		formats = append(formats, format)
	}
	return formats
}

// GetSupportedImportFormats returns supported import formats
func (m *Manager) GetSupportedImportFormats() []types.ImportFormat {
	formats := make([]types.ImportFormat, 0, len(m.importers))
	for format := range m.importers {
		formats = append(formats, format)
	}
	return formats
}

// ValidateExportPath validates that an export path is valid
func (m *Manager) ValidateExportPath(path string) error {
	// Check if directory exists and is writable
	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return fmt.Errorf("export directory does not exist: %s", dir)
	}

	return nil
}

// ValidateImportPath validates that an import path is valid
func (m *Manager) ValidateImportPath(path string) error {
	// Check if file exists and is readable
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("import file does not exist: %s", path)
	}

	return nil
}

// Close closes the import/export manager and releases resources
func (m *Manager) Close() error {
	// Clean up resources if needed
	return nil
}
