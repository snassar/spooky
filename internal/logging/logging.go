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

	spookyschemas "spooky/internal/schemas"
	spookytypeslogging "spooky/internal/types/logging"
	spookytypesschemas "spooky/internal/types/schemas"
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

	// Schema validators
	schemaDrivenValidator *spookyschemas.SchemaDrivenValidator
	enhancedValidator     *spookyschemas.EnhancedValidator

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

	// Default configuration - for CLI, don't output logs to terminal by default
	config := &spookytypeslogging.LogConfig{
		Level:  spookytypeslogging.LogLevelError, // Only show errors by default
		Format: "json",
		Output: "null", // Don't output logs to terminal by default
	}

	// Create default slog logger that discards logs by default
	logger := slog.New(slog.NewJSONHandler(io.Discard, &slog.HandlerOptions{
		Level: slog.LevelError, // Only show errors by default
	}))

	// Create schema validators for logging configuration validation
	schemaDrivenConfig := &spookyschemas.SchemaDrivenValidationConfig{
		UseEmbeddedSchemas: true,
		StrictValidation:   true,
		AllowUnknownFields: false,
		DetailedErrors:     true,
	}
	schemaDrivenValidator := spookyschemas.NewSchemaDrivenValidator(nil, schemaDrivenConfig)

	enhancedConfig := &spookyschemas.ValidationConfig{
		Mode: spookyschemas.ValidationModeStrict,
		ErrorHandling: &spookyschemas.ErrorHandlingConfig{
			StopOnFirstError:   false,
			MaxErrors:          100,
			IncludeWarnings:    true,
			IncludeContext:     true,
			IncludeSuggestions: true,
		},
		Evolution: &spookyschemas.EvolutionConfig{
			EnableTracking:  true,
			AllowDeprecated: true,
			WarnDeprecated:  true,
			AllowBreaking:   false,
		},
	}
	enhancedValidator := spookyschemas.NewEnhancedValidator(enhancedConfig)

	return &LogManager{
		config:                config,
		logger:                logger,
		writer:                io.Discard, // Don't output logs to terminal by default
		loggers:               make(map[string]*Logger),
		schemaDrivenValidator: schemaDrivenValidator,
		enhancedValidator:     enhancedValidator,
		ctx:                   ctx,
		cancel:                cancel,
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
		return lm.configureFileOutput(config)
	case "null":
		lm.writer = io.Discard
	default:
		return fmt.Errorf("unsupported output destination: %s", config.Output)
	}

	return nil
}

// configureFileOutput configures file output destination
func (lm *LogManager) configureFileOutput(config *spookytypeslogging.LogConfig) error {
	if config.File == nil || config.File.Path == "" {
		return fmt.Errorf("file output requires file.path configuration")
	}

	if err := lm.ensureLogDirectory(config.File.Path); err != nil {
		return err
	}

	file, err := lm.openLogFile(config.File)
	if err != nil {
		return err
	}

	lm.outputFile = file
	lm.writer = file
	return nil
}

// ensureLogDirectory ensures the log directory exists
func (lm *LogManager) ensureLogDirectory(filePath string) error {
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("failed to create log directory: %w", err)
	}
	return nil
}

// openLogFile opens the log file with appropriate flags and permissions
func (lm *LogManager) openLogFile(fileConfig *spookytypeslogging.LogFileConfig) (*os.File, error) {
	flag := lm.determineFileFlags(fileConfig)
	perm := lm.determineFilePermissions(fileConfig)

	file, err := os.OpenFile(fileConfig.Path, flag, perm)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	return file, nil
}

// determineFileFlags determines the file open flags
func (lm *LogManager) determineFileFlags(fileConfig *spookytypeslogging.LogFileConfig) int {
	flag := os.O_CREATE | os.O_WRONLY
	if fileConfig.Append {
		flag |= os.O_APPEND
	} else {
		flag |= os.O_TRUNC
	}
	return flag
}

// determineFilePermissions determines the file permissions
func (lm *LogManager) determineFilePermissions(fileConfig *spookytypeslogging.LogFileConfig) os.FileMode {
	perm := os.FileMode(0o644)
	if fileConfig.Permissions != "" {
		if parsed, err := parseOctalPermissions(fileConfig.Permissions); err == nil {
			perm = parsed
		}
	}
	return perm
}

