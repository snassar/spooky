package project

import (
	"fmt"
	"path/filepath"
	"strings"
)

// ProjectDependenciesEngine manages project dependencies and imports
type ProjectDependenciesEngine struct{}

// NewProjectDependenciesEngine creates a new project dependencies engine
func NewProjectDependenciesEngine() *ProjectDependenciesEngine {
	return &ProjectDependenciesEngine{}
}

// DependencyNode represents a node in the dependency graph
type DependencyNode struct {
	Project      *Project
	Path         string
	Dependencies []string
	Imports      []string
	Visited      bool
	InStack      bool
}

// DependencyGraph represents the project dependency graph
type DependencyGraph struct {
	Nodes map[string]*DependencyNode
	Edges map[string][]string
}

// DependencyResult represents the result of dependency analysis
type DependencyResult struct {
	Valid        bool
	Circular     bool
	Errors       []ValidationError
	Warnings     []ValidationError
	Graph        *DependencyGraph
	Order        []string // Topological order
	CircularPath []string // Path of circular dependency if found
}

// ImportResult represents the result of an import operation
type ImportResult struct {
	Success     bool
	Imported    []string
	Failed      []string
	Errors      []ValidationError
	Warnings    []ValidationError
	SharedVars  map[string]interface{}
	SharedFacts map[string]interface{}
}

// BuildDependencyGraph builds a dependency graph for projects
func (de *ProjectDependenciesEngine) BuildDependencyGraph(projects map[string]*Project) *DependencyGraph {
	graph := &DependencyGraph{
		Nodes: make(map[string]*DependencyNode),
		Edges: make(map[string][]string),
	}

	// Create nodes for all projects
	for path, project := range projects {
		node := &DependencyNode{
			Project:      project,
			Path:         path,
			Dependencies: make([]string, 0),
			Imports:      make([]string, 0),
			Visited:      false,
			InStack:      false,
		}

		// Add dependencies from project configuration
		if project.Dependencies != nil {
			node.Imports = project.Dependencies.Imports
		}

		graph.Nodes[path] = node
	}

	// Build edges based on imports
	for path, node := range graph.Nodes {
		for _, importPath := range node.Imports {
			// Resolve import path to actual project path
			resolvedPath := de.resolveImportPath(path, importPath)
			if resolvedPath != "" {
				graph.Edges[path] = append(graph.Edges[path], resolvedPath)
				node.Dependencies = append(node.Dependencies, resolvedPath)
			}
		}
	}

	return graph
}

// ValidateDependencies validates project dependencies for circular references
func (de *ProjectDependenciesEngine) ValidateDependencies(projects map[string]*Project) *DependencyResult {
	result := &DependencyResult{
		Valid:    true,
		Circular: false,
		Graph:    de.BuildDependencyGraph(projects),
	}

	// Check for circular dependencies using DFS
	for path := range result.Graph.Nodes {
		if !result.Graph.Nodes[path].Visited {
			if de.hasCircularDependency(result.Graph, path, result) {
				result.Valid = false
				result.Circular = true
				break
			}
		}
	}

	// Generate topological order if no circular dependencies
	if result.Valid {
		result.Order = de.topologicalSort(result.Graph)
	}

	return result
}

