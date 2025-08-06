package cli

import (
	"spooky/internal/cli/commands"
	"spooky/internal/cli/completion"
	"spooky/internal/cli/flags"
	"spooky/internal/cli/help"
	"spooky/internal/cli/types"
	"spooky/internal/logging"
)

// NewCLIManager creates a new CLI manager with all dependencies wired together
func NewCLIManager(config *types.Config, logger logging.Logger) CLIManager {
	// Create concrete implementations
	commandBuilder := commands.NewBuilder()
	commandExecutor := commands.NewExecutor()
	helpRenderer := help.NewRenderer()
	flagsParser := flags.NewParser()

	// Create sub-managers
	commandsManager := commands.NewManager(
		config.CommandsConfig,
		commandBuilder,
		commandExecutor,
		logger,
	)

	// Create completion manager (will be updated when root command is available)
	completionManager := completion.NewManager(
		config.CompletionConfig,
		nil, // generator will be set later
		nil, // root command will be set later
		logger,
	)

	// Create help manager (will be updated when root command is available)
	helpManager := help.NewManager(
		config.HelpConfig,
		helpRenderer,
		nil, // root command will be set later
		logger,
	)

	// Create flags manager
	flagsManager := flags.NewManager(
		config.FlagsConfig,
		flagsParser,
		logger,
	)

	// Create main CLI manager
	cliManager := NewManager(
		config,
		commandsManager,
		completionManager,
		helpManager,
		flagsManager,
		logger,
	)

	return cliManager
}
