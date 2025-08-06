package local

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	spookyfactscollectors "spooky/internal/facts/collectors"
	spookyfactstypes "spooky/internal/facts/types"

	machinecpu "github.com/shirou/gopsutil/v4/cpu"
	machinedisk "github.com/shirou/gopsutil/v4/disk"
	machinehost "github.com/shirou/gopsutil/v4/host"
	machineload "github.com/shirou/gopsutil/v4/load"
	machinemem "github.com/shirou/gopsutil/v4/mem"
	machinenet "github.com/shirou/gopsutil/v4/net"
	machineprocess "github.com/shirou/gopsutil/v4/process"
	machinesensors "github.com/shirou/gopsutil/v4/sensors"
)

// Collector collects facts from the local system
type Collector struct {
	spookyfactscollectors.BaseCollector
}

// NewCollector creates a new local fact collector
func NewCollector() *Collector {
	return &Collector{
		BaseCollector: *spookyfactscollectors.NewBaseCollector(spookyfactstypes.SourceLocal, spookyfactstypes.MergePolicyReplace),
	}
}

// Collect gathers all available facts from the local system
func (c *Collector) Collect(server string) (*spookyfactstypes.FactCollection, error) {
	collection := &spookyfactstypes.FactCollection{
		Server:    server,
		Timestamp: time.Now(),
		Facts:     make(map[string]*spookyfactstypes.Fact),
	}

	// Collect system facts
	if err := c.collectSystemFacts(collection); err != nil {
		return nil, fmt.Errorf("failed to collect system facts: %w", err)
	}

	// Collect OS facts
	if err := c.collectOSFacts(collection); err != nil {
		return nil, fmt.Errorf("failed to collect OS facts: %w", err)
	}

	// Collect hardware facts
	if err := c.collectHardwareFacts(collection); err != nil {
		return nil, fmt.Errorf("failed to collect hardware facts: %w", err)
	}

	// Collect network facts
	if err := c.collectNetworkFacts(collection); err != nil {
		return nil, fmt.Errorf("failed to collect network facts: %w", err)
	}

	// Collect environment facts
	if err := c.collectEnvironmentFacts(collection); err != nil {
		return nil, fmt.Errorf("failed to collect environment facts: %w", err)
	}

	// Collect enhanced facts using gopsutil
	if err := c.collectEnhancedFacts(collection); err != nil {
		return nil, fmt.Errorf("failed to collect enhanced facts: %w", err)
	}

	// Collect Spooky minimal facts
	if err := c.collectSpookyMinimalFacts(collection); err != nil {
		return nil, fmt.Errorf("failed to collect Spooky minimal facts: %w", err)
	}

	return collection, nil
}

// CollectSpecific collects only the specified facts
func (c *Collector) CollectSpecific(server string, keys []string) (*spookyfactstypes.FactCollection, error) {
	collection := &spookyfactstypes.FactCollection{
		Server:    server,
		Timestamp: time.Now(),
		Facts:     make(map[string]*spookyfactstypes.Fact),
	}

	for _, key := range keys {
		if err := c.collectSpecificFact(collection, key); err != nil {
			return nil, fmt.Errorf("failed to collect fact %s: %w", key, err)
		}
	}

	return collection, nil
}

// GetFact retrieves a single fact
func (c *Collector) GetFact(server, key string) (*spookyfactstypes.Fact, error) {
	collection := &spookyfactstypes.FactCollection{
		Server:    server,
		Timestamp: time.Now(),
		Facts:     make(map[string]*spookyfactstypes.Fact),
	}

	if err := c.collectSpecificFact(collection, key); err != nil {
		return nil, fmt.Errorf("failed to collect fact %s: %w", key, err)
	}

	fact, exists := collection.Facts[key]
	if !exists {
		return nil, fmt.Errorf("fact %s not found", key)
	}

	return fact, nil
}

// Validate validates the collector configuration
func (c *Collector) Validate() error {
	// Local collector doesn't require special validation
	return nil
}

// collectSystemFacts collects basic system information
func (c *Collector) collectSystemFacts(collection *spookyfactstypes.FactCollection) error {
	// Hostname
	if hostname, err := os.Hostname(); err == nil {
		c.createFact(collection, "hostname", hostname)
	}

	// Machine ID
	if machineID, err := c.getMachineID(); err == nil {
		c.createFact(collection, "machine_id", machineID)
	}

	// FQDN
	if fqdn, err := c.getFQDN(); err == nil {
		c.createFact(collection, "fqdn", fqdn)
	}

	return nil
}

// collectOSFacts collects operating system information
func (c *Collector) collectOSFacts(collection *spookyfactstypes.FactCollection) error {
	// OS name
	c.createFact(collection, "os_name", runtime.GOOS)

	// OS architecture
	c.createFact(collection, "os_arch", runtime.GOARCH)

	// OS version (platform specific)
	if version, err := c.getOSVersion(); err == nil {
		c.createFact(collection, "os_version", version)
	}

	// OS distribution (platform specific)
	if distro, err := c.getOSDistro(); err == nil {
		c.createFact(collection, "os_distro", distro)
	}

	// Kernel version (platform specific)
	if kernel, err := c.getKernelVersion(); err == nil {
		c.createFact(collection, "kernel_version", kernel)
	}

	return nil
}

// collectHardwareFacts collects hardware information
func (c *Collector) collectHardwareFacts(collection *spookyfactstypes.FactCollection) error {
	// CPU cores
	c.createFactWithValue(collection, "cpu_cores", runtime.NumCPU())

	// CPU model (platform specific)
	if model, err := c.getCPUModel(); err == nil {
		c.createFact(collection, "cpu_model", model)
	}

	// Memory information (platform specific)
	if memInfo, err := c.getMemoryInfo(); err == nil {
		c.createFactWithValue(collection, "memory_total", memInfo.Total)
		c.createFactWithValue(collection, "memory_used", memInfo.Used)
		c.createFactWithValue(collection, "memory_available", memInfo.Available)
	}

	// Disk information (platform specific)
	if diskInfo, err := c.getDiskInfo(); err == nil {
		c.createFactWithValue(collection, "disk_total", diskInfo.Total)
		c.createFactWithValue(collection, "disk_used", diskInfo.Used)
		c.createFactWithValue(collection, "disk_available", diskInfo.Available)
	}

	return nil
}

