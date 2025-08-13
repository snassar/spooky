// Package ssh provides file transfer capabilities for the spooky codebase.
// This package implements SFTP and SCP file transfer operations with progress tracking and error handling.
package ssh

import (
	"bufio"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"

	spookytypes "spooky/internal/types"
	spookytypeslogging "spooky/internal/types/logging"
	spookytypesssh "spooky/internal/types/ssh"
)

// FileTransferManager manages file transfer operations via SSH
type FileTransferManager struct {
	client *Client
	logger spookytypeslogging.Logger
	mu     sync.RWMutex
}

// NewFileTransferManager creates a new file transfer manager
func NewFileTransferManager(client *Client, logger spookytypeslogging.Logger) *FileTransferManager {
	return &FileTransferManager{
		client: client,
		logger: logger,
	}
}

// TransferFile performs a file transfer operation using the specified mode
func (ft *FileTransferManager) TransferFile(ctx context.Context, connection *spookytypes.Connection, transfer *spookytypesssh.FileTransfer) (*spookytypesssh.FileTransferResult, error) {
	startTime := time.Now()

	// Validate transfer request
	if err := ft.validateTransferRequest(transfer); err != nil {
		return ft.createFailedResult(transfer, err, startTime), nil
	}

	// Get connection from pool
	pooledConn, err := ft.client.connectionPool.GetConnection(connection.Host, connection.Port, connection.User)
	if err != nil {
		return ft.createFailedResult(transfer, fmt.Errorf("failed to get SSH connection: %w", err), startTime), nil
	}
	defer ft.client.connectionPool.ReturnConnection(pooledConn)

	// Perform transfer based on mode
	var result *spookytypesssh.FileTransferResult
	switch transfer.Mode {
	case spookytypesssh.TransferModeSFTP:
		result, err = ft.transferViaSFTP(ctx, pooledConn, transfer, startTime)
	case spookytypesssh.TransferModeSCP:
		result, err = ft.transferViaSCP(ctx, pooledConn, transfer, startTime)
	default:
		err = fmt.Errorf("unsupported transfer mode: %s", transfer.Mode)
		result = ft.createFailedResult(transfer, err, startTime)
	}

	if err != nil {
		ft.logger.Error("File transfer failed", err, map[string]interface{}{
			"local_path":  transfer.LocalPath,
			"remote_path": transfer.RemotePath,
			"mode":        transfer.Mode,
			"direction":   transfer.Direction,
		})
	}

	return result, err
}

// validateTransferRequest validates the file transfer request
func (ft *FileTransferManager) validateTransferRequest(transfer *spookytypesssh.FileTransfer) error {
	if transfer.LocalPath == "" {
		return fmt.Errorf("local path is required")
	}

	if transfer.RemotePath == "" {
		return fmt.Errorf("remote path is required")
	}

	// Validate local file exists for upload
	if transfer.Direction == spookytypesssh.TransferDirectionUpload {
		if _, err := os.Stat(transfer.LocalPath); os.IsNotExist(err) {
			return fmt.Errorf("local file does not exist: %s", transfer.LocalPath)
		}
	}

	return nil
}

// transferViaSFTP performs file transfer using SFTP
func (ft *FileTransferManager) transferViaSFTP(ctx context.Context, pooledConn *PooledConnection, transfer *spookytypesssh.FileTransfer, startTime time.Time) (*spookytypesssh.FileTransferResult, error) {
	// Create SFTP client
	sftpClient, err := sftp.NewClient(pooledConn.Client)
	if err != nil {
		return ft.createFailedResult(transfer, fmt.Errorf("failed to create SFTP client: %w", err), startTime), nil
	}
	defer sftpClient.Close()

	// Create progress tracker
	progress := ft.createProgressTracker(transfer)

	// Perform transfer based on direction
	var result *spookytypesssh.FileTransferResult
	switch transfer.Direction {
	case spookytypesssh.TransferDirectionUpload:
		result, err = ft.uploadViaSFTP(ctx, sftpClient, transfer, progress, startTime)
	case spookytypesssh.TransferDirectionDownload:
		result, err = ft.downloadViaSFTP(ctx, sftpClient, transfer, progress, startTime)
	default:
		err = fmt.Errorf("unsupported transfer direction: %s", transfer.Direction)
		result = ft.createFailedResult(transfer, err, startTime)
	}

	return result, err
}

