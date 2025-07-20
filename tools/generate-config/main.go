package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
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
			Password: fmt.Sprintf("hardware-secure-pass-%d", hardwarePerDC+i),
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

	// VMs
	vmCount := 0
	vmsPerDC := scale.VMs / 2
	vmsPerType := vmsPerDC / 4 // 4 types: database, web, workload, storage

	// FRA00 VMs
	// Database VMs
	for i := 1; i <= vmsPerType; i++ {
		metadata := fmt.Sprintf("vm-fra00-db-%d-debian-debian12-database-production", i)
		id := generateGitStyleID(metadata)
		var dbType string
		switch i % 3 {
		case 1:
			dbType = "postgresql"
		case 2:
			dbType = "mysql"
		default:
			dbType = "mongodb"
		}
		tier := "production"
		if i > vmsPerType*4/5 {
			tier = "staging"
		}
		machine := Machine{
			ID:       fmt.Sprintf("vm-%s", id),
			Host:     fmt.Sprintf("10.1.10.%d", i),
			Port:     22,
			User:     "debian",
			Password: fmt.Sprintf("vm-secure-pass-%04d", vmCount+1),
			Tags: map[string]string{
				"datacenter": "FRA00",
				"type":       "vm",
				"role":       "database",
				"os":         "debian12",
				"tier":       tier,
				"db_type":    dbType,
			},
		}
		machines = append(machines, machine)
		vmCount++
		machineCount++
	}

	// Web VMs
	for i := 1; i <= vmsPerType; i++ {
		metadata := fmt.Sprintf("vm-fra00-web-%d-debian-debian12-web-production", i)
		id := generateGitStyleID(metadata)
		webType := "nginx"
		if i%2 == 0 {
			webType = "apache"
		}
		tier := "production"
		if i > vmsPerType*4/5 {
			tier = "staging"
		}
		machine := Machine{
			ID:       fmt.Sprintf("vm-%s", id),
			Host:     fmt.Sprintf("10.1.20.%d", i),
			Port:     22,
			User:     "debian",
			Password: fmt.Sprintf("vm-secure-pass-%04d", vmCount+1),
			Tags: map[string]string{
				"datacenter": "FRA00",
				"type":       "vm",
				"role":       "web",
				"os":         "debian12",
				"tier":       tier,
				"web_type":   webType,
			},
		}
		machines = append(machines, machine)
		vmCount++
		machineCount++
	}

	// Workload VMs
	for i := 1; i <= vmsPerType; i++ {
		metadata := fmt.Sprintf("vm-fra00-workload-%d-debian-debian12-workload-production", i)
		id := generateGitStyleID(metadata)
		workloadType := "compute"
		if i%2 == 0 {
			workloadType = "batch"
		}
		tier := "production"
		if i > vmsPerType*4/5 {
			tier = "staging"
		}
		machine := Machine{
			ID:       fmt.Sprintf("vm-%s", id),
			Host:     fmt.Sprintf("10.1.30.%d", i),
			Port:     22,
			User:     "debian",
			Password: fmt.Sprintf("vm-secure-pass-%04d", vmCount+1),
			Tags: map[string]string{
				"datacenter":    "FRA00",
				"type":          "vm",
				"role":          "workload",
				"os":            "debian12",
				"tier":          tier,
				"workload_type": workloadType,
			},
		}
		machines = append(machines, machine)
		vmCount++
		machineCount++
	}

	// Storage VMs
	for i := 1; i <= vmsPerType; i++ {
		metadata := fmt.Sprintf("vm-fra00-storage-%d-debian-debian12-storage-production", i)
		id := generateGitStyleID(metadata)
		storageType := "block"
		if i%2 == 0 {
			storageType = "object"
		}
		tier := "production"
		if i > vmsPerType*4/5 {
			tier = "staging"
		}
		machine := Machine{
			ID:       fmt.Sprintf("vm-%s", id),
			Host:     fmt.Sprintf("10.1.40.%d", i),
			Port:     22,
			User:     "debian",
			Password: fmt.Sprintf("vm-secure-pass-%04d", vmCount+1),
			Tags: map[string]string{
				"datacenter":   "FRA00",
				"type":         "vm",
				"role":         "storage",
				"os":           "debian12",
				"tier":         tier,
				"storage_type": storageType,
			},
		}
		machines = append(machines, machine)
		vmCount++
		machineCount++
	}

	// BER0 VMs (same pattern as FRA00)
	// Database VMs
	for i := 1; i <= vmsPerType; i++ {
		metadata := fmt.Sprintf("vm-ber0-db-%d-debian-debian12-database-production", i)
		id := generateGitStyleID(metadata)
		var dbType string
		switch i % 3 {
		case 1:
			dbType = "postgresql"
		case 2:
			dbType = "mysql"
		default:
			dbType = "mongodb"
		}
		tier := "production"
		if i > vmsPerType*4/5 {
			tier = "staging"
		}
		machine := Machine{
			ID:       fmt.Sprintf("vm-%s", id),
			Host:     fmt.Sprintf("10.2.10.%d", i),
			Port:     22,
			User:     "debian",
			Password: fmt.Sprintf("vm-secure-pass-%04d", vmCount+1),
			Tags: map[string]string{
				"datacenter": "BER0",
				"type":       "vm",
				"role":       "database",
				"os":         "debian12",
				"tier":       tier,
				"db_type":    dbType,
			},
		}
		machines = append(machines, machine)
		vmCount++
		machineCount++
	}

	// Web VMs
	for i := 1; i <= vmsPerType; i++ {
		metadata := fmt.Sprintf("vm-ber0-web-%d-debian-debian12-web-production", i)
		id := generateGitStyleID(metadata)
		webType := "nginx"
		if i%2 == 0 {
			webType = "apache"
		}
		tier := "production"
		if i > vmsPerType*4/5 {
			tier = "staging"
		}
		machine := Machine{
			ID:       fmt.Sprintf("vm-%s", id),
			Host:     fmt.Sprintf("10.2.20.%d", i),
			Port:     22,
			User:     "debian",
			Password: fmt.Sprintf("vm-secure-pass-%04d", vmCount+1),
			Tags: map[string]string{
				"datacenter": "BER0",
				"type":       "vm",
				"role":       "web",
				"os":         "debian12",
				"tier":       tier,
				"web_type":   webType,
			},
		}
		machines = append(machines, machine)
		vmCount++
		machineCount++
	}

	// Workload VMs
	for i := 1; i <= vmsPerType; i++ {
		metadata := fmt.Sprintf("vm-ber0-workload-%d-debian-debian12-workload-production", i)
		id := generateGitStyleID(metadata)
		workloadType := "compute"
		if i%2 == 0 {
			workloadType = "batch"
		}
		tier := "production"
		if i > vmsPerType*4/5 {
			tier = "staging"
		}
		machine := Machine{
			ID:       fmt.Sprintf("vm-%s", id),
			Host:     fmt.Sprintf("10.2.30.%d", i),
			Port:     22,
			User:     "debian",
			Password: fmt.Sprintf("vm-secure-pass-%04d", vmCount+1),
			Tags: map[string]string{
				"datacenter":    "BER0",
				"type":          "vm",
				"role":          "workload",
				"os":            "debian12",
				"tier":          tier,
				"workload_type": workloadType,
			},
		}
		machines = append(machines, machine)
		vmCount++
		machineCount++
	}

	// Storage VMs
	for i := 1; i <= vmsPerType; i++ {
		metadata := fmt.Sprintf("vm-ber0-storage-%d-debian-debian12-storage-production", i)
		id := generateGitStyleID(metadata)
		storageType := "block"
		if i%2 == 0 {
			storageType = "object"
		}
		tier := "production"
		if i > vmsPerType*4/5 {
			tier = "staging"
		}
		machine := Machine{
			ID:       fmt.Sprintf("vm-%s", id),
			Host:     fmt.Sprintf("10.2.40.%d", i),
			Port:     22,
			User:     "debian",
			Password: fmt.Sprintf("vm-secure-pass-%04d", vmCount+1),
			Tags: map[string]string{
				"datacenter":   "BER0",
				"type":         "vm",
				"role":         "storage",
				"os":           "debian12",
				"tier":         tier,
				"storage_type": storageType,
			},
		}
		machines = append(machines, machine)
		vmCount++
		machineCount++
	}

	fmt.Printf("Generated %d machines (%d hardware + %d VMs) for %s scale\n", machineCount, scale.Hardware, scale.VMs, scale.Name)
	return machines
}

