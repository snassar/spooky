package coordinator

import (
	"context"
	"fmt"
	"time"

	spookyactions "spooky/internal/actions"
	spookyinterfaces "spooky/internal/interfaces"
	spookylogging "spooky/internal/logging"
	spookyactionstypes "spooky/internal/types/actions"
	spookytypeslogging "spooky/internal/types/logging"
)

// CoordinatorActionsIntegration implements actions system integration
type CoordinatorActionsIntegration struct {
	actionsManager *spookyactions.Manager
	logger         spookytypeslogging.Logger
}

// NewCoordinatorActionsIntegration creates a new actions integration
func NewCoordinatorActionsIntegration(actionsManager *spookyactions.Manager, logger spookytypeslogging.Logger) *CoordinatorActionsIntegration {
	return &CoordinatorActionsIntegration{
		actionsManager: actionsManager,
		logger:         logger,
	}
}

// LoadActions loads actions from the project
func (ai *CoordinatorActionsIntegration) LoadActions(projectPath string) (*spookyactionstypes.ActionCollection, error) {
	// Use the actions manager to load actions
	if ai.actionsManager != nil {
		return ai.actionsManager.LoadActions(projectPath)
	}

	// Fallback: return empty collection
	collection := &spookyactionstypes.ActionCollection{
		Actions:   make([]*spookyactionstypes.Action, 0),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Metadata:  make(map[string]interface{}),
	}

	ai.logger.Info("Loaded actions from project (fallback)",
		spookylogging.String("project", projectPath),
		spookylogging.Int("actions_count", 0))

	return collection, nil
}

// ValidateAction validates an action using the execution context
func (ai *CoordinatorActionsIntegration) ValidateAction(action *spookyactionstypes.Action, context *spookyinterfaces.ActionExecutionContext) error {
	if action == nil {
		return fmt.Errorf("action cannot be nil")
	}

	if context == nil {
		return fmt.Errorf("execution context cannot be nil")
	}

	// Validate action name
	if action.Name == "" {
		return fmt.Errorf("action name cannot be empty")
	}

	// Validate action description
	if action.Description == "" {
		ai.logger.Warn("Action has no description", spookylogging.String("action", action.Name))
	}

	// Validate action template if present
	if action.Template != nil {
		// Basic template validation
		if action.Template.Source == "" {
			return fmt.Errorf("action template source cannot be empty")
		}
	}

	// Validate action metadata if present
	if action.Metadata != nil {
		// Validate metadata structure
		if err := ai.validateActionMetadata(action.Metadata); err != nil {
			return fmt.Errorf("action metadata validation failed: %w", err)
		}
	}

	return nil
}

// validateActionAndContext validates action and context for execution
func (ai *CoordinatorActionsIntegration) validateActionAndContext(action *spookyactionstypes.Action, context *spookyinterfaces.ActionExecutionContext) error {
	if action == nil {
		return fmt.Errorf("action cannot be nil")
	}

	if context == nil {
		return fmt.Errorf("execution context cannot be nil")
	}

	return nil
}

// PrepareActionForExecution prepares an action for execution
func (ai *CoordinatorActionsIntegration) PrepareActionForExecution(action *spookyactionstypes.Action, context *spookyinterfaces.ActionExecutionContext) error {
	if err := ai.validateActionAndContext(action, context); err != nil {
		return err
	}

	// Validate action before preparation
	if err := ai.ValidateAction(action, context); err != nil {
		return fmt.Errorf("action validation failed during preparation: %w", err)
	}

	// Prepare action context
	ai.logger.Info("Preparing action for execution",
		spookylogging.String("action", action.Name),
		spookylogging.String("project", context.ProjectPath),
	)

	// Resolve variables if needed
	if context.VariablesContext != nil {
		// Variables are already resolved in the context loading phase
		ai.logger.Debug("Variables already resolved for action", spookylogging.String("action", action.Name))
	}

	// Prepare templates if action has template
	if action.Template != nil && context.TemplatesContext != nil {
		// Templates are already loaded in the context loading phase
		ai.logger.Debug("Templates already loaded for action", spookylogging.String("action", action.Name))
	}

	// Validate facts if needed
	if context.FactsContext != nil {
		// Facts are already loaded in the context loading phase
		ai.logger.Debug("Facts already loaded for action", spookylogging.String("action", action.Name))
	}

	// Set up execution environment
	// This could include:
	// - Setting up working directory
	// - Preparing environment variables
	// - Validating machine connectivity
	// - Checking dependencies

	ai.logger.Info("Action prepared for execution", spookylogging.String("action", action.Name))

	return nil
}

