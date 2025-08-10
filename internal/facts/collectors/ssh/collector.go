package ssh

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	spookyfactstypes "spooky/internal/types/facts"
	spookylogging "spooky/internal/logging"
	spookyssh "spooky/internal/ssh"
	spookysshtypes "spooky/internal/ssh/types"
)

// Collector implements SSH-based fact collection
type Collector struct {
	sshManager spookyssh.SSHManager
	logger     spookylogging.Logger
}

// NewCollector creates a new SSH fact collector
func NewCollector(logger spookylogging.Logger) *Collector {
	return &Collector{
		sshManager: spookyssh.NewDefaultManager(logger),
		logger:     logger,
	}
}

// Collect gathers facts from a remote machine via SSH
func (c *Collector) Collect(server string, config *spookysshtypes.SSHConfig) (*spookyfactstypes.FactCollection, error) {
	collection := &spookyfactstypes.FactCollection{
		Server:    server,
		Facts:     make(map[string]*spookyfactstypes.Fact),
		Timestamp: time.Now(),
	}

	// Establish SSH connection
	connection, err := c.sshManager.Connect(server, config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to %s: %w", server, err)
	}
	defer c.sshManager.CloseConnection(connection)

	// Detect platform first
	platform, err := c.detectPlatform(connection)
	if err != nil {
		return nil, fmt.Errorf("failed to detect platform: %w", err)
	}

	// Collect system facts
	if err := c.collectSystemFacts(collection, connection, platform); err != nil {
		return nil, fmt.Errorf("failed to collect system facts: %w", err)
	}

	// Collect OS facts
	if err := c.collectOSFacts(collection, connection, platform); err != nil {
		return nil, fmt.Errorf("failed to collect OS facts: %w", err)
	}

	// Collect hardware facts
	if err := c.collectHardwareFacts(collection, connection, platform); err != nil {
		return nil, fmt.Errorf("failed to collect hardware facts: %w", err)
	}

	// Collect network facts
	if err := c.collectNetworkFacts(collection, connection, platform); err != nil {
		return nil, fmt.Errorf("failed to collect network facts: %w", err)
	}

	// Collect user facts
	if err := c.collectUserFacts(collection, connection, platform); err != nil {
		return nil, fmt.Errorf("failed to collect user facts: %w", err)
	}

	// Collect environment facts
	if err := c.collectEnvironmentFacts(collection, connection, platform); err != nil {
		return nil, fmt.Errorf("failed to collect environment facts: %w", err)
	}

	return collection, nil
}

// Platform represents the detected operating system
type Platform struct {
	OS      string // "linux", "darwin", "windows", "freebsd", "openbsd", "netbsd"
	Arch    string // "amd64", "arm64", "386", etc.
	Version string // OS version
}

// detectPlatform detects the operating system and architecture
func (c *Collector) detectPlatform(connection *spookysshtypes.SSHConnection) (*Platform, error) {
	// Get OS type
	osResult, err := c.sshManager.ExecuteCommand(connection, "uname -s")
	if err != nil {
		return nil, fmt.Errorf("failed to detect OS: %w", err)
	}
	os := strings.ToLower(strings.TrimSpace(osResult.Stdout))

	// Get architecture
	archResult, err := c.sshManager.ExecuteCommand(connection, "uname -m")
	if err != nil {
		return nil, fmt.Errorf("failed to detect architecture: %w", err)
	}
	arch := strings.TrimSpace(archResult.Stdout)

	// Get kernel version
	versionResult, err := c.sshManager.ExecuteCommand(connection, "uname -r")
	if err != nil {
		return nil, fmt.Errorf("failed to detect kernel version: %w", err)
	}
	version := strings.TrimSpace(versionResult.Stdout)

	// Normalize OS names
	switch os {
	case "linux":
		os = "linux"
	case "darwin":
		os = "darwin"
	case "freebsd":
		os = "freebsd"
	case "openbsd":
		os = "openbsd"
	case "netbsd":
		os = "netbsd"
	default:
		os = "unknown"
	}

	return &Platform{
		OS:      os,
		Arch:    arch,
		Version: version,
	}, nil
}

// executeCommand executes a command via SSH and returns the stdout
func (c *Collector) executeCommand(connection *spookysshtypes.SSHConnection, command string) (string, error) {
	result, err := c.sshManager.ExecuteCommand(connection, command)
	if err != nil {
		return "", fmt.Errorf("failed to execute command '%s': %w", command, err)
	}
	return strings.TrimSpace(result.Stdout), nil
}

