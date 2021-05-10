package conf

type SentinelConfig struct {
	// Map of name => host:port addresses of ring shards.
	Addrs map[string]string `json:"addrs"`

	// Map of name => password of ring shards, to allow different shards to have
	// different passwords. It will be ignored if the Password field is set.
	Passwords string `json:"passwords"`

	// Frequency of PING commands sent to check shards availability.
	// Shard is considered down after 3 subsequent failed checks.
	HeartbeatFrequency int64 `json:"heartbeat_frequency"`

	HashReplicas int `json:"hash_replicas"`

	DB       int    `json:"db"`
	Password string `json:"password"`

	MaxRetries      int   `json:"max_retries"`
	MinRetryBackoff int64 `json:"min_retry_backoff"`
	MaxRetryBackoff int64 `json:"max_retry_backoff"`

	DialTimeout  int64 `json:"dial_timeout"`
	ReadTimeout  int64 `json:"read_timeout"`
	WriteTimeout int64 `json:"write_timeout"`

	PoolSize           int   `json:"pool_size"`
	MinIdleConns       int   `json:"min_idle_conns"`
	MaxConnAge         int64 `json:"max_conn_age"`
	PoolTimeout        int64 `json:"pool_timeout"`
	IdleTimeout        int64 `json:"idle_timeout"`
	IdleCheckFrequency int64 `json:"idle_check_frequency"`
}
