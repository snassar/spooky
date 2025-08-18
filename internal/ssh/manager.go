// Package ssh provides SSH management functionality for the spooky codebase.
// This package implements the SSHManager interface for coordinating SSH operations.
package ssh

import (
	"context"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"

	spookyinterfaces "spooky/internal/interfaces"
	spookytypes "spooky/internal/types"
	spookytypesactions "spooky/internal/types/actions"
	spookytypeslogging "spooky/internal/types/logging"
	spookytypesmachines "spooky/internal/types/machines"
	spookyschemas "spooky/internal/types/schemas"
	spookytypesssh "spooky/internal/types/ssh"
)

// Status constants for action results
const (
	StatusSuccess = "success"
	StatusFailed  = "failed"
)

// Manager implements the SSHManager interface
type Manager struct {
	client *ReusableSSHClient
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

	// Create ReusableSSHClient for connection reuse
	client := NewReusableSSHClient(config, logger)

	return &Manager{
		client: client,
		logger: logger,
	}
}

// CreateClient creates a new SSH client with the given configuration
func (m *Manager) CreateClient(_ context.Context, config *spookytypes.ClientConfig) (*spookytypes.Client, error) {
	// Create a new ReusableSSHClient with the provided configuration
	_ = NewReusableSSHClient(config, m.logger)

	// Convert to the interface type
	sshClient := &spookytypes.Client{
		Config:    config,
		Status:    spookytypesssh.ClientStatusInitialized,
		Connected: false,
	}

	return sshClient, nil
}

// Connect establishes an SSH connection to the given host
func (m *Manager) Connect(_ context.Context, request *spookytypes.ConnectionRequest) (*spookytypes.ConnectionResult, error) {
	m.logger.Debug("Establishing SSH connection", map[string]interface{}{
		"host": request.Host,
		"port": request.Port,
		"user": request.User,
	})

	// Convert ConnectionRequest to Machine for the ReusableSSHClient
	machine := &spookytypes.Machine{
		Hostname:   request.Host,
		Host:       request.Host,
		Port:       request.Port,
		User:       request.User,
		Password:   request.Password,
		KeyFile:    request.KeyPath,
		Passphrase: "", // TODO: Add passphrase support to ConnectionRequest
	}

	// Test connection by running a simple command
	ctx, cancel := context.WithTimeout(context.Background(), request.Timeout)
	defer cancel()

	_, err := m.client.RunCommand(ctx, machine, "echo 'connection test'")
	if err != nil {
		return &spookytypes.ConnectionResult{
			Request:       request,
			Success:       false,
			Error:         err.Error(),
			ConnectTime:   request.Timeout, // We don't have exact timing here
			RetryAttempts: 0,
			CompletedAt:   time.Now(),
		}, nil
	}

	// Create connection result
	startTime := time.Now()
	connection := &spookytypes.Connection{
		Host:          request.Host,
		Port:          request.Port,
		User:          request.User,
		Status:        spookytypesssh.ConnectionStatusConnected,
		ConnectedAt:   &startTime,
		ClientVersion: "spooky-reusable-client",
		ServerVersion: "unknown", // We don't have this info from the reusable client
		Latency:       0,         // We don't have exact latency info
	}

	result := &spookytypes.ConnectionResult{
		Connection:  connection,
		Request:     request,
		Success:     true,
		ConnectTime: time.Since(startTime),
		CompletedAt: time.Now(),
	}

	return result, nil
}

// Authenticate authenticates with the given credentials
func (m *Manager) Authenticate(_ context.Context, _ *spookytypes.Connection, _ *spookytypes.Authentication) (*spookytypes.AuthenticationResult, error) {
	// Authentication is handled during connection establishment in the ReusableSSHClient
	// The ReusableSSHClient caches authentication info and reuses it
	return &spookytypes.AuthenticationResult{
		Success: true,
	}, nil
}

// CreateSession creates a new SSH session
func (m *Manager) CreateSession(_ context.Context, connection *spookytypes.Connection) (*spookytypes.Session, error) {
	m.logger.Debug("Creating SSH session", map[string]interface{}{
		"host": connection.Host,
		"port": connection.Port,
		"user": connection.User,
	})

	// Create a session using the connection information
	// Note: The ReusableSSHClient manages sessions internally
	session := &spookytypes.Session{
		SessionID:  fmt.Sprintf("%s-%d-%d", connection.Host, connection.Port, time.Now().UnixNano()),
		Connection: connection,
		Status:     spookytypesssh.SessionStatusCreated,
		StartedAt:  time.Now(),
	}

	return session, nil
}

