package logging

import (
	"context"
	"testing"
)

func TestLogger_Info(t *testing.T) {
	logger := GetLogger()

	// Test basic info logging
	logger.Info("Test info message",
		String("test_key", "test_value"),
		Int("test_number", 42),
		Bool("test_bool", true),
	)

	// Basic assertion to use the parameter
	if logger == nil {
		t.Fatal("Logger should not be nil")
	}
}

func TestLogger_Error(t *testing.T) {
	logger := GetLogger()

	// Test error logging
	err := &testError{message: "test error"}
	logger.Error("Test error message", err,
		String("error_context", "test"),
		Int("error_code", 500),
	)

	// Basic assertion to use the parameter
	if logger == nil {
		t.Fatal("Logger should not be nil")
	}
}

func TestLogger_WithFields(t *testing.T) {
	logger := GetLogger()

	// Test logger with fields
	loggerWithFields := logger.WithFields(
		String("request_id", "req-123"),
		String("user", "testuser"),
	)

	loggerWithFields.Info("Message with context",
		String("action", "test_action"),
	)

	// Basic assertion to use the parameter
	if loggerWithFields == nil {
		t.Fatal("Logger with fields should not be nil")
	}
}

func TestLogger_Context(t *testing.T) {
	logger := GetLogger()
	ctx := context.Background()

	// Test context integration
	ctxWithLogger := WithContext(ctx, logger)
	retrievedLogger := FromContext(ctxWithLogger)

	if retrievedLogger == nil {
		t.Error("Expected logger from context, got nil")
	}
}

func TestNewLogger_TextFormat(t *testing.T) {
	// Test text format logger
	logger := NewLogger(Config{
		Level:     DebugLevel,
		Format:    "text",
		Output:    "stdout",
		Timestamp: true,
	})

	logger.Debug("Debug message in text format",
		String("debug_key", "debug_value"),
	)

	// Basic assertion to use the parameter
	if logger == nil {
		t.Fatal("Text format logger should not be nil")
	}
}

func TestNewLogger_JSONFormat(t *testing.T) {
	// Test JSON format logger
	logger := NewLogger(Config{
		Level:     InfoLevel,
		Format:    "json",
		Output:    "stdout",
		Timestamp: true,
	})

	logger.Info("Info message in JSON format",
		String("json_key", "json_value"),
		Int("json_number", 123),
	)

	// Basic assertion to use the parameter
	if logger == nil {
		t.Fatal("JSON format logger should not be nil")
	}
}

// testError implements error interface for testing
type testError struct {
	message string
}

func (e *testError) Error() string {
	return e.message
}
