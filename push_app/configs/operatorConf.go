package configs

// OperatorConf Global Operator Configs
type OperatorConf struct {
	ChanBufSize    int    `yaml:"chan_buf_size"`
	DataSourceType string `yaml:"data_source_type"`
	TimeInterval   uint8  `yaml:"time_interval"`
	ClientId       string `yaml:"client_id"`
}
