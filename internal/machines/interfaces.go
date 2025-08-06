package machines

import (
	configtypes "spooky/internal/config/types"
	"spooky/internal/machines/types"
)

// MachineManager defines the main interface for machine operations
type MachineManager interface {
	// Core operations
	BuildIndexes(machines []configtypes.Machine) error
	UpdateIndexes(machines []configtypes.Machine) error
	GetState() *types.IndexManagerState
	Stop() error

	// Lookup operations
	LookupByName(name string) (*configtypes.Machine, bool)
	LookupByHost(host string) (*configtypes.Machine, bool)
	LookupByTag(tagKey string) ([]*configtypes.Machine, bool)
	LookupByTagValue(tagKey, tagValue string) ([]*configtypes.Machine, bool)
	LookupByNetwork(networkType string) ([]*configtypes.Machine, bool)
	LookupBySubnet(subnet string) ([]*configtypes.Machine, bool)
	FilterByTags(criteria map[string]string) []*configtypes.Machine

	// Performance and optimization
	GetIndexMetrics() *types.IndexMetrics
	GetIndexPerformance() *types.IndexPerformanceStats
	OptimizeIndexes() error
	CleanupIndexes() error
	ValidateIndexes() error
}
