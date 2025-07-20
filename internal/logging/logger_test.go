package logging

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEnsureLogDirectory(t *testing.T) {
	// Test with stdout (no directory needed)
	err := ensureLogDirectory("stdout")
	assert.NoError(t, err)

	// Test with stderr (no directory needed)
	err = ensureLogDirectory("stderr")
	assert.NoError(t, err)

	// Test with empty string (no directory needed)
	err = ensureLogDirectory("")
	assert.NoError(t, err)

	// Test with file path (should create directory)
	tempDir := t.TempDir()
	logFile := filepath.Join(tempDir, "logs", "test.log")

	err = ensureLogDirectory(logFile)
	assert.NoError(t, err)

	// Verify directory was created
	logDir := filepath.Dir(logFile)
	_, err = os.Stat(logDir)
	assert.NoError(t, err)
}

func TestConfigureLogger(t *testing.T) {
	// Test with debug level
	ConfigureLogger("debug", "json", "stdout", false, false)
	logger := GetLogger()
	assert.NotNil(t, logger)

	// Test with info level
	ConfigureLogger("info", "json", "stdout", false, false)
	logger = GetLogger()
	assert.NotNil(t, logger)

	// Test with warn level
	ConfigureLogger("warn", "json", "stdout", false, false)
	logger = GetLogger()
	assert.NotNil(t, logger)

	// Test with error level
	ConfigureLogger("error", "json", "stdout", false, false)
	logger = GetLogger()
	assert.NotNil(t, logger)

	// Test with invalid level (should default to info)
	ConfigureLogger("invalid", "json", "stdout", false, false)
	logger = GetLogger()
	assert.NotNil(t, logger)

	// Test with text format
	ConfigureLogger("info", "text", "stdout", false, false)
	logger = GetLogger()
	assert.NotNil(t, logger)

	// Test with quiet flag
	ConfigureLogger("info", "json", "stdout", true, false)
	logger = GetLogger()
	assert.NotNil(t, logger)

	// Test with verbose flag
	ConfigureLogger("info", "json", "stdout", false, true)
	logger = GetLogger()
	assert.NotNil(t, logger)

	// Test with file output
	tempDir := t.TempDir()
	logFile := filepath.Join(tempDir, "test.log")

	ConfigureLogger("info", "json", logFile, false, false)
	logger = GetLogger()
	assert.NotNil(t, logger)

	// Verify log file was created
	_, err := os.Stat(logFile)
	assert.NoError(t, err)
}

func TestConfigureLoggerWithInvalidDirectory(t *testing.T) {
	// Test with invalid directory (should fall back to stdout)
	ConfigureLogger("info", "json", "/invalid/path/test.log", false, false)
	logger := GetLogger()
	assert.NotNil(t, logger)
}

func TestNewLogger(t *testing.T) {
	// Test with debug level
	config := Config{
		Level:     DebugLevel,
		Format:    "json",
		Output:    "stdout",
		Timestamp: true,
	}
	logger := NewLogger(config)
	assert.NotNil(t, logger)

	// Test with info level
	config.Level = InfoLevel
	logger = NewLogger(config)
	assert.NotNil(t, logger)

	// Test with warn level
	config.Level = WarnLevel
	logger = NewLogger(config)
	assert.NotNil(t, logger)

	// Test with error level
	config.Level = ErrorLevel
	logger = NewLogger(config)
	assert.NotNil(t, logger)

	// Test with invalid level (should default to info)
	config.Level = LogLevel("invalid")
	logger = NewLogger(config)
	assert.NotNil(t, logger)

	// Test with text format
	config.Level = InfoLevel
	config.Format = "text"
	logger = NewLogger(config)
	assert.NotNil(t, logger)

	// Test with file output
	tempDir := t.TempDir()
	logFile := filepath.Join(tempDir, "test.log")
	config.Output = logFile
	config.Format = "json"

	logger = NewLogger(config)
	assert.NotNil(t, logger)

	// Verify log file was created
	_, err := os.Stat(logFile)
	assert.NoError(t, err)
}

func TestLoggerMethods(t *testing.T) {
	config := Config{
		Level:     DebugLevel,
		Format:    "json",
		Output:    "stdout",
		Timestamp: true,
	}
	logger := NewLogger(config)

	// Test Debug method
	logger.Debug("debug message", String("key", "value"))

	// Test Info method
	logger.Info("info message", String("key", "value"))

	// Test Warn method
	logger.Warn("warn message", String("key", "value"))

	// Test Error method
	logger.Error("error message", assert.AnError, String("key", "value"))

	// Test WithContext method
	ctx := context.Background()
	newLogger := logger.WithContext(ctx)
	assert.NotNil(t, newLogger)

	// Test WithFields method
	newLogger = logger.WithFields(String("field1", "value1"), Int("field2", 42))
	assert.NotNil(t, newLogger)

	// Note: Sync method is not part of the Logger interface
	// but is available on the underlying zapLogger implementation
}

