package cli

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// Global configuration variables
var (
	// Global flags
	configDir  string
	sshKeyPath string
	logLevel   string
	logFile    string
	cacheDir   string
	dryRun     bool
	verbose    bool
	quiet      bool
)

// GlobalConfig holds the global configuration
type GlobalConfig struct {
	ConfigDir  string
	SSHKeyPath string
	LogLevel   string
	LogFile    string
	CacheDir   string
	DryRun     bool
	Verbose    bool
	Quiet      bool
}

// GetGlobalConfig returns the current global configuration
func GetGlobalConfig() GlobalConfig {
	return GlobalConfig{
		ConfigDir:  configDir,
		SSHKeyPath: sshKeyPath,
		LogLevel:   logLevel,
		LogFile:    logFile,
		CacheDir:   cacheDir,
		DryRun:     dryRun,
		Verbose:    verbose,
		Quiet:      quiet,
	}
}

// AddGlobalFlags adds global flags to the root command
func AddGlobalFlags(rootCmd *cobra.Command) {
	// Get default values from environment variables
	defaultConfigDir := getEnvOrDefault("SPOOKY_CONFIG_DIR", ".")
	defaultSSHKeyPath := getEnvOrDefault("SPOOKY_SSH_KEY_PATH", "~/.ssh/")
	defaultLogLevel := getEnvOrDefault("SPOOKY_LOG_LEVEL", "error") // Changed from "info" to "error"
	defaultLogFile := getEnvOrDefault("SPOOKY_LOG_FILE", getDefaultLogFile())
	defaultCacheDir := getEnvOrDefault("SPOOKY_CACHE_DIR", getDefaultCacheDir())

	// Add global flags
	rootCmd.PersistentFlags().StringVar(&configDir, "config-dir", defaultConfigDir, "Directory containing configuration files")
	rootCmd.PersistentFlags().StringVar(&sshKeyPath, "ssh-key-path", defaultSSHKeyPath, "Path to SSH private key or directory")
	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", defaultLogLevel, "Log level: debug, info, warn, error")
	rootCmd.PersistentFlags().StringVar(&logFile, "log-file", defaultLogFile, "Log file path")
	rootCmd.PersistentFlags().StringVar(&cacheDir, "cache-dir", defaultCacheDir, "Directory for cache files")
	rootCmd.PersistentFlags().BoolVar(&dryRun, "dry-run", false, "Show what would be done without making changes")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
	rootCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "Suppress all output except errors")

	// Note: config-dir is optional with a default value, so we don't mark it as required
}

// getEnvOrDefault gets an environment variable or returns a default value
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getDefaultLogFile returns the default log file path
func getDefaultLogFile() string {
	// Try to use XDG_STATE_HOME if available
	if xdgStateHome := os.Getenv("XDG_STATE_HOME"); xdgStateHome != "" {
		return filepath.Join(xdgStateHome, "spooky", "logs", "spooky.log")
	}

	// Fallback to home directory
	if homeDir, err := os.UserHomeDir(); err == nil {
		return filepath.Join(homeDir, ".local", "state", "spooky", "logs", "spooky.log")
	}

	// Final fallback
	return "./spooky.log"
}

// getDefaultCacheDir returns the default cache directory path
func getDefaultCacheDir() string {
	// Try to use XDG_CACHE_HOME if available
	if xdgCacheHome := os.Getenv("XDG_CACHE_HOME"); xdgCacheHome != "" {
		return filepath.Join(xdgCacheHome, "spooky")
	}

	// Fallback to home directory
	if homeDir, err := os.UserHomeDir(); err == nil {
		return filepath.Join(homeDir, ".cache", "spooky")
	}

	// Final fallback
	return "./.cache/spooky"
}

// getFactsDBPath returns the path to the facts database
func getFactsDBPath() string {
	// Check for environment variable first
	if factsPath := os.Getenv("SPOOKY_FACTS_PATH"); factsPath != "" {
		return factsPath
	}

	// Default to XDG_STATE_HOME or ~/.local/state/spooky/facts.db for production
	if xdgStateHome := os.Getenv("XDG_STATE_HOME"); xdgStateHome != "" {
		return filepath.Join(xdgStateHome, "spooky", "facts.db")
	}

	// Fallback to home directory
	if homeDir, err := os.UserHomeDir(); err == nil {
		return filepath.Join(homeDir, ".local", "state", "spooky", "facts.db")
	}

	// Final fallback
	return "./facts.db"
}
