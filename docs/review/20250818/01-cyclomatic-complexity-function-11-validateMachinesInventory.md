# Function Improvement Plan: validateMachinesInventory

**Function:** `validateMachinesInventory`  
**File:** `internal/machines/validator.go:XXX`  
**Current Complexity:** 9  
**Target Complexity:** < 8  
**Priority:** Medium

## Current Function Analysis

### Function Signature
```go
func (v *Validator) validateMachinesInventory(inventory *spookytypesmachines.MachinesInventory) (*spookytypesmachines.ValidationResult, error)
```

### Current Issues
1. **Complex validation logic** - Multiple conditions for different validation types and scenarios
2. **Mixed responsibilities** - Inventory validation, machine validation, result aggregation, and error handling
3. **Deep nesting** - Multiple levels of conditional logic for different validation checks
4. **Complex result processing** - Multiple conditions for validation result creation and aggregation
5. **Repetitive validation patterns** - Similar logic repeated for different validation types

### Complexity Breakdown
- 9 cyclomatic complexity points from:
  - 1 main function entry
  - 2 inventory validation checks
  - 3 machine validation conditions
  - 2 error handling paths
  - 1 result processing condition

## Refactoring Strategy

### Phase 1: Extract Input Validation (Immediate - 1 day)

#### Extract Input Validation
```go
func (v *Validator) validateInventoryInput(inventory *spookytypesmachines.MachinesInventory) error {
    if inventory == nil {
        return fmt.Errorf("inventory cannot be nil")
    }
    
    if inventory.Machines == nil {
        return fmt.Errorf("inventory machines cannot be nil")
    }
    
    if len(inventory.Machines) == 0 {
        return fmt.Errorf("inventory must contain at least one machine")
    }
    
    return nil
}
```

#### Extract Result Initialization
```go
func (v *Validator) initializeValidationResult() *spookytypesmachines.ValidationResult {
    return &spookytypesmachines.ValidationResult{
        Valid:       true,
        Errors:      []string{},
        Warnings:    []string{},
        ValidatedAt: time.Now(),
    }
}
```

### Phase 2: Extract Machine Validation (Day 2)

#### Extract Individual Machine Validation
```go
func (v *Validator) validateIndividualMachines(machines []*spookytypesmachines.Machine) ([]string, []string) {
    var errors []string
    var warnings []string
    
    for i, machine := range machines {
        machineErrors, machineWarnings := v.validateMachine(machine, i)
        errors = append(errors, machineErrors...)
        warnings = append(warnings, machineWarnings...)
    }
    
    return errors, warnings
}

func (v *Validator) validateMachine(machine *spookytypesmachines.Machine, index int) ([]string, []string) {
    var errors []string
    var warnings []string
    
    // Validate required fields
    if err := v.validateRequiredFields(machine, index); err != nil {
        errors = append(errors, err.Error())
    }
    
    // Validate field formats
    if warnings := v.validateFieldFormats(machine, index); len(warnings) > 0 {
        warnings = append(warnings, warnings...)
    }
    
    // Validate connectivity
    if err := v.validateConnectivity(machine, index); err != nil {
        errors = append(errors, err.Error())
    }
    
    // Validate authentication
    if err := v.validateAuthentication(machine, index); err != nil {
        errors = append(errors, err.Error())
    }
    
    return errors, warnings
}

func (v *Validator) validateRequiredFields(machine *spookytypesmachines.Machine, index int) error {
    if machine.Name == "" {
        return fmt.Errorf("machine[%d]: name is required", index)
    }
    
    if machine.Hostname == "" {
        return fmt.Errorf("machine[%d]: hostname is required", index)
    }
    
    if machine.User == "" {
        return fmt.Errorf("machine[%d]: user is required", index)
    }
    
    return nil
}

func (v *Validator) validateFieldFormats(machine *spookytypesmachines.Machine, index int) []string {
    var warnings []string
    
    // Validate hostname format
    if !v.isValidHostname(machine.Hostname) {
        warnings = append(warnings, fmt.Sprintf("machine[%d]: hostname '%s' may not be valid", index, machine.Hostname))
    }
    
    // Validate port range
    if machine.Port <= 0 || machine.Port > 65535 {
        warnings = append(warnings, fmt.Sprintf("machine[%d]: port %d is outside valid range (1-65535)", index, machine.Port))
    }
    
    // Validate user format
    if !v.isValidUsername(machine.User) {
        warnings = append(warnings, fmt.Sprintf("machine[%d]: username '%s' may not be valid", index, machine.User))
    }
    
    return warnings
}

func (v *Validator) validateConnectivity(machine *spookytypesmachines.Machine, index int) error {
    // Check if hostname is reachable
    if !v.isHostnameReachable(machine.Hostname) {
        return fmt.Errorf("machine[%d]: hostname '%s' is not reachable", index, machine.Hostname)
    }
    
    // Check if port is open
    if !v.isPortOpen(machine.Hostname, machine.Port) {
        return fmt.Errorf("machine[%d]: port %d is not open on '%s'", index, machine.Port, machine.Hostname)
    }
    
    return nil
}

func (v *Validator) validateAuthentication(machine *spookytypesmachines.Machine, index int) error {
    if machine.Authentication == nil {
        return fmt.Errorf("machine[%d]: authentication configuration is required", index)
    }
    
    // Validate authentication method
    if err := v.validateAuthMethod(machine.Authentication, index); err != nil {
        return err
    }
    
    // Validate authentication credentials
    if err := v.validateAuthCredentials(machine.Authentication, index); err != nil {
        return err
    }
    
    return nil
}
```