// ImportProject imports a project and its dependencies
func (de *ProjectDependenciesEngine) ImportProject(projectPath string, importPath string, projects map[string]*Project) *ImportResult {
	result := &ImportResult{
		Success:     true,
		Imported:    make([]string, 0),
		Failed:      make([]string, 0),
		SharedVars:  make(map[string]interface{}),
		SharedFacts: make(map[string]interface{}),
	}

	// Validate that the project exists
	project, exists := projects[projectPath]
	if !exists {
		result.Success = false
		result.Errors = append(result.Errors, ValidationError{
			Field:    "project",
			Message:  fmt.Sprintf("project not found: %s", projectPath),
			Severity: "error",
		})
		return result
	}

	// Validate that the import path exists
	resolvedImportPath := de.resolveImportPath(projectPath, importPath)
	if resolvedImportPath == "" {
		result.Success = false
		result.Errors = append(result.Errors, ValidationError{
			Field:    "import",
			Message:  fmt.Sprintf("import path not found: %s", importPath),
			Severity: "error",
		})
		return result
	}

	_, exists = projects[resolvedImportPath]
	if !exists {
		result.Success = false
		result.Errors = append(result.Errors, ValidationError{
			Field:    "import",
			Message:  fmt.Sprintf("imported project not found: %s", resolvedImportPath),
			Severity: "error",
		})
		return result
	}

	// Check for circular dependencies
	if de.wouldCreateCircularDependency(projectPath, resolvedImportPath, projects) {
		result.Success = false
		result.Errors = append(result.Errors, ValidationError{
			Field:    "import",
			Message:  fmt.Sprintf("import would create circular dependency: %s -> %s", projectPath, resolvedImportPath),
			Severity: "error",
		})
		return result
	}

	// Add import to project dependencies
	if project.Dependencies == nil {
		project.Dependencies = &ProjectDependencies{}
	}
	project.Dependencies.Imports = append(project.Dependencies.Imports, importPath)

	// Share variables and facts if specified
	if project.Dependencies.SharedVariables != nil {
		for _, varName := range project.Dependencies.SharedVariables {
			// This would need to be implemented with actual variable sharing logic
			result.SharedVars[varName] = fmt.Sprintf("shared_from_%s", resolvedImportPath)
		}
	}

	if project.Dependencies.SharedFacts != nil {
		for _, factName := range project.Dependencies.SharedFacts {
			// This would need to be implemented with actual facts sharing logic
			result.SharedFacts[factName] = fmt.Sprintf("shared_from_%s", resolvedImportPath)
		}
	}

	result.Imported = append(result.Imported, resolvedImportPath)

	return result
}

// GetDependencyOrder returns the topological order of project dependencies
func (de *ProjectDependenciesEngine) GetDependencyOrder(projects map[string]*Project) ([]string, error) {
	// Check for circular dependencies first
	result := de.ValidateDependencies(projects)
	if !result.Valid {
		return nil, fmt.Errorf("circular dependency detected: %v", result.CircularPath)
	}

	return result.Order, nil
}

// GetProjectDependencies returns all dependencies for a specific project
func (de *ProjectDependenciesEngine) GetProjectDependencies(projectPath string, projects map[string]*Project) ([]string, error) {
	graph := de.BuildDependencyGraph(projects)

	node, exists := graph.Nodes[projectPath]
	if !exists {
		return nil, fmt.Errorf("project not found: %s", projectPath)
	}

	return node.Dependencies, nil
}

// GetProjectDependents returns all projects that depend on a specific project
func (de *ProjectDependenciesEngine) GetProjectDependents(projectPath string, projects map[string]*Project) ([]string, error) {
	graph := de.BuildDependencyGraph(projects)

	dependents := make([]string, 0)
	for path, edges := range graph.Edges {
		for _, edge := range edges {
			if edge == projectPath {
				dependents = append(dependents, path)
				break
			}
		}
	}

	return dependents, nil
}

// hasCircularDependency checks for circular dependencies using DFS
func (de *ProjectDependenciesEngine) hasCircularDependency(graph *DependencyGraph, path string, result *DependencyResult) bool {
	node := graph.Nodes[path]

	if node.InStack {
		// Circular dependency found
		result.CircularPath = append(result.CircularPath, path)
		return true
	}

	if node.Visited {
		return false
	}

	node.Visited = true
	node.InStack = true

	// Check all dependencies
	for _, depPath := range graph.Edges[path] {
		if de.hasCircularDependency(graph, depPath, result) {
			result.CircularPath = append(result.CircularPath, path)
			return true
		}
	}

	node.InStack = false
	return false
}

// wouldCreateCircularDependency checks if adding an import would create a circular dependency
func (de *ProjectDependenciesEngine) wouldCreateCircularDependency(projectPath, importPath string, projects map[string]*Project) bool {
	// Create a temporary graph with the new import
	tempProjects := make(map[string]*Project)
	for k, v := range projects {
		tempProjects[k] = v
	}

	// Add the import temporarily
	if tempProjects[projectPath].Dependencies == nil {
		tempProjects[projectPath].Dependencies = &ProjectDependencies{}
	}
	tempProjects[projectPath].Dependencies.Imports = append(tempProjects[projectPath].Dependencies.Imports, importPath)

	// Check for circular dependencies
	result := de.ValidateDependencies(tempProjects)
	return !result.Valid
}

