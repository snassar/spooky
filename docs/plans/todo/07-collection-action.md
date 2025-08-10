# Implementation Plan: Action Orchestration

## Overview
Implement action orchestration functionality that can run multiple actions with internal planning, dependency resolution, parallel orchestration where possible, rollback on failure, and comprehensive result aggregation.

## Task Details
- **Task ID**: 3.1
- **Priority**: High
- **File**: `internal/actions/manager.go`
- **Functions**: `Run`

## Current State Analysis

### Existing Patterns
1. **Action Structure**: Actions use `types.Action` with dependencies, priority, and run settings
2. **Context System**: `types.ActionContext` provides run context with facts, variables, machines
3. **Result Structure**: `types.RunResult` captures run results with output, error, exit code
4. **State Management**: Actions track state via `types.ActionState`
5. **Error Handling**: Consistent error wrapping with context
6. **Dependency Management**: Actions have dependency arrays for run order

### Existing Implementation Examples
- **Action Manager**: `internal/actions/manager.go` provides action management
- **Action Types**: `internal/types/actions/action.go` defines action structure with dependencies
- **Context System**: `internal/types/actions/context.go` provides run context

## Implementation Requirements

### Interface Compliance
The action orchestration system must:
1. **Run actions** with internal planning, dependency resolution, and optimization
3. **Handle failures** and implement rollback strategies
4. **Support parallel orchestration** where dependencies allow
5. **Implement timeout** handling for the entire collection
6. **Support partial orchestration** and failure recovery
7. **Track orchestration state** and progress
8. **Handle circular dependencies** and validation

### Required Dependencies
- Action manager for individual action running
- Dependency resolution system
- Rollback and recovery mechanisms
- Result aggregation system
- State tracking and progress monitoring

## Detailed Implementation Plan

### Step 1: Analyze Current Action Orchestration Method

**Current Implementation** (to be replaced):
```go
func (m *Manager) Run(ctx context.Context, actions []*types.Action, context *types.ActionContext) (*types.OrchestrationResult, error) {
    // TODO: Implement actual action orchestration with internal planning
    return &types.OrchestrationResult{
        Status: types.OrchestrationStatusCompleted,
        Output: "Actions orchestrated successfully (placeholder)",
    }, nil
}
```

### Step 2: Implement Internal Planning

#### 2.1 Internal Planning Logic
```go
// planOrchestration analyzes dependencies and creates an optimized orchestration plan
func (m *Manager) planOrchestration(actions []*types.Action, context *types.ActionContext) (*types.OrchestrationPlan, error) {
    m.logger.Debug("Planning action orchestration",
        logging.Int("action_count", len(actions)))

    // Build dependency graph
    graph, err := m.buildDependencyGraph(actions)
    if err != nil {
        return nil, fmt.Errorf("failed to build dependency graph: %w", err)
    }

    // Determine orchestration order
    orchestrationOrder, err := graph.topologicalSort()
    if err != nil {
        return nil, fmt.Errorf("failed to determine orchestration order: %w", err)
    }

    // Optimize orchestration strategy
    optimizedPlan, err := m.optimizeOrchestration(actions, orchestrationOrder, graph, context)
    if err != nil {
        return nil, fmt.Errorf("failed to optimize orchestration: %w", err)
    }

    return optimizedPlan, nil
}
```

#### 2.2 Orchestration Plan Structure
```go
// OrchestrationPlan represents the complete orchestration strategy
type OrchestrationPlan struct {
    Actions            []*types.Action
    OrchestrationOrder []string
    ParallelGroups     [][]string      // Groups that can run in parallel
    SequentialOrder    []string        // Actions that must run sequentially
    ConcurrencyLevel   int             // Optimal concurrency for parallel groups
    EstimatedDuration  time.Duration
    ResourceRequirements map[string]interface{}
    RollbackStrategy   *RollbackStrategy
}

// RollbackStrategy defines how to handle failures
type RollbackStrategy struct {
    EnableRollback     bool
    RollbackOrder      []string        // Reverse orchestration order for rollback
    CriticalActions    []string        // Actions that trigger rollback on failure
}
```

