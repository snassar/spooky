package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type Commander interface {
	Output(name string, args ...string) ([]byte, error)
	Run(name string, args ...string) error
}

type RealCommander struct{}

func (RealCommander) Output(name string, args ...string) ([]byte, error) {
	return exec.Command(name, args...).Output()
}
func (RealCommander) Run(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Allow test override
var stdPrintln = fmt.Println
var stdPrintf = fmt.Printf

func main() {
	if err := runPreCommitChecks(RealCommander{}); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

// runPreCommitChecks performs the pre-commit coverage checks
func runPreCommitChecks(cmdr Commander) error {
	if _, err := stdPrintln("Running pre-commit coverage checks..."); err != nil {
		return fmt.Errorf("failed to print status: %w", err)
	}

	// Check if we're in a git repository
	if _, err := cmdr.Output("git", "rev-parse", "--git-dir"); err != nil {
		errorMsg := "❌ Not in a git repository"
		if _, printErr := stdPrintln(errorMsg); printErr != nil {
			return fmt.Errorf("failed to print error: %w", printErr)
		}
		return errors.New(errorMsg)
	}

	// Get staged Go files
	output, err := cmdr.Output("git", "diff", "--cached", "--name-only", "--diff-filter=ACM")
	if err != nil {
		errorMsg := fmt.Sprintf("❌ Failed to get staged files: %v", err)
		if _, printErr := stdPrintln(errorMsg); printErr != nil {
			return fmt.Errorf("failed to print error: %w", printErr)
		}
		return errors.New(errorMsg)
	}

	stagedFiles := strings.Split(strings.TrimSpace(string(output)), "\n")
	var goFiles []string
	for _, file := range stagedFiles {
		if strings.HasSuffix(file, ".go") && !strings.HasSuffix(file, "_test.go") {
			goFiles = append(goFiles, file)
		}
	}

	if len(goFiles) == 0 {
		if _, err := stdPrintln("No Go files staged, skipping coverage check"); err != nil {
			return fmt.Errorf("failed to print status: %w", err)
		}
		return nil
	}

	if _, err := stdPrintln("Staged Go files:"); err != nil {
		return fmt.Errorf("failed to print status: %w", err)
	}
	for _, file := range goFiles {
		if _, err := stdPrintf("  %s\n", file); err != nil {
			return fmt.Errorf("failed to print file: %w", err)
		}
	}

	// Generate coverage profile
	if _, err := stdPrintln("Generating coverage profile..."); err != nil {
		return fmt.Errorf("failed to print status: %w", err)
	}
	if err := cmdr.Run("go", "test", "./...", "-coverprofile=./tests/coverage.out", "-covermode=atomic", "-coverpkg=./...", "-v"); err != nil {
		errorMsg := fmt.Sprintf("❌ Test execution failed: %v", err)
		if _, printErr := stdPrintln(errorMsg); printErr != nil {
			return fmt.Errorf("failed to print error: %w", printErr)
		}
		return errors.New(errorMsg)
	}

	// Run coverage check
	if _, err := stdPrintln("Running coverage check..."); err != nil {
		return fmt.Errorf("failed to print status: %w", err)
	}
	if err := cmdr.Run("go", "run", "github.com/vladopajic/go-test-coverage/v2@latest", "--config=./tests/testcoverage.yml"); err != nil {
		if _, printErr := stdPrintln("❌ Coverage thresholds not met"); printErr != nil {
			return fmt.Errorf("failed to print error: %w", printErr)
		}
		if _, printErr := stdPrintln("Please add tests to improve coverage before committing"); printErr != nil {
			return fmt.Errorf("failed to print error: %w", printErr)
		}
		if _, printErr := stdPrintln("Run 'make check-coverage' for detailed coverage report"); printErr != nil {
			return fmt.Errorf("failed to print error: %w", printErr)
		}
		return errors.New("coverage thresholds not met")
	}

	if _, err := stdPrintln("✅ Coverage thresholds met"); err != nil {
		return fmt.Errorf("failed to print status: %w", err)
	}
	return nil
}
