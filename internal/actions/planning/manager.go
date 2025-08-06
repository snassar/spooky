package planning

import (
	"fmt"
	"sync"
	"time"

	"spooky/internal/actions/types"
	"spooky/internal/logging"
)

// Manager implements the PlanningManager interface
type Manager struct {
	// Configuration
	defaultStrategy     types.PlanningStrategy
	defaultOptimization types.PlanningOptimization

	// State
	planners map[string]Planner
	plans    map[string]*types.ActionPlan
	logger   logging.Logger
	mu       sync.RWMutex
}

// NewManager creates a new PlanningManager
func NewManager(logger logging.Logger) *Manager {
	return &Manager{
		defaultStrategy:     types.PlanningStrategySequential,
		defaultOptimization: types.PlanningOptimizationNone,
		planners:            make(map[string]Planner),
		plans:               make(map[string]*types.ActionPlan),
		logger:              logger,
	}
}

// PlanAction plans a single action
func (m *Manager) PlanAction(action *types.Action, context *types.ActionContext) (*types.ActionPlan, error) {
	if action == nil {
		return nil, fmt.Errorf("action cannot be nil")
	}

	m.logger.Info("Planning action", logging.String("action", action.Name))

	planner, err := m.CreatePlanner(action)
	if err != nil {
		return nil, fmt.Errorf("failed to create planner for action %s: %w", action.Name, err)
	}

	plan, err := planner.Plan(context)
	if err != nil {
		return nil, fmt.Errorf("failed to plan action %s: %w", action.Name, err)
	}

	// Store the plan
	m.mu.Lock()
	m.plans[plan.PlanID] = plan
	m.mu.Unlock()

	m.logger.Info("Successfully planned action",
		logging.String("action", action.Name),
		logging.String("plan_id", plan.PlanID),
		logging.Int("steps", len(plan.Steps)))

	return plan, nil
}

// PlanActionCollection plans a collection of actions
func (m *Manager) PlanActionCollection(collection *types.ActionCollection, context *types.ActionContext) (*types.ActionPlan, error) {
	if collection == nil {
		return nil, fmt.Errorf("action collection cannot be nil")
	}

	m.logger.Info("Planning action collection", logging.Int("actions_count", len(collection.Actions)))

	// Create a combined plan
	planID := generatePlanID()
	plan := &types.ActionPlan{
		PlanID:       planID,
		ActionName:   "collection",
		Steps:        make([]*types.PlanStep, 0),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		Strategy:     m.defaultStrategy,
		Optimization: m.defaultOptimization,
		Metadata:     make(map[string]interface{}),
	}

	// Plan each action and combine steps
	for _, action := range collection.Actions {
		actionPlan, err := m.PlanAction(action, context)
		if err != nil {
			return nil, fmt.Errorf("failed to plan action %s in collection: %w", action.Name, err)
		}

		// Add steps from this action's plan
		for _, step := range actionPlan.Steps {
			plan.Steps = append(plan.Steps, step)
		}
	}

	// Store the combined plan
	m.mu.Lock()
	m.plans[plan.PlanID] = plan
	m.mu.Unlock()

	m.logger.Info("Successfully planned action collection",
		logging.String("plan_id", plan.PlanID),
		logging.Int("total_steps", len(plan.Steps)))

	return plan, nil
}

// ValidatePlan validates an action plan
func (m *Manager) ValidatePlan(plan *types.ActionPlan) error {
	if plan == nil {
		return fmt.Errorf("plan cannot be nil")
	}

	m.logger.Info("Validating plan", logging.String("plan_id", plan.PlanID))

	// Create a validator for this plan
	validator := NewPlanValidator(m.logger)

	// Validate the plan
	if err := validator.ValidatePlan(plan); err != nil {
		return fmt.Errorf("plan validation failed: %w", err)
	}

	m.logger.Info("Plan validation successful", logging.String("plan_id", plan.PlanID))
	return nil
}

// CreatePlanner creates a new planner for an action
func (m *Manager) CreatePlanner(action *types.Action) (Planner, error) {
	if action == nil {
		return nil, fmt.Errorf("action cannot be nil")
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// Check if planner already exists
	if planner, exists := m.planners[action.Name]; exists {
		return planner, nil
	}

	// Create new planner
	planner := NewPlanner(action, m.logger)
	m.planners[action.Name] = planner

	m.logger.Debug("Created planner for action", logging.String("action", action.Name))
	return planner, nil
}

// GetPlanner gets an existing planner for an action
func (m *Manager) GetPlanner(action *types.Action) (Planner, error) {
	if action == nil {
		return nil, fmt.Errorf("action cannot be nil")
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	planner, exists := m.planners[action.Name]
	if !exists {
		return nil, fmt.Errorf("planner not found for action %s", action.Name)
	}

	return planner, nil
}

// GetPlan gets a plan by ID
func (m *Manager) GetPlan(planID string) (*types.ActionPlan, error) {
	if planID == "" {
		return nil, fmt.Errorf("plan ID cannot be empty")
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	plan, exists := m.plans[planID]
	if !exists {
		return nil, fmt.Errorf("plan not found: %s", planID)
	}

	return plan, nil
}

// ListPlans lists all stored plans
func (m *Manager) ListPlans() ([]*types.ActionPlan, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	plans := make([]*types.ActionPlan, 0, len(m.plans))
	for _, plan := range m.plans {
		plans = append(plans, plan)
	}

	return plans, nil
}

// DeletePlan deletes a plan by ID
func (m *Manager) DeletePlan(planID string) error {
	if planID == "" {
		return fmt.Errorf("plan ID cannot be empty")
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.plans[planID]; !exists {
		return fmt.Errorf("plan not found: %s", planID)
	}

	delete(m.plans, planID)
	m.logger.Info("Deleted plan", logging.String("plan_id", planID))

	return nil
}

// SetDefaultStrategy sets the default planning strategy
func (m *Manager) SetDefaultStrategy(strategy types.PlanningStrategy) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.defaultStrategy = strategy
	m.logger.Info("Set default planning strategy", logging.String("strategy", string(strategy)))
}

// SetDefaultOptimization sets the default planning optimization
func (m *Manager) SetDefaultOptimization(optimization types.PlanningOptimization) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.defaultOptimization = optimization
	m.logger.Info("Set default planning optimization", logging.String("optimization", string(optimization)))
}

// generatePlanID generates a unique plan ID
func generatePlanID() string {
	return fmt.Sprintf("plan_%d", time.Now().UnixNano())
}
