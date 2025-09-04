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
	"spooky/internal/sync"
	"spooky/internal/templates"
	"spooky/internal/utilities"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/spf13/cobra"
	"github.com/zclconf/go-cty/cty"
)

var (
	actionsCmd = &cobra.Command{
		Use:   schemas.ResourceTypeActions,
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

// runAction executes a specific action across target machines.
// It performs comprehensive action execution with proper error handling and logging.
//
// Parameters:
//   - cmd: Cobra command instance
//   - args: Command arguments (first argument should be action name)
//
// Returns:
//   - error: Any error that occurred during action execution
//
// The function:
//  1. Validates input parameters
//  2. Loads all required configurations
//  3. Determines target machines
//  4. Creates necessary managers and executors
//  5. Executes the action with proper error handling
//  6. Displays comprehensive results
func runAction(cmd *cobra.Command, args []string) error {
	// Input validation
	if len(args) == 0 {
		return fmt.Errorf("action name is required")
	}

	actionName := args[0]
	if actionName == "" {
		return fmt.Errorf("action name cannot be empty")
	}

	logger := logging.GetGlobalLogger()
	logger.Info("🔧 Running action", slog.String("action_name", actionName))

	// Load project configuration
	projectConfig, err := loadProjectConfig()
	if err != nil {
		return utilities.WrapErrorf(err, "failed to load project configuration")
	}

	// Load actions configuration
	actions, err := loadActionsConfig()
	if err != nil {
		return utilities.WrapErrorf(err, "failed to load actions configuration")
	}

	// Find the action
	action, found := findAction(actions, actionName)
	if !found {
		return fmt.Errorf("action '%s' not found in actions configuration", actionName)
	}

	// Load machines configuration
	machines, err := loadMachinesConfig()
	if err != nil {
		return utilities.WrapErrorf(err, "failed to load machines configuration")
	}

	// Determine target machines
	targetMachines, err := determineTargetMachines(action, machines, actionTargets)
	if err != nil {
		return utilities.WrapErrorf(err, "failed to determine target machines")
	}

	if len(targetMachines) == 0 {
		return fmt.Errorf("no target machines found for action '%s'", actionName)
	}

	logger.Info("🎯 Target machines", slog.String("machines", formatMachineList(targetMachines)))

	// Load SSH configuration
	sshConfig, err := loadSSHConfig()
	if err != nil {
		return utilities.WrapErrorf(err, "failed to load SSH configuration")
	}

	// Create SSH manager with encryption
	ageEncryption, err := encryption.NewAgeEncryption("", "")
	if err != nil {
		logger.Warn("failed to initialize age encryption, continuing without encryption support",
			slog.String("error", err.Error()))
		// Continue with nil encryption - SSH manager will handle this gracefully
	}

	sshManager := ssh.NewSSHManager(ageEncryption, sshConfig)

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
	logger := logging.GetGlobalLogger()

	logger.Info("📋 Available Actions")
	logger.Info("===================")

	// Load actions configuration
	actions, err := loadActionsConfig()
	if err != nil {
		return utilities.WrapErrorf(err, "failed to load actions configuration")
	}

	if len(actions) == 0 {
		logger.Info("No actions found in actions.hcl")
		return nil
	}

	for _, action := range actions {
		logger.Info("🔧 Action details",
			slog.String("description", action.Description),
			slog.String("type", action.Type))

		if len(action.Targets) > 0 {
			logger.Info("   Targets", slog.String("targets", strings.Join(action.Targets, ", ")))
		}
		if action.Timeout > 0 {
			logger.Info("   Timeout", slog.Int("timeout_seconds", action.Timeout))
		}
		if action.Retries > 0 {
			logger.Info("   Retries", slog.Int("retries", action.Retries))
		}
	}

	return nil
}

// showAction shows details of a specific action
func showAction(cmd *cobra.Command, args []string) error {
	// Input validation
	if len(args) == 0 {
		return fmt.Errorf("action name is required")
	}

	actionName := args[0]
	if actionName == "" {
		return fmt.Errorf("action name cannot be empty")
	}

	logger := logging.GetGlobalLogger()

	// Load actions configuration
	actions, err := loadActionsConfig()
	if err != nil {
		return utilities.WrapErrorf(err, "failed to load actions configuration")
	}

	// Find the action
	action, found := findAction(actions, actionName)
	if !found {
		return fmt.Errorf("action '%s' not found in actions configuration", actionName)
	}

	// Display action details
	displayActionBasicInfo(logger, actionName, action)
	displayActionTypeSpecificInfo(logger, action)
	displayActionExecutionConfig(logger, action)
	displayActionConditionalExecution(logger, action)

	return nil
}

// displayActionBasicInfo displays basic action information
func displayActionBasicInfo(logger *logging.Logger, actionName string, action *schemas.ActionsActionV1) {
	logger.Info("🔧 Action details",
		slog.String("action_name", actionName),
		slog.String("description", action.Description),
		slog.String("type", action.Type))

	if len(action.Tags) > 0 {
		logger.Info("Action tags", slog.String("tags", strings.Join(action.Tags, ", ")))
	}
}

// displayActionTypeSpecificInfo displays type-specific action information
func displayActionTypeSpecificInfo(logger *logging.Logger, action *schemas.ActionsActionV1) {
	switch action.Type {
	case "command":
		logger.Info("Command action details", slog.String("command", action.Command))

	case "template_deploy":
		displayTemplateDeployInfo(logger, action)
	case "file_sync":
		displayFileSyncInfo(logger, action)
	case "service_control":
		displayServiceControlInfo(logger, action)
	}
}

// displayTemplateDeployInfo displays template deploy specific information
func displayTemplateDeployInfo(logger *logging.Logger, action *schemas.ActionsActionV1) {
	logger.Info("Template deploy action details",
		slog.String("source", action.Source),
		slog.String("destination", action.Destination),
		slog.Bool("validate", action.Validate),
		slog.Bool("backup", action.Backup))

	if action.Permissions != "" {
		logger.Info("Template permissions", slog.String("permissions", action.Permissions))
	}
	if action.Owner != "" {
		logger.Info("Template owner", slog.String("owner", action.Owner))
	}
	if action.Group != "" {
		logger.Info("Template group", slog.String("group", action.Group))
	}
}

// displayFileSyncInfo displays file sync specific information
func displayFileSyncInfo(logger *logging.Logger, action *schemas.ActionsActionV1) {
	logger.Info("File sync action details",
		slog.String("source", action.SyncSource),
		slog.String("destination", action.SyncDestination),
		slog.Bool("delete", action.SyncDelete),
		slog.Bool("preserve", action.SyncPreserve))
}

// displayServiceControlInfo displays service control specific information
func displayServiceControlInfo(logger *logging.Logger, action *schemas.ActionsActionV1) {
	logger.Info("Service control action details",
		slog.String("service", action.ServiceName),
		slog.String("action", action.ServiceAction))
}

// displayActionExecutionConfig displays action execution configuration
func displayActionExecutionConfig(logger *logging.Logger, action *schemas.ActionsActionV1) {
	if len(action.Targets) > 0 {
		logger.Info("Action targets", slog.String("targets", strings.Join(action.Targets, ", ")))
	}
	if action.RunAs != "" {
		logger.Info("Action run as", slog.String("run_as", action.RunAs))
	}
	logger.Info("Action sudo configuration", slog.Bool("sudo", action.Sudo))
	logger.Info("Action timeout configuration",
		slog.Int("timeout_seconds", action.Timeout),
		slog.Int("retries", action.Retries),
		slog.Int("retry_delay_seconds", action.RetryDelay))
}

// displayActionConditionalExecution displays action conditional execution settings
func displayActionConditionalExecution(logger *logging.Logger, action *schemas.ActionsActionV1) {
	if action.When != "" {
		logger.Info("Action when condition", slog.String("when", action.When))
	}
	if action.OnlyIf != "" {
		logger.Info("Action only_if condition", slog.String("only_if", action.OnlyIf))
	}
	if action.Unless != "" {
		logger.Info("Action unless condition", slog.String("unless", action.Unless))
	}
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

	// Use the simplified validator to parse and validate the actions configuration
	validator := schemas.NewSimpleValidator()
	result, err := validator.ValidateHCLContent(schemas.ResourceTypeActions, string(content))
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

	// Extract actions block
	actionsBlock, err := extractActionsBlock(file)
	if err != nil {
		return nil, err
	}

	// Extract action blocks
	actionContent, err := extractActionBlocks(actionsBlock)
	if err != nil {
		return nil, err
	}

	logger.Debug("found action blocks", slog.Int("count", len(actionContent.Blocks)))

	// Process action blocks
	actions, err := processActionBlocks(actionContent)
	if err != nil {
		return nil, err
	}

	logger.Debug("loadActionsConfig completed", slog.Int("action_count", len(actions)))
	return actions, nil
}

// extractActionsBlock extracts the actions block from the parsed HCL file
func extractActionsBlock(file *hcl.File) (*hcl.Block, error) {
	return extractResourceBlock(file, schemas.ResourceTypeActions, "actions.hcl")
}

// extractActionBlocks extracts action blocks from the actions block
func extractActionBlocks(actionsBlock *hcl.Block) (*hcl.BodyContent, error) {
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

	return actionContent, nil
}

// processActionBlocks processes all action blocks and returns the parsed actions
func processActionBlocks(actionContent *hcl.BodyContent) ([]*schemas.ActionsActionV1, error) {
	var actions []*schemas.ActionsActionV1

	// Process each action block
	for _, actionBlock := range actionContent.Blocks {
		action, err := parseActionBlock(actionBlock)
		if err != nil {
			return nil, err
		}
		actions = append(actions, action)
	}

	return actions, nil
}

// parseActionBlock parses a single action block into an ActionsActionV1 struct
func parseActionBlock(actionBlock *hcl.Block) (*schemas.ActionsActionV1, error) {
	actionName := actionBlock.Labels[0]
	logger := logging.GetGlobalLogger()
	logger.Debug("processing action", slog.String("action_name", actionName))

	// Create a new action struct with defaults
	action := createDefaultAction()

	// Extract action attributes
	actionAttrContent, err := extractActionAttributes(actionBlock)
	if err != nil {
		return nil, fmt.Errorf("failed to decode action %s attributes: %v", actionName, err)
	}

	// Parse all attributes
	if err := parseActionAttributes(action, actionAttrContent, actionName); err != nil {
		return nil, err
	}

	// Store the action name in the description for now (since the struct doesn't have a name field)
	// In a full implementation, we'd add a Name field to the struct
	action.Description = fmt.Sprintf("%s: %s", actionName, action.Description)

	logger.Debug("successfully parsed action",
		slog.String("action_name", actionName),
		slog.String("description", action.Description))

	return action, nil
}

// createDefaultAction creates an action with default values
func createDefaultAction() *schemas.ActionsActionV1 {
	return &schemas.ActionsActionV1{
		Description: "",
		Type:        "",
		Tags:        []string{},
		Targets:     []string{},
		Timeout:     300, // Default timeout
		Retries:     0,   // Default retries
		RetryDelay:  5,   // Default retry delay
	}
}

// extractActionAttributes extracts attributes from an action block
func extractActionAttributes(actionBlock *hcl.Block) (*hcl.BodyContent, error) {
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
			// Template deploy specific fields
			{Name: "source", Required: false},
			{Name: "destination", Required: false},
			{Name: "validate", Required: false},
			{Name: "backup", Required: false},
			{Name: "permissions", Required: false},
			{Name: "owner", Required: false},
			{Name: "group", Required: false},
			// File sync specific fields
			{Name: "sync_source", Required: false},
			{Name: "sync_destination", Required: false},
			{Name: "sync_delete", Required: false},
			{Name: "sync_preserve", Required: false},
		},
	}

	// Extract action attributes
	actionAttrContent, diags := actionBlock.Body.Content(actionAttrSchema)
	if diags.HasErrors() {
		return nil, fmt.Errorf("failed to decode action attributes: %v", diags)
	}

	return actionAttrContent, nil
}

