package main

// Main file of log pushing operator

import (
	"fmt"
	"io/ioutil"
	"log"
	"push_app/components"
	"sync"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"gopkg.in/yaml.v3"
)

// PRODUCER_TMO is the timeout for producers
var PRODUCER_TMO int

// operator struct
type operator struct {
	parser       *components.LogParser
	preprocessor *components.Preprocessor
	producers    []*components.KafkaProducer
}

// Main function of the operator
func main() {
	PRODUCER_TMO = 15 * 1000 // This is an arbitrarily chosen timeout
	conf := getConfig()
	fmt.Printf("global conf: %v\n", conf.OpConf)
	pushOperator := makePushOperator(conf, PRODUCER_TMO)
	pushOperator.run()
}

// Create operator object
func makePushOperator(conf Config, producerTimeout int) *operator {
	clusterConfList := conf.ClusterConf.toClusterConfList()
	fmt.Printf("Generated %v configs\n", clusterConfList)
	pushOperator := operator{producers: make([]*components.KafkaProducer, len(clusterConfList))}
	for idx, clusterConf := range clusterConfList {
		msgChan := make(chan components.DataShare, conf.OpConf.ChanBufSize)
		pushOperator.producers[idx] = components.MakeKafkaProducer(&clusterConf, msgChan, producerTimeout)
	}
	pushOperator.parser = components.MakeParser(conf.OpConf.ChanBufSize) // TODO: This will change
	DataShareChan := make(chan components.DataShare, conf.OpConf.ChanBufSize)
	pushOperator.preprocessor = components.MakePreprocessor(len(clusterConfList), pushOperator.parser.LogChan, DataShareChan, pushOperator.producers)
	return &pushOperator
}

// Main thread of the operator
func (o *operator) run() {
	// Start parser goroutine
	go o.parser.ParseLoop()

	// Start preprocessor goroutine
	go o.preprocessor.PreprocessLoop()

	// Start the producer goroutines
	var wg sync.WaitGroup
	for _, producer := range o.producers {
		wg.Add(1)
		go producer.ProduceLoop(&wg)
	}
	wg.Wait()
}

// Clusters Config of Kafka servers
type Clusters struct {
	Confs map[string]kafka.ConfigMap `yaml:"Clusters"`
}

// Turn Kafka server config map into slice required by Kafka-go package
func (c *Clusters) toClusterConfList() []kafka.ConfigMap {
	configSlice := make([]kafka.ConfigMap, len(c.Confs))
	idx := 0
	for _, conf := range c.Confs {
		//fmt.Printf(" server: %v\n", conf)
		configSlice[idx] = conf
		idx++
	}
	//fmt.Printf("Generated %v configs\n", configSlice)
	//fmt.Printf("Generated %v configs\n", len(configSlice))
	return configSlice
}

// OperatorConf Global Operator Configs
type OperatorConf struct {
	ChanBufSize int `yaml:"chan_buf_size"`
}

// Config General config structure
type Config struct {
	ClusterConf Clusters
	OpConf      OperatorConf
}

// Get config from config file
func getConfig() Config {
	config := Config{}
	yamlFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	config.OpConf = getGlobalConfig(yamlFile)
	config.ClusterConf = getKafkaConfig(yamlFile)
	return config
}

// Get kafka cluster config
func getKafkaConfig(yamlFile []byte) Clusters {
	yamlFile, _ = ioutil.ReadFile("config.yaml")
	clustersConf := Clusters{}
	err := yaml.Unmarshal(yamlFile, &clustersConf)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	fmt.Printf("clusters conf: %v\n", clustersConf.Confs)
	return clustersConf
}

// Get global operator config
func getGlobalConfig(yamlFile []byte) OperatorConf {
	opConf := OperatorConf{ChanBufSize: 0}
	err := yaml.Unmarshal(yamlFile, &opConf)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	return opConf
}
