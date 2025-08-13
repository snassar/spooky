package ssh

import (
	"context"
	"testing"

	spookylogging "spooky/internal/logging"
	spookytypes "spooky/internal/types"
	spookytypesssh "spooky/internal/types/ssh"
)

func TestNewManager(t *testing.T) {
	// Create a logger
	logManager := spookylogging.NewLogManager()
	logger := logManager.GetLogger("test")

	// Create SSH manager
	manager := NewManager(logger)
	if manager == nil {
		t.Fatal("Expected SSH manager to be created, got nil")
	}
}

func TestManager_ValidateConnection(t *testing.T) {
	// Create a logger
	logManager := spookylogging.NewLogManager()
	logger := logManager.GetLogger("test")

	// Create SSH manager
	manager := NewManager(logger)

	// Test valid connection request
	validRequest := &spookytypes.ConnectionRequest{
		Host:     "localhost",
		Port:     22,
		User:     "testuser",
		KeyPath:  "/path/to/key",
		Password: "",
	}

	result, err := manager.ValidateConnection(context.Background(), validRequest)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if !result.Valid {
		t.Fatalf("Expected validation to pass, got errors: %v", result.Errors)
	}

	// Test invalid connection request (missing host)
	invalidRequest := &spookytypes.ConnectionRequest{
		Host:     "",
		Port:     22,
		User:     "testuser",
		KeyPath:  "/path/to/key",
		Password: "",
	}

	result, err = manager.ValidateConnection(context.Background(), invalidRequest)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result.Valid {
		t.Fatal("Expected validation to fail, but it passed")
	}

	if len(result.Errors) == 0 {
		t.Fatal("Expected validation errors, got none")
	}
}

func TestManager_ValidateAuthentication(t *testing.T) {
	// Create a logger
	logManager := spookylogging.NewLogManager()
	logger := logManager.GetLogger("test")

	// Create SSH manager
	manager := NewManager(logger)

	// Test valid authentication (public key)
	validAuth := &spookytypes.Authentication{
		Method:  spookytypesssh.AuthMethodPublicKey,
		KeyPath: "/path/to/key",
	}

	result, err := manager.ValidateAuthentication(context.Background(), validAuth)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if !result.Valid {
		t.Fatalf("Expected validation to pass, got errors: %v", result.Errors)
	}

	// Test invalid authentication (missing key path)
	invalidAuth := &spookytypes.Authentication{
		Method:  spookytypesssh.AuthMethodPublicKey,
		KeyPath: "",
	}

	result, err = manager.ValidateAuthentication(context.Background(), invalidAuth)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result.Valid {
		t.Fatal("Expected validation to fail, but it passed")
	}

	if len(result.Errors) == 0 {
		t.Fatal("Expected validation errors, got none")
	}
}
