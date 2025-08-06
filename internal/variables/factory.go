package variables

import (
	"spooky/internal/logging"
	"spooky/internal/variables/importexport"
	"spooky/internal/variables/loading"
	"spooky/internal/variables/resolution"
	"spooky/internal/variables/types"
	"spooky/internal/variables/validation"
)

// NewVariableManager creates a new VariableManager with all its dependencies
func NewVariableManager(logger logging.Logger) VariableManager {
	// Create default configuration
	config := &types.Config{
		LoadingConfig: &types.LoadingConfig{
			DefaultEncoding:   "utf-8",
			MaxFileSize:       1024 * 1024, // 1MB
			AllowedExtensions: []string{".hcl", ".json"},
		},
		ResolutionConfig: &types.ResolutionConfig{
			MaxRecursionDepth: 10,
			DefaultValues:     make(map[string]interface{}),
			StrictMode:        false,
		},
		ValidationConfig: &types.ValidationConfig{
			ValidationRules:     &types.ValidationRules{},
			StrictValidation:    false,
			MaxValidationErrors: 100,
		},
		ImportExportConfig: &types.ImportExportConfig{
			ExportOptions: &types.ExportOptions{
				IncludeMetadata: true,
				PrettyPrint:     true,
			},
			ImportOptions: &types.ImportOptions{
				MergePolicy: "overwrite",
				Overwrite:   false,
			},
			DefaultFormat: types.ExportFormatHCL,
		},
	}

	// Create sub-managers
	loadingManager := loading.NewManager(config.LoadingConfig, logger)
	resolutionManager := resolution.NewManager(config.ResolutionConfig, logger)
	validationManager := validation.NewManager(config.ValidationConfig, logger)
	importExportManager := importexport.NewManager(config.ImportExportConfig, logger)

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
