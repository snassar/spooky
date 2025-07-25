package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)

// Machine represents a machine configuration
type Machine struct {
	ID       string
	Host     string
	Port     int
	User     string
	Password string
	Tags     map[string]string
}

// Action represents an action configuration
type Action struct {
	Name        string
	Description string
	Command     string
	Script      string
	Tags        []string
	Machines    []string
	Parallel    bool
	Timeout     int
}

// ScaleConfig defines the parameters for different scale configurations
type ScaleConfig struct {
	Name        string
	Hardware    int
	VMs         int
	Containers  int
	Description string
}

// generateID creates a timestamp-based ID with random suffix
func generateID() string {
	now := time.Now()
	timestamp := now.Format("20060102")

	// Generate 4 random alphanumeric characters
	randomBytes := make([]byte, 2)
	if _, err := rand.Read(randomBytes); err != nil {
		panic(fmt.Errorf("failed to generate random bytes: %w", err))
	}
	randomStr := hex.EncodeToString(randomBytes)[:4]

	return fmt.Sprintf("%s+%s", timestamp, randomStr)
}

// generateGitStyleID creates a Git-style ID from metadata
func generateGitStyleID(metadata string) string {
	hash := sha256.Sum256([]byte(metadata))
	return hex.EncodeToString(hash[:8]) // 16-character short ID
}

