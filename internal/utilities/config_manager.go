package utilities

import (
	"os"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
)

// ConfigManager manages spooky configuration files
type ConfigManager struct {
	config *PathConfig
}

// NewConfigManager creates a new configuration manager
func NewConfigManager() (*ConfigManager, error) {
	config, err := GetPathConfig("spooky")
	if err != nil {
		return nil, errors.Wrap(err, "failed to get path configuration")
	}

	// Ensure directories exist
	if err := EnsureDirectories(config); err != nil {
		return nil, errors.Wrap(err, "failed to ensure directories")
	}

	return &ConfigManager{
		config: config,
	}, nil
}

// GetConfigPath returns the path to the main configuration file
func (cm *ConfigManager) GetConfigPath() string {
	return cm.config.ConfigFile
}

// GetConfigDir returns the configuration directory
func (cm *ConfigManager) GetConfigDir() string {
	return cm.config.ConfigDir
}

// WriteConfig writes a configuration file
func (cm *ConfigManager) WriteConfig(content string) error {
	return cm.WriteConfigFile(cm.config.ConfigFile, content)
}

// WriteConfigFile writes a configuration file to a specific path
func (cm *ConfigManager) WriteConfigFile(path, content string) error {
	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return errors.Wrapf(err, "failed to create config directory: %s", dir)
	}

	// Write file
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return errors.Wrapf(err, "failed to write config file: %s", path)
	}

	return nil
}

// ReadConfig reads the main configuration file
func (cm *ConfigManager) ReadConfig() (string, error) {
	return cm.ReadConfigFile(cm.config.ConfigFile)
}

// ReadConfigFile reads a configuration file from a specific path
func (cm *ConfigManager) ReadConfigFile(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", errors.Wrapf(err, "failed to read config file: %s", path)
	}

	return string(content), nil
}

// ConfigExists checks if the main configuration file exists
func (cm *ConfigManager) ConfigExists() bool {
	return cm.ConfigFileExists(cm.config.ConfigFile)
}

// ConfigFileExists checks if a configuration file exists
func (cm *ConfigManager) ConfigFileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// CreateDefaultConfig creates a default configuration file
func (cm *ConfigManager) CreateDefaultConfig() error {
	defaultConfig := `# Spooky Configuration File
# Generated automatically on first run

# Logging configuration
logging {
  level  = "info"
  format = "json"
  output = "file"
  
  file {
    path = "${log_file}"
  }
  
  structured {
    include_timestamp = true
    include_level     = true
    include_source    = false
  }
}

# Project configuration
project {
  default_timeout = "30m"
  max_parallel    = 4
  retry_attempts  = 3
}

# SSH configuration
ssh {
  default_user     = ""
  default_port     = 22
  connection_timeout = "10s"
  keepalive_interval = "30s"
}

# Template configuration
templates {
  default_engine = "go"
  cache_enabled  = true
  cache_ttl      = "1h"
}

# Variables configuration
variables {
  auto_load = true
  env_prefix = "SPOOKY_"
}

# Facts configuration
facts {
  collection_timeout = "5m"
  cache_enabled     = true
  cache_ttl         = "1h"
}

# Actions configuration
actions {
  default_timeout = "10m"
  max_retries     = 3
  parallel_limit  = 10
}

# Machines configuration
machines {
  inventory_timeout = "30s"
  connection_pool_size = 10
  max_connections_per_host = 5
}
`

	return cm.WriteConfig(defaultConfig)
}

// GetConfigInfo returns information about the configuration
func (cm *ConfigManager) GetConfigInfo() (*ConfigInfo, error) {
	info := &ConfigInfo{
		ConfigDir:  cm.config.ConfigDir,
		ConfigFile: cm.config.ConfigFile,
		Exists:     cm.ConfigExists(),
	}

	if info.Exists {
		// Get file info
		if stat, err := os.Stat(cm.config.ConfigFile); err == nil {
			info.Size = stat.Size()
			info.ModTime = stat.ModTime().Format(time.RFC3339)
		}
	}

	return info, nil
}

// ConfigInfo represents configuration file information
type ConfigInfo struct {
	ConfigDir  string `json:"config_dir"`
	ConfigFile string `json:"config_file"`
	Exists     bool   `json:"exists"`
	Size       int64  `json:"size,omitempty"`
	ModTime    string `json:"mod_time,omitempty"`
}

// ListConfigFiles lists all configuration files in the config directory
func (cm *ConfigManager) ListConfigFiles() ([]string, error) {
	entries, err := os.ReadDir(cm.config.ConfigDir)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read config directory: %s", cm.config.ConfigDir)
	}

	var configFiles []string
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".hcl" {
			configFiles = append(configFiles, filepath.Join(cm.config.ConfigDir, entry.Name()))
		}
	}

	return configFiles, nil
}

// BackupConfig creates a backup of the current configuration
func (cm *ConfigManager) BackupConfig() (string, error) {
	if !cm.ConfigExists() {
		return "", errors.New("no configuration file to backup")
	}

	content, err := cm.ReadConfig()
	if err != nil {
		return "", err
	}

	backupPath := cm.config.ConfigFile + ".backup"
	if err := cm.WriteConfigFile(backupPath, content); err != nil {
		return "", errors.Wrap(err, "failed to create backup")
	}

	return backupPath, nil
}

// RestoreConfig restores configuration from backup
func (cm *ConfigManager) RestoreConfig() error {
	backupPath := cm.config.ConfigFile + ".backup"
	
	if !cm.ConfigFileExists(backupPath) {
		return errors.New("no backup file found")
	}

	content, err := cm.ReadConfigFile(backupPath)
	if err != nil {
		return err
	}

	return cm.WriteConfig(content)
}

// ValidateConfig validates the configuration file
func (cm *ConfigManager) ValidateConfig() error {
	if !cm.ConfigExists() {
		return errors.New("configuration file does not exist")
	}

	content, err := cm.ReadConfig()
	if err != nil {
		return err
	}

	// TODO: Add HCL validation here
	// For now, just check if it's not empty
	if len(content) == 0 {
		return errors.New("configuration file is empty")
	}

	return nil
}
