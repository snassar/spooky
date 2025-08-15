package spookyactions

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	spookyinterfaces "spooky/internal/interfaces"
	spookyschemas "spooky/internal/schemas"
	spookytypes "spooky/internal/types"
	spookytypesactions "spooky/internal/types/actions"
	spookytypeslogging "spooky/internal/types/logging"

	"github.com/hashicorp/hcl/v2/hclsimple"
)

// Manager implements the ActionsIntegration interface
type Manager struct {
	logger          spookytypeslogging.Logger
	validator       spookyinterfaces.ActionValidator
	sshManager      spookyinterfaces.SSHManager
	schemaValidator *spookyschemas.Validator
}

// NewManager creates a new actions manager
func NewManager(
	logger spookytypeslogging.Logger,
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
func (m *Manager) loadActionsFromFile(_ context.Context, filePath string) ([]spookytypes.Action, error) {
	m.logger.Debug("Loading actions from file", map[string]interface{}{
		"file": filePath,
	})

	// Validate file against schema
	if _, err := m.schemaValidator.ValidateFile(filePath, "actions"); err != nil {
		return nil, fmt.Errorf("schema validation failed for %s: %w", filePath, err)
	}

	// Parse HCL file
	var actionContainer struct {
		Actions []*spookytypesactions.Action `hcl:"action,block"`
	}
	if err := hclsimple.DecodeFile(filePath, nil, &actionContainer); err != nil {
		return nil, fmt.Errorf("failed to parse HCL content from %s: %w", filePath, err)
	}

	// Convert to spookytypes.Action slice
	var actions []spookytypes.Action
	for _, action := range actionContainer.Actions {
		// Convert *spookytypesactions.Action to spookytypes.Action
		actions = append(actions, *action)
	}

	m.logger.Info("Actions loaded from file", map[string]interface{}{
		"file":    filePath,
		"actions": len(actions),
	})

	return actions, nil
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
func (m *Manager) createActionPlan(_ context.Context, actions []spookytypes.Action, machines []spookytypes.Machine) (*spookytypesactions.ActionPlan, error) {
	m.logger.Debug("Creating action plan", map[string]interface{}{
		"actions":  len(actions),
		"machines": len(machines),
	})

	// Convert actions to internal format
	var internalActions []*spookytypesactions.Action
	for idx := range actions {
		action := &actions[idx]
		// Convert spookytypes.Action to spookytypesactions.Action
		internalAction := &spookytypesactions.Action{
			Name:         action.Name,
			Description:  action.Description,
			Type:         action.Type,
			Machines:     action.Machines,
			Tags:         action.Tags,
			Dependencies: action.Dependencies,
			// Add other fields as needed
		}
		internalActions = append(internalActions, internalAction)
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

		stepResults, err := m.runActionStep(ctx, session, stepActions, plan.Actions, plan)
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
func (m *Manager) runActionStep(ctx context.Context, session *spookytypesactions.ActingSession, actionNames []string, allActions []*spookytypesactions.Action, _ *spookytypesactions.ActionPlan) ([]spookytypes.ActingResult, error) {
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

// runAction runs a single action on target machines
func (m *Manager) runAction(ctx context.Context, session *spookytypesactions.ActingSession, action *spookytypesactions.Action) ([]spookytypes.ActingResult, error) {
	m.logger.Debug("Running action", map[string]interface{}{
		"action": action.Name,
		"type":   action.Type,
	})

	// Get target machines
	targetMachines := m.getTargetMachines(action, session, session.MachineInventory)

	var results []spookytypes.ActingResult

	// Run action on each target machine
	for i := range targetMachines {
		result, err := m.runActionOnMachine(ctx, session, action, &targetMachines[i])
		if err != nil {
			m.logger.Error("Failed to run action on machine", err, map[string]interface{}{
				"action":  action.Name,
				"machine": targetMachines[i].Hostname,
			})
			continue
		}

		results = append(results, result)
	}

	return results, nil
}

// runActionOnMachine runs a single action on a specific machine
func (m *Manager) runActionOnMachine(ctx context.Context, _ *spookytypesactions.ActingSession, action *spookytypesactions.Action, machine *spookytypes.Machine) (spookytypes.ActingResult, error) {
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

// setFilePermissions sets file permissions on a remote machine via SSH
func (m *Manager) setFilePermissions(ctx context.Context, session *spookytypes.Session, machine *spookytypes.Machine, filePath, permissions string, action *spookytypesactions.Action) error {
	chmodCommand := fmt.Sprintf("chmod %s %s", permissions, filePath)

	sshCommand := &spookytypes.SSHCommand{
		Command:       chmodCommand,
		WorkingDir:    action.WorkingDir,
		Environment:   action.Environment,
		Timeout:       time.Duration(action.Timeout) * time.Second,
		CaptureOutput: true,
	}

	commandResult, err := m.sshManager.RunCommand(ctx, session, sshCommand)
	if err != nil {
		return fmt.Errorf("failed to set file permissions on %s: %w", machine.Hostname, err)
	}

	if !commandResult.Success {
		return fmt.Errorf("failed to set file permissions on %s: %s", machine.Hostname, commandResult.Error)
	}

	return nil
}

// runCommandAction runs a command action on a machine
func (m *Manager) runCommandAction(ctx context.Context, action *spookytypesactions.Action, machine *spookytypes.Machine, result *spookytypesactions.ActingResult) error {
	m.logger.Debug("Running command action", map[string]interface{}{
		"action":  action.Name,
		"machine": machine.Hostname,
		"command": action.CommandString,
	})

	// Create SSH connection
	connectionRequest := &spookytypes.ConnectionRequest{
		Host:     machine.Host,
		Port:     machine.Port,
		User:     machine.User,
		Password: machine.Password,
		KeyPath:  machine.KeyFile,
		Timeout:  time.Duration(action.Timeout) * time.Second,
	}

	connectionResult, err := m.sshManager.Connect(ctx, connectionRequest)
	if err != nil {
		return fmt.Errorf("failed to connect to %s: %w", machine.Hostname, err)
	}

	// Create SSH session
	session, err := m.sshManager.CreateSession(ctx, connectionResult.Connection)
	if err != nil {
		return fmt.Errorf("failed to create session on %s: %w", machine.Hostname, err)
	}

	// Prepare command
	commandStr := action.CommandString
	if action.Command != nil {
		commandStr = action.Command.Command
	}

	// Create SSH command
	sshCommand := &spookytypes.SSHCommand{
		Command:       commandStr,
		Args:          action.Command.Args,
		WorkingDir:    action.WorkingDir,
		Environment:   action.Environment,
		Timeout:       time.Duration(action.Timeout) * time.Second,
		CaptureOutput: true,
	}

	// Run command
	startTime := time.Now()
	commandResult, err := m.sshManager.RunCommand(ctx, session, sshCommand)
	endTime := time.Now()

	// Update result
	result.StartTime = startTime
	result.EndTime = endTime
	result.Duration = endTime.Sub(startTime)

	if err != nil {
		result.Status = "failed"
		result.Error = err.Error()
		return err
	}

	result.Status = "success"
	result.ExitCode = commandResult.ExitCode
	result.Stdout = commandResult.Stdout
	result.Stderr = commandResult.Stderr

	return nil
}

// runScriptAction runs a script action on a machine
func (m *Manager) runScriptAction(ctx context.Context, action *spookytypesactions.Action, machine *spookytypes.Machine, result *spookytypesactions.ActingResult) error {
	m.logger.Debug("Running script action", map[string]interface{}{
		"action":  action.Name,
		"machine": machine.Hostname,
		"script":  action.Script.Script,
	})

	// Create SSH connection
	connectionRequest := &spookytypes.ConnectionRequest{
		Host:     machine.Host,
		Port:     machine.Port,
		User:     machine.User,
		Password: machine.Password,
		KeyPath:  machine.KeyFile,
		Timeout:  time.Duration(action.Timeout) * time.Second,
	}

	connectionResult, err := m.sshManager.Connect(ctx, connectionRequest)
	if err != nil {
		return fmt.Errorf("failed to connect to %s: %w", machine.Hostname, err)
	}

	// Create SSH session
	session, err := m.sshManager.CreateSession(ctx, connectionResult.Connection)
	if err != nil {
		return fmt.Errorf("failed to create session on %s: %w", machine.Hostname, err)
	}

	// Prepare script command
	scriptCommand := action.Script.Script
	if action.Script.Shell != "" {
		scriptCommand = fmt.Sprintf("%s -c '%s'", action.Script.Shell, scriptCommand)
	}

	// Create SSH command
	sshCommand := &spookytypes.SSHCommand{
		Command:       scriptCommand,
		Args:          action.Script.Args,
		WorkingDir:    action.WorkingDir,
		Environment:   action.Environment,
		Timeout:       time.Duration(action.Timeout) * time.Second,
		CaptureOutput: true,
	}

	// Run script
	startTime := time.Now()
	commandResult, err := m.sshManager.RunCommand(ctx, session, sshCommand)
	endTime := time.Now()

	// Update result
	result.StartTime = startTime
	result.EndTime = endTime
	result.Duration = endTime.Sub(startTime)

	if err != nil {
		result.Status = "failed"
		result.Error = err.Error()
		return err
	}

	result.Status = "success"
	result.ExitCode = commandResult.ExitCode
	result.Stdout = commandResult.Stdout
	result.Stderr = commandResult.Stderr

	return nil
}

// runTemplateAction runs a template deployment action on a machine
func (m *Manager) runTemplateAction(ctx context.Context, action *spookytypesactions.Action, machine *spookytypes.Machine, result *spookytypesactions.ActingResult) error {
	return m.runFileTransferAction(ctx, action, machine, result, action.Template.Source, action.Template.Destination, action.Template.Permissions, "File transferred")
}

// runFileCopyAction runs a file copy action on a machine
func (m *Manager) runFileCopyAction(ctx context.Context, action *spookytypesactions.Action, machine *spookytypes.Machine, result *spookytypesactions.ActingResult) error {
	return m.runFileTransferAction(ctx, action, machine, result, action.FileCopy.Source, action.FileCopy.Destination, action.FileCopy.Permissions, "File copied")
}

// runFileTransferAction is a helper function for file transfer operations
func (m *Manager) runFileTransferAction(ctx context.Context, action *spookytypesactions.Action, machine *spookytypes.Machine, result *spookytypesactions.ActingResult, source, destination, permissions, successMessage string) error {
	m.logger.Debug("Running file transfer action", map[string]interface{}{
		"action":      action.Name,
		"machine":     machine.Hostname,
		"source":      source,
		"destination": destination,
	})

	// Create SSH connection
	connectionRequest := &spookytypes.ConnectionRequest{
		Host:     machine.Host,
		Port:     machine.Port,
		User:     machine.User,
		Password: machine.Password,
		KeyPath:  machine.KeyFile,
		Timeout:  time.Duration(action.Timeout) * time.Second,
	}

	connectionResult, err := m.sshManager.Connect(ctx, connectionRequest)
	if err != nil {
		return fmt.Errorf("failed to connect to %s: %w", machine.Hostname, err)
	}

	// Create SSH session
	session, err := m.sshManager.CreateSession(ctx, connectionResult.Connection)
	if err != nil {
		return fmt.Errorf("failed to create session on %s: %w", machine.Hostname, err)
	}

	// Create file transfer
	transfer := &spookytypes.FileTransfer{
		LocalPath:  source,
		RemotePath: destination,
		Direction:  "upload",
		Mode:       "sftp",
	}

	// Transfer file
	startTime := time.Now()
	_, err = m.sshManager.TransferFile(ctx, session, transfer)
	endTime := time.Now()

	// Update result
	result.StartTime = startTime
	result.EndTime = endTime
	result.Duration = endTime.Sub(startTime)

	if err != nil {
		result.Status = "failed"
		result.Error = err.Error()
		return err
	}

	// Set file permissions if specified
	if permissions != "" {
		if err := m.setFilePermissions(ctx, session, machine, destination, permissions, action); err != nil {
			result.Status = "failed"
			result.Error = fmt.Sprintf("file transfer succeeded but permission setting failed: %v", err)
			return err
		}
	}

	result.Status = "success"
	result.ExitCode = 0
	result.Stdout = fmt.Sprintf("%s successfully: %s -> %s", successMessage, source, destination)

	return nil
}

// runServiceControlAction runs a service control action on a machine
func (m *Manager) runServiceControlAction(ctx context.Context, action *spookytypesactions.Action, machine *spookytypes.Machine, result *spookytypesactions.ActingResult) error {
	m.logger.Debug("Running service control action", map[string]interface{}{
		"action":         action.Name,
		"machine":        machine.Hostname,
		"service":        action.ServiceControl.Service,
		"service_action": action.ServiceControl.Action,
	})

	// Create SSH connection
	connectionRequest := &spookytypes.ConnectionRequest{
		Host:     machine.Host,
		Port:     machine.Port,
		User:     machine.User,
		Password: machine.Password,
		KeyPath:  machine.KeyFile,
		Timeout:  time.Duration(action.Timeout) * time.Second,
	}

	connectionResult, err := m.sshManager.Connect(ctx, connectionRequest)
	if err != nil {
		return fmt.Errorf("failed to connect to %s: %w", machine.Hostname, err)
	}

	// Create SSH session
	session, err := m.sshManager.CreateSession(ctx, connectionResult.Connection)
	if err != nil {
		return fmt.Errorf("failed to create session on %s: %w", machine.Hostname, err)
	}

	// Prepare service command
	serviceCommand := fmt.Sprintf("systemctl %s %s", action.ServiceControl.Action, action.ServiceControl.Service)

	// Create SSH command
	sshCommand := &spookytypes.SSHCommand{
		Command:       serviceCommand,
		WorkingDir:    action.WorkingDir,
		Environment:   action.Environment,
		Timeout:       time.Duration(action.Timeout) * time.Second,
		CaptureOutput: true,
	}

	// Run service command
	startTime := time.Now()
	commandResult, err := m.sshManager.RunCommand(ctx, session, sshCommand)
	endTime := time.Now()

	// Update result
	result.StartTime = startTime
	result.EndTime = endTime
	result.Duration = endTime.Sub(startTime)

	if err != nil {
		result.Status = "failed"
		result.Error = err.Error()
		return err
	}

	result.Status = "success"
	result.ExitCode = commandResult.ExitCode
	result.Stdout = commandResult.Stdout
	result.Stderr = commandResult.Stderr

	return nil
}

// getTargetMachines determines which machines should run the action
func (m *Manager) getTargetMachines(action *spookytypesactions.Action, session *spookytypesactions.ActingSession, availableMachines []spookytypes.Machine) []spookytypes.Machine {
	// If session has machine inventory, use it for better performance
	if session != nil && len(session.MachineInventory) > 0 {
		return m.getTargetMachinesFromSession(action, session)
	}

	// Fallback to the passed availableMachines parameter
	// If action has specific machines defined, filter by name
	if len(action.Machines) > 0 {
		var targetMachines []spookytypes.Machine
		var missingMachines []string

		for _, targetName := range action.Machines {
			found := false
			for i := range availableMachines {
				if availableMachines[i].Hostname == targetName {
					targetMachines = append(targetMachines, availableMachines[i])
					found = true
					break
				}
			}
			if !found {
				missingMachines = append(missingMachines, targetName)
			}
		}

		// Log warning for missing machines
		if len(missingMachines) > 0 {
			m.logger.Warn("Some target machines not found in inventory", map[string]interface{}{
				"missing_machines": missingMachines,
				"action":           action.Name,
			})
		}

		return targetMachines
	}

	// If action has tags defined, filter machines by tags
	if len(action.Tags) > 0 {
		var targetMachines []spookytypes.Machine
		for i := range availableMachines {
			if m.machineHasTags(&availableMachines[i], action.Tags) {
				targetMachines = append(targetMachines, availableMachines[i])
			}
		}

		// Log warning if no machines match tags
		if len(targetMachines) == 0 {
			m.logger.Warn("No machines found matching action tags", map[string]interface{}{
				"action_tags": action.Tags,
				"action":      action.Name,
			})
		}

		return targetMachines
	}

	// If no specific targeting, return all available machines
	return availableMachines
}

// getTargetMachinesFromSession uses the session's machine inventory for efficient lookup
func (m *Manager) getTargetMachinesFromSession(action *spookytypesactions.Action, session *spookytypesactions.ActingSession) []spookytypes.Machine {
	// Convert session machine inventory to the expected type
	availableMachines := make([]spookytypes.Machine, len(session.MachineInventory))
	for i := range session.MachineInventory {
		machine := &session.MachineInventory[i]
		availableMachines[i] = spookytypes.Machine(*machine)
	}

	// Use the same logic as getTargetMachines but with session inventory
	if len(action.Machines) > 0 {
		var targetMachines []spookytypes.Machine
		var missingMachines []string

		for _, targetName := range action.Machines {
			found := false
			for i := range availableMachines {
				if availableMachines[i].Hostname == targetName {
					targetMachines = append(targetMachines, availableMachines[i])
					found = true
					break
				}
			}
			if !found {
				missingMachines = append(missingMachines, targetName)
			}
		}

		// Log warning for missing machines
		if len(missingMachines) > 0 {
			m.logger.Warn("Some target machines not found in session inventory", map[string]interface{}{
				"missing_machines": missingMachines,
				"action":           action.Name,
			})
		}

		return targetMachines
	}

	// If action has tags defined, filter machines by tags
	if len(action.Tags) > 0 {
		var targetMachines []spookytypes.Machine
		for i := range availableMachines {
			if m.machineHasTags(&availableMachines[i], action.Tags) {
				targetMachines = append(targetMachines, availableMachines[i])
			}
		}

		// Log warning if no machines match tags
		if len(targetMachines) == 0 {
			m.logger.Warn("No machines found matching action tags in session inventory", map[string]interface{}{
				"action_tags": action.Tags,
				"action":      action.Name,
			})
		}

		return targetMachines
	}

	// If no specific targeting, return all available machines
	return availableMachines
}

// machineHasTags checks if a machine has the specified tags
// Supports both key=value and key-only tag matching
func (m *Manager) machineHasTags(machine *spookytypes.Machine, tags []string) bool {
	if len(tags) == 0 {
		return true
	}

	// Handle case where machine has no tags
	if len(machine.Tags) == 0 {
		return false
	}

	// Check each required tag
	for _, requiredTag := range tags {
		tagMatched := false

		// Check if tag is in key=value format
		if strings.Contains(requiredTag, "=") {
			parts := strings.SplitN(requiredTag, "=", 2)
			if len(parts) == 2 {
				key, value := parts[0], parts[1]
				if machineValue, exists := machine.Tags[key]; exists && machineValue == value {
					tagMatched = true
				}
			}
		} else {
			// Key-only format - check if the key exists in machine tags
			if _, exists := machine.Tags[requiredTag]; exists {
				tagMatched = true
			}
		}

		// If any required tag doesn't match, machine doesn't have all required tags
		if !tagMatched {
			return false
		}
	}

	return true
}

// GetSSHManager returns the SSH manager for authentication testing
func (m *Manager) GetSSHManager() spookyinterfaces.SSHManager {
	return m.sshManager
}

// resolveDependencies resolves action dependencies and creates running order
func (m *Manager) resolveDependencies(actions []*spookytypesactions.Action) (runningOrder [][]string, dependencies map[string][]string, err error) {
	// Build dependency graph
	dependencies = make(map[string][]string)
	for idx := range actions {
		action := actions[idx]
		dependencies[action.Name] = action.Dependencies
	}

	// Check for circular dependencies
	if err := m.checkCircularDependencies(dependencies); err != nil {
		return nil, nil, err
	}

	// Create running order using topological sort
	runningOrder = m.topologicalSort(dependencies)
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
