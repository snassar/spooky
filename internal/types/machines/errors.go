// Package machines provides machine inventory error types for the spooky codebase.
package machines

import (
	"fmt"
	"time"

	spookytypescommon "spooky/internal/types/common"
)

// MachineError represents a machine-related error
type MachineError struct {
	spookytypescommon.ErrorDetails

	MachineHostname string `json:"machine_hostname" hcl:"machine_hostname"`
	Operation       string `json:"operation" hcl:"operation"`
	ErrorType       string `json:"error_type" hcl:"error_type"`
}

// NewMachineError creates a new machine error
func NewMachineError(hostname, operation, errorType, message string) *MachineError {
	return &MachineError{
		ErrorDetails: spookytypescommon.ErrorDetails{
			Code:        "MACHINE_ERROR",
			Message:     message,
			Timestamp:   time.Now(),
			Recoverable: true,
		},
		MachineHostname: hostname,
		Operation:       operation,
		ErrorType:       errorType,
	}
}

// Error implements the error interface
func (e *MachineError) Error() string {
	return fmt.Sprintf("machine error [%s]: %s - %s", e.MachineHostname, e.Operation, e.Message)
}

// Unwrap returns the underlying error
func (e *MachineError) Unwrap() error {
	return nil
}

// MachineValidationError represents a machine validation error
type MachineValidationError struct {
	MachineError
	Field string      `json:"field" hcl:"field"`
	Value interface{} `json:"value,omitempty" hcl:"value,optional"`
	Rule  string      `json:"rule" hcl:"rule"`
}

// NewMachineValidationError creates a new machine validation error
func NewMachineValidationError(hostname, field string, value interface{}, rule, message string) *MachineValidationError {
	return &MachineValidationError{
		MachineError: *NewMachineError(hostname, "validation", "validation_error", message),
		Field:        field,
		Value:        value,
		Rule:         rule,
	}
}

// MachineConnectionError represents a machine connection error
type MachineConnectionError struct {
	MachineError
	ConnectionType string `json:"connection_type" hcl:"connection_type"`
	Latency        int    `json:"latency,omitempty" hcl:"latency,optional"`
	RetryCount     int    `json:"retry_count" hcl:"retry_count"`
}

// NewMachineConnectionError creates a new machine connection error
func NewMachineConnectionError(hostname, connectionType string, retryCount int, message string) *MachineConnectionError {
	return &MachineConnectionError{
		MachineError:   *NewMachineError(hostname, "connection", "connection_error", message),
		ConnectionType: connectionType,
		RetryCount:     retryCount,
	}
}

// MachineAuthenticationError represents a machine authentication error
type MachineAuthenticationError struct {
	MachineError
	AuthMethod string `json:"auth_method" hcl:"auth_method"`
	Username   string `json:"username" hcl:"username"`
}

// NewMachineAuthenticationError creates a new machine authentication error
func NewMachineAuthenticationError(hostname, authMethod, username, message string) *MachineAuthenticationError {
	return &MachineAuthenticationError{
		MachineError: *NewMachineError(hostname, "authentication", "authentication_error", message),
		AuthMethod:   authMethod,
		Username:     username,
	}
}

// MachineLoadError represents a machine loading error
type MachineLoadError struct {
	MachineError
	SourcePath string `json:"source_path" hcl:"source_path"`
	FileType   string `json:"file_type" hcl:"file_type"`
}

// NewMachineLoadError creates a new machine load error
func NewMachineLoadError(hostname, sourcePath, fileType, message string) *MachineLoadError {
	return &MachineLoadError{
		MachineError: *NewMachineError(hostname, "loading", "load_error", message),
		SourcePath:   sourcePath,
		FileType:     fileType,
	}
}

// MachinePingError represents a machine ping error
type MachinePingError struct {
	MachineError
	PingType string `json:"ping_type" hcl:"ping_type"`
	Latency  int    `json:"latency,omitempty" hcl:"latency,optional"`
}

// NewMachinePingError creates a new machine ping error
func NewMachinePingError(hostname, pingType string, latency int, message string) *MachinePingError {
	return &MachinePingError{
		MachineError: *NewMachineError(hostname, "ping", "ping_error", message),
		PingType:     pingType,
		Latency:      latency,
	}
}