### Step 3: Implement Dependency Resolution

#### 3.1 Dependency Graph Construction
```go
// buildDependencyGraph builds a dependency graph for actions
func (m *Manager) buildDependencyGraph(actions []*types.Action) (*dependencyGraph, error) {
    graph := &dependencyGraph{
        nodes: make(map[string]*actionNode),
        edges: make(map[string][]string),
    }
    
    // Create nodes for all actions
    for _, action := range actions {
        graph.nodes[action.Name] = &actionNode{
            action: action,
            state:  types.ActionStatusPending,
        }
    }
    
    // Build edges based on dependencies
    for _, action := range actions {
        for _, dep := range action.Dependencies {
            if _, exists := graph.nodes[dep]; !exists {
                return nil, fmt.Errorf("dependency '%s' not found for action '%s'", dep, action.Name)
            }
            graph.edges[action.Name] = append(graph.edges[action.Name], dep)
        }
    }
    
    // Check for circular dependencies
    if err := graph.detectCircularDependencies(); err != nil {
        return nil, fmt.Errorf("circular dependency detected: %w", err)
    }
    
    return graph, nil
}
```

#### 3.2 Dependency Graph Structure
```go
// dependencyGraph represents the dependency graph for actions
type dependencyGraph struct {
    nodes map[string]*actionNode
    edges map[string][]string // action -> dependencies
}

// actionNode represents a node in the dependency graph
type actionNode struct {
    action *types.Action
    state  types.ActionStatus
    result *types.RunResult
}

// detectCircularDependencies detects circular dependencies using DFS
func (g *dependencyGraph) detectCircularDependencies() error {
    visited := make(map[string]bool)
    recStack := make(map[string]bool)
    
    for actionName := range g.nodes {
        if !visited[actionName] {
            if g.isCyclicUtil(actionName, visited, recStack) {
                return fmt.Errorf("circular dependency detected")
            }
        }
    }
    
    return nil
}

// isCyclicUtil is a utility function for cycle detection
func (g *dependencyGraph) isCyclicUtil(actionName string, visited, recStack map[string]bool) bool {
    visited[actionName] = true
    recStack[actionName] = true
    
    for _, dep := range g.edges[actionName] {
        if !visited[dep] {
            if g.isCyclicUtil(dep, visited, recStack) {
                return true
            }
        } else if recStack[dep] {
            return true
        }
    }
    
    recStack[actionName] = false
    return false
}
```

### Step 4: Implement Topological Sorting

#### 4.1 Topological Sort for Orchestration Order
```go
// topologicalSort performs topological sorting to determine orchestration order
func (g *dependencyGraph) topologicalSort() ([]string, error) {
    var result []string
    visited := make(map[string]bool)
    temp := make(map[string]bool)
    
    var visit func(string) error
    visit = func(actionName string) error {
        if temp[actionName] {
            return fmt.Errorf("circular dependency detected")
        }
        if visited[actionName] {
            return nil
        }
        
        temp[actionName] = true
        
        for _, dep := range g.edges[actionName] {
            if err := visit(dep); err != nil {
                return err
            }
        }
        
        temp[actionName] = false
        visited[actionName] = true
        result = append(result, actionName)
        return nil
    }
    
    for actionName := range g.nodes {
        if !visited[actionName] {
            if err := visit(actionName); err != nil {
                return nil, err
            }
        }
    }
    
    return result, nil
}
```

### Step 5: Implement Orchestration Optimization

#### 5.1 Optimization Method
```go
// optimizeOrchestration creates an optimized orchestration strategy
func (m *Manager) optimizeOrchestration(actions []*types.Action, orchestrationOrder []string, graph *dependencyGraph, context *types.ActionContext) (*types.OrchestrationPlan, error) {
    plan := &OrchestrationPlan{
        Actions:        actions,
        OrchestrationOrder: orchestrationOrder,
        ConcurrencyLevel: context.MaxConcurrent,
    }
    
                    // Group actions by dependency level for parallel orchestration
        parallelGroups := m.groupActionsByLevel(orchestrationOrder, graph)
        plan.ParallelGroups = parallelGroups
        
        // Determine sequential vs parallel orchestration
        if context.Parallel && len(parallelGroups) > 1 {
            plan.SequentialOrder = orchestrationOrder // Fallback to sequential if parallel fails
        } else {
            plan.SequentialOrder = orchestrationOrder
        }
    
    // Estimate duration
    plan.EstimatedDuration = m.estimateRunDuration(actions, parallelGroups)
    
    // Define rollback strategy
            plan.RollbackStrategy = m.createRollbackStrategy(actions, orchestrationOrder)
    
    return plan, nil
}
```