// updateSlogLogger updates the slog logger with current configuration
func (lm *LogManager) updateSlogLogger() {
	// Create handler options
	opts := &slog.HandlerOptions{
		Level: slog.Level(lm.config.Level.ToSlogLevel()),
	}

	// Create handler based on format
	var handler slog.Handler
	const jsonFormat = "json"
	switch lm.config.Format {
	case jsonFormat:
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

// ValidateLogConfig validates logging configuration using schema validators
func (lm *LogManager) ValidateLogConfig(ctx context.Context, config *spookytypeslogging.LogConfig) (*spookytypesschemas.ValidationResult, error) {
	// Get logging schema for enhanced validation
	loggingSchema, err := lm.getLoggingSchema()
	if err != nil {
		return nil, fmt.Errorf("failed to get logging schema: %w", err)
	}

	// Use enhanced validator for comprehensive logging configuration validation
	result, err := lm.enhancedValidator.ValidateWithEnhancedFeatures(ctx, loggingSchema, config)
	if err != nil {
		return nil, fmt.Errorf("failed to validate logging config with enhanced validator: %w", err)
	}

	// Add additional custom validation for logging-specific rules
	lm.addCustomLoggingValidation(config, result)

	return result, nil
}

// getLoggingSchema gets the logging schema for validation
func (lm *LogManager) getLoggingSchema() (*spookytypesschemas.Schema, error) {
	// Try to get schema from embedded schemas first
	if schema, err := lm.schemaDrivenValidator.GetEmbeddedSchema("logging"); err == nil {
		return schema, nil
	}

	// Fallback: create a basic logging schema
	return &spookytypesschemas.Schema{
		Name:        "logging",
		Type:        "hcl",
		Version:     "1.0",
		Description: "Logging configuration schema",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Content:     "", // Will be loaded from file if needed
		Metadata:    make(map[string]interface{}),
	}, nil
}

// addCustomLoggingValidation adds custom validation rules specific to logging
func (lm *LogManager) addCustomLoggingValidation(config *spookytypeslogging.LogConfig, result *spookytypesschemas.ValidationResult) {
	// Validate log level
	if config.Level == "" {
		lm.addSchemaError(result, "missing_log_level", "Log level is required", "error")
	}

	// Validate log format
	if config.Format == "" {
		lm.addSchemaError(result, "missing_log_format", "Log format is required", "error")
	}

	// Validate output configuration
	if config.Output == "" {
		lm.addSchemaError(result, "missing_log_output", "Log output is required", "error")
	}

	// Validate file path if output is file
	if config.Output == "file" && (config.File == nil || config.File.Path == "") {
		lm.addSchemaError(result, "missing_file_path", "File path is required when output is file", "error")
	}
}

// addSchemaError adds a schema error to the validation result
func (lm *LogManager) addSchemaError(result *spookytypesschemas.ValidationResult, code, message, severity string) {
	schemaError := spookytypesschemas.SchemaError{
		Code:     code,
		Message:  message,
		Severity: severity,
	}
	result.Errors = append(result.Errors, schemaError)
	result.Valid = false
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
func (l *Logger) shouldFilter(msg string, _ ...map[string]interface{}) bool {
	if l.manager.config.Filtering == nil {
		return false
	}

	if l.shouldFilterByComponent() {
		return true
	}

	if l.shouldFilterByPatterns(msg) {
		return true
	}

	return false
}

// shouldFilterByComponent checks component-specific filtering
func (l *Logger) shouldFilterByComponent() bool {
	if l.manager.config.Filtering.Components == nil {
		return false
	}

	componentLevel, exists := l.manager.config.Filtering.Components[l.component]
	if !exists {
		return false
	}

	return !l.shouldLog(componentLevel)
}

// shouldFilterByPatterns checks pattern-based filtering
func (l *Logger) shouldFilterByPatterns(msg string) bool {
	if l.manager.config.Filtering.Patterns == nil {
		return false
	}

	if l.shouldFilterByIncludePatterns(msg) {
		return true
	}

	return l.shouldFilterByExcludePatterns(msg)
}

// shouldFilterByIncludePatterns checks include pattern filtering
func (l *Logger) shouldFilterByIncludePatterns(msg string) bool {
	if len(l.manager.config.Filtering.Patterns.Include) == 0 {
		return false
	}

	for _, pattern := range l.manager.config.Filtering.Patterns.Include {
		if isMatch, _ := regexp.MatchString(pattern, msg); isMatch {
			return false // Don't filter if pattern matches
		}
	}

	return true // Filter if no patterns match
}

// shouldFilterByExcludePatterns checks exclude pattern filtering
func (l *Logger) shouldFilterByExcludePatterns(msg string) bool {
	for _, pattern := range l.manager.config.Filtering.Patterns.Exclude {
		if matched, _ := regexp.MatchString(pattern, msg); matched {
			return true // Filter if pattern matches
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
func (h *StructuredHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.level
}

// Handle handles the Record
//
//nolint:gocritic // slog.Handler interface requires slog.Record by value
func (h *StructuredHandler) Handle(_ context.Context, r slog.Record) error {
	// Build structured log entry
	entry := h.buildLogEntry(&r)

	// Format and write
	formatted := h.formatEntry(entry)
	_, err := h.writer.Write([]byte(formatted + "\n"))
	return err
}

// WithAttrs returns a new handler with the given attributes
func (h *StructuredHandler) WithAttrs(_ []slog.Attr) slog.Handler {
	// For simplicity, return the same handler
	// In a full implementation, you'd want to store the attributes
	return h
}

// WithGroup returns a new handler with the given group
func (h *StructuredHandler) WithGroup(_ string) slog.Handler {
	// For simplicity, return the same handler
	// In a full implementation, you'd want to handle groups
	return h
}

// buildLogEntry builds a log entry from a slog record
func (h *StructuredHandler) buildLogEntry(r *slog.Record) *spookytypeslogging.LogEntry {
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
	const jsonFormat = "json"
	format := jsonFormat // default format

	switch format {
	case jsonFormat:
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
			pairs = append(pairs, fmt.Sprintf(`%q:%v`, k, v))
		}
		fieldsJSON = "{" + strings.Join(pairs, ",") + "}"
	}

	return fmt.Sprintf(`{"timestamp":%q,"level":%q,"message":%q,"fields":%s}`,
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

	// Add timestamp, level, and message
	pairs = append(pairs,
		fmt.Sprintf("time=%s", entry.Timestamp.Format(time.RFC3339)),
		fmt.Sprintf("level=%s", strings.ToLower(string(entry.Level))),
		fmt.Sprintf("msg=%q", entry.Message),
	)

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
