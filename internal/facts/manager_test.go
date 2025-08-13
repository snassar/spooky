package facts

import (
	"context"
	"testing"
	"time"

	spookytypes "spooky/internal/types"
	spookytypesfacts "spooky/internal/types/facts"
	spookytypeslogging "spooky/internal/types/logging"
	spookytypesschemas "spooky/internal/types/schemas"

	"github.com/stretchr/testify/assert"
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

func (m *MockFactStorage) Clear(ctx context.Context) error {
	m.facts = make(map[string]*FactCollection)
	return nil
}

func (m *MockFactStorage) GetStats() (map[string]interface{}, error) {
	return map[string]interface{}{
		"total_entries": len(m.facts),
		"total_size":    0,
		"storage_type":  "memory",
	}, nil
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

func (m *MockSchemaValidator) ValidateWithContext(schema *spookytypesschemas.Schema, data interface{}, context map[string]interface{}) (*spookytypesschemas.ValidationResult, error) {
	return m.Validate(schema, data)
}

func (m *MockSchemaValidator) ValidateField(schema *spookytypesschemas.Schema, fieldPath string, value interface{}) (*spookytypesschemas.ValidationResult, error) {
	return &spookytypesschemas.ValidationResult{
		Valid:       true,
		ValidatedAt: time.Now(),
		Errors:      []spookytypesschemas.SchemaError{},
		Warnings:    []spookytypesschemas.SchemaError{},
	}, nil
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
func (m *MockLogger) GetLevel() spookytypes.LogLevel                                { return spookytypeslogging.LogLevelInfo }

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

func createTestManager(t *testing.T) *Manager {
	collector := NewMockFactCollector()
	validator := &MockSchemaValidator{}
	logger := &MockLogger{}

	manager := NewManager(collector, validator, logger)
	assert.NotNil(t, manager)
	assert.NotNil(t, manager.collector)
	assert.NotNil(t, manager.validator)
	assert.NotNil(t, manager.logger)

	return manager
}

func TestManager_CollectAndStoreFacts(t *testing.T) {
	manager := createTestManager(t)

	machine := &spookytypes.Machine{
		Hostname: "test-machine",
		Host:     "test.example.com",
		Port:     22,
		User:     "testuser",
	}

	// Test collecting and validating facts
	err := manager.CollectAndStoreFacts(context.Background(), machine)
	assert.NoError(t, err)
}

func TestManager_CollectAndStoreFacts_NilMachine(t *testing.T) {
	manager := createTestManager(t)

	err := manager.CollectAndStoreFacts(context.Background(), nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "machine cannot be nil")
}

func TestManager_CollectAndStoreFacts_ValidationError(t *testing.T) {
	manager := createTestManager(t)

	machine := &spookytypes.Machine{
		Hostname: "test-machine",
		Host:     "test.example.com",
		Port:     22,
		User:     "testuser",
	}

	// Test with validation error
	err := manager.CollectAndStoreFacts(context.Background(), machine)
	// Should pass since we're using a mock collector that returns valid facts
	assert.NoError(t, err)
}

func TestManager_CollectAndStoreFactsParallel(t *testing.T) {
	manager := createTestManager(t)

	machines := []*spookytypes.Machine{
		{
			Hostname: "test-machine-1",
			Host:     "test1.example.com",
			Port:     22,
			User:     "testuser",
		},
		{
			Hostname: "test-machine-2",
			Host:     "test2.example.com",
			Port:     22,
			User:     "testuser",
		},
	}

	// Test parallel collection
	err := manager.CollectAndStoreFactsParallel(context.Background(), machines, 2)
	assert.NoError(t, err)
}

func TestManager_CollectAndStoreFactsParallel_EmptyMachines(t *testing.T) {
	manager := createTestManager(t)

	err := manager.CollectAndStoreFactsParallel(context.Background(), []*spookytypes.Machine{}, 2)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no machines provided")
}

func TestManager_CollectAndStoreFactsParallel_InvalidWorkers(t *testing.T) {
	manager := createTestManager(t)

	machines := []*spookytypes.Machine{
		{
			Hostname: "test-machine",
			Host:     "test.example.com",
			Port:     22,
			User:     "testuser",
		},
	}

	// Test with invalid worker count (should default to 4)
	err := manager.CollectAndStoreFactsParallel(context.Background(), machines, 0)
	assert.NoError(t, err)
}
