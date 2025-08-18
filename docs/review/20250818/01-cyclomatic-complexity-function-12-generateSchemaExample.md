# Function Improvement Plan: generateSchemaExample

**Function:** `generateSchemaExample`  
**File:** `internal/schemas/manager.go:XXX`  
**Current Complexity:** 9  
**Target Complexity:** < 8  
**Priority:** Medium

## Current Function Analysis

### Function Signature
```go
func (m *Manager) generateSchemaExample(schema *spookytypesschemas.Schema) (string, error)
```

### Current Issues
1. **Complex example generation logic** - Multiple conditions for different schema types and field types
2. **Mixed responsibilities** - Schema validation, example generation, formatting, and error handling
3. **Deep nesting** - Multiple levels of conditional logic for different schema components
4. **Complex field processing** - Multiple conditions for different field types and validation rules
5. **Repetitive generation patterns** - Similar logic repeated for different field types

### Complexity Breakdown
- 9 cyclomatic complexity points from:
  - 1 main function entry
  - 2 schema validation checks
  - 3 field type conditions
  - 2 error handling paths
  - 1 example generation condition

## Refactoring Strategy

### Phase 1: Extract Schema Validation (Immediate - 1 day)

#### Extract Schema Validation
```go
func (m *Manager) validateSchemaForExample(schema *spookytypesschemas.Schema) error {
    if schema == nil {
        return fmt.Errorf("schema cannot be nil")
    }
    
    if schema.Validation == nil {
        return fmt.Errorf("schema must have validation configuration")
    }
    
    if len(schema.Validation.Fields) == 0 {
        return fmt.Errorf("schema must have at least one field")
    }
    
    return nil
}
```

#### Extract Example Initialization
```go
func (m *Manager) initializeExampleGeneration() *ExampleGenerator {
    return &ExampleGenerator{
        Buffer:    &strings.Builder{},
        Indent:    0,
        MaxDepth:  5,
    }
}

type ExampleGenerator struct {
    Buffer   *strings.Builder
    Indent   int
    MaxDepth int
}

func (g *ExampleGenerator) writeLine(format string, args ...interface{}) {
    indent := strings.Repeat("  ", g.Indent)
    g.Buffer.WriteString(indent)
    g.Buffer.WriteString(fmt.Sprintf(format, args...))
    g.Buffer.WriteString("\n")
}
```

### Phase 2: Extract Field Processing (Day 2)

#### Extract Field Type Processing
```go
func (m *Manager) generateFieldExample(field *spookytypesschemas.Field, generator *ExampleGenerator) error {
    switch field.Type {
    case "string":
        return m.generateStringExample(field, generator)
    case "number":
        return m.generateNumberExample(field, generator)
    case "boolean":
        return m.generateBooleanExample(field, generator)
    case "array":
        return m.generateArrayExample(field, generator)
    case "object":
        return m.generateObjectExample(field, generator)
    default:
        return fmt.Errorf("unsupported field type: %s", field.Type)
    }
}

func (m *Manager) generateStringExample(field *spookytypesschemas.Field, generator *ExampleGenerator) error {
    example := m.generateStringValue(field)
    generator.writeLine("%s = %q", field.Name, example)
    return nil
}

func (m *Manager) generateNumberExample(field *spookytypesschemas.Field, generator *ExampleGenerator) error {
    example := m.generateNumberValue(field)
    generator.writeLine("%s = %v", field.Name, example)
    return nil
}

func (m *Manager) generateBooleanExample(field *spookytypesschemas.Field, generator *ExampleGenerator) error {
    example := m.generateBooleanValue(field)
    generator.writeLine("%s = %v", field.Name, example)
    return nil
}

func (m *Manager) generateArrayExample(field *spookytypesschemas.Field, generator *ExampleGenerator) error {
    generator.writeLine("%s = [", field.Name)
    generator.Indent++
    
    // Generate array elements
    if err := m.generateArrayElements(field, generator); err != nil {
        return err
    }
    
    generator.Indent--
    generator.writeLine("]")
    return nil
}

func (m *Manager) generateObjectExample(field *spookytypesschemas.Field, generator *ExampleGenerator) error {
    generator.writeLine("%s = {", field.Name)
    generator.Indent++
    
    // Generate object properties
    if err := m.generateObjectProperties(field, generator); err != nil {
        return err
    }
    
    generator.Indent--
    generator.writeLine("}")
    return nil
}
```

