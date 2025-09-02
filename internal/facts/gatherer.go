package facts

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"spooky/internal/logging"
	"spooky/internal/schemas"
	"spooky/internal/ssh"
	"spooky/internal/utilities"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/pkg/errors"
)

// Gatherer collects facts from remote machines
type Gatherer struct {
	sshManager *ssh.SimpleSSHManager
	config     *schemas.ProjectV1
}

// NewGatherer creates a new facts gatherer
func NewGatherer(sshManager *ssh.SimpleSSHManager, config *schemas.ProjectV1) *Gatherer {
	return &Gatherer{
		sshManager: sshManager,
		config:     config,
	}
}

// MachineFacts represents all facts collected from a single machine
type MachineFacts struct {
	Machine       *schemas.MachinesMachineV1
	BasicFacts    *schemas.BasicFactsV1
	EnhancedFacts *schemas.EnhancedFactsV1
	CustomFacts   *schemas.CustomFactsV1
	CollectedAt   time.Time
	Error         error
}

// GatherFactsFromMachine collects all facts from a single machine
func (g *Gatherer) GatherFactsFromMachine(ctx context.Context, machine *schemas.MachinesMachineV1) (*MachineFacts, error) {
	result := &MachineFacts{
		Machine:     machine,
		CollectedAt: time.Now(),
	}

	// Collect basic facts via SSH commands
	basicFacts, err := g.gatherBasicFacts(ctx, machine)
	if err != nil {
		result.Error = errors.Wrap(err, "failed to gather basic facts")
		return result, result.Error
	}
	result.BasicFacts = basicFacts

	// Collect enhanced facts if spooky-facts is available
	enhancedFacts, err := g.gatherEnhancedFacts(ctx, machine)
	if err != nil {
		// Enhanced facts are optional, so we don't fail the entire operation
		logger := logging.GetGlobalLogger()
		logger.Warn("failed to gather enhanced facts, continuing with basic facts only",
			slog.String("machine", machine.Hostname),
			slog.String("error", err.Error()))
	}
	result.EnhancedFacts = enhancedFacts

	// Collect custom facts (age-encrypted)
	customFacts, err := g.gatherCustomFacts(ctx, machine)
	if err != nil {
		// Custom facts are optional, so we don't fail the entire operation
		logger := logging.GetGlobalLogger()
		logger.Warn("failed to gather custom facts, continuing with basic facts only",
			slog.String("machine", machine.Hostname),
			slog.String("error", err.Error()))
	}
	result.CustomFacts = customFacts

	return result, nil
}

