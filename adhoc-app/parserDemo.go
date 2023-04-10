package main

/*
// Main file of log pushing adHocOperator

import (
	"fmt"
	"push_app/components"
)

// Run this to test the Data Source and Parser
func Main() {
	config := components.LogConfig{}
	err := config.LoadConfig()

	if err != nil {
		panic(err)
	}

	lokiDataSource := components.MakeLokiDataSource(config)
	logParser := components.MakeParser(10, lokiDataSource.GetLogChan(), config)
	outChan := logParser.GetLogChan()

	go lokiDataSource.Run()
	go logParser.Run()

	for log := range outChan {
		if log.EOF {
			break
		}

		fmt.Printf("Output %s\n", string(log.Val))
	}
}
*/
