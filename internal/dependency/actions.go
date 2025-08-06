package dependency

import (
	"fmt"
	"strings"
)

// ActionInfo represents basic information about an action for dependency management
type ActionInfo struct {
	Dependencies []string
	SourceFile   string
	SourceLine   int
}

// ActionDependencyManager extends the dependency system for actions
type ActionDependencyManager struct {
	graph *DependencyGraph
}

// NewActionDependencyManager creates a new action dependency manager
func NewActionDependencyManager() *ActionDependencyManager {
	return &ActionDependencyManager{
		graph: NewDependencyGraph(),
	}
}

// AddAction adds an action to the dependency graph
func (adm *ActionDependencyManager) AddAction(name string, dependencies []string, sourceFile string, sourceLine int) error {
	// Add the action node to the graph
	adm.graph.AddNode(name, "action", sourceFile, sourceLine)

	// Add dependencies if they exist
	if len(dependencies) > 0 {
		for _, dep := range dependencies {
			if err := adm.graph.AddDependency(name, dep); err != nil {
				return fmt.Errorf("failed to add dependency %s -> %s: %w", name, dep, err)
			}
		}
	}

	return nil
}

// AddActionCollection adds all actions from a collection to the dependency graph
func (adm *ActionDependencyManager) AddActionCollection(actions map[string]ActionInfo) error {
	// First, add all nodes to the graph
	for name, info := range actions {
		adm.graph.AddNode(name, "action", info.SourceFile, info.SourceLine)
	}

	// Then, add all dependencies
	for name, info := range actions {
		if len(info.Dependencies) > 0 {
			for _, dep := range info.Dependencies {
				if err := adm.graph.AddDependency(name, dep); err != nil {
					return fmt.Errorf("failed to add dependency %s -> %s: %w", name, dep, err)
				}
			}
		}
	}

	return nil
}

// ValidateDependencies validates all action dependencies
func (adm *ActionDependencyManager) ValidateDependencies() *ActionDependencyValidationResult {
	result := &ActionDependencyValidationResult{
		Valid:        true,
		Errors:       []ActionDependencyError{},
		Warnings:     []ActionDependencyWarning{},
		CircularRefs: []string{},
		MissingRefs:  []string{},
		SelfRefs:     []string{},
	}

	// Check for circular references
	if circular := adm.graph.DetectCircularRefs(); circular != nil {
		result.Valid = false
		result.CircularRefs = circular
		result.Errors = append(result.Errors, ActionDependencyError{
			Type:    "circular_reference",
			Message: fmt.Sprintf("Circular dependency detected: %s", strings.Join(circular, " -> ")),
			Details: circular,
		})
	}

	// Check for missing references
	missingRefs := adm.detectMissingReferences()
	if len(missingRefs) > 0 {
		result.Valid = false
		result.MissingRefs = missingRefs
		for _, missing := range missingRefs {
			result.Errors = append(result.Errors, ActionDependencyError{
				Type:    "missing_reference",
				Message: fmt.Sprintf("Action '%s' references non-existent action", missing),
				Details: []string{missing},
			})
		}
	}

	// Check for self-references
	selfRefs := adm.detectSelfReferences()
	if len(selfRefs) > 0 {
		result.Valid = false
		result.SelfRefs = selfRefs
		for _, selfRef := range selfRefs {
			result.Errors = append(result.Errors, ActionDependencyError{
				Type:    "self_reference",
				Message: fmt.Sprintf("Action '%s' depends on itself", selfRef),
				Details: []string{selfRef},
			})
		}
	}

	// Check for orphaned actions (actions with no dependents)
	orphaned := adm.detectOrphanedActions()
	if len(orphaned) > 0 {
		result.Warnings = append(result.Warnings, ActionDependencyWarning{
			Type:    "orphaned_action",
			Message: fmt.Sprintf("Orphaned actions found: %s", strings.Join(orphaned, ", ")),
			Details: orphaned,
		})
	}

	return result
}

// GetExecutionOrder returns the topological order for action execution
func (adm *ActionDependencyManager) GetExecutionOrder() ([]string, error) {
	// Validate dependencies first
	validation := adm.ValidateDependencies()
	if !validation.Valid {
		return nil, fmt.Errorf("dependency validation failed: %v", validation.Errors)
	}

	// Get topological order
	order, err := adm.graph.ResolveOrder()
	if err != nil {
		return nil, fmt.Errorf("failed to resolve execution order: %w", err)
	}

	return order, nil
}

// GetParallelGroups returns groups of actions that can be executed in parallel
func (adm *ActionDependencyManager) GetParallelGroups() ([][]string, error) {
	// Get execution order first
	order, err := adm.GetExecutionOrder()
	if err != nil {
		return nil, err
	}

	// Group actions by their dependency level
	groups := make([][]string, 0)
	levelMap := make(map[int][]string)
	maxLevel := 0

	// Calculate dependency level for each action
	for _, actionName := range order {
		level := adm.calculateDependencyLevel(actionName)
		levelMap[level] = append(levelMap[level], actionName)
		if level > maxLevel {
			maxLevel = level
		}
	}

	// Create groups in level order
	for level := 0; level <= maxLevel; level++ {
		if actions, exists := levelMap[level]; exists && len(actions) > 0 {
			groups = append(groups, actions)
		}
	}

	return groups, nil
}

