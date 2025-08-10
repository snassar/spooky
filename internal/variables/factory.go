package variables

import (
	spookyinterfaces "spooky/internal/interfaces"
	spookytypes "spooky/internal/types"
)

// NewVariableManager creates a new VariableManager with all its dependencies
func NewVariableManager(logger spookyinterfaces.Logger) spookyinterfaces.VariableManager {
	// Create default configuration
	config := &spookytypes.VariableValidationConfig{
		// TODO: Configure validation settings when types are finalized
	}

	// Create sub-managers with nil implementations for now
	// TODO: Implement proper sub-managers when interfaces are finalized
	var loadingManager spookyinterfaces.VariableLoadingManager = nil
	var resolutionManager spookyinterfaces.ResolutionManager = nil
	var validationManager spookyinterfaces.VariableValidationManager = nil
	var importExportManager spookyinterfaces.ImportExportManager = nil

	// Create main manager
	return NewManager(
		config,
		loadingManager,
		resolutionManager,
		validationManager,
		importExportManager,
		logger,
	)
}
