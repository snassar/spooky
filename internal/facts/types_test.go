package facts

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFactStruct(t *testing.T) {
	// Test Fact struct creation and access
	now := time.Now()
	fact := &Fact{
		Key:       "os.name",
		Value:     "linux",
		Source:    "ssh",
		Server:    "test-server",
		Timestamp: now,
		TTL:       30 * time.Minute,
		Metadata: map[string]interface{}{
			"test": "value",
		},
	}

	assert.Equal(t, "os.name", fact.Key)
	assert.Equal(t, "linux", fact.Value)
	assert.Equal(t, "ssh", fact.Source)
	assert.Equal(t, "test-server", fact.Server)
	assert.Equal(t, now, fact.Timestamp)
	assert.Equal(t, 30*time.Minute, fact.TTL)
	assert.Equal(t, "value", fact.Metadata["test"])
}

func TestFactCollectionStruct(t *testing.T) {
	// Test FactCollection struct creation and access
	now := time.Now()
	collection := &FactCollection{
		Server:    "test-server",
		Timestamp: now,
		Facts: map[string]*Fact{
			"os.name": {
				Key:       "os.name",
				Value:     "linux",
				Source:    "ssh",
				Server:    "test-server",
				Timestamp: now,
			},
		},
	}

	assert.Equal(t, "test-server", collection.Server)
	assert.Equal(t, now, collection.Timestamp)
	assert.Len(t, collection.Facts, 1)
	assert.Equal(t, "linux", collection.Facts["os.name"].Value)
}

func TestFactCollectionClone(t *testing.T) {
	// Test FactCollection cloning
	now := time.Now()
	original := &FactCollection{
		Server:    "test-server",
		Timestamp: now,
		Facts: map[string]*Fact{
			"os.name": {
				Key:       "os.name",
				Value:     "linux",
				Source:    "ssh",
				Server:    "test-server",
				Timestamp: now,
				Metadata: map[string]interface{}{
					"test": "value",
				},
			},
			"hardware.cpu": {
				Key:       "hardware.cpu",
				Value:     "Intel i7",
				Source:    "local",
				Server:    "test-server",
				Timestamp: now,
			},
		},
	}

	cloned := original.Clone()

	assert.NotNil(t, cloned)
	assert.Equal(t, original.Server, cloned.Server)
	assert.Equal(t, original.Timestamp, cloned.Timestamp)
	assert.NotEqual(t, original.Facts, cloned.Facts) // Should be different maps
	assert.Equal(t, len(original.Facts), len(cloned.Facts))

	// Test that modifying cloned doesn't affect original
	cloned.Facts["os.name"].Value = "windows"
	assert.Equal(t, "linux", original.Facts["os.name"].Value)
	assert.Equal(t, "windows", cloned.Facts["os.name"].Value)

	// Test metadata cloning
	cloned.Facts["os.name"].Metadata["test"] = "new-value"
	assert.Equal(t, "value", original.Facts["os.name"].Metadata["test"])
	assert.Equal(t, "new-value", cloned.Facts["os.name"].Metadata["test"])
}

func TestFactCollectionCloneNil(t *testing.T) {
	// Test cloning nil collection
	var collection *FactCollection
	cloned := collection.Clone()
	assert.Nil(t, cloned)
}

func TestSystemFactsStruct(t *testing.T) {
	// Test SystemFacts struct creation and access
	systemFacts := &SystemFacts{
		MachineID: "test-machine-id",
		Hostname:  "test-host",
		FQDN:      "test-host.example.com",
		OS: OSInfo{
			Name:         "linux",
			Version:      "20.04",
			Distribution: "ubuntu",
			Architecture: "x86_64",
			Kernel:       "5.4.0",
		},
		Hardware: HardwareInfo{
			CPU: CPUInfo{
				Cores:     4,
				Model:     "Intel i7",
				Arch:      "x86_64",
				Frequency: "3.7GHz",
			},
			Memory: MemoryInfo{
				Total:     16 * 1024 * 1024 * 1024, // 16GB
				Used:      8 * 1024 * 1024 * 1024,  // 8GB
				Available: 8 * 1024 * 1024 * 1024,  // 8GB
			},
			Storage: StorageInfo{
				Disks: []DiskInfo{
					{
						Device:     "/dev/sda",
						MountPoint: "/",
						Total:      500 * 1024 * 1024 * 1024, // 500GB
						Used:       200 * 1024 * 1024 * 1024, // 200GB
						Available:  300 * 1024 * 1024 * 1024, // 300GB
						Filesystem: "ext4",
					},
				},
			},
		},
		Network: NetworkInfo{
			Interfaces: []InterfaceInfo{
				{
					Name:      "eth0",
					Addresses: []string{"192.168.1.100"},
					MAC:       "00:11:22:33:44:55",
					MTU:       1500,
					State:     "up",
				},
			},
			DNS: DNSInfo{
				Nameservers: []string{"8.8.8.8", "8.8.4.4"},
				Search:      []string{"example.com"},
			},
		},
		Environment: map[string]string{
			"PATH": "/usr/local/bin:/usr/bin:/bin",
			"HOME": "/home/user",
		},
	}

	assert.Equal(t, "test-machine-id", systemFacts.MachineID)
	assert.Equal(t, "test-host", systemFacts.Hostname)
	assert.Equal(t, "test-host.example.com", systemFacts.FQDN)
	assert.Equal(t, "linux", systemFacts.OS.Name)
	assert.Equal(t, "20.04", systemFacts.OS.Version)
	assert.Equal(t, 4, systemFacts.Hardware.CPU.Cores)
	assert.Equal(t, "Intel i7", systemFacts.Hardware.CPU.Model)
	assert.Equal(t, uint64(16*1024*1024*1024), systemFacts.Hardware.Memory.Total)
	assert.Len(t, systemFacts.Hardware.Storage.Disks, 1)
	assert.Equal(t, "/dev/sda", systemFacts.Hardware.Storage.Disks[0].Device)
	assert.Len(t, systemFacts.Network.Interfaces, 1)
	assert.Equal(t, "eth0", systemFacts.Network.Interfaces[0].Name)
	assert.Len(t, systemFacts.Network.DNS.Nameservers, 2)
	assert.Equal(t, "8.8.8.8", systemFacts.Network.DNS.Nameservers[0])
	assert.Len(t, systemFacts.Environment, 2)
	assert.Equal(t, "/usr/local/bin:/usr/bin:/bin", systemFacts.Environment["PATH"])
}

