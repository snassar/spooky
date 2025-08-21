package main

import (
	"log/slog"
	"os"
	"spooky/commands"
	"spooky/internal/logging"
)

func main() {
	// Initialize global logger
	logger, err := logging.NewLogger(nil) // Use default config
	if err != nil {
		// Fallback to basic logging
		slog.Error("failed to initialize logger", "error", err)
		os.Exit(1)
	}
	logging.SetGlobalLogger(logger)
	defer logger.Close()

	logger.Info("starting spooky application",
		slog.String("component", "main"),
		slog.String("version", "0.1.0"))

	if err := commands.Execute(); err != nil {
		logger.Error("application failed",
			slog.String("component", "main"),
			slog.String("error", err.Error()))
		os.Exit(1)
	}

	logger.Info("spooky application completed successfully",
		slog.String("component", "main"))
}
