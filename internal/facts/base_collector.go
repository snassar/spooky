package facts

import (
	"fmt"
	"time"
)

// BaseCollector provides common functionality for collectors
type BaseCollector struct{}

// NewBaseCollector creates a new base collector
func NewBaseCollector() *BaseCollector {
	return &BaseCollector{}
}

// CollectAll provides a common implementation for collecting all facts
func (b *BaseCollector) CollectAll(server string, collector interface {
	collectSystemFacts(*FactCollection) error
	collectOSFacts(*FactCollection) error
	collectHardwareFacts(*FactCollection) error
	collectNetworkFacts(*FactCollection) error
	collectEnvironmentFacts(*FactCollection) error
}) (*FactCollection, error) {
	collection := &FactCollection{
		Server:    server,
		Timestamp: time.Now(),
		Facts:     make(map[string]*Fact),
	}

	// Collect system facts
	if err := collector.collectSystemFacts(collection); err != nil {
		return nil, fmt.Errorf("failed to collect system facts: %w", err)
	}

	// Collect OS facts
	if err := collector.collectOSFacts(collection); err != nil {
		return nil, fmt.Errorf("failed to collect OS facts: %w", err)
	}

	// Collect hardware facts
	if err := collector.collectHardwareFacts(collection); err != nil {
		return nil, fmt.Errorf("failed to collect hardware facts: %w", err)
	}

	// Collect network facts
	if err := collector.collectNetworkFacts(collection); err != nil {
		return nil, fmt.Errorf("failed to collect network facts: %w", err)
	}

	// Collect environment facts
	if err := collector.collectEnvironmentFacts(collection); err != nil {
		return nil, fmt.Errorf("failed to collect environment facts: %w", err)
	}

	return collection, nil
}
