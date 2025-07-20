package ssh

import (
	"spooky/internal/config"

	gossh "golang.org/x/crypto/ssh"
)

// SSHClient represents an SSH client connection
//
//revive:disable:exported
type SSHClient struct {
	config *config.Machine
	client *gossh.Client
}

//revive:enable:exported