// generateMachines creates machines for a specific scale configuration
func generateMachines(scale ScaleConfig) []Machine {
	var machines []Machine
	machineCount := 0

	// Hardware machines
	hardwarePerDC := (scale.Hardware + 1) / 2 // Round up for first datacenter
	hardwareCount := 0
	// FRA00 Hardware machines
	for i := 1; i <= hardwarePerDC && hardwareCount < scale.Hardware; i++ {
		metadata := fmt.Sprintf("hardware-fra00-%d-admin-debian12-vm-host-high", i)
		id := generateGitStyleID(metadata)
		machine := Machine{
			ID:       fmt.Sprintf("machine-%s", id),
			Host:     fmt.Sprintf("10.1.1.%d", i),
			Port:     22,
			User:     "admin",
			Password: fmt.Sprintf("hardware-secure-pass-%03d", i),
			Tags: map[string]string{
				"datacenter": "FRA00",
				"type":       "hardware",
				"role":       "vm-host",
				"os":         "debian12",
				"capacity":   "high",
			},
		}
		machines = append(machines, machine)
		machineCount++
		hardwareCount++
	}

	// BER0 Hardware machines
	ber0Hardware := scale.Hardware - hardwareCount // Remaining hardware for BER0
	for i := 1; i <= ber0Hardware && hardwareCount < scale.Hardware; i++ {
		metadata := fmt.Sprintf("hardware-ber0-%d-admin-debian12-vm-host-high", i)
		id := generateGitStyleID(metadata)
		machine := Machine{
			ID:       fmt.Sprintf("machine-%s", id),
			Host:     fmt.Sprintf("10.2.1.%d", i),
			Port:     22,
			User:     "admin",
			Password: fmt.Sprintf("hardware-secure-pass-%03d", i+hardwarePerDC),
			Tags: map[string]string{
				"datacenter": "BER0",
				"type":       "hardware",
				"role":       "vm-host",
				"os":         "debian12",
				"capacity":   "high",
			},
		}
		machines = append(machines, machine)
		machineCount++
		hardwareCount++
	}

	// VM machines - FRA00
	vmTypes := []struct {
		role     string
		os       string
		capacity string
		subnet   int
	}{
		{"web-server", "ubuntu22", "medium", 10},
		{"database", "ubuntu22", "high", 20},
		{"cache", "ubuntu22", "medium", 30},
		{"monitoring", "ubuntu22", "low", 40},
	}

	vmPerType := (scale.VMs + 7) / 8 // 4 types per datacenter, round up to ensure we get enough machines
	vmCount := 0
	for _, vmType := range vmTypes {
		for i := 1; i <= vmPerType && vmCount < scale.VMs; i++ {
			metadata := fmt.Sprintf("vm-fra00-%s-%d-admin-%s-%s", vmType.role, i, vmType.os, vmType.capacity)
			id := generateGitStyleID(metadata)
			machine := Machine{
				ID:       fmt.Sprintf("machine-%s", id),
				Host:     fmt.Sprintf("10.1.%d.%d", vmType.subnet, i),
				Port:     22,
				User:     "admin",
				Password: fmt.Sprintf("vm-secure-pass-%03d", machineCount),
				Tags: map[string]string{
					"datacenter": "FRA00",
					"type":       "vm",
					"role":       vmType.role,
					"os":         vmType.os,
					"capacity":   vmType.capacity,
				},
			}
			machines = append(machines, machine)
			machineCount++
			vmCount++
		}
	}

	// VM machines - BER0
	for _, vmType := range vmTypes {
		for i := 1; i <= vmPerType && vmCount < scale.VMs; i++ {
			metadata := fmt.Sprintf("vm-ber0-%s-%d-admin-%s-%s", vmType.role, i, vmType.os, vmType.capacity)
			id := generateGitStyleID(metadata)
			machine := Machine{
				ID:       fmt.Sprintf("machine-%s", id),
				Host:     fmt.Sprintf("10.2.%d.%d", vmType.subnet, i),
				Port:     22,
				User:     "admin",
				Password: fmt.Sprintf("vm-secure-pass-%03d", machineCount),
				Tags: map[string]string{
					"datacenter": "BER0",
					"type":       "vm",
					"role":       vmType.role,
					"os":         vmType.os,
					"capacity":   vmType.capacity,
				},
			}
			machines = append(machines, machine)
			machineCount++
			vmCount++
		}
	}

	// Container machines - FRA00
	containerTypes := []struct {
		role     string
		os       string
		capacity string
		subnet   int
	}{
		{"app-server", "alpine", "medium", 50},
		{"api-gateway", "alpine", "high", 60},
		{"worker", "alpine", "low", 70},
		{"redis", "alpine", "medium", 80},
	}

	containerPerType := (scale.Containers + 7) / 8 // 4 types per datacenter, round up to ensure we get enough machines
	containerCount := 0
	for _, containerType := range containerTypes {
		for i := 1; i <= containerPerType && containerCount < scale.Containers; i++ {
			metadata := fmt.Sprintf("container-fra00-%s-%d-admin-%s-%s", containerType.role, i, containerType.os, containerType.capacity)
			id := generateGitStyleID(metadata)
			machine := Machine{
				ID:       fmt.Sprintf("machine-%s", id),
				Host:     fmt.Sprintf("10.1.%d.%d", containerType.subnet, i),
				Port:     22,
				User:     "admin",
				Password: fmt.Sprintf("container-secure-pass-%03d", machineCount),
				Tags: map[string]string{
					"datacenter": "FRA00",
					"type":       "container",
					"role":       containerType.role,
					"os":         containerType.os,
					"capacity":   containerType.capacity,
				},
			}
			machines = append(machines, machine)
			machineCount++
			containerCount++
		}
	}

	// Container machines - BER0
	for _, containerType := range containerTypes {
		for i := 1; i <= containerPerType && containerCount < scale.Containers; i++ {
			metadata := fmt.Sprintf("container-ber0-%s-%d-admin-%s-%s", containerType.role, i, containerType.os, containerType.capacity)
			id := generateGitStyleID(metadata)
			machine := Machine{
				ID:       fmt.Sprintf("machine-%s", id),
				Host:     fmt.Sprintf("10.2.%d.%d", containerType.subnet, i),
				Port:     22,
				User:     "admin",
				Password: fmt.Sprintf("container-secure-pass-%03d", machineCount),
				Tags: map[string]string{
					"datacenter": "BER0",
					"type":       "container",
					"role":       containerType.role,
					"os":         containerType.os,
					"capacity":   containerType.capacity,
				},
			}
			machines = append(machines, machine)
			machineCount++
			containerCount++
		}
	}

	return machines
}

