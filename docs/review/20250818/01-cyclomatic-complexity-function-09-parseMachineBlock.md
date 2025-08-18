# Function Improvement Plan: parseMachineBlock

**Function:** `parseMachineBlock`  
**File:** `internal/machines/loader.go:XXX`  
**Current Complexity:** 9  
**Target Complexity:** < 8  
**Priority:** Medium

## Current Function Analysis

### Function Signature
```go
func (l *Loader) parseMachineBlock(block *hcl.Block) (*spookytypesmachines.Machine, error)
```

### Current Issues
1. **Complex block parsing logic** - Multiple conditions for different block attributes and nested blocks
2. **Mixed responsibilities** - Block parsing, validation, machine creation, and error handling
3. **Deep nesting** - Multiple levels of conditional logic for different block components
4. **Complex attribute extraction** - Multiple conditions for required and optional attributes
5. **Repetitive validation patterns** - Similar logic repeated for different attribute types

### Complexity Breakdown
- 9 cyclomatic complexity points from:
  - 1 main function entry
  - 2 block validation checks
  - 3 attribute extraction conditions
  - 2 error handling paths
  - 1 nested block processing condition

## Refactoring Strategy

### Phase 1: Extract Block Validation (Immediate - 1 day)

#### Extract Block Validation
```go
func (l *Loader) validateMachineBlock(block *hcl.Block) error {
    if block == nil {
        return fmt.Errorf("block cannot be nil")
    }
    
    if block.Type != "machine" {
        return fmt.Errorf("expected machine block, got %s", block.Type)
    }
    
    if len(block.Labels) == 0 {
        return fmt.Errorf("machine block must have a name label")
    }
    
    return nil
}
```

#### Extract Machine Name Extraction
```go
func (l *Loader) extractMachineName(block *hcl.Block) (string, error) {
    if len(block.Labels) == 0 {
        return "", fmt.Errorf("machine block missing name label")
    }
    
    name := block.Labels[0]
    if name == "" {
        return "", fmt.Errorf("machine name cannot be empty")
    }
    
    return name, nil
}
```

### Phase 2: Extract Attribute Processing (Day 2)

#### Extract Required Attribute Processing
```go
func (l *Loader) extractRequiredAttributes(block *hcl.Block) (map[string]interface{}, error) {
    requiredAttrs := []string{"hostname", "port", "user"}
    attributes := make(map[string]interface{})
    
    for _, attrName := range requiredAttrs {
        attr, exists := block.Body.Attributes[attrName]
        if !exists {
            return nil, fmt.Errorf("required attribute '%s' missing", attrName)
        }
        
        value, err := l.extractAttributeValue(attr)
        if err != nil {
            return nil, fmt.Errorf("failed to extract %s: %w", attrName, err)
        }
        
        attributes[attrName] = value
    }
    
    return attributes, nil
}

func (l *Loader) extractOptionalAttributes(block *hcl.Block) map[string]interface{} {
    optionalAttrs := []string{"description", "tags", "metadata"}
    attributes := make(map[string]interface{})
    
    for _, attrName := range optionalAttrs {
        if attr, exists := block.Body.Attributes[attrName]; exists {
            if value, err := l.extractAttributeValue(attr); err == nil {
                attributes[attrName] = value
            }
        }
    }
    
    return attributes
}

func (l *Loader) extractAttributeValue(attr *hcl.Attribute) (interface{}, error) {
    switch attr.Expr.(type) {
    case *hcl.LiteralValueExpr:
        return l.extractLiteralValue(attr.Expr.(*hcl.LiteralValueExpr))
    case *hcl.TupleConsExpr:
        return l.extractTupleValue(attr.Expr.(*hcl.TupleConsExpr))
    case *hcl.ObjectConsExpr:
        return l.extractObjectValue(attr.Expr.(*hcl.ObjectConsExpr))
    default:
        return nil, fmt.Errorf("unsupported expression type: %T", attr.Expr)
    }
}

func (l *Loader) extractLiteralValue(expr *hcl.LiteralValueExpr) (interface{}, error) {
    switch expr.Val.Type() {
    case cty.String:
        return expr.Val.AsString(), nil
    case cty.Number:
        return expr.Val.AsBigFloat().String(), nil
    case cty.Bool:
        return expr.Val.True(), nil
    default:
        return nil, fmt.Errorf("unsupported literal type: %s", expr.Val.Type())
    }
}

func (l *Loader) extractTupleValue(expr *hcl.TupleConsExpr) ([]interface{}, error) {
    var values []interface{}
    
    for _, elem := range expr.Exprs {
        if literal, ok := elem.(*hcl.LiteralValueExpr); ok {
            value, err := l.extractLiteralValue(literal)
            if err != nil {
                return nil, err
            }
            values = append(values, value)
        } else {
            return nil, fmt.Errorf("unsupported tuple element type: %T", elem)
        }
    }
    
    return values, nil
}

func (l *Loader) extractObjectValue(expr *hcl.ObjectConsExpr) (map[string]interface{}, error) {
    values := make(map[string]interface{})
    
    for _, item := range expr.Items {
        key, err := l.extractObjectKey(item.KeyExpr)
        if err != nil {
            return nil, err
        }
        
        value, err := l.extractObjectValue(item.ValExpr)
        if err != nil {
            return nil, err
        }
        
        values[key] = value
    }
    
    return values, nil
}

func (l *Loader) extractObjectKey(expr hcl.Expression) (string, error) {
    if literal, ok := expr.(*hcl.LiteralValueExpr); ok {
        if literal.Val.Type() == cty.String {
            return literal.Val.AsString(), nil
        }
    }
    return "", fmt.Errorf("object key must be a string literal")
}

func (l *Loader) extractObjectValue(expr hcl.Expression) (interface{}, error) {
    switch v := expr.(type) {
    case *hcl.LiteralValueExpr:
        return l.extractLiteralValue(v)
    case *hcl.TupleConsExpr:
        return l.extractTupleValue(v)
    case *hcl.ObjectConsExpr:
        return l.extractObjectValue(v)
    default:
        return nil, fmt.Errorf("unsupported object value type: %T", expr)
    }
}
```

