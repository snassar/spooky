// Package config provides configuration management functionality for the spooky codebase.
package config

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/hashicorp/hcl/v2/hclsimple"

	spookyinterfaces "spooky/internal/interfaces"
	spookytypes "spooky/internal/types"
	spookytypeslogging "spooky/internal/types/logging"
	spookytypesschemas "spooky/internal/types/schemas"
)

// Integration implements the ConfigIntegration interface
type Integration struct {
	logger spookytypeslogging.Logger
}

// NewIntegration creates a new config integration
func NewIntegration(logger spookytypeslogging.Logger) spookyinterfaces.ConfigIntegration {
	return &Integration{
		logger: logger,
	}
}

// LoadConfig loads configuration from the given source
func (i *Integration) LoadConfig(ctx context.Context, source string) (*spookytypes.Config, error) {
	if source == "" {
		return nil, fmt.Errorf("config source cannot be empty")
	}

	// Check if source file exists
	if _, err := os.Stat(source); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file not found: %s", source)
	}

	// Parse HCL configuration
	var config spookytypes.Config
	if err := hclsimple.DecodeFile(source, nil, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file %s: %w", source, err)
	}

	i.logger.Info("Configuration loaded successfully", map[string]interface{}{
		"source": source,
	})

	return &config, nil
}

// ValidateConfig validates configuration
func (i *Integration) ValidateConfig(ctx context.Context, config *spookytypes.Config) (*spookytypes.ValidationResult, error) {
	if config == nil {
		return &spookytypes.ValidationResult{
			Valid:    false,
			Errors:   []spookytypesschemas.SchemaError{{Message: "config cannot be nil"}},
			Warnings: []spookytypesschemas.SchemaError{},
		}, nil
	}

	// Basic validation
	var errors []spookytypesschemas.SchemaError
	var warnings []spookytypesschemas.SchemaError

	// Validate global configuration
	if config.Global != nil {
		if config.Global.DefaultProjectPath == "" {
			warnings = append(warnings, spookytypesschemas.SchemaError{
				Message: "default project path is recommended",
			})
		}
	}

	// Validate logging configuration
	if config.Logging != nil {
		if config.Logging.Level == "" {
			warnings = append(warnings, spookytypesschemas.SchemaError{
				Message: "logging level is recommended",
			})
		}
	}

	valid := len(errors) == 0

	i.logger.Info("Configuration validation completed", map[string]interface{}{
		"valid":    valid,
		"errors":   len(errors),
		"warnings": len(warnings),
	})

	return &spookytypes.ValidationResult{
		Valid:    valid,
		Errors:   errors,
		Warnings: warnings,
	}, nil
}

// SaveConfig saves configuration to the given destination
func (i *Integration) SaveConfig(ctx context.Context, config *spookytypes.Config, destination string) error {
	if config == nil {
		return fmt.Errorf("config cannot be nil")
	}

	if destination == "" {
		return fmt.Errorf("destination cannot be empty")
	}

	// Ensure destination directory exists
	destDir := filepath.Dir(destination)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	// Save configuration to file
	// TODO: Implement HCL marshaling
	// For now, return an error indicating this is not implemented
	return fmt.Errorf("config saving is not yet implemented")

	i.logger.Info("Configuration saved successfully", map[string]interface{}{
		"destination": destination,
	})

	return nil
}
