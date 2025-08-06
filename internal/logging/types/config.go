package types

// Config holds logging configuration
type Config struct {
	Level     LogLevel `env:"LOG_LEVEL" default:"info"`
	Format    string   `env:"LOG_FORMAT" default:"json"`
	Output    string   `env:"LOG_OUTPUT" default:"stdout"`
	Timestamp bool     `env:"LOG_TIMESTAMP" default:"true"`
}

// LogLevel represents the logging level
type LogLevel string

const (
	DebugLevel LogLevel = "debug"
	InfoLevel  LogLevel = "info"
	WarnLevel  LogLevel = "warn"
	ErrorLevel LogLevel = "error"
)
