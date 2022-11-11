package components

import (
	"encoding/json"
	"fmt"
)

//Queries Loki

// LogParser queries Loki logs, structure the log, and puts it into a channel shared with preprocessor
type LogParser struct {
	ParsedChan chan Log              // Channel to communicate with preprocessor goroutine
	InLogChan  chan PrometheusMetric // Channel to communicate with data source goroutine
	config     LogConfig             // Configuration for the log parser
}

// MakeParser creates LogParser object, this definitely will change
func MakeParser(chanBufSize int, inLogChan chan PrometheusMetric, config LogConfig) *LogParser {
	//TODO: This will change
	parser := LogParser{ParsedChan: make(chan Log, chanBufSize), InLogChan: inLogChan, config: config}
	fmt.Printf("parser channel buffer size: %v\n", cap(parser.ParsedChan))
	return &parser
}

// Main goroutine for LogParser, checks for incoming log from InLogChan and parses it.
func (parser *LogParser) Run() {
	exit := false
	for !exit {
		select {
		case log := <-parser.InLogChan:
			fmt.Printf("parser data: %v\n", log)
			if log.EOF {
				fmt.Println("Parser EOF")
				parser.ParsedChan <- Log{
					EOF: true,
				}
				exit = true
			} else {
				parsedMap, err := parser.ParseLog(log)

				if err == nil {
					//log.Val = processedLog
					parser.ParsedChan <- Log{
						EOF:   false,
						Val:   parsedMap,
						Topic: parser.config.ParserConfig.Topic,
						Key:   log.Key,
					}
				} else {
					fmt.Println(err)
				}
			}
		}
	}
}

// ParseLog takes a log and matches the schema of the log with config	and returns the processed log.
func (parser *LogParser) ParseLog(log PrometheusMetric) (map[string]interface{}, error) {
	/*
		var parsedLog map[string]interface{}
		err := json.Unmarshal(log, &parsedLog)

		if err != nil {
			return nil, err
		}
	*/
	fmt.Printf("parse log PrometheusMetric: %v\n", log)
	var rawMetrics map[string]interface{}
	data, _ := json.Marshal(log.Metric)
	json.Unmarshal(data, &rawMetrics)
	rawMetrics["value"] = log.Values
	fmt.Printf("parse log rawMetrics: %v\n", rawMetrics)
	processedLog := make(map[string]interface{})
	fmt.Printf("parse log processedlog: %v\n", processedLog)
	fmt.Printf("parse log logschema: %v\n", parser.config.ParserConfig.LogSchema)
	for _, schema := range parser.config.ParserConfig.LogSchema {
		if val, ok := rawMetrics[schema.Key]; ok {
			processedLog[schema.Key] = val
		} else {
			if schema.Nullable {
				processedLog[schema.Key] = nil
			} else {
				return nil, fmt.Errorf("Missing required field: %v", schema.Key)
			}
		}
	}

	/*
		processedLogBytes, err := json.Marshal(processedLog)

		if err != nil {
			return nil, err
		}
	*/

	return processedLog, nil
}

func (parser *LogParser) GetLogChan() chan Log {
	return parser.ParsedChan
}

/*
// ParseLoop Iteratively parse logs using Grafana Loki
func (parser *LogParser) ParseLoop() {

		//TODO: This is currently generating naive test data, will need to change
		//for i := 0; i < 10; i++ {
		//	parser.ParsedChan <- Log{
		//		EOF:   false,
		//		Val:   []byte(strconv.Itoa(i)),
		//		Topic: "naive_test",
		//	}
		//}
		parser.ParsedChan <- Log{EOF: true}

	for metric := range parser.InLogChan {
		parsedMap, err := parser.ParseLog(metric)
		if err != nil {
			log.Fatalf("Parse log error in parse loop: %v", err)
		}
		parser.ParsedChan <- parsedMap
	}
	parser.ParsedChan <- nil
}
*/
