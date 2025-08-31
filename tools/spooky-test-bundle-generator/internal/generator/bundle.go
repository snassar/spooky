package generator

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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

	// Generate container configurations
	if err := bg.generateContainers(profile, outputPath); err != nil {
		return fmt.Errorf("failed to generate container configurations: %w", err)
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

// generateContainers creates container configurations for all OSes
func (bg *BundleGenerator) generateContainers(profile *profiles.Profile, outputPath string) error {
	fmt.Println("📦 Generating container configurations...")

	containersDir := filepath.Join(outputPath, "containers")
	if err := os.MkdirAll(containersDir, 0o755); err != nil {
		return fmt.Errorf("failed to create containers directory: %w", err)
	}

	// Generate container configs for each OS
	osConfigs := []struct {
		name    string
		config  *profiles.ContainerConfig
		machine *profiles.MachineConfig
	}{
		{"debian13", profile.Containers.Debian13, profile.Project.Machines.Debian13},
		{"fedora42", profile.Containers.Fedora42, profile.Project.Machines.Fedora42},
		{"arch", profile.Containers.Arch, profile.Project.Machines.Arch},
		{"alpine319", profile.Containers.Alpine319, profile.Project.Machines.Alpine319},
		{"opensuse156", profile.Containers.Opensuse156, profile.Project.Machines.Opensuse156},
	}

	for _, osConfig := range osConfigs {
		if osConfig.config != nil && osConfig.machine != nil {
			if err := bg.generateContainerConfig(osConfig.name, osConfig.config, osConfig.machine, containersDir); err != nil {
				return fmt.Errorf("failed to generate %s container config: %w", osConfig.name, err)
			}
		}
	}

	return nil
}

// generateContainerConfig creates container configuration for a specific OS
func (bg *BundleGenerator) generateContainerConfig(osName string, container *profiles.ContainerConfig, machine *profiles.MachineConfig, containersDir string) error {
	osDir := filepath.Join(containersDir, osName)
	if err := os.MkdirAll(osDir, 0o755); err != nil {
		return fmt.Errorf("failed to create %s directory: %w", osName, err)
	}

	// Generate Containerfile
	if err := bg.generateContainerfile(osName, container, osDir); err != nil {
		return fmt.Errorf("failed to generate Containerfile for %s: %w", osName, err)
	}

	// Generate SSH config
	if err := bg.generateSSHConfig(osName, container, osDir); err != nil {
		return fmt.Errorf("failed to generate SSH config for %s: %w", osName, err)
	}

	return nil
}

// generateContainerfile creates a Containerfile for a specific OS
func (bg *BundleGenerator) generateContainerfile(osName string, container *profiles.ContainerConfig, osDir string) error {
	containerfilePath := filepath.Join(osDir, "Containerfile")

	var content strings.Builder
	content.WriteString(fmt.Sprintf("FROM %s\n\n", container.BaseImage))

	// Generate package installation commands based on OS
	switch osName {
	case "debian13":
		content.WriteString("# Update package list and install packages\n")
		content.WriteString("RUN apt-get update && apt-get install -y \\\n")
		for i, pkg := range container.Packages {
			if i == len(container.Packages)-1 {
				content.WriteString(fmt.Sprintf("    %s\n\n", pkg))
			} else {
				content.WriteString(fmt.Sprintf("    %s \\\n", pkg))
			}
		}
	case "fedora42":
		content.WriteString("# Install packages\n")
		content.WriteString("RUN dnf install -y \\\n")
		for i, pkg := range container.Packages {
			if i == len(container.Packages)-1 {
				content.WriteString(fmt.Sprintf("    %s\n\n", pkg))
			} else {
				content.WriteString(fmt.Sprintf("    %s \\\n", pkg))
			}
		}
	case "arch":
		content.WriteString("# Install packages\n")
		content.WriteString("RUN pacman -Syu --noconfirm \\\n")
		for i, pkg := range container.Packages {
			if i == len(container.Packages)-1 {
				content.WriteString(fmt.Sprintf("    %s\n\n", pkg))
			} else {
				content.WriteString(fmt.Sprintf("    %s \\\n", pkg))
			}
		}
	case "alpine319":
		content.WriteString("# Install packages\n")
		content.WriteString("RUN apk add --no-cache \\\n")
		for i, pkg := range container.Packages {
			if i == len(container.Packages)-1 {
				content.WriteString(fmt.Sprintf("    %s\n\n", pkg))
			} else {
				content.WriteString(fmt.Sprintf("    %s \\\n", pkg))
			}
		}
	case "opensuse156":
		content.WriteString("# Install packages\n")
		content.WriteString("RUN zypper install -y \\\n")
		for i, pkg := range container.Packages {
			if i == len(container.Packages)-1 {
				content.WriteString(fmt.Sprintf("    %s\n\n", pkg))
			} else {
				content.WriteString(fmt.Sprintf("    %s \\\n", pkg))
			}
		}
		// Add vim package for xxd command
		content.WriteString("# Install additional tools\n")
		content.WriteString("RUN zypper install -y vim\n\n")
	}

	// Generate machine ID
	content.WriteString("# Generate machine ID\n")
	if osName == "opensuse156" {
		content.WriteString("RUN echo \"$(head -c 16 /dev/urandom | hexdump -e '16/1 \"%02x\"')\" > /etc/machine-id\n\n")
	} else {
		content.WriteString("RUN echo \"$(head -c 16 /dev/urandom | xxd -p)\" > /etc/machine-id\n\n")
	}

	// Copy SSH config
	content.WriteString("# Copy SSH configuration\n")
	content.WriteString("COPY sshd_config /etc/ssh/sshd_config\n\n")

	// Generate SSH host keys
	content.WriteString("# Generate SSH host keys\n")
	content.WriteString("RUN ssh-keygen -A\n\n")

	// Setup SSH service
	content.WriteString("# Setup SSH service\n")
	if osName == "alpine319" {
		// Alpine containers don't need OpenRC, just create sshd directory
		content.WriteString("RUN mkdir -p /run/sshd\n")
	} else {
		// For systemd-based systems, we'll start sshd manually in containers
		content.WriteString("RUN mkdir -p /run/sshd\n")
	}
	content.WriteString("\n")

	// Expose SSH port
	content.WriteString(fmt.Sprintf("EXPOSE %d\n\n", container.SSHPort))

	// Start command
	content.WriteString("# Start SSH service\n")
	content.WriteString("CMD [\"/usr/sbin/sshd\", \"-D\"]\n")

	return os.WriteFile(containerfilePath, []byte(content.String()), 0o644)
}

// generateSSHConfig creates SSH configuration for a specific OS
func (bg *BundleGenerator) generateSSHConfig(osName string, container *profiles.ContainerConfig, osDir string) error {
	sshConfigPath := filepath.Join(osDir, "sshd_config")

	var content strings.Builder
	content.WriteString("# SSH Configuration for Spooky Testing\n")
	content.WriteString(fmt.Sprintf("Port %d\n", container.SSHPort))
	content.WriteString("Protocol 2\n")
	content.WriteString("HostKey /etc/ssh/ssh_host_rsa_key\n")
	content.WriteString("HostKey /etc/ssh/ssh_host_ecdsa_key\n")
	content.WriteString("HostKey /etc/ssh/ssh_host_ed25519_key\n")

	if container.SSHConfig.PasswordAuth {
		content.WriteString("PasswordAuthentication yes\n")
	} else {
		content.WriteString("PasswordAuthentication no\n")
	}

	if container.SSHConfig.PubkeyAuth {
		content.WriteString("PubkeyAuthentication yes\n")
	} else {
		content.WriteString("PubkeyAuthentication no\n")
	}

	content.WriteString(fmt.Sprintf("PermitRootLogin %s\n", container.SSHConfig.PermitRootLogin))

	if container.SSHConfig.StrictModes {
		content.WriteString("StrictModes yes\n")
	} else {
		content.WriteString("StrictModes no\n")
	}

	if container.SSHConfig.MaxAuthTries != nil {
		content.WriteString(fmt.Sprintf("MaxAuthTries %d\n", *container.SSHConfig.MaxAuthTries))
	}

	if container.SSHConfig.LoginGraceTime != nil {
		content.WriteString(fmt.Sprintf("LoginGraceTime %d\n", *container.SSHConfig.LoginGraceTime))
	}

	content.WriteString("X11Forwarding no\n")
	content.WriteString("PrintMotd no\n")
	content.WriteString("AcceptEnv LANG LC_*\n")

	// Use different sftp subsystem paths for different OSes
	if osName == "alpine319" {
		content.WriteString("Subsystem sftp /usr/lib/openssh/sftp-server\n")
	} else {
		content.WriteString("Subsystem sftp /usr/libexec/openssh/sftp-server\n")
	}

	return os.WriteFile(sshConfigPath, []byte(content.String()), 0o644)
}

// generateProject creates the Spooky project structure
func (bg *BundleGenerator) generateProject(profile *profiles.Profile, outputPath string) error {
	fmt.Println("📁 Generating Spooky project...")

	projectDir := filepath.Join(outputPath, "spooky-project")
	if err := os.MkdirAll(projectDir, 0o755); err != nil {
		return fmt.Errorf("failed to create project directory: %w", err)
	}

	// Get the current working directory (where the generator is running from)
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	// Look for spooky binary in the generator directory
	spookyPath := filepath.Join(currentDir, "spooky")
	if _, err := os.Stat(spookyPath); os.IsNotExist(err) {
		return fmt.Errorf("spooky binary not found at %s. Please build it first: just build-spooky-test-bundle-generator", spookyPath)
	}

	// Use spooky project init to generate all base HCL files
	cmd := exec.Command(spookyPath, "project", "init",
		"--name", profile.Project.Name,
		"--description", profile.Project.Description,
		projectDir)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to run spooky project init: %w\nOutput: %s", string(output), err)
	}

	// Now customize the generated files for test bundle specific needs
	if err := bg.customizeProjectHCL(profile, projectDir); err != nil {
		return fmt.Errorf("failed to customize project.hcl: %w", err)
	}

	if err := bg.customizeMachinesHCL(profile, projectDir); err != nil {
		return fmt.Errorf("failed to customize machines.hcl: %w", err)
	}

	// Generate templates directory
	if err := bg.generateTemplates(profile, projectDir); err != nil {
		return fmt.Errorf("failed to generate templates: %w", err)
	}

	// Generate files directory
	if err := bg.generateFiles(profile, projectDir); err != nil {
		return fmt.Errorf("failed to generate files: %w", err)
	}

	// Generate Spooky runner container
	if err := bg.generateSpookyRunner(profile, outputPath); err != nil {
		return fmt.Errorf("failed to generate Spooky runner: %w", err)
	}

	return nil
}