// uploadViaSFTP uploads a file using SFTP
func (ft *FileTransferManager) uploadViaSFTP(ctx context.Context, sftpClient *sftp.Client, transfer *spookytypesssh.FileTransfer, progress *ProgressTracker, startTime time.Time) (*spookytypesssh.FileTransferResult, error) {
	// Open local file
	localFile, err := os.Open(transfer.LocalPath)
	if err != nil {
		return ft.createFailedResult(transfer, fmt.Errorf("failed to open local file: %w", err), startTime), nil
	}
	defer localFile.Close()

	// Get file info
	fileInfo, err := localFile.Stat()
	if err != nil {
		return ft.createFailedResult(transfer, fmt.Errorf("failed to get file info: %w", err), startTime), nil
	}

	// Create remote directory if needed
	remoteDir := filepath.Dir(transfer.RemotePath)
	if err := sftpClient.MkdirAll(remoteDir); err != nil {
		return ft.createFailedResult(transfer, fmt.Errorf("failed to create remote directory: %w", err), startTime), nil
	}

	// Create remote file
	remoteFile, err := sftpClient.Create(transfer.RemotePath)
	if err != nil {
		return ft.createFailedResult(transfer, fmt.Errorf("failed to create remote file: %w", err), startTime), nil
	}
	defer remoteFile.Close()

	// Set file permissions if specified
	if transfer.Permissions != 0 {
		if err := remoteFile.Chmod(transfer.Permissions); err != nil {
			ft.logger.Warn("Failed to set remote file permissions", map[string]interface{}{
				"remote_path": transfer.RemotePath,
				"permissions": transfer.Permissions,
				"error":       err.Error(),
			})
		}
	}

	// Copy file content with progress tracking
	bytesTransferred, err := ft.copyWithProgress(ctx, localFile, remoteFile, fileInfo.Size(), progress)
	if err != nil {
		return ft.createFailedResult(transfer, fmt.Errorf("failed to copy file content: %w", err), startTime), nil
	}

	// Verify transfer if requested
	if transfer.Verify {
		if err := ft.verifyTransfer(sftpClient, transfer, fileInfo.Size()); err != nil {
			return ft.createFailedResult(transfer, fmt.Errorf("transfer verification failed: %w", err), startTime), nil
		}
	}

	endTime := time.Now()
	duration := endTime.Sub(startTime)
	transferRate := float64(bytesTransferred) / duration.Seconds()

	ft.logger.Info("SFTP upload completed", map[string]interface{}{
		"local_path":        transfer.LocalPath,
		"remote_path":       transfer.RemotePath,
		"bytes_transferred": bytesTransferred,
		"duration":          duration,
		"transfer_rate":     transferRate,
	})

	return &spookytypesssh.FileTransferResult{
		Transfer:         transfer,
		Success:          true,
		BytesTransferred: bytesTransferred,
		StartTime:        startTime,
		EndTime:          endTime,
		Duration:         duration,
		TransferRate:     transferRate,
		LocalPath:        transfer.LocalPath,
		RemotePath:       transfer.RemotePath,
		Permissions:      transfer.Permissions,
	}, nil
}

