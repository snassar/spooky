// Package variables provides variable loading and management functionality for the spooky codebase.
package variables

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/zclconf/go-cty/cty"

	spookytypeslogging "spooky/internal/types/logging"
	spookytypesvariables "spooky/internal/types/variables"
)

// Loader provides functionality to load variables from HCL files
type Loader struct {
	logger spookytypeslogging.Logger
}

// NewLoader creates a new variable loader instance
func NewLoader(logger spookytypeslogging.Logger) *Loader {
	return &Loader{
		logger: logger,
	}
}

// LoadVariablesFromFile loads variables from a single HCL file
func (l *Loader) LoadVariablesFromFile(ctx context.Context, filePath string) (map[string]*spookytypesvariables.Variable, error) {
	l.logger.Debug("Loading variables from file", map[string]interface{}{
		"file_path": filePath,
	})

	// Read file content
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	// Parse HCL content
	parser := hclparse.NewParser()
	file, diags := parser.ParseHCL(data, filePath)
	if diags.HasErrors() {
		return nil, fmt.Errorf("failed to parse HCL file %s: %s", filePath, diags.Error())
	}

	// Extract variables block
	content, diags := file.Body.Content(&hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{
			{
				Type: "variables",
			},
		},
	})
	if diags.HasErrors() {
		return nil, fmt.Errorf("failed to parse variables block in %s: %s", filePath, diags.Error())
	}

	if len(content.Blocks) == 0 {
		return nil, fmt.Errorf("no variables block found in %s", filePath)
	}

	variablesBlock := content.Blocks[0]
	variables, err := l.parseVariablesBlock(variablesBlock, filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse variables block: %w", err)
	}

	l.logger.Info("Variables loaded successfully", map[string]interface{}{
		"source":    filePath,
		"count":     len(variables),
		"variables": getVariableNames(variables),
	})

	return variables, nil
}

// LoadVariablesFromDirectory loads variables from all HCL files in a directory
func (l *Loader) LoadVariablesFromDirectory(ctx context.Context, dirPath string) (map[string]*spookytypesvariables.Variable, error) {
	l.logger.Debug("Loading variables from directory", map[string]interface{}{
		"dir_path": dirPath,
	})

	// Read directory entries
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory %s: %w", dirPath, err)
	}

	var allVariables = make(map[string]*spookytypesvariables.Variable)
	var loadErrors []string

	// Process each HCL file in the directory
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		if !strings.HasSuffix(entry.Name(), ".hcl") {
			continue
		}

		filePath := filepath.Join(dirPath, entry.Name())
		variables, err := l.LoadVariablesFromFile(ctx, filePath)
		if err != nil {
			loadErrors = append(loadErrors, fmt.Sprintf("%s: %v", entry.Name(), err))
			continue
		}

		// Merge variables from this file
		for name, variable := range variables {
			allVariables[name] = variable
		}
	}

	if len(loadErrors) > 0 {
		l.logger.Warn("Some variable files failed to load", map[string]interface{}{
			"dir_path": dirPath,
			"errors":   loadErrors,
		})
	}

	l.logger.Info("Variables loaded from directory", map[string]interface{}{
		"dir_path": dirPath,
		"count":    len(allVariables),
		"errors":   len(loadErrors),
	})

	return allVariables, nil
}

// parseVariablesBlock parses the variables block and extracts variable definitions
func (l *Loader) parseVariablesBlock(block *hcl.Block, sourceFile string) (map[string]*spookytypesvariables.Variable, error) {
	variables := make(map[string]*spookytypesvariables.Variable)

	// Parse variable blocks within the variables block
	content, diags := block.Body.Content(&hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{
			{
				Type:       "variable",
				LabelNames: []string{"name"},
			},
		},
	})
	if diags.HasErrors() {
		return nil, fmt.Errorf("failed to parse variable blocks: %s", diags.Error())
	}

	for _, variableBlock := range content.Blocks {
		variable, err := l.parseVariableBlock(variableBlock, sourceFile)
		if err != nil {
			l.logger.Warn("Failed to parse variable block", map[string]interface{}{
				"variable_name": variableBlock.Labels[0],
				"source_file":   sourceFile,
				"error":         err.Error(),
			})
			continue
		}
		variables[variable.Name] = variable
	}

	return variables, nil
}

