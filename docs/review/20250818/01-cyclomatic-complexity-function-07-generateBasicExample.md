# Function Improvement Plan: generateBasicExample

**Function:** `generateBasicExample`  
**File:** `internal/schemas/manager.go:XXX`  
**Current Complexity:** 10  
**Target Complexity:** < 8  
**Priority:** High

## Current Function Analysis

### Function Signature
```go
func (m *Manager) generateBasicExample(schema *spookytypesschemas.Schema) (string, error)
```

### Current Issues
1. **Complex example generation logic** - Multiple conditions for different schema types and fields
2. **Mixed responsibilities** - Example generation, formatting, validation, and error handling
3. **Deep nesting** - Multiple levels of conditional logic for different schema components
4. **Complex string building** - Multiple string concatenation operations with conditional logic
5. **Repetitive field processing patterns** - Similar logic repeated for different field types

### Complexity Breakdown
- 10 cyclomatic complexity points from:
  - 1 main function entry
  - 2 schema validation checks
  - 3 field type processing conditions
  - 2 error handling paths
  - 2 string building conditions

## Refactoring Strategy

### Phase 1: Extract Schema Validation (Immediate - 1 day)

#### Extract Schema Validation
```go
func (m *Manager) validateSchemaForExample(schema *spookytypesschemas.Schema) error {
    if schema == nil {
        return fmt.Errorf("schema cannot be nil")
    }
    
    if schema.Name == "" {
        return fmt.Errorf("schema name is required")
    }
    
    if schema.Validation == nil || schema.Validation.Fields == nil {
        return fmt.Errorf("schema must have validation fields")
    }
    
    return nil
}
```

#### Extract Example Header Generation
```go
func (m *Manager) generateExampleHeader(schema *spookytypesschemas.Schema) string {
    var header strings.Builder
    
    header.WriteString(fmt.Sprintf("# Example: %s\n\n", schema.Name))
    header.WriteString(fmt.Sprintf("This is a basic example for the %s schema.\n", schema.Name))
    header.WriteString("Modify the values according to your needs.\n\n")
    
    return header.String()
}
```

### Phase 2: Extract Field Processing (Day 2)

#### Extract Field Example Generation
```go
func (m *Manager) generateFieldExamples(fields map[string]*spookytypesschemas.Field) string {
    var examples strings.Builder
    
    examples.WriteString("## Fields\n\n")
    
    for fieldName, field := range fields {
        example := m.generateFieldExample(fieldName, field)
        examples.WriteString(example)
    }
    
    return examples.String()
}

func (m *Manager) generateFieldExample(fieldName string, field *spookytypesschemas.Field) string {
    var example strings.Builder
    
    example.WriteString(fmt.Sprintf("### %s\n", fieldName))
    example.WriteString(fmt.Sprintf("**Type:** %s\n", field.Type))
    example.WriteString(fmt.Sprintf("**Required:** %t\n", field.Required))
    
    if field.Description != "" {
        example.WriteString(fmt.Sprintf("**Description:** %s\n", field.Description))
    }
    
    example.WriteString(fmt.Sprintf("**Example Value:** %s\n\n", m.generateExampleValue(field)))
    
    return example.String()
}

func (m *Manager) generateExampleValue(field *spookytypesschemas.Field) string {
    switch field.Type {
    case "string":
        return m.generateStringExample(field)
    case "number":
        return m.generateNumberExample(field)
    case "boolean":
        return m.generateBooleanExample(field)
    case "array":
        return m.generateArrayExample(field)
    case "object":
        return m.generateObjectExample(field)
    default:
        return fmt.Sprintf("/* %s value */", field.Type)
    }
}

func (m *Manager) generateStringExample(field *spookytypesschemas.Field) string {
    if field.Example != "" {
        return fmt.Sprintf(`"%s"`, field.Example)
    }
    
    // Generate contextual examples based on field name
    switch {
    case strings.Contains(strings.ToLower(field.Name), "name"):
        return `"example-name"`
    case strings.Contains(strings.ToLower(field.Name), "description"):
        return `"Example description"`
    case strings.Contains(strings.ToLower(field.Name), "url"):
        return `"https://example.com"`
    case strings.Contains(strings.ToLower(field.Name), "email"):
        return `"user@example.com"`
    case strings.Contains(strings.ToLower(field.Name), "path"):
        return `"/path/to/resource"`
    default:
        return `"example-string"`
    }
}

func (m *Manager) generateNumberExample(field *spookytypesschemas.Field) string {
    if field.Example != "" {
        return field.Example
    }
    
    // Generate contextual examples based on field name
    switch {
    case strings.Contains(strings.ToLower(field.Name), "port"):
        return "8080"
    case strings.Contains(strings.ToLower(field.Name), "timeout"):
        return "30"
    case strings.Contains(strings.ToLower(field.Name), "count"):
        return "5"
    case strings.Contains(strings.ToLower(field.Name), "size"):
        return "1024"
    default:
        return "42"
    }
}

func (m *Manager) generateBooleanExample(field *spookytypesschemas.Field) string {
    if field.Example != "" {
        return field.Example
    }
    
    // Generate contextual examples based on field name
    switch {
    case strings.Contains(strings.ToLower(field.Name), "enabled"):
        return "true"
    case strings.Contains(strings.ToLower(field.Name), "disabled"):
        return "false"
    case strings.Contains(strings.ToLower(field.Name), "required"):
        return "true"
    case strings.Contains(strings.ToLower(field.Name), "optional"):
        return "false"
    default:
        return "true"
    }
}

func (m *Manager) generateArrayExample(field *spookytypesschemas.Field) string {
    if field.Example != "" {
        return field.Example
    }
    
    return `["item1", "item2", "item3"]`
}

func (m *Manager) generateObjectExample(field *spookytypesschemas.Field) string {
    if field.Example != "" {
        return field.Example
    }
    
    return `{
  "key1": "value1",
  "key2": "value2"
}`
}
```

