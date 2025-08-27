package profiles

import (
	"fmt"
	"os"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

// Parser handles parsing of HCL profile files
type Parser struct{}

// NewParser creates a new profile parser
func NewParser() *Parser {
	return &Parser{}
}

// ParseFile parses a profile from an HCL file
func (p *Parser) ParseFile(filePath string) (*Profile, error) {
	// Read the file
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read profile file %s: %w", filePath, err)
	}

	// Parse HCL
	file, diags := hclsyntax.ParseConfig(fileContent, filePath, hcl.Pos{Line: 1, Column: 1})
	if diags.HasErrors() {
		return nil, fmt.Errorf("failed to parse HCL file %s: %v", filePath, diags)
	}

	// Create schema context for profile blocks
	schema := &hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{
			{
				Type:       "profile",
				LabelNames: []string{"name", "description"},
			},
		},
	}

	// Decode the body to get the profile block
	bodyContent, diags := file.Body.Content(schema)
	if diags.HasErrors() {
		return nil, fmt.Errorf("failed to decode HCL body: %v", diags)
	}

	if len(bodyContent.Blocks) == 0 {
		return nil, fmt.Errorf("no profile block found in %s", filePath)
	}

	if len(bodyContent.Blocks) > 1 {
		return nil, fmt.Errorf("multiple profile blocks found in %s", filePath)
	}

	profileBlock := bodyContent.Blocks[0]

	// Decode the profile block body into Profile struct
	var profile Profile
	if diags := gohcl.DecodeBody(profileBlock.Body, nil, &profile); diags.HasErrors() {
		return nil, fmt.Errorf("failed to decode profile from %s: %v", filePath, diags)
	}

	// Set the name and description from the block labels
	profile.Name = profileBlock.Labels[0]
	if len(profileBlock.Labels) > 1 {
		profile.Description = profileBlock.Labels[1]
	}

	fmt.Printf("📄 Read profile file: %s (%d bytes)\n", filePath, len(fileContent))

	return &profile, nil
}

// ValidateProfile validates a parsed profile
func (p *Parser) ValidateProfile(profile *Profile) error {
	if profile.Name == "" {
		return fmt.Errorf("profile name is required")
	}

	if profile.Description == "" {
		return fmt.Errorf("profile description is required")
	}

	if profile.Containers == nil {
		return fmt.Errorf("containers configuration is required")
	}

	// Validate that at least one container is configured
	containerCount := 0
	if profile.Containers.Debian13 != nil {
		containerCount++
		if err := p.validateContainer(profile.Containers.Debian13, "debian13"); err != nil {
			return err
		}
	}
	if profile.Containers.Fedora42 != nil {
		containerCount++
		if err := p.validateContainer(profile.Containers.Fedora42, "fedora42"); err != nil {
			return err
		}
	}
	if profile.Containers.Arch != nil {
		containerCount++
		if err := p.validateContainer(profile.Containers.Arch, "arch"); err != nil {
			return err
		}
	}
	if profile.Containers.Alpine319 != nil {
		containerCount++
		if err := p.validateContainer(profile.Containers.Alpine319, "alpine319"); err != nil {
			return err
		}
	}
	if profile.Containers.Opensuse156 != nil {
		containerCount++
		if err := p.validateContainer(profile.Containers.Opensuse156, "opensuse156"); err != nil {
			return err
		}
	}

	if containerCount == 0 {
		return fmt.Errorf("at least one container configuration is required")
	}

	if profile.Project == nil {
		return fmt.Errorf("project configuration is required")
	}

	if profile.Project.Name == "" {
		return fmt.Errorf("project name is required")
	}

	if profile.Project.Machines == nil {
		return fmt.Errorf("machines configuration is required")
	}

	// Validate that machines match configured containers
	if profile.Containers.Debian13 != nil && profile.Project.Machines.Debian13 == nil {
		return fmt.Errorf("debian13 container configured but no corresponding machine")
	}
	if profile.Containers.Fedora42 != nil && profile.Project.Machines.Fedora42 == nil {
		return fmt.Errorf("fedora42 container configured but no corresponding machine")
	}
	if profile.Containers.Arch != nil && profile.Project.Machines.Arch == nil {
		return fmt.Errorf("arch container configured but no corresponding machine")
	}
	if profile.Containers.Alpine319 != nil && profile.Project.Machines.Alpine319 == nil {
		return fmt.Errorf("alpine319 container configured but no corresponding machine")
	}
	if profile.Containers.Opensuse156 != nil && profile.Project.Machines.Opensuse156 == nil {
		return fmt.Errorf("opensuse156 container configured but no corresponding machine")
	}

	return nil
}

// validateContainer validates a single container configuration
func (p *Parser) validateContainer(container *ContainerConfig, name string) error {
	if container.BaseImage == "" {
		return fmt.Errorf("%s container base_image is required", name)
	}

	if container.StaticIP == "" {
		return fmt.Errorf("%s container static_ip is required", name)
	}

	if container.SSHPort == 0 {
		return fmt.Errorf("%s container ssh_port is required", name)
	}

	if len(container.Packages) == 0 {
		return fmt.Errorf("%s container packages are required", name)
	}

	if container.SSHConfig == nil {
		return fmt.Errorf("%s container ssh_config is required", name)
	}

	return nil
}
