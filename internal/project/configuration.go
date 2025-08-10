package project

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"spooky/internal/config"
	"spooky/internal/types/config"
	"spooky/internal/logging"
)

// ProjectConfigurationEngine manages project-specific configuration
type ProjectConfigurationEngine struct {
	configManager config.ConfigManager
	loader        config.Loader
}

// NewProjectConfigurationEngine creates a new project configuration engine
func NewProjectConfigurationEngine() *ProjectConfigurationEngine {
	return &ProjectConfigurationEngine{
		loader: nil, // Will be initialized when needed
	}
}

// ConfigurationContext represents the resolved configuration context for a project
type ConfigurationContext struct {
	Project     *Project
	Global      *types.GlobalConfig
	Environment map[string]string
	CLI         map[string]interface{}
	Resolved    *ResolvedConfiguration
}

// ResolvedConfiguration represents the final resolved configuration with precedence applied
type ResolvedConfiguration struct {
	// Project settings
	Name        string
	Description string
	Version     string
	Environment string
	Region      string
	Tags        []string

	// Execution settings
	DefaultTimeout        int
	MaxParallel           int
	DryRunDefault         bool
	ValidateBeforeExecute bool
	BackupBeforeChanges   bool

	// Structure settings
	TemplatesDir string
	DataDir      string
	ScriptsDir   string
	LogsDir      string
	BackupsDir   string

	// Isolation settings
	IsolationEnabled bool
	FactsScope       string
	VariablesScope   string
	MachineAccess    string
	AllowedMachines  []string
	AllowedTags      []string

	// SSH settings (inherited from global)
	SSHDefaultUser           string
	SSHDefaultPort           int
	SSHConnectionTimeout     int
	SSHCommandTimeout        int
	SSHRetryAttempts         int
	SSHRetryDelay            int
	SSHKeyPath               string
	SSHStrictHostKeyChecking bool
	SSHAllowPasswordAuth     bool

	// Facts settings
	FactsStorageType string
	FactsStoragePath string

	// Template settings
	TemplateDataDirectory      string
	TemplateTemplatesDirectory string
	TemplateAutoLoadData       bool
	TemplateStrictValidation   bool

	// Security settings
	SecurityAllowInsecureConnections bool
	SecurityAllowPasswordAuth        bool
	SecurityRequireHTTPSImports      bool
	SecurityMaxFileSize              string
}

// LoadProjectConfiguration loads and resolves project configuration with precedence
func (ce *ProjectConfigurationEngine) LoadProjectConfiguration(projectPath string) (*ConfigurationContext, error) {
	logger := logging.GetLogger()

	// 1. Load global configuration (base)
	globalConfig, err := ce.loadGlobalConfiguration()
	if err != nil {
		return nil, fmt.Errorf("failed to load global config: %w", err)
	}

	// 2. Load project configuration
	project, err := ce.loadProjectFile(projectPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load project config: %w", err)
	}

	// 3. Load environment variables
	envVars := ce.loadEnvironmentVariables()

	// 4. Create configuration context
	context := &ConfigurationContext{
		Project:     project,
		Global:      globalConfig,
		Environment: envVars,
		CLI:         make(map[string]interface{}), // Will be populated by CLI
	}

	// 5. Resolve configuration with precedence
	resolved, err := ce.resolveConfiguration(context)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve configuration: %w", err)
	}

	context.Resolved = resolved

	logger.Info("Project configuration loaded",
		logging.String("project", project.Name),
		logging.String("environment", resolved.Environment))

	return context, nil
}