// parseActionAttributes parses all attributes for an action
func parseActionAttributes(action *schemas.ActionsActionV1, attrContent *hcl.BodyContent, actionName string) error {
	// Parse basic attributes
	if err := parseBasicAttributes(action, attrContent, actionName); err != nil {
		return err
	}

	// Parse template deploy specific attributes
	if err := parseTemplateDeployAttributes(action, attrContent, actionName); err != nil {
		return err
	}

	// Parse file sync specific attributes
	if err := parseFileSyncAttributes(action, attrContent, actionName); err != nil {
		return err
	}

	return nil
}

// parseBasicAttributes parses basic action attributes
func parseBasicAttributes(action *schemas.ActionsActionV1, attrContent *hcl.BodyContent, actionName string) error {
	if err := parseStringAttribute(attrContent, "description", &action.Description, actionName); err != nil {
		return err
	}
	if err := parseStringAttribute(attrContent, "type", &action.Type, actionName); err != nil {
		return err
	}
	if err := parseStringAttribute(attrContent, "command", &action.Command, actionName); err != nil {
		return err
	}
	if err := parseStringSliceAttribute(attrContent, "targets", &action.Targets, actionName); err != nil {
		return err
	}
	if err := parseIntAttribute(attrContent, "timeout", &action.Timeout, actionName); err != nil {
		return err
	}
	if err := parseIntAttribute(attrContent, "retries", &action.Retries, actionName); err != nil {
		return err
	}
	if err := parseIntAttribute(attrContent, "retry_delay", &action.RetryDelay, actionName); err != nil {
		return err
	}
	if err := parseBoolAttribute(attrContent, "sudo", &action.Sudo, actionName); err != nil {
		return err
	}
	return nil
}

