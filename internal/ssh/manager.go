package ssh

import (
	"fmt"
	"time"

	spookyinterfaces "spooky/internal/interfaces"
	spookytypes "spooky/internal/types"
)

// Manager implements SSHManager interface
type Manager struct {
	config                *spookytypes.SSHClientConfig
	clientManager         spookyinterfaces.ClientManager
	authenticationManager spookyinterfaces.AuthenticationEngine
	connectionPoolManager spookyinterfaces.ConnectionPool
	actingManager         spookyinterfaces.ActingEngine
	keyManager            spookyinterfaces.SSHKeyManager
	logger                spookyinterfaces.Logger
}

// NewManager creates a new SSH manager
func NewManager(
	config *spookytypes.SSHClientConfig,
	clientManager spookyinterfaces.ClientManager,
	authenticationManager spookyinterfaces.AuthenticationEngine,
	connectionPoolManager spookyinterfaces.ConnectionPool,
	actingManager spookyinterfaces.ActingEngine,
	keyManager spookyinterfaces.SSHKeyManager,
	logger spookyinterfaces.Logger,
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
func (m *Manager) Connect(host string, config *spookytypes.SSHConfig) (*spookytypes.SSHConnection, error) {
	return m.clientManager.Connect(host, config)
}

// ExecuteCommand executes a command on the SSH connection
func (m *Manager) ExecuteCommand(connection *spookytypes.SSHConnection, command string) (*spookytypes.CommandResult, error) {
	return m.clientManager.ExecuteCommand(connection, command)
}

// ExecuteScript executes a script on the SSH connection
func (m *Manager) ExecuteScript(connection *spookytypes.SSHConnection, script string) (*spookytypes.CommandResult, error) {
	return m.clientManager.ExecuteScript(connection, script)
}

// CloseConnection closes an SSH connection
func (m *Manager) CloseConnection(connection *spookytypes.SSHConnection) error {
	return m.clientManager.CloseConnection(connection)
}

// GetConnection gets a connection from the pool
func (m *Manager) GetConnection(host string) (*spookytypes.SSHConnection, error) {
	return m.connectionPoolManager.GetConnection(host)
}

// ReturnConnection returns a connection to the pool
func (m *Manager) ReturnConnection(connection *spookytypes.SSHConnection) error {
	return m.connectionPoolManager.ReturnConnection(connection)
}

// CloseAllConnections closes all connections in the pool
func (m *Manager) CloseAllConnections() error {
	return m.connectionPoolManager.CloseAllConnections()
}

// Authenticate authenticates an SSH connection
func (m *Manager) Authenticate(connection *spookytypes.SSHConnection, auth *spookytypes.AuthenticationConfig) error {
	return m.authenticationManager.Authenticate(connection, auth)
}

// ValidateAuthentication validates authentication configuration
func (m *Manager) ValidateAuthentication(auth *spookytypes.AuthenticationConfig) error {
	return m.authenticationManager.ValidateAuthentication(auth)
}

// ExecuteAction executes an action on the SSH connection
func (m *Manager) ExecuteAction(connection *spookytypes.SSHConnection, action *spookytypes.SSHAction) (*spookytypes.ActionResult, error) {
	return m.actingManager.ExecuteAction(connection, action)
}

// ExecuteTemplate executes a template action on the SSH connection
func (m *Manager) ExecuteTemplate(connection *spookytypes.SSHConnection, template *spookytypes.TemplateAction) (*spookytypes.ActionResult, error) {
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
func (m *Manager) GetConnectionStats() *spookytypes.ConnectionStats {
	return &spookytypes.ConnectionStats{
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
func (m *Manager) ConnectToMachine(machine *spookytypes.SSHMachine) (*spookytypes.SSHConnection, error) {
	// Create SSH config from machine
	config := &spookytypes.SSHConfig{
		Host:     machine.Host,
		Port:     machine.Port,
		Username: machine.Username,
		Timeout:  m.config.DefaultTimeout,
	}

	return m.Connect(machine.Host, config)
}

func (m *Manager) ExecuteActionOnMachine(machine *spookytypes.SSHMachine, action *spookytypes.SSHAction) (*spookytypes.ActionResult, error) {
	// Connect to machine
	connection, err := m.ConnectToMachine(machine)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to machine: %w", err)
	}
	defer m.CloseConnection(connection)

	// Execute action
	return m.ExecuteAction(connection, action)
}
