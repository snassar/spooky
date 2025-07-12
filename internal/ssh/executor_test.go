package ssh

import (
	"errors"
	"strings"
	"testing"

	"spooky/internal/config"
)

// --- Mock SSHClient for executor tests ---
type mockSSHClient struct {
	failConnect bool
	failExec    bool
	output      string
}

func (m *mockSSHClient) ExecuteCommand(cmd string) (string, error) {
	if m.failExec {
		return "", errors.New("execution failed")
	}
	return m.output, nil
}

func (m *mockSSHClient) ExecuteScript(script string) (string, error) {
	if m.failExec {
		return "", errors.New("script execution failed")
	}
	return m.output, nil
}

func (m *mockSSHClient) Close() error { return nil }

func TestExecuteConfig_EmptyActions(t *testing.T) {
	cfg := &config.Config{
		Servers: []config.Server{{
			Name:     "testserver",
			Host:     "localhost",
			User:     "testuser",
			Password: "testpass",
		}},
		Actions: []config.Action{},
	}

	err := ExecuteConfig(cfg)
	if err != nil {
		t.Errorf("expected no error for empty actions, got: %v", err)
	}
}

func TestExecuteConfig_InvalidServer(t *testing.T) {
	cfg := &config.Config{
		Servers: []config.Server{{
			Name:     "testserver",
			Host:     "localhost",
			User:     "testuser",
			Password: "testpass",
		}},
		Actions: []config.Action{{
			Name:    "testaction",
			Command: "echo hello",
			Servers: []string{"nonexistent"},
		}},
	}

	err := ExecuteConfig(cfg)
	if err == nil {
		t.Error("expected error for non-existent server")
	}
	if !strings.Contains(err.Error(), "failed to get servers for action") {
		t.Errorf("expected server error, got: %v", err)
	}
}

func TestExecuteActionSequential_EmptyServers(t *testing.T) {
	action := &config.Action{
		Name:    "testaction",
		Command: "echo hello",
	}
	servers := []*config.Server{} // Empty servers list

	err := executeActionSequential(action, servers)
	if err != nil {
		t.Errorf("expected no error for empty servers, got: %v", err)
	}
}

func TestExecuteActionParallel_EmptyServers(t *testing.T) {
	action := &config.Action{
		Name:    "testaction",
		Command: "echo hello",
	}
	servers := []*config.Server{} // Empty servers list

	err := executeActionParallel(action, servers)
	if err != nil {
		t.Errorf("expected no error for empty servers, got: %v", err)
	}
}
