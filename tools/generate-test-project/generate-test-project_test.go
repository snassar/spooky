package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGenerateID(t *testing.T) {
	id1 := generateID()
	id2 := generateID()

	if id1 == "" {
		t.Error("generateID returned empty string")
	}

	if id2 == "" {
		t.Error("generateID returned empty string")
	}

	// IDs should be different (due to random component)
	if id1 == id2 {
		t.Error("generateID returned same ID twice")
	}
}

func TestGenerateGitStyleID(t *testing.T) {
	metadata1 := "test-metadata-1"
	metadata2 := "test-metadata-2"

	id1 := generateGitStyleID(metadata1)
	id2 := generateGitStyleID(metadata2)
	id3 := generateGitStyleID(metadata1) // Same metadata

	if id1 == "" {
		t.Error("generateGitStyleID returned empty string")
	}

	if len(id1) != 16 {
		t.Errorf("generateGitStyleID returned ID of length %d, expected 16", len(id1))
	}

	// Same metadata should produce same ID
	if id1 != id3 {
		t.Error("generateGitStyleID returned different IDs for same metadata")
	}

	// Different metadata should produce different IDs
	if id1 == id2 {
		t.Error("generateGitStyleID returned same ID for different metadata")
	}
}

func TestGenerateMachines(t *testing.T) {
	scale := ScaleConfig{
		Name:     "test",
		Hardware: 4,
		VMs:      8,
	}

	machines := generateMachines(scale)

	expectedTotal := scale.Hardware + scale.VMs
	if len(machines) != expectedTotal {
		t.Errorf("Expected %d machines, got %d", expectedTotal, len(machines))
	}

	// Check that machines have required fields
	for i, machine := range machines {
		if machine.ID == "" {
			t.Errorf("Machine %d has empty ID", i)
		}
		if machine.Host == "" {
			t.Errorf("Machine %d has empty Host", i)
		}
		if machine.User == "" {
			t.Errorf("Machine %d has empty User", i)
		}
		if machine.Port == 0 {
			t.Errorf("Machine %d has invalid Port", i)
		}
	}
}

func TestGenerateActions(t *testing.T) {
	actions := generateActions()

	if len(actions) == 0 {
		t.Error("generateActions returned no actions")
	}

	// Check that actions have required fields
	for i, action := range actions {
		if action.Name == "" {
			t.Errorf("Action %d has empty Name", i)
		}
		if action.Description == "" {
			t.Errorf("Action %d has empty Description", i)
		}
		if action.Command == "" && action.Script == "" {
			t.Errorf("Action %d has neither Command nor Script", i)
		}
		if action.Command != "" && action.Script != "" {
			t.Errorf("Action %d has both Command and Script", i)
		}
	}
}

func TestWriteMachine(t *testing.T) {
	// Create a temporary file
	tmpfile, err := os.CreateTemp("", "test-machine")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())
	defer tmpfile.Close()

	machine := &Machine{
		ID:       "test-machine",
		Host:     "192.168.1.1",
		Port:     22,
		User:     "admin",
		Password: "password",
		Tags: map[string]string{
			"datacenter": "FRA00",
			"type":       "hardware",
		},
	}

	writeMachine(tmpfile, machine)

	// Read the file content
	content, err := os.ReadFile(tmpfile.Name())
	if err != nil {
		t.Fatal(err)
	}

	contentStr := string(content)

	// Check that the content contains expected fields
	expectedFields := []string{
		"machine \"test-machine\"",
		"host     = \"192.168.1.1\"",
		"port     = 22",
		"user     = \"admin\"",
		"password = \"password\"",
		"datacenter = \"FRA00\"",
		"type = \"hardware\"",
	}

	for _, field := range expectedFields {
		if !contains(contentStr, field) {
			t.Errorf("Generated content missing field: %s", field)
		}
	}
}

