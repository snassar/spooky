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
	base := NewBaseCollector()
	return base.CollectAll(server, c)
}

// CollectSpecific collects only the specified facts
func (c *SSHCollector) CollectSpecific(server string, keys []string) (*FactCollection, error) {
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

	return nil, fmt.Errorf("fact %s not found for server %s", key, server)
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
	switch key {
	case FactMachineID, FactHostname, FactFQDN:
		return c.collectSystemFacts(collection)
	case FactOSName, FactOSVersion, FactOSDistro, FactOSArch, FactOSKernel:
		return c.collectOSFacts(collection)
	case FactCPUCores, FactCPUModel, FactMemoryTotal, FactMemoryUsed, FactMemoryAvail, FactDiskTotal, FactDiskUsed, FactDiskAvail:
		return c.collectHardwareFacts(collection)
	case FactNetworkIPs, FactNetworkMACs, FactDNS:
		return c.collectNetworkFacts(collection)
	case FactEnvironment:
		return c.collectEnvironmentFacts(collection)
	default:
		return fmt.Errorf("unknown fact key: %s", key)
	}
}

// executeCommand executes a command on the remote server
func (c *SSHCollector) executeCommand(_ string) (string, error) {
	// This will need to be implemented using the existing SSH client
	// For now, return a placeholder
	return "", fmt.Errorf("SSH command execution not yet implemented")
}

// Helper methods for parsing command output
func (c *SSHCollector) parseOSRelease(_ string) OSInfo {
	// TODO: Implement parsing of /etc/os-release
	return OSInfo{}
}

func (c *SSHCollector) parseMemInfo(_ string) MemoryInfo {
	// TODO: Implement parsing of /proc/meminfo
	return MemoryInfo{}
}

func (c *SSHCollector) parseDiskInfo(_ string) DiskInfo {
	// TODO: Implement parsing of df output
	return DiskInfo{}
}

func (c *SSHCollector) parseIPAddresses(_ string) []string {
	// TODO: Implement parsing of ip addr output
	return []string{}
}

func (c *SSHCollector) parseMACAddresses(_ string) []string {
	// TODO: Implement parsing of ip link output
	return []string{}
}

func (c *SSHCollector) parseDNSConfig(_ string) DNSInfo {
	// TODO: Implement parsing of /etc/resolv.conf
	return DNSInfo{}
}

func (c *SSHCollector) parseEnvironment(_ string) map[string]string {
	// TODO: Implement parsing of env output
	return map[string]string{}
}
