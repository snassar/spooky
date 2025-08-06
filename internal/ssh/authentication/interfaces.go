package authentication

import (
	"spooky/internal/ssh/types"
)

// AuthenticationEngine defines the interface for authentication operations
type AuthenticationEngine interface {
	Authenticate(connection *types.SSHConnection, auth *types.AuthenticationConfig) error
	ValidateAuthentication(auth *types.AuthenticationConfig) error
	GetSupportedMethods() []string
}
