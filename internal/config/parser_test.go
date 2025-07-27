package config

import (
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

		err := os.WriteFile(configFile, []byte(configContent), 0644)
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
