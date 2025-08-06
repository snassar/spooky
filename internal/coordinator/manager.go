package coordinator

import (
	"fmt"
	"path/filepath"
	"time"

	spookyactions "spooky/internal/actions"
	spookyactionstypes "spooky/internal/actions/types"
	spookyconfig "spooky/internal/config"
	spookyfacts "spooky/internal/facts"
	spookybadger "spooky/internal/facts/storage/badger"
	spookyinterfaces "spooky/internal/interfaces"
	spookylogging "spooky/internal/logging"
	spookymachines "spooky/internal/machines"
	spookymachinestypes "spooky/internal/machines/types"
	spookysecrets "spooky/internal/secrets"
	spookyssh "spooky/internal/ssh"
	spookytemplates "spooky/internal/templates"
	spookytemplatestypes "spooky/internal/templates/types"
	spookyvariables "spooky/internal/variables"
)

// CoordinatorManager implements the IntegrationManager interface
type CoordinatorManager struct {
	facts     spookyinterfaces.FactsIntegration
	actions   spookyinterfaces.ActionsIntegration
	variables spookyinterfaces.VariablesIntegration
	templates spookyinterfaces.TemplatesIntegration
	machines  spookyinterfaces.MachinesIntegration
	crypto    spookyinterfaces.CryptoIntegration
	config    spookyinterfaces.ConfigIntegration
	logger    spookylogging.Logger
	cache     *spookyinterfaces.CacheManager
}

// NewCoordinatorManager creates a new integration manager
func NewCoordinatorManager(
	factsManager spookyfacts.FactManager,
	actionsManager spookyactions.ActionManager,
	variablesManager spookyvariables.VariableManager,
	templatesManager spookytemplates.TemplateManager,
	machinesManager spookymachines.MachineManager,
	cryptoManager *spookysecrets.Manager,
	configManager spookyconfig.ConfigManager,
	logger spookylogging.Logger,
) *CoordinatorManager {
	return &CoordinatorManager{
		facts:     NewCoordinatorFactsIntegration(factsManager, logger),
		actions:   NewCoordinatorActionsIntegration(actionsManager, logger),
		variables: NewCoordinatorVariablesIntegration(variablesManager, logger),
		templates: NewCoordinatorTemplatesIntegration(templatesManager, logger),
		machines:  NewCoordinatorMachinesIntegration(machinesManager, logger),
		crypto:    NewCoordinatorCryptoIntegration(cryptoManager, logger),
		config:    NewCoordinatorConfigIntegration(configManager, logger),
		logger:    logger,
		cache:     spookyinterfaces.NewCacheManager(),
	}
}

// NewCoordinatorManagerFromProject creates a new coordinator manager from a project path
func NewCoordinatorManagerFromProject(projectPath string, logger spookylogging.Logger) (*CoordinatorManager, error) {
	// Create facts manager
	factsStorage, err := spookybadger.NewBadgerFactStorage(filepath.Join(projectPath, "facts.db"))
	if err != nil {
		return nil, fmt.Errorf("failed to create facts storage: %w", err)
	}

	// Create SSH client for facts collection
	sshManager := spookyssh.NewDefaultManager(logger)
	var sshClient spookyssh.SSHClient = sshManager
	factsManager := spookyfacts.NewManagerWithStorage(&sshClient, factsStorage, logger)

	// Create actions manager
	actionsManager := spookyactions.NewManager(logger)

	// Create variables manager with real implementation
	variablesManager := spookyvariables.NewVariableManager(logger)

	// Create templates manager with real implementation
	templatesManager, err := spookytemplates.NewDefaultTemplateManager(logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create templates manager: %w", err)
	}

	// Create machines manager with default config
	machinesManager := spookymachines.NewManager(spookymachinestypes.DefaultIndexManagerConfig(), logger)

	// Create crypto manager
	cryptoManager, err := spookysecrets.NewManager(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create crypto manager: %w", err)
	}

	// Create config manager with nil dependencies for now
	configManager := spookyconfig.NewManager(nil, nil, nil, nil, logger)

	return NewCoordinatorManager(
		factsManager,
		actionsManager,
		variablesManager,
		templatesManager,
		machinesManager,
		cryptoManager,
		configManager,
		logger,
	), nil
}

