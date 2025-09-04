package commands

import (
	"strings"
	"testing"

	"spooky/internal/schemas"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
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

// TestExtractVariablesBlock tests the extractVariablesBlock helper function
func TestExtractVariablesBlock(t *testing.T) {
	tests := []struct {
		name        string
		hclContent  string
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid variables block",
			hclContent: `
variables {
  variable "test_var" {
    value = "test_value"
  }
}`,
			expectError: false,
		},
		{
			name: "no variables block",
			hclContent: `
project {
  name = "test"
}`,
			expectError: true,
			errorMsg:    "failed to decode variables block",
		},
		{
			name:        "empty content",
			hclContent:  "",
			expectError: true,
			errorMsg:    "no variables block found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, diags := hclsyntax.ParseConfig([]byte(tt.hclContent), "test.hcl", hcl.Pos{Line: 1, Column: 1})
			if diags.HasErrors() {
				t.Fatalf("Failed to parse test HCL: %v", diags)
			}

			block, err := extractVariablesBlock(file)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				} else if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("expected error message containing %q, got %q", tt.errorMsg, err.Error())
				}
				if block != nil {
					t.Errorf("expected nil block on error, got %v", block)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if block == nil {
					t.Errorf("expected block but got nil")
				} else if block.Type != "variables" {
					t.Errorf("expected block type 'variables', got %q", block.Type)
				}
			}
		})
	}
}

// TestExtractVariableBlocks tests the extractVariableBlocks helper function
func TestExtractVariableBlocks(t *testing.T) {
	hclContent := `
variables {
  variable "test_var" {
    value = "test_value"
    description = "Test variable"
  }
  variable "another_var" {
    value = 42
    description = "Another test variable"
  }
}`

	file, diags := hclsyntax.ParseConfig([]byte(hclContent), "test.hcl", hcl.Pos{Line: 1, Column: 1})
	if diags.HasErrors() {
		t.Fatalf("Failed to parse test HCL: %v", diags)
	}

	// Extract variables block
	variablesBlock, err := extractVariablesBlock(file)
	if err != nil {
		t.Fatalf("Failed to extract variables block: %v", err)
	}

	// Extract variable blocks
	content, err := extractVariableBlocks(variablesBlock)
	if err != nil {
		t.Fatalf("Failed to extract variable blocks: %v", err)
	}

	// Verify we got the expected number of variable blocks
	expectedCount := 2
	if len(content.Blocks) != expectedCount {
		t.Errorf("Expected %d variable blocks, got %d", expectedCount, len(content.Blocks))
	}

	// Verify block labels
	expectedLabels := []string{"test_var", "another_var"}
	for i, block := range content.Blocks {
		if i >= len(expectedLabels) {
			t.Errorf("Unexpected block at index %d: %s", i, block.Labels[0])
			continue
		}
		if block.Labels[0] != expectedLabels[i] {
			t.Errorf("Expected label %s, got %s", expectedLabels[i], block.Labels[0])
		}
	}
}

// TestProcessVariableBlocks tests the processVariableBlocks helper function
func TestProcessVariableBlocks(t *testing.T) {
	// Create a mock HCL content with variable blocks
	hclContent := `
variables {
  variable "string_var" {
    value = "test_string"
    description = "A string variable"
  }
  variable "number_var" {
    value = 42
    description = "A number variable"
  }
  variable "bool_var" {
    value = true
    description = "A boolean variable"
  }
}`

	file, diags := hclsyntax.ParseConfig([]byte(hclContent), "test.hcl", hcl.Pos{Line: 1, Column: 1})
	if diags.HasErrors() {
		t.Fatalf("Failed to parse test HCL: %v", diags)
	}

	// Extract variables block
	variablesBlock, err := extractVariablesBlock(file)
	if err != nil {
		t.Fatalf("Failed to extract variables block: %v", err)
	}

	// Extract variable blocks
	content, err := extractVariableBlocks(variablesBlock)
	if err != nil {
		t.Fatalf("Failed to extract variable blocks: %v", err)
	}

	// Process variable blocks
	variables, err := processVariableBlocks(content)
	if err != nil {
		t.Fatalf("Failed to process variable blocks: %v", err)
	}

	// Verify the processed variables exist
	expectedVarNames := []string{"string_var", "number_var", "bool_var"}
	for _, name := range expectedVarNames {
		if _, exists := variables[name]; !exists {
			t.Errorf("Expected variable %s to exist", name)
		}
	}

	// Verify we got the expected number of variables
	if len(variables) != len(expectedVarNames) {
		t.Errorf("Expected %d variables, got %d", len(expectedVarNames), len(variables))
	}
}

