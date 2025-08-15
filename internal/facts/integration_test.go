package facts

import (
	"context"
	"testing"
	"time"

	spookyinterfaces "spooky/internal/interfaces"
	spookytypes "spooky/internal/types"
	spookytypesfacts "spooky/internal/types/facts"
)

func TestNewIntegration(t *testing.T) {
	// Create a real manager with mock dependencies
	mockCollector := &MockFactCollector{}
	mockValidator := &MockSchemaValidator{}
	mockLogger := &MockLogger{}

	manager := NewManager(mockCollector, mockValidator, mockLogger)

	integration := NewIntegration(manager)

	if integration == nil {
		t.Fatal("NewIntegration returned nil")
	}

	// Test that the integration implements the interface
	var _ spookyinterfaces.FactsIntegration = integration
}

func TestIntegrationCollectFacts(t *testing.T) {
	// Create a real manager with mock dependencies
	mockCollector := &MockFactCollector{}
	mockValidator := &MockSchemaValidator{}
	mockLogger := &MockLogger{}

	manager := NewManager(mockCollector, mockValidator, mockLogger)
	integration := NewIntegration(manager)

	machine := &spookytypes.Machine{
		Hostname: "test-server",
		Port:     22,
		User:     "test-user",
	}

	facts, err := integration.CollectFacts(context.Background(), machine)
	if err != nil {
		t.Fatalf("CollectFacts failed: %v", err)
	}

	if facts == nil {
		t.Fatal("CollectFacts returned nil facts")
	}

	// Verify facts structure
	if factCollection, ok := facts.(*FactCollection); ok {
		if factCollection.MachineID != machine.Hostname {
			t.Errorf("Expected machine ID %s, got %s", machine.Hostname, factCollection.MachineID)
		}
	} else {
		t.Fatal("CollectFacts did not return *FactCollection")
	}
}

func TestIntegrationStoreFacts(t *testing.T) {
	// Create a real manager with mock dependencies
	mockCollector := &MockFactCollector{}
	mockValidator := &MockSchemaValidator{}
	mockLogger := &MockLogger{}

	manager := NewManager(mockCollector, mockValidator, mockLogger)
	integration := NewIntegration(manager)

	facts := &FactCollection{
		MachineID:   "test-server",
		CollectedAt: time.Now(),
		Facts: &spookytypesfacts.Facts{
			System: &spookytypesfacts.SystemFacts{
				OS: &spookytypesfacts.OSFacts{
					Name:    "Linux",
					Version: "Ubuntu 20.04",
				},
			},
		},
	}

	err := integration.StoreFacts(context.Background(), facts)
	if err != nil {
		t.Fatalf("StoreFacts failed: %v", err)
	}
}

func TestIntegrationLoadFacts(t *testing.T) {
	// Create a real manager with mock dependencies
	mockCollector := &MockFactCollector{}
	mockValidator := &MockSchemaValidator{}
	mockLogger := &MockLogger{}

	manager := NewManager(mockCollector, mockValidator, mockLogger)
	integration := NewIntegration(manager)

	facts, err := integration.LoadFacts(context.Background())
	if err != nil {
		t.Fatalf("LoadFacts failed: %v", err)
	}

	// LoadFacts returns nil for in-memory storage
	if facts != nil {
		t.Fatal("LoadFacts should return nil for in-memory storage")
	}
}

func TestIntegrationValidateFacts(t *testing.T) {
	// Create a real manager with mock dependencies
	mockCollector := &MockFactCollector{}
	mockValidator := &MockSchemaValidator{}
	mockLogger := &MockLogger{}

	manager := NewManager(mockCollector, mockValidator, mockLogger)
	integration := NewIntegration(manager)

	facts := &FactCollection{
		MachineID:   "test-server",
		CollectedAt: time.Now(),
		Facts: &spookytypesfacts.Facts{
			System: &spookytypesfacts.SystemFacts{
				OS: &spookytypesfacts.OSFacts{
					Name:    "Linux",
					Version: "Ubuntu 20.04",
				},
			},
		},
	}

	result, err := integration.ValidateFacts(context.Background(), facts)
	if err != nil {
		t.Fatalf("ValidateFacts failed: %v", err)
	}

	if result == nil {
		t.Fatal("ValidateFacts returned nil result")
	}
}

func TestIntegrationGetManager(t *testing.T) {
	// Create a real manager with mock dependencies
	mockCollector := &MockFactCollector{}
	mockValidator := &MockSchemaValidator{}
	mockLogger := &MockLogger{}

	manager := NewManager(mockCollector, mockValidator, mockLogger)
	integration := NewIntegration(manager)

	managerInterface := integration.GetManager()
	if managerInterface == nil {
		t.Fatal("GetManager returned nil")
	}

	// Verify it's the same manager
	if managerInterface != manager {
		t.Fatal("GetManager returned different manager")
	}
}
