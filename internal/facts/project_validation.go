package facts

import (
	"fmt"
	"regexp"
	"time"

	"spooky/internal/types/facts"
)

// ProjectFactsValidator provides validation for project facts using schema
type ProjectFactsValidator struct {
	utils *ValidationUtils
}

// NewProjectFactsValidator creates a new project facts validator
func NewProjectFactsValidator() *ProjectFactsValidator {
	return &ProjectFactsValidator{
		utils: NewValidationUtils(),
	}
}

// ValidateProjectFacts validates facts against project schema structure
func (pfv *ProjectFactsValidator) ValidateProjectFacts(facts *types.FactCollection) *ValidationResult {
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
	if err := pfv.validateMachineID(facts); err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, ValidationError{
			Field:    "machine_id",
			Message:  err.Error(),
			Severity: "error",
		})
	}

	// Validate project_id (should be present in facts metadata)
	if err := pfv.validateProjectID(facts); err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, ValidationError{
			Field:    "project_id",
			Message:  err.Error(),
			Severity: "error",
		})
	}

	// Validate collected_at timestamp
	if err := pfv.validateTimestamp(facts); err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, ValidationError{
			Field:    "collected_at",
			Message:  err.Error(),
			Severity: "error",
		})
	}

	// Validate TTL format
	if err := pfv.validateTTLFormat(facts); err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, ValidationError{
			Field:    "ttl",
			Message:  err.Error(),
			Severity: "error",
		})
	}

	// Validate facts structure
	if err := pfv.validateFactsStructure(facts); err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, ValidationError{
			Field:    "facts",
			Message:  err.Error(),
			Severity: "error",
		})
	}

	return result
}

// ValidateProjectFactsConstraints validates project facts constraints
func (pfv *ProjectFactsValidator) ValidateProjectFactsConstraints(facts *types.FactCollection) *ValidationResult {
	validators := []func(*types.FactCollection) error{
		pfv.validateProjectTTL,
		pfv.validateProjectIsolation,
		pfv.validateProjectScope,
	}
	return pfv.utils.ValidateConstraints(facts, "facts-constraints", validators)
}

// validateMachineID validates machine_id format and presence
func (pfv *ProjectFactsValidator) validateMachineID(facts *types.FactCollection) error {
	// Look for machine_id in facts
	machineIDFact, exists := facts.Facts["machine_id"]
	if !exists {
		return fmt.Errorf("machine_id is required in project facts")
	}

	machineID, ok := machineIDFact.Value.(string)
	if !ok {
		return fmt.Errorf("machine_id must be a string")
	}

	return pfv.utils.ValidateMachineID(machineID)
}

// validateProjectID validates project_id presence and format
func (pfv *ProjectFactsValidator) validateProjectID(facts *types.FactCollection) error {
	// Check if project_id is present in any fact's metadata
	projectIDFound := false
	var projectID string

	for _, fact := range facts.Facts {
		if fact.Metadata != nil {
			if pid, exists := fact.Metadata["project_id"]; exists {
				projectIDFound = true
				if pidStr, ok := pid.(string); ok {
					projectID = pidStr
					break
				}
			}
		}
	}

	if !projectIDFound {
		return fmt.Errorf("project_id is required in project facts metadata")
	}

	// Validate project_id format
	if projectID == "" {
		return fmt.Errorf("project_id cannot be empty")
	}

	// Project ID should be alphanumeric with underscores and hyphens
	projectIDPattern := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	if !projectIDPattern.MatchString(projectID) {
		return fmt.Errorf("project_id contains invalid characters: %s", projectID)
	}

	return nil
}

// validateTimestamp validates collected_at timestamp
func (pfv *ProjectFactsValidator) validateTimestamp(facts *types.FactCollection) error {
	// Check if timestamp is set
	if facts.Timestamp.IsZero() {
		return fmt.Errorf("collected_at timestamp is required")
	}

	return pfv.utils.ValidateTimestamp(facts.Timestamp)
}

// validateTTLFormat validates TTL format
func (pfv *ProjectFactsValidator) validateTTLFormat(facts *types.FactCollection) error {
	// Check each fact's TTL
	for key, fact := range facts.Facts {
		if fact.TTL > 0 {
			if err := pfv.utils.ValidateTTL(fact.TTL); err != nil {
				return fmt.Errorf("invalid TTL for fact %s: %w", key, err)
			}
		}
	}

	return nil
}

// validateFactsStructure validates the structure of facts
func (pfv *ProjectFactsValidator) validateFactsStructure(facts *types.FactCollection) error {
	return pfv.utils.ValidateFactStructure(facts.Facts)
}

// validateProjectTTL validates project facts TTL constraints
func (pfv *ProjectFactsValidator) validateProjectTTL(facts *types.FactCollection) error {
	// Project facts should have shorter TTL (1h default)
	for key, fact := range facts.Facts {
		if fact.TTL > 0 && fact.TTL > 7*24*time.Hour {
			return fmt.Errorf("project fact %s has very long TTL (%s), consider using shorter TTL for project facts",
				key, fact.TTL.String())
		}
	}

	return nil
}

