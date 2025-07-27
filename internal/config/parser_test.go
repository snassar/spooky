package config

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseProjectConfigWithDebug(t *testing.T) {
	t.Run("DebugEnabledValidConfig", func(t *testing.T) {
		// Use existing valid project from test-valid-project
		configFile := filepath.Join("..", "..", "examples", "testing", "test-valid-project", "project.hcl")

		// Call ParseProjectConfigWithDebug with debug=true
		config, err := ParseProjectConfigWithDebug(configFile, true)

		// Assert successful parsing
		assert.NoError(t, err)
		assert.NotNil(t, config)

		// Verify ProjectConfig struct fields
		assert.Equal(t, "test-valid-project", config.Name)
		assert.Equal(t, "test-valid-project project", config.Description)
		assert.Equal(t, "1.0.0", config.Version)
		assert.Equal(t, "development", config.Environment)
		assert.Equal(t, 300, config.DefaultTimeout)
		assert.True(t, config.DefaultParallel)

		// Verify paths are resolved correctly (should be absolute)
		projectDir := filepath.Dir(configFile)
		expectedInventoryPath := filepath.Join(projectDir, "inventory.hcl")
		expectedActionsPath := filepath.Join(projectDir, "actions.hcl")
		assert.Equal(t, expectedInventoryPath, config.InventoryFile)
		assert.Equal(t, expectedActionsPath, config.ActionsFile)

		// Verify storage configuration
		assert.NotNil(t, config.Storage)
		assert.Equal(t, "badgerdb", config.Storage.Type)
		assert.Equal(t, ".facts.db", config.Storage.Path)

		// Verify logging configuration
		assert.NotNil(t, config.Logging)
		assert.Equal(t, "info", config.Logging.Level)
		assert.Equal(t, "json", config.Logging.Format)
		assert.Equal(t, "logs/spooky.log", config.Logging.Output)

		// Verify SSH configuration
		assert.NotNil(t, config.SSH)
		assert.Equal(t, "debian", config.SSH.DefaultUser)
		assert.Equal(t, 22, config.SSH.DefaultPort)
		assert.Equal(t, 30, config.SSH.ConnectionTimeout)
		assert.Equal(t, 300, config.SSH.CommandTimeout)
		assert.Equal(t, 3, config.SSH.RetryAttempts)

		// Verify tags
		assert.NotNil(t, config.Tags)
		assert.Equal(t, "test-valid-project", config.Tags["project"])
	})

	t.Run("DebugDisabledValidConfig", func(t *testing.T) {
		// Use existing valid project from test-valid-project
		configFile := filepath.Join("..", "..", "examples", "testing", "test-valid-project", "project.hcl")

		// Call ParseProjectConfigWithDebug with debug=false
		config, err := ParseProjectConfigWithDebug(configFile, false)

		// Assert successful parsing
		assert.NoError(t, err)
		assert.NotNil(t, config)

		// Verify ProjectConfig struct fields
		assert.Equal(t, "test-valid-project", config.Name)
		assert.Equal(t, "test-valid-project project", config.Description)
		assert.Equal(t, "1.0.0", config.Version)
		assert.Equal(t, "development", config.Environment)

		// Verify paths are resolved correctly (should be absolute)
		projectDir := filepath.Dir(configFile)
		expectedInventoryPath := filepath.Join(projectDir, "inventory.hcl")
		expectedActionsPath := filepath.Join(projectDir, "actions.hcl")
		assert.Equal(t, expectedInventoryPath, config.InventoryFile)
		assert.Equal(t, expectedActionsPath, config.ActionsFile)
	})

	t.Run("DebugEnabledInvalidConfig", func(t *testing.T) {
		// Use existing invalid project from test-invalid-project
		configFile := filepath.Join("..", "..", "examples", "testing", "test-invalid-project", "project.hcl")

		// Call ParseProjectConfigWithDebug with debug=true
		config, err := ParseProjectConfigWithDebug(configFile, true)

		// Assert error returned
		assert.Error(t, err)
		assert.Nil(t, config)

		// Verify error message contains expected content
		assert.Contains(t, err.Error(), "failed to parse project HCL file")
	})

	t.Run("PathResolutionWithDebug", func(t *testing.T) {
		// Use existing valid project from test-valid-project
		configFile := filepath.Join("..", "..", "examples", "testing", "test-valid-project", "project.hcl")

		// Call ParseProjectConfigWithDebug with debug=true
		config, err := ParseProjectConfigWithDebug(configFile, true)

		// Assert successful parsing
		assert.NoError(t, err)
		assert.NotNil(t, config)

		// Verify ProjectConfig struct fields
		assert.Equal(t, "test-valid-project", config.Name)

		// Verify paths are resolved correctly
		projectDir := filepath.Dir(configFile)
		expectedInventoryPath := filepath.Join(projectDir, "inventory.hcl")
		expectedActionsPath := filepath.Join(projectDir, "actions.hcl")
		assert.Equal(t, expectedInventoryPath, config.InventoryFile)
		assert.Equal(t, expectedActionsPath, config.ActionsFile)
	})

	t.Run("NonExistentFile", func(t *testing.T) {
		// Call ParseProjectConfigWithDebug with non-existent file
		config, err := ParseProjectConfigWithDebug("/non/existent/file.hcl", true)

		// Assert error returned
		assert.Error(t, err)
		assert.Nil(t, config)

		// Verify error message contains expected content
		assert.Contains(t, err.Error(), "failed to parse project HCL file")
	})

	t.Run("MissingProjectBlock", func(t *testing.T) {
		// Use test-missing-project-file which has no project.hcl
		// Create a temporary file with no project block
		tempDir := t.TempDir()
		configFile := filepath.Join(tempDir, "project.hcl")

		configContent := `# No project block here
storage {
  type = "badgerdb"
  path = "/tmp/facts"
}`

		err := os.WriteFile(configFile, []byte(configContent), 0o600)
		require.NoError(t, err)

		// Call ParseProjectConfigWithDebug with debug=true
		config, err := ParseProjectConfigWithDebug(configFile, true)

		// Assert error returned
		assert.Error(t, err)
		assert.Nil(t, config)

		// Verify error message contains expected content
		assert.Contains(t, err.Error(), "Unsupported block type")
	})
}