#### Extract Nested Block Processing
```go
func (l *Loader) extractNestedBlocks(block *hcl.Block) (map[string]interface{}, error) {
    nestedBlocks := make(map[string]interface{})
    
    for _, nestedBlock := range block.Body.Blocks {
        switch nestedBlock.Type {
        case "authentication":
            auth, err := l.extractAuthenticationBlock(nestedBlock)
            if err != nil {
                return nil, fmt.Errorf("failed to extract authentication: %w", err)
            }
            nestedBlocks["authentication"] = auth
        case "connection":
            conn, err := l.extractConnectionBlock(nestedBlock)
            if err != nil {
                return nil, fmt.Errorf("failed to extract connection: %w", err)
            }
            nestedBlocks["connection"] = conn
        default:
            return nil, fmt.Errorf("unsupported nested block type: %s", nestedBlock.Type)
        }
    }
    
    return nestedBlocks, nil
}

func (l *Loader) extractAuthenticationBlock(block *hcl.Block) (map[string]interface{}, error) {
    auth := make(map[string]interface{})
    
    for name, attr := range block.Body.Attributes {
        value, err := l.extractAttributeValue(attr)
        if err != nil {
            return nil, fmt.Errorf("failed to extract authentication.%s: %w", name, err)
        }
        auth[name] = value
    }
    
    return auth, nil
}

func (l *Loader) extractConnectionBlock(block *hcl.Block) (map[string]interface{}, error) {
    conn := make(map[string]interface{})
    
    for name, attr := range block.Body.Attributes {
        value, err := l.extractAttributeValue(attr)
        if err != nil {
            return nil, fmt.Errorf("failed to extract connection.%s: %w", name, err)
        }
        conn[name] = value
    }
    
    return conn, nil
}
```

### Phase 3: Extract Machine Creation (Day 3)

