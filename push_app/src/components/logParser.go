package components

//Queries Loki

// LogParser queries Loki logs, structure the log, and puts it into a channel shared with preprocessor
type LogParser struct {
	LogChan chan Log //Channel to communicate with preprocessor goroutine
}

//Parsed log
type Log struct {
	//TODO: Fill this
}

// MakeParser creates LogParser object
func MakeParser() *LogParser {
	parser := LogParser{LogChan: make(chan Log)}
	return &parser
}

//Iteratively parse logs using Grafana Loki
func (parser *LogParser) parseLoop() {
	//TODO: Fill this
}
