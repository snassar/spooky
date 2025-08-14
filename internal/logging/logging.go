// Package logging provides logging functionality for the spooky codebase.
// This package implements logging operations using Go's slog package with schema-driven configuration.
package logging

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"time"

	spookytypeslogging "spooky/internal/types/logging"
)

// LogManager implements the LogManager interface using slog
type LogManager struct {
	// Configuration
	config *spookytypeslogging.LogConfig

	// slog logger
	logger *slog.Logger

	// Output writer
	writer io.Writer

	// Loggers by component
	loggers map[string]*Logger

	// Mutex for thread safety
	mutex sync.RWMutex

	// Output file (if configured)
	outputFile *os.File

	// Context for cancellation
	ctx    context.Context
	cancel context.CancelFunc
}

// NewLogManager creates a new log manager with default configuration
func NewLogManager() *LogManager {
	ctx, cancel := context.WithCancel(context.Background())

	// Default configuration
	config := &spookytypeslogging.LogConfig{
		Level:  spookytypeslogging.LogLevelInfo,
		Format: "json",
		Output: "stderr",
	}

	// Create default slog logger
	logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	return &LogManager{
		config:  config,
		logger:  logger,
		writer:  os.Stderr,
		loggers: make(map[string]*Logger),
		ctx:     ctx,
		cancel:  cancel,
	}
}

// GetLogger returns a logger for the given component
func (lm *LogManager) GetLogger(component string) spookytypeslogging.Logger {
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
		slogger:   lm.logger.With("component", component),
	}

	lm.loggers[component] = logger
	return logger
}

// SetLevel sets the log level for all loggers
func (lm *LogManager) SetLevel(level spookytypeslogging.LogLevel) {
	lm.mutex.Lock()
	defer lm.mutex.Unlock()

	lm.config.Level = level

	// Update slog logger level
	lm.updateSlogLogger()

	// Update all existing loggers
	for _, logger := range lm.loggers {
		logger.level = level
	}
}

// GetLevel returns the current log level
func (lm *LogManager) GetLevel() spookytypeslogging.LogLevel {
	lm.mutex.RLock()
	defer lm.mutex.RUnlock()
	return lm.config.Level
}

// Configure configures logging with the given configuration
func (lm *LogManager) Configure(config *spookytypeslogging.LogConfig) error {
	lm.mutex.Lock()
	defer lm.mutex.Unlock()

	// Close existing output file if different
	if lm.outputFile != nil && lm.config.Output != config.Output {
		lm.outputFile.Close()
		lm.outputFile = nil
	}

	// Configure output
	if err := lm.configureOutput(config); err != nil {
		return fmt.Errorf("failed to configure output: %w", err)
	}

	lm.config = config

	// Update slog logger
	lm.updateSlogLogger()

	// Update all existing loggers
	for _, logger := range lm.loggers {
		logger.level = config.Level
		logger.slogger = lm.logger.With("component", logger.component)
	}

	return nil
}

// configureOutput configures the output destination
func (lm *LogManager) configureOutput(config *spookytypeslogging.LogConfig) error {
	switch config.Output {
	case "stdout":
		lm.writer = os.Stdout
	case "stderr":
		lm.writer = os.Stderr
	case "file":
		if config.File == nil || config.File.Path == "" {
			return fmt.Errorf("file output requires file.path configuration")
		}

		// Ensure directory exists
		dir := filepath.Dir(config.File.Path)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create log directory: %w", err)
		}

		// Determine file mode
		flag := os.O_CREATE | os.O_WRONLY
		if config.File.Append {
			flag |= os.O_APPEND
		} else {
			flag |= os.O_TRUNC
		}

		// Parse permissions
		perm := os.FileMode(0644)
		if config.File.Permissions != "" {
			if parsed, err := parseOctalPermissions(config.File.Permissions); err == nil {
				perm = parsed
			}
		}

		file, err := os.OpenFile(config.File.Path, flag, perm)
		if err != nil {
			return fmt.Errorf("failed to open log file: %w", err)
		}

		lm.outputFile = file
		lm.writer = file

	case "null":
		lm.writer = io.Discard
	default:
		return fmt.Errorf("unsupported output destination: %s", config.Output)
	}

	return nil
}

