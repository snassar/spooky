# Function Improvement Plan: GenerateDocumentation

**Function:** `GenerateDocumentation`  
**File:** `internal/schemas/manager.go:787`  
**Current Complexity:** 15  
**Target Complexity:** < 10  
**Priority:** Critical

## Current Function Analysis

### Function Signature
```go
func (m *Manager) GenerateDocumentation(schema *spookytypesschemas.Schema) (string, error)
```

### Current Issues
1. **Multiple nested null checks** - Deep nesting for validation, evolution, and field checks
2. **Complex string building logic** - Multiple string concatenation operations with conditional logic
3. **Mixed responsibilities** - Documentation generation, formatting, validation, and iteration
4. **Deep nesting** - Multiple levels of conditional logic for different documentation sections
5. **Repetitive documentation patterns** - Similar patterns repeated for different schema sections

### Complexity Breakdown
- 15 cyclomatic complexity points from:
  - 1 main function entry
  - 3 null checks (schema, validation, evolution)
  - 4 nested conditional blocks for different sections
  - 3 iteration loops with nested formatting
  - 2 error handling paths
  - 2 string building conditions

## Refactoring Strategy

### Phase 1: Extract Header Generation (Immediate - 1 day)

#### Extract Schema Header Generation
```go
func (m *Manager) generateSchemaHeader(schema *spookytypesschemas.Schema) string {
    var header strings.Builder
    
    header.WriteString(fmt.Sprintf("# %s\n\n", schema.Name))
    header.WriteString(fmt.Sprintf("**Type:** %s\n", schema.Type))
    header.WriteString(fmt.Sprintf("**Version:** %s\n", schema.Version))
    header.WriteString(fmt.Sprintf("**Description:** %s\n\n", schema.Description))
    
    return header.String()
}
```

#### Extract Validation Documentation
```go
func (m *Manager) generateValidationDocumentation(validation *spookytypesschemas.Validation) string {
    if validation == nil {
        return ""
    }
    
    var doc strings.Builder
    
    doc.WriteString("## Validation\n\n")
    doc.WriteString(fmt.Sprintf("**Mode:** %s\n", validation.Mode))
    doc.WriteString(fmt.Sprintf("**Enabled:** %t\n\n", validation.Enabled))
    
    // Generate field documentation
    if validation.Fields != nil {
        doc.WriteString(m.generateFieldDocumentation(validation.Fields))
    }
    
    // Generate rule documentation
    if len(validation.Rules) > 0 {
        doc.WriteString(m.generateRuleDocumentation(validation.Rules))
    }
    
    return doc.String()
}

func (m *Manager) generateFieldDocumentation(fields map[string]*spookytypesschemas.Field) string {
    var doc strings.Builder
    
    doc.WriteString("### Fields\n\n")
    for fieldName, field := range fields {
        doc.WriteString(fmt.Sprintf("#### %s\n", fieldName))
        doc.WriteString(fmt.Sprintf("**Type:** %s\n", field.Type))
        doc.WriteString(fmt.Sprintf("**Required:** %t\n", field.Required))
        
        if field.Description != "" {
            doc.WriteString(fmt.Sprintf("**Description:** %s\n", field.Description))
        }
        doc.WriteString("\n")
    }
    
    return doc.String()
}

func (m *Manager) generateRuleDocumentation(rules []spookytypesschemas.ValidationRule) string {
    var doc strings.Builder
    
    doc.WriteString("### Validation Rules\n\n")
    for i := range rules {
        rule := &rules[i]
        doc.WriteString(fmt.Sprintf("#### %s\n", rule.Name))
        doc.WriteString(fmt.Sprintf("**Type:** %s\n", rule.Type))
        doc.WriteString(fmt.Sprintf("**Severity:** %s\n", rule.Severity))
        doc.WriteString(fmt.Sprintf("**Message:** %s\n\n", rule.Message))
    }
    
    return doc.String()
}
```

### Phase 2: Extract Evolution Documentation (Day 2)

