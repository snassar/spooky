// Package machines provides machine loading functionality for the spooky codebase.
package machines

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/zclconf/go-cty/cty"

	spookytypes "spooky/internal/types"
	spookytypeslogging "spooky/internal/types/logging"
	spookytypesmachines "spooky/internal/types/machines"
)

// Loader provides machine loading functionality
type Loader struct {
	logger spookytypeslogging.Logger
}

// NewLoader creates a new machine loader
func NewLoader(logger spookytypeslogging.Logger) *Loader {
	return &Loader{
		logger: logger,
	}
}

// LoadMachinesFromFile loads machines from a single HCL file
func (l *Loader) LoadMachinesFromFile(ctx context.Context, filePath string) ([]spookytypes.Machine, error) {
	l.logger.Debug("Loading machines from file", map[string]interface{}{
		"file_path": filePath,
	})

	// Read the file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	// Parse HCL content
	parser := hclparse.NewParser()
	file, diags := parser.ParseHCL(data, filePath)
	if diags.HasErrors() {
		return nil, fmt.Errorf("failed to parse HCL file %s: %w", filePath, diags.Error())
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
		return nil, fmt.Errorf("failed to parse machines block in %s: %w", filePath, diags.Error())
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
		return nil, fmt.Errorf("failed to parse machine blocks: %w", diags.Error())
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
		return nil, fmt.Errorf("failed to parse machine block: %w", diags.Error())
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
			return nil, fmt.Errorf("invalid hostname: %w", diags.Error())
		}
		if val.Type() == cty.String {
			machine.Hostname = val.AsString()
		}
	}

	if attr, exists := attrs["host"]; exists {
		val, diags := attr.Expr.Value(nil)
		if diags.HasErrors() {
			return nil, fmt.Errorf("invalid host: %w", diags.Error())
		}
		if val.Type() == cty.String {
			machine.Host = val.AsString()
		}
	}

	if attr, exists := attrs["port"]; exists {
		val, diags := attr.Expr.Value(nil)
		if diags.HasErrors() {
			return nil, fmt.Errorf("invalid port: %w", diags.Error())
		}
		if val.Type() == cty.Number {
			portInt, _ := val.AsBigFloat().Int64()
			machine.Port = int(portInt)
		}
	}

	if attr, exists := attrs["user"]; exists {
		val, diags := attr.Expr.Value(nil)
		if diags.HasErrors() {
			return nil, fmt.Errorf("invalid user: %w", diags.Error())
		}
		if val.Type() == cty.String {
			machine.User = val.AsString()
		}
	}

	// Parse SSH authentication
	if attr, exists := attrs["key_file"]; exists {
		val, diags := attr.Expr.Value(nil)
		if diags.HasErrors() {
			return nil, fmt.Errorf("invalid key_file: %w", diags.Error())
		}
		if val.Type() == cty.String {
			machine.KeyFile = val.AsString()
		}
	}

	if attr, exists := attrs["passphrase"]; exists {
		val, diags := attr.Expr.Value(nil)
		if diags.HasErrors() {
			return nil, fmt.Errorf("invalid passphrase: %w", diags.Error())
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

	// Parse SSH connection configuration
	if attr, exists := attrs["connection_timeout"]; exists {
		val, diags := attr.Expr.Value(nil)
		if diags.HasErrors() {
			return nil, fmt.Errorf("invalid connection_timeout: %w", diags.Error())
		}
		if val.Type() == cty.Number {
			timeoutInt, _ := val.AsBigFloat().Int64()
			machine.ConnectionTimeout = int(timeoutInt)
		}
	}

	if attr, exists := attrs["command_timeout"]; exists {
		val, diags := attr.Expr.Value(nil)
		if diags.HasErrors() {
			return nil, fmt.Errorf("invalid command_timeout: %w", diags.Error())
		}
		if val.Type() == cty.Number {
			timeoutInt, _ := val.AsBigFloat().Int64()
			machine.CommandTimeout = int(timeoutInt)
		}
	}

	if attr, exists := attrs["max_connections"]; exists {
		val, diags := attr.Expr.Value(nil)
		if diags.HasErrors() {
			return nil, fmt.Errorf("invalid max_connections: %w", diags.Error())
		}
		if val.Type() == cty.Number {
			maxInt, _ := val.AsBigFloat().Int64()
			machine.MaxConnections = int(maxInt)
		}
	}

	if attr, exists := attrs["retry_attempts"]; exists {
		val, diags := attr.Expr.Value(nil)
		if diags.HasErrors() {
			return nil, fmt.Errorf("invalid retry_attempts: %w", diags.Error())
		}
		if val.Type() == cty.Number {
			retryInt, _ := val.AsBigFloat().Int64()
			machine.RetryAttempts = int(retryInt)
		}
	}

	if attr, exists := attrs["retry_delay"]; exists {
		val, diags := attr.Expr.Value(nil)
		if diags.HasErrors() {
			return nil, fmt.Errorf("invalid retry_delay: %w", diags.Error())
		}
		if val.Type() == cty.Number {
			delayInt, _ := val.AsBigFloat().Int64()
			machine.RetryDelay = int(delayInt)
		}
	}

	// Parse resources block
	for _, resourcesBlock := range content.Blocks {
		if resourcesBlock.Type == "resources" {
			resources, err := l.parseResourcesBlock(resourcesBlock)
			if err != nil {
				return nil, fmt.Errorf("invalid resources block: %w", err)
			}
			machine.Resources = resources
		}
	}

	// Parse metadata block
	for _, metadataBlock := range content.Blocks {
		if metadataBlock.Type == "metadata" {
			metadata, err := l.parseMetadataBlock(metadataBlock)
			if err != nil {
				return nil, fmt.Errorf("invalid metadata block: %w", err)
			}
			// Merge with existing metadata
			if machine.MachineMetadata == nil {
				machine.MachineMetadata = metadata
			} else {
				// Merge custom fields
				if metadata.CustomFields != nil {
					if machine.MachineMetadata.CustomFields == nil {
						machine.MachineMetadata.CustomFields = make(map[string]string)
					}
					for k, v := range metadata.CustomFields {
						machine.MachineMetadata.CustomFields[k] = v
					}
				}
				// Merge other fields (only if not already set)
				if machine.MachineMetadata.Environment == "" && metadata.Environment != "" {
					machine.MachineMetadata.Environment = metadata.Environment
				}
				if machine.MachineMetadata.Datacenter == "" && metadata.Datacenter != "" {
					machine.MachineMetadata.Datacenter = metadata.Datacenter
				}
				if machine.MachineMetadata.Rack == "" && metadata.Rack != "" {
					machine.MachineMetadata.Rack = metadata.Rack
				}
				if machine.MachineMetadata.Location == "" && metadata.Location != "" {
					machine.MachineMetadata.Location = metadata.Location
				}
				if machine.MachineMetadata.Owner == "" && metadata.Owner != "" {
					machine.MachineMetadata.Owner = metadata.Owner
				}
				if machine.MachineMetadata.Department == "" && metadata.Department != "" {
					machine.MachineMetadata.Department = metadata.Department
				}
				if machine.MachineMetadata.CostCenter == "" && metadata.CostCenter != "" {
					machine.MachineMetadata.CostCenter = metadata.CostCenter
				}
			}
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

	if !val.Type().IsObjectType() {
		return nil, fmt.Errorf("expected object, got %s", val.Type().FriendlyName())
	}

	result := make(map[string]string)
	for key, value := range val.AsValueMap() {
		if value.Type() == cty.String {
			result[key] = value.AsString()
		}
	}

	return result, nil
}

// parseArrayAttribute parses an array attribute into a []string
func (l *Loader) parseArrayAttribute(attr *hcl.Attribute) ([]string, error) {
	val, diags := attr.Expr.Value(nil)
	if diags.HasErrors() {
		return nil, fmt.Errorf("failed to parse array attribute: %s", diags.Error())
	}

	// Handle both list and tuple types (HCL can represent arrays as either)
	if !val.Type().IsListType() && !val.Type().IsTupleType() {
		return nil, fmt.Errorf("expected list or tuple, got %s", val.Type().FriendlyName())
	}

	var result []string
	for _, item := range val.AsValueSlice() {
		if item.Type() == cty.String {
			result = append(result, item.AsString())
		}
	}

	return result, nil
}

// parseResourcesBlock parses a resources block
func (l *Loader) parseResourcesBlock(block *hcl.Block) (*spookytypesmachines.MachineResources, error) {
	attrs, diags := block.Body.JustAttributes()
	if diags.HasErrors() {
		return nil, fmt.Errorf("failed to parse resources attributes: %w", diags.Error())
	}

	resources := &spookytypesmachines.MachineResources{}

	if attr, exists := attrs["cpu_cores"]; exists {
		val, diags := attr.Expr.Value(nil)
		if diags.HasErrors() {
			return nil, fmt.Errorf("invalid cpu_cores: %w", diags.Error())
		}
		if val.Type() == cty.Number {
			cpuInt, _ := val.AsBigFloat().Int64()
			resources.CPUCores = int(cpuInt)
		}
	}

	if attr, exists := attrs["memory_gb"]; exists {
		val, diags := attr.Expr.Value(nil)
		if diags.HasErrors() {
			return nil, fmt.Errorf("invalid memory_gb: %w", diags.Error())
		}
		if val.Type() == cty.Number {
			memInt, _ := val.AsBigFloat().Int64()
			resources.MemoryGB = int(memInt)
		}
	}

	if attr, exists := attrs["disk_gb"]; exists {
		val, diags := attr.Expr.Value(nil)
		if diags.HasErrors() {
			return nil, fmt.Errorf("invalid disk_gb: %w", diags.Error())
		}
		if val.Type() == cty.Number {
			diskInt, _ := val.AsBigFloat().Int64()
			resources.DiskGB = int(diskInt)
		}
	}

	if attr, exists := attrs["network_speed"]; exists {
		val, diags := attr.Expr.Value(nil)
		if diags.HasErrors() {
			return nil, fmt.Errorf("invalid network_speed: %w", diags.Error())
		}
		if val.Type() == cty.String {
			resources.NetworkSpeed = val.AsString()
		}
	}

	return resources, nil
}

// parseMetadataBlock parses a metadata block
func (l *Loader) parseMetadataBlock(block *hcl.Block) (*spookytypesmachines.MachineMetadata, error) {
	attrs, diags := block.Body.JustAttributes()
	if diags.HasErrors() {
		return nil, fmt.Errorf("failed to parse metadata attributes: %w", diags.Error())
	}

	metadata := &spookytypesmachines.MachineMetadata{}

	// Parse string attributes
	stringFields := []string{"environment", "datacenter", "rack", "location", "owner", "department", "cost_center"}
	for _, field := range stringFields {
		if attr, exists := attrs[field]; exists {
			val, diags := attr.Expr.Value(nil)
			if diags.HasErrors() {
				return nil, fmt.Errorf("invalid %s: %w", field, diags.Error())
			}
			if val.Type() == cty.String {
				switch field {
				case "environment":
					metadata.Environment = val.AsString()
				case "datacenter":
					metadata.Datacenter = val.AsString()
				case "rack":
					metadata.Rack = val.AsString()
				case "location":
					metadata.Location = val.AsString()
				case "owner":
					metadata.Owner = val.AsString()
				case "department":
					metadata.Department = val.AsString()
				case "cost_center":
					metadata.CostCenter = val.AsString()
				}
			}
		}
	}

	// Parse custom fields (any remaining attributes)
	metadata.CustomFields = make(map[string]string)
	for name, attr := range attrs {
		// Skip already parsed fields
		found := false
		for _, field := range stringFields {
			if name == field {
				found = true
				break
			}
		}
		if found {
			continue
		}

		val, diags := attr.Expr.Value(nil)
		if diags.HasErrors() {
			continue // Skip invalid attributes
		}
		if val.Type() == cty.String {
			metadata.CustomFields[name] = val.AsString()
		}
	}

	return metadata, nil
}

// LoadMachinesFromDirectory loads machines from a directory containing HCL files
func (l *Loader) LoadMachinesFromDirectory(ctx context.Context, dirPath string) ([]spookytypes.Machine, error) {
	l.logger.Debug("Loading machines from directory", map[string]interface{}{
		"dir_path": dirPath,
	})

	// Read directory entries
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory %s: %w", dirPath, err)
	}

	var allMachines []spookytypes.Machine

	// Process each HCL file in the directory
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		if filepath.Ext(entry.Name()) == ".hcl" {
			filePath := filepath.Join(dirPath, entry.Name())
			machines, err := l.LoadMachinesFromFile(ctx, filePath)
			if err != nil {
				l.logger.Warn("Failed to load machines from file", map[string]interface{}{
					"file_path": filePath,
					"error":     err.Error(),
				})
				continue
			}
			allMachines = append(allMachines, machines...)
		}
	}

	l.logger.Info("Machines loaded from directory", map[string]interface{}{
		"dir_path": dirPath,
		"count":    len(allMachines),
	})

	return allMachines, nil
}