// topologicalSort performs topological sorting of the dependency graph
func (de *ProjectDependenciesEngine) topologicalSort(graph *DependencyGraph) []string {
	// Reset visited flags
	for _, node := range graph.Nodes {
		node.Visited = false
	}

	order := make([]string, 0)

	// Perform DFS for each unvisited node
	for path := range graph.Nodes {
		if !graph.Nodes[path].Visited {
			de.topologicalSortDFS(graph, path, &order)
		}
	}

	// Reverse the order to get topological sort
	for i, j := 0, len(order)-1; i < j; i, j = i+1, j-1 {
		order[i], order[j] = order[j], order[i]
	}

	return order
}

// topologicalSortDFS performs DFS for topological sorting
func (de *ProjectDependenciesEngine) topologicalSortDFS(graph *DependencyGraph, path string, order *[]string) {
	node := graph.Nodes[path]
	node.Visited = true

	// Visit all dependencies first
	for _, depPath := range graph.Edges[path] {
		if !graph.Nodes[depPath].Visited {
			de.topologicalSortDFS(graph, depPath, order)
		}
	}

	// Add current node to order
	*order = append(*order, path)
}

// resolveImportPath resolves an import path to an actual project path
func (de *ProjectDependenciesEngine) resolveImportPath(projectPath, importPath string) string {
	// Handle relative paths
	if strings.HasPrefix(importPath, "./") || strings.HasPrefix(importPath, "../") {
		return filepath.Join(filepath.Dir(projectPath), importPath)
	}

	// Handle absolute paths
	if filepath.IsAbs(importPath) {
		return importPath
	}

	// Handle module-style paths (e.g., "github.com/user/repo")
	// This would need to be implemented based on the module resolution strategy
	if strings.Contains(importPath, "/") {
		// For now, assume it's a relative path from the project root
		return filepath.Join(filepath.Dir(projectPath), importPath)
	}

	// Handle simple names (assume same directory)
	return filepath.Join(filepath.Dir(projectPath), importPath)
}

// ValidateImportPath validates that an import path is valid
func (de *ProjectDependenciesEngine) ValidateImportPath(importPath string) error {
	if importPath == "" {
		return fmt.Errorf("import path cannot be empty")
	}

	// Check for invalid characters
	if strings.Contains(importPath, "..") && !strings.HasPrefix(importPath, "../") {
		return fmt.Errorf("import path contains invalid relative path: %s", importPath)
	}

	// Check for absolute paths that might be problematic
	if filepath.IsAbs(importPath) {
		// Could add additional validation here
	}

	return nil
}

// GetSharedVariables returns variables that should be shared between projects
func (de *ProjectDependenciesEngine) GetSharedVariables(projectPath string, projects map[string]*Project) (map[string]interface{}, error) {
	return de.getSharedItems(projectPath, projects, func(deps *ProjectDependencies) []string {
		return deps.SharedVariables
	}, "shared_variable_")
}

// GetSharedFacts returns facts that should be shared between projects
func (de *ProjectDependenciesEngine) GetSharedFacts(projectPath string, projects map[string]*Project) (map[string]interface{}, error) {
	return de.getSharedItems(projectPath, projects, func(deps *ProjectDependencies) []string {
		return deps.SharedFacts
	}, "shared_fact_")
}

// getSharedItems is a helper function to reduce code duplication
func (de *ProjectDependenciesEngine) getSharedItems(projectPath string, projects map[string]*Project,
	getItems func(*ProjectDependencies) []string, prefix string) (map[string]interface{}, error) {
	project, exists := projects[projectPath]
	if !exists {
		return nil, fmt.Errorf("project not found: %s", projectPath)
	}

	if project.Dependencies == nil || len(getItems(project.Dependencies)) == 0 {
		return make(map[string]interface{}), nil
	}

	sharedItems := make(map[string]interface{})

	// This would need to be implemented with actual resolution logic
	for _, itemName := range getItems(project.Dependencies) {
		sharedItems[itemName] = fmt.Sprintf("%s%s", prefix, itemName)
	}

	return sharedItems, nil
}