#### Extract Machine Creation Logic
```go
func (l *Loader) createMachineFromData(name string, requiredAttrs, optionalAttrs, nestedBlocks map[string]interface{}) (*spookytypesmachines.Machine, error) {
    machine := &spookytypesmachines.Machine{
        Name:        name,
        Hostname:    requiredAttrs["hostname"].(string),
        Port:        l.extractPort(requiredAttrs["port"]),
        User:        requiredAttrs["user"].(string),
        Description: l.extractString(optionalAttrs["description"]),
        Tags:        l.extractTags(optionalAttrs["tags"]),
        Metadata:    l.extractMetadata(optionalAttrs["metadata"]),
    }
    
    // Apply nested blocks
    if auth, exists := nestedBlocks["authentication"]; exists {
        machine.Authentication = l.createAuthentication(auth.(map[string]interface{}))
    }
    
    if conn, exists := nestedBlocks["connection"]; exists {
        machine.Connection = l.createConnection(conn.(map[string]interface{}))
    }
    
    return machine, nil
}

func (l *Loader) extractPort(value interface{}) int {
    if str, ok := value.(string); ok {
        if port, err := strconv.Atoi(str); err == nil {
            return port
        }
    }
    return 22 // Default SSH port
}

func (l *Loader) extractString(value interface{}) string {
    if str, ok := value.(string); ok {
        return str
    }
    return ""
}

func (l *Loader) extractTags(value interface{}) []string {
    if tags, ok := value.([]interface{}); ok {
        var result []string
        for _, tag := range tags {
            if str, ok := tag.(string); ok {
                result = append(result, str)
            }
        }
        return result
    }
    return nil
}

func (l *Loader) extractMetadata(value interface{}) map[string]interface{} {
    if metadata, ok := value.(map[string]interface{}); ok {
        return metadata
    }
    return nil
}

func (l *Loader) createAuthentication(auth map[string]interface{}) *spookytypesmachines.Authentication {
    return &spookytypesmachines.Authentication{
        Method:  l.extractString(auth["method"]),
        KeyPath: l.extractString(auth["key_path"]),
        Password: l.extractString(auth["password"]),
    }
}

func (l *Loader) createConnection(conn map[string]interface{}) *spookytypesmachines.Connection {
    return &spookytypesmachines.Connection{
        Timeout: l.extractInt(conn["timeout"]),
        Retries: l.extractInt(conn["retries"]),
    }
}

func (l *Loader) extractInt(value interface{}) int {
    if str, ok := value.(string); ok {
        if i, err := strconv.Atoi(str); err == nil {
            return i
        }
    }
    return 0
}
```

### Phase 4: Refactored Main Function (Day 4)

#### Final Refactored Function
```go
func (l *Loader) parseMachineBlock(block *hcl.Block) (*spookytypesmachines.Machine, error) {
    // Validate block
    if err := l.validateMachineBlock(block); err != nil {
        return nil, err
    }
    
    // Extract machine name
    name, err := l.extractMachineName(block)
    if err != nil {
        return nil, err
    }
    
    // Extract required attributes
    requiredAttrs, err := l.extractRequiredAttributes(block)
    if err != nil {
        return nil, err
    }
    
    // Extract optional attributes
    optionalAttrs := l.extractOptionalAttributes(block)
    
    // Extract nested blocks
    nestedBlocks, err := l.extractNestedBlocks(block)
    if err != nil {
        return nil, err
    }
    
    // Create machine
    return l.createMachineFromData(name, requiredAttrs, optionalAttrs, nestedBlocks)
}
```

## Complexity Reduction Metrics

### Before Refactoring
- **Cyclomatic Complexity:** 9
- **Lines of Code:** ~60
- **Nesting Levels:** 4-5
- **Responsibilities:** 4 (validation, extraction, creation, error handling)

### After Refactoring
- **Cyclomatic Complexity:** 3 (main function) + 1-3 per helper function
- **Lines of Code:** ~200 (distributed across multiple focused functions)
- **Nesting Levels:** 1-2
- **Responsibilities:** 1 per function (single responsibility principle)

## Testing Strategy

