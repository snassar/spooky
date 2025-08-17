// Package templates provides template management functionality for the spooky codebase.
package templates

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	texttemplate "text/template"
	"time"

	spookyinterfaces "spooky/internal/interfaces"
	spookyschemas "spooky/internal/schemas"
	spookytypes "spooky/internal/types"
	spookytypescommon "spooky/internal/types/common"
	spookytypeslogging "spooky/internal/types/logging"
	spookytypesschemas "spooky/internal/types/schemas"
	spookytypestemplates "spooky/internal/types/templates"
)

// Manager provides enhanced template management functionality
type Manager struct {
	logger             spookytypeslogging.Logger
	cache              TemplateCache
	functions          TemplateFunctionRegistry
	contextResolver    TemplateContextResolver
	metadataManager    TemplateMetadataManager
	validator          TemplateValidator
	securityManager    TemplateSecurityManager
	performanceManager TemplatePerformanceManager
	mu                 sync.RWMutex
}

// TemplateCache provides template caching functionality
type TemplateCache struct {
	cache map[string]CacheEntry
	mu    sync.RWMutex
	ttl   time.Duration
}

// CacheEntry represents a cached template
type CacheEntry struct {
	Template    *spookytypes.Template
	ExpiresAt   time.Time
	AccessCount int
}

// TemplateFunctionRegistry provides secure template function management
type TemplateFunctionRegistry struct {
	functions map[string]TemplateFunction
	security  FunctionSecurityManager
	cache     FunctionResultCache
	mu        sync.RWMutex
}

// TemplateFunction represents a template function with security restrictions
type TemplateFunction struct {
	Name        string
	Function    interface{}
	Security    FunctionSecurity
	Description string
	Examples    []FunctionExample
}

// FunctionSecurity defines security restrictions for template functions
type FunctionSecurity struct {
	AllowedInRestricted bool
	MaxRunningTime      time.Duration
	MaxMemoryUsage      int64
	AllowedPatterns     []string
	ForbiddenPatterns   []string
}

// FunctionExample represents an example of function usage
type FunctionExample struct {
	Input       string
	Output      string
	Description string
}

// FunctionResultCache caches function results
type FunctionResultCache struct {
	cache map[string]FunctionCacheEntry
	ttl   time.Duration
}

// FunctionCacheEntry represents a cached function result
type FunctionCacheEntry struct {
	Result    interface{}
	ExpiresAt time.Time
}

// TemplateContextResolver resolves template context with facts, variables, and machines data
type TemplateContextResolver struct {
	factsManager     spookyinterfaces.FactsIntegration
	variablesManager spookyinterfaces.VariablesIntegration
	machinesManager  spookyinterfaces.MachinesIntegration
	cache            ContextCache
	validator        ContextValidator
	mu               sync.RWMutex
}

// ContextCache caches resolved contexts
type ContextCache struct {
	cache map[string]ContextCacheEntry
	mu    sync.RWMutex
	ttl   time.Duration
}

// ContextCacheEntry represents a cached context
type ContextCacheEntry struct {
	Context   *spookytypes.TemplateContext
	ExpiresAt time.Time
}

// ContextValidator validates template contexts
type ContextValidator struct {
	logger spookytypeslogging.Logger
}

// ResolveContext resolves template context with facts, variables, and machines data
func (c *TemplateContextResolver) ResolveContext(ctx context.Context, template *spookytypes.Template, data map[string]interface{}) (*spookytypes.TemplateContext, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Check cache first
	cacheKey := c.generateCacheKey(template, data)
	if cached, exists := c.getCachedContext(cacheKey); exists {
		return cached, nil
	}

	// Resolve project information
	project, err := c.resolveProjectContext(ctx, template)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve project context: %w", err)
	}

	// Resolve facts
	facts, err := c.resolveFactsContext(ctx, template)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve facts context: %w", err)
	}

	// Resolve machines
	machines, err := c.resolveMachinesContext(ctx, template)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve machines context: %w", err)
	}

	// Resolve environment variables
	environment, err := c.resolveEnvironmentContext(ctx, template)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve environment context: %w", err)
	}

	// Resolve variables
	variables, err := c.resolveVariablesContext(ctx, template)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve variables context: %w", err)
	}

	// Create context
	context := &spookytypes.TemplateContext{
		Project:     project,
		Facts:       facts,
		Machines:    machines,
		Environment: environment,
		CustomData:  data,
		Variables:   variables,
	}

	// Validate context
	if err := c.ValidateContext(ctx, context); err != nil {
		return nil, fmt.Errorf("context validation failed: %w", err)
	}

	// Cache the context
	c.cacheContext(cacheKey, context)

	return context, nil
}

// ValidateContext validates a template context
func (c *TemplateContextResolver) ValidateContext(_ context.Context, context *spookytypes.TemplateContext) error {
	// Validate project context
	if context.Project == nil {
		return fmt.Errorf("project context is required")
	}

	// Validate facts context
	if context.Facts == nil {
		// Facts are optional, but if present should be valid
		context.Facts = make(map[string]interface{})
	}

	// Validate machines context
	if context.Machines == nil {
		// Machines are optional, but if present should be valid
		context.Machines = make([]map[string]interface{}, 0)
	}

	// Validate environment context
	if context.Environment == nil {
		// Environment is optional, but if present should be valid
		context.Environment = make(map[string]string)
	}

	// Validate variables context
	if context.Variables == nil {
		// Variables are optional, but if present should be valid
		context.Variables = make(map[string]interface{})
	}

	return nil
}

// TransformContext transforms a template context
func (c *TemplateContextResolver) TransformContext(_ context.Context, context *spookytypes.TemplateContext) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Transform project context
	if err := c.transformProjectContext(context.Project); err != nil {
		return fmt.Errorf("failed to transform project context: %w", err)
	}

	// Transform facts context
	if err := c.transformFactsContext(context.Facts); err != nil {
		return fmt.Errorf("failed to transform facts context: %w", err)
	}

	// Transform machines context
	if err := c.transformMachinesContext(context.Machines); err != nil {
		return fmt.Errorf("failed to transform machines context: %w", err)
	}

	// Transform environment context
	if err := c.transformEnvironmentContext(context.Environment); err != nil {
		return fmt.Errorf("failed to transform environment context: %w", err)
	}

	// Transform variables context
	if err := c.transformVariablesContext(context.Variables); err != nil {
		return fmt.Errorf("failed to transform variables context: %w", err)
	}

	return nil
}

// TemplateMetadataManager manages template metadata
type TemplateMetadataManager struct {
	validator MetadataValidator
	indexer   *EnhancedMetadataIndexer
	cache     MetadataCache
	mu        sync.RWMutex
}

// MetadataValidator validates template metadata
type MetadataValidator struct {
	logger spookytypeslogging.Logger
}

// LoadMetadata loads template metadata
func (m *TemplateMetadataManager) LoadMetadata(ctx context.Context, templatePath string) (*spookytypestemplates.TemplateMetadata, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Check cache first
	if cached, exists := m.getCachedMetadata(templatePath); exists {
		return cached, nil
	}

	// Load metadata from file
	metadata, err := m.loadMetadataFromFile(ctx, templatePath)
	if err != nil {
		return nil, fmt.Errorf("failed to load metadata from file: %w", err)
	}

	// Validate metadata
	if err := m.ValidateMetadata(ctx, metadata); err != nil {
		return nil, fmt.Errorf("metadata validation failed: %w", err)
	}

	// Index metadata
	if err := m.IndexMetadata(ctx, metadata); err != nil {
		return nil, fmt.Errorf("failed to index metadata: %w", err)
	}

	// Cache metadata
	m.cacheMetadata(templatePath, metadata)

	return metadata, nil
}

