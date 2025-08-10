package templates

import (
	"context"
	"testing"
	"time"

	"spooky/internal/logging"
	loggingTypes "spooky/internal/types/logging"
	"spooky/internal/templates/engine"
	"spooky/internal/templates/functions"
	"spooky/internal/templates/secrets"
	"spooky/internal/templates/types"
	"spooky/internal/templates/validation"
)

// createTestLogger creates a logger for testing
func createTestLogger() logging.Logger {
	return logging.NewLogger(loggingTypes.Config{
		Level:     loggingTypes.InfoLevel,
		Format:    "text",
		Output:    "stdout",
		Timestamp: false,
	})
}

func TestNewManager(t *testing.T) {
	logger := createTestLogger()
	config := &types.Config{
		DefaultTimeout:  30,
		MaxTemplateSize: 1024 * 1024,
	}

	engineManager := engine.NewManager(config.EngineConfig, logger)
	functionsManager := functions.NewManager(config.FunctionsConfig, logger)
	validationManager := validation.NewManager(config.ValidationConfig, logger)
	secretsManager := secrets.NewManager(config.SecretsConfig, logger)

	manager := NewManager(
		config,
		engineManager,
		functionsManager,
		validationManager,
		secretsManager,
		logger,
	)

	if manager == nil {
		t.Fatal("Expected manager to be created, got nil")
	}
}

func TestNewTemplateManager(t *testing.T) {
	logger := createTestLogger()
	config := &types.Config{
		DefaultTimeout:  30,
		MaxTemplateSize: 1024 * 1024,
	}

	manager, err := NewTemplateManager(config, logger)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if manager == nil {
		t.Fatal("Expected manager to be created, got nil")
	}
}

func TestNewDefaultTemplateManager(t *testing.T) {
	logger := createTestLogger()

	manager, err := NewDefaultTemplateManager(logger)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if manager == nil {
		t.Fatal("Expected manager to be created, got nil")
	}
}

func TestManager_NewTemplateContext(t *testing.T) {
	logger := createTestLogger()
	manager, err := NewDefaultTemplateManager(logger)
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}

	projectPath := "/test/project"
	context, err := manager.NewTemplateContext(projectPath)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if context == nil {
		t.Fatal("Expected context to be created, got nil")
	}

	if context.ProjectPath != projectPath {
		t.Errorf("Expected project path %s, got %s", projectPath, context.ProjectPath)
	}

	if context.Data == nil {
		t.Fatal("Expected data map to be initialized")
	}

	if context.Functions == nil {
		t.Fatal("Expected functions map to be initialized")
	}
}

func TestManager_GetTemplateFunctions(t *testing.T) {
	logger := createTestLogger()
	manager, err := NewDefaultTemplateManager(logger)
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}

	context := &types.TemplateContext{
		ProjectPath: "/test/project",
		Data:        make(map[string]interface{}),
		Functions:   make(map[string]interface{}),
	}

	functions := manager.GetTemplateFunctions(context)
	if functions == nil {
		t.Fatal("Expected functions to be returned, got nil")
	}

	// Check for some expected built-in functions
	expectedFunctions := []string{"upper", "lower", "trim", "custom", "system", "env"}
	for _, fnName := range expectedFunctions {
		if functions[fnName] == nil {
			t.Errorf("Expected function %s to be available", fnName)
		}
	}
}

func TestManager_SetDefaultTimeout(t *testing.T) {
	logger := createTestLogger()
	manager, err := NewDefaultTemplateManager(logger)
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}

	timeout := 60 * time.Second
	manager.SetDefaultTimeout(timeout)

	// Note: In a real implementation, we would verify the timeout was set
	// For now, we just ensure the method doesn't panic
}

func TestManager_SetMaxTemplateSize(t *testing.T) {
	logger := createTestLogger()
	manager, err := NewDefaultTemplateManager(logger)
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}

	maxSize := int64(2 * 1024 * 1024) // 2MB
	manager.SetMaxTemplateSize(maxSize)

	// Note: In a real implementation, we would verify the size was set
	// For now, we just ensure the method doesn't panic
}

func TestManager_RegisterCustomFunction(t *testing.T) {
	logger := createTestLogger()
	manager, err := NewDefaultTemplateManager(logger)
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}

	// Test registering a valid function
	customFn := func(s string) string { return "custom_" + s }
	err = manager.RegisterCustomFunction("custom", customFn)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Test registering an invalid function (nil)
	err = manager.RegisterCustomFunction("invalid", nil)
	if err == nil {
		t.Error("Expected error for nil function, got nil")
	}
}

func TestManager_ValidateTemplates(t *testing.T) {
	logger := createTestLogger()
	manager, err := NewDefaultTemplateManager(logger)
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}

	errors, err := manager.ValidateTemplates("/test/project")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// errors can be nil if no templates are found or no validation errors occur
	if errors == nil {
		// This is acceptable - no validation errors
		return
	}
}

func TestManager_Close(t *testing.T) {
	logger := createTestLogger()
	manager, err := NewDefaultTemplateManager(logger)
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}

	err = manager.Close()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestManager_RenderTemplateWithTimeout(t *testing.T) {
	logger := createTestLogger()
	manager, err := NewDefaultTemplateManager(logger)
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}

	ctx := context.Background()
	additionalData := map[string]interface{}{
		"test": "value",
	}

	// This would fail because we don't have a real template file
	// But we can test that the method exists and doesn't panic
	_, err = manager.RenderTemplateWithTimeout(ctx, "nonexistent.tmpl", "/test/project", additionalData)
	if err == nil {
		t.Error("Expected error for nonexistent template, got nil")
	}
}
