# Templates System Implementation Plan

## Overview

This document provides a comprehensive implementation plan for the templates system in spooky. It covers the current state analysis, implementation phases, and detailed technical specifications for building a robust template rendering and management system.

**Schema Integration**: This implementation plan follows the schema validation patterns and template configuration definitions defined in the existing template system design.

**Architecture Integration**: Templates integrate with the overall spooky architecture, providing dynamic content generation and configuration management for all system components.

## Current State Analysis

### What We Have ✅
- **Basic template types** in `internal/types/templates/template.go`
- **Template manager** with basic loading and rendering in `internal/templates/manager.go`
- **Template integration** implementing `TemplatesIntegration` interface
- **Template schemas** defined in `internal/schemas/schemas/template-structure.schema.hcl`
- **Template system design documentation** in `docs/plans/template-system-design.md`
- **Integration manager** that includes `GetTemplatesIntegration()` method

### What We Need 🔄
- **CLI commands** for template management
- **Enhanced template engine** with Go templates and custom functions
- **Template context management** with facts, variables, and machine data
- **Template validation** with comprehensive error reporting
- **Template preview and debugging** capabilities
- **Template caching** and performance optimization
- **Template security** and sandboxing

## Implementation Plan

### Phase 1: Core Template Engine Enhancement

#### 1.1 Enhanced Template Manager
**File**: `internal/templates/manager.go`
- **Replace simple string replacement** with proper Go template engine
- **Add custom template functions** for facts, variables, environment data
- **Implement template context creation** with project, facts, machines data
- **Add template caching** for performance
- **Implement template validation** with syntax checking and function validation

#### 1.2 Template Context Management
**File**: `internal/templates/context.go` (new)
- **Create TemplateContext struct** with all available data
- **Implement context creation** from project, facts, variables, machines
- **Add template functions** for data access (`custom()`, `system()`, `env()`, `var()`, etc.)
- **Implement nested data access** with path-based lookups
- **Add environment variable integration**

#### 1.3 Template Engine Implementation
**File**: `internal/templates/engine.go` (new)
- **Implement Go template engine** with custom functions
- **Add template preprocessing** and validation
- **Implement error handling** with line numbers and context
- **Add template function registration** and validation
- **Implement template caching** with TTL

### Phase 2: CLI Commands Implementation

#### 2.1 Template CLI Commands
**File**: `cmd/templates.go` (new)
- **Implement `spooky templates list`** command
- **Implement `spooky templates validate`** command
- **Implement `spooky templates render`** command
- **Implement `spooky templates preview`** command
- **Add proper flag handling** and argument validation
- **Follow established CLI patterns** from other commands

#### 2.2 Template Command Handlers
**File**: `cmd/templates.go` (continued)
- **Implement `handleTemplatesList`** function
- **Implement `handleTemplatesValidate`** function
- **Implement `handleTemplatesRender`** function
- **Implement `handleTemplatesPreview`** function
- **Add project path validation** and error handling
- **Integrate with template manager** through integration manager

### Phase 3: Template Validation and Error Handling

#### 3.1 Enhanced Template Validation
**File**: `internal/templates/validation.go` (new)
- **Implement syntax validation** using Go template parser
- **Add function validation** for custom template functions
- **Implement variable validation** for required template variables
- **Add security validation** for dangerous patterns
- **Implement comprehensive error reporting** with line numbers

#### 3.2 Template Error Types
**File**: `internal/types/templates/errors.go` (new)
- **Define TemplateError struct** with type, message, field, line, column
- **Implement TemplateValidationResult** with errors and warnings
- **Add error categorization** (syntax, function, variable, security)
- **Implement error formatting** for CLI output

### Phase 4: Template Integration and Context

#### 4.1 Facts Integration
**Files**: `internal/templates/context.go`, `internal/templates/functions.go`
- **Integrate with facts system** through integration manager
- **Add facts access functions** (`system()`, `custom()`)
- **Implement facts caching** for performance
- **Add machine-specific facts** support

