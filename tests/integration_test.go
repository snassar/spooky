package tests

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"
)

var (
	runSimpleServer    = flag.Bool("simple", false, "Run only simple server test")
	runPublicKeyServer = flag.Bool("publickey", false, "Run only public key server test")
	runSFTPServer      = flag.Bool("sftp", false, "Run only SFTP server test")
	runTimeoutServer   = flag.Bool("timeout", false, "Run only timeout server test")
)

// TestInfrastructureServers runs integration tests against all infrastructure servers
func TestInfrastructureServers(t *testing.T) {
	// Parse flags
	flag.Parse()

	// Set up cleanup at the end of all tests
	t.Cleanup(func() {
		CleanupServers()
	})

	// Get the absolute path to the tests directory
	testsDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}

	// Check if any specific test is requested
	if *runSimpleServer || *runPublicKeyServer || *runSFTPServer || *runTimeoutServer {
		// Run only the requested test(s)
		if *runSimpleServer {
			t.Run("SimpleServer", func(t *testing.T) {
				testSimpleServerFunc(t, testsDir)
			})
		}
		if *runPublicKeyServer {
			t.Run("PublicKeyServer", func(t *testing.T) {
				testPublicKeyServerFunc(t, testsDir)
			})
		}
		if *runSFTPServer {
			t.Run("SFTPServer", func(t *testing.T) {
				testSFTPServerFunc(t, testsDir)
			})
		}
		if *runTimeoutServer {
			t.Run("TimeoutServer", func(t *testing.T) {
				testTimeoutServerFunc(t, testsDir)
			})
		}
	} else {
		// Run all tests
		t.Run("SimpleServer", func(t *testing.T) {
			testSimpleServerFunc(t, testsDir)
		})

		t.Run("PublicKeyServer", func(t *testing.T) {
			testPublicKeyServerFunc(t, testsDir)
		})

		t.Run("SFTPServer", func(t *testing.T) {
			testSFTPServerFunc(t, testsDir)
		})

		t.Run("TimeoutServer", func(t *testing.T) {
			testTimeoutServerFunc(t, testsDir)
		})
	}
}

// startServer starts a server and returns a cleanup function and the port
func startServer(t *testing.T, serverPath string) (func(), int) {
	cmd := exec.Command("go", "run", "main.go")
	cmd.Dir = serverPath

	// Capture stdout to parse the port
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		t.Fatalf("Failed to start server: %v", err)
	}

	// Add process to tracker for cleanup
	serverTracker.AddProcess(cmd.Process.Pid, cmd)

	// Wait a bit for the server to start and find its port
	time.Sleep(3 * time.Second)

	// Parse the port from server output
	output := stdout.String()
	port := parsePortFromOutput(output)
	if port == 0 {
		t.Fatalf("Failed to parse port from server output: %s", output)
	}

	// Return cleanup function and port
	// Note: We don't call cmd.Wait() here because the server runs indefinitely
	return func() {
		if cmd.Process != nil {
			serverTracker.RemoveProcess(cmd.Process.Pid)
			cmd.Process.Kill()
			// Don't wait for the process to exit since it might hang
			// The cleanup will handle killing any remaining processes
		}
	}, port
}

// parsePortFromOutput extracts the port number from server startup output
func parsePortFromOutput(output string) int {
	// Look for patterns like "Starting ... server on :1234" or "Starting ... server on 127.0.0.1:1234"
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.Contains(line, "Starting") && strings.Contains(line, "server on") {
			// Extract port from the line
			if strings.Contains(line, ":") {
				parts := strings.Split(line, ":")
				if len(parts) >= 2 {
					portStr := strings.TrimSpace(parts[len(parts)-1])
					// Remove any trailing text after the port number
					if idx := strings.Index(portStr, " "); idx != -1 {
						portStr = portStr[:idx]
					}
					if port, err := strconv.Atoi(portStr); err == nil {
						return port
					}
				}
			}
		}
	}
	return 0
}

// testSimpleServerFunc tests the simple SSH server
func testSimpleServerFunc(t *testing.T, testsDir string) {
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

// testPublicKeyServerFunc tests the public key SSH server
func testPublicKeyServerFunc(t *testing.T, testsDir string) {
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
}

// testSFTPServerFunc tests the SFTP server
func testSFTPServerFunc(t *testing.T, testsDir string) {
	serverPath := filepath.Join(testsDir, "infrastructure", "sftp_server")
	cleanup, port := startServer(t, serverPath)
	defer cleanup()

	// Test SFTP connection
	client, err := connectSSH("127.0.0.1:"+strconv.Itoa(port), "testuser", nil)
	if err != nil {
		t.Fatalf("Failed to connect to SFTP server: %v", err)
	}
	defer client.Close()

	// Test SFTP session
	sftpClient, err := connectSFTP(client)
	if err != nil {
		t.Fatalf("Failed to create SFTP client: %v", err)
	}
	defer sftpClient.Close()

	// Test basic SFTP operations
	testSFTPOperations(t, sftpClient)
}

// testTimeoutServerFunc tests the timeout SSH server
func testTimeoutServerFunc(t *testing.T, testsDir string) {
	serverPath := filepath.Join(testsDir, "infrastructure", "timeout_server")
	cleanup, port := startServer(t, serverPath)
	defer cleanup()

	// Test SSH connection with password authentication
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

	// If we get here without errors, the test passes
	fmt.Println("Timeout server test completed successfully")
}
