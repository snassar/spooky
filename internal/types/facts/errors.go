package facts

import (
	"fmt"
)

// FactError represents a fact-related error
type FactError struct {
	Operation string
	Server    string
	Key       string
	Message   string
	Cause     error
}

// Error implements the error interface
func (e *FactError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("fact error [%s] for server %s, key %s: %s: %v",
			e.Operation, e.Server, e.Key, e.Message, e.Cause)
	}
	return fmt.Sprintf("fact error [%s] for server %s, key %s: %s",
		e.Operation, e.Server, e.Key, e.Message)
}

// Unwrap returns the underlying error
func (e *FactError) Unwrap() error {
	return e.Cause
}

// CollectionError represents a fact collection error
type CollectionError struct {
	Server  string
	Source  string
	Message string
	Cause   error
}

// Error implements the error interface
func (e *CollectionError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("collection error for server %s from %s: %s: %v",
			e.Server, e.Source, e.Message, e.Cause)
	}
	return fmt.Sprintf("collection error for server %s from %s: %s",
		e.Server, e.Source, e.Message)
}

// Unwrap returns the underlying error
func (e *CollectionError) Unwrap() error {
	return e.Cause
}

// StorageError represents a storage-related error
type StorageError struct {
	Operation string
	Path      string
	Message   string
	Cause     error
}

// Error implements the error interface
func (e *StorageError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("storage error [%s] at %s: %s: %v",
			e.Operation, e.Path, e.Message, e.Cause)
	}
	return fmt.Sprintf("storage error [%s] at %s: %s",
		e.Operation, e.Path, e.Message)
}

// Unwrap returns the underlying error
func (e *StorageError) Unwrap() error {
	return e.Cause
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string
	Value   interface{}
	Message string
	Cause   error
}

// Error implements the error interface
func (e *ValidationError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("validation error for field %s with value %v: %s: %v",
			e.Field, e.Value, e.Message, e.Cause)
	}
	return fmt.Sprintf("validation error for field %s with value %v: %s",
		e.Field, e.Value, e.Message)
}

// Unwrap returns the underlying error
func (e *ValidationError) Unwrap() error {
	return e.Cause
}
