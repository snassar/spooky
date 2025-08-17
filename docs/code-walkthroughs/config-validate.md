# Code Walkthrough: `spooky config validate`

## Overview

This document provides a detailed function-by-function walkthrough of the `spooky config validate` command execution flow. This command validates the global spooky configuration file.

## Command Structure

```bash
spooky config validate [--config <file>]
```

## Execution Flow

### 1. Entry Point: `main()`

**File**: `main.go`
```go
func main() {
    cmd.Execute()
}
```

**Purpose**: Entry point that calls the CLI execution system.

### 2. CLI Setup: `cmd.Execute()`

**File**: `cmd/root.go`
```go
func Execute() {
    err := RootCmd.Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }
}
```

**Purpose**: Executes the root Cobra command with all subcommands.

### 3. Pre-Run Setup: `RootCmd.PersistentPreRunE`

**File**: `cmd/root.go`
```go
PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
    // Skip auto-setup for version and help commands
    if cmd.Name() == "version" || cmd.Name() == "help" {
        return nil
    }

    // Auto-setup configuration for all other commands
    if err := spookyconfig.AutoSetupConfig(); err != nil {
        return fmt.Errorf("configuration setup failed: %w", err)
    }

    // Initialize integration manager for commands that need it
    if err := InitializeIntegrationsDependencies(); err != nil {
        return fmt.Errorf("integration initialization failed: %w", err)
    }

    return nil
}
```

**Purpose**: 
- Sets up global configuration automatically
- Initializes integration dependencies
- Runs before any command execution

### 4. Command Routing: `configValidateCmd.RunE`

**File**: `cmd/config.go`
```go
var configValidateCmd = &cobra.Command{
    Use:   "validate",
    Short: "Validate spooky configuration file",
    Long: `Validate the spooky configuration file for syntax and settings.

This command validates the global spooky configuration file located at
$XDG_CONFIG_HOME/spooky/spooky.hcl (defaulting to $HOME/.config/spooky/spooky.hcl).

Examples:
  spooky config validate
  spooky config validate --config /path/to/spooky.hcl`,
    RunE: func(cmd *cobra.Command, args []string) error {
        configFile, _ := cmd.Flags().GetString("config")
        return handleConfigValidate(configFile)
    },
}
```

**Purpose**: Routes the command to `handleConfigValidate` with optional config file path.

### 5. Main Handler: `handleConfigValidate()`

**File**: `cmd/config.go`
```go
func handleConfigValidate(configFile string) error {
    ctx := context.Background()

    // Initialize dependencies if not already done
    if configManager == nil {
        if err := InitializeConfigDependencies(); err != nil {
            return fmt.Errorf("failed to initialize config dependencies: %w", err)
        }
    }

    // Determine config file path
    if configFile == "" {
        configFile = spookyconfig.GetDefaultConfigPath()
    }

    fmt.Printf("🔍 Validating configuration file: %s\n", configFile)

    // Get config integration
    configIntegration := getIntegrationManager().GetConfigIntegration()
    if configIntegration == nil {
        return fmt.Errorf("config integration not available")
    }

    // Load configuration
    config, err := configIntegration.LoadConfig(ctx, configFile)
    if err != nil {
        return fmt.Errorf("failed to load configuration: %w", err)
    }

    // Validate configuration
    result, err := configIntegration.ValidateConfig(ctx, config)
    if err != nil {
        return fmt.Errorf("validation failed: %w", err)
    }

    // Display validation results
    fmt.Printf("\n✅ Validation Results:\n")
    fmt.Printf("%s\n", strings.Repeat("─", 50))

    if len(result.Errors) == 0 && len(result.Warnings) == 0 {
        fmt.Printf("🎉 Configuration is valid!\n")
        return nil
    }

    // Display errors
    if len(result.Errors) > 0 {
        fmt.Printf("❌ Errors (%d):\n", len(result.Errors))
        for i := range result.Errors {
            err := &result.Errors[i]
            fmt.Printf("  %d. %s\n", i+1, err.Message)
            if err.Context != nil {
                fmt.Printf("     Context: %v\n", err.Context)
            }
        }
        fmt.Printf("\n")
    }

    // Display warnings
    if len(result.Warnings) > 0 {
        fmt.Printf("⚠️  Warnings (%d):\n", len(result.Warnings))
        for i := range result.Warnings {
            warning := &result.Warnings[i]
            fmt.Printf("  %d. %s\n", i+1, warning.Message)
            if warning.Context != nil {
                fmt.Printf("     Context: %v\n", warning.Context)
            }
        }
        fmt.Printf("\n")
    }

    // Summary
    if len(result.Errors) > 0 {
        fmt.Printf("❌ Validation failed with %d errors\n", len(result.Errors))
        return fmt.Errorf("configuration validation failed")
    }

    fmt.Printf("✅ Validation completed successfully\n")
    return nil
}
```

