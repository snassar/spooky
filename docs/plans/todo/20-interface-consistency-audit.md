# Implementation Plan: Interface Consistency Audit

## Overview
Audit and standardize all interfaces across the spooky codebase to ensure consistency, proper implementation, and maintainable architecture.

## Task Details
- **Task ID**: 8.1
- **Priority**: Medium
- **Files**: All `internal/` packages with interfaces
- **Functions**: Interface auditing, standardization, implementation verification

## Current State Analysis

### Existing Interface Patterns
1. **Interface-First Design**: Most packages use interface-based design
2. **Integration Interfaces**: Coordinator uses integration interfaces for subsystem communication
3. **Type Conversion**: Complex type conversion logic scattered across packages
4. **Error Handling**: Inconsistent error wrapping patterns

### Interface Inconsistencies Found
1. **Type Conversion Complexity**: Multiple packages have complex interface-to-concrete type conversions
2. **Interface Implementation**: Some interfaces are not consistently implemented
3. **Error Handling**: Inconsistent error wrapping and propagation patterns
4. **Method Signatures**: Some interfaces have inconsistent method signatures

## Implementation Requirements

### Interface Compliance
The interface consistency audit must:
1. **Audit all interfaces** across the codebase
2. **Standardize interface patterns** and naming conventions
3. **Ensure consistent implementation** of all interfaces
4. **Create centralized type conversion utilities**
5. **Standardize error handling patterns**
6. **Document interface contracts** and usage patterns

### Required Dependencies
- All existing packages and interfaces
- Type conversion utilities
- Error handling standardization
- Documentation system

## Detailed Implementation Plan

### Step 1: Interface Audit and Documentation

#### 1.1 Create Interface Inventory
```go
// internal/interfaces/audit/inventory.go
package audit

import (
    "reflect"
    "spooky/internal/interfaces"
)

// InterfaceInventory represents the complete interface inventory
type InterfaceInventory struct {
    Interfaces map[string]*InterfaceInfo
    Packages   map[string]*PackageInfo
}

// InterfaceInfo contains information about an interface
type InterfaceInfo struct {
    Name           string
    Package        string
    Methods        []MethodInfo
    Implementations []string
    Usage          []string
    Status         InterfaceStatus
}

// MethodInfo contains information about an interface method
type MethodInfo struct {
    Name     string
    Signature string
    Returns  []string
    Params   []ParamInfo
}

// ParamInfo contains information about a method parameter
type ParamInfo struct {
    Name string
    Type string
}

// InterfaceStatus represents the status of an interface
type InterfaceStatus string

const (
    InterfaceStatusComplete    InterfaceStatus = "complete"
    InterfaceStatusIncomplete  InterfaceStatus = "incomplete"
    InterfaceStatusDeprecated  InterfaceStatus = "deprecated"
    InterfaceStatusUnused      InterfaceStatus = "unused"
)

// NewInterfaceInventory creates a new interface inventory
func NewInterfaceInventory() *InterfaceInventory {
    return &InterfaceInventory{
        Interfaces: make(map[string]*InterfaceInfo),
        Packages:   make(map[string]*PackageInfo),
    }
}

// AuditInterfaces audits all interfaces in the codebase
func (ii *InterfaceInventory) AuditInterfaces() error {
    // Audit interfaces in each package
    packages := []string{
        "actions", "facts", "machines", "templates", 
        "variables", "storage", "ssh", "config", "cli"
    }
    
    for _, pkg := range packages {
        if err := ii.auditPackage(pkg); err != nil {
            return fmt.Errorf("failed to audit package %s: %w", pkg, err)
        }
    }
    
    return nil
}

// auditPackage audits interfaces in a specific package
func (ii *InterfaceInventory) auditPackage(packageName string) error {
    // Scan package for interfaces
    interfaces, err := ii.scanPackageInterfaces(packageName)
    if err != nil {
        return err
    }
    
    // Analyze each interface
    for _, iface := range interfaces {
        info := &InterfaceInfo{
            Name:    iface.Name,
            Package: packageName,
            Methods: ii.extractMethods(iface),
        }
        
        // Check implementations
        info.Implementations = ii.findImplementations(iface)
        
        // Check usage
        info.Usage = ii.findUsage(iface)
        
        // Determine status
        info.Status = ii.determineStatus(info)
        
        ii.Interfaces[iface.Name] = info
    }
    
    return nil
}
```

