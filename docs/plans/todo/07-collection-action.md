# Implementation Plan: Action Orchestration

## Overview
Implement action orchestration functionality that can run multiple actions with internal planning, dependency resolution, parallel orchestration where possible, rollback on failure, and comprehensive result aggregation.

## Task Details
- **Task ID**: 3.1
- **Priority**: High
- **File**: `internal/actions/actor.go` (primary) and `internal/actions/manager.go` (secondary)
- **Functions**: `Run` (in actor.go), `ExecuteActionCollection` (in manager.go)

## Current State Analysis

### Existing Patterns
1. **Action Structure**: Actions use `spookyactionstypes.Action` with dependencies, priority, and run settings
2. **Context System**: `spookyactionstypes.ActionContext` provides run context with facts, variables, machines
3. **Result Structure**: `spookyactionstypes.RunResult` captures run results with output, error, exit code
4. **State Management**: Actions track state via `spookyactionstypes.ActionState`
5. **Error Handling**: Consistent error wrapping with context
6. **Dependency Management**: Actions have dependency arrays for run order

### Existing Implementation Examples
- **Action Manager**: `internal/actions/manager.go` provides action management with `ExecuteAction` and `ExecuteActionCollection`
- **Acting Manager**: `internal/actions/acting/manager.go` provides comprehensive acting system
- **Actor System**: `internal/actions/acting/actor.go` provides individual action execution
- **Action Types**: `internal/types/actions/action.go` defines action structure with dependencies
- **Context System**: `internal/types/actions/context.go` provides run context
- **Orchestration Types**: `internal/types/actions/orchestration.go` defines orchestration structures

### Current Implementation Status
- **Basic Framework**: `Run` method exists in `internal/actions/actor.go` with basic structure
- **Type Definitions**: `OrchestrationPlan`, `OrchestrationResult`, `RollbackStrategy` are defined
- **Acting System**: Complete acting system with managers, actors, and sessions
- **Planning System**: Basic planning infrastructure exists
- **Missing Components**: Dependency graph, topological sort, orchestration optimization, execution logic

## Implementation Requirements

### Interface Compliance
The action orchestration system must:
1. **Run actions** with internal planning, dependency resolution, and optimization
2. **Handle failures** and implement rollback strategies
3. **Support parallel orchestration** where dependencies allow
4. **Implement timeout** handling for the entire collection
5. **Support partial orchestration** and failure recovery
6. **Track orchestration state** and progress
7. **Handle circular dependencies** and validation
8. **Integrate with existing acting system** for individual action execution

### Required Dependencies
- Acting manager for individual action execution
- Dependency resolution system
- Rollback and recovery mechanisms
- Result aggregation system
- State tracking and progress monitoring

## Detailed Implementation Plan

### Step 1: Complete Dependency Graph Implementation

**Current Implementation** (in `internal/actions/actor.go`):
```go
// buildDependencyGraph builds a dependency graph from actions
func (m *Manager) buildDependencyGraph(actions []*spookytypes.Action) (*dependencyGraph, error) {
	// TODO: Implement dependency graph building
	return &dependencyGraph{}, nil
}
```