// ValidateMetadata validates template metadata
func (m *TemplateMetadataManager) ValidateMetadata(_ context.Context, metadata *spookytypestemplates.TemplateMetadata) error {
	// Basic validation
	if metadata == nil {
		return fmt.Errorf("metadata cannot be nil")
	}

	// Validate name
	if metadata.Name == "" {
		return fmt.Errorf("metadata name is required")
	}

	// Validate version format
	if metadata.Version != "" {
		if !isValidVersion(metadata.Version) {
			return fmt.Errorf("invalid version format: %s", metadata.Version)
		}
	}

	// Validate tags
	if metadata.Tags != nil {
		for i, tag := range metadata.Tags {
			if tag == "" {
				return fmt.Errorf("tag at index %d cannot be empty", i)
			}
		}
	}

	return nil
}

// IndexMetadata indexes template metadata
func (m *TemplateMetadataManager) IndexMetadata(_ context.Context, metadata *spookytypestemplates.TemplateMetadata) error {
	return m.indexer.IndexMetadata(metadata)
}

// SearchMetadata searches for template metadata
func (m *TemplateMetadataManager) SearchMetadata(_ context.Context, query string) ([]*spookytypestemplates.TemplateMetadata, error) {
	results, err := m.indexer.Search(query, &SearchFilters{})
	if err != nil {
		return nil, err
	}

	var metadata []*spookytypestemplates.TemplateMetadata
	for _, result := range results {
		metadata = append(metadata, result.Metadata)
	}

	return metadata, nil
}

// loadMetadataFromFile loads metadata from a file
func (m *TemplateMetadataManager) loadMetadataFromFile(_ context.Context, templatePath string) (*spookytypestemplates.TemplateMetadata, error) {
	// Check if metadata file exists
	metadataPath := templatePath + ".meta"
	if _, err := os.Stat(metadataPath); os.IsNotExist(err) {
		// Create default metadata if file doesn't exist
		return m.createDefaultMetadata(templatePath), nil
	}

	// Load metadata from file
	// For now, create basic metadata
	// This will be enhanced to load from HCL/JSON files
	metadata := &spookytypestemplates.TemplateMetadata{
		Name:        filepath.Base(templatePath),
		Description: fmt.Sprintf("Template loaded from %s", templatePath),
		Version:     "1.0.0",
		Tags:        []string{},
	}

	return metadata, nil
}

// createDefaultMetadata creates default metadata for a template
func (m *TemplateMetadataManager) createDefaultMetadata(templatePath string) *spookytypestemplates.TemplateMetadata {
	return &spookytypestemplates.TemplateMetadata{
		Name:        filepath.Base(templatePath),
		Description: fmt.Sprintf("Template loaded from %s", templatePath),
		Version:     "1.0.0",
		Tags:        []string{},
	}
}

// getCachedMetadata gets cached metadata
func (m *TemplateMetadataManager) getCachedMetadata(key string) (*spookytypestemplates.TemplateMetadata, bool) {
	m.cache.mu.RLock()
	defer m.cache.mu.RUnlock()

	entry, exists := m.cache.cache[key]
	if !exists {
		return nil, false
	}

	// Check if cache entry has expired
	if time.Now().After(entry.ExpiresAt) {
		delete(m.cache.cache, key)
		return nil, false
	}

	return entry.Metadata, true
}

// cacheMetadata caches metadata
func (m *TemplateMetadataManager) cacheMetadata(key string, metadata *spookytypestemplates.TemplateMetadata) {
	m.cache.mu.Lock()
	defer m.cache.mu.Unlock()

	m.cache.cache[key] = MetadataCacheEntry{
		Metadata:  metadata,
		ExpiresAt: time.Now().Add(m.cache.ttl),
	}
}

// MetadataCache caches template metadata
type MetadataCache struct {
	cache map[string]MetadataCacheEntry
	mu    sync.RWMutex
	ttl   time.Duration
}

// MetadataCacheEntry represents cached metadata
type MetadataCacheEntry struct {
	Metadata  *spookytypestemplates.TemplateMetadata
	ExpiresAt time.Time
}

// TemplateValidator provides comprehensive template validation
type TemplateValidator struct {
	schemaValidator   spookytypesschemas.SchemaValidator
	schemaManager     *spookyschemas.Manager
	functionValidator FunctionValidator
	contextValidator  ContextValidator
	metadataValidator MetadataValidator
	cache             ValidationCache
	logger            spookytypeslogging.Logger
}

// ValidateTemplateComprehensive validates template against schemas
func (v *TemplateValidator) ValidateTemplateComprehensive(ctx context.Context, template *spookytypes.Template) (*spookytypesschemas.ValidationResult, error) {
	if v.schemaValidator == nil {
		return &spookytypesschemas.ValidationResult{
			Valid:       false,
			ValidatedAt: time.Now(),
			Errors: []spookytypesschemas.SchemaError{
				{
					Message:   "schema validator not configured",
					FieldPath: "schema_validator",
					Severity:  "error",
				},
			},
			Warnings: []spookytypesschemas.SchemaError{},
			Info:     []spookytypesschemas.SchemaError{},
			Details:  make(map[string]interface{}),
		}, nil
	}

	// Create a comprehensive validation result
	result := &spookytypesschemas.ValidationResult{
		Valid:       true,
		ValidatedAt: time.Now(),
		Errors:      []spookytypesschemas.SchemaError{},
		Warnings:    []spookytypesschemas.SchemaError{},
		Info:        []spookytypesschemas.SchemaError{},
		Details:     make(map[string]interface{}),
	}

	// Validate template structure against template-structure schema
	if err := v.validateTemplateStructure(ctx, template, result); err != nil {
		v.logger.Error("Template structure validation failed", err, map[string]interface{}{
			"template_id": template.ID,
		})
		result.Valid = false
	}

	// Validate template context against template-context schema
	if err := v.validateTemplateContext(ctx, template, result); err != nil {
		v.logger.Error("Template context validation failed", err, map[string]interface{}{
			"template_id": template.ID,
		})
		result.Valid = false
	}

	// Validate template functions against template-functions schema
	if err := v.validateTemplateFunctions(ctx, template, result); err != nil {
		v.logger.Error("Template functions validation failed", err, map[string]interface{}{
			"template_id": template.ID,
		})
		result.Valid = false
	}

	// Validate template metadata against template-metadata schema
	if err := v.validateTemplateMetadata(ctx, template, result); err != nil {
		v.logger.Error("Template metadata validation failed", err, map[string]interface{}{
			"template_id": template.ID,
		})
		result.Valid = false
	}

	return result, nil
}

// validateTemplateStructure validates template against template-structure schema
func (v *TemplateValidator) validateTemplateStructure(_ context.Context, template *spookytypes.Template, result *spookytypesschemas.ValidationResult) error {
	// Load template-metadata-validation schema
	schema, err := v.loadTemplateSchema("template-metadata-validation")
	if err != nil {
		return fmt.Errorf("failed to load template-metadata-validation schema: %w", err)
	}

	// Convert template metadata to map for validation
	metadataData := v.templateMetadataToMap(template)

	// Validate against schema
	validationResult, err := v.schemaValidator.Validate(schema, metadataData)
	if err != nil {
		return fmt.Errorf("template metadata validation failed: %w", err)
	}

	// Merge validation results
	v.mergeValidationResults(result, validationResult)

	return nil
}

