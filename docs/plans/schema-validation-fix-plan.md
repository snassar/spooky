# Schema Validation Fix Plan

## Overview

The schema validation system in the spooky codebase has been **significantly improved** but still has some remaining issues that need to be addressed. This document outlines the current state and provides a plan to complete the implementation.

## Current State Analysis

### ✅ What's Working Well

1. **Schema Parser Implementation**: A comprehensive `SchemaParser` has been implemented in `internal/schemas/parser.go` that:
   - Parses HCL schema content into structured validation rules
   - Populates `schema.Validation.Fields` correctly
   - Handles nested structures and complex validation rules
   - Supports various constraint types (min/max length, patterns, enums, etc.)

2. **Schema Registry**: A robust `Registry` system exists in `internal/schemas/registry.go` that:
   - Manages schema registration and discovery
   - Supports schema versioning and compatibility
   - Provides search and statistics functionality

3. **Enhanced Validator**: A comprehensive `EnhancedValidator` in `internal/schemas/enhanced_validator.go` that:
   - Validates field constraints properly
   - Supports cross-field validation
   - Handles custom validation rules
   - Provides detailed error reporting

4. **Schema-Driven Validator**: A `SchemaDrivenValidator` in `internal/schemas/schema_driven_validator.go` that:
   - Validates configuration files against schemas
   - Supports project structure validation
   - Handles embedded schemas

5. **Schema Files**: Comprehensive schema files exist in `internal/schemas/schemas/` for:
   - Project configuration (`project.schema.hcl`)
   - Machine inventory (`machines.schema.hcl`)
   - Actions (`actions.schema.hcl`)
   - Variables (`variables-*.schema.hcl`)
   - Templates (`template-*.schema.hcl`)
   - Facts (`facts.schema.hcl`)
   - Logging (`logging.schema.hcl`)
   - Global configuration (`spooky.schema.hcl`)

### ⚠️ Current Issues

1. **Schema Loading Integration**: While the parser works, schema loading in `manager.go` calls the parser but there may be integration issues with the validation flow.

2. **Embedded Schema Loading**: The `SchemaDrivenValidator` has code to load embedded schemas but it may not be fully integrated with the main validation flow.

3. **Validation Result Consistency**: Some validation results may still show inconsistent statistics (e.g., `Valid: true` but `Errors: 1`).

4. **Schema-Specific Validation**: The schema-specific validation methods in `SchemaDrivenValidator` are mostly stubs that need implementation.

5. **Testing Coverage**: While some tests exist, comprehensive testing of the complete validation workflow is needed.

## Remaining Work

### Phase 1: Complete Schema Integration (High Priority)

#### 1.1 Fix Schema Loading Integration
**File**: `internal/schemas/manager.go`

The schema loading process calls the parser correctly, but we need to ensure the validation flow uses the parsed schemas properly:

```go
// Current implementation is mostly correct, but verify integration
func (m *Manager) Load(filePath string) (*spookytypesschemas.Schema, error) {
    // ... existing code ...
    
    // Parse validation rules from HCL content
    schemaParser := NewSchemaParser(m.logger)
    if err := schemaParser.ParseValidationRules(schema); err != nil {
        return nil, fmt.Errorf("failed to parse validation rules: %w", err)
    }
    
    // Ensure validation structure is properly initialized
    if schema.Validation == nil {
        schema.Validation = &spookytypesschemas.SchemaValidation{
            Enabled: true,
            Mode:    "strict",
            Fields:  make(map[string]*spookytypesschemas.FieldValidation),
        }
    }
    
    return schema, nil
}
```

#### 1.2 Complete Embedded Schema Integration
**File**: `internal/schemas/schema_driven_validator.go`

The embedded schema loading needs to be properly integrated:

```go
// Ensure embedded schemas are loaded at initialization
func NewSchemaDrivenValidator(logger spookytypeslogging.Logger, config *SchemaDrivenValidationConfig) *SchemaDrivenValidator {
    validator := &SchemaDrivenValidator{
        logger:          logger,
        registry:        nil, // Will be set later
        parser:          hclparse.NewParser(),
        embeddedSchemas: make(map[string]*spookytypesschemas.Schema),
        config:          config,
    }
    
    // Load embedded schemas immediately
    validator.loadEmbeddedSchemas()
    
    return validator
}
```

### Phase 2: Implement Schema-Specific Validation (High Priority)

#### 2.1 Complete Schema-Specific Validation Methods
**File**: `internal/schemas/schema_driven_validator.go`

The schema-specific validation methods are currently stubs. Implement them:

```go
func (v *SchemaDrivenValidator) validateProjectSchema(data interface{}, result *spookytypesschemas.ValidationResult) error {
    // Get project schema
    schema, err := v.getEmbeddedSchema("project")
    if err != nil {
        return fmt.Errorf("failed to get project schema: %w", err)
    }
    
    // Use enhanced validator to validate against schema
    enhancedValidator := NewEnhancedValidator(&ValidationConfig{
        Mode: ValidationModeStrict,
        ErrorHandling: &ErrorHandlingConfig{
            StopOnFirstError:   false,
            MaxErrors:          100,
            IncludeWarnings:    true,
            IncludeContext:     true,
            IncludeSuggestions: true,
        },
    })
    
    validationResult, err := enhancedValidator.ValidateWithEnhancedFeatures(context.Background(), schema, data)
    if err != nil {
        return fmt.Errorf("failed to validate project schema: %w", err)
    }
    
    // Merge validation results
    result.Errors = append(result.Errors, validationResult.Errors...)
    result.Warnings = append(result.Warnings, validationResult.Warnings...)
    result.Info = append(result.Info, validationResult.Info...)
    
    if !validationResult.Valid {
        result.Valid = false
    }
    
    return nil
}
```

#### 2.2 Implement Similar Methods for Other Schema Types
- `validateMachinesSchema`
- `validateActionsSchema`
- `validateVariablesSchema`
- `validateTemplatesSchema`
- `validateLoggingSchema`
- `validateSSHSchema`

### Phase 3: Fix Validation Result Consistency (Medium Priority)

#### 3.1 Ensure Consistent Validation Statistics
**File**: `internal/schemas/enhanced_validator.go`

Fix the validation statistics to be consistent:

```go
// Update validation statistics properly
func (v *EnhancedValidator) ValidateWithEnhancedFeatures(ctx context.Context, schema *spookytypesschemas.Schema, data interface{}) (*spookytypesschemas.ValidationResult, error) {
    // ... existing code ...
    
    // Update statistics correctly
    result.Statistics.Duration = time.Since(start)
    result.Statistics.ValidFields = result.Statistics.TotalFields - result.Statistics.InvalidFields
    
    // Ensure valid flag is set correctly
    result.Valid = len(result.Errors) == 0
    
    return result, nil
}
```

### Phase 4: Add Comprehensive Testing (Medium Priority)

#### 4.1 Add Integration Tests
**File**: `internal/schemas/integration_test.go`

Create comprehensive integration tests:

```go
func TestEndToEndSchemaValidation(t *testing.T) {
    // Test complete validation workflow
    logger := &MockLogger{}
    manager := NewManager(logger)
    
    // Load a real schema file
    schema, err := manager.Load("internal/schemas/schemas/project.schema.hcl")
    require.NoError(t, err)
    
    // Validate test data against schema
    testData := map[string]interface{}{
        "project": map[string]interface{}{
            "name": "test-project",
            "description": "A test project",
        },
    }
    
    validator := NewEnhancedValidator(&ValidationConfig{
        Mode: ValidationModeStrict,
    })
    
    result, err := validator.ValidateWithEnhancedFeatures(context.Background(), schema, testData)
    require.NoError(t, err)
    
    assert.True(t, result.Valid)
    assert.Equal(t, 0, len(result.Errors))
    assert.Greater(t, result.Statistics.TotalFields, 0)
    assert.Greater(t, result.Statistics.RulesProcessed, 0)
}
```