#### 5.2 Action Level Grouping
```go
// groupActionsByLevel groups actions by dependency level
func (m *Manager) groupActionsByLevel(orchestrationOrder []string, graph *dependencyGraph) [][]string {
    levels := make([][]string, 0)
    currentLevel := make([]string, 0)
    completed := make(map[string]bool)
    
    for _, actionName := range orchestrationOrder {
        // Check if all dependencies are completed
        allDepsCompleted := true
        for _, dep := range graph.edges[actionName] {
            if !completed[dep] {
                allDepsCompleted = false
                break
            }
        }
        
        if allDepsCompleted {
            currentLevel = append(currentLevel, actionName)
        } else {
            // Start new level
            if len(currentLevel) > 0 {
                levels = append(levels, currentLevel)
                currentLevel = make([]string, 0)
            }
            currentLevel = append(currentLevel, actionName)
        }
        
        completed[actionName] = true
    }
    
    // Add final level
    if len(currentLevel) > 0 {
        levels = append(levels, currentLevel)
    }
    
    return levels
}
```

### Step 6: Implement Run Method

#### 6.1 Main Run Method
```go
// Run orchestrates actions with internal planning
func (m *Manager) Run(ctx context.Context, actions []*types.Action, context *types.ActionContext) (*types.OrchestrationResult, error) {
    m.logger.Debug("Starting action orchestration",
        logging.Int("action_count", len(actions)))

    // Create orchestration plan internally
    plan, err := m.planOrchestration(actions, context)
    if err != nil {
        return nil, fmt.Errorf("failed to create orchestration plan: %w", err)
    }

    // Prepare run context
    runContext := &collectionRunContext{
        plan:    plan,
        context: context,
        logger:  m.logger,
        results: make(map[string]*types.RunResult),
    }

    // Run based on plan strategy
    if len(plan.ParallelGroups) > 1 && context.Parallel {
        return m.runParallel(ctx, runContext)
    } else {
        return m.runSequential(ctx, runContext)
    }
}
```

#### 6.2 Action Orchestration Context
```go
// collectionRunContext holds context for action orchestration
type collectionRunContext struct {
    plan           *types.OrchestrationPlan
    context        *types.ActionContext
    logger         logging.Logger
    results        map[string]*types.RunResult
    rollbackStack  []*rollbackEntry
}
```

### Step 7: Implement Sequential Orchestration

#### 7.1 Sequential Orchestration Method
```go
// runSequential runs actions sequentially
func (m *Manager) runSequential(ctx context.Context, runContext *collectionRunContext) (*types.OrchestrationResult, error) {
    var allResults []*types.RunResult
    var errors []error

            for _, actionName := range runContext.plan.OrchestrationOrder {
        action := m.findActionByName(actionName, runContext.plan.Actions)
        if action == nil {
            continue
        }
        
                // Run action
        result, err := m.runAction(ctx, action, runContext.context)
        if err != nil {
            errors = append(errors, fmt.Errorf("action %s: %w", actionName, err))
            
            // Handle failure based on action configuration
            if action.Critical {
                return m.handleCollectionFailure(ctx, runContext, errors, allResults)
            } else if !action.AllowFailure {
                return m.handleCollectionFailure(ctx, runContext, errors, allResults)
            }
        }
        
        allResults = append(allResults, result)
        runContext.results[actionName] = result
        
        // Add to rollback stack if needed
        if action.Critical || !action.AllowFailure {
            runContext.rollbackStack = append(runContext.rollbackStack, &rollbackEntry{
                action: action,
                result: result,
            })
        }
    }

    return m.createOrchestrationResult(allResults, errors, types.OrchestrationStatusCompleted)
}
```

