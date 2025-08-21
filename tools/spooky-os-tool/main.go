package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"spooky/internal/utilities"
)

var rootCmd = &cobra.Command{
	Use:   "spooky-os",
	Short: "OS detection and path configuration tool for spooky",
	Long: `spooky-os is a utility tool that detects the current operating system
and provides OS-specific path configurations for the spooky automation platform.

It supports Linux, macOS, Windows 11+, and BSD systems.`,
}

var detectCmd = &cobra.Command{
	Use:   "detect",
	Short: "Detect current operating system information",
	Run: func(cmd *cobra.Command, args []string) {
		osInfo, err := utilities.DetectOS()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error detecting OS: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("🖥️  OS Detection Results:\n")
		fmt.Printf("   OS: %s\n", osInfo.OS)
		fmt.Printf("   Version: %s\n", osInfo.Version)
		fmt.Printf("   Architecture: %s\n", osInfo.Arch)
		fmt.Printf("   Distribution: %s\n", osInfo.Distro)
		fmt.Printf("   Kernel: %s\n", osInfo.Kernel)
		fmt.Printf("   Container: %v\n", osInfo.IsContainer)
		fmt.Printf("   WSL: %v\n", osInfo.IsWSL)
		fmt.Printf("   Root: %v\n", osInfo.IsRoot)
		fmt.Printf("   Supported: %v\n", osInfo.IsSupported())
	},
}

var pathsCmd = &cobra.Command{
	Use:   "paths [app-name]",
	Short: "Get OS-specific path configuration",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		appName := "spooky"
		if len(args) > 0 {
			appName = args[0]
		}

		config, err := utilities.GetPathConfig(appName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting path config: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("📁 Path Configuration for '%s':\n", appName)
		fmt.Printf("   Config Directory: %s\n", config.ConfigDir)
		fmt.Printf("   Log Directory: %s\n", config.LogDir)
		fmt.Printf("   Cache Directory: %s\n", config.CacheDir)
		fmt.Printf("   Data Directory: %s\n", config.DataDir)
		fmt.Printf("   Temp Directory: %s\n", config.TempDir)
		fmt.Printf("   User Home: %s\n", config.UserHomeDir)
		fmt.Printf("\n📄 File Paths:\n")
		fmt.Printf("   Config File: %s\n", config.ConfigFile)
		fmt.Printf("   Log File: %s\n", config.LogFile)
		fmt.Printf("   Cache File: %s\n", config.CacheFile)
	},
}

var setupCmd = &cobra.Command{
	Use:   "setup [app-name]",
	Short: "Create OS-specific directories for the application",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		appName := "spooky"
		if len(args) > 0 {
			appName = args[0]
		}

		config, err := utilities.GetPathConfig(appName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting path config: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("🔧 Setting up directories for '%s':\n", appName)
		
		err = utilities.EnsureDirectories(config)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating directories: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("✅ Successfully created directories:\n")
		fmt.Printf("   %s\n", config.ConfigDir)
		fmt.Printf("   %s\n", config.DataDir)
		fmt.Printf("   %s\n", config.CacheDir)
		fmt.Printf("   %s\n", config.LogDir)
	},
}

func init() {
	rootCmd.AddCommand(detectCmd)
	rootCmd.AddCommand(pathsCmd)
	rootCmd.AddCommand(setupCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
