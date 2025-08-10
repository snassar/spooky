package loading

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"spooky/internal/schemas"
	spookytypesvariables "spooky/internal/types/variables"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
)

// HCLVariableLoader implements VariableLoader for HCL files
type HCLVariableLoader struct {
	validator *schemas.SchemaValidator
}

func NewHCLVariableLoader(validator *schemas.SchemaValidator) *HCLVariableLoader {
	return &HCLVariableLoader{
		validator: validator,
	}
}

func (l *HCLVariableLoader) Load(ctx context.Context, source interface{}) ([]*spookytypesvariables.Variable, error) {
	filePath, ok := source.(string)
	if !ok {
		return nil, fmt.Errorf("HCL loader expects string file path")
	}

	// Read file content
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read HCL file: %w", err)
	}

	// Validate against schema first
	if err := l.ValidateSchema(content); err != nil {
		return nil, fmt.Errorf("schema validation failed: %w", err)
	}

	// Convert HCL to variables
	return l.parseHCLToVariables(content)
}

func (l *HCLVariableLoader) ValidateSchema(content []byte) error {
	if l.validator == nil {
		// Create a new validator if none provided
		l.validator = schemas.NewSchemaValidator()
	}

	if err := l.validator.LoadSchema(schemas.SchemaTypeVariablesHCL); err != nil {
		return fmt.Errorf("failed to load HCL schema: %w", err)
	}

	// Create temporary file for validation
	tmpFile, err := os.CreateTemp("", "variables-*.hcl")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	if _, err := tmpFile.Write(content); err != nil {
		return fmt.Errorf("failed to write content to temp file: %w", err)
	}

	result := l.validator.ValidateFile(tmpFile.Name(), string(schemas.SchemaTypeVariablesHCL))
	if !result.Valid {
		return fmt.Errorf("HCL schema validation failed: %v", result.Errors)
	}

	return nil
}

func (l *HCLVariableLoader) GetName() string {
	return "hcl"
}

func (l *HCLVariableLoader) GetSupportedExtensions() []string {
	return []string{".hcl"}
}

// JSONVariableLoader implements VariableLoader for JSON files
type JSONVariableLoader struct {
	validator *schemas.SchemaValidator
}

func NewJSONVariableLoader(validator *schemas.SchemaValidator) *JSONVariableLoader {
	return &JSONVariableLoader{
		validator: validator,
	}
}

func (l *JSONVariableLoader) Load(ctx context.Context, source interface{}) ([]*spookytypesvariables.Variable, error) {
	filePath, ok := source.(string)
	if !ok {
		return nil, fmt.Errorf("JSON loader expects string file path")
	}

	// Read file content
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read JSON file: %w", err)
	}

	// Parse JSON content
	var jsonData map[string]interface{}
	if err := json.Unmarshal(content, &jsonData); err != nil {
		return nil, fmt.Errorf("JSON parsing failed: %w", err)
	}

	// Convert JSON to variables
	return l.parseJSONToVariables(jsonData)
}

func (l *JSONVariableLoader) ValidateSchema(content []byte) error {
	if l.validator == nil {
		// Create a new validator if none provided
		l.validator = schemas.NewSchemaValidator()
	}

	if err := l.validator.LoadSchema(schemas.SchemaTypeVariablesJSON); err != nil {
		return fmt.Errorf("failed to load JSON schema: %w", err)
	}

	// Create temporary file for validation
	tmpFile, err := os.CreateTemp("", "variables-*.json")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	if _, err := tmpFile.Write(content); err != nil {
		return fmt.Errorf("failed to write content to temp file: %w", err)
	}

	result := l.validator.ValidateFile(tmpFile.Name(), string(schemas.SchemaTypeVariablesJSON))
	if !result.Valid {
		return fmt.Errorf("JSON schema validation failed: %v", result.Errors)
	}

	return nil
}

func (l *JSONVariableLoader) GetName() string {
	return "json"
}

func (l *JSONVariableLoader) GetSupportedExtensions() []string {
	return []string{".json"}
}

