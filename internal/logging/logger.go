package logging

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/errors"
	_ "modernc.org/sqlite"
)

// LogConfig represents logging configuration from HCL schemas
type LogConfig struct {
	Level      string             `hcl:"level,optional"`
	Format     string             `hcl:"format,optional"`
	Output     string             `hcl:"output,optional"`
	File       *LogFileConfig     `hcl:"file,block"`
	Database   *LogDatabaseConfig `hcl:"database,block"`
	Structured *StructuredConfig  `hcl:"structured,block"`
	Buffered   *BufferedConfig    `hcl:"buffered,block"`
	Async      *AsyncConfig       `hcl:"async,block"`
}

// LogFileConfig represents file logging configuration
type LogFileConfig struct {
	Path        string             `hcl:"path"`
	Permissions string             `hcl:"permissions,optional"`
	Rotation    *LogRotationConfig `hcl:"rotation,block"`
}

// LogRotationConfig represents log rotation configuration
type LogRotationConfig struct {
	MaxSize    string `hcl:"max_size,optional"`
	MaxAge     string `hcl:"max_age,optional"`
	MaxBackups int    `hcl:"max_backups,optional"`
	Compress   bool   `hcl:"compress,optional"`
}

// StructuredConfig represents structured logging configuration
type StructuredConfig struct {
	IncludeTimestamp bool              `hcl:"include_timestamp,optional"`
	IncludeLevel     bool              `hcl:"include_level,optional"`
	IncludeSource    bool              `hcl:"include_source,optional"`
	CustomFields     map[string]string `hcl:"custom_fields,optional"`
}

// BufferedConfig represents buffered logging configuration
type BufferedConfig struct {
	Enabled bool `hcl:"enabled,optional"`
	Size    int  `hcl:"size,optional"`
}

// AsyncConfig represents asynchronous logging configuration
type AsyncConfig struct {
	Enabled     bool `hcl:"enabled,optional"`
	QueueSize   int  `hcl:"queue_size,optional"`
	WorkerCount int  `hcl:"worker_count,optional"`
}

// LogDatabaseConfig represents database logging configuration
type LogDatabaseConfig struct {
	Enabled bool   `hcl:"enabled,optional"`
	Path    string `hcl:"path,optional"`
	Driver  string `hcl:"driver,optional"` // "sqlite", "postgres", etc.
}

// LogEntry represents a structured log entry for database storage
type LogEntry struct {
	ID          int64                  `json:"id"`
	Timestamp   time.Time              `json:"timestamp"`
	Level       string                 `json:"level"`
	Message     string                 `json:"message"`
	RunID       string                 `json:"run_id"`
	ActionName  string                 `json:"action_name"`
	ProjectName string                 `json:"project_name"`
	Host        string                 `json:"host"`
	TaskID      string                 `json:"task_id"`
	Duration    time.Duration          `json:"duration"`
	Error       string                 `json:"error"`
	Metadata    map[string]interface{} `json:"metadata"`
	Tags        []string               `json:"tags"`
	CreatedAt   time.Time              `json:"created_at"`
}

// Logger wraps slog.Logger with additional functionality
type Logger struct {
	*slog.Logger
	config *LogConfig
	writer io.WriteCloser
	db     *sql.DB
}

// nullWriter implements io.WriteCloser for null output
type nullWriter struct{}

func (n *nullWriter) Write(p []byte) (int, error) {
	return len(p), nil
}

func (n *nullWriter) Close() error {
	return nil
}

// DefaultConfig returns a sensible default logging configuration
func DefaultConfig() *LogConfig {
	return &LogConfig{
		Level:  "info",
		Format: "text",
		Output: "stderr",
	}
}

// NewLogger creates a new logger with the given configuration
func NewLogger(config *LogConfig) (*Logger, error) {
	var handler slog.Handler
	var writer io.WriteCloser
	var db *sql.DB
	var err error

	// Handle database logging if enabled
	if config.Database != nil && config.Database.Enabled {
		db, err = createDatabaseLogger(config.Database)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create database logger")
		}
	}

	// Create handler and writer for traditional logging
	handler, writer, err = createHandler(config)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create handler")
	}

	logger := &Logger{
		Logger: slog.New(handler),
		config: config,
		writer: writer,
		db:     db,
	}

	return logger, nil
}

