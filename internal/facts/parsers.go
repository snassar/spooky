package facts

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"spooky/internal/schemas"
)

// ============================================================================
// COMMAND OUTPUT PARSERS
// ============================================================================

// SystemFactsParser parses system and architecture related commands
type SystemFactsParser struct{}

// ParseSystemFacts parses system commands and returns structured facts
func (p *SystemFactsParser) ParseSystemFacts() map[string]*schemas.FactV1 {
	facts := make(map[string]*schemas.FactV1)

	// Parse uname information
	if uname, err := p.parseUname(); err == nil {
		for k, v := range uname {
			facts[k] = v
		}
	}

	// Parse hostname information
	if hostname, err := p.parseHostname(); err == nil {
		for k, v := range hostname {
			facts[k] = v
		}
	}

	// Parse virtualization information
	if virt, err := p.parseVirtualization(); err == nil {
		for k, v := range virt {
			facts[k] = v
		}
	}

	return facts
}

// parseUname parses uname command output
func (p *SystemFactsParser) parseUname() (map[string]*schemas.FactV1, error) {
	facts := make(map[string]*schemas.FactV1)

	// Read /proc/version for kernel info
	if data, err := os.ReadFile("/proc/version"); err == nil {
		version := strings.TrimSpace(string(data))
		facts["kernel_version"] = &schemas.FactV1{
			Value:       version,
			Type:        "string",
			Description: "Kernel version and build information",
		}
	}

	// Read /proc/sys/kernel/ostype for OS type
	if data, err := os.ReadFile("/proc/sys/kernel/ostype"); err == nil {
		osType := strings.TrimSpace(string(data))
		facts["system"] = &schemas.FactV1{
			Value:       osType,
			Type:        "string",
			Description: "Operating system type",
		}
	}

	// Read /proc/sys/kernel/arch for architecture
	if data, err := os.ReadFile("/proc/sys/kernel/arch"); err == nil {
		arch := strings.TrimSpace(string(data))
		facts["architecture"] = &schemas.FactV1{
			Value:       arch,
			Type:        "string",
			Description: "System architecture",
		}
		facts["machine"] = &schemas.FactV1{
			Value:       arch,
			Type:        "string",
			Description: "Machine architecture",
		}
	}

	return facts, nil
}

// parseHostname parses hostname related information
func (p *SystemFactsParser) parseHostname() (map[string]*schemas.FactV1, error) {
	facts := make(map[string]*schemas.FactV1)

	// Read hostname
	if hostname, err := os.ReadFile("/proc/sys/kernel/hostname"); err == nil {
		host := strings.TrimSpace(string(hostname))
		facts["hostname"] = &schemas.FactV1{
			Value:       host,
			Type:        "string",
			Description: "System hostname",
		}
		facts["nodename"] = &schemas.FactV1{
			Value:       host,
			Type:        "string",
			Description: "System nodename",
		}
	}

	// Read domainname
	if domain, err := os.ReadFile("/proc/sys/kernel/domainname"); err == nil {
		domainStr := strings.TrimSpace(string(domain))
		if domainStr != "(none)" {
			facts["domain"] = &schemas.FactV1{
				Value:       domainStr,
				Type:        "string",
				Description: "System domain name",
			}
		}
	}

	// Construct FQDN
	if hostname, exists := facts["hostname"]; exists {
		if domain, exists := facts["domain"]; exists {
			fqdn := fmt.Sprintf("%s.%s", hostname.Value, domain.Value)
			facts["fqdn"] = &schemas.FactV1{
				Value:       fqdn,
				Type:        "string",
				Description: "Fully qualified domain name",
			}
		} else {
			facts["fqdn"] = facts["hostname"]
		}
	}

	return facts, nil
}

