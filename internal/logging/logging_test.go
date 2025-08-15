package logging

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"log/slog"

	spookytypeslogging "spooky/internal/types/logging"
)

func TestLogManager(t *testing.T) {
	// Create a new log manager
	manager := NewLogManager()
	defer manager.Close()

	// Test basic logger creation
	logger := manager.GetLogger("test-component")
	if logger == nil {
		t.Fatal("Failed to create logger")
	}

	// Test log level setting
	manager.SetLevel(spookytypeslogging.LogLevelDebug)
	if manager.GetLevel() != spookytypeslogging.LogLevelDebug {
		t.Errorf("Expected log level %s, got %s", spookytypeslogging.LogLevelDebug, manager.GetLevel())
	}

	// Test logging at different levels
	logger.Info("Test info message", map[string]interface{}{"test": "value"})
	logger.Warn("Test warning message", map[string]interface{}{"warning": "test"})
	logger.Error("Test error message", fmt.Errorf("test error"), map[string]interface{}{"error": "test"})
}

func TestLogManagerWithFileOutput(t *testing.T) {
	// Create temporary log file
	tempFile := t.TempDir() + "/test.log"

	// Create log manager with file output
	config := &spookytypeslogging.LogConfig{
		Level:  spookytypeslogging.LogLevelInfo,
		Format: "json",
		Output: "file",
		File: &spookytypeslogging.LogFileConfig{
			Path:        tempFile,
			Permissions: "0644",
			Append:      false,
		},
	}

	manager := NewLogManager()
	defer manager.Close()

	// Configure logging
	if err := manager.Configure(config); err != nil {
		t.Fatalf("Failed to configure logging: %v", err)
	}

	// Get logger and log some messages
	logger := manager.GetLogger("file-test")
	logger.Info("File test message", map[string]interface{}{"test": "file"})
	logger.Warn("File warning message")

	// Check if file was created and contains logs
	time.Sleep(100 * time.Millisecond) // Give time for writes

	content, err := os.ReadFile(tempFile)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	if len(content) == 0 {
		t.Error("Log file is empty")
	}

	// Check if content contains expected JSON
	contentStr := string(content)
	if !strings.Contains(contentStr, "File test message") {
		t.Error("Log file does not contain expected message")
	}
}

func TestLoggerWithFields(_ *testing.T) {
	manager := NewLogManager()
	defer manager.Close()

	logger := manager.GetLogger("fields-test")

	// Test WithFields
	fieldLogger := logger.WithFields(map[string]interface{}{
		"user_id": "12345",
		"session": "abc123",
	})

	fieldLogger.Info("Message with fields", map[string]interface{}{"additional": "data"})
}

func TestLoggerWithComponent(_ *testing.T) {
	manager := NewLogManager()
	defer manager.Close()

	logger := manager.GetLogger("original")

	// Test WithComponent
	componentLogger := logger.WithComponent("new-component")
	componentLogger.Info("Message from new component")
}

func TestLoggerWithOperation(_ *testing.T) {
	manager := NewLogManager()
	defer manager.Close()

	logger := manager.GetLogger("operation-test")

	// Test WithOperation
	operationLogger := logger.WithOperation("user-authentication")
	operationLogger.Info("User login attempt", map[string]interface{}{"user": "testuser"})
}

func TestLogLevelConversion(t *testing.T) {
	// Test LogLevel to slog.Level conversion
	testCases := []struct {
		logLevel spookytypeslogging.LogLevel
		expected slog.Level
	}{
		{spookytypeslogging.LogLevelDebug, slog.LevelDebug},
		{spookytypeslogging.LogLevelInfo, slog.LevelInfo},
		{spookytypeslogging.LogLevelWarn, slog.LevelWarn},
		{spookytypeslogging.LogLevelError, slog.LevelError},
		{spookytypeslogging.LogLevelFatal, slog.LevelError}, // slog doesn't have fatal
	}

	for _, tc := range testCases {
		result := tc.logLevel.ToSlogLevel()
		if result != tc.expected {
			t.Errorf("Expected %v for %s, got %v", tc.expected, tc.logLevel, result)
		}
	}

	// Test slog.Level to LogLevel conversion
	slogTestCases := []struct {
		slogLevel slog.Level
		expected  spookytypeslogging.LogLevel
	}{
		{slog.LevelDebug, spookytypeslogging.LogLevelDebug},
		{slog.LevelInfo, spookytypeslogging.LogLevelInfo},
		{slog.LevelWarn, spookytypeslogging.LogLevelWarn},
		{slog.LevelError, spookytypeslogging.LogLevelError},
	}

	for _, tc := range slogTestCases {
		result := spookytypeslogging.FromSlogLevel(tc.slogLevel)
		if result != tc.expected {
			t.Errorf("Expected %v for %v, got %v", tc.expected, tc.slogLevel, result)
		}
	}
}

