package facts

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"spooky/internal/ssh"
)

// SSHCollector collects facts from remote servers via SSH
type SSHCollector struct {
	sshClient *ssh.SSHClient
}

// NewSSHCollector creates a new SSH-based fact collector
func NewSSHCollector(sshClient *ssh.SSHClient) *SSHCollector {
	return &SSHCollector{
		sshClient: sshClient,
	}
}

// Collect gathers all available facts from a remote server
func (c *SSHCollector) Collect(server string) (*FactCollection, error) {
	if c.sshClient == nil {
		return nil, fmt.Errorf("SSH client is nil")
	}
	base := NewBaseCollector()
	return base.CollectAll(server, c)
}

// CollectSpecific collects only the specified facts
func (c *SSHCollector) CollectSpecific(server string, keys []string) (*FactCollection, error) {
	if c.sshClient == nil {
		return nil, fmt.Errorf("SSH client is nil")
	}

	collection := &FactCollection{
		Server:    server,
		Timestamp: time.Now(),
		Facts:     make(map[string]*Fact),
	}

	for _, key := range keys {
		if err := c.collectSpecificFact(collection, key); err != nil {
			return nil, fmt.Errorf("failed to collect fact %s: %w", key, err)
		}
	}

	return collection, nil
}

// GetFact retrieves a single fact
func (c *SSHCollector) GetFact(server, key string) (*Fact, error) {
	if c.sshClient == nil {
		return nil, fmt.Errorf("SSH client is nil")
	}

	collection := &FactCollection{
		Server:    server,
		Timestamp: time.Now(),
		Facts:     make(map[string]*Fact),
	}

	if err := c.collectSpecificFact(collection, key); err != nil {
		return nil, err
	}

	if fact, exists := collection.Facts[key]; exists {
		return fact, nil
	}

	return nil, ErrFactNotFound(key, server)
}

// Helper function to create a fact
func (c *SSHCollector) createFact(collection *FactCollection, key, value string) {
	collection.Facts[key] = &Fact{
		Key:       key,
		Value:     value,
		Source:    string(SourceSSH),
		Server:    collection.Server,
		Timestamp: collection.Timestamp,
		TTL:       DefaultTTL,
	}
}

// Helper function to create a fact with any value type
func (c *SSHCollector) createFactWithValue(collection *FactCollection, key string, value interface{}) {
	collection.Facts[key] = &Fact{
		Key:       key,
		Value:     value,
		Source:    string(SourceSSH),
		Server:    collection.Server,
		Timestamp: collection.Timestamp,
		TTL:       DefaultTTL,
	}
}

// collectSystemFacts collects basic system identification facts
func (c *SSHCollector) collectSystemFacts(collection *FactCollection) error {
	// Machine ID
	if machineID, err := c.executeCommand("cat /etc/machine-id"); err == nil {
		c.createFact(collection, FactMachineID, strings.TrimSpace(machineID))
	}

	// Hostname
	if hostname, err := c.executeCommand("hostname"); err == nil {
		c.createFact(collection, FactHostname, strings.TrimSpace(hostname))
	}

	// FQDN
	if fqdn, err := c.executeCommand("hostname -f"); err == nil {
		c.createFact(collection, FactFQDN, strings.TrimSpace(fqdn))
	}

	return nil
}

// collectOSFacts collects operating system information
func (c *SSHCollector) collectOSFacts(collection *FactCollection) error {
	// OS release information
	if osRelease, err := c.executeCommand("cat /etc/os-release"); err == nil {
		osInfo := c.parseOSRelease(osRelease)

		c.createFact(collection, FactOSName, osInfo.Name)
		c.createFact(collection, FactOSVersion, osInfo.Version)
		c.createFact(collection, FactOSDistro, osInfo.Distribution)
	}

	// Architecture
	if arch, err := c.executeCommand("uname -m"); err == nil {
		c.createFact(collection, FactOSArch, strings.TrimSpace(arch))
	}

	// Kernel version
	if kernel, err := c.executeCommand("uname -r"); err == nil {
		c.createFact(collection, FactOSKernel, strings.TrimSpace(kernel))
	}

	return nil
}

