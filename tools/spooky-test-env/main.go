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

func main() {
	var rootCmd = &cobra.Command{
		Use:   "spooky-test-env",
		Short: "Manage spooky test environment with Podman",
		Long: `spooky-test-env is a tool for managing the Podman-based test environment for spooky.

It provides commands to check prerequisites, start/stop containers, and clean up resources.`,
	}

	rootCmd.AddCommand(preflightCmd)
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(stopCmd)
	rootCmd.AddCommand(cleanupCmd)
	rootCmd.AddCommand(statusCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

var quietFlag bool

var preflightCmd = &cobra.Command{
	Use:   "preflight",
	Short: "Check prerequisites for the test environment",
	Long:  `Check that all required tools (podman, systemd) are available and working.`,
	RunE: func(_ *cobra.Command, _ []string) error {
		return runPreflight(quietFlag)
	},
}

func init() {
	preflightCmd.Flags().BoolVarP(&quietFlag, "quiet", "q", false, "Suppress output, only return exit code")
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the test environment",
	Long:  `Start the Podman containers for the spooky test environment.`,
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
	Long:  `Show the current status of containers and networks in the test environment.`,
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
	quadletPath := "/usr/libexec/podman/quadlet"
	if _, err := os.Stat(quadletPath); err == nil {
		quadletStatus = "available"
	} else if err := exec.Command("podman", "quadlet", "--help").Run(); err == nil {
		quadletStatus = "available"
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

		fmt.Println()
		if allGood {
			fmt.Println("spooky-test-env requirements satisfied.")
		} else {
			fmt.Println("❌ Some requirements are not met.")
		}
	}

	// Return appropriate exit code
	if allGood {
		return nil
	}
	return fmt.Errorf("requirements not satisfied")
}

func runStart() error {
	fmt.Println("Starting spooky test environment...")

	testEnvDir := getTestEnvDir()
	if err := os.Chdir(testEnvDir); err != nil {
		return fmt.Errorf("failed to change to test environment directory: %w", err)
	}

	// Create the network if it doesn't exist
	fmt.Println("Creating spooky-test network...")
	if err := cmd.Run("podman", "network", "create", "spooky-test"); err != nil {
		fmt.Println("Network spooky-test already exists or creation failed (continuing)...")
	}

	// Start the containers using podman run
	fmt.Println("Starting containers with podman...")
	containers := []string{"spooky-server1", "spooky-server2", "spooky-server3"}

	for _, container := range containers {
		fmt.Printf("Starting %s...\n", container)
		// For now, use direct podman run commands instead of Quadlet
		// TODO: Implement proper Quadlet integration
		if err := startContainerDirectly(container); err != nil {
			return fmt.Errorf("failed to start %s: %w", container, err)
		}
	}

	// Wait a moment for containers to fully start
	fmt.Println("Waiting for containers to start...")
	time.Sleep(10 * time.Second)

	// Get container IPs
	fmt.Println("Getting container IPs...")
	for _, container := range containers {
		ip, err := getContainerIP(container)
		if err != nil {
			fmt.Printf("Warning: failed to get IP for %s: %v\n", container, err)
			continue
		}
		fmt.Printf("%s: %s\n", container, ip)
	}

	fmt.Println("Test environment ready!")
	return nil
}

func runStop() error {
	fmt.Println("Stopping spooky test environment...")

	testEnvDir := getTestEnvDir()
	if err := os.Chdir(testEnvDir); err != nil {
		return fmt.Errorf("failed to change to test environment directory: %w", err)
	}

	// Stop the containers using podman
	fmt.Println("Stopping containers with podman...")
	containers := []string{"spooky-server1", "spooky-server2", "spooky-server3"}

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

	// Remove the network
	fmt.Println("Removing spooky-test network...")
	if err := cmd.Run("podman", "network", "rm", "spooky-test"); err != nil {
		fmt.Println("Network spooky-test not found or already removed")
	}

	// Remove any remaining containers
	containers := []string{"spooky-server1", "spooky-server2", "spooky-server3"}
	for _, container := range containers {
		fmt.Printf("Removing container %s...\n", container)
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

	// Check network status
	fmt.Println("Network:")
	if err := cmd.Run("podman", "network", "ls", "--filter", "name=spooky-test"); err != nil {
		fmt.Println("  spooky-test: not found")
	} else {
		fmt.Println("  spooky-test: exists")
	}
	fmt.Println()

	// Check container status
	fmt.Println("Containers:")
	containers := []string{"spooky-server1", "spooky-server2", "spooky-server3"}
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

func startContainerDirectly(containerName string) error {
	// Define container configurations
	containerConfigs := map[string]struct {
		port     string
		user     string
		password string
		command  string
	}{
		"spooky-server1": {
			port:     "2221:22",
			user:     "root",
			password: "password",
			command:  "bash -c \"echo 'nameserver 8.8.8.8' > /etc/resolv.conf && echo 'nameserver 1.1.1.1' >> /etc/resolv.conf && apt-get update && apt-get install -y openssh-server sudo && mkdir -p /var/run/sshd && echo 'root:password' | chpasswd && sed -i 's/#PermitRootLogin prohibit-password/PermitRootLogin yes/' /etc/ssh/sshd_config && sed -i 's/#PasswordAuthentication yes/PasswordAuthentication yes/' /etc/ssh/sshd_config && /usr/sbin/sshd -D\"",
		},
		"spooky-server2": {
			port:     "2222:22",
			user:     "admin",
			password: "adminpass",
			command:  "bash -c \"echo 'nameserver 8.8.8.8' > /etc/resolv.conf && echo 'nameserver 1.1.1.1' >> /etc/resolv.conf && apt-get update && apt-get install -y openssh-server sudo && mkdir -p /var/run/sshd && useradd -m -s /bin/bash admin && echo 'admin:adminpass' | chpasswd && mkdir -p /home/admin/.ssh && chown -R admin:admin /home/admin/.ssh && chmod 700 /home/admin/.ssh && sed -i 's/#PubkeyAuthentication yes/PubkeyAuthentication yes/' /etc/ssh/sshd_config && sed -i 's/#PasswordAuthentication yes/PasswordAuthentication yes/' /etc/ssh/sshd_config && /usr/sbin/sshd -D\"",
		},
		"spooky-server3": {
			port:     "2223:22",
			user:     "user",
			password: "userpass",
			command:  "bash -c \"echo 'nameserver 8.8.8.8' > /etc/resolv.conf && echo 'nameserver 1.1.1.1' >> /etc/resolv.conf && apt-get update && apt-get install -y openssh-server sudo openssh-sftp-server && mkdir -p /var/run/sshd && useradd -m -s /bin/bash user && echo 'user:userpass' | chpasswd && mkdir -p /home/user/.ssh && chown -R user:user /home/user/.ssh && chmod 700 /home/user/.ssh && mkdir -p /home/user/sftp && chown -R user:user /home/user/sftp && sed -i 's/#PubkeyAuthentication yes/PubkeyAuthentication yes/' /etc/ssh/sshd_config && sed -i 's/#PasswordAuthentication yes/PasswordAuthentication yes/' /etc/ssh/sshd_config && sed -i 's/#Subsystem sftp \\/usr\\/lib\\/openssh\\/sftp-server/Subsystem sftp internal-sftp/' /etc/ssh/sshd_config && echo 'Match User user' >> /etc/ssh/sshd_config && echo '  ChrootDirectory /home/user/sftp' >> /etc/ssh/sshd_config && echo '  ForceCommand internal-sftp' >> /etc/ssh/sshd_config && echo '  AllowTcpForwarding no' >> /etc/ssh/sshd_config && echo '  X11Forwarding no' >> /etc/ssh/sshd_config && /usr/sbin/sshd -D\"",
		},
	}

	config, exists := containerConfigs[containerName]
	if !exists {
		return fmt.Errorf("unknown container: %s", containerName)
	}

	// Run the container with DNS configuration
	args := []string{
		"run", "-d",
		"--name", containerName,
		"--hostname", containerName,
		"--network", "spooky-test",
		"--ip", fmt.Sprintf("10.89.1.%s", strings.TrimPrefix(containerName, "spooky-server")),
		"-p", config.port,
		"-e", "DEBIAN_FRONTEND=noninteractive",
		"--dns", "8.8.8.8",
		"--dns", "8.8.4.4",
		"debian:12-slim",
		"bash", "-c", config.command,
	}

	return cmd.Run("podman", args...)
}

func getContainerIP(containerName string) (string, error) {
	output, err := cmd.Output("podman", "inspect", containerName, "--format", "{{.NetworkSettings.Networks.spooky-test.IPAddress}}")
	if err != nil {
		return "", fmt.Errorf("failed to inspect container %s: %w", containerName, err)
	}

	ip := strings.TrimSpace(string(output))
	if ip == "" {
		return "", fmt.Errorf("no IP address found for container %s", containerName)
	}

	return ip, nil
}
