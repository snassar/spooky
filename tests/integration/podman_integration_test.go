package integration

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

var (
	runPodmanTests = flag.Bool("podman", false, "Run Podman-based integration tests")
	testConfigFile = flag.String("config", "examples/test-environment/test-config.hcl", "Test configuration file")
)

// TestPodmanIntegration runs integration tests against the Podman-based test environment
func TestPodmanIntegration(t *testing.T) {
	// Parse flags
	flag.Parse()

	// Skip if not specifically requested
	if !*runPodmanTests {
		t.Skip("Skipping Podman integration tests. Use -podman flag to run.")
	}

	// Set up cleanup at the end of all tests
	t.Cleanup(func() {
		cleanupPodmanEnvironment(t)
	})

	// Get the absolute path to the project root
	projectRoot, err := getProjectRoot()
	if err != nil {
		t.Fatalf("Failed to get project root: %v", err)
	}

	// Start the Podman test environment
	t.Run("SetupEnvironment", func(t *testing.T) {
		setupPodmanEnvironment(t, projectRoot)
	})

	// Test basic authentication
	t.Run("TestAuthentication", func(t *testing.T) {
		testPodmanAuthentication(t, projectRoot)
	})

	// Test system information gathering
	t.Run("TestSystemInfo", func(t *testing.T) {
		testPodmanSystemInfo(t, projectRoot)
	})

	// Test file operations
	t.Run("TestFileOperations", func(t *testing.T) {
		testPodmanFileOperations(t, projectRoot)
	})

	// Test SFTP operations
	t.Run("TestSFTPOperations", func(t *testing.T) {
		testPodmanSFTPOperations(t, projectRoot)
	})

	// Test tag-based targeting
	t.Run("TestTagBasedTargeting", func(t *testing.T) {
		testPodmanTagBasedTargeting(t, projectRoot)
	})

	// Test concurrent operations
	t.Run("TestConcurrentOperations", func(t *testing.T) {
		testPodmanConcurrentOperations(t, projectRoot)
	})

	// Test error handling
	t.Run("TestErrorHandling", func(t *testing.T) {
		testPodmanErrorHandling(t, projectRoot)
	})

	// Test network connectivity
	t.Run("TestNetworkConnectivity", func(t *testing.T) {
		testPodmanNetworkConnectivity(t, projectRoot)
	})
}

// TestPodmanEnvironmentSetup tests basic environment setup and teardown
func TestPodmanEnvironmentSetup(t *testing.T) {
	// Parse flags
	flag.Parse()

	// Skip if not specifically requested
	if !*runPodmanTests {
		t.Skip("Skipping Podman integration tests. Use -podman flag to run.")
	}

	// Get the absolute path to the project root
	projectRoot, err := getProjectRoot()
	if err != nil {
		t.Fatalf("Failed to get project root: %v", err)
	}

	// Test environment setup
	t.Run("EnvironmentSetup", func(t *testing.T) {
		t.Log("Testing basic environment setup...")

		// Change to project root
		if err := os.Chdir(projectRoot); err != nil {
			t.Fatalf("Failed to change to project root: %v", err)
		}

		// Test SSH connectivity to verify container is ready
		t.Log("Testing SSH connectivity...")
		cmd := exec.Command("ssh", "-o", "StrictHostKeyChecking=no", "-o", "ConnectTimeout=5",
			"-p", "2221", "testuser@localhost", "echo 'SSH connection successful'")
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("SSH connectivity test failed: %v\nOutput: %s", err, string(output))
		}
		t.Logf("SSH connectivity verified: %s", string(output))

		// Test configuration validation
		t.Log("Testing configuration validation...")
		//nolint:gosec // testConfigFile is controlled by test runner, not user input
		cmd = exec.Command("go", "run", "main.go", "validate", *testConfigFile)
		output, err = cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Configuration validation failed: %v\nOutput: %s", err, string(output))
		}
		t.Logf("Configuration validation successful: %s", string(output))

		// Test listing servers and actions
		t.Log("Testing list command...")
		//nolint:gosec // testConfigFile is controlled by test runner, not user input
		cmd = exec.Command("go", "run", "main.go", "list", *testConfigFile)
		output, err = cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("List command failed: %v\nOutput: %s", err, string(output))
		}
		t.Logf("List command successful: %s", string(output))

		t.Log("Environment setup test completed successfully")
	})
}

