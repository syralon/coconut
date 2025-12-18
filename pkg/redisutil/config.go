package redisutil

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/syralon/coconut/toolkit/text"
)

type Config struct {
	// Either a single address or a seed list of host:port addresses
	// of cluster/sentinel nodes.
	Addrs []string `json:"addrs" yaml:"addrs"`

	// ClientName will execute the `CLIENT SETNAME ClientName` command for each conn.
	ClientName string `json:"client_name" yaml:"client_name"`

	// Common options.

	Protocol int    `json:"protocol" yaml:"protocol"`
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`

	SentinelUsername string `json:"sentinel_username" yaml:"sentinel_username"`
	SentinelPassword string `json:"sentinel_password" yaml:"sentinel_password"`

	MaxRetries      int           `json:"max_retries" yaml:"max_retries"`
	MinRetryBackoff text.Duration `json:"min_retry_backoff" yaml:"min_retry_backoff"`
	MaxRetryBackoff text.Duration `json:"max_retry_backoff" yaml:"max_retry_backoff"`

	DialTimeout           text.Duration `json:"dial_timeout" yaml:"dial_timeout"`
	ReadTimeout           text.Duration `json:"read_timeout" yaml:"read_timeout"`
	WriteTimeout          text.Duration `json:"write_timeout" yaml:"write_timeout"`
	ContextTimeoutEnabled bool          `json:"context_timeout_enabled" yaml:"context_timeout_enabled"`

	// ReadBufferSize is the size of the bufio.Reader buffer for each connection.
	// Larger buffers can improve performance for commands that return large responses.
	// Smaller buffers can improve memory usage for larger pools.
	//
	// default: 32KiB (32768 bytes)
	ReadBufferSize int `json:"read_buffer_size" yaml:"read_buffer_size"`

	// WriteBufferSize is the size of the bufio.Writer buffer for each connection.
	// Larger buffers can improve performance for large pipelines and commands with many arguments.
	// Smaller buffers can improve memory usage for larger pools.
	//
	// default: 32KiB (32768 bytes)
	WriteBufferSize int `json:"write_buffer_size" yaml:"write_buffer_size"`

	// PoolFIFO uses FIFO mode for each node connection pool GET/PUT (default LIFO).
	PoolFIFO bool `json:"pool_fifo" yaml:"pool_fifo"`

	PoolSize        int           `json:"pool_size" yaml:"pool_size"`
	PoolTimeout     text.Duration `json:"pool_timeout" yaml:"pool_timeout"`
	MinIdleConns    int           `json:"min_idle_conns" yaml:"min_idle_conns"`
	MaxIdleConns    int           `json:"max_idle_conns" yaml:"max_idle_conns"`
	MaxActiveConns  int           `json:"max_active_conns" yaml:"max_active_conns"`
	ConnMaxIdleTime text.Duration `json:"conn_max_idle_time" yaml:"conn_max_idle_time"`
	ConnMaxLifetime text.Duration `json:"conn_max_lifetime" yaml:"conn_max_lifetime"`

	// Only cluster clients.

	MaxRedirects   int  `json:"max_redirects" yaml:"max_redirects"`
	ReadOnly       bool `json:"read_only" yaml:"read_only"`
	RouteByLatency bool `json:"route_by_latency" yaml:"route_by_latency"`
	RouteRandomly  bool `json:"route_randomly" yaml:"route_randomly"`

	// MasterName is the sentinel master name.
	// Only for failover clients.
	MasterName string `json:"master_name" yaml:"master_name"`

	// DisableIdentity is used to disable CLIENT SETINFO command on connect.
	// default: false
	DisableIdentity bool `json:"disable_identity" yaml:"disable_identity"`

	IdentitySuffix string `json:"identity_suffix" yaml:"identity_suffix"`

	// FailingTimeoutSeconds is the timeout in seconds for marking a cluster node as failing.
	// When a node is marked as failing, it will be avoided for this duration.
	// Only applies to cluster clients. Default is 15 seconds.
	FailingTimeoutSeconds int `json:"failing_timeout_seconds" yaml:"failing_timeout_seconds"`

	UnstableResp3 bool `json:"unstable_resp_3" yaml:"unstable_resp_3"`

	// IsClusterMode can be used when only one Addrs is provided (e.g. Elasticache supports setting up cluster mode with configuration endpoint).
	IsClusterMode bool `json:"is_cluster_mode" yaml:"is_cluster_mode"`
}

func NewClient(c *Config) (redis.UniversalClient, error) {
	client := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:                 c.Addrs,
		ClientName:            c.ClientName,
		Protocol:              c.Protocol,
		Username:              c.Username,
		Password:              c.Password,
		SentinelUsername:      c.SentinelUsername,
		SentinelPassword:      c.SentinelPassword,
		MaxRetries:            c.MaxRetries,
		MinRetryBackoff:       c.MinRetryBackoff.Duration(),
		MaxRetryBackoff:       c.MaxRetryBackoff.Duration(),
		DialTimeout:           c.DialTimeout.Duration(),
		ReadTimeout:           c.ReadTimeout.Duration(),
		WriteTimeout:          c.WriteTimeout.Duration(),
		ContextTimeoutEnabled: c.ContextTimeoutEnabled,
		ReadBufferSize:        c.ReadBufferSize,
		WriteBufferSize:       c.WriteBufferSize,
		PoolFIFO:              c.PoolFIFO,
		PoolSize:              c.PoolSize,
		PoolTimeout:           c.PoolTimeout.Duration(),
		MinIdleConns:          c.MinIdleConns,
		MaxIdleConns:          c.MaxIdleConns,
		MaxActiveConns:        c.MaxActiveConns,
		ConnMaxIdleTime:       c.ConnMaxIdleTime.Duration(),
		ConnMaxLifetime:       c.ConnMaxLifetime.Duration(),
		MaxRedirects:          c.MaxRedirects,
		ReadOnly:              c.ReadOnly,
		RouteByLatency:        c.RouteByLatency,
		RouteRandomly:         c.RouteRandomly,
		MasterName:            c.MasterName,
		DisableIdentity:       c.DisableIdentity,
		IdentitySuffix:        c.IdentitySuffix,
		FailingTimeoutSeconds: c.FailingTimeoutSeconds,
		UnstableResp3:         c.UnstableResp3,
		IsClusterMode:         c.IsClusterMode,
	})
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}
	return client, nil
}
