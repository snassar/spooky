package loading

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"spooky/internal/logging"
	"spooky/internal/schemas"
	"spooky/internal/types"
	spookytypeslogging "spooky/internal/types/logging"
	spookytypesvariables "spooky/internal/types/variables"
)

// VariableLoader defines the interface for specific format loaders
type VariableLoader interface {
	Load(ctx context.Context, source interface{}) ([]*types.Variable, error)
	GetName() string
	GetSupportedExtensions() []string
	ValidateSchema(content []byte) error
}

// Logger defines the interface for logging operations
type Logger interface {
	Debug(msg string, fields ...spookytypeslogging.Field)
	Info(msg string, fields ...spookytypeslogging.Field)
	Warn(msg string, fields ...spookytypeslogging.Field)
	Error(msg string, fields ...spookytypeslogging.Field)
	Fatal(msg string, fields ...spookytypeslogging.Field)
}

// Manager implements LoadingManager interface
type Manager struct {
	config  *spookytypesvariables.LoadingConfig
	loaders map[string]VariableLoader
	logger  Logger
}

// NewManager creates a new loading manager
func NewManager(config *spookytypesvariables.LoadingConfig, logger Logger) *Manager {
	manager := &Manager{
		config:  config,
		loaders: make(map[string]VariableLoader),
		logger:  logger,
	}

	// Register default loaders
	validator := schemas.NewSchemaValidator()
	manager.RegisterLoader(".hcl", NewHCLVariableLoader(validator))
	manager.RegisterLoader(".json", NewJSONVariableLoader(validator))

	return manager
}

// RegisterLoader registers a loader for a specific file extension
func (m *Manager) RegisterLoader(extension string, loader VariableLoader) {
	m.loaders[extension] = loader
}

// LoadFromFile loads variables from a single file
func (m *Manager) LoadFromFile(ctx context.Context, path string) ([]*types.Variable, error) {
	// 1. Validate file exists and is readable
	if err := m.ValidateFile(path); err != nil {
		return nil, fmt.Errorf("file validation failed: %w", err)
	}

	// 2. Determine file format and get appropriate loader
	ext := strings.ToLower(filepath.Ext(path))
	loader, exists := m.loaders[ext]
	if !exists {
		return nil, fmt.Errorf("unsupported file format: %s", ext)
	}

	// 3. Load schema for validation
	if err := m.loadSchemaForFormat(ext); err != nil {
		return nil, fmt.Errorf("failed to load schema: %w", err)
	}

	// 4. Validate file against schema
	if err := m.validateFileAgainstSchema(path, ext); err != nil {
		return nil, fmt.Errorf("schema validation failed: %w", err)
	}

	// 5. Load variables using appropriate loader
	return loader.Load(ctx, path)
}

// LoadFromDirectory loads variables from all supported files in a directory
func (m *Manager) LoadFromDirectory(ctx context.Context, dirPath string) ([]*types.Variable, error) {
	m.logger.Debug("Loading variables from directory", logging.String("path", dirPath))

	var allVariables []*types.Variable

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		if ext != ".hcl" && ext != ".json" {
			return nil
		}

		variables, err := m.LoadFromFile(ctx, path)
		if err != nil {
			m.logger.Warn("Failed to load variables from file",
				logging.String("path", path),
				logging.Error(err))
			return nil // Continue with other files
		}

		allVariables = append(allVariables, variables...)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk directory %s: %w", dirPath, err)
	}

	return allVariables, nil
}

// LoadFromHCL loads variables from HCL content
func (m *Manager) LoadFromHCL(ctx context.Context, content []byte) ([]*types.Variable, error) {
	loader, exists := m.loaders[".hcl"]
	if !exists {
		return nil, fmt.Errorf("HCL loader not available")
	}

	return loader.Load(ctx, content)
}

// LoadFromJSON loads variables from JSON content
func (m *Manager) LoadFromJSON(ctx context.Context, content []byte) ([]*types.Variable, error) {
	loader, exists := m.loaders[".json"]
	if !exists {
		return nil, fmt.Errorf("JSON loader not available")
	}

	return loader.Load(ctx, content)
}

