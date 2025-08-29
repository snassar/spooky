package utilities

import (
	"fmt"

	"spooky/internal/schemas"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
)

// ParseAuthenticationBlock parses an authentication block from HCL content
func ParseAuthenticationBlock(authBlock interface{}, machineName string) (*schemas.MachinesMachineAuthenticationV1, error) {
	auth := &schemas.MachinesMachineAuthenticationV1{}

	// Type assertion to get the block
	block, ok := authBlock.(*hcl.Block)
	if !ok {
		return nil, fmt.Errorf("invalid authentication block type for machine %s", machineName)
	}

	// Validate that the authentication block has the required method label
	if len(block.Labels) == 0 {
		return nil, fmt.Errorf("machine %s: authentication block missing required method label (password, publickey, or certificate)", machineName)
	}

	authMethod := block.Labels[0]

	switch authMethod {
	case "password":
		var passwordData struct {
			Password struct {
				Value     string `hcl:"value,attr"`
				Encrypted bool   `hcl:"encrypted,attr"`
			} `hcl:"password,block"`
		}

		if diags := gohcl.DecodeBody(block.Body, nil, &passwordData); diags.HasErrors() {
			return nil, fmt.Errorf("failed to decode password auth for machine %s: %v", machineName, diags)
		}

		auth.Password = schemas.MachinesMachineAuthenticationPasswordV1{
			Value:     passwordData.Password.Value,
			Encrypted: passwordData.Password.Encrypted,
		}

	case "publickey":
		var publicKeyData struct {
			PublicKeyPath string `hcl:"public_key_path,attr"`
			Passphrase    struct {
				Value     string `hcl:"value,attr"`
				Encrypted bool   `hcl:"encrypted,attr"`
			} `hcl:"passphrase,block"`
		}

		if diags := gohcl.DecodeBody(block.Body, nil, &publicKeyData); diags.HasErrors() {
			return nil, fmt.Errorf("failed to decode publickey auth for machine %s: %v", machineName, diags)
		}

		auth.PublicKeyPath = publicKeyData.PublicKeyPath
		auth.Passphrase = schemas.MachinesMachineAuthenticationPassphraseV1{
			Value:     publicKeyData.Passphrase.Value,
			Encrypted: publicKeyData.Passphrase.Encrypted,
		}

	case "certificate":
		var certificateData struct {
			PrivateKeyPath        string `hcl:"private_key_path,attr"`
			CertificatePath       string `hcl:"certificate_path,attr"`
			CertificatePassphrase struct {
				Value     string `hcl:"value,attr"`
				Encrypted bool   `hcl:"encrypted,attr"`
			} `hcl:"certificate_passphrase,block"`
		}

		if diags := gohcl.DecodeBody(block.Body, nil, &certificateData); diags.HasErrors() {
			return nil, fmt.Errorf("failed to decode certificate auth for machine %s: %v", machineName, diags)
		}

		auth.PrivateKeyPath = certificateData.PrivateKeyPath
		auth.CertificatePath = certificateData.CertificatePath
		auth.CertificatePassphrase = schemas.MachinesMachineAuthenticationPassphraseV1{
			Value:     certificateData.CertificatePassphrase.Value,
			Encrypted: certificateData.CertificatePassphrase.Encrypted,
		}

	default:
		return nil, fmt.Errorf("machine %s: unsupported authentication method '%s' (supported: password, publickey, certificate)", machineName, authMethod)
	}

	return auth, nil
}
