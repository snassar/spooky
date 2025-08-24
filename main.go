package main

import (
	"log/slog"
	"os"
	"spooky/commands"
	"spooky/internal/logging"
)

func main() {
	// Check for logging flags before command execution
	logLevel := "warn" // default to warn level to avoid INFO noise

	// Simple flag checking for logging control
	for i, arg := range os.Args[1:] {
		switch arg {
		case "--quiet", "-q":
			logLevel = "error"
		case "--verbose", "-v":
			logLevel = "debug"
		case "--log-level":
			if i+1 < len(os.Args[1:]) {
				logLevel = os.Args[i+2]
			}
		}
	}

	// Create logging configuration based on flags
	logConfig := &logging.LogConfig{
		Level:  logLevel,
		Format: "text", // Use text format for better CLI readability
		Output: "stderr",
	}

	// Initialize global logger
	logger, err := logging.NewLogger(logConfig)
	if err != nil {
		// Fallback to basic logging
		slog.Error("failed to initialize logger", "error", err)
		os.Exit(1)
	}
	logging.SetGlobalLogger(logger)
	defer logger.Close()

	if err := commands.Execute(); err != nil {
		logger.Error("application failed",
			slog.String("component", "main"),
			slog.String("error", err.Error()))
		os.Exit(1)
	}
}
