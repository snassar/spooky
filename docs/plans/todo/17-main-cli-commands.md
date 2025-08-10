# Implementation Plan: Main CLI Commands

## Overview
Implement the main CLI commands for config management, logging, completion, key management, and command acting, following the established CLI patterns and command structure from the workplan and CLI commands documentation.

## Task Details
- **Task ID**: 5.1
- **Priority**: Low
- **File**: `internal/cli/commands/`, `main.go`
- **Functions**: Config, logging, completion, key management, command acting commands

## Current State Analysis

### Existing Patterns
1. **CLI Structure**: CLI commands use `spf13/cobra` framework with consistent patterns
2. **Command Context**: Commands receive context with project, machines, and configuration
3. **Flag Handling**: Consistent flag patterns with `--machines` (plural), `--tags`, no `--filter`
4. **Output Formatting**: Structured output with JSON and HCL formats
5. **Error Handling**: Consistent error wrapping with context
6. **Command Structure**: "spooky noun-verb" pattern for most commands

### Existing Implementation Examples
- **CLI Commands**: `internal/cli/commands/` provides command implementations
- **CLI Types**: `internal/cli/types/command.go` defines command structure
- **Output Formatting**: `internal/cli/help/renderer.go` provides output formatting
- **Flag Management**: `internal/cli/flags/manager.go` provides flag handling

### Unaddressed TODOs
1. **main.go:204**: `// TODO: Implement actual config validation` - Replace placeholder with real config validation using config manager
2. **internal/cli/commands/executor.go:94**: `// TODO: Add more specific validation logic based on command type and arguments` - Implement command-specific validation rules
3. **internal/cli/commands/executor.go:101**: `// TODO: Implement actions command execution` - Implement actions command execution using coordinator's actions manager
4. **internal/cli/commands/executor.go:108**: `// TODO: Implement facts command execution` - Implement facts command execution using coordinator's facts manager
5. **internal/cli/commands/executor.go:115**: `// TODO: Implement machines command execution` - Implement machines command execution using coordinator's machines manager
6. **internal/cli/commands/executor.go:122**: `// TODO: Implement project command execution` - Implement project command execution using coordinator's project manager
7. **internal/cli/commands/executor.go:129**: `// TODO: Implement templates command execution` - Implement templates command execution using coordinator's templates manager
8. **internal/cli/commands/executor.go:136**: `// TODO: Implement config command execution` - Implement config command execution using coordinator's config manager
9. **internal/cli/commands/executor.go:143**: `// TODO: Implement logs command execution` - Implement logs command execution using coordinator's logging manager
10. **internal/cli/commands/executor.go:150**: `// TODO: Implement completion command execution` - Implement completion command execution using coordinator's completion manager
11. **internal/cli/commands/executor.go:157**: `// TODO: Implement keys command execution` - Implement keys command execution using coordinator's secrets manager

### Completed Items
- ✅ **File Renaming**: `internal/actions/manager.go` → `internal/actions/actor.go` for consistency with "act", "acting", "action" terminology
- ✅ **Terminology Consistency**: Updated all action-related methods to use "act", "acting", "action" instead of "execute", "execution"
- ✅ **Collection Planning**: Full implementation of collection planning functionality as specified in `01-4-collection-planning.md`

## Implementation Requirements

### Interface Compliance
The main CLI commands must:
1. **Follow CLI conventions** from `cli-commands.md`
2. **Support global flags** (--config, --verbose, --quiet, --log-level, --log-file, --version)
3. **Implement config management** (show, validate, encrypt)
4. **Implement logging commands** (show with project parameter)
5. **Implement completion generation** (bash, zsh, fish with --output flag)
6. **Implement command acting** for all command types (actions, facts, machines, project, templates, config, logs, completion, keys)
7. **Remove deprecated flags** and commands
8. **Support proper output formats** (JSON, HCL where applicable)
9. **Follow command structure** from workplan
10. **Implement actual config validation** in main application

