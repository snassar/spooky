# Spooky Schema Embedder

A Go program that embeds HCL schemas and validation rules directly into the binary using Go's `embed` directive.

## Features

- **Schema Embedding**: Embeds all HCL schema files into the binary
- **Validation Rules**: Includes comprehensive validation rules for each schema
- **Metadata Support**: Embeds schema metadata for documentation
- **Easy Access**: Simple API to retrieve schemas and rules by name
- **Categorized Access**: Separate access to structure schemas, validation rules, and metadata

## Structure

```
internal/schemas/
├── schema_embedder.go     # Main embedder implementation
└── schemafiles/          # Embedded HCL files
    ├── structure/        # Core schema definitions
    ├── validation/       # Validation rules
    └── metadata/         # Schema metadata
```

## Embedded Schemas

### Structure Schemas (12 files)
- `project.hcl` - Project configuration structure
- `machines.hcl` - Machine inventory and SSH connectivity
- `actions.hcl` - Action definitions and execution
- `facts.hcl` - Fact collection and storage
- `variables.hcl` - Variable definitions and resolution
- `logging.hcl` - Logging configuration
- `templates.hcl` - Template definitions
- `template-context.hcl` - Template context and variables
- `template-functions.hcl` - Template functions and helpers
- `template-metadata.hcl` - Template metadata and documentation
- `spooky.hcl` - Global spooky configuration
- `project-directory.hcl` - Project directory structure

### Validation Rules (9 files)
Each structure schema has corresponding validation rules:
- `project-rules.hcl`
- `machines-rules.hcl`
- `actions-rules.hcl`
- `facts-rules.hcl`
- `variables-rules.hcl`
- `templates-rules.hcl`
- `logging-rules.hcl`
- `project-directory-rules.hcl`
- `template-metadata-rules.hcl`

### Metadata (1 file)
- `schema-metadata.hcl` - Meta-schema for schema validation

## Usage

```go
package main

import (
    "fmt"
    "log"
    "spooky/internal/schemas"
)

func main() {
    // Initialize the schema embedder
    embedder, err := schemas.NewSchemaEmbedder()
    if err != nil {
        log.Fatal(err)
    }

    // Get a specific schema
    if schema, exists := embedder.GetSchema("project"); exists {
        fmt.Printf("Project schema: %d bytes\n", len(schema))
    }

    // Get validation rules
    if rules, exists := embedder.GetValidationRules("project"); exists {
        fmt.Printf("Project rules: %d bytes\n", len(rules))
    }

    // List all available schemas
    schemas := embedder.ListSchemas()
    for _, name := range schemas {
        fmt.Printf("- %s\n", name)
    }

    // Print summary
    embedder.PrintSchemaSummary()
}
```

## Running

```bash
# Build and run
go run main.go

# Build binary
go build -o spooky-schemas main.go

# Run binary
./spooky-schemas
```

## Output Example

```
🔮 Spooky Schema Embedder
==========================
=== Embedded Schemas Summary ===

📁 Structure Schemas (12):
  - machines (26085 bytes)
  - project (10496 bytes)
  - actions (14876 bytes)
  ...

✅ Validation Rules (9):
  - project (4701 bytes)
  - actions (6448 bytes)
  - machines (4675 bytes)
  ...

📋 Metadata (1):
  - schema-metadata (10293 bytes)

📊 Total: 22 files embedded
```

## Benefits

1. **Self-contained**: All schemas are embedded in the binary
2. **No external dependencies**: No need to ship schema files separately
3. **Version consistency**: Schemas are always in sync with the binary
4. **Fast access**: Schemas are loaded in memory at startup
5. **Type safety**: Go's embed directive ensures compile-time validation

## Development

The schema embedder is designed to be easily extensible:

- Add new schemas to `internal/schemas/schemafiles/structure/`
- Add validation rules to `internal/schemas/schemafiles/validation/`
- Add metadata to `internal/schemas/schemafiles/metadata/`
- The embedder will automatically pick up new files

## Dependencies

- Go 1.24+ (for embed directive support)
- `github.com/hashicorp/hcl/v2` - HCL parsing
- `github.com/zclconf/go-cty` - Configuration type system
