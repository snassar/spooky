package commands

import (
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

// BenchmarkExtractVariablesBlock benchmarks the extractVariablesBlock function
func BenchmarkExtractVariablesBlock(b *testing.B) {
	hclContent := `
variables {
  variable "test_var1" {
    value = "test_value1"
    description = "Test variable 1"
  }
  variable "test_var2" {
    value = 42
    description = "Test variable 2"
  }
  variable "test_var3" {
    value = true
    description = "Test variable 3"
  }
  variable "test_var4" {
    encrypted_value = {
      data = "AGE1-ENCRYPTED-DATA"
      format = "armored"
    }
  }
}`

	file, diags := hclsyntax.ParseConfig([]byte(hclContent), "test.hcl", hcl.Pos{Line: 1, Column: 1})
	if diags.HasErrors() {
		b.Fatalf("Failed to parse test HCL: %v", diags)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := extractVariablesBlock(file)
		if err != nil {
			b.Fatalf("extractVariablesBlock failed: %v", err)
		}
	}
}

// BenchmarkExtractVariableBlocks benchmarks the extractVariableBlocks function
func BenchmarkExtractVariableBlocks(b *testing.B) {
	hclContent := `
variables {
  variable "test_var1" {
    value = "test_value1"
    description = "Test variable 1"
  }
  variable "test_var2" {
    value = 42
    description = "Test variable 2"
  }
  variable "test_var3" {
    value = true
    description = "Test variable 3"
  }
  variable "test_var4" {
    encrypted_value = {
      data = "AGE1-ENCRYPTED-DATA"
      format = "armored"
    }
  }
}`

	file, diags := hclsyntax.ParseConfig([]byte(hclContent), "test.hcl", hcl.Pos{Line: 1, Column: 1})
	if diags.HasErrors() {
		b.Fatalf("Failed to parse test HCL: %v", diags)
	}

	variablesBlock, err := extractVariablesBlock(file)
	if err != nil {
		b.Fatalf("Failed to extract variables block: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := extractVariableBlocks(variablesBlock)
		if err != nil {
			b.Fatalf("extractVariableBlocks failed: %v", err)
		}
	}
}

// BenchmarkProcessVariableBlocks benchmarks the processVariableBlocks function
func BenchmarkProcessVariableBlocks(b *testing.B) {
	hclContent := `
variables {
  variable "test_var1" {
    value = "test_value1"
    description = "Test variable 1"
  }
  variable "test_var2" {
    value = 42
    description = "Test variable 2"
  }
  variable "test_var3" {
    value = true
    description = "Test variable 3"
  }
  variable "test_var4" {
    encrypted_value = {
      data = "AGE1-ENCRYPTED-DATA"
      format = "armored"
    }
  }
}`

	file, diags := hclsyntax.ParseConfig([]byte(hclContent), "test.hcl", hcl.Pos{Line: 1, Column: 1})
	if diags.HasErrors() {
		b.Fatalf("Failed to parse test HCL: %v", diags)
	}

	variablesBlock, err := extractVariablesBlock(file)
	if err != nil {
		b.Fatalf("Failed to extract variables block: %v", err)
	}

	content, err := extractVariableBlocks(variablesBlock)
	if err != nil {
		b.Fatalf("Failed to extract variable blocks: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := processVariableBlocks(content)
		if err != nil {
			b.Fatalf("processVariableBlocks failed: %v", err)
		}
	}
}

// BenchmarkParseVariableAttributes benchmarks the parseVariableAttributes function
func BenchmarkParseVariableAttributes(b *testing.B) {
	hclContent := `
variables {
  variable "test_var" {
    value = "test_string"
    description = "Test description"
    encrypted = false
  }
}`

	file, diags := hclsyntax.ParseConfig([]byte(hclContent), "test.hcl", hcl.Pos{Line: 1, Column: 1})
	if diags.HasErrors() {
		b.Fatalf("Failed to parse test HCL: %v", diags)
	}

	variablesBlock, err := extractVariablesBlock(file)
	if err != nil {
		b.Fatalf("Failed to extract variables block: %v", err)
	}

	content, err := extractVariableBlocks(variablesBlock)
	if err != nil {
		b.Fatalf("Failed to extract variable blocks: %v", err)
	}

	attrContent, err := extractVariableAttributes(content.Blocks[0])
	if err != nil {
		b.Fatalf("Failed to extract variable attributes: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := parseVariableAttributes(attrContent, "test_var")
		if err != nil {
			b.Fatalf("parseVariableAttributes failed: %v", err)
		}
	}
}

// BenchmarkFullVariableProcessing benchmarks the complete variable processing pipeline
func BenchmarkFullVariableProcessing(b *testing.B) {
	hclContent := `
variables {
  variable "test_var1" {
    value = "test_value1"
    description = "Test variable 1"
  }
  variable "test_var2" {
    value = 42
    description = "Test variable 2"
  }
  variable "test_var3" {
    value = true
    description = "Test variable 3"
  }
  variable "test_var4" {
    encrypted_value = {
      data = "AGE1-ENCRYPTED-DATA"
      format = "armored"
    }
  }
}`

	file, diags := hclsyntax.ParseConfig([]byte(hclContent), "test.hcl", hcl.Pos{Line: 1, Column: 1})
	if diags.HasErrors() {
		b.Fatalf("Failed to parse test HCL: %v", diags)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Complete pipeline: extract -> process -> parse
		variablesBlock, err := extractVariablesBlock(file)
		if err != nil {
			b.Fatalf("extractVariablesBlock failed: %v", err)
		}

		content, err := extractVariableBlocks(variablesBlock)
		if err != nil {
			b.Fatalf("extractVariableBlocks failed: %v", err)
		}

		_, err = processVariableBlocks(content)
		if err != nil {
			b.Fatalf("processVariableBlocks failed: %v", err)
		}
	}
}
