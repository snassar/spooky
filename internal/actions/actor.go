package actions

import (
	spookytypes "spooky/internal/types"
	"context"
	"fmt"
	"sync"
	"time"

	spookyfacts "spooky/internal/facts"
	spookylogging "spooky/internal/logging"
	spookytypes "spooky/internal/types"
)

// collectionRunContext holds context for action orchestration
type collectionRunContext struct {
	plan          *spookytypes.OrchestrationPlan
	context       *spookytypes.ActionContext
	logger        spookylogging.Logger
	results       map[string]*spookytypes.RunResult
	rollbackStack []*rollbackEntry
}

// rollbackEntry represents an entry in the rollback stack
type rollbackEntry struct {
	action *spookytypes.Action
	result *spookytypes.RunResult
}

// actionOrchestrationResult holds result of action orchestration
type actionOrchestrationResult struct {
	actionName string
	result     *spookytypes.RunResult
}

// Manager implements the ActionManager interface
type Manager struct {
	// Configuration
	defaultTimeout   time.Duration
	defaultParallel  bool
	customValidators map[string]ActionValidator

	// Dependencies
	factsManager spookyfacts.FactManager

	// State
	actions map[string]*spookytypes.Action
	logger  spookylogging.Logger
	mu      sync.RWMutex
}

// NewManager creates a new ActionManager
func NewManager(logger spookylogging.Logger) *Manager {
	return &Manager{
		defaultTimeout:   30 * time.Minute,
		defaultParallel:  false,
		customValidators: make(map[string]ActionValidator),
		factsManager:     nil,
		actions:          make(map[string]*spookytypes.Action),
		logger:           logger,
	}
}

// NewManagerWithFacts creates a new ActionManager with facts manager
func NewManagerWithFacts(logger spookylogging.Logger, factsManager spookyfacts.FactManager) *Manager {
	return &Manager{
		defaultTimeout:   30 * time.Minute,
		defaultParallel:  false,
		customValidators: make(map[string]ActionValidator),
		factsManager:     factsManager,
		actions:          make(map[string]*spookytypes.Action),
		logger:           logger,
	}
}

// Run orchestrates actions with internal planning
func (m *Manager) Run(ctx context.Context, actions []*spookytypes.Action, context *spookytypes.ActionContext) (*spookytypes.OrchestrationResult, error) {
	m.logger.Debug("Starting action orchestration",
		spookylogging.Int("action_count", len(actions)))

	// Create orchestration plan internally
	plan, err := m.planOrchestration(actions, context)
	if err != nil {
		return nil, fmt.Errorf("failed to create orchestration plan: %w", err)
	}

	// Prepare run context
	runContext := &collectionRunContext{
		plan:    plan,
		context: context,
		logger:  m.logger,
		results: make(map[string]*spookytypes.RunResult),
	}

	// Run based on plan strategy
	if len(plan.ParallelGroups) > 1 && context.Parallel {
		return m.runParallel(ctx, runContext)
	} else {
		return m.runSequential(ctx, runContext)
	}
}

// planOrchestration analyzes dependencies and creates an optimized orchestration plan
func (m *Manager) planOrchestration(actions []*spookytypes.Action, context *spookytypes.ActionContext) (*spookytypes.OrchestrationPlan, error) {
	m.logger.Debug("Planning action orchestration",
		spookylogging.Int("action_count", len(actions)))

	// Build dependency graph
	graph, err := m.buildDependencyGraph(actions)
	if err != nil {
		return nil, fmt.Errorf("failed to build dependency graph: %w", err)
	}

	// Determine orchestration order
	orchestrationOrder, err := graph.topologicalSort()
	if err != nil {
		return nil, fmt.Errorf("failed to determine orchestration order: %w", err)
	}

	// Optimize orchestration strategy
	optimizedPlan, err := m.optimizeOrchestration(actions, orchestrationOrder, graph, context)
	if err != nil {
		return nil, fmt.Errorf("failed to optimize orchestration: %w", err)
	}

	return optimizedPlan, nil
}

// buildDependencyGraph builds a dependency graph from actions
func (m *Manager) buildDependencyGraph(actions []*spookytypes.Action) (*dependencyGraph, error) {
	// TODO: Implement dependency graph building
	return &dependencyGraph{}, nil
}

// optimizeOrchestration optimizes the orchestration strategy
func (m *Manager) optimizeOrchestration(actions []*spookytypes.Action, order []string, graph *dependencyGraph, context *spookytypes.ActionContext) (*spookytypes.OrchestrationPlan, error) {
	// TODO: Implement orchestration optimization
	return &spookytypes.OrchestrationPlan{
		SequentialOrder: order,
		ParallelGroups:  [][]string{order},
	}, nil
}

// runSequential runs actions sequentially
func (m *Manager) runSequential(ctx context.Context, runContext *collectionRunContext) (*spookytypes.OrchestrationResult, error) {
	// TODO: Implement sequential execution
	return &spookytypes.OrchestrationResult{}, nil
}

// runParallel runs actions in parallel
func (m *Manager) runParallel(ctx context.Context, runContext *collectionRunContext) (*spookytypes.OrchestrationResult, error) {
	// TODO: Implement parallel execution
	return &spookytypes.OrchestrationResult{}, nil
}