// updateSlogLogger updates the slog logger with current configuration
func (lm *LogManager) updateSlogLogger() {
	// Create handler options
	opts := &slog.HandlerOptions{
		Level: slog.Level(lm.config.Level.ToSlogLevel()),
	}

	// Create handler based on format
	var handler slog.Handler
	switch lm.config.Format {
	case "json":
		handler = slog.NewJSONHandler(lm.writer, opts)
	case "text":
		handler = slog.NewTextHandler(lm.writer, opts)
	case "structured":
		handler = lm.createStructuredHandler(opts)
	default:
		handler = slog.NewJSONHandler(lm.writer, opts)
	}

	lm.logger = slog.New(handler)
}

// createStructuredHandler creates a custom structured handler
func (lm *LogManager) createStructuredHandler(opts *slog.HandlerOptions) slog.Handler {
	if lm.config.Structured == nil {
		return slog.NewJSONHandler(lm.writer, opts)
	}

	// Create custom handler with structured configuration
	return &StructuredHandler{
		writer:   lm.writer,
		level:    slog.Level(lm.config.Level.ToSlogLevel()),
		config:   lm.config.Structured,
		replacer: lm.createReplacer(),
	}
}

// createReplacer creates a replacer for sensitive fields
func (lm *LogManager) createReplacer() *strings.Replacer {
	if lm.config.Structured == nil || lm.config.Structured.Fields == nil || lm.config.Structured.Fields.Filter == nil {
		return nil
	}

	var pairs []string
	for _, field := range lm.config.Structured.Fields.Filter.Sensitive {
		pairs = append(pairs, field, "[REDACTED]")
	}

	if len(pairs) > 0 {
		return strings.NewReplacer(pairs...)
	}

	return nil
}

// Flush flushes all pending log entries
func (lm *LogManager) Flush() error {
	lm.mutex.Lock()
	defer lm.mutex.Unlock()

	// Flush any buffered writers
	if lm.outputFile != nil {
		// Sync the file to ensure all data is written to disk
		if err := lm.outputFile.Sync(); err != nil {
			return fmt.Errorf("failed to sync output file: %w", err)
		}
	}

	// Flush any buffered loggers
	for _, logger := range lm.loggers {
		// For now, individual loggers don't need flushing
		// The slog logger handles its own buffering
		_ = logger // Use logger to avoid unused variable warning
	}

	return nil
}

// Close closes the log manager
func (lm *LogManager) Close() error {
	lm.mutex.Lock()
	defer lm.mutex.Unlock()

	lm.cancel()

	if lm.outputFile != nil {
		return lm.outputFile.Close()
	}

	return nil
}

// Logger implements the Logger interface using slog
type Logger struct {
	// Component name
	component string

	// Log manager reference
	manager *LogManager

	// Log level
	level spookytypeslogging.LogLevel

	// slog logger
	slogger *slog.Logger

	// Additional fields
	fields map[string]interface{}

	// Operation context
	operation string
}

// Debug logs a debug message
func (l *Logger) Debug(msg string, fields ...map[string]interface{}) {
	if l.shouldLog(spookytypeslogging.LogLevelDebug) {
		l.log(slog.LevelDebug, msg, fields...)
	}
}

// Info logs an info message
func (l *Logger) Info(msg string, fields ...map[string]interface{}) {
	if l.shouldLog(spookytypeslogging.LogLevelInfo) {
		l.log(slog.LevelInfo, msg, fields...)
	}
}

// Warn logs a warning message
func (l *Logger) Warn(msg string, fields ...map[string]interface{}) {
	if l.shouldLog(spookytypeslogging.LogLevelWarn) {
		l.log(slog.LevelWarn, msg, fields...)
	}
}

// Error logs an error message
func (l *Logger) Error(msg string, err error, fields ...map[string]interface{}) {
	if l.shouldLog(spookytypeslogging.LogLevelError) {
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

		l.log(slog.LevelError, msg, errorFields)
	}
}

// Fatal logs a fatal message and exits
func (l *Logger) Fatal(msg string, err error, fields ...map[string]interface{}) {
	l.Error(msg, err, fields...)
	os.Exit(1)
}

// WithFields returns a logger with additional fields
func (l *Logger) WithFields(fields map[string]interface{}) spookytypeslogging.Logger {
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

	// Update slogger with fields
	attrs := make([]any, 0, len(newLogger.fields)*2)
	for k, v := range newLogger.fields {
		attrs = append(attrs, k, v)
	}
	newLogger.slogger = l.slogger.With(attrs...)

	return &newLogger
}