### Phase 3: Refactored Main Function (Day 3)

#### Final Refactored Function
```go
func (m *Manager) generateBasicExample(schema *spookytypesschemas.Schema) (string, error) {
    // Validate schema
    if err := m.validateSchemaForExample(schema); err != nil {
        return "", err
    }
    
    var example strings.Builder
    
    // Generate header
    example.WriteString(m.generateExampleHeader(schema))
    
    // Generate field examples
    example.WriteString(m.generateFieldExamples(schema.Validation.Fields))
    
    // Generate usage section
    example.WriteString(m.generateUsageSection(schema))
    
    return example.String(), nil
}

func (m *Manager) generateUsageSection(schema *spookytypesschemas.Schema) string {
    var usage strings.Builder
    
    usage.WriteString("## Usage\n\n")
    usage.WriteString("To use this schema, create a configuration file with the following structure:\n\n")
    usage.WriteString("```hcl\n")
    usage.WriteString(fmt.Sprintf("%s {\n", strings.ToLower(schema.Name)))
    
    for fieldName, field := range schema.Validation.Fields {
        exampleValue := m.generateExampleValue(field)
        usage.WriteString(fmt.Sprintf("  %s = %s\n", fieldName, exampleValue))
    }
    
    usage.WriteString("}\n")
    usage.WriteString("```\n")
    
    return usage.String()
}
```

## Complexity Reduction Metrics

### Before Refactoring
- **Cyclomatic Complexity:** 10
- **Lines of Code:** ~70
- **Nesting Levels:** 4-5
- **Responsibilities:** 4 (validation, header generation, field processing, formatting)

### After Refactoring
- **Cyclomatic Complexity:** 3 (main function) + 1-3 per helper function
- **Lines of Code:** ~120 (distributed across multiple focused functions)
- **Nesting Levels:** 1-2
- **Responsibilities:** 1 per function (single responsibility principle)

## Testing Strategy

### Unit Tests for Extracted Functions
```go
func TestValidateSchemaForExample(t *testing.T) {
    tests := []struct {
        name    string
        schema  *spookytypesschemas.Schema
        wantErr bool
    }{
        {
            name:    "nil schema",
            schema:  nil,
            wantErr: true,
        },
        {
            name: "empty name",
            schema: &spookytypesschemas.Schema{
                Name: "",
                Validation: &spookytypesschemas.Validation{
                    Fields: map[string]*spookytypesschemas.Field{},
                },
            },
            wantErr: true,
        },
        {
            name: "nil validation",
            schema: &spookytypesschemas.Schema{
                Name:       "test",
                Validation: nil,
            },
            wantErr: true,
        },
        {
            name: "valid schema",
            schema: &spookytypesschemas.Schema{
                Name: "test",
                Validation: &spookytypesschemas.Validation{
                    Fields: map[string]*spookytypesschemas.Field{
                        "testField": {
                            Type:     "string",
                            Required: true,
                        },
                    },
                },
            },
            wantErr: false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := validateSchemaForExample(tt.schema)
            if (err != nil) != tt.wantErr {
                t.Errorf("validateSchemaForExample() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}

func TestGenerateStringExample(t *testing.T) {
    tests := []struct {
        name  string
        field *spookytypesschemas.Field
        want  string
    }{
        {
            name: "with example",
            field: &spookytypesschemas.Field{
                Name:    "test",
                Type:    "string",
                Example: "custom-example",
            },
            want: `"custom-example"`,
        },
        {
            name: "name field",
            field: &spookytypesschemas.Field{
                Name: "server_name",
                Type: "string",
            },
            want: `"example-name"`,
        },
        {
            name: "description field",
            field: &spookytypesschemas.Field{
                Name: "description",
                Type: "string",
            },
            want: `"Example description"`,
        },
        {
            name: "default string",
            field: &spookytypesschemas.Field{
                Name: "test",
                Type: "string",
            },
            want: `"example-string"`,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := generateStringExample(tt.field)
            if got != tt.want {
                t.Errorf("generateStringExample() = %v, want %v", got, tt.want)
            }
        })
    }
}

func TestGenerateNumberExample(t *testing.T) {
    tests := []struct {
        name  string
        field *spookytypesschemas.Field
        want  string
    }{
        {
            name: "with example",
            field: &spookytypesschemas.Field{
                Name:    "test",
                Type:    "number",
                Example: "123",
            },
            want: "123",
        },
        {
            name: "port field",
            field: &spookytypesschemas.Field{
                Name: "server_port",
                Type: "number",
            },
            want: "8080",
        },
        {
            name: "default number",
            field: &spookytypesschemas.Field{
                Name: "test",
                Type: "number",
            },
            want: "42",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := generateNumberExample(tt.field)
            if got != tt.want {
                t.Errorf("generateNumberExample() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

### Integration Tests
```go
func TestGenerateBasicExampleIntegration(t *testing.T) {
    schema := &spookytypesschemas.Schema{
        Name: "TestSchema",
        Validation: &spookytypesschemas.Validation{
            Fields: map[string]*spookytypesschemas.Field{
                "name": {
                    Type:        "string",
                    Required:    true,
                    Description: "The name of the resource",
                },
                "port": {
                    Type:     "number",
                    Required: false,
                },
                "enabled": {
                    Type:     "boolean",
                    Required: true,
                },
            },
        },
    }
    
    example, err := generateBasicExample(schema)
    if err != nil {
        t.Errorf("generateBasicExample() error = %v", err)
        return
    }
    
    // Verify example contains expected sections
    if !strings.Contains(example, "# Example: TestSchema") {
        t.Error("Example missing header")
    }
    if !strings.Contains(example, "### name") {
        t.Error("Example missing name field")
    }
    if !strings.Contains(example, "### port") {
        t.Error("Example missing port field")
    }
    if !strings.Contains(example, "### enabled") {
        t.Error("Example missing enabled field")
    }
    if !strings.Contains(example, "## Usage") {
        t.Error("Example missing usage section")
    }
}
```

## Implementation Timeline

### Day 1: Extract Schema Validation
- [ ] Extract `validateSchemaForExample`
- [ ] Extract `generateExampleHeader`
- [ ] Add unit tests for validation functions
- [ ] Verify schema validation works correctly

### Day 2: Extract Field Processing
- [ ] Extract `generateFieldExamples`
- [ ] Extract `generateFieldExample`
- [ ] Extract `generateExampleValue`
- [ ] Extract type-specific example generators
- [ ] Add unit tests for field processing functions
- [ ] Verify field processing works correctly

### Day 3: Complete Refactoring
- [ ] Refactor main `generateBasicExample` function
- [ ] Add `generateUsageSection`
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
- [ ] Easy to modify individual example generation components
- [ ] Clear function names and purposes
- [ ] Well-documented functions
- [ ] Consistent example generation patterns

## Risk Mitigation

### Potential Risks
1. **Example format changes** - May affect user expectations
2. **Field type handling changes** - May affect example generation
3. **String building changes** - May affect output formatting

### Mitigation Strategies
1. **Preserve output format** - Ensure identical example output
2. **Comprehensive testing** - Test all field types and scenarios
3. **Format validation** - Ensure example formatting remains consistent
4. **Gradual migration** - Implement changes incrementally with validation

## Verification

### Complexity Verification
```bash
# Verify complexity reduction
gocyclo internal/schemas/manager.go | grep generateBasicExample
```

### Functionality Verification
```bash
# Test example generation
go test ./internal/schemas -run TestGenerateBasicExample
```

### Performance Verification
- [ ] No performance regression
- [ ] Memory usage remains stable
- [ ] String building performance remains acceptable

## Conclusion

This refactoring will reduce the cyclomatic complexity of `generateBasicExample` from 10 to approximately 3, while improving code maintainability, testability, and readability. The extracted functions will follow single responsibility principles and be easily testable in isolation.

**Expected Outcome:** A clean, maintainable function with clear separation of concerns and comprehensive test coverage for schema example generation operations.
