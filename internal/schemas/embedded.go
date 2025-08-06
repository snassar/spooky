//go:build !test
// +build !test

// Package schemas provides embedded schema composition for spooky
package schemas

import (
	"embed"
	"fmt"
)

//go:embed schemas/facts-structure.schema.hcl
var factsStructureFS embed.FS

//go:embed schemas/facts-storage.schema.hcl
var factsStorageFS embed.FS

//go:embed schemas/actions.hcl
var actionsFS embed.FS

//go:embed schemas/machines.hcl
var machinesFS embed.FS

//go:embed schemas/variables-structure.hcl
var variablesStructureFS embed.FS

//go:embed schemas/variables-hcl.hcl
var variablesHCLFS embed.FS

//go:embed schemas/variables-json.hcl
var variablesJSONFS embed.FS

//go:embed schemas/project.hcl
var projectFS embed.FS

//go:embed schemas/spooky.hcl
var spookyFS embed.FS

//go:embed schemas/template-structure.hcl
var templateStructureFS embed.FS

//go:embed schemas/template-metadata.hcl
var templateMetadataFS embed.FS

//go:embed schemas/template-context.hcl
var templateContextFS embed.FS

//go:embed schemas/template-functions.hcl
var templateFunctionsFS embed.FS

//go:embed schemas/custom-facts-hcl.hcl
var customFactsHCLFS embed.FS

// SchemaType constants for different schema types
const (
	SchemaTypeProject            SchemaType = "project"
	SchemaTypeMachines           SchemaType = "machines"
	SchemaTypeActions            SchemaType = "actions"
	SchemaTypeSpooky             SchemaType = "spooky"
	SchemaTypeVariables          SchemaType = "variables"
	SchemaTypeVariablesStructure SchemaType = "variables-structure"
	SchemaTypeVariablesHCL       SchemaType = "variables-hcl"
	SchemaTypeVariablesJSON      SchemaType = "variables-json"
	SchemaTypeTemplateStructure  SchemaType = "template-structure"
	SchemaTypeTemplateMetadata   SchemaType = "template-metadata"
	SchemaTypeTemplateContext    SchemaType = "template-context"
	SchemaTypeTemplateFunctions  SchemaType = "template-functions"
	SchemaTypeCustomFactsHCL     SchemaType = "custom-facts-hcl"
	SchemaTypeFactsStructure     SchemaType = "facts-structure"
	SchemaTypeFactsStorage       SchemaType = "facts-storage"
	// Enhanced composed schema types
	SchemaTypeActionsComposed           SchemaType = "actions-composed"
	SchemaTypeMachinesComposed          SchemaType = "machines-composed"
	SchemaTypeVariablesComposed         SchemaType = "variables-composed"
	SchemaTypeProjectComposed           SchemaType = "project-composed"
	SchemaTypeTemplateMetadataComposed  SchemaType = "template-metadata-composed"
	SchemaTypeTemplateContextComposed   SchemaType = "template-context-composed"
	SchemaTypeTemplateFunctionsComposed SchemaType = "template-functions-composed"
)

// GetSchema returns the embedded schema for the given type
func GetSchema(schemaType SchemaType) ([]byte, error) {
	switch schemaType {
	case SchemaTypeProject:
		return projectFS.ReadFile("schemas/project.hcl")
	case SchemaTypeMachines:
		return machinesFS.ReadFile("schemas/machines.hcl")
	case SchemaTypeActions:
		return actionsFS.ReadFile("schemas/actions.hcl")
	case SchemaTypeSpooky:
		return spookyFS.ReadFile("schemas/spooky.hcl")
	case SchemaTypeTemplateStructure:
		return templateStructureFS.ReadFile("schemas/template-structure.hcl")
	case SchemaTypeTemplateMetadata:
		return templateMetadataFS.ReadFile("schemas/template-metadata.hcl")
	case SchemaTypeTemplateContext:
		return templateContextFS.ReadFile("schemas/template-context.hcl")
	case SchemaTypeTemplateFunctions:
		return templateFunctionsFS.ReadFile("schemas/template-functions.hcl")
	case SchemaTypeVariablesStructure:
		return variablesStructureFS.ReadFile("schemas/variables-structure.hcl")
	case SchemaTypeVariablesHCL:
		return variablesHCLFS.ReadFile("schemas/variables-hcl.hcl")
	case SchemaTypeVariablesJSON:
		return variablesJSONFS.ReadFile("schemas/variables-json.hcl")
	case SchemaTypeCustomFactsHCL:
		return customFactsHCLFS.ReadFile("schemas/custom-facts-hcl.hcl")
	case SchemaTypeFactsStructure:
		return factsStructureFS.ReadFile("schemas/facts-structure.schema.hcl")
	case SchemaTypeFactsStorage:
		return factsStorageFS.ReadFile("schemas/facts-storage.schema.hcl")
	default:
		return []byte(""), fmt.Errorf("schema not found: %s", schemaType)
	}
}

// GetComposedSchema returns a composed schema by name
func GetComposedSchema(_ string) (string, error) {
	return "", fmt.Errorf("composed schemas not supported")
}

// ListComposedSchemas returns a list of available composed schemas
func ListComposedSchemas() ([]string, error) {
	return []string{}, nil
}
