package ssh

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"spooky/internal/config"
)

func TestNewTemplateActionExecutor(t *testing.T) {
	executor := NewTemplateActionExecutor()
	assert.NotNil(t, executor)
}

func TestTemplateActionExecutor_ExecuteAction_NilAction(t *testing.T) {
	executor := NewTemplateActionExecutor()
	machines := []*config.Machine{
		{
			Name:     "test-server",
			Host:     "192.168.1.100",
			Port:     22,
			User:     "testuser",
			Password: "testpass",
		},
	}

	err := executor.ExecuteAction(nil, machines)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "action cannot be nil")
}

func TestTemplateActionExecutor_ExecuteAction_NilTemplate(t *testing.T) {
	executor := NewTemplateActionExecutor()
	action := &config.Action{
		Name: "test-action",
		Type: "template_deploy",
		// Template is nil
	}
	machines := []*config.Machine{
		{
			Name:     "test-server",
			Host:     "192.168.1.100",
			Port:     22,
			User:     "testuser",
			Password: "testpass",
		},
	}

	err := executor.ExecuteAction(action, machines)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "template configuration is required")
}

func TestTemplateActionExecutor_ExecuteAction_UnsupportedType(t *testing.T) {
	executor := NewTemplateActionExecutor()
	action := &config.Action{
		Name: "test-action",
		Type: "unsupported_type",
		Template: &config.TemplateConfig{
			Source:      "test.tmpl",
			Destination: "/tmp/test.conf",
		},
	}
	machines := []*config.Machine{
		{
			Name:     "test-server",
			Host:     "192.168.1.100",
			Port:     22,
			User:     "testuser",
			Password: "testpass",
		},
	}

	err := executor.ExecuteAction(action, machines)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported template action type")
}