// parseVariableBlock parses a single variable block
func (l *Loader) parseVariableBlock(block *hcl.Block, sourceFile string) (*spookytypesvariables.Variable, error) {
	variableName := block.Labels[0]

	// Parse both attributes and blocks
	content, diags := block.Body.Content(&hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{
			{Name: "type", Required: false},
			{Name: "default", Required: false},
			{Name: "description", Required: false},
			{Name: "sensitive", Required: false},
			{Name: "validation", Required: false},
			{Name: "constraints", Required: false},
			{Name: "scope", Required: false},
			{Name: "dependencies", Required: false},
		},
		Blocks: []hcl.BlockHeaderSchema{
			{Type: "validation"},
			{Type: "constraints"},
		},
	})
	if diags.HasErrors() {
		return nil, fmt.Errorf("failed to parse variable block: %s", diags.Error())
	}

	variable := &spookytypesvariables.Variable{
		Name: variableName,
	}

	// Add source file information
	variable.SourceFile = sourceFile
	variable.SourceLine = block.DefRange.Start.Line
	variable.Metadata = make(map[string]interface{})
	variable.Metadata["source_file"] = sourceFile

	// Parse attributes
	attrs := content.Attributes

	// Parse type
	if attr, exists := attrs["type"]; exists {
		val, diags := attr.Expr.Value(nil)
		if diags.HasErrors() {
			return nil, fmt.Errorf("invalid type: %s", diags.Error())
		}
		if val.Type() == cty.String {
			variable.Type = spookytypesvariables.VariableType(val.AsString())
		}
	}

	// Parse default value
	if attr, exists := attrs["default"]; exists {
		val, diags := attr.Expr.Value(nil)
		if diags.HasErrors() {
			return nil, fmt.Errorf("invalid default: %s", diags.Error())
		}
		variable.Default = val
	}

	// Parse description
	if attr, exists := attrs["description"]; exists {
		val, diags := attr.Expr.Value(nil)
		if diags.HasErrors() {
			return nil, fmt.Errorf("invalid description: %s", diags.Error())
		}
		if val.Type() == cty.String {
			variable.Description = val.AsString()
		}
	}

	// Parse sensitive flag
	if attr, exists := attrs["sensitive"]; exists {
		val, diags := attr.Expr.Value(nil)
		if diags.HasErrors() {
			return nil, fmt.Errorf("invalid sensitive: %s", diags.Error())
		}
		if val.Type() == cty.Bool {
			variable.Sensitive = val.True()
		}
	}

	// Parse scope
	if attr, exists := attrs["scope"]; exists {
		val, diags := attr.Expr.Value(nil)
		if diags.HasErrors() {
			return nil, fmt.Errorf("invalid scope: %s", diags.Error())
		}
		if val.Type() == cty.String {
			variable.Scope = spookytypesvariables.VariableScope(val.AsString())
		}
	}

	// Parse dependencies
	if attr, exists := attrs["dependencies"]; exists {
		dependencies, err := l.parseArrayAttribute(attr)
		if err != nil {
			return nil, fmt.Errorf("invalid dependencies: %w", err)
		}
		variable.Dependencies = dependencies
	}

	// Parse blocks
	for _, block := range content.Blocks {
		switch block.Type {
		case "validation":
			validation, err := l.parseValidationBlock(block)
			if err != nil {
				return nil, fmt.Errorf("failed to parse validation block: %w", err)
			}
			variable.Validation = validation
		case "constraints":
			constraints, err := l.parseConstraintsBlock(block)
			if err != nil {
				return nil, fmt.Errorf("failed to parse constraints block: %w", err)
			}
			variable.Constraints = constraints
		}
	}

	return variable, nil
}

// parseArrayAttribute parses an array attribute into a []string
func (l *Loader) parseArrayAttribute(attr *hcl.Attribute) ([]string, error) {
	val, diags := attr.Expr.Value(nil)
	if diags.HasErrors() {
		return nil, fmt.Errorf("failed to parse array attribute: %s", diags.Error())
	}

	if val.Type() != cty.List(cty.String) {
		return nil, fmt.Errorf("expected list of strings, got %s", val.Type().FriendlyName())
	}

	var result []string
	for _, item := range val.AsValueSlice() {
		result = append(result, item.AsString())
	}

	return result, nil
}