// collectNetworkFacts collects network information
func (c *Collector) collectNetworkFacts(collection *spookyfactstypes.FactCollection) error {
	// IP addresses
	if ips, err := c.getIPAddresses(); err == nil {
		c.createFactWithValue(collection, "ip_addresses", ips)
	}

	// MAC addresses
	if macs, err := c.getMACAddresses(); err == nil {
		c.createFactWithValue(collection, "mac_addresses", macs)
	}

	// DNS configuration
	if dns, err := c.getDNSConfig(); err == nil {
		c.createFactWithValue(collection, "dns_servers", dns.Nameservers)
		c.createFactWithValue(collection, "dns_search", dns.Search)
	}

	return nil
}

// collectEnvironmentFacts collects environment information
func (c *Collector) collectEnvironmentFacts(collection *spookyfactstypes.FactCollection) error {
	// Environment variables
	env := make(map[string]string)
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		if len(pair) == 2 {
			env[pair[0]] = pair[1]
		}
	}
	c.createFactWithValue(collection, "environment", env)

	return nil
}

// collectEnhancedFacts collects comprehensive system information using gopsutil
func (c *Collector) collectEnhancedFacts(collection *spookyfactstypes.FactCollection) error {
	// Collect load averages
	if loadAvg, err := machineload.Avg(); err == nil {
		c.createFactWithValue(collection, "load_average", map[string]interface{}{
			"load1":  loadAvg.Load1,
			"load5":  loadAvg.Load5,
			"load15": loadAvg.Load15,
		})
	}

	// Collect host information
	if hostInfo, err := machinehost.Info(); err == nil {
		c.createFactWithValue(collection, "host_info", map[string]interface{}{
			"hostname":              hostInfo.Hostname,
			"uptime":                hostInfo.Uptime,
			"boot_time":             hostInfo.BootTime,
			"procs":                 hostInfo.Procs,
			"os":                    hostInfo.OS,
			"platform":              hostInfo.Platform,
			"platform_family":       hostInfo.PlatformFamily,
			"platform_version":      hostInfo.PlatformVersion,
			"kernel_version":        hostInfo.KernelVersion,
			"kernel_arch":           hostInfo.KernelArch,
			"virtualization_system": hostInfo.VirtualizationSystem,
			"virtualization_role":   hostInfo.VirtualizationRole,
			"host_id":               hostInfo.HostID,
		})
	}

	// Collect detailed CPU information
	if cpuInfo, err := machinecpu.Info(); err == nil && len(cpuInfo) > 0 {
		c.createFactWithValue(collection, "cpu_detailed", cpuInfo)
	}

	// Collect CPU times
	if cpuTimes, err := machinecpu.Times(false); err == nil {
		c.createFactWithValue(collection, "cpu_times", cpuTimes)
	}

	// Collect CPU percentages
	if cpuPercent, err := machinecpu.Percent(0, false); err == nil {
		c.createFactWithValue(collection, "cpu_percent", cpuPercent)
	}

	// Collect detailed memory information
	if memInfo, err := machinemem.VirtualMemory(); err == nil {
		c.createFactWithValue(collection, "memory_detailed", map[string]interface{}{
			"total":            memInfo.Total,
			"available":        memInfo.Available,
			"used":             memInfo.Used,
			"used_percent":     memInfo.UsedPercent,
			"free":             memInfo.Free,
			"active":           memInfo.Active,
			"inactive":         memInfo.Inactive,
			"wired":            memInfo.Wired,
			"laundry":          memInfo.Laundry,
			"buffers":          memInfo.Buffers,
			"cached":           memInfo.Cached,
			"writeback":        memInfo.WriteBack,
			"dirty":            memInfo.Dirty,
			"writeback_tmp":    memInfo.WriteBackTmp,
			"shared":           memInfo.Shared,
			"slab":             memInfo.Slab,
			"sreclaimable":     memInfo.Sreclaimable,
			"sunreclaim":       memInfo.Sunreclaim,
			"page_tables":      memInfo.PageTables,
			"swap_cached":      memInfo.SwapCached,
			"commit_limit":     memInfo.CommitLimit,
			"committed_as":     memInfo.CommittedAS,
			"high_total":       memInfo.HighTotal,
			"high_free":        memInfo.HighFree,
			"low_total":        memInfo.LowTotal,
			"low_free":         memInfo.LowFree,
			"swap_total":       memInfo.SwapTotal,
			"swap_free":        memInfo.SwapFree,
			"mapped":           memInfo.Mapped,
			"vmalloc_total":    memInfo.VmallocTotal,
			"vmalloc_used":     memInfo.VmallocUsed,
			"vmalloc_chunk":    memInfo.VmallocChunk,
			"huge_pages_total": memInfo.HugePagesTotal,
			"huge_pages_free":  memInfo.HugePagesFree,
			"huge_page_size":   memInfo.HugePageSize,
		})
	}

	// Collect swap memory information
	if swapInfo, err := machinemem.SwapMemory(); err == nil {
		c.createFactWithValue(collection, "swap_memory", map[string]interface{}{
			"total":        swapInfo.Total,
			"used":         swapInfo.Used,
			"free":         swapInfo.Free,
			"used_percent": swapInfo.UsedPercent,
			"sin":          swapInfo.Sin,
			"sout":         swapInfo.Sout,
		})
	}

	// Collect disk partitions
	if partitions, err := machinedisk.Partitions(false); err == nil {
		c.createFactWithValue(collection, "disk_partitions", partitions)
	}

	// Collect disk I/O statistics
	if diskIO, err := machinedisk.IOCounters(); err == nil {
		c.createFactWithValue(collection, "disk_io", diskIO)
	}

	// Collect network interface statistics
	if netIO, err := machinenet.IOCounters(false); err == nil {
		c.createFactWithValue(collection, "network_io", netIO)
	}

	// Collect network connections
	if connections, err := machinenet.Connections("all"); err == nil {
		c.createFactWithValue(collection, "network_connections", connections)
	}

	// Collect network interfaces
	if interfaces, err := machinenet.Interfaces(); err == nil {
		c.createFactWithValue(collection, "network_interfaces", interfaces)
	}

	// Collect process information (top processes)
	if processes, err := machineprocess.Processes(); err == nil {
		// Limit to top 20 processes by CPU usage to avoid overwhelming the system
		processInfos := make([]map[string]interface{}, 0, 20)
		for i, proc := range processes {
			if i >= 20 {
				break
			}

			if name, err := proc.Name(); err == nil {
				if cpuPercent, err := proc.CPUPercent(); err == nil {
					if memPercent, err := proc.MemoryPercent(); err == nil {
						processInfos = append(processInfos, map[string]interface{}{
							"pid":            proc.Pid,
							"name":           name,
							"cpu_percent":    cpuPercent,
							"memory_percent": memPercent,
						})
					}
				}
			}
		}
		c.createFactWithValue(collection, "top_processes", processInfos)
	}

	// Collect sensors information (temperature, fans, etc.)
	if sensors, err := machinesensors.SensorsTemperatures(); err == nil {
		c.createFactWithValue(collection, "sensors_temperature", sensors)
	}

	// Collect user information
	if users, err := machinehost.Users(); err == nil {
		c.createFactWithValue(collection, "users", users)
	}

	return nil
}

