package main

import (
	"fmt"
	"os"
	"time"

	"spooky/internal/utilities"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "spooky-demo",
	Short: "Demonstration tool for spooky project logging and configuration",
	Long: `spooky-demo demonstrates how spooky handles:
- Project-specific logging with timestamps
- OS-appropriate configuration file management
- Cross-platform path handling`,
}

var demoProjectCmd = &cobra.Command{
	Use:   "project [project-name]",
	Short: "Demonstrate project logging functionality",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		projectName := args[0]

		fmt.Printf("🚀 Starting project logging demo for: %s\n\n", projectName)

		// Create project logger
		projectLogger, err := utilities.NewProjectLogger(projectName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating project logger: %v\n", err)
			os.Exit(1)
		}
		defer projectLogger.Close()

		// Log project start
		startTime := time.Now()
		projectLogger.LogProjectStart("/path/to/project", []string{"--verbose", "--dry-run"})

		fmt.Printf("📝 Project logger created:\n")
		fmt.Printf("   Run ID: %s\n", projectLogger.GetRunID())
		fmt.Printf("   Log File: %s\n", projectLogger.GetLogFile())
		fmt.Printf("   Timestamp: %s\n\n", projectLogger.GetTimestamp().Format(time.RFC3339))

		// Simulate some tasks
		fmt.Println("🔧 Simulating project tasks...")

		// Task 1: Schema validation
		taskStart := time.Now()
		time.Sleep(100 * time.Millisecond) // Simulate work
		projectLogger.LogTask("validate-schemas", "validation", time.Since(taskStart), true, map[string]interface{}{
			"schemas_validated": 12,
			"errors_found":      0,
		})

		// Task 2: HCL parsing
		taskStart = time.Now()
		time.Sleep(200 * time.Millisecond) // Simulate work
		projectLogger.LogTask("parse-hcl", "parsing", time.Since(taskStart), true, map[string]interface{}{
			"files_parsed": 5,
			"blocks_found": 23,
		})

		// Task 3: Template rendering (with error)
		taskStart = time.Now()
		time.Sleep(150 * time.Millisecond) // Simulate work
		projectLogger.LogTask("render-templates", "rendering", time.Since(taskStart), false, map[string]interface{}{
			"templates_rendered": 3,
			"errors":             1,
			"error_details":      "Template 'nginx.conf' not found",
		})

		// Log an error
		projectLogger.LogError(fmt.Errorf("template not found: nginx.conf"), "template rendering", map[string]interface{}{
			"template_name": "nginx.conf",
			"search_paths":  []string{"/templates", "/etc/templates"},
		})

		// Log project end
		duration := time.Since(startTime)
		projectLogger.LogProjectEnd(duration, false) // false because of the error

		fmt.Printf("✅ Project logging demo completed!\n")
		fmt.Printf("   Duration: %v\n", duration)
		fmt.Printf("   Log file: %s\n\n", projectLogger.GetLogFile())

		// Show log summary
		summary, err := utilities.GetProjectLogSummary(projectName)
		if err != nil {
			fmt.Printf("⚠️  Could not get log summary: %v\n", err)
		} else {
			fmt.Printf("📊 Log Summary:\n")
			fmt.Printf("   Project: %s\n", summary.ProjectName)
			fmt.Printf("   Total Runs: %d\n", summary.TotalRuns)
			if summary.LatestLogFile != "" {
				fmt.Printf("   Latest Log: %s\n", summary.LatestLogFile)
			}
		}
	},
}

var demoConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Demonstrate configuration management functionality",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("⚙️  Starting configuration management demo...\n")

		// Create config manager
		configManager, err := utilities.NewConfigManager()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating config manager: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("📁 Configuration paths:\n")
		fmt.Printf("   Config Directory: %s\n", configManager.GetConfigDir())
		fmt.Printf("   Config File: %s\n\n", configManager.GetConfigPath())

				// Show effective config info
		effectiveInfo, err := configManager.GetEffectiveConfigInfo("")
		if err != nil {
			fmt.Printf("⚠️  Could not get effective config info: %v\n", err)
		} else {
			fmt.Printf("📄 Effective Config Info:\n")
			fmt.Printf("   Source: %s\n", effectiveInfo.Source)
			fmt.Printf("   Size: %d bytes\n", effectiveInfo.Size)
			if effectiveInfo.ModTime != "" {
				fmt.Printf("   Modified: %s\n", effectiveInfo.ModTime)
			}
		}

		// Check if user config exists
		if configManager.ConfigExists() {
			fmt.Println("✅ User configuration file exists")
			
			// Show user config info
			info, err := configManager.GetConfigInfo()
			if err != nil {
				fmt.Printf("⚠️  Could not get config info: %v\n", err)
			} else {
				fmt.Printf("📄 User Config Info:\n")
				fmt.Printf("   Size: %d bytes\n", info.Size)
				fmt.Printf("   Modified: %s\n", info.ModTime)
			}
			
			// Create backup
			backupPath, err := configManager.BackupConfig()
			if err != nil {
				fmt.Printf("⚠️  Could not create backup: %v\n", err)
			} else {
				fmt.Printf("💾 Backup created: %s\n", backupPath)
			}
		} else {
			fmt.Println("📝 No user configuration found - using embedded default")
			
			// Show that we can create a user config
			fmt.Println("💡 You can create a user config with: spooky-demo config create")
		}

		// List config files
		configFiles, err := configManager.ListConfigFiles()
		if err != nil {
			fmt.Printf("⚠️  Could not list config files: %v\n", err)
		} else {
			fmt.Printf("\n📋 Configuration files:\n")
			for _, file := range configFiles {
				fmt.Printf("   %s\n", file)
			}
		}

		// Validate config
		if err := configManager.ValidateConfig(""); err != nil {
			fmt.Printf("⚠️  Config validation failed: %v\n", err)
		} else {
			fmt.Println("\n✅ Configuration is valid!")
		}
	},
}

