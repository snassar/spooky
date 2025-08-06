package ssh

import (
	"spooky/internal/logging"
	"spooky/internal/ssh/acting"
	"spooky/internal/ssh/authentication"
	"spooky/internal/ssh/client"
	"spooky/internal/ssh/connection_pool"
	"spooky/internal/ssh/keys"
	"spooky/internal/ssh/types"
	"time"
)

// NewManagerWithConfig creates a new SSH manager with custom configuration
func NewManagerWithConfig(
	config *types.Config,
	logger logging.Logger,
) *Manager {
	// Create subpackage managers with default implementations
	clientManager := client.NewManager(
		config.ClientConfig,
		nil, // connectionManager - will be implemented
		nil, // executionManager - will be implemented
		nil, // hostKeyManager - will be implemented
		logger,
	)

	authenticationManager := authentication.NewManager(
		config.AuthenticationConfig,
		logger,
	)

	connectionPoolManager := connection_pool.NewManager(
		config.PoolConfig,
		logger,
	)

	actingManager := acting.NewManager(
		config.ActingConfig,
		logger,
	)

	keyManager := keys.NewManager(
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
	config *types.Config,
	clientManager client.ClientManager,
	authenticationManager authentication.AuthenticationEngine,
	connectionPoolManager connection_pool.ConnectionPool,
	actingManager acting.ActingEngine,
	keyManager keys.SSHKeyManager,
	logger logging.Logger,
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
func NewDefaultManager(logger logging.Logger) *Manager {
	config := &types.Config{
		DefaultTimeout:          30 * time.Second,
		MaxConnections:          10,
		EnableConnectionPooling: true,
		ClientConfig: &types.ClientConfig{
			DefaultTimeout:  30 * time.Second,
			MaxRetries:      3,
			HostKeyChecking: true,
		},
		AuthenticationConfig: &types.AuthenticationConfig{
			Method:  types.AuthMethodKey,
			Timeout: 30 * time.Second,
		},
		PoolConfig: &types.PoolConfig{
			MaxConnections:      10,
			MaxIdleTime:         5 * time.Minute,
			ConnectionTimeout:   30 * time.Second,
			HealthCheckInterval: 1 * time.Minute,
		},
		ActingConfig: &types.ActingConfig{
			DefaultTimeout: 30 * time.Second,
			MaxParallel:    5,
			EnableRetries:  true,
			MaxRetries:     3,
			RetryDelay:     1 * time.Second,
		},
		KeysConfig: &types.KeysConfig{
			DefaultKeyPath: "~/.ssh/id_rsa",
			KeyCacheTTL:    1 * time.Hour,
			EnableCaching:  true,
		},
	}

	return NewManagerWithConfig(config, logger)
}
