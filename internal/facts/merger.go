package facts

import (
	"fmt"
	"strings"
	"time"

	spookyinterfaces "spooky/internal/interfaces"
	"spooky/internal/types"
	spookyfactstypes "spooky/internal/types/facts"
	spookytypesfacts "spooky/internal/types/facts"
)

// FactMerger handles merging of fact collections
type FactMerger struct {
	logger      spookyinterfaces.Logger
	mergePolicy spookyfactstypes.MergePolicy
}

// NewFactMerger creates a new fact merger with the specified policy
func NewFactMerger(policy spookyfactstypes.MergePolicy) *FactMerger {
	return &FactMerger{
		mergePolicy: policy,
	}
}

// MergeCollections merges multiple fact collections according to the merge policy
func (m *FactMerger) MergeCollections(existing, incoming *spookytypesfacts.FactCollection) (*spookytypesfacts.FactCollection, error) {
	if existing == nil {
		return incoming, nil
	}
	if incoming == nil {
		return existing, nil
	}

	switch m.mergePolicy {
	case types.MergePolicyReplace:
		return m.mergeReplace(existing, incoming)
	case types.MergePolicyMerge:
		return m.mergeCombine(existing, incoming)
	case types.MergePolicySkip:
		return m.mergeSkip(existing, incoming)
	case types.MergePolicyAppend:
		return m.mergeAppend(existing, incoming)
	default:
		return nil, fmt.Errorf("unknown merge policy: %s", m.mergePolicy)
	}
}

// mergeReplace replaces existing facts with new ones
func (m *FactMerger) mergeReplace(existing, incoming *types.FactCollection) (*types.FactCollection, error) {
	merged := &types.FactCollection{
		Server:    existing.Server,
		Timestamp: time.Now(),
		Facts:     make(map[string]*types.Fact),
	}

	// Start with existing facts
	for key, fact := range existing.Facts {
		merged.Facts[key] = fact
	}

	// Replace/add incoming facts
	for key, fact := range incoming.Facts {
		merged.Facts[key] = fact
	}

	return merged, nil
}

// mergeCombine combines facts, preferring newer ones on conflicts
func (m *FactMerger) mergeCombine(existing, incoming *types.FactCollection) (*types.FactCollection, error) {
	merged := &types.FactCollection{
		Server:    existing.Server,
		Timestamp: time.Now(),
		Facts:     make(map[string]*types.Fact),
	}

	// Start with existing facts
	for key, fact := range existing.Facts {
		merged.Facts[key] = fact
	}

	// Add incoming facts, preferring newer ones on conflicts
	for key, newFact := range incoming.Facts {
		if existingFact, exists := merged.Facts[key]; exists {
			// Conflict resolution: prefer newer fact
			if newFact.Timestamp.After(existingFact.Timestamp) {
				merged.Facts[key] = newFact
			}
			// If timestamps are equal, prefer the one with more metadata
			if newFact.Timestamp.Equal(existingFact.Timestamp) && len(newFact.Metadata) > len(existingFact.Metadata) {
				merged.Facts[key] = newFact
			}
		} else {
			merged.Facts[key] = newFact
		}
	}

	return merged, nil
}

// mergeSkip skips facts that already exist
func (m *FactMerger) mergeSkip(existing, incoming *types.FactCollection) (*types.FactCollection, error) {
	merged := &types.FactCollection{
		Server:    existing.Server,
		Timestamp: time.Now(),
		Facts:     make(map[string]*types.Fact),
	}

	// Start with existing facts
	for key, fact := range existing.Facts {
		merged.Facts[key] = fact
	}

	// Only add incoming facts that don't already exist
	for key, fact := range incoming.Facts {
		if _, exists := merged.Facts[key]; !exists {
			merged.Facts[key] = fact
		}
	}

	return merged, nil
}

// mergeAppend appends new facts with suffixes to avoid conflicts
func (m *FactMerger) mergeAppend(existing, incoming *types.FactCollection) (*types.FactCollection, error) {
	merged := &types.FactCollection{
		Server:    existing.Server,
		Timestamp: time.Now(),
		Facts:     make(map[string]*types.Fact),
	}

	// Start with existing facts
	for key, fact := range existing.Facts {
		merged.Facts[key] = fact
	}

	// Add incoming facts, appending suffixes for conflicts
	for key, fact := range incoming.Facts {
		newKey := key
		suffix := 1

		// Find a unique key
		for {
			if _, exists := merged.Facts[newKey]; !exists {
				break
			}
			newKey = fmt.Sprintf("%s_%d", key, suffix)
			suffix++
		}

		// Create a copy of the fact with the new key
		appendedFact := &types.Fact{
			Key:       newKey,
			Value:     fact.Value,
			Source:    fact.Source,
			Server:    fact.Server,
			Timestamp: fact.Timestamp,
			TTL:       fact.TTL,
			Metadata:  make(map[string]interface{}),
		}

		// Copy metadata and add append information
		for k, v := range fact.Metadata {
			appendedFact.Metadata[k] = v
		}
		appendedFact.Metadata["original_key"] = key
		appendedFact.Metadata["appended"] = true
		appendedFact.Metadata["suffix"] = suffix - 1

		merged.Facts[newKey] = appendedFact
	}

	return merged, nil
}

