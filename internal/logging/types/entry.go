package types

import (
	"time"
)

// LogEntry represents a single log entry
type LogEntry struct {
	Timestamp time.Time
	Level     LogLevel
	Message   string
	Fields    []Field
	Error     error
}

// Field represents a structured logging field
type Field struct {
	Key   string
	Value interface{}
}
