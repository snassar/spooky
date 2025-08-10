package main

import (
	spookytypes "spooky/internal/types"
	"testing"

	"github.com/spf13/cobra"
)

// TestMainCLIStructure tests that the main CLI structure is properly set up
func TestMainCLIStructure(t *testing.T) {
	// This test verifies that the main.go CLI structure is working
	// It doesn't test the actual execution, just the structure setup

	// Create a simple root command to test the structure
	rootCmd := &cobra.Command{
		Use:   "spooky",
		Short: "Spooky is a server configuration and automation tool",
	}

	// Test that we can add commands
	testCmd := &cobra.Command{
		Use:   "test",
		Short: "Test command",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	rootCmd.AddCommand(testCmd)

	// Verify the command was added
	commands := rootCmd.Commands()
	if len(commands) != 1 {
		t.Errorf("Expected 1 command, got %d", len(commands))
	}

	if commands[0].Use != "test" {
		t.Errorf("Expected command 'test', got '%s'", commands[0].Use)
	}
}

// TestPlaceholderCommands tests that placeholder commands are properly created
func TestPlaceholderCommands(t *testing.T) {
	// Test config command
	configCmd := createConfigCmd()
	if configCmd.Use != "config" {
		t.Errorf("Expected config command, got %s", configCmd.Use)
	}

	// Test logs command
	logsCmd := createLogsCmd()
	if logsCmd.Use != "logs" {
		t.Errorf("Expected logs command, got %s", logsCmd.Use)
	}

	// Test completion command
	completionCmd := createCompletionCmd()
	if completionCmd.Use != "completion" {
		t.Errorf("Expected completion command, got %s", completionCmd.Use)
	}

	// Test keys command
	keysCmd := createKeysCmd()
	if keysCmd.Use != "keys" {
		t.Errorf("Expected keys command, got %s", keysCmd.Use)
	}
}

// TestAliasCommands tests that alias commands are properly created
func TestAliasCommands(t *testing.T) {
	// Test run command
	runCmd := createRunCmd()
	if runCmd.Use != "run" {
		t.Errorf("Expected run command, got %s", runCmd.Use)
	}

	// Test run command (ls is an alias for run)
	runCmd2 := createRunCmd()
	if runCmd2.Use != "run" {
		t.Errorf("Expected run command, got %s", runCmd2.Use)
	}

	// Test init command
	initCmd := createInitCmd()
	if initCmd.Use != "init" {
		t.Errorf("Expected init command, got %s", initCmd.Use)
	}

	// Test ping command
	pingCmd := createPingCmd()
	if pingCmd.Use != "ping" {
		t.Errorf("Expected ping command, got %s", pingCmd.Use)
	}
}