// loadTemplateSchema loads a template schema by name
func (v *TemplateValidator) loadTemplateSchema(schemaName string) (*spookytypesschemas.Schema, error) {
	schemaPath := fmt.Sprintf("internal/schemas/schemas/%s.schema.hcl", schemaName)

	// Use schema manager to properly load and parse the schema
	if v.schemaManager != nil {
		return v.schemaManager.Load(schemaPath)
	}

	// Fallback to manual loading if no schema manager available
	data, err := os.ReadFile(schemaPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read schema file %s: %w", schemaPath, err)
	}

	// Create schema object
	schema := &spookytypesschemas.Schema{
		Name:        schemaName,
		Description: fmt.Sprintf("Template %s schema", schemaName),
		Content:     string(data),
		Type:        "hcl",
		Version:     "1.0",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Metadata:    make(map[string]interface{}),
	}

	return schema, nil
}

// templateContextToMap converts template context to a map for validation
func (v *TemplateValidator) templateContextToMap(template *spookytypes.Template) map[string]interface{} {
	if template.ContextData == nil {
		return make(map[string]interface{})
	}

	return template.ContextData
}

// mergeValidationResults merges validation results
func (v *TemplateValidator) mergeValidationResults(target, source *spookytypesschemas.ValidationResult) {
	// Merge errors
	target.Errors = append(target.Errors, source.Errors...)

	// Merge warnings
	target.Warnings = append(target.Warnings, source.Warnings...)

	// Merge info
	target.Info = append(target.Info, source.Info...)

	// Update validity
	if !source.Valid {
		target.Valid = false
	}

	// Merge details
	for key, value := range source.Details {
		target.Details[key] = value
	}
}

// validateTemplateContext validates template against template-context schema
func (v *TemplateValidator) validateTemplateContext(_ context.Context, template *spookytypes.Template, result *spookytypesschemas.ValidationResult) error {
	// Load template-context schema
	schema, err := v.loadTemplateSchema("template-context")
	if err != nil {
		return fmt.Errorf("failed to load template-context schema: %w", err)
	}

	// Convert template context to map for validation
	contextData := v.templateContextToMap(template)

	// Validate against schema
	validationResult, err := v.schemaValidator.Validate(schema, contextData)
	if err != nil {
		return fmt.Errorf("template context validation failed: %w", err)
	}

	// Merge validation results
	v.mergeValidationResults(result, validationResult)

	return nil
}

// validateTemplateFunctions validates template against template-functions schema
func (v *TemplateValidator) validateTemplateFunctions(_ context.Context, template *spookytypes.Template, result *spookytypesschemas.ValidationResult) error {
	// Load template-functions schema
	schema, err := v.loadTemplateSchema("template-functions")
	if err != nil {
		return fmt.Errorf("failed to load template-functions schema: %w", err)
	}

	// Convert template functions to map for validation
	functionsData := v.templateFunctionsToMap(template)

	// Validate against schema
	validationResult, err := v.schemaValidator.Validate(schema, functionsData)
	if err != nil {
		return fmt.Errorf("template functions validation failed: %w", err)
	}

	// Merge validation results
	v.mergeValidationResults(result, validationResult)

	return nil
}

// validateTemplateMetadata validates template against template-metadata schema
func (v *TemplateValidator) validateTemplateMetadata(_ context.Context, template *spookytypes.Template, result *spookytypesschemas.ValidationResult) error {
	// Load template-metadata schema
	schema, err := v.loadTemplateSchema("template-metadata")
	if err != nil {
		return fmt.Errorf("failed to load template-metadata schema: %w", err)
	}

	// Convert template metadata to map for validation
	metadataData := v.templateMetadataToMap(template)

	// Validate against schema
	validationResult, err := v.schemaValidator.Validate(schema, metadataData)
	if err != nil {
		return fmt.Errorf("template metadata validation failed: %w", err)
	}

	// Merge validation results
	v.mergeValidationResults(result, validationResult)

	return nil
}

// templateFunctionsToMap converts template functions to map for validation
func (v *TemplateValidator) templateFunctionsToMap(template *spookytypes.Template) map[string]interface{} {
	return map[string]interface{}{
		"template_functions": map[string]interface{}{
			"include":       "template-structure",
			"scope":         "functions",
			"template_type": "functions",
			"functions":     template.Functions,
		},
	}
}

// templateMetadataToMap converts template metadata to map for validation
func (v *TemplateValidator) templateMetadataToMap(template *spookytypes.Template) map[string]interface{} {
	metadata := map[string]interface{}{
		"template_metadata": map[string]interface{}{
			"include":       "template-structure",
			"scope":         "metadata",
			"template_type": "metadata",
		},
	}

	if template.Metadata != nil {
		metadata["template_metadata"].(map[string]interface{})["metadata"] = map[string]interface{}{
			"name":        template.Metadata.Name,
			"description": template.Metadata.Description,
			"author":      template.Metadata.Author,
			"version":     template.Metadata.Version,
			"tags":        template.Metadata.Tags,
			"license":     template.Metadata.License,
		}
	}

	return metadata
}

// FunctionValidator validates template functions
type FunctionValidator struct {
	logger spookytypeslogging.Logger
}

// ValidationCache caches validation results
type ValidationCache struct {
	cache map[string]ValidationCacheEntry
	ttl   time.Duration
}

// ValidationCacheEntry represents cached validation results
type ValidationCacheEntry struct {
	Result    *spookytypesschemas.ValidationResult
	ExpiresAt time.Time
}

// TemplateSecurityManager provides template security and sandboxing
type TemplateSecurityManager struct {
	sandbox         TemplateSandbox
	accessControl   AccessController
	patternFilter   PatternFilter
	resourceMonitor ResourceMonitor
	auditLogger     AuditLogger
	mu              sync.RWMutex
}

// TemplateSandbox provides template running sandboxing
type TemplateSandbox struct {
	enabled          bool
	maxRunningTime   time.Duration
	maxMemoryUsage   int64
	allowedFunctions map[string]bool
}

// AccessController manages template access control
type AccessController struct {
	allowedUsers map[string]bool
	allowedPaths map[string]bool
	mu           sync.RWMutex
}

// PatternFilter filters dangerous patterns in templates
type PatternFilter struct {
	forbiddenPatterns []string
	allowedPatterns   []string
	mu                sync.RWMutex
}

// ResourceMonitor monitors template resource usage
type ResourceMonitor struct {
	maxMemoryUsage int64
	maxRunningTime time.Duration
}

// AuditLogger logs template operations for security
type AuditLogger struct {
	logger spookytypeslogging.Logger
	mu     sync.RWMutex
}

// ValidateTemplateSecurity validates template security
func (s *TemplateSecurityManager) ValidateTemplateSecurity(ctx context.Context, template *spookytypes.Template) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check for dangerous patterns
	if err := s.FilterDangerousPatterns(ctx, template); err != nil {
		return fmt.Errorf("dangerous patterns detected: %w", err)
	}

	// Check access control
	if err := s.CheckAccessControl(ctx, template, "default"); err != nil {
		return fmt.Errorf("access control check failed: %w", err)
	}

	// Validate sandbox settings
	if err := s.validateSandboxSettings(template); err != nil {
		return fmt.Errorf("sandbox validation failed: %w", err)
	}

	return nil
}

