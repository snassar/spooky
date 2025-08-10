package formatters

import (
	"fmt"
	"spooky/internal/types/logging"
)

// TextFormatter implements text formatting for log entries
type TextFormatter struct{}

// NewTextFormatter creates a new text formatter
func NewTextFormatter() *TextFormatter {
	return &TextFormatter{}
}

// Format formats a log entry as text
func (f *TextFormatter) Format(entry *types.LogEntry) ([]byte, error) {
	// Format timestamp
	timestamp := entry.Timestamp.Format("2006-01-02 15:04:05.000")

	// Format level with padding
	level := fmt.Sprintf("%-5s", string(entry.Level))

	// Start with basic log line
	logLine := fmt.Sprintf("%s %s %s", timestamp, level, entry.Message)

	// Add fields if present
	if len(entry.Fields) > 0 {
		fieldStr := ""
		for i, field := range entry.Fields {
			if i > 0 {
				fieldStr += " "
			}
			fieldStr += fmt.Sprintf("%s=%v", field.Key, field.Value)
		}
		logLine += " " + fieldStr
	}

	// Add error if present
	if entry.Error != nil {
		logLine += fmt.Sprintf(" error=%s", entry.Error.Error())
	}

	return []byte(logLine + "\n"), nil
}

// GetName returns the formatter name
func (f *TextFormatter) GetName() string {
	return "text"
}
