package importexport

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	spookyconfigtypes "spooky/internal/types/config"
	spookylogging "spooky/internal/logging"
	spookymachinestypes "spooky/internal/machines/types"
	spookyschemas "spooky/internal/schemas"
)

// Manager implements the ImportExportManager interface
type Manager struct {
	config    *spookymachinestypes.ExportConfig
	validator ImportExportValidator
	backend   ImportExportBackend
	logger    spookylogging.Logger
}

// NewManager creates a new import/export manager with dependency injection
func NewManager(
	config *spookymachinestypes.ExportConfig,
	validator ImportExportValidator,
	backend ImportExportBackend,
	logger spookylogging.Logger,
) *Manager {
	return &Manager{
		config:    config,
		validator: validator,
		backend:   backend,
		logger:    logger,
	}
}

// ExportMachines exports machines to the specified output writer
func (m *Manager) ExportMachines(ctx context.Context, machines []*spookyconfigtypes.Machine, format spookymachinestypes.ExportFormat, output io.Writer) error {
	m.logger.Debug("exporting machines",
		spookylogging.Int("count", len(machines)),
		spookylogging.String("format", string(format)))

	// Validate export format
	if err := m.validator.ValidateExportFormat(ctx, format); err != nil {
		return fmt.Errorf("invalid export format: %w", err)
	}

	// Validate export data if enabled
	if m.config.ValidateOnExport {
		if err := m.validator.ValidateExportData(ctx, machines); err != nil {
			return fmt.Errorf("export data validation failed: %w", err)
		}
	}

	// Serialize machines
	data, err := m.backend.SerializeMachines(ctx, machines, format)
	if err != nil {
		return fmt.Errorf("failed to serialize machines: %w", err)
	}

	// Write to output
	if _, err := output.Write(data); err != nil {
		return fmt.Errorf("failed to write export data: %w", err)
	}

	m.logger.Debug("machines exported successfully",
		spookylogging.Int("count", len(machines)),
		spookylogging.String("format", string(format)),
		spookylogging.Int("bytes", len(data)))

	return nil
}

// ExportMachinesToFile exports machines to a file
func (m *Manager) ExportMachinesToFile(ctx context.Context, machines []*spookyconfigtypes.Machine, format spookymachinestypes.ExportFormat, filePath string) error {
	m.logger.Info("exporting machines to file",
		spookylogging.Int("count", len(machines)),
		spookylogging.String("format", string(format)),
		spookylogging.String("file", filePath))

	// Validate file path
	if err := m.backend.ValidateFilePath(ctx, filePath); err != nil {
		return fmt.Errorf("invalid file path: %w", err)
	}

	// Serialize machines
	data, err := m.backend.SerializeMachines(ctx, machines, format)
	if err != nil {
		return fmt.Errorf("failed to serialize machines: %w", err)
	}

	// Write to file
	if err := m.backend.WriteExportFile(ctx, filePath, data); err != nil {
		return fmt.Errorf("failed to write export file: %w", err)
	}

	m.logger.Info("machines exported to file successfully",
		spookylogging.Int("count", len(machines)),
		spookylogging.String("format", string(format)),
		spookylogging.String("file", filePath))

	return nil
}

// ExportMachinesFiltered exports machines with filtering
func (m *Manager) ExportMachinesFiltered(ctx context.Context, filter *spookymachinestypes.MachineFilter, format spookymachinestypes.ExportFormat, output io.Writer) (*spookymachinestypes.ExportResult, error) {
	m.logger.Debug("exporting machines with filter",
		spookylogging.String("format", string(format)))

	// This is a placeholder implementation
	// In a real implementation, this would filter machines based on the filter criteria
	// For now, we'll return an empty result
	result := &spookymachinestypes.ExportResult{
		Machines: []spookyconfigtypes.Machine{},
		Stats: spookymachinestypes.ExportStats{
			TotalMachines:    0,
			ExportedMachines: 0,
			FilteredMachines: 0,
			ExportTime:       0,
		},
		Options: spookymachinestypes.ExportOptions{
			Format: format,
			Filter: filter,
		},
	}

	m.logger.Debug("filtered export completed",
		spookylogging.String("format", string(format)),
		spookylogging.Int("exported", result.Stats.ExportedMachines))

	return result, nil
}

