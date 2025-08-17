# Code Walkthrough: `spooky secrets validate`

## Overview

This document provides a detailed function-by-function walkthrough of the `spooky secrets validate` command execution flow. This command validates age configuration and keys for a project.

## Command Structure

```bash
spooky secrets validate <project-path>
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

### 4. Command Routing: `secretsValidateCmd.RunE`

**File**: `cmd/secrets.go`
```go
var secretsValidateCmd = &cobra.Command{
    Use:   "validate [project-path]",
    Short: "Validate age configuration and keys for project",
    Long: `Validate age configuration and keys for a project.

This command validates:
- Age configuration in spooky.hcl
- Identity files and permissions
- Recipients file format
- Encrypted values in project files

Examples:
  spooky secrets validate ./my-project`,
    Args: cobra.ExactArgs(1),
    RunE: func(_ *cobra.Command, args []string) error {
        return handleSecretsValidate(args[0])
    },
}
```

**Purpose**: Routes the command to `handleSecretsValidate` with the project path.

### 5. Main Handler: `handleSecretsValidate()`

**File**: `cmd/secrets.go`
```go
func handleSecretsValidate(projectPath string) error {
    ctx := context.Background()

    // Initialize dependencies if not already done
    if secretsManager == nil {
        if err := InitializeSecretsDependencies(); err != nil {
            return fmt.Errorf("failed to initialize secrets dependencies: %w", err)
        }
    }

    fmt.Printf("🔍 Validating secrets configuration in project: %s\n", projectPath)

    // Get secrets integration
    secretsIntegration := getIntegrationManager().GetSecretsIntegration()
    if secretsIntegration == nil {
        return fmt.Errorf("secrets integration not available")
    }

    // Load project configuration
    config, err := loadProjectConfig(projectPath)
    if err != nil {
        return fmt.Errorf("failed to load project configuration: %w", err)
    }

    // Validate secrets configuration
    result, err := secretsIntegration.ValidateSecrets(ctx, config)
    if err != nil {
        return fmt.Errorf("validation failed: %w", err)
    }

    // Display validation results
    fmt.Printf("\n✅ Validation Results:\n")
    fmt.Printf("%s\n", strings.Repeat("─", 50))

    if len(result.Errors) == 0 && len(result.Warnings) == 0 {
        fmt.Printf("🎉 All secrets configuration is valid!\n")
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
        return fmt.Errorf("secrets validation failed")
    }

    fmt.Printf("✅ Validation completed successfully\n")
    return nil
}
```

**Purpose**: 
- Initializes secrets dependencies
- Loads project configuration
- Validates secrets configuration
- Reports detailed results to the user

### 6. Dependency Initialization: `InitializeSecretsDependencies()`

**File**: `cmd/secrets.go`
```go
func InitializeSecretsDependencies() error {
    // Create log manager for secrets component
    logManager := spookylogging.NewLogManager()
    secretsLogger = logManager.GetLogger("secrets")

    // Initialize secrets components
    keyManager := spookysecrets.NewKeyManager(secretsLogger)
    validator := spookysecrets.NewValidator(secretsLogger)
    manager := spookysecrets.NewManager(keyManager, validator, secretsLogger)

    // Create secrets integration
    secretsManager = spookysecrets.NewIntegration(manager)

    return nil
}
```

**Purpose**: Creates logger, key manager, validator, and manager instances for dependency injection.

### 7. Project Configuration Loading: `loadProjectConfig()`

**File**: `cmd/secrets.go`
```go
func loadProjectConfig(projectPath string) (*spookytypes.Config, error) {
    // Load project configuration from spooky.hcl
    configFile := filepath.Join(projectPath, "spooky.hcl")
    
    if _, err := os.Stat(configFile); os.IsNotExist(err) {
        return nil, fmt.Errorf("spooky.hcl not found in project: %s", projectPath)
    }

    // Parse HCL configuration
    parser := hclparse.NewParser()
    file, diags := parser.ParseHCLFile(configFile)
    if diags.HasErrors() {
        return nil, fmt.Errorf("failed to parse spooky.hcl: %v", diags)
    }

    // Extract configuration
    var config spookytypes.Config
    if err := hcl.DecodeBody(file.Body, nil, &config); err != nil {
        return nil, fmt.Errorf("failed to decode spooky.hcl: %w", err)
    }

    return &config, nil
}
```

**Purpose**: 
- Loads spooky.hcl configuration file
- Parses HCL content
- Extracts configuration structure

### 8. Secrets Validation: `secretsIntegration.ValidateSecrets()`

**File**: `internal/secrets/integration.go`
```go
func (i *SecretsIntegration) ValidateSecrets(ctx context.Context, config *spookytypes.Config) (*spookytypes.ValidationResult, error) {
    i.logger.Info("Validating secrets configuration", map[string]interface{}{
        "config_path": config.Path,
    })

    // Use the manager to validate secrets
    result, err := i.manager.ValidateSecrets(ctx, config)
    if err != nil {
        return nil, fmt.Errorf("failed to validate secrets: %w", err)
    }

    if result.Valid {
        i.logger.Info("Secrets validation passed", map[string]interface{}{
            "config_path": config.Path,
        })
    } else {
        i.logger.Warn("Secrets validation failed", map[string]interface{}{
            "config_path": config.Path,
            "errors":      len(result.Errors),
            "warnings":    len(result.Warnings),
        })
    }

    return result, nil
}
```

**Purpose**: Delegates secrets validation to the secrets manager.

### 9. Deep Secrets Validation: `manager.ValidateSecrets()`

**File**: `internal/secrets/manager.go`
```go
func (m *Manager) ValidateSecrets(ctx context.Context, config *spookytypes.Config) (*spookytypes.ValidationResult, error) {
    m.logger.Info("Validating secrets configuration", map[string]interface{}{
        "config_path": config.Path,
    })

    // Use the validator to validate secrets
    result, err := m.validator.ValidateSecrets(ctx, config)
    if err != nil {
        return nil, fmt.Errorf("failed to validate secrets: %w", err)
    }

    if result.Valid {
        m.logger.Info("Secrets validation passed", map[string]interface{}{
            "config_path": config.Path,
        })
    } else {
        m.logger.Warn("Secrets validation failed", map[string]interface{}{
            "config_path": config.Path,
            "errors":      len(result.Errors),
            "warnings":    len(result.Warnings),
        })
    }

    return result, nil
}
```

**Purpose**: Delegates secrets validation to the secrets validator.

### 10. Secrets Validator: `validator.ValidateSecrets()`

**File**: `internal/secrets/validator.go`
```go
func (v *Validator) ValidateSecrets(ctx context.Context, config *spookytypes.Config) (*spookytypes.ValidationResult, error) {
    v.logger.Info("Validating secrets configuration", map[string]interface{}{
        "config_path": config.Path,
    })

    result := &spookytypes.ValidationResult{
        Valid:       true,
        ValidatedAt: time.Now(),
        Errors:      []spookytypesschemas.SchemaError{},
        Warnings:    []spookytypesschemas.SchemaError{},
    }

    // Validate age configuration
    if err := v.validateAgeConfiguration(config.Age); err != nil {
        result.Valid = false
        result.Errors = append(result.Errors, spookytypesschemas.SchemaError{
            Code:     "invalid_age_configuration",
            Message:  fmt.Sprintf("Invalid age configuration: %s", err.Error()),
            Severity: "error",
        })
    }

    // Validate identity files
    if err := v.validateIdentityFiles(config.Age.Identities); err != nil {
        result.Valid = false
        result.Errors = append(result.Errors, spookytypesschemas.SchemaError{
            Code:     "invalid_identity_files",
            Message:  fmt.Sprintf("Invalid identity files: %s", err.Error()),
            Severity: "error",
        })
    }

    // Validate recipients file
    if err := v.validateRecipientsFile(config.Age.RecipientsFile); err != nil {
        result.Valid = false
        result.Errors = append(result.Errors, spookytypesschemas.SchemaError{
            Code:     "invalid_recipients_file",
            Message:  fmt.Sprintf("Invalid recipients file: %s", err.Error()),
            Severity: "error",
        })
    }

    // Validate encrypted values in project files
    if err := v.validateEncryptedValues(config.Path); err != nil {
        result.Warnings = append(result.Warnings, spookytypesschemas.SchemaError{
            Code:     "encrypted_values_validation_failed",
            Message:  fmt.Sprintf("Encrypted values validation failed: %s", err.Error()),
            Severity: "warning",
        })
    }

    return result, nil
}
```

**Purpose**: 
- Validates age configuration
- Validates identity files
- Validates recipients file
- Validates encrypted values

### 11. Age Configuration Validation: `validateAgeConfiguration()`

**File**: `internal/secrets/validator.go`
```go
func (v *Validator) validateAgeConfiguration(ageConfig *spookytypes.AgeConfig) error {
    if ageConfig == nil {
        return fmt.Errorf("age configuration is missing")
    }

    // Validate age version
    if ageConfig.Version == "" {
        return fmt.Errorf("age version is required")
    }

    // Validate age binary path
    if ageConfig.BinaryPath != "" {
        if _, err := os.Stat(ageConfig.BinaryPath); os.IsNotExist(err) {
            return fmt.Errorf("age binary not found at: %s", ageConfig.BinaryPath)
        }
    }

    // Validate identities configuration
    if len(ageConfig.Identities) == 0 {
        return fmt.Errorf("at least one identity file is required")
    }

    // Validate recipients file
    if ageConfig.RecipientsFile == "" {
        return fmt.Errorf("recipients file is required")
    }

    return nil
}
```

**Purpose**: 
- Validates age configuration structure
- Validates age binary path
- Validates identities configuration
- Validates recipients file

### 12. Identity Files Validation: `validateIdentityFiles()`

**File**: `internal/secrets/validator.go`
```go
func (v *Validator) validateIdentityFiles(identities []string) error {
    for i, identityPath := range identities {
        // Check if identity file exists
        if _, err := os.Stat(identityPath); os.IsNotExist(err) {
            return fmt.Errorf("identity file not found: %s", identityPath)
        }

        // Validate identity file permissions
        if err := v.validateIdentityFilePermissions(identityPath); err != nil {
            return fmt.Errorf("identity file %d permissions invalid: %w", i+1, err)
        }

        // Validate identity file format
        if err := v.validateIdentityFileFormat(identityPath); err != nil {
            return fmt.Errorf("identity file %d format invalid: %w", i+1, err)
        }
    }

    return nil
}
```

**Purpose**: 
- Validates identity file existence
- Validates identity file permissions
- Validates identity file format

### 13. Recipients File Validation: `validateRecipientsFile()`

**File**: `internal/secrets/validator.go`
```go
func (v *Validator) validateRecipientsFile(recipientsFile string) error {
    // Check if recipients file exists
    if _, err := os.Stat(recipientsFile); os.IsNotExist(err) {
        return fmt.Errorf("recipients file not found: %s", recipientsFile)
    }

    // Read recipients file
    content, err := os.ReadFile(recipientsFile)
    if err != nil {
        return fmt.Errorf("failed to read recipients file: %w", err)
    }

    // Parse recipients
    recipients := strings.Split(string(content), "\n")
    var validRecipients []string

    for i, recipient := range recipients {
        recipient = strings.TrimSpace(recipient)
        if recipient == "" {
            continue
        }

        // Validate recipient format
        if err := v.validateRecipientFormat(recipient); err != nil {
            return fmt.Errorf("invalid recipient at line %d: %w", i+1, err)
        }

        validRecipients = append(validRecipients, recipient)
    }

    if len(validRecipients) == 0 {
        return fmt.Errorf("no valid recipients found in file")
    }

    return nil
}
```

**Purpose**: 
- Validates recipients file existence
- Validates recipient format
- Validates recipient content

### 14. Encrypted Values Validation: `validateEncryptedValues()`

**File**: `internal/secrets/validator.go`
```go
func (v *Validator) validateEncryptedValues(projectPath string) error {
    // Scan project files for encrypted values
    encryptedFiles, err := v.findEncryptedFiles(projectPath)
    if err != nil {
        return fmt.Errorf("failed to scan for encrypted files: %w", err)
    }

    for _, file := range encryptedFiles {
        if err := v.validateEncryptedFile(file); err != nil {
            return fmt.Errorf("encrypted file validation failed for %s: %w", file, err)
        }
    }

    return nil
}
```

**Purpose**: 
- Scans project for encrypted files
- Validates encrypted file format
- Validates encrypted content

## Key Components

### Secrets Integration
- **Purpose**: Provides secrets functionality to the CLI
- **Responsibilities**: Validation, encryption, decryption
- **Interface**: `spookyinterfaces.SecretsIntegration`

### Secrets Manager
- **Purpose**: Coordinates secrets operations
- **Responsibilities**: Validation, key management, encryption
- **Interface**: `spookyinterfaces.SecretsManager`

### Secrets Validator
- **Purpose**: Validates secrets configurations
- **Responsibilities**: Age configuration validation, key validation, file validation
- **Features**: Comprehensive secrets validation

### Key Manager
- **Purpose**: Manages age keys and identities
- **Responsibilities**: Key validation, identity management, permissions
- **Features**: Secure key handling

## Validation Process

1. **Project Configuration Loading**
   - Loads spooky.hcl configuration
   - Parses HCL content
   - Extracts age configuration

2. **Age Configuration Validation**
   - Validates age version
   - Validates age binary path
   - Validates configuration structure

3. **Identity Files Validation**
   - Validates identity file existence
   - Validates identity file permissions
   - Validates identity file format

4. **Recipients File Validation**
   - Validates recipients file existence
   - Validates recipient format
   - Validates recipient content

5. **Encrypted Values Validation**
   - Scans project for encrypted files
   - Validates encrypted file format
   - Validates encrypted content

## Error Handling

- **Comprehensive Validation**: Validates all aspects of secrets configuration
- **Detailed Error Reporting**: Provides specific error messages with context
- **Warning System**: Reports non-critical issues as warnings
- **File System Validation**: Validates file existence and permissions

## Output Format

### Success Case
```
🔍 Validating secrets configuration in project: /path/to/project

✅ Validation Results:
──────────────────────────────────────────────────
🎉 All secrets configuration is valid!
✅ Validation completed successfully
```

### Error Case
```
🔍 Validating secrets configuration in project: /path/to/project

✅ Validation Results:
──────────────────────────────────────────────────
❌ Errors (2):
  1. Identity file not found: ~/.age/identity.txt
     Context: age.identities[0]
  2. Recipients file not found: ~/.age/recipients.txt
     Context: age.recipients_file

⚠️  Warnings (1):
  1. Encrypted values validation failed: no encrypted files found
     Context: project_scan

❌ Validation failed with 2 errors
```

## Architecture Patterns

- **Interface-Based Design**: Uses interfaces for loose coupling
- **Dependency Injection**: Injects dependencies through constructors
- **Layered Architecture**: Separates concerns across multiple layers
- **Error Aggregation**: Collects and reports all validation issues
- **File System Integration**: Validates file existence and permissions
- **Security-First Validation**: Prioritizes security in validation