### Required Dependencies
- CLI framework for command handling
- Config management system
- Logging system
- Completion generation system
- Coordinator system for command acting

## Detailed Implementation Plan

### Step 1: Implement Config Validation in Main Application

**File**: `main.go`

```go
// Update the config validation TODO in main.go
func createConfigCmd() *cobra.Command {
    // ... existing code ...

    validateCmd := &cobra.Command{
        Use:   "validate",
        Short: "Validate config file",
        Long:  "Validate a configuration file for correctness.",
        RunE: func(cmd *cobra.Command, args []string) error {
            configPath := getConfigPath(cmd)

            if _, err := os.Stat(configPath); os.IsNotExist(err) {
                return fmt.Errorf("configuration file does not exist: %s", configPath)
            }

            // Implement actual config validation
            return validateConfigFile(configPath)
        },
    }

    // ... rest of existing code ...
}

// validateConfigFile implements actual config validation
func validateConfigFile(configPath string) error {
    // Load configuration using config manager
    configManager := config.NewManager()
    
    // Load and parse configuration
    cfg, err := configManager.LoadConfig(configPath)
    if err != nil {
        return fmt.Errorf("failed to load configuration: %w", err)
    }

    // Validate configuration structure
    if err := configManager.ValidateConfig(cfg); err != nil {
        return fmt.Errorf("configuration validation failed: %w", err)
    }

    fmt.Printf("Configuration file %s is valid\n", configPath)
    return nil
}
```

### Step 2: Implement Command Acting System

**File**: `internal/cli/commands/executor.go`

