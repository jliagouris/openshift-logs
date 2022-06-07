package main

// Main file of log pushing operator

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"push_app/src/components"
	"sync"
)

var PRODUCER_TMO int

type operator struct {
	parser        *components.LogParser
	preprocessor  *components.Preprocessor
	dataShareChan chan components.DataShare
	producers     []*components.KafkaProducer
}

// Main function of the operator
func main() {
	PRODUCER_TMO = 15 * 1000
	confList := getKafkaConfigMap()
	pushOperator := makePushOperator(confList, PRODUCER_TMO)
	pushOperator.run()
}

// Create operator object
func makePushOperator(confList []kafka.ConfigMap, producerTimeout int) *operator {
	pushOperator := operator{producers: make([]*components.KafkaProducer, len(confList))}
	for idx, conf := range confList {
		msgChan := make(chan components.ProducerMessage)
		pushOperator.producers[idx] = components.MakeKafkaProducer(&conf, msgChan, producerTimeout)
	}

	pushOperator.parser = components.MakeParser() // TODO: This will change
	pushOperator.preprocessor = components.MakePreprocessor(len(confList), pushOperator.parser.LogChan)
	pushOperator.dataShareChan = pushOperator.preprocessor.DataShareChan
	return &pushOperator
}

func (o *operator) run() {
	// Start parser goroutine
	go o.parser.ParseLoop()

	// Start preprocessor goroutine
	go o.preprocessor.PreprocessLoop()

	// Start dataShare dispatch goroutine
	go o.dispatchDataShareLoop()

	// Start the producer goroutines
	var wg sync.WaitGroup
	for _, producer := range o.producers {
		wg.Add(1)
		go producer.ProduceLoop(&wg)
	}
	wg.Wait()
}

func (o *operator) dispatchDataShareLoop() {
	for dataShare := range o.dataShareChan {
		kafkaMsg := o.dataShare2ProducerMsg(dataShare)
		if dataShare.EOF {
			for _, producer := range o.producers {
				producer.MsgChan <- kafkaMsg
			}
		} else {
			for _, producerId := range dataShare.ProducerIdArr {
				o.producers[producerId].MsgChan <- kafkaMsg
			}
		}
	}
}

func (o *operator) dataShare2ProducerMsg(dataShare components.DataShare) components.ProducerMessage {
	//TODO: Fill this
	return components.ProducerMessage{}
}

// Gets the configurations that allow the operator talk to the correct Kafka brokers
func getKafkaConfigMap() []kafka.ConfigMap {
	//TODO: Fill this, need to know how to get info
	return nil
}
