package machines

import "time"

// IndexManagerConfig holds configuration for the index manager
type IndexManagerConfig struct {
	CacheTTL             time.Duration
	OptimizationInterval time.Duration
	MaxIndexSize         int
	EnableMetrics        bool
	EnableOptimization   bool
}

// DefaultIndexManagerConfig returns default configuration
func DefaultIndexManagerConfig() *IndexManagerConfig {
	return &IndexManagerConfig{
		CacheTTL:             5 * time.Minute,
		OptimizationInterval: 10 * time.Minute,
		MaxIndexSize:         10000,
		EnableMetrics:        true,
		EnableOptimization:   true,
	}
}

// ConnectivityTestOptions represents options for connectivity testing
type ConnectivityTestOptions struct {
	Timeout    time.Duration
	Parallel   int
	Phases     []ConnectivityTestPhase
	SSHCommand string
	DNSTimeout time.Duration
}

// DefaultConnectivityTestOptions returns default connectivity test options
func DefaultConnectivityTestOptions() *ConnectivityTestOptions {
	return &ConnectivityTestOptions{
		Timeout:    30 * time.Second,
		Parallel:   10,
		Phases:     []ConnectivityTestPhase{PhaseDNS, PhaseSSH},
		SSHCommand: "echo 'SSH connectivity test'",
		DNSTimeout: 5 * time.Second,
	}
}

// ExportOptions contains export configuration options
type ExportOptions struct {
	Format          ExportFormat
	OutputFile      string
	Filter          *MachineFilter
	SortBy          string
	IncludeMetadata bool
}