// SandboxTemplate creates a sandbox for template running
func (s *TemplateSecurityManager) SandboxTemplate(_ context.Context, template *spookytypes.Template) (*TemplateSandbox, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Create sandbox with template-specific settings
	sandbox := &TemplateSandbox{
		allowedFunctions: make(map[string]bool),
		enabled:          s.sandbox.enabled,
		maxRunningTime:   s.sandbox.maxRunningTime,
		maxMemoryUsage:   s.sandbox.maxMemoryUsage,
	}

	// Configure sandbox based on template security level
	switch template.SecurityLevel {
	case "restricted":
		sandbox.enabled = true
		sandbox.maxRunningTime = 10 * time.Second
		sandbox.maxMemoryUsage = 50 * 1024 * 1024 // 50MB
	case "unrestricted":
		sandbox.enabled = false
		sandbox.maxRunningTime = 60 * time.Second
		sandbox.maxMemoryUsage = 500 * 1024 * 1024 // 500MB
	default:
		// Use default settings
	}

	// Log sandbox creation
	s.auditLogger.LogSecurityEvent("sandbox_created", map[string]interface{}{
		"template_id":      template.ID,
		"security_level":   template.SecurityLevel,
		"enabled":          sandbox.enabled,
		"max_running_time": sandbox.maxRunningTime,
		"max_memory_usage": sandbox.maxMemoryUsage,
	})

	return sandbox, nil
}

// CheckAccessControl checks access control for a template
func (s *TemplateSecurityManager) CheckAccessControl(_ context.Context, template *spookytypes.Template, user string) error {
	s.accessControl.mu.RLock()
	defer s.accessControl.mu.RUnlock()

	// Check if user is allowed
	if !s.accessControl.allowedUsers[user] {
		// Log access denied
		s.auditLogger.LogSecurityEvent("access_denied", map[string]interface{}{
			"template_id": template.ID,
			"user":        user,
		})
		return fmt.Errorf("access denied for user: %s", user)
	}

	// Log access granted
	s.auditLogger.LogSecurityEvent("access_granted", map[string]interface{}{
		"template_id": template.ID,
		"user":        user,
	})

	return nil
}

// FilterDangerousPatterns filters dangerous patterns in templates
func (s *TemplateSecurityManager) FilterDangerousPatterns(_ context.Context, template *spookytypes.Template) error {
	s.patternFilter.mu.RLock()
	defer s.patternFilter.mu.RUnlock()

	// Check template content for forbidden patterns
	for _, pattern := range s.patternFilter.forbiddenPatterns {
		if strings.Contains(strings.ToLower(template.Content), strings.ToLower(pattern)) {
			// Log security violation
			s.auditLogger.LogSecurityEvent("dangerous_pattern_detected", map[string]interface{}{
				"template_id": template.ID,
				"pattern":     pattern,
			})
			return fmt.Errorf("dangerous pattern detected: %s", pattern)
		}
	}

	// Check template variables for forbidden patterns
	if template.Variables != nil {
		for key, value := range template.Variables {
			if strValue, ok := value.(string); ok {
				for _, pattern := range s.patternFilter.forbiddenPatterns {
					if strings.Contains(strings.ToLower(strValue), strings.ToLower(pattern)) {
						// Log security violation
						s.auditLogger.LogSecurityEvent("dangerous_pattern_in_variable", map[string]interface{}{
							"template_id":  template.ID,
							"variable_key": key,
							"pattern":      pattern,
						})
						return fmt.Errorf("dangerous pattern in variable %s: %s", key, pattern)
					}
				}
			}
		}
	}

	return nil
}

// validateSandboxSettings validates sandbox settings
func (s *TemplateSecurityManager) validateSandboxSettings(template *spookytypes.Template) error {
	// Validate security level
	validLevels := []string{"restricted", "standard", "unrestricted"}
	valid := false
	for _, level := range validLevels {
		if template.SecurityLevel == level {
			valid = true
			break
		}
	}
	if !valid {
		return fmt.Errorf("invalid security level: %s", template.SecurityLevel)
	}

	// Validate template type
	validTypes := []string{"script", "config", "documentation", "deployment"}
	valid = false
	for _, t := range validTypes {
		if template.Type == t {
			valid = true
			break
		}
	}
	if !valid {
		return fmt.Errorf("invalid template type: %s", template.Type)
	}

	return nil
}

// LogSecurityEvent logs a security event
func (a *AuditLogger) LogSecurityEvent(eventType string, details map[string]interface{}) {
	a.mu.Lock()
	defer a.mu.Unlock()

	// Add timestamp
	details["timestamp"] = time.Now().UTC()
	details["event_type"] = eventType

	// Log the security event
	a.logger.Info("Security event", details)
}

// TemplatePerformanceManager provides performance optimization
type TemplatePerformanceManager struct {
	cache        TemplateResultCache
	compiler     TemplateCompiler
	parallelizer ParallelProcessor
	monitor      PerformanceMonitor
}

// TemplateResultCache caches template rendering results
type TemplateResultCache struct {
	cache map[string]ResultCacheEntry
	mu    sync.RWMutex
	ttl   time.Duration
}

// ResultCacheEntry represents a cached rendering result
type ResultCacheEntry struct {
	Result    string
	ExpiresAt time.Time
}

// TemplateCompiler optimizes template compilation
type TemplateCompiler struct {
	optimizationLevel int
	cache             map[string]*texttemplate.Template
	mu                sync.RWMutex
}

// ParallelProcessor handles parallel template processing
type ParallelProcessor struct {
	maxWorkers int
	semaphore  chan struct{}
}

// PerformanceMonitor monitors template performance
type PerformanceMonitor struct {
	metrics map[string]Metric
	mu      sync.RWMutex
}

// Metric represents a performance metric
type Metric struct {
	Name       string
	Value      float64
	Count      int64
	Min        float64
	Max        float64
	LastUpdate time.Time
}

// NewManager creates a new enhanced template manager
func NewManager(
	logger spookytypeslogging.Logger,
) *Manager {
	manager := &Manager{
		logger: logger,
		cache: TemplateCache{
			cache: make(map[string]CacheEntry),
			ttl:   5 * time.Minute,
		},
		functions: TemplateFunctionRegistry{
			functions: make(map[string]TemplateFunction),
			security: FunctionSecurityManager{
				restrictedMode:   false,
				allowedFunctions: make(map[string]bool),
			},
			cache: FunctionResultCache{
				cache: make(map[string]FunctionCacheEntry),
				ttl:   1 * time.Minute,
			},
		},
		contextResolver: TemplateContextResolver{
			cache: ContextCache{
				cache: make(map[string]ContextCacheEntry),
				ttl:   2 * time.Minute,
			},
			validator: ContextValidator{logger: logger},
		},
		metadataManager: TemplateMetadataManager{
			validator: MetadataValidator{logger: logger},
			indexer:   NewEnhancedMetadataIndexer(),
			cache: MetadataCache{
				cache: make(map[string]MetadataCacheEntry),
				ttl:   10 * time.Minute,
			},
		},
		validator: TemplateValidator{
			functionValidator: FunctionValidator{logger: logger},
			contextValidator:  ContextValidator{logger: logger},
			metadataValidator: MetadataValidator{logger: logger},
			cache: ValidationCache{
				cache: make(map[string]ValidationCacheEntry),
				ttl:   5 * time.Minute,
			},
			logger: logger,
		},
		securityManager: TemplateSecurityManager{
			sandbox: TemplateSandbox{
				enabled:          true,
				maxRunningTime:   30 * time.Second,
				maxMemoryUsage:   100 * 1024 * 1024, // 100MB
				allowedFunctions: make(map[string]bool),
			},
			accessControl: AccessController{
				allowedUsers: make(map[string]bool),
				allowedPaths: make(map[string]bool),
			},
			patternFilter: PatternFilter{
				forbiddenPatterns: []string{
					"exec", "system", "eval", "shell",
					"password", "secret", "key", "token",
				},
				allowedPatterns: []string{},
			},
			resourceMonitor: ResourceMonitor{
				maxMemoryUsage: 100 * 1024 * 1024, // 100MB
				maxRunningTime: 30 * time.Second,
			},
			auditLogger: AuditLogger{logger: logger},
		},
		performanceManager: TemplatePerformanceManager{
			cache: TemplateResultCache{
				cache: make(map[string]ResultCacheEntry),
				ttl:   5 * time.Minute,
			},
			compiler: TemplateCompiler{
				optimizationLevel: 1,
				cache:             make(map[string]*texttemplate.Template),
			},
			parallelizer: ParallelProcessor{
				maxWorkers: 4,
				semaphore:  make(chan struct{}, 4),
			},
			monitor: PerformanceMonitor{
				metrics: make(map[string]Metric),
			},
		},
	}

	// Register built-in template functions
	manager.registerBuiltInFunctions()

	return manager
}

