package components

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"math/rand"
	"strconv"
	"time"
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
		secretStrArr := createDataShares(log.Val)
		//var map1 map[string]interface{}
		//json.Unmarshal(secretStrArr[0], &map1)
		keyStringArr := getKeyString(log.Key)
		fmt.Printf("Generated key: %v\n", string(keyStringArr))
		//fmt.Printf("map1: %v\n", binary.BigEndian.Uint64(secretStrArr[0]))
		val0 := keyStringArr + " " + "0" + " " + secretStrArr[0]
		share0 := DataShare{
			ProducerIdArr: []int{1, 2},
			Message: kafka.Message{
				TopicPartition: kafka.TopicPartition{Topic: &log.Topic, Partition: kafka.PartitionAny},
				Value:          []byte(val0),
				Key:            []byte("0"),
			},
			EOF: false,
		}
		val1 := keyStringArr + " " + "1" + " " + secretStrArr[1]
		share1 := DataShare{
			ProducerIdArr: []int{0, 2},
			Message: kafka.Message{
				TopicPartition: kafka.TopicPartition{Topic: &log.Topic, Partition: kafka.PartitionAny},
				Value:          []byte(val1),
				Key:            []byte("1"),
			},
			EOF: false,
		}
		val2 := keyStringArr + " " + "2" + " " + secretStrArr[2]
		share2 := DataShare{
			ProducerIdArr: []int{0, 1},
			Message: kafka.Message{
				TopicPartition: kafka.TopicPartition{Topic: &log.Topic, Partition: kafka.PartitionAny},
				Value:          []byte(val2),
				Key:            []byte("2"),
			},
			EOF: false,
		}
		return []DataShare{share0, share1, share2}
	} else {
		share := DataShare{
			ProducerIdArr: []int{0, 1, 2},
			EOF:           true,
		}
		return []DataShare{share}
	}
}

func getKeyString(key DataShareKey) string {
	str := ""
	str += key.ClientId + " " + strconv.Itoa(int(key.QueryId)) + " " + strconv.Itoa(int(key.SeqNum))
	return str
}

func createDataShares(metrics map[string]interface{}) []string {
	fmt.Printf("metrics: %v\n", metrics)
	//valBytes, _ := GetBytes(metrics["value"])
	//valByteArr := generateRandomIntShares(valBytes)
	shareArr := make([]string, 3)
	shareIntArr := generateRandomIntShares(metrics["value"].(int))
	for idx, num := range shareIntArr {
		shareArr[idx] = strconv.Itoa(num)
	}
	return shareArr
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
/*
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
*/

func generateRandomIntShares(num int) []int {
	//data1 := binary.BigEndian.Uint64(intBytes)
	//fmt.Printf("int bytes: %v\n", intBytes)
	shares := make([]int, 3)
	rand.Seed(time.Now().UnixNano())
	shares[0] = rand.Int()
	shares[1] = rand.Int()
	shares[2] = num - shares[1] - shares[0]
	//fmt.Printf("share2 int: %v\n", data-(share0+share1))
	//fmt.Printf("share2: %v\n", binary.BigEndian.Uint64(shares[2]))
	//fmt.Printf("sum data: %v\n", share0+share1+binary.BigEndian.Uint64(shares[2]))
	return shares
}