// parseVirtualization parses virtualization information
func (p *SystemFactsParser) parseVirtualization() (map[string]*schemas.FactV1, error) {
	facts := make(map[string]*schemas.FactV1)

	// Check for common virtualization indicators
	virtType := "physical"
	virtRole := "host"

	// Check /proc/cpuinfo for hypervisor info
	if data, err := os.ReadFile("/proc/cpuinfo"); err == nil {
		scanner := bufio.NewScanner(strings.NewReader(string(data)))
		for scanner.Scan() {
			line := scanner.Text()
			if strings.Contains(line, "hypervisor") {
				virtType = "kvm"
				virtRole = "guest"
				break
			}
		}
	}

	// Check for Xen
	if _, err := os.Stat("/proc/xen"); err == nil {
		virtType = "xen"
		virtRole = "guest"
	}

	// Check for VMware
	if data, err := os.ReadFile("/proc/scsi/scsi"); err == nil {
		if strings.Contains(string(data), "VMware") {
			virtType = "vmware"
			virtRole = "guest"
		}
	}

	// Check for VirtualBox
	if data, err := os.ReadFile("/proc/scsi/scsi"); err == nil {
		if strings.Contains(string(data), "VBOX") {
			virtType = "virtualbox"
			virtRole = "guest"
		}
	}

	facts["virtualization_type"] = &schemas.FactV1{
		Value:       virtType,
		Type:        "string",
		Description: "Virtualization type",
	}

	facts["virtualization_role"] = &schemas.FactV1{
		Value:       virtRole,
		Type:        "string",
		Description: "Virtualization role (host/guest)",
	}

	return facts, nil
}

// HardwareFactsParser parses hardware related commands
type HardwareFactsParser struct{}

// ParseHardwareFacts parses hardware commands and returns structured facts
func (p *HardwareFactsParser) ParseHardwareFacts() map[string]*schemas.FactV1 {
	facts := make(map[string]*schemas.FactV1)

	// Parse CPU information
	if cpu, err := p.parseCPUInfo(); err == nil {
		for k, v := range cpu {
			facts[k] = v
		}
	}

	// Parse memory information
	if memory, err := p.parseMemoryInfo(); err == nil {
		for k, v := range memory {
			facts[k] = v
		}
	}

	// Parse storage information
	if storage, err := p.parseStorageInfo(); err == nil {
		for k, v := range storage {
			facts[k] = v
		}
	}

	// Parse load averages
	if load, err := p.parseLoadAverage(); err == nil {
		for k, v := range load {
			facts[k] = v
		}
	}

	// Parse BIOS and hardware information
	if bios, err := p.parseBIOSInfo(); err == nil {
		for k, v := range bios {
			facts[k] = v
		}
	}

	return facts
}

// parseCPUInfo parses /proc/cpuinfo
func (p *HardwareFactsParser) parseCPUInfo() (map[string]*schemas.FactV1, error) {
	facts := make(map[string]*schemas.FactV1)

	data, err := os.ReadFile("/proc/cpuinfo")
	if err != nil {
		return facts, err
	}

	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	processorCount := 0
	physicalCores := 0
	threadsPerCore := 1
	var processorInfo []string

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "processor") {
			processorCount++
		} else if strings.HasPrefix(line, "cpu cores") {
			if parts := strings.Split(line, ":"); len(parts) == 2 {
				if cores, err := strconv.Atoi(strings.TrimSpace(parts[1])); err == nil {
					physicalCores = cores
				}
			}
		} else if strings.HasPrefix(line, "model name") {
			if parts := strings.Split(line, ":"); len(parts) == 2 {
				model := strings.TrimSpace(parts[1])
				processorInfo = append(processorInfo, model)
			}
		}
	}

	// Calculate derived values
	if physicalCores > 0 && processorCount > 0 {
		threadsPerCore = processorCount / physicalCores
	}

	facts["processor_count"] = &schemas.FactV1{
		Value:       physicalCores,
		Type:        "number",
		Description: "Number of physical CPU cores",
	}

	facts["processor_nproc"] = &schemas.FactV1{
		Value:       processorCount,
		Type:        "number",
		Description: "Number of logical processors",
	}

	facts["processor_threads_per_core"] = &schemas.FactV1{
		Value:       threadsPerCore,
		Type:        "number",
		Description: "Number of threads per CPU core",
	}

	facts["processor_vcpus"] = &schemas.FactV1{
		Value:       processorCount,
		Type:        "number",
		Description: "Total number of virtual CPUs",
	}

	if len(processorInfo) > 0 {
		facts["processor"] = &schemas.FactV1{
			Value:       processorInfo,
			Type:        "array",
			Description: "Array of processor model names",
		}
	}

	return facts, nil
}