func TestWriteAction(t *testing.T) {
	// Create a temporary file
	tmpfile, err := os.CreateTemp("", "test-action")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())
	defer tmpfile.Close()

	action := &Action{
		Name:        "test-action",
		Description: "Test action description",
		Command:     "echo 'test'",
		Tags:        []string{"test", "example"},
		Machines:    []string{"all"},
		Parallel:    true,
		Timeout:     300,
	}

	writeAction(tmpfile, action)

	// Read the file content
	content, err := os.ReadFile(tmpfile.Name())
	if err != nil {
		t.Fatal(err)
	}

	contentStr := string(content)

	// Check that the content contains expected fields
	expectedFields := []string{
		"action \"test-action\"",
		"description = \"Test action description\"",
		"command = \"echo 'test'\"",
		"\"test\"",
		"\"example\"",
		"\"all\"",
		"parallel    = true",
		"timeout     = 300",
	}

	for _, field := range expectedFields {
		if !contains(contentStr, field) {
			t.Errorf("Generated content missing field: %s", field)
		}
	}
}

func TestWriteInventoryFile(t *testing.T) {
	// Create a temporary directory
	tmpdir, err := os.MkdirTemp("", "test-inventory")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpdir)

	machines := []Machine{
		{
			ID:       "machine-1",
			Host:     "192.168.1.1",
			Port:     22,
			User:     "admin",
			Password: "password1",
			Tags: map[string]string{
				"datacenter": "FRA00",
			},
		},
		{
			ID:       "machine-2",
			Host:     "192.168.1.2",
			Port:     22,
			User:     "admin",
			Password: "password2",
			Tags: map[string]string{
				"datacenter": "BER0",
			},
		},
	}

	err = writeInventoryFile(tmpdir, machines)
	if err != nil {
		t.Fatalf("writeInventoryFile failed: %v", err)
	}

	// Check that the file was created
	inventoryFile := filepath.Join(tmpdir, "inventory.hcl")
	if _, err := os.Stat(inventoryFile); os.IsNotExist(err) {
		t.Error("inventory.hcl file was not created")
	}

	// Read the file content
	content, err := os.ReadFile(inventoryFile)
	if err != nil {
		t.Fatal(err)
	}

	contentStr := string(content)

	// Check that the content contains expected elements
	expectedElements := []string{
		"# Inventory file for test project",
		"# Total machines: 2",
		"machine \"machine-1\"",
		"machine \"machine-2\"",
		"192.168.1.1",
		"192.168.1.2",
	}

	for _, element := range expectedElements {
		if !contains(contentStr, element) {
			t.Errorf("Generated inventory missing element: %s", element)
		}
	}
}

func TestWriteActionsFile(t *testing.T) {
	// Create a temporary directory
	tmpdir, err := os.MkdirTemp("", "test-actions")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpdir)

	actions := []Action{
		{
			Name:        "action-1",
			Description: "Test action 1",
			Command:     "echo 'test1'",
			Tags:        []string{"test"},
			Machines:    []string{"all"},
			Parallel:    true,
			Timeout:     300,
		},
		{
			Name:        "action-2",
			Description: "Test action 2",
			Script:      "scripts/test.sh",
			Tags:        []string{"test", "script"},
			Machines:    []string{"tag:type=hardware"},
			Parallel:    false,
			Timeout:     600,
		},
	}

	err = writeActionsFile(tmpdir, actions)
	if err != nil {
		t.Fatalf("writeActionsFile failed: %v", err)
	}

	// Check that the file was created
	actionsFile := filepath.Join(tmpdir, "actions.hcl")
	if _, err := os.Stat(actionsFile); os.IsNotExist(err) {
		t.Error("actions.hcl file was not created")
	}

	// Read the file content
	content, err := os.ReadFile(actionsFile)
	if err != nil {
		t.Fatal(err)
	}

	contentStr := string(content)

	// Check that the content contains expected elements
	expectedElements := []string{
		"# Actions file for test project",
		"action \"action-1\"",
		"action \"action-2\"",
		"echo 'test1'",
		"scripts/test.sh",
		"\"test\"",
		"\"script\"",
		"\"all\"",
		"\"tag:type=hardware\"",
	}

	for _, element := range expectedElements {
		if !contains(contentStr, element) {
			t.Errorf("Generated actions missing element: %s", element)
		}
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || containsSubstring(s, substr)))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