// LoadTemplate loads a template from the given path with caching
func (m *Manager) LoadTemplate(ctx context.Context, templatePath string) (*spookytypes.Template, error) {
	m.logger.Info("Loading template", map[string]interface{}{
		"path": templatePath,
	})

	// Check cache first
	if cached := m.getCachedTemplate(templatePath); cached != nil {
		m.logger.Debug("Returning cached template", map[string]interface{}{
			"path": templatePath,
		})
		return cached, nil
	}

	// Read template file
	data, err := os.ReadFile(templatePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read template file: %w", err)
	}

	// Create template structure
	tmplData := &spookytypes.Template{
		SourcePath: templatePath,
		Content:    string(data),
		ID:         filepath.Base(templatePath),
	}

	// Load and validate metadata
	if err := m.loadTemplateMetadata(ctx, tmplData); err != nil {
		m.logger.Warn("Failed to load template metadata", map[string]interface{}{
			"path":  templatePath,
			"error": err.Error(),
		})
	}

	// Validate template
	if err := m.validateTemplate(ctx, tmplData); err != nil {
		return nil, fmt.Errorf("template validation failed: %w", err)
	}

	// Cache the template
	m.cacheTemplate(templatePath, tmplData)

	m.logger.Info("Template loaded successfully", map[string]interface{}{
		"path": templatePath,
		"size": len(data),
	})

	return tmplData, nil
}

// RenderTemplate renders a template with enhanced features
func (m *Manager) RenderTemplate(ctx context.Context, tmplData *spookytypes.Template, data map[string]interface{}) (string, error) {
	start := time.Now()

	m.logger.Info("Rendering template", map[string]interface{}{
		"template":  tmplData.ID,
		"data_keys": len(data),
	})

	// Check result cache first
	if cached := m.getCachedResult(tmplData, data); cached != "" {
		m.logger.Debug("Returning cached result", map[string]interface{}{
			"template": tmplData.ID,
		})
		return cached, nil
	}

	// Resolve template context
	resolvedData, err := m.resolveTemplateContext(ctx, tmplData, data)
	if err != nil {
		return "", fmt.Errorf("failed to resolve template context: %w", err)
	}

	// Validate security
	if err := m.validateTemplateSecurity(ctx, tmplData); err != nil {
		return "", fmt.Errorf("template security validation failed: %w", err)
	}

	// Compile template with optimization
	compiledTemplate, err := m.compileTemplate(ctx, tmplData)
	if err != nil {
		return "", fmt.Errorf("failed to compile template: %w", err)
	}

	// Run the template with the resolved data
	var buf bytes.Buffer
	if err := compiledTemplate.Execute(&buf, resolvedData); err != nil {
		m.logger.Error("Failed to run template", err, map[string]interface{}{
			"template": tmplData.ID,
			"error":    err.Error(),
		})
		return "", fmt.Errorf("failed to run template %s: %w", tmplData.ID, err)
	}

	result := buf.String()

	// Cache the result
	m.cacheResult(tmplData, data, result)

	// Record performance metrics
	duration := time.Since(start)
	m.recordPerformanceMetric("template_render_time", float64(duration.Milliseconds()))
	m.recordPerformanceMetric("template_render_size", float64(len(result)))

	m.logger.Info("Template rendered successfully", map[string]interface{}{
		"template":      tmplData.ID,
		"result_length": len(result),
		"duration_ms":   duration.Milliseconds(),
	})

	return result, nil
}

// RegisterTemplateFunctions registers custom template functions
func (m *Manager) RegisterTemplateFunctions(functions map[string]interface{}) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for name, fn := range functions {
		if err := m.registerFunction(name, fn); err != nil {
			return fmt.Errorf("failed to register function %s: %w", name, err)
		}
	}

	m.logger.Info("Registered template functions", map[string]interface{}{
		"count": len(functions),
	})

	return nil
}

// ResolveTemplateContext resolves template context with facts, variables, and machines data
func (m *Manager) ResolveTemplateContext(ctx context.Context, template *spookytypes.Template, data map[string]interface{}) (*spookytypes.TemplateContext, error) {
	return m.contextResolver.ResolveContext(ctx, template, data)
}

// ValidateTemplateWithSchema validates template against schemas
func (m *Manager) ValidateTemplateWithSchema(ctx context.Context, template *spookytypes.Template) (*spookytypesschemas.ValidationResult, error) {
	return m.validator.ValidateTemplateComprehensive(ctx, template)
}

// GetTemplateMetadata gets template metadata
func (m *Manager) GetTemplateMetadata(ctx context.Context, templatePath string) (*spookytypestemplates.TemplateMetadata, error) {
	return m.metadataManager.LoadMetadata(ctx, templatePath)
}

// SetFactsIntegration sets the facts integration for context resolution
func (m *Manager) SetFactsIntegration(factsIntegration spookyinterfaces.FactsIntegration) {
	m.contextResolver.factsManager = factsIntegration
}

// SetVariablesIntegration sets the variables integration for context resolution
func (m *Manager) SetVariablesIntegration(variablesIntegration spookyinterfaces.VariablesIntegration) {
	m.contextResolver.variablesManager = variablesIntegration
}

// SetMachinesIntegration sets the machines integration for context resolution
func (m *Manager) SetMachinesIntegration(machinesIntegration spookyinterfaces.MachinesIntegration) {
	m.contextResolver.machinesManager = machinesIntegration
}

// SetSchemaValidator sets the schema validator for template validation
func (m *Manager) SetSchemaValidator(schemaValidator spookytypesschemas.SchemaValidator) {
	m.validator.schemaValidator = schemaValidator
}

// SetSchemaManager sets the schema manager for template validation
func (m *Manager) SetSchemaManager(schemaManager *spookyschemas.Manager) {
	m.validator.schemaManager = schemaManager
}

// Private helper methods

func (m *Manager) getCachedTemplate(templatePath string) *spookytypes.Template {
	m.cache.mu.RLock()
	defer m.cache.mu.RUnlock()

	if entry, exists := m.cache.cache[templatePath]; exists && time.Now().Before(entry.ExpiresAt) {
		entry.AccessCount++
		m.cache.cache[templatePath] = entry
		return entry.Template
	}

	return nil
}

func (m *Manager) cacheTemplate(templatePath string, template *spookytypes.Template) {
	m.cache.mu.Lock()
	defer m.cache.mu.Unlock()

	m.cache.cache[templatePath] = CacheEntry{
		Template:    template,
		ExpiresAt:   time.Now().Add(m.cache.ttl),
		AccessCount: 1,
	}
}

func (m *Manager) getCachedResult(template *spookytypes.Template, data map[string]interface{}) string {
	// Create cache key from template ID and data hash
	cacheKey := fmt.Sprintf("%s:%v", template.ID, data)

	m.performanceManager.cache.mu.RLock()
	defer m.performanceManager.cache.mu.RUnlock()

	if entry, exists := m.performanceManager.cache.cache[cacheKey]; exists && time.Now().Before(entry.ExpiresAt) {
		return entry.Result
	}

	return ""
}

