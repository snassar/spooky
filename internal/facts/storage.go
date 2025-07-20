package facts

import (
	"fmt"
	"io"
	"strings"
	"time"
)

// FactStorage defines the interface for persistent fact storage
type FactStorage interface {
	GetMachineFacts(machineID string) (*MachineFacts, error)
	SetMachineFacts(machineID string, facts *MachineFacts) error
	QueryFacts(query *FactQuery) ([]*MachineFacts, error)
	DeleteFacts(query *FactQuery) (int, error) // Delete facts matching query, returns count
	DeleteMachineFacts(machineID string) error // Delete specific machine facts
	ExportToJSON(w io.Writer) error
	ImportFromJSON(r io.Reader) error
	ExportToJSONWithEncryption(w io.Writer, opts ExportOptions) error
	ImportFromJSONWithDecryption(r io.Reader, identityFile string) error
	Close() error
}

// MachineFacts represents persistent machine facts for storage
type MachineFacts struct {
	MachineID   string            `json:"machine_id"`   // Machine ID as UUID
	MachineName string            `json:"machine_name"` // Human-readable name from HCL
	ActionFile  string            `json:"action_file"`  // Source action file path
	ProjectName string            `json:"project_name"` // Project name (portable across systems)
	ProjectPath string            `json:"project_path"` // Absolute path (for reference)
	Hostname    string            `json:"hostname"`
	IPAddresses []string          `json:"ip_addresses"` // All IP addresses
	PrimaryIP   string            `json:"primary_ip"`   // Primary/primary IP address
	OS          string            `json:"os"`
	OSVersion   string            `json:"os_version"`
	CPU         CPUInfo           `json:"cpu"`
	Memory      MemoryInfo        `json:"memory"`
	SystemID    string            `json:"system_id"` // Actual system ID (/etc/machine-id)
	Tags        map[string]string `json:"tags"`      // Team, environment, etc.
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

// ExportOptions defines options for encrypted export
type ExportOptions struct {
	EncryptSensitive    bool   // Whether to encrypt sensitive fields
	PublicKeyFile       string // Path to age public key file
	SanitizeHostnames   bool   // Whether to sanitize hostnames
	SanitizeMachineID   bool   // Whether to sanitize machine IDs
	SanitizeIPAddresses bool   // Whether to sanitize IP addresses
}

// EncryptedFact represents a fact with encrypted sensitive fields
type EncryptedFact struct {
	MachineID   string            `json:"machine_id"`
	MachineName string            `json:"machine_name"`
	ActionFile  string            `json:"action_file"`
	ProjectName string            `json:"project_name"`
	ProjectPath string            `json:"project_path"`
	Hostname    string            `json:"hostname,omitempty"`
	IPAddress   string            `json:"ip_address,omitempty"`
	OS          string            `json:"os"`
	OSVersion   string            `json:"os_version"`
	CPU         CPUInfo           `json:"cpu"`
	Memory      MemoryInfo        `json:"memory"`
	SystemID    string            `json:"system_id,omitempty"`
	Tags        map[string]string `json:"tags"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`

	// Encrypted fields
	EncryptedEnvironment string `json:"encrypted_environment,omitempty"`
	EncryptedNetwork     string `json:"encrypted_network,omitempty"`
	EncryptedDNS         string `json:"encrypted_dns,omitempty"`
}

// FactQuery defines query parameters for searching facts
type FactQuery struct {
	MachineName   string            // Query by human-readable machine name
	ActionFile    string            // Query by action file
	ProjectName   string            // Query by project name (portable)
	ProjectPath   string            // Query by absolute project path
	Tags          map[string]string // Query by tags
	OS            string            // Query by OS
	Environment   string            // Query by environment tag
	Limit         int               // Limit results
	SearchQuery   string            // Text search query (supports regex)
	SearchField   string            // Field to search in
	UpdatedBefore *time.Time        // Filter by update time
	UpdatedAfter  *time.Time        // Filter by update time
}

// StorageType defines the type of storage backend
type StorageType string

const (
	StorageTypeBadger StorageType = "badger"
	StorageTypeJSON   StorageType = "json"
)

// StorageOptions defines configuration for fact storage
type StorageOptions struct {
	Type StorageType
	Path string
}

// NewFactStorage creates a new fact storage instance
func NewFactStorage(opts StorageOptions) (FactStorage, error) {
	switch opts.Type {
	case StorageTypeJSON:
		return NewJSONFactStorage(opts.Path)
	case StorageTypeBadger, "": // Default to BadgerDB
		return NewBadgerFactStorage(opts.Path)
	default:
		return nil, fmt.Errorf("unsupported storage type: %s", opts.Type)
	}
}

// ConvertFactCollectionToMachineFacts converts a FactCollection to MachineFacts
func ConvertFactCollectionToMachineFacts(machineID string, collection *FactCollection) *MachineFacts {
	facts := &MachineFacts{
		MachineID:   machineID,
		MachineName: collection.Server,
		CreatedAt:   collection.Timestamp,
		UpdatedAt:   time.Now(),
		Tags:        make(map[string]string),
	}

	// Extract facts and convert to storage format
	for key, fact := range collection.Facts {
		if converter, exists := getFactConverter(key); exists {
			converter(facts, fact.Value)
		}
	}

	return facts
}

// factConverter is a function type that converts a fact value to the appropriate field in MachineFacts
type factConverter func(*MachineFacts, interface{})

// getFactConverter returns the appropriate fact converter for the given key
func getFactConverter(key string) (factConverter, bool) {
	converters := map[string]factConverter{
		"hostname":         convertHostname,
		"machine_id":       convertMachineID,
		"os.name":          convertOSName,
		"os.version":       convertOSVersion,
		"cpu.cores":        convertCPUCores,
		"cpu.model":        convertCPUModel,
		"cpu.arch":         convertCPUArch,
		"cpu.frequency":    convertCPUFrequency,
		"memory.total":     convertMemoryTotal,
		"memory.used":      convertMemoryUsed,
		"memory.available": convertMemoryAvailable,
		"network.ips":      convertNetworkIPs,
	}

	converter, exists := converters[key]
	return converter, exists
}

// Individual fact converter functions
func convertHostname(facts *MachineFacts, value interface{}) {
	if str, ok := value.(string); ok {
		facts.Hostname = str
	}
}

func convertMachineID(facts *MachineFacts, value interface{}) {
	if str, ok := value.(string); ok {
		facts.SystemID = str
	}
}

func convertOSName(facts *MachineFacts, value interface{}) {
	if str, ok := value.(string); ok {
		facts.OS = str
	}
}

func convertOSVersion(facts *MachineFacts, value interface{}) {
	if str, ok := value.(string); ok {
		facts.OSVersion = str
	}
}

func convertCPUCores(facts *MachineFacts, value interface{}) {
	if cores, ok := value.(int); ok {
		facts.CPU.Cores = cores
	}
}

func convertCPUModel(facts *MachineFacts, value interface{}) {
	if str, ok := value.(string); ok {
		facts.CPU.Model = str
	}
}

func convertCPUArch(facts *MachineFacts, value interface{}) {
	if str, ok := value.(string); ok {
		facts.CPU.Arch = str
	}
}

func convertCPUFrequency(facts *MachineFacts, value interface{}) {
	if str, ok := value.(string); ok {
		facts.CPU.Frequency = str
	}
}

func convertMemoryTotal(facts *MachineFacts, value interface{}) {
	if total, ok := value.(uint64); ok {
		facts.Memory.Total = total
	}
}

func convertMemoryUsed(facts *MachineFacts, value interface{}) {
	if used, ok := value.(uint64); ok {
		facts.Memory.Used = used
	}
}

func convertMemoryAvailable(facts *MachineFacts, value interface{}) {
	if avail, ok := value.(uint64); ok {
		facts.Memory.Available = avail
	}
}

func convertNetworkIPs(facts *MachineFacts, value interface{}) {
	if ips, ok := value.([]string); ok && len(ips) > 0 {
		facts.IPAddresses = ips
		// Set primary IP to first non-loopback address, or first address if all are loopback
		for _, ip := range ips {
			if !strings.HasPrefix(ip, "127.") && !strings.HasPrefix(ip, "::1") {
				facts.PrimaryIP = ip
				break
			}
		}
		if facts.PrimaryIP == "" && len(ips) > 0 {
			facts.PrimaryIP = ips[0]
		}
	}
}

// matchesQuery checks if facts match the query criteria
func matchesQuery(facts *MachineFacts, query *FactQuery) bool {
	// Check machine name filter
	if query.MachineName != "" && facts.MachineName != query.MachineName {
		return false
	}

	// Check action file filter
	if query.ActionFile != "" && facts.ActionFile != query.ActionFile {
		return false
	}

	// Check project name filter
	if query.ProjectName != "" && facts.ProjectName != query.ProjectName {
		return false
	}

	// Check OS filter
	if query.OS != "" && facts.OS != query.OS {
		return false
	}

	// Check environment filter
	if query.Environment != "" {
		if env, exists := facts.Tags["environment"]; !exists || env != query.Environment {
			return false
		}
	}

	// Check tag filters
	for key, value := range query.Tags {
		if tagValue, exists := facts.Tags[key]; !exists || tagValue != value {
			return false
		}
	}

	// Check time filters
	if query.UpdatedBefore != nil && facts.UpdatedAt.After(*query.UpdatedBefore) {
		return false
	}

	if query.UpdatedAfter != nil && facts.UpdatedAt.Before(*query.UpdatedAfter) {
		return false
	}

	return true
}

// ConvertMachineFactsToFactCollection converts MachineFacts back to FactCollection
func ConvertMachineFactsToFactCollection(facts *MachineFacts) *FactCollection {
	collection := &FactCollection{
		Server:    facts.Hostname,
		Timestamp: facts.CreatedAt,
		Facts:     make(map[string]*Fact),
	}

	// Convert back to individual facts
	if facts.Hostname != "" {
		collection.Facts["hostname"] = &Fact{
			Key:       "hostname",
			Value:     facts.Hostname,
			Source:    "storage",
			Server:    facts.MachineID,
			Timestamp: facts.UpdatedAt,
			TTL:       DefaultTTL,
		}
	}

	if facts.SystemID != "" {
		collection.Facts["machine_id"] = &Fact{
			Key:       "machine_id",
			Value:     facts.SystemID,
			Source:    "storage",
			Server:    facts.MachineID,
			Timestamp: facts.UpdatedAt,
			TTL:       DefaultTTL,
		}
	}

	if facts.OS != "" {
		collection.Facts["os.name"] = &Fact{
			Key:       "os.name",
			Value:     facts.OS,
			Source:    "storage",
			Server:    facts.MachineID,
			Timestamp: facts.UpdatedAt,
			TTL:       DefaultTTL,
		}
	}

	if facts.OSVersion != "" {
		collection.Facts["os.version"] = &Fact{
			Key:       "os.version",
			Value:     facts.OSVersion,
			Source:    "storage",
			Server:    facts.MachineID,
			Timestamp: facts.UpdatedAt,
			TTL:       DefaultTTL,
		}
	}

	if facts.CPU.Cores > 0 {
		collection.Facts["cpu.cores"] = &Fact{
			Key:       "cpu.cores",
			Value:     facts.CPU.Cores,
			Source:    "storage",
			Server:    facts.MachineID,
			Timestamp: facts.UpdatedAt,
			TTL:       DefaultTTL,
		}
	}

	if facts.CPU.Model != "" {
		collection.Facts["cpu.model"] = &Fact{
			Key:       "cpu.model",
			Value:     facts.CPU.Model,
			Source:    "storage",
			Server:    facts.MachineID,
			Timestamp: facts.UpdatedAt,
			TTL:       DefaultTTL,
		}
	}

	if facts.CPU.Arch != "" {
		collection.Facts["cpu.arch"] = &Fact{
			Key:       "cpu.arch",
			Value:     facts.CPU.Arch,
			Source:    "storage",
			Server:    facts.MachineID,
			Timestamp: facts.UpdatedAt,
			TTL:       DefaultTTL,
		}
	}

	if facts.CPU.Frequency != "" {
		collection.Facts["cpu.frequency"] = &Fact{
			Key:       "cpu.frequency",
			Value:     facts.CPU.Frequency,
			Source:    "storage",
			Server:    facts.MachineID,
			Timestamp: facts.UpdatedAt,
			TTL:       DefaultTTL,
		}
	}

	if facts.Memory.Total > 0 {
		collection.Facts["memory.total"] = &Fact{
			Key:       "memory.total",
			Value:     facts.Memory.Total,
			Source:    "storage",
			Server:    facts.MachineID,
			Timestamp: facts.UpdatedAt,
			TTL:       DefaultTTL,
		}
	}

	if facts.Memory.Used > 0 {
		collection.Facts["memory.used"] = &Fact{
			Key:       "memory.used",
			Value:     facts.Memory.Used,
			Source:    "storage",
			Server:    facts.MachineID,
			Timestamp: facts.UpdatedAt,
			TTL:       DefaultTTL,
		}
	}

	if facts.Memory.Available > 0 {
		collection.Facts["memory.available"] = &Fact{
			Key:       "memory.available",
			Value:     facts.Memory.Available,
			Source:    "storage",
			Server:    facts.MachineID,
			Timestamp: facts.UpdatedAt,
			TTL:       DefaultTTL,
		}
	}

	if len(facts.IPAddresses) > 0 {
		collection.Facts["network.ips"] = &Fact{
			Key:       "network.ips",
			Value:     facts.IPAddresses,
			Source:    "storage",
			Server:    facts.MachineID,
			Timestamp: facts.UpdatedAt,
			TTL:       DefaultTTL,
		}
	}

	return collection
}
