package components

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"os"
)

// KafkaProducer Producer that sends queried and preprocessed logs to Kafka broker
type KafkaProducer struct {
	producer *kafka.Producer
}

// MakeKafkaProducer Creates producer object
func MakeKafkaProducer(configMap *kafka.ConfigMap) *KafkaProducer {
	p, err := kafka.NewProducer(configMap)

	if err != nil {
		fmt.Printf("Failed to create producer: %s", err)
		os.Exit(1)
	}

	kp := KafkaProducer{producer: p}

	return &kp
}
