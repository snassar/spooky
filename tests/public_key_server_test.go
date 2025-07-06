package tests

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

// TestPublicKeyServer tests the public key SSH server functionality
func TestPublicKeyServer(t *testing.T) {
	// Set up cleanup at the end of test
	t.Cleanup(func() {
		CleanupServers()
	})

	testsDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}

	serverPath := filepath.Join(testsDir, "infrastructure", "public_key_server")
	cleanup, port := startServer(t, serverPath)
	defer cleanup()

	// Generate test key pair using the main package's key generation
	privateKey, _, err := generateTestKeyPairFromMain()
	if err != nil {
		t.Fatalf("Failed to generate test key pair: %v", err)
	}

	// Test SSH connection with public key
	client, err := connectSSHWithKey("127.0.0.1:"+strconv.Itoa(port), "testuser", privateKey)
	if err != nil {
		t.Fatalf("Failed to connect to public key server: %v", err)
	}
	defer client.Close()

	// Test session
	session, err := client.NewSession()
	if err != nil {
		t.Fatalf("Failed to create session: %v", err)
	}
	defer session.Close()

	output, err := session.Output("")
	if err != nil {
		t.Fatalf("Failed to execute command: %v", err)
	}

	// Check that the output contains the public key
	if len(output) == 0 {
		t.Error("Expected non-empty output containing public key")
	}

	// Verify the output format
	expectedPrefix := "public key used by testuser:"
	if !strings.HasPrefix(string(output), expectedPrefix) {
		t.Errorf("Expected output to start with %q, got %q", expectedPrefix, string(output))
	}
}
