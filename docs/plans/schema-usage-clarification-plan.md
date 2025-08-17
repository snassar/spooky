# Schema Usage Clarification Plan

## Overview

The current schema system has become confusing due to mixed purposes, inconsistent naming, and unclear loading logic. This plan establishes a clear separation between schema structure definitions, validation rules, and metadata, with explicit file organization and loading patterns.

## Current Problems

### 1. Mixed Purposes
- Schema files are used for both structure definitions AND validation rules
- Code inconsistently parses validation rules from structure files
- Unclear whether a schema file defines structure or validation logic

### 2. Inconsistent File Naming
- Some files: `variables-structure.schema.hcl`
- Others: `project.schema.hcl`
- No clear pattern for what each file type contains

### 3. Ambiguous Loading Logic
- `SchemaParser.ParseValidationRules()` tries to parse validation rules from structure files
- `loadEmbeddedSchemas()` loads all `.hcl` files without distinction
- Unclear what each loaded schema is used for

### 4. Confusing Method Names
- `validateProjectSchema()` vs `validateProjectStructure()`
- Methods don't clearly indicate what they're validating

## Proposed Solution

### 1. Clear File Organization

```
internal/schemas/schemas/
├── structure/           # Schema structure definitions
│   ├── project.hcl
│   ├── machines.hcl
│   ├── actions.hcl
│   ├── variables.hcl
│   ├── templates.hcl
│   ├── logging.hcl
│   ├── facts.hcl
│   └── spooky.hcl
├── validation/          # Custom validation rules (if needed)
│   ├── project-rules.hcl
│   ├── machines-rules.hcl
│   ├── actions-rules.hcl
│   └── ...
└── metadata/           # Schema metadata and versioning
    ├── schema-metadata.hcl
    └── compatibility.hcl
```

### 2. Clear Purpose Separation

#### Structure Schemas (`structure/`)
- Define allowed fields, types, and constraints
- Used for field-level validation
- Example: `project.hcl` defines that `name` is a required string

#### Validation Schemas (`validation/`)
- Define custom business logic validation rules
- Used for cross-field validation and complex rules
- Example: `project-rules.hcl` defines that project name must be unique

#### Metadata Schemas (`metadata/`)
- Define schema versioning and compatibility
- Used for schema evolution tracking
- Example: `schema-metadata.hcl` defines version compatibility matrix

### 3. Clear Loading Logic

```go
// SchemaDrivenValidator changes
type SchemaDrivenValidator struct {
    logger          spookytypeslogging.Logger
    registry        *SchemaRegistry
    parser          *hclparse.Parser
    structureSchemas map[string]*spookytypesschemas.Schema  // structure/ schemas
    validationSchemas map[string]*spookytypesschemas.Schema // validation/ schemas
    metadataSchemas map[string]*spookytypesschemas.Schema   // metadata/ schemas
    config          *SchemaDrivenValidationConfig
}

func (v *SchemaDrivenValidator) loadStructureSchemas() {
    // Load from structure/ directory
    structureDir := "schemas/structure"
    // ... load structure schemas
}

func (v *SchemaDrivenValidator) loadValidationSchemas() {
    // Load from validation/ directory (if any)
    validationDir := "schemas/validation"
    // ... load validation schemas
}

func (v *SchemaDrivenValidator) loadMetadataSchemas() {
    // Load from metadata/ directory
    metadataDir := "schemas/metadata"
    // ... load metadata schemas
}
```

### 4. Clear Method Names

```go
// Structure validation methods
func (v *SchemaDrivenValidator) validateProjectStructure(data interface{}, result *spookytypesschemas.ValidationResult) error
func (v *SchemaDrivenValidator) validateMachinesStructure(data interface{}, result *spookytypesschemas.ValidationResult) error
func (v *SchemaDrivenValidator) validateActionsStructure(data interface{}, result *spookytypesschemas.ValidationResult) error

// Custom validation methods
func (v *SchemaDrivenValidator) applyProjectValidationRules(data interface{}, result *spookytypesschemas.ValidationResult) error
func (v *SchemaDrivenValidator) applyMachinesValidationRules(data interface{}, result *spookytypesschemas.ValidationResult) error
func (v *SchemaDrivenValidator) applyActionsValidationRules(data interface{}, result *spookytypesschemas.ValidationResult) error

// Combined validation methods
func (v *SchemaDrivenValidator) validateProject(data interface{}, result *spookytypesschemas.ValidationResult) error {
    // 1. Validate structure
    if err := v.validateProjectStructure(data, result); err != nil {
        return err
    }
    
    // 2. Apply custom validation rules
    if err := v.applyProjectValidationRules(data, result); err != nil {
        return err
    }
    
    return nil
}
```

