// Package logging provides logging functionality for the spooky codebase.
// This package implements logging operations and log management.
package logging

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	spookytypes "spooky/internal/types"
)

// LogManager implements the LogManager interface
type LogManager struct {
	// Configuration
	config *spookytypes.LogConfig

	// Loggers by component
	loggers map[string]spookytypes.Logger

	// Mutex for thread safety
	mutex sync.RWMutex

	// Output file (if configured)
	outputFile *os.File
}

// NewLogManager creates a new log manager
func NewLogManager() spookyinterfaces.LogManager {
	return &LogManager{
		loggers: make(map[string]spookytypes.Logger),
		config: &spookytypes.LogConfig{
			Level:            spookytypes.LogLevelInfo,
			Format:           "text",
			Output:           "stdout",
			IncludeTimestamp: true,
		},
	}
}

// GetLogger returns a logger for the given component
func (lm *LogManager) GetLogger(component string) spookytypes.Logger {
	lm.mutex.RLock()
	if logger, exists := lm.loggers[component]; exists {
		lm.mutex.RUnlock()
		return logger
	}
	lm.mutex.RUnlock()

	// Create new logger
	lm.mutex.Lock()
	defer lm.mutex.Unlock()

	// Double-check after acquiring write lock
	if logger, exists := lm.loggers[component]; exists {
		return logger
	}

	logger := &Logger{
		component: component,
		manager:   lm,
		level:     lm.config.Level,
	}

	lm.loggers[component] = logger
	return logger
}

// SetLevel sets the log level for all loggers
func (lm *LogManager) SetLevel(level spookytypes.LogLevel) {
	lm.mutex.Lock()
	defer lm.mutex.Unlock()

	lm.config.Level = level

	// Update all existing loggers
	for _, logger := range lm.loggers {
		if l, ok := logger.(*Logger); ok {
			l.level = level
		}
	}
}

// GetLevel returns the current log level
func (lm *LogManager) GetLevel() spookytypes.LogLevel {
	lm.mutex.RLock()
	defer lm.mutex.RUnlock()
	return lm.config.Level
}

// Configure configures logging with the given configuration
func (lm *LogManager) Configure(config *spookytypes.LogConfig) error {
	lm.mutex.Lock()
	defer lm.mutex.Unlock()

	// Close existing output file if different
	if lm.outputFile != nil && lm.config.Output != config.Output {
		lm.outputFile.Close()
		lm.outputFile = nil
	}

	// Open new output file if needed
	if config.Output == "file" && config.FilePath != "" {
		file, err := os.OpenFile(config.FilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return fmt.Errorf("failed to open log file: %w", err)
		}
		lm.outputFile = file
	}

	lm.config = config

	// Update all existing loggers
	for _, logger := range lm.loggers {
		if l, ok := logger.(*Logger); ok {
			l.level = config.Level
		}
	}

	return nil
}

// Flush flushes all pending log entries
func (lm *LogManager) Flush() error {
	// For now, no buffering, so nothing to flush
	return nil
}

// Close closes the log manager
func (lm *LogManager) Close() error {
	lm.mutex.Lock()
	defer lm.mutex.Unlock()

	if lm.outputFile != nil {
		return lm.outputFile.Close()
	}

	return nil
}

// Logger implements the Logger interface
type Logger struct {
	// Component name
	component string

	// Log manager reference
	manager *LogManager

	// Log level
	level spookytypes.LogLevel

	// Additional fields
	fields map[string]interface{}
}

// Debug logs a debug message
func (l *Logger) Debug(msg string, fields ...map[string]interface{}) {
	if l.level == spookytypes.LogLevelDebug {
		l.log(spookytypes.LogLevelDebug, msg, fields...)
	}
}

// Info logs an info message
func (l *Logger) Info(msg string, fields ...map[string]interface{}) {
	if l.shouldLog(spookytypes.LogLevelInfo) {
		l.log(spookytypes.LogLevelInfo, msg, fields...)
	}
}

