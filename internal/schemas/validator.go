package schemas

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/pkg/errors"
)

// SchemaValidator provides functionality to validate HCL content against embedded schemas
type SchemaValidator struct {
	embedder *SchemaEmbedder
	logger   *slog.Logger
}

// ValidationResult contains the result of schema validation
type ValidationResult struct {
	IsValid    bool
	Errors     []ValidationError
	Warnings   []ValidationWarning
	SchemaName string
}

// ValidationError represents a schema validation error
type ValidationError struct {
	Field    string
	Value    interface{}
	Message  string
	Severity string
}

// ValidationWarning represents a schema validation warning
type ValidationWarning struct {
	Field   string
	Value   interface{}
	Message string
}

// NewSchemaValidator creates a new schema validator
func NewSchemaValidator() (*SchemaValidator, error) {
	embedder, err := NewSchemaEmbedder()
	if err != nil {
		return nil, errors.Wrap(err, "failed to create schema embedder")
	}

	return &SchemaValidator{
		embedder: embedder,
		logger:   slog.Default(),
	}, nil
}

// ValidateContent validates HCL content against a specific schema
func (sv *SchemaValidator) ValidateContent(schemaName, content string) (*ValidationResult, error) {
	result := &ValidationResult{
		IsValid:    true,
		Errors:     []ValidationError{},
		Warnings:   []ValidationWarning{},
		SchemaName: schemaName,
	}

	// First, validate basic HCL syntax
	_, diags := hclsyntax.ParseConfig([]byte(content), schemaName+".hcl", hcl.Pos{Line: 1, Column: 1})
	if diags.HasErrors() {
		result.IsValid = false
		for _, diag := range diags.Errs() {
			result.Errors = append(result.Errors, ValidationError{
				Message:  diag.Error(),
				Severity: "error",
			})
		}
		return result, nil
	}

	// Get the actual schema to understand what's required
	schema, exists := sv.embedder.GetSchema(schemaName)
	if !exists {
		return nil, NewSchemaNotFoundError(schemaName, sv.embedder.ListSchemas())
	}

	// Parse the schema to understand the structure
	schemaData, err := sv.ParseHCLContent(schema)
	if err != nil {
		result.IsValid = false
		result.Errors = append(result.Errors, ValidationError{
			Message:  fmt.Sprintf("Failed to parse schema: %v", err),
			Severity: "error",
		})
		return result, nil
	}

	// Parse the content to validate
	contentData, err := sv.ParseHCLContent(content)
	if err != nil {
		result.IsValid = false
		result.Errors = append(result.Errors, ValidationError{
			Message:  fmt.Sprintf("Failed to parse content: %v", err),
			Severity: "error",
		})
		return result, nil
	}

	// Validate the content against the schema structure
	if err := sv.validateAgainstSchemaStructure(schemaName, contentData, schemaData, result); err != nil {
		result.IsValid = false
		result.Errors = append(result.Errors, ValidationError{
			Message:  err.Error(),
			Severity: "error",
		})
	}

	// Update validity based on errors
	if len(result.Errors) > 0 {
		result.IsValid = false
	}

	return result, nil
}

// ParseHCLContent parses HCL content into a structured map
func (sv *SchemaValidator) ParseHCLContent(content string) (map[string]interface{}, error) {
	// Parse HCL into AST
	file, diags := hclsyntax.ParseConfig([]byte(content), "content.hcl", hcl.Pos{Line: 1, Column: 1})

	// Check for any parsing errors
	if diags.HasErrors() {
		var errorMsgs []string
		for _, diag := range diags.Errs() {
			errorMsgs = append(errorMsgs, diag.Error())
		}
		return nil, fmt.Errorf("HCL parsing failed: %s", strings.Join(errorMsgs, "; "))
	}

	// Convert AST to structured data
	data := make(map[string]interface{})

	for _, block := range file.Body.(*hclsyntax.Body).Blocks {
		blockData := sv.parseBlock(block)
		data[block.Type] = blockData
	}

	return data, nil
}

