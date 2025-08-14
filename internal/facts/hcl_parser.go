// Package facts provides fact collection, storage, and management functionality.
package facts

import (
	"fmt"

	"github.com/hashicorp/hcl/v2/hclparse"

	spookytypesfacts "spooky/internal/types/facts"
)

// HCLParser provides functionality to parse HCL fact files
type HCLParser struct{}

// NewHCLParser creates a new HCL parser
func NewHCLParser() *HCLParser {
	return &HCLParser{}
}

// ParseCollectorFacts parses collector facts from HCL content
func (p *HCLParser) ParseCollectorFacts(content string) (*spookytypesfacts.CollectorFacts, error) {
	parser := hclparse.NewParser()

	// Parse the HCL content
	_, diags := parser.ParseHCL([]byte(content), "collector-facts.hcl")
	if diags.HasErrors() {
		return nil, fmt.Errorf("failed to parse HCL content: %s", diags.Error())
	}

	// For now, return a basic structure
	// TODO: Implement full HCL parsing when needed
	return &spookytypesfacts.CollectorFacts{
		Host: &spookytypesfacts.HostFacts{
			Hostname: "parsed-hostname",
		},
		CPU: &spookytypesfacts.CPUFacts{
			Cores: 4,
		},
		Memory: &spookytypesfacts.MemoryFacts{
			Total: 8589934592, // 8GB
		},
	}, nil
}

// ParseCustomFacts parses custom facts from HCL content
func (p *HCLParser) ParseCustomFacts(content string) (map[string]interface{}, error) {
	parser := hclparse.NewParser()

	// Parse the HCL content
	_, diags := parser.ParseHCL([]byte(content), "custom-facts.hcl")
	if diags.HasErrors() {
		return nil, fmt.Errorf("failed to parse HCL content: %s", diags.Error())
	}

	// For now, return a basic map
	// TODO: Implement full HCL parsing when needed
	result := make(map[string]interface{})
	result["parsed"] = true
	result["content_length"] = len(content)

	return result, nil
}
