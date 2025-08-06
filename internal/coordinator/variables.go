package coordinator

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"text/template"
	"time"

	"spooky/internal/interfaces"
	"spooky/internal/logging"
	"spooky/internal/variables"
	"spooky/internal/variables/types"
)

// CoordinatorVariablesIntegration implements variables system integration
type CoordinatorVariablesIntegration struct {
	variablesManager variables.VariableManager
	logger           logging.Logger
}

// NewCoordinatorVariablesIntegration creates a new variables integration
func NewCoordinatorVariablesIntegration(variablesManager variables.VariableManager, logger logging.Logger) *CoordinatorVariablesIntegration {
	return &CoordinatorVariablesIntegration{
		variablesManager: variablesManager,
		logger:           logger,
	}
}

// LoadVariables loads variables from the project
func (vi *CoordinatorVariablesIntegration) LoadVariables(projectPath string) (*interfaces.VariablesContext, error) {
	variablesContext := &interfaces.VariablesContext{
		BaseContext: interfaces.BaseContext{
			ProjectPath: projectPath,
			Timestamp:   time.Now(),
		},
		ResolvedVariables: make(map[string]interface{}),
		VariableContext:   make(map[string]interface{}),
		ResolutionContext: make(map[string]interface{}),
	}

	// Load variables from project using variables manager
	if vi.variablesManager != nil {
		collection, err := vi.variablesManager.LoadVariablesForProject(projectPath)
		if err != nil {
			return nil, fmt.Errorf("failed to load variables from project: %w", err)
		}

		// Convert variable collection to interface context
		for _, variable := range collection.Variables {
			if variable != nil {
				variablesContext.VariableContext[variable.Name] = variable.Value
			}
		}

		// Create variable context for resolution
		varContext, err := vi.variablesManager.CreateContext(context.Background(), collection.Variables)
		if err != nil {
			return nil, fmt.Errorf("failed to create variable context: %w", err)
		}

		// Resolve all variables
		err = vi.variablesManager.ResolveVariablesForContext(varContext)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve variables: %w", err)
		}

		// Convert resolved variables to interface context
		for name, variable := range varContext.Variables {
			if variable != nil {
				variablesContext.ResolvedVariables[name] = variable.Value
			}
		}
	}

	vi.logger.Info("Loaded variables from project",
		logging.String("project", projectPath),
		logging.Int("variables_count", len(variablesContext.ResolvedVariables)))

	return variablesContext, nil
}

// ResolveVariables resolves variables using facts context with advanced features
func (vi *CoordinatorVariablesIntegration) ResolveVariables(variablesContext *interfaces.VariablesContext, factsContext *interfaces.FactsContext) error {
	if variablesContext == nil {
		return fmt.Errorf("variables context cannot be nil")
	}

	if factsContext == nil {
		return fmt.Errorf("facts context cannot be nil")
	}

	// Resolve variables using facts if variables manager is available
	if vi.variablesManager != nil {
		// Convert facts context to map for resolution
		factsMap := make(map[string]interface{})

		// Add machine facts with proper scoping
		for machine, factCollection := range factsContext.MachineFacts {
			if factCollection != nil && factCollection.Facts != nil {
				for key, value := range factCollection.Facts {
					// Create scoped variable names
					factsMap[fmt.Sprintf("%s.%s", machine, key)] = value
					factsMap[fmt.Sprintf("machines.%s.%s", machine, key)] = value
				}
			}
		}

		// Add global facts with scoping
		if factsContext.GlobalFacts != nil && factsContext.GlobalFacts.Facts != nil {
			for key, value := range factsContext.GlobalFacts.Facts {
				factsMap[fmt.Sprintf("global.%s", key)] = value
				factsMap[key] = value // Also add without prefix for backward compatibility
			}
		}

		// Add project facts with scoping
		if factsContext.ProjectFacts != nil && factsContext.ProjectFacts.Facts != nil {
			for key, value := range factsContext.ProjectFacts.Facts {
				factsMap[fmt.Sprintf("project.%s", key)] = value
			}
		}

		// Add computed variables
		computedVars := vi.computeDerivedVariables(factsMap)
		for key, value := range computedVars {
			factsMap[key] = value
		}

		// Create variable context from resolved variables
		varContext := &types.VariableContext{
			Variables:   make(map[string]*types.Variable),
			ProjectPath: variablesContext.ProjectPath,
			Environment: make(map[string]string),
		}

		// Add resolved variables to context
		for name, value := range variablesContext.ResolvedVariables {
			varContext.Variables[name] = &types.Variable{
				Name:  name,
				Value: value,
				Type:  "string", // Default type
			}
		}

		// Resolve variables using the new manager
		if err := vi.variablesManager.ResolveVariablesForContext(varContext); err != nil {
			return fmt.Errorf("failed to resolve variables: %w", err)
		}

		// Update variables context with resolved values
		for name, variable := range varContext.Variables {
			variablesContext.ResolvedVariables[name] = variable.Value
		}
	}

	vi.logger.Info("Resolved variables with facts context",
		logging.Int("resolved_count", len(variablesContext.ResolvedVariables)))

	return nil
}

