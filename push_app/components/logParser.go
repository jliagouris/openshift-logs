package components

import (
	"fmt"
	"strconv"
)

//Queries Loki

// LogParser queries Loki logs, structure the log, and puts it into a channel shared with preprocessor
type LogParser struct {
	LogChan chan Log // Channel to communicate with preprocessor goroutine
	//TODO: Fill this
}

// Log Parsed log
type Log struct {
	EOF bool
	//TODO: Fill this
	val   []byte
	topic string
}

// MakeParser creates LogParser object, this definitely will change
func MakeParser(chanBufSize int) *LogParser {
	//TODO: This will change
	parser := LogParser{LogChan: make(chan Log, chanBufSize)}
	fmt.Printf("parser channel buffer size: %v\n", cap(parser.LogChan))
	return &parser
}

// ParseLoop Iteratively parse logs using Grafana Loki
func (parser *LogParser) ParseLoop() {
	//TODO: //TODO: This is currently generating naive test data, will need to change
	for i := 0; i < 10; i++ {
		parser.LogChan <- Log{
			EOF:   false,
			val:   []byte(strconv.Itoa(i)),
			topic: "naive_test",
		}
	}
	parser.LogChan <- Log{EOF: true}
}
