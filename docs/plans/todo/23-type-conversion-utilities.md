# Implementation Plan: Type Conversion Utilities

## Overview
Create centralized type conversion utilities to eliminate code duplication and ensure consistent type conversion patterns across the spooky codebase.

## Task Details
- **Task ID**: 8.4
- **Priority**: Medium
- **Files**: `internal/types/converter/`
- **Functions**: Type conversion utilities, interface-to-concrete conversion, reflection utilities

## Current State Analysis

### Existing Type Conversion Patterns
1. **Scattered Logic**: Type conversion logic is scattered across multiple packages
2. **Code Duplication**: Similar conversion logic repeated in multiple places
3. **Inconsistent Patterns**: Different packages use different conversion approaches
4. **Complex Conversions**: Some conversions are complex and error-prone

### Type Conversion Issues Found
1. **Interface-to-Concrete**: Complex interface-to-concrete type conversions
2. **HCL Value Conversion**: Repeated HCL value conversion logic
3. **Map Conversions**: Inconsistent map type conversions
4. **Reflection Usage**: Inconsistent use of reflection for conversions

## Implementation Requirements

### Type Conversion Compliance
The type conversion utilities must:
1. **Centralize conversion logic** in a dedicated package
2. **Provide consistent interfaces** for type conversions
3. **Handle common conversion patterns** efficiently
4. **Support reflection-based conversions** when needed
5. **Include comprehensive testing** for all conversions
6. **Document conversion patterns** and usage

### Required Dependencies
- All existing packages
- Reflection utilities
- Testing framework

## Detailed Implementation Plan

### Step 1: Core Conversion Package

#### 1.1 Base Converter Interface
```go
// internal/types/converter/converter.go
package converter

import (
    "fmt"
    "reflect"
    "spooky/internal/facts/types"
    "spooky/internal/interfaces"
    "spooky/internal/variables/types"
)

// Converter provides centralized type conversion utilities
type Converter struct{}

// NewConverter creates a new converter
func NewConverter() *Converter {
    return &Converter{}
}

// ConvertFactsToConcrete converts interface facts to concrete types
func (c *Converter) ConvertFactsToConcrete(factsContext interfaces.FactsContext) map[string]interface{} {
    if factsContext == nil {
        return make(map[string]interface{})
    }
    
    concreteFacts := make(map[string]interface{})
    
    // Handle different fact context types
    switch v := factsContext.(type) {
    case map[string]interface{}:
        concreteFacts = v
    case *types.FactCollection:
        if v != nil && v.Facts != nil {
            for key, fact := range v.Facts {
                if fact != nil {
                    concreteFacts[key] = fact.Value
                }
            }
        }
    case interfaces.FactsContext:
        // Handle custom facts context
        if machineFacts := v.GetMachineFacts(); machineFacts != nil {
            for machine, factCollection := range machineFacts {
                if factCollection != nil && factCollection.Facts != nil {
                    for key, fact := range factCollection.Facts {
                        if fact != nil {
                            concreteFacts[fmt.Sprintf("%s.%s", machine, key)] = fact.Value
                        }
                    }
                }
            }
        }
        
        // Handle global facts
        if globalFacts := v.GetGlobalFacts(); globalFacts != nil && globalFacts.Facts != nil {
            for key, fact := range globalFacts.Facts {
                if fact != nil {
                    concreteFacts[key] = fact.Value
                }
            }
        }
    default:
        // Try reflection-based conversion
        if converted := c.convertViaReflection(factsContext); converted != nil {
            concreteFacts = converted
        }
    }
    
    return concreteFacts
}

// ConvertVariablesToConcrete converts interface variables to concrete types
func (c *Converter) ConvertVariablesToConcrete(variablesContext interfaces.VariablesContext) map[string]interface{} {
    if variablesContext == nil {
        return make(map[string]interface{})
    }
    
    concreteVariables := make(map[string]interface{})
    
    // Handle different variable context types
    switch v := variablesContext.(type) {
    case map[string]interface{}:
        concreteVariables = v
    case *types.VariableCollection:
        if v != nil && v.Variables != nil {
            for _, variable := range v.Variables {
                if variable != nil {
                    concreteVariables[variable.Name] = variable.Value
                }
            }
        }
    case interfaces.VariablesContext:
        // Handle custom variables context
        if resolvedVars := v.GetResolvedVariables(); resolvedVars != nil {
            concreteVariables = resolvedVars
        }
    default:
        // Try reflection-based conversion
        if converted := c.convertViaReflection(variablesContext); converted != nil {
            concreteVariables = converted
        }
    }
    
    return concreteVariables
}

// convertViaReflection converts interface to map using reflection
func (c *Converter) convertViaReflection(value interface{}) map[string]interface{} {
    if value == nil {
        return nil
    }
    
    v := reflect.ValueOf(value)
    if v.Kind() == reflect.Ptr {
        v = v.Elem()
    }
    
    if v.Kind() == reflect.Map {
        result := make(map[string]interface{})
        for _, key := range v.MapKeys() {
            result[key.String()] = v.MapIndex(key).Interface()
        }
        return result
    }
    
    return nil
}
```

