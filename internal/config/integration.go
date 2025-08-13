// Package config provides configuration management functionality for the spooky codebase.
package config

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"

	spookyinterfaces "spooky/internal/interfaces"
	spookytypes "spooky/internal/types"
	spookytypeslogging "spooky/internal/types/logging"
	spookytypesschemas "spooky/internal/types/schemas"
)

// Integration implements the ConfigIntegration interface
type Integration struct {
	logger spookytypeslogging.Logger
}

// NewIntegration creates a new config integration
func NewIntegration(logger spookytypeslogging.Logger) spookyinterfaces.ConfigIntegration {
	return &Integration{
		logger: logger,
	}
}

// LoadConfig loads configuration from the given source
func (i *Integration) LoadConfig(_ context.Context, source string) (*spookytypes.Config, error) {
	if source == "" {
		return nil, fmt.Errorf("config source cannot be empty")
	}

	// Check if source file exists
	if _, err := os.Stat(source); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file not found: %s", source)
	}

	// Parse HCL configuration
	var config spookytypes.Config
	if err := hclsimple.DecodeFile(source, nil, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file %s: %w", source, err)
	}

	i.logger.Info("Configuration loaded successfully", map[string]interface{}{
		"source": source,
	})

	return &config, nil
}

// ValidateConfig validates configuration
func (i *Integration) ValidateConfig(_ context.Context, config *spookytypes.Config) (*spookytypes.ValidationResult, error) {
	if config == nil {
		return &spookytypes.ValidationResult{
			Valid:    false,
			Errors:   []spookytypesschemas.SchemaError{{Message: "config cannot be nil"}},
			Warnings: []spookytypesschemas.SchemaError{},
		}, nil
	}

	// Basic validation
	var errors []spookytypesschemas.SchemaError
	var warnings []spookytypesschemas.SchemaError

	// Validate global configuration
	if config.Global != nil {
		if config.Global.DefaultProjectPath == "" {
			warnings = append(warnings, spookytypesschemas.SchemaError{
				Message: "default project path is recommended",
			})
		}
	}

	// Validate logging configuration
	if config.Logging != nil {
		if config.Logging.Level == "" {
			warnings = append(warnings, spookytypesschemas.SchemaError{
				Message: "logging level is recommended",
			})
		}
	}

	valid := len(errors) == 0

	i.logger.Info("Configuration validation completed", map[string]interface{}{
		"valid":    valid,
		"errors":   len(errors),
		"warnings": len(warnings),
	})

	return &spookytypes.ValidationResult{
		Valid:    valid,
		Errors:   errors,
		Warnings: warnings,
	}, nil
}