func TestGetLogger(t *testing.T) {
	logger := GetLogger()
	assert.NotNil(t, logger)
}

func TestSetLogger(t *testing.T) {
	config := Config{
		Level:     InfoLevel,
		Format:    "json",
		Output:    "stdout",
		Timestamp: true,
	}
	newLogger := NewLogger(config)

	SetLogger(newLogger)

	logger := GetLogger()
	assert.Equal(t, newLogger, logger)
}

func TestFromContext(t *testing.T) {
	// Test with context that has logger
	config := Config{
		Level:     InfoLevel,
		Format:    "json",
		Output:    "stdout",
		Timestamp: true,
	}
	logger := NewLogger(config)

	ctx := WithContext(context.Background(), logger)
	retrievedLogger := FromContext(ctx)
	assert.Equal(t, logger, retrievedLogger)

	// Test with context that doesn't have logger
	emptyCtx := context.Background()
	retrievedLogger = FromContext(emptyCtx)
	assert.NotNil(t, retrievedLogger) // Should return default logger
}

func TestWithContext(t *testing.T) {
	config := Config{
		Level:     InfoLevel,
		Format:    "json",
		Output:    "stdout",
		Timestamp: true,
	}
	logger := NewLogger(config)

	ctx := context.Background()
	newCtx := WithContext(ctx, logger)

	assert.NotEqual(t, ctx, newCtx)

	// Verify logger is in context
	retrievedLogger := FromContext(newCtx)
	assert.Equal(t, logger, retrievedLogger)
}

func TestConvertFields(t *testing.T) {
	config := Config{
		Level:     InfoLevel,
		Format:    "json",
		Output:    "stdout",
		Timestamp: true,
	}
	logger := NewLogger(config)

	zapLogger, ok := logger.(*zapLogger)
	require.True(t, ok)

	// Test converting various field types
	fields := []Field{
		String("string_key", "string_value"),
		Int("int_key", 42),
		Float64("float_key", 3.14),
		Bool("bool_key", true),
	}

	zapFields := zapLogger.convertFields(fields)
	assert.Len(t, zapFields, 4)
}

func TestLoggerWithFileOutput(t *testing.T) {
	tempDir := t.TempDir()
	logFile := filepath.Join(tempDir, "test.log")

	config := Config{
		Level:     InfoLevel,
		Format:    "json",
		Output:    logFile,
		Timestamp: true,
	}
	logger := NewLogger(config)

	// Write some log messages
	logger.Info("test message 1", String("key1", "value1"))
	logger.Warn("test message 2", String("key2", "value2"))
	logger.Error("test message 3", assert.AnError, String("key3", "value3"))

	// Note: Sync method is not part of the Logger interface
	// but is available on the underlying zapLogger implementation

	// Verify log file contains messages
	content, err := os.ReadFile(logFile)
	assert.NoError(t, err)
	assert.NotEmpty(t, content)

	// Verify content contains our messages
	contentStr := string(content)
	assert.Contains(t, contentStr, "test message 1")
	assert.Contains(t, contentStr, "test message 2")
	assert.Contains(t, contentStr, "test message 3")
}

func TestLoggerWithTextFormat(t *testing.T) {
	tempDir := t.TempDir()
	logFile := filepath.Join(tempDir, "test.log")

	config := Config{
		Level:     InfoLevel,
		Format:    "text",
		Output:    logFile,
		Timestamp: true,
	}
	logger := NewLogger(config)

	// Write some log messages
	logger.Info("test message", String("key", "value"))

	// Note: Sync method is not part of the Logger interface
	// but is available on the underlying zapLogger implementation

	// Verify log file contains messages
	content, err := os.ReadFile(logFile)
	assert.NoError(t, err)
	assert.NotEmpty(t, content)

	// Verify content contains our message
	contentStr := string(content)
	assert.Contains(t, contentStr, "test message")
}

func TestLoggerLevels(t *testing.T) {
	tempDir := t.TempDir()
	logFile := filepath.Join(tempDir, "test.log")

	// Test with Error level (should only log errors)
	config := Config{
		Level:     ErrorLevel,
		Format:    "json",
		Output:    logFile,
		Timestamp: true,
	}
	logger := NewLogger(config)

	logger.Debug("debug message")
	logger.Info("info message")
	logger.Warn("warn message")
	logger.Error("error message", assert.AnError)

	// Note: Sync method is not part of the Logger interface
	// but is available on the underlying zapLogger implementation

	content, err := os.ReadFile(logFile)
	assert.NoError(t, err)
	contentStr := string(content)

	// Should only contain error message
	assert.NotContains(t, contentStr, "debug message")
	assert.NotContains(t, contentStr, "info message")
	assert.NotContains(t, contentStr, "warn message")
	assert.Contains(t, contentStr, "error message")
}
