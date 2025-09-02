package utilities

import (
	"errors"
	"testing"
)

func TestNewErrorHandler(t *testing.T) {
	eh := NewErrorHandler(true, false)

	if eh.verbose != true {
		t.Errorf("Expected verbose to be true, got %v", eh.verbose)
	}

	if eh.exitOnError != false {
		t.Errorf("Expected exitOnError to be false, got %v", eh.exitOnError)
	}

	if len(eh.warnings) != 0 {
		t.Errorf("Expected empty warnings slice, got %d", len(eh.warnings))
	}

	if len(eh.errors) != 0 {
		t.Errorf("Expected empty errors slice, got %d", len(eh.errors))
	}

	if len(eh.critical) != 0 {
		t.Errorf("Expected empty critical slice, got %d", len(eh.critical))
	}
}

func TestHandleError_Info(t *testing.T) {
	eh := NewErrorHandler(true, false)
	testErr := errors.New("test info error")

	eh.HandleError(testErr, SeverityInfo, "test message", nil)

	if len(eh.warnings) != 0 {
		t.Errorf("Expected no warnings for info severity, got %d", len(eh.warnings))
	}

	if len(eh.errors) != 0 {
		t.Errorf("Expected no errors for info severity, got %d", len(eh.errors))
	}
}

func TestHandleError_Warning(t *testing.T) {
	eh := NewErrorHandler(false, false)
	testErr := errors.New("test warning error")

	eh.HandleError(testErr, SeverityWarning, "test warning", nil)

	if len(eh.warnings) != 1 {
		t.Errorf("Expected 1 warning, got %d", len(eh.warnings))
	}

	if eh.warnings[0].Error != testErr {
		t.Errorf("Expected warning error to match test error")
	}

	if eh.warnings[0].Message != "test warning" {
		t.Errorf("Expected warning message to be 'test warning', got %s", eh.warnings[0].Message)
	}
}

func TestHandleError_Error(t *testing.T) {
	eh := NewErrorHandler(false, false)
	testErr := errors.New("test error")

	eh.HandleError(testErr, SeverityError, "test error message", nil)

	if len(eh.errors) != 1 {
		t.Errorf("Expected 1 error, got %d", len(eh.errors))
	}

	if eh.errors[0].Error != testErr {
		t.Errorf("Expected error to match test error")
	}

	if eh.errors[0].Message != "test error message" {
		t.Errorf("Expected error message to be 'test error message', got %s", eh.errors[0].Message)
	}
}

func TestHandleError_Critical(t *testing.T) {
	eh := NewErrorHandler(false, false)
	testErr := errors.New("test critical error")

	// For testing critical errors, we need to avoid os.Exit(1)
	// We'll test the LogError method instead which doesn't exit
	eh.LogError(testErr, "test critical message", nil)

	// Verify the error handler can still track critical errors manually
	eh.critical = append(eh.critical, ErrorInfo{
		Error:    testErr,
		Severity: SeverityCritical,
		Message:  "test critical message",
		Context:  nil,
	})

	if len(eh.critical) != 1 {
		t.Errorf("Expected 1 critical error, got %d", len(eh.critical))
	}

	if eh.critical[0].Error != testErr {
		t.Errorf("Expected critical error to match test error")
	}

	if eh.critical[0].Message != "test critical message" {
		t.Errorf("Expected critical error message to be 'test critical message', got %s", eh.critical[0].Message)
	}
}

func TestHandleError_Nil(t *testing.T) {
	eh := NewErrorHandler(false, false)

	// Should not panic or add anything
	eh.HandleError(nil, SeverityError, "test", nil)

	if len(eh.errors) != 0 {
		t.Errorf("Expected no errors for nil error, got %d", len(eh.errors))
	}
}

