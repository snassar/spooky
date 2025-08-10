# Execute to Act Terminology Migration Plan

## Overview

This plan outlines the systematic migration from "execute" terminology to "act" terminology throughout the Spooky codebase, as per user preference for consistency.

## Migration Rules

- `Execute` → `Run`
- `execute` → `run`
- `Execution` → `Acting`
- `execution` → `acting`
- `exec` → `run` (when referring to action execution, not os/exec)
- `ValidateBeforeExecute` → `ValidateBeforeRunning`
- `validate_before_execute` → `validate_before_running`

## File-by-File Migration Tasks

### 1. Main Application Files

#### `main.go`
- **Task**: Update documentation strings
- **Changes**:
  - Line 25: "Execute commands and scripts" → "Act commands and scripts"
  - Line 27: "Support for parallel execution" → "Support for parallel acting"
  - Line 67: "Application execution failed" → "Application acting failed"

#### `main_simple_test.go`
- **Task**: Update test comments
- **Changes**:
  - Line 12: "actual execution" → "actual acting"

### 2. CLI Layer

#### `internal/cli/manager.go`
- **Task**: Update method names and documentation
- **Changes**:
  - Line 58: "Execute commands and scripts" → "Act commands and scripts"
  - Line 60: "parallel execution" → "parallel acting"
  - Line 75: Comment "ExecuteCommand executes" → "ActCommand acts"
  - Line 76: `ExecuteCommand` → `ActCommand`
  - Line 82: `rootCommand.Execute()` → `rootCommand.Act()` (if applicable)

#### `internal/cli/commands/actions.go`
- **Task**: Update command descriptions and method calls
- **Changes**:
  - Line 19: "Manage and execute actions" → "Manage and act actions"
  - Line 20: "execute actions on remote machines" → "act actions on remote machines"
  - Line 93: "Execute actions on remote machines" → "Act actions on remote machines"
  - Line 129: Comment "Execute action" → "Act action"
  - Line 143: Comment "Create execution context" → "Create acting context"
  - Line 144: `execContext` → `actContext`
  - Line 154: "Would execute action" → "Would act action"
  - Line 158: "Executing action" → "Acting action"
  - Line 159: `ExecuteAction` → `ActAction`
  - Line 161: Comment "Execute all actions" → "Act all actions"
  - Line 163: "No actions found to execute" → "No actions found to act"
  - Line 168: "Would execute %d actions" → "Would act %d actions"
  - Line 172: "Executing %d actions" → "Acting %d actions"
  - Line 174: `execContext` → `actContext`
  - Line 183: `ExecuteAction` → `ActAction`
  - Line 184: "Failed to execute action" → "Failed to act action"
  - Line 194: "Action name to execute" → "Action name to act"
  - Line 198: "Number of parallel executions" → "Number of parallel actings"
  - Line 246: Comment "Create execution context" → "Create acting context"
  - Line 247: `execContext` → `actContext`

#### `internal/cli/interfaces.go`
- **Task**: Update interface method names
- **Changes**:
  - Line 12: `ExecuteCommand` → `ActCommand`

#### `internal/cli/commands/interfaces.go`
- **Task**: Update interface method names
- **Changes**:
  - Line 43: `ExecuteCommand` → `ActCommand`

### 3. Actions Layer

#### `internal/actions/interfaces.go`
- **Task**: Update interface documentation and method names
- **Changes**:
  - Line 9: "action management and execution" → "action management and acting"
  - Line 19: `ExecuteAction` → `ActAction`
  - Line 20: `ExecuteActionCollection` → `ActActionCollection`

#### `internal/actions/actor.go`
- **Task**: Update comments and TODO items
- **Changes**:
  - Line 145: "sequential execution" → "sequential acting"
  - Line 151: "parallel execution" → "parallel acting"
  - Line 157: "single action execution" → "single action acting"
  - Line 315: "parallel execution setting" → "parallel acting setting"