// parseTemplateDeployAttributes parses template deploy specific attributes
func parseTemplateDeployAttributes(action *schemas.ActionsActionV1, attrContent *hcl.BodyContent, actionName string) error {
	if err := parseStringAttribute(attrContent, "source", &action.Source, actionName); err != nil {
		return err
	}
	if err := parseStringAttribute(attrContent, "destination", &action.Destination, actionName); err != nil {
		return err
	}
	if err := parseBoolAttribute(attrContent, "validate", &action.Validate, actionName); err != nil {
		return err
	}
	if err := parseBoolAttribute(attrContent, "backup", &action.Backup, actionName); err != nil {
		return err
	}
	if err := parseStringAttribute(attrContent, "permissions", &action.Permissions, actionName); err != nil {
		return err
	}
	if err := parseStringAttribute(attrContent, "owner", &action.Owner, actionName); err != nil {
		return err
	}
	if err := parseStringAttribute(attrContent, "group", &action.Group, actionName); err != nil {
		return err
	}
	return nil
}

// parseFileSyncAttributes parses file sync specific attributes
func parseFileSyncAttributes(action *schemas.ActionsActionV1, attrContent *hcl.BodyContent, actionName string) error {
	if err := parseStringAttribute(attrContent, "sync_source", &action.SyncSource, actionName); err != nil {
		return err
	}
	if err := parseStringAttribute(attrContent, "sync_destination", &action.SyncDestination, actionName); err != nil {
		return err
	}
	if err := parseBoolAttribute(attrContent, "sync_delete", &action.SyncDelete, actionName); err != nil {
		return err
	}
	if err := parseBoolAttribute(attrContent, "sync_preserve", &action.SyncPreserve, actionName); err != nil {
		return err
	}
	return nil
}

