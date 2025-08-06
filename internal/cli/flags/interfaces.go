package flags

// FlagsManager defines the interface for flag management
type FlagsManager interface {
	SetGlobalFlags(flags map[string]interface{}) error
	SetCommandFlags(commandName string, flags map[string]interface{}) error
	GetGlobalFlags() map[string]interface{}
	GetCommandFlags(commandName string) map[string]interface{}
}

// FlagsParser defines the interface for flag parsing
type FlagsParser interface {
	ParseFlags(cmd interface{}) error
	ValidateFlags(flags map[string]interface{}) error
	GetFlagValue(name string) interface{}
}
