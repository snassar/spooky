package utilities

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"spooky/internal/logging"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/pkg/errors"
)

// HCLValidator provides functionality to validate HCL files and content
type HCLValidator struct {
	// Configuration options for validation
	allowUnknownBlocks bool
	strictMode         bool
}

// ValidationResult contains the result of HCL validation
type ValidationResult struct {
	IsValid    bool
	Errors     []ValidationError
	Warnings   []ValidationWarning
	FileSize   int64
	BlockCount int
}

// ValidationError represents a validation error
type ValidationError struct {
	Message  string
	Line     int
	Column   int
	File     string
	Severity string // "error", "warning"
}

// ValidationWarning represents a validation warning
type ValidationWarning struct {
	Message string
	Line    int
	Column  int
	File    string
}

// NewHCLValidator creates a new HCL validator with default settings
func NewHCLValidator() *HCLValidator {
	return &HCLValidator{
		allowUnknownBlocks: false,
		strictMode:         true,
	}
}

// WithAllowUnknownBlocks configures the validator to allow unknown blocks
func (v *HCLValidator) WithAllowUnknownBlocks(allow bool) *HCLValidator {
	v.allowUnknownBlocks = allow
	return v
}

// WithStrictMode configures the validator to use strict mode
func (v *HCLValidator) WithStrictMode(strict bool) *HCLValidator {
	v.strictMode = strict
	return v
}

// ValidateFile validates an HCL file at the given path
func (v *HCLValidator) ValidateFile(filePath string) (*ValidationResult, error) {
	logger := logging.GetGlobalLogger()

	logger.Debug("validating HCL file",
		slog.String("component", "hcl_validator"),
		slog.String("operation", "validate_file"),
		slog.String("file_path", filePath))

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		logger.Error("file does not exist",
			slog.String("component", "hcl_validator"),
			slog.String("operation", "validate_file"),
			slog.String("file_path", filePath))
		return &ValidationResult{
			IsValid: false,
			Errors: []ValidationError{
				{
					Message:  fmt.Sprintf("File does not exist: %s", filePath),
					Severity: "error",
				},
			},
		}, NewHCLFileError(filePath, "stat", "file does not exist")
	}

	// Check file extension
	if !strings.HasSuffix(strings.ToLower(filePath), ".hcl") {
		return &ValidationResult{
			IsValid: false,
			Errors: []ValidationError{
				{
					Message:  fmt.Sprintf("File does not have .hcl extension: %s", filePath),
					Severity: "error",
				},
			},
		}, NewHCLFileError(filePath, "validate", "file does not have .hcl extension")
	}

	// Read file content
	content, err := os.ReadFile(filePath)
	if err != nil {
		return &ValidationResult{
			IsValid: false,
			Errors: []ValidationError{
				{
					Message:  fmt.Sprintf("Failed to read file: %v", err),
					Severity: "error",
				},
			},
		}, NewHCLFileError(filePath, "read", err.Error())
	}

	// Validate content
	return v.ValidateContent(string(content), filePath)
}

// ValidateContent validates HCL content as a string
func (v *HCLValidator) ValidateContent(content, fileName string) (*ValidationResult, error) {
	result := &ValidationResult{
		IsValid:    true,
		Errors:     []ValidationError{},
		Warnings:   []ValidationWarning{},
		FileSize:   int64(len(content)),
		BlockCount: 0,
	}

	// Parse HCL content
	_, diags := hclsyntax.ParseConfig([]byte(content), fileName, hcl.Pos{Line: 1, Column: 1})
	if diags.HasErrors() {
		result.IsValid = false
		for _, diag := range diags.Errs() {
			result.Errors = append(result.Errors, ValidationError{
				Message:  diag.Error(),
				Line:     0, // We'll get this from the error message if needed
				Column:   0,
				File:     fileName,
				Severity: "error",
			})
		}
		return result, nil
	}

	// Basic validation passed
	result.BlockCount = 1 // Simplified for now

	return result, nil
}

// ValidateDirectory validates all HCL files in a directory
func (v *HCLValidator) ValidateDirectory(dirPath string) (map[string]*ValidationResult, error) {
	results := make(map[string]*ValidationResult)

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return errors.Wrapf(err, "failed to walk directory %s", path)
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Only process .hcl files
		if !strings.HasSuffix(strings.ToLower(path), ".hcl") {
			return nil
		}

		// Validate the file
		result, err := v.ValidateFile(path)
		if err != nil {
			return errors.Wrapf(err, "failed to validate file %s", path)
		}

		// Use relative path as key
		relPath, err := filepath.Rel(dirPath, path)
		if err != nil {
			relPath = path
		}
		results[relPath] = result

		return nil
	})

	return results, errors.Wrap(err, "failed to validate directory")
}

// FormatValidationResult formats a validation result for display
func FormatValidationResult(result *ValidationResult, fileName string) string {
	var output strings.Builder

	output.WriteString(fmt.Sprintf("📄 Validation Results for: %s\n", fileName))
	output.WriteString(fmt.Sprintf("📊 File Size: %d bytes\n", result.FileSize))
	output.WriteString(fmt.Sprintf("📦 Block Count: %d\n", result.BlockCount))

	if result.IsValid {
		output.WriteString("✅ File is valid HCL\n")
	} else {
		output.WriteString("❌ File has validation errors\n")
	}

	if len(result.Errors) > 0 {
		output.WriteString(fmt.Sprintf("\n🚨 Errors (%d):\n", len(result.Errors)))
		for i, err := range result.Errors {
			if err.Line > 0 {
				output.WriteString(fmt.Sprintf("  %d. Line %d:%d - %s\n", i+1, err.Line, err.Column, err.Message))
			} else {
				output.WriteString(fmt.Sprintf("  %d. %s\n", i+1, err.Message))
			}
		}
	}

	if len(result.Warnings) > 0 {
		output.WriteString(fmt.Sprintf("\n⚠️  Warnings (%d):\n", len(result.Warnings)))
		for i, warning := range result.Warnings {
			if warning.Line > 0 {
				output.WriteString(fmt.Sprintf("  %d. Line %d:%d - %s\n", i+1, warning.Line, warning.Column, warning.Message))
			} else {
				output.WriteString(fmt.Sprintf("  %d. %s\n", i+1, warning.Message))
			}
		}
	}

	return output.String()
}

// IsValidHCL is a convenience function for quick validation
func IsValidHCL(content string) bool {
	validator := NewHCLValidator()
	result, err := validator.ValidateContent(content, "unknown")
	if err != nil {
		return false
	}
	return result.IsValid
}

// ValidateHCLFile performs quick validation of an HCL file for syntax errors.
//
// Parameters:
//   - filePath: Path to the HCL file to validate
//
// Returns:
//   - bool: True if file is valid HCL, false if syntax errors found
//   - error: File system or parsing errors
//
// Dependencies: github.com/hashicorp/hcl/v2 for HCL parsing
//
// Example usage:
//
//	valid, err := ValidateHCLFile("config.hcl")
//	if err != nil {
//	    return fmt.Errorf("validation failed: %w", err)
//	}
//	if !valid {
//	    fmt.Println("HCL file has syntax errors")
//	}
//
// Performance: ~1-10ms for typical HCL files
func ValidateHCLFile(filePath string) (bool, error) {
	validator := NewHCLValidator()
	result, err := validator.ValidateFile(filePath)
	if err != nil {
		return false, err
	}
	return result.IsValid, nil
}
