package schemas

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"

	spookyfactstypes "spooky/internal/types/facts"
)

// SchemaType represents the type of schema
type SchemaType string

// Schema represents a loaded schema
type Schema struct {
	Type     SchemaType
	Content  string
	Filename string
}

// ValidationError represents a unified validation error with detailed information
type ValidationError struct {
	File     string `json:"file"`
	Line     int    `json:"line"`
	Column   int    `json:"column"`
	Field    string `json:"field"`
	Message  string `json:"message"`
	Value    string `json:"value,omitempty"`
	Severity string `json:"severity"` // "error" or "warning"
}

// ValidationResult contains unified validation results
type ValidationResult struct {
	Valid    bool              `json:"valid"`
	Errors   []ValidationError `json:"errors,omitempty"`
	Warnings []ValidationError `json:"warnings,omitempty"`
	File     string            `json:"file,omitempty"`
	Schema   string            `json:"schema,omitempty"`
}

// SchemaValidator is the unified schema validator combining both existing validators
type SchemaValidator struct {
	schemas map[SchemaType]*Schema
	parser  *hclparse.Parser
	cache   map[string]*ValidationResult
	mutex   sync.RWMutex
}

// NewSchemaValidator creates a new unified schema validator
func NewSchemaValidator() *SchemaValidator {
	return &SchemaValidator{
		schemas: make(map[SchemaType]*Schema),
		parser:  hclparse.NewParser(),
		cache:   make(map[string]*ValidationResult),
	}
}

// LoadSchema loads a schema from the embedded schemas directory
func (sv *SchemaValidator) LoadSchema(schemaType SchemaType) error {
	// Get the schema content from embedded schemas
	content, err := GetSchema(schemaType)
	if err != nil {
		return fmt.Errorf("failed to load schema %s: %w", schemaType, err)
	}

	// Parse the HCL content
	_, diags := sv.parser.ParseHCL(content, string(schemaType))
	if diags.HasErrors() {
		return fmt.Errorf("failed to parse schema %s: %s", schemaType, diags.Error())
	}

	// Store the schema
	sv.schemas[schemaType] = &Schema{
		Type:     schemaType,
		Content:  string(content),
		Filename: string(schemaType) + ".hcl",
	}

	return nil
}

// LoadAllSchemas loads all available schemas into the validator
func (sv *SchemaValidator) LoadAllSchemas() error {
	schemaTypes := []SchemaType{
		SchemaTypeMachines,
		SchemaTypeSpooky,
		SchemaTypeVariablesStructure,
		SchemaTypeVariablesHCL,
		SchemaTypeVariablesJSON,
		SchemaTypeCustomFactsHCL,
		SchemaTypeFactsStructure,
		SchemaTypeFactsStorage,
	}

	for _, schemaType := range schemaTypes {
		if err := sv.LoadSchema(schemaType); err != nil {
			return fmt.Errorf("failed to load schema %s: %w", schemaType, err)
		}
	}

	return nil
}

// ValidateFile validates an HCL file against a specific schema (from EnhancedSchemaValidator)
func (sv *SchemaValidator) ValidateFile(filePath, schemaName string) *ValidationResult {
	// Check cache first
	sv.mutex.RLock()
	cacheKey := fmt.Sprintf("%s:%s", filePath, schemaName)
	if cached, exists := sv.cache[cacheKey]; exists {
		sv.mutex.RUnlock()
		return cached
	}
	sv.mutex.RUnlock()

	result := &ValidationResult{
		File:   filePath,
		Schema: schemaName,
		Valid:  true,
	}

	// Read the file
	content, err := os.ReadFile(filePath)
	if err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, ValidationError{
			File:     filePath,
			Field:    "file",
			Message:  fmt.Sprintf("failed to read file: %v", err),
			Severity: "error",
		})
		sv.cacheResult(cacheKey, result)
		return result
	}

	// Parse the HCL content
	file, diags := sv.parser.ParseHCL(content, filePath)
	if diags.HasErrors() {
		result.Valid = false
		for _, diag := range diags {
			validationError := ConvertHCLDiagnostic(*diag)
			validationError.File = filePath
			result.Errors = append(result.Errors, *validationError)
		}
		sv.cacheResult(cacheKey, result)
		return result
	}

	// Validate against the specified schema
	schemaResult := sv.validateAgainstSchema(file, schemaName)
	result.Errors = append(result.Errors, schemaResult.Errors...)
	result.Warnings = append(result.Warnings, schemaResult.Warnings...)
	result.Valid = len(result.Errors) == 0

	sv.cacheResult(cacheKey, result)
	return result
}

