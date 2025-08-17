// Package facts provides fact collection and in-memory management functionality.
package facts

import (
	"context"
	"crypto/sha256"
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
	logger     spookytypes.Logger
}

// NewSystemFactCollector creates a new system fact collector
func NewSystemFactCollector(sshManager spookyinterfaces.SSHManager, logger spookytypes.Logger) *SystemFactCollector {
	return &SystemFactCollector{
		name:       "system",
		sshManager: sshManager,
		logger:     logger,
	}
}

// GetName returns the collector name
func (c *SystemFactCollector) GetName() string {
	return c.name
}

// Collect collects facts from the given machine
func (c *SystemFactCollector) Collect(ctx context.Context, machine interface{}) (*spookytypes.FactCollection, error) {
	// Type assert to get machine details
	machineObj, ok := machine.(*spookytypes.Machine)
	if !ok {
		return nil, fmt.Errorf("machine must be of type *spookytypes.Machine")
	}

	// Get machine ID
	machineID, err := c.getMachineID(machineObj)
	if err != nil {
		return nil, fmt.Errorf("failed to get machine ID: %w", err)
	}

	// Collect system facts
	facts := &spookytypesfacts.Facts{
		System: &spookytypesfacts.SystemFacts{},
	}

	// Collect system facts via SSH
	systemFacts, err := c.collectSystemFacts(ctx, machineObj)
	if err != nil {
		return nil, fmt.Errorf("failed to collect system facts: %w", err)
	}
	facts.System = systemFacts

	// Collect collector facts via SSH (if available)
	collectorFacts, err := c.collectCollectorFacts(ctx, machineObj)
	if err != nil {
		// Collector facts are optional, log but don't fail
		fmt.Printf("Warning: failed to collect collector facts: %v\n", err)
	} else {
		facts.Collector = collectorFacts
	}

	// Collect custom facts via SSH (if available)
	customFacts, err := c.collectCustomFacts(ctx, machineObj)
	if err != nil {
		// Custom facts are optional, log but don't fail
		fmt.Printf("Warning: failed to collect custom facts: %v\n", err)
	} else {
		facts.Custom = customFacts
	}

	// Create fact collection
	collection := &spookytypes.FactCollection{
		MachineID:   machineID,
		CollectedAt: time.Now(),
		Facts:       facts,
		Metadata: map[string]interface{}{
			"collector": c.name,
			"machine":   machineObj.Hostname,
		},
	}

	return collection, nil
}

// CollectViaSSH collects facts from remote machine via SSH
func (c *SystemFactCollector) CollectViaSSH(ctx context.Context, machine *spookytypes.Machine) (*spookytypes.FactCollection, error) {
	c.logger.Info("Starting SSH-based fact collection", map[string]interface{}{
		"machine": machine.Hostname,
		"host":    machine.Host,
	})

	// Validate machine configuration for SSH
	if err := c.validateSSHConfiguration(machine); err != nil {
		return nil, fmt.Errorf("invalid SSH configuration: %w", err)
	}

	// Get machine ID (with fallback to generated ID if /etc/machine-id doesn't exist)
	machineID, err := c.getMachineID(machine)
	if err != nil {
		return nil, fmt.Errorf("failed to get machine ID: %w", err)
	}

	// Create fact collection with machine information
	factCollection := &spookytypesfacts.FactCollection{
		MachineID:   machineID,
		CollectedAt: time.Now(),
		Facts:       &spookytypesfacts.Facts{},
	}

	// Collect system facts via SSH
	systemFacts, err := c.collectSystemFactsViaSSH(ctx, machine)
	if err != nil {
		c.logger.Error("Failed to collect system facts via SSH", err, map[string]interface{}{
			"machine": machine.Hostname,
		})
		return nil, fmt.Errorf("failed to collect system facts: %w", err)
	}
	factCollection.Facts.System = systemFacts

	// Collect custom facts from /etc/spooky/custom.hcl via SSH
	customFacts, err := c.collectCustomFactsViaSSH(ctx, machine)
	if err != nil {
		c.logger.Warn("Failed to collect custom facts via SSH", map[string]interface{}{
			"machine": machine.Hostname,
			"error":   err.Error(),
		})
		// Don't fail the entire collection for custom facts
	} else {
		factCollection.Facts.Custom = customFacts
	}

	c.logger.Info("Completed SSH-based fact collection", map[string]interface{}{
		"machine": machine.Hostname,
		"host":    machine.Host,
	})

	return factCollection, nil
}