// parseMemoryInfo parses /proc/meminfo
func (p *HardwareFactsParser) parseMemoryInfo() (map[string]*schemas.FactV1, error) {
	facts := make(map[string]*schemas.FactV1)

	data, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		return facts, err
	}

	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}

		key := strings.TrimSuffix(parts[0], ":")
		value, err := strconv.ParseUint(parts[1], 10, 64)
		if err != nil {
			continue
		}

		// Convert KB to MB for consistency with Ansible
		valueMB := value / 1024

		switch key {
		case "MemTotal":
			facts["memory_total_mb"] = &schemas.FactV1{
				Value:       valueMB,
				Type:        "number",
				Description: "Total memory in MB",
			}
		case "MemAvailable":
			facts["memory_free_mb"] = &schemas.FactV1{
				Value:       valueMB,
				Type:        "number",
				Description: "Available memory in MB",
			}
		case "SwapTotal":
			facts["swap_total_mb"] = &schemas.FactV1{
				Value:       valueMB,
				Type:        "number",
				Description: "Total swap space in MB",
			}
		case "SwapFree":
			facts["swap_free_mb"] = &schemas.FactV1{
				Value:       valueMB,
				Type:        "number",
				Description: "Free swap space in MB",
			}
		}
	}

	// Calculate used memory
	if total, exists := facts["memory_total_mb"]; exists {
		if free, exists := facts["memory_free_mb"]; exists {
			if totalVal, ok := total.Value.(uint64); ok {
				if freeVal, ok := free.Value.(uint64); ok {
					used := totalVal - freeVal
					facts["memory_used_mb"] = &schemas.FactV1{
						Value:       used,
						Type:        "number",
						Description: "Used memory in MB",
					}
				}
			}
		}
	}

	return facts, nil
}

// parseStorageInfo parses storage device information
func (p *HardwareFactsParser) parseStorageInfo() (map[string]*schemas.FactV1, error) {
	facts := make(map[string]*schemas.FactV1)

	// Parse /proc/partitions for basic partition info
	if data, err := os.ReadFile("/proc/partitions"); err == nil {
		scanner := bufio.NewScanner(strings.NewReader(string(data)))
		// Skip header lines
		for i := 0; i < 2; i++ {
			scanner.Scan()
		}

		var devices []string
		for scanner.Scan() {
			line := scanner.Text()
			parts := strings.Fields(line)
			if len(parts) >= 4 {
				device := parts[3]
				if strings.HasPrefix(device, "sd") || strings.HasPrefix(device, "hd") ||
					strings.HasPrefix(device, "nvme") || strings.HasPrefix(device, "vd") {
					devices = append(devices, device)
				}
			}
		}

		if len(devices) > 0 {
			facts["devices"] = &schemas.FactV1{
				Value:       devices,
				Type:        "array",
				Description: "List of storage devices",
			}
		}
	}

	// Parse /proc/mounts for filesystem information
	if data, err := os.ReadFile("/proc/mounts"); err == nil {
		scanner := bufio.NewScanner(strings.NewReader(string(data)))
		var mounts []map[string]interface{}

		for scanner.Scan() {
			line := scanner.Text()
			parts := strings.Fields(line)
			if len(parts) >= 4 {
				mount := map[string]interface{}{
					"device":     parts[0],
					"mountpoint": parts[1],
					"fstype":     parts[2],
					"options":    parts[3],
				}
				mounts = append(mounts, mount)
			}
		}

		if len(mounts) > 0 {
			facts["mounts"] = &schemas.FactV1{
				Value:       mounts,
				Type:        "array",
				Description: "Mounted filesystems",
			}
		}
	}

	return facts, nil
}

// parseBIOSInfo parses BIOS and hardware information
func (p *HardwareFactsParser) parseBIOSInfo() (map[string]*schemas.FactV1, error) {
	facts := make(map[string]*schemas.FactV1)

	// Try to read BIOS information from /sys/class/dmi/id/
	biosPaths := map[string]string{
		"bios_date":       "/sys/class/dmi/id/bios_date",
		"bios_version":    "/sys/class/dmi/id/bios_version",
		"product_name":    "/sys/class/dmi/id/product_name",
		"product_serial":  "/sys/class/dmi/id/product_serial",
		"product_uuid":    "/sys/class/dmi/id/product_uuid",
		"product_version": "/sys/class/dmi/id/product_version",
		"system_vendor":   "/sys/class/dmi/id/sys_vendor",
		"machine_id":      "/etc/machine-id",
	}

	for factName, path := range biosPaths {
		if data, err := os.ReadFile(path); err == nil {
			value := strings.TrimSpace(string(data))
			if value != "" && value != "None" && value != "To be filled by O.E.M." {
				facts[factName] = &schemas.FactV1{
					Value:       value,
					Type:        "string",
					Description: fmt.Sprintf("BIOS/Hardware fact: %s", factName),
				}
			}
		}
	}

	// Detect form factor
	formFactor := p.detectFormFactor()
	if formFactor != "" {
		facts["form_factor"] = &schemas.FactV1{
			Value:       formFactor,
			Type:        "string",
			Description: "System form factor",
		}
	}

	return facts, nil
}

