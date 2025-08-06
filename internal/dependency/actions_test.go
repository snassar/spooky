package dependency

import (
	"strings"
	"testing"
)

func TestNewActionDependencyManager(t *testing.T) {
	adm := NewActionDependencyManager()

	if adm == nil {
		t.Fatal("NewActionDependencyManager returned nil")
	}

	if adm.graph == nil {
		t.Error("Dependency graph should not be nil")
	}
}

func TestAddAction(t *testing.T) {
	adm := NewActionDependencyManager()

	// Test adding a simple action
	err := adm.AddAction("action1", []string{}, "actions.hcl", 1)
	if err != nil {
		t.Fatalf("Failed to add action: %v", err)
	}

	// Test adding action with dependencies
	err = adm.AddAction("action2", []string{"action1"}, "actions.hcl", 2)
	if err != nil {
		t.Fatalf("Failed to add action with dependency: %v", err)
	}

	// Test adding action with non-existent dependency
	err = adm.AddAction("action3", []string{"missing_action"}, "actions.hcl", 3)
	if err == nil {
		t.Error("Expected error for non-existent dependency")
	}
}

func TestAddActionCollection(t *testing.T) {
	adm := NewActionDependencyManager()

	// Create action collection
	actions := map[string]ActionInfo{
		"action1": {
			Dependencies: []string{},
			SourceFile:   "actions.hcl",
			SourceLine:   1,
		},
		"action2": {
			Dependencies: []string{"action1"},
			SourceFile:   "actions.hcl",
			SourceLine:   2,
		},
		"action3": {
			Dependencies: []string{"action1", "action2"},
			SourceFile:   "actions.hcl",
			SourceLine:   3,
		},
	}

	err := adm.AddActionCollection(actions)
	if err != nil {
		t.Fatalf("Failed to add action collection: %v", err)
	}

	// Verify all actions were added
	if len(adm.graph.Nodes) != 3 {
		t.Errorf("Expected 3 nodes, got %d", len(adm.graph.Nodes))
	}
}

func TestValidateDependencies(t *testing.T) {
	adm := NewActionDependencyManager()

	// Test valid dependencies
	actions := map[string]ActionInfo{
		"action1": {
			Dependencies: []string{},
			SourceFile:   "actions.hcl",
			SourceLine:   1,
		},
		"action2": {
			Dependencies: []string{"action1"},
			SourceFile:   "actions.hcl",
			SourceLine:   2,
		},
	}

	adm.AddActionCollection(actions)
	result := adm.ValidateDependencies()

	if !result.Valid {
		t.Errorf("Expected valid dependencies, got errors: %v", result.Errors)
	}
}

func TestValidateDependenciesWithCircularRefs(t *testing.T) {
	adm := NewActionDependencyManager()

	// Test circular dependencies
	actions := map[string]ActionInfo{
		"action1": {
			Dependencies: []string{"action2"},
			SourceFile:   "actions.hcl",
			SourceLine:   1,
		},
		"action2": {
			Dependencies: []string{"action1"},
			SourceFile:   "actions.hcl",
			SourceLine:   2,
		},
	}

	adm.AddActionCollection(actions)
	result := adm.ValidateDependencies()

	if result.Valid {
		t.Error("Expected invalid dependencies due to circular reference")
	}

	if len(result.CircularRefs) == 0 {
		t.Error("Expected circular reference to be detected")
	}

	foundCircularError := false
	for _, err := range result.Errors {
		if err.Type == "circular_reference" {
			foundCircularError = true
			break
		}
	}
	if !foundCircularError {
		t.Error("Expected circular reference error")
	}
}

func TestValidateDependenciesWithMissingRefs(t *testing.T) {
	adm := NewActionDependencyManager()

	// Test missing dependencies - add action with missing dependency directly to graph
	adm.graph.AddNode("action1", "action", "actions.hcl", 1)
	adm.graph.Nodes["action1"].Dependencies = []string{"missing_action"}

	result := adm.ValidateDependencies()

	if result.Valid {
		t.Error("Expected invalid dependencies due to missing reference")
	}

	if len(result.MissingRefs) == 0 {
		t.Error("Expected missing reference to be detected")
	}

	foundMissingError := false
	for _, err := range result.Errors {
		if err.Type == "missing_reference" {
			foundMissingError = true
			break
		}
	}
	if !foundMissingError {
		t.Error("Expected missing reference error")
	}
}