#### Extract Value Generation
```go
func (m *Manager) generateStringValue(field *spookytypesschemas.Field) string {
    // Check for predefined examples
    if field.Example != "" {
        return field.Example
    }
    
    // Check for validation constraints
    if field.Validation != nil && field.Validation.Pattern != "" {
        return m.generatePatternExample(field.Validation.Pattern)
    }
    
    // Generate based on field name
    return m.generateStringFromName(field.Name)
}

func (m *Manager) generateNumberValue(field *spookytypesschemas.Field) interface{} {
    // Check for predefined examples
    if field.Example != "" {
        if num, err := strconv.ParseFloat(field.Example, 64); err == nil {
            return num
        }
    }
    
    // Check for validation constraints
    if field.Validation != nil {
        if field.Validation.Min != nil {
            return *field.Validation.Min
        }
        if field.Validation.Max != nil {
            return *field.Validation.Max
        }
    }
    
    // Generate based on field name
    return m.generateNumberFromName(field.Name)
}

func (m *Manager) generateBooleanValue(field *spookytypesschemas.Field) bool {
    // Check for predefined examples
    if field.Example != "" {
        if field.Example == "true" {
            return true
        }
        if field.Example == "false" {
            return false
        }
    }
    
    // Generate based on field name
    return m.generateBooleanFromName(field.Name)
}

func (m *Manager) generateArrayElements(field *spookytypesschemas.Field, generator *ExampleGenerator) error {
    // Generate 1-3 example elements
    count := m.determineArrayElementCount(field)
    
    for i := 0; i < count; i++ {
        if err := m.generateArrayElement(field, generator, i); err != nil {
            return err
        }
    }
    
    return nil
}

func (m *Manager) generateObjectProperties(field *spookytypesschemas.Field, generator *ExampleGenerator) error {
    // Generate example properties based on field name
    properties := m.generateObjectPropertiesFromName(field.Name)
    
    for key, value := range properties {
        generator.writeLine("%s = %q", key, value)
    }
    
    return nil
}
```

#### Extract Helper Functions
```go
func (m *Manager) generateStringFromName(name string) string {
    switch strings.ToLower(name) {
    case "name", "title":
        return "example-name"
    case "description", "desc":
        return "Example description"
    case "hostname", "host":
        return "example.com"
    case "email":
        return "user@example.com"
    case "url", "uri":
        return "https://example.com"
    case "path", "file":
        return "/path/to/file"
    case "port":
        return "8080"
    case "version":
        return "1.0.0"
    case "id", "uuid":
        return "12345678-1234-1234-1234-123456789abc"
    default:
        return "example-value"
    }
}

func (m *Manager) generateNumberFromName(name string) interface{} {
    switch strings.ToLower(name) {
    case "port":
        return 8080
    case "timeout", "duration":
        return 30
    case "count", "size", "limit":
        return 100
    case "version":
        return 1.0
    case "percentage":
        return 75.5
    default:
        return 42
    }
}

func (m *Manager) generateBooleanFromName(name string) bool {
    switch strings.ToLower(name) {
    case "enabled", "active", "running", "online":
        return true
    case "disabled", "inactive", "stopped", "offline":
        return false
    default:
        return true
    }
}

func (m *Manager) determineArrayElementCount(field *spookytypesschemas.Field) int {
    if field.Validation != nil {
        if field.Validation.MinItems != nil && *field.Validation.MinItems > 0 {
            return *field.Validation.MinItems
        }
        if field.Validation.MaxItems != nil && *field.Validation.MaxItems < 3 {
            return *field.Validation.MaxItems
        }
    }
    return 2 // Default to 2 elements
}

func (m *Manager) generateArrayElement(field *spookytypesschemas.Field, generator *ExampleGenerator, index int) error {
    // Generate simple array elements
    switch field.Items.Type {
    case "string":
        generator.writeLine("%q,", fmt.Sprintf("item-%d", index+1))
    case "number":
        generator.writeLine("%d,", index+1)
    case "boolean":
        generator.writeLine("%v,", index%2 == 0)
    default:
        generator.writeLine("%q,", "example-item")
    }
    return nil
}

func (m *Manager) generateObjectPropertiesFromName(name string) map[string]string {
    switch strings.ToLower(name) {
    case "config", "settings":
        return map[string]string{
            "enabled": "true",
            "timeout": "30",
            "retries": "3",
        }
    case "metadata", "info":
        return map[string]string{
            "version": "1.0.0",
            "author":  "example-author",
            "created": "2024-01-01",
        }
    case "connection", "network":
        return map[string]string{
            "host": "example.com",
            "port": "8080",
            "ssl":  "true",
        }
    default:
        return map[string]string{
            "key1": "value1",
            "key2": "value2",
        }
    }
}

func (m *Manager) generatePatternExample(pattern string) string {
    // Simple pattern example generation
    switch pattern {
    case `^[a-zA-Z0-9_-]+$`:
        return "example-name"
    case `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`:
        return "user@example.com"
    case `^https?://`:
        return "https://example.com"
    case `^[0-9]+$`:
        return "12345"
    default:
        return "example-value"
    }
}
```

### Phase 3: Extract Example Assembly (Day 3)

#### Extract Example Assembly Logic
```go
func (m *Manager) assembleSchemaExample(schema *spookytypesschemas.Schema, generator *ExampleGenerator) error {
    // Write header
    generator.writeLine("# Example configuration for %s", schema.Name)
    generator.writeLine("# Generated from schema: %s", schema.Description)
    generator.writeLine("")
    
    // Write example block
    generator.writeLine("%s {", schema.Name)
    generator.Indent++
    
    // Generate field examples
    if err := m.generateFieldExamples(schema, generator); err != nil {
        return err
    }
    
    generator.Indent--
    generator.writeLine("}")
    
    return nil
}

