package logging

import (
	"context"
	"os"
	"path/filepath"
	spookytypeslogging "spooky/internal/types/logging"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// zapLogger wraps Zap's logger to implement our Logger interface
type zapLogger struct {
	logger *zap.Logger
	fields []spookytypeslogging.Field
}

// Global logger instance
var globalLogger spookytypeslogging.Logger

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
	var logLevel spookytypeslogging.LogLevel
	switch level {
	case "debug":
		logLevel = spookytypeslogging.DebugLevel
	case "info":
		logLevel = spookytypeslogging.InfoLevel
	case "warn":
		logLevel = spookytypeslogging.WarnLevel
	case "error":
		logLevel = spookytypeslogging.ErrorLevel
	default:
		logLevel = spookytypeslogging.InfoLevel
	}

	// Override level based on verbose/quiet flags
	if quiet {
		logLevel = spookytypeslogging.ErrorLevel // Only show errors when quiet
	} else if verbose {
		logLevel = spookytypeslogging.DebugLevel // Show debug when verbose
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
	globalLogger = NewLogger(spookytypeslogging.Config{
		Level:     logLevel,
		Format:    logFormat,
		Output:    logOutput,
		Timestamp: true,
	})
}

func init() {
	// Initialize with a default logger for tests and other cases
	// This will be overridden when ConfigureLogger is called
	globalLogger = NewLogger(spookytypeslogging.Config{
		Level:     spookytypeslogging.InfoLevel,
		Format:    "json",
		Output:    "stdout",
		Timestamp: true,
	})
}

// NewLogger creates a new logger with the given configuration
func NewLogger(config spookytypeslogging.Config) spookytypeslogging.Logger {
	var zapConfig zap.Config

	// Set log level
	level := zap.NewAtomicLevel()
	switch config.Level {
	case spookytypeslogging.DebugLevel:
		level.SetLevel(zapcore.DebugLevel)
	case spookytypeslogging.InfoLevel:
		level.SetLevel(zapcore.InfoLevel)
	case spookytypeslogging.WarnLevel:
		level.SetLevel(zapcore.WarnLevel)
	case spookytypeslogging.ErrorLevel:
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
		fields: []spookytypeslogging.Field{},
	}
}

// GetLogger returns the global logger instance
func GetLogger() spookytypeslogging.Logger {
	return globalLogger
}

// SetLogger sets the global logger instance
func SetLogger(logger spookytypeslogging.Logger) {
	globalLogger = logger
}

// FromContext returns a logger from context, or the global logger if not found
func FromContext(ctx context.Context) spookytypeslogging.Logger {
	if logger, ok := ctx.Value(contextKey{}).(spookytypeslogging.Logger); ok {
		return logger
	}
	return globalLogger
}

// WithContext adds a logger to the context
func WithContext(ctx context.Context, logger spookytypeslogging.Logger) context.Context {
	return context.WithValue(ctx, contextKey{}, logger)
}

// convertFields converts our Field type to Zap fields
func (l *zapLogger) convertFields(fields []spookytypeslogging.Field) []zap.Field {
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
func (l *zapLogger) Debug(msg string, fields ...spookytypeslogging.Field) {
	l.logger.Debug(msg, l.convertFields(fields)...)
}

// Info logs an info message
func (l *zapLogger) Info(msg string, fields ...spookytypeslogging.Field) {
	l.logger.Info(msg, l.convertFields(fields)...)
}

// Warn logs a warning message
func (l *zapLogger) Warn(msg string, fields ...spookytypeslogging.Field) {
	l.logger.Warn(msg, l.convertFields(fields)...)
}

// Error logs an error message
func (l *zapLogger) Error(msg string, fields ...spookytypeslogging.Field) {
	l.logger.Error(msg, l.convertFields(fields)...)
}

// WithContext returns a logger with context (no-op for Zap implementation)
func (l *zapLogger) WithContext(_ context.Context) spookytypeslogging.Logger {
	// For Zap, we don't need to modify the logger for context
	// Context handling is done at the application level
	return l
}



// Sync flushes any buffered log entries
func (l *zapLogger) Sync() error {
	return l.logger.Sync()
}

// Close closes the logger
func (l *zapLogger) Close() error {
	return l.logger.Sync()
}

// WithField returns a logger with a single field
func (l *zapLogger) WithField(key string, value interface{}) spookytypeslogging.Logger {
	return &zapLogger{
		logger: l.logger,
		fields: append(l.fields, spookytypeslogging.Field{Key: key, Value: value}),
	}
}

// WithFields returns a logger with additional fields from a map
func (l *zapLogger) WithFields(fields map[string]interface{}) spookytypeslogging.Logger {
	newFields := make([]spookytypeslogging.Field, 0, len(fields))
	for key, value := range fields {
		newFields = append(newFields, spookytypeslogging.Field{Key: key, Value: value})
	}
	return &zapLogger{
		logger: l.logger,
		fields: append(l.fields, newFields...),
	}
}

// WithError returns a logger with an error field
func (l *zapLogger) WithError(err error) spookytypeslogging.Logger {
	return l.WithField("error", err.Error())
}

// SetLevel sets the log level
func (l *zapLogger) SetLevel(level spookytypeslogging.LogLevel) error {
	// Zap doesn't support changing log level at runtime for individual loggers
	// This would require recreating the logger
	return nil
}

// GetLevel returns the current log level
func (l *zapLogger) GetLevel() spookytypeslogging.LogLevel {
	// Zap doesn't expose the current level easily
	// Return a default value
	return spookytypeslogging.InfoLevel
}

// SetFormatter sets the formatter
func (l *zapLogger) SetFormatter(formatter interface{}) error {
	// Zap doesn't support changing formatters at runtime
	return nil
}

// SetOutput sets the output
func (l *zapLogger) SetOutput(output interface{}) error {
	// Zap doesn't support changing output at runtime
	return nil
}

// Fatal logs a fatal message and exits
func (l *zapLogger) Fatal(msg string, fields ...spookytypeslogging.Field) {
	l.logger.Fatal(msg, l.convertFields(fields)...)
}
