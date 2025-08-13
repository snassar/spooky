package spookyactions

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	spookyinterfaces "spooky/internal/interfaces"
	spookylogging "spooky/internal/logging"
	spookyschemas "spooky/internal/schemas"
	spookytypes "spooky/internal/types"
	spookytypesactions "spooky/internal/types/actions"
)

// Manager implements the ActionsIntegration interface
type Manager struct {
	logger          spookylogging.Logger
	validator       spookyinterfaces.ActionValidator
	sshManager      spookyinterfaces.SSHManager
	schemaValidator *spookyschemas.Validator
}

// NewManager creates a new actions manager
func NewManager(
	logger spookylogging.Logger,
	validator spookyinterfaces.ActionValidator,
	sshManager spookyinterfaces.SSHManager,
	schemaValidator *spookyschemas.Validator,
) spookyinterfaces.ActionsIntegration {
	return &Manager{
		logger:          logger,
		validator:       validator,
		sshManager:      sshManager,
		schemaValidator: schemaValidator,
	}
}

// LoadActions loads actions from the given source (project path)
func (m *Manager) LoadActions(ctx context.Context, source string) ([]spookytypes.Action, error) {
	m.logger.Info("Loading actions", map[string]interface{}{"source": source})

	// Load actions from actions.hcl file
	actionsFile := filepath.Join(source, "actions.hcl")
	actions, err := m.loadActionsFromFile(ctx, actionsFile)
	if err != nil {
		m.logger.Warn("Failed to load actions from actions.hcl", map[string]interface{}{"error": err.Error()})
		actions = []spookytypes.Action{}
	}

	// Load actions from actions/ directory
	actionsDir := filepath.Join(source, "actions")
	dirActions, err := m.loadActionsFromDirectory(ctx, actionsDir)
	if err != nil {
		m.logger.Warn("Failed to load actions from actions/ directory", map[string]interface{}{"error": err.Error()})
	} else {
		actions = append(actions, dirActions...)
	}

	// Validate loaded actions
	if len(actions) > 0 {
		validationResult, err := m.ValidateActions(ctx, actions)
		if err != nil {
			m.logger.Error("Failed to validate actions", err)
			return nil, fmt.Errorf("failed to validate actions: %w", err)
		}

		if !validationResult.Valid {
			m.logger.Error("Action validation failed", fmt.Errorf("validation failed"), map[string]interface{}{"errors": validationResult.Errors})
			return nil, fmt.Errorf("action validation failed: %v", validationResult.Errors)
		}
	}

	m.logger.Info("Successfully loaded actions", map[string]interface{}{"count": len(actions)})
	return actions, nil
}

// ValidateActions validates the given actions
func (m *Manager) ValidateActions(ctx context.Context, actions []spookytypes.Action) (*spookytypes.ValidationResult, error) {
	m.logger.Info("Validating actions", map[string]interface{}{"count": len(actions)})

	if m.validator == nil {
		return &spookytypes.ValidationResult{
			Valid:    true,
			Errors:   []spookytypes.SchemaError{},
			Warnings: []spookytypes.SchemaError{},
		}, nil
	}

	validationResult, err := m.validator.ValidateActions(ctx, actions)
	if err != nil {
		m.logger.Error("Failed to validate actions", err)
		return nil, fmt.Errorf("failed to validate actions: %w", err)
	}

	m.logger.Info("Action validation completed", map[string]interface{}{
		"valid":  validationResult.Valid,
		"errors": len(validationResult.Errors),
	})
	return validationResult, nil
}

