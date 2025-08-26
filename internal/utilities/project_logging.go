package utilities

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"spooky/internal/logging"

	"github.com/pkg/errors"
)

// ProjectLogger manages project-specific logging with timestamps
type ProjectLogger struct {
	projectName string
	runID       string
	timestamp   time.Time
	config      *PathConfig
	logger      *logging.Logger
	logFile     string
}

// NewProjectLogger creates a new project logger with timestamped run
func NewProjectLogger(projectName string) (*ProjectLogger, error) {
	// Get OS-specific path configuration
	config, err := GetPathConfig("spooky")
	if err != nil {
		return nil, errors.Wrap(err, "failed to get path configuration")
	}

	// Ensure directories exist
	if err := EnsureDirectories(config); err != nil {
		return nil, errors.Wrap(err, "failed to ensure directories")
	}

	// Generate run ID and timestamp
	runID := generateRunID()
	timestamp := time.Now()

	// Create project-specific log directory
	projectLogDir := filepath.Join(config.LogDir, "projects", projectName)
	if err := os.MkdirAll(projectLogDir, 0o755); err != nil {
		return nil, errors.Wrapf(err, "failed to create project log directory: %s", projectLogDir)
	}

	// Create timestamped log file
	logFileName := fmt.Sprintf("%s_%s.log", runID, timestamp.Format("20060102-150405"))
	logFile := filepath.Join(projectLogDir, logFileName)

	// Create logging configuration for this project
	logConfig := &logging.LogConfig{
		Level:  "info",
		Format: "json", // Use JSON for better parsing
		Output: "file",
		File: &logging.LogFileConfig{
			Path: logFile,
		},
		Structured: &logging.StructuredConfig{
			IncludeTimestamp: true,
			IncludeLevel:     true,
			CustomFields: map[string]string{
				"project":   projectName,
				"run_id":    runID,
				"timestamp": timestamp.Format(time.RFC3339),
			},
		},
	}

	// Create logger
	logger, err := logging.NewLogger(logConfig)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create project logger")
	}

	return &ProjectLogger{
		projectName: projectName,
		runID:       runID,
		timestamp:   timestamp,
		config:      config,
		logger:      logger,
		logFile:     logFile,
	}, nil
}

// generateRunID creates a unique run identifier
func generateRunID() string {
	return fmt.Sprintf("run_%s", time.Now().Format("20060102_150405_000"))
}

// GetLogFile returns the path to the current log file
func (pl *ProjectLogger) GetLogFile() string {
	return pl.logFile
}

// GetRunID returns the current run ID
func (pl *ProjectLogger) GetRunID() string {
	return pl.runID
}

// GetTimestamp returns the run timestamp
func (pl *ProjectLogger) GetTimestamp() time.Time {
	return pl.timestamp
}

// Logger returns the underlying logger
func (pl *ProjectLogger) Logger() *logging.Logger {
	return pl.logger
}

// Close closes the project logger
func (pl *ProjectLogger) Close() error {
	if pl.logger != nil {
		return pl.logger.Close()
	}
	return nil
}

// LogProjectStart logs the start of a project run
func (pl *ProjectLogger) LogProjectStart(projectPath string, args []string) {
	pl.logger.Info("project run started",
		slog.String("project_name", pl.projectName),
		slog.String("project_path", projectPath),
		slog.String("run_id", pl.runID),
		slog.String("log_file", pl.logFile),
		slog.Any("arguments", args))
}

// LogProjectEnd logs the end of a project run
func (pl *ProjectLogger) LogProjectEnd(duration time.Duration, success bool) {
	if success {
		pl.logger.Info("project run completed",
			slog.String("project_name", pl.projectName),
			slog.String("run_id", pl.runID),
			slog.Duration("duration", duration),
			slog.Bool("success", success))
	} else {
		pl.logger.Error("project run completed",
			slog.String("project_name", pl.projectName),
			slog.String("run_id", pl.runID),
			slog.Duration("duration", duration),
			slog.Bool("success", success))
	}
}

// LogTask logs a task execution
func (pl *ProjectLogger) LogTask(taskName, taskType string, duration time.Duration, success bool, details map[string]interface{}) {
	logger := pl.logger.With(
		slog.String("task_name", taskName),
		slog.String("task_type", taskType),
		slog.Duration("duration", duration),
		slog.Bool("success", success),
	)

	// Add custom details
	for k, v := range details {
		logger = logger.With(slog.Any(k, v))
	}

	if success {
		logger.Info("task executed")
	} else {
		logger.Error("task executed")
	}
}

// LogError logs an error with context
func (pl *ProjectLogger) LogError(err error, context string, details map[string]interface{}) {
	logger := pl.logger.WithError(err).With(
		slog.String("context", context),
	)

	// Add custom details
	for k, v := range details {
		logger = logger.With(slog.Any(k, v))
	}

	logger.Error("error occurred")
}

// ListProjectLogs lists all log files for a project
func ListProjectLogs(projectName string) ([]string, error) {
	config, err := GetPathConfig("spooky")
	if err != nil {
		return nil, errors.Wrap(err, "failed to get path configuration")
	}

	projectLogDir := filepath.Join(config.LogDir, "projects", projectName)

	// Check if directory exists
	if _, err := os.Stat(projectLogDir); os.IsNotExist(err) {
		return []string{}, nil // No logs yet
	}

	// Read directory
	entries, err := os.ReadDir(projectLogDir)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read project log directory: %s", projectLogDir)
	}

	var logFiles []string
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".log" {
			logFiles = append(logFiles, filepath.Join(projectLogDir, entry.Name()))
		}
	}

	return logFiles, nil
}

// GetProjectLogSummary returns a summary of project logs
func GetProjectLogSummary(projectName string) (*ProjectLogSummary, error) {
	logFiles, err := ListProjectLogs(projectName)
	if err != nil {
		return nil, err
	}

	summary := &ProjectLogSummary{
		ProjectName: projectName,
		TotalRuns:   len(logFiles),
		LogFiles:    logFiles,
	}

	// Get latest log file
	if len(logFiles) > 0 {
		summary.LatestLogFile = logFiles[len(logFiles)-1]
	}

	return summary, nil
}

// ProjectLogSummary represents a summary of project logs
type ProjectLogSummary struct {
	ProjectName   string   `json:"project_name"`
	TotalRuns     int      `json:"total_runs"`
	LogFiles      []string `json:"log_files"`
	LatestLogFile string   `json:"latest_log_file,omitempty"`
}
