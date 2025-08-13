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
	spookytypesssh "spooky/internal/types/ssh"
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

// setFilePermissions sets file permissions on a remote machine via SSH
func (m *Manager) setFilePermissions(ctx context.Context, session *spookytypes.Session, machine spookytypes.Machine, filePath, permissions string, action *spookytypesactions.Action) error {
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

// runCommandAction runs a command action
func (m *Manager) runCommandAction(ctx context.Context, action *spookytypesactions.Action, machine spookytypes.Machine, result *spookytypesactions.ActingResult) error {
	m.logger.Debug("Running command action via SSH", map[string]interface{}{
		"action":  action.Name,
		"machine": machine.Hostname,
		"command": action.Command,
	})

	// Create SSH connection request
	connectionRequest := &spookytypes.ConnectionRequest{
		Host:     machine.Hostname,
		Port:     machine.Port,
		User:     machine.User,
		Timeout:  time.Duration(action.Timeout) * time.Second,
		KeyPath:  machine.KeyFile,
		Password: machine.Password,
	}

	// Establish SSH connection
	connectionResult, err := m.sshManager.Connect(ctx, connectionRequest)
	if err != nil {
		return fmt.Errorf("failed to establish SSH connection to %s: %w", machine.Hostname, err)
	}

	if !connectionResult.Success {
		return fmt.Errorf("SSH connection failed to %s: %s", machine.Hostname, connectionResult.Error)
	}

	// Create SSH session
	session, err := m.sshManager.CreateSession(ctx, connectionResult.Connection)
	if err != nil {
		return fmt.Errorf("failed to create SSH session for %s: %w", machine.Hostname, err)
	}

	// Create SSH command
	sshCommand := &spookytypes.SSHCommand{
		Command:       action.Command.Command,
		Args:          action.Command.Args,
		WorkingDir:    action.WorkingDir,
		Environment:   action.Environment,
		Timeout:       time.Duration(action.Timeout) * time.Second,
		CaptureOutput: true,
	}

	// Run command via SSH
	commandResult, err := m.sshManager.RunCommand(ctx, session, sshCommand)
	if err != nil {
		return fmt.Errorf("failed to run command on %s: %w", machine.Hostname, err)
	}

	// Update result with SSH command results
	result.ExitCode = commandResult.ExitCode
	result.Stdout = commandResult.Stdout
	result.Stderr = commandResult.Stderr
	result.StartTime = commandResult.StartTime
	result.EndTime = commandResult.EndTime
	result.Duration = commandResult.Duration

	if !commandResult.Success {
		result.Status = "failure"
		result.ErrorType = "command_execution"
		result.Error = commandResult.Error
	}

	m.logger.Info("Command action completed", map[string]interface{}{
		"action":    action.Name,
		"machine":   machine.Hostname,
		"exit_code": commandResult.ExitCode,
		"duration":  commandResult.Duration,
		"success":   commandResult.Success,
	})

	return nil
}

// runScriptAction runs a script action
func (m *Manager) runScriptAction(ctx context.Context, action *spookytypesactions.Action, machine spookytypes.Machine, result *spookytypesactions.ActingResult) error {
	m.logger.Debug("Running script action via SSH", map[string]interface{}{
		"action":  action.Name,
		"machine": machine.Hostname,
		"script":  action.Script,
	})

	// Create SSH connection request
	connectionRequest := &spookytypes.ConnectionRequest{
		Host:     machine.Hostname,
		Port:     machine.Port,
		User:     machine.User,
		Timeout:  time.Duration(action.Timeout) * time.Second,
		KeyPath:  machine.KeyFile,
		Password: machine.Password,
	}

	// Establish SSH connection
	connectionResult, err := m.sshManager.Connect(ctx, connectionRequest)
	if err != nil {
		return fmt.Errorf("failed to establish SSH connection to %s: %w", machine.Hostname, err)
	}

	if !connectionResult.Success {
		return fmt.Errorf("SSH connection failed to %s: %s", machine.Hostname, connectionResult.Error)
	}

	// Create SSH session
	session, err := m.sshManager.CreateSession(ctx, connectionResult.Connection)
	if err != nil {
		return fmt.Errorf("failed to create SSH session for %s: %w", machine.Hostname, err)
	}

	// Determine shell based on action configuration
	shell := action.Script.Shell
	if shell == "" {
		shell = "/bin/bash"
	}

	// Create script command
	scriptCommand := fmt.Sprintf("%s << 'EOF'\n%s\nEOF", shell, action.Script.Script)

	// Create SSH command
	sshCommand := &spookytypes.SSHCommand{
		Command:       scriptCommand,
		WorkingDir:    action.WorkingDir,
		Environment:   action.Environment,
		Timeout:       time.Duration(action.Timeout) * time.Second,
		CaptureOutput: true,
	}

	// Run script via SSH
	commandResult, err := m.sshManager.RunCommand(ctx, session, sshCommand)
	if err != nil {
		return fmt.Errorf("failed to run script on %s: %w", machine.Hostname, err)
	}

	// Update result with SSH command results
	result.ExitCode = commandResult.ExitCode
	result.Stdout = commandResult.Stdout
	result.Stderr = commandResult.Stderr
	result.StartTime = commandResult.StartTime
	result.EndTime = commandResult.EndTime
	result.Duration = commandResult.Duration

	if !commandResult.Success {
		result.Status = "failure"
		result.ErrorType = "script_execution"
		result.Error = commandResult.Error
	}

	m.logger.Info("Script action completed", map[string]interface{}{
		"action":    action.Name,
		"machine":   machine.Hostname,
		"exit_code": commandResult.ExitCode,
		"duration":  commandResult.Duration,
		"success":   commandResult.Success,
	})

	return nil
}

// runTemplateAction runs a template deployment action
func (m *Manager) runTemplateAction(ctx context.Context, action *spookytypesactions.Action, machine spookytypes.Machine, result *spookytypesactions.ActingResult) error {
	m.logger.Debug("Running template deployment action via SSH", map[string]interface{}{
		"action":   action.Name,
		"machine":  machine.Hostname,
		"template": action.Template,
	})

	if action.Template == nil {
		return fmt.Errorf("template action requires template configuration")
	}

	// Create SSH connection request
	connectionRequest := &spookytypes.ConnectionRequest{
		Host:     machine.Hostname,
		Port:     machine.Port,
		User:     machine.User,
		Timeout:  time.Duration(action.Timeout) * time.Second,
		KeyPath:  machine.KeyFile,
		Password: machine.Password,
	}

	// Establish SSH connection
	connectionResult, err := m.sshManager.Connect(ctx, connectionRequest)
	if err != nil {
		return fmt.Errorf("failed to establish SSH connection to %s: %w", machine.Hostname, err)
	}

	if !connectionResult.Success {
		return fmt.Errorf("SSH connection failed to %s: %s", machine.Hostname, connectionResult.Error)
	}

	// Create SSH session
	session, err := m.sshManager.CreateSession(ctx, connectionResult.Connection)
	if err != nil {
		return fmt.Errorf("failed to create SSH session for %s: %w", machine.Hostname, err)
	}

	// For template deployment, we need to:
	// 1. Transfer the template file to the remote machine
	// 2. Execute the template on the remote machine

	// Create file transfer for template
	transfer := &spookytypes.FileTransfer{
		LocalPath:  action.Template.Source,
		RemotePath: action.Template.Destination,
		Direction:  spookytypesssh.TransferDirectionUpload,
		Mode:       spookytypesssh.TransferModeSFTP,
	}

	// Transfer template file
	transferResult, err := m.sshManager.TransferFile(ctx, session, transfer)
	if err != nil {
		return fmt.Errorf("failed to transfer template file to %s: %w", machine.Hostname, err)
	}

	if !transferResult.Success {
		return fmt.Errorf("template file transfer failed to %s: %s", machine.Hostname, transferResult.Error)
	}

	// Set file permissions if specified
	if action.Template.Permissions != "" {
		err = m.setFilePermissions(ctx, session, machine, action.Template.Destination, action.Template.Permissions, action)
		if err != nil {
			return fmt.Errorf("failed to set template file permissions on %s: %w", machine.Hostname, err)
		}
	}

	// Execute template if it's executable (we'll assume it's executable for now)
	// In a real implementation, you might check file permissions or have a separate flag
	sshCommand := &spookytypes.SSHCommand{
		Command:       action.Template.Destination,
		WorkingDir:    action.WorkingDir,
		Environment:   action.Environment,
		Timeout:       time.Duration(action.Timeout) * time.Second,
		CaptureOutput: true,
	}

	commandResult, err := m.sshManager.RunCommand(ctx, session, sshCommand)
	if err != nil {
		return fmt.Errorf("failed to execute template on %s: %w", machine.Hostname, err)
	}

	// Update result with SSH command results
	result.ExitCode = commandResult.ExitCode
	result.Stdout = commandResult.Stdout
	result.Stderr = commandResult.Stderr
	result.StartTime = commandResult.StartTime
	result.EndTime = commandResult.EndTime
	result.Duration = commandResult.Duration

	if !commandResult.Success {
		result.Status = "failure"
		result.ErrorType = "template_execution"
		result.Error = commandResult.Error
	}

	m.logger.Info("Template action completed", map[string]interface{}{
		"action":      action.Name,
		"machine":     machine.Hostname,
		"template":    action.Template.Source,
		"destination": action.Template.Destination,
		"success":     result.Status == "success",
	})

	return nil
}

// runFileCopyAction runs a file copy action
func (m *Manager) runFileCopyAction(ctx context.Context, action *spookytypesactions.Action, machine spookytypes.Machine, result *spookytypesactions.ActingResult) error {
	m.logger.Debug("Running file copy action via SSH", map[string]interface{}{
		"action":    action.Name,
		"machine":   machine.Hostname,
		"file_copy": action.FileCopy,
	})

	if action.FileCopy == nil {
		return fmt.Errorf("file copy action requires file copy configuration")
	}

	// Create SSH connection request
	connectionRequest := &spookytypes.ConnectionRequest{
		Host:     machine.Hostname,
		Port:     machine.Port,
		User:     machine.User,
		Timeout:  time.Duration(action.Timeout) * time.Second,
		KeyPath:  machine.KeyFile,
		Password: machine.Password,
	}

	// Establish SSH connection
	connectionResult, err := m.sshManager.Connect(ctx, connectionRequest)
	if err != nil {
		return fmt.Errorf("failed to establish SSH connection to %s: %w", machine.Hostname, err)
	}

	if !connectionResult.Success {
		return fmt.Errorf("SSH connection failed to %s: %s", machine.Hostname, connectionResult.Error)
	}

	// Create SSH session
	session, err := m.sshManager.CreateSession(ctx, connectionResult.Connection)
	if err != nil {
		return fmt.Errorf("failed to create SSH session for %s: %w", machine.Hostname, err)
	}

	// Create file transfer
	transfer := &spookytypes.FileTransfer{
		LocalPath:  action.FileCopy.Source,
		RemotePath: action.FileCopy.Destination,
		Direction:  spookytypesssh.TransferDirectionUpload,
		Mode:       spookytypesssh.TransferModeSFTP,
	}

	// Transfer file
	transferResult, err := m.sshManager.TransferFile(ctx, session, transfer)
	if err != nil {
		return fmt.Errorf("failed to transfer file to %s: %w", machine.Hostname, err)
	}

	if !transferResult.Success {
		return fmt.Errorf("file transfer failed to %s: %s", machine.Hostname, transferResult.Error)
	}

	// Set file permissions if specified
	if action.FileCopy.Permissions != "" {
		err = m.setFilePermissions(ctx, session, machine, action.FileCopy.Destination, action.FileCopy.Permissions, action)
		if err != nil {
			return fmt.Errorf("failed to set file permissions on %s: %w", machine.Hostname, err)
		}
	}

	// Set file ownership if specified
	if action.FileCopy.Owner != "" || action.FileCopy.Group != "" {
		chownCommand := fmt.Sprintf("chown %s:%s %s",
			action.FileCopy.Owner,
			action.FileCopy.Group,
			action.FileCopy.Destination)

		sshCommand := &spookytypes.SSHCommand{
			Command:       chownCommand,
			WorkingDir:    action.WorkingDir,
			Environment:   action.Environment,
			Timeout:       time.Duration(action.Timeout) * time.Second,
			CaptureOutput: true,
		}

		commandResult, err := m.sshManager.RunCommand(ctx, session, sshCommand)
		if err != nil {
			return fmt.Errorf("failed to set file ownership on %s: %w", machine.Hostname, err)
		}

		if !commandResult.Success {
			return fmt.Errorf("failed to set file ownership on %s: %s", machine.Hostname, commandResult.Error)
		}
	}

	// Update result
	result.Stdout = fmt.Sprintf("File copied successfully from %s to %s",
		action.FileCopy.Source, action.FileCopy.Destination)
	result.StartTime = time.Now()
	result.EndTime = time.Now()

	m.logger.Info("File copy action completed", map[string]interface{}{
		"action":      action.Name,
		"machine":     machine.Hostname,
		"source":      action.FileCopy.Source,
		"destination": action.FileCopy.Destination,
		"success":     true,
	})

	return nil
}

// runServiceControlAction runs a service control action
func (m *Manager) runServiceControlAction(ctx context.Context, action *spookytypesactions.Action, machine spookytypes.Machine, result *spookytypesactions.ActingResult) error {
	m.logger.Debug("Running service control action via SSH", map[string]interface{}{
		"action":          action.Name,
		"machine":         machine.Hostname,
		"service_control": action.ServiceControl,
	})

	if action.ServiceControl == nil {
		return fmt.Errorf("service control action requires service control configuration")
	}

	// Create SSH connection request
	connectionRequest := &spookytypes.ConnectionRequest{
		Host:     machine.Hostname,
		Port:     machine.Port,
		User:     machine.User,
		Timeout:  time.Duration(action.Timeout) * time.Second,
		KeyPath:  machine.KeyFile,
		Password: machine.Password,
	}

	// Establish SSH connection
	connectionResult, err := m.sshManager.Connect(ctx, connectionRequest)
	if err != nil {
		return fmt.Errorf("failed to establish SSH connection to %s: %w", machine.Hostname, err)
	}

	if !connectionResult.Success {
		return fmt.Errorf("SSH connection failed to %s: %s", machine.Hostname, connectionResult.Error)
	}

	// Create SSH session
	session, err := m.sshManager.CreateSession(ctx, connectionResult.Connection)
	if err != nil {
		return fmt.Errorf("failed to create SSH session for %s: %w", machine.Hostname, err)
	}

	// Determine service control command based on action
	var serviceCommand string
	switch action.ServiceControl.Action {
	case "start":
		serviceCommand = fmt.Sprintf("systemctl start %s", action.ServiceControl.Service)
	case "stop":
		serviceCommand = fmt.Sprintf("systemctl stop %s", action.ServiceControl.Service)
	case "restart":
		serviceCommand = fmt.Sprintf("systemctl restart %s", action.ServiceControl.Service)
	case "reload":
		serviceCommand = fmt.Sprintf("systemctl reload %s", action.ServiceControl.Service)
	case "enable":
		serviceCommand = fmt.Sprintf("systemctl enable %s", action.ServiceControl.Service)
	case "disable":
		serviceCommand = fmt.Sprintf("systemctl disable %s", action.ServiceControl.Service)
	case "status":
		serviceCommand = fmt.Sprintf("systemctl status %s", action.ServiceControl.Service)
	default:
		return fmt.Errorf("unsupported service control action: %s", action.ServiceControl.Action)
	}

	// Add sudo if required (using systemd flag as proxy for sudo requirement)
	if action.ServiceControl.Systemd {
		serviceCommand = fmt.Sprintf("sudo %s", serviceCommand)
	}

	// Create SSH command
	sshCommand := &spookytypes.SSHCommand{
		Command:       serviceCommand,
		WorkingDir:    action.WorkingDir,
		Environment:   action.Environment,
		Timeout:       time.Duration(action.Timeout) * time.Second,
		CaptureOutput: true,
	}

	// Run service control command via SSH
	commandResult, err := m.sshManager.RunCommand(ctx, session, sshCommand)
	if err != nil {
		return fmt.Errorf("failed to run service control command on %s: %w", machine.Hostname, err)
	}

	// Update result with SSH command results
	result.ExitCode = commandResult.ExitCode
	result.Stdout = commandResult.Stdout
	result.Stderr = commandResult.Stderr
	result.StartTime = commandResult.StartTime
	result.EndTime = commandResult.EndTime
	result.Duration = commandResult.Duration

	if !commandResult.Success {
		result.Status = "failure"
		result.ErrorType = "service_control"
		result.Error = commandResult.Error
	}

	m.logger.Info("Service control action completed", map[string]interface{}{
		"action":      action.Name,
		"machine":     machine.Hostname,
		"service":     action.ServiceControl.Service,
		"action_type": action.ServiceControl.Action,
		"exit_code":   commandResult.ExitCode,
		"duration":    commandResult.Duration,
		"success":     commandResult.Success,
	})

	return nil
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
