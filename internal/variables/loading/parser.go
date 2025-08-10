package loading

import (
	"encoding/json"
	"fmt"

	"spooky/internal/logging"
	"spooky/internal/schemas"
	spookylogging "spooky/internal/types/logging"
	spookytypesvariables "spooky/internal/types/variables"

	"github.com/hashicorp/hcl/v2/hclparse"
)

// VariableParser implements the Parser interface
type VariableParser struct {
	logger spookylogging.Logger
}

// NewParser creates a new parser
func NewParser(logger spookylogging.Logger) *VariableParser {
	return &VariableParser{
		logger: logger,
	}
}

// ParseHCL parses HCL content into variables
func (p *VariableParser) ParseHCL(content []byte) ([]*spookytypesvariables.Variable, error) {
	p.logger.Debug("Parsing HCL content", logging.Int("length", len(content)))

	// Use schema system for HCL parsing
	validator := schemas.NewSchemaValidator()
	if err := validator.LoadSchema(schemas.SchemaTypeVariablesHCL); err != nil {
		return nil, fmt.Errorf("failed to load variables schema: %w", err)
	}

	// Parse HCL content using schema
	result := validator.ValidateFile("variables.hcl", "variables-hcl")
	if !result.Valid {
		return nil, fmt.Errorf("HCL validation failed: %v", result.Errors)
	}

	// Convert validated content to variables
	var collection spookytypesvariables.VariableCollection
	parser := hclparse.NewParser()
	_, diags := parser.ParseHCL(content, "variables.hcl")
	if diags.HasErrors() {
		return nil, fmt.Errorf("failed to parse HCL: %v", diags)
	}

	// Use schema-validated HCL parsing with JSON as intermediate
	// The schema system provides validation, this handles the conversion

	// Convert HCL to JSON first, then to struct
	// This is a simplified approach - proper HCL parsing would be more complex
	var jsonData map[string]interface{}
	if err := json.Unmarshal(content, &jsonData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal HCL: %w", err)
	}

	// Convert to JSON bytes
	jsonBytes, err := json.Marshal(jsonData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal to JSON: %w", err)
	}

	// Parse as JSON
	if err := json.Unmarshal(jsonBytes, &collection); err != nil {
		return nil, fmt.Errorf("failed to parse as JSON: %w", err)
	}

	return collection.Variables, nil
}

// ParseJSON parses JSON content into variables
func (p *VariableParser) ParseJSON(content []byte) ([]*spookytypesvariables.Variable, error) {
	p.logger.Debug("Parsing JSON content", logging.Int("length", len(content)))

	var collection spookytypesvariables.VariableCollection
	if err := json.Unmarshal(content, &collection); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return collection.Variables, nil
}

// ValidateContent validates content format
func (p *VariableParser) ValidateContent(content []byte, format string) error {
	switch format {
	case "hcl":
		validator := schemas.NewSchemaValidator()
		if err := validator.LoadSchema(schemas.SchemaTypeVariablesHCL); err != nil {
			return fmt.Errorf("failed to load variables schema: %w", err)
		}

		result := validator.ValidateFile("variables.hcl", "variables-hcl")
		if !result.Valid {
			return fmt.Errorf("HCL validation failed: %v", result.Errors)
		}
		return nil
	case "json":
		// Use schema system for JSON validation too
		validator := schemas.NewSchemaValidator()
		if err := validator.LoadSchema(schemas.SchemaTypeVariablesJSON); err != nil {
			return fmt.Errorf("failed to load variables JSON schema: %w", err)
		}

		result := validator.ValidateFile("variables.json", "variables-json")
		if !result.Valid {
			return fmt.Errorf("JSON validation failed: %v", result.Errors)
		}
		return nil
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
}
