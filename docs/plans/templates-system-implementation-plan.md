# Templates System Implementation Plan

## Overview

This document provides a comprehensive implementation plan for the templates system in spooky. It covers the current state analysis, implementation phases, and detailed technical specifications for building a robust template rendering and management system.

**Schema Integration**: This implementation plan follows the schema validation patterns and template configuration definitions defined in the existing template system design.

**Architecture Integration**: Templates integrate with the overall spooky architecture, providing dynamic content generation and configuration management for all system components.

## Current State Analysis

### What We Have ✅
- **Basic template types** in `internal/types/templates/template.go`
- **Template manager** with basic loading and rendering in `internal/templates/manager.go`
- **Template integration** implementing `TemplatesIntegration` interface in `internal/templates/integration.go`
- **Template schemas** defined in `internal/schemas/schemas/template-structure.schema.hcl`
- **Template system design documentation** in `docs/plans/template-system-design.md`
- **Integration manager** that includes `GetTemplatesIntegration()` method
- **Integration factory** with placeholder for templates integration creation

### What We Need 🔄
- **CLI commands** for template management (`cmd/templates.go` missing)
- **Enhanced template engine** with Go templates and custom functions (currently uses simple string replacement)
- **Template context management** with facts, variables, and machine data
- **Template validation** with comprehensive error reporting
- **Template preview and debugging** capabilities
- **Template caching** and performance optimization
- **Template security** and sandboxing

### Critical Issues ❌
- **Integration factory returns nil** for templates integration in `internal/integration/factory.go`
- **Template manager uses string replacement** instead of proper Go template engine
- **No CLI commands** exist for template management
- **No template context management** system exists
- **No placeholder implementations** are acceptable - all code must be fully functional

## Implementation Plan

### Phase 0: Critical Fixes (Immediate)

#### 0.1 Fix Integration Factory
**File**: `internal/integration/factory.go`
- **Fix `createTemplatesIntegration()` method** to return proper integration instead of nil
- **Create template manager instance** with proper logger dependency
- **Create template integration instance** and return it
- **Add proper logging** for templates integration creation
- **Ensure integration manager** can access templates integration

**Implementation** (fully functional, no placeholders):
```go
import (
    spookyinterfaces "spooky/internal/interfaces"
    spookylogging "spooky/internal/logging"
    spookytemplates "spooky/internal/templates"
    spookytypes "spooky/internal/types"
)

func (f *Factory) createTemplatesIntegration() spookyinterfaces.TemplatesIntegration {
    // Create log manager for templates
    logManager := spookylogging.NewLogManager()
    templatesLogger := logManager.GetLogger("templates")

    // Create templates manager with full functionality
    manager := spookytemplates.NewManager(templatesLogger)

    // Create templates integration with complete implementation
    integration := spookytemplates.NewIntegration(manager)

    f.logger.Info("Templates integration created successfully")
    return integration
}
```

### Phase 1: Core Template Engine Enhancement

#### 1.1 Enhanced Template Manager
**File**: `internal/templates/manager.go`
- **Replace simple string replacement** with proper Go template engine
- **Add custom template functions** for facts, variables, environment data
- **Implement template context creation** with project, facts, machines data
- **Add template caching** for performance
- **Implement template validation** with syntax checking and function validation

**Current Implementation Issues**:
```go
// CURRENT - Simple string replacement (needs replacement)
result := template.Content
for key, value := range data {
    placeholder := fmt.Sprintf("{{.%s}}", key)
    result = strings.ReplaceAll(result, placeholder, fmt.Sprintf("%v", value))
}
```

**Required Implementation** (complete, no placeholders):
```go
// NEW - Proper Go template engine with full functionality
tmpl, err := template.New(template.ID).Funcs(template.FuncMap{
    "custom": func(path string) interface{} { 
        return getNestedValue(data.Custom, path) 
    },
    "system": func(path string) interface{} { 
        return getNestedValue(data.System, path) 
    },
    "env": func(key string) string { 
        return data.Environment[key] 
    },
    "var": func(name string) interface{} { 
        return data.Variables[name] 
    },
    "varOrDefault": func(name string, defaultValue interface{}) interface{} {
        if value, exists := data.Variables[name]; exists {
            return value
        }
        return defaultValue
    },
    "machines": func() []Machine {
        return data.Machines
    },
    "machine": func(name string) *Machine {
        for _, m := range data.Machines {
            if m.Name == name {
                return &m
            }
        }
        return nil
    },
    "data": func(path string) interface{} {
        return getNestedValue(data.CustomData, path)
    }
}).Parse(template.Content)
```

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

**CLI Structure** (following established patterns):
```go
// templatesCmd represents the templates command
var templatesCmd = &cobra.Command{
    Use:   "templates",
    Short: "Manage project templates",
    Long:  `Manage and render project templates with dynamic data integration.`,
}

// templatesListCmd represents the templates list command
var templatesListCmd = &cobra.Command{
    Use:   "list [project-path]",
    Short: "List project templates",
    Long:  `List all templates in the specified project.`,
    Args:  cobra.ExactArgs(1),
    RunE:  handleTemplatesList,
}

// templatesValidateCmd represents the templates validate command
var templatesValidateCmd = &cobra.Command{
    Use:   "validate [project-path]",
    Short: "Validate project templates",
    Long:  `Validate template syntax and configuration.`,
    Args:  cobra.ExactArgs(1),
    RunE:  handleTemplatesValidate,
}

// templatesRenderCmd represents the templates render command
var templatesRenderCmd = &cobra.Command{
    Use:   "render [project-path] [template-path]",
    Short: "Render a template",
    Long:  `Render a template with project data and facts.`,
    Args:  cobra.ExactArgs(2),
    RunE:  handleTemplatesRender,
}
```

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
- **Template integration initialization** already handled in `InitializeIntegrationsDependencies()`
- **Template manager instance** creation handled in integration factory
- **Template integration registration** handled in integration manager
- **Template logging** setup handled in integration factory

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

