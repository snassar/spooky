package config

import (
	"testing"

	"spooky/internal/config/environment"
	"spooky/internal/config/loading"
	"spooky/internal/config/validation"
	"spooky/internal/logging"
	spookyconfigtypes "spooky/internal/types/config"
	loggingtypes "spooky/internal/types/logging"
)

func TestNewManager(t *testing.T) {
	// Create logger
	logger := logging.NewLogger(loggingtypes.Config{})

	// Create sub-managers
	loadingManager := loading.NewManager(&spookyconfigtypes.LoadingConfig{}, logger)
	validationManager := validation.NewManager(&spookyconfigtypes.ValidationConfig{}, logger)
	environmentManager := environment.NewManager(&spookyconfigtypes.EnvironmentConfig{}, logger)

	// Create config manager
	manager := NewManager(nil, loadingManager, validationManager, environmentManager, logger)

	if manager == nil {
		t.Fatal("Expected manager to be created, got nil")
	}
}

func TestManagerCreation(t *testing.T) {
	// Create logger
	logger := logging.NewLogger(loggingtypes.Config{})

	// Create sub-managers
	loadingManager := loading.NewManager(&spookyconfigtypes.LoadingConfig{}, logger)
	validationManager := validation.NewManager(&spookyconfigtypes.ValidationConfig{}, logger)
	environmentManager := environment.NewManager(&spookyconfigtypes.EnvironmentConfig{}, logger)

	// Create config manager
	manager := NewManager(nil, loadingManager, validationManager, environmentManager, logger)

	// Verify manager is created successfully
	if manager == nil {
		t.Fatal("Expected manager to be created, got nil")
	}
}
