package logging

import (
	"context"
)

// Logger interface defines the core logging functionality
type Logger interface {
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, err error, fields ...Field)

	WithContext(ctx context.Context) Logger
	WithFields(fields ...Field) Logger
}

// Field represents a structured logging field
type Field struct {
	Key   string
	Value interface{}
}

// LogLevel represents the logging level
type LogLevel string

const (
	DebugLevel LogLevel = "debug"
	InfoLevel  LogLevel = "info"
	WarnLevel  LogLevel = "warn"
	ErrorLevel LogLevel = "error"
)

// Config holds logging configuration
type Config struct {
	Level     LogLevel `env:"LOG_LEVEL" default:"info"`
	Format    string   `env:"LOG_FORMAT" default:"json"`
	Output    string   `env:"LOG_OUTPUT" default:"stdout"`
	Timestamp bool     `env:"LOG_TIMESTAMP" default:"true"`
}
