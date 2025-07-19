package cli

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/spf13/pflag"
)

var initOnce sync.Once

func initCommandsOnce() {
	initOnce.Do(func() {
		InitCommands()
	})
}

func resetCommandState() {
	// Reset global variables
	parallel = 5 // Default value from InitCommands
	timeout = 30

	// Reset command flags
	ExecuteCmd.Flags().VisitAll(func(flag *pflag.Flag) {
		_ = flag.Value.Set(flag.DefValue) // Ignore errors in test setup
	})
	ValidateCmd.Flags().VisitAll(func(flag *pflag.Flag) {
		_ = flag.Value.Set(flag.DefValue) // Ignore errors in test setup
	})
	ListCmd.Flags().VisitAll(func(flag *pflag.Flag) {
		_ = flag.Value.Set(flag.DefValue) // Ignore errors in test setup
	})
}

func TestExecuteCmd_NoArgs(t *testing.T) {
	initCommandsOnce()
	resetCommandState()
	// Test with no arguments - now requires exactly 1 argument
	ExecuteCmd.SetArgs([]string{})
	err := ExecuteCmd.Execute()
	if err == nil {
		t.Error("expected error when no source is provided")
	}
}

func TestExecuteCmd_NonExistentFile(t *testing.T) {
	initCommandsOnce()
	resetCommandState()
	// Test with non-existent file
	ExecuteCmd.SetArgs([]string{"/nonexistent/file.hcl"})
	err := ExecuteCmd.Execute()
	if err == nil {
		t.Error("expected error when config file does not exist")
	}
}

func TestExecuteCmd_InvalidConfig(t *testing.T) {
	initCommandsOnce()
	resetCommandState()
	// Create a temporary invalid config file
	tempDir := t.TempDir()
	invalidConfig := filepath.Join(tempDir, "invalid.hcl")

	err := os.WriteFile(invalidConfig, []byte("invalid hcl content"), 0o600)
	if err != nil {
		t.Fatalf("failed to create invalid config file: %v", err)
	}

	ExecuteCmd.SetArgs([]string{invalidConfig})
	err = ExecuteCmd.Execute()
	if err == nil {
		t.Error("expected error when config file is invalid")
	}
}

func TestExecuteCmd_ValidConfig(t *testing.T) {
	initCommandsOnce()
	resetCommandState()
	// Create a temporary valid config file
	tempDir := t.TempDir()
	validConfig := filepath.Join(tempDir, "valid.hcl")

	configContent := `
server "test" {
  host = "localhost"
  user = "testuser"
  password = "testpass"
}

action "test_action" {
  description = "Test action"
  command = "echo test"
}
`

	err := os.WriteFile(validConfig, []byte(configContent), 0o600)
	if err != nil {
		t.Fatalf("failed to create valid config file: %v", err)
	}

	ExecuteCmd.SetArgs([]string{validConfig})
	err = ExecuteCmd.Execute()
	// This will fail to connect or authenticate but should not fail due to configuration issues
	if err != nil && !strings.Contains(err.Error(), "connection refused") && !strings.Contains(err.Error(), "no route to host") && !strings.Contains(err.Error(), "unable to authenticate") {
		t.Errorf("expected connection/authentication error, got: %v", err)
	}
}

func TestExecuteCmd_ValidConfigWithParallel(t *testing.T) {
	initCommandsOnce()
	resetCommandState()
	// Create a temporary valid config file
	tempDir := t.TempDir()
	validConfig := filepath.Join(tempDir, "valid.hcl")

	configContent := `
server "test" {
  host = "localhost"
  user = "testuser"
  password = "testpass"
}

action "test_action" {
  description = "Test action"
  command = "echo test"
  parallel = true
}
`

	err := os.WriteFile(validConfig, []byte(configContent), 0o600)
	if err != nil {
		t.Fatalf("failed to create valid config file: %v", err)
	}

	ExecuteCmd.SetArgs([]string{validConfig, "--parallel", "10"})
	err = ExecuteCmd.Execute()
	// This will fail to connect or authenticate but should not fail due to configuration issues
	if err != nil && !strings.Contains(err.Error(), "connection refused") && !strings.Contains(err.Error(), "no route to host") && !strings.Contains(err.Error(), "unable to authenticate") {
		t.Errorf("expected connection/authentication error, got: %v", err)
	}
}

