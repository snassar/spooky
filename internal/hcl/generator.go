package hcl

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

// Generator is a generic utility for converting Go structs to HCL format
type Generator struct {
	// Configuration options
	UseDefaults bool // Whether to include fields with default values
	IndentSize  int  // Number of spaces for indentation

	// Shared utilities
	sharedUtils *SharedHCLUtils
}

// NewGenerator creates a new HCL generator with default settings
func NewGenerator() *Generator {
	return &Generator{
		UseDefaults: true,
		IndentSize:  2,
		sharedUtils: NewSharedHCLUtils(),
	}
}

// ToHCL converts any Go struct to HCL string with a custom block name
func (hg *Generator) ToHCL(config interface{}, blockName string) (string, error) {
	// Convert the config to HCL using the HCL library
	hclFile := hclwrite.NewEmptyFile()

	// Convert the config to cty.Value
	ctyValue, err := hg.structToCty(config)
	if err != nil {
		return "", fmt.Errorf("failed to convert config to cty.Value: %v", err)
	}

	// Create the root body
	rootBody := hclFile.Body()

	// Convert to HCL block
	err = hg.ctyValueToHCL(ctyValue, rootBody, blockName)
	if err != nil {
		return "", fmt.Errorf("failed to convert cty.Value to HCL: %v", err)
	}

	return string(hclFile.Bytes()), nil
}

// structToCty converts a Go struct to cty.Value, respecting schema tags
func (hg *Generator) structToCty(config interface{}) (cty.Value, error) {
	v := reflect.ValueOf(config)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return cty.NilVal, fmt.Errorf("config must be a struct, got %s", v.Kind())
	}

	// Convert struct to map[string]cty.Value
	valueMap := make(map[string]cty.Value)
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		// Get the JSON tag name
		jsonTag := fieldType.Tag.Get("json")
		if jsonTag == "" || jsonTag == "-" {
			continue
		}

		// Handle comma-separated JSON tags (e.g., "basic_facts,omitempty")
		if commaIndex := strings.Index(jsonTag, ","); commaIndex != -1 {
			jsonTag = jsonTag[:commaIndex]
		}

		// Convert field value to cty.Value
		ctyVal, err := hg.fieldToCty(field, fieldType)
		if err != nil {
			return cty.NilVal, fmt.Errorf("failed to convert field %s: %v", jsonTag, err)
		}

		// Include field if it has a value or if we're using defaults
		if !ctyVal.IsNull() || hg.shouldIncludeField(field, fieldType) {
			valueMap[jsonTag] = ctyVal
		}
	}

	return cty.ObjectVal(valueMap), nil
}

// fieldToCty converts a struct field to cty.Value, respecting schema tags
func (hg *Generator) fieldToCty(field reflect.Value, fieldType reflect.StructField) (cty.Value, error) {
	switch field.Kind() {
	case reflect.Ptr:
		return hg.handlePointerField(field, fieldType)
	case reflect.Interface:
		return hg.handleInterfaceField(field, fieldType)
	case reflect.String:
		return hg.handleStringField(field, fieldType)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return hg.handleIntField(field, fieldType)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return hg.handleUintField(field, fieldType)
	case reflect.Float32, reflect.Float64:
		return hg.handleFloatField(field, fieldType)
	case reflect.Bool:
		return hg.handleBoolField(field, fieldType)
	case reflect.Struct:
		return hg.handleStructField(field)
	case reflect.Slice:
		return hg.handleSliceField(field, fieldType)
	case reflect.Map:
		return hg.handleMapField(field, fieldType)
	default:
		return cty.NilVal, fmt.Errorf("unsupported field type: %v", field.Kind())
	}
}

// handlePointerField handles pointer type fields
func (hg *Generator) handlePointerField(field reflect.Value, fieldType reflect.StructField) (cty.Value, error) {
	if field.IsNil() {
		// Check if there's a default value in the tag
		if defaultValue := fieldType.Tag.Get("default"); defaultValue != "" {
			return hg.parseDefaultValue(defaultValue, fieldType)
		}
		return cty.NilVal, nil
	}
	return hg.fieldToCty(field.Elem(), fieldType)
}