// parseStringAttribute parses a string attribute
func parseStringAttribute(attrContent *hcl.BodyContent, attrName string, target *string, actionName string) error {
	if attr, exists := attrContent.Attributes[attrName]; exists {
		var value string
		if diags := gohcl.DecodeExpression(attr.Expr, nil, &value); diags.HasErrors() {
			return fmt.Errorf("failed to decode %s for action %s: %v", attrName, actionName, diags)
		}
		*target = value
	}
	return nil
}

// parseStringSliceAttribute parses a string slice attribute
func parseStringSliceAttribute(attrContent *hcl.BodyContent, attrName string, target *[]string, actionName string) error {
	if attr, exists := attrContent.Attributes[attrName]; exists {
		var value []string
		if diags := gohcl.DecodeExpression(attr.Expr, nil, &value); diags.HasErrors() {
			return fmt.Errorf("failed to decode %s for action %s: %v", attrName, actionName, diags)
		}
		*target = value
	}
	return nil
}

// parseIntAttribute parses an integer attribute
func parseIntAttribute(attrContent *hcl.BodyContent, attrName string, target *int, actionName string) error {
	if attr, exists := attrContent.Attributes[attrName]; exists {
		var value int
		if diags := gohcl.DecodeExpression(attr.Expr, nil, &value); diags.HasErrors() {
			return fmt.Errorf("failed to decode %s for action %s: %v", attrName, actionName, diags)
		}
		*target = value
	}
	return nil
}

