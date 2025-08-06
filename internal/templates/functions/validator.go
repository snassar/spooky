package functions

import (
	"fmt"
	"reflect"
)

// Validator implements FunctionValidator interface
type Validator struct{}

// NewFunctionValidator creates a new function validator
func NewFunctionValidator() FunctionValidator {
	return &Validator{}
}

// ValidateFunction validates a function
func (v *Validator) ValidateFunction(name string, fn interface{}) error {
	if fn == nil {
		return fmt.Errorf("function '%s' is nil", name)
	}

	// Check if it's a function
	fnType := reflect.TypeOf(fn)
	if fnType.Kind() != reflect.Func {
		return fmt.Errorf("function '%s' is not a function type", name)
	}

	// Validate function signature
	if err := v.ValidateFunctionSignature(fn); err != nil {
		return fmt.Errorf("function '%s' signature validation failed: %w", name, err)
	}

	return nil
}

// ValidateFunctionSignature validates function signature
func (v *Validator) ValidateFunctionSignature(fn interface{}) error {
	fnType := reflect.TypeOf(fn)

	// Check number of inputs (should be reasonable)
	if fnType.NumIn() > 10 {
		return fmt.Errorf("function has too many input parameters (%d)", fnType.NumIn())
	}

	// Check number of outputs (should be 1-2 for template functions)
	if fnType.NumOut() < 1 || fnType.NumOut() > 2 {
		return fmt.Errorf("function has invalid number of outputs (%d)", fnType.NumOut())
	}

	// Check output types (second output should be error if present)
	if fnType.NumOut() == 2 {
		errorType := reflect.TypeOf((*error)(nil)).Elem()
		if !fnType.Out(1).Implements(errorType) {
			return fmt.Errorf("function second output must be error type")
		}
	}

	return nil
}

// ValidateFunctionReturnType validates function return type
func (v *Validator) ValidateFunctionReturnType(fn interface{}) error {
	fnType := reflect.TypeOf(fn)

	// Check first return type (should be a valid template function return type)
	firstReturnType := fnType.Out(0)
	validTypes := []reflect.Kind{
		reflect.String,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64,
		reflect.Bool,
		reflect.Interface,
		reflect.Map,
		reflect.Slice,
	}

	isValid := false
	for _, validType := range validTypes {
		if firstReturnType.Kind() == validType {
			isValid = true
			break
		}
	}

	if !isValid {
		return fmt.Errorf("function return type '%s' is not valid for template functions", firstReturnType.Kind())
	}

	return nil
}