```go
// Update the executor to implement all command acting TODOs
package commands

import (
    "context"
    "fmt"
    "spooky/internal/cli/types"
    "spooky/internal/coordinator"
    "spooky/internal/logging"
)

// Executor handles command acting
type Executor struct {
    coordinator *coordinator.Coordinator
    logger      logging.Logger
}

// NewExecutor creates a new executor
func NewExecutor(coordinator *coordinator.Coordinator, logger logging.Logger) *Executor {
    return &Executor{
        coordinator: coordinator,
        logger:      logger,
    }
}

// ValidateActing validates command acting with specific logic
func (e *Executor) ValidateActing(command *spookyclitypes.Command, args []string) error {
    if command == nil {
        return fmt.Errorf("command cannot be nil")
    }

    // Basic validation
    if command.Name == "" {
        return fmt.Errorf("command name cannot be empty")
    }

    if command.Use == "" {
        return fmt.Errorf("command use cannot be empty")
    }

    // Add more specific validation logic based on command type and arguments
    switch command.Type {
    case "actions":
        return e.validateActionsCommand(command, args)
    case "facts":
        return e.validateFactsCommand(command, args)
    case "machines":
        return e.validateMachinesCommand(command, args)
    case "project":
        return e.validateProjectCommand(command, args)
    case "templates":
        return e.validateTemplatesCommand(command, args)
    case "config":
        return e.validateConfigCommand(command, args)
    case "logs":
        return e.validateLogsCommand(command, args)
    case "completion":
        return e.validateCompletionCommand(command, args)
    case "keys":
        return e.validateKeysCommand(command, args)
    default:
        return fmt.Errorf("unknown command type: %s", command.Type)
    }
}

// actActionsCommand implements actions command acting
func (e *Executor) actActionsCommand(command *spookyclitypes.Command, args []string) error {
    e.logger.Info("Acting actions command",
        logging.String("command", command.Name),
        logging.Strings("args", args))

    // Delegate to coordinator's actions manager
    return e.coordinator.ActActionsCommand(command, args)
}

// actFactsCommand implements facts command acting
func (e *Executor) actFactsCommand(command *spookyclitypes.Command, args []string) error {
    e.logger.Info("Acting facts command",
        logging.String("command", command.Name),
        logging.Strings("args", args))

    // Delegate to coordinator's facts manager
    return e.coordinator.ActFactsCommand(command, args)
}

// actMachinesCommand implements machines command acting
func (e *Executor) actMachinesCommand(command *spookyclitypes.Command, args []string) error {
    e.logger.Info("Acting machines command",
        logging.String("command", command.Name),
        logging.Strings("args", args))

    // Delegate to coordinator's machines manager
    return e.coordinator.ActMachinesCommand(command, args)
}

// actProjectCommand implements project command acting
func (e *Executor) actProjectCommand(command *spookyclitypes.Command, args []string) error {
    e.logger.Info("Acting project command",
        logging.String("command", command.Name),
        logging.Strings("args", args))

    // Delegate to coordinator's project manager
    return e.coordinator.ActProjectCommand(command, args)
}

// actTemplatesCommand implements templates command acting
func (e *Executor) actTemplatesCommand(command *spookyclitypes.Command, args []string) error {
    e.logger.Info("Acting templates command",
        logging.String("command", command.Name),
        logging.Strings("args", args))

    // Delegate to coordinator's templates manager
    return e.coordinator.ActTemplatesCommand(command, args)
}

// actConfigCommand implements config command acting
func (e *Executor) actConfigCommand(command *spookyclitypes.Command, args []string) error {
    e.logger.Info("Executing config command",
        logging.String("command", command.Name),
        logging.Strings("args", args))

    // Delegate to coordinator's config manager
    return e.coordinator.ActConfigCommand(command, args)
}

// actLogsCommand implements logs command acting
func (e *Executor) actLogsCommand(command *spookyclitypes.Command, args []string) error {
    e.logger.Info("Executing logs command",
        logging.String("command", command.Name),
        logging.Strings("args", args))

    // Delegate to coordinator's logging manager
    return e.coordinator.ActLogsCommand(command, args)
}

// actCompletionCommand implements completion command acting
func (e *Executor) actCompletionCommand(command *spookyclitypes.Command, args []string) error {
    e.logger.Info("Executing completion command",
        logging.String("command", command.Name),
        logging.Strings("args", args))

    // Delegate to coordinator's completion manager
    return e.coordinator.ActCompletionCommand(command, args)
}

// actKeysCommand implements keys command acting
func (e *Executor) actKeysCommand(command *spookyclitypes.Command, args []string) error {
    e.logger.Info("Executing keys command",
        logging.String("command", command.Name),
        logging.Strings("args", args))

    // Delegate to coordinator's secrets manager
    return e.coordinator.ActKeysCommand(command, args)
}

// Command-specific validation methods
func (e *Executor) validateActionsCommand(command *spookyclitypes.Command, args []string) error {
    // Validate actions-specific requirements
    if len(args) < 1 {
        return fmt.Errorf("actions command requires at least one argument")
    }
    return nil
}

func (e *Executor) validateFactsCommand(command *spookyclitypes.Command, args []string) error {
    // Validate facts-specific requirements
    return nil
}

func (e *Executor) validateMachinesCommand(command *spookyclitypes.Command, args []string) error {
    // Validate machines-specific requirements
    return nil
}

func (e *Executor) validateProjectCommand(command *spookyclitypes.Command, args []string) error {
    // Validate project-specific requirements
    return nil
}

func (e *Executor) validateTemplatesCommand(command *spookyclitypes.Command, args []string) error {
    // Validate templates-specific requirements
    return nil
}

func (e *Executor) validateConfigCommand(command *spookyclitypes.Command, args []string) error {
    // Validate config-specific requirements
    return nil
}

func (e *Executor) validateLogsCommand(command *spookyclitypes.Command, args []string) error {
    // Validate logs-specific requirements
    return nil
}

func (e *Executor) validateCompletionCommand(command *spookyclitypes.Command, args []string) error {
    // Validate completion-specific requirements
    return nil
}

func (e *Executor) validateKeysCommand(command *spookyclitypes.Command, args []string) error {
    // Validate keys-specific requirements
    return nil
}
```

