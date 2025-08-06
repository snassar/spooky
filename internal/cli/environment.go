package cli

import (
	"fmt"
	"os"

	"spooky/internal/logging"
)

// InitializeSpookyEnvironment initializes the Spooky environment
func InitializeSpookyEnvironment() error {
	logger := logging.GetLogger()

	// Create necessary directories
	dirs := []string{
		"$HOME/.config/spooky",
		"$HOME/.local/share/spooky",
		"$HOME/.cache/spooky",
	}

	for _, dir := range dirs {
		expandedDir := os.ExpandEnv(dir)
		if err := os.MkdirAll(expandedDir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", expandedDir, err)
		}
		logger.Info("Created directory", logging.String("path", expandedDir))
	}

	// Create default config file if it doesn't exist
	configPath := os.ExpandEnv("$HOME/.config/spooky/spooky.hcl")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		if err := createDefaultConfig(configPath); err != nil {
			return fmt.Errorf("failed to create default config: %w", err)
		}
		logger.Info("Created default config file", logging.String("path", configPath))
	}

	return nil
}

func createDefaultConfig(configPath string) error {
	// Create minimal default config
	defaultConfig := `# Spooky CLI Configuration

# Enable features
enable_completion = true
enable_help = true

# Commands configuration
commands {
  auto_initialize = true
  validate_commands = true
}

# Completion configuration
completion {
  enabled_shells = ["bash", "zsh", "fish"]
}

# Help configuration
help {
  enable_examples = true
  enable_usage = true
}
`
	return os.WriteFile(configPath, []byte(defaultConfig), 0644)
}
