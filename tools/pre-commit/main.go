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
	if err := printStatus("Running pre-commit coverage checks..."); err != nil {
		return err
	}

	if err := checkGitRepository(cmdr); err != nil {
		return err
	}

	goFiles, err := getStagedGoFiles(cmdr)
	if err != nil {
		return err
	}

	if len(goFiles) == 0 {
		return printStatus("No Go files staged, skipping coverage check")
	}

	if err := printStagedFiles(goFiles); err != nil {
		return err
	}

	if err := runLinting(cmdr); err != nil {
		return err
	}

	if err := generateCoverageProfile(cmdr); err != nil {
		return err
	}

	if err := runCoverageCheck(cmdr); err != nil {
		return err
	}

	return printStatus("✅ Linting and coverage checks passed")
}

// printStatus prints a status message
func printStatus(message string) error {
	if _, err := stdPrintln(message); err != nil {
		return fmt.Errorf("failed to print status: %w", err)
	}
	return nil
}

// checkGitRepository verifies we're in a git repository
func checkGitRepository(cmdr Commander) error {
	if _, err := cmdr.Output("git", "rev-parse", "--git-dir"); err != nil {
		errorMsg := "❌ Not in a git repository"
		if _, printErr := stdPrintln(errorMsg); printErr != nil {
			return fmt.Errorf("failed to print error: %w", printErr)
		}
		return errors.New(errorMsg)
	}
	return nil
}

// getStagedGoFiles gets the list of staged Go files
func getStagedGoFiles(cmdr Commander) ([]string, error) {
	output, err := cmdr.Output("git", "diff", "--cached", "--name-only", "--diff-filter=ACM")
	if err != nil {
		errorMsg := fmt.Sprintf("❌ Failed to get staged files: %v", err)
		if _, printErr := stdPrintln(errorMsg); printErr != nil {
			return nil, fmt.Errorf("failed to print error: %w", printErr)
		}
		return nil, errors.New(errorMsg)
	}

	stagedFiles := strings.Split(strings.TrimSpace(string(output)), "\n")
	var goFiles []string
	for _, file := range stagedFiles {
		if strings.HasSuffix(file, ".go") && !strings.HasSuffix(file, "_test.go") {
			goFiles = append(goFiles, file)
		}
	}

	return goFiles, nil
}

// printStagedFiles prints the list of staged Go files
func printStagedFiles(goFiles []string) error {
	if err := printStatus("Staged Go files:"); err != nil {
		return err
	}
	for _, file := range goFiles {
		if _, err := stdPrintf("  %s\n", file); err != nil {
			return fmt.Errorf("failed to print file: %w", err)
		}
	}
	return nil
}

// runLinting runs golangci-lint
func runLinting(cmdr Commander) error {
	if err := printStatus("Running golangci-lint..."); err != nil {
		return err
	}
	if err := cmdr.Run("golangci-lint", "run"); err != nil {
		errorMsg := fmt.Sprintf("❌ Linting failed: %v", err)
		if _, printErr := stdPrintln(errorMsg); printErr != nil {
			return fmt.Errorf("failed to print error: %w", printErr)
		}
		return errors.New(errorMsg)
	}
	return nil
}

// generateCoverageProfile generates the coverage profile
func generateCoverageProfile(cmdr Commander) error {
	if err := printStatus("Generating coverage profile..."); err != nil {
		return err
	}
	if err := cmdr.Run("go", "test", "./...", "-coverprofile=./tests/coverage.out", "-covermode=atomic", "-coverpkg=./...", "-v"); err != nil {
		errorMsg := fmt.Sprintf("❌ Test run failed: %v", err)
		if _, printErr := stdPrintln(errorMsg); printErr != nil {
			return fmt.Errorf("failed to print error: %w", printErr)
		}
		return errors.New(errorMsg)
	}
	return nil
}

// runCoverageCheck runs the coverage check
func runCoverageCheck(cmdr Commander) error {
	if err := printStatus("Running coverage check..."); err != nil {
		return err
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
	return nil
}
