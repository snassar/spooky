package schemas

import (
	"fmt"
	"log/slog"
	"regexp"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/zclconf/go-cty/cty"
)

// ValidationRange defines a range for validation rules
type ValidationRange struct {
	Min interface{} `json:"min,omitempty"`
	Max interface{} `json:"max,omitempty"`
}

// EnhancedValidationRule defines a business rule for enhanced validation
type EnhancedValidationRule struct {
	Name      string           `json:"name"`
	Type      string           `json:"type"`
	Message   string           `json:"message"`
	Condition string           `json:"condition,omitempty"`
	Pattern   string           `json:"pattern,omitempty"`
	Enum      string           `json:"enum,omitempty"`
	Range     *ValidationRange `json:"range,omitempty"`
	Severity  string           `json:"severity"`
}

// Validator provides comprehensive validation combining schema validation and business rules
type Validator struct {
	logger *slog.Logger
	// Support for multiple schema versions
	supportedVersions []string
	// Enhanced validation rules
	enhancedRules map[string][]EnhancedValidationRule
}

// NewValidator creates a new unified validator with both schema and enhanced validation
func NewValidator() *Validator {
	validator := &Validator{
		logger:            slog.Default(),
		supportedVersions: SupportedVersions,
		enhancedRules:     make(map[string][]EnhancedValidationRule),
	}

	// Load enhanced validation rules
	if err := validator.loadEnhancedRules(); err != nil {
		// Log error but continue without enhanced validation
		validator.logger.Warn("Failed to load enhanced validation rules", "error", err)
	}

	return validator
}

// ValidateHCLContent validates HCL content against both schema and enhanced rules
func (v *Validator) ValidateHCLContent(schemaName, content string) (*ValidationResult, error) {
	// First parse HCL into structured data
	parsedData, err := v.ParseHCLContent(content)
	if err != nil {
		return &ValidationResult{
			IsValid:    false,
			Errors:     []ValidationError{{Message: fmt.Sprintf("HCL parsing failed: %v", err), Severity: "error"}},
			Warnings:   []ValidationWarning{},
			SchemaName: schemaName,
		}, nil
	}

	// Then validate against both schema and enhanced rules
	return v.ValidateData(schemaName, parsedData), nil
}

// ValidateData validates parsed data against both schema and enhanced rules
func (v *Validator) ValidateData(schemaName string, data map[string]interface{}) *ValidationResult {
	result := &ValidationResult{
		IsValid:    true,
		Errors:     []ValidationError{},
		Warnings:   []ValidationWarning{},
		SchemaName: schemaName,
	}

	// Step 1: Basic schema validation
	v.validateSchema(schemaName, data, result)

	// Step 2: Enhanced business rule validation
	v.validateEnhancedRules(schemaName, data, result)

	// Update validity based on errors
	if len(result.Errors) > 0 {
		result.IsValid = false
	}

	return result
}

// validateSchema performs basic schema validation
func (v *Validator) validateSchema(schemaName string, data map[string]interface{}, result *ValidationResult) {
	// Validate against all supported versions
	for _, version := range v.supportedVersions {
		switch version {
		case "1":
			switch schemaName {
			case "project":
				v.validateProjectV1(data, result)
			case "machines":
				v.validateMachinesV1(data, result)
			case "actions":
				v.validateActionsV1(data, result)
			case "variables":
				v.validateVariablesV1(data, result)
			case "logging":
				v.validateLoggingV1(data, result)
			case "facts":
				v.validateFactsV1(data, result)
			case "spooky":
				v.validateSpookyV1(data, result)
			case "project-directory":
				v.validateProjectDirectoryV1(data, result)
			}
		}
	}
}

// validateEnhancedRules performs enhanced business rule validation
func (v *Validator) validateEnhancedRules(schemaName string, data map[string]interface{}, result *ValidationResult) {
	rules, exists := v.enhancedRules[schemaName]
	if !exists {
		return
	}

	switch schemaName {
	case "actions":
		v.validateActionsEnhanced(data, rules, result)
	case "machines":
		v.validateMachinesEnhanced(data, rules, result)
	case "variables":
		v.validateVariablesEnhanced(data, rules, result)
	}
}

// ParseHCLContent parses HCL content into a structured map
func (v *Validator) ParseHCLContent(content string) (map[string]interface{}, error) {
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
		blockData := v.parseBlock(block)

		// Handle project blocks with labels (project "name" {})
		if block.Type == "project" && len(block.Labels) > 0 {
			// Add the name as a label to the block data for validation
			blockData["name"] = block.Labels[0]
		}

		data[block.Type] = blockData
	}

	return data, nil
}

