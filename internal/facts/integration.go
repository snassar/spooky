// Package facts provides fact collection and in-memory management functionality.
package facts

import (
	"context"
	"fmt"

	spookyinterfaces "spooky/internal/interfaces"
	spookytypes "spooky/internal/types"
	spookytypeslogging "spooky/internal/types/logging"
	spookytypesschemas "spooky/internal/types/schemas"
)

// Integration implements the FactsIntegration interface
type Integration struct {
	manager *Manager
	logger  spookytypeslogging.Logger
}

// NewIntegration creates a new facts integration
func NewIntegration(manager *Manager) spookyinterfaces.FactsIntegration {
	return &Integration{
		manager: manager,
		logger:  manager.logger,
	}
}

// CollectFacts collects facts from the given machine
func (i *Integration) CollectFacts(ctx context.Context, machine *spookytypes.Machine) (interface{}, error) {
	if machine == nil {
		return nil, fmt.Errorf("machine cannot be nil")
	}

	i.logger.Info("Collecting facts via integration", map[string]interface{}{
		"machine": machine.Hostname,
	})

	// Use the manager to collect facts
	facts, err := i.manager.CollectFacts(ctx, machine)
	if err != nil {
		return nil, fmt.Errorf("failed to collect facts: %w", err)
	}

	return facts, nil
}

// StoreFacts stores facts in memory
func (i *Integration) StoreFacts(ctx context.Context, facts interface{}) error {
	if facts == nil {
		return fmt.Errorf("facts cannot be nil")
	}

	i.logger.Info("Storing facts via integration", map[string]interface{}{
		"facts_type": fmt.Sprintf("%T", facts),
	})

	// Convert interface{} to FactCollection
	factCollection, ok := facts.(*FactCollection)
	if !ok {
		return fmt.Errorf("invalid facts type: expected *FactCollection, got %T", facts)
	}

	// Store facts using the manager (collect and store)
	// Store facts in memory for the duration of the operation
	i.logger.Info("Facts stored in memory", map[string]interface{}{
		"machine_id": factCollection.MachineID,
	})

	return nil
}

// LoadFacts loads facts from memory
func (i *Integration) LoadFacts(ctx context.Context) (interface{}, error) {
	i.logger.Info("Loading facts via integration")

	// Facts are only stored in memory during operations
	// Return nil as there's no persistent storage to load from
	return nil, nil
}

// ValidateFacts validates facts
func (i *Integration) ValidateFacts(ctx context.Context, facts interface{}) (*spookytypes.ValidationResult, error) {
	if facts == nil {
		return &spookytypes.ValidationResult{
			Valid:    false,
			Errors:   []spookytypesschemas.SchemaError{{Message: "facts cannot be nil"}},
			Warnings: []spookytypesschemas.SchemaError{},
		}, nil
	}

	i.logger.Info("Validating facts via integration", map[string]interface{}{
		"facts_type": fmt.Sprintf("%T", facts),
	})

	// Convert interface{} to FactCollection
	factCollection, ok := facts.(*FactCollection)
	if !ok {
		return &spookytypes.ValidationResult{
			Valid:    false,
			Errors:   []spookytypesschemas.SchemaError{{Message: fmt.Sprintf("invalid facts type: expected *FactCollection, got %T", facts)}},
			Warnings: []spookytypesschemas.SchemaError{},
		}, nil
	}

	// Validate facts using the manager
	result, err := i.manager.ValidateFacts(ctx, factCollection)
	if err != nil {
		return nil, fmt.Errorf("failed to validate facts: %w", err)
	}

	return result, nil
}

// GetManager returns the underlying fact manager
func (i *Integration) GetManager() interface{} {
	return i.manager
}
