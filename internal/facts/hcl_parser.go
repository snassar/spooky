// Package facts provides fact collection and in-memory management functionality.
package facts

import (
	"fmt"

	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/hcl/v2/hclsyntax"

	spookytypesfacts "spooky/internal/types/facts"
)

// HCL type constants
const (
	HCLTypeNumber = "number"
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

	// Parse all blocks
	if err := p.parseAllBlocks(body, result); err != nil {
		return nil, err
	}

	return result, nil
}

func (p *HCLParser) parseAllBlocks(body *hclsyntax.Body, result *spookytypesfacts.CollectorFacts) error {
	blockParsers := map[string]func(*hclsyntax.Block, *spookytypesfacts.CollectorFacts) error{
		"host":    p.parseHostBlockWrapper,
		"cpu":     p.parseCPUBlockWrapper,
		"memory":  p.parseMemoryBlockWrapper,
		"network": p.parseNetworkBlockWrapper,
	}

	// Parse single blocks
	for blockType, parser := range blockParsers {
		if block := p.findBlock(body, blockType); block != nil {
			if err := parser(block, result); err != nil {
				return fmt.Errorf("failed to parse %s block: %w", blockType, err)
			}
		}
	}

	// Parse multiple disk blocks
	if err := p.parseDiskBlocks(body, result); err != nil {
		return fmt.Errorf("failed to parse disk blocks: %w", err)
	}

	return nil
}

func (p *HCLParser) parseHostBlockWrapper(block *hclsyntax.Block, result *spookytypesfacts.CollectorFacts) error {
	return p.parseHostBlock(block, result.Host)
}

func (p *HCLParser) parseCPUBlockWrapper(block *hclsyntax.Block, result *spookytypesfacts.CollectorFacts) error {
	return p.parseCPUBlock(block, result.CPU)
}

func (p *HCLParser) parseMemoryBlockWrapper(block *hclsyntax.Block, result *spookytypesfacts.CollectorFacts) error {
	return p.parseMemoryBlock(block, result.Memory)
}

func (p *HCLParser) parseNetworkBlockWrapper(block *hclsyntax.Block, result *spookytypesfacts.CollectorFacts) error {
	return p.parseNetworkBlock(block, result.Network)
}

func (p *HCLParser) parseDiskBlocks(body *hclsyntax.Body, result *spookytypesfacts.CollectorFacts) error {
	for _, diskBlock := range p.findBlocks(body, "disk") {
		disk := &spookytypesfacts.DiskFacts{}
		if err := p.parseDiskBlock(diskBlock, disk); err != nil {
			return fmt.Errorf("failed to parse disk block: %w", err)
		}
		result.Disks = append(result.Disks, disk)
	}
	return nil
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
	parsers := map[string]func(hclsyntax.Expression, *spookytypesfacts.HostFacts) error{
		"hostname":              p.parseHostHostname,
		"uptime":                p.parseHostUptime,
		"boot_time":             p.parseHostBootTime,
		"os":                    p.parseHostOS,
		"platform":              p.parseHostPlatform,
		"platform_family":       p.parseHostPlatformFamily,
		"platform_version":      p.parseHostPlatformVersion,
		"kernel_version":        p.parseHostKernelVersion,
		"kernel_arch":           p.parseHostKernelArch,
		"virtualization_system": p.parseHostVirtualizationSystem,
		"virtualization_role":   p.parseHostVirtualizationRole,
	}

	for name, attr := range block.Body.Attributes {
		if parser, exists := parsers[name]; exists {
			if err := parser(attr.Expr, host); err != nil {
				return fmt.Errorf("failed to parse host attribute %s: %w", name, err)
			}
		}
	}
	return nil
}

func (p *HCLParser) parseHostHostname(expr hclsyntax.Expression, host *spookytypesfacts.HostFacts) error {
	if val, err := p.parseStringValue(expr); err == nil {
		host.Hostname = val
	}
	return nil
}