// downloadViaSFTP downloads a file using SFTP
func (ft *FileTransferManager) downloadViaSFTP(ctx context.Context, sftpClient *sftp.Client, transfer *spookytypesssh.FileTransfer, progress *ProgressTracker, startTime time.Time) (*spookytypesssh.FileTransferResult, error) {
	// Open remote file
	remoteFile, err := sftpClient.Open(transfer.RemotePath)
	if err != nil {
		return ft.createFailedResult(transfer, fmt.Errorf("failed to open remote file: %w", err), startTime), nil
	}
	defer remoteFile.Close()

	// Get remote file info
	fileInfo, err := remoteFile.Stat()
	if err != nil {
		return ft.createFailedResult(transfer, fmt.Errorf("failed to get remote file info: %w", err), startTime), nil
	}

	// Create local directory if needed
	localDir := filepath.Dir(transfer.LocalPath)
	if err := os.MkdirAll(localDir, 0755); err != nil {
		return ft.createFailedResult(transfer, fmt.Errorf("failed to create local directory: %w", err), startTime), nil
	}

	// Create local file
	localFile, err := os.Create(transfer.LocalPath)
	if err != nil {
		return ft.createFailedResult(transfer, fmt.Errorf("failed to create local file: %w", err), startTime), nil
	}
	defer localFile.Close()

	// Set file permissions if specified
	if transfer.Permissions != 0 {
		if err := localFile.Chmod(transfer.Permissions); err != nil {
			ft.logger.Warn("Failed to set local file permissions", map[string]interface{}{
				"local_path":  transfer.LocalPath,
				"permissions": transfer.Permissions,
				"error":       err.Error(),
			})
		}
	}

	// Copy file content with progress tracking
	bytesTransferred, err := ft.copyWithProgress(ctx, remoteFile, localFile, fileInfo.Size(), progress)
	if err != nil {
		return ft.createFailedResult(transfer, fmt.Errorf("failed to copy file content: %w", err), startTime), nil
	}

	// Verify transfer if requested
	if transfer.Verify {
		if err := ft.verifyTransfer(sftpClient, transfer, fileInfo.Size()); err != nil {
			return ft.createFailedResult(transfer, fmt.Errorf("transfer verification failed: %w", err), startTime), nil
		}
	}

	endTime := time.Now()
	duration := endTime.Sub(startTime)
	transferRate := float64(bytesTransferred) / duration.Seconds()

	ft.logger.Info("SFTP download completed", map[string]interface{}{
		"local_path":        transfer.LocalPath,
		"remote_path":       transfer.RemotePath,
		"bytes_transferred": bytesTransferred,
		"duration":          duration,
		"transfer_rate":     transferRate,
	})

	return &spookytypesssh.FileTransferResult{
		Transfer:         transfer,
		Success:          true,
		BytesTransferred: bytesTransferred,
		StartTime:        startTime,
		EndTime:          endTime,
		Duration:         duration,
		TransferRate:     transferRate,
		LocalPath:        transfer.LocalPath,
		RemotePath:       transfer.RemotePath,
		Permissions:      transfer.Permissions,
	}, nil
}

// transferViaSCP performs file transfer using SCP
func (ft *FileTransferManager) transferViaSCP(ctx context.Context, pooledConn *PooledConnection, transfer *spookytypesssh.FileTransfer, startTime time.Time) (*spookytypesssh.FileTransferResult, error) {
	// Create SCP session
	session, err := pooledConn.Client.NewSession()
	if err != nil {
		return ft.createFailedResult(transfer, fmt.Errorf("failed to create SCP session: %w", err), startTime), nil
	}
	defer session.Close()

	// Perform transfer based on direction
	var result *spookytypesssh.FileTransferResult
	switch transfer.Direction {
	case spookytypesssh.TransferDirectionUpload:
		result, err = ft.uploadViaSCP(ctx, session, transfer, startTime)
	case spookytypesssh.TransferDirectionDownload:
		result, err = ft.downloadViaSCP(ctx, session, transfer, startTime)
	default:
		err = fmt.Errorf("unsupported transfer direction: %s", transfer.Direction)
		result = ft.createFailedResult(transfer, err, startTime)
	}

	return result, err
}

