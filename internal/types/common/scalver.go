package common

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// ScalVer represents a Scalable Versioning version
type ScalVer struct {
	Major int
	Date  string // YYYY, YYYYMM, or YYYYMMDD
	Patch int
}

// ParseScalVer parses a ScalVer string into its components
func ParseScalVer(version string) (*ScalVer, error) {
	if version == "" {
		return nil, fmt.Errorf("version string cannot be empty")
	}

	// Check for development suffix (e.g., 0.20241215.0-dev-abc123)
	if idx := strings.Index(version, "-dev-"); idx != -1 {
		version = version[:idx]
	}

	// Parse the main version components
	parts := strings.Split(version, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid ScalVer format: expected MAJOR.DATE.PATCH, got %d parts", len(parts))
	}

	// Parse major version
	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, fmt.Errorf("invalid major version: %s", parts[0])
	}

	// Parse date component
	date := parts[1]
	if !isValidDateComponent(date) {
		return nil, fmt.Errorf("invalid date component: %s (must be YYYY, YYYYMM, or YYYYMMDD)", date)
	}

	// Parse patch version
	patch, err := strconv.Atoi(parts[2])
	if err != nil {
		return nil, fmt.Errorf("invalid patch version: %s", parts[2])
	}

	if patch < 0 {
		return nil, fmt.Errorf("patch version cannot be negative: %d", patch)
	}

	return &ScalVer{
		Major: major,
		Date:  date,
		Patch: patch,
	}, nil
}

// String returns the ScalVer as a string
func (sv *ScalVer) String() string {
	return fmt.Sprintf("%d.%s.%d", sv.Major, sv.Date, sv.Patch)
}

// IsDevelopment returns true if this is a development version
func (sv *ScalVer) IsDevelopment() bool {
	// Development versions typically have major version 0
	return sv.Major == 0
}

// IsStable returns true if this is a stable version
func (sv *ScalVer) IsStable() bool {
	return sv.Major >= 1
}

// GetDatePrecision returns the precision of the date component
func (sv *ScalVer) GetDatePrecision() string {
	switch len(sv.Date) {
	case 4:
		return "yearly"
	case 6:
		return "monthly"
	case 8:
		return "daily"
	default:
		return "unknown"
	}
}

// Compare compares this ScalVer with another ScalVer
func (sv *ScalVer) Compare(other *ScalVer) int {
	// Compare major versions
	if sv.Major != other.Major {
		if sv.Major < other.Major {
			return -1
		}
		return 1
	}

	// Compare dates (lexicographic comparison works for YYYY, YYYYMM, YYYYMMDD)
	if sv.Date != other.Date {
		if sv.Date < other.Date {
			return -1
		}
		return 1
	}

	// Compare patch versions
	if sv.Patch != other.Patch {
		if sv.Patch < other.Patch {
			return -1
		}
		return 1
	}

	return 0
}

// LessThan returns true if this ScalVer is less than the other
func (sv *ScalVer) LessThan(other *ScalVer) bool {
	return sv.Compare(other) < 0
}

// GreaterThan returns true if this ScalVer is greater than the other
func (sv *ScalVer) GreaterThan(other *ScalVer) bool {
	return sv.Compare(other) > 0
}

// Equal returns true if this ScalVer equals the other
func (sv *ScalVer) Equal(other *ScalVer) bool {
	return sv.Compare(other) == 0
}

// isValidDateComponent checks if a date component is valid
func isValidDateComponent(date string) bool {
	if len(date) != 4 && len(date) != 6 && len(date) != 8 {
		return false
	}

	// Check that all characters are digits
	for _, char := range date {
		if char < '0' || char > '9' {
			return false
		}
	}

	// Validate the date components
	switch len(date) {
	case 4: // YYYY
		year, _ := strconv.Atoi(date)
		return year >= 1900 && year <= 2100
	case 6: // YYYYMM
		year, _ := strconv.Atoi(date[:4])
		month, _ := strconv.Atoi(date[4:])
		return year >= 1900 && year <= 2100 && month >= 1 && month <= 12
	case 8: // YYYYMMDD
		year, _ := strconv.Atoi(date[:4])
		month, _ := strconv.Atoi(date[4:6])
		day, _ := strconv.Atoi(date[6:])
		return year >= 1900 && year <= 2100 && month >= 1 && month <= 12 && day >= 1 && day <= 31
	default:
		return false
	}
}

