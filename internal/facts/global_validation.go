package facts

import (
	"fmt"
	"time"

	"spooky/internal/facts/types"
)

// GlobalFactsValidator provides validation for global facts using schema
type GlobalFactsValidator struct {
	utils *ValidationUtils
}

// NewGlobalFactsValidator creates a new global facts validator
func NewGlobalFactsValidator() *GlobalFactsValidator {
	return &GlobalFactsValidator{
		utils: NewValidationUtils(),
	}
}

// ValidationResult represents validation result
type ValidationResult struct {
	Valid  bool
	Schema string
	Errors []ValidationError
}

// ValidationError represents a validation error
type ValidationError struct {
	Field    string
	Message  string
	Severity string
}

// ValidateGlobalFacts validates facts against global schema structure
func (gfv *GlobalFactsValidator) ValidateGlobalFacts(facts *types.FactCollection) *ValidationResult {
	result := &ValidationResult{
		Valid:  true,
		Schema: "facts-structure",
	}

	if facts == nil {
		result.Valid = false
		result.Errors = append(result.Errors, ValidationError{
			Field:    "facts",
			Message:  "facts collection cannot be nil",
			Severity: "error",
		})
		return result
	}

	// Validate machine_id (should be present in facts)
	if err := gfv.validateMachineID(facts); err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, ValidationError{
			Field:    "machine_id",
			Message:  err.Error(),
			Severity: "error",
		})
	}

	// Validate collected_at timestamp
	if err := gfv.validateTimestamp(facts); err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, ValidationError{
			Field:    "collected_at",
			Message:  err.Error(),
			Severity: "error",
		})
	}

	// Validate TTL format
	if err := gfv.validateTTLFormat(facts); err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, ValidationError{
			Field:    "ttl",
			Message:  err.Error(),
			Severity: "error",
		})
	}

	// Validate facts structure
	if err := gfv.validateFactsStructure(facts); err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, ValidationError{
			Field:    "facts",
			Message:  err.Error(),
			Severity: "error",
		})
	}

	return result
}

// ValidateGlobalFactsConstraints validates global facts constraints
func (gfv *GlobalFactsValidator) ValidateGlobalFactsConstraints(facts *types.FactCollection) *ValidationResult {
	validators := []func(*types.FactCollection) error{
		gfv.validateGlobalTTL,
		gfv.validatePersistence,
		gfv.validateSharedAccess,
	}
	return gfv.utils.ValidateConstraints(facts, "facts-constraints", validators)
}

// validateMachineID validates machine_id format and presence
func (gfv *GlobalFactsValidator) validateMachineID(facts *types.FactCollection) error {
	// Look for machine_id in facts
	machineIDFact, exists := facts.Facts["machine_id"]
	if !exists {
		return fmt.Errorf("machine_id is required in global facts")
	}

	machineID, ok := machineIDFact.Value.(string)
	if !ok {
		return fmt.Errorf("machine_id must be a string")
	}

	return gfv.utils.ValidateMachineID(machineID)
}

// validateTimestamp validates collected_at timestamp
func (gfv *GlobalFactsValidator) validateTimestamp(facts *types.FactCollection) error {
	// Check if timestamp is set
	if facts.Timestamp.IsZero() {
		return fmt.Errorf("collected_at timestamp is required")
	}

	return gfv.utils.ValidateTimestamp(facts.Timestamp)
}

// validateTTLFormat validates TTL format
func (gfv *GlobalFactsValidator) validateTTLFormat(facts *types.FactCollection) error {
	// Check each fact's TTL
	for key, fact := range facts.Facts {
		if fact.TTL > 0 {
			if err := gfv.utils.ValidateTTL(fact.TTL); err != nil {
				return fmt.Errorf("invalid TTL for fact %s: %w", key, err)
			}
		}
	}

	return nil
}

// validateFactsStructure validates the structure of facts
func (gfv *GlobalFactsValidator) validateFactsStructure(facts *types.FactCollection) error {
	return gfv.utils.ValidateFactStructure(facts.Facts)
}

// validateGlobalTTL validates global facts TTL constraints
func (gfv *GlobalFactsValidator) validateGlobalTTL(facts *types.FactCollection) error {
	// Global facts should have longer TTL (24h default)
	defaultGlobalTTL := 24 * time.Hour

	for key, fact := range facts.Facts {
		if fact.TTL > 0 && fact.TTL < defaultGlobalTTL {
			return fmt.Errorf("global fact %s has short TTL (%s), consider using longer TTL for global facts",
				key, fact.TTL.String())
		}
	}

	return nil
}

// validatePersistence validates persistence constraints
func (gfv *GlobalFactsValidator) validatePersistence(facts *types.FactCollection) error {
	// Global facts should be persistent (not have very short TTL)
	for key, fact := range facts.Facts {
		if fact.TTL > 0 && fact.TTL < time.Hour {
			return fmt.Errorf("global fact %s has very short TTL (%s), global facts should be persistent",
				key, fact.TTL.String())
		}
	}

	return nil
}

// validateSharedAccess validates shared access constraints
func (gfv *GlobalFactsValidator) validateSharedAccess(facts *types.FactCollection) error {
	// Global facts should not be project-specific
	// This is more of a logical validation - global facts should be accessible across projects
	// Implementation depends on how project isolation is handled

	// For now, just ensure facts don't have project-specific metadata
	for key, fact := range facts.Facts {
		if fact.Metadata != nil {
			if _, exists := fact.Metadata["project_id"]; exists {
				return fmt.Errorf("global fact %s contains project-specific metadata, global facts should be shared", key)
			}
		}
	}

	return nil
}

// ValidateMachineIDUniqueness validates that a machine ID is unique across all facts
func (gfv *GlobalFactsValidator) ValidateMachineIDUniqueness(machineID string, existingCollections []*types.FactCollection) error {
	for _, collection := range existingCollections {
		if existingMachineID, exists := collection.Facts["machine_id"]; exists {
			if id, ok := existingMachineID.Value.(string); ok && id == machineID {
				return fmt.Errorf("machine ID %s is already used by %s", machineID, collection.Server)
			}
		}
	}
	return nil
}

// ValidateGlobalFactsImport validates facts for import into global storage
func (gfv *GlobalFactsValidator) ValidateGlobalFactsImport(facts *types.FactCollection, importSource string) *ValidationResult {
	result := gfv.ValidateGlobalFacts(facts)

	// Additional import-specific validations
	if importSource != "" {
		// Validate import source format
		if err := gfv.validateImportSource(importSource); err != nil {
			result.Valid = false
			result.Errors = append(result.Errors, ValidationError{
				Field:    "import_source",
				Message:  err.Error(),
				Severity: "error",
			})
		}
	}

	return result
}

// validateImportSource validates import source format
func (gfv *GlobalFactsValidator) validateImportSource(source string) error {
	return gfv.utils.ValidateImportSource(source)
}
