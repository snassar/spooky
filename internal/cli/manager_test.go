package cli

import (
	"testing"

	spookyclitypes "spooky/internal/types/cli"
	spookylogging "spooky/internal/logging"
)

func TestNewCLIManager(t *testing.T) {
	// Create a test configuration
	config := &spookyclitypes.Config{
		CommandsConfig: &spookyclitypes.CommandsConfig{
			AutoInitialize:   true,
			ValidateCommands: true,
		},
		CompletionConfig: &spookyclitypes.CompletionConfig{
			EnabledShells: []string{"bash", "zsh"},
		},
		HelpConfig: &spookyclitypes.HelpConfig{
			EnableExamples: true,
			EnableUsage:    true,
		},
		FlagsConfig: &spookyclitypes.FlagsConfig{
			GlobalFlags: make(map[string]interface{}),
		},
		EnableCompletion: true,
		EnableHelp:       true,
	}

	// Create a test logger
	logger := spookylogging.GetLogger()

	// Create CLI manager
	cliManager := NewCLIManager(config, logger)

	// Verify the manager was created
	if cliManager == nil {
		t.Fatal("CLI manager should not be nil")
	}

	// Test initialization
	if err := cliManager.InitializeCommands(); err != nil {
		t.Fatalf("Failed to initialize commands: %v", err)
	}

	// Test getting root command
	rootCmd := cliManager.GetRootCommand()
	if rootCmd == nil {
		t.Fatal("Root command should not be nil after initialization")
	}

	// Test listing commands
	commands := cliManager.ListCommands()
	if len(commands) != 0 {
		t.Logf("Found %d registered commands", len(commands))
	}

	// Test closing
	if err := cliManager.Close(); err != nil {
		t.Fatalf("Failed to close CLI manager: %v", err)
	}
}

func TestCLIManager_RegisterCommand(t *testing.T) {
	config := &spookyclitypes.Config{
		CommandsConfig: &spookyclitypes.CommandsConfig{
			AutoInitialize:   true,
			ValidateCommands: true,
		},
	}
	logger := spookylogging.GetLogger()
	cliManager := NewCLIManager(config, logger)

	// Create a test command
	testCommand := &spookyclitypes.Command{
		Name:  "test",
		Use:   "test",
		Short: "Test command",
		Long:  "A test command for testing",
	}

	// Register the command
	if err := cliManager.RegisterCommand(testCommand); err != nil {
		t.Fatalf("Failed to register command: %v", err)
	}

	// Verify the command was registered
	commands := cliManager.ListCommands()
	if len(commands) != 1 {
		t.Fatalf("Expected 1 command, got %d", len(commands))
	}

	if commands[0].Name != "test" {
		t.Fatalf("Expected command name 'test', got '%s'", commands[0].Name)
	}

	// Test getting the command
	retrievedCommand, err := cliManager.GetCommand("test")
	if err != nil {
		t.Fatalf("Failed to get command: %v", err)
	}

	if retrievedCommand.Name != "test" {
		t.Fatalf("Expected command name 'test', got '%s'", retrievedCommand.Name)
	}

	// Test unregistering the command
	if err := cliManager.UnregisterCommand("test"); err != nil {
		t.Fatalf("Failed to unregister command: %v", err)
	}

	// Verify the command was unregistered
	commands = cliManager.ListCommands()
	if len(commands) != 0 {
		t.Fatalf("Expected 0 commands after unregistering, got %d", len(commands))
	}
}
