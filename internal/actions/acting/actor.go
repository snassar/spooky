package acting

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	spookytypes "spooky/internal/types"
	spookylogging "spooky/internal/logging"
	spookyschemas "spooky/internal/schemas"
)

// actorImpl implements the Actor interface
type actorImpl struct {
	// Action and context
	action  *spookyactionstypes.Action
	context *spookyactionstypes.ActionContext

	// State
	state    spookyactionstypes.ActingState
	status   spookyactionstypes.ActingStatus
	progress float64

	// Configuration
	timeout  time.Duration
	parallel bool

	// Dependencies
	logger spookylogging.Logger
}

// Execute executes the action
func (a *actorImpl) Execute(ctx context.Context, context *spookyactionstypes.ActionContext) (*spookyactionstypes.ActingResult, error) {
	if context == nil {
		return nil, fmt.Errorf("context cannot be nil")
	}

	a.logger.Info("Executing actor", spookylogging.String("action", a.action.Name))

	// Update state
	a.state = spookyactionstypes.ActingStateRunning
	a.status = spookyactionstypes.ActingStatusRunning
	a.progress = 0.0

	// Create result
	result := spookyactionstypes.NewActingResult(a.action.Name, "default")

	// Set start time
	now := time.Now()
	result.StartTime = &now

	// Execute based on action type
	var err error
	switch a.action.Type {
	case "command":
		err = a.executeCommand(ctx, context, result)
	case "script":
		err = a.executeScript(ctx, context, result)
	case "template_deploy":
		err = a.executeTemplateDeploy(ctx, context, result)
	case "template_evaluate":
		err = a.executeTemplateEvaluate(ctx, context, result)
	case "template_validate":
		err = a.executeTemplateValidate(ctx, context, result)
	case "template_cleanup":
		err = a.executeTemplateCleanup(ctx, context, result)
	default:
		err = fmt.Errorf("unsupported action type: %s", a.action.Type)
	}

	// Set end time and duration
	endTime := time.Now()
	result.EndTime = &endTime
	if result.StartTime != nil {
		result.Duration = result.EndTime.Sub(*result.StartTime)
	}

	// Update result based on execution
	if err != nil {
		result.Status = spookyactionstypes.ActingStatusFailed
		result.Error = err.Error()
		a.state = spookyactionstypes.ActingStateFailed
		a.status = spookyactionstypes.ActingStatusFailed
	} else {
		result.Status = spookyactionstypes.ActingStatusCompleted
		result.ExitCode = 0
		a.state = spookyactionstypes.ActingStateCompleted
		a.status = spookyactionstypes.ActingStatusCompleted
	}

	a.progress = 100.0

	a.logger.Info("Actor execution completed",
		spookylogging.String("action", a.action.Name),
		spookylogging.String("status", string(result.Status)))

	return result, err
}

// Prepare prepares the actor for execution
func (a *actorImpl) Prepare(context *spookyactionstypes.ActionContext) error {
	if context == nil {
		return fmt.Errorf("context cannot be nil")
	}

	a.logger.Info("Preparing actor", spookylogging.String("action", a.action.Name))

	// Validate action
	if err := a.validateAction(); err != nil {
		return fmt.Errorf("action validation failed: %w", err)
	}

	// Prepare context
	if err := a.prepareContext(context); err != nil {
		return fmt.Errorf("context preparation failed: %w", err)
	}

	// Update state
	a.state = spookyactionstypes.ActingStatePending
	a.status = spookyactionstypes.ActingStatusPending
	a.progress = 0.0

	a.logger.Info("Actor prepared", spookylogging.String("action", a.action.Name))
	return nil
}

// Cancel cancels the actor execution
func (a *actorImpl) Cancel() error {
	a.logger.Info("Cancelling actor", spookylogging.String("action", a.action.Name))

	// Update state
	a.state = spookyactionstypes.ActingStateCancelled
	a.status = spookyactionstypes.ActingStatusCancelled

	a.logger.Info("Actor cancelled", spookylogging.String("action", a.action.Name))
	return nil
}

// GetState returns the current state
func (a *actorImpl) GetState() spookyactionstypes.ActingState {
	return a.state
}

// GetProgress returns the current progress
func (a *actorImpl) GetProgress() float64 {
	return a.progress
}

// GetStatus returns the current status
func (a *actorImpl) GetStatus() spookyactionstypes.ActingStatus {
	return a.status
}

// SetTimeout sets the timeout
func (a *actorImpl) SetTimeout(timeout time.Duration) {
	a.timeout = timeout
}

// SetParallel sets the parallel flag
func (a *actorImpl) SetParallel(parallel bool) {
	a.parallel = parallel
}

// validateAction validates the action
func (a *actorImpl) validateAction() error {
	if a.action == nil {
		return fmt.Errorf("action cannot be nil")
	}

	if a.action.Name == "" {
		return fmt.Errorf("action name cannot be empty")
	}

	// Validate action type
	switch a.action.Type {
	case "command":
		if a.action.Command == "" {
			return fmt.Errorf("command cannot be empty for command action")
		}
	case "script":
		if a.action.Script == "" {
			return fmt.Errorf("script cannot be empty for script action")
		}
	case "template_deploy", "template_evaluate", "template_validate", "template_cleanup":
		if a.action.Template == nil {
			return fmt.Errorf("template cannot be nil for template action")
		}
		if a.action.Template.Source == "" {
			return fmt.Errorf("template source cannot be empty")
		}
		if a.action.Template.Destination == "" {
			return fmt.Errorf("template destination cannot be empty")
		}
	default:
		return fmt.Errorf("unsupported action type: %s", a.action.Type)
	}

	return nil
}