// generateActions creates a set of test actions
func generateActions() []Action {
	return []Action{
		{
			Name:        "system-update",
			Description: "Update system packages on all machines",
			Command:     "apt update && apt upgrade -y",
			Tags:        []string{"maintenance", "system"},
			Machines:    []string{"all"},
			Parallel:    true,
			Timeout:     1800,
		},
		{
			Name:        "security-scan",
			Description: "Run security vulnerability scan",
			Command:     "nmap -sV --script vuln localhost",
			Tags:        []string{"security", "monitoring"},
			Machines:    []string{"tag:type=hardware"},
			Parallel:    false,
			Timeout:     3600,
		},
		{
			Name:        "backup-databases",
			Description: "Create database backups",
			Command:     "pg_dump -Fc /var/lib/postgresql/backup.sql",
			Tags:        []string{"backup", "database"},
			Machines:    []string{"tag:role=database"},
			Parallel:    true,
			Timeout:     7200,
		},
		{
			Name:        "restart-services",
			Description: "Restart critical services",
			Command:     "systemctl restart nginx postgresql redis",
			Tags:        []string{"maintenance", "services"},
			Machines:    []string{"tag:role=web-server", "tag:role=database", "tag:role=cache"},
			Parallel:    false,
			Timeout:     300,
		},
		{
			Name:        "check-disk-space",
			Description: "Check available disk space",
			Command:     "df -h",
			Tags:        []string{"monitoring", "disk"},
			Machines:    []string{"all"},
			Parallel:    true,
			Timeout:     60,
		},
		{
			Name:        "update-firewall",
			Description: "Update firewall rules",
			Script:      "scripts/update-firewall.sh",
			Tags:        []string{"security", "network"},
			Machines:    []string{"tag:type=hardware"},
			Parallel:    false,
			Timeout:     600,
		},
		{
			Name:        "monitor-logs",
			Description: "Check system logs for errors",
			Command:     "journalctl --since '1 hour ago' --level=error",
			Tags:        []string{"monitoring", "logs"},
			Machines:    []string{"all"},
			Parallel:    true,
			Timeout:     120,
		},
		{
			Name:        "cleanup-temp",
			Description: "Clean up temporary files",
			Command:     "find /tmp -type f -mtime +7 -delete",
			Tags:        []string{"maintenance", "cleanup"},
			Machines:    []string{"all"},
			Parallel:    true,
			Timeout:     300,
		},
		{
			Name:        "check-memory",
			Description: "Check memory usage",
			Command:     "free -h",
			Tags:        []string{"monitoring", "memory"},
			Machines:    []string{"all"},
			Parallel:    true,
			Timeout:     30,
		},
		{
			Name:        "update-ssl-certs",
			Description: "Update SSL certificates",
			Script:      "scripts/update-ssl.sh",
			Tags:        []string{"security", "ssl"},
			Machines:    []string{"tag:role=web-server"},
			Parallel:    true,
			Timeout:     900,
		},
		{
			Name:        "database-maintenance",
			Description: "Perform database maintenance tasks",
			Script:      "scripts/db-maintenance.sh",
			Tags:        []string{"maintenance", "database"},
			Machines:    []string{"tag:role=database"},
			Parallel:    false,
			Timeout:     3600,
		},
		{
			Name:        "network-test",
			Description: "Test network connectivity",
			Command:     "ping -c 3 8.8.8.8",
			Tags:        []string{"monitoring", "network"},
			Machines:    []string{"all"},
			Parallel:    true,
			Timeout:     60,
		},
		{
			Name:        "update-monitoring",
			Description: "Update monitoring agents",
			Script:      "scripts/update-monitoring.sh",
			Tags:        []string{"monitoring", "update"},
			Machines:    []string{"tag:role=monitoring"},
			Parallel:    true,
			Timeout:     600,
		},
		{
			Name:        "backup-configs",
			Description: "Backup configuration files",
			Command:     "tar -czf /backup/configs-$(date +%Y%m%d).tar.gz /etc",
			Tags:        []string{"backup", "config"},
			Machines:    []string{"all"},
			Parallel:    true,
			Timeout:     1800,
		},
		{
			Name:        "health-check",
			Description: "Perform comprehensive health check",
			Script:      "scripts/health-check.sh",
			Tags:        []string{"monitoring", "health"},
			Machines:    []string{"all"},
			Parallel:    true,
			Timeout:     300,
		},
	}
}

