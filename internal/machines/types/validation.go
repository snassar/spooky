package types

import (
	"time"
)

// ValidationConfig represents configuration for machine validation
type ValidationConfig struct {
	// Validation settings
	Enabled          bool `json:"enabled"`
	StrictValidation bool `json:"strict_validation"`
	ValidateOnLoad   bool `json:"validate_on_load"`
	ValidateOnSave   bool `json:"validate_on_save"`

	// Validation types
	ValidateConfiguration bool `json:"validate_configuration"`
	ValidateConnectivity  bool `json:"validate_connectivity"`
	ValidateSecurity      bool `json:"validate_security"`
	ValidateStructure     bool `json:"validate_structure"`

	// Performance settings
	Timeout       time.Duration `json:"timeout"`
	MaxRetries    int           `json:"max_retries"`
	RetryDelay    time.Duration `json:"retry_delay"`
	ParallelLimit int           `json:"parallel_limit"`

	// Error handling
	FailFast      bool `json:"fail_fast"`
	CollectErrors bool `json:"collect_errors"`
	MaxErrors     int  `json:"max_errors"`
}

// ValidationResult represents the result of a validation operation
type ValidationResult struct {
	// Basic information
	MachineName string        `json:"machine_name"`
	Valid       bool          `json:"valid"`
	Timestamp   time.Time     `json:"timestamp"`
	Duration    time.Duration `json:"duration"`

	// Validation details
	Errors   []ValidationError `json:"errors,omitempty"`
	Warnings []ValidationError `json:"warnings,omitempty"`
	Info     []ValidationError `json:"info,omitempty"`

	// Validation types
	ConfigurationValid bool `json:"configuration_valid"`
	ConnectivityValid  bool `json:"connectivity_valid"`
	SecurityValid      bool `json:"security_valid"`
	StructureValid     bool `json:"structure_valid"`

	// Metadata
	ValidationType string            `json:"validation_type"`
	Metadata       map[string]string `json:"metadata,omitempty"`
}

// ValidationStatus represents the status of the validation system
type ValidationStatus struct {
	Enabled               bool               `json:"enabled"`
	TotalMachines         int                `json:"total_machines"`
	ValidMachines         int                `json:"valid_machines"`
	InvalidMachines       int                `json:"invalid_machines"`
	PendingValidations    int                `json:"pending_validations"`
	LastValidationTime    *time.Time         `json:"last_validation_time,omitempty"`
	AverageValidationTime time.Duration      `json:"average_validation_time"`
	ValidationErrors      int                `json:"validation_errors"`
	ConfigurationValid    bool               `json:"configuration_valid"`
	ConnectivityValid     bool               `json:"connectivity_valid"`
	SecurityValid         bool               `json:"security_valid"`
	StructureValid        bool               `json:"structure_valid"`
	RecentResults         []ValidationResult `json:"recent_results,omitempty"`
}

// ValidationHistory represents validation history for a machine
type ValidationHistory struct {
	MachineName    string            `json:"machine_name"`
	Timestamp      time.Time         `json:"timestamp"`
	Valid          bool              `json:"valid"`
	Duration       time.Duration     `json:"duration"`
	ValidationType string            `json:"validation_type"`
	Errors         []ValidationError `json:"errors,omitempty"`
	Warnings       []ValidationError `json:"warnings,omitempty"`
	Info           []ValidationError `json:"info,omitempty"`
	Metadata       map[string]string `json:"metadata,omitempty"`
}

// ValidationSummary represents a summary of validation results
type ValidationSummary struct {
	TotalMachines   int                `json:"total_machines"`
	ValidMachines   int                `json:"valid_machines"`
	InvalidMachines int                `json:"invalid_machines"`
	TotalErrors     int                `json:"total_errors"`
	TotalWarnings   int                `json:"total_warnings"`
	TotalInfo       int                `json:"total_info"`
	ValidationTime  time.Duration      `json:"validation_time"`
	Results         []ValidationResult `json:"results,omitempty"`
	ErrorSummary    map[string]int     `json:"error_summary,omitempty"`
	WarningSummary  map[string]int     `json:"warning_summary,omitempty"`
}
