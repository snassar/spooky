package ssh

import (
	"spooky/internal/config"

	"golang.org/x/crypto/ssh"
)

// SSHClient represents an SSH client connection
// Wraps the underlying *ssh.Client and the associated *config.Server
// for convenience and method extensions.
//
//nolint:revive // SSHClient is a reasonable name that clearly indicates its purpose
type SSHClient struct {
	Client *ssh.Client
	Server *config.Server
}
