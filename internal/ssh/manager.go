package ssh

import (
	"fmt"
	"time"

	"spooky/internal/logging"
	"spooky/internal/ssh/acting"
	"spooky/internal/ssh/authentication"
	"spooky/internal/ssh/client"
	"spooky/internal/ssh/connection_pool"
	"spooky/internal/ssh/keys"
	"spooky/internal/ssh/types"
)

// Manager implements SSHManager interface
type Manager struct {
	config                *types.Config
	clientManager         client.ClientManager
	authenticationManager authentication.AuthenticationEngine
	connectionPoolManager connection_pool.ConnectionPool
	actingManager         acting.ActingEngine
	keyManager            keys.SSHKeyManager
	logger                logging.Logger
}

// NewManager creates a new SSH manager
func NewManager(
	config *types.Config,
	clientManager client.ClientManager,
	authenticationManager authentication.AuthenticationEngine,
	connectionPoolManager connection_pool.ConnectionPool,
	actingManager acting.ActingEngine,
	keyManager keys.SSHKeyManager,
	logger logging.Logger,
) *Manager {
	return &Manager{
		config:                config,
		clientManager:         clientManager,
		authenticationManager: authenticationManager,
		connectionPoolManager: connectionPoolManager,
		actingManager:         actingManager,
		keyManager:            keyManager,
		logger:                logger,
	}
}

// Connect establishes an SSH connection
func (m *Manager) Connect(host string, config *types.SSHConfig) (*types.SSHConnection, error) {
	return m.clientManager.Connect(host, config)
}

// ExecuteCommand executes a command on the SSH connection
func (m *Manager) ExecuteCommand(connection *types.SSHConnection, command string) (*types.CommandResult, error) {
	return m.clientManager.ExecuteCommand(connection, command)
}

// ExecuteScript executes a script on the SSH connection
func (m *Manager) ExecuteScript(connection *types.SSHConnection, script string) (*types.CommandResult, error) {
	return m.clientManager.ExecuteScript(connection, script)
}

// CloseConnection closes an SSH connection
func (m *Manager) CloseConnection(connection *types.SSHConnection) error {
	return m.clientManager.CloseConnection(connection)
}

// GetConnection gets a connection from the pool
func (m *Manager) GetConnection(host string) (*types.SSHConnection, error) {
	return m.connectionPoolManager.GetConnection(host)
}

// ReturnConnection returns a connection to the pool
func (m *Manager) ReturnConnection(connection *types.SSHConnection) error {
	return m.connectionPoolManager.ReturnConnection(connection)
}

// CloseAllConnections closes all connections in the pool
func (m *Manager) CloseAllConnections() error {
	return m.connectionPoolManager.CloseAllConnections()
}

// Authenticate authenticates an SSH connection
func (m *Manager) Authenticate(connection *types.SSHConnection, auth *types.AuthenticationConfig) error {
	return m.authenticationManager.Authenticate(connection, auth)
}

// ValidateAuthentication validates authentication configuration
func (m *Manager) ValidateAuthentication(auth *types.AuthenticationConfig) error {
	return m.authenticationManager.ValidateAuthentication(auth)
}

// ExecuteAction executes an action on the SSH connection
func (m *Manager) ExecuteAction(connection *types.SSHConnection, action *types.SSHAction) (*types.ActionResult, error) {
	return m.actingManager.ExecuteAction(connection, action)
}

// ExecuteTemplate executes a template action on the SSH connection
func (m *Manager) ExecuteTemplate(connection *types.SSHConnection, template *types.TemplateAction) (*types.ActionResult, error) {
	return m.actingManager.ExecuteTemplate(connection, template)
}

// SetDefaultTimeout sets the default timeout for connections
func (m *Manager) SetDefaultTimeout(timeout time.Duration) error {
	if timeout <= 0 {
		return fmt.Errorf("timeout must be positive")
	}
	m.config.DefaultTimeout = timeout
	return nil
}

// SetMaxConnections sets the maximum number of connections
func (m *Manager) SetMaxConnections(max int) error {
	if max <= 0 {
		return fmt.Errorf("max connections must be positive")
	}
	m.config.MaxConnections = max
	return nil
}

// EnableConnectionPooling enables or disables connection pooling
func (m *Manager) EnableConnectionPooling(enabled bool) error {
	m.config.EnableConnectionPooling = enabled
	return nil
}

// TestConnection tests SSH connectivity
func (m *Manager) TestConnection(host string) error {
	return m.clientManager.TestConnection(host)
}

// GetConnectionStats gets connection statistics
func (m *Manager) GetConnectionStats() *types.ConnectionStats {
	return &types.ConnectionStats{
		PoolStats: m.connectionPoolManager.GetStats(),
	}
}

// Close closes the SSH manager
func (m *Manager) Close() error {
	// Close all connections
	if err := m.CloseAllConnections(); err != nil {
		return fmt.Errorf("failed to close connections: %w", err)
	}

	// Close client manager
	if err := m.clientManager.Close(); err != nil {
		return fmt.Errorf("failed to close client manager: %w", err)
	}

	m.logger.Info("SSH manager closed")
	return nil
}

// Coordinator integration methods
func (m *Manager) ConnectToMachine(machine *types.Machine) (*types.SSHConnection, error) {
	// Create SSH config from machine
	config := &types.SSHConfig{
		Host:     machine.Host,
		Port:     machine.Port,
		Username: machine.Username,
		Timeout:  m.config.DefaultTimeout,
	}

	return m.Connect(machine.Host, config)
}

func (m *Manager) ExecuteActionOnMachine(machine *types.Machine, action *types.SSHAction) (*types.ActionResult, error) {
	// Connect to machine
	connection, err := m.ConnectToMachine(machine)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to machine: %w", err)
	}
	defer m.CloseConnection(connection)

	// Execute action
	return m.ExecuteAction(connection, action)
}
