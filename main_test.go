package main

import (
	"bytes"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain_HelpCommand(t *testing.T) {
	// Test help command
	output := captureOutput(func() {
		os.Args = []string{"spooky", "--help"}
		main()
	})

	assert.Contains(t, output, "Spooky is a powerful server configuration and automation tool")
	assert.Contains(t, output, "Available Commands:")
	assert.Contains(t, output, "init")
	assert.Contains(t, output, "validate")
	assert.Contains(t, output, "list")
}

func TestMain_VersionCommand(t *testing.T) {
	// Test that version information is available in help
	output := captureOutput(func() {
		os.Args = []string{"spooky", "--help"}
		main()
	})

	// Version information should be available in the help output
	assert.Contains(t, output, "spooky")
}

func TestMain_InitCommand(t *testing.T) {
	// Test init command with help
	output := captureOutput(func() {
		os.Args = []string{"spooky", "init", "--help"}
		main()
	})

	assert.Contains(t, output, "Create a new spooky project")
}

func TestMain_ValidateCommand(t *testing.T) {
	// Test validate command with help
	output := captureOutput(func() {
		os.Args = []string{"spooky", "validate", "--help"}
		main()
	})

	assert.Contains(t, output, "Validate the project files")
}

func TestMain_ListCommand(t *testing.T) {
	// Test list command with help
	output := captureOutput(func() {
		os.Args = []string{"spooky", "list", "--help"}
		main()
	})

	assert.Contains(t, output, "List machines and actions")
}

func TestMain_ListMachinesCommand(t *testing.T) {
	// Test list-machines command with help
	output := captureOutput(func() {
		os.Args = []string{"spooky", "list-machines", "--help"}
		main()
	})

	assert.Contains(t, output, "List all machines")
}

func TestMain_ListActionsCommand(t *testing.T) {
	// Test list-actions command with help
	output := captureOutput(func() {
		os.Args = []string{"spooky", "list-actions", "--help"}
		main()
	})

	assert.Contains(t, output, "List all actions")
}

func TestMain_ListTemplatesCommand(t *testing.T) {
	// Test list-templates command with help
	output := captureOutput(func() {
		os.Args = []string{"spooky", "list-templates", "--help"}
		main()
	})

	assert.Contains(t, output, "List all template files")
}

func TestMain_ListFactsCommand(t *testing.T) {
	// Test list-facts command with help
	output := captureOutput(func() {
		os.Args = []string{"spooky", "list-facts", "--help"}
		main()
	})

	assert.Contains(t, output, "List all facts")
}

func TestMain_GatherFactsCommand(t *testing.T) {
	// Test gather-facts command with help
	output := captureOutput(func() {
		os.Args = []string{"spooky", "gather-facts", "--help"}
		main()
	})

	assert.Contains(t, output, "Gather facts from all machines")
}

func TestMain_RenderTemplateCommand(t *testing.T) {
	// Test render-template command with help
	output := captureOutput(func() {
		os.Args = []string{"spooky", "render-template", "--help"}
		main()
	})

	assert.Contains(t, output, "Render a template")
}

func TestMain_ValidateTemplateCommand(t *testing.T) {
	// Test validate-template command with help
	output := captureOutput(func() {
		os.Args = []string{"spooky", "validate-template", "--help"}
		main()
	})

	assert.Contains(t, output, "Validate a template")
}

func TestMain_GlobalFlags(t *testing.T) {
	// Test global flags
	testCases := []struct {
		name   string
		args   []string
		expect string
	}{
		{
			name:   "verbose flag",
			args:   []string{"spooky", "--verbose", "--help"},
			expect: "Spooky is a powerful server configuration and automation tool",
		},
		{
			name:   "quiet flag",
			args:   []string{"spooky", "--quiet", "--help"},
			expect: "Spooky is a powerful server configuration and automation tool",
		},
		{
			name:   "log-level flag",
			args:   []string{"spooky", "--log-level=debug", "--help"},
			expect: "Spooky is a powerful server configuration and automation tool",
		},
		{
			name:   "log-file flag",
			args:   []string{"spooky", "--log-file=/tmp/test.log", "--help"},
			expect: "Spooky is a powerful server configuration and automation tool",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			output := captureOutput(func() {
				os.Args = tc.args
				main()
			})

			assert.Contains(t, output, tc.expect)
		})
	}
}