## Implementation Plan

### Phase 1: File Reorganization

#### 1.1 Create New Directory Structure
```bash
mkdir -p internal/schemas/schemas/structure
mkdir -p internal/schemas/schemas/validation
mkdir -p internal/schemas/schemas/metadata
```

#### 1.2 Move Existing Files
```bash
# Move structure files
mv internal/schemas/schemas/project.schema.hcl internal/schemas/schemas/structure/project.hcl
mv internal/schemas/schemas/machines.schema.hcl internal/schemas/schemas/structure/machines.hcl
mv internal/schemas/schemas/actions.schema.hcl internal/schemas/schemas/structure/actions.hcl
mv internal/schemas/schemas/variables-structure.schema.hcl internal/schemas/schemas/structure/variables.hcl
mv internal/schemas/schemas/template-structure.schema.hcl internal/schemas/schemas/structure/templates.hcl
mv internal/schemas/schemas/logging.schema.hcl internal/schemas/schemas/structure/logging.hcl
mv internal/schemas/schemas/facts.schema.hcl internal/schemas/schemas/structure/facts.hcl
mv internal/schemas/schemas/spooky.schema.hcl internal/schemas/schemas/structure/spooky.hcl

# Move metadata files
mv internal/schemas/schemas/schema-metadata.schema.hcl internal/schemas/schemas/metadata/schema-metadata.hcl

# Move validation files (if any)
mv internal/schemas/schemas/template-metadata-validation.schema.hcl internal/schemas/schemas/validation/template-metadata-rules.hcl

# Remove remaining old files
rm internal/schemas/schemas/*.schema.hcl
```

### Phase 2: Code Changes

#### 2.1 Update SchemaDrivenValidator

**File: `internal/schemas/schema_driven_validator.go`**

