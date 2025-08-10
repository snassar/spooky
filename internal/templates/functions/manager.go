package functions

import (
	"fmt"
	"strings"
	"text/template"
	"time"

	spookyinterfaces "spooky/internal/interfaces"
	spookytypes "spooky/internal/types"
)

// Manager implements FunctionsManager interface
type Manager struct {
	config           *spookytypes.FunctionsConfig
	builtinFunctions template.FuncMap
	customFunctions  template.FuncMap
	validator        spookyinterfaces.FunctionValidator
	logger           spookyinterfaces.Logger
}

// NewManager creates a new functions manager
func NewManager(config *spookytypes.FunctionsConfig, logger spookyinterfaces.Logger) *Manager {
	return &Manager{
		config:           config,
		builtinFunctions: make(template.FuncMap),
		customFunctions:  make(template.FuncMap),
		validator:        NewFunctionValidator(),
		logger:           logger,
	}
}

// GetBuiltinFunctions returns built-in template functions
// Aligns with template-functions.hcl allowed_functions list
func (m *Manager) GetBuiltinFunctions() template.FuncMap {
	functions := make(template.FuncMap)

	// Add built-in functions from template-functions.hcl schema
	functions["custom"] = func(_ string) interface{} { return "mock-value" }
	functions["system"] = func(_ string) interface{} { return "mock-value" }
	functions["env"] = func(_ string) string { return "mock-value" }
	functions["data"] = func() map[string]interface{} { return make(map[string]interface{}) }
	functions["var"] = func(_ string) interface{} { return "mock-value" }
	functions["varOrDefault"] = func(_ string, _ interface{}) interface{} { return "mock-value" }
	functions["project"] = func() interface{} { return nil }
	functions["projectName"] = func() string { return "mock-project" }
	functions["machines"] = func() []interface{} { return []interface{}{} }
	functions["machine"] = func(_ string) interface{} { return nil }
	functions["facts"] = func() map[string]interface{} { return make(map[string]interface{}) }
	functions["fact"] = func(_ string) interface{} { return "mock-value" }

	// Enhanced functions
	functions["upper"] = strings.ToUpper
	functions["lower"] = strings.ToLower
	functions["trim"] = strings.TrimSpace
	functions["add"] = func(a, b int) int { return a + b }
	functions["sub"] = func(a, b int) int { return a - b }
	functions["mul"] = func(a, b int) int { return a * b }
	functions["div"] = func(a, b int) int { return a / b }

	return functions
}

// RegisterCustomFunction registers a custom function
func (m *Manager) RegisterCustomFunction(name string, fn interface{}) error {
	// 1. Validate function against template-functions.hcl schema
	if err := m.ValidateFunction(name, fn); err != nil {
		return fmt.Errorf("function validation failed: %w", err)
	}

	// 2. Register function
	m.customFunctions[name] = fn

	return nil
}

// ValidateFunction validates a function
func (m *Manager) ValidateFunction(name string, fn interface{}) error {
	// Validate function signature
	if err := m.validator.ValidateFunction(name, fn); err != nil {
		return fmt.Errorf("function signature validation failed: %w", err)
	}

	return nil
}

// GetFunction gets a function by name
func (m *Manager) GetFunction(name string) (interface{}, bool) {
	// Check custom functions first
	if fn, exists := m.customFunctions[name]; exists {
		return fn, true
	}

	// Check builtin functions
	if fn, exists := m.builtinFunctions[name]; exists {
		return fn, true
	}

	return nil, false
}

// ListFunctions lists all available functions
func (m *Manager) ListFunctions() []string {
	var functions []string

	// Add builtin functions
	for name := range m.builtinFunctions {
		functions = append(functions, name)
	}

	// Add custom functions
	for name := range m.customFunctions {
		functions = append(functions, name)
	}

	return functions
}

// RemoveFunction removes a custom function
func (m *Manager) RemoveFunction(name string) error {
	if _, exists := m.customFunctions[name]; !exists {
		return fmt.Errorf("function '%s' not found", name)
	}

	delete(m.customFunctions, name)
	return nil
}

// EnableBuiltinFunctions enables or disables builtin functions
func (m *Manager) EnableBuiltinFunctions(enabled bool) error {
	if m.config == nil {
		m.config = &spookytypes.FunctionsConfig{}
	}
	m.config.BuiltinFunctions = enabled
	return nil
}

// SetFunctionTimeout sets function timeout
func (m *Manager) SetFunctionTimeout(timeout time.Duration) error {
	if m.config == nil {
		m.config = &spookytypes.FunctionsConfig{}
	}
	m.config.FunctionTimeout = timeout
	return nil
}

// Close closes the functions manager
func (m *Manager) Close() error {
	// Cleanup resources if needed
	return nil
}
