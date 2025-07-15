package ssh

import (
	"spooky/internal/config"

	"golang.org/x/crypto/ssh"
)

// SSHClient represents an SSH client connection
// Wraps the underlying *ssh.Client and the associated *config.Server
// for convenience and method extensions.
type SSHClient struct {
	Client *ssh.Client
	Server *config.Server
}