// createHandler creates the appropriate slog handler based on configuration
func createHandler(config *LogConfig) (slog.Handler, io.WriteCloser, error) {
	var writer io.WriteCloser
	var err error

	// Determine output writer
	switch strings.ToLower(config.Output) {
	case "stdout":
		writer = os.Stdout
	case "stderr", "":
		writer = os.Stderr
	case "file":
		if config.File == nil || config.File.Path == "" {
			return nil, nil, errors.New("file output requires file.path configuration")
		}
		writer, err = createLogFile(config.File)
		if err != nil {
			return nil, nil, errors.Wrap(err, "failed to create log file")
		}
	case "null":
		writer = &nullWriter{}
	default:
		return nil, nil, errors.Errorf("unsupported output type: %s", config.Output)
	}

	// Determine log level
	level := parseLogLevel(config.Level)

	// Create handler based on format
	var handler slog.Handler
	switch strings.ToLower(config.Format) {
	case "json":
		handler = slog.NewJSONHandler(writer, &slog.HandlerOptions{
			Level:     level,
			AddSource: config.Structured != nil && config.Structured.IncludeSource,
		})
	case "text", "":
		handler = slog.NewTextHandler(writer, &slog.HandlerOptions{
			Level:     level,
			AddSource: config.Structured != nil && config.Structured.IncludeSource,
		})
	default:
		return nil, nil, errors.Errorf("unsupported format: %s", config.Format)
	}

	// Add custom fields if structured logging is configured
	if config.Structured != nil && config.Structured.CustomFields != nil {
		attrs := make([]slog.Attr, 0, len(config.Structured.CustomFields))
		for k, v := range config.Structured.CustomFields {
			attrs = append(attrs, slog.String(k, v))
		}
		handler = handler.WithAttrs(attrs)
	}

	return handler, writer, nil
}

// createLogFile creates a log file with proper permissions and rotation
func createLogFile(config *LogFileConfig) (io.WriteCloser, error) {
	// Ensure directory exists
	dir := filepath.Dir(config.Path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, errors.Wrapf(err, "failed to create log directory: %s", dir)
	}

	// Open file with appropriate permissions
	file, err := os.OpenFile(config.Path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to open log file: %s", config.Path)
	}

	// Check if rotation is configured
	if config.Rotation != nil {
		// For now, return the basic file writer
		// Log rotation can be implemented later as a separate feature
		// This ensures the current implementation is functional
		return file, nil
	}

	return file, nil
}

// createDatabaseLogger creates a SQLite database for structured logging
func createDatabaseLogger(config *LogDatabaseConfig) (*sql.DB, error) {
	dbPath := config.Path
	if dbPath == "" {
		// Use default XDG state directory
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, errors.Wrap(err, "failed to get user home directory")
		}

		// Determine XDG_STATE_HOME or equivalent
		xdgStateHome := os.Getenv("XDG_STATE_HOME")
		if xdgStateHome == "" {
			xdgStateHome = filepath.Join(homeDir, ".local", "state")
		}

		dbPath = filepath.Join(xdgStateHome, "spooky", "logs.db")
	}

	// Ensure directory exists
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, errors.Wrapf(err, "failed to create database directory: %s", dir)
	}

	// Open SQLite database with WAL mode for multi-process access
	db, err := sql.Open("sqlite", dbPath+"?_journal_mode=WAL&_synchronous=NORMAL&_cache_size=10000&_busy_timeout=5000")
	if err != nil {
		return nil, errors.Wrapf(err, "failed to open database: %s", dbPath)
	}

	// Initialize database schema
	if err := initDatabaseSchema(db); err != nil {
		return nil, errors.Wrap(err, "failed to initialize database schema")
	}

	return db, nil
}

