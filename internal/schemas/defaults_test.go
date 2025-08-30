package schemas

import (
	"strings"
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
	if loggingDefaults.FilePerms != "0o644" {
		t.Errorf("Expected file permissions to be '0o644', got '%s'", loggingDefaults.FilePerms)
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
	hcl, err := dcg.ToHCLWithBlockName(sshDefaults, "ssh")
	if err != nil {
		t.Errorf("Expected no error generating HCL, got %v", err)
	}

	// Check that HCL contains expected values
	if hcl == "" {
		t.Error("Expected non-empty HCL output")
	}

	// Note: print the actual HCL output for test verification
	t.Logf("Generated HCL:\n%s", hcl)

	// Check for specific values in HCL (now in block format)
	if !contains(hcl, "timeout") || !contains(hcl, "30") {
		t.Error("Expected HCL to contain 'timeout' and '30'")
	}
	if !contains(hcl, "keepalive_interval") || !contains(hcl, "60") {
		t.Error("Expected HCL to contain 'keepalive_interval' and '60'")
	}
	if !contains(hcl, "known_hosts_strict") || !contains(hcl, "true") {
		t.Error("Expected HCL to contain 'known_hosts_strict' and 'true'")
	}
	// Check for block format
	if !contains(hcl, "ssh {") {
		t.Error("Expected HCL to contain 'ssh {' block")
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

func TestDefaultConfigGenerator_ToHCL(t *testing.T) {
	dcg := NewDefaultConfigGenerator()

	// Test with SpookySSH config
	sshConfig := dcg.GetDefaultSpookySSH()
	hclOutput, err := dcg.ToHCLWithBlockName(sshConfig, "ssh")
	if err != nil {
		t.Fatalf("ToHCL failed: %v", err)
	}

	// Verify the HCL output contains expected fields
	expectedFields := []string{
		"timeout",
		"keepalive_interval",
		"keepalive_count",
		"key_scan_timeout",
		"known_hosts_strict",
		"connection_pool_size",
		"compression",
		"compression_level",
		"tcp_keepalive",
		"tcp_keepalive_count",
	}

	for _, field := range expectedFields {
		if !strings.Contains(hclOutput, field) {
			t.Errorf("HCL output missing expected field: %s", field)
		}
	}

	// Verify the output is valid HCL format
	if !strings.Contains(hclOutput, "ssh") {
		t.Error("HCL output should contain the ssh block")
	}

	t.Logf("Generated HCL:\n%s", hclOutput)
}

func TestDefaultConfigGenerator_ToHCL_ComplexStruct(t *testing.T) {
	dcg := NewDefaultConfigGenerator()

	// Test with full Spooky config
	spookyConfig := dcg.GetDefaultSpookyConfig()
	hclOutput, err := dcg.ToHCLWithBlockName(spookyConfig, "spooky")
	if err != nil {
		t.Fatalf("ToHCL failed: %v", err)
	}

	// Verify nested blocks are present
	expectedBlocks := []string{
		"ssh",
		"security",
		"age",
		"logging",
	}

	for _, block := range expectedBlocks {
		if !strings.Contains(hclOutput, block) {
			t.Errorf("HCL output missing expected block: %s", block)
		}
	}

	// Verify the output is valid HCL format
	if !strings.Contains(hclOutput, "spooky") {
		t.Error("HCL output should contain the spooky block")
	}

	t.Logf("Generated HCL:\n%s", hclOutput)
}

func TestDefaultConfigGenerator_ToHCL_EmptyStruct(t *testing.T) {
	dcg := NewDefaultConfigGenerator()

	// Test with empty struct
	type EmptyStruct struct{}
	empty := EmptyStruct{}

	hclOutput, err := dcg.ToHCLWithBlockName(empty, "empty")
	if err != nil {
		t.Fatalf("ToHCL failed: %v", err)
	}

	// Should generate an empty block
	if !strings.Contains(hclOutput, "empty") {
		t.Error("HCL output should contain the empty block")
	}

	t.Logf("Generated HCL for empty struct:\n%s", hclOutput)
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
