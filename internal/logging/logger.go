package logging

import (
	"context"
	"os"
	"path/filepath"
	"spooky/internal/logging/types"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// zapLogger wraps Zap's logger to implement our Logger interface
type zapLogger struct {
	logger *zap.Logger
	fields []types.Field
}

// Global logger instance
var globalLogger Logger

// contextKey for storing logger in context
type contextKey struct{}

// ensureLogDirectory ensures the log directory exists
func ensureLogDirectory(logFile string) error {
	if logFile == "" || logFile == "stdout" || logFile == "stderr" {
		return nil // No directory needed for stdout/stderr
	}

	// Get the directory from the log file path
	logDir := filepath.Dir(logFile)

	// Create the directory and all parent directories
	return os.MkdirAll(logDir, 0o755)
}

// ConfigureLogger configures the global logger based on the provided settings
func ConfigureLogger(level, format, output string, quiet, verbose bool) {
	// Determine log level based on flags
	var logLevel types.LogLevel
	switch level {
	case "debug":
		logLevel = types.DebugLevel
	case "info":
		logLevel = types.InfoLevel
	case "warn":
		logLevel = types.WarnLevel
	case "error":
		logLevel = types.ErrorLevel
	default:
		logLevel = types.InfoLevel
	}

	// Override level based on verbose/quiet flags
	if quiet {
		logLevel = types.ErrorLevel // Only show errors when quiet
	} else if verbose {
		logLevel = types.DebugLevel // Show debug when verbose
	}

	// Determine format
	logFormat := "json"
	if format == "text" {
		logFormat = "text"
	}

	// Determine output
	logOutput := "stdout"
	if output != "" {
		logOutput = output
	}

	// Ensure log directory exists
	if err := ensureLogDirectory(logOutput); err != nil {
		// If we can't create the log directory, fall back to stdout
		logOutput = "stdout"
	}

	// Create and set the logger
	globalLogger = NewLogger(types.Config{
		Level:     logLevel,
		Format:    logFormat,
		Output:    logOutput,
		Timestamp: true,
	})
}

func init() {
	// Initialize with a default logger for tests and other cases
	// This will be overridden when ConfigureLogger is called
	globalLogger = NewLogger(types.Config{
		Level:     types.InfoLevel,
		Format:    "json",
		Output:    "stdout",
		Timestamp: true,
	})
}

// NewLogger creates a new logger with the given configuration
func NewLogger(config types.Config) Logger {
	var zapConfig zap.Config

	// Set log level
	level := zap.NewAtomicLevel()
	switch config.Level {
	case types.DebugLevel:
		level.SetLevel(zapcore.DebugLevel)
	case types.InfoLevel:
		level.SetLevel(zapcore.InfoLevel)
	case types.WarnLevel:
		level.SetLevel(zapcore.WarnLevel)
	case types.ErrorLevel:
		level.SetLevel(zapcore.ErrorLevel)
	default:
		level.SetLevel(zapcore.InfoLevel)
	}

	// Configure based on format
	if config.Format == "text" {
		zapConfig = zap.NewDevelopmentConfig()
		zapConfig.Level = level
		zapConfig.OutputPaths = []string{config.Output}
		zapConfig.ErrorOutputPaths = []string{config.Output}
	} else {
		// JSON format (default)
		zapConfig = zap.NewProductionConfig()
		zapConfig.Level = level
		zapConfig.OutputPaths = []string{config.Output}
		zapConfig.ErrorOutputPaths = []string{config.Output}

		// Customize timestamp format
		if config.Timestamp {
			zapConfig.EncoderConfig.TimeKey = "timestamp"
			zapConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		}
	}

	zapLoggerInstance, err := zapConfig.Build()
	if err != nil {
		// Fallback to basic logger if configuration fails
		zapLoggerInstance, _ = zap.NewProduction()
	}

	return &zapLogger{
		logger: zapLoggerInstance,
		fields: []types.Field{},
	}
}

// GetLogger returns the global logger instance
func GetLogger() Logger {
	return globalLogger
}

// SetLogger sets the global logger instance
func SetLogger(logger Logger) {
	globalLogger = logger
}

// FromContext returns a logger from context, or the global logger if not found
func FromContext(ctx context.Context) Logger {
	if logger, ok := ctx.Value(contextKey{}).(Logger); ok {
		return logger
	}
	return globalLogger
}

// WithContext adds a logger to the context
func WithContext(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, contextKey{}, logger)
}

// convertFields converts our Field type to Zap fields
func (l *zapLogger) convertFields(fields []types.Field) []zap.Field {
	zapFields := make([]zap.Field, 0, len(fields))

	for _, field := range fields {
		zapFields = append(zapFields, zap.Any(field.Key, field.Value))
	}

	// Add any fields from the logger instance
	for _, field := range l.fields {
		zapFields = append(zapFields, zap.Any(field.Key, field.Value))
	}

	return zapFields
}

// Debug logs a debug message
func (l *zapLogger) Debug(msg string, fields ...types.Field) {
	l.logger.Debug(msg, l.convertFields(fields)...)
}

// Info logs an info message
func (l *zapLogger) Info(msg string, fields ...types.Field) {
	l.logger.Info(msg, l.convertFields(fields)...)
}

// Warn logs a warning message
func (l *zapLogger) Warn(msg string, fields ...types.Field) {
	l.logger.Warn(msg, l.convertFields(fields)...)
}

// Error logs an error message
func (l *zapLogger) Error(msg string, err error, fields ...types.Field) {
	// Add error field if provided
	if err != nil {
		fields = append(fields, types.Field{Key: "error", Value: err.Error()})
	}
	l.logger.Error(msg, l.convertFields(fields)...)
}

// WithContext returns a logger with context (no-op for Zap implementation)
func (l *zapLogger) WithContext(_ context.Context) Logger {
	// For Zap, we don't need to modify the logger for context
	// Context handling is done at the application level
	return l
}

// WithFields returns a logger with additional fields
func (l *zapLogger) WithFields(fields ...types.Field) Logger {
	return &zapLogger{
		logger: l.logger,
		fields: append(l.fields, fields...),
	}
}

// Sync flushes any buffered log entries
func (l *zapLogger) Sync() error {
	return l.logger.Sync()
}
