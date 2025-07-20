package facts

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
)

// LocalCollector collects facts from the local machine
type LocalCollector struct{}

// NewLocalCollector creates a new local fact collector
func NewLocalCollector() *LocalCollector {
	return &LocalCollector{}
}

// Helper function to create a fact
func (c *LocalCollector) createFact(collection *FactCollection, key, value string) {
	collection.Facts[key] = &Fact{
		Key:       key,
		Value:     value,
		Source:    string(SourceLocal),
		Server:    collection.Server,
		Timestamp: collection.Timestamp,
		TTL:       DefaultTTL,
	}
}

// Helper function to create a fact with any value type
func (c *LocalCollector) createFactWithValue(collection *FactCollection, key string, value interface{}) {
	collection.Facts[key] = &Fact{
		Key:       key,
		Value:     value,
		Source:    string(SourceLocal),
		Server:    collection.Server,
		Timestamp: collection.Timestamp,
		TTL:       DefaultTTL,
	}
}

// Collect gathers all available facts from the local machine
func (c *LocalCollector) Collect(server string) (*FactCollection, error) {
	base := NewBaseCollector()
	return base.CollectAll(server, c)
}

// CollectSpecific collects only the specified facts
func (c *LocalCollector) CollectSpecific(server string, keys []string) (*FactCollection, error) {
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
func (c *LocalCollector) GetFact(server, key string) (*Fact, error) {
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
func (c *LocalCollector) collectSystemFacts(collection *FactCollection) error {
	// Machine ID
	if machineID, err := c.readFile("/etc/machine-id"); err == nil {
		c.createFact(collection, FactMachineID, strings.TrimSpace(machineID))
	}

	// Hostname
	if hostname, err := os.Hostname(); err == nil {
		c.createFact(collection, FactHostname, hostname)
	}

	// FQDN
	if fqdn, err := c.executeCommand("hostname", "-f"); err == nil {
		c.createFact(collection, FactFQDN, strings.TrimSpace(fqdn))
	}

	return nil
}

// collectOSFacts collects operating system information
func (c *LocalCollector) collectOSFacts(collection *FactCollection) error {
	// Get host info using gopsutil
	if hostInfo, err := host.Info(); err == nil {
		c.createFact(collection, FactOSName, hostInfo.OS)

		c.createFact(collection, FactOSVersion, hostInfo.PlatformVersion)

		c.createFact(collection, FactOSDistro, hostInfo.Platform)

		c.createFact(collection, FactOSKernel, hostInfo.KernelVersion)
	}

	// Architecture
	c.createFact(collection, FactOSArch, runtime.GOARCH)

	return nil
}

// collectHardwareFacts collects hardware information
func (c *LocalCollector) collectHardwareFacts(collection *FactCollection) error {
	// CPU information
	if cpuInfo, err := cpu.Info(); err == nil && len(cpuInfo) > 0 {
		// Calculate total cores across all CPUs
		totalCores := 0
		for i := range cpuInfo {
			totalCores += int(cpuInfo[i].Cores)
		}
		c.createFactWithValue(collection, FactCPUCores, totalCores)

		c.createFact(collection, FactCPUModel, cpuInfo[0].ModelName)

		c.createFact(collection, FactCPUArch, runtime.GOARCH)

		// Get CPU frequency (current frequency in MHz)
		if freq, err := cpu.Percent(0, false); err == nil && len(freq) > 0 {
			// Get current frequency from /proc/cpuinfo as gopsutil doesn't provide it directly
			if freqStr, err := c.getCPUFrequency(); err == nil {
				c.createFact(collection, FactCPUFreq, freqStr)
			}
		}
	}

	// Memory information
	if memInfo, err := mem.VirtualMemory(); err == nil {
		c.createFactWithValue(collection, FactMemoryTotal, memInfo.Total)

		c.createFactWithValue(collection, FactMemoryUsed, memInfo.Used)

		c.createFactWithValue(collection, FactMemoryAvail, memInfo.Available)
	}

	// Disk information
	if diskInfo, err := disk.Usage("/"); err == nil {
		c.createFactWithValue(collection, FactDiskTotal, diskInfo.Total)

		c.createFactWithValue(collection, FactDiskUsed, diskInfo.Used)

		c.createFactWithValue(collection, FactDiskAvail, diskInfo.Free)
	}

	return nil
}

// collectNetworkFacts collects network information
func (c *LocalCollector) collectNetworkFacts(collection *FactCollection) error {
	// Network interfaces
	if interfaces, err := net.Interfaces(); err == nil {
		var ips []string
		var macs []string

		for _, iface := range interfaces {
			if iface.HardwareAddr != "" {
				macs = append(macs, iface.HardwareAddr)
			}
			for _, addr := range iface.Addrs {
				ips = append(ips, addr.Addr)
			}
		}

		c.createFactWithValue(collection, FactNetworkIPs, ips)

		c.createFactWithValue(collection, FactNetworkMACs, macs)
	}

	// DNS configuration
	if resolvConf, err := c.readFile("/etc/resolv.conf"); err == nil {
		dns := c.parseDNSConfig(resolvConf)
		c.createFactWithValue(collection, FactDNS, dns)
	}

	return nil
}

// collectEnvironmentFacts collects environment variables
func (c *LocalCollector) collectEnvironmentFacts(collection *FactCollection) error {
	env := make(map[string]string)
	for _, envVar := range os.Environ() {
		if parts := strings.SplitN(envVar, "=", 2); len(parts) == 2 {
			env[parts[0]] = parts[1]
		}
	}

	c.createFactWithValue(collection, FactEnvironment, env)

	return nil
}

// collectSpecificFact collects a specific fact based on the key
func (c *LocalCollector) collectSpecificFact(collection *FactCollection, key string) error {
	switch key {
	case FactMachineID, FactHostname, FactFQDN:
		return c.collectSystemFacts(collection)
	case FactOSName, FactOSVersion, FactOSDistro, FactOSArch, FactOSKernel:
		return c.collectOSFacts(collection)
	case FactCPUCores, FactCPUModel, FactCPUArch, FactCPUFreq, FactMemoryTotal, FactMemoryUsed, FactMemoryAvail, FactDiskTotal, FactDiskUsed, FactDiskAvail:
		return c.collectHardwareFacts(collection)
	case FactNetworkIPs, FactNetworkMACs, FactDNS:
		return c.collectNetworkFacts(collection)
	case FactEnvironment:
		return c.collectEnvironmentFacts(collection)
	default:
		return fmt.Errorf("unknown fact key: %s", key)
	}
}

// readFile reads a file and returns its contents
func (c *LocalCollector) readFile(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// executeCommand executes a local command
func (c *LocalCollector) executeCommand(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

// getCPUFrequency gets the current CPU frequency from /proc/cpuinfo
func (c *LocalCollector) getCPUFrequency() (string, error) {
	content, err := c.readFile("/proc/cpuinfo")
	if err != nil {
		return "", err
	}

	lines := strings.Split(content, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "cpu MHz") {
			parts := strings.Split(line, ":")
			if len(parts) == 2 {
				freq := strings.TrimSpace(parts[1])
				return freq + " MHz", nil
			}
		}
	}

	return "", fmt.Errorf("CPU frequency not found in /proc/cpuinfo")
}

// parseDNSConfig parses /etc/resolv.conf content
func (c *LocalCollector) parseDNSConfig(content string) DNSInfo {
	dns := DNSInfo{}
	lines := strings.Split(content, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "nameserver ") {
			nameserver := strings.TrimPrefix(line, "nameserver ")
			dns.Nameservers = append(dns.Nameservers, nameserver)
		} else if strings.HasPrefix(line, "search ") {
			search := strings.TrimPrefix(line, "search ")
			dns.Search = append(dns.Search, search)
		}
	}

	return dns
}
