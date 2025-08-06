package schemas

import (
	"testing"
)

func TestSimpleSchemaComposition(t *testing.T) {
	// Simple test to verify basic functionality
	composer := NewSchemaComposer()

	// Test that we can create a composer
	if composer == nil {
		t.Fatal("Failed to create schema composer")
	}

	// Test that cache is empty initially
	if len(composer.cache) != 0 {
		t.Error("Cache should be empty initially")
	}

	// Test that we can clear cache
	composer.ClearCache()
	if len(composer.cache) != 0 {
		t.Error("Cache should be empty after clear")
	}

	t.Log("✅ Basic schema composer functionality works")
}
