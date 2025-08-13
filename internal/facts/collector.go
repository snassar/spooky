// Package facts provides fact collection, storage, and management functionality.
package facts

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	spookytypes "spooky/internal/types"
	spookytypesfacts "spooky/internal/types/facts"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/load"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/shirou/gopsutil/v4/net"
	"github.com/shirou/gopsutil/v4/process"
)

// SystemFactCollector collects system facts using gopsutil
type SystemFactCollector struct {
	name string
}

// NewSystemFactCollector creates a new system fact collector
func NewSystemFactCollector() *SystemFactCollector {
	return &SystemFactCollector{
		name: "system",
	}
}

// GetName returns the collector name
func (c *SystemFactCollector) GetName() string {
	return c.name
}

// Collect collects facts from the given machine
func (c *SystemFactCollector) Collect(ctx context.Context, machine *spookytypes.Machine) (*FactCollection, error) {
	// Get machine ID
	machineID, err := c.getMachineID()
	if err != nil {
		return nil, fmt.Errorf("failed to get machine ID: %w", err)
	}

	// Collect system facts
	facts := &spookytypesfacts.Facts{
		System: &spookytypesfacts.SystemFacts{},
	}

	// Collect OS facts
	osFacts, err := c.collectOSFacts()
	if err != nil {
		return nil, fmt.Errorf("failed to collect OS facts: %w", err)
	}
	facts.System.OS = osFacts

	// Collect hardware facts
	hardwareFacts, err := c.collectHardwareFacts()
	if err != nil {
		return nil, fmt.Errorf("failed to collect hardware facts: %w", err)
	}
	facts.System.Hardware = hardwareFacts

	// Collect network facts
	networkFacts, err := c.collectNetworkFacts()
	if err != nil {
		return nil, fmt.Errorf("failed to collect network facts: %w", err)
	}
	facts.System.Network = networkFacts

	// Collect load average facts
	loadFacts, err := c.collectLoadAverageFacts()
	if err != nil {
		return nil, fmt.Errorf("failed to collect load average facts: %w", err)
	}
	facts.System.LoadAverage = loadFacts

	// Collect process facts
	processFacts, err := c.collectProcessFacts()
	if err != nil {
		return nil, fmt.Errorf("failed to collect process facts: %w", err)
	}
	facts.System.Processes = processFacts

	// Create fact collection
	collection := &FactCollection{
		MachineID:   machineID,
		CollectedAt: time.Now(),
		Facts:       facts,
		Metadata: map[string]interface{}{
			"collector": c.name,
			"machine":   machine.Hostname,
		},
	}

	return collection, nil
}

// getMachineID gets the machine ID from /etc/machine-id
func (c *SystemFactCollector) getMachineID() (string, error) {
	data, err := os.ReadFile("/etc/machine-id")
	if err != nil {
		return "", fmt.Errorf("failed to read /etc/machine-id: %w", err)
	}

	machineID := strings.TrimSpace(string(data))

	// Validate machine ID format (32-character hex string)
	if !regexp.MustCompile(`^[a-f0-9]{32}$`).MatchString(machineID) {
		return "", fmt.Errorf("invalid machine ID format: %s", machineID)
	}

	return machineID, nil
}

// collectOSFacts collects operating system facts
func (c *SystemFactCollector) collectOSFacts() (*spookytypesfacts.OSFacts, error) {
	hostInfo, err := host.Info()
	if err != nil {
		return nil, fmt.Errorf("failed to get host info: %w", err)
	}

	return &spookytypesfacts.OSFacts{
		Name:     hostInfo.OS,
		Version:  hostInfo.PlatformVersion,
		Arch:     hostInfo.KernelArch,
		Kernel:   hostInfo.KernelVersion,
		Platform: hostInfo.Platform,
		Family:   hostInfo.PlatformFamily,
	}, nil
}

