---
description: Focused implementation plan for error handling standardization using existing infrastructure
globs: ["internal/**/*.go", "cmd/**/*.go"]
alwaysApply: false
---

# Error Handling Standardization - Focused Implementation Plan

**Generated:** 2025-08-18  
**Issue Type:** Code Quality - Error Handling Standardization  
**Priority:** High  
**Status:** Implementation Planning

## Overview

This focused implementation plan addresses error handling inconsistencies by leveraging the existing comprehensive error handling infrastructure. Rather than building new systems, we'll systematically migrate existing code to use the established error types and patterns.

## Current State Analysis

### Existing Infrastructure ✅
The codebase already has excellent error handling infrastructure:

1. **Common Error Types** (`internal/types/common/common.go`):
   - `ErrorDetails` struct with context, stack traces, and recoverability
   - `ValidationError` with field information
   - Structured error information with timestamps

2. **Domain-Specific Error Types**:
   - SSH errors (`internal/types/ssh/errors.go`)
   - Action errors (`internal/types/actions/errors.go`) 
   - Facts errors (`internal/types/facts/errors.go`)
   - Machine errors (`internal/types/machines/errors.go`)

3. **Error Handling Patterns**:
   - Structured error wrapping
   - Error categorization
   - Context preservation
   - Recovery information

### Current Issues ❌
- **Inconsistent Usage**: Code uses `fmt.Errorf` and `errors.New` instead of existing structured types
- **Missing Context**: Errors lack operation and component information
- **Pattern Inconsistency**: Error handling varies across packages

## Focused Implementation Strategy

### Core Principle: Leverage Existing Infrastructure
**Do NOT create new error handling systems.** Instead:
1. **Use existing error types** from `internal/types/common/common.go`
2. **Utilize domain-specific error types** (SSH, Actions, Facts, Machines)
3. **Build on ErrorDetails** structure for context and stack traces
4. **Follow established patterns** for error wrapping and categorization

### Systematic Migration Approach

#### Phase 1: Foundation (Week 1)
**Goal**: Enhance existing error types with additional context

**Task 1.1: Audit Current Error Usage**
- [ ] Scan all Go files for `fmt.Errorf` and `errors.New` usage
- [ ] Identify patterns of missing context
- [ ] Document current error handling inconsistencies

**Task 1.2: Enhance Error Wrapping Utilities**
- [ ] Create error wrapping utilities that use existing `ErrorDetails`
- [ ] Add operation and component context helpers
- [ ] Implement consistent error categorization

**Deliverable**: Enhanced error wrapping utilities that build on existing infrastructure

#### Phase 2: Package Migration (Week 2-3)
**Goal**: Migrate packages to use existing error types consistently

**Task 2.1: SSH Package Migration**
- [ ] Update `internal/ssh/client.go` to use existing SSH error types
- [ ] Replace `fmt.Errorf` with `spookytypesssh.NewConnectionError`
- [ ] Add operation context to all SSH errors

**Task 2.2: Facts Package Migration**
- [ ] Update `internal/facts/manager.go` to use existing fact error types
- [ ] Replace generic errors with `spookytypesfacts.NewValidationError`
- [ ] Add machine context to fact collection errors

**Task 2.3: Actions Package Migration**
- [ ] Update `internal/actions/manager.go` to use existing action error types
- [ ] Replace inconsistent error wrapping with `spookytypesactions.NewActionError`
- [ ] Add action context to all action errors

**Task 2.4: Config Package Migration**
- [ ] Update `internal/config/loader.go` to use existing config error types
- [ ] Replace generic errors with `spookytypesconfig.NewConfigurationError`
- [ ] Add file path context to configuration errors

**Deliverable**: All packages using existing error types consistently

#### Phase 3: Error Context Enhancement (Week 4)
**Goal**: Add rich context to all errors using existing patterns

**Task 3.1: Operation Context**
- [ ] Add operation names to all error creation
- [ ] Include component information in error context
- [ ] Add timestamp and stack trace information

**Task 3.2: Recovery Information**
- [ ] Add recovery suggestions to error types
- [ ] Implement error categorization (validation, connection, etc.)
- [ ] Add severity levels to errors

**Deliverable**: All errors include rich context and recovery information

#### Phase 4: Validation and Testing (Week 5)
**Goal**: Ensure consistent error handling across the codebase