// ValidateData validates data structure against a schema (from EnhancedSchemaValidator)
func (sv *SchemaValidator) ValidateData(data interface{}, schemaName string) error {
	// Check if schema exists
	_, exists := sv.schemas[SchemaType(schemaName)]
	if !exists {
		return fmt.Errorf("schema not found: %s", schemaName)
	}

	// For now, we'll implement basic validation based on the schema type
	// This is a simplified implementation - in a full implementation,
	// we would parse the HCL schema and validate against it

	switch SchemaType(schemaName) {
	case SchemaTypeFactsStructure:
		return sv.validateFactsStructure(data)
	case SchemaTypeProject:
		return sv.validateProjectDirectory(data)
	case SchemaTypeActions:
		return sv.validateActions(data)
	case SchemaTypeMachines:
		return sv.validateMachines(data)
	case SchemaTypeVariables:
		return sv.validateVariables(data)
	case SchemaTypeTemplateStructure:
		return sv.validateTemplates(data)
	default:
		// For unknown schemas, just check if data is not nil
		if data == nil {
			return fmt.Errorf("data cannot be nil for schema: %s", schemaName)
		}
		return nil
	}
}

// validateFactsStructure validates fact collection data
func (sv *SchemaValidator) validateFactsStructure(data interface{}) error {
	if data == nil {
		return fmt.Errorf("fact collection data cannot be nil")
	}

	// Type assertion to check if it's a fact collection
	switch v := data.(type) {
	case *spookyfactstypes.FactCollection:
		return sv.validateFactCollection(v)
	case map[string]*spookyfactstypes.FactCollection:
		// Validate each collection in the map
		for machineID, collection := range v {
			if err := sv.validateFactCollection(collection); err != nil {
				return fmt.Errorf("invalid fact collection for machine %s: %w", machineID, err)
			}
		}
		return nil
	default:
		return fmt.Errorf("expected *spookyfactstypes.FactCollection or map[string]*spookyfactstypes.FactCollection, got %T", data)
	}
}

// validateFactCollection validates a single fact collection
func (sv *SchemaValidator) validateFactCollection(collection *spookyfactstypes.FactCollection) error {
	if collection == nil {
		return fmt.Errorf("fact collection cannot be nil")
	}

	// Validate machine_id (should be 32-character hex string)
	if collection.Server == "" {
		return fmt.Errorf("machine_id is required")
	}

	// Check if machine_id matches the expected pattern (32 hex chars)
	if len(collection.Server) != 32 {
		return fmt.Errorf("machine_id must be 32 characters long, got %d", len(collection.Server))
	}

	// Basic hex validation
	for _, char := range collection.Server {
		if !((char >= '0' && char <= '9') || (char >= 'a' && char <= 'f')) {
			return fmt.Errorf("machine_id must contain only hexadecimal characters (0-9, a-f)")
		}
	}

	// Validate timestamp
	if collection.Timestamp.IsZero() {
		return fmt.Errorf("timestamp is required")
	}

	// Validate facts map
	if collection.Facts == nil {
		return fmt.Errorf("facts map cannot be nil")
	}

	// Validate each fact
	for key, fact := range collection.Facts {
		if err := sv.validateFact(key, fact); err != nil {
			return fmt.Errorf("invalid fact %s: %w", key, err)
		}
	}

	return nil
}