func (m *Manager) generateFieldExamples(schema *spookytypesschemas.Schema, generator *ExampleGenerator) error {
    // Sort fields for consistent output
    fieldNames := m.getSortedFieldNames(schema.Validation.Fields)
    
    for _, fieldName := range fieldNames {
        field := schema.Validation.Fields[fieldName]
        
        // Add field comment if description exists
        if field.Description != "" {
            generator.writeLine("# %s", field.Description)
        }
        
        // Generate field example
        if err := m.generateFieldExample(field, generator); err != nil {
            return fmt.Errorf("failed to generate example for field %s: %w", fieldName, err)
        }
        
        // Add blank line between fields
        generator.writeLine("")
    }
    
    return nil
}

func (m *Manager) getSortedFieldNames(fields map[string]*spookytypesschemas.Field) []string {
    var names []string
    for name := range fields {
        names = append(names, name)
    }
    sort.Strings(names)
    return names
}
```

### Phase 4: Refactored Main Function (Day 4)

#### Final Refactored Function
```go
func (m *Manager) generateSchemaExample(schema *spookytypesschemas.Schema) (string, error) {
    // Validate schema
    if err := m.validateSchemaForExample(schema); err != nil {
        return "", err
    }
    
    // Initialize example generator
    generator := m.initializeExampleGeneration()
    
    // Assemble schema example
    if err := m.assembleSchemaExample(schema, generator); err != nil {
        return "", err
    }
    
    return generator.Buffer.String(), nil
}
```

## Complexity Reduction Metrics

### Before Refactoring
- **Cyclomatic Complexity:** 9
- **Lines of Code:** ~70
- **Nesting Levels:** 4-5
- **Responsibilities:** 4 (validation, generation, formatting, error handling)

### After Refactoring
- **Cyclomatic Complexity:** 3 (main function) + 1-3 per helper function
- **Lines of Code:** ~200 (distributed across multiple focused functions)
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
            name: "nil validation",
            schema: &spookytypesschemas.Schema{
                Validation: nil,
            },
            wantErr: true,
        },
        {
            name: "empty fields",
            schema: &spookytypesschemas.Schema{
                Validation: &spookytypesschemas.Validation{
                    Fields: map[string]*spookytypesschemas.Field{},
                },
            },
            wantErr: true,
        },
        {
            name: "valid schema",
            schema: &spookytypesschemas.Schema{
                Validation: &spookytypesschemas.Validation{
                    Fields: map[string]*spookytypesschemas.Field{
                        "name": {Type: "string"},
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

func TestGenerateStringFromName(t *testing.T) {
    tests := []struct {
        name string
        want string
    }{
        {
            name: "name",
            want: "example-name",
        },
        {
            name: "hostname",
            want: "example.com",
        },
        {
            name: "email",
            want: "user@example.com",
        },
        {
            name: "port",
            want: "8080",
        },
        {
            name: "unknown",
            want: "example-value",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := generateStringFromName(tt.name)
            if got != tt.want {
                t.Errorf("generateStringFromName() = %v, want %v", got, tt.want)
            }
        })
    }
}

func TestGenerateNumberFromName(t *testing.T) {
    tests := []struct {
        name string
        want interface{}
    }{
        {
            name: "port",
            want: 8080,
        },
        {
            name: "timeout",
            want: 30,
        },
        {
            name: "count",
            want: 100,
        },
        {
            name: "percentage",
            want: 75.5,
        },
        {
            name: "unknown",
            want: 42,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := generateNumberFromName(tt.name)
            if got != tt.want {
                t.Errorf("generateNumberFromName() = %v, want %v", got, tt.want)
            }
        })
    }
}

func TestGenerateBooleanFromName(t *testing.T) {
    tests := []struct {
        name string
        want bool
    }{
        {
            name: "enabled",
            want: true,
        },
        {
            name: "disabled",
            want: false,
        },
        {
            name: "active",
            want: true,
        },
        {
            name: "inactive",
            want: false,
        },
        {
            name: "unknown",
            want: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := generateBooleanFromName(tt.name)
            if got != tt.want {
                t.Errorf("generateBooleanFromName() = %v, want %v", got, tt.want)
            }
        })
    }
}

func TestGenerateFieldExample(t *testing.T) {
    tests := []struct {
        name    string
        field   *spookytypesschemas.Field
        want    string
        wantErr bool
    }{
        {
            name: "string field",
            field: &spookytypesschemas.Field{
                Name: "name",
                Type: "string",
            },
            want:    "name = \"example-name\"\n",
            wantErr: false,
        },
        {
            name: "number field",
            field: &spookytypesschemas.Field{
                Name: "port",
                Type: "number",
            },
            want:    "port = 8080\n",
            wantErr: false,
        },
        {
            name: "boolean field",
            field: &spookytypesschemas.Field{
                Name: "enabled",
                Type: "boolean",
            },
            want:    "enabled = true\n",
            wantErr: false,
        },
        {
            name: "unsupported type",
            field: &spookytypesschemas.Field{
                Name: "unknown",
                Type: "unknown",
            },
            want:    "",
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            generator := initializeExampleGeneration()
            err := generateFieldExample(tt.field, generator)
            if (err != nil) != tt.wantErr {
                t.Errorf("generateFieldExample() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if !tt.wantErr {
                got := generator.Buffer.String()
                if got != tt.want {
                    t.Errorf("generateFieldExample() = %v, want %v", got, tt.want)
                }
            }
        })
    }
}
```