#### `internal/actions/manager.go`
- **Task**: Update method names and comments
- **Changes**:
  - Line 139: Comment "ExecuteAction executes" → "ActAction acts"
  - Line 140: `ExecuteAction` → `ActAction`
  - Line 182: `executeCommandAction` → `actCommandAction`
  - Line 184: `executeScriptAction` → `actScriptAction`
  - Line 186: `executeTemplateAction` → `actTemplateAction`
  - Line 207: Comment "ExecuteActionCollection executes" → "ActActionCollection acts"
  - Line 208: `ExecuteActionCollection` → `ActActionCollection`
  - Line 234: Comment "executeCommandAction executes" → "actCommandAction acts"
  - Line 235: `executeCommandAction` → `actCommandAction`
  - Line 263: Comment "executeScriptAction executes" → "actScriptAction acts"
  - Line 264: `executeScriptAction` → `actScriptAction`
  - Line 292: Comment "executeTemplateAction executes" → "actTemplateAction acts"
  - Line 293: `executeTemplateAction` → `actTemplateAction`

#### `internal/actions/acting/actor.go`
- **Task**: Update method names and comments
- **Changes**:
  - Line 32: Comment "Execute executes" → "Act acts"
  - Line 33: `Execute` → `Act`
  - Line 38: "Executing actor" → "Acting actor"
  - Line 52: Comment "Execute based on" → "Act based on"
  - Line 56: `executeCommand` → `actCommand`
  - Line 58: `executeScript` → `actScript`
  - Line 60: `executeTemplateDeploy` → `actTemplateDeploy`
  - Line 62: `executeTemplateEvaluate` → `actTemplateEvaluate`
  - Line 64: `executeTemplateValidate` → `actTemplateValidate`
  - Line 66: `executeTemplateCleanup` → `actTemplateCleanup`
  - Line 78: Comment "Update result based on execution" → "Update result based on acting"
  - Line 93: "Actor execution completed" → "Actor acting completed"
  - Line 100: Comment "Prepare prepares the actor for execution" → "Prepare prepares the actor for acting"
  - Line 127: Comment "Cancel cancels the actor execution" → "Cancel cancels the actor acting"
  - Line 226: Comment "executeCommand executes" → "actCommand acts"
  - Line 227: `executeCommand` → `actCommand`
  - Line 228: "Executing command action" → "Acting command action"
  - Line 232: "actual command execution" → "actual command acting"
  - Line 235: "Executing command on each machine" → "Acting command on each machine"
  - Line 239: "Command executed successfully" → "Command acted successfully"
  - Line 245: Comment "executeScript executes" → "actScript acts"
  - Line 246: `executeScript` → `actScript`
  - Line 247: "Executing script action" → "Acting script action"
  - Line 250: "actual script execution" → "actual script acting"
  - Line 254: "Executing script on each machine" → "Acting script on each machine"
  - Line 258: "Script executed successfully" → "Script acted successfully"
  - Line 264: Comment "executeTemplateDeploy executes" → "actTemplateDeploy acts"
  - Line 265: `executeTemplateDeploy` → `actTemplateDeploy`
  - Line 266: "Executing template deploy action" → "Acting template deploy action"
  - Line 284: Comment "executeTemplateEvaluate executes" → "actTemplateEvaluate acts"
  - Line 285: `executeTemplateEvaluate` → `actTemplateEvaluate`
  - Line 286: "Executing template evaluate action" → "Acting template evaluate action"
  - Line 302: Comment "executeTemplateValidate executes" → "actTemplateValidate acts"
  - Line 303: `executeTemplateValidate` → `actTemplateValidate`

#### `internal/actions/acting/manager.go`
- **Task**: Update method names and comments
- **Changes**:
  - Line 62: Comment "ExecuteAction executes" → "ActAction acts"
  - Line 63: `ExecuteAction` → `ActAction`
  - Line 124: Comment "ExecuteActionCollection executes" → "ActActionCollection acts"
  - Line 125: `ExecuteActionCollection` → `ActActionCollection`
  - Line 170: `executeActionWithRetry` → `actActionWithRetry`
  - Line 195: `executeActionWithRetry` → `actActionWithRetry`
  - Line 358: Comment "executeActionWithRetry executes" → "actActionWithRetry acts"
  - Line 359: `executeActionWithRetry` → `actActionWithRetry`

