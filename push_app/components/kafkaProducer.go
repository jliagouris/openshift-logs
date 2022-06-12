package components

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"os"
	"sync"
)

// KafkaProducer Producer that sends queried and preprocessed logs to Kafka broker
type KafkaProducer struct {
	producer *kafka.Producer // producer in confluent_kafka_go
	//TODO: Buffered channels for better performance? Maybe?
	MsgChan chan DataShare // Get asynchronously generated messages through this channel
	timeout int            // Producer timeout
}

// ProduceLoop Goroutine that loops to push messages to Secrecy kafka broker
func (p *KafkaProducer) ProduceLoop(wg *sync.WaitGroup) {

	// Goroutine that handles producer event messages asynchronously
	go func() {
		for e := range p.events() {
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
	for dataShare := range p.MsgChan {
		if dataShare.EOF {
			break
		}
		p.produce(&dataShare.Message)
	}

	// Wait for all messages to be delivered
	p.flush()
	p.Close()
	wg.Done()
}

// Wrapper of confluent_kafka produce method
func (p *KafkaProducer) produce(msg *kafka.Message) {
	//TODO: Do we need a delivery channel?
	err := p.producer.Produce(msg, nil)
	if err != nil {
		fmt.Printf("An message failed to be sent\n")
	}
}

// events Wrapper of confluent_kafka events method
func (p *KafkaProducer) events() chan kafka.Event {
	return p.producer.Events()
}

// flush Wrapper of confluent_kafka flush() method, waits until all messages are acked or timeout
func (p *KafkaProducer) flush() {
	//TODO: What do we do if some messages fail?
	unsuccessfulMsgCnt := p.producer.Flush(p.timeout)
	if unsuccessfulMsgCnt > 0 {
		fmt.Printf("%v messages failed to be delievered\n", unsuccessfulMsgCnt)
	} else {
		fmt.Printf("All messages successfully delievered\n")
	}
}

// Close Wrapper of confluent_kafka Close() method
func (p *KafkaProducer) Close() {
	p.producer.Close()
}

// MakeKafkaProducer Creates producer object
func MakeKafkaProducer(configMap *kafka.ConfigMap, msgChan chan DataShare, timeout int) *KafkaProducer {
	//fmt.Println("Called MakeKafkaProducer")

	p, err := kafka.NewProducer(configMap)

	if err != nil {
		fmt.Printf("Failed to create producer: %s", err)
		os.Exit(1)
	}

	kp := KafkaProducer{producer: p,
		MsgChan: msgChan,
		timeout: timeout,
	}

	return &kp
}
