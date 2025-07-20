package cli

import (
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"

	"spooky/internal/facts"
	"spooky/internal/logging"
)

var initOnce sync.Once

func initCommandsOnce() {
	initOnce.Do(func() {
		InitCommands()
	})
}

func resetCommandState() {
	// Reset command flags
	ValidateCmd.Flags().VisitAll(func(flag *pflag.Flag) {
		_ = flag.Value.Set(flag.DefValue) // Ignore errors in test setup
	})
	ListCmd.Flags().VisitAll(func(flag *pflag.Flag) {
		_ = flag.Value.Set(flag.DefValue) // Ignore errors in test setup
	})
}

func TestValidateCmd_NoArgs(t *testing.T) {
	initCommandsOnce()
	resetCommandState()
	// Test with no arguments - requires exactly 1 argument
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
machine "test" {
  host = "localhost"
  user = "testuser"
  password = "testpass"
}

action "test_action" {
  description = "Test action"
  command = "echo test"
  machines = ["test"]
}
`

	err := os.WriteFile(validConfig, []byte(configContent), 0o600)
	if err != nil {
		t.Fatalf("failed to create valid config file: %v", err)
	}

	ValidateCmd.SetArgs([]string{validConfig})
	err = ValidateCmd.Execute()
	if err != nil {
		t.Errorf("unexpected error with valid config: %v", err)
	}
}

func TestValidateCmd_ConfigWithKeyFile(t *testing.T) {
	initCommandsOnce()
	resetCommandState()
	// Create a temporary valid config file with key file
	tempDir := t.TempDir()
	validConfig := filepath.Join(tempDir, "valid.hcl")

	configContent := `
machine "test" {
  host = "localhost"
  user = "testuser"
  key_file = "~/.ssh/id_rsa"
}

action "test_action" {
  description = "Test action"
  command = "echo test"
  machines = ["test"]
}
`

	err := os.WriteFile(validConfig, []byte(configContent), 0o600)
	if err != nil {
		t.Fatalf("failed to create valid config file: %v", err)
	}

	ValidateCmd.SetArgs([]string{validConfig})
	err = ValidateCmd.Execute()
	if err != nil {
		t.Errorf("unexpected error with valid config: %v", err)
	}
}

func TestValidateCmd_ConfigWithTags(t *testing.T) {
	initCommandsOnce()
	resetCommandState()
	// Create a temporary valid config file with tags
	tempDir := t.TempDir()
	validConfig := filepath.Join(tempDir, "valid.hcl")

	configContent := `
machine "test" {
  host = "localhost"
  user = "testuser"
  password = "testpass"
  tags = {
    environment = "dev"
    type = "test"
  }
}

action "test_action" {
  description = "Test action"
  command = "echo test"
  machines = ["test"]
  tags = ["test"]
}
`

	err := os.WriteFile(validConfig, []byte(configContent), 0o600)
	if err != nil {
		t.Fatalf("failed to create valid config file: %v", err)
	}

	ValidateCmd.SetArgs([]string{validConfig})
	err = ValidateCmd.Execute()
	if err != nil {
		t.Errorf("unexpected error with valid config: %v", err)
	}
}

func TestListCmd_NoArgs(t *testing.T) {
	initCommandsOnce()
	resetCommandState()
	// Test with no arguments - should require a resource type or config file
	ListCmd.SetArgs([]string{})
	err := ListCmd.Execute()
	// List command should require arguments now
	if err == nil {
		t.Error("expected error when no args provided")
	}
}

func TestListCmd_NonExistentFile(t *testing.T) {
	initCommandsOnce()
	resetCommandState()
	// Test with non-existent file
	ListCmd.SetArgs([]string{"/nonexistent/file.hcl"})
	err := ListCmd.Execute()
	if err == nil {
		t.Error("expected error when config file does not exist")
	}
}

func TestListCmd_InvalidConfig(t *testing.T) {
	initCommandsOnce()
	resetCommandState()
	// Create a temporary invalid config file
	tempDir := t.TempDir()
	invalidConfig := filepath.Join(tempDir, "invalid.hcl")

	err := os.WriteFile(invalidConfig, []byte("invalid hcl content"), 0o600)
	if err != nil {
		t.Fatalf("failed to create invalid config file: %v", err)
	}

	ListCmd.SetArgs([]string{invalidConfig})
	err = ListCmd.Execute()
	if err == nil {
		t.Error("expected error when config file is invalid")
	}
}

func TestListCmd_ValidConfig(t *testing.T) {
	initCommandsOnce()
	resetCommandState()
	// Create a temporary valid config file
	tempDir := t.TempDir()
	validConfig := filepath.Join(tempDir, "valid.hcl")

	configContent := `
machine "test" {
  host = "localhost"
  user = "testuser"
  password = "testpass"
}

action "test_action" {
  description = "Test action"
  command = "echo test"
  machines = ["test"]
}
`

	err := os.WriteFile(validConfig, []byte(configContent), 0o600)
	if err != nil {
		t.Fatalf("failed to create valid config file: %v", err)
	}

	ListCmd.SetArgs([]string{validConfig})
	err = ListCmd.Execute()
	if err != nil {
		t.Errorf("unexpected error with valid config: %v", err)
	}
}

func TestListCmd_ConfigWithMultipleServersAndActions(t *testing.T) {
	initCommandsOnce()
	resetCommandState()
	// Create a temporary valid config file with multiple servers and actions
	tempDir := t.TempDir()
	validConfig := filepath.Join(tempDir, "valid.hcl")

	configContent := `
machine "server1" {
  host = "server1.example.com"
  user = "admin"
  password = "pass1"
}

machine "server2" {
  host = "server2.example.com"
  user = "admin"
  password = "pass2"
}

action "action1" {
  description = "First action"
  command = "echo action1"
  machines = ["server1"]
}

action "action2" {
  description = "Second action"
  command = "echo action2"
  machines = ["server2"]
}
`

	err := os.WriteFile(validConfig, []byte(configContent), 0o600)
	if err != nil {
		t.Fatalf("failed to create valid config file: %v", err)
	}

	ListCmd.SetArgs([]string{validConfig})
	err = ListCmd.Execute()
	if err != nil {
		t.Errorf("unexpected error with valid config: %v", err)
	}
}

func TestListCmd_ConfigWithOnlyServers(t *testing.T) {
	initCommandsOnce()
	resetCommandState()
	// Create a temporary valid config file with only servers
	tempDir := t.TempDir()
	validConfig := filepath.Join(tempDir, "valid.hcl")

	configContent := `
machine "server1" {
  host = "server1.example.com"
  user = "admin"
  password = "pass1"
}

machine "server2" {
  host = "server2.example.com"
  user = "admin"
  password = "pass2"
}
`

	err := os.WriteFile(validConfig, []byte(configContent), 0o600)
	if err != nil {
		t.Fatalf("failed to create valid config file: %v", err)
	}

	ListCmd.SetArgs([]string{validConfig})
	err = ListCmd.Execute()
	if err != nil {
		t.Errorf("unexpected error with valid config: %v", err)
	}
}

func TestInitCommands(t *testing.T) {
	initCommandsOnce()
	resetCommandState()

	// Test that ValidateCmd has expected flags
	if ValidateCmd.Flags().Lookup("schema") == nil {
		t.Error("schema flag not found on ValidateCmd")
	}
	if ValidateCmd.Flags().Lookup("strict") == nil {
		t.Error("strict flag not found on ValidateCmd")
	}
	if ValidateCmd.Flags().Lookup("format") == nil {
		t.Error("format flag not found on ValidateCmd")
	}

	// Test that ListCmd has expected flags
	if ListCmd.Flags().Lookup("format") == nil {
		t.Error("format flag not found on ListCmd")
	}
	if ListCmd.Flags().Lookup("filter") == nil {
		t.Error("filter flag not found on ListCmd")
	}
	if ListCmd.Flags().Lookup("sort") == nil {
		t.Error("sort flag not found on ListCmd")
	}
	if ListCmd.Flags().Lookup("reverse") == nil {
		t.Error("reverse flag not found on ListCmd")
	}
}

func TestIsLocalFile(t *testing.T) {
	tests := []struct {
		name     string
		source   string
		expected bool
	}{
		{"local file", "config.hcl", true},
		{"local file with path", "/path/to/config.hcl", true},
		{"github url", "github.com/user/repo", false},
		{"git url", "git://github.com/user/repo", false},
		{"s3 url", "s3://bucket/path", false},
		{"http url", "http://example.com/config.hcl", false},
		{"https url", "https://example.com/config.hcl", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isLocalFile(tt.source)
			if result != tt.expected {
				t.Errorf("isLocalFile(%s) = %v, expected %v", tt.source, result, tt.expected)
			}
		})
	}
}

func TestListFromConfigFile(t *testing.T) {
	// Create a temporary valid config file
	tempDir := t.TempDir()
	validConfig := filepath.Join(tempDir, "valid.hcl")

	configContent := `
machine "server1" {
  host = "server1.example.com"
  user = "admin"
  password = "pass1"
}

machine "server2" {
  host = "server2.example.com"
  user = "admin"
  password = "pass2"
}

action "action1" {
  description = "First action"
  command = "echo action1"
  machines = ["server1"]
}

action "action2" {
  description = "Second action"
  command = "echo action2"
  machines = ["server2"]
}
`

	err := os.WriteFile(validConfig, []byte(configContent), 0o600)
	if err != nil {
		t.Fatalf("failed to create valid config file: %v", err)
	}

	logger := logging.GetLogger()
	err = listFromConfigFile(logger, validConfig)
	if err != nil {
		t.Errorf("unexpected error listing from config file: %v", err)
	}
}

func TestListFromConfigFileEmptyConfig(t *testing.T) {
	// Create a temporary empty config file
	tempDir := t.TempDir()
	emptyConfig := filepath.Join(tempDir, "empty.hcl")

	configContent := `
# Empty configuration file
`

	err := os.WriteFile(emptyConfig, []byte(configContent), 0o600)
	if err != nil {
		t.Fatalf("failed to create empty config file: %v", err)
	}

	logger := logging.GetLogger()
	err = listFromConfigFile(logger, emptyConfig)
	// Empty config should fail validation since it requires at least one machine
	if err == nil {
		t.Error("expected error when config file has no machines")
	}
}

func TestListMachines(t *testing.T) {
	logger := logging.GetLogger()
	err := listMachines(logger)
	// This function is not yet implemented, so it should return an error
	if err == nil {
		t.Error("expected error from unimplemented listMachines function")
	}
}

func TestListFacts(t *testing.T) {
	logger := logging.GetLogger()
	err := listFacts(logger)
	// This function may fail if no facts are available, which is expected
	// We just test that it doesn't panic
	if err != nil {
		// Error is expected if no facts are available
		t.Logf("listFacts returned expected error: %v", err)
	}
}

func TestListFactsWithJSONFormat(t *testing.T) {
	// Set format to JSON
	format = "json"
	defer func() { format = "table" }() // Reset after test

	logger := logging.GetLogger()
	err := listFacts(logger)
	// This function may fail if no facts are available, which is expected
	// We just test that it doesn't panic
	if err != nil {
		// Error is expected if no facts are available
		t.Logf("listFacts returned expected error: %v", err)
	}
}

func TestGetUniqueServers(t *testing.T) {
	// Create test facts
	testFacts := []*facts.Fact{
		{Server: "server1", Key: "os", Value: "linux"},
		{Server: "server1", Key: "version", Value: "1.0"},
		{Server: "server2", Key: "os", Value: "linux"},
		{Server: "server3", Key: "os", Value: "windows"},
	}

	uniqueServers := getUniqueServers(testFacts)
	expected := []string{"server1", "server2", "server3"}

	assert.ElementsMatch(t, expected, uniqueServers)
}

func TestGetUniqueServersEmpty(t *testing.T) {
	uniqueServers := getUniqueServers([]*facts.Fact{})
	assert.Empty(t, uniqueServers)
}

func TestListTemplates(t *testing.T) {
	logger := logging.GetLogger()
	err := listTemplates(logger)
	// This function is not yet implemented, so it should return an error
	if err == nil {
		t.Error("expected error from unimplemented listTemplates function")
	}
}

func TestListConfigs(t *testing.T) {
	logger := logging.GetLogger()
	err := listConfigs(logger)
	// This function is not yet implemented, so it should return an error
	if err == nil {
		t.Error("expected error from unimplemented listConfigs function")
	}
}

func TestListActions(t *testing.T) {
	logger := logging.GetLogger()
	err := listActions(logger)
	// This function is not yet implemented, so it should return an error
	if err == nil {
		t.Error("expected error from unimplemented listActions function")
	}
}

func TestValidateCmdWithInvalidSource(t *testing.T) {
	initCommandsOnce()
	resetCommandState()

	// Test with invalid source (remote URL)
	cmd := ValidateCmd
	cmd.SetArgs([]string{"https://example.com/config.hcl"})
	err := cmd.Execute()
	if err == nil {
		t.Error("expected error when remote source is provided")
	}
}

func TestValidateCmdWithNonExistentFile(t *testing.T) {
	initCommandsOnce()
	resetCommandState()

	// Test with non-existent file
	cmd := ValidateCmd
	cmd.SetArgs([]string{"/nonexistent/file.hcl"})
	err := cmd.Execute()
	if err == nil {
		t.Error("expected error when config file does not exist")
	}
}

func TestListCmdWithConfigFile(t *testing.T) {
	initCommandsOnce()
	resetCommandState()

	// Create a temporary valid config file
	tempDir := t.TempDir()
	validConfig := filepath.Join(tempDir, "valid.hcl")

	configContent := `
machine "test" {
  host = "localhost"
  user = "testuser"
  password = "testpass"
}

action "test_action" {
  description = "Test action"
  command = "echo test"
  machines = ["test"]
}
`

	err := os.WriteFile(validConfig, []byte(configContent), 0o600)
	if err != nil {
		t.Fatalf("failed to create valid config file: %v", err)
	}

	// Test list command with config file
	ListCmd.SetArgs([]string{validConfig})
	err = ListCmd.Execute()
	if err != nil {
		t.Errorf("unexpected error with valid config: %v", err)
	}
}

func TestListCmdWithInvalidResourceType(t *testing.T) {
	initCommandsOnce()
	resetCommandState()

	// Test list command with invalid resource type
	ListCmd.SetArgs([]string{"invalid_resource"})
	err := ListCmd.Execute()
	// This should error since invalid resource types are not supported
	if err == nil {
		t.Error("expected error with invalid resource type")
	}
}

func TestListCmdWithValidResourceType(t *testing.T) {
	initCommandsOnce()
	resetCommandState()

	// Test list command with valid resource type
	ListCmd.SetArgs([]string{"machines"})
	err := ListCmd.Execute()
	// This may error if no machines are configured, which is expected
	if err != nil {
		t.Logf("listCmd with machines returned expected error: %v", err)
	}
}

func TestInitCommandsDuplicate(_ *testing.T) {
	// Test that InitCommands can be called multiple times without issues
	InitCommands()
	InitCommands()
	// Should not panic or cause issues
}