// setupPodmanEnvironment starts the Podman test environment
func setupPodmanEnvironment(t *testing.T, projectRoot string) {
	t.Log("Setting up Podman test environment...")

	// Change to project root
	if err := os.Chdir(projectRoot); err != nil {
		t.Fatalf("Failed to change to project root: %v", err)
	}

	// The workflow already handles container setup, so we just verify connectivity
	t.Log("Verifying SSH connectivity...")

	// Test SSH connection to verify container is ready
	cmd := exec.Command("ssh", "-o", "StrictHostKeyChecking=no", "-o", "ConnectTimeout=5",
		"-p", "2221", "testuser@localhost", "echo 'SSH connection successful'")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("SSH connectivity test failed: %v\nOutput: %s", err, string(output))
	}

	t.Logf("SSH connectivity verified: %s", string(output))
}

// cleanupPodmanEnvironment cleans up the Podman test environment
func cleanupPodmanEnvironment(t *testing.T) {
	t.Log("Cleaning up Podman test environment...")

	// The workflow handles cleanup, so we just log the cleanup
	t.Log("Cleanup handled by workflow")
}

// testPodmanAuthentication tests basic authentication across all servers
func testPodmanAuthentication(t *testing.T, projectRoot string) {
	t.Log("Testing authentication across all servers...")

	if err := os.Chdir(projectRoot); err != nil {
		t.Fatalf("Failed to change to project root: %v", err)
	}

	// Test configuration validation first
	t.Log("Validating configuration...")
	//nolint:gosec // testConfigFile is controlled by test runner, not user input
	cmd := exec.Command("go", "run", "main.go", "validate", *testConfigFile)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Configuration validation failed: %v\nOutput: %s", err, string(output))
	}
	t.Logf("Configuration validation output:\n%s", string(output))

	// Test listing servers and actions
	t.Log("Listing servers and actions...")
	//nolint:gosec // testConfigFile is controlled by test runner, not user input
	cmd = exec.Command("go", "run", "main.go", "list", *testConfigFile)
	output, err = cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("List command failed: %v\nOutput: %s", err, string(output))
	}
	t.Logf("List command output:\n%s", string(output))

	// Execute all actions in the configuration
	t.Log("Executing all actions...")
	//nolint:gosec // testConfigFile is controlled by test runner, not user input
	cmd = exec.Command("go", "run", "main.go", "execute", *testConfigFile)
	output, err = cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Action execution failed: %v\nOutput: %s", err, string(output))
	}

	t.Logf("Action execution output:\n%s", string(output))
}

// testPodmanSystemInfo tests system information gathering
func testPodmanSystemInfo(t *testing.T, projectRoot string) {
	t.Log("Testing system information gathering...")

	if err := os.Chdir(projectRoot); err != nil {
		t.Fatalf("Failed to change to project root: %v", err)
	}

	// Execute the configuration (which includes system info commands)
	//nolint:gosec // testConfigFile is controlled by test runner, not user input
	cmd := exec.Command("go", "run", "main.go", "execute", *testConfigFile)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("System info test failed: %v\nOutput: %s", err, string(output))
	}

	t.Logf("System info test output:\n%s", string(output))
}

// testPodmanFileOperations tests basic file operations
func testPodmanFileOperations(t *testing.T, projectRoot string) {
	t.Log("Testing file operations...")

	if err := os.Chdir(projectRoot); err != nil {
		t.Fatalf("Failed to change to project root: %v", err)
	}

	// Execute the configuration (which includes file operations)
	//nolint:gosec // testConfigFile is controlled by test runner, not user input
	cmd := exec.Command("go", "run", "main.go", "execute", *testConfigFile)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("File operations test failed: %v\nOutput: %s", err, string(output))
	}

	t.Logf("File operations test output:\n%s", string(output))
}

