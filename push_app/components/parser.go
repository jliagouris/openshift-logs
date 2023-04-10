package components

import (
	"encoding/json"
	"fmt"
)

// LogParser queries logs, structure the log, and puts it into a channel shared with preprocessor
type LogParser struct {
	ParsedChan chan Log              // Channel to communicate with preprocessor goroutine
	InLogChan  chan PrometheusMetric // Channel to communicate with data source goroutine
	config     LogConfig             // Configuration for the log parser
}

// MakeParser creates LogParser object, this definitely will change
func MakeParser(chanBufSize int, inLogChan chan PrometheusMetric, config LogConfig) *LogParser {
	parser := LogParser{ParsedChan: make(chan Log, chanBufSize), InLogChan: inLogChan, config: config}
	return &parser
}

// Main goroutine for LogParser, checks for incoming log from InLogChan and parses it.
func (parser *LogParser) Run() {
	for log := range parser.InLogChan {
		if log.EOF {
			parser.ParsedChan <- Log{
				EOF: true,
			}
		} else {
			parsedMap, err := parser.ParseLog(log)

			if err == nil {
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

// ParseLog takes a log and matches the schema of the log with config	and returns the processed log.
func (parser *LogParser) ParseLog(log PrometheusMetric) (map[string]interface{}, error) {
	var rawMetrics map[string]interface{}
	data, _ := json.Marshal(log.Metric)
	json.Unmarshal(data, &rawMetrics)
	rawMetrics["value"] = log.Values
	processedLog := make(map[string]interface{})
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

	return processedLog, nil
}

func (parser *LogParser) GetLogChan() chan Log {
	return parser.ParsedChan
}
