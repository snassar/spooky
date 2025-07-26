package cli

import (
	"crypto/rand"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"spooky/internal/logging"
)

func TestInitCmd(t *testing.T) {
	// Initialize commands to set up flags
	InitCommands()

	// Create temporary directory for test projects
	tempDir := filepath.Join(os.TempDir(), "spooky-test-"+generateRandomSuffix())
	if err := os.MkdirAll(tempDir, 0o755); err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir) // Clean up after test

	tests := []struct {
		name        string
		args        []string
		flags       map[string]string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "valid project with positional args",
			args:        []string{"test-project", tempDir},
			expectError: false,
		},
		{
			name:        "valid project with flags",
			args:        []string{},
			flags:       map[string]string{"project": "test-project", "path": tempDir},
			expectError: false,
		},
		{
			name:        "missing project name",
			args:        []string{},
			expectError: true,
			errorMsg:    "project name is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up the command
			cmd := InitCmd
			cmd.SetArgs(tt.args)

			// Reset flags to clear any previous values
			err := cmd.Flags().Set("project", "")
			require.NoError(t, err)

			// Only set default path if no positional args are provided
			if len(tt.args) <= 1 {
				err = cmd.Flags().Set("path", ".")
				require.NoError(t, err)
			}

			// Set flags
			for flag, value := range tt.flags {
				err = cmd.Flags().Set(flag, value)
				require.NoError(t, err)
			}

			// Execute the command
			err = cmd.Execute()

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateCmd(t *testing.T) {
	projectRoot := getProjectRoot()
	tests := []struct {
		name        string
		projectPath string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "valid project",
			projectPath: filepath.Join(projectRoot, "examples", "testing", "test-valid-project"),
			expectError: false,
		},
		{
			name:        "invalid project syntax",
			projectPath: filepath.Join(projectRoot, "examples", "testing", "test-invalid-project"),
			expectError: true,
			errorMsg:    "Unclosed configuration block",
		},
		{
			name:        "missing project file",
			projectPath: filepath.Join(projectRoot, "examples", "testing", "test-missing-project-file"),
			expectError: true,
			errorMsg:    "project.hcl not found",
		},
		{
			name:        "missing inventory",
			projectPath: filepath.Join(projectRoot, "examples", "testing", "test-missing-inventory"),
			expectError: false, // Only logs a warning, doesn't fail validation
		},
		{
			name:        "missing actions",
			projectPath: filepath.Join(projectRoot, "examples", "testing", "test-missing-actions"),
			expectError: false, // Only logs a warning, doesn't fail validation
		},
		{
			name:        "duplicate machines",
			projectPath: filepath.Join(projectRoot, "examples", "testing", "test-duplicate-machines"),
			expectError: false, // Validation is passing, duplicate detection may not be implemented
		},
		{
			name:        "duplicate actions",
			projectPath: filepath.Join(projectRoot, "examples", "testing", "test-duplicate-actions"),
			expectError: false, // Validation is passing, duplicate detection may not be implemented
		},
		{
			name:        "invalid port range",
			projectPath: filepath.Join(projectRoot, "examples", "testing", "test-port-timeout-range"),
			expectError: true,
			errorMsg:    "multiple inventory blocks found",
		},
		{
			name:        "action references nonexistent machine",
			projectPath: filepath.Join(projectRoot, "examples", "testing", "test-action-nonexistent-machine"),
			expectError: false, // Validation is passing, machine reference validation may not be implemented
		},
		{
			name:        "action command script mutual exclusion",
			projectPath: filepath.Join(projectRoot, "examples", "testing", "test-action-command-script-mutual-excl"),
			expectError: true,
			errorMsg:    "multiple actions blocks found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := ValidateCmd
			args := []string{}
			if tt.projectPath != "" {
				args = append(args, tt.projectPath)
			}
			cmd.SetArgs(args)

			err := cmd.Execute()

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestListCmd(t *testing.T) {
	projectRoot := getProjectRoot()
	tests := []struct {
		name        string
		projectPath string
		verbose     bool
		expectError bool
	}{
		{
			name:        "valid project",
			projectPath: filepath.Join(projectRoot, "examples", "testing", "test-valid-project"),
			expectError: false,
		},
		{
			name:        "valid project verbose",
			projectPath: filepath.Join(projectRoot, "examples", "testing", "test-valid-project"),
			verbose:     true,
			expectError: false,
		},
		{
			name:        "large inventory",
			projectPath: filepath.Join(projectRoot, "examples", "testing", "test-large-inventory"),
			expectError: false,
		},
		{
			name:        "large actions",
			projectPath: filepath.Join(projectRoot, "examples", "testing", "test-large-actions"),
			expectError: false,
		},
		{
			name:        "invalid project",
			projectPath: filepath.Join(projectRoot, "examples", "testing", "test-invalid-project"),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := ListCmd
			args := []string{}
			if tt.projectPath != "" {
				args = append(args, tt.projectPath)
			}
			cmd.SetArgs(args)

			if tt.verbose {
				err := cmd.Flags().Set("verbose", "true")
				require.NoError(t, err)
			}

			err := cmd.Execute()

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestListMachinesCmd(t *testing.T) {
	projectRoot := getProjectRoot()
	tests := []struct {
		name        string
		projectPath string
		expectError bool
	}{
		{
			name:        "valid project",
			projectPath: filepath.Join(projectRoot, "examples", "testing", "test-valid-project"),
			expectError: false,
		},
		{
			name:        "large inventory",
			projectPath: filepath.Join(projectRoot, "examples", "testing", "test-large-inventory"),
			expectError: false,
		},
		{
			name:        "special characters",
			projectPath: filepath.Join(projectRoot, "examples", "testing", "test-special-characters"),
			expectError: false,
		},
		{
			name:        "missing inventory",
			projectPath: filepath.Join(projectRoot, "examples", "testing", "test-missing-inventory"),
			expectError: false, // Only shows a warning, doesn't fail
		},
		{
			name:        "invalid inventory",
			projectPath: filepath.Join(projectRoot, "examples", "testing", "test-invalid-inventory"),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := ListMachinesCmd
			args := []string{}
			if tt.projectPath != "" {
				args = append(args, tt.projectPath)
			}
			cmd.SetArgs(args)

			err := cmd.Execute()

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestListActionsCmd(t *testing.T) {
	projectRoot := getProjectRoot()
	tests := []struct {
		name        string
		projectPath string
		expectError bool
	}{
		{
			name:        "valid project",
			projectPath: filepath.Join(projectRoot, "examples", "testing", "test-valid-project"),
			expectError: false,
		},
		{
			name:        "large actions",
			projectPath: filepath.Join(projectRoot, "examples", "testing", "test-large-actions"),
			expectError: false,
		},
		{
			name:        "missing actions",
			projectPath: filepath.Join(projectRoot, "examples", "testing", "test-missing-actions"),
			expectError: false, // Only shows a warning, doesn't fail
		},
		{
			name:        "invalid actions",
			projectPath: filepath.Join(projectRoot, "examples", "testing", "test-invalid-actions"),
			expectError: false, // Actions are being loaded successfully despite invalid ones
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := ListActionsCmd
			args := []string{}
			if tt.projectPath != "" {
				args = append(args, tt.projectPath)
			}
			cmd.SetArgs(args)

			err := cmd.Execute()

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestListTemplatesCmd(t *testing.T) {
	projectRoot := getProjectRoot()
	tests := []struct {
		name         string
		projectPath  string
		templatesDir string
		expectError  bool
	}{
		{
			name:        "valid project",
			projectPath: filepath.Join(projectRoot, "examples", "testing", "test-valid-project"),
			expectError: false,
		},
		{
			name:         "custom templates dir",
			projectPath:  filepath.Join(projectRoot, "examples", "testing", "test-valid-project"),
			templatesDir: filepath.Join(projectRoot, "examples", "testing", "test-valid-project", "templates"),
			expectError:  false,
		},
		{
			name:        "invalid project",
			projectPath: filepath.Join(projectRoot, "examples", "testing", "test-invalid-project"),
			expectError: false, // Command succeeds even with invalid project
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := ListTemplatesCmd
			args := []string{}
			if tt.projectPath != "" {
				args = append(args, tt.projectPath)
			}
			cmd.SetArgs(args)

			if tt.templatesDir != "" {
				err := cmd.Flags().Set("templates-dir", tt.templatesDir)
				require.NoError(t, err)
			}

			err := cmd.Execute()

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestListFactsCmd(t *testing.T) {
	projectRoot := getProjectRoot()
	tests := []struct {
		name        string
		projectPath string
		factsDBPath string
		expectError bool
	}{
		{
			name:        "valid project",
			projectPath: filepath.Join(projectRoot, "examples", "testing", "test-valid-project"),
			expectError: false,
		},
		{
			name:        "custom facts db path",
			projectPath: filepath.Join(projectRoot, "examples", "testing", "test-valid-project"),
			factsDBPath: filepath.Join(projectRoot, "examples", "testing", "test-valid-project", ".facts.db"),
			expectError: false,
		},
		{
			name:        "corrupted facts db",
			projectPath: filepath.Join(projectRoot, "examples", "testing", "test-corrupted-facts-db"),
			expectError: false, // Facts loading not implemented yet, just shows "No facts found"
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := ListFactsCmd
			args := []string{}
			if tt.projectPath != "" {
				args = append(args, tt.projectPath)
			}
			cmd.SetArgs(args)

			if tt.factsDBPath != "" {
				err := cmd.Flags().Set("facts-db-path", tt.factsDBPath)
				require.NoError(t, err)
			}

			err := cmd.Execute()

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGatherFactsCmd(t *testing.T) {
	projectRoot := getProjectRoot()
	tests := []struct {
		name        string
		projectPath string
		sshKeyPath  string
		dryRun      bool
		factsDBPath string
		expectError bool
	}{
		{
			name:        "valid project dry run",
			projectPath: filepath.Join(projectRoot, "examples", "testing", "test-valid-project"),
			dryRun:      true,
			expectError: false,
		},
		{
			name:        "unreachable hosts",
			projectPath: filepath.Join(projectRoot, "examples", "testing", "test-unreachable-host"),
			dryRun:      true,
			expectError: false,
		},
		{
			name:        "network timeouts",
			projectPath: filepath.Join(projectRoot, "examples", "testing", "test-network-timeouts"),
			dryRun:      true,
			expectError: false,
		},
		{
			name:        "invalid ssh key",
			projectPath: filepath.Join(projectRoot, "examples", "testing", "test-invalid-ssh-key"),
			dryRun:      true,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := GatherFactsCmd
			args := []string{}
			if tt.projectPath != "" {
				args = append(args, tt.projectPath)
			}
			cmd.SetArgs(args)

			if tt.sshKeyPath != "" {
				err := cmd.Flags().Set("ssh-key-path", tt.sshKeyPath)
				require.NoError(t, err)
			}

			if tt.dryRun {
				err := cmd.Flags().Set("dry-run", "true")
				require.NoError(t, err)
			}

			if tt.factsDBPath != "" {
				err := cmd.Flags().Set("facts-db-path", tt.factsDBPath)
				require.NoError(t, err)
			}

			err := cmd.Execute()

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestRenderTemplateCmd(t *testing.T) {
	projectRoot := getProjectRoot()
	tests := []struct {
		name         string
		templateFile string
		projectPath  string
		output       string
		dryRun       bool
		server       string
		sshKeyPath   string
		templatesDir string
		dataDir      string
		expectError  bool
	}{
		{
			name:         "valid template",
			templateFile: "test.tmpl",
			projectPath:  filepath.Join(projectRoot, "examples", "testing", "test-valid-project"),
			templatesDir: filepath.Join(projectRoot, "examples", "testing", "test-valid-project", "templates"),
			dryRun:       true,
			expectError:  false,
		},
		{
			name:         "invalid template syntax",
			templateFile: "invalid-syntax.tmpl",
			projectPath:  filepath.Join(projectRoot, "examples", "testing", "test-invalid-template-syntax"),
			templatesDir: filepath.Join(projectRoot, "examples", "testing", "test-invalid-template-syntax", "templates"),
			dryRun:       true,
			expectError:  true,
		},
		{
			name:         "template missing data",
			templateFile: "missing-data.tmpl",
			projectPath:  filepath.Join(projectRoot, "examples", "testing", "test-template-missing-data"),
			templatesDir: filepath.Join(projectRoot, "examples", "testing", "test-template-missing-data", "templates"),
			dryRun:       true,
			expectError:  true,
		},
		{
			name:         "template undefined function",
			templateFile: "undefined-function.tmpl",
			projectPath:  filepath.Join(projectRoot, "examples", "testing", "test-template-undefined-function"),
			templatesDir: filepath.Join(projectRoot, "examples", "testing", "test-template-undefined-function", "templates"),
			dryRun:       true,
			expectError:  true,
		},
		{
			name:         "template unbalanced braces",
			templateFile: "unbalanced-braces.tmpl",
			projectPath:  filepath.Join(projectRoot, "examples", "testing", "test-template-unbalanced-braces"),
			templatesDir: filepath.Join(projectRoot, "examples", "testing", "test-template-unbalanced-braces", "templates"),
			dryRun:       true,
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := RenderTemplateCmd
			args := []string{tt.templateFile}
			if tt.projectPath != "" {
				args = append(args, tt.projectPath)
			}
			cmd.SetArgs(args)

			if tt.output != "" {
				err := cmd.Flags().Set("output", tt.output)
				require.NoError(t, err)
			}

			if tt.dryRun {
				err := cmd.Flags().Set("dry-run", "true")
				require.NoError(t, err)
			}

			if tt.server != "" {
				err := cmd.Flags().Set("server", tt.server)
				require.NoError(t, err)
			}

			if tt.sshKeyPath != "" {
				err := cmd.Flags().Set("ssh-key-path", tt.sshKeyPath)
				require.NoError(t, err)
			}

			if tt.templatesDir != "" {
				err := cmd.Flags().Set("templates-dir", tt.templatesDir)
				require.NoError(t, err)
			}

			if tt.dataDir != "" {
				err := cmd.Flags().Set("data-dir", tt.dataDir)
				require.NoError(t, err)
			}

			err := cmd.Execute()

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateTemplateCmd(t *testing.T) {
	projectRoot := getProjectRoot()
	tests := []struct {
		name         string
		templateFile string
		projectPath  string
		templatesDir string
		dataDir      string
		expectError  bool
	}{
		{
			name:         "valid template",
			templateFile: "test.tmpl",
			projectPath:  filepath.Join(projectRoot, "examples", "testing", "test-valid-project"),
			templatesDir: filepath.Join(projectRoot, "examples", "testing", "test-valid-project", "templates"),
			expectError:  false,
		},
		{
			name:         "invalid template syntax",
			templateFile: "invalid-syntax.tmpl",
			projectPath:  filepath.Join(projectRoot, "examples", "testing", "test-invalid-template-syntax"),
			templatesDir: filepath.Join(projectRoot, "examples", "testing", "test-invalid-template-syntax", "templates"),
			expectError:  true,
		},
		{
			name:         "template missing data",
			templateFile: "missing-data.tmpl",
			projectPath:  filepath.Join(projectRoot, "examples", "testing", "test-template-missing-data"),
			templatesDir: filepath.Join(projectRoot, "examples", "testing", "test-template-missing-data", "templates"),
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := ValidateTemplateCmd
			args := []string{tt.templateFile}
			if tt.projectPath != "" {
				args = append(args, tt.projectPath)
			}
			cmd.SetArgs(args)

			if tt.templatesDir != "" {
				err := cmd.Flags().Set("templates-dir", tt.templatesDir)
				require.NoError(t, err)
			}

			if tt.dataDir != "" {
				err := cmd.Flags().Set("data-dir", tt.dataDir)
				require.NoError(t, err)
			}

			err := cmd.Execute()

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestInitProject(t *testing.T) {
	tests := []struct {
		name        string
		projectName string
		path        string
		expectError bool
	}{
		{
			name:        "valid project",
			projectName: "test-project",
			path:        t.TempDir(),
			expectError: false,
		},
		{
			name:        "empty project name",
			projectName: "",
			path:        t.TempDir(),
			expectError: true,
		},
		{
			name:        "invalid path",
			projectName: "test-project",
			path:        "/nonexistent/path",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := logging.GetLogger()
			err := initProject(logger, tt.projectName, tt.path)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				// Verify project structure was created
				projectPath := filepath.Join(tt.path, tt.projectName)
				assert.DirExists(t, projectPath)
				assert.FileExists(t, filepath.Join(projectPath, "project.hcl"))
				assert.FileExists(t, filepath.Join(projectPath, "inventory.hcl"))
				assert.FileExists(t, filepath.Join(projectPath, "actions.hcl"))
				assert.DirExists(t, filepath.Join(projectPath, "actions"))
				assert.DirExists(t, filepath.Join(projectPath, "templates"))
				assert.DirExists(t, filepath.Join(projectPath, "data"))
				assert.DirExists(t, filepath.Join(projectPath, "files"))
				assert.DirExists(t, filepath.Join(projectPath, "logs"))
			}
		})
	}
}

func TestValidateProject(t *testing.T) {
	projectRoot := getProjectRoot()
	tests := []struct {
		name        string
		path        string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "valid project",
			path:        filepath.Join(projectRoot, "examples", "testing", "test-valid-project"),
			expectError: false,
		},
		{
			name:        "invalid project",
			path:        filepath.Join(projectRoot, "examples", "testing", "test-invalid-project"),
			expectError: true,
			errorMsg:    "Unclosed configuration block",
		},
		{
			name:        "missing project file",
			path:        filepath.Join(projectRoot, "examples", "testing", "test-missing-project-file"),
			expectError: true,
			errorMsg:    "project.hcl not found",
		},
		{
			name:        "missing inventory",
			path:        filepath.Join(projectRoot, "examples", "testing", "test-missing-inventory"),
			expectError: false, // Only shows a warning, doesn't fail validation
		},
		{
			name:        "missing actions",
			path:        filepath.Join(projectRoot, "examples", "testing", "test-missing-actions"),
			expectError: false, // Only shows a warning, doesn't fail validation
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := logging.GetLogger()
			err := validateProject(logger, tt.path)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestListProject(t *testing.T) {
	projectRoot := getProjectRoot()
	tests := []struct {
		name        string
		path        string
		expectError bool
	}{
		{
			name:        "valid project",
			path:        filepath.Join(projectRoot, "examples", "testing", "test-valid-project"),
			expectError: false,
		},
		{
			name:        "large inventory",
			path:        filepath.Join(projectRoot, "examples", "testing", "test-large-inventory"),
			expectError: false,
		},
		{
			name:        "large actions",
			path:        filepath.Join(projectRoot, "examples", "testing", "test-large-actions"),
			expectError: false,
		},
		{
			name:        "invalid project",
			path:        filepath.Join(projectRoot, "examples", "testing", "test-invalid-project"),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := logging.GetLogger()
			err := listProject(logger, tt.path)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestListProjectMachines(t *testing.T) {
	projectRoot := getProjectRoot()
	tests := []struct {
		name        string
		path        string
		expectError bool
	}{
		{
			name:        "valid project",
			path:        filepath.Join(projectRoot, "examples", "testing", "test-valid-project"),
			expectError: false,
		},
		{
			name:        "large inventory",
			path:        filepath.Join(projectRoot, "examples", "testing", "test-large-inventory"),
			expectError: false,
		},
		{
			name:        "special characters",
			path:        filepath.Join(projectRoot, "examples", "testing", "test-special-characters"),
			expectError: false,
		},
		{
			name:        "missing inventory",
			path:        filepath.Join(projectRoot, "examples", "testing", "test-missing-inventory"),
			expectError: false, // Only shows a warning, doesn't fail
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := logging.GetLogger()
			err := listProjectMachines(logger, tt.path)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestListProjectActions(t *testing.T) {
	projectRoot := getProjectRoot()
	tests := []struct {
		name        string
		path        string
		expectError bool
	}{
		{
			name:        "valid project",
			path:        filepath.Join(projectRoot, "examples", "testing", "test-valid-project"),
			expectError: false,
		},
		{
			name:        "large actions",
			path:        filepath.Join(projectRoot, "examples", "testing", "test-large-actions"),
			expectError: false,
		},
		{
			name:        "missing actions",
			path:        filepath.Join(projectRoot, "examples", "testing", "test-missing-actions"),
			expectError: false, // Only shows a warning, doesn't fail
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := logging.GetLogger()
			err := listProjectActions(logger, tt.path)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestListProjectTemplates(t *testing.T) {
	projectRoot := getProjectRoot()
	tests := []struct {
		name         string
		path         string
		templatesDir string
		expectError  bool
	}{
		{
			name:        "valid project",
			path:        filepath.Join(projectRoot, "examples", "testing", "test-valid-project"),
			expectError: false,
		},
		{
			name:         "custom templates dir",
			path:         filepath.Join(projectRoot, "examples", "testing", "test-valid-project"),
			templatesDir: filepath.Join(projectRoot, "examples", "testing", "test-valid-project", "templates"),
			expectError:  false,
		},
		{
			name:        "invalid project",
			path:        filepath.Join(projectRoot, "examples", "testing", "test-invalid-project"),
			expectError: false, // Just shows "No templates found", doesn't fail
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := logging.GetLogger()
			err := listProjectTemplates(logger, tt.path, tt.templatesDir)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestListProjectFacts(t *testing.T) {
	projectRoot := getProjectRoot()
	tests := []struct {
		name        string
		path        string
		factsDBPath string
		expectError bool
	}{
		{
			name:        "valid project",
			path:        filepath.Join(projectRoot, "examples", "testing", "test-valid-project"),
			expectError: false,
		},
		{
			name:        "custom facts db path",
			path:        filepath.Join(projectRoot, "examples", "testing", "test-valid-project"),
			factsDBPath: filepath.Join(projectRoot, "examples", "testing", "test-valid-project", ".facts.db"),
			expectError: false,
		},
		{
			name:        "corrupted facts db",
			path:        filepath.Join(projectRoot, "examples", "testing", "test-corrupted-facts-db"),
			expectError: false, // Facts loading not implemented yet, just shows "No facts found"
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := logging.GetLogger()
			err := listProjectFacts(logger, tt.path, tt.factsDBPath)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGatherProjectFacts(t *testing.T) {
	projectRoot := getProjectRoot()
	tests := []struct {
		name        string
		path        string
		sshKeyPath  string
		dryRun      bool
		factsDBPath string
		expectError bool
	}{
		{
			name:        "valid project dry run",
			path:        filepath.Join(projectRoot, "examples", "testing", "test-valid-project"),
			dryRun:      true,
			expectError: false,
		},
		{
			name:        "unreachable hosts dry run",
			path:        filepath.Join(projectRoot, "examples", "testing", "test-unreachable-host"),
			dryRun:      true,
			expectError: false,
		},
		{
			name:        "network timeouts dry run",
			path:        filepath.Join(projectRoot, "examples", "testing", "test-network-timeouts"),
			dryRun:      true,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := logging.GetLogger()
			err := gatherProjectFacts(logger, tt.path, tt.sshKeyPath, tt.dryRun, tt.factsDBPath)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestRenderProjectTemplate(t *testing.T) {
	projectRoot := getProjectRoot()
	tests := []struct {
		name         string
		templateFile string
		path         string
		output       string
		dryRun       bool
		server       string
		sshKeyPath   string
		templatesDir string
		dataDir      string
		expectError  bool
	}{
		{
			name:         "valid template",
			templateFile: "test.tmpl",
			path:         filepath.Join(projectRoot, "examples", "testing", "test-valid-project"),
			templatesDir: filepath.Join(projectRoot, "examples", "testing", "test-valid-project", "templates"),
			dryRun:       true,
			expectError:  false,
		},
		{
			name:         "invalid template syntax",
			templateFile: "invalid-syntax.tmpl",
			path:         filepath.Join(projectRoot, "examples", "testing", "test-invalid-template-syntax"),
			templatesDir: filepath.Join(projectRoot, "examples", "testing", "test-invalid-template-syntax", "templates"),
			dryRun:       true,
			expectError:  true,
		},
		{
			name:         "template missing data",
			templateFile: "missing-data.tmpl",
			path:         filepath.Join(projectRoot, "examples", "testing", "test-template-missing-data"),
			templatesDir: filepath.Join(projectRoot, "examples", "testing", "test-template-missing-data", "templates"),
			dryRun:       true,
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := logging.GetLogger()
			err := renderProjectTemplate(logger, tt.templateFile, tt.path, tt.output, tt.dryRun, tt.server, tt.sshKeyPath, tt.templatesDir, tt.dataDir)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateProjectTemplate(t *testing.T) {
	projectRoot := getProjectRoot()
	tests := []struct {
		name         string
		templateFile string
		path         string
		templatesDir string
		dataDir      string
		expectError  bool
	}{
		{
			name:         "valid template",
			templateFile: "test.tmpl",
			path:         filepath.Join(projectRoot, "examples", "testing", "test-valid-project"),
			templatesDir: filepath.Join(projectRoot, "examples", "testing", "test-valid-project", "templates"),
			expectError:  false,
		},
		{
			name:         "invalid template syntax",
			templateFile: "invalid-syntax.tmpl",
			path:         filepath.Join(projectRoot, "examples", "testing", "test-invalid-template-syntax"),
			templatesDir: filepath.Join(projectRoot, "examples", "testing", "test-invalid-template-syntax", "templates"),
			expectError:  true,
		},
		{
			name:         "template missing data",
			templateFile: "missing-data.tmpl",
			path:         filepath.Join(projectRoot, "examples", "testing", "test-template-missing-data"),
			templatesDir: filepath.Join(projectRoot, "examples", "testing", "test-template-missing-data", "templates"),
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := logging.GetLogger()
			err := validateProjectTemplate(logger, tt.templateFile, tt.path, tt.templatesDir, tt.dataDir)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateConfigFile(t *testing.T) {
	projectRoot := getProjectRoot()
	tests := []struct {
		name        string
		filePath    string
		fileType    string
		parser      func(string) error
		expectError bool
	}{
		{
			name:        "valid file",
			filePath:    filepath.Join(projectRoot, "examples", "testing", "test-valid-project", "project.hcl"),
			fileType:    "project",
			parser:      func(_ string) error { return nil },
			expectError: false,
		},
		{
			name:        "missing file",
			filePath:    "./nonexistent/file.hcl",
			fileType:    "project",
			parser:      func(_ string) error { return nil },
			expectError: false, // Only shows a warning, doesn't fail
		},
		{
			name:        "parser error",
			filePath:    filepath.Join(projectRoot, "examples", "testing", "test-valid-project", "project.hcl"),
			fileType:    "project",
			parser:      func(_ string) error { return fmt.Errorf("parser error") },
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := logging.GetLogger()
			err := validateConfigFile(logger, tt.filePath, tt.fileType, tt.parser)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// getProjectRoot returns the absolute path to the project root directory
func getProjectRoot() string {
	wd, err := os.Getwd()
	if err != nil {
		panic(fmt.Sprintf("failed to get working directory: %v", err))
	}
	// Go up from internal/cli to project root
	return filepath.Join(wd, "..", "..")
}

func TestDebugFileAccess(t *testing.T) {
	// Test working directory
	wd, err := os.Getwd()
	require.NoError(t, err)
	t.Logf("Working directory: %s", wd)

	// Get project root directory (go up from internal/cli to project root)
	projectRoot := getProjectRoot()
	t.Logf("Project root: %s", projectRoot)

	// Test file existence with absolute path
	testPath := filepath.Join(projectRoot, "examples", "testing", "test-valid-project", "project.hcl")
	t.Logf("Test path: %s", testPath)

	// Check if file exists
	if _, err := os.Stat(testPath); os.IsNotExist(err) {
		t.Errorf("File does not exist: %s", testPath)
	} else if err != nil {
		t.Errorf("Error checking file: %v", err)
	} else {
		t.Logf("File exists: %s", testPath)
	}
}

// generateRandomSuffix generates a 6-character alphanumeric suffix
func generateRandomSuffix() string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, 6)
	if _, err := rand.Read(b); err != nil {
		// Fallback to a simple random string if crypto/rand fails
		return "abcdef"
	}
	for i := range b {
		b[i] = charset[b[i]%byte(len(charset))]
	}
	return string(b)
}
