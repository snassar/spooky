# Templates System Implementation Plan

## Current Status Analysis

### Existing Implementation
- **Template manager** implementing basic template operations in `internal/templates/manager.go`
- **Template integration** implementing `TemplatesIntegration` interface in `internal/templates/integration.go`
- **Template types** defined in `internal/types/templates/template.go`
- **Integration manager** that includes `GetTemplatesIntegration()` method
- **Integration factory** with `createTemplatesIntegration()` method

### Current Capabilities
- Basic template loading from files
- Simple template rendering using Go's text/template engine
- Basic template validation (content checks, syntax validation)
- Integration with the main system through `TemplatesIntegration` interface

### Missing Components
- Advanced template functions and security restrictions
- Template context management with facts, variables, and machines data
- Template metadata management and validation
- Template caching and performance optimization
- Template security sandboxing and access control
- Template composition and inheritance patterns

## Implementation Plan

### Phase 1: Core Template System Enhancement

#### 1.1 Template Manager Enhancement
**File**: `internal/templates/manager.go`

**Enhancements**:
- Add template caching with TTL support
- Implement template function registry with security restrictions
- Add template context resolution with facts, variables, and machines data
- Implement template metadata management
- Add template validation against schemas

**Key Features**:
```go
type Manager struct {
    logger spookytypeslogging.Logger
    cache  TemplateCache
    functions TemplateFunctionRegistry
    contextResolver TemplateContextResolver
    metadataManager TemplateMetadataManager
    validator TemplateValidator
}

// New methods to add:
func (m *Manager) RegisterTemplateFunctions(functions map[string]interface{}) error
func (m *Manager) ResolveTemplateContext(ctx context.Context, template *spookytypes.Template, data map[string]interface{}) (*spookytypes.TemplateContext, error)
func (m *Manager) ValidateTemplateWithSchema(ctx context.Context, template *spookytypes.Template) (*spookytypesschemas.ValidationResult, error)
func (m *Manager) GetTemplateMetadata(ctx context.Context, templatePath string) (*spookytypes.TemplateMetadata, error)
```

#### 1.2 Template Function Registry
**File**: `internal/templates/functions.go` (new)

**Implementation**:
- Implement secure template function registry
- Add built-in functions for data access (var, facts, env, etc.)
- Add string manipulation functions (upper, lower, trim, etc.)
- Add mathematical functions (add, sub, mul, div)
- Add array manipulation functions (length, index, first, last, etc.)
- Implement function security restrictions and sandboxing

**Key Features**:
```go
type TemplateFunctionRegistry struct {
    functions map[string]TemplateFunction
    security  FunctionSecurityManager
    cache     FunctionResultCache
}

type TemplateFunction struct {
    Name        string
    Function    interface{}
    Security    FunctionSecurity
    Description string
    Examples    []FunctionExample
}

type FunctionSecurity struct {
    AllowedInRestricted bool
    MaxExecutionTime    time.Duration
    MaxMemoryUsage      int64
    AllowedPatterns     []string
    ForbiddenPatterns   []string
}
```

#### 1.3 Template Context Resolver
**File**: `internal/templates/context.go` (new)

**Implementation**:
- Implement context resolution for facts, variables, machines, and environment data
- Add context caching and lazy loading
- Implement context validation and transformation
- Add context composition and inheritance patterns

**Key Features**:
```go
type TemplateContextResolver struct {
    factsManager    spookyinterfaces.FactsIntegration
    variablesManager spookyinterfaces.VariablesIntegration
    machinesManager spookyinterfaces.MachinesIntegration
    cache          ContextCache
    validator      ContextValidator
}

func (r *TemplateContextResolver) ResolveContext(ctx context.Context, template *spookytypes.Template, data map[string]interface{}) (*spookytypes.TemplateContext, error)
func (r *TemplateContextResolver) ValidateContext(ctx context.Context, context *spookytypes.TemplateContext) error
func (r *TemplateContextResolver) TransformContext(ctx context.Context, context *spookytypes.TemplateContext) error
```

#### 1.4 Template Metadata Manager
**File**: `internal/templates/metadata.go` (new)

**Implementation**:
- Implement template metadata management
- Add metadata validation against schemas
- Implement metadata inheritance and composition
- Add metadata indexing and search capabilities

**Key Features**:
```go
type TemplateMetadataManager struct {
    validator MetadataValidator
    indexer   MetadataIndexer
    cache     MetadataCache
}

func (m *TemplateMetadataManager) LoadMetadata(ctx context.Context, templatePath string) (*spookytypes.TemplateMetadata, error)
func (m *TemplateMetadataManager) ValidateMetadata(ctx context.Context, metadata *spookytypes.TemplateMetadata) error
func (m *TemplateMetadataManager) IndexMetadata(ctx context.Context, metadata *spookytypes.TemplateMetadata) error
func (m *TemplateMetadataManager) SearchMetadata(ctx context.Context, query string) ([]*spookytypes.TemplateMetadata, error)
```

### Phase 2: Schema Integration and Validation

#### 2.1 Schema-Driven Template Validation
**File**: `internal/templates/validator.go` (enhance existing)