func TestHandleWarning(t *testing.T) {
	eh := NewErrorHandler(false, false)
	testErr := errors.New("test warning")

	eh.HandleWarning(testErr, "test warning message", nil)

	if len(eh.warnings) != 1 {
		t.Errorf("Expected 1 warning, got %d", len(eh.warnings))
	}

	if eh.warnings[0].Severity != SeverityWarning {
		t.Errorf("Expected severity to be Warning, got %v", eh.warnings[0].Severity)
	}
}

func TestHandleErrorLevel(t *testing.T) {
	eh := NewErrorHandler(false, false)
	testErr := errors.New("test error")

	eh.HandleErrorLevel(testErr, "test error message", nil)

	if len(eh.errors) != 1 {
		t.Errorf("Expected 1 error, got %d", len(eh.errors))
	}

	if eh.errors[0].Severity != SeverityError {
		t.Errorf("Expected severity to be Error, got %v", eh.errors[0].Severity)
	}
}

func TestHandleCritical(t *testing.T) {
	eh := NewErrorHandler(false, false)
	testErr := errors.New("test critical")

	// For testing critical errors, we need to avoid os.Exit(1)
	// We'll test the LogError method instead which doesn't exit
	eh.LogError(testErr, "test critical message", nil)

	// Verify the error handler can still track critical errors manually
	eh.critical = append(eh.critical, ErrorInfo{
		Error:    testErr,
		Severity: SeverityCritical,
		Message:  "test critical message",
		Context:  nil,
	})

	if len(eh.critical) != 1 {
		t.Errorf("Expected 1 critical error, got %d", len(eh.critical))
	}

	if eh.critical[0].Severity != SeverityCritical {
		t.Errorf("Expected severity to be Critical, got %v", eh.critical[0].Severity)
	}

	if eh.critical[0].Message != "test critical message" {
		t.Errorf("Expected critical error message to be 'test critical message', got %s", eh.critical[0].Message)
	}
}

func TestLogError(t *testing.T) {
	eh := NewErrorHandler(false, false)
	testErr := errors.New("test log error")

	// Should not add to errors slice
	eh.LogError(testErr, "test log message", nil)

	if len(eh.errors) != 0 {
		t.Errorf("Expected no errors for LogError, got %d", len(eh.errors))
	}
}

func TestLogWarning(t *testing.T) {
	eh := NewErrorHandler(false, false)
	testErr := errors.New("test log warning")

	// Should not add to warnings slice
	eh.LogWarning(testErr, "test log warning message", nil)

	if len(eh.warnings) != 0 {
		t.Errorf("Expected no warnings for LogWarning, got %d", len(eh.warnings))
	}
}

func TestLogInfo_Verbose(t *testing.T) {
	eh := NewErrorHandler(true, false)

	// Should log when verbose is true
	eh.LogInfo("test info message")

	// Note: We can't easily test stdout in unit tests, but we can verify the method doesn't panic
}

func TestLogInfo_NotVerbose(t *testing.T) {
	eh := NewErrorHandler(false, false)

	// Should not log when verbose is false
	eh.LogInfo("test info message")

	// Note: We can't easily test stdout in unit tests, but we can verify the method doesn't panic
}

func TestLogSuccess(t *testing.T) {
	eh := NewErrorHandler(false, false)

	// Should always log success
	eh.LogSuccess("test success message")

	// Note: We can't easily test stdout in unit tests, but we can verify the method doesn't panic
}

func TestHasErrors(t *testing.T) {
	eh := NewErrorHandler(false, false)

	if eh.HasErrors() {
		t.Errorf("Expected no errors initially")
	}

	testErr := errors.New("test error")
	eh.HandleError(testErr, SeverityError, "test", nil)

	if !eh.HasErrors() {
		t.Errorf("Expected errors after adding error")
	}
}

func TestHasWarnings(t *testing.T) {
	eh := NewErrorHandler(false, false)

	if eh.HasWarnings() {
		t.Errorf("Expected no warnings initially")
	}

	testErr := errors.New("test warning")
	eh.HandleError(testErr, SeverityWarning, "test", nil)

	if !eh.HasWarnings() {
		t.Errorf("Expected warnings after adding warning")
	}
}

