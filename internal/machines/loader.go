// Package machines provides machine inventory loading and management functionality.
package machines

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/zclconf/go-cty/cty"

	spookytypes "spooky/internal/types"
	spookytypeslogging "spooky/internal/types/logging"
	spookytypesmachines "spooky/internal/types/machines"
)

// AttributeParser defines the interface for parsing individual attributes
type AttributeParser interface {
	Parse(attr *hcl.Attribute, machine *spookytypes.Machine) error
	GetFieldName() string
}

// AttributeParserRegistry manages all attribute parsers
type AttributeParserRegistry struct {
	parsers map[string]AttributeParser
}

// NewAttributeParserRegistry creates a new attribute parser registry
func NewAttributeParserRegistry() *AttributeParserRegistry {
	return &AttributeParserRegistry{
		parsers: map[string]AttributeParser{
			"hostname":           &StringAttributeParser{fieldName: "Host"},
			"host":               &StringAttributeParser{fieldName: "Host"},
			"port":               &IntAttributeParser{fieldName: "Port"},
			"user":               &StringAttributeParser{fieldName: "User"},
			"key_file":           &StringAttributeParser{fieldName: "KeyFile"},
			"passphrase":         &StringAttributeParser{fieldName: "Passphrase"},
			"tags":               &ObjectAttributeParser{fieldName: "Tags"},
			"groups":             &ArrayAttributeParser{fieldName: "Groups"},
			"roles":              &ArrayAttributeParser{fieldName: "Roles"},
			"classes":            &ArrayAttributeParser{fieldName: "Classes"},
			"connection_timeout": &IntAttributeParser{fieldName: "ConnectionTimeout"},
			"command_timeout":    &IntAttributeParser{fieldName: "CommandTimeout"},
			"max_connections":    &IntAttributeParser{fieldName: "MaxConnections"},
			"retry_attempts":     &IntAttributeParser{fieldName: "RetryAttempts"},
			"retry_delay":        &IntAttributeParser{fieldName: "RetryDelay"},
		},
	}
}

// StringAttributeParser handles string attributes
type StringAttributeParser struct {
	fieldName string
}

func (s *StringAttributeParser) Parse(attr *hcl.Attribute, machine *spookytypes.Machine) error {
	val, diags := attr.Expr.Value(nil)
	if diags.HasErrors() {
		return fmt.Errorf("invalid %s: %s", s.fieldName, diags.Error())
	}
	if val.Type() == cty.String {
		return s.setField(machine, val.AsString())
	}
	return fmt.Errorf("expected string for %s, got %s", s.fieldName, val.Type().FriendlyName())
}

func (s *StringAttributeParser) setField(machine *spookytypes.Machine, value string) error {
	switch s.fieldName {
	case "Hostname":
		machine.Hostname = value
	case "Host":
		machine.Host = value
	case "User":
		machine.User = value
	case "KeyFile":
		machine.KeyFile = value
	case "Passphrase":
		machine.Passphrase = value
	default:
		return fmt.Errorf("unknown string field: %s", s.fieldName)
	}
	return nil
}

func (s *StringAttributeParser) GetFieldName() string {
	return s.fieldName
}

// IntAttributeParser handles integer attributes
type IntAttributeParser struct {
	fieldName string
}

func (i *IntAttributeParser) Parse(attr *hcl.Attribute, machine *spookytypes.Machine) error {
	val, diags := attr.Expr.Value(nil)
	if diags.HasErrors() {
		return fmt.Errorf("invalid %s: %s", i.fieldName, diags.Error())
	}
	if val.Type() == cty.Number {
		intVal, _ := val.AsBigFloat().Int64()
		return i.setField(machine, int(intVal))
	}
	return fmt.Errorf("expected number for %s, got %s", i.fieldName, val.Type().FriendlyName())
}

func (i *IntAttributeParser) setField(machine *spookytypes.Machine, value int) error {
	switch i.fieldName {
	case "Port":
		machine.Port = value
	case "ConnectionTimeout":
		machine.ConnectionTimeout = value
	case "CommandTimeout":
		machine.CommandTimeout = value
	case "MaxConnections":
		machine.MaxConnections = value
	case "RetryAttempts":
		machine.RetryAttempts = value
	case "RetryDelay":
		machine.RetryDelay = value
	default:
		return fmt.Errorf("unknown int field: %s", i.fieldName)
	}
	return nil
}

func (i *IntAttributeParser) GetFieldName() string {
	return i.fieldName
}

// ObjectAttributeParser handles object attributes (maps)
type ObjectAttributeParser struct {
	fieldName string
}

