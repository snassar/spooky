# Function Improvement Plan: extractMetadataFromString

**Function:** `extractMetadataFromString`  
**File:** `internal/schemas/manager.go:1192`  
**Current Complexity:** 15  
**Target Complexity:** < 10  
**Priority:** Critical

## Current Function Analysis

### Function Signature
```go
func (m *Manager) extractMetadataFromString(content string) (map[string]interface{}, error)
```

### Current Issues
1. **Complex string parsing logic** - Multiple nested conditions for brace counting and content extraction
2. **Mixed responsibilities** - Parsing, validation, extraction, and error handling
3. **Deep nesting** - Multiple levels of conditional logic for string manipulation
4. **Complex error handling** - Multiple error conditions with different handling paths
5. **Repetitive parsing patterns** - Similar logic repeated for different metadata fields

### Complexity Breakdown
- 15 cyclomatic complexity points from:
  - 1 main function entry
  - 3 metadata block boundary checks
  - 4 brace counting loop conditions
  - 3 key-value parsing conditions
  - 2 error handling paths
  - 2 string manipulation conditions

## Refactoring Strategy

### Phase 1: Extract Block Boundary Detection (Immediate - 1 day)

#### Extract Metadata Block Detection
```go
func (m *Manager) findMetadataBlock(content string) (int, int, error) {
    metadataStart := strings.Index(content, "metadata {")
    if metadataStart == -1 {
        return -1, -1, nil // No metadata block found
    }

    metadataEnd, err := m.findBlockEnd(content, metadataStart)
    if err != nil {
        return -1, -1, fmt.Errorf("failed to find metadata block end: %w", err)
    }

    return metadataStart, metadataEnd, nil
}

func (m *Manager) findBlockEnd(content string, start int) (int, error) {
    braceCount := 0
    for i := start; i < len(content); i++ {
        if content[i] == '{' {
            braceCount++
        } else if content[i] == '}' {
            braceCount--
            if braceCount == 0 {
                return i + 1, nil
            }
        }
    }
    return -1, fmt.Errorf("unclosed metadata block")
}
```

#### Extract Content Extraction
```go
func (m *Manager) extractMetadataContent(content string, start, end int) string {
    return content[start:end]
}
```

### Phase 2: Extract Key-Value Parsing (Day 2)

#### Extract Line Processing
```go
func (m *Manager) parseMetadataLines(metadataContent string) (map[string]interface{}, error) {
    metadata := make(map[string]interface{})
    lines := strings.Split(metadataContent, "\n")

    for _, line := range lines {
        line = strings.TrimSpace(line)
        if m.isMetadataBoundary(line) {
            continue
        }

        if key, value, ok := m.parseKeyValue(line); ok {
            metadata[key] = value
        }
    }

    return metadata, nil
}

func (m *Manager) isMetadataBoundary(line string) bool {
    return line == "" || line == "metadata {" || line == "}"
}

func (m *Manager) parseKeyValue(line string) (string, string, bool) {
    if !strings.Contains(line, "=") {
        return "", "", false
    }

    parts := strings.SplitN(line, "=", 2)
    if len(parts) != 2 {
        return "", "", false
    }

    key := strings.TrimSpace(parts[0])
    value := strings.TrimSpace(parts[1])

    // Remove quotes if present
    value = m.removeQuotes(value)

    return key, value, true
}

func (m *Manager) removeQuotes(value string) string {
    if strings.HasPrefix(value, `"`) && strings.HasSuffix(value, `"`) {
        return value[1 : len(value)-1]
    }
    return value
}
```

### Phase 3: Refactored Main Function (Day 3)

#### Final Refactored Function
```go
func (m *Manager) extractMetadataFromString(content string) (map[string]interface{}, error) {
    // Find metadata block boundaries
    start, end, err := m.findMetadataBlock(content)
    if err != nil {
        return nil, err
    }

    if start == -1 {
        return nil, nil // No metadata block found
    }

    // Extract metadata content
    metadataContent := m.extractMetadataContent(content, start, end)

    // Parse metadata content
    return m.parseMetadataLines(metadataContent)
}
```

## Complexity Reduction Metrics

### Before Refactoring
- **Cyclomatic Complexity:** 15
- **Lines of Code:** ~50
- **Nesting Levels:** 4-5
- **Responsibilities:** 4 (boundary detection, content extraction, parsing, error handling)

### After Refactoring
- **Cyclomatic Complexity:** 3 (main function) + 1-3 per helper function
- **Lines of Code:** ~80 (distributed across multiple focused functions)
- **Nesting Levels:** 1-2
- **Responsibilities:** 1 per function (single responsibility principle)

## Testing Strategy

