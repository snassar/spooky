package version

import (
	"testing"
	"time"
)

func TestScalVer_String(t *testing.T) {
	tests := []struct {
		name     string
		scalVer  ScalVer
		expected string
	}{
		{
			name:     "yearly format",
			scalVer:  ScalVer{Major: 1, Date: "2025", Patch: 0},
			expected: "1.2025.0",
		},
		{
			name:     "monthly format",
			scalVer:  ScalVer{Major: 1, Date: "202503", Patch: 2},
			expected: "1.202503.2",
		},
		{
			name:     "daily format",
			scalVer:  ScalVer{Major: 0, Date: "20250905", Patch: 0},
			expected: "0.20250905.0",
		},
		{
			name:     "alpha version",
			scalVer:  ScalVer{Major: 0, Date: "20250905", Patch: 1},
			expected: "0.20250905.1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.scalVer.String()
			if result != tt.expected {
				t.Errorf("ScalVer.String() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestScalVer_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		scalVer  ScalVer
		expected bool
	}{
		{
			name:     "valid yearly format",
			scalVer:  ScalVer{Major: 1, Date: "2025", Patch: 0},
			expected: true,
		},
		{
			name:     "valid monthly format",
			scalVer:  ScalVer{Major: 1, Date: "202503", Patch: 2},
			expected: true,
		},
		{
			name:     "valid daily format",
			scalVer:  ScalVer{Major: 0, Date: "20250905", Patch: 0},
			expected: true,
		},
		{
			name:     "invalid date format - too short",
			scalVer:  ScalVer{Major: 1, Date: "25", Patch: 0},
			expected: false,
		},
		{
			name:     "invalid date format - invalid month",
			scalVer:  ScalVer{Major: 1, Date: "202513", Patch: 0},
			expected: false,
		},
		{
			name:     "invalid date format - invalid day",
			scalVer:  ScalVer{Major: 1, Date: "20250332", Patch: 0},
			expected: false,
		},
		{
			name:     "invalid major - negative",
			scalVer:  ScalVer{Major: -1, Date: "2025", Patch: 0},
			expected: false,
		},
		{
			name:     "invalid patch - negative",
			scalVer:  ScalVer{Major: 1, Date: "2025", Patch: -1},
			expected: false,
		},
		{
			name:     "valid with zero values",
			scalVer:  ScalVer{Major: 0, Date: "2025", Patch: 0},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.scalVer.IsValid()
			if result != tt.expected {
				t.Errorf("ScalVer.IsValid() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestScalVer_Compare(t *testing.T) {
	tests := []struct {
		name     string
		v1       ScalVer
		v2       ScalVer
		expected int
	}{
		{
			name:     "equal versions",
			v1:       ScalVer{Major: 1, Date: "2025", Patch: 0},
			v2:       ScalVer{Major: 1, Date: "2025", Patch: 0},
			expected: 0,
		},
		{
			name:     "different major - v1 less",
			v1:       ScalVer{Major: 1, Date: "2025", Patch: 0},
			v2:       ScalVer{Major: 2, Date: "2025", Patch: 0},
			expected: -1,
		},
		{
			name:     "different major - v1 greater",
			v1:       ScalVer{Major: 2, Date: "2025", Patch: 0},
			v2:       ScalVer{Major: 1, Date: "2025", Patch: 0},
			expected: 1,
		},
		{
			name:     "different date - v1 less",
			v1:       ScalVer{Major: 1, Date: "2025", Patch: 0},
			v2:       ScalVer{Major: 1, Date: "2026", Patch: 0},
			expected: -1,
		},
		{
			name:     "different date - v1 greater",
			v1:       ScalVer{Major: 1, Date: "2026", Patch: 0},
			v2:       ScalVer{Major: 1, Date: "2025", Patch: 0},
			expected: 1,
		},
		{
			name:     "different patch - v1 less",
			v1:       ScalVer{Major: 1, Date: "2025", Patch: 0},
			v2:       ScalVer{Major: 1, Date: "2025", Patch: 1},
			expected: -1,
		},
		{
			name:     "different patch - v1 greater",
			v1:       ScalVer{Major: 1, Date: "2025", Patch: 1},
			v2:       ScalVer{Major: 1, Date: "2025", Patch: 0},
			expected: 1,
		},
		{
			name:     "yearly vs monthly - same year",
			v1:       ScalVer{Major: 1, Date: "2025", Patch: 0},
			v2:       ScalVer{Major: 1, Date: "202503", Patch: 0},
			expected: -1,
		},
		{
			name:     "monthly vs daily - same month",
			v1:       ScalVer{Major: 1, Date: "202503", Patch: 0},
			v2:       ScalVer{Major: 1, Date: "20250301", Patch: 0},
			expected: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.v1.Compare(tt.v2)
			if result != tt.expected {
				t.Errorf("ScalVer.Compare() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestScalVer_Less(t *testing.T) {
	tests := []struct {
		name     string
		v1       ScalVer
		v2       ScalVer
		expected bool
	}{
		{
			name:     "v1 less than v2",
			v1:       ScalVer{Major: 1, Date: "2025", Patch: 0},
			v2:       ScalVer{Major: 1, Date: "2025", Patch: 1},
			expected: true,
		},
		{
			name:     "v1 equal to v2",
			v1:       ScalVer{Major: 1, Date: "2025", Patch: 0},
			v2:       ScalVer{Major: 1, Date: "2025", Patch: 0},
			expected: false,
		},
		{
			name:     "v1 greater than v2",
			v1:       ScalVer{Major: 1, Date: "2025", Patch: 1},
			v2:       ScalVer{Major: 1, Date: "2025", Patch: 0},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.v1.Less(tt.v2)
			if result != tt.expected {
				t.Errorf("ScalVer.Less() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestScalVer_Greater(t *testing.T) {
	tests := []struct {
		name     string
		v1       ScalVer
		v2       ScalVer
		expected bool
	}{
		{
			name:     "v1 greater than v2",
			v1:       ScalVer{Major: 1, Date: "2025", Patch: 1},
			v2:       ScalVer{Major: 1, Date: "2025", Patch: 0},
			expected: true,
		},
		{
			name:     "v1 equal to v2",
			v1:       ScalVer{Major: 1, Date: "2025", Patch: 0},
			v2:       ScalVer{Major: 1, Date: "2025", Patch: 0},
			expected: false,
		},
		{
			name:     "v1 less than v2",
			v1:       ScalVer{Major: 1, Date: "2025", Patch: 0},
			v2:       ScalVer{Major: 1, Date: "2025", Patch: 1},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.v1.Greater(tt.v2)
			if result != tt.expected {
				t.Errorf("ScalVer.Greater() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestScalVer_Equal(t *testing.T) {
	tests := []struct {
		name     string
		v1       ScalVer
		v2       ScalVer
		expected bool
	}{
		{
			name:     "equal versions",
			v1:       ScalVer{Major: 1, Date: "2025", Patch: 0},
			v2:       ScalVer{Major: 1, Date: "2025", Patch: 0},
			expected: true,
		},
		{
			name:     "different major",
			v1:       ScalVer{Major: 1, Date: "2025", Patch: 0},
			v2:       ScalVer{Major: 2, Date: "2025", Patch: 0},
			expected: false,
		},
		{
			name:     "different date",
			v1:       ScalVer{Major: 1, Date: "2025", Patch: 0},
			v2:       ScalVer{Major: 1, Date: "2026", Patch: 0},
			expected: false,
		},
		{
			name:     "different patch",
			v1:       ScalVer{Major: 1, Date: "2025", Patch: 0},
			v2:       ScalVer{Major: 1, Date: "2025", Patch: 1},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.v1.Equal(tt.v2)
			if result != tt.expected {
				t.Errorf("ScalVer.Equal() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestParseScalVer(t *testing.T) {
	tests := []struct {
		name        string
		version     string
		expected    ScalVer
		expectError bool
	}{
		{
			name:        "valid yearly format",
			version:     "1.2025.0",
			expected:    ScalVer{Major: 1, Date: "2025", Patch: 0},
			expectError: false,
		},
		{
			name:        "valid monthly format",
			version:     "1.202503.2",
			expected:    ScalVer{Major: 1, Date: "202503", Patch: 2},
			expectError: false,
		},
		{
			name:        "valid daily format",
			version:     "0.20250905.0",
			expected:    ScalVer{Major: 0, Date: "20250905", Patch: 0},
			expectError: false,
		},
		{
			name:        "invalid format - too few parts",
			version:     "1.2025",
			expected:    ScalVer{},
			expectError: true,
		},
		{
			name:        "invalid format - too many parts",
			version:     "1.2025.0.1",
			expected:    ScalVer{},
			expectError: true,
		},
		{
			name:        "invalid major - not a number",
			version:     "a.2025.0",
			expected:    ScalVer{},
			expectError: true,
		},
		{
			name:        "invalid patch - not a number",
			version:     "1.2025.a",
			expected:    ScalVer{},
			expectError: true,
		},
		{
			name:        "invalid date format",
			version:     "1.25.0",
			expected:    ScalVer{},
			expectError: true,
		},
		{
			name:        "valid with pre-release identifier",
			version:     "0.20250905.0-dev-abc123",
			expected:    ScalVer{Major: 0, Date: "20250905", Patch: 0},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseScalVer(tt.version)
			if tt.expectError {
				if err == nil {
					t.Errorf("ParseScalVer() expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("ParseScalVer() unexpected error: %v", err)
				}
				if result != tt.expected {
					t.Errorf("ParseScalVer() = %v, want %v", result, tt.expected)
				}
			}
		})
	}
}

func TestNewScalVer(t *testing.T) {
	// Test that NewScalVer creates a version with current date
	version := NewScalVer(1, 5)

	if version.Major != 1 {
		t.Errorf("NewScalVer() Major = %v, want 1", version.Major)
	}
	if version.Patch != 5 {
		t.Errorf("NewScalVer() Patch = %v, want 5", version.Patch)
	}

	// Check that date is in YYYYMMDD format
	expectedDate := time.Now().UTC().Format("20060102")
	if version.Date != expectedDate {
		t.Errorf("NewScalVer() Date = %v, want %v", version.Date, expectedDate)
	}

	// Check that the version is valid
	if !version.IsValid() {
		t.Errorf("NewScalVer() created invalid version: %v", version)
	}
}

func TestNewScalVerWithDate(t *testing.T) {
	tests := []struct {
		name     string
		major    int
		date     string
		patch    int
		expected ScalVer
	}{
		{
			name:     "yearly format",
			major:    1,
			date:     "2025",
			patch:    0,
			expected: ScalVer{Major: 1, Date: "2025", Patch: 0},
		},
		{
			name:     "monthly format",
			major:    1,
			date:     "202503",
			patch:    2,
			expected: ScalVer{Major: 1, Date: "202503", Patch: 2},
		},
		{
			name:     "daily format",
			major:    0,
			date:     "20250905",
			patch:    0,
			expected: ScalVer{Major: 0, Date: "20250905", Patch: 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewScalVerWithDate(tt.major, tt.date, tt.patch)
			if result != tt.expected {
				t.Errorf("NewScalVerWithDate() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGetVersionString(t *testing.T) {
	// Test default version (should not have build info)
	version := GetVersionString()
	if version == "" {
		t.Errorf("GetVersionString() returned empty string")
	}

	// The version should be in ScalVer format
	_, err := ParseScalVer(version)
	if err != nil {
		t.Errorf("GetVersionString() returned invalid ScalVer format: %v", err)
	}
}

func TestGetFullVersionInfo(t *testing.T) {
	info := GetFullVersionInfo()

	// Check that we get a valid version
	if !info.Version.IsValid() {
		t.Errorf("GetFullVersionInfo() returned invalid version: %v", info.Version)
	}

	// Check that the string representation is not empty
	infoStr := info.String()
	if infoStr == "" {
		t.Errorf("GetFullVersionInfo().String() returned empty string")
	}

	// Should contain "spooky" and the version
	if !contains(infoStr, "spooky") {
		t.Errorf("GetFullVersionInfo().String() should contain 'spooky', got: %s", infoStr)
	}
}

func TestSetBuildInfo(t *testing.T) {
	// Save original values
	originalBuildInfo := buildInfo
	originalBuildTime := buildTime
	originalGitCommit := gitCommit
	originalGitBranch := gitBranch
	originalGoVersion := goVersion

	// Set test values
	SetBuildInfo("test-build", "2025-01-01T00:00:00Z", "abc123", "main", "go1.21.0")

	// Check that values were set
	if buildInfo != "test-build" {
		t.Errorf("SetBuildInfo() buildInfo = %v, want test-build", buildInfo)
	}
	if buildTime != "2025-01-01T00:00:00Z" {
		t.Errorf("SetBuildInfo() buildTime = %v, want 2025-01-01T00:00:00Z", buildTime)
	}
	if gitCommit != "abc123" {
		t.Errorf("SetBuildInfo() gitCommit = %v, want abc123", gitCommit)
	}
	if gitBranch != "main" {
		t.Errorf("SetBuildInfo() gitBranch = %v, want main", gitBranch)
	}
	if goVersion != "go1.21.0" {
		t.Errorf("SetBuildInfo() goVersion = %v, want go1.21.0", goVersion)
	}

	// Restore original values
	SetBuildInfo(originalBuildInfo, originalBuildTime, originalGitCommit, originalGitBranch, originalGoVersion)
}

func TestGetBuildInfo(t *testing.T) {
	// Test with no build info
	originalBuildInfo := buildInfo
	originalGitCommit := gitCommit

	SetBuildInfo("", "", "", "", "")
	info := getBuildInfo()
	if info != "" {
		t.Errorf("getBuildInfo() with no build info = %v, want empty string", info)
	}

	// Test with build info
	SetBuildInfo("dev", "", "", "", "")
	info = getBuildInfo()
	if info != "dev" {
		t.Errorf("getBuildInfo() with build info = %v, want dev", info)
	}

	// Test with git commit fallback
	SetBuildInfo("", "", "abcdef123456", "", "")
	info = getBuildInfo()
	if info != "dev-abcdef1" {
		t.Errorf("getBuildInfo() with git commit = %v, want dev-abcdef1", info)
	}

	// Restore original values
	SetBuildInfo(originalBuildInfo, "", originalGitCommit, "", "")
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || (len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || containsInMiddle(s, substr))))
}

func containsInMiddle(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