#### 4.2 Variables Integration
**Files**: `internal/templates/context.go`, `internal/templates/functions.go`
- **Integrate with variables system** through integration manager
- **Add variable access functions** (`var()`, `varOrDefault()`, `varRequired()`)
- **Implement variable resolution** with context
- **Add variable validation** for template requirements

#### 4.3 Machines Integration
**Files**: `internal/templates/context.go`, `internal/templates/functions.go`
- **Integrate with machines system** through integration manager
- **Add machine access functions** (`machines()`, `machine()`)
- **Implement machine filtering** by tags and patterns
- **Add machine-specific context** for rendering

### Phase 5: Template Security and Performance

#### 5.1 Template Security
**File**: `internal/templates/security.go` (new)
- **Implement function whitelist** for allowed template functions
- **Add pattern blacklist** for dangerous patterns
- **Implement execution timeouts** for template rendering
- **Add memory usage limits** for template processing
- **Implement sandboxing** for template execution

#### 5.2 Template Performance
**File**: `internal/templates/cache.go` (new)
- **Implement template caching** with TTL
- **Add parsed template caching** for repeated rendering
- **Implement context caching** for expensive data lookups
- **Add performance metrics** and monitoring
- **Implement lazy loading** for template data

### Phase 6: Template Preview and Debugging

#### 6.1 Template Preview
**File**: `internal/templates/preview.go` (new)
- **Implement template analysis** with variable detection
- **Add mock data generation** for preview rendering
- **Implement template structure analysis** with dependency detection
- **Add template complexity metrics**
- **Implement template documentation generation**

#### 6.2 Template Debugging
**File**: `internal/templates/debug.go` (new)
- **Implement step-by-step template rendering**
- **Add variable value tracing** during rendering
- **Implement template function call logging**
- **Add rendering performance profiling**
- **Implement error context enhancement**

### Phase 7: Integration and Testing

#### 7.1 Integration Manager Updates
**File**: `cmd/integrations.go`
- **Add template integration initialization** to `InitializeIntegrationsDependencies()`
- **Create template manager instance** with proper dependencies
- **Register template integration** with integration manager
- **Add template logging** setup

#### 7.2 Template Testing
**Files**: `internal/templates/*_test.go`
- **Add unit tests** for template manager
- **Add integration tests** for template rendering
- **Add CLI command tests** for template commands
- **Add performance tests** for template caching
- **Add security tests** for template validation

#### 7.3 Template Examples
**Directory**: `examples/templates/`
- **Create example templates** for common use cases
- **Add template documentation** with usage examples
- **Create test data** for template validation
- **Add template best practices** documentation

## Implementation Order

### Priority 1: Core Functionality
1. **Enhanced template manager** with Go template engine
2. **Template context management** with data integration
3. **Basic CLI commands** (list, validate, render)
4. **Template validation** with error reporting

### Priority 2: Integration
1. **Facts integration** for template data access
2. **Variables integration** for dynamic configuration
3. **Machines integration** for machine-specific rendering
4. **Integration manager updates**

### Priority 3: Advanced Features
1. **Template preview** and analysis
2. **Template debugging** capabilities
3. **Template security** and sandboxing
4. **Template caching** and performance optimization

### Priority 4: Testing and Documentation
1. **Comprehensive testing** suite
2. **Template examples** and documentation
3. **CLI help** and usage documentation
4. **Performance benchmarking**

## Key Implementation Details

### Template Functions to Implement
```go
// Core template functions
"custom": func(path string) interface{}     // Access custom facts
"system": func(path string) interface{}     // Access system facts
"env": func(key string) string              // Access environment variables
"var": func(name string) interface{}        // Access project variables
"varOrDefault": func(name, default) interface{}
"machines": func() []Machine                 // Access machine inventory
"machine": func(name string) *Machine       // Access specific machine
"data": func(path string) interface{}       // Access additional data
```

### CLI Command Structure
```bash
spooky templates list <project-path> [--verbose] [--pattern]
spooky templates validate <project-path> [--template] [--verbose]
spooky templates render <project-path> <template> [--data] [--machine] [--output] [--preview] [--dry-run]
spooky templates preview <project-path> <template> [--data]
```