// testPodmanSFTPOperations tests SFTP-specific operations
func testPodmanSFTPOperations(t *testing.T, projectRoot string) {
	t.Log("Testing SFTP operations...")

	if err := os.Chdir(projectRoot); err != nil {
		t.Fatalf("Failed to change to project root: %v", err)
	}

	// Execute the configuration (which includes SFTP operations)
	//nolint:gosec // testConfigFile is controlled by test runner, not user input
	cmd := exec.Command("go", "run", "main.go", "execute", *testConfigFile)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("SFTP operations test failed: %v\nOutput: %s", err, string(output))
	}

	t.Logf("SFTP operations test output:\n%s", string(output))
}

// testPodmanTagBasedTargeting tests tag-based server targeting
func testPodmanTagBasedTargeting(t *testing.T, projectRoot string) {
	t.Log("Testing tag-based targeting...")

	if err := os.Chdir(projectRoot); err != nil {
		t.Fatalf("Failed to change to project root: %v", err)
	}

	// Execute the configuration (which includes tag-based targeting)
	//nolint:gosec // testConfigFile is controlled by test runner, not user input
	cmd := exec.Command("go", "run", "main.go", "execute", *testConfigFile)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Tag-based targeting test failed: %v\nOutput: %s", err, string(output))
	}

	t.Logf("Tag-based targeting test output:\n%s", string(output))
}

// testPodmanConcurrentOperations tests concurrent operations
func testPodmanConcurrentOperations(t *testing.T, projectRoot string) {
	t.Log("Testing concurrent operations...")

	if err := os.Chdir(projectRoot); err != nil {
		t.Fatalf("Failed to change to project root: %v", err)
	}

	// Execute the configuration with parallel flag
	//nolint:gosec // testConfigFile is controlled by test runner, not user input
	cmd := exec.Command("go", "run", "main.go", "execute", *testConfigFile, "--parallel")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Concurrent operations test failed: %v\nOutput: %s", err, string(output))
	}

	t.Logf("Concurrent operations test output:\n%s", string(output))
}

// testPodmanErrorHandling tests error handling scenarios
func testPodmanErrorHandling(t *testing.T, projectRoot string) {
	t.Log("Testing error handling...")

	if err := os.Chdir(projectRoot); err != nil {
		t.Fatalf("Failed to change to project root: %v", err)
	}

	// Test with invalid configuration file
	t.Run("InvalidConfig", func(t *testing.T) {
		cmd := exec.Command("go", "run", "main.go", "validate", "nonexistent-config.hcl")
		output, err := cmd.CombinedOutput()
		// This should fail, which is expected
		if err == nil {
			t.Error("Expected validation to fail for nonexistent config file")
		}
		t.Logf("Invalid config test output:\n%s", string(output))
	})

	// Test with valid configuration
	t.Run("ValidConfig", func(t *testing.T) {
		//nolint:gosec // testConfigFile is controlled by test runner, not user input
		cmd := exec.Command("go", "run", "main.go", "validate", *testConfigFile)
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Valid config validation failed: %v\nOutput: %s", err, string(output))
		}
		t.Logf("Valid config test output:\n%s", string(output))
	})
}

// testPodmanNetworkConnectivity tests network connectivity between servers
func testPodmanNetworkConnectivity(t *testing.T, projectRoot string) {
	t.Log("Testing network connectivity...")

	if err := os.Chdir(projectRoot); err != nil {
		t.Fatalf("Failed to change to project root: %v", err)
	}

	// Execute the configuration (which includes network connectivity tests)
	//nolint:gosec // testConfigFile is controlled by test runner, not user input
	cmd := exec.Command("go", "run", "main.go", "execute", *testConfigFile)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Network connectivity test failed: %v\nOutput: %s", err, string(output))
	}

	t.Logf("Network connectivity test output:\n%s", string(output))
}

// getProjectRoot returns the absolute path to the project root
func getProjectRoot() (string, error) {
	// Start from the current directory and look for go.mod
	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// Walk up the directory tree to find go.mod
	for {
		if _, err := os.Stat(filepath.Join(currentDir, "go.mod")); err == nil {
			return currentDir, nil
		}

		parent := filepath.Dir(currentDir)
		if parent == currentDir {
			return "", fmt.Errorf("could not find go.mod file")
		}
		currentDir = parent
	}
}