func TestTemplateActionExecutor_ValidateTemplateSyntax(t *testing.T) {
	executor := NewTemplateActionExecutor()

	tests := []struct {
		name        string
		content     []byte
		expectError bool
	}{
		{
			name:        "valid template",
			content:     []byte(`Hello {{.name}}!`),
			expectError: false,
		},
		{
			name:        "valid template with functions",
			content:     []byte(`Hostname: {{hostname}}`),
			expectError: false,
		},
		{
			name:        "invalid template - unbalanced braces",
			content:     []byte(`Hello {{.name}!`),
			expectError: true,
		},
		{
			name:        "invalid template - undefined function",
			content:     []byte(`Hello {{undefined_function}}!`),
			expectError: true,
		},
		{
			name:        "empty template",
			content:     []byte(``),
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := executor.validateTemplateSyntax(tt.content)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestTemplateActionExecutor_ValidateTemplateSyntax_WithTestExamples(t *testing.T) {
	executor := NewTemplateActionExecutor()
	examplesDir := "../../examples/testing"

	// Test valid template
	validTemplatePath := filepath.Join(examplesDir, "test-valid-project", "templates", "test.tmpl")
	if _, err := os.Stat(validTemplatePath); err == nil {
		content, err := os.ReadFile(validTemplatePath)
		require.NoError(t, err)

		err = executor.validateTemplateSyntax(content)
		assert.NoError(t, err, "Valid template should pass validation")
	}

	// Test invalid template (unbalanced braces)
	invalidTemplatePath := filepath.Join(examplesDir, "test-template-unbalanced-braces", "templates", "unbalanced.tmpl")
	if _, err := os.Stat(invalidTemplatePath); err == nil {
		content, err := os.ReadFile(invalidTemplatePath)
		require.NoError(t, err)

		err = executor.validateTemplateSyntax(content)
		assert.Error(t, err, "Invalid template should fail validation")
		assert.Contains(t, err.Error(), "template syntax")
	}

	// Test template with undefined function
	undefinedFuncPath := filepath.Join(examplesDir, "test-template-undefined-function", "templates", "undefined-function.tmpl")
	if _, err := os.Stat(undefinedFuncPath); err == nil {
		content, err := os.ReadFile(undefinedFuncPath)
		require.NoError(t, err)

		err = executor.validateTemplateSyntax(content)
		assert.Error(t, err, "Template with undefined function should fail validation")
		assert.Contains(t, err.Error(), "template syntax")
	}
}

func TestTemplateActionExecutor_CreateFuncMap(t *testing.T) {
	executor := NewTemplateActionExecutor()

	// Test with test values
	funcMap := executor.createValidationFuncMap(true)
	assert.NotNil(t, funcMap)

	// Verify that functions return test values
	if machineIDFunc, ok := funcMap["machineID"].(func() string); ok {
		result := machineIDFunc()
		assert.Equal(t, "test", result)
	}

	if hostnameFunc, ok := funcMap["hostname"].(func() string); ok {
		result := hostnameFunc()
		assert.Equal(t, "test", result)
	}

	if fileExistsFunc, ok := funcMap["fileExists"].(func(string) bool); ok {
		result := fileExistsFunc("any-path")
		assert.True(t, result)
	}

	// Test with empty values
	funcMap = executor.createValidationFuncMap(false)
	assert.NotNil(t, funcMap)

	if machineIDFunc, ok := funcMap["machineID"].(func() string); ok {
		result := machineIDFunc()
		assert.Equal(t, "", result)
	}

	if fileExistsFunc, ok := funcMap["fileExists"].(func(string) bool); ok {
		result := fileExistsFunc("any-path")
		assert.False(t, result)
	}
}

func TestTemplateActionExecutor_CreateFuncMapWithValues(t *testing.T) {
	// Test with specific values
	funcMap := createFuncMapWithValues(
		"test-machine-id",
		"test-os-version",
		"test-hostname",
		"test-ip",
		"test-disk",
		"test-memory",
		true,
		"test-content",
		"test-size",
		"test-owner",
	)

	assert.NotNil(t, funcMap)

	// Test machineID function
	if machineIDFunc, ok := funcMap["machineID"].(func() string); ok {
		result := machineIDFunc()
		assert.Equal(t, "test-machine-id", result)
	}

	// Test osVersion function
	if osVersionFunc, ok := funcMap["osVersion"].(func() string); ok {
		result := osVersionFunc()
		assert.Equal(t, "test-os-version", result)
	}

	// Test hostname function
	if hostnameFunc, ok := funcMap["hostname"].(func() string); ok {
		result := hostnameFunc()
		assert.Equal(t, "test-hostname", result)
	}

	// Test ipAddress function
	if ipAddressFunc, ok := funcMap["ipAddress"].(func() string); ok {
		result := ipAddressFunc()
		assert.Equal(t, "test-ip", result)
	}

	// Test diskSpace function
	if diskSpaceFunc, ok := funcMap["diskSpace"].(func() string); ok {
		result := diskSpaceFunc()
		assert.Equal(t, "test-disk", result)
	}

	// Test memoryInfo function
	if memoryInfoFunc, ok := funcMap["memoryInfo"].(func() string); ok {
		result := memoryInfoFunc()
		assert.Equal(t, "test-memory", result)
	}

	// Test fileExists function
	if fileExistsFunc, ok := funcMap["fileExists"].(func(string) bool); ok {
		result := fileExistsFunc("any-path")
		assert.True(t, result)
	}

	// Test fileContent function
	if fileContentFunc, ok := funcMap["fileContent"].(func(string) string); ok {
		result := fileContentFunc("any-path")
		assert.Equal(t, "test-content", result)
	}

	// Test fileSize function
	if fileSizeFunc, ok := funcMap["fileSize"].(func(string) string); ok {
		result := fileSizeFunc("any-path")
		assert.Equal(t, "test-size", result)
	}

	// Test fileOwner function
	if fileOwnerFunc, ok := funcMap["fileOwner"].(func(string) string); ok {
		result := fileOwnerFunc("any-path")
		assert.Equal(t, "test-owner", result)
	}
}

func TestTemplateActionExecutor_ExecuteTemplateDeploy_FileNotExists(t *testing.T) {
	executor := NewTemplateActionExecutor()
	action := &config.Action{
		Name: "test-action",
		Type: "template_deploy",
		Template: &config.TemplateConfig{
			Source:      "/nonexistent/template.tmpl",
			Destination: "/tmp/test.conf",
		},
	}
	machines := []*config.Machine{
		{
			Name:     "test-server",
			Host:     "192.168.1.100",
			Port:     22,
			User:     "testuser",
			Password: "testpass",
		},
	}

	err := executor.executeTemplateDeploy(action, machines)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "template file does not exist")
}

func TestTemplateActionExecutor_ExecuteTemplateDeploy_ValidTemplate(t *testing.T) {
	// Create a temporary template file
	tempDir := t.TempDir()
	templateFile := filepath.Join(tempDir, "test.tmpl")
	templateContent := `Hello {{.name}}!
Hostname: {{hostname}}
OS: {{osVersion}}`

	err := os.WriteFile(templateFile, []byte(templateContent), 0o600)
	require.NoError(t, err)

	executor := NewTemplateActionExecutor()
	action := &config.Action{
		Name: "test-action",
		Type: "template_deploy",
		Template: &config.TemplateConfig{
			Source:      templateFile,
			Destination: "/tmp/test.conf",
		},
	}
	machines := []*config.Machine{
		{
			Name:     "test-server",
			Host:     "192.168.1.100",
			Port:     22,
			User:     "testuser",
			Password: "testpass",
		},
	}

	// This will fail because we can't actually connect to the server
	// but we can test that the template validation passes
	err = executor.executeTemplateDeploy(action, machines)
	// Should fail due to SSH connection, not template validation
	assert.Error(t, err)
	// But it shouldn't be a template syntax error
	assert.NotContains(t, err.Error(), "template syntax validation failed")
}

func TestTemplateActionExecutor_ExecuteTemplateDeploy_InvalidTemplate(t *testing.T) {
	// Create a temporary template file with invalid syntax
	tempDir := t.TempDir()
	templateFile := filepath.Join(tempDir, "invalid.tmpl")
	templateContent := `Hello {{.name}!
Hostname: {{hostname}
OS: {{osVersion}`

	err := os.WriteFile(templateFile, []byte(templateContent), 0o600)
	require.NoError(t, err)

	executor := NewTemplateActionExecutor()
	action := &config.Action{
		Name: "test-action",
		Type: "template_deploy",
		Template: &config.TemplateConfig{
			Source:      templateFile,
			Destination: "/tmp/test.conf",
		},
	}
	machines := []*config.Machine{
		{
			Name:     "test-server",
			Host:     "192.168.1.100",
			Port:     22,
			User:     "testuser",
			Password: "testpass",
		},
	}

	err = executor.executeTemplateDeploy(action, machines)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "template syntax validation failed")
}

func TestTemplateActionExecutor_ExecuteTemplateOperation(t *testing.T) {
	executor := NewTemplateActionExecutor()
	action := &config.Action{
		Name: "test-action",
		Type: "template_validate",
		Template: &config.TemplateConfig{
			Source:      "test.tmpl",
			Destination: "/tmp/test.conf",
		},
	}
	machines := []*config.Machine{
		{
			Name:     "test-server",
			Host:     "192.168.1.100",
			Port:     22,
			User:     "testuser",
			Password: "testpass",
		},
	}

	// Test with a simple operation that always succeeds
	operation := func(_ *SSHClient, _ *config.Action) error {
		return nil
	}

	err := executor.executeTemplateOperation(action, machines, "testing", "tested", operation)
	// Should fail due to SSH connection, not the operation itself
	assert.Error(t, err)
	// But it shouldn't be an operation error
	assert.NotContains(t, err.Error(), "testing failed")
}

func TestTemplateActionExecutor_ExecuteTemplateOperation_WithFailingSSH(t *testing.T) {
	executor := NewTemplateActionExecutor()
	action := &config.Action{
		Name: "test-action",
		Type: "template_validate",
		Template: &config.TemplateConfig{
			Source:      "test.tmpl",
			Destination: "/tmp/test.conf",
		},
	}
	machines := []*config.Machine{
		{
			Name:     "test-server",
			Host:     "192.168.1.100",
			Port:     22,
			User:     "testuser",
			Password: "testpass",
		},
	}

	// Test with an operation that would fail if SSH worked
	operation := func(_ *SSHClient, _ *config.Action) error {
		return assert.AnError
	}

	err := executor.executeTemplateOperation(action, machines, "testing", "tested", operation)
	// Should fail due to SSH connection, not the operation itself
	assert.Error(t, err)
	// But it shouldn't be the operation error we defined
	assert.NotEqual(t, assert.AnError, err)
}

func TestTemplateActionExecutor_WithTestDataFromExamples(t *testing.T) {
	// Use test data from examples/testing where available
	examplesDir := "../../examples/testing"

	// Test with valid project configuration
	validProjectPath := filepath.Join(examplesDir, "test-valid-project")
	if _, err := os.Stat(validProjectPath); os.IsNotExist(err) {
		t.Skip("Test data directory not found, skipping test")
	}

	executor := NewTemplateActionExecutor()

	// Create an action similar to what would be in test data
	action := &config.Action{
		Name: "deploy-config",
		Type: "template_deploy",
		Template: &config.TemplateConfig{
			Source:      filepath.Join(validProjectPath, "templates", "test.tmpl"),
			Destination: "/tmp/test.conf",
			Permissions: "644",
			Owner:       "root",
			Group:       "root",
		},
		Machines: []string{"example-server"},
		Parallel: true,
		Timeout:  300,
	}

	machines := []*config.Machine{
		{
			Name:     "example-server",
			Host:     "192.168.1.100",
			Port:     22,
			User:     "testuser",
			Password: "testpass",
			Tags: map[string]string{
				"environment": "test",
				"role":        "web",
			},
		},
	}

	// This will fail due to SSH connection, but we can test the action structure
	err := executor.ExecuteAction(action, machines)
	assert.Error(t, err)
	// But it shouldn't be a template validation error if the template exists
	if _, err := os.Stat(action.Template.Source); err == nil {
		assert.NotContains(t, err.Error(), "template syntax validation failed")
	}
}

func TestTemplateActionExecutor_TemplateActionTypes(t *testing.T) {
	executor := NewTemplateActionExecutor()

	// Create a temporary template file
	tempDir := t.TempDir()
	templateFile := filepath.Join(tempDir, "test.tmpl")
	templateContent := `Hello {{.name}}!`

	err := os.WriteFile(templateFile, []byte(templateContent), 0o600)
	require.NoError(t, err)

	machines := []*config.Machine{
		{
			Name:     "test-server",
			Host:     "192.168.1.100",
			Port:     22,
			User:     "testuser",
			Password: "testpass",
		},
	}

	tests := []struct {
		name        string
		actionType  string
		description string
	}{
		{
			name:        "template_deploy",
			actionType:  "template_deploy",
			description: "Deploy template files to target servers",
		},
		{
			name:        "template_evaluate",
			actionType:  "template_evaluate",
			description: "Evaluate templates on servers with server-specific facts",
		},
		{
			name:        "template_validate",
			actionType:  "template_validate",
			description: "Validate templates on servers",
		},
		{
			name:        "template_cleanup",
			actionType:  "template_cleanup",
			description: "Remove template files from servers",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			action := &config.Action{
				Name: tt.name,
				Type: tt.actionType,
				Template: &config.TemplateConfig{
					Source:      templateFile,
					Destination: "/tmp/test.conf",
				},
			}

			err := executor.ExecuteAction(action, machines)
			// All should fail due to SSH connection, but not due to unsupported type
			assert.Error(t, err)
			assert.NotContains(t, err.Error(), "unsupported template action type")
		})
	}
}