// detectFormFactor attempts to determine the system form factor
func (p *HardwareFactsParser) detectFormFactor() string {
	// Check for common form factor indicators
	if data, err := os.ReadFile("/sys/class/dmi/id/chassis_type"); err == nil {
		chassisType := strings.TrimSpace(string(data))
		switch chassisType {
		case "1", "2", "3", "4", "5", "6", "7":
			return "desktop"
		case "8", "9", "10", "11", "12", "13", "14":
			return "laptop"
		case "15", "16", "17", "18", "19", "20", "21", "22", "23":
			return "server"
		case "30", "31", "32", "33", "34":
			return "tablet"
		}
	}

	// Fallback detection based on system characteristics
	if _, err := os.Stat("/sys/class/power_supply/BAT0"); err == nil {
		return "laptop" // Has battery
	}

	// Check if it's a server by looking for server-specific indicators
	if data, err := os.ReadFile("/proc/cpuinfo"); err == nil {
		if strings.Contains(string(data), "Xeon") || strings.Contains(string(data), "EPYC") {
			return "server"
		}
	}

	return "desktop" // Default fallback
}

// parseLoadAverage parses /proc/loadavg
func (p *HardwareFactsParser) parseLoadAverage() (map[string]*schemas.FactV1, error) {
	facts := make(map[string]*schemas.FactV1)

	data, err := os.ReadFile("/proc/loadavg")
	if err != nil {
		return facts, err
	}

	parts := strings.Fields(strings.TrimSpace(string(data)))
	if len(parts) >= 3 {
		if load1, err := strconv.ParseFloat(parts[0], 64); err == nil {
			facts["load_1"] = &schemas.FactV1{
				Value:       load1,
				Type:        "number",
				Description: "1-minute load average",
			}
		}
		if load5, err := strconv.ParseFloat(parts[1], 64); err == nil {
			facts["load_5"] = &schemas.FactV1{
				Value:       load5,
				Type:        "number",
				Description: "5-minute load average",
			}
		}
		if load15, err := strconv.ParseFloat(parts[2], 64); err == nil {
			facts["load_15"] = &schemas.FactV1{
				Value:       load15,
				Type:        "number",
				Description: "15-minute load average",
			}
		}
	}

	return facts, nil
}

// NetworkFactsParser parses network related commands
type NetworkFactsParser struct{}

// ParseNetworkFacts parses network commands and returns structured facts
func (p *NetworkFactsParser) ParseNetworkFacts() map[string]*schemas.FactV1 {
	facts := make(map[string]*schemas.FactV1)

	// Parse network interfaces
	if interfaces, err := p.parseNetworkInterfaces(); err == nil {
		for k, v := range interfaces {
			facts[k] = v
		}
	}

	// Parse routing information
	if routing, err := p.parseRoutingInfo(); err == nil {
		for k, v := range routing {
			facts[k] = v
		}
	}

	// Parse SSH host keys
	if sshKeys, err := p.parseSSHHostKeys(); err == nil {
		for k, v := range sshKeys {
			facts[k] = v
		}
	}

	return facts
}

// parseNetworkInterfaces parses network interface information
func (p *NetworkFactsParser) parseNetworkInterfaces() (map[string]*schemas.FactV1, error) {
	facts := make(map[string]*schemas.FactV1)

	// Read /proc/net/dev for interface list
	if data, err := os.ReadFile("/proc/net/dev"); err == nil {
		scanner := bufio.NewScanner(strings.NewReader(string(data)))
		// Skip header lines
		for i := 0; i < 2; i++ {
			scanner.Scan()
		}

		var interfaces []string
		var allIPv4 []string
		var allIPv6 []string

		for scanner.Scan() {
			line := scanner.Text()
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				interfaceName := strings.TrimSuffix(parts[0], ":")
				if interfaceName != "lo" { // Skip loopback
					interfaces = append(interfaces, interfaceName)

					// Get IP addresses for this interface
					if ipv4, err := p.getInterfaceIPv4(interfaceName); err == nil {
						allIPv4 = append(allIPv4, ipv4...)
					}
					if ipv6, err := p.getInterfaceIPv6(interfaceName); err == nil {
						allIPv6 = append(allIPv6, ipv6...)
					}
				}
			}
		}

		if len(interfaces) > 0 {
			facts["interfaces"] = &schemas.FactV1{
				Value:       interfaces,
				Type:        "array",
				Description: "List of network interfaces",
			}
		}

		if len(allIPv4) > 0 {
			facts["all_ipv4_addresses"] = &schemas.FactV1{
				Value:       allIPv4,
				Type:        "array",
				Description: "All IPv4 addresses",
			}
		}

		if len(allIPv6) > 0 {
			facts["all_ipv6_addresses"] = &schemas.FactV1{
				Value:       allIPv6,
				Type:        "array",
				Description: "All IPv6 addresses",
			}
		}
	}

	return facts, nil
}

