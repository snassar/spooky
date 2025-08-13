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

// LoadActions loads actions from the specified source
func (m *Manager) LoadActions(ctx context.Context, source string) ([]spookytypes.Action, error) {
	m.logger.Info("Loading actions", map[string]interface{}{
		"source": source,
	})

	// Check if source is a directory
	if info, err := os.Stat(source); err == nil && info.IsDir() {
		return m.loadActionsFromDirectory(ctx, source)
	}

	// Check if source is a file
	if _, err := os.Stat(source); err == nil {
		return m.loadActionsFromFile(ctx, source)
	}

	return nil, fmt.Errorf("source not found: %s", source)
}

// loadActionsFromFile loads actions from a single HCL file
func (m *Manager) loadActionsFromFile(ctx context.Context, filePath string) ([]spookytypes.Action, error) {
	m.logger.Debug("Loading actions from file", map[string]interface{}{
		"file": filePath,
	})

	// Validate file against schema
	if _, err := m.schemaValidator.ValidateFile(filePath, "actions"); err != nil {
		return nil, fmt.Errorf("schema validation failed for %s: %w", filePath, err)
	}

	// TODO: Implement actual HCL parsing
	// For now, return empty slice
	m.logger.Info("Actions loaded from file", map[string]interface{}{
		"file":    filePath,
		"actions": 0,
	})

	return []spookytypes.Action{}, nil
}

// loadActionsFromDirectory loads actions from multiple HCL files in a directory
func (m *Manager) loadActionsFromDirectory(ctx context.Context, dirPath string) ([]spookytypes.Action, error) {
	m.logger.Debug("Loading actions from directory", map[string]interface{}{
		"directory": dirPath,
	})

	var allActions []spookytypes.Action

	// Read directory entries
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory %s: %w", dirPath, err)
	}

	// Process each .hcl file
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

	m.logger.Info("Actions loaded from directory", map[string]interface{}{
		"directory": dirPath,
		"files":     len(entries),
		"actions":   len(allActions),
	})

	return allActions, nil
}

// ValidateActions validates a collection of actions
func (m *Manager) ValidateActions(ctx context.Context, actions []spookytypes.Action) (*spookytypes.ValidationResult, error) {
	m.logger.Info("Validating actions", map[string]interface{}{
		"count": len(actions),
	})

	return m.validator.ValidateActions(ctx, actions)
}

// RunActions runs the specified actions on the target machines
func (m *Manager) RunActions(ctx context.Context, actions []spookytypes.Action, machines []spookytypes.Machine) ([]spookytypes.ActingResult, error) {
	m.logger.Info("Running actions", map[string]interface{}{
		"actions":  len(actions),
		"machines": len(machines),
	})

	// Create acting session
	session := &spookytypesactions.ActingSession{
		SessionID:    "session-" + fmt.Sprintf("%d", ctx.Value("session_id")),
		Status:       "active",
		TotalActions: len(actions),
	}

	// Plan action running
	plan, err := m.createActionPlan(ctx, actions, machines)
	if err != nil {
		m.logger.Error("Failed to create action plan", err, map[string]interface{}{
			"actions":  len(actions),
			"machines": len(machines),
		})
		return nil, fmt.Errorf("failed to create action plan: %w", err)
	}

	// Run action plan
	results, err := m.runActionPlan(ctx, session, plan)
	if err != nil {
		m.logger.Error("Failed to run action plan", err, map[string]interface{}{
			"plan": plan.Name,
		})
		return nil, fmt.Errorf("failed to run action plan: %w", err)
	}

	m.logger.Info("Action running completed", map[string]interface{}{
		"actions": len(actions),
		"results": len(results),
	})

	return results, nil
}

// createActionPlan creates a running plan for the given actions and machines
func (m *Manager) createActionPlan(ctx context.Context, actions []spookytypes.Action, machines []spookytypes.Machine) (*spookytypesactions.ActionPlan, error) {
	m.logger.Debug("Creating action plan", map[string]interface{}{
		"actions":  len(actions),
		"machines": len(machines),
	})

	// Convert actions to internal format
	var internalActions []*spookytypesactions.Action
	for _, action := range actions {
		internalActions = append(internalActions, &action)
	}

	// Create plan
	plan := &spookytypesactions.ActionPlan{
		Name:          "action-running-plan",
		Description:   "Generated action running plan",
		Actions:       internalActions,
		Machines:      machines,
		Parallel:      true,
		MaxConcurrent: 4,
		Timeout:       300,
	}

	// Resolve dependencies and create running order
	runningOrder, dependencies, err := m.resolveDependencies(plan.Actions)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve dependencies: %w", err)
	}

	plan.RunOrder = runningOrder
	plan.Dependencies = dependencies

	m.logger.Info("Action plan created", map[string]interface{}{
		"plan":          plan.Name,
		"actions":       len(plan.Actions),
		"running_steps": len(runningOrder),
		"dependencies":  len(dependencies),
	})

	return plan, nil
}