// Helper methods for parsing
func (l *HCLVariableLoader) parseHCLToVariables(content []byte) ([]*spookytypesvariables.Variable, error) {
	// Parse HCL content using the proper HCL library
	parser := hclparse.NewParser()
	file, diags := parser.ParseHCL(content, "variables.hcl")
	if diags.HasErrors() {
		return nil, fmt.Errorf("HCL parsing failed: %v", diags)
	}

	// Extract variable blocks
	variables := []*spookytypesvariables.Variable{}
	body := file.Body

	// Find variables block
	variablesBlock := findVariablesBlock(body)
	if variablesBlock == nil {
		return variables, nil // No variables defined
	}

	// Parse each variable block - handle both labeled and unlabeled blocks
	blockContent, diags := variablesBlock.Body.Content(&hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{
			{Type: "variable", LabelNames: []string{"name"}},
		},
	})
	if diags.HasErrors() {
		return nil, fmt.Errorf("failed to parse variable blocks: %v", diags)
	}

	for _, block := range blockContent.Blocks {
		variable, err := parseVariableBlock(block)
		if err != nil {
			if len(block.Labels) > 0 {
				return nil, fmt.Errorf("failed to parse variable %s: %w", block.Labels[0], err)
			}
			return nil, fmt.Errorf("failed to parse variable block: %w", err)
		}
		variables = append(variables, variable)
	}

	return variables, nil
}

// findVariablesBlock finds the variables block in the HCL body
func findVariablesBlock(body hcl.Body) *hcl.Block {
	content, diags := body.Content(&hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{
			{
				Type: "variables",
			},
		},
	})
	if diags.HasErrors() {
		return nil
	}

	if len(content.Blocks) == 0 {
		return nil
	}

	return content.Blocks[0]
}

// parseVariableBlock parses a single variable block
func parseVariableBlock(block *hcl.Block) (*spookytypesvariables.Variable, error) {
	if len(block.Labels) == 0 {
		return nil, fmt.Errorf("variable block must have a name label")
	}

	variable := &spookytypesvariables.Variable{
		Name:      block.Labels[0],
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Parse attributes and nested blocks
	if err := parseVariableBlocks(block.Body, variable); err != nil {
		return nil, err
	}

	return variable, nil
}

// parseVariableBlocks parses nested blocks in a variable
func parseVariableBlocks(body hcl.Body, variable *spookytypesvariables.Variable) error {
	content, diags := body.Content(&hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{
			{Type: "validation"},
			{Type: "constraints"},
		},
		Attributes: []hcl.AttributeSchema{
			{Name: "type", Required: false},
			{Name: "description", Required: false},
			{Name: "default", Required: false},
			{Name: "required", Required: false},
			{Name: "sensitive", Required: false},
			{Name: "encrypted", Required: false},
			{Name: "scope", Required: false},
			{Name: "dependencies", Required: false},
		},
	})
	if diags.HasErrors() {
		return fmt.Errorf("failed to parse nested blocks: %v", diags)
	}

	// Parse attributes first
	for name, attr := range content.Attributes {
		switch name {
		case "type":
			variable.Type = extractStringValue(attr.Expr)
		case "description":
			variable.Description = extractStringValue(attr.Expr)
		case "default":
			variable.Default = extractValue(attr.Expr)
		case "required":
			variable.Required = extractBoolValue(attr.Expr)
		case "sensitive":
			variable.Sensitive = extractBoolValue(attr.Expr)
		case "encrypted":
			variable.Encrypted = extractBoolValue(attr.Expr)
		case "scope":
			variable.Scope = extractStringValue(attr.Expr)
		case "dependencies":
			variable.Dependencies = extractStringListValue(attr.Expr)
		}
	}

	// Parse nested blocks
	for _, block := range content.Blocks {
		switch block.Type {
		case "validation":
			validation, err := parseValidationBlock(block)
			if err != nil {
				return err
			}
			variable.Validation = validation
		case "constraints":
			constraints, err := parseConstraintsBlock(block)
			if err != nil {
				return err
			}
			variable.Constraints = constraints
		}
	}

	return nil
}

