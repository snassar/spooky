package dependency

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"spooky/internal/logging"
)

// DependencyManager provides high-level dependency management
type DependencyManager struct {
	graph  *DependencyGraph
	logger logging.Logger
	cache  map[string]interface{}
	mutex  sync.RWMutex
}

// DependencyError represents a dependency-related error
type DependencyError struct {
	Type        string   // "circular", "missing", "invalid"
	Message     string
	Nodes       []string
	File        string
	Line        int
	Suggestions []string
}

func (e *DependencyError) Error() string {
	return e.Message
}

// NewDependencyManager creates a new dependency manager
func NewDependencyManager(logger logging.Logger) *DependencyManager {
	return &DependencyManager{
		graph:  NewDependencyGraph(),
		logger: logger,
		cache:  make(map[string]interface{}),
	}
}

// AddVariable adds a variable to the dependency graph
func (dm *DependencyManager) AddVariable(name, file string, line int, dependencies []string) error {
	dm.logger.Debug("Adding variable to dependency graph",
		logging.String("variable", name),
		logging.String("file", file),
		logging.Int("line", line))

	// Add the variable node
	dm.graph.AddNode(name, "variable", file, line)

	// Add dependencies
	for _, dep := range dependencies {
		if err := dm.graph.AddDependency(name, dep); err != nil {
			return fmt.Errorf("failed to add dependency %s -> %s: %w", name, dep, err)
		}
	}

	dm.logger.Debug("Variable added successfully", logging.String("variable", name))
	return nil
}

// AddAction adds an action to the dependency graph
func (dm *DependencyManager) AddAction(name, file string, line int, dependencies []string) error {
	dm.logger.Debug("Adding action to dependency graph",
		logging.String("action", name),
		logging.String("file", file),
		logging.Int("line", line))

	// Add the action node
	dm.graph.AddNode(name, "action", file, line)

	// Add dependencies
	for _, dep := range dependencies {
		if err := dm.graph.AddDependency(name, dep); err != nil {
			return fmt.Errorf("failed to add dependency %s -> %s: %w", name, dep, err)
		}
	}

	dm.logger.Debug("Action added successfully", logging.String("action", name))
	return nil
}

// ValidateDependencies validates the entire dependency graph
func (dm *DependencyManager) ValidateDependencies() []*DependencyError {
	dm.logger.Info("Validating dependency graph")
	errors := []*DependencyError{}

	// Check for circular dependencies
	if circular := dm.graph.DetectCircularRefs(); circular != nil {
		errors = append(errors, &DependencyError{
			Type:    "circular",
			Message: fmt.Sprintf("Circular dependency detected: %s", strings.Join(circular, " -> ")),
			Nodes:   circular,
			Suggestions: []string{
				"Review the dependency chain and remove circular references",
				"Consider using computed values instead of direct dependencies",
				"Break the cycle by introducing intermediate variables",
			},
		})
	}

	// Check for missing dependencies
	if missing := dm.graph.ValidateDependencies(); len(missing) > 0 {
		for _, miss := range missing {
			errors = append(errors, &DependencyError{
				Type:    "missing",
				Message: fmt.Sprintf("Missing dependency: %s", miss),
				Nodes:   []string{miss},
				Suggestions: []string{
					"Define the missing dependency",
					"Check for typos in dependency names",
					"Ensure all referenced variables/actions are defined",
				},
			})
		}
	}

	// Check for self-references
	nodes := dm.graph.GetAllNodes()
	for name, node := range nodes {
		for _, dep := range node.Dependencies {
			if dep == name {
				errors = append(errors, &DependencyError{
					Type:    "self_reference",
					Message: fmt.Sprintf("Self-reference detected: %s depends on itself", name),
					Nodes:   []string{name},
					File:    node.File,
					Line:    node.Line,
					Suggestions: []string{
						"Remove the self-reference",
						"Use a different variable name",
						"Consider if this dependency is necessary",
					},
				})
			}
		}
	}

	dm.logger.Info("Dependency validation completed",
		logging.Int("errors", len(errors)),
		logging.Int("nodes", len(nodes)))

	return errors
}

// GetResolutionOrder returns the order in which variables should be resolved
func (dm *DependencyManager) GetResolutionOrder() ([]string, error) {
	dm.logger.Debug("Calculating resolution order")

	order, err := dm.graph.ResolveOrder()
	if err != nil {
		return nil, fmt.Errorf("failed to determine resolution order: %w", err)
	}

	dm.logger.Debug("Resolution order calculated",
		logging.Int("variables", len(order)))

	return order, nil
}