func TestGetExecutionOrder(t *testing.T) {
	adm := NewActionDependencyManager()

	// Create valid dependency chain
	actions := map[string]ActionInfo{
		"action1": {
			Dependencies: []string{},
			SourceFile:   "actions.hcl",
			SourceLine:   1,
		},
		"action2": {
			Dependencies: []string{"action1"},
			SourceFile:   "actions.hcl",
			SourceLine:   2,
		},
		"action3": {
			Dependencies: []string{"action1", "action2"},
			SourceFile:   "actions.hcl",
			SourceLine:   3,
		},
	}

	adm.AddActionCollection(actions)
	order, err := adm.GetExecutionOrder()
	if err != nil {
		t.Fatalf("Failed to get execution order: %v", err)
	}

	if len(order) != 3 {
		t.Errorf("Expected 3 actions in order, got %d", len(order))
	}

	// action1 should come first (no dependencies)
	if order[0] != "action1" {
		t.Errorf("Expected action1 first, got %s", order[0])
	}

	// action2 should come before action3
	action2Index := -1
	action3Index := -1
	for i, action := range order {
		if action == "action2" {
			action2Index = i
		}
		if action == "action3" {
			action3Index = i
		}
	}

	if action2Index == -1 || action3Index == -1 {
		t.Error("Expected to find both action2 and action3 in order")
	}

	if action2Index >= action3Index {
		t.Error("Expected action2 to come before action3")
	}
}

func TestGetExecutionOrderWithCircularDeps(t *testing.T) {
	adm := NewActionDependencyManager()

	// Create circular dependencies
	actions := map[string]ActionInfo{
		"action1": {
			Dependencies: []string{"action2"},
			SourceFile:   "actions.hcl",
			SourceLine:   1,
		},
		"action2": {
			Dependencies: []string{"action1"},
			SourceFile:   "actions.hcl",
			SourceLine:   2,
		},
	}

	adm.AddActionCollection(actions)
	_, err := adm.GetExecutionOrder()
	if err == nil {
		t.Error("Expected error for circular dependencies")
	}

	if !strings.Contains(err.Error(), "dependency validation failed") {
		t.Errorf("Expected dependency validation error, got: %v", err)
	}
}

func TestGetParallelGroups(t *testing.T) {
	adm := NewActionDependencyManager()

	// Create actions with dependencies
	actions := map[string]ActionInfo{
		"action1": {
			Dependencies: []string{},
			SourceFile:   "actions.hcl",
			SourceLine:   1,
		},
		"action2": {
			Dependencies: []string{},
			SourceFile:   "actions.hcl",
			SourceLine:   2,
		},
		"action3": {
			Dependencies: []string{"action1", "action2"},
			SourceFile:   "actions.hcl",
			SourceLine:   3,
		},
		"action4": {
			Dependencies: []string{"action1"},
			SourceFile:   "actions.hcl",
			SourceLine:   4,
		},
	}

	adm.AddActionCollection(actions)
	groups, err := adm.GetParallelGroups()
	if err != nil {
		t.Fatalf("Failed to get parallel groups: %v", err)
	}

	if len(groups) < 2 {
		t.Errorf("Expected at least 2 parallel groups, got %d", len(groups))
	}

	// First group should contain action1 and action2 (no dependencies)
	foundAction1 := false
	foundAction2 := false
	for _, action := range groups[0] {
		if action == "action1" {
			foundAction1 = true
		}
		if action == "action2" {
			foundAction2 = true
		}
	}

	if !foundAction1 || !foundAction2 {
		t.Error("Expected action1 and action2 in first parallel group")
	}
}

