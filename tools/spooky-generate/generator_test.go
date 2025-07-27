package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewGenerator(t *testing.T) {
	// Set test values
	actionsCount = 50
	machinesCount = 25
	outputPath = "/tmp/test"
	projectName = "test-project"

	generator := NewGenerator()

	assert.Equal(t, 50, generator.actionsCount)
	assert.Equal(t, 25, generator.machinesCount)
	assert.Equal(t, "/tmp/test", generator.outputPath)
	assert.Equal(t, "test-project", generator.projectName)
}

func TestGenerateProjectConfig(t *testing.T) {
	generator := &Generator{
		projectName: "test-project",
	}

	config, err := generator.generateProjectConfig()
	require.NoError(t, err)

	assert.Contains(t, config, `project "test-project"`)
	assert.Contains(t, config, `description = "Generated spooky project for testing"`)
	assert.Contains(t, config, `inventory_file = "inventory.hcl"`)
	assert.Contains(t, config, `actions_file = "actions.hcl"`)
	assert.Contains(t, config, `log_file = "logs/spooky.log"`)
	assert.Contains(t, config, `facts_db = ".facts.db"`)
}

func TestGenerateInventory(t *testing.T) {
	generator := &Generator{
		machinesCount: 5,
	}

	inventory, err := generator.generateInventory()
	require.NoError(t, err)

	// Check basic structure
	assert.Contains(t, inventory, "inventory {")
	assert.Contains(t, inventory, "# Generated inventory for testing")
	assert.Contains(t, inventory, "# Contains 5 machines")

	// Check machine generation
	assert.Contains(t, inventory, `machine "server-1"`)
	assert.Contains(t, inventory, `machine "server-5"`)
	assert.Contains(t, inventory, `host = "10.0.1.1"`)
	assert.Contains(t, inventory, `host = "10.0.1.5"`)

	// Check tags
	assert.Contains(t, inventory, `role = "web"`)
	assert.Contains(t, inventory, `role = "database"`)
	assert.Contains(t, inventory, `environment = "production"`)
	assert.Contains(t, inventory, `region = "us-east"`)

	// Count machines
	machineCount := strings.Count(inventory, `machine "server-`)
	assert.Equal(t, 5, machineCount)
}

func TestGenerateInventoryWithLargeCount(t *testing.T) {
	generator := &Generator{
		machinesCount: 100,
	}

	inventory, err := generator.generateInventory()
	require.NoError(t, err)

	// Count machines
	machineCount := strings.Count(inventory, `machine "server-`)
	assert.Equal(t, 100, machineCount)

	// Check IP distribution
	assert.Contains(t, inventory, `host = "10.0.1.1"`)
	assert.Contains(t, inventory, `host = "10.0.1.100"`)
}

func TestGenerateActionFiles(t *testing.T) {
	tests := []struct {
		name          string
		actionsCount  int
		expectedFiles int
	}{
		{"small", 25, 1},
		{"medium", 100, 3},
		{"large", 300, 5},
		{"very-large", 600, 8},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			generator := &Generator{
				actionsCount: tt.actionsCount,
			}

			actionFiles, err := generator.generateActionFiles()
			require.NoError(t, err)

			assert.Len(t, actionFiles, tt.expectedFiles)

			// Check file names
			expectedFiles := []string{"monitoring.hcl"}
			if tt.expectedFiles >= 3 {
				expectedFiles = append(expectedFiles, "deployment.hcl", "security.hcl")
			}
			if tt.expectedFiles >= 5 {
				expectedFiles = append(expectedFiles, "maintenance.hcl", "backup.hcl")
			}
			if tt.expectedFiles >= 8 {
				expectedFiles = append(expectedFiles, "network.hcl", "database.hcl", "services.hcl")
			}

			for _, expectedFile := range expectedFiles {
				assert.Contains(t, actionFiles, expectedFile)
			}

			// Check total actions
			totalActions := 0
			for _, content := range actionFiles {
				actionCount := strings.Count(content, `action "`)
				totalActions += actionCount
			}
			assert.Equal(t, tt.actionsCount, totalActions)
		})
	}
}

