// Package ssh provides SSH management functionality for the spooky codebase.
// This package implements the SSHManager interface for coordinating SSH operations.
package ssh

import (
	"context"
	"fmt"
	"net"
	"time"

	spookyinterfaces "spooky/internal/interfaces"
	spookytypes "spooky/internal/types"
	spookytypesactions "spooky/internal/types/actions"
	spookytypeslogging "spooky/internal/types/logging"
	spookytypesmachines "spooky/internal/types/machines"
	spookyschemas "spooky/internal/types/schemas"
	spookytypesssh "spooky/internal/types/ssh"
)

// Manager implements the SSHManager interface
type Manager struct {
	client *Client
	logger spookytypeslogging.Logger
}

// NewManager creates a new SSH manager
func NewManager(logger spookytypeslogging.Logger) spookyinterfaces.SSHManager {
	// Create a default client configuration
	config := &spookytypes.ClientConfig{
		DefaultPort:        22,
		DefaultTimeout:     30 * time.Second,
		MaxConnections:     10,
		MaxRetryAttempts:   3,
		RetryDelay:         5 * time.Second,
		KeepaliveInterval:  60 * time.Second,
		KeepaliveCount:     3,
		EnableCompression:  false,
		EnableKeepalive:    true,
		StrictHostKeyCheck: true,
	}

	// Create Client directly with the logger
	client := NewClient(config, logger)

	return &Manager{
		client: client,
		logger: logger,
	}
}

// CreateClient creates a new SSH client with the given configuration
func (m *Manager) CreateClient(ctx context.Context, config *spookytypes.ClientConfig) (*spookytypes.Client, error) {
	// For now, we'll use the existing Client
	// In a more sophisticated implementation, we might create different client types
	// based on the configuration
	_ = NewClient(config, m.logger)

	// Convert Client to the interface type
	sshClient := &spookytypes.Client{
		Config:    config,
		Status:    spookytypesssh.ClientStatusInitialized,
		Connected: false,
	}

	return sshClient, nil
}

// Connect establishes an SSH connection to the given host
func (m *Manager) Connect(ctx context.Context, request *spookytypes.ConnectionRequest) (*spookytypes.ConnectionResult, error) {
	m.logger.Debug("Establishing SSH connection", map[string]interface{}{
		"host": request.Host,
		"port": request.Port,
		"user": request.User,
	})

	// Use the Client to establish the connection
	connectionResult, err := m.client.Connect(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to establish SSH connection: %w", err)
	}

	return connectionResult, nil
}

// Authenticate authenticates with the given credentials
func (m *Manager) Authenticate(ctx context.Context, connection *spookytypes.Connection, auth *spookytypes.Authentication) (*spookytypes.AuthenticationResult, error) {
	// For now, authentication is handled during connection establishment
	// In a more sophisticated implementation, we might support re-authentication
	return &spookytypes.AuthenticationResult{
		Success: true,
	}, nil
}

// CreateSession creates a new SSH session
func (m *Manager) CreateSession(ctx context.Context, connection *spookytypes.Connection) (*spookytypes.Session, error) {
	m.logger.Debug("Creating SSH session", map[string]interface{}{
		"host": connection.Host,
		"port": connection.Port,
		"user": connection.User,
	})

	// Create a session using the connection information
	session := &spookytypes.Session{
		SessionID:  fmt.Sprintf("%s-%d-%d", connection.Host, connection.Port, time.Now().UnixNano()),
		Connection: connection,
		Status:     spookytypesssh.SessionStatusCreated,
		StartedAt:  time.Now(),
	}

	return session, nil
}

// RunCommand runs a command via SSH
func (m *Manager) RunCommand(ctx context.Context, session *spookytypes.Session, command *spookytypes.SSHCommand) (*spookytypes.SSHCommandResult, error) {
	m.logger.Debug("Running SSH command", map[string]interface{}{
		"session_id": session.SessionID,
		"command":    command.Command,
		"host":       session.Connection.Host,
	})

	// Use the Client to run the command
	commandResult, err := m.client.RunCommand(ctx, session.Connection, command)
	if err != nil {
		return nil, fmt.Errorf("failed to run SSH command: %w", err)
	}

	return commandResult, nil
}

