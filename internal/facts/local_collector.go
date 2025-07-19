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

// Collect gathers all available facts from the local machine
func (c *LocalCollector) Collect(server string) (*FactCollection, error) {
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
		collection.Facts[FactMachineID] = &Fact{
			Key:       FactMachineID,
			Value:     strings.TrimSpace(machineID),
			Source:    string(SourceLocal),
			Server:    collection.Server,
			Timestamp: collection.Timestamp,
			TTL:       DefaultTTL,
		}
	}

	// Hostname
	if hostname, err := os.Hostname(); err == nil {
		collection.Facts[FactHostname] = &Fact{
			Key:       FactHostname,
			Value:     hostname,
			Source:    string(SourceLocal),
			Server:    collection.Server,
			Timestamp: collection.Timestamp,
			TTL:       DefaultTTL,
		}
	}

	// FQDN
	if fqdn, err := c.executeCommand("hostname", "-f"); err == nil {
		collection.Facts[FactFQDN] = &Fact{
			Key:       FactFQDN,
			Value:     strings.TrimSpace(fqdn),
			Source:    string(SourceLocal),
			Server:    collection.Server,
			Timestamp: collection.Timestamp,
			TTL:       DefaultTTL,
		}
	}

	return nil
}

// collectOSFacts collects operating system information
func (c *LocalCollector) collectOSFacts(collection *FactCollection) error {
	// Get host info using gopsutil
	if hostInfo, err := host.Info(); err == nil {
		collection.Facts[FactOSName] = &Fact{
			Key:       FactOSName,
			Value:     hostInfo.OS,
			Source:    string(SourceLocal),
			Server:    collection.Server,
			Timestamp: collection.Timestamp,
			TTL:       DefaultTTL,
		}

		collection.Facts[FactOSVersion] = &Fact{
			Key:       FactOSVersion,
			Value:     hostInfo.PlatformVersion,
			Source:    string(SourceLocal),
			Server:    collection.Server,
			Timestamp: collection.Timestamp,
			TTL:       DefaultTTL,
		}

		collection.Facts[FactOSDistro] = &Fact{
			Key:       FactOSDistro,
			Value:     hostInfo.Platform,
			Source:    string(SourceLocal),
			Server:    collection.Server,
			Timestamp: collection.Timestamp,
			TTL:       DefaultTTL,
		}

		collection.Facts[FactOSKernel] = &Fact{
			Key:       FactOSKernel,
			Value:     hostInfo.KernelVersion,
			Source:    string(SourceLocal),
			Server:    collection.Server,
			Timestamp: collection.Timestamp,
			TTL:       DefaultTTL,
		}
	}

	// Architecture
	collection.Facts[FactOSArch] = &Fact{
		Key:       FactOSArch,
		Value:     runtime.GOARCH,
		Source:    string(SourceLocal),
		Server:    collection.Server,
		Timestamp: collection.Timestamp,
		TTL:       DefaultTTL,
	}

	return nil
}

