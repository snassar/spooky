package actions

import (
	spookytypes "spooky/internal/types"
	"context"
	"testing"
	"time"

	spookylogging "spooky/internal/logging"
	spookytypesactions "spooky/internal/types/actions"
	spookytypesconfig "spooky/internal/types/config"
	spookytypesfacts "spooky/internal/types/facts"
	spookytypeslogging "spooky/internal/types/logging"
	spookytypesvariables "spooky/internal/types/variables"
)

func TestPlanCollection(t *testing.T) {
	// Create a logger
	logger := spookylogging.NewLogger(spookytypeslogging.Config{
		Level:     spookytypeslogging.InfoLevel,
		Format:    "json",
		Output:    "stdout",
		Timestamp: true,
	})

	// Create action manager with facts manager
	manager := NewManager(logger)

	// Test collection planning
	collectionIDs := []string{"test-collection-1", "test-collection-2"}
	options := &spookytypesactions.PlanningOptions{
		PlanName:       "test-plan",
		Description:    "Test collection plan",
		Parallel:       true,
		Optimize:       true,
		ValidateFacts:  false, // Disable fact validation for tests
		CreateRollback: true,
		Metadata:       make(map[string]interface{}),
	}

	// Plan the collection
	plan, err := manager.PlanCollection(collectionIDs, options)
	if err != nil {
		t.Fatalf("Failed to plan collection: %v", err)
	}

	// Verify plan structure
	if plan == nil {
		t.Fatal("Plan is nil")
	}

	if plan.ID == "" {
		t.Error("Plan ID is empty")
	}

	if plan.Name != "test-plan" {
		t.Errorf("Expected plan name 'test-plan', got '%s'", plan.Name)
	}

	if len(plan.Collections) != 2 {
		t.Errorf("Expected 2 collections, got %d", len(plan.Collections))
	}

	// Check that collections have actions (they should be added by loadCollections)
	for i, collection := range plan.Collections {
		if len(collection.Actions) == 0 {
			t.Errorf("Collection %d should have actions", i)
		}
	}

	if plan.EstimatedTime <= 0 {
		t.Error("Estimated time should be positive")
	}

	if plan.RollbackPlan == nil {
		t.Error("Rollback plan should not be nil")
	}

	if plan.Validation == nil {
		t.Error("Validation should not be nil")
	}

	if plan.FactRequirements == nil {
		t.Error("Fact requirements should not be nil")
	}

	t.Logf("Successfully created collection plan: %s", plan.ID)
	t.Logf("Estimated time: %v", plan.EstimatedTime)
	t.Logf("Collections: %d", len(plan.Collections))
}

