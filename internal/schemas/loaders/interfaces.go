package loaders

import (
	"spooky/internal/schemas/types"
)

// Loader interface defines schema loading strategies
type Loader interface {
	LoadSchema(schemaType types.SchemaType) (*types.Schema, error)
	LoadAllSchemas() (map[types.SchemaType]*types.Schema, error)
	GetName() string
}