func (m *Manager) cacheResult(template *spookytypes.Template, data map[string]interface{}, result string) {
	// Create cache key from template ID and data hash
	cacheKey := fmt.Sprintf("%s:%v", template.ID, data)

	m.performanceManager.cache.mu.Lock()
	defer m.performanceManager.cache.mu.Unlock()

	m.performanceManager.cache.cache[cacheKey] = ResultCacheEntry{
		Result:    result,
		ExpiresAt: time.Now().Add(m.performanceManager.cache.ttl),
	}
}

func (m *Manager) resolveTemplateContext(_ context.Context, _ *spookytypes.Template, data map[string]interface{}) (map[string]interface{}, error) {
	// For now, return the original data
	// This will be enhanced with facts, variables, and machines data
	return data, nil
}

func (m *Manager) validateTemplateSecurity(_ context.Context, template *spookytypes.Template) error {
	// For now, basic security check
	if strings.Contains(template.Content, "exec") || strings.Contains(template.Content, "system") {
		return fmt.Errorf("template contains forbidden patterns")
	}
	return nil
}

func (m *Manager) compileTemplate(_ context.Context, template *spookytypes.Template) (*texttemplate.Template, error) {
	// Check compilation cache
	m.performanceManager.compiler.mu.RLock()
	if cached, exists := m.performanceManager.compiler.cache[template.ID]; exists {
		m.performanceManager.compiler.mu.RUnlock()
		return cached, nil
	}
	m.performanceManager.compiler.mu.RUnlock()

	// Compile template
	tmpl, err := texttemplate.New(template.ID).Parse(template.Content)
	if err != nil {
		return nil, err
	}

	// Cache compiled template
	m.performanceManager.compiler.mu.Lock()
	m.performanceManager.compiler.cache[template.ID] = tmpl
	m.performanceManager.compiler.mu.Unlock()

	return tmpl, nil
}

func (m *Manager) loadTemplateMetadata(_ context.Context, template *spookytypes.Template) error {
	// For now, create basic metadata
	// This will be enhanced to load from metadata files
	template.Metadata = &spookytypestemplates.TemplateMetadata{
		Name:        template.ID,
		Description: fmt.Sprintf("Template loaded from %s", template.SourcePath),
		Version:     "1.0.0",
		Tags:        []string{},
	}
	return nil
}

func (m *Manager) validateTemplate(_ context.Context, template *spookytypes.Template) error {
	// Basic validation
	if template.Content == "" {
		return fmt.Errorf("template content is empty")
	}

	// Validate template syntax
	if _, err := texttemplate.New(template.ID).Parse(template.Content); err != nil {
		return fmt.Errorf("invalid template syntax: %w", err)
	}

	return nil
}

func (m *Manager) registerBuiltInFunctions() {
	// Register basic built-in functions
	basicFunctions := map[string]interface{}{
		"upper": strings.ToUpper,
		"lower": strings.ToLower,
		"trim":  strings.TrimSpace,
		"len":   func(s string) int { return len(s) },
	}

	for name, fn := range basicFunctions {
		if err := m.registerFunction(name, fn); err != nil {
			m.logger.Error("Failed to register function", err, map[string]interface{}{
				"function": name,
			})
		}
	}
}

func (m *Manager) registerFunction(name string, fn interface{}) error {
	m.functions.mu.Lock()
	defer m.functions.mu.Unlock()

	m.functions.functions[name] = TemplateFunction{
		Name:        name,
		Function:    fn,
		Security:    FunctionSecurity{AllowedInRestricted: true},
		Description: fmt.Sprintf("Function: %s", name),
	}

	return nil
}

func (m *Manager) recordPerformanceMetric(name string, value float64) {
	m.performanceManager.monitor.mu.Lock()
	defer m.performanceManager.monitor.mu.Unlock()

	metric, exists := m.performanceManager.monitor.metrics[name]
	if !exists {
		metric = Metric{
			Name: name,
			Min:  value,
			Max:  value,
		}
	}

	metric.Value = value
	metric.Count++
	metric.LastUpdate = time.Now()

	if value < metric.Min {
		metric.Min = value
	}
	if value > metric.Max {
		metric.Max = value
	}

	m.performanceManager.monitor.metrics[name] = metric
}

// TemplateContextResolver helper methods

// generateCacheKey generates a cache key for the context
func (c *TemplateContextResolver) generateCacheKey(template *spookytypes.Template, _ map[string]interface{}) string {
	// Simple cache key generation - in a real implementation, this would be more sophisticated
	return fmt.Sprintf("%s:%s:%s", template.ID, template.SourcePath, template.Type)
}

// getCachedContext gets a cached context
func (c *TemplateContextResolver) getCachedContext(key string) (*spookytypes.TemplateContext, bool) {
	c.cache.mu.RLock()
	defer c.cache.mu.RUnlock()

	entry, exists := c.cache.cache[key]
	if !exists {
		return nil, false
	}

	// Check if cache entry has expired
	if time.Now().After(entry.ExpiresAt) {
		delete(c.cache.cache, key)
		return nil, false
	}

	return entry.Context, true
}

// cacheContext caches a context
func (c *TemplateContextResolver) cacheContext(key string, context *spookytypes.TemplateContext) {
	c.cache.mu.Lock()
	defer c.cache.mu.Unlock()

	c.cache.cache[key] = ContextCacheEntry{
		Context:   context,
		ExpiresAt: time.Now().Add(c.cache.ttl),
	}
}

// resolveProjectContext resolves project information
func (c *TemplateContextResolver) resolveProjectContext(_ context.Context, template *spookytypes.Template) (map[string]interface{}, error) {
	project := make(map[string]interface{})

	// Add basic project information
	if template.Metadata != nil {
		project["name"] = template.Metadata.Name
		project["description"] = template.Metadata.Description
		project["author"] = template.Metadata.Author
		project["version"] = template.Metadata.Version
		project["tags"] = template.Metadata.Tags
		project["license"] = template.Metadata.License
	}

	// Add template-specific information
	project["template_id"] = template.ID
	project["template_type"] = template.Type
	project["template_scope"] = template.Scope
	project["template_security_level"] = template.SecurityLevel
	project["template_engine"] = template.Engine

	return project, nil
}

// resolveFactsContext resolves facts data
func (c *TemplateContextResolver) resolveFactsContext(ctx context.Context, _ *spookytypes.Template) (map[string]interface{}, error) {
	if c.factsManager == nil {
		return make(map[string]interface{}), nil
	}

	// Load facts from the facts manager
	facts, err := c.factsManager.LoadFacts(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load facts: %w", err)
	}

	// Convert facts to map[string]interface{}
	factsMap, ok := facts.(map[string]interface{})
	if !ok {
		// If facts is not a map, create an empty map
		factsMap = make(map[string]interface{})
	}

	return factsMap, nil
}

// resolveMachinesContext resolves machines data
func (c *TemplateContextResolver) resolveMachinesContext(ctx context.Context, _ *spookytypes.Template) ([]map[string]interface{}, error) {
	if c.machinesManager == nil {
		return make([]map[string]interface{}, 0), nil
	}

	// Load machines from the machines manager
	machines, err := c.machinesManager.GetFullInventory(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load machines: %w", err)
	}

	// Convert machines to []map[string]interface{}
	machinesMap := make([]map[string]interface{}, len(machines))
	for i := range machines {
		machineMap := make(map[string]interface{})
		machineMap["hostname"] = machines[i].Hostname
		machineMap["port"] = machines[i].Port
		machineMap["user"] = machines[i].User
		machineMap["tags"] = machines[i].Tags
		// Add other machine fields as needed
		machinesMap[i] = machineMap
	}

	return machinesMap, nil
}