// WithComponent returns a logger with a component name
func (l *Logger) WithComponent(component string) spookytypeslogging.Logger {
	newLogger := *l
	newLogger.component = component
	newLogger.slogger = l.slogger.With("component", component)
	return &newLogger
}

// WithOperation returns a logger with an operation name
func (l *Logger) WithOperation(operation string) spookytypeslogging.Logger {
	newLogger := *l
	newLogger.operation = operation
	newLogger.slogger = l.slogger.With("operation", operation)
	return &newLogger
}

// SetLevel sets the log level
func (l *Logger) SetLevel(level spookytypeslogging.LogLevel) {
	l.level = level
}

// GetLevel returns the current log level
func (l *Logger) GetLevel() spookytypeslogging.LogLevel {
	return l.level
}

// shouldLog determines if a message should be logged at the given level
func (l *Logger) shouldLog(level spookytypeslogging.LogLevel) bool {
	levelOrder := map[spookytypeslogging.LogLevel]int{
		spookytypeslogging.LogLevelDebug: 0,
		spookytypeslogging.LogLevelInfo:  1,
		spookytypeslogging.LogLevelWarn:  2,
		spookytypeslogging.LogLevelError: 3,
		spookytypeslogging.LogLevelFatal: 4,
	}

	return levelOrder[level] >= levelOrder[l.level]
}

// log logs a message at the given level
func (l *Logger) log(level slog.Level, msg string, fields ...map[string]interface{}) {
	// Check if we should filter this log
	if l.shouldFilter(msg, fields...) {
		return
	}

	// Build attributes
	attrs := make([]any, 0, len(l.fields)*2+len(fields)*2)

	// Add logger fields
	for k, v := range l.fields {
		attrs = append(attrs, k, v)
	}

	// Add provided fields
	if len(fields) > 0 {
		for k, v := range fields[0] {
			attrs = append(attrs, k, v)
		}
	}

	// Add caller information if configured
	if l.manager.config.Structured != nil && l.manager.config.Structured.Caller != nil && l.manager.config.Structured.Caller.Enabled {
		if caller := l.getCaller(); caller != nil {
			attrs = append(attrs, "caller", caller)
		}
	}

	// Log using slog
	l.slogger.Log(context.Background(), level, msg, attrs...)
}

// shouldFilter determines if this log should be filtered
func (l *Logger) shouldFilter(msg string, fields ...map[string]interface{}) bool {
	if l.manager.config.Filtering == nil {
		return false
	}

	// Check component-specific filtering
	if l.manager.config.Filtering.Components != nil {
		if componentLevel, exists := l.manager.config.Filtering.Components[l.component]; exists {
			if !l.shouldLog(componentLevel) {
				return true
			}
		}
	}

	// Check pattern-based filtering
	if l.manager.config.Filtering.Patterns != nil {
		// Check include patterns
		if len(l.manager.config.Filtering.Patterns.Include) > 0 {
			matched := false
			for _, pattern := range l.manager.config.Filtering.Patterns.Include {
				if matched, _ := regexp.MatchString(pattern, msg); matched {
					matched = true
					break
				}
			}
			if !matched {
				return true
			}
		}

		// Check exclude patterns
		for _, pattern := range l.manager.config.Filtering.Patterns.Exclude {
			if matched, _ := regexp.MatchString(pattern, msg); matched {
				return true
			}
		}
	}

	return false
}

// getCaller gets caller information
func (l *Logger) getCaller() *spookytypeslogging.LogCaller {
	skip := 3 // Default skip frames
	if l.manager.config.Structured != nil && l.manager.config.Structured.Caller != nil {
		skip = l.manager.config.Structured.Caller.SkipFrames
	}

	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return nil
	}

	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return nil
	}

	caller := &spookytypeslogging.LogCaller{
		File:     filepath.Base(file),
		Line:     line,
		Function: fn.Name(),
	}

	// Include package information if configured
	if l.manager.config.Structured != nil && l.manager.config.Structured.Caller != nil && l.manager.config.Structured.Caller.IncludePackage {
		parts := strings.Split(fn.Name(), ".")
		if len(parts) > 1 {
			caller.Package = parts[0]
		}
	}

	return caller
}

// StructuredHandler implements a custom slog handler for structured logging
type StructuredHandler struct {
	writer   io.Writer
	level    slog.Level
	config   *spookytypeslogging.LogStructuredConfig
	replacer *strings.Replacer
}

// Enabled reports whether the handler handles records at the given level
func (h *StructuredHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= h.level
}

