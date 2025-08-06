package storage

import (
	"fmt"
	"io"

	spookyschemas "spooky/internal/schemas"
)

// ExportToJSON exports data from storage to JSON format
func ExportToJSON(storage ExportableStorage, w io.Writer) error {
	// Validate against schema first
	validator := spookyschemas.NewSchemaValidator()
	if err := validator.LoadSchema(spookyschemas.SchemaTypeFactsStructure); err != nil {
		return fmt.Errorf("failed to load facts schema: %w", err)
	}

	// Use the storage's own export method
	return storage.ExportToJSON(w)
}

// ExportToHCL exports data from storage to HCL format
func ExportToHCL(storage ExportableStorage, w io.Writer) error {
	// Validate against schema first
	validator := spookyschemas.NewSchemaValidator()
	if err := validator.LoadSchema(spookyschemas.SchemaTypeFactsStructure); err != nil {
		return fmt.Errorf("failed to load facts schema: %w", err)
	}

	// Use the storage's own export method
	return storage.ExportToHCL(w)
}

// ExportToJSONWithValidation exports data with custom validation
func ExportToJSONWithValidation(storage ExportableStorage, w io.Writer, schemaType spookyschemas.SchemaType) error {
	// Validate against specified schema
	validator := spookyschemas.NewSchemaValidator()
	if err := validator.LoadSchema(schemaType); err != nil {
		return fmt.Errorf("failed to load schema: %w", err)
	}

	// Use the storage's own export method
	return storage.ExportToJSON(w)
}

// ExportToHCLWithValidation exports data with custom validation
func ExportToHCLWithValidation(storage ExportableStorage, w io.Writer, schemaType spookyschemas.SchemaType) error {
	// Validate against specified schema
	validator := spookyschemas.NewSchemaValidator()
	if err := validator.LoadSchema(schemaType); err != nil {
		return fmt.Errorf("failed to load schema: %w", err)
	}

	// Use the storage's own export method
	return storage.ExportToHCL(w)
}
