package loading

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"

	spookyconfigtypes "spooky/internal/config/types"
)

// Parser implements ConfigParser interface
type Parser struct {
	hclParser *hclparse.Parser
}

// NewConfigParser creates a new config parser
func NewConfigParser() *Parser {
	return &Parser{
		hclParser: hclparse.NewParser(),
	}
}

// ParseHCL parses HCL content using proper HCL library
func (p *Parser) ParseHCL(content []byte) (interface{}, error) {
	// Parse HCL content using proper HCL library
	file, diags := p.hclParser.ParseHCL(content, "config.hcl")
	if diags.HasErrors() {
		return nil, fmt.Errorf("HCL parsing failed: %v", diags)
	}

	// Extract configuration from HCL body
	result := make(map[string]interface{})
	attrs, diags := file.Body.JustAttributes()
	if diags.HasErrors() {
		return nil, fmt.Errorf("failed to parse HCL attributes: %v", diags)
	}

	for name, attr := range attrs {
		value, err := extractHCLValue(attr.Expr)
		if err != nil {
			return nil, fmt.Errorf("failed to extract value for %s: %w", name, err)
		}
		result[name] = value
	}

	return result, nil
}

// extractHCLValue extracts a value from an HCL expression
func extractHCLValue(expr hcl.Expression) (interface{}, error) {
	val, diags := expr.Value(nil)
	if diags.HasErrors() {
		return nil, fmt.Errorf("failed to evaluate expression: %v", diags)
	}

	// Convert HCL value to Go interface{}
	if val.IsNull() {
		return nil, nil
	}

	switch {
	case val.Type().IsPrimitiveType() && val.Type().FriendlyName() == "string":
		return val.AsString(), nil
	case val.Type().IsPrimitiveType() && val.Type().FriendlyName() == "bool":
		return val.True(), nil
	case val.Type().IsPrimitiveType() && val.Type().FriendlyName() == "number":
		if val.CanIterateElements() {
			// Handle list length
			return val.LengthInt(), nil
		}
		// Try to get as float64
		if f := val.AsBigFloat(); f != nil {
			if f64, _ := f.Float64(); f64 != 0 {
				return f64, nil
			}
		}
		return val.AsBigFloat(), nil
	case val.Type().IsListType():
		var result []interface{}
		for _, item := range val.AsValueSlice() {
			if item.Type().IsPrimitiveType() && item.Type().FriendlyName() == "string" {
				result = append(result, item.AsString())
			} else if item.Type().IsPrimitiveType() && item.Type().FriendlyName() == "number" {
				if f := item.AsBigFloat(); f != nil {
					if f64, _ := f.Float64(); f64 != 0 {
						result = append(result, f64)
					}
				}
			} else {
				result = append(result, item)
			}
		}
		return result, nil
	default:
		return val, nil
	}
}

// ParseJSON parses JSON content using proper JSON parsing
func (p *Parser) ParseJSON(content []byte) (interface{}, error) {
	// Parse JSON content using proper JSON parsing
	var result interface{}
	if err := json.Unmarshal(content, &result); err != nil {
		return nil, fmt.Errorf("JSON parsing failed: %w", err)
	}

	return result, nil
}

