# Function Improvement Plan: Load

**Function:** `Load`  
**File:** `internal/schemas/manager.go:73`  
**Current Complexity:** 14  
**Target Complexity:** < 10  
**Priority:** High

## Current Function Analysis

### Function Signature
```go
func (m *Manager) Load(filePath string) (*spookytypesschemas.Schema, error)
```

### Current Issues
1. **Multiple validation checks** - Complex validation logic for file path and content
2. **Complex metadata extraction and application** - Multiple steps for metadata handling
3. **Mixed responsibilities** - Loading, validation, metadata processing, and registration
4. **Deep nesting** - Multiple levels of conditional logic for different operations
5. **Complex error handling** - Multiple error conditions with different handling paths

### Complexity Breakdown
- 14 cyclomatic complexity points from:
  - 1 main function entry
  - 3 validation checks (file path, content, metadata)
  - 3 metadata processing conditions
  - 3 error handling paths
  - 2 registration conditions
  - 2 type checking conditions

## Refactoring Strategy

### Phase 1: Extract Validation Logic (Immediate - 1 day)

#### Extract File Path Validation
```go
func (m *Manager) validateFilePath(filePath string) error {
    if filePath == "" {
        return fmt.Errorf("file path cannot be empty")
    }

    if _, err := os.Stat(filePath); os.IsNotExist(err) {
        return fmt.Errorf("schema file does not exist: %s", filePath)
    }

    return nil
}
```

#### Extract File Content Reading
```go
func (m *Manager) readFileContent(filePath string) ([]byte, error) {
    data, err := os.ReadFile(filePath)
    if err != nil {
        return nil, fmt.Errorf("failed to read schema file: %w", err)
    }
    return data, nil
}
```

#### Extract Base Schema Creation
```go
func (m *Manager) createBaseSchema(filePath string, data []byte) *spookytypesschemas.Schema {
    return &spookytypesschemas.Schema{
        Version:     "1.0",
        Type:        "hcl",
        Name:        filepath.Base(filePath),
        Description: fmt.Sprintf("Schema loaded from %s", filePath),
        CreatedAt:   time.Now(),
        UpdatedAt:   time.Now(),
        Content:     string(data),
        Metadata:    make(map[string]interface{}),
    }
}
```

### Phase 2: Extract Metadata Processing (Day 2)

#### Extract Metadata Application
```go
func (m *Manager) applyMetadata(schema *spookytypesschemas.Schema, data []byte) error {
    metadata, err := m.extractMetadataFromString(string(data))
    if err != nil {
        return fmt.Errorf("failed to extract metadata: %w", err)
    }

    if metadata != nil {
        m.applyMetadataFields(schema, metadata)
    }

    return nil
}

func (m *Manager) applyMetadataFields(schema *spookytypesschemas.Schema, metadata map[string]interface{}) {
    if version, ok := metadata["schema_version"].(string); ok {
        schema.Version = version
    }
    if schemaType, ok := metadata["schema_type"].(string); ok {
        schema.Type = schemaType
    }
    if name, ok := metadata["schema_name"].(string); ok {
        schema.Name = name
    }
    if description, ok := metadata["description"].(string); ok {
        schema.Description = description
    }
    if lastUpdated, ok := metadata["last_updated"].(string); ok {
        if parsed, err := time.Parse(time.RFC3339, lastUpdated); err == nil {
            schema.UpdatedAt = parsed
        }
    }
}
```

#### Extract Validation and Registration
```go
func (m *Manager) validateAndRegister(schema *spookytypesschemas.Schema, data []byte, filePath string) error {
    if err := m.validateSchemaContent(data, filePath); err != nil {
        return fmt.Errorf("schema content validation failed: %w", err)
    }

    if err := m.registry.Register(schema); err != nil {
        return fmt.Errorf("failed to register schema: %w", err)
    }

    return nil
}
```

### Phase 3: Refactored Main Function (Day 3)

#### Final Refactored Function
```go
func (m *Manager) Load(filePath string) (*spookytypesschemas.Schema, error) {
    // Validate file path
    if err := m.validateFilePath(filePath); err != nil {
        return nil, err
    }

    // Read file content
    data, err := m.readFileContent(filePath)
    if err != nil {
        return nil, err
    }

    // Create base schema
    schema := m.createBaseSchema(filePath, data)

    // Apply metadata
    if err := m.applyMetadata(schema, data); err != nil {
        return nil, err
    }

    // Validate and register
    if err := m.validateAndRegister(schema, data, filePath); err != nil {
        return nil, err
    }

    return schema, nil
}
```

## Complexity Reduction Metrics

### Before Refactoring
- **Cyclomatic Complexity:** 14
- **Lines of Code:** ~60
- **Nesting Levels:** 4-5
- **Responsibilities:** 4 (validation, loading, metadata processing, registration)

### After Refactoring
- **Cyclomatic Complexity:** 3 (main function) + 1-3 per helper function
- **Lines of Code:** ~100 (distributed across multiple focused functions)
- **Nesting Levels:** 1-2
- **Responsibilities:** 1 per function (single responsibility principle)

## Testing Strategy

