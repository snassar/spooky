package utilities

import (
	"errors"
	"log/slog"
	"os"
	"testing"
)

func TestWrapError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		context  string
		expected string
	}{
		{
			name:     "nil error",
			err:      nil,
			context:  "test context",
			expected: "",
		},
		{
			name:     "simple error",
			err:      errors.New("original error"),
			context:  "test context",
			expected: "test context: original error",
		},
		{
			name:     "empty context",
			err:      errors.New("original error"),
			context:  "",
			expected: ": original error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := WrapError(tt.err, tt.context)

			if tt.err == nil {
				if result != nil {
					t.Errorf("Expected nil result for nil error, got %v", result)
				}
				return
			}

			if result == nil {
				t.Errorf("Expected non-nil result for non-nil error")
				return
			}

			if result.Error() != tt.expected {
				t.Errorf("Expected error message '%s', got '%s'", tt.expected, result.Error())
			}
		})
	}
}

func TestWrapErrorf(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		format   string
		args     []interface{}
		expected string
	}{
		{
			name:     "nil error",
			err:      nil,
			format:   "failed to %s",
			args:     []interface{}{"process"},
			expected: "",
		},
		{
			name:     "formatted context",
			err:      errors.New("original error"),
			format:   "failed to %s %s",
			args:     []interface{}{"process", "file"},
			expected: "failed to process file: original error",
		},
		{
			name:     "no format args",
			err:      errors.New("original error"),
			format:   "simple context",
			args:     []interface{}{},
			expected: "simple context: original error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := WrapErrorf(tt.err, tt.format, tt.args...)

			if tt.err == nil {
				if result != nil {
					t.Errorf("Expected nil result for nil error, got %v", result)
				}
				return
			}

			if result == nil {
				t.Errorf("Expected non-nil result for non-nil error")
				return
			}

			if result.Error() != tt.expected {
				t.Errorf("Expected error message '%s', got '%s'", tt.expected, result.Error())
			}
		})
	}
}

func TestNewError(t *testing.T) {
	tests := []struct {
		name     string
		message  string
		expected string
	}{
		{
			name:     "simple message",
			message:  "test error",
			expected: "test error",
		},
		{
			name:     "empty message",
			message:  "",
			expected: "",
		},
		{
			name:     "message with special characters",
			message:  "error: file not found: /path/to/file",
			expected: "error: file not found: /path/to/file",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewError(tt.message)

			if result == nil {
				t.Errorf("Expected non-nil result")
				return
			}

			if result.Error() != tt.expected {
				t.Errorf("Expected error message '%s', got '%s'", tt.expected, result.Error())
			}
		})
	}
}

func TestNewErrorf(t *testing.T) {
	tests := []struct {
		name     string
		format   string
		args     []interface{}
		expected string
	}{
		{
			name:     "formatted message",
			format:   "failed to %s %s",
			args:     []interface{}{"process", "file"},
			expected: "failed to process file",
		},
		{
			name:     "no format args",
			format:   "simple error",
			args:     []interface{}{},
			expected: "simple error",
		},
		{
			name:     "multiple format args",
			format:   "error %d: %s in %s",
			args:     []interface{}{404, "not found", "/path"},
			expected: "error 404: not found in /path",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewErrorf(tt.format, tt.args...)

			if result == nil {
				t.Errorf("Expected non-nil result")
				return
			}

			if result.Error() != tt.expected {
				t.Errorf("Expected error message '%s', got '%s'", tt.expected, result.Error())
			}
		})
	}
}