// GatherFactsFromMachines collects facts from multiple machines in parallel
func (g *Gatherer) GatherFactsFromMachines(ctx context.Context, machines []*schemas.MachinesMachineV1) ([]*MachineFacts, error) {
	maxParallel := g.config.FactsParallelCollection
	if maxParallel <= 0 {
		maxParallel = 10 // Default parallel limit
	}

	// Create a semaphore to limit concurrent connections
	semaphore := make(chan struct{}, maxParallel)
	results := make([]*MachineFacts, len(machines))
	errors := make(chan error, len(machines))

	// Start gathering facts for each machine
	for i, machine := range machines {
		go func(index int, m *schemas.MachinesMachineV1) {
			semaphore <- struct{}{}        // Acquire semaphore
			defer func() { <-semaphore }() // Release semaphore

			facts, err := g.GatherFactsFromMachine(ctx, m)
			if err != nil {
				errors <- fmt.Errorf("failed to gather facts from machine %s: %w", m.Hostname, err)
				return
			}
			results[index] = facts
		}(i, machine)
	}

	// Wait for all goroutines to complete
	timeout := time.Duration(g.config.FactsTimeout) * time.Second
	if timeout <= 0 {
		timeout = 30 * time.Second // Default timeout
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Wait for completion or timeout
	completed := 0
	for completed < len(machines) {
		select {
		case <-ctx.Done():
			return results, fmt.Errorf("facts gathering timed out: %w", ctx.Err())
		case err := <-errors:
			return results, err
		default:
			completed = 0
			for _, result := range results {
				if result != nil {
					completed++
				}
			}
			time.Sleep(100 * time.Millisecond)
		}
	}

	return results, nil
}

// createEmptyBasicFacts creates an empty BasicFactsV1 structure with initialized fact maps
func createEmptyBasicFacts() *schemas.BasicFactsV1 {
	return &schemas.BasicFactsV1{
		SystemFacts:   &schemas.SystemFactsV1{Facts: make(map[string]*schemas.FactV1)},
		HardwareFacts: &schemas.HardwareFactsV1{Facts: make(map[string]*schemas.FactV1)},
		NetworkFacts:  &schemas.NetworkFactsV1{Facts: make(map[string]*schemas.FactV1)},
		OSFacts:       &schemas.OSFactsV1{Facts: make(map[string]*schemas.FactV1)},
		UserFacts:     &schemas.UserFactsV1{Facts: make(map[string]*schemas.FactV1)},
		RuntimeFacts:  &schemas.RuntimeFactsV1{Facts: make(map[string]*schemas.FactV1)},
	}
}

// gatherBasicFacts collects basic system facts via SSH commands
func (g *Gatherer) gatherBasicFacts(ctx context.Context, machine *schemas.MachinesMachineV1) (*schemas.BasicFactsV1, error) {
	basicFacts := createEmptyBasicFacts()

	// Execute all fact categories
	if err := g.executeSystemFacts(ctx, machine, basicFacts); err != nil {
		return nil, err
	}
	if err := g.executeHardwareFacts(ctx, machine, basicFacts); err != nil {
		return nil, err
	}
	if err := g.executeNetworkFacts(ctx, machine, basicFacts); err != nil {
		return nil, err
	}
	if err := g.executeOSFacts(ctx, machine, basicFacts); err != nil {
		return nil, err
	}
	if err := g.executeUserFacts(ctx, machine, basicFacts); err != nil {
		return nil, err
	}
	if err := g.executeRuntimeFacts(ctx, machine, basicFacts); err != nil {
		return nil, err
	}

	return basicFacts, nil
}

// executeSystemFacts executes system-related fact commands
func (g *Gatherer) executeSystemFacts(ctx context.Context, machine *schemas.MachinesMachineV1, basicFacts *schemas.BasicFactsV1) error {
	systemCommands := map[string]string{
		"hostname":               "hostname",
		"os":                     "uname -s",
		"architecture":           "uname -m",
		"kernel_version":         "uname -r",
		"nodename":               "uname -n",
		"domain":                 "hostname -d",
		"fqdn":                   "hostname -f",
		"machine":                "uname -m",
		"userspace_architecture": "uname -m",
		"userspace_bits":         "getconf LONG_BIT",
	}

	return g.executeFactCommands(ctx, machine, systemCommands, basicFacts.SystemFacts.Facts)
}

// executeHardwareFacts executes hardware-related fact commands
func (g *Gatherer) executeHardwareFacts(ctx context.Context, machine *schemas.MachinesMachineV1, basicFacts *schemas.BasicFactsV1) error {
	hardwareCommands := map[string]string{
		"cpu_count":                  "nproc",
		"memory_total":               "free -b | grep '^Mem:' | awk '{print $2}'",
		"memory_used":                "free -b | grep '^Mem:' | awk '{print $3}'",
		"memory_free":                "free -b | grep '^Mem:' | awk '{print $4}'",
		"swap_total":                 "free -b | grep '^Swap:' | awk '{print $2}'",
		"swap_free":                  "free -b | grep '^Swap:' | awk '{print $4}'",
		"load_1":                     "cat /proc/loadavg | awk '{print $1}'",
		"load_5":                     "cat /proc/loadavg | awk '{print $2}'",
		"load_15":                    "cat /proc/loadavg | awk '{print $3}'",
		"processor_cores":            "grep -c '^processor' /proc/cpuinfo",
		"processor_vcpus":            "grep -c '^processor' /proc/cpuinfo",
		"processor_count":            "grep -c '^processor' /proc/cpuinfo",
		"processor_nproc":            "nproc",
		"processor_threads_per_core": "grep '^cpu cores' /proc/cpuinfo | head -1 | cut -d: -f2 | tr -d ' '",
	}

	return g.executeFactCommands(ctx, machine, hardwareCommands, basicFacts.HardwareFacts.Facts)
}

// executeNetworkFacts executes network-related fact commands
func (g *Gatherer) executeNetworkFacts(ctx context.Context, machine *schemas.MachinesMachineV1, basicFacts *schemas.BasicFactsV1) error {
	networkCommands := map[string]string{
		"interfaces":         "ip -o link show | awk '{print $2}' | sed 's/://' | tr '\n' ','",
		"default_ip":         "ip route get 8.8.8.8 | awk '{print $7}'",
		"all_ipv4":           "ip -4 addr show | grep 'inet ' | awk '{print $2}' | cut -d'/' -f1 | tr '\n' ','",
		"dns":                "cat /etc/resolv.conf | grep '^nameserver' | awk '{print $2}' | tr '\n' ','",
		"all_ipv4_addresses": "ip -4 addr show | grep 'inet ' | awk '{print $2}' | cut -d'/' -f1 | tr '\n' ','",
		"all_ipv6_addresses": "ip -6 addr show | grep 'inet6 ' | awk '{print $2}' | cut -d'/' -f1 | tr '\n' ','",
	}

	return g.executeFactCommands(ctx, machine, networkCommands, basicFacts.NetworkFacts.Facts)
}

// executeOSFacts executes operating system-related fact commands
func (g *Gatherer) executeOSFacts(ctx context.Context, machine *schemas.MachinesMachineV1, basicFacts *schemas.BasicFactsV1) error {
	osCommands := map[string]string{
		"os_version":                 "uname -r",
		"os_family":                  "cat /etc/os-release | grep '^ID=' | cut -d'=' -f2 | tr -d '\"'",
		"pkg_mgr":                    "which apt && echo 'apt' || which yum && echo 'yum' || which dnf && echo 'dnf' || which pacman && echo 'pacman' || echo 'unknown'",
		"distribution":               "cat /etc/os-release | grep '^ID=' | cut -d'=' -f2 | tr -d '\"'",
		"distribution_version":       "cat /etc/os-release | grep '^VERSION_ID=' | cut -d'=' -f2 | tr -d '\"'",
		"distribution_major_version": "cat /etc/os-release | grep '^VERSION_ID=' | cut -d'=' -f2 | tr -d '\"' | cut -d'.' -f1",
		"distribution_release":       "cat /etc/os-release | grep '^VERSION_CODENAME=' | cut -d'=' -f2 | tr -d '\"'",
		"service_mgr":                "systemctl --version >/dev/null 2>&1 && echo 'systemd' || echo 'unknown'",
		"lsb":                        "lsb_release -a 2>/dev/null || echo 'lsb_release not available'",
	}

	return g.executeFactCommands(ctx, machine, osCommands, basicFacts.OSFacts.Facts)
}

// executeUserFacts executes user-related fact commands
func (g *Gatherer) executeUserFacts(ctx context.Context, machine *schemas.MachinesMachineV1, basicFacts *schemas.BasicFactsV1) error {
	userCommands := map[string]string{
		"user":          "whoami",
		"home":          "echo $HOME",
		"shell":         "echo $SHELL",
		"user_id":       "id -u",
		"group_id":      "id -g",
		"real_user_id":  "id -u",
		"real_group_id": "id -g",
		"user_uid":      "id -u",
		"user_gid":      "id -g",
		"user_dir":      "echo $HOME",
		"user_shell":    "echo $SHELL",
		"user_gecos":    "id -P | cut -d: -f5",
	}

	return g.executeFactCommands(ctx, machine, userCommands, basicFacts.UserFacts.Facts)
}

// executeRuntimeFacts executes runtime-related fact commands
func (g *Gatherer) executeRuntimeFacts(ctx context.Context, machine *schemas.MachinesMachineV1, basicFacts *schemas.BasicFactsV1) error {
	runtimeCommands := map[string]string{
		"uptime":         "cat /proc/uptime | awk '{print $1}'",
		"date":           "date -Iseconds",
		"cmdline":        "cat /proc/cmdline",
		"uptime_seconds": "cat /proc/uptime | awk '{print $1}'",
		"date_time":      "date -Iseconds",
	}

	return g.executeFactCommands(ctx, machine, runtimeCommands, basicFacts.RuntimeFacts.Facts)
}

// executeFactCommands executes a set of fact commands and stores results
func (g *Gatherer) executeFactCommands(ctx context.Context, machine *schemas.MachinesMachineV1, commands map[string]string, facts map[string]*schemas.FactV1) error {
	for factName, command := range commands {
		factType := g.determineFactType(factName)
		if fact := g.runCommand(ctx, machine, command, factName, factType); fact != nil {
			facts[factName] = fact
		}
	}
	return nil
}

// determineFactType determines the appropriate fact type based on the fact name
func (g *Gatherer) determineFactType(factName string) string {
	// Define numeric fact types
	numericFacts := map[string]bool{
		"cpu_count": true, "memory_total": true, "memory_used": true, "memory_free": true,
		"swap_total": true, "swap_free": true, "load_1": true, "load_5": true, "load_15": true,
		"processor_cores": true, "processor_vcpus": true, "processor_count": true,
		"processor_nproc": true, "processor_threads_per_core": true,
	}

	if numericFacts[factName] {
		return "number"
	}
	return "string"
}

// runCommand executes a command on a machine and returns a FactV1 if successful
func (g *Gatherer) runCommand(ctx context.Context, machine *schemas.MachinesMachineV1, command, factName, factType string) *schemas.FactV1 {
	result, err := utilities.RunCommand(ctx, machine.Hostname, command, machine, g.sshManager, 0, 0, 0)
	if err != nil {
		// Skip failed commands, but log them
		logger := logging.GetGlobalLogger()
		logger.Warn("failed to execute command, skipping fact collection",
			slog.String("command", command),
			slog.String("machine", machine.Hostname),
			slog.String("error", err.Error()))
		return nil
	}

	if result.ExitCode == 0 {
		value := strings.TrimSpace(result.Stdout)
		if value != "" {
			return &schemas.FactV1{
				Value:       value,
				Type:        factType,
				Description: fmt.Sprintf("Basic system fact: %s", factName),
			}
		}
	}

	return nil
}

// gatherEnhancedFacts collects enhanced facts from spooky-facts tool
func (g *Gatherer) gatherEnhancedFacts(ctx context.Context, machine *schemas.MachinesMachineV1) (*schemas.EnhancedFactsV1, error) {
	// Check if spooky-facts is available
	result, err := g.sshManager.RunCommandOnMachine(ctx, machine, "which spooky-facts")
	if err != nil || result.ExitCode != 0 {
		return nil, errors.New("spooky-facts not available on remote machine")
	}

	// Run spooky-facts to generate facts
	result, err = g.sshManager.RunCommandOnMachine(ctx, machine, "spooky-facts gather")
	if err != nil {
		return nil, errors.Wrap(err, "failed to run spooky-facts")
	}

	if result.ExitCode != 0 {
		return nil, errors.Errorf("spooky-facts failed: %s", result.Stderr)
	}

	// Read the generated facts file
	result, err = g.sshManager.RunCommandOnMachine(ctx, machine, "cat /etc/spooky/facts.hcl")
	if err != nil {
		// Try alternative location
		result, err = g.sshManager.RunCommandOnMachine(ctx, machine, "cat ~/.config/spooky/facts.hcl")
		if err != nil {
			return nil, errors.Wrap(err, "failed to read facts file")
		}
	}

	if result.ExitCode != 0 {
		return nil, errors.Errorf("failed to read facts file: %s", result.Stderr)
	}

	// Parse the HCL content
	enhancedFacts, err := g.parseFactsHCL(result.Stdout)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse facts HCL")
	}

	return enhancedFacts, nil
}