#### `internal/actions/acting/interfaces.go`
- **Task**: Update interface method names
- **Changes**:
  - Line 12: `ExecuteAction` → `ActAction`
  - Line 13: `ExecuteActionCollection` → `ActActionCollection`
  - Line 51: `ExecuteCommand` → `ActCommand`
  - Line 52: `ExecuteScript` → `ActScript`
  - Line 53: `ExecuteTemplate` → `ActTemplate`
  - Line 83: `ExecuteCommand` → `ActCommand`
  - Line 84: `ExecuteScript` → `ActScript`

### 4. SSH Layer

#### `internal/ssh/manager.go`
- **Task**: Update method names and comments
- **Changes**:
  - Line 52: Comment "ExecuteCommand executes" → "ActCommand acts"
  - Line 53: `ExecuteCommand` → `ActCommand`
  - Line 54: `ExecuteCommand` → `ActCommand`
  - Line 57: Comment "ExecuteScript executes" → "ActScript acts"
  - Line 58: `ExecuteScript` → `ActScript`
  - Line 59: `ExecuteScript` → `ActScript`
  - Line 92: Comment "ExecuteAction executes" → "ActAction acts"
  - Line 93: `ExecuteAction` → `ActAction`
  - Line 94: `ExecuteAction` → `ActAction`
  - Line 97: Comment "ExecuteTemplate executes" → "ActTemplate acts"
  - Line 98: `ExecuteTemplate` → `ActTemplate`
  - Line 99: `ExecuteTemplate` → `ActTemplate`
  - Line 167: `ExecuteActionOnMachine` → `ActActionOnMachine`
  - Line 176: `ExecuteAction` → `ActAction`

#### `internal/ssh/interfaces.go`
- **Task**: Update interface method names
- **Changes**:
  - Line 12: `ExecuteCommand` → `ActCommand`
  - Line 13: `ExecuteScript` → `ActScript`
  - Line 26: `ExecuteAction` → `ActAction`
  - Line 27: `ExecuteTemplate` → `ActTemplate`
  - Line 43: `ExecuteCommand` → `ActCommand`
  - Line 44: `ExecuteScript` → `ActScript`
  - Line 66: `ExecuteAction` → `ActAction`
  - Line 67: `ExecuteTemplate` → `ActTemplate`

#### `internal/ssh/client/manager.go`
- **Task**: Update method names, comments, and error messages
- **Changes**:
  - Line 14: `executionManager` → `actingManager`
  - Line 23: `executionManager` → `actingManager`
  - Line 30: `executionManager:` → `actingManager:`
  - Line 64: Comment "ExecuteCommand executes" → "ActCommand acts"
  - Line 65: `ExecuteCommand` → `ActCommand`
  - Line 67: `validateExecutionParams` → `validateActingParams`
  - Line 68: "execution validation failed" → "acting validation failed"
  - Line 72: `executionManager.ExecuteCommand` → `actingManager.ActCommand`
  - Line 74: "command execution failed" → "command acting failed"
  - Line 85: Comment "ExecuteScript executes" → "ActScript acts"
  - Line 86: `ExecuteScript` → `ActScript`
  - Line 88: `validateExecutionParams` → `validateActingParams`
  - Line 89: "execution validation failed" → "acting validation failed"
  - Line 93: `executionManager.ExecuteScript` → `actingManager.ActScript`
  - Line 95: "script execution failed" → "script acting failed"
  - Line 200: `validateExecutionParams` → `validateActingParams`

#### `internal/ssh/client/interfaces.go`
- **Task**: Update interface names and method names
- **Changes**:
  - Line 12: `ExecuteCommand` → `ActCommand`
  - Line 13: `ExecuteScript` → `ActScript`
  - Line 35: Comment "ExecutionManager defines" → "ActingManager defines"
  - Line 36: `ExecutionManager` → `ActingManager`
  - Line 37: `ExecuteCommand` → `ActCommand`
  - Line 38: `ExecuteScript` → `ActScript`

