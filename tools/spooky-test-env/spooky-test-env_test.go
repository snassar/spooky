package main

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

type mockCommander struct {
	output map[string][]byte
	runErr map[string]error
}

func (m *mockCommander) Output(name string, args ...string) ([]byte, error) {
	key := name + " " + joinArgs(args)
	if out, ok := m.output[key]; ok {
		return out, nil
	}
	return nil, errors.New("mock output not found")
}

func (m *mockCommander) Run(name string, args ...string) error {
	key := name + " " + joinArgs(args)
	if err, ok := m.runErr[key]; ok {
		return err
	}
	return nil
}

func joinArgs(args []string) string {
	return " " + filepath.Join(args...)
}

func TestGetTestEnvDir(t *testing.T) {
	cwd, _ := os.Getwd()
	dir := getTestEnvDir()
	if dir == "" || dir == "." {
		t.Errorf("Expected a non-empty test env dir, got %q", dir)
	}
	if !filepath.IsAbs(dir) && dir != filepath.Join(cwd, "..", "..", "spooky-test-env") {
		t.Errorf("Unexpected test env dir: %q", dir)
	}
}

func TestGetContainerIP_Error(t *testing.T) {
	cmd := &mockCommander{output: map[string][]byte{}, runErr: map[string]error{}}
	_, err := getContainerIPWithCmd(cmd, "nonexistent")
	if err == nil {
		t.Error("Expected error for nonexistent container")
	}
}

func getContainerIPWithCmd(cmd TestEnvCommander, containerName string) (string, error) {
	output, err := cmd.Output("podman", "inspect", containerName, "--format", "{{.NetworkSettings.Networks.spooky-test.IPAddress}}")
	if err != nil {
		return "", err
	}
	return string(output), nil
}
