package logging

import (
	"log/slog"
)

// LogLevel constants for consistent usage across the project
const (
	LevelDebug = "debug"
	LevelInfo  = "info"
	LevelWarn  = "warn"
	LevelError = "error"
	LevelFatal = "fatal"
)

// LogLevelMap maps log level names to their numeric values.
var LogLevelMap = map[string]slog.Level{
	LevelDebug: slog.LevelDebug,
	LevelInfo:  slog.LevelInfo,
	LevelWarn:  slog.LevelWarn,
	LevelError: slog.LevelError,
	LevelFatal: slog.LevelError, // slog doesn't have fatal, use error
}

// Standardized error context keys for consistent logging
const (
	// Component identification
	KeyComponent = "component"
	KeyPackage   = "package"
	KeyFunction  = "function"

	// Error context
	KeyError        = "error"
	KeyErrorType    = "error_type"
	KeyErrorCode    = "error_code"
	KeyErrorMessage = "error_message"

	// Operation context
	KeyOperation = "operation"
	KeyAction    = "action"
	KeyTarget    = "target"
	KeySource    = "source"

	// Resource context
	KeyResource = "resource"
	KeyPath     = "path"
	KeyFile     = "file"
	KeyHost     = "host"
	KeyPort     = "port"
	KeyUser     = "user"

	// Performance context
	KeyDuration   = "duration"
	KeyBytes      = "bytes"
	KeyCount      = "count"
	KeySize       = "size"
	KeyOperations = "operations"

	// Security context
	KeyEncrypted = "encrypted"
	KeyKeyType   = "key_type"
	KeyAlgorithm = "algorithm"

	// Network context
	KeyConnection = "connection"
	KeyProtocol   = "protocol"
	KeyTimeout    = "timeout"
	KeyRetries    = "retries"
)

// Standardized error logging functions for consistent usage

// LogError logs an error with standard context
func LogError(logger *Logger, message string, err error, attrs ...any) {
	if err != nil {
		// Add standard error context
		errorAttrs := []any{
			KeyError, err.Error(),
			KeyErrorType, getErrorType(err),
		}

		// Combine with additional attributes
		allAttrs := append(errorAttrs, attrs...)

		logger.Error(message, allAttrs...)
	} else {
		// Log without error context
		logger.Error(message, attrs...)
	}
}

// LogWarn logs a warning with standard context
func LogWarn(logger *Logger, message string, attrs ...any) {
	logger.Warn(message, attrs...)
}

// LogInfo logs info with standard context
func LogInfo(logger *Logger, message string, attrs ...any) {
	logger.Info(message, attrs...)
}

// LogDebug logs debug info with standard context
func LogDebug(logger *Logger, message string, attrs ...any) {
	logger.Debug(message, attrs...)
}

// Helper function to get error type
func getErrorType(err error) string {
	if err == nil {
		return "nil"
	}

	// Check for common error types
	if err.Error() == "" {
		return "empty"
	}

	// Use reflection to get type name
	return "error"
}

// Standardized attribute builders for common contexts

// WithComponent adds component context
func WithComponent(component string) slog.Attr {
	return slog.String(KeyComponent, component)
}

// WithPackage adds package context
func WithPackage(pkg string) slog.Attr {
	return slog.String(KeyPackage, pkg)
}

// WithFunction adds function context
func WithFunction(function string) slog.Attr {
	return slog.String(KeyFunction, function)
}

// WithOperation adds operation context
func WithOperation(operation string) slog.Attr {
	return slog.String(KeyOperation, operation)
}

// WithResource adds resource context
func WithResource(resource string) slog.Attr {
	return slog.String(KeyResource, resource)
}

// WithPath adds path context
func WithPath(path string) slog.Attr {
	return slog.String(KeyPath, path)
}

// WithHost adds host context
func WithHost(host string) slog.Attr {
	return slog.String(KeyHost, host)
}

// WithDuration adds duration context
func WithDuration(duration slog.Value) slog.Attr {
	return slog.Attr{Key: KeyDuration, Value: duration}
}

// WithBytes adds bytes context
func WithBytes(bytes int64) slog.Attr {
	return slog.Int64(KeyBytes, bytes)
}

// WithCount adds count context
func WithCount(count int) slog.Attr {
	return slog.Int(KeyCount, count)
}
