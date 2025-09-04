package version

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// ScalVer represents a Scalable Calendar Version following the ScalVer specification
// Format: <MAJOR>.<DATE>.<PATCH> where DATE can be YYYY, YYYYMM, or YYYYMMDD
type ScalVer struct {
	Major int
	Date  string // YYYY, YYYYMM, or YYYYMMDD
	Patch int
}

// String returns the string representation of the ScalVer
func (v ScalVer) String() string {
	return fmt.Sprintf("%d.%s.%d", v.Major, v.Date, v.Patch)
}

// IsValid checks if the ScalVer is valid according to ScalVer specification
func (v ScalVer) IsValid() bool {
	// Check if major and patch are non-negative
	if v.Major < 0 || v.Patch < 0 {
		return false
	}

	// Check if date format is valid (YYYY, YYYYMM, or YYYYMMDD)
	dateLen := len(v.Date)
	if dateLen != 4 && dateLen != 6 && dateLen != 8 {
		return false
	}

	// Validate year (must be 4 digits starting with 20 or 2)
	if dateLen >= 4 {
		year := v.Date[:4]
		if !regexp.MustCompile(`^(20\d{2}|2\d{3})$`).MatchString(year) {
			return false
		}
	}

	// Validate month (01-12)
	if dateLen >= 6 {
		month := v.Date[4:6]
		if !regexp.MustCompile(`^(0[1-9]|1[0-2])$`).MatchString(month) {
			return false
		}
	}

	// Validate day (01-31)
	if dateLen == 8 {
		day := v.Date[6:8]
		if !regexp.MustCompile(`^(0[1-9]|[12]\d|3[01])$`).MatchString(day) {
			return false
		}
	}

	return true
}

// Compare compares two ScalVer versions
// Returns -1 if v < other, 0 if v == other, 1 if v > other
func (v ScalVer) Compare(other ScalVer) int {
	// Compare major first
	if v.Major != other.Major {
		if v.Major < other.Major {
			return -1
		}
		return 1
	}

	// Compare date (lexicographically works for YYYY, YYYYMM, YYYYMMDD)
	if v.Date != other.Date {
		if v.Date < other.Date {
			return -1
		}
		return 1
	}

	// Compare patch
	if v.Patch != other.Patch {
		if v.Patch < other.Patch {
			return -1
		}
		return 1
	}

	return 0
}

// Less returns true if v is less than other
func (v ScalVer) Less(other ScalVer) bool {
	return v.Compare(other) < 0
}

// Greater returns true if v is greater than other
func (v ScalVer) Greater(other ScalVer) bool {
	return v.Compare(other) > 0
}

// Equal returns true if v equals other
func (v ScalVer) Equal(other ScalVer) bool {
	return v.Compare(other) == 0
}

// ParseScalVer parses a string into a ScalVer
// Supports SemVer pre-release identifiers (e.g., "1.20250905.0-dev-abc123")
func ParseScalVer(version string) (ScalVer, error) {
	// Remove pre-release identifiers for parsing (they don't affect the core version)
	coreVersion := version
	if dashIndex := strings.Index(version, "-"); dashIndex != -1 {
		coreVersion = version[:dashIndex]
	}

	parts := strings.Split(coreVersion, ".")
	if len(parts) != 3 {
		return ScalVer{}, fmt.Errorf("invalid ScalVer format: expected MAJOR.DATE.PATCH, got %s", version)
	}

	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return ScalVer{}, fmt.Errorf("invalid major version: %s", parts[0])
	}

	patch, err := strconv.Atoi(parts[2])
	if err != nil {
		return ScalVer{}, fmt.Errorf("invalid patch version: %s", parts[2])
	}

	date := parts[1]

	scalVer := ScalVer{
		Major: major,
		Date:  date,
		Patch: patch,
	}

	if !scalVer.IsValid() {
		return ScalVer{}, fmt.Errorf("invalid ScalVer: %s", version)
	}

	return scalVer, nil
}

// NewScalVer creates a new ScalVer with the current date
func NewScalVer(major int, patch int) ScalVer {
	now := time.Now().UTC()
	date := now.Format("20060102") // YYYYMMDD format

	return ScalVer{
		Major: major,
		Date:  date,
		Patch: patch,
	}
}

