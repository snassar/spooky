package facts

import (
	"fmt"
	"time"
)

// StubCollector provides a base implementation for collectors that are not yet implemented
type StubCollector struct {
	name string
}

// NewStubCollector creates a new stub collector
func NewStubCollector(name string) *StubCollector {
	return &StubCollector{name: name}
}

// Collect gathers facts (stub implementation)
func (c *StubCollector) Collect(server string) (*FactCollection, error) {
	return &FactCollection{
		Server:    server,
		Timestamp: time.Now(),
		Facts:     make(map[string]*Fact),
	}, nil
}

// CollectSpecific collects only the specified facts (stub implementation)
func (c *StubCollector) CollectSpecific(server string, _ []string) (*FactCollection, error) {
	return &FactCollection{
		Server:    server,
		Timestamp: time.Now(),
		Facts:     make(map[string]*Fact),
	}, nil
}

// GetFact retrieves a single fact (stub implementation)
func (c *StubCollector) GetFact(_, _ string) (*Fact, error) {
	return nil, fmt.Errorf("%s fact collection not yet implemented", c.name)
}