func (o *ObjectAttributeParser) Parse(attr *hcl.Attribute, machine *spookytypes.Machine) error {
	val, diags := attr.Expr.Value(nil)
	if diags.HasErrors() {
		return fmt.Errorf("failed to parse %s: %s", o.fieldName, diags.Error())
	}

	if val.Type() != cty.Map(cty.String) {
		return fmt.Errorf("expected map of strings for %s, got %s", o.fieldName, val.Type().FriendlyName())
	}

	result := make(map[string]string)
	val.ForEachElement(func(key, value cty.Value) bool {
		result[key.AsString()] = value.AsString()
		return false
	})

	return o.setField(machine, result)
}

func (o *ObjectAttributeParser) setField(machine *spookytypes.Machine, value map[string]string) error {
	switch o.fieldName {
	case "Tags":
		machine.Tags = value
	default:
		return fmt.Errorf("unknown object field: %s", o.fieldName)
	}
	return nil
}

func (o *ObjectAttributeParser) GetFieldName() string {
	return o.fieldName
}

// ArrayAttributeParser handles array attributes (lists)
type ArrayAttributeParser struct {
	fieldName string
}

func (a *ArrayAttributeParser) Parse(attr *hcl.Attribute, machine *spookytypes.Machine) error {
	val, diags := attr.Expr.Value(nil)
	if diags.HasErrors() {
		return fmt.Errorf("failed to parse %s: %s", a.fieldName, diags.Error())
	}

	if val.Type() != cty.List(cty.String) {
		return fmt.Errorf("expected list of strings for %s, got %s", a.fieldName, val.Type().FriendlyName())
	}

	var result []string
	for _, item := range val.AsValueSlice() {
		result = append(result, item.AsString())
	}

	return a.setField(machine, result)
}

func (a *ArrayAttributeParser) setField(machine *spookytypes.Machine, value []string) error {
	switch a.fieldName {
	case "Groups":
		machine.Groups = value
	case "Roles":
		machine.Roles = value
	case "Classes":
		machine.Classes = value
	default:
		return fmt.Errorf("unknown array field: %s", a.fieldName)
	}
	return nil
}

func (a *ArrayAttributeParser) GetFieldName() string {
	return a.fieldName
}

// BlockParser defines the interface for parsing blocks
type BlockParser interface {
	Parse(block *hcl.Block, machine *spookytypes.Machine) error
	GetBlockType() string
}

// BlockParserRegistry manages all block parsers
type BlockParserRegistry struct {
	parsers map[string]BlockParser
}

// NewBlockParserRegistry creates a new block parser registry
func NewBlockParserRegistry(loader *Loader) *BlockParserRegistry {
	return &BlockParserRegistry{
		parsers: map[string]BlockParser{
			"resources": &ResourcesBlockParser{loader: loader},
			"metadata":  &MetadataBlockParser{loader: loader},
		},
	}
}

// ResourcesBlockParser handles resources blocks
type ResourcesBlockParser struct {
	loader *Loader
}

func (r *ResourcesBlockParser) Parse(block *hcl.Block, machine *spookytypes.Machine) error {
	resources, err := r.loader.parseResourcesBlock(block)
	if err != nil {
		return fmt.Errorf("failed to parse resources block: %w", err)
	}
	machine.Resources = resources
	return nil
}

func (r *ResourcesBlockParser) GetBlockType() string {
	return "resources"
}

// MetadataBlockParser handles metadata blocks
type MetadataBlockParser struct {
	loader *Loader
}

func (m *MetadataBlockParser) Parse(block *hcl.Block, machine *spookytypes.Machine) error {
	metadata, err := m.loader.parseMetadataBlock(block)
	if err != nil {
		return fmt.Errorf("failed to parse metadata block: %w", err)
	}
	machine.MachineMetadata = metadata
	return nil
}

func (m *MetadataBlockParser) GetBlockType() string {
	return "metadata"
}

// Loader provides functionality to load machine inventory from HCL files
type Loader struct {
	logger spookytypeslogging.Logger
}

// NewLoader creates a new machine loader instance
func NewLoader(logger spookytypeslogging.Logger) *Loader {
	return &Loader{
		logger: logger,
	}
}

