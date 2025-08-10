package acting

import (
	spookylogging "spooky/internal/logging"
	spookysshtypes "spooky/internal/ssh/types"
)

// Manager implements ActingEngine interface
type Manager struct {
	config *spookysshtypes.ActingConfig
	logger spookylogging.Logger
}

// NewManager creates a new acting manager
func NewManager(config *spookysshtypes.ActingConfig, logger spookylogging.Logger) *Manager {
	return &Manager{
		config: config,
		logger: logger,
	}
}

// RunAction runs an action on the SSH connection
func (m *Manager) RunAction(connection *spookysshtypes.SSHConnection, action *spookysshtypes.SSHAction) (*spookysshtypes.ActionResult, error) {
	// TODO: Implement action running logic
	m.logger.Info("Action run", spookylogging.String("host", connection.Host), spookylogging.String("action", action.Name))
	return &spookysshtypes.ActionResult{}, nil
}

// RunTemplate runs a template action on the SSH connection
func (m *Manager) RunTemplate(connection *spookysshtypes.SSHConnection, template *spookysshtypes.TemplateAction) (*spookysshtypes.ActionResult, error) {
	// TODO: Implement template running logic
	m.logger.Info("Template run", spookylogging.String("host", connection.Host), spookylogging.String("template", template.Name))
	return &spookysshtypes.ActionResult{}, nil
}

// RunSequential runs actions sequentially
func (m *Manager) RunSequential(connection *spookysshtypes.SSHConnection, actions []*spookysshtypes.SSHAction) (*spookysshtypes.ActionResult, error) {
	// TODO: Implement sequential running logic
	m.logger.Info("Sequential running requested", spookylogging.String("host", connection.Host), spookylogging.Int("action_count", len(actions)))
	return &spookysshtypes.ActionResult{}, nil
}

// RunParallel runs actions in parallel
func (m *Manager) RunParallel(connection *spookysshtypes.SSHConnection, actions []*spookysshtypes.SSHAction) (*spookysshtypes.ActionResult, error) {
	// TODO: Implement parallel running logic
	m.logger.Info("Parallel running requested", spookylogging.String("host", connection.Host), spookylogging.Int("action_count", len(actions)))
	return &spookysshtypes.ActionResult{}, nil
}