// collectHardwareFacts collects hardware information
func (c *SSHCollector) collectHardwareFacts(collection *FactCollection) error {
	// CPU cores
	if cores, err := c.executeCommand("nproc"); err == nil {
		if coreCount, err := strconv.Atoi(strings.TrimSpace(cores)); err == nil {
			c.createFactWithValue(collection, FactCPUCores, coreCount)
		}
	}

	// CPU model
	if model, err := c.executeCommand("cat /proc/cpuinfo | grep 'model name' | head -1 | cut -d: -f2"); err == nil {
		c.createFact(collection, FactCPUModel, strings.TrimSpace(model))
	}

	// Memory information
	if memInfo, err := c.executeCommand("cat /proc/meminfo"); err == nil {
		memory := c.parseMemInfo(memInfo)

		c.createFactWithValue(collection, FactMemoryTotal, memory.Total)
		c.createFactWithValue(collection, FactMemoryUsed, memory.Used)
		c.createFactWithValue(collection, FactMemoryAvail, memory.Available)
	}

	// Disk information
	if dfOutput, err := c.executeCommand("df -B1 /"); err == nil {
		disk := c.parseDiskInfo(dfOutput)

		c.createFactWithValue(collection, FactDiskTotal, disk.Total)
		c.createFactWithValue(collection, FactDiskUsed, disk.Used)
		c.createFactWithValue(collection, FactDiskAvail, disk.Available)
	}

	return nil
}

// collectNetworkFacts collects network information
func (c *SSHCollector) collectNetworkFacts(collection *FactCollection) error {
	// IP addresses
	if ipOutput, err := c.executeCommand("ip -json addr show"); err == nil {
		ips := c.parseIPAddresses(ipOutput)
		c.createFactWithValue(collection, FactNetworkIPs, ips)
	}

	// MAC addresses
	if macOutput, err := c.executeCommand("ip -json link show"); err == nil {
		macs := c.parseMACAddresses(macOutput)
		c.createFactWithValue(collection, FactNetworkMACs, macs)
	}

	// DNS configuration
	if resolvConf, err := c.executeCommand("cat /etc/resolv.conf"); err == nil {
		dns := c.parseDNSConfig(resolvConf)
		c.createFactWithValue(collection, FactDNS, dns)
	}

	return nil
}

// collectEnvironmentFacts collects environment variables
func (c *SSHCollector) collectEnvironmentFacts(collection *FactCollection) error {
	if envOutput, err := c.executeCommand("env"); err == nil {
		env := c.parseEnvironment(envOutput)
		c.createFactWithValue(collection, FactEnvironment, env)
	}

	return nil
}

// collectSpecificFact collects a specific fact based on the key
func (c *SSHCollector) collectSpecificFact(collection *FactCollection, key string) error {
	collector, exists := c.getFactCollector(key)
	if !exists {
		return fmt.Errorf("unknown fact key: %s", key)
	}

	return collector(collection)
}

// factCollector is a function type that collects a specific fact
type factCollector func(*FactCollection) error

// getFactCollector returns the appropriate fact collector for the given key
func (c *SSHCollector) getFactCollector(key string) (factCollector, bool) {
	collectors := map[string]factCollector{
		FactMachineID:   c.collectMachineID,
		FactHostname:    c.collectHostname,
		FactFQDN:        c.collectFQDN,
		FactOSName:      c.collectOSName,
		FactOSVersion:   c.collectOSVersion,
		FactOSDistro:    c.collectOSDistro,
		FactOSArch:      c.collectOSArch,
		FactOSKernel:    c.collectOSKernel,
		FactCPUCores:    c.collectCPUCores,
		FactCPUModel:    c.collectCPUModel,
		FactMemoryTotal: c.collectMemoryTotal,
		FactMemoryUsed:  c.collectMemoryUsed,
		FactMemoryAvail: c.collectMemoryAvail,
		FactDiskTotal:   c.collectDiskTotal,
		FactDiskUsed:    c.collectDiskUsed,
		FactDiskAvail:   c.collectDiskAvail,
		FactNetworkIPs:  c.collectNetworkIPs,
		FactNetworkMACs: c.collectNetworkMACs,
		FactDNS:         c.collectDNS,
		FactEnvironment: c.collectEnvironment,
	}

	collector, exists := collectors[key]
	return collector, exists
}

// Individual fact collector methods
func (c *SSHCollector) collectMachineID(collection *FactCollection) error {
	if machineID, err := c.executeCommand("cat /etc/machine-id"); err == nil {
		c.createFact(collection, FactMachineID, strings.TrimSpace(machineID))
	}
	return nil
}

func (c *SSHCollector) collectHostname(collection *FactCollection) error {
	if hostname, err := c.executeCommand("hostname"); err == nil {
		c.createFact(collection, FactHostname, strings.TrimSpace(hostname))
	}
	return nil
}

func (c *SSHCollector) collectFQDN(collection *FactCollection) error {
	if fqdn, err := c.executeCommand("hostname -f"); err == nil {
		c.createFact(collection, FactFQDN, strings.TrimSpace(fqdn))
	}
	return nil
}

