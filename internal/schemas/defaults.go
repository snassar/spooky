package schemas

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

// DefaultConfigGenerator extracts default values from struct tags and creates default instances
type DefaultConfigGenerator struct{}

// NewDefaultConfigGenerator creates a new default config generator
func NewDefaultConfigGenerator() *DefaultConfigGenerator {
	return &DefaultConfigGenerator{}
}

// GetDefaultSpookyConfig returns a default SpookyV1 configuration
func (dcg *DefaultConfigGenerator) GetDefaultSpookyConfig() *SpookyV1 {
	return &SpookyV1{
		SSH:      *dcg.GetDefaultSpookySSH(),
		Security: *dcg.GetDefaultSpookySecurity(),
		Age:      *dcg.GetDefaultSpookyAge(),
		Logging:  *dcg.GetDefaultSpookyLogging(),
	}
}

// GetDefaultSpookySSH returns a default SpookySSHV1 configuration
func (dcg *DefaultConfigGenerator) GetDefaultSpookySSH() *SpookySSHV1 {
	return &SpookySSHV1{
		Timeout:            dcg.extractIntDefault(reflect.TypeOf(SpookySSHV1{}).Field(0), 30),
		KeepaliveInterval:  dcg.extractIntDefault(reflect.TypeOf(SpookySSHV1{}).Field(1), 60),
		KeepaliveCount:     dcg.extractIntDefault(reflect.TypeOf(SpookySSHV1{}).Field(2), 3),
		KeyScanTimeout:     dcg.extractIntDefault(reflect.TypeOf(SpookySSHV1{}).Field(3), 10),
		KnownHostsStrict:   dcg.extractBoolDefault(reflect.TypeOf(SpookySSHV1{}).Field(4), true),
		ConnectionPoolSize: dcg.extractIntDefault(reflect.TypeOf(SpookySSHV1{}).Field(5), 10),

		// Proxy configuration (no defaults)
		ProxyCommand: "",
		ProxyJump:    "",

		// Compression configuration
		Compression:      dcg.extractBoolDefault(reflect.TypeOf(SpookySSHV1{}).Field(7), false),
		CompressionLevel: dcg.extractIntDefault(reflect.TypeOf(SpookySSHV1{}).Field(8), 6),

		// TCP keepalive configuration
		TCPKeepAlive:              dcg.extractBoolDefault(reflect.TypeOf(SpookySSHV1{}).Field(9), true),
		TCPKeepAliveCount:         dcg.extractIntDefault(reflect.TypeOf(SpookySSHV1{}).Field(10), 3),
		TCPKeepAliveIdle:          60 * time.Second,
		TCPKeepAliveInterval:      10 * time.Second,
		TCPKeepAliveProbeInterval: 5 * time.Second,
	}
}

// GetDefaultSpookySecurity returns a default SpookySecurityV1 configuration
func (dcg *DefaultConfigGenerator) GetDefaultSpookySecurity() *SpookySecurityV1 {
	return &SpookySecurityV1{
		AllowUnsafeCommands: dcg.extractBoolDefault(reflect.TypeOf(SpookySecurityV1{}).Field(0), false),
		AuditLogging:        dcg.extractBoolDefault(reflect.TypeOf(SpookySecurityV1{}).Field(1), true),
	}
}

// GetDefaultSpookyAge returns a default SpookyAgeV1 configuration
func (dcg *DefaultConfigGenerator) GetDefaultSpookyAge() *SpookyAgeV1 {
	return &SpookyAgeV1{
		Identities: "", // No default path
		Recipients: "", // No default path
	}
}

