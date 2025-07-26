package logging

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewLogger(t *testing.T) {
	tests := []struct {
		name   string
		config Config
	}{
		{
			name: "default json logger",
			config: Config{
				Level:     InfoLevel,
				Format:    "json",
				Output:    "stdout",
				Timestamp: true,
			},
		},
		{
			name: "debug level logger",
			config: Config{
				Level:     DebugLevel,
				Format:    "json",
				Output:    "stdout",
				Timestamp: true,
			},
		},
		{
			name: "warn level logger",
			config: Config{
				Level:     WarnLevel,
				Format:    "json",
				Output:    "stdout",
				Timestamp: true,
			},
		},
		{
			name: "error level logger",
			config: Config{
				Level:     ErrorLevel,
				Format:    "json",
				Output:    "stdout",
				Timestamp: true,
			},
		},
		{
			name: "text format logger",
			config: Config{
				Level:     InfoLevel,
				Format:    "text",
				Output:    "stdout",
				Timestamp: true,
			},
		},
		{
			name: "no timestamp logger",
			config: Config{
				Level:     InfoLevel,
				Format:    "json",
				Output:    "stdout",
				Timestamp: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := NewLogger(tt.config)
			assert.NotNil(t, logger)

			// Test that we can log messages
			logger.Info("test message", String("test_key", "test_value"))
			logger.Debug("debug message", Int("debug_key", 42))
			logger.Warn("warning message", Bool("warn_key", true))
			logger.Error("error message", assert.AnError, Float64("error_key", 3.14))
		})
	}
}

func TestConfigureLogger(t *testing.T) {
	tests := []struct {
		name    string
		level   string
		format  string
		output  string
		quiet   bool
		verbose bool
	}{
		{
			name:    "default configuration",
			level:   "info",
			format:  "json",
			output:  "stdout",
			quiet:   false,
			verbose: false,
		},
		{
			name:    "debug level with verbose flag",
			level:   "info",
			format:  "json",
			output:  "stdout",
			quiet:   false,
			verbose: true,
		},
		{
			name:    "error level with quiet flag",
			level:   "info",
			format:  "json",
			output:  "stdout",
			quiet:   true,
			verbose: false,
		},
		{
			name:    "text format",
			level:   "info",
			format:  "text",
			output:  "stdout",
			quiet:   false,
			verbose: false,
		},
		{
			name:    "custom output file",
			level:   "info",
			format:  "json",
			output:  "/tmp/test.log",
			quiet:   false,
			verbose: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ConfigureLogger(tt.level, tt.format, tt.output, tt.quiet, tt.verbose)

			logger := GetLogger()
			assert.NotNil(t, logger)

			// Test that we can log messages
			logger.Info("configured logger test", String("test", "value"))
		})
	}
}

func TestLoggerWithFields(_ *testing.T) {
	logger := NewLogger(Config{
		Level:     DebugLevel,
		Format:    "json",
		Output:    "stdout",
		Timestamp: true,
	})

	// Test WithFields
	loggerWithFields := logger.WithFields(
		String("global_key", "global_value"),
		Int("global_number", 42),
	)

	// Test that the logger with fields works
	loggerWithFields.Info("message with global fields", String("local_key", "local_value"))

	// Test chaining WithFields
	loggerWithMoreFields := loggerWithFields.WithFields(
		Bool("chained_key", true),
		Float64("chained_number", 3.14),
	)

	loggerWithMoreFields.Debug("message with chained fields")
}

func TestLoggerWithContext(t *testing.T) {
	logger := NewLogger(Config{
		Level:     InfoLevel,
		Format:    "json",
		Output:    "stdout",
		Timestamp: true,
	})

	ctx := context.Background()

	// Test WithContext (should be no-op for Zap implementation)
	loggerWithContext := logger.WithContext(ctx)
	assert.NotNil(t, loggerWithContext)

	// Test that we can still log
	loggerWithContext.Info("message with context", String("ctx_key", "ctx_value"))
}

func TestGlobalLoggerFunctions(t *testing.T) {
	// Test GetLogger
	logger := GetLogger()
	assert.NotNil(t, logger)

	// Test SetLogger
	newLogger := NewLogger(Config{
		Level:     DebugLevel,
		Format:    "json",
		Output:    "stdout",
		Timestamp: true,
	})

	SetLogger(newLogger)

	// Verify the logger was set
	currentLogger := GetLogger()
	assert.Equal(t, newLogger, currentLogger)
}

func TestContextLoggerFunctions(t *testing.T) {
	logger := NewLogger(Config{
		Level:     InfoLevel,
		Format:    "json",
		Output:    "stdout",
		Timestamp: true,
	})

	ctx := context.Background()

	// Test WithContext
	ctxWithLogger := WithContext(ctx, logger)
	assert.NotNil(t, ctxWithLogger)

	// Test FromContext
	loggerFromContext := FromContext(ctxWithLogger)
	assert.Equal(t, logger, loggerFromContext)

	// Test FromContext with no logger in context (should return global)
	emptyCtx := context.Background()
	globalLogger := FromContext(emptyCtx)
	assert.NotNil(t, globalLogger)
}