func TestGetDependencyChain(t *testing.T) {
	adm := NewActionDependencyManager()

	// Create dependency chain
	actions := map[string]ActionInfo{
		"action1": {
			Dependencies: []string{},
			SourceFile:   "actions.hcl",
			SourceLine:   1,
		},
		"action2": {
			Dependencies: []string{"action1"},
			SourceFile:   "actions.hcl",
			SourceLine:   2,
		},
		"action3": {
			Dependencies: []string{"action2"},
			SourceFile:   "actions.hcl",
			SourceLine:   3,
		},
	}

	adm.AddActionCollection(actions)
	chain, err := adm.GetDependencyChain("action3")
	if err != nil {
		t.Fatalf("Failed to get dependency chain: %v", err)
	}

	// Should be: action1, action2, action3
	expected := []string{"action1", "action2", "action3"}
	if len(chain) != len(expected) {
		t.Errorf("Expected chain length %d, got %d", len(expected), len(chain))
	}

	for i, expectedAction := range expected {
		if chain[i] != expectedAction {
			t.Errorf("Expected %s at position %d, got %s", expectedAction, i, chain[i])
		}
	}
}

func TestGetDependents(t *testing.T) {
	adm := NewActionDependencyManager()

	// Create dependency relationships
	actions := map[string]ActionInfo{
		"action1": {
			Dependencies: []string{},
			SourceFile:   "actions.hcl",
			SourceLine:   1,
		},
		"action2": {
			Dependencies: []string{"action1"},
			SourceFile:   "actions.hcl",
			SourceLine:   2,
		},
		"action3": {
			Dependencies: []string{"action1"},
			SourceFile:   "actions.hcl",
			SourceLine:   3,
		},
	}

	adm.AddActionCollection(actions)
	dependents, err := adm.GetDependents("action1")
	if err != nil {
		t.Fatalf("Failed to get dependents: %v", err)
	}

	if len(dependents) != 2 {
		t.Errorf("Expected 2 dependents, got %d", len(dependents))
	}

	// Check that both action2 and action3 are dependents
	foundAction2 := false
	foundAction3 := false
	for _, dep := range dependents {
		if dep == "action2" {
			foundAction2 = true
		}
		if dep == "action3" {
			foundAction3 = true
		}
	}

	if !foundAction2 || !foundAction3 {
		t.Error("Expected action2 and action3 to be dependents of action1")
	}
}

func TestGetDependencies(t *testing.T) {
	adm := NewActionDependencyManager()

	// Create dependency relationships
	actions := map[string]ActionInfo{
		"action1": {
			Dependencies: []string{},
			SourceFile:   "actions.hcl",
			SourceLine:   1,
		},
		"action2": {
			Dependencies: []string{"action1"},
			SourceFile:   "actions.hcl",
			SourceLine:   2,
		},
		"action3": {
			Dependencies: []string{"action1", "action2"},
			SourceFile:   "actions.hcl",
			SourceLine:   3,
		},
	}

	adm.AddActionCollection(actions)
	dependencies, err := adm.GetDependencies("action3")
	if err != nil {
		t.Fatalf("Failed to get dependencies: %v", err)
	}

	if len(dependencies) != 2 {
		t.Errorf("Expected 2 dependencies, got %d", len(dependencies))
	}

	// Check that both action1 and action2 are dependencies
	foundAction1 := false
	foundAction2 := false
	for _, dep := range dependencies {
		if dep == "action1" {
			foundAction1 = true
		}
		if dep == "action2" {
			foundAction2 = true
		}
	}

	if !foundAction1 || !foundAction2 {
		t.Error("Expected action1 and action2 to be dependencies of action3")
	}
}

func TestGetDependencyStats(t *testing.T) {
	adm := NewActionDependencyManager()

	// Create actions with dependencies
	actions := map[string]ActionInfo{
		"action1": {
			Dependencies: []string{},
			SourceFile:   "actions.hcl",
			SourceLine:   1,
		},
		"action2": {
			Dependencies: []string{"action1"},
			SourceFile:   "actions.hcl",
			SourceLine:   2,
		},
		"action3": {
			Dependencies: []string{"action1", "action2"},
			SourceFile:   "actions.hcl",
			SourceLine:   3,
		},
	}

	adm.AddActionCollection(actions)
	stats := adm.GetDependencyStats()

	if stats["total_actions"] != 3 {
		t.Errorf("Expected 3 total actions, got %v", stats["total_actions"])
	}

	if stats["total_nodes"] != 3 {
		t.Errorf("Expected 3 total nodes, got %v", stats["total_nodes"])
	}

	if stats["total_edges"] != 3 {
		t.Errorf("Expected 3 total edges, got %v", stats["total_edges"])
	}
}

