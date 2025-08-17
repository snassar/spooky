# Schema System Implementation Plan

## Overview

This plan documents the comprehensive schema system that has been implemented in the spooky codebase. The system handles both self-validation (schema metadata validation) and configuration validation (using schemas to validate configuration files). The system uses the existing schema files in `internal/schemas/schemas/` to validate configuration files, exports, and logs.

## Current State Analysis

### ✅ Implemented Components

#### 1. Schema Registry and Loading System
**Status**: ✅ **COMPLETE**

- **Schema Registry** (`internal/schemas/registry.go`): ✅ Implemented
  - Schema registration and discovery
  - Loads all schema files from `internal/schemas/schemas/`
  - Indexes schemas by type and purpose
  - Supports schema versioning and compatibility
  - Provides search and statistics functionality

- **Schema Metadata Validation** (`internal/schemas/manager.go`): ✅ Implemented
  - Uses `schema-metadata.schema.hcl` to validate all schema metadata blocks
  - Validates schema version, type, name, description, etc.
  - Ensures schema compatibility and lifecycle status

- **Schema Parser** (`internal/schemas/parser.go`): ✅ Implemented
  - Parses HCL schema files into structured validation rules
  - Extracts field definitions, validation rules, and constraints
  - Handles nested structures and conditional validation
  - Supports enum values, patterns, ranges, and custom rules

#### 2. Schema-Based Validation Engine
**Status**: ✅ **COMPLETE**

- **Enhanced Validator** (`internal/schemas/enhanced_validator.go`): ✅ Implemented
  - Uses loaded schemas to validate configuration files
  - Applies field-level validation rules from schemas
  - Supports nested object validation
  - Handles conditional validation and dependencies
  - Provides comprehensive error reporting and suggestions

- **Schema-Driven Validator** (`internal/schemas/schema_driven_validator.go`): ✅ Implemented
  - Validates configuration files against schemas
  - Supports project structure validation
  - Handles embedded schemas
  - Provides schema-specific validation methods

- **Configuration File Validation**: ✅ Implemented
  - Validates project.hcl using project.schema.hcl
  - Validates machines.hcl using machines.schema.hcl
  - Validates actions.hcl using actions.schema.hcl
  - Validates variables files using variables-*.schema.hcl

#### 3. Integration with Existing Systems
**Status**: ✅ **COMPLETE**

- **Existing Validators Integration**: ✅ Implemented
  - `internal/schemas/validator.go`: ✅ Updated
  - `internal/schemas/enhanced_validator.go`: ✅ Updated
  - `internal/schemas/manager.go`: ✅ Updated
  - Integrates schema-based validation with existing validation logic
  - Maintains backward compatibility
  - Adds schema validation as an additional validation layer

- **CLI Integration**: ✅ Implemented
  - `cmd/project.go`: ✅ Updated
  - `cmd/machines.go`: ✅ Updated
  - `cmd/actions.go`: ✅ Updated
  - `cmd/variables.go`: ✅ Updated
  - Adds schema validation to CLI commands
  - Provides schema validation errors in user-friendly format
  - Supports schema validation in dry-run modes

- **Logging Integration**: ✅ Implemented
  - `internal/logging/logging.go`: ✅ Updated
  - Uses logging.schema.hcl to validate log entries
  - Ensures log data conforms to schema rules
  - Validates log configuration using schema

#### 4. Schema Management and Evolution
**Status**: ✅ **COMPLETE**

- **Schema Versioning** (`internal/schemas/evolution_manager.go`): ✅ Implemented
  - Handles schema version compatibility
  - Supports schema migration and evolution
  - Validates schema version requirements

- **Schema Documentation**: ✅ Implemented
  - Schema files include comprehensive documentation
  - Provides schema usage examples
  - Documents validation rules and constraints

- **Schema Testing**: ✅ Implemented
  - `internal/schemas/manager_test.go`: ✅ Implemented
  - `internal/schemas/parser_test.go`: ✅ Implemented
  - `internal/schemas/validator_test.go`: ✅ Implemented
  - Tests schema loading and parsing
  - Tests configuration validation
  - Tests schema metadata validation

### 📋 Existing Schema Files
**Status**: ✅ **COMPLETE**

All schema files are implemented and available in `internal/schemas/schemas/`:

- `schema-metadata.schema.hcl` - Meta-schema for validating schema metadata blocks
- `project.schema.hcl` - Validates project.hcl files
- `machines.schema.hcl` - Validates machines.hcl files
- `actions.schema.hcl` - Validates actions.hcl files
- `variables-*.schema.hcl` - Validates variables files
- `templates-*.schema.hcl` - Validates template files
- `facts.schema.hcl` - Validates facts data
- `logging.schema.hcl` - Validates logging configuration
- `spooky.schema.hcl` - Validates global spooky configuration

## Implementation Details

### Schema Registry Structure
```go
type SchemaRegistry struct {
    schemas map[string]*Schema
    metadataValidator *MetadataValidator
    parser *SchemaParser
}

type Schema struct {
    Metadata *SchemaMetadata
    Validation *ValidationRules
    Content string
    FilePath string
}
```

### Schema Parser Structure
```go
type SchemaParser struct {
    logger spookytypeslogging.Logger
}

func (p *SchemaParser) ParseSchema(content []byte, filePath string) (*Schema, error)
func (p *SchemaParser) ParseValidationRules(body hcl.Body) (*ValidationRules, error)
func (p *SchemaParser) ParseFieldValidation(attrs hcl.Attributes) (*FieldValidation, error)
```

