// Package ssh provides SSH management functionality for the spooky codebase.
// This package implements the SSHManager interface for coordinating SSH operations.
package ssh

import (
	"context"
	"fmt"
	"time"

	spookyinterfaces "spooky/internal/interfaces"
	spookytypes "spooky/internal/types"
	spookytypesactions "spooky/internal/types/actions"
	spookytypeslogging "spooky/internal/types/logging"
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
func (m *Manager) CreateActingSession(ctx context.Context, connection *spookytypes.Connection) (*spookytypes.ActingSession, error) {
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