#### Extract Evolution Documentation Generation
```go
func (m *Manager) generateEvolutionDocumentation(evolution *spookytypesschemas.Evolution) string {
    if evolution == nil {
        return ""
    }
    
    var doc strings.Builder
    
    doc.WriteString("## Evolution\n\n")
    
    if len(evolution.Deprecations) > 0 {
        doc.WriteString(m.generateDeprecationDocumentation(evolution.Deprecations))
    }
    
    if len(evolution.BreakingChanges) > 0 {
        doc.WriteString(m.generateBreakingChangesDocumentation(evolution.BreakingChanges))
    }
    
    return doc.String()
}

func (m *Manager) generateDeprecationDocumentation(deprecations []spookytypesschemas.Deprecation) string {
    var doc strings.Builder
    
    doc.WriteString("### Deprecations\n\n")
    for i := range deprecations {
        deprecation := &deprecations[i]
        doc.WriteString(fmt.Sprintf("- **%s:** %s\n", deprecation.Field, deprecation.Reason))
        
        if deprecation.Replacement != "" {
            doc.WriteString(fmt.Sprintf("  - **Replacement:** %s\n", deprecation.Replacement))
        }
        doc.WriteString("\n")
    }
    
    return doc.String()
}

func (m *Manager) generateBreakingChangesDocumentation(breakingChanges []spookytypesschemas.BreakingChange) string {
    var doc strings.Builder
    
    doc.WriteString("### Breaking Changes\n\n")
    for i := range breakingChanges {
        breakingChange := &breakingChanges[i]
        doc.WriteString(fmt.Sprintf("- **%s:** %s\n", breakingChange.Field, breakingChange.Description))
        doc.WriteString(fmt.Sprintf("  - **Impact:** %s\n", breakingChange.Impact))
        
        if breakingChange.Mitigation != "" {
            doc.WriteString(fmt.Sprintf("  - **Mitigation:** %s\n", breakingChange.Mitigation))
        }
        doc.WriteString("\n")
    }
    
    return doc.String()
}
```

### Phase 3: Refactored Main Function (Day 3)

#### Final Refactored Function
```go
func (m *Manager) GenerateDocumentation(schema *spookytypesschemas.Schema) (string, error) {
    if schema == nil {
        return "", fmt.Errorf("schema cannot be nil")
    }
    
    var doc strings.Builder
    
    // Generate schema header
    doc.WriteString(m.generateSchemaHeader(schema))
    
    // Generate validation documentation
    if schema.Validation != nil {
        doc.WriteString(m.generateValidationDocumentation(schema.Validation))
    }
    
    // Generate evolution documentation
    if schema.Evolution != nil {
        doc.WriteString(m.generateEvolutionDocumentation(schema.Evolution))
    }
    
    return doc.String(), nil
}
```

## Complexity Reduction Metrics

### Before Refactoring
- **Cyclomatic Complexity:** 15
- **Lines of Code:** ~80
- **Nesting Levels:** 4-5
- **Responsibilities:** 4 (header generation, validation docs, evolution docs, formatting)

### After Refactoring
- **Cyclomatic Complexity:** 3 (main function) + 1-2 per helper function
- **Lines of Code:** ~120 (distributed across multiple focused functions)
- **Nesting Levels:** 1-2
- **Responsibilities:** 1 per function (single responsibility principle)

## Testing Strategy

### Unit Tests for Extracted Functions
```go
func TestGenerateSchemaHeader(t *testing.T) {
    schema := &spookytypesschemas.Schema{
        Name:        "TestSchema",
        Type:        "hcl",
        Version:     "1.0.0",
        Description: "Test schema for documentation",
    }
    
    header := generateSchemaHeader(schema)
    
    expected := "# TestSchema\n\n**Type:** hcl\n**Version:** 1.0.0\n**Description:** Test schema for documentation\n\n"
    if header != expected {
        t.Errorf("generateSchemaHeader() = %v, want %v", header, expected)
    }
}

func TestGenerateFieldDocumentation(t *testing.T) {
    fields := map[string]*spookytypesschemas.Field{
        "testField": {
            Type:        "string",
            Required:    true,
            Description: "Test field description",
        },
    }
    
    doc := generateFieldDocumentation(fields)
    
    if !strings.Contains(doc, "#### testField") {
        t.Error("generateFieldDocumentation() missing field name")
    }
    if !strings.Contains(doc, "**Type:** string") {
        t.Error("generateFieldDocumentation() missing field type")
    }
    if !strings.Contains(doc, "**Required:** true") {
        t.Error("generateFieldDocumentation() missing required status")
    }
    if !strings.Contains(doc, "**Description:** Test field description") {
        t.Error("generateFieldDocumentation() missing description")
    }
}

func TestGenerateDeprecationDocumentation(t *testing.T) {
    deprecations := []spookytypesschemas.Deprecation{
        {
            Field:       "oldField",
            Reason:      "Replaced by newField",
            Replacement: "newField",
        },
    }
    
    doc := generateDeprecationDocumentation(deprecations)
    
    if !strings.Contains(doc, "- **oldField:** Replaced by newField") {
        t.Error("generateDeprecationDocumentation() missing deprecation info")
    }
    if !strings.Contains(doc, "  - **Replacement:** newField") {
        t.Error("generateDeprecationDocumentation() missing replacement info")
    }
}
```

