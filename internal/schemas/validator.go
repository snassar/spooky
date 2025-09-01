package schemas

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/zclconf/go-cty/cty"
)

// SimpleValidator provides focused, essential validation without over-engineering
type SimpleValidator struct {
	logger *slog.Logger
}

// NewSimpleValidator creates a new simplified validator
func NewSimpleValidator() *SimpleValidator {
	return &SimpleValidator{
		logger: slog.Default(),
	}
}

// ValidateHCLContent validates HCL content with essential validation only
func (v *SimpleValidator) ValidateHCLContent(schemaName, content string) (*ValidationResult, error) {
	// Parse HCL into structured data
	parsedData, err := v.ParseHCLContent(content)
	if err != nil {
		return &ValidationResult{
			IsValid:    false,
			Errors:     []ValidationError{{Message: fmt.Sprintf("HCL parsing failed: %v", err), Severity: "error"}},
			Warnings:   []ValidationWarning{},
			SchemaName: schemaName,
		}, nil
	}

	// Validate with essential rules only
	result := &ValidationResult{
		IsValid:    true,
		Errors:     []ValidationError{},
		Warnings:   []ValidationWarning{},
		SchemaName: schemaName,
	}
	v.validateEssential(schemaName, parsedData, result)
	return result, nil
}

// ValidateData validates parsed data with essential validation only
func (v *SimpleValidator) ValidateData(schemaName string, data map[string]interface{}) *ValidationResult {
	result := &ValidationResult{
		IsValid:    true,
		Errors:     []ValidationError{},
		Warnings:   []ValidationWarning{},
		SchemaName: schemaName,
	}

	// Essential validation only
	v.validateEssential(schemaName, data, result)

	// Update validity based on errors
	if len(result.Errors) > 0 {
		result.IsValid = false
	}

	return result
}

// validateEssential performs essential validation without over-engineering
func (v *SimpleValidator) validateEssential(schemaName string, data map[string]interface{}, result *ValidationResult) {
	switch schemaName {
	case "project":
		v.validateProjectEssential(data, result)
	case "machines":
		v.validateMachinesEssential(data, result)
	case "actions":
		v.validateActionsEssential(data, result)
	case "variables":
		v.validateVariablesEssential(data, result)
	case "project-directory":
		v.validateProjectDirectoryEssential(data, result)
	}
}

// validateProjectEssential validates project with essential rules only
func (v *SimpleValidator) validateProjectEssential(data map[string]interface{}, result *ValidationResult) {
	// Check for required project block
	projectBlock, exists := data["project"]
	if !exists {
		result.Errors = append(result.Errors, ValidationError{
			Field:    "project",
			Message:  "missing required project block",
			Severity: "error",
		})
		return
	}

	// Validate project block structure
	if projectMap, ok := projectBlock.(map[string]interface{}); ok {
		v.validateProjectBlockEssential(projectMap, result)
	} else {
		result.Errors = append(result.Errors, ValidationError{
			Field:    "project",
			Message:  "project block must be a configuration block",
			Severity: "error",
		})
	}
}

// validateProjectBlockEssential validates project block contents with essential rules
func (v *SimpleValidator) validateProjectBlockEssential(project map[string]interface{}, result *ValidationResult) {
	// Required fields
	if name, exists := project["name"]; !exists || name == "" {
		result.Errors = append(result.Errors, ValidationError{
			Field:    "project.name",
			Message:  "project name is required",
			Severity: "error",
		})
	}

	// Description length validation
	if desc, exists := project["description"]; exists {
		if descStr, ok := desc.(string); ok && len(descStr) > 1024 {
			result.Errors = append(result.Errors, ValidationError{
				Field:    "project.description",
				Value:    descStr,
				Message:  "description must be at most 1024 characters",
				Severity: "error",
			})
		}
	}

	// Range validation for numeric fields
	if maxParallel, exists := project["run_max_parallel"]; exists {
		if val, ok := maxParallel.(float64); ok {
			if val < 1 || val > 100 {
				result.Errors = append(result.Errors, ValidationError{
					Field:    "project.run_max_parallel",
					Value:    val,
					Message:  "run_max_parallel must be between 1 and 100",
					Severity: "error",
				})
			}
		}
	}

	if factsTimeout, exists := project["facts_timeout"]; exists {
		if val, ok := factsTimeout.(float64); ok {
			if val < 1 || val > 3600 {
				result.Errors = append(result.Errors, ValidationError{
					Field:    "project.facts_timeout",
					Value:    val,
					Message:  "facts_timeout must be between 1 and 3600",
					Severity: "error",
				})
			}
		}
	}
}

