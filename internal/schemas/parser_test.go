package schemas

import (
	"testing"

	"github.com/stretchr/testify/assert"

	spookytypesschemas "spooky/internal/types/schemas"
)

func TestSchemaParser_EmptyContent(t *testing.T) {
	logger := &MockLogger{}
	parser := NewSchemaParser(logger)

	schema := &spookytypesschemas.Schema{
		Version:     "1.0",
		Type:        "hcl",
		Name:        "empty-schema",
		Description: "Empty schema",
		Content:     "",
	}

	err := parser.ParseValidationRules(schema)
	assert.NoError(t, err)
	assert.Nil(t, schema.Validation)
}

func TestSchemaParser_InvalidHCL(t *testing.T) {
	logger := &MockLogger{}
	parser := NewSchemaParser(logger)

	schema := &spookytypesschemas.Schema{
		Version:     "1.0",
		Type:        "hcl",
		Name:        "invalid-schema",
		Description: "Invalid schema",
		Content:     `project { name = { type = "string" required = true }`, // Missing closing brace
	}

	err := parser.ParseValidationRules(schema)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse HCL content")
}