**Enhancements**:
- Integrate with template schemas for validation
- Add validation for template structure, context, functions, and metadata
- Implement validation result aggregation and reporting
- Add validation caching and performance optimization

**Key Features**:
```go
type TemplateValidator struct {
    schemaValidator spookyschemas.SchemaValidator
    functionValidator FunctionValidator
    contextValidator ContextValidator
    metadataValidator MetadataValidator
    cache           ValidationCache
}

func (v *TemplateValidator) ValidateTemplateComprehensive(ctx context.Context, template *spookytypes.Template) (*spookytypesschemas.ValidationResult, error)
func (v *TemplateValidator) ValidateTemplateStructure(ctx context.Context, template *spookytypes.Template) error
func (v *TemplateValidator) ValidateTemplateFunctions(ctx context.Context, template *spookytypes.Template) error
func (v *TemplateValidator) ValidateTemplateContext(ctx context.Context, template *spookytypes.Template) error
func (v *TemplateValidator) ValidateTemplateMetadata(ctx context.Context, template *spookytypes.Template) error
```

#### 2.2 Template Schema Integration
**Integration Points**:
- Use `template-structure.schema.hcl` for base template validation
- Use `template-context.schema.hcl` for context validation
- Use `template-functions.schema.hcl` for function validation
- Use `template-metadata.schema.hcl` for metadata validation

**Implementation**:
```go
func (v *TemplateValidator) validateAgainstSchema(ctx context.Context, template *spookytypes.Template, schemaType string) error {
    schema, err := v.schemaValidator.GetSchema(schemaType)
    if err != nil {
        return fmt.Errorf("failed to get schema %s: %w", schemaType, err)
    }
    
    return v.schemaValidator.ValidateAgainstSchema(template, schema)
}
```

### Phase 3: Security and Performance

#### 3.1 Template Security Manager
**File**: `internal/templates/security.go` (new)

**Implementation**:
- Implement template sandboxing and access control
- Add pattern-based security filtering
- Implement resource limiting and monitoring
- Add audit logging for template operations

**Key Features**:
```go
type TemplateSecurityManager struct {
    sandbox     TemplateSandbox
    accessControl AccessController
    patternFilter PatternFilter
    resourceMonitor ResourceMonitor
    auditLogger AuditLogger
}

func (s *TemplateSecurityManager) ValidateTemplateSecurity(ctx context.Context, template *spookytypes.Template) error
func (s *TemplateSecurityManager) SandboxTemplate(ctx context.Context, template *spookytypes.Template) (*TemplateSandbox, error)
func (s *TemplateSecurityManager) CheckAccessControl(ctx context.Context, template *spookytypes.Template, user string) error
func (s *TemplateSecurityManager) FilterDangerousPatterns(ctx context.Context, template *spookytypes.Template) error
```

#### 3.2 Template Performance Optimization
**File**: `internal/templates/performance.go` (new)

**Implementation**:
- Implement template result caching
- Add template compilation optimization
- Implement parallel template processing
- Add performance monitoring and metrics

**Key Features**:
```go
type TemplatePerformanceManager struct {
    cache       TemplateResultCache
    compiler    TemplateCompiler
    parallelizer ParallelProcessor
    monitor     PerformanceMonitor
}

func (p *TemplatePerformanceManager) CacheTemplateResult(ctx context.Context, template *spookytypes.Template, data map[string]interface{}, result string) error
func (p *TemplatePerformanceManager) GetCachedResult(ctx context.Context, template *spookytypes.Template, data map[string]interface{}) (string, bool)
func (p *TemplatePerformanceManager) CompileTemplate(ctx context.Context, template *spookytypes.Template) (*CompiledTemplate, error)
func (p *TemplatePerformanceManager) ProcessTemplatesParallel(ctx context.Context, templates []*spookytypes.Template, data map[string]interface{}) ([]TemplateResult, error)
```

### Phase 4: Advanced Features

#### 4.1 Template Composition and Inheritance
**File**: `internal/templates/composition.go` (new)

**Implementation**:
- Implement template inheritance patterns
- Add template composition and merging
- Implement template override mechanisms
- Add template dependency resolution

**Key Features**:
```go
type TemplateCompositionManager struct {
    inheritanceManager InheritanceManager
    compositionManager CompositionManager
    dependencyResolver DependencyResolver
}

func (c *TemplateCompositionManager) InheritTemplate(ctx context.Context, baseTemplate *spookytypes.Template, overrides map[string]interface{}) (*spookytypes.Template, error)
func (c *TemplateCompositionManager) ComposeTemplates(ctx context.Context, templates []*spookytypes.Template) (*spookytypes.Template, error)
func (c *TemplateCompositionManager) ResolveDependencies(ctx context.Context, template *spookytypes.Template) ([]*spookytypes.Template, error)
```

#### 4.2 Template Discovery and Management
**File**: `internal/templates/discovery.go` (new)

**Implementation**:
- Implement template discovery and indexing
- Add template search and filtering
- Implement template categorization and tagging
- Add template versioning and lifecycle management

