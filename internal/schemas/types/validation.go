package types

// ValidationError represents a unified validation error with detailed information
type ValidationError struct {
	File     string `json:"file"`
	Line     int    `json:"line"`
	Column   int    `json:"column"`
	Field    string `json:"field"`
	Message  string `json:"message"`
	Value    string `json:"value,omitempty"`
	Severity string `json:"severity"` // "error" or "warning"
}

// ValidationResult contains unified validation results
type ValidationResult struct {
	Valid    bool              `json:"valid"`
	Errors   []ValidationError `json:"errors,omitempty"`
	Warnings []ValidationError `json:"warnings,omitempty"`
	File     string            `json:"file,omitempty"`
	Schema   string            `json:"schema,omitempty"`
}
