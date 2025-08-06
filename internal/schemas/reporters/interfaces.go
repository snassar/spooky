package reporters

import (
	"spooky/internal/schemas/types"
)

// Reporter interface defines validation reporting strategies
type Reporter interface {
	Report(result *types.ValidationResult) error
	GetName() string
}