#### 1.2 Interface Documentation Generator
```go
// internal/interfaces/audit/documentation.go
package audit

import (
    "fmt"
    "os"
    "path/filepath"
    "text/template"
)

// DocumentationGenerator generates interface documentation
type DocumentationGenerator struct {
    inventory *InterfaceInventory
}

// NewDocumentationGenerator creates a new documentation generator
func NewDocumentationGenerator(inventory *InterfaceInventory) *DocumentationGenerator {
    return &DocumentationGenerator{
        inventory: inventory,
    }
}

// GenerateInterfaceDocumentation generates documentation for all interfaces
func (dg *DocumentationGenerator) GenerateInterfaceDocumentation(outputPath string) error {
    // Create output directory
    if err := os.MkdirAll(outputPath, 0755); err != nil {
        return fmt.Errorf("failed to create output directory: %w", err)
    }
    
    // Generate main interface documentation
    if err := dg.generateMainDocumentation(outputPath); err != nil {
        return err
    }
    
    // Generate package-specific documentation
    if err := dg.generatePackageDocumentation(outputPath); err != nil {
        return err
    }
    
    return nil
}

// generateMainDocumentation generates the main interface documentation
func (dg *DocumentationGenerator) generateMainDocumentation(outputPath string) error {
    templateContent := `# Interface Documentation

## Overview
This document provides comprehensive documentation for all interfaces in the spooky codebase.

## Interface Status Summary
{{range $name, $info := .Interfaces}}
### {{$name}} ({{$info.Package}})
- **Status**: {{$info.Status}}
- **Methods**: {{len $info.Methods}}
- **Implementations**: {{len $info.Implementations}}
- **Usage**: {{len $info.Usage}}

{{range $info.Methods}}
#### {{.Name}}
```go
{{.Signature}}
```
{{end}}
{{end}}
`
    
    tmpl, err := template.New("interfaces").Parse(templateContent)
    if err != nil {
        return fmt.Errorf("failed to parse template: %w", err)
    }
    
    file, err := os.Create(filepath.Join(outputPath, "interfaces.md"))
    if err != nil {
        return fmt.Errorf("failed to create documentation file: %w", err)
    }
    defer file.Close()
    
    return tmpl.Execute(file, dg.inventory)
}
```

### Step 2: Interface Standardization

#### 2.1 Create Interface Standards
```go
// internal/interfaces/standards/standards.go
package standards

// InterfaceStandards defines standards for interface design
type InterfaceStandards struct {
    NamingConventions    NamingConventions
    MethodPatterns       MethodPatterns
    ErrorHandling        ErrorHandling
    Documentation        Documentation
}

// NamingConventions defines naming conventions for interfaces
type NamingConventions struct {
    InterfaceSuffixes []string // e.g., ["Manager", "Integration", "Provider"]
    MethodPrefixes    []string // e.g., ["Get", "Set", "Create", "Delete"]
    PackagePrefixes   []string // e.g., ["spooky"]
}

// MethodPatterns defines standard method patterns
type MethodPatterns struct {
    CRUDOperations    bool // Create, Read, Update, Delete
    ValidationMethods bool // Validate, ValidateAll
    ContextSupport    bool // Context.Context parameter
    ErrorWrapping     bool // Consistent error wrapping
}

// ErrorHandling defines error handling standards
type ErrorHandling struct {
    WrappingPattern string // e.g., "failed to %s: %w"
    ErrorTypes      []string // e.g., ["ValidationError", "NotFoundError"]
    ContextErrors   bool // Include context in errors
}

// Documentation defines documentation standards
type Documentation struct {
    RequiredComments bool // All methods must have comments
    ExamplesRequired bool // Examples required for public interfaces
    UsagePatterns    bool // Document usage patterns
}

// NewInterfaceStandards creates new interface standards
func NewInterfaceStandards() *InterfaceStandards {
    return &InterfaceStandards{
        NamingConventions: NamingConventions{
            InterfaceSuffixes: []string{"Manager", "Integration", "Provider", "Client", "Service"},
            MethodPrefixes:    []string{"Get", "Set", "Create", "Delete", "Update", "Validate"},
            PackagePrefixes:   []string{"spooky"},
        },
        MethodPatterns: MethodPatterns{
            CRUDOperations:    true,
            ValidationMethods: true,
            ContextSupport:    true,
            ErrorWrapping:     true,
        },
        ErrorHandling: ErrorHandling{
            WrappingPattern: "failed to %s: %w",
            ErrorTypes:      []string{"ValidationError", "NotFoundError", "ConfigurationError"},
            ContextErrors:   true,
        },
        Documentation: Documentation{
            RequiredComments: true,
            ExamplesRequired: true,
            UsagePatterns:    true,
        },
    }
}
```