// getInterfaceIPv4 gets IPv4 addresses for a specific interface
func (p *NetworkFactsParser) getInterfaceIPv4(interfaceName string) ([]string, error) {
	var ips []string

	// Read /proc/net/fib_trie for IPv4 addresses
	if data, err := os.ReadFile("/proc/net/fib_trie"); err == nil {
		scanner := bufio.NewScanner(strings.NewReader(string(data)))
		for scanner.Scan() {
			line := scanner.Text()
			if strings.Contains(line, interfaceName) && strings.Contains(line, "/32") {
				// Extract IP address from line
				if ip := p.extractIPFromLine(line); ip != "" {
					ips = append(ips, ip)
				}
			}
		}
	}

	return ips, nil
}

// getInterfaceIPv6 gets IPv6 addresses for a specific interface
func (p *NetworkFactsParser) getInterfaceIPv6(interfaceName string) ([]string, error) {
	var ips []string

	// Read /proc/net/if_inet6 for IPv6 addresses
	if data, err := os.ReadFile("/proc/net/if_inet6"); err == nil {
		scanner := bufio.NewScanner(strings.NewReader(string(data)))
		for scanner.Scan() {
			line := scanner.Text()
			parts := strings.Fields(line)
			if len(parts) >= 6 && parts[5] == interfaceName {
				// Convert hex IPv6 to readable format
				if ip := p.hexToIPv6(parts[0]); ip != "" {
					ips = append(ips, ip)
				}
			}
		}
	}

	return ips, nil
}

// extractIPFromLine extracts IP address from a line
func (p *NetworkFactsParser) extractIPFromLine(line string) string {
	// Simple regex-like extraction for IPv4
	parts := strings.Fields(line)
	for _, part := range parts {
		if strings.Contains(part, ".") && strings.Count(part, ".") == 3 {
			return part
		}
	}
	return ""
}

// hexToIPv6 converts hex IPv6 to readable format
func (p *NetworkFactsParser) hexToIPv6(hex string) string {
	// Simple conversion - in production you'd want a proper IPv6 parser
	if len(hex) == 32 {
		// Convert 32 hex chars to IPv6 format
		segments := make([]string, 8)
		for i := 0; i < 8; i++ {
			segments[i] = hex[i*4 : i*4+4]
		}
		return strings.Join(segments, ":")
	}
	return ""
}

// parseRoutingInfo parses routing information
func (p *NetworkFactsParser) parseRoutingInfo() (map[string]*schemas.FactV1, error) {
	facts := make(map[string]*schemas.FactV1)

	// Read /proc/net/route for routing info
	if data, err := os.ReadFile("/proc/net/route"); err == nil {
		scanner := bufio.NewScanner(strings.NewReader(string(data)))
		// Skip header line
		scanner.Scan()

		for scanner.Scan() {
			line := scanner.Text()
			parts := strings.Fields(line)
			if len(parts) >= 8 {
				destination := parts[1]
				gateway := parts[2]
				iface := parts[0]

				// Look for default route (destination 00000000)
				if destination == "00000000" && gateway != "00000000" {
					facts["default_ipv4"] = &schemas.FactV1{
						Value: map[string]interface{}{
							"gateway":   gateway,
							"interface": iface,
						},
						Type:        "object",
						Description: "Default IPv4 route information",
					}
					break
				}
			}
		}
	}

	return facts, nil
}

// parseSSHHostKeys parses SSH host keys
func (p *NetworkFactsParser) parseSSHHostKeys() (map[string]*schemas.FactV1, error) {
	facts := make(map[string]*schemas.FactV1)

	// SSH host key paths
	sshKeyPaths := map[string]string{
		"ssh_host_key_rsa_public":     "/etc/ssh/ssh_host_rsa_key.pub",
		"ssh_host_key_ecdsa_public":   "/etc/ssh/ssh_host_ecdsa_key.pub",
		"ssh_host_key_ed25519_public": "/etc/ssh/ssh_host_ed25519_key.pub",
	}

	for factName, path := range sshKeyPaths {
		if data, err := os.ReadFile(path); err == nil {
			keyContent := strings.TrimSpace(string(data))
			if keyContent != "" {
				facts[factName] = &schemas.FactV1{
					Value:       keyContent,
					Type:        "string",
					Description: fmt.Sprintf("SSH host key: %s", factName),
				}
			}
		}
	}

	return facts, nil
}

