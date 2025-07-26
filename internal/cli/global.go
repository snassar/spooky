package cli

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// Global configuration variables
var (
	logLevel string
	logFile  string
	verbose  bool
	quiet    bool
)

// GlobalConfig represents the global configuration
type GlobalConfig struct {
	LogLevel string
	LogFile  string
	Verbose  bool
	Quiet    bool
}

// GetGlobalConfig returns the current global configuration
func GetGlobalConfig() GlobalConfig {
	// Use default values if not set
	configLogLevel := logLevel
	if configLogLevel == "" {
		configLogLevel = getEnvOrDefault("SPOOKY_LOG_LEVEL", "error")
	}

	configLogFile := logFile
	if configLogFile == "" {
		configLogFile = getEnvOrDefault("SPOOKY_LOG_FILE", getDefaultLogFile())
	}

	return GlobalConfig{
		LogLevel: configLogLevel,
		LogFile:  configLogFile,
		Verbose:  verbose,
		Quiet:    quiet,
	}
}

// AddGlobalFlags adds global flags to the root command
func AddGlobalFlags(rootCmd *cobra.Command) {
	// Get default values from environment variables
	defaultLogLevel := getEnvOrDefault("SPOOKY_LOG_LEVEL", "error") // Changed from "info" to "error"
	defaultLogFile := getEnvOrDefault("SPOOKY_LOG_FILE", getDefaultLogFile())

	// Add global flags
	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", defaultLogLevel, "Log level: debug, info, warn, error")
	rootCmd.PersistentFlags().StringVar(&logFile, "log-file", defaultLogFile, "Log file path")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
	rootCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "Suppress all output except errors")
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
