# Function Improvement Plan: calculateCoverage

**Function:** `calculateCoverage`  
**File:** `internal/schemas/manager.go:XXX`  
**Current Complexity:** 10  
**Target Complexity:** < 8  
**Priority:** High

## Current Function Analysis

### Function Signature
```go
func (m *Manager) calculateCoverage(schema *spookytypesschemas.Schema, data interface{}) (*spookytypesschemas.CoverageResult, error)
```

### Current Issues
1. **Complex data traversal logic** - Multiple conditions for different data types and structures
2. **Mixed responsibilities** - Data traversal, field matching, calculation, and error handling
3. **Deep nesting** - Multiple levels of conditional logic for different data types
4. **Complex field matching** - Multiple conditions for field validation and coverage calculation
5. **Repetitive calculation patterns** - Similar logic repeated for different data structures

### Complexity Breakdown
- 10 cyclomatic complexity points from:
  - 1 main function entry
  - 2 data type validation checks
  - 3 field matching conditions
  - 2 error handling paths
  - 2 calculation conditions

## Refactoring Strategy

### Phase 1: Extract Data Validation (Immediate - 1 day)

#### Extract Data Validation
```go
func (m *Manager) validateDataForCoverage(data interface{}) error {
    if data == nil {
        return fmt.Errorf("data cannot be nil")
    }
    
    switch data.(type) {
    case map[string]interface{}, []interface{}, string, int, float64, bool:
        return nil
    default:
        return fmt.Errorf("unsupported data type: %T", data)
    }
}
```

#### Extract Schema Validation
```go
func (m *Manager) validateSchemaForCoverage(schema *spookytypesschemas.Schema) error {
    if schema == nil {
        return fmt.Errorf("schema cannot be nil")
    }
    
    if schema.Validation == nil || schema.Validation.Fields == nil {
        return fmt.Errorf("schema must have validation fields")
    }
    
    return nil
}
```

### Phase 2: Extract Coverage Calculation (Day 2)

#### Extract Coverage Calculation Logic
```go
func (m *Manager) calculateFieldCoverage(schema *spookytypesschemas.Schema, data interface{}) (*spookytypesschemas.CoverageResult, error) {
    result := &spookytypesschemas.CoverageResult{
        TotalFields:    len(schema.Validation.Fields),
        CoveredFields:  0,
        MissingFields:  []string{},
        ExtraFields:    []string{},
        Coverage:       0.0,
    }
    
    // Calculate coverage based on data type
    switch v := data.(type) {
    case map[string]interface{}:
        return m.calculateMapCoverage(schema, v, result)
    case []interface{}:
        return m.calculateArrayCoverage(schema, v, result)
    default:
        return m.calculatePrimitiveCoverage(schema, data, result)
    }
}

func (m *Manager) calculateMapCoverage(schema *spookytypesschemas.Schema, data map[string]interface{}, result *spookytypesschemas.CoverageResult) (*spookytypesschemas.CoverageResult, error) {
    // Check for covered fields
    for fieldName := range schema.Validation.Fields {
        if _, exists := data[fieldName]; exists {
            result.CoveredFields++
        } else {
            result.MissingFields = append(result.MissingFields, fieldName)
        }
    }
    
    // Check for extra fields
    for fieldName := range data {
        if _, exists := schema.Validation.Fields[fieldName]; !exists {
            result.ExtraFields = append(result.ExtraFields, fieldName)
        }
    }
    
    // Calculate coverage percentage
    result.Coverage = m.calculateCoveragePercentage(result.CoveredFields, result.TotalFields)
    
    return result, nil
}

func (m *Manager) calculateArrayCoverage(schema *spookytypesschemas.Schema, data []interface{}, result *spookytypesschemas.CoverageResult) (*spookytypesschemas.CoverageResult, error) {
    // For arrays, check if the schema supports array validation
    if len(data) > 0 {
        // Validate first element against schema
        firstElement := data[0]
        elementResult, err := m.calculateFieldCoverage(schema, firstElement)
        if err != nil {
            return nil, err
        }
        
        // Aggregate results for array
        result.CoveredFields = elementResult.CoveredFields
        result.MissingFields = elementResult.MissingFields
        result.ExtraFields = elementResult.ExtraFields
        result.Coverage = elementResult.Coverage
    }
    
    return result, nil
}

func (m *Manager) calculatePrimitiveCoverage(schema *spookytypesschemas.Schema, data interface{}, result *spookytypesschemas.CoverageResult) (*spookytypesschemas.CoverageResult, error) {
    // For primitive types, check if schema has a single field
    if len(schema.Validation.Fields) == 1 {
        // Assume the single field is covered
        result.CoveredFields = 1
        result.Coverage = 100.0
    } else {
        // Multiple fields but primitive data - mark all as missing
        for fieldName := range schema.Validation.Fields {
            result.MissingFields = append(result.MissingFields, fieldName)
        }
        result.Coverage = 0.0
    }
    
    return result, nil
}

func (m *Manager) calculateCoveragePercentage(covered, total int) float64 {
    if total == 0 {
        return 0.0
    }
    return float64(covered) / float64(total) * 100.0
}
```

