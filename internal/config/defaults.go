package config

const (
	// DefaultSSHPort is the default SSH port if not specified
	DefaultSSHPort = 22

	// DefaultTimeout is the default timeout for SSH connections in seconds
	DefaultTimeout = 30

	// DefaultPasswordLength is the default length for generated passwords
	DefaultPasswordLength = 25

	// MaxKeyDirectories is the maximum number of key directories per day
	MaxKeyDirectories = 1000
)

// SetDefaults applies default values to a configuration
func SetDefaults(config *Config) {
	for i := range config.Servers {
		if config.Servers[i].Port == 0 {
			config.Servers[i].Port = DefaultSSHPort
		}
	}

	for i := range config.Actions {
		if config.Actions[i].Timeout == 0 {
			config.Actions[i].Timeout = DefaultTimeout
		}
	}
}