func (c *SSHCollector) collectOSName(collection *FactCollection) error {
	if osRelease, err := c.executeCommand("cat /etc/os-release"); err == nil {
		osInfo := c.parseOSRelease(osRelease)
		c.createFact(collection, FactOSName, osInfo.Name)
	}
	return nil
}

func (c *SSHCollector) collectOSVersion(collection *FactCollection) error {
	if osRelease, err := c.executeCommand("cat /etc/os-release"); err == nil {
		osInfo := c.parseOSRelease(osRelease)
		c.createFact(collection, FactOSVersion, osInfo.Version)
	}
	return nil
}

func (c *SSHCollector) collectOSDistro(collection *FactCollection) error {
	if osRelease, err := c.executeCommand("cat /etc/os-release"); err == nil {
		osInfo := c.parseOSRelease(osRelease)
		c.createFact(collection, FactOSDistro, osInfo.Distribution)
	}
	return nil
}

func (c *SSHCollector) collectOSArch(collection *FactCollection) error {
	if arch, err := c.executeCommand("uname -m"); err == nil {
		c.createFact(collection, FactOSArch, strings.TrimSpace(arch))
	}
	return nil
}

func (c *SSHCollector) collectOSKernel(collection *FactCollection) error {
	if kernel, err := c.executeCommand("uname -r"); err == nil {
		c.createFact(collection, FactOSKernel, strings.TrimSpace(kernel))
	}
	return nil
}

func (c *SSHCollector) collectCPUCores(collection *FactCollection) error {
	if cores, err := c.executeCommand("nproc"); err == nil {
		if coreCount, err := strconv.Atoi(strings.TrimSpace(cores)); err == nil {
			c.createFactWithValue(collection, FactCPUCores, coreCount)
		}
	}
	return nil
}

func (c *SSHCollector) collectCPUModel(collection *FactCollection) error {
	if model, err := c.executeCommand("cat /proc/cpuinfo | grep 'model name' | head -1 | cut -d: -f2"); err == nil {
		c.createFact(collection, FactCPUModel, strings.TrimSpace(model))
	}
	return nil
}

func (c *SSHCollector) collectMemoryTotal(collection *FactCollection) error {
	if memInfo, err := c.executeCommand("cat /proc/meminfo"); err == nil {
		memory := c.parseMemInfo(memInfo)
		c.createFactWithValue(collection, FactMemoryTotal, memory.Total)
	}
	return nil
}

func (c *SSHCollector) collectMemoryUsed(collection *FactCollection) error {
	if memInfo, err := c.executeCommand("cat /proc/meminfo"); err == nil {
		memory := c.parseMemInfo(memInfo)
		c.createFactWithValue(collection, FactMemoryUsed, memory.Used)
	}
	return nil
}

func (c *SSHCollector) collectMemoryAvail(collection *FactCollection) error {
	if memInfo, err := c.executeCommand("cat /proc/meminfo"); err == nil {
		memory := c.parseMemInfo(memInfo)
		c.createFactWithValue(collection, FactMemoryAvail, memory.Available)
	}
	return nil
}

func (c *SSHCollector) collectDiskTotal(collection *FactCollection) error {
	if dfOutput, err := c.executeCommand("df -B1 /"); err == nil {
		disk := c.parseDiskInfo(dfOutput)
		c.createFactWithValue(collection, FactDiskTotal, disk.Total)
	}
	return nil
}

func (c *SSHCollector) collectDiskUsed(collection *FactCollection) error {
	if dfOutput, err := c.executeCommand("df -B1 /"); err == nil {
		disk := c.parseDiskInfo(dfOutput)
		c.createFactWithValue(collection, FactDiskUsed, disk.Used)
	}
	return nil
}

func (c *SSHCollector) collectDiskAvail(collection *FactCollection) error {
	if dfOutput, err := c.executeCommand("df -B1 /"); err == nil {
		disk := c.parseDiskInfo(dfOutput)
		c.createFactWithValue(collection, FactDiskAvail, disk.Available)
	}
	return nil
}

func (c *SSHCollector) collectNetworkIPs(collection *FactCollection) error {
	if ipOutput, err := c.executeCommand("ip -json addr show"); err == nil {
		ips := c.parseIPAddresses(ipOutput)
		c.createFactWithValue(collection, FactNetworkIPs, ips)
	}
	return nil
}

func (c *SSHCollector) collectNetworkMACs(collection *FactCollection) error {
	if macOutput, err := c.executeCommand("ip -json link show"); err == nil {
		macs := c.parseMACAddresses(macOutput)
		c.createFactWithValue(collection, FactNetworkMACs, macs)
	}
	return nil
}