// collectSystemFacts collects basic system information
func (c *Collector) collectSystemFacts(collection *spookyfactstypes.FactCollection, connection *spookysshtypes.SSHConnection, platform *Platform) error {
	// System type
	c.createFact(collection, "spooky_system", platform.OS)

	// Architecture
	c.createFact(collection, "spooky_architecture", platform.Arch)
	c.createFact(collection, "spooky_machine", platform.Arch)

	// Kernel version
	c.createFact(collection, "spooky_kernel", platform.Version)

	// Hostname
	hostname, err := c.executeCommand(connection, "hostname")
	if err == nil {
		c.createFact(collection, "spooky_hostname", hostname)
	}

	// FQDN
	fqdn, err := c.executeCommand(connection, "hostname -f")
	if err == nil {
		c.createFact(collection, "spooky_fqdn", fqdn)
	}

	// Domain (extract from FQDN)
	if fqdn != "" {
		parts := strings.Split(fqdn, ".")
		if len(parts) > 1 {
			domain := strings.Join(parts[1:], ".")
			c.createFact(collection, "spooky_domain", domain)
		}
	}

	return nil
}

// collectOSFacts collects operating system information
func (c *Collector) collectOSFacts(collection *spookyfactstypes.FactCollection, connection *spookysshtypes.SSHConnection, platform *Platform) error {
	switch platform.OS {
	case "linux":
		return c.collectLinuxOSFacts(collection, connection)
	case "darwin":
		return c.collectDarwinOSFacts(collection, connection)
	case "freebsd", "openbsd", "netbsd":
		return c.collectBSDOSFacts(collection, connection, platform.OS)
	default:
		return fmt.Errorf("unsupported OS: %s", platform.OS)
	}
}

// collectHardwareFacts collects hardware information
func (c *Collector) collectHardwareFacts(collection *spookyfactstypes.FactCollection, connection *spookysshtypes.SSHConnection, platform *Platform) error {
	switch platform.OS {
	case "linux":
		return c.collectLinuxHardwareFacts(collection, connection)
	case "darwin":
		return c.collectDarwinHardwareFacts(collection, connection)
	case "freebsd", "openbsd", "netbsd":
		return c.collectBSDHardwareFacts(collection, connection)
	default:
		return fmt.Errorf("unsupported OS: %s", platform.OS)
	}
}

// collectNetworkFacts collects network information
func (c *Collector) collectNetworkFacts(collection *spookyfactstypes.FactCollection, connection *spookysshtypes.SSHConnection, platform *Platform) error {
	switch platform.OS {
	case "linux":
		return c.collectLinuxNetworkFacts(collection, connection)
	case "darwin":
		return c.collectDarwinNetworkFacts(collection, connection)
	case "freebsd", "openbsd", "netbsd":
		return c.collectBSDNetworkFacts(collection, connection)
	default:
		return fmt.Errorf("unsupported OS: %s", platform.OS)
	}
}

// collectUserFacts collects user information
func (c *Collector) collectUserFacts(collection *spookyfactstypes.FactCollection, connection *spookysshtypes.SSHConnection, platform *Platform) error {
	// User ID
	userID, err := c.executeCommand(connection, "whoami")
	if err == nil {
		c.createFact(collection, "spooky_user_id", userID)
	}

	// User directory
	userDir, err := c.executeCommand(connection, "echo $HOME")
	if err == nil {
		c.createFact(collection, "spooky_user_dir", userDir)
	}

	// User shell
	userShell, err := c.executeCommand(connection, "echo $SHELL")
	if err == nil {
		c.createFact(collection, "spooky_user_shell", userShell)
	}

	return nil
}

// collectEnvironmentFacts collects environment information
func (c *Collector) collectEnvironmentFacts(collection *spookyfactstypes.FactCollection, connection *spookysshtypes.SSHConnection, platform *Platform) error {
	// Environment variables
	envOutput, err := c.executeCommand(connection, "env")
	if err == nil {
		env := c.parseEnvironment(envOutput)
		c.createFactWithValue(collection, "spooky_env", env)
	}

	return nil
}

// Platform-specific collection methods

func (c *Collector) collectLinuxOSFacts(collection *spookyfactstypes.FactCollection, connection *spookysshtypes.SSHConnection) error {
	// Read /etc/os-release
	osRelease, err := c.executeCommand(connection, "cat /etc/os-release")
	if err == nil {
		osInfo := c.parseOSRelease(osRelease)
		c.createFact(collection, "spooky_os_name", osInfo.Name)
		c.createFact(collection, "spooky_os_version", osInfo.Version)
		c.createFact(collection, "spooky_os_family", osInfo.Distribution)
		c.createFact(collection, "spooky_distribution", osInfo.Distribution)
		c.createFact(collection, "spooky_distribution_version", osInfo.Version)
	}

	return nil
}

