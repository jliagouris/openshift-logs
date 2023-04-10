package configs

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

type PrometheusConf struct {
	Route     string `yaml:"route"`
	Token     string `yaml:"token"`
	Query     string `yaml:"query"`
	NumWorker int    `yaml:"numWorker"`
	ClientId  string
}

func (ps *PrometheusConf) LoadConfig() error {
	if data, err := ioutil.ReadFile("prom_config.yaml"); err != nil {
		fmt.Printf("Error reading prom config%v\n", err)
		return err
	} else {
		//fmt.Printf("data: _____________________%v\n", data)
		if err := yaml.Unmarshal(data, &ps); err != nil {
			return err
		}
	}

	return nil
}