**Updated Implementation**:
```go
// dependencyGraph represents the dependency graph for actions
type dependencyGraph struct {
	nodes map[string]*actionNode
	edges map[string][]string // action -> dependencies
}

// actionNode represents a node in the dependency graph
type actionNode struct {
	action *spookytypes.Action
	state  spookytypes.ActionStatus
	result *spookytypes.RunResult
}

// buildDependencyGraph builds a dependency graph from actions
func (m *Manager) buildDependencyGraph(actions []*spookytypes.Action) (*dependencyGraph, error) {
	graph := &dependencyGraph{
		nodes: make(map[string]*actionNode),
		edges: make(map[string][]string),
	}
	
	// Create nodes for all actions
	for _, action := range actions {
		graph.nodes[action.Name] = &actionNode{
			action: action,
			state:  spookytypes.ActionStatusPending,
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

### Step 2: Implement Topological Sorting

**Current Implementation**:
```go
// topologicalSort performs topological sorting on the dependency graph
func (g *dependencyGraph) topologicalSort() ([]string, error) {
	// TODO: Implement topological sorting
	return []string{}, nil
}
```

**Updated Implementation**:
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

### Step 3: Complete Orchestration Optimization

**Current Implementation**:
```go
// optimizeOrchestration optimizes the orchestration strategy
func (m *Manager) optimizeOrchestration(actions []*spookytypes.Action, order []string, graph *dependencyGraph, context *spookytypes.ActionContext) (*spookytypes.OrchestrationPlan, error) {
	// TODO: Implement orchestration optimization
	return &spookytypes.OrchestrationPlan{
		SequentialOrder: order,
		ParallelGroups:  [][]string{order},
	}, nil
}
```

**Updated Implementation**:
```go
// optimizeOrchestration creates an optimized orchestration strategy
func (m *Manager) optimizeOrchestration(actions []*spookytypes.Action, orchestrationOrder []string, graph *dependencyGraph, context *spookytypes.ActionContext) (*spookytypes.OrchestrationPlan, error) {
	plan := &spookytypes.OrchestrationPlan{
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

### Step 4: Complete Sequential Orchestration

**Current Implementation**:
```go
// runSequential runs actions sequentially
func (m *Manager) runSequential(ctx context.Context, runContext *collectionRunContext) (*spookytypes.OrchestrationResult, error) {
	// TODO: Implement sequential execution
	return &spookytypes.OrchestrationResult{}, nil
}
```

**Updated Implementation**:
```go
// runSequential runs actions sequentially
func (m *Manager) runSequential(ctx context.Context, runContext *collectionRunContext) (*spookytypes.OrchestrationResult, error) {
	var allResults []*spookytypes.RunResult
	var errors []error

	for _, actionName := range runContext.plan.OrchestrationOrder {
		action := m.findActionByName(actionName, runContext.plan.Actions)
		if action == nil {
			continue
		}
		
		// Run action using acting system
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

	return m.createOrchestrationResult(allResults, errors, spookytypes.OrchestrationStatusCompleted)
}
```

### Step 5: Complete Parallel Orchestration

**Current Implementation**:
```go
// runParallel runs actions in parallel
func (m *Manager) runParallel(ctx context.Context, runContext *collectionRunContext) (*spookytypes.OrchestrationResult, error) {
	// TODO: Implement parallel execution
	return &spookytypes.OrchestrationResult{}, nil
}
```

**Updated Implementation**:
```go
// runParallel runs actions in parallel where possible
func (m *Manager) runParallel(ctx context.Context, runContext *collectionRunContext) (*spookytypes.OrchestrationResult, error) {
	var allResults []*spookytypes.RunResult
	var errors []error
	
	// Run each level in parallel
	for levelIndex, level := range runContext.plan.ParallelGroups {
		m.logger.Debug("Orchestrating action level",
			spookylogging.Int("level", levelIndex+1),
			spookylogging.Int("action_count", len(level)))
		
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
	
	return m.createOrchestrationResult(allResults, errors, spookytypes.OrchestrationStatusCompleted)
}

// runActionLevel orchestrates a level of actions in parallel
func (m *Manager) runActionLevel(ctx context.Context, level []string, runContext *collectionRunContext) ([]*spookytypes.RunResult, []error) {
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
	var levelResults []*spookytypes.RunResult
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

### Step 6: Complete Action Execution

**Current Implementation**:
```go
// runAction runs a single action
func (m *Manager) runAction(ctx context.Context, action *spookytypes.Action, context *spookytypes.ActionContext) (*spookytypes.RunResult, error) {
	// TODO: Implement single action execution
	return &spookytypes.RunResult{}, nil
}
```

**Updated Implementation**:
```go
// runAction runs a single action using the acting system
func (m *Manager) runAction(ctx context.Context, action *spookytypes.Action, context *spookytypes.ActionContext) (*spookytypes.RunResult, error) {
	m.logger.Debug("Orchestrating action",
		spookylogging.String("action", action.Name),
		spookylogging.String("type", action.Type))

	// Create actor for action orchestration
	actor := acting.NewActor(action, m.logger)
	
	// Orchestrate action
	result := &spookytypes.RunResult{
		ActionName: action.Name,
		Status:     spookytypes.RunStatusRunning,
		StartTime:  time.Now(),
	}
	
	err := actor.Act(ctx, context, result)
	if err != nil {
		result.Status = spookytypes.RunStatusFailed
		result.Error = err.Error()
		result.EndTime = time.Now()
		result.Duration = result.EndTime.Sub(*result.StartTime)
		return result, err
	}
	
	result.Status = spookytypes.RunStatusCompleted
	result.EndTime = time.Now()
	result.Duration = result.EndTime.Sub(*result.StartTime)
	
	return result, nil
}
```

### Step 7: Implement Rollback and Recovery

```go
// handleCollectionFailure handles action orchestration failure and rollback
func (m *Manager) handleCollectionFailure(ctx context.Context, runContext *collectionRunContext, errors []error, results []*spookytypes.RunResult) (*spookytypes.OrchestrationResult, error) {
	m.logger.Error("Action orchestration failed, initiating rollback",
		spookylogging.Int("error_count", len(errors)))
	
	// Perform rollback
	rollbackErrors := m.performRollback(ctx, runContext)
	
	// Combine original errors with rollback errors
	allErrors := append(errors, rollbackErrors...)
	
	return m.createOrchestrationResult(results, allErrors, spookytypes.OrchestrationStatusFailed)
}

// performRollback performs rollback of orchestrated actions
func (m *Manager) performRollback(ctx context.Context, runContext *collectionRunContext) []error {
	var rollbackErrors []error
	
	// Run rollback in reverse order
	for i := len(runContext.rollbackStack) - 1; i >= 0; i-- {
		entry := runContext.rollbackStack[i]
		
		m.logger.Info("Rolling back action",
			spookylogging.String("action", entry.action.Name))
		
		// Run rollback action if defined
		if entry.action.RollbackAction != "" {
			rollbackResult, err := m.runRollbackAction(ctx, entry.action, runContext.context)
			if err != nil {
				rollbackErrors = append(rollbackErrors, fmt.Errorf("rollback for action %s failed: %w", entry.action.Name, err))
			} else {
				m.logger.Info("Rollback completed",
					spookylogging.String("action", entry.action.Name),
					spookylogging.Int("exit_code", rollbackResult.ExitCode))
			}
		}
	}
	
	return rollbackErrors
}
```

### Step 8: Implement Result Aggregation

```go
// createOrchestrationResult creates the final orchestration result
func (m *Manager) createOrchestrationResult(actionResults []*spookytypes.RunResult, errors []error, status spookytypes.OrchestrationStatus) (*spookytypes.OrchestrationResult, error) {
	result := &spookytypes.OrchestrationResult{
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
		if actionResult.Status == spookytypes.RunStatusCompleted {
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
		result.Status = spookytypes.OrchestrationStatusCompleted
		result.ExitCode = 0
	} else if successCount > 0 {
		result.Status = spookytypes.OrchestrationStatusCompleted
		result.ExitCode = 1 // Partial success
	} else {
		result.Status = spookytypes.OrchestrationStatusFailed
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

### Step 9: Add Supporting Structures

```go
// Helper Methods
func (m *Manager) findActionByName(name string, actions []*spookytypes.Action) *spookytypes.Action {
	for _, action := range actions {
		if action.Name == name {
			return action
		}
	}
	return nil
}

func (m *Manager) estimateRunDuration(actions []*spookytypes.Action, parallelGroups [][]string) time.Duration {
	// Simple estimation based on action count and parallel groups
	// In a real implementation, this would analyze action types and historical data
	baseTime := time.Duration(len(actions)) * 5 * time.Second
	if len(parallelGroups) > 1 {
		baseTime = baseTime / time.Duration(len(parallelGroups))
	}
	return baseTime
}

func (m *Manager) createRollbackStrategy(actions []*spookytypes.Action, orchestrationOrder []string) *spookytypes.RollbackStrategy {
	strategy := &spookytypes.RollbackStrategy{
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

1. ✅ **Types defined** - OrchestrationPlan, OrchestrationResult, RollbackStrategy
2. ✅ **Basic framework** - Run method structure exists
3. 🔄 **Dependency graph construction** - TODO: Implement buildDependencyGraph
4. 🔄 **Circular dependency detection** - TODO: Implement detectCircularDependencies
5. 🔄 **Topological sorting** - TODO: Implement topologicalSort
6. 🔄 **Orchestration optimization** - TODO: Implement optimizeOrchestration
7. 🔄 **Action orchestration context** - TODO: Complete collectionRunContext
8. 🔄 **Sequential orchestration** - TODO: Implement runSequential
9. 🔄 **Parallel orchestration** - TODO: Implement runParallel
10. 🔄 **Action execution integration** - TODO: Implement runAction with acting system
11. 🔄 **Rollback and recovery** - TODO: Implement rollback logic
12. 🔄 **Result aggregation** - TODO: Implement createOrchestrationResult
13. 🔄 **Supporting structures** - TODO: Add helper methods
14. 🔄 **Comprehensive tests** - TODO: Write tests
15. 🔄 **Performance optimization** - TODO: Optimize performance
16. 🔄 **Documentation and cleanup** - TODO: Update documentation

## API Summary

The final API will be clean and simple:

```go
type Manager struct {
    // Run orchestrates actions with internal planning
    Run(ctx context.Context, actions []*spookytypes.Action, context *spookytypes.ActionContext) (*spookytypes.OrchestrationResult, error)
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

## Integration with Existing Systems

### Acting System Integration
The orchestration system integrates with the existing acting system:
- Uses `ActingManager` for individual action execution
- Leverages `Actor` interface for action execution
- Utilizes `ActingSession` for session management
- Integrates with `ActingResult` for result handling

### Planning System Integration
The orchestration system leverages the existing planning system:
- Uses `PlanningManager` for action planning
- Integrates with `ActionPlan` for plan management
- Utilizes `PlanValidator` for validation
- Leverages `PlanningOptimization` for optimization

### Type System Integration
The orchestration system uses the unified type system:
- Uses `spookytypes.Action` for action definitions
- Leverages `spookytypes.ActionContext` for context
- Utilizes `spookytypes.OrchestrationPlan` for planning
- Integrates with `spookytypes.OrchestrationResult` for results