// calculateDependencyLevel calculates the dependency level of an action
func (adm *ActionDependencyManager) calculateDependencyLevel(actionName string) int {
	node, exists := adm.graph.Nodes[actionName]
	if !exists {
		return 0
	}

	if len(node.Dependencies) == 0 {
		return 0
	}

	maxLevel := 0
	for _, dep := range node.Dependencies {
		depLevel := adm.calculateDependencyLevel(dep)
		if depLevel >= maxLevel {
			maxLevel = depLevel + 1
		}
	}

	return maxLevel
}

// GetDependencyChain returns the dependency chain for a specific action
func (adm *ActionDependencyManager) GetDependencyChain(actionName string) ([]string, error) {
	chain, err := adm.graph.GetDependencyChain(actionName)
	if err != nil {
		return nil, fmt.Errorf("failed to get dependency chain for %s: %w", actionName, err)
	}
	return chain, nil
}

// GetDependents returns all actions that depend on the specified action
func (adm *ActionDependencyManager) GetDependents(actionName string) ([]string, error) {
	dependents, err := adm.graph.GetDependents(actionName)
	if err != nil {
		return nil, fmt.Errorf("failed to get dependents for %s: %w", actionName, err)
	}
	return dependents, nil
}

// GetDependencies returns all actions that the specified action depends on
func (adm *ActionDependencyManager) GetDependencies(actionName string) ([]string, error) {
	dependencies, err := adm.graph.GetDependencies(actionName)
	if err != nil {
		return nil, fmt.Errorf("failed to get dependencies for %s: %w", actionName, err)
	}
	return dependencies, nil
}

// GetDependencyStats returns statistics about the dependency graph
func (adm *ActionDependencyManager) GetDependencyStats() map[string]interface{} {
	stats := adm.graph.GetDependencyStats()

	// Add action-specific statistics
	stats["total_actions"] = len(adm.graph.Nodes)

	// Count actions by dependency level
	levelCounts := make(map[int]int)
	for _, node := range adm.graph.Nodes {
		level := len(node.Dependencies)
		levelCounts[level]++
	}
	stats["dependency_levels"] = levelCounts

	return stats
}

// VisualizeDependencies returns a string representation of the dependency graph
func (adm *ActionDependencyManager) VisualizeDependencies() string {
	return adm.graph.VisualizeDependencies()
}

// detectMissingReferences detects actions that reference non-existent actions
func (adm *ActionDependencyManager) detectMissingReferences() []string {
	missing := make([]string, 0)
	existingActions := make(map[string]bool)

	// Build set of existing actions
	for actionName := range adm.graph.Nodes {
		existingActions[actionName] = true
	}

	// Check each action's dependencies
	for actionName, node := range adm.graph.Nodes {
		for _, dep := range node.Dependencies {
			if !existingActions[dep] {
				missing = append(missing, fmt.Sprintf("%s -> %s", actionName, dep))
			}
		}
	}

	return missing
}

// detectSelfReferences detects actions that depend on themselves
func (adm *ActionDependencyManager) detectSelfReferences() []string {
	selfRefs := make([]string, 0)

	for actionName, node := range adm.graph.Nodes {
		for _, dep := range node.Dependencies {
			if dep == actionName {
				selfRefs = append(selfRefs, actionName)
				break
			}
		}
	}

	return selfRefs
}

// detectOrphanedActions detects actions that have no dependents
func (adm *ActionDependencyManager) detectOrphanedActions() []string {
	orphaned := make([]string, 0)

	for actionName, node := range adm.graph.Nodes {
		if len(node.Dependents) == 0 {
			orphaned = append(orphaned, actionName)
		}
	}

	return orphaned
}

// allDependenciesProcessed checks if all dependencies of an action have been processed
func (adm *ActionDependencyManager) allDependenciesProcessed(actionName string, processed map[string]bool) bool {
	node, exists := adm.graph.Nodes[actionName]
	if !exists {
		return true
	}

	for _, dep := range node.Dependencies {
		if !processed[dep] {
			return false
		}
	}

	return true
}

// ActionDependencyValidationResult represents the result of dependency validation
type ActionDependencyValidationResult struct {
	Valid        bool                      `json:"valid"`
	Errors       []ActionDependencyError   `json:"errors,omitempty"`
	Warnings     []ActionDependencyWarning `json:"warnings,omitempty"`
	CircularRefs []string                  `json:"circular_refs,omitempty"`
	MissingRefs  []string                  `json:"missing_refs,omitempty"`
	SelfRefs     []string                  `json:"self_refs,omitempty"`
}

// ActionDependencyError represents a dependency validation error
type ActionDependencyError struct {
	Type    string   `json:"type"`
	Message string   `json:"message"`
	Details []string `json:"details,omitempty"`
}

// ActionDependencyWarning represents a dependency validation warning
type ActionDependencyWarning struct {
	Type    string   `json:"type"`
	Message string   `json:"message"`
	Details []string `json:"details,omitempty"`
}