### Step 1: Implement Config Management Commands

**File**: `internal/cli/commands/config.go`

```go
package commands

import (
    "context"
    "fmt"
    "os"
    "path/filepath"

    "github.com/spf13/cobra"
    "spooky/internal/cli/types"
    "spooky/internal/config"
    "spooky/internal/logging"
)

// ConfigManager handles config management commands
type ConfigManager struct {
    configManager config.Manager
    logger        logging.Logger
}

// NewConfigManager creates a new config manager
func NewConfigManager(configManager config.Manager, logger logging.Logger) *ConfigManager {
    return &ConfigManager{
        configManager: configManager,
        logger:        logger,
    }
}

// NewConfigShowCommand creates the config show command
func NewConfigShowCommand(manager *ConfigManager) *cobra.Command {
    var configPath string

    cmd := &cobra.Command{
        Use:   "show",
        Short: "Show config path",
        Long: `Show the current configuration file path.

Examples:
  spooky config show
  spooky config show --config /path/to/config/`,
        RunE: func(cmd *cobra.Command, args []string) error {
            return manager.ShowConfig(cmd.Context(), configPath)
        },
    }

    // Add flags
    cmd.Flags().StringVarP(&configPath, "config", "c", "", "Override config file path")

    return cmd
}

// NewConfigValidateCommand creates the config validate command
func NewConfigValidateCommand(manager *ConfigManager) *cobra.Command {
    var configPath string

    cmd := &cobra.Command{
        Use:   "validate",
        Short: "Validate config file",
        Long: `Validate a configuration file for correctness.

Examples:
  spooky config validate --config /path/to/config/
  spooky config validate`,
        RunE: func(cmd *cobra.Command, args []string) error {
            return manager.ValidateConfig(cmd.Context(), configPath)
        },
    }

    // Add flags
    cmd.Flags().StringVarP(&configPath, "config", "c", "", "Config file to validate")

    return cmd
}

// NewConfigEncryptCommand creates the config encrypt command
func NewConfigEncryptCommand(manager *ConfigManager) *cobra.Command {
    cmd := &cobra.Command{
        Use:   "encrypt",
        Short: "Encrypt config data",
        Long: `Encrypt sensitive data in the configuration file.

Examples:
  spooky config encrypt`,
        RunE: func(cmd *cobra.Command, args []string) error {
            return manager.EncryptConfig(cmd.Context())
        },
    }

    return cmd
}

// ShowConfig shows the current configuration path
func (m *ConfigManager) ShowConfig(ctx context.Context, overridePath string) error {
    var configPath string
    var err error

    if overridePath != "" {
        configPath = overridePath
    } else {
        configPath, err = m.configManager.GetConfigPath()
        if err != nil {
            return fmt.Errorf("failed to get config path: %w", err)
        }
    }

    // Resolve to absolute path
    absPath, err := filepath.Abs(configPath)
    if err != nil {
        return fmt.Errorf("failed to resolve config path: %w", err)
    }

    // Check if config file exists
    if _, err := os.Stat(absPath); os.IsNotExist(err) {
        fmt.Printf("Config path: %s (file does not exist)\n", absPath)
    } else {
        fmt.Printf("Config path: %s\n", absPath)
    }

    return nil
}

// ValidateConfig validates a configuration file
func (m *ConfigManager) ValidateConfig(ctx context.Context, configPath string) error {
    m.logger.Info("Validating configuration",
        logging.String("config_path", configPath))

    var configToValidate string
    var err error

    if configPath != "" {
        configToValidate = configPath
    } else {
        configToValidate, err = m.configManager.GetConfigPath()
        if err != nil {
            return fmt.Errorf("failed to get config path: %w", err)
        }
    }

    // Load and validate configuration
    cfg, err := m.configManager.LoadConfig(configToValidate)
    if err != nil {
        return fmt.Errorf("failed to load configuration: %w", err)
    }

    // Validate configuration
    validationResult, err := m.configManager.ValidateConfig(cfg)
    if err != nil {
        return fmt.Errorf("configuration validation failed: %w", err)
    }

    if validationResult.IsValid {
        fmt.Println("✓ Configuration is valid")
        return nil
    } else {
        fmt.Println("✗ Configuration validation failed:")
        for _, error := range validationResult.Errors {
            fmt.Printf("  - %s\n", error.Message)
        }
        return fmt.Errorf("configuration validation failed")
    }
}

// EncryptConfig encrypts sensitive data in the configuration
func (m *ConfigManager) EncryptConfig(ctx context.Context) error {
    m.logger.Info("Encrypting configuration data")

    // Get config path
    configPath, err := m.configManager.GetConfigPath()
    if err != nil {
        return fmt.Errorf("failed to get config path: %w", err)
    }

    // Load current configuration
    cfg, err := m.configManager.LoadConfig(configPath)
    if err != nil {
        return fmt.Errorf("failed to load configuration: %w", err)
    }

    // Encrypt sensitive data
    err = m.configManager.EncryptConfig(cfg)
    if err != nil {
        return fmt.Errorf("failed to encrypt configuration: %w", err)
    }

    // Save encrypted configuration
    err = m.configManager.SaveConfig(cfg, configPath)
    if err != nil {
        return fmt.Errorf("failed to save encrypted configuration: %w", err)
    }

    fmt.Println("✓ Configuration encrypted successfully")
    return nil
}
```

