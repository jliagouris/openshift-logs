package main

//Main file of log pushing operator

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"push_app/src/components"
)

type operator struct {
	producers []components.Kafka_producer
}

//Main function of the operator
func main() {
	//TODO: Fill this

}

//Gets the configurations that allow the operator talk to the correct Kafka brokers
func getKafkaConfigMap() []kafka.ConfigMap {
	//TODO: Fill this
	return nil
}