// GetDefaultSpookyLogging returns a default SpookyLoggingV1 configuration
func (dcg *DefaultConfigGenerator) GetDefaultSpookyLogging() *SpookyLoggingV1 {
	return &SpookyLoggingV1{
		Level:      dcg.extractStringDefault(reflect.TypeOf(SpookyLoggingV1{}).Field(0), "info"),
		Format:     dcg.extractStringDefault(reflect.TypeOf(SpookyLoggingV1{}).Field(1), "json"),
		Output:     dcg.extractStringDefault(reflect.TypeOf(SpookyLoggingV1{}).Field(2), "stderr"),
		FilePath:   "", // Will be computed based on OS
		FilePerms:  dcg.extractStringDefault(reflect.TypeOf(SpookyLoggingV1{}).Field(4), "0o644"),
		FileAppend: dcg.extractBoolDefault(reflect.TypeOf(SpookyLoggingV1{}).Field(5), true),
	}
}

// GetDefaultProjectConfig returns a default ProjectV1 configuration
func (dcg *DefaultConfigGenerator) GetDefaultProjectConfig() *ProjectV1 {
	return &ProjectV1{
		RunMaxParallel:          dcg.extractIntDefault(reflect.TypeOf(ProjectV1{}).Field(2), 10),
		RunDryRunDefault:        dcg.extractBoolDefault(reflect.TypeOf(ProjectV1{}).Field(3), false),
		RunValidateBeforeRun:    dcg.extractBoolDefault(reflect.TypeOf(ProjectV1{}).Field(4), true),
		RunBackupBeforeChanges:  dcg.extractBoolDefault(reflect.TypeOf(ProjectV1{}).Field(5), false),
		FactsTimeout:            dcg.extractIntDefault(reflect.TypeOf(ProjectV1{}).Field(6), 30),
		FactsParallelCollection: dcg.extractIntDefault(reflect.TypeOf(ProjectV1{}).Field(7), 10),
		FactsRetryAttempts:      dcg.extractIntDefault(reflect.TypeOf(ProjectV1{}).Field(8), 3),
		FactsRetryDelay:         dcg.extractIntDefault(reflect.TypeOf(ProjectV1{}).Field(9), 5),
	}
}

// GetDefaultMachinesConfig returns a default MachinesV1 configuration
func (dcg *DefaultConfigGenerator) GetDefaultMachinesConfig() *MachinesV1 {
	return &MachinesV1{
		Machine: []MachinesMachineV1{},
		Group:   []MachinesGroupV1{},
	}
}

// GetDefaultActionsConfig returns a default ActionsV1 configuration
func (dcg *DefaultConfigGenerator) GetDefaultActionsConfig() *ActionsV1 {
	return &ActionsV1{
		Action: []ActionsActionV1{},
	}
}

// GetDefaultVariablesConfig returns a default VariablesV1 configuration
func (dcg *DefaultConfigGenerator) GetDefaultVariablesConfig() *VariablesV1 {
	return &VariablesV1{
		Variable: []VariablesVariableV1{},
	}
}

// GetDefaultLoggingConfig returns a default LoggingV1 configuration
func (dcg *DefaultConfigGenerator) GetDefaultLoggingConfig() *LoggingV1 {
	return &LoggingV1{
		Level:                          "info",
		Format:                         "json",
		Output:                         "stderr",
		FilePermissions:                "0o644",
		FileAppend:                     true,
		StructuredTimestampEnabled:     true,
		StructuredTimestampFormat:      "RFC3339",
		StructuredTimestampTimezone:    "UTC",
		StructuredLevelKey:             "level",
		StructuredMessageKey:           "message",
		StructuredErrorKey:             "error",
		PerformanceBufferEnabled:       false,
		PerformanceBufferSize:          4096,
		PerformanceBufferFlushInterval: "1s",
		PerformanceAsyncEnabled:        false,
		PerformanceAsyncQueueSize:      1000,
		PerformanceAsyncWorkers:        1,
		PerformanceAsyncDropWhenFull:   false,
		RotationEnabled:                false,
		RotationMaxSize:                "100MB",
		RotationMaxAge:                 "30d",
		RotationMaxBackups:             5,
		RotationCompress:               true,
		RotationLocalTime:              false,
	}
}

