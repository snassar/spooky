// Package facts provides fact collection and in-memory management functionality.
package facts

import (
	"fmt"

	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/hcl/v2/hclsyntax"

	spookytypesfacts "spooky/internal/types/facts"
)

// HCLParser provides functionality to parse HCL fact files
type HCLParser struct{}

// NewHCLParser creates a new HCL parser
func NewHCLParser() *HCLParser {
	return &HCLParser{}
}

// ParseCollectorFacts parses collector facts from HCL content
func (p *HCLParser) ParseCollectorFacts(content string) (*spookytypesfacts.CollectorFacts, error) {
	parser := hclparse.NewParser()

	// Parse the HCL content
	file, diags := parser.ParseHCL([]byte(content), "collector-facts.hcl")
	if diags.HasErrors() {
		return nil, fmt.Errorf("failed to parse HCL content: %s", diags.Error())
	}

	// Create result structure
	result := &spookytypesfacts.CollectorFacts{
		Host:    &spookytypesfacts.HostFacts{},
		CPU:     &spookytypesfacts.CPUFacts{},
		Memory:  &spookytypesfacts.MemoryFacts{},
		Disks:   []*spookytypesfacts.DiskFacts{},
		Network: &spookytypesfacts.NetworkFacts{},
	}

	// Parse the body
	body := file.Body.(*hclsyntax.Body)

	// Parse host block
	if hostBlock := p.findBlock(body, "host"); hostBlock != nil {
		if err := p.parseHostBlock(hostBlock, result.Host); err != nil {
			return nil, fmt.Errorf("failed to parse host block: %w", err)
		}
	}

	// Parse cpu block
	if cpuBlock := p.findBlock(body, "cpu"); cpuBlock != nil {
		if err := p.parseCPUBlock(cpuBlock, result.CPU); err != nil {
			return nil, fmt.Errorf("failed to parse cpu block: %w", err)
		}
	}

	// Parse memory block
	if memoryBlock := p.findBlock(body, "memory"); memoryBlock != nil {
		if err := p.parseMemoryBlock(memoryBlock, result.Memory); err != nil {
			return nil, fmt.Errorf("failed to parse memory block: %w", err)
		}
	}

	// Parse disks blocks
	for _, diskBlock := range p.findBlocks(body, "disk") {
		disk := &spookytypesfacts.DiskFacts{}
		if err := p.parseDiskBlock(diskBlock, disk); err != nil {
			return nil, fmt.Errorf("failed to parse disk block: %w", err)
		}
		result.Disks = append(result.Disks, disk)
	}

	// Parse network block
	if networkBlock := p.findBlock(body, "network"); networkBlock != nil {
		if err := p.parseNetworkBlock(networkBlock, result.Network); err != nil {
			return nil, fmt.Errorf("failed to parse network block: %w", err)
		}
	}

	return result, nil
}

// ParseCustomFacts parses custom facts from HCL content
func (p *HCLParser) ParseCustomFacts(content string) (map[string]interface{}, error) {
	parser := hclparse.NewParser()

	// Parse the HCL content
	file, diags := parser.ParseHCL([]byte(content), "custom-facts.hcl")
	if diags.HasErrors() {
		return nil, fmt.Errorf("failed to parse HCL content: %s", diags.Error())
	}

	// Parse the body and extract all attributes
	body := file.Body.(*hclsyntax.Body)
	result := make(map[string]interface{})

	// Parse all attributes in the custom facts
	for name, attr := range body.Attributes {
		if val, err := p.parseCustomValue(attr.Expr); err == nil {
			result[name] = val
		}
	}

	return result, nil
}

// Helper methods for HCL parsing

// findBlock finds a single block with the given type
func (p *HCLParser) findBlock(body *hclsyntax.Body, blockType string) *hclsyntax.Block {
	for _, block := range body.Blocks {
		if block.Type == blockType {
			return block
		}
	}
	return nil
}

// findBlocks finds all blocks with the given type
func (p *HCLParser) findBlocks(body *hclsyntax.Body, blockType string) []*hclsyntax.Block {
	var blocks []*hclsyntax.Block
	for _, block := range body.Blocks {
		if block.Type == blockType {
			blocks = append(blocks, block)
		}
	}
	return blocks
}