```go
type SchemaDrivenValidator struct {
    logger          spookytypeslogging.Logger
    registry        *SchemaRegistry
    parser          *hclparse.Parser
    structureSchemas map[string]*spookytypesschemas.Schema  // structure/ schemas
    validationSchemas map[string]*spookytypesschemas.Schema // validation/ schemas
    metadataSchemas map[string]*spookytypesschemas.Schema   // metadata/ schemas
    config          *SchemaDrivenValidationConfig
}

func NewSchemaDrivenValidator(logger spookytypeslogging.Logger, config *SchemaDrivenValidationConfig) *SchemaDrivenValidator {
    if config == nil {
        config = &SchemaDrivenValidationConfig{
            UseEmbeddedSchemas: true,
            StrictValidation:   true,
            AllowUnknownFields: false,
            DetailedErrors:     true,
            CustomRules:        make(map[string]CustomValidationRule),
        }
    }

    validator := &SchemaDrivenValidator{
        logger:          logger,
        registry:        nil, // Will be set later
        parser:          hclparse.NewParser(),
        structureSchemas: make(map[string]*spookytypesschemas.Schema),
        validationSchemas: make(map[string]*spookytypesschemas.Schema),
        metadataSchemas: make(map[string]*spookytypesschemas.Schema),
        config:          config,
    }

    // Load all schema types
    validator.loadStructureSchemas()
    validator.loadValidationSchemas()
    validator.loadMetadataSchemas()

    return validator
}

func (v *SchemaDrivenValidator) loadStructureSchemas() {
    structureDir := "schemas/structure"
    
    if info, err := os.Stat(structureDir); err != nil || !info.IsDir() {
        v.logger.Error("Structure schemas directory not found", fmt.Errorf("structure schemas directory not found: %s", structureDir), map[string]interface{}{
            "structure_dir": structureDir,
        })
        return
    }

    entries, err := os.ReadDir(structureDir)
    if err != nil {
        v.logger.Error("Failed to read structure schemas directory", err, map[string]interface{}{
            "structure_dir": structureDir,
        })
        return
    }

    for _, entry := range entries {
        if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".hcl") {
            schemaPath := filepath.Join(structureDir, entry.Name())
            schema, err := v.loadStructureSchema(schemaPath)
            if err != nil {
                v.logger.Error("Failed to load structure schema", err, map[string]interface{}{
                    "schema_path": schemaPath,
                })
                continue
            }

            schemaName := strings.TrimSuffix(entry.Name(), ".hcl")
            v.structureSchemas[schemaName] = schema
            
            v.logger.Info("Loaded structure schema", map[string]interface{}{
                "schema_name": schemaName,
                "schema_path": schemaPath,
            })
        }
    }
}

func (v *SchemaDrivenValidator) loadValidationSchemas() {
    validationDir := "schemas/validation"
    
    if info, err := os.Stat(validationDir); err != nil || !info.IsDir() {
        v.logger.Debug("Validation schemas directory not found", map[string]interface{}{
            "validation_dir": validationDir,
        })
        return
    }

    entries, err := os.ReadDir(validationDir)
    if err != nil {
        v.logger.Error("Failed to read validation schemas directory", err, map[string]interface{}{
            "validation_dir": validationDir,
        })
        return
    }

    for _, entry := range entries {
        if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".hcl") {
            schemaPath := filepath.Join(validationDir, entry.Name())
            schema, err := v.loadValidationSchema(schemaPath)
            if err != nil {
                v.logger.Error("Failed to load validation schema", err, map[string]interface{}{
                    "schema_path": schemaPath,
                })
                continue
            }

            schemaName := strings.TrimSuffix(entry.Name(), ".hcl")
            v.validationSchemas[schemaName] = schema
            
            v.logger.Info("Loaded validation schema", map[string]interface{}{
                "schema_name": schemaName,
                "schema_path": schemaPath,
            })
        }
    }
}

func (v *SchemaDrivenValidator) loadMetadataSchemas() {
    metadataDir := "schemas/metadata"
    
    if info, err := os.Stat(metadataDir); err != nil || !info.IsDir() {
        v.logger.Debug("Metadata schemas directory not found", map[string]interface{}{
            "metadata_dir": metadataDir,
        })
        return
    }

    entries, err := os.ReadDir(metadataDir)
    if err != nil {
        v.logger.Error("Failed to read metadata schemas directory", err, map[string]interface{}{
            "metadata_dir": metadataDir,
        })
        return
    }

    for _, entry := range entries {
        if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".hcl") {
            schemaPath := filepath.Join(metadataDir, entry.Name())
            schema, err := v.loadMetadataSchema(schemaPath)
            if err != nil {
                v.logger.Error("Failed to load metadata schema", err, map[string]interface{}{
                    "schema_path": schemaPath,
                })
                continue
            }

            schemaName := strings.TrimSuffix(entry.Name(), ".hcl")
            v.metadataSchemas[schemaName] = schema
            
            v.logger.Info("Loaded metadata schema", map[string]interface{}{
                "schema_name": schemaName,
                "schema_path": schemaPath,
            })
        }
    }
}
```

#### 2.2 Update Validation Methods