var demoPathsCmd = &cobra.Command{
	Use:   "paths",
	Short: "Show OS-specific paths for spooky",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🖥️  OS-specific paths for spooky:\n")

		// Get OS info
		osInfo, err := utilities.DetectOS()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error detecting OS: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("OS: %s (%s)\n", osInfo.OS, osInfo.Distro)
		fmt.Printf("Architecture: %s\n", osInfo.Arch)
		fmt.Printf("Container: %v\n", osInfo.IsContainer)
		fmt.Printf("WSL: %v\n\n", osInfo.IsWSL)

		// Get path config
		config, err := utilities.GetPathConfig("spooky")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting path config: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("📁 Path Configuration:\n")
		fmt.Printf("   Config: %s\n", config.ConfigDir)
		fmt.Printf("   Logs:   %s\n", config.LogDir)
		fmt.Printf("   Cache:  %s\n", config.CacheDir)
		fmt.Printf("   Data:   %s\n", config.DataDir)
		fmt.Printf("   Temp:   %s\n\n", config.TempDir)

		fmt.Printf("📄 File Paths:\n")
		fmt.Printf("   Config: %s\n", config.ConfigFile)
		fmt.Printf("   Log:    %s\n", config.LogFile)
		fmt.Printf("   Cache:  %s\n", config.CacheFile)
	},
}

var createConfigCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a user configuration file from embedded default",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("📝 Creating user configuration from embedded default...")
		
		// Create config manager
		configManager, err := utilities.NewConfigManager()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating config manager: %v\n", err)
			os.Exit(1)
		}
		
		// Check if user config already exists
		if configManager.ConfigExists() {
			fmt.Println("⚠️  User configuration already exists")
			fmt.Printf("   Path: %s\n", configManager.GetConfigPath())
			fmt.Println("💡 Use 'spooky-demo config backup' to backup first")
			return
		}
		
		// Create default config
		if err := configManager.CreateDefaultConfig(); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating default config: %v\n", err)
			os.Exit(1)
		}
		
		fmt.Println("✅ User configuration created!")
		fmt.Printf("   Path: %s\n", configManager.GetConfigPath())
		fmt.Println("💡 You can now edit this file to customize your settings")
	},
}

var backupConfigCmd = &cobra.Command{
	Use:   "backup",
	Short: "Create a backup of the current user configuration",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("💾 Creating configuration backup...")
		
		// Create config manager
		configManager, err := utilities.NewConfigManager()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating config manager: %v\n", err)
			os.Exit(1)
		}
		
		// Check if user config exists
		if !configManager.ConfigExists() {
			fmt.Println("⚠️  No user configuration to backup")
			fmt.Println("💡 Create one first with 'spooky-demo config create'")
			return
		}
		
		// Create backup
		backupPath, err := configManager.BackupConfig()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating backup: %v\n", err)
			os.Exit(1)
		}
		
		fmt.Println("✅ Configuration backup created!")
		fmt.Printf("   Backup: %s\n", backupPath)
	},
}

var testCustomConfigCmd = &cobra.Command{
	Use:   "test-custom [config-file]",
	Short: "Test configuration with a custom config file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		customConfigFile := args[0]
		fmt.Printf("🧪 Testing custom configuration: %s\n\n", customConfigFile)
		
		// Create config manager
		configManager, err := utilities.NewConfigManager()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating config manager: %v\n", err)
			os.Exit(1)
		}
		
		// Show effective config info with custom file
		effectiveInfo, err := configManager.GetEffectiveConfigInfo(customConfigFile)
		if err != nil {
			fmt.Printf("⚠️  Could not get effective config info: %v\n", err)
		} else {
			fmt.Printf("📄 Effective Config Info:\n")
			fmt.Printf("   Source: %s\n", effectiveInfo.Source)
			fmt.Printf("   Config File: %s\n", effectiveInfo.ConfigFile)
			fmt.Printf("   Size: %d bytes\n", effectiveInfo.Size)
			if effectiveInfo.ModTime != "" {
				fmt.Printf("   Modified: %s\n", effectiveInfo.ModTime)
			}
		}
		
		// Validate custom config
		if err := configManager.ValidateConfig(customConfigFile); err != nil {
			fmt.Printf("⚠️  Custom config validation failed: %v\n", err)
		} else {
			fmt.Println("✅ Custom configuration is valid!")
		}
		
		// Show config priority
		fmt.Printf("\n📋 Configuration Priority:\n")
		fmt.Printf("   1. Custom config file (--config flag)\n")
		fmt.Printf("   2. User config (~/.config/spooky/spooky.hcl)\n")
		fmt.Printf("   3. Embedded default\n")
	},
}

func init() {
	rootCmd.AddCommand(demoProjectCmd)
	rootCmd.AddCommand(demoConfigCmd)
	rootCmd.AddCommand(demoPathsCmd)
	demoConfigCmd.AddCommand(createConfigCmd)
	demoConfigCmd.AddCommand(backupConfigCmd)
	demoConfigCmd.AddCommand(testCustomConfigCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
