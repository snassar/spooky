package reporters

import (
	"fmt"
	"os"
	"spooky/internal/schemas/types"
)

// ConsoleReporter implements console reporting for validation results
type ConsoleReporter struct{}

// NewConsoleReporter creates a new console reporter
func NewConsoleReporter() *ConsoleReporter {
	return &ConsoleReporter{}
}

// Report reports validation results to the console
func (r *ConsoleReporter) Report(result *types.ValidationResult) error {
	if result.Valid {
		fmt.Fprintf(os.Stdout, "✅ Validation passed for %s\n", result.Schema)
		return nil
	}

	fmt.Fprintf(os.Stderr, "❌ Validation failed for %s\n", result.Schema)

	// Report errors
	for _, err := range result.Errors {
		fmt.Fprintf(os.Stderr, "  Error: %s:%d:%d: %s: %s\n",
			err.File, err.Line, err.Column, err.Field, err.Message)
	}

	// Report warnings
	for _, warning := range result.Warnings {
		fmt.Fprintf(os.Stderr, "  Warning: %s:%d:%d: %s: %s\n",
			warning.File, warning.Line, warning.Column, warning.Field, warning.Message)
	}

	return nil
}

// GetName returns the reporter name
func (r *ConsoleReporter) GetName() string {
	return "console"
}