### Step 2: Implement Logging Commands

**File**: `internal/cli/commands/logs.go`

```go
package commands

import (
    "context"
    "fmt"
    "io"
    "os"
    "path/filepath"

    "github.com/spf13/cobra"
    "spooky/internal/cli/types"
    "spooky/internal/logging"
)

// LogManager handles logging commands
type LogManager struct {
    logger logging.Logger
}

// NewLogManager creates a new log manager
func NewLogManager(logger logging.Logger) *LogManager {
    return &LogManager{
        logger: logger,
    }
}

// NewLogsShowCommand creates the logs show command
func NewLogsShowCommand(manager *LogManager) *cobra.Command {
    var logFile string

    cmd := &cobra.Command{
        Use:   "show <project>",
        Short: "Show logs",
        Long: `Show logs for a project.

Examples:
  spooky logs show ./my-project
  spooky logs show ./my-project --log-file logs/app.log`,
        Args: cobra.ExactArgs(1),
        RunE: func(cmd *cobra.Command, args []string) error {
            projectPath := args[0]
            return manager.ShowLogs(cmd.Context(), projectPath, logFile)
        },
    }

    // Add flags
    cmd.Flags().StringVarP(&logFile, "log-file", "l", "", "Log file path")

    return cmd
}

// ShowLogs shows logs for a project
func (m *LogManager) ShowLogs(ctx context.Context, projectPath, logFile string) error {
    m.logger.Info("Showing logs",
        logging.String("project", projectPath),
        logging.String("log_file", logFile))

    // Validate project path
    if err := m.validateProjectPath(projectPath); err != nil {
        return fmt.Errorf("invalid project path: %w", err)
    }

    // Determine log file path
    var logFilePath string
    if logFile != "" {
        // Use absolute path if provided
        if filepath.IsAbs(logFile) {
            logFilePath = logFile
        } else {
            // Relative to project
            logFilePath = filepath.Join(projectPath, logFile)
        }
    } else {
        // Default log file
        logFilePath = filepath.Join(projectPath, "logs", "spooky.log")
    }

    // Check if log file exists
    if _, err := os.Stat(logFilePath); os.IsNotExist(err) {
        return fmt.Errorf("log file does not exist: %s", logFilePath)
    }

    // Open and read log file
    file, err := os.Open(logFilePath)
    if err != nil {
        return fmt.Errorf("failed to open log file: %w", err)
    }
    defer file.Close()

    // Copy log content to stdout
    _, err = io.Copy(os.Stdout, file)
    if err != nil {
        return fmt.Errorf("failed to read log file: %w", err)
    }

    return nil
}

// validateProjectPath validates that the project path exists and is a valid project
func (m *LogManager) validateProjectPath(projectPath string) error {
    // Check if path exists
    if _, err := os.Stat(projectPath); os.IsNotExist(err) {
        return fmt.Errorf("project path does not exist: %s", projectPath)
    }

    // Check if it's a directory
    if stat, err := os.Stat(projectPath); err == nil && !stat.IsDir() {
        return fmt.Errorf("project path is not a directory: %s", projectPath)
    }

    // Check for project.hcl file
    projectFile := filepath.Join(projectPath, "project.hcl")
    if _, err := os.Stat(projectFile); os.IsNotExist(err) {
        return fmt.Errorf("project.hcl not found in: %s", projectPath)
    }

    return nil
}
```

