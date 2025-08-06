package types

import (
	"time"
)

// PooledConnection represents a connection in the pool
type PooledConnection struct {
	Connection *SSHConnection `hcl:"connection"`
	InUse      bool           `hcl:"in_use"`
	CreatedAt  time.Time      `hcl:"created_at"`
	LastUsed   time.Time      `hcl:"last_used"`
	UseCount   int            `hcl:"use_count"`
}

// PoolConfig represents connection pool configuration
type PoolConfig struct {
	MaxConnections      int           `hcl:"max_connections,optional"`
	MaxIdleTime         time.Duration `hcl:"max_idle_time,optional"`
	ConnectionTimeout   time.Duration `hcl:"connection_timeout,optional"`
	HealthCheckInterval time.Duration `hcl:"health_check_interval,optional"`
}

// PoolStats represents connection pool statistics
type PoolStats struct {
	TotalConnections  int       `hcl:"total_connections"`
	ActiveConnections int       `hcl:"active_connections"`
	IdleConnections   int       `hcl:"idle_connections"`
	MaxConnections    int       `hcl:"max_connections"`
	CreatedAt         time.Time `hcl:"created_at"`
	LastHealthCheck   time.Time `hcl:"last_health_check"`
}

// ConnectionEvent represents a connection lifecycle event
type ConnectionEvent struct {
	Type      string    `hcl:"type"`
	Host      string    `hcl:"host"`
	Timestamp time.Time `hcl:"timestamp"`
	Error     string    `hcl:"error,optional"`
}

// ConnectionMetrics represents connection performance metrics
type ConnectionMetrics struct {
	Host                  string        `hcl:"host"`
	TotalConnections      int           `hcl:"total_connections"`
	SuccessfulConnections int           `hcl:"successful_connections"`
	FailedConnections     int           `hcl:"failed_connections"`
	AverageLatency        time.Duration `hcl:"average_latency"`
	LastConnection        time.Time     `hcl:"last_connection"`
}