// Handle handles the Record
func (h *StructuredHandler) Handle(ctx context.Context, r slog.Record) error {
	// Build structured log entry
	entry := h.buildLogEntry(r)

	// Format and write
	formatted := h.formatEntry(entry)
	_, err := h.writer.Write([]byte(formatted + "\n"))
	return err
}

// WithAttrs returns a new handler with the given attributes
func (h *StructuredHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	// For simplicity, return the same handler
	// In a full implementation, you'd want to store the attributes
	return h
}

// WithGroup returns a new handler with the given group
func (h *StructuredHandler) WithGroup(name string) slog.Handler {
	// For simplicity, return the same handler
	// In a full implementation, you'd want to handle groups
	return h
}

// buildLogEntry builds a log entry from a slog record
func (h *StructuredHandler) buildLogEntry(r slog.Record) *spookytypeslogging.LogEntry {
	entry := &spookytypeslogging.LogEntry{
		Timestamp: r.Time,
		Level:     spookytypeslogging.FromSlogLevel(r.Level),
		Message:   r.Message,
		Fields:    make(map[string]interface{}),
	}

	// Add attributes as fields
	r.Attrs(func(attr slog.Attr) bool {
		entry.Fields[attr.Key] = attr.Value.Any()
		return true
	})

	return entry
}

// formatEntry formats a log entry according to configuration
func (h *StructuredHandler) formatEntry(entry *spookytypeslogging.LogEntry) string {
	// Use default JSON format since format is not available in LogStructuredConfig
	// The format is typically set in the main LogConfig
	format := "json" // default format

	switch format {
	case "json":
		return h.formatJSON(entry)
	case "text":
		return h.formatText(entry)
	case "logfmt":
		return h.formatLogfmt(entry)
	default:
		return h.formatJSON(entry)
	}
}

// formatJSON formats entry as JSON
func (h *StructuredHandler) formatJSON(entry *spookytypeslogging.LogEntry) string {
	// Simple JSON formatting
	fieldsJSON := "{}"
	if len(entry.Fields) > 0 {
		// Convert fields to JSON string (simplified)
		var pairs []string
		for k, v := range entry.Fields {
			pairs = append(pairs, fmt.Sprintf(`"%s":"%v"`, k, v))
		}
		fieldsJSON = "{" + strings.Join(pairs, ",") + "}"
	}

	return fmt.Sprintf(`{"timestamp":"%s","level":"%s","message":"%s","fields":%s}`,
		entry.Timestamp.Format(time.RFC3339),
		entry.Level,
		entry.Message,
		fieldsJSON)
}

// formatText formats entry as human-readable text
func (h *StructuredHandler) formatText(entry *spookytypeslogging.LogEntry) string {
	fields := ""
	if len(entry.Fields) > 0 {
		var pairs []string
		for k, v := range entry.Fields {
			pairs = append(pairs, fmt.Sprintf("%s=%v", k, v))
		}
		fields = " " + strings.Join(pairs, " ")
	}

	return fmt.Sprintf("%s [%s] %s%s",
		entry.Timestamp.Format("2006-01-02T15:04:05Z07:00"),
		strings.ToUpper(string(entry.Level)),
		entry.Message,
		fields)
}

// formatLogfmt formats entry in logfmt format
func (h *StructuredHandler) formatLogfmt(entry *spookytypeslogging.LogEntry) string {
	var pairs []string

	// Add timestamp
	pairs = append(pairs, fmt.Sprintf("time=%s", entry.Timestamp.Format(time.RFC3339)))

	// Add level
	pairs = append(pairs, fmt.Sprintf("level=%s", strings.ToLower(string(entry.Level))))

	// Add message
	pairs = append(pairs, fmt.Sprintf("msg=%q", entry.Message))

	// Add fields
	for k, v := range entry.Fields {
		pairs = append(pairs, fmt.Sprintf("%s=%v", k, v))
	}

	return strings.Join(pairs, " ")
}

// parseOctalPermissions parses octal permissions string
func parseOctalPermissions(perm string) (os.FileMode, error) {
	if len(perm) != 3 && len(perm) != 4 {
		return 0, fmt.Errorf("invalid permission format: %s", perm)
	}

	var mode os.FileMode
	for _, c := range perm {
		if c < '0' || c > '7' {
			return 0, fmt.Errorf("invalid permission character: %c", c)
		}
		mode = mode<<3 | os.FileMode(c-'0')
	}

	return mode, nil
}
