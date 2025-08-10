package reporters

import (
	"encoding/json"
	"fmt"
	"os"
	"spooky/internal/types/schemas"
)

// JSONReporter implements JSON reporting for validation results
type JSONReporter struct {
	output *os.File
}

// NewJSONReporter creates a new JSON reporter
func NewJSONReporter(output *os.File) *JSONReporter {
	if output == nil {
		output = os.Stdout
	}
	return &JSONReporter{
		output: output,
	}
}

// Report reports validation results as JSON
func (r *JSONReporter) Report(result *types.ValidationResult) error {
	// Marshal the result to JSON
	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal validation result: %w", err)
	}

	// Write to output
	_, err = r.output.Write(jsonData)
	if err != nil {
		return fmt.Errorf("failed to write JSON output: %w", err)
	}

	// Add newline
	_, err = r.output.WriteString("\n")
	return err
}

// GetName returns the reporter name
func (r *JSONReporter) GetName() string {
	return "json"
}
