package schemas

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidator_IntegrationWithHCLFiles(t *testing.T) {
	validator := NewValidator()

	// Test with valid project HCL file
	t.Run("Valid Project HCL", func(t *testing.T) {
		hclContent := `
metadata {
  version = 1
  description = "A test project configuration"
}

project "test-project" {
  description = "A test automation project"
  version = "1.0.0"
  author = "Test User"
  email = "test@example.com"
  
  run {
    default_timeout = 300
    max_parallel = 10
  }
  
  facts {
    enabled = true
    timeout = 60
  }
  
  logging {
    level = "info"
    format = "text"
    output = "stdout"
  }
}
`
		result, err := validator.ValidateHCLContent("project", hclContent)
		if err != nil {
			t.Fatalf("Failed to validate HCL content: %v", err)
		}

		if !result.IsValid {
			t.Errorf("Expected valid project HCL, got errors: %v", result.Errors)
		}
	})

	// Test with valid machines HCL file
	t.Run("Valid Machines HCL", func(t *testing.T) {
		hclContent := `
metadata {
  version = 1
  description = "A test machines configuration"
}

machines {
  machine {
    name = "web-server-01"
    description = "Primary web server"
    hostname = "web01.example.com"
    port = 22
    user = "admin"
    
    authentication {
      method = "publickey"
      public_key_path = "~/.ssh/id_rsa"
    }
    
    tags = ["web", "production"]
  }
}
`
		result, err := validator.ValidateHCLContent("machines", hclContent)
		if err != nil {
			t.Fatalf("Failed to validate HCL content: %v", err)
		}
		if err != nil {
			t.Fatalf("Failed to validate HCL content: %v", err)
		}
		if !result.IsValid {
			t.Errorf("Expected valid machines HCL, got errors: %v", result.Errors)
		}
	})

	// Test with valid actions HCL file
	t.Run("Valid Actions HCL", func(t *testing.T) {
		hclContent := `
metadata {
  version = 1
  description = "A test actions configuration"
}

actions {
  action {
    name = "deploy-web"
    description = "Deploy web application"
    type = "command"
    command = "echo 'Deploying web app'"
    machines = ["web-server-01"]
  }
}
`
		result, err := validator.ValidateHCLContent("actions", hclContent)
		if err != nil {
			t.Fatalf("Failed to validate HCL content: %v", err)
		}
		if err != nil {
			t.Fatalf("Failed to validate HCL content: %v", err)
		}
		if !result.IsValid {
			t.Errorf("Expected valid actions HCL, got errors: %v", result.Errors)
		}
	})

	// Test with valid variables HCL file
	t.Run("Valid Variables HCL", func(t *testing.T) {
		hclContent := `
metadata {
  version = 1
  description = "A test variables configuration"
}

variables {
  variable {
    name = "database_url"
    description = "Database connection URL"
    value = "postgresql://user:pass@localhost/db"
  }
}
`
		result, err := validator.ValidateHCLContent("variables", hclContent)
		if err != nil {
			t.Fatalf("Failed to validate HCL content: %v", err)
		}
		if err != nil {
			t.Fatalf("Failed to validate HCL content: %v", err)
		}
		if !result.IsValid {
			t.Errorf("Expected valid variables HCL, got errors: %v", result.Errors)
		}
	})

	// Test with invalid project HCL file (missing required fields)
	t.Run("Invalid Project HCL - Missing Required Fields", func(t *testing.T) {
		hclContent := `
metadata {
  version = 1
}

project "test-project" {
  # Missing required description
}
`
		result, err := validator.ValidateHCLContent("project", hclContent)
		if err != nil {
			t.Fatalf("Failed to validate HCL content: %v", err)
		}
		if result.IsValid {
			t.Error("Expected invalid project HCL due to missing required fields")
		}
		if len(result.Errors) == 0 {
			t.Error("Expected validation errors for missing required fields")
		}
	})

	// Test with invalid machines HCL file (missing authentication)
	t.Run("Invalid Machines HCL - Missing Authentication", func(t *testing.T) {
		hclContent := `
metadata {
  version = 1
}

machines {
  machine {
    name = "web-server-01"
    hostname = "web01.example.com"
    user = "admin"
    # Missing authentication block
  }
}
`
		result, err := validator.ValidateHCLContent("machines", hclContent)
		if err != nil {
			t.Fatalf("Failed to validate HCL content: %v", err)
		}

		if result.IsValid {
			t.Error("Expected invalid machines HCL due to missing authentication")
		}
		if len(result.Errors) == 0 {
			t.Error("Expected validation errors for missing authentication")
		}
	})

	// Test with invalid actions HCL file (missing machine targeting)
	t.Run("Invalid Actions HCL - Missing Machine Targeting", func(t *testing.T) {
		hclContent := `
metadata {
  version = 1
}

actions {
  action {
    name = "deploy-web"
    description = "Deploy web application"
    type = "command"
    command = "echo 'Deploying web app'"
    # Missing machines or tags
  }
}
`
		result, err := validator.ValidateHCLContent("actions", hclContent)
		if err != nil {
			t.Fatalf("Failed to validate HCL content: %v", err)
		}
		if result.IsValid {
			t.Error("Expected invalid actions HCL due to missing machine targeting")
		}
		if len(result.Errors) == 0 {
			t.Error("Expected validation errors for missing machine targeting")
		}
	})

	// Test with invalid variables HCL file (invalid name)
	t.Run("Invalid Variables HCL - Invalid Name", func(t *testing.T) {
		hclContent := `
metadata {
  version = 1
}

variables {
  variable {
    name = "invalid@name"
    value = "test-value"
  }
}
`
		result, err := validator.ValidateHCLContent("variables", hclContent)
		if err != nil {
			t.Fatalf("Failed to validate HCL content: %v", err)
		}

		if result.IsValid {
			t.Error("Expected invalid variables HCL due to invalid name")
		}
		if len(result.Errors) == 0 {
			t.Error("Expected validation errors for invalid name")
		}
	})
}