// validateSSHConfiguration validates machine configuration for SSH operations
func (c *SystemFactCollector) validateSSHConfiguration(machine *spookytypes.Machine) error {
	if machine.Host == "" {
		return fmt.Errorf("machine host is required for SSH operations")
	}

	if machine.User == "" {
		return fmt.Errorf("machine user is required for SSH operations")
	}

	if machine.Port <= 0 || machine.Port > 65535 {
		return fmt.Errorf("invalid SSH port: %d", machine.Port)
	}

	// Check that either password or key file is provided
	if machine.Password == "" && machine.KeyFile == "" {
		return fmt.Errorf("either password or key file must be provided for SSH authentication")
	}

	// If key file is provided, check if passphrase is required
	if machine.KeyFile != "" {
		if machine.Passphrase == "" {
			c.logger.Debug("SSH key file provided without passphrase", map[string]interface{}{
				"machine":  machine.Hostname,
				"key_file": machine.KeyFile,
			})
		} else {
			c.logger.Debug("SSH key file provided with passphrase", map[string]interface{}{
				"machine":  machine.Hostname,
				"key_file": machine.KeyFile,
			})
		}
	}

	return nil
}

// collectSystemFactsViaSSH collects system facts from remote machine via SSH
func (c *SystemFactCollector) collectSystemFactsViaSSH(ctx context.Context, machine *spookytypes.Machine) (*spookytypesfacts.SystemFacts, error) {
	systemFacts := &spookytypesfacts.SystemFacts{}

	// Collect OS facts
	osFacts, err := c.collectOSFactsViaSSH(ctx, machine)
	if err != nil {
		return nil, fmt.Errorf("failed to collect OS facts: %w", err)
	}
	systemFacts.OS = osFacts

	// Collect hardware facts
	hardwareFacts, err := c.collectHardwareFactsViaSSH(ctx, machine)
	if err != nil {
		return nil, fmt.Errorf("failed to collect hardware facts: %w", err)
	}
	systemFacts.Hardware = hardwareFacts

	// Collect network facts
	networkFacts, err := c.collectNetworkFactsViaSSH(ctx, machine)
	if err != nil {
		return nil, fmt.Errorf("failed to collect network facts: %w", err)
	}
	systemFacts.Network = networkFacts

	return systemFacts, nil
}

// collectOSFactsViaSSH collects OS facts from remote machine via SSH
func (c *SystemFactCollector) collectOSFactsViaSSH(_ context.Context, machine *spookytypes.Machine) (*spookytypesfacts.OSFacts, error) {
	// Collect OS information
	osInfo, err := c.runSSHCommand(machine, "cat /etc/os-release")
	if err != nil {
		return nil, fmt.Errorf("failed to get OS info: %w", err)
	}

	// Parse OS information
	osInfoMap := c.parseOSRelease(osInfo)

	// Get kernel version
	kernel, err := c.runSSHCommand(machine, "uname -r")
	if err == nil {
		osInfoMap["KERNEL"] = strings.TrimSpace(kernel)
	}

	// Get architecture
	arch, err := c.runSSHCommand(machine, "uname -m")
	if err == nil {
		osInfoMap["ARCH"] = strings.TrimSpace(arch)
	}

	// Create OSFacts struct
	osFacts := &spookytypesfacts.OSFacts{
		Name:     osInfoMap["NAME"],
		Version:  osInfoMap["VERSION"],
		Arch:     osInfoMap["ARCH"],
		Kernel:   osInfoMap["KERNEL"],
		Platform: osInfoMap["ID"],
		Family:   osInfoMap["ID_LIKE"],
	}

	return osFacts, nil
}