// RunActions runs the given actions on the specified machines
func (m *Manager) RunActions(ctx context.Context, actions []spookytypes.Action, machines []spookytypes.Machine) ([]spookytypes.ActingResult, error) {
	m.logger.Info("Running actions", map[string]interface{}{
		"action_count":  len(actions),
		"machine_count": len(machines),
	})

	// Create acting session
	session := &spookytypesactions.ActingSession{
		SessionID:   generateSessionID(),
		ProjectPath: getProjectPathFromContext(ctx),
		Actions:     actions,
		Machines:    machines,
		Status:      string(spookytypesactions.SessionStatusPending),
		StartTime:   &[]time.Time{time.Now()}[0],
		Total:       len(actions) * len(machines),
	}

	// Plan action execution
	plan, err := m.createActionPlan(ctx, actions, machines)
	if err != nil {
		m.logger.Error("Failed to create action plan", err)
		return nil, fmt.Errorf("failed to create action plan: %w", err)
	}

	// Execute actions according to plan
	results, err := m.executeActionPlan(ctx, session, plan)
	if err != nil {
		m.logger.Error("Failed to execute action plan", err)
		return nil, fmt.Errorf("failed to execute action plan: %w", err)
	}

	// Update session status
	endTime := time.Now()
	session.EndTime = &endTime
	session.Status = string(spookytypesactions.SessionStatusCompleted)
	session.Results = results

	m.logger.Info("Action execution completed", map[string]interface{}{
		"session_id": session.SessionID,
		"results":    len(results),
	})
	return results, nil
}

// loadActionsFromFile loads actions from a single HCL file
func (m *Manager) loadActionsFromFile(_ context.Context, filePath string) ([]spookytypes.Action, error) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return []spookytypes.Action{}, nil
	}

	_, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read actions file %s: %w", filePath, err)
	}

	// Validate against actions schema first
	if m.schemaValidator != nil {
		validationResult, err := m.schemaValidator.ValidateFile(filePath, "actions")
		if err != nil {
			return nil, fmt.Errorf("failed to validate actions file %s: %w", filePath, err)
		}

		if !validationResult.Valid {
			return nil, fmt.Errorf("actions file %s failed validation: %v", filePath, validationResult.Errors)
		}
	}

	// TODO: Implement proper HCL parsing for actions using schema-validated data
	// For now, return empty slice as placeholder
	var actions []spookytypes.Action

	m.logger.Debug("Loaded actions from file", map[string]interface{}{
		"file":  filePath,
		"count": len(actions),
	})
	return actions, nil
}

// loadActionsFromDirectory loads actions from all HCL files in a directory
func (m *Manager) loadActionsFromDirectory(ctx context.Context, dirPath string) ([]spookytypes.Action, error) {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		return []spookytypes.Action{}, nil
	}

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read actions directory %s: %w", dirPath, err)
	}

	var allActions []spookytypes.Action
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".hcl") {
			continue
		}

		filePath := filepath.Join(dirPath, entry.Name())
		actions, err := m.loadActionsFromFile(ctx, filePath)
		if err != nil {
			m.logger.Warn("Failed to load actions from file", map[string]interface{}{
				"file":  filePath,
				"error": err.Error(),
			})
			continue
		}

		allActions = append(allActions, actions...)
	}

	m.logger.Debug("Loaded actions from directory", map[string]interface{}{
		"directory": dirPath,
		"count":     len(allActions),
	})
	return allActions, nil
}

// createActionPlan creates an execution plan for the given actions and machines
func (m *Manager) createActionPlan(_ context.Context, actions []spookytypes.Action, machines []spookytypes.Machine) (*spookytypesactions.ActionPlan, error) {
	plan := &spookytypesactions.ActionPlan{
		PlanID:      generatePlanID(),
		PlanName:    "action-execution-plan",
		Description: "Generated action execution plan",
		Actions:     make([]*spookytypesactions.Action, len(actions)),
		Parallel:    true,
		Validated:   false,
	}

	// Convert actions to action pointers
	for i, action := range actions {
		// Since Action is a concrete type, we can cast directly
		plan.Actions[i] = &action
	}

	// Resolve dependencies and create execution order
	executionOrder, dependencies, err := m.resolveDependencies(plan.Actions)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve dependencies: %w", err)
	}

	plan.ExecutionOrder = executionOrder
	plan.Dependencies = dependencies
	plan.Validated = true

	m.logger.Debug("Created action plan", map[string]interface{}{
		"plan_id":         plan.PlanID,
		"execution_steps": len(executionOrder),
	})
	return plan, nil
}

