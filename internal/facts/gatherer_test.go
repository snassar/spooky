package facts

import (
	"context"
	"testing"
	"time"

	"spooky/internal/schemas"
	"spooky/internal/ssh"
)

func TestNewGatherer(t *testing.T) {
	// Test with nil SSH manager
	gatherer := NewGatherer(nil, nil)
	if gatherer == nil {
		t.Fatal("NewGatherer returned nil")
	}
	if gatherer.sshManager != nil {
		t.Error("expected sshManager to be nil when passed nil")
	}
	if gatherer.config != nil {
		t.Error("expected config to be nil when passed nil")
	}

	// Test with actual instances
	sshManager := ssh.NewSimpleSSHManager(nil, nil)
	config := &schemas.ProjectV1{
		Name:                    "test-project",
		FactsTimeout:            30,
		FactsParallelCollection: 10,
	}

	gatherer = NewGatherer(sshManager, config)
	if gatherer == nil {
		t.Fatal("NewGatherer returned nil")
	}
	if gatherer.sshManager == nil {
		t.Error("expected sshManager to be set")
	}
	if gatherer.config == nil {
		t.Error("expected config to be set")
	}
}

func TestGatherer_GatherFactsFromMachine_InputValidation(t *testing.T) {
	gatherer := NewGatherer(nil, nil)
	ctx := context.Background()

	tests := []struct {
		name        string
		machine     *schemas.MachinesMachineV1
		expectError bool
		errorMsg    string
	}{
		{
			name:        "nil machine",
			machine:     nil,
			expectError: true,
			errorMsg:    "machine cannot be nil",
		},
		{
			name: "empty hostname",
			machine: &schemas.MachinesMachineV1{
				Hostname: "",
			},
			expectError: true,
			errorMsg:    "machine hostname cannot be empty",
		},
		{
			name: "valid machine (will fail due to nil SSH manager)",
			machine: &schemas.MachinesMachineV1{
				Hostname: "test-machine",
				User:     "test-user",
				Port:     22,
			},
			expectError: true, // Will fail due to nil SSH manager
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := gatherer.GatherFactsFromMachine(ctx, tt.machine)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				} else if err.Error() != tt.errorMsg {
					t.Errorf("expected error message %q, got %q", tt.errorMsg, err.Error())
				}
				if result != nil {
					t.Errorf("expected nil result on error, got %v", result)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if result == nil {
					t.Errorf("expected result but got nil")
				}
			}
		})
	}
}

func TestGatherer_GatherFactsFromMachines_InputValidation(t *testing.T) {
	gatherer := NewGatherer(nil, nil)
	ctx := context.Background()

	// Test with empty machines list
	machines := []*schemas.MachinesMachineV1{}
	results, err := gatherer.GatherFactsFromMachines(ctx, machines)
	if err != nil {
		t.Errorf("unexpected error with empty machines list: %v", err)
	}
	if results == nil {
		t.Errorf("expected results but got nil")
	}
	if len(results) != 0 {
		t.Errorf("expected empty results, got %d results", len(results))
	}
}

func TestGatherer_GatherFactsFromMachines_WithConfig(t *testing.T) {
	config := &schemas.ProjectV1{
		Name:                    "test-project",
		FactsTimeout:            5, // Short timeout for testing
		FactsParallelCollection: 2,
	}
	gatherer := NewGatherer(nil, config)
	ctx := context.Background()

	// Test with machines that will fail (nil SSH manager)
	machines := []*schemas.MachinesMachineV1{
		{
			Hostname: "test-machine-1",
			User:     "test-user",
			Port:     22,
		},
		{
			Hostname: "test-machine-2",
			User:     "test-user",
			Port:     22,
		},
	}

	results, err := gatherer.GatherFactsFromMachines(ctx, machines)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if results == nil {
		t.Errorf("expected results but got nil")
	}
	if len(results) != len(machines) {
		t.Errorf("expected %d results, got %d", len(machines), len(results))
	}
}

func TestCreateEmptyBasicFacts(t *testing.T) {
	basicFacts := createEmptyBasicFacts()
	if basicFacts == nil {
		t.Fatal("createEmptyBasicFacts returned nil")
	}

	// Check that all fact categories are initialized
	if basicFacts.SystemFacts == nil {
		t.Error("SystemFacts should be initialized")
	}
	if basicFacts.HardwareFacts == nil {
		t.Error("HardwareFacts should be initialized")
	}
	if basicFacts.NetworkFacts == nil {
		t.Error("NetworkFacts should be initialized")
	}
	if basicFacts.OSFacts == nil {
		t.Error("OSFacts should be initialized")
	}
	if basicFacts.UserFacts == nil {
		t.Error("UserFacts should be initialized")
	}
	if basicFacts.RuntimeFacts == nil {
		t.Error("RuntimeFacts should be initialized")
	}

	// Check that fact maps are initialized
	if basicFacts.SystemFacts.Facts == nil {
		t.Error("SystemFacts.Facts should be initialized")
	}
	if basicFacts.HardwareFacts.Facts == nil {
		t.Error("HardwareFacts.Facts should be initialized")
	}
	if basicFacts.NetworkFacts.Facts == nil {
		t.Error("NetworkFacts.Facts should be initialized")
	}
	if basicFacts.OSFacts.Facts == nil {
		t.Error("OSFacts.Facts should be initialized")
	}
	if basicFacts.UserFacts.Facts == nil {
		t.Error("UserFacts.Facts should be initialized")
	}
	if basicFacts.RuntimeFacts.Facts == nil {
		t.Error("RuntimeFacts.Facts should be initialized")
	}
}

func TestMachineFacts_Structure(t *testing.T) {
	machine := &schemas.MachinesMachineV1{
		Hostname: "test-machine",
		User:     "test-user",
		Port:     22,
	}

	machineFacts := &MachineFacts{
		Machine:     machine,
		BasicFacts:  createEmptyBasicFacts(),
		CollectedAt: time.Now(),
	}

	if machineFacts.Machine == nil {
		t.Error("Machine should be set")
	}
	if machineFacts.BasicFacts == nil {
		t.Error("BasicFacts should be set")
	}
	if machineFacts.CollectedAt.IsZero() {
		t.Error("CollectedAt should be set")
	}
}