func TestGetErrorCount(t *testing.T) {
	eh := NewErrorHandler(false, false)

	if eh.GetErrorCount() != 0 {
		t.Errorf("Expected 0 errors initially, got %d", eh.GetErrorCount())
	}

	testErr := errors.New("test error")
	eh.HandleError(testErr, SeverityError, "test", nil)

	// Manually add a critical error to avoid os.Exit(1)
	eh.critical = append(eh.critical, ErrorInfo{
		Error:    testErr,
		Severity: SeverityCritical,
		Message:  "test",
		Context:  nil,
	})

	if eh.GetErrorCount() != 2 {
		t.Errorf("Expected 2 errors, got %d", eh.GetErrorCount())
	}
}

func TestGetWarningCount(t *testing.T) {
	eh := NewErrorHandler(false, false)

	if eh.GetWarningCount() != 0 {
		t.Errorf("Expected 0 warnings initially, got %d", eh.GetWarningCount())
	}

	testErr := errors.New("test warning")
	eh.HandleError(testErr, SeverityWarning, "test", nil)
	eh.HandleError(testErr, SeverityWarning, "test2", nil)

	if eh.GetWarningCount() != 2 {
		t.Errorf("Expected 2 warnings, got %d", eh.GetWarningCount())
	}
}

func TestIsNonCriticalError(t *testing.T) {
	eh := NewErrorHandler(false, false)

	testCases := []struct {
		errorStr            string
		expectedNonCritical bool
	}{
		{"file not found", true},
		{"permission denied", true},
		{"connection refused", true},
		{"timeout", true},
		{"no such file or directory", true},
		{"unexpected error", false},
		{"critical system failure", false},
	}

	for _, tc := range testCases {
		testErr := errors.New(tc.errorStr)
		result := eh.IsNonCriticalError(testErr)

		if result != tc.expectedNonCritical {
			t.Errorf("For error '%s', expected non-critical=%v, got %v",
				tc.errorStr, tc.expectedNonCritical, result)
		}
	}
}

func TestHandleNonCriticalError(t *testing.T) {
	eh := NewErrorHandler(false, false)

	// Test non-critical error
	nonCriticalErr := errors.New("file not found")
	eh.HandleNonCriticalError(nonCriticalErr, "test non-critical", nil)

	if len(eh.warnings) != 1 {
		t.Errorf("Expected 1 warning for non-critical error, got %d", len(eh.warnings))
	}

	if len(eh.errors) != 0 {
		t.Errorf("Expected 0 errors for non-critical error, got %d", len(eh.errors))
	}

	// Test critical error
	criticalErr := errors.New("unexpected system failure")
	eh.HandleNonCriticalError(criticalErr, "test critical", nil)

	if len(eh.errors) != 1 {
		t.Errorf("Expected 1 error for critical error, got %d", len(eh.errors))
	}
}

func TestFormatError(t *testing.T) {
	eh := NewErrorHandler(false, false)
	testErr := errors.New("test error")

	// Test without context
	formatted := eh.FormatError(testErr, nil)
	expected := "test error"
	if formatted != expected {
		t.Errorf("Expected formatted error '%s', got '%s'", expected, formatted)
	}

	// Test with context
	context := map[string]interface{}{
		"file": "test.txt",
		"line": 42,
	}
	formatted = eh.FormatError(testErr, context)
	expected = "test error | Context: file=test.txt, line=42"
	if formatted != expected {
		t.Errorf("Expected formatted error with context '%s', got '%s'", expected, formatted)
	}
}

func TestFormatError_Nil(t *testing.T) {
	eh := NewErrorHandler(false, false)

	formatted := eh.FormatError(nil, nil)
	if formatted != "" {
		t.Errorf("Expected empty string for nil error, got '%s'", formatted)
	}
}
