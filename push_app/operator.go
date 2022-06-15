package main

// Main file of log pushing operator

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"push_app/components"
	"sync"
)

// PRODUCER_TMO is the timeout for producers
var PRODUCER_TMO int

// operator struct
type operator struct {
	parser        *components.LogParser
	preprocessor  *components.Preprocessor
	dataShareChan <-chan components.DataShare //channel shared by operator and preprocessor
	producers     []*components.KafkaProducer
}

// Main function of the operator
func main() {
	PRODUCER_TMO = 15 * 1000 // This is an arbitrarily chosen timeout
	confList := getKafkaConfigMap()
	pushOperator := makePushOperator(confList, PRODUCER_TMO)
	pushOperator.run()
}

// Create operator object
func makePushOperator(confList []kafka.ConfigMap, producerTimeout int) *operator {
	//fmt.Println("Called makePushOperator")
	pushOperator := operator{producers: make([]*components.KafkaProducer, len(confList))}
	for idx, conf := range confList {
		//fmt.Printf("Making producer #%v\n", idx)
		msgChan := make(chan components.DataShare)
		pushOperator.producers[idx] = components.MakeKafkaProducer(&conf, msgChan, producerTimeout)
	}
	pushOperator.parser = components.MakeParser() // TODO: This will change
	DataShareChan := make(chan components.DataShare)
	pushOperator.preprocessor = components.MakePreprocessor(len(confList), pushOperator.parser.LogChan, DataShareChan, pushOperator.producers)
	pushOperator.dataShareChan = DataShareChan
	return &pushOperator
}

// Main thread of the operator
func (o *operator) run() {
	// Start parser goroutine
	go o.parser.ParseLoop()

	// Start preprocessor goroutine
	go o.preprocessor.PreprocessLoop()

	// Start dataShare dispatch goroutine
	//go o.dispatchDataShareLoop()

	// Start the producer goroutines
	var wg sync.WaitGroup
	for _, producer := range o.producers {
		wg.Add(1)
		go producer.ProduceLoop(&wg)
	}
	wg.Wait()
}

/*
// Dispatch data shares to their designated producers
func (o *operator) dispatchDataShareLoop() {
	//dataShareCnt := 0
	for dataShare := range o.dataShareChan {
		//dataShareCnt++
		//fmt.Printf("Datashare cnt: %v\n", dataShareCnt)
		//fmt.Printf("datashare content: %v\n", dataShare)
		if dataShare.EOF {
			for _, producer := range o.producers {
				producer.MsgChan <- dataShare
			}
		} else {
			for _, producerId := range dataShare.ProducerIdArr {
				o.producers[producerId].MsgChan <- dataShare
			}
		}
	}
}
*/

/*
// Turn a data share to Kafka messages
func (o *operator) dataShare2ProducerMsg(dataShare components.DataShare) components.ProducerMessage {
	//TODO: Fill this
	msg := components.ProducerMessage{
		Msg:   dataShare.Message,
		EOF:   false,
		Topic: "",
	}
	return
}
*/

type Clusters struct {
	Confs map[string]kafka.ConfigMap `yaml:"Clusters"`
}

// Gets the configurations that allow the operator talk to the correct Kafka brokers
func getKafkaConfigMap() []kafka.ConfigMap {
	clustersConf := Clusters{}
	yamlFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, &clustersConf)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	configSlice := make([]kafka.ConfigMap, len(clustersConf.Confs))
	idx := 0
	for _, conf := range clustersConf.Confs {
		//fmt.Printf(" server: %v\n", conf)
		configSlice[idx] = conf
		idx++
	}
	//fmt.Printf("Generated %v configs\n", configSlice)
	//fmt.Printf("Generated %v configs\n", len(configSlice))
	return configSlice
}
