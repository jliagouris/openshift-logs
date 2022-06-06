package components

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"os"
)

// KafkaProducer Producer that sends queried and preprocessed logs to Kafka broker
type KafkaProducer struct {
	producer *kafka.Producer      // producer in confluent_kafka_go
	msgChan  chan ProducerMessage // Get asynchronously generated messages through this channel
	timeout  int                  // Producer timeout
}

// ProducerMessage Message sent to the producer. Wraps real data msg to be sent
type ProducerMessage struct {
	msg   kafka.Message // msg contains log and any other information needed by Kafka
	kill  bool          // kill indicates if all logs has been produces, and we can exit the loop
	topic string        // topic the message to be sent to
}

// ProduceLoop Goroutine that loops to push messages to Secrecy kafka broker
func (p *KafkaProducer) ProduceLoop() {

	// Goroutine that handles producer event messages
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Failed to deliver message: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Produced event to topic %s: key = %-10s value = %s\n",
						*ev.TopicPartition.Topic, string(ev.Key), string(ev.Value))
				}
			}
		}
	}()

	// Produce all msgs sent to the channel until receive a close msg
	for producerMsg := range p.msgChan {
		if producerMsg.kill {
			break
		}
		p.produce(producerMsg.topic, &producerMsg.msg)
	}

	// Wait for all messages to be delivered
	p.Flush()
	p.Close()
}

// Wrapper of confluent_kafka produce method
func (p *KafkaProducer) produce(topic string, msg *kafka.Message) {
	//TODO: Do we need a delivery channel?
	err := p.producer.Produce(msg, nil)
	if err != nil {
		return
	}
}

// Events Wrapper of confluent_kafka Events method
func (p *KafkaProducer) Events() chan kafka.Event {
	return p.producer.Events()
}

// Flush Wrapper of confluent_kafka Flush() method
func (p *KafkaProducer) Flush() {
	//TODO: What do we do if some messages fail?
	unsuccessfulMsgCnt := p.producer.Flush(p.timeout)
	if unsuccessfulMsgCnt > 0 {
		fmt.Printf("%v messages failed to deliever\n", unsuccessfulMsgCnt)
	} else {
		fmt.Printf("All messages successfully delievered\n")
	}
}

func (p *KafkaProducer) Close() {
	p.producer.Close()
}

// MakeKafkaProducer Creates producer object
func MakeKafkaProducer(configMap *kafka.ConfigMap, msgChan chan ProducerMessage, timeout int) *KafkaProducer {
	p, err := kafka.NewProducer(configMap)

	if err != nil {
		fmt.Printf("Failed to create producer: %s", err)
		os.Exit(1)
	}

	kp := KafkaProducer{producer: p,
		msgChan: msgChan,
		timeout: timeout,
	}

	return &kp
}