// parseBoolAttribute parses a boolean attribute
func parseBoolAttribute(attrContent *hcl.BodyContent, attrName string, target *bool, actionName string) error {
	if attr, exists := attrContent.Attributes[attrName]; exists {
		var value bool
		if diags := gohcl.DecodeExpression(attr.Expr, nil, &value); diags.HasErrors() {
			return fmt.Errorf("failed to decode %s for action %s: %v", attrName, actionName, diags)
		}
		*target = value
	}
	return nil
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
	targetNames = action.Targets
	if len(overrideTargets) > 0 {
		targetNames = overrideTargets
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
	logger := logging.GetGlobalLogger()

	logger.Info("📊 Action Results")
	logger.Info("================")

	successCount := 0
	failureCount := 0

	for _, result := range results {
		if result.Success {
			successCount++
			logger.Info("✅ Action succeeded",
				slog.String("machine", result.Machine.Hostname),
				slog.String("message", result.Message))
		} else {
			failureCount++
			logger.Error("❌ Action failed",
				slog.String("machine", result.Machine.Hostname),
				slog.String("error", result.Error))
		}
	}

	logger.Info("Action execution summary",
		slog.Int("successful", successCount),
		slog.Int("failed", failureCount))
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
	sshManager    *ssh.Manager
	projectConfig *schemas.ProjectV1
}

// NewActionExecutor creates a new action executor
func NewActionExecutor(sshManager *ssh.Manager, projectConfig *schemas.ProjectV1) *ActionExecutor {
	return &ActionExecutor{
		sshManager:    sshManager,
		projectConfig: projectConfig,
	}
}

// RunAction executes an action across target machines.
// It takes a context, action configuration, list of target machines, and a dry-run flag.
// Returns a slice of action results and any error that occurred during execution.
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
	logger := logging.GetGlobalLogger()
	logger.Info("📄 Executing template_deploy action",
		slog.String("source", action.Source),
		slog.String("destination", action.Destination),
		slog.String("machine", machine.Hostname))

	// Validate required fields
	if action.Source == "" {
		return fmt.Errorf("template_deploy action requires a source field")
	}
	if action.Destination == "" {
		return fmt.Errorf("template_deploy action requires a destination field")
	}

	// Load variables configuration for template rendering
	variables, err := loadVariablesConfig()
	if err != nil {
		return fmt.Errorf("failed to load variables configuration: %w", err)
	}

	// Create template context
	templateContext, err := ae.buildTemplateContext(variables, machine)
	if err != nil {
		return fmt.Errorf("failed to build template context: %w", err)
	}

	// Create transparent decryptor for template rendering
	ageEncryption, err := encryption.NewAgeEncryption("", "")
	if err != nil {
		logger.Warn("failed to initialize age encryption, continuing without decryption support",
			slog.String("error", err.Error()))
		// Continue with nil encryption - template renderer will handle this gracefully
	}

	transparentDecryptor := encryption.NewTransparentDecryptor(ageEncryption)
	templateRenderer := templates.NewTemplateRenderer(transparentDecryptor)

	// Validate template if requested
	if action.Validate {
		logger.Debug("Validating template syntax", slog.String("template_path", action.Source))
		if err := templateRenderer.ValidateTemplate(action.Source); err != nil {
			return fmt.Errorf("template validation failed: %w", err)
		}
		logger.Debug("Template validation successful")
	}

	// Render template
	logger.Debug("Rendering template", slog.String("template_path", action.Source))
	renderedContent, err := templateRenderer.RenderTemplate(action.Source, templateContext)
	if err != nil {
		return fmt.Errorf("failed to render template: %s: %w", action.Source, err)
	}

	// Create backup if requested
	if err := ae.createBackupIfNeeded(ctx, action, machine); err != nil {
		return fmt.Errorf("failed to create backup: %w", err)
	}

	// Upload rendered content to destination
	logger.Debug("Uploading rendered template to destination",
		slog.String("destination", action.Destination),
		slog.String("machine", machine.Hostname))

	// For now, we'll use a simplified approach with echo and redirection
	// In a full implementation, you'd want to use proper SCP or SFTP
	escapedContent := strings.ReplaceAll(renderedContent, "'", "'\"'\"'")
	uploadCmd := fmt.Sprintf("echo '%s' > %s", escapedContent, action.Destination)

	result, err := ae.sshManager.RunCommandOnMachine(ctx, machine, uploadCmd)
	if err != nil {
		return fmt.Errorf("failed to upload template content: %w", err)
	}

	if result.ExitCode != 0 {
		return fmt.Errorf("upload command returned non-zero exit code: %d", result.ExitCode)
	}

	// Set file attributes (permissions, owner, group)
	if err := ae.setFileAttributes(ctx, action, machine); err != nil {
		return fmt.Errorf("failed to set file attributes: %w", err)
	}

	logger.Info("✅ Template deployed successfully",
		slog.String("source", action.Source),
		slog.String("destination", action.Destination),
		slog.String("machine", machine.Hostname))

	return nil
}