// collectSpecificFact collects a specific fact by key
func (c *Collector) collectSpecificFact(collection *spookyfactstypes.FactCollection, key string) error {
	factCollectors := map[string]func(*spookyfactstypes.FactCollection) error{
		"hostname":         c.collectHostname,
		"machine_id":       c.collectMachineID,
		"fqdn":             c.collectFQDN,
		"os_name":          c.collectOSName,
		"os_version":       c.collectOSVersion,
		"os_distro":        c.collectOSDistro,
		"os_arch":          c.collectOSArch,
		"kernel_version":   c.collectKernelVersion,
		"cpu_cores":        c.collectCPUCores,
		"cpu_model":        c.collectCPUModel,
		"memory_total":     c.collectMemoryTotal,
		"memory_used":      c.collectMemoryUsed,
		"memory_available": c.collectMemoryAvailable,
		"disk_total":       c.collectDiskTotal,
		"disk_used":        c.collectDiskUsed,
		"disk_available":   c.collectDiskAvailable,
		"ip_addresses":     c.collectIPAddresses,
		"mac_addresses":    c.collectMACAddresses,
		"dns_servers":      c.collectDNSServers,
		"dns_search":       c.collectDNSSearch,
		"environment":      c.collectEnvironment,
	}

	if collector, exists := factCollectors[key]; exists {
		return collector(collection)
	}

	return fmt.Errorf("unknown fact key: %s", key)
}

// Helper methods for specific fact collection
func (c *Collector) collectHostname(collection *spookyfactstypes.FactCollection) error {
	if hostname, err := os.Hostname(); err == nil {
		c.createFact(collection, "hostname", hostname)
	}
	return nil
}

func (c *Collector) collectMachineID(collection *spookyfactstypes.FactCollection) error {
	if machineID, err := c.getMachineID(); err == nil {
		c.createFact(collection, "machine_id", machineID)
	}
	return nil
}

func (c *Collector) collectFQDN(collection *spookyfactstypes.FactCollection) error {
	if fqdn, err := c.getFQDN(); err == nil {
		c.createFact(collection, "fqdn", fqdn)
	}
	return nil
}

func (c *Collector) collectOSName(collection *spookyfactstypes.FactCollection) error {
	c.createFact(collection, "os_name", runtime.GOOS)
	return nil
}

func (c *Collector) collectOSVersion(collection *spookyfactstypes.FactCollection) error {
	if version, err := c.getOSVersion(); err == nil {
		c.createFact(collection, "os_version", version)
	}
	return nil
}

func (c *Collector) collectOSDistro(collection *spookyfactstypes.FactCollection) error {
	if distro, err := c.getOSDistro(); err == nil {
		c.createFact(collection, "os_distro", distro)
	}
	return nil
}

func (c *Collector) collectOSArch(collection *spookyfactstypes.FactCollection) error {
	c.createFact(collection, "os_arch", runtime.GOARCH)
	return nil
}

func (c *Collector) collectKernelVersion(collection *spookyfactstypes.FactCollection) error {
	if kernel, err := c.getKernelVersion(); err == nil {
		c.createFact(collection, "kernel_version", kernel)
	}
	return nil
}

func (c *Collector) collectCPUCores(collection *spookyfactstypes.FactCollection) error {
	// Use gopsutil for more accurate CPU information
	if cpuInfo, err := machinecpu.Info(); err == nil && len(cpuInfo) > 0 {
		// Calculate total cores across all CPUs
		totalCores := 0
		for i := range cpuInfo {
			totalCores += int(cpuInfo[i].Cores)
		}
		c.createFactWithValue(collection, "cpu_cores", totalCores)
	} else {
		// Fallback to runtime
		cores := runtime.NumCPU()
		c.createFactWithValue(collection, "cpu_cores", cores)
	}
	return nil
}

func (c *Collector) collectCPUModel(collection *spookyfactstypes.FactCollection) error {
	if model, err := c.getCPUModel(); err == nil {
		c.createFact(collection, "cpu_model", model)
	}
	return nil
}

func (c *Collector) collectMemoryTotal(collection *spookyfactstypes.FactCollection) error {
	// Use gopsutil for memory information
	if memInfo, err := machinemem.VirtualMemory(); err == nil {
		c.createFactWithValue(collection, "memory_total", memInfo.Total)
	} else {
		// Fallback to system command parsing
		if memInfo, err := c.getMemoryInfo(); err == nil {
			c.createFactWithValue(collection, "memory_total", memInfo.Total)
		}
	}
	return nil
}

func (c *Collector) collectMemoryUsed(collection *spookyfactstypes.FactCollection) error {
	// Use gopsutil for memory information
	if memInfo, err := machinemem.VirtualMemory(); err == nil {
		c.createFactWithValue(collection, "memory_used", memInfo.Used)
	} else {
		// Fallback to system command parsing
		if memInfo, err := c.getMemoryInfo(); err == nil {
			c.createFactWithValue(collection, "memory_used", memInfo.Used)
		}
	}
	return nil
}

func (c *Collector) collectMemoryAvailable(collection *spookyfactstypes.FactCollection) error {
	// Use gopsutil for memory information
	if memInfo, err := machinemem.VirtualMemory(); err == nil {
		c.createFactWithValue(collection, "memory_available", memInfo.Available)
	} else {
		// Fallback to system command parsing
		if memInfo, err := c.getMemoryInfo(); err == nil {
			c.createFactWithValue(collection, "memory_available", memInfo.Available)
		}
	}
	return nil
}

