package interfaces

import (
	spookytypeslogging "spooky/internal/types/logging"
)

// Logger defines the main interface for logging operations
type Logger interface {
	// Core logging operations
	Debug(msg string, fields ...spookytypeslogging.Field)
	Info(msg string, fields ...spookytypeslogging.Field)
	Warn(msg string, fields ...spookytypeslogging.Field)
	Error(msg string, fields ...spookytypeslogging.Field)
	Fatal(msg string, fields ...spookytypeslogging.Field)

	// Field management
	WithField(key string, value interface{}) Logger
	WithFields(fields map[string]interface{}) Logger
	WithError(err error) Logger

	// Configuration
	SetLevel(level spookytypeslogging.LogLevel) error
	GetLevel() spookytypeslogging.LogLevel
	SetFormatter(formatter Formatter) error
	SetOutput(output Output) error

	// Utility operations
	Sync() error
	Close() error
}

// Formatter defines the interface for log formatting
type Formatter interface {
	Format(entry *spookytypeslogging.LogEntry) ([]byte, error)
	GetName() string
}

// Output defines the interface for log output
type Output interface {
	Write(data []byte) error
	GetName() string
	Close() error
}

// LogManager defines the interface for log management
type LogManager interface {
	// Core management operations
	GetLogger(name string) Logger
	SetDefaultLogger(logger Logger) error
	GetDefaultLogger() Logger

	// Configuration
	SetDefaultLevel(level spookytypeslogging.LogLevel) error
	SetDefaultFormatter(formatter Formatter) error
	SetDefaultOutput(output Output) error

	// Utility operations
	ListLoggers() []string
	Close() error
}