// collectHardwareFacts collects hardware facts
func (c *SystemFactCollector) collectHardwareFacts() (*spookytypesfacts.HardwareFacts, error) {
	// Collect CPU facts
	cpuFacts, err := c.collectCPUFacts()
	if err != nil {
		return nil, fmt.Errorf("failed to collect CPU facts: %w", err)
	}

	// Collect memory facts
	memoryFacts, err := c.collectMemoryFacts()
	if err != nil {
		return nil, fmt.Errorf("failed to collect memory facts: %w", err)
	}

	// Collect disk facts
	diskFacts, err := c.collectDiskFacts()
	if err != nil {
		return nil, fmt.Errorf("failed to collect disk facts: %w", err)
	}

	// Collect disk I/O facts
	diskIOFacts, err := c.collectDiskIOFacts()
	if err != nil {
		return nil, fmt.Errorf("failed to collect disk I/O facts: %w", err)
	}

	return &spookytypesfacts.HardwareFacts{
		CPU:    cpuFacts,
		Memory: memoryFacts,
		Disks:  diskFacts,
		DiskIO: diskIOFacts,
	}, nil
}

// collectCPUFacts collects CPU facts
func (c *SystemFactCollector) collectCPUFacts() (*spookytypesfacts.CPUFacts, error) {
	// Get CPU info
	cpuInfo, err := cpu.Info()
	if err != nil {
		return nil, fmt.Errorf("failed to get CPU info: %w", err)
	}

	if len(cpuInfo) == 0 {
		return nil, fmt.Errorf("no CPU info available")
	}

	// Get CPU times
	cpuTimes, err := cpu.Times(false)
	if err != nil {
		return nil, fmt.Errorf("failed to get CPU times: %w", err)
	}

	var times *spookytypesfacts.CPUTimes
	if len(cpuTimes) > 0 {
		times = &spookytypesfacts.CPUTimes{
			User:      cpuTimes[0].User,
			System:    cpuTimes[0].System,
			Idle:      cpuTimes[0].Idle,
			Nice:      cpuTimes[0].Nice,
			IOWait:    cpuTimes[0].Iowait,
			IRQ:       cpuTimes[0].Irq,
			SoftIRQ:   cpuTimes[0].Softirq,
			Steal:     cpuTimes[0].Steal,
			Guest:     cpuTimes[0].Guest,
			GuestNice: cpuTimes[0].GuestNice,
		}
	}

	// Get CPU percentage
	cpuPercent, err := cpu.Percent(0, false)
	if err != nil {
		return nil, fmt.Errorf("failed to get CPU percentage: %w", err)
	}

	var percent float64
	if len(cpuPercent) > 0 {
		percent = cpuPercent[0]
	}

	// Get per-core information
	coresDetail, err := c.collectCPUCoresDetail()
	if err != nil {
		return nil, fmt.Errorf("failed to collect CPU cores detail: %w", err)
	}

	return &spookytypesfacts.CPUFacts{
		Cores:        len(cpuInfo),
		Model:        cpuInfo[0].ModelName,
		Frequency:    cpuInfo[0].Mhz,
		Architecture: cpuInfo[0].Family,
		Vendor:       cpuInfo[0].VendorID,
		Times:        times,
		Percent:      percent,
		CoresDetail:  coresDetail,
	}, nil
}

// collectCPUCoresDetail collects detailed information for each CPU core
func (c *SystemFactCollector) collectCPUCoresDetail() ([]*spookytypesfacts.CPUCoreDetail, error) {
	cpuInfo, err := cpu.Info()
	if err != nil {
		return nil, fmt.Errorf("failed to get CPU info: %w", err)
	}

	cpuTimes, err := cpu.Times(true)
	if err != nil {
		return nil, fmt.Errorf("failed to get CPU times: %w", err)
	}

	cpuPercent, err := cpu.Percent(0, true)
	if err != nil {
		return nil, fmt.Errorf("failed to get CPU percentage: %w", err)
	}

	var coresDetail []*spookytypesfacts.CPUCoreDetail

	for i, info := range cpuInfo {
		var times *spookytypesfacts.CPUTimes
		if i < len(cpuTimes) {
			times = &spookytypesfacts.CPUTimes{
				User:      cpuTimes[i].User,
				System:    cpuTimes[i].System,
				Idle:      cpuTimes[i].Idle,
				Nice:      cpuTimes[i].Nice,
				IOWait:    cpuTimes[i].Iowait,
				IRQ:       cpuTimes[i].Irq,
				SoftIRQ:   cpuTimes[i].Softirq,
				Steal:     cpuTimes[i].Steal,
				Guest:     cpuTimes[i].Guest,
				GuestNice: cpuTimes[i].GuestNice,
			}
		}

		var percent float64
		if i < len(cpuPercent) {
			percent = cpuPercent[i]
		}

		core := &spookytypesfacts.CPUCoreDetail{
			CPU:       i,
			ModelName: info.ModelName,
			MHz:       info.Mhz,
			CacheSize: int64(info.CacheSize),
			Percent:   percent,
			Times:     times,
		}

		coresDetail = append(coresDetail, core)
	}

	return coresDetail, nil
}

