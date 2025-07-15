package ssh

import (
	"strings"
	"testing"

	"spooky/internal/config"
)

func TestExecuteConfig_EmptyActions(t *testing.T) {
	cfg := &config.Config{
		Servers: []config.Server{
			{Name: "test", Host: "localhost", User: "testuser", Password: "testpass"},
		},
		Actions: []config.Action{},
	}

	err := ExecuteConfig(cfg)
	if err != nil {
		t.Errorf("expected no error for empty actions, got: %v", err)
	}
}

func TestExecuteConfig_InvalidServer(t *testing.T) {
	cfg := &config.Config{
		Servers: []config.Server{
			{Name: "test", Host: "localhost", User: "testuser", Password: "testpass"},
		},
		Actions: []config.Action{
			{
				Name:    "testaction",
				Command: "echo test",
				Servers: []string{"nonexistent"},
			},
		},
	}

	err := ExecuteConfig(cfg)
	if err == nil {
		t.Error("expected error when server does not exist")
	}
	if !strings.Contains(err.Error(), "server 'nonexistent' not found") {
		t.Errorf("expected server not found error, got: %v", err)
	}
}

func TestExecuteConfig_ValidConfig(t *testing.T) {
	cfg := &config.Config{
		Servers: []config.Server{
			{Name: "test", Host: "localhost", User: "testuser", Password: "testpass"},
		},
		Actions: []config.Action{
			{
				Name:    "testaction",
				Command: "echo test",
			},
		},
	}

	// This will fail to connect but should not fail due to configuration issues
	err := ExecuteConfig(cfg)
	if err != nil && !strings.Contains(err.Error(), "connection refused") && !strings.Contains(err.Error(), "no route to host") {
		t.Errorf("expected connection error, got: %v", err)
	}
}

func TestExecuteConfig_ActionWithDescription(t *testing.T) {
	cfg := &config.Config{
		Servers: []config.Server{
			{Name: "test", Host: "localhost", User: "testuser", Password: "testpass"},
		},
		Actions: []config.Action{
			{
				Name:        "testaction",
				Description: "Test action description",
				Command:     "echo test",
			},
		},
	}

	// This will fail to connect but should not fail due to configuration issues
	err := ExecuteConfig(cfg)
	if err != nil && !strings.Contains(err.Error(), "connection refused") && !strings.Contains(err.Error(), "no route to host") {
		t.Errorf("expected connection error, got: %v", err)
	}
}

func TestExecuteConfig_ActionWithScript(t *testing.T) {
	cfg := &config.Config{
		Servers: []config.Server{
			{Name: "test", Host: "localhost", User: "testuser", Password: "testpass"},
		},
		Actions: []config.Action{
			{
				Name:   "testaction",
				Script: "/nonexistent/script.sh",
			},
		},
	}

	// This will fail to connect but should not fail due to configuration issues
	err := ExecuteConfig(cfg)
	if err != nil && !strings.Contains(err.Error(), "connection refused") && !strings.Contains(err.Error(), "no route to host") {
		t.Errorf("expected connection error, got: %v", err)
	}
}

func TestExecuteConfig_ActionWithBothCommandAndScript(t *testing.T) {
	cfg := &config.Config{
		Servers: []config.Server{
			{Name: "test", Host: "localhost", User: "testuser", Password: "testpass"},
		},
		Actions: []config.Action{
			{
				Name:    "testaction",
				Command: "echo test",
				Script:  "/nonexistent/script.sh",
			},
		},
	}

	// This should fail due to validation error (both command and script specified)
	err := ExecuteConfig(cfg)
	if err == nil {
		t.Error("expected error when both command and script are specified")
	}
}

func TestExecuteActionSequential_EmptyServers(t *testing.T) {
	action := &config.Action{
		Name:    "testaction",
		Command: "echo test",
	}

	err := executeActionSequential(action, []*config.Server{})
	if err != nil {
		t.Errorf("expected no error for empty servers list, got: %v", err)
	}
}

func TestExecuteActionSequential_SingleServer(t *testing.T) {
	action := &config.Action{
		Name:    "testaction",
		Command: "echo test",
	}

	servers := []*config.Server{
		{Name: "test", Host: "localhost", User: "testuser", Password: "testpass"},
	}

	// This will fail to connect but should not fail due to configuration issues
	err := executeActionSequential(action, servers)
	if err != nil && !strings.Contains(err.Error(), "connection refused") && !strings.Contains(err.Error(), "no route to host") {
		t.Errorf("expected connection error, got: %v", err)
	}
}

func TestExecuteActionSequential_MultipleServers(t *testing.T) {
	action := &config.Action{
		Name:    "testaction",
		Command: "echo test",
	}

	servers := []*config.Server{
		{Name: "test1", Host: "localhost", User: "testuser", Password: "testpass"},
		{Name: "test2", Host: "localhost", User: "testuser", Password: "testpass"},
	}

	// This will fail to connect but should not fail due to configuration issues
	err := executeActionSequential(action, servers)
	if err != nil && !strings.Contains(err.Error(), "connection refused") && !strings.Contains(err.Error(), "no route to host") {
		t.Errorf("expected connection error, got: %v", err)
	}
}

