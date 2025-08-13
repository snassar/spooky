// Package config provides configuration loading and management for spooky.
package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hashicorp/hcl/v2/hclsimple"

	spookylogging "spooky/internal/logging"
	spookytypeslogging "spooky/internal/types/logging"
)

// LoggingConfigManager handles loading and managing logging configuration
type LoggingConfigManager struct {
	GlobalConfigPath string
}

// NewLoggingConfigManager creates a new logging configuration manager
func NewLoggingConfigManager() *LoggingConfigManager {
	// Determine global config path
	configHome := os.Getenv("XDG_CONFIG_HOME")
	if configHome == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			configHome = "~/.config"
		} else {
			configHome = filepath.Join(homeDir, ".config")
		}
	}

	globalConfigPath := filepath.Join(configHome, "spooky", "logging.hcl")

	return &LoggingConfigManager{
		GlobalConfigPath: globalConfigPath,
	}
}

// LoadGlobalLoggingConfig loads the global logging configuration
func (lcm *LoggingConfigManager) LoadGlobalLoggingConfig() (*spookytypeslogging.LogConfig, error) {
	// Check if global logging config exists
	if _, err := os.Stat(lcm.GlobalConfigPath); os.IsNotExist(err) {
		// Return default configuration if no global config exists
		return lcm.getDefaultLoggingConfig(), nil
	}

	// Load and parse global logging configuration
	config, err := lcm.parseLoggingConfig(lcm.GlobalConfigPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load global logging config from %s: %w", lcm.GlobalConfigPath, err)
	}

	return config, nil
}

// LoadProjectLoggingConfig loads project-specific logging configuration
func (lcm *LoggingConfigManager) LoadProjectLoggingConfig(projectPath string) (*spookytypeslogging.LogConfig, error) {
	projectConfigPath := filepath.Join(projectPath, "project.hcl")

	// Check if project config exists
	if _, err := os.Stat(projectConfigPath); os.IsNotExist(err) {
		return nil, nil // No project config, use global only
	}

	// Load project configuration and extract logging section
	config, err := lcm.parseProjectLoggingConfig(projectConfigPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load project logging config from %s: %w", projectConfigPath, err)
	}

	return config, nil
}

// MergeLoggingConfigs merges global and project-specific logging configurations
func (lcm *LoggingConfigManager) MergeLoggingConfigs(global, project *spookytypeslogging.LogConfig) *spookytypeslogging.LogConfig {
	if global == nil {
		global = lcm.getDefaultLoggingConfig()
	}

	if project == nil {
		return global
	}

	// Create merged configuration
	merged := &spookytypeslogging.LogConfig{
		Level:       global.Level,
		Format:      global.Format,
		Output:      global.Output,
		File:        global.File,
		Structured:  global.Structured,
		Performance: global.Performance,
		Filtering:   global.Filtering,
		Rotation:    global.Rotation,
	}

	// Override with project-specific settings
	if project.Level != "" {
		merged.Level = project.Level
	}
	if project.Format != "" {
		merged.Format = project.Format
	}
	if project.Output != "" {
		merged.Output = project.Output
	}
	if project.File != nil {
		merged.File = project.File
	}
	if project.Structured != nil {
		merged.Structured = project.Structured
	}
	if project.Performance != nil {
		merged.Performance = project.Performance
	}
	if project.Filtering != nil {
		merged.Filtering = project.Filtering
	}
	if project.Rotation != nil {
		merged.Rotation = project.Rotation
	}

	return merged
}

// SetupLogging configures the logging system with appropriate configuration
func (lcm *LoggingConfigManager) SetupLogging(projectPath string) (*spookylogging.LogManager, error) {
	// Load global logging configuration
	globalConfig, err := lcm.LoadGlobalLoggingConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load global logging config: %w", err)
	}

	// Load project-specific logging configuration (if any)
	var projectConfig *spookytypeslogging.LogConfig
	if projectPath != "" {
		projectConfig, err = lcm.LoadProjectLoggingConfig(projectPath)
		if err != nil {
			return nil, fmt.Errorf("failed to load project logging config: %w", err)
		}
	}

	// Merge configurations
	finalConfig := lcm.MergeLoggingConfigs(globalConfig, projectConfig)

	// Create and configure log manager
	logManager := spookylogging.NewLogManager()

	if err := logManager.Configure(finalConfig); err != nil {
		return nil, fmt.Errorf("failed to configure logging: %w", err)
	}

	return logManager, nil
}

