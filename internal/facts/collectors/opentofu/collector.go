package opentofu

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	spookyfactscollectors "spooky/internal/facts/collectors"
	spookyfactstypes "spooky/internal/facts/types"
)

// Collector collects facts from OpenTofu state files and outputs
type Collector struct {
	spookyfactscollectors.BaseCollector
	stateFile string
	configDir string
	outputs   map[string]interface{}
}

// NewCollector creates a new OpenTofu fact collector
func NewCollector(stateFile, configDir string) *Collector {
	return &Collector{
		BaseCollector: *spookyfactscollectors.NewBaseCollector(spookyfactstypes.SourceOpenTofu, spookyfactstypes.MergePolicyMerge),
		stateFile:     stateFile,
		configDir:     configDir,
		outputs:       make(map[string]interface{}),
	}
}

// Collect gathers facts from OpenTofu state and outputs
func (c *Collector) Collect(server string) (*spookyfactstypes.FactCollection, error) {
	collection := &spookyfactstypes.FactCollection{
		Server:    server,
		Timestamp: time.Now(),
		Facts:     make(map[string]*spookyfactstypes.Fact),
	}

	// Collect from state file if available
	if c.stateFile != "" {
		if err := c.collectFromStateFile(collection); err != nil {
			return nil, fmt.Errorf("failed to collect from state file: %w", err)
		}
	}

	// Collect from outputs if config directory is available
	if c.configDir != "" {
		if err := c.collectFromOutputs(collection); err != nil {
			return nil, fmt.Errorf("failed to collect from outputs: %w", err)
		}
	}

	return collection, nil
}

// CollectSpecific collects only the specified facts from OpenTofu
func (c *Collector) CollectSpecific(server string, keys []string) (*spookyfactstypes.FactCollection, error) {
	collection := &spookyfactstypes.FactCollection{
		Server:    server,
		Timestamp: time.Now(),
		Facts:     make(map[string]*spookyfactstypes.Fact),
	}

	// Collect all facts first, then filter
	allFacts, err := c.Collect(server)
	if err != nil {
		return nil, err
	}

	// Filter to only requested keys
	for _, key := range keys {
		if fact, exists := allFacts.Facts[key]; exists {
			collection.Facts[key] = fact
		}
	}

	return collection, nil
}

// GetFact retrieves a single fact from OpenTofu
func (c *Collector) GetFact(server, key string) (*spookyfactstypes.Fact, error) {
	collection := &spookyfactstypes.FactCollection{
		Server:    server,
		Timestamp: time.Now(),
		Facts:     make(map[string]*spookyfactstypes.Fact),
	}

	// Try to get from state file first
	if c.stateFile != "" {
		if err := c.collectFromStateFile(collection); err == nil {
			if fact, exists := collection.Facts[key]; exists {
				return fact, nil
			}
		}
	}

	// Try to get from outputs
	if c.configDir != "" {
		if err := c.collectFromOutputs(collection); err == nil {
			if fact, exists := collection.Facts[key]; exists {
				return fact, nil
			}
		}
	}

	return nil, fmt.Errorf("fact %s not found in OpenTofu state or outputs", key)
}

// Validate validates the collector configuration
func (c *Collector) Validate() error {
	if c.stateFile == "" && c.configDir == "" {
		return fmt.Errorf("neither state file nor config directory is configured")
	}

	if c.stateFile != "" {
		if _, err := os.Stat(c.stateFile); os.IsNotExist(err) {
			return fmt.Errorf("state file does not exist: %s", c.stateFile)
		}
	}

	if c.configDir != "" {
		if _, err := os.Stat(c.configDir); os.IsNotExist(err) {
			return fmt.Errorf("config directory does not exist: %s", c.configDir)
		}
	}

	return nil
}

// SetStateFile sets the state file path
func (c *Collector) SetStateFile(stateFile string) {
	c.stateFile = stateFile
}

// SetConfigDir sets the config directory path
func (c *Collector) SetConfigDir(configDir string) {
	c.configDir = configDir
}

// GetStateFile returns the current state file path
func (c *Collector) GetStateFile() string {
	return c.stateFile
}

// GetConfigDir returns the current config directory path
func (c *Collector) GetConfigDir() string {
	return c.configDir
}

