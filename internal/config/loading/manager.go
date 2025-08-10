package loading

import (
	"fmt"
	"os"
	"path/filepath"

	spookyconfigtypes "spooky/internal/types/config"
	spookytypeslogging "spooky/internal/types/logging"
)

// Manager implements LoadingManager interface
type Manager struct {
	config     *spookyconfigtypes.LoadingConfig
	parser     spookytypesconfig.ConfigParser
	xdgManager spookytypesconfig.XDGManager
	logger     spookytypeslogging.Logger
}

// NewManager creates a new loading manager
func NewManager(config *spookyconfigtypes.LoadingConfig, logger spookytypeslogging.Logger) *Manager {
	return &Manager{
		config:     config,
		parser:     Newspookytypesconfig.ConfigParser(),
		xdgManager: Newspookytypesconfig.XDGManager(),
		logger:     logger,
	}
}

// LoadGlobalConfig loads global configuration
func (m *Manager) LoadGlobalConfig() (*spookyconfigtypes.GlobalConfig, error) {
	// 1. Get config path
	configPath := m.GetConfigPath()

	// 2. Validate config path
	if err := m.ValidateConfigPath(configPath); err != nil {
		return nil, fmt.Errorf("config path validation failed: %w", err)
	}

	// 3. Load and parse config
	content, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	config, err := m.parser.ParseHCL(content)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return config.(*spookyconfigtypes.GlobalConfig), nil
}

// LoadProjectConfig loads project configuration
func (m *Manager) LoadProjectConfig(projectPath string) (*spookyconfigtypes.ProjectConfig, error) {
	configPath := filepath.Join(projectPath, "project.hcl")

	// 1. Validate config path
	if err := m.ValidateConfigPath(configPath); err != nil {
		return nil, fmt.Errorf("config path validation failed: %w", err)
	}

	// 2. Load and parse config
	content, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	config, err := m.parser.ParseHCL(content)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return config.(*spookyconfigtypes.ProjectConfig), nil
}

// LoadFromFile loads configuration from a file
func (m *Manager) LoadFromFile(path string) (interface{}, error) {
	// 1. Validate file path
	if err := m.ValidateConfigPath(path); err != nil {
		return nil, fmt.Errorf("file path validation failed: %w", err)
	}

	// 2. Read file content
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// 3. Determine format and parse
	ext := filepath.Ext(path)
	switch ext {
	case ".hcl":
		return m.parser.ParseHCL(content)
	case ".json":
		return m.parser.ParseJSON(content)
	default:
		return nil, fmt.Errorf("unsupported file format: %s", ext)
	}
}

// LoadFromEnvironment loads environment variables
func (m *Manager) LoadFromEnvironment() (map[string]interface{}, error) {
	envVars := make(map[string]interface{})

	for _, env := range os.Environ() {
		pair := filepath.SplitList(env)
		if len(pair) == 2 {
			envVars[pair[0]] = pair[1]
		}
	}

	return envVars, nil
}

// SetConfigPath sets the configuration path
func (m *Manager) SetConfigPath(path string) error {
	if m.config == nil {
		m.config = &spookyconfigtypes.LoadingConfig{}
	}
	m.config.ConfigPath = path
	return nil
}

// SetDefaultConfig sets the default configuration
func (m *Manager) SetDefaultConfig(defaultConfig *spookyconfigtypes.GlobalConfig) error {
	if m.config == nil {
		m.config = &spookyconfigtypes.LoadingConfig{}
	}
	m.config.DefaultConfig = defaultConfig
	return nil
}

// EnableAutoReload enables or disables auto-reload
func (m *Manager) EnableAutoReload(enabled bool) error {
	if m.config == nil {
		m.config = &spookyconfigtypes.LoadingConfig{}
	}
	m.config.AutoReload = enabled
	return nil
}

// GetConfigPath returns the current configuration path
func (m *Manager) GetConfigPath() string {
	if m.config != nil && m.config.ConfigPath != "" {
		return m.config.ConfigPath
	}
	return m.xdgManager.GetConfigHome() + "/spooky/spooky.hcl"
}

// ValidateConfigPath validates a configuration path
func (m *Manager) ValidateConfigPath(path string) error {
	if path == "" {
		return fmt.Errorf("config path is empty")
	}

	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("config path does not exist: %w", err)
	}

	if info.IsDir() {
		return fmt.Errorf("config path is a directory, expected a file")
	}

	return nil
}

// Close closes the loading manager
func (m *Manager) Close() error {
	// Cleanup resources if needed
	return nil
}
