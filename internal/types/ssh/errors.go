package ssh

// SSHError represents an SSH-specific error
type SSHError struct {
	Code    string `hcl:"code"`
	Message string `hcl:"message"`
	Details string `hcl:"details,optional"`
}

// Error codes
const (
	ErrCodeConnectionFailed     = "connection_failed"
	ErrCodeAuthenticationFailed = "authentication_failed"
	ErrCodeCommandFailed        = "command_failed"
	ErrCodeTimeout              = "timeout"
	ErrCodeInvalidConfig        = "invalid_config"
	ErrCodeKeyNotFound          = "key_not_found"
	ErrCodeInvalidKey           = "invalid_key"
	ErrCodePoolExhausted        = "pool_exhausted"
	ErrCodeConnectionClosed     = "connection_closed"
)

// Error returns the error message
func (e *SSHError) Error() string {
	return e.Message
}

// ConnectionError represents a connection-specific error
type ConnectionError struct {
	Host    string `hcl:"host"`
	Port    int    `hcl:"port"`
	Message string `hcl:"message"`
	Code    string `hcl:"code"`
}

// Error returns the error message
func (e *ConnectionError) Error() string {
	return e.Message
}

// AuthenticationError represents an authentication-specific error
type AuthenticationError struct {
	Method  string `hcl:"method"`
	User    string `hcl:"user"`
	Message string `hcl:"message"`
	Code    string `hcl:"code"`
}

// Error returns the error message
func (e *AuthenticationError) Error() string {
	return e.Message
}

// ExecutionError represents an acting-specific error
type ExecutionError struct {
	Command  string `hcl:"command"`
	ExitCode int    `hcl:"exit_code"`
	Message  string `hcl:"message"`
	Stdout   string `hcl:"stdout"`
	Stderr   string `hcl:"stderr"`
}

// Error returns the error message
func (e *ExecutionError) Error() string {
	return e.Message
}