// customizeProjectHCL customizes the project.hcl file for the test bundle
// We still customize this because we need to use the profile's project name and description
// rather than the generic template from spooky project init
func (bg *BundleGenerator) customizeProjectHCL(profile *profiles.Profile, projectDir string) error {
	projectHCLPath := filepath.Join(projectDir, "project.hcl")

	var content strings.Builder
	content.WriteString("# Spooky Project Configuration\n")
	content.WriteString(fmt.Sprintf("project \"%s\" {\n", profile.Project.Name))
	content.WriteString(fmt.Sprintf("  description = \"%s\"\n", profile.Project.Description))
	content.WriteString("}\n")

	return os.WriteFile(projectHCLPath, []byte(content.String()), 0o644)
}

// customizeMachinesHCL customizes the machines.hcl file for the test bundle
// We still customize this because we need to add test-specific machine configurations
// for the different OS containers in the test bundle
func (bg *BundleGenerator) customizeMachinesHCL(profile *profiles.Profile, projectDir string) error {
	machinesHCLPath := filepath.Join(projectDir, "machines.hcl")

	var content strings.Builder
	content.WriteString("machines {\n")

	// Generate machine configs for each OS
	osConfigs := []struct {
		name    string
		config  *profiles.ContainerConfig
		machine *profiles.MachineConfig
	}{
		{"debian13", profile.Containers.Debian13, profile.Project.Machines.Debian13},
		{"fedora42", profile.Containers.Fedora42, profile.Project.Machines.Fedora42},
		{"arch", profile.Containers.Arch, profile.Project.Machines.Arch},
		{"alpine319", profile.Containers.Alpine319, profile.Project.Machines.Alpine319},
		{"opensuse156", profile.Containers.Opensuse156, profile.Project.Machines.Opensuse156},
	}

	for _, osConfig := range osConfigs {
		if osConfig.config != nil && osConfig.machine != nil {
			content.WriteString(fmt.Sprintf("  machine \"%s\" {\n", osConfig.name))
			content.WriteString(fmt.Sprintf("    hostname = \"%s\"\n", osConfig.machine.IP)) // Use internal container IPs
			content.WriteString(fmt.Sprintf("    port = %d\n", osConfig.machine.Port))

			// User field (required)
			user := "root"
			if osConfig.machine.Auth.Username != nil {
				user = *osConfig.machine.Auth.Username
			}
			content.WriteString(fmt.Sprintf("    user = \"%s\"\n", user))

			// Authentication block (required)
			content.WriteString("    authentication \"password\" {\n")
			content.WriteString("      password {\n")
			content.WriteString(fmt.Sprintf("        value = \"%s\"\n", osConfig.machine.Auth.Password))
			content.WriteString("        encrypted = false\n")
			content.WriteString("      }\n")
			content.WriteString("    }\n")

			content.WriteString("  }\n\n")
		}
	}

	content.WriteString("}\n")

	return os.WriteFile(machinesHCLPath, []byte(content.String()), 0o644)
}