// validateMachinesEssential validates machines with essential rules only
func (v *SimpleValidator) validateMachinesEssential(data map[string]interface{}, result *ValidationResult) {
	// Check for required machines block
	machinesBlock, exists := data["machines"]
	if !exists {
		result.Errors = append(result.Errors, ValidationError{
			Field:    "machines",
			Message:  "missing required machines block",
			Severity: "error",
		})
		return
	}

	// Validate machines block structure
	if machinesMap, ok := machinesBlock.(map[string]interface{}); ok {
		v.validateMachinesBlockEssential(machinesMap, result)
	} else {
		result.Errors = append(result.Errors, ValidationError{
			Field:    "machines",
			Message:  "machines block must be a configuration block",
			Severity: "error",
		})
	}
}

// validateMachinesBlockEssential validates machines block contents with essential rules
func (v *SimpleValidator) validateMachinesBlockEssential(machines map[string]interface{}, result *ValidationResult) {
	// Check for at least one machine
	hasMachine := false
	for machineName, machineValue := range machines {
		if machineName != "group" {
			if machineMap, ok := machineValue.(map[string]interface{}); ok {
				hasMachine = true
				v.validateMachineEssential(machineMap, result)
			}
		}
	}

	if !hasMachine {
		result.Errors = append(result.Errors, ValidationError{
			Field:    "machines",
			Message:  "at least one machine must be defined",
			Severity: "error",
		})
	}
}

// validateMachineEssential validates individual machine with essential rules
func (v *SimpleValidator) validateMachineEssential(machine map[string]interface{}, result *ValidationResult) {
	// Required fields
	if hostname, exists := machine["hostname"]; !exists || hostname == "" {
		result.Errors = append(result.Errors, ValidationError{
			Field:    "machine.hostname",
			Message:  "machine hostname is required",
			Severity: "error",
		})
	}

	// Port validation
	if port, exists := machine["port"]; exists {
		if portVal, ok := port.(float64); ok {
			if portVal < 1 || portVal > 65535 {
				result.Errors = append(result.Errors, ValidationError{
					Field:    "machine.port",
					Value:    portVal,
					Message:  "port must be between 1 and 65535",
					Severity: "error",
				})
			}
		}
	}

	// Basic authentication validation (keep only essential)
	if auth, exists := machine["authentication"]; exists {
		if authMap, ok := auth.(map[string]interface{}); ok {
			v.validateAuthenticationEssential(authMap, result)
		}
	}
}

// validateAuthenticationEssential validates authentication with essential rules only
func (v *SimpleValidator) validateAuthenticationEssential(auth map[string]interface{}, result *ValidationResult) {
	// Method validation
	if method, exists := auth["method"]; exists {
		if methodStr, ok := method.(string); ok {
			validMethods := map[string]bool{
				"password":    true,
				"publickey":   true,
				"certificate": true,
			}
			if !validMethods[methodStr] {
				result.Errors = append(result.Errors, ValidationError{
					Field:    "authentication.method",
					Value:    methodStr,
					Message:  "authentication method must be one of: password, publickey, certificate",
					Severity: "error",
				})
			}
		}
	}
}

// validateActionsEssential validates actions with essential rules only
func (v *SimpleValidator) validateActionsEssential(data map[string]interface{}, result *ValidationResult) {
	// Check for required actions block
	actionsBlock, exists := data["actions"]
	if !exists {
		result.Errors = append(result.Errors, ValidationError{
			Field:    "actions",
			Message:  "missing required actions block",
			Severity: "error",
		})
		return
	}

	// Validate actions block structure
	if actionsMap, ok := actionsBlock.(map[string]interface{}); ok {
		v.validateActionsBlockEssential(actionsMap, result)
	} else {
		result.Errors = append(result.Errors, ValidationError{
			Field:    "actions",
			Message:  "actions block must be a configuration block",
			Severity: "error",
		})
	}
}

// validateActionsBlockEssential validates actions block contents with essential rules
func (v *SimpleValidator) validateActionsBlockEssential(actions map[string]interface{}, result *ValidationResult) {
	// Check for at least one action
	hasAction := false
	for actionName, actionValue := range actions {
		if actionName != "action" {
			continue
		}

		if actionList, ok := actionValue.([]interface{}); ok {
			for _, action := range actionList {
				if actionMap, ok := action.(map[string]interface{}); ok {
					hasAction = true
					v.validateActionEssential(actionMap, result)
				}
			}
		} else if actionMap, ok := actionValue.(map[string]interface{}); ok {
			hasAction = true
			v.validateActionEssential(actionMap, result)
		}
	}

	if !hasAction {
		result.Errors = append(result.Errors, ValidationError{
			Field:    "actions",
			Message:  "at least one action must be defined",
			Severity: "error",
		})
	}
}

