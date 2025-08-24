package schemas

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/zclconf/go-cty/cty"
)

// StructValidator provides functionality to validate HCL content against Go struct schemas
type StructValidator struct {
	// Support for multiple schema versions
	supportedVersions []string
}

// NewStructValidator creates a new struct-based schema validator
func NewStructValidator() *StructValidator {
	return &StructValidator{
		supportedVersions: SupportedVersions,
	}
}

// ValidateHCLContent validates HCL content string against the appropriate schema
func (sv *StructValidator) ValidateHCLContent(schemaName, content string) (*ValidationResult, error) {
	// First parse HCL into structured data
	parsedData, err := sv.ParseHCLContent(content)
	if err != nil {
		return &ValidationResult{
			IsValid:    false,
			Errors:     []ValidationError{{Message: fmt.Sprintf("HCL parsing failed: %v", err), Severity: "error"}},
			Warnings:   []ValidationWarning{},
			SchemaName: schemaName,
		}, nil
	}

	// Then validate against the appropriate schema
	switch schemaName {
	case "project":
		return sv.ValidateProject(parsedData), nil
	case "machines":
		return sv.ValidateMachines(parsedData), nil
	case "actions":
		return sv.ValidateActions(parsedData), nil
	case "variables":
		return sv.ValidateVariables(parsedData), nil
	case "logging":
		return sv.ValidateLogging(parsedData), nil
	case "facts":
		return sv.ValidateFacts(parsedData), nil
	case "spooky":
		return sv.ValidateSpooky(parsedData), nil
	case "project-directory":
		return sv.ValidateProjectDirectory(parsedData), nil
	default:
		return &ValidationResult{
			IsValid:    false,
			Errors:     []ValidationError{{Message: fmt.Sprintf("Unknown schema: %s", schemaName), Severity: "error"}},
			Warnings:   []ValidationWarning{},
			SchemaName: schemaName,
		}, nil
	}
}

// ParseHCLContent parses HCL content into a structured map[string]interface{}
func (sv *StructValidator) ParseHCLContent(content string) (map[string]interface{}, error) {
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
func (sv *StructValidator) parseBlock(block *hclsyntax.Block) map[string]interface{} {
	result := make(map[string]interface{})

	// Handle block labels (e.g., action "name" {} -> name field)
	if len(block.Labels) > 0 {
		result["name"] = block.Labels[0]
	}

	// Special handling for authentication blocks
	if block.Type == "authentication" && len(block.Labels) > 0 {
		result["method"] = block.Labels[0]
	}

	// Parse attributes
	for name, attr := range block.Body.Attributes {
		result[name] = sv.parseExpression(attr.Expr)
	}

	// Parse nested blocks
	for _, nestedBlock := range block.Body.Blocks {
		blockData := sv.parseBlock(nestedBlock)

		// Handle multiple blocks of the same type
		if existing, exists := result[nestedBlock.Type]; exists {
			if list, ok := existing.([]interface{}); ok {
				result[nestedBlock.Type] = append(list, blockData)
			} else {
				// Convert single item to list
				result[nestedBlock.Type] = []interface{}{existing, blockData}
			}
		} else {
			result[nestedBlock.Type] = blockData
		}
	}

	return result
}

// parseExpression converts HCL expressions to Go values
func (sv *StructValidator) parseExpression(expr hclsyntax.Expression) interface{} {
	switch v := expr.(type) {
	case *hclsyntax.LiteralValueExpr:
		return sv.parseLiteralValue(v.Val)
	case *hclsyntax.TupleConsExpr:
		var result []interface{}
		for _, elem := range v.Exprs {
			result = append(result, sv.parseExpression(elem))
		}
		return result
	case *hclsyntax.ObjectConsExpr:
		result := make(map[string]interface{})
		for _, item := range v.Items {
			key := sv.parseExpression(item.KeyExpr)
			value := sv.parseExpression(item.ValueExpr)
			if keyStr, ok := key.(string); ok {
				result[keyStr] = value
			}
		}
		return result
	case *hclsyntax.ScopeTraversalExpr:
		// Handle variable references (e.g., var.name)
		var parts []string
		for _, part := range v.Traversal {
			if name, ok := part.(hcl.TraverseAttr); ok {
				parts = append(parts, name.Name)
			}
		}
		return strings.Join(parts, ".")
	case *hclsyntax.TemplateExpr:
		// Handle template expressions (e.g., "${var.name}")
		var parts []string
		for _, part := range v.Parts {
			if lit, ok := part.(*hclsyntax.LiteralValueExpr); ok {
				parts = append(parts, lit.Val.AsString())
			} else {
				// For embedded expressions, try to evaluate them
				parts = append(parts, fmt.Sprintf("${%v}", sv.parseExpression(part)))
			}
		}
		return strings.Join(parts, "")
	case *hclsyntax.BinaryOpExpr:
		// Handle binary operations (e.g., "hello" + "world")
		left := sv.parseExpression(v.LHS)
		right := sv.parseExpression(v.RHS)
		if leftStr, ok := left.(string); ok {
			if rightStr, ok := right.(string); ok {
				return leftStr + rightStr
			}
		}
		// For other operations, return a representation
		return fmt.Sprintf("%v %s %v", left, v.Op, right)
	case *hclsyntax.ConditionalExpr:
		// Handle conditional expressions
		condition := sv.parseExpression(v.Condition)
		trueVal := sv.parseExpression(v.TrueResult)
		falseVal := sv.parseExpression(v.FalseResult)
		return fmt.Sprintf("if %v then %v else %v", condition, trueVal, falseVal)
	case *hclsyntax.ForExpr:
		// Handle for expressions
		return fmt.Sprintf("for %s in %v: %v", v.KeyVar, sv.parseExpression(v.CollExpr), sv.parseExpression(v.ValExpr))
	case *hclsyntax.SplatExpr:
		// Handle splat expressions
		return fmt.Sprintf("%v[*]", sv.parseExpression(v.Source))
	case *hclsyntax.IndexExpr:
		// Handle index expressions
		collection := sv.parseExpression(v.Collection)
		key := sv.parseExpression(v.Key)
		return fmt.Sprintf("%v[%v]", collection, key)
	case *hclsyntax.RelativeTraversalExpr:
		// Handle relative traversal expressions
		source := sv.parseExpression(v.Source)
		var parts []string
		for _, part := range v.Traversal {
			if name, ok := part.(hcl.TraverseAttr); ok {
				parts = append(parts, name.Name)
			}
		}
		return fmt.Sprintf("%v.%s", source, strings.Join(parts, "."))
	case *hclsyntax.FunctionCallExpr:
		// Handle function calls
		args := make([]string, len(v.Args))
		for i, arg := range v.Args {
			args[i] = fmt.Sprintf("%v", sv.parseExpression(arg))
		}
		return fmt.Sprintf("%s(%s)", v.Name, strings.Join(args, ", "))
	case *hclsyntax.ParenthesesExpr:
		// Handle parenthesized expressions
		return sv.parseExpression(v.Expression)
	default:
		// For any other expression types, try to get a string representation
		// This is a fallback for unknown expression types
		return fmt.Sprintf("{{%s}}", expr.Range().String())
	}
}

// parseLiteralValue converts HCL literal values to Go values
func (sv *StructValidator) parseLiteralValue(val cty.Value) interface{} {
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
		return fmt.Sprintf("{{%s}}", val.Type().GoString())
	}
}

// ValidateProject validates project configuration against the latest schema version
func (sv *StructValidator) ValidateProject(content map[string]interface{}) *ValidationResult {
	result := &ValidationResult{
		IsValid:    true,
		Errors:     []ValidationError{},
		Warnings:   []ValidationWarning{},
		SchemaName: "project",
	}

	// Always validate against the latest version for project init
	sv.validateProjectV1(content, result)

	// Update validity based on errors
	if len(result.Errors) > 0 {
		result.IsValid = false
	}

	return result
}

// ValidateMachines validates machines configuration against supported schema versions
func (sv *StructValidator) ValidateMachines(content map[string]interface{}) *ValidationResult {
	result := &ValidationResult{
		IsValid:    true,
		Errors:     []ValidationError{},
		Warnings:   []ValidationWarning{},
		SchemaName: "machines",
	}

	// Validate against all supported versions
	for _, version := range sv.supportedVersions {
		switch version {
		case "1":
			sv.validateMachinesV1(content, result)
		}
	}

	// Update validity based on errors
	if len(result.Errors) > 0 {
		result.IsValid = false
	}

	return result
}