// Note: variables.hcl is now generated by spooky project init

// Note: actions.hcl is now generated by spooky project init

// generateSpookyRunner creates a minimal Go container to run Spooky commands
func (bg *BundleGenerator) generateSpookyRunner(profile *profiles.Profile, outputPath string) error {
	runnerDir := filepath.Join(outputPath, "spooky-runner")
	if err := os.MkdirAll(runnerDir, 0o755); err != nil {
		return fmt.Errorf("failed to create spooky-runner directory: %w", err)
	}

	// Generate Containerfile for the runner
	containerfilePath := filepath.Join(runnerDir, "Containerfile")
	content := `FROM golang:1.24-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

# Set working directory
WORKDIR /build

# Copy the entire spooky source code
COPY spooky-source /build/spooky

# Build spooky binary
WORKDIR /build/spooky
RUN go build -o /build/spooky-bin ./main.go

# Build spooky-test-bundle-generator
WORKDIR /build/spooky/tools/spooky-test-bundle-generator
RUN go build -o /build/bundle-generator ./main.go

# Final stage
FROM alpine:3.19

# Copy built binaries
COPY --from=builder /build/spooky-bin /usr/local/bin/spooky
COPY --from=builder /build/bundle-generator /usr/local/bin/spooky-test-bundle-generator

# Copy spooky project files
COPY spooky-project /app/spooky-project

# Set working directory
WORKDIR /app

# Set permissions
RUN chmod +x /usr/local/bin/spooky /usr/local/bin/spooky-test-bundle-generator

# Default command
CMD ["/usr/local/bin/spooky"]
`

	if err := os.WriteFile(containerfilePath, []byte(content), 0o644); err != nil {
		return fmt.Errorf("failed to write Containerfile: %w", err)
	}

	return nil
}

