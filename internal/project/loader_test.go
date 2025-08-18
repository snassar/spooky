package project

import (
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	spookytypeslogging "spooky/internal/types/logging"
	spookytypesproject "spooky/internal/types/project"
)

// MockLogger is a mock implementation of Logger for testing
type MockLogger struct{}

func (m *MockLogger) Debug(msg string, fields ...map[string]interface{})                 {}
func (m *MockLogger) Info(msg string, fields ...map[string]interface{})                  {}
func (m *MockLogger) Warn(msg string, fields ...map[string]interface{})                  {}
func (m *MockLogger) Error(msg string, err error, fields ...map[string]interface{})      {}
func (m *MockLogger) Fatal(msg string, err error, fields ...map[string]interface{})      {}
func (m *MockLogger) WithFields(fields map[string]interface{}) spookytypeslogging.Logger { return m }
func (m *MockLogger) WithComponent(component string) spookytypeslogging.Logger           { return m }
func (m *MockLogger) WithOperation(operation string) spookytypeslogging.Logger           { return m }
func (m *MockLogger) SetLevel(level spookytypeslogging.LogLevel)                         {}
func (m *MockLogger) GetLevel() spookytypeslogging.LogLevel                              { return spookytypeslogging.LogLevelInfo }

func TestGetSettingsBlockSchema(t *testing.T) {
	schema := getSettingsBlockSchema()

	assert.NotNil(t, schema)
	assert.Len(t, schema.Attributes, 7)

	expectedAttributes := []string{
		"parallel_workers",
		"timeout_seconds",
		"log_level",
		"default_dry_run",
		"validate_before_run",
		"max_retries",
		"retry_delay_seconds",
	}

	for _, attrName := range expectedAttributes {
		found := false
		for _, attr := range schema.Attributes {
			if attr.Name == attrName {
				found = true
				assert.False(t, attr.Required, "Attribute %s should not be required", attrName)
				break
			}
		}
		assert.True(t, found, "Expected attribute %s not found in schema", attrName)
	}
}

func TestParseBlockContent(t *testing.T) {
	// Test valid settings block
	hclContent := `
	settings {
		parallel_workers = 4
		timeout_seconds = 300
		log_level = "info"
		default_dry_run = true
		validate_before_run = false
		max_retries = 3
		retry_delay_seconds = 5
	}
	`

	parser := hclparse.NewParser()
	file, diags := parser.ParseHCL([]byte(hclContent), "test.hcl")
	require.False(t, diags.HasErrors(), "Failed to parse HCL: %v", diags)

	body := file.Body
	content, diags := body.Content(&hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{
			{Type: "settings"},
		},
	})
	require.False(t, diags.HasErrors(), "Failed to get body content: %v", diags)
	require.Len(t, content.Blocks, 1, "Expected exactly one settings block")

	block := content.Blocks[0]
	parsedContent, err := parseBlockContent(block)

	assert.NoError(t, err)
	assert.NotNil(t, parsedContent)
	assert.Len(t, parsedContent.Attributes, 7)
}

func TestParseIntegerAttribute(t *testing.T) {
	// Create test content with integer attributes
	hclContent := `
	settings {
		parallel_workers = 4
		timeout_seconds = 300
	}
	`

	parser := hclparse.NewParser()
	file, diags := parser.ParseHCL([]byte(hclContent), "test.hcl")
	require.False(t, diags.HasErrors(), "Failed to parse HCL: %v", diags)

	body := file.Body
	content, diags := body.Content(&hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{
			{Type: "settings"},
		},
	})
	require.False(t, diags.HasErrors(), "Failed to get body content: %v", diags)

	block := content.Blocks[0]
	parsedContent, err := parseBlockContent(block)
	require.NoError(t, err)

	// Test existing attribute
	value, err := parseIntegerAttribute(parsedContent, "parallel_workers")
	assert.NoError(t, err)
	assert.Equal(t, 4, value)

	// Test non-existing attribute
	value, err = parseIntegerAttribute(parsedContent, "non_existing")
	assert.NoError(t, err)
	assert.Equal(t, 0, value)
}