### Step 3: Implement Completion Commands

**File**: `internal/cli/commands/completion.go`

```go
package commands

import (
    "context"
    "fmt"
    "os"
    "path/filepath"

    "github.com/spf13/cobra"
    "spooky/internal/cli/completion"
    "spooky/internal/logging"
)

// CompletionManager handles completion commands
type CompletionManager struct {
    completionManager completion.Manager
    logger            logging.Logger
}

// NewCompletionManager creates a new completion manager
func NewCompletionManager(completionManager completion.Manager, logger logging.Logger) *CompletionManager {
    return &CompletionManager{
        completionManager: completionManager,
        logger:            logger,
    }
}

// NewCompletionGenerateCommand creates the completion generate command
func NewCompletionGenerateCommand(manager *CompletionManager) *cobra.Command {
    var outputPath string

    cmd := &cobra.Command{
        Use:   "generate <shell>",
        Short: "Generate completion script",
        Long: `Generate shell completion script for spooky.

Supported shells: bash, zsh, fish

Examples:
  spooky completion generate bash
  spooky completion generate bash --output ~/.bash_completion.d/spooky
  spooky completion generate zsh --output ~/.zsh/completions/_spooky
  spooky completion generate fish --output ~/.config/fish/completions/spooky.fish`,
        Args: cobra.ExactArgs(1),
        RunE: func(cmd *cobra.Command, args []string) error {
            shell := args[0]
            return manager.GenerateCompletion(cmd.Context(), shell, outputPath)
        },
    }

    // Add flags
    cmd.Flags().StringVarP(&outputPath, "output", "o", "", "Output file path")

    return cmd
}

// GenerateCompletion generates completion script for the specified shell
func (m *CompletionManager) GenerateCompletion(ctx context.Context, shell, outputPath string) error {
    m.logger.Info("Generating completion script",
        logging.String("shell", shell),
        logging.String("output", outputPath))

    // Validate shell type
    if err := m.validateShell(shell); err != nil {
        return fmt.Errorf("invalid shell: %w", err)
    }

    // Generate completion script
    script, err := m.completionManager.GenerateCompletion(shell)
    if err != nil {
        return fmt.Errorf("failed to generate completion script: %w", err)
    }

    // Determine output destination
    if outputPath == "" {
        // Output to stdout
        fmt.Print(script)
        return nil
    }

    // Create output directory if it doesn't exist
    outputDir := filepath.Dir(outputPath)
    if err := os.MkdirAll(outputDir, 0755); err != nil {
        return fmt.Errorf("failed to create output directory: %w", err)
    }

    // Write script to file
    err = os.WriteFile(outputPath, []byte(script), 0644)
    if err != nil {
        return fmt.Errorf("failed to write completion script: %w", err)
    }

    fmt.Printf("✓ Completion script generated: %s\n", outputPath)
    return nil
}

// validateShell validates the shell type
func (m *CompletionManager) validateShell(shell string) error {
    validShells := []string{"bash", "zsh", "fish"}
    
    for _, validShell := range validShells {
        if shell == validShell {
            return nil
        }
    }
    
    return fmt.Errorf("unsupported shell: %s (supported: %v)", shell, validShells)
}
```

