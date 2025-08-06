package collectors

import (
	"spooky/internal/facts/types"
)

// Collector defines the interface for fact collection strategies
type Collector interface {
	// Core collection methods
	Collect(server string) (*types.FactCollection, error)
	CollectSpecific(server string, keys []string) (*types.FactCollection, error)
	GetFact(server, key string) (*types.Fact, error)

	// Configuration
	GetSource() types.FactSource
	GetMergePolicy() types.MergePolicy
	SetMergePolicy(policy types.MergePolicy)

	// Validation
	Validate() error
}

// BaseCollector provides common functionality for collectors
type BaseCollector struct {
	source      types.FactSource
	mergePolicy types.MergePolicy
}

// GetSource returns the fact source
func (bc *BaseCollector) GetSource() types.FactSource {
	return bc.source
}

// GetMergePolicy returns the merge policy
func (bc *BaseCollector) GetMergePolicy() types.MergePolicy {
	return bc.mergePolicy
}

// SetMergePolicy sets the merge policy
func (bc *BaseCollector) SetMergePolicy(policy types.MergePolicy) {
	bc.mergePolicy = policy
}

// Validate validates the collector configuration
func (bc *BaseCollector) Validate() error {
	// Base validation - can be overridden by specific collectors
	return nil
}

// NewBaseCollector creates a new base collector
func NewBaseCollector(source types.FactSource, mergePolicy types.MergePolicy) *BaseCollector {
	return &BaseCollector{
		source:      source,
		mergePolicy: mergePolicy,
	}
}

// NewSSHCollector creates a new SSH collector
func NewSSHCollector(sshClient interface{}) Collector {
	// This is a placeholder - the actual implementation would be in the ssh package
	// For now, return nil to avoid compilation errors
	return nil
}