#### 4.2 Add Schema-Specific Tests
Test each schema type with appropriate test data:

```go
func TestProjectSchemaValidation(t *testing.T) {
    // Test project schema validation
}

func TestMachinesSchemaValidation(t *testing.T) {
    // Test machines schema validation
}

func TestActionsSchemaValidation(t *testing.T) {
    // Test actions schema validation
}
```

### Phase 5: Performance Optimization (Low Priority)

#### 5.1 Add Schema Caching
**File**: `internal/schemas/manager.go`

Implement schema caching to avoid repeated parsing:

```go
type SchemaCache struct {
    schemas map[string]*spookytypesschemas.Schema
    mutex   sync.RWMutex
    ttl     time.Duration
}

func (m *Manager) LoadWithCache(filePath string) (*spookytypesschemas.Schema, error) {
    // Check cache first
    if cached := m.cache.Get(filePath); cached != nil {
        return cached, nil
    }
    
    // Load and cache
    schema, err := m.Load(filePath)
    if err != nil {
        return nil, err
    }
    
    m.cache.Set(filePath, schema)
    return schema, nil
}
```

## Implementation Steps

### Step 1: Complete Schema Integration (Week 1)
1. Verify schema loading integration works correctly
2. Complete embedded schema integration
3. Test schema parsing with real schema files
4. Fix any integration issues

### Step 2: Implement Schema-Specific Validation (Week 1)
1. Implement all schema-specific validation methods
2. Test each schema type with appropriate test data
3. Ensure validation results are consistent
4. Add error handling for schema loading failures

### Step 3: Add Comprehensive Testing (Week 2)
1. Create integration tests for complete validation workflow
2. Add schema-specific tests for each schema type
3. Test error scenarios and edge cases
4. Ensure test coverage meets requirements

### Step 4: Performance Optimization (Week 2)
1. Implement schema caching
2. Optimize validation performance
3. Add performance benchmarks
4. Monitor memory usage

## Success Criteria

### Functional Requirements
- [x] Schemas are properly parsed into validation structures
- [x] Validation rules are correctly applied to HCL data
- [ ] All schema-specific validation methods are implemented
- [ ] Validation results are consistent and accurate
- [ ] Embedded schemas are properly integrated
- [ ] All existing schema files work correctly

### Performance Requirements
- [ ] Schema loading time < 100ms per schema
- [ ] Validation time < 50ms for typical HCL files
- [ ] Memory usage remains reasonable
- [ ] Schema caching reduces repeated parsing overhead

### Quality Requirements
- [ ] Test coverage > 90% for validation code
- [ ] All validation tests pass
- [ ] No regression in existing functionality
- [ ] Documentation is updated

## Risk Assessment

### Low Risk
- **Breaking existing functionality**: The core validation system is working well
- **Performance impact**: Schema parsing is already optimized

### Medium Risk
- **Schema integration issues**: May need debugging of integration points
- **Test coverage gaps**: Need to ensure comprehensive testing

### Mitigation Strategies
1. **Incremental testing**: Test each component thoroughly before integration
2. **Backward compatibility**: Maintain support for existing validation methods
3. **Comprehensive testing**: Add extensive tests to catch regressions
4. **Performance monitoring**: Monitor performance impact and optimize as needed

## Dependencies

### Internal Dependencies
- `internal/schemas/` package (mostly complete)
- `internal/types/schemas/` package
- HCL parsing libraries (working)

### External Dependencies
- `github.com/hashicorp/hcl/v2` for HCL parsing
- `github.com/zclconf/go-cty` for type handling

## Timeline

- **Week 1**: Complete schema integration and implement schema-specific validation
- **Week 2**: Add comprehensive testing and performance optimization

## Conclusion

The schema validation system has made significant progress and is mostly functional. The remaining work focuses on completing the integration between components and ensuring comprehensive testing coverage. The core architecture is solid and the implementation is well-structured, making the remaining work straightforward to complete.

This fix will ensure the schema validation system is fully functional and reliable for validating configuration files throughout the spooky codebase.