// writeMachine writes a machine configuration to the file
func writeMachine(f *os.File, machine *Machine) {
	fmt.Fprintf(f, "machine \"%s\" {\n", machine.ID)
	fmt.Fprintf(f, "  host     = \"%s\"\n", machine.Host)
	fmt.Fprintf(f, "  port     = %d\n", machine.Port)
	fmt.Fprintf(f, "  user     = \"%s\"\n", machine.User)
	fmt.Fprintf(f, "  password = \"%s\"\n", machine.Password)

	// Write tags
	if len(machine.Tags) > 0 {
		fmt.Fprintf(f, "  tags = {\n")
		for key, value := range machine.Tags {
			fmt.Fprintf(f, "    %s = \"%s\"\n", key, value)
		}
		fmt.Fprintf(f, "  }\n")
	}
	fmt.Fprintf(f, "}\n\n")
}

// writeAction writes an action configuration to the file
func writeAction(f *os.File, action *Action) {
	fmt.Fprintf(f, "action \"%s\" {\n", action.Name)
	if action.Description != "" {
		fmt.Fprintf(f, "  description = \"%s\"\n", action.Description)
	}
	if action.Command != "" {
		fmt.Fprintf(f, "  command = \"%s\"\n", action.Command)
	}
	if action.Script != "" {
		fmt.Fprintf(f, "  script = \"%s\"\n", action.Script)
	}

	// Write tags
	if len(action.Tags) > 0 {
		fmt.Fprintf(f, "  tags = [\n")
		for _, tag := range action.Tags {
			fmt.Fprintf(f, "    \"%s\",\n", tag)
		}
		fmt.Fprintf(f, "  ]\n")
	}

	// Write machines
	if len(action.Machines) > 0 {
		fmt.Fprintf(f, "  machines = [\n")
		for _, machine := range action.Machines {
			fmt.Fprintf(f, "    \"%s\",\n", machine)
		}
		fmt.Fprintf(f, "  ]\n")
	}
	fmt.Fprintf(f, "  parallel    = %t\n", action.Parallel)
	fmt.Fprintf(f, "  timeout     = %d\n", action.Timeout)
	fmt.Fprintf(f, "}\n\n")
}

// runSpookyProjectInit runs the spooky project init command
func runSpookyProjectInit(projectName, projectPath string) error {
	cmd := exec.Command("./build/spooky", "project", "init", projectName, projectPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// writeInventoryFile writes the inventory.hcl file
func writeInventoryFile(projectPath string, machines []Machine) error {
	filename := filepath.Join(projectPath, "inventory.hcl")

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating inventory file: %v", err)
	}
	defer file.Close()

	// Write header
	fmt.Fprintf(file, "# Inventory file for test project\n")
	fmt.Fprintf(file, "# Generated machines for spooky SSH automation tool\n\n")

	// Write inventory wrapper block
	fmt.Fprintf(file, "inventory {\n")

	// Write all machines
	for i := range machines {
		writeMachine(file, &machines[i])
	}

	// Close inventory wrapper block
	fmt.Fprintf(file, "}\n")

	fmt.Printf("Generated inventory file: %s\n", filename)
	fmt.Printf("Total machines: %d\n", len(machines))
	return nil
}

// writeActionsFile writes the actions.hcl file
func writeActionsFile(projectPath string, actions []Action) error {
	filename := filepath.Join(projectPath, "actions.hcl")

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating actions file: %v", err)
	}
	defer file.Close()

	// Write header
	fmt.Fprintf(file, "# Actions file for test project\n")
	fmt.Fprintf(file, "# Test actions for spooky SSH automation tool\n\n")

	// Write actions wrapper block
	fmt.Fprintf(file, "actions {\n")

	// Write all actions
	for i := range actions {
		writeAction(file, &actions[i])
	}

	// Close actions wrapper block
	fmt.Fprintf(file, "}\n")

	fmt.Printf("Generated actions file: %s\n", filename)
	fmt.Printf("Total actions: %d\n", len(actions))
	return nil
}

