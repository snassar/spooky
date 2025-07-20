package cli

import (
	"fmt"
	"os"
	"time"

	"spooky/internal/config"
	"spooky/internal/logging"
	"spooky/internal/ssh"

	"github.com/spf13/cobra"
)

var MachinesCmd = &cobra.Command{
	Use:   "machines",
	Short: "Manage machine information and operations",
	Long: `Manage machine information and perform operations on machines.

This command group provides tools for listing machines, testing connectivity,
and checking machine status.

Examples:
  # List all machines from a configuration file
  spooky machines list config.hcl

  # Test connectivity to a specific machine
  spooky machines ping machine-name

  # Check status of a machine
  spooky machines status machine-name

  # Execute actions on machines
  spooky machines action action-name config.hcl`,
}

var machinesListCmd = &cobra.Command{
	Use:   "list [config-file]",
	Short: "List all machines",
	Long: `List all machines defined in configuration files.

Examples:
  spooky machines list                    # List machines from default config
  spooky machines list config.hcl         # List machines from specific config
  spooky machines list --format json      # Output in JSON format`,
	Args: cobra.MaximumNArgs(1),
	RunE: runMachinesList,
}

var machinesPingCmd = &cobra.Command{
	Use:   "ping <machine>",
	Short: "Test connectivity to a machine",
	Long: `Test SSH connectivity to a specific machine.

This command attempts to establish an SSH connection to the specified machine
and reports the connection status and response time.

Examples:
  spooky machines ping server1
  spooky machines ping server1 --timeout 10s
  spooky machines ping server1 --config config.hcl`,
	Args: cobra.ExactArgs(1),
	RunE: runMachinesPing,
}

var machinesStatusCmd = &cobra.Command{
	Use:   "status <machine>",
	Short: "Check machine status",
	Long: `Check the status of a specific machine.

This command performs various health checks on the machine including:
- SSH connectivity
- System uptime
- Basic system information
- Service status (if configured)

Examples:
  spooky machines status server1
  spooky machines status server1 --detailed
  spooky machines status server1 --config config.hcl`,
	Args: cobra.ExactArgs(1),
	RunE: runMachinesStatus,
}

var machinesActionCmd = &cobra.Command{
	Use:   "action <source>",
	Short: "Execute actions on machines",
	Long: `Execute actions defined in configuration files on machines.

This command runs actions from configuration files on target machines.
Actions can be executed sequentially or in parallel depending on configuration.

Examples:
  spooky machines action config.hcl
  spooky machines action config.hcl --parallel 5
  spooky machines action config.hcl --timeout 60`,
	Args: cobra.ExactArgs(1),
	RunE: runMachinesAction,
}

func init() {
	// Add flags for machines list command
	machinesListCmd.Flags().String("format", "table", "Output format: table, json, yaml")
	machinesListCmd.Flags().Bool("detailed", false, "Show detailed machine information")

	// Add flags for machines ping command
	machinesPingCmd.Flags().Duration("timeout", 30*time.Second, "Connection timeout")
	machinesPingCmd.Flags().Int("count", 1, "Number of ping attempts")

	// Add flags for machines status command
	machinesStatusCmd.Flags().Bool("detailed", false, "Show detailed status information")
	machinesStatusCmd.Flags().Bool("services", false, "Check service status")

	// Add flags for machines action command
	machinesActionCmd.Flags().String("hosts", "", "Comma-separated list of target hosts")
	machinesActionCmd.Flags().Int("parallel", 5, "Number of parallel executions")
	machinesActionCmd.Flags().Int("timeout", 30, "Execution timeout per host in seconds")
	machinesActionCmd.Flags().Int("retry", 3, "Number of retry attempts")
	machinesActionCmd.Flags().String("tags", "", "Comma-separated list of tags to execute")
	machinesActionCmd.Flags().String("skip-tags", "", "Comma-separated list of tags to skip")

	// Add subcommands to machines command
	MachinesCmd.AddCommand(machinesListCmd)
	MachinesCmd.AddCommand(machinesPingCmd)
	MachinesCmd.AddCommand(machinesStatusCmd)
	MachinesCmd.AddCommand(machinesActionCmd)

	// Note: machinesCmd will be added to root in main.go
}

