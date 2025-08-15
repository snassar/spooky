package cmd

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"

	spookylogging "spooky/internal/logging"
	spookyschemas "spooky/internal/schemas"
)

// schemasCmd represents the schemas command
var schemasCmd = &cobra.Command{
	Use:   "schemas",
	Short: "Manage and validate schemas",
	Long: `Manage and validate HCL schemas for the spooky codebase.

This command provides schema management functionality including:
- Schema validation with detailed error reporting
- Schema listing and discovery`,
}

// schemasValidateCmd represents the schemas validate command
var schemasValidateCmd = &cobra.Command{
	Use:   "validate [schema-file] [data-file]",
	Short: "Validate data against a schema",
	Long: `Validate HCL data files against their corresponding schemas.

This command provides comprehensive validation with detailed error reporting,
including field-level validation, cross-field validation, and custom validation rules.`,
	Args: cobra.ExactArgs(2),
	RunE: func(_ *cobra.Command, args []string) error {
		schemaFile := args[0]
		dataFile := args[1]

		return handleSchemasValidate(schemaFile, dataFile)
	},
}

// schemasListCmd represents the schemas list command
var schemasListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available schemas",
	Long: `List all available schemas in the system.

This command shows:
- All registered schemas
- Schema types and versions
- Schema statistics and metadata`,
	RunE: func(_ *cobra.Command, _ []string) error {
		return handleSchemasList()
	},
}

// handleSchemasValidate handles schema validation
func handleSchemasValidate(schemaFile, dataFile string) error {
	ctx := context.Background()

	// Initialize logging
	logManager := spookylogging.NewLogManager()
	logger := logManager.GetLogger("schemas")

	// Create schema manager
	manager := spookyschemas.NewManager(logger)

	// Load schema
	schema, err := manager.Load(schemaFile)
	if err != nil {
		return fmt.Errorf("failed to load schema: %w", err)
	}

	// Track schema evolution
	if err := manager.TrackSchemaEvolution(schema); err != nil {
		logger.Warn("Failed to track schema evolution", map[string]interface{}{
			"error": err.Error(),
		})
	}

	// Validate data file against schema with enhanced features
	result, err := manager.ValidateWithEnhancedFeatures(ctx, schema, dataFile)
	if err != nil {
		return fmt.Errorf("failed to validate data: %w", err)
	}

	// Report validation results
	fmt.Printf("🔍 Validating %s against schema %s\n", dataFile, schemaFile)

	if result.Valid {
		fmt.Printf("✅ Validation passed\n")
		fmt.Printf("📊 Statistics:\n")
		fmt.Printf("   - Total fields: %d\n", result.Statistics.TotalFields)
		fmt.Printf("   - Valid fields: %d\n", result.Statistics.ValidFields)
		fmt.Printf("   - Rules executed: %d\n", result.Statistics.RulesExecuted)
		fmt.Printf("   - Duration: %v\n", result.Statistics.Duration)
	} else {
		fmt.Printf("❌ Validation failed\n")
		fmt.Printf("📊 Statistics:\n")
		fmt.Printf("   - Total fields: %d\n", result.Statistics.TotalFields)
		fmt.Printf("   - Invalid fields: %d\n", result.Statistics.InvalidFields)
		fmt.Printf("   - Rules failed: %d\n", result.Statistics.RulesFailed)
		fmt.Printf("   - Duration: %v\n", result.Statistics.Duration)

		// Report errors
		if len(result.Errors) > 0 {
			fmt.Printf("\n❌ Errors:\n")
			for i := range result.Errors {
				err := &result.Errors[i]
				fmt.Printf("   %d. %s\n", i+1, err.Message)
				if err.FieldPath != "" {
					fmt.Printf("      Field: %s\n", err.FieldPath)
				}
				if err.Location != nil {
					fmt.Printf("      Location: %s:%d:%d\n",
						err.Location.FilePath, err.Location.Line, err.Location.Column)
				}
				if len(err.Suggestions) > 0 {
					fmt.Printf("      Suggestions:\n")
					for _, suggestion := range err.Suggestions {
						fmt.Printf("        - %s\n", suggestion)
					}
				}
				fmt.Printf("\n")
			}
		}

		// Report warnings
		if len(result.Warnings) > 0 {
			fmt.Printf("⚠️  Warnings:\n")
			for i := range result.Warnings {
				warning := &result.Warnings[i]
				fmt.Printf("   %d. %s\n", i+1, warning.Message)
			}
		}
	}

	// Report recommendations
	if len(result.Recommendations) > 0 {
		fmt.Printf("💡 Recommendations:\n")
		for _, recommendation := range result.Recommendations {
			fmt.Printf("   - %s\n", recommendation)
		}
	}

	return nil
}

// handleSchemasList handles schema listing
func handleSchemasList() error {
	// Initialize logging
	logManager := spookylogging.NewLogManager()
	logger := logManager.GetLogger("schemas")

	// Create schema manager
	manager := spookyschemas.NewManager(logger)

	// Load schemas from the schemas directory
	schemasDir := filepath.Join("internal", "schemas", "schemas")
	if err := manager.LoadSchemas(schemasDir); err != nil {
		return fmt.Errorf("failed to load schemas: %w", err)
	}

	// Get schema statistics
	stats := manager.GetSchemaStatistics()

	// Report schema information
	fmt.Printf("📚 Available Schemas\n\n")
	fmt.Printf("📊 Statistics:\n")
	fmt.Printf("   - Total schemas: %d\n", stats["total_schemas"])

	if byType, ok := stats["by_type"].(map[string]int); ok {
		fmt.Printf("   - By type:\n")
		for schemaType, count := range byType {
			fmt.Printf("     - %s: %d\n", schemaType, count)
		}
	}

	if byVersion, ok := stats["by_version"].(map[string]int); ok {
		fmt.Printf("   - By version:\n")
		for version, count := range byVersion {
			fmt.Printf("     - %s: %d\n", version, count)
		}
	}

	// List all schemas
	schemas := manager.List()
	if len(schemas) > 0 {
		fmt.Printf("\n📋 Schemas:\n")
		for _, schema := range schemas {
			fmt.Printf("   - %s (%s, v%s)\n", schema.Name, schema.Type, schema.Version)
			if schema.Description != "" {
				fmt.Printf("     %s\n", schema.Description)
			}
		}
	}

	return nil
}

func init() {
	// Add subcommands to schemas command
	schemasCmd.AddCommand(schemasValidateCmd)
	schemasCmd.AddCommand(schemasListCmd)

	// Add schemas command to root
	RootCmd.AddCommand(schemasCmd)
}
