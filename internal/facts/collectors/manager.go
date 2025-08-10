package collectors

import (
	"fmt"
	spookyfactstypes "spooky/internal/types/facts"
)

// Manager provides collector management functionality
type Manager struct {
	collector  spookyfactstypes.FactCollector
	collectors map[string]spookyfactstypes.FactCollector
}

// NewManager creates a new collector manager
func NewManager(collector spookyfactstypes.FactCollector) *Manager {
	return &Manager{
		collector:  collector,
		collectors: make(map[string]spookyfactstypes.FactCollector),
	}
}

// GetCollector returns the primary collector
func (m *Manager) GetCollector() spookyfactstypes.FactCollector {
	return m.collector
}

// SetCollector sets the primary collector
func (m *Manager) SetCollector(collector spookyfactstypes.FactCollector) {
	m.collector = collector
}

// RegisterCollector registers a custom collector
func (m *Manager) RegisterCollector(name string, collector spookyfactstypes.FactCollector) {
	if collector == nil {
		return
	}
	m.collectors[name] = collector
}

// GetCustomCollector gets a custom collector by name
func (m *Manager) GetCustomCollector(name string) (spookyfactstypes.FactCollector, bool) {
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