#### Extract Field Validation
```go
func (m *Manager) validateFieldValue(field *spookytypesschemas.Field, value interface{}) error {
    if field.Required && value == nil {
        return fmt.Errorf("required field %s is missing", field.Name)
    }
    
    if value == nil {
        return nil // Optional field can be nil
    }
    
    // Type validation
    if err := m.validateFieldType(field, value); err != nil {
        return fmt.Errorf("field %s: %w", field.Name, err)
    }
    
    return nil
}

func (m *Manager) validateFieldType(field *spookytypesschemas.Field, value interface{}) error {
    switch field.Type {
    case "string":
        return m.validateStringField(value)
    case "number":
        return m.validateNumberField(value)
    case "boolean":
        return m.validateBooleanField(value)
    case "array":
        return m.validateArrayField(value)
    case "object":
        return m.validateObjectField(value)
    default:
        return fmt.Errorf("unsupported field type: %s", field.Type)
    }
}

func (m *Manager) validateStringField(value interface{}) error {
    if _, ok := value.(string); !ok {
        return fmt.Errorf("expected string, got %T", value)
    }
    return nil
}

func (m *Manager) validateNumberField(value interface{}) error {
    switch value.(type) {
    case int, int32, int64, float32, float64:
        return nil
    default:
        return fmt.Errorf("expected number, got %T", value)
    }
}

func (m *Manager) validateBooleanField(value interface{}) error {
    if _, ok := value.(bool); !ok {
        return fmt.Errorf("expected boolean, got %T", value)
    }
    return nil
}

func (m *Manager) validateArrayField(value interface{}) error {
    if _, ok := value.([]interface{}); !ok {
        return fmt.Errorf("expected array, got %T", value)
    }
    return nil
}

func (m *Manager) validateObjectField(value interface{}) error {
    if _, ok := value.(map[string]interface{}); !ok {
        return fmt.Errorf("expected object, got %T", value)
    }
    return nil
}
```

### Phase 3: Refactored Main Function (Day 3)

#### Final Refactored Function
```go
func (m *Manager) calculateCoverage(schema *spookytypesschemas.Schema, data interface{}) (*spookytypesschemas.CoverageResult, error) {
    // Validate inputs
    if err := m.validateSchemaForCoverage(schema); err != nil {
        return nil, err
    }
    
    if err := m.validateDataForCoverage(data); err != nil {
        return nil, err
    }
    
    // Calculate coverage
    return m.calculateFieldCoverage(schema, data)
}
```

## Complexity Reduction Metrics

### Before Refactoring
- **Cyclomatic Complexity:** 10
- **Lines of Code:** ~80
- **Nesting Levels:** 4-5
- **Responsibilities:** 4 (validation, calculation, field matching, error handling)

### After Refactoring
- **Cyclomatic Complexity:** 3 (main function) + 1-3 per helper function
- **Lines of Code:** ~150 (distributed across multiple focused functions)
- **Nesting Levels:** 1-2
- **Responsibilities:** 1 per function (single responsibility principle)

## Testing Strategy

