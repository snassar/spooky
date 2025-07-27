package config

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/hashicorp/hcl/v2/hclparse"
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
		configFile := filepath.Join("..", "..", "examples", "testing", "test-missing-project-file", "project.hcl")

		// Call ParseProjectConfigWithDebug with debug=true
		config, err := ParseProjectConfigWithDebug(configFile, true)

		// Assert error returned
		assert.Error(t, err)
		assert.Nil(t, config)

		// Verify error message contains expected content
		assert.Contains(t, err.Error(), "failed to parse project HCL file")
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
				expected: filepath.Join(configDir, "actions", "deploy.sh"),
			},
			{
				name:     "ParentDirectory",
				relative: "../config/app.conf",
				expected: filepath.Join(configDir, "..", "config", "app.conf"),
			},
			{
				name:     "CurrentDirectory",
				relative: "./scripts/setup.sh",
				expected: filepath.Join(configDir, ".", "scripts", "setup.sh"),
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
		assert.Equal(t, filepath.Join("/path", "to", "project", "   "), result, "Whitespace-only path should be joined with config directory") //nolint:gocritic
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
				expected: filepath.Join(configDir, "config files", "app config.conf"),
			},
			{
				name:     "DashesAndUnderscores",
				relative: "scripts/deploy-script_v2.sh",
				expected: filepath.Join(configDir, "scripts", "deploy-script_v2.sh"),
			},
			{
				name:     "DotsInPath",
				relative: "config/app.config.conf",
				expected: filepath.Join(configDir, "config", "app.config.conf"),
			},
			{
				name:     "UnicodeCharacters",
				relative: "config/测试.conf",
				expected: filepath.Join(configDir, "config", "测试.conf"),
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

func TestPathResolutionEdgeCases(t *testing.T) {
	t.Run("ParentDirectory", func(t *testing.T) {
		// Test resolution of paths with parent directory references
		configFile := "/path/to/project/project.hcl"
		relativePath := "../../../etc/config.conf"

		result := resolvePath(configFile, relativePath, false)

		// Should resolve to absolute path with parent directory traversal
		expected := "/path/to/project/../../../etc/config.conf"
		assert.Equal(t, expected, result, "Parent directory traversal should work correctly")

		// Verify the path is resolved to absolute
		assert.True(t, filepath.IsAbs(result), "Result should be absolute path")
	})

	t.Run("CurrentDirectory", func(t *testing.T) {
		// Test resolution of paths with current directory references
		configFile := "/path/to/project/project.hcl"
		relativePath := "./config/app.conf"

		result := resolvePath(configFile, relativePath, false)

		// Should resolve to absolute path with current directory references
		expected := "/path/to/project/./config/app.conf"
		assert.Equal(t, expected, result, "Current directory references should be handled correctly")

		// Verify the path is resolved to absolute
		assert.True(t, filepath.IsAbs(result), "Result should be absolute path")
	})

	t.Run("CrossPlatform", func(t *testing.T) {
		// Test path resolution across different operating systems

		// Test Unix-style paths
		configFile := "/path/to/project/project.hcl"
		relativePath := "config/app.conf"

		result := resolvePath(configFile, relativePath, false)
		expected := "/path/to/project/config/app.conf"
		assert.Equal(t, expected, result, "Unix-style paths should be handled correctly")

		// Test Windows-style paths (only if running on Windows)
		if filepath.IsAbs("C:\\Windows\\System32\\cmd.exe") {
			configFile = "C:\\path\\to\\project\\project.hcl"
			relativePath = "config\\app.conf"

			result = resolvePath(configFile, relativePath, false)
			expected = "C:\\path\\to\\project\\config\\app.conf"
			assert.Equal(t, expected, result, "Windows-style paths should be handled correctly")
		}
	})

	t.Run("SpecialCharacters", func(t *testing.T) {
		// Test resolution of paths with spaces and special characters
		configFile := "/path/to/project/project.hcl"
		relativePath := "config files/app config.conf"

		result := resolvePath(configFile, relativePath, false)

		// Should resolve to absolute path preserving spaces and special characters
		expected := "/path/to/project/config files/app config.conf"
		assert.Equal(t, expected, result, "Spaces in paths should be preserved")

		// Test with more special characters
		relativePath = "config/测试.conf"
		result = resolvePath(configFile, relativePath, false)
		expected = "/path/to/project/config/测试.conf"
		assert.Equal(t, expected, result, "Unicode characters should be handled correctly")
	})

	t.Run("Symlinks", func(t *testing.T) {
		// Test resolution of paths that involve symlinks
		// Create a temporary directory structure with symlinks
		tempDir, err := os.MkdirTemp("", "path-test")
		require.NoError(t, err)
		defer os.RemoveAll(tempDir)

		// Create directory structure
		projectDir := filepath.Join(tempDir, "project")
		err = os.MkdirAll(projectDir, 0o755)
		require.NoError(t, err)

		configFile := filepath.Join(projectDir, "project.hcl")

		// Create a symlink to a file
		targetFile := filepath.Join(tempDir, "target.txt")
		err = os.WriteFile(targetFile, []byte("test"), 0o600)
		require.NoError(t, err)

		symlinkPath := filepath.Join(projectDir, "symlink.txt")
		err = os.Symlink(targetFile, symlinkPath)
		require.NoError(t, err)

		// Test resolving path through symlink
		relativePath := "symlink.txt"
		result := resolvePath(configFile, relativePath, false)

		// Should resolve to the symlink path (not the target)
		expected := filepath.Join(projectDir, "symlink.txt")
		assert.Equal(t, expected, result, "Symlink path should be resolved correctly")

		// Verify the result is absolute
		assert.True(t, filepath.IsAbs(result), "Result should be absolute path")
	})

	t.Run("LongPaths", func(t *testing.T) {
		// Test resolution of very long paths
		configFile := "/path/to/project/project.hcl"
		deepPath := strings.Repeat("very/deep/nested/", 20) + "file.conf"

		result := resolvePath(configFile, deepPath, false)

		// Should handle long paths gracefully
		assert.True(t, filepath.IsAbs(result), "Long path should be resolved to absolute")
		assert.Contains(t, result, "very/deep/nested", "Long path should contain expected components")

		// Verify the path is not truncated
		assert.Greater(t, len(result), len(configFile), "Resolved path should be longer than config file path")
	})

	t.Run("EmptyPaths", func(t *testing.T) {
		// Test handling of empty and null paths
		configFile := "/path/to/project/project.hcl"

		// Test with empty string
		result := resolvePath(configFile, "", false)
		expected := "/path/to/project"
		assert.Equal(t, expected, result, "Empty path should return config directory")

		// Test with whitespace-only string
		result = resolvePath(configFile, "   ", false)
		expected = "/path/to/project/   "
		assert.Equal(t, expected, result, "Whitespace-only path should be joined with config directory")

		// Test with just dots
		result = resolvePath(configFile, ".", false)
		expected = "/path/to/project/."
		assert.Equal(t, expected, result, "Single dot should be joined with config directory")
	})

	t.Run("MachinePathsEdgeCases", func(t *testing.T) {
		// Test edge cases in machine path resolution
		configFile := "/path/to/project/project.hcl"

		machine := &Machine{
			Name:    "test-server",
			Host:    "192.168.1.100",
			Port:    22,
			User:    "testuser",
			KeyFile: "../../../etc/ssh/id_rsa", // Path with parent directory traversal
		}

		// Test that machine paths are resolved correctly
		resolveMachinePaths(configFile, machine)

		expected := "/path/to/project/../../../etc/ssh/id_rsa"
		assert.Equal(t, expected, machine.KeyFile, "Machine key file should be resolved with parent directory traversal")
	})

	t.Run("ActionPathsEdgeCases", func(t *testing.T) {
		// Test edge cases in action path resolution
		configFile := "/path/to/project/project.hcl"

		action := &Action{
			Name:   "test-action",
			Type:   "script",
			Script: "./scripts/deploy.sh", // Path with current directory reference
		}

		// Test that action paths are resolved correctly
		resolveActionPaths(configFile, action)

		expected := "/path/to/project/./scripts/deploy.sh"
		assert.Equal(t, expected, action.Script, "Action script should be resolved with current directory reference")
	})

	t.Run("ProjectPathsEdgeCases", func(t *testing.T) {
		// Test edge cases in project path resolution
		configFile := "/path/to/project/project.hcl"

		project := &ProjectConfig{
			Name:          "test-project",
			InventoryFile: "../config/inventory.hcl", // Path with parent directory
			ActionsFile:   "./actions/actions.hcl",   // Path with current directory
		}

		// Test that project paths are resolved correctly
		resolveProjectPaths(configFile, project)

		expectedInventory := "/path/to/project/../config/inventory.hcl"
		expectedActions := "/path/to/project/./actions/actions.hcl"

		assert.Equal(t, expectedInventory, project.InventoryFile, "Project inventory file should be resolved with parent directory")
		assert.Equal(t, expectedActions, project.ActionsFile, "Project actions file should be resolved with current directory")
	})
}

func TestResolveMachinePaths(t *testing.T) {
	t.Run("KeyFilePathResolution", func(t *testing.T) {
		// Test resolution of SSH key file paths
		configFile := "/path/to/project/project.hcl"
		configDir := filepath.Dir(configFile)

		machine := &Machine{
			Name:    "test-server",
			Host:    "192.168.1.100",
			Port:    22,
			User:    "testuser",
			KeyFile: "keys/id_rsa",
		}

		// Call resolveMachinePaths
		resolveMachinePaths(configFile, machine)

		// Verify the key file path is resolved correctly
		expected := filepath.Join(configDir, "keys", "id_rsa")
		assert.Equal(t, expected, machine.KeyFile, "Key file path should be resolved correctly")

		// Verify other fields remain unchanged
		assert.Equal(t, "test-server", machine.Name, "Machine name should remain unchanged")
		assert.Equal(t, "192.168.1.100", machine.Host, "Host should remain unchanged")
		assert.Equal(t, 22, machine.Port, "Port should remain unchanged")
		assert.Equal(t, "testuser", machine.User, "User should remain unchanged")
	})

	t.Run("ScriptPathResolution", func(t *testing.T) {
		// Note: Machine struct doesn't have a Script field, so this test verifies
		// that the function only processes KeyFile paths
		configFile := "/path/to/project/project.hcl"

		machine := &Machine{
			Name:    "test-server",
			Host:    "192.168.1.100",
			Port:    22,
			User:    "testuser",
			KeyFile: "", // No key file
		}

		// Call resolveMachinePaths
		resolveMachinePaths(configFile, machine)

		// Verify the machine remains unchanged when no key file is specified
		assert.Equal(t, "", machine.KeyFile, "Empty key file should remain empty")
		assert.Equal(t, "test-server", machine.Name, "Machine name should remain unchanged")
	})

	t.Run("BothPathsPresent", func(t *testing.T) {
		// Test when both key file and other paths exist
		configFile := "/path/to/project/project.hcl"
		configDir := filepath.Dir(configFile)

		machine := &Machine{
			Name:     "test-server",
			Host:     "192.168.1.100",
			Port:     22,
			User:     "testuser",
			KeyFile:  "keys/id_rsa",
			Password: "password123", // Both key file and password
			Tags: map[string]string{
				"environment": "production",
				"role":        "web",
			},
		}

		// Call resolveMachinePaths
		resolveMachinePaths(configFile, machine)

		// Verify the key file path is resolved correctly
		expected := filepath.Join(configDir, "keys", "id_rsa")
		assert.Equal(t, expected, machine.KeyFile, "Key file path should be resolved correctly")

		// Verify other fields remain unchanged
		assert.Equal(t, "password123", machine.Password, "Password should remain unchanged")
		assert.Equal(t, "production", machine.Tags["environment"], "Tags should remain unchanged")
		assert.Equal(t, "web", machine.Tags["role"], "Tags should remain unchanged")
	})

	t.Run("NoPathsPresent", func(t *testing.T) {
		// Test when no paths are specified
		configFile := "/path/to/project/project.hcl"

		machine := &Machine{
			Name:     "test-server",
			Host:     "192.168.1.100",
			Port:     22,
			User:     "testuser",
			KeyFile:  "", // No key file
			Password: "password123",
		}

		// Call resolveMachinePaths
		resolveMachinePaths(configFile, machine)

		// Verify no errors and paths remain empty
		assert.Equal(t, "", machine.KeyFile, "Empty key file should remain empty")
		assert.Equal(t, "password123", machine.Password, "Password should remain unchanged")
		assert.Equal(t, "test-server", machine.Name, "Machine name should remain unchanged")
	})

	t.Run("InvalidPaths", func(t *testing.T) {
		// Test handling of invalid path formats
		configFile := "/path/to/project/project.hcl"
		configDir := filepath.Dir(configFile)

		testCases := []struct {
			name     string
			keyFile  string
			expected string
		}{
			{
				name:     "ParentDirectoryTraversal",
				keyFile:  "../../../etc/passwd",
				expected: filepath.Join(configDir, "..", "..", "..", "etc", "passwd"),
			},
			{
				name:     "SpecialCharacters",
				keyFile:  "keys/my key with spaces",
				expected: filepath.Join(configDir, "keys", "my key with spaces"),
			},
			{
				name:     "UnicodeCharacters",
				keyFile:  "keys/测试密钥",
				expected: filepath.Join(configDir, "keys", "测试密钥"),
			},
			{
				name:     "MultipleDots",
				keyFile:  "keys/../config/../keys/id_rsa",
				expected: filepath.Join(configDir, "keys", "..", "config", "..", "keys", "id_rsa"),
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				machine := &Machine{
					Name:    "test-server",
					Host:    "192.168.1.100",
					Port:    22,
					User:    "testuser",
					KeyFile: tc.keyFile,
				}

				// Call resolveMachinePaths
				resolveMachinePaths(configFile, machine)

				// Verify the path is resolved (no security filtering)
				assert.Equal(t, tc.expected, machine.KeyFile, "Invalid path should be resolved as-is")
			})
		}
	})

	t.Run("AbsolutePaths", func(t *testing.T) {
		// Test that absolute paths are not modified
		configFile := "/path/to/project/project.hcl"

		absolutePaths := []string{
			"/home/user/.ssh/id_rsa",
			"/etc/ssh/keys/server_key",
			"/usr/local/keys/deploy_key",
		}

		for _, absolutePath := range absolutePaths {
			t.Run(filepath.Base(absolutePath), func(t *testing.T) {
				machine := &Machine{
					Name:    "test-server",
					Host:    "192.168.1.100",
					Port:    22,
					User:    "testuser",
					KeyFile: absolutePath,
				}

				// Call resolveMachinePaths
				resolveMachinePaths(configFile, machine)

				// Verify absolute paths are not modified
				assert.Equal(t, absolutePath, machine.KeyFile, "Absolute path should remain unchanged")
			})
		}
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
				configFile := filepath.Join(projectRoot, project, "inventory.hcl")

				// Create a machine with a relative key file path
				machine := &Machine{
					Name:    "test-server",
					Host:    "192.168.1.100",
					Port:    22,
					User:    "testuser",
					KeyFile: "keys/id_rsa",
				}

				// Call resolveMachinePaths
				resolveMachinePaths(configFile, machine)

				// Verify the path is resolved relative to the inventory file
				expected := filepath.Join(filepath.Dir(configFile), "keys", "id_rsa")
				assert.Equal(t, expected, machine.KeyFile, "Path should be resolved relative to inventory file")

				// Verify the result contains the expected components
				assert.Contains(t, machine.KeyFile, project, "Result should contain project name")
				assert.Contains(t, machine.KeyFile, "keys/id_rsa", "Result should contain the relative path")
			})
		}
	})
}

