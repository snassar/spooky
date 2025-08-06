package connectivity

import (
	"context"

	configtypes "spooky/internal/config/types"
	"spooky/internal/machines/types"
)

// ConnectivityManager defines the interface for connectivity operations
type ConnectivityManager interface {
	// Connectivity testing
	TestMachineConnectivity(ctx context.Context, machine *configtypes.Machine) []types.ConnectivityTestResult
	TestMachinesConnectivity(ctx context.Context, machines []*configtypes.Machine) map[string][]types.ConnectivityTestResult
	TestSpecificPhase(ctx context.Context, machine *configtypes.Machine, phase types.ConnectivityTestPhase) types.ConnectivityTestResult
	TestPhases(ctx context.Context, machine *configtypes.Machine, phases []types.ConnectivityTestPhase) []types.ConnectivityTestResult

	// Summary and reporting
	GetTestSummary(results map[string][]types.ConnectivityTestResult) *types.ConnectivityTestSummary
}

// ConnectivityTester defines the interface for connectivity testing
type ConnectivityTester interface {
	// Core testing operations
	TestMachineConnectivity(ctx context.Context, machine *configtypes.Machine) []types.ConnectivityTestResult
	TestMachinesConnectivity(ctx context.Context, machines []*configtypes.Machine) map[string][]types.ConnectivityTestResult
	TestSpecificPhase(ctx context.Context, machine *configtypes.Machine, phase types.ConnectivityTestPhase) types.ConnectivityTestResult
	TestPhases(ctx context.Context, machine *configtypes.Machine, phases []types.ConnectivityTestPhase) []types.ConnectivityTestResult

	// Utility operations
	getRelevantPhases(machine *configtypes.Machine) []types.ConnectivityTestPhase
	testPhase(ctx context.Context, machine *configtypes.Machine, phase types.ConnectivityTestPhase) types.ConnectivityTestResult
	testDNSResolution(ctx context.Context, machine *configtypes.Machine) error
	testSSHConnectivity(ctx context.Context, machine *configtypes.Machine) error
}