// parseBlock recursively parses HCL blocks into structured data
func (sv *SchemaValidator) parseBlock(block *hclsyntax.Block) map[string]interface{} {
	result := make(map[string]interface{})

	// Parse attributes
	for name, attr := range block.Body.Attributes {
		value := sv.parseAttributeValue(attr.Expr)
		result[name] = value
	}

	// Parse nested blocks
	for _, nestedBlock := range block.Body.Blocks {
		nestedData := sv.parseBlock(nestedBlock)

		// Handle blocks with labels (like "action "name" {")
		if len(nestedBlock.Labels) > 0 {
			label := nestedBlock.Labels[0]
			if existing, exists := result[label]; exists {
				// Convert to array if multiple blocks with same label
				if arr, ok := existing.([]map[string]interface{}); ok {
					result[label] = append(arr, nestedData)
				} else {
					result[label] = []map[string]interface{}{existing.(map[string]interface{}), nestedData}
				}
			} else {
				result[label] = nestedData
			}
		} else {
			// Handle blocks without labels
			// If we already have a block of this type, convert to array
			if existing, exists := result[nestedBlock.Type]; exists {
				if arr, ok := existing.([]map[string]interface{}); ok {
					result[nestedBlock.Type] = append(arr, nestedData)
				} else {
					result[nestedBlock.Type] = []map[string]interface{}{existing.(map[string]interface{}), nestedData}
				}
			} else {
				result[nestedBlock.Type] = nestedData
			}
		}
	}

	return result
}

// parseAttributeValue converts HCL expression to Go value
func (sv *SchemaValidator) parseAttributeValue(expr hclsyntax.Expression) interface{} {
	switch v := expr.(type) {
	case *hclsyntax.LiteralValueExpr:
		return v.Val
	case *hclsyntax.TemplateExpr:
		// Handle template expressions - try to extract literal values
		if len(v.Parts) == 1 {
			// Single part - check if it's a literal value
			if lit, ok := v.Parts[0].(*hclsyntax.LiteralValueExpr); ok {
				return lit.Val
			}
		} else if len(v.Parts) == 0 {
			// Empty template - return empty string
			return ""
		}

		// For complex templates, try to extract the first literal value
		for _, part := range v.Parts {
			if lit, ok := part.(*hclsyntax.LiteralValueExpr); ok {
				return lit.Val
			}
		}

		// For complex templates, return a placeholder
		return "complex_expression"
	case *hclsyntax.ObjectConsExpr:
		result := make(map[string]interface{})
		for _, item := range v.Items {
			if key, ok := item.KeyExpr.(*hclsyntax.LiteralValueExpr); ok {
				keyStr := key.Val.AsString()
				result[keyStr] = sv.parseAttributeValue(item.ValueExpr)
			}
		}
		return result
	case *hclsyntax.TupleConsExpr:
		var result []interface{}
		for _, item := range v.Exprs {
			result = append(result, sv.parseAttributeValue(item))
		}
		return result
	default:
		// For complex expressions, return a placeholder
		return "complex_expression"
	}
}

// validateAgainstSchemaStructure validates content against the actual schema structure
func (sv *SchemaValidator) validateAgainstSchemaStructure(schemaName string, contentData, schemaData map[string]interface{}, result *ValidationResult) error {
	// Extract required blocks from the schema
	requiredBlocks := sv.extractRequiredBlocks(schemaData)

	// Validate that all required blocks exist in content
	for _, blockName := range requiredBlocks {
		if _, exists := contentData[blockName]; !exists {
			result.Errors = append(result.Errors, ValidationError{
				Field:    blockName,
				Message:  fmt.Sprintf("missing required block: %s", blockName),
				Severity: "error",
			})
		}
	}

	// Validate each block's structure against the schema
	for blockName, blockData := range contentData {
		if schemaBlock, exists := schemaData[blockName]; exists {
			if err := sv.validateBlockStructure(blockName, blockData, schemaBlock, result); err != nil {
				result.Errors = append(result.Errors, ValidationError{
					Field:    blockName,
					Message:  err.Error(),
					Severity: "error",
				})
			}
		}
	}

	return nil
}