// ValidateActions validates actions configuration against supported schema versions
func (sv *StructValidator) ValidateActions(content map[string]interface{}) *ValidationResult {
	result := &ValidationResult{
		IsValid:    true,
		Errors:     []ValidationError{},
		Warnings:   []ValidationWarning{},
		SchemaName: "actions",
	}

	// Validate against all supported versions
	for _, version := range sv.supportedVersions {
		switch version {
		case "1":
			sv.validateActionsV1(content, result)
		}
	}

	// Update validity based on errors
	if len(result.Errors) > 0 {
		result.IsValid = false
	}

	return result
}

// ValidateVariables validates variables configuration against supported schema versions
func (sv *StructValidator) ValidateVariables(content map[string]interface{}) *ValidationResult {
	result := &ValidationResult{
		IsValid:    true,
		Errors:     []ValidationError{},
		Warnings:   []ValidationWarning{},
		SchemaName: "variables",
	}

	// Validate against all supported versions
	for _, version := range sv.supportedVersions {
		switch version {
		case "1":
			sv.validateVariablesV1(content, result)
		}
	}

	// Update validity based on errors
	if len(result.Errors) > 0 {
		result.IsValid = false
	}

	return result
}

// validateProjectV1 validates project configuration against V1 schema
func (sv *StructValidator) validateProjectV1(content map[string]interface{}, result *ValidationResult) {
	// Check for required project block
	projectBlock, exists := content["project"]
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
		sv.validateProjectBlockV1(projectMap, result)
	} else {
		result.Errors = append(result.Errors, ValidationError{
			Field:    "project",
			Message:  "project block must be a configuration block",
			Severity: "error",
		})
	}
}

// validateProjectBlockV1 validates the project block contents against V1 schema
func (sv *StructValidator) validateProjectBlockV1(project map[string]interface{}, result *ValidationResult) {
	// Required fields
	if name, exists := project["name"]; !exists {
		result.Errors = append(result.Errors, ValidationError{
			Field:    "project.name",
			Message:  "missing required field: name",
			Severity: "error",
		})
	} else {
		sv.validateProjectNameV1(name, result)
	}

	if description, exists := project["description"]; !exists {
		result.Errors = append(result.Errors, ValidationError{
			Field:    "project.description",
			Message:  "missing required field: description",
			Severity: "error",
		})
	} else {
		sv.validateProjectDescriptionV1(description, result)
	}

	// Optional fields with validation
	if runMaxParallel, exists := project["run_max_parallel"]; exists {
		sv.validateProjectRunMaxParallelV1(runMaxParallel, result)
	}

	if factsTimeout, exists := project["facts_timeout"]; exists {
		sv.validateProjectFactsTimeoutV1(factsTimeout, result)
	}
}

