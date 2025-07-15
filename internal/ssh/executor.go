package ssh

import (
	"fmt"
	"os"
	"sync"

	"spooky/internal/config"
)

// ExecuteConfig executes all actions in the configuration
func ExecuteConfig(cfg *config.Config) error {
	fmt.Printf("🚀 Starting execution of %d actions...\n", len(cfg.Actions))

	for _, action := range cfg.Actions {
		fmt.Printf("\n⚡ Executing action: %s\n", action.Name)
		if action.Description != "" {
			fmt.Printf("📝 Description: %s\n", action.Description)
		}

		// Get target servers for this action
		targetServers, err := config.GetServersForAction(&action, cfg)
		if err != nil {
			return fmt.Errorf("failed to get servers for action %s: %w", action.Name, err)
		}

		fmt.Printf("🌐 Target servers: %d\n", len(targetServers))

		// Execute on each server
		if action.Parallel {
			err = executeActionParallel(&action, targetServers)
		} else {
			err = executeActionSequential(&action, targetServers)
		}

		if err != nil {
			return fmt.Errorf("failed to execute action %s: %w", action.Name, err)
		}
	}

	fmt.Println("\n✅ All actions completed successfully!")
	return nil
}

// executeActionSequential executes an action sequentially on all target servers
func executeActionSequential(action *config.Action, servers []*config.Server) error {
	// Validate action before connecting
	if action.Command != "" && action.Script != "" {
		return fmt.Errorf("action %s: both command and script specified", action.Name)
	}
	if action.Command == "" && action.Script == "" {
		return fmt.Errorf("action %s: neither command nor script specified", action.Name)
	}

	for _, server := range servers {
		fmt.Printf("  🔗 Connecting to %s (%s@%s:%d)...\n", server.Name, server.User, server.Host, server.Port)

		// Create SSH client
		client, err := NewSSHClient(server, 30) // Default timeout
		if err != nil {
			fmt.Printf("  ❌ Failed to connect to %s: %v\n", server.Name, err)
			continue
		}
		defer client.Close()

		// Execute the action
		var output string
		if action.Command != "" {
			output, err = client.ExecuteCommand(action.Command)
		} else if action.Script != "" {
			output, err = client.ExecuteScript(action.Script)
		}

		if err != nil {
			fmt.Printf("  ❌ Failed to execute action on %s: %v\n", server.Name, err)
			continue
		}

		fmt.Printf("  ✅ Success on %s\n", server.Name)
		if output != "" {
			fmt.Printf("  📄 Output:\n%s\n", indentOutput(output))
		}
	}

	return nil
}

// executeActionParallel executes an action in parallel on all target servers
func executeActionParallel(action *config.Action, servers []*config.Server) error {
	// Validate action before connecting
	if action.Command != "" && action.Script != "" {
		return fmt.Errorf("action %s: both command and script specified", action.Name)
	}
	if action.Command == "" && action.Script == "" {
		return fmt.Errorf("action %s: neither command nor script specified", action.Name)
	}

	// Validate script file exists before connecting (for parallel execution)
	if action.Script != "" {
		if _, err := os.Stat(action.Script); os.IsNotExist(err) {
			return fmt.Errorf("failed to read script file %s: %w", action.Script, err)
		}
	}

	var wg sync.WaitGroup
	results := make(chan string, len(servers))
	errors := make(chan error, len(servers))

	for _, server := range servers {
		wg.Add(1)
		go func(s *config.Server) {
			defer wg.Done()

			fmt.Printf("  🔗 Connecting to %s (%s@%s:%d)...\n", s.Name, s.User, s.Host, s.Port)

			// Create SSH client
			client, err := NewSSHClient(s, 30) // Default timeout
			if err != nil {
				errors <- fmt.Errorf("failed to connect to %s: %w", s.Name, err)
				return
			}
			defer client.Close()

			// Execute the action
			var output string
			if action.Command != "" {
				output, err = client.ExecuteCommand(action.Command)
			} else if action.Script != "" {
				output, err = client.ExecuteScript(action.Script)
			}

			if err != nil {
				errors <- fmt.Errorf("failed to execute action on %s: %w", s.Name, err)
				return
			}

			results <- fmt.Sprintf("✅ Success on %s\n%s", s.Name, indentOutput(output))
		}(server)
	}

	// Wait for all goroutines to complete
	wg.Wait()
	close(results)
	close(errors)

	// Collect results
	for result := range results {
		fmt.Println(result)
	}

	// Check for errors
	select {
	case err := <-errors:
		return err
	default:
		return nil
	}
}
