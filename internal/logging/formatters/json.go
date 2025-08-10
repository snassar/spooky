package formatters

import (
	"encoding/json"
	"spooky/internal/types/logging"
)

// JSONFormatter implements JSON formatting for log entries
type JSONFormatter struct{}

// NewJSONFormatter creates a new JSON formatter
func NewJSONFormatter() *JSONFormatter {
	return &JSONFormatter{}
}

// Format formats a log entry as JSON
func (f *JSONFormatter) Format(entry *types.LogEntry) ([]byte, error) {
	// Create a structured log entry for JSON output
	logData := map[string]interface{}{
		"timestamp": entry.Timestamp.Format("2006-01-02T15:04:05.000Z07:00"),
		"level":     string(entry.Level),
		"message":   entry.Message,
	}

	// Add fields
	if len(entry.Fields) > 0 {
		fields := make(map[string]interface{})
		for _, field := range entry.Fields {
			fields[field.Key] = field.Value
		}
		logData["fields"] = fields
	}

	// Add error if present
	if entry.Error != nil {
		logData["error"] = entry.Error.Error()
	}

	return json.Marshal(logData)
}

// GetName returns the formatter name
func (f *JSONFormatter) GetName() string {
	return "json"
}
