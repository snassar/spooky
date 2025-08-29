package utilities

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
)

// ErrorAggregator collects and manages multiple errors
type ErrorAggregator struct {
	errors []error
	limit  int
}

// NewErrorAggregator creates a new error aggregator
func NewErrorAggregator(limit int) *ErrorAggregator {
	return &ErrorAggregator{
		errors: make([]error, 0),
		limit:  limit,
	}
}

// Add adds an error to the aggregator
func (ea *ErrorAggregator) Add(err error) {
	if err == nil {
		return
	}

	if len(ea.errors) < ea.limit {
		ea.errors = append(ea.errors, err)
	}
}

// AddWithContext adds an error with context
func (ea *ErrorAggregator) AddWithContext(err error, context string) {
	if err == nil {
		return
	}

	wrappedErr := errors.Wrap(err, context)
	ea.Add(wrappedErr)
}

// HasErrors checks if any errors have been collected.
//
// Returns true if errors exist, false otherwise.
//
// Example:
//
//	if aggregator.HasErrors() {
//	    fmt.Println("Errors found:", aggregator.ErrorCount())
//	}
func (ea *ErrorAggregator) HasErrors() bool {
	return len(ea.errors) > 0
}

// ErrorCount returns the number of errors collected
func (ea *ErrorAggregator) ErrorCount() int {
	return len(ea.errors)
}

// Errors returns all collected errors
func (ea *ErrorAggregator) Errors() []error {
	return ea.errors
}

// Error returns a combined error message
func (ea *ErrorAggregator) Error() string {
	if len(ea.errors) == 0 {
		return ""
	}

	if len(ea.errors) == 1 {
		return ea.errors[0].Error()
	}

	return fmt.Sprintf("multiple errors occurred (%d total): %s",
		len(ea.errors), ea.errors[0].Error())
}

// CombinedError returns a single error combining all collected errors
func (ea *ErrorAggregator) CombinedError() error {
	if len(ea.errors) == 0 {
		return nil
	}

	if len(ea.errors) == 1 {
		return ea.errors[0]
	}

	// Create a combined error with the first error as the base
	combined := errors.Wrap(ea.errors[0], fmt.Sprintf("and %d more errors", len(ea.errors)-1))
	return combined
}

// RetryConfig defines retry behavior
type RetryConfig struct {
	MaxAttempts     int
	InitialDelay    time.Duration
	MaxDelay        time.Duration
	BackoffFactor   float64
	RetryableErrors []error
}

// DefaultRetryConfig returns a sensible default retry configuration
func DefaultRetryConfig() *RetryConfig {
	return &RetryConfig{
		MaxAttempts:   3,
		InitialDelay:  100 * time.Millisecond,
		MaxDelay:      5 * time.Second,
		BackoffFactor: 2.0,
	}
}

// RetryableOperation defines an operation that can be retried
type RetryableOperation func() error

// RetryWithBackoff executes an operation with exponential backoff retry
func RetryWithBackoff(operation RetryableOperation, config *RetryConfig) error {
	if config == nil {
		config = DefaultRetryConfig()
	}

	var lastErr error
	delay := config.InitialDelay

	for attempt := 1; attempt <= config.MaxAttempts; attempt++ {
		err := operation()
		if err == nil {
			return nil
		}

		lastErr = err

		// Check if this is the last attempt
		if attempt == config.MaxAttempts {
			break
		}

		// Check if error is retryable
		if !isRetryableError(err, config.RetryableErrors) {
			break
		}

		// Wait before retry
		time.Sleep(delay)

		// Calculate next delay with exponential backoff
		delay = time.Duration(float64(delay) * config.BackoffFactor)
		if delay > config.MaxDelay {
			delay = config.MaxDelay
		}
	}

	return errors.Wrapf(lastErr, "operation failed after %d attempts", config.MaxAttempts)
}

// isRetryableError checks if an error is retryable
func isRetryableError(err error, retryableErrors []error) bool {
	if len(retryableErrors) == 0 {
		// Default retryable errors
		return isDefaultRetryableError(err)
	}

	for _, retryableErr := range retryableErrors {
		if errors.Is(err, retryableErr) {
			return true
		}
	}

	return false
}

// isDefaultRetryableError checks if an error is retryable by default
func isDefaultRetryableError(err error) bool {
	// Add logic to identify retryable errors
	// For now, we'll retry most errors except validation errors

	// Don't retry validation errors
	if errors.Is(err, ErrHCLValidationFailed) ||
		errors.Is(err, ErrHCLSyntaxError) {
		return false
	}

	// Don't retry file not found errors
	if errors.Is(err, ErrHCLFileNotFound) {
		return false
	}

	return true
}

// BatchProcessor processes items in batches with error aggregation
type BatchProcessor struct {
	BatchSize  int
	Aggregator *ErrorAggregator
}

// NewBatchProcessor creates a new batch processor
func NewBatchProcessor(batchSize int, errorLimit int) *BatchProcessor {
	return &BatchProcessor{
		BatchSize:  batchSize,
		Aggregator: NewErrorAggregator(errorLimit),
	}
}

// ProcessItems processes items in batches
func (bp *BatchProcessor) ProcessItems(items []string, processor func(string) error) error {
	for i := 0; i < len(items); i += bp.BatchSize {
		end := i + bp.BatchSize
		if end > len(items) {
			end = len(items)
		}

		batch := items[i:end]
		for _, item := range batch {
			if err := processor(item); err != nil {
				bp.Aggregator.AddWithContext(err, fmt.Sprintf("processing item: %s", item))
			}
		}
	}

	return bp.Aggregator.CombinedError()
}

// ErrorRecoveryStrategy defines how to recover from errors
type ErrorRecoveryStrategy func(error) error

// RecoverableError represents an error that can be recovered from
type RecoverableError struct {
	OriginalError error
	RecoveryFunc  ErrorRecoveryStrategy
}

func (e *RecoverableError) Error() string {
	return e.OriginalError.Error()
}

func (e *RecoverableError) Unwrap() error {
	return e.OriginalError
}

// NewRecoverableError creates a new recoverable error
func NewRecoverableError(err error, recovery ErrorRecoveryStrategy) error {
	return &RecoverableError{
		OriginalError: err,
		RecoveryFunc:  recovery,
	}
}

// TryRecover attempts to recover from an error
func TryRecover(err error) error {
	var recoverableErr *RecoverableError
	if errors.As(err, &recoverableErr) && recoverableErr.RecoveryFunc != nil {
		return recoverableErr.RecoveryFunc(recoverableErr.OriginalError)
	}
	return err
}
