package facts

import (
	"fmt"
	"strings"

	"spooky/internal/config"
	"spooky/internal/logging"

	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
)

// HCLCollector collects facts from HCL configuration files
type HCLCollector struct {
	filePath    string
	logger      logging.Logger
	mergePolicy MergePolicy
}

// NewHCLCollector creates a new HCL fact collector
func NewHCLCollector(filePath string, logger logging.Logger, mergePolicy MergePolicy) *HCLCollector {
	if logger == nil {
		logger = logging.GetLogger()
	}

	return &HCLCollector{
		filePath:    filePath,
		logger:      logger,
		mergePolicy: mergePolicy,
	}
}

// Collect reads all facts from the HCL configuration file
func (c *HCLCollector) Collect(server string) (*FactCollection, error) {
	return collectFromFile(
		c.filePath,
		server,
		"HCL",
		c.logger,
		c.parseHCLFile,
		func(data interface{}, server string) map[string]*Fact {
			config := data.(*config.Config)
			return c.extractFactsFromConfig(config, server)
		},
	)
}

// CollectSpecific reads specific facts from the HCL configuration file
func (c *HCLCollector) CollectSpecific(server string, keys []string) (*FactCollection, error) {
	return collectSpecificFacts(c, server, keys, c.logger, "HCL file")
}

// GetFact retrieves a single fact from the HCL configuration file
func (c *HCLCollector) GetFact(server, key string) (*Fact, error) {
	return getSpecificFact(c, server, key, c.logger, "HCL file")
}

// parseHCLFile parses the HCL configuration file
func (c *HCLCollector) parseHCLFile() (interface{}, error) {
	parser := hclparse.NewParser()

	// Read the file
	file, diags := parser.ParseHCLFile(c.filePath)
	if diags.HasErrors() {
		return nil, fmt.Errorf("HCL parsing errors: %s", diags.Error())
	}

	// Decode the configuration
	var cfg config.Config
	diags = gohcl.DecodeBody(file.Body, nil, &cfg)
	if diags.HasErrors() {
		return nil, fmt.Errorf("HCL decoding errors: %s", diags.Error())
	}

	return &cfg, nil
}

// extractFactsFromConfig extracts facts from the parsed configuration
func (c *HCLCollector) extractFactsFromConfig(cfg *config.Config, server string) map[string]*Fact {
	facts := make(map[string]*Fact)

	// Extract machine-specific facts
	for i := range cfg.Machines {
		if cfg.Machines[i].Name == server {
			c.extractMachineFacts(&cfg.Machines[i], facts)
			break
		}
	}

	// Extract global configuration facts
	c.extractGlobalFacts(cfg, facts)

	// Extract action facts
	c.extractActionFacts(cfg, facts, server)

	return facts
}

// extractMachineFacts extracts facts from a machine configuration
func (c *HCLCollector) extractMachineFacts(machine *config.Machine, facts map[string]*Fact) {
	// Basic machine facts
	facts["machine.name"] = &Fact{
		Key:    "machine.name",
		Value:  machine.Name,
		Source: string(SourceHCL),
	}

	facts["machine.host"] = &Fact{
		Key:    "machine.host",
		Value:  machine.Host,
		Source: string(SourceHCL),
	}

	facts["machine.port"] = &Fact{
		Key:    "machine.port",
		Value:  machine.Port,
		Source: string(SourceHCL),
	}

	facts["machine.user"] = &Fact{
		Key:    "machine.user",
		Value:  machine.User,
		Source: string(SourceHCL),
	}

	// Authentication facts
	if machine.Password != "" {
		facts["machine.auth_type"] = &Fact{
			Key:    "machine.auth_type",
			Value:  "password",
			Source: string(SourceHCL),
		}
	} else if machine.KeyFile != "" {
		facts["machine.auth_type"] = &Fact{
			Key:    "machine.auth_type",
			Value:  "key_file",
			Source: string(SourceHCL),
		}
		facts["machine.key_file"] = &Fact{
			Key:    "machine.key_file",
			Value:  machine.KeyFile,
			Source: string(SourceHCL),
		}
	}

	// Tag facts
	for key, value := range machine.Tags {
		factKey := fmt.Sprintf("machine.tags.%s", key)
		facts[factKey] = &Fact{
			Key:    factKey,
			Value:  value,
			Source: string(SourceHCL),
		}
	}
}

