package logging

import (
	"spooky/internal/logging/types"
)

// Field helpers for common data types

// String creates a string field
func String(key, value string) types.Field {
	return types.Field{Key: key, Value: value}
}

// Int creates an integer field
func Int(key string, value int) types.Field {
	return types.Field{Key: key, Value: value}
}

// Int64 creates an int64 field
func Int64(key string, value int64) types.Field {
	return types.Field{Key: key, Value: value}
}

// Float64 creates a float64 field
func Float64(key string, value float64) types.Field {
	return types.Field{Key: key, Value: value}
}

// Bool creates a boolean field
func Bool(key string, value bool) types.Field {
	return types.Field{Key: key, Value: value}
}

// Error creates an error field
func Error(err error) types.Field {
	if err == nil {
		return types.Field{Key: "error", Value: "<nil>"}
	}
	return types.Field{Key: "error", Value: err.Error()}
}

// Duration creates a duration field (in milliseconds)
func Duration(key string, durationMs int64) types.Field {
	return types.Field{Key: key, Value: durationMs}
}

// RequestID creates a request ID field
func RequestID(id string) types.Field {
	return types.Field{Key: "request_id", Value: id}
}

// Server creates a server field
func Server(name string) types.Field {
	return types.Field{Key: "server", Value: name}
}

// Action creates an action field
func Action(name string) types.Field {
	return types.Field{Key: "action", Value: name}
}

// Host creates a host field
func Host(host string) types.Field {
	return types.Field{Key: "host", Value: host}
}

// Port creates a port field
func Port(port int) types.Field {
	return types.Field{Key: "port", Value: port}
}

// StringSlice creates a string slice field
func StringSlice(key string, value []string) types.Field {
	return types.Field{Key: key, Value: value}
}