// Warn logs a warning message
func (l *Logger) Warn(msg string, fields ...map[string]interface{}) {
	if l.shouldLog(spookytypes.LogLevelWarn) {
		l.log(spookytypes.LogLevelWarn, msg, fields...)
	}
}

// Error logs an error message
func (l *Logger) Error(msg string, err error, fields ...map[string]interface{}) {
	if l.shouldLog(spookytypes.LogLevelError) {
		// Add error to fields
		errorFields := map[string]interface{}{
			"error": err.Error(),
		}

		// Merge with provided fields
		if len(fields) > 0 {
			for k, v := range fields[0] {
				errorFields[k] = v
			}
		}

		l.log(spookytypes.LogLevelError, msg, errorFields)
	}
}

// Fatal logs a fatal message and exits
func (l *Logger) Fatal(msg string, err error, fields ...map[string]interface{}) {
	l.Error(msg, err, fields...)
	os.Exit(1)
}

// WithFields returns a logger with additional fields
func (l *Logger) WithFields(fields map[string]interface{}) spookytypes.Logger {
	newLogger := *l
	newLogger.fields = make(map[string]interface{})

	// Copy existing fields
	for k, v := range l.fields {
		newLogger.fields[k] = v
	}

	// Add new fields
	for k, v := range fields {
		newLogger.fields[k] = v
	}

	return &newLogger
}

// WithComponent returns a logger with a component name
func (l *Logger) WithComponent(component string) spookytypes.Logger {
	newLogger := *l
	newLogger.component = component
	return &newLogger
}

// WithOperation returns a logger with an operation name
func (l *Logger) WithOperation(operation string) spookytypes.Logger {
	return l.WithFields(map[string]interface{}{
		"operation": operation,
	})
}

// SetLevel sets the log level
func (l *Logger) SetLevel(level spookytypes.LogLevel) {
	l.level = level
}

// GetLevel returns the current log level
func (l *Logger) GetLevel() spookytypes.LogLevel {
	return l.level
}

// shouldLog determines if a message should be logged at the given level
func (l *Logger) shouldLog(level spookytypes.LogLevel) bool {
	levelOrder := map[spookytypes.LogLevel]int{
		spookytypes.LogLevelDebug: 0,
		spookytypes.LogLevelInfo:  1,
		spookytypes.LogLevelWarn:  2,
		spookytypes.LogLevelError: 3,
		spookytypes.LogLevelFatal: 4,
	}

	return levelOrder[level] >= levelOrder[l.level]
}

// log logs a message at the given level
func (l *Logger) log(level spookytypes.LogLevel, msg string, fields ...map[string]interface{}) {
	// Build log entry
	entry := &spookytypes.LogEntry{
		Timestamp: time.Now(),
		Level:     level,
		Message:   msg,
		Component: l.component,
		Fields:    make(map[string]interface{}),
	}

	// Add logger fields
	for k, v := range l.fields {
		entry.Fields[k] = v
	}

	// Add provided fields
	if len(fields) > 0 {
		for k, v := range fields[0] {
			entry.Fields[k] = v
		}
	}

	// Format and output log entry
	l.outputLogEntry(entry)
}

// outputLogEntry outputs a log entry
func (l *Logger) outputLogEntry(entry *spookytypes.LogEntry) {
	// Simple text format for now
	timestamp := entry.Timestamp.Format("2006-01-02T15:04:05Z07:00")
	level := string(entry.Level)
	component := entry.Component
	message := entry.Message

	logLine := fmt.Sprintf("[%s] %s %s: %s", timestamp, level, component, message)

	// Add fields if present
	if len(entry.Fields) > 0 {
		logLine += fmt.Sprintf(" %+v", entry.Fields)
	}

	// Output to appropriate destination
	if l.manager.outputFile != nil {
		fmt.Fprintln(l.manager.outputFile, logLine)
	} else {
		log.Println(logLine)
	}
}
