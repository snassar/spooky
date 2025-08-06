package facts

import (
	"strings"
	"testing"
)

func TestFormatDetector_DetectFormat(t *testing.T) {
	detector := NewFormatDetector()

	tests := []struct {
		name     string
		content  string
		expected string
		hasError bool
	}{
		{
			name: "Valid JSON facts",
			content: `{
				"metadata": {
					"exported_at": "2023-01-01T00:00:00Z",
					"project_id": "test-project",
					"export_format": "json",
					"version": "1.0"
				},
				"global_facts": [
					{
						"machine_id": "test-server",
						"collected_at": "2023-01-01T00:00:00Z",
						"ttl": "24h",
						"facts": {
							"os": "linux",
							"cpu_count": 4
						}
					}
				],
				"project_facts": []
			}`,
			expected: "json",
			hasError: false,
		},
		{
			name: "Valid HCL facts",
			content: `facts_export = {
  metadata = {
    exported_at = "2023-01-01T00:00:00Z"
    project_id = "test-project"
    export_format = "hcl"
    version = "1.0"
  }
  global_facts = [
    {
      machine_id = "test-server"
      collected_at = "2023-01-01T00:00:00Z"
      ttl = "24h"
      facts = {
        system = {
          os = {
            name = "linux"
            version = "unknown"
            arch = "x86_64"
            kernel = "unknown"
          }
        }
        enhanced = {}
      }
    }
  ]
  project_facts = []
}`,
			expected: "hcl",
			hasError: false,
		},
		{
			name:     "Invalid content",
			content:  "This is not JSON or HCL content",
			expected: "",
			hasError: true,
		},
		{
			name:     "Empty content",
			content:  "",
			expected: "",
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.content)
			format, err := detector.DetectFormat(reader)

			if tt.hasError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if format != tt.expected {
				t.Errorf("Expected format %s, got %s", tt.expected, format)
			}
		})
	}
}

func TestFormatDetector_IsValidJSON(t *testing.T) {
	detector := NewFormatDetector()

	tests := []struct {
		name     string
		content  []byte
		expected bool
	}{
		{
			name:     "Valid JSON object",
			content:  []byte(`{"key": "value"}`),
			expected: true,
		},
		{
			name:     "Valid JSON array",
			content:  []byte(`[1, 2, 3]`),
			expected: true,
		},
		{
			name:     "Invalid JSON",
			content:  []byte(`{"key": "value"`),
			expected: false,
		},
		{
			name:     "Plain text",
			content:  []byte(`This is not JSON`),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := detector.isValidJSON(tt.content)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestFormatDetector_IsValidHCL(t *testing.T) {
	detector := NewFormatDetector()

	tests := []struct {
		name     string
		content  []byte
		expected bool
	}{
		{
			name: "Valid HCL with global_facts",
			content: []byte(`global_facts = [
				{
					machine_id = "test"
					collected_at = "2023-01-01T00:00:00Z"
					ttl = "24h"
				}
			]`),
			expected: true,
		},
		{
			name: "Valid HCL with project_facts",
			content: []byte(`project_facts = [
				{
					machine_id = "test"
					project_id = "test-project"
					collected_at = "2023-01-01T00:00:00Z"
					ttl = "24h"
				}
			]`),
			expected: true,
		},
		{
			name:     "Invalid HCL",
			content:  []byte(`This is not HCL content`),
			expected: false,
		},
		{
			name:     "Empty content",
			content:  []byte(``),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := detector.isValidHCL(tt.content)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestFormatDetector_ValidateAgainstJSONSchema(t *testing.T) {
	detector := NewFormatDetector()

	// Test with valid JSON facts content
	validJSON := []byte(`{
		"metadata": {
			"exported_at": "2023-01-01T00:00:00Z",
			"project_id": "test-project",
			"export_format": "json",
			"version": "1.0"
		},
		"global_facts": [
			{
				"machine_id": "test-server",
				"collected_at": "2023-01-01T00:00:00Z",
				"ttl": "24h",
				"facts": {
					"os": "linux",
					"cpu_count": 4
				}
			}
		],
		"project_facts": []
	}`)

	result := detector.validateAgainstJSONSchema(validJSON)
	if !result {
		t.Error("Expected valid JSON to pass schema validation")
	}
}

func TestFormatDetector_ValidateAgainstHCLSchema(t *testing.T) {
	detector := NewFormatDetector()

	// Test with valid HCL facts content
	validHCL := []byte(`facts_export = {
  metadata = {
    exported_at = "2023-01-01T00:00:00Z"
    project_id = "test-project"
    export_format = "hcl"
    version = "1.0"
  }
  global_facts = [
    {
      machine_id = "test-server"
      collected_at = "2023-01-01T00:00:00Z"
      ttl = "24h"
      facts = {
        system = {
          os = {
            name = "linux"
            version = "unknown"
            arch = "x86_64"
            kernel = "unknown"
          }
        }
        enhanced = {}
      }
    }
  ]
  project_facts = []
}`)

	result := detector.validateAgainstHCLSchema(validHCL)
	if !result {
		t.Error("Expected valid HCL to pass schema validation")
	}
}
