package dependency

import (
	"fmt"
	"sort"
	"strings"
	"sync"
)

// DependencyNode represents a node in the dependency graph
type DependencyNode struct {
	Name         string
	Dependencies []string
	Dependents   []string
	File         string
	Line         int
	Type         string // "variable", "action", "template"
	Metadata     map[string]interface{}
}

// DependencyGraph manages the dependency relationships
type DependencyGraph struct {
	Nodes map[string]*DependencyNode
	Edges map[string][]string
	mutex sync.RWMutex
}

// NewDependencyGraph creates a new dependency graph
func NewDependencyGraph() *DependencyGraph {
	return &DependencyGraph{
		Nodes: make(map[string]*DependencyNode),
		Edges: make(map[string][]string),
	}
}

// AddNode adds a node to the dependency graph
func (dg *DependencyGraph) AddNode(name, fileType, file string, line int) *DependencyNode {
	dg.mutex.Lock()
	defer dg.mutex.Unlock()

	node := &DependencyNode{
		Name:         name,
		Dependencies: []string{},
		Dependents:   []string{},
		File:         file,
		Line:         line,
		Type:         fileType,
		Metadata:     make(map[string]interface{}),
	}

	dg.Nodes[name] = node
	dg.Edges[name] = []string{}
	return node
}

// AddDependency adds a dependency relationship between two nodes
func (dg *DependencyGraph) AddDependency(from, to string) error {
	dg.mutex.Lock()
	defer dg.mutex.Unlock()

	// Check if both nodes exist
	if _, exists := dg.Nodes[from]; !exists {
		return fmt.Errorf("source node '%s' does not exist", from)
	}
	if _, exists := dg.Nodes[to]; !exists {
		return fmt.Errorf("target node '%s' does not exist", to)
	}

	// Check for self-reference
	if from == to {
		return fmt.Errorf("self-reference detected: '%s' depends on itself", from)
	}

	// Add dependency
	dg.Nodes[from].Dependencies = append(dg.Nodes[from].Dependencies, to)
	dg.Nodes[to].Dependents = append(dg.Nodes[to].Dependents, from)
	dg.Edges[from] = append(dg.Edges[from], to)

	return nil
}

// DetectCircularRefs detects circular references in the dependency graph
func (dg *DependencyGraph) DetectCircularRefs() []string {
	dg.mutex.RLock()
	defer dg.mutex.RUnlock()

	visited := make(map[string]bool)
	recStack := make(map[string]bool)
	circular := []string{}

	for nodeName := range dg.Nodes {
		if !visited[nodeName] {
			if dg.hasCycle(nodeName, visited, recStack, &circular) {
				// Reverse the circular path to show the correct order
				for i, j := 0, len(circular)-1; i < j; i, j = i+1, j-1 {
					circular[i], circular[j] = circular[j], circular[i]
				}
				return circular
			}
		}
	}
	return nil
}

// hasCycle performs depth-first search to detect cycles
func (dg *DependencyGraph) hasCycle(node string, visited, recStack map[string]bool, circular *[]string) bool {
	visited[node] = true
	recStack[node] = true

	// Check if node exists before accessing its dependencies
	nodeObj, exists := dg.Nodes[node]
	if !exists {
		// Node doesn't exist, no cycle possible
		recStack[node] = false
		return false
	}

	for _, dep := range nodeObj.Dependencies {
		if !visited[dep] {
			if dg.hasCycle(dep, visited, recStack, circular) {
				*circular = append(*circular, node)
				return true
			}
		} else if recStack[dep] {
			*circular = append(*circular, node)
			return true
		}
	}

	recStack[node] = false
	return false
}

// ResolveOrder performs topological sorting to determine resolution order
func (dg *DependencyGraph) ResolveOrder() ([]string, error) {
	dg.mutex.RLock()
	defer dg.mutex.RUnlock()

	// Check for circular dependencies first
	if circular := dg.DetectCircularRefs(); circular != nil {
		return nil, fmt.Errorf("circular dependency detected: %s", strings.Join(circular, " -> "))
	}

	// Calculate in-degrees
	inDegree := make(map[string]int)
	queue := []string{}
	result := []string{}

	// Initialize in-degrees
	for nodeName := range dg.Nodes {
		inDegree[nodeName] = len(dg.Nodes[nodeName].Dependencies)
		if inDegree[nodeName] == 0 {
			queue = append(queue, nodeName)
		}
	}

	// Process queue (topological sort)
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		result = append(result, current)

		// Reduce in-degree for dependents
		for _, dependent := range dg.Nodes[current].Dependents {
			inDegree[dependent]--
			if inDegree[dependent] == 0 {
				queue = append(queue, dependent)
			}
		}
	}

	// Check if all nodes were processed
	if len(result) != len(dg.Nodes) {
		return nil, fmt.Errorf("circular dependency detected (incomplete resolution)")
	}

	return result, nil
}

// GetDependencies returns all dependencies for a given node
func (dg *DependencyGraph) GetDependencies(nodeName string) ([]string, error) {
	dg.mutex.RLock()
	defer dg.mutex.RUnlock()

	if node, exists := dg.Nodes[nodeName]; exists {
		return append([]string{}, node.Dependencies...), nil
	}
	return nil, fmt.Errorf("node '%s' not found", nodeName)
}

// GetDependents returns all dependents for a given node
func (dg *DependencyGraph) GetDependents(nodeName string) ([]string, error) {
	dg.mutex.RLock()
	defer dg.mutex.RUnlock()

	if node, exists := dg.Nodes[nodeName]; exists {
		return append([]string{}, node.Dependents...), nil
	}
	return nil, fmt.Errorf("node '%s' not found", nodeName)
}

