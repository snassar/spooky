# Function Improvement Plan: loadSchemasFromDirectory

**Function:** `loadSchemasFromDirectory`  
**File:** `internal/schemas/manager.go:XXX`  
**Current Complexity:** 10  
**Target Complexity:** < 8  
**Priority:** High

## Current Function Analysis

### Function Signature
```go
func (m *Manager) loadSchemasFromDirectory(dirPath string) ([]*spookytypesschemas.Schema, error)
```

### Current Issues
1. **Complex directory traversal** - Multiple conditions for file filtering and processing
2. **Mixed responsibilities** - Directory reading, file filtering, schema loading, and error handling
3. **Deep nesting** - Multiple levels of conditional logic for different file types
4. **Complex error handling** - Multiple error conditions with different handling paths
5. **Repetitive file processing patterns** - Similar logic repeated for different file types

### Complexity Breakdown
- 10 cyclomatic complexity points from:
  - 1 main function entry
  - 2 directory validation checks
  - 3 file filtering conditions
  - 2 error handling paths
  - 2 file processing conditions

## Refactoring Strategy

### Phase 1: Extract Directory Validation (Immediate - 1 day)

#### Extract Directory Validation
```go
func (m *Manager) validateDirectory(dirPath string) error {
    if dirPath == "" {
        return fmt.Errorf("directory path cannot be empty")
    }
    
    info, err := os.Stat(dirPath)
    if err != nil {
        return fmt.Errorf("failed to stat directory %s: %w", dirPath, err)
    }
    
    if !info.IsDir() {
        return fmt.Errorf("path is not a directory: %s", dirPath)
    }
    
    return nil
}
```

#### Extract File Discovery
```go
func (m *Manager) discoverSchemaFiles(dirPath string) ([]string, error) {
    entries, err := os.ReadDir(dirPath)
    if err != nil {
        return nil, fmt.Errorf("failed to read directory %s: %w", dirPath, err)
    }
    
    var schemaFiles []string
    for _, entry := range entries {
        if entry.IsDir() {
            continue
        }
        
        if m.isSchemaFile(entry.Name()) {
            schemaFiles = append(schemaFiles, filepath.Join(dirPath, entry.Name()))
        }
    }
    
    return schemaFiles, nil
}

func (m *Manager) isSchemaFile(filename string) bool {
    ext := strings.ToLower(filepath.Ext(filename))
    return ext == ".hcl" || ext == ".json" || ext == ".yaml" || ext == ".yml"
}
```

### Phase 2: Extract Schema Loading (Day 2)

#### Extract Individual Schema Loading
```go
func (m *Manager) loadSchemaFromPath(filePath string) (*spookytypesschemas.Schema, error) {
    schema, err := m.loadSchemaFromFile(filePath)
    if err != nil {
        return nil, fmt.Errorf("failed to load schema from %s: %w", filePath, err)
    }
    
    // Apply metadata if present
    if err := m.applySchemaMetadata(schema, filePath); err != nil {
        m.logger.Warn("Failed to apply metadata", map[string]interface{}{
            "file":  filePath,
            "error": err.Error(),
        })
    }
    
    return schema, nil
}

func (m *Manager) applySchemaMetadata(schema *spookytypesschemas.Schema, filePath string) error {
    // Extract metadata from file path or content
    metadata := m.extractPathMetadata(filePath)
    
    if metadata != nil {
        m.applyMetadataFields(schema, metadata)
    }
    
    return nil
}

func (m *Manager) extractPathMetadata(filePath string) map[string]interface{} {
    metadata := make(map[string]interface{})
    
    // Extract directory name as category
    dir := filepath.Dir(filePath)
    if dir != "." {
        metadata["category"] = filepath.Base(dir)
    }
    
    // Extract file name without extension as schema name
    baseName := filepath.Base(filePath)
    nameWithoutExt := strings.TrimSuffix(baseName, filepath.Ext(baseName))
    metadata["schema_name"] = nameWithoutExt
    
    return metadata
}
```

