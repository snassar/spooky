package facts

import (
	"time"
)

// Fact represents a single fact about a system
type Fact struct {
	Key       string                 `json:"key"`
	Value     interface{}            `json:"value"`
	Source    string                 `json:"source"` // "ssh", "local", "hcl", "opentofu"
	Server    string                 `json:"server"` // server name or "local"
	Timestamp time.Time              `json:"timestamp"`
	TTL       time.Duration          `json:"ttl"`      // time to live, 0 = no expiration
	Metadata  map[string]interface{} `json:"metadata"` // additional context
}

// FactCollection represents a collection of facts for a server
type FactCollection struct {
	Server    string           `json:"server"`
	Timestamp time.Time        `json:"timestamp"`
	Facts     map[string]*Fact `json:"facts"`
}

// FactCollector defines the interface for collecting facts
type FactCollector interface {
	Collect(server string) (*FactCollection, error)
	CollectSpecific(server string, keys []string) (*FactCollection, error)
	GetFact(server, key string) (*Fact, error)
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

// FactSource defines where facts come from
type FactSource string

const (
	SourceSSH      FactSource = "ssh"
	SourceLocal    FactSource = "local"
	SourceHCL      FactSource = "hcl"
	SourceOpenTofu FactSource = "opentofu"
)

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
