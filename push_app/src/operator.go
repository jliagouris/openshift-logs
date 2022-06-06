package main

//Main file of log pushing operator

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"push_app/src/components"
)

type operator struct {
	producers []*components.KafkaProducer
}

//Main function of the operator
func main() {
	confList := getKafkaConfigMap()
	pushOperator := makePushOperator(confList)
	//TODO: parser parses log, preprocessor transforms the log, producers asynchronously pushes message to broker. Communicate via channel
}

//Create operator object
func makePushOperator(confList []kafka.ConfigMap) *operator {
	pushOperator := operator{producers: make([]*components.KafkaProducer, len(confList))}
	for idx, conf := range confList {
		pushOperator.producers[idx] = components.MakeKafkaProducer(&conf)
	}
	return &pushOperator
}

//Gets the configurations that allow the operator talk to the correct Kafka brokers
func getKafkaConfigMap() []kafka.ConfigMap {
	//TODO: Fill this
	return nil
}
