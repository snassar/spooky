package facts

import (
	"context"
	"strings"
	"testing"
	"time"

	spookytypes "spooky/internal/types"
	spookytypesfacts "spooky/internal/types/facts"
	spookytypesschemas "spooky/internal/types/schemas"
)

// MockFactStorage implements FactStorage for testing
type MockFactStorage struct {
	facts map[string]*FactCollection
}

func NewMockFactStorage() *MockFactStorage {
	return &MockFactStorage{
		facts: make(map[string]*FactCollection),
	}
}

func (m *MockFactStorage) Store(ctx context.Context, machineID string, facts *FactCollection) error {
	m.facts[machineID] = facts
	return nil
}

func (m *MockFactStorage) Get(ctx context.Context, machineID string) (*FactCollection, error) {
	if facts, exists := m.facts[machineID]; exists {
		return facts, nil
	}
	return nil, nil
}

func (m *MockFactStorage) List(ctx context.Context) ([]string, error) {
	var machineIDs []string
	for machineID := range m.facts {
		machineIDs = append(machineIDs, machineID)
	}
	return machineIDs, nil
}

func (m *MockFactStorage) Delete(ctx context.Context, machineID string) error {
	delete(m.facts, machineID)
	return nil
}

func (m *MockFactStorage) Close() error {
	return nil
}

// MockFactCollector implements FactCollector for testing
type MockFactCollector struct {
	name string
}

func NewMockFactCollector() *MockFactCollector {
	return &MockFactCollector{
		name: "mock-collector",
	}
}