#### Extract Batch Processing
```go
func (m *Manager) processSchemaFiles(filePaths []string) ([]*spookytypesschemas.Schema, error) {
    var schemas []*spookytypesschemas.Schema
    var errors []error
    
    for _, filePath := range filePaths {
        schema, err := m.loadSchemaFromPath(filePath)
        if err != nil {
            errors = append(errors, fmt.Errorf("file %s: %w", filePath, err))
            continue
        }
        
        schemas = append(schemas, schema)
    }
    
    if len(errors) > 0 {
        return schemas, m.aggregateErrors(errors)
    }
    
    return schemas, nil
}

func (m *Manager) aggregateErrors(errors []error) error {
    if len(errors) == 1 {
        return errors[0]
    }
    
    var messages []string
    for _, err := range errors {
        messages = append(messages, err.Error())
    }
    
    return fmt.Errorf("multiple errors occurred:\n%s", strings.Join(messages, "\n"))
}
```

### Phase 3: Refactored Main Function (Day 3)

#### Final Refactored Function
```go
func (m *Manager) loadSchemasFromDirectory(dirPath string) ([]*spookytypesschemas.Schema, error) {
    // Validate directory
    if err := m.validateDirectory(dirPath); err != nil {
        return nil, err
    }
    
    // Discover schema files
    filePaths, err := m.discoverSchemaFiles(dirPath)
    if err != nil {
        return nil, err
    }
    
    if len(filePaths) == 0 {
        m.logger.Info("No schema files found in directory", map[string]interface{}{
            "directory": dirPath,
        })
        return []*spookytypesschemas.Schema{}, nil
    }
    
    // Process schema files
    return m.processSchemaFiles(filePaths)
}
```

## Complexity Reduction Metrics

### Before Refactoring
- **Cyclomatic Complexity:** 10
- **Lines of Code:** ~60
- **Nesting Levels:** 4-5
- **Responsibilities:** 4 (validation, discovery, loading, error handling)

### After Refactoring
- **Cyclomatic Complexity:** 3 (main function) + 1-3 per helper function
- **Lines of Code:** ~100 (distributed across multiple focused functions)
- **Nesting Levels:** 1-2
- **Responsibilities:** 1 per function (single responsibility principle)

## Testing Strategy