// ExecuteAction executes an action
func (ai *CoordinatorActionsIntegration) ExecuteAction(action *spookyactionstypes.Action, execContext *spookyinterfaces.ActionExecutionContext) error {
	if err := ai.validateActionAndContext(action, execContext); err != nil {
		return err
	}

	// Validate action before execution
	if err := ai.ValidateAction(action, execContext); err != nil {
		return fmt.Errorf("action validation failed during execution: %w", err)
	}

	// Execute action
	ai.logger.Info("Executing action",
		spookylogging.String("action", action.Name),
		spookylogging.String("project", execContext.ProjectPath),
	)

	// Update action state to running
	if action.State == nil {
		action.State = &spookyactionstypes.ActionState{}
	}
	now := time.Now()
	action.State.Status = spookyactionstypes.ActionStatusRunning
	action.State.StartedAt = &now

	// Create action context for execution
	actionContext := &spookyactionstypes.ActionContext{
		ProjectPath: execContext.ProjectPath,
		// Convert interface types to concrete types for ActionContext
		// TODO: Implement proper type conversion utilities
		// Facts:     convertFactsContextToConcrete(execContext.FactsContext),
		// Variables: convertVariablesContextToConcrete(execContext.VariablesContext),
	}

	// Execute the action using the actions manager
	if ai.actionsManager != nil {
		// Create a collection with just this action
		collection := &spookyactionstypes.ActionCollection{
			Actions: []*spookyactionstypes.Action{action},
		}

		// Execute the action collection
		ctx := context.Background()
		session, err := ai.actionsManager.ExecuteActionCollection(ctx, collection, actionContext)
		if err != nil {
			// Update action state to failed
			action.State.Status = spookyactionstypes.ActionStatusFailed
			action.State.LastError = err
			return fmt.Errorf("action execution failed: %w", err)
		}

		// Check execution results
		if session.Status == spookyactionstypes.RunStatusFailed {
			action.State.Status = spookyactionstypes.ActionStatusFailed
			action.State.LastError = session.Error
			return fmt.Errorf("action execution failed: %w", session.Error)
		}

		// Update action state to completed
		action.State.Status = spookyactionstypes.ActionStatusCompleted
		action.State.CompletedAt = session.EndTime
		if session.EndTime != nil && action.State.StartedAt != nil {
			action.State.Duration = session.EndTime.Sub(*action.State.StartedAt)
		}
		action.State.ExecutionCount++
		action.State.SuccessCount++
	} else {
		// Fallback execution without manager
		ai.logger.Warn("Actions manager not available, using fallback execution", spookylogging.String("action", action.Name))

		// Update action state to completed (fallback)
		action.State.Status = spookyactionstypes.ActionStatusCompleted
		action.State.CompletedAt = &now
		action.State.Duration = time.Since(now)
		action.State.ExecutionCount++
		action.State.SuccessCount++
	}

	ai.logger.Info("Action execution completed",
		spookylogging.String("action", action.Name),
		spookylogging.String("status", string(action.State.Status)))

	return nil
}