// ApplyCLIOverrides applies CLI flag overrides to the configuration context
func (ce *ProjectConfigurationEngine) ApplyCLIOverrides(context *ConfigurationContext, cliFlags map[string]interface{}) error {
	if context == nil || context.Resolved == nil {
		return fmt.Errorf("configuration context is not initialized")
	}

	// Apply CLI overrides to resolved configuration
	for key, value := range cliFlags {
		switch key {
		case "timeout":
			if timeout, ok := value.(int); ok {
				context.Resolved.DefaultTimeout = timeout
			}
		case "parallel":
			if parallel, ok := value.(int); ok {
				context.Resolved.MaxParallel = parallel
			}
		case "dry-run":
			if dryRun, ok := value.(bool); ok {
				context.Resolved.DryRunDefault = dryRun
			}
		case "environment":
			if env, ok := value.(string); ok {
				context.Resolved.Environment = env
			}
		case "ssh-timeout":
			if timeout, ok := value.(int); ok {
				context.Resolved.SSHConnectionTimeout = timeout
			}
		case "ssh-user":
			if user, ok := value.(string); ok {
				context.Resolved.SSHDefaultUser = user
			}
		case "ssh-port":
			if port, ok := value.(int); ok {
				context.Resolved.SSHDefaultPort = port
			}
		}
	}

	// Store CLI flags for reference
	context.CLI = cliFlags

	return nil
}

// GetConfigurationSummary returns a summary of the resolved configuration
func (ce *ProjectConfigurationEngine) GetConfigurationSummary(context *ConfigurationContext) map[string]interface{} {
	if context == nil || context.Resolved == nil {
		return nil
	}

	summary := map[string]interface{}{
		"project": map[string]interface{}{
			"name":        context.Resolved.Name,
			"description": context.Resolved.Description,
			"version":     context.Resolved.Version,
			"environment": context.Resolved.Environment,
			"region":      context.Resolved.Region,
			"tags":        context.Resolved.Tags,
		},
		"execution": map[string]interface{}{
			"default_timeout":         context.Resolved.DefaultTimeout,
			"max_parallel":            context.Resolved.MaxParallel,
			"dry_run_default":         context.Resolved.DryRunDefault,
			"validate_before_execute": context.Resolved.ValidateBeforeExecute,
			"backup_before_changes":   context.Resolved.BackupBeforeChanges,
		},
		"structure": map[string]interface{}{
			"templates_dir": context.Resolved.TemplatesDir,
			"data_dir":      context.Resolved.DataDir,
			"scripts_dir":   context.Resolved.ScriptsDir,
			"logs_dir":      context.Resolved.LogsDir,
			"backups_dir":   context.Resolved.BackupsDir,
		},
		"isolation": map[string]interface{}{
			"enabled":          context.Resolved.IsolationEnabled,
			"facts_scope":      context.Resolved.FactsScope,
			"variables_scope":  context.Resolved.VariablesScope,
			"machine_access":   context.Resolved.MachineAccess,
			"allowed_machines": context.Resolved.AllowedMachines,
			"allowed_tags":     context.Resolved.AllowedTags,
		},
		"ssh": map[string]interface{}{
			"default_user":             context.Resolved.SSHDefaultUser,
			"default_port":             context.Resolved.SSHDefaultPort,
			"connection_timeout":       context.Resolved.SSHConnectionTimeout,
			"command_timeout":          context.Resolved.SSHCommandTimeout,
			"retry_attempts":           context.Resolved.SSHRetryAttempts,
			"retry_delay":              context.Resolved.SSHRetryDelay,
			"key_path":                 context.Resolved.SSHKeyPath,
			"strict_host_key_checking": context.Resolved.SSHStrictHostKeyChecking,
			"allow_password_auth":      context.Resolved.SSHAllowPasswordAuth,
		},
		"facts": map[string]interface{}{
			"storage_type": context.Resolved.FactsStorageType,
			"storage_path": context.Resolved.FactsStoragePath,
		},
		"templates": map[string]interface{}{
			"data_directory":      context.Resolved.TemplateDataDirectory,
			"templates_directory": context.Resolved.TemplateTemplatesDirectory,
			"auto_load_data":      context.Resolved.TemplateAutoLoadData,
			"strict_validation":   context.Resolved.TemplateStrictValidation,
		},
		"security": map[string]interface{}{
			"allow_insecure_connections": context.Resolved.SecurityAllowInsecureConnections,
			"allow_password_auth":        context.Resolved.SecurityAllowPasswordAuth,
			"require_https_imports":      context.Resolved.SecurityRequireHTTPSImports,
			"max_file_size":              context.Resolved.SecurityMaxFileSize,
		},
	}

	return summary
}

