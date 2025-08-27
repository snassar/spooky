package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"spooky-test-bundle-generator/internal/generator"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "spooky-test-bundle-generator",
	Short: "Generate test bundles for Spooky with containerized SSH servers",
	Long: `Generate test bundles for Spooky automation testing.

This tool creates complete test environments including:
- Containerfiles with SSH servers
- Spooky project configurations
- Network setup and IP assignment
- Test templates and files
- Validation and testing scripts

Each bundle is a self-contained test environment that can be used
to test Spooky's SSH connectivity, facts gathering, template rendering,
file synchronization, and action execution.`,
}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a test bundle from a profile",
	Long: `Generate a test bundle from an HCL profile configuration.

The profile defines the container configuration, SSH settings,
Spooky project structure, and test scenarios.`,
	Args: cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		profilePath, _ := cmd.Flags().GetString("profile")
		outputPath, _ := cmd.Flags().GetString("output")

		if profilePath == "" {
			fmt.Fprintf(os.Stderr, "Error: --profile is required\n")
			os.Exit(1)
		}

		if outputPath == "" {
			fmt.Fprintf(os.Stderr, "Error: --output is required\n")
			os.Exit(1)
		}

		bundleGen := generator.NewBundleGenerator()
		if err := bundleGen.GenerateBundle(profilePath, outputPath); err != nil {
			fmt.Fprintf(os.Stderr, "Error generating bundle: %v\n", err)
			os.Exit(1)
		}
	},
}

var listProfilesCmd = &cobra.Command{
	Use:   "list-profiles",
	Short: "List available test profiles",
	Long: `List all available test profiles in the profiles directory.

Profiles define different test scenarios including:
- Different OS distributions (Ubuntu, CentOS, Alpine)
- Various SSH configurations (permissive, strict)
- Multi-container clusters
- Different authentication methods`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := listProfiles(); err != nil {
			fmt.Fprintf(os.Stderr, "Error listing profiles: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	generateCmd.Flags().String("profile", "", "Path to HCL profile file")
	generateCmd.Flags().String("output", "", "Output directory for generated bundle")
	generateCmd.MarkFlagRequired("profile")
	generateCmd.MarkFlagRequired("output")

	rootCmd.AddCommand(generateCmd)
	rootCmd.AddCommand(listProfilesCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func listProfiles() error {
	profilesDir := "profiles"
	if _, err := os.Stat(profilesDir); os.IsNotExist(err) {
		fmt.Println("No profiles directory found. Create profiles/ directory with .hcl files.")
		return nil
	}

	entries, err := os.ReadDir(profilesDir)
	if err != nil {
		return fmt.Errorf("failed to read profiles directory: %w", err)
	}

	fmt.Println("Available test profiles:")
	fmt.Println()

	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".hcl" {
			profileName := strings.TrimSuffix(entry.Name(), ".hcl")
			fmt.Printf("  %s\n", profileName)
		}
	}

	fmt.Println()
	fmt.Println("To generate a bundle:")
	fmt.Println("  ./spooky-test-bundle-generator generate --profile profiles/PROFILE.hcl --output bundles/PROFILE")

	return nil
}
