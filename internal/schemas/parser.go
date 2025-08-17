// Package schemas provides schema parsing functionality for the spooky codebase.
package schemas

import (
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"time"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/zclconf/go-cty/cty"

	spookytypeslogging "spooky/internal/types/logging"
	spookytypesschemas "spooky/internal/types/schemas"
)

// SchemaParser provides functionality to parse HCL schema content into structured validation rules
type SchemaParser struct {
	logger spookytypeslogging.Logger
	parser *hclparse.Parser
}

// NewSchemaParser creates a new schema parser instance
func NewSchemaParser(logger spookytypeslogging.Logger) *SchemaParser {
	return &SchemaParser{
		logger: logger,
		parser: hclparse.NewParser(),
	}
}

// ParseValidationRules parses HCL schema content and populates the Validation.Fields structure
func (p *SchemaParser) ParseValidationRules(schema *spookytypesschemas.Schema) error {
	if schema == nil {
		return fmt.Errorf("schema cannot be nil")
	}

	if schema.Content == "" {
		p.logger.Debug("Schema has no content to parse", map[string]interface{}{
			"schema_name": schema.Name,
		})
		return nil
	}

	p.logger.Debug("Parsing validation rules from schema content", map[string]interface{}{
		"schema_name":    schema.Name,
		"content_length": len(schema.Content),
	})

	// Parse HCL content
	file, diags := p.parser.ParseHCL([]byte(schema.Content), schema.Name)
	if diags.HasErrors() {
		return fmt.Errorf("failed to parse HCL content: %s", diags.Error())
	}

	// Initialize validation structure if not present
	if schema.Validation == nil {
		schema.Validation = &spookytypesschemas.SchemaValidation{
			Enabled: true,
			Mode:    "strict",
			Fields:  make(map[string]*spookytypesschemas.FieldValidation),
		}
	}

	if schema.Validation.Fields == nil {
		schema.Validation.Fields = make(map[string]*spookytypesschemas.FieldValidation)
	}

	// Parse blocks to extract field validation rules
	if err := p.parseBlocks(file.Body, schema); err != nil {
		return fmt.Errorf("failed to parse blocks: %w", err)
	}

	p.logger.Info("Successfully parsed validation rules", map[string]interface{}{
		"schema_name":   schema.Name,
		"fields_parsed": len(schema.Validation.Fields),
	})

	return nil
}

// parseBlocks recursively parses HCL blocks to extract field validation rules
func (p *SchemaParser) parseBlocks(body hcl.Body, schema *spookytypesschemas.Schema) error {
	// Define allowed block types for validation rules
	bodySchema := &hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{
			{Type: "template_metadata"},
			{Type: "project"},
			{Type: "machines"},
			{Type: "actions"},
			{Type: "variables"},
			{Type: "templates"},
			{Type: "logging"},
			{Type: "facts"},
			{Type: "spooky"},
			{Type: "validation_rules"},
		},
	}

	// Parse the body as HCL syntax with allowed blocks
	content, diags := body.Content(bodySchema)
	if diags.HasErrors() {
		return fmt.Errorf("failed to parse body content: %s", diags.Error())
	}

	// Process each block
	for _, block := range content.Blocks {
		if err := p.parseBlock(block, schema); err != nil {
			return fmt.Errorf("failed to parse block %s: %w", block.Type, err)
		}
	}

	return nil
}

