package components

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
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
		} else {
			for _, producerId := range dataShare.ProducerIdArr {
				p.producers[producerId].MsgChan <- dataShare
			}
		}
	}
}

// Generate Data shares from log
func (p *Preprocessor) log2DataShares(log Log) []DataShare {
	if !log.EOF {
		secretByteArr := createDataShares(log.Val)
		//var map1 map[string]interface{}
		//json.Unmarshal(secretByteArr[0], &map1)
		keyBytes, err := json.Marshal(log.Key)
		if err != nil {
			fmt.Println("error in log 2 datashares getting key byte array:", err)
		}
		fmt.Printf("Generated key: %v\n", string(keyBytes))
		fmt.Printf("map1: %v\n", binary.BigEndian.Uint64(secretByteArr[0]))
		share1 := DataShare{
			ProducerIdArr: []int{1, 2},
			Message: kafka.Message{
				TopicPartition: kafka.TopicPartition{Topic: &log.Topic, Partition: kafka.PartitionAny},
				Value:          secretByteArr[0],
				Key:            keyBytes,
			},
			EOF: false,
		}
		share2 := DataShare{
			ProducerIdArr: []int{0, 2},
			Message: kafka.Message{
				TopicPartition: kafka.TopicPartition{Topic: &log.Topic, Partition: kafka.PartitionAny},
				Value:          secretByteArr[1],
				Key:            keyBytes,
			},
			EOF: false,
		}
		share3 := DataShare{
			ProducerIdArr: []int{0, 1},
			Message: kafka.Message{
				TopicPartition: kafka.TopicPartition{Topic: &log.Topic, Partition: kafka.PartitionAny},
				Value:          secretByteArr[2],
				Key:            keyBytes,
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
	fmt.Printf("metrics: %v\n", metrics)
	valBytes, _ := GetBytes(metrics["value"])
	valByteArr := generateRandomIntShares(valBytes)
	return valByteArr
}

func GetBytes(key interface{}) ([]byte, error) {
	//fmt.Printf("GetBytes input type: %v\n", key)
	switch v := key.(type) {
	case int:
		//fmt.Printf("int type: %v\n", v)
		bs := make([]byte, 8)
		binary.BigEndian.PutUint64(bs, uint64(v))
		return bs, nil
	}
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(key)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

const shareDataSize = 8

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
	//data1 := binary.BigEndian.Uint64(intBytes)
	//fmt.Printf("int bytes: %v\n", intBytes)
	shares := make([][]byte, 3)
	for i := 0; i < 3; i++ {
		shares[i] = make([]byte, shareDataSize)
		shares[i][0] = 0
	}
	if _, err := rand.Read(shares[0]); err != nil {
		fmt.Printf("Sth is wrong with generating random boolean share 0: %v\n", err)
	}
	if _, err := rand.Read(shares[1]); err != nil {
		fmt.Printf("Sth is wrong with generating random boolean share 1: %v\n", err)
	}
	data := binary.BigEndian.Uint64(intBytes)
	//fmt.Printf("data: %v\n", data)
	share0 := binary.BigEndian.Uint64(shares[0])
	//fmt.Printf("share0: %v\n", share0)
	share1 := binary.BigEndian.Uint64(shares[1])
	//fmt.Printf("share1: %v\n", share1)
	binary.BigEndian.PutUint64(shares[2], data-(share0+share1))
	//fmt.Printf("share2 int: %v\n", data-(share0+share1))
	//fmt.Printf("share2: %v\n", binary.BigEndian.Uint64(shares[2]))
	//fmt.Printf("sum data: %v\n", share0+share1+binary.BigEndian.Uint64(shares[2]))
	return shares
}