func TestEnsureLogDirectory(t *testing.T) {
	// Test with stdout (should not create directory)
	err := ensureLogDirectory("stdout")
	assert.NoError(t, err)

	// Test with stderr (should not create directory)
	err = ensureLogDirectory("stderr")
	assert.NoError(t, err)

	// Test with empty string (should not create directory)
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

func TestLoggerWithTestDataFromExamples(t *testing.T) {
	// Use test data from examples/testing where available
	examplesDir := "../../examples/testing"

	// Test with valid project configuration
	validProjectPath := filepath.Join(examplesDir, "test-valid-project")
	if _, err := os.Stat(validProjectPath); os.IsNotExist(err) {
		t.Skip("Test data directory not found, skipping test")
	}

	// Create logger similar to what would be used in a valid project
	logger := NewLogger(Config{
		Level:     InfoLevel,
		Format:    "json",
		Output:    "stdout",
		Timestamp: true,
	})

	// Test logging with project-like fields
	logger.Info("Processing valid project",
		String("project", "test-valid-project"),
		String("environment", "development"),
		String("version", "1.0.0"),
	)

	// Test with unwritable log file scenario
	unwritableProjectPath := filepath.Join(examplesDir, "test-log-file-unwritable")
	if _, err := os.Stat(unwritableProjectPath); err == nil {
		logger.Warn("Attempting to write to unwritable log file",
			String("project", "test-log-file-unwritable"),
			String("log_file", "/root/unwritable.log"),
		)
	}
}

func TestLoggerLevels(_ *testing.T) {
	logger := NewLogger(Config{
		Level:     DebugLevel,
		Format:    "json",
		Output:    "stdout",
		Timestamp: true,
	})

	// Test all log levels
	logger.Debug("debug message", String("level", "debug"))
	logger.Info("info message", String("level", "info"))
	logger.Warn("warning message", String("level", "warn"))
	logger.Error("error message", assert.AnError, String("level", "error"))
}

func TestLoggerErrorHandling(_ *testing.T) {
	logger := NewLogger(Config{
		Level:     ErrorLevel,
		Format:    "json",
		Output:    "stdout",
		Timestamp: true,
	})

	// Test error logging with nil error
	logger.Error("message with nil error", nil, String("test", "value"))

	// Test error logging with actual error
	testError := assert.AnError
	logger.Error("message with error", testError, String("error_type", "test_error"))
}

func TestLoggerFieldConversion(_ *testing.T) {
	logger := NewLogger(Config{
		Level:     DebugLevel,
		Format:    "json",
		Output:    "stdout",
		Timestamp: true,
	})

	// Test various field types
	logger.Info("message with various field types",
		String("string_field", "string_value"),
		Int("int_field", 42),
		Int64("int64_field", 123456789),
		Float64("float64_field", 3.14159),
		Bool("bool_field", true),
		Duration("duration_field", 1500), // 1.5 seconds in ms
		RequestID("req-123"),
		Server("test-server"),
		Action("test-action"),
		Host("192.168.1.100"),
		Port(22),
	)
}

func TestLoggerSync(t *testing.T) {
	logger := NewLogger(Config{
		Level:     InfoLevel,
		Format:    "json",
		Output:    "stdout",
		Timestamp: true,
	})

	// Test Sync method - may fail on stdout, which is expected
	err := logger.(*zapLogger).Sync()
	// Sync can fail on stdout/stderr, which is normal behavior
	if err != nil {
		assert.Contains(t, err.Error(), "invalid argument")
	}
}

func TestLoggerInvalidConfigurations(t *testing.T) {
	tests := []struct {
		name   string
		config Config
	}{
		{
			name: "invalid level",
			config: Config{
				Level:     "invalid",
				Format:    "json",
				Output:    "stdout",
				Timestamp: true,
			},
		},
		{
			name: "invalid format",
			config: Config{
				Level:     InfoLevel,
				Format:    "invalid",
				Output:    "stdout",
				Timestamp: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Should not panic and should create a fallback logger
			logger := NewLogger(tt.config)
			assert.NotNil(t, logger)

			// Should still be able to log
			logger.Info("test message with invalid config")
		})
	}
}

func TestLoggerPerformance(t *testing.T) {
	logger := NewLogger(Config{
		Level:     InfoLevel,
		Format:    "json",
		Output:    "stdout",
		Timestamp: true,
	})

	// Test logging performance with many fields
	start := time.Now()

	for i := 0; i < 100; i++ {
		logger.Info("performance test message",
			Int("iteration", i),
			String("test", "performance"),
			Bool("success", true),
		)
	}

	duration := time.Since(start)
	assert.Less(t, duration, 5*time.Second, "Logging 100 messages should complete quickly")
}

func TestLoggerConcurrency(_ *testing.T) {
	logger := NewLogger(Config{
		Level:     InfoLevel,
		Format:    "json",
		Output:    "stdout",
		Timestamp: true,
	})

	// Test concurrent logging
	done := make(chan bool, 10)

	for i := 0; i < 10; i++ {
		go func(id int) {
			logger.Info("concurrent message", Int("goroutine_id", id))
			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}
}
