package utilities

import (
	"strings"
	"testing"

	"spooky/internal/hcl"
	"spooky/internal/schemas"
)

func TestHCLGenerator_FactsV1(t *testing.T) {
	// Create a sample FactsV1 structure
	facts := &schemas.FactsV1{
		BasicFacts: &schemas.BasicFactsV1{
			SystemFacts: &schemas.SystemFactsV1{
				Facts: map[string]*schemas.FactV1{
					"hostname": {
						Value:       "server01.example.com",
						Type:        "string",
						Description: "System hostname",
					},
					"os": {
						Value:       "linux",
						Type:        "string",
						Description: "Operating system",
					},
				},
			},
			HardwareFacts: &schemas.HardwareFactsV1{
				Facts: map[string]*schemas.FactV1{
					"cpu_count": {
						Value:       8,
						Type:        "number",
						Description: "Number of CPU cores",
					},
				},
			},
		},
		EnhancedFacts: &schemas.EnhancedFactsV1{
			Facts: map[string]*schemas.FactV1{
				"memory_total": {
					Value:       17179869184,
					Type:        "number",
					Description: "Total memory in bytes",
				},
			},
		},
		CustomFacts: &schemas.CustomFactsV1{
			Facts: map[string]*schemas.FactV1{
				"secret_key": {
					EncryptedValue: &schemas.EncryptedValueV1{
						Data:        "AGE1-PL...",
						Format:      "armored",
						Algorithm:   "age",
						Version:     "v1",
						EncryptedAt: "2024-01-15T10:30:00Z",
					},
					Type:        "encrypted",
					Description: "Encrypted secret key",
					Sensitive:   true,
				},
			},
		},
	}

	// Generate HCL
	hcl, err := hcl.GenerateHCL(facts, "facts")
	if err != nil {
		t.Fatalf("Failed to generate HCL: %v", err)
	}

	// Check that the HCL contains the expected structure
	expectedBlocks := []string{
		"facts {",
		"basic_facts {",
		"system_facts {",
		"hardware_facts {",
		"enhanced_facts {",
		"custom_facts {",
	}

	for _, block := range expectedBlocks {
		if !strings.Contains(hcl, block) {
			t.Errorf("HCL output missing expected block: %s", block)
		}
	}

	// Check that it contains some expected fact values
	expectedFacts := []string{
		"hostname",
		"server01.example.com",
		"cpu_count",
		"memory_total",
		"secret_key",
	}

	for _, fact := range expectedFacts {
		if !strings.Contains(hcl, fact) {
			t.Errorf("HCL output missing expected fact: %s", fact)
		}
	}

	t.Logf("Generated HCL for FactsV1:\n%s", hcl)
}
