package facts

import (
	"fmt"

	spookyfactscollectorslocal "spooky/internal/facts/collectors/local"
	spookyfactscollectorsssh "spooky/internal/facts/collectors/ssh"
	spookyfactstypes "spooky/internal/types/facts"
	spookylogging "spooky/internal/logging"
	spookysshtypes "spooky/internal/ssh/types"
	"time"
)

// NewSSHCollector creates a new SSH collector
func NewSSHCollector(sshClient interface{}) spookyfactstypes.FactCollector {
	// Create an adapter that implements the FactCollector interface
	logger := spookylogging.GetLogger()
	sshCollector := spookyfactscollectorsssh.NewCollector(logger)
	return &sshCollectorAdapter{
		collector: sshCollector,
		logger:    logger,
	}
}

// sshCollectorAdapter adapts the SSH collector to the FactCollector interface
type sshCollectorAdapter struct {
	collector *spookyfactscollectorsssh.Collector
	logger    spookylogging.Logger
}

func (a *sshCollectorAdapter) Collect(server string) (*spookyfactstypes.FactCollection, error) {
	// Create a default SSH config for the server
	config := &spookysshtypes.SSHConfig{
		Host:     server,
		Port:     22,
		Username: "root", // Default user, should be configurable
		Timeout:  30 * time.Second,
	}

	return a.collector.Collect(server, config)
}

func (a *sshCollectorAdapter) CollectSpecific(server string, keys []string) (*spookyfactstypes.FactCollection, error) {
	// For now, collect all facts and filter later
	// TODO: Implement selective collection
	return a.Collect(server)
}

func (a *sshCollectorAdapter) GetFact(server, key string) (*spookyfactstypes.Fact, error) {
	collection, err := a.Collect(server)
	if err != nil {
		return nil, err
	}

	if fact, exists := collection.Facts[key]; exists {
		return fact, nil
	}

	return nil, fmt.Errorf("fact '%s' not found for server '%s'", key, server)
}

// NewLocalCollector creates a new local collector
func NewLocalCollector() spookyfactstypes.FactCollector {
	return spookyfactscollectorslocal.NewCollector()
}

// NewJSONCollector creates a new JSON collector
func NewJSONCollector(source string, mergePolicy spookyfactstypes.MergePolicy) spookyfactstypes.FactCollector {
	// TODO: Implement real JSON collector
	return nil
}

// NewHTTPCollector creates a new HTTP collector
func NewHTTPCollector(source string, headers map[string]string, timeout time.Duration, mergePolicy spookyfactstypes.MergePolicy) spookyfactstypes.FactCollector {
	// TODO: Implement real HTTP collector
	return nil
}

// NewHCLCollector creates a new HCL collector
func NewHCLCollector(filePath string, logger spookylogging.Logger, mergePolicy spookyfactstypes.MergePolicy) spookyfactstypes.FactCollector {
	// TODO: Implement real HCL collector
	return nil
}

// NewOpenTofuCollector creates a new OpenTofu collector
func NewOpenTofuCollector(statePath string, logger spookylogging.Logger, mergePolicy spookyfactstypes.MergePolicy) spookyfactstypes.FactCollector {
	// TODO: Implement real OpenTofu collector
	return nil
}
