package common

import (
	"strings"
	"testing"
)

func TestParseScalVer(t *testing.T) {
	tests := []struct {
		name    string
		version string
		want    *ScalVer
		wantErr bool
	}{
		{
			name:    "valid yearly format",
			version: "0.2025.0",
			want: &ScalVer{
				Major: 0,
				Date:  "2025",
				Patch: 0,
			},
			wantErr: false,
		},
		{
			name:    "valid monthly format",
			version: "0.202508.0",
			want: &ScalVer{
				Major: 0,
				Date:  "202508",
				Patch: 0,
			},
			wantErr: false,
		},
		{
			name:    "valid daily format",
			version: "0.20250812.0",
			want: &ScalVer{
				Major: 0,
				Date:  "20250812",
				Patch: 0,
			},
			wantErr: false,
		},
		{
			name:    "valid development format",
			version: "0.20250812.0-dev-abc123",
			want: &ScalVer{
				Major: 0,
				Date:  "20250812",
				Patch: 0,
			},
			wantErr: false,
		},
		{
			name:    "valid stable version",
			version: "1.2025.0",
			want: &ScalVer{
				Major: 1,
				Date:  "2025",
				Patch: 0,
			},
			wantErr: false,
		},
		{
			name:    "valid with patch",
			version: "0.2025.5",
			want: &ScalVer{
				Major: 0,
				Date:  "2025",
				Patch: 5,
			},
			wantErr: false,
		},
		{
			name:    "empty version",
			version: "",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalid format - too few parts",
			version: "0.2025",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalid format - too many parts",
			version: "0.2025.0.1",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalid major version",
			version: "abc.2025.0",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalid patch version",
			version: "0.2025.abc",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "negative patch version",
			version: "0.2025.-1",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalid date format",
			version: "0.202.0",
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseScalVer(tt.version)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseScalVer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("ParseScalVer() returned nil for valid version")
				return
			}
			if tt.wantErr && got != nil {
				t.Errorf("ParseScalVer() returned non-nil for invalid version")
				return
			}
			if got != nil && tt.want != nil {
				if got.Major != tt.want.Major {
					t.Errorf("ParseScalVer() Major = %v, want %v", got.Major, tt.want.Major)
				}
				if got.Date != tt.want.Date {
					t.Errorf("ParseScalVer() Date = %v, want %v", got.Date, tt.want.Date)
				}
				if got.Patch != tt.want.Patch {
					t.Errorf("ParseScalVer() Patch = %v, want %v", got.Patch, tt.want.Patch)
				}
			}
		})
	}
}

