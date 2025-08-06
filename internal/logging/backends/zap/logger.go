package zap

import (
	"context"
	"spooky/internal/logging"
	"spooky/internal/logging/types"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ZapLogger wraps Zap's logger to implement our Logger interface
type ZapLogger struct {
	logger *zap.Logger
	fields []types.Field
}

// ZapBackend implements the Backend interface for Zap
type ZapBackend struct{}

// NewZapBackend creates a new Zap backend
func NewZapBackend() *ZapBackend {
	return &ZapBackend{}
}

// CreateLogger creates a new Zap logger with the given configuration
func (b *ZapBackend) CreateLogger(config *types.Config) (logging.Logger, error) {
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

	return &ZapLogger{
		logger: zapLoggerInstance,
		fields: []types.Field{},
	}, nil
}

// GetName returns the backend name
func (b *ZapBackend) GetName() string {
	return "zap"
}

// convertFields converts our Field type to Zap fields
func (l *ZapLogger) convertFields(fields []types.Field) []zap.Field {
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
func (l *ZapLogger) Debug(msg string, fields ...types.Field) {
	l.logger.Debug(msg, l.convertFields(fields)...)
}

// Info logs an info message
func (l *ZapLogger) Info(msg string, fields ...types.Field) {
	l.logger.Info(msg, l.convertFields(fields)...)
}

// Warn logs a warning message
func (l *ZapLogger) Warn(msg string, fields ...types.Field) {
	l.logger.Warn(msg, l.convertFields(fields)...)
}

// Error logs an error message
func (l *ZapLogger) Error(msg string, err error, fields ...types.Field) {
	// Add error field if provided
	if err != nil {
		fields = append(fields, types.Field{Key: "error", Value: err.Error()})
	}
	l.logger.Error(msg, l.convertFields(fields)...)
}

// WithContext returns a logger with context (no-op for Zap implementation)
func (l *ZapLogger) WithContext(_ context.Context) logging.Logger {
	// For Zap, we don't need to modify the logger for context
	// Context handling is done at the application level
	return l
}

// WithFields returns a logger with additional fields
func (l *ZapLogger) WithFields(fields ...types.Field) logging.Logger {
	return &ZapLogger{
		logger: l.logger,
		fields: append(l.fields, fields...),
	}
}

// Sync flushes any buffered log entries
func (l *ZapLogger) Sync() error {
	return l.logger.Sync()
}
