package actions

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"spooky/internal/actions/types"
	"spooky/internal/logging"
	"spooky/internal/schemas"
)

// Manager implements the ActionManager interface
type Manager struct {
	// Configuration
	defaultTimeout   time.Duration
	defaultParallel  bool
	customValidators map[string]ActionValidator

	// State
	actions map[string]*types.Action
	logger  logging.Logger
	mu      sync.RWMutex
}

// NewManager creates a new ActionManager
func NewManager(logger logging.Logger) *Manager {
	return &Manager{
		defaultTimeout:   30 * time.Minute,
		defaultParallel:  false,
		customValidators: make(map[string]ActionValidator),
		actions:          make(map[string]*types.Action),
		logger:           logger,
	}
}

// LoadActions loads actions from the project
func (m *Manager) LoadActions(projectPath string) (*types.ActionCollection, error) {
	m.logger.Info("Loading actions from project", logging.String("project", projectPath))

	// For now, return an empty collection
	// In a real implementation, this would load actions from actions.hcl files
	collection := &types.ActionCollection{
		Actions:   make([]*types.Action, 0),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Metadata:  make(map[string]interface{}),
	}

	m.logger.Info("Loaded actions from project",
		logging.String("project", projectPath),
		logging.Int("actions_count", 0))

	return collection, nil
}