### Step 8: Implement Parallel Orchestration

#### 8.1 Parallel Orchestration Method
```go
// runParallel runs actions in parallel where possible
func (m *Manager) runParallel(ctx context.Context, runContext *collectionRunContext) (*types.OrchestrationResult, error) {
    var allResults []*types.RunResult
    var errors []error
    
            // Run each level in parallel
    for levelIndex, level := range runContext.plan.ParallelGroups {
        m.logger.Debug("Orchestrating action level",
            logging.Int("level", levelIndex+1),
            logging.Int("action_count", len(level)))
        
        levelResults, levelErrors := m.runActionLevel(ctx, level, runContext)
        
        // Check for critical failures
        for _, err := range levelErrors {
            errors = append(errors, err)
            // If any action in this level is critical and failed, stop orchestration
            for _, actionName := range level {
                action := m.findActionByName(actionName, runContext.plan.Actions)
                if action != nil && action.Critical {
                    return m.handleCollectionFailure(ctx, runContext, errors, allResults)
                }
            }
        }
        
        allResults = append(allResults, levelResults...)
        
        // Add to rollback stack
        for _, actionName := range level {
            if result, exists := runContext.results[actionName]; exists {
                action := m.findActionByName(actionName, runContext.plan.Actions)
                if action != nil && (action.Critical || !action.AllowFailure) {
                    runContext.rollbackStack = append(runContext.rollbackStack, &rollbackEntry{
                        action: action,
                        result: result,
                    })
                }
            }
        }
    }
    
    return m.createOrchestrationResult(allResults, errors, types.OrchestrationStatusCompleted)
}
```

#### 8.2 Action Level Orchestration
```go
// runActionLevel orchestrates a level of actions in parallel
func (m *Manager) runActionLevel(ctx context.Context, level []string, runContext *collectionRunContext) ([]*types.RunResult, []error) {
    if len(level) == 0 {
        return nil, nil
    }
    
            // Determine concurrency limit
        maxConcurrent := runContext.context.MaxConcurrent
    if maxConcurrent <= 0 {
        maxConcurrent = len(level)
    }
    
    // Create worker pool
    semaphore := make(chan struct{}, maxConcurrent)
            results := make(chan *actionOrchestrationResult, len(level))
    errors := make(chan error, len(level))
    
    // Start workers
    for _, actionName := range level {
        go func(name string) {
            semaphore <- struct{}{} // Acquire semaphore
            defer func() { <-semaphore }() // Release semaphore
            
            action := m.findActionByName(name, runContext.plan.Actions)
            if action == nil {
                errors <- fmt.Errorf("action %s not found", name)
                return
            }
            
            result, err := m.runAction(ctx, action, runContext.context)
            
            if err != nil {
                errors <- fmt.Errorf("action %s: %w", name, err)
            } else {
                results <- &actionOrchestrationResult{
                    actionName: name,
                    result:     result,
                }
            }
        }(actionName)
    }
    
    // Collect results
    var levelResults []*types.RunResult
    var levelErrors []error
    
    for i := 0; i < len(level); i++ {
        select {
        case actionResult := <-results:
            levelResults = append(levelResults, actionResult.result)
            runContext.results[actionResult.actionName] = actionResult.result
        case err := <-errors:
            levelErrors = append(levelErrors, err)
        case <-ctx.Done():
            return levelResults, append(levelErrors, fmt.Errorf("action orchestration cancelled: %w", ctx.Err()))
        }
    }
    
    return levelResults, levelErrors
}
```

### Step 9: Implement Action Orchestration

