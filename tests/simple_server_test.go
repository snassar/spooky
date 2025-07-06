package tests

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

// TestSimpleServer tests the simple SSH server functionality
func TestSimpleServer(t *testing.T) {
	// Set up cleanup at the end of test
	t.Cleanup(func() {
		CleanupServers()
	})

	// Get the absolute path to the tests directory
	testsDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}

	serverPath := filepath.Join(testsDir, "infrastructure", "simple_server")
	cleanup, port := startServer(t, serverPath)
	defer cleanup()

	// Test basic SSH connection
	client, err := connectSSH("127.0.0.1:"+strconv.Itoa(port), "testuser", nil)
	if err != nil {
		t.Fatalf("Failed to connect to simple server: %v", err)
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

	expected := "Hello testuser\n"
	if string(output) != expected {
		t.Errorf("Expected output %q, got %q", expected, string(output))
	}
}

// TestSimpleServerMultipleConnections tests multiple concurrent connections
func TestSimpleServerMultipleConnections(t *testing.T) {
	// Set up cleanup at the end of test
	t.Cleanup(func() {
		CleanupServers()
	})

	testsDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}

	serverPath := filepath.Join(testsDir, "infrastructure", "simple_server")
	cleanup, port := startServer(t, serverPath)
	defer cleanup()

	// Test multiple concurrent connections
	for i := 0; i < 3; i++ {
		t.Run(fmt.Sprintf("Connection%d", i), func(t *testing.T) {
			client, err := connectSSH("127.0.0.1:"+strconv.Itoa(port), fmt.Sprintf("user%d", i), nil)
			if err != nil {
				t.Fatalf("Failed to connect: %v", err)
			}
			defer client.Close()

			session, err := client.NewSession()
			if err != nil {
				t.Fatalf("Failed to create session: %v", err)
			}
			defer session.Close()

			output, err := session.Output("")
			if err != nil {
				t.Fatalf("Failed to execute command: %v", err)
			}

			expected := fmt.Sprintf("Hello user%d\n", i)
			if string(output) != expected {
				t.Errorf("Expected output %q, got %q", expected, string(output))
			}
		})
	}
}