// generateTemplates creates the templates directory and files
func (bg *BundleGenerator) generateTemplates(profile *profiles.Profile, projectDir string) error {
	if profile.Project.Templates == nil {
		return nil
	}

	templatesDir := filepath.Join(projectDir, "templates")
	if err := os.MkdirAll(templatesDir, 0o755); err != nil {
		return fmt.Errorf("failed to create templates directory: %w", err)
	}

	// Generate template files
	templates := []struct {
		name     string
		template *profiles.TemplateConfig
	}{
		{"nginx_config", profile.Project.Templates.NginxConfig},
		{"app_config", profile.Project.Templates.AppConfig},
		{"systemd_service", profile.Project.Templates.SystemdService},
	}

	for _, tmpl := range templates {
		if tmpl.template != nil {
			templatePath := filepath.Join(templatesDir, tmpl.template.Name)
			if err := os.WriteFile(templatePath, []byte(tmpl.template.Content), 0o644); err != nil {
				return fmt.Errorf("failed to write template %s: %w", tmpl.name, err)
			}
		}
	}

	return nil
}

// generateFiles creates the files directory and static files
func (bg *BundleGenerator) generateFiles(profile *profiles.Profile, projectDir string) error {
	if profile.Project.Files == nil {
		return nil
	}

	filesDir := filepath.Join(projectDir, "files")
	if err := os.MkdirAll(filesDir, 0o755); err != nil {
		return fmt.Errorf("failed to create files directory: %w", err)
	}

	// Generate static files
	files := []struct {
		name string
		file *profiles.FileConfig
	}{
		{"custom_facts_script", profile.Project.Files.CustomFactsScript},
		{"python_facts_script", profile.Project.Files.PythonFactsScript},
		{"test_data", profile.Project.Files.TestData},
		{"sync_config", profile.Project.Files.SyncConfig},
	}

	for _, file := range files {
		if file.file != nil {
			filePath := filepath.Join(filesDir, file.file.Name)
			if err := os.WriteFile(filePath, []byte(file.file.Content), 0o644); err != nil {
				return fmt.Errorf("failed to write file %s: %w", file.name, err)
			}
		}
	}

	return nil
}

