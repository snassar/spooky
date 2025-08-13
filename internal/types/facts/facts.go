// Package facts provides fact data structures and types.
package facts

import (
	"time"
)

// Facts represents the complete facts structure for a machine
type Facts struct {
	// System-level facts (from gopsutil)
	System *SystemFacts `json:"system" hcl:"system"`

	// Enhanced system information
	Enhanced *EnhancedFacts `json:"enhanced,omitempty" hcl:"enhanced,optional"`

	// Application facts
	Applications *ApplicationFacts `json:"applications,omitempty" hcl:"applications,optional"`

	// Deployment facts
	Deployment *DeploymentFacts `json:"deployment,omitempty" hcl:"deployment,optional"`

	// Environment facts
	Environment *EnvironmentFacts `json:"environment,omitempty" hcl:"environment,optional"`

	// Monitoring facts
	Monitoring *MonitoringFacts `json:"monitoring,omitempty" hcl:"monitoring,optional"`

	// Custom facts (user-defined)
	Custom map[string]interface{} `json:"custom,omitempty" hcl:"custom,optional"`
}

// SystemFacts represents system-level facts from gopsutil
type SystemFacts struct {
	// Operating system facts
	OS *OSFacts `json:"os" hcl:"os"`

	// Hardware facts
	Hardware *HardwareFacts `json:"hardware" hcl:"hardware"`

	// Network facts
	Network *NetworkFacts `json:"network" hcl:"network"`

	// Load average facts
	LoadAverage *LoadAverageFacts `json:"load_average,omitempty" hcl:"load_average,optional"`

	// Process information
	Processes *ProcessFacts `json:"processes,omitempty" hcl:"processes,optional"`
}

// OSFacts represents operating system facts
type OSFacts struct {
	Name     string `json:"name" hcl:"name"`
	Version  string `json:"version" hcl:"version"`
	Arch     string `json:"arch" hcl:"arch"`
	Kernel   string `json:"kernel" hcl:"kernel"`
	Platform string `json:"platform,omitempty" hcl:"platform,optional"`
	Family   string `json:"family,omitempty" hcl:"family,optional"`
}

// HardwareFacts represents hardware facts
type HardwareFacts struct {
	CPU    *CPUFacts    `json:"cpu" hcl:"cpu"`
	Memory *MemoryFacts `json:"memory" hcl:"memory"`
	Disks  []*DiskFacts `json:"disks" hcl:"disks"`
	DiskIO *DiskIOFacts `json:"disk_io,omitempty" hcl:"disk_io,optional"`
}

// CPUFacts represents CPU information
type CPUFacts struct {
	Cores        int              `json:"cores" hcl:"cores"`
	Model        string           `json:"model" hcl:"model"`
	Frequency    float64          `json:"frequency,omitempty" hcl:"frequency,optional"`
	Architecture string           `json:"architecture,omitempty" hcl:"architecture,optional"`
	Vendor       string           `json:"vendor,omitempty" hcl:"vendor,optional"`
	Times        *CPUTimes        `json:"times,omitempty" hcl:"times,optional"`
	Percent      float64          `json:"percent,omitempty" hcl:"percent,optional"`
	CoresDetail  []*CPUCoreDetail `json:"cores_detail,omitempty" hcl:"cores_detail,optional"`
}

// CPUTimes represents CPU time breakdown
type CPUTimes struct {
	User      float64 `json:"user,omitempty" hcl:"user,optional"`
	System    float64 `json:"system,omitempty" hcl:"system,optional"`
	Idle      float64 `json:"idle,omitempty" hcl:"idle,optional"`
	Nice      float64 `json:"nice,omitempty" hcl:"nice,optional"`
	IOWait    float64 `json:"iowait,omitempty" hcl:"iowait,optional"`
	IRQ       float64 `json:"irq,omitempty" hcl:"irq,optional"`
	SoftIRQ   float64 `json:"softirq,omitempty" hcl:"softirq,optional"`
	Steal     float64 `json:"steal,omitempty" hcl:"steal,optional"`
	Guest     float64 `json:"guest,omitempty" hcl:"guest,optional"`
	GuestNice float64 `json:"guest_nice,omitempty" hcl:"guest_nice,optional"`
}

