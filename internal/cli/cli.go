// Package cli provides the command-line interface implementation for spooky.
// This package implements the CLI functionality and command management.
package cli

import (
	"context"
	"fmt"
	"os"

	spookyinterfaces "spooky/internal/interfaces"
	spookylogging "spooky/internal/logging"
	spookyschemas "spooky/internal/schemas"
	spookystorage "spooky/internal/storage"
	spookytypes "spooky/internal/types"
)

// CLI represents the main CLI application
type CLI struct {
	// Core managers
	logManager    spookyinterfaces.LogManager
	configManager spookyinterfaces.ConfigManager
	schemaManager spookyinterfaces.SchemaManager
	cliManager    spookyinterfaces.CLIManager

	// Integration manager
	integrationManager spookyinterfaces.IntegrationManager

	// Configuration
	config *spookytypes.Config

	// Logger
	logger spookytypes.Logger
}

// NewCLI creates a new CLI instance
func NewCLI() (*CLI, error) {
	cli := &CLI{}

	// Initialize core managers
	if err := cli.initializeManagers(); err != nil {
		return nil, fmt.Errorf("failed to initialize managers: %w", err)
	}

	// Load configuration
	if err := cli.loadConfiguration(); err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}

	// Configure logging
	if err := cli.configureLogging(); err != nil {
		return nil, fmt.Errorf("failed to configure logging: %w", err)
	}

	// Register commands
	if err := cli.registerCommands(); err != nil {
		return nil, fmt.Errorf("failed to register commands: %w", err)
	}

	return cli, nil
}

// initializeManagers initializes the core managers
func (c *CLI) initializeManagers() error {
	// Initialize log manager
	c.logManager = spookylogging.NewLogManager()

	// Initialize config manager
	c.configManager = spookystorage.NewConfigManager()

	// Initialize schema manager
	c.schemaManager = spookyschemas.NewSchemaManager()

	// Initialize CLI manager
	c.cliManager = NewCLIManager()

	// Initialize integration manager
	c.integrationManager = NewIntegrationManager()

	return nil
}

// loadConfiguration loads the CLI configuration
func (c *CLI) loadConfiguration() error {
	// For now, use default configuration
	c.config = c.configManager.GetDefaultConfig()
	return nil
}

// configureLogging configures logging based on configuration
func (c *CLI) configureLogging() error {
	// Configure logging with default settings
	logConfig := &spookytypes.LogConfig{
		Level:            spookytypes.LogLevelInfo,
		Format:           "text",
		Output:           "stdout",
		IncludeTimestamp: true,
	}

	if err := c.logManager.Configure(logConfig); err != nil {
		return fmt.Errorf("failed to configure logging: %w", err)
	}

	// Get logger for CLI
	c.logger = c.logManager.GetLogger("cli")

	return nil
}

// registerCommands registers all CLI commands
func (c *CLI) registerCommands() error {
	// Register version command
	versionCmd := NewVersionCommand()
	if err := c.cliManager.RegisterCommand(versionCmd); err != nil {
		return fmt.Errorf("failed to register version command: %w", err)
	}

	// Register help command
	helpCmd := NewHelpCommand(c.cliManager)
	if err := c.cliManager.RegisterCommand(helpCmd); err != nil {
		return fmt.Errorf("failed to register help command: %w", err)
	}

	// Register project init command
	projectInitCmd := NewProjectInitCommand(c.integrationManager)
	if err := c.cliManager.RegisterCommand(projectInitCmd); err != nil {
		return fmt.Errorf("failed to register project init command: %w", err)
	}

	// Register project validate command
	projectValidateCmd := NewProjectValidateCommand(c.integrationManager)
	if err := c.cliManager.RegisterCommand(projectValidateCmd); err != nil {
		return fmt.Errorf("failed to register project validate command: %w", err)
	}

	return nil
}

// Run runs the CLI application
func (c *CLI) Run(ctx context.Context, args []string) error {
	c.logger.Info("Starting spooky CLI", map[string]interface{}{
		"args": args,
	})

	// Parse command from arguments
	if len(args) == 0 {
		// Show help if no arguments provided
		return c.cliManager.ShowHelp("")
	}

	commandName := args[0]
	commandArgs := args[1:]

	// Handle special commands
	switch commandName {
	case "--version", "-v":
		return c.cliManager.ShowVersion()
	case "--help", "-h", "help":
		if len(commandArgs) > 0 {
			return c.cliManager.ShowHelp(commandArgs[0])
		}
		return c.cliManager.ShowHelp("")
	}

	// Execute the command
	if err := c.cliManager.ExecuteCommand(ctx, commandName, commandArgs); err != nil {
		c.logger.Error("Command execution failed", err, map[string]interface{}{
			"command": commandName,
			"args":    commandArgs,
		})
		return fmt.Errorf("command '%s' failed: %w", commandName, err)
	}

	c.logger.Info("Command completed successfully", map[string]interface{}{
		"command": commandName,
	})

	return nil
}

// Close closes the CLI application
func (c *CLI) Close() error {
	c.logger.Info("Shutting down spooky CLI")

	// Close log manager
	if err := c.logManager.Close(); err != nil {
		return fmt.Errorf("failed to close log manager: %w", err)
	}

	return nil
}

// Main is the main entry point for the CLI
func Main() {
	// Create CLI instance
	cli, err := NewCLI()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize CLI: %v\n", err)
		os.Exit(1)
	}
	defer cli.Close()

	// Run CLI
	ctx := context.Background()
	if err := cli.Run(ctx, os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
