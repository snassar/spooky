// Package facts provides fact collection and in-memory management functionality.
package facts

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	spookyinterfaces "spooky/internal/interfaces"
	spookytypes "spooky/internal/types"
	spookytypesfacts "spooky/internal/types/facts"
	spookytypesssh "spooky/internal/types/ssh"
)

// SystemFactCollector collects system facts using SSH commands
type SystemFactCollector struct {
	name       string
	sshManager spookyinterfaces.SSHManager
}

// NewSystemFactCollector creates a new system fact collector
func NewSystemFactCollector(sshManager spookyinterfaces.SSHManager) *SystemFactCollector {
	return &SystemFactCollector{
		name:       "system",
		sshManager: sshManager,
	}
}

// GetName returns the collector name
func (c *SystemFactCollector) GetName() string {
	return c.name
}

// Collect collects facts from the given machine
func (c *SystemFactCollector) Collect(ctx context.Context, machine *spookytypes.Machine) (*spookytypesfacts.FactCollection, error) {
	// Get machine ID
	machineID, err := c.getMachineID(machine)
	if err != nil {
		return nil, fmt.Errorf("failed to get machine ID: %w", err)
	}

	// Collect system facts
	facts := &spookytypesfacts.Facts{
		System: &spookytypesfacts.SystemFacts{},
	}

	// Collect system facts via SSH
	systemFacts, err := c.collectSystemFacts(ctx, machine)
	if err != nil {
		return nil, fmt.Errorf("failed to collect system facts: %w", err)
	}
	facts.System = systemFacts

	// Collect collector facts via SSH (if available)
	collectorFacts, err := c.collectCollectorFacts(ctx, machine)
	if err != nil {
		// Collector facts are optional, log but don't fail
		fmt.Printf("Warning: failed to collect collector facts: %v\n", err)
	} else {
		facts.Collector = collectorFacts
	}

	// Collect custom facts via SSH (if available)
	customFacts, err := c.collectCustomFacts(ctx, machine)
	if err != nil {
		// Custom facts are optional, log but don't fail
		fmt.Printf("Warning: failed to collect custom facts: %v\n", err)
	} else {
		facts.Custom = customFacts
	}

	// Create fact collection
	collection := &spookytypesfacts.FactCollection{
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

// getMachineID gets the machine ID from /etc/machine-id via SSH
func (c *SystemFactCollector) getMachineID(machine *spookytypes.Machine) (string, error) {
	// Run command via SSH
	output, err := c.runSSHCommand(machine, "cat /etc/machine-id")
	if err != nil {
		return "", fmt.Errorf("failed to read /etc/machine-id via SSH: %w", err)
	}

	machineID := strings.TrimSpace(output)

	// Validate machine ID format (32-character hex string)
	if !regexp.MustCompile(`^[a-f0-9]{32}$`).MatchString(machineID) {
		return "", fmt.Errorf("invalid machine ID format: %s", machineID)
	}

	return machineID, nil
}

// collectSystemFacts collects operating system facts via SSH commands
func (c *SystemFactCollector) collectSystemFacts(_ context.Context, machine *spookytypes.Machine) (*spookytypesfacts.SystemFacts, error) {
	facts := &spookytypesfacts.SystemFacts{}

	// Collect OS information
	osFacts, err := c.collectOSFacts(machine)
	if err != nil {
		return nil, fmt.Errorf("failed to collect OS facts: %w", err)
	}
	facts.OS = osFacts

	// Collect hardware information
	hardwareFacts, err := c.collectHardwareFacts(machine)
	if err != nil {
		return nil, fmt.Errorf("failed to collect hardware facts: %w", err)
	}
	facts.Hardware = hardwareFacts

	// Collect network information
	networkFacts, err := c.collectNetworkFacts(machine)
	if err != nil {
		return nil, fmt.Errorf("failed to collect network facts: %w", err)
	}
	facts.Network = networkFacts

	// Collect load average
	loadFacts, err := c.collectLoadAverageFacts(machine)
	if err != nil {
		return nil, fmt.Errorf("failed to collect load average facts: %w", err)
	}
	facts.LoadAverage = loadFacts

	// Collect process information
	processFacts, err := c.collectProcessFacts(machine)
	if err != nil {
		return nil, fmt.Errorf("failed to collect process facts: %w", err)
	}
	facts.Processes = processFacts

	return facts, nil
}

// collectOSFacts collects operating system facts via SSH
func (c *SystemFactCollector) collectOSFacts(machine *spookytypes.Machine) (*spookytypesfacts.OSFacts, error) {
	// Get OS release information
	osRelease, err := c.runSSHCommand(machine, "cat /etc/os-release")
	if err != nil {
		return nil, fmt.Errorf("failed to get OS release info: %w", err)
	}

	// Parse OS release information
	osInfo := c.parseOSRelease(osRelease)

	// Get kernel information
	kernelInfo, err := c.runSSHCommand(machine, "uname -a")
	if err != nil {
		return nil, fmt.Errorf("failed to get kernel info: %w", err)
	}

	// Parse kernel information
	kernelParts := strings.Fields(kernelInfo)
	if len(kernelParts) < 3 {
		return nil, fmt.Errorf("invalid kernel info format: %s", kernelInfo)
	}

	return &spookytypesfacts.OSFacts{
		Name:     osInfo["NAME"],
		Version:  osInfo["VERSION"],
		Arch:     kernelParts[11], // Architecture from uname -a
		Kernel:   kernelParts[2],  // Kernel version from uname -a
		Platform: osInfo["ID"],
		Family:   osInfo["ID_LIKE"],
	}, nil
}

// collectHardwareFacts collects hardware facts via SSH
func (c *SystemFactCollector) collectHardwareFacts(machine *spookytypes.Machine) (*spookytypesfacts.HardwareFacts, error) {
	// Collect CPU facts
	cpuFacts, err := c.collectCPUFacts(machine)
	if err != nil {
		return nil, fmt.Errorf("failed to collect CPU facts: %w", err)
	}

	// Collect memory facts
	memoryFacts, err := c.collectMemoryFacts(machine)
	if err != nil {
		return nil, fmt.Errorf("failed to collect memory facts: %w", err)
	}

	// Collect disk facts
	diskFacts, err := c.collectDiskFacts(machine)
	if err != nil {
		return nil, fmt.Errorf("failed to collect disk facts: %w", err)
	}

	return &spookytypesfacts.HardwareFacts{
		CPU:    cpuFacts,
		Memory: memoryFacts,
		Disks:  diskFacts,
	}, nil
}

// collectCPUFacts collects CPU facts via SSH
func (c *SystemFactCollector) collectCPUFacts(machine *spookytypes.Machine) (*spookytypesfacts.CPUFacts, error) {
	// Get CPU info
	cpuInfo, err := c.runSSHCommand(machine, "cat /proc/cpuinfo")
	if err != nil {
		return nil, fmt.Errorf("failed to get CPU info: %w", err)
	}

	// Parse CPU information
	cpuDetails := c.parseCPUInfo(cpuInfo)

	// Get CPU count
	cpuCount, err := c.runSSHCommand(machine, "nproc")
	if err != nil {
		return nil, fmt.Errorf("failed to get CPU count: %w", err)
	}

	cores := 0
	if _, err := fmt.Sscanf(strings.TrimSpace(cpuCount), "%d", &cores); err != nil {
		// If parsing fails, default to 0 cores
		cores = 0
	}

	frequency, _ := strconv.ParseFloat(cpuDetails["cpu MHz"], 64)

	return &spookytypesfacts.CPUFacts{
		Cores:        cores,
		Model:        cpuDetails["model name"],
		Frequency:    frequency,
		Architecture: cpuDetails["cpu family"],
		Vendor:       cpuDetails["vendor_id"],
	}, nil
}

// collectMemoryFacts collects memory facts via SSH
func (c *SystemFactCollector) collectMemoryFacts(machine *spookytypes.Machine) (*spookytypesfacts.MemoryFacts, error) {
	// Get memory information
	memInfo, err := c.runSSHCommand(machine, "free -b")
	if err != nil {
		return nil, fmt.Errorf("failed to get memory info: %w", err)
	}

	// Parse memory information
	memDetails := c.parseMemoryInfo(memInfo)

	return &spookytypesfacts.MemoryFacts{
		Total:     memDetails["total"],
		Available: memDetails["available"],
		Used:      memDetails["used"],
		Free:      memDetails["free"],
	}, nil
}

// collectDiskFacts collects disk facts via SSH
func (c *SystemFactCollector) collectDiskFacts(machine *spookytypes.Machine) ([]*spookytypesfacts.DiskFacts, error) {
	// Get disk usage information
	diskInfo, err := c.runSSHCommand(machine, "df -h")
	if err != nil {
		return nil, fmt.Errorf("failed to get disk info: %w", err)
	}

	// Parse disk information
	diskDetails := c.parseDiskInfo(diskInfo)

	var diskFacts []*spookytypesfacts.DiskFacts
	for _, disk := range diskDetails {
		device, _ := disk["device"].(string)
		mountpoint, _ := disk["mountpoint"].(string)
		fstype, _ := disk["fstype"].(string)
		total, _ := disk["total"].(int64)
		used, _ := disk["used"].(int64)
		free, _ := disk["free"].(int64)

		diskFacts = append(diskFacts, &spookytypesfacts.DiskFacts{
			Device:     device,
			MountPoint: mountpoint,
			Filesystem: fstype,
			Total:      total,
			Used:       used,
			Free:       free,
		})
	}

	return diskFacts, nil
}

// collectNetworkFacts collects network facts via SSH
func (c *SystemFactCollector) collectNetworkFacts(machine *spookytypes.Machine) (*spookytypesfacts.NetworkFacts, error) {
	// Get hostname
	hostname, err := c.runSSHCommand(machine, "hostname")
	if err != nil {
		return nil, fmt.Errorf("failed to get hostname: %w", err)
	}

	// Get network interfaces
	interfaces, err := c.runSSHCommand(machine, "ip addr show")
	if err != nil {
		return nil, fmt.Errorf("failed to get network interfaces: %w", err)
	}

	// Parse network information
	networkDetails := c.parseNetworkInfo(interfaces)

	networkInterfaces, _ := networkDetails["interfaces"].([]*spookytypesfacts.NetworkInterface)
	ipAddresses, _ := networkDetails["ip_addresses"].([]string)
	primaryIP, _ := networkDetails["primary_ip"].(string)

	return &spookytypesfacts.NetworkFacts{
		Hostname:    strings.TrimSpace(hostname),
		Interfaces:  networkInterfaces,
		IPAddresses: ipAddresses,
		PrimaryIP:   primaryIP,
	}, nil
}

// collectLoadAverageFacts collects load average facts via SSH
func (c *SystemFactCollector) collectLoadAverageFacts(machine *spookytypes.Machine) (*spookytypesfacts.LoadAverageFacts, error) {
	// Get load average
	loadAvg, err := c.runSSHCommand(machine, "cat /proc/loadavg")
	if err != nil {
		return nil, fmt.Errorf("failed to get load average: %w", err)
	}

	// Parse load average
	loadParts := strings.Fields(strings.TrimSpace(loadAvg))
	if len(loadParts) < 3 {
		return nil, fmt.Errorf("invalid load average format: %s", loadAvg)
	}

	var load1, load5, load15 float64
	if _, err := fmt.Sscanf(loadParts[0], "%f", &load1); err != nil {
		load1 = 0.0
	}
	if _, err := fmt.Sscanf(loadParts[1], "%f", &load5); err != nil {
		load5 = 0.0
	}
	if _, err := fmt.Sscanf(loadParts[2], "%f", &load15); err != nil {
		load15 = 0.0
	}

	return &spookytypesfacts.LoadAverageFacts{
		Load1:  load1,
		Load5:  load5,
		Load15: load15,
	}, nil
}

// collectProcessFacts collects process facts via SSH
func (c *SystemFactCollector) collectProcessFacts(machine *spookytypes.Machine) (*spookytypesfacts.ProcessFacts, error) {
	// Get process count
	processCount, err := c.runSSHCommand(machine, "ps aux | wc -l")
	if err != nil {
		return nil, fmt.Errorf("failed to get process count: %w", err)
	}

	count := 0
	if _, err := fmt.Sscanf(strings.TrimSpace(processCount), "%d", &count); err != nil {
		// If parsing fails, default to 0
		count = 0
	}

	return &spookytypesfacts.ProcessFacts{
		Count: count,
	}, nil
}

// collectCollectorFacts collects facts from spooky-collector binary
func (c *SystemFactCollector) collectCollectorFacts(_ context.Context, machine *spookytypes.Machine) (*spookytypesfacts.CollectorFacts, error) {
	// Check if collector facts file exists
	checkCmd := "test -f /etc/spooky/facts.hcl && echo 'exists' || echo 'not_found'"
	result, err := c.runSSHCommand(machine, checkCmd)
	if err != nil {
		return nil, fmt.Errorf("failed to check collector facts file on %s: %w", machine.Hostname, err)
	}

	if strings.TrimSpace(result) == "not_found" {
		return nil, fmt.Errorf("collector facts file not found on %s: /etc/spooky/facts.hcl", machine.Hostname)
	}

	// Read collector facts file
	factsContent, err := c.runSSHCommand(machine, "cat /etc/spooky/facts.hcl")
	if err != nil {
		return nil, fmt.Errorf("failed to read collector facts from %s: %w", machine.Hostname, err)
	}

	// Parse HCL content using the HCL parser
	parser := NewHCLParser()
	collectorFacts, err := parser.ParseCollectorFacts(factsContent)
	if err != nil {
		return nil, fmt.Errorf("failed to parse collector facts from %s: %w", machine.Hostname, err)
	}

	return collectorFacts, nil
}

// collectCustomFacts collects custom facts from HCL files
func (c *SystemFactCollector) collectCustomFacts(_ context.Context, machine *spookytypes.Machine) (map[string]interface{}, error) {
	// Check if custom facts file exists (optional)
	checkCmd := "test -f /etc/spooky/custom.hcl && echo 'exists' || echo 'not_found'"
	result, err := c.runSSHCommand(machine, checkCmd)
	if err != nil {
		return nil, fmt.Errorf("failed to check custom facts file on %s: %w", machine.Hostname, err)
	}

	if strings.TrimSpace(result) == "not_found" {
		// Custom facts are optional, return empty map
		return make(map[string]interface{}), nil
	}

	// Read custom facts file
	factsContent, err := c.runSSHCommand(machine, "cat /etc/spooky/custom.hcl")
	if err != nil {
		return nil, fmt.Errorf("failed to read custom facts from %s: %w", machine.Hostname, err)
	}

	// Parse HCL content using the HCL parser
	parser := NewHCLParser()
	customFacts, err := parser.ParseCustomFacts(factsContent)
	if err != nil {
		return nil, fmt.Errorf("failed to parse custom facts from %s: %w", machine.Hostname, err)
	}

	return customFacts, nil
}

// runSSHCommand runs a command via SSH on the target machine
func (c *SystemFactCollector) runSSHCommand(machine *spookytypes.Machine, command string) (string, error) {
	ctx := context.Background()

	// Create connection request with actual machine configuration
	connectionRequest := &spookytypes.ConnectionRequest{
		Host:       machine.Hostname,
		Port:       machine.Port,
		User:       machine.User,
		Password:   machine.Password,
		KeyPath:    machine.KeyFile,
		Passphrase: machine.Passphrase,
		AuthMethod: spookytypesssh.AuthMethodPublicKey, // Default to public key
		Timeout:    30 * time.Second,
	}

	// Establish connection
	connectionResult, err := c.sshManager.Connect(ctx, connectionRequest)
	if err != nil {
		return "", fmt.Errorf("failed to establish SSH connection to %s: %w", machine.Hostname, err)
	}

	// Create session
	session, err := c.sshManager.CreateSession(ctx, connectionResult.Connection)
	if err != nil {
		return "", fmt.Errorf("failed to create SSH session on %s: %w", machine.Hostname, err)
	}

	// Create SSH command
	sshCommand := &spookytypes.SSHCommand{
		Command: command,
		Timeout: 30 * time.Second,
	}

	// Run command
	commandResult, err := c.sshManager.RunCommand(ctx, session, sshCommand)
	if err != nil {
		return "", fmt.Errorf("failed to run SSH command on %s: %w", machine.Hostname, err)
	}

	if !commandResult.Success {
		return "", fmt.Errorf("SSH command failed on %s with exit code %d: %s",
			machine.Hostname, commandResult.ExitCode, commandResult.Stderr)
	}

	return commandResult.Stdout, nil
}

// FactCollectionError represents fact collection errors
type FactCollectionError struct {
	Machine     string
	ErrorType   string
	Message     string
	Recoverable bool
}

func (e *FactCollectionError) Error() string {
	return fmt.Sprintf("fact collection error on %s (%s): %s", e.Machine, e.ErrorType, e.Message)
}

// Helper methods for parsing command output

func (c *SystemFactCollector) parseOSRelease(content string) map[string]string {
	result := make(map[string]string)
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		if strings.Contains(line, "=") {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				key := parts[0]
				value := strings.Trim(parts[1], "\"")
				result[key] = value
			}
		}
	}
	return result
}

