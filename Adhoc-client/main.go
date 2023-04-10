package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
	"net/http"
	"os"
)

type AdHocQuery struct {
	Query string   `yaml:"Query"`
	Keys  []string `yaml:"Keys"`
	Topic string   `yaml:"Topic"`
}

type ListenerConfig struct {
	Urls []string `yaml:"Urls"`
}

func main() {

	queryJSON := getQueryJson()
	if queryJSON == nil {
		fmt.Println("query json is nil")
		return
	} else {
		fmt.Printf("query json: %v\n", string(queryJSON))
	}
	listenerConfig := ListenerConfig{}
	loadListenerConfig(&listenerConfig)
	fmt.Printf("listener config: %v\n", listenerConfig)
	//req, _ := http.NewRequest("POST", listenerConfig.Url, bytes.NewBuffer(queryJSON))
	for _, url := range listenerConfig.Urls {
		req, _ := http.NewRequest("POST", url, bytes.NewBuffer(queryJSON))
		req.Header.Set("Content-Typeype", "application/json")
		fmt.Printf("request url: %v\n", req.URL)
		client := &http.Client{}
		res, e := client.Do(req)
		fmt.Printf("res: %v\n", res == nil)
		if res == nil || e != nil {
			fmt.Printf("Error in req: %v\n", e)
			return
		}

		fmt.Println("response Status:", res.Status)
		// Print the body to the stdout
		fmt.Println("response Body:", res.Status)
	}
}

func getQueryJson() []byte {
	fmt.Println("get query json")
	query := AdHocQuery{}

	if data, err := os.ReadFile("config.yaml"); err != nil {
		fmt.Printf("Error reading prom config%v\n", err)
		return nil
	} else {
		fmt.Printf("data: _____________________%v\n", data)
		if err := yaml.Unmarshal(data, &query); err != nil {
			fmt.Printf("Error in yaml unmarshal%v\n", err)
			return nil
		}
	}
	fmt.Println(query)
	queryJSON, err := json.Marshal(&query)
	if err != nil {
		fmt.Printf("Error in json marshal%v\n", err)
		return nil
	}
	fmt.Println(string(queryJSON))
	return queryJSON
}

func loadListenerConfig(config *ListenerConfig) {
	fmt.Println("load listener config")
	if data, err := os.ReadFile("config.yaml"); err != nil {
		fmt.Printf("Error reading prom config%v\n", err)
		return
	} else {
		fmt.Printf("data: _____________________%v\n", string(data))
		if err := yaml.Unmarshal(data, config); err != nil {
			fmt.Printf("Error in yaml unmarshal%v\n", err)
			return
		}
	}
	fmt.Println(*config)
}