#### 1.2 HCL Value Converter
```go
// internal/types/converter/hcl.go
package converter

import (
    "fmt"
    "github.com/hashicorp/hcl/v2"
    "github.com/zclconf/go-cty/cty"
)

// HCLConverter provides HCL-specific conversion utilities
type HCLConverter struct{}

// NewHCLConverter creates a new HCL converter
func NewHCLConverter() *HCLConverter {
    return &HCLConverter{}
}

// ConvertHCLValue converts HCL value to Go interface{}
func (hc *HCLConverter) ConvertHCLValue(value cty.Value) (interface{}, error) {
    if value.IsNull() {
        return nil, nil
    }
    
    typ := value.Type()
    
    switch {
    case typ == cty.String:
        return value.AsString(), nil
    case typ == cty.Number:
        // Handle single number
        bf := value.AsBigFloat()
        if bf.IsInt() {
            i, _ := bf.Int64()
            return i, nil
        }
        f, _ := bf.Float64()
        return f, nil
    case typ == cty.Bool:
        return value.True(), nil
    case typ.IsListType():
        var result []interface{}
        it := value.ElementIterator()
        for it.Next() {
            _, val := it.Element()
            converted, err := hc.ConvertHCLValue(val)
            if err != nil {
                return nil, fmt.Errorf("failed to convert list element: %w", err)
            }
            result = append(result, converted)
        }
        return result, nil
    case typ.IsTupleType():
        var result []interface{}
        it := value.ElementIterator()
        for it.Next() {
            _, val := it.Element()
            converted, err := hc.ConvertHCLValue(val)
            if err != nil {
                return nil, fmt.Errorf("failed to convert tuple element: %w", err)
            }
            result = append(result, converted)
        }
        return result, nil
    case typ.IsMapType():
        result := make(map[string]interface{})
        it := value.ElementIterator()
        for it.Next() {
            key, val := it.Element()
            converted, err := hc.ConvertHCLValue(val)
            if err != nil {
                return nil, fmt.Errorf("failed to convert map value: %w", err)
            }
            result[key.AsString()] = converted
        }
        return result, nil
    case typ.IsObjectType():
        result := make(map[string]interface{})
        it := value.ElementIterator()
        for it.Next() {
            key, val := it.Element()
            converted, err := hc.ConvertHCLValue(val)
            if err != nil {
                return nil, fmt.Errorf("failed to convert object field: %w", err)
            }
            result[key.AsString()] = converted
        }
        return result, nil
    default:
        return value.GoString(), nil
    }
}

// ConvertHCLExpression converts HCL expression to Go interface{}
func (hc *HCLConverter) ConvertHCLExpression(expr hcl.Expression) (interface{}, error) {
    val, diags := expr.Value(nil)
    if diags.HasErrors() {
        return nil, fmt.Errorf("failed to evaluate expression: %v", diags)
    }
    
    return hc.ConvertHCLValue(val)
}
```

#### 1.3 Map Converter
```go
// internal/types/converter/map.go
package converter

import (
    "fmt"
    "reflect"
)

// MapConverter provides map-specific conversion utilities
type MapConverter struct{}

// NewMapConverter creates a new map converter
func NewMapConverter() *MapConverter {
    return &MapConverter{}
}

// ConvertMapToString converts map[string]interface{} to map[string]string
func (mc *MapConverter) ConvertMapToString(input map[string]interface{}) map[string]string {
    result := make(map[string]string)
    for key, value := range input {
        if str, ok := value.(string); ok {
            result[key] = str
        } else {
            result[key] = fmt.Sprintf("%v", value)
        }
    }
    return result
}

// ConvertMapToInterface converts map[string]string to map[string]interface{}
func (mc *MapConverter) ConvertMapToInterface(input map[string]string) map[string]interface{} {
    result := make(map[string]interface{})
    for key, value := range input {
        result[key] = value
    }
    return result
}

// ConvertSliceToString converts []interface{} to []string
func (mc *MapConverter) ConvertSliceToString(input []interface{}) []string {
    result := make([]string, len(input))
    for i, value := range input {
        if str, ok := value.(string); ok {
            result[i] = str
        } else {
            result[i] = fmt.Sprintf("%v", value)
        }
    }
    return result
}

// ConvertSliceToInterface converts []string to []interface{}
func (mc *MapConverter) ConvertSliceToInterface(input []string) []interface{} {
    result := make([]interface{}, len(input))
    for i, value := range input {
        result[i] = value
    }
    return result
}

// MergeMaps merges multiple maps into a single map
func (mc *MapConverter) MergeMaps(maps ...map[string]interface{}) map[string]interface{} {
    result := make(map[string]interface{})
    for _, m := range maps {
        for key, value := range m {
            result[key] = value
        }
    }
    return result
}
```

