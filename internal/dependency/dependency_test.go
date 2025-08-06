package dependency

import (
	"strings"
	"testing"

	"spooky/internal/logging"
)

func TestDependencyGraph(t *testing.T) {
	graph := NewDependencyGraph()

	// Test adding nodes
	t.Run("AddNode", func(t *testing.T) {
		node := graph.AddNode("var1", "variable", "variables.hcl", 1)
		if node == nil {
			t.Fatal("Expected node to be created")
		}
		if node.Name != "var1" {
			t.Errorf("Expected name 'var1', got '%s'", node.Name)
		}
		if node.Type != "variable" {
			t.Errorf("Expected type 'variable', got '%s'", node.Type)
		}
	})

	// Test adding dependencies
	t.Run("AddDependency", func(t *testing.T) {
		graph.AddNode("var1", "variable", "variables.hcl", 1)
		graph.AddNode("var2", "variable", "variables.hcl", 2)

		err := graph.AddDependency("var1", "var2")
		if err != nil {
			t.Fatalf("Failed to add dependency: %v", err)
		}

		// Check that dependency was added
		deps, err := graph.GetDependencies("var1")
		if err != nil {
			t.Fatalf("Failed to get dependencies: %v", err)
		}
		if len(deps) != 1 || deps[0] != "var2" {
			t.Errorf("Expected dependency 'var2', got %v", deps)
		}

		// Check that dependent was added
		dependents, err := graph.GetDependents("var2")
		if err != nil {
			t.Fatalf("Failed to get dependents: %v", err)
		}
		if len(dependents) != 1 || dependents[0] != "var1" {
			t.Errorf("Expected dependent 'var1', got %v", dependents)
		}
	})

	// Test circular dependency detection
	t.Run("CircularDependency", func(t *testing.T) {
		graph.Clear()
		graph.AddNode("var1", "variable", "variables.hcl", 1)
		graph.AddNode("var2", "variable", "variables.hcl", 2)
		graph.AddNode("var3", "variable", "variables.hcl", 3)

		// Create circular dependency: var1 -> var2 -> var3 -> var1
		graph.AddDependency("var1", "var2")
		graph.AddDependency("var2", "var3")
		graph.AddDependency("var3", "var1")

		circular := graph.DetectCircularRefs()
		if circular == nil {
			t.Fatal("Expected circular dependency to be detected")
		}
		if len(circular) == 0 {
			t.Fatal("Expected circular dependency path to be returned")
		}
	})

	// Test self-reference detection
	t.Run("SelfReference", func(t *testing.T) {
		graph.Clear()
		graph.AddNode("var1", "variable", "variables.hcl", 1)

		err := graph.AddDependency("var1", "var1")
		if err == nil {
			t.Fatal("Expected error for self-reference")
		}
		if !strings.Contains(err.Error(), "self-reference") {
			t.Errorf("Expected self-reference error, got: %v", err)
		}
	})

	// Test resolution order
	t.Run("ResolutionOrder", func(t *testing.T) {
		graph.Clear()
		graph.AddNode("var1", "variable", "variables.hcl", 1)
		graph.AddNode("var2", "variable", "variables.hcl", 2)
		graph.AddNode("var3", "variable", "variables.hcl", 3)

		// var1 depends on var2, var2 depends on var3
		graph.AddDependency("var1", "var2")
		graph.AddDependency("var2", "var3")

		order, err := graph.ResolveOrder()
		if err != nil {
			t.Fatalf("Failed to get resolution order: %v", err)
		}

		// var3 should come first (no dependencies)
		if order[0] != "var3" {
			t.Errorf("Expected var3 first, got %s", order[0])
		}
		// var2 should come second
		if order[1] != "var2" {
			t.Errorf("Expected var2 second, got %s", order[1])
		}
		// var1 should come last
		if order[2] != "var1" {
			t.Errorf("Expected var1 last, got %s", order[2])
		}
	})

	// Test dependency chain
	t.Run("DependencyChain", func(t *testing.T) {
		graph.Clear()
		graph.AddNode("var1", "variable", "variables.hcl", 1)
		graph.AddNode("var2", "variable", "variables.hcl", 2)
		graph.AddNode("var3", "variable", "variables.hcl", 3)

		graph.AddDependency("var1", "var2")
		graph.AddDependency("var2", "var3")

		chain, err := graph.GetDependencyChain("var1")
		if err != nil {
			t.Fatalf("Failed to get dependency chain: %v", err)
		}

		// Should be: var3, var2, var1
		expected := []string{"var3", "var2", "var1"}
		if len(chain) != len(expected) {
			t.Errorf("Expected chain length %d, got %d", len(expected), len(chain))
		}
		for i, expectedVar := range expected {
			if chain[i] != expectedVar {
				t.Errorf("Expected %s at position %d, got %s", expectedVar, i, chain[i])
			}
		}
	})

	// Test statistics
	t.Run("Statistics", func(t *testing.T) {
		graph.Clear()
		graph.AddNode("var1", "variable", "variables.hcl", 1)
		graph.AddNode("var2", "variable", "variables.hcl", 2)
		graph.AddNode("var3", "variable", "variables.hcl", 3)

		graph.AddDependency("var1", "var2")
		graph.AddDependency("var2", "var3")

		stats := graph.GetDependencyStats()
		if stats["total_nodes"] != 3 {
			t.Errorf("Expected 3 nodes, got %v", stats["total_nodes"])
		}
		if stats["total_edges"] != 2 {
			t.Errorf("Expected 2 edges, got %v", stats["total_edges"])
		}
		if stats["root_nodes"] != 1 {
			t.Errorf("Expected 1 root node, got %v", stats["root_nodes"])
		}
		if stats["leaf_nodes"] != 1 {
			t.Errorf("Expected 1 leaf node, got %v", stats["leaf_nodes"])
		}
	})
}