// CreateActingSession creates a new SSH acting session
func (m *Manager) CreateActingSession(ctx context.Context, connection *spookytypes.Connection) (*spookytypesactions.ActingSession, error) {
	m.logger.Debug("Creating SSH acting session", map[string]interface{}{
		"host": connection.Host,
		"port": connection.Port,
		"user": connection.User,
	})

	// Create an acting session using the correct type from actions package
	now := time.Now()
	actingSession := &spookytypesactions.ActingSession{
		SessionID:     fmt.Sprintf("acting-%s-%d-%d", connection.Host, connection.Port, now.UnixNano()),
		CreatedAt:     now,
		ExpiresAt:     now.Add(24 * time.Hour), // Default 24 hour expiry
		Status:        "active",
		StartTime:     &now,
		Parallel:      false,
		MaxConcurrent: 1,
		Timeout:       30 * time.Minute,
		AllowFailures: false,
	}

	return actingSession, nil
}

// RunAction executes an action on a remote machine via SSH
func (m *Manager) RunAction(ctx context.Context, session *spookytypesactions.ActingSession, action *spookytypesactions.Action) (*spookytypesactions.ActingResult, error) {
	m.logger.Debug("Running action via SSH", map[string]interface{}{
		"session_id": session.SessionID,
		"action":     action.Name,
		"command":    action.CommandString,
	})

	// Get the target machine from the session's machine inventory
	var targetMachine *spookytypes.Machine
	if len(session.MachineInventory) > 0 {
		// Use the first machine in the inventory for now
		// In a more sophisticated implementation, this could be based on
		// action targeting rules or load balancing
		targetMachine = &session.MachineInventory[0]
	} else {
		return &spookytypesactions.ActingResult{
			ActionName:  action.Name,
			MachineName: "unknown",
			Status:      "failed",
			Error:       "no target machine available in session",
			StartTime:   time.Now(),
			EndTime:     time.Now(),
		}, nil
	}

	// Create a connection request for the action
	connectionRequest := &spookytypes.ConnectionRequest{
		Host:     targetMachine.Host,
		Port:     targetMachine.Port,
		User:     targetMachine.User,
		Password: targetMachine.Password,
		KeyPath:  targetMachine.KeyFile,
		Timeout:  session.Timeout,
	}

	// Establish connection
	connectionResult, err := m.Connect(ctx, connectionRequest)
	if err != nil {
		return &spookytypesactions.ActingResult{
			ActionName:  action.Name,
			MachineName: targetMachine.Hostname,
			Status:      "failed",
			Error:       fmt.Sprintf("connection failed: %v", err),
			StartTime:   time.Now(),
			EndTime:     time.Now(),
		}, nil
	}

	// Create SSH session
	sshSession, err := m.CreateSession(ctx, connectionResult.Connection)
	if err != nil {
		return &spookytypesactions.ActingResult{
			ActionName:  action.Name,
			MachineName: targetMachine.Hostname,
			Status:      "failed",
			Error:       fmt.Sprintf("session creation failed: %v", err),
			StartTime:   time.Now(),
			EndTime:     time.Now(),
		}, nil
	}

	// Prepare command
	commandStr := action.CommandString
	if action.Command != nil {
		commandStr = action.Command.Command
	}

	// Create SSH command
	sshCommand := &spookytypes.SSHCommand{
		Command:       commandStr,
		Args:          action.Command.Args,
		WorkingDir:    action.WorkingDir,
		Environment:   action.Environment,
		Timeout:       time.Duration(action.Timeout) * time.Second,
		CaptureOutput: true,
	}

	// Run command
	startTime := time.Now()
	commandResult, err := m.RunCommand(ctx, sshSession, sshCommand)
	endTime := time.Now()

	// Create acting result
	actingResult := &spookytypesactions.ActingResult{
		ActionName:  action.Name,
		MachineName: targetMachine.Hostname,
		StartTime:   startTime,
		EndTime:     endTime,
		Duration:    endTime.Sub(startTime),
	}

	if err != nil {
		actingResult.Status = "failed"
		actingResult.Error = err.Error()
	} else {
		actingResult.Status = "success"
		actingResult.ExitCode = commandResult.ExitCode
		actingResult.Stdout = commandResult.Stdout
		actingResult.Stderr = commandResult.Stderr
	}

	return actingResult, nil
}

