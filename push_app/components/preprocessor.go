package components

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"strconv"
)

// Preprocesses logs queried by parser

type Preprocessor struct {
	LogChan       <-chan Log     // Channel shared with logParser, gets queried logs
	DataShareChan chan DataShare // Channel shared with operator
	numProducers  int            // Number of producers, used in producer decision
	producers     []*KafkaProducer
}

// DataShare is the preprocessed share of data to be sent to secrecy servers
type DataShare struct {
	ProducerIdArr []int         // Which producers does this data share go to?
	Message       kafka.Message // Message to the producer
	EOF           bool          // If all logs are already processed and the task can end
}

// MakePreprocessor Creates preprocessor object
func MakePreprocessor(numProducers int, LogChan <-chan Log, DataShareChan chan DataShare, producers []*KafkaProducer) *Preprocessor {
	preprocessor := Preprocessor{
		LogChan:       LogChan,
		DataShareChan: DataShareChan,
		numProducers:  numProducers,
		producers:     producers,
	}
	return &preprocessor
}

// PreprocessLoop Goroutine that iteratively processes logs passed by parser
func (p *Preprocessor) PreprocessLoop() {
	//logCnt := 0
	go p.dispatchDataShareLoop()
	for log := range p.LogChan {
		//logCnt++
		//fmt.Printf("logCnt: %v\n", logCnt)
		dataShares := p.log2DataShares(log)

		// Send data shares to the channel
		for _, dataShare := range dataShares {
			p.DataShareChan <- dataShare
		}
	}
}

func (p *Preprocessor) dispatchDataShareLoop() {
	//dataShareCnt := 0
	for dataShare := range p.DataShareChan {
		//dataShareCnt++
		//fmt.Printf("Datashare cnt: %v\n", dataShareCnt)
		//fmt.Printf("datashare content: %v\n", dataShare)
		if dataShare.EOF {
			for _, producer := range p.producers {
				producer.MsgChan <- dataShare
			}
		} else {
			for _, producerId := range dataShare.ProducerIdArr {
				p.producers[producerId].MsgChan <- dataShare
			}
		}
	}
}

// Generate Data shares from log
func (p *Preprocessor) log2DataShares(log Log) []DataShare {
	//TODO: This is currently generating naive test data, will need to change
	if !log.EOF {
		share1 := DataShare{
			ProducerIdArr: []int{1, 2},
			Message: kafka.Message{
				TopicPartition: kafka.TopicPartition{Topic: &log.topic, Partition: kafka.PartitionAny},
				Value:          log.val,
				Key:            []byte(strconv.Itoa(1)),
			},
			EOF: false,
		}
		share2 := DataShare{
			ProducerIdArr: []int{0, 2},
			Message: kafka.Message{
				TopicPartition: kafka.TopicPartition{Topic: &log.topic, Partition: kafka.PartitionAny},
				Value:          log.val,
				Key:            []byte(strconv.Itoa(2)),
			},
			EOF: false,
		}
		share3 := DataShare{
			ProducerIdArr: []int{0, 1},
			Message: kafka.Message{
				TopicPartition: kafka.TopicPartition{Topic: &log.topic, Partition: kafka.PartitionAny},
				Value:          log.val,
				Key:            []byte(strconv.Itoa(3)),
			},
			EOF: false,
		}
		return []DataShare{share1, share2, share3}
	} else {
		share := DataShare{
			ProducerIdArr: []int{0, 1, 2},
			EOF:           true,
		}
		return []DataShare{share}
	}
}

const shareDataSize = 64

func generateRandomBooleanShares(log Log) [][]byte {
	shares := make([][]byte, 3)
	for i := 0; i < 3; i++ {
		shares[i] = make([]byte, shareDataSize>>3)
	}
	if _, err := rand.Read(shares[0]); err == nil {
		fmt.Printf("Sth is wrong with generating random boolean share 0")
	}
	if _, err := rand.Read(shares[1]); err == nil {
		fmt.Printf("Sth is wrong with generating random boolean share 1")
	}
	for i := range log.val {
		shares[2][i] = log.val[i] ^ (shares[0][i] ^ shares[1][i])
	}
	return shares
}

func generateRandomIntShares(log Log) [][]byte {
	shares := make([][]byte, 3)
	for i := 0; i < 3; i++ {
		shares[i] = make([]byte, shareDataSize>>3)
	}
	if _, err := rand.Read(shares[0]); err == nil {
		fmt.Printf("Sth is wrong with generating random boolean share 0")
	}
	if _, err := rand.Read(shares[1]); err == nil {
		fmt.Printf("Sth is wrong with generating random boolean share 1")
	}
	data := binary.LittleEndian.Uint64(log.val)
	share0 := binary.LittleEndian.Uint64(shares[0])
	share1 := binary.LittleEndian.Uint64(shares[1])
	binary.LittleEndian.PutUint64(shares[2], data-(share0+share1))
	return shares
}