// LoadMachinesFromFile loads machine inventory from a single HCL file
func (l *Loader) LoadMachinesFromFile(_ context.Context, filePath string) ([]spookytypes.Machine, error) {
	// Read file content
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	// Parse HCL content
	parser := hclparse.NewParser()
	file, diags := parser.ParseHCL(data, filePath)
	if diags.HasErrors() {
		return nil, fmt.Errorf("failed to parse HCL file %s: %s", filePath, diags.Error())
	}

	// Extract machines block
	content, diags := file.Body.Content(&hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{
			{
				Type: "machines",
			},
		},
	})
	if diags.HasErrors() {
		return nil, fmt.Errorf("failed to parse machines block in %s: %s", filePath, diags.Error())
	}

	if len(content.Blocks) == 0 {
		return nil, fmt.Errorf("no machines block found in %s", filePath)
	}

	machinesBlock := content.Blocks[0]
	machines, err := l.parseMachinesBlock(machinesBlock, filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse machines block: %w", err)
	}

	l.logger.Info("Machines loaded successfully", map[string]interface{}{
		"source":   filePath,
		"count":    len(machines),
		"machines": getMachineHostnames(machines),
	})

	return machines, nil
}

// parseMachinesBlock parses the machines block and extracts machine definitions
func (l *Loader) parseMachinesBlock(block *hcl.Block, sourceFile string) ([]spookytypes.Machine, error) {
	var machines []spookytypes.Machine

	// Parse machine blocks within the machines block
	content, diags := block.Body.Content(&hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{
			{
				Type:       "machine",
				LabelNames: []string{"name"},
			},
		},
	})
	if diags.HasErrors() {
		return nil, fmt.Errorf("failed to parse machine blocks: %s", diags.Error())
	}

	for _, machineBlock := range content.Blocks {
		machine, err := l.parseMachineBlock(machineBlock, sourceFile)
		if err != nil {
			l.logger.Warn("Failed to parse machine block", map[string]interface{}{
				"machine_name": machineBlock.Labels[0],
				"source_file":  sourceFile,
				"error":        err.Error(),
			})
			continue
		}
		machines = append(machines, *machine)
	}

	return machines, nil
}

// parseMachineBlock parses a single machine block
func (l *Loader) parseMachineBlock(block *hcl.Block, sourceFile string) (*spookytypes.Machine, error) {
	// Step 1: Parse HCL content
	content, err := l.parseBlockContent(block)
	if err != nil {
		return nil, err
	}

	// Step 2: Create base machine
	machine := l.createBaseMachine(block.Labels[0], sourceFile)

	// Step 3: Parse attributes using registry
	if err := l.parseAttributes(content.Attributes, machine); err != nil {
		return nil, err
	}

	// Step 4: Parse blocks using registry
	if err := l.parseBlocks(content.Blocks, machine); err != nil {
		return nil, err
	}

	return machine, nil
}

// parseObjectAttribute parses an object attribute into a map[string]string
func (l *Loader) parseObjectAttribute(attr *hcl.Attribute) (map[string]string, error) {
	val, diags := attr.Expr.Value(nil)
	if diags.HasErrors() {
		return nil, fmt.Errorf("failed to parse object attribute: %s", diags.Error())
	}

	if val.Type() != cty.Map(cty.String) {
		return nil, fmt.Errorf("expected map of strings, got %s", val.Type().FriendlyName())
	}

	result := make(map[string]string)
	val.ForEachElement(func(key, value cty.Value) bool {
		result[key.AsString()] = value.AsString()
		return false
	})

	return result, nil
}

// getMachineHostnames extracts hostnames from a slice of machines for logging
func getMachineHostnames(machines []spookytypes.Machine) []string {
	hostnames := make([]string, len(machines))
	for i, machine := range machines {
		hostnames[i] = machine.Hostname
	}
	return hostnames
}

// parseResourcesBlock parses a resources block
func (l *Loader) parseResourcesBlock(block *hcl.Block) (*spookytypesmachines.MachineResources, error) {
	content, diags := block.Body.Content(&hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{
			{Name: "cpu_cores", Required: false},
			{Name: "memory_mb", Required: false},
			{Name: "disk_gb", Required: false},
			{Name: "network_interfaces", Required: false},
		},
	})
	if diags.HasErrors() {
		return nil, fmt.Errorf("failed to parse resources block: %s", diags.Error())
	}

	resources := &spookytypesmachines.MachineResources{}

	if err := l.parseResourceAttributes(content.Attributes, resources); err != nil {
		return nil, err
	}

	return resources, nil
}

