package cli

import (
	spookyclicommands "spooky/internal/cli/commands"
	spookyclicompletion "spooky/internal/cli/completion"
	spookycliflags "spooky/internal/cli/flags"
	spookyclihelp "spooky/internal/cli/help"
	spookyclitypes "spooky/internal/cli/types"
	spookylogging "spooky/internal/logging"
)

// NewCLIManager creates a new CLI manager with all dependencies wired together
func NewCLIManager(config *spookyclitypes.Config, logger spookylogging.Logger) CLIManager {
	// Create concrete implementations
	commandBuilder := spookyclicommands.NewBuilder()
	commandExecutor := spookyclicommands.NewExecutor()
	helpRenderer := spookyclihelp.NewRenderer()
	flagsParser := spookycliflags.NewParser()

	// Create sub-managers
	commandsManager := spookyclicommands.NewManager(
		config.CommandsConfig,
		commandBuilder,
		commandExecutor,
		logger,
	)

	// Create completion manager (will be updated when root command is available)
	completionManager := spookyclicompletion.NewManager(
		config.CompletionConfig,
		nil, // generator will be set later
		nil, // root command will be set later
		logger,
	)

	// Create help manager (will be updated when root command is available)
	helpManager := spookyclihelp.NewManager(
		config.HelpConfig,
		helpRenderer,
		nil, // root command will be set later
		logger,
	)

	// Create flags manager
	flagsManager := spookycliflags.NewManager(
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