// extractGlobalFacts extracts global configuration facts
func (c *HCLCollector) extractGlobalFacts(cfg *config.Config, facts map[string]*Fact) {
	facts["config.machine_count"] = &Fact{
		Key:    "config.machine_count",
		Value:  len(cfg.Machines),
		Source: string(SourceHCL),
	}

	facts["config.action_count"] = &Fact{
		Key:    "config.action_count",
		Value:  len(cfg.Actions),
		Source: string(SourceHCL),
	}

	// Extract unique tags across all machines
	tagSet := make(map[string]bool)
	for _, machine := range cfg.Machines {
		for tag := range machine.Tags {
			tagSet[tag] = true
		}
	}

	facts["config.unique_tags"] = &Fact{
		Key:    "config.unique_tags",
		Value:  len(tagSet),
		Source: string(SourceHCL),
	}
}

// extractActionFacts extracts facts from actions that apply to the server
func (c *HCLCollector) extractActionFacts(cfg *config.Config, facts map[string]*Fact, server string) {
	actionCount := 0
	applicableActions := make([]string, 0)

	for i := range cfg.Actions {
		action := &cfg.Actions[i]
		// Check if action applies to this server
		applies := false

		// Check explicit machine list
		for _, machineName := range action.Machines {
			if machineName == server {
				applies = true
				break
			}
		}

		// Check tags (if no explicit machines specified)
		if !applies && len(action.Machines) == 0 && len(action.Tags) > 0 {
			// Find machine with matching tags
			for j := range cfg.Machines {
				machine := &cfg.Machines[j]
				if machine.Name == server {
					for _, actionTag := range action.Tags {
						if _, hasTag := machine.Tags[actionTag]; hasTag {
							applies = true
							break
						}
					}
					break
				}
			}
		}

		if applies {
			actionCount++
			applicableActions = append(applicableActions, action.Name)

			// Add action-specific facts
			factKey := fmt.Sprintf("action.%s.name", action.Name)
			facts[factKey] = &Fact{
				Key:    factKey,
				Value:  action.Name,
				Source: string(SourceHCL),
			}

			if action.Description != "" {
				factKey = fmt.Sprintf("action.%s.description", action.Name)
				facts[factKey] = &Fact{
					Key:    factKey,
					Value:  action.Description,
					Source: string(SourceHCL),
				}
			}

			if action.Command != "" {
				factKey = fmt.Sprintf("action.%s.command", action.Name)
				facts[factKey] = &Fact{
					Key:    factKey,
					Value:  action.Command,
					Source: string(SourceHCL),
				}
			}

			if action.Script != "" {
				factKey = fmt.Sprintf("action.%s.script", action.Name)
				facts[factKey] = &Fact{
					Key:    factKey,
					Value:  action.Script,
					Source: string(SourceHCL),
				}
			}

			if action.Timeout > 0 {
				factKey = fmt.Sprintf("action.%s.timeout", action.Name)
				facts[factKey] = &Fact{
					Key:    factKey,
					Value:  action.Timeout,
					Source: string(SourceHCL),
				}
			}

			factKey = fmt.Sprintf("action.%s.parallel", action.Name)
			facts[factKey] = &Fact{
				Key:    factKey,
				Value:  action.Parallel,
				Source: string(SourceHCL),
			}
		}
	}

	facts["server.applicable_actions"] = &Fact{
		Key:    "server.applicable_actions",
		Value:  actionCount,
		Source: string(SourceHCL),
	}

	if len(applicableActions) > 0 {
		facts["server.action_names"] = &Fact{
			Key:    "server.action_names",
			Value:  strings.Join(applicableActions, ","),
			Source: string(SourceHCL),
		}
	}
}