// computeDerivedVariables computes derived variables from facts
func (vi *CoordinatorVariablesIntegration) computeDerivedVariables(facts map[string]interface{}) map[string]interface{} {
	computed := make(map[string]interface{})

	// Compute machine count
	machineCount := 0
	for key := range facts {
		if strings.HasPrefix(key, "machines.") && !strings.Contains(key, ".") {
			machineCount++
		}
	}
	computed["machine_count"] = machineCount

	// Compute environment variables
	if envVars, ok := facts["global.environment"].(map[string]interface{}); ok {
		for key, value := range envVars {
			computed[fmt.Sprintf("env.%s", key)] = value
		}
	}

	// Compute system information
	if sysInfo, ok := facts["global.system"].(map[string]interface{}); ok {
		for key, value := range sysInfo {
			computed[fmt.Sprintf("system.%s", key)] = value
		}
	}

	return computed
}

// ValidateVariables validates variables data
func (vi *CoordinatorVariablesIntegration) ValidateVariables(variablesContext *interfaces.VariablesContext) error {
	if variablesContext == nil {
		return fmt.Errorf("variables context cannot be nil")
	}

	// Basic validation - can be enhanced based on requirements
	// Check for required variables
	// Validate variable types
	// Check for circular references

	return nil
}

// SubstituteVariables substitutes variables in a template string using Go template engine
func (vi *CoordinatorVariablesIntegration) SubstituteVariables(templateStr string, variablesContext *interfaces.VariablesContext) (string, error) {
	if templateStr == "" {
		return "", fmt.Errorf("template string cannot be empty")
	}

	if variablesContext == nil {
		return "", fmt.Errorf("variables context cannot be nil")
	}

	// Create a combined data map for template rendering
	data := make(map[string]interface{})

	// Add resolved variables (highest priority)
	for name, value := range variablesContext.ResolvedVariables {
		data[name] = value
	}

	// Add variable context (fallback)
	for name, value := range variablesContext.VariableContext {
		if _, exists := data[name]; !exists {
			data[name] = value
		}
	}

	// Add resolution context (lowest priority)
	for name, value := range variablesContext.ResolutionContext {
		if _, exists := data[name]; !exists {
			data[name] = value
		}
	}

	// Create Go template with custom functions
	tmpl, err := template.New("variables").Funcs(template.FuncMap{
		"default": func(defaultVal interface{}, val interface{}) interface{} {
			if val == nil || val == "" {
				return defaultVal
			}
			return val
		},
		"join": func(sep string, vals []string) string {
			return strings.Join(vals, sep)
		},
		"split": func(sep, str string) []string {
			return strings.Split(str, sep)
		},
		"upper": func(str string) string {
			return strings.ToUpper(str)
		},
		"lower": func(str string) string {
			return strings.ToLower(str)
		},
		"title": func(str string) string {
			return strings.Title(str)
		},
		"trim": func(str string) string {
			return strings.TrimSpace(str)
		},
		"replace": func(old, new, str string) string {
			return strings.ReplaceAll(str, old, new)
		},
		"contains": func(substr, str string) bool {
			return strings.Contains(str, substr)
		},
		"hasPrefix": func(prefix, str string) bool {
			return strings.HasPrefix(str, prefix)
		},
		"hasSuffix": func(suffix, str string) bool {
			return strings.HasSuffix(str, suffix)
		},
	}).Parse(templateStr)

	if err != nil {
		// If template parsing fails, fall back to simple string replacement
		return vi.simpleVariableSubstitution(templateStr, data), nil
	}

	// Execute template
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		// If template execution fails, fall back to simple string replacement
		return vi.simpleVariableSubstitution(templateStr, data), nil
	}

	result := buf.String()

	vi.logger.Info("Substituted variables in template",
		logging.String("template_length", fmt.Sprintf("%d", len(templateStr))),
		logging.String("result_length", fmt.Sprintf("%d", len(result))))

	return result, nil
}

// simpleVariableSubstitution performs simple {{variable}} substitution
func (vi *CoordinatorVariablesIntegration) simpleVariableSubstitution(templateStr string, data map[string]interface{}) string {
	result := templateStr

	// Substitute variables using {{variable}} syntax
	for name, value := range data {
		placeholder := fmt.Sprintf("{{%s}}", name)
		if strings.Contains(result, placeholder) {
			result = strings.ReplaceAll(result, placeholder, fmt.Sprintf("%v", value))
		}
	}

	return result
}

// GetVariable gets a specific variable by name
func (vi *CoordinatorVariablesIntegration) GetVariable(name string, context *interfaces.VariablesContext) (interface{}, error) {
	if name == "" {
		return nil, fmt.Errorf("variable name cannot be empty")
	}

	if context == nil {
		return nil, fmt.Errorf("variables context cannot be nil")
	}

	// Look up variable in resolved variables
	if value, exists := context.ResolvedVariables[name]; exists {
		return value, nil
	}

	// Look up variable in variable context
	if value, exists := context.VariableContext[name]; exists {
		return value, nil
	}

	// Look up variable in resolution context
	if value, exists := context.ResolutionContext[name]; exists {
		return value, nil
	}

	return nil, fmt.Errorf("variable '%s' not found", name)
}