func (c *Collector) collectDarwinOSFacts(collection *spookyfactstypes.FactCollection, connection *spookysshtypes.SSHConnection) error {
	// OS name
	c.createFact(collection, "spooky_os_name", "Darwin")
	c.createFact(collection, "spooky_os_family", "Darwin")

	// OS version
	version, err := c.executeCommand(connection, "sw_vers -productVersion")
	if err == nil {
		c.createFact(collection, "spooky_os_version", version)
		c.createFact(collection, "spooky_distribution_version", version)
	}

	// Distribution
	c.createFact(collection, "spooky_distribution", "MacOSX")

	return nil
}

func (c *Collector) collectBSDOSFacts(collection *spookyfactstypes.FactCollection, connection *spookysshtypes.SSHConnection, bsdType string) error {
	c.createFact(collection, "spooky_os_name", bsdType)
	c.createFact(collection, "spooky_os_family", "BSD")
	c.createFact(collection, "spooky_distribution", bsdType)

	// Version
	version, err := c.executeCommand(connection, "uname -r")
	if err == nil {
		c.createFact(collection, "spooky_distribution_version", version)
	}

	return nil
}

func (c *Collector) collectLinuxHardwareFacts(collection *spookyfactstypes.FactCollection, connection *spookysshtypes.SSHConnection) error {
	// CPU info
	cpuInfo, err := c.executeCommand(connection, "cat /proc/cpuinfo | grep 'model name' | head -1 | cut -d: -f2")
	if err == nil {
		c.createFactWithValue(collection, "spooky_processor", []string{strings.TrimSpace(cpuInfo)})
	}

	// CPU cores
	cpuCores, err := c.executeCommand(connection, "nproc")
	if err == nil {
		if cores, err := strconv.Atoi(cpuCores); err == nil {
			c.createFactWithValue(collection, "spooky_processor_cores", cores)
			c.createFactWithValue(collection, "spooky_processor_vcpus", cores)
		}
	}

	// Memory info
	memInfo, err := c.executeCommand(connection, "cat /proc/meminfo")
	if err == nil {
		mem := c.parseMemInfo(memInfo)
		c.createFactWithValue(collection, "spooky_memtotal_mb", mem.Total/1024/1024)
	}

	return nil
}

func (c *Collector) collectDarwinHardwareFacts(collection *spookyfactstypes.FactCollection, connection *spookysshtypes.SSHConnection) error {
	// CPU info
	cpuInfo, err := c.executeCommand(connection, "sysctl -n machdep.cpu.brand_string")
	if err == nil {
		c.createFactWithValue(collection, "spooky_processor", []string{strings.TrimSpace(cpuInfo)})
	}

	// CPU cores
	cpuCores, err := c.executeCommand(connection, "sysctl -n hw.ncpu")
	if err == nil {
		if cores, err := strconv.Atoi(cpuCores); err == nil {
			c.createFactWithValue(collection, "spooky_processor_cores", cores)
			c.createFactWithValue(collection, "spooky_processor_vcpus", cores)
		}
	}

	// Memory info
	memTotal, err := c.executeCommand(connection, "sysctl -n hw.memsize")
	if err == nil {
		if total, err := strconv.ParseUint(memTotal, 10, 64); err == nil {
			c.createFactWithValue(collection, "spooky_memtotal_mb", total/1024/1024)
		}
	}

	return nil
}

func (c *Collector) collectBSDHardwareFacts(collection *spookyfactstypes.FactCollection, connection *spookysshtypes.SSHConnection) error {
	// CPU info
	cpuInfo, err := c.executeCommand(connection, "sysctl -n hw.model")
	if err == nil {
		c.createFactWithValue(collection, "spooky_processor", []string{strings.TrimSpace(cpuInfo)})
	}

	// CPU cores
	cpuCores, err := c.executeCommand(connection, "sysctl -n hw.ncpu")
	if err == nil {
		if cores, err := strconv.Atoi(cpuCores); err == nil {
			c.createFactWithValue(collection, "spooky_processor_cores", cores)
			c.createFactWithValue(collection, "spooky_processor_vcpus", cores)
		}
	}

	// Memory info
	memTotal, err := c.executeCommand(connection, "sysctl -n hw.physmem")
	if err == nil {
		if total, err := strconv.ParseUint(memTotal, 10, 64); err == nil {
			c.createFactWithValue(collection, "spooky_memtotal_mb", total/1024/1024)
		}
	}

	return nil
}