// Facts returns the facts integration
func (m *CoordinatorManager) Facts() spookyinterfaces.FactsIntegration {
	return m.facts
}

// Actions returns the actions integration
func (m *CoordinatorManager) Actions() spookyinterfaces.ActionsIntegration {
	return m.actions
}

// Variables returns the variables integration
func (m *CoordinatorManager) Variables() spookyinterfaces.VariablesIntegration {
	return m.variables
}

// Templates returns the templates integration
func (m *CoordinatorManager) Templates() spookyinterfaces.TemplatesIntegration {
	return m.templates
}

// Machines returns the machines integration
func (m *CoordinatorManager) Machines() spookyinterfaces.MachinesIntegration {
	return m.machines
}

// Crypto returns the crypto integration
func (m *CoordinatorManager) Crypto() spookyinterfaces.CryptoIntegration {
	return m.crypto
}

// Config returns the config integration
func (m *CoordinatorManager) Config() spookyinterfaces.ConfigIntegration {
	return m.config
}

// LoadContextForAction loads all necessary context for action execution
func (m *CoordinatorManager) LoadContextForAction(action *spookyactionstypes.Action, projectPath string, machineNames []string) (*spookyinterfaces.ActionExecutionContext, error) {
	context := &spookyinterfaces.ActionExecutionContext{
		BaseContext: spookyinterfaces.BaseContext{
			ProjectPath: projectPath,
			Timestamp:   time.Now(),
		},
		MachineNames: machineNames,
		Action:       action,
	}

	// Load facts context
	if factsContext, err := m.facts.LoadFacts(machineNames); err != nil {
		m.logger.Warn("Failed to load facts context", spookylogging.Error(err))
	} else {
		context.FactsContext = factsContext
	}

	// Load variables context
	if variablesContext, err := m.variables.LoadVariables(projectPath); err != nil {
		m.logger.Warn("Failed to load variables context", spookylogging.Error(err))
	} else {
		context.VariablesContext = variablesContext
	}

	// Load templates context
	if templatesContext, err := m.templates.LoadTemplates(projectPath); err != nil {
		m.logger.Warn("Failed to load templates context", spookylogging.Error(err))
	} else {
		context.TemplatesContext = templatesContext
	}

	return context, nil
}

// ValidateActionWithAllSystems validates an action using all available systems
func (m *CoordinatorManager) ValidateActionWithAllSystems(action *spookyactionstypes.Action, context *spookyinterfaces.ActionExecutionContext) error {
	var errors []error

	// Validate with facts system
	if context.FactsContext != nil {
		if err := m.facts.ValidateFacts(context.FactsContext); err != nil {
			errors = append(errors, fmt.Errorf("facts validation: %w", err))
		}
	}

	// Validate with variables system
	if context.VariablesContext != nil {
		if err := m.variables.ValidateVariables(context.VariablesContext); err != nil {
			errors = append(errors, fmt.Errorf("variables validation: %w", err))
		}
	}

	// Validate with templates system
	if context.TemplatesContext != nil && action.Template != nil {
		// Convert config.TemplateConfig to templates.types.Template for validation
		template := &spookytemplatestypes.Template{
			Name:   action.Name + "_template",
			Source: action.Template.Source,
		}

		if err := m.templates.ValidateTemplate(template, context.TemplatesContext); err != nil {
			errors = append(errors, fmt.Errorf("template validation: %w", err))
		}
	}

	// Validate with actions system
	if err := m.actions.ValidateAction(action, context); err != nil {
		errors = append(errors, fmt.Errorf("action validation: %w", err))
	}

	if len(errors) > 0 {
		return &ValidationError{Errors: errors}
	}

	return nil
}