// resolveEnvironmentContext resolves environment variables
func (c *TemplateContextResolver) resolveEnvironmentContext(_ context.Context, template *spookytypes.Template) (map[string]string, error) {
	environment := make(map[string]string)

	// Add common environment variables
	environment["HOME"] = getEnvOrDefault("HOME", "")
	environment["USER"] = getEnvOrDefault("USER", "")
	environment["PWD"] = getEnvOrDefault("PWD", "")
	environment["PATH"] = getEnvOrDefault("PATH", "")

	// Add template-specific environment variables
	if template.Variables != nil {
		for key, value := range template.Variables {
			if strValue, ok := value.(string); ok {
				environment["TEMPLATE_"+key] = strValue
			}
		}
	}

	return environment, nil
}

// resolveVariablesContext resolves variables data
func (c *TemplateContextResolver) resolveVariablesContext(ctx context.Context, _ *spookytypes.Template) (map[string]interface{}, error) {
	if c.variablesManager == nil {
		return make(map[string]interface{}), nil
	}

	// Load variables from the variables manager
	variables, err := c.variablesManager.LoadVariables(ctx, "")
	if err != nil {
		return nil, fmt.Errorf("failed to load variables: %w", err)
	}

	// Convert variables to map[string]interface{}
	variablesMap := make(map[string]interface{})
	for key, variable := range variables {
		variablesMap[key] = variable.ResolvedValue
	}

	return variablesMap, nil
}

// transformProjectContext transforms project context
func (c *TemplateContextResolver) transformProjectContext(_ map[string]interface{}) error {
	// Add any project-specific transformations here
	return nil
}

// transformFactsContext transforms facts context
func (c *TemplateContextResolver) transformFactsContext(_ map[string]interface{}) error {
	// Add any facts-specific transformations here
	return nil
}

// transformMachinesContext transforms machines context
func (c *TemplateContextResolver) transformMachinesContext(_ []map[string]interface{}) error {
	// Add any machines-specific transformations here
	return nil
}

// transformEnvironmentContext transforms environment context
func (c *TemplateContextResolver) transformEnvironmentContext(_ map[string]string) error {
	// Add any environment-specific transformations here
	return nil
}

// transformVariablesContext transforms variables context
func (c *TemplateContextResolver) transformVariablesContext(_ map[string]interface{}) error {
	// Add any variables-specific transformations here
	return nil
}

// getEnvOrDefault gets an environment variable with a default value
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// Helper functions for metadata management

// isValidVersion checks if a version string is valid ScalVer format
func isValidVersion(version string) bool {
	if version == "" {
		return false
	}

	// Check if it's a valid ScalVer format
	return spookytypescommon.IsValidScalVerFormat(version)
}

// Note: ScalVer functionality is now centralized in spookytypescommon package
// Use spookytypescommon.ScalVer, spookytypescommon.ParseScalVer, etc.

// ValidateTemplate validates a template
func (m *Manager) ValidateTemplate(ctx context.Context, template *spookytypes.Template) (*spookytypes.ValidationResult, error) {
	m.logger.Info("Validating template", map[string]interface{}{
		"template": template.ID,
	})

	// Use comprehensive schema validation
	schemaResult, err := m.ValidateTemplateWithSchema(ctx, template)
	if err != nil {
		m.logger.Error("Schema validation failed", err, map[string]interface{}{
			"template": template.ID,
		})
		return &spookytypes.ValidationResult{
			Valid:    false,
			Errors:   []spookytypes.SchemaError{{Message: fmt.Sprintf("schema validation failed: %v", err)}},
			Warnings: []spookytypes.SchemaError{},
		}, nil
	}

	// Convert schema validation result to interface validation result
	var errors []spookytypes.SchemaError
	var warnings []spookytypes.SchemaError

	// Convert schema errors
	for i := range schemaResult.Errors {
		errors = append(errors, spookytypes.SchemaError{
			Message: schemaResult.Errors[i].Message,
		})
	}

	// Convert schema warnings
	for i := range schemaResult.Warnings {
		warnings = append(warnings, spookytypes.SchemaError{
			Message: schemaResult.Warnings[i].Message,
		})
	}

	m.logger.Info("Template validation completed", map[string]interface{}{
		"template": template.ID,
		"valid":    schemaResult.Valid,
		"errors":   len(errors),
		"warnings": len(warnings),
	})

	return &spookytypes.ValidationResult{
		Valid:    schemaResult.Valid,
		Errors:   errors,
		Warnings: warnings,
	}, nil
}

// EnhancedMetadataIndexer provides multi-dimensional indexing for template metadata
type EnhancedMetadataIndexer struct {
	// Multiple indexes for different search dimensions
	byName         map[string]*spookytypestemplates.TemplateMetadata
	byTags         map[string][]*spookytypestemplates.TemplateMetadata
	byCategory     map[string][]*spookytypestemplates.TemplateMetadata
	bySubcategory  map[string][]*spookytypestemplates.TemplateMetadata
	byAuthor       map[string][]*spookytypestemplates.TemplateMetadata
	byKeywords     map[string][]*spookytypestemplates.TemplateMetadata
	byDependencies map[string][]*spookytypestemplates.TemplateMetadata

	// Full-text search index (simple implementation)
	fullTextIndex map[string][]*spookytypestemplates.TemplateMetadata

	// Priority-based index
	byPriority map[int][]*spookytypestemplates.TemplateMetadata

	mu sync.RWMutex
}

// NewEnhancedMetadataIndexer creates a new enhanced metadata indexer
func NewEnhancedMetadataIndexer() *EnhancedMetadataIndexer {
	return &EnhancedMetadataIndexer{
		byName:         make(map[string]*spookytypestemplates.TemplateMetadata),
		byTags:         make(map[string][]*spookytypestemplates.TemplateMetadata),
		byCategory:     make(map[string][]*spookytypestemplates.TemplateMetadata),
		bySubcategory:  make(map[string][]*spookytypestemplates.TemplateMetadata),
		byAuthor:       make(map[string][]*spookytypestemplates.TemplateMetadata),
		byKeywords:     make(map[string][]*spookytypestemplates.TemplateMetadata),
		byDependencies: make(map[string][]*spookytypestemplates.TemplateMetadata),
		fullTextIndex:  make(map[string][]*spookytypestemplates.TemplateMetadata),
		byPriority:     make(map[int][]*spookytypestemplates.TemplateMetadata),
	}
}

// IndexMetadata indexes template metadata across all dimensions
func (i *EnhancedMetadataIndexer) IndexMetadata(metadata *spookytypestemplates.TemplateMetadata) error {
	i.mu.Lock()
	defer i.mu.Unlock()

	// Index by name
	i.byName[metadata.Name] = metadata

	// Index by tags
	if metadata.Tags != nil {
		for _, tag := range metadata.Tags {
			i.byTags[tag] = append(i.byTags[tag], metadata)
		}
	}

	// Index by category
	if metadata.Category != "" {
		i.byCategory[metadata.Category] = append(i.byCategory[metadata.Category], metadata)
	}

	// Index by subcategory
	if metadata.Subcategory != "" {
		i.bySubcategory[metadata.Subcategory] = append(i.bySubcategory[metadata.Subcategory], metadata)
	}

	// Index by author
	if metadata.Author != "" {
		i.byAuthor[metadata.Author] = append(i.byAuthor[metadata.Author], metadata)
	}

	// Index by keywords
	if metadata.Keywords != nil {
		for _, keyword := range metadata.Keywords {
			i.byKeywords[keyword] = append(i.byKeywords[keyword], metadata)
		}
	}

	// Index by dependencies
	if metadata.Dependencies != nil {
		for _, dep := range metadata.Dependencies {
			i.byDependencies[dep] = append(i.byDependencies[dep], metadata)
		}
	}

	// Index by priority
	i.byPriority[metadata.Priority] = append(i.byPriority[metadata.Priority], metadata)

	// Build full-text index
	i.buildFullTextIndex(metadata)

	return nil
}