func TestParseStringAttribute(t *testing.T) {
	// Create test content with string attributes
	hclContent := `
	settings {
		log_level = "debug"
	}
	`

	parser := hclparse.NewParser()
	file, diags := parser.ParseHCL([]byte(hclContent), "test.hcl")
	require.False(t, diags.HasErrors(), "Failed to parse HCL: %v", diags)

	body := file.Body
	content, diags := body.Content(&hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{
			{Type: "settings"},
		},
	})
	require.False(t, diags.HasErrors(), "Failed to get body content: %v", diags)

	block := content.Blocks[0]
	parsedContent, err := parseBlockContent(block)
	require.NoError(t, err)

	// Test existing attribute
	value, err := parseStringAttribute(parsedContent, "log_level")
	assert.NoError(t, err)
	assert.Equal(t, "debug", value)

	// Test non-existing attribute
	value, err = parseStringAttribute(parsedContent, "non_existing")
	assert.NoError(t, err)
	assert.Equal(t, "", value)
}

func TestParseBooleanAttribute(t *testing.T) {
	// Create test content with boolean attributes
	hclContent := `
	settings {
		default_dry_run = true
		validate_before_run = false
	}
	`

	parser := hclparse.NewParser()
	file, diags := parser.ParseHCL([]byte(hclContent), "test.hcl")
	require.False(t, diags.HasErrors(), "Failed to parse HCL: %v", diags)

	body := file.Body
	content, diags := body.Content(&hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{
			{Type: "settings"},
		},
	})
	require.False(t, diags.HasErrors(), "Failed to get body content: %v", diags)

	block := content.Blocks[0]
	parsedContent, err := parseBlockContent(block)
	require.NoError(t, err)

	// Test existing attributes
	value, err := parseBooleanAttribute(parsedContent, "default_dry_run")
	assert.NoError(t, err)
	assert.True(t, value)

	value, err = parseBooleanAttribute(parsedContent, "validate_before_run")
	assert.NoError(t, err)
	assert.False(t, value)

	// Test non-existing attribute
	value, err = parseBooleanAttribute(parsedContent, "non_existing")
	assert.NoError(t, err)
	assert.False(t, value)
}

func TestPopulateSettings(t *testing.T) {
	// Create test content with all attributes
	hclContent := `
	settings {
		parallel_workers = 4
		timeout_seconds = 300
		log_level = "debug"
		default_dry_run = true
		validate_before_run = false
		max_retries = 3
		retry_delay_seconds = 5
	}
	`

	parser := hclparse.NewParser()
	file, diags := parser.ParseHCL([]byte(hclContent), "test.hcl")
	require.False(t, diags.HasErrors(), "Failed to parse HCL: %v", diags)

	body := file.Body
	content, diags := body.Content(&hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{
			{Type: "settings"},
		},
	})
	require.False(t, diags.HasErrors(), "Failed to get body content: %v", diags)

	block := content.Blocks[0]
	parsedContent, err := parseBlockContent(block)
	require.NoError(t, err)

	// Create settings and populate
	settings := &spookytypesproject.Settings{}
	err = populateSettings(settings, parsedContent)

	assert.NoError(t, err)
	assert.Equal(t, 4, settings.ParallelWorkers)
	assert.Equal(t, 300, settings.TimeoutSeconds)
	assert.Equal(t, "debug", settings.LogLevel)
	assert.True(t, settings.DefaultDryRun)
	assert.False(t, settings.ValidateBeforeRun)
	assert.Equal(t, 3, settings.MaxRetries)
	assert.Equal(t, 5, settings.RetryDelaySeconds)
}

