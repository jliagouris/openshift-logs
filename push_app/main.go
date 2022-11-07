package main

// Main file of log pushing adHocOperator

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"push_app/components"
	"push_app/configs"
	"sync"
)

// PRODUCER_TMO is the timeout for producers
var PRODUCER_TMO int

// adHocOperator struct
type adHocOperator struct {
	parser       *components.LogParser
	preprocessor *components.Preprocessor
	producers    []*components.KafkaProducer
	dataSource   *components.PrometheusDataSource
	listener     *components.AdHocListener
}

// Main function of the adHocOperator
func main() {
	PRODUCER_TMO = 15 * 1000 // This is an arbitrarily chosen timeout
	conf := getConfig()
	fmt.Printf("global conf: %v\n", conf.OpConf)
	pushOperator := makeOperator(conf, PRODUCER_TMO)
	pushOperator.run()
}

// Create adHocOperator object
func makeOperator(conf Config, producerTimeout int) *adHocOperator {
	clusterConfList := conf.KafkaConf.ToClusterConfList()
	fmt.Printf("Generated %v configs\n", clusterConfList)
	ahOperator := adHocOperator{producers: make([]*components.KafkaProducer, len(clusterConfList))}
	for idx, clusterConf := range clusterConfList {
		msgChan := make(chan components.DataShare, conf.OpConf.ChanBufSize)
		ahOperator.producers[idx] = components.MakeKafkaProducer(&clusterConf, msgChan, producerTimeout)
	}
	//Get parser config
	config := components.LogConfig{}
	err := config.LoadConfig()
	if err != nil {
		panic(err)
	}
	//ahOperator.dataSource = components.MakeLokiDataSource(config)
	//datasource := components.MakeDataSource(conf.OpConf)
	ahOperator.listener = components.MakeListener(conf.OpConf)
	ahOperator.dataSource = components.MakePrometheusDataSource(conf.OpConf, ahOperator.listener.GetQueryChan())
	//fmt.Println("Datasource conf: ******************")
	//fmt.Printf("%v\n", *ahOperator.dataSource.Conf)
	ahOperator.parser = components.MakeParser(conf.OpConf.ChanBufSize, ahOperator.dataSource.GetDataChan(), config) // TODO: This will change
	DataShareChan := make(chan components.DataShare, conf.OpConf.ChanBufSize)
	ahOperator.preprocessor = components.MakePreprocessor(len(clusterConfList), ahOperator.parser.ParsedChan, DataShareChan, ahOperator.producers)
	return &ahOperator
}

// Main thread of the adHocOperator
func (o *adHocOperator) run() {

	//Start listener
	go o.listener.Run()

	// Activate data source
	go o.dataSource.Run()

	// Start parser goroutine
	go o.parser.Run()

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

// Config General config structure
type Config struct {
	KafkaConf configs.KafkaClientConf
	OpConf    configs.OperatorConf
}

// Get config from config file
func getConfig() Config {
	config := Config{}
	yamlFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	config.OpConf = getGlobalConfig(yamlFile)
	config.KafkaConf = getKafkaConfig(yamlFile)
	return config
}

// Get kafka cluster config
func getKafkaConfig(yamlFile []byte) configs.KafkaClientConf {
	yamlFile, _ = ioutil.ReadFile("config.yaml")
	clustersConf := configs.KafkaClientConf{}
	err := yaml.Unmarshal(yamlFile, &clustersConf)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	fmt.Printf("clusters conf: %v\n", clustersConf.Confs)
	return clustersConf
}

// Get global adHocOperator config
func getGlobalConfig(yamlFile []byte) configs.OperatorConf {
	opConf := configs.OperatorConf{ChanBufSize: 0}
	err := yaml.Unmarshal(yamlFile, &opConf)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	return opConf
}
