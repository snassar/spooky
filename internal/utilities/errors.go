package utilities

import "github.com/pkg/errors"

// HCL-specific error types
var (
	ErrHCLSyntaxError      = errors.New("HCL syntax error")
	ErrHCLFileNotFound     = errors.New("HCL file not found")
	ErrHCLInvalidBlock     = errors.New("invalid HCL block")
	ErrHCLValidationFailed = errors.New("HCL validation failed")
	ErrHCLParseFailed      = errors.New("failed to parse HCL")
	ErrHCLReadFailed       = errors.New("failed to read HCL file")
	ErrHCLDirectoryError   = errors.New("HCL directory error")
)

// HCLSyntaxError provides detailed syntax error information
type HCLSyntaxError struct {
	File    string
	Line    int
	Column  int
	Message string
	Token   string
}

func (e *HCLSyntaxError) Error() string {
	if e.Line > 0 {
		return errors.Errorf("HCL syntax error at %s:%d:%d: %s",
			e.File, e.Line, e.Column, e.Message).Error()
	}
	return errors.Errorf("HCL syntax error in %s: %s", e.File, e.Message).Error()
}

func (e *HCLSyntaxError) Unwrap() error {
	return ErrHCLSyntaxError
}

// HCLValidationError provides structured validation error information
type HCLValidationError struct {
	File      string
	BlockType string
	Field     string
	Value     interface{}
	Rule      string
	Message   string
	Severity  string
}

func (e *HCLValidationError) Error() string {
	return errors.Errorf("HCL validation error in %s: %s", e.File, e.Message).Error()
}

func (e *HCLValidationError) Unwrap() error {
	return ErrHCLValidationFailed
}

// HCLFileError provides file-related error information
type HCLFileError struct {
	FilePath  string
	Operation string
	Message   string
}

func (e *HCLFileError) Error() string {
	return errors.Errorf("HCL file error for %s during %s: %s",
		e.FilePath, e.Operation, e.Message).Error()
}

func (e *HCLFileError) Unwrap() error {
	return ErrHCLFileNotFound
}

// HCLBlockError provides block-related error information
type HCLBlockError struct {
	File      string
	BlockType string
	BlockName string
	Message   string
}

func (e *HCLBlockError) Error() string {
	return errors.Errorf("HCL block error in %s: %s block '%s': %s",
		e.File, e.BlockType, e.BlockName, e.Message).Error()
}

func (e *HCLBlockError) Unwrap() error {
	return ErrHCLInvalidBlock
}

// HCLDirectoryError provides directory-related error information
type HCLDirectoryError struct {
	Directory string
	Operation string
	Message   string
	FileCount int
}

func (e *HCLDirectoryError) Error() string {
	return errors.Errorf("HCL directory error for %s during %s: %s",
		e.Directory, e.Operation, e.Message).Error()
}

func (e *HCLDirectoryError) Unwrap() error {
	return ErrHCLDirectoryError
}

// Helper functions for creating typed errors
func NewHCLSyntaxError(file, message string, line, column int, token string) error {
	return errors.WithStack(&HCLSyntaxError{
		File:    file,
		Line:    line,
		Column:  column,
		Message: message,
		Token:   token,
	})
}

func NewHCLValidationError(file, blockType, field string, value interface{}, rule, message, severity string) error {
	return errors.WithStack(&HCLValidationError{
		File:      file,
		BlockType: blockType,
		Field:     field,
		Value:     value,
		Rule:      rule,
		Message:   message,
		Severity:  severity,
	})
}

func NewHCLFileError(filePath, operation, message string) error {
	return errors.WithStack(&HCLFileError{
		FilePath:  filePath,
		Operation: operation,
		Message:   message,
	})
}

func NewHCLBlockError(file, blockType, blockName, message string) error {
	return errors.WithStack(&HCLBlockError{
		File:      file,
		BlockType: blockType,
		BlockName: blockName,
		Message:   message,
	})
}

func NewHCLDirectoryError(directory, operation, message string, fileCount int) error {
	return errors.WithStack(&HCLDirectoryError{
		Directory: directory,
		Operation: operation,
		Message:   message,
		FileCount: fileCount,
	})
}