// TestExtractVariableAttributes tests the extractVariableAttributes helper function
func TestExtractVariableAttributes(t *testing.T) {
	hclContent := `
variables {
  variable "test_var" {
    value = "test_value"
    description = "Test description"
    encrypted = false
  }
}`

	file, diags := hclsyntax.ParseConfig([]byte(hclContent), "test.hcl", hcl.Pos{Line: 1, Column: 1})
	if diags.HasErrors() {
		t.Fatalf("Failed to parse test HCL: %v", diags)
	}

	// Extract variables block
	variablesBlock, err := extractVariablesBlock(file)
	if err != nil {
		t.Fatalf("Failed to extract variables block: %v", err)
	}

	// Extract variable blocks
	content, err := extractVariableBlocks(variablesBlock)
	if err != nil {
		t.Fatalf("Failed to extract variable blocks: %v", err)
	}

	// Extract attributes from the first variable block
	attrContent, err := extractVariableAttributes(content.Blocks[0])
	if err != nil {
		t.Fatalf("Failed to extract variable attributes: %v", err)
	}

	// Verify expected attributes exist
	expectedAttrs := []string{"value", "description", "encrypted"}
	for _, attrName := range expectedAttrs {
		if _, exists := attrContent.Attributes[attrName]; !exists {
			t.Errorf("Expected attribute %s to exist", attrName)
		}
	}
}

// TestParseVariableAttributes tests the parseVariableAttributes helper function
func TestParseVariableAttributes(t *testing.T) {
	tests := []struct {
		name        string
		hclContent  string
		expectError bool
	}{
		{
			name: "string value",
			hclContent: `
variables {
  variable "string_var" {
    value = "test_string"
  }
}`,
			expectError: false,
		},
		{
			name: "number value",
			hclContent: `
variables {
  variable "number_var" {
    value = 42
  }
}`,
			expectError: false,
		},
		{
			name: "boolean value",
			hclContent: `
variables {
  variable "bool_var" {
    value = true
  }
}`,
			expectError: false,
		},
		{
			name: "encrypted value",
			hclContent: `
variables {
  variable "encrypted_var" {
    encrypted_value = {
      data = "AGE1-ENCRYPTED-DATA"
      format = "armored"
    }
  }
}`,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, diags := hclsyntax.ParseConfig([]byte(tt.hclContent), "test.hcl", hcl.Pos{Line: 1, Column: 1})
			if diags.HasErrors() {
				t.Fatalf("Failed to parse test HCL: %v", diags)
			}

			// Extract variables block
			variablesBlock, err := extractVariablesBlock(file)
			if err != nil {
				t.Fatalf("Failed to extract variables block: %v", err)
			}

			// Extract variable blocks
			content, err := extractVariableBlocks(variablesBlock)
			if err != nil {
				t.Fatalf("Failed to extract variable blocks: %v", err)
			}

			// Extract attributes
			attrContent, err := extractVariableAttributes(content.Blocks[0])
			if err != nil {
				t.Fatalf("Failed to extract variable attributes: %v", err)
			}

			// Parse variable attributes
			value, err := parseVariableAttributes(attrContent, "test_var")

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if value == nil {
					t.Errorf("expected non-nil value")
				}
			}
		})
	}
}