// generateMachineID generates a machine ID for the given hostname
// Based on systemd machine-id specification: https://www.freedesktop.org/software/systemd/man/latest/machine-id.html
func generateMachineID(hostname string) string {
	// Create a deterministic but unique machine ID based on hostname
	// This follows the systemd machine-id format: 32-character hexadecimal string

	// Use SHA-256 hash of hostname to generate a proper 32-byte value
	hash := sha256.Sum256([]byte(hostname))

	// Convert to 32-character hexadecimal string
	return fmt.Sprintf("%032x", hash)
}

// collectHardwareFactsViaSSH collects hardware facts from remote machine via SSH
func (c *SystemFactCollector) collectHardwareFactsViaSSH(_ context.Context, machine *spookytypes.Machine) (*spookytypesfacts.HardwareFacts, error) {
	hardwareFacts := &spookytypesfacts.HardwareFacts{}

	// Collect CPU information
	cpuInfo, err := c.runSSHCommand(machine, "cat /proc/cpuinfo")
	if err == nil {
		cpuInfoMap := c.parseCPUInfo(cpuInfo)
		hardwareFacts.CPU = &spookytypesfacts.CPUFacts{
			Cores:        1, // Default value
			Model:        cpuInfoMap["model name"],
			Frequency:    0.0, // Default value
			Architecture: cpuInfoMap["cpu family"],
			Vendor:       cpuInfoMap["vendor_id"],
		}
	}

	// Collect memory information
	memInfo, err := c.runSSHCommand(machine, "cat /proc/meminfo")
	if err == nil {
		memInfoMap := c.parseMemInfo(memInfo)
		hardwareFacts.Memory = &spookytypesfacts.MemoryFacts{
			Total: memInfoMap["total"],
			Used:  memInfoMap["used"],
			Free:  memInfoMap["free"],
		}
	}

	// Collect disk information
	diskInfo, err := c.runSSHCommand(machine, "df -h")
	if err == nil {
		hardwareFacts.Disks = c.parseDiskInfo(diskInfo)
	}

	return hardwareFacts, nil
}

// parseMemInfo parses memory information from /proc/meminfo
func (c *SystemFactCollector) parseMemInfo(content string) map[string]int64 {
	result := make(map[string]int64)
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		switch {
		case strings.Contains(line, "MemTotal:"):
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				if val, err := strconv.ParseInt(parts[1], 10, 64); err == nil {
					result["total"] = val * 1024 // Convert KB to bytes
				}
			}
		case strings.Contains(line, "MemAvailable:"):
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				if val, err := strconv.ParseInt(parts[1], 10, 64); err == nil {
					result["available"] = val * 1024 // Convert KB to bytes
				}
			}
		case strings.Contains(line, "MemFree:"):
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				if val, err := strconv.ParseInt(parts[1], 10, 64); err == nil {
					result["free"] = val * 1024 // Convert KB to bytes
				}
			}
		}
	}

	// Calculate used memory
	if total, ok := result["total"]; ok {
		if available, ok := result["available"]; ok {
			result["used"] = total - available
		}
	}

	return result
}

// parseNetworkInterfaces parses network interfaces from ip link show output
func (c *SystemFactCollector) parseNetworkInterfaces(content string) []*spookytypesfacts.NetworkInterface {
	var interfaces []*spookytypesfacts.NetworkInterface

	lines := strings.Split(content, "\n")
	for _, line := range lines {
		if strings.Contains(line, ":") && !strings.Contains(line, "lo:") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				name := strings.TrimSuffix(parts[1], ":")
				interfaces = append(interfaces, &spookytypesfacts.NetworkInterface{
					Name: name,
				})
			}
		}
	}

	return interfaces
}

// parseDiskInfo parses disk information from df command output
func (c *SystemFactCollector) parseDiskInfo(content string) []*spookytypesfacts.DiskFacts {
	var disks []*spookytypesfacts.DiskFacts

	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "Filesystem") {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 6 {
			continue
		}

		// Parse df output: Filesystem Size Used Avail Use% Mounted on
		disk := &spookytypesfacts.DiskFacts{
			Device:     fields[0],
			MountPoint: fields[5],
			Filesystem: fields[0], // Use device as filesystem for now
			Total:      0,         // Would need to parse size string
			Used:       0,         // Would need to parse used string
			Free:       0,         // Would need to parse avail string
		}

		disks = append(disks, disk)
	}

	return disks
}

