package schemas

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
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
		FilePerms:  dcg.extractStringDefault(reflect.TypeOf(SpookyLoggingV1{}).Field(4), "0644"),
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
		FilePermissions:                "0644",
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
		Encrypted: dcg.extractBoolDefault(reflect.TypeOf(FactsV1{}).Field(3), false),
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
	// For now, return a simple HCL representation
	// TODO: Implement proper HCL generation using the HCL library
	return dcg.generateSimpleHCL(config), nil
}

// generateSimpleHCL generates a simple HCL representation
func (dcg *DefaultConfigGenerator) generateSimpleHCL(config interface{}) string {
	var result strings.Builder

	v := reflect.ValueOf(config)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return ""
	}

	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		// Get the JSON tag name
		jsonTag := fieldType.Tag.Get("json")
		if jsonTag == "" || jsonTag == "-" {
			continue
		}

		// Handle nested structs
		if field.Kind() == reflect.Struct {
			if !field.IsZero() {
				result.WriteString(fmt.Sprintf("\n%s {\n", jsonTag))
				result.WriteString(dcg.generateSimpleHCL(field.Interface()))
				result.WriteString("}\n")
			}
			continue
		}

		// Handle slices
		if field.Kind() == reflect.Slice {
			if field.Len() > 0 {
				result.WriteString(fmt.Sprintf("\n%s {\n", jsonTag))
				for j := 0; j < field.Len(); j++ {
					result.WriteString(dcg.generateSimpleHCL(field.Index(j).Interface()))
				}
				result.WriteString("}\n")
			}
			continue
		}

		// Handle basic types
		if !field.IsZero() {
			switch field.Kind() {
			case reflect.String:
				result.WriteString(fmt.Sprintf("  %s = \"%s\"\n", jsonTag, field.String()))
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				result.WriteString(fmt.Sprintf("  %s = %d\n", jsonTag, field.Int()))
			case reflect.Bool:
				result.WriteString(fmt.Sprintf("  %s = %t\n", jsonTag, field.Bool()))
			}
		}
	}

	return result.String()
}
