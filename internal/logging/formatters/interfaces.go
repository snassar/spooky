package formatters

import (
	"spooky/internal/logging/types"
)

// Formatter interface defines log formatting strategies
type Formatter interface {
	Format(entry *types.LogEntry) ([]byte, error)
	GetName() string
}