#### 2.2 Interface Validator
```go
// internal/interfaces/standards/validator.go
package standards

import (
    "fmt"
    "reflect"
    "strings"
)

// InterfaceValidator validates interfaces against standards
type InterfaceValidator struct {
    standards *InterfaceStandards
}

// NewInterfaceValidator creates a new interface validator
func NewInterfaceValidator(standards *InterfaceStandards) *InterfaceValidator {
    return &InterfaceValidator{
        standards: standards,
    }
}

// ValidateInterface validates an interface against standards
func (iv *InterfaceValidator) ValidateInterface(iface interface{}) (*ValidationResult, error) {
    result := &ValidationResult{
        Valid:   true,
        Issues:  make([]ValidationIssue, 0),
        Warnings: make([]ValidationIssue, 0),
    }
    
    ifaceType := reflect.TypeOf(iface)
    if ifaceType.Kind() != reflect.Interface {
        return nil, fmt.Errorf("provided type is not an interface")
    }
    
    // Validate naming conventions
    if err := iv.validateNaming(ifaceType, result); err != nil {
        return nil, err
    }
    
    // Validate method patterns
    if err := iv.validateMethods(ifaceType, result); err != nil {
        return nil, err
    }
    
    // Validate error handling
    if err := iv.validateErrorHandling(ifaceType, result); err != nil {
        return nil, err
    }
    
    return result, nil
}

// ValidationResult represents the result of interface validation
type ValidationResult struct {
    Valid    bool
    Issues   []ValidationIssue
    Warnings []ValidationIssue
}

// ValidationIssue represents a validation issue
type ValidationIssue struct {
    Type        IssueType
    Message     string
    Severity    Severity
    Location    string
    Suggestion  string
}

// IssueType represents the type of validation issue
type IssueType string

const (
    IssueTypeNaming      IssueType = "naming"
    IssueTypeMethod      IssueType = "method"
    IssueTypeError       IssueType = "error"
    IssueTypeDocumentation IssueType = "documentation"
)

// Severity represents the severity of a validation issue
type Severity string

const (
    SeverityError   Severity = "error"
    SeverityWarning Severity = "warning"
    SeverityInfo    Severity = "info"
)

// validateNaming validates interface naming conventions
func (iv *InterfaceValidator) validateNaming(ifaceType reflect.Type, result *ValidationResult) error {
    ifaceName := ifaceType.Name()
    
    // Check for valid suffixes
    hasValidSuffix := false
    for _, suffix := range iv.standards.NamingConventions.InterfaceSuffixes {
        if strings.HasSuffix(ifaceName, suffix) {
            hasValidSuffix = true
            break
        }
    }
    
    if !hasValidSuffix {
        result.Issues = append(result.Issues, ValidationIssue{
            Type:       IssueTypeNaming,
            Message:    fmt.Sprintf("Interface %s does not have a valid suffix", ifaceName),
            Severity:   SeverityWarning,
            Suggestion: fmt.Sprintf("Consider adding one of: %v", iv.standards.NamingConventions.InterfaceSuffixes),
        })
    }
    
    return nil
}

// validateMethods validates interface method patterns
func (iv *InterfaceValidator) validateMethods(ifaceType reflect.Type, result *ValidationResult) error {
    for i := 0; i < ifaceType.NumMethod(); i++ {
        method := ifaceType.Method(i)
        
        // Check for context support
        if iv.standards.MethodPatterns.ContextSupport {
            if !iv.hasContextParameter(method) {
                result.Warnings = append(result.Warnings, ValidationIssue{
                    Type:       IssueTypeMethod,
                    Message:    fmt.Sprintf("Method %s should support context.Context", method.Name),
                    Severity:   SeverityWarning,
                    Suggestion: "Consider adding context.Context as the first parameter",
                })
            }
        }
        
        // Check for error return
        if !iv.hasErrorReturn(method) {
            result.Warnings = append(result.Warnings, ValidationIssue{
                Type:       IssueTypeMethod,
                Message:    fmt.Sprintf("Method %s should return error", method.Name),
                Severity:   SeverityWarning,
                Suggestion: "Consider adding error as the last return value",
            })
        }
    }
    
    return nil
}

// hasContextParameter checks if a method has context.Context parameter
func (iv *InterfaceValidator) hasContextParameter(method reflect.Method) bool {
    methodType := method.Type
    if methodType.NumIn() > 0 {
        firstParam := methodType.In(0)
        return firstParam.String() == "context.Context"
    }
    return false
}

// hasErrorReturn checks if a method returns error
func (iv *InterfaceValidator) hasErrorReturn(method reflect.Method) bool {
    methodType := method.Type
    if methodType.NumOut() > 0 {
        lastReturn := methodType.Out(methodType.NumOut() - 1)
        return lastReturn.String() == "error"
    }
    return false
}
```

