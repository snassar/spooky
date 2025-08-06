package types

import (
	"time"
)

// AuthenticationMethod represents the type of authentication
type AuthenticationMethod string

const (
	AuthMethodPassword AuthenticationMethod = "password"
	AuthMethodKey      AuthenticationMethod = "key"
	AuthMethodMixed    AuthenticationMethod = "mixed"
)

// AuthenticationConfig represents authentication configuration
type AuthenticationConfig struct {
	Method     AuthenticationMethod `hcl:"method"`
	Username   string               `hcl:"username"`
	Password   string               `hcl:"password,optional"`
	KeyFile    string               `hcl:"key_file,optional"`
	Passphrase string               `hcl:"passphrase,optional"`
	Timeout    time.Duration        `hcl:"timeout,optional"`
}

// AuthenticationResult represents the result of authentication
type AuthenticationResult struct {
	Success   bool                 `hcl:"success"`
	Method    AuthenticationMethod `hcl:"method"`
	Username  string               `hcl:"username"`
	Timestamp time.Time            `hcl:"timestamp"`
	Error     string               `hcl:"error,optional"`
}

// AuthenticationCache represents cached authentication results
type AuthenticationCache struct {
	Host      string               `hcl:"host"`
	Username  string               `hcl:"username"`
	Method    AuthenticationMethod `hcl:"method"`
	Success   bool                 `hcl:"success"`
	ExpiresAt time.Time            `hcl:"expires_at"`
	LastUsed  time.Time            `hcl:"last_used"`
}