### Unit Tests for Extracted Functions
```go
func TestValidateDataForCoverage(t *testing.T) {
    tests := []struct {
        name    string
        data    interface{}
        wantErr bool
    }{
        {
            name:    "nil data",
            data:    nil,
            wantErr: true,
        },
        {
            name:    "map data",
            data:    map[string]interface{}{"key": "value"},
            wantErr: false,
        },
        {
            name:    "array data",
            data:    []interface{}{"item1", "item2"},
            wantErr: false,
        },
        {
            name:    "string data",
            data:    "test",
            wantErr: false,
        },
        {
            name:    "number data",
            data:    42,
            wantErr: false,
        },
        {
            name:    "boolean data",
            data:    true,
            wantErr: false,
        },
        {
            name:    "unsupported type",
            data:    make(chan int),
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := validateDataForCoverage(tt.data)
            if (err != nil) != tt.wantErr {
                t.Errorf("validateDataForCoverage() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}

func TestCalculateMapCoverage(t *testing.T) {
    schema := &spookytypesschemas.Schema{
        Validation: &spookytypesschemas.Validation{
            Fields: map[string]*spookytypesschemas.Field{
                "name": {Type: "string", Required: true},
                "port": {Type: "number", Required: false},
                "enabled": {Type: "boolean", Required: true},
            },
        },
    }
    
    data := map[string]interface{}{
        "name":    "test",
        "enabled": true,
        "extra":   "field",
    }
    
    result := &spookytypesschemas.CoverageResult{
        TotalFields: 3,
    }
    
    coverage, err := calculateMapCoverage(schema, data, result)
    if err != nil {
        t.Errorf("calculateMapCoverage() error = %v", err)
        return
    }
    
    if coverage.CoveredFields != 2 {
        t.Errorf("calculateMapCoverage() covered fields = %d, want %d", coverage.CoveredFields, 2)
    }
    if len(coverage.MissingFields) != 1 {
        t.Errorf("calculateMapCoverage() missing fields = %d, want %d", len(coverage.MissingFields), 1)
    }
    if len(coverage.ExtraFields) != 1 {
        t.Errorf("calculateMapCoverage() extra fields = %d, want %d", len(coverage.ExtraFields), 1)
    }
    if coverage.Coverage != 66.66666666666666 {
        t.Errorf("calculateMapCoverage() coverage = %f, want %f", coverage.Coverage, 66.66666666666666)
    }
}

func TestValidateStringField(t *testing.T) {
    tests := []struct {
        name    string
        value   interface{}
        wantErr bool
    }{
        {
            name:    "valid string",
            value:   "test",
            wantErr: false,
        },
        {
            name:    "invalid type",
            value:   42,
            wantErr: true,
        },
        {
            name:    "nil value",
            value:   nil,
            wantErr: false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := validateStringField(tt.value)
            if (err != nil) != tt.wantErr {
                t.Errorf("validateStringField() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

### Integration Tests
```go
func TestCalculateCoverageIntegration(t *testing.T) {
    schema := &spookytypesschemas.Schema{
        Validation: &spookytypesschemas.Validation{
            Fields: map[string]*spookytypesschemas.Field{
                "name": {Type: "string", Required: true},
                "port": {Type: "number", Required: false},
                "enabled": {Type: "boolean", Required: true},
            },
        },
    }
    
    tests := []struct {
        name    string
        data    interface{}
        wantCoverage float64
        wantErr      bool
    }{
        {
            name:    "complete data",
            data:    map[string]interface{}{"name": "test", "port": 8080, "enabled": true},
            wantCoverage: 100.0,
            wantErr:      false,
        },
        {
            name:    "partial data",
            data:    map[string]interface{}{"name": "test", "enabled": true},
            wantCoverage: 66.66666666666666,
            wantErr:      false,
        },
        {
            name:    "empty data",
            data:    map[string]interface{}{},
            wantCoverage: 0.0,
            wantErr:      false,
        },
        {
            name:    "array data",
            data:    []interface{}{map[string]interface{}{"name": "test", "enabled": true}},
            wantCoverage: 66.66666666666666,
            wantErr:      false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            coverage, err := calculateCoverage(schema, tt.data)
            if (err != nil) != tt.wantErr {
                t.Errorf("calculateCoverage() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if coverage.Coverage != tt.wantCoverage {
                t.Errorf("calculateCoverage() coverage = %f, want %f", coverage.Coverage, tt.wantCoverage)
            }
        })
    }
}
```

## Implementation Timeline

### Day 1: Extract Data Validation
- [ ] Extract `validateDataForCoverage`
- [ ] Extract `validateSchemaForCoverage`
- [ ] Add unit tests for validation functions
- [ ] Verify data validation works correctly

### Day 2: Extract Coverage Calculation
- [ ] Extract `calculateFieldCoverage`
- [ ] Extract `calculateMapCoverage`
- [ ] Extract `calculateArrayCoverage`
- [ ] Extract `calculatePrimitiveCoverage`
- [ ] Extract `calculateCoveragePercentage`
- [ ] Extract field validation functions
- [ ] Add unit tests for coverage calculation functions
- [ ] Verify coverage calculation works correctly

### Day 3: Complete Refactoring
- [ ] Refactor main `calculateCoverage` function
- [ ] Add integration tests
- [ ] Verify complexity reduction with gocyclo
- [ ] Code review and documentation
- [ ] Performance testing

## Success Criteria

### Complexity Reduction
- [ ] Main function complexity < 8
- [ ] All extracted functions complexity < 5
- [ ] No function exceeds complexity threshold

### Code Quality
- [ ] Single responsibility principle maintained
- [ ] Clear separation of concerns
- [ ] Comprehensive test coverage (>90%)
- [ ] No regression in functionality

### Maintainability
- [ ] Easy to modify individual coverage calculation components
- [ ] Clear function names and purposes
- [ ] Well-documented functions
- [ ] Consistent coverage calculation patterns

## Risk Mitigation

### Potential Risks
1. **Coverage calculation changes** - May affect coverage accuracy
2. **Data type handling changes** - May affect validation behavior
3. **Field matching changes** - May affect coverage results

### Mitigation Strategies
1. **Comprehensive testing** - Test all data types and scenarios
2. **Coverage validation** - Ensure coverage calculation remains accurate
3. **Field validation** - Ensure field validation remains consistent
4. **Gradual migration** - Implement changes incrementally with validation

## Verification

### Complexity Verification
```bash
# Verify complexity reduction
gocyclo internal/schemas/manager.go | grep calculateCoverage
```

### Functionality Verification
```bash
# Test coverage calculation
go test ./internal/schemas -run TestCalculateCoverage
```

### Performance Verification
- [ ] No performance regression
- [ ] Memory usage remains stable
- [ ] Coverage calculation performance remains acceptable

## Conclusion

This refactoring will reduce the cyclomatic complexity of `calculateCoverage` from 10 to approximately 3, while improving code maintainability, testability, and readability. The extracted functions will follow single responsibility principles and be easily testable in isolation.

**Expected Outcome:** A clean, maintainable function with clear separation of concerns and comprehensive test coverage for schema coverage calculation operations.