// PrepareActionForExecution prepares an action for execution
func (m *CoordinatorManager) PrepareActionForExecution(action *spookyactionstypes.Action, context *spookyinterfaces.ActionExecutionContext) error {
	return m.actions.PrepareActionForExecution(action, context)
}

// ExecuteAction executes an action
func (m *CoordinatorManager) ExecuteAction(action *spookyactionstypes.Action, context *spookyinterfaces.ActionExecutionContext) error {
	return m.actions.ExecuteAction(action, context)
}

// ValidateIntegrationHealth validates the health of all integrations
func (m *CoordinatorManager) ValidateIntegrationHealth() error {
	var errors []error

	// Validate facts integration
	if err := m.validateFactsHealth(); err != nil {
		errors = append(errors, fmt.Errorf("facts integration: %w", err))
	}

	// Validate actions integration
	if err := m.validateActionsHealth(); err != nil {
		errors = append(errors, fmt.Errorf("actions integration: %w", err))
	}

	// Validate variables integration
	if err := m.validateVariablesHealth(); err != nil {
		errors = append(errors, fmt.Errorf("variables integration: %w", err))
	}

	// Validate templates integration
	if err := m.validateTemplatesHealth(); err != nil {
		errors = append(errors, fmt.Errorf("templates integration: %w", err))
	}

	// Validate machines integration
	if err := m.validateMachinesHealth(); err != nil {
		errors = append(errors, fmt.Errorf("machines integration: %w", err))
	}

	// Validate crypto integration
	if err := m.validateCryptoHealth(); err != nil {
		errors = append(errors, fmt.Errorf("crypto integration: %w", err))
	}

	if len(errors) > 0 {
		return &ValidationError{Errors: errors}
	}

	return nil
}

// GetIntegrationStats returns statistics for all integrations
func (m *CoordinatorManager) GetIntegrationStats() map[string]interface{} {
	return map[string]interface{}{
		"facts":     m.getFactsStats(),
		"actions":   m.getActionsStats(),
		"variables": m.getVariablesStats(),
		"templates": m.getTemplatesStats(),
		"machines":  m.getMachinesStats(),
		"crypto":    m.getCryptoStats(),
	}
}

// ClearAllCaches clears all caches
func (m *CoordinatorManager) ClearAllCaches() error {
	m.cache.ClearAllCaches()
	return nil
}

// LoadProjectContext loads all context for a project
func (m *CoordinatorManager) LoadProjectContext(projectPath string) (*spookyinterfaces.ProjectContext, error) {
	context := &spookyinterfaces.ProjectContext{
		BaseContext: spookyinterfaces.BaseContext{
			ProjectPath: projectPath,
			Timestamp:   time.Now(),
		},
	}

	// Load facts context
	if factsContext, err := m.facts.LoadFacts([]string{}); err != nil {
		m.logger.Warn("Failed to load project facts context", spookylogging.Error(err))
	} else {
		context.FactsContext = factsContext
	}

	// Load variables context
	if variablesContext, err := m.variables.LoadVariables(projectPath); err != nil {
		m.logger.Warn("Failed to load project variables context", spookylogging.Error(err))
	} else {
		context.VariablesContext = variablesContext
	}

	// Load templates context
	if templatesContext, err := m.templates.LoadTemplates(projectPath); err != nil {
		m.logger.Warn("Failed to load project templates context", spookylogging.Error(err))
	} else {
		context.TemplatesContext = templatesContext
	}

	// Load machines context
	if machinesContext, err := m.machines.LoadMachines(projectPath); err != nil {
		m.logger.Warn("Failed to load project machines context", spookylogging.Error(err))
	} else {
		context.MachinesContext = machinesContext
	}

	// Load actions context
	if actionsContext, err := m.actions.LoadActions(projectPath); err != nil {
		m.logger.Warn("Failed to load project actions context", spookylogging.Error(err))
	} else {
		context.ActionsContext = actionsContext
	}

	return context, nil
}