// uploadViaSCP uploads a file using SCP
func (ft *FileTransferManager) uploadViaSCP(ctx context.Context, session *ssh.Session, transfer *spookytypesssh.FileTransfer, startTime time.Time) (*spookytypesssh.FileTransferResult, error) {
	// Open local file
	localFile, err := os.Open(transfer.LocalPath)
	if err != nil {
		return ft.createFailedResult(transfer, fmt.Errorf("failed to open local file: %w", err), startTime), nil
	}
	defer localFile.Close()

	// Get file info
	fileInfo, err := localFile.Stat()
	if err != nil {
		return ft.createFailedResult(transfer, fmt.Errorf("failed to get file info: %w", err), startTime), nil
	}

	// Set up SCP command
	scpCommand := fmt.Sprintf("scp -t %s", transfer.RemotePath)
	if err := session.Start(scpCommand); err != nil {
		return ft.createFailedResult(transfer, fmt.Errorf("failed to start SCP command: %w", err), startTime), nil
	}

	// Get stdin pipe
	stdin, err := session.StdinPipe()
	if err != nil {
		return ft.createFailedResult(transfer, fmt.Errorf("failed to get stdin pipe: %w", err), startTime), nil
	}

	// Send file header
	header := fmt.Sprintf("C%04o %d %s\n", fileInfo.Mode().Perm(), fileInfo.Size(), filepath.Base(transfer.RemotePath))
	if _, err := stdin.Write([]byte(header)); err != nil {
		return ft.createFailedResult(transfer, fmt.Errorf("failed to send SCP header: %w", err), startTime), nil
	}

	// Send file content
	bytesTransferred, err := io.Copy(stdin, localFile)
	if err != nil {
		return ft.createFailedResult(transfer, fmt.Errorf("failed to send file content: %w", err), startTime), nil
	}

	// Send end marker
	if _, err := stdin.Write([]byte{0}); err != nil {
		return ft.createFailedResult(transfer, fmt.Errorf("failed to send SCP end marker: %w", err), startTime), nil
	}

	// Wait for completion
	if err := session.Wait(); err != nil {
		return ft.createFailedResult(transfer, fmt.Errorf("SCP command failed: %w", err), startTime), nil
	}

	endTime := time.Now()
	duration := endTime.Sub(startTime)
	transferRate := float64(bytesTransferred) / duration.Seconds()

	ft.logger.Info("SCP upload completed", map[string]interface{}{
		"local_path":        transfer.LocalPath,
		"remote_path":       transfer.RemotePath,
		"bytes_transferred": bytesTransferred,
		"duration":          duration,
		"transfer_rate":     transferRate,
	})

	return &spookytypesssh.FileTransferResult{
		Transfer:         transfer,
		Success:          true,
		BytesTransferred: bytesTransferred,
		StartTime:        startTime,
		EndTime:          endTime,
		Duration:         duration,
		TransferRate:     transferRate,
		LocalPath:        transfer.LocalPath,
		RemotePath:       transfer.RemotePath,
		Permissions:      fileInfo.Mode(),
	}, nil
}