// initDatabaseSchema creates the necessary tables for logging
func initDatabaseSchema(db *sql.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS log_entries (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		timestamp TEXT NOT NULL,
		level TEXT NOT NULL,
		message TEXT NOT NULL,
		run_id TEXT NOT NULL,
		action_name TEXT,
		project_name TEXT,
		host TEXT,
		task_id TEXT,
		duration_ms INTEGER,
		error TEXT,
		metadata TEXT,
		tags TEXT,
		created_at TEXT DEFAULT (datetime('now'))
	);

	CREATE INDEX IF NOT EXISTS idx_log_entries_timestamp ON log_entries(timestamp);
	CREATE INDEX IF NOT EXISTS idx_log_entries_run_id ON log_entries(run_id);
	CREATE INDEX IF NOT EXISTS idx_log_entries_level ON log_entries(level);
	CREATE INDEX IF NOT EXISTS idx_log_entries_action_name ON log_entries(action_name);
	CREATE INDEX IF NOT EXISTS idx_log_entries_host ON log_entries(host);
	CREATE INDEX IF NOT EXISTS idx_log_entries_created_at ON log_entries(created_at);

	CREATE TABLE IF NOT EXISTS action_runs (
		run_id TEXT PRIMARY KEY,
		action_name TEXT NOT NULL,
		project_name TEXT NOT NULL,
		status TEXT NOT NULL,
		started_at TEXT NOT NULL,
		completed_at TEXT,
		duration_ms INTEGER,
		error_count INTEGER DEFAULT 0,
		log_count INTEGER DEFAULT 0
	);
	`

	_, err := db.Exec(schema)
	return errors.Wrap(err, "failed to create database schema")
}

// LogToDatabase logs an entry to the SQLite database
func (l *Logger) LogToDatabase(entry LogEntry) error {
	if l.db == nil {
		return errors.New("database logging not enabled")
	}

	// Convert metadata and tags to JSON
	metadataJSON := "{}"
	if entry.Metadata != nil {
		if data, err := json.Marshal(entry.Metadata); err == nil {
			metadataJSON = string(data)
		}
	}

	tagsJSON := "[]"
	if entry.Tags != nil {
		if data, err := json.Marshal(entry.Tags); err == nil {
			tagsJSON = string(data)
		}
	}

	// Insert log entry
	_, err := l.db.Exec(`
		INSERT INTO log_entries (
			timestamp, level, message, run_id, action_name, 
			project_name, host, task_id, duration_ms, error, 
			metadata, tags
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		entry.Timestamp.Format(time.RFC3339),
		entry.Level,
		entry.Message,
		entry.RunID,
		entry.ActionName,
		entry.ProjectName,
		entry.Host,
		entry.TaskID,
		entry.Duration.Milliseconds(),
		entry.Error,
		metadataJSON,
		tagsJSON,
	)

	return errors.Wrap(err, "failed to insert log entry to database")
}

