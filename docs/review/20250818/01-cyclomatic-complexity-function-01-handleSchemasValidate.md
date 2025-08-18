# Function Improvement Plan: handleSchemasValidate

**Function:** `handleSchemasValidate`  
**File:** `cmd/schemas.go:58`  
**Current Complexity:** 15  
**Target Complexity:** < 10  
**Priority:** Critical

## Current Function Analysis

### Function Signature
```go
func handleSchemasValidate(schemaFile, dataFile string) error
```

### Current Issues
1. **Complex validation result handling** - Multiple nested conditions for success/failure reporting
2. **Mixed responsibilities** - Validation, logging, output formatting, and error reporting
3. **Deep nesting** - Multiple levels of conditional logic for different result types
4. **Repetitive error reporting** - Similar patterns repeated for errors, warnings, and recommendations
5. **Complex iteration logic** - Multiple loops with nested formatting logic

### Complexity Breakdown
- 15 cyclomatic complexity points from:
  - 1 main function entry
  - 3 validation result checks (Valid, Errors, Warnings, Recommendations)
  - 4 nested conditional blocks for different result types
  - 3 iteration loops with nested formatting
  - 4 error handling paths

## Refactoring Strategy

### Phase 1: Extract Result Reporting (Immediate - 1 day)

#### Extract Success Result Reporting
```go
func reportSuccessResults(result *spookytypesschemas.ValidationResult) {
    fmt.Printf("✅ Validation passed\n")
    fmt.Printf("📊 Statistics:\n")
    fmt.Printf("   - Total fields: %d\n", result.Statistics.TotalFields)
    fmt.Printf("   - Valid fields: %d\n", result.Statistics.ValidFields)
    fmt.Printf("   - Rules processed: %d\n", result.Statistics.RulesProcessed)
    fmt.Printf("   - Duration: %v\n", result.Statistics.Duration)
}
```

#### Extract Failure Result Reporting
```go
func reportFailureResults(result *spookytypesschemas.ValidationResult) {
    fmt.Printf("❌ Validation failed\n")
    fmt.Printf("📊 Statistics:\n")
    fmt.Printf("   - Total fields: %d\n", result.Statistics.TotalFields)
    fmt.Printf("   - Invalid fields: %d\n", result.Statistics.InvalidFields)
    fmt.Printf("   - Rules failed: %d\n", result.Statistics.RulesFailed)
    fmt.Printf("   - Duration: %v\n", result.Statistics.Duration)
}
```

#### Extract Error Reporting
```go
func reportValidationErrors(errors []spookytypesschemas.ValidationError) {
    if len(errors) == 0 {
        return
    }
    
    fmt.Printf("\n❌ Errors:\n")
    for i := range errors {
        err := &errors[i]
        reportSingleError(i+1, err)
    }
}

func reportSingleError(index int, err *spookytypesschemas.ValidationError) {
    fmt.Printf("   %d. %s\n", index, err.Message)
    
    if err.FieldPath != "" {
        fmt.Printf("      Field: %s\n", err.FieldPath)
    }
    
    if err.Location != nil {
        fmt.Printf("      Location: %s:%d:%d\n",
            err.Location.FilePath, err.Location.Line, err.Location.Column)
    }
    
    if len(err.Suggestions) > 0 {
        fmt.Printf("      Suggestions:\n")
        for _, suggestion := range err.Suggestions {
            fmt.Printf("        - %s\n", suggestion)
        }
    }
}
```

#### Extract Warning Reporting
```go
func reportValidationWarnings(warnings []spookytypesschemas.ValidationWarning) {
    if len(warnings) == 0 {
        return
    }
    
    fmt.Printf("⚠️  Warnings:\n")
    for i := range warnings {
        warning := &warnings[i]
        fmt.Printf("   %d. %s\n", i+1, warning.Message)
    }
}
```

#### Extract Recommendation Reporting
```go
func reportValidationRecommendations(recommendations []string) {
    if len(recommendations) == 0 {
        return
    }
    
    fmt.Printf("💡 Recommendations:\n")
    for _, recommendation := range recommendations {
        fmt.Printf("   - %s\n", recommendation)
    }
}
```

### Phase 2: Extract Initialization Logic (Day 2)

#### Extract Schema Manager Initialization
```go
func initializeSchemaValidation(schemaFile string) (*spookyschemas.Manager, spookytypeslogging.Logger, error) {
    logManager := spookylogging.NewLogManager()
    logger := logManager.GetLogger("schemas")
    manager := spookyschemas.NewManager(logger)

    schema, err := manager.Load(schemaFile)
    if err != nil {
        return nil, nil, fmt.Errorf("failed to load schema: %w", err)
    }

    if err := manager.TrackSchemaEvolution(schema); err != nil {
        logger.Warn("Failed to track schema evolution", map[string]interface{}{
            "error": err.Error(),
        })
    }

    return manager, logger, nil
}
```

#### Extract Validation Execution
```go
func performSchemaValidation(ctx context.Context, manager *spookyschemas.Manager, schemaFile, dataFile string) (*spookytypesschemas.ValidationResult, error) {
    schema, err := manager.Load(schemaFile)
    if err != nil {
        return nil, fmt.Errorf("failed to load schema: %w", err)
    }

    return manager.ValidateWithEnhancedFeatures(ctx, schema, dataFile)
}
```

