package ssh

import (
	"fmt"
	"os"
	"sync"
	"time"

	"spooky/internal/config"
	"spooky/internal/logging"
)

// ExecuteConfig executes all actions in the configuration
func ExecuteConfig(cfg *config.Config) error {
	logger := logging.GetLogger()

	logger.Info("Starting configuration execution",
		logging.Int("action_count", len(cfg.Actions)),
		logging.Int("server_count", len(cfg.Servers)),
	)

	// Initialize index cache for enterprise-scale performance
	indexCache := &config.IndexCache{}

	for i := range cfg.Actions {
		action := &cfg.Actions[i]
		startTime := time.Now()

		logger.Info("Executing action",
			logging.Action(action.Name),
			logging.String("description", action.Description),
		)

		// Get target servers for this action using optimized lookup
		var targetServers []*config.Server
		var err error

		// Use enterprise-scale lookup for better performance
		index := indexCache.GetIndex(cfg)
		targetServers, err = config.GetServersForActionLarge(cfg, action, index)
		if err != nil {
			logger.Error("Failed to get servers for action", err,
				logging.Action(action.Name),
			)
			return fmt.Errorf("failed to get servers for action %s: %w", action.Name, err)
		}

		logger.Info("Action target servers determined",
			logging.Action(action.Name),
			logging.Int("target_server_count", len(targetServers)),
		)

		// Execute on each server
		if action.Parallel {
			err = executeActionParallel(action, targetServers)
		} else {
			err = executeActionSequential(action, targetServers)
		}

		if err != nil {
			logger.Error("Failed to execute action", err,
				logging.Action(action.Name),
				logging.Duration("duration_ms", time.Since(startTime).Milliseconds()),
			)
			return fmt.Errorf("failed to execute action %s: %w", action.Name, err)
		}

		logger.Info("Action completed successfully",
			logging.Action(action.Name),
			logging.Duration("duration_ms", time.Since(startTime).Milliseconds()),
		)
	}

	logger.Info("All actions completed successfully",
		logging.Int("total_actions", len(cfg.Actions)),
	)
	return nil
}

// executeActionSequential executes an action sequentially on all target servers
func executeActionSequential(action *config.Action, servers []*config.Server) error {
	logger := logging.GetLogger()

	// Validate action before connecting
	if action.Command != "" && action.Script != "" {
		logger.Error("Action validation failed", fmt.Errorf("both command and script specified"),
			logging.Action(action.Name),
		)
		return fmt.Errorf("action %s: both command and script specified", action.Name)
	}
	if action.Command == "" && action.Script == "" {
		logger.Error("Action validation failed", fmt.Errorf("neither command nor script specified"),
			logging.Action(action.Name),
		)
		return fmt.Errorf("action %s: neither command nor script specified", action.Name)
	}

	for _, server := range servers {
		startTime := time.Now()

		logger.Info("Connecting to server",
			logging.Server(server.Name),
			logging.Host(server.Host),
			logging.Port(server.Port),
			logging.String("user", server.User),
		)

		// Create SSH client
		client, err := NewSSHClient(server, 30) // Default timeout
		if err != nil {
			logger.Error("Failed to connect to server", err,
				logging.Server(server.Name),
				logging.Host(server.Host),
				logging.Port(server.Port),
			)
			continue
		}

		// Execute the action
		var output string
		if action.Command != "" {
			output, err = client.ExecuteCommand(action.Command)
		} else if action.Script != "" {
			output, err = client.ExecuteScript(action.Script)
		}

		// Close client after execution
		if closeErr := client.Close(); closeErr != nil {
			logger.Warn("Failed to close SSH connection",
				logging.Server(server.Name),
				logging.Error(closeErr),
			)
		}

		if err != nil {
			logger.Error("Failed to execute action on server", err,
				logging.Server(server.Name),
				logging.Action(action.Name),
				logging.Duration("duration_ms", time.Since(startTime).Milliseconds()),
			)
			continue
		}

		logger.Info("Action executed successfully on server",
			logging.Server(server.Name),
			logging.Action(action.Name),
			logging.Duration("duration_ms", time.Since(startTime).Milliseconds()),
			logging.String("output_length", fmt.Sprintf("%d chars", len(output))),
		)

		if output != "" {
			logger.Debug("Command output",
				logging.Server(server.Name),
				logging.Action(action.Name),
				logging.String("output", output),
			)
		}
	}

	return nil
}

// executeActionParallel executes an action in parallel on all target servers
func executeActionParallel(action *config.Action, servers []*config.Server) error {
	logger := logging.GetLogger()

	// Validate action before connecting
	if action.Command != "" && action.Script != "" {
		logger.Error("Action validation failed", fmt.Errorf("both command and script specified"),
			logging.Action(action.Name),
		)
		return fmt.Errorf("action %s: both command and script specified", action.Name)
	}
	if action.Command == "" && action.Script == "" {
		logger.Error("Action validation failed", fmt.Errorf("neither command nor script specified"),
			logging.Action(action.Name),
		)
		return fmt.Errorf("action %s: neither command nor script specified", action.Name)
	}

	// Validate script file exists before connecting (for parallel execution)
	if action.Script != "" {
		if _, err := os.Stat(action.Script); os.IsNotExist(err) {
			logger.Error("Script file not found", err,
				logging.Action(action.Name),
				logging.String("script_file", action.Script),
			)
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
			startTime := time.Now()

			logger.Info("Connecting to server (parallel)",
				logging.Server(s.Name),
				logging.Host(s.Host),
				logging.Port(s.Port),
				logging.String("user", s.User),
			)

			// Create SSH client
			client, err := NewSSHClient(s, 30) // Default timeout
			if err != nil {
				logger.Error("Failed to connect to server (parallel)", err,
					logging.Server(s.Name),
					logging.Host(s.Host),
					logging.Port(s.Port),
				)
				errors <- fmt.Errorf("failed to connect to %s: %w", s.Name, err)
				return
			}
			// Close client when function returns
			defer func() {
				if closeErr := client.Close(); closeErr != nil {
					logger.Warn("Failed to close SSH connection (parallel)",
						logging.Server(s.Name),
						logging.Error(closeErr),
					)
				}
			}()

			// Execute the action
			var output string
			if action.Command != "" {
				output, err = client.ExecuteCommand(action.Command)
			} else if action.Script != "" {
				output, err = client.ExecuteScript(action.Script)
			}

			if err != nil {
				logger.Error("Failed to execute action on server (parallel)", err,
					logging.Server(s.Name),
					logging.Action(action.Name),
					logging.Duration("duration_ms", time.Since(startTime).Milliseconds()),
				)
				errors <- fmt.Errorf("failed to execute action on %s: %w", s.Name, err)
				return
			}

			logger.Info("Action executed successfully on server (parallel)",
				logging.Server(s.Name),
				logging.Action(action.Name),
				logging.Duration("duration_ms", time.Since(startTime).Milliseconds()),
				logging.String("output_length", fmt.Sprintf("%d chars", len(output))),
			)

			results <- fmt.Sprintf("✅ Success on %s\n%s", s.Name, indentOutput(output))
		}(server)
	}

	// Wait for all goroutines to complete
	wg.Wait()
	close(results)
	close(errors)

	// Collect results
	for result := range results {
		logger.Info("Parallel execution result", logging.String("result", result))
	}

	// Check for errors
	select {
	case err := <-errors:
		logger.Error("Parallel execution failed", err, logging.Action(action.Name))
		return err
	default:
		logger.Info("Parallel execution completed successfully", logging.Action(action.Name))
		return nil
	}
}