// GetNode returns a node by name
func (dg *DependencyGraph) GetNode(nodeName string) (*DependencyNode, bool) {
	dg.mutex.RLock()
	defer dg.mutex.RUnlock()

	node, exists := dg.Nodes[nodeName]
	return node, exists
}

// GetAllNodes returns all nodes in the graph
func (dg *DependencyGraph) GetAllNodes() map[string]*DependencyNode {
	dg.mutex.RLock()
	defer dg.mutex.RUnlock()

	result := make(map[string]*DependencyNode)
	for name, node := range dg.Nodes {
		result[name] = node
	}
	return result
}

// RemoveNode removes a node and all its dependencies
func (dg *DependencyGraph) RemoveNode(nodeName string) error {
	dg.mutex.Lock()
	defer dg.mutex.Unlock()

	if _, exists := dg.Nodes[nodeName]; !exists {
		return fmt.Errorf("node '%s' not found", nodeName)
	}

	// Remove from dependents' dependency lists
	for _, dependent := range dg.Nodes[nodeName].Dependents {
		deps := dg.Nodes[dependent].Dependencies
		newDeps := []string{}
		for _, dep := range deps {
			if dep != nodeName {
				newDeps = append(newDeps, dep)
			}
		}
		dg.Nodes[dependent].Dependencies = newDeps
	}

	// Remove from edges
	delete(dg.Edges, nodeName)

	// Remove the node
	delete(dg.Nodes, nodeName)

	return nil
}

// ValidateDependencies validates that all dependencies exist
func (dg *DependencyGraph) ValidateDependencies() []string {
	dg.mutex.RLock()
	defer dg.mutex.RUnlock()

	missing := []string{}

	for nodeName, node := range dg.Nodes {
		for _, dep := range node.Dependencies {
			if _, exists := dg.Nodes[dep]; !exists {
				missing = append(missing, fmt.Sprintf("%s -> %s (missing)", nodeName, dep))
			}
		}
	}

	return missing
}

// GetDependencyChain returns the dependency chain for a given node
func (dg *DependencyGraph) GetDependencyChain(nodeName string) ([]string, error) {
	dg.mutex.RLock()
	defer dg.mutex.RUnlock()

	if _, exists := dg.Nodes[nodeName]; !exists {
		return nil, fmt.Errorf("node '%s' not found", nodeName)
	}

	visited := make(map[string]bool)
	chain := []string{}

	dg.buildDependencyChain(nodeName, visited, &chain)

	return chain, nil
}

// buildDependencyChain recursively builds the dependency chain
func (dg *DependencyGraph) buildDependencyChain(nodeName string, visited map[string]bool, chain *[]string) {
	if visited[nodeName] {
		return
	}

	visited[nodeName] = true

	// Add dependencies first (depth-first)
	for _, dep := range dg.Nodes[nodeName].Dependencies {
		dg.buildDependencyChain(dep, visited, chain)
	}

	// Add current node
	*chain = append(*chain, nodeName)
}

// GetDependencyStats returns statistics about the dependency graph
func (dg *DependencyGraph) GetDependencyStats() map[string]interface{} {
	dg.mutex.RLock()
	defer dg.mutex.RUnlock()

	stats := map[string]interface{}{
		"total_nodes":     len(dg.Nodes),
		"total_edges":     0,
		"max_dependencies": 0,
		"max_dependents":   0,
		"isolated_nodes":   0,
		"leaf_nodes":       0,
		"root_nodes":       0,
	}

	maxDeps := 0
	maxDependents := 0
	isolated := 0
	leaves := 0
	roots := 0

	for _, node := range dg.Nodes {
		stats["total_edges"] = stats["total_edges"].(int) + len(node.Dependencies)

		if len(node.Dependencies) > maxDeps {
			maxDeps = len(node.Dependencies)
		}
		if len(node.Dependents) > maxDependents {
			maxDependents = len(node.Dependents)
		}

		if len(node.Dependencies) == 0 && len(node.Dependents) == 0 {
			isolated++
		}
		if len(node.Dependents) == 0 {
			leaves++
		}
		if len(node.Dependencies) == 0 {
			roots++
		}
	}

	stats["max_dependencies"] = maxDeps
	stats["max_dependents"] = maxDependents
	stats["isolated_nodes"] = isolated
	stats["leaf_nodes"] = leaves
	stats["root_nodes"] = roots

	return stats
}

// VisualizeDependencies returns a string representation of the dependency graph
func (dg *DependencyGraph) VisualizeDependencies() string {
	dg.mutex.RLock()
	defer dg.mutex.RUnlock()

	var result strings.Builder
	result.WriteString("Dependency Graph:\n")

	// Sort nodes for consistent output
	nodeNames := make([]string, 0, len(dg.Nodes))
	for name := range dg.Nodes {
		nodeNames = append(nodeNames, name)
	}
	sort.Strings(nodeNames)

	for _, name := range nodeNames {
		node := dg.Nodes[name]
		result.WriteString(fmt.Sprintf("  %s (%s:%d)\n", name, node.File, node.Line))
		
		if len(node.Dependencies) > 0 {
			result.WriteString("    Dependencies:\n")
			for _, dep := range node.Dependencies {
				result.WriteString(fmt.Sprintf("      -> %s\n", dep))
			}
		}
		
		if len(node.Dependents) > 0 {
			result.WriteString("    Dependents:\n")
			for _, dep := range node.Dependents {
				result.WriteString(fmt.Sprintf("      <- %s\n", dep))
			}
		}
		result.WriteString("\n")
	}

	return result.String()
}

// Clear removes all nodes and edges from the graph
func (dg *DependencyGraph) Clear() {
	dg.mutex.Lock()
	defer dg.mutex.Unlock()

	dg.Nodes = make(map[string]*DependencyNode)
	dg.Edges = make(map[string][]string)
} 