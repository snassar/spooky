// Package config provides configuration loading and management for spooky.
package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// AutoSetupConfig ensures spooky configuration directory exists and is properly configured
// This function is called before any spooky command run (except --version and --help)
func AutoSetupConfig() error {
	// Determine OS and get appropriate config directory
	configDir, err := getConfigDirectory()
	if err != nil {
		return fmt.Errorf("failed to determine config directory: %w", err)
	}

	// Check if config directory exists
	configExists, err := configDirectoryExists(configDir)
	if err != nil {
		return fmt.Errorf("failed to check config directory: %w", err)
	}

	if !configExists {
		// Create config directory and default files
		if err := createConfigDirectory(configDir); err != nil {
			return fmt.Errorf("failed to create config directory: %w", err)
		}
	} else {
		// Check if required config files exist, create them if missing
		if err := ensureConfigFiles(configDir); err != nil {
			return fmt.Errorf("config file setup failed: %w", err)
		}

		// Validate existing config files
		if err := validateConfigFiles(configDir); err != nil {
			return fmt.Errorf("config validation failed: %w", err)
		}
	}

	return nil
}

// ensureConfigFiles ensures that required config files exist, creating them if missing
func ensureConfigFiles(configDir string) error {
	// Check if spooky.hcl exists
	spookyConfigPath := filepath.Join(configDir, "spooky.hcl")
	if _, err := os.Stat(spookyConfigPath); os.IsNotExist(err) {
		if err := createDefaultSpookyConfig(configDir); err != nil {
			return fmt.Errorf("failed to create default spooky.hcl: %w", err)
		}
	}

	// Check if logging.hcl exists
	loggingConfigPath := filepath.Join(configDir, "logging.hcl")
	if _, err := os.Stat(loggingConfigPath); os.IsNotExist(err) {
		if err := createDefaultLoggingConfig(configDir); err != nil {
			return fmt.Errorf("failed to create default logging.hcl: %w", err)
		}
	}

	return nil
}

// getConfigDirectory returns the appropriate config directory for the current OS
func getConfigDirectory() (string, error) {
	switch runtime.GOOS {
	case "linux", "freebsd", "openbsd", "netbsd", "dragonfly":
		// Use XDG Base Directory Specification
		xdgConfigHome := os.Getenv("XDG_CONFIG_HOME")
		if xdgConfigHome == "" {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return "", fmt.Errorf("failed to get user home directory: %w", err)
			}
			xdgConfigHome = filepath.Join(homeDir, ".config")
		}
		return filepath.Join(xdgConfigHome, "spooky"), nil

	case "darwin":
		// macOS: ~/Library/Application Support/spooky
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to get user home directory: %w", err)
		}
		return filepath.Join(homeDir, "Library", "Application Support", "spooky"), nil

	case "windows":
		// Windows: %APPDATA%\spooky
		appData := os.Getenv("APPDATA")
		if appData == "" {
			return "", fmt.Errorf("APPDATA environment variable not set")
		}
		return filepath.Join(appData, "spooky"), nil

	default:
		return "", fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}
}

// configDirectoryExists checks if the spooky config directory exists
func configDirectoryExists(configDir string) (bool, error) {
	info, err := os.Stat(configDir)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return info.IsDir(), nil
}

// createConfigDirectory creates the spooky config directory and default configuration files
func createConfigDirectory(configDir string) error {
	// Create the config directory
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory %s: %w", configDir, err)
	}

	// Create default spooky.hcl
	if err := createDefaultSpookyConfig(configDir); err != nil {
		return fmt.Errorf("failed to create default spooky.hcl: %w", err)
	}

	// Create default logging.hcl
	if err := createDefaultLoggingConfig(configDir); err != nil {
		return fmt.Errorf("failed to create default logging.hcl: %w", err)
	}

	return nil
}

