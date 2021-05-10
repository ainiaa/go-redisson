package conf

type ClusterConfig struct {
	// A seed list of host:port addresses of cluster nodes.
	Addrs []string `json:"addrs"`

	// The maximum number of retries before giving up. Command is retried
	// on network errors and MOVED/ASK redirects.
	// Default is 8 retries.
	MaxRedirects int `json:"max_redirects"`

	// Enables read-only commands on slave nodes.
	ReadOnly bool `json:"read_only"`
	// Allows routing read-only commands to the closest master or slave node.
	// It automatically enables ReadOnly.
	RouteByLatency bool `json:"route_by_latency"`
	// Allows routing read-only commands to the random master or slave node.
	// It automatically enables ReadOnly.
	RouteRandomly bool `json:"route_randomly"`

	Username string `json:"username"`
	Password string `json:"password"`

	MaxRetries      int   `json:"max_retries"`
	MinRetryBackoff int64 `json:"min_retry_backoff"`
	MaxRetryBackoff int64 `json:"max_retry_backoff"`

	DialTimeout  int64 `json:"dial_timeout"`
	ReadTimeout  int64 `json:"read_timeout"`
	WriteTimeout int64 `json:"write_timeout"`

	// PoolSize applies per cluster node and not for the whole cluster.
	PoolSize           int   `json:"pool_size"`
	MinIdleConns       int   `json:"min_idle_conns"`
	MaxConnAge         int64 `json:"max_conn_age"`
	PoolTimeout        int64 `json:"pool_timeout"`
	IdleTimeout        int64 `json:"idle_timeout"`
	IdleCheckFrequency int64 `json:"idle_check_frequency"`
}