**Key Features**:
```go
type TemplateDiscoveryManager struct {
    indexer     TemplateIndexer
    searcher    TemplateSearcher
    categorizer TemplateCategorizer
    versionManager VersionManager
}

func (d *TemplateDiscoveryManager) IndexTemplates(ctx context.Context, directory string) error
func (d *TemplateDiscoveryManager) SearchTemplates(ctx context.Context, query string, filters map[string]interface{}) ([]*spookytypes.Template, error)
func (d *TemplateDiscoveryManager) CategorizeTemplate(ctx context.Context, template *spookytypes.Template) error
func (d *TemplateDiscoveryManager) ManageTemplateVersion(ctx context.Context, template *spookytypes.Template, version string) error
```

## Integration Points

### 1. Integration Manager Enhancement
**File**: `internal/integration/manager.go`

**Enhancements**:
- Ensure `GetTemplatesIntegration()` returns proper integration
- Add template integration health checks
- Implement template integration metrics

### 2. Integration Factory Enhancement
**File**: `internal/integration/factory.go`

**Enhancements**:
- Fix `createTemplatesIntegration()` method to return proper integration instead of nil
- Add proper dependency injection for template components
- Implement template integration configuration

### 3. CLI Integration
**File**: `cmd/templates.go` (new)

**Implementation**:
- Add `spooky templates` command group
- Implement `spooky templates render` command
- Implement `spooky templates validate` command
- Implement `spooky templates list` command
- Implement `spooky templates search` command

**Commands**:
```bash
spooky templates render <project> <template> --data <file> --output <file>
spooky templates validate <project> [--template <path>]
spooky templates list <project> [--format json|hcl]
spooky templates search <project> <query> [--tags <tags>] [--category <category>]
```

## Testing Strategy

### 1. Unit Tests
- Test template manager functionality
- Test template function registry
- Test template context resolution
- Test template metadata management
- Test template validation

### 2. Integration Tests
- Test template integration with facts system
- Test template integration with variables system
- Test template integration with machines system
- Test template integration with secrets system

### 3. Performance Tests
- Test template rendering performance
- Test template caching effectiveness
- Test parallel template processing
- Test memory usage under load

### 4. Security Tests
- Test template sandboxing
- Test pattern filtering
- Test access control
- Test resource limiting

## Migration Strategy

### 1. Backward Compatibility
- Maintain existing template manager interface
- Ensure existing template files continue to work
- Provide migration tools for template metadata

### 2. Gradual Rollout
- Phase 1: Core enhancements (non-breaking)
- Phase 2: Schema integration (with fallbacks)
- Phase 3: Security features (opt-in)
- Phase 4: Advanced features (opt-in)

### 3. Documentation Updates
- Update template documentation
- Add template examples and tutorials
- Document template security best practices
- Provide template migration guides

## Success Metrics

### 1. Functionality Metrics
- Template rendering success rate
- Template validation accuracy
- Template function coverage
- Template context resolution accuracy

### 2. Performance Metrics
- Template rendering speed
- Template cache hit rate
- Memory usage efficiency
- Parallel processing effectiveness

### 3. Security Metrics
- Security violation detection rate
- Pattern filtering effectiveness
- Access control enforcement
- Audit log completeness

### 4. Usability Metrics
- Template discovery effectiveness
- Search result relevance
- Template categorization accuracy
- User satisfaction with template system

## Timeline

### Phase 1: Core Enhancement (2-3 weeks)
- Week 1: Template manager enhancement
- Week 2: Template function registry
- Week 3: Template context resolver and metadata manager

### Phase 2: Schema Integration (1-2 weeks)
- Week 1: Schema-driven validation
- Week 2: Schema integration testing

### Phase 3: Security and Performance (2-3 weeks)
- Week 1: Security manager implementation
- Week 2: Performance optimization
- Week 3: Security and performance testing

### Phase 4: Advanced Features (2-3 weeks)
- Week 1: Template composition and inheritance
- Week 2: Template discovery and management
- Week 3: CLI integration and testing

### Total Timeline: 7-11 weeks

## Dependencies

### Internal Dependencies
- Facts system integration
- Variables system integration
- Machines system integration
- Secrets system integration
- Schema system integration

### External Dependencies
- Go text/template engine
- HCL parsing and validation
- Caching libraries
- Security libraries
- Performance monitoring libraries

## Risk Mitigation

### 1. Technical Risks
- **Risk**: Template security vulnerabilities
- **Mitigation**: Comprehensive security testing and sandboxing

- **Risk**: Performance degradation
- **Mitigation**: Performance testing and optimization

- **Risk**: Schema compatibility issues
- **Mitigation**: Backward compatibility and migration tools

### 2. Integration Risks
- **Risk**: Integration with existing systems
- **Mitigation**: Comprehensive integration testing

- **Risk**: Breaking existing functionality
- **Mitigation**: Backward compatibility and gradual rollout

### 3. User Adoption Risks
- **Risk**: Complex template system
- **Mitigation**: Comprehensive documentation and examples

- **Risk**: Migration complexity
- **Mitigation**: Migration tools and guides