### Step 2: Advanced Conversion Utilities

#### 2.1 Reflection Utilities
```go
// internal/types/converter/reflection.go
package converter

import (
    "fmt"
    "reflect"
)

// ReflectionConverter provides reflection-based conversion utilities
type ReflectionConverter struct{}

// NewReflectionConverter creates a new reflection converter
func NewReflectionConverter() *ReflectionConverter {
    return &ReflectionConverter{}
}

// ConvertStructToMap converts a struct to map[string]interface{}
func (rc *ReflectionConverter) ConvertStructToMap(value interface{}) (map[string]interface{}, error) {
    if value == nil {
        return nil, fmt.Errorf("value cannot be nil")
    }
    
    v := reflect.ValueOf(value)
    if v.Kind() == reflect.Ptr {
        v = v.Elem()
    }
    
    if v.Kind() != reflect.Struct {
        return nil, fmt.Errorf("value must be a struct, got %s", v.Kind())
    }
    
    result := make(map[string]interface{})
    t := v.Type()
    
    for i := 0; i < v.NumField(); i++ {
        field := v.Field(i)
        fieldType := t.Field(i)
        
        // Skip unexported fields
        if !field.CanInterface() {
            continue
        }
        
        // Get field name (use json tag if available)
        fieldName := fieldType.Name
        if jsonTag := fieldType.Tag.Get("json"); jsonTag != "" && jsonTag != "-" {
            fieldName = jsonTag
        }
        
        result[fieldName] = field.Interface()
    }
    
    return result, nil
}

// ConvertMapToStruct converts map[string]interface{} to struct
func (rc *ReflectionConverter) ConvertMapToStruct(input map[string]interface{}, target interface{}) error {
    if target == nil {
        return fmt.Errorf("target cannot be nil")
    }
    
    v := reflect.ValueOf(target)
    if v.Kind() != reflect.Ptr {
        return fmt.Errorf("target must be a pointer to struct")
    }
    
    v = v.Elem()
    if v.Kind() != reflect.Struct {
        return fmt.Errorf("target must be a pointer to struct")
    }
    
    t := v.Type()
    
    for i := 0; i < v.NumField(); i++ {
        field := v.Field(i)
        fieldType := t.Field(i)
        
        // Skip unexported fields
        if !field.CanSet() {
            continue
        }
        
        // Get field name
        fieldName := fieldType.Name
        if jsonTag := fieldType.Tag.Get("json"); jsonTag != "" && jsonTag != "-" {
            fieldName = jsonTag
        }
        
        // Get value from map
        if value, exists := input[fieldName]; exists {
            if err := rc.setFieldValue(field, value); err != nil {
                return fmt.Errorf("failed to set field %s: %w", fieldName, err)
            }
        }
    }
    
    return nil
}

// setFieldValue sets a field value using reflection
func (rc *ReflectionConverter) setFieldValue(field reflect.Value, value interface{}) error {
    if !field.CanSet() {
        return fmt.Errorf("field cannot be set")
    }
    
    val := reflect.ValueOf(value)
    
    // Handle type conversion
    if val.Type().ConvertibleTo(field.Type()) {
        field.Set(val.Convert(field.Type()))
        return nil
    }
    
    // Handle special cases
    switch field.Kind() {
    case reflect.String:
        field.SetString(fmt.Sprintf("%v", value))
    case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
        if v, ok := value.(int); ok {
            field.SetInt(int64(v))
        } else {
            return fmt.Errorf("cannot convert %v to int", value)
        }
    case reflect.Float32, reflect.Float64:
        if v, ok := value.(float64); ok {
            field.SetFloat(v)
        } else {
            return fmt.Errorf("cannot convert %v to float", value)
        }
    case reflect.Bool:
        if v, ok := value.(bool); ok {
            field.SetBool(v)
        } else {
            return fmt.Errorf("cannot convert %v to bool", value)
        }
    default:
        return fmt.Errorf("unsupported field type: %s", field.Type())
    }
    
    return nil
}
```