// parseBlock parses a single HCL block to extract field validation rules
func (p *SchemaParser) parseBlock(block *hcl.Block, schema *spookytypesschemas.Schema) error {
	blockName := block.Type

	// For validation_rules blocks, use completely permissive parsing
	var content *hcl.BodyContent
	var diags hcl.Diagnostics

	if blockName == "validation_rules" {
		// Use completely permissive parsing for validation_rules
		content, _, diags = block.Body.PartialContent(&hcl.BodySchema{})
	} else {
		// Use normal parsing for other blocks
		content, diags = block.Body.Content(&hcl.BodySchema{
			Attributes: []hcl.AttributeSchema{
				{Name: "*", Required: false}, // Allow any attribute
			},
			Blocks: []hcl.BlockHeaderSchema{
				{Type: "rule"}, // Allow rule blocks
				{Type: "*"},    // Allow any other blocks
			},
		})
	}

	if diags.HasErrors() {
		return fmt.Errorf("failed to parse block %s content: %s", blockName, diags.Error())
	}

	// Process each attribute in the block
	for _, attr := range content.Attributes {
		fieldPath := fmt.Sprintf("%s.%s", blockName, attr.Name)

		// Parse the attribute value as a field validation rule
		if fieldValidation, err := p.parseFieldValidation(attr); err != nil {
			p.logger.Warn("Failed to parse field validation", map[string]interface{}{
				"field_path": fieldPath,
				"error":      err.Error(),
			})
			continue
		} else if fieldValidation != nil {
			schema.Validation.Fields[fieldPath] = fieldValidation
			p.logger.Debug("Parsed field validation", map[string]interface{}{
				"field_path": fieldPath,
				"type":       fieldValidation.Type,
				"required":   fieldValidation.Required,
			})
		}
	}

	// Recursively parse nested blocks
	for _, nestedBlock := range content.Blocks {
		nestedFieldPath := fmt.Sprintf("%s.%s", blockName, nestedBlock.Type)

		// Parse nested block attributes - allow any attributes
		nestedContent, diags := nestedBlock.Body.Content(&hcl.BodySchema{
			Attributes: []hcl.AttributeSchema{
				{Name: "*", Required: false}, // Allow any attribute
			},
		})
		if diags.HasErrors() {
			p.logger.Warn("Failed to parse nested block content", map[string]interface{}{
				"nested_field_path": nestedFieldPath,
				"error":             diags.Error(),
			})
			continue
		}

		// Process nested attributes
		for _, nestedAttr := range nestedContent.Attributes {
			fieldPath := fmt.Sprintf("%s.%s", nestedFieldPath, nestedAttr.Name)

			if fieldValidation, err := p.parseFieldValidation(nestedAttr); err != nil {
				p.logger.Warn("Failed to parse nested field validation", map[string]interface{}{
					"field_path": fieldPath,
					"error":      err.Error(),
				})
				continue
			} else if fieldValidation != nil {
				schema.Validation.Fields[fieldPath] = fieldValidation
				p.logger.Debug("Parsed nested field validation", map[string]interface{}{
					"field_path": fieldPath,
					"type":       fieldValidation.Type,
					"required":   fieldValidation.Required,
				})
			}
		}

		// Handle deeply nested structures (like properties in objects)
		if err := p.parseNestedProperties(nestedBlock, nestedFieldPath, schema); err != nil {
			p.logger.Warn("Failed to parse nested properties", map[string]interface{}{
				"nested_field_path": nestedFieldPath,
				"error":             err.Error(),
			})
		}
	}

	return nil
}

// parseNestedProperties handles deeply nested property structures
func (p *SchemaParser) parseNestedProperties(block *hcl.Block, basePath string, schema *spookytypesschemas.Schema) error {
	content, diags := block.Body.Content(&hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{
			{Type: "properties"},
		},
	})
	if diags.HasErrors() {
		return fmt.Errorf("failed to parse nested properties: %s", diags.Error())
	}

	for _, propertiesBlock := range content.Blocks {
		if propertiesBlock.Type == "properties" {
			propertiesContent, diags := propertiesBlock.Body.Content(&hcl.BodySchema{})
			if diags.HasErrors() {
				return fmt.Errorf("failed to parse properties content: %s", diags.Error())
			}

			for _, attr := range propertiesContent.Attributes {
				fieldPath := fmt.Sprintf("%s.%s", basePath, attr.Name)

				if fieldValidation, err := p.parseFieldValidation(attr); err != nil {
					p.logger.Warn("Failed to parse property field validation", map[string]interface{}{
						"field_path": fieldPath,
						"error":      err.Error(),
					})
					continue
				} else if fieldValidation != nil {
					schema.Validation.Fields[fieldPath] = fieldValidation
					p.logger.Debug("Parsed property field validation", map[string]interface{}{
						"field_path": fieldPath,
						"type":       fieldValidation.Type,
						"required":   fieldValidation.Required,
					})
				}
			}
		}
	}

	return nil
}

// parseFieldValidation parses an HCL attribute into a FieldValidation structure
func (p *SchemaParser) parseFieldValidation(attr *hcl.Attribute) (*spookytypesschemas.FieldValidation, error) {
	// Try to parse as a block (field validation rule)
	if block, ok := attr.Expr.(*hclsyntax.ObjectConsExpr); ok {
		return p.parseFieldValidationBlock(block)
	}

	// Try to parse as a simple value
	if val, diags := attr.Expr.Value(nil); !diags.HasErrors() {
		return p.parseFieldValidationValue(val)
	}

	return nil, fmt.Errorf("unable to parse field validation for attribute %s", attr.Name)
}