// generateJustfile creates the justfile
func (bg *BundleGenerator) generateJustfile(profile *profiles.Profile, outputPath string) error {
	fmt.Println("🔧 Generating justfile...")

	content := fmt.Sprintf(`# Justfile for %s test bundle
# Generated from profile: %s

# Build all containers
build:
    #!/usr/bin/env bash
    echo "Building containers for all OSes..."
    for os in debian13 fedora42 arch alpine319 opensuse156; do
        if [ -d "containers/$os" ]; then
            echo "Building $os container..."
            podman build -t spooky-test-$os containers/$os
        fi
    done

# Start all containers
start:
    #!/usr/bin/env bash
    echo "Starting containers for all OSes..."
    for os in debian13 fedora42 arch alpine319 opensuse156; do
        if [ -d "containers/$os" ]; then
            echo "Starting $os container..."
            podman run -d --name spooky-test-$os --network spooky-test-net spooky-test-$os
        fi
    done

# Stop all containers
stop:
    #!/usr/bin/env bash
    echo "Stopping containers for all OSes..."
    for os in debian13 fedora42 arch alpine319 opensuse156; do
        podman stop spooky-test-$os 2>/dev/null || true
        podman rm spooky-test-$os 2>/dev/null || true
    done

# Clean up all containers and images
cleanup:
    #!/usr/bin/env bash
    echo "Cleaning up containers and images..."
    for os in debian13 fedora42 arch alpine319 opensuse156; do
        podman stop spooky-test-$os 2>/dev/null || true
        podman rm spooky-test-$os 2>/dev/null || true
        podman rmi spooky-test-$os 2>/dev/null || true
    done
    podman rmi spooky-test-runner 2>/dev/null || true
    podman network rm spooky-test-net 2>/dev/null || true

# Setup test network
setup-network:
    #!/usr/bin/env bash
    echo "Setting up test network..."
    podman network create --subnet 10.0.100.0/24 spooky-test-net 2>/dev/null || true

# Build Spooky runner container
build-runner:
    #!/usr/bin/env bash
    echo "Building Spooky runner container..."
    # Copy entire spooky source code to runner directory
    cp -r ../../.. spooky-runner/spooky-source
    # Copy spooky project to runner directory
    cp -r spooky-project spooky-runner/
    # Build the runner container
    podman build -t spooky-test-runner spooky-runner

# Run Spooky tests
test: setup-network build build-runner start
    #!/usr/bin/env bash
    echo "Running Spooky tests..."
    # Run spooky commands inside the runner container
    podman run --rm --network spooky-test-net spooky-test-runner spooky facts gather

# Show container status
status:
    #!/usr/bin/env bash
    echo "Container status:"
    podman ps -a --filter name=spooky-test-

# Show logs for a specific container
logs os:
    #!/usr/bin/env bash
    echo "Logs for $os container:"
    podman logs spooky-test-$os

# Shell into a specific container
shell os:
    #!/usr/bin/env bash
    echo "Opening shell in $os container..."
    podman exec -it spooky-test-$os /bin/bash

# Help
help:
    @echo "Available commands:"
    @echo "  build        - Build all containers"
    @echo "  start        - Start all containers"
    @echo "  stop         - Stop all containers"
    @echo "  cleanup      - Clean up containers and images"
    @echo "  setup-network - Setup test network"
    @echo "  test         - Run complete test suite"
    @echo "  status       - Show container status"
    @echo "  logs <os>    - Show logs for specific OS container"
    @echo "  shell <os>   - Shell into specific OS container"
    @echo "  help         - Show this help"
`, profile.Name, profile.Name)

	return os.WriteFile(filepath.Join(outputPath, "justfile"), []byte(content), 0o755)
}

