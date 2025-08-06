package variables

import (
	spookylogging "spooky/internal/logging"
	spookyvariablesimportexport "spooky/internal/variables/importexport"
	spookyvariablesloading "spooky/internal/variables/loading"
	spookyvariablesresolution "spooky/internal/variables/resolution"
	spookyvariablestypes "spooky/internal/variables/types"
	spookyvariablesvalidation "spooky/internal/variables/validation"
)

// NewVariableManager creates a new VariableManager with all its dependencies
func NewVariableManager(logger spookylogging.Logger) VariableManager {
	// Create default configuration
	config := &spookyvariablestypes.Config{
		LoadingConfig: &spookyvariablestypes.LoadingConfig{
			DefaultEncoding:   "utf-8",
			MaxFileSize:       1024 * 1024, // 1MB
			AllowedExtensions: []string{".hcl", ".json"},
		},
		ResolutionConfig: &spookyvariablestypes.ResolutionConfig{
			MaxRecursionDepth: 10,
			DefaultValues:     make(map[string]interface{}),
			StrictMode:        false,
		},
		ValidationConfig: &spookyvariablestypes.ValidationConfig{
			ValidationRules:     &spookyvariablestypes.ValidationRules{},
			StrictValidation:    false,
			MaxValidationErrors: 100,
		},
		ImportExportConfig: &spookyvariablestypes.ImportExportConfig{
			ExportOptions: &spookyvariablestypes.ExportOptions{
				IncludeMetadata: true,
				PrettyPrint:     true,
			},
			ImportOptions: &spookyvariablestypes.ImportOptions{
				MergePolicy: "overwrite",
				Overwrite:   false,
			},
			DefaultFormat: spookyvariablestypes.ExportFormatHCL,
		},
	}

	// Create sub-managers
	loadingManager := spookyvariablesloading.NewManager(config.LoadingConfig, logger)
	resolutionManager := spookyvariablesresolution.NewManager(config.ResolutionConfig, logger)
	validationManager := spookyvariablesvalidation.NewManager(config.ValidationConfig, logger)
	importExportManager := spookyvariablesimportexport.NewManager(config.ImportExportConfig, logger)

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