func TestActionTemplates(t *testing.T) {
	generator := &Generator{}

	// Test monitoring templates
	monitoringTemplates := generator.getMonitoringTemplates()
	assert.Len(t, monitoringTemplates, 8)
	assert.Equal(t, "check-cpu-usage", monitoringTemplates[0].name)
	assert.Equal(t, "Check CPU usage percentage", monitoringTemplates[0].description)
	assert.Contains(t, monitoringTemplates[0].command, "top -bn1")
	assert.Contains(t, monitoringTemplates[0].tags, "role=web")

	// Test deployment templates
	deploymentTemplates := generator.getDeploymentTemplates()
	assert.Len(t, deploymentTemplates, 8)
	assert.Equal(t, "deploy-frontend", deploymentTemplates[0].name)
	assert.Contains(t, deploymentTemplates[0].command, "git pull")

	// Test security templates
	securityTemplates := generator.getSecurityTemplates()
	assert.Len(t, securityTemplates, 8)
	assert.Equal(t, "audit-user-accounts", securityTemplates[0].name)
	assert.Contains(t, securityTemplates[0].command, "awk -F:")

	// Test all template categories
	categories := []struct {
		name     string
		function func() []actionTemplate
	}{
		{"monitoring", generator.getMonitoringTemplates},
		{"deployment", generator.getDeploymentTemplates},
		{"security", generator.getSecurityTemplates},
		{"maintenance", generator.getMaintenanceTemplates},
		{"backup", generator.getBackupTemplates},
		{"network", generator.getNetworkTemplates},
		{"database", generator.getDatabaseTemplates},
		{"services", generator.getServiceTemplates},
	}

	for _, category := range categories {
		t.Run(category.name, func(t *testing.T) {
			templates := category.function()
			assert.Len(t, templates, 8, "Each category should have 8 templates")

			for i, template := range templates {
				assert.NotEmpty(t, template.name, "Template %d should have a name", i)
				assert.NotEmpty(t, template.description, "Template %d should have a description", i)
				assert.NotEmpty(t, template.command, "Template %d should have a command", i)
				assert.NotEmpty(t, template.tags, "Template %d should have tags", i)
				assert.Greater(t, template.timeout, 0, "Template %d should have a positive timeout", i)
			}
		})
	}
}

func TestCreateProjectStructure(t *testing.T) {
	tempDir := t.TempDir()
	generator := &Generator{}

	err := generator.createProjectStructure(tempDir)
	require.NoError(t, err)

	// Check directories exist
	expectedDirs := []string{
		"actions",
		"files",
		"templates",
		"logs",
		".facts.db",
	}

	for _, dir := range expectedDirs {
		path := filepath.Join(tempDir, dir)
		info, err := os.Stat(path)
		require.NoError(t, err)
		assert.True(t, info.IsDir())
	}

	// Check .gitignore exists
	gitignorePath := filepath.Join(tempDir, ".gitignore")
	content, err := os.ReadFile(gitignorePath)
	require.NoError(t, err)
	assert.Contains(t, string(content), ".facts.db/")
	assert.Contains(t, string(content), "logs/")
	assert.Contains(t, string(content), "*.log")
}