// executeActionPlan executes the actions according to the plan
func (m *Manager) executeActionPlan(ctx context.Context, session *spookytypesactions.ActingSession, plan *spookytypesactions.ActionPlan) ([]spookytypes.ActingResult, error) {
	var allResults []spookytypes.ActingResult

	// Execute actions in dependency order
	for stepIndex, stepActions := range plan.ExecutionOrder {
		m.logger.Info("Executing action step", map[string]interface{}{
			"step":    stepIndex + 1,
			"actions": stepActions,
		})

		stepResults, err := m.executeActionStep(ctx, session, stepActions, plan.Actions)
		if err != nil {
			m.logger.Error("Failed to execute action step", err, map[string]interface{}{
				"step": stepIndex + 1,
			})
			return allResults, fmt.Errorf("failed to execute action step %d: %w", stepIndex+1, err)
		}

		allResults = append(allResults, stepResults...)
	}

	return allResults, nil
}

// executeActionStep executes a single step of actions
func (m *Manager) executeActionStep(ctx context.Context, session *spookytypesactions.ActingSession, actionNames []string, allActions []*spookytypesactions.Action) ([]spookytypes.ActingResult, error) {
	var results []spookytypes.ActingResult

	// Create action map for quick lookup
	actionMap := make(map[string]*spookytypesactions.Action)
	for _, action := range allActions {
		actionMap[action.Name] = action
	}

	// Execute actions in parallel if configured
	for _, actionName := range actionNames {
		action, exists := actionMap[actionName]
		if !exists {
			m.logger.Warn("Action not found in plan", map[string]interface{}{"action": actionName})
			continue
		}

		actionResults, err := m.executeAction(ctx, session, action)
		if err != nil {
			m.logger.Error("Failed to execute action", err, map[string]interface{}{"action": actionName})
			if !action.AllowFailure {
				return results, fmt.Errorf("action %s failed: %w", actionName, err)
			}
		}

		results = append(results, actionResults...)
	}

	return results, nil
}

// executeAction executes a single action on all target machines
func (m *Manager) executeAction(ctx context.Context, session *spookytypesactions.ActingSession, action *spookytypesactions.Action) ([]spookytypes.ActingResult, error) {
	var results []spookytypes.ActingResult

	// Determine target machines
	targetMachines := m.getTargetMachines(action, session.Machines)
	if len(targetMachines) == 0 {
		m.logger.Warn("No target machines found for action", map[string]interface{}{"action": action.Name})
		return results, nil
	}

	// Execute on each target machine
	for _, machine := range targetMachines {
		result, err := m.executeActionOnMachine(ctx, session, action, machine)
		if err != nil {
			m.logger.Error("Failed to execute action on machine", err, map[string]interface{}{
				"action":  action.Name,
				"machine": machine.Name,
			})
			if !action.AllowFailure {
				return results, fmt.Errorf("action %s failed on machine %s: %w", action.Name, machine.Name, err)
			}
		}

		results = append(results, result)
	}

	return results, nil
}

// executeActionOnMachine executes a single action on a specific machine
func (m *Manager) executeActionOnMachine(ctx context.Context, _ *spookytypesactions.ActingSession, action *spookytypesactions.Action, machine spookytypes.Machine) (spookytypes.ActingResult, error) {
	startTime := time.Now()
	result := &spookytypesactions.ActingResult{
		ActionName:  action.Name,
		ActionType:  action.Type,
		MachineName: machine.Name,
		MachineHost: machine.Hostname,
		StartTime:   &startTime,
		Status:      string(spookytypesactions.ResultStatusRunning),
		RetryCount:  0,
	}

	m.logger.Info("Executing action on machine", map[string]interface{}{
		"action":  action.Name,
		"machine": machine.Name,
		"type":    action.Type,
	})

	// Execute based on action type
	var err error
	switch action.Type {
	case string(spookytypesactions.ActionTypeCommand):
		err = m.executeCommandAction(ctx, action, machine, result)
	case string(spookytypesactions.ActionTypeScript):
		err = m.executeScriptAction(ctx, action, machine, result)
	case string(spookytypesactions.ActionTypeTemplateDeploy):
		err = m.executeTemplateAction(ctx, action, machine, result)
	case string(spookytypesactions.ActionTypeFileCopy):
		err = m.executeFileCopyAction(ctx, action, machine, result)
	case string(spookytypesactions.ActionTypeServiceControl):
		err = m.executeServiceControlAction(ctx, action, machine, result)
	default:
		err = fmt.Errorf("unsupported action type: %s", action.Type)
	}

	// Update result
	endTime := time.Now()
	result.EndTime = &endTime
	result.Duration = endTime.Sub(startTime)

	if err != nil {
		result.Status = string(spookytypesactions.ResultStatusFailed)
		result.Error = err
		result.ErrorType = "execution"
		result.ErrorMessage = err.Error()
		m.logger.Error("Action execution failed", err, map[string]interface{}{
			"action":  action.Name,
			"machine": machine.Name,
		})
	} else {
		result.Status = string(spookytypesactions.ResultStatusCompleted)
		m.logger.Info("Action execution completed", map[string]interface{}{
			"action":   action.Name,
			"machine":  machine.Name,
			"duration": result.Duration,
		})
	}

	return *result, err
}