// ValidateFormat validates the format of content
func (p *Parser) ValidateFormat(content []byte, format string) error {
	switch format {
	case "hcl":
		_, diags := p.hclParser.ParseHCL(content, "validation.hcl")
		if diags.HasErrors() {
			return fmt.Errorf("HCL validation failed: %v", diags)
		}
		return nil
	case "json":
		_, diags := p.hclParser.ParseJSON(content, "validation.json")
		if diags.HasErrors() {
			return fmt.Errorf("JSON validation failed: %v", diags)
		}
		return nil
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
}

// LoadActionsConfig loads actions configuration from a project
func LoadActionsConfig(projectPath string) (*spookyconfigtypes.ActionsConfig, error) {
	actionsFile := filepath.Join(projectPath, "actions.hcl")

	// Check if actions file exists
	if _, err := os.Stat(actionsFile); os.IsNotExist(err) {
		return &spookyconfigtypes.ActionsConfig{
			Actions: []spookyconfigtypes.Action{},
		}, nil
	}

	// Read and parse the actions file
	content, err := os.ReadFile(actionsFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read actions file: %w", err)
	}

	// Parse HCL content using hclparse
	parser := hclparse.NewParser()
	file, diags := parser.ParseHCL(content, actionsFile)
	if diags.HasErrors() {
		return nil, fmt.Errorf("failed to parse actions file: %v", diags)
	}

	// Decode the parsed content into our struct
	var actionsWrapper spookyconfigtypes.ActionsWrapper
	if diags := gohcl.DecodeBody(file.Body, nil, &actionsWrapper); diags.HasErrors() {
		return nil, fmt.Errorf("failed to decode actions file: %v", diags)
	}

	// If no actions block found, return empty config
	if actionsWrapper.Actions == nil {
		return &spookyconfigtypes.ActionsConfig{
			Actions: []spookyconfigtypes.Action{},
		}, nil
	}

	return actionsWrapper.Actions, nil
}

// ParseInventoryConfig loads inventory configuration from a file
func ParseInventoryConfig(inventoryFile string) (*spookyconfigtypes.InventoryConfig, error) {
	// Check if inventory file exists
	if _, err := os.Stat(inventoryFile); os.IsNotExist(err) {
		return &spookyconfigtypes.InventoryConfig{
			Machines: []spookyconfigtypes.Machine{},
		}, nil
	}

	// Read and parse the inventory file
	content, err := os.ReadFile(inventoryFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read inventory file: %w", err)
	}

	// Parse HCL content using hclparse
	parser := hclparse.NewParser()
	file, diags := parser.ParseHCL(content, inventoryFile)
	if diags.HasErrors() {
		return nil, fmt.Errorf("failed to parse inventory file: %v", diags)
	}

	// Decode the parsed content into our struct
	var inventoryWrapper spookyconfigtypes.InventoryWrapper
	if diags := gohcl.DecodeBody(file.Body, nil, &inventoryWrapper); diags.HasErrors() {
		return nil, fmt.Errorf("failed to decode inventory file: %v", diags)
	}

	// If no inventory block found, return empty config
	if inventoryWrapper.Inventory == nil {
		return &spookyconfigtypes.InventoryConfig{
			Machines: []spookyconfigtypes.Machine{},
		}, nil
	}

	return inventoryWrapper.Inventory, nil
}

// ParseMachinesInventory loads machines inventory from a file
func ParseMachinesInventory(machinesFile string) (*spookyconfigtypes.InventoryConfig, error) {
	// Check if machines file exists
	if _, err := os.Stat(machinesFile); os.IsNotExist(err) {
		return &spookyconfigtypes.InventoryConfig{
			Machines: []spookyconfigtypes.Machine{},
		}, nil
	}

	// Read and parse the machines file
	content, err := os.ReadFile(machinesFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read machines file: %w", err)
	}

	// Parse HCL content using hclparse
	parser := hclparse.NewParser()
	file, diags := parser.ParseHCL(content, machinesFile)
	if diags.HasErrors() {
		return nil, fmt.Errorf("failed to parse machines file: %v", diags)
	}

	// Decode the parsed content into our struct
	var machinesWrapper spookyconfigtypes.MachinesWrapper
	if diags := gohcl.DecodeBody(file.Body, nil, &machinesWrapper); diags.HasErrors() {
		return nil, fmt.Errorf("failed to decode machines file: %v", diags)
	}

	// If no machines block found, return empty config
	if machinesWrapper.Machines == nil {
		return &spookyconfigtypes.InventoryConfig{
			Machines: []spookyconfigtypes.Machine{},
		}, nil
	}

	return machinesWrapper.Machines, nil
}