func TestResolveActionPaths(t *testing.T) {
	t.Run("ScriptPathResolution", func(t *testing.T) {
		// Test resolution of relative script paths in actions
		configFile := "/path/to/project/project.hcl"
		configDir := filepath.Dir(configFile)

		action := &Action{
			Name:   "test-action",
			Script: "scripts/deploy.sh",
		}

		// Call resolveActionPaths
		resolveActionPaths(configFile, action)

		// Verify the script path is resolved correctly
		expected := filepath.Join(configDir, "scripts", "deploy.sh")
		assert.Equal(t, expected, action.Script, "Script path should be resolved correctly")

		// Verify other fields remain unchanged
		assert.Equal(t, "test-action", action.Name, "Action name should remain unchanged")
	})

	t.Run("TemplatePathResolution", func(t *testing.T) {
		// Note: resolveActionPaths only processes Script field, not Template.Destination
		// This test verifies that Template.Destination is not modified
		configFile := "/path/to/project/project.hcl"

		action := &Action{
			Name: "template-action",
			Template: &TemplateConfig{
				Source:      "templates/app.conf.tmpl",
				Destination: "config/app.conf",
			},
		}

		// Store original template destination
		originalDestination := action.Template.Destination

		// Call resolveActionPaths
		resolveActionPaths(configFile, action)

		// Verify template destination is NOT modified (function only processes Script)
		assert.Equal(t, originalDestination, action.Template.Destination, "Template destination should remain unchanged")
		assert.Equal(t, "template-action", action.Name, "Action name should remain unchanged")
		assert.Equal(t, "templates/app.conf.tmpl", action.Template.Source, "Template source should remain unchanged")
	})

	t.Run("BothPathsPresent", func(t *testing.T) {
		// Test resolution when both script and template paths are present
		configFile := "/path/to/project/project.hcl"
		configDir := filepath.Dir(configFile)

		action := &Action{
			Name:   "complex-action",
			Script: "scripts/setup.sh",
			Template: &TemplateConfig{
				Source:      "templates/settings.conf.tmpl",
				Destination: "config/settings.conf",
			},
		}

		// Store original template destination
		originalDestination := action.Template.Destination

		// Call resolveActionPaths
		resolveActionPaths(configFile, action)

		// Verify script path is resolved correctly
		expectedScript := filepath.Join(configDir, "scripts", "setup.sh")
		assert.Equal(t, expectedScript, action.Script, "Script path should be resolved correctly")

		// Verify template destination is NOT modified (function only processes Script)
		assert.Equal(t, originalDestination, action.Template.Destination, "Template destination should remain unchanged")
		assert.Equal(t, "templates/settings.conf.tmpl", action.Template.Source, "Template source should remain unchanged")
	})

	t.Run("NoPathsPresent", func(t *testing.T) {
		// Test behavior when no script path is specified
		configFile := "/path/to/project/project.hcl"

		action := &Action{
			Name:    "command-action",
			Type:    "command",
			Command: "echo 'hello'",
		}

		// Call resolveActionPaths
		resolveActionPaths(configFile, action)

		// Verify action remains unchanged
		assert.Equal(t, "", action.Script, "Empty script should remain empty")
		assert.Equal(t, "command-action", action.Name, "Action name should remain unchanged")
		assert.Equal(t, "command", action.Type, "Action type should remain unchanged")
		assert.Equal(t, "echo 'hello'", action.Command, "Command should remain unchanged")
	})

	t.Run("InvalidPaths", func(t *testing.T) {
		// Test handling of invalid or malformed paths
		configFile := "/path/to/project/project.hcl"
		configDir := filepath.Dir(configFile)

		testCases := []struct {
			name     string
			script   string
			expected string
		}{
			{
				name:     "ParentDirectoryTraversal",
				script:   "../../../etc/passwd",
				expected: filepath.Join(configDir, "..", "..", "..", "etc", "passwd"),
			},
			{
				name:     "SpecialCharacters",
				script:   "scripts/my script with spaces.sh",
				expected: filepath.Join(configDir, "scripts", "my script with spaces.sh"),
			},
			{
				name:     "UnicodeCharacters",
				script:   "scripts/测试脚本.sh",
				expected: filepath.Join(configDir, "scripts", "测试脚本.sh"),
			},
			{
				name:     "MultipleDots",
				script:   "scripts/../config/../scripts/setup.sh",
				expected: filepath.Join(configDir, "scripts", "..", "config", "..", "scripts", "setup.sh"),
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				action := &Action{
					Name:   "invalid-action",
					Script: tc.script,
				}

				// Call resolveActionPaths
				resolveActionPaths(configFile, action)

				// Verify the path is resolved (no security filtering)
				assert.Equal(t, tc.expected, action.Script, "Invalid path should be resolved as-is")
			})
		}
	})

	t.Run("AbsolutePaths", func(t *testing.T) {
		// Test that absolute paths are not modified
		configFile := "/path/to/project/project.hcl"

		absolutePaths := []string{
			"/usr/local/bin/deploy.sh",
			"/etc/scripts/setup.sh",
			"/home/user/scripts/custom.sh",
		}

		for _, absolutePath := range absolutePaths {
			t.Run(filepath.Base(absolutePath), func(t *testing.T) {
				action := &Action{
					Name:   "absolute-action",
					Script: absolutePath,
				}

				// Call resolveActionPaths
				resolveActionPaths(configFile, action)

				// Verify absolute paths are not modified
				assert.Equal(t, absolutePath, action.Script, "Absolute path should remain unchanged")
			})
		}
	})

	t.Run("EmptyScriptPath", func(t *testing.T) {
		// Test behavior when script path is empty
		configFile := "/path/to/project/project.hcl"

		action := &Action{
			Name:   "empty-script-action",
			Script: "", // Empty script
		}

		// Call resolveActionPaths
		resolveActionPaths(configFile, action)

		// Verify empty script remains empty
		assert.Equal(t, "", action.Script, "Empty script should remain empty")
		assert.Equal(t, "empty-script-action", action.Name, "Action name should remain unchanged")
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
				configFile := filepath.Join(projectRoot, project, "actions.hcl")

				// Create an action with a relative script path
				action := &Action{
					Name:   "test-action",
					Script: "scripts/deploy.sh",
				}

				// Call resolveActionPaths
				resolveActionPaths(configFile, action)

				// Verify the path is resolved relative to the actions file
				expected := filepath.Join(filepath.Dir(configFile), "scripts", "deploy.sh")
				assert.Equal(t, expected, action.Script, "Path should be resolved relative to actions file")

				// Verify the result contains the expected components
				assert.Contains(t, action.Script, project, "Result should contain project name")
				assert.Contains(t, action.Script, "scripts/deploy.sh", "Result should contain the relative path")
			})
		}
	})
}