func TestDependencyManager(t *testing.T) {
	logger := logging.GetLogger()
	manager := NewDependencyManager(logger)

	// Test adding variables
	t.Run("AddVariable", func(t *testing.T) {
		// Add the dependency first
		err := manager.AddVariable("var2", "variables.hcl", 2, []string{})
		if err != nil {
			t.Fatalf("Failed to add dependency variable: %v", err)
		}

		// Then add the variable that depends on it
		err = manager.AddVariable("var1", "variables.hcl", 1, []string{"var2"})
		if err != nil {
			t.Fatalf("Failed to add variable: %v", err)
		}
	})

	// Test validation
	t.Run("Validation", func(t *testing.T) {
		manager.Clear()

		// Manually create a graph with missing dependencies
		// Add a node with a dependency that doesn't exist
		node := manager.graph.AddNode("var1", "variable", "variables.hcl", 1)
		node.Dependencies = []string{"missing_var"}
		// Update the edges map
		manager.graph.Edges["var1"] = []string{"missing_var"}

		errors := manager.ValidateDependencies()
		if len(errors) == 0 {
			t.Fatal("Expected validation errors for missing dependency")
		}

		foundMissing := false
		for _, err := range errors {
			if err.Type == "missing" {
				foundMissing = true
				break
			}
		}
		if !foundMissing {
			t.Error("Expected missing dependency error")
		}
	})

	// Test circular dependency validation
	t.Run("CircularDependencyValidation", func(t *testing.T) {
		manager.Clear()

		// Create circular dependency by adding nodes first, then dependencies
		manager.graph.AddNode("var1", "variable", "variables.hcl", 1)
		manager.graph.AddNode("var2", "variable", "variables.hcl", 2)
		manager.graph.AddNode("var3", "variable", "variables.hcl", 3)

		// Create circular dependency: var1 -> var2 -> var3 -> var1
		manager.graph.AddDependency("var1", "var2")
		manager.graph.AddDependency("var2", "var3")
		manager.graph.AddDependency("var3", "var1")

		errors := manager.ValidateDependencies()
		if len(errors) == 0 {
			t.Fatal("Expected validation errors for circular dependency")
		}

		foundCircular := false
		for _, err := range errors {
			if err.Type == "circular" {
				foundCircular = true
				break
			}
		}
		if !foundCircular {
			t.Error("Expected circular dependency error")
		}
	})

	// Test resolution order
	t.Run("ResolutionOrder", func(t *testing.T) {
		manager.Clear()

		manager.AddVariable("var3", "variables.hcl", 3, []string{})
		manager.AddVariable("var2", "variables.hcl", 2, []string{"var3"})
		manager.AddVariable("var1", "variables.hcl", 1, []string{"var2"})

		order, err := manager.GetResolutionOrder()
		if err != nil {
			t.Fatalf("Failed to get resolution order: %v", err)
		}

		// Should be: var3, var2, var1
		expected := []string{"var3", "var2", "var1"}
		if len(order) != len(expected) {
			t.Errorf("Expected order length %d, got %d", len(expected), len(order))
		}
		for i, expectedVar := range expected {
			if order[i] != expectedVar {
				t.Errorf("Expected %s at position %d, got %s", expectedVar, i, order[i])
			}
		}
	})

	// Test impact analysis
	t.Run("ImpactAnalysis", func(t *testing.T) {
		manager.Clear()

		manager.AddVariable("var3", "variables.hcl", 3, []string{})
		manager.AddVariable("var2", "variables.hcl", 2, []string{"var3"})
		manager.AddVariable("var1", "variables.hcl", 1, []string{"var2"})

		analysis, err := manager.GetImpactAnalysis("var3")
		if err != nil {
			t.Fatalf("Failed to get impact analysis: %v", err)
		}

		impactedCount := analysis["impacted_count"].(int)
		if impactedCount != 3 {
			t.Errorf("Expected 3 impacted variables, got %d", impactedCount)
		}

		impactedVars := analysis["impacted_vars"].([]string)
		expected := []string{"var3", "var2", "var1"}
		if len(impactedVars) != len(expected) {
			t.Errorf("Expected %d impacted variables, got %d", len(expected), len(impactedVars))
		}
	})

	// Test dependency report
	t.Run("DependencyReport", func(t *testing.T) {
		manager.Clear()

		manager.AddVariable("var1", "variables.hcl", 1, []string{"var2"})
		manager.AddVariable("var2", "variables.hcl", 2, []string{})

		report := manager.GetDependencyReport()
		if report == nil {
			t.Fatal("Expected dependency report")
		}

		stats := report["statistics"].(map[string]interface{})
		if stats["total_nodes"] != 2 {
			t.Errorf("Expected 2 nodes in report, got %v", stats["total_nodes"])
		}

		validation := report["validation"].(map[string]interface{})
		if !validation["valid"].(bool) {
			t.Error("Expected validation to be valid")
		}
	})

	// Test caching
	t.Run("Caching", func(t *testing.T) {
		manager.Clear()

		// Cache a result
		manager.CacheResult("test_key", "test_value")

		// Retrieve cached result
		value, exists := manager.GetCachedResult("test_key")
		if !exists {
			t.Fatal("Expected cached result to exist")
		}
		if value != "test_value" {
			t.Errorf("Expected 'test_value', got %v", value)
		}

		// Test non-existent key
		_, exists = manager.GetCachedResult("non_existent")
		if exists {
			t.Error("Expected non-existent key to not exist")
		}

		// Clear cache
		manager.ClearCache()
		_, exists = manager.GetCachedResult("test_key")
		if exists {
			t.Error("Expected cached result to be cleared")
		}
	})
}

func TestDependencyError(t *testing.T) {
	// Test error formatting
	t.Run("ErrorFormatting", func(t *testing.T) {
		err := &DependencyError{
			Type:    "circular",
			Message: "Circular dependency detected",
			Nodes:   []string{"var1", "var2", "var3"},
		}

		errorStr := err.Error()
		if errorStr != "Circular dependency detected" {
			t.Errorf("Expected error message, got: %s", errorStr)
		}
	})
} 