// generateTestProject generates a test project with the specified scale
func generateTestProject(scale ScaleConfig, outputDir string) error {
	configID := generateID()
	projectName := fmt.Sprintf("%s-scale-test-%s", scale.Name, configID)
	projectPath := filepath.Join(outputDir, projectName)

	fmt.Printf("Generating %s scale test project: %s\n", scale.Name, projectPath)

	// Run spooky project init - pass the project name and the parent directory
	// This will create: outputDir/projectName/projectName (the actual project directory)
	fmt.Printf("Running spooky project init...\n")
	if err := runSpookyProjectInit(projectName, outputDir); err != nil {
		return fmt.Errorf("failed to run spooky project init: %w", err)
	}

	// The actual project directory is created by spooky project init
	actualProjectPath := projectPath

	// Generate machines and actions
	machines := generateMachines(scale)
	actions := generateActions()

	// Write inventory file to the actual project directory
	if err := writeInventoryFile(actualProjectPath, machines); err != nil {
		return fmt.Errorf("failed to write inventory file: %w", err)
	}

	// Write actions file to the actual project directory
	if err := writeActionsFile(actualProjectPath, actions); err != nil {
		return fmt.Errorf("failed to write actions file: %w", err)
	}

	fmt.Printf("Successfully generated test project: %s\n", actualProjectPath)
	fmt.Printf("Project contains %d machines and %d actions\n", len(machines), len(actions))
	return nil
}

