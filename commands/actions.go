package commands

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"spooky/internal/encryption"
	"spooky/internal/logging"
	"spooky/internal/schemas"
	"spooky/internal/ssh"
	"spooky/internal/utilities"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/spf13/cobra"
)

var (
	actionsCmd = &cobra.Command{
		Use:   "actions",
		Short: "Execute actions across machines",
		Long: `Execute actions across machines defined in your project configuration.

Actions can be of various types:
- command: Execute a single command
- template_deploy: Deploy a template file
- file_sync: Synchronize files
- service_control: Control system services

Examples:
  spooky actions run deploy_webapp
  spooky actions run backup_database --targets web-server-01,db-server-01
  spooky actions list
  spooky actions show deploy_webapp`,
	}

	runActionCmd = &cobra.Command{
		Use:   "run [action-name]",
		Short: "Run a specific action",
		Long: `Run a specific action across target machines.

The action must be defined in your actions.hcl file. You can optionally
override targets using the --targets flag.

Examples:
  spooky actions run deploy_webapp
  spooky actions run backup_database --targets web-server-01,db-server-01
  spooky actions run system_update --dry-run`,
		Args: cobra.ExactArgs(1),
		RunE: runAction,
	}

	listActionsCmd = &cobra.Command{
		Use:   "list",
		Short: "List all available actions",
		Long: `List all actions defined in your actions.hcl file.

Shows action names, descriptions, types, and target information.`,
		RunE: listActions,
	}

	showActionCmd = &cobra.Command{
		Use:   "show [action-name]",
		Short: "Show details of a specific action",
		Long: `Show detailed information about a specific action.

Displays the complete configuration of the action including all fields
and execution parameters.`,
		Args: cobra.ExactArgs(1),
		RunE: showAction,
	}
)

// Action execution flags
var (
	actionTargets []string
	actionDryRun  bool
	actionTimeout int
	actionRetries int
)

func init() {
	// Add subcommands to actions command
	actionsCmd.AddCommand(runActionCmd)
	actionsCmd.AddCommand(listActionsCmd)
	actionsCmd.AddCommand(showActionCmd)

	// Add flags to run command
	runActionCmd.Flags().StringSliceVarP(&actionTargets, "targets", "t", nil, "override action targets (comma-separated list)")
	runActionCmd.Flags().BoolVarP(&actionDryRun, "dry-run", "n", false, "show what would be executed without actually running")
	runActionCmd.Flags().IntVar(&actionTimeout, "timeout", 0, "override action timeout in seconds")
	runActionCmd.Flags().IntVar(&actionRetries, "retries", -1, "override action retry count (-1 to use action default)")

	// Add actions command to root
	RootCmd.AddCommand(actionsCmd)
}

// runAction executes a specific action
func runAction(cmd *cobra.Command, args []string) error {
	actionName := args[0]
	fmt.Printf("🔧 Running action: %s\n", actionName)

	// Load project configuration
	projectConfig, err := loadProjectConfig()
	if err != nil {
		return fmt.Errorf("failed to load project configuration: %w", err)
	}

	// Load actions configuration
	actions, err := loadActionsConfig()
	if err != nil {
		return fmt.Errorf("failed to load actions configuration: %w", err)
	}

	// Find the action
	action, found := findAction(actions, actionName)
	if !found {
		return fmt.Errorf("action '%s' not found in actions configuration", actionName)
	}

	// Load machines configuration
	machines, err := loadMachinesConfig()
	if err != nil {
		return fmt.Errorf("failed to load machines configuration: %w", err)
	}

	// Determine target machines
	targetMachines, err := determineTargetMachines(action, machines, actionTargets)
	if err != nil {
		return fmt.Errorf("failed to determine target machines: %w", err)
	}

	if len(targetMachines) == 0 {
		return fmt.Errorf("no target machines found for action '%s'", actionName)
	}

	fmt.Printf("🎯 Target machines: %s\n", formatMachineList(targetMachines))

	// Load SSH configuration
	sshConfig, err := loadSSHConfig()
	if err != nil {
		return fmt.Errorf("failed to load SSH configuration: %w", err)
	}

	// Create SSH manager with encryption
	ageEncryption, err := encryption.NewAgeEncryption("", "")
	if err != nil {
		logger := logging.GetGlobalLogger()
		logger.Warn("failed to initialize age encryption, continuing without encryption support",
			slog.String("error", err.Error()))
		// Continue with nil encryption - SSH manager will handle this gracefully
	}

	sshManager := ssh.NewSimpleSSHManager(ageEncryption, sshConfig)

	// Create action executor
	executor := NewActionExecutor(sshManager, projectConfig)

	// Override action parameters if specified
	if actionTimeout > 0 {
		action.Timeout = actionTimeout
	}
	if actionRetries >= 0 {
		action.Retries = actionRetries
	}

	// Execute the action
	ctx := context.Background()
	results, err := executor.RunAction(ctx, action, targetMachines, actionDryRun)
	if err != nil {
		return fmt.Errorf("failed to execute action: %w", err)
	}

	// Display results
	displayActionResults(results)

	return nil
}