**Purpose**: 
- Initializes config dependencies
- Determines config file path
- Loads configuration
- Validates configuration
- Reports detailed results to the user

### 6. Dependency Initialization: `InitializeConfigDependencies()`

**File**: `cmd/config.go`
```go
func InitializeConfigDependencies() error {
    // Create log manager for config component
    logManager := spookylogging.NewLogManager()
    configLogger = logManager.GetLogger("config")

    // Initialize config components
    loader := spookyconfig.NewLoader(configLogger)
    validator := spookyconfig.NewValidator(configLogger)
    manager := spookyconfig.NewManager(loader, validator, configLogger)

    // Create config integration
    configManager = spookyconfig.NewIntegration(manager)

    return nil
}
```

**Purpose**: Creates logger, loader, validator, and manager instances for dependency injection.

### 7. Config Path Resolution: `spookyconfig.GetDefaultConfigPath()`

**File**: `internal/config/paths.go`
```go
func GetDefaultConfigPath() string {
    // Check XDG_CONFIG_HOME environment variable
    if xdgConfigHome := os.Getenv("XDG_CONFIG_HOME"); xdgConfigHome != "" {
        return filepath.Join(xdgConfigHome, "spooky", "spooky.hcl")
    }

    // Fallback to default location
    homeDir, err := os.UserHomeDir()
    if err != nil {
        // If we can't get home directory, use current directory
        return "spooky.hcl"
    }

    return filepath.Join(homeDir, ".config", "spooky", "spooky.hcl")
}
```

**Purpose**: 
- Checks XDG_CONFIG_HOME environment variable
- Falls back to default location
- Returns appropriate config file path

### 8. Configuration Loading: `configIntegration.LoadConfig()`

**File**: `internal/config/integration.go`
```go
func (i *ConfigIntegration) LoadConfig(ctx context.Context, configPath string) (*spookytypes.Config, error) {
    i.logger.Info("Loading configuration", map[string]interface{}{
        "config_path": configPath,
    })

    // Use the manager to load configuration
    config, err := i.manager.LoadConfig(ctx, configPath)
    if err != nil {
        return nil, fmt.Errorf("failed to load configuration: %w", err)
    }

    i.logger.Info("Configuration loaded successfully", map[string]interface{}{
        "config_path": configPath,
    })

    return config, nil
}
```

**Purpose**: Delegates configuration loading to the config manager.

### 9. Deep Configuration Loading: `manager.LoadConfig()`

**File**: `internal/config/manager.go`
```go
func (m *Manager) LoadConfig(ctx context.Context, configPath string) (*spookytypes.Config, error) {
    m.logger.Info("Loading configuration", map[string]interface{}{
        "config_path": configPath,
    })

    // Use the loader to load configuration
    config, err := m.loader.LoadConfig(ctx, configPath)
    if err != nil {
        return nil, fmt.Errorf("failed to load configuration: %w", err)
    }

    m.logger.Info("Configuration loaded successfully", map[string]interface{}{
        "config_path": configPath,
    })

    return config, nil
}
```

**Purpose**: Delegates configuration loading to the config loader.

### 10. Configuration File Loading: `loader.LoadConfig()`

**File**: `internal/config/loader.go`
```go
func (l *Loader) LoadConfig(ctx context.Context, configPath string) (*spookytypes.Config, error) {
    l.logger.Info("Loading configuration from file", map[string]interface{}{
        "config_path": configPath,
    })

    // Check if config file exists
    if _, err := os.Stat(configPath); os.IsNotExist(err) {
        return nil, fmt.Errorf("configuration file not found: %s", configPath)
    }

    // Read config file
    content, err := os.ReadFile(configPath)
    if err != nil {
        return nil, fmt.Errorf("failed to read configuration file: %w", err)
    }

    // Parse HCL content
    parser := hclparse.NewParser()
    file, diags := parser.ParseHCL(content, configPath)
    if diags.HasErrors() {
        return nil, fmt.Errorf("failed to parse configuration file: %v", diags)
    }

    // Decode configuration
    var config spookytypes.Config
    if err := hcl.DecodeBody(file.Body, nil, &config); err != nil {
        return nil, fmt.Errorf("failed to decode configuration: %w", err)
    }

    // Set config path
    config.Path = configPath

    l.logger.Info("Configuration loaded successfully", map[string]interface{}{
        "config_path": configPath,
    })

    return &config, nil
}
```

