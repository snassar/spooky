package tests

import (
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

// TestSFTPServer tests the SFTP server functionality
func TestSFTPServer(t *testing.T) {
	// Set up cleanup at the end of test
	t.Cleanup(func() {
		CleanupServers()
	})

	testsDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}

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