// loadGlobalConfiguration loads the global configuration
func (ce *ProjectConfigurationEngine) loadGlobalConfiguration() (*types.GlobalConfig, error) {
	// For now, return a default global config
	// In a real implementation, this would use the config manager
	globalConfig := &types.GlobalConfig{
		LogLevel:   "info",
		Quiet:      false,
		Verbose:    false,
		ConfigPath: "",
	}

	return globalConfig, nil
}

// loadProjectFile loads the project configuration file
func (ce *ProjectConfigurationEngine) loadProjectFile(projectPath string) (*Project, error) {
	// This would typically parse the project.hcl file
	// For now, we'll create a basic project structure
	project := &Project{
		Path:      projectPath,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Set defaults
	project.SetDefaults()

	return project, nil
}

// loadEnvironmentVariables loads environment variable overrides
func (ce *ProjectConfigurationEngine) loadEnvironmentVariables() map[string]string {
	envVars := make(map[string]string)

	// Project-specific environment variables
	envVars["SPOOKY_PROJECT_ENVIRONMENT"] = os.Getenv("SPOOKY_PROJECT_ENVIRONMENT")
	envVars["SPOOKY_PROJECT_REGION"] = os.Getenv("SPOOKY_PROJECT_REGION")
	envVars["SPOOKY_PROJECT_TAGS"] = os.Getenv("SPOOKY_PROJECT_TAGS")

	// Execution settings
	envVars["SPOOKY_DEFAULT_TIMEOUT"] = os.Getenv("SPOOKY_DEFAULT_TIMEOUT")
	envVars["SPOOKY_MAX_PARALLEL"] = os.Getenv("SPOOKY_MAX_PARALLEL")
	envVars["SPOOKY_DRY_RUN"] = os.Getenv("SPOOKY_DRY_RUN")

	// SSH settings
	envVars["SPOOKY_SSH_TIMEOUT"] = os.Getenv("SPOOKY_SSH_TIMEOUT")
	envVars["SPOOKY_SSH_USER"] = os.Getenv("SPOOKY_SSH_USER")
	envVars["SPOOKY_SSH_PORT"] = os.Getenv("SPOOKY_SSH_PORT")
	envVars["SPOOKY_SSH_KEY_PATH"] = os.Getenv("SPOOKY_SSH_KEY_PATH")

	// Facts settings
	envVars["SPOOKY_FACTS_PATH"] = os.Getenv("SPOOKY_FACTS_PATH")
	envVars["SPOOKY_FACTS_FORMAT"] = os.Getenv("SPOOKY_FACTS_FORMAT")

	return envVars
}

// resolveConfiguration resolves configuration with precedence (CLI > Project > Global > Defaults)
func (ce *ProjectConfigurationEngine) resolveConfiguration(context *ConfigurationContext) (*ResolvedConfiguration, error) {
	resolved := &ResolvedConfiguration{}

	// Start with defaults
	ce.applyDefaults(resolved)

	// Apply global configuration
	if err := ce.applyGlobalConfiguration(resolved, context.Global); err != nil {
		return nil, fmt.Errorf("failed to apply global configuration: %w", err)
	}

	// Apply project configuration
	if err := ce.applyProjectConfiguration(resolved, context.Project); err != nil {
		return nil, fmt.Errorf("failed to apply project configuration: %w", err)
	}

	// Apply environment variable overrides
	ce.applyEnvironmentOverrides(resolved, context.Environment)

	// Validate resolved configuration
	if err := ce.validateResolvedConfiguration(resolved); err != nil {
		return nil, fmt.Errorf("resolved configuration validation failed: %w", err)
	}

	return resolved, nil
}

// applyDefaults applies default values to the resolved configuration
func (ce *ProjectConfigurationEngine) applyDefaults(resolved *ResolvedConfiguration) {
	resolved.Environment = DefaultProjectEnvironment
	resolved.DefaultTimeout = DefaultTimeout
	resolved.MaxParallel = DefaultMaxParallel
	resolved.DryRunDefault = false
	resolved.ValidateBeforeExecute = true
	resolved.BackupBeforeChanges = false

	resolved.TemplatesDir = DefaultTemplatesDir
	resolved.DataDir = DefaultDataDir
	resolved.ScriptsDir = DefaultScriptsDir
	resolved.LogsDir = DefaultLogsDir
	resolved.BackupsDir = DefaultBackupsDir

	resolved.IsolationEnabled = true
	resolved.FactsScope = DefaultFactsScope
	resolved.VariablesScope = DefaultVariablesScope
	resolved.MachineAccess = DefaultMachineAccess

	// SSH defaults
	resolved.SSHDefaultUser = "ubuntu"
	resolved.SSHDefaultPort = 22
	resolved.SSHConnectionTimeout = 30
	resolved.SSHCommandTimeout = 300
	resolved.SSHRetryAttempts = 3
	resolved.SSHRetryDelay = 5
	resolved.SSHKeyPath = "~/.ssh/id_rsa"
	resolved.SSHStrictHostKeyChecking = true
	resolved.SSHAllowPasswordAuth = false

	// Facts defaults
	resolved.FactsStorageType = "badgerdb"
	resolved.FactsStoragePath = "facts.db"

	// Template defaults
	resolved.TemplateDataDirectory = "data"
	resolved.TemplateTemplatesDirectory = "templates"
	resolved.TemplateAutoLoadData = true
	resolved.TemplateStrictValidation = true

	// Security defaults
	resolved.SecurityAllowInsecureConnections = false
	resolved.SecurityAllowPasswordAuth = false
	resolved.SecurityRequireHTTPSImports = true
	resolved.SecurityMaxFileSize = "100MB"
}

// applyGlobalConfiguration applies global configuration settings
func (ce *ProjectConfigurationEngine) applyGlobalConfiguration(resolved *ResolvedConfiguration, global *types.GlobalConfig) error {
	if global == nil {
		return nil
	}

	// Apply SSH settings from global config
	if global.SSH != nil {
		if global.SSH.Timeout != 0 {
			resolved.SSHConnectionTimeout = global.SSH.Timeout
		}
		// Note: Global SSH config has limited fields, so we use defaults for most SSH settings
		// Project-specific SSH settings would be handled in project configuration
	}

	// Apply facts settings from global config
	if global.Facts != nil {
		// Note: Global Facts config has timeout and collection settings, not storage settings
		// Storage settings are typically project-specific or use global defaults
	}

	return nil
}

// applyProjectConfiguration applies project-specific configuration
func (ce *ProjectConfigurationEngine) applyProjectConfiguration(resolved *ResolvedConfiguration, project *Project) error {
	if project == nil {
		return nil
	}

	// Apply project metadata
	if project.Name != "" {
		resolved.Name = project.Name
	}
	if project.Description != "" {
		resolved.Description = project.Description
	}
	if project.Version != "" {
		resolved.Version = project.Version
	}
	if project.Environment != "" {
		resolved.Environment = project.Environment
	}
	if project.Region != "" {
		resolved.Region = project.Region
	}
	if len(project.Tags) > 0 {
		resolved.Tags = project.Tags
	}

	// Apply execution settings
	if project.Execution != nil {
		if project.Execution.DefaultTimeout != 0 {
			resolved.DefaultTimeout = project.Execution.DefaultTimeout
		}
		if project.Execution.MaxParallel != 0 {
			resolved.MaxParallel = project.Execution.MaxParallel
		}
		resolved.DryRunDefault = project.Execution.DryRunDefault
		resolved.ValidateBeforeExecute = project.Execution.ValidateBeforeExecute
		resolved.BackupBeforeChanges = project.Execution.BackupBeforeChanges
	}

	// Apply structure settings
	if project.Structure != nil {
		if project.Structure.TemplatesDir != "" {
			resolved.TemplatesDir = project.Structure.TemplatesDir
		}
		if project.Structure.DataDir != "" {
			resolved.DataDir = project.Structure.DataDir
		}
		if project.Structure.ScriptsDir != "" {
			resolved.ScriptsDir = project.Structure.ScriptsDir
		}
		if project.Structure.LogsDir != "" {
			resolved.LogsDir = project.Structure.LogsDir
		}
		if project.Structure.BackupsDir != "" {
			resolved.BackupsDir = project.Structure.BackupsDir
		}
	}

	// Apply isolation settings
	if project.Isolation != nil {
		resolved.IsolationEnabled = project.Isolation.Enabled
		if project.Isolation.FactsScope != "" {
			resolved.FactsScope = project.Isolation.FactsScope
		}
		if project.Isolation.VariablesScope != "" {
			resolved.VariablesScope = project.Isolation.VariablesScope
		}
		if project.Isolation.MachineAccess != "" {
			resolved.MachineAccess = project.Isolation.MachineAccess
		}
		if len(project.Isolation.AllowedMachines) > 0 {
			resolved.AllowedMachines = project.Isolation.AllowedMachines
		}
		if len(project.Isolation.AllowedTags) > 0 {
			resolved.AllowedTags = project.Isolation.AllowedTags
		}
	}

	return nil
}

// applyEnvironmentOverrides applies environment variable overrides
func (ce *ProjectConfigurationEngine) applyEnvironmentOverrides(resolved *ResolvedConfiguration, envVars map[string]string) {
	// Project settings
	if env := envVars["SPOOKY_PROJECT_ENVIRONMENT"]; env != "" {
		resolved.Environment = env
	}
	if region := envVars["SPOOKY_PROJECT_REGION"]; region != "" {
		resolved.Region = region
	}
	if tags := envVars["SPOOKY_PROJECT_TAGS"]; tags != "" {
		resolved.Tags = strings.Split(tags, ",")
	}

	// Execution settings
	if timeout := envVars["SPOOKY_DEFAULT_TIMEOUT"]; timeout != "" {
		if t, err := strconv.Atoi(timeout); err == nil {
			resolved.DefaultTimeout = t
		}
	}
	if parallel := envVars["SPOOKY_MAX_PARALLEL"]; parallel != "" {
		if p, err := strconv.Atoi(parallel); err == nil {
			resolved.MaxParallel = p
		}
	}
	if dryRun := envVars["SPOOKY_DRY_RUN"]; dryRun != "" {
		resolved.DryRunDefault = dryRun == "true" || dryRun == "1"
	}

	// SSH settings
	if sshTimeout := envVars["SPOOKY_SSH_TIMEOUT"]; sshTimeout != "" {
		if t, err := strconv.Atoi(sshTimeout); err == nil {
			resolved.SSHConnectionTimeout = t
		}
	}
	if sshUser := envVars["SPOOKY_SSH_USER"]; sshUser != "" {
		resolved.SSHDefaultUser = sshUser
	}
	if sshPort := envVars["SPOOKY_SSH_PORT"]; sshPort != "" {
		if p, err := strconv.Atoi(sshPort); err == nil {
			resolved.SSHDefaultPort = p
		}
	}
	if sshKeyPath := envVars["SPOOKY_SSH_KEY_PATH"]; sshKeyPath != "" {
		resolved.SSHKeyPath = sshKeyPath
	}

	// Facts settings
	if factsPath := envVars["SPOOKY_FACTS_PATH"]; factsPath != "" {
		resolved.FactsStoragePath = factsPath
	}
	if factsFormat := envVars["SPOOKY_FACTS_FORMAT"]; factsFormat != "" {
		resolved.FactsStorageType = factsFormat
	}
}

// validateResolvedConfiguration validates the resolved configuration
func (ce *ProjectConfigurationEngine) validateResolvedConfiguration(resolved *ResolvedConfiguration) error {
	// Validate execution settings
	if resolved.DefaultTimeout < 1 {
		return fmt.Errorf("default_timeout must be 1 or greater")
	}
	if resolved.MaxParallel < 2 {
		return fmt.Errorf("max_parallel must be 2 or greater")
	}

	// Validate SSH settings
	if resolved.SSHDefaultPort < 1 || resolved.SSHDefaultPort > 65535 {
		return fmt.Errorf("SSH port must be between 1 and 65535")
	}
	if resolved.SSHConnectionTimeout < 1 {
		return fmt.Errorf("SSH connection timeout must be 1 or greater")
	}
	if resolved.SSHCommandTimeout < 1 {
		return fmt.Errorf("SSH command timeout must be 1 or greater")
	}

	// Validate facts storage type
	validFactsTypes := []string{"badgerdb", "json", "hcl"}
	valid := false
	for _, t := range validFactsTypes {
		if resolved.FactsStorageType == t {
			valid = true
			break
		}
	}
	if !valid {
		return fmt.Errorf("invalid facts storage type: %s", resolved.FactsStorageType)
	}

	return nil
}