// runAction runs a single action
func (m *Manager) runAction(ctx context.Context, action *spookytypes.Action, context *spookytypes.ActionContext) (*spookytypes.RunResult, error) {
	// TODO: Implement single action execution
	return &spookytypes.RunResult{}, nil
}

// LoadActions loads actions from the project
func (m *Manager) LoadActions(projectPath string) (*spookytypes.ActionCollection, error) {
	m.logger.Info("Loading actions from project", spookylogging.String("project", projectPath))

	// For now, return an empty collection
	// In a real implementation, this would load actions from actions.hcl files
	collection := &spookytypes.ActionCollection{
		Actions:   make([]*spookytypes.Action, 0),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Metadata:  make(map[string]interface{}),
	}

	m.logger.Info("Loaded actions from project",
		spookylogging.String("project", projectPath),
		spookylogging.Int("actions_count", 0))

	return collection, nil
}

// GetAction gets an action by name
func (m *Manager) GetAction(name string) (*spookytypes.Action, error) {
	if name == "" {
		return nil, fmt.Errorf("action name cannot be empty")
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	action, exists := m.actions[name]
	if !exists {
		return nil, fmt.Errorf("action '%s' not found", name)
	}

	m.logger.Debug("Retrieved action", spookylogging.String("action", name))
	return action, nil
}

// ListActions lists all available actions
func (m *Manager) ListActions() ([]*spookytypes.Action, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	actions := make([]*spookytypes.Action, 0, len(m.actions))
	for _, action := range m.actions {
		actions = append(actions, action)
	}

	m.logger.Debug("Listed actions", spookylogging.Int("count", len(actions)))
	return actions, nil
}

// AddAction adds a new action
func (m *Manager) AddAction(name string, action *spookytypes.Action) error {
	if name == "" {
		return fmt.Errorf("action name cannot be empty")
	}

	if action == nil {
		return fmt.Errorf("action cannot be nil")
	}

	// Validate action before adding
	if err := m.ValidateAction(action); err != nil {
		return fmt.Errorf("action validation failed: %w", err)
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// Check for duplicate action names
	if _, exists := m.actions[name]; exists {
		return fmt.Errorf("action with name '%s' already exists", name)
	}

	m.actions[name] = action
	m.logger.Info("Added action", spookylogging.String("action", name))

	return nil
}

// RemoveAction removes an action
func (m *Manager) RemoveAction(name string) error {
	if name == "" {
		return fmt.Errorf("action name cannot be empty")
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.actions[name]; !exists {
		return fmt.Errorf("action '%s' not found", name)
	}

	delete(m.actions, name)
	m.logger.Info("Removed action", spookylogging.String("action", name))

	return nil
}

// ValidateAction validates an action
func (m *Manager) ValidateAction(action *spookytypes.Action) error {
	if action == nil {
		return fmt.Errorf("action cannot be nil")
	}

	if action.Name == "" {
		return fmt.Errorf("action name is required")
	}

	// Validate action type
	if action.Type == "" {
		return fmt.Errorf("action type is required")
	}

	// Validate action content based on type
	switch action.Type {
	case "command":
		if action.Command == "" {
			return fmt.Errorf("command is required for command actions")
		}
	case "script":
		if action.Script == "" {
			return fmt.Errorf("script is required for script actions")
		}
	case "template_deploy", "template_evaluate", "template_validate", "template_cleanup":
		if action.Template == nil {
			return fmt.Errorf("template configuration is required for template actions")
		}
		if action.Template.Source == "" {
			return fmt.Errorf("template source is required for template actions")
		}
	case "file_copy":
		if action.FileCopy == nil {
			return fmt.Errorf("file copy configuration is required for file copy actions")
		}
		if action.FileCopy.Source == "" {
			return fmt.Errorf("file copy source is required for file copy actions")
		}
		if action.FileCopy.Destination == "" {
			return fmt.Errorf("file copy destination is required for file copy actions")
		}
	default:
		return fmt.Errorf("unsupported action type: %s", action.Type)
	}

	return nil
}

// SetDefaultTimeout sets the default timeout for actions
func (m *Manager) SetDefaultTimeout(timeout time.Duration) {
	m.defaultTimeout = timeout
}

// SetDefaultParallel sets the default parallel execution setting
func (m *Manager) SetDefaultParallel(parallel bool) {
	m.defaultParallel = parallel
}

// RegisterCustomValidator registers a custom validator
func (m *Manager) RegisterCustomValidator(name string, validator ActionValidator) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.customValidators[name] = validator
	m.logger.Debug("Registered custom validator", spookylogging.String("validator", name))
}

// Close closes the manager and cleans up resources
func (m *Manager) Close() error {
	m.logger.Debug("Closing action manager")
	return nil
}

// dependencyGraph represents a dependency graph for actions
type dependencyGraph struct {
	// TODO: Implement dependency graph structure
}

// topologicalSort performs topological sorting on the dependency graph
func (g *dependencyGraph) topologicalSort() ([]string, error) {
	// TODO: Implement topological sorting
	return []string{}, nil
}
