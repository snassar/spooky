package facts

import (
	"fmt"
)

// Error types for fact collection
type FactNotFoundError struct {
	Key    string
	Server string
	Source string
}

func (e *FactNotFoundError) Error() string {
	if e.Source != "" {
		return fmt.Sprintf("fact '%s' not found for server '%s' in %s", e.Key, e.Server, e.Source)
	}
	return fmt.Sprintf("fact '%s' not found for server '%s'", e.Key, e.Server)
}

type InvalidSourceError struct {
	Source string
	Reason string
}

func (e *InvalidSourceError) Error() string {
	return fmt.Sprintf("invalid source '%s': %s", e.Source, e.Reason)
}

type CollectionError struct {
	Collector string
	Server    string
	Reason    string
}

func (e *CollectionError) Error() string {
	return fmt.Sprintf("failed to collect facts from %s for server '%s': %s", e.Collector, e.Server, e.Reason)
}

// Standardized error constructors
func NewFactNotFoundError(key, server, source string) error {
	return &FactNotFoundError{
		Key:    key,
		Server: server,
		Source: source,
	}
}

func NewInvalidSourceError(source, reason string) error {
	return &InvalidSourceError{
		Source: source,
		Reason: reason,
	}
}

func NewCollectionError(collector, server, reason string) error {
	return &CollectionError{
		Collector: collector,
		Server:    server,
		Reason:    reason,
	}
}

// Helper functions for common error patterns
func ErrFactNotFound(key, server string) error {
	return NewFactNotFoundError(key, server, "")
}

func ErrFactNotFoundInSource(key, server, source string) error {
	return NewFactNotFoundError(key, server, source)
}

func ErrInvalidSource(source, reason string) error {
	return NewInvalidSourceError(source, reason)
}

func ErrCollectionFailed(collector, server, reason string) error {
	return NewCollectionError(collector, server, reason)
}