// listActions lists all available actions
func listActions(cmd *cobra.Command, args []string) error {
	fmt.Println("📋 Available Actions")
	fmt.Println("===================")

	// Load actions configuration
	actions, err := loadActionsConfig()
	if err != nil {
		return fmt.Errorf("failed to load actions configuration: %w", err)
	}

	if len(actions) == 0 {
		fmt.Println("No actions found in actions.hcl")
		return nil
	}

	for _, action := range actions {
		fmt.Printf("\n🔧 %s\n", action.Description)
		fmt.Printf("   Type: %s\n", action.Type)
		if len(action.Targets) > 0 {
			fmt.Printf("   Targets: %s\n", strings.Join(action.Targets, ", "))
		}
		if action.Timeout > 0 {
			fmt.Printf("   Timeout: %ds\n", action.Timeout)
		}
		if action.Retries > 0 {
			fmt.Printf("   Retries: %d\n", action.Retries)
		}
	}

	return nil
}

// showAction shows details of a specific action
func showAction(cmd *cobra.Command, args []string) error {
	actionName := args[0]

	// Load actions configuration
	actions, err := loadActionsConfig()
	if err != nil {
		return fmt.Errorf("failed to load actions configuration: %w", err)
	}

	// Find the action
	action, found := findAction(actions, actionName)
	if !found {
		return fmt.Errorf("action '%s' not found in actions configuration", actionName)
	}

	fmt.Printf("🔧 Action: %s\n", actionName)
	fmt.Println("================")
	fmt.Printf("Description: %s\n", action.Description)
	fmt.Printf("Type: %s\n", action.Type)

	if len(action.Tags) > 0 {
		fmt.Printf("Tags: %s\n", strings.Join(action.Tags, ", "))
	}

	// Show type-specific fields
	switch action.Type {
	case "command":
		fmt.Printf("Command: %s\n", action.Command)

	case "template_deploy":
		fmt.Printf("Source: %s\n", action.Source)
		fmt.Printf("Destination: %s\n", action.Destination)
		fmt.Printf("Validate: %t\n", action.Validate)
		fmt.Printf("Backup: %t\n", action.Backup)
		if action.Permissions != "" {
			fmt.Printf("Permissions: %s\n", action.Permissions)
		}
		if action.Owner != "" {
			fmt.Printf("Owner: %s\n", action.Owner)
		}
		if action.Group != "" {
			fmt.Printf("Group: %s\n", action.Group)
		}
	case "file_sync":
		fmt.Printf("Source: %s\n", action.SyncSource)
		fmt.Printf("Destination: %s\n", action.SyncDestination)
		fmt.Printf("Delete: %t\n", action.SyncDelete)
		fmt.Printf("Preserve: %t\n", action.SyncPreserve)
	case "service_control":
		fmt.Printf("Service: %s\n", action.ServiceName)
		fmt.Printf("Action: %s\n", action.ServiceAction)
	}

	// Show execution configuration
	if len(action.Targets) > 0 {
		fmt.Printf("Targets: %s\n", strings.Join(action.Targets, ", "))
	}
	if action.RunAs != "" {
		fmt.Printf("Run As: %s\n", action.RunAs)
	}
	fmt.Printf("Sudo: %t\n", action.Sudo)
	fmt.Printf("Timeout: %ds\n", action.Timeout)
	fmt.Printf("Retries: %d\n", action.Retries)
	fmt.Printf("Retry Delay: %ds\n", action.RetryDelay)

	// Show conditional execution
	if action.When != "" {
		fmt.Printf("When: %s\n", action.When)
	}
	if action.OnlyIf != "" {
		fmt.Printf("Only If: %s\n", action.OnlyIf)
	}
	if action.Unless != "" {
		fmt.Printf("Unless: %s\n", action.Unless)
	}

	return nil
}