### 5. Coordinator Layer

#### `internal/coordinator/manager.go`
- **Task**: Update method names and comments
- **Changes**:
  - Line 227: Comment "ExecuteAction executes" → "ActAction acts"
  - Line 228: `ExecuteAction` → `ActAction`
  - Line 229: `ExecuteAction` → `ActAction`

#### `internal/coordinator/actions.go`
- **Task**: Update method names and comments
- **Changes**:
  - Line 157: Comment "ExecuteAction executes" → "ActAction acts"
  - Line 158: `ExecuteAction` → `ActAction`
  - Line 199: `ExecuteActionCollection` → `ActActionCollection`

#### `internal/coordinator/actions_test.go`
- **Task**: Update test function names
- **Changes**:
  - Line 75: `TestExecuteAction` → `TestActAction`
  - Line 91: `ExecuteAction` → `ActAction`

#### `internal/coordinator/machines.go`
- **Task**: Update method calls
- **Changes**:
  - Line 162: `ExecuteCommand` → `ActCommand`
  - Line 296: `ExecuteCommand` → `ActCommand`

### 6. Facts Collection

#### `internal/facts/collectors/ssh/collector.go`
- **Task**: Update method calls and comments
- **Changes**:
  - Line 92: `ExecuteCommand` → `ActCommand`
  - Line 99: `ExecuteCommand` → `ActCommand`
  - Line 106: `ExecuteCommand` → `ActCommand`
  - Line 135: Comment "executeCommand executes" → "actCommand acts"

### 7. Project Configuration

#### `internal/project/configuration.go`
- **Task**: Update field names and HCL tags
- **Changes**:
  - Line 46: Comment "Execution settings" → "Acting settings"
  - Line 50: `ValidateBeforeExecute` → `ValidateBeforeAct`
  - Line 199: `"execution"` → `"acting"`
  - Line 203: `"validate_before_execute"` → `"validate_before_act"`
  - Line 292: Comment "Execution settings" → "Acting settings"
  - Line 344: `ValidateBeforeExecute` → `ValidateBeforeAct`
  - Line 436: Comment "Apply execution settings" → "Apply acting settings"
  - Line 437: `project.Execution` → `project.Acting`
  - Line 438: `project.Execution.DefaultTimeout` → `project.Acting.DefaultTimeout`
  - Line 439: `resolved.DefaultTimeout` → `resolved.DefaultTimeout`
  - Line 441: `project.Execution.MaxParallel` → `project.Acting.MaxParallel`
  - Line 442: `resolved.MaxParallel` → `resolved.MaxParallel`
  - Line 444: `resolved.DryRunDefault` → `resolved.DryRunDefault`
  - Line 445: `project.Execution.ValidateBeforeExecute` → `project.Acting.ValidateBeforeAct`
  - Line 446: `resolved.BackupBeforeChanges` → `resolved.BackupBeforeChanges`
  - Line 504: Comment "Execution settings" → "Acting settings"
  - Line 548: Comment "Validate execution settings" → "Validate acting settings"

#### `internal/project/manager.go`
- **Task**: Update field references and comments
- **Changes**:
  - Line 134: Comment "Convert execution configuration" → "Convert acting configuration"
  - Line 135: `project.Execution` → `project.Acting`
  - Line 136: `project.Execution` → `project.Acting`
  - Line 139: `project.Execution.DefaultTimeout` → `project.Acting.DefaultTimeout`
  - Line 140: `project.Execution.DryRunDefault` → `project.Acting.DryRunDefault`
  - Line 141: `project.Execution.ValidateBeforeExecute` → `project.Acting.ValidateBeforeAct`
  - Line 142: `project.Execution.BackupBeforeChanges` → `project.Acting.BackupBeforeChanges`
  - Line 146: `project.Execution.MaxParallel` → `project.Acting.MaxParallel`
  - Line 148: `project.Execution.MaxParallel` → `project.Acting.MaxParallel`
  - Line 165: Comment "detailed Execution settings" → "detailed Acting settings"
  - Line 387: `updates.Execution` → `updates.Acting`
  - Line 388: `existingProject.Execution` → `existingProject.Acting`