// downloadViaSCP downloads a file using SCP
func (ft *FileTransferManager) downloadViaSCP(ctx context.Context, session *ssh.Session, transfer *spookytypesssh.FileTransfer, startTime time.Time) (*spookytypesssh.FileTransferResult, error) {
	// Create local file
	localFile, err := os.Create(transfer.LocalPath)
	if err != nil {
		return ft.createFailedResult(transfer, fmt.Errorf("failed to create local file: %w", err), startTime), nil
	}
	defer localFile.Close()

	// Set up SCP command
	scpCommand := fmt.Sprintf("scp -f %s", transfer.RemotePath)
	if err := session.Start(scpCommand); err != nil {
		return ft.createFailedResult(transfer, fmt.Errorf("failed to start SCP command: %w", err), startTime), nil
	}

	// Get stdout and stdin pipes
	stdout, err := session.StdoutPipe()
	if err != nil {
		return ft.createFailedResult(transfer, fmt.Errorf("failed to get stdout pipe: %w", err), startTime), nil
	}

	stdin, err := session.StdinPipe()
	if err != nil {
		return ft.createFailedResult(transfer, fmt.Errorf("failed to get stdin pipe: %w", err), startTime), nil
	}

	// Read file header
	header, err := bufio.NewReader(stdout).ReadString('\n')
	if err != nil {
		return ft.createFailedResult(transfer, fmt.Errorf("failed to read SCP header: %w", err), startTime), nil
	}

	// Parse header (format: C<permissions> <size> <filename>)
	if !strings.HasPrefix(header, "C") {
		return ft.createFailedResult(transfer, fmt.Errorf("invalid SCP header format"), startTime), nil
	}

	header = strings.TrimPrefix(header, "C")
	parts := strings.Fields(header)
	if len(parts) < 2 {
		return ft.createFailedResult(transfer, fmt.Errorf("invalid SCP header: insufficient parts"), startTime), nil
	}

	// Parse permissions and size
	permissions, err := strconv.ParseUint(parts[0], 8, 32)
	if err != nil {
		return ft.createFailedResult(transfer, fmt.Errorf("failed to parse permissions: %w", err), startTime), nil
	}

	_, err = strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return ft.createFailedResult(transfer, fmt.Errorf("failed to parse file size: %w", err), startTime), nil
	}

	// Send acknowledgment
	if _, err := stdin.Write([]byte{0}); err != nil {
		return ft.createFailedResult(transfer, fmt.Errorf("failed to send SCP acknowledgment: %w", err), startTime), nil
	}

	// Read file content
	bytesTransferred, err := io.Copy(localFile, stdout)
	if err != nil {
		return ft.createFailedResult(transfer, fmt.Errorf("failed to read file content: %w", err), startTime), nil
	}

	// Send completion acknowledgment
	if _, err := stdin.Write([]byte{0}); err != nil {
		return ft.createFailedResult(transfer, fmt.Errorf("failed to send completion acknowledgment: %w", err), startTime), nil
	}

	// Wait for completion
	if err := session.Wait(); err != nil {
		return ft.createFailedResult(transfer, fmt.Errorf("SCP command failed: %w", err), startTime), nil
	}

	// Set file permissions
	if err := localFile.Chmod(os.FileMode(permissions)); err != nil {
		ft.logger.Warn("Failed to set local file permissions", map[string]interface{}{
			"local_path":  transfer.LocalPath,
			"permissions": permissions,
			"error":       err.Error(),
		})
	}

	endTime := time.Now()
	duration := endTime.Sub(startTime)
	transferRate := float64(bytesTransferred) / duration.Seconds()

	ft.logger.Info("SCP download completed", map[string]interface{}{
		"local_path":        transfer.LocalPath,
		"remote_path":       transfer.RemotePath,
		"bytes_transferred": bytesTransferred,
		"duration":          duration,
		"transfer_rate":     transferRate,
	})

	return &spookytypesssh.FileTransferResult{
		Transfer:         transfer,
		Success:          true,
		BytesTransferred: bytesTransferred,
		StartTime:        startTime,
		EndTime:          endTime,
		Duration:         duration,
		TransferRate:     transferRate,
		LocalPath:        transfer.LocalPath,
		RemotePath:       transfer.RemotePath,
		Permissions:      os.FileMode(permissions),
	}, nil
}

// ProgressTracker tracks file transfer progress
type ProgressTracker struct {
	mu               sync.Mutex
	totalBytes       int64
	transferredBytes int64
	startTime        time.Time
	lastUpdate       time.Time
	logger           spookytypeslogging.Logger
}

// createProgressTracker creates a new progress tracker
func (ft *FileTransferManager) createProgressTracker(transfer *spookytypesssh.FileTransfer) *ProgressTracker {
	return &ProgressTracker{
		startTime: time.Now(),
		logger:    ft.logger,
	}
}

// Update updates the progress tracker
func (pt *ProgressTracker) Update(bytes int64) {
	pt.mu.Lock()
	defer pt.mu.Unlock()

	pt.transferredBytes += bytes
	now := time.Now()

	// Log progress every 5 seconds
	if now.Sub(pt.lastUpdate) >= 5*time.Second {
		percentage := float64(pt.transferredBytes) / float64(pt.totalBytes) * 100
		elapsed := now.Sub(pt.startTime)
		rate := float64(pt.transferredBytes) / elapsed.Seconds()

		pt.logger.Info("Transfer progress", map[string]interface{}{
			"transferred_bytes": pt.transferredBytes,
			"total_bytes":       pt.totalBytes,
			"percentage":        percentage,
			"transfer_rate":     rate,
			"elapsed_time":      elapsed,
		})

		pt.lastUpdate = now
	}
}