// loadMachinesConfig loads the machines configuration
func loadMachinesConfig() ([]*schemas.MachinesMachineV1, error) {
	// Use the new parsing logic that handles machine names as hostnames
	machinesWithNames, err := getMachinesWithNames()
	if err != nil {
		return nil, err
	}

	// Convert to the expected format
	var machines []*schemas.MachinesMachineV1
	for _, machineWithName := range machinesWithNames {
		machines = append(machines, machineWithName.Machine)
	}

	return machines, nil
}

// loadActionsConfig loads the actions configuration from actions.hcl
func loadActionsConfig() ([]*schemas.ActionsActionV1, error) {
	logger := logging.GetGlobalLogger()
	logger.Debug("loading actions configuration")

	// Look for actions.hcl in current directory
	actionsHCLPath := "actions.hcl"
	if _, err := os.Stat(actionsHCLPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("actions.hcl not found in current directory")
	}

	// Read and parse actions.hcl using the HCL library properly
	content, err := os.ReadFile(actionsHCLPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read actions.hcl: %w", err)
	}

	// Use the validator to parse and validate the actions configuration
	validator := schemas.NewValidator()
	result, err := validator.ValidateHCLContent("actions", string(content))
	if err != nil {
		return nil, fmt.Errorf("failed to validate actions.hcl: %w", err)
	}

	if !result.IsValid {
		return nil, fmt.Errorf("actions.hcl validation failed: %s", result.Errors[0].Message)
	}

	// Parse HCL using the library
	file, diags := hclsyntax.ParseConfig(content, actionsHCLPath, hcl.Pos{Line: 1, Column: 1})
	if diags.HasErrors() {
		return nil, fmt.Errorf("failed to parse actions.hcl: %v", diags)
	}

	// Define the schema for actions block
	schema := &hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{
			{
				Type:       "actions",
				LabelNames: []string{},
			},
		},
	}

	// Extract the actions block
	bodyContent, diags := file.Body.Content(schema)
	if diags.HasErrors() {
		return nil, fmt.Errorf("failed to decode actions block: %v", diags)
	}

	if len(bodyContent.Blocks) == 0 {
		return nil, fmt.Errorf("no actions block found in actions.hcl")
	}

	if len(bodyContent.Blocks) > 1 {
		return nil, fmt.Errorf("multiple actions blocks found in actions.hcl")
	}

	actionsBlock := bodyContent.Blocks[0]

	// Define schema for action blocks inside actions
	actionSchema := &hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{
			{
				Type:       "action",
				LabelNames: []string{"name"},
			},
		},
	}

	// Extract action blocks
	actionContent, diags := actionsBlock.Body.Content(actionSchema)
	if diags.HasErrors() {
		return nil, fmt.Errorf("failed to decode action blocks: %v", diags)
	}

	logger.Debug("found action blocks", slog.Int("count", len(actionContent.Blocks)))

	// Create actions slice
	var actions []*schemas.ActionsActionV1

	// Process each action block
	for _, actionBlock := range actionContent.Blocks {
		actionName := actionBlock.Labels[0]
		logger.Debug("processing action", slog.String("action_name", actionName))

		// Create a new action struct
		action := &schemas.ActionsActionV1{
			Description: "",
			Type:        "",
			Tags:        []string{},
			Targets:     []string{},
			Timeout:     300, // Default timeout
			Retries:     0,   // Default retries
			RetryDelay:  5,   // Default retry delay
		}

		// Define schema for action attributes
		actionAttrSchema := &hcl.BodySchema{
			Attributes: []hcl.AttributeSchema{
				{Name: "description", Required: true},
				{Name: "type", Required: true},
				{Name: "command", Required: false},

				{Name: "targets", Required: false},
				{Name: "timeout", Required: false},
				{Name: "retries", Required: false},
				{Name: "retry_delay", Required: false},
				{Name: "sudo", Required: false},
			},
		}

		// Extract action attributes
		actionAttrContent, diags := actionBlock.Body.Content(actionAttrSchema)
		if diags.HasErrors() {
			return nil, fmt.Errorf("failed to decode action %s attributes: %v", actionName, diags)
		}

		// Extract description
		if descAttr, exists := actionAttrContent.Attributes["description"]; exists {
			var description string
			if diags := gohcl.DecodeExpression(descAttr.Expr, nil, &description); diags.HasErrors() {
				return nil, fmt.Errorf("failed to decode description for action %s: %v", actionName, diags)
			}
			action.Description = description
		}

		// Extract type
		if typeAttr, exists := actionAttrContent.Attributes["type"]; exists {
			var actionType string
			if diags := gohcl.DecodeExpression(typeAttr.Expr, nil, &actionType); diags.HasErrors() {
				return nil, fmt.Errorf("failed to decode type for action %s: %v", actionName, diags)
			}
			action.Type = actionType
		}

		// Extract command
		if cmdAttr, exists := actionAttrContent.Attributes["command"]; exists {
			var command string
			if diags := gohcl.DecodeExpression(cmdAttr.Expr, nil, &command); diags.HasErrors() {
				return nil, fmt.Errorf("failed to decode command for action %s: %v", actionName, diags)
			}
			action.Command = command
		}

		// Extract targets
		if targetsAttr, exists := actionAttrContent.Attributes["targets"]; exists {
			var targets []string
			if diags := gohcl.DecodeExpression(targetsAttr.Expr, nil, &targets); diags.HasErrors() {
				return nil, fmt.Errorf("failed to decode targets for action %s: %v", actionName, diags)
			}
			action.Targets = targets
		}

		// Extract timeout
		if timeoutAttr, exists := actionAttrContent.Attributes["timeout"]; exists {
			var timeout int
			if diags := gohcl.DecodeExpression(timeoutAttr.Expr, nil, &timeout); diags.HasErrors() {
				return nil, fmt.Errorf("failed to decode timeout for action %s: %v", actionName, diags)
			}
			action.Timeout = timeout
		}

		// Extract retries
		if retriesAttr, exists := actionAttrContent.Attributes["retries"]; exists {
			var retries int
			if diags := gohcl.DecodeExpression(retriesAttr.Expr, nil, &retries); diags.HasErrors() {
				return nil, fmt.Errorf("failed to decode retries for action %s: %v", actionName, diags)
			}
			action.Retries = retries
		}

		// Extract retry_delay
		if retryDelayAttr, exists := actionAttrContent.Attributes["retry_delay"]; exists {
			var retryDelay int
			if diags := gohcl.DecodeExpression(retryDelayAttr.Expr, nil, &retryDelay); diags.HasErrors() {
				return nil, fmt.Errorf("failed to decode retry_delay for action %s: %v", actionName, diags)
			}
			action.RetryDelay = retryDelay
		}

		// Extract sudo
		if sudoAttr, exists := actionAttrContent.Attributes["sudo"]; exists {
			var sudo bool
			if diags := gohcl.DecodeExpression(sudoAttr.Expr, nil, &sudo); diags.HasErrors() {
				return nil, fmt.Errorf("failed to decode sudo for action %s: %v", actionName, diags)
			}
			action.Sudo = sudo
		}

		// Store the action name in the description for now (since the struct doesn't have a name field)
		// In a full implementation, we'd add a Name field to the struct
		action.Description = fmt.Sprintf("%s: %s", actionName, action.Description)

		actions = append(actions, action)
		logger.Debug("successfully parsed action",
			slog.String("action_name", actionName),
			slog.String("description", action.Description))
	}

	logger.Debug("loadActionsConfig completed", slog.Int("action_count", len(actions)))
	return actions, nil
}