// CPUCoreDetail represents detailed information for a CPU core
type CPUCoreDetail struct {
	CPU       int       `json:"cpu" hcl:"cpu"`
	ModelName string    `json:"model_name,omitempty" hcl:"model_name,optional"`
	MHz       float64   `json:"mhz,omitempty" hcl:"mhz,optional"`
	CacheSize int64     `json:"cache_size,omitempty" hcl:"cache_size,optional"`
	Percent   float64   `json:"percent,omitempty" hcl:"percent,optional"`
	Times     *CPUTimes `json:"times,omitempty" hcl:"times,optional"`
}

// MemoryFacts represents memory information
type MemoryFacts struct {
	Total         int64               `json:"total" hcl:"total"`
	Available     int64               `json:"available,omitempty" hcl:"available,optional"`
	Used          int64               `json:"used,omitempty" hcl:"used,optional"`
	Free          int64               `json:"free,omitempty" hcl:"free,optional"`
	Buffers       int64               `json:"buffers,omitempty" hcl:"buffers,optional"`
	Cached        int64               `json:"cached,omitempty" hcl:"cached,optional"`
	Shared        int64               `json:"shared,omitempty" hcl:"shared,optional"`
	Slab          int64               `json:"slab,omitempty" hcl:"slab,optional"`
	Swap          *SwapFacts          `json:"swap,omitempty" hcl:"swap,optional"`
	VirtualMemory *VirtualMemoryFacts `json:"virtual_memory,omitempty" hcl:"virtual_memory,optional"`
}

// SwapFacts represents swap memory information
type SwapFacts struct {
	Total   int64   `json:"total,omitempty" hcl:"total,optional"`
	Used    int64   `json:"used,omitempty" hcl:"used,optional"`
	Free    int64   `json:"free,omitempty" hcl:"free,optional"`
	Percent float64 `json:"percent,omitempty" hcl:"percent,optional"`
}

// VirtualMemoryFacts represents virtual memory information
type VirtualMemoryFacts struct {
	Total     int64   `json:"total,omitempty" hcl:"total,optional"`
	Available int64   `json:"available,omitempty" hcl:"available,optional"`
	Used      int64   `json:"used,omitempty" hcl:"used,optional"`
	Free      int64   `json:"free,omitempty" hcl:"free,optional"`
	Percent   float64 `json:"percent,omitempty" hcl:"percent,optional"`
}

// DiskFacts represents disk information
type DiskFacts struct {
	Device     string          `json:"device" hcl:"device"`
	MountPoint string          `json:"mount_point,omitempty" hcl:"mount_point,optional"`
	Total      int64           `json:"total" hcl:"total"`
	Used       int64           `json:"used,omitempty" hcl:"used,optional"`
	Free       int64           `json:"free,omitempty" hcl:"free,optional"`
	Filesystem string          `json:"filesystem,omitempty" hcl:"filesystem,optional"`
	IOCounters *DiskIOCounters `json:"io_counters,omitempty" hcl:"io_counters,optional"`
	Partition  *PartitionFacts `json:"partition,omitempty" hcl:"partition,optional"`
}

// DiskIOCounters represents disk I/O statistics
type DiskIOCounters struct {
	ReadCount  int64 `json:"read_count,omitempty" hcl:"read_count,optional"`
	WriteCount int64 `json:"write_count,omitempty" hcl:"write_count,optional"`
	ReadBytes  int64 `json:"read_bytes,omitempty" hcl:"read_bytes,optional"`
	WriteBytes int64 `json:"write_bytes,omitempty" hcl:"write_bytes,optional"`
	ReadTime   int64 `json:"read_time,omitempty" hcl:"read_time,optional"`
	WriteTime  int64 `json:"write_time,omitempty" hcl:"write_time,optional"`
	IOTime     int64 `json:"io_time,omitempty" hcl:"io_time,optional"`
	WeightedIO int64 `json:"weighted_io,omitempty" hcl:"weighted_io,optional"`
}

// PartitionFacts represents partition information
type PartitionFacts struct {
	Device     string `json:"device,omitempty" hcl:"device,optional"`
	Mountpoint string `json:"mountpoint,omitempty" hcl:"mountpoint,optional"`
	FSType     string `json:"fstype,omitempty" hcl:"fstype,optional"`
	Opts       string `json:"opts,omitempty" hcl:"opts,optional"`
}

