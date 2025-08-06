package interfaces

import (
	"fmt"
)

// IntegrationError represents an integration-specific error
type IntegrationError struct {
	Operation string
	System    string
	Message   string
	Cause     error
}

func (ie *IntegrationError) Error() string {
	if ie.Cause != nil {
		return fmt.Sprintf("%s integration %s failed: %s (caused by: %v)", ie.System, ie.Operation, ie.Message, ie.Cause)
	}
	return fmt.Sprintf("%s integration %s failed: %s", ie.System, ie.Operation, ie.Message)
}

func (ie *IntegrationError) Unwrap() error {
	return ie.Cause
}

// NewIntegrationError creates a new integration error
func NewIntegrationError(system, operation, message string, cause error) *IntegrationError {
	return &IntegrationError{
		Operation: operation,
		System:    system,
		Message:   message,
		Cause:     cause,
	}
}

// ContextError represents a context-related error
type ContextError struct {
	ContextType string
	Message     string
	Cause       error
}

func (ce *ContextError) Error() string {
	if ce.Cause != nil {
		return fmt.Sprintf("context error (%s): %s (caused by: %v)", ce.ContextType, ce.Message, ce.Cause)
	}
	return fmt.Sprintf("context error (%s): %s", ce.ContextType, ce.Message)
}

func (ce *ContextError) Unwrap() error {
	return ce.Cause
}

// NewContextError creates a new context error
func NewContextError(contextType, message string, cause error) *ContextError {
	return &ContextError{
		ContextType: contextType,
		Message:     message,
		Cause:       cause,
	}
}

// ValidationError represents validation errors (already defined in validation.go)
// This is just a reference to avoid duplication

// CacheError represents a cache-related error
type CacheError struct {
	Operation string
	Key       string
	Message   string
	Cause     error
}

func (ce *CacheError) Error() string {
	if ce.Cause != nil {
		return fmt.Sprintf("cache %s failed for key '%s': %s (caused by: %v)", ce.Operation, ce.Key, ce.Message, ce.Cause)
	}
	return fmt.Sprintf("cache %s failed for key '%s': %s", ce.Operation, ce.Key, ce.Message)
}

func (ce *CacheError) Unwrap() error {
	return ce.Cause
}

// NewCacheError creates a new cache error
func NewCacheError(operation, key, message string, cause error) *CacheError {
	return &CacheError{
		Operation: operation,
		Key:       key,
		Message:   message,
		Cause:     cause,
	}
}