// validateActionEssential validates individual action with essential rules
func (v *SimpleValidator) validateActionEssential(action map[string]interface{}, result *ValidationResult) {
	// Required fields
	if description, exists := action["description"]; !exists || description == "" {
		result.Errors = append(result.Errors, ValidationError{
			Field:    "action.description",
			Message:  "action description is required",
			Severity: "error",
		})
	}

	if actionType, exists := action["type"]; !exists || actionType == "" {
		result.Errors = append(result.Errors, ValidationError{
			Field:    "action.type",
			Message:  "action type is required",
			Severity: "error",
		})
	}

	// Type-specific validation
	if actionType, exists := action["type"]; exists {
		if typeStr, ok := actionType.(string); ok {
			switch typeStr {
			case "command":
				if command, exists := action["command"]; !exists || command == "" {
					result.Errors = append(result.Errors, ValidationError{
						Field:    "action.command",
						Message:  "command is required for command actions",
						Severity: "error",
					})
				}
			}
		}
	}

	// Security validation (keep only essential)
	if command, exists := action["command"]; exists {
		if commandStr, ok := command.(string); ok {
			if strings.ContainsAny(commandStr, ";&|$") {
				result.Errors = append(result.Errors, ValidationError{
					Field:    "action.command",
					Value:    commandStr,
					Message:  "shell operators (;&|$) are not allowed in commands for security reasons",
					Severity: "error",
				})
			}
		}
	}

	// Range validation
	if timeout, exists := action["timeout"]; exists {
		if timeoutVal, ok := timeout.(float64); ok {
			if timeoutVal < 1 || timeoutVal > 3600 {
				result.Errors = append(result.Errors, ValidationError{
					Field:    "action.timeout",
					Value:    timeoutVal,
					Message:  "timeout must be between 1 and 3600 seconds",
					Severity: "error",
				})
			}
		}
	}
}

// validateVariablesEssential validates variables with essential rules only
func (v *SimpleValidator) validateVariablesEssential(data map[string]interface{}, result *ValidationResult) {
	// Check for required variables block
	variablesBlock, exists := data["variables"]
	if !exists {
		result.Errors = append(result.Errors, ValidationError{
			Field:    "variables",
			Message:  "missing required variables block",
			Severity: "error",
		})
		return
	}

	// Basic structure validation
	if variablesMap, ok := variablesBlock.(map[string]interface{}); ok {
		v.validateVariablesBlockEssential(variablesMap, result)
	} else {
		result.Errors = append(result.Errors, ValidationError{
			Field:    "variables",
			Message:  "variables block must be a configuration block",
			Severity: "error",
		})
	}
}

// validateVariablesBlockEssential validates variables block contents with essential rules
func (v *SimpleValidator) validateVariablesBlockEssential(variables map[string]interface{}, result *ValidationResult) {
	// Check for at least one variable
	hasVariable := false
	for varName, varValue := range variables {
		if varName == "variable" {
			if varList, ok := varValue.([]interface{}); ok {
				for _, variable := range varList {
					if varMap, ok := variable.(map[string]interface{}); ok {
						hasVariable = true
						v.validateVariableEssential(varMap, result)
					}
				}
			} else if varMap, ok := varValue.(map[string]interface{}); ok {
				hasVariable = true
				v.validateVariableEssential(varMap, result)
			}
		}
	}

	if !hasVariable {
		result.Errors = append(result.Errors, ValidationError{
			Field:    "variables",
			Message:  "at least one variable must be defined",
			Severity: "error",
		})
	}
}

// validateVariableEssential validates individual variable with essential rules
func (v *SimpleValidator) validateVariableEssential(variable map[string]interface{}, result *ValidationResult) {
	// Required fields
	if value, exists := variable["value"]; !exists || value == "" {
		result.Errors = append(result.Errors, ValidationError{
			Field:    "variable.value",
			Message:  "variable value is required",
			Severity: "error",
		})
	}

	// Encryption format validation (keep only essential)
	if encrypted, exists := variable["encrypted"]; exists && encrypted == true {
		if value, exists := variable["value"]; exists {
			if valueStr, ok := value.(string); ok {
				if !strings.HasPrefix(valueStr, "age1") {
					result.Errors = append(result.Errors, ValidationError{
						Field:    "variable.value",
						Value:    valueStr,
						Message:  "encrypted values must be in age format (age1...)",
						Severity: "error",
					})
				}
			}
		}
	}
}

// validateProjectDirectoryEssential validates project directory structure
func (v *SimpleValidator) validateProjectDirectoryEssential(data map[string]interface{}, result *ValidationResult) {
	// This would validate the actual directory structure
	// For now, just ensure the data structure is correct
	if targetDir, exists := data["target_directory"]; exists {
		if dirStr, ok := targetDir.(string); ok {
			// Check if directory exists
			if _, err := os.Stat(dirStr); os.IsNotExist(err) {
				result.Errors = append(result.Errors, ValidationError{
					Field:    "target_directory",
					Value:    dirStr,
					Message:  "target directory does not exist",
					Severity: "error",
				})
			}
		}
	}
}