// OSFactsParser parses operating system related commands
type OSFactsParser struct{}

// ParseOSFacts parses OS commands and returns structured facts
func (p *OSFactsParser) ParseOSFacts() map[string]*schemas.FactV1 {
	facts := make(map[string]*schemas.FactV1)

	// Parse OS release information
	if osInfo, err := p.parseOSRelease(); err == nil {
		for k, v := range osInfo {
			facts[k] = v
		}
	}

	// Parse package manager information
	if pkgMgr, err := p.parsePackageManager(); err == nil {
		for k, v := range pkgMgr {
			facts[k] = v
		}
	}

	// Parse security module information
	if security, err := p.parseSecurityModules(); err == nil {
		for k, v := range security {
			facts[k] = v
		}
	}

	return facts
}

// parseOSRelease parses /etc/os-release
func (p *OSFactsParser) parseOSRelease() (map[string]*schemas.FactV1, error) {
	facts := make(map[string]*schemas.FactV1)

	// Try /etc/os-release first
	data, err := os.ReadFile("/etc/os-release")
	if err != nil {
		// Fallback to /etc/redhat-release or /etc/debian_version
		if data, err = os.ReadFile("/etc/redhat-release"); err == nil {
			facts["distribution"] = &schemas.FactV1{
				Value:       "RedHat",
				Type:        "string",
				Description: "Distribution name",
			}
			facts["os_family"] = &schemas.FactV1{
				Value:       "RedHat",
				Type:        "string",
				Description: "OS family",
			}
			return facts, nil
		}
		if data, err = os.ReadFile("/etc/debian_version"); err == nil {
			content := strings.TrimSpace(string(data))
			facts["distribution"] = &schemas.FactV1{
				Value:       "Debian",
				Type:        "string",
				Description: "Distribution name",
			}
			facts["os_family"] = &schemas.FactV1{
				Value:       "Debian",
				Type:        "string",
				Description: "OS family",
			}
			facts["distribution_version"] = &schemas.FactV1{
				Value:       content,
				Type:        "string",
				Description: "Distribution version",
			}
			return facts, nil
		}
		return facts, err
	}

	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "=") {
			parts := strings.SplitN(line, "=", 2)
			key := parts[0]
			value := strings.Trim(parts[1], "\"")

			switch key {
			case "ID":
				facts["distribution"] = &schemas.FactV1{
					Value:       value,
					Type:        "string",
					Description: "Distribution identifier",
				}
			case "VERSION_ID":
				facts["distribution_version"] = &schemas.FactV1{
					Value:       value,
					Type:        "string",
					Description: "Distribution version",
				}
			case "ID_LIKE":
				facts["os_family"] = &schemas.FactV1{
					Value:       value,
					Type:        "string",
					Description: "OS family",
				}
			}
		}
	}

	return facts, nil
}

// parsePackageManager parses package manager information
func (p *OSFactsParser) parsePackageManager() (map[string]*schemas.FactV1, error) {
	facts := make(map[string]*schemas.FactV1)

	// Check for common package managers
	packageManagers := []string{"apt", "yum", "dnf", "pacman", "zypper", "emerge"}
	var foundManager string

	for _, manager := range packageManagers {
		if _, err := os.Stat(fmt.Sprintf("/usr/bin/%s", manager)); err == nil {
			foundManager = manager
			break
		}
	}

	if foundManager == "" {
		foundManager = "unknown"
	}

	facts["pkg_mgr"] = &schemas.FactV1{
		Value:       foundManager,
		Type:        "string",
		Description: "Package manager",
	}

	return facts, nil
}

