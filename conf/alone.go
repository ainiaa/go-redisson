package conf

type AloneConfig struct {
	// The network type, either tcp or unix.
	// Default is tcp.
	Network string `json:"network"`
	// host:port address.
	Addr     string `json:"addr"`
	Username string `json:"username"`
	Password string `json:"password"`
	DB       int    `json:"db"`

	MaxRetries int `json:"max_retries"`
	// Minimum backoff between each retry.
	// Default is 8 milliseconds; -1 disables backoff.
	MinRetryBackoff int64 `json:"min_retry_backoff"`
	// Maximum backoff between each retry.
	// Default is 512 milliseconds; -1 disables backoff.
	MaxRetryBackoff int64 `json:"max_retry_backoff"`

	// Dial timeout for establishing new connections.
	// Default is 5 seconds.
	DialTimeout int64 `json:"dial_timeout"`
	// Timeout for socket reads. If reached, commands will fail
	// with a timeout instead of blocking. Use value -1 for no timeout and 0 for default.
	// Default is 3 seconds.
	ReadTimeout int64 `json:"read_timeout"`
	// Timeout for socket writes. If reached, commands will fail
	// with a timeout instead of blocking.
	// Default is ReadTimeout.
	WriteTimeout int64 `json:"write_timeout"`

	// Maximum number of socket connections.
	// Default is 10 connections per every CPU as reported by runtime.NumCPU.
	PoolSize int `json:"pool_size"`
	// Minimum number of idle connections which is useful when establishing
	// new connection is slow.
	MinIdleConns int `json:"min_idle_conns"`
	// Connection age at which client retires (closes) the connection.
	// Default is to not close aged connections.
	MaxConnAge int64 `json:"max_conn_age"`
	// Amount of time client waits for connection if all connections
	// are busy before returning an error.
	// Default is ReadTimeout + 1 second.
	PoolTimeout int64 `json:"pool_timeout"`
	// Amount of time after which client closes idle connections.
	// Should be less than server's timeout.
	// Default is 5 minutes. -1 disables idle timeout check.
	IdleTimeout int64 `json:"idle_timeout"`
	// Frequency of idle checks made by idle connections reaper.
	// Default is 1 minute. -1 disables idle connections reaper,
	// but idle connections are still discarded by the client
	// if IdleTimeout is set.
	IdleCheckFrequency int64 `json:"idle_check_frequency"`
}
