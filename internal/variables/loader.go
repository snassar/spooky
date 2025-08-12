// Package variables provides variable loading functionality for the spooky codebase.
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

// Loader provides variable loading functionality
type Loader struct {
	logger spookytypeslogging.Logger
}

// NewLoader creates a new variable loader
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

	// Read the file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	// Parse HCL content
	parser := hclparse.NewParser()
	file, diags := parser.ParseHCL(data, filePath)
	if diags.HasErrors() {
		return nil, fmt.Errorf("failed to parse HCL file %s: %w", filePath, diags.Error())
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
		return nil, fmt.Errorf("failed to parse variables block in %s: %w", filePath, diags.Error())
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
		return nil, fmt.Errorf("failed to parse variable blocks: %w", diags.Error())
	}

	for _, variableBlock := range content.Blocks {
		variable, err := l.parseVariableBlock(variableBlock, sourceFile)
		if err != nil {
			return nil, fmt.Errorf("failed to parse variable block %s: %w", variableBlock.Labels[0], err)
		}

		variables[variable.Name] = variable
	}

	return variables, nil
}

// parseVariableBlock parses a single variable block
func (l *Loader) parseVariableBlock(block *hcl.Block, sourceFile string) (*spookytypesvariables.Variable, error) {
	if len(block.Labels) == 0 {
		return nil, fmt.Errorf("variable block missing name label")
	}

	variableName := block.Labels[0]

	// Parse variable attributes
	content, diags := block.Body.Content(&hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{
			{Name: "type", Required: true},
			{Name: "description", Required: false},
			{Name: "default", Required: false},
			{Name: "required", Required: false},
			{Name: "sensitive", Required: false},
			{Name: "encrypted", Required: false},
			{Name: "scope", Required: false},
			{Name: "dependencies", Required: false},
			{Name: "metadata", Required: false},
		},
		Blocks: []hcl.BlockHeaderSchema{
			{Type: "validation"},
			{Type: "constraints"},
		},
	})
	if diags.HasErrors() {
		return nil, fmt.Errorf("failed to parse variable attributes: %w", diags.Error())
	}

	variable := &spookytypesvariables.Variable{
		Name:       variableName,
		SourceFile: sourceFile,
		SourceLine: block.DefRange.Start.Line,
		Metadata:   make(map[string]interface{}),
	}

	// Parse type attribute
	if attr, exists := content.Attributes["type"]; exists {
		typeVal, diags := attr.Expr.Value(nil)
		if diags.HasErrors() {
			return nil, fmt.Errorf("failed to parse type attribute: %w", diags.Error())
		}
		if typeVal.Type() != cty.String {
			return nil, fmt.Errorf("type attribute must be a string")
		}
		variable.Type = spookytypesvariables.VariableType(typeVal.AsString())
	}

	// Parse description attribute
	if attr, exists := content.Attributes["description"]; exists {
		descVal, diags := attr.Expr.Value(nil)
		if diags.HasErrors() {
			return nil, fmt.Errorf("failed to parse description attribute: %w", diags.Error())
		}
		if descVal.Type() != cty.String {
			return nil, fmt.Errorf("description attribute must be a string")
		}
		variable.Description = descVal.AsString()
	}

	// Parse default attribute
	if attr, exists := content.Attributes["default"]; exists {
		defaultVal, diags := attr.Expr.Value(nil)
		if diags.HasErrors() {
			return nil, fmt.Errorf("failed to parse default attribute: %w", diags.Error())
		}
		variable.Default = ctyToInterface(defaultVal)
	}

	// Parse required attribute
	if attr, exists := content.Attributes["required"]; exists {
		requiredVal, diags := attr.Expr.Value(nil)
		if diags.HasErrors() {
			return nil, fmt.Errorf("failed to parse required attribute: %w", diags.Error())
		}
		if requiredVal.Type() != cty.Bool {
			return nil, fmt.Errorf("required attribute must be a boolean")
		}
		variable.Required = requiredVal.True()
	}

	// Parse sensitive attribute
	if attr, exists := content.Attributes["sensitive"]; exists {
		sensitiveVal, diags := attr.Expr.Value(nil)
		if diags.HasErrors() {
			return nil, fmt.Errorf("failed to parse sensitive attribute: %w", diags.Error())
		}
		if sensitiveVal.Type() != cty.Bool {
			return nil, fmt.Errorf("sensitive attribute must be a boolean")
		}
		variable.Sensitive = sensitiveVal.True()
	}

	// Parse encrypted attribute
	if attr, exists := content.Attributes["encrypted"]; exists {
		encryptedVal, diags := attr.Expr.Value(nil)
		if diags.HasErrors() {
			return nil, fmt.Errorf("failed to parse encrypted attribute: %w", diags.Error())
		}
		if encryptedVal.Type() != cty.Bool {
			return nil, fmt.Errorf("encrypted attribute must be a boolean")
		}
		variable.Encrypted = encryptedVal.True()
	}

	// Parse scope attribute
	if attr, exists := content.Attributes["scope"]; exists {
		scopeVal, diags := attr.Expr.Value(nil)
		if diags.HasErrors() {
			return nil, fmt.Errorf("failed to parse scope attribute: %w", diags.Error())
		}
		if scopeVal.Type() != cty.String {
			return nil, fmt.Errorf("scope attribute must be a string")
		}
		variable.Scope = spookytypesvariables.VariableScope(scopeVal.AsString())
	}

	// Parse dependencies attribute
	if attr, exists := content.Attributes["dependencies"]; exists {
		depsVal, diags := attr.Expr.Value(nil)
		if diags.HasErrors() {
			return nil, fmt.Errorf("failed to parse dependencies attribute: %w", diags.Error())
		}
		if depsVal.Type() != cty.List(cty.String) {
			return nil, fmt.Errorf("dependencies attribute must be a list of strings")
		}
		var deps []string
		for _, dep := range depsVal.AsValueSlice() {
			deps = append(deps, dep.AsString())
		}
		variable.Dependencies = deps
	}

	// Parse metadata attribute
	if attr, exists := content.Attributes["metadata"]; exists {
		metaVal, diags := attr.Expr.Value(nil)
		if diags.HasErrors() {
			return nil, fmt.Errorf("failed to parse metadata attribute: %w", diags.Error())
		}
		if metaVal.Type() != cty.Map(cty.String) {
			return nil, fmt.Errorf("metadata attribute must be a map of strings")
		}
		for key, value := range metaVal.AsValueMap() {
			variable.Metadata[key] = value.AsString()
		}
	}

	// Parse validation block
	for _, validationBlock := range content.Blocks {
		if validationBlock.Type == "validation" {
			validation, err := l.parseValidationBlock(validationBlock)
			if err != nil {
				return nil, fmt.Errorf("failed to parse validation block: %w", err)
			}
			variable.Validation = validation
		}
	}

	// Parse constraints block
	for _, constraintsBlock := range content.Blocks {
		if constraintsBlock.Type == "constraints" {
			constraints, err := l.parseConstraintsBlock(constraintsBlock)
			if err != nil {
				return nil, fmt.Errorf("failed to parse constraints block: %w", err)
			}
			variable.Constraints = constraints
		}
	}

	return variable, nil
}

