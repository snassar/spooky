package logging

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// zapLogger wraps Zap's logger to implement our Logger interface
type zapLogger struct {
	logger *zap.Logger
	fields []Field
}

// Global logger instance
var globalLogger Logger

// contextKey for storing logger in context
type contextKey struct{}

func init() {
	// Initialize with default configuration
	globalLogger = NewLogger(Config{
		Level:     InfoLevel,
		Format:    "json",
		Output:    "stdout",
		Timestamp: true,
	})
}

// NewLogger creates a new logger with the given configuration
func NewLogger(config Config) Logger {
	var zapConfig zap.Config

	// Set log level
	level := zap.NewAtomicLevel()
	switch config.Level {
	case DebugLevel:
		level.SetLevel(zapcore.DebugLevel)
	case InfoLevel:
		level.SetLevel(zapcore.InfoLevel)
	case WarnLevel:
		level.SetLevel(zapcore.WarnLevel)
	case ErrorLevel:
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
		fields: []Field{},
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
func (l *zapLogger) convertFields(fields []Field) []zap.Field {
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
func (l *zapLogger) Debug(msg string, fields ...Field) {
	l.logger.Debug(msg, l.convertFields(fields)...)
}

// Info logs an info message
func (l *zapLogger) Info(msg string, fields ...Field) {
	l.logger.Info(msg, l.convertFields(fields)...)
}

// Warn logs a warning message
func (l *zapLogger) Warn(msg string, fields ...Field) {
	l.logger.Warn(msg, l.convertFields(fields)...)
}

// Error logs an error message
func (l *zapLogger) Error(msg string, err error, fields ...Field) {
	// Add error field if provided
	if err != nil {
		fields = append(fields, Error(err))
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
func (l *zapLogger) WithFields(fields ...Field) Logger {
	return &zapLogger{
		logger: l.logger,
		fields: append(l.fields, fields...),
	}
}

// Sync flushes any buffered log entries
func (l *zapLogger) Sync() error {
	return l.logger.Sync()
}
