# Terminology Compliance Report

## Overview

This report analyzes the `./internal/` directory for compliance with the terminology rules defined in `.cursor/rules/terminology.mdc`. The analysis covers packages, functions, method names, and comments to identify banned terminology and assess alignment with desired terminology.

## Summary

### Current State
- **Banned terminology found**: Extensive use of "execute", "execution", and "executor" variants
- **Desired terminology present**: Some "acting" and "run" terminology already in use
- **Mixed terminology**: Codebase contains both banned and desired terminology patterns

### Key Findings
1. **High usage of banned terms**: 50+ instances of banned terminology across multiple packages
2. **Partial adoption of desired terms**: "Acting" terminology is already established in some areas
3. **Inconsistent patterns**: Same concepts use different terminology in different packages

## Detailed Analysis

### Banned Terminology Found

#### 1. "Execute" Method Names
**Files affected:**
- `internal/cli/manager.go`: `ExecuteCommand()`
- `internal/actions/interfaces.go`: `ExecuteAction()`, `ExecuteActionCollection()`
- `internal/actions/acting/interfaces.go`: `ExecuteAction()`, `ExecuteActionCollection()`
- `internal/ssh/interfaces.go`: `ExecuteCommand()`, `ExecuteScript()`, `ExecuteAction()`, `ExecuteTemplate()`
- `internal/coordinator/manager.go`: `ExecuteAction()`
- `internal/coordinator/actions.go`: `ExecuteAction()`, `ExecuteActionCollection()`

#### 2. "Execution" in Comments and Documentation
**Files affected:**
- `internal/cli/manager.go`: "Execute commands and scripts", "parallel execution"
- `internal/cli/commands/actions.go`: "Execute actions on remote machines", "Create execution context"
- `internal/actions/interfaces.go`: "action management and execution"
- `internal/actions/actor.go`: "sequential execution", "parallel execution"
- `internal/schemas/project.go`: "Validate project execution settings"

#### 3. "Executor" Interface Names
**Files affected:**
- `internal/actions/acting/interfaces.go`: `ActingExecutor`
- `internal/cli/commands/interfaces.go`: `CommandExecutor`
- `internal/cli/factory.go`: `NewExecutor()`

#### 4. "Execution" in Configuration
**Files affected:**
- `internal/project/configuration.go`: `ValidateBeforeExecute`, "execution" HCL tags
- `internal/schemas/project.go`: `ValidateProjectExecution()`
- `internal/dependency/actions.go`: `GetExecutionOrder()`

### Desired Terminology Already Present

#### 1. "Acting" Terminology
**Well-established in:**
- `internal/actions/acting/` package (complete package)
- `internal/actions/interfaces.go`: `ActingSession`, `ActingStatus`
- `internal/types/types.go`: `ActingContext`, `ActingState`, `ActingSession`
- `internal/ssh/acting/` package

#### 2. "Run" Terminology
**Present in:**
- `internal/actions/actor.go`: `Run()` method
- `internal/cli/commands/actions.go`: "Run actions" command
- `internal/ssh/acting/manager.go`: "Action run", "Template run"

#### 3. "Orchestration" Terminology
**Present in:**
- `internal/actions/actor.go`: "Run orchestrates actions"
- `internal/types/actions/orchestration.go`: `OrchestrationResult`

## Package-by-Package Analysis

### 1. `internal/actions/` - HIGH PRIORITY
**Status**: Mixed terminology
- **Banned**: `ExecuteAction()`, `ExecuteActionCollection()` in interfaces
- **Desired**: `Run()` method, "acting" package structure
- **Action needed**: Rename interface methods to `ActAction()`, `ActActionCollection()`

### 2. `internal/cli/` - HIGH PRIORITY
**Status**: Heavy banned terminology usage
- **Banned**: `ExecuteCommand()`, "Execute actions", "execution context"
- **Desired**: "Run actions" command exists
- **Action needed**: Rename methods and update documentation