### Step 3: Testing and Documentation

#### 3.1 Comprehensive Testing
```go
// internal/types/converter/converter_test.go
package converter

import (
    "testing"
    "spooky/internal/facts/types"
    "spooky/internal/interfaces"
)

func TestConverter_ConvertFactsToConcrete(t *testing.T) {
    converter := NewConverter()
    
    tests := []struct {
        name     string
        input    interfaces.FactsContext
        expected map[string]interface{}
    }{
        {
            name:     "nil input",
            input:    nil,
            expected: make(map[string]interface{}),
        },
        {
            name:     "map input",
            input:    map[string]interface{}{"key": "value"},
            expected: map[string]interface{}{"key": "value"},
        },
        {
            name: "fact collection input",
            input: &types.FactCollection{
                Facts: map[string]*types.Fact{
                    "key": {Key: "key", Value: "value"},
                },
            },
            expected: map[string]interface{}{"key": "value"},
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := converter.ConvertFactsToConcrete(tt.input)
            
            if len(result) != len(tt.expected) {
                t.Errorf("expected %d items, got %d", len(tt.expected), len(result))
            }
            
            for key, expectedValue := range tt.expected {
                if actualValue, exists := result[key]; !exists {
                    t.Errorf("expected key %s not found", key)
                } else if actualValue != expectedValue {
                    t.Errorf("expected value %v for key %s, got %v", expectedValue, key, actualValue)
                }
            }
        })
    }
}

func TestConverter_ConvertVariablesToConcrete(t *testing.T) {
    converter := NewConverter()
    
    tests := []struct {
        name     string
        input    interfaces.VariablesContext
        expected map[string]interface{}
    }{
        {
            name:     "nil input",
            input:    nil,
            expected: make(map[string]interface{}),
        },
        {
            name:     "map input",
            input:    map[string]interface{}{"key": "value"},
            expected: map[string]interface{}{"key": "value"},
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := converter.ConvertVariablesToConcrete(tt.input)
            
            if len(result) != len(tt.expected) {
                t.Errorf("expected %d items, got %d", len(tt.expected), len(result))
            }
            
            for key, expectedValue := range tt.expected {
                if actualValue, exists := result[key]; !exists {
                    t.Errorf("expected key %s not found", key)
                } else if actualValue != expectedValue {
                    t.Errorf("expected value %v for key %s, got %v", expectedValue, key, actualValue)
                }
            }
        })
    }
}
```

## Implementation Strategy

### Phase 1: Core Implementation (Week 1)
1. **Implement base converter** - Create core conversion utilities
2. **Add HCL converter** - Implement HCL-specific conversions
3. **Add map converter** - Implement map conversion utilities

### Phase 2: Advanced Features (Week 2)
1. **Implement reflection utilities** - Add reflection-based conversions
2. **Add comprehensive testing** - Test all conversion utilities
3. **Create documentation** - Document conversion patterns

### Phase 3: Integration (Week 3)
1. **Update existing packages** - Replace scattered conversion logic
2. **Test integration** - Test integration with existing code
3. **Performance optimization** - Optimize conversion performance

## Success Criteria

### Functional Requirements
- [ ] Centralized conversion utilities implemented
- [ ] All common conversion patterns supported
- [ ] Comprehensive testing completed
- [ ] Integration with existing code successful

### Quality Requirements
- [ ] Conversion utilities well-documented
- [ ] Performance optimized
- [ ] Error handling robust
- [ ] Code quality high

## Dependencies

### Required Dependencies
- All existing packages
- Reflection utilities
- Testing framework
- HCL library

### Optional Dependencies
- Performance testing tools
- Code coverage tools

## Risk Assessment

### High Risk
- **Breaking Changes**: Conversion changes may break existing code
- **Performance Impact**: Reflection-based conversions may have performance impact

### Medium Risk
- **Integration Complexity**: Integrating utilities across packages
- **Testing Complexity**: Testing all conversion scenarios

### Low Risk
- **Documentation**: Documentation updates are straightforward
- **Tool Integration**: Integration with existing tools

## Next Steps

1. **Start with core utilities** - Begin with basic conversion utilities
2. **Implement gradually** - Add utilities incrementally
3. **Test thoroughly** - Test all conversion scenarios
4. **Integrate carefully** - Integrate with existing code carefully
