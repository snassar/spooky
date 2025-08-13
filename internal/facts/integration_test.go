package facts

import (
	"context"
	"testing"
	"time"

	spookytypes "spooky/internal/types"
	spookytypesfacts "spooky/internal/types/facts"
	spookytypesschemas "spooky/internal/types/schemas"
)

// MockFactManager implements FactManager for testing
type MockFactManager struct {
	facts map[string]*FactCollection
}

func NewMockFactManager() *MockFactManager {
	return &MockFactManager{
		facts: make(map[string]*FactCollection),
	}
}

func (m *MockFactManager) CollectFacts(ctx context.Context, machine *spookytypes.Machine) (*FactCollection, error) {
	return &FactCollection{
		MachineID:   "1234567890abcdef1234567890abcdef",
		CollectedAt: time.Now(),
		Facts: &spookytypesfacts.Facts{
			System: &spookytypesfacts.SystemFacts{
				OS: &spookytypesfacts.OSFacts{
					Name:    "TestOS",
					Version: "1.0.0",
				},
				Hardware: &spookytypesfacts.HardwareFacts{
					CPU: &spookytypesfacts.CPUFacts{
						Cores: 4,
						Model: "Test CPU",
					},
					Memory: &spookytypesfacts.MemoryFacts{
						Total: 8589934592, // 8GB
					},
				},
				Network: &spookytypesfacts.NetworkFacts{
					Hostname:    "test-host",
					IPAddresses: []string{"192.168.1.100"},
					PrimaryIP:   "192.168.1.100",
				},
			},
		},
		Metadata: make(map[string]interface{}),
	}, nil
}

func (m *MockFactManager) StoreFacts(ctx context.Context, machineID string, facts *FactCollection) error {
	m.facts[machineID] = facts
	return nil
}

func (m *MockFactManager) GetFacts(ctx context.Context, machineID string) (*FactCollection, error) {
	if facts, exists := m.facts[machineID]; exists {
		return facts, nil
	}
	return nil, nil
}

func (m *MockFactManager) ListFacts(ctx context.Context) ([]string, error) {
	var machineIDs []string
	for machineID := range m.facts {
		machineIDs = append(machineIDs, machineID)
	}
	return machineIDs, nil
}

func (m *MockFactManager) DeleteFacts(ctx context.Context, machineID string) error {
	delete(m.facts, machineID)
	return nil
}

func (m *MockFactManager) ValidateFacts(ctx context.Context, facts *FactCollection) (*spookytypes.ValidationResult, error) {
	return &spookytypes.ValidationResult{
		Valid:    true,
		Errors:   []spookytypesschemas.SchemaError{},
		Warnings: []spookytypesschemas.SchemaError{},
	}, nil
}

func (m *MockFactManager) ExportFacts(ctx context.Context, machineIDs []string, format string, outputPath string) error {
	return nil
}

func (m *MockFactManager) ImportFacts(ctx context.Context, format string, inputPath string) error {
	return nil
}

func TestNewIntegration(t *testing.T) {
	mockManager := NewMockFactManager()
	integration := NewIntegration(mockManager)

	if integration == nil {
		t.Fatal("NewIntegration returned nil")
	}

	if integration.manager != mockManager {
		t.Error("manager not set correctly")
	}
}

func TestIntegration_CollectFacts(t *testing.T) {
	mockManager := NewMockFactManager()
	integration := NewIntegration(mockManager)

	ctx := context.Background()
	facts, err := integration.CollectFacts(ctx, "test-machine")

	if err != nil {
		t.Fatalf("CollectFacts failed: %v", err)
	}

	if facts == nil {
		t.Fatal("CollectFacts returned nil")
	}

	// Type assert to check the facts
	factCollection, ok := facts.(*FactCollection)
	if !ok {
		t.Fatal("facts is not of type *FactCollection")
	}

	if factCollection.MachineID != "1234567890abcdef1234567890abcdef" {
		t.Errorf("expected MachineID '1234567890abcdef1234567890abcdef', got %s", factCollection.MachineID)
	}

	if factCollection.Facts == nil {
		t.Error("facts.Facts is nil")
	}

	if factCollection.Facts.System == nil {
		t.Error("facts.Facts.System is nil")
	}

	if factCollection.Facts.System.OS == nil {
		t.Error("facts.Facts.System.OS is nil")
	}

	if factCollection.Facts.System.OS.Name != "TestOS" {
		t.Errorf("expected OS name 'TestOS', got %s", factCollection.Facts.System.OS.Name)
	}
}

func TestIntegration_StoreFacts(t *testing.T) {
	mockManager := NewMockFactManager()
	integration := NewIntegration(mockManager)

	facts := createValidTestFacts("1234567890abcdef1234567890abcdef")

	ctx := context.Background()
	err := integration.StoreFacts(ctx, facts, nil)

	if err != nil {
		t.Fatalf("StoreFacts failed: %v", err)
	}

	// Verify facts were stored
	storedFacts, err := mockManager.GetFacts(ctx, "1234567890abcdef1234567890abcdef")
	if err != nil {
		t.Fatalf("Failed to get stored facts: %v", err)
	}

	if storedFacts == nil {
		t.Fatal("Stored facts are nil")
	}

	if storedFacts.MachineID != "1234567890abcdef1234567890abcdef" {
		t.Errorf("expected MachineID '1234567890abcdef1234567890abcdef', got %s", storedFacts.MachineID)
	}
}