// parseHCLFacts parses HCL facts from content
func (c *SystemFactCollector) parseHCLFacts(_ string) (map[string]interface{}, error) {
	// For now, return empty map - would need to implement HCL parsing
	return make(map[string]interface{}), nil
}

// collectNetworkFactsViaSSH collects network facts from remote machine via SSH
func (c *SystemFactCollector) collectNetworkFactsViaSSH(_ context.Context, machine *spookytypes.Machine) (*spookytypesfacts.NetworkFacts, error) {
	networkFacts := &spookytypesfacts.NetworkFacts{}

	// Collect network interface information
	netInfo, err := c.runSSHCommand(machine, "ip addr show")
	if err == nil {
		networkFacts.Interfaces = c.parseNetworkInterfaces(netInfo)
	}

	return networkFacts, nil
}

// collectCustomFactsViaSSH collects custom facts from /etc/spooky/custom.hcl via SSH
func (c *SystemFactCollector) collectCustomFactsViaSSH(_ context.Context, machine *spookytypes.Machine) (map[string]interface{}, error) {
	// Try to read /etc/spooky/custom.hcl
	customFactsContent, err := c.runSSHCommand(machine, "cat /etc/spooky/custom.hcl 2>/dev/null")
	if err != nil {
		// File doesn't exist or can't be read - this is not an error
		return nil, nil
	}

	if strings.TrimSpace(customFactsContent) == "" {
		return nil, nil
	}

	// Parse HCL content
	customFacts, err := c.parseHCLFacts(customFactsContent)
	if err != nil {
		return nil, fmt.Errorf("failed to parse custom facts: %w", err)
	}

	return customFacts, nil
}

// getMachineID gets the machine ID from /etc/machine-id via SSH, with fallback to generated ID
func (c *SystemFactCollector) getMachineID(machine *spookytypes.Machine) (string, error) {
	// First, try to read /etc/machine-id via SSH
	output, err := c.runSSHCommand(machine, "cat /etc/machine-id 2>/dev/null")
	if err != nil {
		// Machine ID file doesn't exist or can't be read - generate one based on hostname
		c.logger.Info("Machine ID not found, generating one based on hostname", map[string]interface{}{
			"machine": machine.Hostname,
			"host":    machine.Host,
		})
		return generateMachineID(machine.Hostname), nil
	}

	machineID := strings.TrimSpace(output)

	// Check if the file exists but is empty
	if machineID == "" {
		// Empty machine ID file - generate one based on hostname
		c.logger.Info("Machine ID file is empty, generating one based on hostname", map[string]interface{}{
			"machine": machine.Hostname,
			"host":    machine.Host,
		})
		return generateMachineID(machine.Hostname), nil
	}

	// Validate machine ID format (32-character hex string)
	if !regexp.MustCompile(`^[a-f0-9]{32}$`).MatchString(machineID) {
		// Invalid machine ID format - generate one based on hostname
		c.logger.Warn("Invalid machine ID format, generating one based on hostname", map[string]interface{}{
			"machine":    machine.Hostname,
			"host":       machine.Host,
			"invalid_id": machineID,
		})
		return generateMachineID(machine.Hostname), nil
	}

	c.logger.Info("Successfully retrieved machine ID from /etc/machine-id", map[string]interface{}{
		"machine":    machine.Hostname,
		"host":       machine.Host,
		"machine_id": machineID,
	})

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
	// Get disk information
	diskInfo, err := c.runSSHCommand(machine, "df -h")
	if err != nil {
		return nil, fmt.Errorf("failed to get disk info: %w", err)
	}

	// Parse disk information using the new method
	return c.parseDiskInfo(diskInfo), nil
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

func (c *SystemFactCollector) parseNetworkInfo(_ string) map[string]interface{} {
	result := make(map[string]interface{})
	// Simplified parsing - would need more sophisticated parsing
	result["interfaces"] = []string{}
	result["ip_addresses"] = []string{}
	result["primary_ip"] = ""
	return result
}