// buildFullTextIndex builds full-text search index
func (i *EnhancedMetadataIndexer) buildFullTextIndex(metadata *spookytypestemplates.TemplateMetadata) {
	// Index name
	if metadata.Name != "" {
		words := strings.Fields(strings.ToLower(metadata.Name))
		for _, word := range words {
			i.fullTextIndex[word] = append(i.fullTextIndex[word], metadata)
		}
	}

	// Index description
	if metadata.Description != "" {
		words := strings.Fields(strings.ToLower(metadata.Description))
		for _, word := range words {
			i.fullTextIndex[word] = append(i.fullTextIndex[word], metadata)
		}
	}

	// Index tags
	if metadata.Tags != nil {
		for _, tag := range metadata.Tags {
			words := strings.Fields(strings.ToLower(tag))
			for _, word := range words {
				i.fullTextIndex[word] = append(i.fullTextIndex[word], metadata)
			}
		}
	}

	// Index keywords
	if metadata.Keywords != nil {
		for _, keyword := range metadata.Keywords {
			words := strings.Fields(strings.ToLower(keyword))
			for _, word := range words {
				i.fullTextIndex[word] = append(i.fullTextIndex[word], metadata)
			}
		}
	}
}

// SearchResult represents a search result with relevance score
type SearchResult struct {
	Metadata *spookytypestemplates.TemplateMetadata
	Score    float64
	Matched  []string // What matched in the search
}

// Search searches for templates with multiple algorithms
func (i *EnhancedMetadataIndexer) Search(query string, filters *SearchFilters) ([]SearchResult, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()

	var results []SearchResult
	seen := make(map[string]bool)

	// Search by name (exact and partial)
	if nameResults := i.searchByName(query); len(nameResults) > 0 {
		for _, result := range nameResults {
			if !seen[result.Metadata.Name] {
				results = append(results, result)
				seen[result.Metadata.Name] = true
			}
		}
	}

	// Search by tags
	if tagResults := i.searchByTags(query); len(tagResults) > 0 {
		for _, result := range tagResults {
			if !seen[result.Metadata.Name] {
				results = append(results, result)
				seen[result.Metadata.Name] = true
			}
		}
	}

	// Search by keywords
	if keywordResults := i.searchByKeywords(query); len(keywordResults) > 0 {
		for _, result := range keywordResults {
			if !seen[result.Metadata.Name] {
				results = append(results, result)
				seen[result.Metadata.Name] = true
			}
		}
	}

	// Search by full-text
	if fullTextResults := i.searchByFullText(query); len(fullTextResults) > 0 {
		for _, result := range fullTextResults {
			if !seen[result.Metadata.Name] {
				results = append(results, result)
				seen[result.Metadata.Name] = true
			}
		}
	}

	// Apply filters
	results = i.applyFilters(results, filters)

	// Sort by score
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	return results, nil
}

// SearchFilters represents search filters
type SearchFilters struct {
	Tags        []string
	Category    string
	Subcategory string
	Author      string
	MinScore    float64
	Limit       int
}

// searchByName searches by template name
func (i *EnhancedMetadataIndexer) searchByName(query string) []SearchResult {
	var results []SearchResult
	queryLower := strings.ToLower(query)

	for name, metadata := range i.byName {
		if strings.Contains(strings.ToLower(name), queryLower) {
			score := 1.0
			if strings.EqualFold(name, query) {
				score = 2.0 // Exact match gets higher score
			}
			results = append(results, SearchResult{
				Metadata: metadata,
				Score:    score,
				Matched:  []string{"name"},
			})
		}
	}

	return results
}

// searchByIndex searches by a specific index with configurable scoring
func (i *EnhancedMetadataIndexer) searchByIndex(index map[string][]*spookytypestemplates.TemplateMetadata, query, matchType string, exactScore, partialScore float64) []SearchResult {
	var results []SearchResult
	queryLower := strings.ToLower(query)

	// Search for exact matches
	if templates, exists := index[query]; exists {
		for _, metadata := range templates {
			results = append(results, SearchResult{
				Metadata: metadata,
				Score:    exactScore,
				Matched:  []string{matchType},
			})
		}
	}

	// Search for partial matches
	for key, templates := range index {
		if strings.Contains(strings.ToLower(key), queryLower) && key != query {
			for _, metadata := range templates {
				results = append(results, SearchResult{
					Metadata: metadata,
					Score:    partialScore,
					Matched:  []string{matchType},
				})
			}
		}
	}

	return results
}

// searchByTags searches by template tags
func (i *EnhancedMetadataIndexer) searchByTags(query string) []SearchResult {
	return i.searchByIndex(i.byTags, query, "tag", 1.5, 1.0)
}

// searchByKeywords searches by template keywords
func (i *EnhancedMetadataIndexer) searchByKeywords(query string) []SearchResult {
	return i.searchByIndex(i.byKeywords, query, "keyword", 1.3, 1.0)
}

// searchByFullText searches using full-text index
func (i *EnhancedMetadataIndexer) searchByFullText(query string) []SearchResult {
	var results []SearchResult
	queryLower := strings.ToLower(query)
	words := strings.Fields(queryLower)

	// Count matches for each template
	templateScores := make(map[string]float64)
	templateMatches := make(map[string][]string)

	for _, word := range words {
		if templates, exists := i.fullTextIndex[word]; exists {
			for _, metadata := range templates {
				templateScores[metadata.Name] += 0.5
				templateMatches[metadata.Name] = append(templateMatches[metadata.Name], word)
			}
		}
	}

	// Convert to results
	for name, score := range templateScores {
		// Find the metadata for this name
		if metadata, exists := i.byName[name]; exists {
			results = append(results, SearchResult{
				Metadata: metadata,
				Score:    score,
				Matched:  templateMatches[name],
			})
		}
	}

	return results
}

// applyFilters applies search filters to results
func (i *EnhancedMetadataIndexer) applyFilters(results []SearchResult, filters *SearchFilters) []SearchResult {
	var filtered []SearchResult

	for _, result := range results {
		// Apply tag filter
		if len(filters.Tags) > 0 {
			hasTag := false
			for _, filterTag := range filters.Tags {
				for _, templateTag := range result.Metadata.Tags {
					if filterTag == templateTag {
						hasTag = true
						break
					}
				}
				if hasTag {
					break
				}
			}
			if !hasTag {
				continue
			}
		}

		// Apply category filter
		if filters.Category != "" && result.Metadata.Category != filters.Category {
			continue
		}

		// Apply subcategory filter
		if filters.Subcategory != "" && result.Metadata.Subcategory != filters.Subcategory {
			continue
		}

		// Apply author filter
		if filters.Author != "" && result.Metadata.Author != filters.Author {
			continue
		}

		// Apply minimum score filter
		if result.Score < filters.MinScore {
			continue
		}

		filtered = append(filtered, result)

		// Apply limit
		if filters.Limit > 0 && len(filtered) >= filters.Limit {
			break
		}
	}

	return filtered
}
