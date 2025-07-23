package facts

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"spooky/internal/logging"
)

// OpenTofuState represents the structure of an OpenTofu state file
type OpenTofuState struct {
	Version          int                    `json:"version"`
	TerraformVersion string                 `json:"terraform_version"`
	Serial           int                    `json:"serial"`
	Lineage          string                 `json:"lineage"`
	Outputs          map[string]interface{} `json:"outputs"`
	Resources        []OpenTofuResource     `json:"resources"`
}

// OpenTofuResource represents a resource in the OpenTofu state
type OpenTofuResource struct {
	Module    string                 `json:"module"`
	Mode      string                 `json:"mode"`
	Type      string                 `json:"type"`
	Name      string                 `json:"name"`
	Provider  string                 `json:"provider"`
	Instances []OpenTofuInstance     `json:"instances"`
	DependsOn []string               `json:"depends_on"`
	Config    map[string]interface{} `json:"config"`
}

// OpenTofuInstance represents an instance of a resource
type OpenTofuInstance struct {
	SchemaVersion int                    `json:"schema_version"`
	Attributes    map[string]interface{} `json:"attributes"`
	Private       string                 `json:"private"`
	Dependencies  []string               `json:"dependencies"`
}

// OpenTofuCollector collects facts from OpenTofu state files and outputs
type OpenTofuCollector struct {
	statePath   string
	logger      logging.Logger
	mergePolicy MergePolicy
}

// NewOpenTofuCollector creates a new OpenTofu fact collector
func NewOpenTofuCollector(statePath string, logger logging.Logger, mergePolicy MergePolicy) *OpenTofuCollector {
	if logger == nil {
		logger = logging.GetLogger()
	}

	return &OpenTofuCollector{
		statePath:   statePath,
		logger:      logger,
		mergePolicy: mergePolicy,
	}
}

// Collect reads all facts from the OpenTofu state file
func (c *OpenTofuCollector) Collect(server string) (*FactCollection, error) {
	return collectFromFile(
		c.statePath,
		server,
		"OpenTofu state",
		c.logger,
		c.parseStateFile,
		func(data interface{}, server string) map[string]*Fact {
			state := data.(*OpenTofuState)
			return c.extractFactsFromState(state, server)
		},
	)
}

// CollectSpecific reads specific facts from the OpenTofu state file
func (c *OpenTofuCollector) CollectSpecific(server string, keys []string) (*FactCollection, error) {
	return collectSpecificFacts(c, server, keys, c.logger, "OpenTofu state")
}

// GetFact retrieves a single fact from the OpenTofu state file
func (c *OpenTofuCollector) GetFact(server, key string) (*Fact, error) {
	return getSpecificFact(c, server, key, c.logger, "OpenTofu state")
}

// parseStateFile parses the OpenTofu state file
func (c *OpenTofuCollector) parseStateFile() (interface{}, error) {
	data, err := os.ReadFile(c.statePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read state file: %w", err)
	}

	var state OpenTofuState
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, fmt.Errorf("failed to unmarshal state file: %w", err)
	}

	return &state, nil
}

// extractFactsFromState extracts facts from the parsed OpenTofu state
func (c *OpenTofuCollector) extractFactsFromState(state *OpenTofuState, server string) map[string]*Fact {
	facts := make(map[string]*Fact)

	// Extract state metadata
	c.extractStateMetadata(state, facts)

	// Extract outputs
	c.extractOutputs(state, facts)

	// Extract resource facts
	c.extractResourceFacts(state, facts, server)

	return facts
}

// extractStateMetadata extracts metadata from the state file
func (c *OpenTofuCollector) extractStateMetadata(state *OpenTofuState, facts map[string]*Fact) {
	facts["opentofu.version"] = &Fact{
		Key:    "opentofu.version",
		Value:  state.Version,
		Source: string(SourceOpenTofu),
	}

	facts["opentofu.terraform_version"] = &Fact{
		Key:    "opentofu.terraform_version",
		Value:  state.TerraformVersion,
		Source: string(SourceOpenTofu),
	}

	facts["opentofu.serial"] = &Fact{
		Key:    "opentofu.serial",
		Value:  state.Serial,
		Source: string(SourceOpenTofu),
	}

	facts["opentofu.lineage"] = &Fact{
		Key:    "opentofu.lineage",
		Value:  state.Lineage,
		Source: string(SourceOpenTofu),
	}

	facts["opentofu.resource_count"] = &Fact{
		Key:    "opentofu.resource_count",
		Value:  len(state.Resources),
		Source: string(SourceOpenTofu),
	}
}