// RunCommand runs a command via SSH with connection reuse
func (m *Manager) RunCommand(ctx context.Context, session *spookytypes.Session, command *spookytypes.SSHCommand) (*spookytypes.SSHCommandResult, error) {
	m.logger.Debug("Running SSH command", map[string]interface{}{
		"session_id": session.SessionID,
		"command":    command.Command,
		"host":       session.Connection.Host,
	})

	// Convert session.Connection to Machine for the ReusableSSHClient
	machine := &spookytypes.Machine{
		Hostname: session.Connection.Host,
		Host:     session.Connection.Host,
		Port:     session.Connection.Port,
		User:     session.Connection.User,
		// Note: We don't have password/key info in the session.Connection
		// This is a limitation of the current interface design
	}

	// Build the full command string
	cmdStr := command.Command
	if len(command.Args) > 0 {
		cmdStr = cmdStr + " " + strings.Join(command.Args, " ")
	}

	// Run command using the ReusableSSHClient
	startTime := time.Now()
	stdout, err := m.client.RunCommand(ctx, machine, cmdStr)
	endTime := time.Now()

	// Create command result
	result := &spookytypes.SSHCommandResult{
		Command:   command,
		Session:   session,
		StartTime: startTime,
		EndTime:   endTime,
		Duration:  endTime.Sub(startTime),
	}

	if err != nil {
		result.Success = false
		result.Error = err.Error()
		result.ExitCode = -1
		result.Stderr = err.Error()
	} else {
		result.Success = true
		result.ExitCode = 0
		result.Stdout = stdout
	}

	return result, nil
}

