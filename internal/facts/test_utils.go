package facts

import (
	"context"
	"time"

	spookytypes "spooky/internal/types"
	spookytypesfacts "spooky/internal/types/facts"
	spookytypeslogging "spooky/internal/types/logging"
	spookytypesschemas "spooky/internal/types/schemas"
)

// MockFactCollector implements FactCollector for testing
type MockFactCollector struct {
	name string
}

func NewMockFactCollector() *MockFactCollector {
	return &MockFactCollector{
		name: "mock-collector",
	}
}

func (m *MockFactCollector) Collect(_ context.Context, machine *spookytypes.Machine) (*spookytypesfacts.FactCollection, error) {
	return &spookytypesfacts.FactCollection{
		MachineID:   machine.Hostname,
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

func (m *MockSchemaValidator) Validate(_ *spookytypesschemas.Schema, _ interface{}) (*spookytypesschemas.ValidationResult, error) {
	return &spookytypesschemas.ValidationResult{
		Valid:       true,
		ValidatedAt: time.Now(),
		Errors:      []spookytypesschemas.SchemaError{},
		Warnings:    []spookytypesschemas.SchemaError{},
	}, nil
}

func (m *MockSchemaValidator) ValidateFile(_ *spookytypesschemas.Schema, _ string) (*spookytypesschemas.ValidationResult, error) {
	return m.Validate(nil, nil)
}

func (m *MockSchemaValidator) ValidateString(_ *spookytypesschemas.Schema, content string) (*spookytypesschemas.ValidationResult, error) {
	return m.Validate(nil, content)
}

func (m *MockSchemaValidator) ValidateBytes(_ *spookytypesschemas.Schema, data []byte) (*spookytypesschemas.ValidationResult, error) {
	return m.Validate(nil, data)
}

func (m *MockSchemaValidator) ValidateWithContext(_ *spookytypesschemas.Schema, data interface{}, _ map[string]interface{}) (*spookytypesschemas.ValidationResult, error) {
	return m.Validate(nil, data)
}

func (m *MockSchemaValidator) ValidateField(_ *spookytypesschemas.Schema, _ string, _ interface{}) (*spookytypesschemas.ValidationResult, error) {
	return &spookytypesschemas.ValidationResult{
		Valid:       true,
		ValidatedAt: time.Now(),
		Errors:      []spookytypesschemas.SchemaError{},
		Warnings:    []spookytypesschemas.SchemaError{},
	}, nil
}

// MockLogger implements Logger for testing
type MockLogger struct{}

func (m *MockLogger) Debug(_ string, _ ...map[string]interface{})            {}
func (m *MockLogger) Info(_ string, _ ...map[string]interface{})             {}
func (m *MockLogger) Warn(_ string, _ ...map[string]interface{})             {}
func (m *MockLogger) Error(_ string, _ error, _ ...map[string]interface{})   {}
func (m *MockLogger) Fatal(_ string, _ error, _ ...map[string]interface{})   {}
func (m *MockLogger) WithFields(_ map[string]interface{}) spookytypes.Logger { return m }
func (m *MockLogger) WithComponent(_ string) spookytypes.Logger              { return m }
func (m *MockLogger) WithOperation(_ string) spookytypes.Logger              { return m }
func (m *MockLogger) SetLevel(_ spookytypes.LogLevel)                        {}
func (m *MockLogger) GetLevel() spookytypes.LogLevel                         { return spookytypeslogging.LogLevelInfo }