// DiskIOFacts represents overall disk I/O statistics
type DiskIOFacts struct {
	ReadCount  int64 `json:"read_count,omitempty" hcl:"read_count,optional"`
	WriteCount int64 `json:"write_count,omitempty" hcl:"write_count,optional"`
	ReadBytes  int64 `json:"read_bytes,omitempty" hcl:"read_bytes,optional"`
	WriteBytes int64 `json:"write_bytes,omitempty" hcl:"write_bytes,optional"`
	ReadTime   int64 `json:"read_time,omitempty" hcl:"read_time,optional"`
	WriteTime  int64 `json:"write_time,omitempty" hcl:"write_time,optional"`
	IOTime     int64 `json:"io_time,omitempty" hcl:"io_time,optional"`
	WeightedIO int64 `json:"weighted_io,omitempty" hcl:"weighted_io,optional"`
}

// NetworkFacts represents network facts
type NetworkFacts struct {
	Hostname           string                   `json:"hostname" hcl:"hostname"`
	Interfaces         []*NetworkInterface      `json:"interfaces" hcl:"interfaces"`
	IPAddresses        []string                 `json:"ip_addresses" hcl:"ip_addresses"`
	PrimaryIP          string                   `json:"primary_ip,omitempty" hcl:"primary_ip,optional"`
	Connections        int                      `json:"connections,omitempty" hcl:"connections,optional"`
	ListeningPorts     []int                    `json:"listening_ports,omitempty" hcl:"listening_ports,optional"`
	BytesSent          int64                    `json:"bytes_sent,omitempty" hcl:"bytes_sent,optional"`
	BytesRecv          int64                    `json:"bytes_recv,omitempty" hcl:"bytes_recv,optional"`
	PacketsSent        int64                    `json:"packets_sent,omitempty" hcl:"packets_sent,optional"`
	PacketsRecv        int64                    `json:"packets_recv,omitempty" hcl:"packets_recv,optional"`
	ErrIn              int64                    `json:"err_in,omitempty" hcl:"err_in,optional"`
	ErrOut             int64                    `json:"err_out,omitempty" hcl:"err_out,optional"`
	DropIn             int64                    `json:"drop_in,omitempty" hcl:"drop_in,optional"`
	DropOut            int64                    `json:"drop_out,omitempty" hcl:"drop_out,optional"`
	Protocols          *NetworkProtocols        `json:"protocols,omitempty" hcl:"protocols,optional"`
	ConnectionDetails  []*NetworkConnection     `json:"connection_details,omitempty" hcl:"connection_details,optional"`
	InterfaceStats     []*NetworkInterfaceStats `json:"interface_stats,omitempty" hcl:"interface_stats,optional"`
	NetfilterConntrack *NetfilterConntrack      `json:"netfilter_conntrack,omitempty" hcl:"netfilter_conntrack,optional"`
}

// NetworkInterface represents a network interface
type NetworkInterface struct {
	Name        string   `json:"name" hcl:"name"`
	MACAddress  string   `json:"mac_address,omitempty" hcl:"mac_address,optional"`
	IPAddresses []string `json:"ip_addresses,omitempty" hcl:"ip_addresses,optional"`
	MTU         int      `json:"mtu,omitempty" hcl:"mtu,optional"`
	Flags       []string `json:"flags,omitempty" hcl:"flags,optional"`
}

// NetworkProtocols represents network protocol statistics
type NetworkProtocols struct {
	TCP  *TCPStats  `json:"tcp,omitempty" hcl:"tcp,optional"`
	UDP  *UDPStats  `json:"udp,omitempty" hcl:"udp,optional"`
	ICMP *ICMPStats `json:"icmp,omitempty" hcl:"icmp,optional"`
}

// TCPStats represents TCP protocol statistics
type TCPStats struct {
	Established int `json:"established,omitempty" hcl:"established,optional"`
	Listen      int `json:"listen,omitempty" hcl:"listen,optional"`
	TimeWait    int `json:"time_wait,omitempty" hcl:"time_wait,optional"`
	CloseWait   int `json:"close_wait,omitempty" hcl:"close_wait,optional"`
	FinWait1    int `json:"fin_wait1,omitempty" hcl:"fin_wait1,optional"`
	FinWait2    int `json:"fin_wait2,omitempty" hcl:"fin_wait2,optional"`
	Closing     int `json:"closing,omitempty" hcl:"closing,optional"`
	LastAck     int `json:"last_ack,omitempty" hcl:"last_ack,optional"`
	SynSent     int `json:"syn_sent,omitempty" hcl:"syn_sent,optional"`
	SynRecv     int `json:"syn_recv,omitempty" hcl:"syn_recv,optional"`
}

// UDPStats represents UDP protocol statistics
type UDPStats struct {
	Established int `json:"established,omitempty" hcl:"established,optional"`
	Listen      int `json:"listen,omitempty" hcl:"listen,optional"`
}