func (l *Loader) parseResourceAttributes(attrs hcl.Attributes, resources *spookytypesmachines.MachineResources) error {
	parsers := map[string]func(*hcl.Attribute, *spookytypesmachines.MachineResources) error{
		"cpu_cores":     l.parseCPUCores,
		"memory_gb":     l.parseMemoryGB,
		"disk_gb":       l.parseDiskGB,
		"network_speed": l.parseNetworkSpeed,
	}

	for name, parser := range parsers {
		if attr, exists := attrs[name]; exists {
			if err := parser(attr, resources); err != nil {
				return fmt.Errorf("failed to parse %s: %w", name, err)
			}
		}
	}

	return nil
}

func (l *Loader) parseCPUCores(attr *hcl.Attribute, resources *spookytypesmachines.MachineResources) error {
	val, diags := attr.Expr.Value(nil)
	if diags.HasErrors() {
		return fmt.Errorf("invalid cpu_cores: %s", diags.Error())
	}
	if val.Type() == cty.Number {
		cpuInt, _ := val.AsBigFloat().Int64()
		resources.CPUCores = int(cpuInt)
	}
	return nil
}

func (l *Loader) parseMemoryGB(attr *hcl.Attribute, resources *spookytypesmachines.MachineResources) error {
	val, diags := attr.Expr.Value(nil)
	if diags.HasErrors() {
		return fmt.Errorf("invalid memory_gb: %s", diags.Error())
	}
	if val.Type() == cty.Number {
		memInt, _ := val.AsBigFloat().Int64()
		resources.MemoryGB = int(memInt)
	}
	return nil
}

func (l *Loader) parseDiskGB(attr *hcl.Attribute, resources *spookytypesmachines.MachineResources) error {
	val, diags := attr.Expr.Value(nil)
	if diags.HasErrors() {
		return fmt.Errorf("invalid disk_gb: %s", diags.Error())
	}
	if val.Type() == cty.Number {
		diskInt, _ := val.AsBigFloat().Int64()
		resources.DiskGB = int(diskInt)
	}
	return nil
}

func (l *Loader) parseNetworkSpeed(attr *hcl.Attribute, resources *spookytypesmachines.MachineResources) error {
	val, diags := attr.Expr.Value(nil)
	if diags.HasErrors() {
		return fmt.Errorf("invalid network_speed: %s", diags.Error())
	}
	if val.Type() == cty.String {
		resources.NetworkSpeed = val.AsString()
	}
	return nil
}

// parseMetadataBlock parses a metadata block
func (l *Loader) parseMetadataBlock(block *hcl.Block) (*spookytypesmachines.MachineMetadata, error) {
	content, diags := block.Body.Content(&hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{
			{Name: "description", Required: false},
			{Name: "environment", Required: false},
			{Name: "location", Required: false},
			{Name: "owner", Required: false},
			{Name: "team", Required: false},
			{Name: "custom_fields", Required: false},
		},
	})
	if diags.HasErrors() {
		return nil, fmt.Errorf("failed to parse metadata block: %s", diags.Error())
	}

	metadata := &spookytypesmachines.MachineMetadata{
		CustomFields: make(map[string]string),
	}

	if err := l.parseMetadataAttributes(content.Attributes, metadata); err != nil {
		return nil, err
	}

	return metadata, nil
}

func (l *Loader) parseMetadataAttributes(attrs hcl.Attributes, metadata *spookytypesmachines.MachineMetadata) error {
	parsers := map[string]func(*hcl.Attribute, *spookytypesmachines.MachineMetadata) error{
		"environment":   l.parseMetadataEnvironment,
		"location":      l.parseMetadataLocation,
		"owner":         l.parseMetadataOwner,
		"custom_fields": l.parseMetadataCustomFields,
	}

	for name, parser := range parsers {
		if attr, exists := attrs[name]; exists {
			if err := parser(attr, metadata); err != nil {
				return fmt.Errorf("failed to parse %s: %w", name, err)
			}
		}
	}

	return nil
}

func (l *Loader) parseMetadataEnvironment(attr *hcl.Attribute, metadata *spookytypesmachines.MachineMetadata) error {
	val, diags := attr.Expr.Value(nil)
	if diags.HasErrors() {
		return fmt.Errorf("invalid environment: %s", diags.Error())
	}
	if val.Type() == cty.String {
		metadata.Environment = val.AsString()
	}
	return nil
}

func (l *Loader) parseMetadataLocation(attr *hcl.Attribute, metadata *spookytypesmachines.MachineMetadata) error {
	val, diags := attr.Expr.Value(nil)
	if diags.HasErrors() {
		return fmt.Errorf("invalid location: %s", diags.Error())
	}
	if val.Type() == cty.String {
		metadata.Location = val.AsString()
	}
	return nil
}