// createDefaultSpookyConfig creates a minimal spooky.hcl file with sane defaults
func createDefaultSpookyConfig(configDir string) error {
	content := `# Spooky CLI Configuration
# This file contains global configuration for the spooky CLI tool

# Global CLI settings
cli {
  # Default timeout for operations (in seconds)
  default_timeout = 300
  
  # Maximum parallel operations
  max_parallel = 10
  
  # Default log level for CLI operations
  log_level = "info"
  
  # Enable colored output (if supported by terminal)
  colored_output = true
  
  # Show progress indicators for long-running operations
  show_progress = true
}

# SSH configuration
ssh {
  # Default SSH timeout (in seconds)
  timeout = 30
  
  # SSH connection retry attempts
  retry_attempts = 3
  
  # Delay between retry attempts (in seconds)
  retry_delay = 5
  
  # Enable SSH connection pooling
  connection_pooling = true
  
  # Maximum number of SSH connections to keep in pool
  max_connections = 20
}

# Facts collection configuration
facts {
  # Default facts collection timeout (in seconds)
  timeout = 60
  
  # Enable automatic facts collection
  auto_collect = false
  
  # Maximum parallel facts collection workers
  max_parallel = 5
  
  # Facts collection retry attempts
  retry_attempts = 3
  
  # Delay between facts collection retries (in seconds)
  retry_delay = 5
}

# Actions configuration
actions {
  # Default action timeout (in seconds)
  default_timeout = 300
  
  # Maximum parallel action runs
  max_parallel = 10
  
  # Enable dry-run mode by default
  dry_run_default = false
  
  # Validate actions before running
  validate_before_run = true
  
  # Create backups before making changes
  backup_before_changes = false
}

# Storage configuration
storage {
  # Default storage format for facts databases
  facts_format = "badgerdb"
  
  # Enable compression for storage
  compression = true
  
  # Enable encryption for sensitive data
  encryption = false
}
`

	configPath := filepath.Join(configDir, "spooky.hcl")
	return os.WriteFile(configPath, []byte(content), 0644)
}

// createDefaultLoggingConfig creates a minimal logging.hcl file with sane defaults
func createDefaultLoggingConfig(configDir string) error {
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

	configPath := filepath.Join(configDir, "logging.hcl")
	return os.WriteFile(configPath, []byte(content), 0644)
}

// validateConfigFiles validates that the existing config files are valid HCL
func validateConfigFiles(configDir string) error {
	// Check if spooky.hcl exists and is valid
	spookyConfigPath := filepath.Join(configDir, "spooky.hcl")
	if err := validateHCLFile(spookyConfigPath, "spooky.hcl"); err != nil {
		return err
	}

	// Check if logging.hcl exists and is valid
	loggingConfigPath := filepath.Join(configDir, "logging.hcl")
	if err := validateHCLFile(loggingConfigPath, "logging.hcl"); err != nil {
		return err
	}

	return nil
}

// validateHCLFile validates that a file exists and contains valid HCL
func validateHCLFile(filePath, fileName string) error {
	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("required configuration file %s does not exist at %s", fileName, filePath)
	}

	// Read file content
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", fileName, err)
	}

	// Basic HCL validation - check for common syntax errors
	if err := validateHCLSyntax(content, fileName); err != nil {
		return fmt.Errorf("invalid HCL syntax in %s: %w", fileName, err)
	}

	return nil
}

// validateHCLSyntax performs basic HCL syntax validation
func validateHCLSyntax(content []byte, fileName string) error {
	contentStr := string(content)

	// Check for balanced braces first (most critical)
	if !hasBalancedBraces(contentStr) {
		return fmt.Errorf("unbalanced braces in HCL content")
	}

	// Check for basic HCL structure
	if !strings.Contains(contentStr, "{") {
		return fmt.Errorf("missing required HCL block structure (no opening braces)")
	}

	// Check for basic assignment syntax
	lines := strings.Split(contentStr, "\n")
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Check for basic assignment or block syntax
		if !isValidHCLSyntax(line) {
			return fmt.Errorf("invalid HCL syntax at line %d: %s", i+1, line)
		}
	}

	return nil
}

// hasBalancedBraces checks if the content has balanced braces
func hasBalancedBraces(content string) bool {
	stack := 0
	for _, char := range content {
		switch char {
		case '{':
			stack++
		case '}':
			stack--
			if stack < 0 {
				return false
			}
		}
	}
	return stack == 0
}

// isValidHCLSyntax checks if a line has valid HCL syntax
func isValidHCLSyntax(line string) bool {
	// Skip comments and empty lines
	if strings.HasPrefix(line, "#") || line == "" {
		return true
	}

	// Check for assignment syntax: key = value
	if strings.Contains(line, "=") {
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			if key != "" && value != "" {
				return true
			}
		}
	}

	// Check for block syntax: block_name "name" {
	if strings.Contains(line, "{") {
		// Basic block validation
		return true
	}

	// Check for block closing
	if strings.TrimSpace(line) == "}" {
		return true
	}

	// Allow lines that are just whitespace or comments
	if strings.TrimSpace(line) == "" {
		return true
	}

	return false
}

// GetConfigDirectory returns the config directory path (for external use)
func GetConfigDirectory() (string, error) {
	return getConfigDirectory()
}

// GetLoggingConfigManager creates and returns a LoggingConfigManager instance
func GetLoggingConfigManager() *LoggingConfigManager {
	return NewLoggingConfigManager()
}