// ICMPStats represents ICMP protocol statistics
type ICMPStats struct {
	InMsgs    int `json:"in_msgs,omitempty" hcl:"in_msgs,optional"`
	OutMsgs   int `json:"out_msgs,omitempty" hcl:"out_msgs,optional"`
	InErrors  int `json:"in_errors,omitempty" hcl:"in_errors,optional"`
	OutErrors int `json:"out_errors,omitempty" hcl:"out_errors,optional"`
}

// NetworkConnection represents a network connection
type NetworkConnection struct {
	FD     int    `json:"fd,omitempty" hcl:"fd,optional"`
	Family int    `json:"family,omitempty" hcl:"family,optional"`
	Type   int    `json:"type,omitempty" hcl:"type,optional"`
	LAddr  string `json:"laddr,omitempty" hcl:"laddr,optional"`
	RAddr  string `json:"raddr,omitempty" hcl:"raddr,optional"`
	Status string `json:"status,omitempty" hcl:"status,optional"`
	PID    int    `json:"pid,omitempty" hcl:"pid,optional"`
}

// NetworkInterfaceStats represents detailed interface statistics
type NetworkInterfaceStats struct {
	Name          string `json:"name" hcl:"name"`
	BytesSent     int64  `json:"bytes_sent,omitempty" hcl:"bytes_sent,optional"`
	BytesRecv     int64  `json:"bytes_recv,omitempty" hcl:"bytes_recv,optional"`
	PacketsSent   int64  `json:"packets_sent,omitempty" hcl:"packets_sent,optional"`
	PacketsRecv   int64  `json:"packets_recv,omitempty" hcl:"packets_recv,optional"`
	ErrIn         int64  `json:"err_in,omitempty" hcl:"err_in,optional"`
	ErrOut        int64  `json:"err_out,omitempty" hcl:"err_out,optional"`
	DropIn        int64  `json:"drop_in,omitempty" hcl:"drop_in,optional"`
	DropOut       int64  `json:"drop_out,omitempty" hcl:"drop_out,optional"`
	FifoIn        int64  `json:"fifo_in,omitempty" hcl:"fifo_in,optional"`
	FifoOut       int64  `json:"fifo_out,omitempty" hcl:"fifo_out,optional"`
	FrameIn       int64  `json:"frame_in,omitempty" hcl:"frame_in,optional"`
	FrameOut      int64  `json:"frame_out,omitempty" hcl:"frame_out,optional"`
	CompressedIn  int64  `json:"compressed_in,omitempty" hcl:"compressed_in,optional"`
	CompressedOut int64  `json:"compressed_out,omitempty" hcl:"compressed_out,optional"`
	MulticastIn   int64  `json:"multicast_in,omitempty" hcl:"multicast_in,optional"`
	MulticastOut  int64  `json:"multicast_out,omitempty" hcl:"multicast_out,optional"`
}

// NetfilterConntrack represents netfilter connection tracking statistics
type NetfilterConntrack struct {
	Entries       int `json:"entries,omitempty" hcl:"entries,optional"`
	Searched      int `json:"searched,omitempty" hcl:"searched,optional"`
	Found         int `json:"found,omitempty" hcl:"found,optional"`
	New           int `json:"new,omitempty" hcl:"new,optional"`
	Invalid       int `json:"invalid,omitempty" hcl:"invalid,optional"`
	Ignore        int `json:"ignore,omitempty" hcl:"ignore,optional"`
	Delete        int `json:"delete,omitempty" hcl:"delete,optional"`
	DeleteList    int `json:"delete_list,omitempty" hcl:"delete_list,optional"`
	Insert        int `json:"insert,omitempty" hcl:"insert,optional"`
	InsertFailed  int `json:"insert_failed,omitempty" hcl:"insert_failed,optional"`
	Drop          int `json:"drop,omitempty" hcl:"drop,optional"`
	EarlyDrop     int `json:"early_drop,omitempty" hcl:"early_drop,optional"`
	IcmpError     int `json:"icmp_error,omitempty" hcl:"icmp_error,optional"`
	ExpectNew     int `json:"expect_new,omitempty" hcl:"expect_new,optional"`
	ExpectCreate  int `json:"expect_create,omitempty" hcl:"expect_create,optional"`
	ExpectDelete  int `json:"expect_delete,omitempty" hcl:"expect_delete,optional"`
	SearchRestart int `json:"search_restart,omitempty" hcl:"search_restart,optional"`
}

