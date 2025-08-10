package config

import (
	"os"
	"path/filepath"
)

const (
	// DefaultSSHPort is the default SSH port if not specified
	DefaultSSHPort = 22

	// DefaultTimeout is the default timeout for SSH connections in seconds
	DefaultTimeout = 30

	// DefaultPasswordLength is the default length for generated passwords
	DefaultPasswordLength = 25

	// MaxKeyDirectories is the maximum number of key directories per day
	MaxKeyDirectories = 1000

	// Default facts collection timeout
	DefaultFactsTimeout = 30

	// Default facts cache TTL
	DefaultFactsCacheTTL = 3600

	// Default parallel facts collection
	DefaultParallelFactsCollection = 10

	// Default retry attempts for facts collection
	DefaultFactsRetryAttempts = 3

	// Default retry delay for facts collection
	DefaultFactsRetryDelay = 5

	// Default SSH keepalive interval
	DefaultSSHKeepaliveInterval = 60

	// Default SSH keepalive count
	DefaultSSHKeepaliveCount = 3

	// Default SSH key scan timeout
	DefaultSSHKeyScanTimeout = 10

	// Default SSH connection pool size
	DefaultSSHConnectionPoolSize = 10

	// Default template max size (1MB)
	DefaultTemplateMaxSize = 1048576

	// Default template timeout
	DefaultTemplateTimeout = 30

	// Default performance parallel limit
	DefaultParallelLimit = 10

	// Default max memory usage (512MB)
	DefaultMaxMemory = 512

	// Default GC interval (5 minutes)
	DefaultGCInterval = 300

	// Default backup retention
	DefaultBackupRetention = 7
)

// SetMachineDefaults applies default values to a machine
func SetMachineDefaults(machine *Machine) {
	if machine == nil {
		return
	}
	if machine.Port == 0 {
		machine.Port = DefaultSSHPort
	}
}

// SetActionDefaults applies default values to an action
func SetActionDefaults(action *Action) {
	if action == nil {
		return
	}
	if action.Timeout == 0 {
		action.Timeout = DefaultTimeout
	}
}

// SetInventoryDefaults applies default values to an inventory configuration
func SetInventoryDefaults(inventory *InventoryConfig) {
	if inventory == nil {
		return
	}
	for i := range inventory.Machines {
		SetMachineDefaults(&inventory.Machines[i])
	}
}

// SetActionsDefaults applies default values to an actions configuration
func SetActionsDefaults(actions *ActionsConfig) {
	if actions == nil {
		return
	}
	for i := range actions.Actions {
		SetActionDefaults(&actions.Actions[i])
	}
}

// SetProjectDefaults applies default values to a project configuration
func SetProjectDefaults(config *ProjectConfig) {
	if config == nil {
		return
	}

	// Set default file paths if not specified
	if config.InventoryFile == "" {
		config.InventoryFile = "inventory.hcl"
	}
	if config.ActionsFile == "" {
		config.ActionsFile = "actions.hcl"
	}

	// Set default timeouts
	if config.DefaultTimeout == 0 {
		config.DefaultTimeout = DefaultTimeout
	}

	// Set default configuration blocks if not present
	if config.Storage == nil {
		config.Storage = &StorageConfig{}
		SetStorageDefaults(config.Storage)
	}
	if config.Logging == nil {
		config.Logging = &LoggingConfig{}
		SetLoggingDefaults(config.Logging)
	}
	if config.SSH == nil {
		config.SSH = &SSHConfig{}
		SetSSHDefaults(config.SSH)
	}
	if config.Isolation == nil {
		config.Isolation = &IsolationConfig{}
		SetIsolationDefaults(config.Isolation)
	}
}

// SetGlobalConfigDefaults applies default values to a global configuration
func SetGlobalConfigDefaults(config *GlobalConfig) {
	if config == nil {
		return
	}

	// Set default configuration blocks if not present
	if config.Storage == nil {
		config.Storage = &StorageConfig{}
		SetStorageDefaults(config.Storage)
	}
	if config.Facts == nil {
		config.Facts = &FactsConfig{}
		SetFactsDefaults(config.Facts)
	}
	if config.SSH == nil {
		config.SSH = &SSHConfig{}
		SetSSHDefaults(config.SSH)
	}
	if config.Templates == nil {
		config.Templates = &TemplatesConfig{}
		SetTemplatesDefaults(config.Templates)
	}
	if config.Security == nil {
		config.Security = &SecurityConfig{}
		SetSecurityDefaults(config.Security)
	}
	if config.Age == nil {
		config.Age = &AgeConfig{}
		SetAgeDefaults(config.Age)
	}
	if config.Logging == nil {
		config.Logging = &LoggingConfig{}
		SetLoggingDefaults(config.Logging)
	}
	if config.Performance == nil {
		config.Performance = &PerformanceConfig{}
		SetPerformanceDefaults(config.Performance)
	}
}