func TestExecuteActionSequential_ServerWithKeyFile(t *testing.T) {
	action := &config.Action{
		Name:    "testaction",
		Command: "echo test",
	}

	servers := []*config.Server{
		{Name: "test", Host: "localhost", User: "testuser", KeyFile: "/nonexistent/key"},
	}

	// This will fail due to key file not found
	err := executeActionSequential(action, servers)
	if err != nil && !strings.Contains(err.Error(), "failed to read key file") {
		t.Errorf("expected key file error, got: %v", err)
	}
}

func TestExecuteActionParallel_EmptyServers(t *testing.T) {
	action := &config.Action{
		Name:     "testaction",
		Command:  "echo test",
		Parallel: true,
	}

	err := executeActionParallel(action, []*config.Server{})
	if err != nil {
		t.Errorf("expected no error for empty servers list, got: %v", err)
	}
}

func TestExecuteActionParallel_SingleServer(t *testing.T) {
	action := &config.Action{
		Name:     "testaction",
		Command:  "echo test",
		Parallel: true,
	}

	servers := []*config.Server{
		{Name: "test", Host: "localhost", User: "testuser", Password: "testpass"},
	}

	// This will fail to connect but should not fail due to configuration issues
	err := executeActionParallel(action, servers)
	if err != nil && !strings.Contains(err.Error(), "connection refused") && !strings.Contains(err.Error(), "no route to host") {
		t.Errorf("expected connection error, got: %v", err)
	}
}

func TestExecuteActionParallel_MultipleServers(t *testing.T) {
	action := &config.Action{
		Name:     "testaction",
		Command:  "echo test",
		Parallel: true,
	}

	servers := []*config.Server{
		{Name: "test1", Host: "localhost", User: "testuser", Password: "testpass"},
		{Name: "test2", Host: "localhost", User: "testuser", Password: "testpass"},
		{Name: "test3", Host: "localhost", User: "testuser", Password: "testpass"},
	}

	// This will fail to connect but should not fail due to configuration issues
	err := executeActionParallel(action, servers)
	if err != nil && !strings.Contains(err.Error(), "connection refused") && !strings.Contains(err.Error(), "no route to host") {
		t.Errorf("expected connection error, got: %v", err)
	}
}

func TestExecuteActionParallel_ServerWithKeyFile(t *testing.T) {
	action := &config.Action{
		Name:     "testaction",
		Command:  "echo test",
		Parallel: true,
	}

	servers := []*config.Server{
		{Name: "test", Host: "localhost", User: "testuser", KeyFile: "/nonexistent/key"},
	}

	// This will fail due to key file not found
	err := executeActionParallel(action, servers)
	if err != nil && !strings.Contains(err.Error(), "failed to read key file") {
		t.Errorf("expected key file error, got: %v", err)
	}
}

func TestExecuteActionSequential_ActionWithScript(t *testing.T) {
	action := &config.Action{
		Name:   "testaction",
		Script: "/nonexistent/script.sh",
	}

	servers := []*config.Server{
		{Name: "test", Host: "localhost", User: "testuser", Password: "testpass"},
	}

	// This will fail due to script file not found
	err := executeActionSequential(action, servers)
	if err != nil && !strings.Contains(err.Error(), "failed to read script file") {
		t.Errorf("expected script file error, got: %v", err)
	}
}

func TestExecuteActionParallel_ActionWithScript(t *testing.T) {
	action := &config.Action{
		Name:     "testaction",
		Script:   "/nonexistent/script.sh",
		Parallel: true,
	}

	servers := []*config.Server{
		{Name: "test", Host: "localhost", User: "testuser", Password: "testpass"},
	}

	// This will fail due to script file not found
	err := executeActionParallel(action, servers)
	if err != nil && !strings.Contains(err.Error(), "failed to read script file") {
		t.Errorf("expected script file error, got: %v", err)
	}
}

func TestExecuteActionSequential_ActionWithNoCommandOrScript(t *testing.T) {
	action := &config.Action{
		Name: "testaction",
		// No command or script specified
	}

	servers := []*config.Server{
		{Name: "test", Host: "localhost", User: "testuser", Password: "testpass"},
	}

	// This should fail due to validation error
	err := executeActionSequential(action, servers)
	if err == nil {
		t.Error("expected error when neither command nor script is specified")
	}
}

func TestExecuteActionParallel_ActionWithNoCommandOrScript(t *testing.T) {
	action := &config.Action{
		Name:     "testaction",
		Parallel: true,
		// No command or script specified
	}

	servers := []*config.Server{
		{Name: "test", Host: "localhost", User: "testuser", Password: "testpass"},
	}

	// This should fail due to validation error
	err := executeActionParallel(action, servers)
	if err == nil {
		t.Error("expected error when neither command nor script is specified")
	}
}