#### 9.1 Action Orchestration Method
```go
// runAction runs a single action
func (m *Manager) runAction(ctx context.Context, action *types.Action, context *types.ActionContext) (*types.RunResult, error) {
    m.logger.Debug("Orchestrating action",
        logging.String("action", action.Name),
        logging.String("type", action.Type))

            // Create actor for action orchestration
    actor := acting.NewActor(action, m.logger)
    
    // Orchestrate action
    result := &types.RunResult{
        ActionName: action.Name,
        Status:     types.RunStatusRunning,
        StartTime:  time.Now(),
    }
    
    err := actor.Act(ctx, context, result)
    if err != nil {
        result.Status = types.RunStatusFailed
        result.Error = err.Error()
        result.EndTime = time.Now()
        result.Duration = result.EndTime.Sub(*result.StartTime)
        return result, err
    }
    
            result.Status = types.RunStatusCompleted
    result.EndTime = time.Now()
    result.Duration = result.EndTime.Sub(*result.StartTime)
    
    return result, nil
}
```

### Step 10: Implement Rollback and Recovery

#### 10.1 Rollback Handling
```go
// handleCollectionFailure handles action orchestration failure and rollback
func (m *Manager) handleCollectionFailure(ctx context.Context, runContext *collectionRunContext, errors []error, results []*types.RunResult) (*types.OrchestrationResult, error) {
    m.logger.Error("Action orchestration failed, initiating rollback",
        logging.Int("error_count", len(errors)))
    
    // Perform rollback
    rollbackErrors := m.performRollback(ctx, runContext)
    
    // Combine original errors with rollback errors
    allErrors := append(errors, rollbackErrors...)
    
    return m.createOrchestrationResult(results, allErrors, types.OrchestrationStatusFailed)
}
```

#### 10.2 Rollback Orchestration
```go
// performRollback performs rollback of orchestrated actions
func (m *Manager) performRollback(ctx context.Context, runContext *collectionRunContext) []error {
    var rollbackErrors []error
    
            // Run rollback in reverse order
    for i := len(runContext.rollbackStack) - 1; i >= 0; i-- {
        entry := runContext.rollbackStack[i]
        
        m.logger.Info("Rolling back action",
            logging.String("action", entry.action.Name))
        
        // Run rollback action if defined
if entry.action.RollbackAction != "" {
            rollbackResult, err := m.runRollbackAction(ctx, entry.action, runContext.context)
            if err != nil {
                rollbackErrors = append(rollbackErrors, fmt.Errorf("rollback for action %s failed: %w", entry.action.Name, err))
            } else {
                m.logger.Info("Rollback completed",
                    logging.String("action", entry.action.Name),
                    logging.Int("exit_code", rollbackResult.ExitCode))
            }
        }
    }
    
    return rollbackErrors
}
```

### Step 11: Implement Result Aggregation

#### 11.1 Orchestration Result Creation
```go
// createOrchestrationResult creates the final orchestration result
func (m *Manager) createOrchestrationResult(actionResults []*types.RunResult, errors []error, status types.OrchestrationStatus) (*types.OrchestrationResult, error) {
    result := &types.OrchestrationResult{
        Status:    status,
        StartTime: time.Now(),
        EndTime:   time.Now(),
    }
    
    // Aggregate action results
    var outputs []string
    var stderrOutputs []string
    var exitCodes []int
    var successCount, failureCount int
    
    for _, actionResult := range actionResults {
        if actionResult.Status == types.RunStatusCompleted {
            successCount++
            if actionResult.ExitCode == 0 {
                outputs = append(outputs, fmt.Sprintf("[%s] %s", actionResult.ActionName, actionResult.Output))
            } else {
                stderrOutputs = append(stderrOutputs, fmt.Sprintf("[%s] %s", actionResult.ActionName, actionResult.Error))
            }
            exitCodes = append(exitCodes, actionResult.ExitCode)
        } else {
            failureCount++
            stderrOutputs = append(stderrOutputs, fmt.Sprintf("[%s] %s", actionResult.ActionName, actionResult.Error))
        }
    }
    
    // Set overall status
    if failureCount == 0 {
        result.Status = types.CollectionStatusCompleted
        result.ExitCode = 0
    } else if successCount > 0 {
        result.Status = types.CollectionStatusCompleted
        result.ExitCode = 1 // Partial success
    } else {
        result.Status = types.CollectionStatusFailed
        result.ExitCode = 1
    }
    
    // Aggregate output
    if len(outputs) > 0 {
        result.Output = strings.Join(outputs, "\n")
    }
    if len(stderrOutputs) > 0 {
        result.Error = strings.Join(stderrOutputs, "\n")
    }
    
    // Add metadata
    result.Metadata = map[string]interface{}{
        "total_actions":     len(actionResults),
        "successful_actions": successCount,
        "failed_actions":    failureCount,
        "exit_codes":        exitCodes,
        "error_count":       len(errors),
    }
    
    return result, nil
}
```