// collectFromStateFile collects facts from OpenTofu state file
func (c *Collector) collectFromStateFile(collection *spookyfactstypes.FactCollection) error {
	if c.stateFile == "" {
		return fmt.Errorf("state file not configured")
	}

	// Read state file
	content, err := os.ReadFile(c.stateFile)
	if err != nil {
		return fmt.Errorf("failed to read state file: %w", err)
	}

	// Parse state file
	var state StateFile
	if err := json.Unmarshal(content, &state); err != nil {
		return fmt.Errorf("failed to parse state file: %w", err)
	}

	// Extract facts from state
	c.extractFactsFromState(collection, &state)

	return nil
}

// collectFromOutputs collects facts from OpenTofu outputs
func (c *Collector) collectFromOutputs(collection *spookyfactstypes.FactCollection) error {
	if c.configDir == "" {
		return fmt.Errorf("config directory not configured")
	}

	// Run `tofu output -json` to get outputs
	cmd := exec.Command("tofu", "output", "-json")
	cmd.Dir = c.configDir

	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to run tofu output: %w", err)
	}

	// Parse outputs
	var outputs map[string]Output
	if err := json.Unmarshal(output, &outputs); err != nil {
		return fmt.Errorf("failed to parse outputs: %w", err)
	}

	// Extract facts from outputs
	c.extractFactsFromOutputs(collection, outputs)

	return nil
}

// extractFactsFromState extracts facts from OpenTofu state
func (c *Collector) extractFactsFromState(collection *spookyfactstypes.FactCollection, state *StateFile) {
	// Extract version
	c.createFact(collection, "tofu_version", state.Version)

	// Extract terraform version
	c.createFact(collection, "terraform_version", state.TerraformVersion)

	// Extract outputs
	for name, output := range state.Outputs {
		c.createFact(collection, fmt.Sprintf("output_%s", name), output.Value)
	}

	// Extract resource information
	for _, module := range state.Modules {
		for _, resource := range module.Resources {
			// Extract resource type and name
			c.createFact(collection, fmt.Sprintf("resource_%s_type", resource.Name), resource.Type)
			c.createFact(collection, fmt.Sprintf("resource_%s_provider", resource.Name), resource.Provider)

			// Extract resource attributes
			if resource.Primary != nil {
				for key, value := range resource.Primary.Attributes {
					c.createFact(collection, fmt.Sprintf("resource_%s_%s", resource.Name, key), value)
				}
			}
		}
	}
}

// extractFactsFromOutputs extracts facts from OpenTofu outputs
func (c *Collector) extractFactsFromOutputs(collection *spookyfactstypes.FactCollection, outputs map[string]Output) {
	for name, output := range outputs {
		c.createFact(collection, fmt.Sprintf("output_%s", name), output.Value)
		c.createFact(collection, fmt.Sprintf("output_%s_sensitive", name), output.Sensitive)
	}
}

// createFact creates a fact in the collection
func (c *Collector) createFact(collection *spookyfactstypes.FactCollection, key string, value interface{}) {
	collection.Facts[key] = &spookyfactstypes.Fact{
		Key:       key,
		Value:     value,
		Source:    string(c.GetSource()),
		Timestamp: time.Now(),
	}
}

// ExportStateToJSON exports the state file to JSON format
func (c *Collector) ExportStateToJSON(w io.Writer) error {
	if c.stateFile == "" {
		return fmt.Errorf("state file not configured")
	}

	content, err := os.ReadFile(c.stateFile)
	if err != nil {
		return fmt.Errorf("failed to read state file: %w", err)
	}

	_, err = w.Write(content)
	return err
}

// ExportOutputsToJSON exports outputs to JSON format
func (c *Collector) ExportOutputsToJSON(w io.Writer) error {
	if c.configDir == "" {
		return fmt.Errorf("config directory not configured")
	}

	cmd := exec.Command("tofu", "output", "-json")
	cmd.Dir = c.configDir

	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to run tofu output: %w", err)
	}

	_, err = w.Write(output)
	return err
}