// parseSecurityModules parses security module information
func (p *OSFactsParser) parseSecurityModules() (map[string]*schemas.FactV1, error) {
	facts := make(map[string]*schemas.FactV1)

	// Check for SELinux
	if _, err := os.Stat("/etc/selinux/config"); err == nil {
		if data, err := os.ReadFile("/etc/selinux/config"); err == nil {
			content := string(data)
			scanner := bufio.NewScanner(strings.NewReader(content))
			selinuxInfo := make(map[string]interface{})

			for scanner.Scan() {
				line := scanner.Text()
				if strings.HasPrefix(line, "SELINUX=") {
					mode := strings.TrimPrefix(line, "SELINUX=")
					selinuxInfo["status"] = mode
					selinuxInfo["mode"] = mode
				} else if strings.HasPrefix(line, "SELINUXTYPE=") {
					selinuxType := strings.TrimPrefix(line, "SELINUXTYPE=")
					selinuxInfo["type"] = selinuxType
				}
			}

			if len(selinuxInfo) > 0 {
				facts["selinux"] = &schemas.FactV1{
					Value:       selinuxInfo,
					Type:        "object",
					Description: "SELinux status and configuration",
				}
			}
		}
	}

	// Check for AppArmor
	if _, err := os.Stat("/sys/kernel/security/apparmor"); err == nil {
		facts["apparmor"] = &schemas.FactV1{
			Value:       "enabled",
			Type:        "string",
			Description: "AppArmor status",
		}
	}

	// Check for FIPS mode
	if data, err := os.ReadFile("/proc/sys/crypto/fips_enabled"); err == nil {
		fipsEnabled := strings.TrimSpace(string(data))
		if fipsEnabled == "1" {
			facts["fips"] = &schemas.FactV1{
				Value:       true,
				Type:        "boolean",
				Description: "FIPS mode enabled",
			}
		}
	}

	return facts, nil
}

// UserFactsParser parses user and environment related commands
type UserFactsParser struct{}

// ParseUserFacts parses user commands and returns structured facts
func (p *UserFactsParser) ParseUserFacts() map[string]*schemas.FactV1 {
	facts := make(map[string]*schemas.FactV1)

	// Parse current user information
	if userInfo, err := p.parseUserInfo(); err == nil {
		for k, v := range userInfo {
			facts[k] = v
		}
	}

	// Parse environment variables
	if envVars, err := p.parseEnvironment(); err == nil {
		for k, v := range envVars {
			facts[k] = v
		}
	}

	return facts
}

// parseUserInfo parses current user information
func (p *UserFactsParser) parseUserInfo() (map[string]*schemas.FactV1, error) {
	facts := make(map[string]*schemas.FactV1)

	// Get current user ID
	if uid := os.Getuid(); uid != -1 {
		facts["user_uid"] = &schemas.FactV1{
			Value:       uid,
			Type:        "number",
			Description: "Current user UID",
		}
	}

	// Get current group ID
	if gid := os.Getgid(); gid != -1 {
		facts["user_gid"] = &schemas.FactV1{
			Value:       gid,
			Type:        "number",
			Description: "Current user GID",
		}
	}

	// Get home directory
	if home := os.Getenv("HOME"); home != "" {
		facts["user_dir"] = &schemas.FactV1{
			Value:       home,
			Type:        "string",
			Description: "User home directory",
		}
	}

	// Get shell
	if shell := os.Getenv("SHELL"); shell != "" {
		facts["user_shell"] = &schemas.FactV1{
			Value:       shell,
			Type:        "string",
			Description: "User shell",
		}
	}

	return facts, nil
}

// parseEnvironment parses environment variables
func (p *UserFactsParser) parseEnvironment() (map[string]*schemas.FactV1, error) {
	facts := make(map[string]*schemas.FactV1)

	// Get important environment variables
	envVars := map[string]string{
		"PATH":     os.Getenv("PATH"),
		"LANG":     os.Getenv("LANG"),
		"TERM":     os.Getenv("TERM"),
		"USER":     os.Getenv("USER"),
		"LOGNAME":  os.Getenv("LOGNAME"),
		"HOSTNAME": os.Getenv("HOSTNAME"),
	}

	// Filter out empty values
	filteredEnv := make(map[string]string)
	for k, v := range envVars {
		if v != "" {
			filteredEnv[k] = v
		}
	}

	if len(filteredEnv) > 0 {
		facts["env"] = &schemas.FactV1{
			Value:       filteredEnv,
			Type:        "object",
			Description: "Environment variables",
		}
	}

	return facts, nil
}

// RuntimeFactsParser parses runtime and execution related commands
type RuntimeFactsParser struct{}

