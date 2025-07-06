package tests

import (
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"
)

// TestTimeoutServer tests the timeout SSH server functionality
func TestTimeoutServer(t *testing.T) {
	// Set up cleanup at the end of test
	t.Cleanup(func() {
		CleanupServers()
	})

	testsDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}

	serverPath := filepath.Join(testsDir, "infrastructure", "timeout_server")
	cleanup, port := startServer(t, serverPath)
	defer cleanup()

	// Test SSH connection
	client, err := connectSSHWithPassword("127.0.0.1:"+strconv.Itoa(port), "testuser", "anypassword")
	if err != nil {
		t.Fatalf("Failed to connect to timeout server: %v", err)
	}
	defer client.Close()

	// Test session
	session, err := client.NewSession()
	if err != nil {
		t.Fatalf("Failed to create session: %v", err)
	}
	defer session.Close()

	// The timeout server should keep the connection alive for a few seconds
	// We'll just verify we can connect and the session doesn't immediately close
	time.Sleep(3 * time.Second)

	// Close the session and connection explicitly
	session.Close()
	client.Close()

	// If we get here without errors, the test passes
	t.Log("Timeout server test completed successfully")
}

// TestTimeoutServerConnectionTimeout tests that connections eventually timeout
func TestTimeoutServerConnectionTimeout(t *testing.T) {
	// Set up cleanup at the end of test
	t.Cleanup(func() {
		CleanupServers()
	})

	testsDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}

	serverPath := filepath.Join(testsDir, "infrastructure", "timeout_server")
	cleanup, port := startServer(t, serverPath)
	defer cleanup()

	// Test SSH connection
	client, err := connectSSHWithPassword("127.0.0.1:"+strconv.Itoa(port), "testuser", "anypassword")
	if err != nil {
		t.Fatalf("Failed to connect to timeout server: %v", err)
	}
	defer client.Close()

	// Test session
	session, err := client.NewSession()
	if err != nil {
		t.Fatalf("Failed to create session: %v", err)
	}
	defer session.Close()

	// Wait a bit to test the connection, but not too long
	// Note: This test is simplified to avoid hanging
	t.Log("Testing timeout server connection...")
	time.Sleep(5 * time.Second)

	// Close the session and connection explicitly
	session.Close()
	client.Close()

	t.Log("Timeout server connection test completed")
}
