package cli

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"spooky/internal/facts"
	"spooky/internal/logging"

	"github.com/spf13/cobra"
)

var (
	FactsCmd = &cobra.Command{
		Use:   "facts [PROJECT_PATH]",
		Short: "Manage server facts and fact collection within a spooky project",
		Long: `Manage server facts and fact collection from multiple sources within a spooky project.

Facts are machine-specific information that can be collected from various sources
including SSH connections, local system information, and external data sources.

This command must be run from within a spooky project directory or with a project path specified.

Examples:
  # Gather facts from all machines in the current project
  spooky facts gather

  # Gather facts from all machines in a specific project
  spooky facts /path/to/project gather

  # List all collected facts in the current project
  spooky facts list

  # Query facts for specific information in the current project
  spooky facts query "os.name == 'linux'"

  # Export facts to JSON format from the current project
  spooky facts export --output facts.json

  # Import facts from external source into the current project
  spooky facts import external-facts.json`,
		Args: cobra.MaximumNArgs(1),
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// Set the project path for all facts subcommands
			projectPath := "."
			if len(args) > 0 {
				projectPath = args[0]
			}
			cmd.SetContext(context.WithValue(cmd.Context(), projectPathKey{}, projectPath))
			return nil
		},
	}

	factsGatherCmd = &cobra.Command{
		Use:   "gather [hosts]",
		Short: "Gather facts from target servers",
		Long:  `Gather facts from target servers or use inventory if hosts not specified.`,
		Args:  cobra.MaximumNArgs(1),
		RunE:  runFactsGather,
	}

	factsImportCmd = &cobra.Command{
		Use:   "import <source>",
		Short: "Import facts from external sources",
		Long:  `Import facts from local JSON file or HTTPS URL (HTTPS required for security).`,
		Args:  cobra.ExactArgs(1),
		RunE:  runFactsImport,
	}

	factsExportCmd = &cobra.Command{
		Use:   "export",
		Short: "Export facts to various formats",
		Long:  `Export facts to JSON, YAML, CSV, or table format.`,
		RunE:  runFactsExport,
	}

	factsValidateCmd = &cobra.Command{
		Use:   "validate",
		Short: "Validate facts against rules and schemas",
		Long:  `Validate facts against validation rules and schemas.`,
		RunE:  runFactsValidate,
	}

	factsQueryCmd = &cobra.Command{
		Use:   "query <expression>",
		Short: "Query facts using expressions and filters",
		Long:  `Query facts using expressions and filters to find specific information.`,
		Args:  cobra.ExactArgs(1),
		RunE:  runFactsQuery,
	}

	factsCacheCmd = &cobra.Command{
		Use:   "cache",
		Short: "Manage fact cache",
		Long:  `Manage the fact collection cache.`,
		RunE:  runFactsCache,
	}

	factsCacheClearCmd = &cobra.Command{
		Use:   "clear",
		Short: "Clear the fact cache",
		Long:  `Clear all cached facts.`,
		RunE:  runFactsCacheClear,
	}

	factsCacheExpiredCmd = &cobra.Command{
		Use:   "clear-expired",
		Short: "Clear expired facts from cache",
		Long:  `Remove expired facts from the cache.`,
		RunE:  runFactsCacheExpired,
	}

	// Flags
	factsSpecificKeys []string
	factsInventory    string
	factsConfig       string
	factsParallel     int
	factsTimeout      int
	factsUpdate       bool
	factsCacheDir     string
	factsMerge        bool
	factsValidate     bool
	factsFormat       string
	factsMapping      string
	factsOutput       string
	factsFilter       string
	factsFields       string
	factsPretty       bool
	factsRules        string
	factsSchema       string
	factsStrict       bool
	factsLimit        int

	// Enhanced import flags
	importSource      string
	importFile        string
	importURL         string
	importMergeMode   string
	importSelectFacts []string
	importOverride    bool
	importDryRun      bool
	importServer      string
)

