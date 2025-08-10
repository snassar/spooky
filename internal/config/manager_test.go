package config

import (
	"testing"

	"spooky/internal/config/environment"
	"spooky/internal/config/loading"
	"spooky/internal/types/config"
	"spooky/internal/config/validation"
	"spooky/internal/logging"
	loggingtypes "spooky/internal/types/logging"
)

func TestNewManager(t *testing.T) {
	// Create logger
	logger := logging.NewLogger(loggingtypes.Config{})

	// Create sub-managers
	loadingManager := loading.NewManager(&types.LoadingConfig{}, logger)
	validationManager := validation.NewManager(&types.ValidationConfig{}, logger)
	environmentManager := environment.NewManager(&types.EnvironmentConfig{}, logger)

	// Create config manager
	manager := NewManager(nil, loadingManager, validationManager, environmentManager, logger)

	if manager == nil {
		t.Fatal("Expected manager to be created, got nil")
	}
}

func TestManagerImplementsConfigManager(t *testing.T) {
	// Create logger
	logger := logging.NewLogger(loggingtypes.Config{})

	// Create sub-managers
	loadingManager := loading.NewManager(&types.LoadingConfig{}, logger)
	validationManager := validation.NewManager(&types.ValidationConfig{}, logger)
	environmentManager := environment.NewManager(&types.EnvironmentConfig{}, logger)

	// Create config manager
	manager := NewManager(nil, loadingManager, validationManager, environmentManager, logger)

	// Verify it implements the interface
	var _ ConfigManager = manager
}