func (c *SystemFactCollector) parseCPUInfo(content string) map[string]string {
	result := make(map[string]string)
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		if strings.Contains(line, ":") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])
				result[key] = value
			}
		}
	}
	return result
}

func (c *SystemFactCollector) parseMemoryInfo(content string) map[string]int64 {
	result := make(map[string]int64)
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		if strings.Contains(line, "Mem:") {
			parts := strings.Fields(line)
			if len(parts) >= 7 {
				var total, used, free, available int64
				if _, err := fmt.Sscanf(parts[1], "%d", &total); err != nil {
					total = 0
				}
				if _, err := fmt.Sscanf(parts[2], "%d", &used); err != nil {
					used = 0
				}
				if _, err := fmt.Sscanf(parts[3], "%d", &free); err != nil {
					free = 0
				}
				if _, err := fmt.Sscanf(parts[6], "%d", &available); err != nil {
					available = 0
				}
				result["total"] = total
				result["used"] = used
				result["free"] = free
				result["available"] = available
			}
		}
	}
	return result
}

func (c *SystemFactCollector) parseDiskInfo(content string) []map[string]interface{} {
	var result []map[string]interface{}
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "/dev/") {
			parts := strings.Fields(line)
			if len(parts) >= 6 {
				disk := map[string]interface{}{
					"device":     parts[0],
					"mountpoint": parts[5],
					"fstype":     parts[1],
					"total":      parts[2],
					"used":       parts[3],
					"free":       parts[4],
				}
				result = append(result, disk)
			}
		}
	}
	return result
}

func (c *SystemFactCollector) parseNetworkInfo(_ string) map[string]interface{} {
	result := make(map[string]interface{})
	// Simplified parsing - would need more sophisticated parsing
	result["interfaces"] = []string{}
	result["ip_addresses"] = []string{}
	result["primary_ip"] = ""
	return result
}