// findAction finds an action by name
func findAction(actions []*schemas.ActionsActionV1, name string) (*schemas.ActionsActionV1, bool) {
	for _, action := range actions {
		// Check if the action name is in the description (temporary workaround)
		if strings.Contains(action.Description, name+":") {
			return action, true
		}
	}
	return nil, false
}

// determineTargetMachines determines which machines to target for an action
func determineTargetMachines(action *schemas.ActionsActionV1, machines []*schemas.MachinesMachineV1, overrideTargets []string) ([]*schemas.MachinesMachineV1, error) {
	var targetNames []string

	// Use override targets if specified
	if len(overrideTargets) > 0 {
		targetNames = overrideTargets
	} else {
		targetNames = action.Targets
	}

	if len(targetNames) == 0 {
		return nil, fmt.Errorf("no targets specified for action")
	}

	var targetMachines []*schemas.MachinesMachineV1

	// Filter machines by target names
	targetMap := make(map[string]bool)
	for _, target := range targetNames {
		targetMap[target] = true
	}

	for _, machine := range machines {
		// For now, we'll match by hostname (which should be the machine name)
		// In a full implementation, we'd also support groups and tags
		if targetMap[machine.Hostname] {
			targetMachines = append(targetMachines, machine)
		}
	}

	if len(targetMachines) == 0 {
		return nil, fmt.Errorf("no machines found matching targets: %v", targetNames)
	}

	return targetMachines, nil
}

