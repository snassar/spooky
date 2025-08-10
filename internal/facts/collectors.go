package facts

import (
	"fmt"

	spookyfactscollectorslocal "spooky/internal/facts/collectors/local"
	spookyinterfaces "spooky/internal/interfaces"
	spookytypesfacts "spooky/internal/types/facts"
	"time"
)

// NewSSHCollector creates a new SSH collector
func NewSSHCollector(sshClient interface{}) spookytypesfacts.FactCollector {
	// TODO: Fix logger interface compatibility issue
	// For now, return a placeholder implementation
	return &placeholderSSHCollector{}
}

// placeholderSSHCollector is a temporary implementation until logger interface is fixed
type placeholderSSHCollector struct{}

func (p *placeholderSSHCollector) Collect(server string) (*spookytypesfacts.FactCollection, error) {
	return &spookytypesfacts.FactCollection{
		Server:    server,
		Timestamp: time.Now(),
		Facts:     make(map[string]*spookytypesfacts.Fact),
	}, nil
}

func (p *placeholderSSHCollector) CollectSpecific(server string, keys []string) (*spookytypesfacts.FactCollection, error) {
	return p.Collect(server)
}

func (p *placeholderSSHCollector) GetFact(server, key string) (*spookytypesfacts.Fact, error) {
	return nil, fmt.Errorf("placeholder SSH collector - not implemented")
}

// NewLocalCollector creates a new local collector
func NewLocalCollector() spookytypesfacts.FactCollector {
	return spookyfactscollectorslocal.NewCollector()
}

// NewJSONCollector creates a new JSON collector
func NewJSONCollector(source string, mergePolicy spookytypesfacts.MergePolicy) spookytypesfacts.FactCollector {
	// TODO: Implement real JSON collector
	return nil
}

// NewHTTPCollector creates a new HTTP collector
func NewHTTPCollector(source string, headers map[string]string, timeout time.Duration, mergePolicy spookytypesfacts.MergePolicy) spookytypesfacts.FactCollector {
	// TODO: Implement real HTTP collector
	return nil
}

// NewHCLCollector creates a new HCL collector
func NewHCLCollector(filePath string, logger spookyinterfaces.Logger, mergePolicy spookytypesfacts.MergePolicy) spookytypesfacts.FactCollector {
	// TODO: Implement real HCL collector
	return nil
}

// NewOpenTofuCollector creates a new OpenTofu collector
func NewOpenTofuCollector(statePath string, logger spookyinterfaces.Logger, mergePolicy spookytypesfacts.MergePolicy) spookytypesfacts.FactCollector {
	// TODO: Implement real OpenTofu collector
	return nil
}
