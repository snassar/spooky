package ssh

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"spooky/internal/config"
)

func TestNewSSHClient_NoAuth(t *testing.T) {
	server := &config.Server{
		Name: "test",
		Host: "localhost",
		User: "testuser",
		Port: 22,
	}

	_, err := NewSSHClient(server, 30)
	if err == nil {
		t.Error("expected error when no authentication method is provided")
	}
	if !strings.Contains(err.Error(), "no authentication method available") {
		t.Errorf("expected authentication error, got: %v", err)
	}
}

func TestNewSSHClient_PasswordAuth(t *testing.T) {
	server := &config.Server{
		Name:     "test",
		Host:     "localhost",
		User:     "testuser",
		Port:     22,
		Password: "testpass",
	}

	_, err := NewSSHClient(server, 30)
	// This will fail to connect, but should not fail due to authentication setup
	if err != nil && !strings.Contains(err.Error(), "connection refused") && !strings.Contains(err.Error(), "no route to host") {
		t.Errorf("expected connection error, got: %v", err)
	}
}

func TestNewSSHClient_KeyFileAuth(t *testing.T) {
	// Create a temporary key file
	tempDir := t.TempDir()
	keyFile := filepath.Join(tempDir, "test_key")

	// Write a dummy key file (this will fail to parse, but tests the file reading logic)
	err := os.WriteFile(keyFile, []byte("invalid key content"), 0600)
	if err != nil {
		t.Fatalf("failed to create test key file: %v", err)
	}

	server := &config.Server{
		Name:    "test",
		Host:    "localhost",
		User:    "testuser",
		Port:    22,
		KeyFile: keyFile,
	}

	_, err = NewSSHClient(server, 30)
	if err == nil {
		t.Error("expected error when key file is invalid")
	}
	if !strings.Contains(err.Error(), "failed to parse private key") {
		t.Errorf("expected key parsing error, got: %v", err)
	}
}

func TestNewSSHClient_BothAuthMethods(t *testing.T) {
	// Create a temporary key file
	tempDir := t.TempDir()
	keyFile := filepath.Join(tempDir, "test_key")

	// Write a dummy key file
	err := os.WriteFile(keyFile, []byte("invalid key content"), 0600)
	if err != nil {
		t.Fatalf("failed to create test key file: %v", err)
	}

	server := &config.Server{
		Name:     "test",
		Host:     "localhost",
		User:     "testuser",
		Port:     22,
		Password: "testpass",
		KeyFile:  keyFile,
	}

	_, err = NewSSHClient(server, 30)
	// Should fail due to key parsing error
	if err == nil {
		t.Error("expected error when key file is invalid")
	}
	if !strings.Contains(err.Error(), "failed to parse private key") {
		t.Errorf("expected key parsing error, got: %v", err)
	}
}

func TestNewSSHClient_KeyFileNotFound(t *testing.T) {
	server := &config.Server{
		Name:    "test",
		Host:    "localhost",
		User:    "testuser",
		Port:    22,
		KeyFile: "/nonexistent/key/file",
	}

	_, err := NewSSHClient(server, 30)
	if err == nil {
		t.Error("expected error when key file does not exist")
	}
	if !strings.Contains(err.Error(), "failed to read key file") {
		t.Errorf("expected file read error, got: %v", err)
	}
}

func TestNewSSHClient_Timeout(t *testing.T) {
	server := &config.Server{
		Name:     "test",
		Host:     "192.0.2.1", // RFC 5737 reserved for documentation/testing
		User:     "testuser",
		Port:     22,
		Password: "testpass",
	}

	start := time.Now()
	_, err := NewSSHClient(server, 1) // 1 second timeout
	duration := time.Since(start)

	if err == nil {
		t.Error("expected error when connecting to unreachable host")
	}

	// Should timeout within reasonable time (not more than 3 seconds)
	if duration > 3*time.Second {
		t.Errorf("connection should have timed out faster, took: %v", duration)
	}
}

func TestExecuteScript_FileReadError(t *testing.T) {
	// Create a mock SSH client
	client := &SSHClient{
		Server: &config.Server{Name: "test"},
	}

	_, err := client.ExecuteScript("/nonexistent/script.sh")
	if err == nil {
		t.Error("expected error when script file does not exist")
	}
	if !strings.Contains(err.Error(), "failed to read script file") {
		t.Errorf("expected file read error, got: %v", err)
	}
}

func TestSSHClient_Close(t *testing.T) {
	// Test that Close doesn't panic when Client is nil
	client := &SSHClient{
		Server: &config.Server{Name: "test"},
		Client: nil,
	}

	// Should not panic
	err := client.Close()
	if err != nil {
		t.Errorf("Close should not return error when Client is nil, got: %v", err)
	}
}

func TestSSHClient_ExecuteCommand_NoConnection(t *testing.T) {
	// Test ExecuteCommand with no SSH connection
	client := &SSHClient{
		Server: &config.Server{Name: "test"},
		Client: nil,
	}

	_, err := client.ExecuteCommand("echo test")
	if err == nil {
		t.Error("expected error when no SSH connection exists")
	}
	if !strings.Contains(err.Error(), "failed to create session") {
		t.Errorf("expected session creation error, got: %v", err)
	}
}