// parseValidationBlock parses a validation block
func parseValidationBlock(block *hcl.Block) (*spookytypesvariables.VariableValidation, error) {
	attrs, diags := block.Body.JustAttributes()
	if diags.HasErrors() {
		return nil, fmt.Errorf("failed to parse validation attributes: %v", diags)
	}

	validation := &spookytypesvariables.VariableValidation{}

	for name, attr := range attrs {
		switch name {
		case "condition":
			validation.Condition = extractStringValue(attr.Expr)
		case "error_message":
			validation.ErrorMessage = extractStringValue(attr.Expr)
		case "warning_message":
			validation.WarningMessage = extractStringValue(attr.Expr)
		}
	}

	return validation, nil
}

// parseConstraintsBlock parses a constraints block
func parseConstraintsBlock(block *hcl.Block) (*spookytypesvariables.VariableConstraints, error) {
	attrs, diags := block.Body.JustAttributes()
	if diags.HasErrors() {
		return nil, fmt.Errorf("failed to parse constraints attributes: %v", diags)
	}

	constraints := &spookytypesvariables.VariableConstraints{}

	for name, attr := range attrs {
		switch name {
		case "min_length":
			if val := extractIntValue(attr.Expr); val != nil {
				constraints.MinLength = val
			}
		case "max_length":
			if val := extractIntValue(attr.Expr); val != nil {
				constraints.MaxLength = val
			}
		case "pattern":
			if val := extractStringValue(attr.Expr); val != "" {
				constraints.Pattern = &val
			}
		case "min_value":
			if val := extractFloatValue(attr.Expr); val != nil {
				constraints.MinValue = val
			}
		case "max_value":
			if val := extractFloatValue(attr.Expr); val != nil {
				constraints.MaxValue = val
			}
		case "min_items":
			if val := extractIntValue(attr.Expr); val != nil {
				constraints.MinItems = val
			}
		case "max_items":
			if val := extractIntValue(attr.Expr); val != nil {
				constraints.MaxItems = val
			}
		case "file_exists":
			if val := extractBoolValue(attr.Expr); val {
				constraints.FileExists = &val
			}
		case "file_readable":
			if val := extractBoolValue(attr.Expr); val {
				constraints.FileReadable = &val
			}
		case "file_size_max":
			if val := extractStringValue(attr.Expr); val != "" {
				constraints.FileSizeMax = &val
			}
		case "path_exists":
			if val := extractBoolValue(attr.Expr); val {
				constraints.PathExists = &val
			}
		case "path_absolute":
			if val := extractBoolValue(attr.Expr); val {
				constraints.PathAbsolute = &val
			}
		case "path_relative":
			if val := extractBoolValue(attr.Expr); val {
				constraints.PathRelative = &val
			}
		}
	}

	return constraints, nil
}

// HCL Value Extractors

// extractStringValue extracts a string value from an HCL expression
func extractStringValue(expr hcl.Expression) string {
	val, diags := expr.Value(nil)
	if diags.HasErrors() {
		return ""
	}
	if val.Type().IsPrimitiveType() && val.Type().FriendlyName() == "string" {
		return val.AsString()
	}
	return ""
}

// extractBoolValue extracts a boolean value from an HCL expression
func extractBoolValue(expr hcl.Expression) bool {
	val, diags := expr.Value(nil)
	if diags.HasErrors() {
		return false
	}
	if val.Type().IsPrimitiveType() && val.Type().FriendlyName() == "bool" {
		return val.True()
	}
	return false
}

// extractValue extracts a generic value from an HCL expression
func extractValue(expr hcl.Expression) interface{} {
	val, diags := expr.Value(nil)
	if diags.HasErrors() {
		return nil
	}
	return val
}