// initFactsCommands initializes the facts command and its subcommands
func initFactsCommands() {
	// Add new subcommands to facts
	FactsCmd.AddCommand(factsGatherCmd)
	FactsCmd.AddCommand(factsImportCmd)
	FactsCmd.AddCommand(factsExportCmd)
	FactsCmd.AddCommand(factsValidateCmd)
	FactsCmd.AddCommand(factsQueryCmd)

	// Add cache subcommand
	FactsCmd.AddCommand(factsCacheCmd)

	// Add subcommands to cache
	factsCacheCmd.AddCommand(factsCacheClearCmd)
	factsCacheCmd.AddCommand(factsCacheExpiredCmd)

	// Add flags for new commands
	factsGatherCmd.Flags().StringVar(&factsInventory, "inventory", "", "Path to inventory file")
	// Removed legacy config file support - use project-based approach instead
	factsGatherCmd.Flags().IntVar(&factsParallel, "parallel", 10, "Number of parallel fact gathering")
	factsGatherCmd.Flags().IntVar(&factsTimeout, "timeout", 60, "Timeout per host in seconds")
	factsGatherCmd.Flags().StringSliceVarP(&factsSpecificKeys, "facts", "f", nil, "Comma-separated list of fact types to gather")
	factsGatherCmd.Flags().BoolVar(&factsUpdate, "update", false, "Update existing facts instead of replacing")
	factsGatherCmd.Flags().StringVar(&factsCacheDir, "cache-dir", "", "Directory for fact caching")

	factsImportCmd.Flags().BoolVar(&factsMerge, "merge", false, "Merge with existing facts instead of replacing")
	factsImportCmd.Flags().BoolVar(&factsValidate, "validate", false, "Validate facts before importing")
	factsImportCmd.Flags().StringVar(&factsFormat, "format", "", "Source format: json, yaml, csv (default: auto-detect)")
	factsImportCmd.Flags().StringVar(&factsMapping, "mapping", "", "Path to field mapping configuration")

	// Enhanced import flags
	factsImportCmd.Flags().StringVar(&importSource, "source", "local", "Import source: local, http")
	factsImportCmd.Flags().StringVar(&importFile, "file", "", "Path to local JSON file")
	factsImportCmd.Flags().StringVar(&importURL, "url", "", "HTTPS URL for remote import (HTTPS required for security)")
	factsImportCmd.Flags().StringVar(&importMergeMode, "merge-mode", "replace", "Merge mode: replace, merge, append, select")
	factsImportCmd.Flags().StringSliceVar(&importSelectFacts, "select-facts", nil, "Comma-separated list of facts to import")
	factsImportCmd.Flags().BoolVar(&importOverride, "override", false, "Allow fact overrides")
	factsImportCmd.Flags().BoolVar(&importDryRun, "dry-run", false, "Show what would be imported without importing")
	factsImportCmd.Flags().StringVar(&importServer, "server", "", "Specific server to import facts for")

	factsExportCmd.Flags().StringVar(&factsFormat, "format", "json", "Output format: json, yaml, csv, table")
	factsExportCmd.Flags().StringVar(&factsOutput, "output", "", "Output file path (default: stdout)")
	factsExportCmd.Flags().StringVar(&factsFilter, "filter", "", "Filter facts by expression")
	factsExportCmd.Flags().StringVar(&factsFields, "fields", "", "Comma-separated list of fields to include")
	factsExportCmd.Flags().BoolVar(&factsPretty, "pretty", false, "Pretty-print JSON output")

	factsValidateCmd.Flags().StringVar(&factsRules, "rules", "", "Path to validation rules file")
	factsValidateCmd.Flags().StringVar(&factsSchema, "schema", "", "Path to schema file")
	factsValidateCmd.Flags().BoolVar(&factsStrict, "strict", false, "Enable strict validation mode")
	factsValidateCmd.Flags().StringVar(&factsFormat, "format", "text", "Output format: text, json, html")
	factsValidateCmd.Flags().StringVar(&factsOutput, "output", "", "Output file path (default: stdout)")

	factsQueryCmd.Flags().StringVar(&factsFormat, "format", "table", "Output format: table, json, yaml")
	factsQueryCmd.Flags().StringVar(&factsOutput, "output", "", "Output file path (default: stdout)")
	factsQueryCmd.Flags().StringVar(&factsFields, "fields", "", "Comma-separated list of fields to include")
	factsQueryCmd.Flags().IntVar(&factsLimit, "limit", 0, "Limit number of results")
	factsQueryCmd.Flags().BoolVar(&factsPretty, "pretty", false, "Pretty-print JSON output")

}

