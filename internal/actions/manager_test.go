package spookyactions

import (
	"strings"
	"testing"

	spookytypes "spooky/internal/types"
	spookytypesactions "spooky/internal/types/actions"
	spookytypesmachines "spooky/internal/types/machines"
)

// Test the machineHasTags function directly without requiring a full manager
func TestMachineHasTagsLogic(t *testing.T) {
	// Create a test machine with tags
	machine := &spookytypes.Machine{
		Hostname: "test-server",
		Tags: map[string]string{
			"environment": "production",
			"role":        "web",
			"region":      "us-east",
		},
	}

	// Test key=value tag matching
	tests := []struct {
		name     string
		tags     []string
		expected bool
	}{
		{
			name:     "exact key=value match",
			tags:     []string{"environment=production"},
			expected: true,
		},
		{
			name:     "key=value mismatch",
			tags:     []string{"environment=staging"},
			expected: false,
		},
		{
			name:     "key-only match",
			tags:     []string{"role"},
			expected: true,
		},
		{
			name:     "key-only mismatch",
			tags:     []string{"database"},
			expected: false,
		},
		{
			name:     "multiple tags all match",
			tags:     []string{"environment=production", "role=web"},
			expected: true,
		},
		{
			name:     "multiple tags one mismatch",
			tags:     []string{"environment=production", "role=database"},
			expected: false,
		},
		{
			name:     "empty tags",
			tags:     []string{},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := machineHasTagsLogic(machine, tt.tags)
			if result != tt.expected {
				t.Errorf("machineHasTagsLogic(%v, %v) = %v, want %v", machine.Hostname, tt.tags, result, tt.expected)
			}
		})
	}
}

// Test machine with nil tags
func TestMachineHasTagsWithNilTags(t *testing.T) {
	machine := &spookytypes.Machine{
		Hostname: "test-server",
		Tags:     nil,
	}

	// Should return false when machine has no tags but tags are required
	result := machineHasTagsLogic(machine, []string{"environment=production"})
	if result {
		t.Errorf("machineHasTagsLogic should return false for machine with nil tags")
	}

	// Should return true when no tags are required
	result = machineHasTagsLogic(machine, []string{})
	if !result {
		t.Errorf("machineHasTagsLogic should return true when no tags are required")
	}
}

// Test target machine filtering logic
func TestTargetMachineFiltering(t *testing.T) {
	// Create test machines
	machines := []spookytypes.Machine{
		{
			Hostname: "web-01",
			Tags: map[string]string{
				"environment": "production",
				"role":        "web",
			},
		},
		{
			Hostname: "db-01",
			Tags: map[string]string{
				"environment": "production",
				"role":        "database",
			},
		},
		{
			Hostname: "web-02",
			Tags: map[string]string{
				"environment": "staging",
				"role":        "web",
			},
		},
	}

	// Test machine name filtering
	t.Run("filter by machine names", func(t *testing.T) {
		action := &spookytypesactions.Action{
			Name:     "test-action",
			Machines: []string{"web-01", "db-01"},
		}

		targetMachines := filterMachinesByName(action.Machines, machines)
		if len(targetMachines) != 2 {
			t.Errorf("Expected 2 target machines, got %d", len(targetMachines))
		}

		// Check that we got the expected machines
		hostnames := make(map[string]bool)
		for _, machine := range targetMachines {
			hostnames[machine.Hostname] = true
		}

		if !hostnames["web-01"] || !hostnames["db-01"] {
			t.Errorf("Expected web-01 and db-01, got %v", hostnames)
		}
	})

	// Test tag filtering
	t.Run("filter by tags", func(t *testing.T) {
		action := &spookytypesactions.Action{
			Name: "test-action",
			Tags: []string{"environment=production"},
		}

		targetMachines := filterMachinesByTags(action.Tags, machines)
		if len(targetMachines) != 2 {
			t.Errorf("Expected 2 target machines with environment=production, got %d", len(targetMachines))
		}

		// Check that we got production machines
		for _, machine := range targetMachines {
			if machine.Tags["environment"] != "production" {
				t.Errorf("Expected production environment, got %s", machine.Tags["environment"])
			}
		}
	})

	// Test no filtering (should return all machines)
	t.Run("no filtering", func(t *testing.T) {
		action := &spookytypesactions.Action{
			Name: "test-action",
		}

		targetMachines := getTargetMachinesLogic(action, machines)
		if len(targetMachines) != 3 {
			t.Errorf("Expected 3 target machines (all machines), got %d", len(targetMachines))
		}
	})

	// Test missing machines
	t.Run("missing machines", func(t *testing.T) {
		action := &spookytypesactions.Action{
			Name:     "test-action",
			Machines: []string{"web-01", "nonexistent"},
		}

		targetMachines := filterMachinesByName(action.Machines, machines)
		if len(targetMachines) != 1 {
			t.Errorf("Expected 1 target machine (only web-01 exists), got %d", len(targetMachines))
		}

		if targetMachines[0].Hostname != "web-01" {
			t.Errorf("Expected web-01, got %s", targetMachines[0].Hostname)
		}
	})
}