// gatherCustomFacts collects custom facts (may contain age-encrypted facts)
func (g *Gatherer) gatherCustomFacts(ctx context.Context, machine *schemas.MachinesMachineV1) (*schemas.CustomFactsV1, error) {
	// Check if custom facts file exists
	result, err := g.sshManager.RunCommandOnMachine(ctx, machine, "test -f /etc/spooky/custom.hcl && echo 'exists'")
	if err != nil || result.ExitCode != 0 {
		// Try alternative location
		result, err = g.sshManager.RunCommandOnMachine(ctx, machine, "test -f ~/.config/spooky/custom.hcl && echo 'exists'")
		if err != nil || result.ExitCode != 0 {
			return nil, errors.New("custom facts file not found")
		}
	}

	// Read the custom facts file
	result, err = g.sshManager.RunCommandOnMachine(ctx, machine, "cat /etc/spooky/custom.hcl")
	if err != nil {
		// Try alternative location
		result, err = g.sshManager.RunCommandOnMachine(ctx, machine, "cat ~/.config/spooky/custom.hcl")
		if err != nil {
			return nil, errors.Wrap(err, "failed to read custom facts file")
		}
	}

	if result.ExitCode != 0 {
		return nil, errors.Errorf("failed to read custom facts file: %s", result.Stderr)
	}

	// Parse the HCL content
	enhancedFacts, err := g.parseFactsHCL(result.Stdout)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse custom facts HCL")
	}

	// Convert to CustomFactsV1
	customFacts := &schemas.CustomFactsV1{
		Facts: enhancedFacts.Facts,
	}

	return customFacts, nil
}

