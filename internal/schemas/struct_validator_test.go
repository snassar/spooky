package schemas

import (
	"testing"
)

func TestStructValidator_ValidateProject(t *testing.T) {
	validator := NewStructValidator()

	tests := []struct {
		name       string
		content    map[string]interface{}
		wantValid  bool
		wantErrors int
	}{
		{
			name: "valid project",
			content: map[string]interface{}{
				"project": map[string]interface{}{
					"name":        "test-project",
					"description": "A test project",
				},
			},
			wantValid:  true,
			wantErrors: 0,
		},
		{
			name: "missing project block",
			content: map[string]interface{}{
				"something_else": "value",
			},
			wantValid:  false,
			wantErrors: 1,
		},
		{
			name: "missing required name",
			content: map[string]interface{}{
				"project": map[string]interface{}{
					"description": "A test project",
				},
			},
			wantValid:  false,
			wantErrors: 1,
		},
		{
			name: "missing required description",
			content: map[string]interface{}{
				"project": map[string]interface{}{
					"name": "test-project",
				},
			},
			wantValid:  false,
			wantErrors: 1,
		},
		{
			name: "invalid name pattern",
			content: map[string]interface{}{
				"project": map[string]interface{}{
					"name":        "123-invalid", // starts with number
					"description": "A test project",
				},
			},
			wantValid:  false,
			wantErrors: 1,
		},
		{
			name: "name too long",
			content: map[string]interface{}{
				"project": map[string]interface{}{
					"name":        string(make([]byte, 129)), // 129 characters
					"description": "A test project",
				},
			},
			wantValid:  false,
			wantErrors: 2, // Both pattern and length validation fail
		},
		{
			name: "description too long",
			content: map[string]interface{}{
				"project": map[string]interface{}{
					"name":        "test-project",
					"description": string(make([]byte, 1025)), // 1025 characters
				},
			},
			wantValid:  false,
			wantErrors: 1,
		},
		{
			name: "invalid run_max_parallel",
			content: map[string]interface{}{
				"project": map[string]interface{}{
					"name":             "test-project",
					"description":      "A test project",
					"run_max_parallel": 0, // must be >= 1
				},
			},
			wantValid:  false,
			wantErrors: 1,
		},
		{
			name: "valid with optional fields",
			content: map[string]interface{}{
				"project": map[string]interface{}{
					"name":             "test-project",
					"description":      "A test project",
					"run_max_parallel": 5,
					"facts_timeout":    60,
				},
			},
			wantValid:  true,
			wantErrors: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validator.ValidateProject(tt.content)

			if result.IsValid != tt.wantValid {
				t.Errorf("ValidateProject() IsValid = %v, want %v", result.IsValid, tt.wantValid)
			}

			if len(result.Errors) != tt.wantErrors {
				t.Errorf("ValidateProject() Errors count = %d, want %d", len(result.Errors), tt.wantErrors)
				for i, err := range result.Errors {
					t.Logf("Error %d: %s - %s", i, err.Field, err.Message)
				}
			}
		})
	}
}

func TestStructValidator_ValidateMachines(t *testing.T) {
	validator := NewStructValidator()

	tests := []struct {
		name       string
		content    map[string]interface{}
		wantValid  bool
		wantErrors int
	}{
		{
			name: "valid machines with one machine",
			content: map[string]interface{}{
				"machines": map[string]interface{}{
					"machine": []interface{}{
						map[string]interface{}{
							"name":     "test-machine",
							"hostname": "192.168.1.100",
							"user":     "admin",
							"authentication": map[string]interface{}{
								"method":          "publickey",
								"public_key_path": "/path/to/key",
							},
						},
					},
				},
			},
			wantValid:  true,
			wantErrors: 0,
		},
		{
			name: "valid machines with one group",
			content: map[string]interface{}{
				"machines": map[string]interface{}{
					"group": []interface{}{
						map[string]interface{}{
							"name":     "test-group",
							"machines": []interface{}{"machine1", "machine2"},
						},
					},
				},
			},
			wantValid:  true,
			wantErrors: 0,
		},
		{
			name: "missing machines block",
			content: map[string]interface{}{
				"something_else": "value",
			},
			wantValid:  false,
			wantErrors: 1,
		},
		{
			name: "missing required machine fields",
			content: map[string]interface{}{
				"machines": map[string]interface{}{
					"machine": []interface{}{
						map[string]interface{}{
							"name": "test-machine",
							// missing hostname, user, authentication
						},
					},
				},
			},
			wantValid:  false,
			wantErrors: 3, // hostname, user, authentication
		},
		{
			name: "invalid authentication method",
			content: map[string]interface{}{
				"machines": map[string]interface{}{
					"machine": []interface{}{
						map[string]interface{}{
							"name":     "test-machine",
							"hostname": "192.168.1.100",
							"user":     "admin",
							"authentication": map[string]interface{}{
								"method": "invalid_method",
							},
						},
					},
				},
			},
			wantValid:  false,
			wantErrors: 1,
		},
		{
			name: "publickey method without public_key_path",
			content: map[string]interface{}{
				"machines": map[string]interface{}{
					"machine": []interface{}{
						map[string]interface{}{
							"name":     "test-machine",
							"hostname": "192.168.1.100",
							"user":     "admin",
							"authentication": map[string]interface{}{
								"method": "publickey",
								// missing public_key_path
							},
						},
					},
				},
			},
			wantValid:  false,
			wantErrors: 1,
		},
		{
			name: "group with empty machines list",
			content: map[string]interface{}{
				"machines": map[string]interface{}{
					"group": []interface{}{
						map[string]interface{}{
							"name":     "test-group",
							"machines": []interface{}{}, // empty list
						},
					},
				},
			},
			wantValid:  false,
			wantErrors: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validator.ValidateMachines(tt.content)

			if result.IsValid != tt.wantValid {
				t.Errorf("ValidateMachines() IsValid = %v, want %v", result.IsValid, tt.wantValid)
			}

			if len(result.Errors) != tt.wantErrors {
				t.Errorf("ValidateMachines() Errors count = %d, want %d", len(result.Errors), tt.wantErrors)
				for i, err := range result.Errors {
					t.Logf("Error %d: %s - %s", i, err.Field, err.Message)
				}
			}
		})
	}
}

