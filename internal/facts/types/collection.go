package types

// FactSource defines where facts come from
type FactSource string

const (
	SourceSSH      FactSource = "ssh"
	SourceLocal    FactSource = "local"
	SourceHCL      FactSource = "hcl"
	SourceOpenTofu FactSource = "opentofu"
	SourceCustom   FactSource = "custom"
)

// MergePolicy defines how facts should be merged
type MergePolicy string

const (
	MergePolicyReplace MergePolicy = "replace" // Replace existing facts
	MergePolicyMerge   MergePolicy = "merge"   // Merge with existing facts
	MergePolicyAppend  MergePolicy = "append"  // Append to existing facts
	MergePolicySkip    MergePolicy = "skip"    // Skip if facts exist
)

// MergeMode defines the mode for merging facts during import
type MergeMode string

const (
	MergeModeReplace MergeMode = "replace" // Replace existing facts
	MergeModeMerge   MergeMode = "merge"   // Merge with existing facts
	MergeModeAppend  MergeMode = "append"  // Append to existing facts
	MergeModeSelect  MergeMode = "select"  // Only import selected facts
)

// CustomFacts represents custom facts structure
type CustomFacts struct {
	Custom    map[string]interface{} `json:"custom"`
	Overrides map[string]interface{} `json:"overrides"`
	Source    string                 `json:"source,omitempty"`
}

// ImportOptions defines options for fact import operations
type ImportOptions struct {
	Source      string    `json:"source"`
	Path        string    `json:"path"`
	MergeMode   MergeMode `json:"merge_mode"`
	SelectFacts []string  `json:"select_facts"`
	Override    bool      `json:"override"`
	Validate    bool      `json:"validate"`
	DryRun      bool      `json:"dry_run"`
	Server      string    `json:"server"`
	Encrypt     bool      `json:"encrypt"` // Encrypt facts during import
	Decrypt     bool      `json:"decrypt"` // Decrypt facts during import
}

// SystemFacts contains common system information
type SystemFacts struct {
	MachineID   string            `json:"machine_id"`
	Hostname    string            `json:"hostname"`
	FQDN        string            `json:"fqdn"`
	OS          OSInfo            `json:"os"`
	Hardware    HardwareInfo      `json:"hardware"`
	Network     NetworkInfo       `json:"network"`
	Environment map[string]string `json:"environment"`
}

// OSInfo contains operating system information
type OSInfo struct {
	Name         string `json:"name"`
	Version      string `json:"version"`
	Distribution string `json:"distribution"`
	Architecture string `json:"architecture"`
	Kernel       string `json:"kernel"`
}

// HardwareInfo contains hardware information
type HardwareInfo struct {
	CPU     CPUInfo     `json:"cpu"`
	Memory  MemoryInfo  `json:"memory"`
	Storage StorageInfo `json:"storage"`
}

// CPUInfo contains CPU information
type CPUInfo struct {
	Cores     int    `json:"cores"`
	Model     string `json:"model"`
	Arch      string `json:"arch"`
	Frequency string `json:"frequency"`
}

// MemoryInfo contains memory information
type MemoryInfo struct {
	Total     uint64 `json:"total"`     // in bytes
	Available uint64 `json:"available"` // in bytes
	Used      uint64 `json:"used"`      // in bytes
}

// StorageInfo contains storage information
type StorageInfo struct {
	Disks []DiskInfo `json:"disks"`
}

// DiskInfo contains disk information
type DiskInfo struct {
	Device     string `json:"device"`
	MountPoint string `json:"mount_point"`
	Total      uint64 `json:"total"`     // in bytes
	Used       uint64 `json:"used"`      // in bytes
	Available  uint64 `json:"available"` // in bytes
	Filesystem string `json:"filesystem"`
}

// NetworkInfo contains network information
type NetworkInfo struct {
	Interfaces []InterfaceInfo `json:"interfaces"`
	DNS        DNSInfo         `json:"dns"`
}

// InterfaceInfo contains network interface information
type InterfaceInfo struct {
	Name      string   `json:"name"`
	Addresses []string `json:"addresses"`
	MAC       string   `json:"mac"`
	MTU       int      `json:"mtu"`
	State     string   `json:"state"`
}

// DNSInfo contains DNS configuration
type DNSInfo struct {
	Nameservers []string `json:"nameservers"`
	Search      []string `json:"search"`
}