// executeCommandAction executes a command action
func (m *Manager) executeCommandAction(_ context.Context, action *spookytypesactions.Action, machine spookytypes.Machine, result *spookytypesactions.ActingResult) error {
	if action.Command == "" {
		return fmt.Errorf("command action requires a command")
	}

	// TODO: Implement proper SSH connection and command execution
	// For now, return placeholder result
	result.ExitCode = 0
	result.Stdout = "Command execution placeholder"
	result.Stderr = ""
	result.Status = string(spookytypesactions.ResultStatusCompleted)

	return nil
}

// executeScriptAction executes a script action
func (m *Manager) executeScriptAction(ctx context.Context, action *spookytypesactions.Action, machine spookytypes.Machine, result *spookytypesactions.ActingResult) error {
	if action.Script == "" {
		return fmt.Errorf("script action requires a script path")
	}

	// TODO: Implement script execution
	// This would involve:
	// 1. Reading the script file
	// 2. Template processing if it's a .tmpl file
	// 3. Uploading to the target machine
	// 4. Executing the script
	// 5. Cleaning up

	return fmt.Errorf("script execution not yet implemented")
}

// executeTemplateAction executes a template deployment action
func (m *Manager) executeTemplateAction(ctx context.Context, action *spookytypesactions.Action, machine spookytypes.Machine, result *spookytypesactions.ActingResult) error {
	if action.Template == nil {
		return fmt.Errorf("template action requires template configuration")
	}

	// TODO: Implement template deployment
	// This would involve:
	// 1. Loading the template
	// 2. Rendering with variables
	// 3. Uploading to target machine
	// 4. Setting permissions and ownership

	return fmt.Errorf("template deployment not yet implemented")
}

// executeFileCopyAction executes a file copy action
func (m *Manager) executeFileCopyAction(ctx context.Context, action *spookytypesactions.Action, machine spookytypes.Machine, result *spookytypesactions.ActingResult) error {
	if action.FileCopy == nil {
		return fmt.Errorf("file copy action requires file copy configuration")
	}

	// TODO: Implement file copy
	// This would involve:
	// 1. Reading the source file
	// 2. Uploading to target machine
	// 3. Setting permissions and ownership

	return fmt.Errorf("file copy not yet implemented")
}

// executeServiceControlAction executes a service control action
func (m *Manager) executeServiceControlAction(ctx context.Context, action *spookytypesactions.Action, machine spookytypes.Machine, result *spookytypesactions.ActingResult) error {
	if action.ServiceControl == nil {
		return fmt.Errorf("service control action requires service control configuration")
	}

	// TODO: Implement service control
	// This would involve:
	// 1. Building the appropriate service command
	// 2. Executing via SSH
	// 3. Waiting for status if configured

	return fmt.Errorf("service control not yet implemented")
}