// GetStateInfo returns information about the state file
func (c *Collector) GetStateInfo() (*StateInfo, error) {
	if c.stateFile == "" {
		return nil, fmt.Errorf("state file not configured")
	}

	content, err := os.ReadFile(c.stateFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read state file: %w", err)
	}

	var state StateFile
	if err := json.Unmarshal(content, &state); err != nil {
		return nil, fmt.Errorf("failed to parse state file: %w", err)
	}

	info := &StateInfo{
		Version:          state.Version,
		TerraformVersion: state.TerraformVersion,
		Serial:           state.Serial,
		Lineage:          state.Lineage,
		OutputCount:      len(state.Outputs),
		ModuleCount:      len(state.Modules),
		ResourceCount:    0,
		StateFilePath:    c.stateFile,
		LastModified:     time.Now(), // We could get this from file stat
	}

	// Count resources
	for _, module := range state.Modules {
		info.ResourceCount += len(module.Resources)
	}

	return info, nil
}

// GetOutputInfo returns information about outputs
func (c *Collector) GetOutputInfo() (*OutputInfo, error) {
	if c.configDir == "" {
		return nil, fmt.Errorf("config directory not configured")
	}

	cmd := exec.Command("tofu", "output", "-json")
	cmd.Dir = c.configDir

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to run tofu output: %w", err)
	}

	var outputs map[string]Output
	if err := json.Unmarshal(output, &outputs); err != nil {
		return nil, fmt.Errorf("failed to parse outputs: %w", err)
	}

	info := &OutputInfo{
		OutputCount:    len(outputs),
		SensitiveCount: 0,
		ConfigDir:      c.configDir,
		LastUpdated:    time.Now(),
	}

	// Count sensitive outputs
	for _, output := range outputs {
		if output.Sensitive {
			info.SensitiveCount++
		}
	}

	return info, nil
}

// ValidateStateFile validates that a state file is valid
func (c *Collector) ValidateStateFile(filePath string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read state file: %w", err)
	}

	var state StateFile
	if err := json.Unmarshal(content, &state); err != nil {
		return fmt.Errorf("failed to parse state file: %w", err)
	}

	// Basic validation
	if state.Version == 0 {
		return fmt.Errorf("invalid state file: version is 0")
	}

	return nil
}

// ValidateConfig validates that a config directory is valid
func (c *Collector) ValidateConfig(configDir string) error {
	// Check if config directory exists
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		return fmt.Errorf("config directory does not exist: %s", configDir)
	}

	// Check if there are .tf files
	tfFiles, err := filepath.Glob(filepath.Join(configDir, "*.tf"))
	if err != nil {
		return fmt.Errorf("failed to check for .tf files: %w", err)
	}

	if len(tfFiles) == 0 {
		return fmt.Errorf("no .tf files found in config directory")
	}

	return nil
}

// OpenTofu state file structures

// StateFile represents an OpenTofu state file
type StateFile struct {
	Version          int               `json:"version"`
	TerraformVersion string            `json:"terraform_version"`
	Serial           int               `json:"serial"`
	Lineage          string            `json:"lineage"`
	Outputs          map[string]Output `json:"outputs"`
	Modules          []Module          `json:"modules"`
}

// Module represents a module in the state
type Module struct {
	Path      []string            `json:"path"`
	Outputs   map[string]Output   `json:"outputs"`
	Resources map[string]Resource `json:"resources"`
}

// Resource represents a resource in the state
type Resource struct {
	Type      string              `json:"type"`
	Name      string              `json:"name"`
	Provider  string              `json:"provider"`
	Primary   *ResourceInstance   `json:"primary"`
	Instances []*ResourceInstance `json:"instances"`
}

// ResourceInstance represents a resource instance
type ResourceInstance struct {
	ID         string                 `json:"id"`
	Attributes map[string]interface{} `json:"attributes"`
}

// Output represents an output value
type Output struct {
	Value     interface{} `json:"value"`
	Sensitive bool        `json:"sensitive"`
}

// StateInfo contains information about a state file
type StateInfo struct {
	Version          int       `json:"version"`
	TerraformVersion string    `json:"terraform_version"`
	Serial           int       `json:"serial"`
	Lineage          string    `json:"lineage"`
	OutputCount      int       `json:"output_count"`
	ModuleCount      int       `json:"module_count"`
	ResourceCount    int       `json:"resource_count"`
	StateFilePath    string    `json:"state_file_path"`
	LastModified     time.Time `json:"last_modified"`
}

// OutputInfo contains information about outputs
type OutputInfo struct {
	OutputCount    int       `json:"output_count"`
	SensitiveCount int       `json:"sensitive_count"`
	ConfigDir      string    `json:"config_dir"`
	LastUpdated    time.Time `json:"last_updated"`
}