// GetDependencyChain returns the dependency chain for a specific variable
func (dm *DependencyManager) GetDependencyChain(name string) ([]string, error) {
	dm.logger.Debug("Getting dependency chain", logging.String("variable", name))

	chain, err := dm.graph.GetDependencyChain(name)
	if err != nil {
		return nil, fmt.Errorf("failed to get dependency chain for %s: %w", name, err)
	}

	dm.logger.Debug("Dependency chain retrieved",
		logging.String("variable", name))

	return chain, nil
}

// GetDependents returns all variables that depend on the given variable
func (dm *DependencyManager) GetDependents(name string) ([]string, error) {
	dm.logger.Debug("Getting dependents", logging.String("variable", name))

	dependents, err := dm.graph.GetDependents(name)
	if err != nil {
		return nil, fmt.Errorf("failed to get dependents for %s: %w", name, err)
	}

	dm.logger.Debug("Dependents retrieved",
		logging.String("variable", name))

	return dependents, nil
}

// GetDependencies returns all dependencies for the given variable
func (dm *DependencyManager) GetDependencies(name string) ([]string, error) {
	dm.logger.Debug("Getting dependencies", logging.String("variable", name))

	dependencies, err := dm.graph.GetDependencies(name)
	if err != nil {
		return nil, fmt.Errorf("failed to get dependencies for %s: %w", name, err)
	}

	dm.logger.Debug("Dependencies retrieved",
		logging.String("variable", name))

	return dependencies, nil
}

// GetDependencyStats returns statistics about the dependency graph
func (dm *DependencyManager) GetDependencyStats() map[string]interface{} {
	dm.logger.Debug("Getting dependency statistics")

	stats := dm.graph.GetDependencyStats()

	dm.logger.Debug("Dependency statistics retrieved",
		logging.Int("total_nodes", stats["total_nodes"].(int)),
		logging.Int("total_edges", stats["total_edges"].(int)))

	return stats
}

// VisualizeDependencies returns a string representation of the dependency graph
func (dm *DependencyManager) VisualizeDependencies() string {
	dm.logger.Debug("Generating dependency visualization")

	visualization := dm.graph.VisualizeDependencies()

	dm.logger.Debug("Dependency visualization generated",
		logging.Int("length", len(visualization)))

	return visualization
}

// Clear clears the dependency graph
func (dm *DependencyManager) Clear() {
	dm.logger.Info("Clearing dependency graph")

	dm.graph.Clear()
	dm.mutex.Lock()
	dm.cache = make(map[string]interface{})
	dm.mutex.Unlock()

	dm.logger.Info("Dependency graph cleared")
}

// CacheResult caches a computation result
func (dm *DependencyManager) CacheResult(key string, value interface{}) {
	dm.mutex.Lock()
	defer dm.mutex.Unlock()

	dm.cache[key] = value
	dm.logger.Debug("Cached result", logging.String("key", key))
}

// GetCachedResult retrieves a cached computation result
func (dm *DependencyManager) GetCachedResult(key string) (interface{}, bool) {
	dm.mutex.RLock()
	defer dm.mutex.RUnlock()

	value, exists := dm.cache[key]
	if exists {
		dm.logger.Debug("Retrieved cached result", logging.String("key", key))
	}
	return value, exists
}

// ClearCache clears the computation cache
func (dm *DependencyManager) ClearCache() {
	dm.mutex.Lock()
	defer dm.mutex.Unlock()

	dm.cache = make(map[string]interface{})
	dm.logger.Debug("Computation cache cleared")
}

// GetNodeInfo returns detailed information about a specific node
func (dm *DependencyManager) GetNodeInfo(name string) (map[string]interface{}, error) {
	dm.logger.Debug("Getting node info", logging.String("node", name))

	node, exists := dm.graph.GetNode(name)
	if !exists {
		return nil, fmt.Errorf("node '%s' not found", name)
	}

	dependencies, _ := dm.graph.GetDependencies(name)
	dependents, _ := dm.graph.GetDependents(name)

	info := map[string]interface{}{
		"name":         node.Name,
		"type":         node.Type,
		"file":         node.File,
		"line":         node.Line,
		"dependencies": dependencies,
		"dependents":   dependents,
		"metadata":     node.Metadata,
	}

	dm.logger.Debug("Node info retrieved",
		logging.String("node", name),
		logging.Int("dependencies", len(dependencies)),
		logging.Int("dependents", len(dependents)))

	return info, nil
}