// LoadAverageFacts represents system load averages
type LoadAverageFacts struct {
	Load1  float64 `json:"load1,omitempty" hcl:"load1,optional"`
	Load5  float64 `json:"load5,omitempty" hcl:"load5,optional"`
	Load15 float64 `json:"load15,omitempty" hcl:"load15,optional"`
}

// ProcessFacts represents process information
type ProcessFacts struct {
	Count             int                    `json:"count,omitempty" hcl:"count,optional"`
	Running           int                    `json:"running,omitempty" hcl:"running,optional"`
	Sleeping          int                    `json:"sleeping,omitempty" hcl:"sleeping,optional"`
	Stopped           int                    `json:"stopped,omitempty" hcl:"stopped,optional"`
	Zombie            int                    `json:"zombie,omitempty" hcl:"zombie,optional"`
	TopByCPU          []*ProcessInfo         `json:"top_by_cpu,omitempty" hcl:"top_by_cpu,optional"`
	TopByMemory       []*ProcessInfo         `json:"top_by_memory,omitempty" hcl:"top_by_memory,optional"`
	DetailedProcesses []*DetailedProcessInfo `json:"detailed_processes,omitempty" hcl:"detailed_processes,optional"`
}

// ProcessInfo represents basic process information
type ProcessInfo struct {
	PID           int     `json:"pid" hcl:"pid"`
	Name          string  `json:"name" hcl:"name"`
	CPUPercent    float64 `json:"cpu_percent" hcl:"cpu_percent"`
	MemoryPercent float64 `json:"memory_percent,omitempty" hcl:"memory_percent,optional"`
	Cmdline       string  `json:"cmdline,omitempty" hcl:"cmdline,optional"`
	MemoryRSS     int64   `json:"memory_rss,omitempty" hcl:"memory_rss,optional"`
	MemoryVMS     int64   `json:"memory_vms,omitempty" hcl:"memory_vms,optional"`
}

// DetailedProcessInfo represents detailed process information
type DetailedProcessInfo struct {
	PID            int                  `json:"pid" hcl:"pid"`
	PPID           int                  `json:"ppid,omitempty" hcl:"ppid,optional"`
	Name           string               `json:"name" hcl:"name"`
	Cmdline        string               `json:"cmdline,omitempty" hcl:"cmdline,optional"`
	Status         string               `json:"status,omitempty" hcl:"status,optional"`
	CreateTime     int64                `json:"create_time,omitempty" hcl:"create_time,optional"`
	CWD            string               `json:"cwd,omitempty" hcl:"cwd,optional"`
	Exe            string               `json:"exe,omitempty" hcl:"exe,optional"`
	UIDs           []int                `json:"uids,omitempty" hcl:"uids,optional"`
	GIDs           []int                `json:"gids,omitempty" hcl:"gids,optional"`
	Terminal       string               `json:"terminal,omitempty" hcl:"terminal,optional"`
	Nice           int                  `json:"nice,omitempty" hcl:"nice,optional"`
	NumFDs         int                  `json:"num_fds,omitempty" hcl:"num_fds,optional"`
	NumCtxSwitches int64                `json:"num_ctx_switches,omitempty" hcl:"num_ctx_switches,optional"`
	NumThreads     int                  `json:"num_threads,omitempty" hcl:"num_threads,optional"`
	CPUTimes       *CPUTimes            `json:"cpu_times,omitempty" hcl:"cpu_times,optional"`
	MemoryInfo     *ProcessMemoryInfo   `json:"memory_info,omitempty" hcl:"memory_info,optional"`
	MemoryMaps     []*ProcessMemoryMap  `json:"memory_maps,omitempty" hcl:"memory_maps,optional"`
	OpenFiles      []*ProcessOpenFile   `json:"open_files,omitempty" hcl:"open_files,optional"`
	Connections    []*NetworkConnection `json:"connections,omitempty" hcl:"connections,optional"`
	CPUAffinity    []int                `json:"cpu_affinity,omitempty" hcl:"cpu_affinity,optional"`
	IOCounters     *ProcessIOCounters   `json:"io_counters,omitempty" hcl:"io_counters,optional"`
	PageFaults     *ProcessPageFaults   `json:"page_faults,omitempty" hcl:"page_faults,optional"`
	Username       string               `json:"username,omitempty" hcl:"username,optional"`
	Children       []int                `json:"children,omitempty" hcl:"children,optional"`
}

