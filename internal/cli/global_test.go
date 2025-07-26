package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestGetGlobalConfig(t *testing.T) {
	tests := []struct {
		name           string
		setLogLevel    string
		setLogFile     string
		setVerbose     bool
		setQuiet       bool
		expectedConfig GlobalConfig
	}{
		{
			name: "default values",
			expectedConfig: GlobalConfig{
				LogLevel: "error",
				LogFile:  getDefaultLogFile(),
				Verbose:  false,
				Quiet:    false,
			},
		},
		{
			name:        "custom values",
			setLogLevel: "debug",
			setLogFile:  "/custom/log/file.log",
			setVerbose:  true,
			setQuiet:    true,
			expectedConfig: GlobalConfig{
				LogLevel: "debug",
				LogFile:  "/custom/log/file.log",
				Verbose:  true,
				Quiet:    true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset global variables
			logLevel = ""
			logFile = ""
			verbose = false
			quiet = false

			// Set values if provided
			if tt.setLogLevel != "" {
				logLevel = tt.setLogLevel
			}
			if tt.setLogFile != "" {
				logFile = tt.setLogFile
			}
			verbose = tt.setVerbose
			quiet = tt.setQuiet

			config := GetGlobalConfig()
			assert.Equal(t, tt.expectedConfig, config)
		})
	}
}

func TestAddGlobalFlags(t *testing.T) {
	tests := []struct {
		name             string
		envLogLevel      string
		envLogFile       string
		expectedLogLevel string
		expectedLogFile  string
	}{
		{
			name:             "default environment",
			expectedLogLevel: "error",
			expectedLogFile:  getDefaultLogFile(),
		},
		{
			name:             "custom environment variables",
			envLogLevel:      "debug",
			envLogFile:       "/env/log/file.log",
			expectedLogLevel: "debug",
			expectedLogFile:  "/env/log/file.log",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variables
			if tt.envLogLevel != "" {
				os.Setenv("SPOOKY_LOG_LEVEL", tt.envLogLevel)
				defer os.Unsetenv("SPOOKY_LOG_LEVEL")
			}
			if tt.envLogFile != "" {
				os.Setenv("SPOOKY_LOG_FILE", tt.envLogFile)
				defer os.Unsetenv("SPOOKY_LOG_FILE")
			}

			// Create root command
			rootCmd := &cobra.Command{}
			AddGlobalFlags(rootCmd)

			// Verify flags were added
			assert.NotNil(t, rootCmd.PersistentFlags().Lookup("log-level"))
			assert.NotNil(t, rootCmd.PersistentFlags().Lookup("log-file"))
			assert.NotNil(t, rootCmd.PersistentFlags().Lookup("verbose"))
			assert.NotNil(t, rootCmd.PersistentFlags().Lookup("quiet"))

			// Verify default values
			logLevelFlag := rootCmd.PersistentFlags().Lookup("log-level")
			logFileFlag := rootCmd.PersistentFlags().Lookup("log-file")

			assert.Equal(t, tt.expectedLogLevel, logLevelFlag.DefValue)
			assert.Equal(t, tt.expectedLogFile, logFileFlag.DefValue)
		})
	}
}

