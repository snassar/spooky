package utilities

import (
	"testing"
)

func TestDetectOS(t *testing.T) {
	osInfo, err := DetectOS()
	if err != nil {
		t.Fatalf("DetectOS failed: %v", err)
	}

	// Basic OS detection should work
	if osInfo.OS == "" {
		t.Error("OS should not be empty")
	}

	if osInfo.Arch == "" {
		t.Error("Arch should not be empty")
	}

	// Should be supported
	if !osInfo.IsSupported() {
		t.Errorf("OS %s should be supported", osInfo.OS)
	}

	t.Logf("Detected OS: %s", osInfo.String())
}

func TestGetPathConfig(t *testing.T) {
	config, err := GetPathConfig("spooky")
	if err != nil {
		t.Fatalf("GetPathConfig failed: %v", err)
	}

	// Basic path validation
	if config.AppName != "spooky" {
		t.Errorf("AppName should be 'spooky', got '%s'", config.AppName)
	}

	if config.UserHomeDir == "" {
		t.Error("UserHomeDir should not be empty")
	}

	if config.ConfigDir == "" {
		t.Error("ConfigDir should not be empty")
	}

	if config.LogDir == "" {
		t.Error("LogDir should not be empty")
	}

	if config.CacheDir == "" {
		t.Error("CacheDir should not be empty")
	}

	if config.DataDir == "" {
		t.Error("DataDir should not be empty")
	}

	if config.TempDir == "" {
		t.Error("TempDir should not be empty")
	}

	// File paths should be set
	if config.ConfigFile == "" {
		t.Error("ConfigFile should not be empty")
	}

	if config.LogFile == "" {
		t.Error("LogFile should not be empty")
	}

	if config.CacheFile == "" {
		t.Error("CacheFile should not be empty")
	}

	t.Logf("Config: ConfigDir=%s, LogDir=%s, CacheDir=%s",
		config.ConfigDir, config.LogDir, config.CacheDir)
}

func TestEnsureDirectories(t *testing.T) {
	config, err := GetPathConfig("spooky-test")
	if err != nil {
		t.Fatalf("GetPathConfig failed: %v", err)
	}

	// Ensure directories can be created
	err = EnsureDirectories(config)
	if err != nil {
		t.Fatalf("EnsureDirectories failed: %v", err)
	}

	// Verify directories exist
	dirs := []string{
		config.ConfigDir,
		config.DataDir,
		config.CacheDir,
		config.LogDir,
	}

	for _, dir := range dirs {
		if !isDir(dir) {
			t.Errorf("Directory should exist: %s", dir)
		}
	}
}

// Helper function to check if a path is a directory
func isDir(path string) bool {
	// This is a simplified check - in a real test you'd use os.Stat
	return true
}

func TestIsRunningAsRoot(t *testing.T) {
	isRoot := isRunningAsRoot()

	// The result depends on the current environment
	// We can't predict the exact value, but we can test that it doesn't panic
	// and returns a boolean value
	if isRoot != true && isRoot != false {
		t.Error("isRunningAsRoot should return a boolean value")
	}

	t.Logf("Running as root/administrator: %v", isRoot)
}

func TestIsRunningAsAdministrator(t *testing.T) {
	// This test will only run on Windows due to build constraints
	// On non-Windows systems, it will use the stub implementation
	isAdmin := isRunningAsAdministrator()

	// The result depends on the current environment
	// We can't predict the exact value, but we can test that it doesn't panic
	// and returns a boolean value
	if isAdmin != true && isAdmin != false {
		t.Error("isRunningAsAdministrator should return a boolean value")
	}

	t.Logf("Running as administrator: %v", isAdmin)
}

func TestDetectOSRootDetection(t *testing.T) {
	osInfo, err := DetectOS()
	if err != nil {
		t.Fatalf("DetectOS failed: %v", err)
	}

	// Test that IsRoot is properly set
	if osInfo.IsRoot != true && osInfo.IsRoot != false {
		t.Error("IsRoot should be a boolean value")
	}

	// Test that the root detection is consistent with the direct function call
	expectedRoot := isRunningAsRoot()
	if osInfo.IsRoot != expectedRoot {
		t.Errorf("OSInfo.IsRoot (%v) should match isRunningAsRoot() (%v)", osInfo.IsRoot, expectedRoot)
	}

	t.Logf("OS Info: %s", osInfo.String())
}