// ProcessMemoryInfo represents memory information for a process
type ProcessMemoryInfo struct {
	RSS    int64 `json:"rss,omitempty" hcl:"rss,optional"`
	VMS    int64 `json:"vms,omitempty" hcl:"vms,optional"`
	Shared int64 `json:"shared,omitempty" hcl:"shared,optional"`
	Text   int64 `json:"text,omitempty" hcl:"text,optional"`
	Lib    int64 `json:"lib,omitempty" hcl:"lib,optional"`
	Data   int64 `json:"data,omitempty" hcl:"data,optional"`
	Dirty  int64 `json:"dirty,omitempty" hcl:"dirty,optional"`
}

// ProcessMemoryMap represents a memory map for a process
type ProcessMemoryMap struct {
	Path string `json:"path,omitempty" hcl:"path,optional"`
	RSS  int64  `json:"rss,omitempty" hcl:"rss,optional"`
	Size int64  `json:"size,omitempty" hcl:"size,optional"`
	PSS  int64  `json:"pss,omitempty" hcl:"pss,optional"`
}

// ProcessOpenFile represents an open file for a process
type ProcessOpenFile struct {
	FD   int    `json:"fd,omitempty" hcl:"fd,optional"`
	Path string `json:"path,omitempty" hcl:"path,optional"`
}

// ProcessIOCounters represents I/O counters for a process
type ProcessIOCounters struct {
	ReadCount  int64 `json:"read_count,omitempty" hcl:"read_count,optional"`
	WriteCount int64 `json:"write_count,omitempty" hcl:"write_count,optional"`
	ReadBytes  int64 `json:"read_bytes,omitempty" hcl:"read_bytes,optional"`
	WriteBytes int64 `json:"write_bytes,omitempty" hcl:"write_bytes,optional"`
}

// ProcessPageFaults represents page fault information for a process
type ProcessPageFaults struct {
	MinorFaults      int64 `json:"minor_faults,omitempty" hcl:"minor_faults,optional"`
	MajorFaults      int64 `json:"major_faults,omitempty" hcl:"major_faults,optional"`
	ChildMinorFaults int64 `json:"child_minor_faults,omitempty" hcl:"child_minor_faults,optional"`
	ChildMajorFaults int64 `json:"child_major_faults,omitempty" hcl:"child_major_faults,optional"`
}

// EnhancedFacts represents enhanced system information
type EnhancedFacts struct {
	Virtualization *VirtualizationFacts `json:"virtualization,omitempty" hcl:"virtualization,optional"`
	PackageManager *PackageManagerFacts `json:"package_manager,omitempty" hcl:"package_manager,optional"`
	ServiceManager *ServiceManagerFacts `json:"service_manager,omitempty" hcl:"service_manager,optional"`
	SELinux        *SELinuxFacts        `json:"selinux,omitempty" hcl:"selinux,optional"`
	SSHKeys        *SSHKeysFacts        `json:"ssh_keys,omitempty" hcl:"ssh_keys,optional"`
	BIOS           *BIOSFacts           `json:"bios,omitempty" hcl:"bios,optional"`
	Sensors        *SensorsFacts        `json:"sensors,omitempty" hcl:"sensors,optional"`
	Docker         *DockerFacts         `json:"docker,omitempty" hcl:"docker,optional"`
}

// VirtualizationFacts represents virtualization information
type VirtualizationFacts struct {
	System string `json:"system,omitempty" hcl:"system,optional"`
	Role   string `json:"role,omitempty" hcl:"role,optional"`
}

// PackageManagerFacts represents package manager information
type PackageManagerFacts struct {
	Type       string `json:"type" hcl:"type"`
	Version    string `json:"version,omitempty" hcl:"version,optional"`
	ConfigPath string `json:"config_path,omitempty" hcl:"config_path,optional"`
}

// ServiceManagerFacts represents service manager information
type ServiceManagerFacts struct {
	Type    string `json:"type" hcl:"type"`
	Version string `json:"version,omitempty" hcl:"version,optional"`
}

// SELinuxFacts represents SELinux status information
type SELinuxFacts struct {
	Enabled bool   `json:"enabled" hcl:"enabled"`
	Mode    string `json:"mode,omitempty" hcl:"mode,optional"`
	Type    string `json:"type,omitempty" hcl:"type,optional"`
}

// SSHKeysFacts represents SSH host keys
type SSHKeysFacts struct {
	RSA     string `json:"rsa,omitempty" hcl:"rsa,optional"`
	ECDSA   string `json:"ecdsa,omitempty" hcl:"ecdsa,optional"`
	ED25519 string `json:"ed25519,omitempty" hcl:"ed25519,optional"`
}

