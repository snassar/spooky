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

	if profile.Container == nil {
		return fmt.Errorf("container configuration is required")
	}

	if profile.Container.BaseImage == "" {
		return fmt.Errorf("container base_image is required")
	}

	if profile.Container.Name == "" {
		return fmt.Errorf("container name is required")
	}

	if profile.Container.IP == "" {
		return fmt.Errorf("container ip is required")
	}

	if profile.SSH == nil {
		return fmt.Errorf("SSH configuration is required")
	}

	if profile.SpookyProject == nil {
		return fmt.Errorf("spooky_project configuration is required")
	}

	return nil
}
