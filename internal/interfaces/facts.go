package interfaces

import (
	"io"
	spookyfactstypes "spooky/internal/facts/types"
	spookystorage "spooky/internal/storage"
	"time"
)

// FactsIntegration defines the interface for facts system integration
type FactsIntegration interface {
	// LoadFacts loads facts for the specified machines
	LoadFacts(machineNames []string) (*FactsContext, error)

	// CollectFacts collects facts from the specified machines
	CollectFacts(machineNames []string) (*FactsContext, error)

	// ValidateFacts validates facts data integrity
	ValidateFacts(facts *FactsContext) error

	// OptimizeFactsGathering optimizes facts gathering performance
	OptimizeFactsGathering(machineNames []string, parallel int) (int, time.Duration, error)

	// GetFactsForMachine gets facts for a specific machine
	GetFactsForMachine(machine string) (*spookyfactstypes.FactCollection, error)

	// CacheFacts caches facts for later use
	CacheFacts(facts *FactsContext) error

	// CollectFactsForAction collects facts needed for action execution
	CollectFactsForAction(action interface{}, machines []string) (*FactsContext, error)

	// ValidateActionWithFacts validates an action using facts data
	ValidateActionWithFacts(action interface{}, factsContext *FactsContext) error

	// GetFactValue gets a specific fact value from the context
	GetFactValue(factKey string, factsContext *FactsContext) (interface{}, error)

	// ExportFacts exports facts to JSON format with encryption support
	ExportFacts(w io.Writer, opts spookystorage.ExportOptions) error

	// ExportToHCL exports facts to HCL format
	ExportToHCL(w io.Writer, query *spookystorage.FactQuery) error

	// ImportFacts imports facts from JSON format
	ImportFacts(r io.Reader) error

	// ImportFromHCL imports facts from HCL format
	ImportFromHCL(r io.Reader) error
}