```go
// Structure validation methods
func (v *SchemaDrivenValidator) validateProjectStructure(data interface{}, result *spookytypesschemas.ValidationResult) error {
    schema, err := v.getStructureSchema("project")
    if err != nil {
        return fmt.Errorf("failed to get project structure schema: %w", err)
    }

    // Use enhanced validator to validate structure
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
        return fmt.Errorf("failed to validate project structure: %w", err)
    }

    // Merge validation results
    result.Errors = append(result.Errors, validationResult.Errors...)
    result.Warnings = append(result.Warnings, validationResult.Warnings...)
    result.Info = append(result.Info, validationResult.Info...)
    
    if !validationResult.Valid {
        result.Valid = false
    }

    // Update statistics
    result.Statistics.TotalFields += validationResult.Statistics.TotalFields
    result.Statistics.ValidFields += validationResult.Statistics.ValidFields
    result.Statistics.InvalidFields += validationResult.Statistics.InvalidFields
    result.Statistics.RulesProcessed += validationResult.Statistics.RulesProcessed
    result.Statistics.RulesFailed += validationResult.Statistics.RulesFailed

    return nil
}

// Custom validation methods
func (v *SchemaDrivenValidator) applyProjectValidationRules(data interface{}, result *spookytypesschemas.ValidationResult) error {
    schema, err := v.getValidationSchema("project-rules")
    if err != nil {
        // No custom validation rules defined, that's okay
        return nil
    }

    // Apply custom validation rules
    // This would implement business logic validation
    // For now, just log that we're applying custom rules
    v.logger.Debug("Applying custom project validation rules", map[string]interface{}{
        "schema_name": "project-rules",
    })

    return nil
}

// Combined validation methods
func (v *SchemaDrivenValidator) validateProject(data interface{}, result *spookytypesschemas.ValidationResult) error {
    // 1. Validate structure
    if err := v.validateProjectStructure(data, result); err != nil {
        return err
    }
    
    // 2. Apply custom validation rules
    if err := v.applyProjectValidationRules(data, result); err != nil {
        return err
    }
    
    return nil
}

// Helper methods
func (v *SchemaDrivenValidator) getStructureSchema(name string) (*spookytypesschemas.Schema, error) {
    if schema, exists := v.structureSchemas[name]; exists {
        return schema, nil
    }
    return nil, fmt.Errorf("structure schema not found: %s", name)
}

func (v *SchemaDrivenValidator) getValidationSchema(name string) (*spookytypesschemas.Schema, error) {
    if schema, exists := v.validationSchemas[name]; exists {
        return schema, nil
    }
    return nil, fmt.Errorf("validation schema not found: %s", name)
}

func (v *SchemaDrivenValidator) getMetadataSchema(name string) (*spookytypesschemas.Schema, error) {
    if schema, exists := v.metadataSchemas[name]; exists {
        return schema, nil
    }
    return nil, fmt.Errorf("metadata schema not found: %s", name)
}
```

#### 2.3 Update Manager

**File: `internal/schemas/manager.go`**

```go
// Update loadEmbeddedSchemas to use new structure
func (m *Manager) loadEmbeddedSchemas() {
    // Load structure schemas
    m.loadStructureSchemas()
    
    // Load validation schemas
    m.loadValidationSchemas()
    
    // Load metadata schemas
    m.loadMetadataSchemas()
}

func (m *Manager) loadStructureSchemas() {
    structureDir := "schemas/structure"
    
    if info, err := os.Stat(structureDir); err != nil || !info.IsDir() {
        m.logger.Error("Structure schemas directory not found", fmt.Errorf("structure schemas directory not found: %s", structureDir), map[string]interface{}{
            "structure_dir": structureDir,
        })
        return
    }

    entries, err := os.ReadDir(structureDir)
    if err != nil {
        m.logger.Error("Failed to read structure schemas directory", err, map[string]interface{}{
            "structure_dir": structureDir,
        })
        return
    }

    for _, entry := range entries {
        if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".hcl") {
            schemaPath := filepath.Join(structureDir, entry.Name())
            schema, err := m.Load(schemaPath)
            if err != nil {
                m.logger.Warn("Failed to load structure schema", map[string]interface{}{
                    "schema_path": schemaPath,
                    "error":       err.Error(),
                })
                continue
            }

            schemaName := strings.TrimSuffix(entry.Name(), ".hcl")
            m.registry.Register(schema)
            
            m.logger.Info("Loaded structure schema", map[string]interface{}{
                "schema_name": schemaName,
                "schema_path": schemaPath,
            })
        }
    }
}

// Similar methods for loadValidationSchemas() and loadMetadataSchemas()
```

### Phase 3: Update Tests

#### 3.1 Update Integration Tests

**File: `internal/schemas/integration_test.go`**