// GetAction gets an action by name
func (m *Manager) GetAction(name string) (*types.Action, error) {
	if name == "" {
		return nil, fmt.Errorf("action name cannot be empty")
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	action, exists := m.actions[name]
	if !exists {
		return nil, fmt.Errorf("action '%s' not found", name)
	}

	m.logger.Debug("Retrieved action", logging.String("action", name))
	return action, nil
}

// ListActions lists all available actions
func (m *Manager) ListActions() ([]*types.Action, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	actions := make([]*types.Action, 0, len(m.actions))
	for _, action := range m.actions {
		actions = append(actions, action)
	}

	m.logger.Debug("Listed actions", logging.Int("count", len(actions)))
	return actions, nil
}

// AddAction adds a new action
func (m *Manager) AddAction(name string, action *types.Action) error {
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
	m.logger.Info("Added action", logging.String("action", name))

	return nil
}

// RemoveAction removes an action
func (m *Manager) RemoveAction(name string) error {
	if name == "" {
		return fmt.Errorf("action name cannot be empty")
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// Check if action exists
	if _, exists := m.actions[name]; !exists {
		return fmt.Errorf("action '%s' not found", name)
	}

	delete(m.actions, name)
	m.logger.Info("Removed action", logging.String("action", name))

	return nil
}

// ExecuteAction executes a single action
func (m *Manager) ExecuteAction(ctx context.Context, action *types.Action, context *types.ActionContext) (*types.ActingSession, error) {
	if action == nil {
		return nil, fmt.Errorf("action cannot be nil")
	}

	if context == nil {
		return nil, fmt.Errorf("context cannot be nil")
	}

	m.logger.Info("Executing action", logging.String("action", action.Name))

	now := time.Now()
	session := &types.ActingSession{
		ActionName: action.Name,
		Status:     types.ActingStatusRunning,
		StartTime:  &now,
		CreatedAt:  now,
		UpdatedAt:  now,
		Results:    make([]*types.ActingResult, 0),
		Metadata:   make(map[string]interface{}),
	}

	// Validate action before execution
	if err := m.ValidateAction(action); err != nil {
		session.Status = types.ActingStatusFailed
		session.Error = err
		return session, fmt.Errorf("action validation failed: %w", err)
	}

	// Prepare action for execution
	if err := m.PrepareAction(action, context); err != nil {
		session.Status = types.ActingStatusFailed
		session.Error = err
		return session, fmt.Errorf("action preparation failed: %w", err)
	}

	// Execute action based on type
	var result *types.ActingResult
	var err error

	switch action.Type {
	case "command":
		result, err = m.executeCommandAction(ctx, action, context)
	case "script":
		result, err = m.executeScriptAction(ctx, action, context)
	case "template_deploy", "template_evaluate", "template_validate", "template_cleanup":
		result, err = m.executeTemplateAction(ctx, action, context)
	default:
		err = fmt.Errorf("unsupported action type: %s", action.Type)
	}

	if err != nil {
		session.Status = types.ActingStatusFailed
		session.Error = err
	} else {
		session.Status = types.ActingStatusCompleted
		session.Results = append(session.Results, result)
	}

	endTime := time.Now()
	session.EndTime = &endTime
	session.Duration = endTime.Sub(now)
	session.UpdatedAt = endTime

	return session, nil
}

// ExecuteActionCollection executes a collection of actions
func (m *Manager) ExecuteActionCollection(ctx context.Context, collection *types.ActionCollection, context *types.ActionContext) (*types.ActingSession, error) {
	if collection == nil {
		return nil, fmt.Errorf("collection cannot be nil")
	}

	if context == nil {
		return nil, fmt.Errorf("context cannot be nil")
	}

	// TODO: Implement actual collection execution
	m.logger.Info("Executing action collection", logging.Int("actions_count", len(collection.Actions)))

	now := time.Now()
	session := &types.ActingSession{
		ActionName: "collection",
		Status:     types.ActingStatusRunning,
		StartTime:  &now,
		CreatedAt:  now,
		UpdatedAt:  now,
		Results:    make([]*types.ActingResult, 0),
		Metadata:   make(map[string]interface{}),
	}

	return session, nil
}

// executeCommandAction executes a command action
func (m *Manager) executeCommandAction(ctx context.Context, action *types.Action, context *types.ActionContext) (*types.ActingResult, error) {
	result := &types.ActingResult{
		ActionName: action.Name,
		Status:     types.ActingStatusRunning,
		StartTime:  &time.Time{},
		CreatedAt:  time.Now(),
	}

	// Execute command on target machines
	if len(context.Machines) == 0 {
		result.Status = types.ActingStatusFailed
		result.Error = "no machines specified for command execution"
		return result, fmt.Errorf("no machines specified for command execution")
	}

	// For now, return a placeholder result
	// In a real implementation, this would execute the command on each machine via SSH
	result.Status = types.ActingStatusCompleted
	result.Output = fmt.Sprintf("Command '%s' executed successfully on %d machines", action.Command, len(context.Machines))
	result.ExitCode = 0

	endTime := time.Now()
	result.EndTime = &endTime
	result.Duration = endTime.Sub(*result.StartTime)

	return result, nil
}

// executeScriptAction executes a script action
func (m *Manager) executeScriptAction(ctx context.Context, action *types.Action, context *types.ActionContext) (*types.ActingResult, error) {
	result := &types.ActingResult{
		ActionName: action.Name,
		Status:     types.ActingStatusRunning,
		StartTime:  &time.Time{},
		CreatedAt:  time.Now(),
	}

	// Execute script on target machines
	if len(context.Machines) == 0 {
		result.Status = types.ActingStatusFailed
		result.Error = "no machines specified for script execution"
		return result, fmt.Errorf("no machines specified for script execution")
	}

	// For now, return a placeholder result
	// In a real implementation, this would upload and execute the script on each machine via SSH
	result.Status = types.ActingStatusCompleted
	result.Output = fmt.Sprintf("Script '%s' executed successfully on %d machines", action.Script, len(context.Machines))
	result.ExitCode = 0

	endTime := time.Now()
	result.EndTime = &endTime
	result.Duration = endTime.Sub(*result.StartTime)

	return result, nil
}

// executeTemplateAction executes a template action
func (m *Manager) executeTemplateAction(ctx context.Context, action *types.Action, context *types.ActionContext) (*types.ActingResult, error) {
	result := &types.ActingResult{
		ActionName: action.Name,
		Status:     types.ActingStatusRunning,
		StartTime:  &time.Time{},
		CreatedAt:  time.Now(),
	}

	// Execute template action
	if action.Template == nil {
		result.Status = types.ActingStatusFailed
		result.Error = "no template configuration specified"
		return result, fmt.Errorf("no template configuration specified")
	}

	// For now, return a placeholder result
	// In a real implementation, this would render and deploy the template
	result.Status = types.ActingStatusCompleted
	result.Output = fmt.Sprintf("Template '%s' processed successfully", action.Template.Source)
	result.ExitCode = 0

	endTime := time.Now()
	result.EndTime = &endTime
	result.Duration = endTime.Sub(*result.StartTime)

	return result, nil
}

// PrepareAction prepares an action for execution
func (m *Manager) PrepareAction(action *types.Action, context *types.ActionContext) error {
	if action == nil {
		return fmt.Errorf("action cannot be nil")
	}

	if context == nil {
		return fmt.Errorf("context cannot be nil")
	}

	// TODO: Implement actual action preparation
	m.logger.Info("Preparing action", logging.String("action", action.Name))
	return nil
}

// PlanAction creates a plan for executing an action
func (m *Manager) PlanAction(action *types.Action, context *types.ActionContext) (*types.ActionPlan, error) {
	if action == nil {
		return nil, fmt.Errorf("action cannot be nil")
	}

	if context == nil {
		return nil, fmt.Errorf("context cannot be nil")
	}

	// TODO: Implement actual action planning
	m.logger.Info("Planning action", logging.String("action", action.Name))

	now := time.Now()
	plan := &types.ActionPlan{
		PlanID:     "plan-" + action.Name,
		ActionName: action.Name,
		Status:     types.PlanningStatusPending,
		CreatedAt:  now,
		UpdatedAt:  now,
		Steps:      make([]*types.PlanStep, 0),
		Metadata:   make(map[string]interface{}),
	}

	return plan, nil
}

// PlanActionCollection creates a plan for executing a collection of actions
func (m *Manager) PlanActionCollection(collection *types.ActionCollection, context *types.ActionContext) (*types.ActionPlan, error) {
	if collection == nil {
		return nil, fmt.Errorf("collection cannot be nil")
	}

	if context == nil {
		return nil, fmt.Errorf("context cannot be nil")
	}

	// TODO: Implement actual collection planning
	m.logger.Info("Planning action collection", logging.Int("actions_count", len(collection.Actions)))

	now := time.Now()
	plan := &types.ActionPlan{
		PlanID:     "plan-collection",
		ActionName: "collection",
		Status:     types.PlanningStatusPending,
		CreatedAt:  now,
		UpdatedAt:  now,
		Steps:      make([]*types.PlanStep, 0),
		Metadata:   make(map[string]interface{}),
	}

	return plan, nil
}

// ValidatePlan validates an action plan
func (m *Manager) ValidatePlan(plan *types.ActionPlan) error {
	if plan == nil {
		return fmt.Errorf("plan cannot be nil")
	}

	// TODO: Implement actual plan validation
	m.logger.Info("Validating plan", logging.String("plan", plan.PlanID))
	return nil
}

// ValidateAction validates an action
func (m *Manager) ValidateAction(action *types.Action) error {
	if action == nil {
		return fmt.Errorf("action cannot be nil")
	}

	// TODO: Implement actual action validation
	m.logger.Info("Validating action", logging.String("action", action.Name))
	return nil
}

// ValidateActionCollection validates a collection of actions
func (m *Manager) ValidateActionCollection(collection *types.ActionCollection) error {
	if collection == nil {
		return fmt.Errorf("collection cannot be nil")
	}

	// TODO: Implement actual collection validation
	m.logger.Info("Validating action collection", logging.Int("actions_count", len(collection.Actions)))
	return nil
}

// ValidateActionContext validates an action context
func (m *Manager) ValidateActionContext(context *types.ActionContext) error {
	if context == nil {
		return fmt.Errorf("context cannot be nil")
	}

	// TODO: Implement actual context validation
	m.logger.Info("Validating action context", logging.String("project", context.ProjectPath))
	return nil
}

// MergeActions merges multiple actions into a collection
func (m *Manager) MergeActions(actions ...*types.Action) (*types.ActionCollection, error) {
	if len(actions) == 0 {
		return nil, fmt.Errorf("at least one action must be provided")
	}

	// TODO: Implement actual action merging
	m.logger.Info("Merging actions", logging.Int("actions_count", len(actions)))

	now := time.Now()
	collection := &types.ActionCollection{
		Actions:   actions,
		CreatedAt: now,
		UpdatedAt: now,
		Metadata:  make(map[string]interface{}),
	}

	return collection, nil
}

// MergeWithPolicy merges actions with a specific policy
func (m *Manager) MergeWithPolicy(existing, new *types.ActionCollection, policy types.MergePolicy) (*types.ActionCollection, error) {
	if existing == nil {
		return nil, fmt.Errorf("existing collection cannot be nil")
	}

	if new == nil {
		return nil, fmt.Errorf("new collection cannot be nil")
	}

	// TODO: Implement actual policy-based merging
	m.logger.Info("Merging collections with policy", logging.String("policy", policy.PolicyName))

	now := time.Now()
	merged := &types.ActionCollection{
		Actions:   append(existing.Actions, new.Actions...),
		CreatedAt: now,
		UpdatedAt: now,
		Metadata:  make(map[string]interface{}),
	}

	return merged, nil
}

// OptimizeAction optimizes an action for performance
func (m *Manager) OptimizeAction(action *types.Action) error {
	if action == nil {
		return fmt.Errorf("action cannot be nil")
	}

	// TODO: Implement actual action optimization
	m.logger.Info("Optimizing action", logging.String("action", action.Name))
	return nil
}

// OptimizeActionCollection optimizes a collection of actions
func (m *Manager) OptimizeActionCollection(collection *types.ActionCollection) error {
	if collection == nil {
		return fmt.Errorf("collection cannot be nil")
	}

	// TODO: Implement actual collection optimization
	m.logger.Info("Optimizing action collection", logging.Int("actions_count", len(collection.Actions)))
	return nil
}

// GetPerformanceMetrics gets performance metrics for an action
func (m *Manager) GetPerformanceMetrics(action *types.Action) (*types.PerformanceMetrics, error) {
	if action == nil {
		return nil, fmt.Errorf("action cannot be nil")
	}

	// TODO: Implement actual performance metrics collection
	m.logger.Info("Getting performance metrics", logging.String("action", action.Name))

	now := time.Now()
	metrics := &types.PerformanceMetrics{
		ActionName: action.Name,
		CreatedAt:  now,
		UpdatedAt:  now,
		Metadata:   make(map[string]interface{}),
	}

	return metrics, nil
}

// SetDefaultTimeout sets the default timeout
func (m *Manager) SetDefaultTimeout(timeout time.Duration) {
	m.defaultTimeout = timeout
}

// SetDefaultParallel sets the default parallel flag
func (m *Manager) SetDefaultParallel(parallel bool) {
	m.defaultParallel = parallel
}

// RegisterCustomValidator registers a custom validator
func (m *Manager) RegisterCustomValidator(name string, validator ActionValidator) {
	if name == "" {
		m.logger.Warn("Cannot register validator with empty name")
		return
	}

	if validator == nil {
		m.logger.Warn("Cannot register nil validator", logging.String("name", name))
		return
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	m.customValidators[name] = validator
	m.logger.Info("Registered custom validator", logging.String("name", name))
}

// Close closes the manager
func (m *Manager) Close() error {
	m.logger.Info("Closing action manager")
	return nil
}

// validateActions validates actions using schema system
func (m *Manager) validateActions(content []byte) error {
	validator := schemas.NewSchemaValidator()
	if err := validator.LoadSchema(schemas.SchemaTypeActions); err != nil {
		return fmt.Errorf("failed to load actions schema: %w", err)
	}

	// Parse content to interface{} for validation
	var data interface{}
	if err := json.Unmarshal(content, &data); err != nil {
		return fmt.Errorf("failed to parse content for validation: %w", err)
	}

	if err := validator.ValidateData(data, "actions"); err != nil {
		return fmt.Errorf("actions validation failed: %w", err)
	}
	return nil
}

// validateActionsComposed validates complex actions using schema system
func (m *Manager) validateActionsComposed(content []byte) error {
	validator := schemas.NewSchemaValidator()
	if err := validator.LoadSchema(schemas.SchemaTypeActionsComposed); err != nil {
		return fmt.Errorf("failed to load composed actions schema: %w", err)
	}

	// Parse content to interface{} for validation
	var data interface{}
	if err := json.Unmarshal(content, &data); err != nil {
		return fmt.Errorf("failed to parse content for validation: %w", err)
	}

	if err := validator.ValidateData(data, "actions-composed"); err != nil {
		return fmt.Errorf("composed actions validation failed: %w", err)
	}
	return nil
}
