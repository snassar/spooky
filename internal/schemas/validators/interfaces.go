package validators

import (
	"spooky/internal/schemas/types"
)

// Validator interface defines validation strategies
type Validator interface {
	Validate(data interface{}, schema *types.Schema) *types.ValidationResult
	GetName() string
}
