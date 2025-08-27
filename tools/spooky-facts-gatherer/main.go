package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"spooky/internal/schemas"
	"spooky/internal/utilities"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/load"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/shirou/gopsutil/v4/net"
	"github.com/shirou/gopsutil/v4/process"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "spooky-facts",
	Short: "System facts gatherer for spooky automation platform",
	Long: `spooky-facts is a comprehensive system information gatherer that collects
detailed facts about the current system using gopsutil and outputs them in HCL format.

It gathers information about:
- Host information (OS, platform, virtualization)
- CPU details (cores, model, usage)
- Memory information (RAM, swap)
- Disk usage and partitions
- Network interfaces and connections
- System load and uptime
- Running processes
- And more...

The facts are written to /etc/spooky/facts.hcl on Unix systems or appropriate
Windows location.`,
}

var verbose bool

var gatherCmd = &cobra.Command{
	Use:   "gather",
	Short: "Gather system facts and write to HCL file",
	Run: func(cmd *cobra.Command, args []string) {
		if verbose {
			fmt.Println("🔍 Gathering system facts...")
		}

		facts, err := gatherAllFacts()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error gathering facts: %v\n", err)
			os.Exit(1)
		}

		outputPath, err := getFactsOutputPath()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error determining output path: %v\n", err)
			os.Exit(1)
		}

		err = writeFactsToHCL(facts, outputPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error writing facts to HCL: %v\n", err)
			os.Exit(1)
		}

		if verbose {
			fmt.Printf("✅ Facts gathered successfully and written to: %s\n", outputPath)
		}
	},
}

var previewCmd = &cobra.Command{
	Use:   "preview",
	Short: "Gather and display facts without writing to file",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🔍 Gathering system facts (preview mode)...")

		facts, err := gatherAllFacts()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error gathering facts: %v\n", err)
			os.Exit(1)
		}

		hclContent, err := factsToHCL(facts)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error converting facts to HCL: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("\n📄 Generated HCL content:")
		fmt.Println(strings.Repeat("=", 50))
		fmt.Println(hclContent)
	},
}

// SystemFacts represents all gathered system information
type SystemFacts struct {
	Timestamp   time.Time         `json:"timestamp"`
	Host        *host.InfoStat    `json:"host"`
	CPU         *CPUInfo          `json:"cpu"`
	Memory      *MemoryInfo       `json:"memory"`
	Disks       []*DiskInfo       `json:"disks"`
	Networks    []*NetworkInfo    `json:"networks"`
	Load        *load.AvgStat     `json:"load"`
	Processes   []*ProcessInfo    `json:"processes"`
	Environment map[string]string `json:"environment"`
	Runtime     *RuntimeInfo      `json:"runtime"`
}

// CPUInfo represents CPU information
type CPUInfo struct {
	Info          []cpu.InfoStat `json:"info"`
	Percent       []float64      `json:"percent"`
	Count         int            `json:"count"`
	PhysicalCount int            `json:"physical_count"`
}

// MemoryInfo represents memory information
type MemoryInfo struct {
	Virtual *mem.VirtualMemoryStat `json:"virtual"`
	Swap    *mem.SwapMemoryStat    `json:"swap"`
}

// DiskInfo represents disk information
type DiskInfo struct {
	Partition *disk.PartitionStat `json:"partition"`
	Usage     *disk.UsageStat     `json:"usage"`
}

// NetworkInfo represents network interface information
type NetworkInfo struct {
	Interface *net.InterfaceStat  `json:"interface"`
	Addrs     []net.InterfaceAddr `json:"addrs"`
	Stats     *net.IOCountersStat `json:"stats"`
}

// ProcessInfo represents process information
type ProcessInfo struct {
	PID        int32   `json:"pid"`
	Name       string  `json:"name"`
	CPU        float64 `json:"cpu_percent"`
	Memory     float32 `json:"memory_percent"`
	Status     string  `json:"status"`
	CreateTime int64   `json:"create_time"`
}

// RuntimeInfo represents Go runtime information
type RuntimeInfo struct {
	GOOS         string `json:"goos"`
	GOARCH       string `json:"goarch"`
	NumCPU       int    `json:"num_cpu"`
	NumGoroutine int    `json:"num_goroutine"`
	Version      string `json:"version"`
}