// validateProjectNameV1 validates project name against V1 rules
func (sv *StructValidator) validateProjectNameV1(name interface{}, result *ValidationResult) {
	nameStr, ok := name.(string)
	if !ok {
		result.Errors = append(result.Errors, ValidationError{
			Field:    "project.name",
			Message:  "name must be a string",
			Severity: "error",
		})
		return
	}

	// Pattern validation: ^[a-zA-Z][a-zA-Z0-9._-]*$
	pattern := regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9._-]*$`)
	if !pattern.MatchString(nameStr) {
		result.Errors = append(result.Errors, ValidationError{
			Field:    "project.name",
			Value:    nameStr,
			Message:  "name must start with a letter and contain only alphanumeric characters, dots, underscores, and hyphens",
			Severity: "error",
		})
	}

	// Length validation: 1-128 characters
	if len(nameStr) < 1 || len(nameStr) > 128 {
		result.Errors = append(result.Errors, ValidationError{
			Field:    "project.name",
			Value:    nameStr,
			Message:  "name must be between 1 and 128 characters",
			Severity: "error",
		})
	}
}

// validateProjectDescriptionV1 validates project description against V1 rules
func (sv *StructValidator) validateProjectDescriptionV1(description interface{}, result *ValidationResult) {
	descStr, ok := description.(string)
	if !ok {
		result.Errors = append(result.Errors, ValidationError{
			Field:    "project.description",
			Message:  "description must be a string",
			Severity: "error",
		})
		return
	}

	// Length validation: max 1024 characters
	if len(descStr) > 1024 {
		result.Errors = append(result.Errors, ValidationError{
			Field:    "project.description",
			Value:    descStr,
			Message:  "description must not exceed 1024 characters",
			Severity: "error",
		})
	}
}

// validateProjectRunMaxParallelV1 validates run_max_parallel against V1 rules
func (sv *StructValidator) validateProjectRunMaxParallelV1(value interface{}, result *ValidationResult) {
	switch v := value.(type) {
	case int:
		if v < 1 || v > 100 {
			result.Errors = append(result.Errors, ValidationError{
				Field:    "project.run_max_parallel",
				Value:    v,
				Message:  "run_max_parallel must be between 1 and 100",
				Severity: "error",
			})
		}
	case float64:
		// HCL might parse numbers as float64
		if v < 1 || v > 100 {
			result.Errors = append(result.Errors, ValidationError{
				Field:    "project.run_max_parallel",
				Value:    v,
				Message:  "run_max_parallel must be between 1 and 100",
				Severity: "error",
			})
		}
	default:
		result.Errors = append(result.Errors, ValidationError{
			Field:    "project.run_max_parallel",
			Value:    value,
			Message:  "run_max_parallel must be a number",
			Severity: "error",
		})
	}
}

// validateProjectFactsTimeoutV1 validates facts_timeout against V1 rules
func (sv *StructValidator) validateProjectFactsTimeoutV1(value interface{}, result *ValidationResult) {
	switch v := value.(type) {
	case int:
		if v < 1 || v > 3600 {
			result.Errors = append(result.Errors, ValidationError{
				Field:    "project.facts_timeout",
				Value:    v,
				Message:  "facts_timeout must be between 1 and 3600 seconds",
				Severity: "error",
			})
		}
	case float64:
		// HCL might parse numbers as float64
		if v < 1 || v > 3600 {
			result.Errors = append(result.Errors, ValidationError{
				Field:    "project.facts_timeout",
				Value:    v,
				Message:  "facts_timeout must be between 1 and 3600 seconds",
				Severity: "error",
			})
		}
	default:
		result.Errors = append(result.Errors, ValidationError{
			Field:    "project.facts_timeout",
			Value:    value,
			Message:  "facts_timeout must be a number",
			Severity: "error",
		})
	}
}

// validateMachinesV1 validates machines configuration against V1 schema
func (sv *StructValidator) validateMachinesV1(content map[string]interface{}, result *ValidationResult) {
	// Check for required machines block
	machinesBlock, exists := content["machines"]
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
		sv.validateMachinesBlockV1(machinesMap, result)
	} else {
		result.Errors = append(result.Errors, ValidationError{
			Field:    "machines",
			Message:  "machines block must be a configuration block",
			Severity: "error",
		})
	}
}

// validateMachinesBlockV1 validates the machines block contents against V1 schema
func (sv *StructValidator) validateMachinesBlockV1(machines map[string]interface{}, result *ValidationResult) {
	// Check for at least one machine or group
	hasMachine := false
	hasGroup := false

	if machineList, exists := machines["machine"]; exists {
		hasMachine = true
		if machineArray, ok := machineList.([]interface{}); ok {
			for i, machine := range machineArray {
				if machineMap, ok := machine.(map[string]interface{}); ok {
					sv.validateMachineV1(machineMap, fmt.Sprintf("machines.machine[%d]", i), result)
				}
			}
		}
	}

	if groupList, exists := machines["group"]; exists {
		hasGroup = true
		if groupArray, ok := groupList.([]interface{}); ok {
			for i, group := range groupArray {
				if groupMap, ok := group.(map[string]interface{}); ok {
					sv.validateGroupV1(groupMap, fmt.Sprintf("machines.group[%d]", i), result)
				}
			}
		}
	}

	if !hasMachine && !hasGroup {
		result.Warnings = append(result.Warnings, ValidationWarning{
			Field:   "machines",
			Message: "machines block should contain at least one machine or group",
		})
	}
}

// validateMachineV1 validates an individual machine against V1 schema
func (sv *StructValidator) validateMachineV1(machine map[string]interface{}, path string, result *ValidationResult) {
	// Required fields
	if name, exists := machine["name"]; !exists {
		result.Errors = append(result.Errors, ValidationError{
			Field:    fmt.Sprintf("%s.name", path),
			Message:  "missing required field: name",
			Severity: "error",
		})
	} else {
		sv.validateMachineNameV1(name, path, result)
	}

	if hostname, exists := machine["hostname"]; !exists {
		result.Errors = append(result.Errors, ValidationError{
			Field:    fmt.Sprintf("%s.hostname", path),
			Message:  "missing required field: hostname",
			Severity: "error",
		})
	} else {
		sv.validateMachineHostnameV1(hostname, path, result)
	}

	if user, exists := machine["user"]; !exists {
		result.Errors = append(result.Errors, ValidationError{
			Field:    fmt.Sprintf("%s.user", path),
			Message:  "missing required field: user",
			Severity: "error",
		})
	} else {
		sv.validateMachineUserV1(user, path, result)
	}

	if auth, exists := machine["authentication"]; !exists {
		result.Errors = append(result.Errors, ValidationError{
			Field:    fmt.Sprintf("%s.authentication", path),
			Message:  "missing required field: authentication",
			Severity: "error",
		})
	} else {
		sv.validateMachineAuthenticationV1(auth, path, result)
	}
}

// validateMachineNameV1 validates machine name against V1 rules
func (sv *StructValidator) validateMachineNameV1(name interface{}, path string, result *ValidationResult) {
	nameStr, ok := name.(string)
	if !ok {
		result.Errors = append(result.Errors, ValidationError{
			Field:    fmt.Sprintf("%s.name", path),
			Message:  "name must be a string",
			Severity: "error",
		})
		return
	}

	// Pattern validation: ^[a-zA-Z0-9_.-]+$
	pattern := regexp.MustCompile(`^[a-zA-Z0-9_.-]+$`)
	if !pattern.MatchString(nameStr) {
		result.Errors = append(result.Errors, ValidationError{
			Field:    fmt.Sprintf("%s.name", path),
			Value:    nameStr,
			Message:  "name must contain only alphanumeric characters, dots, underscores, and hyphens",
			Severity: "error",
		})
	}

	// Length validation: 1-64 characters
	if len(nameStr) < 1 || len(nameStr) > 64 {
		result.Errors = append(result.Errors, ValidationError{
			Field:    fmt.Sprintf("%s.name", path),
			Value:    nameStr,
			Message:  "name must be between 1 and 64 characters",
			Severity: "error",
		})
	}
}

// validateMachineHostnameV1 validates machine hostname against V1 rules
func (sv *StructValidator) validateMachineHostnameV1(hostname interface{}, path string, result *ValidationResult) {
	hostnameStr, ok := hostname.(string)
	if !ok {
		result.Errors = append(result.Errors, ValidationError{
			Field:    fmt.Sprintf("%s.hostname", path),
			Message:  "hostname must be a string",
			Severity: "error",
		})
		return
	}

	// Length validation: 1-253 characters
	if len(hostnameStr) < 1 || len(hostnameStr) > 253 {
		result.Errors = append(result.Errors, ValidationError{
			Field:    fmt.Sprintf("%s.hostname", path),
			Value:    hostnameStr,
			Message:  "hostname must be between 1 and 253 characters",
			Severity: "error",
		})
	}
}

// validateMachineUserV1 validates machine user against V1 rules
func (sv *StructValidator) validateMachineUserV1(user interface{}, path string, result *ValidationResult) {
	userStr, ok := user.(string)
	if !ok {
		result.Errors = append(result.Errors, ValidationError{
			Field:    fmt.Sprintf("%s.user", path),
			Message:  "user must be a string",
			Severity: "error",
		})
		return
	}

	// Pattern validation: ^[a-zA-Z0-9_.-]+$
	pattern := regexp.MustCompile(`^[a-zA-Z0-9_.-]+$`)
	if !pattern.MatchString(userStr) {
		result.Errors = append(result.Errors, ValidationError{
			Field:    fmt.Sprintf("%s.user", path),
			Value:    userStr,
			Message:  "user must contain only alphanumeric characters, dots, underscores, and hyphens",
			Severity: "error",
		})
	}

	// Length validation: 1-32 characters
	if len(userStr) < 1 || len(userStr) > 32 {
		result.Errors = append(result.Errors, ValidationError{
			Field:    fmt.Sprintf("%s.user", path),
			Value:    userStr,
			Message:  "user must be between 1 and 32 characters",
			Severity: "error",
		})
	}
}

// validateMachineAuthenticationV1 validates machine authentication against V1 schema
func (sv *StructValidator) validateMachineAuthenticationV1(auth interface{}, path string, result *ValidationResult) {
	authMap, ok := auth.(map[string]interface{})
	if !ok {
		result.Errors = append(result.Errors, ValidationError{
			Field:    fmt.Sprintf("%s.authentication", path),
			Message:  "authentication must be a configuration block",
			Severity: "error",
		})
		return
	}

	// Required method field
	if method, exists := authMap["method"]; !exists {
		result.Errors = append(result.Errors, ValidationError{
			Field:    fmt.Sprintf("%s.authentication.method", path),
			Message:  "missing required field: method",
			Severity: "error",
		})
	} else {
		sv.validateAuthMethodV1(method, path, result)
	}

	// Validate method-specific fields
	if method, exists := authMap["method"]; exists {
		if methodStr, ok := method.(string); ok {
			switch methodStr {
			case "publickey":
				sv.validatePublicKeyAuthV1(authMap, path, result)
			case "password":
				sv.validatePasswordAuthV1(authMap, path, result)
			case "certificate":
				sv.validateCertificateAuthV1(authMap, path, result)
			}
		}
	}
}

// validateAuthMethodV1 validates authentication method against V1 rules
func (sv *StructValidator) validateAuthMethodV1(method interface{}, path string, result *ValidationResult) {
	methodStr, ok := method.(string)
	if !ok {
		result.Errors = append(result.Errors, ValidationError{
			Field:    fmt.Sprintf("%s.authentication.method", path),
			Message:  "method must be a string",
			Severity: "error",
		})
		return
	}

	// Enum validation
	validMethods := []string{"publickey", "password", "certificate"}
	valid := false
	for _, validMethod := range validMethods {
		if methodStr == validMethod {
			valid = true
			break
		}
	}

	if !valid {
		result.Errors = append(result.Errors, ValidationError{
			Field:    fmt.Sprintf("%s.authentication.method", path),
			Value:    methodStr,
			Message:  fmt.Sprintf("method must be one of: %s", strings.Join(validMethods, ", ")),
			Severity: "error",
		})
	}
}

// validatePublicKeyAuthV1 validates public key authentication fields
func (sv *StructValidator) validatePublicKeyAuthV1(auth map[string]interface{}, path string, result *ValidationResult) {
	if publicKeyPath, exists := auth["public_key_path"]; !exists {
		result.Errors = append(result.Errors, ValidationError{
			Field:    fmt.Sprintf("%s.authentication.public_key_path", path),
			Message:  "public_key_path is required when method is 'publickey'",
			Severity: "error",
		})
	} else {
		// Validate public_key_path is a string
		if _, ok := publicKeyPath.(string); !ok {
			result.Errors = append(result.Errors, ValidationError{
				Field:    fmt.Sprintf("%s.authentication.public_key_path", path),
				Message:  "public_key_path must be a string",
				Severity: "error",
			})
		}
	}
}

// validatePasswordAuthV1 validates password authentication fields
func (sv *StructValidator) validatePasswordAuthV1(auth map[string]interface{}, path string, result *ValidationResult) {
	// Password field is required for password method
	if password, exists := auth["password"]; !exists {
		result.Errors = append(result.Errors, ValidationError{
			Field:    fmt.Sprintf("%s.authentication.password", path),
			Message:  "password is required when method is 'password'",
			Severity: "error",
		})
	} else {
		// Validate password is a configuration block
		if _, ok := password.(map[string]interface{}); !ok {
			result.Errors = append(result.Errors, ValidationError{
				Field:    fmt.Sprintf("%s.authentication.password", path),
				Message:  "password must be a configuration block",
				Severity: "error",
			})
		}
	}
}

// validateCertificateAuthV1 validates certificate authentication fields
func (sv *StructValidator) validateCertificateAuthV1(auth map[string]interface{}, path string, result *ValidationResult) {
	// Certificate authentication requires both private key and certificate
	if privateKeyPath, exists := auth["private_key_path"]; !exists {
		result.Errors = append(result.Errors, ValidationError{
			Field:    fmt.Sprintf("%s.authentication.private_key_path", path),
			Message:  "private_key_path is required when method is 'certificate'",
			Severity: "error",
		})
	} else {
		// Validate private_key_path is a string
		if _, ok := privateKeyPath.(string); !ok {
			result.Errors = append(result.Errors, ValidationError{
				Field:    fmt.Sprintf("%s.authentication.private_key_path", path),
				Message:  "private_key_path must be a string",
				Severity: "error",
			})
		}
	}

	if certPath, exists := auth["certificate_path"]; !exists {
		result.Errors = append(result.Errors, ValidationError{
			Field:    fmt.Sprintf("%s.authentication.certificate_path", path),
			Message:  "certificate_path is required when method is 'certificate'",
			Severity: "error",
		})
	} else {
		// Validate certificate_path is a string
		if _, ok := certPath.(string); !ok {
			result.Errors = append(result.Errors, ValidationError{
				Field:    fmt.Sprintf("%s.authentication.certificate_path", path),
				Message:  "certificate_path must be a string",
				Severity: "error",
			})
		}
	}
}

// validateGroupV1 validates a machine group against V1 schema
func (sv *StructValidator) validateGroupV1(group map[string]interface{}, path string, result *ValidationResult) {
	// Required fields
	if name, exists := group["name"]; !exists {
		result.Errors = append(result.Errors, ValidationError{
			Field:    fmt.Sprintf("%s.name", path),
			Message:  "missing required field: name",
			Severity: "error",
		})
	} else {
		sv.validateMachineNameV1(name, path, result) // Reuse machine name validation
	}

	if machines, exists := group["machines"]; !exists {
		result.Errors = append(result.Errors, ValidationError{
			Field:    fmt.Sprintf("%s.machines", path),
			Message:  "missing required field: machines",
			Severity: "error",
		})
	} else {
		sv.validateGroupMachinesV1(machines, path, result)
	}
}

// validateGroupMachinesV1 validates group machines list against V1 rules
func (sv *StructValidator) validateGroupMachinesV1(machines interface{}, path string, result *ValidationResult) {
	machinesArray, ok := machines.([]interface{})
	if !ok {
		result.Errors = append(result.Errors, ValidationError{
			Field:    fmt.Sprintf("%s.machines", path),
			Message:  "machines must be a list",
			Severity: "error",
		})
		return
	}

	// Must have at least one machine
	if len(machinesArray) < 1 {
		result.Errors = append(result.Errors, ValidationError{
			Field:    fmt.Sprintf("%s.machines", path),
			Message:  "machines list must contain at least one machine",
			Severity: "error",
		})
		return
	}

	// Validate each machine name in the list
	for i, machine := range machinesArray {
		if machineStr, ok := machine.(string); ok {
			// Validate machine name format
			pattern := regexp.MustCompile(`^[a-zA-Z0-9_.-]+$`)
			if !pattern.MatchString(machineStr) {
				result.Errors = append(result.Errors, ValidationError{
					Field:    fmt.Sprintf("%s.machines[%d]", path, i),
					Value:    machineStr,
					Message:  "machine name must contain only alphanumeric characters, dots, underscores, and hyphens",
					Severity: "error",
				})
			}

			// Validate machine name length
			if len(machineStr) < 1 || len(machineStr) > 64 {
				result.Errors = append(result.Errors, ValidationError{
					Field:    fmt.Sprintf("%s.machines[%d]", path, i),
					Value:    machineStr,
					Message:  "machine name must be between 1 and 64 characters",
					Severity: "error",
				})
			}
		} else {
			result.Errors = append(result.Errors, ValidationError{
				Field:    fmt.Sprintf("%s.machines[%d]", path, i),
				Message:  "machine name must be a string",
				Severity: "error",
			})
		}
	}
}

// validateActionsV1 validates actions configuration against V1 schema
func (sv *StructValidator) validateActionsV1(content map[string]interface{}, result *ValidationResult) {
	// Check for required actions block
	actionsBlock, exists := content["actions"]
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
		sv.validateActionsBlockV1(actionsMap, result)
	} else {
		result.Errors = append(result.Errors, ValidationError{
			Field:    "actions",
			Message:  "actions block must be a configuration block",
			Severity: "error",
		})
	}
}

// validateActionsBlockV1 validates the actions block contents against V1 schema
func (sv *StructValidator) validateActionsBlockV1(actions map[string]interface{}, result *ValidationResult) {
	// Check for action list
	if actionList, exists := actions["action"]; exists {
		if actionArray, ok := actionList.([]interface{}); ok {
			for i, action := range actionArray {
				if actionMap, ok := action.(map[string]interface{}); ok {
					sv.validateActionV1(actionMap, fmt.Sprintf("actions.action[%d]", i), result)
				}
			}
		}
	} else {
		result.Warnings = append(result.Warnings, ValidationWarning{
			Field:   "actions",
			Message: "actions block should contain at least one action",
		})
	}
}

// validateActionV1 validates an individual action against V1 schema
func (sv *StructValidator) validateActionV1(action map[string]interface{}, path string, result *ValidationResult) {
	// Required fields
	if name, exists := action["name"]; !exists {
		result.Errors = append(result.Errors, ValidationError{
			Field:    fmt.Sprintf("%s.name", path),
			Message:  "missing required field: name",
			Severity: "error",
		})
	} else {
		sv.validateActionNameV1(name, path, result)
	}

	if description, exists := action["description"]; !exists {
		result.Errors = append(result.Errors, ValidationError{
			Field:    fmt.Sprintf("%s.description", path),
			Message:  "missing required field: description",
			Severity: "error",
		})
	} else {
		sv.validateActionDescriptionV1(description, path, result)
	}

	if actionType, exists := action["type"]; !exists {
		result.Errors = append(result.Errors, ValidationError{
			Field:    fmt.Sprintf("%s.type", path),
			Message:  "missing required field: type",
			Severity: "error",
		})
	} else {
		sv.validateActionTypeV1(actionType, path, result)
	}
}

// validateActionNameV1 validates action name against V1 rules
func (sv *StructValidator) validateActionNameV1(name interface{}, path string, result *ValidationResult) {
	nameStr, ok := name.(string)
	if !ok {
		result.Errors = append(result.Errors, ValidationError{
			Field:    fmt.Sprintf("%s.name", path),
			Message:  "name must be a string",
			Severity: "error",
		})
		return
	}

	// Pattern validation: ^[a-zA-Z0-9_.-]+$
	pattern := regexp.MustCompile(`^[a-zA-Z0-9_.-]+$`)
	if !pattern.MatchString(nameStr) {
		result.Errors = append(result.Errors, ValidationError{
			Field:    fmt.Sprintf("%s.name", path),
			Value:    nameStr,
			Message:  "name must contain only alphanumeric characters, dots, underscores, and hyphens",
			Severity: "error",
		})
	}

	// Length validation: 1-64 characters
	if len(nameStr) < 1 || len(nameStr) > 64 {
		result.Errors = append(result.Errors, ValidationError{
			Field:    fmt.Sprintf("%s.name", path),
			Value:    nameStr,
			Message:  "name must be between 1 and 64 characters",
			Severity: "error",
		})
	}
}

// validateActionDescriptionV1 validates action description against V1 rules
func (sv *StructValidator) validateActionDescriptionV1(description interface{}, path string, result *ValidationResult) {
	descStr, ok := description.(string)
	if !ok {
		result.Errors = append(result.Errors, ValidationError{
			Field:    fmt.Sprintf("%s.description", path),
			Message:  "description must be a string",
			Severity: "error",
		})
		return
	}

	// Length validation: 1-500 characters
	if len(descStr) < 1 || len(descStr) > 500 {
		result.Errors = append(result.Errors, ValidationError{
			Field:    fmt.Sprintf("%s.description", path),
			Value:    descStr,
			Message:  "description must be between 1 and 500 characters",
			Severity: "error",
		})
	}
}

// validateActionTypeV1 validates action type against V1 rules
func (sv *StructValidator) validateActionTypeV1(actionType interface{}, path string, result *ValidationResult) {
	typeStr, ok := actionType.(string)
	if !ok {
		result.Errors = append(result.Errors, ValidationError{
			Field:    fmt.Sprintf("%s.type", path),
			Message:  "type must be a string",
			Severity: "error",
		})
		return
	}

	// Enum validation
	validTypes := []string{"command", "script", "template_deploy", "file_sync", "service_control"}
	valid := false
	for _, validType := range validTypes {
		if typeStr == validType {
			valid = true
			break
		}
	}

	if !valid {
		result.Errors = append(result.Errors, ValidationError{
			Field:    fmt.Sprintf("%s.type", path),
			Value:    typeStr,
			Message:  fmt.Sprintf("type must be one of: %s", strings.Join(validTypes, ", ")),
			Severity: "error",
		})
	}
}

// validateVariablesV1 validates variables configuration against V1 schema
func (sv *StructValidator) validateVariablesV1(content map[string]interface{}, result *ValidationResult) {
	// Check for required variables block
	variablesBlock, exists := content["variables"]
	if !exists {
		result.Errors = append(result.Errors, ValidationError{
			Field:    "variables",
			Message:  "missing required variables block",
			Severity: "error",
		})
		return
	}

	// Validate variables block structure
	if variablesMap, ok := variablesBlock.(map[string]interface{}); ok {
		sv.validateVariablesBlockV1(variablesMap, result)
	} else {
		result.Errors = append(result.Errors, ValidationError{
			Field:    "variables",
			Message:  "variables block must be a configuration block",
			Severity: "error",
		})
	}
}

// validateVariablesBlockV1 validates the variables block contents against V1 schema
func (sv *StructValidator) validateVariablesBlockV1(variables map[string]interface{}, result *ValidationResult) {
	// Check for variable list
	if variableList, exists := variables["variable"]; exists {
		if variableArray, ok := variableList.([]interface{}); ok {
			for i, variable := range variableArray {
				if variableMap, ok := variable.(map[string]interface{}); ok {
					sv.validateVariableV1(variableMap, fmt.Sprintf("variables.variable[%d]", i), result)
				}
			}
		}
	} else {
		result.Warnings = append(result.Warnings, ValidationWarning{
			Field:   "variables",
			Message: "variables block should contain at least one variable",
		})
	}
}

// validateVariableV1 validates an individual variable against V1 schema
func (sv *StructValidator) validateVariableV1(variable map[string]interface{}, path string, result *ValidationResult) {
	// Required fields
	if name, exists := variable["name"]; !exists {
		result.Errors = append(result.Errors, ValidationError{
			Field:    fmt.Sprintf("%s.name", path),
			Message:  "missing required field: name",
			Severity: "error",
		})
	} else {
		sv.validateVariableNameV1(name, path, result)
	}

	if value, exists := variable["value"]; !exists {
		result.Errors = append(result.Errors, ValidationError{
			Field:    fmt.Sprintf("%s.value", path),
			Message:  "missing required field: value",
			Severity: "error",
		})
	} else {
		sv.validateVariableValueV1(value, path, result)
	}
}

// validateVariableNameV1 validates variable name against V1 rules
func (sv *StructValidator) validateVariableNameV1(name interface{}, path string, result *ValidationResult) {
	nameStr, ok := name.(string)
	if !ok {
		result.Errors = append(result.Errors, ValidationError{
			Field:    fmt.Sprintf("%s.name", path),
			Message:  "name must be a string",
			Severity: "error",
		})
		return
	}

	// Pattern validation: ^[a-zA-Z0-9_.-]+$
	pattern := regexp.MustCompile(`^[a-zA-Z0-9_.-]+$`)
	if !pattern.MatchString(nameStr) {
		result.Errors = append(result.Errors, ValidationError{
			Field:    fmt.Sprintf("%s.name", path),
			Value:    nameStr,
			Message:  "name must contain only alphanumeric characters, dots, underscores, and hyphens",
			Severity: "error",
		})
	}

	// Length validation: 1-64 characters
	if len(nameStr) < 1 || len(nameStr) > 64 {
		result.Errors = append(result.Errors, ValidationError{
			Field:    fmt.Sprintf("%s.name", path),
			Value:    nameStr,
			Message:  "name must be between 1 and 64 characters",
			Severity: "error",
		})
	}
}

// validateVariableValueV1 validates variable value against V1 rules
func (sv *StructValidator) validateVariableValueV1(value interface{}, path string, result *ValidationResult) {
	// Value can be any type, just ensure it's not nil
	if value == nil {
		result.Errors = append(result.Errors, ValidationError{
			Field:    fmt.Sprintf("%s.value", path),
			Message:  "value cannot be null",
			Severity: "error",
		})
	}
}

// Helper function to extract validation rules from struct tags
func (sv *StructValidator) extractValidationRules(structType reflect.Type) map[string]ValidationRuleV1 {
	rules := make(map[string]ValidationRuleV1)

	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		tag := field.Tag.Get("json")
		if tag == "" {
			continue
		}

		rule := ValidationRuleV1{}

		// Extract validation rules from struct tags
		if required := field.Tag.Get("required"); required != "" {
			if required == "true" {
				rule.Required = true
			}
		}

		if pattern := field.Tag.Get("pattern"); pattern != "" {
			rule.Pattern = pattern
		}

		if min := field.Tag.Get("min"); min != "" {
			if minVal, err := strconv.Atoi(min); err == nil {
				rule.Min = minVal
			}
		}

		if max := field.Tag.Get("max"); max != "" {
			if maxVal, err := strconv.Atoi(max); err == nil {
				rule.Max = maxVal
			}
		}

		if minLength := field.Tag.Get("min_length"); minLength != "" {
			if minLen, err := strconv.Atoi(minLength); err == nil {
				rule.MinLength = minLen
			}
		}

		if maxLength := field.Tag.Get("max_length"); maxLength != "" {
			if maxLen, err := strconv.Atoi(maxLength); err == nil {
				rule.MaxLength = maxLen
			}
		}

		if enum := field.Tag.Get("enum"); enum != "" {
			rule.Enum = enum
		}

		if defaultValue := field.Tag.Get("default"); defaultValue != "" {
			rule.Default = defaultValue
		}

		if len(rule.Pattern) > 0 || rule.Required || rule.Min > 0 || rule.Max > 0 ||
			rule.MinLength > 0 || rule.MaxLength > 0 || len(rule.Enum) > 0 || len(rule.Default) > 0 {
			rules[tag] = rule
		}
	}

	return rules
}

// ============================================================================
// CONFIGURATION GENERATION FUNCTIONS
// ============================================================================

// GenerateProjectConfigFromStructs generates project.hcl content from Go structs
func GenerateProjectConfigFromStructs(name, description string) string {
	var content strings.Builder

	content.WriteString("# Spooky Project Configuration\n")
	content.WriteString("# Generated from Go struct schemas\n\n")

	content.WriteString("project {\n")
	content.WriteString(fmt.Sprintf("  name = \"%s\"\n", name))
	content.WriteString(fmt.Sprintf("  description = \"%s\"\n", description))

	// Add optional fields with defaults and descriptions
	content.WriteString("  # run_max_parallel = 10  # Maximum parallel action executions for this project\n")
	content.WriteString("  # run_dry_run_default = false  # Default dry-run mode for actions\n")
	content.WriteString("  # run_validate_before_run = true  # Validate project configuration before running\n")
	content.WriteString("  # run_backup_before_changes = false  # Create backups before making changes\n")
	content.WriteString("  # facts_timeout = 30  # Timeout for facts collection in seconds\n")
	content.WriteString("  # facts_parallel_collection = 10  # Number of parallel facts collection workers\n")
	content.WriteString("  # facts_retry_attempts = 3  # Number of retry attempts for failed facts collection\n")
	content.WriteString("  # facts_retry_delay = 5  # Delay between retry attempts in seconds\n")

	content.WriteString("}\n")

	return content.String()
}

// GenerateMachinesConfigFromStructs generates machines.hcl content from Go structs
func GenerateMachinesConfigFromStructs() string {
	var content strings.Builder

	content.WriteString("# Spooky Machines Configuration\n")
	content.WriteString("# Generated from Go struct schemas\n\n")

	content.WriteString("machines {\n")
	content.WriteString("  # Add your machine configurations here\n")
	content.WriteString("  # machine \"server-name\" {\n")
	content.WriteString("  #   description = \"Description of the server\"\n")
	content.WriteString("  #   hostname = \"192.168.1.100\"\n")
	content.WriteString("  #   port = 22\n")
	content.WriteString("  #   user = \"admin\"\n")
	content.WriteString("  #   authentication \"publickey\" {\n")
	content.WriteString("  #     public_key_path = \"/path/to/private/key\"\n")
	content.WriteString("  #   }\n")
	content.WriteString("  # }\n\n")

	content.WriteString("  # Add your machine groups here\n")
	content.WriteString("  # group \"group-name\" {\n")
	content.WriteString("  #   description = \"Description of the group\"\n")
	content.WriteString("  #   machines = [\"machine1\", \"machine2\"]\n")
	content.WriteString("  #   user = \"admin\"\n")
	content.WriteString("  #   port = 22\n")
	content.WriteString("  #   authentication \"publickey\" {\n")
	content.WriteString("  #     public_key_path = \"/path/to/private/key\"\n")
	content.WriteString("  #   }\n")
	content.WriteString("  # }\n")

	content.WriteString("}\n")

	return content.String()
}

// GenerateActionsConfigFromStructs generates actions.hcl content from Go structs
func GenerateActionsConfigFromStructs() string {
	var content strings.Builder

	content.WriteString("# Spooky Actions Configuration\n")
	content.WriteString("# Generated from Go struct schemas\n\n")

	content.WriteString("actions {\n")
	content.WriteString("  # Add your action configurations here\n")
	content.WriteString("  # action \"action-name\" {\n")
	content.WriteString("  #   description = \"Description of the action\"\n")
	content.WriteString("  #   type = \"command\"  # command, script, template_deploy, file_sync, service_control\n")
	content.WriteString("  #   command = \"df -h\"  # for command type\n")
	content.WriteString("  #   script = \"scripts/deploy.sh\"  # for script type\n")
	content.WriteString("  #   targets = [\"machine-group\"]\n")
	content.WriteString("  #   timeout = 60\n")
	content.WriteString("  #   sudo = false\n")
	content.WriteString("  # }\n")

	content.WriteString("}\n")

	return content.String()
}

// GenerateVariablesConfigFromStructs generates variables.hcl content from Go structs
func GenerateVariablesConfigFromStructs() string {
	var content strings.Builder

	content.WriteString("# Spooky Variables Configuration\n")
	content.WriteString("# Generated from Go struct schemas\n\n")

	content.WriteString("variables {\n")
	content.WriteString("  # Add your variable definitions here\n")
	content.WriteString("  # variable \"variable-name\" {\n")
	content.WriteString("  #   type = \"string\"  # string, number, boolean\n")
	content.WriteString("  #   value = \"default-value\"\n")
	content.WriteString("  #   description = \"Description of the variable\"\n")
	content.WriteString("  #   sensitive = false  # Set to true for sensitive data\n")
	content.WriteString("  # }\n")

	content.WriteString("}\n")

	return content.String()
}

// ============================================================================
// VALIDATION FUNCTIONS FOR NEW SCHEMAS
// ============================================================================

// ValidateProjectDirectory validates project directory structure against supported schema versions
func (sv *StructValidator) ValidateProjectDirectory(content map[string]interface{}) *ValidationResult {
	result := &ValidationResult{
		IsValid:    true,
		Errors:     []ValidationError{},
		Warnings:   []ValidationWarning{},
		SchemaName: "project-directory",
	}

	// Validate against all supported versions
	for _, version := range sv.supportedVersions {
		switch version {
		case "1":
			sv.validateProjectDirectoryV1(content, result)
		}
	}

	// Update validity based on errors
	if len(result.Errors) > 0 {
		result.IsValid = false
	}

	return result
}

// validateProjectDirectoryV1 validates project directory structure against V1 schema
func (sv *StructValidator) validateProjectDirectoryV1(content map[string]interface{}, result *ValidationResult) {
	// Check for required project_directory block
	projectDirBlock, exists := content["project_directory"]
	if !exists {
		result.Errors = append(result.Errors, ValidationError{
			Field:    "project_directory",
			Message:  "missing required project_directory block",
			Severity: "error",
		})
		return
	}

	// Validate project_directory block structure
	if projectDirMap, ok := projectDirBlock.(map[string]interface{}); ok {
		sv.validateProjectDirectoryBlockV1(projectDirMap, result)
	} else {
		result.Errors = append(result.Errors, ValidationError{
			Field:    "project_directory",
			Message:  "project_directory block must be a configuration block",
			Severity: "error",
		})
	}
}

// validateProjectDirectoryBlockV1 validates the project_directory block contents against V1 schema
func (sv *StructValidator) validateProjectDirectoryBlockV1(projectDir map[string]interface{}, result *ValidationResult) {
	// Required fields
	if name, exists := projectDir["name"]; !exists {
		result.Errors = append(result.Errors, ValidationError{
			Field:    "project_directory.name",
			Message:  "missing required field: name",
			Severity: "error",
		})
	} else {
		sv.validateProjectDirectoryNameV1(name, result)
	}

	// Validate files if present
	if files, exists := projectDir["file"]; exists {
		if filesArray, ok := files.([]interface{}); ok {
			for i, file := range filesArray {
				if fileMap, ok := file.(map[string]interface{}); ok {
					sv.validateProjectDirectoryFileV1(fileMap, fmt.Sprintf("project_directory.file[%d]", i), result)
				}
			}
		}
	}

	// Validate directories if present
	if dirs, exists := projectDir["directory"]; exists {
		if dirsArray, ok := dirs.([]interface{}); ok {
			for i, dir := range dirsArray {
				if dirMap, ok := dir.(map[string]interface{}); ok {
					sv.validateProjectDirectoryDirV1(dirMap, fmt.Sprintf("project_directory.directory[%d]", i), result)
				}
			}
		}
	}
}

// validateProjectDirectoryNameV1 validates project directory name against V1 rules
func (sv *StructValidator) validateProjectDirectoryNameV1(name interface{}, result *ValidationResult) {
	nameStr, ok := name.(string)
	if !ok {
		result.Errors = append(result.Errors, ValidationError{
			Field:    "project_directory.name",
			Message:  "name must be a string",
			Severity: "error",
		})
		return
	}

	// Pattern validation: ^[a-zA-Z0-9_.-]+$
	pattern := regexp.MustCompile(`^[a-zA-Z0-9_.-]+$`)
	if !pattern.MatchString(nameStr) {
		result.Errors = append(result.Errors, ValidationError{
			Field:    "project_directory.name",
			Value:    nameStr,
			Message:  "name must contain only alphanumeric characters, dots, underscores, and hyphens",
			Severity: "error",
		})
	}
}

// validateProjectDirectoryFileV1 validates a project directory file requirement against V1 rules
func (sv *StructValidator) validateProjectDirectoryFileV1(file map[string]interface{}, path string, result *ValidationResult) {
	// Required fields
	if name, exists := file["name"]; !exists {
		result.Errors = append(result.Errors, ValidationError{
			Field:    fmt.Sprintf("%s.name", path),
			Message:  "missing required field: name",
			Severity: "error",
		})
	} else {
		// Validate name is a string
		if _, ok := name.(string); !ok {
			result.Errors = append(result.Errors, ValidationError{
				Field:    fmt.Sprintf("%s.name", path),
				Message:  "name must be a string",
				Severity: "error",
			})
		}
	}

	if fileType, exists := file["type"]; !exists {
		result.Errors = append(result.Errors, ValidationError{
			Field:    fmt.Sprintf("%s.type", path),
			Message:  "missing required field: type",
			Severity: "error",
		})
	} else {
		// Validate type is "file"
		if typeStr, ok := fileType.(string); ok {
			if typeStr != "file" {
				result.Errors = append(result.Errors, ValidationError{
					Field:    fmt.Sprintf("%s.type", path),
					Value:    typeStr,
					Message:  "type must be 'file'",
					Severity: "error",
				})
			}
		} else {
			result.Errors = append(result.Errors, ValidationError{
				Field:    fmt.Sprintf("%s.type", path),
				Message:  "type must be a string",
				Severity: "error",
			})
		}
	}
}

// validateProjectDirectoryDirV1 validates a project directory directory requirement against V1 rules
func (sv *StructValidator) validateProjectDirectoryDirV1(dir map[string]interface{}, path string, result *ValidationResult) {
	// Required fields
	if name, exists := dir["name"]; !exists {
		result.Errors = append(result.Errors, ValidationError{
			Field:    fmt.Sprintf("%s.name", path),
			Message:  "missing required field: name",
			Severity: "error",
		})
	} else {
		// Validate name is a string
		if _, ok := name.(string); !ok {
			result.Errors = append(result.Errors, ValidationError{
				Field:    fmt.Sprintf("%s.name", path),
				Message:  "name must be a string",
				Severity: "error",
			})
		}
	}

	if dirType, exists := dir["type"]; !exists {
		result.Errors = append(result.Errors, ValidationError{
			Field:    fmt.Sprintf("%s.name", path),
			Message:  "missing required field: type",
			Severity: "error",
		})
	} else {
		// Validate type is "directory"
		if typeStr, ok := dirType.(string); ok {
			if typeStr != "directory" {
				result.Errors = append(result.Errors, ValidationError{
					Field:    fmt.Sprintf("%s.type", path),
					Value:    typeStr,
					Message:  "type must be 'directory'",
					Severity: "error",
				})
			}
		} else {
			result.Errors = append(result.Errors, ValidationError{
				Field:    fmt.Sprintf("%s.type", path),
				Message:  "type must be 'directory'",
				Severity: "error",
			})
		}
	}
}

// ValidateSpooky validates global spooky configuration against supported schema versions
func (sv *StructValidator) ValidateSpooky(content map[string]interface{}) *ValidationResult {
	result := &ValidationResult{
		IsValid:    true,
		Errors:     []ValidationError{},
		Warnings:   []ValidationWarning{},
		SchemaName: "spooky",
	}

	// Validate against all supported versions
	for _, version := range sv.supportedVersions {
		switch version {
		case "1":
			sv.validateSpookyV1(content, result)
		}
	}

	// Update validity based on errors
	if len(result.Errors) > 0 {
		result.IsValid = false
	}

	return result
}

// validateSpookyV1 validates global spooky configuration against V1 schema
func (sv *StructValidator) validateSpookyV1(content map[string]interface{}, result *ValidationResult) {
	// Check for required spooky block
	spookyBlock, exists := content["spooky"]
	if !exists {
		result.Errors = append(result.Errors, ValidationError{
			Field:    "spooky",
			Message:  "missing required spooky block",
			Severity: "error",
		})
		return
	}

	// Validate spooky block structure
	if spookyMap, ok := spookyBlock.(map[string]interface{}); ok {
		sv.validateSpookyBlockV1(spookyMap, result)
	} else {
		result.Errors = append(result.Errors, ValidationError{
			Field:    "spooky",
			Message:  "spooky block must be a configuration block",
			Severity: "error",
		})
	}
}

// validateSpookyBlockV1 validates the spooky block contents against V1 schema
func (sv *StructValidator) validateSpookyBlockV1(spooky map[string]interface{}, result *ValidationResult) {
	// Validate SSH configuration if present
	if ssh, exists := spooky["ssh"]; exists {
		if sshMap, ok := ssh.(map[string]interface{}); ok {
			sv.validateSpookySSHV1(sshMap, result)
		}
	}

	// Validate security configuration if present
	if security, exists := spooky["security"]; exists {
		if securityMap, ok := security.(map[string]interface{}); ok {
			sv.validateSpookySecurityV1(securityMap, result)
		}
	}

	// Validate age configuration if present
	if age, exists := spooky["age"]; exists {
		if ageMap, ok := age.(map[string]interface{}); ok {
			sv.validateSpookyAgeV1(ageMap, result)
		}
	}

	// Validate logging configuration if present
	if logging, exists := spooky["logging"]; exists {
		if loggingMap, ok := logging.(map[string]interface{}); ok {
			sv.validateSpookyLoggingV1(loggingMap, result)
		}
	}
}

// validateSpookySSHV1 validates SSH configuration in spooky config
func (sv *StructValidator) validateSpookySSHV1(ssh map[string]interface{}, result *ValidationResult) {
	// Validate timeout if present
	if timeout, exists := ssh["timeout"]; exists {
		if timeoutVal, ok := timeout.(float64); ok {
			if timeoutVal < 1 || timeoutVal > 300 {
				result.Errors = append(result.Errors, ValidationError{
					Field:    "spooky.ssh.timeout",
					Value:    timeoutVal,
					Message:  "SSH timeout must be between 1 and 300 seconds",
					Severity: "error",
				})
			}
		}
	}

	// Validate keepalive_interval if present
	if keepalive, exists := ssh["keepalive_interval"]; exists {
		if keepaliveVal, ok := keepalive.(float64); ok {
			if keepaliveVal < 1 || keepaliveVal > 300 {
				result.Errors = append(result.Errors, ValidationError{
					Field:    "spooky.ssh.keepalive_interval",
					Value:    keepaliveVal,
					Message:  "SSH keepalive interval must be between 1 and 300 seconds",
					Severity: "error",
				})
			}
		}
	}
}

// validateSpookySecurityV1 validates security configuration in spooky config
func (sv *StructValidator) validateSpookySecurityV1(security map[string]interface{}, result *ValidationResult) {
	// Validate allow_unsafe_commands if present
	if allowUnsafe, exists := security["allow_unsafe_commands"]; exists {
		if _, ok := allowUnsafe.(bool); !ok {
			result.Errors = append(result.Errors, ValidationError{
				Field:    "spooky.security.allow_unsafe_commands",
				Message:  "allow_unsafe_commands must be a boolean",
				Severity: "error",
			})
		}
	}
}

// validateSpookyAgeV1 validates age encryption configuration in spooky config
func (sv *StructValidator) validateSpookyAgeV1(age map[string]interface{}, result *ValidationResult) {
	// Validate identities if present
	if identities, exists := age["identities"]; exists {
		if _, ok := identities.(string); !ok {
			result.Errors = append(result.Errors, ValidationError{
				Field:    "spooky.age.identities",
				Message:  "identities must be a string",
				Severity: "error",
			})
		}
	}
}

// validateSpookyLoggingV1 validates logging configuration in spooky config
func (sv *StructValidator) validateSpookyLoggingV1(logging map[string]interface{}, result *ValidationResult) {
	// Validate level if present
	if level, exists := logging["level"]; exists {
		if levelStr, ok := level.(string); ok {
			validLevels := []string{"debug", "info", "warn", "error", "fatal"}
			valid := false
			for _, validLevel := range validLevels {
				if levelStr == validLevel {
					valid = true
					break
				}
			}
			if !valid {
				result.Errors = append(result.Errors, ValidationError{
					Field:    "spooky.logging.level",
					Value:    levelStr,
					Message:  fmt.Sprintf("level must be one of: %s", strings.Join(validLevels, ", ")),
					Severity: "error",
				})
			}
		}
	}
}

// ValidateFacts validates facts configuration against supported schema versions
func (sv *StructValidator) ValidateFacts(content map[string]interface{}) *ValidationResult {
	result := &ValidationResult{
		IsValid:    true,
		Errors:     []ValidationError{},
		Warnings:   []ValidationWarning{},
		SchemaName: "facts",
	}

	// Validate against all supported versions
	for _, version := range sv.supportedVersions {
		switch version {
		case "1":
			sv.validateFactsV1(content, result)
		}
	}

	// Update validity based on errors
	if len(result.Errors) > 0 {
		result.IsValid = false
	}

	return result
}

// validateFactsV1 validates facts configuration against V1 schema
func (sv *StructValidator) validateFactsV1(content map[string]interface{}, result *ValidationResult) {
	// Check for required facts block
	factsBlock, exists := content["facts"]
	if !exists {
		result.Errors = append(result.Errors, ValidationError{
			Field:    "facts",
			Message:  "missing required facts block",
			Severity: "error",
		})
		return
	}

	// Validate facts block structure
	if factsMap, ok := factsBlock.(map[string]interface{}); ok {
		sv.validateFactsBlockV1(factsMap, result)
	} else {
		result.Errors = append(result.Errors, ValidationError{
			Field:    "facts",
			Message:  "facts block must be a configuration block",
			Severity: "error",
		})
	}
}

// validateFactsBlockV1 validates the facts block contents against V1 schema
func (sv *StructValidator) validateFactsBlockV1(facts map[string]interface{}, result *ValidationResult) {
	// Required fields
	if name, exists := facts["name"]; !exists {
		result.Errors = append(result.Errors, ValidationError{
			Field:    "facts.name",
			Message:  "missing required field: name",
			Severity: "error",
		})
	} else {
		sv.validateFactNameV1(name, result)
	}

	if value, exists := facts["value"]; !exists {
		result.Errors = append(result.Errors, ValidationError{
			Field:    "facts.value",
			Message:  "missing required field: value",
			Severity: "error",
		})
	} else {
		// Value can be any type, just ensure it's not nil
		if value == nil {
			result.Errors = append(result.Errors, ValidationError{
				Field:    "facts.value",
				Message:  "value cannot be null",
				Severity: "error",
			})
		}
	}

	if factType, exists := facts["type"]; !exists {
		result.Errors = append(result.Errors, ValidationError{
			Field:    "facts.type",
			Message:  "missing required field: type",
			Severity: "error",
		})
	} else {
		sv.validateFactTypeV1(factType, result)
	}
}

// validateFactNameV1 validates fact name against V1 rules
func (sv *StructValidator) validateFactNameV1(name interface{}, result *ValidationResult) {
	nameStr, ok := name.(string)
	if !ok {
		result.Errors = append(result.Errors, ValidationError{
			Field:    "facts.name",
			Message:  "name must be a string",
			Severity: "error",
		})
		return
	}

	// Pattern validation: ^[a-zA-Z0-9_.-]+$
	pattern := regexp.MustCompile(`^[a-zA-Z0-9_.-]+$`)
	if !pattern.MatchString(nameStr) {
		result.Errors = append(result.Errors, ValidationError{
			Field:    "facts.name",
			Value:    nameStr,
			Message:  "name must contain only alphanumeric characters, dots, underscores, and hyphens",
			Severity: "error",
		})
	}

	// Length validation: 1-128 characters
	if len(nameStr) < 1 || len(nameStr) > 128 {
		result.Errors = append(result.Errors, ValidationError{
			Field:    "facts.name",
			Value:    nameStr,
			Message:  "name must be between 1 and 128 characters",
			Severity: "error",
		})
	}
}

// validateFactTypeV1 validates fact type against V1 rules
func (sv *StructValidator) validateFactTypeV1(factType interface{}, result *ValidationResult) {
	typeStr, ok := factType.(string)
	if !ok {
		result.Errors = append(result.Errors, ValidationError{
			Field:    "facts.type",
			Message:  "type must be a string",
			Severity: "error",
		})
		return
	}

	// Enum validation
	validTypes := []string{"string", "number", "boolean", "object", "array", "encrypted"}
	valid := false
	for _, validType := range validTypes {
		if typeStr == validType {
			valid = true
			break
		}
	}

	if !valid {
		result.Errors = append(result.Errors, ValidationError{
			Field:    "facts.type",
			Value:    typeStr,
			Message:  fmt.Sprintf("type must be one of: %s", strings.Join(validTypes, ", ")),
			Severity: "error",
		})
	}
}

// ValidateLogging validates logging configuration against supported schema versions
func (sv *StructValidator) ValidateLogging(content map[string]interface{}) *ValidationResult {
	result := &ValidationResult{
		IsValid:    true,
		Errors:     []ValidationError{},
		Warnings:   []ValidationWarning{},
		SchemaName: "logging",
	}

	// Validate against all supported versions
	for _, version := range sv.supportedVersions {
		switch version {
		case "1":
			sv.validateLoggingV1(content, result)
		}
	}

	// Update validity based on errors
	if len(result.Errors) > 0 {
		result.IsValid = false
	}

	return result
}

// validateLoggingV1 validates logging configuration against V1 schema
func (sv *StructValidator) validateLoggingV1(content map[string]interface{}, result *ValidationResult) {
	// Check for required logging block
	loggingBlock, exists := content["logging"]
	if !exists {
		result.Errors = append(result.Errors, ValidationError{
			Field:    "logging",
			Message:  "missing required logging block",
			Severity: "error",
		})
		return
	}

	// Validate logging block structure
	if loggingMap, ok := loggingBlock.(map[string]interface{}); ok {
		sv.validateLoggingBlockV1(loggingMap, result)
	} else {
		result.Errors = append(result.Errors, ValidationError{
			Field:    "logging",
			Message:  "logging block must be a configuration block",
			Severity: "error",
		})
	}
}

// validateLoggingBlockV1 validates the logging block contents against V1 schema
func (sv *StructValidator) validateLoggingBlockV1(logging map[string]interface{}, result *ValidationResult) {
	// Validate level if present
	if level, exists := logging["level"]; exists {
		sv.validateLoggingLevelV1(level, result)
	}

	// Validate format if present
	if format, exists := logging["format"]; exists {
		sv.validateLoggingFormatV1(format, result)
	}

	// Validate output if present
	if output, exists := logging["output"]; exists {
		sv.validateLoggingOutputV1(output, result)
	}

	// Validate file_path if output is file
	if output, exists := logging["output"]; exists {
		if outputStr, ok := output.(string); ok && outputStr == "file" {
			if filePath, exists := logging["file_path"]; !exists {
				result.Errors = append(result.Errors, ValidationError{
					Field:    "logging.file_path",
					Message:  "file_path is required when output is 'file'",
					Severity: "error",
				})
			} else {
				// Validate file_path is a string
				if _, ok := filePath.(string); !ok {
					result.Errors = append(result.Errors, ValidationError{
						Field:    "logging.file_path",
						Message:  "file_path must be a string",
						Severity: "error",
					})
				}
			}
		}
	}
}

// validateLoggingLevelV1 validates logging level against V1 rules
func (sv *StructValidator) validateLoggingLevelV1(level interface{}, result *ValidationResult) {
	levelStr, ok := level.(string)
	if !ok {
		result.Errors = append(result.Errors, ValidationError{
			Field:    "logging.level",
			Message:  "level must be a string",
			Severity: "error",
		})
		return
	}

	// Enum validation
	validLevels := []string{"debug", "info", "warn", "error", "fatal"}
	valid := false
	for _, validLevel := range validLevels {
		if levelStr == validLevel {
			valid = true
			break
		}
	}

	if !valid {
		result.Errors = append(result.Errors, ValidationError{
			Field:    "logging.level",
			Value:    levelStr,
			Message:  fmt.Sprintf("level must be one of: %s", strings.Join(validLevels, ", ")),
			Severity: "error",
		})
	}
}

// validateLoggingFormatV1 validates logging format against V1 rules
func (sv *StructValidator) validateLoggingFormatV1(format interface{}, result *ValidationResult) {
	formatStr, ok := format.(string)
	if !ok {
		result.Errors = append(result.Errors, ValidationError{
			Field:    "logging.format",
			Message:  "format must be a string",
			Severity: "error",
		})
		return
	}

	// Enum validation
	validFormats := []string{"json", "text", "structured"}
	valid := false
	for _, validFormat := range validFormats {
		if formatStr == validFormat {
			valid = true
			break
		}
	}

	if !valid {
		result.Errors = append(result.Errors, ValidationError{
			Field:    "logging.format",
			Value:    formatStr,
			Message:  fmt.Sprintf("format must be one of: %s", strings.Join(validFormats, ", ")),
			Severity: "error",
		})
	}
}

// validateLoggingOutputV1 validates logging output against V1 rules
func (sv *StructValidator) validateLoggingOutputV1(output interface{}, result *ValidationResult) {
	outputStr, ok := output.(string)
	if !ok {
		result.Errors = append(result.Errors, ValidationError{
			Field:    "logging.output",
			Message:  "output must be a string",
			Severity: "error",
		})
		return
	}

	// Enum validation
	validOutputs := []string{"stdout", "stderr", "file", "null"}
	valid := false
	for _, validOutput := range validOutputs {
		if outputStr == validOutput {
			valid = true
			break
		}
	}

	if !valid {
		result.Errors = append(result.Errors, ValidationError{
			Field:    "logging.output",
			Value:    outputStr,
			Message:  fmt.Sprintf("output must be one of: %s", strings.Join(validOutputs, ", ")),
			Severity: "error",
		})
	}
}

// ============================================================================
// GENERATION FUNCTIONS FOR NEW SCHEMAS
// ============================================================================

// GenerateSpookyConfigFromStructs generates spooky.hcl content from Go structs
func GenerateSpookyConfigFromStructs() string {
	var content strings.Builder

	content.WriteString("# Spooky Global Configuration\n")
	content.WriteString("# Generated from Go struct schemas\n\n")

	content.WriteString("spooky {\n")
	content.WriteString("  # SSH Configuration\n")
	content.WriteString("  ssh {\n")
	content.WriteString("    timeout = 30  # SSH connection timeout in seconds\n")
	content.WriteString("    keepalive_interval = 60  # SSH keepalive interval in seconds\n")
	content.WriteString("  }\n\n")

	content.WriteString("  # Security Configuration\n")
	content.WriteString("  security {\n")
	content.WriteString("    allow_unsafe_commands = false  # Allow potentially unsafe commands\n")
	content.WriteString("  }\n\n")

	content.WriteString("  # Age Encryption Configuration\n")
	content.WriteString("  age {\n")
	content.WriteString("    identities = \"~/.config/spooky/age/identities.txt\"  # Path to age identities\n")
	content.WriteString("  }\n\n")

	content.WriteString("  # Logging Configuration\n")
	content.WriteString("  logging {\n")
	content.WriteString("    level = \"info\"  # debug, info, warn, error, fatal\n")
	content.WriteString("    format = \"structured\"  # json, text, structured\n")
	content.WriteString("    output = \"stdout\"  # stdout, stderr, file, null\n")
	content.WriteString("    # file_path = \"/var/log/spooky.log\"  # Required when output is 'file'\n")
	content.WriteString("  }\n")

	content.WriteString("}\n")

	return content.String()
}

// GenerateFactsConfigFromStructs generates facts.hcl content from Go structs
func GenerateFactsConfigFromStructs() string {
	var content strings.Builder

	content.WriteString("# Spooky Facts Configuration\n")
	content.WriteString("# Generated from Go struct schemas\n\n")

	content.WriteString("facts {\n")
	content.WriteString("  # Example fact definition\n")
	content.WriteString("  # fact \"server_os\" {\n")
	content.WriteString("  #   type = \"string\"  # string, number, boolean, object, array, encrypted\n")
	content.WriteString("  #   value = \"ubuntu\"\n")
	content.WriteString("  #   description = \"Operating system of the server\"\n")
	content.WriteString("  #   sensitive = false\n")
	content.WriteString("  # }\n\n")

	content.WriteString("  # Add your fact definitions here\n")
	content.WriteString("  # Each fact should have a name, type, and value\n")
	content.WriteString("  # Use 'encrypted' type for sensitive information\n")

	content.WriteString("}\n")

	return content.String()
}

// GenerateLoggingConfigFromStructs generates logging.hcl content from Go structs
func GenerateLoggingConfigFromStructs() string {
	var content strings.Builder

	content.WriteString("# Spooky Logging Configuration\n")
	content.WriteString("# Generated from Go struct schemas\n\n")

	content.WriteString("logging {\n")
	content.WriteString("  level = \"info\"  # debug, info, warn, error, fatal\n")
	content.WriteString("  format = \"structured\"  # json, text, structured\n")
	content.WriteString("  output = \"stdout\"  # stdout, stderr, file, null\n")
	content.WriteString("  # file_path = \"/var/log/spooky.log\"  # Required when output is 'file'\n")
	content.WriteString("  # max_file_size = \"100MB\"  # Maximum log file size\n")
	content.WriteString("  # max_files = 5  # Maximum number of log files to keep\n")
	content.WriteString("  # compress_old_logs = true  # Compress old log files\n")

	content.WriteString("}\n")

	return content.String()
}
