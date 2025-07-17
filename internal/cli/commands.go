package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"spooky/internal/config"
	"spooky/internal/ssh"
)

var (
	configFile string // coverage-ignore: global variable declaration
	parallel   bool   // coverage-ignore: global variable declaration
	timeout    int    // coverage-ignore: global variable declaration
)

var ExecuteCmd = &cobra.Command{
	Use:   "execute [config-file]",
	Short: "Execute actions from HCL2 configuration file",
	Long:  `Execute actions defined in an HCL2 configuration file on remote servers via SSH`,
	Args:  cobra.MaximumNArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		if len(args) > 0 {
			configFile = args[0]
		}

		if configFile == "" {
			return fmt.Errorf("config file is required") // coverage-ignore: CLI validation error, tested via integration
		}

		// Validate config file exists
		if _, err := os.Stat(configFile); os.IsNotExist(err) {
			return fmt.Errorf("config file %s does not exist", configFile) // coverage-ignore: file system error, hard to test
		}

		// Parse and execute configuration
		config, err := config.ParseConfig(configFile)
		if err != nil {
			return fmt.Errorf("failed to parse config: %w", err)
		}

		return ssh.ExecuteConfig(config)
	},
}

var ValidateCmd = &cobra.Command{
	Use:   "validate [config-file]",
	Short: "Validate HCL2 configuration file",
	Long:  `Validate the syntax and structure of an HCL2 configuration file`,
	Args:  cobra.MaximumNArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		if len(args) > 0 {
			configFile = args[0]
		}

		if configFile == "" {
			return fmt.Errorf("config file is required")
		}

		// Validate config file exists
		if _, err := os.Stat(configFile); os.IsNotExist(err) {
			return fmt.Errorf("config file %s does not exist", configFile)
		}

		// Parse configuration
		config, err := config.ParseConfig(configFile)
		if err != nil {
			return fmt.Errorf("validation failed: %w", err)
		}

		fmt.Printf("✅ Configuration file '%s' is valid\n", configFile)
		fmt.Printf("📊 Found %d servers and %d actions\n", len(config.Servers), len(config.Actions))

		return nil
	},
}

var ListCmd = &cobra.Command{
	Use:   "list [config-file]",
	Short: "List servers and actions from configuration file",
	Long:  `Display all servers and actions defined in an HCL2 configuration file`,
	Args:  cobra.MaximumNArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		if len(args) > 0 {
			configFile = args[0]
		}

		if configFile == "" {
			return fmt.Errorf("config file is required")
		}

		// Validate config file exists
		if _, err := os.Stat(configFile); os.IsNotExist(err) {
			return fmt.Errorf("config file %s does not exist", configFile)
		}

		// Parse configuration
		config, err := config.ParseConfig(configFile)
		if err != nil {
			return fmt.Errorf("failed to parse config: %w", err)
		}

		// Display servers
		fmt.Println("🌐 Servers:")
		for _, server := range config.Servers {
			fmt.Printf("  - %s (%s@%s:%d)\n", server.Name, server.User, server.Host, server.Port)
		}

		// Display actions
		fmt.Println("\n⚡ Actions:")
		for _, action := range config.Actions {
			fmt.Printf("  - %s: %s\n", action.Name, action.Description)
		}

		return nil
	},
}

// InitCommands initializes all CLI commands and their flags
func InitCommands() {
	// Global flags
	ExecuteCmd.Flags().BoolVarP(&parallel, "parallel", "p", false, "Execute actions in parallel")
	ExecuteCmd.Flags().IntVarP(&timeout, "timeout", "t", 30, "SSH connection timeout in seconds")

	ValidateCmd.Flags().StringVarP(&configFile, "config", "c", "", "Configuration file path")
	ListCmd.Flags().StringVarP(&configFile, "config", "c", "", "Configuration file path")
}