// displayMachineFacts displays facts for a single machine
func displayMachineFacts(machineFacts *facts.MachineFacts) {
	fmt.Printf("Machine: %s\n", machineFacts.MachineName)
	fmt.Printf("  System ID: %s\n", machineFacts.SystemID)
	fmt.Printf("  Action File: %s\n", machineFacts.ActionFile)
	fmt.Printf("  Hostname: %s\n", machineFacts.Hostname)
	if len(machineFacts.IPAddresses) > 0 {
		fmt.Printf("  IP Addresses: %v\n", machineFacts.IPAddresses)
		if machineFacts.PrimaryIP != "" {
			fmt.Printf("  Primary IP: %s\n", machineFacts.PrimaryIP)
		}
	}
	if len(machineFacts.Tags) > 0 {
		fmt.Printf("  Tags: %v\n", machineFacts.Tags)
	}
	fmt.Printf("  Updated: %s\n", machineFacts.UpdatedAt.Format("2006-01-02 15:04:05"))
	fmt.Println()
}

// getProjectFactsDBPath returns the project-specific facts database path from context
func getProjectFactsDBPath(cmd *cobra.Command) string {
	if cmd == nil || cmd.Context() == nil {
		return getFactsDBPath() // Fallback to global path
	}

	if projectPath, ok := cmd.Context().Value(projectPathKey{}).(string); ok && projectPath != "" {
		return filepath.Join(projectPath, ".facts.db")
	}

	return getFactsDBPath() // Fallback to global path
}

func runFactsCache(cmd *cobra.Command, _ []string) error {
	// Try to create storage to check if database is corrupted
	storage, err := facts.NewFactStorage(facts.StorageOptions{
		Type: facts.StorageTypeBadger,
		Path: getProjectFactsDBPath(cmd),
	})
	if err != nil {
		return fmt.Errorf("corrupted facts database: %w", err)
	}
	defer storage.Close()

	// Try to load facts to verify database integrity
	manager := facts.NewManagerWithStorage(nil, storage)
	_, err = manager.GetAllFacts()
	if err != nil {
		return fmt.Errorf("corrupted facts database: %w", err)
	}

	fmt.Println("Fact cache is healthy")
	return nil
}

func runFactsCacheClear(cmd *cobra.Command, _ []string) error {
	// Create storage and fact manager
	storage, err := facts.NewFactStorage(facts.StorageOptions{
		Type: facts.StorageTypeBadger,
		Path: getProjectFactsDBPath(cmd),
	})
	if err != nil {
		return fmt.Errorf("corrupted facts database: %w", err)
	}
	defer storage.Close()

	manager := facts.NewManagerWithStorage(nil, storage)

	// Clear the cache
	manager.ClearCache()

	fmt.Println("Fact cache cleared successfully")
	return nil
}

func runFactsCacheExpired(cmd *cobra.Command, _ []string) error {
	// Create storage and fact manager
	storage, err := facts.NewFactStorage(facts.StorageOptions{
		Type: facts.StorageTypeBadger,
		Path: getProjectFactsDBPath(cmd),
	})
	if err != nil {
		return fmt.Errorf("corrupted facts database: %w", err)
	}
	defer storage.Close()

	manager := facts.NewManagerWithStorage(nil, storage)

	// Clear expired facts
	manager.ClearExpiredCache()

	fmt.Println("Expired facts cleared from cache")
	return nil
}

// New fact command functions
func runFactsGather(cmd *cobra.Command, args []string) error {
	logger := logging.GetLogger()

	// Determine target hosts
	hosts, err := determineTargetHosts(args)
	if err != nil {
		return err
	}

	logger.Info("Starting fact gathering",
		logging.Int("parallel", factsParallel),
		logging.Int("timeout", factsTimeout))

	// Create storage and fact manager
	storage, err := facts.NewFactStorage(facts.StorageOptions{
		Type: facts.StorageTypeBadger,    // Default to BadgerDB
		Path: getProjectFactsDBPath(cmd), // Use configured facts database path
	})
	if err != nil {
		return fmt.Errorf("failed to create storage: %w", err)
	}
	defer storage.Close()

	manager := facts.NewManagerWithStorage(nil, storage)
	allCollections, errors := collectFactsFromHosts(manager, hosts, logger)

	// Report results
	if len(errors) > 0 {
		logger.Warn("Some hosts failed fact collection", logging.Int("error_count", len(errors)))
		for _, err := range errors {
			fmt.Printf("Error: %v\n", err)
		}
	}

	if len(allCollections) == 0 {
		return fmt.Errorf("no facts collected from any host")
	}

	// Display results
	displayFactGatheringResults(allCollections, errors)

	logger.Info("Fact gathering completed successfully",
		logging.Int("hosts_processed", len(allCollections)),
		logging.Int("hosts_failed", len(errors)),
		logging.Int("total_facts", getTotalFactCount(allCollections)))

	return nil
}