### Step 4: Implement Version Command

**File**: `internal/cli/commands/version.go`

```go
package commands

import (
    "context"
    "fmt"
    "runtime"

    "github.com/spf13/cobra"
    "spooky/internal/logging"
)

// VersionManager handles version commands
type VersionManager struct {
    version string
    logger  logging.Logger
}

// NewVersionManager creates a new version manager
func NewVersionManager(version string, logger logging.Logger) *VersionManager {
    return &VersionManager{
        version: version,
        logger:  logger,
    }
}

// NewVersionCommand creates the version command
func NewVersionCommand(manager *VersionManager) *cobra.Command {
    cmd := &cobra.Command{
        Use:   "--version",
        Short: "Show version information",
        Long: `Show spooky version information.

Note: No -v short flag is available.

Example:
  spooky --version`,
        RunE: func(cmd *cobra.Command, args []string) error {
            return manager.ShowVersion(cmd.Context())
        },
    }

    return cmd
}

// ShowVersion shows version information
func (m *VersionManager) ShowVersion(ctx context.Context) error {
    fmt.Printf("spooky version %s\n", m.version)
    fmt.Printf("  Go version: %s\n", runtime.Version())
    fmt.Printf("  OS/Arch: %s/%s\n", runtime.GOOS, runtime.GOARCH)
    return nil
}
```

### Step 5: Implement Global Flags

**File**: `internal/cli/commands/root.go`

```go
package commands

import (
    "context"
    "fmt"
    "os"

    "github.com/spf13/cobra"
    "spooky/internal/cli/types"
    "spooky/internal/config"
    "spooky/internal/logging"
)

// RootManager handles root command and global flags
type RootManager struct {
    configManager config.Manager
    logger        logging.Logger
    version       string
}

// NewRootManager creates a new root manager
func NewRootManager(configManager config.Manager, logger logging.Logger, version string) *RootManager {
    return &RootManager{
        configManager: configManager,
        logger:        logger,
        version:       version,
    }
}

// NewRootCommand creates the root command
func NewRootCommand(manager *RootManager) *cobra.Command {
    var (
        configPath string
        verbose    bool
        quiet      bool
        logLevel   string
        logFile    string
    )

    cmd := &cobra.Command{
        Use:   "spooky",
        Short: "Spooky - Infrastructure automation tool",
        Long: `Spooky is a powerful infrastructure automation tool that helps you
manage and automate your infrastructure operations.

For more information, visit: https://github.com/spooky/spooky`,
        PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
            return manager.setupGlobalFlags(configPath, verbose, quiet, logLevel, logFile)
        },
    }

    // Add global flags
    cmd.PersistentFlags().StringVarP(&configPath, "config", "c", "", "Path to config directory")
    cmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
    cmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "Suppress output")
    cmd.PersistentFlags().StringVarP(&logLevel, "log-level", "", "info", "Log level (debug, info, warn, error)")
    cmd.PersistentFlags().StringVarP(&logFile, "log-file", "", "", "Log file path")

    // Add version command
    versionManager := NewVersionManager(manager.version, manager.logger)
    cmd.AddCommand(NewVersionCommand(versionManager))

    return cmd
}

// setupGlobalFlags sets up global flags and configuration
func (m *RootManager) setupGlobalFlags(configPath, verbose, quiet, logLevel, logFile string) error {
    // Set up logging level
    if err := m.setupLogging(verbose, quiet, logLevel, logFile); err != nil {
        return fmt.Errorf("failed to setup logging: %w", err)
    }

    // Set up configuration
    if err := m.setupConfig(configPath); err != nil {
        return fmt.Errorf("failed to setup configuration: %w", err)
    }

    return nil
}