// runActionPlan runs the actions according to the plan
func (m *Manager) runActionPlan(ctx context.Context, session *spookytypesactions.ActingSession, plan *spookytypesactions.ActionPlan) ([]spookytypes.ActingResult, error) {
	m.logger.Info("Running action plan", map[string]interface{}{
		"plan":    plan.Name,
		"actions": len(plan.Actions),
		"steps":   len(plan.RunOrder),
	})

	var allResults []spookytypes.ActingResult

	// Run each step in order
	for stepIndex, stepActions := range plan.RunOrder {
		m.logger.Debug("Running action step", map[string]interface{}{
			"step":    stepIndex + 1,
			"actions": stepActions,
		})

		stepResults, err := m.runActionStep(ctx, session, stepActions, plan.Actions)
		if err != nil {
			m.logger.Error("Failed to run action step", err, map[string]interface{}{
				"step":    stepIndex + 1,
				"actions": stepActions,
			})
			return allResults, fmt.Errorf("failed to run action step %d: %w", stepIndex+1, err)
		}

		allResults = append(allResults, stepResults...)
	}

	m.logger.Info("Action plan completed", map[string]interface{}{
		"plan":    plan.Name,
		"results": len(allResults),
	})

	return allResults, nil
}

// runActionStep runs a single step of actions
func (m *Manager) runActionStep(ctx context.Context, session *spookytypesactions.ActingSession, actionNames []string, allActions []*spookytypesactions.Action) ([]spookytypes.ActingResult, error) {
	var results []spookytypes.ActingResult

	// Find actions by name
	for _, actionName := range actionNames {
		var action *spookytypesactions.Action
		for _, a := range allActions {
			if a.Name == actionName {
				action = a
				break
			}
		}

		if action == nil {
			m.logger.Error("Action not found", fmt.Errorf("action %s not found", actionName), map[string]interface{}{
				"action": actionName,
			})
			continue
		}

		actionResults, err := m.runAction(ctx, session, action)
		if err != nil {
			m.logger.Error("Failed to run action", err, map[string]interface{}{"action": actionName})
			continue
		}

		results = append(results, actionResults...)
	}

	return results, nil
}

// runAction runs a single action on all target machines
func (m *Manager) runAction(ctx context.Context, session *spookytypesactions.ActingSession, action *spookytypesactions.Action) ([]spookytypes.ActingResult, error) {
	m.logger.Debug("Running action", map[string]interface{}{
		"action": action.Name,
		"type":   action.Type,
	})

	// Get target machines
	targetMachines := m.getTargetMachines(action, session)

	var results []spookytypes.ActingResult

	// Run action on each target machine
	for _, machine := range targetMachines {
		result, err := m.runActionOnMachine(ctx, session, action, machine)
		if err != nil {
			m.logger.Error("Failed to run action on machine", err, map[string]interface{}{
				"action":  action.Name,
				"machine": machine.Hostname,
			})
			continue
		}

		results = append(results, result)
	}

	return results, nil
}

// runActionOnMachine runs a single action on a specific machine
func (m *Manager) runActionOnMachine(ctx context.Context, _ *spookytypesactions.ActingSession, action *spookytypesactions.Action, machine spookytypes.Machine) (spookytypes.ActingResult, error) {
	m.logger.Debug("Running action on machine", map[string]interface{}{
		"action":  action.Name,
		"machine": machine.Hostname,
		"type":    action.Type,
	})

	// Create result
	result := &spookytypesactions.ActingResult{
		ActionName:  action.Name,
		MachineName: machine.Hostname,
		Status:      "success",
		StartTime:   time.Now(),
		EndTime:     time.Now(),
		ExitCode:    0,
	}

	// Run action based on type
	var err error
	switch action.Type {
	case "command":
		err = m.runCommandAction(ctx, action, machine, result)
	case "script":
		err = m.runScriptAction(ctx, action, machine, result)
	case "template_deploy":
		err = m.runTemplateAction(ctx, action, machine, result)
	case "file_copy":
		err = m.runFileCopyAction(ctx, action, machine, result)
	case "service_control":
		err = m.runServiceControlAction(ctx, action, machine, result)
	default:
		err = fmt.Errorf("unsupported action type: %s", action.Type)
	}

	if err != nil {
		result.Status = "failure"
		result.ErrorType = "running"
		result.Error = err.Error()
		m.logger.Error("Action running failed", err, map[string]interface{}{
			"action":  action.Name,
			"machine": machine.Hostname,
			"type":    action.Type,
		})
	} else {
		m.logger.Info("Action running completed", map[string]interface{}{
			"action":  action.Name,
			"machine": machine.Hostname,
			"type":    action.Type,
		})
	}

	return *result, err
}

