// Merged test file: contains all tests from main_test.go and pre-commit_test.go
package main

import (
	"bytes"
	"errors"
	"os"
	"strings"
	"testing"
)

type MockCommander struct {
	OutputFunc func(name string, args ...string) ([]byte, error)
	RunFunc    func(name string, args ...string) error
	RunLog     *bytes.Buffer     // for capturing output in tests
	output     map[string][]byte // for main_test.go style
	runErr     map[string]error  // for main_test.go style
}

func (m MockCommander) Output(name string, args ...string) ([]byte, error) {
	if m.OutputFunc != nil {
		return m.OutputFunc(name, args...)
	}
	key := name + " " + joinArgs(args)
	if m.output != nil {
		if out, ok := m.output[key]; ok {
			return out, nil
		}
	}
	return nil, errors.New("mock output not found")
}

func (m MockCommander) Run(name string, args ...string) error {
	if m.RunFunc != nil {
		if m.RunLog != nil {
			m.RunLog.WriteString(name + " " + strings.Join(args, " ") + "\n")
		}
		return m.RunFunc(name, args...)
	}
	key := name + " " + joinArgs(args)
	if m.runErr != nil {
		if err, ok := m.runErr[key]; ok {
			return err
		}
	}
	return nil
}

func joinArgs(args []string) string {
	return " " + join(args)
}

func join(args []string) string {
	result := ""
	for _, a := range args {
		result += a + " "
	}
	return result
}

// captureOutput captures stdout/stderr during function execution
func captureOutput(fn func()) string {
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	fn()
	w.Close()
	os.Stdout = oldStdout
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(r); err != nil {
		panic(err)
	}
	return buf.String()
}

// --- All tests from both files below ---

func TestGitRepositoryDetection(t *testing.T) {
	t.Run("Valid git repository", func(t *testing.T) {
		mock := MockCommander{
			OutputFunc: func(name string, args ...string) ([]byte, error) {
				if name == "git" && args[0] == "rev-parse" {
					return []byte(".git"), nil
				}
				if name == "git" && args[0] == "diff" {
					return []byte("main.go"), nil
				}
				return nil, errors.New("unexpected command")
			},
			RunFunc: func(_ string, _ ...string) error { return nil },
		}
		output := captureOutput(func() {
			_ = runPreCommitChecks(mock) // ignore error, output is captured for assertions
		})
		if !strings.Contains(output, "Running pre-commit coverage checks...") {
			t.Errorf("Expected output to contain 'Running pre-commit coverage checks...', got %q", output)
		}
		if !strings.Contains(output, "✅ Linting and coverage checks passed") {
			t.Errorf("Expected output to contain '✅ Linting and coverage checks passed', got %q", output)
		}
	})

	t.Run("Not in git repository", func(t *testing.T) {
		mock := MockCommander{
			OutputFunc: func(_ string, _ ...string) ([]byte, error) {
				return nil, errors.New("not a git repo")
			},
		}
		output := captureOutput(func() {
			_ = runPreCommitChecks(mock) // ignore error, output is captured for assertions
		})
		if !strings.Contains(output, "❌ Not in a git repository") {
			t.Errorf("Expected output to contain '❌ Not in a git repository', got %q", output)
		}
	})
}

func TestStagedFilesDetection(t *testing.T) {
	t.Run("No staged files", func(t *testing.T) {
		mock := MockCommander{
			OutputFunc: func(name string, args ...string) ([]byte, error) {
				if name == "git" && args[0] == "rev-parse" {
					return []byte(".git"), nil
				}
				if name == "git" && args[0] == "diff" {
					return []byte(""), nil
				}
				return nil, errors.New("unexpected command")
			},
		}
		output := captureOutput(func() {
			_ = runPreCommitChecks(mock) // ignore error, output is captured for assertions
		})
		if !strings.Contains(output, "No Go files staged, skipping coverage check") {
			t.Errorf("Expected output to contain 'No Go files staged, skipping coverage check', got %q", output)
		}
	})

	t.Run("Staged Go files", func(t *testing.T) {
		mock := MockCommander{
			OutputFunc: func(name string, args ...string) ([]byte, error) {
				if name == "git" && args[0] == "rev-parse" {
					return []byte(".git"), nil
				}
				if name == "git" && args[0] == "diff" {
					return []byte("main.go\ninternal/cli/commands.go\nssh.go\nREADME.md"), nil
				}
				return nil, errors.New("unexpected command")
			},
			RunFunc: func(_ string, _ ...string) error { return nil },
		}
		output := captureOutput(func() {
			_ = runPreCommitChecks(mock) // ignore error, output is captured for assertions
		})
		for _, file := range []string{"main.go", "internal/cli/commands.go", "ssh.go"} {
			if !strings.Contains(output, file) {
				t.Errorf("Expected output to contain staged file %q", file)
			}
		}
		if !strings.Contains(output, "Staged Go files:") {
			t.Errorf("Expected output to contain 'Staged Go files:', got %q", output)
		}
		if !strings.Contains(output, "✅ Linting and coverage checks passed") {
			t.Errorf("Expected output to contain '✅ Linting and coverage checks passed', got %q", output)
		}
	})

	t.Run("Staged files with test files", func(t *testing.T) {
		mock := MockCommander{
			OutputFunc: func(name string, args ...string) ([]byte, error) {
				if name == "git" && args[0] == "rev-parse" {
					return []byte(".git"), nil
				}
				if name == "git" && args[0] == "diff" {
					return []byte("main.go\ncommands_test.go\nssh.go\nconfig_test.go"), nil
				}
				return nil, errors.New("unexpected command")
			},
			RunFunc: func(_ string, _ ...string) error { return nil },
		}
		output := captureOutput(func() {
			_ = runPreCommitChecks(mock) // ignore error, output is captured for assertions
		})
		for _, file := range []string{"main.go", "ssh.go"} {
			if !strings.Contains(output, file) {
				t.Errorf("Expected output to contain staged file %q", file)
			}
		}
		if !strings.Contains(output, "Staged Go files:") {
			t.Errorf("Expected output to contain 'Staged Go files:', got %q", output)
		}
		if !strings.Contains(output, "✅ Linting and coverage checks passed") {
			t.Errorf("Expected output to contain '✅ Linting and coverage checks passed', got %q", output)
		}
	})

	t.Run("Git command fails", func(t *testing.T) {
		mock := MockCommander{
			OutputFunc: func(name string, args ...string) ([]byte, error) {
				if name == "git" && args[0] == "rev-parse" {
					return []byte(".git"), nil
				}
				if name == "git" && args[0] == "diff" {
					return nil, errors.New("fail")
				}
				return nil, errors.New("unexpected command")
			},
		}
		output := captureOutput(func() {
			_ = runPreCommitChecks(mock) // ignore error, output is captured for assertions
		})
		if !strings.Contains(output, "❌ Failed to get staged files:") {
			t.Errorf("Expected output to contain '❌ Failed to get staged files:', got %q", output)
		}
	})
}