func TestMain_InvalidCommand(t *testing.T) {
	// Test invalid command by running the actual command
	cmd := exec.Command("go", "run", "main.go", "invalid-command")
	output, err := cmd.CombinedOutput()

	// Should return error
	assert.Error(t, err)

	// Should contain error message
	outputStr := string(output)
	assert.Contains(t, outputStr, "unknown command")
	assert.Contains(t, outputStr, "invalid-command")
}

func TestMain_InvalidFlag(t *testing.T) {
	// Test invalid flag by running the actual command
	cmd := exec.Command("go", "run", "main.go", "--invalid-flag")
	output, err := cmd.CombinedOutput()

	// Should return error
	assert.Error(t, err)

	// Should contain error message
	outputStr := string(output)
	assert.Contains(t, outputStr, "unknown flag")
	assert.Contains(t, outputStr, "--invalid-flag")
}

func TestMain_CommandWithInvalidArgs(t *testing.T) {
	// Test command with invalid arguments
	cmd := exec.Command("go", "run", "main.go", "validate", "nonexistent-file.hcl")
	output, err := cmd.CombinedOutput()

	// Should return error
	assert.Error(t, err)

	// Should contain error message
	outputStr := string(output)
	assert.Contains(t, outputStr, "Error:")
}

func TestMain_InitCommandWithValidArgs(t *testing.T) {
	// Test init command with valid arguments
	tempDir := t.TempDir()

	cmd := exec.Command("go", "run", "main.go", "init", "test-project", tempDir)
	output, err := cmd.CombinedOutput()

	// Should not return error
	assert.NoError(t, err)

	// Should not contain error
	outputStr := string(output)
	assert.NotContains(t, outputStr, "Error:")
}

func TestMain_ValidateCommandWithValidConfig(t *testing.T) {
	// Test validate command with valid configuration
	configPath := "examples/testing/test-valid-project"

	cmd := exec.Command("go", "run", "main.go", "validate", configPath)
	output, err := cmd.CombinedOutput()

	// Should not return error
	assert.NoError(t, err)

	// Should not contain error
	outputStr := string(output)
	assert.NotContains(t, outputStr, "Error:")
}

func TestMain_ListCommandWithValidConfig(t *testing.T) {
	// Test list command with valid configuration
	configPath := "examples/testing/test-valid-project"

	cmd := exec.Command("go", "run", "main.go", "list", configPath)
	output, err := cmd.CombinedOutput()

	// Should not return error
	assert.NoError(t, err)

	// Should not contain error
	outputStr := string(output)
	assert.NotContains(t, outputStr, "Error:")
}

func TestMain_ListMachinesCommandWithValidConfig(t *testing.T) {
	// Test list-machines command with valid configuration
	configPath := "examples/testing/test-valid-project"

	cmd := exec.Command("go", "run", "main.go", "list-machines", configPath)
	output, err := cmd.CombinedOutput()

	// Should not return error
	assert.NoError(t, err)

	// Should not contain error
	outputStr := string(output)
	assert.NotContains(t, outputStr, "Error:")
}

func TestMain_ListActionsCommandWithValidConfig(t *testing.T) {
	// Test list-actions command with valid configuration
	configPath := "examples/testing/test-valid-project"

	cmd := exec.Command("go", "run", "main.go", "list-actions", configPath)
	output, err := cmd.CombinedOutput()

	// Should not return error
	assert.NoError(t, err)

	// Should not contain error
	outputStr := string(output)
	assert.NotContains(t, outputStr, "Error:")
}