### Unit Tests for Extracted Functions
```go
func TestFindMetadataBlock(t *testing.T) {
    tests := []struct {
        name        string
        content     string
        wantStart   int
        wantEnd     int
        wantErr     bool
    }{
        {
            name:      "valid metadata block",
            content:   "metadata { key = value }",
            wantStart: 0,
            wantEnd:   25,
            wantErr:   false,
        },
        {
            name:      "no metadata block",
            content:   "no metadata here",
            wantStart: -1,
            wantEnd:   -1,
            wantErr:   false,
        },
        {
            name:      "unclosed metadata block",
            content:   "metadata { key = value",
            wantStart: -1,
            wantEnd:   -1,
            wantErr:   true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            start, end, err := findMetadataBlock(tt.content)
            if (err != nil) != tt.wantErr {
                t.Errorf("findMetadataBlock() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if start != tt.wantStart || end != tt.wantEnd {
                t.Errorf("findMetadataBlock() = (%v, %v), want (%v, %v)", start, end, tt.wantStart, tt.wantEnd)
            }
        })
    }
}

func TestParseKeyValue(t *testing.T) {
    tests := []struct {
        name     string
        line     string
        wantKey  string
        wantVal  string
        wantOk   bool
    }{
        {
            name:    "valid key value",
            line:    "key = value",
            wantKey: "key",
            wantVal: "value",
            wantOk:  true,
        },
        {
            name:    "quoted value",
            line:    "key = \"quoted value\"",
            wantKey: "key",
            wantVal: "quoted value",
            wantOk:  true,
        },
        {
            name:    "no equals sign",
            line:    "key value",
            wantKey: "",
            wantVal: "",
            wantOk:  false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            key, val, ok := parseKeyValue(tt.line)
            if ok != tt.wantOk {
                t.Errorf("parseKeyValue() ok = %v, want %v", ok, tt.wantOk)
                return
            }
            if key != tt.wantKey || val != tt.wantVal {
                t.Errorf("parseKeyValue() = (%v, %v), want (%v, %v)", key, val, tt.wantKey, tt.wantVal)
            }
        })
    }
}
```

### Integration Tests
```go
func TestExtractMetadataFromStringIntegration(t *testing.T) {
    tests := []struct {
        name    string
        content string
        want    map[string]interface{}
        wantErr bool
    }{
        {
            name: "valid metadata",
            content: `
metadata {
    version = "1.0"
    description = "Test schema"
    author = "test"
}
`,
            want: map[string]interface{}{
                "version":     "1.0",
                "description": "Test schema",
                "author":      "test",
            },
            wantErr: false,
        },
        {
            name:    "no metadata",
            content: "no metadata here",
            want:    nil,
            wantErr: false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := extractMetadataFromString(tt.content)
            if (err != nil) != tt.wantErr {
                t.Errorf("extractMetadataFromString() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("extractMetadataFromString() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

## Implementation Timeline

### Day 1: Extract Block Boundary Detection
- [ ] Extract `findMetadataBlock`
- [ ] Extract `findBlockEnd`
- [ ] Extract `extractMetadataContent`
- [ ] Add unit tests for boundary detection functions
- [ ] Verify boundary detection works correctly

### Day 2: Extract Key-Value Parsing
- [ ] Extract `parseMetadataLines`
- [ ] Extract `isMetadataBoundary`
- [ ] Extract `parseKeyValue`
- [ ] Extract `removeQuotes`
- [ ] Add unit tests for parsing functions
- [ ] Verify parsing works correctly

### Day 3: Complete Refactoring
- [ ] Refactor main `extractMetadataFromString` function
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
- [ ] Easy to modify individual parsing components
- [ ] Clear function names and purposes
- [ ] Well-documented functions
- [ ] Consistent error handling patterns

## Risk Mitigation

### Potential Risks
1. **String parsing edge cases** - May miss edge cases in metadata parsing
2. **Performance impact** - Multiple function calls may impact performance
3. **Error handling changes** - May affect error propagation patterns

### Mitigation Strategies
1. **Comprehensive testing** - Test all edge cases and error conditions
2. **Performance benchmarking** - Measure performance impact and optimize if needed
3. **Error handling validation** - Ensure error handling remains consistent
4. **Gradual migration** - Implement changes incrementally with validation

## Verification

### Complexity Verification
```bash
# Verify complexity reduction
gocyclo internal/schemas/manager.go | grep extractMetadataFromString
```

### Functionality Verification
```bash
# Test with various metadata formats
go test ./internal/schemas -run TestExtractMetadataFromString
```

### Performance Verification
- [ ] No performance regression
- [ ] Memory usage remains stable
- [ ] String parsing performance remains acceptable

## Conclusion

This refactoring will reduce the cyclomatic complexity of `extractMetadataFromString` from 15 to approximately 3, while improving code maintainability, testability, and readability. The extracted functions will follow single responsibility principles and be easily testable in isolation.

**Expected Outcome:** A clean, maintainable function with clear separation of concerns and comprehensive test coverage for metadata parsing operations.
