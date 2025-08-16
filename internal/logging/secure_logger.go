// Package logging provides secure logging functionality with redaction patterns
package logging

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	spookytypeslogging "spooky/internal/types/logging"
)

// RedactionPatterns defines patterns for redacting sensitive data
type RedactionPatterns struct {
	// Age strings are safe to log (they're encrypted)
	AgeStringPattern *regexp.Regexp

	// Decrypted values to redact
	DecryptedValuePattern *regexp.Regexp

	// Sensitive field names to redact
	SensitiveFieldPattern *regexp.Regexp

	// Object/map values to redact
	ObjectValuePattern *regexp.Regexp
}

// NewRedactionPatterns creates new redaction patterns
func NewRedactionPatterns() *RedactionPatterns {
	return &RedactionPatterns{
		AgeStringPattern:      regexp.MustCompile(`^age1[a-zA-Z0-9]+`),
		DecryptedValuePattern: regexp.MustCompile(`(?i)(password|secret|key|token|credential)`),
		SensitiveFieldPattern: regexp.MustCompile(`(?i)(password|secret|key|token|credential|private_key|ssh_key|auth_key)`),
		ObjectValuePattern:    regexp.MustCompile(`\{.*\}`), // Simple object detection
	}
}

// SecureLogger wraps a logger with redaction capabilities
type SecureLogger struct {
	logger     spookytypeslogging.Logger
	patterns   *RedactionPatterns
	redactMode bool
}

// NewSecureLogger creates a new secure logger
func NewSecureLogger(logger spookytypeslogging.Logger, redactMode bool) *SecureLogger {
	return &SecureLogger{
		logger:     logger,
		patterns:   NewRedactionPatterns(),
		redactMode: redactMode,
	}
}

// Info logs an info message with redaction
func (l *SecureLogger) Info(msg string, fields ...map[string]interface{}) {
	l.logger.Info(msg, l.sanitizeFields(fields)...)
}

// Debug logs a debug message with redaction
func (l *SecureLogger) Debug(msg string, fields ...map[string]interface{}) {
	l.logger.Debug(msg, l.sanitizeFields(fields)...)
}

// Warn logs a warning message with redaction
func (l *SecureLogger) Warn(msg string, fields ...map[string]interface{}) {
	l.logger.Warn(msg, l.sanitizeFields(fields)...)
}

// Error logs an error message with redaction
func (l *SecureLogger) Error(msg string, err error, fields ...map[string]interface{}) {
	l.logger.Error(msg, err, l.sanitizeFields(fields)...)
}

// Fatal logs a fatal message and exits
func (l *SecureLogger) Fatal(msg string, err error, fields ...map[string]interface{}) {
	l.logger.Fatal(msg, err, l.sanitizeFields(fields)...)
}

// WithFields returns a logger with additional fields
func (l *SecureLogger) WithFields(fields map[string]interface{}) spookytypeslogging.Logger {
	return NewSecureLogger(l.logger.WithFields(fields), l.redactMode)
}

// WithComponent returns a logger with a component name
func (l *SecureLogger) WithComponent(component string) spookytypeslogging.Logger {
	return NewSecureLogger(l.logger.WithComponent(component), l.redactMode)
}

// WithOperation returns a logger with an operation name
func (l *SecureLogger) WithOperation(operation string) spookytypeslogging.Logger {
	return NewSecureLogger(l.logger.WithOperation(operation), l.redactMode)
}

// SetLevel sets the log level
func (l *SecureLogger) SetLevel(level spookytypeslogging.LogLevel) {
	l.logger.SetLevel(level)
}

// GetLevel returns the current log level
func (l *SecureLogger) GetLevel() spookytypeslogging.LogLevel {
	return l.logger.GetLevel()
}

// sanitizeFields sanitizes log fields to prevent sensitive data leakage
func (l *SecureLogger) sanitizeFields(fields []map[string]interface{}) []map[string]interface{} {
	sanitized := make([]map[string]interface{}, len(fields))

	for i, field := range fields {
		sanitized[i] = l.sanitizeFieldMap(field)
	}

	return sanitized
}

// sanitizeFieldMap sanitizes a field map
func (l *SecureLogger) sanitizeFieldMap(field map[string]interface{}) map[string]interface{} {
	sanitized := make(map[string]interface{})

	for key, value := range field {
		if l.shouldRedactKey(key) {
			sanitized[key] = "[REDACTED]"
		} else {
			sanitized[key] = l.sanitizeValue(value)
		}
	}

	return sanitized
}

// shouldRedactKey determines if a key should be redacted
func (l *SecureLogger) shouldRedactKey(key string) bool {
	// Redact sensitive field names
	return l.patterns.SensitiveFieldPattern.MatchString(key)
}

// sanitizeValue sanitizes a field value
func (l *SecureLogger) sanitizeValue(value interface{}) interface{} {
	if !l.redactMode {
		return value // No redaction needed
	}

	switch v := value.(type) {
	case string:
		if l.patterns.AgeStringPattern.MatchString(v) {
			return v // Age strings are safe
		}
		return "[REDACTED_VALUE]"

	case map[string]interface{}:
		return "[REDACTED_OBJECT]"

	case []interface{}:
		return "[REDACTED_ARRAY]"

	default:
		return value
	}
}

// PostLogScanner scans log files for potential secret leaks
type PostLogScanner struct {
	patterns *RedactionPatterns
	logFile  string
}

// NewPostLogScanner creates a new post-log scanner
func NewPostLogScanner(logFile string) *PostLogScanner {
	return &PostLogScanner{
		patterns: NewRedactionPatterns(),
		logFile:  logFile,
	}
}

// ScanForLeaks scans the log file for potential secret leaks
func (s *PostLogScanner) ScanForLeaks() error {
	// Read the log file
	data, err := os.ReadFile(s.logFile)
	if err != nil {
		return fmt.Errorf("failed to read log file: %w", err)
	}

	lines := strings.Split(string(data), "\n")
	var leaks []string

	for lineNum, line := range lines {
		lineNum++ // Convert to 1-based line numbers

		if s.containsLeakedSecret(line) {
			leaks = append(leaks, fmt.Sprintf("Line %d: %s", lineNum, strings.TrimSpace(line)))
		}
	}

	if len(leaks) > 0 {
		return fmt.Errorf("potential secret leaks detected in log file:\n%s", strings.Join(leaks, "\n"))
	}

	return nil
}

// containsLeakedSecret checks if a line contains leaked secrets
func (s *PostLogScanner) containsLeakedSecret(line string) bool {
	// Look for patterns that suggest decrypted secrets
	if s.patterns.DecryptedValuePattern.MatchString(line) {
		// Check if it's not an age string
		if !s.patterns.AgeStringPattern.MatchString(line) {
			return true
		}
	}

	return false
}

// SetupSecureLogging sets up secure logging with redaction
func SetupSecureLogging(logger spookytypeslogging.Logger, redactMode bool) spookytypeslogging.Logger {
	if redactMode {
		return NewSecureLogger(logger, true)
	}
	return logger
}