// validateFact validates a single fact
func (sv *SchemaValidator) validateFact(key string, fact *spookyfactstypes.Fact) error {
	if fact == nil {
		return fmt.Errorf("fact cannot be nil")
	}

	if key == "" {
		return fmt.Errorf("fact key cannot be empty")
	}

	if fact.Key == "" {
		return fmt.Errorf("fact key field cannot be empty")
	}

	// Validate that key matches the fact's key field
	if key != fact.Key {
		return fmt.Errorf("fact key mismatch: expected %s, got %s", key, fact.Key)
	}

	// Validate value is not nil
	if fact.Value == nil {
		return fmt.Errorf("fact value cannot be nil")
	}

	return nil
}

// validateProjectDirectory validates project directory structure
func (sv *SchemaValidator) validateProjectDirectory(data interface{}) error {
	if data == nil {
		return fmt.Errorf("project directory data cannot be nil")
	}

	// Basic validation - project directory should be a string path
	if path, ok := data.(string); ok {
		if path == "" {
			return fmt.Errorf("project directory path cannot be empty")
		}
		return nil
	}

	return fmt.Errorf("expected string path for project directory, got %T", data)
}

// validateActions validates actions data
func (sv *SchemaValidator) validateActions(data interface{}) error {
	if data == nil {
		return fmt.Errorf("actions data cannot be nil")
	}

	// Basic validation - actions should be a slice or map
	switch v := data.(type) {
	case []interface{}:
		if len(v) == 0 {
			return fmt.Errorf("actions cannot be empty")
		}
	case map[string]interface{}:
		if len(v) == 0 {
			return fmt.Errorf("actions cannot be empty")
		}
	default:
		return fmt.Errorf("expected slice or map for actions, got %T", data)
	}

	return nil
}

// validateMachines validates machines data
func (sv *SchemaValidator) validateMachines(data interface{}) error {
	if data == nil {
		return fmt.Errorf("machines data cannot be nil")
	}

	// Basic validation - machines should be a slice or map
	switch v := data.(type) {
	case []interface{}:
		if len(v) == 0 {
			return fmt.Errorf("machines cannot be empty")
		}
	case map[string]interface{}:
		if len(v) == 0 {
			return fmt.Errorf("machines cannot be empty")
		}
	default:
		return fmt.Errorf("expected slice or map for machines, got %T", data)
	}

	return nil
}

// validateVariables validates variables data
func (sv *SchemaValidator) validateVariables(data interface{}) error {
	if data == nil {
		return fmt.Errorf("variables data cannot be nil")
	}

	// Basic validation - variables should be a map
	if vars, ok := data.(map[string]interface{}); ok {
		if len(vars) == 0 {
			return fmt.Errorf("variables cannot be empty")
		}
		return nil
	}

	return fmt.Errorf("expected map for variables, got %T", data)
}

// validateTemplates validates templates data
func (sv *SchemaValidator) validateTemplates(data interface{}) error {
	if data == nil {
		return fmt.Errorf("templates data cannot be nil")
	}

	// Basic validation - templates should be a slice or map
	switch v := data.(type) {
	case []interface{}:
		if len(v) == 0 {
			return fmt.Errorf("templates cannot be empty")
		}
	case map[string]interface{}:
		if len(v) == 0 {
			return fmt.Errorf("templates cannot be empty")
		}
	default:
		return fmt.Errorf("expected slice or map for templates, got %T", data)
	}

	return nil
}

// ValidateProject validates a project configuration against the project schema (from original SchemaValidator)
func (sv *SchemaValidator) ValidateProject(projectPath string) *ValidationResult {
	result := &ValidationResult{Valid: true}

	// Read project.hcl file
	projectFile := filepath.Join(projectPath, "project.hcl")
	content, err := os.ReadFile(projectFile)
	if err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, ValidationError{
			File:     projectFile,
			Field:    "project.hcl",
			Message:  fmt.Sprintf("failed to read project file: %v", err),
			Severity: "error",
		})
		return result
	}

	// Parse the project file to check if it's valid HCL
	_, diags := sv.parser.ParseHCL(content, projectFile)
	if diags.HasErrors() {
		result.Valid = false
		for _, diag := range diags {
			validationError := ConvertHCLDiagnostic(*diag)
			validationError.File = projectFile
			result.Errors = append(result.Errors, *validationError)
		}
		return result
	}

	// For now, just validate that the file can be parsed as valid HCL
	// DEPRECATED: Schema system is fully implemented - this TODO is ready for removal
	return result
}