// CreateDefaultGlobalLoggingConfig creates a default global logging configuration file
func (lcm *LoggingConfigManager) CreateDefaultGlobalLoggingConfig() error {
	// Ensure directory exists
	configDir := filepath.Dir(lcm.GlobalConfigPath)
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Create default configuration content
	content := `# Global logging configuration for spooky
# This file configures logging behavior for all spooky operations

logging {
  # Log level (debug, info, warn, error, fatal)
  level = "info"
  
  # Output format (json, text, structured)
  format = "json"
  
  # Output destination (stdout, stderr, file, null)
  output = "stderr"
  
  # File output configuration (used when output = "file")
  # file {
  #   path        = "/var/log/spooky/spooky.log"
  #   permissions = "0644"
  #   append      = true
  # }
  
  # Component-specific filtering
  filtering {
    components = {
      # "ssh"     = "debug"
      # "facts"   = "info"
      # "actions" = "warn"
    }
  }
  
  # Performance optimization
  performance {
    buffer {
      enabled        = false
      size           = 4096
      flush_interval = "1s"
    }
    
    async {
      enabled       = false
      queue_size    = 1000
      workers       = 1
      drop_when_full = false
    }
  }
  
  # Log rotation (when using file output)
  # rotation {
  #   enabled      = true
  #   max_size     = "100MB"
  #   max_age      = "30d"
  #   max_backups  = 5
  #   compress     = true
  #   local_time   = false
  # }
}
`

	// Write configuration file
	if err := os.WriteFile(lcm.GlobalConfigPath, []byte(content), 0o600); err != nil {
		return fmt.Errorf("failed to write default logging config: %w", err)
	}

	return nil
}

// parseLoggingConfig parses a logging configuration file
func (lcm *LoggingConfigManager) parseLoggingConfig(configPath string) (*spookytypeslogging.LogConfig, error) {
	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return lcm.getDefaultLoggingConfig(), nil
	}

	// Parse HCL configuration
	var config spookytypeslogging.LogConfig
	if err := hclsimple.DecodeFile(configPath, nil, &config); err != nil {
		return nil, fmt.Errorf("failed to parse logging config file %s: %w", configPath, err)
	}

	// Validate configuration
	if err := lcm.validateLoggingConfig(&config); err != nil {
		return nil, fmt.Errorf("invalid logging configuration: %w", err)
	}

	return &config, nil
}

// validateLoggingConfig validates logging configuration
func (lcm *LoggingConfigManager) validateLoggingConfig(config *spookytypeslogging.LogConfig) error {
	if config == nil {
		return fmt.Errorf("logging configuration cannot be nil")
	}

	// Validate log level
	if config.Level != "" {
		switch config.Level {
		case spookytypeslogging.LogLevelDebug,
			spookytypeslogging.LogLevelInfo,
			spookytypeslogging.LogLevelWarn,
			spookytypeslogging.LogLevelError,
			spookytypeslogging.LogLevelFatal:
			// Valid level
		default:
			return fmt.Errorf("invalid log level: %s", config.Level)
		}
	}

	// Validate format
	if config.Format != "" {
		switch config.Format {
		case "json", "text", "structured":
			// Valid format
		default:
			return fmt.Errorf("invalid log format: %s", config.Format)
		}
	}

	// Validate output
	if config.Output != "" {
		switch config.Output {
		case "stdout", "stderr", "file", "null":
			// Valid output
		default:
			return fmt.Errorf("invalid log output: %s", config.Output)
		}
	}

	// Validate file configuration if output is file
	if config.Output == "file" && config.File != nil {
		if config.File.Path == "" {
			return fmt.Errorf("file path is required when output is file")
		}
	}

	return nil
}

// parseProjectLoggingConfig parses project-specific logging configuration
func (lcm *LoggingConfigManager) parseProjectLoggingConfig(projectConfigPath string) (*spookytypeslogging.LogConfig, error) {
	// Check if project config file exists
	if _, err := os.Stat(projectConfigPath); os.IsNotExist(err) {
		return nil, nil
	}

	// Parse project configuration to extract logging section
	var projectConfig struct {
		Logging *spookytypeslogging.LogConfig `hcl:"logging,optional"`
	}

	if err := hclsimple.DecodeFile(projectConfigPath, nil, &projectConfig); err != nil {
		return nil, fmt.Errorf("failed to parse project config file %s: %w", projectConfigPath, err)
	}

	// Return logging configuration if present
	if projectConfig.Logging != nil {
		// Validate the logging configuration
		if err := lcm.validateLoggingConfig(projectConfig.Logging); err != nil {
			return nil, fmt.Errorf("invalid project logging configuration: %w", err)
		}
		return projectConfig.Logging, nil
	}

	return nil, nil
}

// getDefaultLoggingConfig returns the default logging configuration
func (lcm *LoggingConfigManager) getDefaultLoggingConfig() *spookytypeslogging.LogConfig {
	return &spookytypeslogging.LogConfig{
		Level:  spookytypeslogging.LogLevelInfo,
		Format: "json",
		Output: "stderr",
	}
}