func (p *HCLParser) parseHostUptime(expr hclsyntax.Expression, host *spookytypesfacts.HostFacts) error {
	if val, err := p.parseInt64Value(expr); err == nil {
		host.Uptime = val
	}
	return nil
}

func (p *HCLParser) parseHostBootTime(expr hclsyntax.Expression, host *spookytypesfacts.HostFacts) error {
	if val, err := p.parseInt64Value(expr); err == nil {
		host.BootTime = val
	}
	return nil
}

func (p *HCLParser) parseHostOS(expr hclsyntax.Expression, host *spookytypesfacts.HostFacts) error {
	if val, err := p.parseStringValue(expr); err == nil {
		host.OS = val
	}
	return nil
}

func (p *HCLParser) parseHostPlatform(expr hclsyntax.Expression, host *spookytypesfacts.HostFacts) error {
	if val, err := p.parseStringValue(expr); err == nil {
		host.Platform = val
	}
	return nil
}

func (p *HCLParser) parseHostPlatformFamily(expr hclsyntax.Expression, host *spookytypesfacts.HostFacts) error {
	if val, err := p.parseStringValue(expr); err == nil {
		host.PlatformFamily = val
	}
	return nil
}

func (p *HCLParser) parseHostPlatformVersion(expr hclsyntax.Expression, host *spookytypesfacts.HostFacts) error {
	if val, err := p.parseStringValue(expr); err == nil {
		host.PlatformVersion = val
	}
	return nil
}

func (p *HCLParser) parseHostKernelVersion(expr hclsyntax.Expression, host *spookytypesfacts.HostFacts) error {
	if val, err := p.parseStringValue(expr); err == nil {
		host.KernelVersion = val
	}
	return nil
}

func (p *HCLParser) parseHostKernelArch(expr hclsyntax.Expression, host *spookytypesfacts.HostFacts) error {
	if val, err := p.parseStringValue(expr); err == nil {
		host.KernelArch = val
	}
	return nil
}

func (p *HCLParser) parseHostVirtualizationSystem(expr hclsyntax.Expression, host *spookytypesfacts.HostFacts) error {
	if val, err := p.parseStringValue(expr); err == nil {
		host.VirtualizationSystem = val
	}
	return nil
}

func (p *HCLParser) parseHostVirtualizationRole(expr hclsyntax.Expression, host *spookytypesfacts.HostFacts) error {
	if val, err := p.parseStringValue(expr); err == nil {
		host.VirtualizationRole = val
	}
	return nil
}

// createParser creates a parser function for a specific field and type
func createParser[T any](field *T, parser func(hclsyntax.Expression) (T, error)) func(hclsyntax.Expression) error {
	return func(expr hclsyntax.Expression) error {
		if val, err := parser(expr); err == nil {
			*field = val
		}
		return nil
	}
}

// parseCPUBlock parses CPU facts from a cpu block
func (p *HCLParser) parseCPUBlock(block *hclsyntax.Block, cpu *spookytypesfacts.CPUFacts) error {
	parsers := map[string]func(hclsyntax.Expression) error{
		"cores":     createParser(&cpu.Cores, p.parseIntValue),
		"model":     createParser(&cpu.Model, p.parseStringValue),
		"vendor":    createParser(&cpu.Vendor, p.parseStringValue),
		"frequency": createParser(&cpu.Frequency, p.parseFloat64Value),
	}
	return p.parseBlockAttributes(block, parsers)
}

// parseMemoryBlock parses memory facts from a memory block
func (p *HCLParser) parseMemoryBlock(block *hclsyntax.Block, memory *spookytypesfacts.MemoryFacts) error {
	parsers := map[string]func(hclsyntax.Expression) error{
		"total":     createParser(&memory.Total, p.parseInt64Value),
		"available": createParser(&memory.Available, p.parseInt64Value),
		"used":      createParser(&memory.Used, p.parseInt64Value),
		"free":      createParser(&memory.Free, p.parseInt64Value),
	}
	return p.parseBlockAttributes(block, parsers)
}