func TestValidator_IntegrationWithRealFiles(t *testing.T) {
	validator := NewValidator()

	// Test with actual testdata files
	testDataDir := "testdata"
	if _, err := os.Stat(testDataDir); os.IsNotExist(err) {
		t.Skip("testdata directory not found, skipping real file tests")
	}

	// Test valid project file
	t.Run("Real Valid Project File", func(t *testing.T) {
		content, err := os.ReadFile(filepath.Join(testDataDir, "valid-project.hcl"))
		if err != nil {
			t.Skipf("Could not read valid-project.hcl: %v", err)
		}

		result, err := validator.ValidateHCLContent("project", string(content))
		if err != nil {
			t.Fatalf("Failed to validate HCL content: %v", err)
		}
		if err != nil {
			t.Fatalf("Failed to validate HCL content: %v", err)
		}
		if !result.IsValid {
			t.Errorf("Expected valid project file, got errors: %v", result.Errors)
		}
	})

	// Test valid machines file
	t.Run("Real Valid Machines File", func(t *testing.T) {
		content, err := os.ReadFile(filepath.Join(testDataDir, "valid-machines.hcl"))
		if err != nil {
			t.Skipf("Could not read valid-machines.hcl: %v", err)
		}

		result, err := validator.ValidateHCLContent("machines", string(content))
		if err != nil {
			t.Fatalf("Failed to validate HCL content: %v", err)
		}
		if err != nil {
			t.Fatalf("Failed to validate HCL content: %v", err)
		}
		if !result.IsValid {
			t.Errorf("Expected valid machines file, got errors: %v", result.Errors)
		}
	})

	// Test valid actions file
	t.Run("Real Valid Actions File", func(t *testing.T) {
		content, err := os.ReadFile(filepath.Join(testDataDir, "valid-actions.hcl"))
		if err != nil {
			t.Skipf("Could not read valid-actions.hcl: %v", err)
		}

		result, err := validator.ValidateHCLContent("actions", string(content))
		if err != nil {
			t.Fatalf("Failed to validate HCL content: %v", err)
		}
		if err != nil {
			t.Fatalf("Failed to validate HCL content: %v", err)
		}
		if !result.IsValid {
			t.Errorf("Expected valid actions file, got errors: %v", result.Errors)
		}
	})

	// Test invalid project file
	t.Run("Real Invalid Project File", func(t *testing.T) {
		content, err := os.ReadFile(filepath.Join(testDataDir, "invalid-project.hcl"))
		if err != nil {
			t.Skipf("Could not read invalid-project.hcl: %v", err)
		}

		result, err := validator.ValidateHCLContent("project", string(content))
		if err != nil {
			t.Fatalf("Failed to validate HCL content: %v", err)
		}
		if result.IsValid {
			t.Error("Expected invalid project file")
		}
		if len(result.Errors) == 0 {
			t.Error("Expected validation errors for invalid project file")
		}
	})

	// Test invalid variables file
	t.Run("Real Invalid Variables File", func(t *testing.T) {
		content, err := os.ReadFile(filepath.Join(testDataDir, "invalid-variables.hcl"))
		if err != nil {
			t.Skipf("Could not read invalid-variables.hcl: %v", err)
		}

		result, err := validator.ValidateHCLContent("variables", string(content))
		if err != nil {
			t.Fatalf("Failed to validate HCL content: %v", err)
		}
		if result.IsValid {
			t.Error("Expected invalid variables file")
		}
		if len(result.Errors) == 0 {
			t.Error("Expected validation errors for invalid variables file")
		}
	})
}