### Integration Tests
```go
func TestGenerateDocumentationIntegration(t *testing.T) {
    schema := &spookytypesschemas.Schema{
        Name:        "TestSchema",
        Type:        "hcl",
        Version:     "1.0.0",
        Description: "Test schema",
        Validation: &spookytypesschemas.Validation{
            Mode:    "strict",
            Enabled: true,
            Fields: map[string]*spookytypesschemas.Field{
                "testField": {
                    Type:     "string",
                    Required: true,
                },
            },
        },
        Evolution: &spookytypesschemas.Evolution{
            Deprecations: []spookytypesschemas.Deprecation{
                {
                    Field:       "oldField",
                    Reason:      "Replaced",
                    Replacement: "newField",
                },
            },
        },
    }
    
    doc, err := GenerateDocumentation(schema)
    if err != nil {
        t.Errorf("GenerateDocumentation() error = %v", err)
        return
    }
    
    // Verify all sections are present
    if !strings.Contains(doc, "# TestSchema") {
        t.Error("Documentation missing schema header")
    }
    if !strings.Contains(doc, "## Validation") {
        t.Error("Documentation missing validation section")
    }
    if !strings.Contains(doc, "## Evolution") {
        t.Error("Documentation missing evolution section")
    }
    if !strings.Contains(doc, "### Deprecations") {
        t.Error("Documentation missing deprecations section")
    }
}
```

## Implementation Timeline

### Day 1: Extract Header and Validation Documentation
- [ ] Extract `generateSchemaHeader`
- [ ] Extract `generateValidationDocumentation`
- [ ] Extract `generateFieldDocumentation`
- [ ] Extract `generateRuleDocumentation`
- [ ] Add unit tests for extracted functions
- [ ] Verify header and validation documentation generation

### Day 2: Extract Evolution Documentation
- [ ] Extract `generateEvolutionDocumentation`
- [ ] Extract `generateDeprecationDocumentation`
- [ ] Extract `generateBreakingChangesDocumentation`
- [ ] Add unit tests for evolution documentation functions
- [ ] Verify evolution documentation generation

### Day 3: Complete Refactoring
- [ ] Refactor main `GenerateDocumentation` function
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
- [ ] Easy to modify individual documentation components
- [ ] Clear function names and purposes
- [ ] Well-documented functions
- [ ] Consistent documentation formatting

## Risk Mitigation

### Potential Risks
1. **Documentation format changes** - May affect user expectations
2. **Performance impact** - Multiple string operations may impact performance
3. **Missing sections** - May miss important documentation sections

### Mitigation Strategies
1. **Preserve output format** - Ensure identical documentation output
2. **Performance benchmarking** - Measure performance impact and optimize if needed
3. **Comprehensive testing** - Test all documentation scenarios
4. **Gradual migration** - Implement changes incrementally with validation

## Verification

### Complexity Verification
```bash
# Verify complexity reduction
gocyclo internal/schemas/manager.go | grep GenerateDocumentation
```

### Functionality Verification
```bash
# Test documentation generation
go test ./internal/schemas -run TestGenerateDocumentation
```

### Performance Verification
- [ ] No performance regression
- [ ] Memory usage remains stable
- [ ] String building performance remains acceptable

## Conclusion

This refactoring will reduce the cyclomatic complexity of `GenerateDocumentation` from 15 to approximately 3, while improving code maintainability, testability, and readability. The extracted functions will follow single responsibility principles and be easily testable in isolation.

**Expected Outcome:** A clean, maintainable function with clear separation of concerns and comprehensive test coverage for documentation generation operations.