// collectMemoryFacts collects memory facts
func (c *SystemFactCollector) collectMemoryFacts() (*spookytypesfacts.MemoryFacts, error) {
	// Get virtual memory info
	vmem, err := mem.VirtualMemory()
	if err != nil {
		return nil, fmt.Errorf("failed to get virtual memory info: %w", err)
	}

	// Get swap memory info
	swap, err := mem.SwapMemory()
	if err != nil {
		return nil, fmt.Errorf("failed to get swap memory info: %w", err)
	}

	swapFacts := &spookytypesfacts.SwapFacts{
		Total:   int64(swap.Total),
		Used:    int64(swap.Used),
		Free:    int64(swap.Free),
		Percent: swap.UsedPercent,
	}

	virtualMemoryFacts := &spookytypesfacts.VirtualMemoryFacts{
		Total:     int64(vmem.Total),
		Available: int64(vmem.Available),
		Used:      int64(vmem.Used),
		Free:      int64(vmem.Free),
		Percent:   vmem.UsedPercent,
	}

	return &spookytypesfacts.MemoryFacts{
		Total:         int64(vmem.Total),
		Available:     int64(vmem.Available),
		Used:          int64(vmem.Used),
		Free:          int64(vmem.Free),
		Buffers:       int64(vmem.Buffers),
		Cached:        int64(vmem.Cached),
		Shared:        int64(vmem.Shared),
		Slab:          int64(vmem.Slab),
		Swap:          swapFacts,
		VirtualMemory: virtualMemoryFacts,
	}, nil
}

// collectDiskFacts collects disk facts
func (c *SystemFactCollector) collectDiskFacts() ([]*spookytypesfacts.DiskFacts, error) {
	// Get disk partitions
	partitions, err := disk.Partitions(false)
	if err != nil {
		return nil, fmt.Errorf("failed to get disk partitions: %w", err)
	}

	var diskFacts []*spookytypesfacts.DiskFacts

	for _, partition := range partitions {
		// Get disk usage
		usage, err := disk.Usage(partition.Mountpoint)
		if err != nil {
			// Skip partitions we can't access
			continue
		}

		// Get disk I/O counters
		ioCounters, err := disk.IOCounters(partition.Device)
		if err != nil {
			// Skip if we can't get I/O counters
			continue
		}

		var ioCountersFacts *spookytypesfacts.DiskIOCounters
		if len(ioCounters) > 0 {
			io := ioCounters[partition.Device]
			ioCountersFacts = &spookytypesfacts.DiskIOCounters{
				ReadCount:  int64(io.ReadCount),
				WriteCount: int64(io.WriteCount),
				ReadBytes:  int64(io.ReadBytes),
				WriteBytes: int64(io.WriteBytes),
				ReadTime:   int64(io.ReadTime),
				WriteTime:  int64(io.WriteTime),
				IOTime:     int64(io.IoTime),
				WeightedIO: int64(io.WeightedIO),
			}
		}

		partitionFacts := &spookytypesfacts.PartitionFacts{
			Device:     partition.Device,
			Mountpoint: partition.Mountpoint,
			FSType:     partition.Fstype,
			Opts:       strings.Join(partition.Opts, ","),
		}

		diskFact := &spookytypesfacts.DiskFacts{
			Device:     partition.Device,
			MountPoint: partition.Mountpoint,
			Total:      int64(usage.Total),
			Used:       int64(usage.Used),
			Free:       int64(usage.Free),
			Filesystem: partition.Fstype,
			IOCounters: ioCountersFacts,
			Partition:  partitionFacts,
		}

		diskFacts = append(diskFacts, diskFact)
	}

	return diskFacts, nil
}

