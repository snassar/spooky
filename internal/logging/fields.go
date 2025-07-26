package logging

// Field helpers for common data types

// String creates a string field
func String(key, value string) Field {
	return Field{Key: key, Value: value}
}

// Int creates an integer field
func Int(key string, value int) Field {
	return Field{Key: key, Value: value}
}

// Int64 creates an int64 field
func Int64(key string, value int64) Field {
	return Field{Key: key, Value: value}
}

// Float64 creates a float64 field
func Float64(key string, value float64) Field {
	return Field{Key: key, Value: value}
}

// Bool creates a boolean field
func Bool(key string, value bool) Field {
	return Field{Key: key, Value: value}
}

// Error creates an error field
func Error(err error) Field {
	if err == nil {
		return Field{Key: "error", Value: "<nil>"}
	}
	return Field{Key: "error", Value: err.Error()}
}

// Duration creates a duration field (in milliseconds)
func Duration(key string, durationMs int64) Field {
	return Field{Key: key, Value: durationMs}
}

// RequestID creates a request ID field
func RequestID(id string) Field {
	return Field{Key: "request_id", Value: id}
}

// Server creates a server field
func Server(name string) Field {
	return Field{Key: "server", Value: name}
}

// Action creates an action field
func Action(name string) Field {
	return Field{Key: "action", Value: name}
}

// Host creates a host field
func Host(host string) Field {
	return Field{Key: "host", Value: host}
}

// Port creates a port field
func Port(port int) Field {
	return Field{Key: "port", Value: port}
}
