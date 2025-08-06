package types

import "fmt"

// MachineError represents a machine-related error
type MachineError struct {
	Operation string
	Machine   string
	Cause     error
}

// Error implements the error interface
func (e *MachineError) Error() string {
	if e.Machine != "" {
		return fmt.Sprintf("machine operation '%s' failed for machine '%s': %v", e.Operation, e.Machine, e.Cause)
	}
	return fmt.Sprintf("machine operation '%s' failed: %v", e.Operation, e.Cause)
}

// Unwrap returns the underlying error
func (e *MachineError) Unwrap() error {
	return e.Cause
}

// IndexError represents an index-related error
type IndexError struct {
	Operation string
	IndexType string
	Cause     error
}

// Error implements the error interface
func (e *IndexError) Error() string {
	if e.IndexType != "" {
		return fmt.Sprintf("index operation '%s' failed for index type '%s': %v", e.Operation, e.IndexType, e.Cause)
	}
	return fmt.Sprintf("index operation '%s' failed: %v", e.Operation, e.Cause)
}

// Unwrap returns the underlying error
func (e *IndexError) Unwrap() error {
	return e.Cause
}

// ConnectivityError represents a connectivity-related error
type ConnectivityError struct {
	Operation string
	Machine   string
	Phase     string
	Cause     error
}

// Error implements the error interface
func (e *ConnectivityError) Error() string {
	if e.Machine != "" && e.Phase != "" {
		return fmt.Sprintf("connectivity operation '%s' failed for machine '%s' in phase '%s': %v", e.Operation, e.Machine, e.Phase, e.Cause)
	}
	return fmt.Sprintf("connectivity operation '%s' failed: %v", e.Operation, e.Cause)
}

// Unwrap returns the underlying error
func (e *ConnectivityError) Unwrap() error {
	return e.Cause
}

// ExportError represents an export-related error
type ExportError struct {
	Operation string
	Format    string
	Cause     error
}

// Error implements the error interface
func (e *ExportError) Error() string {
	if e.Format != "" {
		return fmt.Sprintf("export operation '%s' failed for format '%s': %v", e.Operation, e.Format, e.Cause)
	}
	return fmt.Sprintf("export operation '%s' failed: %v", e.Operation, e.Cause)
}

// Unwrap returns the underlying error
func (e *ExportError) Unwrap() error {
	return e.Cause
}
