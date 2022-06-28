package components

import (
	"encoding/json"
	"fmt"
	"strconv"
)

//Queries Loki

// LogParser queries Loki logs, structure the log, and puts it into a channel shared with preprocessor
type LogParser struct {
	LogChan   chan Log  // Channel to communicate with preprocessor goroutine
	InLogChan chan Log  // Channel to communicate with data source goroutine
	config    LogConfig // Configuration for the log parser
}

// MakeParser creates LogParser object, this definitely will change
func MakeParser(chanBufSize int, inLogChan chan Log, config LogConfig) *LogParser {
	//TODO: This will change
	parser := LogParser{LogChan: make(chan Log, chanBufSize), InLogChan: inLogChan, config: config}
	fmt.Printf("parser channel buffer size: %v\n", cap(parser.LogChan))
	return &parser
}

// Main goroutine for LogParser, checks for incoming log from InLogChan and parses it.
func (parser *LogParser) Run() {
	exit := false
	for !exit {
		select {
		case log := <-parser.InLogChan:
			if log.EOF {
				parser.LogChan <- log
				exit = true
			} else {
				processedLog, err := parser.ProcessLog(log.Val)

				if err == nil {
					log.Val = processedLog
					parser.LogChan <- log
				} else {
					fmt.Println(err)
				}
			}
		}
	}
}

// ProcessLog takes a log and matches the schema of the log with config	and returns the processed log.
func (parser *LogParser) ProcessLog(log []byte) ([]byte, error) {
	var parsedLog map[string]interface{}
	err := json.Unmarshal(log, &parsedLog)

	if err != nil {
		return nil, err
	}

	processedLog := make(map[string]interface{})

	for _, schema := range parser.config.ParserConfig.LogSchema {
		if val, ok := parsedLog[schema.Key]; ok {
			processedLog[schema.Key] = val
		} else {
			if schema.Nullable {
				processedLog[schema.Key] = nil
			} else {
				return nil, fmt.Errorf("Missing required field: %v", schema.Key)
			}
		}
	}

	processedLogBytes, err := json.Marshal(processedLog)

	if err != nil {
		return nil, err
	}

	return processedLogBytes, nil
}

func (parser *LogParser) GetLogChan() chan Log {
	return parser.LogChan
}

// ParseLoop Iteratively parse logs using Grafana Loki
func (parser *LogParser) ParseLoop() {
	//TODO: This is currently generating naive test data, will need to change
	for i := 0; i < 10; i++ {
		parser.LogChan <- Log{
			EOF:   false,
			Val:   []byte(strconv.Itoa(i)),
			Topic: "naive_test",
		}
	}
	parser.LogChan <- Log{EOF: true}
}
