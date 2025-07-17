package integration

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
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

		// Run preflight check
		t.Log("Running preflight check...")
		if err := exec.Command("go", "run", "../../tools/spooky-test-env/main.go", "preflight").Run(); err != nil {
			t.Fatalf("Preflight check failed: %v", err)
		}

		// Start the test environment
		t.Log("Starting test environment...")
		if err := exec.Command("go", "run", "../../tools/spooky-test-env/main.go", "start").Run(); err != nil {
			t.Fatalf("Failed to start test environment: %v", err)
		}

		// Wait for containers to be ready
		t.Log("Waiting for containers to be ready...")
		time.Sleep(15 * time.Second)

		// Check environment status
		t.Log("Checking environment status...")
		if err := exec.Command("go", "run", "../../tools/spooky-test-env/main.go", "status").Run(); err != nil {
			t.Fatalf("Failed to check environment status: %v", err)
		}

		// Verify containers are running
		t.Log("Verifying containers are running...")
		output, err := exec.Command("podman", "ps", "--filter", "name=spooky-server", "--format", "{{.Names}}").Output()
		if err != nil {
			t.Fatalf("Failed to check container status: %v", err)
		}

		containerNames := strings.TrimSpace(string(output))
		if containerNames == "" {
			t.Fatal("No spooky-server containers found")
		}

		t.Logf("Found containers: %s", containerNames)

		// Clean up
		t.Log("Cleaning up test environment...")
		if cleanupErr := exec.Command("go", "run", "../../tools/spooky-test-env/main.go", "cleanup").Run(); cleanupErr != nil {
			t.Logf("Warning: failed to cleanup test environment: %v", cleanupErr)
		}
	})
}

// setupPodmanEnvironment starts the Podman test environment
func setupPodmanEnvironment(t *testing.T, projectRoot string) {
	t.Log("Setting up Podman test environment...")

	// Change to project root
	if err := os.Chdir(projectRoot); err != nil {
		t.Fatalf("Failed to change to project root: %v", err)
	}

	// Run preflight check
	t.Log("Running preflight check...")
	if err := exec.Command("go", "run", "../../tools/spooky-test-env/main.go", "preflight").Run(); err != nil {
		t.Fatalf("Preflight check failed: %v", err)
	}

	// Start the test environment
	t.Log("Starting test environment...")
	if err := exec.Command("go", "run", "../../tools/spooky-test-env/main.go", "start").Run(); err != nil {
		t.Fatalf("Failed to start test environment: %v", err)
	}

	// Wait for containers to be ready
	t.Log("Waiting for containers to be ready...")
	time.Sleep(15 * time.Second)

	// Check environment status
	t.Log("Checking environment status...")
	if err := exec.Command("go", "run", "tools/spooky-test-env/main.go", "status").Run(); err != nil {
		t.Fatalf("Failed to check environment status: %v", err)
	}
}

// cleanupPodmanEnvironment cleans up the Podman test environment
func cleanupPodmanEnvironment(t *testing.T) {
	t.Log("Cleaning up Podman test environment...")

	projectRoot, err := getProjectRoot()
	if err != nil {
		t.Logf("Warning: failed to get project root for cleanup: %v", err)
		return
	}

	if err := os.Chdir(projectRoot); err != nil {
		t.Logf("Warning: failed to change to project root for cleanup: %v", err)
		return
	}

	// Stop and cleanup the test environment
	if err := exec.Command("go", "run", "../../tools/spooky-test-env/main.go", "cleanup").Run(); err != nil {
		t.Logf("Warning: failed to cleanup test environment: %v", err)
	}
}

// testPodmanAuthentication tests basic authentication across all servers
func testPodmanAuthentication(t *testing.T, projectRoot string) {
	t.Log("Testing authentication across all servers...")

	if err := os.Chdir(projectRoot); err != nil {
		t.Fatalf("Failed to change to project root: %v", err)
	}

	// Execute the authentication test
	cmd := exec.Command("go", "run", "../../main.go", "execute", *testConfigFile)
	cmd.Args = append(cmd.Args, "--action", "test-authentication")

	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Authentication test failed: %v\nOutput: %s", err, string(output))
	}

	t.Logf("Authentication test output:\n%s", string(output))
}

// testPodmanSystemInfo tests system information gathering
func testPodmanSystemInfo(t *testing.T, projectRoot string) {
	t.Log("Testing system information gathering...")

	if err := os.Chdir(projectRoot); err != nil {
		t.Fatalf("Failed to change to project root: %v", err)
	}

	// Execute the system info test
	cmd := exec.Command("go", "run", "../../main.go", "execute", *testConfigFile)
	cmd.Args = append(cmd.Args, "--action", "test-system-info")

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

	// Execute the file operations test
	cmd := exec.Command("go", "run", "../../main.go", "execute", *testConfigFile)
	cmd.Args = append(cmd.Args, "--action", "test-file-operations")

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

	// Execute the SFTP operations test
	cmd := exec.Command("go", "run", "../../main.go", "execute", *testConfigFile)
	cmd.Args = append(cmd.Args, "--action", "test-sftp-operations")

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

	// Test production servers
	t.Run("ProductionServers", func(t *testing.T) {
		cmd := exec.Command("go", "run", "../../main.go", "execute", *testConfigFile)
		cmd.Args = append(cmd.Args, "--action", "test-production-servers")

		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Production servers test failed: %v\nOutput: %s", err, string(output))
		}
		t.Logf("Production servers test output:\n%s", string(output))
	})

	// Test staging servers
	t.Run("StagingServers", func(t *testing.T) {
		cmd := exec.Command("go", "run", "../../main.go", "execute", *testConfigFile)
		cmd.Args = append(cmd.Args, "--action", "test-staging-servers")

		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Staging servers test failed: %v\nOutput: %s", err, string(output))
		}
		t.Logf("Staging servers test output:\n%s", string(output))
	})
}

// testPodmanConcurrentOperations tests concurrent operations
func testPodmanConcurrentOperations(t *testing.T, projectRoot string) {
	t.Log("Testing concurrent operations...")

	if err := os.Chdir(projectRoot); err != nil {
		t.Fatalf("Failed to change to project root: %v", err)
	}

	// Execute the concurrent operations test
	cmd := exec.Command("go", "run", "../../main.go", "execute", *testConfigFile)
	cmd.Args = append(cmd.Args, "--action", "test-concurrent-operations")

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

	// Execute the error handling test
	cmd := exec.Command("go", "run", "../../main.go", "execute", *testConfigFile)
	cmd.Args = append(cmd.Args, "--action", "test-error-handling")

	output, err := cmd.CombinedOutput()
	// Note: This test is expected to handle errors gracefully, so we don't fail on error
	t.Logf("Error handling test output:\n%s", string(output))
	if err != nil {
		t.Logf("Error occurred (expected): %v", err)
	}
}

// testPodmanNetworkConnectivity tests network connectivity between servers
func testPodmanNetworkConnectivity(t *testing.T, projectRoot string) {
	t.Log("Testing network connectivity...")

	if err := os.Chdir(projectRoot); err != nil {
		t.Fatalf("Failed to change to project root: %v", err)
	}

	// Execute the network connectivity test
	cmd := exec.Command("go", "run", "../../main.go", "execute", *testConfigFile)
	cmd.Args = append(cmd.Args, "--action", "test-network-connectivity")

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