### Phase 3: Refactored Main Function (Day 3)

#### Final Refactored Function
```go
func handleSchemasValidate(schemaFile, dataFile string) error {
    ctx := context.Background()

    // Extract initialization
    manager, logger, err := initializeSchemaValidation(schemaFile)
    if err != nil {
        return err
    }

    // Extract validation
    result, err := performSchemaValidation(ctx, manager, schemaFile, dataFile)
    if err != nil {
        return err
    }

    // Extract result reporting
    reportValidationResults(dataFile, schemaFile, result)

    return nil
}

func reportValidationResults(dataFile, schemaFile string, result *spookytypesschemas.ValidationResult) {
    fmt.Printf("🔍 Validating %s against schema %s\n", dataFile, schemaFile)

    if result.Valid {
        reportSuccessResults(result)
    } else {
        reportFailureResults(result)
        reportValidationErrors(result.Errors)
        reportValidationWarnings(result.Warnings)
    }
    
    reportValidationRecommendations(result.Recommendations)
}
```

## Complexity Reduction Metrics

### Before Refactoring
- **Cyclomatic Complexity:** 15
- **Lines of Code:** ~80
- **Nesting Levels:** 4-5
- **Responsibilities:** 4 (validation, logging, formatting, error handling)

### After Refactoring
- **Cyclomatic Complexity:** 3 (main function) + 1-2 per helper function
- **Lines of Code:** ~120 (distributed across multiple focused functions)
- **Nesting Levels:** 1-2
- **Responsibilities:** 1 per function (single responsibility principle)

## Testing Strategy

### Unit Tests for Extracted Functions
```go
func TestReportSuccessResults(t *testing.T) {
    result := &spookytypesschemas.ValidationResult{
        Valid: true,
        Statistics: &spookytypesschemas.ValidationStatistics{
            TotalFields:    10,
            ValidFields:    10,
            RulesProcessed: 5,
            Duration:       time.Second,
        },
    }
    
    // Capture output and verify formatting
    // Test implementation here
}

func TestReportSingleError(t *testing.T) {
    err := &spookytypesschemas.ValidationError{
        Message:  "Test error",
        FieldPath: "test.field",
        Location: &spookytypesschemas.Location{
            FilePath: "test.hcl",
            Line:     10,
            Column:   5,
        },
        Suggestions: []string{"Fix this", "Try that"},
    }
    
    // Capture output and verify formatting
    // Test implementation here
}
```

### Integration Tests
```go
func TestHandleSchemasValidateIntegration(t *testing.T) {
    // Test with valid schema and data
    // Test with invalid schema
    // Test with invalid data
    // Test with mixed results (errors + warnings)
}
```

## Implementation Timeline

### Day 1: Extract Result Reporting Functions
- [ ] Extract `reportSuccessResults`
- [ ] Extract `reportFailureResults`
- [ ] Extract `reportValidationErrors`
- [ ] Extract `reportSingleError`
- [ ] Extract `reportValidationWarnings`
- [ ] Extract `reportValidationRecommendations`
- [ ] Add unit tests for extracted functions

### Day 2: Extract Initialization and Validation Logic
- [ ] Extract `initializeSchemaValidation`
- [ ] Extract `performSchemaValidation`
- [ ] Add unit tests for extracted functions
- [ ] Update main function to use extracted functions

### Day 3: Complete Refactoring
- [ ] Refactor main `handleSchemasValidate` function
- [ ] Add `reportValidationResults` wrapper
- [ ] Add integration tests
- [ ] Verify complexity reduction with gocyclo
- [ ] Code review and documentation

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
- [ ] Easy to modify individual reporting components
- [ ] Clear function names and purposes
- [ ] Well-documented functions
- [ ] Consistent error handling patterns

## Risk Mitigation

### Potential Risks
1. **Function signature changes** - May affect calling code
2. **Output format changes** - May affect user expectations
3. **Error handling changes** - May affect error propagation

### Mitigation Strategies
1. **Maintain function signature** - Keep same input/output interface
2. **Preserve output format** - Ensure identical user-facing output
3. **Comprehensive testing** - Test all scenarios to prevent regressions
4. **Gradual migration** - Implement changes incrementally

## Verification

### Complexity Verification
```bash
# Verify complexity reduction
gocyclo cmd/schemas.go | grep handleSchemasValidate
```

### Functionality Verification
```bash
# Test with various scenarios
spooky schemas validate test-schema.hcl test-data.hcl
spooky schemas validate invalid-schema.hcl test-data.hcl
spooky schemas validate test-schema.hcl invalid-data.hcl
```

### Performance Verification
- [ ] No performance regression
- [ ] Memory usage remains stable
- [ ] Response time remains acceptable

## Conclusion

This refactoring will reduce the cyclomatic complexity of `handleSchemasValidate` from 15 to approximately 3, while improving code maintainability, testability, and readability. The extracted functions will follow single responsibility principles and be easily testable in isolation.

**Expected Outcome:** A clean, maintainable function with clear separation of concerns and comprehensive test coverage.
