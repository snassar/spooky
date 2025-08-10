package facts

import (
	"time"
)

// Fact represents a single piece of information about a system
type Fact struct {
	Key       string                 `json:"key"`
	Value     interface{}            `json:"value"`
	Source    string                 `json:"source"` // "ssh", "local", "hcl", "opentofu"
	Server    string                 `json:"server"` // server name or "local"
	Timestamp time.Time              `json:"timestamp"`
	TTL       time.Duration          `json:"ttl"`      // time to live, 0 = no expiration
	Metadata  map[string]interface{} `json:"metadata"` // additional context

	// Individual fact encryption
	Encrypted     bool   `json:"encrypted,omitempty"`      // Whether this fact is encrypted
	EncryptedData string `json:"encrypted_data,omitempty"` // Encrypted value if encrypted
}

// FactCollection represents a collection of facts for a server
type FactCollection struct {
	Server      string                            `json:"server"`
	Timestamp   time.Time                         `json:"timestamp"`
	Facts       map[string]*Fact                  `json:"facts"`        // System facts (flat structure)
	CustomFacts map[string]map[string]interface{} `json:"custom_facts"` // Custom facts (hierarchical: filename -> key -> value)

	// Encryption metadata
	EncryptionMetadata *EncryptionMetadata `json:"encryption_metadata,omitempty"`
	EncryptedData      string              `json:"encrypted_data,omitempty"`
}

// EncryptionMetadata contains information about fact encryption
type EncryptionMetadata struct {
	EncryptedAt       string   `json:"encrypted_at"`
	EncryptionVersion string   `json:"encryption_version"`
	Recipients        []string `json:"recipients"`
}

// Clone creates a deep copy of the fact collection
func (fc *FactCollection) Clone() *FactCollection {
	if fc == nil {
		return nil
	}

	clone := &FactCollection{
		Server:             fc.Server,
		Timestamp:          fc.Timestamp,
		Facts:              make(map[string]*Fact),
		CustomFacts:        make(map[string]map[string]interface{}),
		EncryptionMetadata: fc.EncryptionMetadata,
		EncryptedData:      fc.EncryptedData,
	}

	// Clone facts
	for key, fact := range fc.Facts {
		clone.Facts[key] = &Fact{
			Key:           fact.Key,
			Value:         fact.Value,
			Source:        fact.Source,
			Server:        fact.Server,
			Timestamp:     fact.Timestamp,
			TTL:           fact.TTL,
			Metadata:      fact.Metadata,
			Encrypted:     fact.Encrypted,
			EncryptedData: fact.EncryptedData,
		}
	}

	// Clone custom facts
	for filename, facts := range fc.CustomFacts {
		clone.CustomFacts[filename] = make(map[string]interface{})
		for key, value := range facts {
			clone.CustomFacts[filename][key] = value
		}
	}

	return clone
}

// GetCustomFact retrieves a custom fact by filename and key
func (fc *FactCollection) GetCustomFact(filename, key string) (interface{}, bool) {
	if fc.CustomFacts == nil {
		return nil, false
	}
	if facts, exists := fc.CustomFacts[filename]; exists {
		if value, exists := facts[key]; exists {
			return value, true
		}
	}
	return nil, false
}

// SetCustomFact sets a custom fact by filename and key
func (fc *FactCollection) SetCustomFact(filename, key string, value interface{}) {
	if fc.CustomFacts == nil {
		fc.CustomFacts = make(map[string]map[string]interface{})
	}
	if fc.CustomFacts[filename] == nil {
		fc.CustomFacts[filename] = make(map[string]interface{})
	}
	fc.CustomFacts[filename][key] = value
}

// GetCustomFactsByFile retrieves all custom facts for a specific file
func (fc *FactCollection) GetCustomFactsByFile(filename string) (map[string]interface{}, bool) {
	if fc.CustomFacts == nil {
		return nil, false
	}
	if facts, exists := fc.CustomFacts[filename]; exists {
		return facts, true
	}
	return nil, false
}

// GetAllCustomFacts retrieves all custom facts
func (fc *FactCollection) GetAllCustomFacts() map[string]map[string]interface{} {
	if fc.CustomFacts == nil {
		return make(map[string]map[string]interface{})
	}
	return fc.CustomFacts
}

// HasCustomFacts checks if the collection has any custom facts
func (fc *FactCollection) HasCustomFacts() bool {
	return len(fc.CustomFacts) > 0
}

// GetValue returns the fact value, handling encrypted values
func (f *Fact) GetValue() interface{} {
	if f.Encrypted {
		// Return encrypted data if fact is encrypted
		return f.EncryptedData
	}
	return f.Value
}

// SetEncryptedValue sets the fact as encrypted with the provided encrypted data
func (f *Fact) SetEncryptedValue(encryptedData string) {
	f.Encrypted = true
	f.EncryptedData = encryptedData
	f.Value = nil // Clear plain value when encrypted
}

// SetPlainValue sets the fact as plain text with the provided value
func (f *Fact) SetPlainValue(value interface{}) {
	f.Encrypted = false
	f.EncryptedData = ""
	f.Value = value
}

// IsEncrypted checks if the fact is encrypted
func (f *Fact) IsEncrypted() bool {
	return f.Encrypted
}

// FactCollector defines the interface for collecting facts
type FactCollector interface {
	Collect(server string) (*FactCollection, error)
	CollectSpecific(server string, keys []string) (*FactCollection, error)
	GetFact(server, key string) (*Fact, error)
}

// FactKey represents common fact keys
const (
	FactMachineID   = "machine_id"
	FactHostname    = "hostname"
	FactFQDN        = "fqdn"
	FactOSName      = "os.name"
	FactOSVersion   = "os.version"
	FactOSDistro    = "os.distribution"
	FactOSArch      = "os.architecture"
	FactOSKernel    = "os.kernel"
	FactCPUCores    = "cpu.cores"
	FactCPUModel    = "cpu.model"
	FactCPUArch     = "cpu.arch"
	FactCPUFreq     = "cpu.frequency"
	FactMemoryTotal = "memory.total"
	FactMemoryUsed  = "memory.used"
	FactMemoryAvail = "memory.available"
	FactDiskTotal   = "disk.total"
	FactDiskUsed    = "disk.used"
	FactDiskAvail   = "disk.available"
	FactNetworkIPs  = "network.ips"
	FactNetworkMACs = "network.macs"
	FactDNS         = "network.dns"
	FactEnvironment = "environment"
)

// DefaultTTL is the default time-to-live for facts
const DefaultTTL = 1 * time.Hour