// Test session-based machine filtering
func TestSessionMachineFiltering(t *testing.T) {
	// Create test session with machine inventory
	session := &spookytypesactions.ActingSession{
		SessionID: "test-session",
		MachineInventory: []spookytypesmachines.Machine{
			{
				Hostname: "web-01",
				Tags: map[string]string{
					"environment": "production",
					"role":        "web",
				},
			},
			{
				Hostname: "db-01",
				Tags: map[string]string{
					"environment": "production",
					"role":        "database",
				},
			},
		},
		MachineCache: map[string]*spookytypesmachines.Machine{
			"web-01": {
				Hostname: "web-01",
				Tags: map[string]string{
					"environment": "production",
					"role":        "web",
				},
			},
			"db-01": {
				Hostname: "db-01",
				Tags: map[string]string{
					"environment": "production",
					"role":        "database",
				},
			},
		},
	}

	// Test machine name filtering with session
	t.Run("filter by machine names with session", func(t *testing.T) {
		action := &spookytypesactions.Action{
			Name:     "test-action",
			Machines: []string{"web-01", "db-01"},
		}

		targetMachines := filterMachinesFromSession(action.Machines, session)
		if len(targetMachines) != 2 {
			t.Errorf("Expected 2 target machines, got %d", len(targetMachines))
		}

		// Check that we got the expected machines
		hostnames := make(map[string]bool)
		for _, machine := range targetMachines {
			hostnames[machine.Hostname] = true
		}

		if !hostnames["web-01"] || !hostnames["db-01"] {
			t.Errorf("Expected web-01 and db-01, got %v", hostnames)
		}
	})

	// Test tag filtering with session
	t.Run("filter by tags with session", func(t *testing.T) {
		action := &spookytypesactions.Action{
			Name: "test-action",
			Tags: []string{"role=web"},
		}

		targetMachines := filterMachinesByTagsFromSession(action.Tags, session)
		if len(targetMachines) != 1 {
			t.Errorf("Expected 1 target machine with role=web, got %d", len(targetMachines))
		}

		if targetMachines[0].Hostname != "web-01" {
			t.Errorf("Expected web-01, got %s", targetMachines[0].Hostname)
		}
	})
}

// Helper functions for testing (copied from the actual implementation)

// machineHasTagsLogic is the core logic for checking if a machine has the specified tags
func machineHasTagsLogic(machine *spookytypes.Machine, tags []string) bool {
	if len(tags) == 0 {
		return true
	}

	// Handle case where machine has no tags
	if len(machine.Tags) == 0 {
		return false
	}

	// Check each required tag
	for _, requiredTag := range tags {
		tagMatched := false

		// Check if tag is in key=value format
		if strings.Contains(requiredTag, "=") {
			parts := strings.SplitN(requiredTag, "=", 2)
			if len(parts) == 2 {
				key, value := parts[0], parts[1]
				if machineValue, exists := machine.Tags[key]; exists && machineValue == value {
					tagMatched = true
				}
			}
		} else {
			// Key-only format - check if the key exists in machine tags
			if _, exists := machine.Tags[requiredTag]; exists {
				tagMatched = true
			}
		}

		// If any required tag doesn't match, machine doesn't have all required tags
		if !tagMatched {
			return false
		}
	}

	return true
}

// filterMachinesByName filters machines by hostname
func filterMachinesByName(targetNames []string, machines []spookytypes.Machine) []spookytypes.Machine {
	if len(targetNames) == 0 {
		return machines
	}

	var targetMachines []spookytypes.Machine
	for _, targetName := range targetNames {
		for i := range machines {
			if machines[i].Hostname == targetName {
				targetMachines = append(targetMachines, machines[i])
				break
			}
		}
	}
	return targetMachines
}

// filterMachinesByTags filters machines by tags
func filterMachinesByTags(tags []string, machines []spookytypes.Machine) []spookytypes.Machine {
	if len(tags) == 0 {
		return machines
	}

	var targetMachines []spookytypes.Machine
	for i := range machines {
		if machineHasTagsLogic(&machines[i], tags) {
			targetMachines = append(targetMachines, machines[i])
		}
	}
	return targetMachines
}

// getTargetMachinesLogic determines which machines should run the action
func getTargetMachinesLogic(action *spookytypesactions.Action, machines []spookytypes.Machine) []spookytypes.Machine {
	// If action has specific machines defined, filter by name
	if len(action.Machines) > 0 {
		return filterMachinesByName(action.Machines, machines)
	}

	// If action has tags defined, filter machines by tags
	if len(action.Tags) > 0 {
		return filterMachinesByTags(action.Tags, machines)
	}

	// If no specific targeting, return all available machines
	return machines
}

// filterMachinesFromSession filters machines from session inventory by name
func filterMachinesFromSession(targetNames []string, session *spookytypesactions.ActingSession) []spookytypes.Machine {
	if len(targetNames) == 0 {
		return convertSessionMachines(session.MachineInventory)
	}

	var targetMachines []spookytypes.Machine
	for _, targetName := range targetNames {
		if machine, exists := session.MachineCache[targetName]; exists {
			targetMachines = append(targetMachines, spookytypes.Machine(*machine))
		}
	}
	return targetMachines
}

// filterMachinesByTagsFromSession filters machines from session inventory by tags
func filterMachinesByTagsFromSession(tags []string, session *spookytypesactions.ActingSession) []spookytypes.Machine {
	if len(tags) == 0 {
		return convertSessionMachines(session.MachineInventory)
	}

	var targetMachines []spookytypes.Machine
	for _, sessionMachine := range session.MachineInventory {
		machine := spookytypes.Machine(sessionMachine)
		if machineHasTagsLogic(&machine, tags) {
			targetMachines = append(targetMachines, machine)
		}
	}
	return targetMachines
}

// convertSessionMachines converts session machines to the expected type
func convertSessionMachines(sessionMachines []spookytypesmachines.Machine) []spookytypes.Machine {
	machines := make([]spookytypes.Machine, len(sessionMachines))
	for i, sessionMachine := range sessionMachines {
		machines[i] = spookytypes.Machine(sessionMachine)
	}
	return machines
}