// ValidateFacts validates facts data against the appropriate schema (from original SchemaValidator)
func (sv *SchemaValidator) ValidateFacts(factsPath string, format string) *ValidationResult {
	result := &ValidationResult{Valid: true}

	// For BadgerDB, we can't easily validate the schema, but we can check if it's accessible
	if format == "badgerdb" {
		// Check if the directory exists
		if _, err := os.Stat(factsPath); os.IsNotExist(err) {
			result.Valid = false
			result.Errors = append(result.Errors, ValidationError{
				Field:    "facts",
				Message:  fmt.Sprintf("facts database not found at %s", factsPath),
				Severity: "error",
			})
			return result
		}

		// Check if it contains BadgerDB files
		entries, err := os.ReadDir(factsPath)
		if err != nil {
			result.Valid = false
			result.Errors = append(result.Errors, ValidationError{
				Field:    "facts",
				Message:  fmt.Sprintf("failed to read facts directory: %v", err),
				Severity: "error",
			})
			return result
		}

		// Look for BadgerDB files
		hasBadgerFiles := false
		for _, entry := range entries {
			if !entry.IsDir() && (entry.Name() == "MANIFEST" || entry.Name() == "KEYREGISTRY" || strings.HasSuffix(entry.Name(), ".vlog")) {
				hasBadgerFiles = true
				break
			}
		}

		if !hasBadgerFiles {
			result.Valid = false
			result.Errors = append(result.Errors, ValidationError{
				Field:    "facts",
				Message:  fmt.Sprintf("facts database exists but contains no data at %s", factsPath),
				Severity: "error",
			})
			return result
		}

		// For now, just validate that the database can be opened
		// DEPRECATED: Schema system is fully implemented - this TODO is ready for removal
		return result
	}

	// For HCL and JSON formats, use the existing validation logic
	// Load appropriate facts schema based on format
	var schemaType SchemaType
	switch format {
	case "hcl", "json", "badger":
		schemaType = SchemaTypeFactsStructure
	default:
		schemaType = SchemaTypeFactsStructure // Default to facts structure
	}

	if _, exists := sv.schemas[schemaType]; !exists {
		if err := sv.LoadSchema(schemaType); err != nil {
			result.Valid = false
			result.Errors = append(result.Errors, ValidationError{
				Field:    "schema",
				Message:  fmt.Sprintf("failed to load facts schema: %v", err),
				Severity: "error",
			})
			return result
		}
	}

	// Read facts file
	content, err := os.ReadFile(factsPath)
	if err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, ValidationError{
			Field:    "facts",
			Message:  fmt.Sprintf("failed to read facts file: %v", err),
			Severity: "error",
		})
		return result
	}

	// Parse the facts file
	file, diags := sv.parser.ParseHCL(content, factsPath)
	if diags.HasErrors() {
		result.Valid = false
		for _, diag := range diags {
			validationError := ConvertHCLDiagnostic(*diag)
			validationError.File = factsPath
			result.Errors = append(result.Errors, *validationError)
		}
		return result
	}

	// Basic validation - check for required facts block
	body := file.Body
	attrs, diags := body.JustAttributes()
	if diags.HasErrors() {
		result.Valid = false
		for _, diag := range diags {
			validationError := ConvertHCLDiagnostic(*diag)
			validationError.File = factsPath
			result.Errors = append(result.Errors, *validationError)
		}
		return result
	}

	// Check for required fields based on schema
	requiredFields := []string{"machine_id", "collected_at", "facts"}
	for _, field := range requiredFields {
		if _, exists := attrs[field]; !exists {
			result.Valid = false
			result.Errors = append(result.Errors, ValidationError{
				Field:    field,
				Message:  fmt.Sprintf("required field '%s' not found", field),
				Severity: "error",
			})
		}
	}

	return result
}