// generateActions creates test actions
func generateActions() []Action {
	return []Action{
		{
			Name:        "check-production-status",
			Description: "Check status of all production servers",
			Command:     "uptime && df -h && systemctl status --no-pager",
			Tags:        []string{"tier=production"},
			Parallel:    true,
			Timeout:     300,
		},
		{
			Name:        "update-databases",
			Description: "Update all database servers",
			Command:     "apt update && apt upgrade -y",
			Tags:        []string{"role=database"},
			Parallel:    true,
			Timeout:     600,
		},
		{
			Name:        "check-fra00-web",
			Description: "Check FRA00 web servers specifically",
			Command:     "systemctl status nginx apache2 --no-pager",
			Tags:        []string{"datacenter=FRA00", "role=web"},
			Parallel:    true,
			Timeout:     120,
		},
		{
			Name:        "backup-storage",
			Description: "Create backups on all storage servers",
			Script:      "/usr/local/bin/backup-storage.sh",
			Tags:        []string{"role=storage"},
			Parallel:    false,
			Timeout:     1800,
		},
		{
			Name:        "check-hardware",
			Description: "Check hardware server status",
			Command:     "lscpu && free -h && df -h",
			Tags:        []string{"type=hardware"},
			Parallel:    true,
			Timeout:     180,
		},
		{
			Name:        "update-staging",
			Description: "Update all staging servers",
			Command:     "apt update && apt upgrade -y",
			Tags:        []string{"tier=staging"},
			Parallel:    true,
			Timeout:     600,
		},
		{
			Name:        "check-ber0-db",
			Description: "Check BER0 database servers",
			Command:     "systemctl status postgresql mysql mongod --no-pager",
			Tags:        []string{"datacenter=BER0", "role=database"},
			Parallel:    true,
			Timeout:     120,
		},
		{
			Name:        "monitor-compute",
			Description: "Monitor compute workload servers",
			Command:     "htop --batch --iterations=1 && nvidia-smi",
			Tags:        []string{"workload_type=compute"},
			Parallel:    true,
			Timeout:     60,
		},
		{
			Name:        "check-nginx",
			Description: "Check all nginx web servers",
			Command:     "nginx -t && systemctl status nginx --no-pager",
			Tags:        []string{"web_type=nginx"},
			Parallel:    true,
			Timeout:     90,
		},
		{
			Name:        "full-system-check",
			Description: "Comprehensive system check",
			Machines:    []string{"machine-550e8400e29b41d4", "vm-550e8400e29b41da", "vm-550e8400e29b41e6"},
			Command:     "uptime && df -h && free -h && systemctl --failed --no-pager",
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
	fmt.Fprintf(f, "  tags = {\n")
	for k, v := range machine.Tags {
		fmt.Fprintf(f, "    %s = \"%s\"\n", k, v)
	}
	fmt.Fprintf(f, "  }\n")
	fmt.Fprintf(f, "}\n\n")
}

// writeAction writes an action configuration to the file
func writeAction(f *os.File, action *Action) {
	fmt.Fprintf(f, "action \"%s\" {\n", action.Name)
	fmt.Fprintf(f, "  description = \"%s\"\n", action.Description)
	if action.Command != "" {
		fmt.Fprintf(f, "  command     = \"%s\"\n", action.Command)
	}
	if action.Script != "" {
		fmt.Fprintf(f, "  script      = \"%s\"\n", action.Script)
	}
	if len(action.Tags) > 0 {
		fmt.Fprintf(f, "  tags        = [\"%s\"]\n", strings.Join(action.Tags, "\", \""))
	}
	if len(action.Machines) > 0 {
		fmt.Fprintf(f, "  machines     = [\"%s\"]\n", strings.Join(action.Machines, "\", \""))
	}
	fmt.Fprintf(f, "  parallel    = %t\n", action.Parallel)
	fmt.Fprintf(f, "  timeout     = %d\n", action.Timeout)
	fmt.Fprintf(f, "}\n\n")
}

// generateConfigFile generates a configuration file for a specific scale
func generateConfigFile(scale ScaleConfig) error {
	configID := generateID()
	configDir := "../examples/actions"
	filename := filepath.Join(configDir, fmt.Sprintf("%s-scale-example-%s.hcl", scale.Name, configID))

	// Create output directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(filename), 0o755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()

	// Write header
	caser := cases.Title(language.English)
	fmt.Fprintf(file, "# %s configuration for spooky SSH automation tool\n", caser.String(scale.Name))
	fmt.Fprintf(file, "# %s\n", scale.Description)
	fmt.Fprintf(file, "# Data centers: FRA00 (Frankfurt) and BER0 (Berlin)\n")
	fmt.Fprintf(file, "# IP range: 10.0.0.0/8\n")
	fmt.Fprintf(file, "# Generated with Git-style IDs for deterministic identification\n")
	fmt.Fprintf(file, "# Config ID: %s\n\n", configID)

	// Write machine section header
	fmt.Fprintf(file, "# =============================================================================\n")
	fmt.Fprintf(file, "# MACHINES (%d total)\n", scale.Hardware+scale.VMs)
	fmt.Fprintf(file, "# =============================================================================\n\n")

	// Generate and write all machines
	machines := generateMachines(scale)
	fmt.Printf("Writing %d machines to %s...\n", len(machines), filename)
	for i, machine := range machines {
		writeMachine(file, &machine)
		if i%1000 == 0 && i > 0 {
			fmt.Printf("Written %d machines...\n", i)
		}
	}

	// Write action section header
	fmt.Fprintf(file, "# =============================================================================\n")
	fmt.Fprintf(file, "# ACTIONS FOR TESTING\n")
	fmt.Fprintf(file, "# =============================================================================\n\n")

	// Generate and write all actions
	actions := generateActions()
	for i := range actions {
		writeAction(file, &actions[i])
	}

	fmt.Printf("Generated configuration file: %s\n", filename)
	fmt.Printf("Total machines: %d\n", len(machines))
	fmt.Printf("Total actions: %d\n", len(actions))

	return nil
}

func main() {
	// Define scale configurations
	scales := []ScaleConfig{
		{
			Name:        "small",
			Hardware:    10, // 5 per datacenter
			VMs:         30, // 15 per datacenter (4 types, ~4 each)
			Description: "Small hosting provider with 40 servers (10 hardware + 30 VMs)",
		},
		{
			Name:        "medium",
			Hardware:    100, // 50 per datacenter
			VMs:         300, // 150 per datacenter (4 types, ~38 each)
			Description: "Medium hosting provider with 400 servers (100 hardware + 300 VMs)",
		},
		{
			Name:        "large",
			Hardware:    2500, // 1250 per datacenter
			VMs:         7500, // 3750 per datacenter (4 types, ~938 each)
			Description: "Large hosting provider with 10,000 servers (2,500 hardware + 7,500 VMs)",
		},
	}

	// Generate all scale configurations
	for _, scale := range scales {
		fmt.Printf("\n=== Generating %s scale configuration ===\n", scale.Name)
		if err := generateConfigFile(scale); err != nil {
			fmt.Printf("Error generating %s scale config: %v\n", scale.Name, err)
			os.Exit(1)
		}
	}

	fmt.Printf("\n=== All configurations generated successfully ===\n")
}