func (l *Loader) parseMetadataOwner(attr *hcl.Attribute, metadata *spookytypesmachines.MachineMetadata) error {
	val, diags := attr.Expr.Value(nil)
	if diags.HasErrors() {
		return fmt.Errorf("invalid owner: %s", diags.Error())
	}
	if val.Type() == cty.String {
		metadata.Owner = val.AsString()
	}
	return nil
}

func (l *Loader) parseMetadataCustomFields(attr *hcl.Attribute, metadata *spookytypesmachines.MachineMetadata) error {
	customFields, err := l.parseObjectAttribute(attr)
	if err != nil {
		return fmt.Errorf("invalid custom_fields: %w", err)
	}
	metadata.CustomFields = customFields
	return nil
}

// getMachineBlockSchema returns the schema for machine blocks
func (l *Loader) getMachineBlockSchema() *hcl.BodySchema {
	return &hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{
			{Name: "hostname", Required: false},
			{Name: "host", Required: false},
			{Name: "port", Required: false},
			{Name: "user", Required: false},
			{Name: "key_file", Required: false},
			{Name: "passphrase", Required: false},
			{Name: "tags", Required: false},
			{Name: "groups", Required: false},
			{Name: "roles", Required: false},
			{Name: "classes", Required: false},
			{Name: "connection_timeout", Required: false},
			{Name: "command_timeout", Required: false},
			{Name: "max_connections", Required: false},
			{Name: "retry_attempts", Required: false},
			{Name: "retry_delay", Required: false},
		},
		Blocks: []hcl.BlockHeaderSchema{
			{Type: "resources"},
			{Type: "metadata"},
		},
	}
}

// parseBlockContent parses the HCL content of a block
func (l *Loader) parseBlockContent(block *hcl.Block) (*hcl.BodyContent, error) {
	content, diags := block.Body.Content(l.getMachineBlockSchema())
	if diags.HasErrors() {
		return nil, fmt.Errorf("failed to parse machine block: %s", diags.Error())
	}
	return content, nil
}

// createBaseMachine creates a base machine with default values
func (l *Loader) createBaseMachine(machineName, sourceFile string) *spookytypes.Machine {
	machine := &spookytypes.Machine{
		Hostname: machineName,
		Port:     22, // Default SSH port
	}

	// Initialize metadata
	if machine.MachineMetadata == nil {
		machine.MachineMetadata = &spookytypesmachines.MachineMetadata{}
	}
	machine.MachineMetadata.CustomFields = make(map[string]string)
	machine.MachineMetadata.CustomFields["source_file"] = sourceFile

	return machine
}

// parseAttributes parses all attributes using the attribute parser registry
func (l *Loader) parseAttributes(attrs hcl.Attributes, machine *spookytypes.Machine) error {
	registry := NewAttributeParserRegistry()

	for attrName, attr := range attrs {
		parser, exists := registry.parsers[attrName]
		if !exists {
			return fmt.Errorf("unknown attribute: %s", attrName)
		}

		if err := parser.Parse(attr, machine); err != nil {
			return fmt.Errorf("failed to parse %s: %w", attrName, err)
		}
	}

	return nil
}

// parseBlocks parses all blocks using the block parser registry
func (l *Loader) parseBlocks(blocks hcl.Blocks, machine *spookytypes.Machine) error {
	registry := NewBlockParserRegistry(l)

	for _, block := range blocks {
		parser, exists := registry.parsers[block.Type]
		if !exists {
			return fmt.Errorf("unknown block type: %s", block.Type)
		}

		if err := parser.Parse(block, machine); err != nil {
			return fmt.Errorf("failed to parse %s block: %w", block.Type, err)
		}
	}

	return nil
}

// LoadMachinesFromDirectory loads machine inventory from all HCL files in a directory
func (l *Loader) LoadMachinesFromDirectory(ctx context.Context, dirPath string) ([]spookytypes.Machine, error) {
	var allMachines []spookytypes.Machine

	// Walk through directory looking for HCL files
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Only process HCL files
		if !strings.HasSuffix(info.Name(), ".hcl") {
			return nil
		}

		// Load machines from this file
		machines, err := l.LoadMachinesFromFile(ctx, path)
		if err != nil {
			l.logger.Warn("Failed to load machines from file", map[string]interface{}{
				"file":  path,
				"error": err.Error(),
			})
			return nil // Continue with other files
		}

		allMachines = append(allMachines, machines...)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk directory %s: %w", dirPath, err)
	}

	l.logger.Info("Machines loaded from directory", map[string]interface{}{
		"directory": dirPath,
		"count":     len(allMachines),
	})

	return allMachines, nil
}
