package components

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"os"
)

//Producer that sends queried and preprocessed logs to Kafka broker
type Kafka_producer struct {
	producer *kafka.Producer
}

func make_kafka_producer(configMap *kafka.ConfigMap) *Kafka_producer {
	p, err := kafka.NewProducer(configMap)

	if err != nil {
		fmt.Printf("Failed to create producer: %s", err)
		os.Exit(1)
	}

	kp := Kafka_producer{producer: p}

	return &kp
}
