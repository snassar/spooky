package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"spooky/internal/config"
	"spooky/internal/facts"
	"spooky/internal/logging"
	"spooky/internal/ssh"
)

var (
	parallel int // coverage-ignore: global variable declaration
	timeout  int // coverage-ignore: global variable declaration

	// Execute command flags
	hosts    string
	retry    int
	tags     string
	skipTags string

	// Validate command flags
	schema string
	strict bool
	format string

	// List command flags
	filter  string
	sort    string
	reverse bool
)

var ExecuteCmd = &cobra.Command{
	Use:   "execute <source>",
	Short: "Execute configuration files or remote sources",
	Long:  `Execute actions defined in configuration files or remote sources on remote servers via SSH`,
	Args:  cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		logger := logging.GetLogger()
		source := args[0]

		// TODO: Support remote sources (Git, S3, HTTP)
		// For now, only support local files
		if !isLocalFile(source) {
			return fmt.Errorf("remote sources not yet supported: %s", source)
		}

		// Validate config file exists
		if _, err := os.Stat(source); os.IsNotExist(err) {
			logger.Error("Config file not found", err, logging.String("config_file", source))
			return fmt.Errorf("config file %s does not exist", source) // coverage-ignore: file system error, hard to test
		}

		// Parse and execute configuration
		config, err := config.ParseConfig(source)
		if err != nil {
			logger.Error("Failed to parse configuration", err, logging.String("config_file", source))
			return fmt.Errorf("failed to parse config: %w", err)
		}

		logger.Info("Starting configuration execution",
			logging.String("config_file", source),
			logging.Int("action_count", len(config.Actions)),
			logging.Int("server_count", len(config.Servers)),
			logging.Int("parallel", parallel),
			logging.Int("timeout", timeout),
		)

		return ssh.ExecuteConfig(config)
	},
}

var ValidateCmd = &cobra.Command{
	Use:   "validate <source>",
	Short: "Validate configuration files and templates",
	Long:  `Validate the syntax and structure of configuration files and templates`,
	Args:  cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		logger := logging.GetLogger()
		source := args[0]

		// TODO: Support remote sources (Git, S3, HTTP)
		// For now, only support local files
		if !isLocalFile(source) {
			return fmt.Errorf("remote sources not yet supported: %s", source)
		}

		// Validate config file exists
		if _, err := os.Stat(source); os.IsNotExist(err) {
			logger.Error("Config file not found", err, logging.String("config_file", source))
			return fmt.Errorf("config file %s does not exist", source)
		}

		// Parse configuration
		config, err := config.ParseConfig(source)
		if err != nil {
			logger.Error("Configuration validation failed", err, logging.String("config_file", source))
			return fmt.Errorf("validation failed: %w", err)
		}

		logger.Info("Configuration file validated successfully",
			logging.String("config_file", source),
			logging.Int("server_count", len(config.Servers)),
			logging.Int("action_count", len(config.Actions)),
		)

		return nil
	},
}

var ListCmd = &cobra.Command{
	Use:   "list <resource|config-file>",
	Short: "List resources, facts, or configurations",
	Long:  `Display resources, facts, or configurations based on the specified resource type or from a configuration file`,
	Args:  cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		logger := logging.GetLogger()
		resource := args[0]

		// Check if it's a config file path
		if strings.HasSuffix(resource, ".hcl") || strings.HasSuffix(resource, ".yaml") || strings.HasSuffix(resource, ".yml") {
			return listFromConfigFile(logger, resource)
		}

		// Otherwise treat as resource type
		switch resource {
		case "servers":
			return listServers(logger)
		case "facts":
			return listFacts(logger)
		case "templates":
			return listTemplates(logger)
		case "configs":
			return listConfigs(logger)
		case "actions":
			return listActions(logger)
		default:
			return fmt.Errorf("unknown resource type: %s. Supported types: servers, facts, templates, configs, actions", resource)
		}
	},
}

// InitCommands initializes all CLI commands and their flags
func InitCommands() {
	// Execute command flags
	ExecuteCmd.Flags().StringVar(&hosts, "hosts", "", "Comma-separated list of target hosts")
	ExecuteCmd.Flags().IntVar(&parallel, "parallel", 5, "Number of parallel executions")
	ExecuteCmd.Flags().IntVar(&timeout, "timeout", 30, "Execution timeout per host in seconds")
	ExecuteCmd.Flags().IntVar(&retry, "retry", 3, "Number of retry attempts")
	ExecuteCmd.Flags().StringVar(&tags, "tags", "", "Comma-separated list of tags to execute")
	ExecuteCmd.Flags().StringVar(&skipTags, "skip-tags", "", "Comma-separated list of tags to skip")

	// Validate command flags
	ValidateCmd.Flags().StringVar(&schema, "schema", "", "Path to schema file for validation")
	ValidateCmd.Flags().BoolVar(&strict, "strict", false, "Enable strict validation mode")
	ValidateCmd.Flags().StringVar(&format, "format", "text", "Output format: text, json, yaml")

	// List command flags
	ListCmd.Flags().StringVar(&format, "format", "table", "Output format: table, json, yaml")
	ListCmd.Flags().StringVar(&filter, "filter", "", "Filter results by expression")
	ListCmd.Flags().StringVar(&sort, "sort", "name", "Sort field")
	ListCmd.Flags().BoolVar(&reverse, "reverse", false, "Reverse sort order")

	// Initialize facts commands
	initFactsCommands()
}

// isLocalFile checks if the source is a local file
func isLocalFile(source string) bool {
	// Check if it's a remote source (Git, S3, HTTP)
	if strings.HasPrefix(source, "github.com/") ||
		strings.HasPrefix(source, "git://") ||
		strings.HasPrefix(source, "s3://") ||
		strings.HasPrefix(source, "http://") ||
		strings.HasPrefix(source, "https://") {
		return false
	}
	return true
}