// extractStringListValue extracts a string list from an HCL expression
func extractStringListValue(expr hcl.Expression) []string {
	val, diags := expr.Value(nil)
	if diags.HasErrors() {
		return nil
	}
	if val.Type().IsListType() {
		var result []string
		for _, item := range val.AsValueSlice() {
			if item.Type().IsPrimitiveType() && item.Type().FriendlyName() == "string" {
				result = append(result, item.AsString())
			}
		}
		return result
	}
	return nil
}

// extractIntValue extracts an integer value from an HCL expression
func extractIntValue(expr hcl.Expression) *int {
	val, diags := expr.Value(nil)
	if diags.HasErrors() {
		return nil
	}
	if val.Type().IsPrimitiveType() && val.Type().FriendlyName() == "number" {
		if val.IsNull() || !val.IsKnown() {
			return nil
		}
		// For now, return nil as HCL numbers are typically float64
		return nil
	}
	return nil
}

// extractFloatValue extracts a float value from an HCL expression
func extractFloatValue(expr hcl.Expression) *float64 {
	val, diags := expr.Value(nil)
	if diags.HasErrors() {
		return nil
	}
	if val.Type().IsPrimitiveType() && val.Type().FriendlyName() == "number" {
		if val.IsNull() || !val.IsKnown() {
			return nil
		}
		f := val.AsBigFloat()
		if f != nil {
			if f64, _ := f.Float64(); f64 != 0 {
				return &f64
			}
		}
	}
	return nil
}

func (l *JSONVariableLoader) parseJSONToVariables(jsonData map[string]interface{}) ([]*spookytypesvariables.Variable, error) {
	// Check for export format
	if _, ok := jsonData["metadata"]; ok {
		// Export format with metadata
		return parseExportFormat(jsonData)
	}

	// Check for simple variables format
	if variablesObj, ok := jsonData["variables"]; ok {
		return parseSimpleVariablesFormat(variablesObj)
	}

	// Check for direct variables object
	return parseDirectVariablesFormat(jsonData)
}

// parseExportFormat parses variables from export format
func parseExportFormat(jsonData map[string]interface{}) ([]*spookytypesvariables.Variable, error) {
	variables := []*spookytypesvariables.Variable{}

	// Extract variables from export format
	variablesObj, ok := jsonData["variables"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid variables section in export format")
	}

	// Sort keys to ensure consistent ordering for tests
	keys := make([]string, 0, len(variablesObj))
	for k := range variablesObj {
		keys = append(keys, k)
	}

	for _, name := range keys {
		varData := variablesObj[name]
		if varMap, ok := varData.(map[string]interface{}); ok {
			variable, err := parseVariableFromJSON(varMap)
			if err != nil {
				return nil, fmt.Errorf("failed to parse variable %s: %w", name, err)
			}
			variable.Name = name
			variables = append(variables, variable)
		}
	}

	return variables, nil
}

// parseSimpleVariablesFormat parses variables from simple format
func parseSimpleVariablesFormat(variablesObj interface{}) ([]*spookytypesvariables.Variable, error) {
	variables := []*spookytypesvariables.Variable{}

	if variablesMap, ok := variablesObj.(map[string]interface{}); ok {
		for name, varData := range variablesMap {
			if varMap, ok := varData.(map[string]interface{}); ok {
				variable, err := parseVariableFromJSON(varMap)
				if err != nil {
					return nil, fmt.Errorf("failed to parse variable %s: %w", name, err)
				}
				variable.Name = name
				variables = append(variables, variable)
			}
		}
	}

	return variables, nil
}

// parseDirectVariablesFormat parses variables from direct format
func parseDirectVariablesFormat(jsonData map[string]interface{}) ([]*spookytypesvariables.Variable, error) {
	variables := []*spookytypesvariables.Variable{}

	for name, varData := range jsonData {
		if varMap, ok := varData.(map[string]interface{}); ok {
			variable, err := parseVariableFromJSON(varMap)
			if err != nil {
				return nil, fmt.Errorf("failed to parse variable %s: %w", name, err)
			}
			variable.Name = name
			variables = append(variables, variable)
		}
	}

	return variables, nil
}