// parseValidationBlock parses a validation block
func (l *Loader) parseValidationBlock(block *hcl.Block) (*spookytypesvariables.VariableValidation, error) {
	content, diags := block.Body.Content(&hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{
			{Name: "condition", Required: false},
			{Name: "error_message", Required: false},
			{Name: "warning_message", Required: false},
		},
	})
	if diags.HasErrors() {
		return nil, fmt.Errorf("failed to parse validation block: %s", diags.Error())
	}

	validation := &spookytypesvariables.VariableValidation{}

	attrs := content.Attributes

	if attr, exists := attrs["condition"]; exists {
		val, diags := attr.Expr.Value(nil)
		if diags.HasErrors() {
			return nil, fmt.Errorf("invalid condition: %s", diags.Error())
		}
		if val.Type() == cty.String {
			validation.Condition = val.AsString()
		}
	}

	if attr, exists := attrs["error_message"]; exists {
		val, diags := attr.Expr.Value(nil)
		if diags.HasErrors() {
			return nil, fmt.Errorf("invalid error_message: %s", diags.Error())
		}
		if val.Type() == cty.String {
			validation.ErrorMessage = val.AsString()
		}
	}

	if attr, exists := attrs["warning_message"]; exists {
		val, diags := attr.Expr.Value(nil)
		if diags.HasErrors() {
			return nil, fmt.Errorf("invalid warning_message: %s", diags.Error())
		}
		if val.Type() == cty.String {
			validation.WarningMessage = val.AsString()
		}
	}

	return validation, nil
}

// parseConstraintsBlock parses a constraints block
func (l *Loader) parseConstraintsBlock(block *hcl.Block) (*spookytypesvariables.VariableConstraints, error) {
	content, diags := block.Body.Content(&hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{
			{Name: "min_length", Required: false},
			{Name: "max_length", Required: false},
			{Name: "min_value", Required: false},
			{Name: "max_value", Required: false},
			{Name: "pattern", Required: false},
			{Name: "allowed_values", Required: false},
		},
	})
	if diags.HasErrors() {
		return nil, fmt.Errorf("failed to parse constraints block: %s", diags.Error())
	}

	constraints := &spookytypesvariables.VariableConstraints{}

	attrs := content.Attributes

	if attr, exists := attrs["min_length"]; exists {
		val, diags := attr.Expr.Value(nil)
		if diags.HasErrors() {
			return nil, fmt.Errorf("invalid min_length: %s", diags.Error())
		}
		if val.Type() == cty.Number {
			minInt, _ := val.AsBigFloat().Int64()
			minLength := int(minInt)
			constraints.MinLength = &minLength
		}
	}

	if attr, exists := attrs["max_length"]; exists {
		val, diags := attr.Expr.Value(nil)
		if diags.HasErrors() {
			return nil, fmt.Errorf("invalid max_length: %s", diags.Error())
		}
		if val.Type() == cty.Number {
			maxInt, _ := val.AsBigFloat().Int64()
			maxLength := int(maxInt)
			constraints.MaxLength = &maxLength
		}
	}

	if attr, exists := attrs["min_value"]; exists {
		val, diags := attr.Expr.Value(nil)
		if diags.HasErrors() {
			return nil, fmt.Errorf("invalid min_value: %s", diags.Error())
		}
		if val.Type() == cty.Number {
			minFloat, _ := val.AsBigFloat().Float64()
			constraints.MinValue = &minFloat
		}
	}

	if attr, exists := attrs["max_value"]; exists {
		val, diags := attr.Expr.Value(nil)
		if diags.HasErrors() {
			return nil, fmt.Errorf("invalid max_value: %s", diags.Error())
		}
		if val.Type() == cty.Number {
			maxFloat, _ := val.AsBigFloat().Float64()
			constraints.MaxValue = &maxFloat
		}
	}

	if attr, exists := attrs["pattern"]; exists {
		val, diags := attr.Expr.Value(nil)
		if diags.HasErrors() {
			return nil, fmt.Errorf("invalid pattern: %s", diags.Error())
		}
		if val.Type() == cty.String {
			pattern := val.AsString()
			constraints.Pattern = &pattern
		}
	}

	return constraints, nil
}

// getVariableNames extracts variable names from a map of variables
func getVariableNames(variables map[string]*spookytypesvariables.Variable) []string {
	var names []string
	for name := range variables {
		names = append(names, name)
	}
	return names
}
