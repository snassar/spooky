package config

import (
	"os"
	"testing"
)

func TestParseConfig_ValidFile(t *testing.T) {
	// Create a temporary valid config file
	content := `machine "web1" {
  host = "192.168.1.10"
  user = "admin"
  password = "secret"
}`
	f, err := os.CreateTemp("", "valid_config_*.hcl")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(f.Name())
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("Failed to write test content: %v", err)
	}
	f.Close()

	cfg, err := ParseConfig(f.Name())
	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
	if cfg == nil {
		t.Error("expected config, got nil")
	}
}

func TestParseConfig_NonExistentFile(t *testing.T) {
	_, err := ParseConfig("/nonexistent/file.hcl")
	if err == nil {
		t.Error("expected error for non-existent file")
	}
}

func TestParseConfig_InvalidSyntax(t *testing.T) {
	// Create a temporary invalid config file
	content := `machine "web1" { host = "192.168.1.10" user = "admin" password = "secret"` // missing closing brace
	f, err := os.CreateTemp("", "invalid_config_*.hcl")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(f.Name())
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("Failed to write test content: %v", err)
	}
	f.Close()

	_, err = ParseConfig(f.Name())
	if err == nil {
		t.Error("expected error for invalid syntax")
	}
}