// executeFileSyncAction executes a file sync action
func (ae *ActionExecutor) runFileSyncAction(ctx context.Context, action *schemas.ActionsActionV1, machine *schemas.MachinesMachineV1) error {
	if action.SyncSource == "" || action.SyncDestination == "" {
		return fmt.Errorf("file_sync action requires sync_source and sync_destination fields")
	}

	logger := logging.GetGlobalLogger()
	logger.Info("🔄 Starting file synchronization",
		slog.String("source", action.SyncSource),
		slog.String("destination", action.SyncDestination),
		slog.String("machine", machine.Hostname))

	// Create remote sync engine
	remoteSyncEngine := sync.NewRemoteSyncEngine(ae.sshManager)

	// Determine sync mode based on action configuration
	syncMode := sync.ModeOneWayReplica // Default mode
	if action.SyncPreserve {
		syncMode = sync.ModeOneWaySafe
	}

	// Create sync options
	syncOptions := &sync.RemoteOptions{
		Options: &sync.Options{
			BlockLength:   sync.DefaultBlockLength,
			CreateBackup:  true,
			PreservePerms: action.SyncPreserve,
			PreserveOwner: false, // Could be configurable
			PreserveGroup: false, // Could be configurable
			DryRun:        false,
			Verbose:       true,
			Mode:          syncMode,
		},
		Machine:         machine,
		ProgressReport:  ae.createProgressReporter(action),
		ConflictResolve: sync.ConflictResolutionBackup,
		SyncDelete:      action.SyncDelete,
	}

	// Determine if this is a local-to-remote or remote-to-local sync
	// For now, we'll assume local-to-remote (source is local, destination is remote)
	localPath := action.SyncSource
	remotePath := action.SyncDestination

	// Validate local path exists
	if _, err := os.Stat(localPath); os.IsNotExist(err) {
		return fmt.Errorf("local source path does not exist: %s", localPath)
	}

	// Perform the synchronization
	result, err := remoteSyncEngine.SyncRemoteDirectory(ctx, localPath, remotePath, syncOptions)
	if err != nil {
		return fmt.Errorf("file synchronization failed: %v", err)
	}

	// Log results
	if result.Success {
		logger.Info("✅ File synchronization completed successfully",
			slog.String("machine", machine.Hostname),
			slog.Int64("bytes_transferred", result.BytesTransferred),
			slog.Int64("bytes_saved", result.BytesSaved),
			slog.Int("operations", result.Operations))

		if len(result.Conflicts) > 0 {
			logger.Info("⚠️ Conflicts detected during sync",
				slog.String("machine", machine.Hostname),
				slog.Int("conflict_count", len(result.Conflicts)))
		}
	} else {
		return fmt.Errorf("file synchronization failed: %v", result.Error)
	}

	return nil
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

// loadVariablesConfig loads the variables configuration from variables.hcl
func loadVariablesConfig() (map[string]interface{}, error) {
	logger := logging.GetGlobalLogger()
	logger.Debug("loading variables configuration")

	// Look for variables.hcl in current directory
	variablesHCLPath := "variables.hcl"
	if _, err := os.Stat(variablesHCLPath); os.IsNotExist(err) {
		logger.Debug("variables.hcl not found, using empty variables")
		return make(map[string]interface{}), nil
	}

	// Read and parse variables.hcl using the HCL library
	content, err := os.ReadFile(variablesHCLPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read variables.hcl: %w", err)
	}

	// Use the simplified validator to parse and validate the variables configuration
	validator := schemas.NewSimpleValidator()
	result, err := validator.ValidateHCLContent("variables", string(content))
	if err != nil {
		return nil, fmt.Errorf("failed to validate variables.hcl: %w", err)
	}

	if !result.IsValid {
		return nil, fmt.Errorf("variables.hcl validation failed: %s", result.Errors[0].Message)
	}

	// Parse HCL using the library
	file, diags := hclsyntax.ParseConfig(content, variablesHCLPath, hcl.Pos{Line: 1, Column: 1})
	if diags.HasErrors() {
		return nil, fmt.Errorf("failed to parse variables.hcl: %v", diags)
	}

	// Extract variables block
	variablesBlock, err := extractVariablesBlock(file)
	if err != nil {
		return nil, err
	}

	// Extract variable blocks
	variables, err := extractVariableBlocks(variablesBlock)
	if err != nil {
		return nil, err
	}

	logger.Debug("found variable blocks", slog.Int("count", len(variables.Blocks)))

	// Process variable blocks
	variablesMap, err := processVariableBlocks(variables)
	if err != nil {
		return nil, err
	}

	logger.Debug("loadVariablesConfig completed", slog.Int("variable_count", len(variablesMap)))
	return variablesMap, nil
}

// extractVariablesBlock extracts the variables block from the parsed HCL file
func extractVariablesBlock(file *hcl.File) (*hcl.Block, error) {
	// Parse the file body to find variables block
	body := file.Body
	content, diags := body.Content(&hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{
			{Type: "variables"},
		},
	})

	if diags.HasErrors() {
		return nil, fmt.Errorf("failed to decode variables block: %v", diags)
	}

	if len(content.Blocks) == 0 {
		return nil, fmt.Errorf("no variables block found in variables.hcl")
	}

	return content.Blocks[0], nil
}

// extractVariableBlocks extracts variable blocks from the variables block
func extractVariableBlocks(variablesBlock *hcl.Block) (*hcl.BodyContent, error) {
	// Define schema for variable blocks
	variableBlockSchema := &hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{
			{Type: "variable", LabelNames: []string{"name"}},
		},
	}

	// Extract variable blocks
	content, diags := variablesBlock.Body.Content(variableBlockSchema)
	if diags.HasErrors() {
		return nil, fmt.Errorf("failed to decode variable blocks: %v", diags)
	}

	return content, nil
}

// processVariableBlocks processes variable blocks into a map
func processVariableBlocks(content *hcl.BodyContent) (map[string]interface{}, error) {
	variables := make(map[string]interface{})

	for _, block := range content.Blocks {
		variableName := block.Labels[0]
		logger := logging.GetGlobalLogger()
		logger.Debug("processing variable", slog.String("variable_name", variableName))

		// Extract variable attributes
		variableAttrContent, err := extractVariableAttributes(block)
		if err != nil {
			return nil, fmt.Errorf("failed to decode variable %s attributes: %v", variableName, err)
		}

		// Parse variable attributes
		variableValue, err := parseVariableAttributes(variableAttrContent, variableName)
		if err != nil {
			return nil, err
		}

		variables[variableName] = variableValue
		logger.Debug("successfully parsed variable",
			slog.String("variable_name", variableName),
			slog.String("value_type", fmt.Sprintf("%T", variableValue)))
	}

	return variables, nil
}

// extractVariableAttributes extracts attributes from a variable block
func extractVariableAttributes(variableBlock *hcl.Block) (*hcl.BodyContent, error) {
	// Define schema for variable attributes
	variableAttrSchema := &hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{
			{Name: "value", Required: false},
			{Name: "encrypted_value", Required: false},
			{Name: "description", Required: false},
			{Name: "encrypted", Required: false},
		},
	}

	// Extract variable attributes
	content, diags := variableBlock.Body.Content(variableAttrSchema)
	if diags.HasErrors() {
		return nil, fmt.Errorf("failed to decode variable attributes: %v", diags)
	}

	return content, nil
}

