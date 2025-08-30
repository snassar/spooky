package schemas

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"spooky/internal/hcl"
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
		Timeout:                   dcg.extractIntDefault(reflect.TypeOf(SpookySSHV1{}).Field(0), 30),
		KeepaliveInterval:         dcg.extractIntDefault(reflect.TypeOf(SpookySSHV1{}).Field(1), 60),
		KeepaliveCount:            dcg.extractIntDefault(reflect.TypeOf(SpookySSHV1{}).Field(2), 3),
		KeyScanTimeout:            dcg.extractIntDefault(reflect.TypeOf(SpookySSHV1{}).Field(3), 10),
		KnownHostsStrict:          dcg.extractBoolDefault(reflect.TypeOf(SpookySSHV1{}).Field(4), false), // Deprecated
		KnownHostsMode:            dcg.extractStringDefault(reflect.TypeOf(SpookySSHV1{}).Field(5), "accept-new"),
		ConnectionPoolSize:        dcg.extractIntDefault(reflect.TypeOf(SpookySSHV1{}).Field(6), 10),
		ProxyCommand:              "", // No default
		ProxyJump:                 "", // No default
		Compression:               dcg.extractBoolDefault(reflect.TypeOf(SpookySSHV1{}).Field(9), false),
		CompressionLevel:          dcg.extractIntDefault(reflect.TypeOf(SpookySSHV1{}).Field(10), 6),
		TCPKeepAlive:              dcg.extractBoolDefault(reflect.TypeOf(SpookySSHV1{}).Field(11), true),
		TCPKeepAliveCount:         dcg.extractIntDefault(reflect.TypeOf(SpookySSHV1{}).Field(12), 3),
		TCPKeepAliveIdle:          dcg.extractDurationDefault(reflect.TypeOf(SpookySSHV1{}).Field(13), 60*time.Second),
		TCPKeepAliveInterval:      dcg.extractDurationDefault(reflect.TypeOf(SpookySSHV1{}).Field(14), 10*time.Second),
		TCPKeepAliveProbeInterval: dcg.extractDurationDefault(reflect.TypeOf(SpookySSHV1{}).Field(15), 5*time.Second),
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
		Group: []MachinesGroupV1{},
		// Note: Individual machines are not included in the default config
		// as they should be added by users based on their specific needs
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

// extractDurationDefault extracts a duration default value from a struct field tag
func (dcg *DefaultConfigGenerator) extractDurationDefault(field reflect.StructField, fallback time.Duration) time.Duration {
	if defaultStr := dcg.extractDefaultTag(field); defaultStr != "" {
		if val, err := time.ParseDuration(defaultStr); err == nil {
			return val
		}
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

// GenerateProjectConfigFromStructs generates a project configuration HCL string from Go structs
func GenerateProjectConfigFromStructs(name, description string) string {
	generator := NewDefaultConfigGenerator()
	config := generator.GetDefaultProjectConfig()

	// Override with provided values
	config.Name = name
	config.Description = description

	hcl, err := hcl.GenerateConfigHCL(config, "project")
	if err != nil {
		// This is a critical bug - struct generation should never fail
		panic(fmt.Sprintf("failed to generate project config from struct: %v", err))
	}

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

	hcl, err := hcl.GenerateConfigHCL(config, "machines")
	if err != nil {
		return fmt.Sprintf("machines {\n  # Error generating from struct: %v\n}", err)
	}

	return hcl
}

// GenerateActionsConfigFromStructs generates an actions configuration HCL string from Go structs
func GenerateActionsConfigFromStructs() string {
	generator := NewDefaultConfigGenerator()
	config := generator.GetDefaultActionsConfig()

	hcl, err := hcl.GenerateConfigHCL(config, "actions")
	if err != nil {
		return fmt.Sprintf("actions {\n  # Error generating from struct: %v\n}", err)
	}
	return hcl
}

// GenerateVariablesConfigFromStructs generates a variables configuration HCL string from Go structs
func GenerateVariablesConfigFromStructs() string {
	generator := NewDefaultConfigGenerator()
	config := generator.GetDefaultVariablesConfig()

	hcl, err := hcl.GenerateConfigHCL(config, "variables")
	if err != nil {
		return fmt.Sprintf("variables {\n  # Error generating from struct: %v\n}", err)
	}

	return hcl
}
