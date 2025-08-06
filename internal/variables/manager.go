package variables

import (
	"context"
	"fmt"
	"sync"

	spookylogging "spooky/internal/logging"
	spookyvariablesimportexport "spooky/internal/variables/importexport"
	spookyvariablesloading "spooky/internal/variables/loading"
	spookyvariablesresolution "spooky/internal/variables/resolution"
	spookyvariablestypes "spooky/internal/variables/types"
	spookyvariablesvalidation "spooky/internal/variables/validation"
)

// Manager implements VariableManager interface
type Manager struct {
	config              *spookyvariablestypes.Config
	loadingManager      spookyvariablesloading.LoadingManager
	resolutionManager   spookyvariablesresolution.ResolutionManager
	validationManager   spookyvariablesvalidation.ValidationManager
	importExportManager spookyvariablesimportexport.ImportExportManager
	logger              spookylogging.Logger
	variables           map[string]*spookyvariablestypes.Variable
	mu                  sync.RWMutex
}

// NewManager creates a new variable manager
func NewManager(
	config *spookyvariablestypes.Config,
	loadingManager spookyvariablesloading.LoadingManager,
	resolutionManager spookyvariablesresolution.ResolutionManager,
	validationManager spookyvariablesvalidation.ValidationManager,
	importExportManager spookyvariablesimportexport.ImportExportManager,
	logger spookylogging.Logger,
) VariableManager {
	return &Manager{
		config:              config,
		loadingManager:      loadingManager,
		resolutionManager:   resolutionManager,
		validationManager:   validationManager,
		importExportManager: importExportManager,
		logger:              logger,
		variables:           make(map[string]*spookyvariablestypes.Variable),
	}
}

// LoadVariables loads variables from the specified path
func (m *Manager) LoadVariables(ctx context.Context, path string) (*spookyvariablestypes.VariableCollection, error) {
	// Handle nil dependencies gracefully
	if m.loadingManager == nil {
		// Return empty collection if loading manager is nil
		return &spookyvariablestypes.VariableCollection{
			Variables: []*spookyvariablestypes.Variable{},
			Path:      path,
		}, nil
	}

	// 1. Load variables using loading manager
	variables, err := m.loadingManager.LoadFromDirectory(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("failed to load variables: %w", err)
	}

	// 2. Create collection
	collection := &spookyvariablestypes.VariableCollection{
		Variables: variables,
		Path:      path,
	}

	// 3. Validate collection if validation manager is available
	if m.validationManager != nil {
		result, err := m.validationManager.ValidateCollection(ctx, collection)
		if err != nil {
			return nil, fmt.Errorf("validation failed: %w", err)
		}

		if !result.Valid {
			return nil, fmt.Errorf("variable validation failed: %v", result.Errors)
		}
	}

	// 4. Store variables in memory
	m.mu.Lock()
	for _, variable := range variables {
		m.variables[variable.Name] = variable
	}
	m.mu.Unlock()

	return collection, nil
}

// GetVariable gets a variable by name
func (m *Manager) GetVariable(ctx context.Context, name string) (*spookyvariablestypes.Variable, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	variable, exists := m.variables[name]
	if !exists {
		return nil, fmt.Errorf("variable %s not found", name)
	}

	return variable, nil
}

