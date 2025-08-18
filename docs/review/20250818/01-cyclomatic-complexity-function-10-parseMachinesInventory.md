# Function Improvement Plan: parseMachinesInventory

**Function:** `parseMachinesInventory`  
**File:** `internal/machines/loader.go:XXX`  
**Current Complexity:** 9  
**Target Complexity:** < 8  
**Priority:** Medium

## Current Function Analysis

### Function Signature
```go
func (l *Loader) parseMachinesInventory(content []byte) (*spookytypesmachines.MachinesInventory, error)
```

### Current Issues
1. **Complex HCL parsing logic** - Multiple conditions for different HCL structures and validation
2. **Mixed responsibilities** - HCL parsing, validation, inventory creation, and error handling
3. **Deep nesting** - Multiple levels of conditional logic for different parsing scenarios
4. **Complex block processing** - Multiple conditions for machine block extraction and validation
5. **Repetitive validation patterns** - Similar logic repeated for different validation checks

### Complexity Breakdown
- 9 cyclomatic complexity points from:
  - 1 main function entry
  - 2 content validation checks
  - 3 HCL parsing conditions
  - 2 error handling paths
  - 1 block processing condition

## Refactoring Strategy

### Phase 1: Extract Content Validation (Immediate - 1 day)

#### Extract Content Validation
```go
func (l *Loader) validateInventoryContent(content []byte) error {
    if content == nil {
        return fmt.Errorf("content cannot be nil")
    }
    
    if len(content) == 0 {
        return fmt.Errorf("content cannot be empty")
    }
    
    if len(content) > maxInventorySize {
        return fmt.Errorf("content size exceeds maximum allowed size of %d bytes", maxInventorySize)
    }
    
    return nil
}
```

#### Extract HCL Parsing
```go
func (l *Loader) parseHCLContent(content []byte) (*hcl.File, error) {
    file, diags := hcl.ParseBytes(content, "machines.hcl")
    if diags.HasErrors() {
        return nil, fmt.Errorf("failed to parse HCL content: %v", diags)
    }
    
    return file, nil
}
```

### Phase 2: Extract Block Processing (Day 2)

#### Extract Machine Block Processing
```go
func (l *Loader) extractMachineBlocks(file *hcl.File) ([]*hcl.Block, error) {
    var machineBlocks []*hcl.Block
    
    for _, block := range file.Body.(*hcl.BodySchema).Blocks {
        if block.Type == "machine" {
            machineBlocks = append(machineBlocks, block)
        }
    }
    
    if len(machineBlocks) == 0 {
        return nil, fmt.Errorf("no machine blocks found in inventory")
    }
    
    return machineBlocks, nil
}

func (l *Loader) validateMachineBlocks(blocks []*hcl.Block) error {
    machineNames := make(map[string]bool)
    
    for _, block := range blocks {
        if len(block.Labels) == 0 {
            return fmt.Errorf("machine block missing name label")
        }
        
        name := block.Labels[0]
        if name == "" {
            return fmt.Errorf("machine name cannot be empty")
        }
        
        if machineNames[name] {
            return fmt.Errorf("duplicate machine name: %s", name)
        }
        
        machineNames[name] = true
    }
    
    return nil
}
```

#### Extract Machine Creation
```go
func (l *Loader) createMachinesFromBlocks(blocks []*hcl.Block) ([]*spookytypesmachines.Machine, error) {
    var machines []*spookytypesmachines.Machine
    
    for _, block := range blocks {
        machine, err := l.parseMachineBlock(block)
        if err != nil {
            return nil, fmt.Errorf("failed to parse machine block %s: %w", block.Labels[0], err)
        }
        
        machines = append(machines, machine)
    }
    
    return machines, nil
}
```

### Phase 3: Extract Inventory Creation (Day 3)