func TestLogManagerWithFiltering(t *testing.T) {
	// Create log manager with filtering
	config := &spookytypeslogging.LogConfig{
		Level:  spookytypeslogging.LogLevelDebug,
		Format: "json",
		Output: "stderr",
		Filtering: &spookytypeslogging.LogFilteringConfig{
			Components: map[string]spookytypeslogging.LogLevel{
				"filtered": spookytypeslogging.LogLevelError, // Only errors and above
			},
		},
	}

	manager := NewLogManager()
	defer manager.Close()

	// Configure logging
	if err := manager.Configure(config); err != nil {
		t.Fatalf("Failed to configure logging: %v", err)
	}

	// Get logger for filtered component
	logger := manager.GetLogger("filtered")

	// These should be filtered out (below error level)
	logger.Debug("This debug message should be filtered")
	logger.Info("This info message should be filtered")
	logger.Warn("This warning message should be filtered")

	// This should not be filtered (error level)
	logger.Error("This error message should not be filtered", fmt.Errorf("test error"))
}

func TestLogManagerWithStructuredConfig(t *testing.T) {
	// Create log manager with structured configuration
	config := &spookytypeslogging.LogConfig{
		Level:  spookytypeslogging.LogLevelInfo,
		Format: "structured",
		Output: "stderr",
		Structured: &spookytypeslogging.LogStructuredConfig{
			Timestamp: &spookytypeslogging.LogTimestampConfig{
				Enabled:  true,
				Format:   "RFC3339",
				Timezone: "UTC",
			},
			Component: &spookytypeslogging.LogComponentConfig{
				Key:     "component",
				Enabled: true,
			},
			Caller: &spookytypeslogging.LogCallerConfig{
				Enabled:    true,
				Key:        "caller",
				SkipFrames: 2,
			},
		},
	}

	manager := NewLogManager()
	defer manager.Close()

	// Configure logging
	if err := manager.Configure(config); err != nil {
		t.Fatalf("Failed to configure logging: %v", err)
	}

	// Get logger and test structured logging
	logger := manager.GetLogger("structured-test")
	logger.Info("Structured test message", map[string]interface{}{"structured": "test"})
}

func TestLogManagerNullOutput(t *testing.T) {
	// Test null output (no logging)
	config := &spookytypeslogging.LogConfig{
		Level:  spookytypeslogging.LogLevelInfo,
		Format: "json",
		Output: "null",
	}

	manager := NewLogManager()
	defer manager.Close()

	// Configure logging
	if err := manager.Configure(config); err != nil {
		t.Fatalf("Failed to configure logging: %v", err)
	}

	// Get logger and log (should not output anything)
	logger := manager.GetLogger("null-test")
	logger.Info("This message should not appear anywhere")
	logger.Error("This error should not appear anywhere", fmt.Errorf("test error"))
}

func TestLogManagerInvalidConfig(t *testing.T) {
	manager := NewLogManager()
	defer manager.Close()

	// Test invalid file path
	config := &spookytypeslogging.LogConfig{
		Level:  spookytypeslogging.LogLevelInfo,
		Format: "json",
		Output: "file",
		File: &spookytypeslogging.LogFileConfig{
			Path: "/invalid/path/that/should/not/exist/test.log",
		},
	}

	// This should fail due to invalid path
	if err := manager.Configure(config); err == nil {
		t.Error("Expected error for invalid file path, got nil")
	}
}
