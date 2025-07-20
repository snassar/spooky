package integration

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"spooky/internal/config"
	"spooky/internal/ssh"
)

func TestPodmanCIIntegration(t *testing.T) {
	// Skip if not running in CI environment
	if os.Getenv("SPOOKY_TEST_SSH_HOST") == "" {
		t.Skip("Skipping Podman CI integration test - not in CI environment")
	}

	// Get environment variables
	sshHost := os.Getenv("SPOOKY_TEST_SSH_HOST")
	sshPortStr := os.Getenv("SPOOKY_TEST_SSH_PORT")
	sshUser := os.Getenv("SPOOKY_TEST_SSH_USER")
	sshKey := os.Getenv("SPOOKY_TEST_SSH_KEY")

	if sshHost == "" || sshPortStr == "" || sshUser == "" || sshKey == "" {
		t.Fatal("Missing required environment variables for SSH connection")
	}

	// Expand ~ in SSH key path if present
	if strings.HasPrefix(sshKey, "~") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			t.Fatalf("Failed to get home directory: %v", err)
		}
		sshKey = strings.Replace(sshKey, "~", homeDir, 1)
	}

	// Parse port number
	sshPort, err := strconv.Atoi(sshPortStr)
	if err != nil {
		t.Fatalf("Invalid SSH port number: %s", sshPortStr)
	}

	t.Logf("Testing SSH connection to %s@%s:%d using key %s", sshUser, sshHost, sshPort, sshKey)

	// Test 1: Basic SSH connection
	t.Run("SSH Connection", func(t *testing.T) {
		// Debug: Check if SSH key file exists
		if _, err := os.Stat(sshKey); os.IsNotExist(err) {
			t.Fatalf("SSH key file does not exist: %s", sshKey)
		}
		t.Logf("SSH key file exists: %s", sshKey)

		// Check SSH key permissions
		if info, err := os.Stat(sshKey); err == nil {
			t.Logf("SSH key permissions: %o", info.Mode().Perm())
		}

		// Create a simple machine config for testing
		machine := config.Machine{
			Name:     "test-machine",
			Host:     sshHost,
			Port:     sshPort, // Use port from environment variable
			User:     sshUser,
			KeyFile:  sshKey,
			Password: "", // Use key authentication
		}

		t.Logf("Attempting SSH connection to %s@%s:%d with key %s", machine.User, machine.Host, machine.Port, machine.KeyFile)

		// Test SSH client creation with longer timeout for CI
		client, err := ssh.NewSSHClient(&machine, 60) // Increased timeout for CI
		if err != nil {
			t.Fatalf("Failed to create SSH client: %v", err)
		}
		defer client.Close()

		// Test command execution
		output, err := client.ExecuteCommand("echo 'SSH connection test successful'")
		if err != nil {
			t.Fatalf("Failed to execute command: %v", err)
		}

		expectedOutput := "SSH connection test successful"
		if output != expectedOutput {
			t.Errorf("Expected output '%s', got '%s'", expectedOutput, output)
		}

		t.Logf("SSH connection test passed: %s", output)
	})

	// Test 2: Configuration file execution
	t.Run("Configuration Execution", func(t *testing.T) {
		// Create a temporary config file
		configContent := fmt.Sprintf(`
server "test-server" {
  host = "%s"
  port = %d
  user = "%s"
  key_file = "%s"
}

action "test-action" {
  description = "Test action in CI environment"
  command = "echo 'Configuration test successful' && whoami && pwd"
  servers = ["test-server"]
  parallel = false
  timeout = 30
}
`, sshHost, sshPort, sshUser, sshKey)

		// Write config to temporary file
		tmpDir := t.TempDir()
		configFile := filepath.Join(tmpDir, "test-config.hcl")
		if err := os.WriteFile(configFile, []byte(configContent), 0o600); err != nil {
			t.Fatalf("Failed to write config file: %v", err)
		}

		// Parse configuration
		cfg, err := config.ParseConfig(configFile)
		if err != nil {
			t.Fatalf("Failed to parse config: %v", err)
		}

		// Execute configuration
		err = ssh.ExecuteConfig(cfg)
		if err != nil {
			t.Fatalf("Failed to execute config: %v", err)
		}

		t.Log("Configuration execution test passed")
	})

	// Test 3: Script execution
	t.Run("Script Execution", func(t *testing.T) {
		// Create a test script
		scriptContent := `#!/bin/bash
echo "Script execution test started"
echo "Current user: $(whoami)"
echo "Current directory: $(pwd)"
echo "Script execution test completed"
`

		// Write script to temporary file
		tmpDir := t.TempDir()
		scriptFile := filepath.Join(tmpDir, "test-script.sh")
		if err := os.WriteFile(scriptFile, []byte(scriptContent), 0o600); err != nil {
			t.Fatalf("Failed to write script file: %v", err)
		}

		// Create machine config
		machine := config.Machine{
			Name:     "test-machine",
			Host:     sshHost,
			Port:     sshPort,
			User:     sshUser,
			KeyFile:  sshKey,
			Password: "",
		}

		// Create SSH client
		client, err := ssh.NewSSHClient(&machine, 30)
		if err != nil {
			t.Fatalf("Failed to create SSH client: %v", err)
		}
		defer client.Close()

		// Execute script
		output, err := client.ExecuteScript(scriptFile)
		if err != nil {
			t.Fatalf("Failed to execute script: %v", err)
		}

		// Verify output contains expected content
		if output == "" {
			t.Error("Script execution produced no output")
		}

		t.Logf("Script execution test passed. Output: %s", output)
	})

	// Test 4: Parallel execution
	t.Run("Parallel Execution", func(t *testing.T) {
		// Create a config with parallel execution
		configContent := fmt.Sprintf(`
server "test-server" {
  host = "%s"
  port = %d
  user = "%s"
  key_file = "%s"
}

action "parallel-test" {
  description = "Test parallel execution"
  command = "sleep 1 && echo 'Parallel execution test'"
  servers = ["test-server"]
  parallel = true
  timeout = 30
}
`, sshHost, sshPort, sshUser, sshKey)

		// Write config to temporary file
		tmpDir := t.TempDir()
		configFile := filepath.Join(tmpDir, "parallel-config.hcl")
		if err := os.WriteFile(configFile, []byte(configContent), 0o600); err != nil {
			t.Fatalf("Failed to write config file: %v", err)
		}

		// Parse and execute configuration
		cfg, err := config.ParseConfig(configFile)
		if err != nil {
			t.Fatalf("Failed to parse config: %v", err)
		}

		startTime := time.Now()
		err = ssh.ExecuteConfig(cfg)
		duration := time.Since(startTime)

		if err != nil {
			t.Fatalf("Failed to execute parallel config: %v", err)
		}

		// Parallel execution should be faster than sequential
		if duration > 5*time.Second {
			t.Logf("Parallel execution took %v (expected to be fast)", duration)
		}

		t.Log("Parallel execution test passed")
	})

	// Test 5: Error handling
	t.Run("Error Handling", func(t *testing.T) {
		// Create machine config
		machine := config.Machine{
			Name:     "test-machine",
			Host:     sshHost,
			Port:     sshPort,
			User:     sshUser,
			KeyFile:  sshKey,
			Password: "",
		}

		// Create SSH client
		client, err := ssh.NewSSHClient(&machine, 30)
		if err != nil {
			t.Fatalf("Failed to create SSH client: %v", err)
		}
		defer client.Close()

		// Test command that should fail
		_, err = client.ExecuteCommand("nonexistentcommand")
		if err == nil {
			t.Error("Expected error for nonexistent command, but got none")
		} else {
			t.Logf("Error handling test passed: %v", err)
		}
	})

	// Test 6: CLI command execution
	t.Run("CLI Commands", func(t *testing.T) {
		// Test validate command
		tmpDir := t.TempDir()
		configFile := filepath.Join(tmpDir, "cli-test-config.hcl")
		configContent := fmt.Sprintf(`
server "test-server" {
  host = "%s"
  port = %d
  user = "%s"
  key_file = "%s"
}

action "test-action" {
  description = "Test action"
  command = "echo 'test'"
  servers = ["test-server"]
}
`, sshHost, sshPort, sshUser, sshKey)

		if err := os.WriteFile(configFile, []byte(configContent), 0o600); err != nil {
			t.Fatalf("Failed to write config file: %v", err)
		}

		// Check if binary exists and is executable
		binaryPath := "./build/spooky"
		if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
			t.Fatalf("Spooky binary not found: %s", binaryPath)
		}

		// Test validate command using built binary
		cmd := exec.Command(binaryPath, "validate", configFile)
		cmd.Dir = "." // Ensure we're in the right directory
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Validate command failed: %v\nOutput: %s", err, output)
		}

		t.Logf("Validate command passed: %s", output)

		// Test list command using built binary
		cmd = exec.Command(binaryPath, "list", configFile)
		cmd.Dir = "." // Ensure we're in the right directory
		output, err = cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("List command failed: %v\nOutput: %s", err, output)
		}

		t.Logf("List command passed: %s", output)
	})
}