// parseVariableAttributes parses all attributes for a variable
func parseVariableAttributes(attrContent *hcl.BodyContent, variableName string) (interface{}, error) {
	// Check for encrypted_value first (highest priority)
	if attr, exists := attrContent.Attributes["encrypted_value"]; exists {
		return parseEncryptedValueAttribute(attr, variableName)
	}

	// Check for plain value
	if attr, exists := attrContent.Attributes["value"]; exists {
		return parseValueAttribute(attr, variableName)
	}

	// Check for description
	if attr, exists := attrContent.Attributes["description"]; exists {
		return parseDescriptionAttribute(attr, variableName)
	}

	// Return empty string if no value found
	return "", nil
}

// parseEncryptedValueAttribute parses encrypted_value attribute
func parseEncryptedValueAttribute(attr *hcl.Attribute, variableName string) (interface{}, error) {
	val, diags := attr.Expr.Value(nil)
	if diags.HasErrors() {
		return nil, fmt.Errorf("failed to evaluate encrypted_value for variable %s: %v", variableName, diags)
	}

	if val.IsNull() {
		return nil, nil
	}

	// For encrypted values, we expect an object/map
	if val.Type().IsObjectType() {
		return parseEncryptedObjectValue(val)
	}

	// Fallback to string representation
	return val.AsString(), nil
}

// parseEncryptedObjectValue parses encrypted object value into map
func parseEncryptedObjectValue(val cty.Value) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	for key, value := range val.AsValueMap() {
		if value.Type() == cty.String {
			result[key] = value.AsString()
		} else {
			result[key] = value.AsString() // Fallback to string representation
		}
	}
	return result, nil
}

// parseValueAttribute parses value attribute with type conversion
func parseValueAttribute(attr *hcl.Attribute, variableName string) (interface{}, error) {
	val, diags := attr.Expr.Value(nil)
	if diags.HasErrors() {
		return nil, fmt.Errorf("failed to evaluate value for variable %s: %v", variableName, diags)
	}

	if val.IsNull() {
		return nil, nil
	}

	return convertCtyValueToGoType(val)
}

// convertCtyValueToGoType converts cty.Value to appropriate Go type
func convertCtyValueToGoType(val cty.Value) (interface{}, error) {
	switch val.Type() {
	case cty.String:
		return val.AsString(), nil
	case cty.Number:
		// Convert number to float64 for simplicity
		floatVal, _ := val.AsBigFloat().Float64()
		return floatVal, nil
	case cty.Bool:
		return val.True(), nil
	case cty.List(cty.String):
		return convertStringList(val)
	default:
		return val.AsString(), nil // Fallback to string representation
	}
}

// convertStringList converts cty list of strings to Go string slice
func convertStringList(val cty.Value) ([]string, error) {
	var result []string
	for _, item := range val.AsValueSlice() {
		result = append(result, item.AsString())
	}
	return result, nil
}

// parseDescriptionAttribute parses description attribute
func parseDescriptionAttribute(attr *hcl.Attribute, variableName string) (string, error) {
	var description string
	if err := gohcl.DecodeExpression(attr.Expr, nil, &description); err != nil {
		return "", fmt.Errorf("failed to decode description for variable %s: %v", variableName, err)
	}
	return description, nil
}

// buildTemplateContext creates a template context with all necessary data
func (ae *ActionExecutor) buildTemplateContext(variables map[string]interface{}, machine *schemas.MachinesMachineV1) (*templates.TemplateContext, error) {
	templateContext := &templates.TemplateContext{
		Variables:   variables,
		Machines:    make(map[string]interface{}),
		Facts:       make(map[string]interface{}),
		Environment: ae.buildEnvironmentContext(machine),
		Project:     ae.buildProjectContext(),
	}

	// Add machine information to context
	templateContext.Machines[machine.Hostname] = ae.buildMachineContext(machine)

	return templateContext, nil
}

// buildEnvironmentContext creates the environment context for template rendering
func (ae *ActionExecutor) buildEnvironmentContext(machine *schemas.MachinesMachineV1) map[string]string {
	return map[string]string{
		"HOSTNAME": machine.Hostname,
		"USER":     machine.User,
	}
}

// buildProjectContext creates the project context for template rendering
func (ae *ActionExecutor) buildProjectContext() map[string]interface{} {
	return map[string]interface{}{
		"name":        ae.projectConfig.Name,
		"description": ae.projectConfig.Description,
	}
}

// buildMachineContext creates the machine context for template rendering
func (ae *ActionExecutor) buildMachineContext(machine *schemas.MachinesMachineV1) map[string]interface{} {
	// Note: machine.Variables is not implemented in the current schema
	// This is a placeholder for future enhancement
	return map[string]interface{}{
		"hostname": machine.Hostname,
		"user":     machine.User,
		"port":     machine.Port,
	}
}