### Step 3: Implementation Verification

#### 3.1 Implementation Checker
```go
// internal/interfaces/audit/implementation.go
package audit

import (
    "fmt"
    "reflect"
    "strings"
)

// ImplementationChecker checks interface implementations
type ImplementationChecker struct {
    inventory *InterfaceInventory
}

// NewImplementationChecker creates a new implementation checker
func NewImplementationChecker(inventory *InterfaceInventory) *ImplementationChecker {
    return &ImplementationChecker{
        inventory: inventory,
    }
}

// CheckImplementations checks all interface implementations
func (ic *ImplementationChecker) CheckImplementations() (*ImplementationReport, error) {
    report := &ImplementationReport{
        Implementations: make(map[string]*ImplementationInfo),
        Issues:          make([]ImplementationIssue, 0),
    }
    
    for ifaceName, ifaceInfo := range ic.inventory.Interfaces {
        implInfo := &ImplementationInfo{
            InterfaceName: ifaceName,
            Implementations: make([]string, 0),
            Issues:         make([]ImplementationIssue, 0),
        }
        
        // Find all implementations
        implementations, err := ic.findImplementations(ifaceName)
        if err != nil {
            return nil, fmt.Errorf("failed to find implementations for %s: %w", ifaceName, err)
        }
        
        implInfo.Implementations = implementations
        
        // Check each implementation
        for _, impl := range implementations {
            if err := ic.checkImplementation(ifaceName, impl, implInfo); err != nil {
                implInfo.Issues = append(implInfo.Issues, ImplementationIssue{
                    Type:    "check_error",
                    Message: err.Error(),
                    Target:  impl,
                })
            }
        }
        
        report.Implementations[ifaceName] = implInfo
    }
    
    return report, nil
}

// ImplementationReport represents the implementation check report
type ImplementationReport struct {
    Implementations map[string]*ImplementationInfo
    Issues          []ImplementationIssue
}

// ImplementationInfo contains information about interface implementations
type ImplementationInfo struct {
    InterfaceName  string
    Implementations []string
    Issues         []ImplementationIssue
}

// ImplementationIssue represents an implementation issue
type ImplementationIssue struct {
    Type    string
    Message string
    Target  string
}

// checkImplementation checks a specific implementation
func (ic *ImplementationChecker) checkImplementation(ifaceName, implName string, info *ImplementationInfo) error {
    // Get interface type
    ifaceType, err := ic.getInterfaceType(ifaceName)
    if err != nil {
        return err
    }
    
    // Get implementation type
    implType, err := ic.getImplementationType(implName)
    if err != nil {
        return err
    }
    
    // Check if implementation satisfies interface
    if !implType.Implements(ifaceType) {
        info.Issues = append(info.Issues, ImplementationIssue{
            Type:    "interface_not_implemented",
            Message: fmt.Sprintf("Type %s does not implement interface %s", implName, ifaceName),
            Target:  implName,
        })
    }
    
    return nil
}
```