### Integration Tests
```go
func TestGenerateSchemaExampleIntegration(t *testing.T) {
    tests := []struct {
        name    string
        schema  *spookytypesschemas.Schema
        want    string
        wantErr bool
    }{
        {
            name: "simple schema",
            schema: &spookytypesschemas.Schema{
                Name:        "test",
                Description: "Test schema",
                Validation: &spookytypesschemas.Validation{
                    Fields: map[string]*spookytypesschemas.Field{
                        "name": {
                            Type:        "string",
                            Description: "The name of the resource",
                        },
                        "port": {
                            Type:        "number",
                            Description: "The port number",
                        },
                        "enabled": {
                            Type:        "boolean",
                            Description: "Whether the resource is enabled",
                        },
                    },
                },
            },
            want: `# Example configuration for test
# Generated from schema: Test schema

# The name of the resource
name = "example-name"

# The port number
port = 8080

# Whether the resource is enabled
enabled = true

`,
            wantErr: false,
        },
        {
            name: "nil schema",
            schema: nil,
            want:    "",
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := generateSchemaExample(tt.schema)
            if (err != nil) != tt.wantErr {
                t.Errorf("generateSchemaExample() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if got != tt.want {
                t.Errorf("generateSchemaExample() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

## Implementation Timeline

### Day 1: Extract Schema Validation
- [ ] Extract `validateSchemaForExample`
- [ ] Extract `initializeExampleGeneration`
- [ ] Add unit tests for validation functions
- [ ] Verify schema validation works correctly

### Day 2: Extract Field Processing
- [ ] Extract `generateFieldExample`
- [ ] Extract field type generation functions
- [ ] Extract value generation functions
- [ ] Extract helper functions
- [ ] Add unit tests for field processing functions
- [ ] Verify field processing works correctly

### Day 3: Extract Example Assembly
- [ ] Extract `assembleSchemaExample`
- [ ] Extract `generateFieldExamples`
- [ ] Extract `getSortedFieldNames`
- [ ] Add unit tests for example assembly functions
- [ ] Verify example assembly works correctly

### Day 4: Complete Refactoring
- [ ] Refactor main `generateSchemaExample` function
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
1. **Example generation changes** - May affect example quality and consistency
2. **Field type handling changes** - May affect example generation for different types
3. **Formatting changes** - May affect example readability

### Mitigation Strategies
1. **Comprehensive testing** - Test all schema types and field combinations
2. **Example validation** - Ensure generated examples are valid and useful
3. **Formatting consistency** - Ensure consistent example formatting
4. **Gradual migration** - Implement changes incrementally with validation

## Verification

### Complexity Verification
```bash
# Verify complexity reduction
gocyclo internal/schemas/manager.go | grep generateSchemaExample
```

### Functionality Verification
```bash
# Test schema example generation
go test ./internal/schemas -run TestGenerateSchemaExample
```

### Performance Verification
- [ ] No performance regression
- [ ] Memory usage remains stable
- [ ] Example generation performance remains acceptable

## Conclusion

This refactoring will reduce the cyclomatic complexity of `generateSchemaExample` from 9 to approximately 3, while improving code maintainability, testability, and readability. The extracted functions will follow single responsibility principles and be easily testable in isolation.

**Expected Outcome:** A clean, maintainable function with clear separation of concerns and comprehensive test coverage for schema example generation operations.