**Purpose**: 
- Checks config file existence
- Reads config file content
- Parses HCL content
- Decodes configuration structure
- Sets config path

### 11. Configuration Validation: `configIntegration.ValidateConfig()`

**File**: `internal/config/integration.go`
```go
func (i *ConfigIntegration) ValidateConfig(ctx context.Context, config *spookytypes.Config) (*spookytypes.ValidationResult, error) {
    i.logger.Info("Validating configuration", map[string]interface{}{
        "config_path": config.Path,
    })

    // Use the manager to validate configuration
    result, err := i.manager.ValidateConfig(ctx, config)
    if err != nil {
        return nil, fmt.Errorf("failed to validate configuration: %w", err)
    }

    if result.Valid {
        i.logger.Info("Configuration validation passed", map[string]interface{}{
            "config_path": config.Path,
        })
    } else {
        i.logger.Warn("Configuration validation failed", map[string]interface{}{
            "config_path": config.Path,
            "errors":      len(result.Errors),
            "warnings":    len(result.Warnings),
        })
    }

    return result, nil
}
```

**Purpose**: Delegates configuration validation to the config manager.

### 12. Deep Configuration Validation: `manager.ValidateConfig()`

**File**: `internal/config/manager.go`
```go
func (m *Manager) ValidateConfig(ctx context.Context, config *spookytypes.Config) (*spookytypes.ValidationResult, error) {
    m.logger.Info("Validating configuration", map[string]interface{}{
        "config_path": config.Path,
    })

    // Use the validator to validate configuration
    result, err := m.validator.ValidateConfig(ctx, config)
    if err != nil {
        return nil, fmt.Errorf("failed to validate configuration: %w", err)
    }

    if result.Valid {
        m.logger.Info("Configuration validation passed", map[string]interface{}{
            "config_path": config.Path,
        })
    } else {
        m.logger.Warn("Configuration validation failed", map[string]interface{}{
            "config_path": config.Path,
            "errors":      len(result.Errors),
            "warnings":    len(result.Warnings),
        })
    }

    return result, nil
}
```

**Purpose**: Delegates configuration validation to the config validator.

### 13. Configuration Validator: `validator.ValidateConfig()`

**File**: `internal/config/validator.go`
```go
func (v *Validator) ValidateConfig(ctx context.Context, config *spookytypes.Config) (*spookytypes.ValidationResult, error) {
    v.logger.Info("Validating configuration", map[string]interface{}{
        "config_path": config.Path,
    })

    result := &spookytypes.ValidationResult{
        Valid:       true,
        ValidatedAt: time.Now(),
        Errors:      []spookytypesschemas.SchemaError{},
        Warnings:    []spookytypesschemas.SchemaError{},
    }

    // Validate logging configuration
    if err := v.validateLoggingConfig(config.Logging); err != nil {
        result.Valid = false
        result.Errors = append(result.Errors, spookytypesschemas.SchemaError{
            Code:     "invalid_logging_config",
            Message:  fmt.Sprintf("Invalid logging configuration: %s", err.Error()),
            Severity: "error",
        })
    }

    // Validate SSH configuration
    if err := v.validateSSHConfig(config.SSH); err != nil {
        result.Valid = false
        result.Errors = append(result.Errors, spookytypesschemas.SchemaError{
            Code:     "invalid_ssh_config",
            Message:  fmt.Sprintf("Invalid SSH configuration: %s", err.Error()),
            Severity: "error",
        })
    }

    // Validate facts configuration
    if err := v.validateFactsConfig(config.Facts); err != nil {
        result.Valid = false
        result.Errors = append(result.Errors, spookytypesschemas.SchemaError{
            Code:     "invalid_facts_config",
            Message:  fmt.Sprintf("Invalid facts configuration: %s", err.Error()),
            Severity: "error",
        })
    }

    // Validate age configuration
    if err := v.validateAgeConfig(config.Age); err != nil {
        result.Valid = false
        result.Errors = append(result.Errors, spookytypesschemas.SchemaError{
            Code:     "invalid_age_config",
            Message:  fmt.Sprintf("Invalid age configuration: %s", err.Error()),
            Severity: "error",
        })
    }

    // Validate performance configuration
    if err := v.validatePerformanceConfig(config.Performance); err != nil {
        result.Warnings = append(result.Warnings, spookytypesschemas.SchemaError{
            Code:     "invalid_performance_config",
            Message:  fmt.Sprintf("Invalid performance configuration: %s", err.Error()),
            Severity: "warning",
        })
    }

    return result, nil
}
```