func (m *MockFactCollector) Collect(ctx context.Context, machine *spookytypes.Machine) (*FactCollection, error) {
	return &FactCollection{
		MachineID:   "1234567890abcdef1234567890abcdef",
		CollectedAt: time.Now(),
		Facts: &spookytypesfacts.Facts{
			System: &spookytypesfacts.SystemFacts{
				OS: &spookytypesfacts.OSFacts{
					Name:    "TestOS",
					Version: "1.0.0",
					Arch:    "x86_64",
					Kernel:  "5.0.0",
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

func (m *MockFactCollector) GetName() string {
	return m.name
}

// MockSchemaValidator implements SchemaValidator for testing
type MockSchemaValidator struct{}

func (m *MockSchemaValidator) Validate(schema *spookytypesschemas.Schema, data interface{}) (*spookytypesschemas.ValidationResult, error) {
	return &spookytypesschemas.ValidationResult{
		Valid:       true,
		ValidatedAt: time.Now(),
		Errors:      []spookytypesschemas.SchemaError{},
		Warnings:    []spookytypesschemas.SchemaError{},
	}, nil
}

func (m *MockSchemaValidator) ValidateFile(schema *spookytypesschemas.Schema, filePath string) (*spookytypesschemas.ValidationResult, error) {
	return m.Validate(schema, nil)
}

func (m *MockSchemaValidator) ValidateString(schema *spookytypesschemas.Schema, content string) (*spookytypesschemas.ValidationResult, error) {
	return m.Validate(schema, content)
}

func (m *MockSchemaValidator) ValidateBytes(schema *spookytypesschemas.Schema, data []byte) (*spookytypesschemas.ValidationResult, error) {
	return m.Validate(schema, data)
}

// MockLogger implements Logger for testing
type MockLogger struct{}

func (m *MockLogger) Debug(msg string, fields ...map[string]interface{})            {}
func (m *MockLogger) Info(msg string, fields ...map[string]interface{})             {}
func (m *MockLogger) Warn(msg string, fields ...map[string]interface{})             {}
func (m *MockLogger) Error(msg string, err error, fields ...map[string]interface{}) {}
func (m *MockLogger) Fatal(msg string, err error, fields ...map[string]interface{}) {}
func (m *MockLogger) WithFields(fields map[string]interface{}) spookytypes.Logger   { return m }
func (m *MockLogger) WithComponent(component string) spookytypes.Logger             { return m }
func (m *MockLogger) WithOperation(operation string) spookytypes.Logger             { return m }
func (m *MockLogger) SetLevel(level spookytypes.LogLevel)                           {}
func (m *MockLogger) GetLevel() spookytypes.LogLevel                                { return spookytypes.LogLevelInfo }

// createValidTestFacts creates a valid FactCollection for testing
func createValidTestFacts(machineID string) *FactCollection {
	return &FactCollection{
		MachineID:   machineID,
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
	}
}

func TestNewManager(t *testing.T) {
	collector := NewMockFactCollector()
	validator := &MockSchemaValidator{}
	logger := &MockLogger{}

	manager := NewManager(collector, validator, logger)

	if manager == nil {
		t.Fatal("NewManager returned nil")
	}

	if manager.storage == nil {
		t.Error("storage not set correctly")
	}

	if manager.collector != collector {
		t.Error("collector not set correctly")
	}

	if manager.validator != validator {
		t.Error("validator not set correctly")
	}

	if manager.logger != logger {
		t.Error("logger not set correctly")
	}
}

func TestManager_CollectFacts(t *testing.T) {
	collector := NewMockFactCollector()
	validator := &MockSchemaValidator{}
	logger := &MockLogger{}

	manager := NewManager(collector, validator, logger)

	machine := &spookytypes.Machine{
		Hostname: "test-machine",
		Host:     "test-machine",
		Port:     22,
		User:     "test-user",
	}

	ctx := context.Background()
	facts, err := manager.CollectFacts(ctx, machine)

	if err != nil {
		t.Fatalf("CollectFacts failed: %v", err)
	}

	if facts == nil {
		t.Fatal("CollectFacts returned nil facts")
	}

	if facts.MachineID != "1234567890abcdef1234567890abcdef" {
		t.Errorf("expected MachineID '1234567890abcdef1234567890abcdef', got %s", facts.MachineID)
	}

	if facts.Facts == nil {
		t.Error("facts.Facts is nil")
	}

	if facts.Facts.System == nil {
		t.Error("facts.Facts.System is nil")
	}

	if facts.Facts.System.OS == nil {
		t.Error("facts.Facts.System.OS is nil")
	}

	if facts.Facts.System.OS.Name != "TestOS" {
		t.Errorf("expected OS name 'TestOS', got %s", facts.Facts.System.OS.Name)
	}
}

func TestManager_StoreFacts(t *testing.T) {
	collector := NewMockFactCollector()
	validator := &MockSchemaValidator{}
	logger := &MockLogger{}

	manager := NewManager(collector, validator, logger)

	facts := &FactCollection{
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
	}

	ctx := context.Background()
	err := manager.StoreFacts(ctx, "test-machine", facts)

	if err != nil {
		t.Fatalf("StoreFacts failed: %v", err)
	}

	// Verify facts were stored
	storedFacts, err := manager.storage.Get(ctx, "test-machine")
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

func TestManager_GetFacts(t *testing.T) {
	collector := NewMockFactCollector()
	validator := &MockSchemaValidator{}
	logger := &MockLogger{}

	manager := NewManager(collector, validator, logger)

	// Store some facts first
	facts := createValidTestFacts("1234567890abcdef1234567890abcdef")

	ctx := context.Background()
	err := manager.StoreFacts(ctx, "1234567890abcdef1234567890abcdef", facts)
	if err != nil {
		t.Fatalf("Failed to store facts: %v", err)
	}

	// Retrieve facts
	retrievedFacts, err := manager.GetFacts(ctx, "1234567890abcdef1234567890abcdef")
	if err != nil {
		t.Fatalf("GetFacts failed: %v", err)
	}

	if retrievedFacts == nil {
		t.Fatal("GetFacts returned nil")
	}

	if retrievedFacts.MachineID != "1234567890abcdef1234567890abcdef" {
		t.Errorf("expected MachineID '1234567890abcdef1234567890abcdef', got %s", retrievedFacts.MachineID)
	}
}

func TestManager_ListFacts(t *testing.T) {
	collector := NewMockFactCollector()
	validator := &MockSchemaValidator{}
	logger := &MockLogger{}

	manager := NewManager(collector, validator, logger)

	// Store facts for multiple machines
	ctx := context.Background()

	facts1 := createValidTestFacts("11111111111111111111111111111111")
	facts2 := createValidTestFacts("22222222222222222222222222222222")

	err := manager.StoreFacts(ctx, "11111111111111111111111111111111", facts1)
	if err != nil {
		t.Fatalf("Failed to store facts for machine-1: %v", err)
	}

	err = manager.StoreFacts(ctx, "22222222222222222222222222222222", facts2)
	if err != nil {
		t.Fatalf("Failed to store facts for machine-2: %v", err)
	}

	// List facts
	machineIDs, err := manager.ListFacts(ctx)
	if err != nil {
		t.Fatalf("ListFacts failed: %v", err)
	}

	if len(machineIDs) != 2 {
		t.Errorf("expected 2 machine IDs, got %d", len(machineIDs))
	}

	// Check that both machine IDs are present
	found1, found2 := false, false
	for _, id := range machineIDs {
		if id == "11111111111111111111111111111111" {
			found1 = true
		}
		if id == "22222222222222222222222222222222" {
			found2 = true
		}
	}

	if !found1 {
		t.Error("machine-1 not found in list")
	}

	if !found2 {
		t.Error("machine-2 not found in list")
	}
}

func TestManager_DeleteFacts(t *testing.T) {
	collector := NewMockFactCollector()
	validator := &MockSchemaValidator{}
	logger := &MockLogger{}

	manager := NewManager(collector, validator, logger)

	// Store facts first
	facts := createValidTestFacts("1234567890abcdef1234567890abcdef")

	ctx := context.Background()
	err := manager.StoreFacts(ctx, "1234567890abcdef1234567890abcdef", facts)
	if err != nil {
		t.Fatalf("Failed to store facts: %v", err)
	}

	// Verify facts exist
	storedFacts, err := manager.GetFacts(ctx, "1234567890abcdef1234567890abcdef")
	if err != nil {
		t.Fatalf("Failed to get facts: %v", err)
	}

	if storedFacts == nil {
		t.Fatal("Facts not found after storing")
	}

	// Delete facts
	err = manager.DeleteFacts(ctx, "1234567890abcdef1234567890abcdef")
	if err != nil {
		t.Fatalf("DeleteFacts failed: %v", err)
	}

	// Verify facts are deleted
	_, err = manager.GetFacts(ctx, "1234567890abcdef1234567890abcdef")
	if err == nil {
		t.Error("Expected error when getting deleted facts")
	}

	if !strings.Contains(err.Error(), "facts not found") {
		t.Errorf("Expected 'facts not found' error, got: %v", err)
	}
}

func TestManager_ValidateFacts(t *testing.T) {
	collector := NewMockFactCollector()
	validator := &MockSchemaValidator{}
	logger := &MockLogger{}

	manager := NewManager(collector, validator, logger)

	facts := createValidTestFacts("1234567890abcdef1234567890abcdef")

	ctx := context.Background()
	result, err := manager.ValidateFacts(ctx, facts)

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

func TestManager_CollectFacts_NilMachine(t *testing.T) {
	collector := NewMockFactCollector()
	validator := &MockSchemaValidator{}
	logger := &MockLogger{}

	manager := NewManager(collector, validator, logger)

	ctx := context.Background()
	_, err := manager.CollectFacts(ctx, nil)

	if err == nil {
		t.Error("Expected error when machine is nil")
	}

	if err.Error() != "machine cannot be nil" {
		t.Errorf("Expected 'machine cannot be nil' error, got: %v", err)
	}
}

func TestManager_StoreFacts_EmptyMachineID(t *testing.T) {
	collector := NewMockFactCollector()
	validator := &MockSchemaValidator{}
	logger := &MockLogger{}

	manager := NewManager(collector, validator, logger)

	facts := &FactCollection{
		MachineID:   "test-machine",
		CollectedAt: time.Now(),
		Facts:       &spookytypesfacts.Facts{},
		Metadata:    make(map[string]interface{}),
	}

	ctx := context.Background()
	err := manager.StoreFacts(ctx, "", facts)

	if err == nil {
		t.Error("Expected error when machine ID is empty")
	}

	if err.Error() != "machine ID cannot be empty" {
		t.Errorf("Expected 'machine ID cannot be empty' error, got: %v", err)
	}
}

func TestManager_StoreFacts_NilFacts(t *testing.T) {
	collector := NewMockFactCollector()
	validator := &MockSchemaValidator{}
	logger := &MockLogger{}

	manager := NewManager(collector, validator, logger)

	ctx := context.Background()
	err := manager.StoreFacts(ctx, "test-machine", nil)

	if err == nil {
		t.Error("Expected error when facts is nil")
	}

	if err.Error() != "facts cannot be nil" {
		t.Errorf("Expected 'facts cannot be nil' error, got: %v", err)
	}
}

func TestManager_GetFacts_EmptyMachineID(t *testing.T) {
	collector := NewMockFactCollector()
	validator := &MockSchemaValidator{}
	logger := &MockLogger{}

	manager := NewManager(collector, validator, logger)

	ctx := context.Background()
	_, err := manager.GetFacts(ctx, "")

	if err == nil {
		t.Error("Expected error when machine ID is empty")
	}

	if err.Error() != "machine ID cannot be empty" {
		t.Errorf("Expected 'machine ID cannot be empty' error, got: %v", err)
	}
}

func TestManager_ExportFacts(t *testing.T) {
	collector := NewMockFactCollector()
	validator := &MockSchemaValidator{}
	logger := &MockLogger{}

	manager := NewManager(collector, validator, logger)

	// Create test facts with proper structure
	facts := &FactCollection{
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
						Total: 8589934592,
					},
				},
				Network: &spookytypesfacts.NetworkFacts{
					Hostname: "test-host",
					Interfaces: []*spookytypesfacts.NetworkInterface{
						{
							Name:        "eth0",
							IPAddresses: []string{"192.168.1.100"},
						},
					},
				},
			},
		},
		Metadata: make(map[string]interface{}),
	}

	ctx := context.Background()

	// Store facts first
	if err := manager.StoreFacts(ctx, facts.MachineID, facts); err != nil {
		t.Fatalf("Failed to store facts: %v", err)
	}

	// Test JSON export
	if err := manager.ExportFacts(ctx, []string{facts.MachineID}, "json", "test-export.json"); err != nil {
		t.Fatalf("Failed to export JSON: %v", err)
	}

	// Test HCL export
	if err := manager.ExportFacts(ctx, []string{facts.MachineID}, "hcl", "test-export.hcl"); err != nil {
		t.Fatalf("Failed to export HCL: %v", err)
	}

	// Test unsupported format
	if err := manager.ExportFacts(ctx, []string{facts.MachineID}, "xml", "test-export.xml"); err == nil {
		t.Error("Expected error for unsupported format")
	}

	// Note: Files are cleaned up after test completion
	// Uncomment the following lines to keep files for inspection:
	// fmt.Printf("JSON export file: test-export.json\n")
	// fmt.Printf("HCL export file: test-export.hcl\n")
}