// parseFactsHCL parses HCL content and returns facts structure
func (g *Gatherer) parseFactsHCL(hclContent string) (*schemas.EnhancedFactsV1, error) {
	// Parse HCL content
	file, diags := hclsyntax.ParseConfig([]byte(hclContent), "facts.hcl", hcl.Pos{Line: 1, Column: 1})
	if diags.HasErrors() {
		return nil, errors.Wrapf(errors.New("HCL parsing failed"), "failed to parse facts HCL: %v", diags)
	}

	// Define schema for facts blocks
	schema := &hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{
			{
				Type:       "fact",
				LabelNames: []string{"name"},
			},
		},
	}

	// Extract fact blocks
	bodyContent, diags := file.Body.Content(schema)
	if diags.HasErrors() {
		return nil, errors.Wrapf(errors.New("HCL decoding failed"), "failed to decode facts HCL: %v", diags)
	}

	// Create enhanced facts structure
	enhancedFacts := &schemas.EnhancedFactsV1{
		Facts: make(map[string]*schemas.FactV1),
	}

	// Process each fact block
	for _, block := range bodyContent.Blocks {
		factName := block.Labels[0]

		// Define schema for fact attributes
		factSchema := &hcl.BodySchema{
			Attributes: []hcl.AttributeSchema{
				{Name: "value", Required: true},
				{Name: "type", Required: false},
				{Name: "description", Required: false},
			},
		}

		// Extract fact attributes
		factContent, diags := block.Body.Content(factSchema)
		if diags.HasErrors() {
			continue // Skip this fact if we can't parse it
		}

		// Get the fact value
		if valueAttr, exists := factContent.Attributes["value"]; exists {
			var factValue string
			if diags := gohcl.DecodeExpression(valueAttr.Expr, nil, &factValue); diags.HasErrors() {
				continue // Skip if we can't decode the value
			}

			// Create fact structure
			fact := &schemas.FactV1{
				Value: factValue,
			}

			// Get optional type
			if typeAttr, exists := factContent.Attributes["type"]; exists {
				var factType string
				if diags := gohcl.DecodeExpression(typeAttr.Expr, nil, &factType); !diags.HasErrors() {
					fact.Type = factType
				}
			}

			// Get optional description
			if descAttr, exists := factContent.Attributes["description"]; exists {
				var description string
				if diags := gohcl.DecodeExpression(descAttr.Expr, nil, &description); !diags.HasErrors() {
					fact.Description = description
				}
			}

			enhancedFacts.Facts[factName] = fact
		}
	}

	return enhancedFacts, nil
}