func TestOSInfoStruct(t *testing.T) {
	// Test OSInfo struct
	osInfo := OSInfo{
		Name:         "linux",
		Version:      "20.04",
		Distribution: "ubuntu",
		Architecture: "x86_64",
		Kernel:       "5.4.0",
	}

	assert.Equal(t, "linux", osInfo.Name)
	assert.Equal(t, "20.04", osInfo.Version)
	assert.Equal(t, "ubuntu", osInfo.Distribution)
	assert.Equal(t, "x86_64", osInfo.Architecture)
	assert.Equal(t, "5.4.0", osInfo.Kernel)
}

func TestHardwareInfoStruct(t *testing.T) {
	// Test HardwareInfo struct
	hardware := HardwareInfo{
		CPU: CPUInfo{
			Cores:     4,
			Model:     "Intel i7",
			Arch:      "x86_64",
			Frequency: "3.7GHz",
		},
		Memory: MemoryInfo{
			Total:     16 * 1024 * 1024 * 1024,
			Used:      8 * 1024 * 1024 * 1024,
			Available: 8 * 1024 * 1024 * 1024,
		},
		Storage: StorageInfo{
			Disks: []DiskInfo{
				{
					Device:     "/dev/sda",
					MountPoint: "/",
					Total:      500 * 1024 * 1024 * 1024,
					Used:       200 * 1024 * 1024 * 1024,
					Available:  300 * 1024 * 1024 * 1024,
					Filesystem: "ext4",
				},
			},
		},
	}

	assert.Equal(t, 4, hardware.CPU.Cores)
	assert.Equal(t, "Intel i7", hardware.CPU.Model)
	assert.Equal(t, uint64(16*1024*1024*1024), hardware.Memory.Total)
	assert.Len(t, hardware.Storage.Disks, 1)
	assert.Equal(t, "/dev/sda", hardware.Storage.Disks[0].Device)
}

func TestNetworkInfoStruct(t *testing.T) {
	// Test NetworkInfo struct
	network := NetworkInfo{
		Interfaces: []InterfaceInfo{
			{
				Name:      "eth0",
				Addresses: []string{"192.168.1.100", "fe80::1234:5678:9abc:def0"},
				MAC:       "00:11:22:33:44:55",
				MTU:       1500,
				State:     "up",
			},
			{
				Name:      "lo",
				Addresses: []string{"127.0.0.1", "::1"},
				MAC:       "",
				MTU:       65536,
				State:     "up",
			},
		},
		DNS: DNSInfo{
			Nameservers: []string{"8.8.8.8", "8.8.4.4", "1.1.1.1"},
			Search:      []string{"example.com", "local"},
		},
	}

	assert.Len(t, network.Interfaces, 2)
	assert.Equal(t, "eth0", network.Interfaces[0].Name)
	assert.Len(t, network.Interfaces[0].Addresses, 2)
	assert.Equal(t, "192.168.1.100", network.Interfaces[0].Addresses[0])
	assert.Equal(t, "00:11:22:33:44:55", network.Interfaces[0].MAC)
	assert.Equal(t, 1500, network.Interfaces[0].MTU)
	assert.Equal(t, "up", network.Interfaces[0].State)

	assert.Len(t, network.DNS.Nameservers, 3)
	assert.Equal(t, "8.8.8.8", network.DNS.Nameservers[0])
	assert.Len(t, network.DNS.Search, 2)
	assert.Equal(t, "example.com", network.DNS.Search[0])
}