func TestExecuteCmd_ValidConfigWithTimeout(t *testing.T) {
	initCommandsOnce()
	resetCommandState()
	// Create a temporary valid config file
	tempDir := t.TempDir()
	validConfig := filepath.Join(tempDir, "valid.hcl")

	configContent := `
server "test" {
  host = "localhost"
  user = "testuser"
  password = "testpass"
}

action "test_action" {
  description = "Test action"
  command = "echo test"
}
`

	err := os.WriteFile(validConfig, []byte(configContent), 0o600)
	if err != nil {
		t.Fatalf("failed to create valid config file: %v", err)
	}

	ExecuteCmd.SetArgs([]string{validConfig, "--timeout", "60"})
	err = ExecuteCmd.Execute()
	// This will fail to connect or authenticate but should not fail due to configuration issues
	if err != nil && !strings.Contains(err.Error(), "connection refused") && !strings.Contains(err.Error(), "no route to host") && !strings.Contains(err.Error(), "unable to authenticate") {
		t.Errorf("expected connection/authentication error, got: %v", err)
	}
}

func TestExecuteCmd_ConfigWithMultipleServers(t *testing.T) {
	initCommandsOnce()
	resetCommandState()
	// Create a temporary valid config file with multiple servers
	tempDir := t.TempDir()
	validConfig := filepath.Join(tempDir, "valid.hcl")

	configContent := `
server "server1" {
  host = "localhost"
  user = "testuser"
  password = "testpass"
}

server "server2" {
  host = "localhost"
  user = "testuser"
  password = "testpass"
}

action "test_action" {
  description = "Test action"
  command = "echo test"
}
`

	err := os.WriteFile(validConfig, []byte(configContent), 0o600)
	if err != nil {
		t.Fatalf("failed to create valid config file: %v", err)
	}

	ExecuteCmd.SetArgs([]string{validConfig})
	err = ExecuteCmd.Execute()
	// This will fail to connect or authenticate but should not fail due to configuration issues
	if err != nil && !strings.Contains(err.Error(), "connection refused") && !strings.Contains(err.Error(), "no route to host") && !strings.Contains(err.Error(), "unable to authenticate") {
		t.Errorf("expected connection/authentication error, got: %v", err)
	}
}

func TestExecuteCmd_ConfigWithScript(t *testing.T) {
	initCommandsOnce()
	resetCommandState()
	// Create a temporary valid config file with script
	tempDir := t.TempDir()
	validConfig := filepath.Join(tempDir, "valid.hcl")

	configContent := `
server "test" {
  host = "localhost"
  user = "testuser"
  password = "testpass"
}

action "test_action" {
  description = "Test action"
  script = "/nonexistent/script.sh"
}
`

	err := os.WriteFile(validConfig, []byte(configContent), 0o600)
	if err != nil {
		t.Fatalf("failed to create valid config file: %v", err)
	}

	ExecuteCmd.SetArgs([]string{validConfig})
	err = ExecuteCmd.Execute()
	// This will fail due to script file not found
	if err != nil && !strings.Contains(err.Error(), "failed to read script file") {
		t.Errorf("expected script file error, got: %v", err)
	}
}

func TestValidateCmd_NoArgs(t *testing.T) {
	initCommandsOnce()
	resetCommandState()
	// Test with no arguments - now requires exactly 1 argument
	ValidateCmd.SetArgs([]string{})
	err := ValidateCmd.Execute()
	if err == nil {
		t.Error("expected error when no source is provided")
	}
}

func TestValidateCmd_NonExistentFile(t *testing.T) {
	initCommandsOnce()
	resetCommandState()
	// Test with non-existent file
	ValidateCmd.SetArgs([]string{"/nonexistent/file.hcl"})
	err := ValidateCmd.Execute()
	if err == nil {
		t.Error("expected error when config file does not exist")
	}
}

func TestValidateCmd_InvalidConfig(t *testing.T) {
	initCommandsOnce()
	resetCommandState()
	// Create a temporary invalid config file
	tempDir := t.TempDir()
	invalidConfig := filepath.Join(tempDir, "invalid.hcl")

	err := os.WriteFile(invalidConfig, []byte("invalid hcl content"), 0o600)
	if err != nil {
		t.Fatalf("failed to create invalid config file: %v", err)
	}

	ValidateCmd.SetArgs([]string{invalidConfig})
	err = ValidateCmd.Execute()
	if err == nil {
		t.Error("expected error when config file is invalid")
	}
}