// NewScalVerWithDate creates a new ScalVer with a specific date
func NewScalVerWithDate(major int, date string, patch int) ScalVer {
	return ScalVer{
		Major: major,
		Date:  date,
		Patch: patch,
	}
}

// GetCurrentVersion returns the current version of spooky
// This can be overridden at build time using -ldflags
var GetCurrentVersion = func() ScalVer {
	// Use build-time version if available
	if buildVersion != "" {
		if version, err := ParseScalVer(buildVersion); err == nil {
			return version
		}
	}

	// Default to current date with major 0 (alpha/experimental)
	return NewScalVer(0, 0)
}

// GetVersionString returns the version string with optional build metadata
func GetVersionString() string {
	version := GetCurrentVersion()

	// Add build metadata if available
	buildInfo := getBuildInfo()
	if buildInfo != "" {
		return fmt.Sprintf("%s+%s", version.String(), buildInfo)
	}

	return version.String()
}

// GetFullVersionInfo returns detailed version information
func GetFullVersionInfo() VersionInfo {
	return VersionInfo{
		Version:   GetCurrentVersion(),
		BuildInfo: getBuildInfo(),
		BuildTime: getBuildTime(),
		GitCommit: getGitCommit(),
		GitBranch: getGitBranch(),
		GoVersion: getGoVersion(),
	}
}

// VersionInfo contains detailed version information
type VersionInfo struct {
	Version   ScalVer `json:"version"`
	BuildInfo string  `json:"build_info,omitempty"`
	BuildTime string  `json:"build_time,omitempty"`
	GitCommit string  `json:"git_commit,omitempty"`
	GitBranch string  `json:"git_branch,omitempty"`
	GoVersion string  `json:"go_version,omitempty"`
}

// String returns a formatted version info string
func (vi VersionInfo) String() string {
	var parts []string

	parts = append(parts, fmt.Sprintf("spooky %s", vi.Version.String()))
	parts = append(parts, "Automation and configuration management tool")

	if vi.BuildInfo != "" {
		parts = append(parts, fmt.Sprintf("Build: %s", vi.BuildInfo))
	}

	if vi.BuildTime != "" {
		parts = append(parts, fmt.Sprintf("Built: %s", vi.BuildTime))
	}

	if vi.GitCommit != "" {
		parts = append(parts, fmt.Sprintf("Commit: %s", vi.GitCommit))
	}

	if vi.GitBranch != "" {
		parts = append(parts, fmt.Sprintf("Branch: %s", vi.GitBranch))
	}

	if vi.GoVersion != "" {
		parts = append(parts, fmt.Sprintf("Go: %s", vi.GoVersion))
	}

	return strings.Join(parts, "\n")
}

// Build-time variables (set via -ldflags)
var (
	buildInfo    string
	buildTime    string
	gitCommit    string
	gitBranch    string
	goVersion    string
	buildVersion string // Set at build time to override GetCurrentVersion
)

// getBuildInfo returns build information if available
func getBuildInfo() string {
	if buildInfo != "" {
		return buildInfo
	}

	// Try to get git commit hash as fallback
	if gitCommit != "" {
		// Return short commit hash
		if len(gitCommit) > 7 {
			return fmt.Sprintf("dev-%s", gitCommit[:7])
		}
		return fmt.Sprintf("dev-%s", gitCommit)
	}

	return ""
}

// getBuildTime returns build time if available
func getBuildTime() string {
	return buildTime
}

// getGitCommit returns git commit hash if available
func getGitCommit() string {
	return gitCommit
}

// getGitBranch returns git branch if available
func getGitBranch() string {
	return gitBranch
}

// getGoVersion returns Go version if available
func getGoVersion() string {
	return goVersion
}

// SetBuildInfo sets build information (typically called at build time)
func SetBuildInfo(info, time, commit, branch, goVer string) {
	buildInfo = info
	buildTime = time
	gitCommit = commit
	gitBranch = branch
	goVersion = goVer
}
