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
	collection := &FactCollection{
		Server:    server,
		Timestamp: time.Now(),
		Facts:     make(map[string]*Fact),
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

	return collection, nil
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

// collectSystemFacts collects basic system identification facts
func (c *SSHCollector) collectSystemFacts(collection *FactCollection) error {
	// Machine ID
	if machineID, err := c.executeCommand("cat /etc/machine-id"); err == nil {
		collection.Facts[FactMachineID] = &Fact{
			Key:       FactMachineID,
			Value:     strings.TrimSpace(machineID),
			Source:    string(SourceSSH),
			Server:    collection.Server,
			Timestamp: collection.Timestamp,
			TTL:       DefaultTTL,
		}
	}

	// Hostname
	if hostname, err := c.executeCommand("hostname"); err == nil {
		collection.Facts[FactHostname] = &Fact{
			Key:       FactHostname,
			Value:     strings.TrimSpace(hostname),
			Source:    string(SourceSSH),
			Server:    collection.Server,
			Timestamp: collection.Timestamp,
			TTL:       DefaultTTL,
		}
	}

	// FQDN
	if fqdn, err := c.executeCommand("hostname -f"); err == nil {
		collection.Facts[FactFQDN] = &Fact{
			Key:       FactFQDN,
			Value:     strings.TrimSpace(fqdn),
			Source:    string(SourceSSH),
			Server:    collection.Server,
			Timestamp: collection.Timestamp,
			TTL:       DefaultTTL,
		}
	}

	return nil
}

// collectOSFacts collects operating system information
func (c *SSHCollector) collectOSFacts(collection *FactCollection) error {
	// OS release information
	if osRelease, err := c.executeCommand("cat /etc/os-release"); err == nil {
		osInfo := c.parseOSRelease(osRelease)

		collection.Facts[FactOSName] = &Fact{
			Key:       FactOSName,
			Value:     osInfo.Name,
			Source:    string(SourceSSH),
			Server:    collection.Server,
			Timestamp: collection.Timestamp,
			TTL:       DefaultTTL,
		}

		collection.Facts[FactOSVersion] = &Fact{
			Key:       FactOSVersion,
			Value:     osInfo.Version,
			Source:    string(SourceSSH),
			Server:    collection.Server,
			Timestamp: collection.Timestamp,
			TTL:       DefaultTTL,
		}

		collection.Facts[FactOSDistro] = &Fact{
			Key:       FactOSDistro,
			Value:     osInfo.Distribution,
			Source:    string(SourceSSH),
			Server:    collection.Server,
			Timestamp: collection.Timestamp,
			TTL:       DefaultTTL,
		}
	}

	// Architecture
	if arch, err := c.executeCommand("uname -m"); err == nil {
		collection.Facts[FactOSArch] = &Fact{
			Key:       FactOSArch,
			Value:     strings.TrimSpace(arch),
			Source:    string(SourceSSH),
			Server:    collection.Server,
			Timestamp: collection.Timestamp,
			TTL:       DefaultTTL,
		}
	}

	// Kernel version
	if kernel, err := c.executeCommand("uname -r"); err == nil {
		collection.Facts[FactOSKernel] = &Fact{
			Key:       FactOSKernel,
			Value:     strings.TrimSpace(kernel),
			Source:    string(SourceSSH),
			Server:    collection.Server,
			Timestamp: collection.Timestamp,
			TTL:       DefaultTTL,
		}
	}

	return nil
}

// collectHardwareFacts collects hardware information
func (c *SSHCollector) collectHardwareFacts(collection *FactCollection) error {
	// CPU cores
	if cores, err := c.executeCommand("nproc"); err == nil {
		if coreCount, err := strconv.Atoi(strings.TrimSpace(cores)); err == nil {
			collection.Facts[FactCPUCores] = &Fact{
				Key:       FactCPUCores,
				Value:     coreCount,
				Source:    string(SourceSSH),
				Server:    collection.Server,
				Timestamp: collection.Timestamp,
				TTL:       DefaultTTL,
			}
		}
	}

	// CPU model
	if model, err := c.executeCommand("cat /proc/cpuinfo | grep 'model name' | head -1 | cut -d: -f2"); err == nil {
		collection.Facts[FactCPUModel] = &Fact{
			Key:       FactCPUModel,
			Value:     strings.TrimSpace(model),
			Source:    string(SourceSSH),
			Server:    collection.Server,
			Timestamp: collection.Timestamp,
			TTL:       DefaultTTL,
		}
	}

	// Memory information
	if memInfo, err := c.executeCommand("cat /proc/meminfo"); err == nil {
		memory := c.parseMemInfo(memInfo)

		collection.Facts[FactMemoryTotal] = &Fact{
			Key:       FactMemoryTotal,
			Value:     memory.Total,
			Source:    string(SourceSSH),
			Server:    collection.Server,
			Timestamp: collection.Timestamp,
			TTL:       DefaultTTL,
		}

		collection.Facts[FactMemoryUsed] = &Fact{
			Key:       FactMemoryUsed,
			Value:     memory.Used,
			Source:    string(SourceSSH),
			Server:    collection.Server,
			Timestamp: collection.Timestamp,
			TTL:       DefaultTTL,
		}

		collection.Facts[FactMemoryAvail] = &Fact{
			Key:       FactMemoryAvail,
			Value:     memory.Available,
			Source:    string(SourceSSH),
			Server:    collection.Server,
			Timestamp: collection.Timestamp,
			TTL:       DefaultTTL,
		}
	}

	// Disk information
	if dfOutput, err := c.executeCommand("df -B1 /"); err == nil {
		disk := c.parseDiskInfo(dfOutput)

		collection.Facts[FactDiskTotal] = &Fact{
			Key:       FactDiskTotal,
			Value:     disk.Total,
			Source:    string(SourceSSH),
			Server:    collection.Server,
			Timestamp: collection.Timestamp,
			TTL:       DefaultTTL,
		}

		collection.Facts[FactDiskUsed] = &Fact{
			Key:       FactDiskUsed,
			Value:     disk.Used,
			Source:    string(SourceSSH),
			Server:    collection.Server,
			Timestamp: collection.Timestamp,
			TTL:       DefaultTTL,
		}

		collection.Facts[FactDiskAvail] = &Fact{
			Key:       FactDiskAvail,
			Value:     disk.Available,
			Source:    string(SourceSSH),
			Server:    collection.Server,
			Timestamp: collection.Timestamp,
			TTL:       DefaultTTL,
		}
	}

	return nil
}

// collectNetworkFacts collects network information
func (c *SSHCollector) collectNetworkFacts(collection *FactCollection) error {
	// IP addresses
	if ipOutput, err := c.executeCommand("ip -json addr show"); err == nil {
		ips := c.parseIPAddresses(ipOutput)
		collection.Facts[FactNetworkIPs] = &Fact{
			Key:       FactNetworkIPs,
			Value:     ips,
			Source:    string(SourceSSH),
			Server:    collection.Server,
			Timestamp: collection.Timestamp,
			TTL:       DefaultTTL,
		}
	}

	// MAC addresses
	if macOutput, err := c.executeCommand("ip -json link show"); err == nil {
		macs := c.parseMACAddresses(macOutput)
		collection.Facts[FactNetworkMACs] = &Fact{
			Key:       FactNetworkMACs,
			Value:     macs,
			Source:    string(SourceSSH),
			Server:    collection.Server,
			Timestamp: collection.Timestamp,
			TTL:       DefaultTTL,
		}
	}

	// DNS configuration
	if resolvConf, err := c.executeCommand("cat /etc/resolv.conf"); err == nil {
		dns := c.parseDNSConfig(resolvConf)
		collection.Facts[FactDNS] = &Fact{
			Key:       FactDNS,
			Value:     dns,
			Source:    string(SourceSSH),
			Server:    collection.Server,
			Timestamp: collection.Timestamp,
			TTL:       DefaultTTL,
		}
	}

	return nil
}

// collectEnvironmentFacts collects environment variables
func (c *SSHCollector) collectEnvironmentFacts(collection *FactCollection) error {
	if envOutput, err := c.executeCommand("env"); err == nil {
		env := c.parseEnvironment(envOutput)
		collection.Facts[FactEnvironment] = &Fact{
			Key:       FactEnvironment,
			Value:     env,
			Source:    string(SourceSSH),
			Server:    collection.Server,
			Timestamp: collection.Timestamp,
			TTL:       DefaultTTL,
		}
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