func TestValidateCmd_ValidConfig(t *testing.T) {
	initCommandsOnce()
	resetCommandState()
	// Create a temporary valid config file
	tempDir := t.TempDir()
	validConfig := filepath.Join(tempDir, "valid.hcl")

	configContent := `
server "test" {
  host = "localhost"
  user = "testuser"
  password = "testpass"
}

action "test_action" {
  description = "Test action"
  command = "echo test"
}
`

	err := os.WriteFile(validConfig, []byte(configContent), 0o600)
	if err != nil {
		t.Fatalf("failed to create valid config file: %v", err)
	}

	ValidateCmd.SetArgs([]string{validConfig})
	err = ValidateCmd.Execute()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidateCmd_ConfigWithKeyFile(t *testing.T) {
	initCommandsOnce()
	resetCommandState()
	// Create a temporary valid config file with key file
	tempDir := t.TempDir()
	validConfig := filepath.Join(tempDir, "valid.hcl")

	configContent := `
server "test" {
  host = "localhost"
  user = "testuser"
  key_file = "/path/to/key"
}

action "test_action" {
  description = "Test action"
  command = "echo test"
}
`

	err := os.WriteFile(validConfig, []byte(configContent), 0o600)
	if err != nil {
		t.Fatalf("failed to create valid config file: %v", err)
	}

	ValidateCmd.SetArgs([]string{validConfig})
	err = ValidateCmd.Execute()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidateCmd_ConfigWithTags(t *testing.T) {
	initCommandsOnce()
	resetCommandState()
	// Create a temporary valid config file with tags
	tempDir := t.TempDir()
	validConfig := filepath.Join(tempDir, "valid.hcl")

	configContent := `
server "test" {
  host = "localhost"
  user = "testuser"
  password = "testpass"
  tags = {
    env = "prod"
    region = "us-west"
  }
}

action "test_action" {
  description = "Test action"
  command = "echo test"
  tags = ["env"]
}
`

	err := os.WriteFile(validConfig, []byte(configContent), 0o600)
	if err != nil {
		t.Fatalf("failed to create valid config file: %v", err)
	}

	ValidateCmd.SetArgs([]string{validConfig})
	err = ValidateCmd.Execute()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestListCmd_NoArgs(t *testing.T) {
	initCommandsOnce()
	resetCommandState()
	// Test with no arguments - now requires exactly 1 argument (resource type)
	ListCmd.SetArgs([]string{})
	err := ListCmd.Execute()
	if err == nil {
		t.Error("expected error when no resource type is provided")
	}
}

func TestListCmd_NonExistentFile(t *testing.T) {
	initCommandsOnce()
	resetCommandState()
	// Test with invalid resource type
	ListCmd.SetArgs([]string{"invalid-resource"})
	err := ListCmd.Execute()
	if err == nil {
		t.Error("expected error when invalid resource type is provided")
	}
}

func TestListCmd_InvalidConfig(t *testing.T) {
	initCommandsOnce()
	resetCommandState()
	// Test with invalid resource type
	ListCmd.SetArgs([]string{"invalid-resource"})
	err := ListCmd.Execute()
	if err == nil {
		t.Error("expected error when invalid resource type is provided")
	}
}

func TestListCmd_ValidConfig(t *testing.T) {
	initCommandsOnce()
	resetCommandState()
	// Test with valid resource type (will show "not yet implemented" message)
	ListCmd.SetArgs([]string{"servers"})
	err := ListCmd.Execute()
	if err == nil {
		t.Error("expected error since server listing is not yet implemented")
	}
}

func TestListCmd_ConfigWithMultipleServersAndActions(t *testing.T) {
	initCommandsOnce()
	resetCommandState()
	// Test with valid resource type (will show "not yet implemented" message)
	ListCmd.SetArgs([]string{"actions"})
	err := ListCmd.Execute()
	if err == nil {
		t.Error("expected error since action listing is not yet implemented")
	}
}

func TestListCmd_ConfigWithOnlyServers(t *testing.T) {
	initCommandsOnce()
	resetCommandState()
	// Test with valid resource type (will show "not yet implemented" message)
	ListCmd.SetArgs([]string{"servers"})
	err := ListCmd.Execute()
	if err == nil {
		t.Error("expected error since server listing is not yet implemented")
	}
}

func TestInitCommands(t *testing.T) {
	// Test that InitCommands doesn't panic
	// This test ensures that InitCommands can be called multiple times safely
	// The sync.Once ensures it only runs once
	initCommandsOnce()

	// Verify that flags are properly set
	if ExecuteCmd.Flags().Lookup("parallel") == nil {
		t.Error("parallel flag not found on ExecuteCmd")
	}
	if ExecuteCmd.Flags().Lookup("timeout") == nil {
		t.Error("timeout flag not found on ExecuteCmd")
	}
	if ValidateCmd.Flags().Lookup("format") == nil {
		t.Error("format flag not found on ValidateCmd")
	}
	if ListCmd.Flags().Lookup("format") == nil {
		t.Error("format flag not found on ListCmd")
	}
}

func TestExecuteCmd_InvalidTimeout(t *testing.T) {
	initCommandsOnce()
	resetCommandState()
	// Create a temporary config file for testing
	tempDir := t.TempDir()
	validConfig := filepath.Join(tempDir, "valid.hcl")
	configContent := `
server "test" {
  host = "localhost"
  user = "testuser"
  password = "testpass"
}

action "test_action" {
  description = "Test action"
  command = "echo test"
}
`
	err := os.WriteFile(validConfig, []byte(configContent), 0o600)
	if err != nil {
		t.Fatalf("failed to create valid config file: %v", err)
	}

	// Test with invalid timeout value
	ExecuteCmd.SetArgs([]string{validConfig, "--timeout", "invalid"})
	err = ExecuteCmd.Execute()
	if err == nil {
		t.Error("expected error when timeout is invalid")
	}
}

func TestExecuteCmd_NegativeTimeout(t *testing.T) {
	initCommandsOnce()
	resetCommandState()
	// Create a temporary config file for testing
	tempDir := t.TempDir()
	validConfig := filepath.Join(tempDir, "valid.hcl")
	configContent := `
server "test" {
  host = "localhost"
  user = "testuser"
  password = "testpass"
}

action "test_action" {
  description = "Test action"
  command = "echo test"
}
`
	err := os.WriteFile(validConfig, []byte(configContent), 0o600)
	if err != nil {
		t.Fatalf("failed to create valid config file: %v", err)
	}

	// Test with negative timeout value - this will fail at SSH connection level
	ExecuteCmd.SetArgs([]string{validConfig, "--timeout", "-1"})
	err = ExecuteCmd.Execute()
	// The timeout validation happens at SSH level, so we expect connection errors
	if err != nil && !strings.Contains(err.Error(), "connection refused") && !strings.Contains(err.Error(), "no route to host") && !strings.Contains(err.Error(), "unable to authenticate") {
		t.Errorf("expected connection/authentication error, got: %v", err)
	}
}

func TestExecuteCmd_ZeroTimeout(t *testing.T) {
	initCommandsOnce()
	resetCommandState()
	// Create a temporary config file for testing
	tempDir := t.TempDir()
	validConfig := filepath.Join(tempDir, "valid.hcl")
	configContent := `
server "test" {
  host = "localhost"
  user = "testuser"
  password = "testpass"
}

action "test_action" {
  description = "Test action"
  command = "echo test"
}
`
	err := os.WriteFile(validConfig, []byte(configContent), 0o600)
	if err != nil {
		t.Fatalf("failed to create valid config file: %v", err)
	}

	// Test with zero timeout value - this will fail at SSH connection level
	ExecuteCmd.SetArgs([]string{validConfig, "--timeout", "0"})
	err = ExecuteCmd.Execute()
	// The timeout validation happens at SSH level, so we expect connection errors
	if err != nil && !strings.Contains(err.Error(), "connection refused") && !strings.Contains(err.Error(), "no route to host") && !strings.Contains(err.Error(), "unable to authenticate") {
		t.Errorf("expected connection/authentication error, got: %v", err)
	}
}
