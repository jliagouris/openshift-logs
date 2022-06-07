package components

//Queries Loki

// LogParser queries Loki logs, structure the log, and puts it into a channel shared with preprocessor
type LogParser struct {
	LogChan chan Log // Channel to communicate with preprocessor goroutine
}

// Log Parsed log
type Log struct {
	EOF bool
	//TODO: Fill this
}

// MakeParser creates LogParser object
func MakeParser() *LogParser {
	parser := LogParser{LogChan: make(chan Log)}
	return &parser
}

// ParseLoop Iteratively parse logs using Grafana Loki
func (parser *LogParser) ParseLoop() {
	//TODO: Fill this
}