// generateReadme creates the README.md file
func (bg *BundleGenerator) generateReadme(profile *profiles.Profile, outputPath string) error {
	fmt.Println("📖 Generating README...")

	var content strings.Builder
	content.WriteString(fmt.Sprintf("# %s Test Bundle\n\n", profile.Name))
	content.WriteString(fmt.Sprintf("%s\n\n", profile.DescriptionText))
	content.WriteString("## Overview\n\n")
	content.WriteString(fmt.Sprintf("This test bundle contains cross-OS testing configurations for the Spooky project, testing %s functionality across multiple operating systems.\n\n", profile.Name))

	content.WriteString("## Supported Operating Systems\n\n")
	content.WriteString("- **Debian 13 (Trixie)** - Debian-based Linux with systemd\n")
	content.WriteString("- **Fedora 42** - Red Hat-based Linux with systemd\n")
	content.WriteString("- **Arch Linux** - Rolling release Linux with systemd\n")
	content.WriteString("- **Alpine 3.19** - Lightweight Linux with OpenRC\n")
	content.WriteString("- **openSUSE Leap 15.6** - SUSE-based Linux with systemd\n\n")

	content.WriteString("## Quick Start\n\n")
	content.WriteString("1. **Setup and run tests:**\n")
	content.WriteString("   ```bash\n")
	content.WriteString("   just test\n")
	content.WriteString("   ```\n\n")

	content.WriteString("2. **Build containers only:**\n")
	content.WriteString("   ```bash\n")
	content.WriteString("   just build\n")
	content.WriteString("   ```\n\n")

	content.WriteString("3. **Start containers:**\n")
	content.WriteString("   ```bash\n")
	content.WriteString("   just start\n")
	content.WriteString("   ```\n\n")

	content.WriteString("4. **Stop containers:**\n")
	content.WriteString("   ```bash\n")
	content.WriteString("   just stop\n")
	content.WriteString("   ```\n\n")

	content.WriteString("5. **Clean up everything:**\n")
	content.WriteString("   ```bash\n")
	content.WriteString("   just cleanup\n")
	content.WriteString("   ```\n\n")

	content.WriteString("## Container Management\n\n")
	content.WriteString("### View Status\n")
	content.WriteString("```bash\n")
	content.WriteString("just status\n")
	content.WriteString("```\n\n")

	content.WriteString("### View Logs\n")
	content.WriteString("```bash\n")
	content.WriteString("just logs debian13    # View Debian 13 logs\n")
	content.WriteString("just logs fedora42    # View Fedora 42 logs\n")
	content.WriteString("just logs arch        # View Arch Linux logs\n")
	content.WriteString("just logs alpine319   # View Alpine 3.19 logs\n")
	content.WriteString("just logs opensuse156 # View openSUSE Leap 15.6 logs\n")
	content.WriteString("```\n\n")

	content.WriteString("### Shell Access\n")
	content.WriteString("```bash\n")
	content.WriteString("just shell debian13    # Shell into Debian 13 container\n")
	content.WriteString("just shell fedora42    # Shell into Fedora 42 container\n")
	content.WriteString("just shell arch        # Shell into Arch Linux container\n")
	content.WriteString("just shell alpine319   # Shell into Alpine 3.19 container\n")
	content.WriteString("just shell opensuse156 # Shell into openSUSE Leap 15.6 container\n")
	content.WriteString("```\n\n")

	content.WriteString("## Network Configuration\n\n")
	content.WriteString("All containers are connected to a custom Podman network (`spooky-test-net`) with the subnet `10.0.100.0/24`.\n\n")

	content.WriteString("### Container IPs\n")
	content.WriteString("- Debian 13: 10.0.100.11\n")
	content.WriteString("- Fedora 42: 10.0.100.12\n")
	content.WriteString("- Arch Linux: 10.0.100.13\n")
	content.WriteString("- Alpine 3.19: 10.0.100.14\n")
	content.WriteString("- openSUSE Leap 15.6: 10.0.100.15\n\n")

	content.WriteString("## Spooky Project\n\n")
	content.WriteString("The Spooky project configuration is located in the `spooky-project/` directory:\n\n")
	content.WriteString("- `project.hcl` - Main project configuration\n")
	content.WriteString("- `machines.hcl` - Machine definitions for all OSes\n")
	content.WriteString("- `variables.hcl` - Test variables and configuration\n")
	content.WriteString("- `actions.hcl` - Test actions and workflows\n")
	content.WriteString("- `templates/` - Template files for testing\n")
	content.WriteString("- `files/` - Static files for testing\n\n")

	content.WriteString("## Running Spooky Commands\n\n")
	content.WriteString("```bash\n")
	content.WriteString("cd spooky-project\n\n")
	content.WriteString("# Gather facts from all machines\n")
	content.WriteString("SPOOKY_BIN=\"../../../spooky\"\n")
	content.WriteString("\"$SPOOKY_BIN\" facts gather\n\n")
	content.WriteString("# Export facts to HCL format\n")
	content.WriteString("\"$SPOOKY_BIN\" facts export\n\n")
	content.WriteString("# Validate project configuration\n")
	content.WriteString("\"$SPOOKY_BIN\" project validate .\n")
	content.WriteString("```\n\n")

	content.WriteString("## Troubleshooting\n\n")
	content.WriteString("### Container Won't Start\n")
	content.WriteString("1. Check if the network exists: `podman network ls`\n")
	content.WriteString("2. Create network: `just setup-network`\n")
	content.WriteString("3. Check container logs: `just logs <os>`\n\n")

	content.WriteString("### SSH Connection Issues\n")
	content.WriteString("1. Verify container is running: `just status`\n")
	content.WriteString("2. Check SSH service: `just shell <os> && systemctl status sshd`\n")
	content.WriteString("3. Verify SSH config: `just shell <os> && cat /etc/ssh/sshd_config`\n\n")

	content.WriteString("### Spooky Command Failures\n")
	content.WriteString("1. Verify all containers are running: `just status`\n")
	content.WriteString("2. Check machine connectivity: `cd spooky-project && SPOOKY_BIN=\"../../../spooky\" && \"$SPOOKY_BIN\" facts gather --verbose`\n")
	content.WriteString("3. Validate project configuration: `cd spooky-project && SPOOKY_BIN=\"../../../spooky\" && \"$SPOOKY_BIN\" project validate .`\n\n")

	content.WriteString("## Cleanup\n\n")
	content.WriteString("To completely clean up all test artifacts:\n\n")
	content.WriteString("```bash\n")
	content.WriteString("just cleanup\n")
	content.WriteString("```\n\n")
	content.WriteString("This will:\n")
	content.WriteString("- Stop and remove all test containers\n")
	content.WriteString("- Remove all test images\n")
	content.WriteString("- Remove the test network\n\n")

	content.WriteString("## Profile Information\n\n")
	content.WriteString(fmt.Sprintf("- **Profile Name**: %s\n", profile.Name))
	content.WriteString(fmt.Sprintf("- **Description**: %s\n", profile.Description))
	content.WriteString("- **Generated**: 2025-01-27\n")
	if profile.Project.Variables.TestType != nil {
		content.WriteString(fmt.Sprintf("- **Test Type**: %s\n\n", *profile.Project.Variables.TestType))
	} else {
		content.WriteString("- **Test Type**: Not specified\n\n")
	}

	content.WriteString("## License\n\n")
	content.WriteString("This test bundle is part of the Spooky project testing infrastructure.\n")

	return os.WriteFile(filepath.Join(outputPath, "README.md"), []byte(content.String()), 0o644)
}

// validateGeneratedProject validates the generated Spooky project
func (bg *BundleGenerator) validateGeneratedProject(outputPath string) error {
	fmt.Println("🔍 Validating generated project...")

	// Get the current working directory (where the generator is running from)
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	// Look for spooky binary in the generator directory
	spookyPath := filepath.Join(currentDir, "spooky")
	if _, err := os.Stat(spookyPath); os.IsNotExist(err) {
		return fmt.Errorf("spooky binary not found at %s. Please build it first: just build-spooky-test-bundle-generator", spookyPath)
	}

	// Validate the project configuration
	projectDir := filepath.Join(outputPath, "spooky-project")
	cmd := exec.Command(spookyPath, "project", "validate", projectDir)
	// Don't change working directory - run from current directory

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("project validation failed: %w\nOutput: %s", err, string(output))
	}

	fmt.Printf("✅ Project validation successful: %s", string(output))
	return nil
}