func TestStructValidator_ValidateActions(t *testing.T) {
	validator := NewStructValidator()

	tests := []struct {
		name       string
		content    map[string]interface{}
		wantValid  bool
		wantErrors int
	}{
		{
			name: "valid actions with one action",
			content: map[string]interface{}{
				"actions": map[string]interface{}{
					"action": []interface{}{
						map[string]interface{}{
							"name":        "test-action",
							"description": "A test action",
							"type":        "command",
						},
					},
				},
			},
			wantValid:  true,
			wantErrors: 0,
		},
		{
			name: "missing actions block",
			content: map[string]interface{}{
				"something_else": "value",
			},
			wantValid:  false,
			wantErrors: 1,
		},
		{
			name: "missing required action fields",
			content: map[string]interface{}{
				"actions": map[string]interface{}{
					"action": []interface{}{
						map[string]interface{}{
							"name": "test-action",
							// missing description and type
						},
					},
				},
			},
			wantValid:  false,
			wantErrors: 2, // description and type
		},
		{
			name: "invalid action type",
			content: map[string]interface{}{
				"actions": map[string]interface{}{
					"action": []interface{}{
						map[string]interface{}{
							"name":        "test-action",
							"description": "A test action",
							"type":        "invalid_type",
						},
					},
				},
			},
			wantValid:  false,
			wantErrors: 1,
		},
		{
			name: "description too short",
			content: map[string]interface{}{
				"actions": map[string]interface{}{
					"action": []interface{}{
						map[string]interface{}{
							"name":        "test-action",
							"description": "", // empty description
							"type":        "command",
						},
					},
				},
			},
			wantValid:  false,
			wantErrors: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validator.ValidateActions(tt.content)

			if result.IsValid != tt.wantValid {
				t.Errorf("ValidateActions() IsValid = %v, want %v", result.IsValid, tt.wantValid)
			}

			if len(result.Errors) != tt.wantErrors {
				t.Errorf("ValidateActions() Errors count = %d, want %d", len(result.Errors), tt.wantErrors)
				for i, err := range result.Errors {
					t.Logf("Error %d: %s - %s", i, err.Field, err.Message)
				}
			}
		})
	}
}

func TestStructValidator_ValidateVariables(t *testing.T) {
	validator := NewStructValidator()

	tests := []struct {
		name       string
		content    map[string]interface{}
		wantValid  bool
		wantErrors int
	}{
		{
			name: "valid variables with one variable",
			content: map[string]interface{}{
				"variables": map[string]interface{}{
					"variable": []interface{}{
						map[string]interface{}{
							"name":  "test_var",
							"value": "test_value",
						},
					},
				},
			},
			wantValid:  true,
			wantErrors: 0,
		},
		{
			name: "missing variables block",
			content: map[string]interface{}{
				"something_else": "value",
			},
			wantValid:  false,
			wantErrors: 1,
		},
		{
			name: "missing required variable fields",
			content: map[string]interface{}{
				"variables": map[string]interface{}{
					"variable": []interface{}{
						map[string]interface{}{
							"name": "test_var",
							// missing value
						},
					},
				},
			},
			wantValid:  false,
			wantErrors: 1,
		},
		{
			name: "variable with null value",
			content: map[string]interface{}{
				"variables": map[string]interface{}{
					"variable": []interface{}{
						map[string]interface{}{
							"name":  "test_var",
							"value": nil,
						},
					},
				},
			},
			wantValid:  false,
			wantErrors: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validator.ValidateVariables(tt.content)

			if result.IsValid != tt.wantValid {
				t.Errorf("ValidateVariables() IsValid = %v, want %v", result.IsValid, tt.wantValid)
			}

			if len(result.Errors) != tt.wantErrors {
				t.Errorf("ValidateVariables() Errors count = %d, want %d", len(result.Errors), tt.wantErrors)
				for i, err := range result.Errors {
					t.Logf("Error %d: %s - %s", i, err.Field, err.Message)
				}
			}
		})
	}
}