// formatMachineList formats a list of machines for display
func formatMachineList(machines []*schemas.MachinesMachineV1) string {
	var names []string
	for _, machine := range machines {
		names = append(names, machine.Hostname)
	}
	return strings.Join(names, ", ")
}

// displayActionResults displays the results of action execution
func displayActionResults(results []*ActionResult) {
	fmt.Println("\n📊 Action Results")
	fmt.Println("================")

	successCount := 0
	failureCount := 0

	for _, result := range results {
		if result.Success {
			successCount++
			fmt.Printf("✅ %s: %s\n", result.Machine.Hostname, result.Message)
		} else {
			failureCount++
			fmt.Printf("❌ %s: %v\n", result.Machine.Hostname, result.Error)
		}
	}

	fmt.Printf("\nSummary: %d successful, %d failed\n", successCount, failureCount)
}

// ActionResult represents the result of executing an action on a machine
type ActionResult struct {
	Machine *schemas.MachinesMachineV1
	Success bool
	Message string
	Error   string
	Output  string
}

// ActionExecutor handles the execution of actions across machines
type ActionExecutor struct {
	sshManager    *ssh.SimpleSSHManager
	projectConfig *schemas.ProjectV1
}

// NewActionExecutor creates a new action executor
func NewActionExecutor(sshManager *ssh.SimpleSSHManager, projectConfig *schemas.ProjectV1) *ActionExecutor {
	return &ActionExecutor{
		sshManager:    sshManager,
		projectConfig: projectConfig,
	}
}

// ExecuteAction executes an action across target machines
func (ae *ActionExecutor) RunAction(ctx context.Context, action *schemas.ActionsActionV1, machines []*schemas.MachinesMachineV1, dryRun bool) ([]*ActionResult, error) {
	var results []*ActionResult

	for _, machine := range machines {
		result := &ActionResult{
			Machine: machine,
		}

		if dryRun {
			result.Success = true
			result.Message = fmt.Sprintf("Would execute %s action", action.Type)
		} else {
			// Execute the action based on its type
			err := ae.runActionOnMachine(ctx, action, machine)
			if err != nil {
				result.Success = false
				result.Error = err.Error()
			} else {
				result.Success = true
				result.Message = fmt.Sprintf("Successfully executed %s action", action.Type)
			}
		}

		results = append(results, result)
	}

	return results, nil
}

// executeActionOnMachine executes an action on a specific machine
func (ae *ActionExecutor) runActionOnMachine(ctx context.Context, action *schemas.ActionsActionV1, machine *schemas.MachinesMachineV1) error {
	// Check conditions first
	if err := ae.checkConditions(ctx, action, machine); err != nil {
		return fmt.Errorf("condition check failed: %w", err)
	}

	// Execute based on action type
	switch action.Type {
	case "command":
		return ae.runCommandAction(ctx, action, machine)

	case "template_deploy":
		return ae.runTemplateDeployAction(ctx, action, machine)
	case "file_sync":
		return ae.runFileSyncAction(ctx, action, machine)
	case "service_control":
		return ae.runServiceControlAction(ctx, action, machine)
	default:
		return fmt.Errorf("unsupported action type: %s", action.Type)
	}
}