func runMachinesList(cmd *cobra.Command, args []string) error {
	configFile := "."
	if len(args) > 0 {
		configFile = args[0]
	}

	format, _ := cmd.Flags().GetString("format")
	detailed, _ := cmd.Flags().GetBool("detailed")

	// Load configuration
	cfg, err := config.ParseConfig(configFile)
	if err != nil {
		return fmt.Errorf("error loading configuration: %w", err)
	}

	if len(cfg.Machines) == 0 {
		fmt.Println("No machines found in configuration")
		return nil
	}

	// Output based on format
	switch format {
	case "json":
		return outputMachinesJSON(cfg.Machines, detailed)
	case "yaml":
		return outputMachinesYAML(cfg.Machines, detailed)
	default:
		return outputMachinesTable(cfg.Machines, detailed)
	}
}

func runMachinesPing(cmd *cobra.Command, args []string) error {
	machineName := args[0]
	timeout, _ := cmd.Flags().GetDuration("timeout")
	count, _ := cmd.Flags().GetInt("count")

	// Load configuration to find machine
	cfg, err := config.ParseConfig(".")
	if err != nil {
		return fmt.Errorf("error loading configuration: %w", err)
	}

	// Find machine in configuration
	var targetMachine *config.Machine
	for _, machine := range cfg.Machines {
		if machine.Name == machineName {
			targetMachine = &machine
			break
		}
	}

	if targetMachine == nil {
		return fmt.Errorf("machine '%s' not found in configuration", machineName)
	}

	fmt.Printf("Pinging machine: %s (%s)\n", machineName, targetMachine.Host)

	// Perform ping attempts
	successCount := 0
	var totalTime time.Duration

	for i := 0; i < count; i++ {
		if count > 1 {
			fmt.Printf("Attempt %d/%d: ", i+1, count)
		}

		start := time.Now()

		// Create SSH client
		client, err := ssh.NewSSHClient(targetMachine, int(timeout.Seconds()))

		if err != nil {
			fmt.Printf("Failed: %v\n", err)
			continue
		}

		// Test connection with a simple command
		_, err = client.ExecuteCommand("echo 'ping test'")
		client.Close()

		duration := time.Since(start)

		if err != nil {
			fmt.Printf("Failed: %v (%.2fms)\n", err, float64(duration.Microseconds())/1000)
		} else {
			fmt.Printf("Success: %.2fms\n", float64(duration.Microseconds())/1000)
			successCount++
			totalTime += duration
		}

		if count > 1 && i < count-1 {
			time.Sleep(500 * time.Millisecond)
		}
	}

	// Summary
	fmt.Printf("\nPing Summary:\n")
	fmt.Printf("  Successful: %d/%d\n", successCount, count)
	if successCount > 0 {
		avgTime := totalTime / time.Duration(successCount)
		fmt.Printf("  Average time: %.2fms\n", float64(avgTime.Microseconds())/1000)
	}

	if successCount == 0 {
		return fmt.Errorf("all ping attempts failed")
	}

	return nil
}