### Unit Tests for Extracted Functions
```go
func TestValidateDirectory(t *testing.T) {
    tests := []struct {
        name    string
        dirPath string
        wantErr bool
    }{
        {
            name:    "empty path",
            dirPath: "",
            wantErr: true,
        },
        {
            name:    "non-existent directory",
            dirPath: "/nonexistent/dir",
            wantErr: true,
        },
        {
            name:    "file instead of directory",
            dirPath: "testdata/test-file.txt",
            wantErr: true,
        },
        {
            name:    "valid directory",
            dirPath: "testdata",
            wantErr: false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := validateDirectory(tt.dirPath)
            if (err != nil) != tt.wantErr {
                t.Errorf("validateDirectory() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}

func TestIsSchemaFile(t *testing.T) {
    tests := []struct {
        name     string
        filename string
        want     bool
    }{
        {
            name:     "HCL file",
            filename: "test.hcl",
            want:     true,
        },
        {
            name:     "JSON file",
            filename: "test.json",
            want:     true,
        },
        {
            name:     "YAML file",
            filename: "test.yaml",
            want:     true,
        },
        {
            name:     "YML file",
            filename: "test.yml",
            want:     true,
        },
        {
            name:     "text file",
            filename: "test.txt",
            want:     false,
        },
        {
            name:     "no extension",
            filename: "test",
            want:     false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := isSchemaFile(tt.filename)
            if got != tt.want {
                t.Errorf("isSchemaFile() = %v, want %v", got, tt.want)
            }
        })
    }
}

func TestExtractPathMetadata(t *testing.T) {
    tests := []struct {
        name     string
        filePath string
        want     map[string]interface{}
    }{
        {
            name:     "file in subdirectory",
            filePath: "schemas/project/test.hcl",
            want: map[string]interface{}{
                "category":    "project",
                "schema_name": "test",
            },
        },
        {
            name:     "file in root",
            filePath: "test.hcl",
            want: map[string]interface{}{
                "schema_name": "test",
            },
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := extractPathMetadata(tt.filePath)
            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("extractPathMetadata() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

### Integration Tests
```go
func TestLoadSchemasFromDirectoryIntegration(t *testing.T) {
    tests := []struct {
        name     string
        dirPath  string
        wantCount int
        wantErr  bool
    }{
        {
            name:      "valid directory with schemas",
            dirPath:   "testdata/schemas",
            wantCount: 3,
            wantErr:   false,
        },
        {
            name:      "empty directory",
            dirPath:   "testdata/empty",
            wantCount: 0,
            wantErr:   false,
        },
        {
            name:      "non-existent directory",
            dirPath:   "testdata/nonexistent",
            wantCount: 0,
            wantErr:   true,
        },
        {
            name:      "directory with mixed files",
            dirPath:   "testdata/mixed",
            wantCount: 2,
            wantErr:   false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            schemas, err := loadSchemasFromDirectory(tt.dirPath)
            if (err != nil) != tt.wantErr {
                t.Errorf("loadSchemasFromDirectory() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if len(schemas) != tt.wantCount {
                t.Errorf("loadSchemasFromDirectory() returned %d schemas, want %d", len(schemas), tt.wantCount)
            }
        })
    }
}
```

## Implementation Timeline

### Day 1: Extract Directory Validation
- [ ] Extract `validateDirectory`
- [ ] Extract `discoverSchemaFiles`
- [ ] Extract `isSchemaFile`
- [ ] Add unit tests for directory validation functions
- [ ] Verify directory validation works correctly

### Day 2: Extract Schema Loading
- [ ] Extract `loadSchemaFromPath`
- [ ] Extract `applySchemaMetadata`
- [ ] Extract `extractPathMetadata`
- [ ] Extract `processSchemaFiles`
- [ ] Extract `aggregateErrors`
- [ ] Add unit tests for schema loading functions
- [ ] Verify schema loading works correctly

### Day 3: Complete Refactoring
- [ ] Refactor main `loadSchemasFromDirectory` function
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
- [ ] Easy to modify individual directory processing components
- [ ] Clear function names and purposes
- [ ] Well-documented functions
- [ ] Consistent error handling patterns

## Risk Mitigation

### Potential Risks
1. **Directory traversal changes** - May affect file discovery behavior
2. **File filtering changes** - May affect which files are processed
3. **Error handling changes** - May affect error propagation patterns

### Mitigation Strategies
1. **Comprehensive testing** - Test all directory scenarios and edge cases
2. **File filtering validation** - Ensure file filtering remains consistent
3. **Error handling validation** - Ensure error handling remains consistent
4. **Gradual migration** - Implement changes incrementally with validation

## Verification

### Complexity Verification
```bash
# Verify complexity reduction
gocyclo internal/schemas/manager.go | grep loadSchemasFromDirectory
```

### Functionality Verification
```bash
# Test directory schema loading
go test ./internal/schemas -run TestLoadSchemasFromDirectory
```

### Performance Verification
- [ ] No performance regression
- [ ] Memory usage remains stable
- [ ] Directory traversal performance remains acceptable

## Conclusion

This refactoring will reduce the cyclomatic complexity of `loadSchemasFromDirectory` from 10 to approximately 3, while improving code maintainability, testability, and readability. The extracted functions will follow single responsibility principles and be easily testable in isolation.

**Expected Outcome:** A clean, maintainable function with clear separation of concerns and comprehensive test coverage for directory-based schema loading operations.
