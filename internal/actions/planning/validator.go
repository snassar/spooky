package planning

import (
	"fmt"

	spookyactionstypes "spooky/internal/actions/types"
	spookylogging "spooky/internal/logging"
)

// ActionPlanValidator implements the PlanValidator interface
type ActionPlanValidator struct {
	strictMode      bool
	validationLevel spookyactionstypes.ValidationLevel
	rules           []PlanValidationRule
	logger          spookylogging.Logger
}

// NewPlanValidator creates a new PlanValidator
func NewPlanValidator(logger spookylogging.Logger) PlanValidator {
	return &ActionPlanValidator{
		strictMode:      false,
		validationLevel: spookyactionstypes.ValidationLevelBasic,
		rules:           make([]PlanValidationRule, 0),
		logger:          logger,
	}
}

// ValidatePlan validates an action plan
func (v *ActionPlanValidator) ValidatePlan(plan *spookyactionstypes.ActionPlan) error {
	if plan == nil {
		return fmt.Errorf("plan cannot be nil")
	}

	v.logger.Info("Validating plan", spookylogging.String("plan_id", plan.PlanID))

	// Basic validation
	if err := v.validatePlanStructure(plan); err != nil {
		return fmt.Errorf("plan structure validation failed: %w", err)
	}

	// Validate each step
	for i, step := range plan.Steps {
		if err := v.ValidatePlanStep(step); err != nil {
			return fmt.Errorf("step %d validation failed: %w", i, err)
		}
	}

	// Validate dependencies
	if err := v.ValidatePlanDependencies(plan); err != nil {
		return fmt.Errorf("plan dependencies validation failed: %w", err)
	}

	// Run custom validation rules
	for _, rule := range v.rules {
		if err := rule.Validate(plan); err != nil {
			if v.strictMode {
				return fmt.Errorf("custom rule '%s' validation failed: %w", rule.Name(), err)
			}
			v.logger.Warn("Custom validation rule failed",
				spookylogging.String("rule", rule.Name()),
				spookylogging.String("error", err.Error()))
		}
	}

	v.logger.Info("Plan validation successful", spookylogging.String("plan_id", plan.PlanID))
	return nil
}

// ValidatePlanStep validates a plan step
func (v *ActionPlanValidator) ValidatePlanStep(step *spookyactionstypes.PlanStep) error {
	if step == nil {
		return fmt.Errorf("step cannot be nil")
	}

	// Basic step validation
	if step.StepID == "" {
		return fmt.Errorf("step ID cannot be empty")
	}

	if step.Action == nil {
		return fmt.Errorf("step action cannot be nil")
	}

	if step.Action.Name == "" {
		return fmt.Errorf("step action name cannot be empty")
	}

	// Run custom validation rules for steps
	for _, rule := range v.rules {
		if err := rule.ValidateStep(step); err != nil {
			if v.strictMode {
				return fmt.Errorf("custom rule '%s' step validation failed: %w", rule.Name(), err)
			}
			v.logger.Warn("Custom step validation rule failed",
				spookylogging.String("rule", rule.Name()),
				spookylogging.String("step_id", step.StepID),
				spookylogging.String("error", err.Error()))
		}
	}

	return nil
}

// ValidatePlanDependencies validates plan dependencies
func (v *ActionPlanValidator) ValidatePlanDependencies(plan *spookyactionstypes.ActionPlan) error {
	if plan == nil {
		return fmt.Errorf("plan cannot be nil")
	}

	// Check for circular dependencies
	visited := make(map[string]bool)
	recStack := make(map[string]bool)

	for _, step := range plan.Steps {
		if !visited[step.StepID] {
			if v.hasCircularDependency(step, plan.Steps, visited, recStack) {
				return fmt.Errorf("circular dependency detected in plan")
			}
		}
	}

	return nil
}

// AddValidationRule adds a custom validation rule
func (v *ActionPlanValidator) AddValidationRule(rule PlanValidationRule) error {
	if rule == nil {
		return fmt.Errorf("validation rule cannot be nil")
	}

	v.rules = append(v.rules, rule)
	v.logger.Debug("Added validation rule", spookylogging.String("rule", rule.Name()))
	return nil
}

// RemoveValidationRule removes a validation rule by name
func (v *ActionPlanValidator) RemoveValidationRule(name string) error {
	for i, rule := range v.rules {
		if rule.Name() == name {
			v.rules = append(v.rules[:i], v.rules[i+1:]...)
			v.logger.Debug("Removed validation rule", spookylogging.String("rule", name))
			return nil
		}
	}

	return fmt.Errorf("validation rule not found: %s", name)
}

// GetValidationRules returns all validation rules
func (v *ActionPlanValidator) GetValidationRules() ([]PlanValidationRule, error) {
	return v.rules, nil
}

// SetStrictMode sets the strict validation mode
func (v *ActionPlanValidator) SetStrictMode(strict bool) {
	v.strictMode = strict
	v.logger.Debug("Set strict validation mode", spookylogging.Bool("strict", strict))
}

// SetValidationLevel sets the validation level
func (v *ActionPlanValidator) SetValidationLevel(level spookyactionstypes.ValidationLevel) {
	v.validationLevel = level
	v.logger.Debug("Set validation level", spookylogging.String("level", string(level)))
}

// validatePlanStructure validates the basic structure of a plan
func (v *ActionPlanValidator) validatePlanStructure(plan *spookyactionstypes.ActionPlan) error {
	if plan.PlanID == "" {
		return fmt.Errorf("plan ID cannot be empty")
	}

	if len(plan.Steps) == 0 {
		return fmt.Errorf("plan must have at least one step")
	}

	return nil
}

// hasCircularDependency checks for circular dependencies using DFS
func (v *ActionPlanValidator) hasCircularDependency(step *spookyactionstypes.PlanStep, allSteps []*spookyactionstypes.PlanStep, visited, recStack map[string]bool) bool {
	visited[step.StepID] = true
	recStack[step.StepID] = true

	for _, depID := range step.Dependencies {
		// Find the dependent step
		var depStep *spookyactionstypes.PlanStep
		for _, s := range allSteps {
			if s.StepID == depID {
				depStep = s
				break
			}
		}

		if depStep == nil {
			continue // Skip if dependency not found
		}

		if !visited[depStep.StepID] {
			if v.hasCircularDependency(depStep, allSteps, visited, recStack) {
				return true
			}
		} else if recStack[depStep.StepID] {
			return true
		}
	}

	recStack[step.StepID] = false
	return false
}