#### Extract Validation Helper Functions
```go
func (v *Validator) isValidHostname(hostname string) bool {
    // Basic hostname validation
    if len(hostname) == 0 || len(hostname) > 253 {
        return false
    }
    
    // Check for valid characters
    hostnameRegex := regexp.MustCompile(`^[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?$`)
    return hostnameRegex.MatchString(hostname)
}

func (v *Validator) isValidUsername(username string) bool {
    // Basic username validation
    if len(username) == 0 || len(username) > 32 {
        return false
    }
    
    // Check for valid characters
    usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
    return usernameRegex.MatchString(username)
}

func (v *Validator) isHostnameReachable(hostname string) bool {
    // Try to resolve hostname
    _, err := net.LookupHost(hostname)
    return err == nil
}

func (v *Validator) isPortOpen(hostname string, port int) bool {
    // Try to connect to port
    conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", hostname, port), 5*time.Second)
    if err != nil {
        return false
    }
    defer conn.Close()
    return true
}

func (v *Validator) validateAuthMethod(auth *spookytypesmachines.Authentication, index int) error {
    if auth.Method == "" {
        return fmt.Errorf("machine[%d]: authentication method is required", index)
    }
    
    validMethods := []string{"ssh_key", "password", "certificate"}
    isValid := false
    for _, method := range validMethods {
        if auth.Method == method {
            isValid = true
            break
        }
    }
    
    if !isValid {
        return fmt.Errorf("machine[%d]: invalid authentication method '%s'", index, auth.Method)
    }
    
    return nil
}

func (v *Validator) validateAuthCredentials(auth *spookytypesmachines.Authentication, index int) error {
    switch auth.Method {
    case "ssh_key":
        return v.validateSSHKey(auth, index)
    case "password":
        return v.validatePassword(auth, index)
    case "certificate":
        return v.validateCertificate(auth, index)
    default:
        return fmt.Errorf("machine[%d]: unsupported authentication method '%s'", index, auth.Method)
    }
}

func (v *Validator) validateSSHKey(auth *spookytypesmachines.Authentication, index int) error {
    if auth.KeyPath == "" {
        return fmt.Errorf("machine[%d]: SSH key path is required for ssh_key authentication", index)
    }
    
    // Check if key file exists
    if _, err := os.Stat(auth.KeyPath); os.IsNotExist(err) {
        return fmt.Errorf("machine[%d]: SSH key file '%s' does not exist", index, auth.KeyPath)
    }
    
    return nil
}

func (v *Validator) validatePassword(auth *spookytypesmachines.Authentication, index int) error {
    if auth.Password == "" {
        return fmt.Errorf("machine[%d]: password is required for password authentication", index)
    }
    
    return nil
}

func (v *Validator) validateCertificate(auth *spookytypesmachines.Authentication, index int) error {
    if auth.CertificatePath == "" {
        return fmt.Errorf("machine[%d]: certificate path is required for certificate authentication", index)
    }
    
    // Check if certificate file exists
    if _, err := os.Stat(auth.CertificatePath); os.IsNotExist(err) {
        return fmt.Errorf("machine[%d]: certificate file '%s' does not exist", index, auth.CertificatePath)
    }
    
    return nil
}
```

### Phase 3: Extract Result Processing (Day 3)

#### Extract Result Processing Logic
```go
func (v *Validator) processValidationResult(result *spookytypesmachines.ValidationResult, errors, warnings []string) {
    // Add errors and warnings
    result.Errors = errors
    result.Warnings = warnings
    
    // Determine overall validity
    result.Valid = len(errors) == 0
    
    // Set validation summary
    result.Summary = v.generateValidationSummary(result)
}

func (v *Validator) generateValidationSummary(result *spookytypesmachines.ValidationResult) string {
    if result.Valid {
        if len(result.Warnings) == 0 {
            return "All machines are valid"
        } else {
            return fmt.Sprintf("All machines are valid with %d warnings", len(result.Warnings))
        }
    } else {
        return fmt.Sprintf("Validation failed with %d errors and %d warnings", len(result.Errors), len(result.Warnings))
    }
}
```

### Phase 4: Refactored Main Function (Day 4)

#### Final Refactored Function
```go
func (v *Validator) validateMachinesInventory(inventory *spookytypesmachines.MachinesInventory) (*spookytypesmachines.ValidationResult, error) {
    // Validate input
    if err := v.validateInventoryInput(inventory); err != nil {
        return nil, err
    }
    
    // Initialize result
    result := v.initializeValidationResult()
    
    // Validate individual machines
    errors, warnings := v.validateIndividualMachines(inventory.Machines)
    
    // Process validation result
    v.processValidationResult(result, errors, warnings)
    
    return result, nil
}
```

## Complexity Reduction Metrics

### Before Refactoring
- **Cyclomatic Complexity:** 9
- **Lines of Code:** ~80
- **Nesting Levels:** 4-5
- **Responsibilities:** 4 (validation, processing, result creation, error handling)

### After Refactoring
- **Cyclomatic Complexity:** 3 (main function) + 1-3 per helper function
- **Lines of Code:** ~200 (distributed across multiple focused functions)
- **Nesting Levels:** 1-2
- **Responsibilities:** 1 per function (single responsibility principle)

## Testing Strategy

### Unit Tests for Extracted Functions
```go
func TestValidateInventoryInput(t *testing.T) {
    tests := []struct {
        name      string
        inventory *spookytypesmachines.MachinesInventory
        wantErr   bool
    }{
        {
            name:      "nil inventory",
            inventory: nil,
            wantErr:   true,
        },
        {
            name: "nil machines",
            inventory: &spookytypesmachines.MachinesInventory{
                Machines: nil,
            },
            wantErr: true,
        },
        {
            name: "empty machines",
            inventory: &spookytypesmachines.MachinesInventory{
                Machines: []*spookytypesmachines.Machine{},
            },
            wantErr: true,
        },
        {
            name: "valid inventory",
            inventory: &spookytypesmachines.MachinesInventory{
                Machines: []*spookytypesmachines.Machine{
                    {Name: "test", Hostname: "test.example.com"},
                },
            },
            wantErr: false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := validateInventoryInput(tt.inventory)
            if (err != nil) != tt.wantErr {
                t.Errorf("validateInventoryInput() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}

func TestValidateRequiredFields(t *testing.T) {
    tests := []struct {
        name    string
        machine *spookytypesmachines.Machine
        index   int
        wantErr bool
    }{
        {
            name: "missing name",
            machine: &spookytypesmachines.Machine{
                Name:     "",
                Hostname: "test.example.com",
                User:     "admin",
            },
            index:   0,
            wantErr: true,
        },
        {
            name: "missing hostname",
            machine: &spookytypesmachines.Machine{
                Name:     "test",
                Hostname: "",
                User:     "admin",
            },
            index:   0,
            wantErr: true,
        },
        {
            name: "missing user",
            machine: &spookytypesmachines.Machine{
                Name:     "test",
                Hostname: "test.example.com",
                User:     "",
            },
            index:   0,
            wantErr: true,
        },
        {
            name: "valid machine",
            machine: &spookytypesmachines.Machine{
                Name:     "test",
                Hostname: "test.example.com",
                User:     "admin",
            },
            index:   0,
            wantErr: false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := validateRequiredFields(tt.machine, tt.index)
            if (err != nil) != tt.wantErr {
                t.Errorf("validateRequiredFields() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}

func TestIsValidHostname(t *testing.T) {
    tests := []struct {
        name    string
        hostname string
        want    bool
    }{
        {
            name:     "valid hostname",
            hostname: "test.example.com",
            want:     true,
        },
        {
            name:     "empty hostname",
            hostname: "",
            want:     false,
        },
        {
            name:     "invalid characters",
            hostname: "test@example.com",
            want:     false,
        },
        {
            name:     "starts with dash",
            hostname: "-test.example.com",
            want:     false,
        },
        {
            name:     "ends with dash",
            hostname: "test.example.com-",
            want:     false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := isValidHostname(tt.hostname)
            if got != tt.want {
                t.Errorf("isValidHostname() = %v, want %v", got, tt.want)
            }
        })
    }
}

func TestValidateAuthMethod(t *testing.T) {
    tests := []struct {
        name    string
        auth    *spookytypesmachines.Authentication
        index   int
        wantErr bool
    }{
        {
            name: "empty method",
            auth: &spookytypesmachines.Authentication{
                Method: "",
            },
            index:   0,
            wantErr: true,
        },
        {
            name: "invalid method",
            auth: &spookytypesmachines.Authentication{
                Method: "invalid",
            },
            index:   0,
            wantErr: true,
        },
        {
            name: "valid ssh_key",
            auth: &spookytypesmachines.Authentication{
                Method: "ssh_key",
            },
            index:   0,
            wantErr: false,
        },
        {
            name: "valid password",
            auth: &spookytypesmachines.Authentication{
                Method: "password",
            },
            index:   0,
            wantErr: false,
        },
        {
            name: "valid certificate",
            auth: &spookytypesmachines.Authentication{
                Method: "certificate",
            },
            index:   0,
            wantErr: false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := validateAuthMethod(tt.auth, tt.index)
            if (err != nil) != tt.wantErr {
                t.Errorf("validateAuthMethod() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

### Integration Tests
```go
func TestValidateMachinesInventoryIntegration(t *testing.T) {
    tests := []struct {
        name      string
        inventory *spookytypesmachines.MachinesInventory
        want      *spookytypesmachines.ValidationResult
        wantErr   bool
    }{
        {
            name: "valid inventory",
            inventory: &spookytypesmachines.MachinesInventory{
                Machines: []*spookytypesmachines.Machine{
                    {
                        Name:     "test1",
                        Hostname: "test1.example.com",
                        User:     "admin",
                        Authentication: &spookytypesmachines.Authentication{
                            Method: "ssh_key",
                            KeyPath: "/path/to/key",
                        },
                    },
                },
            },
            want: &spookytypesmachines.ValidationResult{
                Valid: true,
            },
            wantErr: false,
        },
        {
            name: "invalid inventory",
            inventory: &spookytypesmachines.MachinesInventory{
                Machines: []*spookytypesmachines.Machine{
                    {
                        Name:     "",
                        Hostname: "test.example.com",
                        User:     "admin",
                    },
                },
            },
            want: &spookytypesmachines.ValidationResult{
                Valid: false,
            },
            wantErr: false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := validateMachinesInventory(tt.inventory)
            if (err != nil) != tt.wantErr {
                t.Errorf("validateMachinesInventory() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if result.Valid != tt.want.Valid {
                t.Errorf("validateMachinesInventory() valid = %v, want %v", result.Valid, tt.want.Valid)
            }
        })
    }
}
```

## Implementation Timeline

### Day 1: Extract Input Validation
- [ ] Extract `validateInventoryInput`
- [ ] Extract `initializeValidationResult`
- [ ] Add unit tests for input validation functions
- [ ] Verify input validation works correctly

### Day 2: Extract Machine Validation
- [ ] Extract `validateIndividualMachines`
- [ ] Extract `validateMachine`
- [ ] Extract `validateRequiredFields`
- [ ] Extract `validateFieldFormats`
- [ ] Extract `validateConnectivity`
- [ ] Extract `validateAuthentication`
- [ ] Extract validation helper functions
- [ ] Add unit tests for machine validation functions
- [ ] Verify machine validation works correctly

### Day 3: Extract Result Processing
- [ ] Extract `processValidationResult`
- [ ] Extract `generateValidationSummary`
- [ ] Add unit tests for result processing functions
- [ ] Verify result processing works correctly

### Day 4: Complete Refactoring
- [ ] Refactor main `validateMachinesInventory` function
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
- [ ] Easy to modify individual validation components
- [ ] Clear function names and purposes
- [ ] Well-documented functions
- [ ] Consistent validation patterns

## Risk Mitigation

### Potential Risks
1. **Validation logic changes** - May affect validation behavior
2. **Connectivity checks** - May affect network-dependent validations
3. **Authentication validation** - May affect security-related validations

### Mitigation Strategies
1. **Comprehensive testing** - Test all validation scenarios and edge cases
2. **Validation consistency** - Ensure validation behavior remains consistent
3. **Network isolation** - Mock network calls for testing
- [ ] **Gradual migration** - Implement changes incrementally with validation

## Verification

### Complexity Verification
```bash
# Verify complexity reduction
gocyclo internal/machines/validator.go | grep validateMachinesInventory
```

### Functionality Verification
```bash
# Test inventory validation
go test ./internal/machines -run TestValidateMachinesInventory
```

### Performance Verification
- [ ] No performance regression
- [ ] Memory usage remains stable
- [ ] Validation performance remains acceptable

## Conclusion

This refactoring will reduce the cyclomatic complexity of `validateMachinesInventory` from 9 to approximately 3, while improving code maintainability, testability, and readability. The extracted functions will follow single responsibility principles and be easily testable in isolation.

**Expected Outcome:** A clean, maintainable function with clear separation of concerns and comprehensive test coverage for machines inventory validation operations.