// validateProjectIsolation validates project isolation constraints
func (pfv *ProjectFactsValidator) validateProjectIsolation(facts *types.FactCollection) error {
	// All facts in a project should have the same project_id
	var projectID string
	projectIDSet := false

	for key, fact := range facts.Facts {
		if fact.Metadata != nil {
			if pid, exists := fact.Metadata["project_id"]; exists {
				if pidStr, ok := pid.(string); ok {
					if !projectIDSet {
						projectID = pidStr
						projectIDSet = true
					} else if pidStr != projectID {
						return fmt.Errorf("project fact %s has different project_id (%s vs %s), all facts must belong to the same project",
							key, pidStr, projectID)
					}
				}
			}
		}
	}

	return nil
}

// validateProjectScope validates project scope constraints
func (pfv *ProjectFactsValidator) validateProjectScope(facts *types.FactCollection) error {
	// Project facts should not be shared across projects
	// This is more of a logical validation - project facts should be isolated

	// For now, just ensure all facts have project_id metadata
	for key, fact := range facts.Facts {
		if fact.Metadata == nil {
			return fmt.Errorf("project fact %s missing metadata", key)
		}

		if _, exists := fact.Metadata["project_id"]; !exists {
			return fmt.Errorf("project fact %s missing project_id in metadata", key)
		}
	}

	return nil
}

// ValidateProjectMachineIDUniqueness validates machine ID uniqueness within a project
func (pfv *ProjectFactsValidator) ValidateProjectMachineIDUniqueness(machineID, projectID string, existingCollections []*types.FactCollection) error {
	for _, collection := range existingCollections {
		// Check if collection belongs to this project
		if projectFact, exists := collection.Facts["project.id"]; exists {
			if project, ok := projectFact.Value.(string); ok && project == projectID {
				// Check machine ID
				if existingMachineID, exists := collection.Facts["machine_id"]; exists {
					if id, ok := existingMachineID.Value.(string); ok && id == machineID {
						return fmt.Errorf("machine ID %s is already used by %s in project %s", machineID, collection.Server, projectID)
					}
				}
			}
		}
	}
	return nil
}

// ValidateProjectFactsImport validates facts for import into project storage
func (pfv *ProjectFactsValidator) ValidateProjectFactsImport(facts *types.FactCollection, projectID, importSource string) *ValidationResult {
	result := pfv.ValidateProjectFacts(facts)

	// Additional import-specific validations
	if projectID != "" {
		// Ensure all facts have the correct project_id
		for _, fact := range facts.Facts {
			if fact.Metadata == nil {
				fact.Metadata = make(map[string]interface{})
			}
			fact.Metadata["project_id"] = projectID
		}
	}

	if importSource != "" {
		// Validate import source format
		if err := pfv.validateImportSource(importSource); err != nil {
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
func (pfv *ProjectFactsValidator) validateImportSource(source string) error {
	return pfv.utils.ValidateImportSource(source)
}

// ValidateProjectFactsConsistency validates consistency across project facts
func (pfv *ProjectFactsValidator) ValidateProjectFactsConsistency(facts *types.FactCollection) *ValidationResult {
	result := &ValidationResult{
		Valid:  true,
		Schema: "facts-consistency",
	}

	// Check for consistent server identifiers
	if err := pfv.validateServerConsistency(facts); err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, ValidationError{
			Field:    "server_consistency",
			Message:  err.Error(),
			Severity: "warning",
		})
	}

	// Check for consistent timestamps
	if err := pfv.validateTimestampConsistency(facts); err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, ValidationError{
			Field:    "timestamp_consistency",
			Message:  err.Error(),
			Severity: "warning",
		})
	}

	return result
}

// validateServerConsistency validates server identifier consistency
func (pfv *ProjectFactsValidator) validateServerConsistency(facts *types.FactCollection) error {
	var serverID string
	serverIDSet := false

	for key, fact := range facts.Facts {
		if !serverIDSet {
			serverID = fact.Server
			serverIDSet = true
		} else if fact.Server != serverID {
			return fmt.Errorf("inconsistent server identifiers: %s vs %s in fact %s",
				serverID, fact.Server, key)
		}
	}

	return nil
}

// validateTimestampConsistency validates timestamp consistency
func (pfv *ProjectFactsValidator) validateTimestampConsistency(facts *types.FactCollection) error {
	// Check if all facts have similar timestamps (within 5 minutes)
	if len(facts.Facts) < 2 {
		return nil // No consistency check needed for single fact
	}

	var timestamps []time.Time
	for _, fact := range facts.Facts {
		timestamps = append(timestamps, fact.Timestamp)
	}

	// Check if timestamps are within reasonable range
	for i := 1; i < len(timestamps); i++ {
		diff := timestamps[i].Sub(timestamps[0])
		if diff < 0 {
			diff = -diff
		}
		if diff > 5*time.Minute {
			return fmt.Errorf("fact timestamps are inconsistent (difference: %s)", diff.String())
		}
	}

	return nil
}
