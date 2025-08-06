package ssh

import (
	"testing"
	"time"

	"spooky/internal/logging"
)

func TestNewDefaultManager(t *testing.T) {
	logger := logging.GetLogger()
	manager := NewDefaultManager(logger)

	if manager == nil {
		t.Fatal("Expected manager to be created, got nil")
	}

	if manager.config == nil {
		t.Fatal("Expected config to be set, got nil")
	}

	if manager.config.DefaultTimeout != 30*time.Second {
		t.Errorf("Expected default timeout to be 30s, got %v", manager.config.DefaultTimeout)
	}

	if manager.config.MaxConnections != 10 {
		t.Errorf("Expected max connections to be 10, got %d", manager.config.MaxConnections)
	}
}

func TestManager_SetDefaultTimeout(t *testing.T) {
	logger := logging.GetLogger()
	manager := NewDefaultManager(logger)

	timeout := 60 * time.Second
	err := manager.SetDefaultTimeout(timeout)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if manager.config.DefaultTimeout != timeout {
		t.Errorf("Expected timeout to be %v, got %v", timeout, manager.config.DefaultTimeout)
	}
}

func TestManager_SetMaxConnections(t *testing.T) {
	logger := logging.GetLogger()
	manager := NewDefaultManager(logger)

	maxConnections := 20
	err := manager.SetMaxConnections(maxConnections)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if manager.config.MaxConnections != maxConnections {
		t.Errorf("Expected max connections to be %d, got %d", maxConnections, manager.config.MaxConnections)
	}
}

func TestManager_EnableConnectionPooling(t *testing.T) {
	logger := logging.GetLogger()
	manager := NewDefaultManager(logger)

	err := manager.EnableConnectionPooling(false)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if manager.config.EnableConnectionPooling {
		t.Error("Expected connection pooling to be disabled")
	}
}
