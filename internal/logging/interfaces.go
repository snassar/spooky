package logging

import (
	"context"
	"spooky/internal/logging/types"
)

// Logger interface defines the core logging functionality
type Logger interface {
	Debug(msg string, fields ...types.Field)
	Info(msg string, fields ...types.Field)
	Warn(msg string, fields ...types.Field)
	Error(msg string, err error, fields ...types.Field)

	WithContext(ctx context.Context) Logger
	WithFields(fields ...types.Field) Logger
	Sync() error
}

// Formatter interface defines log formatting strategies
type Formatter interface {
	Format(entry *types.LogEntry) ([]byte, error)
	GetName() string
}

// Output interface defines log output destinations
type Output interface {
	Write(data []byte) error
	Close() error
	GetName() string
}

// Backend interface defines logging backend implementations
type Backend interface {
	CreateLogger(config *types.Config) (Logger, error)
	GetName() string
}

// Manager interface defines the logging manager functionality
type Manager interface {
	Configure(config *types.Config) error
	GetLogger() Logger
	SetLogger(logger Logger)
	FromContext(ctx context.Context) Logger
	WithContext(ctx context.Context, logger Logger) context.Context
	Close() error
}
