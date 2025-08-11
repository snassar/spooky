// Package logging provides types for logging operations in the spooky codebase.
// These types define the structure for logging configuration and log entries.
package logging

import (
	"time"
)

// LogLevel represents the logging level
type LogLevel string

const (
	// LogLevelDebug represents debug level logging
	LogLevelDebug LogLevel = "debug"

	// LogLevelInfo represents info level logging
	LogLevelInfo LogLevel = "info"

	// LogLevelWarn represents warning level logging
	LogLevelWarn LogLevel = "warn"

	// LogLevelError represents error level logging
	LogLevelError LogLevel = "error"

	// LogLevelFatal represents fatal level logging
	LogLevelFatal LogLevel = "fatal"
)

// LogConfig represents logging configuration
type LogConfig struct {
	// Log level
	Level LogLevel `json:"level" hcl:"level"`

	// Log format (json, text, structured)
	Format string `json:"format" hcl:"format"`

	// Output destination (stdout, stderr, file)
	Output string `json:"output" hcl:"output"`

	// Log file path (if output is file)
	FilePath string `json:"file_path,omitempty" hcl:"file_path,optional"`

	// Whether to include timestamps
	IncludeTimestamp bool `json:"include_timestamp" hcl:"include_timestamp"`

	// Whether to include caller information
	IncludeCaller bool `json:"include_caller" hcl:"include_caller"`

	// Whether to include stack traces
	IncludeStackTrace bool `json:"include_stack_trace" hcl:"include_stack_trace"`

	// Log rotation settings
	Rotation *LogRotation `json:"rotation,omitempty" hcl:"rotation,optional"`

	// Log filtering settings
	Filtering *LogFiltering `json:"filtering,omitempty" hcl:"filtering,optional"`

	// Log performance settings
	Performance *LogPerformance `json:"performance,omitempty" hcl:"performance,optional"`
}

// LogRotation represents log rotation configuration
type LogRotation struct {
	// Whether rotation is enabled
	Enabled bool `json:"enabled" hcl:"enabled"`

	// Maximum file size in bytes
	MaxSize int64 `json:"max_size" hcl:"max_size"`

	// Maximum number of backup files
	MaxBackups int `json:"max_backups" hcl:"max_backups"`

	// Maximum age of backup files
	MaxAge time.Duration `json:"max_age" hcl:"max_age"`

	// Whether to compress backup files
	Compress bool `json:"compress" hcl:"compress"`
}

// LogFiltering represents log filtering configuration
type LogFiltering struct {
	// Whether filtering is enabled
	Enabled bool `json:"enabled" hcl:"enabled"`

	// Include patterns (regex)
	IncludePatterns []string `json:"include_patterns,omitempty" hcl:"include_patterns,optional"`

	// Exclude patterns (regex)
	ExcludePatterns []string `json:"exclude_patterns,omitempty" hcl:"exclude_patterns,optional"`

	// Minimum log level for filtering
	MinLevel LogLevel `json:"min_level" hcl:"min_level"`

	// Whether to filter by component
	FilterByComponent bool `json:"filter_by_component" hcl:"filter_by_component"`

	// Component filters
	ComponentFilters map[string]LogLevel `json:"component_filters,omitempty" hcl:"component_filters,optional"`
}

// LogPerformance represents log performance configuration
type LogPerformance struct {
	// Whether to use buffered logging
	Buffered bool `json:"buffered" hcl:"buffered"`

	// Buffer size in bytes
	BufferSize int `json:"buffer_size" hcl:"buffer_size"`

	// Flush interval
	FlushInterval time.Duration `json:"flush_interval" hcl:"flush_interval"`

	// Whether to use async logging
	Async bool `json:"async" hcl:"async"`

	// Async queue size
	QueueSize int `json:"queue_size" hcl:"queue_size"`

	// Whether to drop logs when queue is full
	DropWhenFull bool `json:"drop_when_full" hcl:"drop_when_full"`
}

// LogEntry represents a log entry
type LogEntry struct {
	// Log timestamp
	Timestamp time.Time `json:"timestamp" hcl:"timestamp"`

	// Log level
	Level LogLevel `json:"level" hcl:"level"`

	// Log message
	Message string `json:"message" hcl:"message"`

	// Log component
	Component string `json:"component,omitempty" hcl:"component,optional"`

	// Log operation
	Operation string `json:"operation,omitempty" hcl:"operation,optional"`

	// Log fields
	Fields map[string]interface{} `json:"fields,omitempty" hcl:"fields,optional"`

	// Error information
	Error *LogError `json:"error,omitempty" hcl:"error,optional"`

	// Caller information
	Caller *LogCaller `json:"caller,omitempty" hcl:"caller,optional"`

	// Stack trace
	StackTrace []string `json:"stack_trace,omitempty" hcl:"stack_trace,optional"`
}

// LogError represents error information in a log entry
type LogError struct {
	// Error message
	Message string `json:"message" hcl:"message"`

	// Error type
	Type string `json:"type" hcl:"type"`

	// Error code
	Code string `json:"code,omitempty" hcl:"code,optional"`

	// Error context
	Context map[string]interface{} `json:"context,omitempty" hcl:"context,optional"`

	// Whether error is recoverable
	Recoverable bool `json:"recoverable" hcl:"recoverable"`
}

// LogCaller represents caller information in a log entry
type LogCaller struct {
	// File name
	File string `json:"file" hcl:"file"`

	// Line number
	Line int `json:"line" hcl:"line"`

	// Function name
	Function string `json:"function" hcl:"function"`

	// Package name
	Package string `json:"package" hcl:"package"`
}

// Logger provides logging functionality
type Logger interface {
	// Debug logs a debug message
	Debug(msg string, fields ...map[string]interface{})

	// Info logs an info message
	Info(msg string, fields ...map[string]interface{})

	// Warn logs a warning message
	Warn(msg string, fields ...map[string]interface{})

	// Error logs an error message
	Error(msg string, err error, fields ...map[string]interface{})

	// Fatal logs a fatal message and exits
	Fatal(msg string, err error, fields ...map[string]interface{})

	// WithFields returns a logger with additional fields
	WithFields(fields map[string]interface{}) Logger

	// WithComponent returns a logger with a component name
	WithComponent(component string) Logger

	// WithOperation returns a logger with an operation name
	WithOperation(operation string) Logger

	// SetLevel sets the log level
	SetLevel(level LogLevel)

	// GetLevel returns the current log level
	GetLevel() LogLevel
}

// LogManager provides logging management functionality
type LogManager interface {
	// GetLogger returns a logger for the given component
	GetLogger(component string) Logger

	// SetLevel sets the log level for all loggers
	SetLevel(level LogLevel)

	// GetLevel returns the current log level
	GetLevel() LogLevel

	// Configure configures logging with the given configuration
	Configure(config *LogConfig) error

	// Flush flushes all pending log entries
	Flush() error

	// Close closes the log manager
	Close() error
}
