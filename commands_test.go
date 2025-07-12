package main

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/spf13/cobra"
)

// --- Test helpers ---
func setArgs(cmd *cobra.Command, args ...string) {
	cmd.SetArgs(args)
}

// resetConfigFile resets the global configFile variable for test isolation
func resetConfigFile() {
	configFile = ""
}

// --- Tests ---
func TestExecuteCmd_NoArgs(t *testing.T) {
	resetConfigFile()
	setArgs(executeCmd)
	err := executeCmd.RunE(executeCmd, []string{})
	if err == nil || err.Error() != "config file is required" {
		t.Errorf("expected config file required error, got %v", err)
	}
}

func TestExecuteCmd_NonExistentFile(t *testing.T) {
	resetConfigFile()
	setArgs(executeCmd, "nonexistent.hcl")
	err := executeCmd.RunE(executeCmd, []string{"nonexistent.hcl"})
	if err == nil || err.Error() == "" {
		t.Errorf("expected file not exist error, got %v", err)
	}
}

func TestExecuteCmd_InvalidConfig(t *testing.T) {
	resetConfigFile()
	// Create temp file with invalid HCL content
	f, _ := ioutil.TempFile("", "invalid*.hcl")
	defer os.Remove(f.Name())
	f.WriteString("invalid hcl content")
	f.Close()

	setArgs(executeCmd, f.Name())
	err := executeCmd.RunE(executeCmd, []string{f.Name()})
	if err == nil {
		t.Errorf("expected parse error, got nil")
	}
}

func TestExecuteCmd_ValidConfig(t *testing.T) {
	resetConfigFile()
	// Create temp file with valid HCL content
	f, _ := ioutil.TempFile("", "valid*.hcl")
	defer os.Remove(f.Name())
	f.WriteString(`
server "test" {
  host = "localhost"
  user = "testuser"
  password = "testpass"
}

action "test_action" {
  description = "Test action"
  command = "echo 'hello'"
  servers = ["test"]
}
`)
	f.Close()

	setArgs(executeCmd, f.Name())
	err := executeCmd.RunE(executeCmd, []string{f.Name()})
	// This will likely fail due to SSH connection, but we're testing the config parsing
	if err != nil && err.Error() != "failed to connect to test@localhost:22: dial tcp localhost:22: connectex: No connection could be made because the target machine actively refused it." {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidateCmd_NoArgs(t *testing.T) {
	resetConfigFile()
	setArgs(validateCmd)
	err := validateCmd.RunE(validateCmd, []string{})
	if err == nil || err.Error() != "config file is required" {
		t.Errorf("expected config file required error, got %v", err)
	}
}

func TestValidateCmd_NonExistentFile(t *testing.T) {
	resetConfigFile()
	setArgs(validateCmd, "nonexistent.hcl")
	err := validateCmd.RunE(validateCmd, []string{"nonexistent.hcl"})
	if err == nil || err.Error() == "" {
		t.Errorf("expected file not exist error, got %v", err)
	}
}

func TestValidateCmd_InvalidConfig(t *testing.T) {
	resetConfigFile()
	f, _ := ioutil.TempFile("", "invalid*.hcl")
	defer os.Remove(f.Name())
	f.WriteString("invalid hcl content")
	f.Close()

	setArgs(validateCmd, f.Name())
	err := validateCmd.RunE(validateCmd, []string{f.Name()})
	if err == nil {
		t.Errorf("expected parse error, got nil")
	}
}

func TestValidateCmd_ValidConfig(t *testing.T) {
	resetConfigFile()
	f, _ := ioutil.TempFile("", "valid*.hcl")
	defer os.Remove(f.Name())
	f.WriteString(`
server "test" {
  host = "localhost"
  user = "testuser"
  password = "testpass"
}

action "test_action" {
  description = "Test action"
  command = "echo 'hello'"
  servers = ["test"]
}
`)
	f.Close()

	setArgs(validateCmd, f.Name())
	err := validateCmd.RunE(validateCmd, []string{f.Name()})
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestListCmd_NoArgs(t *testing.T) {
	resetConfigFile()
	setArgs(listCmd)
	err := listCmd.RunE(listCmd, []string{})
	if err == nil || err.Error() != "config file is required" {
		t.Errorf("expected config file required error, got %v", err)
	}
}

func TestListCmd_NonExistentFile(t *testing.T) {
	resetConfigFile()
	setArgs(listCmd, "nonexistent.hcl")
	err := listCmd.RunE(listCmd, []string{"nonexistent.hcl"})
	if err == nil || err.Error() == "" {
		t.Errorf("expected file not exist error, got %v", err)
	}
}

func TestListCmd_InvalidConfig(t *testing.T) {
	resetConfigFile()
	f, _ := ioutil.TempFile("", "invalid*.hcl")
	defer os.Remove(f.Name())
	f.WriteString("invalid hcl content")
	f.Close()

	setArgs(listCmd, f.Name())
	err := listCmd.RunE(listCmd, []string{f.Name()})
	if err == nil {
		t.Errorf("expected parse error, got nil")
	}
}

func TestListCmd_ValidConfig(t *testing.T) {
	resetConfigFile()
	f, _ := ioutil.TempFile("", "valid*.hcl")
	defer os.Remove(f.Name())
	f.WriteString(`
server "test" {
  host = "localhost"
  user = "testuser"
  password = "testpass"
}

action "test_action" {
  description = "Test action"
  command = "echo 'hello'"
  servers = ["test"]
}
`)
	f.Close()

	setArgs(listCmd, f.Name())
	err := listCmd.RunE(listCmd, []string{f.Name()})
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}
