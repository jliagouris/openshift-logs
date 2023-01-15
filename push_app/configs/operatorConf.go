package configs

// OperatorConf Global Operator Configs
type OperatorConf struct {
	ChanBufSize    int    `yaml:"chan_buf_size"`
	DataSourceType string `yaml:"data_source_type"`
	TimeInterval   string `yaml:"time_interval"`
	ClientId       string
}