// parseBlockAttributes is a generic helper for parsing block attributes
func (p *HCLParser) parseBlockAttributes(block *hclsyntax.Block, handlers map[string]func(hclsyntax.Expression) error) error {
	for name, attr := range block.Body.Attributes {
		if handler, exists := handlers[name]; exists {
			if err := handler(attr.Expr); err != nil {
				return fmt.Errorf("failed to parse attribute %s: %w", name, err)
			}
		}
	}
	return nil
}

// parseHostBlock parses host facts from a host block
func (p *HCLParser) parseHostBlock(block *hclsyntax.Block, host *spookytypesfacts.HostFacts) error {
	for _, attr := range block.Body.Attributes {
		switch attr.Name {
		case "hostname":
			if val, err := p.parseStringValue(attr.Expr); err == nil {
				host.Hostname = val
			}
		case "uptime":
			if val, err := p.parseInt64Value(attr.Expr); err == nil {
				host.Uptime = val
			}
		case "boot_time":
			if val, err := p.parseInt64Value(attr.Expr); err == nil {
				host.BootTime = val
			}
		case "os":
			if val, err := p.parseStringValue(attr.Expr); err == nil {
				host.OS = val
			}
		case "platform":
			if val, err := p.parseStringValue(attr.Expr); err == nil {
				host.Platform = val
			}
		case "platform_family":
			if val, err := p.parseStringValue(attr.Expr); err == nil {
				host.PlatformFamily = val
			}
		case "platform_version":
			if val, err := p.parseStringValue(attr.Expr); err == nil {
				host.PlatformVersion = val
			}
		case "kernel_version":
			if val, err := p.parseStringValue(attr.Expr); err == nil {
				host.KernelVersion = val
			}
		case "kernel_arch":
			if val, err := p.parseStringValue(attr.Expr); err == nil {
				host.KernelArch = val
			}
		case "virtualization_system":
			if val, err := p.parseStringValue(attr.Expr); err == nil {
				host.VirtualizationSystem = val
			}
		case "virtualization_role":
			if val, err := p.parseStringValue(attr.Expr); err == nil {
				host.VirtualizationRole = val
			}
		}
	}
	return nil
}

// parseCPUBlock parses CPU facts from a cpu block
func (p *HCLParser) parseCPUBlock(block *hclsyntax.Block, cpu *spookytypesfacts.CPUFacts) error {
	return p.parseBlockAttributes(block, map[string]func(hclsyntax.Expression) error{
		"cores": func(expr hclsyntax.Expression) error {
			if val, err := p.parseIntValue(expr); err == nil {
				cpu.Cores = val
			}
			return nil
		},
		"model": func(expr hclsyntax.Expression) error {
			if val, err := p.parseStringValue(expr); err == nil {
				cpu.Model = val
			}
			return nil
		},
		"vendor": func(expr hclsyntax.Expression) error {
			if val, err := p.parseStringValue(expr); err == nil {
				cpu.Vendor = val
			}
			return nil
		},
		"frequency": func(expr hclsyntax.Expression) error {
			if val, err := p.parseFloat64Value(expr); err == nil {
				cpu.Frequency = val
			}
			return nil
		},
	})
}

// parseMemoryBlock parses memory facts from a memory block
func (p *HCLParser) parseMemoryBlock(block *hclsyntax.Block, memory *spookytypesfacts.MemoryFacts) error {
	return p.parseBlockAttributes(block, map[string]func(hclsyntax.Expression) error{
		"total": func(expr hclsyntax.Expression) error {
			if val, err := p.parseInt64Value(expr); err == nil {
				memory.Total = val
			}
			return nil
		},
		"available": func(expr hclsyntax.Expression) error {
			if val, err := p.parseInt64Value(expr); err == nil {
				memory.Available = val
			}
			return nil
		},
		"used": func(expr hclsyntax.Expression) error {
			if val, err := p.parseInt64Value(expr); err == nil {
				memory.Used = val
			}
			return nil
		},
		"free": func(expr hclsyntax.Expression) error {
			if val, err := p.parseInt64Value(expr); err == nil {
				memory.Free = val
			}
			return nil
		},
	})
}

