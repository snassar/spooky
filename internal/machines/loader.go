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
func (l *Loader) LoadMachinesFromFile(ctx context.Context, filePath string) ([]spookytypes.Machine, error) {
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
	machineName := block.Labels[0]

	// Parse both attributes and blocks
	content, diags := block.Body.Content(&hcl.BodySchema{
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
	})
	if diags.HasErrors() {
		return nil, fmt.Errorf("failed to parse machine block: %s", diags.Error())
	}

	machine := &spookytypes.Machine{
		Hostname: machineName,
		Port:     22, // Default SSH port
	}

	// Add source file information to machine metadata
	if machine.MachineMetadata == nil {
		machine.MachineMetadata = &spookytypesmachines.MachineMetadata{}
	}
	machine.MachineMetadata.CustomFields = make(map[string]string)
	machine.MachineMetadata.CustomFields["source_file"] = sourceFile

	// Parse attributes
	attrs := content.Attributes

	// Parse basic attributes
	if attr, exists := attrs["hostname"]; exists {
		val, diags := attr.Expr.Value(nil)
		if diags.HasErrors() {
			return nil, fmt.Errorf("invalid hostname: %s", diags.Error())
		}
		if val.Type() == cty.String {
			machine.Hostname = val.AsString()
		}
	}

	if attr, exists := attrs["host"]; exists {
		val, diags := attr.Expr.Value(nil)
		if diags.HasErrors() {
			return nil, fmt.Errorf("invalid host: %s", diags.Error())
		}
		if val.Type() == cty.String {
			machine.Host = val.AsString()
		}
	}

	if attr, exists := attrs["port"]; exists {
		val, diags := attr.Expr.Value(nil)
		if diags.HasErrors() {
			return nil, fmt.Errorf("invalid port: %s", diags.Error())
		}
		if val.Type() == cty.Number {
			portInt, _ := val.AsBigFloat().Int64()
			machine.Port = int(portInt)
		}
	}

	if attr, exists := attrs["user"]; exists {
		val, diags := attr.Expr.Value(nil)
		if diags.HasErrors() {
			return nil, fmt.Errorf("invalid user: %s", diags.Error())
		}
		if val.Type() == cty.String {
			machine.User = val.AsString()
		}
	}

	// Parse SSH authentication
	if attr, exists := attrs["key_file"]; exists {
		val, diags := attr.Expr.Value(nil)
		if diags.HasErrors() {
			return nil, fmt.Errorf("invalid key_file: %s", diags.Error())
		}
		if val.Type() == cty.String {
			machine.KeyFile = val.AsString()
		}
	}

	if attr, exists := attrs["passphrase"]; exists {
		val, diags := attr.Expr.Value(nil)
		if diags.HasErrors() {
			return nil, fmt.Errorf("invalid passphrase: %s", diags.Error())
		}
		if val.Type() == cty.String {
			machine.Passphrase = val.AsString()
		}
	}

	// Parse tags
	if attr, exists := attrs["tags"]; exists {
		tags, err := l.parseObjectAttribute(attr)
		if err != nil {
			return nil, fmt.Errorf("invalid tags: %w", err)
		}
		machine.Tags = tags
	}

	// Parse groups
	if attr, exists := attrs["groups"]; exists {
		groups, err := l.parseArrayAttribute(attr)
		if err != nil {
			return nil, fmt.Errorf("invalid groups: %w", err)
		}
		machine.Groups = groups
	}

	// Parse roles
	if attr, exists := attrs["roles"]; exists {
		roles, err := l.parseArrayAttribute(attr)
		if err != nil {
			return nil, fmt.Errorf("invalid roles: %w", err)
		}
		machine.Roles = roles
	}

	// Parse classes
	if attr, exists := attrs["classes"]; exists {
		classes, err := l.parseArrayAttribute(attr)
		if err != nil {
			return nil, fmt.Errorf("invalid classes: %w", err)
		}
		machine.Classes = classes
	}

	// Parse timeouts
	if attr, exists := attrs["connection_timeout"]; exists {
		val, diags := attr.Expr.Value(nil)
		if diags.HasErrors() {
			return nil, fmt.Errorf("invalid connection_timeout: %s", diags.Error())
		}
		if val.Type() == cty.Number {
			timeoutInt, _ := val.AsBigFloat().Int64()
			machine.ConnectionTimeout = int(timeoutInt)
		}
	}

	if attr, exists := attrs["command_timeout"]; exists {
		val, diags := attr.Expr.Value(nil)
		if diags.HasErrors() {
			return nil, fmt.Errorf("invalid command_timeout: %s", diags.Error())
		}
		if val.Type() == cty.Number {
			timeoutInt, _ := val.AsBigFloat().Int64()
			machine.CommandTimeout = int(timeoutInt)
		}
	}

	// Parse connection settings
	if attr, exists := attrs["max_connections"]; exists {
		val, diags := attr.Expr.Value(nil)
		if diags.HasErrors() {
			return nil, fmt.Errorf("invalid max_connections: %s", diags.Error())
		}
		if val.Type() == cty.Number {
			maxConnInt, _ := val.AsBigFloat().Int64()
			machine.MaxConnections = int(maxConnInt)
		}
	}

	if attr, exists := attrs["retry_attempts"]; exists {
		val, diags := attr.Expr.Value(nil)
		if diags.HasErrors() {
			return nil, fmt.Errorf("invalid retry_attempts: %s", diags.Error())
		}
		if val.Type() == cty.Number {
			retryInt, _ := val.AsBigFloat().Int64()
			machine.RetryAttempts = int(retryInt)
		}
	}

	if attr, exists := attrs["retry_delay"]; exists {
		val, diags := attr.Expr.Value(nil)
		if diags.HasErrors() {
			return nil, fmt.Errorf("invalid retry_delay: %s", diags.Error())
		}
		if val.Type() == cty.Number {
			delayInt, _ := val.AsBigFloat().Int64()
			machine.RetryDelay = int(delayInt)
		}
	}

	// Parse blocks
	for _, block := range content.Blocks {
		switch block.Type {
		case "resources":
			resources, err := l.parseResourcesBlock(block)
			if err != nil {
				return nil, fmt.Errorf("failed to parse resources block: %w", err)
			}
			machine.Resources = resources
		case "metadata":
			metadata, err := l.parseMetadataBlock(block)
			if err != nil {
				return nil, fmt.Errorf("failed to parse metadata block: %w", err)
			}
			machine.MachineMetadata = metadata
		}
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

// parseArrayAttribute parses an array attribute into a []string
func (l *Loader) parseArrayAttribute(attr *hcl.Attribute) ([]string, error) {
	val, diags := attr.Expr.Value(nil)
	if diags.HasErrors() {
		return nil, fmt.Errorf("failed to parse array attribute: %s", diags.Error())
	}

	if val.Type() != cty.List(cty.String) {
		return nil, fmt.Errorf("expected list of strings, got %s", val.Type().FriendlyName())
	}

	var result []string
	for _, item := range val.AsValueSlice() {
		result = append(result, item.AsString())
	}

	return result, nil
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

	attrs := content.Attributes

	if attr, exists := attrs["cpu_cores"]; exists {
		val, diags := attr.Expr.Value(nil)
		if diags.HasErrors() {
			return nil, fmt.Errorf("invalid cpu_cores: %s", diags.Error())
		}
		if val.Type() == cty.Number {
			cpuInt, _ := val.AsBigFloat().Int64()
			resources.CPUCores = int(cpuInt)
		}
	}

	if attr, exists := attrs["memory_gb"]; exists {
		val, diags := attr.Expr.Value(nil)
		if diags.HasErrors() {
			return nil, fmt.Errorf("invalid memory_gb: %s", diags.Error())
		}
		if val.Type() == cty.Number {
			memInt, _ := val.AsBigFloat().Int64()
			resources.MemoryGB = int(memInt)
		}
	}

	if attr, exists := attrs["disk_gb"]; exists {
		val, diags := attr.Expr.Value(nil)
		if diags.HasErrors() {
			return nil, fmt.Errorf("invalid disk_gb: %s", diags.Error())
		}
		if val.Type() == cty.Number {
			diskInt, _ := val.AsBigFloat().Int64()
			resources.DiskGB = int(diskInt)
		}
	}

	if attr, exists := attrs["network_speed"]; exists {
		val, diags := attr.Expr.Value(nil)
		if diags.HasErrors() {
			return nil, fmt.Errorf("invalid network_speed: %s", diags.Error())
		}
		if val.Type() == cty.String {
			resources.NetworkSpeed = val.AsString()
		}
	}

	return resources, nil
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

	attrs := content.Attributes

	if attr, exists := attrs["environment"]; exists {
		val, diags := attr.Expr.Value(nil)
		if diags.HasErrors() {
			return nil, fmt.Errorf("invalid environment: %s", diags.Error())
		}
		if val.Type() == cty.String {
			metadata.Environment = val.AsString()
		}
	}

	if attr, exists := attrs["location"]; exists {
		val, diags := attr.Expr.Value(nil)
		if diags.HasErrors() {
			return nil, fmt.Errorf("invalid location: %s", diags.Error())
		}
		if val.Type() == cty.String {
			metadata.Location = val.AsString()
		}
	}

	if attr, exists := attrs["owner"]; exists {
		val, diags := attr.Expr.Value(nil)
		if diags.HasErrors() {
			return nil, fmt.Errorf("invalid owner: %s", diags.Error())
		}
		if val.Type() == cty.String {
			metadata.Owner = val.AsString()
		}
	}

	if attr, exists := attrs["custom_fields"]; exists {
		customFields, err := l.parseObjectAttribute(attr)
		if err != nil {
			return nil, fmt.Errorf("invalid custom_fields: %w", err)
		}
		metadata.CustomFields = customFields
	}

	return metadata, nil
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
