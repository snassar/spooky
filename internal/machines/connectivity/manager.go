package connectivity

import (
	"context"
	"fmt"
	"net"
	"os"
	"time"

	configtypes "spooky/internal/types/config"
	"spooky/internal/machines/types"

	"golang.org/x/crypto/ssh"
)

// Manager provides comprehensive connectivity testing capabilities
type Manager struct {
	options *types.ConnectivityTestOptions
}

// NewManager creates a new connectivity manager with the given options
func NewManager(options *types.ConnectivityTestOptions) *Manager {
	if options == nil {
		options = types.DefaultConnectivityTestOptions()
	}

	// Validate and fix options
	if options.Parallel <= 0 {
		options.Parallel = 10
	}
	if options.Timeout <= 0 {
		options.Timeout = 30 * time.Second
	}
	if options.DNSTimeout <= 0 {
		options.DNSTimeout = 5 * time.Second
	}
	if options.SSHCommand == "" {
		options.SSHCommand = "echo 'SSH connectivity test'"
	}
	if len(options.Phases) == 0 {
		options.Phases = []types.ConnectivityTestPhase{types.PhaseDNS, types.PhaseSSH}
	}

	return &Manager{
		options: options,
	}
}

// TestMachineConnectivity performs comprehensive connectivity testing on a single machine
func (cm *Manager) TestMachineConnectivity(ctx context.Context, machine *configtypes.Machine) []types.ConnectivityTestResult {
	var results []types.ConnectivityTestResult

	// Get relevant phases for this machine
	phases := cm.getRelevantPhases(machine)

	// Test each phase
	for _, phase := range phases {
		result := cm.testPhase(ctx, machine, phase)
		results = append(results, result)
	}

	return results
}

// TestMachinesConnectivity performs connectivity testing on multiple machines
func (cm *Manager) TestMachinesConnectivity(ctx context.Context, machines []*configtypes.Machine) map[string][]types.ConnectivityTestResult {
	results := make(map[string][]types.ConnectivityTestResult)

	// Test each machine
	for _, machine := range machines {
		machineResults := cm.TestMachineConnectivity(ctx, machine)
		results[machine.Name] = machineResults
	}

	return results
}

// TestSpecificPhase tests a specific connectivity phase for a machine
func (cm *Manager) TestSpecificPhase(ctx context.Context, machine *configtypes.Machine, phase types.ConnectivityTestPhase) types.ConnectivityTestResult {
	return cm.testPhase(ctx, machine, phase)
}

// TestPhases tests specific phases for a machine
func (cm *Manager) TestPhases(ctx context.Context, machine *configtypes.Machine, phases []types.ConnectivityTestPhase) []types.ConnectivityTestResult {
	var results []types.ConnectivityTestResult

	for _, phase := range phases {
		result := cm.testPhase(ctx, machine, phase)
		results = append(results, result)
	}

	return results
}

// GetTestSummary provides a summary of connectivity test results
func (cm *Manager) GetTestSummary(results map[string][]types.ConnectivityTestResult) *types.ConnectivityTestSummary {
	summary := &types.ConnectivityTestSummary{
		PhaseResults: make(map[types.ConnectivityTestPhase]*types.PhaseSummary),
	}

	for _, machineResults := range results {
		for _, result := range machineResults {
			summary.TotalTests++
			if result.Success {
				summary.SuccessfulTests++
			} else {
				summary.FailedTests++
			}
			summary.TotalDuration += result.Duration

			// Update phase-specific summary
			if summary.PhaseResults[result.Phase] == nil {
				summary.PhaseResults[result.Phase] = &types.PhaseSummary{
					Phase: result.Phase,
				}
			}
			phaseSummary := summary.PhaseResults[result.Phase]
			phaseSummary.TotalTests++
			if result.Success {
				phaseSummary.SuccessfulTests++
			} else {
				phaseSummary.FailedTests++
			}
			phaseSummary.AverageDuration = (phaseSummary.AverageDuration + result.Duration) / 2
		}
	}

	return summary
}

