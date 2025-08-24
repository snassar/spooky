package schemas

import (
	"testing"
)

func TestDefaultConfigGenerator(t *testing.T) {
	dcg := NewDefaultConfigGenerator()

	// Test Spooky SSH defaults
	sshDefaults := dcg.GetDefaultSpookySSH()
	if sshDefaults.Timeout != 30 {
		t.Errorf("Expected SSH timeout to be 30, got %d", sshDefaults.Timeout)
	}
	if sshDefaults.KeepaliveInterval != 60 {
		t.Errorf("Expected SSH keepalive interval to be 60, got %d", sshDefaults.KeepaliveInterval)
	}
	if sshDefaults.KeepaliveCount != 3 {
		t.Errorf("Expected SSH keepalive count to be 3, got %d", sshDefaults.KeepaliveCount)
	}
	if sshDefaults.KeyScanTimeout != 10 {
		t.Errorf("Expected SSH key scan timeout to be 10, got %d", sshDefaults.KeyScanTimeout)
	}
	if !sshDefaults.KnownHostsStrict {
		t.Errorf("Expected SSH known hosts strict to be true, got %t", sshDefaults.KnownHostsStrict)
	}
	if sshDefaults.ConnectionPoolSize != 10 {
		t.Errorf("Expected SSH connection pool size to be 10, got %d", sshDefaults.ConnectionPoolSize)
	}

	// Test Spooky Security defaults
	securityDefaults := dcg.GetDefaultSpookySecurity()
	if securityDefaults.AllowUnsafeCommands {
		t.Errorf("Expected allow unsafe commands to be false, got %t", securityDefaults.AllowUnsafeCommands)
	}
	if !securityDefaults.AuditLogging {
		t.Errorf("Expected audit logging to be true, got %t", securityDefaults.AuditLogging)
	}

	// Test Spooky Logging defaults
	loggingDefaults := dcg.GetDefaultSpookyLogging()
	if loggingDefaults.Level != "info" {
		t.Errorf("Expected log level to be 'info', got '%s'", loggingDefaults.Level)
	}
	if loggingDefaults.Format != "json" {
		t.Errorf("Expected log format to be 'json', got '%s'", loggingDefaults.Format)
	}
	if loggingDefaults.Output != "stderr" {
		t.Errorf("Expected log output to be 'stderr', got '%s'", loggingDefaults.Output)
	}
	if loggingDefaults.FilePerms != "0644" {
		t.Errorf("Expected file permissions to be '0644', got '%s'", loggingDefaults.FilePerms)
	}
	if !loggingDefaults.FileAppend {
		t.Errorf("Expected file append to be true, got %t", loggingDefaults.FileAppend)
	}

	// Test Project defaults
	projectDefaults := dcg.GetDefaultProjectConfig()
	if projectDefaults.RunMaxParallel != 10 {
		t.Errorf("Expected run max parallel to be 10, got %d", projectDefaults.RunMaxParallel)
	}
	if projectDefaults.RunDryRunDefault {
		t.Errorf("Expected run dry run default to be false, got %t", projectDefaults.RunDryRunDefault)
	}
	if !projectDefaults.RunValidateBeforeRun {
		t.Errorf("Expected run validate before run to be true, got %t", projectDefaults.RunValidateBeforeRun)
	}
	if projectDefaults.RunBackupBeforeChanges {
		t.Errorf("Expected run backup before changes to be false, got %t", projectDefaults.RunBackupBeforeChanges)
	}
	if projectDefaults.FactsTimeout != 30 {
		t.Errorf("Expected facts timeout to be 30, got %d", projectDefaults.FactsTimeout)
	}
	if projectDefaults.FactsParallelCollection != 10 {
		t.Errorf("Expected facts parallel collection to be 10, got %d", projectDefaults.FactsParallelCollection)
	}
	if projectDefaults.FactsRetryAttempts != 3 {
		t.Errorf("Expected facts retry attempts to be 3, got %d", projectDefaults.FactsRetryAttempts)
	}
	if projectDefaults.FactsRetryDelay != 5 {
		t.Errorf("Expected facts retry delay to be 5, got %d", projectDefaults.FactsRetryDelay)
	}
}

func TestDefaultConfigGeneratorHCL(t *testing.T) {
	dcg := NewDefaultConfigGenerator()

	// Test HCL generation
	sshDefaults := dcg.GetDefaultSpookySSH()
	hcl, err := dcg.ToHCL(sshDefaults)
	if err != nil {
		t.Errorf("Expected no error generating HCL, got %v", err)
	}

	// Check that HCL contains expected values
	if hcl == "" {
		t.Error("Expected non-empty HCL output")
	}

	// Check for specific values in HCL
	if !contains(hcl, "timeout = 30") {
		t.Error("Expected HCL to contain 'timeout = 30'")
	}
	if !contains(hcl, "keepalive_interval = 60") {
		t.Error("Expected HCL to contain 'keepalive_interval = 60'")
	}
	if !contains(hcl, "known_hosts_strict = true") {
		t.Error("Expected HCL to contain 'known_hosts_strict = true'")
	}
}

func TestDefaultConfigGeneratorComplete(t *testing.T) {
	dcg := NewDefaultConfigGenerator()

	// Test complete Spooky config
	spookyDefaults := dcg.GetDefaultSpookyConfig()

	// Verify SSH section
	if spookyDefaults.SSH.Timeout != 30 {
		t.Errorf("Expected SSH timeout to be 30, got %d", spookyDefaults.SSH.Timeout)
	}

	// Verify Security section
	if !spookyDefaults.Security.AuditLogging {
		t.Errorf("Expected audit logging to be true, got %t", spookyDefaults.Security.AuditLogging)
	}

	// Verify Logging section
	if spookyDefaults.Logging.Level != "info" {
		t.Errorf("Expected log level to be 'info', got '%s'", spookyDefaults.Logging.Level)
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		(len(s) > len(substr) && (s[:len(substr)] == substr ||
			s[len(s)-len(substr):] == substr ||
			containsSubstring(s, substr))))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