// GetDefaultFactsConfig returns a default FactsV1 configuration
func (dcg *DefaultConfigGenerator) GetDefaultFactsConfig() *FactsV1 {
	return &FactsV1{
		BasicFacts: &BasicFactsV1{
			SystemFacts:   &SystemFactsV1{Facts: make(map[string]*FactV1)},
			HardwareFacts: &HardwareFactsV1{Facts: make(map[string]*FactV1)},
			NetworkFacts:  &NetworkFactsV1{Facts: make(map[string]*FactV1)},
			OSFacts:       &OSFactsV1{Facts: make(map[string]*FactV1)},
			UserFacts:     &UserFactsV1{Facts: make(map[string]*FactV1)},
			RuntimeFacts:  &RuntimeFactsV1{Facts: make(map[string]*FactV1)},
		},
		EnhancedFacts: &EnhancedFactsV1{Facts: make(map[string]*FactV1)},
		CustomFacts:   &CustomFactsV1{Facts: make(map[string]*FactV1)},
	}
}

// extractIntDefault extracts an integer default value from a struct field tag
func (dcg *DefaultConfigGenerator) extractIntDefault(field reflect.StructField, fallback int) int {
	if defaultStr := dcg.extractDefaultTag(field); defaultStr != "" {
		if val, err := strconv.Atoi(defaultStr); err == nil {
			return val
		}
	}
	return fallback
}

// extractBoolDefault extracts a boolean default value from a struct field tag
func (dcg *DefaultConfigGenerator) extractBoolDefault(field reflect.StructField, fallback bool) bool {
	if defaultStr := dcg.extractDefaultTag(field); defaultStr != "" {
		if val, err := strconv.ParseBool(defaultStr); err == nil {
			return val
		}
	}
	return fallback
}

// extractStringDefault extracts a string default value from a struct field tag
func (dcg *DefaultConfigGenerator) extractStringDefault(field reflect.StructField, fallback string) string {
	if defaultStr := dcg.extractDefaultTag(field); defaultStr != "" {
		return defaultStr
	}
	return fallback
}

// extractDefaultTag extracts the default value from a struct field tag
func (dcg *DefaultConfigGenerator) extractDefaultTag(field reflect.StructField) string {
	tag := field.Tag.Get("json")
	if tag == "" {
		return ""
	}

	// Look for default tag in the same struct tag
	defaultTag := field.Tag.Get("default")
	if defaultTag != "" {
		return defaultTag
	}

	return ""
}

// ToHCL converts the default config to HCL string
func (dcg *DefaultConfigGenerator) ToHCL(config interface{}) (string, error) {
	// Convert the config to HCL using the HCL library
	hclFile := hclwrite.NewEmptyFile()

	// Convert the config to cty.Value
	ctyValue, err := dcg.structToCty(config)
	if err != nil {
		return "", fmt.Errorf("failed to convert config to cty.Value: %v", err)
	}

	// Create the root body
	rootBody := hclFile.Body()

	// Get the config type name for the root block
	configType := reflect.TypeOf(config)
	if configType.Kind() == reflect.Ptr {
		configType = configType.Elem()
	}

	// Convert to HCL block
	err = dcg.ctyValueToHCL(ctyValue, rootBody, configType.Name())
	if err != nil {
		return "", fmt.Errorf("failed to convert cty.Value to HCL: %v", err)
	}

	return string(hclFile.Bytes()), nil
}

// structToCty converts a Go struct to cty.Value
func (dcg *DefaultConfigGenerator) structToCty(config interface{}) (cty.Value, error) {
	v := reflect.ValueOf(config)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return cty.NilVal, fmt.Errorf("config must be a struct, got %s", v.Kind())
	}

	// Convert struct to map[string]cty.Value
	valueMap := make(map[string]cty.Value)
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		// Get the JSON tag name
		jsonTag := fieldType.Tag.Get("json")
		if jsonTag == "" || jsonTag == "-" {
			continue
		}

		// Handle comma-separated JSON tags (e.g., "basic_facts,omitempty")
		if commaIndex := strings.Index(jsonTag, ","); commaIndex != -1 {
			jsonTag = jsonTag[:commaIndex]
		}

		// Convert field value to cty.Value
		ctyVal, err := dcg.fieldToCty(field)
		if err != nil {
			return cty.NilVal, fmt.Errorf("failed to convert field %s: %v", jsonTag, err)
		}

		if !ctyVal.IsNull() {
			valueMap[jsonTag] = ctyVal
		}
	}

	return cty.ObjectVal(valueMap), nil
}

