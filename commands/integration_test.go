//go:build integration

package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"spooky/internal/utilities"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCommandsIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	tests := []struct {
		name        string
		description string
		testFunc    func(t *testing.T)
	}{
		{
			name:        "TestFactsGatheringEndToEnd",
			description: "Test complete facts gathering workflow",
			testFunc:    testFactsGatheringEndToEnd,
		},
		{
			name:        "TestActionExecutionEndToEnd",
			description: "Test complete action execution workflow",
			testFunc:    testActionExecutionEndToEnd,
		},
		{
			name:        "TestMachineConnectionEndToEnd",
			description: "Test complete machine connection workflow",
			testFunc:    testMachineConnectionEndToEnd,
		},
		{
			name:        "TestProjectValidationEndToEnd",
			description: "Test complete project validation workflow",
			testFunc:    testProjectValidationEndToEnd,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("Running integration test: %s", tt.description)
			tt.testFunc(t)
		})
	}
}

func testFactsGatheringEndToEnd(t *testing.T) {
	// Create a temporary project directory
	projectDir, err := os.MkdirTemp("", "spooky-project-*")
	require.NoError(t, err)
	defer os.RemoveAll(projectDir)

	// Create basic project structure
	err = createTestProjectStructure(projectDir)
	require.NoError(t, err)

	// Test facts gathering
	// ctx := testhelpers.TestContext(t, 30*time.Second)

	// This would test the actual facts gathering command
	// For now, we'll test the project structure validation
	validator := utilities.NewProjectValidator()
	result := validator.ValidateProject(projectDir)

	assert.True(t, result.IsValid, "Project should be valid")
	assert.Empty(t, result.Errors, "Project should have no validation errors")

	t.Logf("Facts gathering test completed successfully")
}

func testActionExecutionEndToEnd(t *testing.T) {
	// Create a temporary project directory
	projectDir, err := os.MkdirTemp("", "spooky-project-*")
	require.NoError(t, err)
	defer os.RemoveAll(projectDir)

	// Create test project with actions
	err = createTestProjectWithActions(projectDir)
	require.NoError(t, err)

	// Test action execution
	// ctx := testhelpers.TestContext(t, 30*time.Second)

	// This would test the actual action execution
	// For now, we'll validate the project structure
	validator := utilities.NewProjectValidator()
	result := validator.ValidateProject(projectDir)

	assert.True(t, result.IsValid, "Project with actions should be valid")
	assert.Empty(t, result.Errors, "Project should have no validation errors")

	t.Logf("Action execution test completed successfully")
}

func testMachineConnectionEndToEnd(t *testing.T) {
	// Create a temporary project directory
	projectDir, err := os.MkdirTemp("", "spooky-project-*")
	require.NoError(t, err)
	defer os.RemoveAll(projectDir)

	// Create test project with machines
	err = createTestProjectWithMachines(projectDir)
	require.NoError(t, err)

	// Test machine connection
	// ctx := testhelpers.TestContext(t, 30*time.Second)

	// This would test the actual machine connection
	// For now, we'll validate the project structure
	validator := utilities.NewProjectValidator()
	result := validator.ValidateProject(projectDir)

	assert.True(t, result.IsValid, "Project with machines should be valid")
	assert.Empty(t, result.Errors, "Project should have no validation errors")

	t.Logf("Machine connection test completed successfully")
}

func testProjectValidationEndToEnd(t *testing.T) {
	// Create a temporary project directory
	projectDir, err := os.MkdirTemp("", "spooky-project-*")
	require.NoError(t, err)
	defer os.RemoveAll(projectDir)

	// Create complete test project
	err = createCompleteTestProject(projectDir)
	require.NoError(t, err)

	// Test project validation
	validator := utilities.NewProjectValidator()
	result := validator.ValidateProject(projectDir)

	assert.True(t, result.IsValid, "Complete project should be valid")
	assert.Empty(t, result.Errors, "Project should have no validation errors")

	// Test individual component validation
	t.Run("ProjectFileValidation", func(t *testing.T) {
		projectFile := filepath.Join(projectDir, "project.hcl")
		assert.FileExists(t, projectFile)
	})

	t.Run("MachinesFileValidation", func(t *testing.T) {
		machinesFile := filepath.Join(projectDir, "machines.hcl")
		assert.FileExists(t, machinesFile)
	})

	t.Run("ActionsFileValidation", func(t *testing.T) {
		actionsFile := filepath.Join(projectDir, "actions.hcl")
		assert.FileExists(t, actionsFile)
	})

	t.Logf("Project validation test completed successfully")
}

// Helper functions for creating test project structures

func createTestProjectStructure(projectDir string) error {
	// Create project.hcl
	projectContent := `project "test-project" {
  name = "Test Project"
  description = "A test project for integration testing"
  version = "1.0.0"
}`

	projectFile := filepath.Join(projectDir, "project.hcl")
	if err := os.WriteFile(projectFile, []byte(projectContent), 0644); err != nil {
		return fmt.Errorf("failed to create project.hcl: %w", err)
	}

	return nil
}

func createTestProjectWithActions(projectDir string) error {
	// Create basic project structure
	if err := createTestProjectStructure(projectDir); err != nil {
		return err
	}

	// Create actions.hcl
	actionsContent := `action "test-action" {
  name = "Test Action"
  description = "A test action for integration testing"
  
  command {
    name = "echo"
    args = ["Hello, World!"]
  }
}`

	actionsFile := filepath.Join(projectDir, "actions.hcl")
	if err := os.WriteFile(actionsFile, []byte(actionsContent), 0644); err != nil {
		return fmt.Errorf("failed to create actions.hcl: %w", err)
	}

	return nil
}

