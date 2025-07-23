package facts

import (
	"fmt"
	"time"

	"spooky/internal/logging"
)

// CollectorType represents the type of fact collector
type CollectorType string

const (
	CollectorTypeLocal    CollectorType = "local"
	CollectorTypeSSH      CollectorType = "ssh"
	CollectorTypeJSON     CollectorType = "json"
	CollectorTypeHTTP     CollectorType = "http"
	CollectorTypeHCL      CollectorType = "hcl"
	CollectorTypeOpenTofu CollectorType = "opentofu"
	CollectorTypeStub     CollectorType = "stub"
)

// CollectorConfig holds configuration for any fact collector
type CollectorConfig struct {
	Type        CollectorType
	Source      string                 // File path, URL, etc.
	Headers     map[string]string      // HTTP headers
	Timeout     time.Duration          // Timeout for operations
	MergePolicy MergePolicy            // How to merge facts
	Logger      logging.Logger         // Logger instance
	Metadata    map[string]interface{} // Additional metadata
	Parameters  map[string]interface{} // Parameters for specific collectors
}

// NewCollectorConfig creates a new collector configuration with defaults
func NewCollectorConfig(collectorType CollectorType, source string) *CollectorConfig {
	return &CollectorConfig{
		Type:        collectorType,
		Source:      source,
		Headers:     make(map[string]string),
		Timeout:     30 * time.Second,
		MergePolicy: MergePolicyReplace,
		Logger:      logging.GetLogger(),
		Metadata:    make(map[string]interface{}),
		Parameters:  make(map[string]interface{}),
	}
}

// WithHeaders adds HTTP headers to the configuration
func (c *CollectorConfig) WithHeaders(headers map[string]string) *CollectorConfig {
	c.Headers = headers
	return c
}

// WithTimeout sets the timeout for the collector
func (c *CollectorConfig) WithTimeout(timeout time.Duration) *CollectorConfig {
	c.Timeout = timeout
	return c
}

// WithMergePolicy sets the merge policy for the collector
func (c *CollectorConfig) WithMergePolicy(policy MergePolicy) *CollectorConfig {
	c.MergePolicy = policy
	return c
}

// WithLogger sets the logger for the collector
func (c *CollectorConfig) WithLogger(logger logging.Logger) *CollectorConfig {
	c.Logger = logger
	return c
}

// WithMetadata adds metadata to the collector configuration
func (c *CollectorConfig) WithMetadata(metadata map[string]interface{}) *CollectorConfig {
	for k, v := range metadata {
		c.Metadata[k] = v
	}
	return c
}

// Validate validates the collector configuration
func (c *CollectorConfig) Validate() error {
	if c.Type == "" {
		return ErrInvalidSource("collector type", "cannot be empty")
	}

	if c.Source == "" {
		return ErrInvalidSource("source", "cannot be empty")
	}

	if c.Timeout <= 0 {
		return ErrInvalidSource("timeout", "must be positive")
	}

	if err := ValidateMergePolicy(c.MergePolicy); err != nil {
		return ErrInvalidSource("merge policy", err.Error())
	}

	return nil
}

// NewCollector creates a new fact collector based on the configuration
func NewCollector(config *CollectorConfig) (FactCollector, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	switch config.Type {
	case CollectorTypeLocal:
		return NewLocalCollector(), nil
	case CollectorTypeSSH:
		// SSH collector needs SSH client, which should be passed separately
		return nil, ErrInvalidSource("SSH collector", "requires SSH client - use NewSSHCollector directly")
	case CollectorTypeJSON:
		return NewJSONCollector(config.Source, config.MergePolicy), nil
	case CollectorTypeHTTP:
		return NewHTTPCollector(config.Source, config.Headers, config.Timeout, config.MergePolicy), nil
	case CollectorTypeHCL:
		filePath, ok := config.Parameters["file_path"].(string)
		if !ok {
			return nil, fmt.Errorf("HCL collector requires 'file_path' parameter")
		}
		return NewHCLCollector(filePath, config.Logger, config.MergePolicy), nil
	case CollectorTypeOpenTofu:
		statePath, ok := config.Parameters["state_path"].(string)
		if !ok {
			return nil, fmt.Errorf("OpenTofu collector requires 'state_path' parameter")
		}
		return NewOpenTofuCollector(statePath, config.Logger, config.MergePolicy), nil
	case CollectorTypeStub:
		return NewStubCollector("stub"), nil
	default:
		return nil, ErrInvalidSource(string(config.Type), "unknown collector type")
	}
}
