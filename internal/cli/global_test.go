package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestGetGlobalConfig(t *testing.T) {
	// Test with default values
	config := GetGlobalConfig()

	assert.Equal(t, configDir, config.ConfigDir)
	assert.Equal(t, sshKeyPath, config.SSHKeyPath)
	assert.Equal(t, logLevel, config.LogLevel)
	assert.Equal(t, logFile, config.LogFile)
	assert.Equal(t, dryRun, config.DryRun)
	assert.Equal(t, verbose, config.Verbose)
	assert.Equal(t, quiet, config.Quiet)
}

func TestAddGlobalFlags(t *testing.T) {
	// Create a test command
	rootCmd := &cobra.Command{
		Use:   "test",
		Short: "Test command",
	}

	// Add global flags
	AddGlobalFlags(rootCmd)

	// Verify flags were added
	assert.NotNil(t, rootCmd.PersistentFlags().Lookup("config-dir"))
	assert.NotNil(t, rootCmd.PersistentFlags().Lookup("ssh-key-path"))
	assert.NotNil(t, rootCmd.PersistentFlags().Lookup("log-level"))
	assert.NotNil(t, rootCmd.PersistentFlags().Lookup("log-file"))
	assert.NotNil(t, rootCmd.PersistentFlags().Lookup("dry-run"))
	assert.NotNil(t, rootCmd.PersistentFlags().Lookup("verbose"))
	assert.NotNil(t, rootCmd.PersistentFlags().Lookup("quiet"))
}

func TestGetEnvOrDefault(t *testing.T) {
	// Test with environment variable set
	testKey := "TEST_ENV_VAR"
	testValue := "test_value"
	os.Setenv(testKey, testValue)
	defer os.Unsetenv(testKey)

	result := getEnvOrDefault(testKey, "default")
	assert.Equal(t, testValue, result)

	// Test with environment variable not set
	result = getEnvOrDefault("NONEXISTENT_VAR", "default_value")
	assert.Equal(t, "default_value", result)
}

func TestGetDefaultLogFile(t *testing.T) {
	// Test with XDG_STATE_HOME set
	xdgStateHome := "/tmp/xdg_state_home"
	os.Setenv("XDG_STATE_HOME", xdgStateHome)
	defer os.Unsetenv("XDG_STATE_HOME")

	expectedPath := filepath.Join(xdgStateHome, "spooky", "logs", "spooky.log")
	result := getDefaultLogFile()
	assert.Equal(t, expectedPath, result)

	// Test without XDG_STATE_HOME (should fall back to home directory)
	os.Unsetenv("XDG_STATE_HOME")

	// Note: Cannot easily mock os.UserHomeDir in tests
	// The function will use the actual home directory
	result = getDefaultLogFile()
	assert.NotEmpty(t, result)
}

func TestGlobalFlagsWithEnvironmentVariables(t *testing.T) {
	// Set environment variables
	testConfigDir := "/tmp/test_config"
	testSSHKeyPath := "/tmp/test_ssh"
	testLogLevel := "debug"
	testLogFile := "/tmp/test.log"

	os.Setenv("SPOOKY_CONFIG_DIR", testConfigDir)
	os.Setenv("SPOOKY_SSH_KEY_PATH", testSSHKeyPath)
	os.Setenv("SPOOKY_LOG_LEVEL", testLogLevel)
	os.Setenv("SPOOKY_LOG_FILE", testLogFile)

	defer func() {
		os.Unsetenv("SPOOKY_CONFIG_DIR")
		os.Unsetenv("SPOOKY_SSH_KEY_PATH")
		os.Unsetenv("SPOOKY_LOG_LEVEL")
		os.Unsetenv("SPOOKY_LOG_FILE")
	}()

	// Create a test command
	rootCmd := &cobra.Command{
		Use:   "test",
		Short: "Test command",
	}

	// Add global flags
	AddGlobalFlags(rootCmd)

	// Verify default values are set from environment
	configDirFlag := rootCmd.PersistentFlags().Lookup("config-dir")
	assert.Equal(t, testConfigDir, configDirFlag.DefValue)

	sshKeyPathFlag := rootCmd.PersistentFlags().Lookup("ssh-key-path")
	assert.Equal(t, testSSHKeyPath, sshKeyPathFlag.DefValue)

	logLevelFlag := rootCmd.PersistentFlags().Lookup("log-level")
	assert.Equal(t, testLogLevel, logLevelFlag.DefValue)

	logFileFlag := rootCmd.PersistentFlags().Lookup("log-file")
	assert.Equal(t, testLogFile, logFileFlag.DefValue)
}

func TestGlobalFlagsWithoutEnvironmentVariables(t *testing.T) {
	// Unset environment variables to test defaults
	os.Unsetenv("SPOOKY_CONFIG_DIR")
	os.Unsetenv("SPOOKY_SSH_KEY_PATH")
	os.Unsetenv("SPOOKY_LOG_LEVEL")
	os.Unsetenv("SPOOKY_LOG_FILE")

	// Note: Cannot easily mock os.UserHomeDir in tests
	// The function will use the actual home directory

	// Create a test command
	rootCmd := &cobra.Command{
		Use:   "test",
		Short: "Test command",
	}

	// Add global flags
	AddGlobalFlags(rootCmd)

	// Verify default values
	configDirFlag := rootCmd.PersistentFlags().Lookup("config-dir")
	assert.Equal(t, ".", configDirFlag.DefValue)

	sshKeyPathFlag := rootCmd.PersistentFlags().Lookup("ssh-key-path")
	assert.Equal(t, "~/.ssh/", sshKeyPathFlag.DefValue)

	logLevelFlag := rootCmd.PersistentFlags().Lookup("log-level")
	assert.Equal(t, "error", logLevelFlag.DefValue)

	logFileFlag := rootCmd.PersistentFlags().Lookup("log-file")
	// Note: Cannot easily mock home directory, so just verify it's not empty
	assert.NotEmpty(t, logFileFlag.DefValue)
}