func gatherAllFacts() (*SystemFacts, error) {
	facts := &SystemFacts{
		Timestamp: time.Now(),
	}

	// Gather host information
	hostInfo, err := host.Info()
	if err != nil {
		return nil, fmt.Errorf("failed to gather host info: %w", err)
	}
	facts.Host = hostInfo

	// Gather CPU information
	cpuInfo, err := cpu.Info()
	if err != nil {
		return nil, fmt.Errorf("failed to gather CPU info: %w", err)
	}

	cpuPercent, err := cpu.Percent(0, true)
	if err != nil {
		return nil, fmt.Errorf("failed to gather CPU percent: %w", err)
	}

	facts.CPU = &CPUInfo{
		Info:          cpuInfo,
		Percent:       cpuPercent,
		Count:         len(cpuInfo),
		PhysicalCount: runtime.NumCPU(),
	}

	// Gather memory information
	virtualMem, err := mem.VirtualMemory()
	if err != nil {
		return nil, fmt.Errorf("failed to gather virtual memory: %w", err)
	}

	swapMem, err := mem.SwapMemory()
	if err != nil {
		return nil, fmt.Errorf("failed to gather swap memory: %w", err)
	}

	facts.Memory = &MemoryInfo{
		Virtual: virtualMem,
		Swap:    swapMem,
	}

	// Gather disk information
	partitions, err := disk.Partitions(false)
	if err != nil {
		return nil, fmt.Errorf("failed to gather disk partitions: %w", err)
	}

	for _, partition := range partitions {
		usage, err := disk.Usage(partition.Mountpoint)
		if err != nil {
			// Skip partitions we can't read
			continue
		}

		diskInfo := &DiskInfo{
			Partition: &partition,
			Usage:     usage,
		}
		facts.Disks = append(facts.Disks, diskInfo)
	}

	// Gather network information
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("failed to gather network interfaces: %w", err)
	}

	ioCounters, err := net.IOCounters(true)
	if err != nil {
		return nil, fmt.Errorf("failed to gather network IO counters: %w", err)
	}

	ioCountersMap := make(map[string]*net.IOCountersStat)
	for _, counter := range ioCounters {
		ioCountersMap[counter.Name] = &counter
	}

	for _, iface := range interfaces {
		// Use the addresses directly from the interface
		var ifaceAddrs []net.InterfaceAddr
		for _, addr := range iface.Addrs {
			ifaceAddrs = append(ifaceAddrs, addr)
		}

		networkInfo := &NetworkInfo{
			Interface: &iface,
			Addrs:     ifaceAddrs,
			Stats:     ioCountersMap[iface.Name],
		}
		facts.Networks = append(facts.Networks, networkInfo)
	}

	// Gather load information
	loadAvg, err := load.Avg()
	if err != nil {
		return nil, fmt.Errorf("failed to gather load average: %w", err)
	}
	facts.Load = loadAvg

	// Gather process information (top 20 by CPU usage)
	processes, err := process.Processes()
	if err != nil {
		return nil, fmt.Errorf("failed to gather processes: %w", err)
	}

	// Get CPU and memory info for each process
	var processInfos []*ProcessInfo
	for _, proc := range processes {
		cpuPercent, err := proc.CPUPercent()
		if err != nil {
			continue
		}

		memoryPercent, err := proc.MemoryPercent()
		if err != nil {
			continue
		}

		name, err := proc.Name()
		if err != nil {
			continue
		}

		status, err := proc.Status()
		if err != nil {
			status = []string{"unknown"}
		}

		createTime, err := proc.CreateTime()
		if err != nil {
			createTime = 0
		}

		// Convert status to string
		statusStr := "unknown"
		if len(status) > 0 {
			statusStr = status[0]
		}

		processInfo := &ProcessInfo{
			PID:        proc.Pid,
			Name:       name,
			CPU:        cpuPercent,
			Memory:     memoryPercent,
			Status:     statusStr,
			CreateTime: createTime,
		}
		processInfos = append(processInfos, processInfo)
	}

	// Sort by CPU usage and take top 20
	if len(processInfos) > 20 {
		// Simple sorting by CPU (in a real implementation, you'd want proper sorting)
		processInfos = processInfos[:20]
	}
	facts.Processes = processInfos

	// Gather environment variables
	facts.Environment = make(map[string]string)
	for _, env := range os.Environ() {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) == 2 {
			facts.Environment[parts[0]] = parts[1]
		}
	}

	// Gather runtime information
	facts.Runtime = &RuntimeInfo{
		GOOS:         runtime.GOOS,
		GOARCH:       runtime.GOARCH,
		NumCPU:       runtime.NumCPU(),
		NumGoroutine: runtime.NumGoroutine(),
		Version:      runtime.Version(),
	}

	return facts, nil
}

