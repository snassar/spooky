package ssh

import (
	"time"

	spookyinterfaces "spooky/internal/interfaces"
	spookysshacting "spooky/internal/ssh/acting"
	spookysshauthentication "spooky/internal/ssh/authentication"
	spookysshclient "spooky/internal/ssh/client"
	spookysshconnectionpool "spooky/internal/ssh/connection_pool"
	spookysshkeys "spooky/internal/ssh/keys"
	spookytypes "spooky/internal/types"
)

// NewManagerWithConfig creates a new SSH manager with custom configuration
func NewManagerWithConfig(
	config *spookytypes.SSHClientConfig,
	logger spookyinterfaces.Logger,
) *Manager {
	// Create subpackage managers with default implementations
	clientManager := spookysshclient.NewManager(
		config.ClientConfig,
		nil, // connectionManager - will be implemented
		nil, // executionManager - will be implemented
		nil, // hostKeyManager - will be implemented
		logger,
	)

	authenticationManager := spookysshauthentication.NewManager(
		config.AuthenticationConfig,
		logger,
	)

	connectionPoolManager := spookysshconnectionpool.NewManager(
		config.PoolConfig,
		logger,
	)

	actingManager := spookysshacting.NewManager(
		config.ActingConfig,
		logger,
	)

	keyManager := spookysshkeys.NewManager(
		config.KeysConfig,
		logger,
	)

	return NewManagerWithDependencies(
		config,
		clientManager,
		authenticationManager,
		connectionPoolManager,
		actingManager,
		keyManager,
		logger,
	)
}

// NewManagerWithDependencies creates a new SSH manager with custom dependencies
func NewManagerWithDependencies(
	config *spookytypes.SSHClientConfig,
	clientManager spookyinterfaces.ClientManager,
	authenticationManager spookyinterfaces.AuthenticationEngine,
	connectionPoolManager spookyinterfaces.ConnectionPool,
	actingManager spookyinterfaces.ActingEngine,
	keyManager spookyinterfaces.SSHKeyManager,
	logger spookyinterfaces.Logger,
) *Manager {
	return &Manager{
		config:                config,
		clientManager:         clientManager,
		authenticationManager: authenticationManager,
		connectionPoolManager: connectionPoolManager,
		actingManager:         actingManager,
		keyManager:            keyManager,
		logger:                logger,
	}
}

// NewDefaultManager creates a new SSH manager with sensible defaults
func NewDefaultManager(logger spookyinterfaces.Logger) *Manager {
	config := &spookytypes.SSHClientConfig{
		DefaultTimeout:          30 * time.Second,
		MaxConnections:          10,
		EnableConnectionPooling: true,
		ClientConfig: &spookytypes.ClientConfig{
			DefaultTimeout:  30 * time.Second,
			MaxRetries:      3,
			HostKeyChecking: true,
		},
		AuthenticationConfig: &spookytypes.AuthenticationConfig{
			Method:  spookytypes.AuthMethodKey,
			Timeout: 30 * time.Second,
		},
		PoolConfig: &spookytypes.PoolConfig{
			MaxConnections:      10,
			MaxIdleTime:         5 * time.Minute,
			ConnectionTimeout:   30 * time.Second,
			HealthCheckInterval: 1 * time.Minute,
		},
		ActingConfig: &spookytypes.ActingConfig{
			DefaultTimeout: 30 * time.Second,
			MaxParallel:    5,
			EnableRetries:  true,
			MaxRetries:     3,
			RetryDelay:     1 * time.Second,
		},
		KeysConfig: &spookytypes.KeysConfig{
			DefaultKeyPath: "~/.ssh/id_rsa",
			KeyCacheTTL:    1 * time.Hour,
			EnableCaching:  true,
		},
	}

	return NewManagerWithConfig(config, logger)
}
