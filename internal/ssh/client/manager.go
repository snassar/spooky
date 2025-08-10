package client

import (
	"fmt"
	"time"

	spookyinterfaces "spooky/internal/interfaces"
	spookylogging "spooky/internal/logging"
	spookytypes "spooky/internal/types"
)

// Manager implements ClientManager interface
type Manager struct {
	config            *spookytypes.ClientConfig
	connectionManager spookyinterfaces.ConnectionManager
	executionManager  spookyinterfaces.ExecutionManager
	hostKeyManager    spookyinterfaces.HostKeyManager
	logger            spookyinterfaces.Logger
}

// NewManager creates a new client manager
func NewManager(
	config *spookytypes.ClientConfig,
	connectionManager spookyinterfaces.ConnectionManager,
	executionManager spookyinterfaces.ExecutionManager,
	hostKeyManager spookyinterfaces.HostKeyManager,
	logger spookyinterfaces.Logger,
) *Manager {
	return &Manager{
		config:            config,
		connectionManager: connectionManager,
		executionManager:  executionManager,
		hostKeyManager:    hostKeyManager,
		logger:            logger,
	}
}

// Connect establishes an SSH connection
func (m *Manager) Connect(host string, config *spookytypes.SSHConfig) (*spookytypes.SSHConnection, error) {
	// 1. Validate host and config
	if err := m.validateConnectionParams(host, config); err != nil {
		return nil, fmt.Errorf("connection validation failed: %w", err)
	}

	// 2. Set default timeout if not specified
	if config.Timeout == 0 {
		config.Timeout = m.config.DefaultTimeout
	}

	// 3. Establish connection
	connection, err := m.connectionManager.Connect(host, config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to %s: %w", host, err)
	}

	// 4. Test connection
	if err := m.connectionManager.TestConnection(connection); err != nil {
		m.connectionManager.CloseConnection(connection)
		return nil, fmt.Errorf("connection test failed: %w", err)
	}

	m.logger.Info("SSH connection established", spookylogging.String("host", host))
	return connection, nil
}

// ExecuteCommand executes a command on the SSH connection
func (m *Manager) ExecuteCommand(connection *spookytypes.SSHConnection, command string) (*spookytypes.CommandResult, error) {
	// 1. Validate connection and command
	if err := m.validateExecutionParams(connection, command); err != nil {
		return nil, fmt.Errorf("execution validation failed: %w", err)
	}

	// 2. Execute command
	result, err := m.executionManager.ExecuteCommand(connection, command)
	if err != nil {
		return nil, fmt.Errorf("command execution failed: %w", err)
	}

	m.logger.Info("Command executed successfully",
		spookylogging.String("host", connection.Host),
		spookylogging.String("command", command),
		spookylogging.Int("exit_code", result.ExitCode))

	return result, nil
}

// ExecuteScript executes a script on the SSH connection
func (m *Manager) ExecuteScript(connection *spookytypes.SSHConnection, script string) (*spookytypes.CommandResult, error) {
	// 1. Validate connection and script
	if err := m.validateExecutionParams(connection, script); err != nil {
		return nil, fmt.Errorf("execution validation failed: %w", err)
	}

	// 2. Execute script
	result, err := m.executionManager.ExecuteScript(connection, script)
	if err != nil {
		return nil, fmt.Errorf("script execution failed: %w", err)
	}

	m.logger.Info("Script executed successfully",
		spookylogging.String("host", connection.Host),
		spookylogging.Int("exit_code", result.ExitCode))

	return result, nil
}

// CloseConnection closes an SSH connection
func (m *Manager) CloseConnection(connection *spookytypes.SSHConnection) error {
	if connection == nil {
		return nil
	}

	if err := m.connectionManager.CloseConnection(connection); err != nil {
		return fmt.Errorf("failed to close connection: %w", err)
	}

	m.logger.Info("SSH connection closed", spookylogging.String("host", connection.Host))
	return nil
}

// SetDefaultTimeout sets the default timeout for connections
func (m *Manager) SetDefaultTimeout(timeout time.Duration) error {
	if timeout <= 0 {
		return fmt.Errorf("timeout must be positive")
	}
	m.config.DefaultTimeout = timeout
	return nil
}

// SetMaxRetries sets the maximum number of retries
func (m *Manager) SetMaxRetries(max int) error {
	if max < 0 {
		return fmt.Errorf("max retries cannot be negative")
	}
	m.config.MaxRetries = max
	return nil
}

// EnableHostKeyChecking enables or disables host key checking
func (m *Manager) EnableHostKeyChecking(enabled bool) error {
	m.config.HostKeyChecking = enabled
	return nil
}

// TestConnection tests SSH connectivity
func (m *Manager) TestConnection(host string) error {
	// Create a temporary config for testing
	config := &spookytypes.SSHConfig{
		Host:    host,
		Port:    22,
		Timeout: m.config.DefaultTimeout,
	}

	// Try to connect
	connection, err := m.Connect(host, config)
	if err != nil {
		return fmt.Errorf("connection test failed: %w", err)
	}

	// Close the test connection
	defer m.CloseConnection(connection)

	m.logger.Info("Connection test successful", spookylogging.String("host", host))
	return nil
}

// GetConnectionInfo gets connection information
func (m *Manager) GetConnectionInfo(connection *spookytypes.SSHConnection) *spookytypes.ConnectionInfo {
	if connection == nil {
		return nil
	}

	return &spookytypes.ConnectionInfo{
		Host:      connection.Host,
		Port:      connection.Port,
		Username:  connection.Username,
		Connected: connection.Connected,
		CreatedAt: connection.CreatedAt,
		LastUsed:  connection.LastUsed,
	}
}

// Close closes the client manager
func (m *Manager) Close() error {
	m.logger.Info("SSH client manager closed")
	return nil
}

// Helper methods
func (m *Manager) validateConnectionParams(host string, config *spookytypes.SSHConfig) error {
	if host == "" {
		return fmt.Errorf("host cannot be empty")
	}

	if config == nil {
		return fmt.Errorf("SSH config cannot be nil")
	}

	return nil
}

func (m *Manager) validateExecutionParams(connection *spookytypes.SSHConnection, command string) error {
	if connection == nil {
		return fmt.Errorf("connection cannot be nil")
	}

	if command == "" {
		return fmt.Errorf("command cannot be empty")
	}

	return nil
}