// CollectResults collects results from multiple acting sessions
func (m *Manager) CollectResults(ctx context.Context, sessions []*spookytypesactions.ActingSession) ([]*spookytypesactions.ActingResult, error) {
	var results []*spookytypesactions.ActingResult

	for _, session := range sessions {
		// Collect results from completed sessions
		if session.Status == "completed" {
			// Create a result based on session information
			result := &spookytypesactions.ActingResult{
				ActionName:  session.Metadata["action_name"],  // Get from metadata
				MachineName: session.Metadata["machine_name"], // Get from metadata
				Status:      "success",
			}

			// Handle time fields (they are pointers)
			if session.StartTime != nil {
				result.StartTime = *session.StartTime
			} else {
				result.StartTime = time.Now()
			}

			if session.EndTime != nil {
				result.EndTime = *session.EndTime
				if session.StartTime != nil {
					result.Duration = session.EndTime.Sub(*session.StartTime)
				}
			} else {
				result.EndTime = time.Now()
			}

			results = append(results, result)
		}
	}

	return results, nil
}

// PingMachine tests SSH connectivity to a machine
func (m *Manager) PingMachine(ctx context.Context, machine *spookytypes.Machine) (*spookytypes.MachineStatus, error) {
	m.logger.Debug("Pinging machine", map[string]interface{}{
		"hostname": machine.Hostname,
		"host":     machine.Host,
		"port":     machine.Port,
		"user":     machine.User,
	})

	// Create connection request
	connectionRequest := &spookytypes.ConnectionRequest{
		Host:     machine.Host,
		Port:     machine.Port,
		User:     machine.User,
		Password: machine.Password,
		KeyPath:  machine.KeyFile,
		Timeout:  10 * time.Second, // Short timeout for ping
	}

	// Test DNS resolution first
	if err := m.testDNSResolution(machine.Host); err != nil {
		return &spookytypes.MachineStatus{
			Machine:   machine,
			Status:    "unreachable",
			Error:     fmt.Sprintf("DNS resolution failed: %v", err),
			LastCheck: time.Now(),
		}, nil
	}

	// Test ICMP reachability (if supported)
	if err := m.testICMPReachability(machine.Host); err != nil {
		m.logger.Debug("ICMP test failed, continuing with SSH test", map[string]interface{}{
			"host":  machine.Host,
			"error": err.Error(),
		})
	}

	// Test SSH connectivity
	startTime := time.Now()
	_, err := m.Connect(ctx, connectionRequest)
	endTime := time.Now()

	latency := int(endTime.Sub(startTime).Milliseconds())

	status := &spookytypes.MachineStatus{
		Machine:   machine,
		LastCheck: time.Now(),
		Latency:   latency,
	}

	if err != nil {
		status.Status = "unreachable"
		status.Error = err.Error()
	} else {
		status.Status = "reachable"
		// Update machine connectivity status
		if machine.Connectivity == nil {
			machine.Connectivity = &spookytypesmachines.MachineConnectivity{}
		}
		machine.Connectivity.SSHReachable = true
		machine.Connectivity.LastSSHCheck = time.Now()
		machine.Connectivity.SSHLatency = latency
	}

	return status, nil
}

// testDNSResolution tests DNS resolution for a host
func (m *Manager) testDNSResolution(host string) error {
	// Use net.LookupHost to test DNS resolution
	_, err := net.LookupHost(host)
	if err != nil {
		return fmt.Errorf("DNS resolution failed for %s: %w", host, err)
	}
	return nil
}

// testICMPReachability tests ICMP reachability for a host
func (m *Manager) testICMPReachability(host string) error {
	// Use net.DialTimeout to test basic connectivity
	// This is a lightweight alternative to ICMP ping
	conn, err := net.DialTimeout("tcp", host+":80", 5*time.Second)
	if err != nil {
		// Try port 22 (SSH) as fallback
		conn, err = net.DialTimeout("tcp", host+":22", 5*time.Second)
		if err != nil {
			return fmt.Errorf("connectivity test failed for %s: %w", host, err)
		}
	}
	if conn != nil {
		conn.Close()
	}
	return nil
}