// getRelevantPhases determines which phases are relevant for a machine
func (cm *Manager) getRelevantPhases(machine *configtypes.Machine) []types.ConnectivityTestPhase {
	var phases []types.ConnectivityTestPhase

	// Always test DNS if the machine has a hostname
	if machine.Host != "" {
		phases = append(phases, types.PhaseDNS)
	}

	// Test SSH if the machine has SSH configuration
	if machine.User != "" && (machine.Password != "" || machine.KeyFile != "") {
		phases = append(phases, types.PhaseSSH)
	}

	return phases
}

// testPhase tests a specific connectivity phase
func (cm *Manager) testPhase(ctx context.Context, machine *configtypes.Machine, phase types.ConnectivityTestPhase) types.ConnectivityTestResult {
	start := time.Now()
	result := types.ConnectivityTestResult{
		Machine:   machine,
		Phase:     phase,
		Timestamp: start,
	}

	var err error
	switch phase {
	case types.PhaseDNS:
		err = cm.testDNSResolution(ctx, machine)
	case types.PhaseSSH:
		err = cm.testSSHConnectivity(ctx, machine)
	default:
		err = fmt.Errorf("unknown test phase: %s", phase)
	}

	result.Duration = time.Since(start)
	result.Success = err == nil
	if err != nil {
		result.Error = err.Error()
	}

	return result
}

// testDNSResolution tests DNS resolution for a machine
func (cm *Manager) testDNSResolution(ctx context.Context, machine *configtypes.Machine) error {
	if machine.Host == "" {
		return fmt.Errorf("no hostname specified for machine")
	}

	// Skip DNS resolution for IP addresses
	if isIPAddress(machine.Host) {
		return nil
	}

	// Create a context with timeout
	timeoutCtx, cancel := context.WithTimeout(ctx, cm.options.DNSTimeout)
	defer cancel()

	// Resolve the hostname
	_, err := net.DefaultResolver.LookupHost(timeoutCtx, machine.Host)
	if err != nil {
		return fmt.Errorf("DNS resolution failed for %s: %w", machine.Host, err)
	}

	return nil
}

// testSSHConnectivity tests SSH connectivity to a machine
func (cm *Manager) testSSHConnectivity(ctx context.Context, machine *configtypes.Machine) error {
	if machine.User == "" {
		return fmt.Errorf("no user specified for machine")
	}

	if machine.Password == "" && machine.KeyFile == "" {
		return fmt.Errorf("no authentication method specified (password or key file)")
	}

	// Create SSH client configuration
	config := &ssh.ClientConfig{
		User:            machine.User,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         cm.options.Timeout,
	}

	// Set up authentication
	if machine.Password != "" {
		config.Auth = []ssh.AuthMethod{
			ssh.Password(machine.Password),
		}
	} else if machine.KeyFile != "" {
		key, err := loadPrivateKey(machine.KeyFile)
		if err != nil {
			return fmt.Errorf("failed to load private key: %w", err)
		}
		config.Auth = []ssh.AuthMethod{
			ssh.PublicKeys(key),
		}
	}

	// Connect to the machine
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", machine.Host, machine.Port), config)
	if err != nil {
		return fmt.Errorf("SSH connection failed: %w", err)
	}
	defer client.Close()

	// Execute a simple command to verify connectivity
	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create SSH session: %w", err)
	}
	defer session.Close()

	// Execute the test command
	err = session.Run(cm.options.SSHCommand)
	if err != nil {
		return fmt.Errorf("SSH command execution failed: %w", err)
	}

	return nil
}

// isIPAddress checks if the given host string is an IP address
func isIPAddress(host string) bool {
	return net.ParseIP(host) != nil
}

// loadPrivateKey loads a private key from file
func loadPrivateKey(keyFile string) (ssh.Signer, error) {
	keyBytes, err := os.ReadFile(keyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read key file: %w", err)
	}

	key, err := ssh.ParsePrivateKey(keyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	return key, nil
}