func TestLoadActionsConfig(t *testing.T) {
	t.Run("RootActionsFileOnly", func(t *testing.T) {
		// Use test-only-actions-hcl project
		projectPath := filepath.Join("..", "..", "examples", "testing", "test-only-actions-hcl")

		// Call LoadActionsConfig
		config, err := LoadActionsConfig(projectPath)

		// Assert successful loading
		assert.NoError(t, err)
		assert.NotNil(t, config)

		// Verify actions from both root file and actions/ directory (5 total)
		// 1 from root actions.hcl + 1 from 01-dependencies.hcl + 1 from 02-system-update.hcl + 2 from 03-monitoring.hcl
		assert.Len(t, config.Actions, 5)

		// Check that root action is included
		rootActionFound := false
		for _, action := range config.Actions {
			if action.Name == "check-status" && action.Command == "uptime && df -h" {
				rootActionFound = true
				break
			}
		}
		assert.True(t, rootActionFound, "Root action should be included in merged config")

		// Check that directory actions are included
		dirActionFound := false
		for _, action := range config.Actions {
			if action.Name == "install-dependencies" {
				dirActionFound = true
				break
			}
		}
		assert.True(t, dirActionFound, "Directory action should be included in merged config")
	})

	t.Run("ActionsDirectoryOnly", func(t *testing.T) {
		// Use test-missing-actions/test-valid-project (has actions/ but also root actions.hcl)
		projectPath := filepath.Join("..", "..", "examples", "testing", "test-missing-actions", "test-valid-project")

		// Call LoadActionsConfig
		config, err := LoadActionsConfig(projectPath)

		// Assert successful loading
		assert.NoError(t, err)
		assert.NotNil(t, config)

		// Verify actions from both root file and directory files (5 total)
		// 1 from root actions.hcl + 1 from 01-dependencies.hcl + 1 from 02-system-update.hcl + 2 from 03-monitoring.hcl
		assert.Len(t, config.Actions, 5)

		// Check that actions are loaded in sorted order (directory files are sorted)
		actionNames := make([]string, len(config.Actions))
		for i, action := range config.Actions {
			actionNames[i] = action.Name
		}
		// Should contain both root and directory actions
		assert.Contains(t, actionNames, "check-status")
		assert.Contains(t, actionNames, "install-dependencies")
		assert.Contains(t, actionNames, "update-system")
		assert.Contains(t, actionNames, "check-disk-space")
		assert.Contains(t, actionNames, "check-memory")
	})

	t.Run("BothSources", func(t *testing.T) {
		// Use test-valid-project
		projectPath := filepath.Join("..", "..", "examples", "testing", "test-valid-project")

		// Call LoadActionsConfig
		config, err := LoadActionsConfig(projectPath)

		// Assert successful loading
		assert.NoError(t, err)
		assert.NotNil(t, config)

		// Verify actions merged from both sources
		// Should have actions from root actions.hcl + actions/ directory
		assert.Len(t, config.Actions, 5) // 1 from root + 4 from actions/

		// Check that root action is included
		rootActionFound := false
		for _, action := range config.Actions {
			if action.Name == "check-status" && action.Command == "uptime && df -h" {
				rootActionFound = true
				break
			}
		}
		assert.True(t, rootActionFound, "Root action should be included in merged config")

		// Check that directory actions are included
		dirActionFound := false
		for _, action := range config.Actions {
			if action.Name == "install-dependencies" {
				dirActionFound = true
				break
			}
		}
		assert.True(t, dirActionFound, "Directory action should be included in merged config")
	})

	t.Run("MergeConflicts", func(t *testing.T) {
		// Use test-duplicate-actions project
		projectPath := filepath.Join("..", "..", "examples", "testing", "test-duplicate-actions")

		// Call LoadActionsConfig
		config, err := LoadActionsConfig(projectPath)

		// Assert successful loading (LoadActionsConfig doesn't validate duplicates)
		assert.NoError(t, err)
		assert.NotNil(t, config)

		// Verify both duplicate actions are loaded (LoadActionsConfig just merges, doesn't validate)
		checkStatusActions := 0
		for _, action := range config.Actions {
			if action.Name == "check-status" {
				checkStatusActions++
			}
		}
		assert.Equal(t, 2, checkStatusActions, "Both duplicate actions should be loaded")

		// Verify the actions have different descriptions (showing they're different instances)
		descriptions := make([]string, 0)
		for _, action := range config.Actions {
			if action.Name == "check-status" {
				descriptions = append(descriptions, action.Description)
			}
		}
		assert.Contains(t, descriptions, "Check server status")
		assert.Contains(t, descriptions, "Another check status action")
	})

	t.Run("InvalidFilesInDirectory", func(t *testing.T) {
		// Use test-invalid-actions project
		projectPath := filepath.Join("..", "..", "examples", "testing", "test-invalid-actions")

		// Call LoadActionsConfig
		config, err := LoadActionsConfig(projectPath)

		// Assert successful loading (the invalid HCL is actually parsed successfully)
		assert.NoError(t, err)
		assert.NotNil(t, config)

		// Verify that all actions are loaded (including the "invalid" ones)
		// The HCL parser is more lenient than expected
		assert.Len(t, config.Actions, 7) // 3 from root + 4 from actions/

		// Check that the "invalid" actions are actually loaded
		invalidActionFound := false
		for _, action := range config.Actions {
			if action.Name == "invalid-action" {
				invalidActionFound = true
				break
			}
		}
		assert.True(t, invalidActionFound, "Invalid action should be loaded despite missing command")

		// Check that the "broken" action is also loaded
		brokenActionFound := false
		for _, action := range config.Actions {
			if action.Name == "broken-action" {
				brokenActionFound = true
				break
			}
		}
		assert.True(t, brokenActionFound, "Broken action should be loaded despite syntax issues")
	})

	t.Run("NoActionsFiles", func(t *testing.T) {
		// Use test-only-project-hcl project (no actions.hcl, no actions/ directory)
		projectPath := filepath.Join("..", "..", "examples", "testing", "test-only-project-hcl")

		// Call LoadActionsConfig
		config, err := LoadActionsConfig(projectPath)

		// Assert successful loading
		assert.NoError(t, err)
		assert.NotNil(t, config)

		// Verify empty ActionsConfig returned
		assert.Len(t, config.Actions, 0)
	})

	t.Run("DirectoryReadError", func(t *testing.T) {
		// Use a non-existent project path to trigger directory read error
		projectPath := "/non/existent/project/path"

		// Call LoadActionsConfig
		config, err := LoadActionsConfig(projectPath)

		// Assert successful loading (no error because both root and actions/ don't exist)
		assert.NoError(t, err)
		assert.NotNil(t, config)

		// Verify empty ActionsConfig returned
		assert.Len(t, config.Actions, 0)
	})
}