func (c *Collector) collectDiskTotal(collection *spookyfactstypes.FactCollection) error {
	// Use gopsutil for disk information
	if diskInfo, err := machinedisk.Usage("/"); err == nil {
		c.createFactWithValue(collection, "disk_total", diskInfo.Total)
	} else {
		// Fallback to system command parsing
		if diskInfo, err := c.getDiskInfo(); err == nil {
			c.createFactWithValue(collection, "disk_total", diskInfo.Total)
		}
	}
	return nil
}

func (c *Collector) collectDiskUsed(collection *spookyfactstypes.FactCollection) error {
	// Use gopsutil for disk information
	if diskInfo, err := machinedisk.Usage("/"); err == nil {
		c.createFactWithValue(collection, "disk_used", diskInfo.Used)
	} else {
		// Fallback to system command parsing
		if diskInfo, err := c.getDiskInfo(); err == nil {
			c.createFactWithValue(collection, "disk_used", diskInfo.Used)
		}
	}
	return nil
}

func (c *Collector) collectDiskAvailable(collection *spookyfactstypes.FactCollection) error {
	// Use gopsutil for disk information
	if diskInfo, err := machinedisk.Usage("/"); err == nil {
		c.createFactWithValue(collection, "disk_available", diskInfo.Free)
	} else {
		// Fallback to system command parsing
		if diskInfo, err := c.getDiskInfo(); err == nil {
			c.createFactWithValue(collection, "disk_available", diskInfo.Available)
		}
	}
	return nil
}

func (c *Collector) collectIPAddresses(collection *spookyfactstypes.FactCollection) error {
	if ips, err := c.getIPAddresses(); err == nil {
		c.createFactWithValue(collection, "ip_addresses", ips)
	}
	return nil
}

func (c *Collector) collectMACAddresses(collection *spookyfactstypes.FactCollection) error {
	if macs, err := c.getMACAddresses(); err == nil {
		c.createFactWithValue(collection, "mac_addresses", macs)
	}
	return nil
}

func (c *Collector) collectDNSServers(collection *spookyfactstypes.FactCollection) error {
	if dns, err := c.getDNSConfig(); err == nil {
		c.createFactWithValue(collection, "dns_servers", dns.Nameservers)
	}
	return nil
}

func (c *Collector) collectDNSSearch(collection *spookyfactstypes.FactCollection) error {
	if dns, err := c.getDNSConfig(); err == nil {
		c.createFactWithValue(collection, "dns_search", dns.Search)
	}
	return nil
}

func (c *Collector) collectEnvironment(collection *spookyfactstypes.FactCollection) error {
	env := make(map[string]string)
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		if len(pair) == 2 {
			env[pair[0]] = pair[1]
		}
	}
	c.createFactWithValue(collection, "environment", env)
	return nil
}

// Platform-specific helper methods
func (c *Collector) getMachineID() (string, error) {
	// Try to read machine ID from /etc/machine-id
	if data, err := os.ReadFile("/etc/machine-id"); err == nil {
		return strings.TrimSpace(string(data)), nil
	}

	// Fallback to hostname
	return os.Hostname()
}

func (c *Collector) getFQDN() (string, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return "", err
	}

	// Try to get FQDN using hostname command
	if output, err := exec.Command("hostname", "-f").Output(); err == nil {
		return strings.TrimSpace(string(output)), nil
	}

	return hostname, nil
}

func (c *Collector) getOSVersion() (string, error) {
	switch runtime.GOOS {
	case "linux":
		// Try to read from /etc/os-release
		if data, err := os.ReadFile("/etc/os-release"); err == nil {
			lines := strings.Split(string(data), "\n")
			for _, line := range lines {
				if strings.HasPrefix(line, "VERSION=") {
					version := strings.TrimPrefix(line, "VERSION=")
					return strings.Trim(version, `"`), nil
				}
			}
		}
	case "darwin":
		if output, err := exec.Command("sw_vers", "-productVersion").Output(); err == nil {
			return strings.TrimSpace(string(output)), nil
		}
	}

	return runtime.GOOS, nil
}

func (c *Collector) getOSDistro() (string, error) {
	switch runtime.GOOS {
	case "linux":
		// Try to read from /etc/os-release
		if data, err := os.ReadFile("/etc/os-release"); err == nil {
			lines := strings.Split(string(data), "\n")
			for _, line := range lines {
				if strings.HasPrefix(line, "ID=") {
					distro := strings.TrimPrefix(line, "ID=")
					return strings.Trim(distro, `"`), nil
				}
			}
		}
	case "darwin":
		return "macos", nil
	}

	return runtime.GOOS, nil
}

func (c *Collector) getKernelVersion() (string, error) {
	switch runtime.GOOS {
	case "linux":
		if output, err := exec.Command("uname", "-r").Output(); err == nil {
			return strings.TrimSpace(string(output)), nil
		}
	case "darwin":
		if output, err := exec.Command("uname", "-r").Output(); err == nil {
			return strings.TrimSpace(string(output)), nil
		}
	}

	return "", fmt.Errorf("kernel version not available for %s", runtime.GOOS)
}

func (c *Collector) getCPUModel() (string, error) {
	switch runtime.GOOS {
	case "linux":
		if data, err := os.ReadFile("/proc/cpuinfo"); err == nil {
			lines := strings.Split(string(data), "\n")
			for _, line := range lines {
				if strings.HasPrefix(line, "model name") {
					parts := strings.Split(line, ":")
					if len(parts) >= 2 {
						return strings.TrimSpace(parts[1]), nil
					}
				}
			}
		}
	case "darwin":
		if output, err := exec.Command("sysctl", "-n", "machdep.cpu.brand_string").Output(); err == nil {
			return strings.TrimSpace(string(output)), nil
		}
	}

	return "unknown", nil
}

func (c *Collector) getMemoryInfo() (*spookyfactstypes.MemoryInfo, error) {
	switch runtime.GOOS {
	case "linux":
		if data, err := os.ReadFile("/proc/meminfo"); err == nil {
			return c.parseMemInfo(string(data)), nil
		}
	case "darwin":
		// Use vm_stat command on macOS
		if output, err := exec.Command("vm_stat").Output(); err == nil {
			return c.parseDarwinMemInfo(string(output)), nil
		}
	}

	return &spookyfactstypes.MemoryInfo{}, fmt.Errorf("memory info not available for %s", runtime.GOOS)
}

