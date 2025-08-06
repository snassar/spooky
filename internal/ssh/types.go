package ssh

import (
	spookyconfig "spooky/internal/config"

	gossh "golang.org/x/crypto/ssh"
)

// SSHClient represents an SSH client connection
//
//revive:disable:exported
type SSHClient struct {
	config *spookyconfig.Machine
	client *gossh.Client
}

//revive:enable:exported