// runCommandAction runs a command action
func (m *Manager) runCommandAction(_ context.Context, action *spookytypesactions.Action, machine spookytypes.Machine, result *spookytypesactions.ActingResult) error {
	// TODO: Implement proper SSH connection and command running
	result.Stdout = "Command running placeholder"
	return nil
}

// runScriptAction runs a script action
func (m *Manager) runScriptAction(ctx context.Context, action *spookytypesactions.Action, machine spookytypes.Machine, result *spookytypesactions.ActingResult) error {
	// TODO: Implement script running
	return fmt.Errorf("script running not yet implemented")
}

// runTemplateAction runs a template deployment action
func (m *Manager) runTemplateAction(ctx context.Context, action *spookytypesactions.Action, machine spookytypes.Machine, result *spookytypesactions.ActingResult) error {
	// TODO: Implement template deployment
	return fmt.Errorf("template deployment not yet implemented")
}

// runFileCopyAction runs a file copy action
func (m *Manager) runFileCopyAction(ctx context.Context, action *spookytypesactions.Action, machine spookytypes.Machine, result *spookytypesactions.ActingResult) error {
	// TODO: Implement file copy
	return fmt.Errorf("file copy not yet implemented")
}

// runServiceControlAction runs a service control action
func (m *Manager) runServiceControlAction(ctx context.Context, action *spookytypesactions.Action, machine spookytypes.Machine, result *spookytypesactions.ActingResult) error {
	if action.ServiceControl == nil {
		return fmt.Errorf("service control action requires service control configuration")
	}

	// TODO: Implement service control
	return fmt.Errorf("service control not yet implemented")
}

// getTargetMachines determines which machines should run the action
func (m *Manager) getTargetMachines(action *spookytypesactions.Action, session *spookytypesactions.ActingSession) []spookytypes.Machine {
	// TODO: Implement machine targeting logic
	return []spookytypes.Machine{}
}

// resolveDependencies resolves action dependencies and creates running order
func (m *Manager) resolveDependencies(actions []*spookytypesactions.Action) ([][]string, map[string][]string, error) {
	// Build dependency graph
	dependencies := make(map[string][]string)
	for _, action := range actions {
		dependencies[action.Name] = action.Dependencies
	}

	// Check for circular dependencies
	if err := m.checkCircularDependencies(dependencies); err != nil {
		return nil, nil, err
	}

	// Create running order using topological sort
	runningOrder := m.topologicalSort(dependencies)
	return runningOrder, dependencies, nil
}

// checkCircularDependencies checks for circular dependencies in the action graph
func (m *Manager) checkCircularDependencies(dependencies map[string][]string) error {
	visited := make(map[string]bool)
	recStack := make(map[string]bool)

	for action := range dependencies {
		if !visited[action] {
			if m.hasCycle(action, dependencies, visited, recStack) {
				return fmt.Errorf("circular dependency detected")
			}
		}
	}

	return nil
}

// hasCycle checks if there's a cycle in the dependency graph
func (m *Manager) hasCycle(action string, dependencies map[string][]string, visited, recStack map[string]bool) bool {
	visited[action] = true
	recStack[action] = true

	for _, dep := range dependencies[action] {
		if !visited[dep] {
			if m.hasCycle(dep, dependencies, visited, recStack) {
				return true
			}
		} else if recStack[dep] {
			return true
		}
	}

	recStack[action] = false
	return false
}

// topologicalSort performs topological sort on the dependency graph
func (m *Manager) topologicalSort(dependencies map[string][]string) [][]string {
	// Calculate in-degrees
	inDegree := make(map[string]int)
	for action := range dependencies {
		inDegree[action] = 0
	}

	for _, deps := range dependencies {
		for _, dep := range deps {
			inDegree[dep]++
		}
	}

	// Find nodes with no dependencies
	var queue []string
	for action, degree := range inDegree {
		if degree == 0 {
			queue = append(queue, action)
		}
	}

	var runningOrder [][]string
	for len(queue) > 0 {
		var currentLevel []string
		var nextQueue []string

		for _, action := range queue {
			currentLevel = append(currentLevel, action)

			// Reduce in-degree for dependent actions
			for _, dep := range dependencies[action] {
				inDegree[dep]--
				if inDegree[dep] == 0 {
					nextQueue = append(nextQueue, dep)
				}
			}
		}

		runningOrder = append(runningOrder, currentLevel)
		queue = nextQueue
	}

	return runningOrder
}

// machineHasTags checks if a machine has the specified tags
func (m *Manager) machineHasTags(machine *spookytypes.Machine, tags []string) bool {
	if len(tags) == 0 {
		return true
	}

	machineTags := make(map[string]bool)
	for _, tag := range machine.Tags {
		machineTags[tag] = true
	}

	for _, requiredTag := range tags {
		if !machineTags[requiredTag] {
			return false
		}
	}

	return true
}
