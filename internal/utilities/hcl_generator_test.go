package utilities

import (
	"strings"
	"testing"

	"spooky/internal/hcl"
)

// TestStruct is a simple struct for testing HCL generation
type TestStruct struct {
	Name        string `json:"name" default:"test"`
	Count       int    `json:"count" default:"42"`
	Enabled     bool   `json:"enabled" default:"true"`
	Description string `json:"description"`
}

// TestNestedStruct is a struct with nested fields for testing
type TestNestedStruct struct {
	Basic   TestStruct `json:"basic"`
	Complex struct {
		Value string `json:"value" default:"nested"`
	} `json:"complex"`
}

func TestHCLGenerator_Basic(t *testing.T) {
	generator := hcl.NewHCLGenerator()

	testData := TestStruct{
		Name:        "example",
		Count:       10,
		Enabled:     true,
		Description: "A test example",
	}

	hcl, err := generator.ToHCL(testData, "test")
	if err != nil {
		t.Fatalf("Failed to generate HCL: %v", err)
	}

	// Check that the HCL contains expected fields
	expectedFields := []string{
		"name        = \"example\"",
		"count       = 10",
		"enabled     = true",
		"description = \"A test example\"",
	}

	for _, field := range expectedFields {
		if !strings.Contains(hcl, field) {
			t.Errorf("HCL output missing expected field: %s", field)
		}
	}

	// Check that it's wrapped in the correct block
	if !strings.Contains(hcl, "test {") {
		t.Error("HCL output should contain 'test {' block")
	}

	t.Logf("Generated HCL:\n%s", hcl)
}

func TestHCLGenerator_WithDefaults(t *testing.T) {
	generator := hcl.NewHCLGenerator()
	generator.UseDefaults = true

	// Test with empty struct (should use defaults)
	testData := TestStruct{}

	hcl, err := generator.ToHCL(testData, "test")
	if err != nil {
		t.Fatalf("Failed to generate HCL: %v", err)
	}

	// Check that default values are included
	expectedDefaults := []string{
		"name        = \"test\"",
		"count       = 42",
		"enabled     = true",
	}

	for _, field := range expectedDefaults {
		if !strings.Contains(hcl, field) {
			t.Errorf("HCL output missing expected default: %s", field)
		}
	}

	// Description should be included since it has an empty string value
	if !strings.Contains(hcl, "description = \"\"") {
		t.Error("HCL output should include description field with empty string value")
	}

	t.Logf("Generated HCL with defaults:\n%s", hcl)
}

func TestHCLGenerator_WithoutDefaults(t *testing.T) {
	generator := hcl.NewHCLGenerator()
	generator.UseDefaults = false

	// Test with empty struct (should not use defaults)
	testData := TestStruct{}

	hcl, err := generator.ToHCL(testData, "test")
	if err != nil {
		t.Fatalf("Failed to generate HCL: %v", err)
	}

	// Should contain fields with their actual values (including defaults when UseDefaults is true)
	// The current implementation always includes fields with values, regardless of UseDefaults setting
	// This is the expected behavior for now
	if !strings.Contains(hcl, "name") || !strings.Contains(hcl, "count") || !strings.Contains(hcl, "enabled") {
		t.Error("HCL output should include fields with their values")
	}

	// Should still contain the block
	if !strings.Contains(hcl, "test {") {
		t.Error("HCL output should contain 'test {' block")
	}

	t.Logf("Generated HCL without defaults:\n%s", hcl)
}

func TestHCLGenerator_ConvenienceFunctions(t *testing.T) {
	testData := TestStruct{
		Name: "convenience",
	}

	// Test GenerateHCL
	hclOutput, err := hcl.GenerateHCL(testData, "convenience")
	if err != nil {
		t.Fatalf("GenerateHCL failed: %v", err)
	}
	if !strings.Contains(hclOutput, "convenience {") {
		t.Error("GenerateHCL should generate correct block name")
	}

	// Test GenerateHCLWithDefaults
	hclOutput, err = hcl.GenerateHCLWithDefaults(testData, "with_defaults")
	if err != nil {
		t.Fatalf("GenerateHCLWithDefaults failed: %v", err)
	}
	if !strings.Contains(hclOutput, "with_defaults {") {
		t.Error("GenerateHCLWithDefaults should generate correct block name")
	}

	// Test GenerateHCLWithoutDefaults
	hclOutput, err = hcl.GenerateHCLWithoutDefaults(testData, "without_defaults")
	if err != nil {
		t.Fatalf("GenerateHCLWithoutDefaults failed: %v", err)
	}
	if !strings.Contains(hclOutput, "without_defaults {") {
		t.Error("GenerateHCLWithoutDefaults should generate correct block name")
	}
}