**Purpose**: 
- Validates logging configuration
- Validates SSH configuration
- Validates facts configuration
- Validates age configuration
- Validates performance configuration

### 14. Logging Configuration Validation: `validateLoggingConfig()`

**File**: `internal/config/validator.go`
```go
func (v *Validator) validateLoggingConfig(loggingConfig *spookytypes.LoggingConfig) error {
    if loggingConfig == nil {
        return fmt.Errorf("logging configuration is missing")
    }

    // Validate log level
    if loggingConfig.Level != "" {
        validLevels := []string{"debug", "info", "warn", "error"}
        isValid := false
        for _, level := range validLevels {
            if loggingConfig.Level == level {
                isValid = true
                break
            }
        }
        if !isValid {
            return fmt.Errorf("invalid log level: %s (valid levels: %v)", loggingConfig.Level, validLevels)
        }
    }

    // Validate log format
    if loggingConfig.Format != "" {
        validFormats := []string{"text", "json"}
        isValid := false
        for _, format := range validFormats {
            if loggingConfig.Format == format {
                isValid = true
                break
            }
        }
        if !isValid {
            return fmt.Errorf("invalid log format: %s (valid formats: %v)", loggingConfig.Format, validFormats)
        }
    }

    // Validate log file path if specified
    if loggingConfig.File != "" {
        logDir := filepath.Dir(loggingConfig.File)
        if _, err := os.Stat(logDir); os.IsNotExist(err) {
            return fmt.Errorf("log directory does not exist: %s", logDir)
        }
    }

    return nil
}
```

**Purpose**: 
- Validates log level format
- Validates log format
- Validates log file path
- Validates logging configuration structure

## Key Components

### Config Integration
- **Purpose**: Provides configuration functionality to the CLI
- **Responsibilities**: Loading, validation, management
- **Interface**: `spookyinterfaces.ConfigIntegration`

### Config Manager
- **Purpose**: Coordinates configuration operations
- **Responsibilities**: Loading, validation, management
- **Interface**: `spookyinterfaces.ConfigManager`

### Config Validator
- **Purpose**: Validates configuration settings
- **Responsibilities**: Syntax validation, setting validation, path validation
- **Features**: Comprehensive configuration validation

### Config Loader
- **Purpose**: Loads configuration from disk
- **Responsibilities**: HCL parsing, file loading, path resolution
- **Features**: XDG standard support

## Validation Process

1. **Config Path Resolution**
   - Checks XDG_CONFIG_HOME environment variable
   - Falls back to default location
   - Validates config file existence

2. **Configuration Loading**
   - Reads config file content
   - Parses HCL content
   - Decodes configuration structure

3. **Configuration Validation**
   - Validates logging configuration
   - Validates SSH configuration
   - Validates facts configuration
   - Validates age configuration
   - Validates performance configuration

4. **Setting Validation**
   - Validates log levels and formats
   - Validates SSH settings
   - Validates facts storage settings
   - Validates age encryption settings
   - Validates performance tuning settings

## Error Handling

- **Comprehensive Validation**: Validates all configuration sections
- **Detailed Error Reporting**: Provides specific error messages with context
- **Warning System**: Reports non-critical issues as warnings
- **Path Validation**: Validates file paths and directory existence

## Output Format

### Success Case
```
🔍 Validating configuration file: /home/user/.config/spooky/spooky.hcl

✅ Validation Results:
──────────────────────────────────────────────────
🎉 Configuration is valid!
✅ Validation completed successfully
```

### Error Case
```
🔍 Validating configuration file: /home/user/.config/spooky/spooky.hcl

✅ Validation Results:
──────────────────────────────────────────────────
❌ Errors (2):
  1. Invalid log level: verbose (valid levels: [debug info warn error])
     Context: logging.level
  2. Invalid SSH configuration: port must be between 1 and 65535
     Context: ssh.port

⚠️  Warnings (1):
  1. Invalid performance configuration: max_workers must be positive
     Context: performance.max_workers

❌ Validation failed with 2 errors
```

## Architecture Patterns

- **Interface-Based Design**: Uses interfaces for loose coupling
- **Dependency Injection**: Injects dependencies through constructors
- **Layered Architecture**: Separates concerns across multiple layers
- **Error Aggregation**: Collects and reports all validation issues
- **XDG Standard Support**: Follows XDG base directory specification
- **Comprehensive Validation**: Validates all configuration aspects
