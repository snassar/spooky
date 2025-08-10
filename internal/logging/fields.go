package logging

import (
	spookytypeslogging "spooky/internal/types/logging"
)

// Field helpers for common data types

// String creates a string field
func String(key, value string) spookytypeslogging.Field {
	return spookytypeslogging.Field{Key: key, Value: value}
}

// Int creates an integer field
func Int(key string, value int) spookytypeslogging.Field {
	return spookytypeslogging.Field{Key: key, Value: value}
}

// Int64 creates an int64 field
func Int64(key string, value int64) spookytypeslogging.Field {
	return spookytypeslogging.Field{Key: key, Value: value}
}

// Float64 creates a float64 field
func Float64(key string, value float64) spookytypeslogging.Field {
	return spookytypeslogging.Field{Key: key, Value: value}
}

// Bool creates a boolean field
func Bool(key string, value bool) spookytypeslogging.Field {
	return spookytypeslogging.Field{Key: key, Value: value}
}

// Error creates an error field
func Error(err error) spookytypeslogging.Field {
	if err == nil {
		return spookytypeslogging.Field{Key: "error", Value: "<nil>"}
	}
	return spookytypeslogging.Field{Key: "error", Value: err.Error()}
}

// Duration creates a duration field (in milliseconds)
func Duration(key string, durationMs int64) spookytypeslogging.Field {
	return spookytypeslogging.Field{Key: key, Value: durationMs}
}

// RequestID creates a request ID field
func RequestID(id string) spookytypeslogging.Field {
	return spookytypeslogging.Field{Key: "request_id", Value: id}
}

// Server creates a server field
func Server(name string) spookytypeslogging.Field {
	return spookytypeslogging.Field{Key: "server", Value: name}
}

// Action creates an action field
func Action(name string) spookytypeslogging.Field {
	return spookytypeslogging.Field{Key: "action", Value: name}
}

// Host creates a host field
func Host(host string) spookytypeslogging.Field {
	return spookytypeslogging.Field{Key: "host", Value: host}
}

// Port creates a port field
func Port(port int) spookytypeslogging.Field {
	return spookytypeslogging.Field{Key: "port", Value: port}
}

// StringSlice creates a string slice field
func StringSlice(key string, value []string) spookytypeslogging.Field {
	return spookytypeslogging.Field{Key: key, Value: value}
}