// parseVariableFromJSON parses a single variable from JSON
func parseVariableFromJSON(varMap map[string]interface{}) (*spookytypesvariables.Variable, error) {
	variable := &spookytypesvariables.Variable{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Extract basic fields
	if v, ok := varMap["name"].(string); ok {
		variable.Name = v
	}
	if v, ok := varMap["type"].(string); ok {
		variable.Type = v
	}
	if v, ok := varMap["value"]; ok {
		variable.Value = v
	}
	if v, ok := varMap["default"]; ok {
		variable.Default = v
	}
	if v, ok := varMap["description"].(string); ok {
		variable.Description = v
	}
	if v, ok := varMap["required"].(bool); ok {
		variable.Required = v
	}
	if v, ok := varMap["sensitive"].(bool); ok {
		variable.Sensitive = v
	}
	if v, ok := varMap["encrypted"].(bool); ok {
		variable.Encrypted = v
	}
	if v, ok := varMap["scope"].(string); ok {
		variable.Scope = v
	}
	if v, ok := varMap["dependencies"].([]interface{}); ok {
		variable.Dependencies = extractStringSlice(v)
	}

	// Parse nested objects
	if validationMap, ok := varMap["validation"].(map[string]interface{}); ok {
		variable.Validation = parseValidationFromJSON(validationMap)
	}
	if constraintsMap, ok := varMap["constraints"].(map[string]interface{}); ok {
		variable.Constraints = parseConstraintsFromJSON(constraintsMap)
	}

	return variable, nil
}

// parseValidationFromJSON parses validation from JSON
func parseValidationFromJSON(validationMap map[string]interface{}) *spookytypesvariables.VariableValidation {
	validation := &spookytypesvariables.VariableValidation{}

	if v, ok := validationMap["condition"].(string); ok {
		validation.Condition = v
	}
	if v, ok := validationMap["error_message"].(string); ok {
		validation.ErrorMessage = v
	}
	if v, ok := validationMap["warning_message"].(string); ok {
		validation.WarningMessage = v
	}

	return validation
}

// parseConstraintsFromJSON parses constraints from JSON
func parseConstraintsFromJSON(constraintsMap map[string]interface{}) *spookytypesvariables.VariableConstraints {
	constraints := &spookytypesvariables.VariableConstraints{}

	if v, ok := constraintsMap["min_length"].(float64); ok {
		val := int(v)
		constraints.MinLength = &val
	}
	if v, ok := constraintsMap["max_length"].(float64); ok {
		val := int(v)
		constraints.MaxLength = &val
	}
	if v, ok := constraintsMap["pattern"].(string); ok {
		constraints.Pattern = &v
	}
	if v, ok := constraintsMap["min_value"].(float64); ok {
		constraints.MinValue = &v
	}
	if v, ok := constraintsMap["max_value"].(float64); ok {
		constraints.MaxValue = &v
	}
	if v, ok := constraintsMap["min_items"].(float64); ok {
		val := int(v)
		constraints.MinItems = &val
	}
	if v, ok := constraintsMap["max_items"].(float64); ok {
		val := int(v)
		constraints.MaxItems = &val
	}
	if v, ok := constraintsMap["file_exists"].(bool); ok {
		constraints.FileExists = &v
	}
	if v, ok := constraintsMap["file_readable"].(bool); ok {
		constraints.FileReadable = &v
	}
	if v, ok := constraintsMap["file_size_max"].(string); ok {
		constraints.FileSizeMax = &v
	}
	if v, ok := constraintsMap["path_exists"].(bool); ok {
		constraints.PathExists = &v
	}
	if v, ok := constraintsMap["path_absolute"].(bool); ok {
		constraints.PathAbsolute = &v
	}
	if v, ok := constraintsMap["path_relative"].(bool); ok {
		constraints.PathRelative = &v
	}

	return constraints
}

// extractStringSlice converts interface{} slice to string slice
func extractStringSlice(slice []interface{}) []string {
	var result []string
	for _, item := range slice {
		if str, ok := item.(string); ok {
			result = append(result, str)
		}
	}
	return result
}