// ValidateVariableDependencies validates dependencies for a specific variable
func (dm *DependencyManager) ValidateVariableDependencies(name string) []*DependencyError {
	dm.logger.Debug("Validating variable dependencies", logging.String("variable", name))

	errors := []*DependencyError{}

	// Check if variable exists
	node, exists := dm.graph.GetNode(name)
	if !exists {
		errors = append(errors, &DependencyError{
			Type:    "missing",
			Message: fmt.Sprintf("Variable '%s' not found in dependency graph", name),
			Nodes:   []string{name},
			Suggestions: []string{
				"Ensure the variable is defined in variables.hcl or variables/*.hcl",
				"Check for typos in the variable name",
				"Verify the variable file is being loaded",
			},
		})
		return errors
	}

	// Check dependencies
	dependencies, err := dm.graph.GetDependencies(name)
	if err != nil {
		errors = append(errors, &DependencyError{
			Type:    "invalid",
			Message: fmt.Sprintf("Failed to get dependencies for '%s': %v", name, err),
			Nodes:   []string{name},
			File:    node.File,
			Line:    node.Line,
		})
		return errors
	}

	// Check each dependency
	for _, dep := range dependencies {
		if _, exists := dm.graph.GetNode(dep); !exists {
			errors = append(errors, &DependencyError{
				Type:    "missing_dependency",
				Message: fmt.Sprintf("Variable '%s' depends on undefined variable '%s'", name, dep),
				Nodes:   []string{name, dep},
				File:    node.File,
				Line:    node.Line,
				Suggestions: []string{
					fmt.Sprintf("Define variable '%s' in variables.hcl or variables/*.hcl", dep),
					"Check for typos in the dependency name",
					"Ensure the dependency variable is loaded before this variable",
				},
			})
		}
	}

	dm.logger.Debug("Variable dependency validation completed",
		logging.String("variable", name),
		logging.Int("errors", len(errors)))

	return errors
}

// GetImpactAnalysis returns the impact of changing a specific variable
func (dm *DependencyManager) GetImpactAnalysis(name string) (map[string]interface{}, error) {
	dm.logger.Debug("Performing impact analysis", logging.String("variable", name))

	// Check if variable exists
	if _, exists := dm.graph.GetNode(name); !exists {
		return nil, fmt.Errorf("variable '%s' not found", name)
	}

	// Get all dependents (direct and indirect)
	impacted := dm.getImpactedVariables(name, make(map[string]bool))

	// Get dependency chain
	chain, _ := dm.graph.GetDependencyChain(name)

	analysis := map[string]interface{}{
		"variable":        name,
		"impacted_count":  len(impacted),
		"impacted_vars":   impacted,
		"dependency_chain": chain,
		"impact_level":    dm.calculateImpactLevel(len(impacted)),
	}

	dm.logger.Debug("Impact analysis completed",
		logging.String("variable", name),
		logging.Int("impacted_count", len(impacted)))

	return analysis, nil
}

// getImpactedVariables recursively finds all variables impacted by a change
func (dm *DependencyManager) getImpactedVariables(name string, visited map[string]bool) []string {
	if visited[name] {
		return []string{}
	}

	visited[name] = true
	impacted := []string{name}

	dependents, err := dm.graph.GetDependents(name)
	if err != nil {
		return impacted
	}

	for _, dep := range dependents {
		impacted = append(impacted, dm.getImpactedVariables(dep, visited)...)
	}

	return impacted
}

// calculateImpactLevel determines the impact level based on the number of affected variables
func (dm *DependencyManager) calculateImpactLevel(impactedCount int) string {
	switch {
	case impactedCount == 1:
		return "low"
	case impactedCount <= 5:
		return "medium"
	case impactedCount <= 20:
		return "high"
	default:
		return "critical"
	}
}

// GetDependencyReport generates a comprehensive dependency report
func (dm *DependencyManager) GetDependencyReport() map[string]interface{} {
	dm.logger.Info("Generating dependency report")

	stats := dm.GetDependencyStats()
	errors := dm.ValidateDependencies()

	report := map[string]interface{}{
		"generated_at":    time.Now().Format(time.RFC3339),
		"statistics":      stats,
		"validation": map[string]interface{}{
			"valid":   len(errors) == 0,
			"errors":  len(errors),
			"details": errors,
		},
		"visualization": dm.VisualizeDependencies(),
	}

	dm.logger.Info("Dependency report generated",
		logging.Int("total_nodes", stats["total_nodes"].(int)),
		logging.Int("validation_errors", len(errors)))

	return report
} 