// TestPodmanContainerSetup verifies that the Podman container is properly set up
func TestPodmanContainerSetup(t *testing.T) {
	// Skip if not in CI environment
	if os.Getenv("SPOOKY_TEST_SSH_HOST") == "" {
		t.Skip("Skipping Podman container setup test - not in CI environment")
	}

	sshHost := os.Getenv("SPOOKY_TEST_SSH_HOST")
	sshPort := os.Getenv("SPOOKY_TEST_SSH_PORT")
	sshUser := os.Getenv("SPOOKY_TEST_SSH_USER")

	t.Run("Container Connectivity", func(t *testing.T) {
		// Test basic connectivity to the container
		cmd := exec.Command("ssh", "-o", "StrictHostKeyChecking=no", "-o", "ConnectTimeout=10",
			"-p", sshPort, fmt.Sprintf("%s@%s", sshUser, sshHost), "echo 'connectivity test'")

		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Failed to connect to container: %v\nOutput: %s", err, output)
		}

		expectedOutput := "connectivity test\n"
		if string(output) != expectedOutput {
			t.Errorf("Expected output '%s', got '%s'", expectedOutput, string(output))
		}

		t.Log("Container connectivity test passed")
	})

	t.Run("Container Environment", func(t *testing.T) {
		// Test container environment
		cmd := exec.Command("ssh", "-o", "StrictHostKeyChecking=no", "-o", "ConnectTimeout=10",
			"-p", sshPort, fmt.Sprintf("%s@%s", sshUser, sshHost), "uname -a && whoami && pwd")

		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Failed to check container environment: %v\nOutput: %s", err, output)
		}

		t.Logf("Container environment: %s", output)
	})
}