// SetDefaultEncoding sets the default encoding for file operations
func (m *Manager) SetDefaultEncoding(encoding string) error {
	if m.config == nil {
		m.config = &spookytypesvariables.LoadingConfig{}
	}
	m.config.DefaultEncoding = encoding
	return nil
}

// SetMaxFileSize sets the maximum file size for loading
func (m *Manager) SetMaxFileSize(maxSize int64) error {
	if m.config == nil {
		m.config = &spookytypesvariables.LoadingConfig{}
	}
	m.config.MaxFileSize = maxSize
	return nil
}

// SetAllowedExtensions sets the allowed file extensions
func (m *Manager) SetAllowedExtensions(extensions []string) error {
	if m.config == nil {
		m.config = &spookytypesvariables.LoadingConfig{}
	}
	m.config.AllowedExtensions = extensions
	return nil
}

// ValidateFile validates that a file exists and is readable
func (m *Manager) ValidateFile(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("file does not exist or is not accessible: %w", err)
	}

	if info.IsDir() {
		return fmt.Errorf("path is a directory, expected a file")
	}

	if m.config != nil && m.config.MaxFileSize > 0 && info.Size() > m.config.MaxFileSize {
		return fmt.Errorf("file size %d exceeds maximum allowed size %d", info.Size(), m.config.MaxFileSize)
	}

	return nil
}

// GetSupportedFormats returns the list of supported file formats
func (m *Manager) GetSupportedFormats() []string {
	formats := make([]string, 0, len(m.loaders))
	for ext := range m.loaders {
		formats = append(formats, ext)
	}
	return formats
}

// Close closes the loading manager and releases resources
func (m *Manager) Close() error {
	// Clean up resources if needed
	return nil
}

// Helper methods for schema validation
func (m *Manager) loadSchemaForFormat(format string) error {
	// Load appropriate schema based on format
	switch format {
	case ".hcl":
		validator := schemas.NewSchemaValidator()
		return validator.LoadSchema(schemas.SchemaTypeVariablesHCL)
	case ".json":
		validator := schemas.NewSchemaValidator()
		return validator.LoadSchema(schemas.SchemaTypeVariablesJSON)
	default:
		return fmt.Errorf("unsupported format for schema loading: %s", format)
	}
}

func (m *Manager) validateFileAgainstSchema(filePath, format string) error {
	// Validate file against loaded schema
	schemaName := ""
	switch format {
	case ".hcl":
		schemaName = "variables-hcl"
	case ".json":
		schemaName = "variables-json"
	default:
		return fmt.Errorf("unsupported format for validation: %s", format)
	}

	validator := schemas.NewSchemaValidator()
	result := validator.ValidateFile(filePath, schemaName)
	if !result.Valid {
		return fmt.Errorf("schema validation failed: %v", result.Errors)
	}

	return nil
}

// ValidateFileAgainstSchema validates a file against a specific schema
func (m *Manager) ValidateFileAgainstSchema(path string, schemaType schemas.SchemaType) error {
	validator := schemas.NewSchemaValidator()
	if err := validator.LoadSchema(schemaType); err != nil {
		return fmt.Errorf("failed to load schema: %w", err)
	}

	result := validator.ValidateFile(path, string(schemaType))
	if !result.Valid {
		return fmt.Errorf("schema validation failed: %v", result.Errors)
	}

	return nil
}

// ValidateContentAgainstSchema validates content against a specific schema
func (m *Manager) ValidateContentAgainstSchema(content []byte, schemaType schemas.SchemaType) error {
	// For content validation, we need to write to a temporary file first
	// since the validator expects a file path
	tmpFile, err := os.CreateTemp("", "variables-*.hcl")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	if _, err := tmpFile.Write(content); err != nil {
		return fmt.Errorf("failed to write content to temp file: %w", err)
	}

	validator := schemas.NewSchemaValidator()
	if err := validator.LoadSchema(schemaType); err != nil {
		return fmt.Errorf("failed to load schema: %w", err)
	}

	result := validator.ValidateFile(tmpFile.Name(), string(schemaType))
	if !result.Valid {
		return fmt.Errorf("schema validation failed: %v", result.Errors)
	}

	return nil
}
