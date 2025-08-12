// Package logging provides types for logging operations in the spooky codebase.
// These types define the structure for logging configuration and log entries.
package logging

import (
	"log/slog"
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

// ToSlogLevel converts LogLevel to slog.Level
func (l LogLevel) ToSlogLevel() slog.Level {
	switch l {
	case LogLevelDebug:
		return slog.LevelDebug
	case LogLevelInfo:
		return slog.LevelInfo
	case LogLevelWarn:
		return slog.LevelWarn
	case LogLevelError:
		return slog.LevelError
	case LogLevelFatal:
		return slog.LevelError // slog doesn't have fatal, use error
	default:
		return slog.LevelInfo
	}
}

// FromSlogLevel converts slog.Level to LogLevel
func FromSlogLevel(level slog.Level) LogLevel {
	switch level {
	case slog.LevelDebug:
		return LogLevelDebug
	case slog.LevelInfo:
		return LogLevelInfo
	case slog.LevelWarn:
		return LogLevelWarn
	case slog.LevelError:
		return LogLevelError
	default:
		return LogLevelInfo
	}
}

// LogConfig represents logging configuration
type LogConfig struct {
	// Log level
	Level LogLevel `json:"level" hcl:"level"`

	// Log format (json, text, structured)
	Format string `json:"format" hcl:"format"`

	// Output destination (stdout, stderr, file, null)
	Output string `json:"output" hcl:"output"`

	// File output configuration
	File *LogFileConfig `json:"file,omitempty" hcl:"file,optional"`

	// Structured logging configuration
	Structured *LogStructuredConfig `json:"structured,omitempty" hcl:"structured,optional"`

	// Performance configuration
	Performance *LogPerformanceConfig `json:"performance,omitempty" hcl:"performance,optional"`

	// Filtering configuration
	Filtering *LogFilteringConfig `json:"filtering,omitempty" hcl:"filtering,optional"`

	// Rotation configuration
	Rotation *LogRotationConfig `json:"rotation,omitempty" hcl:"rotation,optional"`
}

// LogFileConfig represents file output configuration
type LogFileConfig struct {
	// Path to log file
	Path string `json:"path" hcl:"path"`

	// File permissions in octal format
	Permissions string `json:"permissions,omitempty" hcl:"permissions,optional"`

	// Whether to append to existing file
	Append bool `json:"append,omitempty" hcl:"append,optional"`
}

// LogStructuredConfig represents structured logging configuration
type LogStructuredConfig struct {
	// Timestamp configuration
	Timestamp *LogTimestampConfig `json:"timestamp,omitempty" hcl:"timestamp,optional"`

	// Level configuration
	Level *LogLevelConfig `json:"level,omitempty" hcl:"level,optional"`

	// Message configuration
	Message *LogMessageConfig `json:"message,omitempty" hcl:"message,optional"`

	// Component configuration
	Component *LogComponentConfig `json:"component,omitempty" hcl:"component,optional"`

	// Operation configuration
	Operation *LogOperationConfig `json:"operation,omitempty" hcl:"operation,optional"`

	// Error configuration
	Error *LogErrorConfig `json:"error,omitempty" hcl:"error,optional"`

	// Caller configuration
	Caller *LogCallerConfig `json:"caller,omitempty" hcl:"caller,optional"`

	// Fields configuration
	Fields *LogFieldsConfig `json:"fields,omitempty" hcl:"fields,optional"`
}

// LogTimestampConfig represents timestamp configuration
type LogTimestampConfig struct {
	// Whether to include timestamps
	Enabled bool `json:"enabled,omitempty" hcl:"enabled,optional"`

	// Timestamp format
	Format string `json:"format,omitempty" hcl:"format,optional"`

	// Timezone for timestamps
	Timezone string `json:"timezone,omitempty" hcl:"timezone,optional"`
}

// LogLevelConfig represents level configuration
type LogLevelConfig struct {
	// Field key for log level
	Key string `json:"key,omitempty" hcl:"key,optional"`

	// Whether to use uppercase level names
	Uppercase bool `json:"uppercase,omitempty" hcl:"uppercase,optional"`

	// Whether to include color codes
	Color bool `json:"color,omitempty" hcl:"color,optional"`
}

// LogMessageConfig represents message configuration
type LogMessageConfig struct {
	// Field key for log message
	Key string `json:"key,omitempty" hcl:"key,optional"`

	// Maximum message length (0 for no truncation)
	Truncate int `json:"truncate,omitempty" hcl:"truncate,optional"`
}

// LogComponentConfig represents component configuration
type LogComponentConfig struct {
	// Field key for component name
	Key string `json:"key,omitempty" hcl:"key,optional"`

	// Whether to include component information
	Enabled bool `json:"enabled,omitempty" hcl:"enabled,optional"`

	// Whether to include package information
	IncludePackage bool `json:"include_package,omitempty" hcl:"include_package,optional"`
}

// LogOperationConfig represents operation configuration
type LogOperationConfig struct {
	// Field key for operation name
	Key string `json:"key,omitempty" hcl:"key,optional"`

	// Whether to include operation information
	Enabled bool `json:"enabled,omitempty" hcl:"enabled,optional"`

	// Whether to include operation ID
	IncludeID bool `json:"include_id,omitempty" hcl:"include_id,optional"`
}

// LogErrorConfig represents error configuration
type LogErrorConfig struct {
	// Field key for error information
	Key string `json:"key,omitempty" hcl:"key,optional"`

	// Whether to include stack traces
	IncludeStack bool `json:"include_stack,omitempty" hcl:"include_stack,optional"`

	// Whether to include error type information
	IncludeType bool `json:"include_type,omitempty" hcl:"include_type,optional"`

	// Whether to include error codes
	IncludeCode bool `json:"include_code,omitempty" hcl:"include_code,optional"`
}

// LogCallerConfig represents caller configuration
type LogCallerConfig struct {
	// Whether to include caller information
	Enabled bool `json:"enabled,omitempty" hcl:"enabled,optional"`

	// Field key for caller information
	Key string `json:"key,omitempty" hcl:"key,optional"`

	// Whether to include package information
	IncludePackage bool `json:"include_package,omitempty" hcl:"include_package,optional"`

	// Number of stack frames to skip
	SkipFrames int `json:"skip_frames,omitempty" hcl:"skip_frames,optional"`
}

// LogFieldsConfig represents fields configuration
type LogFieldsConfig struct {
	// Global fields to include in all log entries
	Global map[string]interface{} `json:"global,omitempty" hcl:"global,optional"`

	// Order of fields in log output
	Order []string `json:"order,omitempty" hcl:"order,optional"`

	// Field filtering configuration
	Filter *LogFieldFilterConfig `json:"filter,omitempty" hcl:"filter,optional"`
}

// LogFieldFilterConfig represents field filtering configuration
type LogFieldFilterConfig struct {
	// Field names to exclude from output
	Exclude []string `json:"exclude,omitempty" hcl:"exclude,optional"`

	// Field names to include in output
	Include []string `json:"include,omitempty" hcl:"include,optional"`

	// Sensitive field names to mask or redact
	Sensitive []string `json:"sensitive,omitempty" hcl:"sensitive,optional"`
}

// LogPerformanceConfig represents performance configuration
type LogPerformanceConfig struct {
	// Buffering configuration
	Buffer *LogBufferConfig `json:"buffer,omitempty" hcl:"buffer,optional"`

	// Async configuration
	Async *LogAsyncConfig `json:"async,omitempty" hcl:"async,optional"`
}

// LogBufferConfig represents buffering configuration
type LogBufferConfig struct {
	// Whether to use buffered logging
	Enabled bool `json:"enabled,omitempty" hcl:"enabled,optional"`

	// Buffer size in bytes
	Size int `json:"size,omitempty" hcl:"size,optional"`

	// Flush interval
	FlushInterval string `json:"flush_interval,omitempty" hcl:"flush_interval,optional"`
}

// LogAsyncConfig represents async configuration
type LogAsyncConfig struct {
	// Whether to use asynchronous logging
	Enabled bool `json:"enabled,omitempty" hcl:"enabled,optional"`

	// Queue size for async logging
	QueueSize int `json:"queue_size,omitempty" hcl:"queue_size,optional"`

	// Number of worker goroutines
	Workers int `json:"workers,omitempty" hcl:"workers,optional"`

	// Whether to drop logs when queue is full
	DropWhenFull bool `json:"drop_when_full,omitempty" hcl:"drop_when_full,optional"`
}

// LogFilteringConfig represents filtering configuration
type LogFilteringConfig struct {
	// Component-specific log level configuration
	Components map[string]LogLevel `json:"components,omitempty" hcl:"components,optional"`

	// Pattern-based filtering configuration
	Patterns *LogPatternConfig `json:"patterns,omitempty" hcl:"patterns,optional"`
}

// LogPatternConfig represents pattern-based filtering configuration
type LogPatternConfig struct {
	// Regex patterns to include
	Include []string `json:"include,omitempty" hcl:"include,optional"`

	// Regex patterns to exclude
	Exclude []string `json:"exclude,omitempty" hcl:"exclude,optional"`
}

// LogRotationConfig represents rotation configuration
type LogRotationConfig struct {
	// Whether to enable log file rotation
	Enabled bool `json:"enabled,omitempty" hcl:"enabled,optional"`

	// Maximum file size before rotation
	MaxSize string `json:"max_size,omitempty" hcl:"max_size,optional"`

	// Maximum age of rotated files
	MaxAge string `json:"max_age,omitempty" hcl:"max_age,optional"`

	// Maximum number of backup files to keep
	MaxBackups int `json:"max_backups,omitempty" hcl:"max_backups,optional"`

	// Whether to compress rotated log files
	Compress bool `json:"compress,omitempty" hcl:"compress,optional"`

	// Whether to use local time for rotation timestamps
	LocalTime bool `json:"local_time,omitempty" hcl:"local_time,optional"`
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