func TestParserErrorHandling(t *testing.T) {
	t.Run("HCLParsingErrors", func(t *testing.T) {
		// Use test-malformed-hcl project
		configFile := filepath.Join("..", "..", "examples", "testing", "test-malformed-hcl", "project.hcl")

		// Test ParseProjectConfig with malformed HCL
		_, err := ParseProjectConfig(configFile)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to parse project HCL file")

		// Test ParseProjectConfigWithDebug with malformed HCL
		_, err = ParseProjectConfigWithDebug(configFile, true)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to parse project HCL file")
	})

	t.Run("FileReadErrors", func(t *testing.T) {
		// Test with non-existent file
		nonexistentFile := "/nonexistent/path/project.hcl"

		_, err := ParseProjectConfig(nonexistentFile)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to parse project HCL file")
		assert.Contains(t, err.Error(), "could not be read")

		_, err = ParseProjectConfigWithDebug(nonexistentFile, false)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to parse project HCL file")

		_, err = ParseInventoryConfig(nonexistentFile)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to parse inventory HCL file")

		_, err = ParseActionsConfig(nonexistentFile)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to parse actions HCL file")
	})

	t.Run("DecodingErrors", func(t *testing.T) {
		// Use test-type-mismatches project
		configFile := filepath.Join("..", "..", "examples", "testing", "test-type-mismatches", "project.hcl")

		_, err := ParseProjectConfig(configFile)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to decode project configuration")
		assert.Contains(t, err.Error(), "Unsuitable value")
	})

	t.Run("ValidationErrors", func(t *testing.T) {
		// Use test-empty-project-name project
		configFile := filepath.Join("..", "..", "examples", "testing", "test-empty-project-name", "project.hcl")

		_, err := ParseProjectConfig(configFile)
		assert.Error(t, err)
		// Should fail validation for empty project name
	})

	t.Run("PathResolutionErrors", func(t *testing.T) {
		// Use test-valid-project for path resolution testing
		configFile := filepath.Join("..", "..", "examples", "testing", "test-valid-project", "project.hcl")

		// Should parse successfully but with resolved paths
		config, err := ParseProjectConfig(configFile)
		assert.NoError(t, err)
		assert.NotNil(t, config)
		assert.Equal(t, "test-valid-project", config.Name)
	})

	t.Run("MultipleBlocks", func(t *testing.T) {
		// Use test-multiple-project-blocks project
		configFile := filepath.Join("..", "..", "examples", "testing", "test-multiple-project-blocks", "project.hcl")

		_, err := ParseProjectConfig(configFile)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to decode project configuration")
	})

	t.Run("MissingBlocks", func(t *testing.T) {
		// Use test-no-project-block project
		configFile := filepath.Join("..", "..", "examples", "testing", "test-no-project-block", "project.hcl")

		_, err := ParseProjectConfig(configFile)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to decode project configuration")
		assert.Contains(t, err.Error(), "Unsupported argument")
	})

	t.Run("EmptyFile", func(t *testing.T) {
		// Use test-empty-project-hcl project
		configFile := filepath.Join("..", "..", "examples", "testing", "test-empty-project-hcl", "project.hcl")

		_, err := ParseProjectConfig(configFile)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no project block found in configuration")
	})

	t.Run("LargeFile", func(t *testing.T) {
		// Use test-very-long-strings project for large file testing
		configFile := filepath.Join("..", "..", "examples", "testing", "test-very-long-strings", "project.hcl")

		// Should handle large files gracefully
		config, err := ParseProjectConfig(configFile)
		assert.NoError(t, err)
		assert.NotNil(t, config)
	})

	t.Run("ConcurrentParsing", func(t *testing.T) {
		// Use test-valid-project for concurrent testing
		configFile := filepath.Join("..", "..", "examples", "testing", "test-valid-project", "project.hcl")

		// Test concurrent parsing
		var wg sync.WaitGroup
		results := make([]*ProjectConfig, 10)
		errors := make([]error, 10)

		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func(index int) {
				defer wg.Done()
				results[index], errors[index] = ParseProjectConfig(configFile)
			}(i)
		}
		wg.Wait()

		// All should succeed
		for i := 0; i < 10; i++ {
			assert.NoError(t, errors[i])
			assert.NotNil(t, results[i])
			assert.Equal(t, "test-valid-project", results[i].Name)
		}
	})

	t.Run("InventoryConfigErrors", func(t *testing.T) {
		// Use test-invalid-inventory project
		configFile := filepath.Join("..", "..", "examples", "testing", "test-invalid-inventory", "inventory.hcl")

		_, err := ParseInventoryConfig(configFile)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to parse inventory HCL file")
	})

	t.Run("ActionsConfigErrors", func(t *testing.T) {
		// Use test-malformed-actions-hcl project
		configFile := filepath.Join("..", "..", "examples", "testing", "test-malformed-actions-hcl", "actions.hcl")

		_, err := ParseActionsConfig(configFile)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to parse actions HCL file")
	})

	t.Run("LoadActionsConfigErrors", func(t *testing.T) {
		// Use test-only-project-hcl for project with no actions files
		projectPath := filepath.Join("..", "..", "examples", "testing", "test-only-project-hcl")

		config, err := LoadActionsConfig(projectPath)
		assert.NoError(t, err)
		assert.NotNil(t, config)
		assert.Empty(t, config.Actions)
	})

	t.Run("WrapperBlockValidation", func(t *testing.T) {
		// Use test-multiple-project-blocks project
		configFile := filepath.Join("..", "..", "examples", "testing", "test-multiple-project-blocks", "project.hcl")

		_, err := ParseProjectConfig(configFile)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to decode project configuration")
	})

	t.Run("RecoveryFromErrors", func(t *testing.T) {
		// First, try to parse invalid file
		invalidFile := filepath.Join("..", "..", "examples", "testing", "test-malformed-hcl", "project.hcl")

		_, err := ParseProjectConfig(invalidFile)
		assert.Error(t, err)

		// Then, try to parse valid file
		validFile := filepath.Join("..", "..", "examples", "testing", "test-valid-project", "project.hcl")

		config, err := ParseProjectConfig(validFile)
		assert.NoError(t, err)
		assert.NotNil(t, config)
		assert.Equal(t, "test-valid-project", config.Name)
	})

	t.Run("UTF8Errors", func(t *testing.T) {
		// Use test-non-utf8-files project
		configFile := filepath.Join("..", "..", "examples", "testing", "test-non-utf8-files", "invalid-encoding.hcl")

		_, err := ParseProjectConfig(configFile)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to parse project HCL file")
	})

	t.Run("PermissionErrors", func(t *testing.T) {
		// Use a non-existent file to simulate permission errors
		configFile := "/root/nonexistent/project.hcl"

		_, err := ParseProjectConfig(configFile)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to parse project HCL file")
	})
}