### Step 12: Add Supporting Structures

#### 12.1 Rollback Entry Structure
```go
// rollbackEntry represents an entry in the rollback stack
type rollbackEntry struct {
    action *types.Action
    result *types.RunResult
}

// actionOrchestrationResult holds result of action orchestration
type actionOrchestrationResult struct {
    actionName string
    result     *types.RunResult
}
```

#### 12.2 Helper Methods
```go
// findActionByName finds an action by name in the actions slice
func (m *Manager) findActionByName(name string, actions []*types.Action) *types.Action {
    for _, action := range actions {
        if action.Name == name {
            return action
        }
    }
    return nil
}

// estimateRunDuration estimates the total run duration
func (m *Manager) estimateRunDuration(actions []*types.Action, parallelGroups [][]string) time.Duration {
    // Simple estimation based on action count and parallel groups
    // In a real implementation, this would analyze action types and historical data
    baseTime := time.Duration(len(actions)) * 5 * time.Second
    if len(parallelGroups) > 1 {
        baseTime = baseTime / time.Duration(len(parallelGroups))
    }
    return baseTime
}

// createRollbackStrategy creates a rollback strategy for the actions
func (m *Manager) createRollbackStrategy(actions []*types.Action, orchestrationOrder []string) *RollbackStrategy {
    strategy := &RollbackStrategy{
        EnableRollback: true,
        RollbackOrder:  make([]string, len(orchestrationOrder)),
        CriticalActions: make([]string, 0),
    }
    
    // Copy orchestration order for rollback (reverse)
    copy(strategy.RollbackOrder, orchestrationOrder)
    
    // Identify critical actions
    for _, action := range actions {
        if action.Critical {
            strategy.CriticalActions = append(strategy.CriticalActions, action.Name)
        }
    }
    
    return strategy
}
```

## Configuration Options

### Supported Options
- **Parallel**: Enable parallel orchestration
- **MaxConcurrent**: Limit concurrent orchestrations
- **Timeout**: Action orchestration timeout
- **AllowFailure**: Continue on individual failures
- **Rollback**: Enable automatic rollback

## Dependencies

### Internal Dependencies
- `spooky/internal/types/actions`
- `spooky/internal/actions/acting`
- `spooky/internal/logging`

### External Dependencies
- `context` (standard library)
- `time` (standard library)
- `strings` (standard library)
- `fmt` (standard library)

## Implementation Order

1. Implement dependency graph construction
2. Add circular dependency detection
3. Implement topological sorting
4. Create orchestration plan structure
5. Implement orchestration optimization
6. Create action orchestration context
7. Implement sequential orchestration
8. Implement parallel orchestration
9. Add rollback and recovery
10. Implement result aggregation
11. Add supporting structures
12. Write comprehensive tests
13. Performance optimization
14. Documentation and cleanup

## API Summary

The final API will be clean and simple:

```go
type Manager struct {
    // Run orchestrates actions with internal planning
    Run(ctx context.Context, actions []*types.Action, context *types.ActionContext) (*types.OrchestrationResult, error)
}
```

This provides a stateless approach where planning happens internally as part of the orchestration process.

## CLI Workflow

### Debugging and Inspection
```bash
# See what would happen in details
spooky actions run <project> --debug --dry-run

# Save plan to file for manual inspection
spooky actions run <project> --debug --dry-run --output plan.json

# Log a run with high detail (<project>/logs/)
spooky actions run <project> --debug
```

### Normal Usage
```bash
# Preview execution
spooky actions run <project> --dry-run

# Execute orchestration
spooky actions run <project>
```