// parseValidationBlock parses a validation block
func (l *Loader) parseValidationBlock(block *hcl.Block) (*spookytypesvariables.VariableValidation, error) {
	content, diags := block.Body.Content(&hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{
			{Name: "condition", Required: true},
			{Name: "error_message", Required: true},
			{Name: "warning_message", Required: false},
		},
	})
	if diags.HasErrors() {
		return nil, fmt.Errorf("failed to parse validation block attributes: %w", diags.Error())
	}

	validation := &spookytypesvariables.VariableValidation{}

	// Parse condition attribute
	if attr, exists := content.Attributes["condition"]; exists {
		conditionVal, diags := attr.Expr.Value(nil)
		if diags.HasErrors() {
			return nil, fmt.Errorf("failed to parse condition attribute: %w", diags.Error())
		}
		if conditionVal.Type() != cty.String {
			return nil, fmt.Errorf("condition attribute must be a string")
		}
		validation.Condition = conditionVal.AsString()
	}

	// Parse error_message attribute
	if attr, exists := content.Attributes["error_message"]; exists {
		errorMsgVal, diags := attr.Expr.Value(nil)
		if diags.HasErrors() {
			return nil, fmt.Errorf("failed to parse error_message attribute: %w", diags.Error())
		}
		if errorMsgVal.Type() != cty.String {
			return nil, fmt.Errorf("error_message attribute must be a string")
		}
		validation.ErrorMessage = errorMsgVal.AsString()
	}

	// Parse warning_message attribute
	if attr, exists := content.Attributes["warning_message"]; exists {
		warningMsgVal, diags := attr.Expr.Value(nil)
		if diags.HasErrors() {
			return nil, fmt.Errorf("failed to parse warning_message attribute: %w", diags.Error())
		}
		if warningMsgVal.Type() != cty.String {
			return nil, fmt.Errorf("warning_message attribute must be a string")
		}
		validation.WarningMessage = warningMsgVal.AsString()
	}

	return validation, nil
}

