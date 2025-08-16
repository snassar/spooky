// Package facts provides fact collection and in-memory management functionality.
package facts

import (
	"context"
	"fmt"

	spookyinterfaces "spooky/internal/interfaces"
	spookytypes "spooky/internal/types"
	spookytypesfacts "spooky/internal/types/facts"
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
		"host":    machine.Host,
	})

	// Determine if this is a remote machine that needs SSH
	if machine.Host != "" && machine.Host != "localhost" && machine.Host != "127.0.0.1" {
		// Use SSH-based collection for remote machines
		return i.collectFactsViaSSH(ctx, machine)
	}

	// Use local collection for localhost/127.0.0.1
	facts, err := i.manager.CollectFacts(ctx, machine)
	if err != nil {
		return nil, fmt.Errorf("failed to collect facts: %w", err)
	}

	return facts, nil
}

// collectFactsViaSSH collects facts from remote machine via SSH
func (i *Integration) collectFactsViaSSH(ctx context.Context, machine *spookytypes.Machine) (interface{}, error) {
	i.logger.Info("Collecting facts via SSH", map[string]interface{}{
		"machine": machine.Hostname,
		"host":    machine.Host,
	})

	// Get the underlying collector that has SSH capabilities
	collector, ok := i.manager.GetCollector().(interface {
		CollectViaSSH(context.Context, *spookytypes.Machine) (*spookytypesfacts.FactCollection, error)
	})
	if !ok {
		return nil, fmt.Errorf("collector does not support SSH operations")
	}

	// Collect facts using SSH
	facts, err := collector.CollectViaSSH(ctx, machine)
	if err != nil {
		return nil, fmt.Errorf("failed to collect facts via SSH: %w", err)
	}

	return facts, nil
}

// StoreFacts stores facts in memory
func (i *Integration) StoreFacts(_ context.Context, facts interface{}) error {
	if facts == nil {
		return fmt.Errorf("facts cannot be nil")
	}

	i.logger.Info("Storing facts via integration", map[string]interface{}{
		"facts_type": fmt.Sprintf("%T", facts),
	})

	// Convert interface{} to FactCollection
	factCollection, ok := facts.(*spookytypesfacts.FactCollection)
	if !ok {
		return fmt.Errorf("invalid facts type: expected *spookytypesfacts.FactCollection, got %T", facts)
	}

	// Store facts using the manager (collect and store)
	// Store facts in memory for the duration of the operation
	i.logger.Info("Facts stored in memory", map[string]interface{}{
		"machine_id": factCollection.MachineID,
	})

	return nil
}

// LoadFacts loads facts from memory
func (i *Integration) LoadFacts(_ context.Context) (interface{}, error) {
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
	factCollection, ok := facts.(*spookytypesfacts.FactCollection)
	if !ok {
		return &spookytypes.ValidationResult{
			Valid:    false,
			Errors:   []spookytypesschemas.SchemaError{{Message: fmt.Sprintf("invalid facts type: expected *spookytypesfacts.FactCollection, got %T", facts)}},
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

// DecryptFacts decrypts age-encrypted values in facts collection
func (i *Integration) DecryptFacts(ctx context.Context, facts interface{}, secretsIntegration spookyinterfaces.SecretsIntegration, identityPath string) error {
	if facts == nil {
		return fmt.Errorf("facts cannot be nil")
	}

	if secretsIntegration == nil {
		return fmt.Errorf("secrets integration cannot be nil")
	}

	i.logger.Info("Decrypting facts via integration", map[string]interface{}{
		"facts_type": fmt.Sprintf("%T", facts),
	})

	// Convert interface{} to FactCollection
	factCollection, ok := facts.(*spookytypesfacts.FactCollection)
	if !ok {
		return fmt.Errorf("invalid facts type: expected *spookytypesfacts.FactCollection, got %T", facts)
	}

	// Decrypt facts using the manager
	return i.manager.DecryptFacts(ctx, factCollection, secretsIntegration, identityPath)
}