func TestIntegration_StoreFacts_NilFacts(t *testing.T) {
	mockManager := NewMockFactManager()
	integration := NewIntegration(mockManager)

	ctx := context.Background()
	err := integration.StoreFacts(ctx, nil, nil)

	if err == nil {
		t.Error("Expected error when facts is nil")
	}

	if err.Error() != "facts cannot be nil" {
		t.Errorf("Expected 'facts cannot be nil' error, got: %v", err)
	}
}

func TestIntegration_StoreFacts_InvalidType(t *testing.T) {
	mockManager := NewMockFactManager()
	integration := NewIntegration(mockManager)

	ctx := context.Background()
	err := integration.StoreFacts(ctx, "invalid-type", nil)

	if err == nil {
		t.Error("Expected error when facts is invalid type")
	}

	if err.Error() != "invalid facts type" {
		t.Errorf("Expected 'invalid facts type' error, got: %v", err)
	}
}

func TestIntegration_LoadFacts(t *testing.T) {
	mockManager := NewMockFactManager()
	integration := NewIntegration(mockManager)

	// Store some facts first
	facts := createValidTestFacts("1234567890abcdef1234567890abcdef")
	ctx := context.Background()
	err := mockManager.StoreFacts(ctx, "1234567890abcdef1234567890abcdef", facts)
	if err != nil {
		t.Fatalf("Failed to store facts: %v", err)
	}

	// Load facts
	loadedFacts, err := integration.LoadFacts(ctx, nil)
	if err != nil {
		t.Fatalf("LoadFacts failed: %v", err)
	}

	if loadedFacts == nil {
		t.Fatal("LoadFacts returned nil")
	}

	// Type assert to check the facts
	factCollection, ok := loadedFacts.(*FactCollection)
	if !ok {
		t.Fatal("loadedFacts is not of type *FactCollection")
	}

	if factCollection.MachineID != "1234567890abcdef1234567890abcdef" {
		t.Errorf("expected MachineID '1234567890abcdef1234567890abcdef', got %s", factCollection.MachineID)
	}
}

func TestIntegration_LoadFacts_Empty(t *testing.T) {
	mockManager := NewMockFactManager()
	integration := NewIntegration(mockManager)

	ctx := context.Background()
	loadedFacts, err := integration.LoadFacts(ctx, nil)
	if err != nil {
		t.Fatalf("LoadFacts failed: %v", err)
	}

	if loadedFacts != nil {
		t.Error("LoadFacts should return nil when no facts exist")
	}
}

func TestIntegration_ValidateFacts(t *testing.T) {
	mockManager := NewMockFactManager()
	integration := NewIntegration(mockManager)

	facts := createValidTestFacts("1234567890abcdef1234567890abcdef")

	ctx := context.Background()
	result, err := integration.ValidateFacts(ctx, facts)

	if err != nil {
		t.Fatalf("ValidateFacts failed: %v", err)
	}

	if result == nil {
		t.Fatal("ValidationResult is nil")
	}

	if !result.Valid {
		t.Error("Validation should pass with mock validator")
	}
}

func TestIntegration_ValidateFacts_NilFacts(t *testing.T) {
	mockManager := NewMockFactManager()
	integration := NewIntegration(mockManager)

	ctx := context.Background()
	result, err := integration.ValidateFacts(ctx, nil)

	if err != nil {
		t.Fatalf("ValidateFacts failed: %v", err)
	}

	if result == nil {
		t.Fatal("ValidationResult is nil")
	}

	if result.Valid {
		t.Error("Validation should fail when facts is nil")
	}

	if len(result.Errors) == 0 {
		t.Error("Validation should have errors when facts is nil")
	}
}

func TestIntegration_ValidateFacts_InvalidType(t *testing.T) {
	mockManager := NewMockFactManager()
	integration := NewIntegration(mockManager)

	ctx := context.Background()
	result, err := integration.ValidateFacts(ctx, "invalid-type")

	if err != nil {
		t.Fatalf("ValidateFacts failed: %v", err)
	}

	if result == nil {
		t.Fatal("ValidationResult is nil")
	}

	if result.Valid {
		t.Error("Validation should fail when facts is invalid type")
	}

	if len(result.Errors) == 0 {
		t.Error("Validation should have errors when facts is invalid type")
	}
}

func TestIntegration_ExportFacts(t *testing.T) {
	mockManager := NewMockFactManager()
	integration := NewIntegration(mockManager)

	facts := createValidTestFacts("1234567890abcdef1234567890abcdef")

	ctx := context.Background()
	err := integration.ExportFacts(ctx, facts, "json", "/tmp/test-facts.json")

	if err != nil {
		t.Fatalf("ExportFacts failed: %v", err)
	}
}

func TestIntegration_ExportFacts_NilFacts(t *testing.T) {
	mockManager := NewMockFactManager()
	integration := NewIntegration(mockManager)

	ctx := context.Background()
	err := integration.ExportFacts(ctx, nil, "json", "/tmp/test-facts.json")

	if err == nil {
		t.Error("Expected error when facts is nil")
	}

	if err.Error() != "facts cannot be nil" {
		t.Errorf("Expected 'facts cannot be nil' error, got: %v", err)
	}
}

func TestIntegration_ExportFacts_UnsupportedFormat(t *testing.T) {
	mockManager := NewMockFactManager()
	integration := NewIntegration(mockManager)

	facts := createValidTestFacts("1234567890abcdef1234567890abcdef")

	ctx := context.Background()
	err := integration.ExportFacts(ctx, facts, "unsupported", "/tmp/test-facts.unsupported")

	if err == nil {
		t.Error("Expected error for unsupported format")
	}

	if err.Error() != "unsupported export format: unsupported" {
		t.Errorf("Expected 'unsupported export format' error, got: %v", err)
	}
}