// parseDiskBlock parses disk facts from a disk block
func (p *HCLParser) parseDiskBlock(block *hclsyntax.Block, disk *spookytypesfacts.DiskFacts) error {
	for _, attr := range block.Body.Attributes {
		if err := p.parseDiskAttribute(attr, disk); err != nil {
			return fmt.Errorf("failed to parse disk attribute %s: %w", attr.Name, err)
		}
	}
	return nil
}

func (p *HCLParser) parseDiskAttribute(attr *hclsyntax.Attribute, disk *spookytypesfacts.DiskFacts) error {
	parsers := map[string]func(hclsyntax.Expression, *spookytypesfacts.DiskFacts) error{
		"device":      p.parseDiskDevice,
		"mount_point": p.parseDiskMountPoint,
		"filesystem":  p.parseDiskFilesystem,
		"total":       p.parseDiskTotal,
		"free":        p.parseDiskFree,
		"used":        p.parseDiskUsed,
	}

	if parser, exists := parsers[attr.Name]; exists {
		return parser(attr.Expr, disk)
	}
	return nil
}

func (p *HCLParser) parseDiskDevice(expr hclsyntax.Expression, disk *spookytypesfacts.DiskFacts) error {
	if val, err := p.parseStringValue(expr); err == nil {
		disk.Device = val
	}
	return nil
}

func (p *HCLParser) parseDiskMountPoint(expr hclsyntax.Expression, disk *spookytypesfacts.DiskFacts) error {
	if val, err := p.parseStringValue(expr); err == nil {
		disk.MountPoint = val
	}
	return nil
}

func (p *HCLParser) parseDiskFilesystem(expr hclsyntax.Expression, disk *spookytypesfacts.DiskFacts) error {
	if val, err := p.parseStringValue(expr); err == nil {
		disk.Filesystem = val
	}
	return nil
}

func (p *HCLParser) parseDiskTotal(expr hclsyntax.Expression, disk *spookytypesfacts.DiskFacts) error {
	if val, err := p.parseInt64Value(expr); err == nil {
		disk.Total = val
	}
	return nil
}

func (p *HCLParser) parseDiskFree(expr hclsyntax.Expression, disk *spookytypesfacts.DiskFacts) error {
	if val, err := p.parseInt64Value(expr); err == nil {
		disk.Free = val
	}
	return nil
}

func (p *HCLParser) parseDiskUsed(expr hclsyntax.Expression, disk *spookytypesfacts.DiskFacts) error {
	if val, err := p.parseInt64Value(expr); err == nil {
		disk.Used = val
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
		if lit.Val.Type().IsPrimitiveType() && lit.Val.Type().FriendlyName() == HCLTypeNumber {
			val, _ := lit.Val.AsBigFloat().Int64()
			return int(val), nil
		}
	}
	return 0, fmt.Errorf("expected integer value")
}

func (p *HCLParser) parseInt64Value(expr hclsyntax.Expression) (int64, error) {
	if lit, ok := expr.(*hclsyntax.LiteralValueExpr); ok {
		if lit.Val.Type().IsPrimitiveType() && lit.Val.Type().FriendlyName() == HCLTypeNumber {
			val, _ := lit.Val.AsBigFloat().Int64()
			return val, nil
		}
	}
	return 0, fmt.Errorf("expected integer value")
}

func (p *HCLParser) parseFloat64Value(expr hclsyntax.Expression) (float64, error) {
	if lit, ok := expr.(*hclsyntax.LiteralValueExpr); ok {
		if lit.Val.Type().IsPrimitiveType() && lit.Val.Type().FriendlyName() == HCLTypeNumber {
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
			case HCLTypeNumber:
				val, _ := lit.Val.AsBigFloat().Int64()
				return val, nil
			case "bool":
				return lit.Val.True(), nil
			}
		}
	}
	return nil, fmt.Errorf("unsupported value type")
}
