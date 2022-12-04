package configs

import "github.com/confluentinc/confluent-kafka-go/kafka"

// KafkaClientConf Config of Kafka servers
type KafkaClientConf struct {
	Confs map[string]kafka.ConfigMap `yaml:"Clusters"`
}

// ToClusterConfList Turn Kafka server config map into slice required by Kafka-go package
func (c *KafkaClientConf) ToClusterConfList() []kafka.ConfigMap {
	configSlice := make([]kafka.ConfigMap, len(c.Confs))
	idx := 0
	for _, conf := range c.Confs {
		//fmt.Printf(" server: %v\n", conf)
		configSlice[idx] = conf
		idx++
	}
	return configSlice
}
