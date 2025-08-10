package logging

import (
	"time"
)

// LogEntry represents a single log entry
type LogEntry struct {
	Timestamp time.Time
	Level     LogLevel
	Message   string
	Fields    []Field
	Error     error
}

// Field represents a structured logging field
type Field struct {
	Key   string
	Value interface{}
}

// Logger defines the main interface for logging operations
type Logger interface {
	// Core logging operations
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	Fatal(msg string, fields ...Field)

	// Field management
	WithField(key string, value interface{}) Logger
	WithFields(fields map[string]interface{}) Logger
	WithError(err error) Logger

	// Configuration
	SetLevel(level LogLevel) error
	GetLevel() LogLevel
	SetFormatter(formatter interface{}) error
	SetOutput(output interface{}) error

	// Utility operations
	Sync() error
	Close() error
}