// handleInterfaceField handles interface{} type fields by examining the actual value
func (hg *Generator) handleInterfaceField(field reflect.Value, fieldType reflect.StructField) (cty.Value, error) {
	if field.IsNil() {
		// Check if there's a default value in the tag
		if defaultValue := fieldType.Tag.Get("default"); defaultValue != "" {
			return hg.parseDefaultValue(defaultValue, fieldType)
		}
		return cty.NilVal, nil
	}

	// Get the actual value stored in the interface
	actualValue := field.Elem()

	// Handle the actual type of the value
	switch actualValue.Kind() {
	case reflect.String:
		return cty.StringVal(actualValue.String()), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return cty.NumberIntVal(actualValue.Int()), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return cty.NumberUIntVal(actualValue.Uint()), nil
	case reflect.Float32, reflect.Float64:
		return cty.NumberFloatVal(actualValue.Float()), nil
	case reflect.Bool:
		return cty.BoolVal(actualValue.Bool()), nil
	case reflect.Slice:
		return hg.handleSliceField(actualValue, fieldType)
	case reflect.Map:
		return hg.handleMapField(actualValue, fieldType)
	case reflect.Struct:
		return hg.handleStructField(actualValue)
	case reflect.Ptr:
		if actualValue.IsNil() {
			return cty.NilVal, nil
		}
		return hg.handleInterfaceField(actualValue, fieldType)
	default:
		// For unsupported types, try to convert to string as a fallback
		return cty.StringVal(fmt.Sprintf("%v", actualValue.Interface())), nil
	}
}

// handleStringField handles string type fields
func (hg *Generator) handleStringField(field reflect.Value, fieldType reflect.StructField) (cty.Value, error) {
	value := field.String()
	if value == "" {
		// Check for default value
		if defaultValue := fieldType.Tag.Get("default"); defaultValue != "" {
			return cty.StringVal(defaultValue), nil
		}
	}
	return cty.StringVal(value), nil
}

// handleIntField handles integer type fields
func (hg *Generator) handleIntField(field reflect.Value, fieldType reflect.StructField) (cty.Value, error) {
	value := field.Int()
	if value == 0 {
		// Check for default value
		if defaultValue := fieldType.Tag.Get("default"); defaultValue != "" {
			return hg.parseDefaultValue(defaultValue, fieldType)
		}
	}
	return cty.NumberIntVal(value), nil
}

// handleUintField handles unsigned integer type fields
func (hg *Generator) handleUintField(field reflect.Value, fieldType reflect.StructField) (cty.Value, error) {
	value := field.Uint()
	if value == 0 {
		// Check for default value
		if defaultValue := fieldType.Tag.Get("default"); defaultValue != "" {
			return hg.parseDefaultValue(defaultValue, fieldType)
		}
	}
	return cty.NumberUIntVal(value), nil
}

// handleFloatField handles float type fields
func (hg *Generator) handleFloatField(field reflect.Value, fieldType reflect.StructField) (cty.Value, error) {
	value := field.Float()
	if value == 0 {
		// Check for default value
		if defaultValue := fieldType.Tag.Get("default"); defaultValue != "" {
			return hg.parseDefaultValue(defaultValue, fieldType)
		}
	}
	return cty.NumberFloatVal(value), nil
}

// handleBoolField handles boolean type fields
func (hg *Generator) handleBoolField(field reflect.Value, fieldType reflect.StructField) (cty.Value, error) {
	value := field.Bool()
	if !value {
		// Check for default value
		if defaultValue := fieldType.Tag.Get("default"); defaultValue != "" {
			return hg.parseDefaultValue(defaultValue, fieldType)
		}
	}
	return cty.BoolVal(value), nil
}

// handleStructField handles struct type fields
func (hg *Generator) handleStructField(field reflect.Value) (cty.Value, error) {
	// Handle nested structs
	if field.CanInterface() {
		return hg.structToCty(field.Interface())
	}
	return cty.NilVal, fmt.Errorf("cannot convert unexported struct field")
}

