package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
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
	hardwarePerDC := scale.Hardware / 2
	// FRA00 Hardware machines
	for i := 1; i <= hardwarePerDC; i++ {
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
	}

	// BER0 Hardware machines
	for i := 1; i <= hardwarePerDC; i++ {
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

	vmPerType := scale.VMs / 8 // 4 types per datacenter
	for _, vmType := range vmTypes {
		for i := 1; i <= vmPerType; i++ {
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
		}
	}

	// VM machines - BER0
	for _, vmType := range vmTypes {
		for i := 1; i <= vmPerType; i++ {
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
	fmt.Fprintf(file, "# Generated with Git-style IDs for deterministic identification\n")
	fmt.Fprintf(file, "# Total machines: %d\n\n", len(machines))

	// Write all machines
	for i, machine := range machines {
		writeMachine(file, &machine)
		if i%1000 == 0 && i > 0 {
			fmt.Printf("Written %d machines to inventory...\n", i)
		}
	}

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

	// Write all actions
	for i := range actions {
		writeAction(file, &actions[i])
	}

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
	if len(os.Args) < 2 {
		fmt.Println("Usage: generate-test-project <scale> [output-dir]")
		fmt.Println("Scales: small, medium, large, or custom")
		fmt.Println("Example: generate-test-project medium ./test-projects")
		fmt.Println("Example: generate-test-project custom:100:300 ./test-projects")
		os.Exit(1)
	}

	scaleArg := os.Args[1]
	outputDir := "./test-projects"
	if len(os.Args) > 2 {
		outputDir = os.Args[2]
	}

	// Create output directory if it doesn't exist
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		fmt.Printf("Error creating output directory: %v\n", err)
		os.Exit(1)
	}

	var scale ScaleConfig

	// Parse scale argument
	if strings.HasPrefix(scaleArg, "custom:") {
		// Custom scale format: custom:hardware:vms
		parts := strings.Split(scaleArg, ":")
		if len(parts) != 3 {
			fmt.Println("Custom scale format: custom:hardware:vms")
			fmt.Println("Example: custom:50:150")
			os.Exit(1)
		}

		hardware, err := strconv.Atoi(parts[1])
		if err != nil {
			fmt.Printf("Invalid hardware count: %s\n", parts[1])
			os.Exit(1)
		}

		vms, err := strconv.Atoi(parts[2])
		if err != nil {
			fmt.Printf("Invalid VM count: %s\n", parts[2])
			os.Exit(1)
		}

		scale = ScaleConfig{
			Name:        "custom",
			Hardware:    hardware,
			VMs:         vms,
			Description: fmt.Sprintf("Custom scale with %d hardware + %d VMs", hardware, vms),
		}
	} else {
		// Predefined scales
		switch scaleArg {
		case "small":
			scale = ScaleConfig{
				Name:        "small",
				Hardware:    10,
				VMs:         30,
				Description: "Small hosting provider with 40 servers (10 hardware + 30 VMs)",
			}
		case "medium":
			scale = ScaleConfig{
				Name:        "medium",
				Hardware:    100,
				VMs:         300,
				Description: "Medium hosting provider with 400 servers (100 hardware + 300 VMs)",
			}
		case "large":
			scale = ScaleConfig{
				Name:        "large",
				Hardware:    2500,
				VMs:         7500,
				Description: "Large hosting provider with 10,000 servers (2,500 hardware + 7,500 VMs)",
			}
		default:
			fmt.Printf("Unknown scale: %s\n", scaleArg)
			fmt.Println("Available scales: small, medium, large, or custom:hardware:vms")
			os.Exit(1)
		}
	}

	// Generate the test project
	if err := generateTestProject(scale, outputDir); err != nil {
		fmt.Printf("Error generating test project: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\n=== Test project generated successfully ===\n")
	fmt.Printf("Scale: %s\n", scale.Name)
	fmt.Printf("Hardware machines: %d\n", scale.Hardware)
	fmt.Printf("VM machines: %d\n", scale.VMs)
	fmt.Printf("Total machines: %d\n", scale.Hardware+scale.VMs)
	fmt.Printf("Output directory: %s\n", outputDir)
}
