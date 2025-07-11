package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	fmt.Println("Running pre-commit coverage checks...")

	// Check if we're in a git repository
	if _, err := exec.Command("git", "rev-parse", "--git-dir").Output(); err != nil {
		fmt.Println("❌ Not in a git repository")
		os.Exit(1)
	}

	// Get staged Go files
	cmd := exec.Command("git", "diff", "--cached", "--name-only", "--diff-filter=ACM")
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("❌ Failed to get staged files: %v\n", err)
		os.Exit(1)
	}

	stagedFiles := strings.Split(strings.TrimSpace(string(output)), "\n")
	var goFiles []string
	for _, file := range stagedFiles {
		if strings.HasSuffix(file, ".go") && !strings.HasSuffix(file, "_test.go") {
			goFiles = append(goFiles, file)
		}
	}

	if len(goFiles) == 0 {
		fmt.Println("No Go files staged, skipping coverage check")
		os.Exit(0)
	}

	fmt.Println("Staged Go files:")
	for _, file := range goFiles {
		fmt.Printf("  %s\n", file)
	}

	// Generate coverage profile
	fmt.Println("Generating coverage profile...")
	coverageCmd := exec.Command("go", "test", "./...", "-coverprofile=./cover.out", "-covermode=atomic", "-coverpkg=./...", "-v")
	coverageCmd.Stdout = os.Stdout
	coverageCmd.Stderr = os.Stderr
	if err := coverageCmd.Run(); err != nil {
		fmt.Printf("❌ Test execution failed: %v\n", err)
		os.Exit(1)
	}

	// Run coverage check
	fmt.Println("Running coverage check...")
	checkCmd := exec.Command("go", "run", "github.com/vladopajic/go-test-coverage/v2@latest", "--config=./tests/testcoverage.yml")
	checkCmd.Stdout = os.Stdout
	checkCmd.Stderr = os.Stderr
	if err := checkCmd.Run(); err != nil {
		fmt.Println("❌ Coverage thresholds not met")
		fmt.Println("Please add tests to improve coverage before committing")
		fmt.Println("Run 'make check-coverage' for detailed coverage report")
		os.Exit(1)
	}

	fmt.Println("✅ Coverage thresholds met")
}
