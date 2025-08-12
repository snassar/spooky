// Package config provides configuration loading and management for spooky.
package config

import (
	"fmt"
	"os"
	"path/filepath"

	spookylogging "spooky/internal/logging"
	spookytypeslogging "spooky/internal/types/logging"
)

// LoggingConfigManager handles loading and managing logging configuration
type LoggingConfigManager struct {
	GlobalConfigPath  string
	projectConfigPath string
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
	if err := os.MkdirAll(configDir, 0755); err != nil {
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
	if err := os.WriteFile(lcm.GlobalConfigPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write default logging config: %w", err)
	}

	return nil
}

// parseLoggingConfig parses a logging configuration file
func (lcm *LoggingConfigManager) parseLoggingConfig(configPath string) (*spookytypeslogging.LogConfig, error) {
	// TODO: Implement HCL parsing for logging configuration
	// For now, return default configuration
	return lcm.getDefaultLoggingConfig(), nil
}

// parseProjectLoggingConfig parses project-specific logging configuration
func (lcm *LoggingConfigManager) parseProjectLoggingConfig(projectConfigPath string) (*spookytypeslogging.LogConfig, error) {
	// TODO: Implement HCL parsing for project configuration
	// Extract logging section from project.hcl
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