// createBackupIfNeeded creates a backup of the destination file if backup is requested
func (ae *ActionExecutor) createBackupIfNeeded(ctx context.Context, action *schemas.ActionsActionV1, machine *schemas.MachinesMachineV1) error {
	if !action.Backup {
		return nil
	}

	logger := logging.GetGlobalLogger()
	logger.Debug("Creating backup of existing file", slog.String("destination", action.Destination))

	backupCmd := fmt.Sprintf("cp %s %s.backup.$(date +%%Y%%m%%d_%%H%%M%%S)", action.Destination, action.Destination)
	result, err := ae.sshManager.RunCommandOnMachine(ctx, machine, backupCmd)
	if err != nil {
		logger.Warn("failed to create backup, continuing with deployment",
			slog.String("error", err.Error()))
	} else if result.ExitCode != 0 {
		logger.Warn("backup command returned non-zero exit code, continuing with deployment",
			slog.Int("exit_code", result.ExitCode))
	}

	return nil
}

// setFileAttributes sets file permissions, owner, and group if specified
func (ae *ActionExecutor) setFileAttributes(ctx context.Context, action *schemas.ActionsActionV1, machine *schemas.MachinesMachineV1) error {
	if action.Permissions == "" && action.Owner == "" && action.Group == "" {
		return nil
	}

	logger := logging.GetGlobalLogger()
	logger.Debug("Setting file attributes",
		slog.String("permissions", action.Permissions),
		slog.String("owner", action.Owner),
		slog.String("group", action.Group))

	// Set permissions
	if err := ae.setFilePermissions(ctx, action, machine); err != nil {
		return err
	}

	// Set owner
	if err := ae.setFileOwner(ctx, action, machine); err != nil {
		return err
	}

	// Set group
	if err := ae.setFileGroup(ctx, action, machine); err != nil {
		return err
	}

	return nil
}

// setFilePermissions sets file permissions if specified
func (ae *ActionExecutor) setFilePermissions(ctx context.Context, action *schemas.ActionsActionV1, machine *schemas.MachinesMachineV1) error {
	if action.Permissions == "" {
		return nil
	}

	return ae.executeFileAttributeCommand(ctx, machine, "chmod", action.Permissions, action.Destination, "permissions")
}

// setFileOwner sets file owner if specified
func (ae *ActionExecutor) setFileOwner(ctx context.Context, action *schemas.ActionsActionV1, machine *schemas.MachinesMachineV1) error {
	if action.Owner == "" {
		return nil
	}

	return ae.executeFileAttributeCommand(ctx, machine, "chown", action.Owner, action.Destination, "owner")
}

// setFileGroup sets file group if specified
func (ae *ActionExecutor) setFileGroup(ctx context.Context, action *schemas.ActionsActionV1, machine *schemas.MachinesMachineV1) error {
	if action.Group == "" {
		return nil
	}

	return ae.executeFileAttributeCommand(ctx, machine, "chgrp", action.Group, action.Destination, "group")
}

// executeFileAttributeCommand is a shared utility for executing file attribute commands
func (ae *ActionExecutor) executeFileAttributeCommand(ctx context.Context, machine *schemas.MachinesMachineV1, command, value, destination, attributeType string) error {
	logger := logging.GetGlobalLogger()
	cmd := fmt.Sprintf("%s %s %s", command, value, destination)
	result, err := ae.sshManager.RunCommandOnMachine(ctx, machine, cmd)
	if err != nil {
		logger.Warn(fmt.Sprintf("failed to set file %s", attributeType), slog.String("error", err.Error()))
	} else if result.ExitCode != 0 {
		logger.Warn(fmt.Sprintf("%s command returned non-zero exit code", command), slog.Int("exit_code", result.ExitCode))
	}

	return nil
}

// createProgressReporter creates a progress reporting function for file synchronization
func (ae *ActionExecutor) createProgressReporter(action *schemas.ActionsActionV1) func(progress *sync.Progress) {
	return func(progress *sync.Progress) {
		logger := logging.GetGlobalLogger()

		// Calculate percentage
		if progress.TotalFiles > 0 {
			progress.Percentage = float64(progress.FilesProcessed) / float64(progress.TotalFiles) * 100
		}

		// Log progress at appropriate intervals
		if progress.FilesProcessed%10 == 0 || progress.FilesProcessed == progress.TotalFiles {
			logger.Info("🔄 File sync progress",
				slog.String("action", action.Description),
				slog.String("current_file", progress.CurrentFile),
				slog.String("operation", progress.CurrentOperation),
				slog.Int64("files_processed", progress.FilesProcessed),
				slog.Int64("total_files", progress.TotalFiles),
				slog.Float64("percentage", progress.Percentage),
				slog.Int64("bytes_transferred", progress.BytesTransferred),
				slog.Int64("bytes_saved", progress.BytesSaved))
		}
	}
}
