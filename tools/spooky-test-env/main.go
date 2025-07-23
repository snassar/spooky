package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

type TestEnvCommander interface {
	Output(name string, args ...string) ([]byte, error)
	Run(name string, args ...string) error
}

type TestEnvRealCommander struct{}

func (TestEnvRealCommander) Output(name string, args ...string) ([]byte, error) {
	return exec.Command(name, args...).Output()
}

func (TestEnvRealCommander) Run(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

var cmd TestEnvCommander = TestEnvRealCommander{}

// Global flags
var sshKeyPath string

func main() {
	var rootCmd = &cobra.Command{
		Use:   "spooky-test-env",
		Short: "Manage spooky test environment with Podman and Quadlet",
		Long: `spooky-test-env is a tool for managing the Podman-based test environment for spooky.

It provides commands to check prerequisites, start/stop containers using Quadlet,
and clean up resources. The test environment includes SSH-enabled containers
for testing spooky's remote management capabilities.`,
		PersistentPreRunE: func(_ *cobra.Command, _ []string) error {
			// Set default SSH key path if not provided
			if sshKeyPath == "" {
				sshKeyPath = getDefaultSSHKeyPath()
			}
			return nil
		},
	}

	// Global flags
	rootCmd.PersistentFlags().StringVarP(&sshKeyPath, "ssh-key", "k", "", "Path to SSH private key (default: ~/.ssh/id_ed25519 or SPOOKY_TEST_SSH_KEY env var)")

	rootCmd.AddCommand(preflightCmd)
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(stopCmd)
	rootCmd.AddCommand(cleanupCmd)
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(buildCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func getDefaultSSHKeyPath() string {
	// Check environment variable first
	if envPath := os.Getenv("SPOOKY_TEST_SSH_KEY"); envPath != "" {
		return envPath
	}

	// Default to ~/.ssh/id_ed25519
	return filepath.Join(os.Getenv("HOME"), ".ssh", "id_ed25519")
}

var quietFlag bool

var preflightCmd = &cobra.Command{
	Use:   "preflight",
	Short: "Check prerequisites for the test environment",
	Long:  `Check that all required tools (podman, systemd, quadlet) are available and working.`,
	RunE: func(_ *cobra.Command, _ []string) error {
		return runPreflight(quietFlag)
	},
}

func init() {
	preflightCmd.Flags().BoolVarP(&quietFlag, "quiet", "q", false, "Suppress output, only return exit code")
}

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build container images for the test environment",
	Long:  `Build the container images needed for the test environment.`,
	RunE: func(_ *cobra.Command, _ []string) error {
		return runBuild()
	},
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the test environment",
	Long:  `Start the Podman containers for the spooky test environment using Quadlet.`,
	RunE: func(_ *cobra.Command, _ []string) error {
		return runStart()
	},
}

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop the test environment",
	Long:  `Stop the Podman containers for the spooky test environment.`,
	RunE: func(_ *cobra.Command, _ []string) error {
		return runStop()
	},
}

var cleanupCmd = &cobra.Command{
	Use:   "cleanup",
	Short: "Clean up test environment resources",
	Long:  `Remove containers, networks, and other resources created by the test environment.`,
	RunE: func(_ *cobra.Command, _ []string) error {
		return runCleanup()
	},
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show test environment status",
	Long:  `Show the current status of the test environment containers and networks.`,
	RunE: func(_ *cobra.Command, _ []string) error {
		return runStatus()
	},
}

func runPreflight(quiet bool) error {
	// Check programs
	podmanAvailable := "not found"
	podmanVersion := ""
	if _, err := exec.LookPath("podman"); err == nil {
		if out, err := exec.Command("podman", "--version").Output(); err == nil {
			podmanAvailable = "available"
			podmanVersion = strings.TrimSpace(string(out))
		} else {
			podmanAvailable = "found but not working"
		}
	}

	// Check systemd (Linux only)
	systemdAvailable := "not applicable (not Linux)"
	systemdVersion := ""
	if runtime.GOOS == "linux" {
		if _, err := exec.LookPath("systemctl"); err == nil {
			if out, err := exec.Command("systemctl", "--version").Output(); err == nil {
				systemdAvailable = "available"
				// systemctl --version outputs multiple lines, take the first line
				lines := strings.Split(strings.TrimSpace(string(out)), "\n")
				if len(lines) > 0 {
					systemdVersion = strings.TrimSpace(lines[0])
				}
			} else {
				systemdAvailable = "found but not working"
			}
		} else {
			systemdAvailable = "not found"
		}
	}

	// Check rootless Podman (Linux only)
	rootlessStatus := "not applicable (not Linux)"
	if runtime.GOOS == "linux" {
		if out, err := exec.Command("podman", "info", "--log-level=error", "--format", "{{.Host.Security.Rootless}}").Output(); err == nil {
			if strings.TrimSpace(string(out)) == "true" {
				rootlessStatus = "yes"
			} else {
				rootlessStatus = "no"
			}
		} else {
			rootlessStatus = "failed to check"
		}
	}

	// Check Quadlet support
	quadletStatus := "not found"
	if err := exec.Command("podman", "quadlet", "--help").Run(); err == nil {
		quadletStatus = "available"
	}

	// Check SSH key
	sshKeyStatus := "not found"
	if _, err := os.Stat(sshKeyPath); err == nil {
		sshKeyStatus = "available"
	}

	// Determine overall status
	allGood := true
	if podmanAvailable != "available" {
		allGood = false
	}
	if runtime.GOOS == "linux" {
		if rootlessStatus != "yes" || quadletStatus == "not found" {
			allGood = false
		}
	}
	if sshKeyStatus == "not found" {
		allGood = false
	}

	// Output results (unless quiet mode)
	if !quiet {
		fmt.Println("Checking requirements for testing:")
		fmt.Println()

		fmt.Println("Programs:")
		fmt.Printf("* podman: %s;\t\t%s\n", podmanAvailable, podmanVersion)
		fmt.Printf("* systemd: %s;\t\t%s\n", systemdAvailable, systemdVersion)

		fmt.Println()
		fmt.Println("Environment:")
		fmt.Printf("* user can run podman rootless: %s;\n", rootlessStatus)
		fmt.Printf("* quadlet support: %s;\n", quadletStatus)
		fmt.Printf("* SSH key: %s (%s);\n", sshKeyStatus, sshKeyPath)

		fmt.Println()
		if allGood {
			fmt.Println("spooky-test-env requirements satisfied.")
		} else {
			fmt.Println("❌ Some requirements are not met.")
			if sshKeyStatus == "not found" {
				fmt.Printf("   Generate SSH key with: ssh-keygen -t ed25519 -f %s -N ''\n", sshKeyPath)
			}
		}
	}

	// Return appropriate exit code
	if allGood {
		return nil
	}
	return fmt.Errorf("requirements not satisfied")
}

func runBuild() error {
	fmt.Println("Building container images for test environment...")

	testEnvDir := getTestEnvDir()
	if err := os.Chdir(testEnvDir); err != nil {
		return fmt.Errorf("failed to change to test environment directory: %w", err)
	}

	// Generate SSH key if it doesn't exist
	if _, err := os.Stat(sshKeyPath); os.IsNotExist(err) {
		fmt.Printf("Generating SSH key at %s...\n", sshKeyPath)
		if err := cmd.Run("ssh-keygen", "-t", "ed25519", "-f", sshKeyPath, "-N", ""); err != nil {
			return fmt.Errorf("failed to generate SSH key: %w", err)
		}
	}

	// Copy SSH public key content to build context as authorized_keys
	sshPubKeyPath := sshKeyPath + ".pub"
	if err := cmd.Run("cp", sshPubKeyPath, "authorized_keys"); err != nil {
		return fmt.Errorf("failed to copy SSH public key: %w", err)
	}

	// Build SSH container image
	fmt.Println("Building SSH container image...")
	if err := cmd.Run("podman", "build", "-f", "Containerfile.ssh", "-t", "spooky-test-ssh", "."); err != nil {
		return fmt.Errorf("failed to build SSH container: %w", err)
	}

	// Build no-SSH container image
	fmt.Println("Building no-SSH container image...")
	if err := cmd.Run("podman", "build", "-f", "Containerfile.no-ssh", "-t", "spooky-test-no-ssh", "."); err != nil {
		return fmt.Errorf("failed to build no-SSH container: %w", err)
	}

	// Build SSH-no-key container image
	fmt.Println("Building SSH-no-key container image...")
	if err := cmd.Run("podman", "build", "-f", "Containerfile.ssh-no-key", "-t", "spooky-test-ssh-no-key", "."); err != nil {
		return fmt.Errorf("failed to build SSH-no-key container: %w", err)
	}

	// Clean up the copied public key file
	if err := os.Remove("authorized_keys"); err != nil && !os.IsNotExist(err) {
		fmt.Printf("Warning: failed to remove temporary public key file: %v\n", err)
	}

	fmt.Println("All container images built successfully!")
	return nil
}

func runStart() error {
	fmt.Println("Starting spooky test environment...")

	testEnvDir := getTestEnvDir()
	if err := os.Chdir(testEnvDir); err != nil {
		return fmt.Errorf("failed to change to test environment directory: %w", err)
	}

	// Check if images exist
	images := []string{"spooky-test-ssh", "spooky-test-no-ssh", "spooky-test-ssh-no-key"}
	for _, image := range images {
		if out, err := cmd.Output("podman", "images", "--format", "{{.Repository}}:{{.Tag}}", image); err != nil || strings.TrimSpace(string(out)) == "" {
			fmt.Printf("Image %s not found, building first...\n", image)
			if err := runBuild(); err != nil {
				return fmt.Errorf("failed to build images: %w", err)
			}
			break
		}
	}

	// Start containers using Quadlet
	fmt.Println("Starting containers with Quadlet...")

	// Start 7 working SSH containers
	for i := 1; i <= 7; i++ {
		port := 2220 + i
		containerName := fmt.Sprintf("spooky-test-server-%d", i)

		fmt.Printf("Starting working SSH container %d on port %d...\n", i, port)
		if err := startContainerWithQuadlet(containerName, "spooky-test-ssh", port); err != nil {
			return fmt.Errorf("failed to start %s: %w", containerName, err)
		}
	}

	// Start container with no SSH running (port 2228)
	fmt.Println("Starting container with no SSH on port 2228...")
	if err := startContainerWithQuadlet("spooky-test-no-ssh", "spooky-test-no-ssh", 2228); err != nil {
		return fmt.Errorf("failed to start spooky-test-no-ssh: %w", err)
	}

	// Start container with SSH but no authorized key (port 2229)
	fmt.Println("Starting container with SSH but no key on port 2229...")
	if err := startContainerWithQuadlet("spooky-test-ssh-no-key", "spooky-test-ssh-no-key", 2229); err != nil {
		return fmt.Errorf("failed to start spooky-test-ssh-no-key: %w", err)
	}

	// Wait for all containers to be ready
	fmt.Println("Waiting for all containers to be ready...")
	for i := 1; i <= 30; i++ {
		workingCount := 0
		noSSHCount := 0
		noKeyCount := 0

		// Count working containers
		if out, err := cmd.Output("podman", "ps", "--filter", "name=spooky-test-server", "--format", "{{.Status}}"); err == nil {
			workingCount = strings.Count(string(out), "Up")
		}

		// Count no-SSH container
		if out, err := cmd.Output("podman", "ps", "--filter", "name=spooky-test-no-ssh", "--format", "{{.Status}}"); err == nil {
			noSSHCount = strings.Count(string(out), "Up")
		}

		// Count SSH-no-key container
		if out, err := cmd.Output("podman", "ps", "--filter", "name=spooky-test-ssh-no-key", "--format", "{{.Status}}"); err == nil {
			noKeyCount = strings.Count(string(out), "Up")
		}

		totalRunning := workingCount + noSSHCount + noKeyCount

		if totalRunning == 9 {
			fmt.Printf("All 9 containers are running (attempt %d/30)\n", i)
			break
		}

		fmt.Printf("Waiting for containers to start... (%d/9 running, attempt %d/30)\n", totalRunning, i)
		time.Sleep(2 * time.Second)
	}

	// Show container IPs
	fmt.Println("Container IPs:")
	containers := []string{"spooky-test-server-1", "spooky-test-server-2", "spooky-test-server-3", "spooky-test-server-4", "spooky-test-server-5", "spooky-test-server-6", "spooky-test-server-7", "spooky-test-no-ssh", "spooky-test-ssh-no-key"}
	for _, container := range containers {
		ip, err := getContainerIP(container)
		if err != nil {
			fmt.Printf("  %s: not available\n", container)
		} else {
			fmt.Printf("  %s: %s\n", container, ip)
		}
	}

	fmt.Println("Test environment ready!")
	return nil
}

func startContainerWithQuadlet(containerName, imageName string, port int) error {
	// Create the spooky-test network if it doesn't exist
	_ = cmd.Run("podman", "network", "create", "spooky-test")
	// Network might already exist, that's okay

	// Assign static IP based on container number
	var staticIP string
	switch containerName {
	case "spooky-test-server-1":
		staticIP = "10.1.10.1"
	case "spooky-test-server-2":
		staticIP = "10.1.10.2"
	case "spooky-test-server-3":
		staticIP = "10.1.10.3"
	case "spooky-test-server-4":
		staticIP = "10.1.10.4"
	case "spooky-test-server-5":
		staticIP = "10.1.10.5"
	case "spooky-test-server-6":
		staticIP = "10.1.10.6"
	case "spooky-test-server-7":
		staticIP = "10.1.10.7"
	case "spooky-test-no-ssh":
		staticIP = "10.1.10.8"
	case "spooky-test-ssh-no-key":
		staticIP = "10.1.10.9"
	default:
		staticIP = ""
	}

	args := []string{
		"run", "-d",
		"--name", containerName,
		"--network", "spooky-test",
	}

	// Add static IP if specified
	if staticIP != "" {
		args = append(args, "--ip", staticIP)
	}

	args = append(args,
		"-p", fmt.Sprintf("%d:22", port),
		"--dns", "8.8.8.8",
		"--dns", "8.8.4.4",
		fmt.Sprintf("localhost/%s:latest", imageName),
	)

	if err := cmd.Run("podman", args...); err != nil {
		return fmt.Errorf("failed to start container: %w", err)
	}

	return nil
}

func runStop() error {
	fmt.Println("Stopping spooky test environment...")

	testEnvDir := getTestEnvDir()
	if err := os.Chdir(testEnvDir); err != nil {
		return fmt.Errorf("failed to change to test environment directory: %w", err)
	}

	// Stop containers using podman
	containers := []string{"spooky-test-server-1", "spooky-test-server-2", "spooky-test-server-3", "spooky-test-server-4", "spooky-test-server-5", "spooky-test-server-6", "spooky-test-server-7", "spooky-test-no-ssh", "spooky-test-ssh-no-key"}

	for _, container := range containers {
		fmt.Printf("Stopping %s...\n", container)
		if err := cmd.Run("podman", "stop", container); err != nil {
			fmt.Printf("Warning: failed to stop %s: %v\n", container, err)
		}
	}

	fmt.Println("Test environment stopped!")
	return nil
}

func runCleanup() error {
	fmt.Println("Cleaning up spooky test environment...")

	// Stop containers first
	if err := runStop(); err != nil {
		fmt.Printf("Warning: failed to stop containers: %v\n", err)
	}

	// Remove containers
	containers := []string{"spooky-test-server-1", "spooky-test-server-2", "spooky-test-server-3", "spooky-test-server-4", "spooky-test-server-5", "spooky-test-server-6", "spooky-test-server-7", "spooky-test-no-ssh", "spooky-test-ssh-no-key"}

	for _, container := range containers {
		fmt.Printf("Removing container %s...\n", container)

		// Force remove the container
		if err := cmd.Run("podman", "rm", "-f", container); err != nil {
			fmt.Printf("Container %s not found or already removed\n", container)
		}
	}

	fmt.Println("Cleanup completed!")
	return nil
}

func runStatus() error {
	fmt.Println("Spooky test environment status:")
	fmt.Println()

	// Show container status
	fmt.Println("Container Status:")
	containers := []string{"spooky-test-server-1", "spooky-test-server-2", "spooky-test-server-3", "spooky-test-server-4", "spooky-test-server-5", "spooky-test-server-6", "spooky-test-server-7", "spooky-test-no-ssh", "spooky-test-ssh-no-key"}

	for _, container := range containers {
		output, err := cmd.Output("podman", "ps", "--filter", "name="+container, "--format", "{{.Names}}: {{.Status}}")
		if err != nil || strings.TrimSpace(string(output)) == "" {
			fmt.Printf("  %s: not running\n", container)
		} else {
			fmt.Printf("  %s\n", strings.TrimSpace(string(output)))
		}
	}
	fmt.Println()

	// Show container IPs if running
	fmt.Println("Container IPs:")
	for _, container := range containers {
		ip, err := getContainerIP(container)
		if err != nil {
			fmt.Printf("  %s: not available\n", container)
		} else {
			fmt.Printf("  %s: %s\n", container, ip)
		}
	}

	return nil
}

func getTestEnvDir() string {
	// Get the current working directory (where the command is run from)
	currentDir, err := os.Getwd()
	if err != nil {
		// Fallback to current directory
		return "."
	}

	// Navigate to the test environment directory in examples
	return filepath.Join(currentDir, "examples", "test-environment")
}

func getContainerIP(containerName string) (string, error) {
	// Use jq to parse the JSON output since the template syntax doesn't work with hyphens
	output, err := cmd.Output("podman", "inspect", containerName)
	if err != nil {
		return "", fmt.Errorf("failed to inspect container %s: %w", containerName, err)
	}

	// Use jq to extract the IP address
	cmd2 := exec.Command("jq", "-r", ".[0].NetworkSettings.Networks[\"spooky-test\"].IPAddress")
	cmd2.Stdin = strings.NewReader(string(output))
	ipOutput, err := cmd2.Output()
	if err != nil {
		return "", fmt.Errorf("failed to parse IP address for container %s: %w", containerName, err)
	}

	ip := strings.TrimSpace(string(ipOutput))
	if ip == "" || ip == "null" {
		return "", fmt.Errorf("no IP address found for container %s", containerName)
	}

	return ip, nil
}