### Unit Tests for Extracted Functions
```go
func TestValidateMachineBlock(t *testing.T) {
    tests := []struct {
        name    string
        block   *hcl.Block
        wantErr bool
    }{
        {
            name:    "nil block",
            block:   nil,
            wantErr: true,
        },
        {
            name: "wrong block type",
            block: &hcl.Block{
                Type:   "wrong",
                Labels: []string{"test"},
            },
            wantErr: true,
        },
        {
            name: "missing name label",
            block: &hcl.Block{
                Type:   "machine",
                Labels: []string{},
            },
            wantErr: true,
        },
        {
            name: "valid block",
            block: &hcl.Block{
                Type:   "machine",
                Labels: []string{"test-machine"},
            },
            wantErr: false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := validateMachineBlock(tt.block)
            if (err != nil) != tt.wantErr {
                t.Errorf("validateMachineBlock() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}

func TestExtractMachineName(t *testing.T) {
    tests := []struct {
        name    string
        block   *hcl.Block
        want    string
        wantErr bool
    }{
        {
            name: "valid name",
            block: &hcl.Block{
                Labels: []string{"test-machine"},
            },
            want:    "test-machine",
            wantErr: false,
        },
        {
            name: "empty name",
            block: &hcl.Block{
                Labels: []string{""},
            },
            want:    "",
            wantErr: true,
        },
        {
            name: "no labels",
            block: &hcl.Block{
                Labels: []string{},
            },
            want:    "",
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := extractMachineName(tt.block)
            if (err != nil) != tt.wantErr {
                t.Errorf("extractMachineName() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if got != tt.want {
                t.Errorf("extractMachineName() = %v, want %v", got, tt.want)
            }
        })
    }
}

func TestExtractLiteralValue(t *testing.T) {
    tests := []struct {
        name    string
        expr    *hcl.LiteralValueExpr
        want    interface{}
        wantErr bool
    }{
        {
            name: "string value",
            expr: &hcl.LiteralValueExpr{
                Val: cty.StringVal("test"),
            },
            want:    "test",
            wantErr: false,
        },
        {
            name: "number value",
            expr: &hcl.LiteralValueExpr{
                Val: cty.NumberIntVal(42),
            },
            want:    "42",
            wantErr: false,
        },
        {
            name: "boolean value",
            expr: &hcl.LiteralValueExpr{
                Val: cty.BoolVal(true),
            },
            want:    true,
            wantErr: false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := extractLiteralValue(tt.expr)
            if (err != nil) != tt.wantErr {
                t.Errorf("extractLiteralValue() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("extractLiteralValue() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

### Integration Tests
```go
func TestParseMachineBlockIntegration(t *testing.T) {
    tests := []struct {
        name    string
        block   *hcl.Block
        want    *spookytypesmachines.Machine
        wantErr bool
    }{
        {
            name: "complete machine block",
            block: &hcl.Block{
                Type:   "machine",
                Labels: []string{"test-machine"},
                Body: &hcl.BodySchema{
                    Attributes: map[string]*hcl.AttributeSchema{
                        "hostname": {Required: true},
                        "port":     {Required: true},
                        "user":     {Required: true},
                        "description": {Required: false},
                    },
                },
            },
            want: &spookytypesmachines.Machine{
                Name:        "test-machine",
                Hostname:    "test.example.com",
                Port:        22,
                User:        "admin",
                Description: "Test machine",
            },
            wantErr: false,
        },
        {
            name: "missing required attributes",
            block: &hcl.Block{
                Type:   "machine",
                Labels: []string{"test-machine"},
                Body: &hcl.BodySchema{
                    Attributes: map[string]*hcl.AttributeSchema{},
                },
            },
            want:    nil,
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := parseMachineBlock(tt.block)
            if (err != nil) != tt.wantErr {
                t.Errorf("parseMachineBlock() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("parseMachineBlock() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

## Implementation Timeline

### Day 1: Extract Block Validation
- [ ] Extract `validateMachineBlock`
- [ ] Extract `extractMachineName`
- [ ] Add unit tests for validation functions
- [ ] Verify block validation works correctly

### Day 2: Extract Attribute Processing
- [ ] Extract `extractRequiredAttributes`
- [ ] Extract `extractOptionalAttributes`
- [ ] Extract `extractAttributeValue`
- [ ] Extract value extraction functions
- [ ] Add unit tests for attribute processing functions
- [ ] Verify attribute processing works correctly

### Day 3: Extract Nested Block Processing
- [ ] Extract `extractNestedBlocks`
- [ ] Extract `extractAuthenticationBlock`
- [ ] Extract `extractConnectionBlock`
- [ ] Add unit tests for nested block processing functions
- [ ] Verify nested block processing works correctly

### Day 4: Complete Refactoring
- [ ] Extract `createMachineFromData`
- [ ] Extract machine creation helper functions
- [ ] Refactor main `parseMachineBlock` function
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
- [ ] Easy to modify individual parsing components
- [ ] Clear function names and purposes
- [ ] Well-documented functions
- [ ] Consistent parsing patterns

## Risk Mitigation

### Potential Risks
1. **Parsing logic changes** - May affect machine block parsing behavior
2. **Attribute handling changes** - May affect attribute extraction
3. **Nested block changes** - May affect nested block processing

### Mitigation Strategies
1. **Comprehensive testing** - Test all block types and scenarios
2. **Parsing validation** - Ensure parsing behavior remains consistent
3. **Attribute validation** - Ensure attribute extraction remains consistent
4. **Gradual migration** - Implement changes incrementally with validation

## Verification

### Complexity Verification
```bash
# Verify complexity reduction
gocyclo internal/machines/loader.go | grep parseMachineBlock
```

### Functionality Verification
```bash
# Test machine block parsing
go test ./internal/machines -run TestParseMachineBlock
```

### Performance Verification
- [ ] No performance regression
- [ ] Memory usage remains stable
- [ ] Block parsing performance remains acceptable

## Conclusion

This refactoring will reduce the cyclomatic complexity of `parseMachineBlock` from 9 to approximately 3, while improving code maintainability, testability, and readability. The extracted functions will follow single responsibility principles and be easily testable in isolation.

**Expected Outcome:** A clean, maintainable function with clear separation of concerns and comprehensive test coverage for machine block parsing operations.