### 8. Schemas and Validation

#### `internal/schemas/project.go`
- **Task**: Update validation function names and comments
- **Changes**:
  - Line 121: Comment "Validate project execution settings" → "Validate project acting settings"
  - Line 122: `ValidateProjectExecution` → `ValidateProjectActing`
  - Line 125: `"execution"` → `"acting"`
  - Line 127: `fmt.Sprintf("%v", proj.Execution)` → `fmt.Sprintf("%v", proj.Acting)`
  - Line 345: Comment "ValidateProjectExecution validates" → "ValidateProjectActing validates"
  - Line 346: `ValidateProjectExecution` → `ValidateProjectActing`
  - Line 347: `execution` → `acting`
  - Line 348: Comment "Execution settings are optional" → "Acting settings are optional"
  - Line 352: `execution.DefaultTimeout` → `acting.DefaultTimeout`
  - Line 357: `execution.MaxParallel` → `acting.MaxParallel`

#### `internal/schemas/errors.go`
- **Task**: Update error messages
- **Changes**:
  - Line 72: `"action_exec"` → `"action_act"`
  - Line 79: No change needed (scriptfile is correct)

#### `internal/schemas/evolution.go`
- **Task**: Update method names and comments
- **Changes**:
  - Line 464: Comment "ExecuteSchemaEvolutionWorkflow executes" → "ActSchemaEvolutionWorkflow acts"
  - Line 465: `ExecuteSchemaEvolutionWorkflow` → `ActSchemaEvolutionWorkflow`
  - Line 477: `executedSteps` → `actedSteps`
  - Line 479: Comment "Execute each step" → "Act each step"
  - Line 484: Comment "Execute step based on" → "Act step based on"
  - Line 492: `executedSteps` → `actedSteps`
  - Line 496: `"executed_steps"` → `"acted_steps"`
  - Line 511: `executedSteps` → `actedSteps`
  - Line 515: `"executed_steps"` → `"acted_steps"`
  - Line 525: `executedSteps` → `actedSteps`
  - Line 529: `"executed_steps"` → `"acted_steps"`
  - Line 540: `executedSteps` → `actedSteps`
  - Line 545: `"executed_steps"` → `"acted_steps"`

### 9. Templates

#### `internal/templates/engine/manager.go`
- **Task**: Update method names and comments
- **Changes**:
  - Line 87: Comment "SetMaxExecutionTime sets" → "SetMaxActingTime sets"
  - Line 88: `SetMaxExecutionTime` → `SetMaxActingTime`
  - Line 92: `MaxExecutionTime` → `MaxActingTime`

#### `internal/templates/engine/interfaces.go`
- **Task**: Update interface method names
- **Changes**:
  - Line 16: `SetMaxExecutionTime` → `SetMaxActingTime`

#### `internal/templates/factory.go`
- **Task**: Update field names
- **Changes**:
  - Line 39: `MaxExecutionTime:` → `MaxActingTime:`

#### `internal/templates/engine/renderer.go`
- **Task**: Update method calls
- **Changes**:
  - Line 21: `tmpl.Execute` → `tmpl.Act` (if applicable, otherwise keep as is)

### 10. Dependency Management

#### `internal/dependency/actions.go`
- **Task**: Update method names and comments
- **Changes**:
  - Line 127: Comment "GetExecutionOrder returns" → "GetActingOrder returns"
  - Line 128: `GetExecutionOrder` → `GetActingOrder`
  - Line 138: "failed to resolve execution order" → "failed to resolve acting order"
  - Line 144: Comment "GetParallelGroups returns groups" → "GetParallelGroups returns groups"
  - Line 146: Comment "Get execution order first" → "Get acting order first"
  - Line 147: `GetExecutionOrder` → `GetActingOrder`

