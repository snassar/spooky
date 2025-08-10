package storage

import (
	"fmt"

	spookyfactstypes "spooky/internal/types/facts"
	spookyschemas "spooky/internal/schemas"
)

// ValidateFactCollection validates a fact collection against the schema
func ValidateFactCollection(collection *spookyfactstypes.FactCollection) error {
	validator := spookyschemas.NewSchemaValidator()
	if err := validator.LoadSchema(spookyschemas.SchemaTypeFactsStructure); err != nil {
		return fmt.Errorf("failed to load facts schema: %w", err)
	}

	return validator.ValidateData(collection, string(spookyschemas.SchemaTypeFactsStructure))
}

// ValidateFactCollections validates multiple fact collections
func ValidateFactCollections(collections []*spookyfactstypes.FactCollection) error {
	validator := spookyschemas.NewSchemaValidator()
	if err := validator.LoadSchema(spookyschemas.SchemaTypeFactsStructure); err != nil {
		return fmt.Errorf("failed to load facts schema: %w", err)
	}

	for i, collection := range collections {
		if err := validator.ValidateData(collection, string(spookyschemas.SchemaTypeFactsStructure)); err != nil {
			return fmt.Errorf("failed to validate fact collection %d: %w", i, err)
		}
	}

	return nil
}

// ValidateData validates data against a specific schema
func ValidateData(data interface{}, schemaType spookyschemas.SchemaType) error {
	validator := spookyschemas.NewSchemaValidator()
	if err := validator.LoadSchema(schemaType); err != nil {
		return fmt.Errorf("failed to load schema: %w", err)
	}

	return validator.ValidateData(data, string(schemaType))
}
