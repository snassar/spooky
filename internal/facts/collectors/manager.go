package collectors

import (
	"fmt"
)

// Manager provides collector management functionality
type Manager struct {
	collector  Collector
	collectors map[string]Collector
}

// NewManager creates a new collector manager
func NewManager(collector Collector) *Manager {
	return &Manager{
		collector:  collector,
		collectors: make(map[string]Collector),
	}
}

// GetCollector returns the primary collector
func (m *Manager) GetCollector() Collector {
	return m.collector
}

// SetCollector sets the primary collector
func (m *Manager) SetCollector(collector Collector) {
	m.collector = collector
}

// RegisterCollector registers a custom collector
func (m *Manager) RegisterCollector(name string, collector Collector) {
	if collector == nil {
		return
	}
	m.collectors[name] = collector
}

// GetCustomCollector gets a custom collector by name
func (m *Manager) GetCustomCollector(name string) (Collector, bool) {
	collector, exists := m.collectors[name]
	return collector, exists
}

// RemoveCollector removes a custom collector
func (m *Manager) RemoveCollector(name string) error {
	if _, exists := m.collectors[name]; !exists {
		return fmt.Errorf("collector %s not found", name)
	}
	delete(m.collectors, name)
	return nil
}