// TransferFile transfers a file via SSH
func (m *Manager) TransferFile(ctx context.Context, session *spookytypes.Session, transfer *spookytypes.FileTransfer) (*spookytypes.FileTransferResult, error) {
	m.logger.Debug("Transferring file via SSH", map[string]interface{}{
		"session_id":  session.SessionID,
		"local_path":  transfer.LocalPath,
		"remote_path": transfer.RemotePath,
		"direction":   transfer.Direction,
		"mode":        transfer.Mode,
	})

	// Use the Client's file transfer manager
	fileTransferManager := m.client.GetFileTransferManager()
	transferResult, err := fileTransferManager.TransferFile(ctx, session.Connection, transfer)
	if err != nil {
		return nil, fmt.Errorf("failed to transfer file: %w", err)
	}

	return transferResult, nil
}

// ValidateConnection validates SSH connection parameters
func (m *Manager) ValidateConnection(ctx context.Context, request *spookytypes.ConnectionRequest) (*spookytypes.ValidationResult, error) {
	// Basic validation of connection parameters
	if request.Host == "" {
		return &spookytypes.ValidationResult{
			Valid: false,
			Errors: []spookyschemas.SchemaError{
				*spookyschemas.NewSchemaError("connection", "ssh", "host is required"),
			},
		}, nil
	}

	if request.Port <= 0 || request.Port > 65535 {
		return &spookytypes.ValidationResult{
			Valid: false,
			Errors: []spookyschemas.SchemaError{
				*spookyschemas.NewSchemaError("connection", "ssh", "port must be between 1 and 65535"),
			},
		}, nil
	}

	if request.User == "" {
		return &spookytypes.ValidationResult{
			Valid: false,
			Errors: []spookyschemas.SchemaError{
				*spookyschemas.NewSchemaError("connection", "ssh", "user is required"),
			},
		}, nil
	}

	// Check if authentication method is provided
	if request.Password == "" && request.KeyPath == "" {
		return &spookytypes.ValidationResult{
			Valid: false,
			Errors: []spookyschemas.SchemaError{
				*spookyschemas.NewSchemaError("connection", "ssh", "either password or key_path must be provided"),
			},
		}, nil
	}

	return &spookytypes.ValidationResult{
		Valid: true,
	}, nil
}

// ValidateAuthentication validates SSH authentication parameters
func (m *Manager) ValidateAuthentication(ctx context.Context, auth *spookytypes.Authentication) (*spookytypes.ValidationResult, error) {
	// Basic validation of authentication parameters
	if auth == nil {
		return &spookytypes.ValidationResult{
			Valid: false,
			Errors: []spookyschemas.SchemaError{
				*spookyschemas.NewSchemaError("authentication", "ssh", "authentication configuration is required"),
			},
		}, nil
	}

	// Validate based on authentication method
	switch auth.Method {
	case spookytypesssh.AuthMethodPassword:
		if auth.Password == "" {
			return &spookytypes.ValidationResult{
				Valid: false,
				Errors: []spookyschemas.SchemaError{
					*spookyschemas.NewSchemaError("authentication", "ssh", "password is required for password authentication"),
				},
			}, nil
		}
	case spookytypesssh.AuthMethodPublicKey:
		if auth.KeyPath == "" {
			return &spookytypes.ValidationResult{
				Valid: false,
				Errors: []spookyschemas.SchemaError{
					*spookyschemas.NewSchemaError("authentication", "ssh", "key_path is required for public key authentication"),
				},
			}, nil
		}
	default:
		return &spookytypes.ValidationResult{
			Valid: false,
			Errors: []spookyschemas.SchemaError{
				*spookyschemas.NewSchemaError("authentication", "ssh", fmt.Sprintf("unsupported authentication method: %s", auth.Method)),
			},
		}, nil
	}

	return &spookytypes.ValidationResult{
		Valid: true,
	}, nil
}

// GetConnectionPool returns the connection pool
func (m *Manager) GetConnectionPool() *spookytypes.ConnectionPool {
	// Return the connection pool from the Client
	// This is a simplified implementation - in a real implementation,
	// we might want to expose more pool management functionality
	return &spookytypes.ConnectionPool{
		MaxConnections:     m.client.config.MaxConnections,
		MaxIdleConnections: m.client.config.MaxIdleConnections,
		IdleTimeout:        m.client.config.IdleTimeout,
		ConnectionTimeout:  m.client.config.DefaultTimeout,
	}
}

// Close closes all SSH connections
func (m *Manager) Close(ctx context.Context) error {
	m.logger.Info("Closing SSH manager and all connections")

	// Close the Client
	err := m.client.Close(ctx)
	if err != nil {
		return fmt.Errorf("failed to close SSH client: %w", err)
	}

	return nil
}
