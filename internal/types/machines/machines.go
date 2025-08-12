// Package machines provides machine inventory types for the spooky codebase.
// This package defines the data structures for machine inventory management.
package machines

import (
	"time"

	spookytypescommon "spooky/internal/types/common"
)

// Machine represents a single machine in the inventory
type Machine struct {
	spookytypescommon.CompleteEntity

	// Basic identification
	Hostname string `json:"hostname" hcl:"hostname"`
	Host     string `json:"host" hcl:"host"`
	Port     int    `json:"port,omitempty" hcl:"port,optional" default:"22"`

	// Authentication
	User       string `json:"user" hcl:"user"`
	Password   string `json:"password,omitempty" hcl:"password,optional" sensitive:"true"`
	KeyFile    string `json:"key_file,omitempty" hcl:"key_file,optional"`
	Passphrase string `json:"passphrase,omitempty" hcl:"passphrase,optional" sensitive:"true"`

	// Organization
	Tags    map[string]string `json:"tags,omitempty" hcl:"tags,optional"`
	Groups  []string          `json:"groups,omitempty" hcl:"groups,optional"`
	Roles   []string          `json:"roles,omitempty" hcl:"roles,optional"`
	Classes []string          `json:"classes,omitempty" hcl:"classes,optional"`

	// SSH connection configuration
	ConnectionTimeout int `json:"connection_timeout,omitempty" hcl:"connection_timeout,optional" default:"30"`
	CommandTimeout    int `json:"command_timeout,omitempty" hcl:"command_timeout,optional" default:"300"`
	MaxConnections    int `json:"max_connections,omitempty" hcl:"max_connections,optional" default:"10"`
	RetryAttempts     int `json:"retry_attempts,omitempty" hcl:"retry_attempts,optional" default:"3"`
	RetryDelay        int `json:"retry_delay,omitempty" hcl:"retry_delay,optional" default:"5"`

	// Resource specifications
	Resources *MachineResources `json:"resources,omitempty" hcl:"resources,optional"`

	// Machine metadata
	MachineMetadata *MachineMetadata `json:"metadata,omitempty" hcl:"metadata,optional"`

	// Connectivity status
	Connectivity *MachineConnectivity `json:"connectivity,omitempty" hcl:"connectivity,optional"`
}

// MachineResources represents machine resource specifications
type MachineResources struct {
	CPUCores     int    `json:"cpu_cores,omitempty" hcl:"cpu_cores,optional"`
	MemoryGB     int    `json:"memory_gb,omitempty" hcl:"memory_gb,optional"`
	DiskGB       int    `json:"disk_gb,omitempty" hcl:"disk_gb,optional"`
	NetworkSpeed string `json:"network_speed,omitempty" hcl:"network_speed,optional"`
}

// MachineMetadata represents machine-specific metadata
type MachineMetadata struct {
	Environment  string            `json:"environment,omitempty" hcl:"environment,optional"`
	Datacenter   string            `json:"datacenter,omitempty" hcl:"datacenter,optional"`
	Rack         string            `json:"rack,omitempty" hcl:"rack,optional"`
	Location     string            `json:"location,omitempty" hcl:"location,optional"`
	Owner        string            `json:"owner,omitempty" hcl:"owner,optional"`
	Department   string            `json:"department,omitempty" hcl:"department,optional"`
	CostCenter   string            `json:"cost_center,omitempty" hcl:"cost_center,optional"`
	CustomFields map[string]string `json:"custom_fields,omitempty" hcl:"custom_fields,optional"`
}

// MachineConnectivity represents machine connectivity status
type MachineConnectivity struct {
	LastPing       time.Time `json:"last_ping,omitempty" hcl:"last_ping,optional"`
	PingLatency    int       `json:"ping_latency,omitempty" hcl:"ping_latency,optional"` // milliseconds
	SSHReachable   bool      `json:"ssh_reachable" hcl:"ssh_reachable"`
	LastSSHCheck   time.Time `json:"last_ssh_check,omitempty" hcl:"last_ssh_check,optional"`
	SSHLatency     int       `json:"ssh_latency,omitempty" hcl:"ssh_latency,optional"` // milliseconds
	ConnectionPool int       `json:"connection_pool,omitempty" hcl:"connection_pool,optional"`
}

// MachineStatus represents the status of a machine
type MachineStatus struct {
	Machine   *Machine               `json:"machine" hcl:"machine"`
	Status    string                 `json:"status" hcl:"status"` // online, offline, error, unknown
	LastCheck time.Time              `json:"last_check" hcl:"last_check"`
	Error     string                 `json:"error,omitempty" hcl:"error,optional"`
	Latency   int                    `json:"latency,omitempty" hcl:"latency,optional"` // milliseconds
	Details   map[string]interface{} `json:"details,omitempty" hcl:"details,optional"`
}

// MachineCollection represents a collection of machines
type MachineCollection struct {
	spookytypescommon.TimestampedEntity
	spookytypescommon.NamedEntity

	Machines []*Machine          `json:"machines" hcl:"machines"`
	Groups   map[string][]string `json:"groups,omitempty" hcl:"groups,optional"`
	Tags     map[string][]string `json:"tags,omitempty" hcl:"tags,optional"`
}

// MachineFilter represents filtering criteria for machines
type MachineFilter struct {
	Hostnames []string          `json:"hostnames,omitempty" hcl:"hostnames,optional"`
	Groups    []string          `json:"groups,omitempty" hcl:"groups,optional"`
	Roles     []string          `json:"roles,omitempty" hcl:"roles,optional"`
	Tags      map[string]string `json:"tags,omitempty" hcl:"tags,optional"`
	Patterns  []string          `json:"patterns,omitempty" hcl:"patterns,optional"`
}

// MachineQuery represents a query for machines
type MachineQuery struct {
	Filter    *MachineFilter `json:"filter,omitempty" hcl:"filter,optional"`
	Limit     int            `json:"limit,omitempty" hcl:"limit,optional"`
	Offset    int            `json:"offset,omitempty" hcl:"offset,optional"`
	SortBy    string         `json:"sort_by,omitempty" hcl:"sort_by,optional"`
	SortOrder string         `json:"sort_order,omitempty" hcl:"sort_order,optional"`
}

// MachineResult represents the result of a machine query
type MachineResult struct {
	Machines []*Machine `json:"machines" hcl:"machines"`
	Total    int        `json:"total" hcl:"total"`
	Limit    int        `json:"limit" hcl:"limit"`
	Offset   int        `json:"offset" hcl:"offset"`
}