// determineTargetHosts determines which hosts to collect facts from
func determineTargetHosts(args []string) ([]string, error) {
	logger := logging.GetLogger()

	// Check if we're in a valid project context
	projectPath := getProjectFactsDBPath(nil) // Pass nil to use global path
	if projectPath == "" {
		return nil, fmt.Errorf("invalid project: no project path found")
	}

	// Check if project directory exists and is valid
	if _, err := os.Stat(projectPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("invalid project: project directory does not exist")
	}

	switch {
	case len(args) > 0:
		// Use provided hosts argument
		hosts := strings.Split(args[0], ",")
		for i, host := range hosts {
			hosts[i] = strings.TrimSpace(host)
		}
		return hosts, nil

	case factsConfig != "":
		// Legacy config file support removed - use project-based approach instead
		return nil, fmt.Errorf("legacy config file support has been removed. Use project-based approach with --inventory flag")

	case factsInventory != "":
		// TODO: Load hosts from inventory file
		logger.Info("Inventory file support not yet implemented", logging.String("inventory", factsInventory))
		return nil, fmt.Errorf("inventory file support not yet implemented")

	default:
		// Default to local host
		return []string{"local"}, nil
	}
}

// collectFactsFromHosts collects facts from all specified hosts
func collectFactsFromHosts(manager *facts.Manager, hosts []string, logger logging.Logger) ([]*facts.FactCollection, []error) {
	var allCollections []*facts.FactCollection
	var errors []error

	// For now, collect sequentially (parallel implementation would require goroutines and channels)
	for _, host := range hosts {
		logger.Info("Collecting facts from host", logging.String("host", host))

		collection, err := collectFactsFromHost(manager, host)
		if err != nil {
			logger.Error("Failed to collect facts from host", err, logging.String("host", host))
			errors = append(errors, fmt.Errorf("host %s: %w", host, err))
			continue
		}

		allCollections = append(allCollections, collection)
		logger.Info("Successfully collected facts from host",
			logging.String("host", host),
			logging.Int("fact_count", len(collection.Facts)))
	}

	return allCollections, errors
}

// collectFactsFromHost collects facts from a single host
func collectFactsFromHost(manager *facts.Manager, host string) (*facts.FactCollection, error) {
	if len(factsSpecificKeys) > 0 {
		// Collect specific facts
		return manager.CollectSpecificFacts(host, factsSpecificKeys)
	}
	// Collect all facts and persist to storage
	return manager.GatherAndPersistFacts(host)
}

// displayFactGatheringResults displays the results of fact gathering
func displayFactGatheringResults(allCollections []*facts.FactCollection, errors []error) {
	// Display summary
	fmt.Printf("Fact Gathering Summary:\n")
	fmt.Printf("Hosts processed: %d\n", len(allCollections))
	fmt.Printf("Hosts failed: %d\n", len(errors))
	fmt.Printf("Total facts collected: %d\n", getTotalFactCount(allCollections))
	fmt.Println()

	// Display facts for each host
	for _, collection := range allCollections {
		displayHostFacts(collection)
	}
}

// displayHostFacts displays facts for a single host
func displayHostFacts(collection *facts.FactCollection) {
	fmt.Printf("Host: %s (%d facts)\n", collection.Server, len(collection.Facts))
	fmt.Printf("Collected at: %s\n", collection.Timestamp.Format("2006-01-02 15:04:05"))
	fmt.Println()

	// Show first few facts as preview
	count := 0
	for key, fact := range collection.Facts {
		if count >= 5 { // Show only first 5 facts as preview
			fmt.Printf("... and %d more facts\n", len(collection.Facts)-5)
			break
		}

		valueStr := fmt.Sprintf("%v", fact.Value)
		if len(valueStr) > 40 {
			valueStr = valueStr[:37] + "..."
		}

		fmt.Printf("  %-25s: %s\n", key, valueStr)
		count++
	}
	fmt.Println()
}

// Helper function to get total fact count from collections
func getTotalFactCount(collections []*facts.FactCollection) int {
	total := 0
	for _, collection := range collections {
		total += len(collection.Facts)
	}
	return total
}

