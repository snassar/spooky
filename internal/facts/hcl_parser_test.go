package facts

import (
	"testing"

	"github.com/hashicorp/hcl/v2/hclsyntax"
)

func TestHCLParser_parseCPUBlock(t *testing.T) {
	parser := NewHCLParser()

	// Test that the parser can be created
	if parser == nil {
		t.Fatal("NewHCLParser returned nil")
	}
}

func TestHCLParser_parseMemoryBlock(t *testing.T) {
	parser := NewHCLParser()

	// Test that the parser can be created
	if parser == nil {
		t.Fatal("NewHCLParser returned nil")
	}
}

func TestCreateParser(t *testing.T) {
	// Test that the createParser function works correctly
	var testValue int
	parser := createParser(&testValue, func(_ hclsyntax.Expression) (int, error) {
		return 42, nil
	})

	// This is a basic test to ensure the function can be created
	if parser == nil {
		t.Fatal("createParser returned nil")
	}
}