## Updated Implementation Order

### Priority 0: Critical Fixes (Immediate)
1. **Fix integration factory** - Implement proper templates integration creation (fully functional)
2. **Verify integration manager** - Ensure templates integration is accessible
3. **Remove all placeholder code** - Replace any "not implemented" or stub functions with complete implementations
4. **Follow import alias patterns** - Use `spookytemplates` and `spookytypestemplates` consistently

### Priority 1: Core Functionality
1. **Enhanced template manager** - Replace string replacement with complete Go template engine implementation
2. **Template context management** - Implement complete context creation and data integration using `spookytypes.TemplatesContext`
3. **CLI commands** - Create `cmd/templates.go` with fully functional template commands
4. **Template validation** - Implement comprehensive validation with complete error reporting using `spookytypes.TemplateValidationError`

### Priority 2: Integration
1. **Facts integration** - Complete template data access through facts integration using `spookyinterfaces.FactsIntegration`
2. **Variables integration** - Complete template data access through variables integration using `spookyinterfaces.VariablesIntegration`
3. **Machines integration** - Complete machine-specific rendering context using `spookyinterfaces.MachinesIntegration`
4. **Template error types** - Complete dedicated error types and validation results using `spookytypes.TemplateError`

### Priority 3: Advanced Features
1. **Template preview** - Complete template analysis and mock rendering functionality
2. **Template debugging** - Complete step-by-step rendering and error context functionality
3. **Template security** - Complete function whitelist and pattern blacklist implementation
4. **Template caching** - Complete performance optimization with TTL implementation

### Priority 4: Testing and Documentation
1. **Comprehensive testing** - Complete unit, integration, and performance test suites
2. **Template examples** - Complete example templates and documentation
3. **CLI help** - Complete usage documentation and examples
4. **Performance benchmarking** - Complete template rendering performance test suite

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
- ✅ **Template rendering** works with all data sources (no placeholder implementations)
- ✅ **CLI commands** provide complete template management (fully functional)
- ✅ **Template validation** catches all syntax and security issues (comprehensive implementation)
- ✅ **Template preview** shows accurate analysis and mock rendering (complete functionality)
- ✅ **Integration** works seamlessly with all other systems (no stub implementations)
- ✅ **Import aliases** follow established patterns (`spookytemplates`, `spookytypestemplates`)
- ✅ **Type safety** maintained through proper interface usage and type composition

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
- ✅ **Comprehensive testing** covers all functionality (no untested code)
- ✅ **Error handling** provides actionable feedback (complete error handling)
- ✅ **Documentation** is complete and accurate (no incomplete documentation)
- ✅ **Code quality** meets project standards (no placeholder or stub code)

## Risk Mitigation

### Technical Risks
- **Template engine complexity** - Mitigation: Complete implementation with thorough testing (no partial implementations)
- **Performance bottlenecks** - Mitigation: Complete caching and optimization strategies (no placeholder optimizations)
- **Security vulnerabilities** - Mitigation: Complete security validation and sandboxing (no stub security)
- **Integration issues** - Mitigation: Complete interface-based design with proper error handling (no incomplete integrations)

### Implementation Risks
- **Scope creep** - Mitigation: Strict adherence to implementation phases with complete functionality at each phase
- **Testing complexity** - Mitigation: Complete test scenarios and automation (no partial testing)
- **Performance issues** - Mitigation: Complete performance testing and optimization (no placeholder performance)
- **Integration challenges** - Mitigation: Complete integration with existing systems (no stub integrations)

## Critical Dependencies

### Phase 0 Dependencies
- **Integration factory fix** must be completed before any other work (fully functional implementation)
- **Integration manager** must properly initialize templates integration (complete initialization)

### Phase 1 Dependencies
- **Template manager enhancement** depends on integration factory fix (complete enhancement, no placeholders)
- **CLI commands** depend on template manager enhancement (fully functional commands)
- **Template context** depends on facts, variables, and machines integrations (complete context implementation)

### Phase 2 Dependencies
- **Template validation** depends on template engine implementation (complete validation)
- **Template error types** depend on validation implementation (complete error handling)
- **Integration features** depend on context management (complete integration features)

## Conclusion

This updated implementation plan addresses the critical integration issue and provides a clear roadmap for building a robust templates system. The phased approach ensures that core functionality is delivered first, followed by advanced features and optimizations.

**Critical Implementation Principle**: No placeholder implementations are acceptable. Every component must be fully functional and complete before moving to the next phase. This includes:
- No "not implemented" errors or stub functions
- No TODO comments without immediate implementation
- No partial implementations that return errors or empty results
- Complete functionality at each phase before proceeding

The plan maintains consistency with existing codebase patterns, follows established architectural principles, and provides clear success criteria for implementation validation. By following this plan, the templates system will provide powerful dynamic content generation capabilities while maintaining security, performance, and usability standards.

**Critical Note**: The integration factory fix (Phase 0) must be completed before any other work can proceed effectively, as it's a blocking dependency for all subsequent phases. All implementations must be complete and functional.
