# Function Improvement Plan: loadSchemaFromFile

**Function:** `loadSchemaFromFile`  
**File:** `internal/schemas/manager.go:XXX`  
**Current Complexity:** 12  
**Target Complexity:** < 8  
**Priority:** High

## Current Function Analysis

### Function Signature
```go
func (m *Manager) loadSchemaFromFile(filePath string) (*spookytypesschemas.Schema, error)
```

### Current Issues
1. **Complex file type detection** - Multiple conditions for different file extensions
2. **Mixed responsibilities** - File loading, type detection, validation, and error handling
3. **Deep nesting** - Multiple levels of conditional logic for different file types
4. **Complex error handling** - Multiple error conditions with different handling paths
5. **Repetitive validation patterns** - Similar validation logic for different file types

### Complexity Breakdown
- 12 cyclomatic complexity points from:
  - 1 main function entry
  - 3 file type detection conditions
  - 3 validation checks for different file types
  - 3 error handling paths
  - 2 file processing conditions

## Refactoring Strategy

### Phase 1: Extract File Type Detection (Immediate - 1 day)

#### Extract File Type Detection
```go
func (m *Manager) detectFileType(filePath string) (string, error) {
    ext := strings.ToLower(filepath.Ext(filePath))
    
    switch ext {
    case ".hcl":
        return "hcl", nil
    case ".json":
        return "json", nil
    case ".yaml", ".yml":
        return "yaml", nil
    default:
        return "", fmt.Errorf("unsupported file type: %s", ext)
    }
}
```

#### Extract File Content Reading
```go
func (m *Manager) readFileContent(filePath string) ([]byte, error) {
    data, err := os.ReadFile(filePath)
    if err != nil {
        return nil, fmt.Errorf("failed to read file %s: %w", filePath, err)
    }
    return data, nil
}
```

#### Extract File Validation
```go
func (m *Manager) validateFile(filePath string) error {
    if filePath == "" {
        return fmt.Errorf("file path cannot be empty")
    }
    
    if _, err := os.Stat(filePath); os.IsNotExist(err) {
        return fmt.Errorf("file does not exist: %s", filePath)
    }
    
    return nil
}
```

### Phase 2: Extract Schema Creation (Day 2)

#### Extract Schema Creation by Type
```go
func (m *Manager) createSchemaFromContent(filePath, fileType string, content []byte) (*spookytypesschemas.Schema, error) {
    switch fileType {
    case "hcl":
        return m.createHCLSchema(filePath, content)
    case "json":
        return m.createJSONSchema(filePath, content)
    case "yaml":
        return m.createYAMLSchema(filePath, content)
    default:
        return nil, fmt.Errorf("unsupported file type: %s", fileType)
    }
}

func (m *Manager) createHCLSchema(filePath string, content []byte) (*spookytypesschemas.Schema, error) {
    schema := &spookytypesschemas.Schema{
        Name:        filepath.Base(filePath),
        Type:        "hcl",
        Version:     "1.0",
        Content:     string(content),
        CreatedAt:   time.Now(),
        UpdatedAt:   time.Now(),
        Metadata:    make(map[string]interface{}),
    }
    
    // Parse HCL content
    if err := m.parseHCLContent(schema, content); err != nil {
        return nil, fmt.Errorf("failed to parse HCL content: %w", err)
    }
    
    return schema, nil
}

func (m *Manager) createJSONSchema(filePath string, content []byte) (*spookytypesschemas.Schema, error) {
    schema := &spookytypesschemas.Schema{
        Name:        filepath.Base(filePath),
        Type:        "json",
        Version:     "1.0",
        Content:     string(content),
        CreatedAt:   time.Now(),
        UpdatedAt:   time.Now(),
        Metadata:    make(map[string]interface{}),
    }
    
    // Parse JSON content
    if err := m.parseJSONContent(schema, content); err != nil {
        return nil, fmt.Errorf("failed to parse JSON content: %w", err)
    }
    
    return schema, nil
}

func (m *Manager) createYAMLSchema(filePath string, content []byte) (*spookytypesschemas.Schema, error) {
    schema := &spookytypesschemas.Schema{
        Name:        filepath.Base(filePath),
        Type:        "yaml",
        Version:     "1.0",
        Content:     string(content),
        CreatedAt:   time.Now(),
        UpdatedAt:   time.Now(),
        Metadata:    make(map[string]interface{}),
    }
    
    // Parse YAML content
    if err := m.parseYAMLContent(schema, content); err != nil {
        return nil, fmt.Errorf("failed to parse YAML content: %w", err)
    }
    
    return schema, nil
}
```

### Phase 3: Refactored Main Function (Day 3)

#### Final Refactored Function
```go
func (m *Manager) loadSchemaFromFile(filePath string) (*spookytypesschemas.Schema, error) {
    // Validate file
    if err := m.validateFile(filePath); err != nil {
        return nil, err
    }
    
    // Detect file type
    fileType, err := m.detectFileType(filePath)
    if err != nil {
        return nil, err
    }
    
    // Read file content
    content, err := m.readFileContent(filePath)
    if err != nil {
        return nil, err
    }
    
    // Create schema from content
    return m.createSchemaFromContent(filePath, fileType, content)
}
```

## Complexity Reduction Metrics

### Before Refactoring
- **Cyclomatic Complexity:** 12
- **Lines of Code:** ~50
- **Nesting Levels:** 4-5
- **Responsibilities:** 4 (validation, type detection, loading, creation)