func TestCoverageGeneration(t *testing.T) {
	t.Run("Tests pass", func(t *testing.T) {
		mock := MockCommander{
			OutputFunc: func(name string, args ...string) ([]byte, error) {
				if name == "git" && args[0] == "rev-parse" {
					return []byte(".git"), nil
				}
				if name == "git" && args[0] == "diff" {
					return []byte("main.go"), nil
				}
				return nil, errors.New("unexpected command")
			},
			RunFunc: func(_ string, _ ...string) error { return nil },
		}
		output := captureOutput(func() {
			_ = runPreCommitChecks(mock) // ignore error, output is captured for assertions
		})
		if !strings.Contains(output, "Generating coverage profile...") {
			t.Errorf("Expected output to contain 'Generating coverage profile...', got %q", output)
		}
		if !strings.Contains(output, "✅ Linting and coverage checks passed") {
			t.Errorf("Expected output to contain '✅ Linting and coverage checks passed', got %q", output)
		}
	})

	t.Run("Tests fail", func(t *testing.T) {
		mock := MockCommander{
			OutputFunc: func(name string, args ...string) ([]byte, error) {
				if name == "git" && args[0] == "rev-parse" {
					return []byte(".git"), nil
				}
				if name == "git" && args[0] == "diff" {
					return []byte("main.go"), nil
				}
				return nil, errors.New("unexpected command")
			},
			RunFunc: func(name string, args ...string) error {
				if name == "go" && args[0] == "test" {
					return errors.New("test execution failed")
				}
				return nil
			},
		}
		output := captureOutput(func() {
			_ = runPreCommitChecks(mock) // ignore error, output is captured for assertions
		})
		if !strings.Contains(output, "❌ Test execution failed:") {
			t.Errorf("Expected output to contain '❌ Test execution failed:', got %q", output)
		}
	})
}

func TestCoverageCheck(t *testing.T) {
	t.Run("Coverage passes", func(t *testing.T) {
		mock := MockCommander{
			OutputFunc: func(name string, args ...string) ([]byte, error) {
				if name == "git" && args[0] == "rev-parse" {
					return []byte(".git"), nil
				}
				if name == "git" && args[0] == "diff" {
					return []byte("main.go"), nil
				}
				return nil, errors.New("unexpected command")
			},
			RunFunc: func(_ string, _ ...string) error { return nil },
		}
		output := captureOutput(func() {
			_ = runPreCommitChecks(mock) // ignore error, output is captured for assertions
		})
		if !strings.Contains(output, "✅ Linting and coverage checks passed") {
			t.Errorf("Expected output to contain '✅ Linting and coverage checks passed', got %q", output)
		}
	})

	t.Run("Coverage fails", func(t *testing.T) {
		mock := MockCommander{
			OutputFunc: func(name string, args ...string) ([]byte, error) {
				if name == "git" && args[0] == "rev-parse" {
					return []byte(".git"), nil
				}
				if name == "git" && args[0] == "diff" {
					return []byte("main.go"), nil
				}
				return nil, errors.New("unexpected command")
			},
			RunFunc: func(name string, args ...string) error {
				if name == "go" && args[0] == "test" {
					return nil
				}
				if name == "go" && args[0] == "run" {
					return errors.New("coverage thresholds not met")
				}
				return nil
			},
		}
		output := captureOutput(func() {
			_ = runPreCommitChecks(mock) // ignore error, output is captured for assertions
		})
		if !strings.Contains(output, "❌ Coverage thresholds not met") {
			t.Errorf("Expected output to contain '❌ Coverage thresholds not met', got %q", output)
		}
		if !strings.Contains(output, "Please add tests to improve coverage before committing") {
			t.Errorf("Expected output to contain help message, got %q", output)
		}
	})
}

func TestRunPreCommitChecks_Error(t *testing.T) {
	cmd := &MockCommander{runErr: map[string]error{"go test ./... -coverprofile=./tests/coverage.out -covermode=atomic -coverpkg=./...": errors.New("fail")}}
	err := runPreCommitChecks(cmd)
	if err == nil {
		t.Error("Expected error from runPreCommitChecks")
	}
}
