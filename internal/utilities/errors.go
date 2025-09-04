package utilities

import (
	"fmt"
	"log/slog"
)

// Common error types for the application
var (
	ErrHCLValidationFailed = fmt.Errorf("HCL validation failed")
	ErrHCLSyntaxError      = fmt.Errorf("HCL syntax error")
	ErrHCLFileNotFound     = fmt.Errorf("HCL file not found")
)

// WrapError wraps an error with a simple context message
func WrapError(err error, context string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", context, err)
}

// WrapErrorf wraps an error with a formatted context message
func WrapErrorf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	context := fmt.Sprintf(format, args...)
	return fmt.Errorf("%s: %w", context, err)
}

// NewError creates a new error with a simple message
func NewError(message string) error {
	return fmt.Errorf("%s", message)
}

// NewErrorf creates a new error with a formatted message
func NewErrorf(format string, args ...interface{}) error {
	return fmt.Errorf(format, args...)
}

// LogError logs an error with context using structured logging
func LogError(err error, context string, logger *slog.Logger) {
	if err == nil {
		return
	}

	if logger == nil {
		logger = slog.Default()
	}

	logger.Error("operation failed",
		slog.String("context", context),
		slog.String("error", err.Error()))
}

// LogErrorf logs an error with formatted context using structured logging
func LogErrorf(err error, logger *slog.Logger, format string, args ...interface{}) {
	if err == nil {
		return
	}

	if logger == nil {
		logger = slog.Default()
	}

	context := fmt.Sprintf(format, args...)
	logger.Error("operation failed",
		slog.String("context", context),
		slog.String("error", err.Error()))
}

// LogWarning logs a warning with context using structured logging
func LogWarning(err error, context string, logger *slog.Logger) {
	if err == nil {
		return
	}

	if logger == nil {
		logger = slog.Default()
	}

	logger.Warn("operation warning",
		slog.String("context", context),
		slog.String("error", err.Error()))
}

// LogWarningf logs a warning with formatted context using structured logging
func LogWarningf(err error, logger *slog.Logger, format string, args ...interface{}) {
	if err == nil {
		return
	}

	if logger == nil {
		logger = slog.Default()
	}

	context := fmt.Sprintf(format, args...)
	logger.Warn("operation warning",
		slog.String("context", context),
		slog.String("error", err.Error()))
}

// NewHCLFileError creates a new HCL file error with context
func NewHCLFileError(filePath, operation, message string) error {
	return fmt.Errorf("HCL file error in %s for %s: %s", operation, filePath, message)
}

// HandleCleanupError handles cleanup errors with appropriate logging
func HandleCleanupError(err error, resource string, operation string) {
	if err == nil {
		return
	}

	logger := slog.Default()
	logger.Warn("cleanup warning",
		slog.String("resource", resource),
		slog.String("operation", operation),
		slog.String("error", err.Error()))
}
