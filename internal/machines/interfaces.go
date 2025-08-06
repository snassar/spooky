package machines

import (
	spookyconfigtypes "spooky/internal/config/types"
	spookymachinestypes "spooky/internal/machines/types"
)

// MachineManager defines the main interface for machine operations
type MachineManager interface {
	// Core operations
	BuildIndexes(machines []spookyconfigtypes.Machine) error
	UpdateIndexes(machines []spookyconfigtypes.Machine) error
	GetState() *spookymachinestypes.IndexManagerState
	Stop() error

	// Lookup operations
	LookupByName(name string) (*spookyconfigtypes.Machine, bool)
	LookupByHost(host string) (*spookyconfigtypes.Machine, bool)
	LookupByTag(tagKey string) ([]*spookyconfigtypes.Machine, bool)
	LookupByTagValue(tagKey, tagValue string) ([]*spookyconfigtypes.Machine, bool)
	LookupByNetwork(networkType string) ([]*spookyconfigtypes.Machine, bool)
	LookupBySubnet(subnet string) ([]*spookyconfigtypes.Machine, bool)
	FilterByTags(criteria map[string]string) []*spookyconfigtypes.Machine

	// Performance and optimization
	GetIndexMetrics() *spookymachinestypes.IndexMetrics
	GetIndexPerformance() *spookymachinestypes.IndexPerformanceStats
	OptimizeIndexes() error
	CleanupIndexes() error
	ValidateIndexes() error
}
