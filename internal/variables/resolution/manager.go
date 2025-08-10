package resolution

import (
	"context"
	"fmt"
	spookyinterfaces "spooky/internal/interfaces"
	spookytypeslogging "spooky/internal/types/logging"

	spookytypesvariables "spooky/internal/types/variables"
)

// Manager implements ResolutionManager interface
type Manager struct {
	config    *spookytypesvariables.ResolutionConfig
	resolvers map[string]spookyinterfaces.VariableResolver
	logger    spookytypeslogging.Logger
}

// NewManager creates a new resolution manager
func NewManager(config *spookytypesvariables.ResolutionConfig, logger spookytypeslogging.Logger) *Manager {
	return &Manager{
		config:    config,
		resolvers: make(map[string]spookyinterfaces.VariableResolver),
		logger:    logger,
	}
}

// ResolveVariable resolves a single variable
func (m *Manager) ResolveVariable(ctx context.Context, variable *spookytypesvariables.Variable, context *spookytypesvariables.VariableContext) error {
	// 1. Validate variable is not already resolved
	if variable.Resolved {
		return nil
	}

	// 2. Get appropriate resolver for variable type
	resolver, exists := m.resolvers[variable.Type]
	if !exists {
		return fmt.Errorf("no resolver found for variable type: %s", variable.Type)
	}

	// 3. Resolve variable
	if err := resolver.Resolve(ctx, variable, context); err != nil {
		return fmt.Errorf("failed to resolve variable %s: %w", variable.Name, err)
	}

	// 4. Mark as resolved
	variable.Resolved = true
	return nil
}

// ResolveDependencies resolves all variables with dependency management
func (m *Manager) ResolveDependencies(ctx context.Context, variables []*spookytypesvariables.Variable) error {
	// 1. Validate dependencies
	if err := m.ValidateDependencies(ctx, variables); err != nil {
		return fmt.Errorf("dependency validation failed: %w", err)
	}

	// 2. Detect circular dependencies
	if err := m.DetectCircularDependencies(variables); err != nil {
		return fmt.Errorf("circular dependency detected: %w", err)
	}

	// 3. Get resolution order
	order, err := m.GetResolutionOrder(variables)
	if err != nil {
		return fmt.Errorf("failed to determine resolution order: %w", err)
	}

	// 4. Resolve variables in order
	context := &spookytypesvariables.VariableContext{
		Variables: make(map[string]*spookytypesvariables.Variable),
	}

	for _, variable := range order {
		if err := m.ResolveVariable(ctx, variable, context); err != nil {
			return fmt.Errorf("failed to resolve variable %s: %w", variable.Name, err)
		}

		// Add resolved variable to context
		context.Variables[variable.Name] = variable
	}

	return nil
}

// ResolveContext resolves all variables in a context
func (m *Manager) ResolveContext(ctx context.Context, context *spookytypesvariables.VariableContext) error {
	// Get all variables from context
	variables := make([]*spookytypesvariables.Variable, 0, len(context.Variables))
	for _, variable := range context.Variables {
		variables = append(variables, variable)
	}

	// Resolve dependencies
	return m.ResolveDependencies(ctx, variables)
}

// ValidateDependencies validates variable dependencies
func (m *Manager) ValidateDependencies(ctx context.Context, variables []*spookytypesvariables.Variable) error {
	// Validate all dependencies exist
	variableMap := make(map[string]*spookytypesvariables.Variable)
	for _, variable := range variables {
		variableMap[variable.Name] = variable
	}

	for _, variable := range variables {
		for _, depName := range variable.Dependencies {
			if _, exists := variableMap[depName]; !exists {
				return fmt.Errorf("variable %s depends on undefined variable %s", variable.Name, depName)
			}
		}
	}

	return nil
}

