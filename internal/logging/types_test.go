package logging

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogLevelConstants(t *testing.T) {
	// Test that all log level constants are defined correctly
	assert.Equal(t, LogLevel("debug"), DebugLevel)
	assert.Equal(t, LogLevel("info"), InfoLevel)
	assert.Equal(t, LogLevel("warn"), WarnLevel)
	assert.Equal(t, LogLevel("error"), ErrorLevel)
}

func TestFieldStruct(t *testing.T) {
	// Test Field struct creation and access
	field := Field{
		Key:   "test_key",
		Value: "test_value",
	}

	assert.Equal(t, "test_key", field.Key)
	assert.Equal(t, "test_value", field.Value)
}

func TestConfigStruct(t *testing.T) {
	// Test Config struct creation and access
	config := Config{
		Level:     InfoLevel,
		Format:    "json",
		Output:    "stdout",
		Timestamp: true,
	}

	assert.Equal(t, InfoLevel, config.Level)
	assert.Equal(t, "json", config.Format)
	assert.Equal(t, "stdout", config.Output)
	assert.Equal(t, true, config.Timestamp)
}

func TestLoggerInterface(t *testing.T) {
	// Test that our logger implementation satisfies the Logger interface
	logger := NewLogger(Config{
		Level:     InfoLevel,
		Format:    "json",
		Output:    "stdout",
		Timestamp: true,
	})

	// Test all interface methods
	logger.Debug("debug message", String("key", "value"))
	logger.Info("info message", Int("number", 42))
	logger.Warn("warning message", Bool("flag", true))
	logger.Error("error message", assert.AnError, Float64("pi", 3.14))

	// Test WithContext
	ctx := context.Background()
	loggerWithContext := logger.WithContext(ctx)
	assert.NotNil(t, loggerWithContext)

	// Test WithFields
	loggerWithFields := logger.WithFields(String("global", "value"))
	assert.NotNil(t, loggerWithFields)
}

func TestLoggerInterfaceChaining(_ *testing.T) {
	logger := NewLogger(Config{
		Level:     DebugLevel,
		Format:    "json",
		Output:    "stdout",
		Timestamp: true,
	})

	// Test chaining WithContext and WithFields
	ctx := context.Background()

	loggerWithContext := logger.WithContext(ctx)
	loggerWithContextAndFields := loggerWithContext.WithFields(
		String("chained_key", "chained_value"),
		Int("chained_number", 123),
	)

	// Should be able to log with the chained logger
	loggerWithContextAndFields.Info("chained logger test", String("local", "value"))
}

func TestConfigDefaultValues(t *testing.T) {
	// Test Config with default values
	config := Config{}

	// These should be zero values
	assert.Equal(t, LogLevel(""), config.Level)
	assert.Equal(t, "", config.Format)
	assert.Equal(t, "", config.Output)
	assert.Equal(t, false, config.Timestamp)
}

func TestConfigWithCustomValues(t *testing.T) {
	// Test Config with custom values
	config := Config{
		Level:     DebugLevel,
		Format:    "text",
		Output:    "/var/log/app.log",
		Timestamp: false,
	}

	assert.Equal(t, DebugLevel, config.Level)
	assert.Equal(t, "text", config.Format)
	assert.Equal(t, "/var/log/app.log", config.Output)
	assert.Equal(t, false, config.Timestamp)
}

func TestFieldWithVariousTypes(t *testing.T) {
	// Test Field with various value types
	tests := []struct {
		name  string
		key   string
		value interface{}
	}{
		{"string", "string_key", "string_value"},
		{"int", "int_key", 42},
		{"float", "float_key", 3.14},
		{"bool", "bool_key", true},
		{"slice", "slice_key", []string{"a", "b", "c"}},
		{"map", "map_key", map[string]int{"a": 1, "b": 2}},
		{"nil", "nil_key", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			field := Field{
				Key:   tt.key,
				Value: tt.value,
			}

			assert.Equal(t, tt.key, field.Key)
			assert.Equal(t, tt.value, field.Value)
		})
	}
}

