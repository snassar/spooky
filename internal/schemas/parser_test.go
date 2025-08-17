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

// Helper functions

func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}

func float64Ptr(f float64) *float64 {
	return &f
}

func assertFieldValidationEqual(t *testing.T, expected, actual *spookytypesschemas.FieldValidation, fieldPath string) {
	assert.Equal(t, expected.Type, actual.Type, "Field %s type mismatch", fieldPath)
	assert.Equal(t, expected.Required, actual.Required, "Field %s required mismatch", fieldPath)
	assert.Equal(t, expected.Description, actual.Description, "Field %s description mismatch", fieldPath)
	assert.Equal(t, expected.Default, actual.Default, "Field %s default mismatch", fieldPath)

	if expected.Constraints != nil {
		assert.NotNil(t, actual.Constraints, "Field %s constraints should not be nil", fieldPath)
		if actual.Constraints != nil {
			assert.Equal(t, expected.Constraints.MinLength, actual.Constraints.MinLength, "Field %s min_length mismatch", fieldPath)
			assert.Equal(t, expected.Constraints.MaxLength, actual.Constraints.MaxLength, "Field %s max_length mismatch", fieldPath)
			assert.Equal(t, expected.Constraints.Pattern, actual.Constraints.Pattern, "Field %s pattern mismatch", fieldPath)
			assert.Equal(t, expected.Constraints.Format, actual.Constraints.Format, "Field %s format mismatch", fieldPath)
			assert.Equal(t, expected.Constraints.Min, actual.Constraints.Min, "Field %s min mismatch", fieldPath)
			assert.Equal(t, expected.Constraints.Max, actual.Constraints.Max, "Field %s max mismatch", fieldPath)
			assert.Equal(t, expected.Constraints.Enum, actual.Constraints.Enum, "Field %s enum mismatch", fieldPath)
		}
	} else {
		assert.Nil(t, actual.Constraints, "Field %s constraints should be nil", fieldPath)
	}
}