// MergeFacts merges individual facts according to the merge policy
func (m *FactMerger) MergeFacts(existing, incoming *types.Fact) (*types.Fact, error) {
	if existing == nil {
		return incoming, nil
	}
	if incoming == nil {
		return existing, nil
	}

	switch m.mergePolicy {
	case types.MergePolicyReplace:
		return incoming, nil
	case types.MergePolicyMerge:
		return m.mergeFactCombine(existing, incoming)
	case types.MergePolicySkip:
		return existing, nil // Skip incoming fact
	case types.MergePolicyAppend:
		return m.mergeFactAppend(existing, incoming)
	default:
		return nil, fmt.Errorf("unknown merge policy: %s", m.mergePolicy)
	}
}

// mergeFactCombine combines two facts, preferring the newer one
func (m *FactMerger) mergeFactCombine(existing, incoming *types.Fact) (*types.Fact, error) {
	// Prefer newer fact
	if incoming.Timestamp.After(existing.Timestamp) {
		return incoming, nil
	}

	// If timestamps are equal, prefer the one with more metadata
	if incoming.Timestamp.Equal(existing.Timestamp) && len(incoming.Metadata) > len(existing.Metadata) {
		return incoming, nil
	}

	return existing, nil
}

// mergeFactAppend creates a new fact with an appended key
func (m *FactMerger) mergeFactAppend(_, incoming *types.Fact) (*types.Fact, error) {
	// Create a new fact with appended key
	appendedFact := &types.Fact{
		Key:       fmt.Sprintf("%s_appended", incoming.Key),
		Value:     incoming.Value,
		Source:    incoming.Source,
		Server:    incoming.Server,
		Timestamp: incoming.Timestamp,
		TTL:       incoming.TTL,
		Metadata:  make(map[string]interface{}),
	}

	// Copy metadata and add append information
	for k, v := range incoming.Metadata {
		appendedFact.Metadata[k] = v
	}
	appendedFact.Metadata["original_key"] = incoming.Key
	appendedFact.Metadata["appended"] = true

	return appendedFact, nil
}

// DetectConflicts detects conflicts between existing and incoming facts
func (m *FactMerger) DetectConflicts(existing, incoming *types.FactCollection) []string {
	var conflicts []string

	for key := range incoming.Facts {
		if _, exists := existing.Facts[key]; exists {
			conflicts = append(conflicts, key)
		}
	}

	return conflicts
}

// ValidateMergePolicy validates that a merge policy is supported
func ValidateMergePolicy(policy types.MergePolicy) error {
	validPolicies := []types.MergePolicy{
		types.MergePolicyReplace,
		types.MergePolicyMerge,
		types.MergePolicySkip,
		types.MergePolicyAppend,
	}

	for _, valid := range validPolicies {
		if policy == valid {
			return nil
		}
	}

	return fmt.Errorf("invalid merge policy: %s. Valid policies are: %s",
		policy, strings.Join([]string{
			string(types.MergePolicyReplace),
			string(types.MergePolicyMerge),
			string(types.MergePolicySkip),
			string(types.MergePolicyAppend),
		}, ", "))
}

// DeepMerge performs deep merging of nested structures
func DeepMerge(existing, custom interface{}) interface{} {
	if existing == nil {
		return custom
	}
	if custom == nil {
		return existing
	}

	switch existingVal := existing.(type) {
	case map[string]interface{}:
		if customMap, ok := custom.(map[string]interface{}); ok {
			merged := make(map[string]interface{})

			// Copy existing values
			for k, v := range existingVal {
				merged[k] = v
			}

			// Merge custom values
			for k, v := range customMap {
				if existingVal, exists := existingVal[k]; exists {
					merged[k] = DeepMerge(existingVal, v)
				} else {
					merged[k] = v
				}
			}

			return merged
		}
	case []interface{}:
		if customSlice, ok := custom.([]interface{}); ok {
			// For arrays, append custom values
			return append(existingVal, customSlice...)
		}
	}

	// For primitive types, prefer custom value
	return custom
}

// ApplyOverrides applies overrides to existing facts
func ApplyOverrides(facts *types.FactCollection, overrides map[string]interface{}) *types.FactCollection {
	if overrides == nil {
		return facts
	}

	merged := facts.Clone()

	for category, categoryOverrides := range overrides {
		if categoryMap, ok := categoryOverrides.(map[string]interface{}); ok {
			for key, value := range categoryMap {
				factKey := fmt.Sprintf("%s.%s", category, key)
				merged.Facts[factKey] = &types.Fact{
					Key:       factKey,
					Value:     value,
					Source:    string(types.SourceCustom),
					Server:    facts.Server,
					Timestamp: time.Now(),
					TTL:       types.DefaultTTL,
					Metadata:  map[string]interface{}{"override": true, "category": category},
				}
			}
		}
	}

	return merged
}