### Step 4: Centralized Type Conversion Utilities

#### 4.1 Type Conversion Package
```go
// internal/types/converter/converter.go
package converter

import (
    "fmt"
    "reflect"
    "spooky/internal/facts/types"
    "spooky/internal/interfaces"
    "spooky/internal/variables/types"
)

// Converter provides centralized type conversion utilities
type Converter struct{}

// NewConverter creates a new converter
func NewConverter() *Converter {
    return &Converter{}
}

// ConvertFactsToConcrete converts interface facts to concrete types
func (c *Converter) ConvertFactsToConcrete(factsContext interfaces.FactsContext) map[string]interface{} {
    if factsContext == nil {
        return make(map[string]interface{})
    }
    
    concreteFacts := make(map[string]interface{})
    
    // Handle different fact context types
    switch v := factsContext.(type) {
    case map[string]interface{}:
        concreteFacts = v
    case *types.FactCollection:
        if v != nil && v.Facts != nil {
            for key, fact := range v.Facts {
                if fact != nil {
                    concreteFacts[key] = fact.Value
                }
            }
        }
    case interfaces.FactsContext:
        // Handle custom facts context
        if machineFacts := v.GetMachineFacts(); machineFacts != nil {
            for machine, factCollection := range machineFacts {
                if factCollection != nil && factCollection.Facts != nil {
                    for key, fact := range factCollection.Facts {
                        if fact != nil {
                            concreteFacts[fmt.Sprintf("%s.%s", machine, key)] = fact.Value
                        }
                    }
                }
            }
        }
        
        // Handle global facts
        if globalFacts := v.GetGlobalFacts(); globalFacts != nil && globalFacts.Facts != nil {
            for key, fact := range globalFacts.Facts {
                if fact != nil {
                    concreteFacts[key] = fact.Value
                }
            }
        }
    default:
        // Try reflection-based conversion
        if converted := c.convertViaReflection(factsContext); converted != nil {
            concreteFacts = converted
        }
    }
    
    return concreteFacts
}

// ConvertVariablesToConcrete converts interface variables to concrete types
func (c *Converter) ConvertVariablesToConcrete(variablesContext interfaces.VariablesContext) map[string]interface{} {
    if variablesContext == nil {
        return make(map[string]interface{})
    }
    
    concreteVariables := make(map[string]interface{})
    
    // Handle different variable context types
    switch v := variablesContext.(type) {
    case map[string]interface{}:
        concreteVariables = v
    case *types.VariableCollection:
        if v != nil && v.Variables != nil {
            for _, variable := range v.Variables {
                if variable != nil {
                    concreteVariables[variable.Name] = variable.Value
                }
            }
        }
    case interfaces.VariablesContext:
        // Handle custom variables context
        if resolvedVars := v.GetResolvedVariables(); resolvedVars != nil {
            concreteVariables = resolvedVars
        }
    default:
        // Try reflection-based conversion
        if converted := c.convertViaReflection(variablesContext); converted != nil {
            concreteVariables = converted
        }
    }
    
    return concreteVariables
}

// convertViaReflection converts interface to map using reflection
func (c *Converter) convertViaReflection(value interface{}) map[string]interface{} {
    if value == nil {
        return nil
    }
    
    v := reflect.ValueOf(value)
    if v.Kind() == reflect.Ptr {
        v = v.Elem()
    }
    
    if v.Kind() == reflect.Map {
        result := make(map[string]interface{})
        for _, key := range v.MapKeys() {
            result[key.String()] = v.MapIndex(key).Interface()
        }
        return result
    }
    
    return nil
}
```

### Step 5: Error Handling Standardization