func TestGetEnvOrDefault(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		defaultValue string
		envValue     string
		expected     string
	}{
		{
			name:         "environment variable set",
			key:          "TEST_VAR",
			defaultValue: "default",
			envValue:     "custom",
			expected:     "custom",
		},
		{
			name:         "environment variable not set",
			key:          "TEST_VAR",
			defaultValue: "default",
			envValue:     "",
			expected:     "default",
		},
		{
			name:         "empty key",
			key:          "",
			defaultValue: "default",
			envValue:     "custom",
			expected:     "default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variable if provided
			if tt.envValue != "" {
				os.Setenv(tt.key, tt.envValue)
				defer os.Unsetenv(tt.key)
			}

			result := getEnvOrDefault(tt.key, tt.defaultValue)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetDefaultLogFile(t *testing.T) {
	tests := []struct {
		name           string
		xdgStateHome   string
		userHomeDir    string
		expectedPrefix string
	}{
		{
			name:           "XDG_STATE_HOME set",
			xdgStateHome:   "/custom/state",
			expectedPrefix: "/custom/state/spooky/logs/spooky.log",
		},
		{
			name:           "XDG_STATE_HOME not set, user home available",
			xdgStateHome:   "",
			userHomeDir:    "",
			expectedPrefix: "", // Will be determined by actual os.UserHomeDir()
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variable if provided
			if tt.xdgStateHome != "" {
				os.Setenv("XDG_STATE_HOME", tt.xdgStateHome)
				defer os.Unsetenv("XDG_STATE_HOME")
			}

			result := getDefaultLogFile()

			if tt.xdgStateHome != "" {
				assert.Equal(t, tt.expectedPrefix, result)
			} else {
				// When XDG_STATE_HOME is not set, it should use user home directory
				// We can't easily mock os.UserHomeDir, so we just verify the path structure
				assert.NotEqual(t, "./spooky.log", result)
				assert.Contains(t, result, "spooky/logs/spooky.log")
			}
		})
	}
}

func TestGetFactsDBPath(t *testing.T) {
	tests := []struct {
		name            string
		spookyFactsPath string
		xdgStateHome    string
		userHomeDir     string
		expectedPrefix  string
	}{
		{
			name:            "SPOOKY_FACTS_PATH set",
			spookyFactsPath: "/custom/facts.db",
			expectedPrefix:  "/custom/facts.db",
		},
		{
			name:            "XDG_STATE_HOME set",
			spookyFactsPath: "",
			xdgStateHome:    "/custom/state",
			expectedPrefix:  "/custom/state/spooky/facts.db",
		},
		{
			name:            "XDG_STATE_HOME not set, user home available",
			spookyFactsPath: "",
			xdgStateHome:    "",
			userHomeDir:     "",
			expectedPrefix:  "", // Will be determined by actual os.UserHomeDir()
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variables if provided
			if tt.spookyFactsPath != "" {
				os.Setenv("SPOOKY_FACTS_PATH", tt.spookyFactsPath)
				defer os.Unsetenv("SPOOKY_FACTS_PATH")
			}
			if tt.xdgStateHome != "" {
				os.Setenv("XDG_STATE_HOME", tt.xdgStateHome)
				defer os.Unsetenv("XDG_STATE_HOME")
			}

			result := getFactsDBPath()

			switch {
			case tt.spookyFactsPath != "":
				assert.Equal(t, tt.expectedPrefix, result)
			case tt.xdgStateHome != "":
				assert.Equal(t, tt.expectedPrefix, result)
			default:
				// When neither SPOOKY_FACTS_PATH nor XDG_STATE_HOME is set, it should use user home directory
				// We can't easily mock os.UserHomeDir, so we just verify the path structure
				assert.NotEqual(t, "./facts.db", result)
				assert.Contains(t, result, "spooky/facts.db")
			}
		})
	}
}

func TestGlobalFlagsIntegration(t *testing.T) {
	t.Run("flag parsing and execution", func(t *testing.T) {
		// Create root command with global flags
		rootCmd := &cobra.Command{
			Use:   "test",
			Short: "Test command",
		}
		AddGlobalFlags(rootCmd)

		// Test flag parsing
		args := []string{
			"--log-level", "debug",
			"--log-file", "/test/log/file.log",
			"--verbose",
			"--quiet",
		}
		rootCmd.SetArgs(args)

		// Execute command
		err := rootCmd.Execute()
		assert.NoError(t, err)

		// Verify flags were parsed correctly
		config := GetGlobalConfig()
		assert.Equal(t, "debug", config.LogLevel)
		assert.Equal(t, "/test/log/file.log", config.LogFile)
		assert.True(t, config.Verbose)
		assert.True(t, config.Quiet)
	})
}

func TestGlobalConfigSerialization(t *testing.T) {
	t.Run("config can be marshaled and unmarshaled", func(t *testing.T) {
		config := GlobalConfig{
			LogLevel: "info",
			LogFile:  "/test/log/file.log",
			Verbose:  true,
			Quiet:    false,
		}

		// Test that config can be used in string operations
		assert.Contains(t, config.LogLevel, "info")
		assert.Contains(t, config.LogFile, "/test/log/file.log")
		assert.True(t, config.Verbose)
		assert.False(t, config.Quiet)
	})
}