```go
func TestStructureSchemaValidation(t *testing.T) {
    logger := &MockLogger{}
    validator := NewSchemaDrivenValidator(logger, nil)
    
    // Test structure validation
    testData := map[string]interface{}{
        "project": map[string]interface{}{
            "name": "test-project",
            "description": "Test project",
        },
    }
    
    result := &spookytypesschemas.ValidationResult{
        Valid: true,
        Errors: []*spookytypesschemas.ValidationError{},
        Warnings: []*spookytypesschemas.ValidationError{},
        Info: []*spookytypesschemas.ValidationError{},
        Statistics: &spookytypesschemas.ValidationStatistics{},
    }
    
    err := validator.validateProjectStructure(testData, result)
    require.NoError(t, err)
    assert.True(t, result.Valid)
}

func TestCustomValidationRules(t *testing.T) {
    logger := &MockLogger{}
    validator := NewSchemaDrivenValidator(logger, nil)
    
    // Test custom validation rules
    testData := map[string]interface{}{
        "project": map[string]interface{}{
            "name": "test-project",
        },
    }
    
    result := &spookytypesschemas.ValidationResult{
        Valid: true,
        Errors: []*spookytypesschemas.ValidationError{},
        Warnings: []*spookytypesschemas.ValidationError{},
        Info: []*spookytypesschemas.ValidationError{},
        Statistics: &spookytypesschemas.ValidationStatistics{},
    }
    
    err := validator.applyProjectValidationRules(testData, result)
    require.NoError(t, err)
    // Custom rules might add warnings or errors
}

func TestCombinedValidation(t *testing.T) {
    logger := &MockLogger{}
    validator := NewSchemaDrivenValidator(logger, nil)
    
    // Test combined validation (structure + custom rules)
    testData := map[string]interface{}{
        "project": map[string]interface{}{
            "name": "test-project",
            "description": "Test project",
        },
    }
    
    result := &spookytypesschemas.ValidationResult{
        Valid: true,
        Errors: []*spookytypesschemas.ValidationError{},
        Warnings: []*spookytypesschemas.ValidationError{},
        Info: []*spookytypesschemas.ValidationError{},
        Statistics: &spookytypesschemas.ValidationStatistics{},
    }
    
    err := validator.validateProject(testData, result)
    require.NoError(t, err)
    assert.True(t, result.Valid)
}
```

### Phase 4: Documentation Updates

#### 4.1 Update Schema Documentation

**File: `docs/schemas/README.md`**

```markdown
# Schema System Documentation

## Overview

The schema system is organized into three clear categories:

### Structure Schemas (`schemas/structure/`)
Define the allowed fields, types, and constraints for configuration files.
- `project.hcl` - Project configuration structure
- `machines.hcl` - Machine inventory structure
- `actions.hcl` - Action definitions structure
- `variables.hcl` - Variable definitions structure
- `templates.hcl` - Template definitions structure
- `logging.hcl` - Logging configuration structure
- `facts.hcl` - Fact collection structure
- `spooky.hcl` - Global spooky configuration structure

### Validation Schemas (`schemas/validation/`)
Define custom business logic validation rules.
- `project-rules.hcl` - Custom project validation rules
- `machines-rules.hcl` - Custom machine validation rules
- `actions-rules.hcl` - Custom action validation rules

### Metadata Schemas (`schemas/metadata/`)
Define schema versioning and compatibility information.
- `schema-metadata.hcl` - Schema metadata and versioning
- `compatibility.hcl` - Version compatibility matrix

## Usage

### Structure Validation
```go
validator := NewSchemaDrivenValidator(logger, nil)
err := validator.validateProjectStructure(data, result)
```

### Custom Validation Rules
```go
err := validator.applyProjectValidationRules(data, result)
```

### Combined Validation
```go
err := validator.validateProject(data, result) // Structure + Custom Rules
```
```

## Benefits

### 1. Clear Separation of Concerns
- Structure validation is separate from business logic validation
- Each schema type has a clear, single purpose
- No more confusion about what each file contains

### 2. Improved Maintainability
- Easy to find and modify structure definitions
- Custom validation rules are clearly separated
- Schema metadata is organized and versioned

### 3. Better Testing
- Can test structure validation independently
- Can test custom validation rules independently
- Clear test organization matches schema organization

### 4. Enhanced Flexibility
- Can add custom validation rules without modifying structure
- Can evolve structure without affecting custom rules
- Clear implementation path for schema changes

## Implementation Strategy

### 1. Direct Implementation
- Reorganize files immediately
- Update code to use new structure
- Update tests for new organization
- Update documentation

### 2. Testing Strategy
- Comprehensive tests for each schema type
- Integration tests for combined validation
- Performance tests for large schemas

## Success Criteria

1. **Clarity**: No confusion about what each schema file contains
2. **Maintainability**: Easy to find and modify schema definitions
3. **Testability**: Clear separation allows independent testing
4. **Performance**: No performance regression from reorganization
5. **Functionality**: All existing functionality preserved with clearer structure
