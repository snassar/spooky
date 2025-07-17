package integration

import (
	"flag"
	"os/exec"
	"strings"
	"testing"
	"time"
)

var runPodmanBasicTests = flag.Bool("podman-basic", false, "Run basic Podman environment tests")

// TestPodmanBasicEnvironment tests basic Podman environment functionality
func TestPodmanBasicEnvironment(t *testing.T) {
	// Parse flags
	flag.Parse()

	// Skip if not specifically requested
	if !*runPodmanBasicTests {
		t.Skip("Skipping basic Podman tests. Use -podman-basic flag to run.")
	}

	// Test basic environment operations
	t.Run("BasicEnvironment", func(t *testing.T) {

		// Test 1: Preflight check
		t.Log("Testing preflight check...")
		preflightCmd := exec.Command("go", "run", "../../tools/spooky-test-env/main.go", "preflight")
		preflightOutput, err := preflightCmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Preflight check failed: %v\nOutput: %s", err, string(preflightOutput))
		}
		t.Logf("Preflight check passed:\n%s", string(preflightOutput))

		// Test 2: Start environment
		t.Log("Testing environment start...")
		startCmd := exec.Command("go", "run", "../../tools/spooky-test-env/main.go", "start")
		startOutput, err := startCmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Environment start failed: %v\nOutput: %s", err, string(startOutput))
		}
		t.Logf("Environment started:\n%s", string(startOutput))

		// Test 3: Wait for containers to be ready
		t.Log("Waiting for containers to be ready...")
		time.Sleep(10 * time.Second)

		// Test 4: Check container status
		t.Log("Checking container status...")
		psCmd := exec.Command("podman", "ps", "--filter", "name=spooky-server", "--format", "{{.Names}} {{.Status}}")
		psOutput, err := psCmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Failed to check container status: %v", err)
		}

		containerInfo := strings.TrimSpace(string(psOutput))
		if containerInfo == "" {
			t.Fatal("No spooky-server containers found")
		}

		t.Logf("Container status:\n%s", containerInfo)

		// Test 5: Check if containers are running
		lines := strings.Split(containerInfo, "\n")
		for _, line := range lines {
			if !strings.Contains(line, "Up") {
				t.Errorf("Container not running: %s", line)
			}
		}

		// Test 6: Check environment status
		t.Log("Testing environment status...")
		statusCmd := exec.Command("go", "run", "../../tools/spooky-test-env/main.go", "status")
		statusOutput, err := statusCmd.CombinedOutput()
		if err != nil {
			t.Logf("Status check failed (non-critical): %v\nOutput: %s", err, string(statusOutput))
		} else {
			t.Logf("Environment status:\n%s", string(statusOutput))
		}

		// Test 7: Cleanup
		t.Log("Testing environment cleanup...")
		cleanupCmd := exec.Command("go", "run", "../../tools/spooky-test-env/main.go", "cleanup")
		cleanupOutput, err := cleanupCmd.CombinedOutput()
		if err != nil {
			t.Logf("Cleanup failed (non-critical): %v\nOutput: %s", err, string(cleanupOutput))
		} else {
			t.Logf("Environment cleaned up:\n%s", string(cleanupOutput))
		}

		// Test 8: Verify cleanup
		t.Log("Verifying cleanup...")
		psAfterCmd := exec.Command("podman", "ps", "--filter", "name=spooky-server", "--format", "{{.Names}}")
		psAfterOutput, err := psAfterCmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Failed to check container status after cleanup: %v", err)
		}

		remainingContainers := strings.TrimSpace(string(psAfterOutput))
		if remainingContainers != "" {
			t.Errorf("Containers still running after cleanup: %s", remainingContainers)
		} else {
			t.Log("All containers successfully cleaned up")
		}
	})
}
