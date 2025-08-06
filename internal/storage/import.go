package storage

import (
	"fmt"
	"io"

	spookyschemas "spooky/internal/schemas"
)

// ImportFromJSON imports data from JSON format into storage
func ImportFromJSON(storage ExportableStorage, r io.Reader) error {
	// Validate against schema first
	validator := spookyschemas.NewSchemaValidator()
	if err := validator.LoadSchema(spookyschemas.SchemaTypeFactsStructure); err != nil {
		return fmt.Errorf("failed to load facts schema: %w", err)
	}

	// Use the storage's own import method
	return storage.ImportFromJSON(r)
}

// ImportFromHCL imports data from HCL format into storage
func ImportFromHCL(storage ExportableStorage, r io.Reader) error {
	// Validate against schema first
	validator := spookyschemas.NewSchemaValidator()
	if err := validator.LoadSchema(spookyschemas.SchemaTypeFactsStructure); err != nil {
		return fmt.Errorf("failed to load facts schema: %w", err)
	}

	// Use the storage's own import method
	return storage.ImportFromHCL(r)
}

// ImportFromJSONWithValidation imports data with custom validation
func ImportFromJSONWithValidation(storage ExportableStorage, r io.Reader, schemaType spookyschemas.SchemaType) error {
	// Validate against specified schema
	validator := spookyschemas.NewSchemaValidator()
	if err := validator.LoadSchema(schemaType); err != nil {
		return fmt.Errorf("failed to load schema: %w", err)
	}

	// Use the storage's own import method
	return storage.ImportFromJSON(r)
}

// ImportFromHCLWithValidation imports data with custom validation
func ImportFromHCLWithValidation(storage ExportableStorage, r io.Reader, schemaType spookyschemas.SchemaType) error {
	// Validate against specified schema
	validator := spookyschemas.NewSchemaValidator()
	if err := validator.LoadSchema(schemaType); err != nil {
		return fmt.Errorf("failed to load schema: %w", err)
	}

	// Use the storage's own import method
	return storage.ImportFromHCL(r)
}