func runFactsImport(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("import source file is required")
	}

	source := args[0]

	// Create storage and fact manager
	storage, err := facts.NewFactStorage(facts.StorageOptions{
		Type: facts.StorageTypeBadger,    // Default to BadgerDB
		Path: getProjectFactsDBPath(cmd), // Use configured facts database path
	})
	if err != nil {
		return fmt.Errorf("failed to create storage: %w", err)
	}
	defer storage.Close()

	manager := facts.NewManagerWithStorage(nil, storage)

	// Determine source path based on flags
	var sourcePath string
	switch importSource {
	case "local":
		if importFile != "" {
			sourcePath = importFile
		} else {
			sourcePath = source
		}
	case "http":
		if importURL != "" {
			sourcePath = importURL
		} else {
			sourcePath = source
		}
	default:
		sourcePath = source
	}

	// Create import options
	options := &facts.ImportOptions{
		Source:      importSource,
		Path:        sourcePath,
		MergeMode:   facts.MergeMode(importMergeMode),
		SelectFacts: importSelectFacts,
		Override:    importOverride,
		Validate:    factsValidate || importOverride, // Validate if override is enabled
		DryRun:      importDryRun,
		Server:      importServer,
	}

	// Check if it's a custom source (JSON file or HTTP URL)
	if isCustomSource(sourcePath) {
		// Use enhanced custom fact import
		if err := manager.ImportCustomFactsWithOptions(sourcePath, options); err != nil {
			return fmt.Errorf("failed to import custom facts: %w", err)
		}

		if importDryRun {
			fmt.Println("DRY RUN: Facts would be imported successfully")
		} else {
			fmt.Printf("Custom facts imported successfully from %s\n", sourcePath)
		}
		return nil
	}

	// Fall back to storage import for existing format
	file, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer file.Close()

	// Import facts
	if err := storage.ImportFromJSON(file); err != nil {
		return fmt.Errorf("failed to import facts: %w", err)
	}

	fmt.Printf("Facts imported successfully from %s\n", sourcePath)
	return nil
}

// isCustomSource checks if the source is a custom format (JSON file or HTTPS URL)
func isCustomSource(source string) bool {
	// Empty string is not a valid source
	if source == "" {
		return false
	}

	// Check if it's an HTTP URL (not allowed for security)
	if len(source) > 7 && source[:7] == "http://" {
		return false
	}

	// Check if it's an FTP URL (not supported)
	if len(source) > 6 && source[:6] == "ftp://" {
		return false
	}

	// Check if it's an HTTPS URL (allowed)
	if len(source) > 8 && source[:8] == "https://" {
		return true
	}

	// Check if it's a JSON file (but not an HTTP URL)
	if len(source) > 5 && source[len(source)-5:] == ".json" {
		return true
	}

	// Only HTTPS URLs and JSON files are considered custom sources
	return false
}

func runFactsExport(cmd *cobra.Command, _ []string) error {
	// Create storage and fact manager
	storage, err := facts.NewFactStorage(facts.StorageOptions{
		Type: facts.StorageTypeBadger,    // Default to BadgerDB
		Path: getProjectFactsDBPath(cmd), // Use configured facts database path
	})
	if err != nil {
		return fmt.Errorf("corrupted facts database: %w", err)
	}
	defer storage.Close()

	// Try to load facts to verify database integrity
	manager := facts.NewManagerWithStorage(nil, storage)
	_, err = manager.GetAllFacts()
	if err != nil {
		return fmt.Errorf("corrupted facts database: %w", err)
	}

	// Determine output destination
	var output io.Writer = os.Stdout
	if factsOutput != "" {
		file, err := os.Create(factsOutput)
		if err != nil {
			return fmt.Errorf("failed to create output file: %w", err)
		}
		defer file.Close()
		output = file
	}

	// Export facts
	if err := storage.ExportToJSON(output); err != nil {
		return fmt.Errorf("failed to export facts: %w", err)
	}

	fmt.Fprintf(os.Stderr, "Facts exported successfully\n")
	return nil
}