// getTargetMachines determines which machines should execute the action
func (m *Manager) getTargetMachines(action *spookytypesactions.Action, allMachines []spookytypes.Machine) []spookytypes.Machine {
	var targetMachines []spookytypes.Machine

	// If specific machines are specified, use those
	if len(action.Machines) > 0 {
		machineMap := make(map[string]spookytypes.Machine)
		for _, machine := range allMachines {
			machineMap[machine.Name] = machine
		}

		for _, machineName := range action.Machines {
			if machine, exists := machineMap[machineName]; exists {
				targetMachines = append(targetMachines, machine)
			}
		}
		return targetMachines
	}

	// If tags are specified, filter by tags
	if len(action.Tags) > 0 {
		for _, machine := range allMachines {
			if m.machineHasTags(&machine, action.Tags) {
				targetMachines = append(targetMachines, machine)
			}
		}
		return targetMachines
	}

	// If no targeting specified, use all machines
	return allMachines
}

// machineHasTags checks if a machine has any of the specified tags
func (m *Manager) machineHasTags(machine *spookytypes.Machine, tags []string) bool {
	machineTags := make(map[string]bool)
	for _, tag := range machine.Tags {
		machineTags[tag] = true
	}

	for _, tag := range tags {
		if machineTags[tag] {
			return true
		}
	}
	return false
}

// resolveDependencies resolves action dependencies and creates execution order
func (m *Manager) resolveDependencies(actions []*spookytypesactions.Action) ([][]string, map[string][]string, error) {
	// Build dependency graph
	dependencies := make(map[string][]string)
	reverseDeps := make(map[string][]string)
	actionMap := make(map[string]*spookytypesactions.Action)

	for _, action := range actions {
		actionMap[action.Name] = action
		dependencies[action.Name] = action.Dependencies
		for _, dep := range action.Dependencies {
			reverseDeps[dep] = append(reverseDeps[dep], action.Name)
		}
	}

	// Check for circular dependencies
	if err := m.checkCircularDependencies(dependencies); err != nil {
		return nil, nil, err
	}

	// Create execution order using topological sort
	executionOrder := m.topologicalSort(dependencies)

	return executionOrder, dependencies, nil
}

// checkCircularDependencies checks for circular dependencies in the action graph
func (m *Manager) checkCircularDependencies(dependencies map[string][]string) error {
	visited := make(map[string]bool)
	recStack := make(map[string]bool)

	var dfs func(string) error
	dfs = func(node string) error {
		visited[node] = true
		recStack[node] = true

		for _, dep := range dependencies[node] {
			if !visited[dep] {
				if err := dfs(dep); err != nil {
					return err
				}
			} else if recStack[dep] {
				return fmt.Errorf("circular dependency detected involving %s", node)
			}
		}

		recStack[node] = false
		return nil
	}

	for node := range dependencies {
		if !visited[node] {
			if err := dfs(node); err != nil {
				return err
			}
		}
	}

	return nil
}

// topologicalSort performs topological sorting of actions based on dependencies
func (m *Manager) topologicalSort(dependencies map[string][]string) [][]string {
	inDegree := make(map[string]int)
	for action, deps := range dependencies {
		if _, exists := inDegree[action]; !exists {
			inDegree[action] = 0
		}
		for _, dep := range deps {
			inDegree[dep]++
		}
	}

	var queue []string
	for action, degree := range inDegree {
		if degree == 0 {
			queue = append(queue, action)
		}
	}

	var executionOrder [][]string
	for len(queue) > 0 {
		var currentLevel []string
		var nextQueue []string

		for _, action := range queue {
			currentLevel = append(currentLevel, action)
			for _, dep := range dependencies[action] {
				inDegree[dep]--
				if inDegree[dep] == 0 {
					nextQueue = append(nextQueue, dep)
				}
			}
		}

		if len(currentLevel) > 0 {
			executionOrder = append(executionOrder, currentLevel)
		}
		queue = nextQueue
	}

	return executionOrder
}

// Helper functions
func generateSessionID() string {
	return fmt.Sprintf("session-%d", time.Now().UnixNano())
}

func generatePlanID() string {
	return fmt.Sprintf("plan-%d", time.Now().UnixNano())
}

func getProjectPathFromContext(ctx context.Context) string {
	// TODO: Extract project path from context
	return ""
}
