package hcl

// GenerateConfigHCL is a convenience function for configuration structs
func GenerateConfigHCL(config interface{}, blockName string) (string, error) {
	generator := NewConfigGenerator()
	return generator.ToHCLWithBlockName(config, blockName)
}