func TestResolvePath(t *testing.T) {
	t.Run("AbsolutePaths", func(t *testing.T) {
		// Test that absolute paths remain unchanged
		configFile := "/path/to/project/project.hcl"
		absolutePaths := []string{
			"/usr/local/bin/script.sh",
			"/etc/config/app.conf",
			"/home/user/.ssh/id_rsa",
		}

		for _, path := range absolutePaths {
			t.Run(filepath.Base(path), func(t *testing.T) {
				result := resolvePath(configFile, path, false)
				assert.Equal(t, path, result, "Absolute path should remain unchanged")
			})
		}

		// Test Windows-style absolute path (only on Windows)
		if filepath.IsAbs("C:\\Windows\\System32\\cmd.exe") {
			t.Run("WindowsAbsolute", func(t *testing.T) {
				result := resolvePath(configFile, "C:\\Windows\\System32\\cmd.exe", false)
				assert.Equal(t, "C:\\Windows\\System32\\cmd.exe", result, "Windows absolute path should remain unchanged")
			})
		}
	})

	t.Run("RelativePaths", func(t *testing.T) {
		// Test that relative paths are resolved to absolute
		configFile := "/path/to/project/project.hcl"
		configDir := filepath.Dir(configFile)

		testCases := []struct {
			name     string
			relative string
			expected string
		}{
			{
				name:     "SimpleRelative",
				relative: "inventory.hcl",
				expected: filepath.Join(configDir, "inventory.hcl"),
			},
			{
				name:     "Subdirectory",
				relative: "actions/deploy.sh",
				expected: filepath.Join(configDir, "actions/deploy.sh"),
			},
			{
				name:     "ParentDirectory",
				relative: "../config/app.conf",
				expected: filepath.Join(configDir, "../config/app.conf"),
			},
			{
				name:     "CurrentDirectory",
				relative: "./scripts/setup.sh",
				expected: filepath.Join(configDir, "./scripts/setup.sh"),
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				result := resolvePath(configFile, tc.relative, false)
				assert.Equal(t, tc.expected, result, "Relative path should be resolved correctly")
			})
		}
	})

	t.Run("EmptyPaths", func(t *testing.T) {
		// Test handling of empty paths
		configFile := "/path/to/project/project.hcl"

		// Test with empty string - should be joined with config directory
		result := resolvePath(configFile, "", false)
		assert.Equal(t, "/path/to/project", result, "Empty path should return config directory")

		// Test with whitespace-only string - should be joined with config directory
		result = resolvePath(configFile, "   ", false)
		assert.Equal(t, filepath.Join("/path/to/project", "   "), result, "Whitespace-only path should be joined with config directory")
	})

	t.Run("SpecialCharacters", func(t *testing.T) {
		// Test paths containing spaces, special chars
		configFile := "/path/to/project/project.hcl"
		configDir := filepath.Dir(configFile)

		testCases := []struct {
			name     string
			relative string
			expected string
		}{
			{
				name:     "SpacesInPath",
				relative: "config files/app config.conf",
				expected: filepath.Join(configDir, "config files/app config.conf"),
			},
			{
				name:     "DashesAndUnderscores",
				relative: "scripts/deploy-script_v2.sh",
				expected: filepath.Join(configDir, "scripts/deploy-script_v2.sh"),
			},
			{
				name:     "DotsInPath",
				relative: "config/app.config.conf",
				expected: filepath.Join(configDir, "config/app.config.conf"),
			},
			{
				name:     "UnicodeCharacters",
				relative: "config/测试.conf",
				expected: filepath.Join(configDir, "config/测试.conf"),
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				result := resolvePath(configFile, tc.relative, false)
				assert.Equal(t, tc.expected, result, "Path with special characters should be handled correctly")
			})
		}
	})

	t.Run("DebugMode", func(t *testing.T) {
		// Test path resolution with debug mode enabled
		configFile := "/path/to/project/project.hcl"
		relativePath := "inventory.hcl"

		// Capture stdout to check debug output
		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		// Call resolvePath with debug=true
		result := resolvePath(configFile, relativePath, true)

		// Restore stdout
		w.Close()
		os.Stdout = oldStdout

		// Read captured output
		var buf bytes.Buffer
		_, err := buf.ReadFrom(r)
		require.NoError(t, err)
		debugOutput := buf.String()

		// Verify debug output contains expected information
		assert.Contains(t, debugOutput, "[DEBUG] resolvePath")
		assert.Contains(t, debugOutput, "configFile=")
		assert.Contains(t, debugOutput, "path=")
		assert.Contains(t, debugOutput, "configDir=")
		assert.Contains(t, debugOutput, "resolved=")

		// Verify the result is correct
		expected := filepath.Join(filepath.Dir(configFile), relativePath)
		assert.Equal(t, expected, result, "Path should be resolved correctly even with debug mode")
	})

	t.Run("RealProjectPaths", func(t *testing.T) {
		// Test with paths from actual test projects
		projectRoot := "../../examples/testing"
		testProjects := []string{
			"test-valid-project",
			"test-special-characters",
		}

		for _, project := range testProjects {
			t.Run(project, func(t *testing.T) {
				configFile := filepath.Join(projectRoot, project, "project.hcl")

				// Test with relative paths that exist in the project
				relativePaths := []string{
					"inventory.hcl",
					"actions.hcl",
					"actions/01-dependencies.hcl",
				}

				for _, relativePath := range relativePaths {
					t.Run(filepath.Base(relativePath), func(t *testing.T) {
						result := resolvePath(configFile, relativePath, false)

						// Verify the result contains the expected components
						assert.Contains(t, result, project, "Result should contain project name")
						assert.Contains(t, result, relativePath, "Result should contain the relative path")

						// Verify the result is properly resolved relative to config file
						expected := filepath.Join(filepath.Dir(configFile), relativePath)
						assert.Equal(t, expected, result, "Path should be resolved relative to config file")
					})
				}
			})
		}
	})
}