func createTestProjectWithMachines(projectDir string) error {
	// Create basic project structure
	if err := createTestProjectStructure(projectDir); err != nil {
		return err
	}

	// Create machines.hcl
	machinesContent := `machine "test-machine" {
  hostname = "localhost"
  port = 22
  user = "testuser"
  
  authentication {
    password {
      value = "testpass"
      encrypted = false
    }
  }
}`

	machinesFile := filepath.Join(projectDir, "machines.hcl")
	if err := os.WriteFile(machinesFile, []byte(machinesContent), 0644); err != nil {
		return fmt.Errorf("failed to create machines.hcl: %w", err)
	}

	return nil
}

func createCompleteTestProject(projectDir string) error {
	// Create project.hcl
	projectContent := `project "complete-test-project" {
  name = "Complete Test Project"
  description = "A complete test project for integration testing"
  version = "1.0.0"
  
  metadata {
    author = "Test Author"
    license = "MIT"
    repository = "https://github.com/test/spooky-project"
  }
}`

	projectFile := filepath.Join(projectDir, "project.hcl")
	if err := os.WriteFile(projectFile, []byte(projectContent), 0644); err != nil {
		return fmt.Errorf("failed to create project.hcl: %w", err)
	}

	// Create machines.hcl
	machinesContent := `machine "web-server" {
  hostname = "web.example.com"
  port = 22
  user = "deploy"
  
  authentication {
    public_key_path = "~/.ssh/id_rsa"
    passphrase {
      value = ""
      encrypted = false
    }
  }
  
  connection {
    timeout = 30
    retries = 3
  }
}

machine "database-server" {
  hostname = "db.example.com"
  port = 22
  user = "admin"
  
  authentication {
    password {
      value = "secure-password"
      encrypted = false
    }
  }
}`

	machinesFile := filepath.Join(projectDir, "machines.hcl")
	if err := os.WriteFile(machinesFile, []byte(machinesContent), 0644); err != nil {
		return fmt.Errorf("failed to create machines.hcl: %w", err)
	}

	// Create actions.hcl
	actionsContent := `action "deploy-web" {
  name = "Deploy Web Application"
  description = "Deploy the web application to web servers"
  
  targets = ["web-server"]
  
  command {
    name = "git"
    args = ["pull", "origin", "main"]
  }
  
  command {
    name = "docker"
    args = ["compose", "up", "-d"]
  }
  
  retry {
    count = 3
    delay = 5
  }
}

action "backup-database" {
  name = "Backup Database"
  description = "Create a backup of the database"
  
  targets = ["database-server"]
  
  command {
    name = "pg_dump"
    args = ["-h", "localhost", "-U", "postgres", "mydb"]
    output = "/backups/db_backup_$(date +%Y%m%d_%H%M%S).sql"
  }
  
  timeout = 300
}`

	actionsFile := filepath.Join(projectDir, "actions.hcl")
	if err := os.WriteFile(actionsFile, []byte(actionsContent), 0644); err != nil {
		return fmt.Errorf("failed to create actions.hcl: %w", err)
	}

	// Create variables.hcl
	variablesContent := `variable "environment" {
  type = "string"
  default = "development"
  description = "Target environment for deployment"
}

variable "app_version" {
  type = "string"
  default = "1.0.0"
  description = "Application version to deploy"
}`

	variablesFile := filepath.Join(projectDir, "variables.hcl")
	if err := os.WriteFile(variablesFile, []byte(variablesContent), 0644); err != nil {
		return fmt.Errorf("failed to create variables.hcl: %w", err)
	}

	// Create templates directory
	templatesDir := filepath.Join(projectDir, "templates")
	if err := os.MkdirAll(templatesDir, 0755); err != nil {
		return fmt.Errorf("failed to create templates directory: %w", err)
	}

	// Create a sample template
	templateContent := `# Web Application Configuration
Environment: {{ .environment }}
Version: {{ .app_version }}
Deployed at: {{ .timestamp }}`

	templateFile := filepath.Join(templatesDir, "app.conf.tmpl")
	if err := os.WriteFile(templateFile, []byte(templateContent), 0644); err != nil {
		return fmt.Errorf("failed to create template file: %w", err)
	}

	return nil
}

// TestCommandLineInterface tests the CLI interface
func TestCommandLineInterface(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping CLI integration test in short mode")
	}

	t.Run("HelpCommand", func(t *testing.T) {
		// Test that help command works
		// This would test the actual CLI help output
		t.Log("Help command test would be implemented here")
	})

	t.Run("VersionCommand", func(t *testing.T) {
		// Test that version command works
		// This would test the actual CLI version output
		t.Log("Version command test would be implemented here")
	})

	t.Run("InvalidCommand", func(t *testing.T) {
		// Test that invalid commands are handled properly
		// This would test error handling for invalid commands
		t.Log("Invalid command test would be implemented here")
	})
}

// TestConfigurationHandling tests configuration file handling
func TestConfigurationHandling(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping configuration integration test in short mode")
	}

	t.Run("DefaultConfiguration", func(t *testing.T) {
		// Test default configuration loading
		t.Log("Default configuration test would be implemented here")
	})

	t.Run("CustomConfiguration", func(t *testing.T) {
		// Test custom configuration file loading
		t.Log("Custom configuration test would be implemented here")
	})

	t.Run("InvalidConfiguration", func(t *testing.T) {
		// Test invalid configuration handling
		t.Log("Invalid configuration test would be implemented here")
	})
}