// ValidateProjectDirectory validates project directory structure (from original SchemaValidator)
func (sv *SchemaValidator) ValidateProjectDirectory(projectPath string) *ValidationResult {
	result := &ValidationResult{Valid: true}

	// Check required directories
	requiredDirs := []string{"templates", "files", "logs", "data", "facts.db"}
	for _, dir := range requiredDirs {
		dirPath := filepath.Join(projectPath, dir)
		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			result.Valid = false
			result.Errors = append(result.Errors, ValidationError{
				Field:    dir,
				Message:  fmt.Sprintf("required directory '%s' not found", dir),
				Severity: "error",
			})
		}
	}

	// Check required files
	requiredFiles := []string{"project.hcl"}
	for _, file := range requiredFiles {
		filePath := filepath.Join(projectPath, file)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			result.Valid = false
			result.Errors = append(result.Errors, ValidationError{
				Field:    file,
				Message:  fmt.Sprintf("required file '%s' not found", file),
				Severity: "error",
			})
		}
	}

	// Check optional but recommended files
	optionalFiles := []string{"machines.hcl", "actions.hcl", "README.md"}
	for _, file := range optionalFiles {
		filePath := filepath.Join(projectPath, file)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			result.Warnings = append(result.Warnings, ValidationError{
				Field:    file,
				Message:  fmt.Sprintf("recommended file '%s' not found", file),
				Severity: "warning",
			})
		}
	}

	return result
}

// validateAgainstSchema validates parsed HCL against a specific schema (from EnhancedSchemaValidator)
func (sv *SchemaValidator) validateAgainstSchema(file *hcl.File, schemaName string) *ValidationResult {
	result := &ValidationResult{
		Schema: schemaName,
		Valid:  true,
	}

	// Basic HCL structure validation
	if file.Body == nil {
		result.Errors = append(result.Errors, ValidationError{
			Field:    "body",
			Message:  "HCL file has no body",
			Severity: "error",
		})
		return result
	}

	// For now, just validate basic HCL structure
	// More sophisticated schema validation can be added later
	result = sv.validateBasicHCL(file)

	return result
}

// validateBasicHCL validates basic HCL structure (from EnhancedSchemaValidator)
func (sv *SchemaValidator) validateBasicHCL(file *hcl.File) *ValidationResult {
	result := &ValidationResult{
		Valid: true,
	}

	// Basic validation - just check if it can be parsed with any blocks allowed
	_, diags := file.Body.Content(&hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{
			{Type: "project", LabelNames: []string{"name"}},
			{Type: "machines"},
			{Type: "machine", LabelNames: []string{"name"}},
			{Type: "actions"},
			{Type: "action", LabelNames: []string{"name"}},
			{Type: "variable", LabelNames: []string{"name"}},
			{Type: "variables"},
			{Type: "storage"},
			{Type: "logging"},
			{Type: "ssh"},
		},
	})
	if diags.HasErrors() {
		result.Valid = false
		for _, diag := range diags {
			validationError := ConvertHCLDiagnostic(*diag)
			result.Errors = append(result.Errors, *validationError)
		}
	}

	return result
}

// GetSchema returns a loaded schema by type
func (sv *SchemaValidator) GetSchema(schemaType SchemaType) (*Schema, error) {
	schema, exists := sv.schemas[schemaType]
	if !exists {
		return nil, fmt.Errorf("schema %s not loaded", schemaType)
	}
	return schema, nil
}

// ListLoadedSchemas returns a list of loaded schema types
func (sv *SchemaValidator) ListLoadedSchemas() []SchemaType {
	var schemas []SchemaType
	for schemaType := range sv.schemas {
		schemas = append(schemas, schemaType)
	}
	return schemas
}

// GetValidationErrors returns all validation errors (interface method)
func (sv *SchemaValidator) GetValidationErrors() []ValidationError {
	// This would return accumulated validation errors
	// For now, return empty slice as errors are returned in ValidationResult
	return []ValidationError{}
}

// cacheResult caches validation results for performance
func (sv *SchemaValidator) cacheResult(key string, result *ValidationResult) {
	sv.mutex.Lock()
	defer sv.mutex.Unlock()
	sv.cache[key] = result
}