func TestVisualizeDependencies(t *testing.T) {
	adm := NewActionDependencyManager()

	// Create simple dependency structure
	actions := map[string]ActionInfo{
		"action1": {
			Dependencies: []string{},
			SourceFile:   "actions.hcl",
			SourceLine:   1,
		},
		"action2": {
			Dependencies: []string{"action1"},
			SourceFile:   "actions.hcl",
			SourceLine:   2,
		},
	}

	adm.AddActionCollection(actions)
	visualization := adm.VisualizeDependencies()

	if visualization == "" {
		t.Error("Expected non-empty visualization")
	}

	if !strings.Contains(visualization, "action1") {
		t.Error("Expected visualization to contain action1")
	}

	if !strings.Contains(visualization, "action2") {
		t.Error("Expected visualization to contain action2")
	}

	if !strings.Contains(visualization, "Dependency Graph:") {
		t.Error("Expected visualization to contain header")
	}
}

func TestDetectMissingReferences(t *testing.T) {
	adm := NewActionDependencyManager()

	// Add action with missing dependency
	adm.graph.AddNode("action1", "action", "actions.hcl", 1)
	adm.graph.Nodes["action1"].Dependencies = []string{"missing_action"}

	missing := adm.detectMissingReferences()
	if len(missing) == 0 {
		t.Error("Expected missing reference to be detected")
	}

	if !strings.Contains(missing[0], "action1 -> missing_action") {
		t.Errorf("Expected missing reference format, got: %s", missing[0])
	}
}

func TestDetectSelfReferences(t *testing.T) {
	adm := NewActionDependencyManager()

	// Add action with self-reference
	adm.graph.AddNode("action1", "action", "actions.hcl", 1)
	adm.graph.Nodes["action1"].Dependencies = []string{"action1"}

	selfRefs := adm.detectSelfReferences()
	if len(selfRefs) == 0 {
		t.Error("Expected self-reference to be detected")
	}

	if selfRefs[0] != "action1" {
		t.Errorf("Expected self-reference for action1, got: %s", selfRefs[0])
	}
}

func TestDetectOrphanedActions(t *testing.T) {
	adm := NewActionDependencyManager()

	// Add orphaned action (no dependents)
	adm.graph.AddNode("action1", "action", "actions.hcl", 1)

	// Add action with dependents
	adm.graph.AddNode("action2", "action", "actions.hcl", 2)
	adm.graph.AddNode("action3", "action", "actions.hcl", 3)
	adm.graph.Nodes["action3"].Dependencies = []string{"action2"}
	adm.graph.Nodes["action2"].Dependents = []string{"action3"}

	orphaned := adm.detectOrphanedActions()
	if len(orphaned) == 0 {
		t.Error("Expected orphaned action to be detected")
	}

	foundAction1 := false
	for _, action := range orphaned {
		if action == "action1" {
			foundAction1 = true
			break
		}
	}

	if !foundAction1 {
		t.Error("Expected action1 to be detected as orphaned")
	}
}

func TestAllDependenciesProcessed(t *testing.T) {
	adm := NewActionDependencyManager()

	// Add actions with dependencies
	adm.graph.AddNode("action1", "action", "actions.hcl", 1)
	adm.graph.AddNode("action2", "action", "actions.hcl", 2)
	adm.graph.Nodes["action2"].Dependencies = []string{"action1"}

	processed := make(map[string]bool)

	// action1 should be ready (no dependencies)
	if !adm.allDependenciesProcessed("action1", processed) {
		t.Error("Expected action1 to be ready for processing")
	}

	// action2 should not be ready (action1 not processed)
	if adm.allDependenciesProcessed("action2", processed) {
		t.Error("Expected action2 to not be ready for processing")
	}

	// Mark action1 as processed
	processed["action1"] = true

	// action2 should now be ready
	if !adm.allDependenciesProcessed("action2", processed) {
		t.Error("Expected action2 to be ready for processing after action1")
	}
}
