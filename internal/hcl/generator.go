package hcl

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

// HCLGenerator is a generic utility for converting Go structs to HCL format
type HCLGenerator struct {
	// Configuration options
	UseDefaults bool // Whether to include fields with default values
	IndentSize  int  // Number of spaces for indentation
}

// NewHCLGenerator creates a new HCL generator with default settings
func NewHCLGenerator() *HCLGenerator {
	return &HCLGenerator{
		UseDefaults: true,
		IndentSize:  2,
	}
}

// ToHCL converts any Go struct to HCL string with a custom block name
func (hg *HCLGenerator) ToHCL(config interface{}, blockName string) (string, error) {
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
func (hg *HCLGenerator) structToCty(config interface{}) (cty.Value, error) {
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
func (hg *HCLGenerator) fieldToCty(field reflect.Value, fieldType reflect.StructField) (cty.Value, error) {
	switch field.Kind() {
	case reflect.Ptr:
		// Handle pointers
		if field.IsNil() {
			// Check if there's a default value in the tag
			if defaultValue := fieldType.Tag.Get("default"); defaultValue != "" {
				return hg.parseDefaultValue(defaultValue, fieldType)
			}
			return cty.NilVal, nil
		}
		return hg.fieldToCty(field.Elem(), fieldType)
	case reflect.String:
		value := field.String()
		if value == "" {
			// Check for default value
			if defaultValue := fieldType.Tag.Get("default"); defaultValue != "" {
				return cty.StringVal(defaultValue), nil
			}
		}
		return cty.StringVal(value), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		value := field.Int()
		if value == 0 {
			// Check for default value
			if defaultValue := fieldType.Tag.Get("default"); defaultValue != "" {
				return hg.parseDefaultValue(defaultValue, fieldType)
			}
		}
		return cty.NumberIntVal(value), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		value := field.Uint()
		if value == 0 {
			// Check for default value
			if defaultValue := fieldType.Tag.Get("default"); defaultValue != "" {
				return hg.parseDefaultValue(defaultValue, fieldType)
			}
		}
		return cty.NumberUIntVal(value), nil
	case reflect.Float32, reflect.Float64:
		value := field.Float()
		if value == 0 {
			// Check for default value
			if defaultValue := fieldType.Tag.Get("default"); defaultValue != "" {
				return hg.parseDefaultValue(defaultValue, fieldType)
			}
		}
		return cty.NumberFloatVal(value), nil
	case reflect.Bool:
		value := field.Bool()
		if !value {
			// Check for default value
			if defaultValue := fieldType.Tag.Get("default"); defaultValue != "" {
				return hg.parseDefaultValue(defaultValue, fieldType)
			}
		}
		return cty.BoolVal(value), nil
	case reflect.Struct:
		// Handle nested structs
		if field.CanInterface() {
			return hg.structToCty(field.Interface())
		}
		return cty.NilVal, fmt.Errorf("cannot convert unexported struct field")
	case reflect.Slice:
		// Handle slices
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
	case reflect.Map:
		// Handle maps
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
	case reflect.Interface:
		// Handle interface{} types
		if field.IsNil() {
			return cty.NilVal, nil
		}

		// Get the concrete value
		concreteValue := field.Elem()
		return hg.fieldToCty(concreteValue, fieldType)
	default:
		return cty.NilVal, fmt.Errorf("unsupported field type: %s", field.Kind())
	}
}

// parseDefaultValue parses a default value string based on the field type
func (hg *HCLGenerator) parseDefaultValue(defaultValue string, fieldType reflect.StructField) (cty.Value, error) {
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
func (hg *HCLGenerator) shouldIncludeField(field reflect.Value, fieldType reflect.StructField) bool {
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

	return false
}

// ctyValueToHCL converts a cty.Value to HCL and writes it to the body
func (hg *HCLGenerator) ctyValueToHCL(value cty.Value, body *hclwrite.Body, blockName string) error {
	if !value.IsKnown() {
		return fmt.Errorf("cannot convert unknown value to HCL")
	}

	if value.IsNull() {
		return nil
	}

	valueType := value.Type()
	if valueType.IsObjectType() {
		// Create a block for objects
		block := body.AppendNewBlock(blockName, nil)
		blockBody := block.Body()

		// Add each field to the block
		for key, val := range value.AsValueMap() {
			if !val.IsNull() {
				err := hg.ctyValueToHCL(val, blockBody, key)
				if err != nil {
					return fmt.Errorf("failed to convert field %s: %v", key, err)
				}
			}
		}

	} else if valueType.IsListType() {
		// Handle lists by creating blocks for each element
		values := value.AsValueSlice()
		for _, val := range values {
			if !val.IsNull() {
				err := hg.ctyValueToHCL(val, body, blockName)
				if err != nil {
					return fmt.Errorf("failed to convert list element: %v", err)
				}
			}
		}

	} else if valueType.IsMapType() {
		// Handle maps by creating blocks for each key-value pair
		valueMap := value.AsValueMap()
		for key, val := range valueMap {
			if !val.IsNull() {
				// For maps, we create a block with the key as the label
				block := body.AppendNewBlock(blockName, []string{key})
				blockBody := block.Body()

				err := hg.ctyValueToHCL(val, blockBody, "")
				if err != nil {
					return fmt.Errorf("failed to convert map value for key %s: %v", key, err)
				}
			}
		}

	} else {
		// For primitive types, set as attribute
		body.SetAttributeValue(blockName, value)
	}

	return nil
}

// GenerateHCL is a convenience function for quick HCL generation
func GenerateHCL(config interface{}, blockName string) (string, error) {
	generator := NewHCLGenerator()
	return generator.ToHCL(config, blockName)
}

// GenerateHCLWithDefaults is a convenience function that includes default values
func GenerateHCLWithDefaults(config interface{}, blockName string) (string, error) {
	generator := NewHCLGenerator()
	generator.UseDefaults = true
	return generator.ToHCL(config, blockName)
}

// GenerateHCLWithoutDefaults is a convenience function that excludes default values
func GenerateHCLWithoutDefaults(config interface{}, blockName string) (string, error) {
	generator := NewHCLGenerator()
	generator.UseDefaults = false
	return generator.ToHCL(config, blockName)
}