// extractRequiredBlocks extracts the names of required blocks from the schema
func (sv *SchemaValidator) extractRequiredBlocks(schemaData map[string]interface{}) []string {
	var requiredBlocks []string

	// Look for top-level blocks in the schema - these are the required ones
	for blockName, blockData := range schemaData {
		// Skip internal schema metadata
		if blockName == "schema" || blockName == "metadata" {
			continue
		}

		// Only treat blocks as required if they don't look like schema definition blocks
		// Schema definition blocks contain attributes like "type", "required", "description", etc.
		if blockMap, ok := blockData.(map[string]interface{}); ok {
			if !sv.isSchemaDefinitionBlock(blockMap) {
				requiredBlocks = append(requiredBlocks, blockName)
			}
		}
	}

	return requiredBlocks
}

// validateBlockStructure validates a block's structure against its schema definition
func (sv *SchemaValidator) validateBlockStructure(blockName string, blockData, schemaBlock interface{}, result *ValidationResult) error {
	// Convert blockData to map if it's not already
	var blockMap map[string]interface{}
	switch v := blockData.(type) {
	case map[string]interface{}:
		blockMap = v
	case []map[string]interface{}:
		// Handle array of blocks
		for i, item := range v {
			if err := sv.validateBlockStructure(fmt.Sprintf("%s[%d]", blockName, i), item, schemaBlock, result); err != nil {
				return err
			}
		}
		return nil
	case string:
		// Skip string values (likely schema content)
		return nil
	default:
		// Skip unknown types for now
		return nil
	}

	// Convert schemaBlock to map if it's not already
	var schemaMap map[string]interface{}
	switch v := schemaBlock.(type) {
	case map[string]interface{}:
		schemaMap = v
	default:
		// Skip if schema block is not a map
		return nil
	}

	// Skip validation if this looks like a schema definition block
	if sv.isSchemaDefinitionBlock(blockMap) {
		return nil
	}

	// Validate required attributes from schema
	if requiredAttrs := sv.extractRequiredAttributes(schemaMap); len(requiredAttrs) > 0 {
		for _, attrName := range requiredAttrs {
			if _, exists := blockMap[attrName]; !exists {
				result.Errors = append(result.Errors, ValidationError{
					Field:    fmt.Sprintf("%s.%s", blockName, attrName),
					Message:  fmt.Sprintf("missing required attribute: %s", attrName),
					Severity: "error",
				})
			}
		}
	}

	// Validate nested blocks
	for nestedBlockName, nestedBlockData := range blockMap {
		if schemaNestedBlock, exists := schemaMap[nestedBlockName]; exists {
			if err := sv.validateBlockStructure(
				fmt.Sprintf("%s.%s", blockName, nestedBlockName),
				nestedBlockData,
				schemaNestedBlock,
				result,
			); err != nil {
				return err
			}
		}
	}

	return nil
}

// isSchemaDefinitionBlock checks if a block appears to be a schema definition
func (sv *SchemaValidator) isSchemaDefinitionBlock(blockData map[string]interface{}) bool {
	// Check if this looks like a schema definition by looking for schema-specific keys
	schemaKeys := []string{"type", "required", "pattern", "min_length", "max_length", "min", "max", "default", "description"}

	for _, key := range schemaKeys {
		if _, exists := blockData[key]; exists {
			return true
		}
	}

	return false
}

// extractRequiredAttributes extracts required attribute names from a schema block
func (sv *SchemaValidator) extractRequiredAttributes(schemaBlock map[string]interface{}) []string {
	var requiredAttrs []string

	// Look for attributes in the schema block
	for attrName, attrData := range schemaBlock {
		// Skip nested blocks
		if _, isMap := attrData.(map[string]interface{}); isMap {
			continue
		}
		if _, isArray := attrData.([]interface{}); isArray {
			continue
		}

		// For now, assume all attributes are required
		// This could be enhanced to read actual required/optional markers from schema
		requiredAttrs = append(requiredAttrs, attrName)
	}

	return requiredAttrs
}