// copyWithProgress copies data with progress tracking
func (ft *FileTransferManager) copyWithProgress(ctx context.Context, src io.Reader, dst io.Writer, totalSize int64, progress *ProgressTracker) (int64, error) {
	progress.totalBytes = totalSize

	buffer := make([]byte, 32*1024) // 32KB buffer
	var totalBytes int64

	for {
		select {
		case <-ctx.Done():
			return totalBytes, ctx.Err()
		default:
		}

		n, err := src.Read(buffer)
		if n > 0 {
			_, writeErr := dst.Write(buffer[:n])
			if writeErr != nil {
				return totalBytes, writeErr
			}
			totalBytes += int64(n)
			progress.Update(int64(n))
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return totalBytes, err
		}
	}

	return totalBytes, nil
}

// verifyTransfer verifies the transferred file
func (ft *FileTransferManager) verifyTransfer(sftpClient *sftp.Client, transfer *spookytypesssh.FileTransfer, expectedSize int64) error {
	// Get remote file info
	fileInfo, err := sftpClient.Stat(transfer.RemotePath)
	if err != nil {
		return fmt.Errorf("failed to get remote file info: %w", err)
	}

	// Check file size
	if fileInfo.Size() != expectedSize {
		return fmt.Errorf("file size mismatch: expected %d, got %d", expectedSize, fileInfo.Size())
	}

	// Calculate checksum if requested
	if transfer.Checksum != "" {
		remoteChecksum, err := ft.calculateRemoteChecksum(sftpClient, transfer.RemotePath)
		if err != nil {
			return fmt.Errorf("failed to calculate remote checksum: %w", err)
		}

		if remoteChecksum != transfer.Checksum {
			return fmt.Errorf("checksum mismatch: expected %s, got %s", transfer.Checksum, remoteChecksum)
		}
	}

	return nil
}

// calculateRemoteChecksum calculates the checksum of a remote file
func (ft *FileTransferManager) calculateRemoteChecksum(sftpClient *sftp.Client, remotePath string) (string, error) {
	// Open remote file
	file, err := sftpClient.Open(remotePath)
	if err != nil {
		return "", fmt.Errorf("failed to open remote file: %w", err)
	}
	defer file.Close()

	// Calculate SHA256 hash
	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", fmt.Errorf("failed to calculate hash: %w", err)
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

// createFailedResult creates a failed transfer result
func (ft *FileTransferManager) createFailedResult(transfer *spookytypesssh.FileTransfer, err error, startTime time.Time) *spookytypesssh.FileTransferResult {
	endTime := time.Now()
	return &spookytypesssh.FileTransferResult{
		Transfer:   transfer,
		Success:    false,
		Error:      err.Error(),
		StartTime:  startTime,
		EndTime:    endTime,
		Duration:   endTime.Sub(startTime),
		LocalPath:  transfer.LocalPath,
		RemotePath: transfer.RemotePath,
	}
}

// BatchTransfer performs multiple file transfers
func (ft *FileTransferManager) BatchTransfer(ctx context.Context, connection *spookytypes.Connection, transfers []*spookytypesssh.FileTransfer) ([]*spookytypesssh.FileTransferResult, error) {
	var results []*spookytypesssh.FileTransferResult
	var mu sync.Mutex
	var wg sync.WaitGroup

	// Create semaphore for concurrency control
	semaphore := make(chan struct{}, 5) // Limit to 5 concurrent transfers

	for _, transfer := range transfers {
		wg.Add(1)
		go func(t *spookytypesssh.FileTransfer) {
			defer wg.Done()

			// Acquire semaphore
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			// Perform transfer
			result, err := ft.TransferFile(ctx, connection, t)
			if err != nil {
				ft.logger.Error("Batch transfer failed", err, map[string]interface{}{
					"local_path":  t.LocalPath,
					"remote_path": t.RemotePath,
				})
			}

			// Add result to slice
			mu.Lock()
			results = append(results, result)
			mu.Unlock()
		}(transfer)
	}

	wg.Wait()

	return results, nil
}