// ImportMachines imports machines from the specified input reader
func (m *Manager) ImportMachines(ctx context.Context, input io.Reader, format spookymachinestypes.ExportFormat) ([]*spookyconfigtypes.Machine, error) {
	m.logger.Debug("importing machines",
		spookylogging.String("format", string(format)))

	// Validate import format
	if err := m.validator.ValidateImportFormat(ctx, format); err != nil {
		return nil, fmt.Errorf("invalid import format: %w", err)
	}

	// Read input data
	data, err := io.ReadAll(input)
	if err != nil {
		return nil, fmt.Errorf("failed to read input data: %w", err)
	}

	// Validate data integrity
	if err := m.validator.ValidateDataIntegrity(ctx, data, format); err != nil {
		return nil, fmt.Errorf("data integrity validation failed: %w", err)
	}

	// Deserialize machines
	machines, err := m.backend.DeserializeMachines(ctx, data, format)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize machines: %w", err)
	}

	m.logger.Debug("machines imported successfully",
		spookylogging.Int("count", len(machines)),
		spookylogging.String("format", string(format)))

	return machines, nil
}

// ImportMachinesFromFile imports machines from a file
func (m *Manager) ImportMachinesFromFile(ctx context.Context, filePath string, format spookymachinestypes.ExportFormat) ([]*spookyconfigtypes.Machine, error) {
	m.logger.Info("importing machines from file",
		spookylogging.String("format", string(format)),
		spookylogging.String("file", filePath))

	// Validate file path
	if err := m.backend.ValidateFilePath(ctx, filePath); err != nil {
		return nil, fmt.Errorf("invalid file path: %w", err)
	}

	// Read file data
	data, err := m.backend.ReadImportFile(ctx, filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read import file: %w", err)
	}

	// Validate data integrity
	if err := m.validator.ValidateDataIntegrity(ctx, data, format); err != nil {
		return nil, fmt.Errorf("data integrity validation failed: %w", err)
	}

	// Deserialize machines
	machines, err := m.backend.DeserializeMachines(ctx, data, format)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize machines: %w", err)
	}

	m.logger.Info("machines imported from file successfully",
		spookylogging.Int("count", len(machines)),
		spookylogging.String("format", string(format)),
		spookylogging.String("file", filePath))

	return machines, nil
}

// ImportMachinesWithValidation imports machines with optional validation
func (m *Manager) ImportMachinesWithValidation(ctx context.Context, input io.Reader, format spookymachinestypes.ExportFormat, validate bool) ([]*spookyconfigtypes.Machine, error) {
	m.logger.Debug("importing machines with validation",
		spookylogging.String("format", string(format)),
		spookylogging.Bool("validate", validate))

	// Import machines
	machines, err := m.ImportMachines(ctx, input, format)
	if err != nil {
		return nil, err
	}

	// Validate if requested
	if validate {
		for _, machine := range machines {
			if err := m.validator.ValidateMachineData(ctx, machine); err != nil {
				return nil, fmt.Errorf("machine validation failed for %s: %w", machine.Name, err)
			}
		}
	}

	m.logger.Debug("machines imported with validation successfully",
		spookylogging.Int("count", len(machines)),
		spookylogging.String("format", string(format)),
		spookylogging.Bool("validated", validate))

	return machines, nil
}

// ConfigureExport configures the export system
func (m *Manager) ConfigureExport(ctx context.Context, config *spookymachinestypes.ExportConfig) error {
	m.logger.Info("configuring export system")

	// Store configuration
	if err := m.backend.StoreExportConfig(ctx, config); err != nil {
		return fmt.Errorf("failed to store export configuration: %w", err)
	}

	m.config = config
	m.logger.Info("export system configured successfully")
	return nil
}

// GetExportConfig retrieves the current export configuration
func (m *Manager) GetExportConfig(ctx context.Context) (*spookymachinestypes.ExportConfig, error) {
	return m.backend.LoadExportConfig(ctx)
}

// ValidateExportFormat validates an export format
func (m *Manager) ValidateExportFormat(ctx context.Context, format spookymachinestypes.ExportFormat) error {
	return m.validator.ValidateExportFormat(ctx, format)
}

// GetSupportedFormats returns the list of supported export formats
func (m *Manager) GetSupportedFormats(ctx context.Context) []spookymachinestypes.ExportFormat {
	return []spookymachinestypes.ExportFormat{
		spookymachinestypes.ExportFormatJSON,
		spookymachinestypes.ExportFormatHCL,
	}
}

// validateMachines validates machines using schema system
func (m *Manager) validateMachines(content []byte) error {
	validator := spookyschemas.NewSchemaValidator()
	if err := validator.LoadSchema(spookyschemas.SchemaTypeMachines); err != nil {
		return fmt.Errorf("failed to load machines schema: %w", err)
	}

	// Parse content to interface{} for validation
	var data interface{}
	if err := json.Unmarshal(content, &data); err != nil {
		return fmt.Errorf("failed to parse machines for validation: %w", err)
	}

	if err := validator.ValidateData(data, "machines"); err != nil {
		return fmt.Errorf("machines validation failed: %w", err)
	}
	return nil
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