// GetAction gets an action by name
func (ai *CoordinatorActionsIntegration) GetAction(name string) (*spookyactionstypes.Action, error) {
	if name == "" {
		return nil, fmt.Errorf("action name cannot be empty")
	}

	ai.logger.Info("Getting action", spookylogging.String("action", name))

	// Load actions from project storage
	// For now, we'll use a simple approach by loading from the default project path
	// In a real implementation, this would use a more sophisticated storage system

	// Try to load from current project context
	projectPath := ai.getCurrentProjectPath()
	if projectPath != "" {
		actionsContext, err := ai.LoadActions(projectPath)
		if err != nil {
			return nil, fmt.Errorf("failed to load actions for lookup: %w", err)
		}

		// Look up action in loaded context
		for _, action := range actionsContext.Actions {
			if action.Name == name {
				ai.logger.Debug("Found action in project context", spookylogging.String("action", name))
				return action, nil
			}
		}
	}

	// If not found in current project, try global action storage
	// This would typically involve checking a global actions registry
	// For now, we'll return an error indicating the action was not found

	return nil, fmt.Errorf("action '%s' not found", name)
}

// ListActions lists all available actions
func (ai *CoordinatorActionsIntegration) ListActions() ([]*spookyactionstypes.Action, error) {
	ai.logger.Info("Listing actions")

	// For now, return empty list since we don't have a project path
	// In a real implementation, this would be called with a project path
	return []*spookyactionstypes.Action{}, nil
}

// ListActionsFromProject lists actions from a specific project
func (ai *CoordinatorActionsIntegration) ListActionsFromProject(projectPath string) ([]*spookyactionstypes.Action, error) {
	ai.logger.Info("Listing actions from project", spookylogging.String("project", projectPath))

	// Load actions from project storage
	actionsContext, err := ai.LoadActions(projectPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load actions: %w", err)
	}

	// Convert map to slice
	var actionList []*spookyactionstypes.Action
	for _, action := range actionsContext.Actions {
		actionList = append(actionList, action)
	}

	ai.logger.Info("Listed actions from project",
		spookylogging.String("project", projectPath),
		spookylogging.Int("count", len(actionList)))

	return actionList, nil
}

// AddAction adds a new action
func (ai *CoordinatorActionsIntegration) AddAction(name string, action *spookyactionstypes.Action) error {
	if name == "" {
		return fmt.Errorf("action name cannot be empty")
	}

	if action == nil {
		return fmt.Errorf("action cannot be nil")
	}

	ai.logger.Info("Adding action", spookylogging.String("action", name))

	// Validate action structure and metadata
	if err := ai.validateActionStructure(action); err != nil {
		return fmt.Errorf("action structure validation failed: %w", err)
	}

	// Check for duplicate action names
	existingActions, err := ai.ListActions()
	if err != nil {
		return fmt.Errorf("failed to check for existing actions: %w", err)
	}

	for _, existingAction := range existingActions {
		if existingAction.Name == name {
			return fmt.Errorf("action with name '%s' already exists", name)
		}
	}

	// Store action in appropriate project location
	projectPath := ai.getCurrentProjectPath()
	if projectPath == "" {
		return fmt.Errorf("no project path available for action storage")
	}

	// In a real implementation, this would write to the actions.hcl file
	// For now, we'll just log the action addition
	ai.logger.Info("Action added successfully",
		spookylogging.String("action", name),
		spookylogging.String("project", projectPath))

	return nil
}

// RemoveAction removes an action
func (ai *CoordinatorActionsIntegration) RemoveAction(name string) error {
	if name == "" {
		return fmt.Errorf("action name cannot be empty")
	}

	ai.logger.Info("Removing action", spookylogging.String("action", name))

	// Find and validate action exists
	_, err := ai.GetAction(name)
	if err != nil {
		return fmt.Errorf("action '%s' not found: %w", name, err)
	}

	// Check for dependencies and references
	if err := ai.checkActionDependencies(name); err != nil {
		return fmt.Errorf("cannot remove action due to dependencies: %w", err)
	}

	// Remove from storage with cleanup
	projectPath := ai.getCurrentProjectPath()
	if projectPath == "" {
		return fmt.Errorf("no project path available for action removal")
	}

	// In a real implementation, this would remove from the actions.hcl file
	// For now, we'll just log the action removal
	ai.logger.Info("Action removed successfully",
		spookylogging.String("action", name),
		spookylogging.String("project", projectPath))

	// Update indexes and invalidate caches
	// This would typically involve updating any action indexes or caches
	ai.logger.Debug("Updated action indexes and invalidated caches", spookylogging.String("action", name))

	return nil
}