### Unit Tests for Extracted Functions
```go
func TestValidateFilePath(t *testing.T) {
    tests := []struct {
        name    string
        filePath string
        wantErr bool
    }{
        {
            name:     "empty path",
            filePath: "",
            wantErr:  true,
        },
        {
            name:     "non-existent file",
            filePath: "/nonexistent/file.hcl",
            wantErr:  true,
        },
        {
            name:     "valid file path",
            filePath: "testdata/valid-schema.hcl",
            wantErr:  false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := validateFilePath(tt.filePath)
            if (err != nil) != tt.wantErr {
                t.Errorf("validateFilePath() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}

func TestCreateBaseSchema(t *testing.T) {
    filePath := "test-schema.hcl"
    data := []byte("schema { }")
    
    schema := createBaseSchema(filePath, data)
    
    if schema.Name != "test-schema.hcl" {
        t.Errorf("createBaseSchema() name = %v, want %v", schema.Name, "test-schema.hcl")
    }
    if schema.Type != "hcl" {
        t.Errorf("createBaseSchema() type = %v, want %v", schema.Type, "hcl")
    }
    if schema.Content != "schema { }" {
        t.Errorf("createBaseSchema() content = %v, want %v", schema.Content, "schema { }")
    }
}

func TestApplyMetadataFields(t *testing.T) {
    schema := &spookytypesschemas.Schema{
        Version:     "1.0",
        Type:        "hcl",
        Name:        "original",
        Description: "original description",
    }
    
    metadata := map[string]interface{}{
        "schema_version": "2.0",
        "schema_type":    "json",
        "schema_name":    "updated",
        "description":    "updated description",
    }
    
    applyMetadataFields(schema, metadata)
    
    if schema.Version != "2.0" {
        t.Errorf("applyMetadataFields() version = %v, want %v", schema.Version, "2.0")
    }
    if schema.Type != "json" {
        t.Errorf("applyMetadataFields() type = %v, want %v", schema.Type, "json")
    }
    if schema.Name != "updated" {
        t.Errorf("applyMetadataFields() name = %v, want %v", schema.Name, "updated")
    }
    if schema.Description != "updated description" {
        t.Errorf("applyMetadataFields() description = %v, want %v", schema.Description, "updated description")
    }
}
```

### Integration Tests
```go
func TestLoadIntegration(t *testing.T) {
    tests := []struct {
        name     string
        filePath string
        wantErr  bool
    }{
        {
            name:     "valid schema file",
            filePath: "testdata/valid-schema.hcl",
            wantErr:  false,
        },
        {
            name:     "invalid schema file",
            filePath: "testdata/invalid-schema.hcl",
            wantErr:  true,
        },
        {
            name:     "non-existent file",
            filePath: "testdata/nonexistent.hcl",
            wantErr:  true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            manager := NewManager(nil)
            schema, err := manager.Load(tt.filePath)
            if (err != nil) != tt.wantErr {
                t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if !tt.wantErr && schema == nil {
                t.Error("Load() returned nil schema for valid file")
            }
        })
    }
}
```

## Implementation Timeline

### Day 1: Extract Validation Logic
- [ ] Extract `validateFilePath`
- [ ] Extract `readFileContent`
- [ ] Extract `createBaseSchema`
- [ ] Add unit tests for validation functions
- [ ] Verify validation logic works correctly

### Day 2: Extract Metadata Processing
- [ ] Extract `applyMetadata`
- [ ] Extract `applyMetadataFields`
- [ ] Extract `validateAndRegister`
- [ ] Add unit tests for metadata processing functions
- [ ] Verify metadata processing works correctly

### Day 3: Complete Refactoring
- [ ] Refactor main `Load` function
- [ ] Add integration tests
- [ ] Verify complexity reduction with gocyclo
- [ ] Code review and documentation
- [ ] Performance testing

## Success Criteria

### Complexity Reduction
- [ ] Main function complexity < 10
- [ ] All extracted functions complexity < 5
- [ ] No function exceeds complexity threshold

### Code Quality
- [ ] Single responsibility principle maintained
- [ ] Clear separation of concerns
- [ ] Comprehensive test coverage (>90%)
- [ ] No regression in functionality

### Maintainability
- [ ] Easy to modify individual loading components
- [ ] Clear function names and purposes
- [ ] Well-documented functions
- [ ] Consistent error handling patterns

## Risk Mitigation

### Potential Risks
1. **File I/O changes** - May affect file handling behavior
2. **Metadata processing changes** - May affect schema metadata handling
3. **Error handling changes** - May affect error propagation patterns

### Mitigation Strategies
1. **Comprehensive testing** - Test all file I/O scenarios
2. **Metadata validation** - Ensure metadata processing remains consistent
3. **Error handling validation** - Ensure error handling remains consistent
4. **Gradual migration** - Implement changes incrementally with validation

## Verification

### Complexity Verification
```bash
# Verify complexity reduction
gocyclo internal/schemas/manager.go | grep "Load "
```

### Functionality Verification
```bash
# Test schema loading
go test ./internal/schemas -run TestLoad
```

### Performance Verification
- [ ] No performance regression
- [ ] Memory usage remains stable
- [ ] File I/O performance remains acceptable

## Conclusion

This refactoring will reduce the cyclomatic complexity of `Load` from 14 to approximately 3, while improving code maintainability, testability, and readability. The extracted functions will follow single responsibility principles and be easily testable in isolation.

**Expected Outcome:** A clean, maintainable function with clear separation of concerns and comprehensive test coverage for schema loading operations.