func (c *Collector) getDiskInfo() (*spookyfactstypes.DiskInfo, error) {
	switch runtime.GOOS {
	case "linux":
		if output, err := exec.Command("df", "-h", "/").Output(); err == nil {
			return c.parseDiskInfo(string(output)), nil
		}
	case "darwin":
		if output, err := exec.Command("df", "-h", "/").Output(); err == nil {
			return c.parseDiskInfo(string(output)), nil
		}
	}

	return &spookyfactstypes.DiskInfo{}, fmt.Errorf("disk info not available for %s", runtime.GOOS)
}

func (c *Collector) getIPAddresses() ([]string, error) {
	switch runtime.GOOS {
	case "linux":
		if output, err := exec.Command("hostname", "-I").Output(); err == nil {
			return c.parseIPAddresses(string(output)), nil
		}
	case "darwin":
		if output, err := exec.Command("ifconfig").Output(); err == nil {
			return c.parseDarwinIPAddresses(string(output)), nil
		}
	}

	return []string{}, fmt.Errorf("IP addresses not available for %s", runtime.GOOS)
}

func (c *Collector) getMACAddresses() ([]string, error) {
	switch runtime.GOOS {
	case "linux":
		if output, err := exec.Command("ip", "link", "show").Output(); err == nil {
			return c.parseMACAddresses(string(output)), nil
		}
	case "darwin":
		if output, err := exec.Command("ifconfig").Output(); err == nil {
			return c.parseDarwinMACAddresses(string(output)), nil
		}
	}

	return []string{}, fmt.Errorf("MAC addresses not available for %s", runtime.GOOS)
}

func (c *Collector) getDNSConfig() (*spookyfactstypes.DNSInfo, error) {
	switch runtime.GOOS {
	case "linux":
		if data, err := os.ReadFile("/etc/resolv.conf"); err == nil {
			return c.parseDNSConfig(string(data)), nil
		}
	case "darwin":
		if output, err := exec.Command("scutil", "--dns").Output(); err == nil {
			return c.parseDarwinDNSConfig(string(output)), nil
		}
	}

	return &spookyfactstypes.DNSInfo{}, fmt.Errorf("DNS config not available for %s", runtime.GOOS)
}

// Helper methods for creating facts
func (c *Collector) createFact(collection *spookyfactstypes.FactCollection, key, value string) {
	collection.Facts[key] = &spookyfactstypes.Fact{
		Key:       key,
		Value:     value,
		Source:    string(c.GetSource()),
		Timestamp: time.Now(),
	}
}

func (c *Collector) createFactWithValue(collection *spookyfactstypes.FactCollection, key string, value interface{}) {
	collection.Facts[key] = &spookyfactstypes.Fact{
		Key:       key,
		Value:     value,
		Source:    string(c.GetSource()),
		Timestamp: time.Now(),
	}
}

// Parsing helper methods
func (c *Collector) parseMemInfo(memInfo string) *spookyfactstypes.MemoryInfo {
	info := &spookyfactstypes.MemoryInfo{}
	lines := strings.Split(memInfo, "\n")

	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}

		switch fields[0] {
		case "MemTotal:":
			if val, err := strconv.ParseUint(fields[1], 10, 64); err == nil {
				info.Total = val * 1024 // Convert KB to bytes
			}
		case "MemAvailable:":
			if val, err := strconv.ParseUint(fields[1], 10, 64); err == nil {
				info.Available = val * 1024
			}
		case "MemFree:":
			if val, err := strconv.ParseUint(fields[1], 10, 64); err == nil {
				// Store free memory in Available if not already set
				if info.Available == 0 {
					info.Available = val * 1024
				}
			}
		}
	}

	info.Used = info.Total - info.Available
	return info
}

func (c *Collector) parseDarwinMemInfo(vmStat string) *spookyfactstypes.MemoryInfo {
	info := &spookyfactstypes.MemoryInfo{}
	lines := strings.Split(vmStat, "\n")

	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 3 {
			continue
		}

		if fields[0] == "Pages" && fields[1] == "free:" {
			if val, err := strconv.ParseUint(fields[2], 10, 64); err == nil {
				info.Available = val * 4096 // Convert pages to bytes
			}
		} else if fields[0] == "Pages" && fields[1] == "active:" {
			if val, err := strconv.ParseUint(fields[2], 10, 64); err == nil {
				info.Used += val * 4096
			}
		} else if fields[0] == "Pages" && fields[1] == "inactive:" {
			if val, err := strconv.ParseUint(fields[2], 10, 64); err == nil {
				info.Used += val * 4096
			}
		}
	}

	// Estimate total memory (this is approximate on macOS)
	info.Total = info.Used + info.Available
	return info
}

func (c *Collector) parseDiskInfo(dfOutput string) *spookyfactstypes.DiskInfo {
	info := &spookyfactstypes.DiskInfo{}
	lines := strings.Split(dfOutput, "\n")

	// Skip header line
	if len(lines) < 2 {
		return info
	}

	fields := strings.Fields(lines[1])
	if len(fields) < 4 {
		return info
	}

	// Parse size fields (in KB)
	if val, err := strconv.ParseUint(fields[1], 10, 64); err == nil {
		info.Total = val * 1024
	}
	if val, err := strconv.ParseUint(fields[2], 10, 64); err == nil {
		info.Used = val * 1024
	}
	if val, err := strconv.ParseUint(fields[3], 10, 64); err == nil {
		info.Available = val * 1024
	}

	return info
}

func (c *Collector) parseIPAddresses(output string) []string {
	var ips []string
	fields := strings.Fields(output)

	for _, field := range fields {
		// Basic IP validation
		if strings.Contains(field, ".") || strings.Contains(field, ":") {
			ips = append(ips, field)
		}
	}

	return ips
}

func (c *Collector) parseDarwinIPAddresses(ifconfigOutput string) []string {
	var ips []string
	lines := strings.Split(ifconfigOutput, "\n")

	for _, line := range lines {
		if strings.Contains(line, "inet ") {
			fields := strings.Fields(line)
			for i, field := range fields {
				if field == "inet" && i+1 < len(fields) {
					ip := fields[i+1]
					if ip != "127.0.0.1" && ip != "::1" {
						ips = append(ips, ip)
					}
				}
			}
		}
	}

	return ips
}

