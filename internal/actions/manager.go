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
	spookytypesmachines "spooky/internal/types/machines"
	spookytypesssh "spooky/internal/types/ssh"

	"github.com/hashicorp/hcl/v2/hclsimple"
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
	for _, action := range actions {
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
func (m *Manager) runActionStep(ctx context.Context, session *spookytypesactions.ActingSession, actionNames []string, allActions []*spookytypesactions.Action, plan *spookytypesactions.ActionPlan) ([]spookytypes.ActingResult, error) {
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

// runCommandAction runs a command action
func (m *Manager) runCommandAction(ctx context.Context, action *spookytypesactions.Action, machine *spookytypes.Machine, result *spookytypesactions.ActingResult) error {
	m.logger.Debug("Running command action via SSH", map[string]interface{}{
		"action":  action.Name,
		"machine": machine.Hostname,
		"command": action.CommandString,
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

	// Determine command to run
	var commandToRun string
	var args []string
	var workingDir string
	var environment map[string]string
	var timeout time.Duration

	if action.CommandString != "" {
		// Use CommandString if provided
		commandToRun = action.CommandString
		workingDir = action.WorkingDir
		environment = action.Environment
		timeout = time.Duration(action.Timeout) * time.Second
	} else if action.Command != nil {
		// Fall back to Command config if available
		commandToRun = action.Command.Command
		args = action.Command.Args
		workingDir = action.Command.WorkingDir
		if workingDir == "" {
			workingDir = action.WorkingDir
		}
		environment = action.Command.Environment
		if environment == nil {
			environment = action.Environment
		}
		if action.Command.Timeout > 0 {
			timeout = time.Duration(action.Command.Timeout) * time.Second
		} else {
			timeout = time.Duration(action.Timeout) * time.Second
		}
	} else {
		return fmt.Errorf("no command specified for action %s", action.Name)
	}

	// Create SSH command
	sshCommand := &spookytypes.SSHCommand{
		Command:       commandToRun,
		Args:          args,
		WorkingDir:    workingDir,
		Environment:   environment,
		Timeout:       timeout,
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
func (m *Manager) runScriptAction(ctx context.Context, action *spookytypesactions.Action, machine *spookytypes.Machine, result *spookytypesactions.ActingResult) error {
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
func (m *Manager) runTemplateAction(ctx context.Context, action *spookytypesactions.Action, machine *spookytypes.Machine, result *spookytypesactions.ActingResult) error {
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
func (m *Manager) runFileCopyAction(ctx context.Context, action *spookytypesactions.Action, machine *spookytypes.Machine, result *spookytypesactions.ActingResult) error {
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
func (m *Manager) runServiceControlAction(ctx context.Context, action *spookytypesactions.Action, machine *spookytypes.Machine, result *spookytypesactions.ActingResult) error {
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
	for i, machine := range session.MachineInventory {
		availableMachines[i] = spookytypes.Machine(machine)
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

// populateSessionWithMachineInventory loads and populates the session with machine inventory
func (m *Manager) populateSessionWithMachineInventory(ctx context.Context, session *spookytypesactions.ActingSession) error {
	if session == nil {
		return fmt.Errorf("session cannot be nil")
	}

	if session.ProjectPath == "" {
		return fmt.Errorf("session project path is required for machine inventory loading")
	}

	// Load machines from the project
	machines, err := m.loadMachinesFromProject(ctx, session.ProjectPath)
	if err != nil {
		return fmt.Errorf("failed to load machine inventory: %w", err)
	}

	// Convert to session machine type and populate inventory
	session.MachineInventory = make([]spookytypesmachines.Machine, len(machines))
	session.MachineCache = make(map[string]*spookytypesmachines.Machine)

	for i, machine := range machines {
		// Convert spookytypes.Machine to spookytypesmachines.Machine
		sessionMachine := spookytypesmachines.Machine(machine)
		session.MachineInventory[i] = sessionMachine
		session.MachineCache[machine.Hostname] = &session.MachineInventory[i]
	}

	m.logger.Debug("Populated session with machine inventory", map[string]interface{}{
		"session_id":    session.SessionID,
		"project_path":  session.ProjectPath,
		"machine_count": len(session.MachineInventory),
		"cache_entries": len(session.MachineCache),
	})

	return nil
}

// loadMachinesFromProject loads machines from a project path
// This is a helper method that would typically use the machines integration
func (m *Manager) loadMachinesFromProject(_ context.Context, projectPath string) ([]spookytypes.Machine, error) {
	// This is a placeholder implementation
	// In a real implementation, this would use the machines integration
	// For now, return empty slice to avoid compilation errors
	return []spookytypes.Machine{}, nil
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
	for _, action := range actions {
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