#### 5.1 Error Handling Package
```go
// internal/errors/standards/standards.go
package standards

import (
    "fmt"
    "strings"
)

// ErrorStandards defines error handling standards
type ErrorStandards struct {
    WrappingPattern string
    ErrorTypes      map[string]ErrorType
    ContextErrors   bool
}

// ErrorType defines an error type
type ErrorType struct {
    Name        string
    Description string
    Code        string
    Severity    Severity
}

// Severity represents error severity
type Severity string

const (
    SeverityCritical Severity = "critical"
    SeverityError    Severity = "error"
    SeverityWarning  Severity = "warning"
    SeverityInfo     Severity = "info"
)

// NewErrorStandards creates new error standards
func NewErrorStandards() *ErrorStandards {
    return &ErrorStandards{
        WrappingPattern: "failed to %s: %w",
        ErrorTypes: map[string]ErrorType{
            "ValidationError": {
                Name:        "ValidationError",
                Description: "Validation failed",
                Code:        "VALIDATION_ERROR",
                Severity:    SeverityError,
            },
            "NotFoundError": {
                Name:        "NotFoundError",
                Description: "Resource not found",
                Code:        "NOT_FOUND",
                Severity:    SeverityError,
            },
            "ConfigurationError": {
                Name:        "ConfigurationError",
                Description: "Configuration error",
                Code:        "CONFIG_ERROR",
                Severity:    SeverityError,
            },
        },
        ContextErrors: true,
    }
}

// ErrorWrapper provides standardized error wrapping
type ErrorWrapper struct {
    standards *ErrorStandards
}

// NewErrorWrapper creates a new error wrapper
func NewErrorWrapper(standards *ErrorStandards) *ErrorWrapper {
    return &ErrorWrapper{
        standards: standards,
    }
}

// Wrap wraps an error with context
func (ew *ErrorWrapper) Wrap(err error, context string) error {
    if err == nil {
        return nil
    }
    
    return fmt.Errorf(ew.standards.WrappingPattern, context, err)
}

// WrapWithContext wraps an error with additional context
func (ew *ErrorWrapper) WrapWithContext(err error, context string, additionalContext map[string]interface{}) error {
    if err == nil {
        return nil
    }
    
    // Add additional context to the error message
    contextParts := []string{context}
    for key, value := range additionalContext {
        contextParts = append(contextParts, fmt.Sprintf("%s=%v", key, value))
    }
    
    fullContext := strings.Join(contextParts, " ")
    return fmt.Errorf(ew.standards.WrappingPattern, fullContext, err)
}
```

## Implementation Strategy

### Phase 1: Audit and Documentation (Week 1)
1. **Create interface inventory** - Scan all packages for interfaces
2. **Generate documentation** - Create comprehensive interface documentation
3. **Identify inconsistencies** - Document all interface issues

### Phase 2: Standardization (Week 2)
1. **Define interface standards** - Create interface design standards
2. **Implement validators** - Create interface validation tools
3. **Update existing interfaces** - Apply standards to existing interfaces

### Phase 3: Implementation Verification (Week 3)
1. **Check implementations** - Verify all interface implementations
2. **Fix implementation issues** - Address any implementation problems
3. **Create type conversion utilities** - Centralize type conversion logic

### Phase 4: Error Handling (Week 4)
1. **Standardize error handling** - Implement consistent error patterns
2. **Update error handling** - Apply standards across codebase
3. **Document error patterns** - Document error handling standards

## Success Criteria

### Functional Requirements
- [ ] All interfaces audited and documented
- [ ] Interface standards defined and implemented
- [ ] All interfaces comply with standards
- [ ] Type conversion utilities centralized
- [ ] Error handling standardized

### Quality Requirements
- [ ] Interface documentation complete and accurate
- [ ] All interface implementations verified
- [ ] Type conversion logic centralized and reusable
- [ ] Error handling patterns consistent
- [ ] Code quality improved

## Dependencies

### Required Dependencies
- All existing packages and interfaces
- Documentation generation tools
- Code analysis tools
- Testing framework

### Optional Dependencies
- Static analysis tools
- Code coverage tools
- Documentation hosting

## Risk Assessment

### High Risk
- **Breaking Changes**: Interface changes may break existing code
- **Implementation Complexity**: Some interfaces may be complex to standardize

### Medium Risk
- **Documentation Maintenance**: Interface documentation needs regular updates
- **Performance Impact**: Type conversion utilities may have performance impact

### Low Risk
- **Tool Integration**: Integration with existing tools may be straightforward
- **Testing**: Testing interface compliance is relatively straightforward

## Next Steps

1. **Start with audit** - Begin with interface inventory and documentation
2. **Define standards** - Create interface design standards
3. **Implement gradually** - Apply standards incrementally to avoid breaking changes
4. **Monitor and adjust** - Monitor implementation and adjust as needed