// CreateActingSession creates a new SSH acting session
func (m *Manager) CreateActingSession(_ context.Context, connection *spookytypes.Connection) (*spookytypes.ActingSession, error) {
	m.logger.Debug("Creating SSH acting session", map[string]interface{}{
		"host": connection.Host,
		"port": connection.Port,
		"user": connection.User,
	})

	// Create an acting session using the correct type from actions package
	now := time.Now()
	actingSession := &spookytypes.ActingSession{
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

// RunAction runs an action on a remote machine via SSH with connection reuse
func (m *Manager) RunAction(ctx context.Context, session *spookytypesactions.ActingSession, action *spookytypesactions.Action) (*spookytypesactions.ActingResult, error) {
	// Log action orchestration start
	startTime := time.Now()
	m.logActionOrchestration(action, &spookytypes.Machine{Hostname: "unknown"}, session, startTime)

	// Get the target machine from the session's machine inventory
	var targetMachine *spookytypes.Machine
	if len(session.MachineInventory) > 0 {
		// Use the first machine in the inventory for now
		// In a more sophisticated implementation, this could be based on
		// action targeting rules or load balancing
		targetMachine = &session.MachineInventory[0]
	} else {
		// Create error for no target machine
		actionError := spookytypesssh.NewActionOrchestrationError(
			action.Name,
			"unknown",
			session.SessionID,
			action.CommandString,
			"no target machine available in session",
			0,
			"",
			"",
		)

		// Log the error
		m.logger.Error("Action orchestration failed", actionError, map[string]interface{}{
			"action_name": action.Name,
			"session_id":  session.SessionID,
			"error_type":  actionError.BaseError.ErrorType,
			"error_code":  actionError.BaseError.ErrorCode,
		})

		return &spookytypesactions.ActingResult{
			ActionName:  action.Name,
			MachineName: "unknown",
			Status:      "failed",
			Error:       actionError.BaseError.ErrorMessage,
			StartTime:   startTime,
			EndTime:     time.Now(),
		}, actionError
	}

	// Prepare command
	commandStr := action.CommandString
	if action.Command != nil {
		commandStr = action.Command.Command
		if len(action.Command.Args) > 0 {
			commandStr = commandStr + " " + strings.Join(action.Command.Args, " ")
		}
	}

	// Run command using the ReusableSSHClient
	stdout, err := m.client.RunCommand(ctx, targetMachine, commandStr)
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
		// Create detailed action orchestration error
		actionError := spookytypesssh.NewActionOrchestrationError(
			action.Name,
			targetMachine.Hostname,
			session.SessionID,
			commandStr,
			err.Error(),
			1, // Default exit code for SSH errors
			stdout,
			"", // stderr not available from RunCommand
		)

		// Set additional error context
		actionError.OrchestrationTime = endTime.Sub(startTime)
		actionError.ActionFinished = &endTime
		actionError.WorkingDir = action.WorkingDir
		actionError.Timeout = time.Duration(action.Timeout) * time.Second

		// Log the detailed error
		m.logger.Error("Action orchestration failed", actionError, map[string]interface{}{
			"action_name":        actionError.ActionName,
			"machine_name":       actionError.MachineName,
			"session_id":         actionError.SessionID,
			"command":            actionError.CommandString,
			"exit_code":          actionError.ExitCode,
			"stdout":             actionError.Stdout,
			"stderr":             actionError.Stderr,
			"orchestration_time": actionError.OrchestrationTime,
			"error_type":         actionError.BaseError.ErrorType,
			"error_code":         actionError.BaseError.ErrorCode,
			"retryable":          actionError.BaseError.Retryable,
			"recoverable":        actionError.BaseError.Recoverable,
		})

		actingResult.Status = StatusFailed
		actingResult.Error = actionError.Error()
		actingResult.ExitCode = actionError.ExitCode
		actingResult.Stdout = actionError.Stdout
		actingResult.Stderr = actionError.Stderr
	} else {
		actingResult.Status = StatusSuccess
		actingResult.ExitCode = 0
		actingResult.Stdout = stdout
	}

	// Log action completion
	m.logActionCompletion(action, targetMachine, actingResult)

	return actingResult, nil
}

// CollectResults collects results from multiple acting sessions
func (m *Manager) CollectResults(_ context.Context, sessions []*spookytypesactions.ActingSession) ([]*spookytypesactions.ActingResult, error) {
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

// PingMachine tests SSH connectivity to a machine with connection reuse
func (m *Manager) PingMachine(ctx context.Context, machine *spookytypes.Machine) (*spookytypes.MachineStatus, error) {
	m.logger.Debug("Pinging machine", map[string]interface{}{
		"hostname": machine.Hostname,
		"host":     machine.Host,
		"port":     machine.Port,
		"user":     machine.User,
	})

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

	// Test SSH connectivity using the ReusableSSHClient
	startTime := time.Now()
	_, err := m.client.RunCommand(ctx, machine, "echo 'ping test'")
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

// TransferFile transfers a file to/from a remote machine via SSH
func (m *Manager) TransferFile(ctx context.Context, session *spookytypes.Session, transfer *spookytypes.FileTransfer) (*spookytypes.FileTransferResult, error) {
	m.logger.Info("Starting file transfer", map[string]interface{}{
		"session_id":  session.SessionID,
		"local_path":  transfer.LocalPath,
		"remote_path": transfer.RemotePath,
		"direction":   transfer.Direction,
	})

	// Convert session.Connection to Machine for the ReusableSSHClient
	machine := &spookytypes.Machine{
		Hostname: session.Connection.Host,
		Host:     session.Connection.Host,
		Port:     session.Connection.Port,
		User:     session.Connection.User,
	}

	// Run file transfer
	startTime := time.Now()
	var bytesTransferred int64
	var err error

	switch transfer.Direction {
	case spookytypesssh.TransferDirectionUpload:
		bytesTransferred, err = m.uploadFile(ctx, machine, transfer)
	case spookytypesssh.TransferDirectionDownload:
		bytesTransferred, err = m.downloadFile(ctx, machine, transfer)
	default:
		err = fmt.Errorf("unsupported transfer direction: %s", transfer.Direction)
	}

	endTime := time.Now()

	// Create file transfer result
	result := &spookytypes.FileTransferResult{
		Transfer:  transfer,
		Session:   session,
		StartTime: startTime,
		EndTime:   endTime,
		Duration:  endTime.Sub(startTime),
	}

	if err != nil {
		result.Success = false
		result.Error = err.Error()
		result.BytesTransferred = 0
	} else {
		result.Success = true
		result.BytesTransferred = bytesTransferred
	}

	return result, nil
}

// uploadFile uploads a file from local to remote machine
func (m *Manager) uploadFile(ctx context.Context, machine *spookytypes.Machine, transfer *spookytypes.FileTransfer) (int64, error) {
	// Read source file
	sourceData, err := os.ReadFile(transfer.LocalPath)
	if err != nil {
		return 0, fmt.Errorf("failed to read source file %s: %w", transfer.LocalPath, err)
	}

	// Create destination directory if needed
	dirCommand := fmt.Sprintf("mkdir -p $(dirname %s)", transfer.RemotePath)
	_, err = m.client.RunCommand(ctx, machine, dirCommand)
	if err != nil {
		return 0, fmt.Errorf("failed to create destination directory: %w", err)
	}

	// Upload file using cat command with stdin
	uploadCommand := fmt.Sprintf("cat > %s", transfer.RemotePath)
	_, err = m.client.RunCommandWithStdin(ctx, machine, uploadCommand, string(sourceData))
	if err != nil {
		return 0, fmt.Errorf("failed to upload file to %s: %w", transfer.RemotePath, err)
	}

	return int64(len(sourceData)), nil
}

// downloadFile downloads a file from remote to local machine
func (m *Manager) downloadFile(ctx context.Context, machine *spookytypes.Machine, transfer *spookytypes.FileTransfer) (int64, error) {
	// Read file from remote machine
	remoteData, err := m.client.RunCommand(ctx, machine, fmt.Sprintf("cat %s", transfer.RemotePath))
	if err != nil {
		return 0, fmt.Errorf("failed to read remote file %s: %w", transfer.RemotePath, err)
	}

	// Create destination directory if needed
	destDir := filepath.Dir(transfer.LocalPath)
	if err := os.MkdirAll(destDir, 0o755); err != nil {
		return 0, fmt.Errorf("failed to create destination directory: %w", err)
	}

	// Write to local destination
	err = os.WriteFile(transfer.LocalPath, []byte(remoteData), 0o600)
	if err != nil {
		return 0, fmt.Errorf("failed to write local file %s: %w", transfer.LocalPath, err)
	}

	return int64(len(remoteData)), nil
}

// logActionOrchestration logs action orchestration details
func (m *Manager) logActionOrchestration(action *spookytypesactions.Action, machine *spookytypes.Machine, session *spookytypesactions.ActingSession, startTime time.Time) {
	m.logger.Info("Starting action orchestration", map[string]interface{}{
		"action_name":  action.Name,
		"machine_name": machine.Hostname,
		"machine_host": machine.Host,
		"session_id":   session.SessionID,
		"command":      action.CommandString,
		"working_dir":  action.WorkingDir,
		"timeout":      action.Timeout,
		"parallel":     action.Parallel,
		"start_time":   startTime,
	})
}

// logActionCompletion logs action completion details
func (m *Manager) logActionCompletion(action *spookytypesactions.Action, machine *spookytypes.Machine, result *spookytypesactions.ActingResult) {
	m.logger.Info("Completed action orchestration", map[string]interface{}{
		"action_name":  action.Name,
		"machine_name": machine.Hostname,
		"status":       result.Status,
		"duration":     result.Duration,
		"exit_code":    result.ExitCode,
		"error":        result.Error,
	})
}

// ValidateConnection validates SSH connection parameters
func (m *Manager) ValidateConnection(_ context.Context, request *spookytypes.ConnectionRequest) (*spookytypes.ValidationResult, error) {
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
func (m *Manager) ValidateAuthentication(_ context.Context, auth *spookytypes.Authentication) (*spookytypes.ValidationResult, error) {
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

// GetConnectionPool returns the connection pool statistics
func (m *Manager) GetConnectionPool() *spookytypes.ConnectionPool {
	// Get connection stats from the ReusableSSHClient
	stats := m.client.GetConnectionStats()

	return &spookytypes.ConnectionPool{
		MaxConnections:     m.client.config.MaxConnections,
		MaxIdleConnections: m.client.config.MaxIdleConnections,
		IdleTimeout:        m.client.config.IdleTimeout,
		ConnectionTimeout:  m.client.config.DefaultTimeout,
		// Add additional stats from the reusable client
		ActiveConnections: stats["healthy_connections"].(int),
		IdleConnections:   stats["unhealthy_connections"].(int),
		TotalConnections:  stats["total_connections"].(int),
	}
}

// CleanupSession cleans up SSH resources for an acting session
func (m *Manager) CleanupSession(_ context.Context, session *spookytypesactions.ActingSession) error {
	if session == nil {
		return nil
	}

	m.logger.Debug("Cleaning up SSH session", map[string]interface{}{
		"session_id": session.SessionID,
	})

	// Close SSH connection if it exists
	if session.SSHConnection != nil {
		m.closeSSHConnection(session.SSHConnection)
		session.SSHConnection = nil
	}

	// Clear SSH metadata
	session.SSHMetadata = nil

	return nil
}

// closeSSHConnection closes an SSH connection
func (m *Manager) closeSSHConnection(connection *spookytypes.Connection) {
	if connection == nil {
		return
	}

	m.logger.Debug("Closing SSH connection", map[string]interface{}{
		"host": connection.Host,
		"port": connection.Port,
	})

	// Mark connection as closed
	connection.Status = spookytypesssh.ConnectionStatusClosed

	// Note: The ReusableSSHClient handles actual connection cleanup
	// This method just updates the connection status
}

// Close closes all SSH connections
func (m *Manager) Close(_ context.Context) error {
	m.logger.Info("Closing SSH manager and all connections")

	// Close the ReusableSSHClient
	err := m.client.Close()
	if err != nil {
		return fmt.Errorf("failed to close SSH client: %w", err)
	}

	return nil
}