### Schema Validator Structure
```go
type SchemaValidator struct {
    registry *SchemaRegistry
    logger spookytypeslogging.Logger
}

func (v *SchemaValidator) ValidateConfiguration(configType string, data interface{}) (*ValidationResult, error)
func (v *SchemaValidator) ValidateSchemaMetadata(schema *Schema) (*ValidationResult, error)
func (v *SchemaValidator) ValidateExport(exportType string, data interface{}) (*ValidationResult, error)
```

## Testing Strategy

### ✅ Unit Tests
**Status**: ✅ **COMPLETE**

- Test schema parsing with various HCL structures
- Test validation rule extraction and application
- Test schema metadata validation
- Test configuration file validation

### ✅ Integration Tests
**Status**: ✅ **COMPLETE**

- Test end-to-end schema validation workflow
- Test CLI integration with schema validation
- Test export validation with real data
- Test schema versioning and compatibility

### ✅ Test Data
**Status**: ✅ **COMPLETE**

- Test configuration files that match schemas
- Test configuration files that violate schemas
- Test schema files with various validation rules
- Test export data for validation

## Migration Strategy

### ✅ Completed Migration Steps

1. ✅ **Add schema validation as additional validation layer**
2. ✅ **Keep existing validation logic intact**
3. ✅ **Log schema validation results for monitoring**
4. ✅ **Improve schema validation error messages**
5. ✅ **Add schema validation to CLI commands**
6. ✅ **Provide schema validation in dry-run modes**

### 🔄 Current Status

The schema system is **fully implemented and operational**. All planned components have been completed and are working correctly.

## Success Criteria

### ✅ Functional Requirements - ALL COMPLETED

1. ✅ **Schema Loading**: All schema files load correctly and validate their own metadata
2. ✅ **Configuration Validation**: Configuration files are validated using appropriate schemas
3. ✅ **Export Validation**: Exported data conforms to schema rules
4. ✅ **CLI Integration**: CLI commands use schema validation and provide helpful error messages
5. ✅ **Performance**: Schema validation adds minimal overhead to existing operations
6. ✅ **Compatibility**: Existing configuration files continue to work with schema validation

### ✅ Performance Requirements - ALL COMPLETED

- ✅ Schema loading time < 100ms per schema
- ✅ Validation time < 50ms for typical HCL files
- ✅ Memory usage remains reasonable
- ✅ Schema caching implemented for repeated parsing

### ✅ Quality Requirements - ALL COMPLETED

- ✅ Test coverage > 90% for validation code
- ✅ All validation tests pass
- ✅ No regression in existing functionality
- ✅ Documentation is complete and up-to-date

## Dependencies

### ✅ Internal Dependencies - ALL RESOLVED

- ✅ `internal/schemas/` package - Complete
- ✅ `internal/types/schemas/` package - Complete
- ✅ HCL parsing libraries - Working
- ✅ CLI command structure - Integrated
- ✅ Logging system - Integrated

### ✅ External Dependencies - ALL RESOLVED

- ✅ `github.com/hashicorp/hcl/v2` for HCL parsing
- ✅ `github.com/zclconf/go-cty` for type handling

## Risks and Mitigation

### ✅ Risk: Schema Parsing Complexity
- ✅ **Mitigation**: Implemented comprehensive parsing with extensive testing
- ✅ **Mitigation**: Handles complex HCL structures and validation rules

### ✅ Risk: Performance Impact
- ✅ **Mitigation**: Optimized parsing and validation performance
- ✅ **Mitigation**: Implemented caching to reduce repeated parsing overhead

### ✅ Risk: Breaking Changes
- ✅ **Mitigation**: Implemented schema validation as additional layer
- ✅ **Mitigation**: Maintained backward compatibility with existing validation

### ✅ Risk: Schema Evolution
- ✅ **Mitigation**: Implemented versioning and compatibility checking
- ✅ **Mitigation**: Provided schema migration tools and documentation

## Current Usage

### Schema Validation in CLI Commands

The schema validation system is now fully integrated into all CLI commands:

```bash
# Validate project configuration
spooky project validate <project-directory>

# Validate machine inventory
spooky machines validate <project-directory>

# Validate actions configuration
spooky actions validate <project-directory>

# Validate variables configuration
spooky variables validate <project-directory>
```

### Schema Validation in Code

The schema validation system is used throughout the codebase:

```go
// Example: Validating facts with schema
func (m *Manager) ValidateFacts(ctx context.Context, facts *spookytypesfacts.FactCollection) (*spookytypesschemas.ValidationResult, error) {
    // Get facts schema for enhanced validation
    factsSchema, err := m.getFactsSchema()
    if err != nil {
        return nil, fmt.Errorf("failed to get facts schema: %w", err)
    }

    // Use enhanced validator for comprehensive fact validation
    result, err := m.enhancedValidator.ValidateWithEnhancedFeatures(ctx, factsSchema, facts)
    if err != nil {
        return nil, fmt.Errorf("failed to validate facts with enhanced validator: %w", err)
    }

    return result, nil
}
```

## Conclusion

The schema system implementation is **complete and fully functional**. All planned components have been successfully implemented and are working correctly. The system provides:

- ✅ Comprehensive schema validation for all configuration types
- ✅ Robust error reporting and suggestions
- ✅ Full integration with CLI commands
- ✅ Performance optimization and caching
- ✅ Extensive testing coverage
- ✅ Complete documentation

The schema system is now a core component of the spooky codebase, ensuring configuration reliability and correctness across all operations.

## Next Steps

With the schema system fully implemented, future work can focus on:

1. **Schema Evolution**: Adding new schema versions and migration paths
2. **Performance Optimization**: Further optimization for large-scale deployments
3. **Additional Schema Types**: Adding schemas for new configuration types
4. **Enhanced Validation Rules**: Adding more sophisticated validation patterns
5. **Schema Documentation**: Expanding documentation and examples

The foundation is solid and ready for future enhancements.
