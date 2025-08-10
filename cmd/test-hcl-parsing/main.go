package main

import (
	spookytypes "spooky/internal/types"
	"fmt"
	"log"
	"os"
	"strings"

	"spooky/internal/config/loading"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s <project-path>\n", os.Args[0])
		fmt.Printf("Example: %s testdata\n", os.Args[0])
		os.Exit(1)
	}

	projectPath := os.Args[1]

	// Test enhanced HCL parsing with validation
	fmt.Printf("Testing enhanced HCL parsing for project: %s\n", projectPath)
	fmt.Println("=" + strings.Repeat("=", len(projectPath)+50))

	// Load and validate actions configuration
	actionsConfig, err := loading.LoadActionsConfig(projectPath)
	if err != nil {
		log.Fatalf("Failed to load actions config: %v", err)
	}

	fmt.Printf("Successfully loaded %d actions\n\n", len(actionsConfig.Actions))

	// Display loaded actions
	for i, action := range actionsConfig.Actions {
		fmt.Printf("Action %d: %s\n", i+1, action.Name)
		fmt.Printf("  Type: %s\n", action.Type)
		fmt.Printf("  Description: %s\n", action.Description)

		switch action.Type {
		case "command":
			fmt.Printf("  Command: %s\n", action.Command)
		case "script":
			fmt.Printf("  Script: %s\n", action.Script)
		case "template_deploy", "template_evaluate", "template_validate", "template_cleanup":
			if action.Template != nil {
				fmt.Printf("  Template Source: %s\n", action.Template.Source)
				fmt.Printf("  Template Destination: %s\n", action.Template.Destination)
			}
		case "file_copy":
			fmt.Printf("  File Copy: (not supported in config Action type)\n")
		}

		if action.Timeout > 0 {
			fmt.Printf("  Timeout: %d seconds\n", action.Timeout)
		}

		if len(action.Dependencies) > 0 {
			fmt.Printf("  Dependencies: %v\n", action.Dependencies)
		}

		if action.Critical {
			fmt.Printf("  Critical: true\n")
		}

		if action.Retries > 0 {
			fmt.Printf("  Retries: %d (delay: %d seconds)\n", action.Retries, action.RetryDelay)
		}

		if len(action.Environment) > 0 {
			fmt.Printf("  Environment: %v\n", action.Environment)
		}

		if action.ResourceLimits != nil {
			fmt.Printf("  Resource Limits: Memory=%dMB, CPU=%d%%, Disk=%dMB\n",
				action.ResourceLimits.MemoryMB,
				action.ResourceLimits.CPUPercent,
				action.ResourceLimits.DiskMB)
		}

		fmt.Println()
	}

	fmt.Println("Enhanced HCL parsing and validation completed successfully!")
}
