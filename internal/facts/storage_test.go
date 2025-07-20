package facts

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"
	"time"
)

func TestBadgerFactStorage(t *testing.T) {
	// Create temporary directory for test database
	tempDir, err := os.MkdirTemp("", "badger-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	dbPath := tempDir + "/test.db"

	// Test NewBadgerFactStorage
	storage, err := NewBadgerFactStorage(dbPath)
	if err != nil {
		t.Fatalf("Failed to create BadgerFactStorage: %v", err)
	}
	defer storage.Close()

	// Test data
	machineID := "test-machine-1"
	facts := &MachineFacts{
		MachineID:   machineID,
		MachineName: "test-machine",
		Hostname:    "test-host",
		OS:          "Linux",
		OSVersion:   "Ubuntu 22.04",
		IPAddresses: []string{"192.168.1.100", "10.0.0.1"},
		PrimaryIP:   "192.168.1.100",
		CPU: CPUInfo{
			Cores:     4,
			Model:     "Intel i7",
			Arch:      "x86_64",
			Frequency: "2.4 GHz",
		},
		Memory: MemoryInfo{
			Total:     8589934592, // 8GB
			Used:      4294967296, // 4GB
			Available: 4294967296, // 4GB
		},
		SystemID:  "test-system-id",
		Tags:      map[string]string{"environment": "test", "team": "dev"},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Test SetMachineFacts
	err = storage.SetMachineFacts(machineID, facts)
	if err != nil {
		t.Fatalf("Failed to set machine facts: %v", err)
	}

	// Test GetMachineFacts
	retrievedFacts, err := storage.GetMachineFacts(machineID)
	if err != nil {
		t.Fatalf("Failed to get machine facts: %v", err)
	}

	if retrievedFacts == nil {
		t.Fatal("Retrieved facts is nil")
	}

	// Verify facts match
	if retrievedFacts.MachineID != facts.MachineID {
		t.Errorf("MachineID mismatch: expected %s, got %s", facts.MachineID, retrievedFacts.MachineID)
	}
	if retrievedFacts.Hostname != facts.Hostname {
		t.Errorf("Hostname mismatch: expected %s, got %s", facts.Hostname, retrievedFacts.Hostname)
	}
	if retrievedFacts.CPU.Cores != facts.CPU.Cores {
		t.Errorf("CPU cores mismatch: expected %d, got %d", facts.CPU.Cores, retrievedFacts.CPU.Cores)
	}

	// Test QueryFacts
	query := &FactQuery{
		MachineName: "test-machine",
		OS:          "Linux",
		Limit:       10,
	}
	results, err := storage.QueryFacts(query)
	if err != nil {
		t.Fatalf("Failed to query facts: %v", err)
	}

	if len(results) != 1 {
		t.Errorf("Expected 1 result, got %d", len(results))
	}

	// Test query with no matches
	noMatchQuery := &FactQuery{
		MachineName: "non-existent",
	}
	noResults, err := storage.QueryFacts(noMatchQuery)
	if err != nil {
		t.Fatalf("Failed to query facts: %v", err)
	}

	if len(noResults) != 0 {
		t.Errorf("Expected 0 results, got %d", len(noResults))
	}

	// Test DeleteFacts
	deletedCount, err := storage.DeleteFacts(query)
	if err != nil {
		t.Fatalf("Failed to delete facts: %v", err)
	}

	if deletedCount != 1 {
		t.Errorf("Expected 1 deleted fact, got %d", deletedCount)
	}

	// Verify fact is deleted
	_, err = storage.GetMachineFacts(machineID)
	if err == nil {
		t.Error("Expected error when getting deleted facts")
	}
}

func TestJSONFactStorage(t *testing.T) {
	// Create temporary file for test
	tempFile, err := os.CreateTemp("", "json-test-*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	// Write empty JSON object to file
	_, err = tempFile.WriteString("{}")
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tempFile.Close()

	// Test NewJSONFactStorage
	storage, err := NewJSONFactStorage(tempFile.Name())
	if err != nil {
		t.Fatalf("Failed to create JSONFactStorage: %v", err)
	}
	defer storage.Close()

	// Test data
	machineID := "test-machine-2"
	facts := &MachineFacts{
		MachineID:   machineID,
		MachineName: "test-machine-2",
		Hostname:    "test-host-2",
		OS:          "Linux",
		OSVersion:   "CentOS 8",
		IPAddresses: []string{"192.168.1.101", "10.0.0.2"},
		PrimaryIP:   "192.168.1.101",
		CPU: CPUInfo{
			Cores:     8,
			Model:     "AMD Ryzen",
			Arch:      "x86_64",
			Frequency: "3.2 GHz",
		},
		Memory: MemoryInfo{
			Total:     17179869184, // 16GB
			Used:      8589934592,  // 8GB
			Available: 8589934592,  // 8GB
		},
		SystemID:  "test-system-id-2",
		Tags:      map[string]string{"environment": "staging", "team": "qa"},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Test SetMachineFacts
	err = storage.SetMachineFacts(machineID, facts)
	if err != nil {
		t.Fatalf("Failed to set machine facts: %v", err)
	}

	// Test GetMachineFacts
	retrievedFacts, err := storage.GetMachineFacts(machineID)
	if err != nil {
		t.Fatalf("Failed to get machine facts: %v", err)
	}

	if retrievedFacts == nil {
		t.Fatal("Retrieved facts is nil")
	}

	// Verify facts match
	if retrievedFacts.MachineID != facts.MachineID {
		t.Errorf("MachineID mismatch: expected %s, got %s", facts.MachineID, retrievedFacts.MachineID)
	}
	if retrievedFacts.CPU.Cores != facts.CPU.Cores {
		t.Errorf("CPU cores mismatch: expected %d, got %d", facts.CPU.Cores, retrievedFacts.CPU.Cores)
	}

	// Test QueryFacts
	query := &FactQuery{
		OS:    "Linux",
		Limit: 10,
	}
	results, err := storage.QueryFacts(query)
	if err != nil {
		t.Fatalf("Failed to query facts: %v", err)
	}

	if len(results) != 1 {
		t.Errorf("Expected 1 result, got %d", len(results))
	}

	// Test DeleteFacts
	deletedCount, err := storage.DeleteFacts(query)
	if err != nil {
		t.Fatalf("Failed to delete facts: %v", err)
	}

	if deletedCount != 1 {
		t.Errorf("Expected 1 deleted fact, got %d", deletedCount)
	}

	// Verify fact is deleted
	_, err = storage.GetMachineFacts(machineID)
	if err == nil {
		t.Error("Expected error when getting deleted facts")
	}
}

func TestFactQueryMatching(t *testing.T) {
	facts := &MachineFacts{
		MachineID:   "test-machine",
		MachineName: "test-machine",
		ActionFile:  "/path/to/action.hcl",
		ProjectName: "test-project",
		OS:          "Linux",
		Tags: map[string]string{
			"environment": "production",
			"team":        "ops",
		},
		UpdatedAt: time.Now(),
	}

	tests := []struct {
		name     string
		query    *FactQuery
		expected bool
	}{
		{
			name: "exact machine name match",
			query: &FactQuery{
				MachineName: "test-machine",
			},
			expected: true,
		},
		{
			name: "wrong machine name",
			query: &FactQuery{
				MachineName: "wrong-machine",
			},
			expected: false,
		},
		{
			name: "exact OS match",
			query: &FactQuery{
				OS: "Linux",
			},
			expected: true,
		},
		{
			name: "wrong OS",
			query: &FactQuery{
				OS: "Windows",
			},
			expected: false,
		},
		{
			name: "environment tag match",
			query: &FactQuery{
				Environment: "production",
			},
			expected: true,
		},
		{
			name: "wrong environment",
			query: &FactQuery{
				Environment: "development",
			},
			expected: false,
		},
		{
			name: "multiple tag match",
			query: &FactQuery{
				Tags: map[string]string{
					"environment": "production",
					"team":        "ops",
				},
			},
			expected: true,
		},
		{
			name: "partial tag match",
			query: &FactQuery{
				Tags: map[string]string{
					"environment": "production",
					"team":        "wrong-team",
				},
			},
			expected: false,
		},
		{
			name: "time filter - updated before",
			query: &FactQuery{
				UpdatedBefore: func() *time.Time { t := time.Now().Add(time.Hour); return &t }(),
			},
			expected: true,
		},
		{
			name: "time filter - updated after",
			query: &FactQuery{
				UpdatedAfter: func() *time.Time { t := time.Now().Add(-time.Hour); return &t }(),
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := matchesQuery(facts, tt.query)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestExportImportJSON(t *testing.T) {
	// Create temporary file for test
	tempFile, err := os.CreateTemp("", "export-test-*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	// Write empty JSON object to file
	_, err = tempFile.WriteString("{}")
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tempFile.Close()

	// Create storage
	storage, err := NewJSONFactStorage(tempFile.Name())
	if err != nil {
		t.Fatalf("Failed to create JSONFactStorage: %v", err)
	}
	defer storage.Close()

	// Add test data
	facts1 := &MachineFacts{
		MachineID:   "machine-1",
		MachineName: "test-1",
		Hostname:    "host-1",
		OS:          "Linux",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	facts2 := &MachineFacts{
		MachineID:   "machine-2",
		MachineName: "test-2",
		Hostname:    "host-2",
		OS:          "Windows",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err = storage.SetMachineFacts("machine-1", facts1)
	if err != nil {
		t.Fatalf("Failed to set facts1: %v", err)
	}

	err = storage.SetMachineFacts("machine-2", facts2)
	if err != nil {
		t.Fatalf("Failed to set facts2: %v", err)
	}

	// Test ExportToJSON
	var buf bytes.Buffer
	err = storage.ExportToJSON(&buf)
	if err != nil {
		t.Fatalf("Failed to export to JSON: %v", err)
	}

	// Verify exported data
	var exportedData map[string]*MachineFacts
	err = json.Unmarshal(buf.Bytes(), &exportedData)
	if err != nil {
		t.Fatalf("Failed to unmarshal exported JSON: %v", err)
	}

	if len(exportedData) != 2 {
		t.Errorf("Expected 2 machines in export, got %d", len(exportedData))
	}

	// Test ImportFromJSON
	newStorage, err := NewJSONFactStorage(tempFile.Name() + ".import")
	if err != nil {
		t.Fatalf("Failed to create new storage: %v", err)
	}
	defer newStorage.Close()

	err = newStorage.ImportFromJSON(&buf)
	if err != nil {
		t.Fatalf("Failed to import from JSON: %v", err)
	}

	// Verify imported data
	importedFacts1, err := newStorage.GetMachineFacts("machine-1")
	if err != nil {
		t.Fatalf("Failed to get imported facts1: %v", err)
	}

	if importedFacts1.MachineName != facts1.MachineName {
		t.Errorf("Imported machine name mismatch: expected %s, got %s", facts1.MachineName, importedFacts1.MachineName)
	}
}

func TestStorageInterface(_ *testing.T) {
	// Test that both storage implementations satisfy the interface
	var _ FactStorage = (*BadgerFactStorage)(nil)
	var _ FactStorage = (*JSONFactStorage)(nil)
}

func TestNewFactStorage(t *testing.T) {
	// Test BadgerDB storage creation
	tempDir, err := os.MkdirTemp("", "storage-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Test BadgerDB
	badgerStorage, err := NewFactStorage(StorageOptions{
		Type: StorageTypeBadger,
		Path: tempDir + "/badger.db",
	})
	if err != nil {
		t.Fatalf("Failed to create BadgerDB storage: %v", err)
	}
	defer badgerStorage.Close()

	// Test JSON storage
	tempFile, err := os.CreateTemp("", "json-storage-test-*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	// Write empty JSON object to file
	_, err = tempFile.WriteString("{}")
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tempFile.Close()

	jsonStorage, err := NewFactStorage(StorageOptions{
		Type: StorageTypeJSON,
		Path: tempFile.Name(),
	})
	if err != nil {
		t.Fatalf("Failed to create JSON storage: %v", err)
	}
	defer jsonStorage.Close()

	// Test invalid storage type
	_, err = NewFactStorage(StorageOptions{
		Type: "invalid",
		Path: "/tmp/test",
	})
	if err == nil {
		t.Error("Expected error for invalid storage type")
	}
}

func TestFactConversion(t *testing.T) {
	// Test data
	machineID := "test-conversion"
	collection := &FactCollection{
		Server:    "test-server",
		Timestamp: time.Now(),
		Facts: map[string]*Fact{
			"hostname": {
				Key:       "hostname",
				Value:     "test-host",
				Source:    string(SourceLocal),
				Timestamp: time.Now(),
			},
			"os.name": {
				Key:       "os.name",
				Value:     "Linux",
				Source:    string(SourceLocal),
				Timestamp: time.Now(),
			},
			"cpu.cores": {
				Key:       "cpu.cores",
				Value:     8,
				Source:    string(SourceLocal),
				Timestamp: time.Now(),
			},
			"cpu.frequency": {
				Key:       "cpu.frequency",
				Value:     "2.4 GHz",
				Source:    string(SourceLocal),
				Timestamp: time.Now(),
			},
			"network.ips": {
				Key:       "network.ips",
				Value:     []string{"192.168.1.100", "10.0.0.1"},
				Source:    string(SourceLocal),
				Timestamp: time.Now(),
			},
		},
	}

	// Test conversion to MachineFacts
	machineFacts := ConvertFactCollectionToMachineFacts(machineID, collection)
	if machineFacts == nil {
		t.Fatal("ConvertFactCollectionToMachineFacts returned nil")
	}

	if machineFacts.MachineID != machineID {
		t.Errorf("MachineID mismatch: expected %s, got %s", machineID, machineFacts.MachineID)
	}

	if machineFacts.Hostname != "test-host" {
		t.Errorf("Hostname mismatch: expected test-host, got %s", machineFacts.Hostname)
	}

	if machineFacts.OS != "Linux" {
		t.Errorf("OS mismatch: expected Linux, got %s", machineFacts.OS)
	}

	if machineFacts.CPU.Cores != 8 {
		t.Errorf("CPU cores mismatch: expected 8, got %d", machineFacts.CPU.Cores)
	}

	if machineFacts.CPU.Frequency != "2.4 GHz" {
		t.Errorf("CPU frequency mismatch: expected 2.4 GHz, got %s", machineFacts.CPU.Frequency)
	}

	if len(machineFacts.IPAddresses) != 2 {
		t.Errorf("IP addresses count mismatch: expected 2, got %d", len(machineFacts.IPAddresses))
	}

	if machineFacts.PrimaryIP != "192.168.1.100" {
		t.Errorf("Primary IP mismatch: expected 192.168.1.100, got %s", machineFacts.PrimaryIP)
	}

	// Test conversion back to FactCollection
	convertedCollection := ConvertMachineFactsToFactCollection(machineFacts)
	if convertedCollection == nil {
		t.Fatal("ConvertMachineFactsToFactCollection returned nil")
	}

	if convertedCollection.Server != machineFacts.Hostname {
		t.Errorf("Server mismatch: expected %s, got %s", machineFacts.Hostname, convertedCollection.Server)
	}

	// Verify some key facts are present
	if fact, exists := convertedCollection.Facts["hostname"]; !exists {
		t.Error("hostname fact not found in converted collection")
	} else if fact.Value != machineFacts.Hostname {
		t.Errorf("hostname fact value mismatch: expected %s, got %v", machineFacts.Hostname, fact.Value)
	}
}