// ParseRuntimeFacts parses runtime commands and returns structured facts
func (p *RuntimeFactsParser) ParseRuntimeFacts() map[string]*schemas.FactV1 {
	facts := make(map[string]*schemas.FactV1)

	// Parse uptime information
	if uptime, err := p.parseUptime(); err == nil {
		for k, v := range uptime {
			facts[k] = v
		}
	}

	// Parse date/time information
	if dateTime, err := p.parseDateTime(); err == nil {
		for k, v := range dateTime {
			facts[k] = v
		}
	}

	// Parse Python information
	if python, err := p.parsePythonInfo(); err == nil {
		for k, v := range python {
			facts[k] = v
		}
	}

	return facts
}

// parseUptime parses /proc/uptime
func (p *RuntimeFactsParser) parseUptime() (map[string]*schemas.FactV1, error) {
	facts := make(map[string]*schemas.FactV1)

	data, err := os.ReadFile("/proc/uptime")
	if err != nil {
		return facts, err
	}

	parts := strings.Fields(strings.TrimSpace(string(data)))
	if len(parts) >= 1 {
		if uptime, err := strconv.ParseFloat(parts[0], 64); err == nil {
			facts["uptime_seconds"] = &schemas.FactV1{
				Value:       uptime,
				Type:        "number",
				Description: "System uptime in seconds",
			}
		}
	}

	return facts, nil
}

// parseDateTime parses current date/time information
func (p *RuntimeFactsParser) parseDateTime() (map[string]*schemas.FactV1, error) {
	facts := make(map[string]*schemas.FactV1)

	// Get current timestamp
	timestamp := time.Now().Format(time.RFC3339)
	facts["date_time"] = &schemas.FactV1{
		Value: map[string]interface{}{
			"iso8601": timestamp,
			"year":    time.Now().Year(),
			"month":   int(time.Now().Month()),
			"weekday": time.Now().Weekday().String(),
			"day":     time.Now().Day(),
			"hour":    time.Now().Hour(),
			"minute":  time.Now().Minute(),
			"second":  time.Now().Second(),
		},
		Type:        "object",
		Description: "Current date and time information",
	}

	return facts, nil
}

// parsePythonInfo parses Python information
func (p *RuntimeFactsParser) parsePythonInfo() (map[string]*schemas.FactV1, error) {
	facts := make(map[string]*schemas.FactV1)

	// Check for Python executables
	pythonVersions := []string{"python3", "python", "python2"}
	var foundPython string
	var pythonPath string

	for _, version := range pythonVersions {
		if path, err := exec.LookPath(version); err == nil {
			foundPython = version
			pythonPath = path
			break
		}
	}

	if foundPython != "" {
		pythonInfo := map[string]interface{}{
			"executable": pythonPath,
			"type":       "CPython",
		}

		// Try to get Python version
		if cmd := exec.Command(foundPython, "--version"); cmd != nil {
			if output, err := cmd.Output(); err == nil {
				version := strings.TrimSpace(string(output))
				pythonInfo["version"] = version

				// Parse version into components
				if versionParts := p.parsePythonVersion(version); len(versionParts) > 0 {
					pythonInfo["version_info"] = versionParts

					// Add major, minor, micro components
					if len(versionParts) >= 3 {
						pythonInfo["major"] = versionParts[0]
						pythonInfo["minor"] = versionParts[1]
						pythonInfo["micro"] = versionParts[2]
					}
				}
			}
		}

		// Check for SSL context support
		if cmd := exec.Command(foundPython, "-c", "import ssl; print('ssl' in dir(ssl))"); cmd != nil {
			if output, err := cmd.Output(); err == nil {
				hasSSL := strings.TrimSpace(string(output)) == "True"
				pythonInfo["has_sslcontext"] = hasSSL
			}
		}

		facts["python"] = &schemas.FactV1{
			Value:       pythonInfo,
			Type:        "object",
			Description: "Python executable information",
		}

		// Also add simple version string
		if version, exists := pythonInfo["version"]; exists {
			facts["python_version"] = &schemas.FactV1{
				Value:       version,
				Type:        "string",
				Description: "Python version string",
			}
		}
	}

	return facts, nil
}

// parsePythonVersion parses Python version string into components
func (p *RuntimeFactsParser) parsePythonVersion(version string) []interface{} {
	// Remove "Python " prefix if present
	version = strings.TrimPrefix(version, "Python ")

	// Split by dots and convert to numbers
	parts := strings.Split(version, ".")
	var result []interface{}

	for _, part := range parts {
		// Remove any non-numeric suffix
		part = strings.Fields(part)[0]
		if num, err := strconv.Atoi(part); err == nil {
			result = append(result, num)
		} else {
			result = append(result, part)
		}
	}

	return result
}