// SetVariable sets a variable
func (m *Manager) SetVariable(ctx context.Context, variable *spookyvariablestypes.Variable) error {
	// Validate variable
	result, err := m.validationManager.ValidateVariable(ctx, variable)
	if err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	if !result.Valid {
		return fmt.Errorf("variable validation failed: %v", result.Errors)
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	m.variables[variable.Name] = variable
	return nil
}

// DeleteVariable deletes a variable by name
func (m *Manager) DeleteVariable(ctx context.Context, name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.variables[name]; !exists {
		return fmt.Errorf("variable %s not found", name)
	}

	delete(m.variables, name)
	return nil
}

// ListVariables lists all variables
func (m *Manager) ListVariables(ctx context.Context) ([]*spookyvariablestypes.Variable, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	variables := make([]*spookyvariablestypes.Variable, 0, len(m.variables))
	for _, variable := range m.variables {
		variables = append(variables, variable)
	}

	return variables, nil
}

// CreateContext creates a variable context from variables
func (m *Manager) CreateContext(ctx context.Context, variables []*spookyvariablestypes.Variable) (*spookyvariablestypes.VariableContext, error) {
	context := &spookyvariablestypes.VariableContext{
		Variables: make(map[string]*spookyvariablestypes.Variable),
	}

	for _, variable := range variables {
		context.Variables[variable.Name] = variable
	}

	return context, nil
}

// ResolveContext resolves all variables in a context
func (m *Manager) ResolveContext(ctx context.Context, context *spookyvariablestypes.VariableContext) error {
	// 1. Get all variables from context
	variables := make([]*spookyvariablestypes.Variable, 0, len(context.Variables))
	for _, variable := range context.Variables {
		variables = append(variables, variable)
	}

	// 2. Resolve dependencies
	if err := m.resolutionManager.ResolveDependencies(ctx, variables); err != nil {
		return fmt.Errorf("failed to resolve dependencies: %w", err)
	}

	// 3. Update context with resolved variables
	for _, variable := range variables {
		context.Variables[variable.Name] = variable
	}

	return nil
}

// ExportVariables exports variables to a file
func (m *Manager) ExportVariables(ctx context.Context, format spookyvariablestypes.ExportFormat, path string) error {
	// Get all variables
	variables, err := m.ListVariables(ctx)
	if err != nil {
		return fmt.Errorf("failed to list variables: %w", err)
	}

	// Export based on format
	switch format {
	case spookyvariablestypes.ExportFormatHCL:
		return m.importExportManager.ExportToHCL(ctx, variables, path)
	case spookyvariablestypes.ExportFormatJSON:
		return m.importExportManager.ExportToJSON(ctx, variables, path)
	default:
		return fmt.Errorf("unsupported export format: %s", format)
	}
}

// ImportVariables imports variables from a file
func (m *Manager) ImportVariables(ctx context.Context, format spookyvariablestypes.ImportFormat, path string) error {
	var variables []*spookyvariablestypes.Variable
	var err error

	// Import based on format
	switch format {
	case spookyvariablestypes.ImportFormatHCL:
		variables, err = m.importExportManager.ImportFromHCL(ctx, path)
	case spookyvariablestypes.ImportFormatJSON:
		variables, err = m.importExportManager.ImportFromJSON(ctx, path)
	default:
		return fmt.Errorf("unsupported import format: %s", format)
	}

	if err != nil {
		return fmt.Errorf("failed to import variables: %w", err)
	}

	// Add imported variables
	for _, variable := range variables {
		if err := m.SetVariable(ctx, variable); err != nil {
			return fmt.Errorf("failed to set imported variable %s: %w", variable.Name, err)
		}
	}

	return nil
}

// ValidateVariables validates a list of variables
func (m *Manager) ValidateVariables(ctx context.Context, variables []*spookyvariablestypes.Variable) (*spookyvariablestypes.ValidationResult, error) {
	collection := &spookyvariablestypes.VariableCollection{
		Variables: variables,
	}

	return m.validationManager.ValidateCollection(ctx, collection)
}

// ValidateContext validates a variable context
func (m *Manager) ValidateContext(ctx context.Context, context *spookyvariablestypes.VariableContext) (*spookyvariablestypes.ValidationResult, error) {
	return m.validationManager.ValidateContext(ctx, context)
}

// Coordinator integration methods
func (m *Manager) LoadVariablesForProject(projectPath string) (*spookyvariablestypes.VariableCollection, error) {
	// Load variables from project directory
	variablesPath := fmt.Sprintf("%s/variables", projectPath)
	return m.LoadVariables(context.Background(), variablesPath)
}

func (m *Manager) ResolveVariablesForContext(variableContext *spookyvariablestypes.VariableContext) error {
	return m.ResolveContext(context.Background(), variableContext)
}

func (m *Manager) ValidateVariablesForProject(projectPath string) (*spookyvariablestypes.ValidationResult, error) {
	// Load and validate variables for project
	collection, err := m.LoadVariablesForProject(projectPath)
	if err != nil {
		return nil, err
	}

	return m.validationManager.ValidateCollection(context.Background(), collection)
}

func (m *Manager) ExportVariablesForProject(projectPath string, format spookyvariablestypes.ExportFormat, outputPath string) error {
	// Load variables and export them
	collection, err := m.LoadVariablesForProject(projectPath)
	if err != nil {
		return err
	}

	switch format {
	case spookyvariablestypes.ExportFormatHCL:
		return m.importExportManager.ExportToHCL(context.Background(), collection.Variables, outputPath)
	case spookyvariablestypes.ExportFormatJSON:
		return m.importExportManager.ExportToJSON(context.Background(), collection.Variables, outputPath)
	default:
		return fmt.Errorf("unsupported export format: %s", format)
	}
}
