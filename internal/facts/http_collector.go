package facts

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"spooky/internal/logging"
)

// HTTPCollector collects facts from HTTP endpoints
type HTTPCollector struct {
	client      *http.Client
	baseURL     string
	headers     map[string]string
	mergePolicy MergePolicy
	logger      logging.Logger
}

// NewHTTPCollector creates a new HTTP-based fact collector
func NewHTTPCollector(baseURL string, headers map[string]string, timeout time.Duration, mergePolicy MergePolicy) *HTTPCollector {
	client := &http.Client{
		Timeout: timeout,
	}

	return &HTTPCollector{
		client:      client,
		baseURL:     baseURL,
		headers:     headers,
		mergePolicy: mergePolicy,
		logger:      logging.GetLogger(),
	}
}

// Validate validates the collector configuration
func (c *HTTPCollector) Validate() error {
	if c.baseURL == "" {
		return ErrInvalidSource("HTTP endpoint", "base URL cannot be empty")
	}

	// Basic URL validation
	if !strings.HasPrefix(c.baseURL, "http://") && !strings.HasPrefix(c.baseURL, "https://") {
		return ErrInvalidSource(c.baseURL, "URL must start with http:// or https://")
	}

	if c.client.Timeout <= 0 {
		return ErrInvalidSource("timeout", "timeout must be positive")
	}

	if err := ValidateMergePolicy(c.mergePolicy); err != nil {
		return ErrInvalidSource("merge policy", err.Error())
	}

	return nil
}

// Collect fetches facts from the HTTP endpoint
func (c *HTTPCollector) Collect(server string) (*FactCollection, error) {
	// Validate inputs
	if err := validateServer(server); err != nil {
		return nil, err
	}
	if err := c.Validate(); err != nil {
		return nil, err
	}

	c.logger.Debug("Starting HTTP fact collection",
		logging.Field{Key: "url", Value: c.baseURL},
		logging.Field{Key: "server", Value: server},
		logging.Field{Key: "merge_policy", Value: c.mergePolicy})

	url := c.baseURL
	if server != "local" {
		// If server-specific endpoint is needed, append server name
		url = fmt.Sprintf("%s/%s", c.baseURL, server)
	}

	collection, err := c.fetchFromURL(url, server)
	if err != nil {
		c.logger.Error("Failed to collect facts from HTTP endpoint", err,
			logging.Field{Key: "url", Value: url},
			logging.Field{Key: "server", Value: server})
		return nil, err
	}

	c.logger.Info("Successfully collected facts from HTTP endpoint",
		logging.Field{Key: "url", Value: url},
		logging.Field{Key: "server", Value: server},
		logging.Field{Key: "fact_count", Value: len(collection.Facts)})

	return collection, nil
}

// CollectSpecific fetches specific facts from the HTTP endpoint
func (c *HTTPCollector) CollectSpecific(server string, keys []string) (*FactCollection, error) {
	return collectSpecificFacts(c, server, keys, c.logger, "HTTP endpoint")
}

// GetFact retrieves a single fact from the HTTP endpoint
func (c *HTTPCollector) GetFact(server, key string) (*Fact, error) {
	return getSpecificFact(c, server, key, c.logger, "HTTP response")
}

// fetchFromURL makes an HTTP request and parses the response
func (c *HTTPCollector) fetchFromURL(url, server string) (*FactCollection, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.client.Timeout)
	defer cancel()

	c.logger.Debug("Making HTTP request",
		logging.Field{Key: "url", Value: url},
		logging.Field{Key: "timeout", Value: c.client.Timeout})

	req, err := http.NewRequestWithContext(ctx, "GET", url, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Add headers
	for key, value := range c.headers {
		req.Header.Set(key, value)
	}

	// Add default headers if not present
	if req.Header.Get("Accept") == "" {
		req.Header.Set("Accept", "application/json")
	}
	if req.Header.Get("User-Agent") == "" {
		req.Header.Set("User-Agent", "spooky-facts-collector/1.0")
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.logger.Error("HTTP request failed", fmt.Errorf("status %d: %s", resp.StatusCode, resp.Status),
			logging.Field{Key: "url", Value: url},
			logging.Field{Key: "status_code", Value: resp.StatusCode})
		return nil, fmt.Errorf("HTTP request failed with status %d: %s", resp.StatusCode, resp.Status)
	}

	c.logger.Debug("HTTP request successful",
		logging.Field{Key: "url", Value: url},
		logging.Field{Key: "status_code", Value: resp.StatusCode})

	return c.parseHTTPResponse(resp.Body, server)
}

// parseHTTPResponse parses the HTTP response body
func (c *HTTPCollector) parseHTTPResponse(body io.Reader, server string) (*FactCollection, error) {
	sourceInfo := buildStandardMetadata("http", c.baseURL, "json")
	return parseJSONFromReader(body, server, sourceInfo)
}
