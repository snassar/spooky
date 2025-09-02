package utilities

import (
	"fmt"
	"os"
	"strings"
)

// ErrorSeverity represents the severity level of an error
type ErrorSeverity int

const (
	// SeverityInfo represents informational messages (non-fatal)
	SeverityInfo ErrorSeverity = iota
	// SeverityWarning represents warnings (non-fatal, but concerning)
	SeverityWarning
	// SeverityError represents errors (fatal for the current operation)
	SeverityError
	// SeverityCritical represents critical errors (fatal for the entire program)
	SeverityCritical
)

// ErrorInfo contains error details and severity information
type ErrorInfo struct {
	Error    error
	Severity ErrorSeverity
	Message  string
	Context  map[string]interface{}
	ExitCode int
}

// ErrorHandler provides standardized error handling for tools
type ErrorHandler struct {
	verbose     bool
	warnings    []ErrorInfo
	errors      []ErrorInfo
	critical    []ErrorInfo
	exitOnError bool
}

// NewErrorHandler creates a new error handler with the specified options
func NewErrorHandler(verbose bool, exitOnError bool) *ErrorHandler {
	return &ErrorHandler{
		verbose:     verbose,
		exitOnError: exitOnError,
		warnings:    make([]ErrorInfo, 0),
		errors:      make([]ErrorInfo, 0),
		critical:    make([]ErrorInfo, 0),
	}
}

// HandleError processes an error with the specified severity
func (eh *ErrorHandler) HandleError(err error, severity ErrorSeverity, message string, context map[string]interface{}) {
	if err == nil {
		return
	}

	errorInfo := ErrorInfo{
		Error:    err,
		Severity: severity,
		Message:  message,
		Context:  context,
	}

	switch severity {
	case SeverityInfo:
		if eh.verbose {
			fmt.Printf("ℹ️  %s: %v\n", message, err)
		}
	case SeverityWarning:
		eh.warnings = append(eh.warnings, errorInfo)
		fmt.Fprintf(os.Stderr, "⚠️  %s: %v\n", message, err)
	case SeverityError:
		eh.errors = append(eh.errors, errorInfo)
		fmt.Fprintf(os.Stderr, "❌ %s: %v\n", message, err)
		if eh.exitOnError {
			os.Exit(1)
		}
	case SeverityCritical:
		eh.critical = append(eh.critical, errorInfo)
		fmt.Fprintf(os.Stderr, "💥 %s: %v\n", message, err)
		os.Exit(1)
	}
}

// HandleWarning handles a warning-level error
func (eh *ErrorHandler) HandleWarning(err error, message string, context map[string]interface{}) {
	eh.HandleError(err, SeverityWarning, message, context)
}

// HandleError handles an error-level error
func (eh *ErrorHandler) HandleErrorLevel(err error, message string, context map[string]interface{}) {
	eh.HandleError(err, SeverityError, message, context)
}

// HandleCritical handles a critical-level error
func (eh *ErrorHandler) HandleCritical(err error, message string, context map[string]interface{}) {
	eh.HandleError(err, SeverityCritical, message, context)
}

// LogError logs an error without exiting (for non-critical errors)
func (eh *ErrorHandler) LogError(err error, message string, context map[string]interface{}) {
	if err == nil {
		return
	}
	fmt.Fprintf(os.Stderr, "❌ %s: %v\n", message, err)
}

// LogWarning logs a warning without exiting
func (eh *ErrorHandler) LogWarning(err error, message string, context map[string]interface{}) {
	if err == nil {
		return
	}
	fmt.Fprintf(os.Stderr, "⚠️  %s: %v\n", message, err)
}

// LogInfo logs an informational message
func (eh *ErrorHandler) LogInfo(message string, args ...interface{}) {
	if eh.verbose {
		fmt.Printf("ℹ️  %s\n", fmt.Sprintf(message, args...))
	}
}

// LogSuccess logs a success message
func (eh *ErrorHandler) LogSuccess(message string, args ...interface{}) {
	fmt.Printf("✅ %s\n", fmt.Sprintf(message, args...))
}

// HasErrors returns true if there are any errors or critical errors
func (eh *ErrorHandler) HasErrors() bool {
	return len(eh.errors) > 0 || len(eh.critical) > 0
}

// HasWarnings returns true if there are any warnings
func (eh *ErrorHandler) HasWarnings() bool {
	return len(eh.warnings) > 0
}

// GetErrorCount returns the total count of errors and critical errors
func (eh *ErrorHandler) GetErrorCount() int {
	return len(eh.errors) + len(eh.critical)
}

// GetWarningCount returns the count of warnings
func (eh *ErrorHandler) GetWarningCount() int {
	return len(eh.warnings)
}

// PrintSummary prints a summary of all errors and warnings
func (eh *ErrorHandler) PrintSummary() {
	if eh.HasWarnings() || eh.HasErrors() {
		fmt.Println("\n📊 Error Summary:")
		fmt.Println("=================")

		if eh.HasWarnings() {
			fmt.Printf("⚠️  Warnings: %d\n", eh.GetWarningCount())
		}

		if eh.HasErrors() {
			fmt.Printf("❌ Errors: %d\n", eh.GetErrorCount())
		}

		if eh.HasErrors() {
			fmt.Println("\n💡 Some operations completed with errors. Check the output above for details.")
		}
	}
}

// ExitWithCode exits with the appropriate exit code based on errors
func (eh *ErrorHandler) ExitWithCode() {
	if eh.HasErrors() {
		os.Exit(1)
	}
	if eh.HasWarnings() {
		os.Exit(0) // Warnings don't cause failure
	}
	os.Exit(0)
}

// FormatError formats an error with context information
func (eh *ErrorHandler) FormatError(err error, context map[string]interface{}) string {
	if err == nil {
		return ""
	}

	var parts []string
	parts = append(parts, err.Error())

	if context != nil && len(context) > 0 {
		var contextParts []string
		for key, value := range context {
			contextParts = append(contextParts, fmt.Sprintf("%s=%v", key, value))
		}
		if len(contextParts) > 0 {
			parts = append(parts, fmt.Sprintf("Context: %s", strings.Join(contextParts, ", ")))
		}
	}

	return strings.Join(parts, " | ")
}

// IsNonCriticalError checks if an error should be treated as non-critical
func (eh *ErrorHandler) IsNonCriticalError(err error) bool {
	if err == nil {
		return false
	}

	// Check for common non-critical errors
	errorStr := err.Error()
	nonCriticalPatterns := []string{
		"file not found",
		"permission denied",
		"connection refused",
		"timeout",
		"no such file or directory",
	}

	for _, pattern := range nonCriticalPatterns {
		if strings.Contains(strings.ToLower(errorStr), pattern) {
			return true
		}
	}

	return false
}

// HandleNonCriticalError handles an error that should not cause the program to exit
func (eh *ErrorHandler) HandleNonCriticalError(err error, message string, context map[string]interface{}) {
	if eh.IsNonCriticalError(err) {
		eh.HandleWarning(err, message, context)
	} else {
		eh.HandleErrorLevel(err, message, context)
	}
}