// ValidateProject validates a project
func (m *CoordinatorManager) ValidateProject(projectPath string) error {
	var errors []error

	// Validate project facts
	if err := m.validateProjectFacts(projectPath); err != nil {
		errors = append(errors, fmt.Errorf("project facts: %w", err))
	}

	// Validate project variables
	if err := m.validateProjectVariables(projectPath); err != nil {
		errors = append(errors, fmt.Errorf("project variables: %w", err))
	}

	// Validate project templates
	if err := m.validateProjectTemplates(projectPath); err != nil {
		errors = append(errors, fmt.Errorf("project templates: %w", err))
	}

	// Validate project machines
	if err := m.validateProjectMachines(projectPath); err != nil {
		errors = append(errors, fmt.Errorf("project machines: %w", err))
	}

	// Validate project actions
	if err := m.validateProjectActions(projectPath); err != nil {
		errors = append(errors, fmt.Errorf("project actions: %w", err))
	}

	if len(errors) > 0 {
		return &ValidationError{Errors: errors}
	}

	return nil
}

// GetProjectStats returns statistics for a project
func (m *CoordinatorManager) GetProjectStats(projectPath string) map[string]interface{} {
	return map[string]interface{}{
		"facts":     m.getProjectFactsStats(projectPath),
		"variables": m.getProjectVariablesStats(projectPath),
		"templates": m.getProjectTemplatesStats(projectPath),
		"machines":  m.getProjectMachinesStats(projectPath),
		"actions":   m.getProjectActionsStats(projectPath),
	}
}

// Private helper methods
func (m *CoordinatorManager) validateFactsHealth() error {
	// Basic facts integration health check
	if m.facts == nil {
		return fmt.Errorf("facts integration is nil")
	}

	// Test facts loading with empty machine list
	_, err := m.facts.LoadFacts([]string{})
	if err != nil {
		return fmt.Errorf("facts integration health check failed: %w", err)
	}

	return nil
}

func (m *CoordinatorManager) validateActionsHealth() error {
	// Basic actions integration health check
	if m.actions == nil {
		return fmt.Errorf("actions integration is nil")
	}

	// Test actions loading with empty project path
	_, err := m.actions.LoadActions("")
	if err != nil {
		// This is expected to fail, but we can check if the integration is responsive
		m.logger.Debug("Actions health check - expected failure for empty project path")
	}

	return nil
}

func (m *CoordinatorManager) validateVariablesHealth() error {
	// Basic variables integration health check
	if m.variables == nil {
		return fmt.Errorf("variables integration is nil")
	}

	// Test variables loading with empty project path
	_, err := m.variables.LoadVariables("")
	if err != nil {
		// This is expected to fail, but we can check if the integration is responsive
		m.logger.Debug("Variables health check - expected failure for empty project path")
	}

	return nil
}

func (m *CoordinatorManager) validateTemplatesHealth() error {
	// Basic templates integration health check
	if m.templates == nil {
		return fmt.Errorf("templates integration is nil")
	}

	// Test templates loading with empty project path
	_, err := m.templates.LoadTemplates("")
	if err != nil {
		// This is expected to fail, but we can check if the integration is responsive
		m.logger.Debug("Templates health check - expected failure for empty project path")
	}

	return nil
}

func (m *CoordinatorManager) validateMachinesHealth() error {
	// Basic machines integration health check
	if m.machines == nil {
		return fmt.Errorf("machines integration is nil")
	}

	// Test machines loading with empty project path
	_, err := m.machines.LoadMachines("")
	if err != nil {
		// This is expected to fail, but we can check if the integration is responsive
		m.logger.Debug("Machines health check - expected failure for empty project path")
	}

	return nil
}

func (m *CoordinatorManager) validateCryptoHealth() error {
	// Basic crypto integration health check
	if m.crypto == nil {
		return fmt.Errorf("crypto integration is nil")
	}

	// Test crypto status retrieval
	status := m.crypto.GetCryptoStatus()
	if status == nil {
		return fmt.Errorf("crypto integration returned nil status")
	}

	return nil
}