// parseConstraintsBlock parses a constraints block
func (l *Loader) parseConstraintsBlock(block *hcl.Block) (*spookytypesvariables.VariableConstraints, error) {
	content, diags := block.Body.Content(&hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{
			{Name: "min_length", Required: false},
			{Name: "max_length", Required: false},
			{Name: "pattern", Required: false},
			{Name: "min_value", Required: false},
			{Name: "max_value", Required: false},
			{Name: "min_items", Required: false},
			{Name: "max_items", Required: false},
			{Name: "file_exists", Required: false},
			{Name: "file_readable", Required: false},
			{Name: "file_size_max", Required: false},
			{Name: "path_exists", Required: false},
			{Name: "path_absolute", Required: false},
			{Name: "path_relative", Required: false},
		},
	})
	if diags.HasErrors() {
		return nil, fmt.Errorf("failed to parse constraints block attributes: %w", diags.Error())
	}

	constraints := &spookytypesvariables.VariableConstraints{}

	// Parse string constraints
	if attr, exists := content.Attributes["min_length"]; exists {
		minLengthVal, diags := attr.Expr.Value(nil)
		if diags.HasErrors() {
			return nil, fmt.Errorf("failed to parse min_length attribute: %w", diags.Error())
		}
		if minLengthVal.Type() != cty.Number {
			return nil, fmt.Errorf("min_length attribute must be a number")
		}
		minLengthInt, _ := minLengthVal.AsBigFloat().Int64()
		minLength := int(minLengthInt)
		constraints.MinLength = &minLength
	}

	if attr, exists := content.Attributes["max_length"]; exists {
		maxLengthVal, diags := attr.Expr.Value(nil)
		if diags.HasErrors() {
			return nil, fmt.Errorf("failed to parse max_length attribute: %w", diags.Error())
		}
		if maxLengthVal.Type() != cty.Number {
			return nil, fmt.Errorf("max_length attribute must be a number")
		}
		maxLengthInt, _ := maxLengthVal.AsBigFloat().Int64()
		maxLength := int(maxLengthInt)
		constraints.MaxLength = &maxLength
	}

	if attr, exists := content.Attributes["pattern"]; exists {
		patternVal, diags := attr.Expr.Value(nil)
		if diags.HasErrors() {
			return nil, fmt.Errorf("failed to parse pattern attribute: %w", diags.Error())
		}
		if patternVal.Type() != cty.String {
			return nil, fmt.Errorf("pattern attribute must be a string")
		}
		pattern := patternVal.AsString()
		constraints.Pattern = &pattern
	}

	// Parse numeric constraints
	if attr, exists := content.Attributes["min_value"]; exists {
		minValueVal, diags := attr.Expr.Value(nil)
		if diags.HasErrors() {
			return nil, fmt.Errorf("failed to parse min_value attribute: %w", diags.Error())
		}
		if minValueVal.Type() != cty.Number {
			return nil, fmt.Errorf("min_value attribute must be a number")
		}
		minValue, _ := minValueVal.AsBigFloat().Float64()
		constraints.MinValue = &minValue
	}

	if attr, exists := content.Attributes["max_value"]; exists {
		maxValueVal, diags := attr.Expr.Value(nil)
		if diags.HasErrors() {
			return nil, fmt.Errorf("failed to parse max_value attribute: %w", diags.Error())
		}
		if maxValueVal.Type() != cty.Number {
			return nil, fmt.Errorf("max_value attribute must be a number")
		}
		maxValue, _ := maxValueVal.AsBigFloat().Float64()
		constraints.MaxValue = &maxValue
	}

	// Parse list constraints
	if attr, exists := content.Attributes["min_items"]; exists {
		minItemsVal, diags := attr.Expr.Value(nil)
		if diags.HasErrors() {
			return nil, fmt.Errorf("failed to parse min_items attribute: %w", diags.Error())
		}
		if minItemsVal.Type() != cty.Number {
			return nil, fmt.Errorf("min_items attribute must be a number")
		}
		minItemsInt, _ := minItemsVal.AsBigFloat().Int64()
		minItems := int(minItemsInt)
		constraints.MinItems = &minItems
	}

	if attr, exists := content.Attributes["max_items"]; exists {
		maxItemsVal, diags := attr.Expr.Value(nil)
		if diags.HasErrors() {
			return nil, fmt.Errorf("failed to parse max_items attribute: %w", diags.Error())
		}
		if maxItemsVal.Type() != cty.Number {
			return nil, fmt.Errorf("max_items attribute must be a number")
		}
		maxItemsInt, _ := maxItemsVal.AsBigFloat().Int64()
		maxItems := int(maxItemsInt)
		constraints.MaxItems = &maxItems
	}

	// Parse file constraints
	if attr, exists := content.Attributes["file_exists"]; exists {
		fileExistsVal, diags := attr.Expr.Value(nil)
		if diags.HasErrors() {
			return nil, fmt.Errorf("failed to parse file_exists attribute: %w", diags.Error())
		}
		if fileExistsVal.Type() != cty.Bool {
			return nil, fmt.Errorf("file_exists attribute must be a boolean")
		}
		fileExists := fileExistsVal.True()
		constraints.FileExists = &fileExists
	}

	if attr, exists := content.Attributes["file_readable"]; exists {
		fileReadableVal, diags := attr.Expr.Value(nil)
		if diags.HasErrors() {
			return nil, fmt.Errorf("failed to parse file_readable attribute: %w", diags.Error())
		}
		if fileReadableVal.Type() != cty.Bool {
			return nil, fmt.Errorf("file_readable attribute must be a boolean")
		}
		fileReadable := fileReadableVal.True()
		constraints.FileReadable = &fileReadable
	}

	if attr, exists := content.Attributes["file_size_max"]; exists {
		fileSizeMaxVal, diags := attr.Expr.Value(nil)
		if diags.HasErrors() {
			return nil, fmt.Errorf("failed to parse file_size_max attribute: %w", diags.Error())
		}
		if fileSizeMaxVal.Type() != cty.String {
			return nil, fmt.Errorf("file_size_max attribute must be a string")
		}
		fileSizeMax := fileSizeMaxVal.AsString()
		constraints.FileSizeMax = &fileSizeMax
	}

	// Parse path constraints
	if attr, exists := content.Attributes["path_exists"]; exists {
		pathExistsVal, diags := attr.Expr.Value(nil)
		if diags.HasErrors() {
			return nil, fmt.Errorf("failed to parse path_exists attribute: %w", diags.Error())
		}
		if pathExistsVal.Type() != cty.Bool {
			return nil, fmt.Errorf("path_exists attribute must be a boolean")
		}
		pathExists := pathExistsVal.True()
		constraints.PathExists = &pathExists
	}

	if attr, exists := content.Attributes["path_absolute"]; exists {
		pathAbsoluteVal, diags := attr.Expr.Value(nil)
		if diags.HasErrors() {
			return nil, fmt.Errorf("failed to parse path_absolute attribute: %w", diags.Error())
		}
		if pathAbsoluteVal.Type() != cty.Bool {
			return nil, fmt.Errorf("path_absolute attribute must be a boolean")
		}
		pathAbsolute := pathAbsoluteVal.True()
		constraints.PathAbsolute = &pathAbsolute
	}

	if attr, exists := content.Attributes["path_relative"]; exists {
		pathRelativeVal, diags := attr.Expr.Value(nil)
		if diags.HasErrors() {
			return nil, fmt.Errorf("failed to parse path_relative attribute: %w", diags.Error())
		}
		if pathRelativeVal.Type() != cty.Bool {
			return nil, fmt.Errorf("path_relative attribute must be a boolean")
		}
		pathRelative := pathRelativeVal.True()
		constraints.PathRelative = &pathRelative
	}

	return constraints, nil
}

// ctyToInterface converts a cty.Value to interface{}
func ctyToInterface(val cty.Value) interface{} {
	if val.IsNull() {
		return nil
	}

	switch val.Type() {
	case cty.String:
		return val.AsString()
	case cty.Number:
		// Convert to float64 for numbers
		if val.IsKnown() && !val.IsNull() {
			f, _ := val.AsBigFloat().Float64()
			return f
		}
		return 0.0
	case cty.Bool:
		return val.True()
	case cty.List(cty.String):
		var result []string
		for _, item := range val.AsValueSlice() {
			result = append(result, item.AsString())
		}
		return result
	case cty.Map(cty.String):
		result := make(map[string]string)
		for key, value := range val.AsValueMap() {
			result[key] = value.AsString()
		}
		return result
	default:
		// For complex types, return as string representation
		return val.GoString()
	}
}

// getVariableNames returns a list of variable names
func getVariableNames(variables map[string]*spookytypesvariables.Variable) []string {
	var names []string
	for name := range variables {
		names = append(names, name)
	}
	return names
}