// SetStorageDefaults applies default values to storage configuration
func SetStorageDefaults(storage *StorageConfig) {
	if storage == nil {
		return
	}
	if storage.Type == "" {
		storage.Type = "badgerdb"
	}
	if storage.Path == "" {
		storage.Path = getDefaultStoragePath()
	}
	if storage.BackupRetention == 0 {
		storage.BackupRetention = DefaultBackupRetention
	}
}

// SetFactsDefaults applies default values to facts configuration
func SetFactsDefaults(facts *FactsConfig) {
	if facts == nil {
		return
	}
	if facts.Timeout == 0 {
		facts.Timeout = DefaultFactsTimeout
	}
	if facts.CacheTTL == 0 {
		facts.CacheTTL = DefaultFactsCacheTTL
	}
	if facts.ParallelCollection == 0 {
		facts.ParallelCollection = DefaultParallelFactsCollection
	}
	if facts.RetryAttempts == 0 {
		facts.RetryAttempts = DefaultFactsRetryAttempts
	}
	if facts.RetryDelay == 0 {
		facts.RetryDelay = DefaultFactsRetryDelay
	}
}

// SetSSHDefaults applies default values to SSH configuration
func SetSSHDefaults(ssh *SSHConfig) {
	if ssh == nil {
		return
	}
	if ssh.Timeout == 0 {
		ssh.Timeout = DefaultTimeout
	}
	if ssh.KeepaliveInterval == 0 {
		ssh.KeepaliveInterval = DefaultSSHKeepaliveInterval
	}
	if ssh.KeepaliveCount == 0 {
		ssh.KeepaliveCount = DefaultSSHKeepaliveCount
	}
	if ssh.KeyScanTimeout == 0 {
		ssh.KeyScanTimeout = DefaultSSHKeyScanTimeout
	}
	if ssh.ConnectionPoolSize == 0 {
		ssh.ConnectionPoolSize = DefaultSSHConnectionPoolSize
	}
}

// SetTemplatesDefaults applies default values to templates configuration
func SetTemplatesDefaults(templates *TemplatesConfig) {
	if templates == nil {
		return
	}
	if templates.MaxSize == 0 {
		templates.MaxSize = DefaultTemplateMaxSize
	}
	if templates.Timeout == 0 {
		templates.Timeout = DefaultTemplateTimeout
	}
}

// SetSecurityDefaults applies default values to security configuration
func SetSecurityDefaults(security *SecurityConfig) {
	if security == nil {
		return
	}
	// Security defaults are mostly boolean flags that default to false
	// for safety reasons
}

// SetAgeDefaults applies default values to age configuration
func SetAgeDefaults(age *AgeConfig) {
	if age == nil {
		return
	}
	// Age defaults are mostly empty strings and false booleans
	// for security reasons
}

// SetLoggingDefaults applies default values to logging configuration
func SetLoggingDefaults(logging *LoggingConfig) {
	if logging == nil {
		return
	}
	if logging.Level == "" {
		logging.Level = "info"
	}
	if logging.Format == "" {
		logging.Format = "text"
	}
}

// SetPerformanceDefaults applies default values to performance configuration
func SetPerformanceDefaults(performance *PerformanceConfig) {
	if performance == nil {
		return
	}
	if performance.DefaultParallel == 0 {
		performance.DefaultParallel = DefaultParallelLimit
	}
	if performance.MaxMemory == 0 {
		performance.MaxMemory = DefaultMaxMemory
	}
	if performance.GCInterval == 0 {
		performance.GCInterval = DefaultGCInterval
	}
}

// SetIsolationDefaults applies default values to isolation configuration
func SetIsolationDefaults(isolation *IsolationConfig) {
	if isolation == nil {
		return
	}
	// Isolation defaults are mostly boolean flags that default to false
	// for compatibility reasons
}

// getDefaultStoragePath returns the default storage path
func getDefaultStoragePath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ".spooky"
	}
	return filepath.Join(homeDir, ".local", "state", "spooky")
}
