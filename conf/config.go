package conf

type Config struct {
	ConnType string         `json:"conn_type"` // cluster alone  sentinel
	Alone    AloneConfig    `json:"alone"`
	Cluster  ClusterConfig  `json:"cluster"`
	Sentinel SentinelConfig `json:"sentinel"`
}