func runMachinesStatus(cmd *cobra.Command, args []string) error {
	machineName := args[0]
	detailed, _ := cmd.Flags().GetBool("detailed")
	checkServices, _ := cmd.Flags().GetBool("services")

	// Load configuration to find machine
	cfg, err := config.ParseConfig(".")
	if err != nil {
		return fmt.Errorf("error loading configuration: %w", err)
	}

	// Find machine in configuration
	var targetMachine *config.Machine
	for _, machine := range cfg.Machines {
		if machine.Name == machineName {
			targetMachine = &machine
			break
		}
	}

	if targetMachine == nil {
		return fmt.Errorf("machine '%s' not found in configuration", machineName)
	}

	fmt.Printf("Checking status of machine: %s (%s)\n", machineName, targetMachine.Host)

	// Create SSH client
	client, err := ssh.NewSSHClient(targetMachine, 30)

	if err != nil {
		fmt.Printf("❌ Connection failed: %v\n", err)
		return fmt.Errorf("failed to connect to machine: %w", err)
	}
	defer client.Close()

	fmt.Printf("✅ SSH connection: OK\n")

	// Basic system information
	if detailed {
		fmt.Println("\nSystem Information:")

		// Uptime
		uptime, err := client.ExecuteCommand("uptime")
		if err == nil {
			fmt.Printf("  Uptime: %s", uptime)
		}

		// OS information
		osInfo, err := client.ExecuteCommand("cat /etc/os-release | grep PRETTY_NAME")
		if err == nil {
			fmt.Printf("  OS: %s", osInfo)
		}

		// Memory usage
		memInfo, err := client.ExecuteCommand("free -h | grep Mem")
		if err == nil {
			fmt.Printf("  Memory: %s", memInfo)
		}

		// Disk usage
		diskInfo, err := client.ExecuteCommand("df -h / | tail -1")
		if err == nil {
			fmt.Printf("  Disk: %s", diskInfo)
		}
	}

	// Service status check
	if checkServices {
		fmt.Println("\nService Status:")
		services := []string{"sshd", "systemd", "cron"}

		for _, service := range services {
			status, err := client.ExecuteCommand(fmt.Sprintf("systemctl is-active %s", service))
			if err == nil {
				if status == "active" {
					fmt.Printf("  %s: ✅ Active\n", service)
				} else {
					fmt.Printf("  %s: ❌ %s\n", service, status)
				}
			} else {
				fmt.Printf("  %s: ❓ Unknown\n", service)
			}
		}
	}

	fmt.Println("\n✅ Machine status check completed")
	return nil
}

func outputMachinesTable(machines []config.Machine, detailed bool) error {
	fmt.Printf("Found %d machine(s):\n\n", len(machines))

	if detailed {
		fmt.Printf("%-20s %-15s %-10s %-15s %-20s\n", "NAME", "HOST", "PORT", "USER", "KEY FILE")
		fmt.Println(string(make([]byte, 80)))
		for i := range machines {
			keyFile := "password"
			if machines[i].KeyFile != "" {
				keyFile = machines[i].KeyFile
			}
			fmt.Printf("%-20s %-15s %-10d %-15s %-20s\n",
				machines[i].Name,
				machines[i].Host,
				machines[i].Port,
				machines[i].User,
				keyFile)
		}
	} else {
		for _, machine := range machines {
			fmt.Printf("  %s (%s:%d)\n", machine.Name, machine.Host, machine.Port)
		}
	}

	return nil
}

func outputMachinesJSON(_ []config.Machine, _ bool) error {
	// Placeholder for JSON output
	fmt.Println("JSON output would be implemented here")
	return nil
}

func outputMachinesYAML(_ []config.Machine, _ bool) error {
	// Placeholder for YAML output
	fmt.Println("YAML output would be implemented here")
	return nil
}

func runMachinesAction(_ *cobra.Command, args []string) error {
	source := args[0]
	logger := logging.GetLogger()

	// Check if file exists
	if _, err := os.Stat(source); os.IsNotExist(err) {
		return fmt.Errorf("configuration file does not exist: %s", source)
	}

	// Parse configuration
	cfg, err := config.ParseConfig(source)
	if err != nil {
		return fmt.Errorf("failed to parse configuration: %w", err)
	}

	logger.Info("Starting machine actions execution",
		logging.String("config_file", source),
		logging.Int("action_count", len(cfg.Actions)),
		logging.Int("machine_count", len(cfg.Machines)),
	)

	// Execute configuration using SSH executor
	return ssh.ExecuteConfig(cfg)
}
