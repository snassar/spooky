package project

import (
	"os"
	"path/filepath"
	"testing"
)

func TestProjectManager_CreateProject(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "spooky-project-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create project manager
	pm := NewProjectManager()

	// Create a test project
	projectConfig := &Project{
		Name:        "test-project",
		Description: "A test project",
		Version:     "1.0.0",
		Environment: "development",
	}

	// Create project
	err = pm.CreateProject(tempDir, projectConfig)
	if err != nil {
		t.Fatalf("Failed to create project: %v", err)
	}

	// Verify project structure was created
	projectHCLPath := filepath.Join(tempDir, "project.hcl")
	if _, err := os.Stat(projectHCLPath); os.IsNotExist(err) {
		t.Errorf("project.hcl file was not created")
	}

	factsDBPath := filepath.Join(tempDir, "facts.db")
	if _, err := os.Stat(factsDBPath); os.IsNotExist(err) {
		t.Errorf("facts.db directory was not created")
	}

	readmePath := filepath.Join(tempDir, "README.md")
	if _, err := os.Stat(readmePath); os.IsNotExist(err) {
		t.Errorf("README.md file was not created")
	}
}

func TestProjectManager_ValidateProject(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "spooky-project-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create project manager
	pm := NewProjectManager()

	// Create a test project
	projectConfig := &Project{
		Name:        "test-project",
		Description: "A test project",
		Version:     "1.0.0",
		Environment: "development",
	}

	// Create project
	err = pm.CreateProject(tempDir, projectConfig)
	if err != nil {
		t.Fatalf("Failed to create project: %v", err)
	}

	// Validate project
	result := pm.ValidateProject(tempDir)
	if !result.Valid {
		t.Errorf("Project validation failed: %v", result.Errors)
	}
}

func TestProjectIdentityManager_GenerateProjectName(t *testing.T) {
	pim := NewProjectIdentityManager()

	// Test with a simple directory name
	name, err := pim.GenerateProjectName("/path/to/my-project")
	if err != nil {
		t.Fatalf("Failed to generate project name: %v", err)
	}
	if name != "my-project" {
		t.Errorf("Expected 'my-project', got '%s'", name)
	}

	// Test with a complex directory name
	name, err = pim.GenerateProjectName("/path/to/My Complex Project!")
	if err != nil {
		t.Fatalf("Failed to generate project name: %v", err)
	}
	if name != "my-complex-project" {
		t.Errorf("Expected 'my-complex-project', got '%s'", name)
	}
}

func TestProjectIdentityManager_ValidateProjectName(t *testing.T) {
	pim := NewProjectIdentityManager()

	// Test valid project name
	result := pim.ValidateProjectName("valid-project", "/tmp")
	if !result.Valid {
		t.Errorf("Valid project name failed validation: %v", result.Errors)
	}

	// Test invalid project name (starts with number)
	result = pim.ValidateProjectName("123project", "/tmp")
	if result.Valid {
		t.Error("Invalid project name (starts with number) should have failed validation")
	}

	// Test invalid project name (too long)
	longName := "a" + string(make([]byte, 100))
	result = pim.ValidateProjectName(longName, "/tmp")
	if result.Valid {
		t.Error("Invalid project name (too long) should have failed validation")
	}
}

func TestProjectStructureEngine_ValidateProjectStructure(t *testing.T) {
	pse := NewProjectStructureEngine()

	// Test with non-existent path
	result := pse.ValidateProjectStructure("/non/existent/path")
	if result.Valid {
		t.Error("Non-existent path should have failed validation")
	}

	// Test with existing directory but no project.hcl
	tempDir, err := os.MkdirTemp("", "spooky-project-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	result = pse.ValidateProjectStructure(tempDir)
	if result.Valid {
		t.Error("Directory without project.hcl should have failed validation")
	}
}
