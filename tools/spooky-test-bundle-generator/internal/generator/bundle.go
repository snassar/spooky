package generator

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"spooky-test-bundle-generator/internal/profiles"
)

// BundleGenerator handles the complete bundle generation process
type BundleGenerator struct {
	parser *profiles.Parser
}

// NewBundleGenerator creates a new bundle generator
func NewBundleGenerator() *BundleGenerator {
	return &BundleGenerator{
		parser: profiles.NewParser(),
	}
}

// GenerateBundle generates a complete test bundle from a profile
func (bg *BundleGenerator) GenerateBundle(profilePath, outputPath string) error {
	fmt.Printf("🔧 Generating test bundle from profile: %s\n", profilePath)
	fmt.Printf("📁 Output directory: %s\n", outputPath)

	// Parse and validate profile
	profile, err := bg.parser.ParseFile(profilePath)
	if err != nil {
		return fmt.Errorf("failed to parse profile: %w", err)
	}

	if err := bg.parser.ValidateProfile(profile); err != nil {
		return fmt.Errorf("profile validation failed: %w", err)
	}

	fmt.Printf("✅ Profile '%s' parsed and validated\n", profile.Name)

	// Create output directory
	if err := os.MkdirAll(outputPath, 0o755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Generate Containerfile
	if err := bg.generateContainerfile(profile, outputPath); err != nil {
		return fmt.Errorf("failed to generate Containerfile: %w", err)
	}

	// Generate SSH configuration
	if err := bg.generateSSHConfig(profile, outputPath); err != nil {
		return fmt.Errorf("failed to generate SSH configuration: %w", err)
	}

	// Generate Spooky project
	if err := bg.generateProject(profile, outputPath); err != nil {
		return fmt.Errorf("failed to generate Spooky project: %w", err)
	}

	// Generate justfile
	if err := bg.generateJustfile(profile, outputPath); err != nil {
		return fmt.Errorf("failed to generate justfile: %w", err)
	}

	// Generate README
	if err := bg.generateReadme(profile, outputPath); err != nil {
		return fmt.Errorf("failed to generate README: %w", err)
	}

	// Validate generated project
	if err := bg.validateGeneratedProject(outputPath); err != nil {
		return fmt.Errorf("generated project validation failed: %w", err)
	}

	fmt.Println("✅ Bundle generation completed successfully!")
	fmt.Printf("💡 To test: cd %s && just test\n", outputPath)

	return nil
}

// generateContainerfile creates the Containerfile
func (bg *BundleGenerator) generateContainerfile(profile *profiles.Profile, outputPath string) error {
	fmt.Println("📦 Generating Containerfile...")

	containerfilePath := filepath.Join(outputPath, "Containerfile")
	content := fmt.Sprintf(`FROM %s

# Install SSH server
RUN apt-get update && apt-get install -y openssh-server

# Copy SSH configuration
COPY sshd_config /etc/ssh/sshd_config

# Create test user
RUN useradd -m -s /bin/bash testuser
RUN mkdir -p /home/testuser/.ssh

# Copy test SSH keys (if available)
COPY test-keys/id_rsa.pub /home/testuser/.ssh/authorized_keys || true
RUN chown -R testuser:testuser /home/testuser/.ssh
RUN chmod 700 /home/testuser/.ssh
RUN chmod 600 /home/testuser/.ssh/authorized_keys || true

# Start SSH service
EXPOSE %d
CMD ["/usr/sbin/sshd", "-D"]
`, profile.Container.BaseImage, profile.SSH.Port)

	return os.WriteFile(containerfilePath, []byte(content), 0o644)
}

// generateSSHConfig creates the SSH configuration
func (bg *BundleGenerator) generateSSHConfig(profile *profiles.Profile, outputPath string) error {
	fmt.Println("🔐 Generating SSH configuration...")

	sshConfigPath := filepath.Join(outputPath, "sshd_config")
	content := fmt.Sprintf(`# SSH Configuration for %s
Port %d
Protocol 2
HostKey /etc/ssh/ssh_host_rsa_key
HostKey /etc/ssh/ssh_host_ecdsa_key
HostKey /etc/ssh/ssh_host_ed25519_key

# Authentication settings
PermitRootLogin %s
PasswordAuthentication %s
PubkeyAuthentication %s
AuthorizedKeysFile .ssh/authorized_keys

# Connection settings
ClientAliveInterval 60
ClientAliveCountMax 3
MaxSessions 10

# Logging
LogLevel INFO
SyslogFacility AUTH
`, profile.Container.Name, profile.SSH.Port,
		boolToYesNo(profile.SSH.Settings != nil && profile.SSH.Settings.PermitRootLogin),
		boolToYesNo(profile.SSH.Settings != nil && profile.SSH.Settings.PasswordAuthentication),
		boolToYesNo(profile.SSH.Settings != nil && profile.SSH.Settings.PubkeyAuthentication))

	return os.WriteFile(sshConfigPath, []byte(content), 0o644)
}

// generateProject creates the Spooky project files
func (bg *BundleGenerator) generateProject(profile *profiles.Profile, outputPath string) error {
	fmt.Println("📁 Generating Spooky project...")

	projectPath := filepath.Join(outputPath, "spooky-project")
	if err := os.MkdirAll(projectPath, 0o755); err != nil {
		return fmt.Errorf("failed to create project directory: %w", err)
	}

	// Create project.hcl
	projectHCL := fmt.Sprintf(`project "%s" {
  description = "%s"
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
`, profile.SpookyProject.Name, profile.SpookyProject.Description)

	if err := os.WriteFile(filepath.Join(projectPath, "project.hcl"), []byte(projectHCL), 0o644); err != nil {
		return fmt.Errorf("failed to write project.hcl: %w", err)
	}

	// Create machines.hcl
	machinesHCL := fmt.Sprintf(`machines {
  machine "%s" {
    hostname = "%s"
    port = %d
    user = "%s"
    authentication {
      method = "%s"
      %s
    }
    tags = %s
  }
}
`,
		profile.SpookyProject.Machines.Machines[0].Name,
		profile.SpookyProject.Machines.Machines[0].Hostname,
		profile.SpookyProject.Machines.Machines[0].Port,
		profile.SpookyProject.Machines.Machines[0].User,
		profile.SpookyProject.Machines.Machines[0].Authentication.Method,
		bg.generateAuthConfig(profile.SpookyProject.Machines.Machines[0].Authentication),
		bg.generateTags(profile.SpookyProject.Machines.Machines[0].Tags))

	if err := os.WriteFile(filepath.Join(projectPath, "machines.hcl"), []byte(machinesHCL), 0o644); err != nil {
		return fmt.Errorf("failed to write machines.hcl: %w", err)
	}

	// Create variables.hcl
	var variablesHCL strings.Builder
	variablesHCL.WriteString("variables {\n")
	for _, variable := range profile.SpookyProject.Variables.Variables {
		variablesHCL.WriteString(fmt.Sprintf(`  variable "%s" {
    value = "%s"
    description = "%s"
  }
`, variable.Name, variable.Value, variable.Description))
	}
	variablesHCL.WriteString("}\n")

	if err := os.WriteFile(filepath.Join(projectPath, "variables.hcl"), []byte(variablesHCL.String()), 0o644); err != nil {
		return fmt.Errorf("failed to write variables.hcl: %w", err)
	}

	// Create actions.hcl
	var actionsHCL strings.Builder
	actionsHCL.WriteString("actions {\n")
	for _, action := range profile.SpookyProject.Actions.Actions {
		actionsHCL.WriteString(fmt.Sprintf(`  action "%s" {
    description = "%s"
    type = "%s"
    %s
    tags = %s
  }
`, action.Name, action.Description, action.Type, bg.generateActionConfig(action), bg.generateTags(action.Tags)))
	}
	actionsHCL.WriteString("}\n")

	if err := os.WriteFile(filepath.Join(projectPath, "actions.hcl"), []byte(actionsHCL.String()), 0o644); err != nil {
		return fmt.Errorf("failed to write actions.hcl: %w", err)
	}

	// Create directories
	dirs := []string{"templates", "files", "logs"}
	for _, dir := range dirs {
		if err := os.MkdirAll(filepath.Join(projectPath, dir), 0o755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// Generate template files from profile
	if err := bg.generateTemplates(profile, projectPath); err != nil {
		return fmt.Errorf("failed to generate templates: %w", err)
	}

	// Generate file files from profile
	if err := bg.generateFiles(profile, projectPath); err != nil {
		return fmt.Errorf("failed to generate files: %w", err)
	}

	return nil
}

// generateAuthConfig generates the authentication configuration string
func (bg *BundleGenerator) generateAuthConfig(auth *profiles.AuthenticationConfig) string {
	if auth.Method == "password" && auth.Password != nil {
		return fmt.Sprintf(`password {
      value = "%s"
      encrypted = %t
    }`, auth.Password.Value, auth.Password.Encrypted)
	} else if auth.Method == "publickey" && auth.PublicKey != nil {
		return fmt.Sprintf(`public_key {
      public_key_path = "%s"
    }`, auth.PublicKey.PublicKeyPath)
	}
	return ""
}

// generateActionConfig generates the action configuration string
func (bg *BundleGenerator) generateActionConfig(action *profiles.ActionConfig) string {
	var config strings.Builder

	if action.Command != nil {
		config.WriteString(fmt.Sprintf(`    command = "%s"
    `, *action.Command))
	}
	if action.Source != nil {
		config.WriteString(fmt.Sprintf(`    source = "%s"
    `, *action.Source))
	}
	if action.Destination != nil {
		config.WriteString(fmt.Sprintf(`    destination = "%s"
    `, *action.Destination))
	}
	if action.Validate != nil {
		config.WriteString(fmt.Sprintf(`    validate = %t
    `, *action.Validate))
	}
	if action.SyncSource != nil {
		config.WriteString(fmt.Sprintf(`    sync_source = "%s"
    `, *action.SyncSource))
	}
	if action.SyncDestination != nil {
		config.WriteString(fmt.Sprintf(`    sync_destination = "%s"
    `, *action.SyncDestination))
	}
	if action.SyncMode != nil {
		config.WriteString(fmt.Sprintf(`    sync_mode = "%s"
    `, *action.SyncMode))
	}

	return config.String()
}

// generateTags generates the tags string
func (bg *BundleGenerator) generateTags(tags []string) string {
	if len(tags) == 0 {
		return "[]"
	}

	var tagStrings []string
	for _, tag := range tags {
		tagStrings = append(tagStrings, fmt.Sprintf(`"%s"`, tag))
	}
	return fmt.Sprintf("[%s]", strings.Join(tagStrings, ", "))
}

// generateTemplates generates template files from the profile
func (bg *BundleGenerator) generateTemplates(profile *profiles.Profile, projectPath string) error {
	if len(profile.Templates) == 0 {
		return nil
	}

	templatesDir := filepath.Join(projectPath, "templates")

	for _, templatesConfig := range profile.Templates {
		for _, template := range templatesConfig.Templates {
			templatePath := filepath.Join(templatesDir, template.Name)
			if err := os.WriteFile(templatePath, []byte(template.Content), 0o644); err != nil {
				return fmt.Errorf("failed to write template %s: %w", template.Name, err)
			}
		}
	}

	return nil
}

// generateFiles generates static files from the profile
func (bg *BundleGenerator) generateFiles(profile *profiles.Profile, projectPath string) error {
	if len(profile.Files) == 0 {
		return nil
	}

	filesDir := filepath.Join(projectPath, "files")

	for _, filesConfig := range profile.Files {
		for _, file := range filesConfig.Files {
			// Create directory structure if needed
			filePath := filepath.Join(filesDir, file.Name)
			fileDir := filepath.Dir(filePath)
			if err := os.MkdirAll(fileDir, 0o755); err != nil {
				return fmt.Errorf("failed to create file directory %s: %w", fileDir, err)
			}

			if err := os.WriteFile(filePath, []byte(file.Content), 0o644); err != nil {
				return fmt.Errorf("failed to write file %s: %w", file.Name, err)
			}

			// Set permissions if specified
			if file.Permissions != nil {
				// Parse permissions (e.g., "755" -> 0755)
				perm, err := strconv.ParseUint(*file.Permissions, 8, 32)
				if err == nil {
					if err := os.Chmod(filePath, os.FileMode(perm)); err != nil {
						return fmt.Errorf("failed to set permissions for %s: %w", file.Name, err)
					}
				}
			}
		}
	}

	return nil
}

// generateJustfile creates the justfile for testing
func (bg *BundleGenerator) generateJustfile(profile *profiles.Profile, outputPath string) error {
	fmt.Println("⚡ Generating justfile...")

	justfilePath := filepath.Join(outputPath, "justfile")
	content := fmt.Sprintf(`# Test Bundle Justfile for %s

default:
    @just --list

start:
    #!/usr/bin/env bash
    echo "Starting %s test container..."
    podman run -d --name %s --network spooky-test-net --ip %s -p %d:22 %s:latest

stop:
    #!/usr/bin/env bash
    echo "Stopping %s test container..."
    podman stop %s || true
    podman rm %s || true

build:
    #!/usr/bin/env bash
    echo "Building %s test container..."
    podman build -t %s:latest .

test-ssh:
    #!/usr/bin/env bash
    echo "Testing SSH connectivity..."
    ssh -o ConnectTimeout=5 -o StrictHostKeyChecking=no -p %d testuser@localhost "echo SSH connection successful"

test-validate:
    #!/usr/bin/env bash
    echo "Validating Spooky project..."
    cd spooky-project
    ../../spooky project validate

test: test-ssh test-validate
    #!/usr/bin/env bash
    echo "All tests completed"

clean: stop
    #!/usr/bin/env bash
    echo "Cleaning up test environment..."
    podman rmi %s:latest || true
`,
		profile.Container.Name,
		profile.Container.Name, profile.Container.Name, profile.Container.IP, profile.SSH.Port, profile.Container.Name,
		profile.Container.Name, profile.Container.Name, profile.Container.Name,
		profile.Container.Name, profile.Container.Name,
		profile.SSH.Port,
		profile.Container.Name)

	return os.WriteFile(justfilePath, []byte(content), 0o644)
}

// generateReadme creates the README file
func (bg *BundleGenerator) generateReadme(profile *profiles.Profile, outputPath string) error {
	fmt.Println("📖 Generating README...")

	readmePath := filepath.Join(outputPath, "README.md")
	content := fmt.Sprintf(`# %s Test Bundle

%s

## Quick Start

1. Build the container:
   just build

2. Start the container:
   just start

3. Test SSH connectivity:
   just test-ssh

4. Run Spooky tests:
   just test

5. Clean up:
   just clean

## Container Details

- Image: %s
- Name: %s
- IP: %s
- SSH Port: %d
- SSH User: testuser

## Spooky Project

The generated Spooky project is located in the spooky-project/ directory and includes:

- project.hcl - Project configuration
- machines.hcl - Machine inventory
- variables.hcl - Project variables
- actions.hcl - Test actions
- templates/ - Template files
- files/ - Static files
- logs/ - Log directory

## Testing

This bundle is designed to test:

- SSH connectivity and authentication
- Facts gathering from remote machines
- Template rendering and deployment
- File synchronization
- Action execution

## Network Setup

The container uses the spooky-test-net network. If it doesn't exist, create it:

podman network create --subnet 10.0.100.0/24 spooky-test-net
`,
		profile.Container.Name,
		profile.Description,
		profile.Container.BaseImage,
		profile.Container.Name,
		profile.Container.IP,
		profile.SSH.Port)

	return os.WriteFile(readmePath, []byte(content), 0o644)
}

// validateGeneratedProject validates the generated Spooky project
func (bg *BundleGenerator) validateGeneratedProject(projectDir string) error {
	// Check if spooky binary exists
	spookyPath := "./spooky"
	if _, err := os.Stat(spookyPath); os.IsNotExist(err) {
		return fmt.Errorf("spooky binary not found. Please build it first: go build -o spooky ../../main.go")
	}

	fmt.Println("🔍 Validating generated Spooky project...")

	// Change to project directory for validation
	originalDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}
	defer os.Chdir(originalDir)

	projectPath := filepath.Join(projectDir, "spooky-project")
	if err := os.Chdir(projectPath); err != nil {
		return fmt.Errorf("failed to change to project directory: %w", err)
	}

	// Run spooky project validate
	cmd := exec.Command("../../spooky", "project", "validate")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("spooky project validation failed: %w", err)
	}

	fmt.Println("✅ Generated project validation passed")
	return nil
}

// boolToYesNo converts a boolean to "yes" or "no"
func boolToYesNo(b bool) string {
	if b {
		return "yes"
	}
	return "no"
}
