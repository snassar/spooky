package main

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"testing"
)

type MockCommander struct {
	OutputFunc func(name string, args ...string) ([]byte, error)
	RunFunc    func(name string, args ...string) error
	RunLog     *bytes.Buffer // for capturing output in tests
}

func (m MockCommander) Output(name string, args ...string) ([]byte, error) {
	return m.OutputFunc(name, args...)
}
func (m MockCommander) Run(name string, args ...string) error {
	if m.RunFunc != nil {
		return m.RunFunc(name, args...)
	}
	if m.RunLog != nil {
		fmt.Fprintf(m.RunLog, "Run: %s %s\n", name, strings.Join(args, " "))
	}
	return nil
}

// Helper to capture output from runPreCommitChecks
func captureOutput(fn func()) string {
	var buf bytes.Buffer
	stdPrintln = func(a ...interface{}) (int, error) {
		return fmt.Fprintln(&buf, a...)
	}
	stdPrintf = func(format string, a ...interface{}) (int, error) {
		return fmt.Fprintf(&buf, format, a...)
	}
	defer func() {
		stdPrintln = fmt.Println
		stdPrintf = fmt.Printf
	}()
	fn()
	return buf.String()
}

func TestGitRepositoryDetection(t *testing.T) {
	t.Run("Valid git repository", func(t *testing.T) {
		mock := MockCommander{
			OutputFunc: func(name string, args ...string) ([]byte, error) {
				if name == "git" && args[0] == "rev-parse" {
					return []byte(".git"), nil
				}
				if name == "git" && args[0] == "diff" {
					return []byte("main.go\ncommands.go"), nil
				}
				return nil, errors.New("unexpected command")
			},
			RunFunc: func(name string, args ...string) error { return nil },
		}
		output := captureOutput(func() {
			runPreCommitChecks(mock)
		})
		if !strings.Contains(output, "Running pre-commit coverage checks...") {
			t.Errorf("Expected output to contain 'Running pre-commit coverage checks...', got %q", output)
		}
		if !strings.Contains(output, "✅ Coverage thresholds met") {
			t.Errorf("Expected output to contain '✅ Coverage thresholds met', got %q", output)
		}
	})

	t.Run("Not in git repository", func(t *testing.T) {
		mock := MockCommander{
			OutputFunc: func(name string, args ...string) ([]byte, error) {
				return nil, errors.New("not a git repo")
			},
		}
		output := captureOutput(func() {
			runPreCommitChecks(mock)
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
			runPreCommitChecks(mock)
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
					return []byte("main.go\ncommands.go\nssh.go\nREADME.md"), nil
				}
				return nil, errors.New("unexpected command")
			},
			RunFunc: func(name string, args ...string) error { return nil },
		}
		output := captureOutput(func() {
			runPreCommitChecks(mock)
		})
		for _, file := range []string{"main.go", "commands.go", "ssh.go"} {
			if !strings.Contains(output, file) {
				t.Errorf("Expected output to contain staged file %q", file)
			}
		}
		if !strings.Contains(output, "Staged Go files:") {
			t.Errorf("Expected output to contain 'Staged Go files:', got %q", output)
		}
		if !strings.Contains(output, "✅ Coverage thresholds met") {
			t.Errorf("Expected output to contain '✅ Coverage thresholds met', got %q", output)
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
			RunFunc: func(name string, args ...string) error { return nil },
		}
		output := captureOutput(func() {
			runPreCommitChecks(mock)
		})
		for _, file := range []string{"main.go", "ssh.go"} {
			if !strings.Contains(output, file) {
				t.Errorf("Expected output to contain staged file %q", file)
			}
		}
		if !strings.Contains(output, "Staged Go files:") {
			t.Errorf("Expected output to contain 'Staged Go files:', got %q", output)
		}
		if !strings.Contains(output, "✅ Coverage thresholds met") {
			t.Errorf("Expected output to contain '✅ Coverage thresholds met', got %q", output)
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
			runPreCommitChecks(mock)
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
			RunFunc: func(name string, args ...string) error { return nil },
		}
		output := captureOutput(func() {
			runPreCommitChecks(mock)
		})
		if !strings.Contains(output, "Generating coverage profile...") {
			t.Errorf("Expected output to contain 'Generating coverage profile...', got %q", output)
		}
		if !strings.Contains(output, "✅ Coverage thresholds met") {
			t.Errorf("Expected output to contain '✅ Coverage thresholds met', got %q", output)
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
					return errors.New("fail")
				}
				return nil
			},
		}
		output := captureOutput(func() {
			runPreCommitChecks(mock)
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
			RunFunc: func(name string, args ...string) error { return nil },
		}
		output := captureOutput(func() {
			runPreCommitChecks(mock)
		})
		if !strings.Contains(output, "✅ Coverage thresholds met") {
			t.Errorf("Expected output to contain '✅ Coverage thresholds met', got %q", output)
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
				if name == "go" && args[0] == "run" {
					return errors.New("fail")
				}
				return nil
			},
		}
		output := captureOutput(func() {
			runPreCommitChecks(mock)
		})
		if !strings.Contains(output, "❌ Coverage thresholds not met") {
			t.Errorf("Expected output to contain '❌ Coverage thresholds not met', got %q", output)
		}
	})
}