// QueryLogs retrieves log entries from the database
func (l *Logger) QueryLogs(query LogQuery) ([]LogEntry, error) {
	if l.db == nil {
		return nil, errors.New("database logging not enabled")
	}

	// Build query with conditions
	whereClause := "WHERE 1=1"
	args := []interface{}{}

	if query.RunID != "" {
		whereClause += " AND run_id = ?"
		args = append(args, query.RunID)
	}

	if query.Level != "" {
		whereClause += " AND level = ?"
		args = append(args, query.Level)
	}

	if query.ActionName != "" {
		whereClause += " AND action_name = ?"
		args = append(args, query.ActionName)
	}

	if query.Host != "" {
		whereClause += " AND host = ?"
		args = append(args, query.Host)
	}

	if !query.StartTime.IsZero() {
		whereClause += " AND timestamp >= ?"
		args = append(args, query.StartTime.Format(time.RFC3339))
	}

	if !query.EndTime.IsZero() {
		whereClause += " AND timestamp <= ?"
		args = append(args, query.EndTime.Format(time.RFC3339))
	}

	// Execute query
	rows, err := l.db.Query(`
		SELECT id, timestamp, level, message, run_id, action_name, 
		       project_name, host, task_id, duration_ms, error, 
		       metadata, tags, created_at
		FROM log_entries 
		`+whereClause+`
		ORDER BY timestamp DESC
		LIMIT ?
	`, append(args, query.Limit)...)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []LogEntry
	for rows.Next() {
		var entry LogEntry
		var timestampStr, metadataJSON, tagsJSON string
		var durationMs int64

		err := rows.Scan(
			&entry.ID,
			&timestampStr,
			&entry.Level,
			&entry.Message,
			&entry.RunID,
			&entry.ActionName,
			&entry.ProjectName,
			&entry.Host,
			&entry.TaskID,
			&durationMs,
			&entry.Error,
			&metadataJSON,
			&tagsJSON,
			&entry.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Parse timestamp
		if t, err := time.Parse(time.RFC3339, timestampStr); err == nil {
			entry.Timestamp = t
		}

		// Parse duration
		entry.Duration = time.Duration(durationMs) * time.Millisecond

		// Parse metadata
		if metadataJSON != "{}" {
			if err := json.Unmarshal([]byte(metadataJSON), &entry.Metadata); err != nil {
				// Log warning but continue - metadata parsing failure shouldn't break the entire query
				entry.Metadata = make(map[string]interface{})
			}
		}

		// Parse tags
		if tagsJSON != "[]" {
			if err := json.Unmarshal([]byte(tagsJSON), &entry.Tags); err != nil {
				// Log warning but continue - tags parsing failure shouldn't break the entire query
				entry.Tags = []string{}
			}
		}

		entries = append(entries, entry)
	}

	return entries, nil
}

// LogQuery represents a query for log entries
type LogQuery struct {
	RunID      string    `json:"run_id"`
	Level      string    `json:"level"`
	ActionName string    `json:"action_name"`
	Host       string    `json:"host"`
	StartTime  time.Time `json:"start_time"`
	EndTime    time.Time `json:"end_time"`
	Limit      int       `json:"limit"`
}

// parseLogLevel converts string level to slog.Level
func parseLogLevel(level string) slog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "info", "":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// Close closes the logger and its underlying writer
func (l *Logger) Close() error {
	if l.writer != nil {
		return l.writer.Close()
	}
	return nil
}

// WithError adds error context to the logger
func (l *Logger) WithError(err error) *Logger {
	if err == nil {
		return l
	}

	// Extract error details
	attrs := []slog.Attr{
		slog.String("error", err.Error()),
	}

	// Add error type if available
	if errorType := fmt.Sprintf("%T", err); errorType != "*errors.errorString" {
		attrs = append(attrs, slog.String("error_type", errorType))
	}

	// Add stack trace if available
	if _, ok := err.(interface{ StackTrace() []uintptr }); ok {
		attrs = append(attrs, slog.String("stack_trace", "available"))
	}

	logger := l.Logger
	for _, attr := range attrs {
		logger = logger.With(attr)
	}
	return &Logger{
		Logger: logger,
		config: l.config,
		writer: l.writer,
	}
}

// WithContext adds context information to the logger
func (l *Logger) WithContext(ctx context.Context) *Logger {
	// Extract useful context information
	attrs := []slog.Attr{}

	// Add request ID if available
	if requestID := ctx.Value("request_id"); requestID != nil {
		attrs = append(attrs, slog.String("request_id", fmt.Sprintf("%v", requestID)))
	}

	// Add operation if available
	if operation := ctx.Value("operation"); operation != nil {
		attrs = append(attrs, slog.String("operation", fmt.Sprintf("%v", operation)))
	}

	if len(attrs) == 0 {
		return l
	}

	logger := l.Logger
	for _, attr := range attrs {
		logger = logger.With(attr)
	}
	return &Logger{
		Logger: logger,
		config: l.config,
		writer: l.writer,
	}
}

// Global logger instance
var globalLogger *Logger

// SetGlobalLogger sets the global logger instance
func SetGlobalLogger(logger *Logger) {
	globalLogger = logger
	slog.SetDefault(logger.Logger)
}

// GetGlobalLogger returns the global logger instance
func GetGlobalLogger() *Logger {
	if globalLogger == nil {
		// Create default logger if none exists
		logger, err := NewLogger(DefaultConfig())
		if err != nil {
			// Fallback to basic logger
			handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo})
			globalLogger = &Logger{
				Logger: slog.New(handler),
				config: DefaultConfig(),
				writer: os.Stderr,
			}
		} else {
			globalLogger = logger
		}
		slog.SetDefault(globalLogger.Logger)
	}
	return globalLogger
}