### Template Context Structure
```go
type TemplateContext struct {
    Project     *ProjectConfig
    Facts       map[string]interface{}
    Machines    []*Machine
    Variables   *VariableContext
    Environment map[string]string
    CustomData  map[string]interface{}
}
```

## Technical Specifications

### Template Engine Requirements
- **Go template engine** with custom function support
- **Template caching** with configurable TTL
- **Error handling** with line numbers and context
- **Security sandboxing** for safe template execution
- **Performance monitoring** and metrics collection

### Template Validation Requirements
- **Syntax validation** using Go template parser
- **Function validation** for allowed template functions
- **Variable validation** for required template variables
- **Security validation** for dangerous patterns
- **Comprehensive error reporting** with actionable messages

### CLI Integration Requirements
- **Follow established CLI patterns** from other commands
- **Proper flag handling** and argument validation
- **Integration with existing systems** (facts, variables, machines)
- **Error handling** with user-friendly messages
- **Help documentation** and usage examples

### Performance Requirements
- **Template caching** for repeated rendering
- **Context caching** for expensive data lookups
- **Lazy loading** for template data
- **Memory usage limits** for template processing
- **Execution timeouts** for template rendering

## Integration Points

### Facts System Integration
- **Access to machine facts** through facts integration
- **Custom facts support** for project-specific data
- **Facts caching** for performance optimization
- **Machine-specific facts** for targeted rendering

### Variables System Integration
- **Project variables access** through variables integration
- **Variable resolution** with proper precedence
- **Variable validation** for template requirements
- **Dynamic variable substitution** during rendering

### Machines System Integration
- **Machine inventory access** through machines integration
- **Machine filtering** by tags and patterns
- **Machine-specific context** for targeted rendering
- **SSH integration** for remote template rendering

### Configuration System Integration
- **Template settings** from global configuration
- **Security policies** and restrictions
- **Performance settings** and limits
- **Logging configuration** for template operations

## Success Criteria

### Functional Requirements
- ✅ **Template rendering** works with all data sources
- ✅ **CLI commands** provide full template management
- ✅ **Template validation** catches all syntax and security issues
- ✅ **Template preview** shows accurate analysis and mock rendering
- ✅ **Integration** works seamlessly with all other systems

### Performance Requirements
- ✅ **Template caching** improves rendering performance
- ✅ **Context caching** reduces data lookup overhead
- ✅ **Memory usage** stays within defined limits
- ✅ **Execution time** meets performance benchmarks

### Security Requirements
- ✅ **Function whitelist** prevents dangerous operations
- ✅ **Pattern blacklist** blocks malicious templates
- ✅ **Execution timeouts** prevent infinite loops
- ✅ **Memory limits** prevent resource exhaustion

### Quality Requirements
- ✅ **Comprehensive testing** covers all functionality
- ✅ **Error handling** provides actionable feedback
- ✅ **Documentation** is complete and accurate
- ✅ **Code quality** meets project standards

## Risk Mitigation

### Technical Risks
- **Template engine complexity** - Mitigation: Incremental implementation with thorough testing
- **Performance bottlenecks** - Mitigation: Caching and optimization strategies
- **Security vulnerabilities** - Mitigation: Comprehensive security validation and sandboxing
- **Integration issues** - Mitigation: Interface-based design with proper error handling

### Implementation Risks
- **Scope creep** - Mitigation: Strict adherence to implementation phases
- **Testing complexity** - Mitigation: Comprehensive test scenarios and automation
- **Performance issues** - Mitigation: Early performance testing and optimization
- **Integration challenges** - Mitigation: Incremental integration with existing systems

## Conclusion

This implementation plan provides a comprehensive roadmap for building a robust templates system that integrates seamlessly with the existing spooky architecture. The phased approach ensures that core functionality is delivered first, followed by advanced features and optimizations.

The plan maintains consistency with existing codebase patterns, follows established architectural principles, and provides clear success criteria for implementation validation. By following this plan, the templates system will provide powerful dynamic content generation capabilities while maintaining security, performance, and usability standards.
