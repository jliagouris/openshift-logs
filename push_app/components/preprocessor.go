package components

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log"
	"strconv"
)

// Preprocesses logs queried by parser

type Preprocessor struct {
	ParsedChan    <-chan Log     // Channel shared with logParser, gets queried logs
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
		ParsedChan:    LogChan,
		DataShareChan: DataShareChan,
		numProducers:  numProducers,
		producers:     producers,
	}
	return &preprocessor
}

// PreprocessLoop Goroutine that iteratively processes logs passed by parser
func (p *Preprocessor) PreprocessLoop() {
	logCnt := 0
	go p.dispatchDataShareLoop()
	for log := range p.ParsedChan {
		logCnt++
		fmt.Printf("logCnt: %v\n", logCnt)
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
			fmt.Println("Preprocessor EOF")
			for _, producer := range p.producers {
				producer.MsgChan <- dataShare
			}
			//break
		} else {
			for _, producerId := range dataShare.ProducerIdArr {
				p.producers[producerId].MsgChan <- dataShare
			}
		}
	}
}

// Generate Data shares from log
func (p *Preprocessor) log2DataShares(log Log) []DataShare {
	//TODO: This is currently using only values of metrics, will need to change to support integer
	if !log.EOF {
		secretByteArr := createDataShares(log.Val)
		share1 := DataShare{
			ProducerIdArr: []int{1, 2},
			Message: kafka.Message{
				TopicPartition: kafka.TopicPartition{Topic: &log.Topic, Partition: kafka.PartitionAny},
				Value:          secretByteArr[0],
				Key:            []byte(strconv.Itoa(1)),
			},
			EOF: false,
		}
		share2 := DataShare{
			ProducerIdArr: []int{0, 2},
			Message: kafka.Message{
				TopicPartition: kafka.TopicPartition{Topic: &log.Topic, Partition: kafka.PartitionAny},
				Value:          secretByteArr[1],
				Key:            []byte(strconv.Itoa(2)),
			},
			EOF: false,
		}
		share3 := DataShare{
			ProducerIdArr: []int{0, 1},
			Message: kafka.Message{
				TopicPartition: kafka.TopicPartition{Topic: &log.Topic, Partition: kafka.PartitionAny},
				Value:          secretByteArr[2],
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

func createDataShares(metrics map[string]interface{}) [][]byte {
	mapArr := make([]map[string][]byte, 3)
	mapArr[0] = make(map[string][]byte)
	mapArr[1] = make(map[string][]byte)
	mapArr[2] = make(map[string][]byte)
	for key, val := range metrics {
		valBytes, _ := GetBytes(val)
		valByteArr := generateRandomIntShares(valBytes)
		mapArr[0][key] = valByteArr[0]
		mapArr[1][key] = valByteArr[1]
		mapArr[2][key] = valByteArr[2]
	}
	var bytesArr [][]byte
	for valMap := range mapArr {
		jsonBytes, err := json.Marshal(valMap)
		if err != nil {
			log.Println("Problem in preprocessor: fail to serialize values in createDataShares")
		}
		bytesArr = append(bytesArr, jsonBytes)
	}
	return bytesArr
}

func GetBytes(key interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(key)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

const shareDataSize = 64

/*
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
	for i := range log.Val {
		shares[2][i] = log.Val[i] ^ (shares[0][i] ^ shares[1][i])
	}
	return shares
}
*/
func generateRandomIntShares(intBytes []byte) [][]byte {
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
	data := binary.LittleEndian.Uint64(intBytes)
	share0 := binary.LittleEndian.Uint64(shares[0])
	share1 := binary.LittleEndian.Uint64(shares[1])
	binary.LittleEndian.PutUint64(shares[2], data-(share0+share1))
	return shares
}