// BIOSFacts represents BIOS information
type BIOSFacts struct {
	Vendor       string `json:"vendor,omitempty" hcl:"vendor,optional"`
	Version      string `json:"version,omitempty" hcl:"version,optional"`
	Date         string `json:"date,omitempty" hcl:"date,optional"`
	Release      string `json:"release,omitempty" hcl:"release,optional"`
	BoardVendor  string `json:"board_vendor,omitempty" hcl:"board_vendor,optional"`
	BoardName    string `json:"board_name,omitempty" hcl:"board_name,optional"`
	BoardVersion string `json:"board_version,omitempty" hcl:"board_version,optional"`
}

// SensorsFacts represents hardware sensor information
type SensorsFacts struct {
	Temperatures *TemperatureSensors `json:"temperatures,omitempty" hcl:"temperatures,optional"`
	Fans         *FanSensors         `json:"fans,omitempty" hcl:"fans,optional"`
}

// TemperatureSensors represents temperature sensors
type TemperatureSensors struct {
	CPU         float64   `json:"cpu,omitempty" hcl:"cpu,optional"`
	GPU         float64   `json:"gpu,omitempty" hcl:"gpu,optional"`
	Motherboard float64   `json:"motherboard,omitempty" hcl:"motherboard,optional"`
	CoreTemp    []float64 `json:"core_temp,omitempty" hcl:"core_temp,optional"`
}

// FanSensors represents fan speeds
type FanSensors struct {
	CPUFan    int `json:"cpu_fan,omitempty" hcl:"cpu_fan,optional"`
	SystemFan int `json:"system_fan,omitempty" hcl:"system_fan,optional"`
	CaseFan   int `json:"case_fan,omitempty" hcl:"case_fan,optional"`
}

// DockerFacts represents Docker information
type DockerFacts struct {
	ContainerID   string               `json:"container_id,omitempty" hcl:"container_id,optional"`
	CgroupsCPU    *DockerCgroupsCPU    `json:"cgroups_cpu,omitempty" hcl:"cgroups_cpu,optional"`
	CgroupsMemory *DockerCgroupsMemory `json:"cgroups_memory,omitempty" hcl:"cgroups_memory,optional"`
}

// DockerCgroupsCPU represents Docker cgroups CPU information
type DockerCgroupsCPU struct {
	User   float64 `json:"user,omitempty" hcl:"user,optional"`
	System float64 `json:"system,omitempty" hcl:"system,optional"`
}

// DockerCgroupsMemory represents Docker cgroups memory information
type DockerCgroupsMemory struct {
	Usage    int64 `json:"usage,omitempty" hcl:"usage,optional"`
	Limit    int64 `json:"limit,omitempty" hcl:"limit,optional"`
	MaxUsage int64 `json:"max_usage,omitempty" hcl:"max_usage,optional"`
}

// ApplicationFacts represents application-specific facts
type ApplicationFacts struct {
	Versions *ApplicationVersions `json:"versions,omitempty" hcl:"versions,optional"`
	Config   *ApplicationConfig   `json:"config,omitempty" hcl:"config,optional"`
}

// ApplicationVersions represents application version information
type ApplicationVersions struct {
	Nginx      string `json:"nginx,omitempty" hcl:"nginx,optional"`
	Apache     string `json:"apache,omitempty" hcl:"apache,optional"`
	PostgreSQL string `json:"postgresql,omitempty" hcl:"postgresql,optional"`
	MySQL      string `json:"mysql,omitempty" hcl:"mysql,optional"`
	Redis      string `json:"redis,omitempty" hcl:"redis,optional"`
	Docker     string `json:"docker,omitempty" hcl:"docker,optional"`
	Kubernetes string `json:"kubernetes,omitempty" hcl:"kubernetes,optional"`
}

// ApplicationConfig represents application configuration facts
type ApplicationConfig struct {
	ConfigPaths *ConfigPaths `json:"config_paths,omitempty" hcl:"config_paths,optional"`
	LogPaths    *LogPaths    `json:"log_paths,omitempty" hcl:"log_paths,optional"`
}

// ConfigPaths represents configuration file paths
type ConfigPaths struct {
	NginxConf    string `json:"nginx_conf,omitempty" hcl:"nginx_conf,optional"`
	ApacheConf   string `json:"apache_conf,omitempty" hcl:"apache_conf,optional"`
	PostgresConf string `json:"postgres_conf,omitempty" hcl:"postgres_conf,optional"`
	RedisConf    string `json:"redis_conf,omitempty" hcl:"redis_conf,optional"`
}