// parseDiskBlock parses disk facts from a disk block
func (p *HCLParser) parseDiskBlock(block *hclsyntax.Block, disk *spookytypesfacts.DiskFacts) error {
	for _, attr := range block.Body.Attributes {
		switch attr.Name {
		case "device":
			if val, err := p.parseStringValue(attr.Expr); err == nil {
				disk.Device = val
			}
		case "mount_point":
			if val, err := p.parseStringValue(attr.Expr); err == nil {
				disk.MountPoint = val
			}
		case "filesystem":
			if val, err := p.parseStringValue(attr.Expr); err == nil {
				disk.Filesystem = val
			}
		case "total":
			if val, err := p.parseInt64Value(attr.Expr); err == nil {
				disk.Total = val
			}
		case "free":
			if val, err := p.parseInt64Value(attr.Expr); err == nil {
				disk.Free = val
			}
		case "used":
			if val, err := p.parseInt64Value(attr.Expr); err == nil {
				disk.Used = val
			}
		}
	}
	return nil
}

// parseNetworkBlock parses network facts from a network block
func (p *HCLParser) parseNetworkBlock(block *hclsyntax.Block, network *spookytypesfacts.NetworkFacts) error {
	for _, attr := range block.Body.Attributes {
		switch attr.Name {
		case "hostname":
			if val, err := p.parseStringValue(attr.Expr); err == nil {
				network.Hostname = val
			}
		case "primary_ip":
			if val, err := p.parseStringValue(attr.Expr); err == nil {
				network.PrimaryIP = val
			}
		}
	}
	return nil
}

// Value parsing helpers

func (p *HCLParser) parseStringValue(expr hclsyntax.Expression) (string, error) {
	if lit, ok := expr.(*hclsyntax.LiteralValueExpr); ok {
		if lit.Val.Type().IsPrimitiveType() && lit.Val.Type().FriendlyName() == "string" {
			return lit.Val.AsString(), nil
		}
	}
	return "", fmt.Errorf("expected string value")
}

func (p *HCLParser) parseIntValue(expr hclsyntax.Expression) (int, error) {
	if lit, ok := expr.(*hclsyntax.LiteralValueExpr); ok {
		if lit.Val.Type().IsPrimitiveType() && lit.Val.Type().FriendlyName() == "number" {
			val, _ := lit.Val.AsBigFloat().Int64()
			return int(val), nil
		}
	}
	return 0, fmt.Errorf("expected integer value")
}

func (p *HCLParser) parseInt64Value(expr hclsyntax.Expression) (int64, error) {
	if lit, ok := expr.(*hclsyntax.LiteralValueExpr); ok {
		if lit.Val.Type().IsPrimitiveType() && lit.Val.Type().FriendlyName() == "number" {
			val, _ := lit.Val.AsBigFloat().Int64()
			return val, nil
		}
	}
	return 0, fmt.Errorf("expected integer value")
}

func (p *HCLParser) parseFloat64Value(expr hclsyntax.Expression) (float64, error) {
	if lit, ok := expr.(*hclsyntax.LiteralValueExpr); ok {
		if lit.Val.Type().IsPrimitiveType() && lit.Val.Type().FriendlyName() == "number" {
			val, _ := lit.Val.AsBigFloat().Float64()
			return val, nil
		}
	}
	return 0, fmt.Errorf("expected number value")
}

// parseCustomValue parses any HCL value for custom facts
func (p *HCLParser) parseCustomValue(expr hclsyntax.Expression) (interface{}, error) {
	if lit, ok := expr.(*hclsyntax.LiteralValueExpr); ok {
		if lit.Val.Type().IsPrimitiveType() {
			switch lit.Val.Type().FriendlyName() {
			case "string":
				return lit.Val.AsString(), nil
			case "number":
				val, _ := lit.Val.AsBigFloat().Int64()
				return val, nil
			case "bool":
				return lit.Val.True(), nil
			}
		}
	}
	return nil, fmt.Errorf("unsupported value type")
}