// ParseHCLContent parses HCL content into structured data (simplified)
func (v *SimpleValidator) ParseHCLContent(content string) (map[string]interface{}, error) {
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

	// Convert AST to structured data (simplified)
	data := make(map[string]interface{})

	for _, block := range file.Body.(*hclsyntax.Body).Blocks {
		blockData := v.parseBlockSimple(block)

		// Handle project blocks with labels (project "name" {})
		if block.Type == "project" && len(block.Labels) > 0 {
			blockData["name"] = block.Labels[0]
		}

		data[block.Type] = blockData
	}

	return data, nil
}

// parseBlockSimple parses HCL blocks into structured data (simplified)
func (v *SimpleValidator) parseBlockSimple(block *hclsyntax.Block) map[string]interface{} {
	result := make(map[string]interface{})

	// Parse attributes (simplified)
	for name, attr := range block.Body.Attributes {
		value := v.parseAttributeValueSimple(attr.Expr)
		result[name] = value
	}

	// Parse nested blocks (simplified)
	for _, nestedBlock := range block.Body.Blocks {
		nestedData := v.parseBlockSimple(nestedBlock)

		// Handle blocks with labels (like "action "name" {")
		if len(nestedBlock.Labels) > 0 {
			label := nestedBlock.Labels[0]
			nestedData["name"] = label

			// For action blocks, store under "action" key
			if nestedBlock.Type == "action" {
				if existing, exists := result["action"]; exists {
					if arr, ok := existing.([]map[string]interface{}); ok {
						result["action"] = append(arr, nestedData)
					} else {
						result["action"] = []map[string]interface{}{existing.(map[string]interface{}), nestedData}
					}
				} else {
					result["action"] = nestedData
				}
			} else {
				// For other block types, use the label as key
				if existing, exists := result[label]; exists {
					if arr, ok := existing.([]map[string]interface{}); ok {
						result[label] = append(arr, nestedData)
					} else {
						result[label] = []map[string]interface{}{existing.(map[string]interface{}), nestedData}
					}
				} else {
					result[label] = nestedData
				}
			}
		} else {
			// Handle blocks without labels
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

// parseAttributeValueSimple converts HCL expression to Go value (simplified)
func (v *SimpleValidator) parseAttributeValueSimple(expr hclsyntax.Expression) interface{} {
	switch expr := expr.(type) {
	case *hclsyntax.LiteralValueExpr:
		return v.ctyValueToInterfaceSimple(expr.Val)
	case *hclsyntax.TemplateExpr:
		// Handle template expressions - try to extract literal values
		if len(expr.Parts) == 1 {
			if lit, ok := expr.Parts[0].(*hclsyntax.LiteralValueExpr); ok {
				return v.ctyValueToInterfaceSimple(lit.Val)
			}
		}
		// For complex templates, return as string
		return fmt.Sprintf("%v", expr)
	case *hclsyntax.TupleConsExpr:
		// Handle lists/arrays
		var result []interface{}
		for _, elem := range expr.Exprs {
			result = append(result, v.parseAttributeValueSimple(elem))
		}
		return result
	case *hclsyntax.ObjectConsExpr:
		// Handle maps/objects
		result := make(map[string]interface{})
		for _, item := range expr.Items {
			key := v.parseAttributeValueSimple(item.KeyExpr)
			value := v.parseAttributeValueSimple(item.ValueExpr)
			if keyStr, ok := key.(string); ok {
				result[keyStr] = value
			}
		}
		return result
	default:
		// For other expression types, return as string
		return fmt.Sprintf("%v", expr)
	}
}

// ctyValueToInterfaceSimple converts cty.Value to Go interface{} (simplified)
func (v *SimpleValidator) ctyValueToInterfaceSimple(val cty.Value) interface{} {
	if val.IsNull() {
		return nil
	}

	switch val.Type() {
	case cty.String:
		return val.AsString()
	case cty.Number:
		// Always return as float64 for simplicity
		floatVal, _ := val.AsBigFloat().Float64()
		return float64(floatVal)
	case cty.Bool:
		return val.True()
	default:
		// For complex types, try to iterate if possible
		if val.CanIterateElements() {
			var result []interface{}
			it := val.ElementIterator()
			for it.Next() {
				_, elemVal := it.Element()
				result = append(result, v.ctyValueToInterfaceSimple(elemVal))
			}
			return result
		}
		return fmt.Sprintf("%v", val)
	}
}