// fieldToCty converts a struct field to cty.Value
func (dcg *DefaultConfigGenerator) fieldToCty(field reflect.Value) (cty.Value, error) {
	switch field.Kind() {
	case reflect.Ptr:
		// Handle pointers
		if field.IsNil() {
			return cty.NilVal, nil
		}
		return dcg.fieldToCty(field.Elem())
	case reflect.String:
		return cty.StringVal(field.String()), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return cty.NumberIntVal(field.Int()), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return cty.NumberUIntVal(field.Uint()), nil
	case reflect.Float32, reflect.Float64:
		return cty.NumberFloatVal(field.Float()), nil
	case reflect.Bool:
		return cty.BoolVal(field.Bool()), nil
	case reflect.Struct:
		// Handle nested structs
		if field.CanInterface() {
			return dcg.structToCty(field.Interface())
		}
		return cty.NilVal, fmt.Errorf("cannot convert unexported struct field")
	case reflect.Slice:
		// Handle slices
		if field.Len() == 0 {
			return cty.ListValEmpty(cty.DynamicPseudoType), nil
		}

		values := make([]cty.Value, field.Len())
		for i := 0; i < field.Len(); i++ {
			val, err := dcg.fieldToCty(field.Index(i))
			if err != nil {
				return cty.NilVal, err
			}
			values[i] = val
		}
		return cty.ListVal(values), nil
	case reflect.Map:
		// Handle maps
		if field.Len() == 0 {
			return cty.MapValEmpty(cty.DynamicPseudoType), nil
		}

		valueMap := make(map[string]cty.Value)
		iter := field.MapRange()
		for iter.Next() {
			key := iter.Key()
			val := iter.Value()

			if key.Kind() != reflect.String {
				return cty.NilVal, fmt.Errorf("map keys must be strings")
			}

			ctyVal, err := dcg.fieldToCty(val)
			if err != nil {
				return cty.NilVal, err
			}
			valueMap[key.String()] = ctyVal
		}
		// Create a map with dynamic type to allow mixed value types
		return cty.MapVal(valueMap), nil
	case reflect.Interface:
		// Handle interface{} types
		if field.IsNil() {
			return cty.NilVal, nil
		}

		// Get the concrete value
		concreteValue := field.Elem()
		return dcg.fieldToCty(concreteValue)
	default:
		return cty.NilVal, fmt.Errorf("unsupported field type: %s", field.Kind())
	}
}

// ctyValueToHCL converts a cty.Value to HCL and writes it to the body
func (dcg *DefaultConfigGenerator) ctyValueToHCL(value cty.Value, body *hclwrite.Body, blockName string) error {
	if !value.IsKnown() {
		return fmt.Errorf("cannot convert unknown value to HCL")
	}

	if value.IsNull() {
		return nil
	}

	valueType := value.Type()
	if valueType.IsObjectType() {
		// Create a block for objects
		block := body.AppendNewBlock(blockName, nil)
		blockBody := block.Body()

		// Add each field to the block
		for key, val := range value.AsValueMap() {
			if !val.IsNull() {
				err := dcg.ctyValueToHCL(val, blockBody, key)
				if err != nil {
					return fmt.Errorf("failed to convert field %s: %v", key, err)
				}
			}
		}

	} else if valueType.IsListType() {
		// Handle lists by creating blocks for each element
		values := value.AsValueSlice()
		for _, val := range values {
			if !val.IsNull() {
				err := dcg.ctyValueToHCL(val, body, blockName)
				if err != nil {
					return fmt.Errorf("failed to convert list element: %v", err)
				}
			}
		}

	} else if valueType.IsMapType() {
		// Handle maps by creating blocks for each key-value pair
		valueMap := value.AsValueMap()
		for key, val := range valueMap {
			if !val.IsNull() {
				// For maps, we create a block with the key as the label
				block := body.AppendNewBlock(blockName, []string{key})
				blockBody := block.Body()

				err := dcg.ctyValueToHCL(val, blockBody, "")
				if err != nil {
					return fmt.Errorf("failed to convert map value for key %s: %v", key, err)
				}
			}
		}

	} else {
		// For primitive types, set as attribute
		body.SetAttributeValue(blockName, value)
	}

	return nil
}