// parseFieldValidationBlock parses a field validation block
func (p *SchemaParser) parseFieldValidationBlock(block *hclsyntax.ObjectConsExpr) (*spookytypesschemas.FieldValidation, error) {
	fieldValidation := &spookytypesschemas.FieldValidation{
		Constraints: &spookytypesschemas.FieldConstraints{},
	}

	// Parse each key-value pair in the block
	for _, item := range block.Items {
		key, diags := item.KeyExpr.Value(nil)
		if diags.HasErrors() {
			continue
		}

		value, diags := item.ValueExpr.Value(nil)
		if diags.HasErrors() {
			continue
		}

		keyStr := key.AsString()
		switch keyStr {
		case "type":
			fieldValidation.Type = value.AsString()
		case "required":
			fieldValidation.Required = value.True()
		case "description":
			fieldValidation.Description = value.AsString()
		case "default":
			fieldValidation.Default = p.convertCtyValue(value)
		case "min_length":
			if value.Type() == cty.Number {
				if intVal, acc := value.AsBigFloat().Int64(); acc == big.Exact {
					minLength := int(intVal)
					fieldValidation.Constraints.MinLength = &minLength
				}
			}
		case "max_length":
			if value.Type() == cty.Number {
				if intVal, acc := value.AsBigFloat().Int64(); acc == big.Exact {
					maxLength := int(intVal)
					fieldValidation.Constraints.MaxLength = &maxLength
				}
			}
		case "pattern":
			pattern := value.AsString()
			fieldValidation.Constraints.Pattern = &pattern
		case "format":
			format := value.AsString()
			fieldValidation.Constraints.Format = &format
		case "min":
			if value.Type() == cty.Number {
				if floatVal, acc := value.AsBigFloat().Float64(); acc == big.Exact || acc == big.Below {
					fieldValidation.Constraints.Min = &floatVal
				}
			}
		case "max":
			if value.Type() == cty.Number {
				if floatVal, acc := value.AsBigFloat().Float64(); acc == big.Exact || acc == big.Above {
					fieldValidation.Constraints.Max = &floatVal
				}
			}
		case "enum":
			if value.Type().IsListType() {
				var enumValues []interface{}
				for _, enumVal := range value.AsValueSlice() {
					enumValues = append(enumValues, p.convertCtyValue(enumVal))
				}
				fieldValidation.Constraints.Enum = enumValues
			}
		}
	}

	// Set default type if not specified
	if fieldValidation.Type == "" {
		fieldValidation.Type = "string"
	}

	return fieldValidation, nil
}

// parseFieldValidationValue parses a simple field validation value
func (p *SchemaParser) parseFieldValidationValue(val cty.Value) (*spookytypesschemas.FieldValidation, error) {
	fieldValidation := &spookytypesschemas.FieldValidation{
		Type:        p.inferTypeFromValue(val),
		Required:    false,
		Constraints: &spookytypesschemas.FieldConstraints{},
	}

	return fieldValidation, nil
}

// convertCtyValue converts a cty.Value to a Go interface{}
func (p *SchemaParser) convertCtyValue(val cty.Value) interface{} {
	switch {
	case val.Type() == cty.String:
		return val.AsString()
	case val.Type() == cty.Number:
		if floatVal, acc := val.AsBigFloat().Float64(); acc == big.Exact || acc == big.Below {
			return floatVal
		}
		return val.AsBigFloat().String()
	case val.Type() == cty.Bool:
		return val.True()
	case val.Type().IsListType():
		var result []interface{}
		for _, item := range val.AsValueSlice() {
			result = append(result, p.convertCtyValue(item))
		}
		return result
	case val.Type().IsMapType():
		result := make(map[string]interface{})
		for key, value := range val.AsValueMap() {
			result[key] = p.convertCtyValue(value)
		}
		return result
	default:
		return val.GoString()
	}
}

// inferTypeFromValue infers the field type from a cty.Value
func (p *SchemaParser) inferTypeFromValue(val cty.Value) string {
	switch {
	case val.Type() == cty.String:
		return "string"
	case val.Type() == cty.Number:
		return "number"
	case val.Type() == cty.Bool:
		return "boolean"
	case val.Type().IsListType():
		return "array"
	case val.Type().IsMapType():
		return "object"
	default:
		return "string"
	}
}

// ParseValidationRulesFromFile parses validation rules from a schema file
func (p *SchemaParser) ParseValidationRulesFromFile(filePath string) (*spookytypesschemas.Schema, error) {
	// Read file content
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read schema file: %w", err)
	}

	// Create schema
	schema := &spookytypesschemas.Schema{
		Version:     "1.0",
		Type:        "hcl",
		Name:        filepath.Base(filePath),
		Description: fmt.Sprintf("Schema loaded from %s", filePath),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Content:     string(data),
		Metadata:    make(map[string]interface{}),
	}

	// Parse validation rules
	if err := p.ParseValidationRules(schema); err != nil {
		return nil, fmt.Errorf("failed to parse validation rules: %w", err)
	}

	return schema, nil
}