// setupLogging sets up logging configuration
func (m *RootManager) setupLogging(verbose, quiet bool, logLevel, logFile string) error {
    // Determine log level
    if verbose {
        logLevel = "debug"
    } else if quiet {
        logLevel = "error"
    }

    // Configure logger
    config := &logging.Config{
        Level:  logLevel,
        Output: "console",
    }

    if logFile != "" {
        config.Output = "file"
        config.FilePath = logFile
    }

    return m.logger.Configure(config)
}

// setupConfig sets up configuration
func (m *RootManager) setupConfig(configPath string) error {
    if configPath != "" {
        return m.configManager.SetConfigPath(configPath)
    }
    return nil
}
```

### Step 6: Implement Command Registration

**File**: `internal/cli/commands/register.go`

```go
package commands

import (
    "spooky/internal/cli/completion"
    "spooky/internal/cli/types"
    "spooky/internal/config"
    "spooky/internal/logging"
)

// RegisterCommands registers all CLI commands
func RegisterCommands(
    configManager config.Manager,
    logger logging.Logger,
    version string,
) (*types.RootCommand, error) {
    
    // Create managers
    rootManager := NewRootManager(configManager, logger, version)
    configManager := NewConfigManager(configManager, logger)
    logManager := NewLogManager(logger)
    completionManager := NewCompletionManager(completion.NewManager(), logger)

    // Create root command
    rootCmd := NewRootCommand(rootManager)

    // Add config commands
    configCmd := &cobra.Command{
        Use:   "config",
        Short: "Manage configuration",
        Long:  "Manage spooky configuration files and settings.",
    }
    configCmd.AddCommand(
        NewConfigShowCommand(configManager),
        NewConfigValidateCommand(configManager),
        NewConfigEncryptCommand(configManager),
    )
    rootCmd.AddCommand(configCmd)

    // Add logs commands
    logsCmd := &cobra.Command{
        Use:   "logs",
        Short: "Manage logs",
        Long:  "Manage and view spooky logs.",
    }
    logsCmd.AddCommand(
        NewLogsShowCommand(logManager),
    )
    rootCmd.AddCommand(logsCmd)

    // Add completion commands
    completionCmd := &cobra.Command{
        Use:   "completion",
        Short: "Generate completion scripts",
        Long:  "Generate shell completion scripts for spooky.",
    }
    completionCmd.AddCommand(
        NewCompletionGenerateCommand(completionManager),
    )
    rootCmd.AddCommand(completionCmd)

    return &types.RootCommand{
        Command: rootCmd,
    }, nil
}
```







### File Operations
- Efficient file reading
- Streaming for large files
- Proper resource cleanup
- Memory usage optimization

## Configuration Options

### Supported Options
- **Config path**: Override default config path
- **Log level**: Debug, info, warn, error
- **Log file**: Custom log file path
- **Verbose**: Enable verbose output
- **Quiet**: Suppress output
- **Output**: Completion script output path

## Dependencies

### Internal Dependencies
- `spooky/internal/cli/types`
- `spooky/internal/cli/completion`
- `spooky/internal/config`
- `spooky/internal/logging`

### External Dependencies
- `github.com/spf13/cobra`
- `context` (standard library)
- `os` (standard library)
- `path/filepath` (standard library)


3. **Global Flags**: Proper global flag handling
4. **Error Handling**: Comprehensive error handling
5. **Performance**: Efficient command execution
6. **Testing**: Comprehensive test coverage
7. **Documentation**: Clear code documentation

## Implementation Order

1. Implement config management commands
2. Add logging commands
3. Implement completion commands
4. Add version command
5. Implement global flags
6. Add command registration
7. Write comprehensive tests
8. Performance optimization
9. Documentation and cleanup



### Integration Risks
- **Command structure changes**: Follow workplan specifications
- **Flag handling changes**: Maintain consistency
- **Error propagation**: Ensure consistent error handling
- **User experience**: Maintain intuitive interface
