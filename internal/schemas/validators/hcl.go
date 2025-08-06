package validators

import (
	"fmt"
	"spooky/internal/schemas/types"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
)

// HCLValidator implements HCL validation
type HCLValidator struct {
	parser *hclparse.Parser
}

// NewHCLValidator creates a new HCL validator
func NewHCLValidator() *HCLValidator {
	return &HCLValidator{
		parser: hclparse.NewParser(),
	}
}

// Validate validates HCL data against a schema
func (v *HCLValidator) Validate(data interface{}, schema *types.Schema) *types.ValidationResult {
	result := &types.ValidationResult{
		Valid:  true,
		Schema: string(schema.Type),
	}

	// Parse the HCL content
	file, diags := v.parser.ParseHCL([]byte(schema.Content), schema.Filename)
	if diags.HasErrors() {
		result.Valid = false
		for _, diag := range diags {
			result.Errors = append(result.Errors, types.ValidationError{
				File:     schema.Filename,
				Line:     diag.Subject.Start.Line,
				Column:   diag.Subject.Start.Column,
				Message:  diag.Summary,
				Severity: "error",
			})
		}
		return result
	}

	// Basic HCL validation
	if err := v.validateBasicHCL(file); err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, types.ValidationError{
			File:     schema.Filename,
			Message:  err.Error(),
			Severity: "error",
		})
	}

	return result
}

// GetName returns the validator name
func (v *HCLValidator) GetName() string {
	return "hcl"
}

// validateBasicHCL performs basic HCL validation
func (v *HCLValidator) validateBasicHCL(file *hcl.File) error {
	// Basic validation logic would go here
	// For now, just check if the file is not nil
	if file == nil {
		return fmt.Errorf("nil HCL file")
	}
	return nil
}
