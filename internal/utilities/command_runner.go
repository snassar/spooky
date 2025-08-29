package utilities

import (
	"context"
	"fmt"
	"io"
	"os/exec"
	"time"

	"spooky/internal/schemas"

	"github.com/pkg/errors"
)

// RemoteCommandRunner defines the interface for remote command execution
type RemoteCommandRunner interface {
	RunCommandOnMachine(ctx context.Context, machine *schemas.MachinesMachineV1, command string) (*schemas.CommandResult, error)
}

// RunCommand is a unified low-level function that handles command execution
// with timeouts, retries, and local/remote execution as parameters
func RunCommand(ctx context.Context, hostname string, command string, machine *schemas.MachinesMachineV1, remoteRunner RemoteCommandRunner, timeout time.Duration, retries int, retryDelay time.Duration) (*schemas.CommandResult, error) {
	// Apply defaults
	if timeout == 0 {
		timeout = 30 * time.Second
	}
	if retryDelay == 0 {
		retryDelay = 1 * time.Second
	}

	// Apply timeout to context
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var lastResult *schemas.CommandResult
	var lastErr error

	// Execute with retries
	for attempt := 0; attempt <= retries; attempt++ {
		if attempt > 0 {
			// Wait before retry
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(retryDelay):
			}
		}

		// Execute the command
		var result *schemas.CommandResult
		var err error

		if hostname == "localhost" {
			// Local execution - pass raw command string (like SSH)
			result, err = runLocalCommand(ctx, command)
		} else {
			// Remote execution requires machine config
			if machine == nil {
				lastErr = errors.New("machine configuration required for remote execution")
				continue
			}
			if remoteRunner == nil {
				lastErr = errors.New("remote runner required for remote command execution")
				continue
			}
			result, err = remoteRunner.RunCommandOnMachine(ctx, machine, command)
		}

		if err != nil {
			lastErr = err
			continue
		}

		if result.ExitCode == 0 {
			return result, nil // Success
		}

		lastResult = result
		lastErr = fmt.Errorf("command returned non-zero exit code: %d", result.ExitCode)
	}

	// Return the last result if we have one, otherwise return the last error
	if lastResult != nil {
		return lastResult, lastErr
	}
	return nil, lastErr
}

// runLocalCommand executes a command on the local machine
func runLocalCommand(ctx context.Context, command string) (*schemas.CommandResult, error) {
	// Execute command directly (like SSH does)
	cmd := exec.CommandContext(ctx, "sh", "-c", command)

	// Capture output
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, errors.Wrap(err, "failed to create stdout pipe")
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, errors.Wrap(err, "failed to create stderr pipe")
	}

	// Start command
	if err := cmd.Start(); err != nil {
		return &schemas.CommandResult{
			ExitCode: 1,
			Error:    errors.Wrap(err, "failed to start command"),
		}, nil
	}

	// Read output
	stdoutBytes, err := readPipe(stdout)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read stdout")
	}

	stderrBytes, err := readPipe(stderr)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read stderr")
	}

	// Wait for completion
	if err := cmd.Wait(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return &schemas.CommandResult{
				ExitCode: exitErr.ExitCode(),
				Stdout:   string(stdoutBytes),
				Stderr:   string(stderrBytes),
				Error:    err,
			}, nil
		}
		return &schemas.CommandResult{
			ExitCode: 1,
			Stdout:   string(stdoutBytes),
			Stderr:   string(stderrBytes),
			Error:    err,
		}, nil
	}

	return &schemas.CommandResult{
		ExitCode: 0,
		Stdout:   string(stdoutBytes),
		Stderr:   string(stderrBytes),
	}, nil
}

// readPipe reads all data from a pipe
func readPipe(pipe io.Reader) ([]byte, error) {
	return io.ReadAll(pipe)
}
