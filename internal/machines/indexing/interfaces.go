package indexing

import (
	configtypes "spooky/internal/config/types"
	"spooky/internal/machines/types"
)

// IndexManager defines the interface for indexing operations
type IndexManager interface {
	// Index operations
	BuildIndexes(machines []configtypes.Machine) error
	UpdateIndexes(machines []configtypes.Machine) error
	SynchronizeIndexes(syncData *types.IndexSyncData) error

	// Lookup operations
	LookupByName(name string) (*configtypes.Machine, bool)
	LookupByHost(host string) (*configtypes.Machine, bool)
	LookupByTag(tagKey string) ([]*configtypes.Machine, bool)
	LookupByTagValue(tagKey, tagValue string) ([]*configtypes.Machine, bool)
	LookupByNetwork(networkType string) ([]*configtypes.Machine, bool)
	LookupBySubnet(subnet string) ([]*configtypes.Machine, bool)
	LookupByGroup(groupName string) ([]*configtypes.Machine, bool)
	LookupByMetadata(metadataKey string) ([]*configtypes.Machine, bool)

	// Filtering operations
	FilterByGroups(groups []string) []*configtypes.Machine
	FilterByTags(criteria map[string]string) []*configtypes.Machine

	// Performance operations
	GetMetrics() *types.IndexMetrics
	OptimizeIndexes()
}

// IndexEngine defines the interface for the index engine
type IndexEngine interface {
	// Core operations
	BuildIndexes(machines []configtypes.Machine) error
	buildIndex(indexType types.IndexType, machines []configtypes.Machine) error

	// Lookup operations
	LookupByName(name string) (*configtypes.Machine, bool)
	LookupByHost(host string) (*configtypes.Machine, bool)
	LookupByTag(tagKey string) ([]*configtypes.Machine, bool)
	LookupByTagValue(tagKey, tagValue string) ([]*configtypes.Machine, bool)
	LookupByNetwork(networkType string) ([]*configtypes.Machine, bool)
	LookupBySubnet(subnet string) ([]*configtypes.Machine, bool)
	LookupByGroup(groupName string) ([]*configtypes.Machine, bool)
	LookupByMetadata(metadataKey string) ([]*configtypes.Machine, bool)

	// Filtering operations
	FilterByGroups(groups []string) []*configtypes.Machine
	FilterByTags(criteria map[string]string) []*configtypes.Machine

	// Performance operations
	GetMetrics() *types.IndexMetrics
	OptimizeIndexes()
	optimizeIndex(indexType types.IndexType, index *types.MachineIndex)
	updateLookupMetrics(indexType types.IndexType, hit bool)
}

// IndexCache defines the interface for index caching
type IndexCache interface {
	Get(key string) (*types.MachineIndex, bool)
	Set(key string, index *types.MachineIndex)
	Clear()
	IsExpired() bool
}
