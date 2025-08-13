package cmd

import (
	"strings"
	"testing"

	spookytypes "spooky/internal/types"
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

	expectedCommands := []string{"export"}
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

func TestFilterMachinesAdvanced(t *testing.T) {
	// Create test machines
	machines := []spookytypes.Machine{
		{
			Hostname: "web-001",
			Tags: map[string]string{
				"environment": "production",
				"role":        "web",
			},
			Groups: []string{"webservers", "production"},
		},
		{
			Hostname: "db-001",
			Tags: map[string]string{
				"environment": "production",
				"role":        "database",
			},
			Groups: []string{"database", "production"},
		},
		{
			Hostname: "web-002",
			Tags: map[string]string{
				"environment": "staging",
				"role":        "web",
			},
			Groups: []string{"webservers", "staging"},
		},
	}

	// Test machine filter
	filtered := filterMachinesAdvanced(machines, "web", nil, nil)
	if len(filtered) != 2 {
		t.Errorf("expected 2 machines with 'web' in hostname, got %d", len(filtered))
	}

	// Test tag filter (key=value)
	filtered = filterMachinesAdvanced(machines, "", []string{"environment=production"}, nil)
	if len(filtered) != 2 {
		t.Errorf("expected 2 machines with environment=production, got %d", len(filtered))
	}

	// Test tag filter (key only)
	filtered = filterMachinesAdvanced(machines, "", []string{"role"}, nil)
	if len(filtered) != 3 {
		t.Errorf("expected 3 machines with role tag, got %d", len(filtered))
	}

	// Test group filter
	filtered = filterMachinesAdvanced(machines, "", nil, []string{"webservers"})
	if len(filtered) != 2 {
		t.Errorf("expected 2 machines in webservers group, got %d", len(filtered))
	}

	// Test combined filters
	filtered = filterMachinesAdvanced(machines, "web", []string{"environment=production"}, []string{"webservers"})
	if len(filtered) != 1 {
		t.Errorf("expected 1 machine matching all filters, got %d", len(filtered))
	}
	if filtered[0].Hostname != "web-001" {
		t.Errorf("expected web-001, got %s", filtered[0].Hostname)
	}
}

func TestMatchesMachineFilter(t *testing.T) {
	machine := spookytypes.Machine{Hostname: "web-server-001"}

	// Test empty filter
	if !matchesMachineFilter(machine, "") {
		t.Error("empty filter should match all machines")
	}

	// Test matching filter
	if !matchesMachineFilter(machine, "web") {
		t.Error("machine should match 'web' filter")
	}

	// Test non-matching filter
	if matchesMachineFilter(machine, "db") {
		t.Error("machine should not match 'db' filter")
	}
}

func TestMatchesTagsFilter(t *testing.T) {
	machine := spookytypes.Machine{
		Tags: map[string]string{
			"environment": "production",
			"role":        "web",
		},
	}

	// Test empty filter
	if !matchesTagsFilter(machine, nil) {
		t.Error("empty filter should match all machines")
	}

	// Test key=value filter
	if !matchesTagsFilter(machine, []string{"environment=production"}) {
		t.Error("machine should match environment=production")
	}

	// Test key-only filter
	if !matchesTagsFilter(machine, []string{"role"}) {
		t.Error("machine should match role tag")
	}

	// Test non-matching filter
	if matchesTagsFilter(machine, []string{"environment=staging"}) {
		t.Error("machine should not match environment=staging")
	}
}

func TestMatchesGroupsFilter(t *testing.T) {
	machine := spookytypes.Machine{
		Groups: []string{"webservers", "production"},
	}

	// Test empty filter
	if !matchesGroupsFilter(machine, nil) {
		t.Error("empty filter should match all machines")
	}

	// Test matching group
	if !matchesGroupsFilter(machine, []string{"webservers"}) {
		t.Error("machine should match webservers group")
	}

	// Test non-matching group
	if matchesGroupsFilter(machine, []string{"database"}) {
		t.Error("machine should not match database group")
	}
}