// collectHardwareFacts collects hardware information
func (c *LocalCollector) collectHardwareFacts(collection *FactCollection) error {
	// CPU information
	if cpuInfo, err := cpu.Info(); err == nil && len(cpuInfo) > 0 {
		collection.Facts[FactCPUCores] = &Fact{
			Key:       FactCPUCores,
			Value:     int(cpuInfo[0].Cores),
			Source:    string(SourceLocal),
			Server:    collection.Server,
			Timestamp: collection.Timestamp,
			TTL:       DefaultTTL,
		}

		collection.Facts[FactCPUModel] = &Fact{
			Key:       FactCPUModel,
			Value:     cpuInfo[0].ModelName,
			Source:    string(SourceLocal),
			Server:    collection.Server,
			Timestamp: collection.Timestamp,
			TTL:       DefaultTTL,
		}

		collection.Facts[FactCPUArch] = &Fact{
			Key:       FactCPUArch,
			Value:     runtime.GOARCH,
			Source:    string(SourceLocal),
			Server:    collection.Server,
			Timestamp: collection.Timestamp,
			TTL:       DefaultTTL,
		}
	}

	// Memory information
	if memInfo, err := mem.VirtualMemory(); err == nil {
		collection.Facts[FactMemoryTotal] = &Fact{
			Key:       FactMemoryTotal,
			Value:     memInfo.Total,
			Source:    string(SourceLocal),
			Server:    collection.Server,
			Timestamp: collection.Timestamp,
			TTL:       DefaultTTL,
		}

		collection.Facts[FactMemoryUsed] = &Fact{
			Key:       FactMemoryUsed,
			Value:     memInfo.Used,
			Source:    string(SourceLocal),
			Server:    collection.Server,
			Timestamp: collection.Timestamp,
			TTL:       DefaultTTL,
		}

		collection.Facts[FactMemoryAvail] = &Fact{
			Key:       FactMemoryAvail,
			Value:     memInfo.Available,
			Source:    string(SourceLocal),
			Server:    collection.Server,
			Timestamp: collection.Timestamp,
			TTL:       DefaultTTL,
		}
	}

	// Disk information
	if diskInfo, err := disk.Usage("/"); err == nil {
		collection.Facts[FactDiskTotal] = &Fact{
			Key:       FactDiskTotal,
			Value:     diskInfo.Total,
			Source:    string(SourceLocal),
			Server:    collection.Server,
			Timestamp: collection.Timestamp,
			TTL:       DefaultTTL,
		}

		collection.Facts[FactDiskUsed] = &Fact{
			Key:       FactDiskUsed,
			Value:     diskInfo.Used,
			Source:    string(SourceLocal),
			Server:    collection.Server,
			Timestamp: collection.Timestamp,
			TTL:       DefaultTTL,
		}

		collection.Facts[FactDiskAvail] = &Fact{
			Key:       FactDiskAvail,
			Value:     diskInfo.Free,
			Source:    string(SourceLocal),
			Server:    collection.Server,
			Timestamp: collection.Timestamp,
			TTL:       DefaultTTL,
		}
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

		collection.Facts[FactNetworkIPs] = &Fact{
			Key:       FactNetworkIPs,
			Value:     ips,
			Source:    string(SourceLocal),
			Server:    collection.Server,
			Timestamp: collection.Timestamp,
			TTL:       DefaultTTL,
		}

		collection.Facts[FactNetworkMACs] = &Fact{
			Key:       FactNetworkMACs,
			Value:     macs,
			Source:    string(SourceLocal),
			Server:    collection.Server,
			Timestamp: collection.Timestamp,
			TTL:       DefaultTTL,
		}
	}

	// DNS configuration
	if resolvConf, err := c.readFile("/etc/resolv.conf"); err == nil {
		dns := c.parseDNSConfig(resolvConf)
		collection.Facts[FactDNS] = &Fact{
			Key:       FactDNS,
			Value:     dns,
			Source:    string(SourceLocal),
			Server:    collection.Server,
			Timestamp: collection.Timestamp,
			TTL:       DefaultTTL,
		}
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

	collection.Facts[FactEnvironment] = &Fact{
		Key:       FactEnvironment,
		Value:     env,
		Source:    string(SourceLocal),
		Server:    collection.Server,
		Timestamp: collection.Timestamp,
		TTL:       DefaultTTL,
	}

	return nil
}

// collectSpecificFact collects a specific fact based on the key
func (c *LocalCollector) collectSpecificFact(collection *FactCollection, key string) error {
	switch key {
	case FactMachineID, FactHostname, FactFQDN:
		return c.collectSystemFacts(collection)
	case FactOSName, FactOSVersion, FactOSDistro, FactOSArch, FactOSKernel:
		return c.collectOSFacts(collection)
	case FactCPUCores, FactCPUModel, FactCPUArch, FactMemoryTotal, FactMemoryUsed, FactMemoryAvail, FactDiskTotal, FactDiskUsed, FactDiskAvail:
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