// handleSliceField handles slice type fields
func (hg *Generator) handleSliceField(field reflect.Value, fieldType reflect.StructField) (cty.Value, error) {
	if field.Len() == 0 {
		return cty.ListValEmpty(cty.String), nil
	}

	values := make([]cty.Value, field.Len())
	for i := 0; i < field.Len(); i++ {
		val, err := hg.fieldToCty(field.Index(i), fieldType)
		if err != nil {
			return cty.NilVal, err
		}
		values[i] = val
	}
	return cty.ListVal(values), nil
}

// handleMapField handles map type fields
func (hg *Generator) handleMapField(field reflect.Value, fieldType reflect.StructField) (cty.Value, error) {
	if field.Len() == 0 {
		return cty.MapValEmpty(cty.String), nil
	}

	valueMap := make(map[string]cty.Value)
	iter := field.MapRange()
	for iter.Next() {
		key := iter.Key().String()
		val, err := hg.fieldToCty(iter.Value(), fieldType)
		if err != nil {
			return cty.NilVal, err
		}
		valueMap[key] = val
	}
	return cty.MapVal(valueMap), nil
}

// parseDefaultValue parses a default value string based on the field type
func (hg *Generator) parseDefaultValue(defaultValue string, fieldType reflect.StructField) (cty.Value, error) {
	// Handle special default values
	switch defaultValue {
	case "true":
		return cty.BoolVal(true), nil
	case "false":
		return cty.BoolVal(false), nil
	case "null", "":
		return cty.NilVal, nil
	}

	// Parse based on field type
	switch fieldType.Type.Kind() {
	case reflect.String:
		return cty.StringVal(defaultValue), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if i, err := strconv.ParseInt(defaultValue, 10, 64); err == nil {
			return cty.NumberIntVal(i), nil
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if u, err := strconv.ParseUint(defaultValue, 10, 64); err == nil {
			return cty.NumberUIntVal(u), nil
		}
	case reflect.Float32, reflect.Float64:
		if f, err := strconv.ParseFloat(defaultValue, 64); err == nil {
			return cty.NumberFloatVal(f), nil
		}
	case reflect.Bool:
		if b, err := strconv.ParseBool(defaultValue); err == nil {
			return cty.BoolVal(b), nil
		}
	}

	// If we can't parse it, treat as string
	return cty.StringVal(defaultValue), nil
}

// shouldIncludeField determines if a field should be included in HCL output
func (hg *Generator) shouldIncludeField(field reflect.Value, fieldType reflect.StructField) bool {
	if !hg.UseDefaults {
		return false
	}

	// Check if field has a default value
	if defaultValue := fieldType.Tag.Get("default"); defaultValue != "" {
		return true
	}

	// Check if field is required
	if required := fieldType.Tag.Get("required"); required == "true" {
		return true
	}

	// Check if field has a non-zero value (for primitive types)
	switch field.Kind() {
	case reflect.String:
		return field.String() != ""
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return field.Int() != 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return field.Uint() != 0
	case reflect.Float32, reflect.Float64:
		return field.Float() != 0
	case reflect.Bool:
		return field.Bool()
	case reflect.Slice, reflect.Map:
		return field.Len() > 0
	case reflect.Ptr, reflect.Interface:
		return !field.IsNil()
	}

	return false
}

// ctyValueToHCL converts a cty.Value to HCL and writes it to the body
func (hg *Generator) ctyValueToHCL(value cty.Value, body *hclwrite.Body, blockName string) error {
	return hg.sharedUtils.CtyValueToHCL(value, body, blockName)
}

// GenerateHCL is a convenience function for quick HCL generation
func GenerateHCL(config interface{}, blockName string) (string, error) {
	generator := NewGenerator()
	return generator.ToHCL(config, blockName)
}

// GenerateHCLWithDefaults is a convenience function that includes default values
func GenerateHCLWithDefaults(config interface{}, blockName string) (string, error) {
	generator := NewGenerator()
	generator.UseDefaults = true
	return generator.ToHCL(config, blockName)
}

// GenerateHCLWithoutDefaults is a convenience function that excludes default values
func GenerateHCLWithoutDefaults(config interface{}, blockName string) (string, error) {
	generator := NewGenerator()
	generator.UseDefaults = false
	return generator.ToHCL(config, blockName)
}