func (c *SSHCollector) collectDNS(collection *FactCollection) error {
	if resolvConf, err := c.executeCommand("cat /etc/resolv.conf"); err == nil {
		dns := c.parseDNSConfig(resolvConf)
		c.createFactWithValue(collection, FactDNS, dns)
	}
	return nil
}

func (c *SSHCollector) collectEnvironment(collection *FactCollection) error {
	if envOutput, err := c.executeCommand("env"); err == nil {
		env := c.parseEnvironment(envOutput)
		c.createFactWithValue(collection, FactEnvironment, env)
	}
	return nil
}

// executeCommand executes a command on the remote server
func (c *SSHCollector) executeCommand(command string) (string, error) {
	if c.sshClient == nil {
		return "", fmt.Errorf("SSH client is nil")
	}

	// Execute the command using the SSH client
	output, err := c.sshClient.ExecuteCommand(command)
	if err != nil {
		return "", fmt.Errorf("failed to execute command '%s': %w", command, err)
	}

	return output, nil
}

// Helper methods for parsing command output
func (c *SSHCollector) parseOSRelease(osRelease string) OSInfo {
	osInfo := OSInfo{}

	lines := strings.Split(osRelease, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Remove quotes from values
		line = strings.Trim(line, `"'`)

		// Extract the prefix for switch statement
		switch {
		case strings.HasPrefix(line, "NAME="):
			osInfo.Name = strings.Trim(strings.TrimPrefix(line, "NAME="), `"'`)
		case strings.HasPrefix(line, "VERSION="):
			osInfo.Version = strings.Trim(strings.TrimPrefix(line, "VERSION="), `"'`)
		case strings.HasPrefix(line, "ID="):
			osInfo.Distribution = strings.Trim(strings.TrimPrefix(line, "ID="), `"'`)
		}
	}

	return osInfo
}

func (c *SSHCollector) parseMemInfo(memInfo string) MemoryInfo {
	memory := MemoryInfo{}

	lines := strings.Split(memInfo, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}

		value, err := strconv.ParseUint(parts[1], 10, 64)
		if err != nil {
			continue
		}

		// Convert from kB to bytes
		value *= 1024

		switch parts[0] {
		case "MemTotal:":
			memory.Total = value
		case "MemAvailable:":
			memory.Available = value
		}
	}

	// Calculate used memory
	if memory.Total > 0 && memory.Available > 0 {
		memory.Used = memory.Total - memory.Available
	}

	return memory
}

func (c *SSHCollector) parseDiskInfo(dfOutput string) DiskInfo {
	disk := DiskInfo{}

	lines := strings.Split(dfOutput, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "Filesystem") {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) < 4 {
			continue
		}

		// Parse the root filesystem line
		if parts[5] == "/" {
			if total, err := strconv.ParseUint(parts[1], 10, 64); err == nil {
				disk.Total = total
			}
			if used, err := strconv.ParseUint(parts[2], 10, 64); err == nil {
				disk.Used = used
			}
			if available, err := strconv.ParseUint(parts[3], 10, 64); err == nil {
				disk.Available = available
			}
			disk.Device = parts[0]
			disk.MountPoint = parts[5]
			break
		}
	}

	return disk
}

func (c *SSHCollector) parseIPAddresses(_ string) []string {
	var ips []string

	// For now, return a simple mock response
	// In a real implementation, this would parse the JSON output from "ip -json addr"
	ips = append(ips, "127.0.0.1/8", "192.168.1.100/24")

	return ips
}

func (c *SSHCollector) parseMACAddresses(_ string) []string {
	var macs []string

	// For now, return a simple mock response
	// In a real implementation, this would parse the JSON output from "ip -json link"
	macs = append(macs, "00:00:00:00:00:00", "52:54:00:12:34:56")

	return macs
}

func (c *SSHCollector) parseDNSConfig(resolvConf string) DNSInfo {
	dns := DNSInfo{}

	lines := strings.Split(resolvConf, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if strings.HasPrefix(line, "nameserver") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				dns.Nameservers = append(dns.Nameservers, parts[1])
			}
		} else if strings.HasPrefix(line, "search") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				dns.Search = append(dns.Search, parts[1])
			}
		}
	}

	return dns
}

func (c *SSHCollector) parseEnvironment(envOutput string) map[string]string {
	env := make(map[string]string)

	lines := strings.Split(envOutput, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			env[parts[0]] = parts[1]
		}
	}

	return env
}