// parseBlock recursively parses HCL blocks into structured data
func (v *Validator) parseBlock(block *hclsyntax.Block) map[string]interface{} {
	result := make(map[string]interface{})

	// Parse attributes
	for name, attr := range block.Body.Attributes {
		value := v.parseAttributeValue(attr.Expr)
		result[name] = value
	}

	// Parse nested blocks
	for _, nestedBlock := range block.Body.Blocks {
		nestedData := v.parseBlock(nestedBlock)

		// Handle blocks with labels (like "action "name" {")
		if len(nestedBlock.Labels) > 0 {
			label := nestedBlock.Labels[0]
			// Add the name as a field to the nested data for validation
			nestedData["name"] = label

			// For action blocks, store under "action" key, not the label
			if nestedBlock.Type == "action" {
				if existing, exists := result["action"]; exists {
					// Convert to array if multiple blocks with same label
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
					// Convert to array if multiple blocks with same label
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
func (v *Validator) parseAttributeValue(expr hclsyntax.Expression) interface{} {
	switch expr := expr.(type) {
	case *hclsyntax.LiteralValueExpr:
		return v.ctyValueToInterface(expr.Val)
	case *hclsyntax.TemplateExpr:
		// Handle template expressions - try to extract literal values
		if len(expr.Parts) == 1 {
			// Single part - check if it's a literal value
			if lit, ok := expr.Parts[0].(*hclsyntax.LiteralValueExpr); ok {
				return v.ctyValueToInterface(lit.Val)
			}
		}
		// For complex templates, return as string for now
		return fmt.Sprintf("%v", expr)
	case *hclsyntax.TupleConsExpr:
		// Handle lists/arrays
		var result []interface{}
		for _, elem := range expr.Exprs {
			result = append(result, v.parseAttributeValue(elem))
		}
		return result
	case *hclsyntax.ObjectConsExpr:
		// Handle maps/objects
		result := make(map[string]interface{})
		for _, item := range expr.Items {
			key := v.parseAttributeValue(item.KeyExpr)
			value := v.parseAttributeValue(item.ValueExpr)
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

// ctyValueToInterface converts cty.Value to Go interface{}
func (v *Validator) ctyValueToInterface(val cty.Value) interface{} {
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
				result = append(result, v.ctyValueToInterface(elemVal))
			}
			return result
		}
		return val.GoString()
	}
}

// loadEnhancedRules loads the built-in enhanced validation rules
func (v *Validator) loadEnhancedRules() error {
	// Actions validation rules
	v.enhancedRules["actions"] = []EnhancedValidationRule{
		{
			Name:      "machine_targeting_required",
			Type:      "cross_field",
			Message:   "Action must target machines via 'machines' or 'tags' field",
			Condition: "machines == null && tags == null",
			Severity:  "error",
		},
		{
			Name:      "command_type_validation",
			Type:      "type_validation",
			Message:   "Command actions must have a 'command' field",
			Condition: "type == 'command' && command == null",
			Severity:  "error",
		},

		{
			Name:      "template_deploy_validation",
			Type:      "type_validation",
			Message:   "Template deploy actions must have 'source' and 'destination' fields",
			Condition: "type == 'template_deploy' && (source == null || destination == null)",
			Severity:  "error",
		},
		{
			Name:      "file_sync_validation",
			Type:      "type_validation",
			Message:   "File sync actions must have 'source' and 'destination' fields",
			Condition: "type == 'file_sync' && (source == null || destination == null)",
			Severity:  "error",
		},
		{
			Name:      "service_control_validation",
			Type:      "type_validation",
			Message:   "Service control actions must have 'service' and 'action' fields",
			Condition: "type == 'service_control' && (service == null || action == null)",
			Severity:  "error",
		},
		{
			Name:     "shell_operators_forbidden",
			Type:     "pattern",
			Message:  "Shell operators (;&|$) are not allowed in commands for security reasons",
			Pattern:  `[;&|\x24]`,
			Severity: "error",
		},
		{
			Name:      "critical_and_allow_failure_conflict",
			Type:      "cross_field",
			Message:   "Action cannot be both 'critical' and 'allow_failure'",
			Condition: "critical == true && allow_failure == true",
			Severity:  "error",
		},
	}

	// Machines validation rules
	v.enhancedRules["machines"] = []EnhancedValidationRule{
		{
			Name:      "user_and_sudo_conflict",
			Type:      "cross_field",
			Message:   "Machine cannot have both 'user' and 'sudo' fields",
			Condition: "user != null && sudo == true",
			Severity:  "error",
		},
		{
			Name:      "ssh_key_method_validation",
			Type:      "cross_field",
			Message:   "SSH key authentication requires 'public_key_path' field",
			Condition: "authentication.method == 'publickey' && authentication.public_key_path == null",
			Severity:  "error",
		},
		{
			Name:      "password_method_validation",
			Type:      "cross_field",
			Message:   "Password authentication requires 'password' field",
			Condition: "authentication.method == 'password' && authentication.password == null",
			Severity:  "error",
		},
		{
			Name:      "certificate_method_validation",
			Type:      "cross_field",
			Message:   "Certificate authentication requires 'private_key_path' and 'certificate_path' fields",
			Condition: "authentication.method == 'certificate' && (authentication.private_key_path == null || authentication.certificate_path == null)",
			Severity:  "error",
		},
	}

	// Variables validation rules
	v.enhancedRules["variables"] = []EnhancedValidationRule{
		{
			Name:      "encrypted_value_format",
			Type:      "pattern",
			Message:   "Encrypted values must be in age format (age1...)",
			Pattern:   `^age1[a-z0-9]+$`,
			Condition: "encrypted == true",
			Severity:  "error",
		},
		{
			Name:    "tags_count_limit",
			Type:    "range",
			Message: "Tags count must be between 0 and 10",
			Range: &ValidationRange{
				Min: 0,
				Max: 10,
			},
			Severity: "warning",
		},
	}

	return nil
}

// evaluateRule evaluates a single enhanced validation rule
func (v *Validator) evaluateRule(rule EnhancedValidationRule, data map[string]interface{}) bool {
	switch rule.Type {
	case "cross_field":
		return v.evaluateCrossFieldRule(rule, data)
	case "pattern":
		return v.evaluatePatternRule(rule, data)
	case "enum":
		return v.evaluateEnumRule(rule, data)
	case "range":
		return v.evaluateRangeRule(rule, data)
	case "type_validation":
		return v.evaluateTypeValidationRule(rule, data)
	default:
		return true // Unknown rule type, pass validation
	}
}

// evaluateCrossFieldRule evaluates cross-field validation rules
func (v *Validator) evaluateCrossFieldRule(rule EnhancedValidationRule, data map[string]interface{}) bool {
	// Use the condition field to determine which rule to apply
	condition := rule.Condition

	// Simple condition evaluation for common patterns
	switch condition {
	case "(machines == null || machines.isEmpty()) && (tags == null || tags.isEmpty())":
		return v.checkMachineTargetingRequired(data)
	case "machines == null && tags == null":
		// Only apply this rule if both machines and tags are actually null
		if machines, exists := data["machines"]; exists && machines != nil {
			return false // No violation - machines exist
		}
		if tags, exists := data["tags"]; exists && tags != nil {
			return false // No violation - tags exist
		}
		return true // Violation - neither machines nor tags exist
	case "critical == true && allow_failure == true":
		// Only apply this rule if both critical and allow_failure are actually true
		critical, criticalExists := data["critical"]
		allowFailure, allowFailureExists := data["allow_failure"]
		if criticalExists && critical == true && allowFailureExists && allowFailure == true {
			return true // Violation - both are true
		}
		return false // No violation
	case "user != null && sudo == true":
		return v.checkUserAndSudoConflict(data)
	case "authentication.method == 'publickey' && authentication.public_key_path == null":
		// Only apply this rule if the method is actually 'publickey'
		if auth, exists := data["authentication"]; exists {
			if authMap, ok := auth.(map[string]interface{}); ok {
				if method, exists := authMap["method"]; exists && method == "publickey" {
					return v.checkSSHKeyMethodValidation(data)
				}
			}
		}
		return false
	case "authentication.method == 'password' && authentication.password == null":
		// Only apply this rule if the method is actually 'password'
		if auth, exists := data["authentication"]; exists {
			if authMap, ok := auth.(map[string]interface{}); ok {
				if method, exists := authMap["method"]; exists && method == "password" {
					return v.checkPasswordMethodValidation(data)
				}
			}
		}
		return false
	case "authentication.method == 'certificate' && (authentication.private_key_path == null || authentication.certificate_path == null)":
		// Only apply this rule if the method is actually 'certificate'
		if auth, exists := data["authentication"]; exists {
			if authMap, ok := auth.(map[string]interface{}); ok {
				if method, exists := authMap["method"]; exists && method == "certificate" {
					return v.checkCertificateMethodValidation(data)
				}
			}
		}
		return false
	case "encrypted == true":
		// Only apply this rule if the variable is actually encrypted
		if encrypted, exists := data["encrypted"]; exists && encrypted == true {
			return true // Apply the rule
		}
		return false // Don't apply the rule
	default:
		// For unknown conditions, return false (no violation)
		return false
	}
}

// evaluatePatternRule evaluates pattern-based validation rules
func (v *Validator) evaluatePatternRule(rule EnhancedValidationRule, data map[string]interface{}) bool {
	pattern := rule.Pattern
	if pattern == "" {
		return false // No pattern, no violation
	}
	regex, err := regexp.Compile(pattern)
	if err != nil {
		return false // Invalid regex, no violation
	}

	switch rule.Name {
	case "shell_operators_forbidden":
		if command, exists := data["command"]; exists {
			if commandStr, ok := command.(string); ok {
				return regex.MatchString(commandStr) // Return true if pattern matches (violation)
			}
		}
		return false // No violation
	case "encrypted_value_format":
		if encrypted, exists := data["encrypted"]; exists && encrypted == true {
			if value, exists := data["value"]; exists {
				if valueStr, ok := value.(string); ok {
					return !regex.MatchString(valueStr) // Return true if pattern doesn't match (violation)
				}
			}
		}
		return false // No violation
	default:
		return false // No violation
	}
}

// evaluateEnumRule evaluates enum-based validation rules
func (v *Validator) evaluateEnumRule(rule EnhancedValidationRule, data map[string]interface{}) bool {
	enum := rule.Enum
	if len(enum) == 0 {
		return false // No enum, no violation
	}

	// Get the field value to check
	fieldValue := v.getFieldValue(rule.Name, data)
	if fieldValue == nil {
		return false // No field value, no violation
	}

	// Check if the value is in the enum
	for _, enumValue := range enum {
		if fieldValue == enumValue {
			return false // Value is in enum, no violation
		}
	}

	return true // Value is not in enum, violation
}

// evaluateRangeRule evaluates range-based validation rules
func (v *Validator) evaluateRangeRule(rule EnhancedValidationRule, data map[string]interface{}) bool {
	if rule.Range == nil {
		return false // No range, no violation
	}

	switch rule.Name {
	case "tags_count_limit":
		if tags, exists := data["tags"]; exists {
			if tagsSlice, ok := tags.([]interface{}); ok {
				count := len(tagsSlice)
				min := 0
				max := 10

				if rule.Range.Min != nil {
					if minInt, ok := rule.Range.Min.(int); ok {
						min = minInt
					}
				}
				if rule.Range.Max != nil {
					if maxInt, ok := rule.Range.Max.(int); ok {
						max = maxInt
					}
				}

				return count < min || count > max // Return true if count is out of range (violation)
			}
		}
		return false // No tags, no violation
	default:
		return false // No violation
	}
}

// evaluateTypeValidationRule evaluates type-specific validation rules
func (v *Validator) evaluateTypeValidationRule(rule EnhancedValidationRule, data map[string]interface{}) bool {
	switch rule.Name {
	case "command_type_validation":
		return v.checkCommandTypeValidation(data)

	case "template_deploy_validation":
		return v.checkTemplateDeployValidation(data)
	case "file_sync_validation":
		return v.checkFileSyncValidation(data)
	case "service_control_validation":
		return v.checkServiceControlValidation(data)
	default:
		return true
	}
}

// Helper methods for cross-field validation
func (v *Validator) checkMachineTargetingRequired(data map[string]interface{}) bool {
	// Check if machines or tags field exists and is not empty
	if machines, exists := data["machines"]; exists && machines != nil {
		if machinesSlice, ok := machines.([]interface{}); ok && len(machinesSlice) > 0 {
			return false // No violation - machines are specified
		}
	}
	if tags, exists := data["tags"]; exists && tags != nil {
		if tagsSlice, ok := tags.([]interface{}); ok && len(tagsSlice) > 0 {
			return false // No violation - tags are specified
		}
	}
	return true // Violation - neither machines nor tags are specified
}

func (v *Validator) checkCriticalAndAllowFailureConflict(data map[string]interface{}) bool {
	critical, criticalExists := data["critical"]
	allowFailure, allowFailureExists := data["allow_failure"]

	// Return true if there's a conflict (both critical and allow_failure are true)
	return criticalExists && critical == true && allowFailureExists && allowFailure == true
}

func (v *Validator) checkUserAndSudoConflict(data map[string]interface{}) bool {
	user, userExists := data["user"]
	sudo, sudoExists := data["sudo"]

	// Return true if there's a conflict (user exists and sudo is true)
	return userExists && user != "" && sudoExists && sudo == true
}

func (v *Validator) checkSSHKeyMethodValidation(data map[string]interface{}) bool {
	if auth, exists := data["authentication"]; exists {
		if authMap, ok := auth.(map[string]interface{}); ok {
			if method, exists := authMap["method"]; exists && method == "publickey" {
				if _, hasKey := authMap["public_key_path"]; !hasKey {
					return true // Return true if there's a violation (missing key)
				}
			}
		}
	}
	return false // Return false if there's no violation
}

func (v *Validator) checkPasswordMethodValidation(data map[string]interface{}) bool {
	if auth, exists := data["authentication"]; exists {
		if authMap, ok := auth.(map[string]interface{}); ok {
			if method, exists := authMap["method"]; exists && method == "password" {
				if _, hasPassword := authMap["password"]; !hasPassword {
					return true // Return true if there's a violation (missing password)
				}
			}
		}
	}
	return false // Return false if there's no violation
}

func (v *Validator) checkCertificateMethodValidation(data map[string]interface{}) bool {
	if auth, exists := data["authentication"]; exists {
		if authMap, ok := auth.(map[string]interface{}); ok {
			if method, exists := authMap["method"]; exists && method == "certificate" {
				if _, hasCert := authMap["certificate_path"]; !hasCert {
					return true // Return true if there's a violation (missing certificate)
				}
				if _, hasKey := authMap["private_key_path"]; !hasKey {
					return true // Return true if there's a violation (missing key)
				}
			}
		}
	}
	return false // Return false if there's no violation
}

// Helper methods for type validation
func (v *Validator) checkCommandTypeValidation(data map[string]interface{}) bool {
	if actionType, exists := data["type"]; exists && actionType == "command" {
		if _, hasCommand := data["command"]; !hasCommand {
			return true // Return true if there's a violation (missing command)
		}
	}
	return false // Return false if there's no violation
}

func (v *Validator) checkTemplateDeployValidation(data map[string]interface{}) bool {
	if actionType, exists := data["type"]; exists && actionType == "template_deploy" {
		if _, hasSource := data["source"]; !hasSource {
			return true // Return true if there's a violation (missing source)
		}
		if _, hasDest := data["destination"]; !hasDest {
			return true // Return true if there's a violation (missing destination)
		}
	}
	return false // Return false if there's no violation
}

func (v *Validator) checkFileSyncValidation(data map[string]interface{}) bool {
	if actionType, exists := data["type"]; exists && actionType == "file_sync" {
		if _, hasSource := data["source"]; !hasSource {
			return true // Return true if there's a violation (missing source)
		}
		if _, hasDest := data["destination"]; !hasDest {
			return true // Return true if there's a violation (missing destination)
		}
	}
	return false // Return false if there's no violation
}

func (v *Validator) checkServiceControlValidation(data map[string]interface{}) bool {
	if actionType, exists := data["type"]; exists && actionType == "service_control" {
		if _, hasService := data["service"]; !hasService {
			return true // Return true if there's a violation (missing service)
		}
		if _, hasAction := data["action"]; !hasAction {
			return true // Return true if there's a violation (missing action)
		}
	}
	return false // Return false if there's no violation
}

// getFieldValue extracts field value from data based on rule name
func (v *Validator) getFieldValue(ruleName string, data map[string]interface{}) interface{} {
	switch ruleName {

	case "tags_count_limit_exceeded":
		if tags, exists := data["tags"]; exists {
			if tagsSlice, ok := tags.([]interface{}); ok {
				return len(tagsSlice)
			}
		}
		return 0
	default:
		return data[ruleName]
	}
}

// Schema-specific enhanced validation methods
func (v *Validator) validateActionsEnhanced(data map[string]interface{}, rules []EnhancedValidationRule, result *ValidationResult) {
	// Get the actions block
	actionsBlock, exists := data["actions"]
	if !exists {
		return
	}

	// Get the action list from the actions block
	actionsMap, ok := actionsBlock.(map[string]interface{})
	if !ok {
		return
	}

	actionList, exists := actionsMap["action"]
	if !exists {
		return
	}

	// Handle action array
	var actions []map[string]interface{}
	if actionArray, ok := actionList.([]interface{}); ok {
		for _, action := range actionArray {
			if actionMap, ok := action.(map[string]interface{}); ok {
				actions = append(actions, actionMap)
			}
		}
	} else if actionMap, ok := actionList.(map[string]interface{}); ok {
		// Single action
		actions = []map[string]interface{}{actionMap}
	} else {
		return
	}

	// Validate each action
	for i, action := range actions {
		for _, rule := range rules {
			if v.evaluateRule(rule, action) {
				error := ValidationError{
					Field:    fmt.Sprintf("actions.action[%d]", i),
					Value:    action,
					Message:  rule.Message,
					Severity: rule.Severity,
					Context:  rule.Name,
				}
				result.Errors = append(result.Errors, error)
			}
		}
	}
}

func (v *Validator) validateMachinesEnhanced(data map[string]interface{}, rules []EnhancedValidationRule, result *ValidationResult) {
	// Get the machines block
	machinesBlock, exists := data["machines"]
	if !exists {
		return // No machines block, let basic validation handle this
	}

	machinesMap, ok := machinesBlock.(map[string]interface{})
	if !ok {
		return // Invalid structure, let basic validation handle this
	}

	// Get the machine list
	machineList, exists := machinesMap["machine"]
	if !exists {
		return // No machines, let basic validation handle this
	}

	machines, ok := machineList.([]interface{})
	if !ok {
		return // Invalid structure, let basic validation handle this
	}

	// Validate each machine individually
	for i, machine := range machines {
		machineMap, ok := machine.(map[string]interface{})
		if !ok {
			continue
		}

		// Apply enhanced validation rules to this machine
		for _, rule := range rules {
			if v.evaluateRule(rule, machineMap) {
				if rule.Severity == "error" {
					result.Errors = append(result.Errors, ValidationError{
						Message:  rule.Message,
						Severity: rule.Severity,
						Context:  rule.Name,
						Field:    fmt.Sprintf("machines.machine[%d]", i),
					})
					result.IsValid = false
				} else if rule.Severity == "warning" {
					result.Warnings = append(result.Warnings, ValidationWarning{
						Message: rule.Message,
						Context: rule.Name,
						Field:   fmt.Sprintf("machines.machine[%d]", i),
					})
				}
			}
		}
	}
}

func (v *Validator) validateVariablesEnhanced(data map[string]interface{}, rules []EnhancedValidationRule, result *ValidationResult) {
	// Get the variables block
	variablesBlock, exists := data["variables"]
	if !exists {
		return
	}

	// Get the variable list from the variables block
	variablesMap, ok := variablesBlock.(map[string]interface{})
	if !ok {
		return
	}

	variableList, exists := variablesMap["variable"]
	if !exists {
		return
	}

	// Handle variable array
	var variables []map[string]interface{}
	if variableArray, ok := variableList.([]interface{}); ok {
		for _, variable := range variableArray {
			if variableMap, ok := variable.(map[string]interface{}); ok {
				variables = append(variables, variableMap)
			}
		}
	} else if variableMap, ok := variableList.(map[string]interface{}); ok {
		// Single variable
		variables = []map[string]interface{}{variableMap}
	} else {
		return
	}

	// Validate each variable
	for i, variable := range variables {
		for _, rule := range rules {
			if v.evaluateRule(rule, variable) {
				error := ValidationError{
					Field:    fmt.Sprintf("variables.variable[%d]", i),
					Value:    variable,
					Message:  rule.Message,
					Severity: rule.Severity,
					Context:  rule.Name,
				}
				result.Errors = append(result.Errors, error)
			}
		}
	}
}

// Schema validation methods (V1) - migrated from struct_validator.go
func (v *Validator) validateProjectV1(content map[string]interface{}, result *ValidationResult) {
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
		v.validateProjectBlockV1(projectMap, result)
	} else {
		result.Errors = append(result.Errors, ValidationError{
			Field:    "project",
			Message:  "project block must be a configuration block",
			Severity: "error",
		})
	}
}

// validateProjectBlockV1 validates the project block contents against V1 schema
func (v *Validator) validateProjectBlockV1(project map[string]interface{}, result *ValidationResult) {
	// Required fields - name is now a label, validated during parsing
	if name, exists := project["name"]; exists {
		v.validateProjectNameV1(name, result)
	}

	if description, exists := project["description"]; !exists {
		result.Errors = append(result.Errors, ValidationError{
			Field:    "project.description",
			Message:  "missing required field: description",
			Severity: "error",
		})
	} else {
		v.validateProjectDescriptionV1(description, result)
	}

	// Optional fields with validation
	if runMaxParallel, exists := project["run_max_parallel"]; exists {
		v.validateProjectRunMaxParallelV1(runMaxParallel, result)
	}

	if factsTimeout, exists := project["facts_timeout"]; exists {
		v.validateProjectFactsTimeoutV1(factsTimeout, result)
	}
}

// validateProjectNameV1 validates project name against V1 rules
func (v *Validator) validateProjectNameV1(name interface{}, result *ValidationResult) {
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
func (v *Validator) validateProjectDescriptionV1(description interface{}, result *ValidationResult) {
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
			Message:  "description must be at most 1024 characters",
			Severity: "error",
		})
	}
}

// validateProjectRunMaxParallelV1 validates run_max_parallel against V1 rules
func (v *Validator) validateProjectRunMaxParallelV1(value interface{}, result *ValidationResult) {
	// Type validation
	val, ok := value.(float64)
	if !ok {
		result.Errors = append(result.Errors, ValidationError{
			Field:    "project.run_max_parallel",
			Value:    value,
			Message:  "run_max_parallel must be a number",
			Severity: "error",
		})
		return
	}

	// Range validation: 1-100
	if val < 1 || val > 100 {
		result.Errors = append(result.Errors, ValidationError{
			Field:    "project.run_max_parallel",
			Value:    val,
			Message:  "run_max_parallel must be between 1 and 100",
			Severity: "error",
		})
	}
}

// validateProjectFactsTimeoutV1 validates facts_timeout against V1 rules
func (v *Validator) validateProjectFactsTimeoutV1(value interface{}, result *ValidationResult) {
	// Type validation
	val, ok := value.(float64)
	if !ok {
		result.Errors = append(result.Errors, ValidationError{
			Field:    "project.facts_timeout",
			Value:    value,
			Message:  "facts_timeout must be a number",
			Severity: "error",
		})
		return
	}

	// Range validation: 1-3600
	if val < 1 || val > 3600 {
		result.Errors = append(result.Errors, ValidationError{
			Field:    "project.facts_timeout",
			Value:    val,
			Message:  "facts_timeout must be between 1 and 3600",
			Severity: "error",
		})
	}
}

// validateMachinesV1 validates machines against V1 schema
func (v *Validator) validateMachinesV1(data map[string]interface{}, result *ValidationResult) {
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
		v.validateMachinesBlockV1(machinesMap, result)
	} else {
		result.Errors = append(result.Errors, ValidationError{
			Field:    "machines",
			Message:  "machines block must be a configuration block",
			Severity: "error",
		})
	}
}

// validateMachinesBlockV1 validates the machines block contents against V1 schema
func (v *Validator) validateMachinesBlockV1(machines map[string]interface{}, result *ValidationResult) {
	// Check for at least one machine or group
	hasMachine := false
	hasGroup := false

	// In HCL, machine blocks are identified by their block labels (machine names)
	// So we iterate through all keys in the machines map
	for machineName, machineValue := range machines {
		if machineName == "group" {
			// Handle group blocks
			hasGroup = true
			if groupArray, ok := machineValue.([]interface{}); ok {
				for i, group := range groupArray {
					if groupMap, ok := group.(map[string]interface{}); ok {
						v.validateGroupV1(groupMap, fmt.Sprintf("machines.group[%d]", i), result)
					}
				}
			} else if groupMap, ok := machineValue.(map[string]interface{}); ok {
				v.validateGroupV1(groupMap, "machines.group[0]", result)
			}
		} else {
			// This is a machine block (machine name is the key)
			hasMachine = true
			if machineMap, ok := machineValue.(map[string]interface{}); ok {
				v.validateMachineV1(machineMap, fmt.Sprintf("machines.%s", machineName), result)
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
func (v *Validator) validateMachineV1(machine map[string]interface{}, path string, result *ValidationResult) {
	// Required fields
	// Note: machine name (hostname) comes from the block label, not a field
	// The hostname field is optional since it can be derived from the block label

	// Validate hostname if provided (optional)
	if hostname, exists := machine["hostname"]; exists {
		v.validateMachineHostnameV1(hostname, path, result)
	}

	if user, exists := machine["user"]; !exists {
		result.Errors = append(result.Errors, ValidationError{
			Field:    fmt.Sprintf("%s.user", path),
			Message:  "missing required field: user",
			Severity: "error",
		})
	} else {
		v.validateMachineUserV1(user, path, result)
	}

	// Check for authentication - it could be a block or flattened fields
	if auth, exists := machine["authentication"]; exists {
		v.validateMachineAuthenticationV1(auth, path, result)
	} else if _, exists := machine["password"]; exists {
		// Authentication is flattened - consider this valid since we have authentication data
	} else {
		result.Errors = append(result.Errors, ValidationError{
			Field:    fmt.Sprintf("%s.authentication", path),
			Message:  "missing required field: authentication",
			Severity: "error",
		})
	}
}

// validateActionNameV1 validates action name against V1 rules
func (v *Validator) validateActionNameV1(name interface{}, path string, result *ValidationResult) {
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
func (v *Validator) validateActionDescriptionV1(description interface{}, path string, result *ValidationResult) {
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
func (v *Validator) validateActionTypeV1(actionType interface{}, path string, result *ValidationResult) {
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
	validTypes := []string{"command", "template_deploy", "file_sync", "service_control"}
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

// validateVariableNameV1 validates variable name against V1 rules
func (v *Validator) validateVariableNameV1(name interface{}, path string, result *ValidationResult) {
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
func (v *Validator) validateVariableValueV1(value interface{}, path string, result *ValidationResult) {
	// Value can be any type, just ensure it's not nil
	if value == nil {
		result.Errors = append(result.Errors, ValidationError{
			Field:    fmt.Sprintf("%s.value", path),
			Message:  "value cannot be null",
			Severity: "error",
		})
	}
}

// validateGroupV1 validates an individual group against V1 schema
func (v *Validator) validateGroupV1(group map[string]interface{}, path string, result *ValidationResult) {
	// Note: Group validation not yet implemented
}

// validateMachineNameV1 validates machine name against V1 rules
func (v *Validator) validateMachineNameV1(name interface{}, path string, result *ValidationResult) {
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
func (v *Validator) validateMachineHostnameV1(hostname interface{}, path string, result *ValidationResult) {
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
func (v *Validator) validateMachineUserV1(user interface{}, path string, result *ValidationResult) {
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
func (v *Validator) validateMachineAuthenticationV1(auth interface{}, path string, result *ValidationResult) {
	authMap, ok := auth.(map[string]interface{})
	if !ok {
		result.Errors = append(result.Errors, ValidationError{
			Field:    fmt.Sprintf("%s.authentication", path),
			Message:  "authentication must be a configuration block",
			Severity: "error",
		})
		return
	}

	// In HCL, authentication method is a block label, not a field
	// Look for authentication method blocks
	hasValidMethod := false
	for authMethod, authValue := range authMap {
		if authMethod == "password" || authMethod == "publickey" || authMethod == "certificate" {
			hasValidMethod = true
			// Validate method-specific fields
			if authMethodMap, ok := authValue.(map[string]interface{}); ok {
				switch authMethod {
				case "publickey":
					v.validatePublicKeyAuthV1(authMethodMap, path, result)
				case "password":
					v.validatePasswordAuthV1(authMethodMap, path, result)
				case "certificate":
					v.validateCertificateAuthV1(authMethodMap, path, result)
				}
			}
		}
	}

	if !hasValidMethod {
		result.Errors = append(result.Errors, ValidationError{
			Field:    fmt.Sprintf("%s.authentication", path),
			Message:  "authentication must have a valid method block (password, publickey, or certificate)",
			Severity: "error",
		})
	}
}

// validateAuthMethodV1 validates authentication method against V1 rules
func (v *Validator) validateAuthMethodV1(method interface{}, path string, result *ValidationResult) {
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
func (v *Validator) validatePublicKeyAuthV1(auth map[string]interface{}, path string, result *ValidationResult) {
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
func (v *Validator) validatePasswordAuthV1(auth map[string]interface{}, path string, result *ValidationResult) {
	// Password field is required for password method
	if password, exists := auth["password"]; !exists {
		result.Errors = append(result.Errors, ValidationError{
			Field:    fmt.Sprintf("%s.authentication.password", path),
			Message:  "password is required when method is 'password'",
			Severity: "error",
		})
	} else {
		// Validate password is a string
		if _, ok := password.(string); !ok {
			result.Errors = append(result.Errors, ValidationError{
				Field:    fmt.Sprintf("%s.authentication.password", path),
				Message:  "password must be a string",
				Severity: "error",
			})
		}
	}
}

// validateCertificateAuthV1 validates certificate authentication fields
func (v *Validator) validateCertificateAuthV1(auth map[string]interface{}, path string, result *ValidationResult) {
	// Certificate fields are required for certificate method
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

	if certificatePath, exists := auth["certificate_path"]; !exists {
		result.Errors = append(result.Errors, ValidationError{
			Field:    fmt.Sprintf("%s.authentication.certificate_path", path),
			Message:  "certificate_path is required when method is 'certificate'",
			Severity: "error",
		})
	} else {
		// Validate certificate_path is a string
		if _, ok := certificatePath.(string); !ok {
			result.Errors = append(result.Errors, ValidationError{
				Field:    fmt.Sprintf("%s.authentication.certificate_path", path),
				Message:  "certificate_path must be a string",
				Severity: "error",
			})
		}
	}
}

func (v *Validator) validateActionsV1(content map[string]interface{}, result *ValidationResult) {
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
		v.validateActionsBlockV1(actionsMap, result)
	} else {
		result.Errors = append(result.Errors, ValidationError{
			Field:    "actions",
			Message:  "actions block must be a configuration block",
			Severity: "error",
		})
	}
}

// validateActionsBlockV1 validates the actions block contents against V1 schema
func (v *Validator) validateActionsBlockV1(actions map[string]interface{}, result *ValidationResult) {
	// Check for action list
	if actionList, exists := actions["action"]; exists {
		if actionArray, ok := actionList.([]interface{}); ok {
			for i, action := range actionArray {
				if actionMap, ok := action.(map[string]interface{}); ok {
					v.validateActionV1(actionMap, fmt.Sprintf("actions.action[%d]", i), result)
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
func (v *Validator) validateActionV1(action map[string]interface{}, path string, result *ValidationResult) {
	// Required fields
	if name, exists := action["name"]; !exists {
		result.Errors = append(result.Errors, ValidationError{
			Field:    fmt.Sprintf("%s.name", path),
			Message:  "missing required field: name",
			Severity: "error",
		})
	} else {
		v.validateActionNameV1(name, path, result)
	}

	if description, exists := action["description"]; !exists {
		result.Errors = append(result.Errors, ValidationError{
			Field:    fmt.Sprintf("%s.description", path),
			Message:  "missing required field: description",
			Severity: "error",
		})
	} else {
		v.validateActionDescriptionV1(description, path, result)
	}

	if actionType, exists := action["type"]; !exists {
		result.Errors = append(result.Errors, ValidationError{
			Field:    fmt.Sprintf("%s.type", path),
			Message:  "missing required field: type",
			Severity: "error",
		})
	} else {
		v.validateActionTypeV1(actionType, path, result)
	}
}

func (v *Validator) validateVariablesV1(content map[string]interface{}, result *ValidationResult) {
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
		v.validateVariablesBlockV1(variablesMap, result)
	} else {
		result.Errors = append(result.Errors, ValidationError{
			Field:    "variables",
			Message:  "variables block must be a configuration block",
			Severity: "error",
		})
	}
}

// validateVariablesBlockV1 validates the variables block contents against V1 schema
func (v *Validator) validateVariablesBlockV1(variables map[string]interface{}, result *ValidationResult) {
	// Check for variable list
	if variableList, exists := variables["variable"]; exists {
		if variableArray, ok := variableList.([]interface{}); ok {
			// Handle array of variables
			for i, variable := range variableArray {
				if variableMap, ok := variable.(map[string]interface{}); ok {
					v.validateVariableV1(variableMap, fmt.Sprintf("variables.variable[%d]", i), result)
				}
			}
		} else if variableMap, ok := variableList.(map[string]interface{}); ok {
			// Handle single variable
			v.validateVariableV1(variableMap, "variables.variable[0]", result)
		}
	} else {
		result.Warnings = append(result.Warnings, ValidationWarning{
			Field:   "variables",
			Message: "variables block should contain at least one variable",
		})
	}
}

// validateVariableV1 validates an individual variable against V1 schema
func (v *Validator) validateVariableV1(variable map[string]interface{}, path string, result *ValidationResult) {
	// Required fields
	if name, exists := variable["name"]; !exists {
		result.Errors = append(result.Errors, ValidationError{
			Field:    fmt.Sprintf("%s.name", path),
			Message:  "missing required field: name",
			Severity: "error",
		})
	} else {
		v.validateVariableNameV1(name, path, result)
	}

	if value, exists := variable["value"]; !exists {
		result.Errors = append(result.Errors, ValidationError{
			Field:    fmt.Sprintf("%s.value", path),
			Message:  "missing required field: value",
			Severity: "error",
		})
	} else {
		v.validateVariableValueV1(value, path, result)
	}
}

func (v *Validator) validateLoggingV1(content map[string]interface{}, result *ValidationResult) {
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
		v.validateLoggingBlockV1(loggingMap, result)
	} else {
		result.Errors = append(result.Errors, ValidationError{
			Field:    "logging",
			Message:  "logging block must be a configuration block",
			Severity: "error",
		})
	}
}

// validateLoggingBlockV1 validates the logging block contents against V1 schema
func (v *Validator) validateLoggingBlockV1(logging map[string]interface{}, result *ValidationResult) {
	// Validate level if present
	if level, exists := logging["level"]; exists {
		v.validateLoggingLevelV1(level, result)
	}

	// Validate format if present
	if format, exists := logging["format"]; exists {
		v.validateLoggingFormatV1(format, result)
	}

	// Validate output if present
	if output, exists := logging["output"]; exists {
		v.validateLoggingOutputV1(output, result)
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

func (v *Validator) validateFactsV1(content map[string]interface{}, result *ValidationResult) {
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
		v.validateFactsBlockV1(factsMap, result)
	} else {
		result.Errors = append(result.Errors, ValidationError{
			Field:    "facts",
			Message:  "facts block must be a configuration block",
			Severity: "error",
		})
	}
}

// validateFactsBlockV1 validates the hierarchical facts block contents against V1 schema
func (v *Validator) validateFactsBlockV1(facts map[string]interface{}, result *ValidationResult) {
	// Validate basic_facts if present
	if basicFacts, exists := facts["basic_facts"]; exists {
		if basicFactsMap, ok := basicFacts.(map[string]interface{}); ok {
			v.validateBasicFactsV1(basicFactsMap, result)
		} else {
			result.Errors = append(result.Errors, ValidationError{
				Field:    "facts.basic_facts",
				Message:  "basic_facts must be a configuration block",
				Severity: "error",
			})
		}
	}

	// Validate enhanced_facts if present
	if enhancedFacts, exists := facts["enhanced_facts"]; exists {
		if enhancedFactsMap, ok := enhancedFacts.(map[string]interface{}); ok {
			v.validateEnhancedFactsV1(enhancedFactsMap, result)
		} else {
			result.Errors = append(result.Errors, ValidationError{
				Field:    "facts.enhanced_facts",
				Message:  "enhanced_facts must be a configuration block",
				Severity: "error",
			})
		}
	}

	// Validate custom_facts if present
	if customFacts, exists := facts["custom_facts"]; exists {
		if customFactsMap, ok := customFacts.(map[string]interface{}); ok {
			v.validateCustomFactsV1(customFactsMap, result)
		} else {
			result.Errors = append(result.Errors, ValidationError{
				Field:    "facts.custom_facts",
				Message:  "custom_facts must be a configuration block",
				Severity: "error",
			})
		}
	}
}

// validateBasicFactsV1 validates basic facts gathered via SSH commands
func (v *Validator) validateBasicFactsV1(basicFacts map[string]interface{}, result *ValidationResult) {
	// Validate each fact category
	if systemFacts, exists := basicFacts["system_facts"]; exists {
		if systemFactsMap, ok := systemFacts.(map[string]interface{}); ok {
			if facts, exists := systemFactsMap["facts"]; exists {
				if factsMap, ok := facts.(map[string]interface{}); ok {
					for factName, factData := range factsMap {
						if factMap, ok := factData.(map[string]interface{}); ok {
							v.validateIndividualFactV1(factName, factMap, result, "basic_facts.system_facts")
						}
					}
				}
			}
		}
	}

	if hardwareFacts, exists := basicFacts["hardware_facts"]; exists {
		if hardwareFactsMap, ok := hardwareFacts.(map[string]interface{}); ok {
			if facts, exists := hardwareFactsMap["facts"]; exists {
				if factsMap, ok := facts.(map[string]interface{}); ok {
					for factName, factData := range factsMap {
						if factMap, ok := factData.(map[string]interface{}); ok {
							v.validateIndividualFactV1(factName, factMap, result, "basic_facts.hardware_facts")
						}
					}
				}
			}
		}
	}

	if networkFacts, exists := basicFacts["network_facts"]; exists {
		if networkFactsMap, ok := networkFacts.(map[string]interface{}); ok {
			if facts, exists := networkFactsMap["facts"]; exists {
				if factsMap, ok := facts.(map[string]interface{}); ok {
					for factName, factData := range factsMap {
						if factMap, ok := factData.(map[string]interface{}); ok {
							v.validateIndividualFactV1(factName, factMap, result, "basic_facts.network_facts")
						}
					}
				}
			}
		}
	}

	if osFacts, exists := basicFacts["os_facts"]; exists {
		if osFactsMap, ok := osFacts.(map[string]interface{}); ok {
			if facts, exists := osFactsMap["facts"]; exists {
				if factsMap, ok := facts.(map[string]interface{}); ok {
					for factName, factData := range factsMap {
						if factMap, ok := factData.(map[string]interface{}); ok {
							v.validateIndividualFactV1(factName, factMap, result, "basic_facts.os_facts")
						}
					}
				}
			}
		}
	}

	if userFacts, exists := basicFacts["user_facts"]; exists {
		if userFactsMap, ok := userFacts.(map[string]interface{}); ok {
			if facts, exists := userFactsMap["facts"]; exists {
				if factsMap, ok := facts.(map[string]interface{}); ok {
					for factName, factData := range factsMap {
						if factMap, ok := factData.(map[string]interface{}); ok {
							v.validateIndividualFactV1(factName, factMap, result, "basic_facts.user_facts")
						}
					}
				}
			}
		}
	}

	if runtimeFacts, exists := basicFacts["runtime_facts"]; exists {
		if runtimeFactsMap, ok := runtimeFacts.(map[string]interface{}); ok {
			if facts, exists := runtimeFactsMap["facts"]; exists {
				if factsMap, ok := facts.(map[string]interface{}); ok {
					for factName, factData := range factsMap {
						if factMap, ok := factData.(map[string]interface{}); ok {
							v.validateIndividualFactV1(factName, factMap, result, "basic_facts.runtime_facts")
						}
					}
				}
			}
		}
	}
}

// validateEnhancedFactsV1 validates enhanced facts gathered by spooky-facts
func (v *Validator) validateEnhancedFactsV1(enhancedFacts map[string]interface{}, result *ValidationResult) {
	// Validate individual facts
	for factName, factData := range enhancedFacts {
		if factMap, ok := factData.(map[string]interface{}); ok {
			v.validateIndividualFactV1(factName, factMap, result, "enhanced_facts")
		} else {
			result.Errors = append(result.Errors, ValidationError{
				Field:    fmt.Sprintf("facts.enhanced_facts.%s", factName),
				Message:  "fact must be a configuration block",
				Severity: "error",
			})
		}
	}
}

// validateCustomFactsV1 validates custom facts (including age-encrypted ones)
func (v *Validator) validateCustomFactsV1(customFacts map[string]interface{}, result *ValidationResult) {
	// Validate individual facts
	for factName, factData := range customFacts {
		if factMap, ok := factData.(map[string]interface{}); ok {
			v.validateIndividualFactV1(factName, factMap, result, "custom_facts")
		} else {
			result.Errors = append(result.Errors, ValidationError{
				Field:    fmt.Sprintf("facts.custom_facts.%s", factName),
				Message:  "fact must be a configuration block",
				Severity: "error",
			})
		}
	}
}

// validateIndividualFactV1 validates a single fact block
func (v *Validator) validateIndividualFactV1(factName string, fact map[string]interface{}, result *ValidationResult, context string) {
	// Validate fact name (from block label)
	v.validateFactNameV1(factName, result)

	// Required fields
	if value, exists := fact["value"]; !exists {
		result.Errors = append(result.Errors, ValidationError{
			Field:    fmt.Sprintf("facts.%s.%s.value", context, factName),
			Message:  "missing required field: value",
			Severity: "error",
		})
	} else {
		// Value can be any type, just ensure it's not nil
		if value == nil {
			result.Errors = append(result.Errors, ValidationError{
				Field:    fmt.Sprintf("facts.%s.%s.value", context, factName),
				Message:  "value cannot be null",
				Severity: "error",
			})
		}
	}

	if factType, exists := fact["type"]; !exists {
		result.Errors = append(result.Errors, ValidationError{
			Field:    fmt.Sprintf("facts.%s.%s.type", context, factName),
			Message:  "missing required field: type",
			Severity: "error",
		})
	} else {
		v.validateFactTypeV1(factType, result)
	}

	// Optional fields
	if encrypted, exists := fact["encrypted"]; exists {
		if _, ok := encrypted.(bool); !ok {
			result.Errors = append(result.Errors, ValidationError{
				Field:    fmt.Sprintf("facts.%s.%s.encrypted", context, factName),
				Message:  "encrypted must be a boolean",
				Severity: "error",
			})
		}
	}

	// Validate encrypted_value if present
	if encryptedValue, exists := fact["encrypted_value"]; exists {
		if encryptedValueMap, ok := encryptedValue.(map[string]interface{}); ok {
			v.validateEncryptedValueV1(encryptedValueMap, result, fmt.Sprintf("facts.%s.%s.encrypted_value", context, factName))
		} else {
			result.Errors = append(result.Errors, ValidationError{
				Field:    fmt.Sprintf("facts.%s.%s.encrypted_value", context, factName),
				Message:  "encrypted_value must be a configuration block",
				Severity: "error",
			})
		}
	}

	// Validate sensitive flag if present
	if sensitive, exists := fact["sensitive"]; exists {
		if _, ok := sensitive.(bool); !ok {
			result.Errors = append(result.Errors, ValidationError{
				Field:    fmt.Sprintf("facts.%s.%s.sensitive", context, factName),
				Message:  "sensitive must be a boolean",
				Severity: "error",
			})
		}
	}
}

func (v *Validator) validateSpookyV1(content map[string]interface{}, result *ValidationResult) {
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
		v.validateSpookyBlockV1(spookyMap, result)
	} else {
		result.Errors = append(result.Errors, ValidationError{
			Field:    "spooky",
			Message:  "spooky block must be a configuration block",
			Severity: "error",
		})
	}
}

// validateSpookyBlockV1 validates the spooky block contents against V1 schema
func (v *Validator) validateSpookyBlockV1(spooky map[string]interface{}, result *ValidationResult) {
	// Validate SSH configuration if present
	if ssh, exists := spooky["ssh"]; exists {
		if sshMap, ok := ssh.(map[string]interface{}); ok {
			v.validateSpookySSHV1(sshMap, result)
		}
	}

	// Validate security configuration if present
	if security, exists := spooky["security"]; exists {
		if securityMap, ok := security.(map[string]interface{}); ok {
			v.validateSpookySecurityV1(securityMap, result)
		}
	}

	// Validate age configuration if present
	if age, exists := spooky["age"]; exists {
		if ageMap, ok := age.(map[string]interface{}); ok {
			v.validateSpookyAgeV1(ageMap, result)
		}
	}

	// Validate logging configuration if present
	if logging, exists := spooky["logging"]; exists {
		if loggingMap, ok := logging.(map[string]interface{}); ok {
			v.validateSpookyLoggingV1(loggingMap, result)
		}
	}
}

func (v *Validator) validateProjectDirectoryV1(content map[string]interface{}, result *ValidationResult) {
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
		v.validateProjectDirectoryBlockV1(projectDirMap, result)
	} else {
		result.Errors = append(result.Errors, ValidationError{
			Field:    "project_directory",
			Message:  "project_directory block must be a configuration block",
			Severity: "error",
		})
	}
}

// validateProjectDirectoryBlockV1 validates the project_directory block contents against V1 schema
func (v *Validator) validateProjectDirectoryBlockV1(projectDir map[string]interface{}, result *ValidationResult) {
	// Validate required name field
	if name, exists := projectDir["name"]; !exists {
		result.Errors = append(result.Errors, ValidationError{
			Field:    "project_directory.name",
			Message:  "missing required field: name",
			Severity: "error",
		})
	} else {
		v.validateProjectDirectoryNameV1(name, result)
	}

	// Validate files array
	if files, exists := projectDir["files"]; exists {
		if filesArray, ok := files.([]interface{}); ok {
			for i, file := range filesArray {
				if fileMap, ok := file.(map[string]interface{}); ok {
					v.validateProjectDirectoryFileV1(fileMap, fmt.Sprintf("project_directory.files[%d]", i), result)
				}
			}
		}
	}

	// Validate directories array
	if directories, exists := projectDir["directories"]; exists {
		if dirsArray, ok := directories.([]interface{}); ok {
			for i, dir := range dirsArray {
				if dirMap, ok := dir.(map[string]interface{}); ok {
					v.validateProjectDirectoryDirV1(dirMap, fmt.Sprintf("project_directory.directories[%d]", i), result)
				}
			}
		}
	}
}

// validateLoggingLevelV1 validates logging level against V1 rules
func (v *Validator) validateLoggingLevelV1(level interface{}, result *ValidationResult) {
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
func (v *Validator) validateLoggingFormatV1(format interface{}, result *ValidationResult) {
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
func (v *Validator) validateLoggingOutputV1(output interface{}, result *ValidationResult) {
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

// validateFactNameV1 validates fact name against V1 rules
func (v *Validator) validateFactNameV1(name interface{}, result *ValidationResult) {
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
func (v *Validator) validateFactTypeV1(factType interface{}, result *ValidationResult) {
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

// validateEncryptedValueV1 validates an encrypted value structure
func (v *Validator) validateEncryptedValueV1(encryptedValue map[string]interface{}, result *ValidationResult, context string) {
	// Validate data field (required)
	if data, exists := encryptedValue["data"]; !exists {
		result.Errors = append(result.Errors, ValidationError{
			Field:    fmt.Sprintf("%s.data", context),
			Message:  "missing required field: data",
			Severity: "error",
		})
	} else {
		if dataStr, ok := data.(string); !ok || dataStr == "" {
			result.Errors = append(result.Errors, ValidationError{
				Field:    fmt.Sprintf("%s.data", context),
				Message:  "data must be a non-empty string",
				Severity: "error",
			})
		}
	}

	// Validate format field (optional, defaults to "base64")
	if format, exists := encryptedValue["format"]; exists {
		if formatStr, ok := format.(string); !ok {
			result.Errors = append(result.Errors, ValidationError{
				Field:    fmt.Sprintf("%s.format", context),
				Message:  "format must be a string",
				Severity: "error",
			})
		} else {
			validFormats := []string{"base64", "armored", "compact"}
			valid := false
			for _, validFormat := range validFormats {
				if formatStr == validFormat {
					valid = true
					break
				}
			}
			if !valid {
				result.Errors = append(result.Errors, ValidationError{
					Field:    fmt.Sprintf("%s.format", context),
					Value:    formatStr,
					Message:  "format must be one of: base64, armored, compact",
					Severity: "error",
				})
			}
		}
	}

	// Validate algorithm field (optional, defaults to "age")
	if algorithm, exists := encryptedValue["algorithm"]; exists {
		if algorithmStr, ok := algorithm.(string); !ok || algorithmStr == "" {
			result.Errors = append(result.Errors, ValidationError{
				Field:    fmt.Sprintf("%s.algorithm", context),
				Message:  "algorithm must be a non-empty string",
				Severity: "error",
			})
		}
	}

	// Validate version field (optional, defaults to "v1")
	if version, exists := encryptedValue["version"]; exists {
		if versionStr, ok := version.(string); !ok || versionStr == "" {
			result.Errors = append(result.Errors, ValidationError{
				Field:    fmt.Sprintf("%s.version", context),
				Message:  "version must be a non-empty string",
				Severity: "error",
			})
		}
	}

	// Validate encrypted_at field (optional)
	if encryptedAt, exists := encryptedValue["encrypted_at"]; exists {
		if encryptedAtStr, ok := encryptedAt.(string); !ok || encryptedAtStr == "" {
			result.Errors = append(result.Errors, ValidationError{
				Field:    fmt.Sprintf("%s.encrypted_at", context),
				Message:  "encrypted_at must be a non-empty string",
				Severity: "error",
			})
		}
	}
}

// validateSpookySSHV1 validates SSH configuration in spooky config
func (v *Validator) validateSpookySSHV1(ssh map[string]interface{}, result *ValidationResult) {
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
func (v *Validator) validateSpookySecurityV1(security map[string]interface{}, result *ValidationResult) {
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
func (v *Validator) validateSpookyAgeV1(age map[string]interface{}, result *ValidationResult) {
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
func (v *Validator) validateSpookyLoggingV1(logging map[string]interface{}, result *ValidationResult) {
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

// validateProjectDirectoryNameV1 validates project directory name against V1 rules
func (v *Validator) validateProjectDirectoryNameV1(name interface{}, result *ValidationResult) {
	nameStr, ok := name.(string)
	if !ok {
		result.Errors = append(result.Errors, ValidationError{
			Field:    "project_directory.name",
			Message:  "name must be a string",
			Severity: "error",
		})
		return
	}

	if nameStr == "" {
		result.Errors = append(result.Errors, ValidationError{
			Field:    "project_directory.name",
			Message:  "name cannot be empty",
			Severity: "error",
		})
	}
}

// validateProjectDirectoryFileV1 validates a project directory file against V1 schema
func (v *Validator) validateProjectDirectoryFileV1(file map[string]interface{}, path string, result *ValidationResult) {
	// Required fields
	if name, exists := file["name"]; !exists {
		result.Errors = append(result.Errors, ValidationError{
			Field:    fmt.Sprintf("%s.name", path),
			Message:  "missing required field: name",
			Severity: "error",
		})
	} else {
		v.validateProjectDirectoryFileNameV1(name, path, result)
	}

	if fileType, exists := file["type"]; !exists {
		result.Errors = append(result.Errors, ValidationError{
			Field:    fmt.Sprintf("%s.type", path),
			Message:  "missing required field: type",
			Severity: "error",
		})
	} else {
		v.validateProjectDirectoryFileTypeV1(fileType, path, result)
	}
}

// validateProjectDirectoryFileNameV1 validates project directory file name
func (v *Validator) validateProjectDirectoryFileNameV1(name interface{}, path string, result *ValidationResult) {
	nameStr, ok := name.(string)
	if !ok {
		result.Errors = append(result.Errors, ValidationError{
			Field:    fmt.Sprintf("%s.name", path),
			Message:  "name must be a string",
			Severity: "error",
		})
		return
	}

	if nameStr == "" {
		result.Errors = append(result.Errors, ValidationError{
			Field:    fmt.Sprintf("%s.name", path),
			Message:  "name cannot be empty",
			Severity: "error",
		})
	}
}

// validateProjectDirectoryFileTypeV1 validates project directory file type
func (v *Validator) validateProjectDirectoryFileTypeV1(fileType interface{}, path string, result *ValidationResult) {
	typeStr, ok := fileType.(string)
	if !ok {
		result.Errors = append(result.Errors, ValidationError{
			Field:    fmt.Sprintf("%s.type", path),
			Message:  "type must be a string",
			Severity: "error",
		})
		return
	}

	if typeStr != "file" {
		result.Errors = append(result.Errors, ValidationError{
			Field:    fmt.Sprintf("%s.type", path),
			Value:    typeStr,
			Message:  "type must be 'file'",
			Severity: "error",
		})
	}
}

// validateProjectDirectoryDirV1 validates a project directory directory against V1 schema
func (v *Validator) validateProjectDirectoryDirV1(dir map[string]interface{}, path string, result *ValidationResult) {
	// Required fields
	if name, exists := dir["name"]; !exists {
		result.Errors = append(result.Errors, ValidationError{
			Field:    fmt.Sprintf("%s.name", path),
			Message:  "missing required field: name",
			Severity: "error",
		})
	} else {
		v.validateProjectDirectoryDirNameV1(name, path, result)
	}

	if dirType, exists := dir["type"]; !exists {
		result.Errors = append(result.Errors, ValidationError{
			Field:    fmt.Sprintf("%s.type", path),
			Message:  "missing required field: type",
			Severity: "error",
		})
	} else {
		v.validateProjectDirectoryDirTypeV1(dirType, path, result)
	}
}

// validateProjectDirectoryDirNameV1 validates project directory directory name
func (v *Validator) validateProjectDirectoryDirNameV1(name interface{}, path string, result *ValidationResult) {
	nameStr, ok := name.(string)
	if !ok {
		result.Errors = append(result.Errors, ValidationError{
			Field:    fmt.Sprintf("%s.name", path),
			Message:  "name must be a string",
			Severity: "error",
		})
		return
	}

	if nameStr == "" {
		result.Errors = append(result.Errors, ValidationError{
			Field:    fmt.Sprintf("%s.name", path),
			Message:  "name cannot be empty",
			Severity: "error",
		})
	}
}

// validateProjectDirectoryDirTypeV1 validates project directory directory type
func (v *Validator) validateProjectDirectoryDirTypeV1(dirType interface{}, path string, result *ValidationResult) {
	typeStr, ok := dirType.(string)
	if !ok {
		result.Errors = append(result.Errors, ValidationError{
			Field:    fmt.Sprintf("%s.type", path),
			Message:  "type must be a string",
			Severity: "error",
		})
		return
	}

	if typeStr != "directory" {
		result.Errors = append(result.Errors, ValidationError{
			Field:    fmt.Sprintf("%s.type", path),
			Value:    typeStr,
			Message:  "type must be 'directory'",
			Severity: "error",
		})
	}
}