**Task 4.1: Error Handling Validation**
- [ ] Create linting rules for error handling consistency
- [ ] Validate all errors use existing error types
- [ ] Check for proper error context inclusion

**Task 4.2: Error Handling Tests**
- [ ] Add tests for error type usage
- [ ] Test error context preservation
- [ ] Validate error categorization

**Deliverable**: Comprehensive error handling validation and testing

## Implementation Examples

### Before: Generic Error Handling
```go
// PROBLEMATIC - Generic error handling
func (c *Client) Connect() error {
    if c.config == nil {
        return fmt.Errorf("config is nil") // Generic error
    }
    
    client, err := ssh.Dial("tcp", c.config.Host, c.config.ClientConfig)
    if err != nil {
        return err // Raw error without context
    }
    
    return nil
}
```

### After: Using Existing Error Types
```go
// IMPROVED - Using existing SSH error types
func (c *Client) Connect() error {
    if c.config == nil {
        return spookytypesssh.NewConnectionError(
            "SSH_CLIENT_NIL_CONFIG",
            "SSH client configuration is nil",
            "ssh_client_connect",
            map[string]interface{}{
                "component": "ssh_client",
                "operation": "connect",
            },
        )
    }
    
    client, err := ssh.Dial("tcp", c.config.Host, c.config.ClientConfig)
    if err != nil {
        return spookytypesssh.WrapWithContext(
            err,
            spookytypesssh.ErrorTypeConnection,
            "SSH_CLIENT_CONNECTION_FAILED",
            "SSH connection failed",
            "ssh_client_connect",
            map[string]interface{}{
                "component": "ssh_client",
                "operation": "connect",
                "host":      c.config.Host,
                "port":      c.config.Port,
            },
        )
    }
    
    return nil
}
```

### Before: Missing Context
```go
// PROBLEMATIC - Missing error context
func (m *Manager) CollectFacts(server string) (*spookytypes.FactCollection, error) {
    if server == "" {
        return nil, fmt.Errorf("server is required")
    }
    
    facts, err := m.storage.GetFacts(server)
    if err != nil {
        return nil, err // Missing context about what failed
    }
    
    return facts, nil
}
```

### After: Rich Context Using Existing Types
```go
// IMPROVED - Rich context using existing fact error types
func (m *Manager) CollectFacts(ctx context.Context, machine interface{}) (*spookytypes.FactCollection, error) {
    machineObj, ok := machine.(*spookytypes.Machine)
    if !ok {
        return nil, spookytypesfacts.NewValidationError(
            "FACTS_COLLECTION_INVALID_MACHINE_TYPE",
            "Machine must be of type *spookytypes.Machine",
            "facts_collect",
            map[string]interface{}{
                "component": "facts_manager",
                "operation": "collect_facts",
                "field":     "machine",
                "type":      fmt.Sprintf("%T", machine),
            },
        )
    }

    if machineObj.Hostname == "" {
        return nil, spookytypesfacts.NewValidationError(
            "FACTS_COLLECTION_MISSING_HOSTNAME",
            "Machine hostname is required for fact collection",
            "facts_collect",
            map[string]interface{}{
                "component": "facts_manager",
                "operation": "collect_facts",
                "field":     "hostname",
            },
        )
    }

    facts, err := m.collector.Collect(ctx, machineObj)
    if err != nil {
        return nil, spookytypesfacts.WrapWithContext(
            err,
            "FACTS_COLLECTION_FAILED",
            "Failed to collect facts from machine",
            "facts_collect",
            map[string]interface{}{
                "component": "facts_manager",
                "operation": "collect_facts",
                "machine":   machineObj.Hostname,
                "host":      machineObj.Host,
            },
        )
    }

    return facts, nil
}
```

## Error Wrapping Utilities