// checkConditions checks if the action should be executed based on conditions
func (ae *ActionExecutor) checkConditions(ctx context.Context, action *schemas.ActionsActionV1, machine *schemas.MachinesMachineV1) error {
	// Check OnlyIf condition
	if action.OnlyIf != "" {
		result, err := utilities.RunCommand(ctx, machine.Hostname, action.OnlyIf, machine, ae.sshManager, 0, 0, 0)
		if err != nil {
			return fmt.Errorf("only_if condition failed: %w", err)
		}
		if result.ExitCode != 0 {
			return fmt.Errorf("only_if condition returned non-zero exit code: %d", result.ExitCode)
		}
	}

	// Check Unless condition
	if action.Unless != "" {
		result, err := utilities.RunCommand(ctx, machine.Hostname, action.Unless, machine, ae.sshManager, 0, 0, 0)
		if err != nil {
			return fmt.Errorf("unless condition failed: %w", err)
		}
		if result.ExitCode == 0 {
			return fmt.Errorf("unless condition returned zero exit code")
		}
	}

	return nil
}

// executeCommandAction executes a command action
func (ae *ActionExecutor) runCommandAction(ctx context.Context, action *schemas.ActionsActionV1, machine *schemas.MachinesMachineV1) error {
	if action.Command == "" {
		return fmt.Errorf("command action requires a command field")
	}

	// Build the command with sudo if needed
	command := action.Command
	if action.Sudo {
		command = "sudo " + command
	}

	// Execute with unified command runner
	result, err := utilities.RunCommand(ctx, machine.Hostname, command, machine, ae.sshManager, time.Duration(action.Timeout)*time.Second, action.Retries, time.Duration(action.RetryDelay)*time.Second)
	if err != nil {
		return fmt.Errorf("command execution failed: %w", err)
	}

	if result.ExitCode != 0 {
		return fmt.Errorf("command returned non-zero exit code: %d", result.ExitCode)
	}

	return nil
}

// executeTemplateDeployAction executes a template deploy action
func (ae *ActionExecutor) runTemplateDeployAction(ctx context.Context, action *schemas.ActionsActionV1, machine *schemas.MachinesMachineV1) error {
	// This is a placeholder implementation
	// In a full implementation, we'd:
	// 1. Render the template with variables
	// 2. Upload the rendered file to the destination
	// 3. Set permissions and ownership if specified
	return fmt.Errorf("template_deploy action not yet implemented")
}

// executeFileSyncAction executes a file sync action
func (ae *ActionExecutor) runFileSyncAction(ctx context.Context, action *schemas.ActionsActionV1, machine *schemas.MachinesMachineV1) error {
	// This is a placeholder implementation
	// In a full implementation, we'd use rsync or similar to sync files
	return fmt.Errorf("file_sync action not yet implemented")
}

// executeServiceControlAction executes a service control action
func (ae *ActionExecutor) runServiceControlAction(ctx context.Context, action *schemas.ActionsActionV1, machine *schemas.MachinesMachineV1) error {
	if action.ServiceName == "" || action.ServiceAction == "" {
		return fmt.Errorf("service_control action requires service_name and service_action fields")
	}

	command := fmt.Sprintf("systemctl %s %s", action.ServiceAction, action.ServiceName)
	if action.Sudo {
		command = "sudo " + command
	}

	// Execute with unified command runner
	result, err := utilities.RunCommand(ctx, machine.Hostname, command, machine, ae.sshManager, time.Duration(action.Timeout)*time.Second, action.Retries, time.Duration(action.RetryDelay)*time.Second)
	if err != nil {
		return fmt.Errorf("command execution failed: %w", err)
	}

	if result.ExitCode != 0 {
		return fmt.Errorf("command returned non-zero exit code: %d", result.ExitCode)
	}

	return nil
}