// listFromConfigFile lists resources from a configuration file
func listFromConfigFile(logger logging.Logger, configFile string) error {
	// Validate config file exists
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		logger.Error("Config file not found", err, logging.String("config_file", configFile))
		return fmt.Errorf("config file %s does not exist", configFile)
	}

	// Parse configuration
	config, err := config.ParseConfig(configFile)
	if err != nil {
		logger.Error("Failed to parse configuration", err, logging.String("config_file", configFile))
		return fmt.Errorf("failed to parse config: %w", err)
	}

	logger.Info("Listing resources from configuration file",
		logging.String("config_file", configFile),
		logging.Int("server_count", len(config.Servers)),
		logging.Int("action_count", len(config.Actions)))

	// Display servers
	if len(config.Servers) > 0 {
		fmt.Printf("Servers (%d):\n", len(config.Servers))
		for _, server := range config.Servers {
			fmt.Printf("  - %s (%s@%s:%d)\n", server.Name, server.User, server.Host, server.Port)
		}
		fmt.Println()
	}

	// Display actions
	if len(config.Actions) > 0 {
		fmt.Printf("Actions (%d):\n", len(config.Actions))
		for i := range config.Actions {
			action := &config.Actions[i]
			desc := action.Description
			if desc == "" {
				desc = "No description"
			}
			fmt.Printf("  - %s: %s\n", action.Name, desc)
		}
		fmt.Println()
	}

	// Display summary
	fmt.Printf("Configuration Summary:\n")
	fmt.Printf("  File: %s\n", configFile)
	fmt.Printf("  Servers: %d\n", len(config.Servers))
	fmt.Printf("  Actions: %d\n", len(config.Actions))

	return nil
}

// List functions for different resource types
func listServers(logger logging.Logger) error {
	// TODO: Implement server listing from configuration or inventory
	logger.Info("Listing servers (not yet implemented)")
	return fmt.Errorf("server listing not yet implemented")
}

func listFacts(logger logging.Logger) error {
	// Create fact manager
	manager := facts.NewManager(nil)

	// Try to get cached facts first
	allFacts, err := manager.GetAllFacts()
	if err != nil {
		logger.Error("Failed to retrieve cached facts", err)
		return fmt.Errorf("failed to retrieve cached facts: %w", err)
	}

	// If no cached facts, collect from local server
	if len(allFacts) == 0 {
		logger.Info("No cached facts found, collecting from local server")
		collection, err := manager.CollectAllFacts("local")
		if err != nil {
			logger.Error("Failed to collect facts from local server", err)
			return fmt.Errorf("failed to collect facts from local server: %w", err)
		}

		// Convert collection to slice of facts
		for _, fact := range collection.Facts {
			allFacts = append(allFacts, fact)
		}
	}

	if len(allFacts) == 0 {
		fmt.Println("No facts found. Use 'spooky facts collect <server>' to gather facts first.")
		return nil
	}

	// Check output format
	if format == "json" {
		// Create a structured response for JSON output
		response := map[string]interface{}{
			"summary": map[string]interface{}{
				"total_facts": len(allFacts),
				"servers":     len(getUniqueServers(allFacts)),
			},
			"facts": allFacts,
		}

		jsonData, err := json.MarshalIndent(response, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal facts to JSON: %w", err)
		}
		fmt.Println(string(jsonData))
	} else {
		// Group facts by server
		serverFacts := make(map[string][]*facts.Fact)
		for _, fact := range allFacts {
			serverFacts[fact.Server] = append(serverFacts[fact.Server], fact)
		}

		// Display facts in table format
		fmt.Printf("Facts Summary:\n")
		fmt.Printf("Total facts: %d\n", len(allFacts))
		fmt.Printf("Servers: %d\n\n", len(serverFacts))

		for server, facts := range serverFacts {
			fmt.Printf("Server: %s (%d facts)\n", server, len(facts))
			fmt.Printf("%-30s %-15s %-20s %s\n", "KEY", "SOURCE", "TTL", "VALUE")
			fmt.Printf("%s\n", strings.Repeat("-", 80))

			for _, fact := range facts {
				// Truncate value if too long
				valueStr := fmt.Sprintf("%v", fact.Value)
				if len(valueStr) > 40 {
					valueStr = valueStr[:37] + "..."
				}

				// Format TTL
				ttlStr := "expired"
				if fact.TTL > 0 {
					ttlStr = fact.TTL.String()
				}

				fmt.Printf("%-30s %-15s %-20s %s\n",
					fact.Key,
					fact.Source,
					ttlStr,
					valueStr)
			}
			fmt.Println()
		}
	}

	logger.Info("Facts listed successfully", logging.Int("total_facts", len(allFacts)), logging.Int("servers", len(getUniqueServers(allFacts))))
	return nil
}

// Helper function to get unique servers from facts
func getUniqueServers(facts []*facts.Fact) []string {
	servers := make(map[string]bool)
	for _, fact := range facts {
		servers[fact.Server] = true
	}

	result := make([]string, 0, len(servers))
	for server := range servers {
		result = append(result, server)
	}
	return result
}

func listTemplates(logger logging.Logger) error {
	// TODO: Implement template listing
	logger.Info("Listing templates (not yet implemented)")
	return fmt.Errorf("template listing not yet implemented")
}

func listConfigs(logger logging.Logger) error {
	// TODO: Implement config listing
	logger.Info("Listing configs (not yet implemented)")
	return fmt.Errorf("config listing not yet implemented")
}

func listActions(logger logging.Logger) error {
	// TODO: Implement action listing
	logger.Info("Listing actions (not yet implemented)")
	return fmt.Errorf("action listing not yet implemented")
}
