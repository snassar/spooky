package types

import (
	"time"
)

// ExportConfig represents configuration for machine export operations
type ExportConfig struct {
	// Format settings
	DefaultFormat   ExportFormat `json:"default_format"`
	IncludeMetadata bool         `json:"include_metadata"`
	IncludeSecrets  bool         `json:"include_secrets"`

	// Output settings
	PrettyPrint bool   `json:"pretty_print"`
	IndentSize  int    `json:"indent_size"`
	OutputDir   string `json:"output_dir"`

	// Filtering settings
	DefaultFilter *MachineFilter `json:"default_filter,omitempty"`
	ExcludeFields []string       `json:"exclude_fields,omitempty"`
	IncludeFields []string       `json:"include_fields,omitempty"`

	// Validation settings
	ValidateOnExport bool `json:"validate_on_export"`
	ValidateOnImport bool `json:"validate_on_import"`

	// Performance settings
	BatchSize   int           `json:"batch_size"`
	Timeout     time.Duration `json:"timeout"`
	MaxFileSize int64         `json:"max_file_size"`

	// Security settings
	EncryptExports bool   `json:"encrypt_exports"`
	EncryptionKey  string `json:"encryption_key,omitempty"`
}

// ImportConfig represents configuration for machine import operations
type ImportConfig struct {
	// Validation settings
	ValidateOnImport bool `json:"validate_on_import"`
	StrictValidation bool `json:"strict_validation"`

	// Conflict resolution
	ConflictResolution ConflictResolution `json:"conflict_resolution"`
	BackupExisting     bool               `json:"backup_existing"`
	BackupDir          string             `json:"backup_dir"`

	// Performance settings
	BatchSize int           `json:"batch_size"`
	Timeout   time.Duration `json:"timeout"`

	// Security settings
	DecryptImports bool   `json:"decrypt_imports"`
	DecryptionKey  string `json:"decryption_key,omitempty"`
}

// ConflictResolution represents how to handle import conflicts
type ConflictResolution string

const (
	ConflictResolutionSkip    ConflictResolution = "skip"
	ConflictResolutionReplace ConflictResolution = "replace"
	ConflictResolutionMerge   ConflictResolution = "merge"
	ConflictResolutionRename  ConflictResolution = "rename"
)

// ImportResult represents the result of an import operation
type ImportResult struct {
	TotalMachines    int               `json:"total_machines"`
	ImportedMachines int               `json:"imported_machines"`
	SkippedMachines  int               `json:"skipped_machines"`
	FailedMachines   int               `json:"failed_machines"`
	Conflicts        []ImportConflict  `json:"conflicts,omitempty"`
	Errors           []ImportError     `json:"errors,omitempty"`
	ImportTime       time.Duration     `json:"import_time"`
	ValidationErrors []ValidationError `json:"validation_errors,omitempty"`
}

// ImportConflict represents a conflict during import
type ImportConflict struct {
	MachineName  string                 `json:"machine_name"`
	ConflictType string                 `json:"conflict_type"`
	ExistingData map[string]interface{} `json:"existing_data,omitempty"`
	NewData      map[string]interface{} `json:"new_data,omitempty"`
	Resolution   ConflictResolution     `json:"resolution"`
}

// ImportError represents an error during import
type ImportError struct {
	MachineName string `json:"machine_name,omitempty"`
	Field       string `json:"field,omitempty"`
	Message     string `json:"message"`
	Line        int    `json:"line,omitempty"`
	Column      int    `json:"column,omitempty"`
}

// ValidationError represents a validation error
type ValidationError struct {
	MachineName string `json:"machine_name,omitempty"`
	Field       string `json:"field"`
	Message     string `json:"message"`
	Severity    string `json:"severity"`
}
