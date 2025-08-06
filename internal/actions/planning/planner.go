package planning

import (
	"fmt"
	"time"

	"spooky/internal/actions/types"
	"spooky/internal/logging"
)

// ActionPlanner implements the Planner interface
type ActionPlanner struct {
	action       *types.Action
	strategy     types.PlanningStrategy
	optimization types.PlanningOptimization
	constraints  []types.PlanningConstraint
	logger       logging.Logger
}

// NewPlanner creates a new Planner
func NewPlanner(action *types.Action, logger logging.Logger) Planner {
	return &ActionPlanner{
		action:       action,
		strategy:     types.PlanningStrategySequential,
		optimization: types.PlanningOptimizationNone,
		constraints:  make([]types.PlanningConstraint, 0),
		logger:       logger,
	}
}

// Plan creates a plan for the action
func (p *ActionPlanner) Plan(context *types.ActionContext) (*types.ActionPlan, error) {
	if p.action == nil {
		return nil, fmt.Errorf("action cannot be nil")
	}

	p.logger.Info("Creating plan for action", logging.String("action", p.action.Name))

	// Create a new plan
	plan := &types.ActionPlan{
		PlanID:       generatePlanID(),
		ActionName:   p.action.Name,
		Status:       types.PlanningStatusPending,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		Strategy:     p.strategy,
		Optimization: p.optimization,
		Steps:        make([]*types.PlanStep, 0),
		Constraints:  p.constraints,
		Metadata:     make(map[string]interface{}),
	}

	// Create steps based on the action
	step := &types.PlanStep{
		StepID:    generateStepID(),
		StepName:  p.action.Name,
		StepOrder: 1,
		Action:    p.action,
		Status:    types.PlanStepStatusPending,
		Metadata:  make(map[string]interface{}),
	}

	plan.Steps = append(plan.Steps, step)

	p.logger.Info("Successfully created plan",
		logging.String("action", p.action.Name),
		logging.String("plan_id", plan.PlanID),
		logging.Int("steps", len(plan.Steps)))

	return plan, nil
}

// Validate validates a plan
func (p *ActionPlanner) Validate(plan *types.ActionPlan) error {
	if plan == nil {
		return fmt.Errorf("plan cannot be nil")
	}

	p.logger.Info("Validating plan", logging.String("plan_id", plan.PlanID))

	// Basic validation
	if plan.PlanID == "" {
		return fmt.Errorf("plan ID cannot be empty")
	}

	if len(plan.Steps) == 0 {
		return fmt.Errorf("plan must have at least one step")
	}

	// Validate each step
	for i, step := range plan.Steps {
		if step.StepID == "" {
			return fmt.Errorf("step %d: step ID cannot be empty", i)
		}

		if step.Action == nil {
			return fmt.Errorf("step %d: action cannot be nil", i)
		}
	}

	p.logger.Info("Plan validation successful", logging.String("plan_id", plan.PlanID))
	return nil
}

// Optimize optimizes a plan
func (p *ActionPlanner) Optimize(plan *types.ActionPlan) error {
	if plan == nil {
		return fmt.Errorf("plan cannot be nil")
	}

	p.logger.Info("Optimizing plan", logging.String("plan_id", plan.PlanID))

	// Apply optimization based on strategy
	switch p.optimization {
	case types.PlanningOptimizationBasic:
		// Basic optimization: ensure proper step ordering
		for i, step := range plan.Steps {
			step.StepOrder = i + 1
		}

	case types.PlanningOptimizationAdvanced:
		// Advanced optimization: parallel execution where possible
		for _, step := range plan.Steps {
			if step.Action != nil && step.Action.Parallel {
				step.Parallel = true
			}
		}

	case types.PlanningOptimizationMaximum:
		// Maximum optimization: aggressive parallelization and resource optimization
		for _, step := range plan.Steps {
			step.Parallel = true
			if step.Timeout == 0 {
				step.Timeout = 30 * time.Minute
			}
		}
	}

	p.logger.Info("Plan optimization completed", logging.String("plan_id", plan.PlanID))
	return nil
}

// SetStrategy sets the planning strategy
func (p *ActionPlanner) SetStrategy(strategy types.PlanningStrategy) {
	p.strategy = strategy
	p.logger.Debug("Set planning strategy", logging.String("strategy", string(strategy)))
}

// GetStrategy gets the current planning strategy
func (p *ActionPlanner) GetStrategy() types.PlanningStrategy {
	return p.strategy
}

// SetOptimization sets the optimization level
func (p *ActionPlanner) SetOptimization(optimization types.PlanningOptimization) {
	p.optimization = optimization
	p.logger.Debug("Set optimization level", logging.String("optimization", string(optimization)))
}

// SetConstraints sets the planning constraints
func (p *ActionPlanner) SetConstraints(constraints []types.PlanningConstraint) {
	p.constraints = constraints
	p.logger.Debug("Set planning constraints", logging.Int("constraints_count", len(constraints)))
}

// generateStepID generates a unique step ID
func generateStepID() string {
	return fmt.Sprintf("step_%d", time.Now().UnixNano())
}