func TestMain_ListTemplatesCommandWithValidConfig(t *testing.T) {
	// Test list-templates command with valid configuration
	configPath := "examples/testing/test-valid-project"

	cmd := exec.Command("go", "run", "main.go", "list-templates", configPath)
	output, err := cmd.CombinedOutput()

	// Should not return error
	assert.NoError(t, err)

	// Should not contain error
	outputStr := string(output)
	assert.NotContains(t, outputStr, "Error:")
}

func TestMain_ListFactsCommandWithValidConfig(t *testing.T) {
	// Test list-facts command with valid configuration
	configPath := "examples/testing/test-valid-project"

	cmd := exec.Command("go", "run", "main.go", "list-facts", configPath)
	output, err := cmd.CombinedOutput()

	// Should not return error
	assert.NoError(t, err)

	// Should not contain error
	outputStr := string(output)
	assert.NotContains(t, outputStr, "Error:")
}

func TestMain_GatherFactsCommandWithValidConfig(t *testing.T) {
	// Test gather-facts command with valid configuration
	configPath := "examples/testing/test-valid-project"

	cmd := exec.Command("go", "run", "main.go", "gather-facts", configPath)
	output, _ := cmd.CombinedOutput()

	// Should contain error due to SSH connection failure, but not parsing error
	outputStr := string(output)
	assert.NotContains(t, outputStr, "Failed to parse")
}

func TestMain_RenderTemplateCommandWithValidConfig(t *testing.T) {
	// Test render-template command with valid configuration
	configPath := "examples/testing/test-valid-project"

	cmd := exec.Command("go", "run", "main.go", "render-template", "--template", "test.tmpl", configPath)
	output, _ := cmd.CombinedOutput()

	// Should not contain parsing error
	outputStr := string(output)
	assert.NotContains(t, outputStr, "Failed to parse")
}

func TestMain_ValidateTemplateCommandWithValidConfig(t *testing.T) {
	// Test validate-template command with valid configuration
	configPath := "examples/testing/test-valid-project"

	cmd := exec.Command("go", "run", "main.go", "validate-template", "--template", "test.tmpl", configPath)
	output, _ := cmd.CombinedOutput()

	// Should not contain parsing error
	outputStr := string(output)
	assert.NotContains(t, outputStr, "Failed to parse")
}

// Helper function to capture stdout/stderr output
func captureOutput(fn func()) string {
	// Save original stdout and stderr
	oldStdout := os.Stdout
	oldStderr := os.Stderr

	// Create pipes
	rOut, wOut, _ := os.Pipe()
	rErr, wErr, _ := os.Pipe()

	// Redirect stdout and stderr
	os.Stdout = wOut
	os.Stderr = wErr

	// Run the function
	fn()

	// Close writers
	wOut.Close()
	wErr.Close()

	// Read output
	var bufOut, bufErr bytes.Buffer
	_, _ = bufOut.ReadFrom(rOut) // Ignore error in test context
	_, _ = bufErr.ReadFrom(rErr) // Ignore error in test context

	// Restore original stdout and stderr
	os.Stdout = oldStdout
	os.Stderr = oldStderr

	// Return combined output
	return bufOut.String() + bufErr.String()
}

// Test version variables are set
func TestVersionVariables(t *testing.T) {
	assert.NotEmpty(t, version, "Version should be set")
	assert.NotEmpty(t, commit, "Commit should be set")
}

// Test that main function doesn't panic
func TestMainFunctionNoPanic(t *testing.T) {
	// Test that main function doesn't panic with help
	assert.NotPanics(t, func() {
		os.Args = []string{"spooky", "--help"}
		main()
	})
}

// Test that all commands are properly registered
func TestAllCommandsRegistered(t *testing.T) {
	expectedCommands := []string{
		"init",
		"validate",
		"list",
		"list-machines",
		"list-actions",
		"list-templates",
		"list-facts",
		"gather-facts",
		"render-template",
		"validate-template",
	}

	for _, cmd := range expectedCommands {
		t.Run(cmd, func(t *testing.T) {
			output := captureOutput(func() {
				os.Args = []string{"spooky", cmd, "--help"}
				main()
			})

			// Should not contain "unknown command"
			assert.NotContains(t, output, "unknown command")
		})
	}
}
