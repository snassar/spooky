package logging

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

// ExampleUsage demonstrates how to use the XDG state directory-based logging system
func ExampleUsage() {
	// Example 1: Basic database logging with XDG state directory
	config := &LogConfig{
		Level:  "info",
		Format: "json",
		Output: "stderr", // Also log to stderr for immediate feedback
		Database: &LogDatabaseConfig{
			Enabled: true,
			// Path will default to ~/.local/state/spooky/logs.db on Linux/Unix
			// or %LOCALAPPDATA%\spooky\State\logs.db on Windows
		},
	}

	logger, err := NewLogger(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create logger: %v\n", err)
		return
	}
	defer func() {
		if closeErr := logger.Close(); closeErr != nil {
			// Log the error but don't fail the example since this is cleanup
			fmt.Fprintf(os.Stderr, "Warning: failed to close logger: %v\n", closeErr)
		}
	}()

	// Example 2: Multiple concurrent action runs
	runID1 := uuid.New().String()
	runID2 := uuid.New().String()

	// Simulate multiple Spooky processes running simultaneously
	go simulateActionRun(logger, runID1, "deploy_webserver", "webserver-project")
	go simulateActionRun(logger, runID2, "backup_database", "database-project")

	// Wait for both to complete
	time.Sleep(2 * time.Second)

	// Example 3: Query logs by run ID
	entries, err := logger.QueryLogs(LogQuery{
		RunID: runID1,
		Limit: 100,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to query logs: %v\n", err)
		return
	}

	fmt.Printf("Found %d log entries for run %s\n", len(entries), runID1)
	for _, entry := range entries {
		fmt.Printf("[%s] %s: %s\n", entry.Level, entry.ActionName, entry.Message)
	}

	// Example 4: Query all error logs from the last hour
	errorEntries, err := logger.QueryLogs(LogQuery{
		Level:     "error",
		StartTime: time.Now().Add(-1 * time.Hour),
		Limit:     50,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to query error logs: %v\n", err)
		return
	}

	fmt.Printf("Found %d error entries in the last hour\n", len(errorEntries))
}

// simulateActionRun simulates a Spooky action run with structured logging
func simulateActionRun(logger *Logger, runID, actionName, projectName string) {
	// Log action start
	startTime := time.Now()
	if err := logger.LogToDatabase(LogEntry{
		Timestamp:   startTime,
		Level:       "info",
		Message:     "Starting action execution",
		RunID:       runID,
		ActionName:  actionName,
		ProjectName: projectName,
		Host:        "localhost",
		Tags:        []string{"action_start", "automation"},
	}); err != nil {
		fmt.Printf("Warning: failed to log action start: %v\n", err)
	}

	// Simulate some tasks
	tasks := []string{"validate_config", "deploy_packages", "start_services", "verify_deployment"}
	for i, task := range tasks {
		taskStart := time.Now()

		// Simulate task execution
		time.Sleep(100 * time.Millisecond)

		// Log task completion
		if err := logger.LogToDatabase(LogEntry{
			Timestamp:   time.Now(),
			Level:       "info",
			Message:     fmt.Sprintf("Task %d completed successfully", i+1),
			RunID:       runID,
			ActionName:  actionName,
			ProjectName: projectName,
			TaskID:      task,
			Host:        "localhost",
			Duration:    time.Since(taskStart),
			Tags:        []string{"task_complete", task},
		}); err != nil {
			fmt.Printf("Warning: failed to log task completion: %v\n", err)
		}
	}

	// Log action completion
	duration := time.Since(startTime)
	if err := logger.LogToDatabase(LogEntry{
		Timestamp:   time.Now(),
		Level:       "info",
		Message:     "Action completed successfully",
		RunID:       runID,
		ActionName:  actionName,
		ProjectName: projectName,
		Host:        "localhost",
		Duration:    duration,
		Tags:        []string{"action_complete", "success"},
	}); err != nil {
		fmt.Printf("Warning: failed to log action completion: %v\n", err)
	}
}

// ExampleXDGPaths demonstrates the XDG state directory structure
func ExampleXDGPaths() {
	// On Linux/Unix systems:
	// XDG_STATE_HOME defaults to ~/.local/state
	// Spooky logs will be stored in ~/.local/state/spooky/logs.db

	// On Windows:
	// Equivalent to %LOCALAPPDATA%\spooky\State\logs.db

	homeDir, _ := os.UserHomeDir()

	// Linux/Unix path
	linuxPath := filepath.Join(homeDir, ".local", "state", "spooky", "logs.db")
	fmt.Printf("Linux/Unix SQLite path: %s\n", linuxPath)

	// Windows path
	windowsPath := filepath.Join(homeDir, "AppData", "Local", "spooky", "State", "logs.db")
	fmt.Printf("Windows SQLite path: %s\n", windowsPath)

	// Custom XDG_STATE_HOME
	if err := os.Setenv("XDG_STATE_HOME", "/custom/state/path"); err != nil {
		fmt.Printf("Warning: failed to set XDG_STATE_HOME environment variable: %v\n", err)
	}
	customPath := filepath.Join("/custom/state/path", "spooky", "logs.db")
	fmt.Printf("Custom XDG_STATE_HOME path: %s\n", customPath)
}

// ExampleMultiProcess demonstrates how multiple Spooky processes can write simultaneously
func ExampleMultiProcess() {
	// This example shows how multiple Spooky processes can safely write to the same database

	config := &LogConfig{
		Database: &LogDatabaseConfig{
			Enabled: true,
			// All processes will use the same database file
			// SQLite WAL mode handles concurrent writes safely
		},
	}

	// Process 1: Deploy webserver
	go func() {
		logger, _ := NewLogger(config)
		defer func() {
			if closeErr := logger.Close(); closeErr != nil {
				// Log the error but don't fail the example since this is cleanup
				fmt.Printf("Warning: failed to close logger: %v\n", closeErr)
			}
		}()

		runID := uuid.New().String()
		if err := logger.LogToDatabase(LogEntry{
			Timestamp:   time.Now(),
			Level:       "info",
			Message:     "Deploying webserver",
			RunID:       runID,
			ActionName:  "deploy_webserver",
			ProjectName: "webserver-project",
		}); err != nil {
			fmt.Printf("Warning: failed to log webserver deployment: %v\n", err)
		}
	}()

	// Process 2: Backup database
	go func() {
		logger, _ := NewLogger(config)
		defer func() {
			if closeErr := logger.Close(); closeErr != nil {
				// Log the error but don't fail the example since this is cleanup
				fmt.Printf("Warning: failed to close logger: %v\n", closeErr)
			}
		}()

		runID := uuid.New().String()
		if err := logger.LogToDatabase(LogEntry{
			Timestamp:   time.Now(),
			Level:       "info",
			Message:     "Backing up database",
			RunID:       runID,
			ActionName:  "backup_database",
			ProjectName: "database-project",
		}); err != nil {
			fmt.Printf("Warning: failed to log database backup: %v\n", err)
		}
	}()

	// Process 3: Update monitoring
	go func() {
		logger, _ := NewLogger(config)
		defer func() {
			if closeErr := logger.Close(); closeErr != nil {
				// Log the error but don't fail the example since this is cleanup
				fmt.Printf("Warning: failed to close logger: %v\n", closeErr)
			}
		}()

		runID := uuid.New().String()
		if err := logger.LogToDatabase(LogEntry{
			Timestamp:   time.Now(),
			Level:       "info",
			Message:     "Updating monitoring",
			RunID:       runID,
			ActionName:  "update_monitoring",
			ProjectName: "monitoring-project",
		}); err != nil {
			fmt.Printf("Warning: failed to log monitoring update: %v\n", err)
		}
	}()

	// All three processes can write to the same SQLite database simultaneously
	// without conflicts thanks to WAL mode
}
