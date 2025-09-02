package hcl

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

// ConfigGenerator extracts default values from struct tags and creates default instances
type ConfigGenerator struct {
	sharedUtils *SharedHCLUtils
}

// NewConfigGenerator creates a new config generator
func NewConfigGenerator() *ConfigGenerator {
	return &ConfigGenerator{
		sharedUtils: NewSharedHCLUtils(),
	}
}

// ToHCL converts the default config to HCL string
func (cg *ConfigGenerator) ToHCL(config interface{}) (string, error) {
	return cg.ToHCLWithBlockName(config, "")
}

// ToHCLWithBlockName converts the default config to HCL string with a custom block name
func (cg *ConfigGenerator) ToHCLWithBlockName(config interface{}, blockName string) (string, error) {
	// Convert the config to HCL using the HCL library
	hclFile := hclwrite.NewEmptyFile()

	// Convert the config to cty.Value
	ctyValue, err := cg.structToCty(config)
	if err != nil {
		return "", fmt.Errorf("failed to convert config to cty.Value: %v", err)
	}

	// Create the root body
	rootBody := hclFile.Body()

	// Determine the block name
	if blockName == "" {
		// Get the config type name for the root block
		configType := reflect.TypeOf(config)
		if configType.Kind() == reflect.Ptr {
			configType = configType.Elem()
		}
		blockName = configType.Name()
	}

	// Convert to HCL block
	err = cg.ctyValueToHCL(ctyValue, rootBody, blockName)
	if err != nil {
		return "", fmt.Errorf("failed to convert cty.Value to HCL: %v", err)
	}

	return string(hclFile.Bytes()), nil
}

// structToCty converts a Go struct to cty.Value, respecting schema tags
func (cg *ConfigGenerator) structToCty(config interface{}) (cty.Value, error) {
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
		ctyVal, err := cg.fieldToCty(field)
		if err != nil {
			return cty.NilVal, fmt.Errorf("failed to convert field %s: %v", jsonTag, err)
		}

		if !ctyVal.IsNull() {
			valueMap[jsonTag] = ctyVal
		}
	}

	return cty.ObjectVal(valueMap), nil
}

// fieldToCty converts a struct field to cty.Value
func (cg *ConfigGenerator) fieldToCty(field reflect.Value) (cty.Value, error) {
	switch field.Kind() {
	case reflect.Ptr:
		// Handle pointers
		if field.IsNil() {
			return cty.NilVal, nil
		}
		return cg.fieldToCty(field.Elem())
	case reflect.String:
		return cty.StringVal(field.String()), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return cty.NumberIntVal(field.Int()), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return cty.NumberUIntVal(field.Uint()), nil
	case reflect.Float32, reflect.Float64:
		return cty.NumberFloatVal(field.Float()), nil
	case reflect.Bool:
		return cty.BoolVal(field.Bool()), nil
	case reflect.Slice:
		// Handle slices
		if field.Len() == 0 {
			return cty.ListValEmpty(cty.String), nil
		}
		elements := make([]cty.Value, field.Len())
		for i := 0; i < field.Len(); i++ {
			element, err := cg.fieldToCty(field.Index(i))
			if err != nil {
				return cty.NilVal, fmt.Errorf("failed to convert slice element %d: %v", i, err)
			}
			elements[i] = element
		}
		return cty.ListVal(elements), nil
	case reflect.Struct:
		// Handle nested structs
		if field.Type() == reflect.TypeOf(time.Time{}) {
			// Special handling for time.Time
			return cty.StringVal(field.Interface().(time.Time).Format(time.RFC3339)), nil
		}
		return cg.structToCty(field.Interface())
	case reflect.Map:
		// Handle maps
		if field.Len() == 0 {
			return cty.MapValEmpty(cty.String), nil
		}
		keys := field.MapKeys()
		valueMap := make(map[string]cty.Value)
		for _, key := range keys {
			keyStr := fmt.Sprintf("%v", key.Interface())
			val, err := cg.fieldToCty(field.MapIndex(key))
			if err != nil {
				return cty.NilVal, fmt.Errorf("failed to convert map value for key %s: %v", keyStr, err)
			}
			valueMap[keyStr] = val
		}
		return cty.MapVal(valueMap), nil
	default:
		return cty.NilVal, fmt.Errorf("unsupported field type: %s", field.Kind())
	}
}

// ctyValueToHCL converts a cty.Value to HCL and writes it to the body
func (cg *ConfigGenerator) ctyValueToHCL(value cty.Value, body *hclwrite.Body, blockName string) error {
	return cg.sharedUtils.CtyValueToHCL(value, body, blockName)
}