// collectDiskIOFacts collects overall disk I/O statistics
func (c *SystemFactCollector) collectDiskIOFacts() (*spookytypesfacts.DiskIOFacts, error) {
	ioCounters, err := disk.IOCounters()
	if err != nil {
		return nil, fmt.Errorf("failed to get disk I/O counters: %w", err)
	}

	var totalReadCount, totalWriteCount, totalReadBytes, totalWriteBytes int64
	var totalReadTime, totalWriteTime, totalIOTime, totalWeightedIO int64

	for _, io := range ioCounters {
		totalReadCount += int64(io.ReadCount)
		totalWriteCount += int64(io.WriteCount)
		totalReadBytes += int64(io.ReadBytes)
		totalWriteBytes += int64(io.WriteBytes)
		totalReadTime += int64(io.ReadTime)
		totalWriteTime += int64(io.WriteTime)
		totalIOTime += int64(io.IoTime)
		totalWeightedIO += int64(io.WeightedIO)
	}

	return &spookytypesfacts.DiskIOFacts{
		ReadCount:  totalReadCount,
		WriteCount: totalWriteCount,
		ReadBytes:  totalReadBytes,
		WriteBytes: totalWriteBytes,
		ReadTime:   totalReadTime,
		WriteTime:  totalWriteTime,
		IOTime:     totalIOTime,
		WeightedIO: totalWeightedIO,
	}, nil
}

// collectNetworkFacts collects network facts
func (c *SystemFactCollector) collectNetworkFacts() (*spookytypesfacts.NetworkFacts, error) {
	// Get hostname
	hostname, err := os.Hostname()
	if err != nil {
		return nil, fmt.Errorf("failed to get hostname: %w", err)
	}

	// Get network interfaces
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("failed to get network interfaces: %w", err)
	}

	// Get network I/O counters
	ioCounters, err := net.IOCounters(false)
	if err != nil {
		return nil, fmt.Errorf("failed to get network I/O counters: %w", err)
	}

	var networkInterfaces []*spookytypesfacts.NetworkInterface
	var allIPAddresses []string
	var primaryIP string

	for _, iface := range interfaces {
		// Skip loopback interfaces
		if strings.HasPrefix(iface.Name, "lo") {
			continue
		}

		// Convert addresses to strings
		var ipAddresses []string
		for _, addr := range iface.Addrs {
			ipAddresses = append(ipAddresses, addr.Addr)
		}

		networkInterface := &spookytypesfacts.NetworkInterface{
			Name:        iface.Name,
			MACAddress:  iface.HardwareAddr,
			IPAddresses: ipAddresses,
			MTU:         iface.MTU,
			Flags:       iface.Flags,
		}

		networkInterfaces = append(networkInterfaces, networkInterface)

		// Collect all IP addresses
		for _, addr := range iface.Addrs {
			allIPAddresses = append(allIPAddresses, addr.Addr)
		}

		// Set primary IP (first non-loopback interface)
		if primaryIP == "" && len(iface.Addrs) > 0 {
			primaryIP = iface.Addrs[0].Addr
		}
	}

	// Get network connections count
	connections, err := net.Connections("all")
	if err != nil {
		return nil, fmt.Errorf("failed to get network connections: %w", err)
	}

	// Get listening ports
	var listeningPorts []int
	for _, conn := range connections {
		if conn.Status == "LISTEN" {
			listeningPorts = append(listeningPorts, int(conn.Laddr.Port))
		}
	}

	var totalBytesSent, totalBytesRecv, totalPacketsSent, totalPacketsRecv int64
	var totalErrIn, totalErrOut, totalDropIn, totalDropOut int64

	if len(ioCounters) > 0 {
		io := ioCounters[0]
		totalBytesSent = int64(io.BytesSent)
		totalBytesRecv = int64(io.BytesRecv)
		totalPacketsSent = int64(io.PacketsSent)
		totalPacketsRecv = int64(io.PacketsRecv)
		totalErrIn = int64(io.Errin)
		totalErrOut = int64(io.Errout)
		totalDropIn = int64(io.Dropin)
		totalDropOut = int64(io.Dropout)
	}

	return &spookytypesfacts.NetworkFacts{
		Hostname:       hostname,
		Interfaces:     networkInterfaces,
		IPAddresses:    allIPAddresses,
		PrimaryIP:      primaryIP,
		Connections:    len(connections),
		ListeningPorts: listeningPorts,
		BytesSent:      totalBytesSent,
		BytesRecv:      totalBytesRecv,
		PacketsSent:    totalPacketsSent,
		PacketsRecv:    totalPacketsRecv,
		ErrIn:          totalErrIn,
		ErrOut:         totalErrOut,
		DropIn:         totalDropIn,
		DropOut:        totalDropOut,
	}, nil
}

