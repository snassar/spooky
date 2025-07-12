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
	stdPrintln("Running pre-commit coverage checks...")

	// Check if we're in a git repository
	if _, err := cmdr.Output("git", "rev-parse", "--git-dir"); err != nil {
		errorMsg := "❌ Not in a git repository"
		stdPrintln(errorMsg)
		return errors.New(errorMsg)
	}

	// Get staged Go files
	output, err := cmdr.Output("git", "diff", "--cached", "--name-only", "--diff-filter=ACM")
	if err != nil {
		errorMsg := fmt.Sprintf("❌ Failed to get staged files: %v", err)
		stdPrintln(errorMsg)
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
		stdPrintln("No Go files staged, skipping coverage check")
		return nil
	}

	stdPrintln("Staged Go files:")
	for _, file := range goFiles {
		stdPrintf("  %s\n", file)
	}

	// Generate coverage profile
	stdPrintln("Generating coverage profile...")
	if err := cmdr.Run("go", "test", "./...", "-coverprofile=./cover.out", "-covermode=atomic", "-coverpkg=./...", "-v"); err != nil {
		errorMsg := fmt.Sprintf("❌ Test execution failed: %v", err)
		stdPrintln(errorMsg)
		return errors.New(errorMsg)
	}

	// Run coverage check
	stdPrintln("Running coverage check...")
	if err := cmdr.Run("go", "run", "github.com/vladopajic/go-test-coverage/v2@latest", "--config=./tests/testcoverage.yml"); err != nil {
		stdPrintln("❌ Coverage thresholds not met")
		stdPrintln("Please add tests to improve coverage before committing")
		stdPrintln("Run 'make check-coverage' for detailed coverage report")
		return errors.New("coverage thresholds not met")
	}

	stdPrintln("✅ Coverage thresholds met")
	return nil
}