func (c *Collector) parseMACAddresses(output string) []string {
	var macs []string
	lines := strings.Split(output, "\n")

	for _, line := range lines {
		if strings.Contains(line, "link/ether") {
			fields := strings.Fields(line)
			for i, field := range fields {
				if field == "link/ether" && i+1 < len(fields) {
					macs = append(macs, fields[i+1])
				}
			}
		}
	}

	return macs
}

func (c *Collector) parseDarwinMACAddresses(ifconfigOutput string) []string {
	var macs []string
	lines := strings.Split(ifconfigOutput, "\n")

	for _, line := range lines {
		if strings.Contains(line, "ether ") {
			fields := strings.Fields(line)
			for i, field := range fields {
				if field == "ether" && i+1 < len(fields) {
					macs = append(macs, fields[i+1])
				}
			}
		}
	}

	return macs
}

func (c *Collector) parseDNSConfig(resolvConf string) *spookyfactstypes.DNSInfo {
	info := &spookyfactstypes.DNSInfo{}
	lines := strings.Split(resolvConf, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "nameserver ") {
			server := strings.TrimPrefix(line, "nameserver ")
			info.Nameservers = append(info.Nameservers, server)
		} else if strings.HasPrefix(line, "search ") {
			search := strings.TrimPrefix(line, "search ")
			info.Search = append(info.Search, search)
		}
	}

	return info
}

func (c *Collector) parseDarwinDNSConfig(scutilOutput string) *spookyfactstypes.DNSInfo {
	info := &spookyfactstypes.DNSInfo{}
	lines := strings.Split(scutilOutput, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "nameserver[") {
			// Extract nameserver from line like "nameserver[0] : 8.8.8.8"
			if idx := strings.Index(line, ":"); idx != -1 {
				server := strings.TrimSpace(line[idx+1:])
				info.Nameservers = append(info.Nameservers, server)
			}
		}
	}

	return info
}

// collectSpookyMinimalFacts collects Spooky's minimal facts for local system
func (c *Collector) collectSpookyMinimalFacts(collection *spookyfactstypes.FactCollection) error {
	// Detect platform first
	platform, err := c.detectPlatform()
	if err != nil {
		return fmt.Errorf("failed to detect platform: %w", err)
	}

	// Collect system facts
	if err := c.collectSpookySystemFacts(collection, platform); err != nil {
		return fmt.Errorf("failed to collect Spooky system facts: %w", err)
	}

	// Collect OS facts
	if err := c.collectSpookyOSFacts(collection, platform); err != nil {
		return fmt.Errorf("failed to collect Spooky OS facts: %w", err)
	}

	// Collect hardware facts
	if err := c.collectSpookyHardwareFacts(collection, platform); err != nil {
		return fmt.Errorf("failed to collect Spooky hardware facts: %w", err)
	}

	// Collect network facts
	if err := c.collectSpookyNetworkFacts(collection, platform); err != nil {
		return fmt.Errorf("failed to collect Spooky network facts: %w", err)
	}

	// Collect user facts
	if err := c.collectSpookyUserFacts(collection, platform); err != nil {
		return fmt.Errorf("failed to collect Spooky user facts: %w", err)
	}

	// Collect environment facts
	if err := c.collectSpookyEnvironmentFacts(collection, platform); err != nil {
		return fmt.Errorf("failed to collect Spooky environment facts: %w", err)
	}

	return nil
}

// Platform represents the detected operating system
type Platform struct {
	OS      string // "linux", "darwin", "windows", "freebsd", "openbsd", "netbsd"
	Arch    string // "amd64", "arm64", "386", etc.
	Version string // OS version
}

// detectPlatform detects the operating system and architecture
func (c *Collector) detectPlatform() (*Platform, error) {
	os := runtime.GOOS
	arch := runtime.GOARCH

	// Normalize OS names
	switch os {
	case "linux":
		os = "linux"
	case "darwin":
		os = "darwin"
	case "windows":
		os = "windows"
	case "freebsd":
		os = "freebsd"
	case "openbsd":
		os = "openbsd"
	case "netbsd":
		os = "netbsd"
	default:
		os = "unknown"
	}

	// Get version info
	var version string
	switch os {
	case "linux":
		version, _ = c.getKernelVersion()
	case "darwin":
		version, _ = c.getOSVersion()
	case "windows":
		version = "Windows"
	default:
		version = "unknown"
	}

	return &Platform{
		OS:      os,
		Arch:    arch,
		Version: version,
	}, nil
}

// collectSpookySystemFacts collects basic system information
func (c *Collector) collectSpookySystemFacts(collection *spookyfactstypes.FactCollection, platform *Platform) error {
	// System type
	c.createFact(collection, "spooky_system", platform.OS)

	// Architecture
	c.createFact(collection, "spooky_architecture", platform.Arch)
	c.createFact(collection, "spooky_machine", platform.Arch)

	// Kernel version
	c.createFact(collection, "spooky_kernel", platform.Version)

	// Hostname
	hostname, err := c.getHostname()
	if err == nil {
		c.createFact(collection, "spooky_hostname", hostname)
	}

	// FQDN
	fqdn, err := c.getFQDN()
	if err == nil {
		c.createFact(collection, "spooky_fqdn", fqdn)
	}

	// Domain
	domain, err := c.getDomain()
	if err == nil {
		c.createFact(collection, "spooky_domain", domain)
	}

	return nil
}

// collectSpookyOSFacts collects operating system information
func (c *Collector) collectSpookyOSFacts(collection *spookyfactstypes.FactCollection, platform *Platform) error {
	switch platform.OS {
	case "linux":
		return c.collectSpookyLinuxOSFacts(collection)
	case "darwin":
		return c.collectSpookyDarwinOSFacts(collection)
	case "windows":
		return c.collectSpookyWindowsOSFacts(collection)
	case "freebsd", "openbsd", "netbsd":
		return c.collectSpookyBSDOSFacts(collection, platform.OS)
	default:
		return fmt.Errorf("unsupported OS: %s", platform.OS)
	}
}

