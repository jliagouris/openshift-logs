package components

// Preprocesses logs queried by parser

type Preprocessor struct {
	LogChan       chan Log       // Channel shared with logParser, gets queried logs
	DataShareChan chan DataShare // Channel shared with operator
	numProducers  int            // Number of producers, used in producer decision
}

// DataShare is the preprocessed share of data to be sent to secrecy servers
type DataShare struct {
	ProducerIdArr []int           // Which producers does this data share go to?
	Message       ProducerMessage // Message to the producer
	EOF           bool            // If all logs are already processed and the task can end
}

// MakePreprocessor Creates preprocessor object
func MakePreprocessor(numProducers int, LogChan chan Log) *Preprocessor {
	preprocessor := Preprocessor{
		LogChan:       LogChan,
		DataShareChan: make(chan DataShare),
		numProducers:  numProducers,
	}
	return &preprocessor
}

// PreprocessLoop Goroutine that iteratively processes logs passed by parser
func (p *Preprocessor) PreprocessLoop() {
	for log := range p.LogChan {
		dataShares := p.generateDataShares(log)

		// Send data shares to the channel
		//go func() {
		for _, dataShare := range dataShares {
			p.DataShareChan <- *dataShare
		}
		//}()
	}
}

// Generate Data shares from log
func (p *Preprocessor) generateDataShares(log Log) []*DataShare {
	//TODO: Fill this
	return nil
}
