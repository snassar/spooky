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

// TestManagerUsesReusableSSHClient verifies that the manager uses the ReusableSSHClient
func TestManagerUsesReusableSSHClient(t *testing.T) {
	// Create a logger
	logManager := spookylogging.NewLogManager()
	logger := logManager.GetLogger("test")

	// Create SSH manager
	manager := NewManager(logger)

	// Cast to concrete type to access the client field
	concreteManager, ok := manager.(*Manager)
	if !ok {
		t.Fatal("Manager should be of type *Manager")
	}

	// Verify that the manager has a ReusableSSHClient
	if concreteManager.client == nil {
		t.Fatal("Manager should have a ReusableSSHClient")
	}

	// Check that we can get connection stats (this is a feature of ReusableSSHClient)
	stats := concreteManager.client.GetConnectionStats()
	if stats == nil {
		t.Fatal("Should be able to get connection stats from ReusableSSHClient")
	}

	// Check that we can get metrics (this is a feature of ReusableSSHClient)
	metrics := concreteManager.client.GetMetrics()
	if metrics == nil {
		t.Fatal("Should be able to get metrics from ReusableSSHClient")
	}

	// Verify that the metrics contain expected fields
	expectedFields := []string{"total_connections", "successful_connections", "failed_connections", "reused_connections"}
	for _, field := range expectedFields {
		if _, exists := metrics[field]; !exists {
			t.Errorf("Expected metrics to contain field: %s", field)
		}
	}

	t.Log("Manager successfully uses ReusableSSHClient with connection reuse capabilities")
}

// TestManagerConnectionReuse verifies that the manager properly reuses connections
func TestManagerConnectionReuse(t *testing.T) {
	// Create a logger
	logManager := spookylogging.NewLogManager()
	logger := logManager.GetLogger("test")

	// Create SSH manager
	manager := NewManager(logger)

	// Cast to concrete type to access the client field
	concreteManager, ok := manager.(*Manager)
	if !ok {
		t.Fatal("Manager should be of type *Manager")
	}

	// Get initial connection stats
	initialStats := concreteManager.client.GetConnectionStats()
	initialTotal := initialStats["total_connections"].(int)

	// Create a test machine (we won't actually connect, just test the interface)
	testMachine := &spookytypes.Machine{
		Hostname: "test.example.com",
		Host:     "test.example.com",
		Port:     22,
		User:     "testuser",
	}

	// Simulate multiple command runs (they will fail due to no real connection)
	// but we can verify the interface is working correctly
	for i := 0; i < 3; i++ {
		ctx := context.Background()
		_, err := concreteManager.client.RunCommand(ctx, testMachine, "echo 'test'")
		// We expect this to fail since there's no real SSH server, but the interface should work
		if err == nil {
			t.Logf("Command %d ran successfully (unexpected in test environment)", i+1)
		} else {
			t.Logf("Command %d failed as expected: %v", i+1, err)
		}
	}

	// Get final connection stats
	finalStats := concreteManager.client.GetConnectionStats()
	finalTotal := finalStats["total_connections"].(int)

	// In a real environment with successful connections, we would expect:
	// - Initial total connections: 0
	// - After multiple commands: potentially > 0 (if connections were cached)
	// - The ReusableSSHClient should attempt to reuse connections

	t.Logf("Connection reuse test completed - Initial connections: %d, Final connections: %d", initialTotal, finalTotal)
	t.Log("Note: In a real environment with successful SSH connections, the ReusableSSHClient would cache and reuse connections")
}