// ExportFacts exports collected facts to HCL format
func (g *Gatherer) ExportFacts(machineFacts []*MachineFacts) (*schemas.FactsV1, error) {
	// Combine all machine facts into a single facts structure
	// This is a simplified version - in practice, you might want to organize by machine
	combinedFacts := &schemas.FactsV1{
		BasicFacts:    createEmptyBasicFacts(),
		EnhancedFacts: &schemas.EnhancedFactsV1{Facts: make(map[string]*schemas.FactV1)},
		CustomFacts:   &schemas.CustomFactsV1{Facts: make(map[string]*schemas.FactV1)},
	}

	for _, machineFact := range machineFacts {
		if machineFact.Error != nil {
			continue // Skip machines with errors
		}

		// Add basic facts with machine prefix, organized by category
		if machineFact.BasicFacts != nil {
			// System facts
			if machineFact.BasicFacts.SystemFacts != nil {
				for factName, fact := range machineFact.BasicFacts.SystemFacts.Facts {
					prefixedName := fmt.Sprintf("%s_%s", machineFact.Machine.Hostname, factName)
					combinedFacts.BasicFacts.SystemFacts.Facts[prefixedName] = fact
				}
			}

			// Hardware facts
			if machineFact.BasicFacts.HardwareFacts != nil {
				for factName, fact := range machineFact.BasicFacts.HardwareFacts.Facts {
					prefixedName := fmt.Sprintf("%s_%s", machineFact.Machine.Hostname, factName)
					combinedFacts.BasicFacts.HardwareFacts.Facts[prefixedName] = fact
				}
			}

			// Network facts
			if machineFact.BasicFacts.NetworkFacts != nil {
				for factName, fact := range machineFact.BasicFacts.NetworkFacts.Facts {
					prefixedName := fmt.Sprintf("%s_%s", machineFact.Machine.Hostname, factName)
					combinedFacts.BasicFacts.NetworkFacts.Facts[prefixedName] = fact
				}
			}

			// OS facts
			if machineFact.BasicFacts.OSFacts != nil {
				for factName, fact := range machineFact.BasicFacts.OSFacts.Facts {
					prefixedName := fmt.Sprintf("%s_%s", machineFact.Machine.Hostname, factName)
					combinedFacts.BasicFacts.OSFacts.Facts[prefixedName] = fact
				}
			}

			// User facts
			if machineFact.BasicFacts.UserFacts != nil {
				for factName, fact := range machineFact.BasicFacts.UserFacts.Facts {
					prefixedName := fmt.Sprintf("%s_%s", machineFact.Machine.Hostname, factName)
					combinedFacts.BasicFacts.UserFacts.Facts[prefixedName] = fact
				}
			}

			// Runtime facts
			if machineFact.BasicFacts.RuntimeFacts != nil {
				for factName, fact := range machineFact.BasicFacts.RuntimeFacts.Facts {
					prefixedName := fmt.Sprintf("%s_%s", machineFact.Machine.Hostname, factName)
					combinedFacts.BasicFacts.RuntimeFacts.Facts[prefixedName] = fact
				}
			}
		}

		// Add enhanced facts with machine prefix
		if machineFact.EnhancedFacts != nil {
			for factName, fact := range machineFact.EnhancedFacts.Facts {
				prefixedName := fmt.Sprintf("%s_%s", machineFact.Machine.Hostname, factName)
				combinedFacts.EnhancedFacts.Facts[prefixedName] = fact
			}
		}

		// Add custom facts with machine prefix
		if machineFact.CustomFacts != nil {
			for factName, fact := range machineFact.CustomFacts.Facts {
				prefixedName := fmt.Sprintf("%s_%s", machineFact.Machine.Hostname, factName)
				combinedFacts.CustomFacts.Facts[prefixedName] = fact
			}
		}
	}

	return combinedFacts, nil
}