// SetVariable sets a variable value
func (vi *CoordinatorVariablesIntegration) SetVariable(name string, value interface{}, context *interfaces.VariablesContext) error {
	if name == "" {
		return fmt.Errorf("variable name cannot be empty")
	}

	if context == nil {
		return fmt.Errorf("variables context cannot be nil")
	}

	// Set variable in resolved variables
	context.ResolvedVariables[name] = value

	vi.logger.Debug("Set variable", logging.String("name", name))

	return nil
}

// ListVariables lists all available variables
func (vi *CoordinatorVariablesIntegration) ListVariables(context *interfaces.VariablesContext) (map[string]interface{}, error) {
	if context == nil {
		return nil, fmt.Errorf("variables context cannot be nil")
	}

	// Combine all variable sources
	allVariables := make(map[string]interface{})

	// Add resolved variables
	for name, value := range context.ResolvedVariables {
		allVariables[name] = value
	}

	// Add variable context
	for name, value := range context.VariableContext {
		allVariables[name] = value
	}

	// Add resolution context
	for name, value := range context.ResolutionContext {
		allVariables[name] = value
	}

	return allVariables, nil
}

// ImportVariables imports variables from external sources
func (vi *CoordinatorVariablesIntegration) ImportVariables(source string, format string, context *interfaces.VariablesContext) error {
	if source == "" {
		return fmt.Errorf("source cannot be empty")
	}

	if format == "" {
		return fmt.Errorf("format cannot be empty")
	}

	if context == nil {
		return fmt.Errorf("variables context cannot be nil")
	}

	vi.logger.Info("Importing variables", logging.String("source", source), logging.String("format", format))

	// In a real implementation, this would:
	// - Read from the specified source (file, URL, etc.)
	// - Parse according to the specified format (JSON, HCL, etc.)
	// - Validate the imported variables
	// - Merge with existing variables

	// For now, we'll just log the import attempt
	vi.logger.Debug("Variable import not yet implemented",
		logging.String("source", source),
		logging.String("format", format))

	return nil
}

// ExportVariables exports variables to external formats
func (vi *CoordinatorVariablesIntegration) ExportVariables(format string, context *interfaces.VariablesContext) ([]byte, error) {
	if format == "" {
		return nil, fmt.Errorf("format cannot be empty")
	}

	if context == nil {
		return nil, fmt.Errorf("variables context cannot be nil")
	}

	vi.logger.Info("Exporting variables", logging.String("format", format))

	// In a real implementation, this would:
	// - Convert variables to the specified format
	// - Handle different export formats (JSON, HCL, etc.)
	// - Include metadata and validation information

	// For now, we'll return a simple JSON representation
	allVars := make(map[string]interface{})

	// Combine all variable sources
	for name, value := range context.ResolvedVariables {
		allVars[name] = value
	}

	for name, value := range context.VariableContext {
		if _, exists := allVars[name]; !exists {
			allVars[name] = value
		}
	}

	// Convert to JSON (simplified)
	result := fmt.Sprintf(`{"variables": %d, "format": "%s"}`, len(allVars), format)

	return []byte(result), nil
}

// ValidateVariableReferences validates variable references in templates
func (vi *CoordinatorVariablesIntegration) ValidateVariableReferences(templateStr string, context *interfaces.VariablesContext) ([]string, error) {
	if templateStr == "" {
		return nil, fmt.Errorf("template string cannot be empty")
	}

	if context == nil {
		return nil, fmt.Errorf("variables context cannot be nil")
	}

	var missingVars []string

	// Extract variable references from template
	// This is a simplified approach - in a real implementation, you'd parse the template properly
	references := vi.extractVariableReferences(templateStr)

	// Check if all references exist
	for _, ref := range references {
		if !vi.variableExists(ref, context) {
			missingVars = append(missingVars, ref)
		}
	}

	return missingVars, nil
}

// extractVariableReferences extracts variable references from template string
func (vi *CoordinatorVariablesIntegration) extractVariableReferences(templateStr string) []string {
	var references []string

	// Simple regex-like extraction for {{variable}} patterns
	// In a real implementation, you'd use proper template parsing
	parts := strings.Split(templateStr, "{{")
	for _, part := range parts {
		if strings.Contains(part, "}}") {
			varRef := strings.Split(part, "}}")[0]
			varRef = strings.TrimSpace(varRef)
			if varRef != "" {
				references = append(references, varRef)
			}
		}
	}

	return references
}

// variableExists checks if a variable exists in the context
func (vi *CoordinatorVariablesIntegration) variableExists(name string, context *interfaces.VariablesContext) bool {
	// Check in resolved variables
	if _, exists := context.ResolvedVariables[name]; exists {
		return true
	}

	// Check in variable context
	if _, exists := context.VariableContext[name]; exists {
		return true
	}

	// Check in resolution context
	if _, exists := context.ResolutionContext[name]; exists {
		return true
	}

	return false
}
