package cmd

import (
	"strings"
	"testing"
)

func TestFactsCommandCreation(t *testing.T) {
	// Test that facts command is created correctly
	if factsCmd == nil {
		t.Fatal("factsCmd is nil")
	}

	if factsCmd.Use != "facts" {
		t.Errorf("expected facts command use 'facts', got %s", factsCmd.Use)
	}

	if factsCmd.Short == "" {
		t.Error("facts command should have a short description")
	}

	if factsCmd.Long == "" {
		t.Error("facts command should have a long description")
	}
}

func TestFactsSubcommands(t *testing.T) {
	// Test that all subcommands are created
	subcommands := factsCmd.Commands()

	expectedCommands := []string{"gather", "list", "validate", "export"}
	foundCommands := make(map[string]bool)

	for _, cmd := range subcommands {
		foundCommands[cmd.Name()] = true
	}

	for _, expected := range expectedCommands {
		if !foundCommands[expected] {
			t.Errorf("expected subcommand '%s' not found", expected)
		}
	}
}

func TestFactsGatherCommand(t *testing.T) {
	if factsGatherCmd == nil {
		t.Fatal("factsGatherCmd is nil")
	}

	if !strings.Contains(factsGatherCmd.Use, "gather") {
		t.Errorf("expected gather command use to contain 'gather', got %s", factsGatherCmd.Use)
	}

	if factsGatherCmd.Short == "" {
		t.Error("gather command should have a short description")
	}
}

func TestFactsListCommand(t *testing.T) {
	if factsListCmd == nil {
		t.Fatal("factsListCmd is nil")
	}

	if !strings.Contains(factsListCmd.Use, "list") {
		t.Errorf("expected list command use to contain 'list', got %s", factsListCmd.Use)
	}

	if factsListCmd.Short == "" {
		t.Error("list command should have a short description")
	}
}

func TestFactsValidateCommand(t *testing.T) {
	if factsValidateCmd == nil {
		t.Fatal("factsValidateCmd is nil")
	}

	if !strings.Contains(factsValidateCmd.Use, "validate") {
		t.Errorf("expected validate command use to contain 'validate', got %s", factsValidateCmd.Use)
	}

	if factsValidateCmd.Short == "" {
		t.Error("validate command should have a short description")
	}
}

func TestFactsExportCommand(t *testing.T) {
	if factsExportCmd == nil {
		t.Fatal("factsExportCmd is nil")
	}

	if !strings.Contains(factsExportCmd.Use, "export") {
		t.Errorf("expected export command use to contain 'export', got %s", factsExportCmd.Use)
	}

	if factsExportCmd.Short == "" {
		t.Error("export command should have a short description")
	}
}

func TestFactsCommandRegistration(t *testing.T) {
	// Test that facts command is registered with root command
	rootCommands := RootCmd.Commands()
	factsFound := false

	for _, cmd := range rootCommands {
		if cmd.Name() == "facts" {
			factsFound = true
			break
		}
	}

	if !factsFound {
		t.Error("facts command not found in root command")
	}
}