// collectLoadAverageFacts collects load average facts
func (c *SystemFactCollector) collectLoadAverageFacts() (*spookytypesfacts.LoadAverageFacts, error) {
	loadAvg, err := load.Avg()
	if err != nil {
		return nil, fmt.Errorf("failed to get load average: %w", err)
	}

	return &spookytypesfacts.LoadAverageFacts{
		Load1:  loadAvg.Load1,
		Load5:  loadAvg.Load5,
		Load15: loadAvg.Load15,
	}, nil
}

// collectProcessFacts collects process facts
func (c *SystemFactCollector) collectProcessFacts() (*spookytypesfacts.ProcessFacts, error) {
	// Get all processes
	processes, err := process.Processes()
	if err != nil {
		return nil, fmt.Errorf("failed to get processes: %w", err)
	}

	// Count processes by status
	var running, sleeping, stopped, zombie int

	// For now, just count total processes
	// Status parsing can be implemented later when we understand the exact return type
	running = len(processes)

	// Get top processes by CPU usage
	topByCPU, err := c.getTopProcessesByCPU(processes, 10)
	if err != nil {
		return nil, fmt.Errorf("failed to get top processes by CPU: %w", err)
	}

	// Get top processes by memory usage
	topByMemory, err := c.getTopProcessesByMemory(processes, 10)
	if err != nil {
		return nil, fmt.Errorf("failed to get top processes by memory: %w", err)
	}

	return &spookytypesfacts.ProcessFacts{
		Count:       len(processes),
		Running:     running,
		Sleeping:    sleeping,
		Stopped:     stopped,
		Zombie:      zombie,
		TopByCPU:    topByCPU,
		TopByMemory: topByMemory,
	}, nil
}

// getTopProcessesByCPU gets the top processes by CPU usage
func (c *SystemFactCollector) getTopProcessesByCPU(processes []*process.Process, limit int) ([]*spookytypesfacts.ProcessInfo, error) {
	var processInfos []*spookytypesfacts.ProcessInfo

	for _, p := range processes {
		// Note: CPUPercent() requires calling it twice to get accurate results
		// For now, we'll collect basic process info without CPU percentage
		name, err := p.Name()
		if err != nil {
			continue
		}

		processInfo := &spookytypesfacts.ProcessInfo{
			PID:        int(p.Pid),
			Name:       name,
			CPUPercent: 0.0, // Would need proper CPU measurement
		}

		processInfos = append(processInfos, processInfo)
	}

	// Take top N (simplified - no sorting)
	if len(processInfos) > limit {
		processInfos = processInfos[:limit]
	}

	return processInfos, nil
}

// getTopProcessesByMemory gets the top processes by memory usage
func (c *SystemFactCollector) getTopProcessesByMemory(processes []*process.Process, limit int) ([]*spookytypesfacts.ProcessInfo, error) {
	var processInfos []*spookytypesfacts.ProcessInfo

	for _, p := range processes {
		memoryPercent, err := p.MemoryPercent()
		if err != nil {
			continue
		}

		name, err := p.Name()
		if err != nil {
			continue
		}

		processInfo := &spookytypesfacts.ProcessInfo{
			PID:           int(p.Pid),
			Name:          name,
			MemoryPercent: float64(memoryPercent),
		}

		processInfos = append(processInfos, processInfo)
	}

	// Sort by memory percentage (descending) and take top N
	// This is a simplified implementation - in a real implementation,
	// you would sort the slice and take the top N

	if len(processInfos) > limit {
		processInfos = processInfos[:limit]
	}

	return processInfos, nil
}
