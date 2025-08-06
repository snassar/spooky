package types

import (
	"time"

	"spooky/internal/config/types"
)

// IndexType represents different types of machine indexes
type IndexType string

const (
	IndexTypeName     IndexType = "name"
	IndexTypeHost     IndexType = "host"
	IndexTypeTag      IndexType = "tag"
	IndexTypeGroup    IndexType = "group"
	IndexTypeNetwork  IndexType = "network"
	IndexTypeUser     IndexType = "user"
	IndexTypePort     IndexType = "port"
	IndexTypeMetadata IndexType = "metadata"
)

// ConnectivityTestPhase represents the different phases of connectivity testing
type ConnectivityTestPhase string

const (
	PhaseDNS ConnectivityTestPhase = "dns"
	PhaseSSH ConnectivityTestPhase = "ssh"
)

// ExportFormat represents supported export formats
type ExportFormat string

const (
	ExportFormatHCL  ExportFormat = "hcl"
	ExportFormatJSON ExportFormat = "json"
)

// IndexManagerState tracks the current state of the index manager
type IndexManagerState struct {
	LastBuilt     time.Time
	LastOptimized time.Time
	MachineCount  int
	IndexCount    int
	CacheHitRate  float64
	MemoryUsage   int64
	IsOptimizing  bool
	LastError     error
}

// IndexSyncData represents synchronization data for indexes
type IndexSyncData struct {
	SyncID    string
	Timestamp time.Time
	Machines  []types.Machine
	Checksum  string
}

// IndexMetrics tracks performance metrics for the indexing system
type IndexMetrics struct {
	BuildTime      time.Duration
	LookupTime     time.Duration
	MemoryUsage    int64
	HitRate        float64
	MachineCount   int
	IndexCount     int
	LastUpdated    time.Time
	IndexTypeStats map[IndexType]*IndexTypeStats
}

// IndexTypeStats tracks statistics for each index type
type IndexTypeStats struct {
	BuildTime   time.Duration
	LookupTime  time.Duration
	HitCount    int64
	MissCount   int64
	MemoryUsage int64
	EntryCount  int
}

// IndexPerformanceStats contains performance statistics
type IndexPerformanceStats struct {
	TotalBuildTime   time.Duration
	TotalLookupTime  time.Duration
	TotalMemoryUsage int64
	AverageHitRate   float64
	IndexTypeStats   map[IndexType]*IndexTypePerformanceStats
}

// IndexTypePerformanceStats contains performance stats for each index type
type IndexTypePerformanceStats struct {
	BuildTime   time.Duration
	LookupTime  time.Duration
	MemoryUsage int64
	HitRate     float64
	EntryCount  int
}

// ConnectivityTestResult represents the result of a connectivity test
type ConnectivityTestResult struct {
	Machine   *types.Machine         `json:"machine"`
	Phase     ConnectivityTestPhase  `json:"phase"`
	Success   bool                   `json:"success"`
	Error     string                 `json:"error,omitempty"`
	Duration  time.Duration          `json:"duration"`
	Timestamp time.Time              `json:"timestamp"`
	Details   map[string]interface{} `json:"details,omitempty"`
}

// ConnectivityTestSummary contains summary information for connectivity tests
type ConnectivityTestSummary struct {
	TotalTests      int
	SuccessfulTests int
	FailedTests     int
	TotalDuration   time.Duration
	PhaseResults    map[ConnectivityTestPhase]*PhaseSummary
}

// PhaseSummary contains summary for a specific test phase
type PhaseSummary struct {
	Phase           ConnectivityTestPhase
	TotalTests      int
	SuccessfulTests int
	FailedTests     int
	AverageDuration time.Duration
}

// MachineFilter contains filtering criteria for machine export
type MachineFilter struct {
	MachineNames []string
	Tags         map[string]string
	Groups       []string
	Metadata     map[string]string
	HostPattern  string
	User         string
}

// ExportResult contains the result of an export operation
type ExportResult struct {
	Machines []types.Machine `json:"machines"`
	Stats    ExportStats     `json:"stats"`
	Options  ExportOptions   `json:"options"`
}

// ExportStats contains export statistics
type ExportStats struct {
	TotalMachines    int   `json:"total_machines"`
	ExportedMachines int   `json:"exported_machines"`
	FilteredMachines int   `json:"filtered_machines"`
	ExportTime       int64 `json:"export_time_ms"`
}

// EncryptionResult contains the result of encryption operations
type EncryptionResult struct {
	TotalMachines       int               `json:"total_machines"`
	ProcessedMachines   int               `json:"processed_machines"`
	FilteredMachines    int               `json:"filtered_machines"`
	EncryptedMachines   int               `json:"encrypted_machines,omitempty"`
	DecryptedMachines   int               `json:"decrypted_machines,omitempty"`
	ReencryptedMachines int               `json:"reencrypted_machines,omitempty"`
	FailedMachines      int               `json:"failed_machines"`
	Errors              []EncryptionError `json:"errors,omitempty"`
}

// EncryptionError contains error information for encryption operations
type EncryptionError struct {
	MachineName string `json:"machine_name"`
	Field       string `json:"field"`
	Message     string `json:"message"`
}

// MachineEncryptionStatus contains encryption status for a machine
type MachineEncryptionStatus struct {
	MachineName string                           `json:"machine_name"`
	Fields      map[string]FieldEncryptionStatus `json:"fields"`
}

// FieldEncryptionStatus contains encryption status for a specific field
type FieldEncryptionStatus struct {
	FieldName   string `json:"field_name"`
	IsEmpty     bool   `json:"is_empty"`
	IsEncrypted bool   `json:"is_encrypted"`
}

// MachineIndex provides O(1) lookups for machine data
type MachineIndex struct {
	// Primary indexes
	NameIndex map[string]*types.Machine
	HostIndex map[string]*types.Machine
	UserIndex map[string][]*types.Machine
	PortIndex map[int][]*types.Machine

	// Tag-based indexes
	TagIndex      map[string][]*types.Machine
	TagValueIndex map[string][]*types.Machine

	// Group-based indexes
	GroupIndex map[string][]*types.Machine

	// Network-based indexes
	NetworkIndex map[string][]*types.Machine
	SubnetIndex  map[string][]*types.Machine

	// Metadata indexes
	MetadataIndex map[string][]*types.Machine

	// Reverse indexes for efficient lookups
	MachineTags    map[*types.Machine]map[string]string
	MachineGroups  map[*types.Machine][]string
	MachineNetwork map[*types.Machine]string

	// Performance tracking
	Metrics *IndexMetrics
}