func main() {
	var (
		scale      string
		hardware   int
		vms        int
		containers int
		outputDir  string
	)

	rootCmd := &cobra.Command{
		Use:   "generate-test-project",
		Short: "Generate test projects for spooky with different scale configurations",
		Long: `Generate test projects for spooky with different scale configurations.

This tool creates test projects with varying numbers of machines and actions to test
the spooky SSH automation tool at different scales.

Available predefined scales:
- small: 25 machines (2-5 hardware, rest random mix of VMs and containers)
- medium: 250 machines (20-40 hardware, rest random mix of VMs and containers)
- large: 1500 machines (150-250 hardware, rest random mix of VMs and containers)
- testing: 10000 machines (completely random distribution of hardware, VMs, and containers)

For custom scales, use the --hardware, --vms, and --containers flags.`,
		RunE: func(_ *cobra.Command, _ []string) error {
			// Create output directory if it doesn't exist
			if err := os.MkdirAll(outputDir, 0o755); err != nil {
				return fmt.Errorf("error creating output directory: %w", err)
			}

			var scaleConfig ScaleConfig

			// Parse scale configuration
			if scale == "custom" {
				if hardware <= 0 || vms <= 0 || containers <= 0 {
					return fmt.Errorf("custom scale requires --hardware, --vms, and --containers to be greater than 0")
				}
				scaleConfig = ScaleConfig{
					Name:        "custom",
					Hardware:    hardware,
					VMs:         vms,
					Containers:  containers,
					Description: fmt.Sprintf("Custom scale with %d hardware + %d VMs + %d containers", hardware, vms, containers),
				}
			} else {
				// Predefined scales
				switch scale {
				case "small":
					// 25 total: 2-5 hardware, rest random mix of VMs and containers
					hardwareCount := 3 // Fixed for now, could be randomized
					remaining := 25 - hardwareCount
					vmCount := remaining / 2
					containerCount := remaining - vmCount
					scaleConfig = ScaleConfig{
						Name:        "small",
						Hardware:    hardwareCount,
						VMs:         vmCount,
						Containers:  containerCount,
						Description: fmt.Sprintf("Small hosting provider with 25 machines (%d hardware + %d VMs + %d containers)", hardwareCount, vmCount, containerCount),
					}
				case "medium":
					// 250 total: 20-40 hardware, rest random mix of VMs and containers
					hardwareCount := 30 // Fixed for now, could be randomized
					remaining := 250 - hardwareCount
					vmCount := remaining / 2
					containerCount := remaining - vmCount
					scaleConfig = ScaleConfig{
						Name:        "medium",
						Hardware:    hardwareCount,
						VMs:         vmCount,
						Containers:  containerCount,
						Description: fmt.Sprintf("Medium hosting provider with 250 machines (%d hardware + %d VMs + %d containers)", hardwareCount, vmCount, containerCount),
					}
				case "large":
					// 1500 total: 150-250 hardware, rest random mix of VMs and containers
					hardwareCount := 200 // Fixed for now, could be randomized
					remaining := 1500 - hardwareCount
					vmCount := remaining / 2
					containerCount := remaining - vmCount
					scaleConfig = ScaleConfig{
						Name:        "large",
						Hardware:    hardwareCount,
						VMs:         vmCount,
						Containers:  containerCount,
						Description: fmt.Sprintf("Large hosting provider with 1500 machines (%d hardware + %d VMs + %d containers)", hardwareCount, vmCount, containerCount),
					}
				case "testing":
					// 10000 total: completely random distribution
					hardwareCount := 2000 // Random distribution
					vmCount := 4000
					containerCount := 4000
					scaleConfig = ScaleConfig{
						Name:        "testing",
						Hardware:    hardwareCount,
						VMs:         vmCount,
						Containers:  containerCount,
						Description: fmt.Sprintf("Testing scale with 10000 machines (%d hardware + %d VMs + %d containers)", hardwareCount, vmCount, containerCount),
					}
				default:
					return fmt.Errorf("unknown scale: %s. Available scales: small, medium, large, testing, or custom", scale)
				}
			}

			// Generate the test project
			if err := generateTestProject(scaleConfig, outputDir); err != nil {
				return fmt.Errorf("error generating test project: %w", err)
			}

			fmt.Printf("\n=== Test project generated successfully ===\n")
			fmt.Printf("Scale: %s\n", scaleConfig.Name)
			fmt.Printf("Hardware machines: %d\n", scaleConfig.Hardware)
			fmt.Printf("VM machines: %d\n", scaleConfig.VMs)
			fmt.Printf("Container machines: %d\n", scaleConfig.Containers)
			fmt.Printf("Total machines: %d\n", scaleConfig.Hardware+scaleConfig.VMs+scaleConfig.Containers)
			fmt.Printf("Output directory: %s\n", outputDir)

			return nil
		},
	}

	// Add flags
	rootCmd.Flags().StringVarP(&scale, "scale", "s", "medium", "Scale configuration (small, medium, large, testing, or custom)")
	rootCmd.Flags().IntVarP(&hardware, "hardware", "w", 0, "Number of hardware machines (for custom scale)")
	rootCmd.Flags().IntVarP(&vms, "vms", "v", 0, "Number of VM machines (for custom scale)")
	rootCmd.Flags().IntVarP(&containers, "containers", "c", 0, "Number of container machines (for custom scale)")
	rootCmd.Flags().StringVarP(&outputDir, "output", "o", "./test-projects", "Output directory for generated projects")

	// Mark hardware, vms, and containers as required when scale is custom
	rootCmd.MarkFlagsRequiredTogether("hardware", "vms", "containers")

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