// LogPaths represents log file paths
type LogPaths struct {
	NginxLogs    string `json:"nginx_logs,omitempty" hcl:"nginx_logs,optional"`
	ApacheLogs   string `json:"apache_logs,omitempty" hcl:"apache_logs,optional"`
	PostgresLogs string `json:"postgres_logs,omitempty" hcl:"postgres_logs,optional"`
	RedisLogs    string `json:"redis_logs,omitempty" hcl:"redis_logs,optional"`
}

// DeploymentFacts represents deployment-related facts
type DeploymentFacts struct {
	State    *DeploymentState    `json:"state,omitempty" hcl:"state,optional"`
	Info     *DeploymentInfo     `json:"info,omitempty" hcl:"info,optional"`
	Services *DeploymentServices `json:"services,omitempty" hcl:"services,optional"`
}

// DeploymentState represents deployment state
type DeploymentState struct {
	State string `json:"state,omitempty" hcl:"state,optional"`
}

// DeploymentInfo represents deployment information
type DeploymentInfo struct {
	Version    string    `json:"version,omitempty" hcl:"version,optional"`
	DeployedAt time.Time `json:"deployed_at,omitempty" hcl:"deployed_at,optional"`
	DeployedBy string    `json:"deployed_by,omitempty" hcl:"deployed_by,optional"`
	CommitHash string    `json:"commit_hash,omitempty" hcl:"commit_hash,optional"`
	Branch     string    `json:"branch,omitempty" hcl:"branch,optional"`
}

// DeploymentServices represents service status information
type DeploymentServices struct {
	Nginx      string `json:"nginx,omitempty" hcl:"nginx,optional"`
	Apache     string `json:"apache,omitempty" hcl:"apache,optional"`
	PostgreSQL string `json:"postgresql,omitempty" hcl:"postgresql,optional"`
	Redis      string `json:"redis,omitempty" hcl:"redis,optional"`
}

// EnvironmentFacts represents environment-specific facts
type EnvironmentFacts struct {
	Variables      *EnvironmentVariables `json:"variables,omitempty" hcl:"variables,optional"`
	Infrastructure *InfrastructureFacts  `json:"infrastructure,omitempty" hcl:"infrastructure,optional"`
}

// EnvironmentVariables represents environment variables
type EnvironmentVariables struct {
	NodeEnv     string `json:"NODE_ENV,omitempty" hcl:"NODE_ENV,optional"`
	DatabaseURL string `json:"DATABASE_URL,omitempty" hcl:"DATABASE_URL,optional"`
	RedisURL    string `json:"REDIS_URL,omitempty" hcl:"REDIS_URL,optional"`
	LogLevel    string `json:"LOG_LEVEL,omitempty" hcl:"LOG_LEVEL,optional"`
}

// InfrastructureFacts represents infrastructure information
type InfrastructureFacts struct {
	Datacenter       string `json:"datacenter,omitempty" hcl:"datacenter,optional"`
	Rack             string `json:"rack,omitempty" hcl:"rack,optional"`
	PowerZone        string `json:"power_zone,omitempty" hcl:"power_zone,optional"`
	Region           string `json:"region,omitempty" hcl:"region,optional"`
	AvailabilityZone string `json:"availability_zone,omitempty" hcl:"availability_zone,optional"`
}

// MonitoringFacts represents monitoring configuration facts
type MonitoringFacts struct {
	Endpoints    *MonitoringEndpoints `json:"endpoints,omitempty" hcl:"endpoints,optional"`
	HealthChecks *HealthCheckFacts    `json:"health_checks,omitempty" hcl:"health_checks,optional"`
}

// MonitoringEndpoints represents monitoring endpoints
type MonitoringEndpoints struct {
	PrometheusPort int    `json:"prometheus_port,omitempty" hcl:"prometheus_port,optional"`
	GrafanaPort    int    `json:"grafana_port,omitempty" hcl:"grafana_port,optional"`
	AlertManager   string `json:"alert_manager,omitempty" hcl:"alert_manager,optional"`
}

// HealthCheckFacts represents health check configuration
type HealthCheckFacts struct {
	Enabled  bool   `json:"enabled,omitempty" hcl:"enabled,optional"`
	Port     int    `json:"port,omitempty" hcl:"port,optional"`
	Path     string `json:"path,omitempty" hcl:"path,optional"`
	Interval string `json:"interval,omitempty" hcl:"interval,optional"`
}