func TestLogError(t *testing.T) {
	// Create a test logger that writes to a buffer
	logFile, err := os.CreateTemp("", "test_log_*.log")
	if err != nil {
		t.Fatalf("Failed to create temp log file: %v", err)
	}
	defer func() { _ = os.Remove(logFile.Name()) }()
	defer func() { _ = logFile.Close() }()

	logger := slog.New(slog.NewTextHandler(logFile, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	tests := []struct {
		name    string
		err     error
		context string
		logger  *slog.Logger
	}{
		{
			name:    "nil error",
			err:     nil,
			context: "test context",
			logger:  logger,
		},
		{
			name:    "simple error",
			err:     errors.New("test error"),
			context: "test context",
			logger:  logger,
		},
		{
			name:    "nil logger uses default",
			err:     errors.New("test error"),
			context: "test context",
			logger:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Should not panic
			LogError(tt.err, tt.context, tt.logger)
		})
	}
}

func TestLogErrorf(t *testing.T) {
	// Create a test logger that writes to a buffer
	logFile, err := os.CreateTemp("", "test_log_*.log")
	if err != nil {
		t.Fatalf("Failed to create temp log file: %v", err)
	}
	defer func() { _ = os.Remove(logFile.Name()) }()
	defer func() { _ = logFile.Close() }()

	logger := slog.New(slog.NewTextHandler(logFile, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	tests := []struct {
		name   string
		err    error
		format string
		args   []interface{}
		logger *slog.Logger
	}{
		{
			name:   "nil error",
			err:    nil,
			format: "failed to %s",
			args:   []interface{}{"process"},
			logger: logger,
		},
		{
			name:   "formatted context",
			err:    errors.New("test error"),
			format: "failed to %s %s",
			args:   []interface{}{"process", "file"},
			logger: logger,
		},
		{
			name:   "nil logger uses default",
			err:    errors.New("test error"),
			format: "test %s",
			args:   []interface{}{"context"},
			logger: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Should not panic
			LogErrorf(tt.err, tt.logger, tt.format, tt.args...)
		})
	}
}

func TestLogWarning(t *testing.T) {
	// Create a test logger that writes to a buffer
	logFile, err := os.CreateTemp("", "test_log_*.log")
	if err != nil {
		t.Fatalf("Failed to create temp log file: %v", err)
	}
	defer func() { _ = os.Remove(logFile.Name()) }()
	defer func() { _ = logFile.Close() }()

	logger := slog.New(slog.NewTextHandler(logFile, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	tests := []struct {
		name    string
		err     error
		context string
		logger  *slog.Logger
	}{
		{
			name:    "nil error",
			err:     nil,
			context: "test context",
			logger:  logger,
		},
		{
			name:    "simple warning",
			err:     errors.New("test warning"),
			context: "test context",
			logger:  logger,
		},
		{
			name:    "nil logger uses default",
			err:     errors.New("test warning"),
			context: "test context",
			logger:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Should not panic
			LogWarning(tt.err, tt.context, tt.logger)
		})
	}
}

func TestLogWarningf(t *testing.T) {
	// Create a test logger that writes to a buffer
	logFile, err := os.CreateTemp("", "test_log_*.log")
	if err != nil {
		t.Fatalf("Failed to create temp log file: %v", err)
	}
	defer func() { _ = os.Remove(logFile.Name()) }()
	defer func() { _ = logFile.Close() }()

	logger := slog.New(slog.NewTextHandler(logFile, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	tests := []struct {
		name   string
		err    error
		format string
		args   []interface{}
		logger *slog.Logger
	}{
		{
			name:   "nil error",
			err:    nil,
			format: "warning in %s",
			args:   []interface{}{"process"},
			logger: logger,
		},
		{
			name:   "formatted warning",
			err:    errors.New("test warning"),
			format: "warning in %s %s",
			args:   []interface{}{"process", "file"},
			logger: logger,
		},
		{
			name:   "nil logger uses default",
			err:    errors.New("test warning"),
			format: "test %s",
			args:   []interface{}{"warning"},
			logger: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Should not panic
			LogWarningf(tt.err, tt.logger, tt.format, tt.args...)
		})
	}
}

// Test error unwrapping behavior
func TestErrorUnwrapping(t *testing.T) {
	originalErr := errors.New("original error")
	wrappedErr := WrapError(originalErr, "context")

	// Test that errors.Is works correctly
	if !errors.Is(wrappedErr, originalErr) {
		t.Errorf("Expected wrapped error to be unwrappable to original error")
	}

	// Test that the wrapped error contains the original error
	if wrappedErr.Error() == originalErr.Error() {
		t.Errorf("Expected wrapped error to have additional context")
	}

	// Test that the wrapped error message contains the context
	expectedContext := "context: original error"
	if wrappedErr.Error() != expectedContext {
		t.Errorf("Expected wrapped error message '%s', got '%s'", expectedContext, wrappedErr.Error())
	}
}

// Benchmark tests for performance
func BenchmarkWrapError(b *testing.B) {
	err := errors.New("test error")
	context := "test context"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = WrapError(err, context)
	}
}

func BenchmarkWrapErrorf(b *testing.B) {
	err := errors.New("test error")
	format := "failed to %s %s"
	args := []interface{}{"process", "file"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = WrapErrorf(err, format, args...)
	}
}

func BenchmarkNewError(b *testing.B) {
	message := "test error message"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewError(message)
	}
}

func BenchmarkNewErrorf(b *testing.B) {
	format := "error %d: %s"
	args := []interface{}{404, "not found"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewErrorf(format, args...)
	}
}