#### `internal/dependency/actions_test.go`
- **Task**: Update test function names and method calls
- **Changes**:
  - Line 168: `TestGetExecutionOrder` → `TestGetActingOrder`
  - Line 191: `GetExecutionOrder` → `GetActingOrder`
  - Line 193: "Failed to get execution order" → "Failed to get acting order"
  - Line 226: `TestGetExecutionOrderWithCircularDeps` → `TestGetActingOrderWithCircularDeps`
  - Line 244: `GetExecutionOrder` → `GetActingOrder`

### 11. Interfaces

#### `internal/interfaces/actions.go`
- **Task**: Update method names and comments
- **Changes**:
  - Line 9: Comment "ValidateAction validates" → "ValidateAction validates"
  - Line 10: `ValidateAction` → `ValidateAction` (keep as is)
  - Line 12: Comment "PrepareActionForExecution" → "PrepareActionForActing"
  - Line 13: `PrepareActionForExecution` → `PrepareActionForActing`
  - Line 16: `ExecuteAction` → `ActAction`

#### `internal/interfaces/manager.go`
- **Task**: Update method names and comments
- **Changes**:
  - Line 15: `LoadContextForAction` → `LoadContextForAction` (keep as is)
  - Line 16: `ValidateActionWithAllSystems` → `ValidateActionWithAllSystems` (keep as is)
  - Line 17: `PrepareActionForExecution` → `PrepareActionForActing`
  - Line 18: `ExecuteAction` → `ActAction`

#### `internal/interfaces/contexts.go`
- **Task**: Update struct documentation
- **Changes**:
  - Line 50: Comment "ActionExecutionContext provides" → "ActionActingContext provides"
  - Line 51: `ActionExecutionContext` → `ActionActingContext`

### 12. Types

#### `internal/types/types.go`
- **Task**: Update type aliases
- **Changes**:
  - Line 221: `ExecuteConfig` → `ActConfig`
  - Line 226: `ExecutionError` → `ActingError`
  - Line 308: `ExecutionContext` → `ActingContext`

### 13. Tools

#### `tools/generate-test-project/main.go`
- **Task**: Update method calls
- **Changes**:
  - Line 683: `rootCmd.Execute()` → `rootCmd.Act()` (if applicable)

#### `tools/pre-commit/main.go`
- **Task**: Update error messages
- **Changes**:
  - Line 104: "Test execution failed" → "Test acting failed"

#### `tools/spooky-generate/main.go`
- **Task**: Update method calls
- **Changes**:
  - Line 35: `rootCmd.Execute()` → `rootCmd.Act()` (if applicable)

### 14. Machine Connectivity

#### `internal/machines/connectivity/manager.go`
- **Task**: Update comments and error messages
- **Changes**:
  - Line 238: Comment "Execute a simple command" → "Act a simple command"
  - Line 245: Comment "Execute the test command" → "Act the test command"
  - Line 248: "SSH command execution failed" → "SSH command acting failed"

## Implementation Order

1. **Start with interfaces** - Update interface definitions first
2. **Update type definitions** - Change struct fields and type aliases
3. **Update implementation files** - Modify concrete implementations
4. **Update tests** - Fix test function names and method calls
5. **Update documentation** - Fix comments and error messages
6. **Update configuration** - Change HCL tags and field names

## Testing Strategy

1. Run `go build` to check for compilation errors
2. Run `go test ./...` to ensure tests pass
3. Run `golangci-lint run` to check for any remaining issues
4. Test CLI commands to ensure functionality works
5. Test project validation to ensure configuration changes work

## Notes

- Keep `os/exec` package references unchanged
- Keep `tmpl.Execute()` calls unchanged (these are Go template execution)
- Keep `rootCmd.Execute()` calls unchanged (these are Cobra command execution)
- Focus on action/execution terminology, not general execution concepts

## Risk Assessment

- **High Impact**: Interface changes will require updating all implementations
- **Medium Risk**: Configuration changes may break existing project files
- **Low Risk**: Comment and documentation changes are safe

## Rollback Plan

- Keep git history with clear commit messages
- Test thoroughly before merging
- Consider feature flag for gradual migration if needed
