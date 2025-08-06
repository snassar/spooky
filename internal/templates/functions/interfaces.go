package functions

import (
	"text/template"
	"time"
)

// FunctionsManager defines the interface for template functions
type FunctionsManager interface {
	// Core functions operations
	GetBuiltinFunctions() template.FuncMap
	RegisterCustomFunction(name string, fn interface{}) error
	ValidateFunction(name string, fn interface{}) error

	// Function management
	GetFunction(name string) (interface{}, bool)
	ListFunctions() []string
	RemoveFunction(name string) error

	// Configuration
	EnableBuiltinFunctions(enabled bool) error
	SetFunctionTimeout(timeout time.Duration) error

	// Utility operations
	Close() error
}

// FunctionValidator defines the interface for function validation
type FunctionValidator interface {
	ValidateFunction(name string, fn interface{}) error
	ValidateFunctionSignature(fn interface{}) error
	ValidateFunctionReturnType(fn interface{}) error
}