// SaveConfig saves configuration to the given destination
func (i *Integration) SaveConfig(_ context.Context, config *spookytypes.Config, destination string) error {
	if config == nil {
		return fmt.Errorf("config cannot be nil")
	}

	if destination == "" {
		return fmt.Errorf("destination cannot be empty")
	}

	// Ensure destination directory exists
	destDir := filepath.Dir(destination)
	if err := os.MkdirAll(destDir, 0o755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	// Create HCL file
	file := hclwrite.NewEmptyFile()
	rootBody := file.Body()

	// Add configuration sections
	i.addConfigMetadata(rootBody, config)
	i.addGlobalConfig(rootBody, config)
	i.addCLIConfig(rootBody, config)
	i.addLoggingConfig(rootBody, config)
	i.addSSHConfig(rootBody, config)
	i.addStorageConfig(rootBody, config)
	i.addSecurityConfig(rootBody, config)

	// Write HCL content to file
	if err := os.WriteFile(destination, file.Bytes(), 0o600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	i.logger.Info("Configuration saved successfully", map[string]interface{}{
		"destination": destination,
	})

	return nil
}

// addConfigMetadata adds configuration metadata to the HCL body
func (i *Integration) addConfigMetadata(rootBody *hclwrite.Body, config *spookytypes.Config) {
	if config.Version != "" {
		rootBody.SetAttributeValue("version", cty.StringVal(config.Version))
	}
	if !config.CreatedAt.IsZero() {
		rootBody.SetAttributeValue("created_at", cty.StringVal(config.CreatedAt.Format(time.RFC3339)))
	}
	if !config.UpdatedAt.IsZero() {
		rootBody.SetAttributeValue("updated_at", cty.StringVal(config.UpdatedAt.Format(time.RFC3339)))
	}
}

// addGlobalConfig adds global configuration to the HCL body
func (i *Integration) addGlobalConfig(rootBody *hclwrite.Body, config *spookytypes.Config) {
	if config.Global == nil {
		return
	}

	globalBlock := rootBody.AppendNewBlock("global", nil)
	globalBody := globalBlock.Body()

	if config.Global.DefaultProjectPath != "" {
		globalBody.SetAttributeValue("default_project_path", cty.StringVal(config.Global.DefaultProjectPath))
	}
	if config.Global.DefaultParallelWorkers > 0 {
		globalBody.SetAttributeValue("default_parallel_workers", cty.NumberIntVal(int64(config.Global.DefaultParallelWorkers)))
	}
	if config.Global.DefaultTimeout > 0 {
		globalBody.SetAttributeValue("default_timeout", cty.NumberIntVal(int64(config.Global.DefaultTimeout)))
	}
	if config.Global.DefaultLogLevel != "" {
		globalBody.SetAttributeValue("default_log_level", cty.StringVal(config.Global.DefaultLogLevel))
	}
	if config.Global.DefaultDryRun {
		globalBody.SetAttributeValue("default_dry_run", cty.BoolVal(config.Global.DefaultDryRun))
	}
	if config.Global.DefaultVerbose {
		globalBody.SetAttributeValue("default_verbose", cty.BoolVal(config.Global.DefaultVerbose))
	}
}

// addCLIConfig adds CLI configuration to the HCL body
func (i *Integration) addCLIConfig(rootBody *hclwrite.Body, config *spookytypes.Config) {
	if config.CLI == nil {
		return
	}

	cliBlock := rootBody.AppendNewBlock("cli", nil)
	cliBody := cliBlock.Body()

	if config.CLI.Theme != "" {
		cliBody.SetAttributeValue("theme", cty.StringVal(config.CLI.Theme))
	}
	if config.CLI.Colors {
		cliBody.SetAttributeValue("colors", cty.BoolVal(config.CLI.Colors))
	}
	if config.CLI.ProgressBars {
		cliBody.SetAttributeValue("progress_bars", cty.BoolVal(config.CLI.ProgressBars))
	}
	if config.CLI.ConfirmPrompts {
		cliBody.SetAttributeValue("confirm_prompts", cty.BoolVal(config.CLI.ConfirmPrompts))
	}
	if config.CLI.OutputFormat != "" {
		cliBody.SetAttributeValue("output_format", cty.StringVal(config.CLI.OutputFormat))
	}
}

// addLoggingConfig adds logging configuration to the HCL body
func (i *Integration) addLoggingConfig(rootBody *hclwrite.Body, config *spookytypes.Config) {
	if config.Logging == nil {
		return
	}

	loggingBlock := rootBody.AppendNewBlock("logging", nil)
	loggingBody := loggingBlock.Body()

	if config.Logging.Level != "" {
		loggingBody.SetAttributeValue("level", cty.StringVal(config.Logging.Level))
	}
	if config.Logging.Format != "" {
		loggingBody.SetAttributeValue("format", cty.StringVal(config.Logging.Format))
	}
	if config.Logging.Output != "" {
		loggingBody.SetAttributeValue("output", cty.StringVal(config.Logging.Output))
	}
	if config.Logging.FilePath != "" {
		loggingBody.SetAttributeValue("file_path", cty.StringVal(config.Logging.FilePath))
	}
	if config.Logging.FilePermissions != "" {
		loggingBody.SetAttributeValue("file_permissions", cty.StringVal(config.Logging.FilePermissions))
	}
	if config.Logging.FileMaxSize > 0 {
		loggingBody.SetAttributeValue("file_max_size", cty.NumberIntVal(int64(config.Logging.FileMaxSize)))
	}
	if config.Logging.FileMaxAge > 0 {
		loggingBody.SetAttributeValue("file_max_age", cty.NumberIntVal(int64(config.Logging.FileMaxAge)))
	}
	if config.Logging.FileMaxBackups > 0 {
		loggingBody.SetAttributeValue("file_max_backups", cty.NumberIntVal(int64(config.Logging.FileMaxBackups)))
	}
	if config.Logging.FileCompress {
		loggingBody.SetAttributeValue("file_compress", cty.BoolVal(config.Logging.FileCompress))
	}
}

// addSSHConfig adds SSH configuration to the HCL body
func (i *Integration) addSSHConfig(rootBody *hclwrite.Body, config *spookytypes.Config) {
	if config.SSH == nil {
		return
	}

	sshBlock := rootBody.AppendNewBlock("ssh", nil)
	sshBody := sshBlock.Body()

	if config.SSH.DefaultPort > 0 {
		sshBody.SetAttributeValue("default_port", cty.NumberIntVal(int64(config.SSH.DefaultPort)))
	}
	if config.SSH.DefaultUser != "" {
		sshBody.SetAttributeValue("default_user", cty.StringVal(config.SSH.DefaultUser))
	}
	if config.SSH.DefaultKeyPath != "" {
		sshBody.SetAttributeValue("default_key_path", cty.StringVal(config.SSH.DefaultKeyPath))
	}
	if config.SSH.ConnectionTimeout > 0 {
		sshBody.SetAttributeValue("connection_timeout", cty.NumberIntVal(int64(config.SSH.ConnectionTimeout)))
	}
	if config.SSH.CommandTimeout > 0 {
		sshBody.SetAttributeValue("command_timeout", cty.NumberIntVal(int64(config.SSH.CommandTimeout)))
	}
	if config.SSH.ConnectionPoolSize > 0 {
		sshBody.SetAttributeValue("connection_pool_size", cty.NumberIntVal(int64(config.SSH.ConnectionPoolSize)))
	}
	if config.SSH.ConnectionPoolTimeout > 0 {
		sshBody.SetAttributeValue("connection_pool_timeout", cty.NumberIntVal(int64(config.SSH.ConnectionPoolTimeout)))
	}
	if config.SSH.EnableConnectionPool {
		sshBody.SetAttributeValue("enable_connection_pool", cty.BoolVal(config.SSH.EnableConnectionPool))
	}
	if config.SSH.EnableHostKeyVerification {
		sshBody.SetAttributeValue("enable_host_key_verification", cty.BoolVal(config.SSH.EnableHostKeyVerification))
	}
	if config.SSH.KnownHostsPath != "" {
		sshBody.SetAttributeValue("known_hosts_path", cty.StringVal(config.SSH.KnownHostsPath))
	}
}

// addStorageConfig adds storage configuration to the HCL body
func (i *Integration) addStorageConfig(rootBody *hclwrite.Body, config *spookytypes.Config) {
	if config.Storage == nil {
		return
	}

	storageBlock := rootBody.AppendNewBlock("storage", nil)
	storageBody := storageBlock.Body()

	if config.Storage.Type != "" {
		storageBody.SetAttributeValue("type", cty.StringVal(config.Storage.Type))
	}
	if config.Storage.Path != "" {
		storageBody.SetAttributeValue("path", cty.StringVal(config.Storage.Path))
	}
	if config.Storage.Format != "" {
		storageBody.SetAttributeValue("format", cty.StringVal(config.Storage.Format))
	}
	if config.Storage.Compression {
		storageBody.SetAttributeValue("compression", cty.BoolVal(config.Storage.Compression))
	}
	if config.Storage.Encryption {
		storageBody.SetAttributeValue("encryption", cty.BoolVal(config.Storage.Encryption))
	}
	if config.Storage.EncryptionKeyPath != "" {
		storageBody.SetAttributeValue("encryption_key_path", cty.StringVal(config.Storage.EncryptionKeyPath))
	}
	if config.Storage.BackupEnabled {
		storageBody.SetAttributeValue("backup_enabled", cty.BoolVal(config.Storage.BackupEnabled))
	}
	if config.Storage.BackupPath != "" {
		storageBody.SetAttributeValue("backup_path", cty.StringVal(config.Storage.BackupPath))
	}
	if config.Storage.BackupRetention > 0 {
		storageBody.SetAttributeValue("backup_retention", cty.NumberIntVal(int64(config.Storage.BackupRetention)))
	}
}

// addSecurityConfig adds security configuration to the HCL body
func (i *Integration) addSecurityConfig(rootBody *hclwrite.Body, config *spookytypes.Config) {
	if config.Security == nil {
		return
	}

	securityBlock := rootBody.AppendNewBlock("security", nil)
	securityBody := securityBlock.Body()

	if config.Security.AuditLogging {
		securityBody.SetAttributeValue("audit_logging", cty.BoolVal(config.Security.AuditLogging))
	}
	if config.Security.AuditLogPath != "" {
		securityBody.SetAttributeValue("audit_log_path", cty.StringVal(config.Security.AuditLogPath))
	}
	if config.Security.SensitiveDataMasking {
		securityBody.SetAttributeValue("sensitive_data_masking", cty.BoolVal(config.Security.SensitiveDataMasking))
	}
	if len(config.Security.SensitiveDataPatterns) > 0 {
		patterns := make([]cty.Value, len(config.Security.SensitiveDataPatterns))
		for i, pattern := range config.Security.SensitiveDataPatterns {
			patterns[i] = cty.StringVal(pattern)
		}
		securityBody.SetAttributeValue("sensitive_data_patterns", cty.ListVal(patterns))
	}
	if config.Security.CertificateVerification {
		securityBody.SetAttributeValue("certificate_verification", cty.BoolVal(config.Security.CertificateVerification))
	}
	if config.Security.CAPath != "" {
		securityBody.SetAttributeValue("ca_path", cty.StringVal(config.Security.CAPath))
	}
}