// GenerateProjectConfigFromStructs generates a project configuration HCL string from Go structs
func GenerateProjectConfigFromStructs(name, description string) string {
	generator := NewDefaultConfigGenerator()
	config := generator.GetDefaultProjectConfig()

	// Override with provided values
	config.Name = name
	config.Description = description

	hcl, err := generator.ToHCL(config)
	if err != nil {
		// This is a critical bug - struct generation should never fail
		panic(fmt.Sprintf("failed to generate project config from struct: %v", err))
	}

	// The ToHCL method generates "ProjectV1" as the block name, but we need "project"
	// So we need to replace the block name in the generated HCL
	hcl = strings.Replace(hcl, "ProjectV1", "project", 1)

	// Convert from project { name = "..." } to project "..." { }
	// Remove the name attribute and add it as a quoted label
	hcl = strings.Replace(hcl, fmt.Sprintf(`  name                      = "%s"`, name), "", 1)
	hcl = strings.Replace(hcl, fmt.Sprintf(`  name = "%s"`, name), "", 1)
	hcl = strings.Replace(hcl, "project {", fmt.Sprintf(`project "%s" {`, name), 1)

	return hcl
}

// GenerateMachinesConfigFromStructs generates a machines configuration HCL string from Go structs
func GenerateMachinesConfigFromStructs() string {
	generator := NewDefaultConfigGenerator()
	config := generator.GetDefaultMachinesConfig()

	hcl, err := generator.ToHCL(config)
	if err != nil {
		return fmt.Sprintf("machines {\n  # Error generating from struct: %v\n}", err)
	}

	// The ToHCL method generates "MachinesV1" as the block name, but we need "machines"
	// So we need to replace the block name in the generated HCL
	hcl = strings.Replace(hcl, "MachinesV1", "machines", 1)

	return hcl
}

// GenerateActionsConfigFromStructs generates an actions configuration HCL string from Go structs
func GenerateActionsConfigFromStructs() string {
	generator := NewDefaultConfigGenerator()
	config := generator.GetDefaultActionsConfig()

	hcl, err := generator.ToHCL(config)
	if err != nil {
		return fmt.Sprintf("actions {\n  # Error generating from struct: %v\n}", err)
	}

	// The ToHCL method generates "ActionsV1" as the block name, but we need "actions"
	// So we need to replace the block name in the generated HCL
	hcl = strings.Replace(hcl, "ActionsV1", "actions", 1)

	return hcl
}

// GenerateVariablesConfigFromStructs generates a variables configuration HCL string from Go structs
func GenerateVariablesConfigFromStructs() string {
	generator := NewDefaultConfigGenerator()
	config := generator.GetDefaultVariablesConfig()

	hcl, err := generator.ToHCL(config)
	if err != nil {
		return fmt.Sprintf("variables {\n  # Error generating from struct: %v\n}", err)
	}

	// The ToHCL method generates "VariablesV1" as the block name, but we need "variables"
	// So we need to replace the block name in the generated HCL
	hcl = strings.Replace(hcl, "VariablesV1", "variables", 1)

	return hcl
}
