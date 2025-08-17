package facts

import (
	"context"
	"testing"
	"time"

	spookytypes "spooky/internal/types"
)

func TestNewIntegration(t *testing.T) {
	mockCollector := &MockSystemFactCollector{}
	mockLogger := &MockLogger{}

	// Create facts manager
	manager := NewManager(mockCollector, mockLogger)

	integration := NewIntegration(manager)

	if integration == nil {
		t.Error("Expected integration to be created, got nil")
	}
}

func TestIntegrationCollectFacts(t *testing.T) {
	mockCollector := NewMockSystemFactCollector()
	mockLogger := &MockLogger{}

	// Create facts manager
	manager := NewManager(mockCollector, mockLogger)
	integration := NewIntegration(manager)

	machine := &spookytypes.Machine{
		Hostname: "test-host",
		Host:     "localhost",
	}

	ctx := context.Background()
	result, err := integration.CollectFacts(ctx, machine)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result == nil {
		t.Error("Expected result, got nil")
	}

	// Verify the result is the correct type
	if _, ok := result.(*spookytypes.FactCollection); !ok {
		t.Errorf("Expected *spookytypes.FactCollection, got %T", result)
	}
}

func TestIntegrationCollectFactsViaSSH(t *testing.T) {
	mockCollector := NewMockSystemFactCollector()
	mockLogger := &MockLogger{}

	// Create facts manager
	manager := NewManager(mockCollector, mockLogger)
	integration := NewIntegration(manager)

	machine := &spookytypes.Machine{
		Hostname: "test-host",
		Host:     "192.168.1.100", // Remote machine
	}

	ctx := context.Background()
	result, err := integration.CollectFacts(ctx, machine)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result == nil {
		t.Error("Expected result, got nil")
	}

	// Verify the result is the correct type
	if _, ok := result.(*spookytypes.FactCollection); !ok {
		t.Errorf("Expected *spookytypes.FactCollection, got %T", result)
	}
}

func TestIntegrationCollectFactsError(t *testing.T) {
	// This test would need a mock that can return errors
	// For now, we'll skip this test since the mock doesn't support error injection
	t.Skip("Mock doesn't support error injection")
}

func TestIntegrationStoreFacts(t *testing.T) {
	mockCollector := NewMockSystemFactCollector()
	mockLogger := &MockLogger{}

	// Create facts manager
	manager := NewManager(mockCollector, mockLogger)
	integration := NewIntegration(manager)

	facts := &spookytypes.FactCollection{
		MachineID:   "test-machine",
		CollectedAt: time.Now(),
		Facts:       &spookytypes.Facts{},
	}

	ctx := context.Background()
	err := integration.StoreFacts(ctx, facts)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestIntegrationStoreFactsNil(t *testing.T) {
	mockCollector := NewMockSystemFactCollector()
	mockLogger := &MockLogger{}

	// Create facts manager
	manager := NewManager(mockCollector, mockLogger)
	integration := NewIntegration(manager)

	ctx := context.Background()
	err := integration.StoreFacts(ctx, nil)

	if err == nil {
		t.Error("Expected error for nil facts, got nil")
	}
}