// IsValidScalVerFormat checks if a string is in valid ScalVer format
func IsValidScalVerFormat(version string) bool {
	// Check for development suffix
	if idx := strings.Index(version, "-dev-"); idx != -1 {
		version = version[:idx]
	}

	// Parse the version to validate it
	_, err := ParseScalVer(version)
	return err == nil
}

// GenerateScalVer generates a ScalVer version string
func GenerateScalVer(major int, datePrecision string, patch int) (string, error) {
	if major < 0 {
		return "", fmt.Errorf("major version cannot be negative")
	}

	if patch < 0 {
		return "", fmt.Errorf("patch version cannot be negative")
	}

	var date string
	now := time.Now().UTC()

	switch datePrecision {
	case "yearly":
		date = now.Format("2006")
	case "monthly":
		date = now.Format("200601")
	case "daily":
		date = now.Format("20060102")
	default:
		return "", fmt.Errorf("invalid date precision: %s (must be yearly, monthly, or daily)", datePrecision)
	}

	return fmt.Sprintf("%d.%s.%d", major, date, patch), nil
}

// GenerateDevelopmentScalVer generates a development ScalVer version with git commit
func GenerateDevelopmentScalVer(gitCommit string) (string, error) {
	if gitCommit == "" {
		gitCommit = "unknown"
	}

	// Use short commit hash if it's longer than 7 characters
	if len(gitCommit) > 7 {
		gitCommit = gitCommit[:7]
	}

	baseVersion, err := GenerateScalVer(0, "daily", 0)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s-dev-%s", baseVersion, gitCommit), nil
}

// GetScalVerInfo returns detailed information about a ScalVer version
func GetScalVerInfo(version string) (map[string]interface{}, error) {
	scalver, err := ParseScalVer(version)
	if err != nil {
		return nil, err
	}

	info := map[string]interface{}{
		"version":        scalver.String(),
		"major":          scalver.Major,
		"date":           scalver.Date,
		"patch":          scalver.Patch,
		"is_development": scalver.IsDevelopment(),
		"is_stable":      scalver.IsStable(),
		"date_precision": scalver.GetDatePrecision(),
		"format":         "scalver",
	}

	// Add date information if we can parse it
	if date, err := parseDateComponent(scalver.Date); err == nil {
		info["date_parsed"] = date
		info["date_formatted"] = formatDateComponent(scalver.Date)
	}

	return info, nil
}

// parseDateComponent parses a date component into a time.Time
func parseDateComponent(date string) (time.Time, error) {
	switch len(date) {
	case 4: // YYYY
		return time.Parse("2006", date)
	case 6: // YYYYMM
		return time.Parse("200601", date)
	case 8: // YYYYMMDD
		return time.Parse("20060102", date)
	default:
		return time.Time{}, fmt.Errorf("invalid date component length: %d", len(date))
	}
}

// formatDateComponent formats a date component for display
func formatDateComponent(date string) string {
	switch len(date) {
	case 4: // YYYY
		return date
	case 6: // YYYYMM
		year := date[:4]
		month := date[4:]
		monthInt, _ := strconv.Atoi(month)
		monthName := time.Month(monthInt).String()
		return fmt.Sprintf("%s %s", monthName, year)
	case 8: // YYYYMMDD
		year := date[:4]
		month := date[4:6]
		day := date[6:]
		monthInt, _ := strconv.Atoi(month)
		monthName := time.Month(monthInt).String()
		return fmt.Sprintf("%s %s, %s", monthName, day, year)
	default:
		return date
	}
}

// ValidateScalVerCompatibility checks if two ScalVer versions are compatible
func ValidateScalVerCompatibility(version1, version2 string) (bool, error) {
	scalver1, err := ParseScalVer(version1)
	if err != nil {
		return false, fmt.Errorf("invalid first version: %w", err)
	}

	scalver2, err := ParseScalVer(version2)
	if err != nil {
		return false, fmt.Errorf("invalid second version: %w", err)
	}

	// Major version must match for compatibility
	if scalver1.Major != scalver2.Major {
		return false, nil
	}

	// Development versions are compatible with each other
	if scalver1.IsDevelopment() && scalver2.IsDevelopment() {
		return true, nil
	}

	// Stable versions are compatible if major version matches
	if scalver1.IsStable() && scalver2.IsStable() {
		return true, nil
	}

	return false, nil
}