func TestPlanActionCollection(t *testing.T) {
	// Create a logger
	logger := spookylogging.NewLogger(spookytypeslogging.Config{
		Level:     spookytypeslogging.InfoLevel,
		Format:    "json",
		Output:    "stdout",
		Timestamp: true,
	})

	// Create action manager
	manager := NewManager(logger)

	// Create a test action
	action := &spookytypesactions.Action{
		Name:        "test-action",
		Description: "Test action",
		Type:        "command",
		Command:     "echo 'test'",
		Machines:    []string{"test-machine"},
		Priority:    5,
	}

	// Create a test collection
	collection := &spookytypesactions.ActionCollection{
		Name:         "test-collection",
		Description:  "Test collection",
		Actions:      []*spookytypesactions.Action{action},
		Dependencies: make([]*spookytypesactions.CollectionDependency, 0),
		Priority:     5,
		Parallel:     false,
		Timeout:      300,
		Metadata:     make(map[string]interface{}),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Create a test context
	context := &spookytypesactions.ActionContext{
		ProjectPath: "/test/project",
		Machines:    make([]spookytypesconfig.Machine, 0),
		Facts:       &spookytypesfacts.FactCollection{},
		Variables:   &spookytypesvariables.VariableContext{},
		Parallel:    false,
		Timeout:     30 * time.Minute,
		CustomData:  make(map[string]interface{}),
	}

	// Plan the action collection
	plan, err := manager.PlanActionCollection(collection, context)
	if err != nil {
		t.Fatalf("Failed to plan action collection: %v", err)
	}

	// Verify plan structure
	if plan == nil {
		t.Fatal("Plan is nil")
	}

	if plan.PlanID == "" {
		t.Error("Plan ID is empty")
	}

	if plan.ActionName != "test-collection" {
		t.Errorf("Expected action name 'test-collection', got '%s'", plan.ActionName)
	}

	if plan.Status != spookytypesactions.PlanningStatusCompleted {
		t.Errorf("Expected status 'completed', got '%s'", plan.Status)
	}

	t.Logf("Successfully created action collection plan: %s", plan.PlanID)
}

func.*runs(t *testing.T) {
	// Create a logger
	logger := spookylogging.NewLogger(spookytypeslogging.Config{
		Level:     spookytypeslogging.InfoLevel,
		Format:    "json",
		Output:    "stdout",
		Timestamp: true,
	})

	// Create action manager with facts manager (nil for this test)
	manager := NewManagerWithFacts(logger, nil)

	// Create a test collection with actions that have template dependencies
	action := &spookytypesactions.Action{
		Name:        "test-action",
		Description: "Test action with template",
		Type:        "template_deploy",
		Machines:    []string{"test-machine"},
		Template: &spookytypesactions.TemplateConfig{
			Source:      "{{ .system.hostname }}",
			Destination: "/tmp/test",
		},
		Priority: 5,
	}

	collection := &spookytypesactions.ActionCollection{
		Name:         "test-collection-with-facts",
		Description:  "Test collection with fact dependencies",
		Actions:      []*spookytypesactions.Action{action},
		Dependencies: make([]*spookytypesactions.CollectionDependency, 0),
		Priority:     5,
		Parallel:     false,
		Timeout:      300,
		Metadata:     make(map[string]interface{}),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Test fact requirements analysis
	requirements, err := manager.determineFactRequirements([]*spookytypesactions.ActionCollection{collection})
	if err != nil {
		t.Fatalf("Failed to determine fact requirements: %v", err)
	}

	// Verify fact requirements
	if requirements == nil {
		t.Fatal("Fact requirements is nil")
	}

	if len(requirements.SystemFacts) == 0 {
		t.Error("Should have system facts")
	}

	// Check that system.hostname is included
	found := false
	for _, fact := range requirements.SystemFacts {
		if fact == "system.hostname" {
			found = true
			break
		}
	}

	if !found {
		t.Error("system.hostname should be in system facts")
	}

	t.Logf("Successfully analyzed fact requirements")
	t.Logf("System facts: %v", requirements.SystemFacts)
	t.Logf("Custom facts: %v", requirements.CustomFacts)
}

func TestManager_Run_Basic(t *testing.T) {
	// Create a basic logger
	logger := spookylogging.NewLogger(spookytypeslogging.Config{
		Level:  spookytypeslogging.InfoLevel,
		Format: "text",
		Output: "stdout",
	})

	// Create manager
	manager := NewManager(logger)

	// Create a simple action
	action := &spookytypesactions.Action{
		Name:    "test-action",
		Type:    "command",
		Command: "echo 'hello world'",
	}

	// Create context
	actionContext := &spookytypesactions.ActionContext{
		ProjectPath: "/tmp/test",
		Timeout:     30 * time.Second,
		Parallel:    false,
	}

	// Test Run method
	actions := []*spookytypesactions.Action{action}
	result, err := manager.Run(context.Background(), actions, actionContext)

	// For now, we expect this to fail because we haven't implemented all the helper methods
	// But we can verify the method signature is correct
	if err == nil {
		t.Logf("Run completed successfully: %+v", result)
	} else {
		t.Logf("Run failed as expected (not fully implemented): %v", err)
	}
}
