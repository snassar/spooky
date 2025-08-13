package ssh

import (
	"os"
	"path/filepath"
	"testing"

	spookylogging "spooky/internal/logging"
	spookytypes "spooky/internal/types"
)

func TestHostKeyManager(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "spooky-ssh-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test logger
	logManager := spookylogging.NewLogManager()
	logger := logManager.GetLogger("test")

	// Create host key manager with test configuration
	knownHostsPath := filepath.Join(tempDir, "known_hosts")
	hostKeyManager := NewHostKeyManager(knownHostsPath, true, false, logger)

	// Test loading empty known hosts file (should not fail if file doesn't exist)
	err = hostKeyManager.LoadKnownHosts()
	if err != nil {
		t.Fatalf("Failed to load empty known hosts: %v", err)
	}

	// Test saving known hosts (this should create the file)
	err = hostKeyManager.SaveKnownHosts()
	if err != nil {
		t.Fatalf("Failed to save known hosts: %v", err)
	}

	// Test that the file was created after saving
	if _, err := os.Stat(knownHostsPath); os.IsNotExist(err) {
		t.Error("Known hosts file was not created after saving")
	}

	t.Logf("Host key manager test completed successfully")
}

func TestClientWithHostKeyVerification(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "spooky-ssh-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test configuration with relaxed host key checking for testing
	config := &spookytypes.ClientConfig{
		DefaultPort:        22,
		DefaultTimeout:     30,
		MaxConnections:     10,
		MaxRetryAttempts:   3,
		RetryDelay:         5,
		IdleTimeout:        300,
		KnownHostsPath:     filepath.Join(tempDir, "known_hosts"),
		StrictHostKeyCheck: false, // Allow new hosts for testing
		AllowInsecureHosts: true,  // Allow insecure hosts for testing
	}

	// Create a test logger
	logManager := spookylogging.NewLogManager()
	logger := logManager.GetLogger("test")

	// Create client
	client := NewClient(config, logger)

	// Test host key verification
	err = client.TestHostKeyVerification()
	if err != nil {
		t.Fatalf("Host key verification test failed: %v", err)
	}

	t.Logf("Simple client with host key verification test completed successfully")
}

func TestHostKeyVerificationStrictMode(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "spooky-ssh-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test configuration with strict host key checking
	config := &spookytypes.ClientConfig{
		DefaultPort:        22,
		DefaultTimeout:     30,
		MaxConnections:     10,
		MaxRetryAttempts:   3,
		RetryDelay:         5,
		IdleTimeout:        300,
		KnownHostsPath:     filepath.Join(tempDir, "known_hosts"),
		StrictHostKeyCheck: true,  // Strict mode
		AllowInsecureHosts: false, // Don't allow insecure hosts
	}

	// Create a test logger
	logManager := spookylogging.NewLogManager()
	logger := logManager.GetLogger("test")

	// Create client
	client := NewClient(config, logger)

	// Test host key verification - should fail in strict mode
	err = client.TestHostKeyVerification()
	if err == nil {
		t.Error("Host key verification should have failed in strict mode")
	} else {
		t.Logf("Host key verification correctly failed in strict mode: %v", err)
	}

	t.Logf("Host key verification strict mode test completed successfully")
}