// prepareContext prepares the context
func (a *actorImpl) prepareContext(context *spookyactionstypes.ActionContext) error {
	if context == nil {
		return fmt.Errorf("context cannot be nil")
	}

	// Set default values if not provided
	if context.Timeout == 0 {
		context.Timeout = a.timeout
	}

	if context.MaxConcurrent == 0 {
		context.MaxConcurrent = 1
	}

	// Validate machines if specified
	if len(a.action.Machines) > 0 {
		if context.Machines == nil || len(context.Machines) == 0 {
			return fmt.Errorf("machines not available in context")
		}
	}

	return nil
}

// executeCommand executes a command action
func (a *actorImpl) executeCommand(ctx context.Context, context *spookyactionstypes.ActionContext, result *spookyactionstypes.ActingResult) error {
	a.logger.Debug("Executing command action",
		spookylogging.String("action", a.action.Name),
		spookylogging.String("command", a.action.Command))

	// TODO: Implement actual command execution
	// This would involve:
	// 1. Resolving machines from context
	// 2. Executing command on each machine
	// 3. Collecting results
	// 4. Aggregating results

	result.Output = "Command executed successfully (placeholder)"
	result.ExitCode = 0

	return nil
}

// executeScript executes a script action
func (a *actorImpl) executeScript(ctx context.Context, context *spookyactionstypes.ActionContext, result *spookyactionstypes.ActingResult) error {
	a.logger.Debug("Executing script action",
		spookylogging.String("action", a.action.Name))

	// TODO: Implement actual script execution
	// This would involve:
	// 1. Resolving machines from context
	// 2. Uploading script to machines
	// 3. Executing script on each machine
	// 4. Collecting results
	// 5. Aggregating results

	result.Output = "Script executed successfully (placeholder)"
	result.ExitCode = 0

	return nil
}

// executeTemplateDeploy executes a template deploy action
func (a *actorImpl) executeTemplateDeploy(ctx context.Context, context *spookyactionstypes.ActionContext, result *spookyactionstypes.ActingResult) error {
	a.logger.Debug("Executing template deploy action",
		spookylogging.String("action", a.action.Name),
		spookylogging.String("template", a.action.Template.Source))

	// TODO: Implement actual template deployment
	// This would involve:
	// 1. Resolving template variables
	// 2. Rendering template
	// 3. Deploying to target machines
	// 4. Validating deployment
	// 5. Collecting results

	result.Output = "Template deployed successfully (placeholder)"
	result.ExitCode = 0

	return nil
}

// executeTemplateEvaluate executes a template evaluate action
func (a *actorImpl) executeTemplateEvaluate(ctx context.Context, context *spookyactionstypes.ActionContext, result *spookyactionstypes.ActingResult) error {
	a.logger.Debug("Executing template evaluate action",
		spookylogging.String("action", a.action.Name),
		spookylogging.String("template", a.action.Template.Source))

	// TODO: Implement actual template evaluation
	// This would involve:
	// 1. Resolving template variables
	// 2. Evaluating template without deployment
	// 3. Collecting evaluation results

	result.Output = "Template evaluated successfully (placeholder)"
	result.ExitCode = 0

	return nil
}

// executeTemplateValidate executes a template validate action
func (a *actorImpl) executeTemplateValidate(ctx context.Context, context *spookyactionstypes.ActionContext, result *spookyactionstypes.ActingResult) error {
	a.logger.Debug("Executing template validate action",
		spookylogging.String("action", a.action.Name),
		spookylogging.String("template", a.action.Template.Source))

	// Use schema system for template validation
	if err := a.validateTemplate([]byte(a.action.Template.Source)); err != nil {
		result.Status = spookyactionstypes.ActingStatusFailed
		result.Error = err.Error()
		return fmt.Errorf("template validation failed: %w", err)
	}

	result.Output = "Template validated successfully using schema system"
	result.ExitCode = 0

	return nil
}

// executeTemplateCleanup executes a template cleanup action
func (a *actorImpl) executeTemplateCleanup(ctx context.Context, context *spookyactionstypes.ActionContext, result *spookyactionstypes.ActingResult) error {
	a.logger.Debug("Executing template cleanup action",
		spookylogging.String("action", a.action.Name),
		spookylogging.String("template", a.action.Template.Source))

	// TODO: Implement actual template cleanup
	// This would involve:
	// 1. Identifying deployed templates
	// 2. Removing deployed files
	// 3. Cleaning up backups
	// 4. Collecting cleanup results

	result.Output = "Template cleanup completed successfully (placeholder)"
	result.ExitCode = 0

	return nil
}

// validateTemplate validates template using schema system
func (a *actorImpl) validateTemplate(content []byte) error {
	validator := spookyschemas.NewSchemaValidator()
	if err := validator.LoadSchema(spookyschemas.SchemaTypeTemplateMetadata); err != nil {
		return fmt.Errorf("failed to load template metadata schema: %w", err)
	}

	// Parse content to interface{} for validation
	var data interface{}
	if err := json.Unmarshal(content, &data); err != nil {
		return fmt.Errorf("failed to parse template for validation: %w", err)
	}

	if err := validator.ValidateData(data, "template-metadata"); err != nil {
		return fmt.Errorf("template validation failed: %w", err)
	}
	return nil
}