func (c *Collector) collectLinuxNetworkFacts(collection *spookyfactstypes.FactCollection, connection *spookysshtypes.SSHConnection) error {
	// Get default IPv4
	defaultIP, err := c.executeCommand(connection, "ip route get 1.1.1.1 | grep -oP 'src \\K\\S+'")
	if err == nil {
		interfaceName, err := c.executeCommand(connection, "ip route get 1.1.1.1 | grep -oP 'dev \\K\\S+'")
		if err == nil {
			c.createFactWithValue(collection, "spooky_default_ipv4", map[string]interface{}{
				"address":   defaultIP,
				"interface": interfaceName,
			})
		}
	}

	// Get interfaces
	interfaces, err := c.executeCommand(connection, "ip link show | grep -E '^[0-9]+:' | cut -d: -f2 | tr -d ' '")
	if err == nil {
		interfaceList := strings.Split(strings.TrimSpace(interfaces), "\n")
		c.createFactWithValue(collection, "spooky_interfaces", interfaceList)
	}

	return nil
}

func (c *Collector) collectDarwinNetworkFacts(collection *spookyfactstypes.FactCollection, connection *spookysshtypes.SSHConnection) error {
	// Get default IPv4
	defaultIP, err := c.executeCommand(connection, "route -n get 1.1.1.1 | grep 'interface:' | awk '{print $2}'")
	if err == nil {
		interfaceName := strings.TrimSpace(defaultIP)
		ipAddr, err := c.executeCommand(connection, fmt.Sprintf("ifconfig %s | grep 'inet ' | awk '{print $2}'", interfaceName))
		if err == nil {
			c.createFactWithValue(collection, "spooky_default_ipv4", map[string]interface{}{
				"address":   strings.TrimSpace(ipAddr),
				"interface": interfaceName,
			})
		}
	}

	// Get interfaces
	interfaces, err := c.executeCommand(connection, "ifconfig -l")
	if err == nil {
		interfaceList := strings.Fields(strings.TrimSpace(interfaces))
		c.createFactWithValue(collection, "spooky_interfaces", interfaceList)
	}

	return nil
}

func (c *Collector) collectBSDNetworkFacts(collection *spookyfactstypes.FactCollection, connection *spookysshtypes.SSHConnection) error {
	// Get default IPv4
	defaultIP, err := c.executeCommand(connection, "route -n get 1.1.1.1 | grep 'interface:' | awk '{print $2}'")
	if err == nil {
		interfaceName := strings.TrimSpace(defaultIP)
		ipAddr, err := c.executeCommand(connection, fmt.Sprintf("ifconfig %s | grep 'inet ' | awk '{print $2}'", interfaceName))
		if err == nil {
			c.createFactWithValue(collection, "spooky_default_ipv4", map[string]interface{}{
				"address":   strings.TrimSpace(ipAddr),
				"interface": interfaceName,
			})
		}
	}

	// Get interfaces
	interfaces, err := c.executeCommand(connection, "ifconfig -l")
	if err == nil {
		interfaceList := strings.Fields(strings.TrimSpace(interfaces))
		c.createFactWithValue(collection, "spooky_interfaces", interfaceList)
	}

	return nil
}

// Helper methods

func (c *Collector) createFact(collection *spookyfactstypes.FactCollection, key, value string) {
	collection.Facts[key] = &spookyfactstypes.Fact{
		Key:       key,
		Value:     value,
		Source:    "ssh",
		Timestamp: time.Now(),
	}
}

func (c *Collector) createFactWithValue(collection *spookyfactstypes.FactCollection, key string, value interface{}) {
	collection.Facts[key] = &spookyfactstypes.Fact{
		Key:       key,
		Value:     value,
		Source:    "ssh",
		Timestamp: time.Now(),
	}
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

func (c *Collector) parseMemInfo(memInfo string) spookyfactstypes.MemoryInfo {
	lines := strings.Split(memInfo, "\n")
	mem := spookyfactstypes.MemoryInfo{}

	for _, line := range lines {
		if strings.HasPrefix(line, "MemTotal:") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				if kb, err := strconv.ParseUint(parts[1], 10, 64); err == nil {
					mem.Total = kb * 1024 // Convert KB to bytes
				}
			}
		}
	}

	return mem
}

func (c *Collector) parseEnvironment(envOutput string) map[string]string {
	lines := strings.Split(envOutput, "\n")
	env := make(map[string]string)

	for _, line := range lines {
		if idx := strings.Index(line, "="); idx > 0 {
			key := line[:idx]
			value := line[idx+1:]
			env[key] = value
		}
	}

	return env
}