// collectSpookyHardwareFacts collects hardware information
func (c *Collector) collectSpookyHardwareFacts(collection *spookyfactstypes.FactCollection, platform *Platform) error {
	switch platform.OS {
	case "linux":
		return c.collectSpookyLinuxHardwareFacts(collection)
	case "darwin":
		return c.collectSpookyDarwinHardwareFacts(collection)
	case "windows":
		return c.collectSpookyWindowsHardwareFacts(collection)
	case "freebsd", "openbsd", "netbsd":
		return c.collectSpookyBSDHardwareFacts(collection)
	default:
		return fmt.Errorf("unsupported OS: %s", platform.OS)
	}
}

// collectSpookyNetworkFacts collects network information
func (c *Collector) collectSpookyNetworkFacts(collection *spookyfactstypes.FactCollection, platform *Platform) error {
	switch platform.OS {
	case "linux":
		return c.collectSpookyLinuxNetworkFacts(collection)
	case "darwin":
		return c.collectSpookyDarwinNetworkFacts(collection)
	case "windows":
		return c.collectSpookyWindowsNetworkFacts(collection)
	case "freebsd", "openbsd", "netbsd":
		return c.collectSpookyBSDNetworkFacts(collection)
	default:
		return fmt.Errorf("unsupported OS: %s", platform.OS)
	}
}

// collectSpookyUserFacts collects user information
func (c *Collector) collectSpookyUserFacts(collection *spookyfactstypes.FactCollection, platform *Platform) error {
	// User ID
	userID := os.Getenv("USER")
	if userID == "" {
		userID = os.Getenv("USERNAME") // Windows fallback
	}
	if userID != "" {
		c.createFact(collection, "spooky_user_id", userID)
	}

	// User directory
	userDir := os.Getenv("HOME")
	if userDir == "" {
		userDir = os.Getenv("USERPROFILE") // Windows fallback
	}
	if userDir != "" {
		c.createFact(collection, "spooky_user_dir", userDir)
	}

	// User shell
	userShell := os.Getenv("SHELL")
	if userShell == "" {
		userShell = os.Getenv("ComSpec") // Windows fallback
	}
	if userShell != "" {
		c.createFact(collection, "spooky_user_shell", userShell)
	}

	return nil
}

// collectSpookyEnvironmentFacts collects environment information
func (c *Collector) collectSpookyEnvironmentFacts(collection *spookyfactstypes.FactCollection, platform *Platform) error {
	// Environment variables
	env := make(map[string]string)
	for _, envVar := range os.Environ() {
		if idx := strings.Index(envVar, "="); idx > 0 {
			key := envVar[:idx]
			value := envVar[idx+1:]
			env[key] = value
		}
	}
	c.createFactWithValue(collection, "spooky_env", env)

	return nil
}

// Platform-specific Ansible fact collection methods

func (c *Collector) collectSpookyLinuxOSFacts(collection *spookyfactstypes.FactCollection) error {
	// Read /etc/os-release
	if data, err := os.ReadFile("/etc/os-release"); err == nil {
		osInfo := c.parseOSRelease(string(data))
		c.createFact(collection, "spooky_os_name", osInfo.Name)
		c.createFact(collection, "spooky_os_version", osInfo.Version)
		c.createFact(collection, "spooky_os_family", osInfo.Distribution)
		c.createFact(collection, "spooky_distribution", osInfo.Distribution)
		c.createFact(collection, "spooky_distribution_version", osInfo.Version)
	}

	return nil
}

func (c *Collector) collectSpookyDarwinOSFacts(collection *spookyfactstypes.FactCollection) error {
	// OS name
	c.createFact(collection, "spooky_os_name", "Darwin")
	c.createFact(collection, "spooky_os_family", "Darwin")

	// OS version
	if version, err := c.getOSVersion(); err == nil {
		c.createFact(collection, "spooky_os_version", version)
		c.createFact(collection, "spooky_distribution_version", version)
	}

	// Distribution
	c.createFact(collection, "spooky_distribution", "MacOSX")

	return nil
}

func (c *Collector) collectSpookyWindowsOSFacts(collection *spookyfactstypes.FactCollection) error {
	// OS name
	c.createFact(collection, "spooky_os_name", "Windows")
	c.createFact(collection, "spooky_os_family", "Windows")

	// For Windows, we'll use runtime info since we're local
	c.createFact(collection, "spooky_distribution", "Windows")
	c.createFact(collection, "spooky_distribution_version", "Windows")

	return nil
}

func (c *Collector) collectSpookyBSDOSFacts(collection *spookyfactstypes.FactCollection, bsdType string) error {
	c.createFact(collection, "spooky_os_name", bsdType)
	c.createFact(collection, "spooky_os_family", "BSD")
	c.createFact(collection, "spooky_distribution", bsdType)

	// Version
	if version, err := c.getKernelVersion(); err == nil {
		c.createFact(collection, "spooky_distribution_version", version)
	}

	return nil
}

func (c *Collector) collectSpookyLinuxHardwareFacts(collection *spookyfactstypes.FactCollection) error {
	// CPU info
	if cpuModel, err := c.getCPUModel(); err == nil {
		c.createFactWithValue(collection, "spooky_processor", []string{cpuModel})
	}

	// CPU cores
	if cores, err := machinecpu.Counts(false); err == nil {
		c.createFactWithValue(collection, "spooky_processor_cores", cores)
		c.createFactWithValue(collection, "spooky_processor_vcpus", cores)
	}

	// Memory info
	if memInfo, err := machinemem.VirtualMemory(); err == nil {
		c.createFactWithValue(collection, "spooky_memtotal_mb", memInfo.Total/1024/1024)
	}

	return nil
}

func (c *Collector) collectSpookyDarwinHardwareFacts(collection *spookyfactstypes.FactCollection) error {
	// CPU info
	if cpuModel, err := c.getCPUModel(); err == nil {
		c.createFactWithValue(collection, "spooky_processor", []string{cpuModel})
	}

	// CPU cores
	if cores, err := machinecpu.Counts(false); err == nil {
		c.createFactWithValue(collection, "spooky_processor_cores", cores)
		c.createFactWithValue(collection, "spooky_processor_vcpus", cores)
	}

	// Memory info
	if memInfo, err := machinemem.VirtualMemory(); err == nil {
		c.createFactWithValue(collection, "spooky_memtotal_mb", memInfo.Total/1024/1024)
	}

	return nil
}

