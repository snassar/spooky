package utilities

import (
	"fmt"
	"log/slog"

	"github.com/pkg/errors"
)

// ErrorContext provides context for error wrapping
type ErrorContext struct {
	Operation string
	Resource  string
	Component string
	Details   map[string]interface{}
}

// StandardErrorHandler provides consistent error handling patterns
type StandardErrorHandler struct {
	logger *slog.Logger
}

// NewStandardErrorHandler creates a new standard error handler
func NewStandardErrorHandler() *StandardErrorHandler {
	return &StandardErrorHandler{
		logger: slog.Default(),
	}
}

// WrapError wraps an error with consistent context
func (seh *StandardErrorHandler) WrapError(err error, context ErrorContext) error {
	if err == nil {
		return nil
	}

	// Build context message
	contextMsg := context.Operation
	if context.Resource != "" {
		contextMsg += fmt.Sprintf(" for %s", context.Resource)
	}
	if context.Component != "" {
		contextMsg += fmt.Sprintf(" in %s", context.Component)
	}

	// Add details if provided
	if len(context.Details) > 0 {
		for key, value := range context.Details {
			contextMsg += fmt.Sprintf(" (%s: %v)", key, value)
		}
	}

	return errors.Wrap(err, contextMsg)
}

// WrapErrorf wraps an error with formatted context
func (seh *StandardErrorHandler) WrapErrorf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	return errors.Wrapf(err, format, args...)
}

// NewError creates a new error with context
func (seh *StandardErrorHandler) NewError(context ErrorContext, message string) error {
	contextMsg := context.Operation
	if context.Resource != "" {
		contextMsg += fmt.Sprintf(" for %s", context.Resource)
	}
	if context.Component != "" {
		contextMsg += fmt.Sprintf(" in %s", context.Component)
	}

	if message != "" {
		contextMsg += ": " + message
	}

	return errors.New(contextMsg)
}

// NewErrorf creates a new formatted error with context
func (seh *StandardErrorHandler) NewErrorf(context ErrorContext, format string, args ...interface{}) error {
	contextMsg := context.Operation
	if context.Resource != "" {
		contextMsg += fmt.Sprintf(" for %s", context.Resource)
	}
	if context.Component != "" {
		contextMsg += fmt.Sprintf(" in %s", context.Component)
	}

	fullMessage := fmt.Sprintf("%s: %s", contextMsg, fmt.Sprintf(format, args...))
	return errors.New(fullMessage)
}

// HandleCleanupError handles cleanup errors with appropriate severity
func (seh *StandardErrorHandler) HandleCleanupError(err error, resource string, operation string) {
	if err == nil {
		return
	}

	// Determine if this is a critical cleanup failure
	isCritical := seh.isCriticalCleanupError(err, resource, operation)

	if isCritical {
		seh.logger.Error("critical cleanup failure",
			slog.String("resource", resource),
			slog.String("operation", operation),
			slog.String("error", err.Error()))
	} else {
		seh.logger.Warn("cleanup warning",
			slog.String("resource", resource),
			slog.String("operation", operation),
			slog.String("error", err.Error()))
	}
}

// isCriticalCleanupError determines if a cleanup error is critical
func (seh *StandardErrorHandler) isCriticalCleanupError(err error, resource string, operation string) bool {
	if err == nil {
		return false
	}

	errorStr := err.Error()

	// Critical cleanup failures
	criticalPatterns := []string{
		"permission denied",
		"access denied",
		"device or resource busy",
		"directory not empty",
		"file system error",
		"disk full",
		"quota exceeded",
	}

	for _, pattern := range criticalPatterns {
		if containsIgnoreCase(errorStr, pattern) {
			return true
		}
	}

	// Non-critical cleanup failures (temporary files, network connections, etc.)
	nonCriticalPatterns := []string{
		"no such file or directory",
		"connection reset",
		"broken pipe",
		"connection refused",
		"timeout",
	}

	for _, pattern := range nonCriticalPatterns {
		if containsIgnoreCase(errorStr, pattern) {
			return false
		}
	}

	// Default to non-critical for unknown errors
	return false
}

// containsIgnoreCase checks if a string contains a substring (case insensitive)
func containsIgnoreCase(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr ||
			len(s) > len(substr) &&
				(s[:len(substr)] == substr ||
					s[len(s)-len(substr):] == substr ||
					containsSubstring(s, substr)))
}

// containsSubstring checks if s contains substr (case insensitive)
func containsSubstring(s, substr string) bool {
	if len(substr) == 0 {
		return true
	}
	if len(s) < len(substr) {
		return false
	}

	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// Global error handler instance
var globalErrorHandler = NewStandardErrorHandler()

// WrapError is a convenience function for global error wrapping
func WrapError(err error, context ErrorContext) error {
	return globalErrorHandler.WrapError(err, context)
}

// WrapErrorf is a convenience function for global error wrapping with format
func WrapErrorf(err error, format string, args ...interface{}) error {
	return globalErrorHandler.WrapErrorf(err, format, args...)
}

// NewError is a convenience function for creating new errors
func NewError(context ErrorContext, message string) error {
	return globalErrorHandler.NewError(context, message)
}

// NewErrorf is a convenience function for creating formatted errors
func NewErrorf(context ErrorContext, format string, args ...interface{}) error {
	return globalErrorHandler.NewErrorf(context, format, args...)
}

// HandleCleanupError is a convenience function for handling cleanup errors
func HandleCleanupError(err error, resource string, operation string) {
	globalErrorHandler.HandleCleanupError(err, resource, operation)
}