func TestGetProjectPath(t *testing.T) {
	tests := []struct {
		name       string
		outputPath string
		expected   string
	}{
		{"default", "", "./spooky-project"},
		{"custom", "/tmp/my-project", "/tmp/my-project"},
		{"relative", "my-project", "my-project"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			generator := &Generator{
				outputPath: tt.outputPath,
			}

			result := generator.getProjectPath()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestWriteFile(t *testing.T) {
	tempDir := t.TempDir()
	generator := &Generator{}

	// Test writing to existing directory
	testFile := filepath.Join(tempDir, "test.txt")
	testContent := "Hello, World!"

	err := generator.writeFile(testFile, testContent)
	require.NoError(t, err)

	content, err := os.ReadFile(testFile)
	require.NoError(t, err)
	assert.Equal(t, testContent, string(content))

	// Test writing to non-existent directory
	nestedFile := filepath.Join(tempDir, "nested", "deep", "test.txt")
	nestedContent := "Nested content"

	err = generator.writeFile(nestedFile, nestedContent)
	require.NoError(t, err)

	content, err = os.ReadFile(nestedFile)
	require.NoError(t, err)
	assert.Equal(t, nestedContent, string(content))
}

func TestGenerateInventoryOnly(t *testing.T) {
	tempDir := t.TempDir()
	generator := &Generator{
		machinesCount: 10,
		outputPath:    tempDir,
		projectName:   "test-project",
	}

	err := generator.GenerateInventoryOnly()
	require.NoError(t, err)

	// Check project structure
	expectedFiles := []string{
		"project.hcl",
		"inventory.hcl",
		".gitignore",
	}

	for _, file := range expectedFiles {
		path := filepath.Join(tempDir, file)
		_, err := os.Stat(path)
		assert.NoError(t, err, "File %s should exist", file)
	}

	// Check inventory content
	inventoryPath := filepath.Join(tempDir, "inventory.hcl")
	content, err := os.ReadFile(inventoryPath)
	require.NoError(t, err)

	inventoryContent := string(content)
	assert.Contains(t, inventoryContent, "inventory {")
	machineCount := strings.Count(inventoryContent, `machine "server-`)
	assert.Equal(t, 10, machineCount)
}

func TestGenerateProject(t *testing.T) {
	tempDir := t.TempDir()
	generator := &Generator{
		actionsCount:  50,
		machinesCount: 20,
		outputPath:    tempDir,
		projectName:   "test-project",
	}

	err := generator.GenerateProject()
	require.NoError(t, err)

	// Check project structure
	expectedFiles := []string{
		"project.hcl",
		"inventory.hcl",
		".gitignore",
	}

	for _, file := range expectedFiles {
		path := filepath.Join(tempDir, file)
		_, err := os.Stat(path)
		assert.NoError(t, err, "File %s should exist", file)
	}

	// Check actions directory
	actionsDir := filepath.Join(tempDir, "actions")
	entries, err := os.ReadDir(actionsDir)
	require.NoError(t, err)
	assert.Len(t, entries, 1) // Should have 1 action file for 50 actions (threshold is 51-200 for 3 files)

	// Check action files
	totalActions := 0
	for _, entry := range entries {
		if strings.HasSuffix(entry.Name(), ".hcl") {
			content, err := os.ReadFile(filepath.Join(actionsDir, entry.Name()))
			require.NoError(t, err)
			actionCount := strings.Count(string(content), `action "`)
			totalActions += actionCount
		}
	}
	assert.Equal(t, 50, totalActions)
}

func TestActionTemplateValidation(t *testing.T) {
	generator := &Generator{}

	// Test that all templates have valid HCL syntax
	categories := []struct {
		name     string
		function func() []actionTemplate
	}{
		{"monitoring", generator.getMonitoringTemplates},
		{"deployment", generator.getDeploymentTemplates},
		{"security", generator.getSecurityTemplates},
		{"maintenance", generator.getMaintenanceTemplates},
		{"backup", generator.getBackupTemplates},
		{"network", generator.getNetworkTemplates},
		{"database", generator.getDatabaseTemplates},
		{"services", generator.getServiceTemplates},
	}

	for _, category := range categories {
		t.Run(category.name, func(t *testing.T) {
			templates := category.function()

			for i, template := range templates {
				// Test that command doesn't contain unescaped quotes
				assert.NotContains(t, template.command, `"`, "Template %d command should not contain unescaped quotes", i)

				// Test that description doesn't contain unescaped quotes
				assert.NotContains(t, template.description, `"`, "Template %d description should not contain unescaped quotes", i)

				// Test that tags are properly formatted
				for _, tag := range template.tags {
					assert.Contains(t, tag, "=", "Tag should contain '=' separator")
					parts := strings.Split(tag, "=")
					assert.Len(t, parts, 2, "Tag should have exactly one '=' separator")
					assert.NotEmpty(t, parts[0], "Tag key should not be empty")
					assert.NotEmpty(t, parts[1], "Tag value should not be empty")
				}

				// Test timeout is reasonable
				assert.Greater(t, template.timeout, 0, "Timeout should be positive")
				assert.LessOrEqual(t, template.timeout, 3600, "Timeout should be reasonable (<=3600s)")
			}
		})
	}
}