### Enhanced Error Wrapping (Building on Existing Infrastructure)
```go
// Enhanced error wrapping utilities that use existing ErrorDetails
func WrapWithContext(err error, category ErrorCategory, code, message, operation string, context map[string]interface{}) *EnhancedErrorDetails {
    // Use existing ErrorDetails as foundation
    baseError := &spookytypescommon.ErrorDetails{
        Code:        code,
        Message:     message,
        Context:     context,
        Timestamp:   time.Now(),
        Recoverable: true,
    }
    
    // Enhance with additional context
    enhancedError := &EnhancedErrorDetails{
        ErrorDetails: *baseError,
        Category:     category,
        Severity:     getSeverityForCategory(category),
        Operation:    operation,
        Component:    getComponentFromOperation(operation),
        Suggestions:  getSuggestionsForCategory(category),
    }
    
    // Preserve original error as cause
    if err != nil {
        enhancedError.Context["original_error"] = err.Error()
    }
    
    return enhancedError
}

// Helper functions for consistent error creation
func NewValidationError(code, message, operation string, context map[string]interface{}) *EnhancedErrorDetails {
    return WrapWithContext(nil, ErrorCategoryValidation, code, message, operation, context)
}

func NewConnectionError(code, message, operation string, context map[string]interface{}) *EnhancedErrorDetails {
    return WrapWithContext(nil, ErrorCategoryConnection, code, message, operation, context)
}

func NewAuthenticationError(code, message, operation string, context map[string]interface{}) *EnhancedErrorDetails {
    return WrapWithContext(nil, ErrorCategoryAuthentication, code, message, operation, context)
}
```

## Success Metrics

### Error Handling Consistency
- **Target**: 100% of packages use existing error types
- **Measurement**: Code review and static analysis
- **Baseline**: Current inconsistent usage
- **Monitoring**: Regular code reviews

### Error Context Quality
- **Target**: 100% of errors include operation and component context
- **Measurement**: Error message analysis
- **Baseline**: Current missing context
- **Monitoring**: Error log analysis

### Error Categorization
- **Target**: 100% of errors properly categorized
- **Measurement**: Error type analysis
- **Baseline**: Current uncategorized errors
- **Monitoring**: Error categorization metrics

## Risk Assessment

### Technical Risks
- **Breaking Changes**: Risk of breaking existing error handling
- **Mitigation**: Maintain backward compatibility with existing error types
- **Performance Impact**: Risk of performance impact from enhanced error handling
- **Mitigation**: Optimize error creation and context gathering

### Functional Risks
- **Error Information Loss**: Risk of losing error information during migration
- **Mitigation**: Preserve all existing error information in enhanced types
- **Interface Changes**: Risk of breaking interface contracts
- **Mitigation**: Maintain interface compatibility

## Rollback Plan

### Rollback Triggers
- Error handling performance degradation > 10%
- Error information loss or corruption
- Breaking changes to existing functionality

### Rollback Procedure
1. **Immediate**: Revert to previous error handling patterns
2. **Short-term**: Disable enhanced error handling features
3. **Long-term**: Fix issues and re-enable enhanced error handling

## Dependencies

### Internal Dependencies
- `internal/types/common/common.go` - Existing error types
- `internal/types/*/errors.go` - Domain-specific error types
- `internal/logging/` - Logging infrastructure

### External Dependencies
- No additional external dependencies required
- Uses existing Go standard library error handling

## Timeline

### Week 1: Foundation
- Day 1-2: Audit current error usage
- Day 3-4: Enhance error wrapping utilities
- Day 5: Create error context helpers

### Week 2-3: Package Migration
- Week 2: SSH and Facts package migration
- Week 3: Actions and Config package migration

### Week 4: Error Context Enhancement
- Day 1-3: Add operation and component context
- Day 4-5: Add recovery information and categorization

### Week 5: Validation and Testing
- Day 1-2: Create error handling validation
- Day 3-4: Add comprehensive error handling tests
- Day 5: Final validation and documentation

## Conclusion

This focused implementation plan provides a systematic approach to standardizing error handling across the spooky codebase by leveraging the existing comprehensive error handling infrastructure. The plan focuses on consistency, context, and recovery while maintaining backward compatibility and following established patterns.

**Expected Outcomes:**
- 100% consistent error handling across all packages using existing error types
- Enhanced error context for better debugging
- Proper error categorization and recovery information
- Improved error handling test coverage
- Clear error handling standards and documentation

The implementation will result in a more robust, maintainable, and user-friendly error handling system that builds upon the existing solid foundation rather than replacing it.

**Priority Actions:**
1. **Immediate**: Audit and enhance existing error infrastructure
2. **High**: Migrate packages to use existing error types consistently
3. **Medium**: Add error context and recovery information
4. **Long-term**: Establish comprehensive error handling validation

The plan ensures a smooth transition to standardized error handling while maintaining all existing functionality and adding significant improvements in error context, categorization, and debugging capabilities.