func TestScalVer_String(t *testing.T) {
	tests := []struct {
		name string
		sv   *ScalVer
		want string
	}{
		{
			name: "yearly format",
			sv: &ScalVer{
				Major: 0,
				Date:  "2025",
				Patch: 0,
			},
			want: "0.2025.0",
		},
		{
			name: "monthly format",
			sv: &ScalVer{
				Major: 0,
				Date:  "202508",
				Patch: 0,
			},
			want: "0.202508.0",
		},
		{
			name: "daily format",
			sv: &ScalVer{
				Major: 0,
				Date:  "20250812",
				Patch: 0,
			},
			want: "0.20250812.0",
		},
		{
			name: "with patch",
			sv: &ScalVer{
				Major: 1,
				Date:  "2025",
				Patch: 5,
			},
			want: "1.2025.5",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.sv.String(); got != tt.want {
				t.Errorf("ScalVer.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestScalVer_IsDevelopment(t *testing.T) {
	tests := []struct {
		name string
		sv   *ScalVer
		want bool
	}{
		{
			name: "development version",
			sv: &ScalVer{
				Major: 0,
				Date:  "2025",
				Patch: 0,
			},
			want: true,
		},
		{
			name: "stable version",
			sv: &ScalVer{
				Major: 1,
				Date:  "2025",
				Patch: 0,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.sv.IsDevelopment(); got != tt.want {
				t.Errorf("ScalVer.IsDevelopment() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestScalVer_IsStable(t *testing.T) {
	tests := []struct {
		name string
		sv   *ScalVer
		want bool
	}{
		{
			name: "development version",
			sv: &ScalVer{
				Major: 0,
				Date:  "2025",
				Patch: 0,
			},
			want: false,
		},
		{
			name: "stable version",
			sv: &ScalVer{
				Major: 1,
				Date:  "2025",
				Patch: 0,
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.sv.IsStable(); got != tt.want {
				t.Errorf("ScalVer.IsStable() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestScalVer_GetDatePrecision(t *testing.T) {
	tests := []struct {
		name string
		sv   *ScalVer
		want string
	}{
		{
			name: "yearly precision",
			sv: &ScalVer{
				Major: 0,
				Date:  "2025",
				Patch: 0,
			},
			want: "yearly",
		},
		{
			name: "monthly precision",
			sv: &ScalVer{
				Major: 0,
				Date:  "202508",
				Patch: 0,
			},
			want: "monthly",
		},
		{
			name: "daily precision",
			sv: &ScalVer{
				Major: 0,
				Date:  "20250812",
				Patch: 0,
			},
			want: "daily",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.sv.GetDatePrecision(); got != tt.want {
				t.Errorf("ScalVer.GetDatePrecision() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestScalVer_Compare(t *testing.T) {
	tests := []struct {
		name string
		sv1  *ScalVer
		sv2  *ScalVer
		want int
	}{
		{
			name: "equal versions",
			sv1:  &ScalVer{Major: 0, Date: "2025", Patch: 0},
			sv2:  &ScalVer{Major: 0, Date: "2025", Patch: 0},
			want: 0,
		},
		{
			name: "sv1 less than sv2 - major",
			sv1:  &ScalVer{Major: 0, Date: "2025", Patch: 0},
			sv2:  &ScalVer{Major: 1, Date: "2025", Patch: 0},
			want: -1,
		},
		{
			name: "sv1 greater than sv2 - major",
			sv1:  &ScalVer{Major: 1, Date: "2025", Patch: 0},
			sv2:  &ScalVer{Major: 0, Date: "2025", Patch: 0},
			want: 1,
		},
		{
			name: "sv1 less than sv2 - date",
			sv1:  &ScalVer{Major: 0, Date: "2025", Patch: 0},
			sv2:  &ScalVer{Major: 0, Date: "2026", Patch: 0},
			want: -1,
		},
		{
			name: "sv1 greater than sv2 - date",
			sv1:  &ScalVer{Major: 0, Date: "2026", Patch: 0},
			sv2:  &ScalVer{Major: 0, Date: "2025", Patch: 0},
			want: 1,
		},
		{
			name: "sv1 less than sv2 - patch",
			sv1:  &ScalVer{Major: 0, Date: "2025", Patch: 0},
			sv2:  &ScalVer{Major: 0, Date: "2025", Patch: 1},
			want: -1,
		},
		{
			name: "sv1 greater than sv2 - patch",
			sv1:  &ScalVer{Major: 0, Date: "2025", Patch: 1},
			sv2:  &ScalVer{Major: 0, Date: "2025", Patch: 0},
			want: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.sv1.Compare(tt.sv2); got != tt.want {
				t.Errorf("ScalVer.Compare() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsValidScalVerFormat(t *testing.T) {
	tests := []struct {
		name    string
		version string
		want    bool
	}{
		{
			name:    "valid yearly format",
			version: "0.2025.0",
			want:    true,
		},
		{
			name:    "valid monthly format",
			version: "0.202508.0",
			want:    true,
		},
		{
			name:    "valid daily format",
			version: "0.20250812.0",
			want:    true,
		},
		{
			name:    "valid development format",
			version: "0.20250812.0-dev-abc123",
			want:    true,
		},
		{
			name:    "invalid format",
			version: "0.2025",
			want:    false,
		},
		{
			name:    "empty string",
			version: "",
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidScalVerFormat(tt.version); got != tt.want {
				t.Errorf("IsValidScalVerFormat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenerateScalVer(t *testing.T) {
	tests := []struct {
		name           string
		major          int
		datePrecision  string
		patch          int
		wantErr        bool
		expectedFormat string
	}{
		{
			name:           "valid yearly generation",
			major:          0,
			datePrecision:  "yearly",
			patch:          0,
			wantErr:        false,
			expectedFormat: "0.2025.0", // Will be current year
		},
		{
			name:           "valid monthly generation",
			major:          0,
			datePrecision:  "monthly",
			patch:          0,
			wantErr:        false,
			expectedFormat: "0.202508.0", // Will be current year/month
		},
		{
			name:           "valid daily generation",
			major:          0,
			datePrecision:  "daily",
			patch:          0,
			wantErr:        false,
			expectedFormat: "0.20250812.0", // Will be current date
		},
		{
			name:          "invalid date precision",
			major:         0,
			datePrecision: "invalid",
			patch:         0,
			wantErr:       true,
		},
		{
			name:          "negative major",
			major:         -1,
			datePrecision: "yearly",
			patch:         0,
			wantErr:       true,
		},
		{
			name:          "negative patch",
			major:         0,
			datePrecision: "yearly",
			patch:         -1,
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateScalVer(tt.major, tt.datePrecision, tt.patch)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateScalVer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if !IsValidScalVerFormat(got) {
					t.Errorf("GenerateScalVer() generated invalid format: %s", got)
				}
				// Parse the generated version to verify components
				parsed, err := ParseScalVer(got)
				if err != nil {
					t.Errorf("Generated version could not be parsed: %v", err)
				}
				if parsed.Major != tt.major {
					t.Errorf("Generated major version = %v, want %v", parsed.Major, tt.major)
				}
				if parsed.Patch != tt.patch {
					t.Errorf("Generated patch version = %v, want %v", parsed.Patch, tt.patch)
				}
			}
		})
	}
}

func TestGenerateDevelopmentScalVer(t *testing.T) {
	tests := []struct {
		name      string
		gitCommit string
		wantErr   bool
	}{
		{
			name:      "valid git commit",
			gitCommit: "abc123",
			wantErr:   false,
		},
		{
			name:      "long git commit",
			gitCommit: "abcdef123456789",
			wantErr:   false,
		},
		{
			name:      "empty git commit",
			gitCommit: "",
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateDevelopmentScalVer(tt.gitCommit)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateDevelopmentScalVer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if !strings.Contains(got, "-dev-") {
					t.Errorf("Generated version does not contain development suffix: %s", got)
				}
				if !IsValidScalVerFormat(got) {
					t.Errorf("Generated version is not valid ScalVer format: %s", got)
				}
			}
		})
	}
}

func TestGetScalVerInfo(t *testing.T) {
	tests := []struct {
		name    string
		version string
		wantErr bool
	}{
		{
			name:    "valid version",
			version: "0.2025.0",
			wantErr: false,
		},
		{
			name:    "invalid version",
			version: "invalid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetScalVerInfo(tt.version)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetScalVerInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got["version"] != tt.version {
					t.Errorf("GetScalVerInfo() version = %v, want %v", got["version"], tt.version)
				}
				if got["format"] != "scalver" {
					t.Errorf("GetScalVerInfo() format = %v, want %v", got["format"], "scalver")
				}
			}
		})
	}
}

func TestValidateScalVerCompatibility(t *testing.T) {
	tests := []struct {
		name     string
		version1 string
		version2 string
		want     bool
		wantErr  bool
	}{
		{
			name:     "compatible development versions",
			version1: "0.2025.0",
			version2: "0.2025.1",
			want:     true,
			wantErr:  false,
		},
		{
			name:     "compatible stable versions",
			version1: "1.2025.0",
			version2: "1.2025.1",
			want:     true,
			wantErr:  false,
		},
		{
			name:     "incompatible major versions",
			version1: "0.2025.0",
			version2: "1.2025.0",
			want:     false,
			wantErr:  false,
		},
		{
			name:     "invalid first version",
			version1: "invalid",
			version2: "0.2025.0",
			want:     false,
			wantErr:  true,
		},
		{
			name:     "invalid second version",
			version1: "0.2025.0",
			version2: "invalid",
			want:     false,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ValidateScalVerCompatibility(tt.version1, tt.version2)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateScalVerCompatibility() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ValidateScalVerCompatibility() = %v, want %v", got, tt.want)
			}
		})
	}
}