func runFactsValidate(cmd *cobra.Command, _ []string) error {
	// Create storage and fact manager
	storage, err := facts.NewFactStorage(facts.StorageOptions{
		Type: facts.StorageTypeBadger,    // Default to BadgerDB
		Path: getProjectFactsDBPath(cmd), // Use configured facts database path
	})
	if err != nil {
		return fmt.Errorf("failed to create storage: %w", err)
	}
	defer storage.Close()

	manager := facts.NewManagerWithStorage(nil, storage)

	// Query all facts for validation
	query := &facts.FactQuery{}
	machineFacts, err := manager.QueryMachineFacts(query)
	if err != nil {
		return fmt.Errorf("failed to query facts: %w", err)
	}

	if len(machineFacts) == 0 {
		fmt.Println("No facts found to validate")
		return nil
	}

	// Validate each fact
	var errors []string
	var warnings []string

	for _, fact := range machineFacts {
		// Skip nil facts
		if fact == nil {
			errors = append(errors, "Found nil fact entry")
			continue
		}

		// Basic validation rules
		if fact.MachineName == "" {
			errors = append(errors, fmt.Sprintf("Machine %s: missing machine name", fact.MachineID))
		}

		if fact.SystemID == "" {
			warnings = append(warnings, fmt.Sprintf("Machine %s: missing system ID", fact.MachineID))
		}

		if fact.Hostname == "" {
			warnings = append(warnings, fmt.Sprintf("Machine %s: missing hostname", fact.MachineID))
		}

		if fact.OS == "" {
			warnings = append(warnings, fmt.Sprintf("Machine %s: missing OS information", fact.MachineID))
		}

		// Check for required tags
		if fact.Tags == nil {
			fact.Tags = make(map[string]string)
		}

		if _, hasEnv := fact.Tags["environment"]; !hasEnv {
			warnings = append(warnings, fmt.Sprintf("Machine %s: missing environment tag", fact.MachineID))
		}
	}

	// Report results
	fmt.Printf("Validation Results:\n")
	fmt.Printf("Facts validated: %d\n", len(machineFacts))
	fmt.Printf("Errors: %d\n", len(errors))
	fmt.Printf("Warnings: %d\n", len(warnings))
	fmt.Println()

	if len(errors) > 0 {
		fmt.Println("Errors:")
		for _, err := range errors {
			fmt.Printf("  - %s\n", err)
		}
		fmt.Println()
	}

	if len(warnings) > 0 {
		fmt.Println("Warnings:")
		for _, warning := range warnings {
			fmt.Printf("  - %s\n", warning)
		}
		fmt.Println()
	}

	if len(errors) > 0 {
		return fmt.Errorf("validation failed with %d errors", len(errors))
	}

	fmt.Println("Validation completed successfully")
	return nil
}

func runFactsQuery(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("query expression is required")
	}

	expression := args[0]

	// Parse query expression
	query, err := parseQueryExpression(expression)
	if err != nil {
		return fmt.Errorf("failed to parse query expression: %w", err)
	}

	// Set limit if specified
	if factsLimit > 0 {
		query.Limit = factsLimit
	}

	// Create storage and fact manager
	storage, err := facts.NewFactStorage(facts.StorageOptions{
		Type: facts.StorageTypeBadger,    // Default to BadgerDB
		Path: getProjectFactsDBPath(cmd), // Use configured facts database path
	})
	if err != nil {
		return fmt.Errorf("failed to create storage: %w", err)
	}
	defer storage.Close()

	manager := facts.NewManagerWithStorage(nil, storage)

	// Query facts
	machineFacts, err := manager.QueryMachineFacts(query)
	if err != nil {
		return fmt.Errorf("failed to query facts: %w", err)
	}

	if len(machineFacts) == 0 {
		fmt.Println("No facts found matching query")
		return nil
	}

	// Display results
	fmt.Printf("Found %d facts matching query:\n\n", len(machineFacts))
	for _, fact := range machineFacts {
		displayMachineFacts(fact)
	}

	return nil
}

// parseQueryExpression parses a query expression string into a FactQuery
func parseQueryExpression(expression string) (*facts.FactQuery, error) {
	query := &facts.FactQuery{
		Tags: make(map[string]string),
	}

	// Simple parsing for now - split by comma and parse key=value pairs
	pairs := strings.Split(expression, ",")
	for _, pair := range pairs {
		pair = strings.TrimSpace(pair)
		if pair == "" {
			continue
		}

		parts := strings.SplitN(pair, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid query pair: %s", pair)
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "machine":
			query.MachineName = value
		case "os":
			query.OS = value
		case "project":
			query.ProjectName = value
		case "tag":
			// Handle tag=key:value format
			tagParts := strings.SplitN(value, ":", 2)
			if len(tagParts) == 2 {
				query.Tags[tagParts[0]] = tagParts[1]
			} else {
				query.Tags[value] = "" // Tag exists but no value
			}
		case "limit":
			// Limit is handled separately via flag
		default:
			return nil, fmt.Errorf("unknown query field: %s", key)
		}
	}

	return query, nil
}
