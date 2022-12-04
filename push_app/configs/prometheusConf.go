package configs

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type PrometheusConf struct {
	Route     string `yaml:"route"`
	Token     string `yaml:"token"`
	Query     string `yaml:"query"`
	NumWorker int    `yaml:"numWorker"`
	ClientId  string
	Interval  string
}

func (ps *PrometheusConf) LoadConfig() error {
	if data, err := ioutil.ReadFile("prom_config.yaml"); err != nil {
		fmt.Printf("Error reading prom config%v\n", err)
		return err
	} else {
		if err := yaml.Unmarshal(data, &ps); err != nil {
			return err
		}
	}

	return nil
}