### 3. `internal/ssh/` - MEDIUM PRIORITY
**Status**: Mixed terminology
- **Banned**: `ExecuteCommand()`, `ExecuteScript()`, `ExecuteAction()`
- **Desired**: `acting/` subpackage with proper terminology
- **Action needed**: Rename interface methods

### 4. `internal/project/` - MEDIUM PRIORITY
**Status**: Configuration-level banned terminology
- **Banned**: `ValidateBeforeExecute`, "execution" HCL tags
- **Action needed**: Rename configuration fields and update HCL schemas

### 5. `internal/schemas/` - MEDIUM PRIORITY
**Status**: Validation function names
- **Banned**: `ValidateProjectExecution()`
- **Action needed**: Rename validation functions

### 6. `internal/dependency/` - LOW PRIORITY
**Status**: Internal method names
- **Banned**: `GetExecutionOrder()`
- **Action needed**: Rename to `GetActingOrder()`

## Immediate Action Plan

### CRITICAL: Remove All Banned Terminology NOW

**All banned terminology must be removed immediately. No phasing, no gradual migration.**

#### 1. Interface Renames (IMMEDIATE)
- `internal/actions/interfaces.go`: `ExecuteAction()` → `ActAction()`, `ExecuteActionCollection()` → `ActActionCollection()`
- `internal/cli/interfaces.go`: `ExecuteCommand()` → `ActCommand()`
- `internal/ssh/interfaces.go`: `ExecuteCommand()` → `ActCommand()`, `ExecuteScript()` → `ActScript()`, `ExecuteAction()` → `ActAction()`, `ExecuteTemplate()` → `ActTemplate()`
- `internal/actions/acting/interfaces.go`: `ExecuteAction()` → `ActAction()`, `ExecuteActionCollection()` → `ActActionCollection()`
- `internal/cli/commands/interfaces.go`: `CommandExecutor` → `CommandActor`, `ExecuteCommand()` → `ActCommand()`

#### 2. Implementation Updates (IMMEDIATE)
- `internal/actions/manager.go`: Update all method implementations
- `internal/cli/manager.go`: Update all method implementations  
- `internal/ssh/manager.go`: Update all method implementations
- `internal/coordinator/manager.go`: Update all method implementations
- `internal/coordinator/actions.go`: Update all method implementations

#### 3. Configuration Updates (IMMEDIATE)
- `internal/project/configuration.go`: `ValidateBeforeExecute` → `ValidateBeforeAct`, "execution" → "acting"
- `internal/schemas/project.go`: `ValidateProjectExecution()` → `ValidateProjectActing()`
- `internal/dependency/actions.go`: `GetExecutionOrder()` → `GetActingOrder()`

#### 4. Documentation and Comments (IMMEDIATE)
- Update ALL comments containing "execute", "execution", "executor"
- Update ALL error messages and user-facing text
- Update ALL test function names

## Implementation Strategy

### 1. All-at-Once Approach
- **No phasing, no gradual migration**
- **All changes must be made immediately**
- **No backward compatibility considerations**
- **Break everything if necessary, fix it all**

### 2. Complete Removal
- **Zero tolerance for banned terminology**
- **Replace every instance immediately**
- **No exceptions, no grandfathering**

### 3. Immediate Testing
- **Fix all compilation errors immediately**
- **Update all tests immediately**
- **No broken builds allowed**

## Risk Assessment

### CRITICAL: Immediate Action Required
- **All banned terminology must be removed NOW**
- **No risk assessment needed - this is mandatory**
- **Break everything if necessary, fix it all immediately**

## Conclusion

The codebase contains **extensive banned terminology** that must be **immediately removed**. 

**NO PHASING. NO GRADUAL MIGRATION. NO BACKWARD COMPATIBILITY.**

**ALL BANNED TERMINOLOGY MUST BE REPLACED IMMEDIATELY.**

The desired terminology foundation is already in place, making this a straightforward replacement task that must be completed immediately.