func TestLogLevelComparison(t *testing.T) {
	// Test LogLevel comparison
	assert.Equal(t, DebugLevel, LogLevel("debug"))
	assert.Equal(t, InfoLevel, LogLevel("info"))
	assert.Equal(t, WarnLevel, LogLevel("warn"))
	assert.Equal(t, ErrorLevel, LogLevel("error"))

	// Test inequality
	assert.NotEqual(t, DebugLevel, InfoLevel)
	assert.NotEqual(t, InfoLevel, WarnLevel)
	assert.NotEqual(t, WarnLevel, ErrorLevel)
}

func TestConfigEquality(t *testing.T) {
	// Test Config equality
	config1 := Config{
		Level:     InfoLevel,
		Format:    "json",
		Output:    "stdout",
		Timestamp: true,
	}

	config2 := Config{
		Level:     InfoLevel,
		Format:    "json",
		Output:    "stdout",
		Timestamp: true,
	}

	config3 := Config{
		Level:     DebugLevel,
		Format:    "json",
		Output:    "stdout",
		Timestamp: true,
	}

	// Should be equal
	assert.Equal(t, config1, config2)

	// Should not be equal
	assert.NotEqual(t, config1, config3)
}

func TestFieldEquality(t *testing.T) {
	// Test Field equality
	field1 := Field{
		Key:   "test_key",
		Value: "test_value",
	}

	field2 := Field{
		Key:   "test_key",
		Value: "test_value",
	}

	field3 := Field{
		Key:   "different_key",
		Value: "test_value",
	}

	// Should be equal
	assert.Equal(t, field1, field2)

	// Should not be equal
	assert.NotEqual(t, field1, field3)
}

func TestLoggerInterfaceNilHandling(t *testing.T) {
	// Test that logger interface methods handle nil values gracefully
	logger := NewLogger(Config{
		Level:     InfoLevel,
		Format:    "json",
		Output:    "stdout",
		Timestamp: true,
	})

	// Test with nil context
	loggerWithNilContext := logger.WithContext(context.TODO())
	assert.NotNil(t, loggerWithNilContext)

	// Test with nil fields
	loggerWithNilFields := logger.WithFields()
	assert.NotNil(t, loggerWithNilFields)

	// Test error logging with nil error
	logger.Error("message with nil error", nil, String("key", "value"))
}

func TestLoggerInterfaceWithTestDataFromExamples(_ *testing.T) {
	// Test logger interface with data patterns from examples/testing

	logger := NewLogger(Config{
		Level:     InfoLevel,
		Format:    "json",
		Output:    "stdout",
		Timestamp: true,
	})

	// Simulate logging patterns from test data
	logger.Info("Processing project configuration",
		String("project", "test-valid-project"),
		String("environment", "development"),
		String("version", "1.0.0"),
	)

	logger.Warn("Log file unwritable",
		String("project", "test-log-file-unwritable"),
		String("log_file", "/root/unwritable.log"),
		String("fallback", "stdout"),
	)

	logger.Error("SSH connection failed",
		assert.AnError,
		String("server", "example-server"),
		String("host", "192.168.1.100"),
		Int("port", 22),
		Int("timeout", 30),
	)
}

func TestConfigValidation(t *testing.T) {
	// Test Config validation scenarios
	tests := []struct {
		name   string
		config Config
		valid  bool
	}{
		{
			name: "valid config",
			config: Config{
				Level:     InfoLevel,
				Format:    "json",
				Output:    "stdout",
				Timestamp: true,
			},
			valid: true,
		},
		{
			name: "valid config with text format",
			config: Config{
				Level:     DebugLevel,
				Format:    "text",
				Output:    "/tmp/test.log",
				Timestamp: false,
			},
			valid: true,
		},
		{
			name: "config with empty level",
			config: Config{
				Level:     "",
				Format:    "json",
				Output:    "stdout",
				Timestamp: true,
			},
			valid: true, // Empty level is valid, will use default
		},
		{
			name: "config with empty format",
			config: Config{
				Level:     InfoLevel,
				Format:    "",
				Output:    "stdout",
				Timestamp: true,
			},
			valid: true, // Empty format is valid, will use default
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create logger with config - should not panic
			logger := NewLogger(tt.config)
			assert.NotNil(t, logger)

			// Should be able to log
			logger.Info("config validation test", String("test", "value"))
		})
	}
}