func TestValidator_EnhancedValidationIntegration(t *testing.T) {
	validator := NewValidator()

	// Test enhanced validation with security issues
	t.Run("Enhanced Validation - Security Issues", func(t *testing.T) {
		hclContent := `
metadata {
  version = 1
}

actions {
  action {
    name = "dangerous-action"
    description = "A dangerous action with shell operators"
    type = "command"
    command = "rm -rf / && echo 'dangerous'"
    machines = ["web-server-01"]
  }
}
`
		result, err := validator.ValidateHCLContent("actions", hclContent)
		if err != nil {
			t.Fatalf("Failed to validate HCL content: %v", err)
		}
		if result.IsValid {
			t.Error("Expected enhanced validation to catch security issue")
		}

		// Check for security validation error
		found := false
		for _, err := range result.Errors {
			if err.Message == "Shell operators (;&|$) are not allowed in commands for security reasons" {
				found = true
				break
			}
		}
		if !found {
			t.Error("Expected to find security validation error message")
		}
	})

	// Test enhanced validation with authentication issues
	t.Run("Enhanced Validation - Authentication Issues", func(t *testing.T) {
		hclContent := `
metadata {
  version = 1
}

machines {
  machine {
    name = "web-server-01"
    hostname = "web01.example.com"
    user = "admin"
    
    authentication {
      method = "publickey"
      # Missing public_key_path
    }
  }
}
`
		result, err := validator.ValidateHCLContent("machines", hclContent)
		if err != nil {
			t.Fatalf("Failed to validate HCL content: %v", err)
		}

		if result.IsValid {
			t.Error("Expected enhanced validation to catch authentication issue")
		}

		// Check for authentication validation error
		found := false
		for _, err := range result.Errors {
			if err.Message == "public_key_path is required when method is 'publickey'" {
				found = true
				break
			}
		}
		if !found {
			t.Error("Expected to find authentication validation error message")
		}
	})

	// Test enhanced validation with variable encryption issues
	t.Run("Enhanced Validation - Variable Encryption Issues", func(t *testing.T) {
		hclContent := `
metadata {
  version = 1
}

variables {
  variable {
    name = "database_password"
    value = "plaintext-password"
    encrypted = true
  }
}
`
		result, err := validator.ValidateHCLContent("variables", hclContent)
		if err != nil {
			t.Fatalf("Failed to validate HCL content: %v", err)
		}
		if result.IsValid {
			t.Error("Expected enhanced validation to catch encryption issue")
		}

		// Check for encryption validation error
		found := false
		for _, err := range result.Errors {
			if err.Message == "Encrypted values must be in age format (age1...)" {
				found = true
				break
			}
		}
		if !found {
			t.Error("Expected to find encryption validation error message")
		}
	})
}