// extractOutputs extracts output values from the state
func (c *OpenTofuCollector) extractOutputs(state *OpenTofuState, facts map[string]*Fact) {
	for name, output := range state.Outputs {
		// Handle different output structures
		if outputMap, ok := output.(map[string]interface{}); ok {
			if value, exists := outputMap["value"]; exists {
				factKey := fmt.Sprintf("opentofu.output.%s", name)
				facts[factKey] = &Fact{
					Key:    factKey,
					Value:  value,
					Source: string(SourceOpenTofu),
				}
			}
		} else {
			// Direct value
			factKey := fmt.Sprintf("opentofu.output.%s", name)
			facts[factKey] = &Fact{
				Key:    factKey,
				Value:  output,
				Source: string(SourceOpenTofu),
			}
		}
	}
}

// extractResourceFacts extracts facts from resources that match the server
func (c *OpenTofuCollector) extractResourceFacts(state *OpenTofuState, facts map[string]*Fact, server string) {
	resourceCount := 0
	resourceTypes := make(map[string]int)

	for i := range state.Resources {
		// Check if this resource is relevant to the server
		if c.isResourceRelevantToServer(&state.Resources[i], server) {
			resourceCount++
			resourceTypes[state.Resources[i].Type]++

			// Extract resource-specific facts
			c.extractResourceInstanceFacts(&state.Resources[i], facts)
		}
	}

	facts["opentofu.server.resource_count"] = &Fact{
		Key:    "opentofu.server.resource_count",
		Value:  resourceCount,
		Source: string(SourceOpenTofu),
	}

	// Add resource type counts
	for resourceType, count := range resourceTypes {
		factKey := fmt.Sprintf("opentofu.server.resource_type.%s", resourceType)
		facts[factKey] = &Fact{
			Key:    factKey,
			Value:  count,
			Source: string(SourceOpenTofu),
		}
	}
}

// isResourceRelevantToServer checks if a resource is relevant to the given server
func (c *OpenTofuCollector) isResourceRelevantToServer(resource *OpenTofuResource, server string) bool {
	// Check resource name for server match
	if strings.Contains(strings.ToLower(resource.Name), strings.ToLower(server)) {
		return true
	}

	// Check instances for server-relevant attributes
	for i := range resource.Instances {
		if c.isInstanceRelevantToServer(&resource.Instances[i], server) {
			return true
		}
	}

	return false
}

// isInstanceRelevantToServer checks if an instance is relevant to the given server
func (c *OpenTofuCollector) isInstanceRelevantToServer(instance *OpenTofuInstance, server string) bool {
	// Check common server-related attributes
	serverAttrs := []string{"name", "hostname", "fqdn", "instance_id", "id"}

	for _, attr := range serverAttrs {
		if value, exists := instance.Attributes[attr]; exists {
			if strValue, ok := value.(string); ok {
				if strings.Contains(strings.ToLower(strValue), strings.ToLower(server)) {
					return true
				}
			}
		}
	}

	return false
}

// extractResourceInstanceFacts extracts facts from resource instances
func (c *OpenTofuCollector) extractResourceInstanceFacts(resource *OpenTofuResource, facts map[string]*Fact) {
	for i := range resource.Instances {
		instance := &resource.Instances[i]
		// Extract basic resource facts
		baseKey := fmt.Sprintf("opentofu.resource.%s.%s", resource.Type, resource.Name)

		facts[fmt.Sprintf("%s.type", baseKey)] = &Fact{
			Key:    fmt.Sprintf("%s.type", baseKey),
			Value:  resource.Type,
			Source: string(SourceOpenTofu),
		}

		facts[fmt.Sprintf("%s.provider", baseKey)] = &Fact{
			Key:    fmt.Sprintf("%s.provider", baseKey),
			Value:  resource.Provider,
			Source: string(SourceOpenTofu),
		}

		facts[fmt.Sprintf("%s.mode", baseKey)] = &Fact{
			Key:    fmt.Sprintf("%s.mode", baseKey),
			Value:  resource.Mode,
			Source: string(SourceOpenTofu),
		}

		// Extract instance-specific facts
		instanceKey := fmt.Sprintf("%s.instance.%d", baseKey, i)

		// Extract common attributes
		commonAttrs := []string{"id", "name", "hostname", "fqdn", "ip_address", "private_ip", "public_ip"}
		for _, attr := range commonAttrs {
			if value, exists := instance.Attributes[attr]; exists {
				factKey := fmt.Sprintf("%s.%s", instanceKey, attr)
				facts[factKey] = &Fact{
					Key:    factKey,
					Value:  value,
					Source: string(SourceOpenTofu),
				}
			}
		}

		// Extract tags if present
		if tags, exists := instance.Attributes["tags"]; exists {
			if tagsMap, ok := tags.(map[string]interface{}); ok {
				for tagKey, tagValue := range tagsMap {
					factKey := fmt.Sprintf("%s.tags.%s", instanceKey, tagKey)
					facts[factKey] = &Fact{
						Key:    factKey,
						Value:  tagValue,
						Source: string(SourceOpenTofu),
					}
				}
			}
		}
	}
}
