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
		logging.Int("machine_count", len(cfg.Machines)),
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

		// Get target machines for this action using optimized lookup
		var targetMachines []*config.Machine
		var err error

		// Use enterprise-scale lookup for better performance
		index := indexCache.GetIndex(cfg)
		targetMachines, err = config.GetMachinesForActionLarge(cfg, action, index)
		if err != nil {
			logger.Error("Failed to get machines for action", err,
				logging.Action(action.Name),
			)
			return fmt.Errorf("failed to get machines for action %s: %w", action.Name, err)
		}

		logger.Info("Action target machines determined",
			logging.Action(action.Name),
			logging.Int("target_machine_count", len(targetMachines)),
		)

		// Execute on each machine
		if action.Parallel {
			err = executeActionParallel(action, targetMachines)
		} else {
			err = executeActionSequential(action, targetMachines)
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

// executeActionSequential executes an action sequentially on all target machines
func executeActionSequential(action *config.Action, machines []*config.Machine) error {
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

	for _, machine := range machines {
		startTime := time.Now()

		logger.Info("Connecting to machine",
			logging.Server(machine.Name),
			logging.Host(machine.Host),
			logging.Port(machine.Port),
			logging.String("user", machine.User),
		)

		// Create SSH client
		client, err := NewSSHClient(machine, 30) // Default timeout
		if err != nil {
			logger.Error("Failed to connect to machine", err,
				logging.Server(machine.Name),
				logging.Host(machine.Host),
				logging.Port(machine.Port),
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
				logging.Server(machine.Name),
				logging.Error(closeErr),
			)
		}

		if err != nil {
			logger.Error("Failed to execute action on machine", err,
				logging.Server(machine.Name),
				logging.Action(action.Name),
				logging.Duration("duration_ms", time.Since(startTime).Milliseconds()),
			)
			continue
		}

		logger.Info("Action executed successfully on machine",
			logging.Server(machine.Name),
			logging.Action(action.Name),
			logging.Duration("duration_ms", time.Since(startTime).Milliseconds()),
			logging.String("output_length", fmt.Sprintf("%d chars", len(output))),
		)

		if output != "" {
			logger.Debug("Command output",
				logging.Server(machine.Name),
				logging.Action(action.Name),
				logging.String("output", output),
			)
		}
	}

	return nil
}

// validateActionForParallel validates an action for parallel execution
func validateActionForParallel(action *config.Action) error {
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

	return nil
}

// executeActionOnMachine executes an action on a single machine in a goroutine
func executeActionOnMachine(action *config.Action, machine *config.Machine, results chan<- string, errors chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()
	logger := logging.GetLogger()
	startTime := time.Now()

	logger.Info("Connecting to machine (parallel)",
		logging.Server(machine.Name),
		logging.Host(machine.Host),
		logging.Port(machine.Port),
		logging.String("user", machine.User),
	)

	// Create SSH client
	client, err := NewSSHClient(machine, 30) // Default timeout
	if err != nil {
		logger.Error("Failed to connect to machine (parallel)", err,
			logging.Server(machine.Name),
			logging.Host(machine.Host),
			logging.Port(machine.Port),
		)
		errors <- fmt.Errorf("failed to connect to %s: %w", machine.Name, err)
		return
	}
	// Close client when function returns
	defer func() {
		if closeErr := client.Close(); closeErr != nil {
			logger.Warn("Failed to close SSH connection (parallel)",
				logging.Server(machine.Name),
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
		logger.Error("Failed to execute action on machine (parallel)", err,
			logging.Server(machine.Name),
			logging.Action(action.Name),
			logging.Duration("duration_ms", time.Since(startTime).Milliseconds()),
		)
		errors <- fmt.Errorf("failed to execute action on %s: %w", machine.Name, err)
		return
	}

	logger.Info("Action executed successfully on machine (parallel)",
		logging.Server(machine.Name),
		logging.Action(action.Name),
		logging.Duration("duration_ms", time.Since(startTime).Milliseconds()),
		logging.String("output_length", fmt.Sprintf("%d chars", len(output))),
	)

	results <- fmt.Sprintf("✅ Success on %s\n%s", machine.Name, indentOutput(output))
}

// executeActionParallel executes an action in parallel on all target machines
func executeActionParallel(action *config.Action, machines []*config.Machine) error {
	logger := logging.GetLogger()

	// Validate action before connecting
	if err := validateActionForParallel(action); err != nil {
		return err
	}

	var wg sync.WaitGroup
	results := make(chan string, len(machines))
	errors := make(chan error, len(machines))

	for _, machine := range machines {
		wg.Add(1)
		go executeActionOnMachine(action, machine, results, errors, &wg)
	}

	// Wait for all goroutines to complete
	wg.Wait()
	close(results)
	close(errors)

	// Collect results
	for result := range results {
		logger.Info("Parallel execution result", logging.String("result", result))
	}

	// Collect all errors
	var allErrors []error
	for err := range errors {
		allErrors = append(allErrors, err)
	}

	// Check for errors
	if len(allErrors) > 0 {
		// Log all errors for debugging
		for i, err := range allErrors {
			logger.Error("Parallel execution error", err,
				logging.Action(action.Name),
				logging.Int("error_index", i+1),
				logging.Int("total_errors", len(allErrors)),
			)
		}

		// Return the first error with context about total errors
		if len(allErrors) == 1 {
			logger.Error("Parallel execution failed", allErrors[0], logging.Action(action.Name))
			return allErrors[0]
		}
		combinedError := fmt.Errorf("parallel execution failed on %d servers: %w", len(allErrors), allErrors[0])
		logger.Error("Parallel execution failed", combinedError, logging.Action(action.Name))
		return combinedError
	}

	logger.Info("Parallel execution completed successfully", logging.Action(action.Name))
	return nil
}