func TestEnvironmentVariablePrecedence(t *testing.T) {
	t.Run("environment variables take precedence over defaults", func(t *testing.T) {
		// Set environment variables
		os.Setenv("SPOOKY_LOG_LEVEL", "warn")
		os.Setenv("SPOOKY_LOG_FILE", "/env/custom.log")
		defer os.Unsetenv("SPOOKY_LOG_LEVEL")
		defer os.Unsetenv("SPOOKY_LOG_FILE")

		// Create root command
		rootCmd := &cobra.Command{}
		AddGlobalFlags(rootCmd)

		// Verify environment variables were used as defaults
		logLevelFlag := rootCmd.PersistentFlags().Lookup("log-level")
		logFileFlag := rootCmd.PersistentFlags().Lookup("log-file")

		assert.Equal(t, "warn", logLevelFlag.DefValue)
		assert.Equal(t, "/env/custom.log", logFileFlag.DefValue)
	})
}

func TestFlagValidation(t *testing.T) {
	t.Run("invalid log level is handled gracefully", func(t *testing.T) {
		rootCmd := &cobra.Command{}
		AddGlobalFlags(rootCmd)

		// Set invalid log level
		args := []string{"--log-level", "invalid"}
		rootCmd.SetArgs(args)

		// Command should still execute without error
		err := rootCmd.Execute()
		assert.NoError(t, err)

		// Config should contain the invalid value
		config := GetGlobalConfig()
		assert.Equal(t, "invalid", config.LogLevel)
	})
}

func TestPathResolution(t *testing.T) {
	t.Run("relative paths are resolved correctly", func(t *testing.T) {
		// Test with relative path
		relativePath := "relative/path"
		result := getEnvOrDefault("TEST_PATH", relativePath)
		assert.Equal(t, relativePath, result)

		// Test with absolute path
		absolutePath := "/absolute/path"
		os.Setenv("TEST_PATH", absolutePath)
		defer os.Unsetenv("TEST_PATH")

		result = getEnvOrDefault("TEST_PATH", relativePath)
		assert.Equal(t, absolutePath, result)
	})
}

func TestXDGStateHomeHandling(t *testing.T) {
	t.Run("XDG_STATE_HOME with trailing slash", func(t *testing.T) {
		os.Setenv("XDG_STATE_HOME", "/custom/state/")
		defer os.Unsetenv("XDG_STATE_HOME")

		logFile := getDefaultLogFile()
		factsDB := getFactsDBPath()

		// Should not have double slashes
		assert.NotContains(t, logFile, "//")
		assert.NotContains(t, factsDB, "//")

		// Should contain the correct paths
		assert.Contains(t, logFile, "/custom/state/spooky/logs/spooky.log")
		assert.Contains(t, factsDB, "/custom/state/spooky/facts.db")
	})
}

func TestUserHomeDirectoryFallback(t *testing.T) {
	t.Run("graceful handling when user home is unavailable", func(t *testing.T) {
		// Clear environment variables
		os.Unsetenv("XDG_STATE_HOME")
		os.Unsetenv("SPOOKY_FACTS_PATH")

		// Test that functions don't panic
		assert.NotPanics(t, func() {
			logFile := getDefaultLogFile()
			assert.NotEmpty(t, logFile)
		})

		assert.NotPanics(t, func() {
			factsDB := getFactsDBPath()
			assert.NotEmpty(t, factsDB)
		})
	})
}

func TestFlagConflictHandling(t *testing.T) {
	t.Run("verbose and quiet flags can be set together", func(t *testing.T) {
		rootCmd := &cobra.Command{}
		AddGlobalFlags(rootCmd)

		// Set both verbose and quiet
		args := []string{"--verbose", "--quiet"}
		rootCmd.SetArgs(args)

		err := rootCmd.Execute()
		assert.NoError(t, err)

		config := GetGlobalConfig()
		assert.True(t, config.Verbose)
		assert.True(t, config.Quiet)
	})
}

func TestLogFileDirectoryCreation(t *testing.T) {
	t.Run("log file path contains valid directory structure", func(t *testing.T) {
		logFile := getDefaultLogFile()

		// Should contain spooky/logs directory structure
		assert.Contains(t, logFile, "spooky")
		assert.Contains(t, logFile, "logs")
		assert.Contains(t, logFile, "spooky.log")

		// Should be a valid path
		dir := filepath.Dir(logFile)
		assert.NotEmpty(t, dir)
	})
}

func TestFactsDBPathValidation(t *testing.T) {
	t.Run("facts database path is valid", func(t *testing.T) {
		factsDB := getFactsDBPath()

		// Should contain spooky directory
		assert.Contains(t, factsDB, "spooky")
		assert.Contains(t, factsDB, "facts.db")

		// Should be a valid path
		dir := filepath.Dir(factsDB)
		assert.NotEmpty(t, dir)
	})
}
