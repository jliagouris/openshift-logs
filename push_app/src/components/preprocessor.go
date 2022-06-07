package components

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"strconv"
)

// Preprocesses logs queried by parser

type Preprocessor struct {
	LogChan       <-chan Log       // Channel shared with logParser, gets queried logs
	DataShareChan chan<- DataShare // Channel shared with operator
	numProducers  int              // Number of producers, used in producer decision
}

// DataShare is the preprocessed share of data to be sent to secrecy servers
type DataShare struct {
	ProducerIdArr []int         // Which producers does this data share go to?
	Message       kafka.Message // Message to the producer
	EOF           bool          // If all logs are already processed and the task can end
}

// MakePreprocessor Creates preprocessor object
func MakePreprocessor(numProducers int, LogChan <-chan Log, DataShareChan chan<- DataShare) *Preprocessor {
	preprocessor := Preprocessor{
		LogChan:       LogChan,
		DataShareChan: DataShareChan,
		numProducers:  numProducers,
	}
	return &preprocessor
}

// PreprocessLoop Goroutine that iteratively processes logs passed by parser
func (p *Preprocessor) PreprocessLoop() {
	for log := range p.LogChan {
		dataShares := p.log2DataShares(log)

		// Send data shares to the channel
		for _, dataShare := range dataShares {
			p.DataShareChan <- dataShare
		}
	}
}

// Generate Data shares from log
func (p *Preprocessor) log2DataShares(log Log) []DataShare {
	//TODO: Fill this
	if !log.EOF {
		share1 := DataShare{
			ProducerIdArr: []int{2, 3},
			Message: kafka.Message{
				TopicPartition: kafka.TopicPartition{Topic: &log.topic, Partition: kafka.PartitionAny},
				Value:          []byte(strconv.Itoa(log.val)),
				Key:            []byte(strconv.Itoa(1)),
			},
			EOF: false,
		}
		p.DataShareChan <- share1
		share2 := DataShare{
			ProducerIdArr: []int{1, 3},
			Message: kafka.Message{
				TopicPartition: kafka.TopicPartition{Topic: &log.topic, Partition: kafka.PartitionAny},
				Value:          []byte(strconv.Itoa(log.val)),
				Key:            []byte(strconv.Itoa(2)),
			},
			EOF: false,
		}
		share3 := DataShare{
			ProducerIdArr: []int{2, 3},
			Message: kafka.Message{
				TopicPartition: kafka.TopicPartition{Topic: &log.topic, Partition: kafka.PartitionAny},
				Value:          []byte(strconv.Itoa(log.val)),
				Key:            []byte(strconv.Itoa(3)),
			},
			EOF: false,
		}
		return []DataShare{share1, share2, share3}
	} else {
		share := DataShare{
			ProducerIdArr: []int{1, 2, 3},
			EOF:           true,
		}
		return []DataShare{share}
	}
}