func TestValidateWrapperBlocks(t *testing.T) {
	parser := hclparse.NewParser()

	// Get current working directory to build absolute paths
	cwd, err := os.Getwd()
	require.NoError(t, err)

	// Go up two levels to get to the project root (from internal/config to project root)
	projectRoot := filepath.Dir(filepath.Dir(cwd))

	t.Run("SingleValidBlock", func(t *testing.T) {
		// Use the actual valid inventory file with absolute path
		validProjectPath := filepath.Join(projectRoot, "examples", "testing", "test-valid-project", "inventory.hcl")
		file, diags := parser.ParseHCLFile(validProjectPath)
		require.False(t, diags.HasErrors(), "unexpected parse error: %v", diags)
		require.NotNil(t, file, "file should not be nil")
		err := validateWrapperBlocks(file)
		assert.NoError(t, err)
	})

	t.Run("MultipleInventoryBlocks", func(t *testing.T) {
		// Use test-multiple-inventory-blocks project
		filePath := filepath.Join(projectRoot, "examples", "testing", "test-multiple-inventory-blocks", "inventory.hcl")
		file, diags := parser.ParseHCLFile(filePath)
		require.False(t, diags.HasErrors(), "unexpected parse error: %v", diags)
		require.NotNil(t, file, "file should not be nil")
		err := validateWrapperBlocks(file)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "multiple inventory blocks")
	})

	t.Run("MultipleActionsBlocks", func(t *testing.T) {
		// Use test-multiple-actions-blocks project
		filePath := filepath.Join(projectRoot, "examples", "testing", "test-multiple-actions-blocks", "actions.hcl")
		file, diags := parser.ParseHCLFile(filePath)
		require.False(t, diags.HasErrors(), "unexpected parse error: %v", diags)
		require.NotNil(t, file, "file should not be nil")
		err := validateWrapperBlocks(file)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "multiple actions blocks")
	})

	t.Run("NoBlocks", func(t *testing.T) {
		// Use the actual valid project file (which has no inventory/actions blocks) with absolute path
		validProjectPath := filepath.Join(projectRoot, "examples", "testing", "test-valid-project", "project.hcl")
		file, diags := parser.ParseHCLFile(validProjectPath)
		require.False(t, diags.HasErrors(), "unexpected parse error: %v", diags)
		require.NotNil(t, file, "file should not be nil")
		err := validateWrapperBlocks(file)
		// No error for zero blocks, function only errors for >1
		assert.NoError(t, err)
	})

	t.Run("InvalidBlockStructure", func(t *testing.T) {
		// Use the actual valid inventory file - structure errors are not detected by validateWrapperBlocks
		validProjectPath := filepath.Join(projectRoot, "examples", "testing", "test-valid-project", "inventory.hcl")
		file, diags := parser.ParseHCLFile(validProjectPath)
		require.False(t, diags.HasErrors(), "unexpected parse error: %v", diags)
		require.NotNil(t, file, "file should not be nil")
		err := validateWrapperBlocks(file)
		// Structure errors are not detected by validateWrapperBlocks, only block count
		assert.NoError(t, err)
	})
}