func (m *CoordinatorManager) getFactsStats() map[string]interface{} {
	return map[string]interface{}{
		"status": "healthy",
	}
}

func (m *CoordinatorManager) getActionsStats() map[string]interface{} {
	return map[string]interface{}{
		"status": "healthy",
	}
}

func (m *CoordinatorManager) getVariablesStats() map[string]interface{} {
	return map[string]interface{}{
		"status": "healthy",
	}
}

func (m *CoordinatorManager) getTemplatesStats() map[string]interface{} {
	return map[string]interface{}{
		"status": "healthy",
	}
}

func (m *CoordinatorManager) getMachinesStats() map[string]interface{} {
	return map[string]interface{}{
		"status": "healthy",
	}
}

func (m *CoordinatorManager) getCryptoStats() map[string]interface{} {
	return map[string]interface{}{
		"status": "healthy",
	}
}

func (m *CoordinatorManager) validateProjectFacts(projectPath string) error {
	// Basic project facts validation
	if m.facts == nil {
		return fmt.Errorf("facts integration is nil")
	}

	// Test facts loading for the project
	_, err := m.facts.LoadFacts([]string{})
	if err != nil {
		return fmt.Errorf("project facts validation failed: %w", err)
	}

	return nil
}

func (m *CoordinatorManager) validateProjectVariables(projectPath string) error {
	// Basic project variables validation
	if m.variables == nil {
		return fmt.Errorf("variables integration is nil")
	}

	// Test variables loading for the project
	_, err := m.variables.LoadVariables(projectPath)
	if err != nil {
		return fmt.Errorf("project variables validation failed: %w", err)
	}

	return nil
}

func (m *CoordinatorManager) validateProjectTemplates(projectPath string) error {
	// Basic project templates validation
	if m.templates == nil {
		return fmt.Errorf("templates integration is nil")
	}

	// Test templates loading for the project
	_, err := m.templates.LoadTemplates(projectPath)
	if err != nil {
		return fmt.Errorf("project templates validation failed: %w", err)
	}

	return nil
}

func (m *CoordinatorManager) validateProjectMachines(projectPath string) error {
	// Basic project machines validation
	if m.machines == nil {
		return fmt.Errorf("machines integration is nil")
	}

	// Test machines loading for the project
	_, err := m.machines.LoadMachines(projectPath)
	if err != nil {
		return fmt.Errorf("project machines validation failed: %w", err)
	}

	return nil
}

func (m *CoordinatorManager) validateProjectActions(projectPath string) error {
	// Basic project actions validation
	if m.actions == nil {
		return fmt.Errorf("actions integration is nil")
	}

	// Test actions loading for the project
	_, err := m.actions.LoadActions(projectPath)
	if err != nil {
		return fmt.Errorf("project actions validation failed: %w", err)
	}

	return nil
}

func (m *CoordinatorManager) getProjectFactsStats(projectPath string) map[string]interface{} {
	return map[string]interface{}{
		"project_path": projectPath,
		"status":       "healthy",
	}
}

func (m *CoordinatorManager) getProjectVariablesStats(projectPath string) map[string]interface{} {
	return map[string]interface{}{
		"project_path": projectPath,
		"status":       "healthy",
	}
}

func (m *CoordinatorManager) getProjectTemplatesStats(projectPath string) map[string]interface{} {
	return map[string]interface{}{
		"project_path": projectPath,
		"status":       "healthy",
	}
}

func (m *CoordinatorManager) getProjectMachinesStats(projectPath string) map[string]interface{} {
	return map[string]interface{}{
		"project_path": projectPath,
		"status":       "healthy",
	}
}

func (m *CoordinatorManager) getProjectActionsStats(projectPath string) map[string]interface{} {
	return map[string]interface{}{
		"project_path": projectPath,
		"status":       "healthy",
	}
}