#### Extract Inventory Creation Logic
```go
func (l *Loader) createInventoryFromMachines(machines []*spookytypesmachines.Machine) (*spookytypesmachines.MachinesInventory, error) {
    inventory := &spookytypesmachines.MachinesInventory{
        Machines:  machines,
        Count:     len(machines),
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }
    
    // Calculate inventory statistics
    if err := l.calculateInventoryStats(inventory); err != nil {
        return nil, fmt.Errorf("failed to calculate inventory statistics: %w", err)
    }
    
    // Validate inventory
    if err := l.validateInventory(inventory); err != nil {
        return nil, fmt.Errorf("inventory validation failed: %w", err)
    }
    
    return inventory, nil
}

func (l *Loader) calculateInventoryStats(inventory *spookytypesmachines.MachinesInventory) error {
    stats := &spookytypesmachines.InventoryStats{
        TotalMachines: len(inventory.Machines),
        UniqueHosts:   make(map[string]int),
        UniqueUsers:   make(map[string]int),
        PortUsage:     make(map[int]int),
        TagUsage:      make(map[string]int),
    }
    
    for _, machine := range inventory.Machines {
        // Count unique hosts
        stats.UniqueHosts[machine.Hostname]++
        
        // Count unique users
        stats.UniqueUsers[machine.User]++
        
        // Count port usage
        stats.PortUsage[machine.Port]++
        
        // Count tag usage
        for _, tag := range machine.Tags {
            stats.TagUsage[tag]++
        }
    }
    
    inventory.Stats = stats
    return nil
}

func (l *Loader) validateInventory(inventory *spookytypesmachines.MachinesInventory) error {
    if inventory == nil {
        return fmt.Errorf("inventory cannot be nil")
    }
    
    if len(inventory.Machines) == 0 {
        return fmt.Errorf("inventory must contain at least one machine")
    }
    
    // Validate individual machines
    for i, machine := range inventory.Machines {
        if err := l.validateMachine(machine); err != nil {
            return fmt.Errorf("machine[%d] validation failed: %w", i, err)
        }
    }
    
    return nil
}

func (l *Loader) validateMachine(machine *spookytypesmachines.Machine) error {
    if machine == nil {
        return fmt.Errorf("machine cannot be nil")
    }
    
    if machine.Name == "" {
        return fmt.Errorf("machine name cannot be empty")
    }
    
    if machine.Hostname == "" {
        return fmt.Errorf("machine hostname cannot be empty")
    }
    
    if machine.Port <= 0 || machine.Port > 65535 {
        return fmt.Errorf("machine port must be between 1 and 65535")
    }
    
    if machine.User == "" {
        return fmt.Errorf("machine user cannot be empty")
    }
    
    return nil
}
```

### Phase 4: Refactored Main Function (Day 4)

#### Final Refactored Function
```go
func (l *Loader) parseMachinesInventory(content []byte) (*spookytypesmachines.MachinesInventory, error) {
    // Validate content
    if err := l.validateInventoryContent(content); err != nil {
        return nil, err
    }
    
    // Parse HCL content
    file, err := l.parseHCLContent(content)
    if err != nil {
        return nil, err
    }
    
    // Extract machine blocks
    machineBlocks, err := l.extractMachineBlocks(file)
    if err != nil {
        return nil, err
    }
    
    // Validate machine blocks
    if err := l.validateMachineBlocks(machineBlocks); err != nil {
        return nil, err
    }
    
    // Create machines from blocks
    machines, err := l.createMachinesFromBlocks(machineBlocks)
    if err != nil {
        return nil, err
    }
    
    // Create inventory from machines
    return l.createInventoryFromMachines(machines)
}
```

## Complexity Reduction Metrics

### Before Refactoring
- **Cyclomatic Complexity:** 9
- **Lines of Code:** ~70
- **Nesting Levels:** 4-5
- **Responsibilities:** 4 (validation, parsing, creation, error handling)

### After Refactoring
- **Cyclomatic Complexity:** 3 (main function) + 1-3 per helper function
- **Lines of Code:** ~180 (distributed across multiple focused functions)
- **Nesting Levels:** 1-2
- **Responsibilities:** 1 per function (single responsibility principle)

## Testing Strategy

