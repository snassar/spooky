package logging

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringField(t *testing.T) {
	field := String("test_key", "test_value")

	assert.Equal(t, "test_key", field.Key)
	assert.Equal(t, "test_value", field.Value)
}

func TestIntField(t *testing.T) {
	field := Int("number_key", 42)

	assert.Equal(t, "number_key", field.Key)
	assert.Equal(t, 42, field.Value)
}

func TestInt64Field(t *testing.T) {
	field := Int64("big_number_key", 123456789)

	assert.Equal(t, "big_number_key", field.Key)
	assert.Equal(t, int64(123456789), field.Value)
}

func TestFloat64Field(t *testing.T) {
	field := Float64("float_key", 3.14159)

	assert.Equal(t, "float_key", field.Key)
	assert.Equal(t, 3.14159, field.Value)
}

func TestBoolField(t *testing.T) {
	field := Bool("bool_key", true)

	assert.Equal(t, "bool_key", field.Key)
	assert.Equal(t, true, field.Value)
}

func TestErrorField(t *testing.T) {
	testError := errors.New("test error message")
	field := Error(testError)

	assert.Equal(t, "error", field.Key)
	assert.Equal(t, "test error message", field.Value)
}

func TestDurationField(t *testing.T) {
	field := Duration("duration_key", 1500) // 1.5 seconds in ms

	assert.Equal(t, "duration_key", field.Key)
	assert.Equal(t, int64(1500), field.Value)
}

func TestRequestIDField(t *testing.T) {
	field := RequestID("req-12345")

	assert.Equal(t, "request_id", field.Key)
	assert.Equal(t, "req-12345", field.Value)
}

func TestServerField(t *testing.T) {
	field := Server("web-server-01")

	assert.Equal(t, "server", field.Key)
	assert.Equal(t, "web-server-01", field.Value)
}

func TestActionField(t *testing.T) {
	field := Action("deploy-config")

	assert.Equal(t, "action", field.Key)
	assert.Equal(t, "deploy-config", field.Value)
}

func TestHostField(t *testing.T) {
	field := Host("192.168.1.100")

	assert.Equal(t, "host", field.Key)
	assert.Equal(t, "192.168.1.100", field.Value)
}

func TestPortField(t *testing.T) {
	field := Port(22)

	assert.Equal(t, "port", field.Key)
	assert.Equal(t, 22, field.Value)
}

func TestFieldHelpersWithTestDataFromExamples(t *testing.T) {
	// Test field helpers with data similar to what would be used in examples/testing

	// Simulate server configuration from test data
	serverField := Server("example-server")
	assert.Equal(t, "server", serverField.Key)
	assert.Equal(t, "example-server", serverField.Value)

	// Simulate action configuration from test data
	actionField := Action("deploy-config")
	assert.Equal(t, "action", actionField.Key)
	assert.Equal(t, "deploy-config", actionField.Value)

	// Simulate host configuration from test data
	hostField := Host("192.168.1.100")
	assert.Equal(t, "host", hostField.Key)
	assert.Equal(t, "192.168.1.100", hostField.Value)

	// Simulate port configuration from test data
	portField := Port(22)
	assert.Equal(t, "port", portField.Key)
	assert.Equal(t, 22, portField.Value)
}

func TestFieldHelpersEdgeCases(t *testing.T) {
	// Test empty string values
	emptyStringField := String("empty_key", "")
	assert.Equal(t, "empty_key", emptyStringField.Key)
	assert.Equal(t, "", emptyStringField.Value)

	// Test zero values
	zeroIntField := Int("zero_key", 0)
	assert.Equal(t, "zero_key", zeroIntField.Key)
	assert.Equal(t, 0, zeroIntField.Value)

	zeroFloatField := Float64("zero_float_key", 0.0)
	assert.Equal(t, "zero_float_key", zeroFloatField.Key)
	assert.Equal(t, 0.0, zeroFloatField.Value)

	// Test false boolean
	falseBoolField := Bool("false_key", false)
	assert.Equal(t, "false_key", falseBoolField.Key)
	assert.Equal(t, false, falseBoolField.Value)

	// Test nil error
	nilErrorField := Error(nil)
	assert.Equal(t, "error", nilErrorField.Key)
	assert.Equal(t, "<nil>", nilErrorField.Value)

	// Test zero duration
	zeroDurationField := Duration("zero_duration_key", 0)
	assert.Equal(t, "zero_duration_key", zeroDurationField.Key)
	assert.Equal(t, int64(0), zeroDurationField.Value)
}

func TestFieldHelpersWithSpecialCharacters(t *testing.T) {
	// Test with special characters in keys and values
	specialKeyField := String("key-with-dashes", "value with spaces")
	assert.Equal(t, "key-with-dashes", specialKeyField.Key)
	assert.Equal(t, "value with spaces", specialKeyField.Value)

	// Test with unicode characters
	unicodeField := String("unicode_key", "café")
	assert.Equal(t, "unicode_key", unicodeField.Key)
	assert.Equal(t, "café", unicodeField.Value)

	// Test with numbers in keys
	numericKeyField := Int("key123", 456)
	assert.Equal(t, "key123", numericKeyField.Key)
	assert.Equal(t, 456, numericKeyField.Value)
}

func TestFieldHelpersPerformance(t *testing.T) {
	// Test creating many fields quickly
	for i := 0; i < 1000; i++ {
		field := String("perf_key", "perf_value")
		assert.Equal(t, "perf_key", field.Key)
		assert.Equal(t, "perf_value", field.Value)
	}
}

func TestFieldHelpersConcurrency(t *testing.T) {
	// Test field helpers in concurrent scenarios
	done := make(chan bool, 10)

	for i := 0; i < 10; i++ {
		go func(id int) {
			field := Int("concurrent_key", id)
			assert.Equal(t, "concurrent_key", field.Key)
			assert.Equal(t, id, field.Value)
			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}
}

func TestFieldHelpersCombinedUsage(t *testing.T) {
	// Test combining multiple field helpers in a realistic scenario
	fields := []Field{
		Server("web-server-01"),
		Action("deploy-config"),
		Host("192.168.1.100"),
		Port(22),
		String("environment", "production"),
		Int("timeout", 30),
		Bool("parallel", true),
		Duration("duration", 1500),
		RequestID("req-12345"),
	}

	// Verify all fields
	assert.Len(t, fields, 9)

	// Check specific fields
	assert.Equal(t, "server", fields[0].Key)
	assert.Equal(t, "web-server-01", fields[0].Value)

	assert.Equal(t, "action", fields[1].Key)
	assert.Equal(t, "deploy-config", fields[1].Value)

	assert.Equal(t, "host", fields[2].Key)
	assert.Equal(t, "192.168.1.100", fields[2].Value)

	assert.Equal(t, "port", fields[3].Key)
	assert.Equal(t, 22, fields[3].Value)

	assert.Equal(t, "environment", fields[4].Key)
	assert.Equal(t, "production", fields[4].Value)

	assert.Equal(t, "timeout", fields[5].Key)
	assert.Equal(t, 30, fields[5].Value)

	assert.Equal(t, "parallel", fields[6].Key)
	assert.Equal(t, true, fields[6].Value)

	assert.Equal(t, "duration", fields[7].Key)
	assert.Equal(t, int64(1500), fields[7].Value)

	assert.Equal(t, "request_id", fields[8].Key)
	assert.Equal(t, "req-12345", fields[8].Value)
}
