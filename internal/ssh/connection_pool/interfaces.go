package connection_pool

import (
	"spooky/internal/ssh/types"
)

// ConnectionPool defines the interface for connection pooling
type ConnectionPool interface {
	GetConnection(host string) (*types.SSHConnection, error)
	ReturnConnection(connection *types.SSHConnection) error
	CloseConnection(connection *types.SSHConnection) error
	CloseAllConnections() error
	GetStats() *types.PoolStats
}