func TestFactSourceConstants(t *testing.T) {
	// Test FactSource constants
	assert.Equal(t, FactSource("ssh"), SourceSSH)
	assert.Equal(t, FactSource("local"), SourceLocal)
	assert.Equal(t, FactSource("hcl"), SourceHCL)
	assert.Equal(t, FactSource("opentofu"), SourceOpenTofu)
	assert.Equal(t, FactSource("custom"), SourceCustom)
}

func TestMergePolicyConstants(t *testing.T) {
	// Test MergePolicy constants
	assert.Equal(t, MergePolicy("replace"), MergePolicyReplace)
	assert.Equal(t, MergePolicy("append"), MergePolicyAppend)
	assert.Equal(t, MergePolicy("merge"), MergePolicyMerge)
}

func TestCustomFactsStruct(t *testing.T) {
	// Test CustomFacts struct
	customFacts := &CustomFacts{
		Custom: map[string]interface{}{
			"os": map[string]interface{}{
				"name":    "linux",
				"version": "20.04",
			},
			"hardware": map[string]interface{}{
				"cpu":    "Intel i7",
				"memory": "16GB",
			},
		},
		Overrides: map[string]interface{}{
			"network.hostname": "custom-hostname",
		},
		Source: "test-source",
	}

	assert.Len(t, customFacts.Custom, 2)
	assert.Equal(t, "linux", customFacts.Custom["os"].(map[string]interface{})["name"])
	assert.Equal(t, "Intel i7", customFacts.Custom["hardware"].(map[string]interface{})["cpu"])
	assert.Len(t, customFacts.Overrides, 1)
	assert.Equal(t, "custom-hostname", customFacts.Overrides["network.hostname"])
	assert.Equal(t, "test-source", customFacts.Source)
}

func TestImportOptionsStruct(t *testing.T) {
	// Test ImportOptions struct
	options := &ImportOptions{
		Source:      "file:///path/to/facts.json",
		Path:        "/path/to/facts.json",
		MergeMode:   MergeModeAppend,
		SelectFacts: []string{"os.*", "hardware.*"},
		Override:    true,
		Validate:    true,
		DryRun:      false,
		Server:      "test-server",
	}

	assert.Equal(t, "file:///path/to/facts.json", options.Source)
	assert.Equal(t, "/path/to/facts.json", options.Path)
	assert.Equal(t, MergeModeAppend, options.MergeMode)
	assert.Len(t, options.SelectFacts, 2)
	assert.Equal(t, "os.*", options.SelectFacts[0])
	assert.Equal(t, "hardware.*", options.SelectFacts[1])
	assert.True(t, options.Override)
	assert.True(t, options.Validate)
	assert.False(t, options.DryRun)
	assert.Equal(t, "test-server", options.Server)
}

func TestMergeModeConstants(t *testing.T) {
	// Test MergeMode constants
	assert.Equal(t, MergeMode("append"), MergeModeAppend)
	assert.Equal(t, MergeMode("replace"), MergeModeReplace)
	assert.Equal(t, MergeMode("merge"), MergeModeMerge)
}

func TestFactWithTestDataFromExamples(t *testing.T) {
	// Test with data patterns from examples/testing
	examplesDir := "../../examples/testing"

	// Test with valid project data patterns
	validProjectPath := examplesDir + "/test-valid-project"
	if _, err := os.Stat(validProjectPath); os.IsNotExist(err) {
		t.Skip("Test data directory not found, skipping test")
	}

	// Create facts similar to those in test data
	fact := &Fact{
		Key:       "os.name",
		Value:     "linux",
		Source:    "custom",
		Server:    "example-server",
		Timestamp: time.Now(),
		Metadata: map[string]interface{}{
			"project": "test-valid-project",
			"version": "1.0.0",
		},
	}

	assert.Equal(t, "os.name", fact.Key)
	assert.Equal(t, "linux", fact.Value)
	assert.Equal(t, "custom", fact.Source)
	assert.Equal(t, "example-server", fact.Server)
	assert.Equal(t, "test-valid-project", fact.Metadata["project"])
	assert.Equal(t, "1.0.0", fact.Metadata["version"])

	// Test collection with multiple facts
	collection := &FactCollection{
		Server:    "example-server",
		Timestamp: time.Now(),
		Facts: map[string]*Fact{
			"os.name": fact,
			"hardware.cpu": {
				Key:       "hardware.cpu",
				Value:     "Intel(R) Core(TM) i7-8700K CPU @ 3.70GHz",
				Source:    "custom",
				Server:    "example-server",
				Timestamp: time.Now(),
			},
			"network.hostname": {
				Key:       "network.hostname",
				Value:     "example-server",
				Source:    "custom",
				Server:    "example-server",
				Timestamp: time.Now(),
			},
		},
	}

	assert.Equal(t, "example-server", collection.Server)
	assert.Len(t, collection.Facts, 3)
	assert.Equal(t, "linux", collection.Facts["os.name"].Value)
	assert.Equal(t, "Intel(R) Core(TM) i7-8700K CPU @ 3.70GHz", collection.Facts["hardware.cpu"].Value)
	assert.Equal(t, "example-server", collection.Facts["network.hostname"].Value)
}