### After Refactoring
- **Cyclomatic Complexity:** 3 (main function) + 1-3 per helper function
- **Lines of Code:** ~80 (distributed across multiple focused functions)
- **Nesting Levels:** 1-2
- **Responsibilities:** 1 per function (single responsibility principle)

## Testing Strategy

### Unit Tests for Extracted Functions
```go
func TestDetectFileType(t *testing.T) {
    tests := []struct {
        name     string
        filePath string
        want     string
        wantErr  bool
    }{
        {
            name:     "HCL file",
            filePath: "test.hcl",
            want:     "hcl",
            wantErr:  false,
        },
        {
            name:     "JSON file",
            filePath: "test.json",
            want:     "json",
            wantErr:  false,
        },
        {
            name:     "YAML file",
            filePath: "test.yaml",
            want:     "yaml",
            wantErr:  false,
        },
        {
            name:     "unsupported file",
            filePath: "test.txt",
            want:     "",
            wantErr:  true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := detectFileType(tt.filePath)
            if (err != nil) != tt.wantErr {
                t.Errorf("detectFileType() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if got != tt.want {
                t.Errorf("detectFileType() = %v, want %v", got, tt.want)
            }
        })
    }
}

func TestValidateFile(t *testing.T) {
    tests := []struct {
        name     string
        filePath string
        wantErr  bool
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
            name:     "valid file",
            filePath: "testdata/valid-schema.hcl",
            wantErr:  false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := validateFile(tt.filePath)
            if (err != nil) != tt.wantErr {
                t.Errorf("validateFile() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}

func TestCreateHCLSchema(t *testing.T) {
    filePath := "test.hcl"
    content := []byte("schema { }")
    
    schema, err := createHCLSchema(filePath, content)
    if err != nil {
        t.Errorf("createHCLSchema() error = %v", err)
        return
    }
    
    if schema.Name != "test.hcl" {
        t.Errorf("createHCLSchema() name = %v, want %v", schema.Name, "test.hcl")
    }
    if schema.Type != "hcl" {
        t.Errorf("createHCLSchema() type = %v, want %v", schema.Type, "hcl")
    }
    if schema.Content != "schema { }" {
        t.Errorf("createHCLSchema() content = %v, want %v", schema.Content, "schema { }")
    }
}
```

### Integration Tests
```go
func TestLoadSchemaFromFileIntegration(t *testing.T) {
    tests := []struct {
        name     string
        filePath string
        wantErr  bool
    }{
        {
            name:     "valid HCL file",
            filePath: "testdata/valid-schema.hcl",
            wantErr:  false,
        },
        {
            name:     "valid JSON file",
            filePath: "testdata/valid-schema.json",
            wantErr:  false,
        },
        {
            name:     "valid YAML file",
            filePath: "testdata/valid-schema.yaml",
            wantErr:  false,
        },
        {
            name:     "unsupported file type",
            filePath: "testdata/invalid-schema.txt",
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
            schema, err := loadSchemaFromFile(tt.filePath)
            if (err != nil) != tt.wantErr {
                t.Errorf("loadSchemaFromFile() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if !tt.wantErr && schema == nil {
                t.Error("loadSchemaFromFile() returned nil schema for valid file")
            }
        })
    }
}
```

## Implementation Timeline

### Day 1: Extract File Type Detection
- [ ] Extract `detectFileType`
- [ ] Extract `readFileContent`
- [ ] Extract `validateFile`
- [ ] Add unit tests for file type detection functions
- [ ] Verify file type detection works correctly

### Day 2: Extract Schema Creation
- [ ] Extract `createSchemaFromContent`
- [ ] Extract `createHCLSchema`
- [ ] Extract `createJSONSchema`
- [ ] Extract `createYAMLSchema`
- [ ] Add unit tests for schema creation functions
- [ ] Verify schema creation works correctly

### Day 3: Complete Refactoring
- [ ] Refactor main `loadSchemaFromFile` function
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
- [ ] Easy to modify individual file type handlers
- [ ] Clear function names and purposes
- [ ] Well-documented functions
- [ ] Consistent error handling patterns

## Risk Mitigation

### Potential Risks
1. **File type detection changes** - May affect supported file types
2. **Content parsing changes** - May affect schema parsing behavior
3. **Error handling changes** - May affect error propagation patterns

### Mitigation Strategies
1. **Comprehensive testing** - Test all file types and edge cases
2. **Content validation** - Ensure content parsing remains consistent
3. **Error handling validation** - Ensure error handling remains consistent
4. **Gradual migration** - Implement changes incrementally with validation

## Verification

### Complexity Verification
```bash
# Verify complexity reduction
gocyclo internal/schemas/manager.go | grep loadSchemaFromFile
```

### Functionality Verification
```bash
# Test schema loading from different file types
go test ./internal/schemas -run TestLoadSchemaFromFile
```

### Performance Verification
- [ ] No performance regression
- [ ] Memory usage remains stable
- [ ] File I/O performance remains acceptable

## Conclusion

This refactoring will reduce the cyclomatic complexity of `loadSchemaFromFile` from 12 to approximately 3, while improving code maintainability, testability, and readability. The extracted functions will follow single responsibility principles and be easily testable in isolation.

**Expected Outcome:** A clean, maintainable function with clear separation of concerns and comprehensive test coverage for schema file loading operations.