// DetectCircularDependencies detects circular dependencies in variables
func (m *Manager) DetectCircularDependencies(variables []*spookytypesvariables.Variable) error {
	// Build dependency graph
	graph := make(map[string][]string)
	for _, variable := range variables {
		graph[variable.Name] = variable.Dependencies
	}

	// Use DFS to detect cycles
	visited := make(map[string]bool)
	recStack := make(map[string]bool)

	var dfs func(node string) bool
	dfs = func(node string) bool {
		visited[node] = true
		recStack[node] = true

		for _, neighbor := range graph[node] {
			if !visited[neighbor] {
				if dfs(neighbor) {
					return true
				}
			} else if recStack[neighbor] {
				return true
			}
		}

		recStack[node] = false
		return false
	}

	for node := range graph {
		if !visited[node] {
			if dfs(node) {
				return fmt.Errorf("circular dependency detected")
			}
		}
	}

	return nil
}

// GetDependencyGraph builds a dependency graph from variables
func (m *Manager) GetDependencyGraph(variables []*spookytypesvariables.Variable) (*spookytypesvariables.DependencyGraph, error) {
	graph := &spookytypesvariables.DependencyGraph{
		Nodes: make(map[string]*spookytypesvariables.DependencyNode),
		Edges: make(map[string][]string),
	}

	for _, variable := range variables {
		graph.Nodes[variable.Name] = &spookytypesvariables.DependencyNode{
			Name:         variable.Name,
			Dependencies: variable.Dependencies,
			Resolved:     variable.Resolved,
		}
		graph.Edges[variable.Name] = variable.Dependencies
	}

	return graph, nil
}

// SetMaxRecursionDepth sets the maximum recursion depth for resolution
func (m *Manager) SetMaxRecursionDepth(depth int) error {
	if m.config == nil {
		m.config = &spookytypesvariables.ResolutionConfig{}
	}
	m.config.MaxRecursionDepth = depth
	return nil
}

// SetDefaultValues sets default values for variables
func (m *Manager) SetDefaultValues(defaults map[string]interface{}) error {
	if m.config == nil {
		m.config = &spookytypesvariables.ResolutionConfig{}
	}
	m.config.DefaultValues = defaults
	return nil
}

// EnableStrictMode enables or disables strict mode
func (m *Manager) EnableStrictMode(strict bool) error {
	if m.config == nil {
		m.config = &spookytypesvariables.ResolutionConfig{}
	}
	m.config.StrictMode = strict
	return nil
}

// GetUnresolvedVariables returns variables that are not yet resolved
func (m *Manager) GetUnresolvedVariables(variables []*spookytypesvariables.Variable) []*spookytypesvariables.Variable {
	var unresolved []*spookytypesvariables.Variable
	for _, variable := range variables {
		if !variable.Resolved {
			unresolved = append(unresolved, variable)
		}
	}
	return unresolved
}

// GetResolutionOrder determines the order in which variables should be resolved
func (m *Manager) GetResolutionOrder(variables []*spookytypesvariables.Variable) ([]*spookytypesvariables.Variable, error) {
	// Build dependency graph
	graph := make(map[string][]string)
	variableMap := make(map[string]*spookytypesvariables.Variable)

	for _, variable := range variables {
		variableMap[variable.Name] = variable
		graph[variable.Name] = variable.Dependencies
	}

	// Topological sort
	return m.topologicalSort(graph, variableMap)
}

// Close closes the resolution manager and releases resources
func (m *Manager) Close() error {
	// Clean up resources if needed
	return nil
}

// Helper method for topological sorting
func (m *Manager) topologicalSort(graph map[string][]string, variableMap map[string]*spookytypesvariables.Variable) ([]*spookytypesvariables.Variable, error) {
	// Implementation of topological sort algorithm
	// Returns variables in dependency order
	visited := make(map[string]bool)
	temp := make(map[string]bool)
	var result []*spookytypesvariables.Variable

	var visit func(node string) error
	visit = func(node string) error {
		if temp[node] {
			return fmt.Errorf("circular dependency detected")
		}
		if visited[node] {
			return nil
		}

		temp[node] = true

		for _, neighbor := range graph[node] {
			if err := visit(neighbor); err != nil {
				return err
			}
		}

		temp[node] = false
		visited[node] = true

		if variable, exists := variableMap[node]; exists {
			result = append(result, variable)
		}

		return nil
	}

	for node := range graph {
		if !visited[node] {
			if err := visit(node); err != nil {
				return nil, err
			}
		}
	}

	// Reverse to get correct order
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}

	return result, nil
}