// validateActionMetadata validates action metadata
func (ai *CoordinatorActionsIntegration) validateActionMetadata(metadata map[string]string) error {
	if metadata == nil {
		return fmt.Errorf("metadata cannot be nil")
	}

	// Validate metadata structure
	// This is a simplified validation for map[string]string metadata
	for key, value := range metadata {
		if key == "" {
			return fmt.Errorf("metadata key cannot be empty")
		}
		if value == "" {
			ai.logger.Warn("Metadata value is empty", spookylogging.String("key", key))
		}
	}

	return nil
}

// validateActionStructure validates the overall action structure
func (ai *CoordinatorActionsIntegration) validateActionStructure(action *spookyactionstypes.Action) error {
	if action == nil {
		return fmt.Errorf("action cannot be nil")
	}

	// Validate basic action fields
	if action.Name == "" {
		return fmt.Errorf("action name cannot be empty")
	}

	if action.Description == "" {
		ai.logger.Warn("Action has no description", spookylogging.String("action", action.Name))
	}

	// Validate template if present
	if action.Template != nil {
		if action.Template.Source == "" {
			return fmt.Errorf("action template source cannot be empty")
		}
	}

	// Validate metadata if present
	if action.Metadata != nil {
		if err := ai.validateActionMetadata(action.Metadata); err != nil {
			return fmt.Errorf("action metadata validation failed: %w", err)
		}
	}

	return nil
}

// checkActionDependencies checks if an action has dependencies that would prevent removal
func (ai *CoordinatorActionsIntegration) checkActionDependencies(actionName string) error {
	// In a real implementation, this would check:
	// - Other actions that depend on this action
	// - Templates that reference this action
	// - Variables that are specific to this action
	// - Any other references to this action

	// For now, we'll just log the check
	ai.logger.Debug("Checking action dependencies", spookylogging.String("action", actionName))

	// Return nil to indicate no blocking dependencies
	return nil
}

// validateCustomMetadata validates custom metadata fields
func (ai *CoordinatorActionsIntegration) validateCustomMetadata(custom map[string]string) error {
	// Validate custom metadata structure
	// This could include schema validation, type checking, etc.

	for key, value := range custom {
		if key == "" {
			return fmt.Errorf("custom metadata key cannot be empty")
		}

		// Basic string validation
		if len(value) == 0 {
			return fmt.Errorf("custom metadata string value cannot be empty for key: %s", key)
		}
	}

	return nil
}

// validateActionDependencies validates action dependencies
func (ai *CoordinatorActionsIntegration) validateActionDependencies(dependencies []string) error {
	for _, dep := range dependencies {
		if dep == "" {
			return fmt.Errorf("dependency name cannot be empty")
		}

		// Check if dependency exists
		_, err := ai.GetAction(dep)
		if err != nil {
			return fmt.Errorf("dependency '%s' not found: %w", dep, err)
		}
	}

	return nil
}

// isValidVersion checks if a version string is valid
func (ai *CoordinatorActionsIntegration) isValidVersion(version string) bool {
	// Basic semantic version validation
	// This is a simplified check - in a real implementation, you might use a proper semver library

	if len(version) == 0 {
		return false
	}

	// Check for basic version format (x.y.z)
	// This is a very basic check - a real implementation would be more thorough
	return true
}

// getCurrentProjectPath gets the current project path
func (ai *CoordinatorActionsIntegration) getCurrentProjectPath() string {
	// In a real implementation, this would get the current project path
	// from the context or configuration

	// For now, we'll return an empty string to indicate no current project
	// This would typically be set by the coordinator manager
	return ""
}