### Unit Tests for Extracted Functions
```go
func TestValidateInventoryContent(t *testing.T) {
    tests := []struct {
        name    string
        content []byte
        wantErr bool
    }{
        {
            name:    "nil content",
            content: nil,
            wantErr: true,
        },
        {
            name:    "empty content",
            content: []byte{},
            wantErr: true,
        },
        {
            name:    "valid content",
            content: []byte("machines { }"),
            wantErr: false,
        },
        {
            name:    "content too large",
            content: make([]byte, maxInventorySize+1),
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := validateInventoryContent(tt.content)
            if (err != nil) != tt.wantErr {
                t.Errorf("validateInventoryContent() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}

func TestExtractMachineBlocks(t *testing.T) {
    tests := []struct {
        name    string
        content string
        want    int
        wantErr bool
    }{
        {
            name:    "no machine blocks",
            content: "other { }",
            want:    0,
            wantErr: true,
        },
        {
            name:    "one machine block",
            content: `machine "test" { hostname = "test.example.com" }`,
            want:    1,
            wantErr: false,
        },
        {
            name:    "multiple machine blocks",
            content: `
                machine "test1" { hostname = "test1.example.com" }
                machine "test2" { hostname = "test2.example.com" }
            `,
            want:    2,
            wantErr: false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            file, err := parseHCLContent([]byte(tt.content))
            if err != nil {
                t.Errorf("parseHCLContent() error = %v", err)
                return
            }
            
            blocks, err := extractMachineBlocks(file)
            if (err != nil) != tt.wantErr {
                t.Errorf("extractMachineBlocks() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if len(blocks) != tt.want {
                t.Errorf("extractMachineBlocks() returned %d blocks, want %d", len(blocks), tt.want)
            }
        })
    }
}

func TestValidateMachineBlocks(t *testing.T) {
    tests := []struct {
        name    string
        blocks  []*hcl.Block
        wantErr bool
    }{
        {
            name: "duplicate machine names",
            blocks: []*hcl.Block{
                {Labels: []string{"test"}},
                {Labels: []string{"test"}},
            },
            wantErr: true,
        },
        {
            name: "empty machine name",
            blocks: []*hcl.Block{
                {Labels: []string{""}},
            },
            wantErr: true,
        },
        {
            name: "valid machine blocks",
            blocks: []*hcl.Block{
                {Labels: []string{"test1"}},
                {Labels: []string{"test2"}},
            },
            wantErr: false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := validateMachineBlocks(tt.blocks)
            if (err != nil) != tt.wantErr {
                t.Errorf("validateMachineBlocks() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}

func TestValidateMachine(t *testing.T) {
    tests := []struct {
        name    string
        machine *spookytypesmachines.Machine
        wantErr bool
    }{
        {
            name:    "nil machine",
            machine: nil,
            wantErr: true,
        },
        {
            name: "empty name",
            machine: &spookytypesmachines.Machine{
                Name:     "",
                Hostname: "test.example.com",
                Port:     22,
                User:     "admin",
            },
            wantErr: true,
        },
        {
            name: "empty hostname",
            machine: &spookytypesmachines.Machine{
                Name:     "test",
                Hostname: "",
                Port:     22,
                User:     "admin",
            },
            wantErr: true,
        },
        {
            name: "invalid port",
            machine: &spookytypesmachines.Machine{
                Name:     "test",
                Hostname: "test.example.com",
                Port:     0,
                User:     "admin",
            },
            wantErr: true,
        },
        {
            name: "valid machine",
            machine: &spookytypesmachines.Machine{
                Name:     "test",
                Hostname: "test.example.com",
                Port:     22,
                User:     "admin",
            },
            wantErr: false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := validateMachine(tt.machine)
            if (err != nil) != tt.wantErr {
                t.Errorf("validateMachine() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

### Integration Tests
```go
func TestParseMachinesInventoryIntegration(t *testing.T) {
    tests := []struct {
        name    string
        content string
        want    *spookytypesmachines.MachinesInventory
        wantErr bool
    }{
        {
            name: "valid inventory",
            content: `
                machine "test1" {
                    hostname = "test1.example.com"
                    port = 22
                    user = "admin"
                }
                machine "test2" {
                    hostname = "test2.example.com"
                    port = 2222
                    user = "root"
                }
            `,
            want: &spookytypesmachines.MachinesInventory{
                Count: 2,
            },
            wantErr: false,
        },
        {
            name: "invalid HCL",
            content: `invalid hcl content {`,
            want:    nil,
            wantErr: true,
        },
        {
            name: "no machine blocks",
            content: `other { }`,
            want:    nil,
            wantErr: true,
        },
        {
            name: "duplicate machine names",
            content: `
                machine "test" { hostname = "test1.example.com" }
                machine "test" { hostname = "test2.example.com" }
            `,
            want:    nil,
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            inventory, err := parseMachinesInventory([]byte(tt.content))
            if (err != nil) != tt.wantErr {
                t.Errorf("parseMachinesInventory() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if !tt.wantErr && inventory.Count != tt.want.Count {
                t.Errorf("parseMachinesInventory() returned %d machines, want %d", inventory.Count, tt.want.Count)
            }
        })
    }
}
```

## Implementation Timeline

### Day 1: Extract Content Validation
- [ ] Extract `validateInventoryContent`
- [ ] Extract `parseHCLContent`
- [ ] Add unit tests for validation functions
- [ ] Verify content validation works correctly

### Day 2: Extract Block Processing
- [ ] Extract `extractMachineBlocks`
- [ ] Extract `validateMachineBlocks`
- [ ] Extract `createMachinesFromBlocks`
- [ ] Add unit tests for block processing functions
- [ ] Verify block processing works correctly

### Day 3: Extract Inventory Creation
- [ ] Extract `createInventoryFromMachines`
- [ ] Extract `calculateInventoryStats`
- [ ] Extract `validateInventory`
- [ ] Extract `validateMachine`
- [ ] Add unit tests for inventory creation functions
- [ ] Verify inventory creation works correctly

### Day 4: Complete Refactoring
- [ ] Refactor main `parseMachinesInventory` function
- [ ] Add integration tests
- [ ] Verify complexity reduction with gocyclo
- [ ] Code review and documentation
- [ ] Performance testing

## Success Criteria

### Complexity Reduction
- [ ] Main function complexity < 8
- [ ] All extracted functions complexity < 5
- [ ] No function exceeds complexity threshold

### Code Quality
- [ ] Single responsibility principle maintained
- [ ] Clear separation of concerns
- [ ] Comprehensive test coverage (>90%)
- [ ] No regression in functionality

### Maintainability
- [ ] Easy to modify individual inventory processing components
- [ ] Clear function names and purposes
- [ ] Well-documented functions
- [ ] Consistent inventory processing patterns

## Risk Mitigation

### Potential Risks
1. **HCL parsing changes** - May affect inventory parsing behavior
2. **Block processing changes** - May affect machine block extraction
3. **Validation changes** - May affect inventory validation

### Mitigation Strategies
1. **Comprehensive testing** - Test all inventory scenarios and edge cases
2. **Parsing validation** - Ensure HCL parsing behavior remains consistent
3. **Block validation** - Ensure block processing remains consistent
4. **Gradual migration** - Implement changes incrementally with validation

## Verification

### Complexity Verification
```bash
# Verify complexity reduction
gocyclo internal/machines/loader.go | grep parseMachinesInventory
```

### Functionality Verification
```bash
# Test inventory parsing
go test ./internal/machines -run TestParseMachinesInventory
```

### Performance Verification
- [ ] No performance regression
- [ ] Memory usage remains stable
- [ ] Inventory parsing performance remains acceptable

## Conclusion

This refactoring will reduce the cyclomatic complexity of `parseMachinesInventory` from 9 to approximately 3, while improving code maintainability, testability, and readability. The extracted functions will follow single responsibility principles and be easily testable in isolation.

**Expected Outcome:** A clean, maintainable function with clear separation of concerns and comprehensive test coverage for machines inventory parsing operations.
