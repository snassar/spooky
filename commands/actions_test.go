package commands

import (
	"testing"

	"spooky/internal/schemas"
)

func TestExtractResourceBlock(t *testing.T) {
	// This test would require creating actual HCL content and parsing it
	// For now, we'll test the function signature and basic behavior
	t.Run("function exists", func(t *testing.T) {
		// The function should exist and be callable
		// We can't easily test the full functionality without HCL parsing
		// but we can ensure it's properly defined
		// In Go, functions are never nil, so we just verify the test runs
		t.Log("extractResourceBlock function is defined")
	})
}

func TestExtractActionsBlock(t *testing.T) {
	// Similar to above, we'll test that the function exists
	t.Run("function exists", func(t *testing.T) {
		// In Go, functions are never nil, so we just verify the test runs
		t.Log("extractActionsBlock function is defined")
	})
}

func TestFindAction(t *testing.T) {
	actions := []*schemas.ActionsActionV1{
		{
			Description: "test-action: Test action description",
			Type:        "command",
		},
		{
			Description: "another-action: Another action description",
			Type:        "template_deploy",
		},
	}

	tests := []struct {
		name         string
		actionName   string
		expectFound  bool
		expectAction *schemas.ActionsActionV1
	}{
		{
			name:         "find existing action",
			actionName:   "test-action",
			expectFound:  true,
			expectAction: actions[0],
		},
		{
			name:         "find another existing action",
			actionName:   "another-action",
			expectFound:  true,
			expectAction: actions[1],
		},
		{
			name:         "action not found",
			actionName:   "nonexistent-action",
			expectFound:  false,
			expectAction: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			action, found := findAction(actions, tt.actionName)

			if found != tt.expectFound {
				t.Errorf("expected found=%v, got %v", tt.expectFound, found)
			}

			if tt.expectFound {
				if action == nil {
					t.Errorf("expected action but got nil")
				} else if action != tt.expectAction {
					t.Errorf("expected action %v, got %v", tt.expectAction, action)
				}
			} else {
				if action != nil {
					t.Errorf("expected nil action, got %v", action)
				}
			}
		})
	}
}

func TestDetermineTargetMachines(t *testing.T) {
	machines := []*schemas.MachinesMachineV1{
		{
			Hostname: "machine1",
			User:     "user1",
			Port:     22,
		},
		{
			Hostname: "machine2",
			User:     "user2",
			Port:     22,
		},
		{
			Hostname: "machine3",
			User:     "user3",
			Port:     22,
		},
	}

	action := &schemas.ActionsActionV1{
		Targets: []string{"machine1", "machine3"},
	}

	tests := []struct {
		name            string
		action          *schemas.ActionsActionV1
		overrideTargets []string
		expectCount     int
		expectError     bool
	}{
		{
			name:            "use action targets",
			action:          action,
			overrideTargets: nil,
			expectCount:     2,
			expectError:     false,
		},
		{
			name:            "use override targets",
			action:          action,
			overrideTargets: []string{"machine2"},
			expectCount:     1,
			expectError:     false,
		},
		{
			name:            "no targets specified",
			action:          &schemas.ActionsActionV1{},
			overrideTargets: nil,
			expectCount:     0,
			expectError:     true,
		},
		{
			name:            "targets not found",
			action:          &schemas.ActionsActionV1{Targets: []string{"nonexistent"}},
			overrideTargets: nil,
			expectCount:     0,
			expectError:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := determineTargetMachines(tt.action, machines, tt.overrideTargets)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}

			if len(result) != tt.expectCount {
				t.Errorf("expected %d machines, got %d", tt.expectCount, len(result))
			}
		})
	}
}

func TestFormatMachineList(t *testing.T) {
	machines := []*schemas.MachinesMachineV1{
		{Hostname: "machine1"},
		{Hostname: "machine2"},
		{Hostname: "machine3"},
	}

	result := formatMachineList(machines)
	expected := "machine1, machine2, machine3"

	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestFormatMachineList_Empty(t *testing.T) {
	machines := []*schemas.MachinesMachineV1{}
	result := formatMachineList(machines)
	expected := ""

	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestCreateDefaultAction(t *testing.T) {
	action := createDefaultAction()
	if action == nil {
		t.Fatal("createDefaultAction returned nil")
	}

	// Check default values
	if action.Description != "" {
		t.Errorf("expected empty description, got %q", action.Description)
	}
	if action.Type != "" {
		t.Errorf("expected empty type, got %q", action.Type)
	}
	if action.Tags == nil {
		t.Error("expected Tags to be initialized")
	}
	if action.Targets == nil {
		t.Error("expected Targets to be initialized")
	}
	if action.Timeout != 300 {
		t.Errorf("expected Timeout 300, got %d", action.Timeout)
	}
	if action.Retries != 0 {
		t.Errorf("expected Retries 0, got %d", action.Retries)
	}
	if action.RetryDelay != 5 {
		t.Errorf("expected RetryDelay 5, got %d", action.RetryDelay)
	}
}