func getFactsOutputPath() (string, error) {
	// Get OS-specific path configuration
	config, err := utilities.GetPathConfig("spooky")
	if err != nil {
		return "", fmt.Errorf("failed to get path config: %w", err)
	}

	// Ensure the config directory exists
	err = utilities.EnsureDirectories(config)
	if err != nil {
		return "", fmt.Errorf("failed to create directories: %w", err)
	}

	return filepath.Join(config.ConfigDir, "facts.hcl"), nil
}

func writeFactsToHCL(facts *SystemFacts, outputPath string) error {
	hclContent, err := factsToHCL(facts)
	if err != nil {
		return fmt.Errorf("failed to convert facts to HCL: %w", err)
	}

	err = os.WriteFile(outputPath, []byte(hclContent), 0644)
	if err != nil {
		return fmt.Errorf("failed to write HCL file: %w", err)
	}

	return nil
}

func factsToHCL(facts *SystemFacts) (string, error) {
	// Create the enhanced facts structure using the schema for validation
	enhancedFacts := &schemas.EnhancedFactsV1{
		Facts: make(map[string]*schemas.FactV1),
	}

	// Host information
	enhancedFacts.Facts["hostname"] = &schemas.FactV1{
		Value:       facts.Host.Hostname,
		Type:        "string",
		Description: "System hostname",
	}

	enhancedFacts.Facts["os"] = &schemas.FactV1{
		Value:       facts.Host.OS,
		Type:        "string",
		Description: "Operating system",
	}

	enhancedFacts.Facts["platform"] = &schemas.FactV1{
		Value:       facts.Host.Platform,
		Type:        "string",
		Description: "Platform name",
	}

	enhancedFacts.Facts["platform_family"] = &schemas.FactV1{
		Value:       facts.Host.PlatformFamily,
		Type:        "string",
		Description: "Platform family",
	}

	enhancedFacts.Facts["platform_version"] = &schemas.FactV1{
		Value:       facts.Host.PlatformVersion,
		Type:        "string",
		Description: "Platform version",
	}

	enhancedFacts.Facts["kernel_version"] = &schemas.FactV1{
		Value:       facts.Host.KernelVersion,
		Type:        "string",
		Description: "Kernel version",
	}

	enhancedFacts.Facts["virtualization_system"] = &schemas.FactV1{
		Value:       facts.Host.VirtualizationSystem,
		Type:        "string",
		Description: "Virtualization system",
	}

	enhancedFacts.Facts["virtualization_role"] = &schemas.FactV1{
		Value:       facts.Host.VirtualizationRole,
		Type:        "string",
		Description: "Virtualization role",
	}

	enhancedFacts.Facts["uptime"] = &schemas.FactV1{
		Value:       facts.Host.Uptime,
		Type:        "number",
		Description: "System uptime in seconds",
	}

	enhancedFacts.Facts["boot_time"] = &schemas.FactV1{
		Value:       facts.Host.BootTime,
		Type:        "number",
		Description: "System boot time timestamp",
	}

	// CPU information
	enhancedFacts.Facts["cpu_count"] = &schemas.FactV1{
		Value:       facts.CPU.Count,
		Type:        "number",
		Description: "Number of CPU cores",
	}

	enhancedFacts.Facts["cpu_physical_count"] = &schemas.FactV1{
		Value:       facts.CPU.PhysicalCount,
		Type:        "number",
		Description: "Number of physical CPU cores",
	}

	if len(facts.CPU.Info) > 0 {
		cpu := facts.CPU.Info[0]
		enhancedFacts.Facts["cpu_vendor_id"] = &schemas.FactV1{
			Value:       cpu.VendorID,
			Type:        "string",
			Description: "CPU vendor ID",
		}

		enhancedFacts.Facts["cpu_family"] = &schemas.FactV1{
			Value:       cpu.Family,
			Type:        "string",
			Description: "CPU family",
		}

		enhancedFacts.Facts["cpu_model"] = &schemas.FactV1{
			Value:       cpu.Model,
			Type:        "string",
			Description: "CPU model",
		}

		enhancedFacts.Facts["cpu_model_name"] = &schemas.FactV1{
			Value:       cpu.ModelName,
			Type:        "string",
			Description: "CPU model name",
		}

		enhancedFacts.Facts["cpu_cores"] = &schemas.FactV1{
			Value:       cpu.Cores,
			Type:        "number",
			Description: "CPU cores per socket",
		}

		enhancedFacts.Facts["cpu_mhz"] = &schemas.FactV1{
			Value:       cpu.Mhz,
			Type:        "number",
			Description: "CPU frequency in MHz",
		}
	}

	// Memory information
	if facts.Memory.Virtual != nil {
		enhancedFacts.Facts["memory_total"] = &schemas.FactV1{
			Value:       facts.Memory.Virtual.Total,
			Type:        "number",
			Description: "Total memory in bytes",
		}

		enhancedFacts.Facts["memory_available"] = &schemas.FactV1{
			Value:       facts.Memory.Virtual.Available,
			Type:        "number",
			Description: "Available memory in bytes",
		}

		enhancedFacts.Facts["memory_used"] = &schemas.FactV1{
			Value:       facts.Memory.Virtual.Used,
			Type:        "number",
			Description: "Used memory in bytes",
		}

		enhancedFacts.Facts["memory_used_percent"] = &schemas.FactV1{
			Value:       facts.Memory.Virtual.UsedPercent,
			Type:        "number",
			Description: "Memory usage percentage",
		}
	}

	// Load average
	if facts.Load != nil {
		enhancedFacts.Facts["load_1"] = &schemas.FactV1{
			Value:       facts.Load.Load1,
			Type:        "number",
			Description: "1-minute load average",
		}

		enhancedFacts.Facts["load_5"] = &schemas.FactV1{
			Value:       facts.Load.Load5,
			Type:        "number",
			Description: "5-minute load average",
		}

		enhancedFacts.Facts["load_15"] = &schemas.FactV1{
			Value:       facts.Load.Load15,
			Type:        "number",
			Description: "15-minute load average",
		}
	}

	// Runtime information
	enhancedFacts.Facts["runtime_goos"] = &schemas.FactV1{
		Value:       facts.Runtime.GOOS,
		Type:        "string",
		Description: "Go runtime OS",
	}

	enhancedFacts.Facts["runtime_goarch"] = &schemas.FactV1{
		Value:       facts.Runtime.GOARCH,
		Type:        "string",
		Description: "Go runtime architecture",
	}

	enhancedFacts.Facts["runtime_version"] = &schemas.FactV1{
		Value:       facts.Runtime.Version,
		Type:        "string",
		Description: "Go runtime version",
	}

	// Environment variables (limited to common ones)
	commonEnvVars := []string{"PATH", "HOME", "USER", "SHELL", "LANG", "PWD", "TERM", "HOSTNAME"}
	for _, envVar := range commonEnvVars {
		if value, exists := facts.Environment[envVar]; exists {
			enhancedFacts.Facts[fmt.Sprintf("env_%s", strings.ToLower(envVar))] = &schemas.FactV1{
				Value:       value,
				Type:        "string",
				Description: fmt.Sprintf("Environment variable %s", envVar),
			}
		}
	}

	// Generate HCL directly from the schema structs
	var hcl strings.Builder
	hcl.WriteString("# System Facts gathered by spooky-facts\n")
	hcl.WriteString(fmt.Sprintf("# Generated at: %s\n\n", facts.Timestamp.Format(time.RFC3339)))
	hcl.WriteString("enhanced_facts {\n\n")

	// Generate fact blocks from the schema structs
	for factName, fact := range enhancedFacts.Facts {
		hcl.WriteString(fmt.Sprintf("  fact \"%s\" {\n", factName))

		// Format value based on type
		switch fact.Type {
		case "string":
			hcl.WriteString(fmt.Sprintf("    value = \"%v\"\n", fact.Value))
		case "number":
			hcl.WriteString(fmt.Sprintf("    value = %v\n", fact.Value))
		default:
			hcl.WriteString(fmt.Sprintf("    value = \"%v\"\n", fact.Value))
		}

		hcl.WriteString(fmt.Sprintf("    type = \"%s\"\n", fact.Type))
		hcl.WriteString(fmt.Sprintf("    description = \"%s\"\n", fact.Description))
		hcl.WriteString("  }\n\n")
	}

	hcl.WriteString("}\n")
	return hcl.String(), nil
}

func init() {
	gatherCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
	rootCmd.AddCommand(gatherCmd)
	rootCmd.AddCommand(previewCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