func (c *Collector) collectSpookyWindowsHardwareFacts(collection *spookyfactstypes.FactCollection) error {
	// CPU info
	if cpuModel, err := c.getCPUModel(); err == nil {
		c.createFactWithValue(collection, "spooky_processor", []string{cpuModel})
	}

	// CPU cores
	if cores, err := machinecpu.Counts(false); err == nil {
		c.createFactWithValue(collection, "spooky_processor_cores", cores)
		c.createFactWithValue(collection, "spooky_processor_vcpus", cores)
	}

	// Memory info
	if memInfo, err := machinemem.VirtualMemory(); err == nil {
		c.createFactWithValue(collection, "spooky_memtotal_mb", memInfo.Total/1024/1024)
	}

	return nil
}

func (c *Collector) collectSpookyBSDHardwareFacts(collection *spookyfactstypes.FactCollection) error {
	// CPU info
	if cpuModel, err := c.getCPUModel(); err == nil {
		c.createFactWithValue(collection, "spooky_processor", []string{cpuModel})
	}

	// CPU cores
	if cores, err := machinecpu.Counts(false); err == nil {
		c.createFactWithValue(collection, "spooky_processor_cores", cores)
		c.createFactWithValue(collection, "spooky_processor_vcpus", cores)
	}

	// Memory info
	if memInfo, err := machinemem.VirtualMemory(); err == nil {
		c.createFactWithValue(collection, "spooky_memtotal_mb", memInfo.Total/1024/1024)
	}

	return nil
}

func (c *Collector) collectSpookyLinuxNetworkFacts(collection *spookyfactstypes.FactCollection) error {
	// Get default IPv4
	if addrs, err := machinenet.Interfaces(); err == nil {
		for _, addr := range addrs {
			for _, ip := range addr.Addrs {
				if strings.Contains(ip.Addr, ".") && !strings.HasPrefix(ip.Addr, "127.") {
					c.createFactWithValue(collection, "spooky_default_ipv4", map[string]interface{}{
						"address":   strings.Split(ip.Addr, "/")[0],
						"interface": addr.Name,
					})
					break
				}
			}
		}
	}

	// Get interfaces
	if addrs, err := machinenet.Interfaces(); err == nil {
		var interfaceList []string
		for _, addr := range addrs {
			interfaceList = append(interfaceList, addr.Name)
		}
		c.createFactWithValue(collection, "spooky_interfaces", interfaceList)
	}

	return nil
}

func (c *Collector) collectSpookyDarwinNetworkFacts(collection *spookyfactstypes.FactCollection) error {
	// Get default IPv4
	if addrs, err := machinenet.Interfaces(); err == nil {
		for _, addr := range addrs {
			for _, ip := range addr.Addrs {
				if strings.Contains(ip.Addr, ".") && !strings.HasPrefix(ip.Addr, "127.") {
					c.createFactWithValue(collection, "spooky_default_ipv4", map[string]interface{}{
						"address":   strings.Split(ip.Addr, "/")[0],
						"interface": addr.Name,
					})
					break
				}
			}
		}
	}

	// Get interfaces
	if addrs, err := machinenet.Interfaces(); err == nil {
		var interfaceList []string
		for _, addr := range addrs {
			interfaceList = append(interfaceList, addr.Name)
		}
		c.createFactWithValue(collection, "spooky_interfaces", interfaceList)
	}

	return nil
}

func (c *Collector) collectSpookyWindowsNetworkFacts(collection *spookyfactstypes.FactCollection) error {
	// Get default IPv4
	if addrs, err := machinenet.Interfaces(); err == nil {
		for _, addr := range addrs {
			for _, ip := range addr.Addrs {
				if strings.Contains(ip.Addr, ".") && !strings.HasPrefix(ip.Addr, "127.") {
					c.createFactWithValue(collection, "spooky_default_ipv4", map[string]interface{}{
						"address":   strings.Split(ip.Addr, "/")[0],
						"interface": addr.Name,
					})
					break
				}
			}
		}
	}

	// Get interfaces
	if addrs, err := machinenet.Interfaces(); err == nil {
		var interfaceList []string
		for _, addr := range addrs {
			interfaceList = append(interfaceList, addr.Name)
		}
		c.createFactWithValue(collection, "spooky_interfaces", interfaceList)
	}

	return nil
}

func (c *Collector) collectSpookyBSDNetworkFacts(collection *spookyfactstypes.FactCollection) error {
	// Get default IPv4
	if addrs, err := machinenet.Interfaces(); err == nil {
		for _, addr := range addrs {
			for _, ip := range addr.Addrs {
				if strings.Contains(ip.Addr, ".") && !strings.HasPrefix(ip.Addr, "127.") {
					c.createFactWithValue(collection, "spooky_default_ipv4", map[string]interface{}{
						"address":   strings.Split(ip.Addr, "/")[0],
						"interface": addr.Name,
					})
					break
				}
			}
		}
	}

	// Get interfaces
	if addrs, err := machinenet.Interfaces(); err == nil {
		var interfaceList []string
		for _, addr := range addrs {
			interfaceList = append(interfaceList, addr.Name)
		}
		c.createFactWithValue(collection, "spooky_interfaces", interfaceList)
	}

	return nil
}

// Helper methods for Ansible facts

func (c *Collector) getHostname() (string, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return "", err
	}
	return hostname, nil
}

func (c *Collector) getDomain() (string, error) {
	fqdn, err := c.getFQDN()
	if err != nil {
		return "", err
	}

	parts := strings.Split(fqdn, ".")
	if len(parts) > 1 {
		return strings.Join(parts[1:], "."), nil
	}

	return "", nil
}

func (c *Collector) parseOSRelease(osRelease string) spookyfactstypes.OSInfo {
	lines := strings.Split(osRelease, "\n")
	osInfo := spookyfactstypes.OSInfo{}

	for _, line := range lines {
		if strings.HasPrefix(line, "NAME=") {
			osInfo.Name = strings.Trim(strings.TrimPrefix(line, "NAME="), "\"")
		} else if strings.HasPrefix(line, "VERSION=") {
			osInfo.Version = strings.Trim(strings.TrimPrefix(line, "VERSION="), "\"")
		} else if strings.HasPrefix(line, "ID=") {
			osInfo.Distribution = strings.Trim(strings.TrimPrefix(line, "ID="), "\"")
		}
	}

	return osInfo
}
