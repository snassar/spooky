package facts

import (
	"bytes"
	"io"
	"strings"
	"testing"
	"time"

	"spooky/internal/facts/types"
)

func TestExportToHCL(t *testing.T) {
	// Create test fact collection
	collection := &types.FactCollection{
		Server:    "test-server",
		Timestamp: time.Now(),
		Facts: map[string]*types.Fact{
			"os": {
				Key:   "os",
				Value: "linux",
			},
			"cpu_count": {
				Key:   "cpu_count",
				Value: 4,
			},
		},
	}

	// Create storage with test data
	storage := &MockExporterFactStorage{
		collections: []*types.FactCollection{collection},
	}

	// Create exporter
	exporter := NewExporter(storage)

	// Export to HCL
	var buf bytes.Buffer
	query := &FactQuery{}
	err := exporter.ExportToHCL(&buf, query)
	if err != nil {
		t.Fatalf("ExportToHCL failed: %v", err)
	}

	// Verify HCL output
	output := buf.String()
	if !strings.Contains(output, "facts = [") {
		t.Error("HCL output missing facts array")
	}
	if !strings.Contains(output, "machine_id = \"test-server\"") {
		t.Error("HCL output missing machine_id")
	}
	if !strings.Contains(output, "\"os\" = \"linux\"") {
		t.Error("HCL output missing OS fact")
	}
	if !strings.Contains(output, "\"cpu_count\" = 4") {
		t.Error("HCL output missing CPU fact")
	}
}

func TestFormatHCLValue(t *testing.T) {
	tests := []struct {
		name     string
		value    interface{}
		expected string
	}{
		{"string", "test", "\"test\""},
		{"int", 42, "42"},
		{"float64", 3.14, "3.140000"},
		{"bool", true, "true"},
		{"slice", []interface{}{"a", "b"}, "[\"a\", \"b\"]"},
		{"map", map[string]interface{}{"key": "value"}, "{\n      \"key\" = \"value\"\n    }"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatHCLValue(tt.value)
			if result != tt.expected {
				t.Errorf("formatHCLValue(%v) = %s, want %s", tt.value, result, tt.expected)
			}
		})
	}
}

// MockExporterFactStorage implements FactStorage for testing
type MockExporterFactStorage struct {
	collections []*types.FactCollection
}

func (m *MockExporterFactStorage) SetFactCollection(machineID string, collection *types.FactCollection) error {
	return nil
}

func (m *MockExporterFactStorage) GetFactCollection(machineID string) (*types.FactCollection, error) {
	return nil, nil
}

func (m *MockExporterFactStorage) QueryFactCollections(query *FactQuery) ([]*types.FactCollection, error) {
	return m.collections, nil
}

func (m *MockExporterFactStorage) DeleteFactCollection(machineID string) error {
	return nil
}

func (m *MockExporterFactStorage) DeleteFactCollections(query *FactQuery) (int, error) {
	return 0, nil
}

func (m *MockExporterFactStorage) ImportFromJSON(r io.Reader) error {
	return nil
}

func (m *MockExporterFactStorage) ExportToJSON(w io.Writer) error {
	return nil
}

func (m *MockExporterFactStorage) ImportFromHCL(r io.Reader) error {
	return nil
}

func (m *MockExporterFactStorage) ExportToJSONWithEncryption(w io.Writer, opts types.ExportOptions) error {
	return nil
}

func (m *MockExporterFactStorage) ImportFromJSONWithDecryption(r io.Reader, identityFile string) error {
	return nil
}

func (m *MockExporterFactStorage) Close() error {
	return nil
}

func (m *MockExporterFactStorage) Validate() error {
	return nil
}