func TestParseSettingsBlock(t *testing.T) {
	// Create test content
	hclContent := `
	settings {
		parallel_workers = 4
		timeout_seconds = 300
		log_level = "debug"
		default_dry_run = true
		validate_before_run = false
		max_retries = 3
		retry_delay_seconds = 5
	}
	`

	parser := hclparse.NewParser()
	file, diags := parser.ParseHCL([]byte(hclContent), "test.hcl")
	require.False(t, diags.HasErrors(), "Failed to parse HCL: %v", diags)

	body := file.Body
	content, diags := body.Content(&hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{
			{Type: "settings"},
		},
	})
	require.False(t, diags.HasErrors(), "Failed to get body content: %v", diags)
	require.Len(t, content.Blocks, 1, "Expected exactly one settings block")

	// Test the helper functions directly
	parsedContent, err := parseBlockContent(content.Blocks[0])
	require.NoError(t, err)

	settings := &spookytypesproject.Settings{}
	err = populateSettings(settings, parsedContent)

	assert.NoError(t, err)
	assert.NotNil(t, settings)
	assert.Equal(t, 4, settings.ParallelWorkers)
	assert.Equal(t, 300, settings.TimeoutSeconds)
	assert.Equal(t, "debug", settings.LogLevel)
	assert.True(t, settings.DefaultDryRun)
	assert.False(t, settings.ValidateBeforeRun)
	assert.Equal(t, 3, settings.MaxRetries)
	assert.Equal(t, 5, settings.RetryDelaySeconds)
}

func TestParseSettingsBlockWithDefaults(t *testing.T) {
	// Test with empty settings block (should use defaults)
	hclContent := `
	settings {
	}
	`

	parser := hclparse.NewParser()
	file, diags := parser.ParseHCL([]byte(hclContent), "test.hcl")
	require.False(t, diags.HasErrors(), "Failed to parse HCL: %v", diags)

	body := file.Body
	content, diags := body.Content(&hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{
			{Type: "settings"},
		},
	})
	require.False(t, diags.HasErrors(), "Failed to get body content: %v", diags)
	require.Len(t, content.Blocks, 1, "Expected exactly one settings block")

	// Test the helper functions directly
	parsedContent, err := parseBlockContent(content.Blocks[0])
	require.NoError(t, err)

	settings := &spookytypesproject.Settings{}
	err = populateSettings(settings, parsedContent)

	assert.NoError(t, err)
	assert.NotNil(t, settings)
	// Should have default values
	assert.Equal(t, 0, settings.ParallelWorkers)
	assert.Equal(t, 0, settings.TimeoutSeconds)
	assert.Equal(t, "", settings.LogLevel)
	assert.False(t, settings.DefaultDryRun)
	assert.False(t, settings.ValidateBeforeRun)
	assert.Equal(t, 0, settings.MaxRetries)
	assert.Equal(t, 0, settings.RetryDelaySeconds)
}

func TestParseSettingsBlockPartial(t *testing.T) {
	// Test with only some attributes set
	hclContent := `
	settings {
		parallel_workers = 8
		log_level = "warn"
		default_dry_run = true
	}
	`

	parser := hclparse.NewParser()
	file, diags := parser.ParseHCL([]byte(hclContent), "test.hcl")
	require.False(t, diags.HasErrors(), "Failed to parse HCL: %v", diags)

	body := file.Body
	content, diags := body.Content(&hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{
			{Type: "settings"},
		},
	})
	require.False(t, diags.HasErrors(), "Failed to get body content: %v", diags)
	require.Len(t, content.Blocks, 1, "Expected exactly one settings block")

	// Test the helper functions directly
	parsedContent, err := parseBlockContent(content.Blocks[0])
	require.NoError(t, err)

	settings := &spookytypesproject.Settings{}
	err = populateSettings(settings, parsedContent)

	assert.NoError(t, err)
	assert.NotNil(t, settings)
	// Set values
	assert.Equal(t, 8, settings.ParallelWorkers)
	assert.Equal(t, "warn", settings.LogLevel)
	assert.True(t, settings.DefaultDryRun)
	// Default values for unset attributes
	assert.Equal(t, 0, settings.TimeoutSeconds)
	assert.False(t, settings.ValidateBeforeRun)
	assert.Equal(t, 0, settings.MaxRetries)
	assert.Equal(t, 0, settings.RetryDelaySeconds)
}
